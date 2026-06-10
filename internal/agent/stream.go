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
	Blocks    []StreamBlock
	ToolID    string
	ToolState string
	Status    string
}

type StreamBlockKind string

const (
	StreamBlockText     StreamBlockKind = "text"
	StreamBlockInput    StreamBlockKind = "input"
	StreamBlockOutput   StreamBlockKind = "output"
	StreamBlockDiff     StreamBlockKind = "diff"
	StreamBlockTerminal StreamBlockKind = "terminal"
	StreamBlockImage    StreamBlockKind = "image"
	StreamBlockResource StreamBlockKind = "resource"
)

type StreamBlock struct {
	Kind       StreamBlockKind
	Text       string
	Path       string
	TerminalID string
	MimeType   string
	URI        string
	Name       string
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
		lines := []string{}
		if update.ToolID != "" {
			lines = append(lines, fmt.Sprintf("[TOOL] %s (%s)", title, update.ToolID))
		} else {
			lines = append(lines, fmt.Sprintf("[TOOL] %s", title))
		}
		if update.ToolState != "" {
			lines = append(lines, update.ToolState)
		}
		if text != "" && len(update.Blocks) == 0 {
			lines = append(lines, text)
		}
		lines = append(lines, formatStreamBlocks(update.Blocks)...)
		return strings.Join(lines, "\n") + "\n"
	case StreamUpdateToolUpdated:
		parts := []string{}
		if title != "" {
			parts = append(parts, "[TOOL] "+title)
		} else if update.ToolID != "" {
			parts = append(parts, "[TOOL] "+update.ToolID)
		}
		if update.ToolState != "" {
			parts = append(parts, update.ToolState)
		}
		if text != "" && len(update.Blocks) == 0 {
			parts = append(parts, text)
		}
		parts = append(parts, formatStreamBlocks(update.Blocks)...)
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

func formatStreamBlocks(blocks []StreamBlock) []string {
	lines := []string{}
	for _, block := range blocks {
		switch block.Kind {
		case StreamBlockText:
			if text := strings.TrimRight(block.Text, "\r\n"); text != "" {
				lines = append(lines, text)
			}
		case StreamBlockInput:
			if text := strings.TrimRight(block.Text, "\r\n"); text != "" {
				lines = append(lines, "input: "+text)
			}
		case StreamBlockOutput:
			if text := strings.TrimRight(block.Text, "\r\n"); text != "" {
				lines = append(lines, "output: "+text)
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
			label := firstNonEmpty(block.MimeType, "image")
			if block.URI != "" {
				lines = append(lines, fmt.Sprintf("image: %s %s", label, block.URI))
			} else {
				lines = append(lines, "image: "+label)
			}
		case StreamBlockResource:
			label := firstNonEmpty(block.Name, block.URI, "resource")
			lines = append(lines, "resource: "+label)
		}
	}
	return lines
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
