package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

const providerCapacityText = "Selected model is at capacity. Please try a different model."
const providerPromptStart = `{"jsonrpc":"2.0","id":26,"method":"session/prompt","params":{"sessionId":"fake-thread-0268","prompt":[{"type":"text","text":"work"}]}}
{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"fake-thread-0268","update":{"sessionUpdate":"session_info_update","_meta":{"codex":{"threadStatus":{"type":"active"}}}}}}
`
const providerSystemStatus = `{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"fake-thread-0268","update":{"sessionUpdate":"session_info_update","_meta":{"codex":{"threadStatus":{"type":"systemError"}}}}}}
`
const providerPromptEnd = `{"jsonrpc":"2.0","id":26,"result":{"stopReason":"end_turn","usage":null}}
`
const providerCapacityChunk = `{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"fake-thread-0268","update":{"sessionUpdate":"agent_message_chunk","messageId":"fake-capacity","content":{"type":"text","text":"Selected model is at capacity. Please try a different model."}}}}
`

func TestProviderSystemErrorIsNotAFinishedTurn(t *testing.T) {
	t.Parallel()
	prior := acpxUpdateLine(`{"sessionId":"fake-thread-0268","update":{"sessionUpdate":"agent_message_chunk","messageId":"fake-work","content":{"type":"text","text":"working"}}}`)
	stdout := providerPromptStart + prior + providerSystemStatus + providerCapacityChunk + providerPromptEnd
	run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout})
	var batch *BatchFailureError
	if !errors.As(run.err, &batch) {
		t.Fatalf("want Batch failure, got %T: %v", run.err, run.err)
	}
	want := "Agent Batch failed: provider system error at session/prompt: " + providerCapacityText
	if batch.Reason != ProviderSystemErrorReason || batch.ExitCode != 0 || batch.Err != nil || batch.Stderr != "" {
		t.Fatalf("unexpected Batch fields: %+v", batch)
	}
	if run.err.Error() != want {
		t.Fatalf("want %q, got %q", want, run.err.Error())
	}
	if run.result.Output != stdout || run.result.StopReason != "end_turn" || !reflect.DeepEqual(run.result.Messages, []string{"working", providerCapacityText}) || !reflect.DeepEqual(run.result.Usage, TurnUsage{}) {
		t.Fatalf("lost result: %+v", run.result)
	}
}

func TestProviderSystemErrorWithoutOutputIsASelectionFailure(t *testing.T) {
	t.Parallel()
	run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: providerPromptStart + providerSystemStatus + providerPromptEnd})
	var selection *SelectionFailureError
	if !errors.As(run.err, &selection) {
		t.Fatalf("want Selection failure, got %T: %v", run.err, run.err)
	}
	want := `Agent Selection failed for runtime "codex": provider system error at session/prompt: threadStatus systemError`
	if run.err.Error() != want || selection.Reason != ProviderSystemErrorReason || selection.Err != nil {
		t.Fatalf("selection = %+v, want %q", selection, want)
	}
	description, ok := DescribeProviderSystemError(run.err)
	if !ok || description.Runtime != "codex" || description.Detail != "threadStatus systemError" {
		t.Fatalf("description = %+v, %t", description, ok)
	}
	if run.result.StopReason != "end_turn" || len(run.result.Messages) != 0 {
		t.Fatalf("lost result: %+v", run.result)
	}

	t.Run("provider text counts as Agent output", func(t *testing.T) {
		run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: providerPromptStart + providerSystemStatus + providerCapacityChunk + providerPromptEnd})
		var batch *BatchFailureError
		if !errors.As(run.err, &batch) || errors.As(run.err, &selection) {
			t.Fatalf("provider text must keep Batch classification: %v", run.err)
		}
		if run.err.Error() != "Agent Batch failed: provider system error at session/prompt: "+providerCapacityText {
			t.Fatalf("unexpected error: %v", run.err)
		}
	})
}

