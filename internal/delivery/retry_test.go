// Suite: Delivery Queue item recovery.
// Invariant: a retry derives its stage from item evidence and never mutates a refused queue item.
// Boundary IN: the real SQLite Delivery Queue store and delivery engine transitions.
// Boundary OUT: item inspection, Task Carry-Forward, Git worktrees, and GitHub are represented by fakes.
package delivery

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestAnUnresolvedRunParksWithItsRunID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const (
		gitRoot  = "/repo-unresolved-run-id"
		specSlug = "unresolved-spec"
	)
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{specSlug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.runs[specSlug] = RunResult{RunID: "  run-unresolved-42  ", Outcome: RunOutcomeUnresolved}
	workflow.recordWorkspace = func(branch, worktree string) error {
		_, _, _, err := runStore.RecordDeliveryQueueItemWorktree(ctx, gitRoot, specSlug, branch, worktree)
		return err
	}
	engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}
	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Stage != store.DeliveryStageParked || item.Blocker != BlockerRunUnresolved || item.RunID != "run-unresolved-42" {
		t.Fatalf("unresolved item = %+v, want parked with its trimmed Run ID", item)
	}
}

func TestRetryCarriesForwardBeforeReenteringTheRun(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	item := seedParkedRetryItem(t, ctx, runStore, "/repo-retry-order", "retry-order", BlockerRunUnresolved)
	recovery := &fakeItemRecovery{
		states: []ItemState{
			{Head: "candidate-original"},
			{UnfinishedTasks: []string{"task_02"}, Head: "candidate-original"},
		},
		carryResult: CarryForwardResult{RunID: "run-carried", Carried: []string{"task_01"}},
	}
	engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())

	if _, err := engine.Retry(ctx, "/repo-retry-order", item.SpecSlug); err != nil {
		t.Fatalf("retry Delivery Queue item: %v", err)
	}
	if got, want := recovery.events, []string{"inspect", "carry-forward", "inspect"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("recovery calls = %v, want %v", got, want)
	}
	if recovery.carryWorkDir != item.Worktree || recovery.carrySpecSlug != item.SpecSlug || recovery.carryBranch != item.Branch || recovery.carryRunID != item.RunID {
		t.Fatalf("carry-forward arguments = workdir:%q slug:%q branch:%q run:%q", recovery.carryWorkDir, recovery.carrySpecSlug, recovery.carryBranch, recovery.carryRunID)
	}
}

func TestRetryReentersTheRunWhenATaskIsUnfinished(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-retry-running"
	item := seedParkedRetryItem(t, ctx, runStore, gitRoot, "retry-running", BlockerRunUnresolved)
	recovery := &fakeItemRecovery{
		states: []ItemState{
			{Head: "candidate-original"},
			{UnfinishedTasks: []string{"task_03"}, Head: "current-head"},
		},
		carryResult: CarryForwardResult{RunID: " run-retried ", Carried: []string{"task_01", "task_02"}},
	}
	engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())

	result, err := engine.Retry(ctx, gitRoot, item.SpecSlug)
	if err != nil {
		t.Fatalf("retry Delivery Queue item: %v", err)
	}
	got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if result.Stage != store.DeliveryStageRunning || got.Stage != store.DeliveryStageRunning || got.Blocker != "" || got.RunID != "run-retried" {
		t.Fatalf("retried running item = result:%+v item:%+v", result, got)
	}
	if !reflect.DeepEqual(got.CandidateCommits, []string{"candidate-original"}) {
		t.Fatalf("running candidate commits = %v, want unchanged", got.CandidateCommits)
	}
}

func TestRetryReReviewsTheCurrentHeadWhenNoTaskIsUnfinished(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-retry-reviewing"
	item := seedParkedRetryItem(t, ctx, runStore, gitRoot, "retry-reviewing", BlockerReviewFindings)
	recovery := &fakeItemRecovery{
		states:      []ItemState{{Head: "candidate-original"}, {Head: "current-head"}},
		carryResult: CarryForwardResult{RunID: "run-review", Carried: []string{"task_01"}},
	}
	engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())

	result, err := engine.Retry(ctx, gitRoot, item.SpecSlug)
	if err != nil {
		t.Fatalf("retry Delivery Queue item: %v", err)
	}
	got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if result.Stage != store.DeliveryStageReviewing || got.Stage != store.DeliveryStageReviewing {
		t.Fatalf("retried review stage = result:%q item:%q", result.Stage, got.Stage)
	}
	if want := []string{"candidate-original", "current-head"}; !reflect.DeepEqual(got.CandidateCommits, want) {
		t.Fatalf("review candidate commits = %v, want %v", got.CandidateCommits, want)
	}
}

