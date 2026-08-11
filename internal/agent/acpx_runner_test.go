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
	"sync"
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
	fakeACPXPrompts    = "ROUNDFIX_FAKE_ACPX_PROMPTS"
	fakeACPXStdout     = "ROUNDFIX_FAKE_ACPX_STDOUT"
	fakeACPXStderr     = "ROUNDFIX_FAKE_ACPX_STDERR"
	fakeACPXExitCode   = "ROUNDFIX_FAKE_ACPX_EXIT_CODE"
	fakeACPXStdoutBy   = "ROUNDFIX_FAKE_ACPX_STDOUT_BY_COMMAND"
	fakeACPXStdoutCall = "ROUNDFIX_FAKE_ACPX_STDOUT_BY_CALL"
	fakeACPXStderrBy   = "ROUNDFIX_FAKE_ACPX_STDERR_BY_COMMAND"
	fakeACPXExitBy     = "ROUNDFIX_FAKE_ACPX_EXIT_BY_COMMAND"
	fakeACPXExitByCall = "ROUNDFIX_FAKE_ACPX_EXIT_BY_CALL"
	fakeACPXCanceled   = "ROUNDFIX_FAKE_ACPX_CANCELED"
	fakeACPXClosed     = "ROUNDFIX_FAKE_ACPX_CLOSED"
	fakeACPXPromptDone = "ROUNDFIX_FAKE_ACPX_PROMPT_DONE"
	fakeACPXStarted    = "ROUNDFIX_FAKE_ACPX_STARTED"
	fakeACPXStartEvent = "ROUNDFIX_FAKE_ACPX_STARTED_RUN_EVENT"
	fakeACPXBlock      = "ROUNDFIX_FAKE_ACPX_BLOCK_PROMPT"
	fakeACPXBlockCmd   = "ROUNDFIX_FAKE_ACPX_BLOCK_COMMAND"
	fakeACPXExitCancel = "ROUNDFIX_FAKE_ACPX_EXIT_AFTER_CANCEL"
	fakeACPXCodexPath  = "ROUNDFIX_FAKE_ACPX_CODEX_PATH"
	fakeACPXThoughtLen = "ROUNDFIX_FAKE_ACPX_THOUGHT_LENGTH"
)

func TestMain(m *testing.M) {
	if os.Getenv(fakeACPXEnv) == "1" {
		os.Exit(runFakeACPXProcess())
	}
	os.Exit(m.Run())
}

func TestACPXProbePassesWhenVersionMatchesPin(t *testing.T) {
	t.Parallel()

	invocations, err := runFakeACPXProbe(t, RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, MinimumACPXVersion)

	if err != nil {
		t.Fatalf("expected matching acpx version to pass, got %v", err)
	}
	want := [][]string{{"--version"}}
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("unexpected probe invocations\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestACPXProbeAcceptsNewerVersion(t *testing.T) {
	t.Parallel()

	const foundVersion = "0.12.1"

	invocations, err := runFakeACPXProbe(t, RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, foundVersion)

	if err != nil {
		t.Fatalf("expected newer acpx version to pass, got %v", err)
	}
	want := [][]string{{"--version"}}
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("unexpected probe invocations\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestACPXProbeRejectsMalformedVersion(t *testing.T) {
	t.Parallel()

	const foundVersion = "not-semver"

	_, err := runFakeACPXProbe(t, RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, foundVersion)

	if err == nil {
		t.Fatal("expected malformed acpx version to fail")
	}
	message := err.Error()
	for _, want := range []string{foundVersion, MinimumACPXVersion, acpxInstallCommand()} {
		if !strings.Contains(message, want) {
			t.Fatalf("expected probe error to contain %q, got %q", want, message)
		}
	}
}

func TestACPXProbeMissingBinaryNamesInstallCommand(t *testing.T) {
	t.Parallel()

	err := (ACPXRunner{Command: filepath.Join(t.TempDir(), "missing-acpx")}).Probe(context.Background(), ProbeRequest{
		Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP},
	})

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
	t.Parallel()

	const foundVersion = "0.11.0"

	invocations, err := runFakeACPXProbe(t, RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, foundVersion)

	if err == nil {
		t.Fatal("expected mismatched acpx version to fail")
	}
	message := err.Error()
	for _, want := range []string{foundVersion, MinimumACPXVersion, acpxInstallCommand(), "or newer", "upgrade"} {
		if !strings.Contains(message, want) {
			t.Fatalf("expected probe error to contain %q, got %q", want, message)
		}
	}
	if strings.Contains(message, "downgrade") {
		t.Fatalf("expected probe error not to recommend a downgrade, got %q", message)
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
	t.Parallel()

	runtime := RuntimeSpec{ID: "codex-custom", Protocol: ProtocolStdio, Command: "custom-acp --stdio"}

	invocations, err := runFakeACPXProbe(t, runtime, MinimumACPXVersion)

	if err != nil {
		t.Fatalf("expected command override probe to pass through acpx, got %v", err)
	}
	want := [][]string{{"--version"}}
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("expected command override to probe acpx only\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestResolveAdapterCommandUsesConfigFallbacksAndOverrides(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		runtime RuntimeSpec
		config  string
		want    string
	}{
		{
			name:    "configured agent command",
			runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP},
			config:  `{"agents":{"codex":{"command":"local-codex-acp","args":["--stdio"]}}}`,
			want:    "local-codex-acp --stdio",
		},
		{
			name:    "missing config falls back to default",
			runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP},
			want:    "npx -y @agentclientprotocol/codex-acp",
		},
		{
			name:    "malformed config falls back to default",
			runtime: RuntimeSpec{ID: "claude", Protocol: ProtocolACP},
			config:  `{not json`,
			want:    ClaudeAdapterCommand(),
		},
		{
			name:    "stdio command override",
			runtime: RuntimeSpec{ID: "codex-custom", Protocol: ProtocolStdio, Command: "custom-acp --stdio"},
			config:  `{"agents":{"codex-custom":{"command":"ignored-acp"}}}`,
			want:    "custom-acp --stdio",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			environment := environmentForTest("HOME=" + t.TempDir())
			if tt.config != "" {
				environment = writeACPXConfigForTest(t, tt.config)
			}

			got, err := resolveAdapterCommandWithEnv(tt.runtime, environment)

			if err != nil {
				t.Fatalf("resolve adapter command: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected adapter %q, got %q", tt.want, got)
			}
		})
	}
}

func TestACPXConfigPathUsesOnlyTheExplicitEnvironment(t *testing.T) {
	t.Parallel()

	homeDir := t.TempDir()
	profileDir := t.TempDir()
	tests := []struct {
		name        string
		environment []string
		wantHome    string
		wantErr     string
	}{
		{
			name:        "HOME takes precedence",
			environment: environmentForBase(nil, "HOME="+homeDir, "USERPROFILE="+profileDir),
			wantHome:    homeDir,
		},
		{
			name:        "USERPROFILE supplies the home",
			environment: environmentForBase(nil, "USERPROFILE="+profileDir),
			wantHome:    profileDir,
		},
		{
			name:        "missing explicit home",
			environment: environmentForBase(nil, "PATH="+os.Getenv("PATH")),
			wantErr:     "HOME or USERPROFILE",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := acpxConfigPath(test.environment)
			if test.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("acpx config path error = %v, want diagnostic containing %q", err, test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolve acpx config path: %v", err)
			}
			want := filepath.Join(test.wantHome, ".acpx", "config.json")
			if got != want {
				t.Fatalf("acpx config path = %q, want %q", got, want)
			}
		})
	}
}

func TestCheckAdapterProvesOfficialCodexPackageAndVersion(t *testing.T) {
	t.Parallel()

	command := installFakeVersionAdapter(t, CodexAdapterPackage+" "+PinnedCodexAdapterVersion)
	environment := writeACPXConfigForTest(t, fmt.Sprintf(`{"agents":{"codex":{"command":%q}}}`, command))

	evidence, err := checkAdapter(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, environment)

	if err != nil {
		t.Fatalf("check official Codex adapter: %v", err)
	}
	if evidence.Command != command || evidence.Package != CodexAdapterPackage || evidence.Version != PinnedCodexAdapterVersion {
		t.Fatalf("unexpected adapter evidence: %#v", evidence)
	}
}

func TestCheckAdapterProvesCustomCodexCommandIdentity(t *testing.T) {
	t.Parallel()

	command := installFakeVersionAdapter(t, CodexAdapterPackage+" "+PinnedCodexAdapterVersion)
	runtime := RuntimeSpec{ID: "codex-custom", Protocol: ProtocolStdio, Command: command}

	evidence, err := CheckAdapter(context.Background(), runtime)

	if err != nil {
		t.Fatalf("check custom official Codex adapter: %v", err)
	}
	if evidence.Command != command || evidence.Package != CodexAdapterPackage || evidence.Version != PinnedCodexAdapterVersion {
		t.Fatalf("unexpected custom adapter evidence: %#v", evidence)
	}
}

func TestCheckAdapterClassifiesUnreadyCodexAdapters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		output         string
		wantPackage    string
		wantVersion    string
		wantLineageErr bool
		wantVersionErr bool
	}{
		{
			name:           "legacy lineage",
			output:         "@zed-industries/codex-acp 0.16.0",
			wantPackage:    "@zed-industries/codex-acp",
			wantVersion:    "0.16.0",
			wantLineageErr: true,
		},
		{
			name:           "unknown same-named executable",
			output:         "codex-acp 1.1.5 SECRET_TOKEN=must-not-leak",
			wantPackage:    "",
			wantVersion:    "",
			wantLineageErr: true,
		},
		{
			name:           "unsupported official version",
			output:         "@agentclientprotocol/codex-acp 1.1.4",
			wantPackage:    CodexAdapterPackage,
			wantVersion:    "1.1.4",
			wantVersionErr: true,
		},
		{
			name:           "malformed official version keeps version classification",
			output:         "@agentclientprotocol/codex-acp invalid",
			wantPackage:    CodexAdapterPackage,
			wantVersion:    "invalid",
			wantVersionErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := installFakeVersionAdapter(t, tt.output)
			environment := writeACPXConfigForTest(t, fmt.Sprintf(`{"agents":{"codex":{"command":%q}}}`, command))

			_, err := checkAdapter(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, environment)

			if err == nil {
				t.Fatal("expected adapter readiness failure")
			}
			var lineageErr *AdapterLineageError
			if got := errors.As(err, &lineageErr); got != tt.wantLineageErr {
				t.Fatalf("lineage error = %t, want %t: %v", got, tt.wantLineageErr, err)
			}
			if lineageErr != nil && lineageErr.Classification() != AdapterLineageUnknown {
				t.Fatalf("unexpected lineage classification %q", lineageErr.Classification())
			}
			var versionErr *AdapterVersionError
			if got := errors.As(err, &versionErr); got != tt.wantVersionErr {
				t.Fatalf("version error = %t, want %t: %v", got, tt.wantVersionErr, err)
			}
			if versionErr != nil && versionErr.Classification() != AdapterVersionUnsupported {
				t.Fatalf("unexpected version classification %q", versionErr.Classification())
			}
			if !strings.Contains(err.Error(), command) || !strings.Contains(err.Error(), CodexAdapterInstallCommand()) {
				t.Fatalf("expected effective command and deterministic action, got %q", err)
			}
			if tt.wantPackage != "" && !strings.Contains(err.Error(), tt.wantPackage) {
				t.Fatalf("expected observed package %q, got %q", tt.wantPackage, err)
			}
			if tt.wantVersion != "" && !strings.Contains(err.Error(), tt.wantVersion) {
				t.Fatalf("expected observed version %q, got %q", tt.wantVersion, err)
			}
			if strings.Contains(err.Error(), "SECRET_TOKEN") {
				t.Fatalf("adapter output leaked into diagnostic: %q", err)
			}
		})
	}
}

func TestCheckAdapterProvesOfficialClaudePackageAndVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		runtime     func(*testing.T) RuntimeSpec
		wantVersion string
	}{
		{
			name: "version only with command package identity",
			runtime: func(t *testing.T) RuntimeSpec {
				t.Helper()
				command := installFakeNamedVersionAdapter(t, "npx", PinnedClaudeAdapterVersion)
				return RuntimeSpec{
					ID:       "claude-custom",
					Protocol: ProtocolStdio,
					Command:  command + " -y " + ClaudeAdapterPackage + "@" + PinnedClaudeAdapterVersion,
				}
			},
			wantVersion: PinnedClaudeAdapterVersion,
		},
		{
			name: "newer version",
			runtime: func(t *testing.T) RuntimeSpec {
				t.Helper()
				const newerVersion = "0.64.0"
				command := installFakeNamedVersionAdapter(t, "npx", newerVersion)
				return RuntimeSpec{
					ID:       "claude-custom",
					Protocol: ProtocolStdio,
					Command:  command + " -y " + ClaudeAdapterPackage + "@" + newerVersion,
				}
			},
			wantVersion: "0.64.0",
		},
		{
			name: "version only with resolved path package identity",
			runtime: func(t *testing.T) RuntimeSpec {
				t.Helper()
				command := installSymlinkedPackageAdapter(t, ClaudeAdapterPackage, "claude-agent-acp", PinnedClaudeAdapterVersion)
				return RuntimeSpec{ID: "claude-custom", Protocol: ProtocolStdio, Command: command}
			},
			wantVersion: PinnedClaudeAdapterVersion,
		},
		{
			name: "two field probe",
			runtime: func(t *testing.T) RuntimeSpec {
				t.Helper()
				command := installFakeNamedVersionAdapter(t, "claude-agent-acp", ClaudeAdapterPackage+" "+PinnedClaudeAdapterVersion)
				return RuntimeSpec{ID: "claude-custom", Protocol: ProtocolStdio, Command: command}
			},
			wantVersion: PinnedClaudeAdapterVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runtime := tt.runtime(t)

			evidence, err := CheckAdapter(context.Background(), runtime)

			if err != nil {
				t.Fatalf("check official Claude adapter: %v", err)
			}
			if evidence.Command != runtime.Command || evidence.Package != ClaudeAdapterPackage || evidence.Version != tt.wantVersion {
				t.Fatalf("unexpected Claude adapter evidence: %#v", evidence)
			}
		})
	}
}

func TestCheckAdapterClassifiesUnreadyClaudeAdapters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		runtime         func(*testing.T) RuntimeSpec
		wantPackage     string
		wantVersion     string
		wantLineageErr  bool
		wantVersionErr  bool
		wantMessagePart string
	}{
		{
			name: "unrecognized package lineage with empty probe",
			runtime: func(t *testing.T) RuntimeSpec {
				t.Helper()
				command := installSymlinkedPackageAdapter(t, "@example/claude-code-acp", "claude-code-acp", "")
				return RuntimeSpec{ID: "claude-custom", Protocol: ProtocolStdio, Command: command}
			},
			wantLineageErr:  true,
			wantMessagePart: "did not prove required package lineage",
		},
		{
			name: "unrecognized package lineage at the pinned version",
			runtime: func(t *testing.T) RuntimeSpec {
				t.Helper()
				command := installSymlinkedPackageAdapter(t, "@example/claude-agent-acp", "claude-agent-acp", PinnedClaudeAdapterVersion)
				return RuntimeSpec{ID: "claude-custom", Protocol: ProtocolStdio, Command: command}
			},
			wantLineageErr:  true,
			wantMessagePart: "did not prove required package lineage",
		},
		{
			name: "unproven bare executable",
			runtime: func(t *testing.T) RuntimeSpec {
				t.Helper()
				command := installFakeNamedVersionAdapter(t, "claude-agent-acp", PinnedClaudeAdapterVersion)
				return RuntimeSpec{ID: "claude-custom", Protocol: ProtocolStdio, Command: command}
			},
			wantVersion:     PinnedClaudeAdapterVersion,
			wantLineageErr:  true,
			wantMessagePart: "did not prove",
		},
		{
			name: "unsupported official version",
			runtime: func(t *testing.T) RuntimeSpec {
				t.Helper()
				command := installFakeNamedVersionAdapter(t, "npx", "0.62.9")
				return RuntimeSpec{
					ID:       "claude-custom",
					Protocol: ProtocolStdio,
					Command:  command + " -y " + ClaudeAdapterPackage + "@0.62.9",
				}
			},
			wantPackage:    ClaudeAdapterPackage,
			wantVersion:    "0.62.9",
			wantVersionErr: true,
		},
		{
			name: "malformed version probe",
			runtime: func(t *testing.T) RuntimeSpec {
				t.Helper()
				command := installFakeNamedVersionAdapter(t, "npx", "not-a-version")
				return RuntimeSpec{
					ID:       "claude-custom",
					Protocol: ProtocolStdio,
					Command:  command + " -y " + ClaudeAdapterPackage + "@" + PinnedClaudeAdapterVersion,
				}
			},
			wantPackage:     ClaudeAdapterPackage,
			wantLineageErr:  true,
			wantMessagePart: "did not prove",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runtime := tt.runtime(t)

			_, err := CheckAdapter(context.Background(), runtime)

			if err == nil {
				t.Fatal("expected Claude adapter readiness failure")
			}
			var lineageErr *AdapterLineageError
			if got := errors.As(err, &lineageErr); got != tt.wantLineageErr {
				t.Fatalf("lineage error = %t, want %t: %v", got, tt.wantLineageErr, err)
			}
			if lineageErr != nil {
				if lineageErr.Classification() != AdapterLineageUnknown {
					t.Fatalf("unexpected lineage classification %q", lineageErr.Classification())
				}
				if lineageErr.RequiredPackage != ClaudeAdapterPackage ||
					lineageErr.RequiredVersion != PinnedClaudeAdapterVersion ||
					lineageErr.Install != ClaudeAdapterInstallCommand() {
					t.Fatalf("unexpected Claude lineage contract data: %#v", lineageErr)
				}
			}
			var versionErr *AdapterVersionError
			if got := errors.As(err, &versionErr); got != tt.wantVersionErr {
				t.Fatalf("version error = %t, want %t: %v", got, tt.wantVersionErr, err)
			}
			if versionErr != nil {
				if versionErr.Classification() != AdapterVersionUnsupported {
					t.Fatalf("unexpected version classification %q", versionErr.Classification())
				}
				if versionErr.RequiredPackage != ClaudeAdapterPackage ||
					versionErr.RequiredVersion != PinnedClaudeAdapterVersion ||
					versionErr.Install != ClaudeAdapterInstallCommand() {
					t.Fatalf("unexpected Claude version contract data: %#v", versionErr)
				}
			}
			for _, want := range []string{runtime.Command, ClaudeAdapterInstallCommand(), tt.wantPackage, tt.wantVersion, tt.wantMessagePart} {
				if want != "" && !strings.Contains(err.Error(), want) {
					t.Fatalf("expected error to contain %q, got %q", want, err)
				}
			}
		})
	}
}

