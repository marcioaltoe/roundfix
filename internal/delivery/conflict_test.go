// Boundary: persisted delivery stages and deterministic PR polling. No remote IO.
package delivery

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/store"
)

type fakeConflictResolver struct {
	result ConflictResolution
	calls  int
}

func (fake *fakeConflictResolver) ResolveConflict(_ context.Context, workDir, slug, head string) (ConflictResolution, error) {
	if workDir != "/worktrees/example" || slug != "example" || head != "candidate" {
		return ConflictResolution{}, fmt.Errorf("unexpected conflict target: %s %s %s", workDir, slug, head)
	}
	fake.calls++
	return fake.result, nil
}

func exerciseConflict(t *testing.T, resolver ConflictResolver, reports []CheckReport) (store.DeliveryQueueItem, int) {
	t.Helper()
	ctx := t.Context()
	runStore := openDeliveryEngineStore(t, ctx)
	queue, err := runStore.CreateDeliveryQueue(ctx, "/repo", []string{"example"})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageChecking
	item.Branch, item.Worktree, item.PullRequestNumber = "roundfix/deliver-example", "/worktrees/example", "1"
	item.CandidateCommits = []string{"candidate"}
	recordDeliveryItemWorkspace(t, ctx, runStore, "/repo", &item)
	if err := runStore.UpdateDeliveryQueueItem(ctx, "/repo", item); err != nil {
		t.Fatal(err)
	}
	polls := 0
	boundary := &fakePullRequestBoundary{currentHeadChecks: func(context.Context, string) (CheckReport, error) {
		if polls >= len(reports) {
			t.Fatal("unexpected poll")
		}
		report := reports[polls]
		polls++
		return report, nil
	}}
	clock := &fakeDeliveryClock{now: time.Now()}
	engine := NewEngine(runStore, EngineDependencies{PullRequests: boundary, Conflicts: resolver, Clock: clock, Sleeper: clock, CheckTimeout: time.Second, CheckInterval: time.Second})
	if err := engine.checkCandidate(ctx, "/repo", &item); err != nil {
		t.Fatal(err)
	}
	return readDeliveryQueue(t, ctx, runStore, "/repo").Items[0], polls
}
func conflictReport(state string) CheckReport {
	return CheckReport{HeadSHA: "candidate", Mergeable: state, Checks: []PullRequestCheck{{Bucket: "pass"}}}
}
func TestAConflictingPullRequestStopsTheCheckWaitAtOnce(t *testing.T) {
	item, polls := exerciseConflict(t, nil, []CheckReport{conflictReport("CONFLICTING")})
	if polls != 1 || item.Blocker != BlockerPullRequestConflict || item.Stage != store.DeliveryStageParked {
		t.Fatalf("item=%+v polls=%d", item, polls)
	}
}
func TestAnUnknownMergeableStateKeepsWaiting(t *testing.T) {
	item, polls := exerciseConflict(t, nil, []CheckReport{conflictReport("UNKNOWN"), conflictReport("MERGEABLE")})
	if polls != 2 || item.Stage != store.DeliveryStageMerging {
		t.Fatalf("item=%+v polls=%d", item, polls)
	}
}
func TestADerivedConflictReturnsTheItemToTheGate(t *testing.T) {
	resolver := &fakeConflictResolver{result: ConflictResolution{Head: "derived-merge", Regenerated: []string{"make generate"}}}
	item, polls := exerciseConflict(t, resolver, []CheckReport{conflictReport("CONFLICTING")})
	if polls != 1 || resolver.calls != 1 || item.Stage != store.DeliveryStageGating || !reflect.DeepEqual(item.CandidateCommits, []string{"candidate", "derived-merge"}) {
		t.Fatalf("item=%+v calls=%d", item, resolver.calls)
	}
}
func TestASourceConflictParksWithItsPaths(t *testing.T) {
	resolver := &fakeConflictResolver{result: ConflictResolution{SourcePaths: []string{"source.go"}}}
	item, _ := exerciseConflict(t, resolver, []CheckReport{conflictReport("CONFLICTING")})
	park := ClassifyPark(store.DeliveryQueue{}, item)
	if item.Blocker != BlockerPullRequestConflict+": source.go" || park.Class != ParkClassConflict || !strings.Contains(park.Next, "resolve source.go, commit, then run roundfix deliver retry example") {
		t.Fatalf("item=%+v park=%+v", item, park)
	}
}
func TestALaggingPullRequestHeadIsReadAgain(t *testing.T) {
	ctx := t.Context()
	runStore := openDeliveryEngineStore(t, ctx)
	queue, err := runStore.CreateDeliveryQueue(ctx, "/repo", []string{"example"})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.CandidateCommits = []string{"old", "new"}
	if err := runStore.UpdateDeliveryQueueItem(ctx, "/repo", item); err != nil {
		t.Fatal(err)
	}
	polls := 0
	boundary := &fakePullRequestBoundary{findOrCreatePullRequest: func(context.Context, PullRequestRequest) (PullRequestResult, error) {
		polls++
		head := "old"
		if polls == 2 {
			head = "new"
		}
		return PullRequestResult{PullRequest: PullRequest{Number: "1", HeadBranch: "item", HeadSHA: head}}, nil
	}}
	clock := &fakeDeliveryClock{now: time.Now()}
	engine := NewEngine(runStore, EngineDependencies{PullRequests: boundary, Clock: clock, Sleeper: clock, CheckTimeout: 2 * time.Second, CheckInterval: time.Second})
	pr, err := engine.createPullRequest(ctx, "/repo", "/worktrees/example", "example", Publication{HeadBranch: "item"}, "new")
	if err != nil || polls != 2 || pr.HeadSHA != "new" {
		t.Fatalf("pr=%+v polls=%d err=%v", pr, polls, err)
	}
}
func TestRetryResumesAResolvedConflictAtReview(t *testing.T) {
	ctx := t.Context()
	runStore := openDeliveryEngineStore(t, ctx)
	item := seedParkedRetryItem(t, ctx, runStore, "/repo", "example", BlockerPullRequestConflict+": source.go")
	workflow := newFakeDeliveryWorkflow()
	recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "resolved"}}}
	history := &operatorArchiveHistory{descends: true}
	engine := newRetryDeliveryEngine(runStore, workflow, nil, recovery, newFakeDeliveryBoundary())
	engine.history = history
	result, err := engine.Retry(ctx, "/repo", "example")
	if err != nil || result.Stage != store.DeliveryStageReviewing {
		t.Fatalf("result=%+v err=%v item=%+v", result, err, item)
	}
	retried := readDeliveryQueue(t, ctx, runStore, "/repo").Items[0]
	if retried.CandidateCommits[len(retried.CandidateCommits)-1] != "resolved" {
		t.Fatalf("item=%+v", retried)
	}
}
func TestGitHubCLIReadsTheMergeableState(t *testing.T) {
	script := newScriptedCommandRunner(t, commandStep{name: "gh", args: []string{"pr", "view", "1", "--json", pullRequestJSONFields}, result: CommandResult{Stdout: `{"number":1,"headRefOid":"candidate","mergeable":"CONFLICTING"}`}})
	report, err := (GitHubCLI{WorkDir: "/repo", Runner: script}).CurrentHeadChecks(t.Context(), "1")
	if err != nil || report.Mergeable != "CONFLICTING" || report.HeadSHA != "candidate" {
		t.Fatalf("report=%+v err=%v", report, err)
	}
}
