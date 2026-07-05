package agent

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"roundfix/internal/runevent"

	acp "github.com/coder/acp-go-sdk"
)

const (
	defaultACPXCommand              = "acpx"
	acpxPermissionDeniedStatus      = "permissions_denied"
	acpxExitReasonAgentProtocol     = "agent/protocol error"
	acpxExitReasonTimeout           = "timeout"
	acpxExitReasonPermissionsDenied = "all permissions denied"
	acpxExitReasonUsage             = "usage error"
	acpxExitReasonMissingSession    = "missing session"
)

// ACPXRunner is the acpx-backed invocation core. Later migration tasks wire
// this into Runner after Agent Session lifecycle is available.
type ACPXRunner struct {
	Command string
	Now     func() time.Time
}

// ACPXPromptRequest carries the explicit Agent Session needed by acpx prompt
// invocations while leaving the existing Runner interface unchanged.
type ACPXPromptRequest struct {
	ExecuteRequest
	Session string
}

// BatchFailureError is returned when acpx reports an Agent/Batch-level
// failure. Callers can settle the Batch and continue under ADR 0010.
type BatchFailureError struct {
	ExitCode int
	Reason   string
	Stderr   string
}

func (err *BatchFailureError) Error() string {
	if err == nil {
		return ""
	}
	message := "Agent Batch failed"
	if err.ExitCode != 0 {
		message = fmt.Sprintf("%s after acpx exited with code %d", message, err.ExitCode)
	}
	if err.Reason != "" {
		message += ": " + err.Reason
	}
	if detail := strings.TrimSpace(err.Stderr); detail != "" {
		message += ": " + detail
	}
	return message
}

// InfrastructureError is returned when acpx reports a Roundfix/environment
// bug rather than a user-resolvable Batch failure.
type InfrastructureError struct {
	ExitCode int
	Reason   string
	Stderr   string
}

func (err *InfrastructureError) Error() string {
	if err == nil {
		return ""
	}
	message := "acpx infrastructure error"
	if err.ExitCode != 0 {
		message = fmt.Sprintf("%s after exit code %d", message, err.ExitCode)
	}
	if err.Reason != "" {
		message += ": " + err.Reason
	}
	if detail := strings.TrimSpace(err.Stderr); detail != "" {
		message += ": " + detail
	}
	return message
}

type acpxJSONRPCMessage struct {
	Method string            `json:"method"`
	Params json.RawMessage   `json:"params"`
	Result json.RawMessage   `json:"result"`
	Error  *acpxJSONRPCError `json:"error"`
}

type acpxJSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (runner ACPXRunner) RunPrompt(ctx context.Context, req ACPXPromptRequest, sink runevent.Sink) (ExecuteResult, error) {
	result := ExecuteResult{LogPath: req.LogPath}
	if err := validateACPXPromptRequest(req); err != nil {
		return result, err
	}
	if sink == nil {
		sink = runevent.Discard
	}
	if err := os.MkdirAll(filepath.Dir(req.LogPath), 0o755); err != nil {
		return result, fmt.Errorf("create Agent log directory: %w", err)
	}
	logFile, err := os.OpenFile(req.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return result, fmt.Errorf("open Agent log %q: %w", req.LogPath, err)
	}
	defer func() {
		_ = logFile.Close()
	}()

	if err := ctx.Err(); err != nil {
		return result, StopError{LogPath: req.LogPath, Err: err}
	}
	args, err := acpxPromptArgs(req)
	if err != nil {
		return result, err
	}
	cmd := exec.CommandContext(ctx, runner.command(), args...)
	cmd.Dir = req.GitRoot
	cmd.Stdin = strings.NewReader(req.Prompt)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return result, fmt.Errorf("open acpx stdout: %w", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return result, fmt.Errorf("start acpx prompt: %w", err)
	}

	var output bytes.Buffer
	var streamErr error
	reader := bufio.NewReader(stdout)
	for {
		line, readErr := reader.ReadBytes('\n')
		if len(line) > 0 {
			if _, err := logFile.Write(line); err != nil && streamErr == nil {
				streamErr = fmt.Errorf("write Agent log: %w", err)
			}
			if _, err := output.Write(line); err != nil && streamErr == nil {
				streamErr = fmt.Errorf("capture acpx stdout: %w", err)
			}
			if err := runner.handleStdoutLine(ctx, req.ExecuteRequest, sink, line, &result.StopReason); err != nil && streamErr == nil {
				streamErr = err
			}
		}
		if readErr != nil {
			if !errors.Is(readErr, io.EOF) && streamErr == nil {
				streamErr = fmt.Errorf("read acpx stdout: %w", readErr)
			}
			break
		}
	}
	waitErr := cmd.Wait()
	result.Output = output.String()
	if waitErr != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return result, StopError{LogPath: req.LogPath, Output: result.Output, Err: ctxErr}
		}
		exitCode, ok := commandExitCode(waitErr)
		if !ok {
			return result, fmt.Errorf("wait for acpx prompt: %w", waitErr)
		}
		return result, runner.mapExitCode(ctx, req.ExecuteRequest, sink, exitCode, stderr.String(), result.Output)
	}
	if streamErr != nil {
		return result, &BatchFailureError{Reason: streamErr.Error(), Stderr: stderr.String()}
	}
	if result.StopReason == "" {
		return result, &BatchFailureError{Reason: "missing session/prompt stop reason", Stderr: stderr.String()}
	}
	return result, nil
}

func validateACPXPromptRequest(req ACPXPromptRequest) error {
	if strings.TrimSpace(req.LogPath) == "" {
		return errors.New("Agent log path is required")
	}
	if strings.TrimSpace(req.Session) == "" {
		return errors.New("Agent Session is required")
	}
	if strings.TrimSpace(req.GitRoot) == "" {
		return errors.New("Agent working directory is required")
	}
	return nil
}

