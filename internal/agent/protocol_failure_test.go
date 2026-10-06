package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

// Recorded from acpx 0.19.4 prompt on 2026-10-06 with a local fake ACP
// adapter and disposable subprocess home. Only the temporary cwd is normalized
// to /tmp. Data and request content are inert leak-detection markers.
const protocolInitialize = `{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":1,"clientCapabilities":{"fs":{"readTextFile":true,"writeTextFile":true},"terminal":true},"clientInfo":{"name":"acpx","version":"0.19.4"}}}` + "\n"
const protocolInitialized = `{"jsonrpc":"2.0","id":0,"result":{"protocolVersion":1,"agentCapabilities":{"loadSession":true},"agentInfo":{"name":"fixture","version":"1"},"authMethods":[]}}` + "\n"
const protocolLoad = `{"jsonrpc":"2.0","id":1,"method":"session/load","params":{"sessionId":"test-session","cwd":"/tmp","mcpServers":[],"_meta":{"claudeCode":{"options":{"model":"request-secret"}}}}}` + "\n"
const protocolLoaded = `{"jsonrpc":"2.0","id":1,"result":{"sessionId":"test-session","models":{"currentModelId":"default","availableModels":[{"modelId":"default","name":"Default"},{"modelId":"request-secret","name":"Requested"}]}}}` + "\n"
const protocolSetModel = `{"jsonrpc":"2.0","id":2,"method":"session/set_model","params":{"sessionId":"test-session","modelId":"request-secret"}}` + "\n"
const protocolModelSet = `{"jsonrpc":"2.0","id":2,"result":{}}` + "\n"
const protocolPrompt = `{"jsonrpc":"2.0","id":3,"method":"session/prompt","params":{"sessionId":"test-session","prompt":[{"type":"text","text":"prompt-secret"}]}}` + "\n"

const protocolInitializeFailure = `{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":1,"clientCapabilities":{"fs":{"readTextFile":true,"writeTextFile":true},"terminal":true},"clientInfo":{"name":"acpx","version":"0.19.4"}}}
{"jsonrpc":"2.0","id":0,"error":{"code":-32603,"message":"adapter refused the operation","data":{"secret":"data-secret"}}}
`

const protocolModelFailure = `{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":1,"clientCapabilities":{"fs":{"readTextFile":true,"writeTextFile":true},"terminal":true},"clientInfo":{"name":"acpx","version":"0.19.4"}}}
{"jsonrpc":"2.0","id":0,"result":{"protocolVersion":1,"agentCapabilities":{"loadSession":true},"agentInfo":{"name":"fixture","version":"1"},"authMethods":[]}}
{"jsonrpc":"2.0","id":1,"method":"session/load","params":{"sessionId":"test-session","cwd":"/tmp","mcpServers":[],"_meta":{"claudeCode":{"options":{"model":"request-secret"}}}}}
{"jsonrpc":"2.0","id":1,"result":{"sessionId":"test-session","models":{"currentModelId":"default","availableModels":[{"modelId":"default","name":"Default"},{"modelId":"request-secret","name":"Requested"}]}}}
{"jsonrpc":"2.0","id":2,"method":"session/set_model","params":{"sessionId":"test-session","modelId":"request-secret"}}
{"jsonrpc":"2.0","id":2,"error":{"code":-32603,"message":"adapter refused the operation","data":{"secret":"data-secret"}}}
`

