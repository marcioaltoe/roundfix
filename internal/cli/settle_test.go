package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/daemon"
	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func TestRunSettleCommitsFailedTaskWorktreeWithDaemonMessage(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Recover the completed work",
			taskType:     "backend",
			status:       string(spec.StatusFailed),
			verification: []string{"test -f done.txt"},
		},
	})
	mustWrite(t, filepath.Join(repoDir, "done.txt"), "preserved work\n")
	mustWrite(t, filepath.Join(repoDir, "recovered.txt"), "sibling work\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected exit code 0, got %d (stderr %q)", code, stderr.String())
	}
	assertOnlySettleSurfaceLine(t, stderr.String(), repoDir)
	shortSHA := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", "--short", "HEAD"))
	expectedStdout := "verify test -f done.txt — ok\n" +
		"commit " + filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "task_01.md")) + "\n" +
		"commit done.txt\n" +
		"commit recovered.txt\n" +
		"settled task_01 completed — " + shortSHA + "\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("expected stdout:\n%q\ngot:\n%q", expectedStdout, stdout.String())
	}

	taskPath := implementTaskPath(repoDir, "task_01")
	if content := mustRead(t, taskPath); !strings.Contains(content, "status: completed") {
		t.Fatalf("expected task status completed, got:\n%s", content)
	}
	task := spec.Task{ID: "task_01", Title: "Recover the completed work", Type: "backend"}
	if message := strings.TrimRight(gitSettleOutput(t, repoDir, "log", "-1", "--format=%B"), "\n"); message != daemon.TaskCommitMessage(implementTestSlug, task) {
		t.Fatalf("expected daemon Task commit message %q, got %q", daemon.TaskCommitMessage(implementTestSlug, task), message)
	}
	changed := settleCommitFiles(t, repoDir)
	assertContainsString(t, changed, "done.txt")
	assertContainsString(t, changed, "recovered.txt")
	assertContainsString(t, changed, filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "task_01.md")))
	assertNoRunDatabase(t, homeDir)
}

func TestRunSettleWarnsWhenOtherSpecTasksAreFailed(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Recover the completed work", status: string(spec.StatusFailed), verification: []string{"true"}},
		{id: "task_02", title: "Other failed work", status: string(spec.StatusFailed)},
		{id: "task_03", title: "Completed sibling", status: string(spec.StatusCompleted)},
	})
	mustWrite(t, filepath.Join(repoDir, "done.txt"), "preserved work\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected exit code 0, got %d (stderr %q)", code, stderr.String())
	}
	assertSettleSurfaceLine(t, stderr.String(), repoDir)
	for _, expected := range []string{"warning: other failed Tasks", "task_02", "may have work included in this settle commit"} {
		if !strings.Contains(stderr.String(), expected) {
			t.Fatalf("expected stderr to contain %q, got %q", expected, stderr.String())
		}
	}
	if strings.Contains(stderr.String(), "task_01") || strings.Contains(stderr.String(), "task_03") {
		t.Fatalf("expected warning to name only other failed Tasks, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "commit done.txt\n") || !strings.Contains(stdout.String(), "settled task_01 completed — ") {
		t.Fatalf("expected settle success report with commit path, got %q", stdout.String())
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunSettleNoCommitPrintsNoCommitPathsOrSharedWarning(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Internal fixture should stay untouched"},
	})
	externalRoot := filepath.Join(t.TempDir(), "external-specs")
	writeImplementSpecAtRoot(t, externalRoot, implementTestSlug, []implementSeed{
		{id: "task_01", title: "Recover external task", status: string(spec.StatusFailed), verification: []string{"true"}},
		{id: "task_02", title: "Other failed work", status: string(spec.StatusFailed)},
	})
	configureExternalSpecsRoot(t, repoDir, externalRoot)
	beforeSHA := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", "--short", "HEAD"))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected settle exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	assertOnlySettleSurfaceLine(t, stderr.String(), repoDir)
	expectedStdout := "verify true — ok\n" +
		"settled task_01 completed — " + beforeSHA + "\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("expected stdout:\n%q\ngot:\n%q", expectedStdout, stdout.String())
	}
	if got := gitSettleOutput(t, repoDir, "status", "--porcelain=v1"); got != "" {
		t.Fatalf("expected no repository changes, got %q", got)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunSettleUsesConfiguredExternalSpecRoot(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Internal fixture should stay untouched"},
	})
	externalRoot := filepath.Join(t.TempDir(), "external-specs")
	writeImplementSpecAtRoot(t, externalRoot, implementTestSlug, []implementSeed{
		{
			id:           "task_01",
			title:        "Recover external task",
			status:       string(spec.StatusFailed),
			verification: []string{"true"},
		},
	})
	configureExternalSpecsRoot(t, repoDir, externalRoot)
	mustWrite(t, filepath.Join(repoDir, "recovered.txt"), "preserved work\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected settle exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	assertOnlySettleSurfaceLine(t, stderr.String(), repoDir)
	if !strings.Contains(stdout.String(), "verify true — ok\n") || !strings.Contains(stdout.String(), "settled task_01 completed — ") {
		t.Fatalf("expected settle success output, got %q", stdout.String())
	}
	externalTask := implementTaskPathInRoot(externalRoot, implementTestSlug, "task_01")
	if content := mustRead(t, externalTask); !strings.Contains(content, "status: completed") {
		t.Fatalf("expected external task completed, got:\n%s", content)
	}
	if content := mustRead(t, implementTaskPath(repoDir, "task_01")); !strings.Contains(content, "status: pending") {
		t.Fatalf("expected default-layout fixture untouched, got:\n%s", content)
	}
	task := spec.Task{ID: "task_01", Title: "Recover external task", Type: "backend"}
	if message := strings.TrimRight(gitSettleOutput(t, repoDir, "log", "-1", "--format=%B"), "\n"); message != daemon.TaskCommitMessage(implementTestSlug, task) {
		t.Fatalf("expected daemon Task commit message %q, got %q", daemon.TaskCommitMessage(implementTestSlug, task), message)
	}
	if changed := settleCommitFiles(t, repoDir); strings.Join(changed, "|") != "recovered.txt" {
		t.Fatalf("expected only repository recovery file committed, got %v", changed)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestSettleTaskStatusRetargetsKeptRunWorktreeAndCleansUpAfterIntegration(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Recover kept work",
			taskType:     "backend",
			verification: []string{"test -f done.txt"},
		},
	})
	runner := &implementFakeRunner{
		gitRoot: repoDir,
		onTask: func(req agent.ExecuteRequest, _ string) error {
			if err := os.WriteFile(filepath.Join(req.GitRoot, "done.txt"), []byte("preserved work\n"), 0o644); err != nil {
				return err
			}
			return errors.New("agent failed after writing recoverable work")
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected unresolved implement exit, got %d stderr=%q", code, stderr.String())
	}
	runID := implementRunIDFromStderr(t, stderr.String())
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateUnresolved {
		t.Fatalf("expected unresolved Run before settle, got %s", run.State)
	}
	assertRunWorktreeExists(t, run.WorkDir)
	assertRunBranchExists(t, repoDir, runworktree.BranchName(runID))
	stdout.Reset()
	stderr.Reset()

	code = runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected settle exit 0, got %d stderr=%q", code, stderr.String())
	}
	assertOnlySettleSurfaceLine(t, stderr.String(), run.WorkDir)
	if !strings.Contains(stdout.String(), "verify test -f done.txt — ok\n") {
		t.Fatalf("expected verification line, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "settled task_01 completed — ") {
		t.Fatalf("expected settled line, got %q", stdout.String())
	}
	settled := implementRunFromStore(t, homeDir, runID)
	if settled.State != store.StateUnresolved {
		t.Fatalf("expected settle to preserve the Run's Unresolved outcome, got %s", settled.State)
	}
	assertRunWorktreeRemoved(t, run.WorkDir)
	assertRunBranchRemoved(t, repoDir, runworktree.BranchName(runID))
	if got := mustRead(t, filepath.Join(repoDir, "done.txt")); got != "preserved work\n" {
		t.Fatalf("expected settled work integrated into user checkout, got %q", got)
	}
	if content := mustRead(t, implementTaskPath(repoDir, "task_01")); !strings.Contains(content, "status: completed") {
		t.Fatalf("expected task status completed after settle, got:\n%s", content)
	}
	if status := gitSettleOutput(t, repoDir, "status", "--porcelain=v1"); status != "" {
		t.Fatalf("expected clean user checkout after settle integration, got %q", status)
	}
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}

