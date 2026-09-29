// Suite: Delivery Revalidation.
// Invariant: a queued Spec is checked on its item worktree before its first Run, while premise changes warn without blocking.
// Boundary IN: the real SQLite Delivery Queue store and delivery engine transitions.
// Boundary OUT: item worktrees, revalidation, delivery stages, and GitHub are represented by fakes.
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

func TestRevalidationParksAnItemWhoseStartingMainFailsTheCheck(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const (
		gitRoot = "/repo-revalidation-finding"
		slug    = "finding-spec"
	)
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{slug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.recordWorkspace = func(branch, worktree string) error {
		if _, _, _, err := runStore.RecordDeliveryQueueItemWorktree(ctx, gitRoot, slug, branch, worktree); err != nil {
			return err
		}
		return runStore.SetDeliveryQueueItemWorktreeProvisioned(ctx, gitRoot, slug, true)
	}
	revalidator := &recordingItemRevalidator{results: map[string]Revalidation{
		slug: {Findings: []string{"SC-REF-UNRESOLVED", "SC-TRACE-MISSING"}},
	}}
	var log bytes.Buffer
	engine := newRevalidationTestEngine(runStore, workflow, revalidator, &log)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Stage != store.DeliveryStageParked || item.Blocker != "revalidation-failed: SC-REF-UNRESOLVED, SC-TRACE-MISSING" {
		t.Fatalf("revalidated item = %+v, want parked with finding codes", item)
	}
	if item.Branch == "" || item.Worktree == "" || !item.WorktreeProvisioned {
		t.Fatalf("parked item workspace = branch:%q worktree:%q provisioned:%t", item.Branch, item.Worktree, item.WorktreeProvisioned)
	}
	if workflow.events[slug][0] != "create-branch" || workflow.runs[slug].RunID != "" || slicesContain(workflow.events[slug], "run") {
		t.Fatalf("item events = %v, want worktree creation without Run", workflow.events[slug])
	}
	if log.Len() != 0 {
		t.Fatalf("revalidation finding log = %q, want empty", log.String())
	}
}

func TestAChangedPremiseWarnsAndContinues(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const (
		gitRoot = "/repo-revalidation-warning"
		slug    = "warning-spec"
	)
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{slug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	revalidator := &recordingItemRevalidator{results: map[string]Revalidation{
		slug: {
			ChangedPremises: []string{"internal/a.go", "internal/b.go"},
			ChangedBy:       []string{"merge-one", "merge-two"},
		},
	}}
	var log bytes.Buffer
	engine := newRevalidationTestEngine(runStore, workflow, revalidator, &log)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	warning := "premise-changed: internal/a.go, internal/b.go (merge merge-one, merge-two)"
	if item.Warning != warning || item.Stage != store.DeliveryStageMerged || !slicesContain(workflow.events[slug], "run") {
		t.Fatalf("warned item = %+v events=%v, want warning and completed Run path", item, workflow.events[slug])
	}
	wantLog := "roundfix: warning: Delivery Queue item warning-spec: " + warning + "\n"
	if log.String() != wantLog {
		t.Fatalf("warning log = %q, want %q", log.String(), wantLog)
	}
}

func TestNoChangedPremiseRecordsNoWarning(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const (
		gitRoot = "/repo-revalidation-clean"
		slug    = "clean-spec"
	)
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{slug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	var log bytes.Buffer
	engine := newRevalidationTestEngine(runStore, workflow, &recordingItemRevalidator{}, &log)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Warning != "" || item.Stage != store.DeliveryStageMerged || !slicesContain(workflow.events[slug], "run") {
		t.Fatalf("clean item = %+v events=%v, want merged without warning", item, workflow.events[slug])
	}
	if log.Len() != 0 {
		t.Fatalf("clean revalidation log = %q, want empty", log.String())
	}
}

func TestRevalidationReceivesEarlierMergeCommitsInQueueOrder(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-revalidation-prior-merges"
	queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"merged-one", "will-park", "merged-two", "target"})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	for index, mergeCommit := range map[int]string{0: "merge-one", 2: "merge-two"} {
		item := queue.Items[index]
		item.Stage = store.DeliveryStageMerged
		item.MergeCommit = mergeCommit
		if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
			t.Fatalf("seed merged item %d: %v", index, err)
		}
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.runs["will-park"] = RunResult{RunID: "run-will-park", Outcome: RunOutcomeUnresolved}
	revalidator := &recordingItemRevalidator{}
	engine := newRevalidationTestEngine(runStore, workflow, revalidator, &bytes.Buffer{})

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	if got, want := revalidator.priorMerges["target"], []string{"merge-one", "merge-two"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("target prior merges = %v, want %v", got, want)
	}
	if got := revalidator.priorMerges["will-park"]; !reflect.DeepEqual(got, []string{"merge-one"}) {
		t.Fatalf("unfinished item prior merges = %v, want only earlier merged item", got)
	}
}

func TestARevalidationErrorParksTheItemAsADeliveryError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const (
		gitRoot = "/repo-revalidation-error"
		slug    = "error-spec"
	)
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{slug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	revalidator := &recordingItemRevalidator{errors: map[string]error{slug: errors.New("cannot read merge")}}
	engine := newRevalidationTestEngine(runStore, workflow, revalidator, &bytes.Buffer{})

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Stage != store.DeliveryStageParked || !strings.HasPrefix(item.Blocker, "delivery-error: revalidate item: cannot read merge") {
		t.Fatalf("errored item = %+v, want delivery-error", item)
	}
	if slicesContain(workflow.events[slug], "run") {
		t.Fatalf("item events = %v, want no Run", workflow.events[slug])
	}
}

func TestRetryOfAnItemWithoutARunRefusesWhileFindingsRemain(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-revalidation-retry-refused"
	item := seedNoRunRetryItem(t, ctx, runStore, gitRoot, "retry-refused")
	before := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	recovery := &fakeItemRecovery{states: []ItemState{{UnfinishedTasks: []string{"task_01"}, Head: "candidate"}}}
	revalidator := &recordingItemRevalidator{results: map[string]Revalidation{
		item.SpecSlug: {Findings: []string{"SC-REF-UNRESOLVED", "SC-TASK-UNCOVERED"}},
	}}
	engine := newRetryRevalidationTestEngine(runStore, recovery, revalidator)

	_, err := engine.Retry(ctx, gitRoot, item.SpecSlug)
	if err == nil || !strings.Contains(err.Error(), "SC-REF-UNRESOLVED, SC-TASK-UNCOVERED") {
		t.Fatalf("retry finding error = %v, want both codes", err)
	}
	if got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]; !reflect.DeepEqual(got, before) {
		t.Fatalf("item changed after refused revalidation:\n got: %#v\nwant: %#v", got, before)
	}
	if slicesContain(recovery.events, "carry-forward") {
		t.Fatalf("recovery events = %v, want no carry-forward", recovery.events)
	}
}

