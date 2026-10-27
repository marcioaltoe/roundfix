package daemon

// Suite: authorized executable Task output
// Invariant: only a granted Task authorization keeps an untracked executable, including inside a new directory.
// Boundary IN: Task and Resolve cycles, Task and QA commits, real Git files and modes.
// Boundary OUT: Settle and archive.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func authorizeFixtureExecutables(t *testing.T, fixture *taskCycleFixture, paths ...string) {
	t.Helper()
	record := filepath.Join(fixture.gitRoot, "docs", "specs", taskCycleSlug, "_authorization.md")
	content, err := os.ReadFile(record)
	if err != nil {
		t.Fatal(err)
	}
	declaration := "executables:\n"
	for _, path := range paths {
		declaration += "  - " + path + "\n"
	}
	updated := strings.Replace(string(content), "---\n", "---\n"+declaration, 1)
	mustWriteForTest(t, record, updated)
	commitTaskFixtureSource(t, fixture.gitRoot, "authorize executable output")
}

func writeExecutableForTest(t *testing.T, root, path string) {
	t.Helper()
	absolute := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolute, []byte("#!/bin/sh\nexit 0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(absolute, 0o755); err != nil {
		t.Fatal(err)
	}
}

func executableRefusalForTest(path string, task bool) string {
	suffix := ""
	if task {
		suffix = " unless the Spec's _authorization.md lists the path under executables"
	}
	return fmt.Sprintf("roundfix: refused executable file %s (mode 0755); build artifacts and deliberately executable repository files are not valid Work Item output%s\n", path, suffix)
}

func TestTaskCommitKeepsAnAuthorizedExecutable(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	authorizeFixtureExecutables(t, fixture, "scripts/launch.sh")
	for _, path := range []string{"scripts/launch.sh", "scripts/other.sh"} {
		writeExecutableForTest(t, fixture.gitRoot, path)
	}
	fixture.worktree.snapshots = [][]string{{"scripts/"}}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)
	plan := fixture.plan()
	prepared, err := engine.prepareTaskCommit(context.Background(), plan, fixture.graph.Tasks[0], nil)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(prepared.stageable, "scripts/launch.sh") || slices.Contains(prepared.stageable, "scripts/other.sh") {
		t.Fatalf("stageable = %v", prepared.stageable)
	}
	if len(prepared.dropped) != 1 || prepared.dropped[0].Path != "scripts/other.sh" || !prepared.dropped[0].Lost {
		t.Fatalf("dropped = %+v", prepared.dropped)
	}
	if err := engine.commitTask(context.Background(), plan, fixture.graph.Tasks[0], 1, prepared); err != nil {
		t.Fatal(err)
	}
	if len(committer.paths) != 1 || !slices.Contains(committer.paths[0], "scripts/launch.sh") {
		t.Fatalf("committed paths = %v", committer.paths)
	}
	if !strings.Contains(fixture.progress.String(), executableRefusalForTest("scripts/other.sh", true)) {
		t.Fatalf("progress = %q", fixture.progress.String())
	}
	events := droppedStageEvents(t, fixture.sink)
	if len(events) != 1 || events[0].Summary != "Executable file scripts/other.sh refused with mode 0755; build artifacts and deliberately executable repository files are not valid Work Item output unless the Spec's _authorization.md lists the path under executables." {
		t.Fatalf("events = %+v", events)
	}
	payload := eventPayloadMap(t, events[0])
	if len(payload) != 5 || payload["decision"] != "dropped" || payload["task"] != "task_01" || payload["path"] != "scripts/other.sh" || payload["mode"] != "0755" || payload["reason"] != executableStagePathReason {
		t.Fatalf("payload = %+v", payload)
	}
}

func TestTaskCycleCommitsAnAuthorizedExecutableWithItsMode(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Add launcher"}})
	authorizeFixtureExecutables(t, fixture, "scripts/launch.sh")
	gittest.AppendConfig(t, fixture.gitRoot, "[core]\n\tfileMode = true\n")
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, writeByTask: map[string]string{"task_01": "scripts/launch.sh"}, afterTask: func(string) { writeExecutableForTest(t, fixture.gitRoot, "scripts/launch.sh") }}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(context.Background(), fixture.plan())
	if err != nil {
		t.Fatal(err)
	}
	if result.Completed != 1 || result.Failed != 0 {
		t.Fatalf("result = %+v", result)
	}
	tree := gittest.Run(t, fixture.gitRoot, "ls-tree", "HEAD", "--", "scripts/launch.sh")
	if !strings.HasPrefix(tree, "100755 blob ") || !strings.HasSuffix(tree, "\tscripts/launch.sh\n") {
		t.Fatalf("tree = %q", tree)
	}
}

