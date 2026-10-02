// Suite: GitHub merge readiness and bounded policy recovery.
// Invariant: listed check success cannot bypass a blocked merge state or repeat policy recovery indefinitely.
// Boundary IN: real queue persistence and Engine Run/Retry with a fake clock.
// Boundary OUT: GitHub commands, item workspace and archived evidence.
package delivery

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"roundfix/internal/store"
)

func TestPullRequestReadCarriesTheMergeState(t *testing.T) {
	t.Parallel()
	view := func(state string) commandStep {
		return commandStep{name: "gh", args: []string{"pr", "view", "1", "--json", pullRequestJSONFields}, result: CommandResult{Stdout: fmt.Sprintf(`{"number":1,"headRefOid":"candidate","mergeable":"MERGEABLE","mergeStateStatus":%q}`, state)}}
	}
	runner := newScriptedCommandRunner(t, view(" blocked "), commandStep{name: "gh", args: []string{"pr", "checks", "1", "--json", "bucket,link,name,state,workflow"}, result: CommandResult{Stdout: `[{"name":"test","bucket":"pass"}]`}}, view(" clean "))
	client := GitHubCLI{WorkDir: "/repo", Runner: runner}
	report, err := client.CurrentHeadChecks(t.Context(), "1")
	if err != nil || report.MergeState != "CLEAN" {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	pr, err := parsePullRequest([]byte(`{"number":1,"mergeStateStatus":" blocked "}`))
	if err != nil || pr.MergeState != "BLOCKED" || !strings.Contains(pullRequestJSONFields, "mergeStateStatus") {
		t.Fatalf("PR=%+v fields=%q err=%v", pr, pullRequestJSONFields, err)
	}
}

func TestAPolicyRefusalIsTypedFromGhStderr(t *testing.T) {
	t.Parallel()
	for _, detail := range []string{"the base branch policy prohibits the merge", "THE BASE BRANCH POLICY PROHIBITS THE MERGE", "permission denied"} {
		t.Run(detail, func(t *testing.T) {
			runner := newScriptedCommandRunner(t, commandStep{name: "gh", args: []string{"pr", "view", "1", "--json", pullRequestJSONFields}, result: CommandResult{Stdout: `{"number":1,"headRefOid":"candidate"}`}}, commandStep{name: "gh", args: []string{"pr", "merge", "1", "--squash", "--match-head-commit", "candidate"}, result: CommandResult{ExitCode: 1, Stderr: detail}})
			_, err := (GitHubCLI{WorkDir: "/repo", Runner: runner}).MergePullRequest(t.Context(), "1", "candidate")
			var refusal MergePolicyRefusalError
			wantTyped := strings.Contains(strings.ToLower(detail), "base branch policy prohibits the merge")
			if err == nil || errors.As(err, &refusal) != wantTyped || !strings.Contains(err.Error(), detail) {
				t.Fatalf("merge error=%v typed=%v", err, wantTyped)
			}
		})
	}
}

type mergeStateFixture struct {
	engine   *Engine
	runStore *store.Store
	boundary *fakePullRequestBoundary
	clock    *fakeDeliveryClock
	log      bytes.Buffer
	polls    int
	attempts int
	merges   int
}

func newMergeStateFixture(t *testing.T, states []string, refusals int) *mergeStateFixture {
	t.Helper()
	ctx := t.Context()
	runStore := openDeliveryEngineStore(t, ctx)
	item := seedParkedRetryItem(t, ctx, runStore, "/repo", "example", BlockerDeliveryError)
	item.Stage = store.DeliveryStageChecking
	item.Blocker = ""
	item.PullRequestNumber = "1"
	if err := runStore.UpdateDeliveryQueueItem(ctx, "/repo", item); err != nil {
		t.Fatal(err)
	}
	f := &mergeStateFixture{runStore: runStore, clock: &fakeDeliveryClock{now: time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC)}}
	f.boundary = &fakePullRequestBoundary{
		currentHeadChecks: func(_ context.Context, number string) (CheckReport, error) {
			if number != "1" {
				t.Fatalf("checks PR=%q", number)
			}
			stored := readDeliveryQueue(t, ctx, runStore, "/repo").Items[0]
			if stored.Stage != store.DeliveryStageChecking {
				t.Fatalf("poll stage=%s", stored.Stage)
			}
			state := states[min(f.polls, len(states)-1)]
			f.polls++
			return CheckReport{HeadSHA: "candidate-original", MergeState: state, Checks: []PullRequestCheck{{Name: "test", Bucket: "pass"}}}, nil
		},
		mergePullRequest: func(_ context.Context, number, head string) (MergeResult, error) {
			if number != "1" || head != "candidate-original" {
				t.Fatalf("merge PR=%q head=%q", number, head)
			}
			if f.polls < len(states) {
				t.Fatalf("merge before final state: polls=%d states=%v", f.polls, states)
			}
			f.attempts++
			if f.attempts <= refusals {
				return MergeResult{}, MergePolicyRefusalError{Stderr: "the base branch policy prohibits the merge"}
			}
			f.merges++
			return MergeResult{PullRequest: PullRequest{Number: number, HeadSHA: head, State: "MERGED", MergeCommit: "merge-1"}}, nil
		},
	}
	workflow := newFakeDeliveryWorkflow()
	f.engine = NewEngine(runStore, EngineDependencies{Workspace: workflow, Runner: workflow, Reviewer: workflow, Archiver: workflow, Gate: workflow, Authorizer: workflow, Publication: workflow, PullRequests: f.boundary, Recovery: &fakeItemRecovery{states: []ItemState{{Archived: true, Head: "candidate-original"}}}, Revalidator: cleanTestRevalidator{}, Clock: f.clock, Sleeper: f.clock, Log: &f.log, CheckTimeout: 3 * time.Second, CheckInterval: time.Second})
	return f
}

