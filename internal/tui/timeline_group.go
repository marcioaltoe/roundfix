package tui

import (
	"fmt"
	"strings"
	"time"

	"roundfix/internal/runevent"
)

const (
	timelineExpandedMarker  = "▼"
	timelineCollapsedMarker = "▶"
	// timelineGutterWidth is the aligned timestamp column: "HH:MM:SS "
	// or blanks when the event carries no timestamp.
	timelineGutterWidth = 9
	// timelineSummaryBound caps a summary row at build time; the pane
	// re-bounds every row to its own width at render time.
	timelineSummaryBound = 512
)

func renderedTimelineLines(entries []timelineEntry) []string {
	total := maxTimelineBatch(entries)
	lines := []string{}
	group := []timelineEntry{}
	plain := []timelineEntry{}
	flushPlain := func() {
		if len(plain) == 0 {
			return
		}
		lines = append(lines, timelineEventRows(plain)...)
		plain = nil
	}
	flushGroup := func() {
		if len(group) == 0 {
			return
		}
		lines = append(lines, renderTimelineBatch(group, total)...)
		group = nil
	}
	for _, entry := range entries {
		if entry.event.Batch <= 0 {
			flushGroup()
			plain = append(plain, entry)
			continue
		}
		flushPlain()
		if len(group) > 0 && group[0].event.Batch != entry.event.Batch {
			flushGroup()
		}
		group = append(group, entry)
	}
	flushPlain()
	flushGroup()
	return lines
}

func maxTimelineBatch(entries []timelineEntry) int {
	total := 0
	for _, entry := range entries {
		if entry.event.Batch > total {
			total = entry.event.Batch
		}
	}
	return total
}

// renderTimelineBatch renders one Batch group: a marker-carrying header
// row folding the daemon.batch state and elapsed clock, then the Batch's
// event rows. Settled Batches collapse to the header — their summary row —
// while every other Batch renders expanded. Collapse is state-driven; no
// key toggles it.
func renderTimelineBatch(entries []timelineEntry, total int) []string {
	if len(entries) == 0 {
		return nil
	}
	batch := entries[0].event.Batch
	if total < batch {
		total = batch
	}
	state := timelineBatchState(entries)
	elapsed := timelineBatchElapsed(entries)
	body := make([]timelineEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.kind == runevent.KindDaemonBatch {
			continue
		}
		body = append(body, entry)
	}
	rows := timelineEventRows(body)
	if state == "" && len(rows) == 0 {
		return nil
	}
	settled := timelineBatchSettled(state)
	marker := timelineExpandedMarker
	if settled {
		marker = timelineCollapsedMarker
	}
	parts := []string{marker, fmt.Sprintf("BATCH %03d/%03d", batch, total)}
	if state != "" {
		parts = append(parts, state)
	}
	if elapsed != "" {
		parts = append(parts, elapsed)
	}
	lines := []string{strings.Join(parts, " ")}
	if settled {
		return lines
	}
	return append(lines, rows...)
}

// timelineEventRows renders a run of entries: structured kinds render one
// bounded summary row each behind the timestamp gutter, while chunked
// message/raw text keeps coalescing into whole console lines so streaming
// fragments never render as broken rows.
func timelineEventRows(entries []timelineEntry) []string {
	rows := []string{}
	var stream strings.Builder
	flushStream := func() {
		if stream.Len() == 0 {
			return
		}
		rows = append(rows, splitRenderedLines(stream.String())...)
		stream.Reset()
	}
	for _, entry := range entries {
		if timelineStreamKind(entry.kind) {
			stream.WriteString(entry.text)
			continue
		}
		row := timelineEventRow(entry)
		if row == "" {
			continue
		}
		flushStream()
		rows = append(rows, row)
	}
	flushStream()
	return rows
}

// timelineStreamKind reports the chunked streaming kinds whose text
// coalesces across events; every other known kind renders as one bounded
// summary row.
func timelineStreamKind(kind runevent.Kind) bool {
	return kind == runevent.KindAgentMessage || kind == runevent.KindAgentRaw
}

