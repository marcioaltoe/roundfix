package agent

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
)

const realACPXIntegrationEnv = "ROUNDFIX_REAL_ACPX"

func TestRealACPXCommandOverrideRoundTripCancelCrashResume(t *testing.T) {
	if os.Getenv(realACPXIntegrationEnv) != "1" {
		t.Skipf("set %s=1 to run the real acpx integration test; requires Node and acpx %s", realACPXIntegrationEnv, PinnedACPXVersion)
	}
	if _, err := exec.LookPath("node"); err != nil {
		t.Skipf("Node is required for real acpx integration: %v", err)
	}
	acpxPath, err := exec.LookPath(defaultACPXCommand)
	if err != nil {
		t.Skipf("acpx is required for real acpx integration: %v", err)
	}

	runner := &ACPXRunner{Command: acpxPath}
	if err := runner.Probe(context.Background(), RuntimeSpec{}); err != nil {
		t.Skipf("pinned acpx %s is required for real acpx integration: %v", PinnedACPXVersion, err)
	}

	dir := t.TempDir()
	helper := buildRealACPXEchoAgent(t, dir)
	helperLog := filepath.Join(dir, "echo-agent.log")
	helperState := filepath.Join(dir, "echo-agent-state")
	t.Setenv("ROUNDFIX_ACPX_ECHO_LOG", helperLog)
	t.Setenv("ROUNDFIX_ACPX_ECHO_STATE", helperState)

	runtime := RuntimeSpec{ID: "codex-custom", Protocol: ProtocolStdio, Command: helper}
	session := SessionRef{Name: "roundfix-real-acpx-" + strconv.FormatInt(time.Now().UnixNano(), 36), WorkDir: dir}
	t.Cleanup(func() {
		_ = runner.EndSession(context.Background(), runtime, session)
	})

	roundTripSink := newCaptureSink("")
	roundTrip, err := runner.Run(context.Background(), realACPXRequest(dir, runtime, session, 1, "round-trip hello"), roundTripSink)
	if err != nil {
		t.Fatalf("round-trip prompt through real acpx: %v", err)
	}
	if roundTrip.StopReason != "end_turn" {
		t.Fatalf("expected end_turn stop reason, got %q", roundTrip.StopReason)
	}
	if !sinkContainsPayload(roundTripSink, "echo: round-trip hello") {
		t.Fatalf("expected raw NDJSON payload to journal echo response, got %+v", roundTripSink.Events())
	}
	if content := readFile(t, roundTrip.LogPath); !strings.Contains(content, `"method":"session/update"`) || !strings.Contains(content, "echo: round-trip hello") {
		t.Fatalf("expected raw NDJSON in Agent log, got:\n%s", content)
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	cancelSink := newCaptureSink("prompt-blocked")
	cancelErr := make(chan error, 1)
	go func() {
		_, err := runner.Run(cancelCtx, realACPXRequest(dir, runtime, session, 2, "block-until-cancel"), cancelSink)
		cancelErr <- err
	}()
	select {
	case <-cancelSink.done:
	case <-time.After(10 * time.Second):
		cancel()
		t.Fatal("timed out waiting for blocking prompt to start")
	}
	cancel()
	err = receiveError(t, cancelErr)
	var stopErr StopError
	if !errors.As(err, &stopErr) {
		t.Fatalf("expected cooperative cancel to return StopError, got %T %v", err, err)
	}
	if stopErr.Killed {
		t.Fatalf("expected cooperative cancel without kill fallback, got %#v", stopErr)
	}
	waitForLogLine(t, helperLog, "cancel echo-session")

	_, crashErr := runner.Run(context.Background(), realACPXRequest(dir, runtime, session, 3, "crash-adapter"), newCaptureSink(""))
	if crashErr == nil {
		t.Fatal("expected induced adapter crash to fail the prompt Run")
	}
	waitForLogLine(t, helperLog, "crash echo-session")

	resumeSink := newCaptureSink("")
	resumed, err := runner.Run(context.Background(), realACPXRequest(dir, runtime, session, 4, "after-crash"), resumeSink)
	if err != nil {
		t.Fatalf("expected next Work Item to succeed after adapter crash: %v", err)
	}
	if resumed.StopReason != "end_turn" {
		t.Fatalf("expected resumed prompt stop reason end_turn, got %q", resumed.StopReason)
	}
	if !sinkContainsPayload(resumeSink, "echo: after-crash") {
		t.Fatalf("expected resumed prompt to journal echo response, got %+v", resumeSink.Events())
	}

	if err := runner.EndSession(context.Background(), runtime, session); err != nil {
		t.Fatalf("close real acpx session: %v", err)
	}
	assertNoRealACPXOrphans(t, session.Name, helperLog)
}

func realACPXRequest(dir string, runtime RuntimeSpec, session SessionRef, batch int, prompt string) ExecuteRequest {
	return ExecuteRequest{
		Runtime:   runtime,
		Session:   session,
		RunID:     "run-real-acpx",
		Batch:     rounds.Batch{Number: batch},
		LogPath:   filepath.Join(dir, "runs", "run-real-acpx", "agent", fmt.Sprintf("batch-%03d.log", batch)),
		Prompt:    prompt,
		GitRoot:   dir,
		StopGrace: 2 * time.Second,
	}
}

func sinkContainsPayload(sink *captureSink, needle string) bool {
	for _, event := range sink.Events() {
		if event.Source == runevent.SourceAgent && strings.Contains(string(event.Payload), needle) {
			return true
		}
	}
	return false
}

func buildRealACPXEchoAgent(t *testing.T, dir string) string {
	t.Helper()
	source := filepath.Join(dir, "echo-agent.go")
	binary := filepath.Join(dir, "echo-agent")
	if err := os.WriteFile(source, []byte(realACPXEchoAgentSource), 0o644); err != nil {
		t.Fatalf("write echo agent source: %v", err)
	}
	cmd := exec.Command("go", "build", "-o", binary, source)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build echo agent: %v\n%s", err, string(output))
	}
	return binary
}

