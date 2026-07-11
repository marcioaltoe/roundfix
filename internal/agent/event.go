package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"roundfix/internal/runevent"
)

// statusPayload and rawPayload are the runner-defined small JSON payloads
// for events that have no ACP session update behind them.
type statusPayload struct {
	Status string `json:"status"`
}

type rawPayload struct {
	Text string `json:"text"`
}

func runEventKind(kind StreamUpdateKind) runevent.Kind {
	switch kind {
	case StreamUpdateMessage:
		return runevent.KindAgentMessage
	case StreamUpdateThought:
		return runevent.KindAgentThought
	case StreamUpdateToolStarted:
		return runevent.KindAgentToolStarted
	case StreamUpdateToolUpdated:
		return runevent.KindAgentToolUpdated
	case StreamUpdatePlan:
		return runevent.KindAgentPlan
	case StreamUpdateStatus:
		return runevent.KindAgentStatus
	case StreamUpdateRaw:
		return runevent.KindAgentRaw
	default:
		return ""
	}
}

func streamUpdateKind(kind runevent.Kind) StreamUpdateKind {
	switch kind {
	case runevent.KindAgentMessage:
		return StreamUpdateMessage
	case runevent.KindAgentThought:
		return StreamUpdateThought
	case runevent.KindAgentToolStarted:
		return StreamUpdateToolStarted
	case runevent.KindAgentToolUpdated:
		return StreamUpdateToolUpdated
	case runevent.KindAgentPlan:
		return StreamUpdatePlan
	case runevent.KindAgentStatus:
		return StreamUpdateStatus
	case runevent.KindAgentRaw:
		return StreamUpdateRaw
	default:
		return ""
	}
}

// newAgentRunEvent stamps Run identity from the execute request and bounds
// the summary; payload bytes pass through untouched per ADR 0008.
func newAgentRunEvent(req ExecuteRequest, update StreamUpdate, payload json.RawMessage, at time.Time) runevent.RunEvent {
	summary := ConsoleText(update)
	if summary == "" && update.Kind == StreamUpdateStatus && isSessionLifecycleStatus(update.Status) {
		summary = "SESSION " + strings.ToUpper(update.Status) + "\n"
	}
	return runevent.RunEvent{
		RunID:     req.RunID,
		Batch:     req.Batch.Number,
		Source:    runevent.SourceAgent,
		Kind:      runEventKind(update.Kind),
		ToolID:    update.ToolID,
		ToolState: update.ToolState,
		Summary:   runevent.BoundSummary(summary),
		Time:      at,
		Payload:   payload,
	}
}

func marshalStatusPayload(status string) json.RawMessage {
	payload, err := json.Marshal(statusPayload{Status: status})
	if err != nil {
		return nil
	}
	return payload
}

func marshalRawPayload(text string) json.RawMessage {
	payload, err := json.Marshal(rawPayload{Text: text})
	if err != nil {
		return nil
	}
	return payload
}

// StreamUpdateFromEvent reconstructs the Agent stream model from an
// agent-source Run Event payload. Unknown kinds and undecodable payloads
// report ok=false and must be skipped, never treated as errors.
func StreamUpdateFromEvent(event runevent.RunEvent) (StreamUpdate, bool) {
	kind := streamUpdateKind(event.Kind)
	switch kind {
	case "":
		return StreamUpdate{}, false
	case StreamUpdateStatus:
		var payload statusPayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return StreamUpdate{}, false
		}
		return StreamUpdate{Kind: StreamUpdateStatus, Status: payload.Status}, true
	case StreamUpdateRaw:
		var payload rawPayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return StreamUpdate{}, false
		}
		return StreamUpdate{Kind: StreamUpdateRaw, Text: payload.Text}, true
	default:
		update, ok, err := streamUpdateFromSessionUpdatePayload(event.Payload)
		if err != nil || !ok {
			return StreamUpdate{}, false
		}
		return update, true
	}
}

// WriterSink renders agent-source Run Events as plain text, preserving the
// non-TTY output contract that the Live Run View used to provide.
type WriterSink struct {
	Writer io.Writer
}

func (sink WriterSink) Publish(_ context.Context, event runevent.RunEvent) error {
	update, ok := StreamUpdateFromEvent(event)
	if !ok {
		return nil
	}
	text := ConsoleText(update)
	if text == "" {
		return nil
	}
	_, err := io.WriteString(sink.Writer, text)
	return err
}