// timelineEventRow renders one structured event as its single bounded
// summary row: timestamp gutter, kind label, summary line. Raw payloads
// never render inline. Unknown kinds and suppressed console text (session
// lifecycle statuses) keep the shipped skip policy.
func timelineEventRow(entry timelineEntry) string {
	event := entry.event
	if runevent.IsDaemonKind(event.Kind) {
		summary := timelineRowSummary(entry)
		if summary == "" {
			return ""
		}
		if label := daemonTimelineSection(event.Kind); label != "" {
			summary = label + " " + summary
		}
		return timelineGutter(event.Time) + summary
	}
	if entry.text == "" {
		return ""
	}
	summary := timelineRowSummary(entry)
	if summary == "" {
		return ""
	}
	return timelineGutter(event.Time) + summary
}

// timelineRowSummary bounds the row through the shared summary helper.
// Events journaled without a summary fall back to their reconstructed
// console text so older journals stay viewable, bounded the same way.
func timelineRowSummary(entry timelineEntry) string {
	if record, ok, err := runevent.ProjectSelectionLifecycle(entry.event); ok {
		if err != nil {
			return ""
		}
		return runevent.SelectionLifecycleSummary(record)
	}
	if strings.TrimSpace(entry.event.Summary) != "" {
		return EventSummary(entry.event, timelineSummaryBound)
	}
	fallback := entry.event
	fallback.Summary = entry.text
	return EventSummary(fallback, timelineSummaryBound)
}

// timelineGutter renders the aligned timestamp column; events without a
// timestamp keep the column blank so rows stay aligned across kinds.
func timelineGutter(at time.Time) string {
	if at.IsZero() {
		return strings.Repeat(" ", timelineGutterWidth)
	}
	return at.Format("15:04:05") + " "
}

// timelineBatchSettled reports whether the Batch's journaled state word
// marks it settled; settled Batches collapse to their summary row.
func timelineBatchSettled(state string) bool {
	switch state {
	case "completed", "failed", "stopped":
		return true
	default:
		return false
	}
}

func daemonTimelineSection(kind runevent.Kind) string {
	switch kind {
	case runevent.KindDaemonStatus, runevent.KindDaemonReviewStatus,
		runevent.KindDaemonQuietPeriod, runevent.KindDaemonSelection:
		return "STATUS"
	case runevent.KindDaemonFetch:
		return "FETCH"
	case runevent.KindDaemonVerification:
		return "VERIFY"
	case runevent.KindDaemonCommit:
		return "COMMIT"
	case runevent.KindDaemonPush:
		return "PUSH"
	case runevent.KindDaemonSourceResolution:
		return "SOURCE"
	case runevent.KindDaemonRetry:
		return "RETRY"
	case runevent.KindDaemonOutcome:
		return "OUTCOME"
	case runevent.KindDaemonTask:
		return "TASK"
	case runevent.KindDaemonQA:
		return "QA"
	case runevent.KindDaemonAgentSelectionAttempt, runevent.KindDaemonAgentSelectionActive,
		runevent.KindDaemonAgentSelectionFallback, runevent.KindDaemonAgentSelectionExhausted,
		runevent.KindDaemonAgentSelectionClosed:
		return "SELECTION"
	default:
		return ""
	}
}

func timelineBatchState(entries []timelineEntry) string {
	state := ""
	for _, entry := range entries {
		if entry.kind != runevent.KindDaemonBatch {
			continue
		}
		if next := timelineBatchStateFromSummary(entry.event.Summary); next != "" {
			state = next
		}
	}
	return state
}

func timelineBatchStateFromSummary(summary string) string {
	normalized := strings.ToLower(summary)
	replacer := strings.NewReplacer(".", " ", ",", " ", ";", " ", ":", " ", "(", " ", ")", " ")
	words := strings.Fields(replacer.Replace(normalized))
	for _, state := range []string{"executing", "waiting", "started", "running", "completed", "failed", "stopped"} {
		for _, word := range words {
			if word == state {
				return state
			}
		}
	}
	return ""
}

func timelineBatchElapsed(entries []timelineEntry) string {
	firstIndex := -1
	lastIndex := -1
	for index, entry := range entries {
		if entry.event.Time.IsZero() {
			continue
		}
		if firstIndex < 0 {
			firstIndex = index
		}
		lastIndex = index
	}
	if firstIndex < 0 {
		return ""
	}
	return formatElapsed(entries[lastIndex].event.Time.Sub(entries[firstIndex].event.Time))
}
