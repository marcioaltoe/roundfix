package cli

// Suite: forge readiness.
// Invariant: only definite failures block Doctor; reads never prompt or expose secrets.
// Boundary IN: readiness folds, Doctor output/exit, real exec against local scripts.
// Boundary OUT: live GitHub and delivery publication (delivery's own suite).

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
)

type readinessReply struct {
	out  string
	code int
	err  error
}

func scriptedReadiness(t *testing.T, overrides map[string]readinessReply, missing string) readinessDependencies {
	t.Helper()
	replies := map[string]readinessReply{
		"git remote get-url origin": {out: "git@github.com:owner/repository.git"},
		"gh --version":              {out: "gh version " + readinessGHMinimumVersion + " (date)"},
		"gh auth status --active --hostname github.com --json hosts": {out: `{"hosts":{"github.com":[{"state":"success","login":"operator","active":true}]}}`},
		"gh repo view owner/repository --json viewerPermission":      {out: `{"viewerPermission":"WRITE"}`},
		"git ls-remote origin HEAD":                                  {out: "abcdef HEAD"},
		"git --version":                                              {out: "git version " + readinessGitMinimumVersion},
		"git config --get user.name":                                 {out: "Identity sentinel"},
		"git config --get user.email":                                {out: "identity-sentinel@example.invalid"},
	}
	for key, reply := range overrides {
		replies[key] = reply
	}
	return readinessDependencies{
		timeout: readinessProbeTimeout,
		environ: []string{"GIT_TERMINAL_PROMPT=1", "GH_HOST=unrelated.invalid"},
		resolve: func(name string) (string, error) {
			if name == missing {
				return "", errors.New("not found")
			}
			return "/fake/" + name, nil
		},
		run: func(ctx context.Context, dir string, env []string, name string, args ...string) (string, string, int, error) {
			key := name + " " + strings.Join(args, " ")
			deadline, ok := ctx.Deadline()
			if !ok || time.Until(deadline) > readinessProbeTimeout {
				t.Fatalf("unbounded probe %s", key)
			}
			if strings.Join(env, "\n") == "" || envValue(env, "GIT_TERMINAL_PROMPT") != "0" {
				t.Fatalf("prompt enabled for %s", key)
			}
			if name == "gh" && envValue(env, "GH_HOST") != "github.com" && !strings.Contains(key, "enterprise.example") {
				t.Fatalf("wrong forge for %s", key)
			}
			if strings.Contains(key, "--show-token") {
				t.Fatal("probe requests a token")
			}
			reply, ok := replies[key]
			if !ok {
				t.Fatalf("unexpected probe %s", key)
			}
			return reply.out, "secret stderr must never be rendered", reply.code, reply.err
		},
	}
}

func envValue(env []string, key string) string {
	var value string
	for _, entry := range env {
		if strings.HasPrefix(entry, key+"=") {
			value = strings.TrimPrefix(entry, key+"=")
		}
	}
	return value
}

func readinessLine(t *testing.T, results []CheckResult, name string) CheckResult {
	t.Helper()
	for _, result := range results {
		if result.Name == name {
			return result
		}
	}
	t.Fatalf("missing %s line", name)
	return CheckResult{}
}

