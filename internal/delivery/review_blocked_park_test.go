// Boundary: Delivery Engine parks and their public Park Class.
package delivery

import (
	"testing"

	"roundfix/internal/store"
)

func TestReviewBlockedParkCarriesTheReviewReason(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, reason, blocker string }{
		{"two lines", "  review unavailable  \nsecond line", "review-blocked: review unavailable"},
		{"empty", "", "review-blocked"},
		{"blank first line", " \t\nsecond line", "review-blocked"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			runStore := openDeliveryEngineStore(t, ctx)
			const root = "/repo-review-reason"
			if _, err := runStore.CreateDeliveryQueue(ctx, root, []string{"review-blocked"}); err != nil {
				t.Fatal(err)
			}
			workflow := newFakeDeliveryWorkflow()
			workflow.reviewResults["review-blocked"] = ReviewResult{Outcome: ReviewOutcomeBlocked, Head: "reviewed-review-blocked", Reason: tc.reason}
			engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())
			if _, err := engine.Run(ctx, root); err != nil {
				t.Fatal(err)
			}
			queue := readDeliveryQueue(t, ctx, runStore, root)
			item := queue.Items[0]
			if item.Stage != store.DeliveryStageParked || item.Blocker != tc.blocker {
				t.Fatalf("item=%+v, want parked with %q", item, tc.blocker)
			}
			if got := ClassifyPark(queue, item).Class; got != ParkClassReview {
				t.Fatalf("class=%q, want %q", got, ParkClassReview)
			}
		})
	}
}
