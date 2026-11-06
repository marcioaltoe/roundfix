// Suite: Delivery Queue Turn Records and status.
// Invariant: status reports turns only for a proven matching owner and never creates artifacts.
// Boundary IN: public deliver status dispatch, real queue storage and local Turn Records.
// Boundary OUT: injected owner identity proof; no external delivery processes.
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/store"
)

func TestDeliverStatusReproducesSurfaceTranscriptQueueTurns(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	slugs := []string{"0066-alpha", "0068-bravo", "0070-charlie", "0074-delta", "0080-echo", "0081-golf", "0083-hotel"}
	stages := []store.DeliveryStage{"reviewing", "running", "reviewing", "reviewing", "queued", "queued", "reviewing"}
	writer, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = writer.Close() })
	queue, err := writer.CreateDeliveryQueue(t.Context(), repo, slugs)
	if err != nil {
		t.Fatal(err)
	}
	for index, item := range queue.Items {
		item.Stage = stages[index]
		if item.Stage != store.DeliveryStageQueued {
			item.Branch = "roundfix/deliver-" + item.SpecSlug
			item.Worktree = "/worktrees/" + item.SpecSlug
			if _, _, _, err := writer.RecordDeliveryQueueItemWorktree(t.Context(), repo, item.SpecSlug, item.Branch, item.Worktree); err != nil {
				t.Fatal(err)
			}
		}
		if err := writer.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.ClaimDeliveryQueueOwner(t.Context(), repo, os.Getpid(), "recorded-owner-identity"); err != nil {
		t.Fatal(err)
	}
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) { dependencies.ownerProcesses = &fakeDeliveryOwnerProcesses{} })
	turn := delivery.Turn{ReviewCapacity: 2, Items: []delivery.ItemTurn{
		{SpecSlug: slugs[0], Stage: stages[0], State: delivery.TurnReviewing},
		{SpecSlug: slugs[1], Stage: stages[1], State: delivery.TurnAdvancing},
		{SpecSlug: slugs[2], Stage: stages[2], State: delivery.TurnReviewing},
		{SpecSlug: slugs[3], Stage: stages[3], State: delivery.TurnReviewSlotWait},
		{SpecSlug: slugs[4], Stage: stages[4], State: delivery.TurnPrerequisiteWait, Prerequisites: []string{"0079-foxtrot"}},
		{SpecSlug: slugs[5], Stage: stages[5], State: delivery.TurnWait},
	}}
	path := testTurnRecordPath(t, home, repo, os.Getpid())
	writeTestTurnRecord(t, path, deliveryTurnRecord{Schema: deliveryTurnSchema, OwnerPID: os.Getpid(), OwnerIdentity: "recorded-owner-identity", Turn: turn})
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
	want := "0066-alpha\treviewing\t-\t/worktrees/0066-alpha\n" +
		"0068-bravo\trunning\t-\t/worktrees/0068-bravo\n" +
		"0070-charlie\treviewing\t-\t/worktrees/0070-charlie\n" +
		"0074-delta\treviewing\t-\t/worktrees/0074-delta\n" +
		"0080-echo\tqueued\t-\t-\n" +
		"0081-golf\tqueued\t-\t-\n" +
		"0083-hotel\treviewing\t-\t/worktrees/0083-hotel\n" +
		"Turn: 0066-alpha reviewing; 2 of 2 review slot(s) held\n" +
		"Turn: 0068-bravo advancing at running\n" +
		"Turn: 0070-charlie reviewing; 2 of 2 review slot(s) held\n" +
		"Turn: 0074-delta waiting for a review slot (1 of 1 waiting); 2 of 2 review slot(s) held\n" +
		"Turn: 0080-echo waiting for prerequisites: 0079-foxtrot\n" +
		"Turn: 0081-golf waiting for the owner's turn (1 of 1 waiting)\n" +
		"Turn: 0083-hotel waiting for the owner to read it\n" +
		"Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none\n" +
		"Usage: no Runs recorded\n"
	if code != exitOK || stdout.String() != want || stderr.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q; want stdout=%q", code, stdout.String(), stderr.String(), want)
	}
}

