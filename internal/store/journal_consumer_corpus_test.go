package store

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"roundfix/internal/runevent"
)

// Suite: consumer corpus replay
// Invariant: a journal recorded before the header-projection change replays
// identically through every consumer surface, and the new payload-free header
// projection is a faithful subset of the unchanged full read (ADR 0008).
// Boundary IN: Store appends, the full read, the header projection, and the
// runevent stream consumer.
// Boundary OUT: the stream schema, command output shapes, and any payload
// rewrite.

// consumerCorpusEvent seeds one run with a realistic pre-change vocabulary of
// event kinds and payloads across Batch boundaries.
func consumerCorpusEvent(runID string, index int, cursor int64) runevent.RunEvent {
	base := time.Date(2026, 8, 11, 10, 0, 0, 0, time.UTC).Add(time.Duration(index) * time.Second)
	switch index % 5 {
	case 0:
		return runevent.RunEvent{
			RunID: runID, Batch: 1 + index%2, Source: runevent.SourceAgent,
			Kind: runevent.KindAgentMessage, Summary: "agent chunk",
			Time: base, Payload: []byte(`{"sessionId":"s-1","update":{"sessionUpdate":"agent_message_chunk","content":[{"type":"text","text":"chunk"}]}}`),
		}
	case 1:
		return runevent.RunEvent{
			RunID: runID, Batch: 1 + index%2, Source: runevent.SourceDaemon,
			Kind: runevent.KindDaemonTask, ReviewIssue: "task_02",
			Summary: "Task task_02 started as Batch 001.",
			Time:    base, Payload: []byte(`{"task":"task_02","phase":"started","batch":1}`),
		}
	case 2:
		return runevent.RunEvent{
			RunID: runID, Batch: 1 + index%2, Source: runevent.SourceDaemon,
			Kind: runevent.KindDaemonVerification, ReviewIssue: "task_02",
			Summary: "Verification attempt 1 for Task task_02.",
			Time:    base, Payload: []byte(`{"attempt":1,"phase":"verdict","task":"task_02","verdict":"passed"}`),
		}
	case 3:
		return runevent.RunEvent{
			RunID: runID, Batch: 1 + index%2, Source: runevent.SourceDaemon,
			Kind:    runevent.KindDaemonStatus,
			Summary: "Daemon status.",
			Time:    base, Payload: []byte(`{"phase":"reviewing","task_capacity":2,"verification_capacity":1}`),
		}
	default:
		return runevent.RunEvent{
			RunID: runID, Batch: 1 + index%2, Source: runevent.SourceDaemon,
			Kind:    runevent.KindDaemonOutcome,
			Summary: "Run reached Clean.",
			Time:    base, Payload: []byte(`{"outcome":"clean","state":"clean"}`),
		}
	}
}

// consumerCorpusTotal is the shared fixed-size corpus every consumer replay
// seeds and replays.
const consumerCorpusTotal = 12

func assertCorpus(t *testing.T, events []JournalEvent, wantCursor int64) {
	t.Helper()
	if len(events) == 0 {
		t.Fatal("corpus replay: expected recorded events")
	}
	if int64(len(events)) != wantCursor {
		t.Fatalf("corpus replay: recorded %d events, want %d", len(events), wantCursor)
	}
	for index, entry := range events {
		if entry.Cursor != int64(index+1) {
			t.Fatalf("corpus replay: event %d has cursor %d, want %d", index, entry.Cursor, index+1)
		}
	}
}

// seedConsumerCorpus opens a store, creates a run, and records the shared
// pre-change corpus of consumerCorpusTotal events.
func seedConsumerCorpus(t *testing.T, ctx context.Context) (*Store, string) {
	t.Helper()
	s := openTestStore(t, ctx, t.TempDir())
	t.Cleanup(func() { closeStore(t, s) })
	run, err := s.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	for index := range consumerCorpusTotal {
		if _, err := s.AppendRunEvent(ctx, consumerCorpusEvent(run.ID, index, int64(index+1))); err != nil {
			t.Fatalf("append corpus event %d: %v", index, err)
		}
	}
	return s, run.ID
}

func TestConsumerCorpusFullReadReplaysIdentically(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s, runID := seedConsumerCorpus(t, ctx)

	recorded, err := s.RunEventsAfter(ctx, runID, 0, consumerCorpusTotal)
	if err != nil {
		t.Fatalf("record corpus via full read: %v", err)
	}
	assertCorpus(t, recorded, consumerCorpusTotal)

	paged := []JournalEvent{}
	cursor := int64(0)
	for {
		page, err := s.RunEventsAfter(ctx, runID, cursor, 5)
		if err != nil {
			t.Fatalf("replay corpus page at %d: %v", cursor, err)
		}
		for _, entry := range page {
			if entry.Cursor > cursor {
				cursor = entry.Cursor
			}
			paged = append(paged, entry)
		}
		if len(page) < 5 {
			break
		}
	}
	if len(paged) != len(recorded) {
		t.Fatalf("paged replay count %d != recorded %d", len(paged), len(recorded))
	}
	for index := range recorded {
		if paged[index].Cursor != recorded[index].Cursor {
			t.Fatalf("paged cursor %d != recorded %d", paged[index].Cursor, recorded[index].Cursor)
		}
		if paged[index].Event.Summary != recorded[index].Event.Summary {
			t.Fatalf("paged summary %q != recorded %q", paged[index].Event.Summary, recorded[index].Event.Summary)
		}
		if string(paged[index].Event.Payload) != string(recorded[index].Event.Payload) {
			t.Fatalf("paged payload byte mismatch at cursor %d", recorded[index].Cursor)
		}
	}
}

