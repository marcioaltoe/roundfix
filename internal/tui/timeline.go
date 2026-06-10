package tui

import (
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

// Append renders one Run Event into the timeline. Unknown kinds and
// undecodable payloads are skipped, never fatal, so journals written by
// newer ACP Runtime versions stay viewable.
func (timeline *RunTimeline) Append(event runevent.RunEvent) {
	update, ok := agent.StreamUpdateFromEvent(event)
	if !ok {
		return
	}
	text := agent.ConsoleText(update)
	if text == "" {
		return
	}
	_, _ = timeline.console.Write([]byte(text))
}

// Lines returns the rendered timeline, oldest first, bounded by the ring.
func (timeline *RunTimeline) Lines() []string {
	return timeline.console.Lines()
}