// ConsoleText renders one stream update as console text. It is the shared
// text contract for the non-TTY writer path and the Live Run View timeline.
func ConsoleText(update StreamUpdate) string {
	switch update.Kind {
	case StreamUpdateMessage:
		return update.Text
	case StreamUpdateThought:
		if strings.TrimSpace(update.Text) == "" {
			return ""
		}
		return "THINK " + strings.TrimRight(update.Text, "\r\n") + "\n"
	case StreamUpdateToolStarted:
		return consoleToolStartedText(update)
	case StreamUpdateToolUpdated:
		return consoleToolUpdatedText(update)
	case StreamUpdatePlan:
		if strings.TrimSpace(update.Text) == "" {
			return ""
		}
		return "PLAN\n" + strings.TrimRight(update.Text, "\r\n") + "\n"
	case StreamUpdateStatus:
		if update.Status == "" {
			return ""
		}
		if isSessionLifecycleStatus(update.Status) {
			return ""
		}
		return "SESSION " + strings.ToUpper(update.Status) + "\n"
	case StreamUpdateRaw:
		return update.Text
	default:
		return update.Text
	}
}

func isSessionLifecycleStatus(status string) bool {
	return status == AgentSessionStartedStatus || status == AgentSessionClosedStatus
}

func consoleToolStartedText(update StreamUpdate) string {
	if lines := compactFileBlockLines(update.Blocks); len(lines) > 0 {
		return strings.Join(lines, "\n") + "\n"
	}
	title := consoleToolName(update)
	marker := "[TOOL] " + title
	lines := []string{consoleToolMarker(marker, update.ToolState)}
	lines = appendConsoleToolDetail(lines, update)
	return strings.Join(lines, "\n") + "\n"
}

func consoleToolUpdatedText(update StreamUpdate) string {
	if lines := compactFileBlockLines(update.Blocks); len(lines) > 0 {
		return strings.Join(lines, "\n") + "\n"
	}
	lines := []string{}
	if name := consoleToolName(update); name != "" {
		lines = append(lines, consoleToolMarker("[TOOL] "+name, update.ToolState))
	}
	lines = appendConsoleToolDetail(lines, update)
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

// consoleToolMarker keeps the tool call and its state on one log line.
func consoleToolMarker(marker string, state string) string {
	if strings.TrimSpace(state) == "" {
		return marker
	}
	return marker + " · " + strings.TrimSpace(state)
}

func appendConsoleToolDetail(lines []string, update StreamUpdate) []string {
	return append(lines, consoleBlockLines(update.Blocks)...)
}

func consoleBlockLines(blocks []StreamBlock) []string {
	lines := []string{}
	for _, block := range blocks {
		switch block.Kind {
		case StreamBlockText:
		case StreamBlockInput:
		case StreamBlockOutput:
		case StreamBlockDiff:
		case StreamBlockRead:
			lines = append(lines, compactReadLine(block))
		case StreamBlockEdit:
			lines = append(lines, compactEditLine(block))
		case StreamBlockTerminal:
			if block.TerminalID != "" {
				lines = append(lines, "terminal: "+block.TerminalID)
			}
		case StreamBlockImage:
			lines = append(lines, "image: "+firstNonEmpty(block.MimeType, block.URI, "image"))
		case StreamBlockResource:
			lines = append(lines, "resource: "+firstNonEmpty(block.Name, block.URI, "resource"))
		}
	}
	return lines
}

func compactFileBlockLines(blocks []StreamBlock) []string {
	lines := []string{}
	for _, block := range blocks {
		switch block.Kind {
		case StreamBlockRead:
			lines = append(lines, compactReadLine(block))
		case StreamBlockEdit:
			lines = append(lines, compactEditLine(block))
		}
	}
	return lines
}

func compactReadLine(block StreamBlock) string {
	return fmt.Sprintf("read %s (%d lines)", block.Path, block.LineCount)
}

func compactEditLine(block StreamBlock) string {
	return fmt.Sprintf("edit %s (+%d/-%d)", block.Path, block.NewLineCount, block.OldLineCount)
}

func consoleToolName(update StreamUpdate) string {
	name := firstNonEmpty(update.Title, update.ToolKind, update.ToolID, "tool call")
	if path := primaryUpdatePath(update); path != "" && !strings.Contains(name, path) {
		name += " " + path
	}
	return name
}

func primaryUpdatePath(update StreamUpdate) string {
	for _, block := range update.Blocks {
		if path := strings.TrimSpace(block.Path); path != "" {
			return path
		}
	}
	for _, location := range update.Locations {
		if path := strings.TrimSpace(location.Path); path != "" {
			return path
		}
	}
	return ""
}
