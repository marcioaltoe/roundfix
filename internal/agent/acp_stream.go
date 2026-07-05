package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const acpMethodSessionUpdate = "session/update"

type acpSessionNotificationPayload struct {
	SessionID string          `json:"sessionId"`
	Update    json.RawMessage `json:"update"`
}

type acpSessionUpdateHeader struct {
	SessionUpdate string `json:"sessionUpdate"`
}

type acpContentBlockPayload struct {
	Type         string                  `json:"type"`
	Text         string                  `json:"text"`
	MimeType     string                  `json:"mimeType"`
	URI          string                  `json:"uri"`
	ResourceLink *acpResourceLinkPayload `json:"resourceLink"`
}

type acpResourceLinkPayload struct {
	Name     string  `json:"name"`
	URI      string  `json:"uri"`
	MimeType *string `json:"mimeType"`
}

type acpToolContentPayload struct {
	Type       string                  `json:"type"`
	Text       string                  `json:"text"`
	Path       string                  `json:"path"`
	MimeType   string                  `json:"mimeType"`
	URI        string                  `json:"uri"`
	TerminalID string                  `json:"terminalId"`
	Content    *acpContentBlockPayload `json:"content"`
	Diff       *struct {
		Path string `json:"path"`
	} `json:"diff"`
	Terminal *struct {
		TerminalID string `json:"terminalId"`
	} `json:"terminal"`
	ResourceLink *acpResourceLinkPayload `json:"resourceLink"`
}

type acpPlanEntryPayload struct {
	Content string `json:"content"`
	Status  string `json:"status"`
}

func streamUpdateFromSessionUpdatePayload(payload json.RawMessage) (StreamUpdate, bool, error) {
	updatePayload, ok, err := sessionUpdatePayloadFromPayload(payload)
	if err != nil || !ok {
		return StreamUpdate{}, false, err
	}
	update, err := streamUpdateFromSessionUpdate(updatePayload)
	if err != nil {
		return StreamUpdate{}, false, err
	}
	if update.Kind == "" {
		return StreamUpdate{}, false, nil
	}
	return update, true, nil
}

func sessionUpdatePayloadFromPayload(payload json.RawMessage) (json.RawMessage, bool, error) {
	var note acpSessionNotificationPayload
	if err := json.Unmarshal(payload, &note); err == nil && len(note.Update) > 0 {
		update, err := streamUpdateFromSessionUpdate(note.Update)
		if err != nil {
			return nil, false, err
		}
		if update.Kind != "" {
			return note.Update, true, nil
		}
	}

	var message acpxJSONRPCMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, false, fmt.Errorf("parse acpx session/update JSON-RPC line: %w", err)
	}
	if message.Method != acpMethodSessionUpdate {
		return nil, false, nil
	}
	if len(message.Params) == 0 {
		return nil, false, errors.New("parse acpx session/update: params are required")
	}
	if err := json.Unmarshal(message.Params, &note); err != nil {
		return nil, false, fmt.Errorf("parse acpx session/update params: %w", err)
	}
	if len(note.Update) == 0 {
		return nil, false, errors.New("parse acpx session/update: update is required")
	}
	return note.Update, true, nil
}

