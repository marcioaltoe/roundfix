// Boundary: real item persistence and publication through the fake forge.
package delivery

import (
	"path/filepath"
	"testing"

	"roundfix/internal/store"
)

func TestAGateLaneSkipIsRecordedOnTheItemAndInThePullRequestBody(t *testing.T) {
	t.Parallel()
	const reason = "image golang:1.26 is not present locally"
	checkLaneSkipPublication(t, "premise-changed: internal/example.go", reason,
		"premise-changed: internal/example.go; linux-lane-skipped: "+reason,
		"Delivers Spec spec-one.\n\nLinux Lane: not run. The Delivery Queue gate skipped it: image golang:1.26 is not present locally. This candidate's Go packages were not tested on Linux before this pull request.")
}

func TestAQuietGateLeavesNoLaneWarningAndReplacesAnEarlierOne(t *testing.T) {
	t.Parallel()
	const premise = "premise-changed: internal/example.go"
	const old = "linux-lane-skipped: old reason; with detail"
	for _, tc := range []struct{ name, before, reason, warning, body string }{
		{"quiet", "", "", "", "Delivers Spec spec-one."},
		{"preserve premise", premise, "", premise, "Delivers Spec spec-one."},
		{"clear old", old, "", "", "Delivers Spec spec-one."},
		{"clear old after premise", premise + "; " + old, "", premise, "Delivers Spec spec-one."},
		{"replace old", premise + "; " + old, "docker is not on PATH", premise + "; linux-lane-skipped: docker is not on PATH", "Delivers Spec spec-one.\n\nLinux Lane: not run. The Delivery Queue gate skipped it: docker is not on PATH. This candidate's Go packages were not tested on Linux before this pull request."},
	} {
		t.Run(tc.name, func(t *testing.T) { checkLaneSkipPublication(t, tc.before, tc.reason, tc.warning, tc.body) })
	}
}

func checkLaneSkipPublication(t *testing.T, before, reason, warning, body string) {
	t.Helper()
	ctx := t.Context()
	runStore := openDeliveryEngineStore(t, ctx)
	root := filepath.Join(t.TempDir(), "repo")
	queue, err := runStore.CreateDeliveryQueue(ctx, root, []string{"spec-one"})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageGating
	item.Branch = "roundfix/deliver-spec-one"
	item.Worktree = filepath.Join(t.TempDir(), "item")
	item.CandidateCommits = []string{"archived-spec-one"}
	item.Warning = before
	recordDeliveryItemWorkspace(t, ctx, runStore, root, &item)
	if err := runStore.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
		t.Fatal(err)
	}
	workflow := newFakeDeliveryWorkflow()
	workflow.gates[item.SpecSlug] = GateResult{Passed: true, LaneSkip: reason}
	workflow.publications[item.SpecSlug] = Publication{Remote: "origin", HeadBranch: item.Branch, BaseBranch: "main", Title: item.SpecSlug, Body: "Delivers Spec spec-one."}
	boundary := newFakeDeliveryBoundary()
	calls := 0
	boundary.findPullRequest = func(req PullRequestRequest) PullRequestResult {
		calls++
		if req.Body != body {
			t.Fatalf("PR body=%q want=%q", req.Body, body)
		}
		recorded := readDeliveryQueue(t, ctx, runStore, root).Items[0]
		if recorded.Warning != warning || recorded.Stage != store.DeliveryStagePublishing {
			t.Fatalf("published item=%+v want warning=%q", recorded, warning)
		}
		pr := PullRequest{Number: "1", State: "OPEN", HeadBranch: req.HeadBranch, HeadSHA: boundary.remoteHeads[req.HeadBranch]}
		boundary.pullRequests[req.HeadBranch] = pr
		return PullRequestResult{PullRequest: pr, Created: true}
	}
	engine := newTestDeliveryEngine(runStore, workflow, boundary)
	if _, err := engine.Run(ctx, root); err != nil {
		t.Fatal(err)
	}
	recorded := readDeliveryQueue(t, ctx, runStore, root).Items[0]
	if calls != 1 || recorded.Stage != store.DeliveryStageMerged || recorded.Warning != warning {
		t.Fatalf("calls=%d item=%+v", calls, recorded)
	}
}
