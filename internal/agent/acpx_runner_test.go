package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
)

const (
	fakeACPXEnv        = "ROUNDFIX_FAKE_ACPX"
	fakeACPXArgsPath   = "ROUNDFIX_FAKE_ACPX_ARGS_PATH"
	fakeACPXInvokes    = "ROUNDFIX_FAKE_ACPX_INVOKES"
	fakeACPXPromptPath = "ROUNDFIX_FAKE_ACPX_PROMPT_PATH"
	fakeACPXStdout     = "ROUNDFIX_FAKE_ACPX_STDOUT"
	fakeACPXStderr     = "ROUNDFIX_FAKE_ACPX_STDERR"
	fakeACPXExitCode   = "ROUNDFIX_FAKE_ACPX_EXIT_CODE"
	fakeACPXStdoutBy   = "ROUNDFIX_FAKE_ACPX_STDOUT_BY_COMMAND"
	fakeACPXStderrBy   = "ROUNDFIX_FAKE_ACPX_STDERR_BY_COMMAND"
	fakeACPXExitBy     = "ROUNDFIX_FAKE_ACPX_EXIT_BY_COMMAND"
	fakeACPXCanceled   = "ROUNDFIX_FAKE_ACPX_CANCELED"
	fakeACPXStarted    = "ROUNDFIX_FAKE_ACPX_STARTED"
	fakeACPXBlock      = "ROUNDFIX_FAKE_ACPX_BLOCK_PROMPT"
	fakeACPXExitCancel = "ROUNDFIX_FAKE_ACPX_EXIT_AFTER_CANCEL"
)

func TestMain(m *testing.M) {
	if os.Getenv(fakeACPXEnv) == "1" {
		os.Exit(runFakeACPXProcess())
	}
	os.Exit(m.Run())
}

func TestACPXProbePassesWhenVersionMatchesPin(t *testing.T) {
	err, invocations := runFakeACPXProbe(t, RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, PinnedACPXVersion)

	if err != nil {
		t.Fatalf("expected matching acpx version to pass, got %v", err)
	}
	want := [][]string{{"--version"}}
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("unexpected probe invocations\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestACPXProbeMissingBinaryNamesInstallCommand(t *testing.T) {
	err := (ACPXRunner{Command: filepath.Join(t.TempDir(), "missing-acpx")}).Probe(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP})

	if err == nil {
		t.Fatal("expected missing acpx binary to fail")
	}
	message := err.Error()
	if !strings.Contains(message, acpxInstallCommand()) {
		t.Fatalf("expected install command in probe error, got %q", message)
	}
	if strings.Count(message, acpxInstallCommand()) != 1 {
		t.Fatalf("expected exactly one install command, got %q", message)
	}
}

func TestACPXProbeMismatchedVersionNamesFoundRequiredAndInstallCommand(t *testing.T) {
	const foundVersion = "0.11.0"

	err, invocations := runFakeACPXProbe(t, RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, foundVersion)

	if err == nil {
		t.Fatal("expected mismatched acpx version to fail")
	}
	message := err.Error()
	for _, want := range []string{foundVersion, PinnedACPXVersion, acpxInstallCommand(), "upgrade or downgrade"} {
		if !strings.Contains(message, want) {
			t.Fatalf("expected probe error to contain %q, got %q", want, message)
		}
	}
	if strings.Count(message, acpxInstallCommand()) != 1 {
		t.Fatalf("expected exactly one install command, got %q", message)
	}
	wantInvocations := [][]string{{"--version"}}
	if !reflect.DeepEqual(invocations, wantInvocations) {
		t.Fatalf("unexpected probe invocations\nwant: %#v\ngot:  %#v", wantInvocations, invocations)
	}
}

