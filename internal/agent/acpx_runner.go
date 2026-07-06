package agent

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"roundfix/internal/runevent"
)

const (
	defaultACPXCommand              = "acpx"
	PinnedACPXVersion               = "0.12.0"
	acpxPermissionDeniedStatus      = "permissions_denied"
	acpxExitReasonAgentProtocol     = "agent/protocol error"
	acpxExitReasonTimeout           = "timeout"
	acpxExitReasonPermissionsDenied = "all permissions denied"
	acpxExitReasonUsage             = "usage error"
	acpxExitReasonMissingSession    = "missing session"
	acpxCodexSandboxUnavailable     = "codex_sandbox_full_access_unavailable"
	acpxCodexSandboxModeKey         = "sandbox_mode"
	acpxCodexFullAccessSandbox      = "danger-full-access"
	infrastructureStderrTailLines   = 10
	infrastructureStderrTailBytes   = 1024
	infrastructureStderrDelimiter   = "\n--- acpx stderr tail ---\n"
	infrastructureStderrTruncated   = "[stderr truncated]\n"
)

// ACPXRunner is the acpx-backed invocation core. Later migration tasks wire
// this into Runner after Agent Session lifecycle is available.
type ACPXRunner struct {
	Command         string
	Now             func() time.Time
	warnf           func(string, ...any)
	ensuredSessions map[string]struct{}
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
	if tail, truncated := infrastructureStderrTail(err.Stderr); tail != "" {
		message += infrastructureStderrDelimiter
		if truncated {
			message += infrastructureStderrTruncated
		}
		message += tail
	}
	return message
}

func infrastructureStderrTail(stderr string) (string, bool) {
	tail := strings.TrimSpace(stderr)
	if tail == "" {
		return "", false
	}
	truncated := false
	lines := strings.Split(tail, "\n")
	if len(lines) > infrastructureStderrTailLines {
		lines = lines[len(lines)-infrastructureStderrTailLines:]
		tail = strings.Join(lines, "\n")
		truncated = true
	}
	tail = strings.TrimSpace(tail)
	if len(tail) > infrastructureStderrTailBytes {
		tail = tail[len(tail)-infrastructureStderrTailBytes:]
		truncated = true
	}
	tail = strings.TrimSpace(tail)
	return tail, truncated
}

type ACPXProbeError struct {
	Command         string
	FoundVersion    string
	RequiredVersion string
	Missing         bool
	Err             error
}

func (err ACPXProbeError) Error() string {
	install := acpxInstallCommand()
	if err.Missing {
		return fmt.Sprintf("acpx is required but was not found on PATH; install it with: %s", install)
	}
	if err.FoundVersion != "" {
		return fmt.Sprintf("acpx version mismatch: found %s, but Roundfix requires %s; upgrade or downgrade with: %s", err.FoundVersion, err.RequiredVersion, install)
	}
	message := "acpx --version failed"
	if err.Err != nil {
		message = fmt.Sprintf("%s: %v", message, err.Err)
	}
	return fmt.Sprintf("%s; install it with: %s", message, install)
}

