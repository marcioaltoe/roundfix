package tui

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/runevent"
	"roundfix/internal/store"

	"charm.land/lipgloss/v2"
)

// Suite: cockpit queue fidelity
// Invariant: the header, Phase Row, and Work Queue render through the documented state tokens under a forced color profile, and the same fixtures degrade to marker-only text without color.
// Boundary IN: synchronous cockpit surface renderers with forced token sets.
// Boundary OUT: timeline grouping and Detail Modal styling, covered by the later fidelity tasks.
func newQueueFidelityCockpit(t *testing.T, colorEnabled bool, runState string, events []runevent.RunEvent, clock time.Time) *cockpitModel {
	t.Helper()
	source := &cockpitFakeSource{run: store.Run{ID: "run-review-cafe1234", State: runState}, version: 1}
	for _, event := range events {
		source.addEvent(event)
	}
	model, err := newCockpitModel(context.Background(), CockpitConfig{
		Mode: CockpitAttach,
		View: LiveRunView{
			PRNumber:      "7",
			RunID:         "run-review-cafe1234",
			PipelineState: runState,
			Width:         100,
			Issues:        sampleIssues(3),
			BatchSizes:    []int{2, 1},
		},
		RunID:        "run-review-cafe1234",
		Source:       source,
		ColorEnabled: colorEnabled,
		Now:          func() time.Time { return clock },
	})
	if err != nil {
		t.Fatalf("new cockpit model: %v", err)
	}
	return model
}

func TestCockpitHeaderRendersSectionLabelAndStateChipTokens(t *testing.T) {
	t.Parallel()
	tokens := ResolveTokens(true)
	model := newQueueFidelityCockpit(t, true, store.StateResolvingWithAgent, nil, time.Now())

	got := renderCockpitHeaderLine(model, 100)
	left := tokens.SectionLabel.Render("ROUNDFIX") + tokens.Muted.Render(" // LIVE RUN VIEW // PR #7")
	right := tokens.Muted.Render("RUN cafe1234") + " " + tokens.Running.Render("[RESOLVING WITH AGENT]")
	want := padRightDisplay(left, maxInt(100-displayWidth("RUN cafe1234 [RESOLVING WITH AGENT]"), 1)) + right
	if got != want {
		t.Fatalf("styled header mismatch\ngot:  %q\nwant: %q", got, want)
	}

	chips := []struct {
		name  string
		state string
		token lipgloss.Style
	}{
		{"clean chip green", store.StateClean, tokens.Done},
		{"failed chip red", store.StateFailed, tokens.Failed},
		{"stopped chip red", store.StateStopped, tokens.Failed},
	}
	for _, chip := range chips {
		t.Run(chip.name, func(t *testing.T) {
			model := newQueueFidelityCockpit(t, true, chip.state, nil, time.Now())
			rendered := renderCockpitHeaderLine(model, 120)
			wantChip := chip.token.Render("[" + formatRunStateLabel(chip.state) + "]")
			if !strings.Contains(rendered, wantChip) {
				t.Fatalf("expected header %q to carry state chip %q", rendered, wantChip)
			}
		})
	}
}

func TestCockpitHeaderNoColorRendersMarkerOnlyText(t *testing.T) {
	t.Parallel()
	styled := renderCockpitHeaderLine(newQueueFidelityCockpit(t, true, store.StateResolvingWithAgent, nil, time.Now()), 100)
	plain := renderCockpitHeaderLine(newQueueFidelityCockpit(t, false, store.StateResolvingWithAgent, nil, time.Now()), 100)
	if strings.Contains(plain, "\x1b") {
		t.Fatalf("no-color header carries ANSI: %q", plain)
	}
	if plain != stripANSI(styled) {
		t.Fatalf("no-color header text differs from styled text\nplain:  %q\nstyled: %q", plain, stripANSI(styled))
	}
	for _, marker := range []string{"ROUNDFIX // LIVE RUN VIEW // PR #7", "RUN cafe1234", "[RESOLVING WITH AGENT]"} {
		if !strings.Contains(plain, marker) {
			t.Fatalf("expected no-color header to keep %q, got %q", marker, plain)
		}
	}
}

func TestCockpitPhaseRowColorsMarkersThroughTokens(t *testing.T) {
	t.Parallel()
	tokens := ResolveTokens(true)
	identity := ResolveTokens(false)
	tests := []struct {
		name  string
		phase cockpitPhase
		token lipgloss.Style
	}{
		{"done green", cockpitPhase{name: "FETCH", marker: phaseDone}, tokens.Done},
		{"run amber", cockpitPhase{name: "AGENT", marker: phaseRun}, tokens.Running},
		{"wait amber", cockpitPhase{name: "VERIFY", marker: phaseWait}, tokens.Waiting},
		{"locked red", cockpitPhase{name: "PUSH", marker: phaseLocked}, tokens.Locked},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			styled := renderPhase(tokens, tt.phase)
			want := tt.phase.name + " " + tt.token.Render("["+tt.phase.marker+"]")
			if styled != want {
				t.Fatalf("styled phase mismatch\ngot:  %q\nwant: %q", styled, want)
			}
			plainWant := tt.phase.name + " [" + tt.phase.marker + "]"
			if plain := renderPhase(identity, tt.phase); plain != plainWant {
				t.Fatalf("no-color phase = %q, want marker-only %q", plain, plainWant)
			}
		})
	}
}

