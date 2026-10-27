package cli

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestSettleAdvancesTheRunWorktreeToTheBranchPlan(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{implementQAGateSeed(string(spec.StatusFailed))})
	prd := filepath.Join("docs", "specs", implementTestSlug, "_prd.md")
	mustWrite(t, filepath.Join(repo, prd), mustRead(t, filepath.Join(repo, prd))+"\n## Unreachable Acceptance\n\n- criterion: Success Metric 1\n  reason: no remote\n  satisfied-by: task_qa\n")
	gitImplement(t, repo, "add", prd)
	gitImplement(t, repo, "commit", "-m", "declare first criterion")
	location := configureSettleWorktreeLocation(t, repo, filepath.Join(home, "worktrees"))
	run, _, _ := createImplementRunWorktreeFixture(t, home, repo, location, implementTestSlug, "", store.StateUnresolved)
	reportDir := filepath.Join(run.WorkDir, "docs", "specs", implementTestSlug, "qa")
	mustMkdir(t, reportDir)
	mustWrite(t, filepath.Join(reportDir, "qa-report-2026-09-21.md"), "---\nverdict: partial\nrows_blocked_declared: 2\n---\n\n# QA Report\n\n## Results\n\n| ID | Status |\n| --- | --- |\n| Q1 | blocked (declared: Success Metric 1) |\n| Q2 | blocked (declared: Success Metric 2) |\n")
	mustWrite(t, filepath.Join(repo, prd), strings.Replace(mustRead(t, filepath.Join(repo, prd)), "  reason:", "  criterion: Success Metric 2\n  reason:", 1))
	gitImplement(t, repo, "add", prd)
	gitImplement(t, repo, "commit", "-m", "declare second criterion")
	var out, diag bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_qa"}, &out, &diag)
	if code != exitOK {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &diag)
	}
	advanced, removed := strings.Index(diag.String(), "Run Worktree advanced: "), strings.Index(diag.String(), "Run Worktree removed: ")
	if advanced < 0 || removed <= advanced {
		t.Fatalf("advance/removal order: %s", &diag)
	}
	if !strings.Contains(mustRead(t, implementTaskPath(repo, "task_qa")), "status: completed") {
		t.Fatal("user branch lacks settled Task")
	}
}

func TestSettleRefusesWhenTheBranchChangeOverlapsTheSurface(t *testing.T) {
	t.Parallel()
	checkSettlePlanDrift(t, "overlap")
}

func TestSettleRefusesASpecChangedOnADivergedBranch(t *testing.T) {
	t.Parallel()
	checkSettlePlanDrift(t, "diverged")
}

func TestSettleRefusesASpecChangedAfterATaskWorktreeForked(t *testing.T) {
	t.Parallel()
	checkSettlePlanDrift(t, "task")
}

func TestSettleKeepsTheAdvanceWhenVerificationFails(t *testing.T) {
	t.Parallel()
	checkSettlePlanDrift(t, "verification")
}

func checkSettlePlanDrift(t *testing.T, scenario string) {
	t.Helper()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", taskType: "backend", status: string(spec.StatusFailed), verification: []string{"test -f absent.txt"}}})
	location := configureSettleWorktreeLocation(t, repo, filepath.Join(home, "worktrees"))
	taskID := ""
	if scenario == "task" {
		taskID = "task_01"
	}
	run, _, taskRef := createImplementRunWorktreeFixture(t, home, repo, location, implementTestSlug, taskID, store.StateUnresolved)
	surface := run.WorkDir
	if scenario == "task" {
		surface = taskRef.Path
	}
	changed := filepath.Join("docs", "specs", implementTestSlug, "_prd.md")
	if scenario == "overlap" {
		changed = filepath.Join("docs", "specs", implementTestSlug, "task_01.md")
		mustWrite(t, filepath.Join(surface, changed), mustRead(t, filepath.Join(surface, changed))+"\nlocal work\n")
	}
	if scenario == "diverged" {
		mustWrite(t, filepath.Join(surface, "run.txt"), "Run work\n")
		gitImplement(t, surface, "add", "run.txt")
		gitImplement(t, surface, "commit", "-m", "Run commit")
	}
	mustWrite(t, filepath.Join(repo, changed), mustRead(t, filepath.Join(repo, changed))+"\nbranch amendment\n")
	gitImplement(t, repo, "add", changed)
	gitImplement(t, repo, "commit", "-m", "amend branch plan")
	beforeHead := gitSettleOutput(t, surface, "rev-parse", "HEAD")
	beforeStatus := gitSettleOutput(t, surface, "status", "--porcelain=v1")
	beforeTask := mustRead(t, implementTaskPath(surface, "task_01"))
	beforeFile := mustRead(t, filepath.Join(surface, changed))
	var out, diag bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &out, &diag)
	if scenario == "verification" {
		if code != exitRunFailed || !strings.Contains(diag.String(), "Run Worktree advanced: ") {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &diag)
		}
		if got := gitSettleOutput(t, surface, "rev-parse", "HEAD"); got != gitSettleOutput(t, repo, "rev-parse", "HEAD") {
			t.Fatal("advance lost on verification failure")
		}
	} else {
		if code != exitPreflight || !strings.Contains(diag.String(), filepath.ToSlash(changed)) {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &diag)
		}
		if got := gitSettleOutput(t, surface, "rev-parse", "HEAD"); got != beforeHead {
			t.Fatal("refusal moved HEAD")
		}
		if got := mustRead(t, filepath.Join(surface, changed)); got != beforeFile {
			t.Fatal("refusal changed file")
		}
	}
	if got := mustRead(t, implementTaskPath(surface, "task_01")); got != beforeTask {
		t.Fatal("Task changed")
	}
	if got := gitSettleOutput(t, surface, "status", "--porcelain=v1"); got != beforeStatus {
		t.Fatalf("status changed: %s -> %s", beforeStatus, got)
	}
	assertRunWorktreeExists(t, surface)
}