func TestRunSettleSkipsStaleKeptRunWorktreeAndUsesFailedCheckout(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Recover checkout work",
			taskType:     "backend",
			status:       string(spec.StatusFailed),
			verification: []string{"test -f done.txt"},
		},
	})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	_, olderRef, _ := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "", store.StateUnresolved)
	if err := spec.SetStatus(implementTaskPath(olderRef.Path, "task_01"), spec.StatusFailed); err != nil {
		t.Fatalf("set older Run Worktree task failed: %v", err)
	}
	_, latestRef, _ := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "", store.StateUnresolved)
	if err := spec.SetStatus(implementTaskPath(latestRef.Path, "task_01"), spec.StatusPending); err != nil {
		t.Fatalf("set stale Run Worktree task pending: %v", err)
	}
	mustWrite(t, filepath.Join(repoDir, "done.txt"), "checkout work\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected settle exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	assertOnlySettleSurfaceLine(t, stderr.String(), repoDir)
	if !strings.Contains(stdout.String(), "verify test -f done.txt — ok\n") || !strings.Contains(stdout.String(), "settled task_01 completed — ") {
		t.Fatalf("expected settle success from checkout, got stdout=%q", stdout.String())
	}
	if content := mustRead(t, implementTaskPath(repoDir, "task_01")); !strings.Contains(content, "status: completed") {
		t.Fatalf("expected checkout task settled completed, got:\n%s", content)
	}
	if content := mustRead(t, implementTaskPath(latestRef.Path, "task_01")); !strings.Contains(content, "status: pending") {
		t.Fatalf("expected stale Run Worktree task to remain pending, got:\n%s", content)
	}
}

func TestRunSettleRetargetsKeptTaskWorktreeAndCleansUpAfterIntegration(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Recover concurrent work",
			taskType:     "backend",
			status:       string(spec.StatusFailed),
			verification: []string{"test -f done.txt"},
		},
	})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	run, _, taskRef := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "task_01", store.StateUnresolved)
	if err := os.WriteFile(filepath.Join(taskRef.Path, "done.txt"), []byte("task work\n"), 0o644); err != nil {
		t.Fatalf("write task work: %v", err)
	}
	if err := spec.SetStatus(implementTaskPath(taskRef.Path, "task_01"), spec.StatusFailed); err != nil {
		t.Fatalf("settle task failed in Task Worktree: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected settle exit 0, got %d stderr=%q", code, stderr.String())
	}
	assertOnlySettleSurfaceLine(t, stderr.String(), taskRef.Path)
	shortSHA := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", "--short", "HEAD"))
	expectedStdout := "verify test -f done.txt — ok\n" +
		"commit " + filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "task_01.md")) + "\n" +
		"commit done.txt\n" +
		"settled task_01 completed — " + shortSHA + "\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("expected stdout:\n%q\ngot:\n%q", expectedStdout, stdout.String())
	}
	settled := implementRunFromStore(t, homeDir, run.ID)
	if settled.State != store.StateUnresolved {
		t.Fatalf("expected settle to preserve the Run's Unresolved outcome, got %s", settled.State)
	}
	assertRunWorktreeRemoved(t, run.WorkDir)
	assertRunBranchRemoved(t, repoDir, runworktree.BranchName(run.ID))
	assertRunWorktreeRemoved(t, taskRef.Path)
	assertRunBranchRemoved(t, repoDir, taskRef.Branch)
	if got := mustRead(t, filepath.Join(repoDir, "done.txt")); got != "task work\n" {
		t.Fatalf("expected Task Worktree work integrated into user checkout, got %q", got)
	}
	if content := mustRead(t, implementTaskPath(repoDir, "task_01")); !strings.Contains(content, "status: completed") {
		t.Fatalf("expected task status completed after settle, got:\n%s", content)
	}
	if status := gitSettleOutput(t, repoDir, "status", "--porcelain=v1"); status != "" {
		t.Fatalf("expected clean user checkout after settle integration, got %q", status)
	}
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}