func TestACPXProbeCommandOverrideStillChecksACPXClient(t *testing.T) {
	runtime := RuntimeSpec{ID: "codex-custom", Protocol: ProtocolStdio, Command: "custom-acp --stdio"}

	err, invocations := runFakeACPXProbe(t, runtime, PinnedACPXVersion)

	if err != nil {
		t.Fatalf("expected command override probe to pass through acpx, got %v", err)
	}
	want := [][]string{{"--version"}}
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("expected command override to probe acpx only\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestACPXPromptArgsMatchTechSpecOrder(t *testing.T) {
	tests := []struct {
		name    string
		runtime RuntimeSpec
		want    []string
	}{
		{
			name:    "default adapter",
			runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP},
			want: []string{
				"codex", "prompt",
				"-s", "roundfix-run-1",
				"--cwd", "/repo",
				"--format", "json",
				"--json-strict",
				"--approve-all",
				"-f", "-",
			},
		},
		{
			name:    "model set",
			runtime: RuntimeSpec{ID: "claude", Protocol: ProtocolACP, Model: "opus-test"},
			want: []string{
				"claude", "prompt",
				"-s", "roundfix-run-1",
				"--cwd", "/repo",
				"--format", "json",
				"--json-strict",
				"--approve-all",
				"--model", "opus-test",
				"-f", "-",
			},
		},
		{
			name:    "command override",
			runtime: RuntimeSpec{ID: "codex-custom", Protocol: ProtocolStdio, Command: "custom-acp --stdio"},
			want: []string{
				"--agent", "custom-acp --stdio", "prompt",
				"-s", "roundfix-run-1",
				"--cwd", "/repo",
				"--format", "json",
				"--json-strict",
				"--approve-all",
				"-f", "-",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := acpxPromptArgs(ACPXPromptRequest{
				ExecuteRequest: ExecuteRequest{Runtime: tt.runtime, GitRoot: "/repo"},
				Session:        "roundfix-run-1",
			})
			if err != nil {
				t.Fatalf("acpx prompt args: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("unexpected args\nwant: %#v\ngot:  %#v", tt.want, got)
			}
		})
	}
}

func TestACPXRunPromptSendsPromptOnStdin(t *testing.T) {
	run := runFakeACPXPrompt(t, fakeACPXPrompt{
		runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-test"},
		prompt:  "built prompt",
		stdout:  acpxPromptResponseLine("end_turn"),
	})
	if run.err != nil {
		t.Fatalf("run fake acpx: %v", run.err)
	}
	wantArgs := []string{
		"codex", "prompt",
		"-s", "roundfix-run-1",
		"--cwd", run.gitRoot,
		"--format", "json",
		"--json-strict",
		"--approve-all",
		"--model", "gpt-test",
		"-f", "-",
	}
	if !reflect.DeepEqual(run.args, wantArgs) {
		t.Fatalf("unexpected acpx args\nwant: %#v\ngot:  %#v", wantArgs, run.args)
	}
	if run.prompt != "built prompt" {
		t.Fatalf("expected prompt on stdin, got %q", run.prompt)
	}
	if run.result.StopReason != "end_turn" {
		t.Fatalf("expected stop reason end_turn, got %q", run.result.StopReason)
	}
}

func TestACPXRunEnsuresSessionOncePerRunnerAndSessionName(t *testing.T) {
	harness := newFakeACPXHarness(t)
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP}

	if _, err := harness.run(context.Background(), runtime, "roundfix-run-1"); err != nil {
		t.Fatalf("first run: %v", err)
	}
	if _, err := harness.run(context.Background(), runtime, "roundfix-run-1"); err != nil {
		t.Fatalf("second run: %v", err)
	}
	if _, err := harness.run(context.Background(), runtime, "roundfix-run-2"); err != nil {
		t.Fatalf("third run: %v", err)
	}

	invocations := readJSONInvocations(t, harness.invocationsPath)
	want := [][]string{
		{"codex", "sessions", "ensure", "--name", "roundfix-run-1", "--cwd", harness.gitRoot},
		{"codex", "prompt", "-s", "roundfix-run-1", "--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "-f", "-"},
		{"codex", "prompt", "-s", "roundfix-run-1", "--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "-f", "-"},
		{"codex", "sessions", "ensure", "--name", "roundfix-run-2", "--cwd", harness.gitRoot},
		{"codex", "prompt", "-s", "roundfix-run-2", "--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "-f", "-"},
	}
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("unexpected acpx invocations\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestACPXRunAppliesFullAccessSessionSetup(t *testing.T) {
	tests := []struct {
		name    string
		runtime RuntimeSpec
		want    func(gitRoot string) [][]string
	}{
		{
			name:    "codex mode and sandbox preset",
			runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, FullAccessMode: "full-access"},
			want: func(gitRoot string) [][]string {
				return [][]string{
					{"codex", "sessions", "ensure", "--name", "roundfix-run-1", "--cwd", gitRoot},
					{"codex", "set-mode", "full-access", "-s", "roundfix-run-1"},
					{"codex", "set", "sandbox_mode", "danger-full-access", "-s", "roundfix-run-1"},
					{"codex", "prompt", "-s", "roundfix-run-1", "--cwd", gitRoot, "--format", "json", "--json-strict", "--approve-all", "-f", "-"},
				}
			},
		},
		{
			name:    "claude mode only",
			runtime: RuntimeSpec{ID: "claude", Protocol: ProtocolACP, FullAccessMode: "bypassPermissions"},
			want: func(gitRoot string) [][]string {
				return [][]string{
					{"claude", "sessions", "ensure", "--name", "roundfix-run-1", "--cwd", gitRoot},
					{"claude", "set-mode", "bypassPermissions", "-s", "roundfix-run-1"},
					{"claude", "prompt", "-s", "roundfix-run-1", "--cwd", gitRoot, "--format", "json", "--json-strict", "--approve-all", "-f", "-"},
				}
			},
		},
		{
			name:    "opencode no full access setup",
			runtime: RuntimeSpec{ID: "opencode", Protocol: ProtocolACP},
			want: func(gitRoot string) [][]string {
				return [][]string{
					{"opencode", "sessions", "ensure", "--name", "roundfix-run-1", "--cwd", gitRoot},
					{"opencode", "prompt", "-s", "roundfix-run-1", "--cwd", gitRoot, "--format", "json", "--json-strict", "--approve-all", "-f", "-"},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			harness := newFakeACPXHarness(t)

			if _, err := harness.run(context.Background(), tt.runtime, "roundfix-run-1"); err != nil {
				t.Fatalf("run acpx: %v", err)
			}

			if got, want := readJSONInvocations(t, harness.invocationsPath), tt.want(harness.gitRoot); !reflect.DeepEqual(got, want) {
				t.Fatalf("unexpected acpx invocations\nwant: %#v\ngot:  %#v", want, got)
			}
		})
	}
}

func TestACPXRunFailsSetupWhenFullAccessModeFails(t *testing.T) {
	harness := newFakeACPXHarness(t)
	t.Setenv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"set-mode": 2}))
	t.Setenv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"set-mode": "mode rejected\n"}))

	_, err := harness.run(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP, FullAccessMode: "full-access"}, "roundfix-run-1")

	if err == nil {
		t.Fatal("expected setup failure")
	}
	if !strings.Contains(err.Error(), "set acpx Agent Session mode") || !strings.Contains(err.Error(), "mode rejected") {
		t.Fatalf("expected set-mode context, got %v", err)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	want := [][]string{
		{"codex", "sessions", "ensure", "--name", "roundfix-run-1", "--cwd", harness.gitRoot},
		{"codex", "set-mode", "full-access", "-s", "roundfix-run-1"},
	}
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("unexpected invocations after setup failure\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestACPXRunWarnsWhenCodexSandboxPresetUnavailable(t *testing.T) {
	harness := newFakeACPXHarness(t)
	t.Setenv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"set": 2}))
	t.Setenv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{
		"set": "ACP session session-id does not advertise config option 'sandbox_mode'. Supported config options: mode, model.\n",
	}))
	sink := newCaptureSink("")

	_, err := harness.runWithSink(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP, FullAccessMode: "full-access"}, "roundfix-run-1", sink)

	if err != nil {
		t.Fatalf("expected prompt to continue after unavailable sandbox preset, got %v", err)
	}
	if !sink.HasStatus("codex_sandbox_full_access_unavailable") {
		t.Fatalf("expected sandbox warning status event, got %+v", sink.Events())
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	want := [][]string{
		{"codex", "sessions", "ensure", "--name", "roundfix-run-1", "--cwd", harness.gitRoot},
		{"codex", "set-mode", "full-access", "-s", "roundfix-run-1"},
		{"codex", "set", "sandbox_mode", "danger-full-access", "-s", "roundfix-run-1"},
		{"codex", "prompt", "-s", "roundfix-run-1", "--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "-f", "-"},
	}
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("unexpected invocations\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestACPXRunCancelsPromptCooperatively(t *testing.T) {
	harness := newBlockingFakeACPXHarness(t, true)
	ctx, cancel := context.WithCancel(context.Background())
	resultCh := make(chan error, 1)
	go func() {
		_, err := harness.run(ctx, RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, "roundfix-run-1")
		resultCh <- err
	}()
	waitForFile(t, harness.startedPath)
	cancel()

	err := receiveError(t, resultCh)
	if !IsStopError(err) {
		t.Fatalf("expected StopError, got %T %v", err, err)
	}
	var stopErr StopError
	if !errors.As(err, &stopErr) {
		t.Fatalf("expected StopError details, got %T %v", err, err)
	}
	if stopErr.Killed {
		t.Fatalf("expected cooperative stop without kill, got %#v", stopErr)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if !containsInvocation(invocations, []string{"codex", "cancel", "-s", "roundfix-run-1"}) {
		t.Fatalf("expected cancel invocation, got %#v", invocations)
	}
}

func TestACPXRunKillsPromptAfterCancelGracePeriod(t *testing.T) {
	harness := newBlockingFakeACPXHarness(t, false)
	ctx, cancel := context.WithCancel(context.Background())
	resultCh := make(chan error, 1)
	go func() {
		_, err := harness.run(ctx, RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, "roundfix-run-1")
		resultCh <- err
	}()
	waitForFile(t, harness.startedPath)
	cancel()

	err := receiveError(t, resultCh)
	var stopErr StopError
	if !errors.As(err, &stopErr) {
		t.Fatalf("expected StopError, got %T %v", err, err)
	}
	if !stopErr.Killed {
		t.Fatalf("expected fallback kill after cancel, got %#v", stopErr)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if !containsInvocation(invocations, []string{"codex", "cancel", "-s", "roundfix-run-1"}) {
		t.Fatalf("expected cancel invocation before kill, got %#v", invocations)
	}
}

func TestACPXEndSessionClosesBestEffort(t *testing.T) {
	harness := newFakeACPXHarness(t)
	t.Setenv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions close": 1}))
	t.Setenv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"sessions close": "close rejected\n"}))
	var warnings []string
	harness.runner.warnf = func(format string, args ...any) {
		warnings = append(warnings, fmt.Sprintf(format, args...))
	}

	err := harness.runner.EndSession(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, SessionRef{Name: "roundfix-run-1"})

	if err != nil {
		t.Fatalf("expected close failure to be swallowed, got %v", err)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	want := [][]string{{"codex", "sessions", "close", "-s", "roundfix-run-1"}}
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("unexpected close invocation\nwant: %#v\ngot:  %#v", want, invocations)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "roundfix-run-1") || !strings.Contains(warnings[0], "close rejected") {
		t.Fatalf("expected close warning with session and stderr, got %#v", warnings)
	}
}

func TestACPXRunPromptPublishesUpdateLinesAndCapturesStopReason(t *testing.T) {
	messageLine := acpxUpdateLine(`{"sessionId":"sess-1","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"hello"}}}`)
	thoughtLine := acpxUpdateLine(`{"sessionId":"sess-1","update":{"sessionUpdate":"agent_thought_chunk","content":{"type":"text","text":"thinking"}}}`)
	responseLine := acpxPromptResponseLine("end_turn")
	stdout := messageLine + thoughtLine + responseLine

	run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout})
	if run.err != nil {
		t.Fatalf("run fake acpx: %v", run.err)
	}
	if run.result.StopReason != "end_turn" {
		t.Fatalf("expected stop reason end_turn, got %q", run.result.StopReason)
	}
	events := run.sink.Events()
	if len(events) != 2 {
		t.Fatalf("expected only update events to be journaled, got %d: %+v", len(events), events)
	}
	expectedKinds := []runevent.Kind{runevent.KindAgentMessage, runevent.KindAgentThought}
	expectedPayloads := []string{messageLine, thoughtLine}
	for index, event := range events {
		if event.Kind != expectedKinds[index] {
			t.Fatalf("expected event %d kind %q, got %q", index, expectedKinds[index], event.Kind)
		}
		if event.RunID != "run-acpx" || event.Batch != 7 || event.Source != runevent.SourceAgent {
			t.Fatalf("expected Run identity on event %d, got %+v", index, event)
		}
		if string(event.Payload) != expectedPayloads[index] {
			t.Fatalf("expected byte-identical payload for event %d\nwant: %q\ngot:  %q", index, expectedPayloads[index], string(event.Payload))
		}
	}
	if logContent := readFile(t, run.logPath); logContent != stdout {
		t.Fatalf("expected agent log to contain every stdout line in order\nwant: %q\ngot:  %q", stdout, logContent)
	}
}

