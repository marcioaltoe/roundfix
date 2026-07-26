// Suite: Non-Run sealed ACPX execution
// Invariant: one bounded denied-tools prompt uses a fresh exactly selected session and returns only after proven cleanup.
// Boundary IN: ACPX arguments, exact selection, message extraction, tool rejection, cancellation, and session close.
// Boundary OUT: Baseline proposal schema validation, preferred/fallback policy, Run Events, and Agent logs.

package agent

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
)

func TestSealedACPXPromptReturnsMessageAndClosesSession(t *testing.T) {
	harness := newFakeACPXHarness(t)
	workDir := newPrivateEmptyDirectory(t)
	promptPath := filepath.Join(harness.gitRoot, "sealed-prompt")
	t.Setenv(fakeACPXPromptPath, promptPath)
	t.Setenv(fakeACPXStdoutCall, sealedSelectionFixtures(t, "gpt-5.5", "xhigh"))
	t.Setenv(fakeACPXStdoutBy, mustJSONForTest(t, map[string]string{
		"prompt": acpxUpdateLine(
			`{"sessionId":"sealed","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"{\"ok\":true}"}}}`,
		) + acpxPromptResponseLine("end_turn"),
	}))

	result, err := harness.runner.RunSealedPrompt(context.Background(), SealedPromptRequest{
		Runtime: RuntimeSpec{
			ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh",
		},
		WorkDir: workDir,
		Input:   []byte(`{"sealed":true}`),
	})
	if err != nil {
		t.Fatalf("run sealed prompt: %v", err)
	}
	if string(result.Output) != `{"ok":true}` || result.ToolUsed {
		t.Fatalf("unexpected sealed result: %+v", result)
	}
	if got := string(readFile(t, promptPath)); got != `{"sealed":true}` {
		t.Fatalf("sealed stdin changed: %q", got)
	}
	assertLastInvocationClosesDisposable(t, harness)
}

func TestSealedPromptStreamIgnoresNonTerminalJSONRPCResults(t *testing.T) {
	stream := `{"jsonrpc":"2.0","id":0,"result":{"protocolVersion":1}}` + "\n" +
		`{"jsonrpc":"2.0","id":2,"result":{"sessionId":"sealed"}}` + "\n" +
		acpxUpdateLine(
			`{"sessionId":"sealed","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"{\"ok\":true}"}}}`,
		) +
		acpxPromptResponseLine("end_turn")

	result, err := parseSealedPromptStream([]byte(stream))
	if err != nil {
		t.Fatalf("parse sealed stream with setup results: %v", err)
	}
	if string(result.Output) != `{"ok":true}` || result.ToolUsed {
		t.Fatalf("unexpected sealed result: %+v", result)
	}
}

func TestSealedACPXPromptDiscardsLargeThoughtStreamIncrementally(t *testing.T) {
	harness := newFakeACPXHarness(t)
	workDir := newPrivateEmptyDirectory(t)
	t.Setenv(fakeACPXStdoutCall, sealedSelectionFixtures(t, "gpt-5.5", "xhigh"))
	t.Setenv(fakeACPXThoughtLen, strconv.Itoa(sealedACPStreamMaxBytes+1))
	t.Setenv(fakeACPXStdoutBy, mustJSONForTest(t, map[string]string{
		"prompt": acpxUpdateLine(
			`{"sessionId":"sealed","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"{\"ok\":true}"}}}`,
		) + acpxPromptResponseLine("end_turn"),
	}))

	result, err := harness.runner.RunSealedPrompt(context.Background(), SealedPromptRequest{
		Runtime: RuntimeSpec{
			ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh",
		},
		WorkDir: workDir,
		Input:   []byte(`{"sealed":true}`),
	})
	if err != nil {
		t.Fatalf("run sealed prompt with large thought stream: %v", err)
	}
	if string(result.Output) != `{"ok":true}` || result.ToolUsed {
		t.Fatalf("unexpected sealed result: %+v", result)
	}
	assertLastInvocationClosesDisposable(t, harness)
}