func TestRunSettleTaskWorktreeIntegrationConflictKeepsSurfaces(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Recover conflicting work",
			taskType:     "backend",
			status:       string(spec.StatusFailed),
			verification: []string{"grep -q task shared.txt"},
		},
	})
	mustWrite(t, filepath.Join(repoDir, "shared.txt"), "base\n")
	gitImplement(t, repoDir, "add", "shared.txt")
	gitImplement(t, repoDir, "commit", "-m", "seed shared file")
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	run, runRef, taskRef := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "task_01", store.StateUnresolved)
	if err := os.WriteFile(filepath.Join(taskRef.Path, "shared.txt"), []byte("task\n"), 0o644); err != nil {
		t.Fatalf("write task shared file: %v", err)
	}
	if err := spec.SetStatus(implementTaskPath(taskRef.Path, "task_01"), spec.StatusFailed); err != nil {
		t.Fatalf("settle task failed in Task Worktree: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runRef.Path, "shared.txt"), []byte("run\n"), 0o644); err != nil {
		t.Fatalf("write Run Worktree conflict: %v", err)
	}
	gitImplement(t, runRef.Path, "add", "shared.txt")
	gitImplement(t, runRef.Path, "commit", "-m", "conflicting run branch work")
	runTip := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", runRef.Branch))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected settle conflict exit 1, got %d", code)
	}
	if stdout.String() != "verify grep -q task shared.txt — ok\n" {
		t.Fatalf("expected only verification success on stdout, got %q", stdout.String())
	}
	for _, expected := range []string{"roundfix: settle failed after verification: task worktree integration conflict on shared.txt"} {
		if !strings.Contains(stderr.String(), expected) {
			t.Fatalf("expected stderr to contain %q, got %q", expected, stderr.String())
		}
	}
	if got := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", runRef.Branch)); got != runTip {
		t.Fatalf("expected Run Branch tip %s to stay unchanged, got %s", runTip, got)
	}
	assertRunWorktreeExists(t, run.WorkDir)
	assertRunBranchExists(t, repoDir, runworktree.BranchName(run.ID))
	assertRunWorktreeExists(t, taskRef.Path)
	assertRunBranchExists(t, repoDir, taskRef.Branch)
	if got := mustRead(t, filepath.Join(repoDir, "shared.txt")); got != "base\n" {
		t.Fatalf("expected user checkout unchanged, got %q", got)
	}
	current := implementRunFromStore(t, homeDir, run.ID)
	if current.State != store.StateUnresolved {
		t.Fatalf("expected Run to remain Unresolved, got %s", current.State)
	}
}

func TestSettleAcceptsCompletedTaskWithUncommittedWorkInCheckout(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Recover work a hook refused",
			taskType:     "backend",
			status:       string(spec.StatusCompleted),
			verification: []string{"test -f done.txt"},
		},
	})
	mustWrite(t, filepath.Join(repoDir, "done.txt"), "verified work\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected exit code 0, got %d (stderr %q)", code, stderr.String())
	}
	assertOnlySettleSurfaceLine(t, stderr.String(), repoDir)
	shortSHA := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", "--short", "HEAD"))
	expectedStdout := "verify test -f done.txt — ok\n" +
		"commit done.txt\n" +
		"settled task_01 completed — " + shortSHA + "\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("expected stdout:\n%q\ngot:\n%q", expectedStdout, stdout.String())
	}
	if changed := settleCommitFiles(t, repoDir); strings.Join(changed, "|") != "done.txt" {
		t.Fatalf("expected only the uncommitted work committed, got %v", changed)
	}
	if content := mustRead(t, implementTaskPath(repoDir, "task_01")); !strings.Contains(content, "status: completed") {
		t.Fatalf("expected task status to stay completed, got:\n%s", content)
	}
	if got := gitSettleOutput(t, repoDir, "status", "--porcelain=v1"); got != "" {
		t.Fatalf("expected clean checkout after settle, got %q", got)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestSettleAcceptsCompletedTaskFromKeptTaskWorktreeAfterHookRefusal(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Recover work a hook refused",
			taskType:     "backend",
			status:       string(spec.StatusInProgress),
			verification: []string{"test -f never-written.txt"},
		},
	})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	run, _, taskRef := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "task_01", store.StateUnresolved)
	// The shape a refused commit leaves behind: the Daemon settled the Task
	// completed after its Verification passed, git add ran, and the hook
	// refused the commit, so the verified work stays staged in the surface.
	mustWrite(t, filepath.Join(taskRef.Path, "done.txt"), "task work\n")
	gitImplement(t, taskRef.Path, "add", "done.txt")
	mustWrite(t, implementTaskPath(taskRef.Path, "task_01"), implementTaskContent(implementTestSlug, implementSeed{
		id:           "task_01",
		title:        "Recover work a hook refused",
		taskType:     "backend",
		status:       string(spec.StatusCompleted),
		verification: []string{"test -f done.txt"},
	}))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected settle exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	assertOnlySettleSurfaceLine(t, stderr.String(), taskRef.Path)
	// The Task Worktree's own task file declares this command; the checkout's
	// copy declares one that cannot pass.
	if !strings.Contains(stdout.String(), "verify test -f done.txt — ok\n") {
		t.Fatalf("expected Verification loaded from the Task Worktree, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "settled task_01 completed — ") {
		t.Fatalf("expected settled line, got %q", stdout.String())
	}
	assertRunWorktreeRemoved(t, taskRef.Path)
	assertRunWorktreeRemoved(t, run.WorkDir)
	if got := mustRead(t, filepath.Join(repoDir, "done.txt")); got != "task work\n" {
		t.Fatalf("expected refused work integrated into user checkout, got %q", got)
	}
	if content := mustRead(t, implementTaskPath(repoDir, "task_01")); !strings.Contains(content, "status: completed") {
		t.Fatalf("expected checkout task file completed after settle, got:\n%s", content)
	}
	if status := gitSettleOutput(t, repoDir, "status", "--porcelain=v1"); status != "" {
		t.Fatalf("expected clean user checkout after settle integration, got %q", status)
	}
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}

