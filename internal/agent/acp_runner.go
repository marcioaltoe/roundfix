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
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"roundfix/internal/runevent"

	acp "github.com/coder/acp-go-sdk"
)

// ACPRunner runs ACP Agents and publishes their session updates as Run
// Events. Now overrides the event clock; nil means time.Now.
type ACPRunner struct {
	Now func() time.Time
}

type acpSession struct {
	id           string
	workingDir   string
	allowedRoots []string
}

type acpClient struct {
	runtime RuntimeSpec
	req     ExecuteRequest
	sink    runevent.Sink
	now     func() time.Time
	runCtx  context.Context
	log     io.Writer
	output  *strings.Builder

	pubMu  sync.Mutex
	pubErr error

	mu       sync.Mutex
	conn     *acp.ClientSideConnection
	cmd      *exec.Cmd
	waitCh   chan error
	sessions map[string]*acpSession

	terminalMu   sync.Mutex
	terminalNext int
	terminals    map[string]*acpTerminal
}

type acpTerminal struct {
	id        string
	sessionID string
	cancel    context.CancelFunc
	cmd       *exec.Cmd
	output    *terminalOutput
	done      chan struct{}

	mu       sync.Mutex
	exitCode *int
	signal   *string
}

type terminalOutput struct {
	mu        sync.Mutex
	data      []byte
	limit     int
	truncated bool
}

const defaultTerminalOutputLimit = 10 * 1024 * 1024

func (runner ACPRunner) Probe(ctx context.Context, runtime RuntimeSpec) error {
	if _, err := resolveACPCommand(ctx, runtime, runtimeEffectiveModel(runtime), true); err != nil {
		return err
	}
	return nil
}

func (runner ACPRunner) Run(ctx context.Context, req ExecuteRequest, sink runevent.Sink) (ExecuteResult, error) {
	if strings.TrimSpace(req.LogPath) == "" {
		return ExecuteResult{}, errors.New("Agent log path is required")
	}
	if err := os.MkdirAll(filepath.Dir(req.LogPath), 0o755); err != nil {
		return ExecuteResult{}, fmt.Errorf("create Agent log directory: %w", err)
	}
	logFile, err := os.Create(req.LogPath)
	if err != nil {
		return ExecuteResult{}, fmt.Errorf("create Agent log %q: %w", req.LogPath, err)
	}
	defer func() {
		_ = logFile.Close()
	}()
	if sink == nil {
		sink = runevent.Discard
	}

	var output strings.Builder
	client := &acpClient{
		runtime:  req.Runtime,
		req:      req,
		sink:     sink,
		now:      eventClock(runner.Now),
		runCtx:   ctx,
		log:      logFile,
		output:   &output,
		sessions: make(map[string]*acpSession),
	}
	if err := client.start(ctx); err != nil {
		return ExecuteResult{LogPath: req.LogPath, Output: output.String()}, err
	}
	defer func() {
		_ = client.close(stopGrace(req.StopGrace))
	}()

	sessionID, err := client.newSession(ctx)
	if err != nil {
		return ExecuteResult{LogPath: req.LogPath, Output: output.String()}, err
	}
	promptErr := make(chan error, 1)
	go func() {
		promptErr <- client.prompt(ctx, sessionID, req.Prompt)
	}()

	select {
	case err := <-promptErr:
		if err != nil {
			return ExecuteResult{LogPath: req.LogPath, Output: output.String()}, err
		}
	case <-ctx.Done():
		_ = client.cancelPrompt(sessionID)
		killed := client.stop(stopGrace(req.StopGrace))
		client.emitStatus(context.WithoutCancel(ctx), "stopped")
		return ExecuteResult{LogPath: req.LogPath, Output: output.String()}, StopError{
			LogPath: req.LogPath,
			Output:  output.String(),
			Killed:  killed,
			Err:     ctx.Err(),
		}
	}

	if err := client.publishErr(); err != nil {
		return ExecuteResult{LogPath: req.LogPath, Output: output.String()}, fmt.Errorf("publish Run Events: %w", err)
	}
	return ExecuteResult{LogPath: req.LogPath, Output: output.String()}, nil
}

