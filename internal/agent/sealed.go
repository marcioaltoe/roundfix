package agent

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	SealedPromptMaxInputBytes  = 2 << 20
	SealedPromptMaxOutputBytes = 512 << 10
	SealedPromptTimeout        = 2 * time.Minute

	sealedACPStreamMaxBytes = 2 << 20
	sealedSessionPrefix     = "roundfix-baseline-"
)

var (
	ErrSealedToolUse        = errors.New("sealed ACP attempt emitted a tool event")
	ErrSealedOutputTooLarge = errors.New("sealed ACP output exceeds the configured limit")
)

// SealedPromptRequest contains the complete non-Run ACP prompt boundary. Input
// is passed to stdin unchanged.
type SealedPromptRequest struct {
	Runtime RuntimeSpec
	WorkDir string
	Input   []byte
}

// SealedPromptResult contains only the accumulated Agent message. Raw ACP
// protocol output and thought are discarded.
type SealedPromptResult struct {
	Output   []byte
	ToolUsed bool
}

// RunSealedPrompt opens one fresh exactly selected Agent Session, runs one
// denied-tools prompt, and proves the Session closed before returning.
func (runner *ACPXRunner) RunSealedPrompt(
	ctx context.Context,
	request SealedPromptRequest,
) (SealedPromptResult, error) {
	if err := validateSealedPromptRequest(ctx, request); err != nil {
		return SealedPromptResult{}, err
	}
	attemptCtx, attemptCancel := context.WithTimeout(ctx, SealedPromptTimeout)
	defer attemptCancel()

	adapter, err := CheckAdapter(attemptCtx, request.Runtime)
	if err != nil {
		return SealedPromptResult{}, err
	}
	sessionName, err := sealedAgentSessionName()
	if err != nil {
		return SealedPromptResult{}, err
	}
	codexEnv, err := runner.codexEnvForSession(attemptCtx, request.Runtime, sessionName)
	if err != nil {
		return SealedPromptResult{}, err
	}

	session := SessionRef{Name: sessionName, WorkDir: request.WorkDir}
	capabilities, attemptErr := runner.startSessionSelectionWithEnsure(
		attemptCtx,
		request.Runtime,
		session,
		adapter,
		codexEnv,
		"ensure sealed Agent Session",
		SealedEnsureArguments,
	)
	var result SealedPromptResult
	if attemptErr == nil {
		var proof SelectionProof
		proof, attemptErr = runner.applySessionSelection(
			attemptCtx,
			SessionSelectionRequest{
				Runtime:      request.Runtime,
				Session:      session,
				Capabilities: capabilities,
			},
			codexEnv,
		)
		if attemptErr == nil {
			selected := request
			selected.Runtime.Model = proof.Assignment.AdapterModel
			result, attemptErr = runner.runSealedPromptCommand(
				attemptCtx,
				selected,
				sessionName,
				codexEnv,
			)
		}
	}
	attemptCancel()

	cleanupCtx, cleanupCancel := context.WithTimeout(
		context.WithoutCancel(ctx),
		acpxPreflightCleanupTimeout,
	)
	defer cleanupCancel()
	if ctx.Err() != nil ||
		errors.Is(attemptErr, context.Canceled) ||
		errors.Is(attemptErr, context.DeadlineExceeded) {
		_ = runner.CancelSession(cleanupCtx, request.Runtime, session)
	}
	cleanupErr := runner.closeDisposableSession(
		cleanupCtx,
		request.Runtime,
		sessionName,
		request.WorkDir,
	)
	if cleanupErr != nil {
		if attemptErr != nil {
			return SealedPromptResult{}, errors.Join(attemptErr, cleanupErr)
		}
		return SealedPromptResult{}, cleanupErr
	}
	if err := ctx.Err(); err != nil {
		return SealedPromptResult{}, err
	}
	if attemptErr != nil {
		return SealedPromptResult{}, attemptErr
	}
	return result, nil
}