// A refused commit leaves the Daemon's staging behind, so settle meets a
// surface whose removals are already in the index. Those paths are in neither
// the index nor the worktree, so re-staging them by pathspec matches nothing;
// settle stages the whole surface instead, and the Task commit carries the
// removals onto the Run Branch.
func TestSettleCommitsDeletedAndRenamedWorkFromTaskWorktree(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Recover work a hook refused",
			taskType:     "backend",
			status:       string(spec.StatusInProgress),
			verification: []string{"test -f never-written.txt"},
		},
	})
	mustWrite(t, filepath.Join(repoDir, "removed.txt"), "obsolete\n")
	mustWrite(t, filepath.Join(repoDir, "original.txt"), "moved content\n")
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "seed the files the Task rewrites")
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	run, _, taskRef := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "task_01", store.StateUnresolved)
	if err := os.Remove(filepath.Join(taskRef.Path, "removed.txt")); err != nil {
		t.Fatalf("delete task work: %v", err)
	}
	if err := os.Rename(filepath.Join(taskRef.Path, "original.txt"), filepath.Join(taskRef.Path, "renamed.txt")); err != nil {
		t.Fatalf("rename task work: %v", err)
	}
	mustWrite(t, implementTaskPath(taskRef.Path, "task_01"), implementTaskContent(implementTestSlug, implementSeed{
		id:           "task_01",
		title:        "Recover work a hook refused",
		taskType:     "backend",
		status:       string(spec.StatusCompleted),
		verification: []string{"test -f renamed.txt"},
	}))
	// The shape a refused commit leaves behind: the Daemon staged the whole
	// Task, including both removals, and the hook refused the commit.
	gitImplement(t, taskRef.Path, "add", "-A")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected settle exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	assertOnlySettleSurfaceLine(t, stderr.String(), taskRef.Path)
	shortSHA := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", "--short", "HEAD"))
	expectedStdout := "verify test -f renamed.txt — ok\n" +
		"commit " + filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "task_01.md")) + "\n" +
		"commit original.txt — deleted\n" +
		"commit removed.txt — deleted\n" +
		"commit renamed.txt\n" +
		"settled task_01 completed — " + shortSHA + "\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("expected stdout:\n%q\ngot:\n%q", expectedStdout, stdout.String())
	}
	nameStatus := gitSettleOutput(t, repoDir, "diff-tree", "--no-commit-id", "--name-status", "--no-renames", "-r", "HEAD")
	for _, expected := range []string{"D\toriginal.txt", "D\tremoved.txt", "A\trenamed.txt"} {
		if !strings.Contains(nameStatus, expected) {
			t.Fatalf("expected Task commit to record %q, got:\n%s", expected, nameStatus)
		}
	}
	for _, gone := range []string{"removed.txt", "original.txt"} {
		if _, err := os.Stat(filepath.Join(repoDir, gone)); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected %s removed from the user checkout, got %v", gone, err)
		}
	}
	if got := mustRead(t, filepath.Join(repoDir, "renamed.txt")); got != "moved content\n" {
		t.Fatalf("expected renamed work integrated into user checkout, got %q", got)
	}
	assertRunWorktreeRemoved(t, taskRef.Path)
	assertRunBranchRemoved(t, repoDir, taskRef.Branch)
	assertRunWorktreeRemoved(t, run.WorkDir)
	if status := gitSettleOutput(t, repoDir, "status", "--porcelain=v1"); status != "" {
		t.Fatalf("expected clean user checkout after settle integration, got %q", status)
	}
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}

// The other shape a Task's deletion reaches settle in: the removal is still
// unstaged, and the surface holding it is the user checkout. The stage-all
// records the removal without a pathspec that matches nothing, the Task commit
// carries it, and the report names the path the commit deletes.
func TestSettleCommitsUnstagedDeletionFromCheckoutSurface(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Remove the obsolete file",
			taskType:     "backend",
			status:       string(spec.StatusFailed),
			verification: []string{"test ! -e obsolete.txt"},
		},
	})
	mustWrite(t, filepath.Join(repoDir, "obsolete.txt"), "obsolete\n")
	mustWrite(t, filepath.Join(repoDir, "kept.txt"), "kept\n")
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "seed the files the Task rewrites")
	if err := os.Remove(filepath.Join(repoDir, "obsolete.txt")); err != nil {
		t.Fatalf("delete task work: %v", err)
	}
	mustWrite(t, filepath.Join(repoDir, "kept.txt"), "rewritten\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected settle exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	assertOnlySettleSurfaceLine(t, stderr.String(), repoDir)
	shortSHA := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", "--short", "HEAD"))
	expectedStdout := "verify test ! -e obsolete.txt — ok\n" +
		"commit " + filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "task_01.md")) + "\n" +
		"commit kept.txt\n" +
		"commit obsolete.txt — deleted\n" +
		"settled task_01 completed — " + shortSHA + "\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("expected stdout:\n%q\ngot:\n%q", expectedStdout, stdout.String())
	}
	nameStatus := gitSettleOutput(t, repoDir, "diff-tree", "--no-commit-id", "--name-status", "--no-renames", "-r", "HEAD")
	if !strings.Contains(nameStatus, "D\tobsolete.txt") {
		t.Fatalf("expected Task commit to record the deletion, got:\n%s", nameStatus)
	}
	if _, err := os.Stat(filepath.Join(repoDir, "obsolete.txt")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected obsolete.txt to stay removed after settle, got %v", err)
	}
	if status := gitSettleOutput(t, repoDir, "status", "--porcelain=v1"); status != "" {
		t.Fatalf("expected clean checkout after settle, got %q", status)
	}
	assertNoRunDatabase(t, homeDir)
}

// measuredHookRefusedWork is one of the three findings whose Runs died at the
// commit step: work the authoritative Verification passed and a stricter
// repository hook then refused.
type measuredHookRefusedWork struct {
	path      string
	content   string
	objection string
}

func measuredHookRefusedWorkForTest() []measuredHookRefusedWork {
	var parser strings.Builder
	parser.WriteString("package parser\n\nfunc Parse(input string) []string {\n\tfields := []string{}\n")
	for index := 0; index < 78; index++ {
		fmt.Fprintf(&parser, "\tfields = append(fields, field%02d(input))\n", index)
	}
	parser.WriteString("\treturn fields\n}\n")
	schema := make([]string, 0, 2462)
	schema = append(schema, "// Code generated by the schema generator. DO NOT EDIT.")
	for index := 1; index <= 2461; index++ {
		schema = append(schema, fmt.Sprintf("export type Field%04d = { id: number }", index))
	}
	return []measuredHookRefusedWork{
		{
			path:      "src/parser.go",
			content:   parser.String(),
			objection: "src/parser.go:3: function exceeds the 80-line limit",
		},
		{
			path:      "internal/api/schema.gen.ts",
			content:   strings.Join(schema, "\n") + "\n",
			objection: "internal/api/schema.gen.ts: 2462 lines exceeds the 500-line limit",
		},
		{
			path:      "src/list.ts",
			content:   "export function names(items: Item[]): string[] {\n  return items.map((item) => item.name).sort();\n}\n",
			objection: "src/list.ts:2: use toSorted() instead of sort()",
		},
	}
}