func TestStreamUpdateFromEventAcceptsACPXJSONRPCLinePayload(t *testing.T) {
	line := acpxUpdateLine(`{"sessionId":"sess-1","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"hello from acpx"}}}`)
	update, ok := StreamUpdateFromEvent(runevent.RunEvent{
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentMessage,
		Payload: json.RawMessage(line),
	})
	if !ok {
		t.Fatal("expected acpx JSON-RPC line payload to decode")
	}
	if update.Kind != StreamUpdateMessage || update.Text != "hello from acpx" {
		t.Fatalf("unexpected stream update: %#v", update)
	}
}

func TestACPXExitCodeMapping(t *testing.T) {
	tests := []struct {
		name       string
		exitCode   int
		stdout     string
		stderr     string
		assertErr  func(t *testing.T, err error)
		assertPost func(t *testing.T, run fakeACPXRun)
	}{
		{
			name:     "success",
			exitCode: 0,
			stdout:   acpxPromptResponseLine("end_turn"),
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("expected success, got %v", err)
				}
			},
			assertPost: func(t *testing.T, run fakeACPXRun) {
				t.Helper()
				if run.result.StopReason != "end_turn" {
					t.Fatalf("expected stop reason, got %q", run.result.StopReason)
				}
			},
		},
		{
			name:     "agent protocol failure",
			exitCode: 1,
			stderr:   "protocol exploded\n",
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var batchErr *BatchFailureError
				if !errors.As(err, &batchErr) {
					t.Fatalf("expected BatchFailureError, got %T %v", err, err)
				}
				if batchErr.Reason != acpxExitReasonAgentProtocol {
					t.Fatalf("expected protocol reason, got %q", batchErr.Reason)
				}
				if !strings.Contains(err.Error(), "protocol exploded") {
					t.Fatalf("expected stderr in error context, got %q", err.Error())
				}
			},
		},
		{
			name:     "timeout failure",
			exitCode: 3,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var batchErr *BatchFailureError
				if !errors.As(err, &batchErr) {
					t.Fatalf("expected BatchFailureError, got %T %v", err, err)
				}
				if batchErr.Reason != acpxExitReasonTimeout {
					t.Fatalf("expected timeout reason, got %q", batchErr.Reason)
				}
			},
		},
		{
			name:     "all permissions denied failure",
			exitCode: 5,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var batchErr *BatchFailureError
				if !errors.As(err, &batchErr) {
					t.Fatalf("expected BatchFailureError, got %T %v", err, err)
				}
				if batchErr.Reason != acpxExitReasonPermissionsDenied {
					t.Fatalf("expected permissions denied reason, got %q", batchErr.Reason)
				}
			},
			assertPost: func(t *testing.T, run fakeACPXRun) {
				t.Helper()
				if !run.sink.HasStatus(acpxPermissionDeniedStatus) {
					t.Fatalf("expected loud permission-denied status event, got %+v", run.sink.Events())
				}
			},
		},
		{
			name:     "usage infrastructure error",
			exitCode: 2,
			stderr:   "bad acpx usage\n",
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var infraErr *InfrastructureError
				if !errors.As(err, &infraErr) {
					t.Fatalf("expected InfrastructureError, got %T %v", err, err)
				}
				if infraErr.Reason != acpxExitReasonUsage {
					t.Fatalf("expected usage reason, got %q", infraErr.Reason)
				}
				if !strings.Contains(err.Error(), "bad acpx usage") {
					t.Fatalf("expected stderr in error context, got %q", err.Error())
				}
			},
			assertPost: func(t *testing.T, run fakeACPXRun) {
				t.Helper()
				if strings.Contains(readFile(t, run.logPath), "bad acpx usage") {
					t.Fatal("stderr must not be written to the Agent log")
				}
				if strings.Contains(run.sink.Text(), "bad acpx usage") {
					t.Fatal("stderr must not be journaled as a Run Event")
				}
			},
		},
		{
			name:     "missing session infrastructure error",
			exitCode: 4,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var infraErr *InfrastructureError
				if !errors.As(err, &infraErr) {
					t.Fatalf("expected InfrastructureError, got %T %v", err, err)
				}
				if infraErr.Reason != acpxExitReasonMissingSession {
					t.Fatalf("expected missing session reason, got %q", infraErr.Reason)
				}
			},
		},
		{
			name:     "stop request",
			exitCode: 130,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				if !IsStopError(err) {
					t.Fatalf("expected StopError, got %T %v", err, err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run := runFakeACPXPrompt(t, fakeACPXPrompt{
				stdout:   tt.stdout,
				stderr:   tt.stderr,
				exitCode: tt.exitCode,
			})
			tt.assertErr(t, run.err)
			if tt.assertPost != nil {
				tt.assertPost(t, run)
			}
		})
	}
}

