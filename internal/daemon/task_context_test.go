package daemon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/spec"
	runworktree "roundfix/internal/worktree"
)

func TestAssembleTaskContextBundleReservesExplicitPathsAndCountsOmittedPriorFiles(t *testing.T) {
	workDir := t.TempDir()
	writeBundleStandardFiles(t, workDir, filepath.Join(workDir, "docs", "specs"), taskCycleSlug)
	task := spec.Task{
		ID: "task_01",
		Context: []spec.TaskContextRef{
			{Kind: spec.ContextKindInstruction, Path: ".agents/skills/golang-testing/SKILL.md"},
			{Kind: spec.ContextKindInterface, Path: "internal/spec/task.go"},
		},
	}
	prior := make([]string, 0, 205)
	for index := 204; index >= 0; index-- {
		prior = append(prior, fmt.Sprintf("prior/file_%03d.go", index))
	}

	bundle := assembleTaskContextBundle(TaskPlan{
		WorkDir:   workDir,
		SpecsRoot: filepath.Join(workDir, "docs", "specs"),
		Spec:      spec.Spec{Slug: taskCycleSlug, Dir: filepath.Join(workDir, "docs", "specs", taskCycleSlug)},
	}, task, prior)

	if got := bundlePathCount(bundle); got != specContextBundlePathLimit {
		t.Fatalf("bundle path count = %d, want %d", got, specContextBundlePathLimit)
	}
	for _, expected := range []string{
		"docs/specs/" + taskCycleSlug + "/_prd.md",
		"docs/specs/" + taskCycleSlug + "/_techspec.md",
		"docs/specs/" + taskCycleSlug + "/_tasks.md",
		"AGENTS.md",
		".agents/skills/implement-task/SKILL.md",
		".agents/skills/golang-testing/SKILL.md",
		"internal/spec/task.go",
	} {
		if !bundleContainsPath(bundle, expected) {
			t.Fatalf("expected reserved path %q in bundle, got %+v", expected, bundle)
		}
	}
	if len(bundle.PriorChangedFiles) != specContextBundlePathLimit-7 {
		t.Fatalf("prior file count = %d, want %d", len(bundle.PriorChangedFiles), specContextBundlePathLimit-7)
	}
	if bundle.OmittedPriorFiles != 12 {
		t.Fatalf("OmittedPriorFiles = %d, want 12", bundle.OmittedPriorFiles)
	}
	if bundle.PriorChangedFiles[0] != "prior/file_000.go" {
		t.Fatalf("prior files not sorted, first = %q", bundle.PriorChangedFiles[0])
	}
	if got := bundle.PriorChangedFiles[len(bundle.PriorChangedFiles)-1]; got != "prior/file_192.go" {
		t.Fatalf("prior cap ended at %q, want prior/file_192.go", got)
	}
}

func TestAssembleTaskContextBundleSupportsExternalSpecRoot(t *testing.T) {
	workDir := t.TempDir()
	externalSpecsRoot := filepath.Join(t.TempDir(), "external-specs")
	writeBundleStandardFiles(t, workDir, externalSpecsRoot, taskCycleSlug)

	bundle := assembleTaskContextBundle(TaskPlan{
		WorkDir:   workDir,
		SpecsRoot: externalSpecsRoot,
		Spec:      spec.Spec{Slug: taskCycleSlug, Dir: filepath.Join(externalSpecsRoot, taskCycleSlug)},
	}, spec.Task{ID: "task_01"}, nil)

	wantPRD := filepath.ToSlash(filepath.Join(externalSpecsRoot, taskCycleSlug, "_prd.md"))
	if bundle.PRD != wantPRD {
		t.Fatalf("PRD path = %q, want external path %q", bundle.PRD, wantPRD)
	}
	if !bundleContainsPath(bundle, "AGENTS.md") {
		t.Fatalf("expected repository instruction path retained, got %+v", bundle)
	}
}

