// Suite: Baseline plan preflight CLI
// Invariant: baseline plan exposes deterministic local preflight evidence and actionable exits without prompts or repository writes.
// Boundary IN: plan dispatch, repo/format flags, text and JSON rendering, diagnostics, and exit categories.
// Boundary OUT: preservation decisions, profile alignment, portable Plan Documents, and apply.

package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestBaselinePlanPreflightJSONActionRequired(t *testing.T) {
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, "AGENTS.md", "root policy\n")
	commitBaselinePlanTestRepository(t, repo)
	writeBaselinePlanTestFile(t, repo, "scratch.txt", "unrelated dirty bytes\n")

	before := baselinePlanTestTree(t, repo)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{
		"baseline", "plan", "--repo", repo, "--format", "json",
	}, &stdout, &stderr)
	if code != exitUnverified {
		t.Fatalf("baseline plan exit = %d, want %d stdout=%q stderr=%q", code, exitUnverified, stdout.String(), stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("action-required plan stderr = %q, want empty", stderr.String())
	}
	var result baselinePlanPreflightResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode baseline plan JSON: %v\n%s", err, stdout.String())
	}
	if result.SchemaVersion != baselineResultSchema || result.Operation != "plan" ||
		result.State != "action_required" || result.Category != "decision" ||
		result.Repository.Digest == "" || result.Snapshot.Digest == "" ||
		result.NextAction == "" {
		t.Fatalf("baseline plan result = %+v", result)
	}
	if strings.Contains(stdout.String(), filepath.ToSlash(repo)) {
		t.Fatalf("portable preflight JSON contains absolute checkout path %q:\n%s", repo, stdout.String())
	}
	for _, preimage := range result.Snapshot.Preimages {
		if preimage.Path == "scratch.txt" {
			t.Fatalf("unrelated dirty path entered snapshot: %+v", result.Snapshot.Preimages)
		}
	}
	after := baselinePlanTestTree(t, repo)
	if before != after {
		t.Fatalf("baseline plan changed repository bytes:\nbefore=%s\nafter=%s", before, after)
	}
}

func TestBaselinePlanPreflightText(t *testing.T) {
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, "nested/CLAUDE.md", "nested policy\n")
	commitBaselinePlanTestRepository(t, repo)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{
		"baseline", "plan", "--repo=" + repo, "--format=text",
	}, &stdout, &stderr)
	if code != exitUnverified || stderr.Len() != 0 {
		t.Fatalf("baseline plan text exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		"Baseline plan preflight: action required",
		"Repository identity: sha256:",
		"Bounded snapshot: sha256:",
		"Instruction carriers: 1",
		"baseline.inventory.nested-carrier-conflict: nested/CLAUDE.md",
		"Next action:",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("baseline plan text missing %q:\n%s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), repo) {
		t.Fatalf("baseline plan text contains checkout path %q:\n%s", repo, stdout.String())
	}
}

func TestBaselinePlanPreflightBlocksUnsafeRepository(t *testing.T) {
	repo := newBaselinePlanTestRepository(t)
	if err := os.Symlink("../outside.md", filepath.Join(repo, "AGENTS.md")); err != nil {
		t.Fatalf("create escaping alias: %v", err)
	}
	commitBaselinePlanTestRepository(t, repo)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{
		"baseline", "plan", "--repo", repo, "--format=json",
	}, &stdout, &stderr)
	if code != exitPreflight {
		t.Fatalf("unsafe baseline plan exit = %d, want %d stdout=%q stderr=%q", code, exitPreflight, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "baseline.inventory.unsafe-alias") {
		t.Fatalf("unsafe baseline plan diagnostic = %q", stderr.String())
	}
	var result baselinePlanPreflightResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode unsafe baseline plan JSON: %v\n%s", err, stdout.String())
	}
	if result.State != "failed" || result.Category != "preflight" ||
		len(result.Snapshot.Blocking) == 0 || result.NextAction == "" {
		t.Fatalf("unsafe baseline plan result = %+v", result)
	}
}

func TestBaselinePlanPreflightRejectsUsageAndUncommittedRepository(t *testing.T) {
	tests := []struct {
		name string
		args func(*testing.T) []string
		want string
	}{
		{
			name: "unknown flag",
			args: func(t *testing.T) []string {
				return []string{"baseline", "plan", "--unknown"}
			},
			want: "invalid baseline plan arguments",
		},
		{
			name: "repository without commit",
			args: func(t *testing.T) []string {
				return []string{"baseline", "plan", "--repo", newBaselinePlanTestRepository(t)}
			},
			want: "requires at least one commit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := RunContext(context.Background(), tt.args(t), &stdout, &stderr)
			if code != exitPreflight || stdout.Len() != 0 || !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("baseline plan exit = %d stdout=%q stderr=%q, want %q", code, stdout.String(), stderr.String(), tt.want)
			}
		})
	}
}

func TestBaselinePlanPreflightHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{"baseline", "plan", "--help"}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("baseline plan help exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		"roundfix baseline plan",
		"--repo",
		"--format",
		"read-only",
		"Exit codes:",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("baseline plan help missing %q:\n%s", want, stdout.String())
		}
	}
	for _, forbidden := range []string{"--yes", "--interactive", "--no-input"} {
		if strings.Contains(stdout.String(), forbidden) {
			t.Fatalf("baseline plan help advertises interactive flag %q:\n%s", forbidden, stdout.String())
		}
	}
}

func TestBaselinePlanPreflightRealCLI(t *testing.T) {
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, "AGENTS.md", "real CLI policy\n")
	commitBaselinePlanTestRepository(t, repo)

	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve project root: %v", err)
	}
	binary := filepath.Join(t.TempDir(), "roundfix")
	build := exec.Command("go", "build", "-buildvcs=false", "-o", binary, "./cmd/roundfix")
	build.Dir = projectRoot
	build.Env = os.Environ()
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build real roundfix CLI: %v\n%s", err, output)
	}

	command := exec.Command(binary, "baseline", "plan", "--repo", repo, "--format=json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err = command.Run()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != exitUnverified {
		t.Fatalf("real CLI error = %v stdout=%q stderr=%q, want exit %d", err, stdout.String(), stderr.String(), exitUnverified)
	}
	if stderr.Len() != 0 {
		t.Fatalf("real CLI stderr = %q, want empty", stderr.String())
	}
	var result baselinePlanPreflightResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode real CLI JSON: %v\n%s", err, stdout.String())
	}
	if result.State != "action_required" || result.Repository.Digest == "" || result.Snapshot.Digest == "" {
		t.Fatalf("real CLI result = %+v", result)
	}
}

func newBaselinePlanTestRepository(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	runBaselinePlanTestCommand(t, repo, "git", "init", "--quiet")
	runBaselinePlanTestCommand(t, repo, "git", "config", "user.name", "Roundfix Test")
	runBaselinePlanTestCommand(t, repo, "git", "config", "user.email", "roundfix@example.test")
	runBaselinePlanTestCommand(t, repo, "git", "config", "commit.gpgsign", "false")
	return repo
}

func commitBaselinePlanTestRepository(t *testing.T, repo string) {
	t.Helper()
	runBaselinePlanTestCommand(t, repo, "git", "add", "--all")
	runBaselinePlanTestCommand(t, repo, "git", "commit", "--quiet", "-m", "seed")
}

func writeBaselinePlanTestFile(t *testing.T, repo, relative, content string) {
	t.Helper()
	target := filepath.Join(repo, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatalf("create %s parent: %v", relative, err)
	}
	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relative, err)
	}
}

func runBaselinePlanTestCommand(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	command := exec.Command(name, args...)
	command.Dir = dir
	command.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_OPTIONAL_LOCKS=0")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("run %s %s: %v\n%s", name, strings.Join(args, " "), err, output)
	}
}

func baselinePlanTestTree(t *testing.T, root string) string {
	t.Helper()
	var paths []string
	if err := filepath.WalkDir(root, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, filePath)
		if err != nil {
			return err
		}
		if relative != "." {
			paths = append(paths, filepath.ToSlash(relative))
		}
		return nil
	}); err != nil {
		t.Fatalf("list repository tree: %v", err)
	}
	sort.Strings(paths)
	digest := sha256.New()
	for _, relative := range paths {
		filePath := filepath.Join(root, filepath.FromSlash(relative))
		info, err := os.Lstat(filePath)
		if err != nil {
			t.Fatalf("inspect %s: %v", relative, err)
		}
		fmt.Fprintf(digest, "%d:%s:%d:", len(relative), relative, uint32(info.Mode()))
		switch {
		case info.Mode()&fs.ModeSymlink != 0:
			target, err := os.Readlink(filePath)
			if err != nil {
				t.Fatalf("read link %s: %v", relative, err)
			}
			fmt.Fprintf(digest, "%d:%s", len(target), target)
		case info.Mode().IsRegular():
			data, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("read %s: %v", relative, err)
			}
			fmt.Fprintf(digest, "%d:", len(data))
			if _, err := digest.Write(data); err != nil {
				t.Fatalf("hash %s: %v", relative, err)
			}
		}
	}
	return fmt.Sprintf("%x", digest.Sum(nil))
}
