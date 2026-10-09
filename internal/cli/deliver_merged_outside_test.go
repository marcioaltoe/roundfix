// Suite: Delivery Retry observes merges made outside the queue.
// Invariant: recorded PR evidence takes precedence over local delivery evidence;
// absent evidence preserves normal retry or refuses a closed PR without mutation.
// Boundary IN: CLI dispatch, real local Git and SQLite, and the owner cleanup pass.
// Boundary OUT: gh metadata and detached owner launch, both scripted.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

type outsideMergeRunner struct {
	t       *testing.T
	repo    string
	payload map[string]any
	failure bool
	calls   int
}

func (runner *outsideMergeRunner) Run(_ context.Context, dir, name string, args ...string) (delivery.CommandResult, error) {
	runner.t.Helper()
	if name != "gh" || dir != runner.repo || len(args) != 5 || !reflect.DeepEqual(args[:4], []string{"pr", "view", "404", "--json"}) {
		runner.t.Fatalf("unexpected external command: %s in %s %q", name, dir, args)
	}
	runner.calls++
	if runner.failure {
		return delivery.CommandResult{ExitCode: 1, Stderr: "scripted read failure"}, nil
	}
	data, err := json.Marshal(runner.payload)
	if err != nil {
		runner.t.Fatal(err)
	}
	return delivery.CommandResult{Stdout: string(data)}, nil
}

type outsideMergeFixture struct {
	home, repo, candidate string
	workflow              *commandDeliveryWorkflow
	runner                *outsideMergeRunner
	item                  store.DeliveryQueueItem
	starts                int
}

func newOutsideMergeFixture(t *testing.T, recordedPR bool) *outsideMergeFixture {
	t.Helper()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: "completed"}})
	gitImplement(t, repo, "checkout", "main")
	mustWrite(t, filepath.Join(repo, "docs", "specs", ".gitkeep"), "")
	gitImplement(t, repo, "add", "docs/specs/.gitkeep")
	gitImplement(t, repo, "commit", "-m", "keep Specs Root after archive")
	origin := t.TempDir()
	gittest.InitRepo(t, origin, "--bare", "--initial-branch=main")
	gitImplement(t, repo, "remote", "add", "origin", origin)
	gitImplement(t, repo, "push", "-u", "origin", "main")
	gitImplement(t, repo, "remote", "set-head", "origin", "main")
	db, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	queue, err := db.CreateDeliveryQueue(t.Context(), repo, []string{implementTestSlug})
	if err != nil {
		t.Fatal(err)
	}
	cfg := roundconfig.Builtin()
	cfg.Worktree.Location = filepath.Join(t.TempDir(), "worktrees")
	runner := &outsideMergeRunner{t: t, repo: repo}
	workflow := &commandDeliveryWorkflow{store: db, loaded: roundconfig.Loaded{GitRoot: repo, HomeDir: home, Config: cfg}, git: preflight.ExecGitRunner{}, gh: delivery.GitHubCLI{WorkDir: repo, Runner: runner}}
	branch := deliveryBranchPrefix + implementTestSlug + "-0123456789abcdef"
	ref, err := runworktree.ItemRefFor(repo, cfg.Worktree.Location, branch)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(ref.Path), 0o755); err != nil {
		t.Fatal(err)
	}
	gitImplement(t, repo, "worktree", "add", "-b", branch, ref.Path, "main")
	archive := filepath.Join(ref.Path, "docs", "history", "specs", implementTestSlug)
	if err := os.MkdirAll(filepath.Dir(archive), 0o755); err != nil {
		t.Fatal(err)
	}
	gitImplement(t, ref.Path, "mv", filepath.Join("docs", "specs", implementTestSlug), archive)
	gitImplement(t, ref.Path, "commit", "-m", "archive Spec on item branch")
	candidate := strings.TrimSpace(gitImplementOutput(t, ref.Path, "rev-parse", "HEAD"))
	if _, _, _, err := db.RecordDeliveryQueueItemWorktree(t.Context(), repo, implementTestSlug, branch, ref.Path); err != nil {
		t.Fatal(err)
	}
	if err := db.SetDeliveryQueueItemWorktreeProvisioned(t.Context(), repo, implementTestSlug, true); err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage, item.Blocker = store.DeliveryStageParked, delivery.BlockerChecksFailed
	item.Branch, item.Worktree, item.WorktreeProvisioned = branch, ref.Path, true
	item.CandidateCommits = []string{candidate}
	item.Warning = "premise-changed: recorded warning"
	if recordedPR {
		item.PullRequestNumber = "404"
	}
	if err := db.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	fixture := &outsideMergeFixture{home: home, repo: repo, candidate: candidate, workflow: workflow, runner: runner, item: item}
	runner.payload = map[string]any{"number": 404, "state": "OPEN", "headRefName": branch, "headRefOid": candidate}
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.newDeliveryEngine = func(db *store.Store, _ roundconfig.Loaded) deliveryEngine {
			copyWorkflow := *workflow
			copyWorkflow.store = db
			return outsideMergeEngine(&copyWorkflow)
		}
		deps.startDeliveryOwner = func(_ context.Context, _ roundconfig.Loaded, _ commandEnvironment, stdout, _ io.Writer) int {
			fixture.starts++
			fmt.Fprintln(stdout, "Delivery Owner: test-owner")
			return exitOK
		}
	})
	return fixture
}