const protocolPromptFailure = `{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":1,"clientCapabilities":{"fs":{"readTextFile":true,"writeTextFile":true},"terminal":true},"clientInfo":{"name":"acpx","version":"0.19.4"}}}
{"jsonrpc":"2.0","id":0,"result":{"protocolVersion":1,"agentCapabilities":{"loadSession":true},"agentInfo":{"name":"fixture","version":"1"},"authMethods":[]}}
{"jsonrpc":"2.0","id":1,"method":"session/load","params":{"sessionId":"test-session","cwd":"/tmp","mcpServers":[],"_meta":{"claudeCode":{"options":{"model":"request-secret"}}}}}
{"jsonrpc":"2.0","id":1,"result":{"sessionId":"test-session","models":{"currentModelId":"default","availableModels":[{"modelId":"default","name":"Default"},{"modelId":"request-secret","name":"Requested"}]}}}
{"jsonrpc":"2.0","id":2,"method":"session/set_model","params":{"sessionId":"test-session","modelId":"request-secret"}}
{"jsonrpc":"2.0","id":2,"result":{}}
{"jsonrpc":"2.0","id":3,"method":"session/prompt","params":{"sessionId":"test-session","prompt":[{"type":"text","text":"prompt-secret"}]}}
{"jsonrpc":"2.0","id":3,"error":{"code":-32603,"message":"adapter refused the operation","data":{"secret":"data-secret"}}}
`

const protocolDisconnect = protocolInitialize + protocolInitialized + protocolLoad + protocolLoaded + protocolSetModel + protocolModelSet + protocolPrompt + `{"jsonrpc":"2.0","id":null,"error":{"code":-32603,"message":"ACP agent disconnected during request (connection_close, exit=null, signal=null)","data":{"acpxCode":"RUNTIME","detailCode":"AGENT_DISCONNECTED","origin":"acp","sessionId":"c3548ce1-4860-44a2-af42-f64fd84fb9ce"}}}` + "\n"

func protocolErrorLine(t *testing.T, id any, message string) string {
	t.Helper()
	line, err := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": id, "error": map[string]any{"code": -32603, "message": message, "data": map[string]any{"acpxCode": "AGENT_DISCONNECTED", "secret": "data-secret"}}})
	if err != nil {
		t.Fatal(err)
	}
	return string(line) + "\n"
}

func protocolSelectionFailure(t *testing.T, stdout, stderr string) *SelectionFailureError {
	t.Helper()
	run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout, stderr: stderr, exitCode: 1})
	var failure *SelectionFailureError
	if !errors.As(run.err, &failure) {
		t.Fatalf("want selection failure, got %T: %v", run.err, run.err)
	}
	if failure.Protocol == nil {
		t.Fatal("missing protocol detail")
	}
	return failure
}

func TestACPXPromptNamesTheFailingProtocolStep(t *testing.T) {
	seen := map[string]bool{}
	for _, tc := range []struct {
		step, stdout string
	}{
		{"initialize", protocolInitializeFailure},
		{"session/set_model", protocolModelFailure},
		{"session/prompt", protocolPromptFailure},
	} {
		t.Run(tc.step, func(t *testing.T) {
			failure := protocolSelectionFailure(t, tc.stdout, "")
			want := fmt.Sprintf("Agent Selection failed for runtime %q: agent/protocol error at %s: adapter refused the operation (JSON-RPC -32603)", "codex", tc.step)
			if failure.Error() != want {
				t.Fatalf("want %q, got %q", want, failure.Error())
			}
			if failure.Protocol.Step != tc.step || failure.Protocol.Code != -32603 {
				t.Fatalf("detail: %+v", failure.Protocol)
			}
			if seen[failure.Error()] {
				t.Fatal("step diagnostics are identical")
			}
			seen[failure.Error()] = true
		})
	}
}