type fakeACPXPrompt struct {
	runtime  RuntimeSpec
	prompt   string
	stdout   string
	stderr   string
	exitCode int
}

type fakeACPXRun struct {
	result  ExecuteResult
	err     error
	sink    *captureSink
	gitRoot string
	logPath string
	args    []string
	prompt  string
}

type fakeACPXHarness struct {
	runner          *ACPXRunner
	gitRoot         string
	invocationsPath string
	startedPath     string
}

func newFakeACPXHarness(t *testing.T) *fakeACPXHarness {
	t.Helper()
	dir := t.TempDir()
	invocationsPath := filepath.Join(dir, "invocations.jsonl")
	t.Setenv(fakeACPXEnv, "1")
	t.Setenv(fakeACPXInvokes, invocationsPath)
	t.Setenv(fakeACPXStdoutBy, mustJSONForTest(t, map[string]string{"prompt": acpxPromptResponseLine("end_turn")}))
	return &fakeACPXHarness{
		runner: &ACPXRunner{
			Command: os.Args[0],
			Now: func() time.Time {
				return time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
			},
		},
		gitRoot:         dir,
		invocationsPath: invocationsPath,
	}
}

func newBlockingFakeACPXHarness(t *testing.T, exitAfterCancel bool) *fakeACPXHarness {
	t.Helper()
	harness := newFakeACPXHarness(t)
	harness.startedPath = filepath.Join(harness.gitRoot, "prompt-started")
	t.Setenv(fakeACPXBlock, "1")
	t.Setenv(fakeACPXStarted, harness.startedPath)
	t.Setenv(fakeACPXCanceled, filepath.Join(harness.gitRoot, "canceled"))
	if exitAfterCancel {
		t.Setenv(fakeACPXExitCancel, "1")
	}
	return harness
}

