package tui

import (
	"strings"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
)

// RunTimeline is the one timeline renderer for Run Events: live delivery
// and journal replay are two adapters feeding it. Console memory stays
// bounded by the ring buffer, and high-frequency message chunks coalesce
// into whole lines before rendering.
type RunTimeline struct {
	console StreamBuffer
}

func NewRunTimeline(maxLines int) *RunTimeline {
	return &RunTimeline{console: StreamBuffer{MaxLines: maxLines}}
}

// Append renders one Run Event into the timeline and returns the rendered
// text so follow mode can stream it. Unknown kinds and undecodable payloads
// are skipped (empty return), never fatal, so journals written by newer ACP
// Runtime versions stay viewable.
func (timeline *RunTimeline) Append(event runevent.RunEvent) string {
	text := timelineText(event)
	if text == "" {
		return ""
	}
	_, _ = timeline.console.Write([]byte(text))
	return text
}

// timelineText renders agent events from their raw payloads and daemon
// events from their bounded summaries, so the timeline narrates the whole
// loop. Anything else is an unknown kind and skipped.
func timelineText(event runevent.RunEvent) string {
	if update, ok := agent.StreamUpdateFromEvent(event); ok {
		return agent.ConsoleText(update)
	}
	if !runevent.IsDaemonKind(event.Kind) {
		return ""
	}
	summary := strings.TrimRight(event.Summary, "\r\n")
	if summary == "" {
		return ""
	}
	return summary + "\n"
}

// Lines returns the rendered timeline, oldest first, bounded by the ring.
func (timeline *RunTimeline) Lines() []string {
	return timeline.console.Lines()
}