func (err ACPXProbeError) Unwrap() error {
	return err.Err
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

type acpxStreamResult struct {
	output             string
	stopReason         string
	promptResultParsed bool
	err                error
}

func (runner ACPXRunner) Probe(ctx context.Context, _ RuntimeSpec) error {
	command := runner.command()
	if _, err := exec.LookPath(command); err != nil {
		return ACPXProbeError{Command: command, RequiredVersion: PinnedACPXVersion, Missing: true, Err: err}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, command, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		detail := strings.TrimSpace(string(output))
		if detail == "" {
			detail = err.Error()
		}
		return ACPXProbeError{Command: command, RequiredVersion: PinnedACPXVersion, Err: errors.New(detail)}
	}
	foundVersion := strings.TrimSpace(string(output))
	if foundVersion != PinnedACPXVersion {
		return ACPXProbeError{Command: command, FoundVersion: displayACPXVersion(foundVersion), RequiredVersion: PinnedACPXVersion}
	}
	return nil
}

func (runner *ACPXRunner) Run(ctx context.Context, req ExecuteRequest, sink runevent.Sink) (ExecuteResult, error) {
	result := ExecuteResult{LogPath: req.LogPath}
	if err := runner.ensureSession(ctx, req, sink); err != nil {
		return result, err
	}
	return runner.RunPrompt(ctx, ACPXPromptRequest{
		ExecuteRequest: req,
		Session:        req.Session.Name,
	}, sink)
}

func (runner *ACPXRunner) EndSession(ctx context.Context, runtime RuntimeSpec, session SessionRef) error {
	sessionName := strings.TrimSpace(session.Name)
	if sessionName == "" {
		return nil
	}
	if err := runner.CloseSession(ctx, runtime, session); err != nil {
		runner.warningf("close acpx Agent Session %q: %v", sessionName, err)
	}
	delete(runner.ensuredSessions, sessionName)
	return nil
}

func (runner *ACPXRunner) CloseSession(ctx context.Context, runtime RuntimeSpec, session SessionRef) error {
	sessionName := strings.TrimSpace(session.Name)
	if sessionName == "" {
		return nil
	}
	args, err := acpxCloseArgs(runtime, sessionName, session.WorkDir)
	if err != nil {
		return fmt.Errorf("build acpx close command for Agent Session %q: %w", sessionName, err)
	}
	if err := runner.runACPXCommand(ctx, args); err != nil {
		return fmt.Errorf("close acpx Agent Session %q: %w", sessionName, err)
	}
	return nil
}

func (runner ACPXRunner) ListRoundfixSessions(ctx context.Context, runtime RuntimeSpec, workDir string) ([]RoundfixSession, error) {
	args, err := acpxListSessionsArgs(runtime, workDir)
	if err != nil {
		return nil, fmt.Errorf("build acpx sessions list command: %w", err)
	}
	output, err := runner.runACPXCommandOutput(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("list acpx Agent Sessions: %w", err)
	}
	return ParseRoundfixSessions(output), nil
}

func (runner *ACPXRunner) CancelSession(ctx context.Context, runtime RuntimeSpec, session SessionRef) error {
	sessionName := strings.TrimSpace(session.Name)
	if sessionName == "" {
		return errors.New("agent session name is required")
	}
	args, err := acpxCancelArgs(runtime, sessionName, session.WorkDir)
	if err != nil {
		return fmt.Errorf("build acpx cancel command for Agent Session %q: %w", sessionName, err)
	}
	if err := runner.runACPXCommand(ctx, args); err != nil {
		return fmt.Errorf("cancel acpx Agent Session %q: %w", sessionName, err)
	}
	return nil
}

func (runner *ACPXRunner) ensureSession(ctx context.Context, req ExecuteRequest, sink runevent.Sink) error {
	sessionName := strings.TrimSpace(req.Session.Name)
	if sessionName == "" {
		return errors.New("agent session is required")
	}
	workDir := strings.TrimSpace(req.GitRoot)
	if workDir == "" {
		return errors.New("agent working directory is required")
	}
	if _, ok := runner.ensuredSessions[sessionName]; ok {
		return nil
	}
	args, err := acpxEnsureArgs(req.Runtime, sessionName, workDir)
	if err != nil {
		return err
	}
	if err := runner.runACPXCommand(ctx, args); err != nil {
		return fmt.Errorf("ensure acpx Agent Session %q: %w", sessionName, err)
	}
	if err := runner.applyFullAccess(ctx, req, sink); err != nil {
		return err
	}
	if runner.ensuredSessions == nil {
		runner.ensuredSessions = map[string]struct{}{}
	}
	if err := runner.publishStatus(ctx, req, sink, AgentSessionStartedStatus); err != nil {
		return err
	}
	runner.ensuredSessions[sessionName] = struct{}{}
	return nil
}

func (runner *ACPXRunner) applyFullAccess(ctx context.Context, req ExecuteRequest, sink runevent.Sink) error {
	mode := strings.TrimSpace(req.Runtime.FullAccessMode)
	if mode == "" {
		return nil
	}
	sessionName := strings.TrimSpace(req.Session.Name)
	args, err := acpxSetModeArgs(req.Runtime, mode, sessionName, req.GitRoot)
	if err != nil {
		return err
	}
	if err := runner.runACPXCommand(ctx, args); err != nil {
		return fmt.Errorf("set acpx Agent Session mode %q: %w", mode, err)
	}
	if req.Runtime.ID != "codex" || mode != "full-access" {
		return nil
	}
	args, err = acpxSetConfigArgs(req.Runtime, acpxCodexSandboxModeKey, acpxCodexFullAccessSandbox, sessionName, req.GitRoot)
	if err != nil {
		return err
	}
	if err := runner.runACPXCommand(ctx, args); err != nil {
		if isCodexSandboxUnavailable(err) {
			runner.warningf("codex full-access sandbox preset unavailable for Agent Session %q: %v", sessionName, err)
			if publishErr := runner.publishStatus(ctx, req, sink, acpxCodexSandboxUnavailable); publishErr != nil {
				return publishErr
			}
			return nil
		}
		return fmt.Errorf("set acpx Codex sandbox preset %q: %w", acpxCodexFullAccessSandbox, err)
	}
	return nil
}

func (runner ACPXRunner) RunPrompt(ctx context.Context, req ACPXPromptRequest, sink runevent.Sink) (ExecuteResult, error) {
	result := ExecuteResult{LogPath: req.LogPath}
	if err := validateACPXPromptRequest(req); err != nil {
		return result, err
	}
	if sink == nil {
		sink = runevent.Discard
	}
	logWriter := io.Writer(io.Discard)
	if strings.TrimSpace(req.LogPath) != "" {
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
		logWriter = logFile
	}

	if err := ctx.Err(); err != nil {
		return result, StopError{LogPath: req.LogPath, Err: err}
	}
	args, err := acpxPromptArgs(req)
	if err != nil {
		return result, err
	}
	cmd := exec.Command(runner.command(), args...)
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

	streamCh := make(chan acpxStreamResult, 1)
	go func() {
		streamCh <- runner.readPromptStream(ctx, req.ExecuteRequest, sink, stdout, logWriter)
	}()
	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
	}()

	var waitErr error
	select {
	case waitErr = <-waitCh:
	case <-ctx.Done():
		forceClosed := runner.cancelPrompt(context.WithoutCancel(ctx), req, waitCh)
		stream := waitForACPXStream(streamCh, stopGrace(req.StopGrace))
		result.Output = stream.output
		result.StopReason = stream.stopReason
		_ = runner.publishStatus(context.WithoutCancel(ctx), req.ExecuteRequest, sink, "stopped")
		return result, StopError{LogPath: req.LogPath, Output: result.Output, Killed: forceClosed, Err: ctx.Err()}
	}
	stream := <-streamCh
	result.Output = stream.output
	result.StopReason = stream.stopReason
	if waitErr != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return result, StopError{LogPath: req.LogPath, Output: result.Output, Err: ctxErr}
		}
		exitCode, ok := commandExitCode(waitErr)
		if !ok {
			return result, fmt.Errorf("wait for acpx prompt: %w", waitErr)
		}
		if stream.promptResultParsed && exitCode != 130 {
			result.TransportAnomaly = acpxTransportAnomaly(exitCode, stderr.String())
			return result, nil
		}
		return result, runner.mapExitCode(ctx, req.ExecuteRequest, sink, exitCode, stderr.String(), result.Output)
	}
	if stream.err != nil {
		return result, &BatchFailureError{Reason: stream.err.Error(), Stderr: stderr.String()}
	}
	if result.StopReason == "" {
		return result, &BatchFailureError{Reason: "missing session/prompt stop reason", Stderr: stderr.String()}
	}
	return result, nil
}