func TestCheckAdapterUsesOfficialClaudeInstallHints(t *testing.T) {
	t.Parallel()

	for _, command := range []string{"claude-code-acp", "claude-agent-acp"} {
		t.Run(command, func(t *testing.T) {
			if got := adapterInstallCommand(command); got != ClaudeAdapterInstallCommand() {
				t.Fatalf("install command for %q = %q, want %q", command, got, ClaudeAdapterInstallCommand())
			}
		})
	}
}

func TestCheckAdapterPreservesOpenCodeResolutionWithoutVersionExecution(t *testing.T) {
	t.Parallel()

	command := installFakeVersionAdapter(t, "must not be inspected")
	runtime := RuntimeSpec{ID: "opencode-custom", Protocol: ProtocolStdio, Command: command + " --stdio"}

	evidence, err := CheckAdapter(context.Background(), runtime)

	if err != nil {
		t.Fatalf("check OpenCode adapter: %v", err)
	}
	if evidence.Command != runtime.Command || evidence.Package != "" || evidence.Version != "" {
		t.Fatalf("unexpected OpenCode evidence: %#v", evidence)
	}
}

func TestACPXProbeMissingAdapterNamesInstallCommandBeforeSession(t *testing.T) {
	// Sequential: verifies missing-adapter lookup against the process PATH default.
	harness := newFakeACPXHarness(t)
	environment := writeACPXConfigForTest(t, `{"agents":{"codex":{"command":"codex-acp"}}}`)
	harness.setEnv("HOME", environmentValue(environment, "HOME"))
	t.Setenv("PATH", t.TempDir())
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")

	err := harness.runner.Probe(context.Background(), ProbeRequest{
		Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5"},
		WorkDir: harness.gitRoot,
	})

	if err == nil {
		t.Fatal("expected missing adapter to fail")
	}
	message := err.Error()
	for _, want := range []string{"codex-acp is required but was not found on PATH", CodexAdapterInstallCommand()} {
		if !strings.Contains(message, want) {
			t.Fatalf("expected probe error to contain %q, got %q", want, message)
		}
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	wantInvocations := [][]string{{"--version"}}
	if !reflect.DeepEqual(invocations, wantInvocations) {
		t.Fatalf("expected adapter failure before disposable Agent Session\nwant: %#v\ngot:  %#v", wantInvocations, invocations)
	}
}

func TestACPXProbeMalformedConfigFallsBackToDefaultAdapter(t *testing.T) {
	// Sequential: verifies default-adapter lookup against the process PATH default.
	harness := newFakeACPXHarness(t)
	environment := writeACPXConfigForTest(t, `{not json`)
	harness.setEnv("HOME", environmentValue(environment, "HOME"))
	t.Setenv("PATH", harness.adapterDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": sessionCapabilitySnapshotFixture(t, "gpt-5.5", []string{"gpt-5.5"}, "", "", nil),
	}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5"}

	err := harness.runner.Probe(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})

	if err != nil {
		t.Fatalf("expected malformed config to fall back to default adapter, got %v", err)
	}
}

func TestACPXProbeValidatesSelectionWithDisposableSession(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": sessionCapabilitySnapshotFixture(t, "gpt-5.5", []string{"gpt-5.5"}, "reasoning_effort", "medium", []string{"medium", "xhigh"}),
	}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh"}

	err := harness.runner.Probe(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})

	if err != nil {
		t.Fatalf("expected supported selection to pass, got %v", err)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if len(invocations) != 7 {
		t.Fatalf("expected version, catalogue ensure and show, selection ensure and show, reasoning, close invocations, got %#v", invocations)
	}
	if !reflect.DeepEqual(invocations[0], []string{"--version"}) {
		t.Fatalf("expected version check first, got %#v", invocations[0])
	}
	sessionName := assertDisposableCatalogueEnsureInvocation(t, invocations[1], harness.gitRoot, "codex")
	assertDisposableShowInvocation(t, invocations[2], harness.gitRoot, "codex", sessionName)
	if selectedSession := assertDisposableEnsureInvocation(t, invocations[3], harness.gitRoot, "codex", "gpt-5.5"); selectedSession != sessionName {
		t.Fatalf("selection ensure used session %q, want catalogue session %q", selectedSession, sessionName)
	}
	assertDisposableShowInvocation(t, invocations[4], harness.gitRoot, "codex", sessionName)
	wantReasoning := []string{"--cwd", harness.gitRoot, "--format", "json", "--json-strict", "codex", "set", "reasoning_effort", "xhigh", "-s", sessionName}
	if !reflect.DeepEqual(invocations[5], wantReasoning) {
		t.Fatalf("unexpected reasoning invocation\nwant: %#v\ngot:  %#v", wantReasoning, invocations[5])
	}
	wantClose := []string{"--cwd", harness.gitRoot, "codex", "sessions", "close", sessionName}
	if !reflect.DeepEqual(invocations[6], wantClose) {
		t.Fatalf("unexpected disposable close invocation\nwant: %#v\ngot:  %#v", wantClose, invocations[6])
	}
	if containsCommandKey(invocations, "prompt") {
		t.Fatalf("selection preflight must not send prompts, got %#v", invocations)
	}
}

func TestACPXProbeSkipsEmptyReasoningEffort(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": sessionCapabilitySnapshotFixture(t, "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "", "", nil),
	}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol"}

	err := harness.runner.Probe(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})

	if err != nil {
		t.Fatalf("expected model-managed reasoning selection to pass, got %v", err)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if len(invocations) != 6 {
		t.Fatalf("expected version, catalogue ensure and show, selection ensure and show, close invocations, got %#v", invocations)
	}
	if !reflect.DeepEqual(invocations[0], []string{"--version"}) {
		t.Fatalf("expected version check first, got %#v", invocations[0])
	}
	sessionName := assertDisposableCatalogueEnsureInvocation(t, invocations[1], harness.gitRoot, "codex")
	assertDisposableShowInvocation(t, invocations[2], harness.gitRoot, "codex", sessionName)
	if selectedSession := assertDisposableEnsureInvocation(t, invocations[3], harness.gitRoot, "codex", "gpt-5.6-sol"); selectedSession != sessionName {
		t.Fatalf("selection ensure used session %q, want catalogue session %q", selectedSession, sessionName)
	}
	assertDisposableShowInvocation(t, invocations[4], harness.gitRoot, "codex", sessionName)
	wantClose := []string{"--cwd", harness.gitRoot, "codex", "sessions", "close", sessionName}
	if !reflect.DeepEqual(invocations[5], wantClose) {
		t.Fatalf("unexpected disposable close invocation\nwant: %#v\ngot:  %#v", wantClose, invocations[5])
	}
	if containsCommandKey(invocations, "set reasoning_effort") {
		t.Fatalf("expected no reasoning set call for empty effort, got %#v", invocations)
	}
}

func TestProfileProofAppliesExactReasoningAndClosesDisposableSession(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": sessionCapabilitySnapshotFixture(t, "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", []string{"medium", "high"}),
	}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"}

	err := harness.runner.Probe(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})

	if err != nil {
		t.Fatalf("profile proof failed: %v", err)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if containsCommandKey(invocations, "prompt") {
		t.Fatalf("profile proof must not send prompts, got %#v", invocations)
	}
	if len(invocations) != 7 {
		t.Fatalf("expected version, catalogue ensure and show, selection ensure and show, reasoning, close invocations, got %#v", invocations)
	}
	sessionName := assertDisposableCatalogueEnsureInvocation(t, invocations[1], harness.gitRoot, "codex")
	assertDisposableShowInvocation(t, invocations[2], harness.gitRoot, "codex", sessionName)
	if selectedSession := assertDisposableEnsureInvocation(t, invocations[3], harness.gitRoot, "codex", "gpt-5.6-sol"); selectedSession != sessionName {
		t.Fatalf("selection ensure used session %q, want catalogue session %q", selectedSession, sessionName)
	}
	assertDisposableShowInvocation(t, invocations[4], harness.gitRoot, "codex", sessionName)
	wantReasoning := []string{"--cwd", harness.gitRoot, "--format", "json", "--json-strict", "codex", "set", "reasoning_effort", "high", "-s", sessionName}
	if !reflect.DeepEqual(invocations[5], wantReasoning) {
		t.Fatalf("profile proof did not apply exact reasoning\nwant: %#v\ngot:  %#v", wantReasoning, invocations[5])
	}
	wantClose := []string{"--cwd", harness.gitRoot, "codex", "sessions", "close", sessionName}
	if !reflect.DeepEqual(invocations[6], wantClose) {
		t.Fatalf("profile proof did not close disposable session\nwant: %#v\ngot:  %#v", wantClose, invocations[6])
	}
}

func TestProfileProofUsesSessionSnapshotWhenACPXModelProjectionOmitsOptions(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"set model value=gpt-5.5":          `{"action":"model_set","modelId":"gpt-5.5","resumed":false}`,
		"sessions show":                    sessionCapabilitySnapshotFixture(t, "gpt-5.5", []string{"gpt-5.6-sol", "gpt-5.5"}, "reasoning_effort", "xhigh", []string{"low", "medium", "high", "xhigh"}),
		"set reasoning_effort value=xhigh": selectionStateFixture(t, "reasoning_effort", "xhigh", "gpt-5.5", []string{"gpt-5.6-sol", "gpt-5.5"}, "reasoning_effort", "xhigh", []string{"low", "medium", "high", "xhigh"}),
	}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh"}

	err := harness.runner.Probe(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})

	if err != nil {
		t.Fatalf("expected public config-option response to prove selection: %v", err)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if containsCommandKey(invocations, "set model") {
		t.Fatalf("profile proof used ACPX model projection without config options: %#v", invocations)
	}
	if !containsCommandKey(invocations, "sessions show") {
		t.Fatalf("profile proof did not inspect the public Agent Session snapshot: %#v", invocations)
	}
	if !containsCommandKey(invocations, "set reasoning_effort") {
		t.Fatalf("profile proof did not acquire public config-option evidence: %#v", invocations)
	}
	assertLastInvocationClosesDisposable(t, harness)
}

func TestAcquireSelectionCapabilitiesObservesAfterACPXModelProjectionOmitsOptions(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	models := []string{"future-model", "future-model[high]"}
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"set model value=future-model[high]": `{"action":"model_set","modelId":"future-model[high]","resumed":false}`,
		"sessions show":                      sessionCapabilitySnapshotFixture(t, "future-model[high]", models, "", "", nil),
	}))

	capabilities, err := harness.runner.AcquireSelectionCapabilities(context.Background(), CapabilityAcquisitionRequest{
		Runtime:  RuntimeSpec{ID: "codex", Protocol: ProtocolACP},
		Session:  SessionRef{Name: "roundfix-live", WorkDir: harness.gitRoot},
		ConfigID: "model",
		Value:    "future-model[high]",
		Adapter:  AdapterEvidence{Command: "adapter"},
	})

	if err != nil {
		t.Fatalf("acquire capabilities after model update: %v", err)
	}
	if capabilities.CurrentModel != "future-model[high]" {
		t.Fatalf("current model = %q, want future-model[high]", capabilities.CurrentModel)
	}
	if got, want := selectionCallKeys(readJSONInvocations(t, harness.invocationsPath)), []string{"set model value=future-model[high]"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("selection calls = %#v, want %#v", got, want)
	}
	if !containsCommandKey(readJSONInvocations(t, harness.invocationsPath), "sessions show") {
		t.Fatal("model update did not refresh capabilities from the public Agent Session snapshot")
	}
}

func TestProfileProofClosesDisposableSessionOnSelectionFailure(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": sessionCapabilitySnapshotFixture(t, "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", []string{"medium", "high"}),
	}))
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"set reasoning_effort": 2}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"set reasoning_effort": "reasoning rejected\n"}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"}

	err := harness.runner.Probe(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})

	if err == nil {
		t.Fatal("expected profile proof selection failure")
	}
	var selectionErr *SelectionRejectedError
	if !errors.As(err, &selectionErr) {
		t.Fatalf("expected SelectionRejectedError, got %T %v", err, err)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if len(invocations) < 5 {
		t.Fatalf("expected selection invocations through index 4, got %#v", invocations)
	}
	sessionName := assertDisposableCatalogueEnsureInvocation(t, invocations[1], harness.gitRoot, "codex")
	assertDisposableShowInvocation(t, invocations[2], harness.gitRoot, "codex", sessionName)
	if selectedSession := assertDisposableEnsureInvocation(t, invocations[3], harness.gitRoot, "codex", "gpt-5.6-sol"); selectedSession != sessionName {
		t.Fatalf("selection ensure used session %q, want catalogue session %q", selectedSession, sessionName)
	}
	assertDisposableShowInvocation(t, invocations[4], harness.gitRoot, "codex", sessionName)
	wantClose := []string{"--cwd", harness.gitRoot, "codex", "sessions", "close", sessionName}
	if !reflect.DeepEqual(invocations[len(invocations)-1], wantClose) {
		t.Fatalf("profile proof did not close disposable session after failure\nwant: %#v\ngot:  %#v", wantClose, invocations[len(invocations)-1])
	}
	if containsCommandKey(invocations, "prompt") {
		t.Fatalf("profile proof must not send prompts after failure, got %#v", invocations)
	}
}

func TestApplySessionSelection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		runtime        RuntimeSpec
		initial        string
		stdoutByCall   map[string]string
		exitBy         map[string]int
		wantCalls      []string
		classification string
	}{
		{
			name:    "independent controls apply model before reasoning",
			runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"},
			initial: selectionStateFixture(t, "model", "gpt-5.5", "gpt-5.5", []string{"gpt-5.6-sol", "gpt-5.5"}, "reasoning_effort", "medium", []string{"medium", "high", "xhigh"}),
			stdoutByCall: map[string]string{
				"set model value=gpt-5.6-sol":     selectionStateFixture(t, "model", "gpt-5.6-sol", "gpt-5.6-sol", []string{"gpt-5.6-sol", "gpt-5.5"}, "reasoning_effort", "medium", []string{"medium", "high", "xhigh"}),
				"set reasoning_effort value=high": selectionStateFixture(t, "reasoning_effort", "high", "gpt-5.6-sol", []string{"gpt-5.6-sol", "gpt-5.5"}, "reasoning_effort", "high", []string{"medium", "high", "xhigh"}),
			},
			wantCalls: []string{"set model value=gpt-5.6-sol", "set reasoning_effort value=high"},
		},
		{
			name:    "model variant applies only advertised variant",
			runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "future-model", ReasoningEffort: "high"},
			initial: selectionStateFixture(t, "model", "future-model", "future-model", []string{"future-model", "future-model[high]", "future-model[xhigh]"}, "", "", nil),
			stdoutByCall: map[string]string{
				"set model value=future-model[high]": selectionStateFixture(t, "model", "future-model[high]", "future-model[high]", []string{"future-model", "future-model[high]", "future-model[xhigh]"}, "", "", nil),
			},
			wantCalls: []string{"set model value=future-model[high]"},
		},
		{
			name:    "zero exit mismatched state fails",
			runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"},
			initial: selectionStateFixture(t, "model", "gpt-5.6-sol", "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", []string{"medium", "high"}),
			stdoutByCall: map[string]string{
				"set reasoning_effort value=high": selectionStateFixture(t, "reasoning_effort", "high", "gpt-5.5", []string{"gpt-5.6-sol", "gpt-5.5"}, "reasoning_effort", "high", []string{"medium", "high"}),
			},
			wantCalls:      []string{"set reasoning_effort value=high"},
			classification: EffectiveSelectionMismatch,
		},
		{
			name:           "rejected effort never retries model managed",
			runtime:        RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"},
			initial:        selectionStateFixture(t, "model", "gpt-5.6-sol", "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", []string{"medium", "high"}),
			exitBy:         map[string]int{"set reasoning_effort": 2},
			wantCalls:      []string{"set reasoning_effort value=high"},
			classification: SelectionRejected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			harness := newFakeACPXHarness(t)
			harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, tt.stdoutByCall))
			if tt.exitBy != nil {
				harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, tt.exitBy))
			}
			capabilities, err := ParseSessionConfigOptions([]byte(tt.initial), AdapterEvidence{Command: "adapter"}, SelectionRetention{})
			if err != nil {
				t.Fatalf("parse initial capabilities: %v", err)
			}
			proof, err := harness.runner.ApplySessionSelection(context.Background(), SessionSelectionRequest{
				Runtime:      tt.runtime,
				Session:      SessionRef{Name: "roundfix-live", WorkDir: harness.gitRoot},
				Capabilities: capabilities,
			})
			if tt.classification == "" {
				if err != nil {
					t.Fatalf("apply selection: %v", err)
				}
				if proof.Status != SelectionProofStatusProven || proof.Model != tt.runtime.Model || proof.ReasoningEffort != tt.runtime.ReasoningEffort {
					t.Fatalf("unexpected proof: %#v", proof)
				}
			} else {
				if err == nil {
					t.Fatal("expected selection application failure")
				}
				classified, ok := err.(interface{ Classification() string })
				if !ok || classified.Classification() != tt.classification {
					t.Fatalf("classification = %T %v, want %q", err, err, tt.classification)
				}
			}
			if got := selectionCallKeys(readJSONInvocations(t, harness.invocationsPath)); !reflect.DeepEqual(got, tt.wantCalls) {
				t.Fatalf("selection calls = %#v, want %#v", got, tt.wantCalls)
			}
		})
	}
}

