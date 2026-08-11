package tui

import (
	"context"
	"fmt"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"

	tea "charm.land/bubbletea/v2"
)

// Suite: cockpit Task journal refresh cost
// Invariant: one refresh costs the events that arrived since the last one,
// reads payloads only for the kinds whose payload fields the fold parses, and
// renders exactly what a full rescan rendered.
// Boundary IN: synchronous refreshTaskJournalEvents/Update over a recorded
// in-memory journal.
// Boundary OUT: SQL costs of the header projection (owned by the store
// suites), terminal emulation, daemon event production.

// journalReadRecorder counts what one refresh actually pulled out of the
// journal: header rows, whole-event rows, and the payload bytes those whole
// reads carried.
type journalReadRecorder struct {
	cockpitFakeSource
	headerCursors []int64
	headerRows    int
	fullCursors   []int64
	fullRows      int
	payloadBytes  int
}

func (source *journalReadRecorder) RunEventHeadersAfter(ctx context.Context, runID string, cursor int64) ([]store.RunEventHeader, error) {
	headers, err := source.cockpitFakeSource.RunEventHeadersAfter(ctx, runID, cursor)
	source.headerCursors = append(source.headerCursors, cursor)
	source.headerRows += len(headers)
	return headers, err
}

func (source *journalReadRecorder) RunEventsAfter(ctx context.Context, runID string, cursor int64, limit int) ([]store.JournalEvent, error) {
	page, err := source.cockpitFakeSource.RunEventsAfter(ctx, runID, cursor, limit)
	source.fullCursors = append(source.fullCursors, cursor)
	source.fullRows += len(page)
	for _, entry := range page {
		source.payloadBytes += len(entry.Event.Payload)
	}
	return page, err
}

func (source *journalReadRecorder) reset() {
	source.headerCursors = nil
	source.headerRows = 0
	source.fullCursors = nil
	source.fullRows = 0
	source.payloadBytes = 0
}

// newForwardCursorCockpit builds a two-Task spec Run cockpit over a recording
// source, with both Task files left pending on disk so every settled row the
// assertions read comes from the journal.
func newForwardCursorCockpit(t *testing.T, source *journalReadRecorder) *cockpitModel {
	t.Helper()
	gitRoot := t.TempDir()
	slug := "0081-journal-cost"
	tasks := make([]spec.Task, 0, 2)
	for index, title := range []string{"Build scheduler", "Wire queue"} {
		id := fmt.Sprintf("task_%02d", index+1)
		file := writeCockpitTaskFile(t, gitRoot, slug, id, title, spec.StatusPending)
		tasks = append(tasks, spec.Task{ID: id, File: file, Title: title, Status: spec.StatusPending})
	}
	model, err := newCockpitModel(context.Background(), CockpitConfig{
		Mode: CockpitAttach,
		View: LiveRunView{
			Command:       "implement",
			RunKind:       store.KindImplement,
			SpecSlug:      slug,
			GitRoot:       gitRoot,
			RunID:         "run-1",
			PipelineState: source.run.State,
			Concurrency:   2,
			Tasks:         tasks,
			Width:         120,
		},
		RunID:  "run-1",
		Source: source,
	})
	if err != nil {
		t.Fatalf("new cockpit model: %v", err)
	}
	model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	return model
}

func newForwardCursorSource(t *testing.T, backlog int) *journalReadRecorder {
	t.Helper()
	source := &journalReadRecorder{}
	source.run = store.Run{ID: "run-1", State: store.StateResolvingWithAgent}
	source.version = 1
	source.addTaskEvent(t, "task_01", "started", "", 1)
	source.addTaskEvent(t, "task_02", "started", "", 2)
	for index := 1; index <= backlog; index++ {
		source.addLine(fmt.Sprintf("agent line %04d\n", index))
	}
	return source
}

// TestCockpitRefreshCostTracksNewEvents pins the cost contract: two journals
// whose only difference is backlog size must cost the same refresh once the
// same two events arrive. A rescan from cursor zero would make the wide
// journal cost twenty times the narrow one.
func TestCockpitRefreshCostTracksNewEvents(t *testing.T) {
	t.Parallel()
	type refreshCost struct {
		headerRows   int
		fullRows     int
		payloadBytes int
	}
	costs := map[int]refreshCost{}
	renders := map[int]string{}
	for _, backlog := range []int{20, 400} {
		source := newForwardCursorSource(t, backlog)
		model := newForwardCursorCockpit(t, source)

		// The opening fold pays for the whole journal exactly once — that is
		// the replay every later poll is measured against.
		if source.headerRows != backlog+2 {
			t.Fatalf("backlog %d: expected the opening fold to page %d headers, read %d", backlog, backlog+2, source.headerRows)
		}
		source.reset()

		source.addLine("agent line after backlog\n")
		source.addTaskEvent(t, "task_02", "settled", spec.StatusCompleted, 2)
		model.refreshTaskJournalEvents()

		costs[backlog] = refreshCost{
			headerRows:   source.headerRows,
			fullRows:     source.fullRows,
			payloadBytes: source.payloadBytes,
		}

		source.version++
		model.Update(cockpitTickMsg{})
		renders[backlog] = viewText(model)
	}

	if costs[20] != costs[400] {
		t.Fatalf("expected refresh cost to track new events, got %+v for a 20-event backlog and %+v for a 400-event backlog", costs[20], costs[400])
	}
	// Two events arrived: two headers, and the whole read is spent on the one
	// daemon.task event whose payload fields the fold parses.
	if costs[20].headerRows != 2 || costs[20].fullRows != 1 {
		t.Fatalf("expected 2 header rows and 1 whole-event row for two new events, got %+v", costs[20])
	}

	// The cheaper refresh still folds the events it skipped paying for twice:
	// both journals render the same Work Queue rows.
	for _, rendered := range renders {
		assertTaskQueueRow(t, rendered, "task_01", "[run] "+taskLabelAgentWorking)
		assertTaskQueueRow(t, rendered, "task_02", "[done] "+taskLabelCompleted)
	}
}

