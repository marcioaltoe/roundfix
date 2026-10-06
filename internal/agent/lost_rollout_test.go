package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

// Recorded acpx 0.19.4 stdout from the Spec 0240 investigation appendix.
const lostResumeStdout = `{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":1,"clientCapabilities":{"fs":{"readTextFile":true,"writeTextFile":true},"terminal":true},"clientInfo":{"name":"acpx","version":"0.19.4"}}}
{"jsonrpc":"2.0","id":0,"result":{"protocolVersion":1,"agentCapabilities":{"loadSession":true,"sessionCapabilities":{"resume":{}}},"agentInfo":{"name":"fake-codex","version":"1"},"authMethods":[]}}
{"jsonrpc":"2.0","id":1,"method":"session/resume","params":{"sessionId":"fake-thread-0240","cwd":"/tmp/w","mcpServers":[]}}
{"jsonrpc":"2.0","id":1,"error":{"code":-32603,"message":"Internal error","data":{"details":"no rollout found for thread id 01a1fake-0000-7000-8000-000000000240"}}}
`
const lostPromptStdout = `{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":1,"clientCapabilities":{"fs":{"readTextFile":true,"writeTextFile":true},"terminal":true},"clientInfo":{"name":"acpx","version":"0.19.4"}}}
{"jsonrpc":"2.0","id":0,"result":{"protocolVersion":1,"agentCapabilities":{"loadSession":true,"sessionCapabilities":{"resume":{}}},"agentInfo":{"name":"fake-codex","version":"1"},"authMethods":[]}}
{"jsonrpc":"2.0","id":1,"method":"session/resume","params":{"sessionId":"fake-thread-0240","cwd":"/tmp/w","mcpServers":[]}}
{"jsonrpc":"2.0","id":1,"result":{"sessionId":"fake-thread-0240","models":{"currentModelId":"default","availableModels":[{"modelId":"default","name":"Default"}]}}}
{"jsonrpc":"2.0","id":2,"method":"session/prompt","params":{"sessionId":"fake-thread-0240","prompt":[{"type":"text","text":"second turn"}]}}
{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"fake-thread-0240","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"working"}}}}
{"jsonrpc":"2.0","id":2,"error":{"code":-32603,"message":"Internal error","data":{"details":"no rollout found for thread id 01a1fake-0000-7000-8000-000000000240"}}}
`

func TestLostRolloutIsNamedDuringResume(t *testing.T) {
	run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: lostResumeStdout, exitCode: 1})
	var selection *SelectionFailureError
	if !errors.As(run.err, &selection) || selection.Protocol == nil {
		t.Fatalf("want selection failure: %v", run.err)
	}
	assertLostRollout(t, run.err, "codex", "session/resume", false)
	want := `Agent Selection failed for runtime "codex": agent/protocol error at session/resume: Internal error (JSON-RPC -32603); lost rollout: no rollout found for thread id 01a1fake-0000-7000-8000-000000000240`
	if run.err.Error() != want {
		t.Fatalf("want %q, got %q", want, run.err.Error())
	}
	if selection.Reason != acpxExitReasonAgentProtocol || selection.Err != nil {
		t.Fatalf("changed mapping: %+v", selection)
	}
}

func TestLostRolloutIsNamedMidTurn(t *testing.T) {
	run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: lostPromptStdout, stderr: "stderr cause", exitCode: 1})
	var batch *BatchFailureError
	if !errors.As(run.err, &batch) || batch.Protocol == nil {
		t.Fatalf("want batch failure: %v", run.err)
	}
	assertLostRollout(t, run.err, "", "session/prompt", true)
	want := `Agent Batch failed after acpx exited with code 1: agent/protocol error at session/prompt: Internal error (JSON-RPC -32603); lost rollout: no rollout found for thread id 01a1fake-0000-7000-8000-000000000240: stderr cause`
	if run.err.Error() != want {
		t.Fatalf("want %q, got %q", want, run.err.Error())
	}
	if batch.ExitCode != 1 || batch.Reason != acpxExitReasonAgentProtocol || batch.Stderr != "stderr cause" || batch.Err != nil {
		t.Fatalf("changed mapping: %+v", batch)
	}
}

func assertLostRollout(t *testing.T, err error, runtime, step string, sent bool) {
	t.Helper()
	description, ok := DescribeLostRollout(fmt.Errorf("outer: %w", err))
	if !ok || description.Runtime != runtime || description.Step != step || description.PromptSent != sent || !strings.Contains(description.Detail, LostRolloutPhrase) {
		t.Fatalf("description: %+v, %t", description, ok)
	}
	if !strings.Contains(err.Error(), "); lost rollout: "+description.Detail) {
		t.Fatalf("lost detail missing from text: %v", err)
	}
}

func lostRolloutErrorLine(t *testing.T, code int, message string, data any) string {
	t.Helper()
	line, err := json.Marshal(map[string]any{"id": 3, "error": map[string]any{"code": code, "message": message, "data": data}})
	if err != nil {
		t.Fatal(err)
	}
	return string(line) + "\n"
}