func (c *acpClient) start(ctx context.Context) error {
	modelName := runtimeEffectiveModel(c.runtime)
	command, err := resolveACPCommand(ctx, c.runtime, launchModel(c.runtime, modelName), false)
	if err != nil {
		return err
	}
	if len(command) == 0 {
		return errors.New("empty ACP command")
	}
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = c.req.GitRoot
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("open ACP stdin: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("open ACP stdout: %w", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start ACP Agent command %q: %w", command[0], err)
	}

	c.cmd = cmd
	c.waitCh = make(chan error, 1)
	go func() {
		c.waitCh <- cmd.Wait()
	}()

	conn := acp.NewClientSideConnection(c, stdin, c.interceptUpdates(stdout))
	c.conn = conn
	if _, err := conn.Initialize(ctx, acp.InitializeRequest{
		ProtocolVersion: acp.ProtocolVersionNumber,
		ClientCapabilities: acp.ClientCapabilities{
			Fs: acp.FileSystemCapabilities{
				ReadTextFile:  true,
				WriteTextFile: true,
			},
			Terminal: true,
		},
		ClientInfo: &acp.Implementation{
			Name:    "roundfix",
			Version: "dev",
		},
	}); err != nil {
		_ = c.stop(0)
		detail := strings.TrimSpace(stderr.String())
		if detail != "" {
			return fmt.Errorf("initialize ACP Agent %s: %w: %s", shellCommand(command), err, detail)
		}
		return fmt.Errorf("initialize ACP Agent %s: %w", shellCommand(command), err)
	}
	return nil
}

func (c *acpClient) newSession(ctx context.Context) (acp.SessionId, error) {
	workingDir, err := filepath.Abs(c.req.GitRoot)
	if err != nil {
		return "", fmt.Errorf("resolve Agent working directory: %w", err)
	}
	resp, err := c.conn.NewSession(ctx, acp.NewSessionRequest{
		Cwd:        workingDir,
		McpServers: []acp.McpServer{},
	})
	if err != nil {
		return "", fmt.Errorf("create ACP session: %w", err)
	}
	allowedRoots, err := allowedSessionRoots(workingDir, c.req.AllowAddDirs)
	if err != nil {
		return "", err
	}
	session := &acpSession{
		id:           string(resp.SessionId),
		workingDir:   workingDir,
		allowedRoots: allowedRoots,
	}
	c.mu.Lock()
	c.sessions[session.id] = session
	c.mu.Unlock()

	modelName := runtimeEffectiveModel(c.runtime)
	if modelName != "" && modelName != launchModel(c.runtime, modelName) {
		if err := c.setSessionModel(ctx, resp, modelName); err != nil {
			return "", err
		}
	}
	if c.runtime.FullAccessMode != "" {
		if _, err := c.conn.SetSessionMode(ctx, acp.SetSessionModeRequest{
			SessionId: resp.SessionId,
			ModeId:    acp.SessionModeId(c.runtime.FullAccessMode),
		}); err != nil {
			return "", fmt.Errorf("set ACP session mode %q: %w", c.runtime.FullAccessMode, err)
		}
	}
	return resp.SessionId, nil
}

// setSessionModel selects the runtime model through session config options,
// which replaced the session/set_model method in current ACP revisions.
func (c *acpClient) setSessionModel(ctx context.Context, resp acp.NewSessionResponse, modelName string) error {
	option := modelConfigOption(resp.ConfigOptions)
	if option == nil {
		return fmt.Errorf("set ACP session model %q: agent advertises no model config option", modelName)
	}
	if _, err := c.conn.SetSessionConfigOption(ctx, acp.SetSessionConfigOptionRequest{
		ValueId: &acp.SetSessionConfigOptionValueId{
			SessionId: resp.SessionId,
			ConfigId:  option.Id,
			Value:     acp.SessionConfigValueId(modelName),
		},
	}); err != nil {
		return fmt.Errorf("set ACP session model %q: %w", modelName, err)
	}
	return nil
}

func modelConfigOption(options []acp.SessionConfigOption) *acp.SessionConfigOptionSelect {
	for _, option := range options {
		sel := option.Select
		if sel == nil || sel.Category == nil {
			continue
		}
		if *sel.Category == acp.SessionConfigOptionCategoryModel {
			return sel
		}
	}
	return nil
}

func (c *acpClient) prompt(ctx context.Context, sessionID acp.SessionId, prompt string) error {
	resp, err := c.conn.Prompt(ctx, acp.PromptRequest{
		SessionId: sessionID,
		Prompt:    []acp.ContentBlock{acp.TextBlock(prompt)},
	})
	if err != nil {
		return fmt.Errorf("run ACP prompt: %w", err)
	}
	if resp.StopReason == acp.StopReasonCancelled {
		return context.Canceled
	}
	c.emitStatus(ctx, "completed")
	return nil
}

func (c *acpClient) cancelPrompt(sessionID acp.SessionId) error {
	if c.conn == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return c.conn.Cancel(ctx, acp.CancelNotification{SessionId: sessionID})
}