func TestProveExactSelectionOfficialFixturesNoPrompt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		model  string
		effort string
	}{
		{name: "Sol high", model: "gpt-5.6-sol", effort: "high"},
		{name: "GPT 5.5 xhigh", model: "gpt-5.5", effort: "xhigh"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			harness := newFakeACPXHarness(t)
			models := []string{"gpt-5.6-sol", "gpt-5.6-terra", "gpt-5.6-luna", "gpt-5.5"}
			efforts := []string{"low", "medium", "high", "xhigh", "max", "ultra"}
			harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
				"sessions show":                           sessionCapabilitySnapshotFixture(t, tt.model, models, "reasoning_effort", "medium", efforts),
				"set model value=" + tt.model:             selectionStateFixture(t, "model", tt.model, tt.model, models, "reasoning_effort", "medium", efforts),
				"set reasoning_effort value=" + tt.effort: selectionStateFixture(t, "reasoning_effort", tt.effort, tt.model, models, "reasoning_effort", tt.effort, efforts),
			}))

			proof, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{
				Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: tt.model, ReasoningEffort: tt.effort},
				WorkDir: harness.gitRoot,
			})
			if err != nil {
				t.Fatalf("prove exact selection: %v", err)
			}
			if proof.Status != SelectionProofStatusProven || proof.Model != tt.model || proof.ReasoningEffort != tt.effort {
				t.Fatalf("unexpected proof: %#v", proof)
			}
			invocations := readJSONInvocations(t, harness.invocationsPath)
			if containsCommandKey(invocations, "prompt") {
				t.Fatalf("proof sent an Agent prompt: %#v", invocations)
			}
			if got, want := exactProofCallKeys(invocations), []string{"sessions ensure model=", "sessions show", "sessions ensure model=" + tt.model, "sessions show", "set reasoning_effort value=" + tt.effort, "sessions close"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("proof calls = %#v, want %#v", got, want)
			}
		})
	}
}

func TestProveExactSelectionModelVariant(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	models := []string{"future-model", "future-model[high]", "future-model[xhigh]"}
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show":                      sessionCapabilitySnapshotFixture(t, "future-model", models, "", "", nil),
		"set model value=future-model":       selectionStateFixture(t, "model", "future-model", "future-model", models, "", "", nil),
		"set model value=future-model[high]": selectionStateFixture(t, "model", "future-model[high]", "future-model[high]", models, "", "", nil),
	}))

	proof, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{
		Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "future-model", ReasoningEffort: "high"},
		WorkDir: harness.gitRoot,
	})
	if err != nil {
		t.Fatalf("prove model variant: %v", err)
	}
	if proof.Assignment.Encoding != SelectionEncodingModelVariant || proof.Assignment.AdapterModel != "future-model[high]" {
		t.Fatalf("unexpected variant proof: %#v", proof)
	}
	if containsCommandKey(readJSONInvocations(t, harness.invocationsPath), "set reasoning_effort") {
		t.Fatal("model variant proof must not send an independent reasoning operation")
	}
}

func TestProveExactSelectionEffectiveMismatchCleanup(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	efforts := []string{"medium", "high"}
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show":                   sessionCapabilitySnapshotFixture(t, "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", efforts),
		"set model value=gpt-5.6-sol":     selectionStateFixture(t, "model", "gpt-5.6-sol", "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", efforts),
		"set reasoning_effort value=high": selectionStateFixture(t, "reasoning_effort", "high", "gpt-5.5", []string{"gpt-5.6-sol", "gpt-5.5"}, "reasoning_effort", "high", efforts),
	}))

	_, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"}, WorkDir: harness.gitRoot})
	var mismatch *EffectiveSelectionError
	if !errors.As(err, &mismatch) {
		t.Fatalf("error = %T %v, want *EffectiveSelectionError", err, err)
	}
	assertLastInvocationClosesDisposable(t, harness)
}

func TestProveExactSelectionMalformedEvidenceCleanup(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{"sessions show": `{`}))

	_, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"}, WorkDir: harness.gitRoot})
	var evidence *CapabilityEvidenceError
	if !errors.As(err, &evidence) {
		t.Fatalf("error = %T %v, want *CapabilityEvidenceError", err, err)
	}
	assertLastInvocationClosesDisposable(t, harness)
}

func TestProveExactSelectionCleanupJoinedFailure(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show":               sessionCapabilitySnapshotFixture(t, "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", []string{"medium", "high"}),
		"set model value=gpt-5.6-sol": selectionStateFixture(t, "model", "gpt-5.6-sol", "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", []string{"medium", "high"}),
	}))
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"set reasoning_effort": 2, "sessions close": 2}))

	_, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"}, WorkDir: harness.gitRoot})
	var rejected *SelectionRejectedError
	if !errors.As(err, &rejected) {
		t.Fatalf("joined error lacks selection rejection: %T %v", err, err)
	}
	var cleanup *AgentSessionCleanupError
	if !errors.As(err, &cleanup) {
		t.Fatalf("joined error lacks cleanup failure: %T %v", err, err)
	}
}

// wantCatalogueEnsureRejectionDiagnosis is the diagnosis that a disposable
// Agent Session whose override-free catalogue ensure failed records without a
// close error. It is shared by every missing-session exit shape.
const wantCatalogueEnsureRejectionDiagnosis = `apply Agent Selection "codex"/"gpt-5.6-sol"/"high" during ensure disposable Agent Session without model override: adapter rejected selection: acpx infrastructure error after exit code 2: acpx command failed; recovery: update the ACP Runtime or adapter and retry the exact Agent Selection`

func TestMissingSessionIsRecognisedFromBothExitShapes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		exitCode int
		stderr   string
	}{
		{
			name:     "missing session exit four",
			exitCode: 4,
		},
		{
			name:     "installed missing session exit one",
			exitCode: 1,
			stderr:   "No named session: roundfix-task-08\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			harness := newFakeACPXHarness(t)
			harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{
				"sessions ensure": 2,
				"sessions close":  test.exitCode,
			}))
			if test.stderr != "" {
				harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{
					"sessions close": test.stderr,
				}))
			}
			var warnings []string
			harness.runner.warnf = func(format string, args ...any) {
				warnings = append(warnings, fmt.Sprintf(format, args...))
			}

			_, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{
				Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"},
				WorkDir: harness.gitRoot,
			})
			if err == nil {
				t.Fatal("expected selection failure")
			}
			if got := err.Error(); got != wantCatalogueEnsureRejectionDiagnosis {
				t.Fatalf("selection diagnosis changed or gained a close error\nwant: %q\ngot:  %q", wantCatalogueEnsureRejectionDiagnosis, got)
			}
			var cleanupErr *AgentSessionCleanupError
			if errors.As(err, &cleanupErr) {
				t.Fatalf("missing disposable session close was appended to the diagnosis: %v", err)
			}
			if len(warnings) != 1 || !strings.Contains(warnings[0], "close disposable Agent Session") || !strings.Contains(warnings[0], acpxExitReasonMissingSession) {
				t.Fatalf("missing disposable session close was not recorded: %#v", warnings)
			}
		})
	}
}

func TestUnrelatedExitOneKeepsItsClassification(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions close": 1}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"sessions close": "close rejected\n"}))

	err := harness.runner.CloseSession(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, SessionRef{Name: "roundfix-task-08", WorkDir: harness.gitRoot})
	var infrastructureErr *InfrastructureError
	if !errors.As(err, &infrastructureErr) {
		t.Fatalf("error = %T %v, want *InfrastructureError", err, err)
	}
	const want = "close acpx Agent Session \"roundfix-task-08\": acpx infrastructure error after exit code 1: acpx command failed\n--- acpx stderr tail ---\nclose rejected"
	if got := err.Error(); got != want {
		t.Fatalf("unrelated exit 1 classification changed\nwant: %q\ngot:  %q", want, got)
	}
}

func TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show":               sessionCapabilitySnapshotFixture(t, "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", []string{"medium", "high"}),
		"set model value=gpt-5.6-sol": selectionStateFixture(t, "model", "gpt-5.6-sol", "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", []string{"medium", "high"}),
	}))
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{
		"set reasoning_effort": 2,
		"sessions close":       2,
	}))

	_, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{
		Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"},
		WorkDir: harness.gitRoot,
	})
	if err == nil {
		t.Fatal("expected selection and cleanup failure")
	}
	const wantDiagnosis = `apply Agent Selection "codex"/"gpt-5.6-sol"/"high" during set reasoning_effort: adapter rejected selection: acquire Agent Selection capabilities through acpx: command failed; recovery: update the ACP Runtime or adapter and retry the exact Agent Selection`
	if got := err.Error(); !strings.HasPrefix(got, wantDiagnosis+"\nclose disposable Agent Session") {
		t.Fatalf("selection diagnosis or appended close error changed\nwant prefix: %q\ngot:         %q", wantDiagnosis+"\nclose disposable Agent Session", got)
	}
	var cleanupErr *AgentSessionCleanupError
	if !errors.As(err, &cleanupErr) {
		t.Fatalf("open disposable session close was not recorded in the error chain: %T %v", err, err)
	}
	if got, want := selectionCallKeys(readJSONInvocations(t, harness.invocationsPath)), []string{"sessions ensure model=", "sessions ensure model=gpt-5.6-sol", "set reasoning_effort value=high", "sessions close"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("proof did not open and then attempt to close the disposable session\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestProveExactSelectionCancelCleanup(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	started := filepath.Join(harness.gitRoot, "reasoning-started")
	harness.setEnv(fakeACPXStarted, started)
	harness.setEnv(fakeACPXBlockCmd, "set reasoning_effort")
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show":               sessionCapabilitySnapshotFixture(t, "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", []string{"medium", "high"}),
		"set model value=gpt-5.6-sol": selectionStateFixture(t, "model", "gpt-5.6-sol", "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", []string{"medium", "high"}),
	}))
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := harness.runner.ProveExactSelection(ctx, ProbeRequest{Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"}, WorkDir: harness.gitRoot})
		result <- err
	}()
	waitForFile(t, started)
	cancel()

	if err := receiveError(t, result); !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %T %v, want context.Canceled", err, err)
	}
	assertLastInvocationClosesDisposable(t, harness)
}

func TestProveExactSelectionTimeoutCleanup(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	started := filepath.Join(harness.gitRoot, "session-state-started")
	harness.setEnv(fakeACPXStarted, started)
	harness.setEnv(fakeACPXBlockCmd, "sessions show")
	ctx := newTestDeadlineContext()
	defer ctx.expire()
	result := make(chan error, 1)
	go func() {
		_, err := harness.runner.ProveExactSelection(ctx, ProbeRequest{Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"}, WorkDir: harness.gitRoot})
		result <- err
	}()
	waitForFile(t, started)
	ctx.expire()
	err := receiveError(t, result)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %T %v, want context.DeadlineExceeded", err, err)
	}
	assertLastInvocationClosesDisposable(t, harness)
}

func TestApplySessionSelectionDisposableAndLiveOrder(t *testing.T) {
	t.Parallel()

	newHarness := func(t *testing.T) *fakeACPXHarness {
		harness := newFakeACPXHarness(t)
		efforts := []string{"medium", "high"}
		harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
			"sessions show":                   sessionCapabilitySnapshotFixture(t, "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", efforts),
			"set model value=gpt-5.6-sol":     selectionStateFixture(t, "model", "gpt-5.6-sol", "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", efforts),
			"set reasoning_effort value=high": selectionStateFixture(t, "reasoning_effort", "high", "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "high", efforts),
		}))
		return harness
	}
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"}
	disposable := newHarness(t)
	if _, err := disposable.runner.ProveExactSelection(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: disposable.gitRoot}); err != nil {
		t.Fatalf("prove disposable selection: %v", err)
	}
	live := newHarness(t)
	if err := live.runner.PrepareSession(context.Background(), ExecuteRequest{Runtime: runtime, Session: SessionRef{Name: "roundfix-live"}, GitRoot: live.gitRoot}, runevent.Discard); err != nil {
		t.Fatalf("prepare live selection: %v", err)
	}

	disposableInvocations := readJSONInvocations(t, disposable.invocationsPath)
	if len(disposableInvocations) < 4 {
		t.Fatalf("expected disposable selection invocations through index 3, got %#v", disposableInvocations)
	}
	disposableSession := assertDisposableCatalogueEnsureInvocation(t, disposableInvocations[0], disposable.gitRoot, "codex")
	assertDisposableShowInvocation(t, disposableInvocations[1], disposable.gitRoot, "codex", disposableSession)
	if selectedSession := assertDisposableEnsureInvocation(t, disposableInvocations[2], disposable.gitRoot, "codex", "gpt-5.6-sol"); selectedSession != disposableSession {
		t.Fatalf("selection ensure used session %q, want catalogue session %q", selectedSession, disposableSession)
	}
	assertDisposableShowInvocation(t, disposableInvocations[3], disposable.gitRoot, "codex", disposableSession)
	wantDisposable := []string{"sessions ensure model=", "sessions ensure model=gpt-5.6-sol", "set reasoning_effort value=high", "sessions close"}
	if got := selectionCallKeys(disposableInvocations); !reflect.DeepEqual(got, wantDisposable) {
		t.Fatalf("disposable calls = %#v, want %#v", got, wantDisposable)
	}
	want := []string{"sessions ensure model=gpt-5.6-sol", "set reasoning_effort value=high"}
	if got := selectionCallKeys(readJSONInvocations(t, live.invocationsPath)); !reflect.DeepEqual(got, want) {
		t.Fatalf("live calls = %#v, want %#v", got, want)
	}
}

func TestACPXProbeSelectionSetupUsesBoundedContext(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": sessionCapabilitySnapshotFixture(t, "gpt-5.5", []string{"gpt-5.5"}, "reasoning_effort", "medium", []string{"medium", "xhigh"}),
	}))
	harness.setEnv(codexPathEnv, "/configured/clean/codex")
	probe := &deadlineRecordingCodexProbe{}
	harness.runner.codexSpawn = codexSpawnDependencies{
		goos:       "darwin",
		lookPath:   fakeCodexLookPath("/path/clean/codex", nil),
		quarantine: probe,
		acceptance: probe,
	}
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh"}

	err := harness.runner.Probe(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})

	if err != nil {
		t.Fatalf("expected selection preflight to pass, got %v", err)
	}
	if len(probe.deadlineSeen) == 0 {
		t.Fatal("expected codex setup inspection to receive a context")
	}
	for index, sawDeadline := range probe.deadlineSeen {
		if !sawDeadline {
			t.Fatalf("expected setup context %d to have a deadline", index)
		}
	}
}

func TestACPXProbeSelectionRejectionClosesDisposableSession(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": sessionCapabilitySnapshotFixture(t, "gpt-5.5", []string{"gpt-5.5"}, "reasoning_effort", "medium", []string{"medium", "extreme"}),
	}))
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"set reasoning_effort": 2}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"set reasoning_effort": "reasoning value rejected\n"}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "extreme"}

	err := harness.runner.Probe(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})

	if err == nil {
		t.Fatal("expected rejected reasoning to fail")
	}
	var selectionErr *SelectionRejectedError
	if !errors.As(err, &selectionErr) {
		t.Fatalf("expected SelectionRejectedError, got %T %v", err, err)
	}
	if selectionErr.Assignment.Runtime != "codex" || selectionErr.Assignment.Model != "gpt-5.5" || selectionErr.Assignment.ReasoningEffort != "extreme" {
		t.Fatalf("unexpected selection tuple: %#v", selectionErr)
	}
	var infraErr *InfrastructureError
	if !errors.As(err, &infraErr) {
		t.Fatalf("expected original InfrastructureError in chain, got %T %v", err, err)
	}
	if infraErr.Stderr != "reasoning value rejected\n" {
		t.Fatalf("expected adapter stderr preserved, got %q", infraErr.Stderr)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if len(invocations) != 7 {
		t.Fatalf("expected version, catalogue ensure and show, selection ensure and show, reasoning, close invocations, got %#v", invocations)
	}
	sessionName := assertDisposableCatalogueEnsureInvocation(t, invocations[1], harness.gitRoot, "codex")
	assertDisposableShowInvocation(t, invocations[2], harness.gitRoot, "codex", sessionName)
	if selectedSession := assertDisposableEnsureInvocation(t, invocations[3], harness.gitRoot, "codex", "gpt-5.5"); selectedSession != sessionName {
		t.Fatalf("selection ensure used session %q, want catalogue session %q", selectedSession, sessionName)
	}
	assertDisposableShowInvocation(t, invocations[4], harness.gitRoot, "codex", sessionName)
	wantClose := []string{"--cwd", harness.gitRoot, "codex", "sessions", "close", sessionName}
	if !reflect.DeepEqual(invocations[6], wantClose) {
		t.Fatalf("expected cleanup close after rejection\nwant: %#v\ngot:  %#v", wantClose, invocations[6])
	}
}

func TestACPXProbeModelRejectionSkipsReasoningAndClosesDisposableSession(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions ensure": 2}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"sessions ensure": "missing model metadata\n"}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "xhigh"}

	err := harness.runner.Probe(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})

	if err == nil {
		t.Fatal("expected unavailable model to fail")
	}
	var selectionErr *SelectionRejectedError
	if !errors.As(err, &selectionErr) {
		t.Fatalf("expected SelectionRejectedError, got %T %v", err, err)
	}
	if !strings.Contains(err.Error(), "gpt-5.6-sol") || !strings.Contains(err.Error(), "xhigh") {
		t.Fatalf("expected selection tuple in error, got %v", err)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if len(invocations) != 3 {
		t.Fatalf("expected version, ensure, close only, got %#v", invocations)
	}
	if containsCommandKey(invocations, "set reasoning_effort") || containsCommandKey(invocations, "prompt") {
		t.Fatalf("expected no fallback reasoning or prompt after model rejection, got %#v", invocations)
	}
}

