// Suite: ACP Agent message boundaries
// Invariant: Agent chunks preserve message boundaries and expose the last non-blank message as the answer.
// Boundary IN: parsed StreamUpdates, prompt stream aggregation, and sealed prompt extraction.
// Boundary OUT: review verdict semantics and real ACP adapters.

package agent

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestAgentMessagesSplitAtAThoughtOrToolCallWithoutMessageIDs(t *testing.T) {
	tests := []struct {
		name     string
		boundary StreamUpdate
	}{
		{name: "thought", boundary: StreamUpdate{Kind: StreamUpdateThought, Text: "checking"}},
		{name: "tool call", boundary: StreamUpdate{Kind: StreamUpdateToolStarted, ToolID: "tool-1"}},
		{name: "tool update", boundary: StreamUpdate{Kind: StreamUpdateToolUpdated, ToolID: "tool-1"}},
		{name: "plan", boundary: StreamUpdate{Kind: StreamUpdatePlan, Text: "inspect"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var log agentMessageLog
			log.observe(StreamUpdate{Kind: StreamUpdateMessage, Text: "progress"})
			log.observe(test.boundary)
			log.observe(StreamUpdate{Kind: StreamUpdateMessage, Text: "answer"})

			want := []string{"progress", "answer"}
			if got := log.messages(); !reflect.DeepEqual(got, want) {
				t.Fatalf("messages = %#v, want %#v", got, want)
			}
		})
	}
}

func TestAgentMessagesFollowTheirMessageIDs(t *testing.T) {
	t.Run("different identifiers split consecutive chunks", func(t *testing.T) {
		var log agentMessageLog
		log.observe(StreamUpdate{Kind: StreamUpdateMessage, MessageID: "message-1", Text: "progress"})
		log.observe(StreamUpdate{Kind: StreamUpdateMessage, MessageID: "message-2", Text: "answer"})

		want := []string{"progress", "answer"}
		if got := log.messages(); !reflect.DeepEqual(got, want) {
			t.Fatalf("messages = %#v, want %#v", got, want)
		}
	})

	t.Run("same identifier appends across a tool call", func(t *testing.T) {
		var log agentMessageLog
		log.observe(StreamUpdate{Kind: StreamUpdateMessage, MessageID: "message-1", Text: "first"})
		log.observe(StreamUpdate{Kind: StreamUpdateToolStarted, ToolID: "tool-1"})
		log.observe(StreamUpdate{Kind: StreamUpdateMessage, MessageID: "message-1", Text: " second"})

		want := []string{"first second"}
		if got := log.messages(); !reflect.DeepEqual(got, want) {
			t.Fatalf("messages = %#v, want %#v", got, want)
		}
	})
}

func TestAgentMessageChunksWithoutABoundaryStayOneMessage(t *testing.T) {
	tests := []struct {
		name    string
		between []StreamUpdate
	}{
		{name: "consecutive chunks"},
		{name: "status update", between: []StreamUpdate{{Kind: StreamUpdateStatus, Status: "running"}}},
		{name: "raw update", between: []StreamUpdate{{Kind: StreamUpdateRaw, Text: "metadata"}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var log agentMessageLog
			log.observe(StreamUpdate{Kind: StreamUpdateMessage, Text: "first"})
			for _, update := range test.between {
				log.observe(update)
			}
			log.observe(StreamUpdate{Kind: StreamUpdateMessage, Text: " second"})

			want := []string{"first second"}
			if got := log.messages(); !reflect.DeepEqual(got, want) {
				t.Fatalf("messages = %#v, want %#v", got, want)
			}
			if got := log.final(); got != want[0] {
				t.Fatalf("final message = %q, want %q", got, want[0])
			}
			if got := log.totalBytes(); got != len(want[0]) {
				t.Fatalf("total message bytes = %d, want %d", got, len(want[0]))
			}
		})
	}
}