// measuredHookRefusalScriptForTest enforces the three measured rules over the
// staged content, so the refusal a test observes is computed rather than
// declared.
const measuredHookRefusalScriptForTest = `exec 1>&2
refused=0
for file in $(git diff --cached --name-only --diff-filter=ACM); do
  case "$file" in
    *.go)
      if awk -v path="$file" '/^func /{start=NR} /^}/{if (start && NR-start+1 > 80) {printf "%s:%d: function exceeds the 80-line limit\n", path, start; bad=1} start=0} END{exit !bad}' "$file"; then
        refused=1
      fi
      ;;
    *.ts)
      line=$(grep -n '\.sort(' "$file" | grep -v toSorted | head -1 | cut -d: -f1)
      if [ -n "$line" ]; then
        echo "$file:$line: use toSorted() instead of sort()"
        refused=1
      fi
      ;;
  esac
  lines=$(wc -l < "$file" | tr -d ' ')
  if [ "$lines" -gt 500 ]; then
    echo "$file: $lines lines exceeds the 500-line limit"
    refused=1
  fi
done
if [ "$refused" -ne 0 ]; then
  echo 'husky - pre-commit script failed (code 1)'
  exit 1
fi
exit 0
`

// TestSettleRecoversMeasuredHookRefusedWork settles the three measured
// findings through the settle Command itself. The surface holds work whose
// Verification passes and whose findings a stricter commit hook refuses, which
// is the state the three Run deaths left behind.
//
// The first settle meets the hook still in place: it re-verifies, stages, is
// refused at the commit, and reports the refusal without discarding anything.
// The Supervisor then repairs the misconfiguration the hook-strictness
// invariant names, and the second settle commits the same work untouched.
func TestSettleRecoversMeasuredHookRefusedWork(t *testing.T) {
	t.Parallel()
	work := measuredHookRefusedWorkForTest()
	verification := make([]string, 0, len(work))
	for _, item := range work {
		verification = append(verification, "test -f "+item.path)
	}
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Recover work a hook refused",
			taskType:     "backend",
			status:       string(spec.StatusCompleted),
			verification: verification,
		},
	})
	hooksDir := filepath.Join(repoDir, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatalf("create hooks directory: %v", err)
	}
	// A repository-local hooks path, so a global core.hooksPath cannot decide
	// whether this hook runs.
	gittest.AppendConfig(t, repoDir, "[core]\n\thooksPath = "+hooksDir+"\n")
	hookPath := filepath.Join(hooksDir, "pre-commit")
	if err := os.WriteFile(hookPath, []byte("#!/bin/sh\n"+measuredHookRefusalScriptForTest), 0o755); err != nil {
		t.Fatalf("write pre-commit hook: %v", err)
	}
	for _, item := range work {
		full := filepath.Join(repoDir, filepath.FromSlash(item.path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("create work directory: %v", err)
		}
		mustWrite(t, full, item.content)
	}
	seedHead := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", "HEAD"))

	var refusedStdout bytes.Buffer
	var refusedStderr bytes.Buffer
	refusedCode := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &refusedStdout, &refusedStderr)

	if refusedCode != exitRunFailed {
		t.Fatalf("expected settle to stop on the hook refusal with exit %d, got %d (stderr %q)", exitRunFailed, refusedCode, refusedStderr.String())
	}
	if !strings.Contains(refusedStderr.String(), "pre-commit hook refused the commit") {
		t.Fatalf("expected the refusing hook named, got %q", refusedStderr.String())
	}
	for _, item := range work {
		if !strings.Contains(refusedStderr.String(), item.objection) {
			t.Fatalf("expected the hook's finding %q reported, got %q", item.objection, refusedStderr.String())
		}
		if got := mustRead(t, filepath.Join(repoDir, filepath.FromSlash(item.path))); got != item.content {
			t.Fatalf("expected %s untouched by the refused settle, got %d of %d bytes", item.path, len(got), len(item.content))
		}
	}
	if got := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", "HEAD")); got != seedHead {
		t.Fatalf("expected no commit created while the hook refuses, got %s", got)
	}
	if content := mustRead(t, implementTaskPath(repoDir, "task_01")); !strings.Contains(content, "status: completed") {
		t.Fatalf("expected the Task to stay completed after the refusal, got:\n%s", content)
	}

	// The invariant's repair: the rules that outranked the authoritative
	// Verification leave the hook.
	if err := os.Remove(hookPath); err != nil {
		t.Fatalf("repair the hook configuration: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected settle exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	assertOnlySettleSurfaceLine(t, stderr.String(), repoDir)
	shortSHA := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", "--short", "HEAD"))
	if !strings.Contains(stdout.String(), "settled task_01 completed — "+shortSHA+"\n") {
		t.Fatalf("expected the settled line naming the recovery commit, got %q", stdout.String())
	}
	committed := settleCommitFiles(t, repoDir)
	for _, item := range work {
		assertContainsString(t, committed, item.path)
		if got := gitSettleOutput(t, repoDir, "show", "HEAD:"+item.path); got != item.content {
			t.Fatalf("expected %s recovered byte for byte, got %d of %d bytes", item.path, len(got), len(item.content))
		}
	}
	if got := gitSettleOutput(t, repoDir, "status", "--porcelain=v1"); got != "" {
		t.Fatalf("expected clean checkout after settle, got %q", got)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestSettleRefusesCompletedTaskWhenNoSurfaceHoldsUncommittedWork(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Already committed work",
			taskType:     "backend",
			status:       string(spec.StatusCompleted),
			verification: []string{"true"},
		},
	})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	_, runRef, _ := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "", store.StateUnresolved)
	beforeHeadCount := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-list", "--count", "HEAD"))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected preflight exit 2, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	for _, expected := range []string{
		"no settleable surface",
		runRef.Path + ": status completed (no uncommitted work)",
		repoDir + ": status completed (no uncommitted work)",
		"status is completed with no uncommitted work; nothing to settle",
	} {
		if !strings.Contains(stderr.String(), expected) {
			t.Fatalf("expected stderr to contain %q, got %q", expected, stderr.String())
		}
	}
	if got := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-list", "--count", "HEAD")); got != beforeHeadCount {
		t.Fatalf("expected HEAD count %s, got %s", beforeHeadCount, got)
	}
	assertRunWorktreeExists(t, runRef.Path)
}

func TestRunSettleRefusalEnumeratesCandidateStatuses(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "No failed surface",
			status:       string(spec.StatusCompleted),
			verification: []string{"true"},
		},
	})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	_, runRef, _ := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "", store.StateUnresolved)
	if err := spec.SetStatus(implementTaskPath(runRef.Path, "task_01"), spec.StatusPending); err != nil {
		t.Fatalf("set Run Worktree task pending: %v", err)
	}
	missingTaskRef, err := runworktree.TaskRefFor(runRef, "task_01")
	if err != nil {
		t.Fatalf("derive Task Worktree path: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected preflight exit 2, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	for _, expected := range []string{
		"no settleable surface",
		missingTaskRef.Path + ": path does not exist",
		runRef.Path + ": status pending",
		repoDir + ": status completed (no uncommitted work)",
		"status is pending; run the Implement Command",
	} {
		if !strings.Contains(stderr.String(), expected) {
			t.Fatalf("expected stderr to contain %q, got %q", expected, stderr.String())
		}
	}
}