func TestRetryOfAnItemWithoutARunProceedsOnceTheCheckIsClean(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-revalidation-retry-clean"
	item := seedNoRunRetryItem(t, ctx, runStore, gitRoot, "retry-clean")
	recovery := &fakeItemRecovery{
		states: []ItemState{
			{UnfinishedTasks: []string{"task_01"}, Head: "candidate"},
			{UnfinishedTasks: []string{"task_01"}, Head: "candidate"},
		},
		carryResult: CarryForwardResult{RunID: "run-retried"},
	}
	revalidator := &recordingItemRevalidator{}
	engine := newRetryRevalidationTestEngine(runStore, recovery, revalidator)

	result, err := engine.Retry(ctx, gitRoot, item.SpecSlug)
	if err != nil {
		t.Fatalf("retry Delivery Queue item: %v", err)
	}
	if result.Stage != store.DeliveryStageRunning || !reflect.DeepEqual(recovery.events, []string{"inspect", "carry-forward", "inspect"}) {
		t.Fatalf("retry result = %+v recovery events=%v, want running after carry-forward", result, recovery.events)
	}
	if got := revalidator.priorMerges[item.SpecSlug]; len(got) != 0 {
		t.Fatalf("retry prior merges = %v, want none", got)
	}
}

func TestRetryKeepsTheRecordedPremiseWarning(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-revalidation-retry-warning"
	item := seedNoRunRetryItem(t, ctx, runStore, gitRoot, "retry-warning")
	item.Warning = "premise-changed: internal/old.go (merge prior-merge)"
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("seed warning: %v", err)
	}
	recovery := &fakeItemRecovery{
		states:      []ItemState{{UnfinishedTasks: []string{"task_01"}}, {UnfinishedTasks: []string{"task_01"}}},
		carryResult: CarryForwardResult{RunID: "run-retried"},
	}
	revalidator := &recordingItemRevalidator{results: map[string]Revalidation{
		item.SpecSlug: {ChangedPremises: []string{"internal/new.go"}, ChangedBy: []string{"ignored"}},
	}}
	engine := newRetryRevalidationTestEngine(runStore, recovery, revalidator)

	if _, err := engine.Retry(ctx, gitRoot, item.SpecSlug); err != nil {
		t.Fatalf("retry Delivery Queue item: %v", err)
	}
	if got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0].Warning; got != item.Warning {
		t.Fatalf("retried item warning = %q, want unchanged %q", got, item.Warning)
	}
}

