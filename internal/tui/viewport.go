package tui

import (
	"context"
	"math"
	"strings"

	"roundfix/internal/store"
)

// Window and paging bounds for the cockpit timeline. Internal constants by
// recorded decision: memory is bounded by the window, while the whole
// journal stays reachable by paging.
const (
	defaultTimelineWindowEvents = 500
	defaultTimelinePageSize     = 100
)

// tailCursorSentinel reads "the newest window" via the backward read.
const tailCursorSentinel = int64(math.MaxInt64)

// FollowState names what the cockpit viewport is doing, for the status bar.
type FollowState string

const (
	FollowReplaying FollowState = "replaying"
	FollowFollowing FollowState = "following"
	FollowScrolled  FollowState = "scrolled"
	FollowTerminal  FollowState = "terminal"
)

// TimelineSource is the journal surface the viewport pages. *store.Store
// satisfies it; the cockpit never consumes the live sink (ADR 0009).
type TimelineSource interface {
	RunEventsAfter(ctx context.Context, runID string, cursor int64, limit int) ([]store.JournalEvent, error)
	RunEventsBefore(ctx context.Context, runID string, cursor int64, limit int) ([]store.JournalEvent, error)
}

// timelineEntry keeps the rendered text per event. Texts concatenate before
// line-splitting so partial message chunks coalesce across events exactly
// like the streaming renderer.
type timelineEntry struct {
	cursor int64
	text   string
}

// TimelineViewport is the cockpit's timeline engine: a bounded sliding
// window of journal events rendered to lines, scrolled by line and paged by
// event at the window edges, with the Follow Mode state machine on top.
// It is synchronous and owns no goroutines; the driving program decides
// when to call Poll. Scrolling never affects the Run: the reader simply
// stops paging forward.
type TimelineViewport struct {
	source    TimelineSource
	runID     string
	maxEvents int
	pageSize  int

	entries []timelineEntry
	lines   []string
	atHead  bool
	atTail  bool

	terminal bool
	state    FollowState
	height   int
	scroll   int
	newBelow int
}

// NewTimelineViewport creates a viewport in the REPLAYING state. Non-positive
// bounds fall back to the internal defaults.
func NewTimelineViewport(source TimelineSource, runID string, maxEvents int, pageSize int) *TimelineViewport {
	if maxEvents <= 0 {
		maxEvents = defaultTimelineWindowEvents
	}
	if pageSize <= 0 {
		pageSize = defaultTimelinePageSize
	}
	if pageSize > maxEvents {
		pageSize = maxEvents
	}
	return &TimelineViewport{
		source:    source,
		runID:     runID,
		maxEvents: maxEvents,
		pageSize:  pageSize,
		state:     FollowReplaying,
		height:    24,
	}
}

// Replay loads the newest window from the journal and enters Follow Mode
// (or the terminal state for finished Runs).
func (viewport *TimelineViewport) Replay(ctx context.Context) error {
	page, err := viewport.source.RunEventsBefore(ctx, viewport.runID, tailCursorSentinel, viewport.maxEvents)
	if err != nil {
		return err
	}
	viewport.entries = entriesFromJournal(page)
	viewport.atHead = len(page) < viewport.maxEvents
	viewport.atTail = true
	viewport.newBelow = 0
	viewport.rebuildLines()
	viewport.enterTailState()
	viewport.scrollToBottom()
	return nil
}

// Poll advances the viewport on a change signal. FOLLOWING appends new
// events and keeps the tail pinned; SCROLLED only counts arrivals below the
// frozen viewport; terminal Runs ignore polls.
func (viewport *TimelineViewport) Poll(ctx context.Context) error {
	switch viewport.state {
	case FollowFollowing:
		if err := viewport.drainForward(ctx); err != nil {
			return err
		}
	case FollowTerminal:
		// Late events (the terminal outcome) still land right after a Run
		// ends; keep draining while pinned to the bottom, but never yank a
		// user who is reading history.
		if viewport.atBottom() {
			if err := viewport.drainForward(ctx); err != nil {
				return err
			}
		}
	case FollowScrolled:
		page, err := viewport.source.RunEventsAfter(ctx, viewport.runID, viewport.windowEnd(), viewport.pageSize)
		if err != nil {
			return err
		}
		viewport.newBelow = len(page)
		if len(page) > 0 {
			viewport.atTail = false
		}
	}
	return nil
}

