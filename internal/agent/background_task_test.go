package agent

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/testfixture"
)

func TestClaudeCommandsRefuseBackgroundTasks(t *testing.T) {
	t.Parallel()
	for _, runtimeID := range []string{"claude", "codex"} {
		t.Run(runtimeID, func(t *testing.T) {
			t.Parallel()
			h := newFakeACPXHarness(t)
			path := filepath.Join(h.gitRoot, "background-env.txt")
			h.setEnv("CLAUDE_CODE_DISABLE_BACKGROUND_TASKS", "0")
			h.setEnv("ROUNDFIX_BACKGROUND_ENV_PATH", path)
			h.setEnv("ROUNDFIX_BACKGROUND_TEST_BINARY", h.runner.Command)
			// The Go tool builds the wrapper; the test never writes an executable.
			// It captures the actual environment before entering the shared fake.
			h.runner.Command = testfixture.FixtureBinary(t, "background-acpx", backgroundEnvironmentCaptureSource)
			runtime := selectedRuntime(RuntimeSpec{ID: runtimeID, Protocol: ProtocolACP})
			if _, err := h.run(t.Context(), runtime, "background-env"); err != nil {
				t.Fatal(err)
			}
			session := SessionRef{Name: "background-env", WorkDir: h.gitRoot}
			adapter, err := checkAdapter(t.Context(), runtime, h.runner.Environment)
			if err != nil {
				t.Fatal(err)
			}
			h.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
				"set model": selectionStateFixture(t, "model", runtime.Model, runtime.Model, []string{runtime.Model}, "effort", "high", []string{"high"}),
			}))
			if _, err := h.runner.AcquireSelectionCapabilities(t.Context(), CapabilityAcquisitionRequest{
				Runtime: runtime, Session: session, ConfigID: "model", Value: runtime.Model, Adapter: adapter,
			}); err != nil {
				t.Fatal(err)
			}

			if err := h.runner.CancelSession(t.Context(), runtime, session); err != nil {
				t.Fatal(err)
			}
			if err := h.runner.EndSession(t.Context(), runtime, session); err != nil {
				t.Fatal(err)
			}
			if _, err := h.runner.ListRoundfixSessions(t.Context(), runtime, h.gitRoot); err != nil {
				t.Fatal(err)
			}
			rows := strings.Split(strings.TrimSpace(readFile(t, path)), "\n")
			invocations := readJSONInvocations(t, h.invocationsPath)
			if len(rows) != len(invocations) {
				t.Fatalf("environment rows=%d acpx commands=%d", len(rows), len(invocations))
			}
			want := "absent"
			if runtimeID == "claude" {
				want = "1"
			}
			commands := make(map[string]bool)
			for index, row := range rows {
				command := fakeACPXCommandKey(invocations[index])
				commands[command] = true
				if row != want {
					t.Errorf("%s background refusal=%q want=%q", command, row, want)
				}
			}
			for _, command := range []string{"sessions ensure", "prompt", "cancel", "sessions close", "sessions list", "set model"} {
				if !commands[command] {
					t.Errorf("missing acpx command %s", command)
				}
			}
			if environmentValue(h.runner.Environment, "CLAUDE_CODE_DISABLE_BACKGROUND_TASKS") != "0" {
				t.Fatal("runner base environment was changed")
			}
		})
	}
}

