// Suite: runtime infrastructure parks and uncounted retries.
// Boundary IN: Delivery Engine and the real queue store; execution uses fakes.
package delivery

import (
	"roundfix/internal/store"
	"strings"
	"testing"
)

func TestRuntimeInfrastructureParksTheItem(t *testing.T) {
	ctx := t.Context()
	s := openDeliveryEngineStore(t, ctx)
	const root = "/repo-runtime-infrastructure"
	const slug = "runtime-spec"
	if _, err := s.CreateDeliveryQueue(ctx, root, []string{slug}); err != nil {
		t.Fatal(err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.runs[slug] = RunResult{RunID: "lost-run", Outcome: RunOutcomeUnresolved, RuntimeInfrastructure: "task_03 lost its rollout at session/resume", QAEnvironmentPartial: true}
	workflow.recordWorkspace = func(branch, worktree string) error {
		_, _, _, err := s.RecordDeliveryQueueItemWorktree(ctx, root, slug, branch, worktree)
		return err
	}
	if _, err := newTestDeliveryEngine(s, workflow, newFakeDeliveryBoundary()).Run(ctx, root); err != nil {
		t.Fatal(err)
	}
	item := readDeliveryQueue(t, ctx, s, root).Items[0]
	if item.Stage != store.DeliveryStageParked || item.Blocker != BlockerRuntimeInfrastructure+": task_03 lost its rollout at session/resume" || item.RunID != "lost-run" {
		t.Fatalf("item=%+v", item)
	}
}
func TestClassifyParkRuntimeInfrastructure(t *testing.T) {
	item := store.DeliveryQueueItem{SpecSlug: "runtime-spec", Blocker: BlockerRuntimeInfrastructure + ": qa lost its rollout at session/prompt"}
	got := ClassifyPark(store.DeliveryQueue{}, item)
	if got.Class != ParkClassEnvironment || got.Next != "run roundfix deliver retry runtime-spec; a retry from runtime-infrastructure is not counted" {
		t.Fatalf("classification=%+v", got)
	}
}
func TestRuntimeInfrastructureRetryIsNotCounted(t *testing.T) {
	ctx := t.Context()
	s := openDeliveryEngineStore(t, ctx)
	const root = "/repo-runtime-retry"
	const slug = "runtime-retry"
	queue, err := s.CreateDeliveryQueueWithLimits(ctx, root, []string{slug}, store.DeliveryQueueLimits{MaxRetries: 1})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageParked
	item.Blocker = BlockerRunUnresolved
	item.Branch = "roundfix/deliver-" + slug
	item.Worktree = "/worktrees/" + slug
	item.RunID = "original-run"
	item.CandidateCommits = []string{"candidate-original"}
	recordDeliveryItemWorkspace(t, ctx, s, root, &item)
	if err := s.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
		t.Fatal(err)
	}
	target := item
	target.Stage = store.DeliveryStageRunning
	if _, _, err := s.RetryDeliveryQueueItem(ctx, root, target, item.Blocker); err != nil {
		t.Fatal(err)
	}
	item = readDeliveryQueue(t, ctx, s, root).Items[0]
	item.Stage = store.DeliveryStageParked
	item.Blocker = BlockerRuntimeInfrastructure + ": task_03 lost its rollout at session/prompt"
	if err := s.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
		t.Fatal(err)
	}
	recovery := &fakeItemRecovery{states: []ItemState{{Head: "candidate-original"}, {UnfinishedTasks: []string{"task_03"}, Head: "candidate-original"}}, carryResult: CarryForwardResult{RunID: "carried-run"}}
	engine := newRetryDeliveryEngine(s, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())
	if _, err := engine.Retry(ctx, root, slug); err != nil {
		t.Fatal(err)
	}
	got := readDeliveryQueue(t, ctx, s, root).Items[0]
	if got.RetryCount != 1 || got.Stage != store.DeliveryStageRunning || got.Blocker != "" || strings.Join(recovery.events, ",") != "inspect,carry-forward,inspect" {
		t.Fatalf("item=%+v recovery=%v", got, recovery.events)
	}
}
