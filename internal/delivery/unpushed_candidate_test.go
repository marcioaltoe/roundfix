// Suite: checking a candidate its Pull Request has not received.
// Invariant: only an earlier candidate returns through the gate before publication.
// Boundary IN: the real Delivery Queue store and Delivery Engine.
// Boundary OUT: repository gate, publication and Pull Requests use offline fakes.
package delivery

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestCheckingReturnsACandidateThePullRequestLacksToTheGate(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name string
		head string
	}{
		{name: "earlier candidate", head: "pushed-head"},
		{name: "foreign head", head: "foreign-head"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			runStore := openDeliveryEngineStore(t, ctx)
			const (
				gitRoot = "/repo-unpushed-candidate"
				slug    = "unpushed-candidate"
			)
			queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{slug})
			if err != nil {
				t.Fatalf("create Delivery Queue: %v", err)
			}
			item := queue.Items[0]
			item.Stage = store.DeliveryStageChecking
			item.Branch = "roundfix/deliver-" + slug
			item.Worktree = "/worktrees/" + slug
			item.CandidateCommits = []string{"pushed-head", "resolved-head"}
			item.PullRequestNumber = "318"
			recordDeliveryItemWorkspace(t, ctx, runStore, gitRoot, &item)
			if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
				t.Fatalf("seed checking item: %v", err)
			}

			workflow := newFakeDeliveryWorkflow()
			boundary := newFakeDeliveryBoundary()
			boundary.recordEvent = workflow.recordEvent
			boundary.remoteHeads[item.Branch] = tt.head
			boundary.pullRequests[item.Branch] = PullRequest{
				Number: "318", State: "OPEN", HeadBranch: item.Branch, HeadSHA: tt.head,
			}
			boundary.findPullRequest = func(req PullRequestRequest) PullRequestResult {
				pullRequest := boundary.pullRequests[req.HeadBranch]
				pullRequest.HeadSHA = boundary.remoteHeads[req.HeadBranch]
				boundary.pullRequests[req.HeadBranch] = pullRequest
				return PullRequestResult{PullRequest: pullRequest}
			}
			var log bytes.Buffer
			engine := NewEngine(runStore, EngineDependencies{
				Workspace: workflow, Runner: workflow, Reviewer: workflow,
				Archiver: workflow, Gate: workflow, Authorizer: workflow,
				Publication: workflow, PullRequests: boundary,
				Revalidator: cleanTestRevalidator{}, Log: &log,
			})
			if _, err := engine.Run(ctx, gitRoot); err != nil {
				t.Fatalf("run Delivery Engine: %v", err)
			}
			got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
			if tt.name == "foreign head" {
				if got.Stage != store.DeliveryStageParked || got.Blocker != BlockerReviewStale || boundary.pushEffects != 0 || boundary.mergeEffects != 0 {
					t.Fatalf("foreign-head item = %+v pushes=%d merges=%d, want review-stale without publication or merge", got, boundary.pushEffects, boundary.mergeEffects)
				}
				if boundary.remoteHeads[item.Branch] != tt.head || !reflect.DeepEqual(workflow.events[slug], []string{"checks"}) {
					t.Fatalf("foreign-head remote=%q events=%v, want unchanged remote and checks only", boundary.remoteHeads[item.Branch], workflow.events[slug])
				}
				if strings.Contains(log.String(), "returning to gating to push it") {
					t.Fatalf("foreign head logged a return to gating: %s", log.String())
				}
				return
			}
			if got.Stage != store.DeliveryStageMerged || got.Blocker != "" || got.MergeCommit != "merge-318" || got.PullRequestNumber != "318" || boundary.pushEffects != 1 || boundary.mergeEffects != 1 {
				t.Fatalf("earlier-candidate item = %+v pushes=%d merges=%d, want merged after one push", got, boundary.pushEffects, boundary.mergeEffects)
			}
			if remote := boundary.remoteHeads[item.Branch]; remote != "resolved-head" {
				t.Fatalf("remote head = %q, want resolved-head", remote)
			}
			wantEvents := []string{"checks", "gate", "authorization", "publication", "push", "pull-request", "checks", "authorization", "merge"}
			if !reflect.DeepEqual(workflow.events[slug], wantEvents) {
				t.Fatalf("item actions = %v, want %v", workflow.events[slug], wantEvents)
			}
			wantLog := "roundfix: Delivery Queue item unpushed-candidate: pull request #318 holds earlier candidate pushed-head, not candidate resolved-head; returning to gating to push it\n"
			if !strings.Contains(log.String(), wantLog) {
				t.Fatalf("owner log = %q, want line %q", log.String(), wantLog)
			}
			assertNoUnmatchedDeliveryIntents(t, ctx, runStore, gitRoot)
		})
	}
}
