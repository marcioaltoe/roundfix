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
	"strings"
	"testing"
	"time"

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
		Workspace:    workflow,
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
		"create-branch", "run", "policy", "review", "archive", "gate", "authorization",
		"publication", "push", "pull-request", "checks", "authorization", "merge",
	}
	if got := workflow.events["spec-one"]; !reflect.DeepEqual(got, wantReviewed) {
		t.Fatalf("reviewed item actions = %v, want %v", got, wantReviewed)
	}
	wantOmitted := []string{
		"create-branch", "run", "policy", "review-omitted", "archive", "gate", "authorization",
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
		item.Branch = "deliver/spec-push"
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
		item.Branch = "deliver/spec-pull-request"
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
		item.Branch = "deliver/spec-merge"
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
	slugs := []string{"run-blocked", "review-findings", "review-blocked", "gate-blocked", "checks-blocked", "checks-cancelled", "after-blockers"}
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
	boundary.cancelledChecks["checks-cancelled"] = true
	engine := newTestDeliveryEngine(runStore, workflow, boundary)
	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	queue := readDeliveryQueue(t, ctx, runStore, gitRoot)
	for index, blocker := range []string{BlockerRunUnresolved, BlockerReviewFindings, BlockerReviewBlocked, BlockerGateFailed, BlockerChecksFailed, BlockerChecksFailed} {
		if queue.Items[index].Stage != store.DeliveryStageParked || queue.Items[index].Blocker != blocker {
			t.Fatalf("item %q = %+v, want parked as %q", queue.Items[index].SpecSlug, queue.Items[index], blocker)
		}
	}
	if queue.Items[len(queue.Items)-1].Stage != store.DeliveryStageMerged {
		t.Fatalf("item after blockers stage = %q, want merged", queue.Items[len(queue.Items)-1].Stage)
	}
}

func TestEachItemPublishesFromItsOwnBranch(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-own-branches"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"parked-spec", "merged-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	boundary.failedChecks["parked-spec"] = true
	engine := newTestDeliveryEngine(runStore, workflow, boundary)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	queue := readDeliveryQueue(t, ctx, runStore, gitRoot)
	if queue.Items[0].Branch == "" || queue.Items[1].Branch == "" || queue.Items[0].Branch == queue.Items[1].Branch {
		t.Fatalf("item branches = %q and %q, want two recorded branches", queue.Items[0].Branch, queue.Items[1].Branch)
	}
	if got := workflow.runBranches["parked-spec"]; got != queue.Items[0].Branch {
		t.Fatalf("parked Spec ran on %q, want %q", got, queue.Items[0].Branch)
	}
	if got := workflow.runBranches["merged-spec"]; got != queue.Items[1].Branch {
		t.Fatalf("merged Spec ran on %q, want %q", got, queue.Items[1].Branch)
	}
	parkedPR := boundary.pullRequests[queue.Items[0].Branch]
	if parkedPR.Number == "" {
		t.Fatal("parked item did not create its own pull request")
	}
	if _, merged := boundary.merged[parkedPR.Number]; merged {
		t.Fatalf("parked item's pull request %s was merged", parkedPR.Number)
	}
	if queue.Items[1].Stage != store.DeliveryStageMerged {
		t.Fatalf("second item stage = %q, want merged", queue.Items[1].Stage)
	}
}

func TestAPullRequestOfAnotherItemIsNeverReused(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-pr-ownership"
	queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"first-spec", "second-spec"})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	first := queue.Items[0]
	first.Stage = store.DeliveryStageParked
	first.Branch = "roundfix/deliver-first-spec"
	first.PullRequestNumber = "17"
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, first); err != nil {
		t.Fatalf("seed first item: %v", err)
	}
	second := queue.Items[1]
	second.Stage = store.DeliveryStagePublishing
	second.Branch = "roundfix/deliver-second-spec"
	second.CandidateCommits = []string{"reviewed-second", "archived-second"}
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, second); err != nil {
		t.Fatalf("seed second item: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	boundary.findPullRequest = func(req PullRequestRequest) PullRequestResult {
		return PullRequestResult{PullRequest: PullRequest{
			Number:     "17",
			State:      "OPEN",
			HeadBranch: req.HeadBranch,
			HeadSHA:    "archived-second",
		}}
	}
	engine := newTestDeliveryEngine(runStore, workflow, boundary)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[1]
	if got.Stage != store.DeliveryStageParked || !strings.Contains(got.Blocker, "another Delivery Queue item") {
		t.Fatalf("second item = %+v, want parked for pull request ownership", got)
	}
	if boundary.mergeEffects != 0 {
		t.Fatalf("merge effects = %d, want none", boundary.mergeEffects)
	}
}

