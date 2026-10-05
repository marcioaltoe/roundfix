// Suite: Delivery Retry merge observations.
// Invariant: merge evidence records a parked item before limits or workspace recovery.
// Boundary IN: Engine.Retry and the real SQLite store in a temporary home.
// Boundary OUT: merge observations, workspace, recovery and revalidation use fakes; no GitHub access.
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

type fakeMergeObserver struct {
	observation MergeObservation
	err         error
	calls       int
	root        string
	item        store.DeliveryQueueItem
}

func (fake *fakeMergeObserver) ObserveMerge(_ context.Context, root string, item store.DeliveryQueueItem) (MergeObservation, error) {
	fake.calls++
	fake.root, fake.item = root, item
	return fake.observation, fake.err
}

func mergeRetryEngine(writer *store.Store, observer MergeObserver) (*Engine, *countingRetryWorkspace, *fakeItemRecovery, *recordingItemRevalidator) {
	workflow := newFakeDeliveryWorkflow()
	workspace := &countingRetryWorkspace{fakeDeliveryWorkflow: workflow}
	recovery := &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "candidate-original"}}}
	revalidator := &recordingItemRevalidator{}
	engine := NewEngine(writer, EngineDependencies{Merges: observer, Workspace: workspace, Recovery: recovery, Revalidator: revalidator})
	return engine, workspace, recovery, revalidator
}

func assertMergedRetry(t *testing.T, writer *store.Store, root string, before store.DeliveryQueueItem) {
	t.Helper()
	observation := MergeObservation{Merged: true, MergeCommit: "merge-commit", Head: "merged-head", Evidence: "pull request #404"}
	observer := &fakeMergeObserver{observation: observation}
	engine, workspace, recovery, revalidator := mergeRetryEngine(writer, observer)
	if err := writer.ClaimDeliveryQueueOwner(t.Context(), root, 4242, "owner-identity"); err != nil {
		t.Fatal(err)
	}
	result, err := engine.Retry(t.Context(), root, before.SpecSlug)
	if err != nil {
		t.Fatal(err)
	}
	want := before
	want.Stage, want.Blocker, want.MergeCommit = store.DeliveryStageMerged, "", observation.MergeCommit
	want.CandidateCommits = append(append([]string{}, before.CandidateCommits...), observation.Head)
	got := readDeliveryQueue(t, t.Context(), writer, root).Items[before.Position]
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("merged item = %#v, want %#v", got, want)
	}
	if result.SpecSlug != before.SpecSlug || result.Blocker != before.Blocker || result.Stage != store.DeliveryStageMerged || result.Merge != observation || result.OwnerPID != 4242 || result.OwnerIdentity != "owner-identity" || !reflect.DeepEqual(result.CarriedFrom, CarryForwardResult{}) {
		t.Fatalf("result = %+v", result)
	}
	if observer.calls != 1 || observer.root != root || !reflect.DeepEqual(observer.item, before) {
		t.Fatalf("observer = %+v", observer)
	}
	if workspace.useCalls != 0 || len(recovery.events) != 0 || len(revalidator.priorMerges) != 0 || len(workspace.events) != 0 {
		t.Fatalf("merge performed work: workspace=%d recovery=%v revalidation=%v events=%v", workspace.useCalls, recovery.events, revalidator.priorMerges, workspace.events)
	}
}

func TestARetryRecordsAMergedPullRequestWithoutTouchingTheWorkspace(t *testing.T) {
	t.Parallel()
	for _, blocker := range []string{BlockerChecksFailed, BlockerItemWorktreeMissing, BlockerPrerequisiteUnmerged + ": prerequisite", BlockerCorrectiveSpecRequired + ": archived-spec", BlockerPullRequestConflict} {
		t.Run(blocker, func(t *testing.T) {
			writer := openDeliveryEngineStore(t, t.Context())
			root := "/repo-merged-outside"
			item := seedParkedRetryItem(t, t.Context(), writer, root, "merged-spec", blocker)
			item.PullRequestNumber, item.Warning = "404", "preserved-warning"
			if err := writer.UpdateDeliveryQueueItem(t.Context(), root, item); err != nil {
				t.Fatal(err)
			}
			assertMergedRetry(t, writer, root, item)
		})
	}
}