func TestRecoveredThreadStatusIsAFinishedTurn(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		statuses []string
	}{
		{"active then idle", []string{"active", "idle"}},
		{"active", []string{"active"}},
		{"idle", []string{"idle"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stdout := providerPromptStart + providerSystemStatus + providerCapacityChunk
			for _, status := range tc.statuses {
				stdout += providerStatusLine(status)
			}
			run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout + providerPromptEnd})
			if run.err != nil || run.result.StopReason != "end_turn" || run.result.Answer() != providerCapacityText {
				t.Fatalf("recovered turn: %+v, %v", run.result, run.err)
			}
			sealed, err := parseSealedPromptStream([]byte(stdout + providerPromptEnd))
			if err != nil || string(sealed.Output) != providerCapacityText {
				t.Fatalf("recovered sealed turn: %+v, %v", sealed, err)
			}
		})
	}
}

func TestProviderSystemErrorDetailIsBounded(t *testing.T) {
	t.Parallel()
	longText := "  provider\n\t" + strings.Repeat("界", ProtocolFailureMessageLimit)
	for _, tc := range []struct {
		name, before, after, want string
	}{
		{name: "no chunk", want: "threadStatus systemError"},
		{name: "ignore prior text", before: providerCapacityChunk, want: "threadStatus systemError"},
		{name: "split chunks", after: providerMessageLine(t, "Selected model is at ") + providerMessageLine(t, "capacity."), want: "Selected model is at capacity."},
		{name: "bounded unicode", after: providerMessageLine(t, longText), want: "provider " + strings.Repeat("界", (ProtocolFailureMessageLimit-len("provider "))/len("界"))},
		{name: "metadata without status", after: acpxUpdateLine(`{"sessionId":"fake-thread-0268","update":{"sessionUpdate":"session_info_update","title":"unrelated"}}`), want: "threadStatus systemError"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: providerPromptStart + tc.before + providerSystemStatus + tc.after + providerPromptEnd})
			description, ok := DescribeProviderSystemError(run.err)
			if !ok || description.Detail != tc.want || len(description.Detail) > ProtocolFailureMessageLimit || !utf8.ValidString(description.Detail) || strings.ContainsAny(description.Detail, "\n\r\t") {
				t.Fatalf("detail = %+v, %t; error = %v", description, ok, run.err)
			}
			protocol := protocolFailureInChain(run.err)
			if protocol == nil || !protocol.ProviderSystemError || protocol.Step != "session/prompt" || !protocol.PromptSent || protocol.Message != "" || protocol.Code != 0 || protocol.LostRollout {
				t.Fatalf("protocol = %+v", protocol)
			}
		})
	}
}

func TestSealedPromptRefusesAProviderSystemError(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, chunks, detail string
	}{
		{"capacity", providerCapacityChunk, providerCapacityText},
		{"no text", "", "threadStatus systemError"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseSealedPromptStream([]byte(providerPromptStart + providerSystemStatus + tc.chunks + providerPromptEnd))
			want := "sealed acpx prompt ended with a provider system error: " + tc.detail
			if err == nil || err.Error() != want || len(result.Output) != 0 {
				t.Fatalf("sealed result = %+v, error = %v, want %q and no output", result, err, want)
			}
		})
	}
}

