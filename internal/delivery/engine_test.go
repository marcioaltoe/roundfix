// Suite: durable Spec delivery engine.
// Invariant: every queued Spec advances in order without replaying an observed external effect or publishing an unreviewed head.
// Boundary IN: the real SQLite Delivery Queue store and the delivery engine's stage transitions.
// Boundary OUT: Implement, review, archive, repository gate, authorization files, and GitHub, represented by boundary fakes.
package delivery

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"roundfix/internal/store"
)

func TestDeliveryEngineMergesAQueueWithoutAnOperator(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"spec-one", "spec-two"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}

	workflow := newFakeDeliveryWorkflow()
	workflow.policies["spec-two"] = ReviewPolicyNone
	boundary := newFakeDeliveryBoundary()
	boundary.recordEvent = workflow.recordEvent
	engine := NewEngine(runStore, EngineDependencies{
		Runner:       workflow,
		Reviewer:     workflow,
		Archiver:     workflow,
		Gate:         workflow,
		Authorizer:   workflow,
		Publication:  workflow,
		PullRequests: boundary,
	})
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	queue := readDeliveryQueue(t, ctx, runStore, gitRoot)
	for _, item := range queue.Items {
		if item.Stage != store.DeliveryStageMerged || item.MergeCommit == "" {
			t.Fatalf("item %q = %+v, want merged with a merge commit", item.SpecSlug, item)
		}
	}
	if got, want := workflow.reviews, []string{"spec-one"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("review calls = %v, want %v", got, want)
	}
	if got, want := workflow.omissions, []string{"spec-two"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("recorded review omissions = %v, want %v", got, want)
	}
	if boundary.pullRequestEffects != 2 || boundary.mergeEffects != 2 {
		t.Fatalf("external effects = pull_requests:%d merges:%d, want 2 and 2", boundary.pullRequestEffects, boundary.mergeEffects)
	}
	assertNoUnmatchedDeliveryIntents(t, ctx, runStore, gitRoot)
	wantReviewed := []string{
		"run", "policy", "review", "archive", "gate", "authorization",
		"publication", "push", "pull-request", "checks", "authorization", "merge",
	}
	if got := workflow.events["spec-one"]; !reflect.DeepEqual(got, wantReviewed) {
		t.Fatalf("reviewed item actions = %v, want %v", got, wantReviewed)
	}
	wantOmitted := []string{
		"run", "policy", "review-omitted", "archive", "gate", "authorization",
		"publication", "push", "pull-request", "checks", "authorization", "merge",
	}
	if got := workflow.events["spec-two"]; !reflect.DeepEqual(got, wantOmitted) {
		t.Fatalf("omitted-review item actions = %v, want %v", got, wantOmitted)
	}
}

