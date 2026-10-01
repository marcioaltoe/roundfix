// Suite: Delivery Queue token ceiling.
// Invariant: recorded queue tokens gate queued starts and retries, never an item already running.
// Boundary IN: Delivery Engine and real SQLite usage and queue records.
// Boundary OUT: workspace, Run and publication actions use the existing engine fakes.
package delivery

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"roundfix/internal/store"
)

func seedTokenCeilingQueue(t *testing.T, ceiling, tokens int64) (*store.Store, string, store.DeliveryQueueItem) {
	t.Helper()
	ctx := t.Context()
	writer := openDeliveryEngineStore(t, ctx)
	root := "/repo-token-ceiling"
	queue, err := writer.CreateDeliveryQueueWithLimits(ctx, root, []string{"spent", "next"}, store.DeliveryQueueLimits{MaxTokens: ceiling})
	if err != nil {
		t.Fatal(err)
	}
	run, err := writer.CreateRun(ctx, store.CreateRunRequest{Kind: store.KindImplement, GitRoot: root, LocalBranch: "feat/spent", HeadSHA: "spent-head", SpecSlug: "spent", Agent: "codex"})
	if err != nil {
		t.Fatal(err)
	}
	if err := writer.AppendTokenUsage(ctx, store.TokenUsageRecord{RunID: run.ID, ScopeKind: "task", ScopeID: "task_01", Session: "spent", Basis: "request-sum", TotalTokens: &tokens}); err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage, item.RunID = store.DeliveryStageMerged, run.ID
	if err := writer.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
		t.Fatal(err)
	}
	return writer, root, queue.Items[1]
}

func TestAQueuedItemParksAtTheTokenCeilingWithoutAWorktree(t *testing.T) {
	t.Parallel()
	for _, tokens := range []int64{5000000, 5639755} {
		t.Run(fmt.Sprint(tokens), func(t *testing.T) {
			writer, root, item := seedTokenCeilingQueue(t, 5000000, tokens)
			workflow := newFakeDeliveryWorkflow()
			engine := newTestDeliveryEngine(writer, workflow, newFakeDeliveryBoundary())
			if _, err := engine.Run(t.Context(), root); err != nil {
				t.Fatal(err)
			}
			got := readDeliveryQueue(t, t.Context(), writer, root).Items[1]
			if got.Stage != store.DeliveryStageParked || got.Blocker != BlockerQueueTokenCeiling {
				t.Fatalf("item = %+v, want token-ceiling park", got)
			}
			if got.Branch != "" || got.Worktree != "" || got.WorktreeProvisioned || len(workflow.events[item.SpecSlug]) != 0 {
				t.Fatalf("park created workspace or actions: item=%+v events=%v", got, workflow.events[item.SpecSlug])
			}
		})
	}
}

func TestAQueuedItemStartsBelowTheTokenCeiling(t *testing.T) {
	t.Parallel()
	assertTokenCeilingItemStarts(t, 5000000, 4999999)
}

func TestNoTokenCeilingNeverParks(t *testing.T) {
	t.Parallel()
	assertTokenCeilingItemStarts(t, 0, 5639755)
}

func assertTokenCeilingItemStarts(t *testing.T, ceiling, tokens int64) {
	t.Helper()
	writer, root, item := seedTokenCeilingQueue(t, ceiling, tokens)
	workflow := newFakeDeliveryWorkflow()
	engine := newTestDeliveryEngine(writer, workflow, newFakeDeliveryBoundary())
	if _, err := engine.Run(t.Context(), root); err != nil {
		t.Fatal(err)
	}
	got := readDeliveryQueue(t, t.Context(), writer, root).Items[1]
	events := workflow.events[item.SpecSlug]
	if got.Stage != store.DeliveryStageMerged || len(events) == 0 || events[0] != "create-branch" {
		t.Fatalf("item=%+v events=%v, want started and merged", got, events)
	}
}

