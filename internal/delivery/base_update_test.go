// Suite: bounded Base Updates through persisted delivery stages.
// Boundary IN: Engine checking/Run, queue persistence and a deterministic clock.
// Boundary OUT: forge and default-branch merge are interface fakes.
package delivery

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/store"
)

type baseUpdateResolver struct {
	results []ConflictResolution
	err     error
	calls   []struct{ workDir, slug, head string }
}

func (fake *baseUpdateResolver) ResolveConflict(_ context.Context, workDir, slug, head string) (ConflictResolution, error) {
	fake.calls = append(fake.calls, struct{ workDir, slug, head string }{workDir, slug, head})
	if fake.err != nil {
		return ConflictResolution{}, fake.err
	}
	if len(fake.results) == 0 {
		return ConflictResolution{}, nil
	}
	return fake.results[min(len(fake.calls)-1, len(fake.results)-1)], nil
}

type baseUpdateFixture struct {
	engine   *Engine
	runStore *store.Store
	boundary *fakePullRequestBoundary
	workflow *fakeDeliveryWorkflow
	resolver *baseUpdateResolver
	log      bytes.Buffer
	polls    int
}

func newBaseUpdateFixture(t *testing.T, states ...string) *baseUpdateFixture {
	t.Helper()
	ctx := t.Context()
	// Share the merge-state fixture's persisted item, workflow and fake-clock
	// Check Wait instead of declaring a second timing budget for these polls.
	shared := newMergeStateFixture(t, states, 0)
	workflow, ok := shared.engine.gate.(*fakeDeliveryWorkflow)
	if !ok {
		t.Fatalf("merge-state fixture gate = %T, want *fakeDeliveryWorkflow", shared.engine.gate)
	}
	f := &baseUpdateFixture{engine: shared.engine, runStore: shared.runStore, workflow: workflow, resolver: &baseUpdateResolver{}}
	item := f.item(t)
	item.CandidateCommits = []string{"candidate"}
	if err := f.runStore.UpdateDeliveryQueueItem(ctx, "/repo", item); err != nil {
		t.Fatal(err)
	}
	f.boundary = &fakePullRequestBoundary{
		remoteBranchHead: func(context.Context, string, string) (RemoteHead, bool, error) { return RemoteHead{}, false, nil },
		pushBranch:       func(_ context.Context, _, branch, head string) (RemoteHead, error) { return RemoteHead{SHA: head}, nil },
		findOrCreatePullRequest: func(_ context.Context, req PullRequestRequest) (PullRequestResult, error) {
			stored := f.item(t)
			return PullRequestResult{PullRequest: PullRequest{Number: "1", HeadSHA: stored.CandidateCommits[len(stored.CandidateCommits)-1], HeadBranch: req.HeadBranch}}, nil
		},
		currentHeadChecks: func(_ context.Context, number string) (CheckReport, error) {
			stored := f.item(t)
			if number != "1" || stored.Stage != store.DeliveryStageChecking {
				t.Fatalf("poll number=%s item=%+v", number, stored)
			}
			state := states[min(f.polls, len(states)-1)]
			f.polls++
			return CheckReport{HeadSHA: stored.CandidateCommits[len(stored.CandidateCommits)-1], MergeState: state, Checks: []PullRequestCheck{{Name: "verify", Bucket: "pass"}}}, nil
		}}
	f.engine.pullRequests = f.boundary
	f.engine.conflicts = f.resolver
	// These tests resume checking directly; their reports follow the persisted
	// candidate rather than the shared fixture's fixed recovery head.
	f.engine.recovery = nil
	f.engine.log = &f.log
	return f
}

func (f *baseUpdateFixture) item(t *testing.T) store.DeliveryQueueItem {
	t.Helper()
	return readDeliveryQueue(t, t.Context(), f.runStore, "/repo").Items[0]
}
func (f *baseUpdateFixture) check(t *testing.T) store.DeliveryQueueItem {
	t.Helper()
	item := f.item(t)
	if err := f.engine.checkCandidate(t.Context(), "/repo", &item); err != nil {
		t.Fatal(err)
	}
	return f.item(t)
}
func (f *baseUpdateFixture) resumeChecking(t *testing.T) {
	t.Helper()
	item := f.item(t)
	item.Stage, item.Blocker = store.DeliveryStageChecking, ""
	if err := f.runStore.UpdateDeliveryQueueItem(t.Context(), "/repo", item); err != nil {
		t.Fatal(err)
	}
}

