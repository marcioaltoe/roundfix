package cli

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
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/codex"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
	runworktree "roundfix/internal/worktree"
)

const implementTestSlug = "0001-widget-flow"

const cliTestHelperEnv = "ROUNDFIX_CLI_TEST_HELPER"
const detachTestChildModeEnv = "ROUNDFIX_DETACH_TEST_CHILD"

func TestMain(m *testing.M) {
	if mode := os.Getenv(detachTestChildModeEnv); mode != "" {
		os.Exit(runDetachTestChild(mode))
	}
	if os.Getenv(cliTestHelperEnv) == "1" {
		os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
	}
	os.Exit(m.Run())
}

// implementSeed describes one task file for a test Spec directory. Zero
// values default to status pending, type backend, and one passing
// Verification command.
type implementSeed struct {
	id           string
	title        string
	taskType     string
	status       string
	needs        []string
	verification []string
}

func gitImplement(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmdArgs := append(gitConfigArgsForTest(), args...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Dir = dir
	cmd.Env = isolatedGitEnvForTest()
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output.String())
	}
}

func gitImplementOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmdArgs := append(gitConfigArgsForTest(), args...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Dir = dir
	cmd.Env = isolatedGitEnvForTest()
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output.String())
	}
	return output.String()
}

func runCLIHelper(t *testing.T, dir string, fakeACPX string, extraEnv map[string]string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(os.Args[0], args...)
	cmd.Dir = dir
	cmd.Env = cliHelperEnv(fakeACPX, extraEnv)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), exitCodeFromWait(err)
}

func cliHelperEnv(fakeACPX string, extra map[string]string) []string {
	env := isolatedGitEnvForTest()
	env = withEnvValue(env, cliTestHelperEnv, "1")
	if fakeACPX != "" {
		env = withEnvValue(env, "PATH", filepath.Dir(fakeACPX)+string(os.PathListSeparator)+os.Getenv("PATH"))
	}
	for key, value := range extra {
		env = withEnvValue(env, key, value)
	}
	return env
}

func withEnvValue(env []string, key string, value string) []string {
	prefix := key + "="
	filtered := make([]string, 0, len(env)+1)
	for _, entry := range env {
		if strings.HasPrefix(entry, prefix) {
			continue
		}
		filtered = append(filtered, entry)
	}
	return append(filtered, prefix+value)
}

func fakeACPXCommand(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "acpx")
	body := fmt.Sprintf(`#!/bin/sh
if [ "$1" = "--version" ]; then
  printf '%%s\n' '%s'
  exit 0
fi
case " $* " in
  *" sessions ensure "*)
    model=""
    session=""
    previous=""
    for argument in "$@"; do
      if [ "$previous" = "--model" ]; then
        model="$argument"
      elif [ "$previous" = "--name" ]; then
        session="$argument"
      fi
      previous="$argument"
    done
    printf '%%s' "$model" > "$0.$session.model"
    exit 0
    ;;
  *" sessions show "*)
    session=""
    for argument in "$@"; do
      session="$argument"
    done
    model=$(cat "$0.$session.model")
    printf '{"schema":"acpx.session.v1","acpx":{"current_model_id":"%%s","config_options":[{"id":"model","category":"model","type":"select","currentValue":"%%s","options":[{"value":"%%s"}]},{"id":"reasoning_effort","type":"select","currentValue":"medium","options":[{"value":"low"},{"value":"medium"},{"value":"high"},{"value":"xhigh"},{"value":"max"},{"value":"maximum"},{"value":"ultra"}]}]}}\n' "$model" "$model" "$model"
    exit 0
    ;;
  *" set model "*|*" set reasoning_effort "*|*" set effort "*)
    config_id=""
    config_value=""
    session=""
    previous=""
    for argument in "$@"; do
      if [ "$previous" = "set" ]; then
        config_id="$argument"
      elif [ -n "$config_id" ] && [ -z "$config_value" ]; then
        config_value="$argument"
      elif [ "$previous" = "-s" ]; then
        session="$argument"
      fi
      previous="$argument"
    done
    state_path="$0.$session.model"
    model="$config_value"
    if [ "$config_id" = "model" ]; then
      printf '%%s' "$model" > "$state_path"
      current_reasoning="medium"
    else
      model=$(cat "$state_path")
      current_reasoning="$config_value"
    fi
    printf '{"action":"config_set","configId":"%%s","value":"%%s","configOptions":[{"id":"model","category":"model","type":"select","currentValue":"%%s","options":[{"value":"%%s"}]},{"id":"reasoning_effort","type":"select","currentValue":"%%s","options":[{"value":"low"},{"value":"medium"},{"value":"high"},{"value":"xhigh"},{"value":"max"},{"value":"maximum"},{"value":"ultra"}]}]}\n' "$config_id" "$config_value" "$model" "$model" "$current_reasoning"
    exit 0
    ;;
  *" sessions close "*)
    exit 0
    ;;
  *" prompt "*)
    cat >/dev/null
    if [ -n "$ROUNDFIX_FAKE_ACPX_PROMPT_STARTED" ]; then
      : > "$ROUNDFIX_FAKE_ACPX_PROMPT_STARTED"
    fi
    if [ -n "$ROUNDFIX_FAKE_ACPX_RELEASE" ]; then
      while [ ! -f "$ROUNDFIX_FAKE_ACPX_RELEASE" ]; do
        sleep 0.05
      done
    fi
    printf '{"jsonrpc":"2.0","id":1,"result":{"stopReason":"end_turn"}}\n'
    exit 0
    ;;
esac
exit 0
`, agent.MinimumACPXVersion)
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatalf("write fake acpx: %v", err)
	}
	// The adapter preflight probe LookPaths the agent-command binary; ship a
	// fake codex-acp next to the fake acpx so detach children pass preflight
	// on machines without the real adapter installed (CI runners).
	adapterPath := filepath.Join(filepath.Dir(path), "codex-acp")
	if err := os.WriteFile(adapterPath, []byte("#!/bin/sh\nprintf '%s\\n' '@agentclientprotocol/codex-acp "+agent.PinnedCodexAdapterVersion+"'\n"), 0o755); err != nil {
		t.Fatalf("write fake codex-acp: %v", err)
	}
	return path
}

func parseDetachedReport(t *testing.T, stdout string) (string, string) {
	t.Helper()
	lines := strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")
	if len(lines) != 5 {
		t.Fatalf("expected five detach stdout lines, got %d: %q", len(lines), stdout)
	}
	runID, ok := strings.CutPrefix(lines[0], "Run ID: ")
	if !ok || strings.TrimSpace(runID) == "" {
		t.Fatalf("expected Run ID line, got %q", lines[0])
	}
	consoleLog, ok := strings.CutPrefix(lines[1], "Console Log: ")
	if !ok || strings.TrimSpace(consoleLog) == "" {
		t.Fatalf("expected Console Log line, got %q", lines[1])
	}
	if lines[2] != "Attach: roundfix attach "+runID {
		t.Fatalf("expected Attach line for %s, got %q", runID, lines[2])
	}
	if lines[3] != "Supervisor monitor: roundfix events "+runID+" --follow --filter outcome" {
		t.Fatalf("expected Supervisor monitor line for %s, got %q", runID, lines[3])
	}
	if lines[4] != "Stop: roundfix stop "+runID {
		t.Fatalf("expected stop line for %s, got %q", runID, lines[4])
	}
	return runID, consoleLog
}

func readLineWithTimeout(t *testing.T, reader *bufio.Reader, timeout time.Duration) string {
	t.Helper()
	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)
	go func() {
		line, err := reader.ReadString('\n')
		if err != nil {
			errCh <- err
			return
		}
		lineCh <- line
	}()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case line := <-lineCh:
		return line
	case err := <-errCh:
		t.Fatalf("read line: %v", err)
	case <-timer.C:
		t.Fatalf("timed out waiting for line")
	}
	return ""
}

func waitProcessForTest(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
	}()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case err := <-waitCh:
		return err, true
	case <-timer.C:
		return nil, false
	}
}

func waitForFile(t *testing.T, path string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", path)
}

// waitForFileContains polls a file until it contains substr. A detached child's
// console log is flushed a hair after the store State flips, so reading it once
// right after the state change races under load.
func waitForFileContains(t *testing.T, path string, substr string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	var last string
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(path)
		if err == nil {
			last = string(data)
			if strings.Contains(last, substr) {
				return
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s to contain %q; last content %q", path, substr, last)
}

func waitForRunState(t *testing.T, homeDir string, runID string, state string, timeout time.Duration) store.Run {
	t.Helper()
	deadline := time.Now().Add(timeout)
	var last store.Run
	for time.Now().Before(deadline) {
		reader, err := store.OpenReader(context.Background(), homeDir)
		if err == nil {
			run, found, runErr := reader.Run(context.Background(), runID)
			_ = reader.Close()
			if runErr != nil {
				t.Fatalf("read Run %s: %v", runID, runErr)
			}
			if found {
				last = run
				if run.State == state {
					return run
				}
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for Run %s state %s; last state %s", runID, state, last.State)
	return store.Run{}
}

// waitForCleanOutcomeEvent polls the Run Event Journal until the Clean Daemon
// outcome and its notification receipt are both visible. The store State can
// flip to Clean a hair before those durable events are appended, so a single
// snapshot races under load; polling removes the race without masking a
// product bug.
func waitForCleanOutcomeEvent(t *testing.T, homeDir string, runID string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		events := runEventsForRun(t, homeDir, runID)
		if len(events) >= 2 {
			outcome := events[len(events)-2].Event
			receipt := events[len(events)-1].Event
			var payload runevent.NotificationReceiptPayload
			if outcome.Kind == runevent.KindDaemonOutcome &&
				strings.Contains(string(outcome.Payload), `"Clean"`) &&
				receipt.Kind == runevent.KindDaemonStatus &&
				json.Unmarshal(receipt.Payload, &payload) == nil &&
				strings.HasPrefix(payload.Event, "outcome_notification_") {
				return
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	events := runEventsForRun(t, homeDir, runID)
	if len(events) == 0 {
		t.Fatalf("timed out waiting for journaled Clean outcome for Run %s; no events", runID)
	}
	last := events[len(events)-1].Event
	t.Fatalf("timed out waiting for journaled Clean outcome for Run %s; last kind=%s payload=%s", runID, last.Kind, string(last.Payload))
}

func runEventsForRun(t *testing.T, homeDir string, runID string) []store.JournalEvent {
	t.Helper()
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database reader: %v", err)
	}
	defer func() {
		_ = reader.Close()
	}()
	events, err := reader.RunEventsAfter(context.Background(), runID, 0, 200)
	if err != nil {
		t.Fatalf("read Run Events for %s: %v", runID, err)
	}
	return events
}

func gitConfigArgsForTest() []string {
	return []string{
		"-c", "user.name=Roundfix Test",
		"-c", "user.email=test@example.com",
		"-c", "commit.gpgsign=false",
	}
}

func isolatedGitEnvForTest() []string {
	env := make([]string, 0, len(os.Environ())+2)
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(key, "GIT_CONFIG_") {
			continue
		}
		env = append(env, entry)
	}
	return append(env, "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")
}

// newImplementWorkspace builds a real git repository containing a committed
// Spec directory, on a non-default work branch, with HOME pointed at a temp
// dir and the test chdir'd into the repository. It returns the temp home and
// the symlink-resolved repository root (the root git itself reports).
func newImplementWorkspace(t *testing.T, seeds []implementSeed) (string, string) {
	t.Helper()
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	repoDir := t.TempDir()
	gitImplement(t, repoDir, "init", "--initial-branch=main")
	gitImplement(t, repoDir, "config", "user.name", "Roundfix Test")
	gitImplement(t, repoDir, "config", "user.email", "roundfix-test@example.com")
	gitImplement(t, repoDir, "config", "commit.gpgsign", "false")
	writeImplementSpec(t, repoDir, implementTestSlug, seeds)
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "seed spec")
	gitImplement(t, repoDir, "checkout", "-b", "ma/widget-flow")
	t.Chdir(repoDir)
	resolved, err := filepath.EvalSymlinks(repoDir)
	if err != nil {
		t.Fatalf("resolve repo dir: %v", err)
	}
	return homeDir, resolved
}

func TestLoadCommittedSpecGraphIgnoresDirtyCheckoutTaskMetadata(t *testing.T) {
	_, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", taskType: "backend", status: string(spec.StatusPending)}})
	taskPath := implementTaskPath(repoDir, "task_01")
	dirty := strings.ReplaceAll(mustRead(t, taskPath), "status: pending\ntype: backend", "status: completed\ntype: frontend")
	mustWrite(t, taskPath, dirty)
	head := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))

	graph, _, err := defaultLoadCommittedSpecGraph(context.Background(), repoDir, roundconfig.SpecsRoot{Path: filepath.Join(repoDir, "docs", "specs")}, head, implementTestSlug)
	if err != nil {
		t.Fatalf("load committed graph: %v", err)
	}
	if len(graph.Tasks) != 1 {
		t.Fatalf("len(tasks) = %d, want 1", len(graph.Tasks))
	}
	if graph.Tasks[0].Type != spec.TaskTypeBackend || graph.Tasks[0].Status != spec.StatusPending {
		t.Fatalf("committed graph used dirty task metadata: %+v", graph.Tasks[0])
	}
}

func TestLoadCommittedExternalSpecGraphIgnoresDirtyCheckoutTaskMetadata(t *testing.T) {
	_, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	_, externalRoot := newExternalSpecsRoot(t, implementTestSlug, []implementSeed{{id: "task_01", taskType: "backend", status: string(spec.StatusPending)}})
	taskPath := implementTaskPathInRoot(externalRoot, implementTestSlug, "task_01")
	dirty := strings.ReplaceAll(mustRead(t, taskPath), "status: pending\ntype: backend", "status: completed\ntype: frontend")
	mustWrite(t, taskPath, dirty)

	graph, _, err := defaultLoadCommittedSpecGraph(context.Background(), repoDir, roundconfig.SpecsRoot{Path: externalRoot, External: true}, "project-revision-is-ignored", implementTestSlug)
	if err != nil {
		t.Fatalf("load committed external graph: %v", err)
	}
	if len(graph.Tasks) != 1 {
		t.Fatalf("len(tasks) = %d, want 1", len(graph.Tasks))
	}
	if graph.Tasks[0].Type != spec.TaskTypeBackend || graph.Tasks[0].Status != spec.StatusPending {
		t.Fatalf("external committed graph used dirty task metadata: %+v", graph.Tasks[0])
	}
}

func TestLoadCommittedExternalSpecGraphRequiresExternalGitRepository(t *testing.T) {
	_, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	externalRoot := filepath.Join(t.TempDir(), "external-specs")
	writeImplementSpecAtRoot(t, externalRoot, implementTestSlug, []implementSeed{{id: "task_01"}})

	_, _, err := defaultLoadCommittedSpecGraph(context.Background(), repoDir, roundconfig.SpecsRoot{Path: externalRoot, External: true}, "HEAD", implementTestSlug)
	if err == nil {
		t.Fatal("expected external specs.root without a Git repository to fail")
	}
	if !strings.Contains(err.Error(), "must be committed in its own Git repository") {
		t.Fatalf("expected actionable external Git repository error, got %v", err)
	}
}

func writeUserConfig(t *testing.T, homeDir string, content string) {
	t.Helper()
	path := filepath.Join(homeDir, ".roundfix", "config.yml")
	mustMkdir(t, filepath.Dir(path))
	mustWrite(t, path, content)
}

func configureImplementAutoPush(t *testing.T, repoDir string, enabled bool) {
	t.Helper()
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), fmt.Sprintf("implement:\n  auto_push: %t\n", enabled))
	gitImplement(t, repoDir, "add", ".roundfixrc.yml")
	gitImplement(t, repoDir, "commit", "-m", "configure implement auto push")
}

func configureImplementClaudeReasoning(t *testing.T, repoDir string) {
	t.Helper()
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), "runtimes:\n  claude:\n    reasoning_effort: high\n")
	gitImplement(t, repoDir, "add", ".roundfixrc.yml")
	gitImplement(t, repoDir, "commit", "-m", "configure claude reasoning")
}

func configureExternalSpecsRoot(t *testing.T, repoDir string, specsRoot string) {
	t.Helper()
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), fmt.Sprintf("specs:\n  root: %q\n", specsRoot))
	gitImplement(t, repoDir, "add", ".roundfixrc.yml")
	gitImplement(t, repoDir, "commit", "-m", "configure external Spec Root")
}

func newExternalSpecsRoot(t *testing.T, slug string, seeds []implementSeed) (string, string) {
	t.Helper()
	externalRepo := t.TempDir()
	gitImplement(t, externalRepo, "init", "--initial-branch=main")
	specsRoot := filepath.Join(externalRepo, "docs", "specs")
	writeImplementSpecAtRoot(t, specsRoot, slug, seeds)
	gitImplement(t, externalRepo, "add", "-A")
	gitImplement(t, externalRepo, "commit", "-m", "seed external spec")
	return externalRepo, specsRoot
}

func configureWorktreeBootstrap(t *testing.T, repoDir string, command string, timeout string) {
	t.Helper()
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), fmt.Sprintf("worktree:\n  bootstrap: %q\n  bootstrap_timeout: %s\n", command, timeout))
	gitImplement(t, repoDir, "add", ".roundfixrc.yml")
	gitImplement(t, repoDir, "commit", "-m", "configure worktree bootstrap")
}

func configureImplementUpstream(t *testing.T, repoDir string, remote string, branch string) {
	t.Helper()
	remoteDir := filepath.Join(t.TempDir(), remote+".git")
	mustMkdir(t, remoteDir)
	gitImplement(t, remoteDir, "init", "--bare")
	gitImplement(t, repoDir, "remote", "add", remote, remoteDir)
	gitImplement(t, repoDir, "push", "-u", remote, "HEAD:"+branch)
}

func writeImplementSpec(t *testing.T, repoDir string, slug string, seeds []implementSeed) {
	t.Helper()
	writeImplementSpecAtRoot(t, filepath.Join(repoDir, "docs", "specs"), slug, seeds)
}

func writeImplementSpecAtRoot(t *testing.T, specsRoot string, slug string, seeds []implementSeed) {
	t.Helper()
	specDir := filepath.Join(specsRoot, slug)
	mustMkdir(t, specDir)
	mustWrite(t, filepath.Join(specDir, "_prd.md"), "---\nstatus: active\n---\n\n# PRD\n")

	var manifest strings.Builder
	manifest.WriteString("---\nschema: spec-tasks/v1\ngraph:\n  nodes:\n")
	for _, seed := range seeds {
		manifest.WriteString(fmt.Sprintf("    - id: %s\n      file: %s.md\n", seed.id, seed.id))
		if len(seed.needs) > 0 {
			manifest.WriteString(fmt.Sprintf("      needs: [%s]\n", strings.Join(seed.needs, ", ")))
		}
	}
	manifest.WriteString("---\n\n# Task Graph\n")
	mustWrite(t, filepath.Join(specDir, "_tasks.md"), manifest.String())

	for _, seed := range seeds {
		mustWrite(t, implementTaskPathInRoot(specsRoot, slug, seed.id), implementTaskContent(slug, seed))
	}
}

func implementTaskContent(slug string, seed implementSeed) string {
	status := seed.status
	if status == "" {
		status = string(spec.StatusPending)
	}
	taskType := seed.taskType
	if taskType == "" {
		taskType = "backend"
	}
	title := seed.title
	if title == "" {
		title = "Do the " + seed.id + " work"
	}
	verification := seed.verification
	if len(verification) == 0 {
		verification = []string{"true"}
	}
	var body strings.Builder
	body.WriteString(fmt.Sprintf("---\ntask: %s\nspec: %s\nstatus: %s\ntype: %s\n---\n\n# %s\n\n## Verification\n\n", seed.id, slug, status, taskType, title))
	for _, command := range verification {
		body.WriteString(fmt.Sprintf("- `%s` — expected: passes.\n", command))
	}
	return body.String()
}

func implementTaskPath(repoDir string, taskID string) string {
	return implementTaskPathInRoot(filepath.Join(repoDir, "docs", "specs"), implementTestSlug, taskID)
}

func implementTaskPathInRoot(specsRoot string, slug string, taskID string) string {
	return filepath.Join(specsRoot, slug, taskID+".md")
}

// implementFakeRunner scripts per-Task Agent behavior keyed by the Task id
// parsed from the prompt, writing task statuses through spec.SetStatus the
// way a real Agent edits the task file. A QA prompt writes qaReport as the
// Spec's QA Report, the way the qa-gate Agent does; an empty qaReport
// writes none.
type implementFakeRunner struct {
	mu            sync.Mutex
	gitRoot       string
	statusByTask  map[string]spec.Status
	errByTask     map[string]error
	probeErr      error
	probeRequests []agent.ProbeRequest
	fallback      agent.FallbackSelection
	fallbackOK    bool
	fallbackErr   error
	fallbackSets  []agent.FallbackCandidateSet
	calls         int
	taskIDs       []string
	qaReport      string
	qaCalls       int
	logPaths      []string
	writeLogs     bool
	agentOutput   string
	onTask        func(req agent.ExecuteRequest, taskID string) error
}

func (runner *implementFakeRunner) Probe(_ context.Context, req agent.ProbeRequest) error {
	runner.probeRequests = append(runner.probeRequests, req)
	return runner.probeErr
}

func (runner *implementFakeRunner) ProbeFallback(_ context.Context, _ agent.RuntimeSpec, candidates agent.FallbackCandidateSet) (agent.FallbackSelection, bool, error) {
	runner.fallbackSets = append(runner.fallbackSets, candidates)
	return runner.fallback, runner.fallbackOK, runner.fallbackErr
}