func (c *acpClient) close(grace time.Duration) error {
	if c == nil {
		return nil
	}
	terminalErr := c.closeTerminals(grace)
	killed := c.stop(grace)
	if killed {
		return errors.Join(terminalErr, errors.New("ACP Agent process killed during shutdown"))
	}
	return terminalErr
}

func (c *acpClient) stop(grace time.Duration) bool {
	if c == nil || c.cmd == nil || c.waitCh == nil {
		return false
	}
	if c.cmd.Process == nil {
		return false
	}
	_ = c.cmd.Process.Signal(os.Interrupt)
	if grace <= 0 {
		grace = 3 * time.Second
	}
	timer := time.NewTimer(grace)
	defer timer.Stop()
	select {
	case <-c.waitCh:
		return false
	case <-timer.C:
		_ = c.cmd.Process.Kill()
		// cmd.Wait can outlive the kill: grandchildren that inherited the
		// agent's stderr keep the pipe open and Wait blocks on the copier.
		// The Run must never freeze on runtime teardown, so the wait is
		// bounded and the reaper goroutine is abandoned past it.
		killTimer := time.NewTimer(grace)
		defer killTimer.Stop()
		select {
		case <-c.waitCh:
		case <-killTimer.C:
		}
		return true
	}
}

// interceptUpdates tees the agent's newline-delimited JSON-RPC output so
// session/update notifications are published with their raw params bytes
// (ADR 0008) before the SDK connection processes the line. Handling before
// forwarding keeps event order ahead of prompt completion.
func (c *acpClient) interceptUpdates(stdout io.Reader) io.Reader {
	pipeReader, pipeWriter := io.Pipe()
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			line, err := reader.ReadBytes('\n')
			if len(line) > 0 {
				c.handleWireLine(line)
				if _, writeErr := pipeWriter.Write(line); writeErr != nil {
					// The connection side is gone; drain so the process can exit.
					_, _ = io.Copy(io.Discard, stdout)
					return
				}
			}
			if err != nil {
				if errors.Is(err, io.EOF) {
					_ = pipeWriter.Close()
				} else {
					_ = pipeWriter.CloseWithError(err)
				}
				return
			}
		}
	}()
	return pipeReader
}

type wireEnvelope struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

func (c *acpClient) handleWireLine(line []byte) {
	var envelope wireEnvelope
	if err := json.Unmarshal(line, &envelope); err != nil {
		return
	}
	if envelope.Method != acp.ClientMethodSessionUpdate {
		return
	}
	var note acp.SessionNotification
	if err := json.Unmarshal(envelope.Params, &note); err != nil {
		return
	}
	update := streamUpdateFromACP(note.Update)
	if update.Kind == "" {
		return
	}
	_ = c.emit(c.runCtx, update, json.RawMessage(envelope.Params))
}

func (c *acpClient) emit(ctx context.Context, update StreamUpdate, payload json.RawMessage) error {
	text := formatStreamUpdate(update)
	if text != "" {
		if c.log != nil {
			if _, err := io.WriteString(c.log, text); err != nil {
				return err
			}
		}
		if c.output != nil {
			c.output.WriteString(text)
		}
	}
	event := newAgentRunEvent(c.req, update, payload, c.now())
	if err := c.sink.Publish(ctx, event); err != nil {
		c.recordPublishErr(err)
		return err
	}
	return nil
}

func (c *acpClient) emitStatus(ctx context.Context, status string) {
	update := StreamUpdate{Kind: StreamUpdateStatus, Status: status}
	_ = c.emit(ctx, update, marshalStatusPayload(status))
}

func (c *acpClient) recordPublishErr(err error) {
	c.pubMu.Lock()
	defer c.pubMu.Unlock()
	if c.pubErr == nil {
		c.pubErr = err
	}
}

func (c *acpClient) publishErr() error {
	c.pubMu.Lock()
	defer c.pubMu.Unlock()
	return c.pubErr
}

func (c *acpClient) ReadTextFile(_ context.Context, params acp.ReadTextFileRequest) (acp.ReadTextFileResponse, error) {
	path, err := c.resolveSessionPath(params.SessionId, params.Path)
	if err != nil {
		return acp.ReadTextFileResponse{}, err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return acp.ReadTextFileResponse{}, err
	}
	return acp.ReadTextFileResponse{Content: string(content)}, nil
}