func TestABehindPullRequestReturnsToTheGateWithTheMergedHead(t *testing.T) {
	t.Parallel()
	f := newBaseUpdateFixture(t, "BEHIND")
	f.resolver.results = []ConflictResolution{{Head: "updated"}}
	item := f.check(t)
	if item.Stage != store.DeliveryStageGating || !reflect.DeepEqual(item.CandidateCommits, []string{"candidate", "updated"}) || f.polls != 1 || len(f.resolver.calls) != 1 {
		t.Fatalf("item=%+v calls=%v polls=%d", item, f.resolver.calls, f.polls)
	}
	call := f.resolver.calls[0]
	if call.workDir != item.Worktree || call.slug != item.SpecSlug || call.head != "candidate" {
		t.Fatalf("resolver target=%+v", call)
	}
	want := "roundfix: base update: Delivery Queue item example: pull request #1 was behind its base branch; merged the default branch at updated\n"
	if !strings.Contains(f.log.String(), want) {
		t.Fatalf("log=%q, want %q", f.log.String(), want)
	}
}

func TestABehindPullRequestIsMergedOnlyAtItsUpdatedHead(t *testing.T) {
	t.Parallel()
	f := newBaseUpdateFixture(t, "BEHIND", "CLEAN")
	f.resolver.results = []ConflictResolution{{Head: "updated"}}
	f.boundary.currentHeadChecks = func(_ context.Context, number string) (CheckReport, error) {
		stored := f.item(t)
		if number != "1" || stored.Stage != store.DeliveryStageChecking {
			t.Fatalf("poll number=%s item=%+v", number, stored)
		}
		head := stored.CandidateCommits[len(stored.CandidateCommits)-1]
		state := "BEHIND"
		if head == "updated" {
			state = "CLEAN"
		}
		f.polls++
		return CheckReport{HeadSHA: head, MergeState: state, Checks: []PullRequestCheck{{Name: "verify", Bucket: "pass"}}}, nil
	}
	var merges []string
	f.boundary.mergePullRequest = func(_ context.Context, number, head string) (MergeResult, error) {
		merges = append(merges, head)
		if head == "candidate" {
			return MergeResult{}, MergePolicyRefusalError{Stderr: "the base branch policy prohibits the merge"}
		}
		return MergeResult{PullRequest: PullRequest{Number: number, HeadSHA: head, State: "MERGED", MergeCommit: "squash"}}, nil
	}
	if _, err := f.engine.Run(t.Context(), "/repo"); err != nil {
		t.Fatal(err)
	}
	item := f.item(t)
	if item.Stage != store.DeliveryStageMerged || !reflect.DeepEqual(merges, []string{"updated"}) || len(f.resolver.calls) != 1 {
		t.Fatalf("item=%+v merges=%v calls=%v", item, merges, f.resolver.calls)
	}
	if !reflect.DeepEqual(f.workflow.events["example"], []string{"gate", "authorization", "publication", "authorization"}) {
		t.Fatalf("workflow events=%v", f.workflow.events)
	}
	assertNoUnmatchedDeliveryIntents(t, t.Context(), f.runStore, "/repo")
}

func TestABaseUpdateWithASourceConflictParksWithItsPaths(t *testing.T) {
	t.Parallel()
	f := newBaseUpdateFixture(t, "BEHIND")
	f.resolver.results = []ConflictResolution{{SourcePaths: []string{"source.go"}}}
	item := f.check(t)
	if item.Stage != store.DeliveryStageParked || item.Blocker != "pull-request-conflict: source.go" || ClassifyPark(store.DeliveryQueue{}, item).Class != ParkClassConflict || len(f.resolver.calls) != 1 {
		t.Fatalf("item=%+v calls=%v", item, f.resolver.calls)
	}
}

