// Suite: unresolved Run details flow from the journal into delivery status.
// Boundary IN: real Run and queue stores, workflow Run results and public CLI output.
// Boundary OUT: status starts no owner or provider action.
package cli

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"roundfix/internal/delivery"
	"roundfix/internal/preflight"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

func TestDeliveryRunResultReadsFailedAndSkippedTasks(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	s := openDeliverReleaseStore(t, ctx, t.TempDir())
	root := t.TempDir()
	run, err := s.CreateRun(ctx, store.CreateRunRequest{Kind: store.KindImplement, GitRoot: root, LocalBranch: "fix/unresolved-detail", HeadSHA: "fixture-head", SpecSlug: "0077-example", Agent: "codex"})
	if err != nil {
		t.Fatal(err)
	}
	const reason = "Task commit carries a secret: fixture.txt:2: credential-assignment"
	for _, event := range []struct {
		workItem string
		payload  map[string]any
	}{
		{"other-work-item", map[string]any{"phase": "settled", "task": "task_01", "status": "failed", "reason": reason}},
		{"task_02", map[string]any{"phase": "skipped"}},
		{"other-work-item", map[string]any{"phase": "skipped", "task": "task_03"}},
		{"task_04", map[string]any{"phase": "settled", "status": "completed"}},
		// Legacy secret_scan events have no attempt and do not settle a Task.
		{"task_01", map[string]any{"phase": "secret_scan", "task": "task_01", "matches": []map[string]any{{"path": "fixture.txt", "line": 2, "kind": "credential-assignment"}}}},
	} {
		payload, err := json.Marshal(event.payload)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := s.AppendRunEvent(ctx, runevent.RunEvent{RunID: run.ID, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonTask, ReviewIssue: event.workItem, Time: time.Now(), Payload: payload}); err != nil {
			t.Fatal(err)
		}
	}
	workflow := &commandDeliveryWorkflow{store: s, git: preflight.ExecGitRunner{}}
	run.State = store.StateUnresolved
	got, err := workflow.runResult(ctx, root, "0077-example", roundfixCommandResult{exitCode: exitRunFailed}, "", nil, &run)
	if err != nil {
		t.Fatal(err)
	}
	want := delivery.RunResult{
		RunID: run.ID, Outcome: delivery.RunOutcomeUnresolved,
		FailedTasks:  []delivery.TaskFailure{{Task: "task_01", Reason: reason}},
		SkippedTasks: []string{"task_02", "task_03"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Run result = %+v, want %+v", got, want)
	}
}

func TestDeliverStatusNamesTheSecretRefusalOfARunUnresolvedPark(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runStore, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	queue, err := runStore.CreateDeliveryQueue(t.Context(), repo, []string{"0077-example"})
	if err != nil {
		t.Fatal(err)
	}
	const blocker = "run-unresolved: task_01 failed: Task commit carries a secret: fixture.txt:2: credential-assignment; skipped: task_02, task_03"
	item := queue.Items[0]
	item.Stage = store.DeliveryStageParked
	item.Blocker = blocker
	item.Branch = "roundfix/deliver-0077-example"
	item.Worktree = "/worktrees/0077-example"
	if _, _, _, err := runStore.RecordDeliveryQueueItemWorktree(t.Context(), repo, item.SpecSlug, item.Branch, item.Worktree); err != nil {
		t.Fatal(err)
	}
	if err := runStore.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
	const answer = "if the failed Task's text prescribes the value the blocker locates, amend it in /worktrees/0077-example to one no Secret Pattern matches (letters only, or fewer than 20 characters); then run roundfix deliver retry 0077-example, whose rerun returns any Secret Match to the Agent as Verification Feedback"
	want := "0077-example\tparked\t" + blocker + "\t/worktrees/0077-example\n" +
		"Park: 0077-example finding: " + answer + "\n" +
		"Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none\n" +
		"Usage: no Runs recorded\n" +
		"Pending question: 0077-example parked " + blocker + "\n" +
		"Answer: " + answer + "\n"
	if code != exitOK || stderr.Len() != 0 || stdout.String() != want {
		t.Fatalf("exit=%d stderr=%q stdout=%q want=%q", code, stderr.String(), stdout.String(), want)
	}
}
