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

	"roundfix/internal/codex"
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
	fakeACPXClosed     = "ROUNDFIX_FAKE_ACPX_CLOSED"
	fakeACPXStarted    = "ROUNDFIX_FAKE_ACPX_STARTED"
	fakeACPXBlock      = "ROUNDFIX_FAKE_ACPX_BLOCK_PROMPT"
	fakeACPXExitCancel = "ROUNDFIX_FAKE_ACPX_EXIT_AFTER_CANCEL"
	fakeACPXCodexPath  = "ROUNDFIX_FAKE_ACPX_CODEX_PATH"
)

func TestMain(m *testing.M) {
	if os.Getenv(fakeACPXEnv) == "1" {
		os.Exit(runFakeACPXProcess())
	}
	os.Exit(m.Run())
}

func TestACPXProbePassesWhenVersionMatchesPin(t *testing.T) {
	invocations, err := runFakeACPXProbe(t, RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, PinnedACPXVersion)

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

	invocations, err := runFakeACPXProbe(t, RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, foundVersion)

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

	invocations, err := runFakeACPXProbe(t, runtime, PinnedACPXVersion)

	if err != nil {
		t.Fatalf("expected command override probe to pass through acpx, got %v", err)
	}
	want := [][]string{{"--version"}}
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("expected command override to probe acpx only\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestInfrastructureErrorErrorIncludesBoundedStderrTail(t *testing.T) {
	const (
		stderrDelimiter       = "\n--- acpx stderr tail ---\n"
		stderrTruncatedMarker = "[stderr truncated]\n"
	)
	hundredLineStderr := numberedLinesForTest(100)
	hundredLineTail := numberedLinesRangeForTest(91, 100)
	oversizedStderr := strings.Repeat("x", 1100)

	tests := []struct {
		name            string
		stderr          string
		want            string
		wantSuffix      string
		wantNotContains []string
		wantTailBytes   int
	}{
		{
			name:   "empty stderr keeps existing message",
			stderr: " \n\t",
			want:   "acpx infrastructure error after exit code 2: usage error",
		},
		{
			name:       "multi-line stderr appended as delimited tail",
			stderr:     "\nfirst line\nsecond line\nthird line\n",
			wantSuffix: stderrDelimiter + "first line\nsecond line\nthird line",
		},
		{
			name:            "hundred-line stderr keeps last ten lines",
			stderr:          hundredLineStderr,
			wantSuffix:      stderrDelimiter + stderrTruncatedMarker + hundredLineTail,
			wantNotContains: []string{"line-090"},
		},
		{
			name:          "oversized stderr keeps last kibibyte",
			stderr:        oversizedStderr,
			wantSuffix:    stderrDelimiter + stderrTruncatedMarker + strings.Repeat("x", 1024),
			wantTailBytes: 1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := (&InfrastructureError{
				ExitCode: 2,
				Reason:   acpxExitReasonUsage,
				Stderr:   tt.stderr,
			}).Error()

			if tt.want != "" && message != tt.want {
				t.Fatalf("unexpected message\nwant: %q\ngot:  %q", tt.want, message)
			}
			if tt.wantSuffix != "" && !strings.HasSuffix(message, tt.wantSuffix) {
				t.Fatalf("expected message to end with stderr tail\nwant suffix: %q\ngot:         %q", tt.wantSuffix, message)
			}
			for _, forbidden := range tt.wantNotContains {
				if strings.Contains(message, forbidden) {
					t.Fatalf("expected message to omit %q, got %q", forbidden, message)
				}
			}
			if tt.wantTailBytes > 0 {
				tail := strings.TrimPrefix(strings.TrimPrefix(infrastructureTailFromMessageForTest(t, message), stderrTruncatedMarker), "\n")
				if len(tail) > tt.wantTailBytes {
					t.Fatalf("expected stderr tail <= %d bytes, got %d", tt.wantTailBytes, len(tail))
				}
			}
		})
	}
}