func (viewport *TimelineViewport) drainForward(ctx context.Context) error {
	for {
		page, err := viewport.source.RunEventsAfter(ctx, viewport.runID, viewport.tailCursor(), viewport.pageSize)
		if err != nil {
			return err
		}
		viewport.appendEntries(entriesFromJournal(page))
		if len(page) < viewport.pageSize {
			break
		}
	}
	viewport.atTail = true
	viewport.scrollToBottom()
	return nil
}

// ScrollUp moves the viewport toward older events, freezing Follow Mode and
// paging the window backward at the top edge.
func (viewport *TimelineViewport) ScrollUp(ctx context.Context, lines int) error {
	if lines <= 0 || viewport.state == FollowReplaying {
		return nil
	}
	// Nothing above to scroll to: stay in Follow Mode instead of freezing
	// a viewport that cannot move.
	if viewport.state == FollowFollowing && viewport.maxScroll() == 0 && viewport.atHead {
		return nil
	}
	if viewport.state == FollowFollowing {
		viewport.state = FollowScrolled
	}
	viewport.scroll -= lines
	for viewport.scroll < 0 && !viewport.atHead {
		if err := viewport.pageBackward(ctx); err != nil {
			return err
		}
	}
	if viewport.scroll < 0 {
		viewport.scroll = 0
	}
	return nil
}

// ScrollDown moves the viewport toward newer events, paging forward at the
// bottom edge and resuming Follow Mode when the journal tail is reached.
func (viewport *TimelineViewport) ScrollDown(ctx context.Context, lines int) error {
	if lines <= 0 || viewport.state == FollowFollowing || viewport.state == FollowReplaying {
		return nil
	}
	viewport.scroll += lines
	for viewport.scroll > viewport.maxScroll() && !viewport.atTail {
		if err := viewport.pageForward(ctx); err != nil {
			return err
		}
	}
	if viewport.scroll > viewport.maxScroll() {
		viewport.scroll = viewport.maxScroll()
	}
	if viewport.atBottom() && viewport.atTail {
		viewport.newBelow = 0
		viewport.enterTailState()
	}
	return nil
}

// PageUp and PageDown scroll by one viewport height.
func (viewport *TimelineViewport) PageUp(ctx context.Context) error {
	return viewport.ScrollUp(ctx, maxInt(viewport.height-1, 1))
}

func (viewport *TimelineViewport) PageDown(ctx context.Context) error {
	return viewport.ScrollDown(ctx, maxInt(viewport.height-1, 1))
}

// JumpToTop goes to the top of the loaded window.
func (viewport *TimelineViewport) JumpToTop() {
	if viewport.state == FollowReplaying {
		return
	}
	if viewport.state == FollowFollowing {
		viewport.state = FollowScrolled
	}
	viewport.scroll = 0
}

// JumpToTail reloads the newest window and resumes Follow Mode (End/G).
func (viewport *TimelineViewport) JumpToTail(ctx context.Context) error {
	return viewport.Replay(ctx)
}

// SetTerminal marks the Run as finished: the viewport never follows again,
// but scrollback keeps working. A user mid-scrollback stays where they are.
func (viewport *TimelineViewport) SetTerminal() {
	viewport.terminal = true
	if viewport.state == FollowFollowing || viewport.state == FollowReplaying {
		viewport.state = FollowTerminal
	}
}

// SetHeight resizes the viewport, keeping the tail pinned in Follow Mode.
func (viewport *TimelineViewport) SetHeight(height int) {
	if height < 1 {
		height = 1
	}
	viewport.height = height
	if viewport.state == FollowFollowing || viewport.state == FollowTerminal {
		viewport.scrollToBottom()
		return
	}
	if viewport.scroll > viewport.maxScroll() {
		viewport.scroll = viewport.maxScroll()
	}
}

// VisibleLines returns the viewport-height slice of rendered lines at the
// current scroll position.
func (viewport *TimelineViewport) VisibleLines() []string {
	if len(viewport.lines) == 0 {
		return nil
	}
	start := viewport.scroll
	if start > len(viewport.lines) {
		start = len(viewport.lines)
	}
	end := start + viewport.height
	if end > len(viewport.lines) {
		end = len(viewport.lines)
	}
	return viewport.lines[start:end]
}

// State reports the Follow Mode state and the count of new events below a
// frozen viewport, for the status bar.
func (viewport *TimelineViewport) State() (FollowState, int) {
	return viewport.state, viewport.newBelow
}