func (harness *fakeACPXHarness) run(ctx context.Context, runtime RuntimeSpec, sessionName string) (ExecuteResult, error) {
	return harness.runWithSink(ctx, runtime, sessionName, newCaptureSink(""))
}

func (harness *fakeACPXHarness) runWithSink(ctx context.Context, runtime RuntimeSpec, sessionName string, sink runevent.Sink) (ExecuteResult, error) {
	return harness.runner.Run(ctx, ExecuteRequest{
		Runtime:   runtime,
		RunID:     "run-acpx",
		Batch:     rounds.Batch{Number: 7},
		LogPath:   filepath.Join(harness.gitRoot, "runs", "run-acpx", "agent", "batch-007.log"),
		Prompt:    "prompt",
		GitRoot:   harness.gitRoot,
		StopGrace: 20 * time.Millisecond,
		Session:   SessionRef{Name: sessionName},
	}, sink)
}

func runFakeACPXPrompt(t *testing.T, prompt fakeACPXPrompt) fakeACPXRun {
	t.Helper()
	dir := t.TempDir()
	argsPath := filepath.Join(dir, "args.json")
	promptPath := filepath.Join(dir, "prompt.txt")
	t.Setenv(fakeACPXEnv, "1")
	t.Setenv(fakeACPXArgsPath, argsPath)
	t.Setenv(fakeACPXPromptPath, promptPath)
	t.Setenv(fakeACPXStdout, prompt.stdout)
	t.Setenv(fakeACPXStderr, prompt.stderr)
	t.Setenv(fakeACPXExitCode, strconv.Itoa(prompt.exitCode))

	runtime := prompt.runtime
	if runtime.ID == "" && runtime.Protocol == "" && runtime.Command == "" {
		runtime = RuntimeSpec{ID: "codex", Protocol: ProtocolACP}
	}
	if prompt.prompt == "" {
		prompt.prompt = "prompt"
	}
	sink := newCaptureSink("")
	logPath := filepath.Join(dir, "runs", "run-acpx", "agent", "batch-007.log")
	result, err := (ACPXRunner{
		Command: os.Args[0],
		Now: func() time.Time {
			return time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
		},
	}).RunPrompt(context.Background(), ACPXPromptRequest{
		ExecuteRequest: ExecuteRequest{
			Runtime: runtime,
			RunID:   "run-acpx",
			Batch:   rounds.Batch{Number: 7},
			LogPath: logPath,
			Prompt:  prompt.prompt,
			GitRoot: dir,
		},
		Session: "roundfix-run-1",
	}, sink)
	return fakeACPXRun{
		result:  result,
		err:     err,
		sink:    sink,
		gitRoot: dir,
		logPath: logPath,
		args:    readJSONArgs(t, argsPath),
		prompt:  readFile(t, promptPath),
	}
}

