// Suite: operator-archived recovery without a candidate across parks.
// Boundary IN: real queue persistence and Retry transitions.
// Boundary OUT: Git ancestry and item inspection use existing fakes.
package delivery

import (
	"fmt"
	"reflect"
	"testing"

	"roundfix/internal/store"
)

func TestRetryResumesAnOperatorArchivedItemWithoutACandidateAtReview(t *testing.T) {
	for _, blocker := range []string{BlockerRunUnresolved, BlockerQAEnvironmentPartial, BlockerRunBudgetExceeded} {
		t.Run(blocker, func(t *testing.T) {
			ctx := t.Context()
			runStore := openDeliveryEngineStore(t, ctx)
			const root = "/repo-operator-any-park"
			item := seedParkedRetryItem(t, ctx, runStore, root, "operator-any-park", blocker)
			item.CandidateCommits = nil
			if err := runStore.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
				t.Fatal(err)
			}
			recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, QAOverride: true, Head: "operator-head"}}}
			engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())
			history := &operatorArchiveHistory{descends: true}
			engine.history = history
			result, err := engine.Retry(ctx, root, item.SpecSlug)
			if err != nil {
				t.Fatal(err)
			}
			got := readDeliveryQueue(t, ctx, runStore, root).Items[0]
			if result.Stage != store.DeliveryStageReviewing || got.Stage != store.DeliveryStageReviewing || got.Blocker != "" || !reflect.DeepEqual(got.CandidateCommits, []string{"operator-head"}) {
				t.Fatalf("result=%+v item=%+v", result, got)
			}
			if history.runID != item.RunID || history.ancestor != "run-start" {
				t.Fatalf("Run ID=%q anchor=%q, want %q and Run start", history.runID, history.ancestor, item.RunID)
			}
			if !reflect.DeepEqual(recovery.events, []string{"inspect"}) {
				t.Fatalf("recovery events=%v, want inspect only", recovery.events)
			}
		})
	}
}

func TestRetryRefusesAnArchiveWithoutCandidateOrOverride(t *testing.T) {
	for _, test := range []struct {
		name     string
		override bool
		descends bool
		refusal  string
	}{
		{name: "no override", descends: true, refusal: "candidate head is missing"},
		{name: "not descended", override: true, refusal: `archived item head "operator-head" differs from candidate head ""`},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()
			runStore := openDeliveryEngineStore(t, ctx)
			const root = "/repo-operator-no-candidate-refusal"
			before := seedParkedRetryItem(t, ctx, runStore, root, "operator-no-candidate-refusal", BlockerRunUnresolved)
			before.CandidateCommits = nil
			if err := runStore.UpdateDeliveryQueueItem(ctx, root, before); err != nil {
				t.Fatal(err)
			}
			before = readDeliveryQueue(t, ctx, runStore, root).Items[0]
			recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, QAOverride: test.override, Head: "operator-head"}}}
			engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())
			engine.history = &operatorArchiveHistory{descends: test.descends}
			_, err := engine.Retry(ctx, root, before.SpecSlug)
			want := fmt.Sprintf("retry Delivery Queue item %q: %s", before.SpecSlug, test.refusal)
			if err == nil || err.Error() != want {
				t.Fatalf("refusal=%v, want %q", err, want)
			}
			assertRetryItemUnchanged(t, ctx, runStore, root, before)
		})
	}
}