func TestEngineRefusesToRunWithoutARevalidator(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	engine := NewEngine(runStore, EngineDependencies{
		Workspace:    workflow,
		Runner:       workflow,
		Reviewer:     workflow,
		Archiver:     workflow,
		Gate:         workflow,
		Authorizer:   workflow,
		Publication:  workflow,
		PullRequests: boundary,
		Recovery:     &fakeItemRecovery{},
	})

	if _, err := engine.Run(ctx, "/repo"); err == nil || !strings.Contains(err.Error(), "item revalidator is required") {
		t.Fatalf("Run without revalidator error = %v", err)
	}
	if _, err := engine.Retry(ctx, "/repo", "spec"); err == nil || !strings.Contains(err.Error(), "item revalidator is required") {
		t.Fatalf("Retry without revalidator error = %v", err)
	}
}

type recordingItemRevalidator struct {
	results     map[string]Revalidation
	errors      map[string]error
	priorMerges map[string][]string
}

func (fake *recordingItemRevalidator) Revalidate(
	_ context.Context,
	_ string,
	specSlug string,
	priorMerges []string,
) (Revalidation, error) {
	if fake.priorMerges == nil {
		fake.priorMerges = make(map[string][]string)
	}
	fake.priorMerges[specSlug] = append([]string(nil), priorMerges...)
	return fake.results[specSlug], fake.errors[specSlug]
}

func newRevalidationTestEngine(
	runStore *store.Store,
	workflow *fakeDeliveryWorkflow,
	revalidator ItemRevalidator,
	log *bytes.Buffer,
) *Engine {
	boundary := newFakeDeliveryBoundary()
	boundary.recordEvent = workflow.recordEvent
	return NewEngine(runStore, EngineDependencies{
		Workspace:    workflow,
		Runner:       workflow,
		Reviewer:     workflow,
		Archiver:     workflow,
		Gate:         workflow,
		Authorizer:   workflow,
		Publication:  workflow,
		PullRequests: boundary,
		Revalidator:  revalidator,
		Log:          log,
	})
}

func newRetryRevalidationTestEngine(
	runStore *store.Store,
	recovery ItemRecovery,
	revalidator ItemRevalidator,
) *Engine {
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	return NewEngine(runStore, EngineDependencies{
		Workspace:    workflow,
		Runner:       workflow,
		Reviewer:     workflow,
		Archiver:     workflow,
		Gate:         workflow,
		Authorizer:   workflow,
		Publication:  workflow,
		PullRequests: boundary,
		Recovery:     recovery,
		Revalidator:  revalidator,
	})
}

func seedNoRunRetryItem(
	t *testing.T,
	ctx context.Context,
	runStore *store.Store,
	gitRoot string,
	specSlug string,
) store.DeliveryQueueItem {
	t.Helper()
	item := seedParkedRetryItem(t, ctx, runStore, gitRoot, specSlug, BlockerRevalidationFailed+": SC-REF-UNRESOLVED")
	item.RunID = ""
	item.CandidateCommits = nil
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("clear parked item Run evidence: %v", err)
	}
	return item
}

func slicesContain(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