func TestARunningItemContinuesAboveTheTokenCeiling(t *testing.T) {
	t.Parallel()
	writer, root, item := seedTokenCeilingQueue(t, 5000000, 5639755)
	item.Stage = store.DeliveryStageRunning
	item.Branch, item.Worktree = "roundfix/deliver-next", "/worktrees/next"
	recordDeliveryItemWorkspace(t, t.Context(), writer, root, &item)
	if err := writer.UpdateDeliveryQueueItem(t.Context(), root, item); err != nil {
		t.Fatal(err)
	}
	workflow := newFakeDeliveryWorkflow()
	engine := newTestDeliveryEngine(writer, workflow, newFakeDeliveryBoundary())
	if _, err := engine.Run(t.Context(), root); err != nil {
		t.Fatal(err)
	}
	got := readDeliveryQueue(t, t.Context(), writer, root).Items[1]
	events := workflow.events[item.SpecSlug]
	if got.Stage != store.DeliveryStageMerged || len(events) == 0 || events[0] != "run" {
		t.Fatalf("running item=%+v events=%v, want Run continuing through merge", got, events)
	}
}

func TestARetryIsRefusedAtTheTokenCeiling(t *testing.T) {
	t.Parallel()
	for _, blocker := range []string{BlockerQueueTokenCeiling, BlockerQueueDeadline, BlockerRunUnresolved, BlockerPrerequisiteUnmerged + ": spent"} {
		t.Run(blocker, func(t *testing.T) {
			writer, root, item := seedTokenCeilingQueue(t, 5000000, 5000000)
			item.Stage, item.Blocker = store.DeliveryStageParked, blocker
			if err := writer.UpdateDeliveryQueueItem(t.Context(), root, item); err != nil {
				t.Fatal(err)
			}
			workflow := newFakeDeliveryWorkflow()
			workspace := &countingRetryWorkspace{fakeDeliveryWorkflow: workflow}
			recovery := &fakeItemRecovery{}
			engine := newRetryDeliveryEngine(writer, workflow, workspace, recovery, newFakeDeliveryBoundary())
			_, err := engine.Retry(context.Background(), root, item.SpecSlug)
			want := `retry Delivery Queue item "next": queue token ceiling 5000000 was reached with 5000000 tokens; start a new queue with roundfix deliver start`
			if err == nil || err.Error() != want {
				t.Fatalf("retry error=%v, want %s", err, want)
			}
			got := readDeliveryQueue(t, t.Context(), writer, root).Items[1]
			if !reflect.DeepEqual(got, item) || workspace.useCalls != 0 || len(recovery.events) != 0 || len(workflow.events[item.SpecSlug]) != 0 {
				t.Fatalf("refused retry changed state or performed actions: item=%+v workspace=%d recovery=%v", got, workspace.useCalls, recovery.events)
			}
		})
	}
}

func TestPendingQuestionAnswersTheTokenCeiling(t *testing.T) {
	t.Parallel()
	queue := store.DeliveryQueue{Items: []store.DeliveryQueueItem{{SpecSlug: "next", Stage: store.DeliveryStageParked, Blocker: BlockerQueueTokenCeiling}}}
	got, found := PendingQuestionFor(queue)
	want := PendingQuestion{SpecSlug: "next", Blocker: BlockerQueueTokenCeiling, Answer: "record a new queue for the remaining Specs with roundfix deliver start and a higher --max-tokens"}
	if !found || got != want {
		t.Fatalf("question=%+v found=%v, want %+v", got, found, want)
	}
}

func TestTokenCeilingHasABudgetParkClass(t *testing.T) {
	t.Parallel()
	for _, suffix := range []string{"", ": details"} {
		t.Run(suffix, func(t *testing.T) {
			item := store.DeliveryQueueItem{SpecSlug: "next", Stage: store.DeliveryStageParked, Blocker: BlockerQueueTokenCeiling + suffix}
			queue := store.DeliveryQueue{Items: []store.DeliveryQueueItem{item}}
			park := ClassifyPark(queue, item)
			want := ParkClassification{Class: ParkClassBudget, Next: "record a new queue for the remaining Specs with roundfix deliver start and a higher --max-tokens"}
			question, found := PendingQuestionFor(queue)
			if park != want || !found || question.Answer != want.Next {
				t.Fatalf("classification=%+v question=%+v found=%v, want %+v", park, question, found, want)
			}
		})
	}
}
