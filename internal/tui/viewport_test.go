package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

// fakeTimelineSource serves an in-memory journal through the same paging
// contract as the store.
type fakeTimelineSource struct {
	events []store.JournalEvent
	reads  int
}

func (source *fakeTimelineSource) RunEventsAfter(_ context.Context, _ string, cursor int64, limit int) ([]store.JournalEvent, error) {
	source.reads++
	page := []store.JournalEvent{}
	for _, entry := range source.events {
		if entry.Cursor > cursor && len(page) < limit {
			page = append(page, entry)
		}
	}
	return page, nil
}

func (source *fakeTimelineSource) RunEventsBefore(_ context.Context, _ string, cursor int64, limit int) ([]store.JournalEvent, error) {
	source.reads++
	page := []store.JournalEvent{}
	for index := len(source.events) - 1; index >= 0 && len(page) < limit; index-- {
		if source.events[index].Cursor < cursor {
			page = append(page, source.events[index])
		}
	}
	for left, right := 0, len(page)-1; left < right; left, right = left+1, right-1 {
		page[left], page[right] = page[right], page[left]
	}
	return page, nil
}

func (source *fakeTimelineSource) addLine(text string) {
	cursor := int64(len(source.events) + 1)
	source.events = append(source.events, store.JournalEvent{
		Cursor: cursor,
		Event: runevent.RunEvent{
			Source:  runevent.SourceAgent,
			Kind:    runevent.KindAgentRaw,
			Payload: []byte(`{"text":` + strconv.Quote(text) + `}`),
		},
	})
}

func seededSource(total int) *fakeTimelineSource {
	source := &fakeTimelineSource{}
	for index := 1; index <= total; index++ {
		source.addLine(fmt.Sprintf("line %04d\n", index))
	}
	return source
}

func TestViewportReplayLoadsNewestBoundedWindowAndFollows(t *testing.T) {
	ctx := context.Background()
	source := seededSource(60)
	viewport := NewTimelineViewport(source, "run-1", 50, 10)
	viewport.SetHeight(5)

	if state, _ := viewport.State(); state != FollowReplaying {
		t.Fatalf("expected REPLAYING before replay, got %q", state)
	}
	if err := viewport.Replay(ctx); err != nil {
		t.Fatalf("replay: %v", err)
	}

	if state, below := viewport.State(); state != FollowFollowing || below != 0 {
		t.Fatalf("expected FOLLOWING after replay, got %q %d", state, below)
	}
	visible := viewport.VisibleLines()
	if len(visible) != 5 || visible[4] != "line 0060" || visible[0] != "line 0056" {
		t.Fatalf("expected tail visible after replay, got %v", visible)
	}
}

func TestViewportFollowingAdvancesTailAndKeepsWindowBound(t *testing.T) {
	ctx := context.Background()
	source := seededSource(50)
	viewport := NewTimelineViewport(source, "run-1", 50, 10)
	viewport.SetHeight(3)
	if err := viewport.Replay(ctx); err != nil {
		t.Fatalf("replay: %v", err)
	}

	source.addLine("line 0051\n")
	source.addLine("line 0052\n")
	if err := viewport.Poll(ctx); err != nil {
		t.Fatalf("poll: %v", err)
	}

	visible := viewport.VisibleLines()
	if visible[len(visible)-1] != "line 0052" {
		t.Fatalf("expected tail auto-advance, got %v", visible)
	}
	if len(viewport.entries) != 50 {
		t.Fatalf("expected window bound kept at 50, got %d", len(viewport.entries))
	}
}

func TestViewportScrolledFreezesAndCountsNewEventsBelow(t *testing.T) {
	ctx := context.Background()
	source := seededSource(40)
	viewport := NewTimelineViewport(source, "run-1", 50, 10)
	viewport.SetHeight(4)
	if err := viewport.Replay(ctx); err != nil {
		t.Fatalf("replay: %v", err)
	}

	if err := viewport.ScrollUp(ctx, 10); err != nil {
		t.Fatalf("scroll up: %v", err)
	}
	frozen := strings.Join(viewport.VisibleLines(), "|")
	if state, _ := viewport.State(); state != FollowScrolled {
		t.Fatalf("expected SCROLLED after scrolling up, got %q", state)
	}

	source.addLine("line 0041\n")
	source.addLine("line 0042\n")
	source.addLine("line 0043\n")
	if err := viewport.Poll(ctx); err != nil {
		t.Fatalf("poll: %v", err)
	}

	if got := strings.Join(viewport.VisibleLines(), "|"); got != frozen {
		t.Fatalf("expected frozen viewport while scrolled, got %q want %q", got, frozen)
	}
	if state, below := viewport.State(); state != FollowScrolled || below != 3 {
		t.Fatalf("expected SCROLLED with 3 new below, got %q %d", state, below)
	}

	if err := viewport.JumpToTail(ctx); err != nil {
		t.Fatalf("jump to tail: %v", err)
	}
	if state, below := viewport.State(); state != FollowFollowing || below != 0 {
		t.Fatalf("expected FOLLOWING after End, got %q %d", state, below)
	}
	visible := viewport.VisibleLines()
	if visible[len(visible)-1] != "line 0043" {
		t.Fatalf("expected newest tail after End, got %v", visible)
	}
}

