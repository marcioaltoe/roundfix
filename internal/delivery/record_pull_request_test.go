package delivery

// Suite: Pull Request Record Commit at Checking.
// Invariant: only an exact recorded head reaches checks and merge.
// Boundary IN: durable queue, engine and fake recorder.
// Boundary OUT: forge and repository operations use existing boundary fakes.

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/store"
)

type testPullRequestRecorder struct {
	record func(context.Context, string, string, string, string) (PullRequestRecord, error)
}

func (r testPullRequestRecorder) RecordPullRequest(ctx context.Context, dir, slug, head, number string) (PullRequestRecord, error) {
	return r.record(ctx, dir, slug, head, number)
}

func TestDeliveryRecordsThePullRequestBeforeChecks(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	runStore := openDeliveryEngineStore(t, ctx)
	const repo = "/repo-record"
	if _, err := runStore.CreateDeliveryQueue(ctx, repo, []string{"spec-record"}); err != nil {
		t.Fatal(err)
	}
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	engine := newTestDeliveryEngine(runStore, workflow, boundary)
	var log bytes.Buffer
	engine.log = &log
	engine.recorder = testPullRequestRecorder{record: func(_ context.Context, dir, slug, head, number string) (PullRequestRecord, error) {
		if dir != "/worktrees/spec-record" || slug != "spec-record" || number != "1" {
			t.Fatalf("record args: %s %s %s %s", dir, slug, head, number)
		}
		workflow.recordEvent(slug, "record")
		if head == "recorded-spec-record" {
			return PullRequestRecord{}, nil
		}
		if head != "archived-spec-record" {
			t.Fatalf("candidate=%s", head)
		}
		return PullRequestRecord{Head: "recorded-spec-record", Exact: true}, nil
	}}
	// A forge updates the open pull request's head when Publishing pushes again.
	boundary.findPullRequest = func(req PullRequestRequest) PullRequestResult {
		pr, found := boundary.pullRequests[req.HeadBranch]
		if !found {
			boundary.pullRequestEffects++
			pr = PullRequest{Number: "1", State: "OPEN", HeadBranch: req.HeadBranch}
		}
		pr.HeadSHA = boundary.remoteHeads[req.HeadBranch]
		boundary.pullRequests[req.HeadBranch] = pr
		return PullRequestResult{PullRequest: pr, Created: !found}
	}
	if _, err := engine.Run(ctx, repo); err != nil {
		t.Fatal(err)
	}
	item := readDeliveryQueue(t, ctx, runStore, repo).Items[0]
	wantEvents := []string{"create-branch", "run", "policy", "review", "archive", "gate", "authorization", "publication", "push", "pull-request", "record", "authorization", "publication", "push", "pull-request", "record", "checks", "authorization", "merge"}
	if !reflect.DeepEqual(workflow.events[item.SpecSlug], wantEvents) {
		t.Fatalf("actions=%v want=%v", workflow.events[item.SpecSlug], wantEvents)
	}
	if item.Stage != store.DeliveryStageMerged || !reflect.DeepEqual(item.CandidateCommits, []string{"reviewed-spec-record", "archived-spec-record", "recorded-spec-record"}) {
		t.Fatalf("item=%+v", item)
	}
	if boundary.pullRequestEffects != 1 || boundary.pushEffects != 2 || boundary.mergeEffects != 1 {
		t.Fatalf("effects=%+v", boundary)
	}
	if !strings.Contains(log.String(), "roundfix: Delivery Queue item spec-record: recorded pull request 1 in its Archive Record\n") {
		t.Fatalf("log=%q", log.String())
	}
}

func seedRecordChecking(t *testing.T) (*Engine, *store.Store, *fakeDeliveryWorkflow, *fakeDeliveryBoundary, string) {
	t.Helper()
	ctx := t.Context()
	const repo = "/repo-record-retry"
	runStore := openDeliveryEngineStore(t, ctx)
	queue, err := runStore.CreateDeliveryQueue(ctx, repo, []string{"spec-retry"})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageChecking
	item.Branch = "roundfix/deliver-spec-retry"
	item.Worktree = "/worktrees/spec-retry"
	item.CandidateCommits = []string{"reviewed", "archived"}
	item.PullRequestNumber = "7"
	recordDeliveryItemWorkspace(t, ctx, runStore, repo, &item)
	if err := runStore.UpdateDeliveryQueueItem(ctx, repo, item); err != nil {
		t.Fatal(err)
	}
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	boundary.remoteHeads[item.Branch] = "archived"
	boundary.pullRequests[item.Branch] = PullRequest{Number: "7", State: "OPEN", HeadBranch: item.Branch, HeadSHA: "archived"}
	boundary.findPullRequest = func(req PullRequestRequest) PullRequestResult {
		pr := boundary.pullRequests[req.HeadBranch]
		pr.HeadSHA = boundary.remoteHeads[req.HeadBranch]
		boundary.pullRequests[req.HeadBranch] = pr
		return PullRequestResult{PullRequest: pr}
	}
	return newTestDeliveryEngine(runStore, workflow, boundary), runStore, workflow, boundary, repo
}