func TestDeliveryEngineResumesWithoutDoubleEffects(t *testing.T) {
	t.Parallel()

	t.Run("push intent", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		runStore := openDeliveryEngineStore(t, ctx)
		const gitRoot = "/repo-push"
		queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"spec-push"})
		if err != nil {
			t.Fatalf("create Delivery Queue: %v", err)
		}
		item := queue.Items[0]
		item.Stage = store.DeliveryStagePublishing
		item.CandidateCommits = []string{"reviewed-spec-push", "archived-spec-push"}
		if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
			t.Fatalf("seed publishing item: %v", err)
		}
		if _, err := runStore.RecordDeliveryActionIntent(ctx, gitRoot, item.SpecSlug, store.DeliveryActionPush); err != nil {
			t.Fatalf("record unmatched push intent: %v", err)
		}

		workflow := newFakeDeliveryWorkflow()
		boundary := newFakeDeliveryBoundary()
		boundary.remoteHeads["deliver/spec-push"] = "archived-spec-push"
		engine := newTestDeliveryEngine(runStore, workflow, boundary)
		if _, err := engine.Run(ctx, gitRoot); err != nil {
			t.Fatalf("resume Delivery Engine: %v", err)
		}

		got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
		if got.Stage != store.DeliveryStageMerged {
			t.Fatalf("resumed item stage = %q, want merged", got.Stage)
		}
		if boundary.pushEffects != 0 || boundary.pullRequestEffects != 1 || boundary.mergeEffects != 1 {
			t.Fatalf("resume effects = pushes:%d pull_requests:%d merges:%d, want 0, 1, 1", boundary.pushEffects, boundary.pullRequestEffects, boundary.mergeEffects)
		}
		assertNoUnmatchedDeliveryIntents(t, ctx, runStore, gitRoot)
	})

	t.Run("pull request intent", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		runStore := openDeliveryEngineStore(t, ctx)
		const gitRoot = "/repo-pull-request"
		queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"spec-pull-request"})
		if err != nil {
			t.Fatalf("create Delivery Queue: %v", err)
		}
		item := queue.Items[0]
		item.Stage = store.DeliveryStagePublishing
		item.CandidateCommits = []string{"reviewed-spec-pull-request", "archived-spec-pull-request"}
		if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
			t.Fatalf("seed publishing item: %v", err)
		}
		if _, err := runStore.RecordDeliveryActionIntent(ctx, gitRoot, item.SpecSlug, store.DeliveryActionCreatePullRequest); err != nil {
			t.Fatalf("record unmatched create-pull-request intent: %v", err)
		}

		workflow := newFakeDeliveryWorkflow()
		boundary := newFakeDeliveryBoundary()
		branch := "deliver/spec-pull-request"
		boundary.remoteHeads[branch] = "archived-spec-pull-request"
		boundary.pullRequests[branch] = PullRequest{Number: "1", State: "OPEN", HeadBranch: branch, HeadSHA: "archived-spec-pull-request"}
		boundary.pullRequestEffects = 1
		engine := newTestDeliveryEngine(runStore, workflow, boundary)
		if _, err := engine.Run(ctx, gitRoot); err != nil {
			t.Fatalf("resume Delivery Engine: %v", err)
		}

		got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
		if got.Stage != store.DeliveryStageMerged {
			t.Fatalf("resumed item stage = %q, want merged", got.Stage)
		}
		if boundary.pullRequestEffects != 1 || boundary.mergeEffects != 1 {
			t.Fatalf("resume effects = pull_requests:%d merges:%d, want 1 total each", boundary.pullRequestEffects, boundary.mergeEffects)
		}
		assertNoUnmatchedDeliveryIntents(t, ctx, runStore, gitRoot)
	})

	t.Run("merge intent", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		runStore := openDeliveryEngineStore(t, ctx)
		const gitRoot = "/repo-merge"
		queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"spec-merge"})
		if err != nil {
			t.Fatalf("create Delivery Queue: %v", err)
		}
		item := queue.Items[0]
		item.Stage = store.DeliveryStageMerging
		item.CandidateCommits = []string{"reviewed-spec-merge", "archived-spec-merge"}
		item.PullRequestNumber = "1"
		if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
			t.Fatalf("seed merging item: %v", err)
		}
		if _, err := runStore.RecordDeliveryActionIntent(ctx, gitRoot, item.SpecSlug, store.DeliveryActionMerge); err != nil {
			t.Fatalf("record unmatched merge intent: %v", err)
		}

		workflow := newFakeDeliveryWorkflow()
		boundary := newFakeDeliveryBoundary()
		boundary.merged["1"] = "merge-1"
		boundary.mergeEffects = 1
		engine := newTestDeliveryEngine(runStore, workflow, boundary)
		if _, err := engine.Run(ctx, gitRoot); err != nil {
			t.Fatalf("resume Delivery Engine: %v", err)
		}

		got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
		if got.Stage != store.DeliveryStageMerged || got.MergeCommit != "merge-1" {
			t.Fatalf("resumed item = %+v, want existing merge receipt", got)
		}
		if boundary.pullRequestEffects != 0 || boundary.mergeEffects != 1 {
			t.Fatalf("resume effects = pull_requests:%d merges:%d, want 0 and 1 total", boundary.pullRequestEffects, boundary.mergeEffects)
		}
		assertNoUnmatchedDeliveryIntents(t, ctx, runStore, gitRoot)
	})
}

func TestDeliveryEngineParksAStaleReview(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-stale"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"stale-spec", "next-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.archives["stale-spec"] = ArchiveResult{
		Parent:        "foreign-commit",
		Head:          "archived-stale-spec",
		ExactSpecMove: true,
	}
	boundary := newFakeDeliveryBoundary()
	engine := newTestDeliveryEngine(runStore, workflow, boundary)
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	queue := readDeliveryQueue(t, ctx, runStore, gitRoot)
	if queue.Items[0].Stage != store.DeliveryStageParked || queue.Items[0].Blocker != BlockerReviewStale {
		t.Fatalf("stale item = %+v, want parked as review-stale", queue.Items[0])
	}
	if queue.Items[1].Stage != store.DeliveryStageMerged {
		t.Fatalf("next item stage = %q, want merged", queue.Items[1].Stage)
	}
	if boundary.callsFor("stale-spec") != 0 {
		t.Fatalf("stale item made %d external calls, want none", boundary.callsFor("stale-spec"))
	}
}

