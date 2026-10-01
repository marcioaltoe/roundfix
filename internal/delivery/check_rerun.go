package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type CheckRecovery interface {
	InspectFailedCheck(ctx context.Context, workDir, head string, check PullRequestCheck) (CheckFailure, error)
	RerunFailedCheck(ctx context.Context, workDir string, failure CheckFailure) error
}

type CheckFailure struct {
	RunID         string
	Attempt       int
	Packages      []string
	OutsideChange bool
}

var _ CheckRecovery = GitHubCLI{}

var goFailureSummary = regexp.MustCompile(`(?:^|\s)FAIL\t([^\t\s]+)\t[0-9]+(?:\.[0-9]+)?s(?:\s|$)`)

func checkRunID(link string) string {
	parsed, err := url.Parse(link)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") {
		return ""
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 5 || parts[2] != "actions" || parts[3] != "runs" {
		return ""
	}
	id, err := strconv.ParseUint(parts[4], 10, 64)
	if err != nil || id == 0 {
		return ""
	}
	return parts[4]
}

func (client GitHubCLI) InspectFailedCheck(ctx context.Context, workDir, head string, check PullRequestCheck) (CheckFailure, error) {
	failure := CheckFailure{RunID: checkRunID(check.Link)}
	if failure.RunID == "" {
		return failure, nil
	}
	client.WorkDir = workDir
	attempt, err := client.recoveryCommand(ctx, "gh", "run", "view", failure.RunID, "--json", "attempt")
	if err != nil {
		return failure, err
	}
	var payload struct {
		Attempt int `json:"attempt"`
	}
	if err := json.Unmarshal([]byte(attempt), &payload); err != nil {
		return failure, fmt.Errorf("parse failed check attempt: %w", err)
	}
	failure.Attempt = payload.Attempt
	log, err := client.recoveryCommand(ctx, "gh", "run", "view", failure.RunID, "--log-failed")
	if err != nil {
		return failure, err
	}
	if strings.Contains(log, "[build failed]") || strings.Contains(log, "[setup failed]") {
		return failure, nil
	}
	for _, line := range strings.Split(log, "\n") {
		if match := goFailureSummary.FindStringSubmatch(line); match != nil {
			failure.Packages = append(failure.Packages, match[1])
		}
	}
	slices.Sort(failure.Packages)
	failure.Packages = slices.Compact(failure.Packages)
	if len(failure.Packages) == 0 {
		return failure, nil
	}
	moduleBytes, err := os.ReadFile(filepath.Join(workDir, "go.mod"))
	if err != nil {
		return failure, fmt.Errorf("read failed check module: %w", err)
	}
	module := ""
	for _, line := range strings.Split(string(moduleBytes), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "module" {
			module = strings.Trim(fields[1], `"`)
			break
		}
	}
	if module == "" {
		return failure, nil
	}
	directories := make([]string, 0, len(failure.Packages))
	for _, pkg := range failure.Packages {
		dir := "."
		if pkg != module {
			var ok bool
			dir, ok = strings.CutPrefix(pkg, module+"/")
			if !ok || dir == "" || path.Clean(dir) != dir || strings.HasPrefix(dir, "../") {
				return failure, nil
			}
		}
		directories = append(directories, dir)
	}
	// origin/HEAD is the repository's existing authoritative default-branch reference.
	defaultRef, err := client.recoveryCommand(ctx, "git", "symbolic-ref", "refs/remotes/origin/HEAD")
	if err != nil {
		return failure, err
	}
	branch, ok := strings.CutPrefix(strings.TrimSpace(defaultRef), "refs/remotes/origin/")
	if !ok || branch == "" || strings.HasPrefix(branch, "-") {
		return failure, fmt.Errorf("inspect failed check: invalid remote default branch %q", defaultRef)
	}
	remoteRef := "refs/remotes/origin/" + branch
	if _, err := client.recoveryCommand(ctx, "git", "fetch", "origin", "+refs/heads/"+branch+":"+remoteRef); err != nil {
		return failure, err
	}
	base, err := client.recoveryCommand(ctx, "git", "merge-base", remoteRef, head)
	if err != nil {
		return failure, err
	}
	base = strings.TrimSpace(base)
	if base == "" {
		return failure, fmt.Errorf("inspect failed check: empty merge base")
	}
	// Keep deleted/renamed source paths and unusual filenames attributable.
	changed, err := client.recoveryCommand(ctx, "git", "diff", "--name-only", "--no-renames", "-z", base, head)
	if err != nil {
		return failure, err
	}
	for _, changedPath := range strings.Split(changed, "\x00") {
		if changedPath == "" {
			continue
		}
		for _, dir := range directories {
			if dir == "." || changedPath == dir || strings.HasPrefix(changedPath, dir+"/") {
				return failure, nil
			}
		}
	}
	failure.OutsideChange = true
	return failure, nil
}

func (client GitHubCLI) RerunFailedCheck(ctx context.Context, workDir string, failure CheckFailure) error {
	id, err := strconv.ParseUint(failure.RunID, 10, 64)
	if err != nil || id == 0 || !failure.OutsideChange || failure.Attempt != 1 {
		return fmt.Errorf("re-run failed check: failure is not eligible for one re-run")
	}
	client.WorkDir = workDir
	_, err = client.recoveryCommand(ctx, "gh", "run", "rerun", failure.RunID, "--failed")
	return err
}

func (client GitHubCLI) recoveryCommand(ctx context.Context, name string, args ...string) (string, error) {
	result, err := client.run(ctx, name, args...)
	if err != nil {
		return "", fmt.Errorf("inspect or re-run failed check (%s): %w", name, err)
	}
	if result.ExitCode != 0 {
		return "", commandFailure("inspect or re-run failed check", result)
	}
	return result.Stdout, nil
}
