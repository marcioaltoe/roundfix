// Suite: Delivery Queue retry across every terminal Implement Run of an item.
// Invariant: retry considers the item's Runs newest first and carries each Run's remaining proved Tasks.
// Boundary IN: commandDeliveryWorkflow, the real Run Database, and real local Git worktrees.
// Boundary OUT: the Delivery Engine transition is covered in internal/delivery/retry_test.go.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func TestItemRecoveryCarriesTheNewestRunWhenTheItemRecordsAnOlderOne(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}, {id: "task_02"}})
	workflow := newItemRecoveryWorkflowForRepository(t, homeDir, repoDir)
	older := createDeliveryRetryRun(t, workflow, repoDir, "ma/widget-flow", []implementSeed{{id: "task_01"}}, store.StateBudgetExceeded)
	completeCarryForwardTaskOnCheckout(t, repoDir, "task_01")
	newer := createDeliveryRetryRun(t, workflow, repoDir, "ma/widget-flow", []implementSeed{{id: "task_02"}}, store.StateBudgetExceeded)
	selected, err := workflow.deliveryCarryForwardRuns(t.Context(), implementTestSlug, "ma/widget-flow", older.run.ID)
	if err != nil {
		t.Fatalf("select every item Run: %v", err)
	}
	selectedIDs := make([]string, 0, len(selected))
	for _, run := range selected {
		selectedIDs = append(selectedIDs, run.ID)
	}
	if want := []string{newer.run.ID, older.run.ID}; !slices.Equal(selectedIDs, want) {
		t.Fatalf("selected Runs = %v, want newest-first deduplicated %v", selectedIDs, want)
	}

	result, err := workflow.CarryForward(t.Context(), repoDir, implementTestSlug, "ma/widget-flow", older.run.ID)
	if err != nil {
		t.Fatalf("carry every item Run from stale recorded Run: %v", err)
	}
	wantRuns := []delivery.CarriedRun{{RunID: newer.run.ID, Carried: []string{"task_02"}}}
	if result.RunID != newer.run.ID || !slices.EqualFunc(result.Runs, wantRuns, equalCarriedRun) {
		t.Fatalf("carry-forward result = %+v, want newest Run %q carrying task_02", result, newer.run.ID)
	}
}

func TestItemRecoveryCarriesEveryRunNewestFirst(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}, {id: "task_02"}})
	workflow := newItemRecoveryWorkflowForRepository(t, homeDir, repoDir)
	older := createDeliveryRetryRun(t, workflow, repoDir, "ma/widget-flow", []implementSeed{{id: "task_01"}}, store.StateBudgetExceeded)
	newer := createDeliveryRetryRun(t, workflow, repoDir, "ma/widget-flow", []implementSeed{{id: "task_02"}}, store.StateBudgetExceeded)

	result, err := workflow.CarryForward(t.Context(), repoDir, implementTestSlug, "ma/widget-flow", older.run.ID)
	if err != nil {
		t.Fatalf("carry every item Run: %v", err)
	}
	wantRuns := []delivery.CarriedRun{
		{RunID: newer.run.ID, Carried: []string{"task_02"}},
		{RunID: older.run.ID, Carried: []string{"task_01"}},
	}
	if result.RunID != newer.run.ID || !slices.EqualFunc(result.Runs, wantRuns, equalCarriedRun) {
		t.Fatalf("carry-forward result = %+v, want %+v", result, wantRuns)
	}
}

func TestItemRecoverySkipsACompletedRunWithAGoneWorktree(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	workflow := newItemRecoveryWorkflowForRepository(t, homeDir, repoDir)
	completed := createDeliveryRetryRun(t, workflow, repoDir, "ma/widget-flow", []implementSeed{{id: "task_01"}}, store.StateBudgetExceeded)
	completeCarryForwardTaskOnCheckout(t, repoDir, "task_01")
	gitImplement(t, repoDir, "worktree", "remove", "--force", completed.ref.Path)

	result, err := workflow.CarryForward(t.Context(), repoDir, implementTestSlug, "ma/widget-flow", completed.run.ID)
	if err != nil {
		t.Fatalf("skip completed Run with gone Worktree: %v", err)
	}
	if result.RunID != completed.run.ID || len(result.Runs) != 0 {
		t.Fatalf("completed Run result = %+v, want named Run with nothing carried", result)
	}
}

func TestItemRecoveryRefusalNamesTheRunsAlreadyCarried(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}, {id: "task_02"}})
	workflow := newItemRecoveryWorkflowForRepository(t, homeDir, repoDir)
	older := createDeliveryRetryRun(t, workflow, repoDir, "ma/widget-flow", []implementSeed{{id: "task_01"}}, store.StateBudgetExceeded)
	prdPath := filepath.Join(repoDir, "docs", "specs", implementTestSlug, "_prd.md")
	mustWrite(t, prdPath, mustRead(t, prdPath)+"\nMoved on the item branch.\n")
	gitImplement(t, repoDir, "add", filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "_prd.md")))
	gitImplement(t, repoDir, "commit", "-m", "move older Run input")
	newer := createDeliveryRetryRun(t, workflow, repoDir, "ma/widget-flow", []implementSeed{{id: "task_02"}}, store.StateBudgetExceeded)

	_, err := workflow.CarryForward(t.Context(), repoDir, implementTestSlug, "ma/widget-flow", older.run.ID)
	assertItemRecoveryRefusal(t, err, older.run.ID, repoDir)
	message := err.Error()
	newerIndex := strings.Index(message, newer.run.ID)
	reasonIndex := strings.Index(message, "declared input(s) moved")
	if newerIndex < 0 || reasonIndex < 0 || newerIndex > reasonIndex || !strings.Contains(message, older.run.ID) {
		t.Fatalf("partial carry refusal = %q, want newer and refusing Runs named before the reason", message)
	}
	if got := mustRead(t, filepath.Join(repoDir, "src", "task_02.txt")); got != "task_02 settled\n" {
		t.Fatalf("newer Run implementation = %q, want carried", got)
	}
	if _, statErr := os.Stat(filepath.Join(repoDir, "src", "task_01.txt")); !os.IsNotExist(statErr) {
		t.Fatalf("older refusing Run implementation stat error = %v, want absent", statErr)
	}
}