func (runner *implementFakeRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	runner.mu.Lock()
	runner.calls++
	runner.logPaths = append(runner.logPaths, req.LogPath)
	agentOutput := runner.agentOutput
	writeLogs := runner.writeLogs
	errByTask := runner.errByTask
	onTask := runner.onTask
	statusByTask := runner.statusByTask
	qaReport := runner.qaReport
	runner.mu.Unlock()

	executionRoot := strings.TrimSpace(req.GitRoot)
	if executionRoot == "" {
		executionRoot = runner.gitRoot
	}
	if agentOutput != "" {
		if err := publishFakeAgentOutput(ctx, sink, req, agentOutput); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if writeLogs && strings.TrimSpace(req.LogPath) != "" {
		if err := os.MkdirAll(filepath.Dir(req.LogPath), 0o755); err != nil {
			return agent.ExecuteResult{}, err
		}
		if err := os.WriteFile(req.LogPath, []byte("fake agent output\n"), 0o644); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	taskID := implementTaskIDFromPrompt(req.Prompt)
	if taskID == "" && strings.Contains(req.Prompt, "Spec QA gate") {
		runner.mu.Lock()
		runner.qaCalls++
		runner.mu.Unlock()
		if qaReport != "" {
			reportPath := filepath.Join(implementSpecDirFromPrompt(req.Prompt, executionRoot), "qa", implementQAReportName)
			if err := os.MkdirAll(filepath.Dir(reportPath), 0o755); err != nil {
				return agent.ExecuteResult{}, err
			}
			if err := os.WriteFile(reportPath, []byte(qaReport), 0o644); err != nil {
				return agent.ExecuteResult{}, err
			}
		}
		return agent.ExecuteResult{LogPath: req.LogPath}, nil
	}
	runner.mu.Lock()
	runner.taskIDs = append(runner.taskIDs, taskID)
	runner.mu.Unlock()
	if err := errByTask[taskID]; err != nil {
		return agent.ExecuteResult{}, err
	}
	if onTask != nil {
		if err := onTask(req, taskID); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if status, ok := statusByTask[taskID]; ok {
		if err := spec.SetStatus(implementTaskPathFromPrompt(req.Prompt, executionRoot, taskID), status); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	return agent.ExecuteResult{LogPath: req.LogPath}, nil
}

func (runner *implementFakeRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

const implementQAReportName = "qa-report-2026-01-01.md"

func implementQAReportRelPath() string {
	return filepath.Join("docs", "specs", implementTestSlug, "qa", implementQAReportName)
}

func implementQAReport(verdict string) string {
	return fmt.Sprintf("---\nverdict: %s\n---\n\n# QA Report\n", verdict)
}

func implementTaskIDFromPrompt(prompt string) string {
	for _, line := range strings.Split(prompt, "\n") {
		if strings.HasPrefix(line, "Task: ") {
			return strings.TrimPrefix(line, "Task: ")
		}
	}
	return ""
}

func implementTaskPathFromPrompt(prompt string, executionRoot string, taskID string) string {
	for _, line := range strings.Split(prompt, "\n") {
		path, ok := strings.CutPrefix(line, "Task file: ")
		if !ok {
			continue
		}
		path = strings.TrimSpace(path)
		if filepath.IsAbs(path) {
			return path
		}
		return filepath.Join(executionRoot, path)
	}
	return implementTaskPath(executionRoot, taskID)
}

func implementSpecDirFromPrompt(prompt string, executionRoot string) string {
	for _, line := range strings.Split(prompt, "\n") {
		path, ok := strings.CutPrefix(line, "Spec directory: ")
		if !ok {
			continue
		}
		path = strings.TrimSpace(path)
		if filepath.IsAbs(path) {
			return path
		}
		return filepath.Join(executionRoot, path)
	}
	return filepath.Join(executionRoot, "docs", "specs", implementTestSlug)
}

// withImplementCollaborators wires the standard fake collaborators for an
// implement run and returns the ones tests assert on.
func withImplementCollaborators(t *testing.T, runner agent.Runner) (*fakeCommitter, *fakeVerifier, *fakePusher, *fakeSourceResolver) {
	t.Helper()
	withFakeRunWorktrees(t)
	committer := &fakeCommitter{}
	verifier := &fakeVerifier{}
	pusher := &fakePusher{}
	sourceResolver := &fakeSourceResolver{}
	withAgentRunner(t, runner)
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withFakeWorktree(t)
	withPriorChangedResolver(t, emptyPriorChangedResolver{})
	withPusher(t, pusher)
	withSourceResolver(t, sourceResolver)
	return committer, verifier, pusher, sourceResolver
}

type emptyPriorChangedResolver struct{}

func (emptyPriorChangedResolver) PriorChangedFiles(context.Context, string, string) ([]string, error) {
	return nil, nil
}

func withPriorChangedResolver(t *testing.T, resolver daemon.PriorChangedResolver) {
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.priorChanges = resolver
	})
}

func runCleanImplementForCleanup(t *testing.T, cleanupErr error) (string, string, int, string, string) {
	t.Helper()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
	})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withImplementCollaborators(t, runner)
	cleanupPath := ""
	if cleanupErr != nil {
		oldCleanup := cleanupCleanRunWorktree
		cleanupCleanRunWorktree = func(_ context.Context, ref runworktree.Ref) error {
			cleanupPath = ref.Path
			return cleanupErr
		}
		t.Cleanup(func() {
			cleanupCleanRunWorktree = oldCleanup
		})
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean implement exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if cleanupErr != nil && strings.TrimSpace(cleanupPath) == "" {
		t.Fatal("expected cleanup failure to capture the Run Worktree path")
	}
	return stdout.String(), stderr.String(), code, homeDir, cleanupPath
}

func withFakeRunWorktrees(t *testing.T) {
	t.Helper()
	oldCreate := createRunWorktree
	oldIntegrate := integrateRunWorktree
	oldCleanup := cleanupCleanRunWorktree
	oldPrune := pruneTerminalRunWorktrees
	createRunWorktree = func(_ context.Context, opts runworktree.CreateOptions) (runworktree.Ref, error) {
		userRoot := opts.UserRoot
		runID := opts.RunID
		path := filepath.Join(os.Getenv("HOME"), ".roundfix", "worktrees", "fake", runID)
		if err := copyDir(filepath.Join(userRoot, "docs"), filepath.Join(path, "docs")); err != nil {
			return runworktree.Ref{}, err
		}
		gitImplement(t, path, "init", "--initial-branch=main")
		gitImplement(t, path, "add", "-A")
		gitImplement(t, path, "commit", "-m", "seed fake run worktree")
		gitImplement(t, path, "branch", "-m", runworktree.BranchName(runID))
		return runworktree.Ref{
			RunID:    runID,
			Path:     path,
			Branch:   runworktree.BranchName(runID),
			UserRoot: userRoot,
		}, nil
	}
	integrateRunWorktree = func(_ context.Context, ref runworktree.Ref, _ string, _ string) (runworktree.IntegrationResult, error) {
		if err := copyDir(filepath.Join(ref.Path, "docs"), filepath.Join(ref.UserRoot, "docs")); err != nil {
			return runworktree.IntegrationResult{}, err
		}
		return runworktree.IntegrationResult{Mode: runworktree.ModeFastForwardMerge}, nil
	}
	cleanupCleanRunWorktree = func(_ context.Context, ref runworktree.Ref) error {
		return os.RemoveAll(ref.Path)
	}
	pruneTerminalRunWorktrees = func(context.Context, string, string, runworktree.TerminalRunReconciliationStore, runworktree.TerminalRunLookup) ([]runworktree.PrunedRef, error) {
		return nil, nil
	}
	t.Cleanup(func() {
		createRunWorktree = oldCreate
		integrateRunWorktree = oldIntegrate
		cleanupCleanRunWorktree = oldCleanup
		pruneTerminalRunWorktrees = oldPrune
	})
}

func copyDir(source string, destination string) error {
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", source)
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		destinationPath := filepath.Join(destination, entry.Name())
		if entry.IsDir() {
			if err := copyDir(sourcePath, destinationPath); err != nil {
				return err
			}
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			continue
		}
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(destinationPath, content, info.Mode().Perm()); err != nil {
			return err
		}
	}
	return nil
}

func implementRunIDFromStderr(t *testing.T, stderr string) string {
	t.Helper()
	for _, line := range strings.Split(stderr, "\n") {
		if strings.HasPrefix(line, "Implement Run: ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Implement Run: "))
		}
	}
	t.Fatalf("no Implement Run id in stderr: %q", stderr)
	return ""
}

func implementRunFromStore(t *testing.T, homeDir string, runID string) store.Run {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	run, found, err := runStore.Run(ctx, runID)
	if err != nil {
		t.Fatalf("read run %s: %v", runID, err)
	}
	if !found {
		t.Fatalf("run %s not found in Run Database", runID)
	}
	return run
}

func onlyRunFromStore(t *testing.T, homeDir string) store.Run {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	runIDs, err := runStore.RunIDs(ctx)
	if err != nil {
		t.Fatalf("list runs: %v", err)
	}
	if len(runIDs) != 1 {
		t.Fatalf("expected one Run, got %d: %v", len(runIDs), runIDs)
	}
	run, found, err := runStore.Run(ctx, runIDs[0])
	if err != nil {
		t.Fatalf("read run %s: %v", runIDs[0], err)
	}
	if !found {
		t.Fatalf("run %s not found in Run Database", runIDs[0])
	}
	return run
}

func assertNoActiveRunInGitRoot(t *testing.T, homeDir string, gitRoot string) {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	if blocking, found, err := runStore.ActiveRunInGitRoot(ctx, gitRoot); err != nil {
		t.Fatalf("lookup active run in git root: %v", err)
	} else if found {
		t.Fatalf("expected the Active Run lock released, got %#v", blocking)
	}
}

func assertRunWorktreeExists(t *testing.T, path string) {
	t.Helper()
	if strings.TrimSpace(path) == "" {
		t.Fatal("expected Run Worktree path to be recorded")
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected Run Worktree %q to exist: %v", path, err)
	}
	if !info.IsDir() {
		t.Fatalf("expected Run Worktree %q to be a directory", path)
	}
}

func assertRunWorktreeRemoved(t *testing.T, path string) {
	t.Helper()
	if strings.TrimSpace(path) == "" {
		t.Fatal("expected Run Worktree path to be recorded")
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected Run Worktree %q to be removed, stat error %v", path, err)
	}
}

func assertRunBranchExists(t *testing.T, repoDir string, branch string) {
	t.Helper()
	if strings.TrimSpace(gitImplementOutput(t, repoDir, "branch", "--list", branch)) == "" {
		t.Fatalf("expected Run Branch %q to exist", branch)
	}
}

func assertRunBranchRemoved(t *testing.T, repoDir string, branch string) {
	t.Helper()
	if got := strings.TrimSpace(gitImplementOutput(t, repoDir, "branch", "--list", branch)); got != "" {
		t.Fatalf("expected Run Branch %q to be removed, got %q", branch, got)
	}
}

func integrationCommandFromStderr(t *testing.T, stderr string) string {
	t.Helper()
	for _, line := range strings.Split(stderr, "\n") {
		if command, ok := strings.CutPrefix(line, "Integration command: "); ok {
			command = strings.TrimSpace(command)
			if command == "" {
				t.Fatal("empty integration command")
			}
			return command
		}
	}
	t.Fatalf("expected integration command in stderr, got %q", stderr)
	return ""
}

func runPrintedIntegrationCommand(t *testing.T, repoDir string, command string) {
	t.Helper()
	fields := strings.Fields(command)
	if len(fields) < 2 || fields[0] != "git" {
		t.Fatalf("expected git integration command, got %q", command)
	}
	gitImplement(t, repoDir, fields[1:]...)
}

func TestRunImplementHelpListsExactlyImplementedFlags(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	want := map[string]bool{
		"--spec":              true,
		"--qa":                true,
		"--agent":             true,
		"--model":             true,
		"--reasoning-effort":  true,
		"--agent-command":     true,
		"--agent-full-access": true,
		"--detach":            true,
		"--interactive":       true,
		"--no-agent-console":  true,
		"--no-input":          true,
	}
	got := map[string]bool{}
	for _, flagName := range regexp.MustCompile(`--[a-z][a-z-]*`).FindAllString(stdout.String(), -1) {
		got[flagName] = true
	}
	// --agent-full-access also matches --agent; the set comparison below
	// checks exact flag vocabulary either way.
	for flagName := range want {
		if !got[flagName] {
			t.Fatalf("help is missing flag %s:\n%s", flagName, stdout.String())
		}
	}
	for flagName := range got {
		if !want[flagName] {
			t.Fatalf("help lists unimplemented flag %s:\n%s", flagName, stdout.String())
		}
	}
}

func TestRunImplementVerificationCapacityDoesNotAddFlag(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--help"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	if strings.Contains(stdout.String(), "--verification-concurrency") {
		t.Fatalf("Verification Capacity must remain Config-only:\n%s", stdout.String())
	}
}

func TestRunImplementDetachPrintsReportAndCompletesRun(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	fakeACPX := fakeACPXCommand(t)
	stdout, stderr, code := runCLIHelper(t, repoDir, fakeACPX, nil,
		"implement", "--spec", implementTestSlug, "--agent-command", "codex-acp --stdio", "--detach")

	if code != exitOK {
		t.Fatalf("expected detach caller exit 0, got %d stderr=%q stdout=%q", code, stderr, stdout)
	}
	if stderr != "" {
		t.Fatalf("expected detach caller stderr empty, got %q", stderr)
	}
	runID, consoleLog := parseDetachedReport(t, stdout)
	wantStdout := fmt.Sprintf(
		"Run ID: %s\nConsole Log: %s\nAttach: roundfix attach %s\nSupervisor monitor: roundfix events %s --follow --filter outcome\nStop: roundfix stop %s\n",
		runID,
		consoleLog,
		runID,
		runID,
		runID,
	)
	if stdout != wantStdout {
		t.Fatalf("detach stdout mismatch\nwant: %q\ngot:  %q", wantStdout, stdout)
	}

	run := waitForRunState(t, homeDir, runID, store.StateClean, 90*time.Second)
	if run.Kind != store.KindImplement {
		t.Fatalf("expected Implement Run, got %s", run.Kind)
	}
	if run.WorkDir == "" {
		t.Fatal("expected Run Worktree recorded on detached Run")
	}
	assertRunWorktreeRemoved(t, run.WorkDir)
	waitForFileContains(t, consoleLog, "Implement Run "+runID+" reached Clean", 90*time.Second)
	waitForCleanOutcomeEvent(t, homeDir, runID, 90*time.Second)
}

func TestRunImplementDetachReportsAndRelaysPreflightFailure(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	var foregroundStdout bytes.Buffer
	var foregroundStderr bytes.Buffer
	foregroundCode := RunContext(context.Background(), []string{"implement", "--spec", "9999-missing", "--no-input"}, &foregroundStdout, &foregroundStderr)
	if foregroundCode != exitPreflight {
		t.Fatalf("expected foreground preflight exit 2, got %d stderr=%q", foregroundCode, foregroundStderr.String())
	}

	stdout, stderr, code := runCLIHelper(t, repoDir, "", nil,
		"implement", "--spec", "9999-missing", "--detach")

	if code != exitPreflight {
		t.Fatalf("expected detached preflight exit 2, got %d stderr=%q stdout=%q", code, stderr, stdout)
	}
	if stdout != "" {
		t.Fatalf("expected detached preflight stdout empty, got %q", stdout)
	}
	if !strings.Contains(stderr, "Detached Run child failed Run creation handshake: EOF; child exited (exit status 2); console output follows") {
		t.Fatalf("expected detached preflight diagnostic, got %q", stderr)
	}
	if !strings.Contains(stderr, foregroundStderr.String()) {
		t.Fatalf("detached preflight stderr did not relay child output\nwant output containing: %q\ngot: %q", foregroundStderr.String(), stderr)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunImplementDetachSurvivesCallerProcessGroupKill(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	dir := t.TempDir()
	promptStarted := filepath.Join(dir, "prompt-started")
	releasePrompt := filepath.Join(dir, "release")
	t.Cleanup(func() {
		_ = os.WriteFile(releasePrompt, []byte("release\n"), 0o644)
	})
	fakeACPX := fakeACPXCommand(t)
	cmd := exec.Command(os.Args[0], "implement", "--spec", implementTestSlug, "--agent-command", "codex-acp --stdio", "--detach")
	cmd.Dir = repoDir
	cmd.Env = cliHelperEnv(fakeACPX, map[string]string{
		"ROUNDFIX_FAKE_ACPX_PROMPT_STARTED": promptStarted,
		"ROUNDFIX_FAKE_ACPX_RELEASE":        releasePrompt,
	})
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("open caller stdout pipe: %v", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start detach caller: %v", err)
	}

	firstLine := readLineWithTimeout(t, bufio.NewReader(stdoutPipe), 5*time.Second)
	runID, ok := strings.CutPrefix(strings.TrimSpace(firstLine), "Run ID: ")
	if !ok || strings.TrimSpace(runID) == "" {
		t.Fatalf("expected first detach line with Run id, got %q stderr=%q", firstLine, stderr.String())
	}
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil && !errors.Is(err, syscall.ESRCH) {
		t.Fatalf("kill caller process group: %v", err)
	}
	_, _ = waitProcessForTest(cmd, 2*time.Second)
	waitForFile(t, promptStarted, 60*time.Second)

	var attachStdout bytes.Buffer
	var attachStderr bytes.Buffer
	attachCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	attachCode := RunContext(attachCtx, []string{"attach", runID}, &attachStdout, &attachStderr)
	cancel()
	if attachCode != exitOK {
		t.Fatalf("expected attach to detach cleanly from active Run, got %d stderr=%q stdout=%q", attachCode, attachStderr.String(), attachStdout.String())
	}
	if !strings.Contains(attachStdout.String(), "Detached; Run "+runID+" keeps going.") {
		t.Fatalf("expected attach to prove active Run is attachable, got stdout=%q stderr=%q", attachStdout.String(), attachStderr.String())
	}

	mustWrite(t, releasePrompt, "release\n")
	run := waitForRunState(t, homeDir, runID, store.StateClean, 90*time.Second)
	if run.State != store.StateClean {
		t.Fatalf("expected detached child to reach Clean after caller kill, got %s", run.State)
	}
	waitForCleanOutcomeEvent(t, homeDir, runID, 90*time.Second)
}

func TestRunHelpListsImplementCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "roundfix implement --spec <slug>") {
		t.Fatalf("expected top-level usage to list implement, got %q", stdout.String())
	}
}