func TestPriorChangedFilesUseCurrentWorktreeHeadAndIgnoreSiblingBranch(t *testing.T) {
	ctx := context.Background()
	repoDir := t.TempDir()
	runGitForTest(t, repoDir, "init", "-b", "main")
	mustWriteForTest(t, filepath.Join(repoDir, "base.go"), "package demo\n")
	runGitForTest(t, repoDir, "add", "base.go")
	runGitForTest(t, repoDir, "commit", "-m", "initial")
	initialHead := strings.TrimSpace(runGitForTest(t, repoDir, "rev-parse", "HEAD"))

	mustWriteForTest(t, filepath.Join(repoDir, "integrated.go"), "package demo\n")
	runGitForTest(t, repoDir, "add", "integrated.go")
	runGitForTest(t, repoDir, "commit", "-m", "integrated task")

	taskWorktree := filepath.Join(t.TempDir(), "task-worktree")
	runGitForTest(t, repoDir, "worktree", "add", "-b", "task-worktree", taskWorktree, "HEAD")
	siblingWorktree := filepath.Join(t.TempDir(), "sibling-worktree")
	runGitForTest(t, repoDir, "worktree", "add", "-b", "sibling-worktree", siblingWorktree, initialHead)
	mustWriteForTest(t, filepath.Join(siblingWorktree, "sibling.go"), "package demo\n")
	runGitForTest(t, siblingWorktree, "add", "sibling.go")
	runGitForTest(t, siblingWorktree, "commit", "-m", "unintegrated sibling")

	got, err := runworktree.PriorChangedFiles(ctx, taskWorktree, initialHead)
	if err != nil {
		t.Fatalf("PriorChangedFiles: %v", err)
	}
	if strings.Join(got, "|") != "integrated.go" {
		t.Fatalf("prior changed files = %v, want only integrated.go", got)
	}
}

func TestTaskCyclePromptContainsBundleWithoutReferencedBodies(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Use context"}})
	writePromptReferenceFiles(t, fixture.gitRoot, fixture.specsRoot)
	taskContent := markdownForDaemonTest(`---
task: task_01
spec: 0001-sample-feature
status: pending
type: backend
---

# Task 01: Use context

## Context

- instruction: '.agents/skills/golang-testing/SKILL.md'
- interface: 'internal/source.go'

## Verification

- 'true' — expected: passes.
`)
	mustWriteForTest(t, taskPathFor(fixture.gitRoot, taskCycleSlug, "task_01"), taskContent)
	graph, err := spec.Load(fixture.specsRoot, taskCycleSlug)
	if err != nil {
		t.Fatalf("reload graph: %v", err)
	}
	fixture.graph = graph
	prior := &fakePriorChangedResolver{byWork: map[string][]string{
		fixture.gitRoot: []string{"internal/prior.go", "internal/source.go"},
	}}
	runner := &taskFakeRunner{
		calls:   fixture.calls,
		gitRoot: fixture.gitRoot,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
		},
	}
	engine := fixture.engineWithTaskWorktreesAndPriorChanges(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree, nil, prior)
	plan := fixture.plan()
	plan.HeadSHA = "initial-head"

	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	if result.Completed != 1 {
		t.Fatalf("expected task completed, got %+v", result)
	}
	if len(runner.prompts) != 1 {
		t.Fatalf("prompts = %d, want 1", len(runner.prompts))
	}
	prompt := runner.prompts[0]
	daemonOwnedTaskContent := strings.Replace(taskContent, "status: pending", "status: in_progress", 1)
	if got := strings.Count(prompt, daemonOwnedTaskContent); got != 1 {
		t.Fatalf("expected one complete assigned Task, got %d", got)
	}
	for _, expected := range []string{
		"Spec Context Bundle:",
		"- PRD: docs/specs/" + taskCycleSlug + "/_prd.md",
		"- TechSpec: docs/specs/" + taskCycleSlug + "/_techspec.md",
		"- Task Graph: docs/specs/" + taskCycleSlug + "/_tasks.md",
		"  - AGENTS.md",
		"  - .agents/skills/implement-task/SKILL.md",
		"  - .agents/skills/golang-testing/SKILL.md",
		"  - internal/source.go",
		"  - internal/prior.go",
		"- Omitted prior files: 0",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("expected prompt to contain %q, got:\n%s", expected, prompt)
		}
	}
	for _, forbidden := range []string{
		"PRD BODY SENTINEL",
		"TECHSPEC BODY SENTINEL",
		"AGENTS BODY SENTINEL",
		"IMPLEMENT TASK SKILL SENTINEL",
		"GOLANG TESTING SKILL SENTINEL",
		"SOURCE BODY SENTINEL",
		"diff --git",
		"@@ -",
	} {
		if strings.Contains(prompt, forbidden) {
			t.Fatalf("prompt exposed referenced body marker %q:\n%s", forbidden, prompt)
		}
	}
	calls := prior.seenCalls()
	if len(calls) != 1 || calls[0].workDir != fixture.gitRoot || calls[0].initialHead != "initial-head" {
		t.Fatalf("prior resolver calls = %+v, want one call for fixture root and initial head", calls)
	}
}

