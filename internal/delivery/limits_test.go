// Suite: Delivery Queue limit enforcement.
// Invariant: the engine starts no queued item at its deadline and performs no retry work beyond queue limits.
// Boundary IN: the real SQLite Delivery Queue store and Delivery Engine orchestration.
// Boundary OUT: Git worktrees, recovery, and GitHub actions are represented by existing engine fakes.
package delivery

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/store"
)

func TestAQueuedItemParksAtTheDeadlineWithoutAWorktree(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	deadline := time.Date(2026, time.September, 29, 12, 0, 0, 0, time.UTC)
	const (
		gitRoot  = "/repo-deadline"
		specSlug = "deadline-spec"
	)
	if _, err := runStore.CreateDeliveryQueueWithLimits(ctx, gitRoot, []string{specSlug}, store.DeliveryQueueLimits{Deadline: deadline}); err != nil {
		t.Fatalf("create deadline-limited Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	clock := &fakeDeliveryClock{now: deadline}
	engine := newTestDeliveryEngineWithWait(runStore, workflow, newFakeDeliveryBoundary(), clock, clock)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine at deadline: %v", err)
	}
	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Stage != store.DeliveryStageParked || item.Blocker != BlockerQueueDeadline {
		t.Fatalf("item at deadline = %+v, want parked as %q", item, BlockerQueueDeadline)
	}
	if item.Branch != "" || item.Worktree != "" || item.WorktreeProvisioned {
		t.Fatalf("item at deadline has worktree state: %+v", item)
	}
	if len(workflow.events[specSlug]) != 0 {
		t.Fatalf("item at deadline ran workspace or delivery actions: %v", workflow.events[specSlug])
	}
}

func TestAQueuedItemStartsBeforeTheDeadline(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	deadline := time.Date(2026, time.September, 29, 12, 0, 0, 0, time.UTC)
	const (
		gitRoot  = "/repo-before-deadline"
		specSlug = "before-deadline"
	)
	if _, err := runStore.CreateDeliveryQueueWithLimits(ctx, gitRoot, []string{specSlug}, store.DeliveryQueueLimits{Deadline: deadline}); err != nil {
		t.Fatalf("create deadline-limited Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	clock := &fakeDeliveryClock{now: deadline.Add(-time.Second)}
	engine := newTestDeliveryEngineWithWait(runStore, workflow, newFakeDeliveryBoundary(), clock, clock)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine before deadline: %v", err)
	}
	if events := workflow.events[specSlug]; len(events) == 0 || events[0] != "create-branch" {
		t.Fatalf("events before deadline = %v, want create-branch first", events)
	}
	if item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]; item.Stage != store.DeliveryStageMerged {
		t.Fatalf("item before deadline stage = %q, want merged", item.Stage)
	}
}

func TestAnItemPastQueuedAdvancesAfterTheDeadline(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	deadline := time.Date(2026, time.September, 29, 12, 0, 0, 0, time.UTC)
	const (
		gitRoot  = "/repo-started-before-deadline"
		specSlug = "already-started"
	)
	queue, err := runStore.CreateDeliveryQueueWithLimits(ctx, gitRoot, []string{specSlug}, store.DeliveryQueueLimits{Deadline: deadline})
	if err != nil {
		t.Fatalf("create deadline-limited Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageRunning
	item.Branch = "roundfix/deliver-" + specSlug
	item.Worktree = "/worktrees/" + specSlug
	recordDeliveryItemWorkspace(t, ctx, runStore, gitRoot, &item)
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("seed started Delivery Queue item: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	clock := &fakeDeliveryClock{now: deadline.Add(time.Hour)}
	engine := newTestDeliveryEngineWithWait(runStore, workflow, newFakeDeliveryBoundary(), clock, clock)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("advance started item after deadline: %v", err)
	}
	if events := workflow.events[specSlug]; len(events) == 0 || events[0] != "run" {
		t.Fatalf("started item events after deadline = %v, want run first", events)
	}
	if got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]; got.Stage != store.DeliveryStageMerged {
		t.Fatalf("started item after deadline stage = %q, want merged", got.Stage)
	}
}

func TestRetryRefusesAQueueDeadlineItem(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	deadline := time.Date(2026, time.September, 29, 12, 0, 0, 0, time.UTC)
	const (
		gitRoot  = "/repo-retry-deadline"
		specSlug = "retry-deadline"
	)
	queue, err := runStore.CreateDeliveryQueueWithLimits(ctx, gitRoot, []string{specSlug}, store.DeliveryQueueLimits{Deadline: deadline})
	if err != nil {
		t.Fatalf("create deadline-limited Delivery Queue: %v", err)
	}
	before := queue.Items[0]
	before.Stage = store.DeliveryStageParked
	before.Blocker = BlockerQueueDeadline
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, before); err != nil {
		t.Fatalf("seed deadline-blocked Delivery Queue item: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	workspace := &countingRetryWorkspace{fakeDeliveryWorkflow: workflow}
	recovery := &fakeItemRecovery{states: []ItemState{{Head: "unused"}}}
	engine := newRetryDeliveryEngine(runStore, workflow, workspace, recovery, newFakeDeliveryBoundary())

	_, err = engine.Retry(ctx, gitRoot, specSlug)
	if err == nil || !strings.Contains(err.Error(), deadline.Format(time.RFC3339)) || !strings.Contains(err.Error(), "roundfix deliver start") {
		t.Fatalf("deadline retry error = %v, want deadline and roundfix deliver start", err)
	}
	if workspace.useCalls != 0 || len(recovery.events) != 0 {
		t.Fatalf("deadline retry performed work: use=%d recovery=%v", workspace.useCalls, recovery.events)
	}
	if after := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]; !reflect.DeepEqual(after, before) {
		t.Fatalf("deadline item changed after refused retry:\n got: %#v\nwant: %#v", after, before)
	}
}

func TestRetryRefusesAnItemAtItsRetryLimitBeforeCarryForward(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const (
		gitRoot  = "/repo-engine-retry-limit"
		specSlug = "engine-retry-limit"
	)
	queue, err := runStore.CreateDeliveryQueueWithLimits(ctx, gitRoot, []string{specSlug}, store.DeliveryQueueLimits{MaxRetries: 1})
	if err != nil {
		t.Fatalf("create retry-limited Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageParked
	item.Blocker = BlockerRunUnresolved
	item.Branch = "roundfix/deliver-" + specSlug
	item.Worktree = "/worktrees/" + specSlug
	item.RunID = "run-original"
	item.CandidateCommits = []string{"candidate-original"}
	recordDeliveryItemWorkspace(t, ctx, runStore, gitRoot, &item)
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("seed Delivery Queue item for permitted retry: %v", err)
	}
	target := item
	target.Stage = store.DeliveryStageRunning
	if _, _, err := runStore.RetryDeliveryQueueItem(ctx, gitRoot, target, item.Blocker); err != nil {
		t.Fatalf("use permitted retry: %v", err)
	}
	item = readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	item.Stage = store.DeliveryStageParked
	item.Blocker = BlockerRunUnresolved
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("park Delivery Queue item at retry limit: %v", err)
	}
	before := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	workflow := newFakeDeliveryWorkflow()
	workspace := &countingRetryWorkspace{fakeDeliveryWorkflow: workflow}
	recovery := &fakeItemRecovery{states: []ItemState{{Head: "candidate-original"}}}
	engine := newRetryDeliveryEngine(runStore, workflow, workspace, recovery, newFakeDeliveryBoundary())

	_, err = engine.Retry(ctx, gitRoot, specSlug)
	if !errors.Is(err, store.ErrDeliveryRetryLimit) {
		t.Fatalf("engine retry-limit error = %v, want wrapped ErrDeliveryRetryLimit", err)
	}
	if workspace.useCalls != 0 || len(recovery.events) != 0 {
		t.Fatalf("retry at limit performed work: use=%d recovery=%v", workspace.useCalls, recovery.events)
	}
	if after := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]; !reflect.DeepEqual(after, before) {
		t.Fatalf("retry-limited item changed after refusal:\n got: %#v\nwant: %#v", after, before)
	}
}

type countingRetryWorkspace struct {
	*fakeDeliveryWorkflow
	useCalls int
}

func (workspace *countingRetryWorkspace) UseItemBranch(
	ctx context.Context,
	gitRoot string,
	specSlug string,
	branch string,
	worktree string,
	provisioned bool,
) (string, error) {
	workspace.useCalls++
	return workspace.fakeDeliveryWorkflow.UseItemBranch(ctx, gitRoot, specSlug, branch, worktree, provisioned)
}

var _ ItemWorkspace = (*countingRetryWorkspace)(nil)