func TestTurnRecordsItsAgentBackgroundTasks(t *testing.T) {
	t.Parallel()
	const stream = `{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"s","update":{"sessionUpdate":"tool_call","toolCallId":"bash-call","rawInput":{"run_in_background":true,"description":"Run the load driver"},"_meta":{"claudeCode":{"toolName":"Bash"}}}}}
{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"s","update":{"sessionUpdate":"tool_call_update","toolCallId":"monitor-call","rawInput":{"description":"Watch the driver events"},"_meta":{"claudeCode":{"toolName":"Monitor"}}}}}
{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"s","update":{"sessionUpdate":"tool_call_update","toolCallId":"bash-call","status":"in_progress"}}}
{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"s","update":{"sessionUpdate":"tool_call_update","toolCallId":"bash-call","status":"completed","_meta":{"claudeCode":{"toolName":"Bash","toolResponse":{"backgroundTaskId":"bvlc7cgwx","taskId":"ignored"}}}}}}
{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"s","update":{"sessionUpdate":"tool_call_update","toolCallId":"monitor-call","status":"completed","_meta":{"claudeCode":{"toolName":"Monitor","toolResponse":{"taskId":"bekzicxjj"}}}}}}
{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"s","update":{"sessionUpdate":"tool_call_update","toolCallId":"bash-call","_meta":{"claudeCode":{"toolName":"Bash","toolResponse":{"backgroundTaskId":"bvlc7cgwx"}}}}}}
`
	for _, tc := range []struct {
		name   string
		stdout string
		want   []AgentBackgroundTask
	}{
		{"Bash and Monitor in start order", stream, []AgentBackgroundTask{
			{ID: "bvlc7cgwx", Tool: "Bash", Description: "Run the load driver"},
			{ID: "bekzicxjj", Tool: "Monitor", Description: "Watch the driver events"},
		}},
		{"no background tasks", `{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"s","update":{"sessionUpdate":"tool_call_update","toolCallId":"foreground","rawInput":{"description":"Run in foreground"},"_meta":{"claudeCode":{"toolName":"Bash","toolResponse":{"backgroundTaskId":"","taskId":""}}}}}}` + "\n", nil},
		{"missing description", `{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"s","update":{"sessionUpdate":"tool_call_update","toolCallId":"new-call","_meta":{"claudeCode":{"toolName":"Monitor","toolResponse":{"backgroundTaskId":"","taskId":"monitor-id"}}}}}}` + "\n", []AgentBackgroundTask{{ID: "monitor-id", Tool: "Monitor"}}},
		{"other update metadata is not a task", `{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"s","update":{"sessionUpdate":"session_info_update","_meta":{"claudeCode":{"toolName":"Bash","toolResponse":{"taskId":"not-started"}}}}}}` + "\n", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			run := runFakeACPXPrompt(t, fakeACPXPrompt{runtime: RuntimeSpec{ID: "claude", Protocol: ProtocolACP}, stdout: tc.stdout + acpxPromptResponseLine("end_turn")})
			if run.err != nil {
				t.Fatal(run.err)
			}
			if !reflect.DeepEqual(run.result.BackgroundTasks, tc.want) {
				t.Fatalf("tasks=%+v want=%+v", run.result.BackgroundTasks, tc.want)
			}
		})
	}
}

func TestQAPromptRequiresForegroundCommands(t *testing.T) {
	t.Parallel()
	prompt, err := BuildQAPrompt(sampleQAPromptRequest())
	if err != nil {
		t.Fatal(err)
	}
	const rule = "- Run every command in the foreground and wait for it to exit; never end the turn while a command you started is still running. The end of the turn is the end of the gate.\n"
	if count := strings.Count(prompt, rule); count != 1 {
		t.Fatalf("foreground rule count=%d want=1; prompt=%s", count, prompt)
	}
}

// The compiled wrapper owns and waits for the fake acpx child. Its data record
// contains only the fixed background-task policy, never the full environment.
const backgroundEnvironmentCaptureSource = `package main

import (
 "errors"
 "fmt"
 "os"
 "os/exec"
)

func main() {
 if err := capture(); err != nil {
  fmt.Fprintln(os.Stderr, err)
  os.Exit(2)
 }
 command := exec.Command(os.Getenv("ROUNDFIX_BACKGROUND_TEST_BINARY"), os.Args[1:]...)
 command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr
 if err := command.Run(); err != nil {
  var exit *exec.ExitError
  if errors.As(err, &exit) { os.Exit(exit.ExitCode()) }
  fmt.Fprintln(os.Stderr, err)
  os.Exit(2)
 }
}

func capture() error {
 value, present := os.LookupEnv("CLAUDE_CODE_DISABLE_BACKGROUND_TASKS")
 if !present { value = "absent" }
 file, err := os.OpenFile(os.Getenv("ROUNDFIX_BACKGROUND_ENV_PATH"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
 if err != nil { return err }
 _, writeErr := fmt.Fprintln(file, value)
 return errors.Join(writeErr, file.Close())
}
`
