package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
)

const implementTestSlug = "0001-widget-flow"

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

func writeImplementSpec(t *testing.T, repoDir string, slug string, seeds []implementSeed) {
	t.Helper()
	specDir := filepath.Join(repoDir, "docs", "specs", slug)
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
		mustWrite(t, implementTaskPath(repoDir, seed.id), implementTaskContent(slug, seed))
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
	return filepath.Join(repoDir, "docs", "specs", implementTestSlug, taskID+".md")
}

// implementFakeRunner scripts per-Task Agent behavior keyed by the Task id
// parsed from the prompt, writing task statuses through spec.SetStatus the
// way a real Agent edits the task file. A QA prompt writes qaReport as the
// Spec's QA Report, the way the qa-gate Agent does; an empty qaReport
// writes none.
type implementFakeRunner struct {
	gitRoot      string
	statusByTask map[string]spec.Status
	errByTask    map[string]error
	probeErr     error
	calls        int
	taskIDs      []string
	qaReport     string
	qaCalls      int
	logPaths     []string
	writeLogs    bool
	agentOutput  string
}

func (runner *implementFakeRunner) Probe(context.Context, agent.RuntimeSpec) error {
	return runner.probeErr
}

func (runner *implementFakeRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	runner.calls++
	runner.logPaths = append(runner.logPaths, req.LogPath)
	if runner.agentOutput != "" {
		if err := publishFakeAgentOutput(ctx, sink, req, runner.agentOutput); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if runner.writeLogs {
		if err := os.MkdirAll(filepath.Dir(req.LogPath), 0o755); err != nil {
			return agent.ExecuteResult{}, err
		}
		if err := os.WriteFile(req.LogPath, []byte("fake agent output\n"), 0o644); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	taskID := implementTaskIDFromPrompt(req.Prompt)
	if taskID == "" && strings.Contains(req.Prompt, "Spec QA gate") {
		runner.qaCalls++
		if runner.qaReport != "" {
			reportPath := filepath.Join(runner.gitRoot, implementQAReportRelPath())
			if err := os.MkdirAll(filepath.Dir(reportPath), 0o755); err != nil {
				return agent.ExecuteResult{}, err
			}
			if err := os.WriteFile(reportPath, []byte(runner.qaReport), 0o644); err != nil {
				return agent.ExecuteResult{}, err
			}
		}
		return agent.ExecuteResult{LogPath: req.LogPath}, nil
	}
	runner.taskIDs = append(runner.taskIDs, taskID)
	if err := runner.errByTask[taskID]; err != nil {
		return agent.ExecuteResult{}, err
	}
	if status, ok := runner.statusByTask[taskID]; ok {
		if err := spec.SetStatus(implementTaskPath(runner.gitRoot, taskID), status); err != nil {
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

// withImplementCollaborators wires the standard fake collaborators for an
// implement run and returns the ones tests assert on.
func withImplementCollaborators(t *testing.T, runner agent.Runner) (*fakeCommitter, *fakeVerifier, *fakePusher, *fakeSourceResolver) {
	t.Helper()
	committer := &fakeCommitter{}
	verifier := &fakeVerifier{}
	pusher := &fakePusher{}
	sourceResolver := &fakeSourceResolver{}
	withAgentRunner(t, runner)
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withFakeWorktree(t)
	withPusher(t, pusher)
	withSourceResolver(t, sourceResolver)
	return committer, verifier, pusher, sourceResolver
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
		"--agent-command":     true,
		"--agent-full-access": true,
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

func TestRunHelpListsImplementCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "roundfix implement --spec <slug> --agent <agent>") {
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
			// The workspace has no docs/specs/, so the Spec picker has
			// nothing to offer and fails with the fix instead.
			name:    "missing spec without active Specs",
			args:    []string{"implement", "--agent", "codex"},
			message: "no active Specs to implement",
		},
		{
			name:    "missing spec with no-input",
			args:    []string{"implement", "--agent", "codex", "--no-input"},
			message: "missing required --spec because --no-input disables Interactive Input",
		},
		{
			name:    "interactive without active Specs",
			args:    []string{"implement", "--agent", "codex", "--interactive"},
			message: "no active Specs to implement",
		},
		{
			name:    "interactive with no-input",
			args:    []string{"implement", "--spec", implementTestSlug, "--interactive", "--no-input"},
			message: "--interactive cannot be used with --no-input",
		},
		{
			name:    "interactive with no-agent-console",
			args:    []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--interactive", "--no-agent-console"},
			message: "--interactive cannot be used with --no-agent-console",
		},
		{
			// The built-in config default is codex, so an empty Agent only
			// happens when the flag explicitly clears it.
			name:    "explicitly empty agent with no-input",
			args:    []string{"implement", "--spec", implementTestSlug, "--agent=", "--no-input"},
			message: "missing required --agent because --no-input disables Interactive Input",
		},
		{
			name:    "unsupported agent",
			args:    []string{"implement", "--spec", implementTestSlug, "--agent", "gemini"},
			message: `unsupported Agent "gemini"`,
		},
		{
			name:    "unexpected argument",
			args:    []string{"implement", implementTestSlug, "--agent", "codex"},
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

func TestRunImplementInteractiveInputPicksSpecThroughCollector(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
	})
	mustMkdir(t, filepath.Join(repoDir, "docs", "specs", "0002-broken-prd"))
	mustWrite(t, filepath.Join(repoDir, "docs", "specs", "0002-broken-prd", "_prd.md"), "no frontmatter\n")
	gitImplement(t, repoDir, "add", "docs/specs/0002-broken-prd/_prd.md")
	gitImplement(t, repoDir, "commit", "-m", "add broken spec fixture")
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
		return roundtui.CollectInput(ctx, req, strings.NewReader("1\nclaude\n"), &collected)
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

func TestRunImplementInteractiveInputMergesQAGateChoice(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		input       string
		qaReport    string
		wantQACalls int
	}{
		{name: "scripted yes produces QA Run", args: []string{"implement"}, input: "1\ncodex\ny\n", qaReport: implementQAReport("pass"), wantQACalls: 1},
		{name: "empty input produces non-QA Run", args: []string{"implement"}, input: "1\ncodex\n\n", wantQACalls: 0},
		{name: "qa flag preset keeps QA on with enter", args: []string{"implement", "--qa"}, input: "1\ncodex\n\n", qaReport: implementQAReport("pass"), wantQACalls: 1},
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

func TestRunImplementInteractiveInputRemembersAgentButNotSpecOrQA(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
	})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		qaReport:     implementQAReport("pass"),
	}
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

	// A second invocation with the Agent explicitly cleared reopens the
	// flow: the remembered Agent surfaces as the suggestion, the Spec does
	// not — each Run's target is an explicit choice.
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "keep first run")
	var secondReq roundtui.InputRequest
	var secondCollected strings.Builder
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		secondReq = req
		return roundtui.CollectInput(context.Background(), req, strings.NewReader("1\n\n\n"), &secondCollected)
	})
	stdout.Reset()
	stderr.Reset()

	if code := RunContext(context.Background(), []string{"implement", "--agent="}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected second run exit 0, got %d (stderr %q)", code, stderr.String())
	}
	if secondReq.AgentSuggestion.Value != "claude" || secondReq.AgentSuggestion.Source != "remembered" {
		t.Fatalf("expected the remembered Agent suggestion, got %#v", secondReq.AgentSuggestion)
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

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--interactive"}, &stdout, &stderr)

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

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}, &stdout, &stderr)

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
		t.Fatalf("spec Runs must never push or resolve Review Source threads, got push=%d source=%d", pusher.calls, sourceResolver.calls)
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

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-agent-console", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean implement exit 0, got %d stderr=%q", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "fake implement agent output") {
		t.Fatalf("expected Agent console hidden from stderr, got %q", stderr.String())
	}
	for _, want := range []string{
		"Implement Run:",
		"implement selected Spec",
		"Verification command passed",
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
			mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), fmt.Sprintf("defaults:\n  artifact_dir: %q\n", tt.config))
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

			code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}, &stdout, &stderr)

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

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}, &stdout, &stderr)

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
			wantCode:  0,
			wantState: store.StateClean,
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

			code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}, &stdout, &stderr)

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
			args:     []string{"implement", "--spec", "9999-missing", "--agent", "codex"},
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
			name:  "dirty working tree",
			seeds: []implementSeed{{id: "task_01"}},
			mutate: func(t *testing.T, repoDir string) {
				mustWrite(t, filepath.Join(repoDir, "leftover.txt"), "failed attempt leftovers\n")
			},
			messages: []string{"working tree", "not clean", "commit, stash, or discard"},
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
				args = []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}
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
	})
	if err != nil {
		t.Fatalf("seed blocking run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(ctx, []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d (stderr %q)", code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), blocking.ID) || !strings.Contains(stderr.String(), "roundfix stop "+blocking.ID) {
		t.Fatalf("expected the blocking run id and stop command, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunImplementPreflightProbeFailureCreatesNoRun(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	withImplementCollaborators(t, &implementFakeRunner{gitRoot: repoDir, probeErr: errors.New("codex-acp is not on PATH")})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d (stderr %q)", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "codex-acp is not on PATH") {
		t.Fatalf("expected the probe failure reason, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 0)
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

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}, &stdout, &stderr)

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

func TestRunImplementFailedTaskEndsUnresolvedAndResumeFinishesGraph(t *testing.T) {
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

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d (stderr %q)", code, stderr.String())
	}
	expected := "task_01 completed — Build the widget core\n" +
		"task_02 failed — Wire the widget API\n" +
		"task_03 skipped — Document the widget\n" +
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
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)

	// Resume: clear the failed attempt's leftovers, as the dirty-tree
	// preflight message instructs, then re-run the same command.
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "keep failed attempt")
	secondRunner := &implementFakeRunner{
		gitRoot: repoDir,
		statusByTask: map[string]spec.Status{
			"task_02": spec.StatusCompleted,
			"task_03": spec.StatusCompleted,
		},
	}
	withImplementCollaborators(t, secondRunner)
	stdout.Reset()
	stderr.Reset()

	code = RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected resume exit code 0, got %d (stderr %q)", code, stderr.String())
	}
	if secondRunner.calls != 2 || secondRunner.taskIDs[0] != "task_02" || secondRunner.taskIDs[1] != "task_03" {
		t.Fatalf("expected the resume Run to execute only non-completed Tasks, got %v", secondRunner.taskIDs)
	}
	expected = "task_01 completed — Build the widget core\n" +
		"task_02 completed — Wire the widget API\n" +
		"task_03 completed — Document the widget\n" +
		"Clean: all 3 Task(s) completed.\n"
	if stdout.String() != expected {
		t.Fatalf("expected Clean resume report:\n%q\ngot:\n%q", expected, stdout.String())
	}
	secondRun := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if secondRun.State != store.StateClean {
		t.Fatalf("expected resume Run Clean, got %q", secondRun.State)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 2)
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

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}, &stdout, &stderr)

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

	stoppedCode := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}, &stdout, &stderr)
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

			code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--qa", "--no-input"}, &stdout, &stderr)

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

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--qa", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d (stderr %q)", code, stderr.String())
	}
	expected := "task_01 failed — Build the widget core\n" +
		"task_02 skipped — Wire the widget API\n" +
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

			code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--qa", "--no-input"}, &stdout, &stderr)

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
	if code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--qa", "--no-input"}, &stdout, &stderr); code != 0 {
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

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}, &stdout, &stderr)

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