func TestRetryAfterArchiveReentersGatingWithoutAPullRequest(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-retry-gating"
	item := seedParkedRetryItem(t, ctx, runStore, gitRoot, "retry-gating", BlockerGateFailed)
	recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "candidate-original"}}}
	engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())

	result, err := engine.Retry(ctx, gitRoot, item.SpecSlug)
	if err != nil {
		t.Fatalf("retry archived Delivery Queue item: %v", err)
	}
	if result.Stage != store.DeliveryStageGating || readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0].Stage != store.DeliveryStageGating {
		t.Fatalf("archived retry result = %+v, want gating", result)
	}
	if len(recovery.events) != 1 || recovery.events[0] != "inspect" {
		t.Fatalf("archived recovery calls = %v, want inspect only", recovery.events)
	}
}

func TestRetryAfterArchiveReentersCheckingWithAPullRequest(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-retry-checking"
	item := seedParkedRetryItem(t, ctx, runStore, gitRoot, "retry-checking", BlockerChecksFailed)
	item.PullRequestNumber = "42"
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("record pull request: %v", err)
	}
	recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "candidate-original"}}}
	engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())

	result, err := engine.Retry(ctx, gitRoot, item.SpecSlug)
	if err != nil {
		t.Fatalf("retry archived Delivery Queue item: %v", err)
	}
	if result.Stage != store.DeliveryStageChecking || readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0].Stage != store.DeliveryStageChecking {
		t.Fatalf("archived retry result = %+v, want checking", result)
	}
}

func TestRetriedArchivedItemPublishesAndMergesOnce(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-retry-publish-once"
	item := seedParkedRetryItem(t, ctx, runStore, gitRoot, "retry-publish-once", BlockerGateFailed)
	recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "candidate-original"}}}
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	engine := newRetryDeliveryEngine(runStore, workflow, nil, recovery, boundary)

	if _, err := engine.Retry(ctx, gitRoot, item.SpecSlug); err != nil {
		t.Fatalf("retry archived Delivery Queue item: %v", err)
	}
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run retried Delivery Queue item: %v", err)
	}
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("resume merged Delivery Queue item: %v", err)
	}
	got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if got.Stage != store.DeliveryStageMerged {
		t.Fatalf("retried archived item = %+v, want merged", got)
	}
	if boundary.pushEffects != 1 || boundary.pullRequestEffects != 1 || boundary.mergeEffects != 1 {
		t.Fatalf("external effects = push:%d pull-request:%d merge:%d, want one each", boundary.pushEffects, boundary.pullRequestEffects, boundary.mergeEffects)
	}
}

func TestRetryRefusesAnArchivedItemWhoseHeadMoved(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-retry-moved-archive"
	before := seedParkedRetryItem(t, ctx, runStore, gitRoot, "retry-moved-archive", BlockerReviewStale)
	recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "moved-head"}}}
	engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())

	_, err := engine.Retry(ctx, gitRoot, before.SpecSlug)
	if err == nil || !strings.Contains(err.Error(), "moved-head") || !strings.Contains(err.Error(), "candidate-original") {
		t.Fatalf("moved archived head error = %v, want both heads", err)
	}
	assertRetryItemUnchanged(t, ctx, runStore, gitRoot, before)
}

func TestRetryRefusesWhenCarryForwardRefuses(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-retry-carry-refusal"
	before := seedParkedRetryItem(t, ctx, runStore, gitRoot, "retry-carry-refusal", BlockerRunUnresolved)
	carryErr := errors.New("declared input moved")
	recovery := &fakeItemRecovery{states: []ItemState{{Head: "candidate-original"}}, carryErr: carryErr}
	engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())

	_, err := engine.Retry(ctx, gitRoot, before.SpecSlug)
	if !errors.Is(err, carryErr) {
		t.Fatalf("carry-forward refusal = %v, want wrapped sentinel", err)
	}
	assertRetryItemUnchanged(t, ctx, runStore, gitRoot, before)
}

func TestRetryRefusesAnItemThatIsNotParked(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-retry-not-parked"
	before := seedParkedRetryItem(t, ctx, runStore, gitRoot, "retry-not-parked", BlockerReviewFindings)
	before.Stage = store.DeliveryStageReviewing
	before.Blocker = ""
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, before); err != nil {
		t.Fatalf("seed reviewing item: %v", err)
	}
	recovery := &fakeItemRecovery{states: []ItemState{{Head: "candidate-original"}}}
	engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())

	_, err := engine.Retry(ctx, gitRoot, before.SpecSlug)
	if err == nil || !strings.Contains(err.Error(), `stage "reviewing"`) {
		t.Fatalf("non-parked retry error = %v, want stage refusal", err)
	}
	if len(recovery.events) != 0 {
		t.Fatalf("non-parked retry called recovery: %v", recovery.events)
	}
	assertRetryItemUnchanged(t, ctx, runStore, gitRoot, before)
}

