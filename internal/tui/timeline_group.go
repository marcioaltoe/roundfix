package tui

import (
	"fmt"
	"strings"

	"roundfix/internal/runevent"
)

func renderedTimelineLines(entries []timelineEntry) []string {
	if !hasBatchedTimelineEntries(entries) {
		return splitRenderedLines(concatEntryText(entries))
	}

	total := maxTimelineBatch(entries)
	lines := []string{}
	group := []timelineEntry{}
	flush := func() {
		if len(group) == 0 {
			return
		}
		lines = append(lines, renderTimelineBatch(group, total)...)
		group = nil
	}

	for _, entry := range entries {
		if entry.event.Batch <= 0 {
			flush()
			lines = append(lines, splitRenderedLines(entry.text)...)
			continue
		}
		if len(group) > 0 && group[0].event.Batch != entry.event.Batch {
			flush()
		}
		group = append(group, entry)
	}
	flush()
	return lines
}

func hasBatchedTimelineEntries(entries []timelineEntry) bool {
	for _, entry := range entries {
		if entry.event.Batch > 0 {
			return true
		}
	}
	return false
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
	body := timelineBatchBody(entries)
	if state == "" && len(body) == 0 {
		return nil
	}

	parts := []string{fmt.Sprintf("BATCH %03d/%03d", batch, total)}
	if state != "" {
		parts = append(parts, state)
	}
	if elapsed != "" {
		parts = append(parts, elapsed)
	}
	lines := []string{strings.Join(parts, " ")}
	lines = append(lines, body...)
	return lines
}

func timelineBatchBody(entries []timelineEntry) []string {
	lines := []string{}
	lastSection := ""
	var agentText strings.Builder
	flushAgentText := func() {
		if agentText.Len() == 0 {
			return
		}
		lines = append(lines, splitRenderedLines(agentText.String())...)
		agentText.Reset()
	}
	for _, entry := range entries {
		if entry.text == "" || entry.kind == runevent.KindDaemonBatch {
			continue
		}
		section := daemonTimelineSection(entry.kind)
		if section != "" {
			flushAgentText()
			if section != lastSection {
				lines = append(lines, section)
				lastSection = section
			}
			lines = append(lines, splitRenderedLines(entry.text)...)
		} else {
			lastSection = ""
			agentText.WriteString(entry.text)
		}
	}
	flushAgentText()
	return lines
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
