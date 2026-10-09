package daemon

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func TestATaskThatLeavesItsDeletesPathDoesNotSettle(t *testing.T) {
	t.Parallel()
	testDeletesTaskCycle(t, false)
}

func TestATaskThatRemovesItsDeletesPathSettles(t *testing.T) {
	t.Parallel()
	testDeletesTaskCycle(t, true)
}

func testDeletesTaskCycle(t *testing.T, repairRemoves bool) {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"echo verified"}}})
	taskPath := taskPathFor(fixture.gitRoot, taskCycleSlug, "task_01")
	content, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, taskPath, string(content)+"\n## Context\n\n- deletes: `old/file.txt`\n")
	oldPath := filepath.Join(fixture.gitRoot, "old", "file.txt")
	if err := os.MkdirAll(filepath.Dir(oldPath), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, oldPath, "obsolete")
	runGitForTest(t, fixture.gitRoot, "add", ".")
	runGitForTest(t, fixture.gitRoot, "commit", "-qm", "test: deletion declaration")
	fixture.reloadGraph()
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
	calls := 0
	runner.afterTask = func(string) {
		calls++
		if repairRemoves && calls == 2 {
			if err := os.Remove(oldPath); err != nil {
				t.Fatal(err)
			}
		}
	}
	verifier := &taskFakeVerifier{calls: fixture.calls}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)
	result, err := engine.TaskCycle(context.Background(), fixture.plan())
	if err != nil {
		t.Fatal(err)
	}
	if len(runner.requests) != 2 || runner.requests[0].Session != runner.requests[1].Session {
		t.Fatalf("repair requests = %#v", runner.requests)
	}
	if !strings.Contains(runner.requests[1].Prompt, "deletes: old/file.txt still exists") {
		t.Fatalf("feedback = %s", runner.requests[1].Prompt)
	}
	if !strings.Contains(runner.requests[0].Prompt, "old/file.txt") {
		t.Fatal("deletion missing from Agent context")
	}
	wantStatus := string(spec.StatusFailed)
	if repairRemoves {
		wantStatus = string(spec.StatusCompleted)
		if result.Completed != 1 || result.Failed != 0 || len(verifier.commands) != 1 {
			t.Fatalf("removed path result = %+v, commands %v", result, verifier.commands)
		}
	} else if result.Failed != 1 || result.Completed != 0 || len(verifier.commands) != 0 || len(committer.messages) != 0 {
		t.Fatalf("remaining path result = %+v, commands %v, commits %v", result, verifier.commands, committer.messages)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != wantStatus {
		t.Fatalf("status = %s, want %s", got, wantStatus)
	}
}

func TestUndeletedTaskPathsCountsDanglingSymlinksAndDirectories(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "directory"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing", filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	task := spec.Task{Context: []spec.TaskContextRef{
		{Kind: spec.ContextKindDeletes, Path: "directory"},
		{Kind: spec.ContextKindDeletes, Path: "link"},
		{Kind: spec.ContextKindDeletes, Path: "gone"},
	}}
	remaining, err := undeletedTaskPaths(root, task)
	if err != nil || strings.Join(remaining, ",") != "directory,link" {
		t.Fatalf("remaining = %v, %v", remaining, err)
	}
}