func TestRunSettleRequiresSpecAndTask(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    []string
		message string
	}{
		{
			name:    "missing spec",
			args:    []string{"settle", "--task", "task_01"},
			message: "missing required --spec",
		},
		{
			name:    "missing task",
			args:    []string{"settle", "--spec", implementTestSlug},
			message: "missing required --task",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			setCommandEnvironmentForTest(t, homeDir, t.TempDir())
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runCLIContext(t, context.Background(), tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected exit code 2, got %d", code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.message) {
				t.Fatalf("expected stderr to contain %q, got %q", tt.message, stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunSettlePreflightRefusalsWriteNothing(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		setup   func(t *testing.T) (homeDir string, repoDir string, taskPath string)
		args    []string
		message string
	}{
		{
			name: "repository must resolve",
			setup: func(t *testing.T) (string, string, string) {
				homeDir := t.TempDir()
				repoDir := t.TempDir()
				setCommandEnvironmentForTest(t, homeDir, repoDir)
				return homeDir, repoDir, ""
			},
			args:    []string{"settle", "--spec", implementTestSlug, "--task", "task_01"},
			message: "repository resolves",
		},
		{
			name: "spec must load",
			setup: func(t *testing.T) (string, string, string) {
				homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusFailed)}})
				mustWrite(t, filepath.Join(repoDir, "docs", "specs", implementTestSlug, "_tasks.md"), "not frontmatter\n")
				return homeDir, repoDir, implementTaskPath(repoDir, "task_01")
			},
			args:    []string{"settle", "--spec", implementTestSlug, "--task", "task_01"},
			message: "Spec loads valid",
		},
		{
			name: "task id must exist",
			setup: func(t *testing.T) (string, string, string) {
				homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusFailed)}})
				return homeDir, repoDir, implementTaskPath(repoDir, "task_01")
			},
			args:    []string{"settle", "--spec", implementTestSlug, "--task", "task_missing"},
			message: `Task "task_missing" does not exist in Task Graph`,
		},
		{
			name: "pending task points to implement",
			setup: func(t *testing.T) (string, string, string) {
				homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusPending)}})
				return homeDir, repoDir, implementTaskPath(repoDir, "task_01")
			},
			args:    []string{"settle", "--spec", implementTestSlug, "--task", "task_01"},
			message: "status is pending; run the Implement Command",
		},
		{
			name: "in progress task points to implement",
			setup: func(t *testing.T) (string, string, string) {
				homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusInProgress)}})
				return homeDir, repoDir, implementTaskPath(repoDir, "task_01")
			},
			args:    []string{"settle", "--spec", implementTestSlug, "--task", "task_01"},
			message: "status is in_progress; run the Implement Command",
		},
		{
			name: "completed task with committed work has nothing to settle",
			setup: func(t *testing.T) (string, string, string) {
				homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
				return homeDir, repoDir, implementTaskPath(repoDir, "task_01")
			},
			args:    []string{"settle", "--spec", implementTestSlug, "--task", "task_01"},
			message: "status is completed with no uncommitted work; nothing to settle",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir, taskPath := tt.setup(t)
			var beforeTask string
			if taskPath != "" {
				beforeTask = mustRead(t, taskPath)
			}
			beforeHeadCount := ""
			beforeStatus := ""
			if _, err := os.Stat(filepath.Join(repoDir, ".git")); err == nil {
				beforeHeadCount = strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-list", "--count", "HEAD"))
				beforeStatus = gitSettleOutput(t, repoDir, "status", "--porcelain=v1")
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runCLIContext(t, context.Background(), tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected exit code 2, got %d (stderr %q)", code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.message) {
				t.Fatalf("expected stderr to contain %q, got %q", tt.message, stderr.String())
			}
			if taskPath != "" && mustRead(t, taskPath) != beforeTask {
				t.Fatalf("expected task file unchanged")
			}
			if beforeHeadCount != "" {
				if got := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-list", "--count", "HEAD")); got != beforeHeadCount {
					t.Fatalf("expected HEAD count %s, got %s", beforeHeadCount, got)
				}
				if got := gitSettleOutput(t, repoDir, "status", "--porcelain=v1"); got != beforeStatus {
					t.Fatalf("expected git status unchanged %q, got %q", beforeStatus, got)
				}
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunSettleVerificationFailureLeavesTaskAndTreeUntouched(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:     "task_01",
			status: string(spec.StatusFailed),
			verification: []string{
				"test -f done.txt",
				"test -f missing.txt",
				"touch should-not-run",
			},
		},
	})
	mustWrite(t, filepath.Join(repoDir, "done.txt"), "preserved work\n")
	taskPath := implementTaskPath(repoDir, "task_01")
	beforeTask := mustRead(t, taskPath)
	beforeStatus := gitSettleOutput(t, repoDir, "status", "--porcelain=v1")
	beforeHeadCount := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-list", "--count", "HEAD"))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected exit code 1, got %d (stderr %q)", code, stderr.String())
	}
	diagnosticPath := daemon.VerificationOutputPath(builtinArtifactDirForRepo(t, repoDir), "settle-"+implementTestSlug+"-task_01", 1, 1)
	expectedStdout := "verify test -f done.txt — ok\n" +
		"verify test -f missing.txt — failed (diagnostics: " + diagnosticPath + ")\n" +
		"task_01 stays failed — verification failed\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("expected stdout:\n%q\ngot:\n%q", expectedStdout, stdout.String())
	}
	if _, err := os.Stat(diagnosticPath); err != nil {
		t.Fatalf("expected failed verification diagnostic artifact: %v", err)
	}
	if got := mustRead(t, taskPath); got != beforeTask {
		t.Fatalf("expected task file unchanged")
	}
	if got := gitSettleOutput(t, repoDir, "status", "--porcelain=v1"); got != beforeStatus {
		t.Fatalf("expected git status unchanged %q, got %q", beforeStatus, got)
	}
	if got := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-list", "--count", "HEAD")); got != beforeHeadCount {
		t.Fatalf("expected HEAD count %s, got %s", beforeHeadCount, got)
	}
	if _, err := os.Stat(filepath.Join(repoDir, "should-not-run")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected later verification command not to run, stat error %v", err)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestSettleVerificationRunsSurfaceCommandsVerbatim(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Recover work a hook refused", taskType: "backend", status: string(spec.StatusCompleted)},
	})
	// The commands land in the task file only after the workspace exists, so
	// the first one can name the surface settle must run them in.
	commands := []string{
		"test \"$(pwd -P)\" = " + shellQuoteImplement(repoDir),
		"printf 'a|b' | grep -q 'a|b'",
		"test -f done.txt",
	}
	mustWrite(t, implementTaskPath(repoDir, "task_01"), implementTaskContent(implementTestSlug, implementSeed{
		id:           "task_01",
		title:        "Recover work a hook refused",
		taskType:     "backend",
		status:       string(spec.StatusCompleted),
		verification: commands,
	}))
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "declare settle verification commands")
	mustWrite(t, filepath.Join(repoDir, "done.txt"), "verified work\n")
	runner := &implementFakeRunner{gitRoot: repoDir}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected exit code 0, got %d (stderr %q)", code, stderr.String())
	}
	assertOnlySettleSurfaceLine(t, stderr.String(), repoDir)
	shortSHA := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", "--short", "HEAD"))
	expectedStdout := "verify " + commands[0] + " — ok\n" +
		"verify " + commands[1] + " — ok\n" +
		"verify " + commands[2] + " — ok\n" +
		"commit done.txt\n" +
		"settled task_01 completed — " + shortSHA + "\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("expected stdout:\n%q\ngot:\n%q", expectedStdout, stdout.String())
	}
	if runner.calls != 0 {
		t.Fatalf("expected settle to open no Agent session, got %d", runner.calls)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestSettleVerificationFailureKeepsHookRefusedWorkInTaskWorktree(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Recover work a hook refused",
			taskType:     "backend",
			status:       string(spec.StatusInProgress),
			verification: []string{"test -f never-written.txt"},
		},
	})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	run, _, taskRef := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "task_01", store.StateUnresolved)
	// The shape a refused commit leaves behind, with a Verification that no
	// longer passes in the surface holding the staged work.
	mustWrite(t, filepath.Join(taskRef.Path, "done.txt"), "task work\n")
	gitImplement(t, taskRef.Path, "add", "done.txt")
	taskPath := implementTaskPath(taskRef.Path, "task_01")
	mustWrite(t, taskPath, implementTaskContent(implementTestSlug, implementSeed{
		id:       "task_01",
		title:    "Recover work a hook refused",
		taskType: "backend",
		status:   string(spec.StatusCompleted),
		verification: []string{
			"test -f done.txt",
			"test -f missing.txt",
			"touch should-not-run",
		},
	}))
	beforeTask := mustRead(t, taskPath)
	beforeStatus := gitSettleOutput(t, taskRef.Path, "status", "--porcelain=v1")
	beforeHeadCount := strings.TrimSpace(gitSettleOutput(t, taskRef.Path, "rev-list", "--count", "HEAD"))
	runner := &implementFakeRunner{gitRoot: repoDir}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected exit code 1, got %d (stderr %q)", code, stderr.String())
	}
	assertOnlySettleSurfaceLine(t, stderr.String(), taskRef.Path)
	diagnosticPath := daemon.VerificationOutputPath(builtinArtifactDirForRepo(t, repoDir), "settle-"+implementTestSlug+"-task_01", 1, 1)
	expectedStdout := "verify test -f done.txt — ok\n" +
		"verify test -f missing.txt — failed (diagnostics: " + diagnosticPath + ")\n" +
		"task_01 stays completed — verification failed\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("expected stdout:\n%q\ngot:\n%q", expectedStdout, stdout.String())
	}
	if runner.calls != 0 {
		t.Fatalf("expected settle to open no Agent session, got %d", runner.calls)
	}
	if _, err := os.Stat(filepath.Join(taskRef.Path, "should-not-run")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected later verification command not to run, stat error %v", err)
	}
	assertRunWorktreeExists(t, taskRef.Path)
	assertRunWorktreeExists(t, run.WorkDir)
	if got := mustRead(t, taskPath); got != beforeTask {
		t.Fatalf("expected task file unchanged")
	}
	if got := gitSettleOutput(t, taskRef.Path, "status", "--porcelain=v1"); got != beforeStatus {
		t.Fatalf("expected Task Worktree status unchanged %q, got %q", beforeStatus, got)
	}
	if !strings.Contains(beforeStatus, "A  done.txt") {
		t.Fatalf("expected refused work to stay staged, got %q", beforeStatus)
	}
	if got := strings.TrimSpace(gitSettleOutput(t, taskRef.Path, "rev-list", "--count", "HEAD")); got != beforeHeadCount {
		t.Fatalf("expected HEAD count %s, got %s", beforeHeadCount, got)
	}
	if _, err := os.Stat(filepath.Join(repoDir, "done.txt")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected refused work to stay out of the checkout, stat error %v", err)
	}
}

