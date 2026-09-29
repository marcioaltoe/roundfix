// Suite: Delivery Queue Pending Question derivation.
// Invariant: one persisted parked item is presented for explicit operator action and no automatic pass changes it.
// Boundary IN: Delivery Queue records and one real Delivery Engine pass.
// Boundary OUT: operator retry and replacement-queue commands, exercised by their command suites.
package delivery

import (
	"context"
	"reflect"
	"testing"
	"time"

	"roundfix/internal/store"
)

func TestPendingQuestionIsTheLowestPositionParkedItem(t *testing.T) {
	t.Parallel()
	queue := store.DeliveryQueue{Items: []store.DeliveryQueueItem{
		{SpecSlug: "later-parked", Position: 3, Stage: store.DeliveryStageParked, Blocker: BlockerReviewStale},
		{SpecSlug: "running", Position: 0, Stage: store.DeliveryStageRunning},
		{SpecSlug: "first-parked", Position: 1, Stage: store.DeliveryStageParked, Blocker: BlockerQueueDeadline},
	}}

	question, found := PendingQuestionFor(queue)

	if !found {
		t.Fatal("PendingQuestionFor did not find a parked item")
	}
	if question.SpecSlug != "first-parked" || question.Blocker != BlockerQueueDeadline || question.Waiting != 1 {
		t.Fatalf("Pending Question = %+v, want lowest-position parked item with one waiting", question)
	}
}

func TestPendingQuestionAnswersEachBlockerClass(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		item store.DeliveryQueueItem
		want string
	}{
		{
			name: "revalidation failure",
			item: store.DeliveryQueueItem{
				SpecSlug: "needs-amendment",
				Stage:    store.DeliveryStageParked,
				Blocker:  BlockerRevalidationFailed + ": SC-REF-UNRESOLVED",
				Worktree: "/worktrees/needs-amendment",
			},
			want: "amend the Spec on its item branch in /worktrees/needs-amendment, then run roundfix deliver retry needs-amendment",
		},
		{
			name: "queue deadline",
			item: store.DeliveryQueueItem{
				SpecSlug: "past-deadline",
				Stage:    store.DeliveryStageParked,
				Blocker:  BlockerQueueDeadline,
			},
			want: "record a new queue for the remaining Specs with roundfix deliver start",
		},
		{
			name: "other blocker",
			item: store.DeliveryQueueItem{
				SpecSlug: "needs-retry",
				Stage:    store.DeliveryStageParked,
				Blocker:  BlockerReviewStale,
			},
			want: "resolve the blocker, then run roundfix deliver retry needs-retry",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			question, found := PendingQuestionFor(store.DeliveryQueue{Items: []store.DeliveryQueueItem{tt.item}})
			if !found || question.Answer != tt.want {
				t.Fatalf("Pending Question = %+v found=%t, want answer %q", question, found, tt.want)
			}
		})
	}
}

func TestNoPendingQuestionWithoutAParkedItem(t *testing.T) {
	t.Parallel()
	queue := store.DeliveryQueue{Items: []store.DeliveryQueueItem{
		{SpecSlug: "queued", Position: 0, Stage: store.DeliveryStageQueued},
		{SpecSlug: "merged", Position: 1, Stage: store.DeliveryStageMerged},
	}}

	question, found := PendingQuestionFor(queue)

	if found || question != (PendingQuestion{}) {
		t.Fatalf("Pending Question = %+v found=%t, want none", question, found)
	}
}

func TestAPendingQuestionSurvivesOwnerPassesAndTime(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	deadline := time.Date(2026, time.September, 29, 12, 0, 0, 0, time.UTC)
	const (
		gitRoot  = "/repo-pending-question"
		specSlug = "pending-question"
	)
	queue, err := runStore.CreateDeliveryQueueWithLimits(
		ctx,
		gitRoot,
		[]string{specSlug},
		store.DeliveryQueueLimits{Deadline: deadline},
	)
	if err != nil {
		t.Fatalf("create deadline-limited Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageParked
	item.Blocker = BlockerReviewStale
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("park Delivery Queue item: %v", err)
	}
	beforeQueue := readDeliveryQueue(t, ctx, runStore, gitRoot)
	before, found := PendingQuestionFor(beforeQueue)
	if !found {
		t.Fatal("parked item did not produce a Pending Question")
	}

	clock := &fakeDeliveryClock{now: deadline.Add(-time.Second)}
	clock.now = deadline.Add(time.Hour)
	engine := newTestDeliveryEngineWithWait(runStore, newFakeDeliveryWorkflow(), newFakeDeliveryBoundary(), clock, clock)
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine owner pass after deadline: %v", err)
	}
	after, found := PendingQuestionFor(readDeliveryQueue(t, ctx, runStore, gitRoot))
	if !found || !reflect.DeepEqual(after, before) {
		t.Fatalf("Pending Question after owner pass = %+v found=%t, want unchanged %+v", after, found, before)
	}
}