func TestARetryRecordsAMergeAtTheRetryLimitDeadlineAndTokenCeiling(t *testing.T) {
	t.Parallel()
	t.Run("retry limit and expired deadline", func(t *testing.T) {
		writer := openDeliveryEngineStore(t, t.Context())
		root := "/repo-merged-limits"
		queue, err := writer.CreateDeliveryQueueWithLimits(t.Context(), root, []string{"merged-spec"}, store.DeliveryQueueLimits{MaxRetries: 1, Deadline: time.Now().Add(-time.Hour)})
		if err != nil {
			t.Fatal(err)
		}
		item := queue.Items[0]
		item.Stage, item.Blocker = store.DeliveryStageParked, BlockerChecksFailed
		if err := writer.UpdateDeliveryQueueItem(t.Context(), root, item); err != nil {
			t.Fatal(err)
		}
		item.Stage = store.DeliveryStageQueued
		if _, _, err := writer.RetryDeliveryQueueItem(t.Context(), root, item, item.Blocker); err != nil {
			t.Fatal(err)
		}
		item = readDeliveryQueue(t, t.Context(), writer, root).Items[0]
		item.Stage, item.Blocker = store.DeliveryStageParked, BlockerQueueDeadline
		if err := writer.UpdateDeliveryQueueItem(t.Context(), root, item); err != nil {
			t.Fatal(err)
		}
		assertMergedRetry(t, writer, root, item)
	})
	t.Run("token ceiling", func(t *testing.T) {
		writer, root, item := seedTokenCeilingQueue(t, 5000000, 5000000)
		item.Stage, item.Blocker = store.DeliveryStageParked, BlockerQueueTokenCeiling
		if err := writer.UpdateDeliveryQueueItem(t.Context(), root, item); err != nil {
			t.Fatal(err)
		}
		assertMergedRetry(t, writer, root, item)
	})
}

func TestARetryDoesNotAppendAnAlreadyNewestMergedHead(t *testing.T) {
	t.Parallel()
	writer := openDeliveryEngineStore(t, t.Context())
	root := "/repo-merged-newest"
	item := seedParkedRetryItem(t, t.Context(), writer, root, "merged-spec", BlockerChecksFailed)
	observer := &fakeMergeObserver{observation: MergeObservation{Merged: true, MergeCommit: "merge-commit", Head: "candidate-original"}}
	engine, _, _, _ := mergeRetryEngine(writer, observer)
	if _, err := engine.Retry(t.Context(), root, item.SpecSlug); err != nil {
		t.Fatal(err)
	}
	got := readDeliveryQueue(t, t.Context(), writer, root).Items[0]
	if !reflect.DeepEqual(got.CandidateCommits, item.CandidateCommits) {
		t.Fatalf("candidates = %v, want %v", got.CandidateCommits, item.CandidateCommits)
	}
}

func TestARetryOfAClosedUnmergedPullRequestIsRefusedAndLeavesTheItemParked(t *testing.T) {
	t.Parallel()
	writer := openDeliveryEngineStore(t, t.Context())
	root := "/repo-closed-unmerged"
	item := seedParkedRetryItem(t, t.Context(), writer, root, "closed-spec", BlockerChecksFailed)
	item.PullRequestNumber = "404"
	if err := writer.UpdateDeliveryQueueItem(t.Context(), root, item); err != nil {
		t.Fatal(err)
	}
	engine, workspace, recovery, _ := mergeRetryEngine(writer, &fakeMergeObserver{observation: MergeObservation{ClosedUnmerged: true}})
	_, err := engine.Retry(t.Context(), root, item.SpecSlug)
	want := fmt.Sprintf("retry Delivery Queue item %q: pull request #404 was closed without merging; reopen it or merge the Spec into the default branch, then run roundfix deliver retry %s", item.SpecSlug, item.SpecSlug)
	if err == nil || err.Error() != want {
		t.Fatalf("refusal = %v, want %s", err, want)
	}
	assertRetryItemUnchanged(t, t.Context(), writer, root, item)
	if workspace.useCalls != 0 || len(recovery.events) != 0 {
		t.Fatal("closed observation performed work")
	}
}

func TestARetryWhoseMergeObservationFailsIsRefused(t *testing.T) {
	t.Parallel()
	sentinel := errors.New("cannot read pull request")
	for _, test := range []struct {
		name        string
		observation MergeObservation
		err         error
	}{
		{name: "read failure", err: sentinel},
		{name: "missing merge commit", observation: MergeObservation{Merged: true, Head: "head"}},
		{name: "missing head", observation: MergeObservation{Merged: true, MergeCommit: "commit"}},
		{name: "blank evidence", observation: MergeObservation{Merged: true, MergeCommit: " ", Head: " "}},
	} {
		t.Run(test.name, func(t *testing.T) {
			writer := openDeliveryEngineStore(t, t.Context())
			root := "/repo-observation-failure"
			item := seedParkedRetryItem(t, t.Context(), writer, root, "error-spec", BlockerChecksFailed)
			engine, workspace, recovery, _ := mergeRetryEngine(writer, &fakeMergeObserver{observation: test.observation, err: test.err})
			_, err := engine.Retry(t.Context(), root, item.SpecSlug)
			if err == nil || !strings.HasPrefix(err.Error(), `retry Delivery Queue item "error-spec": observe merge: `) {
				t.Fatalf("refusal = %v", err)
			}
			if test.err != nil && !errors.Is(err, test.err) {
				t.Fatalf("unwrapped observer error: %v", err)
			}
			assertRetryItemUnchanged(t, t.Context(), writer, root, item)
			if workspace.useCalls != 0 || len(recovery.events) != 0 {
				t.Fatal("failed observation performed work")
			}
		})
	}
}