func queueFidelityCard() WorkItem {
	return WorkItem{
		Name:     "Issue #001",
		Title:    "Guard nil cache",
		Severity: "major",
		Ordinal:  1,
		Location: "internal/cache/cache.go:42",
	}
}

func TestCockpitWorkItemCardPinsSelectionStyling(t *testing.T) {
	t.Parallel()
	tokens := ResolveTokens(true)
	model := &cockpitModel{tokens: tokens, selected: 0}
	item := queueFidelityCard()
	const paneWidth = 46
	rowWidth := paneWidth - 4

	unselected := model.workItemBlock(item, "Executing", 1, paneWidth)
	wantUnselected := []string{
		tokens.Running.Render("  [run] MAJOR") + strings.Repeat(" ", rowWidth-len("  [run] MAJOR")-len("#1")) + tokens.Muted.Render("#1"),
		"  Guard nil cache",
		tokens.Muted.Render("  internal/cache/cache.go:42"),
		"",
	}
	if !reflect.DeepEqual(unselected, wantUnselected) {
		t.Fatalf("unselected card mismatch\ngot:  %#v\nwant: %#v", unselected, wantUnselected)
	}

	selected := model.workItemBlock(item, "Executing", 0, paneWidth)
	cardContentWidth := rowWidth - 2
	inner := strings.Join([]string{
		tokens.Running.Render("> [run] MAJOR") + strings.Repeat(" ", cardContentWidth-len("> [run] MAJOR")-len("#1")) + tokens.Muted.Render("#1"),
		"  Guard nil cache",
		tokens.Muted.Render("  internal/cache/cache.go:42"),
	}, "\n")
	wantSelected := append(strings.Split(tokens.Selection.Width(rowWidth).Render(inner), "\n"), "")
	if !reflect.DeepEqual(selected, wantSelected) {
		t.Fatalf("selected card mismatch\ngot:  %#v\nwant: %#v", selected, wantSelected)
	}
}

func TestCockpitWorkItemCardNoColorKeepsMarkerDistinctions(t *testing.T) {
	t.Parallel()
	model := &cockpitModel{tokens: ResolveTokens(false), selected: 0}
	item := queueFidelityCard()

	unselected := strings.Join(model.workItemBlock(item, "Executing", 1, 46), "\n")
	selected := strings.Join(model.workItemBlock(item, "Executing", 0, 46), "\n")
	if strings.Contains(unselected, "\x1b") || strings.Contains(selected, "\x1b") {
		t.Fatalf("no-color card carries ANSI\nunselected: %q\nselected: %q", unselected, selected)
	}
	for _, marker := range []string{"  [run] MAJOR", "#1", "  Guard nil cache", "  internal/cache/cache.go:42"} {
		if !strings.Contains(unselected, marker) {
			t.Fatalf("expected no-color card to keep %q, got:\n%s", marker, unselected)
		}
	}
	if !strings.Contains(selected, "> [run] MAJOR") {
		t.Fatalf("expected no-color selected card to keep the > marker, got:\n%s", selected)
	}
	if !strings.Contains(selected, "┌") || !strings.Contains(selected, "│") {
		t.Fatalf("expected structural selection border without color, got:\n%s", selected)
	}
	if selected == unselected {
		t.Fatal("expected the no-color selected card to stay distinguishable from an unselected card")
	}
}

func TestCockpitBatchSeparatorStampDerivesFromRunEventTimestamps(t *testing.T) {
	t.Parallel()
	tokens := ResolveTokens(true)
	startedAt := time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
	events := []runevent.RunEvent{
		{Batch: 1, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonBatch, Summary: "Batch 001 executing.", Time: startedAt},
		{Batch: 1, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonBatch, Summary: "Batch 001 completed.", Time: startedAt.Add(38 * time.Second)},
	}
	clock := startedAt.Add(5 * time.Minute)
	model := newQueueFidelityCockpit(t, true, store.StateClean, events, clock)

	separator := model.batchSeparator(0, 46)
	lineWidth := 46 - 4
	want := tokens.SectionLabel.Render("BATCH 001/002") +
		strings.Repeat(" ", lineWidth-len("BATCH 001/002")-len("00:38")) +
		tokens.Muted.Render("00:38")
	if separator != want {
		t.Fatalf("settled Batch separator mismatch\ngot:  %q\nwant: %q", separator, want)
	}
	if got := stripANSI(model.batchSeparator(2, 46)); got != "BATCH 002/002" {
		t.Fatalf("expected Batch without timestamped events to carry no stamp, got %q", got)
	}

	plainModel := newQueueFidelityCockpit(t, false, store.StateClean, events, clock)
	plain := plainModel.batchSeparator(0, 46)
	if strings.Contains(plain, "\x1b") {
		t.Fatalf("no-color separator carries ANSI: %q", plain)
	}
	if plain != stripANSI(separator) {
		t.Fatalf("no-color separator text differs from styled text\nplain:  %q\nstyled: %q", plain, stripANSI(separator))
	}
}
