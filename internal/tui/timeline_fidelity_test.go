package tui

import (
	"strings"
	"testing"
	"time"

	"roundfix/internal/runevent"
	"roundfix/internal/store"

	tea "charm.land/bubbletea/v2"
)

// Suite: timeline fidelity (groups, gutter, bounded summaries, indicator)
// Invariant: the Session Timeline groups rows by Batch with state-driven ▼/▶ markers, aligns one timestamp gutter, renders one bounded summary row per structured event, and carries the Live · detail indicator; the no-color render keeps every distinction via markers.
// Boundary IN: renderedTimelineLines, the cockpit timeline pane renderers, and the token row colorizer.
// Boundary OUT: viewport follow/scroll mechanics and Detail Modal styling, covered by the viewport suite and task_04.
func timelineJournal(events ...runevent.RunEvent) []timelineEntry {
	page := make([]store.JournalEvent, 0, len(events))
	for index, event := range events {
		page = append(page, store.JournalEvent{Cursor: int64(index + 1), Event: event})
	}
	return entriesFromJournal(page)
}

func timelineToolEvent(batch int, at time.Time) runevent.RunEvent {
	return runevent.RunEvent{
		Batch:   batch,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentToolUpdated,
		Summary: "[TOOL] apply_patch · completed\n$ git apply\ndiff --git a/main.go b/main.go\n+added line",
		Time:    at,
		Payload: []byte(`{"sessionId":"s","update":{"sessionUpdate":"tool_call_update","toolCallId":"call_9","title":"apply_patch","status":"completed","rawInput":{"command":"git apply"},"content":[{"content":{"type":"text","text":"diff --git a/main.go b/main.go\n+added line"}}]}}`),
	}
}

func TestTimelineSettledBatchCollapsesAndExecutingBatchExpands(t *testing.T) {
	startedAt := time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
	entries := timelineJournal(
		runevent.RunEvent{
			Batch: 1, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonBatch,
			Summary: "Batch 001/002 started with 2 Review Issue(s).", Time: startedAt,
		},
		timelineToolEvent(1, startedAt.Add(10*time.Second)),
		runevent.RunEvent{
			Batch: 1, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonBatch,
			Summary: "Batch 001 completed; 1 Unresolved Review Issue(s) remain.", Time: startedAt.Add(38 * time.Second),
		},
		runevent.RunEvent{
			Batch: 2, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonBatch,
			Summary: "Batch 002/002 started with 1 Review Issue(s).", Time: startedAt.Add(40 * time.Second),
		},
		timelineToolEvent(2, startedAt.Add(50*time.Second)),
	)

	lines := renderedTimelineLines(entries)
	if len(lines) != 3 {
		t.Fatalf("expected collapsed + expanded groups to render 3 rows, got %d:\n%s", len(lines), strings.Join(lines, "\n"))
	}
	if lines[0] != "▶ BATCH 001/002 completed 00:38" {
		t.Fatalf("expected the settled Batch collapsed to its summary row, got %q", lines[0])
	}
	if lines[1] != "▼ BATCH 002/002 started 00:10" {
		t.Fatalf("expected the executing Batch expanded with ▼, got %q", lines[1])
	}
	if lines[2] != "12:00:50 [TOOL] apply_patch · completed" {
		t.Fatalf("expected the executing Batch to render its event rows, got %q", lines[2])
	}
	joined := strings.Join(lines, "\n")
	for _, hidden := range []string{"$ git apply", "diff --git", "+added line"} {
		if strings.Contains(joined, hidden) {
			t.Fatalf("expected tool payload %q to stay behind the summary row, got:\n%s", hidden, joined)
		}
	}
}

func TestTimelineGutterAlignsAcrossKinds(t *testing.T) {
	startedAt := time.Date(2026, 7, 5, 12, 0, 5, 0, time.UTC)
	entries := timelineJournal(
		runevent.RunEvent{
			Source: runevent.SourceAgent, Kind: runevent.KindAgentPlan,
			Summary: "PLAN\npending  Inspect current cockpit render",
			Time:    startedAt,
			Payload: []byte(`{"sessionId":"s","update":{"sessionUpdate":"plan","entries":[{"status":"pending","content":"Inspect current cockpit render"}]}}`),
		},
		timelineToolEvent(0, startedAt.Add(15*time.Second)),
		runevent.RunEvent{
			Source: runevent.SourceDaemon, Kind: runevent.KindDaemonTask,
			Summary: "Task task_01 settled completed.",
		},
	)

	lines := renderedTimelineLines(entries)
	want := []string{
		"12:00:05 PLAN",
		"12:00:20 [TOOL] apply_patch · completed",
		"         TASK Task task_01 settled completed.",
	}
	if len(lines) != len(want) {
		t.Fatalf("expected %d gutter rows, got %d:\n%s", len(want), len(lines), strings.Join(lines, "\n"))
	}
	for index, row := range lines {
		if row != want[index] {
			t.Fatalf("gutter row %d = %q, want %q", index, row, want[index])
		}
		gutter, _, ok := splitTimelineGutter(row)
		if !ok || len([]rune(gutter)) != timelineGutterWidth {
			t.Fatalf("expected row %q to carry the aligned %d-column gutter", row, timelineGutterWidth)
		}
	}
}

