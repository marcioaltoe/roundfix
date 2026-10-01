// Suite: operator-archived recovery after an environment-only QA partial.
// Boundary IN: real queue persistence and engine stages.
// Boundary OUT: Git ancestry, item inspection, and delivery executors use fakes.
package delivery

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/store"
)

type operatorArchiveHistory struct {
	descends        bool
	err             error
	ancestor, runID string
}

func (history *operatorArchiveHistory) Descends(_ context.Context, _, ancestor, _ string) (bool, error) {
	history.ancestor = ancestor
	return history.descends, history.err
}
func (history *operatorArchiveHistory) RunStart(_ context.Context, _, runID string) (string, error) {
	history.runID = runID
	return "run-start", history.err
}

func TestAnEnvironmentOnlyPartialParksAsQAEnvironmentPartial(t *testing.T) {
	for _, partial := range []bool{true, false} {
		t.Run(map[bool]string{true: "environment", false: "unresolved"}[partial], func(t *testing.T) {
			ctx := t.Context()
			runStore := openDeliveryEngineStore(t, ctx)
			const root = "/repo-qa-partial"
			if _, err := runStore.CreateDeliveryQueue(ctx, root, []string{"partial"}); err != nil {
				t.Fatal(err)
			}
			workflow := newFakeDeliveryWorkflow()
			workflow.runs["partial"] = RunResult{RunID: "run-partial", Outcome: RunOutcomeUnresolved, QAEnvironmentPartial: partial}
			engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())
			if _, err := engine.Run(ctx, root); err != nil {
				t.Fatal(err)
			}
			item := readDeliveryQueue(t, ctx, runStore, root).Items[0]
			want := BlockerRunUnresolved
			if partial {
				want = BlockerQAEnvironmentPartial
			}
			if item.Stage != store.DeliveryStageParked || item.Blocker != want || item.RunID != "run-partial" {
				t.Fatalf("item = %+v", item)
			}
			if partial {
				classification := ClassifyPark(store.DeliveryQueue{}, item)
				if classification.Class != ParkClassEnvironment || !strings.Contains(classification.Next, "roundfix reconcile run-partial --carry-forward") || !strings.Contains(classification.Next, "--qa-override") {
					t.Fatalf("classification = %+v", classification)
				}
			}
		})
	}
}

func TestRetryResumesAnOperatorArchivedItemAtReview(t *testing.T) {
	for _, noCandidate := range []bool{false, true} {
		t.Run(map[bool]string{false: "candidate anchor", true: "Run start anchor"}[noCandidate], func(t *testing.T) {
			ctx := t.Context()
			runStore := openDeliveryEngineStore(t, ctx)
			const root = "/repo-operator-retry"
			item := seedParkedRetryItem(t, ctx, runStore, root, "operator-retry", BlockerQAEnvironmentPartial)
			if noCandidate {
				item.CandidateCommits = nil
				if err := runStore.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
					t.Fatal(err)
				}
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
			want := append(item.CandidateCommits, "operator-head")
			if result.Stage != store.DeliveryStageReviewing || !reflect.DeepEqual(got.CandidateCommits, want) || len(recovery.events) != 1 {
				t.Fatalf("result=%+v item=%+v events=%v", result, got, recovery.events)
			}
			anchor := "candidate-original"
			if noCandidate {
				anchor = "run-start"
				if history.runID != item.RunID {
					t.Fatalf("Run ID = %q", history.runID)
				}
			}
			if history.ancestor != anchor {
				t.Fatalf("anchor=%q want %q", history.ancestor, anchor)
			}
			workflow := newFakeDeliveryWorkflow()
			workflow.archives[item.SpecSlug] = ArchiveResult{AlreadyArchived: true, Head: "operator-head"}
			engine.reviewer = workflow
			engine.archiver = workflow
			if err := engine.reviewCandidate(ctx, root, &got); err != nil {
				t.Fatal(err)
			}
			if got.Stage != store.DeliveryStageArchiving {
				t.Fatalf("review stage = %s", got.Stage)
			}
			if err := engine.archiveCandidate(ctx, root, &got); err != nil {
				t.Fatal(err)
			}
			if got.Stage != store.DeliveryStageGating || !reflect.DeepEqual(got.CandidateCommits, want) {
				t.Fatalf("archived retry = %+v", got)
			}

		})
	}
}

func TestRetryRefusesAnOperatorArchiveWithoutAQAOverride(t *testing.T) {
	testOperatorArchiveRefusal(t, false, &operatorArchiveHistory{descends: true})
}
func TestRetryRefusesAnArchivedHeadThatDoesNotDescendFromTheAnchor(t *testing.T) {
	testOperatorArchiveRefusal(t, true, &operatorArchiveHistory{})
}
func TestOperatorArchiveRetryRequiresHistory(t *testing.T) { testOperatorArchiveRefusal(t, true, nil) }
func TestOperatorArchiveRetryRefusesUnreadableHistory(t *testing.T) {
	testOperatorArchiveRefusal(t, true, &operatorArchiveHistory{err: errors.New("history unavailable")})
}
func testOperatorArchiveRefusal(t *testing.T, override bool, history *operatorArchiveHistory) {
	t.Helper()
	ctx := t.Context()
	runStore := openDeliveryEngineStore(t, ctx)
	const root = "/repo-operator-refusal"
	before := seedParkedRetryItem(t, ctx, runStore, root, "operator-refusal", BlockerQAEnvironmentPartial)
	recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, QAOverride: override, Head: "operator-head"}}}
	engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())
	if history != nil {
		engine.history = history
	}
	_, err := engine.Retry(ctx, root, before.SpecSlug)
	if err == nil || !strings.Contains(err.Error(), `archived item head "operator-head" differs from candidate head "candidate-original"`) {
		t.Fatalf("refusal=%v", err)
	}
	assertRetryItemUnchanged(t, ctx, runStore, root, before)
}

func TestTheArchiveStagePassesAnAlreadyArchivedSpec(t *testing.T) {
	ctx := t.Context()
	runStore := openDeliveryEngineStore(t, ctx)
	const root = "/repo-already-archived"
	item := seedParkedRetryItem(t, ctx, runStore, root, "already-archived", BlockerQAEnvironmentPartial)
	workflow := newFakeDeliveryWorkflow()
	workflow.archives[item.SpecSlug] = ArchiveResult{AlreadyArchived: true, Head: "candidate-original"}
	engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())
	item.Stage = store.DeliveryStageArchiving
	if err := engine.archiveCandidate(ctx, root, &item); err != nil {
		t.Fatal(err)
	}
	got := readDeliveryQueue(t, ctx, runStore, root).Items[0]
	if got.Stage != store.DeliveryStageGating || !reflect.DeepEqual(got.CandidateCommits, []string{"candidate-original"}) {
		t.Fatalf("item=%+v", got)
	}
}
