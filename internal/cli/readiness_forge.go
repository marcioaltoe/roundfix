package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"roundfix/internal/baseline"
	roundconfig "roundfix/internal/config"
)

const (
	readinessGHMinimumVersion  = "2.81.0"
	readinessGitMinimumVersion = "2.23.0"
	readinessProbeTimeout      = 10 * time.Second
)

func defaultReadinessDependencies() readinessDependencies {
	return readinessDependencies{
		run: execReadinessRunner, resolve: func(name string) (string, error) {
			path, reason := baseline.ResolveExecutable(name, filepath.SplitList(os.Getenv("PATH")))
			if reason != "" {
				return "", errors.New(reason)
			}
			return path, nil
		}, environ: os.Environ(),
		exists:  func(path string) bool { _, err := os.Stat(path); return err == nil },
		timeout: readinessProbeTimeout,
	}
}

// execReadinessRunner never attaches stdin or forwards child output to a user.
func execReadinessRunner(ctx context.Context, dir string, env []string, name string, args ...string) (string, string, int, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = env
	// Bound pipe draining too, if a child leaves an inherited pipe open.
	cmd.WaitDelay = time.Second
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		code = -1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			code = exitErr.ExitCode()
		}
	}
	if ctx.Err() != nil {
		err = ctx.Err()
	}
	return stdout.String(), stderr.String(), code, err
}

func readinessEnvironment(env []string, key, value string) []string {
	result := make([]string, 0, len(env)+1)
	for _, entry := range env {
		if !strings.HasPrefix(entry, key+"=") {
			result = append(result, entry)
		}
	}
	return append(result, key+"="+value)
}

func (deps readinessDependencies) probe(ctx context.Context, dir, host, name string, args ...string) (string, int, error) {
	timeout := deps.timeout
	if timeout <= 0 {
		timeout = readinessProbeTimeout
	}
	bounded, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	env := readinessEnvironment(deps.environ, "GIT_TERMINAL_PROMPT", "0")
	if host != "" {
		env = readinessEnvironment(env, "GH_HOST", host)
	}
	stdout, _, code, err := deps.run(bounded, dir, env, name, args...)
	if bounded.Err() != nil {
		err = bounded.Err()
	}
	return stdout, code, err
}

func machineReadiness(ctx context.Context, deps readinessDependencies, loaded roundconfig.Loaded) []CheckResult {
	forge := forgeReadiness(ctx, deps, loaded)
	return []CheckResult{forge[0], gitReadiness(ctx, deps, loaded), forge[1], toolchainReadiness(deps, loaded), environmentReadiness(deps)}
}

type readinessRemote struct {
	name, host, repository string
	findings               []readinessFinding
}

var readinessHostPattern = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9.-]*[A-Za-z0-9])?$`)
var readinessRepositoryPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$`)
var readinessVersionPattern = regexp.MustCompile(`\b[0-9]+\.[0-9]+\.[0-9]+\b`)
var readinessHTTP401Pattern = regexp.MustCompile(`(?i)\bHTTP(?:\s+status)?\s*[:=]?\s*401\b`)
var readinessLoginPattern = regexp.MustCompile(`^[A-Za-z0-9-]+$`)

// Keep credentials, query strings and remote userinfo out of all rendered text.
func parseReadinessRemote(raw string) (string, string, bool) {
	raw = strings.TrimSpace(raw)
	var host, path string
	if strings.Contains(raw, "://") {
		parsed, err := url.Parse(raw)
		if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "ssh") || parsed.RawQuery != "" || parsed.Fragment != "" {
			return "", "", false
		}
		host, path = parsed.Hostname(), strings.TrimPrefix(parsed.Path, "/")
	} else {
		userHost, remotePath, ok := strings.Cut(raw, ":")
		_, remoteHost, hasUser := strings.Cut(userHost, "@")
		if !ok || !hasUser {
			return "", "", false
		}
		host, path = remoteHost, remotePath
	}
	host = strings.ToLower(host)
	path = strings.TrimSuffix(path, ".git")
	if !readinessHostPattern.MatchString(host) || !readinessRepositoryPattern.MatchString(path) {
		return "", "", false
	}
	for _, part := range strings.Split(path, "/") {
		if part == "." || part == ".." || strings.HasPrefix(part, "-") {
			return "", "", false
		}
	}
	return host, path, true
}