// SealedEnsureArguments exposes the exact security-sensitive ACPX creation
// contract for deterministic adapter tests.
func SealedEnsureArguments(
	runtime RuntimeSpec,
	sessionName string,
	workDir string,
) ([]string, error) {
	sessionName = strings.TrimSpace(sessionName)
	if sessionName == "" {
		return nil, errors.New("sealed Agent Session name is required")
	}
	args, err := sealedACPXGlobalArguments(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(
		args,
		"sessions",
		"ensure",
		"--name",
		sessionName,
	), nil
}

// SealedPromptArguments exposes the exact security-sensitive ACPX prompt
// contract for deterministic adapter tests.
func SealedPromptArguments(
	runtime RuntimeSpec,
	sessionName string,
	workDir string,
) ([]string, error) {
	sessionName = strings.TrimSpace(sessionName)
	if sessionName == "" {
		return nil, errors.New("sealed Agent Session name is required")
	}
	args, err := sealedACPXGlobalArguments(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(
		args,
		"prompt",
		"-s",
		sessionName,
		"-f",
		"-",
	), nil
}

func sealedACPXGlobalArguments(
	runtime RuntimeSpec,
	workDir string,
) ([]string, error) {
	workDir = strings.TrimSpace(workDir)
	if workDir == "" {
		return nil, errors.New("sealed Agent working directory is required")
	}
	agentArgs, err := acpxAgentArgs(runtime)
	if err != nil {
		return nil, err
	}
	args := []string{
		"--cwd", workDir,
		"--format", "json",
		"--json-strict",
		"--deny-all",
		"--non-interactive-permissions", "fail",
		"--allowed-tools", "",
		"--max-turns", "1",
		"--prompt-retries", "0",
		"--no-terminal",
		"--timeout", strconv.Itoa(int(SealedPromptTimeout.Seconds())),
		"--ttl", strconv.Itoa(int(SealedPromptTimeout.Seconds())),
	}
	model := strings.TrimSpace(runtime.Model)
	if model == "" {
		return nil, errors.New("sealed Agent model is required")
	}
	args = append(args, "--model", model)
	args = append(args, agentArgs...)
	return args, nil
}

func validateSealedPromptRequest(ctx context.Context, request SealedPromptRequest) error {
	if ctx == nil {
		return errors.New("sealed prompt context is required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := validateRuntimeSelection(request.Runtime); err != nil {
		return err
	}
	if len(request.Input) == 0 {
		return errors.New("sealed prompt input is required")
	}
	if len(request.Input) > SealedPromptMaxInputBytes {
		return fmt.Errorf(
			"sealed prompt input is %d bytes; maximum is %d bytes",
			len(request.Input),
			SealedPromptMaxInputBytes,
		)
	}
	workDir := strings.TrimSpace(request.WorkDir)
	if workDir == "" {
		return errors.New("sealed Agent working directory is required")
	}
	absolute, err := filepath.Abs(workDir)
	if err != nil {
		return fmt.Errorf("resolve sealed Agent working directory: %w", err)
	}
	if filepath.Clean(workDir) != absolute {
		return errors.New("sealed Agent working directory must be absolute")
	}
	info, err := os.Stat(absolute)
	if err != nil {
		return fmt.Errorf("inspect sealed Agent working directory: %w", err)
	}
	if !info.IsDir() || info.Mode().Perm()&0o077 != 0 {
		return errors.New("sealed Agent working directory must be a private directory")
	}
	entries, err := os.ReadDir(absolute)
	if err != nil {
		return fmt.Errorf("inspect sealed Agent working directory entries: %w", err)
	}
	if len(entries) != 0 {
		return errors.New("sealed Agent working directory must be empty")
	}
	return nil
}

func sealedAgentSessionName() (string, error) {
	var entropy [8]byte
	if _, err := rand.Read(entropy[:]); err != nil {
		return "", fmt.Errorf("generate sealed Agent Session name: %w", err)
	}
	return fmt.Sprintf("%s%x", sealedSessionPrefix, entropy[:]), nil
}

func (runner *ACPXRunner) runSealedPromptCommand(
	ctx context.Context,
	request SealedPromptRequest,
	sessionName string,
	codexEnv []string,
) (SealedPromptResult, error) {
	args, err := SealedPromptArguments(request.Runtime, sessionName, request.WorkDir)
	if err != nil {
		return SealedPromptResult{}, err
	}
	command := exec.CommandContext(ctx, runner.command(), args...)
	command.Env = acpxCommandEnv(codexEnv)
	command.Dir = request.WorkDir
	command.Stdin = bytes.NewReader(request.Input)
	stdout := newSealedCapture(sealedACPStreamMaxBytes)
	stderr := newSealedCapture(infrastructureStderrTailBytes)
	command.Stdout = stdout
	command.Stderr = stderr
	if err := command.Run(); err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return SealedPromptResult{}, ctxErr
		}
		if errors.Is(err, ErrSealedOutputTooLarge) || stdout.overflow {
			return SealedPromptResult{}, ErrSealedOutputTooLarge
		}
		if exitCode, ok := commandExitCode(err); ok {
			return SealedPromptResult{}, &InfrastructureError{
				ExitCode: exitCode,
				Reason:   "sealed acpx prompt failed",
				Stderr:   stderr.String(),
			}
		}
		return SealedPromptResult{}, fmt.Errorf("run sealed acpx prompt: %w", err)
	}
	if stdout.overflow {
		return SealedPromptResult{}, ErrSealedOutputTooLarge
	}
	return parseSealedPromptStream(stdout.Bytes())
}

func parseSealedPromptStream(stream []byte) (SealedPromptResult, error) {
	var output bytes.Buffer
	var stopReason string
	sawResult := false
	scanner := bufio.NewScanner(bytes.NewReader(stream))
	scanner.Buffer(make([]byte, 64<<10), sealedACPStreamMaxBytes)
	for scanner.Scan() {
		line := append([]byte(nil), scanner.Bytes()...)
		var message acpxJSONRPCMessage
		if err := json.Unmarshal(line, &message); err != nil {
			return SealedPromptResult{}, fmt.Errorf("parse sealed acpx JSON-RPC line: %w", err)
		}
		if message.Error != nil {
			return SealedPromptResult{}, fmt.Errorf(
				"sealed acpx JSON-RPC error %d: %s",
				message.Error.Code,
				message.Error.Message,
			)
		}
		if message.Method == acpMethodSessionUpdate {
			update, ok, err := streamUpdateFromSessionUpdatePayload(line)
			if err != nil {
				return SealedPromptResult{}, err
			}
			if !ok {
				continue
			}
			switch update.Kind {
			case StreamUpdateMessage:
				if output.Len()+len(update.Text) > SealedPromptMaxOutputBytes {
					return SealedPromptResult{}, ErrSealedOutputTooLarge
				}
				_, _ = output.WriteString(update.Text)
			case StreamUpdateThought:
				// Thought remains ephemeral and outside the proposal.
			case StreamUpdateToolStarted, StreamUpdateToolUpdated:
				return SealedPromptResult{ToolUsed: true}, ErrSealedToolUse
			default:
				return SealedPromptResult{}, fmt.Errorf(
					"sealed ACP attempt emitted unsupported update %q",
					update.Kind,
				)
			}
			continue
		}
		if len(message.Result) == 0 {
			continue
		}
		var response struct {
			StopReason string `json:"stopReason"`
		}
		if err := json.Unmarshal(message.Result, &response); err != nil {
			return SealedPromptResult{}, fmt.Errorf(
				"parse sealed acpx prompt response: %w",
				err,
			)
		}
		if sawResult {
			return SealedPromptResult{}, errors.New(
				"sealed acpx prompt returned more than one terminal result",
			)
		}
		sawResult = true
		stopReason = response.StopReason
	}
	if err := scanner.Err(); err != nil {
		return SealedPromptResult{}, fmt.Errorf("read sealed acpx output: %w", err)
	}
	if !sawResult || strings.TrimSpace(stopReason) == "" {
		return SealedPromptResult{}, errors.New(
			"sealed acpx prompt omitted its terminal result",
		)
	}
	if stopReason != "end_turn" {
		return SealedPromptResult{}, fmt.Errorf(
			"sealed acpx prompt stopped with reason %q",
			stopReason,
		)
	}
	return SealedPromptResult{Output: append([]byte(nil), output.Bytes()...)}, nil
}

type sealedCapture struct {
	buffer   bytes.Buffer
	limit    int
	overflow bool
}

func newSealedCapture(limit int) *sealedCapture {
	return &sealedCapture{limit: limit}
}

func (capture *sealedCapture) Write(data []byte) (int, error) {
	remaining := capture.limit - capture.buffer.Len()
	if remaining <= 0 {
		capture.overflow = true
		return 0, ErrSealedOutputTooLarge
	}
	if len(data) <= remaining {
		return capture.buffer.Write(data)
	}
	_, _ = capture.buffer.Write(data[:remaining])
	capture.overflow = true
	return remaining, ErrSealedOutputTooLarge
}

func (capture *sealedCapture) Bytes() []byte {
	return capture.buffer.Bytes()
}

func (capture *sealedCapture) String() string {
	return capture.buffer.String()
}
