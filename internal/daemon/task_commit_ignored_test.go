package daemon

// Suite: Task commit staging
// Invariant: A Task commit stages the same expanded non-ignored files that its recorded-path disclosure names.
// Boundary IN: Daemon Task cycle, real Git index, commit, and Task file.
// Boundary OUT: Settle Command commits and QA scope auditing.

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestTaskCommitNeverStagesAnIgnoredFileInANewDirectory(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Add a package"}})
	ignoredPath := "internal/newpkg/ignored.secret"
	if err := os.WriteFile(filepath.Join(fixture.gitRoot, ".gitignore"), []byte(ignoredPath+"\n"), 0o644); err != nil {
		t.Fatalf("write fixture ignore rule: %v", err)
	}
	commitTaskFixtureSource(t, fixture.gitRoot, "ignore package secret")

	var writeErr error
	runner := &taskFakeRunner{
		calls:       fixture.calls,
		gitRoot:     fixture.gitRoot,
		writeByTask: map[string]string{"task_01": "internal/newpkg/new.go"},
		afterTask: func(taskID string) {
			if taskID != "task_01" {
				return
			}
			writeErr = os.WriteFile(filepath.Join(fixture.gitRoot, filepath.FromSlash(ignoredPath)), []byte("secret\n"), 0o600)
		},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	if writeErr != nil {
		t.Fatalf("write ignored package file: %v", writeErr)
	}
	if result.Completed != 1 || result.Failed != 0 {
		t.Fatalf("TaskCycle result = %+v, want one completed Task", result)
	}
	want := []string{"internal/newpkg/new.go"}
	committed := strings.Fields(gittest.Run(t, fixture.gitRoot, "ls-tree", "--name-only", "-r", "HEAD", "--", "internal/newpkg"))
	if !slices.Equal(committed, want) {
		t.Fatalf("Task commit package paths = %v, want %v", committed, want)
	}
	assertRecordedTaskPaths(t, fixture, "task_01", want)
}