func outsideMergeEngine(workflow *commandDeliveryWorkflow) *delivery.Engine {
	return delivery.NewEngine(workflow.store, delivery.EngineDependencies{
		Merges: workflow, Workspace: workflow, Runner: workflow, Reviewer: workflow,
		Archiver: workflow, Gate: workflow, Authorizer: workflow, Publication: workflow,
		PullRequests: workflow.gh, Checks: workflow.gh, Recovery: workflow,
		History: workflow, Revalidator: workflow, Prerequisites: workflow,
	})
}

func (fixture *outsideMergeFixture) retry(t *testing.T) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"deliver", "retry", implementTestSlug}, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestDeliverRetryRecordsAPullRequestMergedByHandAfterTheItemBranchIsGone(t *testing.T) {
	t.Parallel()
	fixture := newOutsideMergeFixture(t, true)
	mustWrite(t, filepath.Join(fixture.item.Worktree, "ci-fix.txt"), "fixed CI\n")
	gitImplement(t, fixture.item.Worktree, "add", "ci-fix.txt")
	gitImplement(t, fixture.item.Worktree, "commit", "-m", "fix CI after candidate")
	head := strings.TrimSpace(gitImplementOutput(t, fixture.item.Worktree, "rev-parse", "HEAD"))
	gitImplement(t, fixture.repo, "merge", "--no-ff", fixture.item.Branch, "-m", "merge by hand")
	merge := strings.TrimSpace(gitImplementOutput(t, fixture.repo, "rev-parse", "HEAD"))
	gitImplement(t, fixture.repo, "push", "origin", "main")
	gitImplement(t, fixture.repo, "worktree", "remove", fixture.item.Worktree)
	gitImplement(t, fixture.repo, "branch", "-D", fixture.item.Branch)
	fixture.runner.payload["state"], fixture.runner.payload["headRefOid"] = "MERGED", head
	fixture.runner.payload["mergeCommit"] = map[string]string{"oid": merge}
	code, stdout, stderr := fixture.retry(t)
	want := "Merged outside the queue: pull request #404; merge commit " + merge + "\nRetried " + implementTestSlug + ": checks-failed -> merged\nDelivery Owner: test-owner\n"
	if code != exitOK || stdout != want || stderr != "" || fixture.starts != 1 {
		t.Fatalf("retry: exit=%d stdout=%q stderr=%q starts=%d", code, stdout, stderr, fixture.starts)
	}
	item := readDeliveryItemForCLI(t, fixture.workflow.store, fixture.repo)
	if item.Stage != store.DeliveryStageMerged || item.MergeCommit != merge || item.RetryCount != fixture.item.RetryCount || item.Warning != fixture.item.Warning || !reflect.DeepEqual(item.CandidateCommits, []string{fixture.candidate, head}) {
		t.Fatalf("merged item: %+v", item)
	}
	if _, err := outsideMergeEngine(fixture.workflow).Run(t.Context(), fixture.repo); err != nil {
		t.Fatal(err)
	}
	item = readDeliveryItemForCLI(t, fixture.workflow.store, fixture.repo)
	if item.Stage != store.DeliveryStageMerged || item.Blocker != "" {
		t.Fatalf("owner cleanup: %+v", item)
	}
}