func waitForLogLine(t *testing.T, path string, needle string) {
	t.Helper()
	deadline := time.After(10 * time.Second)
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	for {
		if content, err := os.ReadFile(path); err == nil && strings.Contains(string(content), needle) {
			return
		}
		select {
		case <-deadline:
			content, _ := os.ReadFile(path)
			t.Fatalf("timed out waiting for %q in %s\n%s", needle, path, string(content))
		case <-ticker.C:
		}
	}
}

func assertNoRealACPXOrphans(t *testing.T, sessionName string, helperLog string) {
	t.Helper()
	deadline := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		helperPIDs := liveEchoAgentPIDs(helperLog)
		ownerLines := realACPXOwnerLines(sessionName)
		if len(helperPIDs) == 0 && len(ownerLines) == 0 {
			return
		}
		select {
		case <-deadline:
			t.Fatalf("real acpx integration left owner/helper processes for %s; helper pids=%v owner lines=%v", sessionName, helperPIDs, ownerLines)
		case <-ticker.C:
		}
	}
}

func liveEchoAgentPIDs(helperLog string) []int {
	content, err := os.ReadFile(helperLog)
	if err != nil {
		return nil
	}
	seen := map[int]struct{}{}
	for _, line := range strings.Split(string(content), "\n") {
		raw, ok := strings.CutPrefix(line, "pid=")
		if !ok {
			continue
		}
		pid, err := strconv.Atoi(strings.TrimSpace(raw))
		if err != nil {
			continue
		}
		seen[pid] = struct{}{}
	}
	live := []int{}
	for pid := range seen {
		if processAlive(pid) {
			live = append(live, pid)
		}
	}
	return live
}

func processAlive(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}

func realACPXOwnerLines(sessionName string) []string {
	cmd := exec.Command("ps", "-axo", "pid=,command=")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}
	lines := []string{}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "acpx") && strings.Contains(line, sessionName) {
			lines = append(lines, strings.TrimSpace(line))
		}
	}
	return lines
}