func streamUpdateFromSessionUpdate(payload json.RawMessage) (StreamUpdate, error) {
	var header acpSessionUpdateHeader
	if err := json.Unmarshal(payload, &header); err != nil {
		return StreamUpdate{}, fmt.Errorf("parse ACP session update: %w", err)
	}
	switch header.SessionUpdate {
	case "agent_message_chunk":
		var update struct {
			Content acpContentBlockPayload `json:"content"`
		}
		if err := json.Unmarshal(payload, &update); err != nil {
			return StreamUpdate{}, fmt.Errorf("parse ACP message update: %w", err)
		}
		return StreamUpdate{Kind: StreamUpdateMessage, Text: contentBlockText(update.Content)}, nil
	case "agent_thought_chunk":
		var update struct {
			Content acpContentBlockPayload `json:"content"`
		}
		if err := json.Unmarshal(payload, &update); err != nil {
			return StreamUpdate{}, fmt.Errorf("parse ACP thought update: %w", err)
		}
		return StreamUpdate{Kind: StreamUpdateThought, Text: contentBlockText(update.Content)}, nil
	case "tool_call":
		var update struct {
			ToolCallID string                  `json:"toolCallId"`
			Title      string                  `json:"title"`
			Status     string                  `json:"status"`
			RawInput   any                     `json:"rawInput"`
			Content    []acpToolContentPayload `json:"content"`
			RawOutput  any                     `json:"rawOutput"`
		}
		if err := json.Unmarshal(payload, &update); err != nil {
			return StreamUpdate{}, fmt.Errorf("parse ACP tool call update: %w", err)
		}
		blocks := toolContentBlocks(update.Content, update.RawInput, update.RawOutput)
		return StreamUpdate{
			Kind:      StreamUpdateToolStarted,
			Title:     update.Title,
			ToolID:    update.ToolCallID,
			ToolState: update.Status,
			Text:      streamBlockText(blocks),
			Blocks:    blocks,
		}, nil
	case "tool_call_update":
		var update struct {
			ToolCallID string                  `json:"toolCallId"`
			Title      *string                 `json:"title"`
			Status     *string                 `json:"status"`
			RawInput   any                     `json:"rawInput"`
			Content    []acpToolContentPayload `json:"content"`
			RawOutput  any                     `json:"rawOutput"`
		}
		if err := json.Unmarshal(payload, &update); err != nil {
			return StreamUpdate{}, fmt.Errorf("parse ACP tool call delta: %w", err)
		}
		title := ""
		if update.Title != nil {
			title = *update.Title
		}
		state := ""
		if update.Status != nil {
			state = *update.Status
		}
		blocks := toolContentBlocks(update.Content, update.RawInput, update.RawOutput)
		return StreamUpdate{
			Kind:      StreamUpdateToolUpdated,
			Title:     title,
			ToolID:    update.ToolCallID,
			ToolState: state,
			Text:      streamBlockText(blocks),
			Blocks:    blocks,
		}, nil
	case "plan":
		var update struct {
			Entries []acpPlanEntryPayload `json:"entries"`
		}
		if err := json.Unmarshal(payload, &update); err != nil {
			return StreamUpdate{}, fmt.Errorf("parse ACP plan update: %w", err)
		}
		lines := make([]string, 0, len(update.Entries))
		for _, entry := range update.Entries {
			line := strings.TrimSpace(entry.Content)
			if line == "" {
				continue
			}
			if entry.Status != "" {
				line = fmt.Sprintf("%s  %s", entry.Status, line)
			}
			lines = append(lines, line)
		}
		return StreamUpdate{Kind: StreamUpdatePlan, Text: strings.Join(lines, "\n")}, nil
	default:
		return StreamUpdate{}, nil
	}
}

func contentBlockText(block acpContentBlockPayload) string {
	if block.Type == "" && block.Text == "" && block.MimeType == "" && block.URI == "" && block.ResourceLink == nil {
		return ""
	}
	if block.Type == "text" {
		return block.Text
	}
	raw, err := json.Marshal(block)
	if err != nil {
		return ""
	}
	return string(raw)
}

func toolContentBlocks(content []acpToolContentPayload, rawInput any, rawOutput any) []StreamBlock {
	blocks := []StreamBlock{}
	if rawInput != nil {
		blocks = append(blocks, StreamBlock{Kind: StreamBlockInput, Text: formatToolInput(rawInput)})
	}
	for _, item := range content {
		switch {
		case item.Content != nil && item.Content.Type == "text":
			blocks = append(blocks, StreamBlock{Kind: StreamBlockText, Text: item.Content.Text})
		case item.Content != nil && item.Content.Type == "image":
			blocks = append(blocks, StreamBlock{Kind: StreamBlockImage, MimeType: item.Content.MimeType, URI: item.Content.URI})
		case item.Content != nil && item.Content.ResourceLink != nil:
			blocks = append(blocks, resourceBlock(item.Content.ResourceLink))
		case item.Type == "text":
			blocks = append(blocks, StreamBlock{Kind: StreamBlockText, Text: item.Text})
		case item.Type == "image":
			blocks = append(blocks, StreamBlock{Kind: StreamBlockImage, MimeType: item.MimeType, URI: item.URI})
		case item.ResourceLink != nil:
			blocks = append(blocks, resourceBlock(item.ResourceLink))
		case item.Diff != nil:
			blocks = append(blocks, StreamBlock{Kind: StreamBlockDiff, Path: item.Diff.Path})
		case item.Type == "diff":
			blocks = append(blocks, StreamBlock{Kind: StreamBlockDiff, Path: item.Path})
		case item.Terminal != nil:
			blocks = append(blocks, StreamBlock{Kind: StreamBlockTerminal, TerminalID: item.Terminal.TerminalID})
		case item.Type == "terminal":
			blocks = append(blocks, StreamBlock{Kind: StreamBlockTerminal, TerminalID: firstNonEmpty(item.TerminalID, item.URI)})
		}
	}
	if rawOutput != nil {
		blocks = append(blocks, StreamBlock{Kind: StreamBlockOutput, Text: formatToolOutput(rawOutput)})
	}
	return blocks
}