func (viewport *TimelineViewport) enterTailState() {
	if viewport.terminal {
		viewport.state = FollowTerminal
		return
	}
	viewport.state = FollowFollowing
}

func (viewport *TimelineViewport) pageBackward(ctx context.Context) error {
	page, err := viewport.source.RunEventsBefore(ctx, viewport.runID, viewport.windowStart(), viewport.pageSize)
	if err != nil {
		return err
	}
	if len(page) < viewport.pageSize {
		viewport.atHead = true
	}
	if len(page) == 0 {
		return nil
	}
	// Phase 1: prepend and shift the scroll by exactly the lines added at
	// the top, so the visible region stays put.
	before := len(viewport.lines)
	viewport.entries = append(entriesFromJournal(page), viewport.entries...)
	viewport.rebuildLines()
	viewport.scroll += len(viewport.lines) - before
	// Phase 2: evict from the bottom — which never shifts top offsets —
	// but keep every line the viewport can still see.
	for len(viewport.entries) > viewport.maxEvents {
		candidate := viewport.entries[:len(viewport.entries)-1]
		if linesOfEntries(candidate) < viewport.scroll+viewport.height {
			break
		}
		viewport.entries = candidate
		viewport.atTail = false
	}
	viewport.rebuildLines()
	return nil
}

func (viewport *TimelineViewport) pageForward(ctx context.Context) error {
	page, err := viewport.source.RunEventsAfter(ctx, viewport.runID, viewport.windowEnd(), viewport.pageSize)
	if err != nil {
		return err
	}
	if len(page) < viewport.pageSize {
		viewport.atTail = true
	}
	if len(page) == 0 {
		return nil
	}
	viewport.appendEntries(entriesFromJournal(page))
	return nil
}

func (viewport *TimelineViewport) appendEntries(entries []timelineEntry) {
	if len(entries) == 0 {
		return
	}
	viewport.entries = append(viewport.entries, entries...)
	viewport.rebuildLines()
	totalAfterAppend := len(viewport.lines)
	for len(viewport.entries) > viewport.maxEvents {
		viewport.entries = viewport.entries[1:]
		viewport.atHead = false
	}
	viewport.rebuildLines()
	// Head eviction removes lines above the viewport; shift the scroll by
	// the exact number that disappeared so the visible region is stable.
	viewport.scroll -= totalAfterAppend - len(viewport.lines)
	if viewport.scroll < 0 {
		viewport.scroll = 0
	}
}

func (viewport *TimelineViewport) windowStart() int64 {
	if len(viewport.entries) == 0 {
		return tailCursorSentinel
	}
	return viewport.entries[0].cursor
}

func (viewport *TimelineViewport) windowEnd() int64 {
	if len(viewport.entries) == 0 {
		return 0
	}
	return viewport.entries[len(viewport.entries)-1].cursor
}

func (viewport *TimelineViewport) tailCursor() int64 {
	return viewport.windowEnd()
}

func (viewport *TimelineViewport) rebuildLines() {
	viewport.lines = splitRenderedLines(concatEntryText(viewport.entries))
}

func (viewport *TimelineViewport) maxScroll() int {
	if total := len(viewport.lines); total > viewport.height {
		return total - viewport.height
	}
	return 0
}

func (viewport *TimelineViewport) atBottom() bool {
	return viewport.scroll >= viewport.maxScroll()
}

func (viewport *TimelineViewport) scrollToBottom() {
	viewport.scroll = viewport.maxScroll()
}

func entriesFromJournal(page []store.JournalEvent) []timelineEntry {
	entries := make([]timelineEntry, 0, len(page))
	for _, journal := range page {
		entries = append(entries, timelineEntry{
			cursor: journal.Cursor,
			text:   timelineText(journal.Event),
		})
	}
	return entries
}

func concatEntryText(entries []timelineEntry) string {
	var builder strings.Builder
	for _, entry := range entries {
		builder.WriteString(entry.text)
	}
	return builder.String()
}

func linesOfEntries(entries []timelineEntry) int {
	return len(splitRenderedLines(concatEntryText(entries)))
}

// splitRenderedLines splits rendered text into display lines; a trailing
// partial line (no final newline) still occupies one line.
func splitRenderedLines(text string) []string {
	if text == "" {
		return nil
	}
	parts := strings.Split(text, "\n")
	if parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
