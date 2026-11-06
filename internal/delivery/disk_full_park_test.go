// Boundary: real queue store and engine, with fake Run, review and gate failures.
package delivery

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"roundfix/internal/store"
)

func TestAFullDiskParksTheItemAsDiskFull(t *testing.T) {
	t.Parallel()
	cause := errors.New("write /tmp/x/_pkg_.a: no space left on device")
	full := &DiskFullError{Path: "/tmp/x/_pkg_.a", AvailableBytes: 3*1024*1024 + 123, Err: cause}
	if full.Error() != cause.Error() || !errors.Is(full, cause) {
		t.Fatalf("error chain=%v", full)
	}
	wrapped := fmt.Errorf("build item binary: %w", full)
	const wantBlocker = "disk-full: /tmp/x/_pkg_.a: 3 MiB free"
	for _, tc := range []struct {
		name    string
		stepErr error
		gate    bool
		want    string
		class   ParkClass
	}{
		{"Run", wrapped, false, wantBlocker, ParkClassEnvironment},
		{"review", wrapped, false, wantBlocker, ParkClassEnvironment},
		{"gate", nil, true, wantBlocker, ParkClassEnvironment},
		{"other Run error", errors.New("executor unavailable"), false, "delivery-error: run Implement executor: executor unavailable", ParkClassEnvironment},
		{"other review error", errors.New("review unavailable"), false, "delivery-error: run pre-PR review: review unavailable", ParkClassEnvironment},
		{"other gate failure", nil, true, BlockerGateFailed, ParkClassFinding},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			var queue store.DeliveryQueue
			if tc.name == "review" || tc.name == "other review error" {
				fixture := newReviewSlotFixture(t, "widget")
				fixture.seed(t, 0, store.DeliveryStageReviewing)
				fixture.workflow.review = func(context.Context, string) (ReviewOutcome, error) { return "", tc.stepErr }
				if _, err := fixture.engine.Run(ctx, fixture.root); err != nil {
					t.Fatal(err)
				}
				queue = readDeliveryQueue(t, ctx, fixture.store, fixture.root)
			} else {
				s := openDeliveryEngineStore(t, ctx)
				const root = "/repo-full-disk"
				if _, err := s.CreateDeliveryQueue(ctx, root, []string{"widget"}); err != nil {
					t.Fatal(err)
				}
				workflow := newFakeDeliveryWorkflow()
				if tc.gate {
					blocker := ""
					if tc.name == "gate" {
						blocker = full.Blocker()
					}
					workflow.gates["widget"] = GateResult{Reason: cause.Error(), Blocker: blocker}
				} else {
					workflow.runErrors["widget"] = tc.stepErr
				}
				if _, err := newTestDeliveryEngine(s, workflow, newFakeDeliveryBoundary()).Run(ctx, root); err != nil {
					t.Fatal(err)
				}
				queue = readDeliveryQueue(t, ctx, s, root)
			}
			item := queue.Items[0]
			park := ClassifyPark(queue, item)
			if item.Stage != store.DeliveryStageParked || item.Blocker != tc.want || park.Class != tc.class {
				t.Fatalf("item=%+v park=%+v", item, park)
			}
			if tc.want == wantBlocker && park.Next != "free space on the filesystem that holds the path the blocker names: run roundfix gc to trim declared Shared Build Caches and reclaim ended Runs, and remove other files there; then run roundfix deliver retry widget; a retry from disk-full is not counted" {
				t.Fatalf("action=%q", park.Next)
			}
		})
	}
}

func TestADiskFullRetryIsNotCounted(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name              string
		atLimit, archived bool
	}{
		{"below limit", false, false}, {"at limit", true, false}, {"archived gate at limit", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			s := openDeliveryEngineStore(t, ctx)
			const root = "/repo-disk-full-retry"
			const slug = "widget"
			queue, err := s.CreateDeliveryQueueWithLimits(ctx, root, []string{slug}, store.DeliveryQueueLimits{MaxRetries: 1})
			if err != nil {
				t.Fatal(err)
			}
			item := queue.Items[0]
			item.Stage = store.DeliveryStageParked
			item.Blocker = BlockerDeliveryError
			item.Branch = "roundfix/deliver-" + slug
			item.Worktree = "/worktrees/" + slug
			item.RunID = "original-run"
			item.CandidateCommits = []string{"candidate-original"}
			recordDeliveryItemWorkspace(t, ctx, s, root, &item)
			if err := s.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
				t.Fatal(err)
			}
			wantCount := 0
			if tc.atLimit {
				target := item
				target.Stage = store.DeliveryStageRunning
				if _, _, err := s.RetryDeliveryQueueItem(ctx, root, target, item.Blocker); err != nil {
					t.Fatal(err)
				}
				item = readDeliveryQueue(t, ctx, s, root).Items[0]
				wantCount = 1
			}
			item.Stage = store.DeliveryStageParked
			item.Blocker = BlockerDiskFull + ": /tmp/x/_pkg_.a: 0 MiB free"
			if err := s.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
				t.Fatal(err)
			}
			state := ItemState{Head: "candidate-original", Archived: tc.archived}
			recovery := &fakeItemRecovery{states: []ItemState{state, state}, carryResult: CarryForwardResult{RunID: "carried-run"}}
			engine := newRetryDeliveryEngine(s, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())
			if _, err := engine.Retry(ctx, root, slug); err != nil {
				t.Fatal(err)
			}
			got := readDeliveryQueue(t, ctx, s, root).Items[0]
			wantStage := store.DeliveryStageReviewing
			if tc.archived {
				wantStage = store.DeliveryStageGating
			}
			if got.RetryCount != wantCount || got.Blocker != "" || got.Stage != wantStage {
				t.Fatalf("item=%+v want stage=%s count=%d", got, wantStage, wantCount)
			}
		})
	}
}