// TestCockpitTaskJournalRefreshUsesForwardCursorAndHeaderProjection proves the
// mechanism behind the cost: the header cursor advances instead of rewinding,
// payloads load only for the two folded daemon kinds, and folded phases
// survive a refresh that reads nothing new.
func TestCockpitTaskJournalRefreshUsesForwardCursorAndHeaderProjection(t *testing.T) {
	t.Parallel()
	source := newForwardCursorSource(t, 30)
	model := newForwardCursorCockpit(t, source)
	openingCursor := model.taskJournalCursor
	if openingCursor != int64(len(source.events)) {
		t.Fatalf("expected the opening fold to leave the cursor at the journal tail %d, got %d", len(source.events), openingCursor)
	}
	source.reset()

	// Three agent lines and one folded daemon event arrive.
	for index := 1; index <= 3; index++ {
		source.addLine(fmt.Sprintf("later line %d\n", index))
	}
	source.addVerificationEvent(t, "task_01", runevent.VerificationPhaseWaiting, 1, 1)
	model.refreshTaskJournalEvents()

	if len(source.headerCursors) != 1 || source.headerCursors[0] != openingCursor {
		t.Fatalf("expected one header read from the forward cursor %d, got %v", openingCursor, source.headerCursors)
	}
	if source.headerRows != 4 {
		t.Fatalf("expected the refresh to page only the 4 new headers, read %d", source.headerRows)
	}
	// Only the daemon.verification event is read whole; the three agent lines
	// carry payloads the fold never parses and must never be loaded.
	verificationCursor := int64(len(source.events))
	if len(source.fullCursors) != 1 || source.fullCursors[0] != verificationCursor-1 {
		t.Fatalf("expected exactly one whole read, for the folded event at cursor %d, got %v", verificationCursor, source.fullCursors)
	}
	if source.fullRows != 1 {
		t.Fatalf("expected 1 whole-event row, read %d", source.fullRows)
	}
	if model.taskJournalCursor != verificationCursor {
		t.Fatalf("expected the cursor to advance to %d, got %d", verificationCursor, model.taskJournalCursor)
	}

	source.version++
	model.Update(cockpitTickMsg{})
	rendered := viewText(model)
	assertTaskQueueRow(t, rendered, "task_01", "[queued] "+taskLabelWaitingForVerification)
	assertTaskQueueRow(t, rendered, "task_02", "[run] "+taskLabelAgentWorking)

	// A poll with nothing new folds nothing, and the phases folded earlier are
	// still the ones rendered: the refresh accumulates instead of rebuilding.
	source.reset()
	source.version++
	model.Update(cockpitTickMsg{})
	if source.fullRows != 0 {
		t.Fatalf("expected an idle poll to read no whole events, read %d", source.fullRows)
	}
	if quiet := viewText(model); quiet != rendered {
		t.Fatalf("expected an idle poll to leave the render untouched\nbefore:\n%s\n\nafter:\n%s", rendered, quiet)
	}
}

// TestCockpitTaskJournalForwardCursorKeepsSummaryFallback holds the fallback
// the header projection could have quietly replaced: a daemon.task event with
// no payload still settles its row from the summary, and its Task file stays
// pending on disk so the summary is the only source that could have.
func TestCockpitTaskJournalForwardCursorKeepsSummaryFallback(t *testing.T) {
	t.Parallel()
	source := newForwardCursorSource(t, 5)
	model := newForwardCursorCockpit(t, source)
	assertTaskQueueRow(t, viewText(model), "task_01", "[run] "+taskLabelAgentWorking)

	source.addDaemonEvent(runevent.KindDaemonTask, "Task task_01 settled completed.")
	source.version++
	model.Update(cockpitTickMsg{})

	rendered := viewText(model)
	assertTaskQueueRow(t, rendered, "task_01", "[done] "+taskLabelCompleted)
	assertTaskQueueRow(t, rendered, "task_02", "[run] "+taskLabelAgentWorking)

	// The payload-less skip is the same fallback on the other terminal phase.
	source.addDaemonEvent(runevent.KindDaemonTask, "Task task_02 skipped.")
	source.version++
	model.Update(cockpitTickMsg{})
	assertTaskQueueRow(t, viewText(model), "task_02", "[skip] "+taskLabelSkipped)
}