// TestConsumerCorpusEventsStreamReplaysIdentically proves the `events` stream
// consumer — which reads daemon payload fields and hard-fails on a malformed
// selected daemon payload — reproduces the same stream records from the full
// read before and after the header projection ships beside it.
func TestConsumerCorpusEventsStreamReplaysIdentically(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s, runID := seedConsumerCorpus(t, ctx)

	events, err := s.RunEventsAfter(ctx, runID, 0, consumerCorpusTotal)
	if err != nil {
		t.Fatalf("load corpus: %v", err)
	}

	// The fixed corpus is deterministic, so the projected envelope lines are
	// stable. Only daemon kinds project a record; agent chunks and Daemon
	// status events are skipped by the stream consumer.
	wantLines := []string{
		"task-status|started||||task_02",
		"verification|verdict||passed||task_02",
		"outcome||||clean|",
		"task-status|started||||task_02",
		"verification|verdict||passed||task_02",
		"outcome||||clean|",
		"task-status|started||||task_02",
	}
	lines := []string{}
	for _, entry := range events {
		record, ok, err := runevent.ProjectStreamEvent(entry.Cursor, entry.Event, runevent.AllStreamCategories())
		if err != nil {
			t.Fatalf("project stream event %d: %v", entry.Cursor, err)
		}
		if !ok {
			continue
		}
		if record.Schema != runevent.StreamSchema {
			t.Fatalf("stream schema changed: got %q", record.Schema)
		}
		// Stable, order-independent proof the consumer output is unchanged:
		// the envelope fields that are independent of payload content.
		lines = append(lines, strings.Join([]string{
			string(record.Category),
			record.Phase,
			record.Status,
			record.Verdict,
			record.Outcome,
			record.WorkItem,
		}, "|"))
	}
	if len(lines) != len(wantLines) {
		t.Fatalf("events stream produced %d records, want %d", len(lines), len(wantLines))
	}
	for index := range wantLines {
		if lines[index] != wantLines[index] {
			t.Fatalf("events stream record %d = %q, want %q", index, lines[index], wantLines[index])
		}
	}
	// The corpus must exercise every projected category so the replay is not a
	// rehearsal of a single category.
	seen := map[string]bool{}
	for _, line := range lines {
		seen[strings.Split(line, "|")[0]] = true
	}
	for _, category := range []string{
		string(runevent.StreamCategoryTaskStatus),
		string(runevent.StreamCategoryVerification),
		string(runevent.StreamCategoryOutcome),
	} {
		if !seen[category] {
			t.Fatalf("events stream corpus did not exercise category %q (degraded replay)", category)
		}
	}
}

// ReplayCorpusHeaderMatchesFullRead proves the header projection is the
// payload-free subset of the unchanged full read for the same recorded corpus:
// every consumer that needs payload still gets it from the full read, while a
// header-only consumer reads exactly the columns it uses.
func TestReplayCorpusHeaderMatchesFullRead(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s, runID := seedConsumerCorpus(t, ctx)

	events, err := s.RunEventsAfter(ctx, runID, 0, consumerCorpusTotal)
	if err != nil {
		t.Fatalf("load corpus full read: %v", err)
	}
	headers, err := s.RunEventHeadersAfter(ctx, runID, 0, consumerCorpusTotal)
	if err != nil {
		t.Fatalf("load corpus header projection: %v", err)
	}
	if len(headers) != len(events) {
		t.Fatalf("header count %d != full count %d", len(headers), len(events))
	}
	for index := range events {
		header, full := headers[index], events[index]
		if header.Cursor != full.Cursor || header.Batch != full.Event.Batch ||
			header.Source != full.Event.Source || header.Kind != full.Event.Kind ||
			header.Summary != full.Event.Summary || !header.Time.Equal(full.Event.Time) {
			t.Fatalf("header projection diverges at cursor %d: %+v vs %+v", full.Cursor, header, full.Event)
		}
		if len(full.Event.Payload) == 0 {
			t.Fatalf("corpus full read lost a payload at cursor %d", full.Cursor)
		}
	}
}

