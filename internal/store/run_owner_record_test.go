package store

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/runevent"
)

func TestReleaseStaleRunOwnerRecordClearsOnlyTheRecordedTerminalOwner(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	db := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, db)
	request := sampleCreateRunRequest()
	request.GitRoot = t.TempDir()
	request.ArtifactDir = t.TempDir()
	request.OwnerPID = 4242
	request.OwnerIdentity = "test:old"
	run, err := db.CreateRun(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	assertKept := func() {
		t.Helper()
		got, found, err := db.Run(ctx, run.ID)
		if err != nil || !found || !reflect.DeepEqual(got, run) {
			t.Fatalf("Run changed: %+v, %v", got, err)
		}
		events, err := db.RunEventsAfter(ctx, run.ID, 0, 100)
		if err != nil || len(events) != 0 {
			t.Fatalf("unexpected journal: %+v, %v", events, err)
		}
	}
	for _, tt := range []struct {
		name, id, identity, reason string
		pid                        int
	}{
		{name: "empty ID", identity: "test:old", reason: "proof", pid: 4242},
		{name: "empty identity", id: run.ID, reason: "proof", pid: 4242},
		{name: "empty reason", id: run.ID, identity: "test:old", pid: 4242},
		{name: "invalid PID", id: run.ID, identity: "test:old", reason: "proof"},
		{name: "missing Run", id: "absent", identity: "test:old", reason: "proof", pid: 4242},
		{name: "active Run", id: run.ID, identity: "test:old", reason: "proof", pid: 4242},
	} {
		t.Run(tt.name, func(t *testing.T) {
			released, err := db.ReleaseStaleRunOwnerRecord(ctx, tt.id, tt.pid, tt.identity, tt.reason)
			if err == nil || released {
				t.Fatalf("release = %v, %v", released, err)
			}
			if tt.name == "active Run" && err.Error() != fmt.Sprintf("release stale Run owner record of Run %q: Run is %s, not terminal", run.ID, run.State) {
				t.Fatal(err)
			}
			assertKept()
		})
	}
	completed, err := db.CompleteRun(ctx, run.ID, StateFailed)
	if err != nil {
		t.Fatal(err)
	}
	run = completed.Run
	for _, tt := range []struct {
		name     string
		pid      int
		identity string
	}{
		{"different PID", 4243, "test:old"}, {"different identity", 4242, "test:new"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			released, err := db.ReleaseStaleRunOwnerRecord(ctx, run.ID, tt.pid, tt.identity, "proof")
			if err != nil || released {
				t.Fatalf("release = %v, %v", released, err)
			}
			assertKept()
		})
	}
	t.Run("journal failure rolls back owner release", func(t *testing.T) {
		if _, err := db.db.ExecContext(ctx, `CREATE TRIGGER reject_owner_record_event BEFORE INSERT ON run_events BEGIN SELECT RAISE(FAIL, 'journal unavailable'); END`); err != nil {
			t.Fatal(err)
		}
		released, err := db.ReleaseStaleRunOwnerRecord(ctx, run.ID, 4242, "test:old", "proof")
		if released || err == nil || !strings.Contains(err.Error(), "journal stale Run owner record release") {
			t.Fatalf("release = %v, %v", released, err)
		}
		assertKept()
		if _, err := db.db.ExecContext(ctx, `DROP TRIGGER reject_owner_record_event`); err != nil {
			t.Fatal(err)
		}
	})
	released, err := db.ReleaseStaleRunOwnerRecord(ctx, " "+run.ID+" ", 4242, " test:old ", " proof ")
	if err != nil || !released {
		t.Fatalf("release = %v, %v", released, err)
	}
	got, found, err := db.Run(ctx, run.ID)
	if err != nil || !found {
		t.Fatalf("read released Run: %v", err)
	}
	want := run
	want.OwnerPID = nil
	want.OwnerIdentity = ""
	want.OwnerIdentityUnproven = false
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("released Run = %+v, want %+v", got, want)
	}
	events, err := db.RunEventsAfter(ctx, run.ID, 0, 100)
	if err != nil || len(events) != 1 {
		t.Fatalf("journal = %+v, %v", events, err)
	}
	event := events[0].Event
	var payload map[string]any
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		t.Fatal(err)
	}
	if event.Source != runevent.SourceDaemon || event.Kind != runevent.KindDaemonOutcome || event.Summary != "Run released its Stale Run Owner Record: proof." ||
		!reflect.DeepEqual(payload, map[string]any{"event": "stale_run_owner_record_release", "owner_pid": float64(4242), "owner_identity": "test:old", "reason": "proof"}) {
		t.Fatalf("release event = %+v, payload = %+v", event, payload)
	}
	released, err = db.ReleaseStaleRunOwnerRecord(ctx, run.ID, 4242, "test:old", "proof")
	if err != nil || released {
		t.Fatalf("second release = %v, %v", released, err)
	}
	again, err := db.RunEventsAfter(ctx, run.ID, 0, 100)
	if err != nil || !reflect.DeepEqual(again, events) {
		t.Fatalf("second release changed journal: %+v, %v", again, err)
	}
	if strings.TrimSpace(got.State) != StateFailed {
		t.Fatal("release changed outcome")
	}
}