func TestACPXProbeCleanupFailureJoinsSelectionError(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": sessionCapabilitySnapshotFixture(t, "gpt-5.5", []string{"gpt-5.5"}, "reasoning_effort", "medium", []string{"medium", "extreme"}),
	}))
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{
		"set reasoning_effort": 2,
		"sessions close":       2,
	}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{
		"set reasoning_effort": "reasoning value rejected\n",
		"sessions close":       "close rejected\n",
	}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "extreme"}

	err := harness.runner.Probe(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})

	if err == nil {
		t.Fatal("expected selection and cleanup failure")
	}
	var selectionErr *SelectionRejectedError
	if !errors.As(err, &selectionErr) {
		t.Fatalf("expected selection error in joined chain, got %T %v", err, err)
	}
	var cleanupErr *AgentSessionCleanupError
	if !errors.As(err, &cleanupErr) {
		t.Fatalf("expected cleanup error in joined chain, got %T %v", err, err)
	}
	if selectionErr.Classification() != SelectionRejected || cleanupErr.Classification() != SessionCleanupFailed {
		t.Fatalf("unexpected joined classifications: selection=%q cleanup=%q", selectionErr.Classification(), cleanupErr.Classification())
	}
}

func TestACPXProbeCancellationStillClosesDisposableSession(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
	startedPath := filepath.Join(harness.gitRoot, "ensure-started")
	harness.setEnv(fakeACPXStarted, startedPath)
	harness.setEnv(fakeACPXBlockCmd, "sessions ensure")
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh"}
	ctx, cancel := context.WithCancel(context.Background())
	resultCh := make(chan error, 1)
	go func() {
		resultCh <- harness.runner.Probe(ctx, ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})
	}()
	waitForFile(t, startedPath)
	cancel()

	err := receiveError(t, resultCh)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled in error chain, got %T %v", err, err)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if len(invocations) != 3 {
		t.Fatalf("expected version, blocked ensure, cleanup close, got %#v", invocations)
	}
	sessionName := disposableSessionFromEnsure(t, invocations[1])
	wantClose := []string{"--cwd", harness.gitRoot, "codex", "sessions", "close", sessionName}
	if !reflect.DeepEqual(invocations[2], wantClose) {
		t.Fatalf("expected cleanup close after cancellation\nwant: %#v\ngot:  %#v", wantClose, invocations[2])
	}
}

func TestACPXProbeFallbackUsesCandidateAndEffortOrder(t *testing.T) {
	// Sequential: ProbeFallback resolves its working directory from process state.
	harness := newFakeACPXHarness(t)
	t.Chdir(harness.gitRoot)
	harness.setEnv(fakeACPXExitByCall, mustJSONForTest(t, map[string]int{
		"sessions ensure model=gpt-newest": 2,
		"set reasoning_effort value=xhigh": 2,
	}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-failed", ReasoningEffort: "xhigh"}

	selection, ok, err := harness.runner.ProbeFallback(context.Background(), runtime, FallbackCandidateSet{
		Models:  []string{"gpt-failed", "gpt-newest", "gpt-working", "gpt-older"},
		Efforts: []string{"xhigh", "high", "medium"},
	})

	if err != nil {
		t.Fatalf("expected fallback probe to succeed, got %v", err)
	}
	if !ok {
		t.Fatal("expected a proven Fallback Selection")
	}
	if want := (FallbackSelection{Model: "gpt-working", ReasoningEffort: "high"}); selection != want {
		t.Fatalf("unexpected Fallback Selection\nwant: %#v\ngot:  %#v", want, selection)
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	assertFallbackCandidateSessions(t, invocations, []string{"gpt-newest", "gpt-working"})
	if containsInvocationValue(invocations, "--model", "gpt-failed") {
		t.Fatalf("failed Agent Model must not be probed: %#v", invocations)
	}
	if containsInvocationValue(invocations, "--model", "gpt-older") {
		t.Fatalf("candidate walking must stop after the first proven model: %#v", invocations)
	}
	if containsCommandValue(invocations, "set reasoning_effort", "medium") {
		t.Fatalf("effort probing must stop after the highest accepted value: %#v", invocations)
	}
}

func TestACPXProbeFallbackClassifiesModelManagedReasoning(t *testing.T) {
	// Sequential: ProbeFallback resolves its working directory from process state.
	harness := newFakeACPXHarness(t)
	t.Chdir(harness.gitRoot)
	harness.setEnv(fakeACPXExitByCall, mustJSONForTest(t, map[string]int{
		"set reasoning_effort value=xhigh": 2,
		"set reasoning_effort value=high":  2,
	}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-failed", ReasoningEffort: "xhigh"}

	selection, ok, err := harness.runner.ProbeFallback(context.Background(), runtime, FallbackCandidateSet{
		Models:  []string{"gpt-managed"},
		Efforts: []string{"xhigh", "high"},
	})

	if err != nil {
		t.Fatalf("expected model-managed fallback probe to succeed, got %v", err)
	}
	if !ok {
		t.Fatal("expected model assignment to prove a Fallback Selection")
	}
	if want := (FallbackSelection{Model: "gpt-managed"}); selection != want {
		t.Fatalf("unexpected model-managed Fallback Selection\nwant: %#v\ngot:  %#v", want, selection)
	}
	assertFallbackCandidateSessions(t, readJSONInvocations(t, harness.invocationsPath), []string{"gpt-managed"})
}

func TestACPXProbeFallbackReportsNoFallbackWhenModelsAreRejected(t *testing.T) {
	// Sequential: ProbeFallback resolves its working directory from process state.
	harness := newFakeACPXHarness(t)
	t.Chdir(harness.gitRoot)
	harness.setEnv(fakeACPXExitByCall, mustJSONForTest(t, map[string]int{
		"sessions ensure model=gpt-newest": 2,
		"sessions ensure model=gpt-older":  2,
	}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-failed", ReasoningEffort: "xhigh"}

	selection, ok, err := harness.runner.ProbeFallback(context.Background(), runtime, FallbackCandidateSet{
		Models:  []string{"gpt-failed", "gpt-newest", "gpt-older"},
		Efforts: []string{"xhigh", "high"},
	})

	if err != nil {
		t.Fatalf("model rejection is not a probe infrastructure error: %v", err)
	}
	if ok || selection != (FallbackSelection{}) {
		t.Fatalf("expected no fallback, got ok=%t selection=%#v", ok, selection)
	}
	assertFallbackCandidateSessions(t, readJSONInvocations(t, harness.invocationsPath), []string{"gpt-newest", "gpt-older"})
}

func TestACPXProbeFallbackReportsCleanupInfrastructureError(t *testing.T) {
	// Sequential: ProbeFallback resolves its working directory from process state.
	harness := newFakeACPXHarness(t)
	t.Chdir(harness.gitRoot)
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions close": 2}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-failed", ReasoningEffort: "xhigh"}

	selection, ok, err := harness.runner.ProbeFallback(context.Background(), runtime, FallbackCandidateSet{
		Models:  []string{"gpt-working"},
		Efforts: []string{"high"},
	})

	if err == nil {
		t.Fatal("expected disposable session cleanup failure")
	}
	var cleanupErr *AgentSessionCleanupError
	if !errors.As(err, &cleanupErr) {
		t.Fatalf("expected AgentSessionCleanupError, got %T %v", err, err)
	}
	if ok || selection != (FallbackSelection{}) {
		t.Fatalf("infrastructure failure must not return a fallback, got ok=%t selection=%#v", ok, selection)
	}
	assertFallbackCandidateSessions(t, readJSONInvocations(t, harness.invocationsPath), []string{"gpt-working"})
}

func TestACPXProbeFallbackCancellationStillClosesCandidateSession(t *testing.T) {
	// Sequential: ProbeFallback resolves its working directory from process state.
	harness := newFakeACPXHarness(t)
	t.Chdir(harness.gitRoot)
	startedPath := filepath.Join(harness.gitRoot, "fallback-ensure-started")
	harness.setEnv(fakeACPXStarted, startedPath)
	harness.setEnv(fakeACPXBlockCmd, "sessions ensure")
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-failed", ReasoningEffort: "xhigh"}
	ctx, cancel := context.WithCancel(context.Background())
	resultCh := make(chan error, 1)
	go func() {
		_, _, err := harness.runner.ProbeFallback(ctx, runtime, FallbackCandidateSet{Models: []string{"gpt-working"}})
		resultCh <- err
	}()
	waitForFile(t, startedPath)
	cancel()

	err := receiveError(t, resultCh)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled in error chain, got %T %v", err, err)
	}
	assertFallbackCandidateSessions(t, readJSONInvocations(t, harness.invocationsPath), []string{"gpt-working"})
}

func TestInfrastructureErrorErrorIncludesBoundedStderrTail(t *testing.T) {
	t.Parallel()

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

func TestModelNotAdvertisedPromptExitYieldsTypedError(t *testing.T) {
	t.Parallel()

	stderr := "Cannot apply --model gpt-5.6-sol: the ACP agent did not advertise that model.\nAvailable models: gpt-5.5, gpt-5.1\n"
	run := runFakeACPXPrompt(t, fakeACPXPrompt{
		runtime:  RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "xhigh"},
		stderr:   stderr,
		exitCode: 1,
	})

	var selectionErr *SelectionFailure
	if !errors.As(run.err, &selectionErr) {
		t.Fatalf("expected SelectionFailure, got %T %v", run.err, run.err)
	}
	var modelErr *ModelNotAdvertisedError
	if !errors.As(run.err, &modelErr) {
		t.Fatalf("expected ModelNotAdvertisedError in chain, got %T %v", run.err, run.err)
	}
	if modelErr.Runtime != "codex" || modelErr.Model != "gpt-5.6-sol" {
		t.Fatalf("unexpected rejection tuple: %#v", modelErr)
	}
	if !reflect.DeepEqual(modelErr.Advertised, []string{"gpt-5.5", "gpt-5.1"}) {
		t.Fatalf("unexpected advertised list: %#v", modelErr.Advertised)
	}
	for _, want := range []string{"update the ACP Runtime or adapter", "choose an advertised Agent Model", "one-Run --model override"} {
		if !strings.Contains(modelErr.Error(), want) {
			t.Fatalf("expected recovery guidance containing %q, got %q", want, modelErr.Error())
		}
	}
	if strings.Contains(modelErr.Error(), "reasoning_effort") {
		t.Fatalf("model-advertisement recovery must not suggest reasoning changes, got %q", modelErr.Error())
	}
}

func TestModelNotAdvertisedSelectionPreflightYieldsTypedErrorThroughWrapChain(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions ensure": 2}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{
		"sessions ensure": "adapter stderr prefix\nCannot apply --model gpt-5.6-sol: the ACP agent did not advertise that model.\nAvailable models: gpt-5.5, gpt-5.1, opus\nadapter stderr suffix\n",
	}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "xhigh"}

	err := harness.runner.Probe(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})

	var modelErr *ModelNotAdvertisedError
	if !errors.As(err, &modelErr) {
		t.Fatalf("expected ModelNotAdvertisedError through wrap chain, got %T %v", err, err)
	}
	if modelErr.Runtime != "codex" || modelErr.Model != "gpt-5.6-sol" {
		t.Fatalf("unexpected rejection tuple: %#v", modelErr)
	}
	if !reflect.DeepEqual(modelErr.Advertised, []string{"gpt-5.5", "gpt-5.1", "opus"}) {
		t.Fatalf("unexpected advertised list: %#v", modelErr.Advertised)
	}
	var infraErr *InfrastructureError
	if !errors.As(err, &infraErr) {
		t.Fatalf("expected original InfrastructureError in chain, got %T %v", err, err)
	}
}

func TestModelNotAdvertisedGarbageStderrKeepsInfrastructureError(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions ensure": 2}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"sessions ensure": "unparseable garbage\n"}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "xhigh"}

	err := harness.runner.Probe(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})

	var modelErr *ModelNotAdvertisedError
	if errors.As(err, &modelErr) {
		t.Fatalf("garbage stderr must not become ModelNotAdvertisedError: %#v", modelErr)
	}
	var infraErr *InfrastructureError
	if !errors.As(err, &infraErr) {
		t.Fatalf("expected InfrastructureError in chain, got %T %v", err, err)
	}
	want := "acpx infrastructure error after exit code 2: acpx command failed\n--- acpx stderr tail ---\nunparseable garbage"
	if got := infraErr.Error(); got != want {
		t.Fatalf("expected existing infrastructure error unchanged\nwant: %q\ngot:  %q", want, got)
	}
}

func TestACPXPromptArgsPlaceGlobalsBeforeAgentAndSubcommand(t *testing.T) {
	t.Parallel()

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

// The setup turn that raises a runtime_deferred queue owner must be
// structurally unable to perform Agent work. Spec 0089's QA gate measured the
// unrestricted prompt reading directories, reading Specs, and calling a tool
// while the Agent Session still held its default effort, which is the state
// the warm-up exists to leave behind. The restriction is scoped to that one
// invocation: measured against OpenCode 1.18.15 on 2026-08-09, the following
// work prompt on the same Agent Session called tools normally.
func TestACPXPromptArgsInertWarmupWithholdsTools(t *testing.T) {
	t.Parallel()

	contains := func(args []string, want string) bool {
		for _, arg := range args {
			if arg == want {
				return true
			}
		}
		return false
	}
	runtime := RuntimeSpec{ID: "opencode", Protocol: ProtocolACP}

	inert, err := acpxPromptArgs(ACPXPromptRequest{
		ExecuteRequest: ExecuteRequest{Runtime: runtime, GitRoot: "/repo"},
		Session:        "roundfix-run-1",
		Inert:          true,
	})
	if err != nil {
		t.Fatalf("acpx prompt args: %v", err)
	}
	want := []string{
		"--cwd", "/repo",
		"--format", "json",
		"--json-strict",
		"--deny-all",
		"--allowed-tools", "",
		"opencode", "prompt",
		"-s", "roundfix-run-1",
		"-f", "-",
	}
	if !reflect.DeepEqual(inert, want) {
		t.Fatalf("inert warm-up must withhold every tool and deny permissions\nwant: %#v\ngot:  %#v", want, inert)
	}

	work, err := acpxPromptArgs(ACPXPromptRequest{
		ExecuteRequest: ExecuteRequest{Runtime: runtime, GitRoot: "/repo"},
		Session:        "roundfix-run-1",
	})
	if err != nil {
		t.Fatalf("acpx prompt args: %v", err)
	}
	for _, restricted := range []string{"--deny-all", "--allowed-tools"} {
		if contains(work, restricted) {
			t.Fatalf("a work prompt must keep its tools; found %q in %#v", restricted, work)
		}
	}
	if !contains(work, "--approve-all") {
		t.Fatalf("a work prompt must approve its permission requests; got %#v", work)
	}
}

func TestACPXCancelSessionInvokesSessionCancel(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	invocationsPath := filepath.Join(dir, "invocations.jsonl")
	environment := environmentForTest(
		fakeACPXEnv+"=1",
		fakeACPXInvokes+"="+invocationsPath,
	)

	err := (&ACPXRunner{Command: os.Args[0], Environment: environment, codexSpawn: codexSpawnDependencies{goos: "linux"}}).CancelSession(context.Background(), RuntimeSpec{
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
	t.Parallel()

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

func TestParseRoundfixSessionNameRoundTripsGeneratedSessionRefs(t *testing.T) {
	t.Parallel()

	runID := "run_20260717T132519Z_5fc633528e1ff1c5"
	tests := []SessionRef{
		SessionRefForQA(runID, "/repo"),
		SessionRefForReview(runID, 0, "/repo"),
		SessionRefForReview(runID, 7, "/repo"),
	}

	for _, ref := range tests {
		t.Run(ref.Name, func(t *testing.T) {
			got, ok := ParseRoundfixSessionName(ref.Name)
			if !ok {
				t.Fatalf("ParseRoundfixSessionName(%q) did not recognize generated name", ref.Name)
			}
			want := RoundfixSession{Name: ref.Name, RunID: runID}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("unexpected session\nwant: %#v\ngot:  %#v", want, got)
			}
		})
	}
}

func TestACPXListRoundfixSessionsInvokesSessionsList(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdoutBy, mustJSONForTest(t, map[string]string{
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
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions close": 1}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"sessions close": "close rejected\n"}))

	err := harness.runner.CloseSession(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, SessionRef{Name: "roundfix-run-1", WorkDir: harness.gitRoot})

	if err == nil {
		t.Fatal("expected close failure")
	}
	if !strings.Contains(err.Error(), "close acpx Agent Session") || !strings.Contains(err.Error(), "close rejected") {
		t.Fatalf("expected close error context, got %v", err)
	}
}

func TestACPXRunPromptSendsPromptOnStdin(t *testing.T) {
	t.Parallel()

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

func TestACPXRunRequiresModelSelection(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	_, err := harness.runner.Run(context.Background(), ExecuteRequest{
		Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, ReasoningEffort: "xhigh"},
		RunID:   "run-acpx",
		Batch:   rounds.Batch{Number: 7},
		LogPath: filepath.Join(harness.gitRoot, "runs", "run-acpx", "agent", "batch-007.log"),
		Prompt:  "prompt",
		GitRoot: harness.gitRoot,
		Session: SessionRef{Name: "roundfix-run-1"},
	}, newCaptureSink(""))

	if err == nil {
		t.Fatal("expected missing model selection to fail")
	}
	if !strings.Contains(err.Error(), "model") {
		t.Fatalf("expected error containing model, got %q", err.Error())
	}
	if _, statErr := os.Stat(harness.invocationsPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("expected no acpx invocation before selection validation, got stat error %v", statErr)
	}
}

func TestValidateRuntimeSelectionSkipsReasoningConfigKeyForEmptyEffort(t *testing.T) {
	t.Parallel()

	if err := validateRuntimeSelection(RuntimeSpec{ID: "future-runtime", Protocol: ProtocolACP, Model: "future-model"}); err != nil {
		t.Fatalf("expected empty reasoning effort to skip config-key validation, got %v", err)
	}

	err := validateRuntimeSelection(RuntimeSpec{ID: "future-runtime", Protocol: ProtocolACP, Model: "future-model", ReasoningEffort: "high"})
	if err == nil {
		t.Fatal("expected non-empty reasoning effort to require a config key")
	}
	if !strings.Contains(err.Error(), "unsupported ACP Runtime") {
		t.Fatalf("expected unsupported runtime error, got %v", err)
	}
}

func TestACPXRunSkipsEmptyReasoningEffort(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": sessionCapabilitySnapshotFixture(t, "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "", "", nil),
	}))
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol"}

	_, err := harness.runner.Run(context.Background(), ExecuteRequest{
		Runtime:   runtime,
		RunID:     "run-acpx",
		Batch:     rounds.Batch{Number: 7},
		LogPath:   filepath.Join(harness.gitRoot, "runs", "run-acpx", "agent", "batch-007.log"),
		Prompt:    "prompt",
		GitRoot:   harness.gitRoot,
		StopGrace: 20 * time.Millisecond,
		Session:   SessionRef{Name: "roundfix-run-1"},
	}, newCaptureSink(""))

	if err != nil {
		t.Fatalf("expected model-managed reasoning run to pass, got %v", err)
	}
	want := [][]string{
		{"--cwd", harness.gitRoot, "--model", "gpt-5.6-sol", "codex", "sessions", "ensure", "--name", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "gpt-5.6-sol", "codex", "prompt", "-s", "roundfix-run-1", "-f", "-"},
	}
	if got := readJSONInvocations(t, harness.invocationsPath); !reflect.DeepEqual(got, exactSelectionInvocations(want)) {
		t.Fatalf("unexpected acpx invocations\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestACPXRunAppliesSelectionBeforePrompt(t *testing.T) {
	t.Parallel()

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
			// An empty OpenCode effort remains runtime-managed, so a Run issues
			// neither a setup prompt nor an effort set. See ADR-0108.
			name:    "opencode model-managed reasoning issues no effort set",
			runtime: RuntimeSpec{ID: "opencode", Protocol: ProtocolACP, Model: "opencode-model"},
			want: func(gitRoot string) [][]string {
				return [][]string{
					{"--cwd", gitRoot, "--model", "opencode-model", "opencode", "sessions", "ensure", "--name", "roundfix-run-1"},
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

			if got, want := readJSONInvocations(t, harness.invocationsPath), exactSelectionInvocations(tt.want(harness.gitRoot)); !reflect.DeepEqual(got, want) {
				t.Fatalf("unexpected acpx invocations\nwant: %#v\ngot:  %#v", want, got)
			}
		})
	}
}

func TestACPXRunWarmSessionPublishesEffectiveEffortReceipt(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	const model = "openrouter/deepseek/deepseek-v4-pro"
	sink := newCaptureSink("")

	if _, err := harness.runWithSink(context.Background(), RuntimeSpec{
		ID:              "opencode",
		Protocol:        ProtocolACP,
		Model:           model,
		ReasoningEffort: "xhigh",
	}, "roundfix-run-1-task_04", sink); err != nil {
		t.Fatalf("run deferred OpenCode selection: %v", err)
	}

	var receipts []runevent.RunEvent
	for _, event := range sink.Events() {
		if event.Kind == runevent.KindAgentSelectionReceipt {
			receipts = append(receipts, event)
		}
	}
	if len(receipts) != 1 {
		t.Fatalf("selection receipts = %#v, want exactly one", receipts)
	}
	receipt := receipts[0]
	if receipt.RunID != "run-acpx" || receipt.Batch != 7 || receipt.Source != runevent.SourceAgent {
		t.Fatalf("selection receipt identity = %#v", receipt)
	}
	var payload runevent.SelectionReceiptPayload
	if err := json.Unmarshal(receipt.Payload, &payload); err != nil {
		t.Fatalf("decode selection receipt: %v", err)
	}
	if payload.Session != "roundfix-run-1-task_04" || payload.Runtime != "opencode" || payload.Model != model {
		t.Fatalf("selection receipt target = %#v", payload)
	}
	if payload.RequestedReasoningEffort != "xhigh" || payload.ReasoningEffort != "xhigh" || payload.Status != runevent.SelectionReceiptStatusApplied {
		t.Fatalf("selection receipt effort = %#v", payload)
	}
}

func TestACPXRunWarmSessionIsIdempotent(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	runtime := RuntimeSpec{ID: "opencode", Protocol: ProtocolACP, Model: "openrouter/x-ai/grok-4.5", ReasoningEffort: "high"}

	if _, err := harness.run(context.Background(), runtime, "roundfix-run-1-task_04"); err != nil {
		t.Fatalf("first work turn: %v", err)
	}
	if _, err := harness.run(context.Background(), runtime, "roundfix-run-1-task_04"); err != nil {
		t.Fatalf("second work turn: %v", err)
	}

	prompts := readJSONStrings(t, harness.promptsPath)
	if len(prompts) != 3 {
		t.Fatalf("prompt count = %d, want one setup and two work prompts: %#v", len(prompts), prompts)
	}
	warmups := 0
	for _, prompt := range prompts {
		if prompt == acpxDeferredEffortWarmupPrompt {
			warmups++
		}
	}
	if warmups != 1 {
		t.Fatalf("setup prompt count = %d, want 1: %#v", warmups, prompts)
	}
}

func TestACPXRunWarmSessionMismatchStopsBeforeWork(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	const model = "openrouter/x-ai/grok-4.5"
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show":         sessionCapabilitySnapshotFixture(t, model, []string{model}, "effort", "low", []string{"low", "medium", "high"}),
		"set effort value=high": selectionStateFixture(t, "effort", "high", model, []string{model}, "effort", "low", []string{"low", "medium", "high"}),
	}))

	_, err := harness.run(context.Background(), RuntimeSpec{
		ID:              "opencode",
		Protocol:        ProtocolACP,
		Model:           model,
		ReasoningEffort: "high",
	}, "roundfix-run-1-task_04")
	if err == nil || !strings.Contains(err.Error(), CapabilityIssueContradictoryResponse) {
		t.Fatalf("error = %T %v, want contradictory effective effort evidence", err, err)
	}
	prompts := readJSONStrings(t, harness.promptsPath)
	if len(prompts) != 1 || prompts[0] != acpxDeferredEffortWarmupPrompt {
		t.Fatalf("prompt bytes = %#v, want setup only before mismatch", prompts)
	}
}