func TestABaseUpdateWithNothingToMergeKeepsWaiting(t *testing.T) {
	t.Parallel()
	t.Run("merges_once_clean", func(t *testing.T) {
		t.Parallel()
		for _, head := range []string{"", "candidate"} {
			t.Run(fmt.Sprintf("head=%q", head), func(t *testing.T) {
				t.Parallel()
				f := newBaseUpdateFixture(t, "BEHIND", "CLEAN")
				f.resolver.results = []ConflictResolution{{Head: head}}
				item := f.check(t)
				if item.Stage != store.DeliveryStageMerging || f.polls != 2 || len(f.resolver.calls) != 1 || len(item.CandidateCommits) != 1 {
					t.Fatalf("item=%+v polls=%d calls=%v", item, f.polls, f.resolver.calls)
				}
			})
		}
	})
	t.Run("times_out_with_the_merge_state", func(t *testing.T) {
		t.Parallel()
		f := newBaseUpdateFixture(t, "BEHIND")
		item := f.check(t)
		if item.Stage != store.DeliveryStageParked || item.Blocker != "checks-timeout: current-head checks not settled within delivery.check_timeout 3s; merge state BEHIND" || f.polls != 4 || len(item.CandidateCommits) != 1 {
			t.Fatalf("item=%+v polls=%d", item, f.polls)
		}
	})
}

func exhaustBaseUpdates(t *testing.T, f *baseUpdateFixture) store.DeliveryQueueItem {
	t.Helper()
	f.resolver.results = []ConflictResolution{{Head: "updated-1"}, {Head: "updated-2"}, {Head: "updated-3"}, {Head: "updated-4"}}
	for turn := 0; turn < BaseUpdateLimit; turn++ {
		item := f.check(t)
		if item.Stage != store.DeliveryStageGating {
			t.Fatalf("turn=%d item=%+v", turn, item)
		}
		f.resumeChecking(t)
	}
	item := f.check(t)
	park := ClassifyPark(store.DeliveryQueue{}, item)
	wantNext := "run roundfix deliver retry example once the base branch stops moving; the retry resumes at checking, which merges the default branch into the item branch again"
	if item.Stage != store.DeliveryStageParked || item.Blocker != "pull-request-behind: pull request #1 is still behind its base branch after 3 Base Updates" || len(f.resolver.calls) != BaseUpdateLimit || len(item.CandidateCommits) != BaseUpdateLimit+1 || park.Class != ParkClassEnvironment || park.Next != wantNext {
		t.Fatalf("item=%+v park=%+v calls=%v", item, park, f.resolver.calls)
	}
	return item
}

