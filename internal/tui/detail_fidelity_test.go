package tui

import (
	"strings"
	"testing"

	"roundfix/internal/rounds"
	"roundfix/internal/store"

	tea "charm.land/bubbletea/v2"
)

// Suite: Detail Modal fidelity and pane empty states
// Invariant: the modal renders the accent frame, tokened title block, sectioned body, and position footer with a marker-only no-color twin; empty panes explain the Run kind and what populates them.
// Boundary IN: renderDetail/renderCockpitDetailPane and the pane empty-state renderers.
// Boundary OUT: modal open/close/scroll semantics, pinned by the existing cockpit modal suite.
func detailFidelityView() *issueDetailView {
	return &issueDetailView{
		kind:    detailReviewIssue,
		ordinal: 1,
		issue: rounds.Issue{
			Title:     "Guard nil cache",
			Severity:  "major",
			Status:    "pending",
			File:      "internal/cache/cache.go",
			Line:      42,
			SourceRef: "thread:PRRT_1",
		},
		lines: []string{"## Review Comment", "Guard the map lookup."},
	}
}

func TestCockpitDetailModalStyledThroughTokensWithNoColorTwin(t *testing.T) {
	t.Parallel()
	tokens := ResolveTokens(true)
	model := &cockpitModel{tokens: tokens, detail: detailFidelityView()}
	const width, height = 60, 14
	innerWidth := width - 4

	inner := strings.Split(model.renderDetail(width, height), "\n")
	titlePadding := innerWidth - displayWidth("REVIEW.ISSUE  #001") - displayWidth("Esc close · j/k scroll")
	want := map[int]string{
		0: tokens.SectionLabel.Render("REVIEW.ISSUE  #001") + strings.Repeat(" ", titlePadding) + tokens.Muted.Render("Esc close · j/k scroll"),
		1: tokens.Muted.Render(strings.Repeat("-", innerWidth)),
		2: "Guard nil cache",
		3: tokens.Waiting.Render("major") + tokens.Muted.Render(" · ") + tokens.Waiting.Render("pending") + tokens.Muted.Render(" · internal/cache/cache.go:42"),
		4: tokens.Muted.Render("source: thread:PRRT_1"),
		6: tokens.SectionLabel.Render("## Review Comment"),
		7: "Guard the map lookup.",
	}
	for index, expected := range want {
		if inner[index] != expected {
			t.Fatalf("modal line %d mismatch\ngot:  %q\nwant: %q", index, inner[index], expected)
		}
	}
	footer := inner[len(inner)-1]
	if wantFooter := tokens.Muted.Render("Line 1-2 of 2 · PgUp/PgDn page"); footer != wantFooter {
		t.Fatalf("modal position footer mismatch\ngot:  %q\nwant: %q", footer, wantFooter)
	}

	framed := renderCockpitDetailPane(model, width, height)
	if !strings.Contains(framed, ansi256("39")) {
		t.Fatalf("expected the accent frame to carry the cyan border, got %q", framed)
	}
	if !strings.Contains(stripANSI(framed), "┌") {
		t.Fatalf("expected the modal frame structure, got:\n%s", stripANSI(framed))
	}

	plainModel := &cockpitModel{tokens: ResolveTokens(false), detail: detailFidelityView()}
	plain := renderCockpitDetailPane(plainModel, width, height)
	if strings.Contains(plain, "\x1b") {
		t.Fatalf("no-color modal carries ANSI: %q", plain)
	}
	if plain != stripANSI(framed) {
		t.Fatalf("no-color modal text differs from styled text\nplain:\n%s\nstyled:\n%s", plain, stripANSI(framed))
	}
	for _, marker := range []string{"REVIEW.ISSUE  #001", "major · pending · internal/cache/cache.go:42", "## Review Comment", "Line 1-2 of 2"} {
		if !strings.Contains(plain, marker) {
			t.Fatalf("expected no-color modal to keep marker %q, got:\n%s", marker, plain)
		}
	}
}

func TestCockpitFetchRunAttachRendersExplanatoryEmptyStates(t *testing.T) {
	t.Parallel()
	source := &cockpitFakeSource{run: store.Run{ID: "run-fetch", State: store.StateFetched}, version: 1}
	model := newTestCockpit(t, source, LiveRunView{
		Command:       "attach",
		PRNumber:      "123",
		RunKind:       store.KindFetch,
		PipelineState: store.StateFetched,
	})
	model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	rendered := viewText(model)
	for _, expected := range []string{
		"No Work Items.",
		"A Fetch Run writes Review artifacts",
		"to disk and starts no Agent.",
		"No Run Events.",
		"A Fetch Run writes Review artifacts to disk and starts no Agent.",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("expected Fetch Run empty state %q, got:\n%s", expected, rendered)
		}
	}
	for _, bare := range []string{"No Work Items yet.", "No Run Events yet."} {
		if strings.Contains(rendered, bare) {
			t.Fatalf("expected the Fetch Run copy to replace the generic placeholder %q, got:\n%s", bare, rendered)
		}
	}
}

func TestCockpitEmptyStatesNameReviewAndSpecRunExpectations(t *testing.T) {
	t.Parallel()
	reviewSource := &cockpitFakeSource{run: store.Run{ID: "run-review", State: store.StateActive}, version: 1}
	review := newTestCockpit(t, reviewSource, LiveRunView{PipelineState: store.StateActive})
	review.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	reviewRendered := viewText(review)
	for _, expected := range []string{
		"No Work Items yet.",
		"Review Issues queue here once the",
		"Run records them.",
		"No Run Events yet.",
		"Agent and Daemon activity streams here.",
	} {
		if !strings.Contains(reviewRendered, expected) {
			t.Fatalf("expected review Run empty state %q, got:\n%s", expected, reviewRendered)
		}
	}

	specSource := &cockpitFakeSource{run: store.Run{ID: "run-spec", State: store.StateActive}, version: 1}
	specRun := newTestCockpit(t, specSource, LiveRunView{
		PipelineState: store.StateActive,
		RunKind:       store.KindImplement,
		SpecSlug:      "0021-cockpit-visual-fidelity",
		GitRoot:       t.TempDir(),
	})
	specRun.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	specRendered := viewText(specRun)
	for _, expected := range []string{
		"No Tasks yet.",
		"An implement Run queues its Spec's",
		"Task Graph here.",
		"Task execution and verification stream here.",
	} {
		if !strings.Contains(specRendered, expected) {
			t.Fatalf("expected spec Run empty state %q, got:\n%s", expected, specRendered)
		}
	}
}