func TestPushRefusesTheBaseBranch(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-base-branch"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"base-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.publications["base-spec"] = Publication{Remote: "origin", HeadBranch: "main", BaseBranch: "main", Title: "base-spec"}
	boundary := newFakeDeliveryBoundary()
	engine := newTestDeliveryEngine(runStore, workflow, boundary)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Stage != store.DeliveryStageParked || !strings.Contains(item.Blocker, "base branch") {
		t.Fatalf("base-branch item = %+v, want parked refusal", item)
	}
	if boundary.totalCalls != 0 {
		t.Fatalf("base-branch item made %d external calls, want none", boundary.totalCalls)
	}
}

func TestDeliveryWaitsThroughUnreportedChecks(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-unreported-checks"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"unreported-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	boundary.checkReports["unreported-spec"] = []CheckReport{
		{Checks: []PullRequestCheck{}},
		{Checks: []PullRequestCheck{{Name: "verify", Bucket: "pass"}}},
	}
	clock := &fakeDeliveryClock{now: time.Unix(100, 0)}
	engine := newTestDeliveryEngineWithWait(runStore, workflow, boundary, clock, clock)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Stage != store.DeliveryStageMerged {
		t.Fatalf("unreported-check item = %+v, want merged after checks passed", item)
	}
	if clock.sleeps != 1 {
		t.Fatalf("check waits = %d, want 1", clock.sleeps)
	}
}

func TestCheckReadErrorsAreRetriedUntilTheDeadline(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-check-read-error"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"retry-check-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	boundary.checkErrors["retry-check-spec"] = []error{errors.New("temporary check read failure")}
	boundary.checkReports["retry-check-spec"] = []CheckReport{
		{Checks: []PullRequestCheck{{Name: "verify", Bucket: "pass"}}},
	}
	clock := &fakeDeliveryClock{now: time.Unix(100, 0)}
	engine := newTestDeliveryEngineWithWait(runStore, workflow, boundary, clock, clock)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Stage != store.DeliveryStageMerged {
		t.Fatalf("check-read-error item = %+v, want merged after retry", item)
	}
	if clock.sleeps != 1 {
		t.Fatalf("check waits = %d, want 1", clock.sleeps)
	}
}

func TestPersistentCheckReadErrorsParkAtTheDeadline(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-persistent-check-read-error"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"timeout-check-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	boundary.checkErrors["timeout-check-spec"] = []error{
		errors.New("temporary check read failure"),
		errors.New("temporary check read failure"),
		errors.New("temporary check read failure"),
	}
	clock := &fakeDeliveryClock{now: time.Unix(100, 0)}
	engine := newTestDeliveryEngineWithWait(runStore, workflow, boundary, clock, clock)
	engine.checkTimeout = 2 * time.Second

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Stage != store.DeliveryStageParked || item.Blocker != BlockerChecksTimeout {
		t.Fatalf("check-read-error item = %+v, want parked as checks-timeout", item)
	}
	if clock.sleeps != 2 {
		t.Fatalf("check waits = %d, want 2", clock.sleeps)
	}
}

func TestChecksTimeoutParksTheItem(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-check-timeout"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"timeout-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	boundary.checkReports["timeout-spec"] = []CheckReport{
		{Checks: []PullRequestCheck{}},
		{Checks: []PullRequestCheck{{Name: "verify", Bucket: "pending"}}},
		{Checks: []PullRequestCheck{{Name: "verify", Bucket: "pending"}}},
	}
	clock := &fakeDeliveryClock{now: time.Unix(100, 0)}
	engine := newTestDeliveryEngineWithWait(runStore, workflow, boundary, clock, clock)
	engine.checkTimeout = 2 * time.Second

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if item.Stage != store.DeliveryStageParked || item.Blocker != BlockerChecksTimeout {
		t.Fatalf("timed-out item = %+v, want parked as checks-timeout", item)
	}
}