func TestACPXPromptArgsPlaceGlobalsBeforeAgentAndSubcommand(t *testing.T) {
	tests := []struct {
		name    string
		runtime RuntimeSpec
		want    []string
	}{
		{
			name:    "default adapter",
			runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP},
			want: []string{
				"--cwd", "/repo",
				"--format", "json",
				"--json-strict",
				"--approve-all",
				"codex", "prompt",
				"-s", "roundfix-run-1",
				"-f", "-",
			},
		},
		{
			name:    "model set",
			runtime: RuntimeSpec{ID: "claude", Protocol: ProtocolACP, Model: "opus-test"},
			want: []string{
				"--cwd", "/repo",
				"--format", "json",
				"--json-strict",
				"--approve-all",
				"--model", "opus-test",
				"claude", "prompt",
				"-s", "roundfix-run-1",
				"-f", "-",
			},
		},
		{
			name:    "command override",
			runtime: RuntimeSpec{ID: "codex-custom", Protocol: ProtocolStdio, Command: "custom-acp --stdio"},
			want: []string{
				"--cwd", "/repo",
				"--format", "json",
				"--json-strict",
				"--approve-all",
				"--agent", "custom-acp --stdio",
				"prompt",
				"-s", "roundfix-run-1",
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

func TestACPXCancelSessionInvokesSessionCancel(t *testing.T) {
	dir := t.TempDir()
	invocationsPath := filepath.Join(dir, "invocations.jsonl")
	t.Setenv(fakeACPXEnv, "1")
	t.Setenv(fakeACPXInvokes, invocationsPath)

	err := (&ACPXRunner{Command: os.Args[0], codexSpawn: codexSpawnDependencies{goos: "linux"}}).CancelSession(context.Background(), RuntimeSpec{
		ID:       "codex",
		Protocol: ProtocolACP,
	}, SessionRef{
		Name:    "roundfix-run-1",
		WorkDir: "/repo",
	})

	if err != nil {
		t.Fatalf("CancelSession returned error: %v", err)
	}
	want := [][]string{{
		"--cwd", "/repo",
		"codex", "cancel",
		"-s", "roundfix-run-1",
	}}
	if got := readJSONInvocations(t, invocationsPath); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected cancel invocation\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestParseRoundfixSessionsFiltersAndExtractsRunIDs(t *testing.T) {
	output := strings.Join([]string{
		"NAME                          UPDATED",
		"roundfix-run_20260706T120000Z_abcd1234  now",
		"foreign-session               now",
		"| roundfix-run_20260706T120000Z_abcd1234-task_02 | active |",
		`{"sessions":[{"name":"roundfix-run_20260706T120000Z_abcd1234-task_03"},{"name":"other"}]}`,
	}, "\n")

	got := ParseRoundfixSessions(output)
	want := []RoundfixSession{
		{Name: "roundfix-run_20260706T120000Z_abcd1234", RunID: "run_20260706T120000Z_abcd1234"},
		{Name: "roundfix-run_20260706T120000Z_abcd1234-task_02", RunID: "run_20260706T120000Z_abcd1234", TaskID: "task_02"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected sessions\nwant: %#v\ngot:  %#v", want, got)
	}

	jsonOutput := `{"sessions":[{"name":"roundfix-run_20260706T120000Z_abcd1234-task_03"},{"name":"foreign"}]}`
	got = ParseRoundfixSessions(jsonOutput)
	want = []RoundfixSession{{
		Name:   "roundfix-run_20260706T120000Z_abcd1234-task_03",
		RunID:  "run_20260706T120000Z_abcd1234",
		TaskID: "task_03",
	}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected JSON sessions\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestACPXListRoundfixSessionsInvokesSessionsList(t *testing.T) {
	harness := newFakeACPXHarness(t)
	t.Setenv(fakeACPXStdoutBy, mustJSONForTest(t, map[string]string{
		"sessions list": "NAME\nroundfix-run_20260706T120000Z_abcd1234\nforeign\n",
	}))

	got, err := harness.runner.ListRoundfixSessions(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, harness.gitRoot)

	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}
	wantSessions := []RoundfixSession{{
		Name:  "roundfix-run_20260706T120000Z_abcd1234",
		RunID: "run_20260706T120000Z_abcd1234",
	}}
	if !reflect.DeepEqual(got, wantSessions) {
		t.Fatalf("unexpected sessions\nwant: %#v\ngot:  %#v", wantSessions, got)
	}
	wantInvocations := [][]string{{"--cwd", harness.gitRoot, "codex", "sessions", "list"}}
	if invocations := readJSONInvocations(t, harness.invocationsPath); !reflect.DeepEqual(invocations, wantInvocations) {
		t.Fatalf("unexpected invocations\nwant: %#v\ngot:  %#v", wantInvocations, invocations)
	}
}

func TestACPXCloseSessionReturnsCloseFailure(t *testing.T) {
	harness := newFakeACPXHarness(t)
	t.Setenv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions close": 1}))
	t.Setenv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"sessions close": "close rejected\n"}))

	err := harness.runner.CloseSession(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, SessionRef{Name: "roundfix-run-1", WorkDir: harness.gitRoot})

	if err == nil {
		t.Fatal("expected close failure")
	}
	if !strings.Contains(err.Error(), "close acpx Agent Session") || !strings.Contains(err.Error(), "close rejected") {
		t.Fatalf("expected close error context, got %v", err)
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
		"--cwd", run.gitRoot,
		"--format", "json",
		"--json-strict",
		"--approve-all",
		"--model", "gpt-test",
		"codex", "prompt",
		"-s", "roundfix-run-1",
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

func TestACPXRunRequiresConcreteSelection(t *testing.T) {
	tests := []struct {
		name     string
		runtime  RuntimeSpec
		contains string
	}{
		{
			name:     "missing model",
			runtime:  RuntimeSpec{ID: "codex", Protocol: ProtocolACP, ReasoningEffort: "xhigh"},
			contains: "model",
		},
		{
			name:     "missing reasoning effort",
			runtime:  RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-test"},
			contains: "reasoning effort",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			harness := newFakeACPXHarness(t)
			_, err := harness.runner.Run(context.Background(), ExecuteRequest{
				Runtime: tt.runtime,
				RunID:   "run-acpx",
				Batch:   rounds.Batch{Number: 7},
				LogPath: filepath.Join(harness.gitRoot, "runs", "run-acpx", "agent", "batch-007.log"),
				Prompt:  "prompt",
				GitRoot: harness.gitRoot,
				Session: SessionRef{Name: "roundfix-run-1"},
			}, newCaptureSink(""))

			if err == nil {
				t.Fatal("expected missing selection to fail")
			}
			if !strings.Contains(err.Error(), tt.contains) {
				t.Fatalf("expected error containing %q, got %q", tt.contains, err.Error())
			}
			if _, statErr := os.Stat(harness.invocationsPath); !errors.Is(statErr, os.ErrNotExist) {
				t.Fatalf("expected no acpx invocation before selection validation, got stat error %v", statErr)
			}
		})
	}
}

func TestACPXRunAppliesSelectionBeforePrompt(t *testing.T) {
	tests := []struct {
		name    string
		runtime RuntimeSpec
		want    func(gitRoot string) [][]string
	}{
		{
			name:    "codex reasoning_effort",
			runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh"},
			want: func(gitRoot string) [][]string {
				return [][]string{
					{"--cwd", gitRoot, "--model", "gpt-5.5", "codex", "sessions", "ensure", "--name", "roundfix-run-1"},
					{"--cwd", gitRoot, "codex", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-1"},
					{"--cwd", gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "gpt-5.5", "codex", "prompt", "-s", "roundfix-run-1", "-f", "-"},
				}
			},
		},
		{
			name:    "claude effort",
			runtime: RuntimeSpec{ID: "claude", Protocol: ProtocolACP, Model: "opus", ReasoningEffort: "high"},
			want: func(gitRoot string) [][]string {
				return [][]string{
					{"--cwd", gitRoot, "--model", "opus", "claude", "sessions", "ensure", "--name", "roundfix-run-1"},
					{"--cwd", gitRoot, "claude", "set", "effort", "high", "-s", "roundfix-run-1"},
					{"--cwd", gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "opus", "claude", "prompt", "-s", "roundfix-run-1", "-f", "-"},
				}
			},
		},
		{
			name:    "opencode effort",
			runtime: RuntimeSpec{ID: "opencode", Protocol: ProtocolACP, Model: "opencode-model", ReasoningEffort: "maximum"},
			want: func(gitRoot string) [][]string {
				return [][]string{
					{"--cwd", gitRoot, "--model", "opencode-model", "opencode", "sessions", "ensure", "--name", "roundfix-run-1"},
					{"--cwd", gitRoot, "opencode", "set", "effort", "maximum", "-s", "roundfix-run-1"},
					{"--cwd", gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "opencode-model", "opencode", "prompt", "-s", "roundfix-run-1", "-f", "-"},
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

func TestACPXRunReappliesSelectionForFreshRunnerSessionResume(t *testing.T) {
	harness := newFakeACPXHarness(t)
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh"}

	if _, err := harness.run(context.Background(), runtime, "roundfix-run-1"); err != nil {
		t.Fatalf("first run: %v", err)
	}
	harness.runner = &ACPXRunner{Command: os.Args[0], codexSpawn: codexSpawnDependencies{goos: "linux"}}
	if _, err := harness.run(context.Background(), runtime, "roundfix-run-1"); err != nil {
		t.Fatalf("resumed run: %v", err)
	}

	want := [][]string{
		{"--cwd", harness.gitRoot, "--model", "gpt-5.5", "codex", "sessions", "ensure", "--name", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "gpt-5.5", "codex", "prompt", "-s", "roundfix-run-1", "-f", "-"},
		{"--cwd", harness.gitRoot, "--model", "gpt-5.5", "codex", "sessions", "ensure", "--name", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "gpt-5.5", "codex", "prompt", "-s", "roundfix-run-1", "-f", "-"},
	}
	if got := readJSONInvocations(t, harness.invocationsPath); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected resume invocations\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestACPXRunSelectionSetupErrorsPreserveAdapterFailure(t *testing.T) {
	harness := newFakeACPXHarness(t)
	t.Setenv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"set reasoning_effort": 2}))
	t.Setenv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"set reasoning_effort": "reasoning rejected\n"}))

	_, err := harness.run(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh"}, "roundfix-run-1")

	if err == nil {
		t.Fatal("expected reasoning setup failure")
	}
	var infraErr *InfrastructureError
	if !errors.As(err, &infraErr) {
		t.Fatalf("expected InfrastructureError in error chain, got %T %v", err, err)
	}
	if infraErr.Stderr != "reasoning rejected\n" {
		t.Fatalf("expected adapter stderr preserved, got %q", infraErr.Stderr)
	}
	if !strings.Contains(err.Error(), "set acpx Agent Session reasoning_effort") || !strings.Contains(err.Error(), "reasoning rejected") {
		t.Fatalf("expected reasoning operation context, got %v", err)
	}
	want := [][]string{
		{"--cwd", harness.gitRoot, "--model", "gpt-5.5", "codex", "sessions", "ensure", "--name", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-1"},
	}
	if got := readJSONInvocations(t, harness.invocationsPath); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected failed-selection invocations\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestACPXRunCustomCommandRequiresSelectionOptionContract(t *testing.T) {
	harness := newFakeACPXHarness(t)
	t.Setenv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"set reasoning_effort": 2}))
	t.Setenv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{
		"set reasoning_effort": "ACP session session-id does not advertise config option 'reasoning_effort'. Supported config options: model.\n",
	}))
	runtime := RuntimeSpec{ID: "codex-custom", Protocol: ProtocolStdio, Command: "custom-acp --stdio", Model: "gpt-custom", ReasoningEffort: "xhigh"}

	_, err := harness.run(context.Background(), runtime, "roundfix-run-1")

	if err == nil {
		t.Fatal("expected custom command without reasoning option contract to fail")
	}
	if !strings.Contains(err.Error(), "set acpx Agent Session reasoning_effort") || !strings.Contains(err.Error(), "does not advertise config option 'reasoning_effort'") {
		t.Fatalf("expected custom command option-contract failure, got %v", err)
	}
	want := [][]string{
		{"--cwd", harness.gitRoot, "--model", "gpt-custom", "--agent", "custom-acp --stdio", "sessions", "ensure", "--name", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "--agent", "custom-acp --stdio", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-1"},
	}
	if got := readJSONInvocations(t, harness.invocationsPath); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected custom-command invocations\nwant: %#v\ngot:  %#v", want, got)
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
		{"--cwd", harness.gitRoot, "--model", "gpt-test", "codex", "sessions", "ensure", "--name", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "gpt-test", "codex", "prompt", "-s", "roundfix-run-1", "-f", "-"},
		{"--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "gpt-test", "codex", "prompt", "-s", "roundfix-run-1", "-f", "-"},
		{"--cwd", harness.gitRoot, "--model", "gpt-test", "codex", "sessions", "ensure", "--name", "roundfix-run-2"},
		{"--cwd", harness.gitRoot, "codex", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-2"},
		{"--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "gpt-test", "codex", "prompt", "-s", "roundfix-run-2", "-f", "-"},
	}
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("unexpected acpx invocations\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestACPXRunCodexUsesConfiguredCleanPathOnDarwin(t *testing.T) {
	harness := newFakeACPXHarness(t)
	envPath := filepath.Join(harness.gitRoot, "codex-paths.jsonl")
	t.Setenv(fakeACPXCodexPath, envPath)
	t.Setenv(codexPathEnv, "/configured/clean/codex")
	probe := &fakeCodexSpawnProbe{
		accepted: map[string]bool{
			"/configured/clean/codex": true,
			"/path/quarantined/codex": true,
		},
		quarantined: map[string]bool{
			"/path/quarantined/codex": true,
		},
	}
	harness.runner.codexSpawn = codexSpawnDependencies{
		goos:       "darwin",
		lookPath:   fakeCodexLookPath("/path/quarantined/codex", nil),
		quarantine: probe,
		acceptance: probe,
	}

	if _, err := harness.run(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, "roundfix-run-1"); err != nil {
		t.Fatalf("run acpx with configured clean codex: %v", err)
	}

	wantEnv := []string{"/configured/clean/codex", "/configured/clean/codex", "/configured/clean/codex"}
	if got := readJSONStrings(t, envPath); !reflect.DeepEqual(got, wantEnv) {
		t.Fatalf("unexpected CODEX_PATH values\nwant: %#v\ngot:  %#v", wantEnv, got)
	}
	if got, want := probe.quarantineCalls, []string{"/configured/clean/codex"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("expected only configured codex to be inspected\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestACPXRunCodexFallsBackToCleanPathWhenConfiguredPathIsQuarantined(t *testing.T) {
	harness := newFakeACPXHarness(t)
	envPath := filepath.Join(harness.gitRoot, "codex-paths.jsonl")
	t.Setenv(fakeACPXCodexPath, envPath)
	t.Setenv(codexPathEnv, "/configured/quarantined/codex")
	probe := &fakeCodexSpawnProbe{
		accepted: map[string]bool{
			"/configured/quarantined/codex": true,
			"/path/clean/codex":             true,
		},
		quarantined: map[string]bool{
			"/configured/quarantined/codex": true,
		},
	}
	harness.runner.codexSpawn = codexSpawnDependencies{
		goos:       "darwin",
		lookPath:   fakeCodexLookPath("/path/clean/codex", nil),
		quarantine: probe,
		acceptance: probe,
	}

	if _, err := harness.run(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, "roundfix-run-1"); err != nil {
		t.Fatalf("run acpx with fallback clean codex: %v", err)
	}

	wantEnv := []string{"/path/clean/codex", "/path/clean/codex", "/path/clean/codex"}
	if got := readJSONStrings(t, envPath); !reflect.DeepEqual(got, wantEnv) {
		t.Fatalf("unexpected CODEX_PATH values\nwant: %#v\ngot:  %#v", wantEnv, got)
	}
	wantInspections := []string{"/configured/quarantined/codex", "/path/clean/codex"}
	if !reflect.DeepEqual(probe.quarantineCalls, wantInspections) {
		t.Fatalf("unexpected quarantine inspections\nwant: %#v\ngot:  %#v", wantInspections, probe.quarantineCalls)
	}
}

func TestACPXRunCodexSurfacesQuarantinedPathWithoutSpawning(t *testing.T) {
	harness := newFakeACPXHarness(t)
	t.Setenv(codexPathEnv, "")
	probe := &fakeCodexSpawnProbe{
		accepted: map[string]bool{
			"/path/quarantined/codex": true,
		},
		quarantined: map[string]bool{
			"/path/quarantined/codex": true,
		},
	}
	harness.runner.codexSpawn = codexSpawnDependencies{
		goos:       "darwin",
		lookPath:   fakeCodexLookPath("/path/quarantined/codex", nil),
		quarantine: probe,
		acceptance: probe,
	}

	_, err := harness.run(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, "roundfix-run-1")

	if err == nil {
		t.Fatal("expected quarantined codex to fail before spawning acpx")
	}
	message := err.Error()
	for _, want := range []string{"not safe for acpx launch", "quarantined", codex.ReinstallNextAction} {
		if !strings.Contains(message, want) {
			t.Fatalf("expected error to contain %q, got %q", want, message)
		}
	}
	if _, statErr := os.Stat(harness.invocationsPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("expected no acpx invocation file, got stat error %v", statErr)
	}
}

func TestACPXRunCodexInspectsOncePerSessionResolution(t *testing.T) {
	harness := newFakeACPXHarness(t)
	t.Setenv(codexPathEnv, "/configured/clean/codex")
	probe := &fakeCodexSpawnProbe{
		accepted: map[string]bool{"/configured/clean/codex": true},
	}
	harness.runner.codexSpawn = codexSpawnDependencies{
		goos:       "darwin",
		lookPath:   fakeCodexLookPath("/path/clean/codex", nil),
		quarantine: probe,
		acceptance: probe,
	}

	if _, err := harness.run(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, "roundfix-run-1"); err != nil {
		t.Fatalf("first acpx run: %v", err)
	}
	if _, err := harness.run(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, "roundfix-run-1"); err != nil {
		t.Fatalf("second acpx run: %v", err)
	}

	if got, want := probe.quarantineCalls, []string{"/configured/clean/codex"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("expected one quarantine inspection for the session\nwant: %#v\ngot:  %#v", want, got)
	}
	if got, want := probe.acceptanceCalls, []string{"/configured/clean/codex"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("expected one acceptance inspection for the session\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestACPXRunCodexLeavesNonDarwinSpawnUnchanged(t *testing.T) {
	harness := newFakeACPXHarness(t)
	probe := &fakeCodexSpawnProbe{}
	harness.runner.codexSpawn = codexSpawnDependencies{
		goos:       "linux",
		lookPath:   fakeCodexLookPath("/path/clean/codex", nil),
		quarantine: probe,
		acceptance: probe,
	}

	if _, err := harness.run(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, "roundfix-run-1"); err != nil {
		t.Fatalf("run acpx on non-darwin: %v", err)
	}

	if len(probe.quarantineCalls) != 0 || len(probe.acceptanceCalls) != 0 {
		t.Fatalf("expected non-darwin run not to inspect codex, got quarantine=%#v acceptance=%#v", probe.quarantineCalls, probe.acceptanceCalls)
	}
	want := [][]string{
		{"--cwd", harness.gitRoot, "--model", "gpt-test", "codex", "sessions", "ensure", "--name", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "gpt-test", "codex", "prompt", "-s", "roundfix-run-1", "-f", "-"},
	}
	if got := readJSONInvocations(t, harness.invocationsPath); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected non-darwin acpx invocations\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestACPXRunFailsEnsureWithStderrTail(t *testing.T) {
	harness := newFakeACPXHarness(t)
	t.Setenv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions ensure": 2}))
	t.Setenv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"sessions ensure": "ensure rejected\n"}))

	_, err := harness.run(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, "roundfix-run-1")

	if err == nil {
		t.Fatal("expected ensure failure")
	}
	var infraErr *InfrastructureError
	if !errors.As(err, &infraErr) {
		t.Fatalf("expected InfrastructureError, got %T %v", err, err)
	}
	if infraErr.Stderr != "ensure rejected\n" {
		t.Fatalf("expected captured ensure stderr, got %q", infraErr.Stderr)
	}
	if !strings.Contains(err.Error(), "ensure acpx Agent Session") || !strings.Contains(err.Error(), "--- acpx stderr tail ---\nensure rejected") {
		t.Fatalf("expected ensure error with delimited stderr tail, got %q", err.Error())
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
					{"--cwd", gitRoot, "--model", "gpt-test", "codex", "sessions", "ensure", "--name", "roundfix-run-1"},
					{"--cwd", gitRoot, "codex", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-1"},
					{"--cwd", gitRoot, "codex", "set-mode", "full-access", "-s", "roundfix-run-1"},
					{"--cwd", gitRoot, "codex", "set", "sandbox_mode", "danger-full-access", "-s", "roundfix-run-1"},
					{"--cwd", gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "gpt-test", "codex", "prompt", "-s", "roundfix-run-1", "-f", "-"},
				}
			},
		},
		{
			name:    "claude mode only",
			runtime: RuntimeSpec{ID: "claude", Protocol: ProtocolACP, FullAccessMode: "bypassPermissions"},
			want: func(gitRoot string) [][]string {
				return [][]string{
					{"--cwd", gitRoot, "--model", "opus-test", "claude", "sessions", "ensure", "--name", "roundfix-run-1"},
					{"--cwd", gitRoot, "claude", "set", "effort", "high", "-s", "roundfix-run-1"},
					{"--cwd", gitRoot, "claude", "set-mode", "bypassPermissions", "-s", "roundfix-run-1"},
					{"--cwd", gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "opus-test", "claude", "prompt", "-s", "roundfix-run-1", "-f", "-"},
				}
			},
		},
		{
			name:    "opencode no full access setup",
			runtime: RuntimeSpec{ID: "opencode", Protocol: ProtocolACP},
			want: func(gitRoot string) [][]string {
				return [][]string{
					{"--cwd", gitRoot, "--model", "opencode-test", "opencode", "sessions", "ensure", "--name", "roundfix-run-1"},
					{"--cwd", gitRoot, "opencode", "set", "effort", "high", "-s", "roundfix-run-1"},
					{"--cwd", gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "opencode-test", "opencode", "prompt", "-s", "roundfix-run-1", "-f", "-"},
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
		{"--cwd", harness.gitRoot, "--model", "gpt-test", "codex", "sessions", "ensure", "--name", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "set-mode", "full-access", "-s", "roundfix-run-1"},
	}
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("unexpected invocations after setup failure\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestACPXRunWarnsWhenCodexSandboxPresetUnavailable(t *testing.T) {
	harness := newFakeACPXHarness(t)
	t.Setenv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"set sandbox_mode": 2}))
	t.Setenv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{
		"set sandbox_mode": "ACP session session-id does not advertise config option 'sandbox_mode'. Supported config options: mode, model.\n",
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
		{"--cwd", harness.gitRoot, "--model", "gpt-test", "codex", "sessions", "ensure", "--name", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "set-mode", "full-access", "-s", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "set", "sandbox_mode", "danger-full-access", "-s", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "gpt-test", "codex", "prompt", "-s", "roundfix-run-1", "-f", "-"},
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
	if !containsInvocation(invocations, []string{"--cwd", harness.gitRoot, "codex", "cancel", "-s", "roundfix-run-1"}) {
		t.Fatalf("expected cancel invocation, got %#v", invocations)
	}
}

func TestACPXRunClosesSessionAfterCancelGracePeriod(t *testing.T) {
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
		t.Fatalf("expected fallback close after cancel, got %#v", stopErr)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if !containsInvocation(invocations, []string{"--cwd", harness.gitRoot, "codex", "cancel", "-s", "roundfix-run-1"}) {
		t.Fatalf("expected cancel invocation before close, got %#v", invocations)
	}
	if !containsInvocation(invocations, []string{"--cwd", harness.gitRoot, "codex", "sessions", "close", "roundfix-run-1"}) {
		t.Fatalf("expected close invocation after cancel grace, got %#v", invocations)
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

	err := harness.runner.EndSession(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, SessionRef{Name: "roundfix-run-1", WorkDir: harness.gitRoot})

	if err != nil {
		t.Fatalf("expected close failure to be swallowed, got %v", err)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	want := [][]string{{"--cwd", harness.gitRoot, "codex", "sessions", "close", "roundfix-run-1"}}
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("unexpected close invocation\nwant: %#v\ngot:  %#v", want, invocations)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "roundfix-run-1") || !strings.Contains(warnings[0], "close rejected") {
		t.Fatalf("expected close warning with session and stderr, got %#v", warnings)
	}
	if !strings.Contains(warnings[0], "--- acpx stderr tail ---\nclose rejected") {
		t.Fatalf("expected close warning with delimited stderr tail, got %#v", warnings)
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

func TestACPXRunPromptAllowsEmptyLogPathAndStillJournals(t *testing.T) {
	messageLine := acpxUpdateLine(`{"sessionId":"sess-1","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"hello"}}}`)
	responseLine := acpxPromptResponseLine("end_turn")

	run := runFakeACPXPrompt(t, fakeACPXPrompt{
		stdout:     messageLine + responseLine,
		withoutLog: true,
	})

	if run.err != nil {
		t.Fatalf("run fake acpx without Agent log path: %v", run.err)
	}
	if run.result.LogPath != "" {
		t.Fatalf("expected empty Agent log path in result, got %q", run.result.LogPath)
	}
	matches, err := filepath.Glob(filepath.Join(run.gitRoot, "runs", "*", "agent", "*.log"))
	if err != nil {
		t.Fatalf("glob Agent logs: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected no Agent log files, got %#v", matches)
	}
	events := run.sink.Events()
	if len(events) != 1 || events[0].Kind != runevent.KindAgentMessage {
		t.Fatalf("expected Agent message event without log file, got %+v", events)
	}
	if !strings.Contains(string(events[0].Payload), "hello") {
		t.Fatalf("expected raw Agent payload preserved, got %s", events[0].Payload)
	}
}

func TestACPXPromptExitClassificationMatrix(t *testing.T) {
	longStderr := numberedLinesForTest(12)
	bufferErrorLine := `{"jsonrpc":"2.0","error":{"code":-32603,"message":"Message buffer exceeded 10485760 bytes","data":{"acpxCode":"RUNTIME"}}}` + "\n"

	tests := []struct {
		name        string
		stdout      string
		stderr      string
		exitCode    int
		wantStop    string
		wantAnomaly bool
		assertErr   func(t *testing.T, err error)
	}{
		{
			name:     "result exit zero",
			stdout:   acpxPromptResponseLine("end_turn"),
			exitCode: 0,
			wantStop: "end_turn",
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("expected success, got %v", err)
				}
			},
		},
		{
			name:        "result exit one becomes anomaly success",
			stdout:      acpxPromptResponseLine("end_turn"),
			stderr:      longStderr,
			exitCode:    1,
			wantStop:    "end_turn",
			wantAnomaly: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("expected anomaly success, got %T %v", err, err)
				}
			},
		},
		{
			name:     "no result exit one remains Batch failure",
			exitCode: 1,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var batchErr *BatchFailureError
				if !errors.As(err, &batchErr) {
					t.Fatalf("expected BatchFailureError, got %T %v", err, err)
				}
				if got := err.Error(); got != "Agent Batch failed after acpx exited with code 1: agent/protocol error" {
					t.Fatalf("expected byte-identical Batch failure, got %q", got)
				}
			},
		},
		{
			name:     "no result exit two remains infrastructure",
			stderr:   "usage details\n",
			exitCode: 2,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var infraErr *InfrastructureError
				if !errors.As(err, &infraErr) {
					t.Fatalf("expected InfrastructureError, got %T %v", err, err)
				}
				if got := err.Error(); got != "acpx infrastructure error after exit code 2: usage error\n--- acpx stderr tail ---\nusage details" {
					t.Fatalf("expected byte-identical infrastructure error, got %q", got)
				}
			},
		},
		{
			name:     "no result exit four remains infrastructure",
			exitCode: 4,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var infraErr *InfrastructureError
				if !errors.As(err, &infraErr) {
					t.Fatalf("expected InfrastructureError, got %T %v", err, err)
				}
				if got := err.Error(); got != "acpx infrastructure error after exit code 4: missing session" {
					t.Fatalf("expected byte-identical infrastructure error, got %q", got)
				}
			},
		},
		{
			name:        "result exit one with buffer error line becomes anomaly success",
			stdout:      acpxPromptResponseLine("end_turn") + bufferErrorLine,
			stderr:      "Message buffer exceeded 10485760 bytes\n",
			exitCode:    1,
			wantStop:    "end_turn",
			wantAnomaly: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("expected incident-shaped anomaly success, got %T %v", err, err)
				}
			},
		},
		{
			name:     "result exit one thirty remains stop",
			stdout:   acpxPromptResponseLine("end_turn"),
			exitCode: 130,
			wantStop: "end_turn",
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				if !IsStopError(err) {
					t.Fatalf("expected StopError, got %T %v", err, err)
				}
			},
		},
		{
			name:     "partial stream exit one remains Batch failure",
			stdout:   `{"jsonrpc":"2.0","id":1,"result":`,
			exitCode: 1,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var batchErr *BatchFailureError
				if !errors.As(err, &batchErr) {
					t.Fatalf("expected BatchFailureError, got %T %v", err, err)
				}
				if got := err.Error(); got != "Agent Batch failed after acpx exited with code 1: agent/protocol error" {
					t.Fatalf("expected byte-identical Batch failure, got %q", got)
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
			if run.result.StopReason != tt.wantStop {
				t.Fatalf("expected stop reason %q, got %q", tt.wantStop, run.result.StopReason)
			}
			if tt.wantAnomaly {
				if run.result.TransportAnomaly == "" {
					t.Fatal("expected transport anomaly on successful result")
				}
				if !strings.Contains(run.result.TransportAnomaly, fmt.Sprintf("exit code %d", tt.exitCode)) {
					t.Fatalf("expected anomaly to include exit code, got %q", run.result.TransportAnomaly)
				}
				if !strings.Contains(run.result.TransportAnomaly, "line-012") && !strings.Contains(run.result.TransportAnomaly, "Message buffer exceeded 10485760 bytes") {
					t.Fatalf("expected anomaly to include bounded stderr tail, got %q", run.result.TransportAnomaly)
				}
				if strings.Contains(run.result.TransportAnomaly, "line-002") {
					t.Fatalf("expected anomaly stderr tail to be bounded, got %q", run.result.TransportAnomaly)
				}
			} else if run.result.TransportAnomaly != "" {
				t.Fatalf("expected no transport anomaly, got %q", run.result.TransportAnomaly)
			}
		})
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
				if infraErr.Stderr != "bad acpx usage\n" {
					t.Fatalf("expected captured prompt stderr, got %q", infraErr.Stderr)
				}
				if !strings.Contains(err.Error(), "--- acpx stderr tail ---\nbad acpx usage") {
					t.Fatalf("expected delimited stderr tail in error context, got %q", err.Error())
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
	runtime    RuntimeSpec
	prompt     string
	stdout     string
	stderr     string
	exitCode   int
	withoutLog bool
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
			codexSpawn: codexSpawnDependencies{goos: "linux"},
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
	t.Setenv(fakeACPXClosed, filepath.Join(harness.gitRoot, "closed"))
	if exitAfterCancel {
		t.Setenv(fakeACPXExitCancel, "1")
	}
	return harness
}

func selectedRuntime(runtime RuntimeSpec) RuntimeSpec {
	if runtime.ID == "" && runtime.Protocol == "" && runtime.Command == "" {
		runtime = RuntimeSpec{ID: "codex", Protocol: ProtocolACP}
	}
	switch {
	case strings.HasPrefix(runtime.ID, "claude"):
		if runtime.Model == "" {
			runtime.Model = "opus-test"
		}
		if runtime.ReasoningEffort == "" {
			runtime.ReasoningEffort = "high"
		}
	case strings.HasPrefix(runtime.ID, "opencode"):
		if runtime.Model == "" {
			runtime.Model = "opencode-test"
		}
		if runtime.ReasoningEffort == "" {
			runtime.ReasoningEffort = "high"
		}
	default:
		if runtime.Model == "" {
			runtime.Model = "gpt-test"
		}
		if runtime.ReasoningEffort == "" {
			runtime.ReasoningEffort = "xhigh"
		}
	}
	return runtime
}

func (harness *fakeACPXHarness) run(ctx context.Context, runtime RuntimeSpec, sessionName string) (ExecuteResult, error) {
	return harness.runWithSink(ctx, runtime, sessionName, newCaptureSink(""))
}

func (harness *fakeACPXHarness) runWithSink(ctx context.Context, runtime RuntimeSpec, sessionName string, sink runevent.Sink) (ExecuteResult, error) {
	return harness.runner.Run(ctx, ExecuteRequest{
		Runtime:   selectedRuntime(runtime),
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
	if prompt.withoutLog {
		logPath = ""
	}
	runner := &ACPXRunner{
		Command: os.Args[0],
		Now: func() time.Time {
			return time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
		},
		codexSpawn: codexSpawnDependencies{goos: "linux"},
	}
	result, err := runner.RunPrompt(context.Background(), ACPXPromptRequest{
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

func runFakeACPXProbe(t *testing.T, runtime RuntimeSpec, version string) ([][]string, error) {
	t.Helper()
	dir := t.TempDir()
	invocationsPath := filepath.Join(dir, "invocations.jsonl")
	t.Setenv(fakeACPXEnv, "1")
	t.Setenv(fakeACPXInvokes, invocationsPath)
	t.Setenv(fakeACPXStdout, version+"\n")

	err := (ACPXRunner{Command: os.Args[0]}).Probe(context.Background(), runtime)
	return readJSONInvocations(t, invocationsPath), err
}

type fakeCodexSpawnProbe struct {
	quarantined     map[string]bool
	accepted        map[string]bool
	quarantineCalls []string
	acceptanceCalls []string
}

func (probe *fakeCodexSpawnProbe) Quarantined(_ context.Context, path string) (bool, error) {
	probe.quarantineCalls = append(probe.quarantineCalls, path)
	return probe.quarantined[path], nil
}

func (probe *fakeCodexSpawnProbe) Accepted(_ context.Context, path string) (bool, error) {
	probe.acceptanceCalls = append(probe.acceptanceCalls, path)
	accepted, ok := probe.accepted[path]
	if !ok {
		return true, nil
	}
	return accepted, nil
}

func fakeCodexLookPath(path string, err error) codex.LookPathFunc {
	return func(command string) (string, error) {
		if command != codex.BinaryName {
			return "", fmt.Errorf("unexpected look path command %q", command)
		}
		return path, err
	}
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

func readJSONStrings(t *testing.T, path string) []string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read strings: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	values := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var value string
		if err := json.Unmarshal([]byte(line), &value); err != nil {
			t.Fatalf("decode string %q: %v", line, err)
		}
		values = append(values, value)
	}
	return values
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
	if path := os.Getenv(fakeACPXCodexPath); path != "" {
		if err := appendFakeACPXString(path, os.Getenv(codexPathEnv)); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "append codex path: %v\n", err)
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
	if commandKey == "sessions close" {
		if path := os.Getenv(fakeACPXClosed); path != "" {
			if err := os.WriteFile(path, []byte("closed\n"), 0o644); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "write close marker: %v\n", err)
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
			if closed := os.Getenv(fakeACPXClosed); closed != "" {
				if _, err := os.Stat(closed); err == nil {
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

func appendFakeACPXString(path string, value string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, err := file.Write(append(payload, '\n')); err != nil {
		return err
	}
	return nil
}

// fakeACPXCommandKey mirrors the real acpx grammar: program-level global
// options come first (--cwd, --format, --model, --agent, ... take a value;
// --json-strict and --approve-all are booleans), then the agent name (unless
// --agent supplied the raw command), then the subcommand.
func fakeACPXCommandKey(args []string) string {
	valueGlobals := map[string]bool{
		"--cwd":     true,
		"--format":  true,
		"--model":   true,
		"--agent":   true,
		"--timeout": true,
		"--ttl":     true,
	}
	booleanGlobals := map[string]bool{
		"--json-strict": true,
		"--approve-all": true,
	}
	sawAgentOverride := false
	index := 0
globals:
	for index < len(args) {
		switch {
		case valueGlobals[args[index]]:
			if args[index] == "--agent" {
				sawAgentOverride = true
			}
			index += 2
		case booleanGlobals[args[index]]:
			index++
		default:
			break globals
		}
	}
	if !sawAgentOverride {
		index++ // skip the positional agent name
	}
	if index >= len(args) {
		return ""
	}
	if args[index] == "sessions" && index+1 < len(args) {
		return "sessions " + args[index+1]
	}
	if args[index] == "set" && index+1 < len(args) {
		return "set " + args[index+1]
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

func numberedLinesForTest(count int) string {
	return numberedLinesRangeForTest(1, count) + "\n"
}

func numberedLinesRangeForTest(start int, end int) string {
	var builder strings.Builder
	for i := start; i <= end; i++ {
		if builder.Len() > 0 {
			builder.WriteByte('\n')
		}
		fmt.Fprintf(&builder, "line-%03d", i)
	}
	return builder.String()
}

func infrastructureTailFromMessageForTest(t *testing.T, message string) string {
	t.Helper()
	const delimiter = "\n--- acpx stderr tail ---\n"
	_, tail, ok := strings.Cut(message, delimiter)
	if !ok {
		t.Fatalf("message has no stderr tail delimiter: %q", message)
	}
	return tail
}