func runFakeACPXProbe(t *testing.T, runtime RuntimeSpec, version string) (error, [][]string) {
	t.Helper()
	dir := t.TempDir()
	invocationsPath := filepath.Join(dir, "invocations.jsonl")
	t.Setenv(fakeACPXEnv, "1")
	t.Setenv(fakeACPXInvokes, invocationsPath)
	t.Setenv(fakeACPXStdout, version+"\n")

	err := (ACPXRunner{Command: os.Args[0]}).Probe(context.Background(), runtime)
	return err, readJSONInvocations(t, invocationsPath)
}

func readJSONInvocations(t *testing.T, path string) [][]string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read invocations: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	invocations := make([][]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var args []string
		if err := json.Unmarshal([]byte(line), &args); err != nil {
			t.Fatalf("decode invocation %q: %v", line, err)
		}
		invocations = append(invocations, args)
	}
	return invocations
}

func containsInvocation(invocations [][]string, want []string) bool {
	for _, invocation := range invocations {
		if reflect.DeepEqual(invocation, want) {
			return true
		}
	}
	return false
}

func waitForFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.After(5 * time.Second)
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for %s", path)
		case <-ticker.C:
		}
	}
}

func receiveError(t *testing.T, resultCh <-chan error) error {
	t.Helper()
	select {
	case err := <-resultCh:
		return err
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for acpx run")
		return nil
	}
}