func TestACPXRunReappliesSelectionForFreshRunnerSessionResume(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh"}

	if _, err := harness.run(context.Background(), runtime, "roundfix-run-1"); err != nil {
		t.Fatalf("first run: %v", err)
	}
	harness.runner = &ACPXRunner{Command: os.Args[0], Environment: harness.runner.Environment, codexSpawn: codexSpawnDependencies{goos: "linux"}}
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
	if got := readJSONInvocations(t, harness.invocationsPath); !reflect.DeepEqual(got, exactSelectionInvocations(want)) {
		t.Fatalf("unexpected resume invocations\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestACPXRunSelectionSetupErrorsPreserveAdapterFailure(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"set reasoning_effort": 2}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"set reasoning_effort": "reasoning rejected\n"}))

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
	if !strings.Contains(err.Error(), "apply Agent Selection") || !strings.Contains(err.Error(), "set reasoning_effort") {
		t.Fatalf("expected reasoning operation context, got %v", err)
	}
	want := [][]string{
		{"--cwd", harness.gitRoot, "--model", "gpt-5.5", "codex", "sessions", "ensure", "--name", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-1"},
	}
	if got := readJSONInvocations(t, harness.invocationsPath); !reflect.DeepEqual(got, exactSelectionInvocations(want)) {
		t.Fatalf("unexpected failed-selection invocations\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestACPXRunCustomCodexCommandRequiresOfficialAdapterIdentity(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	runtime := RuntimeSpec{
		ID:              "codex-custom",
		Protocol:        ProtocolStdio,
		Command:         filepath.Join(harness.adapterDir, "custom-acp") + " --stdio",
		Model:           "gpt-custom",
		ReasoningEffort: "xhigh",
	}

	_, err := harness.run(context.Background(), runtime, "roundfix-run-1")

	if err == nil {
		t.Fatal("expected unproven custom Codex adapter identity to fail")
	}
	var lineageErr *AdapterLineageError
	if !errors.As(err, &lineageErr) || lineageErr.Command != runtime.Command {
		t.Fatalf("expected custom command lineage failure, got %T %v", err, err)
	}
	if !strings.Contains(err.Error(), CodexAdapterInstallCommand()) {
		t.Fatalf("expected deterministic official adapter action, got %v", err)
	}
	if _, statErr := os.Stat(harness.invocationsPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("adapter identity failure reached acpx session mutation: stat error %v", statErr)
	}
}

func TestACPXRunEnsuresSessionOncePerRunnerAndSessionName(t *testing.T) {
	t.Parallel()

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
	if !reflect.DeepEqual(invocations, exactSelectionInvocations(want)) {
		t.Fatalf("unexpected acpx invocations\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestACPXRunCodexUsesConfiguredCleanPathOnDarwin(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	envPath := filepath.Join(harness.gitRoot, "codex-paths.jsonl")
	harness.setEnv(fakeACPXCodexPath, envPath)
	harness.setEnv(codexPathEnv, "/configured/clean/codex")
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

	wantEnv := []string{"/configured/clean/codex", "/configured/clean/codex", "/configured/clean/codex", "/configured/clean/codex"}
	if got := readJSONStrings(t, envPath); !reflect.DeepEqual(got, wantEnv) {
		t.Fatalf("unexpected CODEX_PATH values\nwant: %#v\ngot:  %#v", wantEnv, got)
	}
	if got, want := probe.quarantineCalls, []string{"/configured/clean/codex"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("expected only configured codex to be inspected\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestACPXRunCodexFallsBackToCleanPathWhenConfiguredPathIsQuarantined(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	envPath := filepath.Join(harness.gitRoot, "codex-paths.jsonl")
	harness.setEnv(fakeACPXCodexPath, envPath)
	harness.setEnv(codexPathEnv, "/configured/quarantined/codex")
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

	wantEnv := []string{"/path/clean/codex", "/path/clean/codex", "/path/clean/codex", "/path/clean/codex"}
	if got := readJSONStrings(t, envPath); !reflect.DeepEqual(got, wantEnv) {
		t.Fatalf("unexpected CODEX_PATH values\nwant: %#v\ngot:  %#v", wantEnv, got)
	}
	wantInspections := []string{"/configured/quarantined/codex", "/path/clean/codex"}
	if !reflect.DeepEqual(probe.quarantineCalls, wantInspections) {
		t.Fatalf("unexpected quarantine inspections\nwant: %#v\ngot:  %#v", wantInspections, probe.quarantineCalls)
	}
}

func TestACPXRunCodexSurfacesQuarantinedPathWithoutSpawning(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(codexPathEnv, "")
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
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(codexPathEnv, "/configured/clean/codex")
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
	t.Parallel()

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
	if got := readJSONInvocations(t, harness.invocationsPath); !reflect.DeepEqual(got, exactSelectionInvocations(want)) {
		t.Fatalf("unexpected non-darwin acpx invocations\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestACPXRunFailsEnsureWithStderrTail(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions ensure": 2}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"sessions ensure": "ensure rejected\n"}))

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
	if !strings.Contains(err.Error(), "ensure live Agent Session") || !strings.Contains(err.Error(), "--- acpx stderr tail ---\nensure rejected") {
		t.Fatalf("expected ensure error with delimited stderr tail, got %q", err.Error())
	}
}

func TestACPXRunAppliesFullAccessSessionSetup(t *testing.T) {
	t.Parallel()

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

			if got, want := readJSONInvocations(t, harness.invocationsPath), exactSelectionInvocations(tt.want(harness.gitRoot)); !reflect.DeepEqual(got, want) {
				t.Fatalf("unexpected acpx invocations\nwant: %#v\ngot:  %#v", want, got)
			}
		})
	}
}

func TestACPXRunFailsSetupWhenFullAccessModeFails(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"set-mode": 2}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"set-mode": "mode rejected\n"}))

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
	if !reflect.DeepEqual(invocations, exactSelectionInvocations(want)) {
		t.Fatalf("unexpected invocations after setup failure\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestACPXRunWarnsWhenCodexSandboxPresetUnavailable(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"set sandbox_mode": 2}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{
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
	if !reflect.DeepEqual(invocations, exactSelectionInvocations(want)) {
		t.Fatalf("unexpected invocations\nwant: %#v\ngot:  %#v", want, invocations)
	}
}

func TestACPXRunCancelsPromptCooperatively(t *testing.T) {
	t.Parallel()

	harness := newBlockingFakeACPXHarness(t, true)
	assertCancellationFixturePaths(t, harness)
	clock := newFakeCancellationClock()
	harness.runner.cancelClock = clock
	ctx, cancel := context.WithCancel(context.Background())
	resultCh := make(chan error, 1)
	go func() {
		_, err := harness.run(ctx, RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, "roundfix-run-1")
		resultCh <- err
	}()
	harness.waitForMilestone(t, "prompt start", harness.milestones.promptStarted)
	cancel()
	harness.waitForMilestone(t, "cancel completion", harness.milestones.cancelCompleted)
	graceTimer := clock.waitForTimer(t, 0)
	harness.waitForMilestone(t, "prompt completion", harness.milestones.promptCompleted)

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
	if graceTimer.Fired() {
		t.Fatal("expected cooperative prompt completion before grace timer fired")
	}
	if !graceTimer.Stopped() {
		t.Fatal("expected cooperative prompt completion to stop grace timer")
	}
	if timers := clock.timersSnapshot(); len(timers) != 1 {
		t.Fatalf("expected only the cooperative grace timer, got %d timers", len(timers))
	}
	assertNoFile(t, "close completion", harness.milestones.closeCompleted)
	invocations := readJSONInvocations(t, harness.invocationsPath)
	want := [][]string{
		{"--cwd", harness.gitRoot, "--model", "gpt-test", "codex", "sessions", "ensure", "--name", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "gpt-test", "codex", "prompt", "-s", "roundfix-run-1", "-f", "-"},
		{"--cwd", harness.gitRoot, "codex", "cancel", "-s", "roundfix-run-1"},
	}
	if !reflect.DeepEqual(invocations, exactSelectionInvocations(want)) {
		t.Fatalf("unexpected cooperative cancellation invocation order\nwant: %#v\ngot:  %#v", exactSelectionInvocations(want), invocations)
	}
}

func TestACPXRunClosesSessionAfterCancelGracePeriod(t *testing.T) {
	t.Parallel()

	harness := newBlockingFakeACPXHarness(t, false)
	assertCancellationFixturePaths(t, harness)
	clock := newFakeCancellationClock()
	harness.runner.cancelClock = clock
	ctx, cancel := context.WithCancel(context.Background())
	resultCh := make(chan error, 1)
	go func() {
		_, err := harness.run(ctx, RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, "roundfix-run-1")
		resultCh <- err
	}()
	harness.waitForMilestone(t, "prompt start", harness.milestones.promptStarted)
	cancel()
	harness.waitForMilestone(t, "cancel completion", harness.milestones.cancelCompleted)
	graceTimer := clock.waitForTimer(t, 0)
	if graceTimer.Duration() != harness.stopGrace {
		t.Fatalf("expected grace timer duration %s, got %s", harness.stopGrace, graceTimer.Duration())
	}
	if !graceTimer.Fire(time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)) {
		t.Fatal("expected first grace timer fire to trigger fallback close")
	}
	harness.waitForMilestone(t, "close completion", harness.milestones.closeCompleted)
	postCloseTimer := clock.waitForTimer(t, 1)

	err := receiveError(t, resultCh)
	var stopErr StopError
	if !errors.As(err, &stopErr) {
		t.Fatalf("expected StopError, got %T %v", err, err)
	}
	if !stopErr.Killed {
		t.Fatalf("expected fallback close after cancel, got %#v", stopErr)
	}
	if !postCloseTimer.Stopped() {
		t.Fatal("expected close milestone to release blocked prompt before post-close timer fired")
	}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	want := [][]string{
		{"--cwd", harness.gitRoot, "--model", "gpt-test", "codex", "sessions", "ensure", "--name", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "set", "reasoning_effort", "xhigh", "-s", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "--format", "json", "--json-strict", "--approve-all", "--model", "gpt-test", "codex", "prompt", "-s", "roundfix-run-1", "-f", "-"},
		{"--cwd", harness.gitRoot, "codex", "cancel", "-s", "roundfix-run-1"},
		{"--cwd", harness.gitRoot, "codex", "sessions", "close", "roundfix-run-1"},
	}
	if !reflect.DeepEqual(invocations, exactSelectionInvocations(want)) {
		t.Fatalf("unexpected cancellation invocation order\nwant: %#v\ngot:  %#v", exactSelectionInvocations(want), invocations)
	}
}