func (c *acpClient) WriteTextFile(_ context.Context, params acp.WriteTextFileRequest) (acp.WriteTextFileResponse, error) {
	path, err := c.resolveSessionPath(params.SessionId, params.Path)
	if err != nil {
		return acp.WriteTextFileResponse{}, err
	}
	mode := os.FileMode(0o600)
	if info, statErr := os.Stat(path); statErr == nil {
		mode = info.Mode().Perm()
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return acp.WriteTextFileResponse{}, fmt.Errorf("stat session file %q: %w", path, statErr)
	}
	if err := os.WriteFile(path, []byte(params.Content), mode); err != nil {
		return acp.WriteTextFileResponse{}, err
	}
	return acp.WriteTextFileResponse{}, nil
}

func (c *acpClient) RequestPermission(_ context.Context, params acp.RequestPermissionRequest) (acp.RequestPermissionResponse, error) {
	if len(params.Options) == 0 {
		return acp.RequestPermissionResponse{Outcome: acp.NewRequestPermissionOutcomeCancelled()}, nil
	}
	return acp.RequestPermissionResponse{Outcome: acp.NewRequestPermissionOutcomeSelected(params.Options[0].OptionId)}, nil
}

// SessionUpdate satisfies the ACP client interface. Publication happens in
// handleWireLine, which sees the raw notification bytes; reacting here too
// would publish every update twice.
func (c *acpClient) SessionUpdate(context.Context, acp.SessionNotification) error {
	return nil
}

func (c *acpClient) resolveSessionPath(sessionID acp.SessionId, rawPath string) (string, error) {
	c.mu.Lock()
	session := c.sessions[string(sessionID)]
	c.mu.Unlock()
	if session == nil {
		return "", fmt.Errorf("received ACP file request for unknown session %q", sessionID)
	}
	path, err := resolvePath(session.workingDir, rawPath)
	if err != nil {
		return "", err
	}
	if !pathInsideAnyRoot(path, session.allowedRoots) {
		return "", fmt.Errorf("ACP file path %q is outside allowed session roots", rawPath)
	}
	return path, nil
}

func streamUpdateFromACP(update acp.SessionUpdate) StreamUpdate {
	switch {
	case update.AgentMessageChunk != nil:
		return StreamUpdate{Kind: StreamUpdateMessage, Text: contentBlockText(update.AgentMessageChunk.Content)}
	case update.AgentThoughtChunk != nil:
		return StreamUpdate{Kind: StreamUpdateThought, Text: contentBlockText(update.AgentThoughtChunk.Content)}
	case update.ToolCall != nil:
		blocks := toolContentBlocks(update.ToolCall.Content, update.ToolCall.RawInput, update.ToolCall.RawOutput)
		return StreamUpdate{
			Kind:      StreamUpdateToolStarted,
			Title:     update.ToolCall.Title,
			ToolID:    string(update.ToolCall.ToolCallId),
			ToolState: string(update.ToolCall.Status),
			Text:      streamBlockText(blocks),
			Blocks:    blocks,
		}
	case update.ToolCallUpdate != nil:
		title := ""
		if update.ToolCallUpdate.Title != nil {
			title = *update.ToolCallUpdate.Title
		}
		state := ""
		if update.ToolCallUpdate.Status != nil {
			state = string(*update.ToolCallUpdate.Status)
		}
		blocks := toolContentBlocks(update.ToolCallUpdate.Content, update.ToolCallUpdate.RawInput, update.ToolCallUpdate.RawOutput)
		return StreamUpdate{
			Kind:      StreamUpdateToolUpdated,
			Title:     title,
			ToolID:    string(update.ToolCallUpdate.ToolCallId),
			ToolState: state,
			Text:      streamBlockText(blocks),
			Blocks:    blocks,
		}
	case update.Plan != nil:
		lines := make([]string, 0, len(update.Plan.Entries))
		for _, entry := range update.Plan.Entries {
			line := strings.TrimSpace(entry.Content)
			if line == "" {
				continue
			}
			if entry.Status != "" {
				line = fmt.Sprintf("%s  %s", entry.Status, line)
			}
			lines = append(lines, line)
		}
		return StreamUpdate{Kind: StreamUpdatePlan, Text: strings.Join(lines, "\n")}
	default:
		return StreamUpdate{}
	}
}

func contentBlockText(block acp.ContentBlock) string {
	if block.Text != nil {
		return block.Text.Text
	}
	raw, err := json.Marshal(block)
	if err != nil {
		return ""
	}
	return string(raw)
}