func TestRunImplementValidationFailures(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		message string
	}{
		{
			// The workspace has no docs/specs/, so Spec Root validation fails
			// before the Spec picker can offer choices.
			name:    "missing spec without active Specs",
			args:    []string{"implement"},
			message: "specs.root resolved to",
		},
		{
			name:    "missing spec with no-input",
			args:    []string{"implement", "--no-input"},
			message: "missing required --spec because --no-input disables Interactive Input",
		},
		{
			name:    "interactive without active Specs",
			args:    []string{"implement", "--interactive"},
			message: "specs.root resolved to",
		},
		{
			name:    "interactive with no-input",
			args:    []string{"implement", "--spec", implementTestSlug, "--interactive", "--no-input"},
			message: "--interactive cannot be used with --no-input",
		},
		{
			name:    "interactive with no-agent-console",
			args:    []string{"implement", "--spec", implementTestSlug, "--interactive", "--no-agent-console"},
			message: "--interactive cannot be used with --no-agent-console",
		},
		{
			name:    "explicitly empty agent with no-input",
			args:    []string{"implement", "--spec", implementTestSlug, "--agent=", "--model", "gpt-5.6-sol", "--reasoning-effort", "high", "--no-input"},
			message: "--agent cannot be empty",
		},
		{
			name:    "unsupported agent",
			args:    []string{"implement", "--spec", implementTestSlug, "--agent", "gemini", "--model", "gemini-pro", "--reasoning-effort", "high"},
			message: `unsupported Agent "gemini"`,
		},
		{
			name:    "unexpected argument",
			args:    []string{"implement", implementTestSlug},
			message: "unexpected argument",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, _ := withCLIWorkspace(t)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != 2 {
				t.Fatalf("expected exit code 2, got %d (stderr %q)", code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.message) {
				t.Fatalf("expected message %q, got %q", tt.message, stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunImplementRejectsInvalidVerificationCapacityBeforeRunCreation(t *testing.T) {
	tests := []struct {
		name     string
		config   string
		contains string
	}{
		{name: "zero", config: "verification:\n  concurrency: 0\n", contains: "verification.concurrency must be greater than 0"},
		{name: "negative", config: "verification:\n  concurrency: -1\n", contains: "verification.concurrency must be greater than 0"},
		{name: "non-integer", config: "verification:\n  concurrency: 1.5\n", contains: "verification.concurrency must be an integer"},
		{name: "unknown key", config: "verification:\n  concurrent: 1\n", contains: "verification.concurrent is not a supported config key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
			mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), tt.config)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

			if code != 2 {
				t.Fatalf("expected validation exit code 2, got %d stderr=%q", code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.contains) {
				t.Fatalf("expected stderr containing %q, got %q", tt.contains, stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunImplementInteractiveInputPicksSpecThroughCollector(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
	})
	mustMkdir(t, filepath.Join(repoDir, "docs", "specs", "0002-broken-prd"))
	mustWrite(t, filepath.Join(repoDir, "docs", "specs", "0002-broken-prd", "_prd.md"), "no frontmatter\n")
	gitImplement(t, repoDir, "add", "docs/specs/0002-broken-prd/_prd.md")
	gitImplement(t, repoDir, "commit", "-m", "add broken spec fixture")
	configureImplementClaudeReasoning(t, repoDir)
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withImplementCollaborators(t, runner)
	var inputReq roundtui.InputRequest
	var collected strings.Builder
	withInteractiveInput(t, func(ctx context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		inputReq = req
		// Drive the real collector synchronously: pick the first listed
		// Spec by number and override the Agent.
		return roundtui.CollectInput(ctx, req, strings.NewReader("1\nclaude\n\n\n"), &collected)
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr %q)", code, stderr.String())
	}
	if len(inputReq.SpecOptions) != 1 || inputReq.SpecOptions[0] != implementTestSlug {
		t.Fatalf("expected the picker to list the active Spec, got %#v", inputReq.SpecOptions)
	}
	for _, expected := range []string{
		"Active Specs:",
		"1. " + implementTestSlug,
		"Pick a Spec by number or slug.",
		"Spec []:",
		"Agent [codex]:",
	} {
		if !strings.Contains(collected.String(), expected) {
			t.Fatalf("expected the Spec picker to show %q, got:\n%s", expected, collected.String())
		}
	}
	if !strings.Contains(stderr.String(), "Interactive Input collected command parameters.") {
		t.Fatalf("expected Interactive Input confirmation, got %q", stderr.String())
	}
	diagnostic := "skipped docs/specs/0002-broken-prd: unreadable _prd.md frontmatter: missing YAML frontmatter opening marker"
	if !strings.Contains(stderr.String(), diagnostic) {
		t.Fatalf("expected skipped Spec diagnostic %q, got %q", diagnostic, stderr.String())
	}
	if strings.Contains(collected.String(), "skipped docs/specs") {
		t.Fatalf("expected picker rendering untouched by skipped diagnostics, got:\n%s", collected.String())
	}
	if !strings.Contains(stdout.String(), "task_01 completed — Build the widget core") {
		t.Fatalf("expected the Run to execute the picked Spec, got %q", stdout.String())
	}
	if strings.Contains(stdout.String(), "skipped docs/specs") {
		t.Fatalf("expected stdout untouched by skipped diagnostics, got %q", stdout.String())
	}
	run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if run.SpecSlug != implementTestSlug {
		t.Fatalf("expected the Run to target the picked Spec, got %q", run.SpecSlug)
	}
	defaults := readInteractiveDefaults(t, homeDir)
	if defaults.Agent != "claude" {
		t.Fatalf("expected the picked Agent remembered, got %#v", defaults)
	}
	if defaults.PRNumber != "" {
		t.Fatalf("expected no PR default from a spec Run, got %#v", defaults)
	}
}

func TestRunImplementUsesConfiguredExternalSpecRootEndToEnd(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Internal fixture should stay untouched"},
	})
	_, externalRoot := newExternalSpecsRoot(t, implementTestSlug, []implementSeed{
		{id: "task_01", title: "Build from external root"},
	})
	configureExternalSpecsRoot(t, repoDir, externalRoot)
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		qaReport:     implementQAReport(spec.VerdictPass),
	}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--qa", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected Clean exit, got %d (stderr %q stdout %q)", code, stderr.String(), stdout.String())
	}
	if !strings.Contains(stderr.String(), "Spec Root: "+externalRoot+"\n") {
		t.Fatalf("expected startup to name external Spec Root %q, got %q", externalRoot, stderr.String())
	}
	externalTask := implementTaskPathInRoot(externalRoot, implementTestSlug, "task_01")
	if content := mustRead(t, externalTask); !strings.Contains(content, "status: completed") {
		t.Fatalf("expected external task completed, got:\n%s", content)
	}
	if content := mustRead(t, implementTaskPath(repoDir, "task_01")); !strings.Contains(content, "status: pending") {
		t.Fatalf("expected default-layout fixture untouched, got:\n%s", content)
	}
	externalReport := filepath.Join(externalRoot, implementTestSlug, "qa", implementQAReportName)
	assertPathExists(t, externalReport)
	if !strings.Contains(stdout.String(), "qa pass — "+externalReport+"\n") {
		t.Fatalf("expected QA output to name external report, got %q", stdout.String())
	}
	run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if run.State != store.StateClean {
		t.Fatalf("expected external-root Run Clean, got %s", run.State)
	}
}

func TestRunImplementInteractiveInputListsConfiguredExternalSpecRoot(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Internal fixture should not be listed"},
	})
	externalSlug := "0002-external-only"
	_, externalRoot := newExternalSpecsRoot(t, externalSlug, []implementSeed{
		{id: "task_01", title: "Build from external picker"},
	})
	configureExternalSpecsRoot(t, repoDir, externalRoot)
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withImplementCollaborators(t, runner)
	var inputReq roundtui.InputRequest
	var collected strings.Builder
	withInteractiveInput(t, func(ctx context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		inputReq = req
		return roundtui.CollectInput(ctx, req, strings.NewReader("1\ncodex\ngpt-5.6-sol\nhigh\n"), &collected)
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected Clean exit, got %d (stderr %q stdout %q)", code, stderr.String(), stdout.String())
	}
	if len(inputReq.SpecOptions) != 1 || inputReq.SpecOptions[0] != externalSlug {
		t.Fatalf("expected picker to list only external Spec %q, got %#v", externalSlug, inputReq.SpecOptions)
	}
	if !strings.Contains(collected.String(), "1. "+externalSlug) {
		t.Fatalf("expected collected prompt to render external Spec, got:\n%s", collected.String())
	}
	externalTask := implementTaskPathInRoot(externalRoot, externalSlug, "task_01")
	if content := mustRead(t, externalTask); !strings.Contains(content, "status: completed") {
		t.Fatalf("expected selected external task completed, got:\n%s", content)
	}
	run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if run.SpecSlug != externalSlug {
		t.Fatalf("expected Run to target external Spec %q, got %q", externalSlug, run.SpecSlug)
	}
}

func TestRunImplementInteractiveInputMergesQAGateChoice(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		input       string
		qaReport    string
		wantQACalls int
	}{
		{name: "scripted yes produces QA Run", args: []string{"implement"}, input: "1\ncodex\ngpt-5.6-sol\nhigh\ny\n", qaReport: implementQAReport("pass"), wantQACalls: 1},
		{name: "empty input produces non-QA Run", args: []string{"implement"}, input: "1\ncodex\ngpt-5.6-sol\nhigh\n\n", wantQACalls: 0},
		{name: "qa flag preset keeps QA on with enter", args: []string{"implement", "--qa"}, input: "1\ncodex\ngpt-5.6-sol\nhigh\n\n", qaReport: implementQAReport("pass"), wantQACalls: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, repoDir := newImplementWorkspace(t, []implementSeed{
				{id: "task_01", title: "Build the widget core"},
			})
			runner := &implementFakeRunner{
				gitRoot:      repoDir,
				statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
				qaReport:     tt.qaReport,
			}
			withImplementCollaborators(t, runner)
			withInteractiveInput(t, func(ctx context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
				return roundtui.CollectInput(ctx, req, strings.NewReader(tt.input), io.Discard)
			})
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("expected exit code 0, got %d (stderr %q)", code, stderr.String())
			}
			if runner.qaCalls != tt.wantQACalls {
				t.Fatalf("expected %d QA call(s), got %d", tt.wantQACalls, runner.qaCalls)
			}
			hasQAVerdict := strings.Contains(stdout.String(), "qa pass — "+implementQAReportRelPath())
			if wantVerdict := tt.wantQACalls > 0; hasQAVerdict != wantVerdict {
				t.Fatalf("expected QA verdict line presence %v, got stdout:\n%s", wantVerdict, stdout.String())
			}
		})
	}
}

func TestRunImplementInteractiveInputPersistsAgentButNotSpecOrQA(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
	})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		qaReport:     implementQAReport("pass"),
	}
	configureImplementClaudeReasoning(t, repoDir)
	withImplementCollaborators(t, runner)
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		values := req.Values
		values.Spec = implementTestSlug
		values.Agent = "claude"
		values.QA = true
		return values, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if code := RunContext(context.Background(), []string{"implement"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected first run exit 0, got %d (stderr %q)", code, stderr.String())
	}
	if defaults := readInteractiveDefaults(t, homeDir); defaults.Agent != "claude" {
		t.Fatalf("expected the Agent remembered after the first Run, got %#v", defaults)
	}

	// A second explicit Interactive Input reopens the flow. The configured
	// Agent remains the suggestion, while the Spec and QA choice are not
	// remembered because each Run's target is an explicit choice.
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "keep first run")
	var secondReq roundtui.InputRequest
	var secondCollected strings.Builder
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		secondReq = req
		return roundtui.CollectInput(context.Background(), req, strings.NewReader("1\n\n\n\n\n"), &secondCollected)
	})
	stdout.Reset()
	stderr.Reset()

	if code := RunContext(context.Background(), []string{"implement", "--interactive"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected second run exit 0, got %d (stderr %q)", code, stderr.String())
	}
	if secondReq.AgentSuggestion.Value != "codex" || secondReq.AgentSuggestion.Source != "config" {
		t.Fatalf("expected the configured Agent suggestion, got %#v", secondReq.AgentSuggestion)
	}
	if secondReq.Values.Spec != "" {
		t.Fatalf("expected the spec slug not remembered across invocations, got %q", secondReq.Values.Spec)
	}
	if secondReq.Values.QA {
		t.Fatal("expected the QA choice not remembered across invocations")
	}
	if !strings.Contains(secondCollected.String(), "QA gate [y/N]:") {
		t.Fatalf("expected the second invocation to show the default QA prompt, got:\n%s", secondCollected.String())
	}
}

func TestRunImplementInteractiveForcedWithFlagsProvidedStillOpensFlow(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
	})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withImplementCollaborators(t, runner)
	var inputReq roundtui.InputRequest
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		inputReq = req
		return req.Values, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--interactive"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr %q)", code, stderr.String())
	}
	if inputReq.Command != "implement" {
		t.Fatalf("expected --interactive to open the flow, got %#v", inputReq)
	}
	if inputReq.Values.Spec != implementTestSlug {
		t.Fatalf("expected the provided Spec pre-filled, got %q", inputReq.Values.Spec)
	}
	run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if run.SpecSlug != implementTestSlug {
		t.Fatalf("expected the Run to target the provided Spec, got %q", run.SpecSlug)
	}
}

func TestRunImplementExecutesSpecEndToEnd(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Write the widget guide", taskType: "docs", verification: []string{"echo docs-check"}},
		{id: "task_02", title: "Build the widget backend", needs: []string{"task_01"}, verification: []string{"echo backend-check"}},
	})
	runner := &implementFakeRunner{
		gitRoot: repoDir,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_02": spec.StatusCompleted,
		},
	}
	committer, verifier, pusher, sourceResolver := withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr %q)", code, stderr.String())
	}
	expected := "task_01 completed — Write the widget guide\n" +
		"task_02 completed — Build the widget backend\n" +
		"Clean: all 2 Task(s) completed.\n"
	if stdout.String() != expected {
		t.Fatalf("expected deterministic stdout report:\n%q\ngot:\n%q", expected, stdout.String())
	}
	if runner.calls != 2 || runner.taskIDs[0] != "task_01" || runner.taskIDs[1] != "task_02" {
		t.Fatalf("expected the Agent to run task_01 then task_02, got %v", runner.taskIDs)
	}
	if committer.calls != 2 {
		t.Fatalf("expected one commit per Task, got %d", committer.calls)
	}
	if !strings.HasPrefix(committer.messages[0], "docs: write the widget guide") || !strings.Contains(committer.messages[0], "Roundfix-Task: task_01") {
		t.Fatalf("expected frontmatter-derived task_01 commit message, got %q", committer.messages[0])
	}
	if !strings.HasPrefix(committer.messages[1], "feat: build the widget backend") || !strings.Contains(committer.messages[1], "Roundfix-Spec: "+implementTestSlug) {
		t.Fatalf("expected frontmatter-derived task_02 commit message, got %q", committer.messages[1])
	}
	if verifier.calls != 2 {
		t.Fatalf("expected the Daemon to run each Task's Verification, got %d call(s)", verifier.calls)
	}
	if pusher.calls != 0 || sourceResolver.calls != 0 {
		t.Fatalf("spec Runs without implement.auto_push must not push or resolve Review Source threads, got push=%d source=%d", pusher.calls, sourceResolver.calls)
	}
	if !strings.Contains(stderr.String(), "reached Clean") {
		t.Fatalf("expected Clean outcome diagnostics on stderr, got %q", stderr.String())
	}
	if strings.Contains(stdout.String(), "Implement Run:") {
		t.Fatalf("Run id belongs on stderr, got stdout %q", stdout.String())
	}
	runID := implementRunIDFromStderr(t, stderr.String())
	run := implementRunFromStore(t, homeDir, runID)
	if run.Kind != store.KindImplement || run.State != store.StateClean {
		t.Fatalf("expected implement Run to reach Clean, got kind=%q state=%q", run.Kind, run.State)
	}
	if run.SpecSlug != implementTestSlug || run.GitRoot != repoDir || run.LocalBranch != "ma/widget-flow" {
		t.Fatalf("unexpected Run row: %#v", run)
	}
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}

func TestRunImplementPassesVerificationCapacityIntoTaskCycle(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", title: "Build the widget backend"}})
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), "worktree:\n  concurrency: 1\nverification:\n  concurrency: 3\n")
	gitImplement(t, repoDir, "add", ".roundfixrc.yml")
	gitImplement(t, repoDir, "commit", "-m", "configure verification capacity")
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected exit code 0, got %d stderr=%q", code, stderr.String())
	}
	runID := implementRunIDFromStderr(t, stderr.String())
	var startPayload map[string]any
	for _, entry := range runEventsForRun(t, homeDir, runID) {
		if entry.Event.Kind != runevent.KindDaemonStatus {
			continue
		}
		var payload map[string]any
		if err := json.Unmarshal(entry.Event.Payload, &payload); err != nil {
			t.Fatalf("decode daemon status payload: %v", err)
		}
		if payload["spec"] == implementTestSlug {
			startPayload = payload
			break
		}
	}
	if startPayload == nil {
		t.Fatal("expected journaled Task-cycle-start event")
	}
	if startPayload["concurrency"] != float64(1) || startPayload["task_capacity"] != float64(1) || startPayload["verification_capacity"] != float64(3) {
		t.Fatalf("unexpected effective capacities: %+v", startPayload)
	}
}

func TestRunImplementCleanCleanupFailureWarnsAndJournalsWithoutChangingReportOrExit(t *testing.T) {
	wantStdout, _, wantCode, _, _ := runCleanImplementForCleanup(t, nil)
	gotStdout, gotStderr, gotCode, homeDir, keptPath := runCleanImplementForCleanup(t, errors.New("forced cleanup failure"))

	if gotCode != wantCode {
		t.Fatalf("expected exit code to stay %d, got %d stderr=%q", wantCode, gotCode, gotStderr)
	}
	if gotStdout != wantStdout {
		t.Fatalf("stdout changed after cleanup failure\nwant:\n%q\ngot:\n%q", wantStdout, gotStdout)
	}
	warning := fmt.Sprintf("roundfix: Run Worktree cleanup failed; kept %s: forced cleanup failure\n", keptPath)
	if strings.Count(gotStderr, warning) != 1 {
		t.Fatalf("expected one cleanup warning %q, got stderr=%q", warning, gotStderr)
	}
	assertCleanCleanupWarningEvent(t, homeDir, gotStderr, keptPath, "forced cleanup failure")
	assertRunWorktreeExists(t, keptPath)
}

func TestRunImplementBootstrapFailureEndsFailedBeforeAgentWork(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:    "task_01",
		title: "Build the widget core",
	}})
	command := "printf bootstrap-output; exit 7"
	configureWorktreeBootstrap(t, repoDir, command, "1s")
	runner := &implementFakeRunner{gitRoot: repoDir}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected bootstrap failure exit %d, got %d (stderr %q)", exitRunFailed, code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout before any Task settles, got %q", stdout.String())
	}
	if runner.calls != 0 {
		t.Fatalf("expected bootstrap failure before Agent work, got %d Agent call(s)", runner.calls)
	}
	if !strings.Contains(stderr.String(), "bootstrap-output") {
		t.Fatalf("expected bootstrap output on stderr, got %q", stderr.String())
	}
	wantFailure := "worktree bootstrap failed: " + command + ": exit status 7"
	if !strings.Contains(stderr.String(), wantFailure) {
		t.Fatalf("expected stderr to contain %q, got %q", wantFailure, stderr.String())
	}

	run := onlyRunFromStore(t, homeDir)
	if run.State != store.StateFailed {
		t.Fatalf("expected Run Failed, got %q", run.State)
	}
	events := runEventsForRun(t, homeDir, run.ID)
	sawBootstrapOutput := false
	sawFailedOutcome := false
	for _, entry := range events {
		if entry.Event.Kind == runevent.KindDaemonBatch {
			t.Fatalf("expected no Batch event before bootstrap failure, got %#v", entry.Event)
		}
		if strings.Contains(entry.Event.Summary, "bootstrap-output") || strings.Contains(string(entry.Event.Payload), "bootstrap-output") {
			sawBootstrapOutput = true
		}
		if entry.Event.Kind == runevent.KindDaemonOutcome && strings.Contains(string(entry.Event.Payload), store.StateFailed) {
			sawFailedOutcome = true
		}
	}
	if !sawBootstrapOutput {
		t.Fatalf("expected bootstrap output in Run Event Journal, got %#v", events)
	}
	if !sawFailedOutcome {
		t.Fatalf("expected Failed outcome in Run Event Journal, got %#v", events)
	}
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}

func TestRunImplementBootstrapRunsBeforeAgentWorkAndVerification(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:           "task_01",
		title:        "Build the widget core",
		verification: []string{"test -f bootstrap.ready"},
	}})
	command := "pwd > bootstrap.cwd && printf ready > bootstrap.ready"
	mustWrite(t, filepath.Join(repoDir, ".gitignore"), "bootstrap.cwd\nbootstrap.ready\n")
	gitImplement(t, repoDir, "add", ".gitignore")
	gitImplement(t, repoDir, "commit", "-m", "ignore bootstrap markers")
	configureWorktreeBootstrap(t, repoDir, command, "1s")
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		onTask: func(req agent.ExecuteRequest, _ string) error {
			if got := mustRead(t, filepath.Join(req.GitRoot, "bootstrap.ready")); got != "ready" {
				return fmt.Errorf("bootstrap marker content = %q", got)
			}
			if got := strings.TrimSpace(mustRead(t, filepath.Join(req.GitRoot, "bootstrap.cwd"))); got != req.GitRoot {
				return fmt.Errorf("bootstrap cwd = %q, want %q", got, req.GitRoot)
			}
			return nil
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected Clean exit, got %d (stderr %q)", code, stderr.String())
	}
	if runner.calls != 1 || runner.taskIDs[0] != "task_01" {
		t.Fatalf("expected one Agent call after bootstrap, got calls=%d tasks=%v", runner.calls, runner.taskIDs)
	}
	if !strings.Contains(stdout.String(), "Clean: all 1 Task(s) completed.") {
		t.Fatalf("expected Clean stdout, got %q", stdout.String())
	}
	run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if run.State != store.StateClean {
		t.Fatalf("expected Run Clean, got %q", run.State)
	}
}

func TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the first slice", verification: []string{"test -f bootstrap.ready"}},
		{id: "task_02", title: "Build the second slice", verification: []string{"test -f bootstrap.ready"}},
	})
	command := "pwd > bootstrap.cwd && printf ready > bootstrap.ready"
	mustWrite(t, filepath.Join(repoDir, ".gitignore"), "bootstrap.cwd\nbootstrap.ready\n")
	gitImplement(t, repoDir, "add", ".gitignore")
	gitImplement(t, repoDir, "commit", "-m", "ignore bootstrap markers")
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), fmt.Sprintf("worktree:\n  concurrency: 2\n  bootstrap: %q\n  bootstrap_timeout: 1s\n", command))
	gitImplement(t, repoDir, "add", ".roundfixrc.yml")
	gitImplement(t, repoDir, "commit", "-m", "configure concurrent worktree bootstrap")
	runner := &implementFakeRunner{
		gitRoot: repoDir,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_02": spec.StatusCompleted,
		},
		onTask: func(req agent.ExecuteRequest, taskID string) error {
			if req.GitRoot == repoDir {
				return fmt.Errorf("%s ran in the user checkout instead of a Task Worktree", taskID)
			}
			ready, err := os.ReadFile(filepath.Join(req.GitRoot, "bootstrap.ready"))
			if err != nil {
				return fmt.Errorf("read %s bootstrap marker: %w", taskID, err)
			}
			if got := string(ready); got != "ready" {
				return fmt.Errorf("%s bootstrap marker content = %q", taskID, got)
			}
			cwd, err := os.ReadFile(filepath.Join(req.GitRoot, "bootstrap.cwd"))
			if err != nil {
				return fmt.Errorf("read %s bootstrap cwd: %w", taskID, err)
			}
			if got := strings.TrimSpace(string(cwd)); got != req.GitRoot {
				return fmt.Errorf("%s bootstrap cwd = %q, want %q", taskID, got, req.GitRoot)
			}
			return nil
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected Clean exit, got %d (stderr %q)", code, stderr.String())
	}
	if runner.calls != 2 {
		t.Fatalf("expected two Agent calls after Task Worktree bootstrap, got %d", runner.calls)
	}
	gotTasks := append([]string(nil), runner.taskIDs...)
	sort.Strings(gotTasks)
	if got := strings.Join(gotTasks, "|"); got != "task_01|task_02" {
		t.Fatalf("expected both Tasks to run, got %q", got)
	}
	if !strings.Contains(stdout.String(), "Clean: all 2 Task(s) completed.") {
		t.Fatalf("expected Clean stdout, got %q", stdout.String())
	}
	run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if run.State != store.StateClean {
		t.Fatalf("expected Run Clean, got %q", run.State)
	}
}

