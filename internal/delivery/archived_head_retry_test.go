// Boundary: real queue persistence and Retry; item inspection and ancestry use
// the existing engine fakes. Corrective parks never inherit ancestry approval.
package delivery

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestRetryOfAnArchivedItemWithACorrectionReturnsToReview(t *testing.T) {
	ctx := t.Context()
	db := openDeliveryEngineStore(t, ctx)
	const root = "/repo-archived-correction"
	before := seedParkedRetryItem(t, ctx, db, root, "correction", BlockerGateFailed)
	before.CandidateCommits = append(before.CandidateCommits, "newest-candidate")
	if err := db.UpdateDeliveryQueueItem(ctx, root, before); err != nil {
		t.Fatal(err)
	}
	recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "correction-head"}}}
	engine := newRetryDeliveryEngine(db, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())
	history := &operatorArchiveHistory{descends: true}
	engine.history = history
	result, err := engine.Retry(ctx, root, before.SpecSlug)
	if err != nil {
		t.Fatal(err)
	}
	got := readDeliveryQueue(t, ctx, db, root).Items[0]
	if result.Stage != store.DeliveryStageReviewing || got.Stage != store.DeliveryStageReviewing || got.Blocker != "" || !reflect.DeepEqual(got.CandidateCommits, []string{"candidate-original", "newest-candidate", "correction-head"}) {
		t.Fatalf("result=%+v item=%+v", result, got)
	}
	if history.ancestor != "newest-candidate" || history.runID != "" || !reflect.DeepEqual(recovery.events, []string{"inspect"}) {
		t.Fatalf("history=%+v recovery=%v", history, recovery.events)
	}
}

func TestRetryOfAnArchivedItemRefusesAHeadThatDoesNotDescend(t *testing.T) {
	for _, tc := range []struct {
		name    string
		history *operatorArchiveHistory
	}{
		{"unrelated", &operatorArchiveHistory{}}, {"no history", nil}, {"unreadable history", &operatorArchiveHistory{err: errors.New("history unavailable")}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			db := openDeliveryEngineStore(t, ctx)
			const root = "/repo-archived-unrelated"
			before := seedParkedRetryItem(t, ctx, db, root, "unrelated", BlockerGateFailed)
			recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "unrelated-head"}}}
			engine := newRetryDeliveryEngine(db, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())
			if tc.history != nil {
				engine.history = tc.history
			}
			_, err := engine.Retry(ctx, root, before.SpecSlug)
			if err == nil || !strings.Contains(err.Error(), `archived item head "unrelated-head" differs from candidate head "candidate-original"`) {
				t.Fatalf("refusal=%v", err)
			}
			assertRetryItemUnchanged(t, ctx, db, root, before)
		})
	}
}

func TestRetryOfACorrectiveSpecParkStillRefusesAMovedHead(t *testing.T) {
	ctx := t.Context()
	db := openDeliveryEngineStore(t, ctx)
	const root = "/repo-corrective-descendant"
	before := seedParkedRetryItem(t, ctx, db, root, "corrective", BlockerCorrectiveSpecRequired+": archived-spec")
	recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "correction-head"}}}
	engine := newRetryDeliveryEngine(db, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())
	history := &operatorArchiveHistory{descends: true}
	engine.history = history
	_, err := engine.Retry(ctx, root, before.SpecSlug)
	if err == nil || !strings.Contains(err.Error(), "author a corrective Spec with its own authorization and QA gate") {
		t.Fatalf("refusal=%v", err)
	}
	if history.ancestor != "" {
		t.Fatalf("corrective retry used ancestry: %+v", history)
	}
	assertRetryItemUnchanged(t, ctx, db, root, before)
}
