package agent

import (
	"fmt"
	"io"
	"strings"
)

type StreamUpdateKind string

const (
	StreamUpdateRaw         StreamUpdateKind = "raw"
	StreamUpdateMessage     StreamUpdateKind = "message"
	StreamUpdateThought     StreamUpdateKind = "thought"
	StreamUpdateToolStarted StreamUpdateKind = "tool_started"
	StreamUpdateToolUpdated StreamUpdateKind = "tool_updated"
	StreamUpdatePlan        StreamUpdateKind = "plan"
	StreamUpdateStatus      StreamUpdateKind = "status"
)

type StreamUpdate struct {
	Kind      StreamUpdateKind
	Title     string
	Text      string
	ToolID    string
	ToolState string
	Status    string
}

type StreamUpdateSink interface {
	HandleAgentUpdate(StreamUpdate)
}

func publishStreamUpdate(dst io.Writer, update StreamUpdate) error {
	if sink, ok := dst.(StreamUpdateSink); ok {
		sink.HandleAgentUpdate(update)
		return nil
	}
	if dst == nil {
		return nil
	}
	text := formatStreamUpdate(update)
	if text == "" {
		return nil
	}
	_, err := io.WriteString(dst, text)
	return err
}

func formatStreamUpdate(update StreamUpdate) string {
	title := strings.TrimSpace(update.Title)
	text := strings.TrimRight(update.Text, "\r\n")
	switch update.Kind {
	case StreamUpdateMessage:
		if text == "" {
			return ""
		}
		return text
	case StreamUpdateThought:
		if text == "" {
			return ""
		}
		return "Thinking: " + text + "\n"
	case StreamUpdateToolStarted:
		if title == "" {
			title = "tool call"
		}
		if update.ToolID != "" {
			return fmt.Sprintf("[TOOL] %s (%s)\n", title, update.ToolID)
		}
		return fmt.Sprintf("[TOOL] %s\n", title)
	case StreamUpdateToolUpdated:
		parts := []string{}
		if title != "" {
			parts = append(parts, title)
		}
		if update.ToolState != "" {
			parts = append(parts, update.ToolState)
		}
		if text != "" {
			parts = append(parts, text)
		}
		if len(parts) == 0 {
			return ""
		}
		return strings.Join(parts, "\n") + "\n"
	case StreamUpdatePlan:
		if text == "" {
			return ""
		}
		return "Plan:\n" + text + "\n"
	case StreamUpdateStatus:
		if update.Status == "" {
			return ""
		}
		return "Session " + update.Status + "\n"
	case StreamUpdateRaw:
		if text == "" {
			return ""
		}
		return text
	default:
		if text == "" {
			return ""
		}
		return text + "\n"
	}
}