func doctorWithReadiness(t *testing.T, results []CheckResult) (int, string, string) {
	t.Helper()
	withCLIWorkspace(t)
	checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: CheckStatusOK}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
	withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: "/fake/repo"}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
		return profileProofResult{}
	})
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.doctor.readiness = func(context.Context, roundconfig.Loaded) []CheckResult { return results }
	})
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestForgeReadinessReportsEachDefiniteFailure(t *testing.T) {
	t.Parallel()
	cases := []struct {
		code, line, missing, key string
		reply                    readinessReply
	}{
		{code: "DR-GH-MISSING", line: HealthCheckGH, missing: "gh"},
		{code: "DR-GH-VERSION", line: HealthCheckGH, key: "gh --version", reply: readinessReply{out: "gh version 2.80.0"}},
		{code: "DR-GH-UNAUTHENTICATED", line: HealthCheckGH, key: "gh auth status --active --hostname github.com --json hosts", reply: readinessReply{out: `{"hosts":{}}`}},
		{code: "DR-GH-TOKEN-REJECTED", line: HealthCheckGH, key: "gh auth status --active --hostname github.com --json hosts", reply: readinessReply{out: `{"hosts":{"github.com":[{"state":"error","error":"HTTP 401: rejected ghp_sentinel; scope-sentinel"}]}}`}},
		{code: "DR-GH-PERMISSION", line: HealthCheckGH, key: "gh repo view owner/repository --json viewerPermission", reply: readinessReply{out: `{"viewerPermission":"READ"}`}},
		{code: "DR-GIT-MISSING", line: HealthCheckGit, missing: "git"},
		{code: "DR-GIT-VERSION", line: HealthCheckGit, key: "git --version", reply: readinessReply{out: "git version 2.22.0"}},
		{code: "DR-GIT-IDENTITY", line: HealthCheckGit, key: "git config --get user.email", reply: readinessReply{code: 1}},
		{code: "DR-REMOTE-MISSING", line: HealthCheckRemote, key: "git remote get-url origin", reply: readinessReply{code: 2}},
		{code: "DR-REMOTE-FORGE", line: HealthCheckRemote, key: "git remote get-url origin", reply: readinessReply{out: "file:///private/repository"}},
	}
	for _, tc := range cases {
		t.Run(tc.code, func(t *testing.T) {
			overrides := map[string]readinessReply{}
			if tc.key != "" {
				overrides[tc.key] = tc.reply
			}
			results := machineReadiness(context.Background(), scriptedReadiness(t, overrides, tc.missing), roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: "/fake/repo"})
			line := readinessLine(t, results, tc.line)
			if line.Status != CheckStatusFailed || !strings.Contains(line.Detail, tc.code+":") || line.NextAction == "" {
				t.Fatalf("missing coded failure: %#v", line)
			}
			code, out, stderr := doctorWithReadiness(t, results)
			if code != exitRunFailed || stderr != "" || !strings.Contains(out, tc.line+": failed (") || !strings.Contains(out, "next: "+line.NextAction) {
				t.Fatalf("Doctor failure contract: exit=%d", code)
			}
			for _, sentinel := range []string{"ghp_sentinel", "scope-sentinel", "identity-sentinel@example.invalid", "secret stderr"} {
				if strings.Contains(out, sentinel) {
					t.Fatal("Doctor leaked private child output")
				}
			}
		})
	}
}

func TestForgeReadinessWarnsWhenTheForgeDoesNotAnswer(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name, key, code, line string
		reply                 readinessReply
	}{
		{"auth deadline", "gh auth status --active --hostname github.com --json hosts", "DR-GH-UNREACHABLE", HealthCheckGH, readinessReply{err: context.DeadlineExceeded}},
		{"auth refused", "gh auth status --active --hostname github.com --json hosts", "DR-GH-UNREACHABLE", HealthCheckGH, readinessReply{code: 1, err: errors.New("connection refused")}},
		{"auth state timeout", "gh auth status --active --hostname github.com --json hosts", "DR-GH-UNREACHABLE", HealthCheckGH, readinessReply{out: `{"hosts":{"github.com":[{"state":"timeout"}]}}`}},
		{"auth non401 error", "gh auth status --active --hostname github.com --json hosts", "DR-GH-UNREACHABLE", HealthCheckGH, readinessReply{out: `{"hosts":{"github.com":[{"state":"error","error":"HTTP 503; ghp_sentinel"}]}}`}},
		{"auth malformed", "gh auth status --active --hostname github.com --json hosts", "DR-GH-UNREACHABLE", HealthCheckGH, readinessReply{out: `invalid JSON`}},
		{"permission deadline", "gh repo view owner/repository --json viewerPermission", "DR-GH-PERMISSION-UNVERIFIED", HealthCheckGH, readinessReply{err: context.DeadlineExceeded}},
		{"permission refused", "gh repo view owner/repository --json viewerPermission", "DR-GH-PERMISSION-UNVERIFIED", HealthCheckGH, readinessReply{code: 1}},
		{"remote deadline", "git ls-remote origin HEAD", "DR-REMOTE-UNREACHABLE", HealthCheckRemote, readinessReply{err: context.DeadlineExceeded}},
		{"remote refused", "git ls-remote origin HEAD", "DR-REMOTE-UNREACHABLE", HealthCheckRemote, readinessReply{code: 1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			results := machineReadiness(context.Background(), scriptedReadiness(t, map[string]readinessReply{tc.key: tc.reply}, ""), roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: "/fake/repo"})
			line := readinessLine(t, results, tc.line)
			if line.Status != CheckStatusWarn || !strings.Contains(line.Detail, tc.code+":") || line.NextAction != "re-run roundfix doctor when github.com is reachable" {
				t.Fatalf("missing warning: %#v", line)
			}
			code, out, stderr := doctorWithReadiness(t, results)
			if code != exitOK || stderr != "" || !strings.Contains(out, "; next: "+line.NextAction+")") {
				t.Fatalf("warning changed Doctor exit/action: %d", code)
			}
		})
	}
}