func TestRenderImplementTaskLinesKeepsGraphOrderWhenCompletionReversed(t *testing.T) {
	_, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build scheduler"},
		{id: "task_02", title: "Wire queue"},
		{id: "task_03", title: "Write docs"},
	})
	specsRoot := filepath.Join(repoDir, "docs", "specs")
	graph, err := spec.Load(specsRoot, implementTestSlug)
	if err != nil {
		t.Fatalf("load spec: %v", err)
	}
	for _, taskID := range []string{"task_03", "task_02", "task_01"} {
		if err := spec.SetStatus(implementTaskPath(repoDir, taskID), spec.StatusCompleted); err != nil {
			t.Fatalf("settle %s: %v", taskID, err)
		}
	}

	report, counts := renderImplementTaskLines(specsRoot, graph, true)

	expected := "task_01 completed — Build scheduler\n" +
		"task_02 completed — Wire queue\n" +
		"task_03 completed — Write docs\n"
	if report != expected {
		t.Fatalf("expected graph-order report after reversed completion:\n%q\ngot:\n%q", expected, report)
	}
	if counts.completed != 3 || counts.failed != 0 || counts.skipped != 0 || counts.pending != 0 {
		t.Fatalf("expected completed counts after reversed completion, got %+v", counts)
	}
}

func TestRenderImplementTaskLinesAddsReasonsForFailedAndSkippedTasks(t *testing.T) {
	_, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build scheduler"},
		{id: "task_02", title: "Wire queue", needs: []string{"task_01"}},
		{id: "task_03", title: "Write docs"},
	})
	specsRoot := filepath.Join(repoDir, "docs", "specs")
	graph, err := spec.Load(specsRoot, implementTestSlug)
	if err != nil {
		t.Fatalf("load spec: %v", err)
	}
	if err := spec.SetStatus(implementTaskPath(repoDir, "task_01"), spec.StatusFailed); err != nil {
		t.Fatalf("settle task_01 failed: %v", err)
	}
	if err := spec.SetStatus(implementTaskPath(repoDir, "task_03"), spec.StatusCompleted); err != nil {
		t.Fatalf("settle task_03 completed: %v", err)
	}
	outcomes := []daemon.TaskOutcome{
		{
			Task:   "task_01",
			Status: string(spec.StatusFailed),
			Reason: `Verification failed: command "make verify" exited with exit status 7; diagnostics: .roundfix/verification.log`,
		},
		{
			Task:   "task_02",
			Status: "skipped",
			Reason: "needs not completed: task_01",
		},
		{
			Task:   "task_03",
			Status: string(spec.StatusCompleted),
		},
	}

	report, counts := renderImplementTaskLinesWithOutcomes(specsRoot, graph, true, outcomes)

	expected := "task_01 failed — Build scheduler\n" +
		"  reason: Verification failed: command \"make verify\" exited with exit status 7; diagnostics: .roundfix/verification.log\n" +
		"task_02 skipped — Wire queue\n" +
		"  reason: needs not completed: task_01\n" +
		"task_03 completed — Write docs\n"
	if report != expected {
		t.Fatalf("expected report with failed/skipped reasons:\n%q\ngot:\n%q", expected, report)
	}
	if counts.completed != 1 || counts.failed != 1 || counts.skipped != 1 || counts.pending != 0 {
		t.Fatalf("expected mixed counts, got %+v", counts)
	}
}

func TestRenderImplementTaskLinesNormalizesMultilineReasons(t *testing.T) {
	_, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build scheduler"},
	})
	specsRoot := filepath.Join(repoDir, "docs", "specs")
	graph, err := spec.Load(specsRoot, implementTestSlug)
	if err != nil {
		t.Fatalf("load spec: %v", err)
	}
	if err := spec.SetStatus(implementTaskPath(repoDir, "task_01"), spec.StatusFailed); err != nil {
		t.Fatalf("settle task_01 failed: %v", err)
	}
	outcomes := []daemon.TaskOutcome{{
		Task:   "task_01",
		Status: string(spec.StatusFailed),
		Reason: "Verification failed:\ncommand \"make verify\" exited\nwith exit status 7",
	}}

	report, counts := renderImplementTaskLinesWithOutcomes(specsRoot, graph, true, outcomes)

	expected := "task_01 failed — Build scheduler\n" +
		"  reason: Verification failed: command \"make verify\" exited with exit status 7\n"
	if report != expected {
		t.Fatalf("expected single-line reason:\n%q\ngot:\n%q", expected, report)
	}
	if counts.failed != 1 || counts.completed != 0 || counts.skipped != 0 || counts.pending != 0 {
		t.Fatalf("expected one failed count, got %+v", counts)
	}
}

func TestRunImplementAutoPushOutcomeMatrix(t *testing.T) {
	tests := []struct {
		name           string
		enableAutoPush bool
		configUpstream bool
		args           []string
		runner         func(string) *implementFakeRunner
		wantCode       int
		wantState      string
		wantPushes     int
		wantStdout     []string
	}{
		{
			name:           "clean qa pass with key pushes",
			enableAutoPush: true,
			configUpstream: true,
			args:           []string{"implement", "--spec", implementTestSlug, "--qa", "--no-input"},
			runner: func(repoDir string) *implementFakeRunner {
				return &implementFakeRunner{
					gitRoot:      repoDir,
					statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
					qaReport:     implementQAReport("pass"),
				}
			},
			wantCode:   0,
			wantState:  store.StateClean,
			wantPushes: 1,
			wantStdout: []string{
				"qa pass — " + implementQAReportRelPath() + "\n",
				"Clean: all 1 Task(s) completed.\n",
				"pushed origin/ma/widget-flow\n",
			},
		},
		{
			name:           "clean without key does not push",
			configUpstream: true,
			args:           []string{"implement", "--spec", implementTestSlug, "--no-input"},
			runner: func(repoDir string) *implementFakeRunner {
				return &implementFakeRunner{
					gitRoot:      repoDir,
					statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
				}
			},
			wantCode:   0,
			wantState:  store.StateClean,
			wantPushes: 0,
			wantStdout: []string{"Clean: all 1 Task(s) completed.\n"},
		},
		{
			name:           "unresolved failed task does not push",
			enableAutoPush: true,
			configUpstream: true,
			args:           []string{"implement", "--spec", implementTestSlug, "--no-input"},
			runner: func(repoDir string) *implementFakeRunner {
				return &implementFakeRunner{
					gitRoot:      repoDir,
					statusByTask: map[string]spec.Status{"task_01": spec.StatusFailed},
				}
			},
			wantCode:   1,
			wantState:  store.StateUnresolved,
			wantPushes: 0,
			wantStdout: []string{"Unresolved: 0 completed, 1 failed, 0 skipped, 0 pending.\n"},
		},
		{
			name:           "stopped does not push",
			enableAutoPush: true,
			configUpstream: true,
			args:           []string{"implement", "--spec", implementTestSlug, "--no-input"},
			runner: func(repoDir string) *implementFakeRunner {
				return &implementFakeRunner{
					gitRoot:   repoDir,
					errByTask: map[string]error{"task_01": agent.StopError{Err: context.Canceled}},
				}
			},
			wantCode:   0,
			wantState:  store.StateStopped,
			wantPushes: 0,
			wantStdout: []string{"Stopped: 0 completed, 0 failed, 0 skipped, 1 pending.\n"},
		},
		{
			name:           "qa fail does not push",
			enableAutoPush: true,
			configUpstream: true,
			args:           []string{"implement", "--spec", implementTestSlug, "--qa", "--no-input"},
			runner: func(repoDir string) *implementFakeRunner {
				return &implementFakeRunner{
					gitRoot:      repoDir,
					statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
					qaReport:     implementQAReport("fail"),
				}
			},
			wantCode:   1,
			wantState:  store.StateUnresolved,
			wantPushes: 0,
			wantStdout: []string{"qa fail — " + implementQAReportRelPath() + "\n", "Unresolved: 1 completed, 0 failed, 0 skipped, 0 pending.\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
				{id: "task_01", title: "Build the widget core"},
			})
			if tt.enableAutoPush {
				configureImplementAutoPush(t, repoDir, true)
			}
			if tt.configUpstream {
				configureImplementUpstream(t, repoDir, "origin", "ma/widget-flow")
			}
			_, _, pusher, _ := withImplementCollaborators(t, tt.runner(repoDir))
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != tt.wantCode {
				t.Fatalf("expected exit code %d, got %d (stderr %q)", tt.wantCode, code, stderr.String())
			}
			if pusher.calls != tt.wantPushes {
				t.Fatalf("expected %d push call(s), got %d args=%v", tt.wantPushes, pusher.calls, pusher.args)
			}
			if tt.wantPushes == 1 {
				if got := strings.Join(pusher.args, " "); got != "push origin HEAD:ma/widget-flow" {
					t.Fatalf("expected push invocation %q, got %q", "push origin HEAD:ma/widget-flow", got)
				}
				if len(pusher.workDirs) != 1 || pusher.workDirs[0] != repoDir {
					t.Fatalf("expected push from git root %q, got %v", repoDir, pusher.workDirs)
				}
			} else if strings.Contains(stdout.String(), "pushed ") {
				t.Fatalf("stdout must not include pushed line without a push, got %q", stdout.String())
			}
			for _, want := range tt.wantStdout {
				if !strings.Contains(stdout.String(), want) {
					t.Fatalf("expected stdout to contain %q, got %q", want, stdout.String())
				}
			}
			runID := implementRunIDFromStderr(t, stderr.String())
			run := implementRunFromStore(t, homeDir, runID)
			if run.State != tt.wantState {
				t.Fatalf("expected Run state %q, got %q", tt.wantState, run.State)
			}
			if tt.wantPushes == 1 {
				_, events := journaledRunEvents(t, homeDir, stderr.String())
				assertImplementPushEvent(t, events, "pushed")
			}
			assertNoActiveRunInGitRoot(t, homeDir, repoDir)
		})
	}
}

func TestRunImplementReportPrintsVerificationFailureReason(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core", verification: []string{"make verify"}},
	})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	_, verifier, _, _ := withImplementCollaborators(t, runner)
	verifier.err = errors.New("exit status 7")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected unresolved implement exit %d, got %d stderr=%q stdout=%q", exitRunFailed, code, stderr.String(), stdout.String())
	}
	wantPrefix := "task_01 failed — Build the widget core\n" +
		"  reason: Verification failed: command \"make verify\" exited with exit status 7; diagnostics: "
	if !strings.Contains(stdout.String(), wantPrefix) {
		t.Fatalf("expected verification reason prefix %q, got stdout=%q", wantPrefix, stdout.String())
	}
	if !strings.Contains(stdout.String(), "Unresolved: 0 completed, 1 failed, 0 skipped, 0 pending.\n") {
		t.Fatalf("expected unresolved outcome line, got stdout=%q", stdout.String())
	}
	runID := implementRunIDFromStderr(t, stderr.String())
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateUnresolved {
		t.Fatalf("expected Run state %q, got %q", store.StateUnresolved, run.State)
	}
}

func TestRunImplementReportPrintsModelNotAdvertisedReason(t *testing.T) {
	const reason = `Agent Model "gpt-5.6-sol" not advertised by runtime "codex"; advertised: gpt-5.5, gpt-5.1`
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core", verification: []string{"make verify"}},
	})
	runner := &implementFakeRunner{
		gitRoot: repoDir,
		errByTask: map[string]error{
			"task_01": &agent.BatchFailureError{
				ExitCode: 1,
				Reason:   "agent/protocol error",
				Err: &agent.ModelNotAdvertisedError{
					Runtime:    "codex",
					Model:      "gpt-5.6-sol",
					Advertised: []string{"gpt-5.5", "gpt-5.1"},
				},
			},
		},
	}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected unresolved implement exit %d, got %d stderr=%q stdout=%q", exitRunFailed, code, stderr.String(), stdout.String())
	}
	want := "task_01 failed — Build the widget core\n" +
		"  reason: " + reason + "\n"
	if !strings.Contains(stdout.String(), want) {
		t.Fatalf("expected model rejection reason line %q, got stdout=%q", want, stdout.String())
	}
	runID := implementRunIDFromStderr(t, stderr.String())
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateUnresolved {
		t.Fatalf("expected Run state %q, got %q", store.StateUnresolved, run.State)
	}
}

func TestRunImplementAutoPushMissingUpstreamWarnsAndStaysClean(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
	})
	configureImplementAutoPush(t, repoDir, true)
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	_, _, pusher, _ := withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected Clean path exit 0, got %d stderr=%q", code, stderr.String())
	}
	if pusher.calls != 0 {
		t.Fatalf("expected no push without upstream, got %d call(s)", pusher.calls)
	}
	note := "Spec Run push skipped: branch ma/widget-flow has no upstream; set an upstream or push manually."
	if strings.Count(stderr.String(), note) != 1 {
		t.Fatalf("expected one missing-upstream note %q, got stderr %q", note, stderr.String())
	}
	if strings.Contains(stdout.String(), "pushed ") {
		t.Fatalf("stdout must not include pushed line without a push, got %q", stdout.String())
	}
	runID, events := journaledRunEvents(t, homeDir, stderr.String())
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateClean {
		t.Fatalf("expected Run to remain Clean, got %q", run.State)
	}
	assertImplementPushEvent(t, events, "skipped")
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}

func TestRunImplementAutoPushFailureEndsFailedAndJournalsPush(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
	})
	configureImplementAutoPush(t, repoDir, true)
	configureImplementUpstream(t, repoDir, "origin", "ma/widget-flow")
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	_, _, pusher, _ := withImplementCollaborators(t, runner)
	pusher.err = errors.New("push failed")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected push failure exit 1, got %d stderr=%q", code, stderr.String())
	}
	if pusher.calls != 1 {
		t.Fatalf("expected one push attempt, got %d", pusher.calls)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout on push failure, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "push Clean spec Run: push failed") {
		t.Fatalf("expected push failure diagnostic, got %q", stderr.String())
	}
	runID, events := journaledRunEvents(t, homeDir, stderr.String())
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateFailed {
		t.Fatalf("expected Run state Failed, got %q", run.State)
	}
	assertImplementPushEvent(t, events, "failed")
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}

func assertImplementPushEvent(t *testing.T, events []store.JournalEvent, decision string) {
	t.Helper()
	needle := fmt.Sprintf("%q:%q", "decision", decision)
	for _, entry := range events {
		if entry.Event.Kind == runevent.KindDaemonPush && strings.Contains(string(entry.Event.Payload), needle) {
			return
		}
	}
	t.Fatalf("expected daemon.push event with decision %q, got %+v", decision, events)
}

func TestRunImplementNoAgentConsoleSuppressesAgentDisplayOnly(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
	})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		agentOutput:  "fake implement agent output\n",
	}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-agent-console", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean implement exit 0, got %d stderr=%q", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "fake implement agent output") {
		t.Fatalf("expected Agent console hidden from stderr, got %q", stderr.String())
	}
	for _, want := range []string{
		"Implement Run:",
		"implement selected Spec",
		"Verification passed (attempt 1).",
		"Task commit created",
		"reached Clean",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to keep daemon/progress line %q, got %q", want, stderr.String())
		}
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	assertJournalContainsAgentAndDaemonEvents(t, events, "fake implement agent output")
}

func TestRunImplementUsesConfiguredArtifactDirectoryForAgentLogs(t *testing.T) {
	tests := []struct {
		name         string
		config       string
		expectedBase func(homeDir string, repoDir string) string
	}{
		{
			name:   "repo relative",
			config: "var/roundfix-artifacts",
			expectedBase: func(_ string, repoDir string) string {
				return filepath.Join(repoDir, "var", "roundfix-artifacts")
			},
		},
		{
			name:   "home relative",
			config: "~/roundfix-artifacts",
			expectedBase: func(homeDir string, _ string) string {
				return filepath.Join(homeDir, "roundfix-artifacts")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
				{id: "task_01", title: "Build the widget core"},
			})
			mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), fmt.Sprintf("defaults:\n  artifact_dir: %q\nlogs:\n  agent: true\n", tt.config))
			gitImplement(t, repoDir, "add", ".roundfixrc.yml")
			gitImplement(t, repoDir, "commit", "-m", "configure artifact dir")
			runner := &implementFakeRunner{
				gitRoot:      repoDir,
				statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
				writeLogs:    true,
			}
			withImplementCollaborators(t, runner)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("expected exit code 0, got %d (stderr %q)", code, stderr.String())
			}
			runID := implementRunIDFromStderr(t, stderr.String())
			expectedLog := agent.LogPath(tt.expectedBase(homeDir, repoDir), runID, 1)
			if len(runner.logPaths) != 1 || runner.logPaths[0] != expectedLog {
				t.Fatalf("expected Agent log path %q, got %v", expectedLog, runner.logPaths)
			}
			if _, err := os.Stat(expectedLog); err != nil {
				t.Fatalf("expected fake Agent log under configured Artifact Directory: %v", err)
			}
			if !strings.Contains(stderr.String(), "Agent log: "+expectedLog) {
				t.Fatalf("expected stderr to name Agent log %q, got %q", expectedLog, stderr.String())
			}
			if _, err := os.Stat(filepath.Join(repoDir, ".roundfix")); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("expected no repo .roundfix directory, got err=%v", err)
			}
		})
	}
}

func TestRunImplementUsesOneAgentSessionPerRunAndCloses(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
		{id: "task_02", title: "Wire the widget API"},
	})
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), "worktree:\n  concurrency: 1\n")
	gitImplement(t, repoDir, "add", ".roundfixrc.yml")
	gitImplement(t, repoDir, "commit", "-m", "configure sequential task cycle")
	inner := &implementFakeRunner{
		gitRoot: repoDir,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_02": spec.StatusCompleted,
		},
	}
	runner := &sessionRecordingRunner{inner: inner}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean implement exit 0, got %d stderr=%q", code, stderr.String())
	}
	if inner.calls != 2 {
		t.Fatalf("expected two Task Agent calls, got %d", inner.calls)
	}
	runID := implementRunIDFromStderr(t, stderr.String())
	assertRecordedOneSessionForRun(t, runner, runID, inner.calls)
	_, events := journaledRunEvents(t, homeDir, "Implement Run: "+runID+"\n")
	assertSessionLifecycleEvents(t, events)
}

func TestRunImplementClosesAgentSessionForTerminalOutcomes(t *testing.T) {
	tests := []struct {
		name      string
		inner     *implementFakeRunner
		commitErr error
		wantCode  int
		wantState string
		closeErr  error
	}{
		{
			name: "clean",
			inner: &implementFakeRunner{
				statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
			},
			wantCode:  1,
			wantState: store.StateFailed,
			closeErr:  errors.New("close failed"),
		},
		{
			name: "unresolved",
			inner: &implementFakeRunner{
				statusByTask: map[string]spec.Status{"task_01": spec.StatusFailed},
			},
			wantCode:  1,
			wantState: store.StateUnresolved,
		},
		{
			name: "failed",
			inner: &implementFakeRunner{
				statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
			},
			commitErr: errors.New("commit failed"),
			wantCode:  1,
			wantState: store.StateFailed,
		},
		{
			name: "stopped",
			inner: &implementFakeRunner{
				errByTask: map[string]error{"task_01": agent.StopError{Err: context.Canceled}},
			},
			wantCode:  0,
			wantState: store.StateStopped,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
				{id: "task_01", title: "Build the widget core"},
			})
			tt.inner.gitRoot = repoDir
			runner := &sessionRecordingRunner{inner: tt.inner, closeErr: tt.closeErr}
			committer, _, _, _ := withImplementCollaborators(t, runner)
			committer.err = tt.commitErr
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

			if code != tt.wantCode {
				t.Fatalf("expected exit code %d, got %d stderr=%q", tt.wantCode, code, stderr.String())
			}
			runID := implementRunIDFromStderr(t, stderr.String())
			run := implementRunFromStore(t, homeDir, runID)
			if run.State != tt.wantState {
				t.Fatalf("expected Run state %q, got %q", tt.wantState, run.State)
			}
			assertRecordedOneSessionForRun(t, runner, runID, len(runner.runSessions))
		})
	}
}

