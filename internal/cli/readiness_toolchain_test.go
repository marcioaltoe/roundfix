package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func writeReadinessManifest(t *testing.T, root string) {
	t.Helper()
	path := filepath.Join(root, doctorSetupManifestPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"decisions":{"verification.gate":{"value":"rtk make verify"},"verification.incremental":{"value":"rtk make verify-incremental"}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestToolchainReadinessNamesEachMissingToolAndItsSource(t *testing.T) {
	root := t.TempDir()
	writeReadinessManifest(t, root)
	cfg := roundconfig.Builtin()
	cfg.Defaults.Verification = "MODE=test rtk make verify"
	cfg.Worktree.Bootstrap = "cd repo && make prepare"
	cfg.Delivery.DerivedPaths = []roundconfig.DerivedPathDeclaration{{Regenerate: "rtk make regenerate"}}
	cfg.Verification.Tools = []string{"rtk", "go"}
	deps := scriptedReadiness(t, nil, "rtk")
	var resolved []string
	resolve := deps.resolve
	deps.resolve = func(name string) (string, error) { resolved = append(resolved, name); return resolve(name) }
	line := toolchainReadiness(deps, roundconfig.Loaded{GitRoot: root, Config: cfg})
	if line.Status != CheckStatusFailed || strings.Count(line.Detail, "DR-TOOL-MISSING:") != 1 {
		t.Fatalf("toolchain = %#v", line)
	}
	for _, phrase := range []string{"rtk is not on PATH", "defaults.verification", "delivery.derived_paths[0].regenerate", "Setup Manifest verification.gate", "Setup Manifest verification.incremental", "verification.tools", "DR-TOOL-UNREAD", cfg.Worktree.Bootstrap} {
		if !strings.Contains(line.Detail, phrase) {
			t.Fatalf("missing %q in %s", phrase, line.Detail)
		}
	}
	if strings.Join(resolved, ",") != "go,rtk" {
		t.Fatalf("tool resolutions = %v", resolved)
	}
	for _, command := range []string{"cd repo", "export X=yes", "set -e", "unset X", "test -f file", "[ -f file ]", "if true", "for x", "while true", "until true", "case x", "exec go", "eval go", "source file", ". file", ":", "true", "false", "command go", "builtin go", "./go test", "/bin/go test", "X=y", "go;evil"} {
		t.Run(command, func(t *testing.T) {
			cfg := roundconfig.Builtin()
			cfg.Defaults.Verification = command
			line := toolchainReadiness(deps, roundconfig.Loaded{Config: cfg})
			if line.Status != CheckStatusWarn || !strings.Contains(line.Detail, "DR-TOOL-UNREAD") {
				t.Fatalf("unread = %#v", line)
			}
		})
	}
}

func TestEnvironmentReadinessNamesAMissingPreloadAndNeverAKeyValue(t *testing.T) {
	deps := scriptedReadiness(t, nil, "")
	deps.exists = func(path string) bool { return path == "/existing.cjs" }
	deps.environ = []string{"NODE_OPTIONS=--require /ignored.cjs", `NODE_OPTIONS=--require "/missing path.cjs" -r /existing.cjs --import package`, "ROUNDFIX_OPENROUTER_API_KEY=key-sentinel", "OPENROUTER_API_KEY=generic-sentinel", "TYPESAFE_API_KEY=generic-sentinel", "ROUNDFIX_UNDECLARED_API_KEY=other-sentinel"}
	line := environmentReadiness(deps)
	if line.Status != CheckStatusWarn || strings.Count(line.Detail, "DR-NODE-PRELOAD-MISSING:") != 1 || !strings.Contains(line.Detail, `"/missing path.cjs"`) {
		t.Fatal("missing preload not reported")
	}
	if !strings.Contains(line.Detail, "ROUNDFIX_OPENROUTER_API_KEY set") || !strings.Contains(line.Detail, "ROUNDFIX_TYPESAFE_API_KEY not set") {
		t.Fatal("key presence not reported")
	}
	_, out, stderr := doctorWithReadiness(t, []CheckResult{line})
	for _, forbidden := range []string{"key-sentinel", "generic-sentinel", "other-sentinel", "/ignored.cjs", "/existing.cjs", "ROUNDFIX_UNDECLARED_API_KEY"} {
		if strings.Contains(out+stderr, forbidden) {
			t.Fatal("private or irrelevant environment data rendered")
		}
	}
	if strings.Contains(strings.ReplaceAll(out, "ROUNDFIX_OPENROUTER_API_KEY", ""), "OPENROUTER_API_KEY") {
		t.Fatal("generic key named")
	}
	deps.environ = []string{"NODE_OPTIONS=-r /existing.cjs --require package"}
	if line := environmentReadiness(deps); line.Status != CheckStatusOK {
		t.Fatal("valid preloads warned")
	}
	deps.environ = []string{"NODE_OPTIONS=-r /missing.cjs --require=/missing.cjs --import=/another.mjs", "ROUNDFIX_OPENROUTER_API_KEY=key-sentinel", "ROUNDFIX_OPENROUTER_API_KEY="}
	line = environmentReadiness(deps)
	if strings.Count(line.Detail, "DR-NODE-PRELOAD-MISSING:") != 2 || !strings.Contains(line.Detail, "ROUNDFIX_OPENROUTER_API_KEY not set") {
		t.Fatal("missing paths must be unique and last key entry must win")
	}
}

func TestDoctorPrintsTheFiveReadinessLines(t *testing.T) {
	root := t.TempDir()
	writeReadinessManifest(t, root)
	deps := scriptedReadiness(t, map[string]readinessReply{
		"gh auth status --active --hostname github.com --json hosts": {out: `{"hosts":{}}`},
		"git config --get user.email":                                {code: 1},
	}, "rtk")
	deps.exists = func(string) bool { return false }
	deps.environ = []string{"NODE_OPTIONS=--require /missing.cjs"}
	results := machineReadiness(context.Background(), deps, roundconfig.Loaded{GitRoot: root, Config: roundconfig.Builtin()})
	code, out, stderr := doctorWithReadiness(t, results)
	if code != exitRunFailed || stderr != "" {
		t.Fatalf("Doctor exit=%d stderr=%q", code, stderr)
	}
	previous := strings.Index(out, "pre-pr-review:")
	for _, name := range []string{"gh", "git", "remote", "toolchain", "environment", "skills"} {
		index := strings.Index(out, "\n"+name+":")
		if index <= previous {
			t.Fatalf("line %s out of order", name)
		}
		previous = index
	}
	for _, phrase := range []string{
		"gh: failed (gh " + readinessGHMinimumVersion + "; DR-GH-UNAUTHENTICATED: gh has no account for github.com; next: gh auth login --hostname github.com)",
		"git: failed (" + readinessGitMinimumVersion + " >= 2.23.0; DR-GIT-IDENTITY: user.email is not set for this repository; next: git config user.email <address>)",
		"remote: ok (origin: github.com/owner/repository; reachable)",
		"DR-TOOL-MISSING: rtk is not on PATH, needed by Setup Manifest verification.gate, Setup Manifest verification.incremental; next: install rtk, or change the command that names it)",
		`environment: warn (DR-NODE-PRELOAD-MISSING: NODE_OPTIONS preload "/missing.cjs" does not exist; spec judge keys: ROUNDFIX_OPENROUTER_API_KEY not set, ROUNDFIX_TYPESAFE_API_KEY not set; next: remove the preload from NODE_OPTIONS where your shell sets it)`,
	} {
		if !strings.Contains(out, phrase) {
			t.Fatalf("Doctor missing transcript phrase %q", phrase)
		}
	}
}
