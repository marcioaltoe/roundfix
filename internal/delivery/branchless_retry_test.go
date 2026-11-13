// Boundary: real delivery engine and SQLite; workspace and item recovery are fakes.
package delivery

import (
	"testing"

	"roundfix/internal/store"
)

func TestRetryRequeuesAnItemParkedBeforeItHasABranch(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-branchless-retry"
	const specSlug = "branchless-retry"
	const blocker = `delivery-error: read prerequisites: refresh prerequisite default branch "main": git fetch origin main: Operation timed out`
	queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{specSlug})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageParked
	item.Blocker = blocker
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatal(err)
	}
	recovery := &fakeItemRecovery{}
	workflow := newFakeDeliveryWorkflow()
	workflow.currentBranch = "workspace-not-used"
	engine := newRetryDeliveryEngine(runStore, workflow, nil, recovery, newFakeDeliveryBoundary())

	result, err := engine.Retry(ctx, gitRoot, specSlug)
	if err != nil {
		t.Fatalf("retry branchless item: %v", err)
	}
	if result.Stage != store.DeliveryStageQueued || result.SpecSlug != specSlug || result.Blocker != blocker {
		t.Fatalf("retry result = %+v, want queued with the original blocker", result)
	}
	stored := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if stored.Stage != store.DeliveryStageQueued || stored.Blocker != "" || stored.RetryCount != 1 || stored.Branch != "" || stored.Worktree != "" || stored.WorktreeProvisioned {
		t.Fatalf("retried item = %+v, want queued, one retry and no workspace", stored)
	}
	if len(recovery.events) != 0 || workflow.currentBranch != "workspace-not-used" {
		t.Fatalf("retry reached recovery %v or used workspace branch %q", recovery.events, workflow.currentBranch)
	}
}