func toolContentBlocks(content []acp.ToolCallContent, rawInput any, rawOutput any) []StreamBlock {
	blocks := []StreamBlock{}
	if rawInput != nil {
		blocks = append(blocks, StreamBlock{Kind: StreamBlockInput, Text: formatToolInput(rawInput)})
	}
	for _, item := range content {
		switch {
		case item.Content != nil && item.Content.Content.Text != nil:
			blocks = append(blocks, StreamBlock{Kind: StreamBlockText, Text: item.Content.Content.Text.Text})
		case item.Content != nil && item.Content.Content.Image != nil:
			blocks = append(blocks, StreamBlock{
				Kind:     StreamBlockImage,
				MimeType: item.Content.Content.Image.MimeType,
				URI:      stringValue(item.Content.Content.Image.Uri),
			})
		case item.Content != nil && item.Content.Content.ResourceLink != nil:
			blocks = append(blocks, StreamBlock{
				Kind:     StreamBlockResource,
				Name:     item.Content.Content.ResourceLink.Name,
				URI:      item.Content.Content.ResourceLink.Uri,
				MimeType: stringValue(item.Content.Content.ResourceLink.MimeType),
			})
		case item.Diff != nil:
			blocks = append(blocks, StreamBlock{Kind: StreamBlockDiff, Path: item.Diff.Path})
		case item.Terminal != nil:
			blocks = append(blocks, StreamBlock{Kind: StreamBlockTerminal, TerminalID: item.Terminal.TerminalId})
		}
	}
	if rawOutput != nil {
		blocks = append(blocks, StreamBlock{Kind: StreamBlockOutput, Text: formatToolOutput(rawOutput)})
	}
	return blocks
}

// formatToolInput renders tool input as a log line, not raw JSON: shell
// invocations show the command itself; anything else falls back to compact
// JSON.
func formatToolInput(rawInput any) string {
	fields, ok := rawInput.(map[string]any)
	if !ok {
		return marshalCompact(rawInput)
	}
	if command := shellCommandFromInput(fields["command"]); command != "" {
		return command
	}
	return marshalCompact(rawInput)
}

// shellCommandFromInput extracts the human command from exec-style inputs
// like ["/bin/zsh","-lc","rtk go test ./..."].
func shellCommandFromInput(value any) string {
	switch command := value.(type) {
	case string:
		return strings.TrimSpace(command)
	case []any:
		parts := make([]string, 0, len(command))
		for _, item := range command {
			text, ok := item.(string)
			if !ok {
				return ""
			}
			parts = append(parts, text)
		}
		if len(parts) == 0 {
			return ""
		}
		last := strings.TrimSpace(parts[len(parts)-1])
		if len(parts) >= 3 && strings.HasPrefix(parts[1], "-") && last != "" {
			return last
		}
		return strings.Join(parts, " ")
	default:
		return ""
	}
}

// toolOutputMaxLines bounds rendered tool output so chatty commands stay
// readable; the raw payload keeps the full data.
const toolOutputMaxLines = 8

// formatToolOutput renders tool output as log text: exec-style payloads show
// their aggregated output, bounded; anything else falls back to compact JSON.
func formatToolOutput(rawOutput any) string {
	fields, ok := rawOutput.(map[string]any)
	if !ok {
		return marshalCompact(rawOutput)
	}
	for _, key := range []string{"aggregated_output", "formatted_output", "stdout", "output"} {
		text, found := fields[key].(string)
		if !found {
			continue
		}
		return boundOutputLines(text)
	}
	return marshalCompact(rawOutput)
}

func boundOutputLines(text string) string {
	trimmed := strings.TrimRight(text, "\r\n")
	if trimmed == "" {
		return ""
	}
	lines := strings.Split(trimmed, "\n")
	if len(lines) <= toolOutputMaxLines {
		return trimmed
	}
	kept := lines[:toolOutputMaxLines]
	return strings.Join(kept, "\n") + fmt.Sprintf("\n… (+%d more line(s))", len(lines)-toolOutputMaxLines)
}