func TestDescribeProviderSystemErrorTraversesFailureChains(t *testing.T) {
	t.Parallel()
	batch := &BatchFailureError{Reason: ProviderSystemErrorReason, Protocol: &ProtocolFailure{
		Step: "session/prompt", PromptSent: true, ProviderSystemError: true, Detail: providerCapacityText,
	}}
	selection := &SelectionFailureError{Runtime: "codex", Err: fmt.Errorf("inner: %w", batch)}
	for _, tc := range []struct {
		name    string
		err     error
		runtime string
	}{
		{"batch", batch, ""},
		{"wrapped batch", fmt.Errorf("outer: %w", batch), ""},
		{"selection", selection, "codex"},
		{"wrapped selection", fmt.Errorf("outer: %w", selection), "codex"},
		{"joined", errors.Join(&BatchFailureError{Protocol: &ProtocolFailure{Message: "other"}}, selection), "codex"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			description, ok := DescribeProviderSystemError(tc.err)
			if !ok || description.Runtime != tc.runtime || description.Detail != providerCapacityText {
				t.Fatalf("description = %+v, %t", description, ok)
			}
		})
	}
	for _, tc := range []struct {
		name string
		err  error
	}{
		{"nil", nil},
		{"plain text", errors.New(ProviderSystemErrorReason)},
		{"nil selection", (*SelectionFailureError)(nil)},
		{"nil batch", (*BatchFailureError)(nil)},
		{"stop pointer", &StopError{Err: batch}},
		{"stop value", StopError{Err: batch}},
		{"joined stop", errors.Join(&StopError{}, selection)},
		{"canceled", errors.Join(context.Canceled, selection)},
		{"deadline", errors.Join(context.DeadlineExceeded, selection)},
		{"lost rollout", &BatchFailureError{Protocol: &ProtocolFailure{LostRollout: true, Detail: LostRolloutPhrase}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if description, ok := DescribeProviderSystemError(tc.err); ok {
				t.Fatalf("unexpected description = %+v", description)
			}
		})
	}
}

func TestProviderSystemErrorKeepsExistingExitPaths(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		exit    int
		access  SessionAccess
		anomaly bool
		refused bool
	}{
		{name: "parsed result anomaly", exit: 1, anomaly: true},
		{name: "read only refusal", exit: 5, access: SessionAccessReadOnly, refused: true},
		{name: "stopped", exit: 130},
	} {
		t.Run(tc.name, func(t *testing.T) {
			run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: providerPromptStart + providerSystemStatus + providerCapacityChunk + providerPromptEnd, exitCode: tc.exit, access: tc.access})
			if description, ok := DescribeProviderSystemError(run.err); ok {
				t.Fatalf("exit mapping became provider failure: %+v", description)
			}
			if (run.result.TransportAnomaly != "") != tc.anomaly || run.result.PermissionRefused != tc.refused {
				t.Fatalf("changed exit mapping: %+v, %v", run.result, run.err)
			}
			if tc.exit == 130 {
				if !IsStopError(run.err) {
					t.Fatalf("want StopError, got %v", run.err)
				}
			} else if run.err != nil {
				t.Fatalf("parsed result must retain precedence: %v", run.err)
			}
		})
	}
}

func TestProviderSystemErrorUsesTheStatusAtTheResult(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, before, after string
		failure             bool
	}{
		{"late recovery", providerSystemStatus + providerCapacityChunk, providerStatusLine("idle"), true},
		{"late error", "", providerSystemStatus + providerCapacityChunk, false},
		{"capacity text alone", providerCapacityChunk, "", false},
		{"new failure resets text", providerSystemStatus + providerCapacityChunk + providerStatusLine("active") + providerSystemStatus, "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stdout := providerPromptStart + tc.before + providerPromptEnd + tc.after
			run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout})
			_, marked := DescribeProviderSystemError(run.err)
			if marked != tc.failure || (run.err != nil) != tc.failure {
				t.Fatalf("error = %v, want provider failure %t", run.err, tc.failure)
			}
			if tc.name == "new failure resets text" {
				description, _ := DescribeProviderSystemError(run.err)
				if description.Detail != "threadStatus systemError" {
					t.Fatalf("stale detail = %q", description.Detail)
				}
			}
			sealed, err := parseSealedPromptStream([]byte(stdout))
			if (err != nil) != tc.failure || (tc.failure && len(sealed.Output) != 0) {
				t.Fatalf("sealed result = %+v, error = %v", sealed, err)
			}
		})
	}
}

func providerStatusLine(status string) string {
	return strings.Replace(providerSystemStatus, ProviderSystemErrorStatus, status, 1)
}

func providerMessageLine(t *testing.T, text string) string {
	t.Helper()
	encoded, err := json.Marshal(text)
	if err != nil {
		t.Fatal(err)
	}
	return acpxUpdateLine(`{"sessionId":"fake-thread-0268","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":` + string(encoded) + `}}}`)
}
