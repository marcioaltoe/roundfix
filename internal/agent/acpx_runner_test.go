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
	fakeACPXPromptPath = "ROUNDFIX_FAKE_ACPX_PROMPT_PATH"
	fakeACPXStdout     = "ROUNDFIX_FAKE_ACPX_STDOUT"
	fakeACPXStderr     = "ROUNDFIX_FAKE_ACPX_STDERR"
	fakeACPXExitCode   = "ROUNDFIX_FAKE_ACPX_EXIT_CODE"
)

func TestMain(m *testing.M) {
	if os.Getenv(fakeACPXEnv) == "1" {
		os.Exit(runFakeACPXProcess())
	}
	os.Exit(m.Run())
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
	if path := os.Getenv(fakeACPXArgsPath); path != "" {
		payload, err := json.Marshal(os.Args[1:])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "marshal args: %v\n", err)
			return 2
		}
		if err := os.WriteFile(path, payload, 0o644); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "write args: %v\n", err)
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
	_, _ = io.WriteString(os.Stdout, os.Getenv(fakeACPXStdout))
	_, _ = io.WriteString(os.Stderr, os.Getenv(fakeACPXStderr))
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