func TestRunImplementPreflightFailures(t *testing.T) {
	tests := []struct {
		name     string
		seeds    []implementSeed
		args     []string
		mutate   func(t *testing.T, repoDir string)
		messages []string
	}{
		{
			name:     "spec not found",
			seeds:    []implementSeed{{id: "task_01"}},
			args:     []string{"implement", "--spec", "9999-missing"},
			messages: []string{`Spec "9999-missing" not found`, "check the slug"},
		},
		{
			name:  "spec not active",
			seeds: []implementSeed{{id: "task_01"}},
			mutate: func(t *testing.T, repoDir string) {
				mustWrite(t, filepath.Join(repoDir, "docs", "specs", implementTestSlug, "_prd.md"), "---\nstatus: draft\n---\n\n# PRD\n")
			},
			messages: []string{"is not active", `expected "active"`},
		},
		{
			name:  "task without verification",
			seeds: []implementSeed{{id: "task_01"}},
			mutate: func(t *testing.T, repoDir string) {
				mustWrite(t, implementTaskPath(repoDir, "task_01"), "---\ntask: task_01\nstatus: pending\ntype: backend\n---\n\n# Title\n")
			},
			messages: []string{`Task "task_01" has no Verification commands`, "## Verification"},
		},
		{
			name: "task graph cycle",
			seeds: []implementSeed{
				{id: "task_01", needs: []string{"task_02"}},
				{id: "task_02", needs: []string{"task_01"}},
			},
			messages: []string{"dependency cycle", "task_01, task_02"},
		},
		{
			name:  "default branch veto",
			seeds: []implementSeed{{id: "task_01"}},
			mutate: func(t *testing.T, repoDir string) {
				gitImplement(t, repoDir, "checkout", "main")
			},
			messages: []string{`current branch "main"`, `default branch "main"`, "name-match", "switch to a work branch"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := newImplementWorkspace(t, tt.seeds)
			withImplementCollaborators(t, &implementFakeRunner{gitRoot: repoDir})
			if tt.mutate != nil {
				tt.mutate(t, repoDir)
			}
			args := tt.args
			if args == nil {
				args = []string{"implement", "--spec", implementTestSlug, "--no-input"}
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), args, &stdout, &stderr)

			if code != 2 {
				t.Fatalf("expected exit code 2, got %d (stderr %q)", code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			for _, message := range tt.messages {
				if !strings.Contains(stderr.String(), message) {
					t.Fatalf("expected preflight message %q, got %q", message, stderr.String())
				}
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestImplementRejectsInvalidTaskTypeBeforeSideEffects(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core", taskType: "Backend"},
	})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d (stderr %q)", code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	for _, want := range []string{
		filepath.Join(repoDir, "docs", "specs", implementTestSlug, "task_01.md"),
		`"Backend"`,
		"backend, frontend, data, infra, docs, test, chore",
		"frontmatter",
		"type",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected invalid Task Type diagnostic to contain %q, got %q", want, stderr.String())
		}
	}
	if len(runner.probeRequests) != 0 {
		t.Fatalf("expected invalid Task Type to skip selection probes, got %#v", runner.probeRequests)
	}
	if len(runner.fallbackSets) != 0 {
		t.Fatalf("expected invalid Task Type to skip fallback probes, got %#v", runner.fallbackSets)
	}
	if runner.calls != 0 {
		t.Fatalf("expected invalid Task Type to skip Agent invocation, got %d call(s)", runner.calls)
	}
	assertNoRunDatabase(t, homeDir)
	worktreeRoot := filepath.Join(homeDir, ".roundfix", "worktrees")
	if _, err := os.Stat(worktreeRoot); err == nil {
		t.Fatalf("expected no Run Worktree root at %s", worktreeRoot)
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stat Run Worktree root: %v", err)
	}
}

func TestRunImplementDirtyWorkingTreePrintsNoteAndRuns(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	withImplementCollaborators(t, &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	})
	mustWrite(t, filepath.Join(repoDir, "leftover.txt"), "failed attempt leftovers\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected dirty working tree to continue, got exit %d (stderr %q)", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "has 1 uncommitted change(s); implement will run in a Run Worktree") {
		t.Fatalf("expected dirty-tree note, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "overlapping local changes end the Run Integration Pending") {
		t.Fatalf("expected Integration Pending note, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Clean: all 1 Task(s) completed.") {
		t.Fatalf("expected clean outcome, got %q", stdout.String())
	}
	runID := implementRunIDFromStderr(t, stderr.String())
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateClean {
		t.Fatalf("expected clean run, got %s", run.State)
	}
	if _, err := os.Stat(filepath.Join(repoDir, "leftover.txt")); err != nil {
		t.Fatalf("expected dirty user file to survive: %v", err)
	}
}

func TestRunImplementRealWorktreeFastForwardsAndCleansPreservingNonOverlappingUserDirt(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:           "task_01",
		title:        "Build isolated work",
		verification: []string{"test -f agent.txt"},
	}})
	mustWrite(t, filepath.Join(repoDir, "local.txt"), "user dirt\n")
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		onTask: func(req agent.ExecuteRequest, _ string) error {
			return os.WriteFile(filepath.Join(req.GitRoot, "agent.txt"), []byte("agent work\n"), 0o644)
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected Clean fast-forward exit, got %d (stderr %q)", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "has 1 uncommitted change(s); implement will run in a Run Worktree") {
		t.Fatalf("expected dirty-tree note, got %q", stderr.String())
	}
	runID := implementRunIDFromStderr(t, stderr.String())
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateClean {
		t.Fatalf("expected Clean Run, got %s", run.State)
	}
	if run.WorkDir == "" {
		t.Fatal("expected Run Worktree path to be recorded")
	}
	assertRunWorktreeRemoved(t, run.WorkDir)
	assertRunBranchRemoved(t, repoDir, runworktree.BranchName(runID))
	if got := mustRead(t, filepath.Join(repoDir, "agent.txt")); got != "agent work\n" {
		t.Fatalf("expected integrated agent file, got %q", got)
	}
	if got := mustRead(t, filepath.Join(repoDir, "local.txt")); got != "user dirt\n" {
		t.Fatalf("expected non-overlapping user dirt preserved, got %q", got)
	}
	if status := gitImplementOutput(t, repoDir, "status", "--porcelain=v1"); !strings.Contains(status, "?? local.txt") {
		t.Fatalf("expected local.txt to remain untracked, got status %q", status)
	}
	if !strings.Contains(stdout.String(), "Clean: all 1 Task(s) completed.") {
		t.Fatalf("expected Clean stdout, got %q", stdout.String())
	}
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}

func TestRunImplementWorktreeIsolationExcludesConcurrentUserCommit(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:           "task_01",
		title:        "Build isolated work",
		verification: []string{"test -f agent.txt"},
	}})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		onTask: func(req agent.ExecuteRequest, _ string) error {
			if err := os.WriteFile(filepath.Join(req.GitRoot, "agent.txt"), []byte("agent work\n"), 0o644); err != nil {
				return err
			}
			mustWrite(t, filepath.Join(repoDir, "user.txt"), "user work\n")
			gitImplement(t, repoDir, "add", "user.txt")
			gitImplement(t, repoDir, "commit", "-m", "user work during run")
			return nil
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected IntegrationPending exit, got %d (stderr %q)", code, stderr.String())
	}
	runID := implementRunIDFromStderr(t, stderr.String())
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateIntegrationPending {
		t.Fatalf("expected IntegrationPending Run, got %s", run.State)
	}
	branch := runworktree.BranchName(runID)
	assertRunWorktreeExists(t, run.WorkDir)
	assertRunBranchExists(t, repoDir, branch)
	taskCommitFiles := gitImplementOutput(t, repoDir, "diff-tree", "--no-commit-id", "--name-only", "-r", branch)
	if !strings.Contains(taskCommitFiles, "agent.txt") {
		t.Fatalf("expected task commit to contain agent file, got %q", taskCommitFiles)
	}
	if strings.Contains(taskCommitFiles, "user.txt") {
		t.Fatalf("expected task commit to exclude user checkout commit, got %q", taskCommitFiles)
	}
	userFiles := gitImplementOutput(t, repoDir, "ls-tree", "-r", "--name-only", "HEAD")
	if !strings.Contains(userFiles, "user.txt") {
		t.Fatalf("expected user commit to survive on user branch, got %q", userFiles)
	}
	if strings.Contains(userFiles, "agent.txt") {
		t.Fatalf("expected unintegrated agent file to stay off user branch, got %q", userFiles)
	}
	if subject := strings.TrimSpace(gitImplementOutput(t, repoDir, "log", "-1", "--format=%s")); subject != "user work during run" {
		t.Fatalf("expected user's concurrent commit at HEAD, got %q", subject)
	}
	if !strings.Contains(stderr.String(), "Integration command: git merge --ff-only "+branch) {
		t.Fatalf("expected integration command for pending Run, got %q", stderr.String())
	}
}

func TestRunImplementOverlapEndsIntegrationPendingAndPrintedCommandWorks(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:           "task_01",
		title:        "Edit shared file",
		verification: []string{"grep -q run shared.txt"},
	}})
	mustWrite(t, filepath.Join(repoDir, "shared.txt"), "base\n")
	gitImplement(t, repoDir, "add", "shared.txt")
	gitImplement(t, repoDir, "commit", "-m", "seed shared file")
	baseHead := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
	mustWrite(t, filepath.Join(repoDir, "shared.txt"), "user\n")
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		onTask: func(req agent.ExecuteRequest, _ string) error {
			return os.WriteFile(filepath.Join(req.GitRoot, "shared.txt"), []byte("run\n"), 0o644)
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected IntegrationPending exit, got %d (stderr %q)", code, stderr.String())
	}
	runID := implementRunIDFromStderr(t, stderr.String())
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateIntegrationPending {
		t.Fatalf("expected IntegrationPending Run, got %s", run.State)
	}
	branch := runworktree.BranchName(runID)
	if head := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD")); head != baseHead {
		t.Fatalf("expected target branch unmoved at %s, got %s", baseHead, head)
	}
	if got := mustRead(t, filepath.Join(repoDir, "shared.txt")); got != "user\n" {
		t.Fatalf("expected overlapping user dirt intact, got %q", got)
	}
	status := gitImplementOutput(t, repoDir, "status", "--porcelain=v1")
	if !strings.Contains(status, " M shared.txt") || strings.Contains(status, "M  shared.txt") {
		t.Fatalf("expected unstaged overlap without phantom staged entries, got status %q", status)
	}
	command := integrationCommandFromStderr(t, stderr.String())
	wantCommand := "git merge --ff-only " + branch
	if command != wantCommand {
		t.Fatalf("expected integration command %q, got %q", wantCommand, command)
	}
	wantOutcome := "IntegrationPending: 1 completed, 0 failed, 0 skipped, 0 pending; integrate with " + wantCommand
	if !strings.Contains(stdout.String(), wantOutcome) {
		t.Fatalf("expected Integration Pending outcome line %q, got %q", wantOutcome, stdout.String())
	}
	gitImplement(t, repoDir, "stash", "push", "-m", "overlap before integration")
	runPrintedIntegrationCommand(t, repoDir, command)
	if got := mustRead(t, filepath.Join(repoDir, "shared.txt")); got != "run\n" {
		t.Fatalf("expected printed command to fast-forward shared file, got %q", got)
	}
}

func TestRunImplementUnresolvedKeepsRealRunWorktreeAndPrintsPath(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:    "task_01",
		title: "Leave unresolved work",
	}})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusFailed},
		onTask: func(req agent.ExecuteRequest, _ string) error {
			return os.WriteFile(filepath.Join(req.GitRoot, "attempt.txt"), []byte("needs user\n"), 0o644)
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected Unresolved exit, got %d (stderr %q)", code, stderr.String())
	}
	runID := implementRunIDFromStderr(t, stderr.String())
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateUnresolved {
		t.Fatalf("expected Unresolved Run, got %s", run.State)
	}
	assertRunWorktreeExists(t, run.WorkDir)
	assertRunBranchExists(t, repoDir, runworktree.BranchName(runID))
	if !strings.Contains(stderr.String(), "Run Worktree kept: "+run.WorkDir) {
		t.Fatalf("expected kept Run Worktree path on stderr, got %q", stderr.String())
	}
	if got := mustRead(t, filepath.Join(run.WorkDir, "attempt.txt")); got != "needs user\n" {
		t.Fatalf("expected unresolved work in kept Run Worktree, got %q", got)
	}
	if !strings.Contains(stdout.String(), "Unresolved: 0 completed, 1 failed") {
		t.Fatalf("expected Unresolved stdout, got %q", stdout.String())
	}
}

func TestRunImplementPreflightReapsEmptyTerminalRunAndTaskWorktrees(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:     "task_01",
		title:  "Build after cleanup",
		status: string(spec.StatusPending),
	}})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	staleRun, _, staleTask := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "task_01", store.StateStopped)
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected implement exit 0, got %d stderr=%q", code, stderr.String())
	}
	for _, expected := range []string{
		"roundfix: reaped terminal Worktree path=" + staleRun.WorkDir + " branch=" + runworktree.BranchName(staleRun.ID),
		"roundfix: reaped terminal Worktree path=" + staleTask.Path + " branch=" + staleTask.Branch,
	} {
		if !strings.Contains(stderr.String(), expected) {
			t.Fatalf("expected stderr to contain %q, got %q", expected, stderr.String())
		}
	}
	assertRunWorktreeRemoved(t, staleRun.WorkDir)
	assertRunBranchRemoved(t, repoDir, runworktree.BranchName(staleRun.ID))
	assertRunWorktreeRemoved(t, staleTask.Path)
	assertRunBranchRemoved(t, repoDir, staleTask.Branch)
	if !strings.Contains(stdout.String(), "Clean: all 1 Task(s) completed.") {
		t.Fatalf("expected Clean implement output, got %q", stdout.String())
	}
}

func TestRunImplementPreflightTerminalReachableChangedBranch(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:     "task_01",
		title:  "Build after reachable cleanup",
		status: string(spec.StatusPending),
	}})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	staleRun, staleRef, _ := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "", store.StateStopped)
	mustWrite(t, filepath.Join(staleRef.Path, "reachable.txt"), "reachable\n")
	gitImplement(t, staleRef.Path, "add", "reachable.txt")
	gitImplement(t, staleRef.Path, "commit", "-m", "reachable terminal work")
	gitImplement(t, repoDir, "merge", "--ff-only", staleRef.Branch)
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected implement exit 0, got %d stderr=%q", code, stderr.String())
	}
	assertRunWorktreeRemoved(t, staleRun.WorkDir)
	assertRunBranchRemoved(t, repoDir, staleRef.Branch)
	if !strings.Contains(stderr.String(), "roundfix: reaped terminal Worktree path="+staleRun.WorkDir+" branch="+staleRef.Branch) {
		t.Fatalf("expected reachable terminal Run cleanup notice, got %q", stderr.String())
	}
}

func TestRunImplementPreflightTerminalUniqueChangedBranch(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:     "task_01",
		title:  "Build while preserving unique work",
		status: string(spec.StatusPending),
	}})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	staleRun, staleRef, _ := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "", store.StateStopped)
	mustWrite(t, filepath.Join(staleRef.Path, "unique.txt"), "unique\n")
	gitImplement(t, staleRef.Path, "add", "unique.txt")
	gitImplement(t, staleRef.Path, "commit", "-m", "unique terminal work")
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected implement exit 0, got %d stderr=%q", code, stderr.String())
	}
	assertRunWorktreeExists(t, staleRun.WorkDir)
	assertRunBranchExists(t, repoDir, staleRef.Branch)
	if strings.Contains(stderr.String(), "reaped terminal Worktree path="+staleRun.WorkDir) {
		t.Fatalf("expected unique terminal Run to be preserved, got %q", stderr.String())
	}
}

func TestRunImplementPreflightClosesTerminalRunSessionsOnly(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:     "task_01",
		title:  "Build after session cleanup",
		status: string(spec.StatusPending),
	}})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	terminalRun, _, terminalTask := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "task_01", store.StateStopped)
	activeOther := createActiveImplementRunForStopInGitRoot(t, homeDir, filepath.Join(repoDir, "other"), "0002-other-spec")
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withAgentRunner(t, runner)
	withRoundfixSessionLister(t, func(_ context.Context, runtime agent.RuntimeSpec, workDir string) ([]agent.RoundfixSession, error) {
		if runtime.ID != "codex" || workDir != repoDir {
			t.Fatalf("unexpected session sweep list runtime=%#v workDir=%q", runtime, workDir)
		}
		return []agent.RoundfixSession{
			{Name: "roundfix-" + terminalRun.ID, RunID: terminalRun.ID},
			{Name: "roundfix-" + terminalRun.ID + "-task_01", RunID: terminalRun.ID, TaskID: "task_01"},
			{Name: "roundfix-" + activeOther.ID, RunID: activeOther.ID},
			{Name: "roundfix-run_20990101T000000Z_unknown", RunID: "run_20990101T000000Z_unknown"},
		}, nil
	})
	closeCalls := []agent.SessionRef{}
	withStopAgentSessionCloser(t, func(_ context.Context, runtime agent.RuntimeSpec, session agent.SessionRef) error {
		if runtime.ID != "codex" {
			t.Fatalf("unexpected close runtime: %#v", runtime)
		}
		closeCalls = append(closeCalls, session)
		return nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected implement exit 0, got %d stderr=%q", code, stderr.String())
	}
	wantCloseCalls := []agent.SessionRef{
		{Name: "roundfix-" + terminalRun.ID, WorkDir: terminalRun.WorkDir},
		{Name: "roundfix-" + terminalRun.ID + "-task_01", WorkDir: terminalTask.Path},
	}
	if len(closeCalls) != len(wantCloseCalls) {
		t.Fatalf("expected close calls %#v, got %#v", wantCloseCalls, closeCalls)
	}
	for i, want := range wantCloseCalls {
		if closeCalls[i] != want {
			t.Fatalf("close call %d\nwant: %#v\ngot:  %#v", i, want, closeCalls[i])
		}
	}
	for _, want := range []string{
		"roundfix: closed session roundfix-" + terminalRun.ID,
		"roundfix: closed session roundfix-" + terminalRun.ID + "-task_01",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
		}
	}
	if strings.Contains(stderr.String(), activeOther.ID) || strings.Contains(stderr.String(), "unknown") {
		t.Fatalf("active and unknown sessions must be untouched, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Clean: all 1 Task(s) completed.") {
		t.Fatalf("expected Clean implement output, got %q", stdout.String())
	}
}

func TestRunImplementPreflightPrunesRetainedRunStorage(t *testing.T) {
	ctx := context.Background()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:     "task_01",
		title:  "Build after retention cleanup",
		status: string(spec.StatusPending),
	}})
	artifactDir := filepath.Join(homeDir, "artifacts")
	writeUserConfig(t, homeDir, fmt.Sprintf("defaults:\n  artifact_dir: %q\nstore:\n  journal_retention: 336h\n", artifactDir))
	fixture := seedGCFixture(t, ctx, homeDir, artifactDir, time.Now().UTC())
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(ctx, []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected implement exit 0, got %d stderr=%q", code, stderr.String())
	}
	if want := "roundfix: pruned Run storage runs=1 journal_rows=2 artifact_bytes=12"; !strings.Contains(stderr.String(), want) {
		t.Fatalf("expected retention sweep summary %q, got %q", want, stderr.String())
	}
	assertRunEvents(t, homeDir, fixture.oldRun.ID, 0)
	assertRunEvents(t, homeDir, fixture.activeRun.ID, 1)
	assertRunEvents(t, homeDir, fixture.recentRun.ID, 1)
	assertPathMissing(t, fixture.oldArtifactDir)
	assertPathExists(t, fixture.activeArtifactDir)
	assertPathExists(t, fixture.recentArtifactDir)
	assertPathExists(t, fixture.orphanDir)
	if !strings.Contains(stdout.String(), "Clean: all 1 Task(s) completed.") {
		t.Fatalf("expected Clean implement output, got %q", stdout.String())
	}
}

func TestRunImplementPreflightRetentionPruneFailureIsNonFatal(t *testing.T) {
	ctx := context.Background()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:     "task_01",
		title:  "Build despite retention warning",
		status: string(spec.StatusPending),
	}})
	artifactDir := filepath.Join(homeDir, "artifacts")
	writeUserConfig(t, homeDir, fmt.Sprintf("defaults:\n  artifact_dir: %q\nstore:\n  journal_retention: 336h\n", artifactDir))
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run store: %v", err)
	}
	oldRun := createGCTestRun(t, ctx, runStore, artifactDir, "old-terminal-file-artifact", 1)
	if _, err := runStore.CompleteRun(ctx, oldRun.ID, store.StateClean); err != nil {
		t.Fatalf("complete old Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run store after seed: %v", err)
	}
	setRunTimestamps(t, homeDir, oldRun.ID, time.Now().Add(-400*time.Hour), time.Now().Add(-400*time.Hour))
	mustMkdir(t, filepath.Join(artifactDir, "runs"))
	mustWrite(t, filepath.Join(artifactDir, "runs", oldRun.ID), "not a directory")
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(ctx, []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected implement exit 0 despite retention warning, got %d stderr=%q", code, stderr.String())
	}
	if want := "roundfix: warning: Journal Retention prune failed:"; !strings.Contains(stderr.String(), want) {
		t.Fatalf("expected retention warning %q, got %q", want, stderr.String())
	}
	assertRunEvents(t, homeDir, oldRun.ID, 1)
	if !strings.Contains(stdout.String(), "Clean: all 1 Task(s) completed.") {
		t.Fatalf("expected Clean implement output, got %q", stdout.String())
	}
}

func TestRunImplementPreflightRetentionZeroSkipsPrune(t *testing.T) {
	ctx := context.Background()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:     "task_01",
		title:  "Build without retention cleanup",
		status: string(spec.StatusPending),
	}})
	artifactDir := filepath.Join(homeDir, "artifacts")
	writeUserConfig(t, homeDir, fmt.Sprintf("defaults:\n  artifact_dir: %q\nstore:\n  journal_retention: 0\n", artifactDir))
	fixture := seedGCFixture(t, ctx, homeDir, artifactDir, time.Now().UTC())
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(ctx, []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected implement exit 0, got %d stderr=%q", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "pruned Run storage") {
		t.Fatalf("expected no retention summary when Journal Retention is zero, got %q", stderr.String())
	}
	assertRunEvents(t, homeDir, fixture.oldRun.ID, 2)
	assertPathExists(t, fixture.oldArtifactDir)
	if !strings.Contains(stdout.String(), "Clean: all 1 Task(s) completed.") {
		t.Fatalf("expected Clean implement output, got %q", stdout.String())
	}
}