func TestACPXRunCancellationCommandFailuresWarnAndContinue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		exitBy      map[string]int
		stderrBy    map[string]string
		wantWarning string
	}{
		{
			name:        "cancel failure still reaches fallback close",
			exitBy:      map[string]int{"cancel": 2},
			stderrBy:    map[string]string{"cancel": "cancel rejected\n"},
			wantWarning: "cancel acpx Agent Session",
		},
		{
			name:        "close failure still waits for prompt termination",
			exitBy:      map[string]int{"sessions close": 2},
			stderrBy:    map[string]string{"sessions close": "close rejected\n"},
			wantWarning: "close acpx Agent Session",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			harness := newBlockingFakeACPXHarness(t, false)
			const promptStartedSummary = "fake acpx prompt started"
			harness.setEnv(fakeACPXStartEvent, promptStartedSummary)
			promptStarted := newCaptureSink(promptStartedSummary)
			clock := newFakeCancellationClock()
			harness.runner.cancelClock = clock
			var warnings []string
			harness.runner.warnf = func(format string, args ...any) {
				warnings = append(warnings, fmt.Sprintf(format, args...))
			}
			harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, tt.exitBy))
			harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, tt.stderrBy))
			ctx, cancel := context.WithCancel(context.Background())
			resultCh := make(chan error, 1)
			go func() {
				_, err := harness.runWithSink(ctx, RuntimeSpec{ID: "codex", Protocol: ProtocolACP}, "roundfix-run-1", promptStarted)
				resultCh <- err
			}()
			select {
			case <-promptStarted.done:
			case <-time.After(agentWaitBudget):
				t.Fatalf("timed out waiting for prompt start event; Agent did not start within %s", agentWaitBudget)
			}
			cancel()
			harness.waitForMilestone(t, "cancel completion", harness.milestones.cancelCompleted)
			if !clock.waitForTimer(t, 0).Fire(time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)) {
				t.Fatal("expected grace timer fire to trigger fallback close")
			}
			harness.waitForMilestone(t, "close completion", harness.milestones.closeCompleted)
			postCloseTimer := clock.waitForTimer(t, 1)

			err := receiveError(t, resultCh)
			var stopErr StopError
			if !errors.As(err, &stopErr) {
				t.Fatalf("expected StopError, got %T %v", err, err)
			}
			if !stopErr.Killed {
				t.Fatalf("expected fallback close result, got %#v", stopErr)
			}
			if !postCloseTimer.Stopped() {
				t.Fatal("expected prompt termination to stop post-close timer")
			}
			if !containsWarning(warnings, tt.wantWarning) {
				t.Fatalf("expected warning containing %q, got %#v", tt.wantWarning, warnings)
			}
			if !containsInvocation(readJSONInvocations(t, harness.invocationsPath), []string{"--cwd", harness.gitRoot, "codex", "sessions", "close", "roundfix-run-1"}) {
				t.Fatalf("expected fallback close invocation after warning, got %#v", readJSONInvocations(t, harness.invocationsPath))
			}
		})
	}
}

func TestACPXRunnerCancellationClockDefaultsToRealTimer(t *testing.T) {
	t.Parallel()

	if got := stopGrace(0); got != 10*time.Second {
		t.Fatalf("expected default stop grace to remain 10s, got %s", got)
	}
	timer := (ACPXRunner{}).cancellationClock().NewTimer(time.Hour)
	defer timer.Stop()
	select {
	case firedAt := <-timer.C():
		t.Fatalf("expected real default cancellation timer not to fire immediately, got %s", firedAt)
	default:
	}
}

func TestFakeCancellationClock(t *testing.T) {
	t.Parallel()

	t.Run("records creation order and waits without firing", func(t *testing.T) {
		clock := newFakeCancellationClock()
		createFirst := make(chan struct{})
		go func() {
			<-createFirst
			clock.NewTimer(10 * time.Second)
		}()

		close(createFirst)
		first := clock.waitForTimer(t, 0)
		clock.NewTimer(20 * time.Second)
		second := clock.waitForTimer(t, 1)
		timers := clock.timersSnapshot()

		if len(timers) != 2 {
			t.Fatalf("expected two created timers, got %d", len(timers))
		}
		if timers[0] != first || timers[1] != second {
			t.Fatalf("timers were not recorded in creation order")
		}
		if first.Duration() != 10*time.Second || second.Duration() != 20*time.Second {
			t.Fatalf("unexpected timer durations: first=%s second=%s", first.Duration(), second.Duration())
		}
		select {
		case firedAt := <-first.C():
			t.Fatalf("waitForTimer fired timer unexpectedly at %s", firedAt)
		default:
		}
	})

	t.Run("fires each timer at most once", func(t *testing.T) {
		clock := newFakeCancellationClock()
		clock.NewTimer(time.Second)
		timer := clock.waitForTimer(t, 0)
		firedAt := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)

		if !timer.Fire(firedAt) {
			t.Fatal("expected first fire to report active timer")
		}
		select {
		case got := <-timer.C():
			if !got.Equal(firedAt) {
				t.Fatalf("unexpected fire time: got %s want %s", got, firedAt)
			}
		default:
			t.Fatal("expected first fire to emit one event")
		}
		if timer.Fire(firedAt.Add(time.Second)) {
			t.Fatal("expected second fire to report inactive timer")
		}
		select {
		case got := <-timer.C():
			t.Fatalf("expected no second timer event, got %s", got)
		default:
		}
	})

	t.Run("stopped timer cannot be fired as active", func(t *testing.T) {
		clock := newFakeCancellationClock()
		clock.NewTimer(time.Second)
		timer := clock.waitForTimer(t, 0)

		if !timer.Stop() {
			t.Fatal("expected first stop to report active timer")
		}
		if !timer.Stopped() {
			t.Fatal("expected stopped timer state")
		}
		if timer.Fire(time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)) {
			t.Fatal("expected stopped timer fire to report inactive timer")
		}
		if timer.Stop() {
			t.Fatal("expected second stop to report inactive timer")
		}
		select {
		case got := <-timer.C():
			t.Fatalf("expected stopped timer to emit no event, got %s", got)
		default:
		}
	})
}

func TestACPXEndSessionClosesBestEffort(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions close": 1}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{"sessions close": "close rejected\n"}))
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
	t.Parallel()

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
	if len(events) != 3 {
		t.Fatalf("expected work-started then two update events, got %d: %+v", len(events), events)
	}
	if events[0].Kind != runevent.KindAgentStatus || !strings.Contains(string(events[0].Payload), AgentWorkStartedStatus) {
		t.Fatalf("expected work-started status before Agent output, got %+v", events[0])
	}
	expectedKinds := []runevent.Kind{runevent.KindAgentMessage, runevent.KindAgentThought}
	expectedPayloads := []string{messageLine, thoughtLine}
	for index, event := range events[1:] {
		if event.Kind != expectedKinds[index] {
			t.Fatalf("expected event %d kind %q, got %q", index+1, expectedKinds[index], event.Kind)
		}
		if event.RunID != "run-acpx" || event.Batch != 7 || event.Source != runevent.SourceAgent {
			t.Fatalf("expected Run identity on event %d, got %+v", index+1, event)
		}
		if string(event.Payload) != expectedPayloads[index] {
			t.Fatalf("expected byte-identical payload for event %d\nwant: %q\ngot:  %q", index+1, expectedPayloads[index], string(event.Payload))
		}
	}
	if logContent := readFile(t, run.logPath); logContent != stdout {
		t.Fatalf("expected agent log to contain every stdout line in order\nwant: %q\ngot:  %q", stdout, logContent)
	}
}

func TestWorkStartedBoundaryPublishesOnFirstAgentOutput(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	messageLine := acpxUpdateLine(`{"sessionId":"sess-1","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"hello"}}}`)
	runner := &ACPXRunner{
		Command: os.Args[0],
		Environment: environmentForTest(
			fakeACPXEnv+"=1",
			fakeACPXStdout+"="+messageLine+acpxPromptResponseLine("end_turn"),
		),
		Now: func() time.Time {
			return time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC)
		},
		codexSpawn: codexSpawnDependencies{goos: "linux"},
	}
	sink := newCaptureSink("")
	req := ACPXPromptRequest{
		ExecuteRequest: ExecuteRequest{
			Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP},
			RunID:   "run-boundary",
			Batch:   rounds.Batch{Number: 2},
			Prompt:  "prompt",
			GitRoot: dir,
		},
		Session: "roundfix-run-boundary",
	}

	for turn := 0; turn < 2; turn++ {
		if _, err := runner.RunPrompt(context.Background(), req, sink); err != nil {
			t.Fatalf("turn %d: RunPrompt() error = %v", turn+1, err)
		}
	}

	events := sink.Events()
	if len(events) != 3 {
		t.Fatalf("events = %+v, want one work-started status followed by two Agent messages", events)
	}
	if events[0].Kind != runevent.KindAgentStatus || !strings.Contains(string(events[0].Payload), AgentWorkStartedStatus) {
		t.Fatalf("first event = %+v, want Agent work-started status", events[0])
	}
	for index := 1; index < len(events); index++ {
		if events[index].Kind != runevent.KindAgentMessage {
			t.Fatalf("event %d kind = %q, want %q", index, events[index].Kind, runevent.KindAgentMessage)
		}
	}
}

func TestWorkStartedBoundaryReportsSelectionFailureWithoutOutput(t *testing.T) {
	t.Parallel()

	run := runFakeACPXPrompt(t, fakeACPXPrompt{
		stderr:   "usage limit exhausted\n",
		exitCode: 1,
	})

	var selectionErr *SelectionFailure
	if !errors.As(run.err, &selectionErr) {
		t.Fatalf("RunPrompt() error = %T %v, want *SelectionFailure", run.err, run.err)
	}
	var batchErr *BatchFailureError
	if errors.As(run.err, &batchErr) {
		t.Fatalf("RunPrompt() error = %T %v, must be distinct from *BatchFailureError", run.err, run.err)
	}
	events := run.sink.Events()
	if len(events) != 1 || events[0].Kind != runevent.KindAgentStatus ||
		!strings.Contains(string(events[0].Payload), AgentSelectionFailedStatus) {
		t.Fatalf("events = %+v, want one Agent selection-failed status", events)
	}
	if run.sink.HasStatus(AgentWorkStartedStatus) {
		t.Fatalf("events = %+v, must not publish Agent work-started without Agent output", events)
	}
}

func TestWorkStartedBoundaryIgnoresInertSessionSetup(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	messageLine := acpxUpdateLine(`{"sessionId":"sess-1","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"hello"}}}`)
	runner := &ACPXRunner{
		Command: os.Args[0],
		Environment: environmentForTest(
			fakeACPXEnv+"=1",
			fakeACPXStdout+"="+messageLine+acpxPromptResponseLine("end_turn"),
		),
		Now:        func() time.Time { return time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC) },
		codexSpawn: codexSpawnDependencies{goos: "linux"},
	}
	sink := newCaptureSink("")
	req := ACPXPromptRequest{
		ExecuteRequest: ExecuteRequest{
			Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP},
			RunID:   "run-boundary",
			Batch:   rounds.Batch{Number: 2},
			Prompt:  acpxDeferredEffortWarmupPrompt,
			GitRoot: dir,
		},
		Session: "roundfix-run-boundary",
		Inert:   true,
	}

	if _, err := runner.RunPrompt(context.Background(), req, sink); err != nil {
		t.Fatalf("inert RunPrompt() error = %v", err)
	}
	if sink.HasStatus(AgentWorkStartedStatus) {
		t.Fatalf("inert Session setup published Agent work-started: %+v", sink.Events())
	}
	req.Inert = false
	req.Prompt = "do the work"
	if _, err := runner.RunPrompt(context.Background(), req, sink); err != nil {
		t.Fatalf("work RunPrompt() error = %v", err)
	}
	if got := countStatusEventsForTest(sink.Events(), AgentWorkStartedStatus); got != 1 {
		t.Fatalf("work-started statuses = %d, want one after inert Session setup", got)
	}
}

func countStatusEventsForTest(events []runevent.RunEvent, status string) int {
	count := 0
	for _, event := range events {
		if event.Kind == runevent.KindAgentStatus && strings.Contains(string(event.Payload), status) {
			count++
		}
	}
	return count
}

func TestACPXRunPromptAllowsEmptyLogPathAndStillJournals(t *testing.T) {
	t.Parallel()

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
	if len(events) != 2 || events[0].Kind != runevent.KindAgentStatus || events[1].Kind != runevent.KindAgentMessage {
		t.Fatalf("expected work-started then Agent message without log file, got %+v", events)
	}
	if !strings.Contains(string(events[1].Payload), "hello") {
		t.Fatalf("expected raw Agent payload preserved, got %s", events[1].Payload)
	}
}

func TestACPXPromptExitClassificationMatrix(t *testing.T) {
	t.Parallel()

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
			name:     "no result exit one becomes selection failure",
			exitCode: 1,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var selectionErr *SelectionFailure
				if !errors.As(err, &selectionErr) {
					t.Fatalf("expected SelectionFailure, got %T %v", err, err)
				}
				if selectionErr.Reason != acpxExitReasonAgentProtocol {
					t.Fatalf("selection failure reason = %q, want %q", selectionErr.Reason, acpxExitReasonAgentProtocol)
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
			name:     "Agent output exit one remains Batch failure",
			stdout:   acpxUpdateLine(`{"sessionId":"sess-1","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"started"}}}`),
			exitCode: 1,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var batchErr *BatchFailureError
				if !errors.As(err, &batchErr) {
					t.Fatalf("expected BatchFailureError after Agent output, got %T %v", err, err)
				}
				var selectionErr *SelectionFailure
				if errors.As(err, &selectionErr) {
					t.Fatalf("error after Agent output = %T %v, must not be SelectionFailure", err, err)
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
			name:     "partial stream without Agent output becomes selection failure",
			stdout:   `{"jsonrpc":"2.0","id":1,"result":`,
			exitCode: 1,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var selectionErr *SelectionFailure
				if !errors.As(err, &selectionErr) {
					t.Fatalf("expected SelectionFailure, got %T %v", err, err)
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
	t.Parallel()

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
	t.Parallel()

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
				var selectionErr *SelectionFailure
				if !errors.As(err, &selectionErr) {
					t.Fatalf("expected SelectionFailure, got %T %v", err, err)
				}
				if selectionErr.Reason != acpxExitReasonAgentProtocol {
					t.Fatalf("expected protocol reason, got %q", selectionErr.Reason)
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
				var selectionErr *SelectionFailure
				if !errors.As(err, &selectionErr) {
					t.Fatalf("expected SelectionFailure, got %T %v", err, err)
				}
				if selectionErr.Reason != acpxExitReasonTimeout {
					t.Fatalf("expected timeout reason, got %q", selectionErr.Reason)
				}
			},
		},
		{
			name:     "all permissions denied failure",
			exitCode: 5,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var selectionErr *SelectionFailure
				if !errors.As(err, &selectionErr) {
					t.Fatalf("expected SelectionFailure, got %T %v", err, err)
				}
				if selectionErr.Reason != acpxExitReasonPermissionsDenied {
					t.Fatalf("expected permissions denied reason, got %q", selectionErr.Reason)
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
	adapterDir      string
	invocationsPath string
	promptsPath     string
	startedPath     string
	stopGrace       time.Duration
	milestones      fakeACPXMilestones
}

func newFakeACPXHarness(t *testing.T) *fakeACPXHarness {
	t.Helper()
	dir := t.TempDir()
	invocationsPath := filepath.Join(dir, "invocations.jsonl")
	promptsPath := filepath.Join(dir, "prompts.jsonl")
	homeDir := filepath.Join(dir, "home")
	if err := os.MkdirAll(homeDir, 0o755); err != nil {
		t.Fatalf("create fake HOME: %v", err)
	}
	adapterDir := filepath.Join(dir, "adapters")
	if err := os.MkdirAll(adapterDir, 0o755); err != nil {
		t.Fatalf("create fake adapter dir: %v", err)
	}
	for _, command := range []string{"codex-acp", "claude-agent-acp", "custom-acp", "opencode", "npx"} {
		installFakeAdapter(t, adapterDir, command)
	}
	config := fmt.Sprintf(`{"agents":{"codex":{"command":%q,"args":["-y",%q]},"claude":{"command":%q,"args":["-y",%q]},"opencode":{"command":%q}}}`,
		filepath.Join(adapterDir, "npx"), CodexAdapterPackage,
		filepath.Join(adapterDir, "npx"), ClaudeAdapterPackage+"@"+PinnedClaudeAdapterVersion,
		filepath.Join(adapterDir, "opencode"),
	)
	writeACPXConfigAtHomeForTest(t, homeDir, config)
	return &fakeACPXHarness{
		runner: &ACPXRunner{
			Command: os.Args[0],
			Environment: environmentForTest(
				"HOME="+homeDir,
				fakeACPXEnv+"=1",
				fakeACPXInvokes+"="+invocationsPath,
				fakeACPXPrompts+"="+promptsPath,
				fakeACPXStdoutBy+"="+mustJSONForTest(t, map[string]string{"prompt": acpxPromptResponseLine("end_turn")}),
			),
			Now: func() time.Time {
				return time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
			},
			codexSpawn: codexSpawnDependencies{goos: "linux"},
		},
		gitRoot:         dir,
		adapterDir:      adapterDir,
		invocationsPath: invocationsPath,
		promptsPath:     promptsPath,
	}
}

func (harness *fakeACPXHarness) setEnv(key string, value string) {
	harness.runner.Environment = environmentForBase(harness.runner.Environment, key+"="+value)
}

func installFakeAdapter(t *testing.T, dir string, command string) {
	t.Helper()
	path := filepath.Join(dir, command)
	content := "#!/bin/sh\nexit 0\n"
	if command == "codex-acp" {
		content = "#!/bin/sh\nprintf '%s\\n' '" + CodexAdapterPackage + " " + PinnedCodexAdapterVersion + "'\n"
	}
	if command == "claude-agent-acp" {
		content = "#!/bin/sh\nprintf '%s\\n' '" + PinnedClaudeAdapterVersion + "'\n"
	}
	if command == "npx" {
		content = "#!/bin/sh\ncase \"$*\" in\n  *claude-agent-acp*) printf '%s\\n' '" + PinnedClaudeAdapterVersion + "' ;;\n  *) printf '%s\\n' '" + CodexAdapterPackage + " " + PinnedCodexAdapterVersion + "' ;;\nesac\n"
	}
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("write fake adapter %s: %v", command, err)
	}
}

func installFakeVersionAdapter(t *testing.T, output string) string {
	t.Helper()
	return installFakeNamedVersionAdapter(t, "codex-acp", output)
}

func installFakeNamedVersionAdapter(t *testing.T, name string, output string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	content := "#!/bin/sh\nprintf '%s\\n' '" + output + "'\n"
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("write fake version adapter: %v", err)
	}
	return path
}

func installSymlinkedPackageAdapter(t *testing.T, packageName string, executableName string, output string) string {
	t.Helper()
	root := t.TempDir()
	packageDir := filepath.Join(root, "node_modules", filepath.FromSlash(packageName), "bin")
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatalf("create fake package adapter directory: %v", err)
	}
	target := filepath.Join(packageDir, "adapter")
	content := "#!/bin/sh\nprintf '%s\\n' '" + output + "'\n"
	if err := os.WriteFile(target, []byte(content), 0o755); err != nil {
		t.Fatalf("write fake package adapter: %v", err)
	}
	binDir := filepath.Join(root, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("create fake adapter bin directory: %v", err)
	}
	command := filepath.Join(binDir, executableName)
	if err := os.Symlink(target, command); err != nil {
		t.Fatalf("symlink fake package adapter: %v", err)
	}
	return command
}

func writeACPXConfigForTest(t *testing.T, content string) []string {
	t.Helper()
	homeDir := t.TempDir()
	writeACPXConfigAtHomeForTest(t, homeDir, content)
	return environmentForTest("HOME=" + homeDir)
}

func writeACPXConfigAtHomeForTest(t *testing.T, homeDir string, content string) {
	t.Helper()
	configPath := filepath.Join(homeDir, ".acpx", "config.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("create acpx config dir: %v", err)
	}
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write acpx config: %v", err)
	}
}

