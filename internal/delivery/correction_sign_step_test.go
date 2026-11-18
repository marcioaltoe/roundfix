// Boundary: real queue persistence and merge/retry transitions; the sign gate
// and forge are injected I/O boundaries.
package delivery

import (
	"context"
	"errors"
	"strings"
	"testing"

	"roundfix/internal/store"
)

type correctionSignGateFunc func(context.Context, string, string, string, string) (CorrectionSignCheck, error)

func (f correctionSignGateFunc) CheckCorrectionSigns(ctx context.Context, root, dir, slug, head string) (CorrectionSignCheck, error) {
	return f(ctx, root, dir, slug, head)
}

func TestTheMergeStepParksAnUndeclaredCorrectionSign(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name                                string
		required, undeclared, fail, nilGate bool
		park, stepError                     bool
	}{
		{name: "on undeclared", required: true, undeclared: true, park: true},
		{name: "on covered", required: true}, {name: "on not-a-correction", required: true},
		{name: "off with sign", undeclared: true}, {name: "nil dependency", nilGate: true},
		{name: "on error", required: true, fail: true, stepError: true}, {name: "off error", fail: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			db := openDeliveryEngineStore(t, ctx)
			root := t.TempDir()
			item := seedParkedRetryItem(t, ctx, db, root, "example", BlockerDefectClassUndeclared)
			item.Stage = store.DeliveryStageMerging
			item.Blocker = ""
			item.PullRequestNumber = "1"
			if err := db.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
				t.Fatal(err)
			}
			workflow := newFakeDeliveryWorkflow()
			boundary := newFakeDeliveryBoundary()
			engine := newTestDeliveryEngine(db, workflow, boundary)
			if !tc.nilGate {
				engine.correctionSigns = correctionSignGateFunc(func(_ context.Context, gotRoot, dir, slug, head string) (CorrectionSignCheck, error) {
					if gotRoot != root || dir != item.Worktree || slug != item.SpecSlug || head != "candidate-original" {
						t.Fatalf("gate arguments: %s %s %s %s", gotRoot, dir, slug, head)
					}
					check := CorrectionSignCheck{Required: tc.required}
					if tc.undeclared {
						check.Undeclared = []string{`commit abc1234 "repair"`}
					}
					if tc.fail {
						return check, errors.New("sign source unavailable")
					}
					return check, nil
				})
			}
			err := engine.mergeCandidate(ctx, root, &item)
			if (err != nil) != tc.stepError {
				t.Fatalf("merge error = %v", err)
			}
			got := readDeliveryQueue(t, ctx, db, root).Items[0]
			if tc.park {
				if got.Stage != store.DeliveryStageParked || got.Blocker != BlockerDefectClassUndeclared+`: commit abc1234 "repair"` {
					t.Fatalf("item=%+v", got)
				}
				park := ClassifyPark(store.DeliveryQueue{}, got)
				for _, part := range []string{"plan", "implementation", "environment", "operator", "unclear", "not-a-correction", `commit abc1234 "repair"`, "roundfix corrections declare --spec example", "roundfix deliver retry example"} {
					if !strings.Contains(park.Next, part) {
						t.Fatalf("next action %q misses %q", park.Next, part)
					}
				}
			} else if !tc.stepError && got.Stage != store.DeliveryStageMerged {
				t.Fatalf("item=%+v", got)
			}
			wantMerges := 1
			if tc.park || tc.stepError {
				wantMerges = 0
			}
			if boundary.mergeEffects != wantMerges {
				t.Fatalf("merges=%d want %d", boundary.mergeEffects, wantMerges)
			}
			if tc.park || tc.stepError {
				assertNoUnmatchedDeliveryIntents(t, ctx, db, root)
			}
		})
	}
}

func TestADefectClassParkRetriesToTheStageItsHeadSupports(t *testing.T) {
	t.Parallel()
	for _, moved := range []bool{false, true} {
		name := "unchanged"
		if moved {
			name = "moved"
		}
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			db := openDeliveryEngineStore(t, ctx)
			root := t.TempDir()
			item := seedParkedRetryItem(t, ctx, db, root, "example", BlockerDefectClassUndeclared+": commit abc1234")
			item.PullRequestNumber = "1"
			if err := db.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
				t.Fatal(err)
			}
			head := "candidate-original"
			want := store.DeliveryStageChecking
			if moved {
				head = "changed-head"
				want = store.DeliveryStageReviewing
			}
			recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: head}}}
			workflow := newFakeDeliveryWorkflow()
			boundary := newFakeDeliveryBoundary()
			boundary.pullRequests[item.Branch] = PullRequest{Number: "1", State: "OPEN", HeadBranch: item.Branch, HeadSHA: head}
			engine := newRetryDeliveryEngine(db, workflow, nil, recovery, boundary)
			engine.history = &operatorArchiveHistory{descends: true}
			// The operator has declared the cause before the retry reaches the gate.
			engine.correctionSigns = correctionSignGateFunc(func(context.Context, string, string, string, string) (CorrectionSignCheck, error) {
				return CorrectionSignCheck{Required: true}, nil
			})
			result, err := engine.Retry(ctx, root, item.SpecSlug)
			if err != nil {
				t.Fatal(err)
			}
			got := readDeliveryQueue(t, ctx, db, root).Items[0]
			if result.Stage != want || got.Stage != want || got.Blocker != "" || got.RetryCount != 1 {
				t.Fatalf("result=%+v item=%+v", result, got)
			}
			if !moved {
				if _, err := engine.Run(ctx, root); err != nil {
					t.Fatal(err)
				}
				got = readDeliveryQueue(t, ctx, db, root).Items[0]
				if got.Stage != store.DeliveryStageMerged || boundary.mergeEffects != 1 {
					t.Fatalf("item=%+v merges=%d", got, boundary.mergeEffects)
				}
			}
		})
	}
}
