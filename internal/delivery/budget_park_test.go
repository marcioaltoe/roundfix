// Suite: Delivery Queue budget parking.
// Invariant: a budget-ended Run parks with recoverable Run evidence without stopping later items.
// Boundary IN: the real SQLite Delivery Queue store and delivery engine transitions.
// Boundary OUT: the Implement runner, item recovery, Git worktrees, and GitHub are represented by fakes.
package delivery

import (
	"context"
	"errors"
	"testing"

	"roundfix/internal/store"
)

func TestDeliveryParksABudgetExceededRunAsARunOutcome(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const (
		gitRoot  = "/repo-budget-run-id"
		specSlug = "budget-spec"
	)
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{specSlug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.runs[specSlug] = RunResult{RunID: "  run-budget-42  ", Outcome: RunOutcomeBudgetExceeded}
	workflow.recordWorkspace = func(branch, worktree string) error {
		_, _, _, err := runStore.RecordDeliveryQueueItemWorktree(ctx, gitRoot, specSlug, branch, worktree)
		return err
	}
	engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}
	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Stage != store.DeliveryStageParked || item.Blocker != BlockerRunBudgetExceeded || item.RunID != "run-budget-42" {
		t.Fatalf("budget-ended item = %+v, want parked with its trimmed Run ID", item)
	}
}

func TestDeliveryContinuesAfterABudgetParkedItem(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-budget-continues"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"budget-spec", "next-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.runs["budget-spec"] = RunResult{RunID: "run-budget", Outcome: RunOutcomeBudgetExceeded}
	engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}
	items := readDeliveryQueue(t, ctx, runStore, gitRoot).Items
	if items[0].Stage != store.DeliveryStageParked || items[0].Blocker != BlockerRunBudgetExceeded {
		t.Fatalf("budget-ended item = %+v, want parked", items[0])
	}
	if items[1].Stage != store.DeliveryStageMerged {
		t.Fatalf("next item stage = %q, want merged", items[1].Stage)
	}
}

func TestDeliveryRetryResumesABudgetParkedItemFromItsRun(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-budget-retry"
	item := seedParkedRetryItem(t, ctx, runStore, gitRoot, "budget-retry", BlockerRunBudgetExceeded)
	item.RunID = "run-budget-original"
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("record budget Run ID: %v", err)
	}
	recovery := &fakeItemRecovery{
		states:      []ItemState{{Head: "candidate-original"}, {UnfinishedTasks: []string{"task_02"}, Head: "candidate-original"}},
		carryResult: CarryForwardResult{RunID: "run-budget-carried", Carried: []string{"task_01"}},
	}
	engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())

	if _, err := engine.Retry(ctx, gitRoot, item.SpecSlug); err != nil {
		t.Fatalf("retry budget-parked Delivery Queue item: %v", err)
	}
	if recovery.carryRunID != "run-budget-original" {
		t.Fatalf("carry-forward Run ID = %q, want recorded budget Run", recovery.carryRunID)
	}
}

func TestDeliveryStillParksAnExecutorErrorAsADeliveryError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const (
		gitRoot  = "/repo-budget-executor-error"
		specSlug = "executor-error-spec"
	)
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{specSlug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.runErrors[specSlug] = errors.New("executor unavailable")
	engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}
	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	want := BlockerDeliveryError + ": run Implement executor: executor unavailable"
	if item.Stage != store.DeliveryStageParked || item.Blocker != want {
		t.Fatalf("executor-error item = %+v, want blocker %q", item, want)
	}
}
