// Suite: corrective Spec delivery parking and retry.
// Invariant: findings on archived Specs name the required correction and retry never carries that work forward implicitly.
// Boundary IN: the real SQLite Delivery Queue store and delivery engine stage transitions.
// Boundary OUT: review execution, item inspection, Git worktrees and corrective Spec authoring are represented by fakes.
package delivery

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestDeliveryEngineParksFindingsAfterArchiveAsCorrectiveSpecRequired(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const (
		gitRoot  = "/repo-corrective-spec"
		specSlug = "delivery-spec"
	)
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{specSlug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.reviewResults[specSlug] = ReviewResult{
		Outcome:       ReviewOutcomeFindings,
		Head:          "reviewed-" + specSlug,
		Reason:        "finding",
		ArchivedSpecs: []string{"0177-two", "0176-one"},
	}
	engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	wantBlocker := BlockerCorrectiveSpecRequired + ": 0176-one, 0177-two"
	if item.Stage != store.DeliveryStageParked || item.Blocker != wantBlocker {
		t.Fatalf("findings-after-archive item = %+v, want parked as %q", item, wantBlocker)
	}
}

func TestDeliveryEngineParksFindingsWithoutArchiveAsReviewFindings(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const (
		gitRoot  = "/repo-review-findings"
		specSlug = "delivery-spec"
	)
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{specSlug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.reviewResults[specSlug] = ReviewResult{
		Outcome: ReviewOutcomeFindings,
		Head:    "reviewed-" + specSlug,
		Reason:  "finding",
	}
	engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Stage != store.DeliveryStageParked || item.Blocker != BlockerReviewFindings {
		t.Fatalf("ordinary findings item = %+v, want parked as %q", item, BlockerReviewFindings)
	}
}

func TestRetryRefusesACorrectiveSpecItemWhoseHeadMoved(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-corrective-moved"
	blocker := BlockerCorrectiveSpecRequired + ": 0176-one, 0177-two"
	before := seedParkedRetryItem(t, ctx, runStore, gitRoot, "delivery-spec", blocker)
	recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "moved-head"}}}
	engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())

	_, err := engine.Retry(ctx, gitRoot, before.SpecSlug)

	if err == nil {
		t.Fatal("moved corrective-Spec retry succeeded, want refusal")
	}
	for _, want := range []string{"0176-one, 0177-two", "moved-head", "candidate-original", "author a corrective Spec with its own authorization and QA gate"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("moved corrective-Spec error = %q, want %q", err, want)
		}
	}
	if !reflect.DeepEqual(recovery.events, []string{"inspect"}) {
		t.Fatalf("moved corrective-Spec recovery calls = %v, want inspect only", recovery.events)
	}
	assertRetryItemUnchanged(t, ctx, runStore, gitRoot, before)
}

func TestRetryReturnsAnUnchangedCorrectiveSpecItemToReviewing(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-corrective-unchanged"
	blocker := BlockerCorrectiveSpecRequired + ": 0176-one"
	before := seedParkedRetryItem(t, ctx, runStore, gitRoot, "delivery-spec", blocker)
	recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "candidate-original"}}}
	engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())

	result, err := engine.Retry(ctx, gitRoot, before.SpecSlug)

	if err != nil {
		t.Fatalf("retry unchanged corrective-Spec item: %v", err)
	}
	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if result.Stage != store.DeliveryStageReviewing || item.Stage != store.DeliveryStageReviewing || item.Blocker != "" {
		t.Fatalf("unchanged corrective-Spec retry = result:%+v item:%+v, want reviewing with no blocker", result, item)
	}
	if !reflect.DeepEqual(item.CandidateCommits, before.CandidateCommits) {
		t.Fatalf("unchanged corrective-Spec candidate commits = %v, want %v", item.CandidateCommits, before.CandidateCommits)
	}
	if !reflect.DeepEqual(recovery.events, []string{"inspect"}) {
		t.Fatalf("unchanged corrective-Spec recovery calls = %v, want inspect only", recovery.events)
	}
}