func TestDeliverStatusPrintsTheChecksAnItemWaitsFor(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	writer, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = writer.Close() })
	queue, err := writer.CreateDeliveryQueue(t.Context(), repo, []string{"0149-example"})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageChecking
	item.Branch, item.Worktree = "roundfix/deliver-0149-example", "/worktrees/0149-example"
	if _, _, _, err := writer.RecordDeliveryQueueItemWorktree(t.Context(), repo, item.SpecSlug, item.Branch, item.Worktree); err != nil {
		t.Fatal(err)
	}
	if err := writer.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	if err := writer.ClaimDeliveryQueueOwner(t.Context(), repo, os.Getpid(), "recorded-owner-identity"); err != nil {
		t.Fatal(err)
	}
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) { dependencies.ownerProcesses = &fakeDeliveryOwnerProcesses{} })
	writeTestTurnRecord(t, testTurnRecordPath(t, home, repo, os.Getpid()), deliveryTurnRecord{
		Schema: deliveryTurnSchema, OwnerPID: os.Getpid(), OwnerIdentity: "recorded-owner-identity",
		Turn: delivery.Turn{ReviewCapacity: 2, Items: []delivery.ItemTurn{{
			SpecSlug: item.SpecSlug, Stage: item.Stage, State: delivery.TurnAdvancing,
			Checks: []string{"CI | Backend / Lint, Type Check, Test & Build", ".github/workflows/ci-frontends.yml (no run yet)"},
		}}},
	})
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
	want := "0149-example\tchecking\t-\t/worktrees/0149-example\n" +
		"Turn: 0149-example advancing at checking; waiting for 2 check(s): CI | Backend / Lint, Type Check, Test & Build; .github/workflows/ci-frontends.yml (no run yet)\n" +
		"Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none\n" +
		"Usage: no Runs recorded\n"
	if code != exitOK || stdout.String() != want || stderr.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q; want stdout=%q", code, stdout.String(), stderr.String(), want)
	}
}

func TestDeliverStatusOrdersWarningsTurnsAndParks(t *testing.T) {
	t.Parallel()
	home, repo := newDeliveryOwnerStatusQueue(t, os.Getpid(), "recorded-owner-identity", store.DeliveryStageChecking, store.DeliveryStageParked)
	writer, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	queue, _, err := writer.DeliveryQueue(t.Context(), repo)
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Warning = "premise-changed: internal/example.go"
	if err := writer.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) { dependencies.ownerProcesses = &fakeDeliveryOwnerProcesses{} })
	writeTestTurnRecord(t, testTurnRecordPath(t, home, repo, os.Getpid()), deliveryTurnRecord{
		Schema: deliveryTurnSchema, OwnerPID: os.Getpid(), OwnerIdentity: "recorded-owner-identity",
		Turn: delivery.Turn{ReviewCapacity: 2, Items: []delivery.ItemTurn{{SpecSlug: item.SpecSlug, Stage: item.Stage, State: delivery.TurnAdvancing}}},
	})
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
	want := "0300-example\tchecking\t-\t-\n0301-example\tparked\tchecks-timeout\t-\n" +
		"Warning: 0300-example premise-changed: internal/example.go\n" +
		"Turn: 0300-example advancing at checking\n" +
		"Park: 0301-example environment: resolve the blocker, then run roundfix deliver retry 0301-example\n" +
		"Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none\n" +
		"Usage: no Runs recorded\n" +
		"Pending question: 0301-example parked checks-timeout\n" +
		"Answer: resolve the blocker, then run roundfix deliver retry 0301-example\n"
	if code != exitOK || stdout.String() != want || stderr.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q want=%q", code, stdout.String(), stderr.String(), want)
	}
}