func (runner ACPXRunner) command() string {
	if strings.TrimSpace(runner.Command) != "" {
		return strings.TrimSpace(runner.Command)
	}
	return defaultACPXCommand
}

func acpxPromptArgs(req ACPXPromptRequest) ([]string, error) {
	agentArgs, err := acpxAgentArgs(req.Runtime)
	if err != nil {
		return nil, err
	}
	args := append([]string{}, agentArgs...)
	args = append(args,
		"prompt",
		"-s", strings.TrimSpace(req.Session),
		"--cwd", strings.TrimSpace(req.GitRoot),
		"--format", "json",
		"--json-strict",
		"--approve-all",
	)
	if model := strings.TrimSpace(req.Runtime.Model); model != "" {
		args = append(args, "--model", model)
	}
	args = append(args, "-f", "-")
	return args, nil
}

func acpxAgentArgs(runtime RuntimeSpec) ([]string, error) {
	if runtime.Protocol == ProtocolStdio {
		command := strings.TrimSpace(runtime.Command)
		if command == "" {
			return nil, errors.New("Agent command override is required")
		}
		return []string{"--agent", command}, nil
	}
	agent := strings.TrimSpace(runtime.ID)
	if agent == "" {
		return nil, errors.New("ACP Runtime id is required")
	}
	return []string{agent}, nil
}

func (runner ACPXRunner) handleStdoutLine(ctx context.Context, req ExecuteRequest, sink runevent.Sink, line []byte, stopReason *string) error {
	var message acpxJSONRPCMessage
	if err := json.Unmarshal(line, &message); err != nil {
		return fmt.Errorf("parse acpx stdout JSON-RPC line: %w", err)
	}
	if message.Error != nil {
		return fmt.Errorf("acpx JSON-RPC error %d: %s", message.Error.Code, message.Error.Message)
	}
	if message.Method == acp.ClientMethodSessionUpdate {
		update, ok, err := streamUpdateFromSessionUpdatePayload(line)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		event := newAgentRunEvent(req, update, json.RawMessage(append([]byte(nil), line...)), eventClock(runner.Now)())
		if err := sink.Publish(ctx, event); err != nil {
			return fmt.Errorf("publish Run Events: %w", err)
		}
		return nil
	}
	if len(message.Result) == 0 {
		return nil
	}
	var response struct {
		StopReason string `json:"stopReason"`
	}
	if err := json.Unmarshal(message.Result, &response); err != nil {
		return fmt.Errorf("parse acpx session/prompt response: %w", err)
	}
	if response.StopReason != "" {
		*stopReason = response.StopReason
	}
	return nil
}

func streamUpdateFromSessionUpdatePayload(payload json.RawMessage) (StreamUpdate, bool, error) {
	note, ok, err := sessionNotificationFromPayload(payload)
	if err != nil || !ok {
		return StreamUpdate{}, false, err
	}
	update := streamUpdateFromACP(note.Update)
	if update.Kind == "" {
		return StreamUpdate{}, false, nil
	}
	return update, true, nil
}

func sessionNotificationFromPayload(payload json.RawMessage) (acp.SessionNotification, bool, error) {
	var note acp.SessionNotification
	if err := json.Unmarshal(payload, &note); err == nil {
		if update := streamUpdateFromACP(note.Update); update.Kind != "" {
			return note, true, nil
		}
	}
	var message acpxJSONRPCMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		return acp.SessionNotification{}, false, fmt.Errorf("parse acpx session/update JSON-RPC line: %w", err)
	}
	if message.Method != acp.ClientMethodSessionUpdate {
		return acp.SessionNotification{}, false, nil
	}
	if len(message.Params) == 0 {
		return acp.SessionNotification{}, false, errors.New("parse acpx session/update: params are required")
	}
	if err := json.Unmarshal(message.Params, &note); err != nil {
		return acp.SessionNotification{}, false, fmt.Errorf("parse acpx session/update params: %w", err)
	}
	return note, true, nil
}

func (runner ACPXRunner) mapExitCode(ctx context.Context, req ExecuteRequest, sink runevent.Sink, exitCode int, stderr string, output string) error {
	switch exitCode {
	case 0:
		return nil
	case 1:
		return &BatchFailureError{ExitCode: exitCode, Reason: acpxExitReasonAgentProtocol, Stderr: stderr}
	case 3:
		return &BatchFailureError{ExitCode: exitCode, Reason: acpxExitReasonTimeout, Stderr: stderr}
	case 5:
		update := StreamUpdate{Kind: StreamUpdateStatus, Status: acpxPermissionDeniedStatus}
		event := newAgentRunEvent(req, update, marshalStatusPayload(acpxPermissionDeniedStatus), eventClock(runner.Now)())
		if err := sink.Publish(ctx, event); err != nil {
			return fmt.Errorf("publish acpx permission-denied Run Event: %w", err)
		}
		return &BatchFailureError{ExitCode: exitCode, Reason: acpxExitReasonPermissionsDenied, Stderr: stderr}
	case 2:
		return &InfrastructureError{ExitCode: exitCode, Reason: acpxExitReasonUsage, Stderr: stderr}
	case 4:
		return &InfrastructureError{ExitCode: exitCode, Reason: acpxExitReasonMissingSession, Stderr: stderr}
	case 130:
		return StopError{LogPath: req.LogPath, Output: output, Err: context.Canceled}
	default:
		return &InfrastructureError{ExitCode: exitCode, Reason: "unexpected acpx exit code", Stderr: stderr}
	}
}

func commandExitCode(err error) (int, bool) {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), true
	}
	return 0, false
}
