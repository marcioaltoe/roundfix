// Suite: bounded per-call Delivery Queue check waits.
// Boundary IN: real queue persistence and fake-clock engine polling.
// Boundary OUT: forge operations use existing fakes; no real-time waits.
package delivery

import (
	"context"
	"testing"
	"time"

	"roundfix/internal/config"
	"roundfix/internal/store"
)

func newCheckWaitFixture(t *testing.T, timeout time.Duration) (*Engine, *store.Store, *store.DeliveryQueueItem, *fakeDeliveryClock, *fakePullRequestBoundary) {
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
	item.CandidateCommits = []string{"earlier", "candidate"}
	recordDeliveryItemWorkspace(t, ctx, runStore, "/repo", &item)
	if err := runStore.UpdateDeliveryQueueItem(ctx, "/repo", item); err != nil {
		t.Fatal(err)
	}
	clock := &fakeDeliveryClock{now: time.Unix(100, 0)}
	boundary := &fakePullRequestBoundary{}
	engine := NewEngine(runStore, EngineDependencies{PullRequests: boundary, Clock: clock, Sleeper: clock, CheckTimeout: timeout})
	return engine, runStore, &item, clock, boundary
}

func TestTheCheckWaitDefaultsToSixtyMinutes(t *testing.T) {
	t.Parallel()
	engine, runStore, item, clock, boundary := newCheckWaitFixture(t, 0)
	polls := 0
	boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) {
		polls++
		return recoveryReport("pending"), nil
	}
	start := clock.Now()
	if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
		t.Fatal(err)
	}
	stored := readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
	elapsed := clock.Now().Sub(start)
	if elapsed != config.DefaultDeliveryCheckTimeout || polls != 241 || stored.Stage != store.DeliveryStageParked || stored.Blocker != "checks-timeout: current-head checks not settled within delivery.check_timeout 1h0m0s" || ClassifyPark(store.DeliveryQueue{}, stored).Class != ParkClassEnvironment {
		t.Fatalf("elapsed=%s polls=%d item=%+v", elapsed, polls, stored)
	}
}

func TestTheCheckWaitUsesTheConfiguredDuration(t *testing.T) {
	t.Parallel()
	engine, runStore, item, clock, boundary := newCheckWaitFixture(t, 45*time.Minute)
	boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) { return recoveryReport("pending"), nil }
	start := clock.Now()
	if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
		t.Fatal(err)
	}
	stored := readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
	if clock.Now().Sub(start) != 45*time.Minute || stored.Stage != store.DeliveryStageParked || stored.Blocker != "checks-timeout: current-head checks not settled within delivery.check_timeout 45m0s" {
		t.Fatalf("elapsed=%s item=%+v", clock.Now().Sub(start), stored)
	}
}

func TestEachCheckWaitStartsItsOwnDeadline(t *testing.T) {
	t.Parallel()
	engine, runStore, item, clock, boundary := newCheckWaitFixture(t, 45*time.Minute)
	polls := 0
	boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) {
		polls++
		if polls == 161 {
			return recoveryReport("pass"), nil
		}
		return recoveryReport("pending"), nil
	}
	start := clock.Now()
	if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
		t.Fatal(err)
	}
	stored := readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
	if clock.Now().Sub(start) != 40*time.Minute || polls != 161 || stored.Stage != store.DeliveryStageMerging {
		t.Fatalf("elapsed=%s polls=%d item=%+v", clock.Now().Sub(start), polls, stored)
	}
	item.Stage = store.DeliveryStageChecking
	if err := runStore.UpdateDeliveryQueueItem(t.Context(), "/repo", *item); err != nil {
		t.Fatal(err)
	}
	boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) { return recoveryReport("pending"), nil }
	start = clock.Now()
	if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
		t.Fatal(err)
	}
	stored = readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
	if clock.Now().Sub(start) != 45*time.Minute || stored.Stage != store.DeliveryStageParked || stored.Blocker != "checks-timeout: current-head checks not settled within delivery.check_timeout 45m0s" {
		t.Fatalf("elapsed=%s item=%+v", clock.Now().Sub(start), stored)
	}
}

func TestThePullRequestHeadWaitUsesTheCheckWait(t *testing.T) {
	t.Parallel()
	engine, _, item, clock, boundary := newCheckWaitFixture(t, 45*time.Minute)
	boundary.findOrCreatePullRequest = func(context.Context, PullRequestRequest) (PullRequestResult, error) {
		return PullRequestResult{PullRequest: PullRequest{Number: "1", HeadBranch: item.Branch, HeadSHA: "earlier"}}, nil
	}
	start := clock.Now()
	_, err := engine.createPullRequest(t.Context(), "/repo", item.Worktree, item.SpecSlug, Publication{HeadBranch: item.Branch, BaseBranch: "main"}, "candidate")
	const want = `find or create pull request: PR Head Branch is at "earlier", expected "candidate"`
	if err == nil || err.Error() != want || clock.Now().Sub(start) != 45*time.Minute {
		t.Fatalf("elapsed=%s error=%v", clock.Now().Sub(start), err)
	}
}