func TestRunImplementPreflightRejectsActiveRunInWorkingTree(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	withImplementCollaborators(t, &implementFakeRunner{gitRoot: repoDir})
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	blocking, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "ma/other-work",
		SpecSlug:    "0002-other-spec",
		OwnerPID:    os.Getpid(),
	})
	if err != nil {
		t.Fatalf("seed blocking run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(ctx, []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d (stderr %q)", code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	wantBlock := fmt.Sprintf(
		"an Active Run already holds working tree %s: run_id=%s kind=%s state=%s; stop it with: roundfix stop %s",
		repoDir,
		blocking.ID,
		blocking.Kind,
		blocking.State,
		blocking.ID,
	)
	if !strings.Contains(stderr.String(), wantBlock) {
		t.Fatalf("expected unchanged live-owner block %q, got %q", wantBlock, stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunImplementPreflightProbeFailureCreatesNoRun(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runner := &implementFakeRunner{gitRoot: repoDir, probeErr: errors.New("codex-acp is not on PATH")}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d (stderr %q)", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "codex-acp is not on PATH") {
		t.Fatalf("expected the probe failure reason, got %q", stderr.String())
	}
	if len(runner.probeRequests) != 1 || runner.probeRequests[0].WorkDir != repoDir {
		t.Fatalf("expected one selection preflight in git root %q, got %#v", repoDir, runner.probeRequests)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 0)
}

func TestImplementProfilePreflightFailureCreatesNoRunWorktreeOrAgentPrompt(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", taskType: "backend"}})
	runner := &implementFakeRunner{gitRoot: repoDir, probeErr: errors.New("adapter rejected configured tuple")}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
	}
	for _, want := range []string{"adapter rejected configured tuple", "backend preferred", "roundfix profiles configure"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("stderr missing %q in %q", want, stderr.String())
		}
	}
	if stdout.Len() != 0 || runner.calls != 0 {
		t.Fatalf("expected no stdout or Agent work, stdout=%q calls=%d", stdout.String(), runner.calls)
	}
	if len(runner.probeRequests) != 1 || runner.probeRequests[0].WorkDir != repoDir {
		t.Fatalf("expected one profile proof in git root %q, got %#v", repoDir, runner.probeRequests)
	}
	if len(runner.fallbackSets) != 0 {
		t.Fatalf("profile preflight must not discover fallback candidates, got %#v", runner.fallbackSets)
	}
	assertNoRunDatabase(t, homeDir)
	if _, err := os.Stat(filepath.Join(homeDir, ".roundfix", "worktrees")); err == nil {
		t.Fatalf("expected no Run Worktree root under %s", homeDir)
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stat worktree root: %v", err)
	}
}

func TestRunImplementSelectionFailureReportsProfileRemediationWithoutCreatingRun(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runner := &implementFakeRunner{
		gitRoot: repoDir,
		probeErr: &agent.SelectionPreflightError{
			Runtime:         "codex",
			Model:           "broken-model",
			ReasoningEffort: "unsupported",
			Err:             errors.New("selection rejected"),
		},
	}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"implement",
		"--spec", implementTestSlug,
		"--agent", "codex",
		"--model", "broken-model",
		"--reasoning-effort", "unsupported",
		"--no-input",
	}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected selection preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
	}
	for _, want := range []string{
		`profile proof failed for runtime "codex", model "broken-model", reasoning_effort "unsupported"`,
		"backend preferred",
		"adapter error: agent selection unavailable",
		"roundfix profiles configure --scope user|project",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
		}
	}
	if stdout.Len() != 0 || runner.calls != 0 {
		t.Fatalf("expected no stdout or Agent work, stdout=%q calls=%d", stdout.String(), runner.calls)
	}
	if len(runner.fallbackSets) != 0 {
		t.Fatalf("profile preflight must not probe dynamic fallback candidates, got %#v", runner.fallbackSets)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunImplementSelectionFailureDoesNotPromptForDynamicFallback(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", title: "Confirm fallback"}})
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	configContent := "runtimes:\n  codex:\n    model: broken-model\n    reasoning_effort: unsupported\n"
	mustWrite(t, configPath, configContent)
	gitImplement(t, repoDir, "add", ".roundfixrc.yml")
	gitImplement(t, repoDir, "commit", "-m", "configure broken selection")
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		probeErr: &agent.SelectionPreflightError{
			Runtime:         "codex",
			Model:           "broken-model",
			ReasoningEffort: "unsupported",
			Err:             errors.New("selection rejected"),
		},
	}
	withImplementCollaborators(t, runner)
	withFallbackConfirmation(t, "y\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"implement",
		"--spec", implementTestSlug,
	}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected selection preflight exit %d, got exit %d stderr=%q stdout=%q", exitPreflight, code, stderr.String(), stdout.String())
	}
	if strings.Contains(stderr.String(), "Fallback Selection:") || strings.Contains(stderr.String(), "Use this Fallback Selection for this Run?") {
		t.Fatalf("pre-Run profile proof must not prompt for dynamic fallback candidates, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "roundfix profiles configure") || !strings.Contains(stderr.String(), "roundfix profiles validate") {
		t.Fatalf("expected profile remediation guidance, got %q", stderr.String())
	}
	if got := mustRead(t, configPath); got != configContent {
		t.Fatalf("failed profile proof must not change Project Config\nwant: %q\n got: %q", configContent, got)
	}
	if stdout.Len() != 0 || runner.calls != 0 || len(runner.fallbackSets) != 0 {
		t.Fatalf("expected no output, Agent work, or dynamic fallback probes; stdout=%q calls=%d fallback=%#v", stdout.String(), runner.calls, runner.fallbackSets)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunImplementPassesOneRunSelectionOverridesToPreflight(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runner := &implementFakeRunner{gitRoot: repoDir, probeErr: errors.New("stop after selection preflight")}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"implement",
		"--spec", implementTestSlug,
		"--agent", "codex",
		"--model", "future-model",
		"--reasoning-effort", "experimental-reasoning",
		"--no-input",
	}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout diagnostics, got %q", stdout.String())
	}
	if len(runner.probeRequests) != 1 {
		t.Fatalf("expected one selection preflight, got %#v", runner.probeRequests)
	}
	got := runner.probeRequests[0].Runtime
	if got.Model != "future-model" || got.ReasoningEffort != "experimental-reasoning" {
		t.Fatalf("expected one-Run selection overrides at preflight, got %#v", got)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 0)
}

func TestRunImplementPersistsEffectiveSelection(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", title: "Store selection"}})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"implement",
		"--spec", implementTestSlug,
		"--agent", "codex",
		"--model", "stored-implement-model",
		"--reasoning-effort", "stored-implement-reasoning",
		"--no-input",
	}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean implement exit, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if run.Model != "stored-implement-model" || run.ReasoningEffort != "stored-implement-reasoning" {
		t.Fatalf("expected stored implement selection, got %#v", run)
	}
	if !strings.Contains(stderr.String(), "Agent Model: stored-implement-model") ||
		!strings.Contains(stderr.String(), "Default Reasoning Effort: stored-implement-reasoning") {
		t.Fatalf("expected implement progress to show concrete selection, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "auto") {
		t.Fatalf("expected no auto selection placeholder, got %q", stderr.String())
	}
}

func TestRunImplementAcceptsExplicitEmptyReasoningEffort(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", title: "Store model-managed selection"}})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"implement",
		"--spec", implementTestSlug,
		"--agent", "codex",
		"--model", "gpt-5.6-sol",
		"--reasoning-effort=",
		"--no-input",
	}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean implement exit, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if run.Model != "gpt-5.6-sol" || run.ReasoningEffort != "" {
		t.Fatalf("expected model-managed implement selection persisted, got %#v", run)
	}
	if !strings.Contains(stderr.String(), "Agent Model: gpt-5.6-sol") ||
		!strings.Contains(stderr.String(), "Default Reasoning Effort: model-managed") {
		t.Fatalf("expected implement progress to show model-managed selection, got %q", stderr.String())
	}
}

func TestRunImplementRejectsExplicitEmptySelectionOverrides(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "model",
			args: []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--model=", "--reasoning-effort", "high", "--no-input"},
			want: "--model cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, _ := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout diagnostics, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("expected stderr to contain %q, got %q", tt.want, stderr.String())
			}
			assertRunCount(t, store.DatabasePath(homeDir), 0)
		})
	}
}

func TestRunImplementSelectionOverrideRejectsPartialBeforeConfigLoad(t *testing.T) {
	homeDir, _ := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
	const invalidConfig = "defaults:\n  agent: [\n"
	writeUserConfig(t, homeDir, invalidConfig)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"implement", "--spec", implementTestSlug, "--reasoning-effort", "high", "--no-input",
	}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout diagnostics, got %q", stdout.String())
	}
	const wantGrammar = "--agent, --model, and --reasoning-effort must be provided together for a one-Run Agent Selection override; omit all three to use Agent Selection Profiles"
	if !strings.Contains(stderr.String(), wantGrammar) {
		t.Fatalf("expected selection grammar error before config load, got %q", stderr.String())
	}
	if got := mustRead(t, configPath); got != invalidConfig {
		t.Fatalf("partial override changed User Config\nwant: %q\n got: %q", invalidConfig, got)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 0)
	if _, err := os.Stat(filepath.Join(homeDir, ".roundfix", "worktrees")); err == nil {
		t.Fatalf("partial override created a Run Worktree root")
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stat Run Worktree root: %v", err)
	}
}

func TestRunImplementAllTasksCompletedReportsWithoutRun(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Write the widget guide", status: string(spec.StatusCompleted)},
		{id: "task_02", title: "Build the widget backend", status: string(spec.StatusCompleted), needs: []string{"task_01"}},
	})
	runner := &implementFakeRunner{gitRoot: repoDir}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr %q)", code, stderr.String())
	}
	expected := "task_01 completed — Write the widget guide\n" +
		"task_02 completed — Build the widget backend\n" +
		"All 2 Task(s) already completed; no Run was created.\n"
	if stdout.String() != expected {
		t.Fatalf("expected no-op report:\n%q\ngot:\n%q", expected, stdout.String())
	}
	if runner.calls != 0 {
		t.Fatalf("expected no Agent invocation, got %d", runner.calls)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunImplementFailedTaskEndsUnresolvedAndKeepsWorktree(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
		{id: "task_02", title: "Wire the widget API", needs: []string{"task_01"}},
		{id: "task_03", title: "Document the widget", taskType: "docs", needs: []string{"task_02"}},
	})
	firstRunner := &implementFakeRunner{
		gitRoot: repoDir,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_02": spec.StatusFailed,
		},
	}
	committer, _, _, _ := withImplementCollaborators(t, firstRunner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d (stderr %q)", code, stderr.String())
	}
	expected := "task_01 completed — Build the widget core\n" +
		"task_02 failed — Wire the widget API\n" +
		"  reason: Agent settled the Task failed\n" +
		"task_03 skipped — Document the widget\n" +
		"  reason: needs not completed: task_02\n" +
		"Unresolved: 1 completed, 1 failed, 1 skipped, 0 pending.\n"
	if stdout.String() != expected {
		t.Fatalf("expected Unresolved report:\n%q\ngot:\n%q", expected, stdout.String())
	}
	if committer.calls != 1 {
		t.Fatalf("expected one commit for the completed Task only, got %d", committer.calls)
	}
	firstRun := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if firstRun.State != store.StateUnresolved {
		t.Fatalf("expected first Run Unresolved, got %q", firstRun.State)
	}
	if firstRun.WorkDir == "" {
		t.Fatal("expected Unresolved Run to record a kept Run Worktree")
	}
	if _, err := os.Stat(firstRun.WorkDir); err != nil {
		t.Fatalf("expected kept Run Worktree to exist: %v", err)
	}
	if !strings.Contains(stderr.String(), "Run Worktree kept: "+firstRun.WorkDir) {
		t.Fatalf("expected kept Run Worktree path on stderr, got %q", stderr.String())
	}
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunImplementResumesStaleInProgressTask(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core", status: string(spec.StatusInProgress)},
	})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr %q)", code, stderr.String())
	}
	if runner.calls != 1 || runner.taskIDs[0] != "task_01" {
		t.Fatalf("expected the stale in_progress Task to execute, got %v", runner.taskIDs)
	}
	if !strings.Contains(stdout.String(), "task_01 completed — Build the widget core") {
		t.Fatalf("expected the stale Task settled completed, got %q", stdout.String())
	}
	run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if run.State != store.StateClean {
		t.Fatalf("expected Clean, got %q", run.State)
	}
}

func TestRunImplementStopRequestEndsStoppedWithInterruptMapping(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
		{id: "task_02", title: "Wire the widget API", needs: []string{"task_01"}},
	})
	runner := &implementFakeRunner{
		gitRoot:   repoDir,
		errByTask: map[string]error{"task_01": agent.StopError{Err: context.Canceled}},
	}
	committer, verifier, _, _ := withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stoppedCode := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)
	code := exitForInterrupt(stoppedCode, true)

	if stoppedCode != 0 {
		t.Fatalf("expected stopped command exit 0 before interrupt mapping, got %d (stderr %q)", stoppedCode, stderr.String())
	}
	if code != 130 {
		t.Fatalf("expected SIGINT exit 130, got %d", code)
	}
	if !strings.Contains(stderr.String(), "reached Stopped") {
		t.Fatalf("expected Stopped diagnostics on stderr, got %q", stderr.String())
	}
	expected := "task_01 pending — Build the widget core\n" +
		"task_02 pending — Wire the widget API\n" +
		"Stopped: 0 completed, 0 failed, 0 skipped, 2 pending.\n"
	if stdout.String() != expected {
		t.Fatalf("expected Stopped report:\n%q\ngot:\n%q", expected, stdout.String())
	}
	if verifier.calls != 0 || committer.calls != 0 {
		t.Fatalf("Stop Request must not verify or commit, got verify=%d commit=%d", verifier.calls, committer.calls)
	}
	run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if run.State != store.StateStopped {
		t.Fatalf("expected Stopped, got %q", run.State)
	}
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}

func TestRunImplementDatabaseStopRequestAfterTaskCommitEndsStoppedAndReleasesLock(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
		{id: "task_02", title: "Wire the widget API", needs: []string{"task_01"}},
	})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted, "task_02": spec.StatusCompleted},
	}
	committer, _, _, _ := withImplementCollaborators(t, runner)
	committer.afterCommit = func(context.Context, daemon.CommitRequest) error {
		return requestStopForActiveRunInGitRoot(homeDir, repoDir)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected stopped command exit 0, got %d (stderr %q)", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached Stopped") {
		t.Fatalf("expected Stopped diagnostics on stderr, got %q", stderr.String())
	}
	expected := "task_01 completed — Build the widget core\n" +
		"task_02 pending — Wire the widget API\n" +
		"Stopped: 1 completed, 0 failed, 0 skipped, 1 pending.\n"
	if stdout.String() != expected {
		t.Fatalf("expected Stopped report:\n%q\ngot:\n%q", expected, stdout.String())
	}
	if committer.calls != 1 {
		t.Fatalf("expected only task_01 committed before stop, got %d commit(s)", committer.calls)
	}
	run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if run.State != store.StateStopped {
		t.Fatalf("expected Stopped, got %q", run.State)
	}
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}

func requestStopForActiveRunInGitRoot(homeDir string, gitRoot string) error {
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		return err
	}
	defer func() {
		_ = runStore.Close()
	}()
	active, found, err := runStore.ActiveRunInGitRoot(ctx, gitRoot)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("no Active Run found for Stop Request")
	}
	return runStore.RequestStop(ctx, active.ID)
}

// implementJournaledQAEvent returns the Run's daemon.qa event from the Run
// Event Journal, and whether one was journaled at all.
func implementJournaledQAEvent(t *testing.T, homeDir string, runID string) (runevent.RunEvent, bool) {
	t.Helper()
	ctx := context.Background()
	reader, err := store.OpenReader(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run Database reader: %v", err)
	}
	defer func() {
		_ = reader.Close()
	}()
	events, err := reader.RunEventsAfter(ctx, runID, 0, 200)
	if err != nil {
		t.Fatalf("list journaled Run Events: %v", err)
	}
	for _, entry := range events {
		if entry.Event.Kind == runevent.KindDaemonQA {
			return entry.Event, true
		}
	}
	return runevent.RunEvent{}, false
}

func TestRunImplementQAVerdictMatrix(t *testing.T) {
	reportRel := implementQAReportRelPath()
	tests := []struct {
		name        string
		report      string
		wantCode    int
		wantVerdict string
		wantState   string
		wantCommit  bool
		wantDetail  string
	}{
		{name: "pass", report: implementQAReport("pass"), wantCode: 0, wantVerdict: "pass", wantState: store.StateClean, wantCommit: true, wantDetail: reportRel},
		{name: "partial", report: implementQAReport("partial"), wantCode: 1, wantVerdict: "partial", wantState: store.StateUnresolved, wantCommit: true, wantDetail: reportRel},
		{name: "fail", report: implementQAReport("fail"), wantCode: 1, wantVerdict: "fail", wantState: store.StateUnresolved, wantCommit: true, wantDetail: reportRel},
		{name: "missing report", report: "", wantCode: 1, wantVerdict: "missing", wantState: store.StateUnresolved, wantCommit: false, wantDetail: "no QA Report found"},
		{name: "unreadable verdict", report: "---\nsummary: no verdict field\n---\n\n# QA Report\n", wantCode: 1, wantVerdict: "unreadable", wantState: store.StateUnresolved, wantCommit: true, wantDetail: reportRel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
				{id: "task_01", title: "Build the widget core"},
			})
			runner := &implementFakeRunner{
				gitRoot:      repoDir,
				statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
				qaReport:     tt.report,
			}
			committer, _, _, _ := withImplementCollaborators(t, runner)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--qa", "--no-input"}, &stdout, &stderr)

			if code != tt.wantCode {
				t.Fatalf("expected exit code %d, got %d (stderr %q)", tt.wantCode, code, stderr.String())
			}
			outcomeLine := "Unresolved: 1 completed, 0 failed, 0 skipped, 0 pending.\n"
			if tt.wantState == store.StateClean {
				outcomeLine = "Clean: all 1 Task(s) completed.\n"
			}
			expected := "task_01 completed — Build the widget core\n" +
				"qa " + tt.wantVerdict + " — " + tt.wantDetail + "\n" +
				outcomeLine
			if stdout.String() != expected {
				t.Fatalf("expected QA report on stdout:\n%q\ngot:\n%q", expected, stdout.String())
			}
			if runner.qaCalls != 1 {
				t.Fatalf("expected one QA Agent invocation, got %d", runner.qaCalls)
			}
			wantCommits := 1
			if tt.wantCommit {
				wantCommits = 2
			}
			if committer.calls != wantCommits {
				t.Fatalf("expected %d commit(s), got %d (%v)", wantCommits, committer.calls, committer.messages)
			}
			if tt.wantCommit {
				wantMessage := "docs: qa report for " + implementTestSlug + " (" + tt.wantVerdict + ")\n\nRoundfix-Spec: " + implementTestSlug
				if committer.messages[len(committer.messages)-1] != wantMessage {
					t.Fatalf("expected the QA Report in its own commit %q, got %v", wantMessage, committer.messages)
				}
			}
			runID := implementRunIDFromStderr(t, stderr.String())
			run := implementRunFromStore(t, homeDir, runID)
			if run.State != tt.wantState {
				t.Fatalf("expected Run state %q, got %q", tt.wantState, run.State)
			}
			event, found := implementJournaledQAEvent(t, homeDir, runID)
			if !found {
				t.Fatalf("expected a daemon.qa event in the Run Event Journal, got none")
			}
			if !strings.Contains(string(event.Payload), fmt.Sprintf("%q", tt.wantVerdict)) {
				t.Fatalf("expected the verdict in the daemon.qa payload, got %s", event.Payload)
			}
			if tt.wantCommit && !strings.Contains(string(event.Payload), reportRel) {
				t.Fatalf("expected the report path in the daemon.qa payload, got %s", event.Payload)
			}
			assertNoActiveRunInGitRoot(t, homeDir, repoDir)
		})
	}
}

func TestRunImplementQAStepSkippedWhenAnyTaskFails(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
		{id: "task_02", title: "Wire the widget API", needs: []string{"task_01"}},
	})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusFailed},
		qaReport:     implementQAReport("pass"),
	}
	committer, _, _, _ := withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--qa", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d (stderr %q)", code, stderr.String())
	}
	expected := "task_01 failed — Build the widget core\n" +
		"  reason: Agent settled the Task failed\n" +
		"task_02 skipped — Wire the widget API\n" +
		"  reason: needs not completed: task_01\n" +
		"Unresolved: 0 completed, 1 failed, 1 skipped, 0 pending.\n"
	if stdout.String() != expected {
		t.Fatalf("expected no QA verdict line with a failed Task:\n%q\ngot:\n%q", expected, stdout.String())
	}
	if runner.qaCalls != 0 {
		t.Fatalf("expected the QA step never invoked with a failed Task, got %d QA call(s)", runner.qaCalls)
	}
	if committer.calls != 0 {
		t.Fatalf("expected no commits, got %d", committer.calls)
	}
	runID := implementRunIDFromStderr(t, stderr.String())
	if _, found := implementJournaledQAEvent(t, homeDir, runID); found {
		t.Fatal("expected no daemon.qa event when the QA step is skipped")
	}
}

