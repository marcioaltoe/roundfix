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

	"roundfix/internal/daemon"
	"roundfix/internal/spec"
	"roundfix/internal/store"
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
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected exit code 0, got %d (stderr %q)", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	shortSHA := strings.TrimSpace(gitSettleOutput(t, repoDir, "rev-parse", "--short", "HEAD"))
	expectedStdout := "verify test -f done.txt — ok\n" +
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
	assertContainsString(t, changed, filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "task_01.md")))
	assertNoRunDatabase(t, homeDir)
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
	expectedStdout := "verify test -f done.txt — ok\n" +
		"verify test -f missing.txt — failed\n" +
		"task_01 stays failed — verification failed\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("expected stdout:\n%q\ngot:\n%q", expectedStdout, stdout.String())
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
		"all current worktree changes plus the task file",
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
