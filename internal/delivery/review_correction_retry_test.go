// Boundary: real queue persistence; correction proof and worktree inspection are fakes.
package delivery

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/store"
)

type fakeReviewCorrectionProver struct {
	proof                    ReviewCorrection
	err                      error
	workDir, candidate, head string
	specs                    []string
}

func (p *fakeReviewCorrectionProver) ProveReviewCorrection(_ context.Context, workDir string, specs []string, candidate, head string) (ReviewCorrection, error) {
	p.workDir, p.specs, p.candidate, p.head = workDir, specs, candidate, head
	return p.proof, p.err
}
func TestRetryReturnsAReviewOnlyCorrectionToReview(t *testing.T) {
	ctx := t.Context()
	db := openDeliveryEngineStore(t, ctx)
	const root = "/repo-review-correction"
	before := seedParkedRetryItem(t, ctx, db, root, "correction", BlockerCorrectiveSpecRequired+": one, two")
	recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "correction-head"}}}
	engine := newRetryDeliveryEngine(db, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())
	prover := &fakeReviewCorrectionProver{proof: ReviewCorrection{Accepted: true}}
	engine.corrections = prover
	result, err := engine.Retry(ctx, root, before.SpecSlug)
	if err != nil {
		t.Fatal(err)
	}
	got := readDeliveryQueue(t, ctx, db, root).Items[0]
	if result.Stage != store.DeliveryStageReviewing || got.Stage != store.DeliveryStageReviewing || got.Blocker != "" || !reflect.DeepEqual(got.CandidateCommits, []string{"candidate-original", "correction-head"}) {
		t.Fatalf("result=%+v item=%+v", result, got)
	}
	if prover.workDir == "" || prover.candidate != "candidate-original" || prover.head != "correction-head" || !reflect.DeepEqual(prover.specs, []string{"one", "two"}) {
		t.Fatalf("proof request=%+v", prover)
	}
}
func TestRetryRefusesACorrectionTheProofRejects(t *testing.T) {
	for _, tc := range []struct {
		name  string
		proof ReviewCorrection
		err   error
	}{
		{name: "refused", proof: ReviewCorrection{Reason: "outside archive"}},
		{name: "error", proof: ReviewCorrection{Accepted: true}, err: errors.New("cannot read proof")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			db := openDeliveryEngineStore(t, ctx)
			const root = "/repo-review-correction-refused"
			before := seedParkedRetryItem(t, ctx, db, root, "correction", BlockerCorrectiveSpecRequired+": one")
			recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "correction-head"}}}
			engine := newRetryDeliveryEngine(db, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())
			engine.corrections = &fakeReviewCorrectionProver{proof: tc.proof, err: tc.err}
			var log bytes.Buffer
			engine.log = &log
			_, err := engine.Retry(ctx, root, before.SpecSlug)
			reason := tc.proof.Reason
			if tc.err != nil {
				reason = tc.err.Error()
			}
			want := `retry Delivery Queue item "correction": archived Specs one were reviewed at parked candidate head "candidate-original", but the item head is "correction-head" (` + reason + `); author a corrective Spec with its own authorization and QA gate`
			if err == nil || err.Error() != want {
				t.Fatalf("refusal=%v want=%s", err, want)
			}
			if tc.err != nil && !strings.Contains(log.String(), reason) {
				t.Fatalf("log=%q", log.String())
			}
			assertRetryItemUnchanged(t, ctx, db, root, before)
		})
	}
}
