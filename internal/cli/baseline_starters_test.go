// Suite: Baseline Starter Files CLI
// Invariant: preview is read-only and only the exact confirmed preview creates repository-owned files.
// Boundary IN: in-process command dispatch, transcripts, JSON, refusals, help and Project Config loading.
// Boundary OUT: catalog validation and filesystem rollback, owned by internal/baseline.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"roundfix/internal/baseline"
	"roundfix/internal/config"
)

func startersCLI(t *testing.T, repo string, args ...string) (string, string, int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	command := append([]string{"baseline", "starters", "--repo", repo}, args...)
	code := runWithContext(context.Background(), command, &stdout, &stderr, commandEnvironment{homeDir: t.TempDir(), workDir: repo})
	return stdout.String(), stderr.String(), code
}

func startersRepository(t *testing.T) string {
	t.Helper()
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, "README.md", "# Starter Files fixture\n")
	commitBaselinePlanTestRepository(t, repo)
	return repo
}

func startersPreview(t *testing.T, repo string) baseline.StartersPayload {
	t.Helper()
	out, diagnostic, code := startersCLI(t, repo, "--profile", "python", "--format", "json")
	if code != exitUnverified || diagnostic != "" {
		t.Fatalf("preview exit=%d stderr=%q stdout=%q", code, diagnostic, out)
	}
	var payload baseline.StartersPayload
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^sha256:[0-9a-f]{64}$`).MatchString(payload.PlanDigest) {
		t.Fatalf("invalid digest: %q", payload.PlanDigest)
	}
	return payload
}

func assertStartersAbsent(t *testing.T, repo string, paths ...string) {
	t.Helper()
	for _, path := range paths {
		if _, err := os.Lstat(filepath.Join(repo, filepath.FromSlash(path))); !os.IsNotExist(err) {
			t.Fatalf("%s exists or cannot be inspected: %v", path, err)
		}
	}
}

func TestBaselineStartersPreviewTranscript(t *testing.T) {
	t.Parallel()
	repo := startersRepository(t)
	payload := startersPreview(t, repo)
	out, diagnostic, code := startersCLI(t, repo, "--profile", "python")
	want := "Baseline starters: python\n- create Makefile\n- create .roundfixrc.yml\n- create .github/workflows/verify.yml\nState: preview\nPlan Digest: " + payload.PlanDigest + "\nNext: roundfix baseline starters --profile python --confirm-plan " + payload.PlanDigest + "\n"
	if out != want || diagnostic != "" || code != exitUnverified {
		t.Fatalf("exit=%d stderr=%q\nstdout=%q\nwant=%q", code, diagnostic, out, want)
	}
	assertStartersAbsent(t, repo, "Makefile", ".roundfixrc.yml", ".github")
}

func TestBaselineStartersApplyTranscript(t *testing.T) {
	t.Parallel()
	repo := startersRepository(t)
	payload := startersPreview(t, repo)
	out, diagnostic, code := startersCLI(t, repo, "--profile", "python", "--confirm-plan", payload.PlanDigest)
	want := "Baseline starters: python\n- create Makefile\n- create .roundfixrc.yml\n- create .github/workflows/verify.yml\nState: applied\nPlan Digest: " + payload.PlanDigest + "\n"
	if out != want || diagnostic != "" || code != exitOK {
		t.Fatalf("exit=%d stderr=%q\nstdout=%q\nwant=%q", code, diagnostic, out, want)
	}
	for _, path := range []string{"Makefile", ".roundfixrc.yml", ".github/workflows/verify.yml"} {
		content, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(path)))
		if err != nil || len(content) == 0 {
			t.Fatalf("written %s: %v (%d bytes)", path, err, len(content))
		}
	}
	out, diagnostic, code = startersCLI(t, repo, "--profile", "python")
	digestLine := regexp.MustCompile(`Plan Digest: sha256:[0-9a-f]{64}\n$`)
	if code != exitOK || diagnostic != "" || !digestLine.MatchString(out) {
		t.Fatalf("current: exit=%d stdout=%q stderr=%q", code, out, diagnostic)
	}
	wantCurrent := "Baseline starters: python\n- keep Makefile\n- keep .roundfixrc.yml\n- keep .github/workflows/verify.yml\nState: current\nPlan Digest: <digest>\n"
	if digestLine.ReplaceAllString(out, "Plan Digest: <digest>\n") != wantCurrent {
		t.Fatalf("current stdout=%q", out)
	}
}

func TestBaselineStartersWritesAProjectConfigRoundfixLoads(t *testing.T) {
	t.Parallel()
	repo := startersRepository(t)
	payload := startersPreview(t, repo)
	out, diagnostic, code := startersCLI(t, repo, "--profile", "python", "--confirm-plan", payload.PlanDigest)
	if code != exitOK || diagnostic != "" {
		t.Fatalf("apply exit=%d stdout=%q stderr=%q", code, out, diagnostic)
	}
	loaded, err := config.Load(config.LoadOptions{HomeDir: t.TempDir(), WorkDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	cfg := loaded.Config
	if cfg.Defaults.Verification != "make verify" || cfg.Worktree.Bootstrap != "uv sync --group dev --frozen" || cfg.Worktree.BootstrapTimeout != 5*time.Minute || !reflect.DeepEqual(cfg.Worktree.Caches, []string{".venv", ".pytest_cache", ".ruff_cache"}) || !reflect.DeepEqual(cfg.Verification.Tools, []string{"uv", "make", "git"}) {
		t.Fatalf("loaded starter values: defaults=%+v worktree=%+v verification=%+v", cfg.Defaults, cfg.Worktree, cfg.Verification)
	}
}

func TestBaselineStartersJSONAndRefusals(t *testing.T) {
	t.Parallel()
	repo := startersRepository(t)
	payload := startersPreview(t, repo)
	if payload.SchemaVersion != baseline.StartersSchemaVersion || payload.Profile != "python" || payload.State != "preview" || len(payload.Starters) != 3 {
		t.Fatalf("payload=%+v", payload)
	}
	for i, path := range []string{"Makefile", ".roundfixrc.yml", ".github/workflows/verify.yml"} {
		entry := payload.Starters[i]
		if entry.Path != path || entry.Action != "create" || !regexp.MustCompile(`^sha256:[0-9a-f]{64}$`).MatchString(entry.ContentIdentity) {
			t.Fatalf("entry=%+v", entry)
		}
	}
	writeBaselinePlanTestFile(t, repo, ".roundfixrc.yml", "# repository owned\n")
	out, diagnostic, code := startersCLI(t, repo, "--profile", "python", "--confirm-plan", payload.PlanDigest, "--format", "json")
	if code != exitUnverified || out != "" || diagnostic != "roundfix baseline starters: starters.plan.stale: the confirmed digest no longer matches; preview again\n" {
		t.Fatalf("stale exit=%d stdout=%q stderr=%q", code, out, diagnostic)
	}
	assertStartersAbsent(t, repo, "Makefile", ".github")
	kept, err := os.ReadFile(filepath.Join(repo, ".roundfixrc.yml"))
	if err != nil || string(kept) != "# repository owned\n" {
		t.Fatalf("kept config=%q err=%v", kept, err)
	}
	for _, args := range [][]string{nil, {"--profile", "unknown-profile"}, {"--profile", "python", "--format", "yaml"}, {"--profile", "python", "extra"}, {"--profile", "python", "--repo", ""}} {
		out, diagnostic, code := startersCLI(t, repo, args...)
		if code != exitPreflight || out != "" || diagnostic == "" {
			t.Fatalf("args=%q exit=%d stdout=%q stderr=%q", args, code, out, diagnostic)
		}
	}
	// A blocked workflow parent must report a write failure and remove earlier writes.
	blocked := startersRepository(t)
	writeBaselinePlanTestFile(t, blocked, ".github", "blocked\n")
	preview := startersPreview(t, blocked)
	out, diagnostic, code = startersCLI(t, blocked, "--profile", "python", "--confirm-plan", preview.PlanDigest)
	if code != exitRunFailed || out != "" || !strings.Contains(diagnostic, "starters.write.failed") {
		t.Fatalf("write failure exit=%d stdout=%q stderr=%q", code, out, diagnostic)
	}
	assertStartersAbsent(t, blocked, "Makefile", ".roundfixrc.yml")
}

func TestBaselineStartersHelp(t *testing.T) {
	t.Parallel()
	out, diagnostic, code := startersCLI(t, t.TempDir(), "--help")
	if code != exitOK || diagnostic != "" || !strings.HasPrefix(out, "Usage:") {
		t.Fatalf("help exit=%d stdout=%q stderr=%q", code, out, diagnostic)
	}
	for _, snippet := range []string{"roundfix baseline starters --profile <id> [--confirm-plan <digest>] [--repo <path>] [--format <text|json>]", "writes only paths that do not exist", "never changes an existing file", "A non-empty preview exits 3", "An empty preview exits 0"} {
		if !strings.Contains(out, snippet) {
			t.Fatalf("help missing %q: %s", snippet, out)
		}
	}
}