func readDeliveryRemote(ctx context.Context, deps readinessDependencies, loaded roundconfig.Loaded) readinessRemote {
	remote := readinessRemote{name: strings.TrimSpace(loaded.Config.Watch.PushRemote)}
	if remote.name == "" {
		remote.name = "origin"
	}
	if loaded.GitRoot == "" {
		return remote
	}
	if _, err := deps.resolve("git"); err != nil {
		return remote
	}
	output, code, err := deps.probe(ctx, loaded.GitRoot, "", "git", "remote", "get-url", remote.name)
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		remote.findings = append(remote.findings, readinessFinding{"DR-REMOTE-UNREACHABLE", CheckStatusWarn,
			"could not read the delivery remote: " + readinessReadReason(err), "re-run roundfix doctor when github.com is reachable"})
		return remote
	}
	if code != 0 || err != nil || strings.TrimSpace(output) == "" {
		remote.findings = append(remote.findings, readinessFinding{"DR-REMOTE-MISSING", CheckStatusFailed,
			remote.name + " is not configured", "git remote add " + remote.name + " <url>"})
		return remote
	}
	var ok bool
	remote.host, remote.repository, ok = parseReadinessRemote(output)
	if !ok {
		remote.findings = append(remote.findings, remote.forgeFinding())
	}
	return remote
}

func (remote readinessRemote) forgeFinding() readinessFinding {
	return readinessFinding{"DR-REMOTE-FORGE", CheckStatusFailed, remote.name + " is not a GitHub forge remote",
		"point " + remote.name + " at the repository's GitHub URL, or set watch.push_remote"}
}

type readinessAuthEntry struct {
	State string `json:"state"`
	Error string `json:"error"`
	Login string `json:"login"`
}

func forgeReadiness(ctx context.Context, deps readinessDependencies, loaded roundconfig.Loaded) []CheckResult {
	remote := readDeliveryRemote(ctx, deps, loaded)
	host := remote.host
	if host == "" {
		host = "github.com"
	}
	gh, knownHost := ghReadiness(ctx, deps, loaded.GitRoot, host, remote.repository)
	remoteResult := CheckResult{Name: HealthCheckRemote, Status: CheckStatusSkipped, Detail: "requires a Git repository"}
	if loaded.GitRoot != "" {
		if _, err := deps.resolve("git"); err != nil {
			remoteResult.Detail = "requires Git"
		} else {
			detail := remote.name
			if remote.host != "" {
				detail += ": " + remote.host + "/" + remote.repository
				if remote.host != "github.com" && !knownHost {
					// An unanswered auth read cannot prove an enterprise host is invalid.
					if gh.Status == CheckStatusWarn {
						remote.findings = append(remote.findings, readinessFinding{"DR-REMOTE-UNREACHABLE", CheckStatusWarn,
							"could not confirm the forge host", "re-run roundfix doctor when " + host + " is reachable"})
					} else {
						remote.findings = append(remote.findings, remote.forgeFinding())
					}
				} else {
					_, code, err := deps.probe(ctx, loaded.GitRoot, "", "git", "ls-remote", remote.name, "HEAD")
					if err != nil || code != 0 {
						remote.findings = append(remote.findings, readinessFinding{"DR-REMOTE-UNREACHABLE", CheckStatusWarn,
							"git ls-remote " + remote.name + " failed: " + readinessReadReason(err), "re-run roundfix doctor when " + host + " is reachable"})
					} else {
						detail += "; reachable"
					}
				}
			}
			remoteResult = readinessResult(HealthCheckRemote, detail, remote.findings)
		}
	}
	return []CheckResult{gh, remoteResult}
}

func readinessReadReason(err error) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return "read timed out"
	}
	if errors.Is(err, context.Canceled) {
		return "read cancelled"
	}
	return "read did not answer"
}