func TestAFourthBehindPullRequestParksPullRequestBehind(t *testing.T) {
	t.Parallel()
	exhaustBaseUpdates(t, newBaseUpdateFixture(t, "BEHIND"))
}
func TestTheBaseUpdateCountStartsOverAfterAPark(t *testing.T) {
	t.Parallel()
	f := newBaseUpdateFixture(t, "BEHIND")
	exhaustBaseUpdates(t, f)
	f.resumeChecking(t)
	item := f.check(t)
	if item.Stage != store.DeliveryStageGating || len(f.resolver.calls) != 4 || item.CandidateCommits[len(item.CandidateCommits)-1] != "updated-4" {
		t.Fatalf("item=%+v calls=%v", item, f.resolver.calls)
	}
}
func TestAnOwnerWithoutABaseUpdateParksPullRequestBehind(t *testing.T) {
	t.Parallel()
	f := newBaseUpdateFixture(t, "BEHIND")
	f.engine.conflicts = nil
	item := f.check(t)
	if item.Stage != store.DeliveryStageParked || item.Blocker != "pull-request-behind: pull request #1 is behind its base branch" || len(f.resolver.calls) != 0 {
		t.Fatalf("item=%+v calls=%v", item, f.resolver.calls)
	}
}
func TestAPolicyRefusalThatComesBackBehindGetsABaseUpdate(t *testing.T) {
	t.Parallel()
	f := newBaseUpdateFixture(t, "CLEAN", "BEHIND")
	f.resolver.results = []ConflictResolution{{Head: "updated"}}
	attempts := 0
	f.boundary.mergePullRequest = func(_ context.Context, _, head string) (MergeResult, error) {
		attempts++
		if head != "candidate" {
			t.Fatalf("merge head=%s", head)
		}
		return MergeResult{}, MergePolicyRefusalError{Stderr: "the base branch policy prohibits the merge"}
	}
	if item := f.check(t); item.Stage != store.DeliveryStageMerging {
		t.Fatalf("item=%+v", item)
	}
	item := f.item(t)
	if err := f.engine.mergeCandidate(t.Context(), "/repo", &item); err != nil {
		t.Fatal(err)
	}
	if item := f.item(t); item.Stage != store.DeliveryStageChecking {
		t.Fatalf("item=%+v", item)
	}
	item = f.check(t)
	if item.Stage != store.DeliveryStageGating || attempts != 1 || len(f.resolver.calls) != 1 || item.CandidateCommits[len(item.CandidateCommits)-1] != "updated" {
		t.Fatalf("item=%+v attempts=%d calls=%v", item, attempts, f.resolver.calls)
	}
}
func TestAConflictingAndBehindPullRequestTakesTheConflictPath(t *testing.T) {
	t.Parallel()
	f := newBaseUpdateFixture(t, "BEHIND")
	f.boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) {
		return CheckReport{HeadSHA: "candidate", Mergeable: "CONFLICTING", MergeState: "BEHIND"}, nil
	}
	f.resolver.results = []ConflictResolution{{Head: "conflict-merge"}}
	item := f.check(t)
	if item.Stage != store.DeliveryStageGating || !reflect.DeepEqual(item.CandidateCommits, []string{"candidate", "conflict-merge"}) || len(f.resolver.calls) != 1 || !strings.Contains(f.log.String(), "roundfix: derived merge:") || strings.Contains(f.log.String(), "roundfix: base update:") {
		t.Fatalf("item=%+v calls=%v log=%q", item, f.resolver.calls, f.log.String())
	}
	// Three subsequent Base Updates must still be available after a conflict merge.
	f.boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) {
		stored := f.item(t)
		return CheckReport{HeadSHA: stored.CandidateCommits[len(stored.CandidateCommits)-1], MergeState: "BEHIND"}, nil
	}
	f.resolver.results = []ConflictResolution{{Head: "unused"}, {Head: "base-1"}, {Head: "base-2"}, {Head: "base-3"}}
	for range BaseUpdateLimit {
		f.resumeChecking(t)
		if item := f.check(t); item.Stage != store.DeliveryStageGating {
			t.Fatalf("item=%+v", item)
		}
	}
}

func TestABaseUpdateErrorKeepsItsCauseAndAParkResetsTheBound(t *testing.T) {
	t.Parallel()
	f := newBaseUpdateFixture(t, "BEHIND")
	f.resolver.results = []ConflictResolution{{Head: "updated-1"}, {Head: "updated-2"}, {Head: "updated-3"}, {Head: "updated-4"}}
	for range 2 {
		f.check(t)
		f.resumeChecking(t)
	}
	cause := errors.New("fetch unavailable")
	f.resolver.err = cause
	item := f.item(t)
	err := f.engine.checkCandidate(t.Context(), "/repo", &item)
	if !errors.Is(err, cause) || err.Error() != "bring pull request up to date: fetch unavailable" {
		t.Fatalf("error=%v", err)
	}
	if stored := f.item(t); stored.Stage != store.DeliveryStageChecking {
		t.Fatalf("item=%+v", stored)
	}
	if _, err := f.engine.Run(t.Context(), "/repo"); err != nil {
		t.Fatal(err)
	}
	if parked := f.item(t); parked.Stage != store.DeliveryStageParked || parked.Blocker != "delivery-error: bring pull request up to date: fetch unavailable" {
		t.Fatalf("step-error park=%+v", parked)
	}
	f.resolver.err = nil
	f.resolver.results = append(f.resolver.results, ConflictResolution{Head: "updated-5"}, ConflictResolution{Head: "updated-6"})
	// There were two updates before the error; two more must fit after its park.
	for range 2 {
		f.resumeChecking(t)
		if item := f.check(t); item.Stage != store.DeliveryStageGating {
			t.Fatalf("item=%+v", item)
		}
	}
}