func streamBlockText(blocks []StreamBlock) string {
	parts := []string{}
	for _, block := range blocks {
		switch block.Kind {
		case StreamBlockText:
			if strings.TrimSpace(block.Text) != "" {
				parts = append(parts, block.Text)
			}
		case StreamBlockInput:
			if strings.TrimSpace(block.Text) != "" {
				parts = append(parts, "input: "+block.Text)
			}
		case StreamBlockOutput:
			if strings.TrimSpace(block.Text) != "" {
				parts = append(parts, "output: "+block.Text)
			}
		case StreamBlockDiff:
			if block.Path != "" {
				parts = append(parts, "diff: "+block.Path)
			}
		case StreamBlockTerminal:
			if block.TerminalID != "" {
				parts = append(parts, "terminal: "+block.TerminalID)
			}
		case StreamBlockImage:
			parts = append(parts, "image: "+firstNonEmpty(block.MimeType, block.URI, "image"))
		case StreamBlockResource:
			parts = append(parts, "resource: "+firstNonEmpty(block.Name, block.URI, "resource"))
		}
	}
	return strings.Join(parts, "\n")
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func marshalCompact(value any) string {
	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(raw)
}

func runtimeEffectiveModel(runtime RuntimeSpec) string {
	if selected := strings.TrimSpace(runtime.Model); selected != "" && !strings.EqualFold(selected, "auto") {
		return normalizeCodexModel(runtime, selected)
	}
	return normalizeCodexModel(runtime, runtime.DefaultModel)
}

func launchModel(runtime RuntimeSpec, modelName string) string {
	if runtime.BootstrapModel {
		return modelName
	}
	return runtime.DefaultModel
}

func normalizeCodexModel(runtime RuntimeSpec, modelName string) string {
	modelName = strings.TrimSpace(modelName)
	if runtime.ID != "codex" {
		return modelName
	}
	if unprefixed, ok := strings.CutPrefix(modelName, "codex/"); ok {
		return strings.TrimSpace(unprefixed)
	}
	if provider, unprefixed, ok := strings.Cut(modelName, "/"); ok && strings.TrimSpace(provider) != "" && !strings.Contains(unprefixed, "/") {
		return strings.TrimSpace(unprefixed)
	}
	return modelName
}

func resolveACPCommand(ctx context.Context, runtime RuntimeSpec, modelName string, verify bool) ([]string, error) {
	launchers := append([]RuntimeLauncher{{Command: runtime.Command, Args: runtime.Args, ProbeArgs: runtime.ProbeArgs}}, runtime.Fallbacks...)
	var failures []string
	for _, launcher := range launchers {
		command := strings.TrimSpace(launcher.Command)
		if command == "" {
			continue
		}
		if _, err := exec.LookPath(command); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", command, err))
			continue
		}
		args := append([]string{}, launcher.Args...)
		args = append(args, runtimeBootstrapArgs(runtime, modelName)...)
		if verify && len(launcher.ProbeArgs) > 0 {
			probe := exec.CommandContext(ctx, command, launcher.ProbeArgs...)
			if output, err := probe.CombinedOutput(); err != nil {
				failures = append(failures, fmt.Sprintf("%s %s: %s", command, strings.Join(launcher.ProbeArgs, " "), strings.TrimSpace(string(output))))
				continue
			}
		}
		return append([]string{command}, args...), nil
	}
	if len(failures) == 0 {
		failures = append(failures, "empty command")
	}
	return nil, ProbeError{Runtime: runtime, Err: errors.New(strings.Join(failures, "; "))}
}

func runtimeBootstrapArgs(runtime RuntimeSpec, modelName string) []string {
	if runtime.ID != "codex" {
		return nil
	}
	args := []string{}
	if strings.TrimSpace(modelName) != "" {
		args = appendCodexConfig(args, "model="+strconv.Quote(modelName))
	}
	args = appendCodexConfig(args,
		"features.code_mode=false",
		"features.code_mode_only=false",
		`approval_policy="never"`,
		// Batches run with the full-access preset (ADR 0011): the
		// workspace-write sandbox blocks all network access, including
		// localhost services that verification commands depend on.
		`sandbox_mode="danger-full-access"`,
	)
	return args
}

func appendCodexConfig(args []string, values ...string) []string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			args = append(args, "-c", value)
		}
	}
	return args
}

func shellCommand(args []string) string {
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, shellArg(arg))
	}
	return strings.Join(parts, " ")
}

func shellArg(value string) string {
	if value == "" {
		return "''"
	}
	if strings.ContainsAny(value, " \t\n\"'\\$`|&;<>*?[]{}()") {
		return "'" + strings.ReplaceAll(value, "'", `'\"'\"'`) + "'"
	}
	return value
}

func allowedSessionRoots(workingDir string, addDirs []string) ([]string, error) {
	roots := []string{workingDir}
	for _, dir := range addDirs {
		if strings.TrimSpace(dir) == "" {
			continue
		}
		abs, err := filepath.Abs(dir)
		if err != nil {
			return nil, fmt.Errorf("resolve ACP add-dir %q: %w", dir, err)
		}
		roots = append(roots, filepath.Clean(abs))
	}
	return roots, nil
}

func resolvePath(workingDir string, rawPath string) (string, error) {
	trimmed := strings.TrimSpace(rawPath)
	if trimmed == "" {
		return "", errors.New("ACP path is required")
	}
	if filepath.IsAbs(trimmed) {
		return filepath.Clean(trimmed), nil
	}
	return filepath.Clean(filepath.Join(workingDir, trimmed)), nil
}