// ReplayCorpusBatchClockMatchesFullEvents proves the moved consumer — the
// cockpit batch-clock refresh, which reads only Batch and Time — computes the
// same per-Batch spans from the header projection and from the full read.
func TestReplayCorpusBatchClockMatchesFullEvents(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s, runID := seedConsumerCorpus(t, ctx)

	events, err := s.RunEventsAfter(ctx, runID, 0, consumerCorpusTotal)
	if err != nil {
		t.Fatalf("load corpus full read: %v", err)
	}
	headers, err := s.RunEventHeadersAfter(ctx, runID, 0, consumerCorpusTotal)
	if err != nil {
		t.Fatalf("load corpus header projection: %v", err)
	}

	type span struct{ first, last time.Time }
	// foldSpans collapses per-Batch first/last times for both the full events
	// and the header projection with one shared lookup shape.
	foldSpans := func(count int, at func(index int) (batch int, when time.Time)) map[int]span {
		spans := map[int]span{}
		for index := range count {
			batch, when := at(index)
			if batch <= 0 || when.IsZero() {
				continue
			}
			current, seen := spans[batch]
			if !seen || when.Before(current.first) {
				current.first = when
			}
			if when.After(current.last) {
				current.last = when
			}
			spans[batch] = current
		}
		return spans
	}
	fullSpans := foldSpans(len(events), func(index int) (int, time.Time) {
		return events[index].Event.Batch, events[index].Event.Time
	})
	headerSpans := foldSpans(len(headers), func(index int) (int, time.Time) {
		return headers[index].Batch, headers[index].Time
	})
	if len(fullSpans) != len(headerSpans) {
		t.Fatalf("batch-clock span count from full %d != header %d", len(fullSpans), len(headerSpans))
	}
	for batch, want := range fullSpans {
		got, ok := headerSpans[batch]
		if !ok {
			t.Fatalf("batch %d missing from header-derived spans", batch)
		}
		if !want.first.Equal(got.first) || !want.last.Equal(got.last) {
			t.Fatalf("batch %d span from headers %+v != full %+v", batch, got, want)
		}
	}
}

// TestJournalConsumerCorpusReplaysEveryConsumer replays one SQLite journal
// written by the pre-Spec build through the current events, Attach, Cockpit,
// reconcile, and GC consumers. The package-specific harnesses are installed
// with Go overlays, so the characterization crosses the real consumer seams
// without adding test hooks to production code or rewriting the fixture.
func TestJournalConsumerCorpusReplaysEveryConsumer(t *testing.T) {
	if testing.Short() {
		t.Skip("consumer corpus replay runs a nested go test; run without -short")
	}
	repositoryRoot := journalConsumerCorpusRepositoryRoot(t)
	testdataDir := filepath.Join(repositoryRoot, "internal", "store", "testdata")
	fixture := filepath.Join(testdataDir, "2026-08-11-prechange-roundfix.db")

	tests := []struct {
		name        string
		packagePath string
		virtualFile string
		harness     string
		testName    string
		environment []string
	}{
		{
			name:        "events Attach reconcile and gc preserve pre-change observations",
			packagePath: "./internal/cli",
			virtualFile: filepath.Join(repositoryRoot, "internal", "cli", "task10_consumer_observation_test.go"),
			harness:     filepath.Join(testdataDir, "task10-cli-consumer-harness_test.go.txt"),
			testName:    "TestTask10PrechangeJournalCLIConsumers",
			environment: []string{
				"TASK10_CORPUS_DB=" + fixture,
				"TASK10_EXPECTATIONS=" + filepath.Join(testdataDir, "2026-08-11-prechange-consumer-expectations.json"),
				"TASK10_RECORD=0",
			},
		},
		{
			name:        "Cockpit rendering preserves the pre-change frame",
			packagePath: "./internal/tui",
			virtualFile: filepath.Join(repositoryRoot, "internal", "tui", "task10_cockpit_observation_test.go"),
			harness:     filepath.Join(testdataDir, "task10-cockpit-consumer-harness_test.go.txt"),
			testName:    "TestTask10PrechangeJournalCockpitConsumer",
			environment: []string{
				"TASK10_CORPUS_DB=" + fixture,
				"TASK10_COCKPIT_EXPECTATION=" + filepath.Join(testdataDir, "2026-08-11-prechange-cockpit.golden"),
				"TASK10_RECORD=0",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			overlay := struct {
				Replace map[string]string `json:"Replace"`
			}{Replace: map[string]string{test.virtualFile: test.harness}}
			overlayData, err := json.Marshal(overlay)
			if err != nil {
				t.Fatalf("encode %s overlay: %v", test.name, err)
			}
			overlayPath := filepath.Join(t.TempDir(), "overlay.json")
			if err := os.WriteFile(overlayPath, overlayData, 0o644); err != nil {
				t.Fatalf("write %s overlay: %v", test.name, err)
			}

			command := exec.CommandContext(
				t.Context(),
				"go", "test",
				"-overlay="+overlayPath,
				"-count=1",
				test.packagePath,
				"-run", "^"+test.testName+"$",
				"-v",
			)
			command.Dir = repositoryRoot
			command.Env = append(os.Environ(), test.environment...)
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("replay %s: %v\n%s", test.name, err, output)
			}
			t.Logf("%s", output)
		})
	}
}

func journalConsumerCorpusRepositoryRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve journal consumer corpus source path")
	}
	root, err := filepath.Abs(filepath.Join(filepath.Dir(filename), "..", ".."))
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}
	return root
}