func TestUnmatchedPushIntentRetriesWhenTheRemoteDiffers(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-retry-push"
	queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"retry-spec"})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStagePublishing
	item.Branch = "roundfix/deliver-retry-spec"
	item.CandidateCommits = []string{"reviewed-retry", "archived-retry"}
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("seed publishing item: %v", err)
	}
	if _, err := runStore.RecordDeliveryActionIntent(ctx, gitRoot, item.SpecSlug, store.DeliveryActionPush); err != nil {
		t.Fatalf("record unmatched push intent: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	boundary.remoteHeads[item.Branch] = "moved-remote-head"
	engine := newTestDeliveryEngine(runStore, workflow, boundary)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if got.Stage != store.DeliveryStageMerged || boundary.pushEffects != 1 {
		t.Fatalf("retried item = %+v pushes=%d, want merged after one retry", got, boundary.pushEffects)
	}
	assertNoUnmatchedDeliveryIntents(t, ctx, runStore, gitRoot)
}

func TestAnItemErrorParksAndTheQueueContinues(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-item-error"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"error-spec", "next-spec"}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.runErrors["error-spec"] = errors.New("persistent executor error")
	boundary := newFakeDeliveryBoundary()
	engine := newTestDeliveryEngine(runStore, workflow, boundary)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	queue := readDeliveryQueue(t, ctx, runStore, gitRoot)
	if queue.Items[0].Stage != store.DeliveryStageParked || !strings.Contains(queue.Items[0].Blocker, "persistent executor error") {
		t.Fatalf("error item = %+v, want parked with reason", queue.Items[0])
	}
	if queue.Items[1].Stage != store.DeliveryStageMerged {
		t.Fatalf("next item stage = %q, want merged", queue.Items[1].Stage)
	}
}

func TestAnAlreadyMergedPullRequestMustCarryTheReviewedHead(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-merged-other-head"
	queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"merged-spec"})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageMerging
	item.Branch = "roundfix/deliver-merged-spec"
	item.CandidateCommits = []string{"reviewed-merged", "archived-merged"}
	item.PullRequestNumber = "41"
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("seed merging item: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	boundary.pullRequests[item.Branch] = PullRequest{Number: "41", State: "MERGED", HeadBranch: item.Branch, HeadSHA: "foreign-head", MergeCommit: "merge-41"}
	boundary.merged["41"] = "merge-41"
	engine := newTestDeliveryEngine(runStore, workflow, boundary)

	if _, err := engine.Run(ctx, gitRoot); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if got.Stage != store.DeliveryStageParked || got.Blocker != BlockerReviewStale || got.MergeCommit != "" {
		t.Fatalf("merged-at-other-head item = %+v, want parked as review-stale", got)
	}
}

func newTestDeliveryEngine(runStore *store.Store, workflow *fakeDeliveryWorkflow, boundary *fakeDeliveryBoundary) *Engine {
	return newTestDeliveryEngineWithWait(runStore, workflow, boundary, nil, nil)
}