func mustJSONForTest[T any](t *testing.T, value T) string {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal test JSON: %v", err)
	}
	return string(payload)
}

func readJSONArgs(t *testing.T, path string) []string {
	t.Helper()
	var args []string
	if err := json.Unmarshal([]byte(readFile(t, path)), &args); err != nil {
		t.Fatalf("decode captured args: %v", err)
	}
	return args
}

func acpxUpdateLine(params string) string {
	return `{"jsonrpc":"2.0","method":"session/update","params":` + params + `}` + "\n"
}

func acpxPromptResponseLine(stopReason string) string {
	return `{"jsonrpc":"2.0","id":1,"result":{"stopReason":"` + stopReason + `"}}` + "\n"
}

func runFakeACPXProcess() int {
	args := os.Args[1:]
	commandKey := fakeACPXCommandKey(args)
	if path := os.Getenv(fakeACPXArgsPath); path != "" {
		payload, err := json.Marshal(args)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "marshal args: %v\n", err)
			return 2
		}
		if err := os.WriteFile(path, payload, 0o644); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "write args: %v\n", err)
			return 2
		}
	}
	if path := os.Getenv(fakeACPXInvokes); path != "" {
		if err := appendFakeACPXInvocation(path, args); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "append invocation: %v\n", err)
			return 2
		}
	}
	if path := os.Getenv(fakeACPXPromptPath); path != "" {
		prompt, err := io.ReadAll(os.Stdin)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "read stdin: %v\n", err)
			return 2
		}
		if err := os.WriteFile(path, prompt, 0o644); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "write prompt: %v\n", err)
			return 2
		}
	}
	if commandKey == "cancel" {
		if path := os.Getenv(fakeACPXCanceled); path != "" {
			if err := os.WriteFile(path, []byte("canceled\n"), 0o644); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "write cancel marker: %v\n", err)
				return 2
			}
		}
	}
	if commandKey == "prompt" && os.Getenv(fakeACPXBlock) == "1" {
		if path := os.Getenv(fakeACPXStarted); path != "" {
			if err := os.WriteFile(path, []byte("started\n"), 0o644); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "write started marker: %v\n", err)
				return 2
			}
		}
		for {
			if canceled := os.Getenv(fakeACPXCanceled); canceled != "" {
				if _, err := os.Stat(canceled); err == nil && os.Getenv(fakeACPXExitCancel) == "1" {
					return 130
				}
			}
			time.Sleep(5 * time.Millisecond)
		}
	}
	stdoutByCommand := fakeACPXStringMap(os.Getenv(fakeACPXStdoutBy))
	stderrByCommand := fakeACPXStringMap(os.Getenv(fakeACPXStderrBy))
	_, _ = io.WriteString(os.Stdout, firstFakeACPXString(stdoutByCommand[commandKey], os.Getenv(fakeACPXStdout)))
	_, _ = io.WriteString(os.Stderr, firstFakeACPXString(stderrByCommand[commandKey], os.Getenv(fakeACPXStderr)))
	exitByCommand := fakeACPXIntMap(os.Getenv(fakeACPXExitBy))
	if exitCode, ok := exitByCommand[commandKey]; ok {
		return exitCode
	}
	if rawExitCode := os.Getenv(fakeACPXExitCode); rawExitCode != "" {
		exitCode, err := strconv.Atoi(rawExitCode)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "parse exit code: %v\n", err)
			return 2
		}
		return exitCode
	}
	return 0
}

func appendFakeACPXInvocation(path string, args []string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	payload, err := json.Marshal(args)
	if err != nil {
		return err
	}
	if _, err := file.Write(append(payload, '\n')); err != nil {
		return err
	}
	return nil
}

func fakeACPXCommandKey(args []string) string {
	index := 1
	if len(args) >= 2 && args[0] == "--agent" {
		index = 2
	}
	if len(args) <= index {
		return ""
	}
	switch args[index] {
	case "sessions":
		if len(args) > index+1 {
			return "sessions " + args[index+1]
		}
	case "prompt", "set-mode", "set", "cancel":
		return args[index]
	}
	if len(args) > index+1 {
		return args[index+1]
	}
	return args[index]
}

func fakeACPXStringMap(raw string) map[string]string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var values map[string]string
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return nil
	}
	return values
}

func fakeACPXIntMap(raw string) map[string]int {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var values map[string]int
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return nil
	}
	return values
}

func firstFakeACPXString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