func TestDeliveryEngineParksUnauthorizedBeforeExternalAction(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-unauthorized"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"unauthorized-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.authorizations["unauthorized-spec"] = Authorization{Operations: []string{"implement", "push", "pull_request"}}
	boundary := newFakeDeliveryBoundary()
	engine := newTestDeliveryEngine(runStore, workflow, boundary)
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Stage != store.DeliveryStageParked || item.Blocker != BlockerUnauthorized {
		t.Fatalf("unauthorized item = %+v, want parked as unauthorized", item)
	}
	if boundary.totalCalls != 0 {
		t.Fatalf("unauthorized item made %d external calls, want none", boundary.totalCalls)
	}
}

func TestDeliveryEngineParksDeclaredBlockersAndContinues(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-blockers"
	slugs := []string{"run-blocked", "review-findings", "review-blocked", "gate-blocked", "checks-blocked", "after-blockers"}
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, slugs); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.runs["run-blocked"] = RunResult{Outcome: RunOutcomeUnresolved, Reason: "Run ended Unresolved"}
	workflow.reviewResults["review-findings"] = ReviewResult{Outcome: ReviewOutcomeFindings, Head: "reviewed-review-findings", Reason: "finding"}
	workflow.reviewResults["review-blocked"] = ReviewResult{Outcome: ReviewOutcomeBlocked, Head: "reviewed-review-blocked", Reason: "review unavailable"}
	workflow.gates["gate-blocked"] = GateResult{Passed: false, Reason: "repository gate failed"}
	boundary := newFakeDeliveryBoundary()
	boundary.failedChecks["checks-blocked"] = true
	engine := newTestDeliveryEngine(runStore, workflow, boundary)
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	queue := readDeliveryQueue(t, ctx, runStore, gitRoot)
	for index, blocker := range []string{BlockerRunUnresolved, BlockerReviewFindings, BlockerReviewBlocked, BlockerGateFailed, BlockerChecksFailed} {
		if queue.Items[index].Stage != store.DeliveryStageParked || queue.Items[index].Blocker != blocker {
			t.Fatalf("item %q = %+v, want parked as %q", queue.Items[index].SpecSlug, queue.Items[index], blocker)
		}
	}
	if queue.Items[len(queue.Items)-1].Stage != store.DeliveryStageMerged {
		t.Fatalf("item after blockers stage = %q, want merged", queue.Items[len(queue.Items)-1].Stage)
	}
}

func newTestDeliveryEngine(runStore *store.Store, workflow *fakeDeliveryWorkflow, boundary *fakeDeliveryBoundary) *Engine {
	boundary.recordEvent = workflow.recordEvent
	return NewEngine(runStore, EngineDependencies{
		Runner:       workflow,
		Reviewer:     workflow,
		Archiver:     workflow,
		Gate:         workflow,
		Authorizer:   workflow,
		Publication:  workflow,
		PullRequests: boundary,
	})
}

func openDeliveryEngineStore(t *testing.T, ctx context.Context) *store.Store {
	t.Helper()
	runStore, err := store.Open(ctx, t.TempDir())
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	t.Cleanup(func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close Run Database: %v", err)
		}
	})
	return runStore
}

func readDeliveryQueue(t *testing.T, ctx context.Context, runStore *store.Store, gitRoot string) store.DeliveryQueue {
	t.Helper()
	queue, found, err := runStore.DeliveryQueue(ctx, gitRoot)
	if err != nil || !found {
		t.Fatalf("read Delivery Queue: found=%v err=%v", found, err)
	}
	return queue
}

func assertNoUnmatchedDeliveryIntents(t *testing.T, ctx context.Context, runStore *store.Store, gitRoot string) {
	t.Helper()
	intents, err := runStore.UnmatchedDeliveryActionIntents(ctx, gitRoot)
	if err != nil {
		t.Fatalf("list unmatched Delivery Action intents: %v", err)
	}
	if len(intents) != 0 {
		t.Fatalf("unmatched Delivery Action intents = %#v, want none", intents)
	}
}