func TestTaskCycleRefusesAnUnlistedExecutableInANewDirectory(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Add helper"}})
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, writeByTask: map[string]string{"task_01": "newdir/helper.sh"}, afterTask: func(string) { writeExecutableForTest(t, fixture.gitRoot, "newdir/helper.sh") }}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(context.Background(), fixture.plan())
	if err != nil {
		t.Fatal(err)
	}
	if result.Failed != 1 || result.Completed != 0 || len(result.Outcomes) != 1 || !strings.Contains(result.Outcomes[0].Reason, "newdir/helper.sh (executable file)") {
		t.Fatalf("result = %+v", result)
	}
	if got := gittest.Run(t, fixture.gitRoot, "ls-tree", "-r", "HEAD", "--", "newdir/"); got != "" {
		t.Fatalf("unexpected commit: %s", got)
	}
}

func TestResolveCycleRefusesAnExecutableInsideANewBatchDirectory(t *testing.T) {
	t.Parallel()
	fixture := newEngineFixture(t)
	gittest.InitRepo(t, fixture.gitRoot)
	gittest.Run(t, fixture.gitRoot, "add", ".")
	gittest.Run(t, fixture.gitRoot, "commit", "-m", "seed batch fixture")
	writeExecutableForTest(t, fixture.gitRoot, "newdir/helper.sh")
	mustWriteForTest(t, filepath.Join(fixture.gitRoot, "newdir", "fixed.go"), "package fixed\n")
	fixture.worktree.snapshots = [][]string{nil, {"newdir/"}}
	runner := &engineFakeRunner{calls: fixture.calls, store: fixture.store}
	verifier := &engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, &engineFakePusher{calls: fixture.calls}, &engineFakeSource{calls: fixture.calls})
	result, err := engine.ResolveCycle(context.Background(), fixture.plan())
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Batches) != 1 || !result.Batches[0].Committed {
		t.Fatalf("result = %+v", result)
	}
	if len(committer.paths) != 1 || !slices.Equal(committer.paths[0], []string{"newdir/fixed.go"}) {
		t.Fatalf("paths = %v", committer.paths)
	}
	if !strings.Contains(fixture.progress.String(), executableRefusalForTest("newdir/helper.sh", false)) || strings.Contains(fixture.progress.String(), "unless the Spec's") {
		t.Fatalf("progress = %q", fixture.progress.String())
	}
}

func TestQAReportCommitRefusesAListedExecutable(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	path := filepath.ToSlash(filepath.Join("docs", "specs", taskCycleSlug, "qa", "evidence", "helper.sh"))
	authorizeFixtureExecutables(t, fixture, path)
	writeExecutableForTest(t, fixture.gitRoot, path)
	fixture.worktree.snapshots = [][]string{{path}}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)
	err := engine.commitQAReport(context.Background(), fixture.plan(), 1, nil, nil, spec.VerdictPass, "", spec.Task{})
	if err != nil {
		t.Fatal(err)
	}
	if len(committer.paths) != 0 {
		t.Fatalf("QA staged executable: %v", committer.paths)
	}
	if fixture.progress.String() != executableRefusalForTest(path, false) {
		t.Fatalf("progress = %q", fixture.progress.String())
	}
	events := droppedStageEvents(t, fixture.sink)
	if len(events) != 1 || strings.Contains(events[0].Summary, "unless the Spec's") {
		t.Fatalf("events = %+v", events)
	}
}

func TestTaskCycleCommitsATrackedFileMadeExecutable(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Make script executable"}})
	gittest.AppendConfig(t, fixture.gitRoot, "[core]\n\tfileMode = true\n")
	mustWriteForTest(t, filepath.Join(fixture.gitRoot, "helper.sh"), "#!/bin/sh\nexit 0\n")
	commitTaskFixtureSource(t, fixture.gitRoot, "track script")
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, afterTask: func(string) {
		if err := os.Chmod(filepath.Join(fixture.gitRoot, "helper.sh"), 0o755); err != nil {
			t.Fatal(err)
		}
	}}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(context.Background(), fixture.plan())
	if err != nil {
		t.Fatal(err)
	}
	if result.Completed != 1 || result.Failed != 0 {
		t.Fatalf("result = %+v", result)
	}
	if tree := gittest.Run(t, fixture.gitRoot, "ls-tree", "HEAD", "--", "helper.sh"); !strings.HasPrefix(tree, "100755 blob ") || !strings.HasSuffix(tree, "\thelper.sh\n") {
		t.Fatalf("tree = %q", tree)
	}
}