func TestARetryOfAnUnmergedItemFollowsTheExistingRules(t *testing.T) {
	t.Parallel()
	for _, withObserver := range []bool{false, true} {
		t.Run(fmt.Sprint(withObserver), func(t *testing.T) {
			writer := openDeliveryEngineStore(t, t.Context())
			root := "/repo-not-merged"
			item := seedParkedRetryItem(t, t.Context(), writer, root, "unmerged-spec", BlockerChecksFailed)
			item.PullRequestNumber = "404"
			if err := writer.UpdateDeliveryQueueItem(t.Context(), root, item); err != nil {
				t.Fatal(err)
			}
			var observer MergeObserver
			fake := &fakeMergeObserver{}
			if withObserver {
				observer = fake
			}
			engine, workspace, recovery, _ := mergeRetryEngine(writer, observer)
			result, err := engine.Retry(t.Context(), root, item.SpecSlug)
			if err != nil {
				t.Fatal(err)
			}
			got := readDeliveryQueue(t, t.Context(), writer, root).Items[0]
			want := item
			want.Stage, want.Blocker, want.RetryCount = store.DeliveryStageChecking, "", 1
			if !reflect.DeepEqual(got, want) || result.Stage != want.Stage || result.Merge != (MergeObservation{}) || workspace.useCalls != 1 || !reflect.DeepEqual(recovery.events, []string{"inspect"}) {
				t.Fatalf("retry = %+v item=%+v workspace=%d recovery=%v", result, got, workspace.useCalls, recovery.events)
			}
			if withObserver && fake.calls != 1 {
				t.Fatalf("observer calls = %d", fake.calls)
			}
		})
	}
}

func TestAMergeObserverIsNeverAskedForAnItemThatIsNotParked(t *testing.T) {
	t.Parallel()
	writer := openDeliveryEngineStore(t, t.Context())
	root := "/repo-observer-stage"
	item := seedParkedRetryItem(t, t.Context(), writer, root, "active-spec", BlockerChecksFailed)
	item.Stage = store.DeliveryStageChecking
	if err := writer.UpdateDeliveryQueueItem(t.Context(), root, item); err != nil {
		t.Fatal(err)
	}
	observer := &fakeMergeObserver{observation: MergeObservation{Merged: true, MergeCommit: "commit", Head: "head"}}
	engine, _, _, _ := mergeRetryEngine(writer, observer)
	if _, err := engine.Retry(t.Context(), root, item.SpecSlug); err == nil {
		t.Fatal("non-parked retry accepted")
	}
	if observer.calls != 0 {
		t.Fatal("observer called before stage refusal")
	}
	assertRetryItemUnchanged(t, t.Context(), writer, root, item)
}

func TestGitHubCLIViewPullRequestReportsMergedAndClosedPullRequests(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"MERGED", "CLOSED"} {
		t.Run(state, func(t *testing.T) {
			merge := "null"
			mergedAt := ""
			if state == "MERGED" {
				merge = `{"oid":"merge-sha"}`
				mergedAt = "2026-10-05T00:00:00Z"
			}
			payload := fmt.Sprintf(`{"number":404,"url":"https://example.test/404","state":%q,"headRefName":"roundfix/deliver-example","headRefOid":"head-sha","mergedAt":%q,"mergeCommit":%s,"mergeable":"MERGEABLE","mergeStateStatus":"CLEAN"}`, state, mergedAt, merge)
			runner := newScriptedCommandRunner(t, commandStep{name: "gh", args: []string{"pr", "view", "404", "--json", pullRequestJSONFields}, result: CommandResult{Stdout: payload}})
			pr, err := (GitHubCLI{WorkDir: "/repo", Runner: runner}).ViewPullRequest(t.Context(), "404")
			if err != nil {
				t.Fatal(err)
			}
			if pr.Number != "404" || pr.State != state || pr.HeadBranch != "roundfix/deliver-example" || pr.HeadSHA != "head-sha" || pr.Merged() != (state == "MERGED") || pr.MergedAt != mergedAt {
				t.Fatalf("pull request = %+v", pr)
			}
			if state == "MERGED" && pr.MergeCommit != "merge-sha" || state == "CLOSED" && pr.MergeCommit != "" {
				t.Fatalf("merge commit = %q", pr.MergeCommit)
			}
		})
	}
}