func TestDeliverStatusIgnoresATurnRecordOfAnotherOwner(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"other identity", "other pid", "owner not running", "no record", "unreadable record", "wrong schema", "unproven identity", "unreadable owner identity"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			pid := os.Getpid()
			if name == "owner not running" {
				pid = 1 << 30
			}
			home, repo := newDeliveryOwnerStatusQueue(t, pid, "recorded-owner-identity", store.DeliveryStageChecking)
			processes := &fakeDeliveryOwnerProcesses{}
			if name == "unproven identity" {
				processes.proveErr = store.ErrOwnerProcessIdentityUnproven
			}
			if name == "unreadable owner identity" {
				processes.proveErr = store.ErrOwnerIdentityUnreadable
			}
			updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) { dependencies.ownerProcesses = processes })
			path := testTurnRecordPath(t, home, repo, pid)
			artifactDir, err := roundconfig.ResolveArtifactDirectory("", repo, home)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(artifactDir); !os.IsNotExist(err) {
				t.Fatalf("artifact dir exists before status: %v", err)
			}
			record := deliveryTurnRecord{Schema: deliveryTurnSchema, OwnerPID: pid, OwnerIdentity: "recorded-owner-identity", Turn: delivery.Turn{ReviewCapacity: 2, Items: []delivery.ItemTurn{{SpecSlug: "0300-example", Stage: store.DeliveryStageChecking, State: delivery.TurnAdvancing}}}}
			switch name {
			case "other identity":
				record.OwnerIdentity = "another-owner-identity"
			case "other pid":
				record.OwnerPID++
			case "wrong schema":
				record.Schema = "roundfix/delivery-turn/v0"
			}
			if name != "no record" {
				writeTestTurnRecord(t, path, record)
			}
			if name == "unreadable record" {
				if err := os.WriteFile(path, []byte("{"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			before := openDeliveryQueueForCLI(t, home, repo)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
			if code != exitOK || strings.Contains(stdout.String(), "Turn:") {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			if after := openDeliveryQueueForCLI(t, home, repo); !reflect.DeepEqual(before, after) {
				t.Fatal("status changed the queue")
			}
			switch name {
			case "unreadable record", "wrong schema":
				if got := stderr.String(); !strings.HasPrefix(got, "roundfix: Delivery Queue turn record "+path+" could not be read: ") || strings.Count(got, "\n") != 1 {
					t.Fatalf("stderr = %q", got)
				}
			case "unreadable owner identity":
				if got, want := stderr.String(), fmt.Sprintf("roundfix: Delivery Queue owner PID %d could not be checked: %v\n", pid, store.ErrOwnerIdentityUnreadable); got != want {
					t.Fatalf("stderr=%q want=%q", got, want)
				}
			default:
				if stderr.Len() != 0 {
					t.Fatalf("stderr=%q", stderr.String())
				}
			}
			if name == "owner not running" {
				want := fmt.Sprintf("Owner: Delivery Queue owner PID %d is not running; run roundfix deliver resume\n", pid)
				if !strings.Contains(stdout.String(), want) {
					t.Fatal(stdout.String())
				}
			}
			if name == "no record" {
				if _, err := os.Stat(artifactDir); !os.IsNotExist(err) {
					t.Fatalf("status created artifact dir: %v", err)
				}
			}
		})
	}
}