func TestACPXPromptRecordsWhetherThePromptWasSent(t *testing.T) {
	for _, tc := range []struct {
		name, stdout, step string
		sent               bool
	}{
		{"initialize", protocolInitialize + protocolErrorLine(t, 0, "init failed"), "initialize", false},
		{"model", protocolInitialize + protocolInitialized + protocolSetModel + protocolErrorLine(t, 2, "model failed"), "session/set_model", false},
		{"prompt", protocolPrompt + protocolErrorLine(t, 3, "turn failed"), "session/prompt", true},
		{"disconnect", protocolDisconnect, "session/prompt", true},
		{"startup", "", "adapter startup", false},
		{"first error", protocolInitialize + protocolErrorLine(t, 0, "first") + protocolPrompt + protocolErrorLine(t, 3, "second"), "initialize", false},
		{"unmatched id", protocolLoad + protocolSetModel + protocolErrorLine(t, 99, "unknown"), "session/set_model", false},
		{"answered setup", protocolInitialize + protocolInitialized, "session setup", false},
		{"answered prompt", protocolPrompt + `{"id":3,"result":{}}` + "\n", "session/prompt", true},
		{"string id", `{"id":"request-1","method":"authenticate"}` + "\n" + protocolErrorLine(t, "request-1", "auth failed"), "authenticate", false},
		{"null request", `{"id":null,"method":"session/prompt"}` + "\n", "adapter startup", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			failure := protocolSelectionFailure(t, tc.stdout, "")
			if tc.name == "first error" && failure.Protocol.Message != "first" {
				t.Fatalf("lost first adapter message: %+v", failure.Protocol)
			}
			if failure.Protocol.Step != tc.step || failure.Protocol.PromptSent != tc.sent {
				t.Fatalf("want %s sent=%t, got %+v", tc.step, tc.sent, failure.Protocol)
			}
			description, ok := DescribeProtocolFailure(fmt.Errorf("wrapped: %w", failure))
			if !ok || !description.FromPrompt || description.Step != tc.step || description.PromptSent != tc.sent {
				t.Fatalf("description: %+v, %t", description, ok)
			}
		})
	}
}

func TestACPXPromptBoundsTheAdapterMessage(t *testing.T) {
	message := "  adapter\n\trefused\u2003" + strings.Repeat("界", ProtocolFailureMessageLimit)
	failure := protocolSelectionFailure(t, protocolSetModel+protocolErrorLine(t, 2, message), "")
	detail := failure.Protocol.Message
	if len(detail) != len("adapter refused ")+3*((ProtocolFailureMessageLimit-len("adapter refused "))/3) || !utf8.ValidString(detail) || strings.ContainsAny(detail, "\n\r\t") || !strings.HasPrefix(detail, "adapter refused ") {
		t.Fatalf("bad message %q (%d bytes)", detail, len(detail))
	}
	for _, secret := range []string{"data-secret", "request-secret", "prompt-secret", "AGENT_DISCONNECTED"} {
		if strings.Contains(failure.Error(), secret) {
			t.Fatalf("diagnostic leaked %s", secret)
		}
	}
	description, ok := DescribeProtocolFailure(failure)
	if !ok || description.Message != detail {
		t.Fatalf("bounded description: %+v %t", description, ok)
	}
	batch := &BatchFailureError{Reason: failure.Reason, Protocol: failure.Protocol}
	if !strings.Contains(batch.Error(), " at session/set_model: "+detail+" (JSON-RPC -32603)") {
		t.Fatalf("batch text: %s", batch)
	}
}

func TestACPXPromptProtocolDetailKeepsTheExitMapping(t *testing.T) {
	failure := protocolSelectionFailure(t, protocolInitialize+protocolErrorLine(t, 0, "stdout diagnostic"), "stderr cause\n")
	if failure.Reason != acpxExitReasonAgentProtocol || failure.Err == nil || failure.Err.Error() != "stderr cause" {
		t.Fatalf("changed mapping: %+v", failure)
	}
	empty := protocolSelectionFailure(t, "", "stderr fallback\n tail")
	if empty.Protocol.Message != "" || empty.Protocol.Code != 0 {
		t.Fatalf("invented error detail: %+v", empty.Protocol)
	}
	description, ok := DescribeProtocolFailure(empty)
	if !ok || description.Message != "stderr fallback tail" {
		t.Fatalf("stderr description: %+v, %t", description, ok)
	}
	for _, code := range []int{2, 4, 130} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: protocolInitialize + protocolErrorLine(t, 0, "diagnostic"), exitCode: code})
			var selection *SelectionFailureError
			var batch *BatchFailureError
			if errors.As(run.err, &selection) || errors.As(run.err, &batch) {
				t.Fatalf("exit %d changed type: %T", code, run.err)
			}
		})
	}
	t.Run("Agent output keeps Batch failure", func(t *testing.T) {
		stdout := protocolPrompt + acpxUpdateLine(`{"sessionId":"test-session","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"started"}}}`) + protocolErrorLine(t, 3, "turn failed")
		run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout, stderr: "stderr cause", exitCode: 1})
		var batch *BatchFailureError
		if !errors.As(run.err, &batch) || batch.Protocol == nil || batch.ExitCode != 1 || batch.Reason != acpxExitReasonAgentProtocol || batch.Stderr != "stderr cause" || batch.Err != nil {
			t.Fatalf("changed Batch mapping: %T %+v", run.err, run.err)
		}
		if !batch.Protocol.PromptSent || batch.Protocol.Step != "session/prompt" || !strings.Contains(batch.Error(), "at session/prompt: turn failed (JSON-RPC -32603): stderr cause") {
			t.Fatalf("Batch protocol detail: %v", batch)
		}
	})
	// Exit-zero stream errors and parsed-result transport anomalies keep their paths.
	run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: protocolInitialize + protocolErrorLine(t, 0, "diagnostic")})
	var selection *SelectionFailureError
	if !errors.As(run.err, &selection) || selection.Protocol != nil {
		t.Fatalf("exit-zero trace attached: %v", run.err)
	}
	run = runFakeACPXPrompt(t, fakeACPXPrompt{stdout: protocolPrompt + `{"id":3,"result":{"stopReason":"end_turn"}}` + "\n", exitCode: 1})
	if run.err != nil || run.result.TransportAnomaly == "" {
		t.Fatalf("parsed-result precedence changed: %v %+v", run.err, run.result)
	}
}