func TestRunImplementQAOnlyRunSettlesOutcomeFromVerdict(t *testing.T) {
	tests := []struct {
		name      string
		verdict   string
		wantCode  int
		wantState string
	}{
		{name: "pass ends Clean", verdict: "pass", wantCode: 0, wantState: store.StateClean},
		{name: "fail ends Unresolved", verdict: "fail", wantCode: 1, wantState: store.StateUnresolved},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
				{id: "task_01", title: "Write the widget guide", status: string(spec.StatusCompleted)},
				{id: "task_02", title: "Build the widget backend", status: string(spec.StatusCompleted), needs: []string{"task_01"}},
			})
			runner := &implementFakeRunner{
				gitRoot:  repoDir,
				qaReport: implementQAReport(tt.verdict),
			}
			withImplementCollaborators(t, runner)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--qa", "--no-input"}, &stdout, &stderr)

			if code != tt.wantCode {
				t.Fatalf("expected exit code %d, got %d (stderr %q)", tt.wantCode, code, stderr.String())
			}
			outcomeLine := "Unresolved: 2 completed, 0 failed, 0 skipped, 0 pending.\n"
			if tt.wantState == store.StateClean {
				outcomeLine = "Clean: all 2 Task(s) completed.\n"
			}
			expected := "task_01 completed — Write the widget guide\n" +
				"task_02 completed — Build the widget backend\n" +
				"qa " + tt.verdict + " — " + implementQAReportRelPath() + "\n" +
				outcomeLine
			if stdout.String() != expected {
				t.Fatalf("expected QA-only report:\n%q\ngot:\n%q", expected, stdout.String())
			}
			if runner.calls != 1 || runner.qaCalls != 1 || len(runner.taskIDs) != 0 {
				t.Fatalf("expected a Run consisting of only the QA step, got calls=%d qaCalls=%d taskIDs=%v", runner.calls, runner.qaCalls, runner.taskIDs)
			}
			run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
			if run.Kind != store.KindImplement || run.State != tt.wantState {
				t.Fatalf("expected a QA-only Run reaching %q, got kind=%q state=%q", tt.wantState, run.Kind, run.State)
			}
			assertNoActiveRunInGitRoot(t, homeDir, repoDir)
		})
	}
}

func TestAttachReplaysCompletedSpecRunReadOnly(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Guide docs", taskType: "docs"},
		{id: "task_02", title: "Build core", needs: []string{"task_01"}},
	})
	runner := &implementFakeRunner{
		gitRoot: repoDir,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_02": spec.StatusCompleted,
		},
		qaReport: implementQAReport("pass"),
	}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--qa", "--no-input"}, &stdout, &stderr); code != 0 {
		t.Fatalf("seed implement run failed: %d stderr=%q", code, stderr.String())
	}
	runID := implementRunIDFromStderr(t, stderr.String())
	// Attach must never probe or start an Agent; a probing attach would
	// fail loudly here.
	withAgentRunner(t, &fakeAgentRunner{probeErr: errors.New("attach must not probe Agents")})
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	stdout.Reset()
	stderr.Reset()

	code := RunContext(context.Background(), []string{"attach", runID}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean attach exit, got %d stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"Roundfix attach",
		"Spec: " + implementTestSlug,
		"Branch: ma/widget-flow",
		"ID: " + runID,
		"State: Clean",
		"Tasks",
		"task_01 completed — Guide docs",
		"task_02 completed — Build core",
		"Task task_01 settled completed.",
		"QA verdict pass for Spec " + implementTestSlug + ".",
		"Run " + runID + " reached Clean; timeline replayed read-only.",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected attach output to contain %q, got:\n%s", expected, output)
		}
	}
	if strings.Contains(output, "Review Issues") {
		t.Fatalf("expected the Task pane instead of Review Issues, got:\n%s", output)
	}
	// Read-only: no new Runs, and the Run's terminal state is untouched.
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateClean {
		t.Fatalf("expected attach to leave the Run state untouched, got %q", run.State)
	}
}

func TestRunImplementInfrastructureFailureEndsFailed(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
	})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	committer, _, _, _ := withImplementCollaborators(t, runner)
	committer.err = errors.New("commit tooling broke")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d (stderr %q)", code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout on infrastructure failure, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "implement failed after Run start") {
		t.Fatalf("expected Run failure diagnostics, got %q", stderr.String())
	}
	run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if run.State != store.StateFailed {
		t.Fatalf("expected Failed, got %q", run.State)
	}
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}

func TestAgentSelectionProfilesMacro(t *testing.T) {
	binary := buildRoundfixBinaryForMacro(t)

	t.Run("mixed profiles configure validate fallback persist and stream", func(t *testing.T) {
		fake := newMacroFakeACPX(t)
		homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
			{id: "task_backend", title: "Build backend API", taskType: "backend"},
			{id: "task_frontend", title: "Build frontend view", taskType: "frontend"},
		})
		sentinels := seedMacroRuntimeOwnedFiles(t, homeDir)
		configPath := filepath.Join(repoDir, ".roundfixrc.yml")
		mustWrite(t, configPath, "worktree:\n  concurrency: 1\n")
		fragmentPath := filepath.Join(t.TempDir(), "profiles.yml")
		mustWrite(t, fragmentPath, macroProfilesYAML())

		stdout, stderr, code := runRoundfixBinaryMacro(t, binary, repoDir, homeDir, fake.binDir, "y\n", fake.env(), "profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json")
		if code != exitOK {
			t.Fatalf("profiles configure failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		configure := decodeMacroConfigureResponse(t, stdout)
		if configure.Schema != profilesConfigureSchema || !configure.Changed || configure.Scope != "project" {
			t.Fatalf("unexpected configure response: %#v stderr=%q", configure, stderr)
		}
		gitImplement(t, repoDir, "add", ".roundfixrc.yml")
		gitImplement(t, repoDir, "commit", "-m", "configure agent selection profiles")
		configBeforeRun := mustRead(t, configPath)

		showBefore := runProfilesShowMacro(t, binary, repoDir, homeDir, fake, "frontend")
		assertMacroFrontendProfileShow(t, showBefore)

		stdout, stderr, code = runRoundfixBinaryMacro(t, binary, repoDir, homeDir, fake.binDir, "", fake.env(), "profiles", "validate", "--json")
		if code != exitOK {
			t.Fatalf("profiles validate failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		validate := decodeMacroValidateResponse(t, stdout)
		if validate.Schema != profilesValidateSchema || !validate.OK {
			t.Fatalf("unexpected validate response: %#v stderr=%q", validate, stderr)
		}
		for _, invocation := range readMacroACPXLog(t, fake.logPath) {
			if invocation.Command == "prompt" {
				t.Fatalf("profiles validate must not prompt, got invocation %#v", invocation)
			}
		}

		env := fake.env()
		env["ROUNDFIX_FAKE_ACPX_FAIL_PREPARE_MODEL"] = "macro-frontend-preferred"
		stdout, stderr, code = runRoundfixBinaryMacro(t, binary, repoDir, homeDir, fake.binDir, "", env, "implement", "--spec", implementTestSlug, "--qa", "--no-input")
		if code != exitOK {
			t.Fatalf("implement macro failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		runID := implementRunIDFromStderr(t, stderr)
		for _, expected := range []string{
			"task_backend completed — Build backend API",
			"task_frontend completed — Build frontend view",
			"qa pass — docs/specs/0001-widget-flow/qa/qa-report-2026-01-01.md",
			"Clean: all 2 Task(s) completed.",
		} {
			if !strings.Contains(stdout, expected) {
				t.Fatalf("implement stdout missing %q:\n%s", expected, stdout)
			}
		}
		if !strings.Contains(stderr, "task task_frontend (frontend) Agent Selection failed") ||
			!strings.Contains(stderr, "activating fallback 1 claude/claude-fable-5/xhigh") {
			t.Fatalf("expected caller-visible fallback notification, got stderr:\n%s", stderr)
		}

		run := implementRunFromStore(t, homeDir, runID)
		if run.Agent != "codex" || run.Model != "macro-general" || run.ReasoningEffort != "high" {
			t.Fatalf("expected spec Run compatibility summary from general profile, got %#v", run)
		}
		assertMacroSelectionAttempts(t, homeDir, runID)
		assertMacroSelectionEventOrder(t, homeDir, runID)
		assertMacroSelectionStream(t, binary, repoDir, homeDir, fake, runID)
		assertMacroACPXActionLog(t, fake.logPath, runID)

		showAfter := runProfilesShowMacro(t, binary, repoDir, homeDir, fake, "frontend")
		if got, want := macroRecommendationSignature(showAfter), macroRecommendationSignature(showBefore); got != want {
			t.Fatalf("recommendations changed during run\nwant: %s\n got: %s", want, got)
		}
		if got, want := macroFallbackSignature(showAfter), macroFallbackSignature(showBefore); got != want {
			t.Fatalf("fallback order changed during run\nwant: %s\n got: %s", want, got)
		}
		if got := mustRead(t, configPath); got != configBeforeRun {
			t.Fatalf("implement mutated Project Config\nwant: %q\n got: %q", configBeforeRun, got)
		}
		assertMacroFilesUnchanged(t, sentinels)
	})

	t.Run("invalid task type blocks every proof and run side effect", func(t *testing.T) {
		fake := newMacroFakeACPX(t)
		homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
			{id: "task_01", title: "Invalid authoring", taskType: "Backend"},
		})

		stdout, stderr, code := runRoundfixBinaryMacro(t, binary, repoDir, homeDir, fake.binDir, "", fake.env(), "implement", "--spec", implementTestSlug, "--no-input")

		if code != exitPreflight {
			t.Fatalf("expected preflight exit %d, got %d stdout=%q stderr=%q", exitPreflight, code, stdout, stderr)
		}
		taskPath := implementTaskPath(repoDir, "task_01")
		for _, expected := range []string{taskPath, "Backend", "backend, frontend, data, infra, docs, test, chore", "update the task frontmatter type"} {
			if !strings.Contains(stderr, expected) {
				t.Fatalf("invalid Task Type stderr missing %q:\n%s", expected, stderr)
			}
		}
		if stdout != "" {
			t.Fatalf("expected no stdout for invalid Task Type, got %q", stdout)
		}
		if invocations := readMacroACPXLog(t, fake.logPath); len(invocations) != 0 {
			t.Fatalf("invalid Task Type must not probe or invoke Agent, got %#v", invocations)
		}
		assertNoRunDatabase(t, homeDir)
		assertNoMacroRunWorktreeOrBranch(t, homeDir, repoDir)
	})

	t.Run("post start failure never activates fallback", func(t *testing.T) {
		fake := newMacroFakeACPX(t)
		homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
			{id: "task_backend", title: "Backend fails after start", taskType: "backend"},
		})
		configPath := filepath.Join(repoDir, ".roundfixrc.yml")
		mustWrite(t, configPath, "worktree:\n  concurrency: 1\n")
		fragmentPath := filepath.Join(t.TempDir(), "profiles.yml")
		mustWrite(t, fragmentPath, macroProfilesYAML())
		stdout, stderr, code := runRoundfixBinaryMacro(t, binary, repoDir, homeDir, fake.binDir, "y\n", fake.env(), "profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json")
		if code != exitOK {
			t.Fatalf("profiles configure failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		gitImplement(t, repoDir, "add", ".roundfixrc.yml")
		gitImplement(t, repoDir, "commit", "-m", "configure agent selection profiles")

		env := fake.env()
		env["ROUNDFIX_FAKE_ACPX_FAIL_PROMPT_MODEL"] = "macro-backend"
		stdout, stderr, code = runRoundfixBinaryMacro(t, binary, repoDir, homeDir, fake.binDir, "", env, "implement", "--spec", implementTestSlug, "--no-input")

		if code != exitRunFailed {
			t.Fatalf("expected unresolved exit %d, got %d stdout=%q stderr=%q", exitRunFailed, code, stdout, stderr)
		}
		runID := implementRunIDFromStderr(t, stderr)
		if !strings.Contains(stdout, "task_backend failed — Backend fails after start") ||
			!strings.Contains(stdout, "Unresolved: 0 completed, 1 failed, 0 skipped, 0 pending.") {
			t.Fatalf("expected task failure report without fallback recovery, got stdout:\n%s", stdout)
		}
		for _, invocation := range readMacroACPXLog(t, fake.logPath) {
			if strings.HasPrefix(invocation.Session, "roundfix-preflight-") {
				continue
			}
			if invocation.Model == "macro-backend-fallback" || strings.Contains(invocation.Session, "fallback") {
				t.Fatalf("post-start failure must not activate fallback, got invocation %#v", invocation)
			}
		}
		for _, entry := range runEventsForRun(t, homeDir, runID) {
			if entry.Event.Kind == runevent.KindDaemonAgentSelectionFallback &&
				strings.Contains(string(entry.Event.Payload), `"scope_id":"task_backend"`) &&
				strings.Contains(string(entry.Event.Payload), `"next_selection"`) {
				t.Fatalf("post-start failure must not publish fallback notification, got event payload %s", string(entry.Event.Payload))
			}
		}
		attempts := agentSelectionAttemptsForRun(t, homeDir, runID)
		for _, attempt := range attempts {
			if attempt.ScopeID == "task_backend" && attempt.SelectionRole == store.AgentSelectionRoleFallback {
				t.Fatalf("post-start failure must not persist fallback attempt, got %#v", attempt)
			}
		}
	})
}

type macroFakeACPX struct {
	binDir    string
	codexPath string
	logPath   string
}

type macroACPXInvocation struct {
	Command    string   `json:"command"`
	Agent      string   `json:"agent"`
	Model      string   `json:"model"`
	Session    string   `json:"session"`
	CWD        string   `json:"cwd"`
	PromptKind string   `json:"prompt_kind"`
	PromptTask string   `json:"prompt_task"`
	Outcome    string   `json:"outcome"`
	Args       []string `json:"args"`
}

type macroCommandResult struct {
	Schema  string `json:"schema"`
	Changed bool   `json:"changed"`
	Scope   string `json:"scope"`
	OK      bool   `json:"ok"`
}

type macroProfilesShowResponse struct {
	Schema   string                     `json:"schema"`
	Profiles []macroProfilesShowProfile `json:"profiles"`
}

type macroProfilesShowProfile struct {
	Category        string                        `json:"category"`
	Preferred       roundconfig.AgentSelection    `json:"preferred"`
	Fallbacks       []roundconfig.AgentSelection  `json:"fallbacks"`
	Recommendations []macroProfilesRecommendation `json:"recommendations"`
}

type macroProfilesRecommendation struct {
	Rank      int                        `json:"rank"`
	Selection roundconfig.AgentSelection `json:"selection"`
}

type macroStreamRecord struct {
	Schema              string `json:"schema"`
	Category            string `json:"category"`
	Cursor              int64  `json:"cursor"`
	ScopeKind           string `json:"scope_kind"`
	ScopeID             string `json:"scope_id"`
	WorkCategory        string `json:"work_category"`
	ProfileSource       string `json:"profile_source"`
	Attempt             int    `json:"attempt"`
	SelectionRole       string `json:"selection_role"`
	FallbackIndex       int    `json:"fallback_index"`
	Runtime             string `json:"runtime"`
	Model               string `json:"model"`
	ReasoningEffort     string `json:"reasoning_effort"`
	NextRuntime         string `json:"next_runtime"`
	NextModel           string `json:"next_model"`
	NextReasoningEffort string `json:"next_reasoning_effort"`
	Status              string `json:"status"`
	ReasonCode          string `json:"reason_code"`
	Reason              string `json:"reason"`
}

func buildRoundfixBinaryForMacro(t *testing.T) string {
	t.Helper()
	repoRoot := repoRootForMacro(t)
	binary := filepath.Join(t.TempDir(), "roundfix")
	cmd := exec.Command("go", "build", "-buildvcs=false", "-o", binary, "./cmd/roundfix")
	cmd.Dir = repoRoot
	cmd.Env = isolatedGitEnvForTest()
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		t.Fatalf("build roundfix binary: %v\n%s", err, output.String())
	}
	return binary
}

func repoRootForMacro(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("git", append(gitConfigArgsForTest(), "rev-parse", "--show-toplevel")...)
	cmd.Env = isolatedGitEnvForTest()
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		t.Fatalf("resolve repo root: %v\n%s", err, output.String())
	}
	return strings.TrimSpace(output.String())
}

func newMacroFakeACPX(t *testing.T) macroFakeACPX {
	t.Helper()
	binDir := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "acpx.jsonl")
	script := strings.ReplaceAll(macroFakeACPXScript, "__PINNED_ACPX_VERSION__", agent.MinimumACPXVersion)
	acpxPath := filepath.Join(binDir, "acpx")
	if err := os.WriteFile(acpxPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake acpx: %v", err)
	}
	for _, adapter := range []string{"codex-acp", "claude-code-acp", "opencode", "npx"} {
		content := "#!/bin/sh\nexit 0\n"
		if adapter == "npx" {
			content = "#!/bin/sh\nprintf '%s\\n' '@agentclientprotocol/codex-acp " + agent.PinnedCodexAdapterVersion + "'\n"
		}
		if err := os.WriteFile(filepath.Join(binDir, adapter), []byte(content), 0o755); err != nil {
			t.Fatalf("write fake adapter %s: %v", adapter, err)
		}
	}
	codexPath := filepath.Join(binDir, "codex")
	buildMacroCodexExecutable(t, codexPath)
	assertMacroCodexExecutableHygiene(t, codexPath)
	return macroFakeACPX{binDir: binDir, codexPath: codexPath, logPath: logPath}
}

func (fake macroFakeACPX) env() map[string]string {
	return map[string]string{
		"CODEX_PATH":             fake.codexPath,
		"ROUNDFIX_FAKE_ACPX_LOG": fake.logPath,
	}
}

func buildMacroCodexExecutable(t *testing.T, destination string) {
	t.Helper()
	if runtime.GOOS != "darwin" {
		if err := os.WriteFile(destination, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
			t.Fatalf("write fake codex executable: %v", err)
		}
		return
	}
	sourceDir := t.TempDir()
	sourcePath := filepath.Join(sourceDir, "main.go")
	if err := os.WriteFile(sourcePath, []byte("package main\n\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatalf("write codex fixture source: %v", err)
	}
	cmd := exec.Command("go", "build", "-buildvcs=false", "-o", destination, sourcePath)
	cmd.Env = isolatedGitEnvForTest()
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		t.Fatalf("build signed codex fixture executable: %v\n%s", err, output.String())
	}
}

func assertMacroCodexExecutableHygiene(t *testing.T, path string) {
	t.Helper()
	result := codex.Inspector{ConfiguredPath: path}.Inspect(context.Background())
	if runtime.GOOS != "darwin" {
		if result.Status != codex.StatusSkipped {
			t.Fatalf("expected non-macOS codex hygiene to be skipped, got %#v", result)
		}
		return
	}
	if result.Status != codex.StatusOK {
		t.Fatalf("macro codex fixture must pass production hygiene, got status=%s detail=%q next=%q", result.Status, result.Detail, result.NextAction)
	}
}