func TestGitReadinessReportsVersionAndIdentity(t *testing.T) {
	t.Parallel()
	for _, key := range []string{"user.name", "user.email"} {
		t.Run(key, func(t *testing.T) {
			result := gitReadiness(context.Background(), scriptedReadiness(t, map[string]readinessReply{"git config --get " + key: {code: 1}}, ""), roundconfig.Loaded{GitRoot: "/fake/repo"})
			if result.Status != CheckStatusFailed || !strings.Contains(result.Detail, key+" is not set") || strings.Contains(result.Detail, "sentinel") {
				t.Fatalf("incorrect identity finding: %#v", result)
			}
		})
	}
	t.Run("newer versions and identity", func(t *testing.T) {
		result := gitReadiness(context.Background(), scriptedReadiness(t, map[string]readinessReply{"git --version": {out: "git version 2.50.1 (Apple Git-155)"}}, ""), roundconfig.Loaded{GitRoot: "/fake/repo"})
		if result.Status != CheckStatusOK || strings.Contains(result.Detail, "sentinel") {
			t.Fatalf("unexpected ready Git: %#v", result)
		}
	})
}

func TestDoctorPrintsForgeWarningsAndExitsZero(t *testing.T) {
	t.Parallel()
	deps := scriptedReadiness(t, map[string]readinessReply{
		"gh auth status --active --hostname github.com --json hosts": {err: context.DeadlineExceeded},
		"git ls-remote origin HEAD":                                  {code: 1},
	}, "")
	results := machineReadiness(context.Background(), deps, roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: "/fake/repo"})
	code, out, stderr := doctorWithReadiness(t, results)
	if code != 0 || stderr != "" {
		t.Fatalf("offline Doctor exit=%d", code)
	}
	for _, text := range []string{"gh: warn (gh " + readinessGHMinimumVersion + "; DR-GH-UNREACHABLE:", "remote: warn (origin: github.com/owner/repository; DR-REMOTE-UNREACHABLE:", "next: re-run roundfix doctor when github.com is reachable"} {
		if !strings.Contains(out, text) {
			t.Fatalf("missing transcript text %q", text)
		}
	}
	if strings.Index(out, "pre-pr-review:") > strings.Index(out, "gh:") || strings.Index(out, "remote:") > strings.Index(out, "skills:") {
		t.Fatal("readiness lines out of order")
	}
}

func TestReadinessFindingsFoldAndPrintNextActions(t *testing.T) {
	t.Parallel()
	for _, statuses := range [][]CheckStatus{{}, {CheckStatusWarn}, {CheckStatusWarn, CheckStatusFailed}, {CheckStatusFailed, CheckStatusWarn}} {
		var findings []readinessFinding
		want := CheckStatusOK
		for i, status := range statuses {
			findings = append(findings, readinessFinding{fmt.Sprintf("DR-%d", i), status, "finding", "action"})
			if status == CheckStatusFailed || want == CheckStatusOK {
				want = status
			}
		}
		result := readinessResult("example", "ready detail", findings)
		if result.Status != want {
			t.Fatalf("fold=%s want=%s", result.Status, want)
		}
		var out bytes.Buffer
		printDoctorResult(&out, result)
		if len(statuses) > 0 && !strings.Contains(out.String(), "next: "+strings.TrimSuffix(strings.Repeat("action && ", len(statuses)), " && ")) {
			t.Fatal("actions were not joined and printed")
		}
	}
}