func TestOtherErrorsAreNotLostRollouts(t *testing.T) {
	for _, tc := range []struct {
		name    string
		code    int
		message string
		data    any
	}{
		{"other detail", -32603, "Internal error", map[string]any{"details": "other cause"}},
		{"other code", -32602, "Internal error", map[string]any{"details": LostRolloutPhrase}},
		{"object detail", -32603, "Internal error", map[string]any{"details": map[string]any{"message": LostRolloutPhrase}}},
		{"other field", -32603, "Internal error", map[string]any{"secret": LostRolloutPhrase}},
		{"non object data", -32603, "Internal error", LostRolloutPhrase},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, output := range []bool{false, true} {
				stdout := protocolPrompt
				if output {
					stdout += acpxUpdateLine(`{"sessionId":"test-session","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"working"}}}`)
				}
				run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout + lostRolloutErrorLine(t, tc.code, tc.message, tc.data), exitCode: 1})
				if description, ok := DescribeLostRollout(run.err); ok {
					t.Fatalf("unexpected loss: %+v", description)
				}
				protocol := protocolFailureInChain(run.err)
				if protocol == nil || protocol.Detail != "" || protocol.LostRollout {
					t.Fatalf("unexpected detail: %+v", protocol)
				}
				prefix := `Agent Selection failed for runtime "codex": agent/protocol error`
				if output {
					prefix = "Agent Batch failed after acpx exited with code 1: agent/protocol error"
				}
				want := fmt.Sprintf("%s at session/prompt: Internal error (JSON-RPC %d)", prefix, tc.code)
				if run.err.Error() != want {
					t.Fatalf("want %q, got %q", want, run.err.Error())
				}
			}
		})
	}
	t.Run("phrase in message", func(t *testing.T) {
		run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: protocolPrompt + lostRolloutErrorLine(t, -32600, LostRolloutPhrase+" example", nil), exitCode: 1})
		assertLostRollout(t, run.err, "codex", "session/prompt", true)
	})
	for _, err := range []error{nil, errors.New(LostRolloutPhrase), context.Canceled, context.DeadlineExceeded, &StopError{}, &SelectionFailureError{}, &BatchFailureError{}} {
		if description, ok := DescribeLostRollout(err); ok {
			t.Fatalf("unexpected loss for %T: %+v", err, description)
		}
	}
}

func TestLostRolloutDetailIsBounded(t *testing.T) {
	detail := "  " + LostRolloutPhrase + "\n\t" + strings.Repeat("界", ProtocolFailureMessageLimit)
	run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: protocolPrompt + lostRolloutErrorLine(t, -32603, "Internal error", map[string]any{"details": detail}), exitCode: 1})
	assertLostRollout(t, run.err, "codex", "session/prompt", true)
	description, _ := DescribeLostRollout(run.err)
	prefix := LostRolloutPhrase + " "
	want := prefix + strings.Repeat("界", (ProtocolFailureMessageLimit-len(prefix))/3)
	if description.Detail != want || !utf8.ValidString(description.Detail) || strings.ContainsAny(description.Detail, "\n\r\t") {
		t.Fatalf("unbounded detail: %q", description.Detail)
	}
}

func TestSessionResumeIsATrackedStep(t *testing.T) {
	stdout := strings.Replace(lostResumeStdout, "no rollout found for thread id 01a1fake-0000-7000-8000-000000000240", "other cause", 1)
	failure := protocolSelectionFailure(t, stdout, "")
	if failure.Protocol.Step != "session/resume" || failure.Protocol.PromptSent || failure.Protocol.Detail != "" || failure.Protocol.LostRollout {
		t.Fatalf("resume detail: %+v", failure.Protocol)
	}
	want := `Agent Selection failed for runtime "codex": agent/protocol error at session/resume: Internal error (JSON-RPC -32603)`
	if failure.Error() != want {
		t.Fatalf("want %q, got %q", want, failure.Error())
	}
}

func TestDescribeLostRolloutTraversesFailureChains(t *testing.T) {
	run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: lostPromptStdout, exitCode: 1})
	selection := &SelectionFailureError{Runtime: "codex", Err: fmt.Errorf("inner: %w", run.err)}
	for _, err := range []error{
		selection,
		fmt.Errorf("outer: %w", selection),
		errors.Join(&BatchFailureError{Protocol: &ProtocolFailure{Message: "other"}}, selection),
	} {
		assertLostRollout(t, err, "codex", "session/prompt", true)
	}
	for _, err := range []error{
		errors.Join(context.Canceled, selection),
		errors.Join(context.DeadlineExceeded, selection),
		errors.Join(&StopError{}, selection),
		(*SelectionFailureError)(nil),
		(*BatchFailureError)(nil),
	} {
		if description, ok := DescribeLostRollout(err); ok {
			t.Fatalf("unexpected loss for %T: %+v", err, description)
		}
	}
}