func TestRetryRefusesAnItemWhoseBranchIsGone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-retry-missing-branch"
	before := seedParkedRetryItem(t, ctx, runStore, gitRoot, "retry-missing-branch", BlockerItemWorktreeMissing)
	workflow := newFakeDeliveryWorkflow()
	workspace := &missingRetryWorkspace{fakeDeliveryWorkflow: workflow, err: ErrItemWorktreeMissing}
	recovery := &fakeItemRecovery{states: []ItemState{{Head: "candidate-original"}}}
	engine := newRetryDeliveryEngine(runStore, workflow, workspace, recovery, newFakeDeliveryBoundary())

	_, err := engine.Retry(ctx, gitRoot, before.SpecSlug)
	if err == nil || !strings.Contains(err.Error(), before.Branch) {
		t.Fatalf("missing item branch error = %v, want branch %q", err, before.Branch)
	}
	if !errors.Is(err, ErrItemWorktreeMissing) {
		t.Fatalf("missing item branch error = %v, want ErrItemWorktreeMissing", err)
	}
	assertRetryItemUnchanged(t, ctx, runStore, gitRoot, before)
}

func TestRetryRefusesAnItemMissingFromTheQueue(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-retry-missing-item"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"other-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	before := readDeliveryQueue(t, ctx, runStore, gitRoot)
	recovery := &fakeItemRecovery{states: []ItemState{{Head: "candidate-original"}}}
	engine := newRetryDeliveryEngine(runStore, newFakeDeliveryWorkflow(), nil, recovery, newFakeDeliveryBoundary())

	_, err := engine.Retry(ctx, gitRoot, "missing-spec")
	if err == nil || !strings.Contains(err.Error(), `"missing-spec"`) || !strings.Contains(err.Error(), "does not contain") {
		t.Fatalf("missing queue item error = %v, want named refusal", err)
	}
	if got := readDeliveryQueue(t, ctx, runStore, gitRoot); !reflect.DeepEqual(got, before) {
		t.Fatalf("queue changed after missing-item refusal:\n got: %#v\nwant: %#v", got, before)
	}
}

func seedParkedRetryItem(
	t *testing.T,
	ctx context.Context,
	runStore *store.Store,
	gitRoot string,
	specSlug string,
	blocker string,
) store.DeliveryQueueItem {
	t.Helper()
	queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{specSlug})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageParked
	item.Blocker = blocker
	item.Branch = "roundfix/deliver-" + specSlug
	item.Worktree = "/worktrees/" + specSlug
	item.RunID = "run-original"
	item.CandidateCommits = []string{"candidate-original"}
	recordDeliveryItemWorkspace(t, ctx, runStore, gitRoot, &item)
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("seed parked Delivery Queue item: %v", err)
	}
	return item
}

func assertRetryItemUnchanged(
	t *testing.T,
	ctx context.Context,
	runStore *store.Store,
	gitRoot string,
	want store.DeliveryQueueItem,
) {
	t.Helper()
	if got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]; !reflect.DeepEqual(got, want) {
		t.Fatalf("Delivery Queue item changed after retry refusal:\n got: %#v\nwant: %#v", got, want)
	}
}

func newRetryDeliveryEngine(
	runStore *store.Store,
	workflow *fakeDeliveryWorkflow,
	workspace ItemWorkspace,
	recovery ItemRecovery,
	boundary *fakeDeliveryBoundary,
) *Engine {
	if workspace == nil {
		workspace = workflow
	}
	boundary.recordEvent = workflow.recordEvent
	return NewEngine(runStore, EngineDependencies{
		Workspace:    workspace,
		Runner:       workflow,
		Reviewer:     workflow,
		Archiver:     workflow,
		Gate:         workflow,
		Authorizer:   workflow,
		Publication:  workflow,
		PullRequests: boundary,
		Recovery:     recovery,
		Revalidator:  cleanTestRevalidator{},
	})
}

type fakeItemRecovery struct {
	states        []ItemState
	inspectCalls  int
	carryResult   CarryForwardResult
	carryErr      error
	events        []string
	carryWorkDir  string
	carrySpecSlug string
	carryBranch   string
	carryRunID    string
}

func (fake *fakeItemRecovery) InspectItem(context.Context, string, string) (ItemState, error) {
	fake.events = append(fake.events, "inspect")
	if len(fake.states) == 0 {
		return ItemState{}, errors.New("no fake item state")
	}
	index := min(fake.inspectCalls, len(fake.states)-1)
	fake.inspectCalls++
	return fake.states[index], nil
}

func (fake *fakeItemRecovery) CarryForward(
	_ context.Context,
	workDir string,
	specSlug string,
	branch string,
	runID string,
) (CarryForwardResult, error) {
	fake.events = append(fake.events, "carry-forward")
	fake.carryWorkDir = workDir
	fake.carrySpecSlug = specSlug
	fake.carryBranch = branch
	fake.carryRunID = runID
	if fake.carryErr != nil {
		return CarryForwardResult{}, fake.carryErr
	}
	return fake.carryResult, nil
}

type missingRetryWorkspace struct {
	*fakeDeliveryWorkflow
	err error
}

func (fake *missingRetryWorkspace) UseItemBranch(
	context.Context,
	string,
	string,
	string,
	string,
	bool,
) (string, error) {
	return "", fake.err
}

var _ ItemRecovery = (*fakeItemRecovery)(nil)
var _ ItemWorkspace = (*missingRetryWorkspace)(nil)