func TestTimelineSelectionProjectionErrorFallsBackToPersistedSummary(t *testing.T) {
	startedAt := time.Date(2026, 7, 5, 12, 0, 5, 0, time.UTC)
	entries := timelineJournal(runevent.RunEvent{
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonAgentSelectionActive,
		Summary: "Agent Selection active from legacy summary.",
		Time:    startedAt,
		Payload: []byte(`{"attempt":1}`),
	})

	lines := renderedTimelineLines(entries)

	want := "12:00:05 SELECTION Agent Selection active from legacy summary."
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("timeline lines = %#v, want %q", lines, want)
	}
}

func TestTimelineRowsStyledThroughTokensAndNoColorTwin(t *testing.T) {
	tokens := ResolveTokens(true)
	styled := &cockpitModel{tokens: tokens}
	plain := &cockpitModel{tokens: ResolveTokens(false)}
	rows := []string{
		"▶ BATCH 001/002 completed 00:38",
		"▼ BATCH 002/002 started 00:10",
		"12:00:50 [TOOL] apply_patch · completed",
		"12:00:32 THINK checking error paths",
		"         TASK Task task_01 settled completed.",
	}

	got := styled.styleTimelineRows(rows, 60)
	want := []string{
		tokens.SectionLabel.Render("▶ BATCH 001/002") + " " + tokens.Done.Render("completed") + " " + tokens.Muted.Render("00:38"),
		tokens.SectionLabel.Render("▼ BATCH 002/002") + " " + tokens.Running.Render("started") + " " + tokens.Muted.Render("00:10"),
		tokens.Muted.Render("12:00:50 ") + tokens.SectionLabel.Render("[TOOL]") + " apply_patch · completed",
		tokens.Muted.Render("12:00:32 ") + tokens.Muted.Render("THINK checking error paths"),
		tokens.Muted.Render("         ") + tokens.SectionLabel.Render("TASK") + " Task task_01 settled completed.",
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("styled timeline row %d mismatch\ngot:  %q\nwant: %q", index, got[index], want[index])
		}
	}

	for index, row := range plain.styleTimelineRows(rows, 60) {
		if strings.Contains(row, "\x1b") {
			t.Fatalf("no-color timeline row carries ANSI: %q", row)
		}
		if row != rows[index] {
			t.Fatalf("no-color timeline row %d = %q, want marker-only %q", index, row, rows[index])
		}
	}
}

func TestTimelineToolPayloadRendersOneBoundedLineInThePane(t *testing.T) {
	startedAt := time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
	events := []runevent.RunEvent{
		{
			Batch: 1, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonBatch,
			Summary: "Batch 001/002 started with 2 Review Issue(s).", Time: startedAt,
		},
		timelineToolEvent(1, startedAt.Add(50*time.Second)),
	}
	model := newQueueFidelityCockpit(t, true, store.StateResolvingWithAgent, events, startedAt.Add(time.Minute))
	model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	rendered := viewText(model)
	if !strings.Contains(rendered, "12:00:50 [TOOL] apply_patch · completed") {
		t.Fatalf("expected the tool event's single summary row, got:\n%s", rendered)
	}
	for _, hidden := range []string{"$ git apply", "diff --git", "+added line"} {
		if strings.Contains(rendered, hidden) {
			t.Fatalf("expected raw tool payload %q to never render inline, got:\n%s", hidden, rendered)
		}
	}
	for _, line := range strings.Split(rendered, "\n") {
		if displayWidth(line) > 120 {
			t.Fatalf("expected every pane row bounded to the pane width, got %d columns: %q", displayWidth(line), line)
		}
	}
}

func TestTimelinePaneHeaderIndicatorFollowsModalState(t *testing.T) {
	model := newQueueFidelityCockpit(t, true, store.StateResolvingWithAgent, nil, time.Now())
	model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	if rendered := viewText(model); !strings.Contains(rendered, "Live · detail hidden") {
		t.Fatalf("expected the pane header to carry the detail-hidden indicator, got:\n%s", rendered)
	}
	pressKey(t, model, "tab")
	pressKey(t, model, "enter")
	rendered := viewText(model)
	if !strings.Contains(rendered, "Live · detail open") {
		t.Fatalf("expected the indicator to flip to detail open with the modal, got:\n%s", rendered)
	}
	if strings.Contains(rendered, "Live · detail hidden") {
		t.Fatalf("expected the hidden indicator to disappear while the modal is open, got:\n%s", rendered)
	}
	pressKey(t, model, "esc")
	if rendered := viewText(model); !strings.Contains(rendered, "Live · detail hidden") {
		t.Fatalf("expected the indicator to return to detail hidden after Esc, got:\n%s", rendered)
	}
}
