package cli

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

func TestDeliverStatusReproducesSurfaceTranscriptNoRecordedOwner(t *testing.T) {
	t.Parallel()
	assertRecoveryStatusTranscript(t, "absent")
}
func TestDeliverStatusReproducesSurfaceTranscriptTasksInProgress(t *testing.T) {
	t.Parallel()
	assertRecoveryStatusTranscript(t, "tasks")
}
func TestDeliverStatusReproducesSurfaceTranscriptNoTaskInProgress(t *testing.T) {
	t.Parallel()
	assertRecoveryStatusTranscript(t, "between")
}

func assertRecoveryStatusTranscript(t *testing.T, state string) {
	t.Helper()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	db, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	queue, err := db.CreateDeliveryQueue(t.Context(), repo, []string{implementTestSlug})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage, item.Branch, item.Worktree = store.DeliveryStageRunning, "roundfix/deliver-example", "<item worktree>"
	if _, _, _, err := db.RecordDeliveryQueueItemWorktree(t.Context(), repo, item.SpecSlug, item.Branch, item.Worktree); err != nil {
		t.Fatal(err)
	}
	run, err := db.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: repo, LocalBranch: item.Branch, SpecSlug: item.SpecSlug})
	if err != nil {
		t.Fatal(err)
	}
	item.RunID = run.ID
	if err := db.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	middle := "Owner: Delivery Queue has no recorded owner and 1 item(s) can still advance; run roundfix deliver resume\nRun: " + item.SpecSlug + " no Active Run\n"
	if state == "absent" {
		if _, err := db.CompleteRun(t.Context(), run.ID, "Failed"); err != nil {
			t.Fatal(err)
		}
	} else {
		if err := db.ClaimDeliveryQueueOwner(t.Context(), repo, os.Getpid(), "recorded-owner-identity"); err != nil {
			t.Fatal(err)
		}
		updateCommandDependenciesForTest(t, func(d *commandDependencies) { d.ownerProcesses = &fakeDeliveryOwnerProcesses{} })
		writeTestTurnRecord(t, testTurnRecordPath(t, home, repo, os.Getpid()), deliveryTurnRecord{Schema: deliveryTurnSchema, OwnerPID: os.Getpid(), OwnerIdentity: "recorded-owner-identity", Turn: delivery.Turn{Items: []delivery.ItemTurn{{SpecSlug: item.SpecSlug, Stage: item.Stage, State: delivery.TurnAdvancing}}}})
		events := []string{`{"task":"task_01","phase":"started"}`, `{"task":"task_01","phase":"settled"}`, `{"task":"task_02","phase":"started"}`, `{"task":"task_02","phase":"skipped"}`, `{"task":"task_03","phase":"started"}`, `{"task":"task_04","phase":"started"}`}
		if state == "between" {
			events = append(events, `{"task":"task_03","phase":"settled"}`, `{"task":"task_04","phase":"settled"}`, `{"task":"task_qa","phase":"settled","status":"failed"}`)
		}
		for _, payload := range events {
			if _, err := db.AppendRunEvent(t.Context(), runevent.RunEvent{RunID: run.ID, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonTask, Payload: []byte(payload)}); err != nil {
				t.Fatal(err)
			}
		}
		progress := "Tasks in progress: task_03, task_04"
		if state == "between" {
			progress = "no Task in progress"
		}
		middle = fmt.Sprintf("Turn: %s advancing at running\nRun: %s %s; %s\n", item.SpecSlug, item.SpecSlug, run.ID, progress)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	var out, diag bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &out, &diag)
	want := item.SpecSlug + "\trunning\t-\t<item worktree>\n" + middle + "Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none\nUsage: no prompts recorded across 1 Run(s); cost not reported\n"
	if code != 0 || out.String() != want || diag.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q want=%q", code, out.String(), diag.String(), want)
	}
}

func TestDeliveryRunLinesSelectsTheActiveItemRunAndLatestTaskStarts(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	db, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	create := func(root, slug, branch string, terminal bool) store.Run {
		t.Helper()
		run, err := db.CreateRunSkippingActiveLock(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: root, SpecSlug: slug, LocalBranch: branch})
		if err != nil {
			t.Fatal(err)
		}
		if terminal {
			if _, err := db.CompleteRun(t.Context(), run.ID, store.StateFailed); err != nil {
				t.Fatal(err)
			}
		}
		return run
	}
	create(repo, "first", "feat/first", false)
	active := create(repo, "first", "feat/first", false)
	create(repo, "first", "feat/other", false)
	create(repo, "other", "feat/first", false)
	create(t.TempDir(), "first", "feat/first", false)
	create(repo, "first", "feat/first", true)
	for _, payload := range []string{
		`{"task":"task_03","phase":"started"}`,
		`{"task":"task_04","phase":"started"}`,
		`{"task":"task_03","phase":"settled"}`,
		`{"task":"task_03","phase":"started"}`,
		`{"task":"task_05","phase":"started"}`,
		`{"task":"task_05","phase":"skipped"}`,
	} {
		if _, err := db.AppendRunEvent(t.Context(), runevent.RunEvent{RunID: active.ID, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonTask, Payload: []byte(payload)}); err != nil {
			t.Fatal(err)
		}
	}
	queue := store.DeliveryQueue{Items: []store.DeliveryQueueItem{
		{SpecSlug: "missing", Stage: store.DeliveryStageRunning, Branch: "feat/missing"},
		{SpecSlug: "first", Stage: store.DeliveryStageRunning, Branch: "feat/first"},
		{SpecSlug: "other", Stage: store.DeliveryStageParked},
	}}
	lines, err := deliveryRunLines(t.Context(), db, repo, queue)
	want := []string{"Run: missing no Active Run", "Run: first " + active.ID + "; Tasks in progress: task_04, task_03"}
	if err != nil || !reflect.DeepEqual(lines, want) {
		t.Fatalf("lines=%q err=%v want=%q", lines, err, want)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := deliveryRunLines(t.Context(), db, repo, queue); err == nil {
		t.Fatal("closed Run Database did not fail the status read")
	}
}
