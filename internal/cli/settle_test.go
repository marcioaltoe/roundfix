package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/daemon"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func TestRunSettleCommitsFailedTaskWorktreeWithDaemonMessage(t *testing.T) {
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

	code := RunContext(context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

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
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Recover the completed work", status: string(spec.StatusFailed), verification: []string{"true"}},
		{id: "task_02", title: "Other failed work", status: string(spec.StatusFailed)},
		{id: "task_03", title: "Completed sibling", status: string(spec.StatusCompleted)},
	})
	mustWrite(t, filepath.Join(repoDir, "done.txt"), "preserved work\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

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

	code := RunContext(context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

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

	code := RunContext(context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

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

func TestRunSettleRetargetsKeptRunWorktreeAndCleansUpAfterIntegration(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Recover kept work",
			taskType:     "backend",
			verification: []string{"test -f done.txt"},
		},
	})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusFailed},
		onTask: func(req agent.ExecuteRequest, _ string) error {
			return os.WriteFile(filepath.Join(req.GitRoot, "done.txt"), []byte("preserved work\n"), 0o644)
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

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

	code = RunContext(context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

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

	code := RunContext(context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

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

	code := RunContext(context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

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

	code := RunContext(context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

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

func TestRunSettleRefusalEnumeratesCandidateStatuses(t *testing.T) {
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

	code := RunContext(context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected preflight exit 2, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	for _, expected := range []string{
		"no failed settle surface",
		missingTaskRef.Path + ": path does not exist",
		runRef.Path + ": status pending",
		repoDir + ": status completed",
		"status is pending; run the Implement Command",
	} {
		if !strings.Contains(stderr.String(), expected) {
			t.Fatalf("expected stderr to contain %q, got %q", expected, stderr.String())
		}
	}
}

func TestRunSettleRequiresSpecAndTask(t *testing.T) {
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
			t.Setenv("HOME", homeDir)
			t.Chdir(t.TempDir())
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

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
				t.Setenv("HOME", homeDir)
				t.Chdir(repoDir)
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
			name: "completed task has nothing to do",
			setup: func(t *testing.T) (string, string, string) {
				homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
				return homeDir, repoDir, implementTaskPath(repoDir, "task_01")
			},
			args:    []string{"settle", "--spec", implementTestSlug, "--task", "task_01"},
			message: "status is completed; nothing to do",
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

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

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

	code := RunContext(context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

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

func TestRunSettleActiveRunOnSameWorkingTreeBlocks(t *testing.T) {
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

	code := RunContext(context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

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
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"settle", "--help"}, &stdout, &stderr)

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
	code = RunContext(context.Background(), []string{"--help"}, &stdout, &stderr)
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