func (f *mergeStateFixture) run(t *testing.T) store.DeliveryQueueItem {
	t.Helper()
	if _, err := f.engine.Run(t.Context(), "/repo"); err != nil {
		t.Fatal(err)
	}
	return readDeliveryQueue(t, t.Context(), f.runStore, "/repo").Items[0]
}

func TestCheckingWaitsWhileTheMergeStateIsBlocked(t *testing.T) {
	t.Parallel()
	f := newMergeStateFixture(t, []string{"BLOCKED", "BLOCKED", "CLEAN"}, 0)
	item := f.run(t)
	if item.Stage != store.DeliveryStageMerged || f.polls != 3 || f.merges != 1 || f.clock.sleeps != 2 {
		t.Fatalf("item=%+v polls=%d merges=%d sleeps=%d", item, f.polls, f.merges, f.clock.sleeps)
	}
	for _, state := range []string{"BLOCKED", "CLEAN"} {
		if strings.Count(f.log.String(), "roundfix: checks: Delivery Queue item example: merge state "+state+"\n") != 1 {
			t.Fatalf("state log=%q", f.log.String())
		}
	}
	assertNoUnmatchedDeliveryIntents(t, t.Context(), f.runStore, "/repo")
}

func TestCheckingWaitsWhileTheMergeStateIsUnknown(t *testing.T) {
	t.Parallel()
	f := newMergeStateFixture(t, []string{"UNKNOWN", "CLEAN"}, 0)
	item := f.run(t)
	if item.Stage != store.DeliveryStageMerged || f.polls != 2 || f.merges != 1 || f.clock.sleeps != 1 {
		t.Fatalf("item=%+v polls=%d merges=%d", item, f.polls, f.merges)
	}
}

func TestABlockedMergeStateParksAtTheChecksTimeout(t *testing.T) {
	t.Parallel()
	f := newMergeStateFixture(t, []string{"BLOCKED"}, 0)
	start := f.clock.Now()
	item := f.run(t)
	if item.Stage != store.DeliveryStageParked || item.Blocker != BlockerChecksTimeout || f.attempts != 0 || f.clock.Now().Sub(start) != 3*time.Second {
		t.Fatalf("item=%+v attempts=%d elapsed=%s", item, f.attempts, f.clock.Now().Sub(start))
	}
}

func TestAnEmptyMergeStateKeepsTheListedChecksDecision(t *testing.T) {
	t.Parallel()
	f := newMergeStateFixture(t, []string{""}, 0)
	item := f.run(t)
	if item.Stage != store.DeliveryStageMerged || f.polls != 1 || f.merges != 1 || f.clock.sleeps != 0 {
		t.Fatalf("item=%+v polls=%d merges=%d", item, f.polls, f.merges)
	}
}

func TestAPolicyRefusedMergeReturnsToChecking(t *testing.T) {
	t.Parallel()
	f := newMergeStateFixture(t, []string{"CLEAN"}, 1)
	item := f.run(t)
	if item.Stage != store.DeliveryStageMerged || f.polls != 2 || f.attempts != 2 || f.merges != 1 || !strings.Contains(f.log.String(), "roundfix: merge refused by branch policy: Delivery Queue item example; checking again\n") {
		t.Fatalf("item=%+v polls=%d attempts=%d log=%q", item, f.polls, f.attempts, f.log.String())
	}
	assertNoUnmatchedDeliveryIntents(t, t.Context(), f.runStore, "/repo")
}

func TestASecondPolicyRefusalAtTheSameHeadParks(t *testing.T) {
	t.Parallel()
	f := newMergeStateFixture(t, []string{"CLEAN"}, 2)
	item := f.run(t)
	if item.Stage != store.DeliveryStageParked || !strings.HasPrefix(item.Blocker, BlockerDeliveryError+":") || f.attempts != 2 || f.merges != 0 || strings.Count(f.log.String(), "checking again") != 1 {
		t.Fatalf("item=%+v attempts=%d log=%q", item, f.attempts, f.log.String())
	}
}

func TestADeliveryErrorRetryAtAnUnchangedArchivedHeadResumesChecking(t *testing.T) {
	t.Parallel()
	f := newMergeStateFixture(t, []string{"CLEAN"}, 2)
	parked := f.run(t)
	if parked.Stage != store.DeliveryStageParked || !strings.HasPrefix(parked.Blocker, BlockerDeliveryError+":") {
		t.Fatalf("park=%+v", parked)
	}
	result, err := f.engine.Retry(t.Context(), "/repo", "example")
	if err != nil {
		t.Fatalf("retry refused: %v", err)
	}
	retried := readDeliveryQueue(t, t.Context(), f.runStore, "/repo").Items[0]
	if result.Stage != store.DeliveryStageChecking || retried.Stage != store.DeliveryStageChecking || retried.RetryCount != 1 || retried.PullRequestNumber != parked.PullRequestNumber {
		t.Fatalf("retry=%+v item=%+v", result, retried)
	}
	item := f.run(t)
	if item.Stage != store.DeliveryStageMerged || f.merges != 1 || f.attempts != 3 {
		t.Fatalf("item=%+v attempts=%d merges=%d", item, f.attempts, f.merges)
	}
	assertNoUnmatchedDeliveryIntents(t, t.Context(), f.runStore, "/repo")
}