func TestViewportScrollUpPagesBackwardAcrossWindowEdgeWithoutGaps(t *testing.T) {
	ctx := context.Background()
	source := seededSource(120)
	viewport := NewTimelineViewport(source, "run-1", 30, 10)
	viewport.SetHeight(4)
	if err := viewport.Replay(ctx); err != nil {
		t.Fatalf("replay: %v", err)
	}

	// Walk to the journal head collecting every visible top line.
	seen := map[string]bool{}
	capture := func() {
		for _, line := range viewport.VisibleLines() {
			seen[line] = true
		}
	}
	capture()
	for step := 0; step < 200; step++ {
		if err := viewport.ScrollUp(ctx, 1); err != nil {
			t.Fatalf("scroll up: %v", err)
		}
		capture()
		if viewport.atHead && viewport.scroll == 0 {
			break
		}
	}

	for index := 1; index <= 120; index++ {
		if !seen[fmt.Sprintf("line %04d", index)] {
			t.Fatalf("expected line %04d visible during backward walk; window slide dropped it", index)
		}
	}
}

func TestViewportScrollDownPagesForwardAndResumesFollowingAtTail(t *testing.T) {
	ctx := context.Background()
	source := seededSource(120)
	viewport := NewTimelineViewport(source, "run-1", 30, 10)
	viewport.SetHeight(4)
	if err := viewport.Replay(ctx); err != nil {
		t.Fatalf("replay: %v", err)
	}
	// Deep scrollback to the journal head first.
	for step := 0; step < 300 && !(viewport.atHead && viewport.scroll == 0); step++ {
		if err := viewport.ScrollUp(ctx, 1); err != nil {
			t.Fatalf("scroll up: %v", err)
		}
	}

	seen := map[string]bool{}
	for step := 0; step < 400; step++ {
		for _, line := range viewport.VisibleLines() {
			seen[line] = true
		}
		if state, _ := viewport.State(); state == FollowFollowing {
			break
		}
		if err := viewport.ScrollDown(ctx, 1); err != nil {
			t.Fatalf("scroll down: %v", err)
		}
	}

	if state, _ := viewport.State(); state != FollowFollowing {
		t.Fatalf("expected FOLLOWING resumed at the journal tail, got %q", state)
	}
	for index := 1; index <= 120; index++ {
		if !seen[fmt.Sprintf("line %04d", index)] {
			t.Fatalf("expected line %04d visible during forward walk", index)
		}
	}
}

func TestViewportTerminalRunsNeverFollow(t *testing.T) {
	ctx := context.Background()
	source := seededSource(20)
	viewport := NewTimelineViewport(source, "run-1", 50, 10)
	viewport.SetHeight(4)
	viewport.SetTerminal()
	if err := viewport.Replay(ctx); err != nil {
		t.Fatalf("replay: %v", err)
	}

	if state, _ := viewport.State(); state != FollowTerminal {
		t.Fatalf("expected terminal state after replay, got %q", state)
	}
	if err := viewport.ScrollUp(ctx, 5); err != nil {
		t.Fatalf("scroll up: %v", err)
	}
	if err := viewport.ScrollDown(ctx, 50); err != nil {
		t.Fatalf("scroll down: %v", err)
	}
	if state, _ := viewport.State(); state != FollowTerminal {
		t.Fatalf("expected terminal Runs to never enter FOLLOWING, got %q", state)
	}
	if err := viewport.Poll(ctx); err != nil {
		t.Fatalf("poll on terminal run: %v", err)
	}
}

func TestViewportCoalescesChunksAcrossEventsAndWindowSlides(t *testing.T) {
	ctx := context.Background()
	source := &fakeTimelineSource{}
	for index := 1; index <= 10; index++ {
		source.addLine(fmt.Sprintf("line %04d\n", index))
	}
	// One logical line split across three chunk events.
	source.addLine("Hel")
	source.addLine("lo ")
	source.addLine("world\n")
	source.addLine("line tail\n")
	viewport := NewTimelineViewport(source, "run-1", 50, 5)
	viewport.SetHeight(3)
	if err := viewport.Replay(ctx); err != nil {
		t.Fatalf("replay: %v", err)
	}

	visible := viewport.VisibleLines()
	joined := strings.Join(visible, "|")
	if !strings.Contains(joined, "Hello world") {
		t.Fatalf("expected chunks coalesced into one line, got %v", visible)
	}
	if strings.Contains(joined, "Hel|") || strings.Contains(joined, "lo |") {
		t.Fatalf("expected no fragment lines, got %v", visible)
	}
}

func TestViewportSkipsUnknownKindsWithoutBlankLines(t *testing.T) {
	ctx := context.Background()
	source := &fakeTimelineSource{}
	source.addLine("line one\n")
	source.events = append(source.events, store.JournalEvent{
		Cursor: int64(len(source.events) + 1),
		Event:  runevent.RunEvent{Source: runevent.SourceDaemon, Kind: "future.unknown", Summary: "future", Payload: []byte(`{}`)},
	})
	source.addLine("line two\n")
	viewport := NewTimelineViewport(source, "run-1", 50, 10)
	viewport.SetHeight(5)
	if err := viewport.Replay(ctx); err != nil {
		t.Fatalf("replay: %v", err)
	}

	visible := strings.Join(viewport.VisibleLines(), "|")
	if visible != "line one|line two" {
		t.Fatalf("expected unknown kinds skipped without blank lines, got %q", visible)
	}
}