func TestTaskCycleParallelTaskPromptUsesTaskWorktreeContextBase(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Integrate first"},
		{id: "task_03", title: "Independent sibling"},
		{id: "task_02", title: "Depends on first", needs: []string{"task_01"}},
	})
	writeBundleStandardFiles(t, fixture.gitRoot, fixture.specsRoot, taskCycleSlug)
	taskWorktrees := newFakeTaskWorktrees()
	prior := &fakePriorChangedResolver{
		resolve: func(workDir string) []string {
			switch {
			case strings.HasSuffix(workDir, ".task_02"):
				return []string{"integrated/task_01.go"}
			case strings.HasSuffix(workDir, ".task_03"):
				return []string{"unintegrated/task_03.go"}
			default:
				return nil
			}
		},
	}
	runner := newTaskSchedulerRunner("task_01", "task_02", "task_03")
	runner.onStart = func(taskID string, req agent.ExecuteRequest) error {
		if taskID != "task_02" {
			return nil
		}
		if req.GitRoot == fixture.gitRoot {
			return errors.New("task_02 did not run in its Task Worktree")
		}
		if !strings.Contains(req.Prompt, "  - integrated/task_01.go") {
			return fmt.Errorf("task_02 prompt missing integrated prior path:\n%s", req.Prompt)
		}
		if strings.Contains(req.Prompt, "unintegrated/task_03.go") {
			return fmt.Errorf("task_02 prompt included unintegrated sibling path:\n%s", req.Prompt)
		}
		return nil
	}
	verifier := &taskSchedulerVerifier{}
	committer := &taskSchedulerCommitter{}
	engine := fixture.engineWithTaskWorktreesAndPriorChanges(t, runner, verifier, committer, fixture.worktree, taskWorktrees, prior)
	plan := fixture.plan()
	plan.Concurrency = 2
	plan.HeadSHA = "run-initial-head"

	resultCh := make(chan struct {
		result TaskCycleResult
		err    error
	}, 1)
	go func() {
		result, err := engine.TaskCycle(context.Background(), plan)
		resultCh <- struct {
			result TaskCycleResult
			err    error
		}{result: result, err: err}
	}()

	assertTaskSet(t, waitSchedulerStarts(t, runner, 2), "task_01", "task_03")
	runner.releaseTask("task_01")
	if got := waitIntegratedTask(t, taskWorktrees); got != "task_01" {
		t.Fatalf("expected task_01 integrated before task_02, got %s", got)
	}
	if got := strings.Join(waitSchedulerStarts(t, runner, 1), "|"); got != "task_02" {
		t.Fatalf("expected task_02 to start after task_01 integration, got %s", got)
	}
	runner.releaseTask("task_02")
	runner.releaseTask("task_03")

	outcome := waitTaskCycleResult(t, resultCh)
	if outcome.err != nil {
		t.Fatalf("TaskCycle: %v", outcome.err)
	}
	if outcome.result.Completed != 3 || outcome.result.Failed != 0 || outcome.result.Skipped != 0 {
		t.Fatalf("expected all Tasks completed, got %+v", outcome.result)
	}
	for _, call := range prior.seenCalls() {
		if call.initialHead != "run-initial-head" {
			t.Fatalf("prior call used initial head %q, want run-initial-head", call.initialHead)
		}
	}
}