func newTestDeliveryEngineWithWait(
	runStore *store.Store,
	workflow *fakeDeliveryWorkflow,
	boundary *fakeDeliveryBoundary,
	clock Clock,
	sleeper Sleeper,
) *Engine {
	boundary.recordEvent = workflow.recordEvent
	return NewEngine(runStore, EngineDependencies{
		Workspace:     workflow,
		Runner:        workflow,
		Reviewer:      workflow,
		Archiver:      workflow,
		Gate:          workflow,
		Authorizer:    workflow,
		Publication:   workflow,
		PullRequests:  boundary,
		Clock:         clock,
		Sleeper:       sleeper,
		CheckTimeout:  10 * time.Minute,
		CheckInterval: time.Second,
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
	runErrors      map[string]error
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
	runBranches    map[string]string
	currentBranch  string
}

func newFakeDeliveryWorkflow() *fakeDeliveryWorkflow {
	return &fakeDeliveryWorkflow{
		runErrors:      map[string]error{},
		runs:           map[string]RunResult{},
		policies:       map[string]ReviewPolicy{},
		reviewResults:  map[string]ReviewResult{},
		archives:       map[string]ArchiveResult{},
		gates:          map[string]GateResult{},
		authorizations: map[string]Authorization{},
		publications:   map[string]Publication{},
		events:         map[string][]string{},
		runBranches:    map[string]string{},
	}
}

func (fake *fakeDeliveryWorkflow) CreateItemBranch(_ context.Context, _ string, slug string) (string, error) {
	branch := "roundfix/deliver-" + slug
	fake.currentBranch = branch
	fake.recordEvent(slug, "create-branch")
	return branch, nil
}

func (fake *fakeDeliveryWorkflow) UseItemBranch(_ context.Context, _ string, branch string) error {
	fake.currentBranch = branch
	return nil
}

func (fake *fakeDeliveryWorkflow) RunSpec(_ context.Context, _ string, slug string) (RunResult, error) {
	fake.recordEvent(slug, "run")
	fake.runBranches[slug] = fake.currentBranch
	if err, ok := fake.runErrors[slug]; ok {
		return RunResult{}, err
	}
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

func (fake *fakeDeliveryWorkflow) Publication(_ context.Context, _ string, slug, branch string) (Publication, error) {
	fake.recordEvent(slug, "publication")
	if publication, ok := fake.publications[slug]; ok {
		return publication, nil
	}
	publication := Publication{Remote: "origin", HeadBranch: branch, BaseBranch: "main", Title: slug, Body: "deliver " + slug}
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
	cancelledChecks    map[string]bool
	checkErrors        map[string][]error
	checkReports       map[string][]CheckReport
	findPullRequest    func(PullRequestRequest) PullRequestResult
	pushEffects        int
	pullRequestEffects int
	mergeEffects       int
	totalCalls         int
	callsBySlug        map[string]int
	recordEvent        func(string, string)
}

func newFakeDeliveryBoundary() *fakeDeliveryBoundary {
	return &fakeDeliveryBoundary{
		remoteHeads:     map[string]string{},
		pullRequests:    map[string]PullRequest{},
		merged:          map[string]string{},
		failedChecks:    map[string]bool{},
		cancelledChecks: map[string]bool{},
		checkErrors:     map[string][]error{},
		checkReports:    map[string][]CheckReport{},
		callsBySlug:     map[string]int{},
	}
}

var _ PullRequestBoundary = (*fakeDeliveryBoundary)(nil)

func (fake *fakeDeliveryBoundary) RemoteBranchHead(_ context.Context, remote, branch string) (RemoteHead, bool, error) {
	fake.recordCall(branch, "observe-push")
	sha, found := fake.remoteHeads[branch]
	return RemoteHead{Remote: remote, Branch: branch, SHA: sha}, found, nil
}

func (fake *fakeDeliveryBoundary) PushBranch(_ context.Context, remote, branch, head string) (RemoteHead, error) {
	fake.recordCall(branch, "push")
	fake.pushEffects++
	fake.remoteHeads[branch] = head
	return RemoteHead{Remote: remote, Branch: branch, SHA: head}, nil
}

func (fake *fakeDeliveryBoundary) FindOrCreatePullRequest(_ context.Context, req PullRequestRequest) (PullRequestResult, error) {
	fake.recordCall(req.HeadBranch, "pull-request")
	if fake.findPullRequest != nil {
		return fake.findPullRequest(req), nil
	}
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
	slug := slugFromBranch(pullRequest.HeadBranch)
	if errs := fake.checkErrors[slug]; len(errs) > 0 {
		err := errs[0]
		fake.checkErrors[slug] = errs[1:]
		return CheckReport{}, err
	}
	if reports := fake.checkReports[slug]; len(reports) > 0 {
		report := reports[0]
		fake.checkReports[slug] = reports[1:]
		report.HeadSHA = pullRequest.HeadSHA
		return report, nil
	}
	if fake.failedChecks[slug] {
		return CheckReport{HeadSHA: pullRequest.HeadSHA, Checks: []PullRequestCheck{{Name: "test", Bucket: "fail"}}}, nil
	}
	if fake.cancelledChecks[slug] {
		return CheckReport{HeadSHA: pullRequest.HeadSHA, Checks: []PullRequestCheck{{Name: "test", Bucket: "cancel"}}}, nil
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
	for _, prefix := range []string{"roundfix/deliver-", "deliver/"} {
		if len(branch) >= len(prefix) && branch[:len(prefix)] == prefix {
			return branch[len(prefix):]
		}
	}
	return branch
}

type fakeDeliveryClock struct {
	now    time.Time
	sleeps int
}

func (fake *fakeDeliveryClock) Now() time.Time {
	return fake.now
}

func (fake *fakeDeliveryClock) Sleep(_ context.Context, duration time.Duration) error {
	fake.sleeps++
	fake.now = fake.now.Add(duration)
	return nil
}
