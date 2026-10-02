// Suite: selected journal reads.
// Invariant: SQL excludes other kinds and Runs while preserving payloads and cursors.
package store

import (
	"bytes"
	"roundfix/internal/runevent"
	"testing"
)

func TestRunEventsOfKindsReturnsOnlyTheNamedKindsInCursorOrder(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	s := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, s)
	req := sampleCreateRunRequest()
	req.GitRoot = t.TempDir()
	run, err := s.CreateRun(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	kinds := []runevent.Kind{runevent.KindAgentMessage, runevent.KindDaemonTask, runevent.KindDaemonVerification, runevent.KindAgentMessage, runevent.KindDaemonTask}
	for i, kind := range kinds {
		e := sampleRunEvent(run.ID, "fixture")
		e.Kind = kind
		e.Payload = []byte(`{ "task": "task_01", "phase": "started" }`)
		if _, err := s.AppendRunEvent(ctx, e); err != nil {
			t.Fatal(err)
		}
		if i == 0 { // an excluded row must never even be scanned
			if _, err := s.db.ExecContext(ctx, `UPDATE run_events SET created_at = 'invalid' WHERE run_id = ? AND cursor = 1`, run.ID); err != nil {
				t.Fatal(err)
			}
		}
	}
	other := req
	other.HeadBranch = "other"
	otherRun, err := s.CreateRun(ctx, other)
	if err != nil {
		t.Fatal(err)
	}
	e := sampleRunEvent(otherRun.ID, "other")
	e.Kind = runevent.KindDaemonTask
	if _, err := s.AppendRunEvent(ctx, e); err != nil {
		t.Fatal(err)
	}
	events, err := s.RunEventsOfKinds(ctx, run.ID, runevent.KindDaemonVerification, runevent.KindDaemonTask, runevent.KindDaemonTask)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 3 {
		t.Fatalf("events=%v", events)
	}
	for i, cursor := range []int64{2, 3, 5} {
		if events[i].Cursor != cursor || events[i].Event.RunID != run.ID || !bytes.Equal(events[i].Event.Payload, []byte(`{ "task": "task_01", "phase": "started" }`)) {
			t.Fatalf("event=%+v", events[i])
		}
	}
	for _, id := range []string{run.ID, "missing"} {
		events, err := s.RunEventsOfKinds(ctx, id)
		if err != nil || len(events) != 0 {
			t.Fatalf("empty selection: %v %v", events, err)
		}
	}
	if _, err := s.RunEventsOfKinds(ctx, ""); err == nil {
		t.Fatal("empty ID accepted")
	}
}