func ghReadiness(ctx context.Context, deps readinessDependencies, dir, host, repository string) (CheckResult, bool) {
	finish := func(detail string, finding readinessFinding, known bool) (CheckResult, bool) {
		return readinessResult(HealthCheckGH, detail, []readinessFinding{finding}), known
	}
	if _, err := deps.resolve("gh"); err != nil {
		return finish("", readinessFinding{"DR-GH-MISSING", CheckStatusFailed, "gh is not on PATH",
			"install GitHub CLI " + readinessGHMinimumVersion + " or newer from https://cli.github.com"}, false)
	}
	output, code, err := deps.probe(ctx, dir, host, "gh", "--version")
	version := readinessVersionPattern.FindString(output)
	detail := "gh " + version
	if err != nil || code != 0 || compareVersions(version, readinessGHMinimumVersion) < 0 {
		return finish(detail, readinessFinding{"DR-GH-VERSION", CheckStatusFailed, "GitHub CLI " + readinessGHMinimumVersion + " or newer is required",
			"upgrade GitHub CLI to " + readinessGHMinimumVersion + " or newer"}, false)
	}
	retry := "re-run roundfix doctor when " + host + " is reachable"
	output, code, err = deps.probe(ctx, dir, host, "gh", "auth", "status", "--active", "--hostname", host, "--json", "hosts")
	var auth struct {
		Hosts map[string][]readinessAuthEntry `json:"hosts"`
	}
	if err != nil || code != 0 || json.Unmarshal([]byte(output), &auth) != nil || auth.Hosts == nil {
		return finish(detail, readinessFinding{"DR-GH-UNREACHABLE", CheckStatusWarn,
			"could not confirm the " + host + " login: " + readinessReadReason(err), retry}, false)
	}
	entries := auth.Hosts[host]
	if len(entries) == 0 {
		return finish(detail, readinessFinding{"DR-GH-UNAUTHENTICATED", CheckStatusFailed, "gh has no account for " + host,
			"gh auth login --hostname " + host}, false)
	}
	entry := entries[0] // --active restricts the response to the active account.
	if entry.State == "error" && readinessHTTP401Pattern.MatchString(entry.Error) {
		return finish(detail, readinessFinding{"DR-GH-TOKEN-REJECTED", CheckStatusFailed, host + " rejected the login (HTTP 401)",
			"gh auth refresh --hostname " + host}, true)
	}
	if entry.State != "success" {
		return finish(detail, readinessFinding{"DR-GH-UNREACHABLE", CheckStatusWarn, "could not confirm the " + host + " login: read did not answer", retry}, true)
	}
	// GitHub login identifiers cannot contain addresses or arbitrary child text.
	if readinessLoginPattern.MatchString(entry.Login) {
		detail += "; " + host + " login=" + entry.Login
	} else {
		detail += "; " + host + " authenticated"
	}
	if repository == "" {
		return readinessResult(HealthCheckGH, detail, nil), true
	}
	output, code, err = deps.probe(ctx, dir, host, "gh", "repo", "view", repository, "--json", "viewerPermission")
	var permission struct {
		ViewerPermission string `json:"viewerPermission"`
	}
	if err != nil || code != 0 || json.Unmarshal([]byte(output), &permission) != nil || permission.ViewerPermission == "" {
		return finish(detail, readinessFinding{"DR-GH-PERMISSION-UNVERIFIED", CheckStatusWarn, "could not confirm write access to " + repository, retry}, true)
	}
	switch permission.ViewerPermission {
	case "ADMIN", "MAINTAIN", "WRITE":
		return readinessResult(HealthCheckGH, detail+"; "+repository+" permission="+permission.ViewerPermission, nil), true
	default:
		return finish(detail, readinessFinding{"DR-GH-PERMISSION", CheckStatusFailed, "account lacks write access to " + repository,
			"ask for write access to " + repository + ", or gh auth switch --hostname " + host}, true)
	}
}

func gitReadiness(ctx context.Context, deps readinessDependencies, loaded roundconfig.Loaded) CheckResult {
	if _, err := deps.resolve("git"); err != nil {
		return readinessResult(HealthCheckGit, "", []readinessFinding{{"DR-GIT-MISSING", CheckStatusFailed,
			"git is not on PATH", "install Git " + readinessGitMinimumVersion + " or newer"}})
	}
	output, code, err := deps.probe(ctx, loaded.GitRoot, "", "git", "--version")
	version := readinessVersionPattern.FindString(output)
	var findings []readinessFinding
	detail := version
	if err != nil || code != 0 || compareVersions(version, readinessGitMinimumVersion) < 0 {
		findings = append(findings, readinessFinding{"DR-GIT-VERSION", CheckStatusFailed, "Git " + readinessGitMinimumVersion + " or newer is required",
			"upgrade Git to " + readinessGitMinimumVersion + " or newer"})
	} else {
		detail = fmt.Sprintf("%s >= %s", version, readinessGitMinimumVersion)
	}
	if loaded.GitRoot != "" {
		for _, key := range []string{"user.name", "user.email"} {
			output, code, err := deps.probe(ctx, loaded.GitRoot, "", "git", "config", "--get", key)
			if strings.TrimSpace(output) == "" || code != 0 || err != nil {
				placeholder := "<name>"
				if key == "user.email" {
					placeholder = "<address>"
				}
				text := key + " is not set for this repository"
				if (err != nil || code != 0) && code != 1 {
					text = key + " could not be read for this repository"
				}
				findings = append(findings, readinessFinding{"DR-GIT-IDENTITY", CheckStatusFailed, text,
					"git config " + key + " " + placeholder})
			}
		}
	}
	return readinessResult(HealthCheckGit, detail, findings)
}