func TestDeliveryRetryAdoptsAPullRequestRecordCommit(t *testing.T) {
	t.Parallel()
	engine, runStore, workflow, boundary, repo := seedRecordChecking(t)
	adopted := 0
	engine.recorder = testPullRequestRecorder{record: func(_ context.Context, _, _, head, number string) (PullRequestRecord, error) {
		if number != "7" {
			t.Fatalf("number=%s", number)
		}
		if head == "recorded" {
			return PullRequestRecord{}, nil
		}
		adopted++
		return PullRequestRecord{Head: "recorded", Exact: true}, nil
	}}
	if _, err := engine.Run(t.Context(), repo); err != nil {
		t.Fatal(err)
	}
	if _, err := engine.Run(t.Context(), repo); err != nil {
		t.Fatal(err)
	}
	item := readDeliveryQueue(t, t.Context(), runStore, repo).Items[0]
	if item.Stage != store.DeliveryStageMerged || adopted != 1 || !reflect.DeepEqual(item.CandidateCommits, []string{"reviewed", "archived", "recorded"}) {
		t.Fatalf("item=%+v adopted=%d", item, adopted)
	}
	if boundary.pushEffects != 1 || boundary.pullRequestEffects != 0 || boundary.mergeEffects != 1 {
		t.Fatalf("effects=%+v", boundary)
	}
	for _, event := range workflow.events[item.SpecSlug] {
		if event == "gate" || event == "run" || event == "archive" {
			t.Fatalf("replayed %s", event)
		}
	}
}

func TestDeliveryParksAnInexactPullRequestRecordCommit(t *testing.T) {
	t.Parallel()
	for _, result := range []PullRequestRecord{{Head: "inexact"}, {Head: "archived", Exact: true}} {
		t.Run(result.Head, func(t *testing.T) {
			engine, runStore, workflow, boundary, repo := seedRecordChecking(t)
			engine.recorder = testPullRequestRecorder{record: func(context.Context, string, string, string, string) (PullRequestRecord, error) { return result, nil }}
			if _, err := engine.Run(t.Context(), repo); err != nil {
				t.Fatal(err)
			}
			item := readDeliveryQueue(t, t.Context(), runStore, repo).Items[0]
			if item.Stage != store.DeliveryStageParked || item.Blocker != BlockerReviewStale || len(item.CandidateCommits) != 2 {
				t.Fatalf("item=%+v", item)
			}
			if boundary.mergeEffects != 0 || len(workflow.events[item.SpecSlug]) != 0 {
				t.Fatalf("unsafe actions=%v merges=%d", workflow.events[item.SpecSlug], boundary.mergeEffects)
			}
		})
	}
}

func TestDeliveryParksWhenThePullRequestCannotBeRecorded(t *testing.T) {
	t.Parallel()
	engine, runStore, workflow, boundary, repo := seedRecordChecking(t)
	engine.recorder = testPullRequestRecorder{record: func(context.Context, string, string, string, string) (PullRequestRecord, error) {
		return PullRequestRecord{}, errors.New("record write refused")
	}}
	if _, err := engine.Run(t.Context(), repo); err != nil {
		t.Fatal(err)
	}
	item := readDeliveryQueue(t, t.Context(), runStore, repo).Items[0]
	if item.Stage != store.DeliveryStageParked || !strings.HasPrefix(item.Blocker, "delivery-error:") || !strings.Contains(item.Blocker, "record write refused") {
		t.Fatalf("item=%+v", item)
	}
	if boundary.mergeEffects != 0 || len(workflow.events[item.SpecSlug]) != 0 {
		t.Fatalf("unsafe actions=%v merges=%d", workflow.events[item.SpecSlug], boundary.mergeEffects)
	}
}
