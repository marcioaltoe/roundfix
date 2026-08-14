package daemon

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"testing"
)

// Suite: shared Verification prober
// Invariant: every command produces one ordered verdict that distinguishes an observed exit from an unobserved command.
// Boundary IN: command classification and per-command verifier requests.
// Boundary OUT: real shell execution and Daemon Run bookkeeping, covered by daemon_test.go and task_engine_test.go.
func TestProbeCommands(t *testing.T) {
	t.Parallel()

	t.Run("classifies ordered command outcomes", func(t *testing.T) {
		const (
			vacuousCommand = "gate already passes"
			honestCommand  = "gate rejects unchanged tree"
			unknownCommand = "gate could not run"
		)
		unknownCause := errors.New("runner unavailable")
		verifier := &probeScriptVerifier{errorsByCommand: map[string]error{
			honestCommand:  &VerificationCommandError{Err: errors.New("exit status 1")},
			unknownCommand: &VerificationUnknownError{Err: unknownCause},
		}}
		commands := []string{vacuousCommand, honestCommand, unknownCommand}
		workDir := t.TempDir()
		outputDir := t.TempDir()

		verdicts, err := ProbeCommands(context.Background(), verifier, workDir, commands, func(index int) string {
			return filepath.Join(outputDir, fmt.Sprintf("command-%d.log", index+1))
		})

		if err != nil {
			t.Fatalf("ProbeCommands() error = %v", err)
		}
		if len(verdicts) != len(commands) {
			t.Fatalf("ProbeCommands() verdict count = %d, want %d", len(verdicts), len(commands))
		}
		for index, command := range commands {
			if verdicts[index].Command != command {
				t.Fatalf("ProbeCommands() verdict %d command = %q, want %q", index, verdicts[index].Command, command)
			}
		}
		if !verdicts[0].Vacuous || verdicts[0].Unknown || verdicts[0].Cause != nil {
			t.Fatalf("passing command verdict = %+v, want vacuous", verdicts[0])
		}
		if verdicts[1].Vacuous || verdicts[1].Unknown || verdicts[1].Cause != nil {
			t.Fatalf("non-zero command verdict = %+v, want honest", verdicts[1])
		}
		if verdicts[2].Vacuous || !verdicts[2].Unknown || !errors.Is(verdicts[2].Cause, unknownCause) {
			t.Fatalf("unobserved command verdict = %+v, want unknown cause %v", verdicts[2], unknownCause)
		}
		var unknownErr *VerificationUnknownError
		if !errors.As(verdicts[2].Cause, &unknownErr) {
			t.Fatalf("unobserved command cause type = %T, want *VerificationUnknownError", verdicts[2].Cause)
		}
		if unknownErr.Command != unknownCommand || unknownErr.DiagnosticPath != filepath.Join(outputDir, "command-3.log") {
			t.Fatalf("unobserved command cause = %+v, want completed command and diagnostic path", unknownErr)
		}

		gotCommands := make([]string, 0, len(verifier.requests))
		gotOutputPaths := make([]string, 0, len(verifier.requests))
		for _, request := range verifier.requests {
			if request.WorkDir != workDir {
				t.Fatalf("verifier WorkDir = %q, want %q", request.WorkDir, workDir)
			}
			gotCommands = append(gotCommands, request.Command)
			gotOutputPaths = append(gotOutputPaths, request.OutputPath)
		}
		if !slices.Equal(gotCommands, commands) {
			t.Fatalf("verifier commands = %q, want %q", gotCommands, commands)
		}
		wantOutputPaths := []string{
			filepath.Join(outputDir, "command-1.log"),
			filepath.Join(outputDir, "command-2.log"),
			filepath.Join(outputDir, "command-3.log"),
		}
		if !slices.Equal(gotOutputPaths, wantOutputPaths) {
			t.Fatalf("verifier output paths = %q, want %q", gotOutputPaths, wantOutputPaths)
		}
	})

	t.Run("returns no verdicts for an empty command list", func(t *testing.T) {
		verifier := &probeScriptVerifier{}
		outputCalls := 0

		verdicts, err := ProbeCommands(context.Background(), verifier, t.TempDir(), nil, func(int) string {
			outputCalls++
			return "unused"
		})

		if err != nil {
			t.Fatalf("ProbeCommands() error = %v", err)
		}
		if len(verdicts) != 0 {
			t.Fatalf("ProbeCommands() verdicts = %+v, want none", verdicts)
		}
		if len(verifier.requests) != 0 || outputCalls != 0 {
			t.Fatalf("empty command list called verifier %d times and outputFor %d times", len(verifier.requests), outputCalls)
		}
	})
}

type probeScriptVerifier struct {
	errorsByCommand map[string]error
	requests        []VerifyRequest
}

func (verifier *probeScriptVerifier) Verify(_ context.Context, request VerifyRequest) (VerifyResult, error) {
	verifier.requests = append(verifier.requests, request)
	return VerifyResult{OutputPath: request.OutputPath}, verifier.errorsByCommand[request.Command]
}