const macroFakeACPXScript = `#!/usr/bin/env python3
import json
import os
import re
import sqlite3
import sys

PINNED = "__PINNED_ACPX_VERSION__"
AGENTS = {"codex", "claude", "opencode"}

def arg_value(argv, flag):
    for index, value in enumerate(argv):
        if value == flag and index + 1 < len(argv):
            return argv[index + 1]
    return ""

def parse(argv):
    if argv == ["--version"]:
        return {"command": "version", "args": argv}
    agent = ""
    agent_index = -1
    for index, value in enumerate(argv):
        if value in AGENTS:
            agent = value
            agent_index = index
            break
    tail = argv[agent_index + 1:] if agent_index >= 0 else []
    command = "unknown"
    session = ""
    if len(tail) >= 2 and tail[0] == "sessions" and tail[1] == "ensure":
        command = "sessions ensure"
        session = arg_value(tail, "--name")
    elif len(tail) >= 2 and tail[0] == "sessions" and tail[1] == "show":
        command = "sessions show"
        session = tail[2] if len(tail) > 2 else ""
    elif len(tail) >= 2 and tail[0] == "sessions" and tail[1] == "close":
        command = "sessions close"
        session = tail[2] if len(tail) > 2 else ""
    elif tail and tail[0] == "set":
        command = "set"
        session = arg_value(tail, "-s")
    elif tail and tail[0] == "prompt":
        command = "prompt"
        session = arg_value(tail, "-s")
    return {
        "command": command,
        "agent": agent,
        "model": arg_value(argv, "--model"),
        "session": session,
        "cwd": arg_value(argv, "--cwd") or os.getcwd(),
        "config_id": tail[1] if command == "set" and len(tail) > 1 else "",
        "config_value": tail[2] if command == "set" and len(tail) > 2 else "",
        "args": argv,
    }

def log(event):
    path = os.environ.get("ROUNDFIX_FAKE_ACPX_LOG", "")
    if not path:
        return
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "a", encoding="utf-8") as handle:
        handle.write(json.dumps(event, sort_keys=True) + "\n")

def session_model(session):
    path = os.environ.get("ROUNDFIX_FAKE_ACPX_LOG", "")
    if not path or not os.path.exists(path):
        return ""
    with open(path, "r", encoding="utf-8") as handle:
        records = handle.readlines()
    for raw in reversed(records):
        record = json.loads(raw)
        if record.get("session") == session and record.get("model"):
            return record["model"]
    return ""

def prompt_field(prompt, label):
    match = re.search(r"^" + re.escape(label) + r":\s*(.+)$", prompt, re.M)
    return match.group(1).strip() if match else ""

def resolve_path(cwd, value):
    value = value.strip()
    if os.path.isabs(value):
        return value
    return os.path.join(cwd, value)

def set_status(path, status):
    with open(path, "r", encoding="utf-8") as handle:
        lines = handle.readlines()
    if not lines or lines[0].strip() != "---":
        raise RuntimeError("missing task frontmatter")
    end = None
    for index in range(1, len(lines)):
        if lines[index].strip() == "---":
            end = index
            break
    if end is None:
        raise RuntimeError("missing task frontmatter terminator")
    replaced = False
    for index in range(1, end):
        if lines[index].startswith("status:"):
            lines[index] = "status: " + status + "\n"
            replaced = True
            break
    if not replaced:
        lines.insert(end, "status: " + status + "\n")
    with open(path, "w", encoding="utf-8") as handle:
        handle.writelines(lines)

def write_qa_report(spec_dir):
    report = os.path.join(spec_dir, "qa", "qa-report-2026-01-01.md")
    os.makedirs(os.path.dirname(report), exist_ok=True)
    with open(report, "w", encoding="utf-8") as handle:
        handle.write("---\nverdict: pass\n---\n\n# QA Report\n")

def require_durable_fallback_notification(event):
    session = event.get("session", "")
    if "-fallback-" not in session:
        return
    home = os.environ.get("HOME", "")
    db_path = os.path.join(home, ".roundfix", "roundfix.db")
    if not os.path.exists(db_path):
        sys.stderr.write("missing durable fallback notification database\n")
        sys.exit(1)
    try:
        connection = sqlite3.connect("file:" + db_path + "?mode=ro", uri=True, timeout=5)
        try:
            cursor = connection.execute(
                "SELECT COUNT(*) FROM run_events WHERE kind = ? AND payload LIKE ? AND payload LIKE ?",
                ("daemon.agent_selection_fallback", "%task_frontend%", "%claude-fable-5%"),
            )
            count = cursor.fetchone()[0]
        finally:
            connection.close()
    except Exception as exc:
        sys.stderr.write("read durable fallback notification failed: " + str(exc) + "\n")
        sys.exit(1)
    if count < 1:
        sys.stderr.write("fallback session prepared before durable fallback notification\n")
        sys.exit(1)

argv = sys.argv[1:]
event = parse(argv)
if event["command"] == "version":
    event["outcome"] = "ok"
    log(event)
    print(PINNED)
    sys.exit(0)

stdin_data = ""
if event["command"] == "prompt":
    stdin_data = sys.stdin.read()
    task = prompt_field(stdin_data, "Task")
    event["prompt_task"] = task
    event["prompt_kind"] = "task" if task else ("qa" if "Spec QA gate" in stdin_data else "unknown")

if event["command"] == "sessions ensure":
    require_durable_fallback_notification(event)
    fail_model = os.environ.get("ROUNDFIX_FAKE_ACPX_FAIL_PREPARE_MODEL", "")
    if event.get("model") == fail_model and not event.get("session", "").startswith("roundfix-preflight-"):
        event["outcome"] = "failed"
        log(event)
        sys.stderr.write("selection start rejected\n")
        sys.exit(1)

if event["command"] == "sessions show":
    model = session_model(event.get("session", ""))
    reasoning_id = "reasoning_effort" if event.get("agent") == "codex" else "effort"
    event["model"] = model
    event["outcome"] = "ok"
    log(event)
    print(json.dumps({
        "schema": "acpx.session.v1",
        "acpx": {
            "current_model_id": model,
            "config_options": [
                {"id": "model", "category": "model", "type": "select", "currentValue": model, "options": [{"value": model}]},
                {"id": reasoning_id, "type": "select", "currentValue": "medium", "options": [{"value": value} for value in ["low", "medium", "high", "xhigh", "max", "maximum", "ultra"]]},
            ],
        },
    }, sort_keys=True))
    sys.exit(0)

if event["command"] == "set":
    config_id = event.get("config_id", "")
    config_value = event.get("config_value", "")
    model = config_value if config_id == "model" else session_model(event.get("session", ""))
    reasoning_id = config_id if config_id in {"reasoning_effort", "effort"} else ("reasoning_effort" if event.get("agent") == "codex" else "effort")
    current_reasoning = config_value if config_id == reasoning_id else "medium"
    event["model"] = model
    event["outcome"] = "ok"
    log(event)
    print(json.dumps({
        "action": "config_set",
        "configId": config_id,
        "value": config_value,
        "configOptions": [
            {"id": "model", "category": "model", "type": "select", "currentValue": model, "options": [{"value": model}]},
            {"id": reasoning_id, "type": "select", "currentValue": current_reasoning, "options": [{"value": value} for value in ["low", "medium", "high", "xhigh", "max", "maximum", "ultra"]]},
        ],
    }, sort_keys=True))
    sys.exit(0)

if event["command"] == "prompt":
    fail_model = os.environ.get("ROUNDFIX_FAKE_ACPX_FAIL_PROMPT_MODEL", "")
    if event.get("model") == fail_model:
        event["outcome"] = "failed"
        log(event)
        sys.stderr.write("agent work rejected\n")
        sys.exit(1)
    if event["prompt_kind"] == "task":
        task_path = resolve_path(event["cwd"], prompt_field(stdin_data, "Task file"))
        set_status(task_path, "completed")
    elif event["prompt_kind"] == "qa":
        spec_dir = resolve_path(event["cwd"], prompt_field(stdin_data, "Spec directory"))
        write_qa_report(spec_dir)
    event["outcome"] = "ok"
    log(event)
    print('{"jsonrpc":"2.0","id":1,"result":{"stopReason":"end_turn"}}')
    sys.exit(0)

event["outcome"] = "ok"
log(event)
sys.exit(0)
`

func macroProfilesYAML() string {
	return `profiles:
  general:
    preferred:
      runtime: codex
      model: macro-general
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: macro-general-fallback
        reasoning_effort: max
  backend:
    preferred:
      runtime: codex
      model: macro-backend
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: macro-backend-fallback
        reasoning_effort: max
  frontend:
    preferred:
      runtime: codex
      model: macro-frontend-preferred
      reasoning_effort: high
    fallbacks:
      - runtime: claude
        model: claude-fable-5
        reasoning_effort: xhigh
  qa:
    preferred:
      runtime: codex
      model: macro-qa
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: macro-qa-fallback
        reasoning_effort: max
  review:
    preferred:
      runtime: codex
      model: macro-review
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: macro-review-fallback
        reasoning_effort: max
`
}

func runRoundfixBinaryMacro(t *testing.T, binary string, dir string, homeDir string, fakeBinDir string, stdin string, extraEnv map[string]string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Dir = dir
	env := isolatedGitEnvForTest()
	env = withEnvValue(env, "HOME", homeDir)
	env = withEnvValue(env, "PATH", fakeBinDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	env = withEnvValue(env, "ROUNDFIX_TUI", "never")
	env = withEnvValue(env, "ROUNDFIX_COLOR", "never")
	env = withEnvValue(env, "NO_COLOR", "1")
	env = withEnvValue(env, "TERM", "dumb")
	for key, value := range extraEnv {
		env = withEnvValue(env, key, value)
	}
	cmd.Env = env
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), exitCodeFromWait(err)
}

func decodeMacroConfigureResponse(t *testing.T, raw string) macroCommandResult {
	t.Helper()
	var response macroCommandResult
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		t.Fatalf("decode configure response %q: %v", raw, err)
	}
	return response
}

func decodeMacroValidateResponse(t *testing.T, raw string) macroCommandResult {
	t.Helper()
	var response macroCommandResult
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		t.Fatalf("decode validate response %q: %v", raw, err)
	}
	return response
}

func runProfilesShowMacro(t *testing.T, binary string, repoDir string, homeDir string, fake macroFakeACPX, category string) macroProfilesShowResponse {
	t.Helper()
	stdout, stderr, code := runRoundfixBinaryMacro(t, binary, repoDir, homeDir, fake.binDir, "", fake.env(), "profiles", "show", "--category", category, "--json")
	if code != exitOK {
		t.Fatalf("profiles show failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	var response macroProfilesShowResponse
	if err := json.Unmarshal([]byte(stdout), &response); err != nil {
		t.Fatalf("decode profiles show response %q: %v", stdout, err)
	}
	return response
}

func assertMacroFrontendProfileShow(t *testing.T, response macroProfilesShowResponse) {
	t.Helper()
	if response.Schema != profilesShowSchema || len(response.Profiles) != 1 {
		t.Fatalf("unexpected profiles show response: %#v", response)
	}
	profile := response.Profiles[0]
	if profile.Category != "frontend" {
		t.Fatalf("expected frontend profile, got %#v", profile)
	}
	if profile.Preferred.Runtime != "codex" || profile.Preferred.Model != "macro-frontend-preferred" || profile.Preferred.ReasoningEffort != "high" {
		t.Fatalf("unexpected frontend preferred selection: %#v", profile.Preferred)
	}
	if len(profile.Fallbacks) != 1 || profile.Fallbacks[0].Runtime != "claude" || profile.Fallbacks[0].Model != "claude-fable-5" || profile.Fallbacks[0].ReasoningEffort != "xhigh" {
		t.Fatalf("unexpected frontend fallback chain: %#v", profile.Fallbacks)
	}
	if len(profile.Recommendations) != 5 {
		t.Fatalf("expected exactly five recommendations, got %#v", profile.Recommendations)
	}
	seen := map[string]bool{}
	for index, recommendation := range profile.Recommendations {
		if recommendation.Rank != index+1 {
			t.Fatalf("recommendation rank order changed at index %d: %#v", index, recommendation)
		}
		key := recommendation.Selection.Runtime + "/" + recommendation.Selection.Model
		if seen[key] {
			t.Fatalf("duplicate recommendation tuple %s in %#v", key, profile.Recommendations)
		}
		seen[key] = true
	}
}

func macroRecommendationSignature(response macroProfilesShowResponse) string {
	if len(response.Profiles) == 0 {
		return ""
	}
	parts := make([]string, 0, len(response.Profiles[0].Recommendations))
	for _, recommendation := range response.Profiles[0].Recommendations {
		selection := recommendation.Selection
		parts = append(parts, fmt.Sprintf("%d:%s/%s/%s", recommendation.Rank, selection.Runtime, selection.Model, selection.ReasoningEffort))
	}
	return strings.Join(parts, ",")
}

func macroFallbackSignature(response macroProfilesShowResponse) string {
	if len(response.Profiles) == 0 {
		return ""
	}
	parts := make([]string, 0, len(response.Profiles[0].Fallbacks))
	for _, selection := range response.Profiles[0].Fallbacks {
		parts = append(parts, selection.Runtime+"/"+selection.Model+"/"+selection.ReasoningEffort)
	}
	return strings.Join(parts, ",")
}

func readMacroACPXLog(t *testing.T, path string) []macroACPXInvocation {
	t.Helper()
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatalf("read fake acpx log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) == 1 && strings.TrimSpace(lines[0]) == "" {
		return nil
	}
	invocations := make([]macroACPXInvocation, 0, len(lines))
	for _, line := range lines {
		var invocation macroACPXInvocation
		if err := json.Unmarshal([]byte(line), &invocation); err != nil {
			t.Fatalf("decode fake acpx log line %q: %v", line, err)
		}
		invocations = append(invocations, invocation)
	}
	return invocations
}

func assertMacroACPXActionLog(t *testing.T, path string, runID string) {
	t.Helper()
	invocations := readMacroACPXLog(t, path)
	prompts := map[string]macroACPXInvocation{}
	closed := map[string]bool{}
	for _, invocation := range invocations {
		if invocation.Command == "prompt" {
			prompts[invocation.PromptKind+":"+invocation.PromptTask] = invocation
		}
		if invocation.Command == "sessions close" {
			closed[invocation.Session] = true
		}
	}
	expectedPrompts := map[string]string{
		"task:task_backend":  "codex/macro-backend",
		"task:task_frontend": "claude/claude-fable-5",
		"qa:":                "codex/macro-qa",
	}
	for key, selection := range expectedPrompts {
		invocation, ok := prompts[key]
		if !ok {
			t.Fatalf("missing prompt invocation for %s in %#v", key, prompts)
		}
		if got := invocation.Agent + "/" + invocation.Model; got != selection {
			t.Fatalf("prompt %s used %s, want %s", key, got, selection)
		}
	}
	for _, session := range []string{
		"roundfix-" + runID + "-task_backend",
		"roundfix-" + runID + "-task_frontend",
		"roundfix-" + runID + "-task_frontend-fallback-01",
		"roundfix-" + runID + "-qa",
	} {
		if !closed[session] {
			t.Fatalf("expected session %s closed; closed sessions=%#v", session, closed)
		}
	}
}

func agentSelectionAttemptsForRun(t *testing.T, homeDir string, runID string) []store.AgentSelectionAttempt {
	t.Helper()
	runStore, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database reader: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close Run Database reader: %v", err)
		}
	}()
	attempts, err := runStore.AgentSelectionAttempts(context.Background(), runID)
	if err != nil {
		t.Fatalf("read Agent Selection attempts: %v", err)
	}
	return attempts
}

func assertMacroSelectionAttempts(t *testing.T, homeDir string, runID string) {
	t.Helper()
	attempts := agentSelectionAttemptsForRun(t, homeDir, runID)
	byScope := map[string][]store.AgentSelectionAttempt{}
	for _, attempt := range attempts {
		byScope[string(attempt.ScopeKind)+":"+attempt.ScopeID] = append(byScope[string(attempt.ScopeKind)+":"+attempt.ScopeID], attempt)
	}
	assertMacroAttemptSequence(t, byScope["task:task_backend"], []string{
		"1|backend|project|preferred|0|codex|macro-backend|high|active||",
		"2|backend|project|preferred|0|codex|macro-backend|high|closed||",
	})
	assertMacroAttemptSequence(t, byScope["task:task_frontend"], []string{
		"1|frontend|project|preferred|0|codex|macro-frontend-preferred|high|failed|selection_start_failed|selection start rejected",
		"2|frontend|project|fallback|1|claude|claude-fable-5|xhigh|active||",
		"3|frontend|project|fallback|1|claude|claude-fable-5|xhigh|closed||",
	})
	assertMacroAttemptSequence(t, byScope["qa:qa"], []string{
		"1|qa|project|preferred|0|codex|macro-qa|high|active||",
		"2|qa|project|preferred|0|codex|macro-qa|high|closed||",
	})
}

func assertMacroAttemptSequence(t *testing.T, attempts []store.AgentSelectionAttempt, expected []string) {
	t.Helper()
	got := make([]string, 0, len(attempts))
	for _, attempt := range attempts {
		reason := attempt.Reason
		if strings.Contains(reason, "selection start rejected") {
			reason = "selection start rejected"
		}
		got = append(got, fmt.Sprintf("%d|%s|%s|%s|%d|%s|%s|%s|%s|%s|%s",
			attempt.Attempt,
			attempt.Category,
			attempt.ProfileSource,
			attempt.SelectionRole,
			attempt.FallbackIndex,
			attempt.Runtime,
			attempt.Model,
			attempt.ReasoningEffort,
			attempt.Status,
			attempt.ReasonCode,
			reason,
		))
	}
	if strings.Join(got, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("selection attempts mismatch\nwant:\n%s\n got:\n%s", strings.Join(expected, "\n"), strings.Join(got, "\n"))
	}
}

func assertMacroSelectionEventOrder(t *testing.T, homeDir string, runID string) {
	t.Helper()
	events := runEventsForRun(t, homeDir, runID)
	failedAttempt := -1
	notification := -1
	fallbackActive := -1
	workStarted := -1
	for index, entry := range events {
		payload := string(entry.Event.Payload)
		switch {
		case entry.Event.Kind == runevent.KindDaemonAgentSelectionFallback &&
			strings.Contains(payload, `"scope_id":"task_frontend"`) &&
			strings.Contains(payload, `"attempt":1`) &&
			strings.Contains(payload, `"status":"failed"`):
			failedAttempt = index
		case entry.Event.Kind == runevent.KindDaemonAgentSelectionFallback &&
			strings.Contains(payload, `"scope_id":"task_frontend"`) &&
			strings.Contains(payload, `"next_selection"`):
			notification = index
		case entry.Event.Kind == runevent.KindDaemonAgentSelectionActive &&
			strings.Contains(payload, `"scope_id":"task_frontend"`) &&
			strings.Contains(payload, `"model":"claude-fable-5"`):
			fallbackActive = index
		case entry.Event.Kind == runevent.KindAgentStatus &&
			entry.Event.Batch == 2 &&
			strings.Contains(payload, agent.AgentWorkStartedStatus):
			workStarted = index
		}
	}
	if failedAttempt < 0 || notification < 0 || fallbackActive < 0 || workStarted < 0 {
		t.Fatalf("missing selection ordering events: failed=%d notification=%d active=%d workStarted=%d", failedAttempt, notification, fallbackActive, workStarted)
	}
	if !(failedAttempt < notification && notification < fallbackActive && fallbackActive < workStarted) {
		t.Fatalf("unexpected fallback ordering: failed=%d notification=%d active=%d workStarted=%d", failedAttempt, notification, fallbackActive, workStarted)
	}
}

func assertMacroSelectionStream(t *testing.T, binary string, repoDir string, homeDir string, fake macroFakeACPX, runID string) {
	t.Helper()
	stdout, stderr, code := runRoundfixBinaryMacro(t, binary, repoDir, homeDir, fake.binDir, "", fake.env(), "events", runID, "--filter", "agent-selection")
	if code != exitOK {
		t.Fatalf("events command failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	records := decodeMacroStreamRecords(t, stdout)
	if len(records) == 0 {
		t.Fatalf("expected Agent Selection stream records, got none")
	}
	forbidden := []string{"prompt", "credential", "token", "cookie", "secret"}
	lower := strings.ToLower(stdout)
	for _, value := range forbidden {
		if strings.Contains(lower, value) {
			t.Fatalf("selection stream leaked forbidden term %q:\n%s", value, stdout)
		}
	}
	notification := -1
	active := -1
	for index, record := range records {
		if record.Schema != runevent.StreamSchema || record.Category != string(runevent.StreamCategorySelection) {
			t.Fatalf("unexpected stream record: %#v", record)
		}
		if record.ScopeID == "task_frontend" && record.Attempt == 0 && record.NextModel == "claude-fable-5" {
			notification = index
			if record.Runtime != "codex" || record.Model != "macro-frontend-preferred" || record.NextRuntime != "claude" || record.Status != "failed" || record.ReasonCode != "selection_start_failed" {
				t.Fatalf("unexpected fallback notification stream record: %#v", record)
			}
		}
		if record.ScopeID == "task_frontend" && record.Attempt == 2 && record.SelectionRole == "fallback" && record.Status == "active" {
			active = index
			if record.Runtime != "claude" || record.Model != "claude-fable-5" || record.ProfileSource != "project" || record.FallbackIndex != 1 {
				t.Fatalf("unexpected fallback active stream record: %#v", record)
			}
		}
	}
	if notification < 0 || active < 0 || notification >= active {
		t.Fatalf("expected notification before fallback active in stream, notification=%d active=%d records=%#v", notification, active, records)
	}
}

func decodeMacroStreamRecords(t *testing.T, raw string) []macroStreamRecord {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) == 1 && strings.TrimSpace(lines[0]) == "" {
		return nil
	}
	records := make([]macroStreamRecord, 0, len(lines))
	for _, line := range lines {
		var record macroStreamRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("decode stream line %q: %v", line, err)
		}
		records = append(records, record)
	}
	return records
}

func seedMacroRuntimeOwnedFiles(t *testing.T, homeDir string) map[string]string {
	t.Helper()
	files := map[string]string{
		filepath.Join(homeDir, ".acpx", "runtime-owned.json"):    `{"runtime":"owned"}` + "\n",
		filepath.Join(homeDir, ".codex", "credentials.json"):     `{"credential":"sentinel"}` + "\n",
		filepath.Join(homeDir, ".claude", "credentials.json"):    `{"credential":"sentinel"}` + "\n",
		filepath.Join(homeDir, ".opencode", "credentials.json"):  `{"credential":"sentinel"}` + "\n",
		filepath.Join(homeDir, ".config", "roundfix-secret.txt"): "secret-sentinel\n",
	}
	for path, content := range files {
		mustMkdir(t, filepath.Dir(path))
		mustWrite(t, path, content)
	}
	return files
}

func assertMacroFilesUnchanged(t *testing.T, files map[string]string) {
	t.Helper()
	for path, expected := range files {
		if got := mustRead(t, path); got != expected {
			t.Fatalf("file %s changed\nwant: %q\n got: %q", path, expected, got)
		}
	}
}

func assertNoMacroRunWorktreeOrBranch(t *testing.T, homeDir string, repoDir string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(homeDir, ".roundfix", "worktrees")); err == nil {
		t.Fatalf("expected no Run Worktree root under %s", homeDir)
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stat Run Worktree root: %v", err)
	}
	branches := gitImplementOutput(t, repoDir, "branch", "--list", "roundfix/*")
	if strings.TrimSpace(branches) != "" {
		t.Fatalf("expected no roundfix Run branches, got %q", branches)
	}
}