type fakeDeliveryWorkflow struct {
	runs           map[string]RunResult
	policies       map[string]ReviewPolicy
	reviewResults  map[string]ReviewResult
	archives       map[string]ArchiveResult
	gates          map[string]GateResult
	authorizations map[string]Authorization
	publications   map[string]Publication
	reviews        []string
	omissions      []string
	events         map[string][]string
}

func newFakeDeliveryWorkflow() *fakeDeliveryWorkflow {
	return &fakeDeliveryWorkflow{
		runs:           map[string]RunResult{},
		policies:       map[string]ReviewPolicy{},
		reviewResults:  map[string]ReviewResult{},
		archives:       map[string]ArchiveResult{},
		gates:          map[string]GateResult{},
		authorizations: map[string]Authorization{},
		publications:   map[string]Publication{},
		events:         map[string][]string{},
	}
}

func (fake *fakeDeliveryWorkflow) RunSpec(_ context.Context, _ string, slug string) (RunResult, error) {
	fake.recordEvent(slug, "run")
	if result, ok := fake.runs[slug]; ok {
		return result, nil
	}
	return RunResult{RunID: "run-" + slug, Outcome: RunOutcomeClean, CandidateCommits: []string{"reviewed-" + slug}}, nil
}

func (fake *fakeDeliveryWorkflow) ReviewPolicy(_ context.Context, _ string, slug string) (ReviewPolicy, error) {
	fake.recordEvent(slug, "policy")
	if policy, ok := fake.policies[slug]; ok {
		return policy, nil
	}
	return ReviewPolicyEnabled, nil
}

func (fake *fakeDeliveryWorkflow) Review(_ context.Context, _ string, slug, head string) (ReviewResult, error) {
	fake.recordEvent(slug, "review")
	fake.reviews = append(fake.reviews, slug)
	if result, ok := fake.reviewResults[slug]; ok {
		return result, nil
	}
	return ReviewResult{Outcome: ReviewOutcomeReviewed, Head: head}, nil
}

func (fake *fakeDeliveryWorkflow) RecordReviewOmission(_ context.Context, _ string, slug, _ string) error {
	fake.recordEvent(slug, "review-omitted")
	fake.omissions = append(fake.omissions, slug)
	return nil
}

func (fake *fakeDeliveryWorkflow) Archive(_ context.Context, _ string, slug, reviewedHead string) (ArchiveResult, error) {
	fake.recordEvent(slug, "archive")
	if result, ok := fake.archives[slug]; ok {
		return result, nil
	}
	return ArchiveResult{Parent: reviewedHead, Head: "archived-" + slug, ExactSpecMove: true}, nil
}

func (fake *fakeDeliveryWorkflow) Gate(_ context.Context, _ string, slug, _ string) (GateResult, error) {
	fake.recordEvent(slug, "gate")
	if result, ok := fake.gates[slug]; ok {
		return result, nil
	}
	return GateResult{Passed: true}, nil
}

func (fake *fakeDeliveryWorkflow) Authorization(_ context.Context, _ string, slug string) (Authorization, error) {
	fake.recordEvent(slug, "authorization")
	if authorization, ok := fake.authorizations[slug]; ok {
		return authorization, nil
	}
	return Authorization{Operations: []string{"push", "pull_request", "merge"}}, nil
}

func (fake *fakeDeliveryWorkflow) Publication(_ context.Context, _ string, slug string) (Publication, error) {
	fake.recordEvent(slug, "publication")
	if publication, ok := fake.publications[slug]; ok {
		return publication, nil
	}
	publication := Publication{Remote: "origin", HeadBranch: "deliver/" + slug, BaseBranch: "main", Title: slug, Body: "deliver " + slug}
	fake.publications[slug] = publication
	return publication, nil
}

func (fake *fakeDeliveryWorkflow) recordEvent(slug, event string) {
	fake.events[slug] = append(fake.events[slug], event)
}

type fakeDeliveryBoundary struct {
	remoteHeads        map[string]string
	pullRequests       map[string]PullRequest
	merged             map[string]string
	failedChecks       map[string]bool
	pushEffects        int
	pullRequestEffects int
	mergeEffects       int
	totalCalls         int
	callsBySlug        map[string]int
	recordEvent        func(string, string)
}

