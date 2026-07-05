package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
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
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output.String())
	}
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
// way a real Agent edits the task file.
type implementFakeRunner struct {
	gitRoot      string
	statusByTask map[string]spec.Status
	errByTask    map[string]error
	probeErr     error
	calls        int
	taskIDs      []string
}

func (runner *implementFakeRunner) Probe(context.Context, agent.RuntimeSpec) error {
	return runner.probeErr
}

func (runner *implementFakeRunner) Run(_ context.Context, req agent.ExecuteRequest, _ runevent.Sink) (agent.ExecuteResult, error) {
	runner.calls++
	taskID := implementTaskIDFromPrompt(req.Prompt)
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
func withImplementCollaborators(t *testing.T, runner *implementFakeRunner) (*fakeCommitter, *fakeVerifier, *fakePusher, *fakeSourceResolver) {
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
		"--agent":             true,
		"--model":             true,
		"--agent-command":     true,
		"--agent-full-access": true,
		"--interactive":       true,
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
	if strings.Contains(stdout.String(), "--qa") {
		t.Fatalf("help must not list --qa before it ships:\n%s", stdout.String())
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
			name:    "missing spec",
			args:    []string{"implement", "--agent", "codex"},
			message: "missing required --spec",
		},
		{
			name:    "missing spec with no-input",
			args:    []string{"implement", "--agent", "codex", "--no-input"},
			message: "missing required --spec",
		},
		{
			name:    "missing spec with interactive",
			args:    []string{"implement", "--agent", "codex", "--interactive"},
			message: "missing required --spec",
		},
		{
			name:    "interactive with no-input",
			args:    []string{"implement", "--spec", implementTestSlug, "--interactive", "--no-input"},
			message: "--interactive cannot be used with --no-input",
		},
		{
			// The built-in config default is codex, so an empty Agent only
			// happens when the flag explicitly clears it.
			name:    "explicitly empty agent",
			args:    []string{"implement", "--spec", implementTestSlug, "--agent="},
			message: "missing required --agent",
		},
		{
			name:    "unsupported agent",
			args:    []string{"implement", "--spec", implementTestSlug, "--agent", "gemini"},
			message: `unsupported Agent "gemini"`,
		},
		{
			name:    "qa flag not shipped",
			args:    []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--qa"},
			message: "flag provided but not defined: -qa",
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
	if !strings.HasPrefix(committer.messages[0], "docs: Write the widget guide") || !strings.Contains(committer.messages[0], "Roundfix-Task: task_01") {
		t.Fatalf("expected frontmatter-derived task_01 commit message, got %q", committer.messages[0])
	}
	if !strings.HasPrefix(committer.messages[1], "feat: Build the widget backend") || !strings.Contains(committer.messages[1], "Roundfix-Spec: "+implementTestSlug) {
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
