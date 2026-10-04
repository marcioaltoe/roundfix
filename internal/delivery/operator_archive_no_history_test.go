// Suite: operator-archived refusal without a candidate or History.
// Boundary IN: real queue persistence and Retry refusals.
// Boundary OUT: item inspection and delivery executors use existing fakes.
package delivery

import (
	"fmt"
	"testing"
)

func TestRetryRefusesAnOperatorArchiveWithoutCandidateOrHistory(t *testing.T) {
	for _, blocker := range []string{BlockerRunUnresolved, BlockerQAEnvironmentPartial} {
		t.Run(blocker, func(t *testing.T) {
			ctx := t.Context()
			runStore := openDeliveryEngineStore(t, ctx)
			const root = "/repo-operator-no-history"
			before := seedParkedRetryItem(t, ctx, runStore, root, "operator-no-history", blocker)
			before.CandidateCommits = nil
			if err := runStore.UpdateDeliveryQueueItem(ctx, root, before); err != nil {
				t.Fatal(err)
			}
			before = readDeliveryQueue(t, ctx, runStore, root).Items[0]
			recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, QAOverride: true, Head: "operator-head"}}}
			engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())
			if engine.history != nil {
				t.Fatal("retry engine has History, want none")
			}
			_, err := engine.Retry(ctx, root, before.SpecSlug)
			want := fmt.Sprintf(`retry Delivery Queue item %q: archived item head "operator-head" differs from candidate head ""`, before.SpecSlug)
			if err == nil || err.Error() != want {
				t.Errorf("refusal=%v, want %q", err, want)
			}
			assertRetryItemUnchanged(t, ctx, runStore, root, before)
		})
	}
}