func environmentForTest(overrides ...string) []string {
	return environmentForBase(os.Environ(), overrides...)
}

func environmentForBase(base []string, overrides ...string) []string {
	keys := make(map[string]struct{}, len(overrides))
	for _, override := range overrides {
		key, _, _ := strings.Cut(override, "=")
		keys[key] = struct{}{}
	}
	environment := make([]string, 0, len(base)+len(overrides))
	for _, entry := range base {
		key, _, _ := strings.Cut(entry, "=")
		if _, replaced := keys[key]; !replaced {
			environment = append(environment, entry)
		}
	}
	return append(environment, overrides...)
}

func newBlockingFakeACPXHarness(t *testing.T, exitAfterCancel bool) *fakeACPXHarness {
	t.Helper()
	harness := newFakeACPXHarness(t)
	harness.milestones = newFakeACPXMilestones(t, harness.gitRoot)
	harness.startedPath = harness.milestones.promptStarted
	harness.stopGrace = 2 * time.Second
	harness.setEnv(fakeACPXBlock, "1")
	harness.setEnv(fakeACPXStarted, harness.startedPath)
	harness.setEnv(fakeACPXCanceled, harness.milestones.cancelCompleted)
	harness.setEnv(fakeACPXClosed, harness.milestones.closeCompleted)
	harness.setEnv(fakeACPXPromptDone, harness.milestones.promptCompleted)
	if exitAfterCancel {
		harness.setEnv(fakeACPXExitCancel, "1")
	}
	return harness
}

type fakeACPXMilestones struct {
	promptStarted   string
	promptCompleted string
	cancelCompleted string
	closeCompleted  string
}

func newFakeACPXMilestones(t *testing.T, dir string) fakeACPXMilestones {
	t.Helper()
	milestoneDir := filepath.Join(dir, "milestones")
	if err := os.MkdirAll(milestoneDir, 0o755); err != nil {
		t.Fatalf("create fake acpx milestone dir: %v", err)
	}
	return fakeACPXMilestones{
		promptStarted:   filepath.Join(milestoneDir, "prompt-started"),
		promptCompleted: filepath.Join(milestoneDir, "prompt-completed"),
		cancelCompleted: filepath.Join(milestoneDir, "cancel-completed"),
		closeCompleted:  filepath.Join(milestoneDir, "close-completed"),
	}
}

func (harness *fakeACPXHarness) waitForMilestone(t *testing.T, name string, path string) {
	t.Helper()
	waitForACPXMilestone(t, name, path, harness.invocationsPath)
}

func assertCancellationFixturePaths(t *testing.T, harness *fakeACPXHarness) {
	t.Helper()
	paths := map[string]string{
		"invocation log":   harness.invocationsPath,
		"prompt milestone": harness.milestones.promptStarted,
		"prompt complete":  harness.milestones.promptCompleted,
		"cancel milestone": harness.milestones.cancelCompleted,
		"close milestone":  harness.milestones.closeCompleted,
	}
	seen := map[string]string{}
	for name, path := range paths {
		rel, err := filepath.Rel(harness.gitRoot, path)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || filepath.IsAbs(rel) {
			t.Fatalf("%s path %q is not isolated under temp dir %q", name, path, harness.gitRoot)
		}
		if previous, ok := seen[path]; ok {
			t.Fatalf("%s path duplicates %s path: %q", name, previous, path)
		}
		seen[path] = name
	}
}

func containsWarning(warnings []string, want string) bool {
	for _, warning := range warnings {
		if strings.Contains(warning, want) {
			return true
		}
	}
	return false
}

func assertNoFile(t *testing.T, name string, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected no fake acpx %s marker at %s", name, path)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat fake acpx %s marker %s: %v", name, path, err)
	}
}

type fakeCancellationClock struct {
	mu      sync.Mutex
	timers  []*fakeCancellationTimer
	created chan struct{}
}

func newFakeCancellationClock() *fakeCancellationClock {
	return &fakeCancellationClock{
		created: make(chan struct{}, 16),
	}
}

func (clock *fakeCancellationClock) NewTimer(duration time.Duration) cancellationTimer {
	clock.mu.Lock()
	timer := &fakeCancellationTimer{
		duration: duration,
		c:        make(chan time.Time, 1),
	}
	clock.timers = append(clock.timers, timer)
	clock.mu.Unlock()

	select {
	case clock.created <- struct{}{}:
	default:
	}
	return timer
}

func (clock *fakeCancellationClock) waitForTimer(t *testing.T, index int) *fakeCancellationTimer {
	t.Helper()
	deadline := time.NewTimer(agentWaitBudget)
	defer deadline.Stop()
	for {
		clock.mu.Lock()
		if index < len(clock.timers) {
			timer := clock.timers[index]
			clock.mu.Unlock()
			return timer
		}
		created := len(clock.timers)
		clock.mu.Unlock()

		select {
		case <-clock.created:
		case <-deadline.C:
			t.Fatalf("timed out waiting for cancellation timer %d; created timers: %d", index+1, created)
		}
	}
}

func (clock *fakeCancellationClock) timersSnapshot() []*fakeCancellationTimer {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	timers := make([]*fakeCancellationTimer, len(clock.timers))
	copy(timers, clock.timers)
	return timers
}

type fakeCancellationTimer struct {
	mu       sync.Mutex
	duration time.Duration
	c        chan time.Time
	fired    bool
	stopped  bool
}

func (timer *fakeCancellationTimer) C() <-chan time.Time {
	return timer.c
}

func (timer *fakeCancellationTimer) Stop() bool {
	timer.mu.Lock()
	defer timer.mu.Unlock()
	if timer.fired || timer.stopped {
		return false
	}
	timer.stopped = true
	return true
}

func (timer *fakeCancellationTimer) Fire(firedAt time.Time) bool {
	timer.mu.Lock()
	if timer.fired || timer.stopped {
		timer.mu.Unlock()
		return false
	}
	timer.fired = true
	timerC := timer.c
	timer.mu.Unlock()

	timerC <- firedAt
	return true
}

func (timer *fakeCancellationTimer) Duration() time.Duration {
	timer.mu.Lock()
	defer timer.mu.Unlock()
	return timer.duration
}

func (timer *fakeCancellationTimer) Fired() bool {
	timer.mu.Lock()
	defer timer.mu.Unlock()
	return timer.fired
}

func (timer *fakeCancellationTimer) Stopped() bool {
	timer.mu.Lock()
	defer timer.mu.Unlock()
	return timer.stopped
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
		// No reasoning-effort default: an empty effort remains runtime-managed.
		// See ADR-0108.
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
	stopGrace := harness.stopGrace
	if stopGrace <= 0 {
		stopGrace = 20 * time.Millisecond
	}
	return harness.runner.Run(ctx, ExecuteRequest{
		Runtime:   selectedRuntime(runtime),
		RunID:     "run-acpx",
		Batch:     rounds.Batch{Number: 7},
		LogPath:   filepath.Join(harness.gitRoot, "runs", "run-acpx", "agent", "batch-007.log"),
		Prompt:    "prompt",
		GitRoot:   harness.gitRoot,
		StopGrace: stopGrace,
		Session:   SessionRef{Name: sessionName},
	}, sink)
}

func runFakeACPXPrompt(t *testing.T, prompt fakeACPXPrompt) fakeACPXRun {
	t.Helper()
	dir := t.TempDir()
	argsPath := filepath.Join(dir, "args.json")
	promptPath := filepath.Join(dir, "prompt.txt")
	environment := environmentForTest(
		fakeACPXEnv+"=1",
		fakeACPXArgsPath+"="+argsPath,
		fakeACPXPromptPath+"="+promptPath,
		fakeACPXStdout+"="+prompt.stdout,
		fakeACPXStderr+"="+prompt.stderr,
		fakeACPXExitCode+"="+strconv.Itoa(prompt.exitCode),
	)

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
		Command:     os.Args[0],
		Environment: environment,
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
	environment := environmentForTest(
		fakeACPXEnv+"=1",
		fakeACPXInvokes+"="+invocationsPath,
		fakeACPXStdout+"="+version+"\n",
	)

	err := (ACPXRunner{Command: os.Args[0], Environment: environment}).Probe(context.Background(), ProbeRequest{Runtime: runtime})
	return readJSONInvocations(t, invocationsPath), err
}

func assertDisposableEnsureInvocation(t *testing.T, got []string, workDir string, runtime string, model string) string {
	t.Helper()
	wantPrefix := []string{"--cwd", workDir, "--model", model, runtime, "sessions", "ensure", "--name"}
	if len(got) != len(wantPrefix)+1 || !reflect.DeepEqual(got[:len(wantPrefix)], wantPrefix) {
		t.Fatalf("unexpected disposable ensure invocation\nwant prefix: %#v + session\ngot:         %#v", wantPrefix, got)
	}
	sessionName := got[len(got)-1]
	if !strings.HasPrefix(sessionName, "roundfix-preflight-") {
		t.Fatalf("expected disposable preflight session name, got %q", sessionName)
	}
	return sessionName
}

func assertDisposableCatalogueEnsureInvocation(t *testing.T, got []string, workDir string, runtime string) string {
	t.Helper()
	wantPrefix := []string{"--cwd", workDir, runtime, "sessions", "ensure", "--name"}
	if len(got) != len(wantPrefix)+1 || !reflect.DeepEqual(got[:len(wantPrefix)], wantPrefix) {
		t.Fatalf("unexpected disposable catalogue ensure invocation\nwant prefix: %#v + session\ngot:         %#v", wantPrefix, got)
	}
	sessionName := got[len(got)-1]
	if !strings.HasPrefix(sessionName, "roundfix-preflight-") {
		t.Fatalf("expected disposable preflight session name, got %q", sessionName)
	}
	return sessionName
}

func assertDisposableShowInvocation(t *testing.T, got []string, workDir string, runtime string, sessionName string) {
	t.Helper()
	want := []string{"--cwd", workDir, "--format", "json", "--json-strict", runtime, "sessions", "show", sessionName}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected disposable session show invocation\nwant: %#v\ngot:  %#v", want, got)
	}
}

func disposableSessionFromEnsure(t *testing.T, got []string) string {
	t.Helper()
	if len(got) < 2 || got[len(got)-2] != "--name" {
		t.Fatalf("expected sessions ensure --name invocation, got %#v", got)
	}
	sessionName := got[len(got)-1]
	if !strings.HasPrefix(sessionName, "roundfix-preflight-") {
		t.Fatalf("expected disposable preflight session name, got %q", sessionName)
	}
	return sessionName
}

func containsCommandKey(invocations [][]string, key string) bool {
	for _, invocation := range invocations {
		if fakeACPXCommandKey(invocation) == key {
			return true
		}
	}
	return false
}

func containsInvocationValue(invocations [][]string, flag string, value string) bool {
	for _, invocation := range invocations {
		for index := 0; index+1 < len(invocation); index++ {
			if invocation[index] == flag && invocation[index+1] == value {
				return true
			}
		}
	}
	return false
}

func containsCommandValue(invocations [][]string, command string, value string) bool {
	for _, invocation := range invocations {
		if fakeACPXCommandKey(invocation) == command && fakeACPXCallKey(invocation) == command+" value="+value {
			return true
		}
	}
	return false
}

func assertFallbackCandidateSessions(t *testing.T, invocations [][]string, wantModels []string) {
	t.Helper()
	models := make([]string, 0, len(wantModels))
	sessions := map[string]int{}
	closed := map[string]int{}
	for _, invocation := range invocations {
		switch fakeACPXCommandKey(invocation) {
		case "sessions ensure":
			models = append(models, invocationValue(invocation, "--model"))
			sessions[disposableSessionFromEnsure(t, invocation)]++
		case "sessions close":
			closed[invocation[len(invocation)-1]]++
		}
	}
	if !reflect.DeepEqual(models, wantModels) {
		t.Fatalf("unexpected candidate order\nwant: %#v\ngot:  %#v", wantModels, models)
	}
	if len(sessions) != len(wantModels) {
		t.Fatalf("expected one unique disposable session per candidate, got %#v", sessions)
	}
	for session, count := range sessions {
		if count != 1 || closed[session] != 1 {
			t.Fatalf("expected disposable session %q to be ensured and closed once, ensured=%d closed=%d", session, count, closed[session])
		}
	}
	if containsCommandKey(invocations, "prompt") {
		t.Fatalf("fallback selection probes must not send prompts: %#v", invocations)
	}
}

func invocationValue(args []string, flag string) string {
	for index := 0; index+1 < len(args); index++ {
		if args[index] == flag {
			return args[index+1]
		}
	}
	return ""
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

type deadlineRecordingCodexProbe struct {
	deadlineSeen []bool
}

func (probe *deadlineRecordingCodexProbe) Quarantined(ctx context.Context, _ string) (bool, error) {
	_, ok := ctx.Deadline()
	probe.deadlineSeen = append(probe.deadlineSeen, ok)
	return false, nil
}

func (probe *deadlineRecordingCodexProbe) Accepted(ctx context.Context, _ string) (bool, error) {
	_, ok := ctx.Deadline()
	probe.deadlineSeen = append(probe.deadlineSeen, ok)
	return true, nil
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

// agentWaitBudget bounds every "this eventually happened" wait in this
// package: file milestones a spawned fake ACPX writes, and results a
// goroutine sends back. It is deliberately far longer than the work needs —
// a passing wait returns the moment its condition holds and pays none of it,
// while a stuck one still fails.
//
// These tests spawn real child processes, and since the ADR-0089 environment
// seam they run beside the package's other parallel tests. Under a saturated
// machine a five-second budget measures load rather than correctness: two
// cancellation tests failed exactly that way inside the full suite while
// passing alone.
//
// Behavioral durations are not wait budgets and are untouched: StopGrace
// still governs how long a stop waits before escalating, and the fake
// cancellation clock keeps its own scale.
const agentWaitBudget = 90 * time.Second

func waitForFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.After(agentWaitBudget)
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

func waitForACPXMilestone(t *testing.T, name string, path string, invocationsPath string) {
	t.Helper()
	deadline := time.After(agentWaitBudget)
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		select {
		case <-deadline:
			invocations := "unavailable"
			if content, err := os.ReadFile(invocationsPath); err == nil {
				invocations = string(content)
			}
			t.Fatalf("timed out waiting for fake acpx %s milestone at %s; invocations so far:\n%s", name, path, invocations)
		case <-ticker.C:
		}
	}
}

func receiveError(t *testing.T, resultCh <-chan error) error {
	t.Helper()
	select {
	case err := <-resultCh:
		return err
	case <-time.After(agentWaitBudget):
		t.Fatal("timed out waiting for acpx run")
		return nil
	}
}

type testDeadlineContext struct {
	context.Context
	done chan struct{}
	once sync.Once
}

func newTestDeadlineContext() *testDeadlineContext {
	return &testDeadlineContext{
		Context: context.Background(),
		done:    make(chan struct{}),
	}
}

func (ctx *testDeadlineContext) Done() <-chan struct{} {
	return ctx.done
}

func (ctx *testDeadlineContext) Err() error {
	select {
	case <-ctx.done:
		return context.DeadlineExceeded
	default:
		return nil
	}
}

func (ctx *testDeadlineContext) expire() {
	ctx.once.Do(func() {
		close(ctx.done)
	})
}

func mustJSONForTest[T any](t *testing.T, value T) string {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal test JSON: %v", err)
	}
	return string(payload)
}

func selectionStateFixture(t *testing.T, configID string, value string, currentModel string, models []string, reasoningID string, currentReasoning string, efforts []string) string {
	t.Helper()
	modelValues := make([]map[string]string, 0, len(models))
	for _, model := range models {
		modelValues = append(modelValues, map[string]string{"value": model})
	}
	options := []map[string]any{{
		"id":           "model",
		"category":     "model",
		"type":         "select",
		"currentValue": currentModel,
		"options":      modelValues,
	}}
	if reasoningID != "" {
		reasoningValues := make([]map[string]string, 0, len(efforts))
		for _, effort := range efforts {
			reasoningValues = append(reasoningValues, map[string]string{"value": effort})
		}
		options = append(options, map[string]any{
			"id":           reasoningID,
			"type":         "select",
			"currentValue": currentReasoning,
			"options":      reasoningValues,
		})
	}
	return mustJSONForTest(t, map[string]any{
		"action":        "config_set",
		"configId":      configID,
		"value":         value,
		"configOptions": options,
	})
}