func TestACPMessageChunkReadsOptionalMessageID(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantID  string
	}{
		{
			name:    "present",
			payload: `{"sessionUpdate":"agent_message_chunk","messageId":"message-1","content":{"type":"text","text":"hello"}}`,
			wantID:  "message-1",
		},
		{
			name:    "absent",
			payload: `{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"hello"}}`,
		},
		{
			name:    "null",
			payload: `{"sessionUpdate":"agent_message_chunk","messageId":null,"content":{"type":"text","text":"hello"}}`,
		},
		{
			name:    "ignored on another update kind",
			payload: `{"sessionUpdate":"agent_thought_chunk","messageId":"thought-1","content":{"type":"text","text":"thinking"}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			update, err := streamUpdateFromSessionUpdate([]byte(test.payload))
			if err != nil {
				t.Fatalf("parse session update: %v", err)
			}
			if update.MessageID != test.wantID {
				t.Fatalf("message id = %q, want %q", update.MessageID, test.wantID)
			}
		})
	}
}

func TestExecuteResultAnswerIsTheLastNonBlankMessage(t *testing.T) {
	tests := []struct {
		name   string
		result ExecuteResult
		want   string
	}{
		{
			name: "last non-blank message",
			result: ExecuteResult{
				Message:  "progress\n\nanswer\n\n   ",
				Messages: []string{"progress", "answer", " \n\t"},
			},
			want: "answer",
		},
		{
			name:   "legacy message fallback",
			result: ExecuteResult{Message: "legacy answer"},
			want:   "legacy answer",
		},
		{
			name:   "all messages blank",
			result: ExecuteResult{Message: "joined transcript", Messages: []string{"", " \n"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.result.Answer(); got != test.want {
				t.Fatalf("Answer() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestACPXRunPromptReportsEachAgentMessage(t *testing.T) {
	t.Run("messages separated by a tool call", func(t *testing.T) {
		stdout := acpxUpdateLine(`{"sessionId":"sess-1","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"progress"}}}`) +
			acpxUpdateLine(`{"sessionId":"sess-1","update":{"sessionUpdate":"tool_call","toolCallId":"tool-1","kind":"read","title":"read","status":"pending"}}`) +
			acpxUpdateLine(`{"sessionId":"sess-1","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"answer"}}}`) +
			acpxPromptResponseLine("end_turn")

		run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout})
		if run.err != nil {
			t.Fatalf("RunPrompt() error = %v", run.err)
		}
		wantMessages := []string{"progress", "answer"}
		if !reflect.DeepEqual(run.result.Messages, wantMessages) {
			t.Fatalf("Messages = %#v, want %#v", run.result.Messages, wantMessages)
		}
		if want := "progress\n\nanswer"; run.result.Message != want {
			t.Fatalf("Message = %q, want %q", run.result.Message, want)
		}
	})

	t.Run("one message remains byte-identical", func(t *testing.T) {
		const message = "first chunk second chunk"
		stdout := acpxUpdateLine(`{"sessionId":"sess-1","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"first chunk"}}}`) +
			acpxUpdateLine(`{"sessionId":"sess-1","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":" second chunk"}}}`) +
			acpxPromptResponseLine("end_turn")

		run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout})
		if run.err != nil {
			t.Fatalf("RunPrompt() error = %v", run.err)
		}
		if want := []string{message}; !reflect.DeepEqual(run.result.Messages, want) {
			t.Fatalf("Messages = %#v, want %#v", run.result.Messages, want)
		}
		if run.result.Message != message {
			t.Fatalf("Message = %q, want %q", run.result.Message, message)
		}
	})
}

func TestSealedPromptOutputIsTheFinalMessage(t *testing.T) {
	tests := []struct {
		name     string
		boundary string
	}{
		{
			name:     "thought boundary",
			boundary: `{"sessionId":"sealed","update":{"sessionUpdate":"agent_thought_chunk","content":{"type":"text","text":"checking"}}}`,
		},
		{
			name:     "plan boundary",
			boundary: `{"sessionId":"sealed","update":{"sessionUpdate":"plan","entries":[{"content":"inspect","status":"in_progress"}]}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stream := acpxUpdateLine(`{"sessionId":"sealed","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"I am checking."}}}`) +
				acpxUpdateLine(test.boundary) +
				acpxUpdateLine(`{"sessionId":"sealed","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"{\"ok\":true}"}}}`) +
				acpxPromptResponseLine("end_turn")

			result, err := parseSealedPromptStream([]byte(stream))
			if err != nil {
				t.Fatalf("parse sealed stream: %v", err)
			}
			if got, want := string(result.Output), `{"ok":true}`; got != want {
				t.Fatalf("sealed output = %q, want %q", got, want)
			}
		})
	}
}

func TestSealedPromptRejectsTotalMessageBytesOverLimit(t *testing.T) {
	commentary := strings.Repeat("x", SealedPromptMaxOutputBytes)
	stream := acpxUpdateLine(fmt.Sprintf(
		`{"sessionId":"sealed","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":%q}}}`,
		commentary,
	)) +
		acpxUpdateLine(`{"sessionId":"sealed","update":{"sessionUpdate":"agent_thought_chunk","content":{"type":"text","text":"checking"}}}`) +
		acpxUpdateLine(`{"sessionId":"sealed","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"{}"}}}`) +
		acpxPromptResponseLine("end_turn")

	_, err := parseSealedPromptStream([]byte(stream))
	if !errors.Is(err, ErrSealedOutputTooLarge) {
		t.Fatalf("sealed stream error = %T %v, want ErrSealedOutputTooLarge", err, err)
	}
}

func TestSealedPromptRejectsToolUpdate(t *testing.T) {
	stream := acpxUpdateLine(`{"sessionId":"sealed","update":{"sessionUpdate":"tool_call_update","toolCallId":"tool-1","status":"completed"}}`) +
		acpxPromptResponseLine("end_turn")

	result, err := parseSealedPromptStream([]byte(stream))
	if !errors.Is(err, ErrSealedToolUse) || !result.ToolUsed {
		t.Fatalf("sealed tool update result=%+v error=%T %v, want ToolUsed and ErrSealedToolUse", result, err, err)
	}
}