func resourceBlock(link *acpResourceLinkPayload) StreamBlock {
	return StreamBlock{
		Kind:     StreamBlockResource,
		Name:     link.Name,
		URI:      link.URI,
		MimeType: stringValue(link.MimeType),
	}
}

// formatToolInput renders tool input as a log line, not raw JSON: shell
// invocations show the command itself; anything else falls back to compact
// JSON.
func formatToolInput(rawInput any) string {
	fields, ok := rawInput.(map[string]any)
	if !ok {
		return marshalCompact(rawInput)
	}
	if command := shellCommandFromInput(fields["command"]); command != "" {
		return command
	}
	return marshalCompact(rawInput)
}

// shellCommandFromInput extracts the human command from exec-style inputs
// like ["/bin/zsh","-lc","rtk go test ./..."].
func shellCommandFromInput(value any) string {
	switch command := value.(type) {
	case string:
		return strings.TrimSpace(command)
	case []any:
		parts := make([]string, 0, len(command))
		for _, item := range command {
			text, ok := item.(string)
			if !ok {
				return ""
			}
			parts = append(parts, text)
		}
		if len(parts) == 0 {
			return ""
		}
		last := strings.TrimSpace(parts[len(parts)-1])
		if len(parts) >= 3 && strings.HasPrefix(parts[1], "-") && last != "" {
			return last
		}
		return strings.Join(parts, " ")
	default:
		return ""
	}
}

// toolOutputMaxLines bounds rendered tool output so chatty commands stay
// readable; the raw payload keeps the full data.
const toolOutputMaxLines = 8

// formatToolOutput renders tool output as log text: exec-style payloads show
// their aggregated output, bounded; anything else falls back to compact JSON.
func formatToolOutput(rawOutput any) string {
	fields, ok := rawOutput.(map[string]any)
	if !ok {
		return marshalCompact(rawOutput)
	}
	for _, key := range []string{"aggregated_output", "formatted_output", "stdout", "output"} {
		text, found := fields[key].(string)
		if !found {
			continue
		}
		return boundOutputLines(text)
	}
	return marshalCompact(rawOutput)
}

func boundOutputLines(text string) string {
	trimmed := strings.TrimRight(text, "\r\n")
	if trimmed == "" {
		return ""
	}
	lines := strings.Split(trimmed, "\n")
	if len(lines) <= toolOutputMaxLines {
		return trimmed
	}
	kept := lines[:toolOutputMaxLines]
	return strings.Join(kept, "\n") + fmt.Sprintf("\n… (+%d more line(s))", len(lines)-toolOutputMaxLines)
}

func streamBlockText(blocks []StreamBlock) string {
	parts := []string{}
	for _, block := range blocks {
		switch block.Kind {
		case StreamBlockText:
			if strings.TrimSpace(block.Text) != "" {
				parts = append(parts, block.Text)
			}
		case StreamBlockInput:
			if strings.TrimSpace(block.Text) != "" {
				parts = append(parts, "input: "+block.Text)
			}
		case StreamBlockOutput:
			if strings.TrimSpace(block.Text) != "" {
				parts = append(parts, "output: "+block.Text)
			}
		case StreamBlockDiff:
			if block.Path != "" {
				parts = append(parts, "diff: "+block.Path)
			}
		case StreamBlockTerminal:
			if block.TerminalID != "" {
				parts = append(parts, "terminal: "+block.TerminalID)
			}
		case StreamBlockImage:
			parts = append(parts, "image: "+firstNonEmpty(block.MimeType, block.URI, "image"))
		case StreamBlockResource:
			parts = append(parts, "resource: "+firstNonEmpty(block.Name, block.URI, "resource"))
		}
	}
	return strings.Join(parts, "\n")
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func marshalCompact(value any) string {
	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(raw)
}