func sessionCapabilitySnapshotFixture(t *testing.T, currentModel string, models []string, reasoningID string, currentReasoning string, efforts []string) string {
	t.Helper()
	var response acpxCapabilityResponse
	if err := json.Unmarshal([]byte(selectionStateFixture(t, "model", currentModel, currentModel, models, reasoningID, currentReasoning, efforts)), &response); err != nil {
		t.Fatalf("decode capability response fixture: %v", err)
	}
	return mustJSONForTest(t, map[string]any{
		"schema": "acpx.session.v1",
		"acpx": map[string]any{
			"current_model_id": currentModel,
			"config_options":   response.ConfigOptions,
		},
	})
}

func selectionCallKeys(invocations [][]string) []string {
	keys := make([]string, 0, len(invocations))
	for _, invocation := range invocations {
		switch fakeACPXCommandKey(invocation) {
		case "sessions ensure", "set model", "set reasoning_effort", "set effort":
			keys = append(keys, fakeACPXCallKey(invocation))
		case "sessions close":
			keys = append(keys, "sessions close")
		}
	}
	return keys
}

func exactProofCallKeys(invocations [][]string) []string {
	keys := make([]string, 0, len(invocations))
	for _, invocation := range invocations {
		switch fakeACPXCommandKey(invocation) {
		case "sessions ensure", "set model", "set reasoning_effort", "set effort":
			keys = append(keys, fakeACPXCallKey(invocation))
		case "sessions show", "sessions close":
			keys = append(keys, fakeACPXCommandKey(invocation))
		}
	}
	return keys
}

func assertLastInvocationClosesDisposable(t *testing.T, harness *fakeACPXHarness) {
	t.Helper()
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if len(invocations) == 0 || fakeACPXCommandKey(invocations[len(invocations)-1]) != "sessions close" {
		t.Fatalf("last invocation did not close disposable Agent Session: %#v", invocations)
	}
}

func exactSelectionInvocations(legacy [][]string) [][]string {
	updated := make([][]string, 0, len(legacy)+2)
	for _, invocation := range legacy {
		command := fakeACPXCommandKey(invocation)
		if command == "set reasoning_effort" || command == "set effort" {
			strict := []string{"--cwd", invocationValue(invocation, "--cwd"), "--format", "json", "--json-strict"}
			strict = append(strict, invocation[2:]...)
			updated = append(updated, strict)
			continue
		}
		updated = append(updated, invocation)
		if command != "sessions ensure" {
			continue
		}
		agentArgs := []string{}
		if commandOverride := invocationValue(invocation, "--agent"); commandOverride != "" {
			agentArgs = []string{"--agent", commandOverride}
		} else {
			for index, value := range invocation {
				if value == "sessions" && index > 0 {
					agentArgs = []string{invocation[index-1]}
					break
				}
			}
		}
		sessionShow := []string{"--cwd", invocationValue(invocation, "--cwd"), "--format", "json", "--json-strict"}
		sessionShow = append(sessionShow, agentArgs...)
		sessionShow = append(sessionShow, "sessions", "show", invocationValue(invocation, "--name"))
		updated = append(updated, sessionShow)
	}
	return updated
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
			fmt.Fprintf(os.Stderr, "marshal args: %v\n", err)
			return 2
		}
		if err := os.WriteFile(path, payload, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write args: %v\n", err)
			return 2
		}
	}
	if path := os.Getenv(fakeACPXInvokes); path != "" {
		if err := appendFakeACPXInvocation(path, args); err != nil {
			fmt.Fprintf(os.Stderr, "append invocation: %v\n", err)
			return 2
		}
	}
	if path := os.Getenv(fakeACPXCodexPath); path != "" {
		if err := appendFakeACPXString(path, os.Getenv(codexPathEnv)); err != nil {
			fmt.Fprintf(os.Stderr, "append codex path: %v\n", err)
			return 2
		}
	}
	promptPath := os.Getenv(fakeACPXPromptPath)
	promptsPath := os.Getenv(fakeACPXPrompts)
	if commandKey == "prompt" && (promptPath != "" || promptsPath != "") {
		prompt, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read stdin: %v\n", err)
			return 2
		}
		if promptPath != "" {
			if err := os.WriteFile(promptPath, prompt, 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "write prompt: %v\n", err)
				return 2
			}
		}
		if promptsPath != "" {
			if err := appendFakeACPXString(promptsPath, string(prompt)); err != nil {
				fmt.Fprintf(os.Stderr, "append prompt: %v\n", err)
				return 2
			}
		}
	}
	if commandKey == "cancel" {
		if path := os.Getenv(fakeACPXCanceled); path != "" {
			if err := os.WriteFile(path, []byte("canceled\n"), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "write cancel marker: %v\n", err)
				return 2
			}
		}
	}
	if commandKey == "sessions close" {
		if path := os.Getenv(fakeACPXClosed); path != "" {
			if err := os.WriteFile(path, []byte("closed\n"), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "write close marker: %v\n", err)
				return 2
			}
		}
	}
	if blockCommand := os.Getenv(fakeACPXBlockCmd); blockCommand != "" && commandKey == blockCommand {
		if path := os.Getenv(fakeACPXStarted); path != "" {
			if err := os.WriteFile(path, []byte("started\n"), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "write started marker: %v\n", err)
				return 2
			}
		}
		for {
			time.Sleep(5 * time.Millisecond)
		}
	}
	if commandKey == "prompt" && os.Getenv(fakeACPXBlock) == "1" {
		if path := os.Getenv(fakeACPXStarted); path != "" {
			if err := os.WriteFile(path, []byte("started\n"), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "write started marker: %v\n", err)
				return 2
			}
		}
		if summary := os.Getenv(fakeACPXStartEvent); summary != "" {
			update := `{"sessionId":"fake","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":` + strconv.Quote(summary) + `}}}`
			if _, err := io.WriteString(os.Stdout, acpxUpdateLine(update)); err != nil {
				fmt.Fprintf(os.Stderr, "write prompt-started Run Event: %v\n", err)
				return 2
			}
		}
		for {
			if canceled := os.Getenv(fakeACPXCanceled); canceled != "" {
				if _, err := os.Stat(canceled); err == nil && os.Getenv(fakeACPXExitCancel) == "1" {
					if err := writeFakeACPXPromptCompletion(); err != nil {
						fmt.Fprintf(os.Stderr, "write prompt completion marker: %v\n", err)
						return 2
					}
					return 130
				}
			}
			if closed := os.Getenv(fakeACPXClosed); closed != "" {
				if _, err := os.Stat(closed); err == nil {
					if err := writeFakeACPXPromptCompletion(); err != nil {
						fmt.Fprintf(os.Stderr, "write prompt completion marker: %v\n", err)
						return 2
					}
					return 130
				}
			}
			time.Sleep(5 * time.Millisecond)
		}
	}
	stdoutByCommand := fakeACPXStringMap(os.Getenv(fakeACPXStdoutBy))
	stdoutByCall := fakeACPXStringMap(os.Getenv(fakeACPXStdoutCall))
	stderrByCommand := fakeACPXStringMap(os.Getenv(fakeACPXStderrBy))
	if commandKey == "prompt" {
		if err := writeFakeACPXThoughtStream(os.Getenv(fakeACPXThoughtLen)); err != nil {
			fmt.Fprintf(os.Stderr, "write thought stream: %v\n", err)
			return 2
		}
	}
	_, _ = io.WriteString(os.Stdout, firstFakeACPXString(stdoutByCall[fakeACPXCallKey(args)], stdoutByCommand[commandKey], fakeSelectionStateOutput(args, os.Getenv(fakeACPXInvokes)), os.Getenv(fakeACPXStdout)))
	_, _ = io.WriteString(os.Stderr, firstFakeACPXString(stderrByCommand[commandKey], os.Getenv(fakeACPXStderr)))
	exitByCommand := fakeACPXIntMap(os.Getenv(fakeACPXExitBy))
	exitByCall := fakeACPXIntMap(os.Getenv(fakeACPXExitByCall))
	if exitCode, ok := exitByCall[fakeACPXCallKey(args)]; ok {
		return exitCode
	}
	if exitCode, ok := exitByCommand[commandKey]; ok {
		return exitCode
	}
	if rawExitCode := os.Getenv(fakeACPXExitCode); rawExitCode != "" {
		exitCode, err := strconv.Atoi(rawExitCode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse exit code: %v\n", err)
			return 2
		}
		return exitCode
	}
	return 0
}

func writeFakeACPXThoughtStream(rawLength string) error {
	if rawLength == "" {
		return nil
	}
	length, err := strconv.Atoi(rawLength)
	if err != nil {
		return fmt.Errorf("parse thought length: %w", err)
	}
	const chunkSize = 1024
	for remaining := length; remaining > 0; {
		size := min(remaining, chunkSize)
		update := `{"sessionId":"sealed","update":{"sessionUpdate":"agent_thought_chunk","content":{"type":"text","text":"` +
			strings.Repeat("x", size) +
			`"}}}`
		if _, err := io.WriteString(os.Stdout, acpxUpdateLine(update)); err != nil {
			return err
		}
		remaining -= size
	}
	return nil
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

func writeFakeACPXPromptCompletion() error {
	path := os.Getenv(fakeACPXPromptDone)
	if path == "" {
		return nil
	}
	return os.WriteFile(path, []byte("completed\n"), 0o644)
}

// fakeACPXCommandKey mirrors the real acpx grammar: program-level global
// options come first (--cwd, --format, --model, --agent, ... take a value;
// permission and terminal controls are booleans), then the agent name (unless
// --agent supplied the raw command), then the subcommand.
func fakeACPXCommandKey(args []string) string {
	valueGlobals := map[string]bool{
		"--cwd":                         true,
		"--format":                      true,
		"--model":                       true,
		"--agent":                       true,
		"--timeout":                     true,
		"--ttl":                         true,
		"--non-interactive-permissions": true,
		"--allowed-tools":               true,
		"--max-turns":                   true,
		"--prompt-retries":              true,
	}
	booleanGlobals := map[string]bool{
		"--json-strict": true,
		"--approve-all": true,
		"--deny-all":    true,
		"--no-terminal": true,
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

func fakeACPXCallKey(args []string) string {
	command := fakeACPXCommandKey(args)
	switch command {
	case "sessions ensure":
		return command + " model=" + invocationValue(args, "--model")
	case "set model", "set reasoning_effort", "set effort":
		for index := 0; index+2 < len(args); index++ {
			if args[index] == "set" {
				return command + " value=" + args[index+2]
			}
		}
	}
	return command
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

func fakeSelectionStateOutput(args []string, invocationsPath string) string {
	command := fakeACPXCommandKey(args)
	if !containsArg(args, "--json-strict") {
		return ""
	}
	if command == "sessions show" {
		model := fakeSessionModel(invocationsPath, args[len(args)-1])
		if model == "" {
			return ""
		}
		reasoningID := "reasoning_effort"
		if strings.Contains(strings.Join(args, " "), " claude ") || strings.Contains(strings.Join(args, " "), " opencode ") {
			reasoningID = "effort"
		}
		response := map[string]any{
			"schema": "acpx.session.v1",
			"acpx": map[string]any{
				"current_model_id": model,
				"config_options": []any{
					map[string]any{"id": "model", "category": "model", "type": "select", "currentValue": model, "options": []any{map[string]string{"value": model}}},
					map[string]any{"id": reasoningID, "type": "select", "currentValue": "medium", "options": []any{
						map[string]string{"value": "low"}, map[string]string{"value": "medium"}, map[string]string{"value": "high"}, map[string]string{"value": "xhigh"}, map[string]string{"value": "maximum"}, map[string]string{"value": "max"}, map[string]string{"value": "ultra"},
					}},
				},
			},
		}
		payload, err := json.Marshal(response)
		if err != nil {
			return ""
		}
		return string(payload)
	}
	if command != "set model" && command != "set reasoning_effort" && command != "set effort" {
		return ""
	}
	configID := strings.TrimPrefix(command, "set ")
	value := ""
	for index := 0; index+2 < len(args); index++ {
		if args[index] == "set" {
			value = args[index+2]
			break
		}
	}
	model := value
	if configID != "model" {
		model = fakeSessionModel(invocationsPath, invocationValue(args, "-s"))
	}
	if model == "" || value == "" {
		return ""
	}
	reasoningID := configID
	currentReasoning := value
	if configID == "model" {
		reasoningID = "reasoning_effort"
		if strings.Contains(strings.Join(args, " "), " claude ") || strings.Contains(strings.Join(args, " "), " opencode ") {
			reasoningID = "effort"
		}
		currentReasoning = "medium"
	}
	response := map[string]any{
		"action":   "config_set",
		"configId": configID,
		"value":    value,
		"configOptions": []any{
			map[string]any{"id": "model", "category": "model", "type": "select", "currentValue": model, "options": []any{map[string]string{"value": model}}},
			map[string]any{"id": reasoningID, "type": "select", "currentValue": currentReasoning, "options": []any{
				map[string]string{"value": "low"}, map[string]string{"value": "medium"}, map[string]string{"value": "high"}, map[string]string{"value": "xhigh"}, map[string]string{"value": "maximum"}, map[string]string{"value": "max"}, map[string]string{"value": "ultra"},
			}},
		},
	}
	payload, err := json.Marshal(response)
	if err != nil {
		return ""
	}
	return string(payload)
}

func fakeSessionModel(invocationsPath string, sessionName string) string {
	content, err := os.ReadFile(invocationsPath)
	if err != nil {
		return ""
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	for index := len(lines) - 1; index >= 0; index-- {
		var args []string
		if json.Unmarshal([]byte(lines[index]), &args) != nil || fakeACPXCommandKey(args) != "sessions ensure" {
			continue
		}
		if invocationValue(args, "--name") == sessionName {
			return invocationValue(args, "--model")
		}
	}
	return ""
}

func containsArg(args []string, want string) bool {
	for _, arg := range args {
		if arg == want {
			return true
		}
	}
	return false
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

// TestReasoningEffortConfigKeyMapsSupportedRuntimes keeps the runtime-specific
// Codex key separate while Claude and OpenCode use the generic effort key.
func TestReasoningEffortConfigKeyMapsSupportedRuntimes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		runtime RuntimeSpec
		wantKey string
	}{
		{name: "codex maps to its own key", runtime: RuntimeSpec{ID: "codex"}, wantKey: acpxCodexReasoningEffortKey},
		{name: "claude maps to the generic key", runtime: RuntimeSpec{ID: "claude"}, wantKey: acpxGenericReasoningEffortKey},
		{name: "opencode maps to the generic key", runtime: RuntimeSpec{ID: "opencode"}, wantKey: acpxGenericReasoningEffortKey},
		{name: "opencode override maps to the generic key", runtime: RuntimeSpec{ID: "opencode-custom"}, wantKey: acpxGenericReasoningEffortKey},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			key, err := acpxReasoningEffortConfigKey(tt.runtime)
			if err != nil {
				t.Fatalf("reasoning key for %q: %v", tt.runtime.ID, err)
			}
			if key != tt.wantKey {
				t.Fatalf("reasoning key = %q, want %q", key, tt.wantKey)
			}
		})
	}
}

// TestValidateRuntimeSelectionAcceptsOpenCodeReasoningEffort closes the
// invocation-override path: a configured effort can reach OpenCode without
// bypassing runtime validation.
func TestValidateRuntimeSelectionAcceptsOpenCodeReasoningEffort(t *testing.T) {
	t.Parallel()

	if err := validateRuntimeSelection(RuntimeSpec{ID: "opencode", Model: "opencode-go/kimi-k3", ReasoningEffort: "max"}); err != nil {
		t.Fatalf("validate OpenCode selection with a non-empty reasoning effort: %v", err)
	}

	if err := validateRuntimeSelection(RuntimeSpec{ID: "opencode", Model: "opencode-go/kimi-k3"}); err != nil {
		t.Fatalf("a model-managed OpenCode selection must validate: %v", err)
	}
}

// TestRuntimeManagedReasoningProvesAgainstAnAdvertisedEffortOption is the
// measured OpenCode shape: the adapter advertises a per-model effort option
// that Roundfix declines to assign, and the selection still proves.
func TestRuntimeManagedReasoningProvesAgainstAnAdvertisedEffortOption(t *testing.T) {
	t.Parallel()

	fixture := `{"action":"config_set","configId":"model","value":"opencode-go/kimi-k3","configOptions":[` +
		`{"id":"model","category":"model","type":"select","currentValue":"opencode-go/kimi-k3","options":[{"value":"opencode-go/kimi-k3"},{"value":"opencode-go/qwen3.8-max"}]},` +
		`{"id":"effort","type":"select","currentValue":"max","options":[{"value":"max"}]}]}`

	capabilities, err := ParseSessionConfigOptions(
		[]byte(fixture),
		AdapterEvidence{Command: "opencode"},
		SelectionRetention{Model: "opencode-go/kimi-k3"},
	)
	if err != nil {
		t.Fatalf("project the OpenCode capability shape: %v", err)
	}

	runtime := RuntimeSpec{ID: "opencode", Model: "opencode-go/kimi-k3"}
	assignment, err := PlanSelectionAssignment(runtime, capabilities)
	if err != nil {
		t.Fatalf("plan a model-managed OpenCode selection: %v", err)
	}
	if assignment.Encoding != SelectionEncodingRuntimeManaged {
		t.Fatalf("encoding = %q, want %q", assignment.Encoding, SelectionEncodingRuntimeManaged)
	}
	if !selectionStateMatches(assignment, capabilities) {
		t.Fatalf("an advertised effort option Roundfix never assigns must not break the proof: %#v", capabilities.ReasoningOption)
	}

	// A runtime Roundfix does control keeps the strict rule.
	claude := RuntimeSpec{ID: "claude", Model: "opencode-go/kimi-k3"}
	claudeAssignment, err := PlanSelectionAssignment(claude, capabilities)
	if err != nil {
		t.Fatalf("plan a model-managed Claude selection: %v", err)
	}
	if claudeAssignment.Encoding != SelectionEncodingModelManaged {
		t.Fatalf("encoding = %q, want %q", claudeAssignment.Encoding, SelectionEncodingModelManaged)
	}
	if selectionStateMatches(claudeAssignment, capabilities) {
		t.Fatal("model_managed must still require the absence of a reasoning option")
	}
}