func (runner ACPXRunner) readPromptStream(ctx context.Context, req ExecuteRequest, sink runevent.Sink, stdout io.Reader, logFile io.Writer) acpxStreamResult {
	var output bytes.Buffer
	var stopReason string
	var promptResultParsed bool
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
			if err := runner.handleStdoutLine(ctx, req, sink, line, &stopReason, &promptResultParsed); err != nil && streamErr == nil {
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
	return acpxStreamResult{output: output.String(), stopReason: stopReason, promptResultParsed: promptResultParsed, err: streamErr}
}

func validateACPXPromptRequest(req ACPXPromptRequest) error {
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

func acpxInstallCommand() string {
	return "npm install -g acpx@" + PinnedACPXVersion
}

func displayACPXVersion(version string) string {
	if strings.TrimSpace(version) == "" {
		return "<empty>"
	}
	return version
}

func (runner ACPXRunner) runACPXCommand(ctx context.Context, args []string) error {
	_, err := runner.runACPXCommandOutput(ctx, args)
	return err
}

func (runner ACPXRunner) runACPXCommandOutput(ctx context.Context, args []string) (string, error) {
	cmd := exec.CommandContext(ctx, runner.command(), args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if exitCode, ok := commandExitCode(err); ok {
			return stdout.String(), &InfrastructureError{ExitCode: exitCode, Reason: "acpx command failed", Stderr: stderr.String()}
		}
		return stdout.String(), fmt.Errorf("run acpx command: %w", err)
	}
	return stdout.String(), nil
}

// acpx CLI grammar: --cwd, --format, --json-strict, --approve-all, --model,
// and --agent are program-level globals and must appear before the agent
// name / subcommand. Subcommands accept only their own options (prompt,
// cancel, set-mode, set take -s <session>; sessions ensure takes --name;
// sessions close takes the session name as a positional argument). Every
// session-scoped invocation carries the global --cwd so session resolution
// is deterministic regardless of the Roundfix process cwd.

func acpxPromptArgs(req ACPXPromptRequest) ([]string, error) {
	agentArgs, err := acpxAgentArgs(req.Runtime)
	if err != nil {
		return nil, err
	}
	args := []string{
		"--cwd", strings.TrimSpace(req.GitRoot),
		"--format", "json",
		"--json-strict",
		"--approve-all",
	}
	if model := strings.TrimSpace(req.Runtime.Model); model != "" {
		args = append(args, "--model", model)
	}
	args = append(args, agentArgs...)
	args = append(args, "prompt", "-s", strings.TrimSpace(req.Session), "-f", "-")
	return args, nil
}

func acpxEnsureArgs(runtime RuntimeSpec, sessionName string, workDir string) ([]string, error) {
	args, err := acpxGlobalArgs(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(args, "sessions", "ensure", "--name", strings.TrimSpace(sessionName)), nil
}

func acpxSetModeArgs(runtime RuntimeSpec, mode string, sessionName string, workDir string) ([]string, error) {
	args, err := acpxGlobalArgs(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(args, "set-mode", strings.TrimSpace(mode), "-s", strings.TrimSpace(sessionName)), nil
}

func acpxSetConfigArgs(runtime RuntimeSpec, key string, value string, sessionName string, workDir string) ([]string, error) {
	args, err := acpxGlobalArgs(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(args, "set", strings.TrimSpace(key), strings.TrimSpace(value), "-s", strings.TrimSpace(sessionName)), nil
}

func acpxCancelArgs(runtime RuntimeSpec, sessionName string, workDir string) ([]string, error) {
	args, err := acpxGlobalArgs(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(args, "cancel", "-s", strings.TrimSpace(sessionName)), nil
}

func acpxCloseArgs(runtime RuntimeSpec, sessionName string, workDir string) ([]string, error) {
	args, err := acpxGlobalArgs(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(args, "sessions", "close", strings.TrimSpace(sessionName)), nil
}

func acpxListSessionsArgs(runtime RuntimeSpec, workDir string) ([]string, error) {
	args, err := acpxGlobalArgs(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(args, "sessions", "list"), nil
}

func acpxGlobalArgs(runtime RuntimeSpec, workDir string) ([]string, error) {
	workDir = strings.TrimSpace(workDir)
	if workDir == "" {
		return nil, errors.New("Agent working directory is required")
	}
	agentArgs, err := acpxAgentArgs(runtime)
	if err != nil {
		return nil, err
	}
	return append([]string{"--cwd", workDir}, agentArgs...), nil
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

func (runner ACPXRunner) handleStdoutLine(ctx context.Context, req ExecuteRequest, sink runevent.Sink, line []byte, stopReason *string, promptResultParsed *bool) error {
	var message acpxJSONRPCMessage
	if err := json.Unmarshal(line, &message); err != nil {
		return fmt.Errorf("parse acpx stdout JSON-RPC line: %w", err)
	}
	if message.Error != nil {
		return fmt.Errorf("acpx JSON-RPC error %d: %s", message.Error.Code, message.Error.Message)
	}
	if message.Method == acpMethodSessionUpdate {
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
		*promptResultParsed = true
	}
	return nil
}

func acpxTransportAnomaly(exitCode int, stderr string) string {
	message := fmt.Sprintf("acpx exited with exit code %d after parsed session/prompt result", exitCode)
	if tail, truncated := infrastructureStderrTail(stderr); tail != "" {
		message += infrastructureStderrDelimiter
		if truncated {
			message += infrastructureStderrTruncated
		}
		message += tail
	}
	return message
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

func (runner ACPXRunner) cancelPrompt(ctx context.Context, req ACPXPromptRequest, waitCh <-chan error) bool {
	grace := stopGrace(req.StopGrace)
	cancelCtx, cancel := context.WithTimeout(ctx, grace)
	defer cancel()
	if args, err := acpxCancelArgs(req.Runtime, req.Session, req.GitRoot); err == nil {
		if err := runner.runACPXCommand(cancelCtx, args); err != nil {
			runner.warningf("cancel acpx Agent Session %q: %v", req.Session, err)
		}
	} else {
		runner.warningf("build acpx cancel command for Agent Session %q: %v", req.Session, err)
	}
	timer := time.NewTimer(grace)
	defer timer.Stop()
	select {
	case <-waitCh:
		return false
	case <-timer.C:
		closeCtx, closeCancel := context.WithTimeout(ctx, grace)
		defer closeCancel()
		if args, err := acpxCloseArgs(req.Runtime, req.Session, req.GitRoot); err == nil {
			if err := runner.runACPXCommand(closeCtx, args); err != nil {
				runner.warningf("close acpx Agent Session %q: %v", req.Session, err)
			}
		} else {
			runner.warningf("build acpx close command for Agent Session %q: %v", req.Session, err)
		}
		closeTimer := time.NewTimer(grace)
		defer closeTimer.Stop()
		select {
		case <-waitCh:
		case <-closeTimer.C:
		}
		return true
	}
}

func waitForACPXStream(streamCh <-chan acpxStreamResult, grace time.Duration) acpxStreamResult {
	timer := time.NewTimer(grace)
	defer timer.Stop()
	select {
	case stream := <-streamCh:
		return stream
	case <-timer.C:
		return acpxStreamResult{}
	}
}

func (runner ACPXRunner) publishStatus(ctx context.Context, req ExecuteRequest, sink runevent.Sink, status string) error {
	if sink == nil {
		sink = runevent.Discard
	}
	update := StreamUpdate{Kind: StreamUpdateStatus, Status: status}
	event := newAgentRunEvent(req, update, marshalStatusPayload(status), eventClock(runner.Now)())
	if err := sink.Publish(ctx, event); err != nil {
		return fmt.Errorf("publish acpx Agent status %q: %w", status, err)
	}
	return nil
}

func (runner ACPXRunner) warningf(format string, args ...any) {
	if runner.warnf != nil {
		runner.warnf(format, args...)
		return
	}
	log.Printf(format, args...)
}

func isCodexSandboxUnavailable(err error) bool {
	text := strings.ToLower(err.Error())
	if !strings.Contains(text, acpxCodexSandboxModeKey) && !strings.Contains(text, acpxCodexFullAccessSandbox) {
		return false
	}
	return strings.Contains(text, "unavailable") ||
		strings.Contains(text, "unknown") ||
		strings.Contains(text, "invalid") ||
		strings.Contains(text, "unsupported") ||
		strings.Contains(text, "does not advertise config option")
}

func commandExitCode(err error) (int, bool) {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), true
	}
	return 0, false
}
