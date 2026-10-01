// Suite: ACP prompt identity capture.
// Boundary: fake acpx stream; no adapter or user session store.
package agent

import "testing"

func TestACPXRunPromptReportsTheACPSessionID(t *testing.T) {
	t.Parallel()
	tests := []struct{ name, stream, want string }{
		{"same id", acpxUpdateLine(`{"sessionId":"s-1","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"hello"}}}`), "s-1"},
		{"no id", acpxUpdateLine(`{"update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"hello"}}}`), ""},
		{"first id wins", acpxUpdateLine(`{"sessionId":"s-1","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"hello"}}}`) + acpxUpdateLine(`{"sessionId":"s-2","update":{"sessionUpdate":"agent_thought_chunk","content":{"type":"text","text":"thinking"}}}`), "s-1"},
		{"first notification lacks id", acpxUpdateLine(`{"update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"hello"}}}`) + acpxUpdateLine(`{"sessionId":"s-2","update":{"sessionUpdate":"agent_thought_chunk","content":{"type":"text","text":"thinking"}}}`), ""},
		{"no notifications", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			stdout := tt.stream + acpxPromptResponseLine("end_turn")
			run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout})
			if run.err != nil {
				t.Fatal(run.err)
			}
			if run.result.ACPSessionID != tt.want {
				t.Fatalf("session id=%q want=%q", run.result.ACPSessionID, tt.want)
			}
			if run.result.Output != stdout || run.result.StopReason != "end_turn" {
				t.Fatalf("existing result fields changed: %+v", run.result)
			}
		})
	}
}