func TestDeliverRetryRecordsASpecArchivedOnTheDefaultBranchWithoutARecordedPullRequest(t *testing.T) {
	t.Parallel()
	fixture := newOutsideMergeFixture(t, false)
	gitImplement(t, fixture.repo, "merge", "--squash", fixture.item.Branch)
	gitImplement(t, fixture.repo, "commit", "-m", "squash delivery by hand")
	merge := strings.TrimSpace(gitImplementOutput(t, fixture.repo, "rev-parse", "HEAD"))
	code, stdout, stderr := fixture.retry(t)
	want := "Merged outside the queue: Spec archived on default branch \"main\"; merge commit " + merge + "\nRetried " + implementTestSlug + ": checks-failed -> merged\nDelivery Owner: test-owner\n"
	if code != exitOK || stdout != want || stderr != "" {
		t.Fatalf("retry: exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	item := readDeliveryItemForCLI(t, fixture.workflow.store, fixture.repo)
	if item.Stage != store.DeliveryStageMerged || item.MergeCommit != merge || !reflect.DeepEqual(item.CandidateCommits, fixture.item.CandidateCommits) || fixture.runner.calls != 0 {
		t.Fatalf("delivery evidence item: %+v gh calls=%d", item, fixture.runner.calls)
	}
}

func TestDeliverRetryRefusesAPullRequestClosedWithoutMerging(t *testing.T) {
	t.Parallel()
	fixture := newOutsideMergeFixture(t, true)
	fixture.runner.payload["state"] = "CLOSED"
	code, stdout, stderr := fixture.retry(t)
	reason := fmt.Sprintf("retry Delivery Queue item %q: pull request #404 was closed without merging; reopen it or merge the Spec into the default branch, then run roundfix deliver retry %s", implementTestSlug, implementTestSlug)
	if code != exitPreflight || stdout != "" || !strings.HasPrefix(stderr, "Retry refused\n") || !strings.Contains(stderr, reason) || !strings.Contains(stderr, "stage: parked; blocker: checks-failed") || fixture.starts != 0 {
		t.Fatalf("refusal: exit=%d stdout=%q stderr=%q starts=%d", code, stdout, stderr, fixture.starts)
	}
	if item := readDeliveryItemForCLI(t, fixture.workflow.store, fixture.repo); !reflect.DeepEqual(item, fixture.item) {
		t.Fatalf("refused item changed: %+v", item)
	}
}

func TestDeliverRetryWithAnOpenPullRequestAndNoMergeEvidenceRetriesAsBefore(t *testing.T) {
	t.Parallel()
	fixture := newOutsideMergeFixture(t, true)
	code, stdout, stderr := fixture.retry(t)
	want := "Retried " + implementTestSlug + ": checks-failed -> checking\nDelivery Owner: test-owner\n"
	if code != exitOK || stdout != want || stderr != "" || fixture.starts != 1 {
		t.Fatalf("open retry: exit=%d stdout=%q stderr=%q starts=%d", code, stdout, stderr, fixture.starts)
	}
	item := readDeliveryItemForCLI(t, fixture.workflow.store, fixture.repo)
	if item.Stage != store.DeliveryStageChecking || item.RetryCount != fixture.item.RetryCount+1 || item.MergeCommit != "" {
		t.Fatalf("normal retry: %+v", item)
	}
}

func TestObserveMergeIgnoresAMergedPullRequestFromAnotherBranch(t *testing.T) {
	t.Parallel()
	fixture := newOutsideMergeFixture(t, true)
	fixture.runner.payload["state"], fixture.runner.payload["headRefName"] = "MERGED", "fix/another-spec"
	fixture.runner.payload["mergeCommit"] = map[string]string{"oid": fixture.candidate}
	observation, err := fixture.workflow.ObserveMerge(t.Context(), fixture.repo, fixture.item)
	if err != nil || observation.Merged || observation.ClosedUnmerged {
		t.Fatalf("unrelated PR: %+v error=%v", observation, err)
	}
}

func TestObserveMergeEvidenceFallbacks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                                                                             string
		recordedPR, closed, readFailure, removeBranch, noCandidate, itemContainsDelivery bool
		wantMerged, wantClosed, wantError                                                bool
	}{
		{name: "closed PR with delivery proof", recordedPR: true, closed: true, wantMerged: true},
		{name: "unreadable PR refuses even with delivery proof", recordedPR: true, readFailure: true, wantError: true},
		{name: "removed branch uses newest candidate", removeBranch: true, wantMerged: true},
		{name: "no branch and no candidate cannot prove merge", removeBranch: true, noCandidate: true},
		{name: "closed PR without an item head stays closed", recordedPR: true, closed: true, removeBranch: true, noCandidate: true, wantClosed: true},
		{name: "branch tip already contains delivery; older candidate does not", itemContainsDelivery: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newOutsideMergeFixture(t, test.recordedPR)
			gitImplement(t, fixture.repo, "merge", "--squash", fixture.item.Branch)
			gitImplement(t, fixture.repo, "commit", "-m", "squash delivery by hand")
			merge := strings.TrimSpace(gitImplementOutput(t, fixture.repo, "rev-parse", "HEAD"))
			if test.closed {
				fixture.runner.payload["state"] = "CLOSED"
			}
			fixture.runner.failure = test.readFailure
			if test.itemContainsDelivery {
				gitImplement(t, fixture.item.Worktree, "merge", "--no-ff", "main", "-m", "bring delivery into item")
			}
			if test.removeBranch {
				gitImplement(t, fixture.repo, "worktree", "remove", fixture.item.Worktree)
				gitImplement(t, fixture.repo, "branch", "-D", fixture.item.Branch)
			}
			if test.noCandidate {
				fixture.item.CandidateCommits = nil
			}
			observation, err := fixture.workflow.ObserveMerge(t.Context(), fixture.repo, fixture.item)
			if (err != nil) != test.wantError || observation.Merged != test.wantMerged || observation.ClosedUnmerged != test.wantClosed {
				t.Fatalf("observation=%+v error=%v", observation, err)
			}
			if test.wantMerged && (observation.Head != fixture.candidate || observation.MergeCommit != merge || observation.Evidence != `Spec archived on default branch "main"`) {
				t.Fatalf("fallback proof=%+v", observation)
			}
		})
	}
}

func TestObserveMergeRequiresCompletePullRequestEvidence(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"headRefOid", "mergeCommit"} {
		t.Run(field, func(t *testing.T) {
			fixture := newOutsideMergeFixture(t, true)
			fixture.runner.payload["state"] = "MERGED"
			fixture.runner.payload["mergeCommit"] = map[string]string{"oid": fixture.candidate}
			delete(fixture.runner.payload, field)
			observation, err := fixture.workflow.ObserveMerge(t.Context(), fixture.repo, fixture.item)
			if err != nil || observation.Merged || observation.ClosedUnmerged {
				t.Fatalf("incomplete PR: %+v error=%v", observation, err)
			}
		})
	}
}