const realACPXEchoAgentSource = `package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type rpcMessage struct {
	JSONRPC string          ` + "`json:\"jsonrpc\"`" + `
	ID      json.RawMessage ` + "`json:\"id\"`" + `
	Method  string          ` + "`json:\"method\"`" + `
	Params  json.RawMessage ` + "`json:\"params\"`" + `
}

type promptParams struct {
	SessionID string ` + "`json:\"sessionId\"`" + `
	Prompt   []struct {
		Type string ` + "`json:\"type\"`" + `
		Text string ` + "`json:\"text\"`" + `
	} ` + "`json:\"prompt\"`" + `
}

type cancelParams struct {
	SessionID string ` + "`json:\"sessionId\"`" + `
}

type echoAgent struct {
	writeMu sync.Mutex
	logMu   sync.Mutex
	mu      sync.Mutex
	cancels map[string]chan struct{}
	writer  *bufio.Writer
	logFile *os.File
}

func main() {
	logPath := os.Getenv("ROUNDFIX_ACPX_ECHO_LOG")
	if logPath == "" {
		fmt.Fprintln(os.Stderr, "ROUNDFIX_ACPX_ECHO_LOG is required")
		os.Exit(2)
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open log: %v\n", err)
		os.Exit(2)
	}
	agent := &echoAgent{
		cancels: map[string]chan struct{}{},
		writer:  bufio.NewWriter(os.Stdout),
		logFile: logFile,
	}
	agent.logf("pid=%d", os.Getpid())
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		var msg rpcMessage
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			agent.logf("decode error %v", err)
			continue
		}
		agent.logf("method %s", msg.Method)
		if msg.Method == "session/prompt" {
			go agent.handlePrompt(msg)
			continue
		}
		agent.handle(msg)
	}
	agent.logf("stdin closed")
}

func (a *echoAgent) handle(msg rpcMessage) {
	switch msg.Method {
	case "initialize":
		a.respond(msg.ID, map[string]any{
			"protocolVersion": 1,
			"authMethods":     []any{},
			"agentCapabilities": map[string]any{
				"loadSession": true,
				"sessionCapabilities": map[string]any{
					"resume": map[string]any{},
					"close":  map[string]any{},
				},
			},
		})
	case "session/new":
		a.persistSession("echo-session")
		a.logf("new echo-session")
		a.respond(msg.ID, map[string]any{"sessionId": "echo-session"})
	case "session/load", "session/resume":
		a.logf("resume echo-session")
		a.respond(msg.ID, map[string]any{})
	case "session/close":
		a.cancel("echo-session")
		a.logf("close echo-session")
		a.respond(msg.ID, map[string]any{})
		go func() {
			time.Sleep(25 * time.Millisecond)
			os.Exit(0)
		}()
	case "session/cancel":
		var params cancelParams
		_ = json.Unmarshal(msg.Params, &params)
		if strings.TrimSpace(params.SessionID) == "" {
			params.SessionID = "echo-session"
		}
		a.cancel(params.SessionID)
		a.logf("cancel %s", params.SessionID)
	default:
		a.respond(msg.ID, map[string]any{})
	}
}

func (a *echoAgent) handlePrompt(msg rpcMessage) {
	var params promptParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		a.respondError(msg.ID, -32602, err.Error())
		return
	}
	sessionID := strings.TrimSpace(params.SessionID)
	if sessionID == "" {
		sessionID = "echo-session"
	}
	text := promptText(params.Prompt)
	cancel := a.startPrompt(sessionID)
	defer a.finishPrompt(sessionID, cancel)
	switch {
	case strings.Contains(text, "block-until-cancel"):
		a.notify(sessionID, "prompt-blocked")
		select {
		case <-cancel:
			a.respond(msg.ID, map[string]any{"stopReason": "cancelled"})
		case <-time.After(30 * time.Second):
			a.respond(msg.ID, map[string]any{"stopReason": "end_turn"})
		}
	case strings.Contains(text, "crash-adapter"):
		a.logf("crash %s", sessionID)
		os.Exit(42)
	default:
		a.notify(sessionID, "echo: "+text)
		a.respond(msg.ID, map[string]any{"stopReason": "end_turn"})
	}
}

func promptText(blocks []struct {
	Type string ` + "`json:\"type\"`" + `
	Text string ` + "`json:\"text\"`" + `
}) string {
	var builder strings.Builder
	for _, block := range blocks {
		if block.Type == "text" {
			builder.WriteString(block.Text)
		}
	}
	return builder.String()
}

func (a *echoAgent) startPrompt(sessionID string) chan struct{} {
	ch := make(chan struct{})
	a.mu.Lock()
	a.cancels[sessionID] = ch
	a.mu.Unlock()
	return ch
}

func (a *echoAgent) finishPrompt(sessionID string, ch chan struct{}) {
	a.mu.Lock()
	if a.cancels[sessionID] == ch {
		delete(a.cancels, sessionID)
	}
	a.mu.Unlock()
}

func (a *echoAgent) cancel(sessionID string) {
	a.mu.Lock()
	ch := a.cancels[sessionID]
	delete(a.cancels, sessionID)
	a.mu.Unlock()
	if ch != nil {
		close(ch)
	}
}

func (a *echoAgent) notify(sessionID string, text string) {
	a.write(map[string]any{
		"jsonrpc": "2.0",
		"method":  "session/update",
		"params": map[string]any{
			"sessionId": sessionID,
			"update": map[string]any{
				"sessionUpdate": "agent_message_chunk",
				"content": map[string]any{
					"type": "text",
					"text": text,
				},
			},
		},
	})
}

func (a *echoAgent) respond(id json.RawMessage, result any) {
	if len(id) == 0 {
		return
	}
	a.write(struct {
		JSONRPC string          ` + "`json:\"jsonrpc\"`" + `
		ID      json.RawMessage ` + "`json:\"id\"`" + `
		Result  any             ` + "`json:\"result\"`" + `
	}{JSONRPC: "2.0", ID: id, Result: result})
}

func (a *echoAgent) respondError(id json.RawMessage, code int, message string) {
	if len(id) == 0 {
		return
	}
	a.write(struct {
		JSONRPC string          ` + "`json:\"jsonrpc\"`" + `
		ID      json.RawMessage ` + "`json:\"id\"`" + `
		Error   map[string]any  ` + "`json:\"error\"`" + `
	}{JSONRPC: "2.0", ID: id, Error: map[string]any{"code": code, "message": message}})
}

func (a *echoAgent) write(value any) {
	a.writeMu.Lock()
	defer a.writeMu.Unlock()
	_ = json.NewEncoder(a.writer).Encode(value)
	_ = a.writer.Flush()
}

func (a *echoAgent) logf(format string, args ...any) {
	a.logMu.Lock()
	defer a.logMu.Unlock()
	_, _ = fmt.Fprintf(a.logFile, format+"\n", args...)
}

func (a *echoAgent) persistSession(sessionID string) {
	if path := os.Getenv("ROUNDFIX_ACPX_ECHO_STATE"); path != "" {
		_ = os.WriteFile(path, []byte(sessionID), 0o644)
	}
}
`
