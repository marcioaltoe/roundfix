package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestSettleRefusesInsideTheRunWorktree(t *testing.T) {
	t.Parallel()
	checkSettleCheckout(t, false, false, "recorded")
}

func TestSettleRefusesInsideATaskWorktree(t *testing.T) {
	t.Parallel()
	checkSettleCheckout(t, true, false, "recorded")
}

func TestSettleFromTheUserCheckoutStillSettlesTheRunWorktree(t *testing.T) {
	t.Parallel()
	checkSettleCheckout(t, false, true, "recorded")
}

func TestSettleCheckoutRecordAndFallback(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"recorded linked checkout", "missing database", "unknown run"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			checkSettleCheckout(t, true, false, mode)
		})
	}
}

func checkSettleCheckout(t *testing.T, taskWorktree, fromUserCheckout bool, recordMode string) {
	t.Helper()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", taskType: "backend", status: string(spec.StatusFailed), verification: []string{"test -f done.txt"}},
	})
	if recordMode == "recorded linked checkout" {
		checkout := filepath.Join(homeDir, "user-checkout")
		gitImplement(t, repoDir, "worktree", "add", "-b", "fix/recovery", checkout)
		repoDir = checkout
	}
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	taskID := ""
	if taskWorktree {
		taskID = "task_01"
	}
	run, runRef, taskRef := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, taskID, store.StateUnresolved)
	surface := runRef.Path
	if taskWorktree {
		surface = taskRef.Path
	}
	if recordMode == "missing database" {
		if err := os.Remove(store.DatabasePath(homeDir)); err != nil {
			t.Fatal(err)
		}
	} else if recordMode == "unknown run" {
		gitImplement(t, surface, "branch", "-m", store.RunBranchPrefix+"unknown-task_01")
	}
	mustWrite(t, filepath.Join(surface, "done.txt"), "recovered work\n")
	taskPath := implementTaskPath(surface, "task_01")
	beforeTask := mustRead(t, taskPath)
	beforeStatus := gitSettleOutput(t, surface, "status", "--porcelain=v1")
	beforeHead := gitSettleOutput(t, surface, "rev-parse", "HEAD")
	branch := strings.TrimSpace(gitSettleOutput(t, surface, "branch", "--show-current"))
	workDir := surface
	if fromUserCheckout {
		workDir = repoDir
	}
	setCommandWorkDirForTest(t, workDir)
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)
	if fromUserCheckout {
		if code != exitOK {
			t.Fatalf("settle exit = %d, stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
		if !strings.Contains(stderr.String(), "Run Worktree removed: "+run.WorkDir+"\n") {
			t.Fatalf("missing removal: %q", stderr.String())
		}
		assertRunWorktreeRemoved(t, run.WorkDir)
		if got := mustRead(t, filepath.Join(repoDir, "done.txt")); got != "recovered work\n" {
			t.Fatalf("integrated work = %q", got)
		}
		return
	}
	if code != exitPreflight {
		t.Fatalf("settle exit = %d, want 2; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	want := fmt.Sprintf("settle runs from the user checkout, not from %s on branch %s; run it from %s: roundfix settle --spec %s --task task_01", surface, branch, repoDir, implementTestSlug)
	if !strings.Contains(stderr.String(), want) {
		t.Fatalf("stderr = %q, want %q", stderr.String(), want)
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	if got := gitSettleOutput(t, surface, "rev-parse", "HEAD"); got != beforeHead {
		t.Fatalf("HEAD changed from %q to %q", beforeHead, got)
	}
	if got := gitSettleOutput(t, surface, "status", "--porcelain=v1"); got != beforeStatus {
		t.Fatalf("status changed from %q to %q", beforeStatus, got)
	}
	if got := mustRead(t, taskPath); got != beforeTask {
		t.Fatal("Task file changed")
	}
	assertRunWorktreeExists(t, run.WorkDir)
	if recordMode == "missing database" {
		assertNoRunDatabase(t, homeDir)
	}
}
