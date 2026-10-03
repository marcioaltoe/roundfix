package daemon

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// ErrVerificationMalformed identifies a command rejected by the shell parser.
var ErrVerificationMalformed = errors.New("verification command malformed")

// CommandVerdict records what one authored command does against a tree where
// no work has happened. Unknown means the command's verdict was not observed.
type CommandVerdict struct {
	Command string
	Vacuous bool
	Unknown bool
	Cause   error
}

type commandProbeError struct {
	command string
	err     error
}

func (err *commandProbeError) Error() string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("probe command %q: %v", err.command, err.err)
}

func (err *commandProbeError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.err
}

// ProbeCommands runs commands in order and classifies their behavior against
// workDir. It owns no Run state or Verification Capacity; outputFor selects
// only the diagnostic destination for each zero-based command index.
func ProbeCommands(
	ctx context.Context,
	verifier Verifier,
	workDir string,
	commands []string,
	outputFor func(index int) string,
) ([]CommandVerdict, error) {
	verdicts := make([]CommandVerdict, 0, len(commands))
	for index, command := range commands {
		outputPath := outputFor(index)
		parser := exec.CommandContext(ctx, "sh", "-n", "-c", command)
		parser.Dir = workDir
		diagnostic, parseErr := parser.CombinedOutput()
		if ctx.Err() != nil {
			return nil, &commandProbeError{command: command, err: ctx.Err()}
		}
		var exitErr *exec.ExitError
		if errors.As(parseErr, &exitErr) {
			verdicts = append(verdicts, CommandVerdict{
				Command: command,
				Unknown: true,
				Cause: completeVerificationUnknownCause(&VerificationUnknownError{
					Err: fmt.Errorf("%w: %s", ErrVerificationMalformed, strings.TrimSpace(string(diagnostic))),
				}, command, outputPath),
			})
			continue
		}
		// A parser that cannot start must not introduce a new refusal.
		result, verifyErr := verifier.Verify(ctx, VerifyRequest{
			WorkDir:    workDir,
			Command:    command,
			OutputPath: outputPath,
		})
		verdict := CommandVerdict{Command: command}
		if verifyErr == nil {
			verdict.Vacuous = true
			verdicts = append(verdicts, verdict)
			continue
		}
		if verificationStopRequested(ctx, verifyErr) {
			return nil, &commandProbeError{command: command, err: verifyErr}
		}
		var commandErr *VerificationCommandError
		if errors.As(verifyErr, &commandErr) {
			if exitCode, couldNotExecute := verificationCommandCouldNotExecute(commandErr); couldNotExecute {
				verdict.Unknown = true
				verdict.Cause = completeVerificationUnknownCause(&VerificationUnknownError{
					Err: fmt.Errorf("command could not be executed (exit %d): %w", exitCode, commandErr.Err),
				}, command, result.OutputPath)
				verdicts = append(verdicts, verdict)
				continue
			}
			verdicts = append(verdicts, verdict)
			continue
		}
		var unknownErr *VerificationUnknownError
		if errors.As(verifyErr, &unknownErr) {
			verdict.Unknown = true
			verdict.Cause = completeVerificationUnknownCause(unknownErr, command, result.OutputPath)
			verdicts = append(verdicts, verdict)
			continue
		}
		return nil, &commandProbeError{command: command, err: verifyErr}
	}
	return verdicts, nil
}

func verificationCommandCouldNotExecute(err error) (int, bool) {
	var exitErr interface{ ExitCode() int }
	if !errors.As(err, &exitErr) {
		return 0, false
	}
	exitCode := exitErr.ExitCode()
	return exitCode, exitCode == 126 || exitCode == 127
}
