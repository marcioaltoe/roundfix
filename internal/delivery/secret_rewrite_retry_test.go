// Boundary: real queue persistence and Retry; Git proof and inspection are fakes.
package delivery

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"roundfix/internal/store"
)

type fakeSecretRewriteProver struct {
	proof                          SecretRewrite
	err                            error
	workDir, base, candidate, head string
	named                          []string
}

func (p *fakeSecretRewriteProver) ProveSecretRewrite(_ context.Context, workDir, base, candidate, head string, named []string) (SecretRewrite, error) {
	p.workDir, p.base, p.candidate, p.head, p.named = workDir, base, candidate, head, named
	return p.proof, p.err
}

func TestRetryAcceptsASecretRewriteOfAnArchivedItem(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	db := openDeliveryEngineStore(t, ctx)
	const root = "/repo-secret-rewrite"
	before := seedParkedRetryItem(t, ctx, db, root, "rewrite", BlockerSecretFound+": a:file.txt:2: credential-assignment (commit abc), other file.txt:7: github-token, a:file.txt:9: credential-assignment (commit def)")
	before.CandidateCommits = append(before.CandidateCommits, "newest-candidate")
	if err := db.UpdateDeliveryQueueItem(ctx, root, before); err != nil {
		t.Fatal(err)
	}
	recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "rewritten-head"}}}
	history := &operatorArchiveHistory{}
	prover := &fakeSecretRewriteProver{proof: SecretRewrite{Accepted: true}}
	engine := NewEngine(db, EngineDependencies{
		Workspace: newFakeDeliveryWorkflow(), Recovery: recovery,
		Revalidator: cleanTestRevalidator{}, History: history, SecretRewrites: prover,
	})
	result, err := engine.Retry(ctx, root, before.SpecSlug)
	if err != nil {
		t.Fatal(err)
	}
	got := readDeliveryQueue(t, ctx, db, root).Items[0]
	if result.Stage != store.DeliveryStageReviewing || got.Stage != store.DeliveryStageReviewing || got.Blocker != "" || !reflect.DeepEqual(got.CandidateCommits, []string{"candidate-original", "newest-candidate", "rewritten-head"}) {
		t.Fatalf("result=%+v item=%+v", result, got)
	}
	if prover.workDir != before.Worktree || prover.base != "run-start" || prover.candidate != "newest-candidate" || prover.head != "rewritten-head" || !reflect.DeepEqual(prover.named, []string{"a:file.txt", "other file.txt"}) || history.runID != before.RunID {
		t.Fatalf("proof request=%+v history=%+v", prover, history)
	}
}

func TestRetryRefusesARejectedSecretRewriteWithItsRule(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		proof   SecretRewrite
		err     error
		history ItemHistory
		reason  string
	}{
		{"rule", SecretRewrite{Reason: "the rewrite changes the line count of named.txt"}, nil, &operatorArchiveHistory{}, "the rewrite changes the line count of named.txt"},
		{"prover error", SecretRewrite{Accepted: true}, errors.New("cannot read objects"), &operatorArchiveHistory{}, "cannot read objects"},
		{"missing history", SecretRewrite{Accepted: true}, nil, nil, "the item has no Run start to bound the rewrite"},
		{"unreadable Run start", SecretRewrite{Accepted: true}, nil, &operatorArchiveHistory{err: errors.New("missing Run")}, "the item has no Run start to bound the rewrite"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := openDeliveryEngineStore(t, ctx)
			const root = "/repo-secret-rewrite-refused"
			before := seedParkedRetryItem(t, ctx, db, root, "rewrite", BlockerSecretFound+": named.txt:2: credential-assignment")
			recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "rewritten-head"}}}
			engine := newRetryDeliveryEngine(db, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())
			engine.history = tc.history
			engine.secretRewrites = &fakeSecretRewriteProver{proof: tc.proof, err: tc.err}
			_, err := engine.Retry(ctx, root, before.SpecSlug)
			want := fmt.Sprintf("retry Delivery Queue item %q: archived item head %q differs from candidate head %q (%s)", before.SpecSlug, "rewritten-head", "candidate-original", tc.reason)
			if err == nil || err.Error() != want {
				t.Fatalf("refusal=%v want=%s", err, want)
			}
			assertRetryItemUnchanged(t, ctx, db, root, before)
		})
	}
}

func TestSecretRewriteRetryKeepsExistingPaths(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, blocker, head string
		descends, noProver  bool
		stage               store.DeliveryStage
	}{
		{name: "other blocker", blocker: BlockerGateFailed, head: "rewritten-head"},
		{name: "lookalike blocker", blocker: BlockerSecretFound + "-other", head: "rewritten-head"},
		{name: "no prover", blocker: BlockerSecretFound, head: "rewritten-head", noProver: true},
		{name: "unchanged candidate", blocker: BlockerSecretFound, head: "candidate-original", stage: store.DeliveryStageGating},
		{name: "descending head", blocker: BlockerSecretFound, head: "rewritten-head", descends: true, stage: store.DeliveryStageReviewing},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := openDeliveryEngineStore(t, ctx)
			const root = "/repo-secret-existing-path"
			before := seedParkedRetryItem(t, ctx, db, root, "existing", tc.blocker)
			recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: tc.head}}}
			engine := newRetryDeliveryEngine(db, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())
			history := &operatorArchiveHistory{descends: tc.descends}
			engine.history = history
			prover := &fakeSecretRewriteProver{proof: SecretRewrite{Accepted: true}}
			if !tc.noProver {
				engine.secretRewrites = prover
			}
			result, err := engine.Retry(ctx, root, before.SpecSlug)
			if tc.stage == "" {
				want := fmt.Sprintf("retry Delivery Queue item %q: archived item head %q differs from candidate head %q", before.SpecSlug, tc.head, "candidate-original")
				if err == nil || err.Error() != want {
					t.Fatalf("refusal=%v want=%s", err, want)
				}
				assertRetryItemUnchanged(t, ctx, db, root, before)
			} else if err != nil || result.Stage != tc.stage {
				t.Fatalf("result=%+v err=%v", result, err)
			}
			if prover.workDir != "" || history.runID != "" {
				t.Fatalf("unexpected rewrite proof=%+v history=%+v", prover, history)
			}
		})
	}
}
