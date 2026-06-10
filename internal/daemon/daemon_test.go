package daemon

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecVerifierRunsConfiguredCommand(t *testing.T) {
	var stream bytes.Buffer
	err := ExecVerifier{}.Verify(context.Background(), VerifyRequest{
		WorkDir: t.TempDir(),
		Command: "printf verified",
		Stream:  &stream,
	})

	if err != nil {
		t.Fatalf("verify command: %v", err)
	}
	if stream.String() != "verified" {
		t.Fatalf("expected verification output, got %q", stream.String())
	}
}

func TestExecVerifierReportsFailedCommand(t *testing.T) {
	var stream bytes.Buffer
	err := ExecVerifier{}.Verify(context.Background(), VerifyRequest{
		WorkDir: t.TempDir(),
		Command: "printf broken; exit 7",
		Stream:  &stream,
	})

	if err == nil {
		t.Fatal("expected verification failure")
	}
	if !strings.Contains(err.Error(), "verification command") {
		t.Fatalf("expected verification error context, got %v", err)
	}
	if stream.String() != "broken" {
		t.Fatalf("expected failed verification output, got %q", stream.String())
	}
}

func TestGitCommitterValidatesRequest(t *testing.T) {
	err := GitCommitter{}.Commit(context.Background(), CommitRequest{})

	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "git root is required") {
		t.Fatalf("expected git root validation error, got %v", err)
	}
}

func TestGitCommitterExcludesProjectConfigFromBatchCommit(t *testing.T) {
	repoDir := t.TempDir()
	runGitForTest(t, repoDir, "init", "-q")
	runGitForTest(t, repoDir, "config", "user.email", "test@example.com")
	runGitForTest(t, repoDir, "config", "user.name", "Test")
	mustWriteForTest(t, filepath.Join(repoDir, "tracked.txt"), "base\n")
	runGitForTest(t, repoDir, "add", "tracked.txt")
	runGitForTest(t, repoDir, "commit", "-q", "-m", "init")

	mustWriteForTest(t, filepath.Join(repoDir, "tracked.txt"), "changed\n")
	mustWriteForTest(t, filepath.Join(repoDir, "created.txt"), "created\n")
	mustWriteForTest(t, filepath.Join(repoDir, ".roundfixrc.yml"), "defaults:\n  agent: codex\n")

	err := GitCommitter{}.Commit(context.Background(), CommitRequest{
		WorkDir:      repoDir,
		Message:      "fix: test batch",
		ExcludePaths: []string{".roundfixrc.yml"},
	})
	if err != nil {
		t.Fatalf("expected Batch commit, got %v", err)
	}

	committed := runGitForTest(t, repoDir, "show", "--name-only", "--format=", "HEAD")
	if !strings.Contains(committed, "tracked.txt") || !strings.Contains(committed, "created.txt") {
		t.Fatalf("expected code changes in commit, got %q", committed)
	}
	if strings.Contains(committed, ".roundfixrc.yml") {
		t.Fatalf("did not expect Project Config in commit, got %q", committed)
	}
	status := runGitForTest(t, repoDir, "status", "--porcelain=v1")
	if !strings.Contains(status, "?? .roundfixrc.yml") {
		t.Fatalf("expected Project Config to remain untracked, got %q", status)
	}
}

func TestGitPusherValidatesRequest(t *testing.T) {
	err := GitPusher{}.Push(context.Background(), PushRequest{})

	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "git root is required") {
		t.Fatalf("expected git root validation error, got %v", err)
	}
}

func runGitForTest(t *testing.T, workDir string, args ...string) string {
	t.Helper()
	cmdArgs := append([]string{"-C", workDir}, args...)
	cmd := exec.Command("git", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, string(output))
	}
	return string(output)
}

func mustWriteForTest(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestSnapshotDiffCommitStagesOnlyAgentChangesInRealRepo(t *testing.T) {
	ctx := context.Background()
	repoDir := t.TempDir()
	runGitForTest(t, repoDir, "init", "-b", "main")
	runGitForTest(t, repoDir, "config", "user.name", "Roundfix Test")
	runGitForTest(t, repoDir, "config", "user.email", "test@example.com")
	mustWriteForTest(t, filepath.Join(repoDir, "tracked.txt"), "original\n")
	runGitForTest(t, repoDir, "add", "tracked.txt")
	runGitForTest(t, repoDir, "commit", "-m", "initial")

	// Pre-existing user work that slipped past Preflight Validation.
	mustWriteForTest(t, filepath.Join(repoDir, "tracked.txt"), "user edit\n")
	mustWriteForTest(t, filepath.Join(repoDir, "user-notes.md"), "wip\n")

	snapshotter := GitWorktreeSnapshotter{}
	before, err := snapshotter.Snapshot(ctx, repoDir)
	if err != nil {
		t.Fatalf("before snapshot: %v", err)
	}

	// The Agent creates a fix.
	mustWriteForTest(t, filepath.Join(repoDir, "agent-fix.go"), "package fix\n")

	after, err := snapshotter.Snapshot(ctx, repoDir)
	if err != nil {
		t.Fatalf("after snapshot: %v", err)
	}
	changed := diffSnapshots(before, after)
	if len(changed) != 1 || changed[0] != "agent-fix.go" {
		t.Fatalf("expected only the Agent-made change in the diff, got %v", changed)
	}

	if err := (GitCommitter{}).Commit(ctx, CommitRequest{
		WorkDir: repoDir,
		Message: BatchCommitMessage(1),
		Paths:   changed,
	}); err != nil {
		t.Fatalf("snapshot-diff commit: %v", err)
	}

	committed := runGitForTest(t, repoDir, "show", "--name-only", "--format=", "HEAD")
	if !strings.Contains(committed, "agent-fix.go") {
		t.Fatalf("expected Agent change committed, got %q", committed)
	}
	if strings.Contains(committed, "tracked.txt") || strings.Contains(committed, "user-notes.md") {
		t.Fatalf("expected user changes kept out of the Batch commit, got %q", committed)
	}
	status := runGitForTest(t, repoDir, "status", "--porcelain=v1")
	if !strings.Contains(status, " M tracked.txt") || !strings.Contains(status, "?? user-notes.md") {
		t.Fatalf("expected user work preserved in the worktree, got %q", status)
	}
}