func pathInsideAnyRoot(path string, roots []string) bool {
	for _, root := range roots {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			continue
		}
		if rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))) {
			return true
		}
	}
	return false
}

func (c *acpClient) CreateTerminal(ctx context.Context, params acp.CreateTerminalRequest) (acp.CreateTerminalResponse, error) {
	session := c.lookupTerminalSession(params.SessionId)
	if session == nil {
		return acp.CreateTerminalResponse{}, fmt.Errorf("received terminal request for unknown session %q", params.SessionId)
	}
	cwd := session.workingDir
	if params.Cwd != nil && strings.TrimSpace(*params.Cwd) != "" {
		resolved, err := resolvePath(session.workingDir, *params.Cwd)
		if err != nil {
			return acp.CreateTerminalResponse{}, err
		}
		cwd = resolved
	}
	if !pathInsideAnyRoot(cwd, session.allowedRoots) {
		return acp.CreateTerminalResponse{}, fmt.Errorf("terminal cwd %q is outside allowed session roots", cwd)
	}
	if err := ctx.Err(); err != nil {
		return acp.CreateTerminalResponse{}, err
	}
	terminalCtx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(terminalCtx, params.Command, params.Args...)
	cmd.Dir = cwd
	cmd.Env = terminalEnv(params.Env)
	output := newTerminalOutput(params.OutputByteLimit)
	cmd.Stdout = output
	cmd.Stderr = output
	if err := cmd.Start(); err != nil {
		cancel()
		return acp.CreateTerminalResponse{}, fmt.Errorf("start terminal command %q: %w", params.Command, err)
	}
	terminal := &acpTerminal{
		id:        c.nextTerminalID(),
		sessionID: string(params.SessionId),
		cancel:    cancel,
		cmd:       cmd,
		output:    output,
		done:      make(chan struct{}),
	}
	c.storeTerminal(terminal)
	go terminal.wait()
	return acp.CreateTerminalResponse{TerminalId: terminal.id}, nil
}

func (c *acpClient) KillTerminal(_ context.Context, params acp.KillTerminalRequest) (acp.KillTerminalResponse, error) {
	terminal, err := c.lookupTerminal(params.SessionId, params.TerminalId)
	if err != nil {
		return acp.KillTerminalResponse{}, err
	}
	terminal.kill()
	return acp.KillTerminalResponse{}, nil
}

func (c *acpClient) TerminalOutput(_ context.Context, params acp.TerminalOutputRequest) (acp.TerminalOutputResponse, error) {
	terminal, err := c.lookupTerminal(params.SessionId, params.TerminalId)
	if err != nil {
		return acp.TerminalOutputResponse{}, err
	}
	output, truncated := terminal.output.snapshot()
	resp := acp.TerminalOutputResponse{Output: output, Truncated: truncated}
	if status := terminal.exitStatus(); status != nil {
		resp.ExitStatus = status
	}
	return resp, nil
}

func (c *acpClient) ReleaseTerminal(ctx context.Context, params acp.ReleaseTerminalRequest) (acp.ReleaseTerminalResponse, error) {
	terminal, err := c.removeTerminal(params.SessionId, params.TerminalId)
	if err != nil {
		return acp.ReleaseTerminalResponse{}, err
	}
	terminal.kill()
	if err := terminal.waitFor(ctx); err != nil {
		return acp.ReleaseTerminalResponse{}, err
	}
	return acp.ReleaseTerminalResponse{}, nil
}

func (c *acpClient) WaitForTerminalExit(ctx context.Context, params acp.WaitForTerminalExitRequest) (acp.WaitForTerminalExitResponse, error) {
	terminal, err := c.lookupTerminal(params.SessionId, params.TerminalId)
	if err != nil {
		return acp.WaitForTerminalExitResponse{}, err
	}
	if err := terminal.waitFor(ctx); err != nil {
		return acp.WaitForTerminalExitResponse{}, err
	}
	exitCode, signal := terminal.exitResult()
	return acp.WaitForTerminalExitResponse{ExitCode: exitCode, Signal: signal}, nil
}

func (c *acpClient) lookupTerminalSession(sessionID acp.SessionId) *acpSession {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sessions[string(sessionID)]
}

func (c *acpClient) nextTerminalID() string {
	c.terminalMu.Lock()
	defer c.terminalMu.Unlock()
	c.terminalNext++
	return "term-" + strconv.Itoa(c.terminalNext)
}

func (c *acpClient) storeTerminal(terminal *acpTerminal) {
	c.terminalMu.Lock()
	defer c.terminalMu.Unlock()
	if c.terminals == nil {
		c.terminals = make(map[string]*acpTerminal)
	}
	c.terminals[terminal.id] = terminal
}