func TestSettleVerificationUnknownVerdictStopsWithoutCommitting(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Recover work a hook refused",
			taskType:     "backend",
			status:       string(spec.StatusCompleted),
			verification: []string{"go test ./...", "touch should-not-run"},
		},
	})
	mustWrite(t, filepath.Join(repoDir, "done.txt"), "verified work\n")
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.verifier = settleUnknownVerdictVerifier{cause: errors.New("fork/exec /bin/sh: permission denied")}
	})
	taskPath := implementTaskPath(repoDir, "task_01")
	beforeTask := mustRead(t, taskPath)
	beforeStatus := gitSettleOutput(t, repoDir, "status", "--porcelain=v1")
	beforeHeadCount := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-list", "--count", "HEAD"))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected exit code 1, got %d (stderr %q)", code, stderr.String())
	}
	expectedStdout := "verify go test ./... — verdict unknown\n" +
		"task_01 stays completed — verification verdict unknown\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("expected stdout:\n%q\ngot:\n%q", expectedStdout, stdout.String())
	}
	assertSettleSurfaceLine(t, stderr.String(), repoDir)
	for _, expected := range []string{
		"settle verification verdict unknown:",
		"fork/exec /bin/sh: permission denied",
	} {
		if !strings.Contains(stderr.String(), expected) {
			t.Fatalf("expected stderr to contain %q, got %q", expected, stderr.String())
		}
	}
	if _, err := os.Stat(filepath.Join(repoDir, "should-not-run")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected later verification command not to run, stat error %v", err)
	}
	if got := mustRead(t, taskPath); got != beforeTask {
		t.Fatalf("expected task file unchanged")
	}
	if got := gitSettleOutput(t, repoDir, "status", "--porcelain=v1"); got != beforeStatus {
		t.Fatalf("expected git status unchanged %q, got %q", beforeStatus, got)
	}
	if got := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-list", "--count", "HEAD")); got != beforeHeadCount {
		t.Fatalf("expected HEAD count %s, got %s", beforeHeadCount, got)
	}
	assertNoRunDatabase(t, homeDir)
}

