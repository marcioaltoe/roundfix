package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"roundfix/internal/runevent"

	acp "github.com/coder/acp-go-sdk"
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
	return runevent.RunEvent{
		RunID:     req.RunID,
		Batch:     req.Batch.Number,
		Source:    runevent.SourceAgent,
		Kind:      runEventKind(update.Kind),
		ToolID:    update.ToolID,
		ToolState: update.ToolState,
		Summary:   runevent.BoundSummary(ConsoleText(update)),
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
		var note acp.SessionNotification
		if err := json.Unmarshal(event.Payload, &note); err != nil {
			return StreamUpdate{}, false
		}
		update := streamUpdateFromACP(note.Update)
		if update.Kind == "" {
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
		return "SESSION " + strings.ToUpper(update.Status) + "\n"
	case StreamUpdateRaw:
		return update.Text
	default:
		return update.Text
	}
}

func consoleToolStartedText(update StreamUpdate) string {
	title := strings.TrimSpace(update.Title)
	if title == "" {
		title = "tool call"
	}
	marker := "[TOOL] " + title
	if update.ToolID != "" {
		marker = fmt.Sprintf("[TOOL] %s (%s)", title, update.ToolID)
	}
	lines := []string{consoleToolMarker(marker, update.ToolState)}
	lines = appendConsoleToolDetail(lines, update)
	return strings.Join(lines, "\n") + "\n"
}

func consoleToolUpdatedText(update StreamUpdate) string {
	lines := []string{}
	if name := firstNonEmpty(update.Title, update.ToolID); name != "" {
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
	if strings.TrimSpace(update.Text) != "" && len(update.Blocks) == 0 {
		lines = append(lines, strings.TrimRight(update.Text, "\r\n"))
	}
	return append(lines, consoleBlockLines(update.Blocks)...)
}

func consoleBlockLines(blocks []StreamBlock) []string {
	lines := []string{}
	for _, block := range blocks {
		switch block.Kind {
		case StreamBlockText:
			if text := strings.TrimRight(block.Text, "\r\n"); text != "" {
				lines = append(lines, text)
			}
		case StreamBlockInput:
			if text := strings.TrimRight(block.Text, "\r\n"); text != "" {
				lines = append(lines, "$ "+text)
			}
		case StreamBlockOutput:
			if text := strings.TrimRight(block.Text, "\r\n"); text != "" {
				lines = append(lines, strings.Split(text, "\n")...)
			}
		case StreamBlockDiff:
			if block.Path != "" {
				lines = append(lines, "diff: "+block.Path)
			}
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
