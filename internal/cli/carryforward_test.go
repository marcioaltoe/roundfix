// Suite: Spec-scoped Task Carry-Forward inspection.
// Invariant: a Spec query reports only carriable terminal Runs in its repository and skips released Run Worktrees.
// Boundary IN: Run Database listings and events, local Git surfaces, and the existing Task Carry-Forward proofs.
// Boundary OUT: Implement Command Preflight selection and reporting, owned by internal/cli/implement_test.go.
package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func TestInspectSpecCarryForwards(t *testing.T) {
	t.Parallel()

	t.Run("a Spec with no prior Runs reports no carry-forwards", func(t *testing.T) {
		t.Parallel()
		homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", title: "Build the core"}})
		runStore := openCarryForwardStore(t, homeDir)

		results, err := inspectSpecCarryForwards(
			context.Background(),
			runStore,
			repoDir,
			builtinCarryForwardSpecsRoot(repoDir),
			implementTestSlug,
		)
		if err != nil {
			t.Fatalf("inspect Spec carry-forwards: %v", err)
		}
		if len(results) != 0 {
			t.Fatalf("Spec carry-forwards = %+v, want none", results)
		}
	})

	t.Run("one carriable Run is reported with its proven candidate count", func(t *testing.T) {
		t.Parallel()
		fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01", title: "Build the core"}})
		addCarryForwardDistractors(t, fixture)
		runStore := openCarryForwardStore(t, fixture.homeDir)

		results, err := inspectSpecCarryForwards(
			context.Background(),
			runStore,
			fixture.repoDir,
			builtinCarryForwardSpecsRoot(fixture.repoDir),
			implementTestSlug,
		)
		if err != nil {
			t.Fatalf("inspect Spec carry-forwards: %v", err)
		}
		if len(results) != 1 {
			t.Fatalf("Spec carry-forwards = %+v, want one scoped Run", results)
		}
		result := results[0]
		if result.Run.ID != fixture.run.ID {
			t.Fatalf("reported Run = %q, want %q", result.Run.ID, fixture.run.ID)
		}
		if len(result.Candidates) != 1 || result.Candidates[0].TaskID != "task_01" {
			t.Fatalf("reported candidates = %+v, want task_01", result.Candidates)
		}
		if result.Candidates[0].Action != "would carry forward with --carry-forward" {
			t.Fatalf("candidate action = %q, want carry-forward action", result.Candidates[0].Action)
		}
		if got := result.carriable(); got != 1 {
			t.Fatalf("carriable candidates = %d, want 1", got)
		}
	})

	t.Run("a Run whose recorded Run Worktree is gone is skipped", func(t *testing.T) {
		t.Parallel()
		fixture := newCarryForwardFixture(t, store.StateStopped, []implementSeed{{id: "task_01", title: "Build the core"}})
		if err := runworktree.CleanupClean(context.Background(), fixture.ref); err != nil {
			t.Fatalf("release carry-forward Run Worktree: %v", err)
		}
		if _, err := os.Stat(fixture.ref.Path); !os.IsNotExist(err) {
			t.Fatalf("released Run Worktree stat error = %v, want not exist", err)
		}
		runStore := openCarryForwardStore(t, fixture.homeDir)

		results, err := inspectSpecCarryForwards(
			context.Background(),
			runStore,
			fixture.repoDir,
			builtinCarryForwardSpecsRoot(fixture.repoDir),
			implementTestSlug,
		)
		if err != nil {
			t.Fatalf("inspect Spec carry-forwards: %v", err)
		}
		if len(results) != 0 {
			t.Fatalf("Spec carry-forwards = %+v, want released Run omitted", results)
		}
	})

	t.Run("a Run whose Tasks all refuse reports zero carriable candidates", func(t *testing.T) {
		t.Parallel()
		fixture := newCarryForwardFixture(t, store.StateStopped, []implementSeed{{id: "task_01", title: "Build the core"}})
		prdPath := filepath.Join(fixture.repoDir, "docs", "specs", implementTestSlug, "_prd.md")
		mustWrite(t, prdPath, mustRead(t, prdPath)+"\nMoved after settlement.\n")
		runStore := openCarryForwardStore(t, fixture.homeDir)

		results, err := inspectSpecCarryForwards(
			context.Background(),
			runStore,
			fixture.repoDir,
			builtinCarryForwardSpecsRoot(fixture.repoDir),
			implementTestSlug,
		)
		if err != nil {
			t.Fatalf("inspect Spec carry-forwards: %v", err)
		}
		if len(results) != 1 || len(results[0].Candidates) != 1 {
			t.Fatalf("Spec carry-forwards = %+v, want one refusing candidate", results)
		}
		result := results[0]
		if got := result.carriable(); got != 0 {
			t.Fatalf("carriable candidates = %d, want 0", got)
		}
		if !strings.Contains(result.Candidates[0].RefusalReason, "declared input(s) moved") {
			t.Fatalf("candidate refusal = %q, want moved-input reason", result.Candidates[0].RefusalReason)
		}
	})
}

func TestCarryForwardExplicitRunStillRefusesAMissingRunWorktree(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateStopped, []implementSeed{{id: "task_01", title: "Build the core"}})
	if err := runworktree.CleanupClean(context.Background(), fixture.ref); err != nil {
		t.Fatalf("release carry-forward Run Worktree: %v", err)
	}
	var stdout strings.Builder
	var stderr strings.Builder

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", fixture.run.ID, "--carry-forward", "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitRunFailed {
		t.Fatalf("carry-forward exit = %d, want %d stderr=%q stdout=%q", code, exitRunFailed, stderr.String(), stdout.String())
	}
	if !strings.Contains(stderr.String(), fixture.run.ID) || !strings.Contains(stderr.String(), "load Run") {
		t.Fatalf("carry-forward stderr = %q, want named-Run load failure", stderr.String())
	}
}

func openCarryForwardStore(t *testing.T, homeDir string) *store.Store {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open carry-forward Run Database: %v", err)
	}
	t.Cleanup(func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close carry-forward Run Database: %v", err)
		}
	})
	return runStore
}

func builtinCarryForwardSpecsRoot(repoDir string) roundconfig.SpecsRoot {
	return roundconfig.SpecsRoot{
		Path:        filepath.Join(repoDir, "docs", "specs"),
		BuiltInRoot: true,
	}
}

func addCarryForwardDistractors(t *testing.T, fixture carryForwardFixture) {
	t.Helper()
	runStore := openCarryForwardStore(t, fixture.homeDir)
	for _, seed := range []struct {
		gitRoot  string
		specSlug string
		state    string
	}{
		{gitRoot: fixture.repoDir, specSlug: "other-spec", state: store.StateStopped},
		{gitRoot: fixture.repoDir + "-other", specSlug: implementTestSlug, state: store.StateStopped},
		{gitRoot: fixture.repoDir, specSlug: implementTestSlug, state: store.StateClean},
	} {
		run, err := runStore.CreateRun(context.Background(), store.CreateRunRequest{
			Kind:        store.KindImplement,
			GitRoot:     seed.gitRoot,
			LocalBranch: "ma/widget-flow",
			HeadSHA:     fixture.run.HeadSHA,
			SpecSlug:    seed.specSlug,
			Agent:       "codex",
		})
		if err != nil {
			t.Fatalf("create carry-forward distractor: %v", err)
		}
		if _, err := runStore.CompleteRun(context.Background(), run.ID, seed.state); err != nil {
			t.Fatalf("complete carry-forward distractor: %v", err)
		}
	}
}