func newFakeDeliveryBoundary() *fakeDeliveryBoundary {
	return &fakeDeliveryBoundary{
		remoteHeads:  map[string]string{},
		pullRequests: map[string]PullRequest{},
		merged:       map[string]string{},
		failedChecks: map[string]bool{},
		callsBySlug:  map[string]int{},
	}
}

var _ PullRequestBoundary = (*fakeDeliveryBoundary)(nil)

func (fake *fakeDeliveryBoundary) RemoteBranchHead(_ context.Context, remote, branch string) (RemoteHead, bool, error) {
	fake.recordCall(branch, "observe-push")
	sha, found := fake.remoteHeads[branch]
	return RemoteHead{Remote: remote, Branch: branch, SHA: sha}, found, nil
}

func (fake *fakeDeliveryBoundary) PushBranch(_ context.Context, remote, branch string) (RemoteHead, error) {
	fake.recordCall(branch, "push")
	fake.pushEffects++
	slug := slugFromBranch(branch)
	sha := "archived-" + slug
	fake.remoteHeads[branch] = sha
	return RemoteHead{Remote: remote, Branch: branch, SHA: sha}, nil
}

func (fake *fakeDeliveryBoundary) FindOrCreatePullRequest(_ context.Context, req PullRequestRequest) (PullRequestResult, error) {
	fake.recordCall(req.HeadBranch, "pull-request")
	if pullRequest, found := fake.pullRequests[req.HeadBranch]; found {
		return PullRequestResult{PullRequest: pullRequest}, nil
	}
	fake.pullRequestEffects++
	number := fmt.Sprintf("%d", fake.pullRequestEffects)
	pullRequest := PullRequest{Number: number, State: "OPEN", HeadBranch: req.HeadBranch, HeadSHA: fake.remoteHeads[req.HeadBranch]}
	fake.pullRequests[req.HeadBranch] = pullRequest
	return PullRequestResult{PullRequest: pullRequest, Created: true}, nil
}

func (fake *fakeDeliveryBoundary) CurrentHeadChecks(_ context.Context, number string) (CheckReport, error) {
	pullRequest, ok := fake.pullRequestByNumber(number)
	if !ok {
		return CheckReport{}, errors.New("pull request not found")
	}
	fake.recordCall(pullRequest.HeadBranch, "checks")
	if fake.failedChecks[slugFromBranch(pullRequest.HeadBranch)] {
		return CheckReport{HeadSHA: pullRequest.HeadSHA, Checks: []PullRequestCheck{{Name: "test", Bucket: "fail"}}}, nil
	}
	return CheckReport{HeadSHA: pullRequest.HeadSHA, Checks: []PullRequestCheck{{Name: "test", Bucket: "pass"}}}, nil
}

func (fake *fakeDeliveryBoundary) MergePullRequest(_ context.Context, number, expectedHead string) (MergeResult, error) {
	pullRequest, ok := fake.pullRequestByNumber(number)
	if !ok {
		pullRequest = PullRequest{Number: number, HeadSHA: expectedHead, HeadBranch: "deliver/spec-merge"}
	}
	fake.recordCall(pullRequest.HeadBranch, "merge")
	if mergeCommit, found := fake.merged[number]; found {
		pullRequest.State = "MERGED"
		pullRequest.MergeCommit = mergeCommit
		return MergeResult{PullRequest: pullRequest, AlreadyMerged: true}, nil
	}
	fake.mergeEffects++
	pullRequest.State = "MERGED"
	pullRequest.MergeCommit = "merge-" + number
	fake.merged[number] = pullRequest.MergeCommit
	return MergeResult{PullRequest: pullRequest}, nil
}

func (fake *fakeDeliveryBoundary) pullRequestByNumber(number string) (PullRequest, bool) {
	for _, pullRequest := range fake.pullRequests {
		if pullRequest.Number == number {
			return pullRequest, true
		}
	}
	return PullRequest{}, false
}

func (fake *fakeDeliveryBoundary) recordCall(branch, event string) {
	fake.totalCalls++
	slug := slugFromBranch(branch)
	fake.callsBySlug[slug]++
	if fake.recordEvent != nil {
		fake.recordEvent(slug, event)
	}
}

func (fake *fakeDeliveryBoundary) callsFor(slug string) int {
	return fake.callsBySlug[slug]
}

func slugFromBranch(branch string) string {
	const prefix = "deliver/"
	if len(branch) >= len(prefix) && branch[:len(prefix)] == prefix {
		return branch[len(prefix):]
	}
	return branch
}