func TestSealedACPXPromptRejectsToolUseAndClosesSession(t *testing.T) {
	harness := newFakeACPXHarness(t)
	workDir := newPrivateEmptyDirectory(t)
	t.Setenv(fakeACPXStdoutCall, sealedSelectionFixtures(t, "gpt-5.5", "xhigh"))
	t.Setenv(fakeACPXStdoutBy, mustJSONForTest(t, map[string]string{
		"prompt": acpxUpdateLine(
			`{"sessionId":"sealed","update":{"sessionUpdate":"tool_call","toolCallId":"tool-1","kind":"read","title":"read","status":"pending"}}`,
		) + acpxPromptResponseLine("end_turn"),
	}))

	_, err := harness.runner.RunSealedPrompt(context.Background(), SealedPromptRequest{
		Runtime: RuntimeSpec{
			ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh",
		},
		WorkDir: workDir,
		Input:   []byte(`{"sealed":true}`),
	})
	if !errors.Is(err, ErrSealedToolUse) {
		t.Fatalf("tool-use error = %T %v, want ErrSealedToolUse", err, err)
	}
	assertLastInvocationClosesDisposable(t, harness)
}

func TestSealedACPXPromptCancellationCancelsAndClosesSession(t *testing.T) {
	harness := newFakeACPXHarness(t)
	workDir := newPrivateEmptyDirectory(t)
	started := filepath.Join(harness.gitRoot, "sealed-started")
	canceled := filepath.Join(harness.gitRoot, "sealed-canceled")
	closed := filepath.Join(harness.gitRoot, "sealed-closed")
	t.Setenv(fakeACPXStarted, started)
	t.Setenv(fakeACPXCanceled, canceled)
	t.Setenv(fakeACPXClosed, closed)
	t.Setenv(fakeACPXBlock, "1")
	t.Setenv(fakeACPXStdoutCall, sealedSelectionFixtures(t, "gpt-5.5", "xhigh"))

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := harness.runner.RunSealedPrompt(ctx, SealedPromptRequest{
			Runtime: RuntimeSpec{
				ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh",
			},
			WorkDir: workDir,
			Input:   []byte(`{"sealed":true}`),
		})
		result <- err
	}()
	waitForFile(t, started)
	cancel()

	if err := receiveError(t, result); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation error = %T %v, want context.Canceled", err, err)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if !containsCommandKey(invocations, "cancel") {
		t.Fatalf("cancellation did not invoke Agent Session cancel: %#v", invocations)
	}
	if len(invocations) == 0 || fakeACPXCommandKey(invocations[len(invocations)-1]) != "sessions close" {
		t.Fatalf("cancellation did not finish with Agent Session close: %#v", invocations)
	}
}

func sealedSelectionFixtures(t *testing.T, model, effort string) string {
	t.Helper()
	models := []string{"gpt-5.6-sol", "gpt-5.5"}
	efforts := []string{"low", "medium", "high", "xhigh", "max", "ultra"}
	return mustJSONForTest(t, map[string]string{
		"sessions show": sessionCapabilitySnapshotFixture(
			t, model, models, "reasoning_effort", "medium", efforts,
		),
		"set model value=" + model: selectionStateFixture(
			t, "model", model, model, models, "reasoning_effort", "medium", efforts,
		),
		"set reasoning_effort value=" + effort: selectionStateFixture(
			t, "reasoning_effort", effort, model, models, "reasoning_effort", effort, efforts,
		),
	})
}

func newPrivateEmptyDirectory(t *testing.T) string {
	t.Helper()
	workDir := filepath.Join(t.TempDir(), "empty")
	if err := os.Mkdir(workDir, 0o700); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(workDir)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("test directory is not private: %v", info.Mode().Perm())
	}
	if entries, err := os.ReadDir(workDir); err != nil || !reflect.DeepEqual(entries, []os.DirEntry{}) {
		t.Fatalf("test directory is not empty: entries=%v err=%v", entries, err)
	}
	return workDir
}
