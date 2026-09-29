// Suite: merged Delivery Queue item cleanup.
// Invariant: a merged item releases its Spec Runs before item-branch removal, while cleanup failures remain retryable.
// Boundary IN: the real SQLite Delivery Queue store and delivery engine cleanup sequencing.
// Boundary OUT: Git-backed Run release and item cleanup, represented by the delivery workflow fake.
package delivery

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestEngineReleasesMergedRunsBeforeRemovingTheItemBranch(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-release-order"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"merged-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}

	workflow := newFakeDeliveryWorkflow()
	engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	if got, want := workflow.cleanupEvents, []string{"release-merged-runs", "remove-item-branch"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("merged cleanup order = %v, want %v", got, want)
	}
}

func TestEngineRecordsACleanupWarningWhenAMergedRunIsKept(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-kept-run"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"merged-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}

	workflow := newFakeDeliveryWorkflow()
	workflow.releaseErrors = []error{errors.New(`keep Run "run-kept": commit abc123 is not represented`)}
	engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if !strings.HasPrefix(item.Blocker, cleanupWarningPrefix) ||
		!strings.Contains(item.Blocker, "run-kept") ||
		!strings.Contains(item.Blocker, "not represented") {
		t.Fatalf("cleanup warning = %q, want kept Run and reason", item.Blocker)
	}
	if workflow.removeCalls != 1 {
		t.Fatalf("item branch removal calls = %d, want one after release failure", workflow.removeCalls)
	}
}

func TestEngineClearsTheCleanupWarningAfterACleanRelease(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-clean-release-retry"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"merged-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}

	workflow := newFakeDeliveryWorkflow()
	workflow.releaseErrors = []error{errors.New(`keep Run "run-kept": not represented`)}
	engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine with kept Run: %v", err)
	}
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("retry merged cleanup: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Blocker != "" {
		t.Fatalf("cleanup warning after clean retry = %q, want empty", item.Blocker)
	}
	if workflow.releaseCalls != 2 || workflow.removeCalls != 2 {
		t.Fatalf("cleanup retry calls = release:%d remove:%d, want two each", workflow.releaseCalls, workflow.removeCalls)
	}
}

func TestEngineJoinsRunReleaseAndItemBranchCleanupFailures(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-joined-cleanup-errors"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"merged-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}

	workflow := newFakeDeliveryWorkflow()
	workflow.releaseErrors = []error{errors.New(`keep Run "run-kept": not represented`)}
	workflow.removeErrors = []error{errors.New("item worktree is busy")}
	engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if !strings.Contains(item.Blocker, "run-kept") || !strings.Contains(item.Blocker, "item worktree is busy") {
		t.Fatalf("joined cleanup warning = %q, want both cleanup failures", item.Blocker)
	}
	if workflow.releaseCalls != 1 || workflow.removeCalls != 1 {
		t.Fatalf("cleanup calls = release:%d remove:%d, want one each", workflow.releaseCalls, workflow.removeCalls)
	}
}

func TestEngineReleasesNothingForAnUnmergedItem(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-unmerged"
	queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"parked-spec"})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageParked
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("park Delivery Queue item: %v", err)
	}

	workflow := newFakeDeliveryWorkflow()
	engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}
	if workflow.releaseCalls != 0 || workflow.removeCalls != 0 {
		t.Fatalf("unmerged cleanup calls = release:%d remove:%d, want zero", workflow.releaseCalls, workflow.removeCalls)
	}
}