// settleUnknownVerdictVerifier stands in for a command runner that never
// observed a verdict — the state daemon.ExecVerifier reports when it cannot
// start the command or retain its diagnostics.
type settleUnknownVerdictVerifier struct {
	cause error
}

func (verifier settleUnknownVerdictVerifier) Verify(_ context.Context, req daemon.VerifyRequest) (daemon.VerifyResult, error) {
	return daemon.VerifyResult{OutputPath: req.OutputPath}, &daemon.VerificationUnknownError{
		Command:        req.Command,
		DiagnosticPath: req.OutputPath,
		Err:            verifier.cause,
	}
}

func TestRunSettleActiveRunOnSameWorkingTreeBlocks(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusFailed)}})
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	active, err := runStore.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "ma/other",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "ma/widget-flow",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	taskPath := implementTaskPath(repoDir, "task_01")
	beforeTask := mustRead(t, taskPath)
	beforeHeadCount := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-list", "--count", "HEAD"))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected exit code 2, got %d (stderr %q)", code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	for _, expected := range []string{active.ID, "working tree", "roundfix stop " + active.ID} {
		if !strings.Contains(stderr.String(), expected) {
			t.Fatalf("expected stderr to contain %q, got %q", expected, stderr.String())
		}
	}
	if got := mustRead(t, taskPath); got != beforeTask {
		t.Fatalf("expected task file unchanged")
	}
	if got := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-list", "--count", "HEAD")); got != beforeHeadCount {
		t.Fatalf("expected HEAD count %s, got %s", beforeHeadCount, got)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database reader: %v", err)
	}
	defer func() { _ = reader.Close() }()
	events, err := reader.RunEventsAfter(context.Background(), active.ID, 0, 10)
	if err != nil {
		t.Fatalf("read Run Events: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected settle to write no Run Event Journal entries, got %#v", events)
	}
	stillActive, found, err := reader.Run(context.Background(), active.ID)
	if err != nil || !found {
		t.Fatalf("read legacy active Run: found=%v err=%v", found, err)
	}
	if stillActive.State != store.StateActive {
		t.Fatalf("expected pid-less legacy Run to stay Active, got %#v", stillActive)
	}
}

func TestRunSettleHelpDocumentsContract(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--help"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	for _, expected := range []string{
		"roundfix settle --spec <slug> --task <task_id>",
		"--spec",
		"--task",
		"Exit codes:",
		"all Run Worktree changes plus",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("expected settle help to contain %q, got:\n%s", expected, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = runCLIContext(t, context.Background(), []string{"--help"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("expected top-level help exit 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "roundfix settle --spec <slug> --task <task_id>") {
		t.Fatalf("expected top-level help to list settle, got:\n%s", stdout.String())
	}
}

func gitSettleOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmdArgs := append(gitConfigArgsForTest(), args...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Dir = dir
	cmd.Env = isolatedGitEnvForTest()
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output.String())
	}
	return output.String()
}

func settleCommitFiles(t *testing.T, repoDir string) []string {
	t.Helper()
	output := gitSettleOutput(t, repoDir, "diff-tree", "--no-commit-id", "--name-only", "-r", "HEAD")
	var files []string
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}
	return files
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func assertContainsString(t *testing.T, values []string, expected string) {
	t.Helper()
	for _, value := range values {
		if value == expected {
			return
		}
	}
	t.Fatalf("expected %q in %v", expected, values)
}

func assertSettleSurfaceLine(t *testing.T, stderr string, surface string) {
	t.Helper()
	want := "Settle surface: " + surface + "\n"
	if !strings.Contains(stderr, want) {
		t.Fatalf("expected stderr to contain %q, got %q", want, stderr)
	}
}

func assertOnlySettleSurfaceLine(t *testing.T, stderr string, surface string) {
	t.Helper()
	want := "Settle surface: " + surface + "\n"
	if stderr != want {
		t.Fatalf("expected stderr %q, got %q", want, stderr)
	}
}

func configureSettleWorktreeLocation(t *testing.T, repoDir string, location string) string {
	t.Helper()
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), "worktree:\n  location: "+filepath.ToSlash(location)+"\n")
	gitImplement(t, repoDir, "add", ".roundfixrc.yml")
	gitImplement(t, repoDir, "commit", "-m", "configure worktree location")
	return location
}

func createImplementRunWorktreeFixture(t *testing.T, homeDir string, repoDir string, location string, specSlug string, taskID string, terminalState string) (store.Run, runworktree.Ref, runworktree.TaskRef) {
	t.Helper()
	ctx := context.Background()
	head := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", "HEAD"))
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close Run Database: %v", err)
		}
	}()
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "ma/widget-flow",
		HeadSHA:     head,
		SpecSlug:    specSlug,
		Agent:       "codex",
		OwnerPID:    os.Getpid(),
	})
	if err != nil {
		t.Fatalf("create implement Run: %v", err)
	}
	runRef, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: repoDir,
		Location: location,
		RunID:    run.ID,
		HeadSHA:  head,
	})
	if err != nil {
		t.Fatalf("create Run Worktree: %v", err)
	}
	runWorktreePath, err := filepath.EvalSymlinks(runRef.Path)
	if err != nil {
		t.Fatalf("resolve Run Worktree path: %v", err)
	}
	runRef.Path = runWorktreePath
	run, err = runStore.SetRunWorkDir(ctx, run.ID, runRef.Path)
	if err != nil {
		t.Fatalf("record Run Worktree: %v", err)
	}
	var taskRef runworktree.TaskRef
	if strings.TrimSpace(taskID) != "" {
		taskRef, err = runworktree.CreateTask(ctx, runRef, taskID, nil)
		if err != nil {
			t.Fatalf("create Task Worktree: %v", err)
		}
	}
	if terminalState != "" {
		completed, completeErr := runStore.CompleteRun(ctx, run.ID, terminalState)
		err = completeErr
		if err != nil {
			t.Fatalf("complete Run %s: %v", terminalState, err)
		}
		run = completed.Run
	}
	return run, runRef, taskRef
}