func TestReadinessRemoteFormsAndDeliverySelection(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{"https://github.com/owner/repository.git", "ssh://git@github.com/owner/repository.git", "git@github.com:owner/repository.git"} {
		t.Run(raw, func(t *testing.T) {
			deps := scriptedReadiness(t, map[string]readinessReply{"git remote get-url publish": {out: raw}, "git ls-remote publish HEAD": {out: "head"}}, "")
			cfg := roundconfig.Builtin()
			cfg.Watch.PushRemote = " publish "
			results := forgeReadiness(context.Background(), deps, roundconfig.Loaded{Config: cfg, GitRoot: "/fake/repo"})
			if results[0].Status != CheckStatusOK || results[1].Status != CheckStatusOK || !strings.Contains(results[1].Detail, "publish: github.com/owner/repository") {
				t.Fatalf("delivery remote not resolved: %#v", results)
			}
		})
	}
	for _, raw := range []string{"/local/path", "file:///repo", "https://github.com/owner/repo/extra", "git@host:../repo"} {
		if _, _, ok := parseReadinessRemote(raw); ok {
			t.Fatalf("accepted non-forge %q", raw)
		}
	}
	t.Run("enterprise", func(t *testing.T) {
		deps := scriptedReadiness(t, map[string]readinessReply{
			"git remote get-url origin": {out: "ssh://git@enterprise.example/owner/repository.git"},
			"gh auth status --active --hostname enterprise.example --json hosts": {out: `{"hosts":{"enterprise.example":[{"state":"success","login":"operator"}]}}`},
		}, "")
		deps.run = enterpriseReadinessRunner(t, deps.run)
		results := forgeReadiness(context.Background(), deps, roundconfig.Loaded{GitRoot: "/fake/repo"})
		if results[0].Status != CheckStatusOK || results[1].Status != CheckStatusOK {
			t.Fatalf("known enterprise host rejected: %#v", results)
		}
	})
	t.Run("unknown host", func(t *testing.T) {
		deps := scriptedReadiness(t, map[string]readinessReply{"git remote get-url origin": {out: "git@other.example:owner/repository.git"}, "gh auth status --active --hostname other.example --json hosts": {out: `{"hosts":{}}`}}, "")
		deps.run = enterpriseReadinessRunner(t, deps.run)
		results := forgeReadiness(context.Background(), deps, roundconfig.Loaded{GitRoot: "/fake/repo"})
		if results[1].Status != CheckStatusFailed || !strings.Contains(results[1].Detail, "DR-REMOTE-FORGE") {
			t.Fatal("unknown host accepted")
		}
	})
	t.Run("outside Git", func(t *testing.T) {
		results := forgeReadiness(context.Background(), scriptedReadiness(t, nil, ""), roundconfig.Loaded{})
		if results[1].Status != CheckStatusSkipped || results[1].Detail != "requires a Git repository" {
			t.Fatal("remote must skip outside Git")
		}
	})
}

// Check the host at the exec boundary before adapting only the scripted fixture's
// environment guard, whose baseline is github.com.
func enterpriseReadinessRunner(t *testing.T, run readinessRunner) readinessRunner {
	return func(ctx context.Context, dir string, env []string, name string, args ...string) (string, string, int, error) {
		if name == "gh" {
			if host := envValue(env, "GH_HOST"); host != "enterprise.example" && host != "other.example" {
				t.Fatal("enterprise read targeted another forge")
			}
		}
		return run(ctx, dir, readinessEnvironment(env, "GH_HOST", "github.com"), name, args...)
	}
}