func TestDescribeProtocolFailurePlacesSessionPreparation(t *testing.T) {
	cause := &InfrastructureError{Stderr: "adapter\n rejected\tselection"}
	for _, tc := range []struct {
		name, step string
		err        error
	}{
		{"ensure", "sessions ensure", &SelectionRejectedError{Operation: "ensure Agent Session", Err: cause}},
		{"model", "sessions ensure", &ModelNotAdvertisedError{Err: cause}},
		{"set model", "set model", &SelectionRejectedError{Operation: "set model", Err: cause}},
		{"set effort", "set reasoning_effort", &SelectionRejectedError{Operation: "set reasoning_effort", Err: cause}},
		{"access", "set-mode", &AccessPolicyError{Err: cause}},
		{"adapter value", "adapter startup", AdapterProbeError{Err: cause}},
		{"adapter pointer", "adapter startup", &AdapterProbeError{Err: cause}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			description, ok := DescribeProtocolFailure(fmt.Errorf("outer: %w", tc.err))
			if !ok || description.Step != tc.step || description.Message != "adapter rejected selection" || description.FromPrompt || description.PromptSent {
				t.Fatalf("description: %+v, %t", description, ok)
			}
		})
	}
	for _, err := range []error{nil, errors.New("plain"), cause, &SelectionRejectedError{Operation: "inspect capabilities"}, &SelectionFailureError{Reason: "unplaced"}} {
		if description, ok := DescribeProtocolFailure(err); ok {
			t.Fatalf("placed %T as %+v", err, description)
		}
	}
	t.Run("preparation stderr is bounded", func(t *testing.T) {
		stderr := "discard-first\n" + strings.Repeat("short\n", infrastructureStderrTailLines-1) + strings.Repeat("界", ProtocolFailureMessageLimit)
		description, ok := DescribeProtocolFailure(&AccessPolicyError{Err: &InfrastructureError{Stderr: stderr}})
		if !ok || len(description.Message) > ProtocolFailureMessageLimit || !utf8.ValidString(description.Message) || strings.ContainsAny(description.Message, "\n\r\t") || strings.Contains(description.Message, "discard-first") {
			t.Fatalf("unbounded stderr description: %+v %t", description, ok)
		}
	})
	// Prompt detail outranks preparation types in the same chain.
	err := &SelectionRejectedError{Operation: "ensure", Err: &BatchFailureError{Protocol: &ProtocolFailure{Step: "initialize", Message: "prompt detail"}}}
	description, ok := DescribeProtocolFailure(err)
	if !ok || !description.FromPrompt || description.Step != "initialize" || description.Message != "prompt detail" {
		t.Fatalf("precedence: %+v %t", description, ok)
	}
}