func TestDeliverRetryPrintsOneLinePerCarriedRun(t *testing.T) {
	var stdout strings.Builder
	printDeliverRetryResult(&stdout, implementTestSlug, delivery.RetryResult{
		Blocker: delivery.BlockerRunBudgetExceeded,
		Stage:   store.DeliveryStageRunning,
		CarriedFrom: delivery.CarryForwardResult{
			RunID: "run_newer",
			Runs: []delivery.CarriedRun{
				{RunID: "run_newer", Carried: []string{"task_02", "task_03"}},
				{RunID: "run_older", Carried: []string{"task_01"}},
			},
		},
	})
	want := "Carried forward from Run run_newer: task_02, task_03\n" +
		"Carried forward from Run run_older: task_01\n" +
		"Retried " + implementTestSlug + ": run-budget-exceeded -> running\n"
	if stdout.String() != want {
		t.Fatalf("retry output = %q, want %q", stdout.String(), want)
	}
}

func equalCarriedRun(left, right delivery.CarriedRun) bool {
	return left.RunID == right.RunID && slices.Equal(left.Carried, right.Carried)
}

func createDeliveryRetryRun(
	t *testing.T,
	workflow *commandDeliveryWorkflow,
	repository string,
	branch string,
	seeds []implementSeed,
	state string,
) carryForwardFixture {
	t.Helper()
	ctx := t.Context()
	head := itemRecoveryHead(t, repository)
	run, err := workflow.store.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repository,
		LocalBranch: branch,
		HeadSHA:     head,
		SpecSlug:    implementTestSlug,
		Agent:       "codex",
	})
	if err != nil {
		t.Fatalf("create delivery retry Run: %v", err)
	}
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: repository,
		Location: t.TempDir(),
		RunID:    run.ID,
		HeadSHA:  head,
	})
	if err != nil {
		t.Fatalf("create delivery retry Run Worktree: %v", err)
	}
	ref.Path, err = filepath.EvalSymlinks(ref.Path)
	if err != nil {
		t.Fatalf("resolve delivery retry Run Worktree: %v", err)
	}
	run, err = workflow.store.SetRunWorkDir(ctx, run.ID, ref.Path)
	if err != nil {
		t.Fatalf("record delivery retry Run Worktree: %v", err)
	}
	commits := make(map[string]string, len(seeds))
	for _, seed := range seeds {
		implementationPath := filepath.Join("src", seed.id+".txt")
		mustMkdir(t, filepath.Join(ref.Path, filepath.Dir(implementationPath)))
		mustWrite(t, filepath.Join(ref.Path, implementationPath), seed.id+" settled\n")
		if err := spec.SetStatus(implementTaskPath(ref.Path, seed.id), spec.StatusCompleted); err != nil {
			t.Fatalf("settle source %s: %v", seed.id, err)
		}
		gitImplement(t, ref.Path, "add", implementationPath, filepath.Join("docs", "specs", implementTestSlug, seed.id+".md"))
		gitImplement(
			t,
			ref.Path,
			"commit",
			"-m", fmt.Sprintf("feat: settle %s", seed.id),
			"-m", fmt.Sprintf("Roundfix-Spec: %s\nRoundfix-Task: %s", implementTestSlug, seed.id),
		)
		commits[seed.id] = itemRecoveryHead(t, ref.Path)
		if _, err := workflow.store.AppendRunEvents(ctx, []runevent.RunEvent{
			{
				RunID:       run.ID,
				Source:      runevent.SourceDaemon,
				Kind:        runevent.KindDaemonVerification,
				ReviewIssue: seed.id,
				Payload:     []byte(fmt.Sprintf(`{"attempt":1,"phase":"verdict","task":%q,"verdict":"passed"}`, seed.id)),
			},
			{
				RunID:       run.ID,
				Source:      runevent.SourceDaemon,
				Kind:        runevent.KindDaemonTask,
				ReviewIssue: seed.id,
				Payload:     []byte(fmt.Sprintf(`{"task":%q,"phase":"settled","status":"completed"}`, seed.id)),
			},
		}); err != nil {
			t.Fatalf("append delivery retry evidence for %s: %v", seed.id, err)
		}
	}
	completed, err := workflow.store.CompleteRun(ctx, run.ID, state)
	if err != nil {
		t.Fatalf("complete delivery retry Run: %v", err)
	}
	return carryForwardFixture{repoDir: repository, run: completed.Run, ref: ref, commits: commits}
}