func TestAssembleTaskContextBundleIsDeterministic(t *testing.T) {
	workDir := t.TempDir()
	writeBundleStandardFiles(t, workDir, filepath.Join(workDir, "docs", "specs"), taskCycleSlug)
	plan := TaskPlan{
		WorkDir:   workDir,
		SpecsRoot: filepath.Join(workDir, "docs", "specs"),
		Spec:      spec.Spec{Slug: taskCycleSlug, Dir: filepath.Join(workDir, "docs", "specs", taskCycleSlug)},
	}
	task := spec.Task{
		ID: "task_01",
		Context: []spec.TaskContextRef{
			{Kind: spec.ContextKindInterface, Path: "internal/b.go"},
			{Kind: spec.ContextKindInstruction, Path: ".agents/skills/golang-testing/SKILL.md"},
		},
	}
	prior := []string{"internal/z.go", "internal/a.go", "internal/b.go"}

	first := assembleTaskContextBundle(plan, task, prior)
	second := assembleTaskContextBundle(plan, task, prior)
	if fmt.Sprintf("%#v", first) != fmt.Sprintf("%#v", second) {
		t.Fatalf("bundle was not deterministic:\nfirst:  %#v\nsecond: %#v", first, second)
	}
	if strings.Join(first.PriorChangedFiles, "|") != "internal/a.go|internal/z.go" {
		t.Fatalf("prior files = %v, want sorted files excluding explicit duplicate", first.PriorChangedFiles)
	}
}

func writeBundleStandardFiles(t *testing.T, workDir string, specsRoot string, slug string) {
	t.Helper()
	for path, content := range map[string]string{
		filepath.Join(workDir, "AGENTS.md"):                                       "AGENTS BODY SENTINEL\n",
		filepath.Join(workDir, ".agents", "skills", "implement-task", "SKILL.md"): "IMPLEMENT TASK SKILL SENTINEL\n",
		filepath.Join(specsRoot, slug, "_prd.md"):                                 "PRD BODY SENTINEL\n",
		filepath.Join(specsRoot, slug, "_techspec.md"):                            "TECHSPEC BODY SENTINEL\n",
		filepath.Join(specsRoot, slug, "_tasks.md"):                               "TASK GRAPH BODY SENTINEL\n",
	} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("create parent for %s: %v", path, err)
		}
		mustWriteForTest(t, path, content)
	}
}

func writePromptReferenceFiles(t *testing.T, workDir string, specsRoot string) {
	t.Helper()
	for path, content := range map[string]string{
		filepath.Join(workDir, "AGENTS.md"):                                       "AGENTS BODY SENTINEL\n",
		filepath.Join(workDir, ".agents", "skills", "implement-task", "SKILL.md"): "IMPLEMENT TASK SKILL SENTINEL\n",
		filepath.Join(specsRoot, taskCycleSlug, "_prd.md"):                        "---\nstatus: active\n---\n\nPRD BODY SENTINEL\n",
		filepath.Join(specsRoot, taskCycleSlug, "_techspec.md"):                   "TECHSPEC BODY SENTINEL\n",
		filepath.Join(workDir, ".agents", "skills", "golang-testing", "SKILL.md"): "GOLANG TESTING SKILL SENTINEL\n",
		filepath.Join(workDir, "internal", "source.go"):                           "package source\n\nconst marker = \"SOURCE BODY SENTINEL\"\n",
		filepath.Join(workDir, "internal", "prior.go"):                            "package prior\n",
	} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("create parent for %s: %v", path, err)
		}
		mustWriteForTest(t, path, content)
	}
}

func bundlePathCount(bundle agent.SpecContextBundle) int {
	count := 0
	for _, path := range []string{bundle.PRD, bundle.TechSpec, bundle.TaskGraph} {
		if strings.TrimSpace(path) != "" {
			count++
		}
	}
	count += len(bundle.Instructions)
	count += len(bundle.Interfaces)
	count += len(bundle.PriorChangedFiles)
	return count
}

func bundleContainsPath(bundle agent.SpecContextBundle, path string) bool {
	path = filepath.ToSlash(filepath.Clean(path))
	paths := []string{bundle.PRD, bundle.TechSpec, bundle.TaskGraph}
	paths = append(paths, bundle.Instructions...)
	paths = append(paths, bundle.Interfaces...)
	paths = append(paths, bundle.PriorChangedFiles...)
	sort.Strings(paths)
	for _, candidate := range paths {
		if candidate == path {
			return true
		}
	}
	return false
}

func markdownForDaemonTest(fixture string) string {
	return strings.ReplaceAll(fixture, "'", "`")
}