func TestTheTurnRecorderWritesBesideTheConsoleLog(t *testing.T) {
	t.Parallel()
	loaded := roundconfig.Loaded{Config: roundconfig.Builtin(), HomeDir: t.TempDir(), GitRoot: t.TempDir()}
	loaded.Config.Defaults.ArtifactDir = filepath.Join(t.TempDir(), "artifacts")
	recorder := &ownerTurnRecorder{loaded: loaded}
	turn := delivery.Turn{ReviewCapacity: 2, Items: []delivery.ItemTurn{
		{SpecSlug: "first", Stage: store.DeliveryStageRunning, State: delivery.TurnAdvancing},
		{SpecSlug: "second", Stage: store.DeliveryStageQueued, State: delivery.TurnPrerequisiteWait, Prerequisites: []string{"first"}},
	}}
	expectedIdentity, err := store.OwnerProcessIdentity(t.Context(), os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(loaded.Config.Defaults.ArtifactDir, "delivery", deliveryOwnerID(loaded.GitRoot, os.Getpid()), "turn.json")
	if got := turnRecordPath(loaded.Config.Defaults.ArtifactDir, loaded.GitRoot, os.Getpid()); got != path {
		t.Fatalf("path = %q want %q", got, path)
	}
	for iteration := range 2 {
		if iteration == 1 {
			turn.ReviewCapacity = 1
			turn.Items = turn.Items[1:]
		}
		if err := recorder.RecordTurn(t.Context(), turn); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var record deliveryTurnRecord
		if err := json.Unmarshal(data, &record); err != nil {
			t.Fatal(err)
		}
		if record.Schema != "roundfix/delivery-turn/v1" || record.OwnerPID != os.Getpid() || record.OwnerIdentity != expectedIdentity || !reflect.DeepEqual(record.Turn, turn) {
			t.Fatalf("record = %+v want turn=%+v identity=%s", record, turn, expectedIdentity)
		}
		var raw struct {
			Schema         string                       `json:"schema"`
			OwnerPID       int                          `json:"owner_pid"`
			OwnerIdentity  string                       `json:"owner_identity"`
			ReviewCapacity int                          `json:"review_capacity"`
			Items          []map[string]json.RawMessage `json:"items"`
		}
		if err := json.Unmarshal(data, &raw); err != nil {
			t.Fatal(err)
		}
		if raw.Schema != deliveryTurnSchema || raw.OwnerPID != os.Getpid() || raw.OwnerIdentity != expectedIdentity || raw.ReviewCapacity != turn.ReviewCapacity || len(raw.Items) != len(turn.Items) {
			t.Fatalf("JSON field names or values changed: %s", data)
		}
		for index, item := range turn.Items {
			for field, want := range map[string]string{"spec": item.SpecSlug, "stage": string(item.Stage), "turn": string(item.State)} {
				var got string
				if err := json.Unmarshal(raw.Items[index][field], &got); err != nil || got != want {
					t.Fatalf("item %d field %s=%q, want %q: %v", index, field, got, want, err)
				}
			}
		}
		if iteration == 0 {
			if _, found := raw.Items[0]["prerequisites"]; found {
				t.Fatalf("empty prerequisites must be omitted: %s", data)
			}
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0o600 {
			t.Fatalf("mode=%o want 600", info.Mode().Perm())
		}
		dirInfo, err := os.Stat(filepath.Dir(path))
		if err != nil {
			t.Fatal(err)
		}
		if dirInfo.Mode().Perm() != 0o700 {
			t.Fatalf("directory mode=%o want 700", dirInfo.Mode().Perm())
		}
		entries, err := os.ReadDir(filepath.Dir(path))
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 1 || entries[0].Name() != "turn.json" {
			t.Fatalf("files = %v", entries)
		}
	}
}

func TestTheDeliveryEngineGetsTheTurnRecorder(t *testing.T) {
	t.Parallel()
	loaded := roundconfig.Loaded{Config: roundconfig.Builtin(), HomeDir: t.TempDir(), GitRoot: t.TempDir()}
	dependencies := commandDeliveryEngineDependencies(nil, loaded)
	recorder, ok := dependencies.TurnRecorder.(*ownerTurnRecorder)
	if !ok || recorder == nil {
		t.Fatalf("TurnRecorder = %T, want *ownerTurnRecorder", dependencies.TurnRecorder)
	}
	if !reflect.DeepEqual(recorder.loaded, loaded) {
		t.Fatal("recorder lost the loaded command configuration")
	}
}

func testTurnRecordPath(t *testing.T, home, repo string, pid int) string {
	t.Helper()
	artifactDir, err := roundconfig.ResolveArtifactDirectory("", repo, home)
	if err != nil {
		t.Fatal(err)
	}
	return turnRecordPath(artifactDir, repo, pid)
}
func writeTestTurnRecord(t *testing.T, path string, record deliveryTurnRecord) {
	t.Helper()
	data, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}