func (c *acpClient) lookupTerminal(sessionID acp.SessionId, terminalID string) (*acpTerminal, error) {
	c.terminalMu.Lock()
	defer c.terminalMu.Unlock()
	terminal := c.terminals[terminalID]
	if terminal == nil {
		return nil, fmt.Errorf("unknown terminal %q", terminalID)
	}
	if terminal.sessionID != string(sessionID) {
		return nil, fmt.Errorf("terminal %q does not belong to session %q", terminalID, sessionID)
	}
	return terminal, nil
}

func (c *acpClient) removeTerminal(sessionID acp.SessionId, terminalID string) (*acpTerminal, error) {
	c.terminalMu.Lock()
	defer c.terminalMu.Unlock()
	terminal := c.terminals[terminalID]
	if terminal == nil {
		return nil, fmt.Errorf("unknown terminal %q", terminalID)
	}
	if terminal.sessionID != string(sessionID) {
		return nil, fmt.Errorf("terminal %q does not belong to session %q", terminalID, sessionID)
	}
	delete(c.terminals, terminalID)
	return terminal, nil
}

func (c *acpClient) closeTerminals(grace time.Duration) error {
	c.terminalMu.Lock()
	terminals := make([]*acpTerminal, 0, len(c.terminals))
	for id, terminal := range c.terminals {
		terminals = append(terminals, terminal)
		delete(c.terminals, id)
	}
	c.terminalMu.Unlock()
	if len(terminals) == 0 {
		return nil
	}
	if grace <= 0 {
		grace = 3 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), grace)
	defer cancel()
	var result error
	for _, terminal := range terminals {
		terminal.kill()
		if err := terminal.waitFor(ctx); err != nil {
			result = errors.Join(result, fmt.Errorf("wait for terminal %s: %w", terminal.id, err))
		}
	}
	return result
}

func (t *acpTerminal) wait() {
	waitErr := t.cmd.Wait()
	t.cancel()
	var exitCode *int
	var signal *string
	if t.cmd.ProcessState != nil {
		code := t.cmd.ProcessState.ExitCode()
		if code >= 0 {
			exitCode = &code
		}
	}
	if exitCode == nil && waitErr != nil {
		message := waitErr.Error()
		signal = &message
	}
	t.mu.Lock()
	t.exitCode = exitCode
	t.signal = signal
	close(t.done)
	t.mu.Unlock()
}

func (t *acpTerminal) kill() {
	if t != nil && t.cancel != nil {
		t.cancel()
	}
}

func (t *acpTerminal) waitFor(ctx context.Context) error {
	if t == nil {
		return errors.New("terminal process is required")
	}
	select {
	case <-t.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *acpTerminal) exitResult() (*int, *string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return cloneInt(t.exitCode), cloneString(t.signal)
}

func (t *acpTerminal) exitStatus() *acp.TerminalExitStatus {
	select {
	case <-t.done:
	default:
		return nil
	}
	exitCode, signal := t.exitResult()
	return &acp.TerminalExitStatus{ExitCode: exitCode, Signal: signal}
}

func cloneInt(src *int) *int {
	if src == nil {
		return nil
	}
	value := *src
	return &value
}

func cloneString(src *string) *string {
	if src == nil {
		return nil
	}
	value := *src
	return &value
}

func newTerminalOutput(limit *int) *terminalOutput {
	resolved := defaultTerminalOutputLimit
	if limit != nil && *limit > 0 {
		resolved = *limit
	}
	return &terminalOutput{limit: resolved}
}

func (output *terminalOutput) Write(payload []byte) (int, error) {
	output.mu.Lock()
	defer output.mu.Unlock()
	output.data = append(output.data, payload...)
	if output.limit > 0 && len(output.data) > output.limit {
		output.data = trimUTF8Tail(output.data, output.limit)
		output.truncated = true
	}
	return len(payload), nil
}

func (output *terminalOutput) snapshot() (string, bool) {
	output.mu.Lock()
	defer output.mu.Unlock()
	return string(append([]byte(nil), output.data...)), output.truncated
}

func trimUTF8Tail(data []byte, limit int) []byte {
	if limit <= 0 || len(data) <= limit {
		return data
	}
	start := len(data) - limit
	for start < len(data) && !utf8.RuneStart(data[start]) {
		start++
	}
	return append([]byte(nil), data[start:]...)
}

func terminalEnv(env []acp.EnvVariable) []string {
	merged := os.Environ()
	for _, item := range env {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		merged = append(merged, name+"="+item.Value)
	}
	return merged
}