func TestForgeProbesAreBoundedAndNeverPromptOrPrintAToken(t *testing.T) {
	// Sequential: the probe's 250 ms bound is the behavior under test, and the child missed it under the race detector while it ran in parallel (0254 QA F4).
	dir := t.TempDir()
	setCommandEnvForTest(t, "PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	setCommandEnvForTest(t, "PROBE_LOG", filepath.Join(dir, "probes"))
	script := `#!/bin/sh
{
 printf 'argv:%s\n' "$*"
 printf 'prompt:%s\n' "$GIT_TERMINAL_PROMPT"
 printf 'host:%s\n' "$GH_HOST"
 input=$(cat)
 printf 'stdin:%s\n' "$input"
} >> "$PROBE_LOG"
case "$*" in
 '--version') case "$0" in *gh.fixture.sh) printf 'gh version %s\n' "$GH_TEST_VERSION";; *) printf 'git version %s\n' "$GIT_TEST_VERSION";; esac;;
 'remote get-url origin') printf 'git@github.com:owner/repository.git\n';;
 'auth status --active --hostname github.com --json hosts') printf '%s\n' '{"hosts":{"github.com":[{"state":"success","login":"operator","token":"token-sentinel","scopes":"scopes-sentinel"}]}}';;
 'repo view owner/repository --json viewerPermission') printf '%s\n' '{"viewerPermission":"WRITE"}';;
 'config --get user.name') printf 'identity-name-sentinel\n';;
 'config --get user.email') printf 'identity-address-sentinel@example.invalid\n';;
 'ls-remote origin HEAD') printf 'abcdef HEAD\n';;
 *) exit 90;;
esac
`
	setCommandEnvForTest(t, "GH_TEST_VERSION", readinessGHMinimumVersion)
	setCommandEnvForTest(t, "GIT_TEST_VERSION", readinessGitMinimumVersion)
	for _, name := range []string{"gh", "git"} {
		writeScriptFixture(t, filepath.Join(dir, name), script)
	}
	deps := readinessDependenciesForCommand(commandEnvironmentForTest(t).dependencies)
	deps.run = func(ctx context.Context, workDir string, env []string, name string, args ...string) (string, string, int, error) {
		return execReadinessRunner(ctx, workDir, env, filepath.Join(dir, name), args...)
	}
	results := machineReadiness(context.Background(), deps, roundconfig.Loaded{GitRoot: dir})
	var out bytes.Buffer
	for _, result := range results {
		if result.Status != CheckStatusOK {
			t.Fatalf("fake executable failed: %s", result.Name)
		}
		printDoctorResult(&out, result)
	}
	log, err := os.ReadFile(filepath.Join(dir, "probes"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(log), "argv:") != 8 || strings.Count(string(log), "prompt:0\n") != 8 || strings.Count(string(log), "stdin:\n") != 8 || strings.Contains(string(log), "--show-token") {
		t.Fatalf("unsafe probe transcript: %s", log)
	}
	for _, value := range []string{"token-sentinel", "scopes-sentinel", "identity-name-sentinel", "identity-address-sentinel@example.invalid"} {
		if strings.Contains(out.String(), value) {
			t.Fatal("secret appeared in readiness output")
		}
	}
	// exec replaces the shell: the process being cancelled is the sleeping child.
	sleeper := `#!/bin/sh
printf '%s\n' "$$" > "$SLEEP_PID"
exec sleep 30
`
	setCommandEnvForTest(t, "SLEEP_PID", filepath.Join(dir, "pid"))
	if err := os.WriteFile(filepath.Join(dir, "gh")+scriptFixtureSuffix, []byte(sleeper), 0o600); err != nil {
		t.Fatal(err)
	}
	deps.environ = commandEnvironmentForTest(t).dependencies.environ()
	deps.timeout = 250 * time.Millisecond
	// Under load the shell can need longer than the bound to record its PID, so
	// each attempt must still time out and cancel, and the liveness check below
	// runs on the first attempt whose child started.
	var pid []byte
	for attempt := 0; attempt < 5 && pid == nil; attempt++ {
		if err := os.Remove(filepath.Join(dir, "pid")); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
		started := time.Now()
		_, _, err = deps.probe(context.Background(), dir, "github.com", "gh", "--version")
		if !errors.Is(err, context.DeadlineExceeded) || time.Since(started) > 3*time.Second {
			t.Fatalf("sleeping child did not cancel: %v", err)
		}
		if recorded, readErr := os.ReadFile(filepath.Join(dir, "pid")); readErr == nil && len(strings.TrimSpace(string(recorded))) > 0 {
			pid = recorded
		}
	}
	if pid == nil {
		t.Fatal("sleeping child never started")
	}
	// kill -0 is a read-only liveness probe, through the same bounded exec runner.
	_, _, code, err := execReadinessRunner(context.Background(), dir, os.Environ(), "/bin/kill", "-0", strings.TrimSpace(string(pid)))
	if code == 0 || err == nil {
		t.Fatal("cancelled child is still alive")
	}
}
