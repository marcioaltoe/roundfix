package worktree

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/store"
)

// Suite: Run Worktree git lifecycle.
// Invariant: Run Worktrees use named Run Branches and integrate through safe git porcelain.
// Boundary IN: internal/worktree behavior and real git temp repositories.
// Boundary OUT: daemon, store, and CLI wiring that will call this package in later tasks.

func TestCreateUsesNamedRunBranchUnderRoundfixHomeAndCopiesFiles(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	location := filepath.Join(homeDir, "configured-worktrees")
	repoDir := initWorktreeRepo(t)
	mustWriteWorktreeTest(t, filepath.Join(repoDir, ".gitignore"), ".env\n")
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "tracked.txt"), "base\n")
	gitWorktreeTest(t, repoDir, "add", ".gitignore", "tracked.txt")
	gitWorktreeTest(t, repoDir, "commit", "-m", "initial")
	headSHA := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "HEAD"))
	mustWriteWorktreeTest(t, filepath.Join(repoDir, ".env"), "SECRET=1\n")

	var ref Ref
	stderr := captureStderr(t, func() {
		var err error
		ref, err = Create(ctx, CreateOptions{
			UserRoot: repoDir,
			Location: location,
			RunID:    "task01",
			HeadSHA:  headSHA,
			CopyList: []string{".env", "missing.env"},
		})
		if err != nil {
			t.Fatalf("create Run Worktree: %v", err)
		}
	})

	if ref.RunID != "task01" {
		t.Fatalf("expected Run ID task01, got %q", ref.RunID)
	}
	if ref.UserRoot != filepath.Clean(repoDir) {
		t.Fatalf("expected user root %q, got %q", repoDir, ref.UserRoot)
	}
	if ref.Branch != "roundfix/run-task01" {
		t.Fatalf("expected named Run Branch, got %q", ref.Branch)
	}
	if filepath.Base(ref.Path) != "task01" {
		t.Fatalf("expected path to end in Run ID, got %q", ref.Path)
	}
	if filepath.Dir(filepath.Dir(ref.Path)) != location {
		t.Fatalf("expected path under configured worktree location %q, got %q", location, ref.Path)
	}
	repoSlug := filepath.Base(filepath.Dir(ref.Path))
	if !strings.HasPrefix(repoSlug, filepath.Base(repoDir)+"-") || !hasLowerHexSuffix(repoSlug, 8) {
		t.Fatalf("expected readable repo slug with 8 hex suffix, got %q", repoSlug)
	}
	if got := strings.TrimSpace(gitWorktreeTest(t, ref.Path, "branch", "--show-current")); got != ref.Branch {
		t.Fatalf("expected worktree branch %q, got %q", ref.Branch, got)
	}
	if got := strings.TrimSpace(gitWorktreeTest(t, repoDir, "branch", "--show-current")); got != "main" {
		t.Fatalf("expected source checkout to stay on main, got %q", got)
	}
	if got := mustReadWorktreeTest(t, filepath.Join(ref.Path, ".env")); got != "SECRET=1\n" {
		t.Fatalf("expected copied env file, got %q", got)
	}
	if !strings.Contains(stderr, "missing.env") || !strings.Contains(stderr, "missing") {
		t.Fatalf("expected missing copy-list note on stderr, got %q", stderr)
	}
	if status := gitStatus(t, ref.Path); status != "" {
		t.Fatalf("expected clean Run Worktree, got %q", status)
	}
}

func TestCreateRunsBootstrapAfterCopyInRunWorktreeRoot(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	location := filepath.Join(homeDir, "configured-worktrees")
	repoDir := initWorktreeRepo(t)
	mustWriteWorktreeTest(t, filepath.Join(repoDir, ".gitignore"), ".env\n")
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "tracked.txt"), "base\n")
	gitWorktreeTest(t, repoDir, "add", ".gitignore", "tracked.txt")
	gitWorktreeTest(t, repoDir, "commit", "-m", "initial")
	headSHA := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "HEAD"))
	mustWriteWorktreeTest(t, filepath.Join(repoDir, ".env"), "SECRET=1\n")
	var output bytes.Buffer

	ref, err := Create(ctx, CreateOptions{
		UserRoot: repoDir,
		Location: location,
		RunID:    "bootstrap-run",
		HeadSHA:  headSHA,
		CopyList: []string{".env"},
		Bootstrap: BootstrapSpec{
			Command: "test -f .env && pwd > bootstrap.pwd && cat .env > bootstrap.env && printf bootstrap-output",
			Timeout: time.Second,
		},
		BootstrapOutput: &output,
	})
	if err != nil {
		t.Fatalf("create Run Worktree: %v", err)
	}

	if got := strings.TrimSpace(mustReadWorktreeTest(t, filepath.Join(ref.Path, "bootstrap.pwd"))); got != ref.Path {
		t.Fatalf("expected bootstrap to run in %q, got %q", ref.Path, got)
	}
	if got := mustReadWorktreeTest(t, filepath.Join(ref.Path, "bootstrap.env")); got != "SECRET=1\n" {
		t.Fatalf("expected bootstrap to see copied .env, got %q", got)
	}
	if output.String() != "bootstrap-output" {
		t.Fatalf("expected bootstrap output to stream, got %q", output.String())
	}
}

func TestCreateTaskRunsBootstrapAfterCopyInTaskWorktreeRoot(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "task-bootstrap")
	mustWriteWorktreeTest(t, filepath.Join(fixture.repoDir, ".env.task"), "TASK_SECRET=1\n")
	var output bytes.Buffer

	task, err := CreateTaskWithOptions(ctx, fixture.ref, "task_01", TaskCreateOptions{
		CopyList: []string{".env.task"},
		Bootstrap: BootstrapSpec{
			Command: "test -f .env.task && pwd > task-bootstrap.pwd && cat .env.task > task-bootstrap.env && printf task-bootstrap-output",
			Timeout: time.Second,
		},
		BootstrapOutput: &output,
	})
	if err != nil {
		t.Fatalf("create Task Worktree: %v", err)
	}

	if got := strings.TrimSpace(mustReadWorktreeTest(t, filepath.Join(task.Path, "task-bootstrap.pwd"))); got != task.Path {
		t.Fatalf("expected bootstrap to run in %q, got %q", task.Path, got)
	}
	if got := mustReadWorktreeTest(t, filepath.Join(task.Path, "task-bootstrap.env")); got != "TASK_SECRET=1\n" {
		t.Fatalf("expected bootstrap to see copied .env.task, got %q", got)
	}
	if output.String() != "task-bootstrap-output" {
		t.Fatalf("expected bootstrap output to stream, got %q", output.String())
	}
}

func TestRunBootstrapReturnsBootstrapErrorOnNonZeroExit(t *testing.T) {
	var output bytes.Buffer
	command := "printf failure-tail; exit 7"

	err := runBootstrap(context.Background(), t.TempDir(), BootstrapSpec{Command: command, Timeout: time.Second}, &output)

	var bootstrapErr *BootstrapError
	if !errors.As(err, &bootstrapErr) {
		t.Fatalf("expected BootstrapError, got %T %[1]v", err)
	}
	if bootstrapErr.Command != command {
		t.Fatalf("expected command %q, got %q", command, bootstrapErr.Command)
	}
	if !strings.Contains(err.Error(), "worktree bootstrap failed: "+command+": exit status 7") {
		t.Fatalf("expected bootstrap failure message, got %q", err.Error())
	}
	if bootstrapErr.Tail != "failure-tail" {
		t.Fatalf("expected captured output tail, got %q", bootstrapErr.Tail)
	}
	if output.String() != "failure-tail" {
		t.Fatalf("expected output streamed before failure, got %q", output.String())
	}
}

func TestRunBootstrapReturnsBootstrapErrorOnTimeout(t *testing.T) {
	command := "sleep 1"

	err := runBootstrap(context.Background(), t.TempDir(), BootstrapSpec{Command: command, Timeout: 10 * time.Millisecond}, io.Discard)

	var bootstrapErr *BootstrapError
	if !errors.As(err, &bootstrapErr) {
		t.Fatalf("expected BootstrapError, got %T %[1]v", err)
	}
	if bootstrapErr.Command != command {
		t.Fatalf("expected command %q, got %q", command, bootstrapErr.Command)
	}
	if !strings.Contains(err.Error(), "worktree bootstrap failed: "+command+": timed out after 10ms") {
		t.Fatalf("expected timeout bootstrap failure, got %q", err.Error())
	}
}

func TestRunBootstrapSkipsEmptyCommand(t *testing.T) {
	var output bytes.Buffer

	err := runBootstrap(context.Background(), t.TempDir(), BootstrapSpec{Timeout: time.Nanosecond}, &output)

	if err != nil {
		t.Fatalf("expected empty bootstrap to skip, got %v", err)
	}
	if output.Len() != 0 {
		t.Fatalf("expected no output for empty bootstrap, got %q", output.String())
	}
}

func TestDeriveRootPathUsesReadableUniqueRepoSlug(t *testing.T) {
	location := t.TempDir()
	parentOne := t.TempDir()
	parentTwo := t.TempDir()
	repoOne := filepath.Join(parentOne, "same-name")
	repoTwo := filepath.Join(parentTwo, "same-name")
	mustMkdirWorktreeTest(t, repoOne)
	mustMkdirWorktreeTest(t, repoTwo)

	pathOne, err := deriveRootPath(location, repoOne, "run-one")
	if err != nil {
		t.Fatalf("derive first path: %v", err)
	}
	pathTwo, err := deriveRootPath(location, repoTwo, "run-one")
	if err != nil {
		t.Fatalf("derive second path: %v", err)
	}

	slugOne := filepath.Base(filepath.Dir(pathOne))
	slugTwo := filepath.Base(filepath.Dir(pathTwo))
	if !strings.HasPrefix(slugOne, "same-name-") || !hasLowerHexSuffix(slugOne, 8) {
		t.Fatalf("expected first readable repo slug with 8 hex suffix, got %q", slugOne)
	}
	if !strings.HasPrefix(slugTwo, "same-name-") || !hasLowerHexSuffix(slugTwo, 8) {
		t.Fatalf("expected second readable repo slug with 8 hex suffix, got %q", slugTwo)
	}
	if slugOne == slugTwo {
		t.Fatalf("expected same-named repos at different paths to get unique slugs, both got %q", slugOne)
	}
	if filepath.Dir(filepath.Dir(pathOne)) != location || filepath.Base(pathOne) != "run-one" {
		t.Fatalf("expected path under location/repo-slug/run-id, got %q", pathOne)
	}
}

func TestTaskWorktreesIntegrateFirstByFastForwardThenCherryPick(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "task-queue")
	base := strings.TrimSpace(gitWorktreeTest(t, fixture.ref.Path, "rev-parse", "HEAD"))
	mustWriteWorktreeTest(t, filepath.Join(fixture.repoDir, ".env.task"), "TASK_SECRET=1\n")
	first, err := CreateTask(ctx, fixture.ref, "task_01", []string{".env.task"})
	if err != nil {
		t.Fatalf("create first Task Worktree: %v", err)
	}
	second, err := CreateTask(ctx, fixture.ref, "task_02", nil)
	if err != nil {
		t.Fatalf("create second Task Worktree: %v", err)
	}
	if filepath.Dir(first.Path) != filepath.Dir(fixture.ref.Path) || filepath.Dir(second.Path) != filepath.Dir(fixture.ref.Path) {
		t.Fatalf("expected Task Worktrees to be siblings of Run Worktree, got run=%q first=%q second=%q", fixture.ref.Path, first.Path, second.Path)
	}
	if strings.HasPrefix(first.Path, fixture.ref.Path+string(filepath.Separator)) || strings.HasPrefix(second.Path, fixture.ref.Path+string(filepath.Separator)) {
		t.Fatalf("expected Task Worktrees not to nest under Run Worktree, got first=%q second=%q run=%q", first.Path, second.Path, fixture.ref.Path)
	}
	if got := mustReadWorktreeTest(t, filepath.Join(first.Path, ".env.task")); got != "TASK_SECRET=1\n" {
		t.Fatalf("expected copy-list file in first Task Worktree, got %q", got)
	}

	firstMessage := "task one subject\n\nTask-Trailer: first\n"
	firstSHA := commitTaskChange(t, first, "first.txt", "first\n", firstMessage)
	secondMessage := "task two subject\n\nTask-Trailer: second\n"
	secondSHA := commitTaskChange(t, second, "second.txt", "second\n", secondMessage)

	firstResult, err := IntegrateTask(ctx, fixture.ref, first)
	if err != nil {
		t.Fatalf("integrate first task: %v", err)
	}
	if firstResult.Mode != ModeTaskFastForward || firstResult.Reason != "" {
		t.Fatalf("expected first Task fast-forward result, got %#v", firstResult)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, fixture.ref.Path, "rev-parse", fixture.ref.Branch)); head != firstSHA {
		t.Fatalf("expected Run Branch at first task SHA %s, got %s", firstSHA, head)
	}

	secondResult, err := IntegrateTask(ctx, fixture.ref, second)
	if err != nil {
		t.Fatalf("integrate second task: %v", err)
	}
	if secondResult.Mode != ModeTaskCherryPick || secondResult.Reason != "" {
		t.Fatalf("expected second Task cherry-pick result, got %#v", secondResult)
	}
	runMessages := gitCommitMessages(t, fixture.ref.Path, base, fixture.ref.Branch)
	if len(runMessages) != 2 {
		t.Fatalf("expected two commits on Run Branch, got %d: %#v", len(runMessages), runMessages)
	}
	if runMessages[0] != firstMessage {
		t.Fatalf("expected first message preserved, got %q want %q", runMessages[0], firstMessage)
	}
	if runMessages[1] != secondMessage {
		t.Fatalf("expected second message preserved, got %q want %q", runMessages[1], secondMessage)
	}
	if taskMessage := gitCommitMessage(t, second.Path, secondSHA); taskMessage != secondMessage {
		t.Fatalf("expected source task message unchanged, got %q", taskMessage)
	}
}

func TestIntegrateTaskReturnsConflictAndLeavesRunBranchUnmoved(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "task-conflict")
	first, err := CreateTask(ctx, fixture.ref, "task_01", nil)
	if err != nil {
		t.Fatalf("create first Task Worktree: %v", err)
	}
	second, err := CreateTask(ctx, fixture.ref, "task_02", nil)
	if err != nil {
		t.Fatalf("create second Task Worktree: %v", err)
	}
	firstSHA := commitTaskChange(t, first, "shared.txt", "first task\n", "first shared\n")
	commitTaskChange(t, second, "shared.txt", "second task\n", "second shared\n")

	firstResult, err := IntegrateTask(ctx, fixture.ref, first)
	if err != nil {
		t.Fatalf("integrate first task: %v", err)
	}
	if firstResult.Mode != ModeTaskFastForward {
		t.Fatalf("expected first Task fast-forward, got %#v", firstResult)
	}
	before := strings.TrimSpace(gitWorktreeTest(t, fixture.ref.Path, "rev-parse", fixture.ref.Branch))
	if before != firstSHA {
		t.Fatalf("expected Run Branch at first SHA %s before conflict, got %s", firstSHA, before)
	}

	conflictResult, err := IntegrateTask(ctx, fixture.ref, second)
	if err != nil {
		t.Fatalf("integrate conflicting task: %v", err)
	}
	if conflictResult.Mode != ModeTaskConflict || !strings.Contains(conflictResult.Reason, "shared.txt") {
		t.Fatalf("expected conflict naming shared.txt, got %#v", conflictResult)
	}
	after := strings.TrimSpace(gitWorktreeTest(t, fixture.ref.Path, "rev-parse", fixture.ref.Branch))
	if after != before {
		t.Fatalf("expected Run Branch tip to stay %s after conflict, got %s", before, after)
	}
	if status := gitStatus(t, fixture.ref.Path); status != "" {
		t.Fatalf("expected Run Worktree clean after cherry-pick abort, got %q", status)
	}
	if _, err := os.Stat(second.Path); err != nil {
		t.Fatalf("expected conflicting Task Worktree kept, stat err=%v", err)
	}
	if !branchExists(t, fixture.repoDir, second.Branch) {
		t.Fatalf("expected conflicting Task Branch %q kept", second.Branch)
	}
}

func TestCleanupTaskRemovesTaskWorktreeAndBranch(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "task-cleanup")
	task, err := CreateTask(ctx, fixture.ref, "task_01", nil)
	if err != nil {
		t.Fatalf("create Task Worktree: %v", err)
	}
	commitTaskChange(t, task, "task.txt", "task\n", "task commit\n")
	result, err := IntegrateTask(ctx, fixture.ref, task)
	if err != nil {
		t.Fatalf("integrate task: %v", err)
	}
	if result.Mode != ModeTaskFastForward {
		t.Fatalf("expected Task fast-forward, got %#v", result)
	}

	if err := CleanupTask(ctx, task); err != nil {
		t.Fatalf("cleanup Task Worktree: %v", err)
	}

	if _, err := os.Stat(task.Path); !os.IsNotExist(err) {
		t.Fatalf("expected Task Worktree removed, stat err=%v", err)
	}
	if branchExists(t, fixture.repoDir, task.Branch) {
		t.Fatalf("expected Task Branch %q removed", task.Branch)
	}
}

func TestCleanupTaskRemovesTaskWorktreeWithUntrackedDebris(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "task-debris-cleanup")
	task, err := CreateTask(ctx, fixture.ref, "task_01", nil)
	if err != nil {
		t.Fatalf("create Task Worktree: %v", err)
	}
	mustWriteWorktreeTest(t, filepath.Join(task.Path, "task.txt"), "task\n")
	gitWorktreeTest(t, task.Path, "add", "task.txt")
	gitWorktreeTest(t, task.Path, "commit", "-m", "task commit")
	result, err := IntegrateTask(ctx, fixture.ref, task)
	if err != nil {
		t.Fatalf("integrate task: %v", err)
	}
	if result.Mode != ModeTaskFastForward {
		t.Fatalf("expected Task fast-forward, got %#v", result)
	}
	mustWriteWorktreeTest(t, filepath.Join(task.Path, ".env.local"), "secret=1\n")
	mustMkdirWorktreeTest(t, filepath.Join(task.Path, "node_modules", "cache"))
	mustWriteWorktreeTest(t, filepath.Join(task.Path, "node_modules", "cache", "entry.txt"), "cache\n")

	if err := CleanupTask(ctx, task); err != nil {
		t.Fatalf("cleanup Task Worktree with debris: %v", err)
	}

	if _, err := os.Stat(task.Path); !os.IsNotExist(err) {
		t.Fatalf("expected Task Worktree removed, stat err=%v", err)
	}
	if branchExists(t, fixture.repoDir, task.Branch) {
		t.Fatalf("expected Task Branch %q removed", task.Branch)
	}
}

func TestPruneTerminalReapsOnlyEmptyTerminalRunAndTaskBranches(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	repoDir := initWorktreeRepo(t)
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "tracked.txt"), "base\n")
	gitWorktreeTest(t, repoDir, "add", "tracked.txt")
	gitWorktreeTest(t, repoDir, "commit", "-m", "initial")
	headSHA := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "HEAD"))
	location := filepath.Join(homeDir, ".roundfix", "worktrees")

	emptyRun, err := Create(ctx, CreateOptions{UserRoot: repoDir, Location: location, RunID: "empty-run", HeadSHA: headSHA})
	if err != nil {
		t.Fatalf("create empty Run Worktree: %v", err)
	}
	emptyTask, err := CreateTask(ctx, emptyRun, "task_01", nil)
	if err != nil {
		t.Fatalf("create empty Task Worktree: %v", err)
	}
	mustWriteWorktreeTest(t, filepath.Join(emptyRun.Path, ".env.local"), "secret=1\n")
	mustMkdirWorktreeTest(t, filepath.Join(emptyRun.Path, "node_modules", "cache"))
	mustWriteWorktreeTest(t, filepath.Join(emptyRun.Path, "node_modules", "cache", "entry.txt"), "cache\n")
	mustWriteWorktreeTest(t, filepath.Join(emptyTask.Path, ".env.local"), "secret=1\n")
	mustMkdirWorktreeTest(t, filepath.Join(emptyTask.Path, "node_modules", "cache"))
	mustWriteWorktreeTest(t, filepath.Join(emptyTask.Path, "node_modules", "cache", "entry.txt"), "cache\n")
	valuableRun, err := Create(ctx, CreateOptions{UserRoot: repoDir, Location: location, RunID: "valuable-run", HeadSHA: headSHA})
	if err != nil {
		t.Fatalf("create valuable Run Worktree: %v", err)
	}
	fixture := integrationFixture{repoDir: repoDir, ref: valuableRun}
	fixture.commitRunChange(t, "valuable-run.txt", "valuable\n")
	valuableTask, err := CreateTask(ctx, valuableRun, "task_01", nil)
	if err != nil {
		t.Fatalf("create valuable Task Worktree: %v", err)
	}
	commitTaskChange(t, valuableTask, "valuable-task.txt", "valuable\n", "valuable task\n")
	nonTerminalRun, err := Create(ctx, CreateOptions{UserRoot: repoDir, Location: location, RunID: "nonterminal-run", HeadSHA: headSHA})
	if err != nil {
		t.Fatalf("create non-terminal Run Worktree: %v", err)
	}
	nonTerminalTask, err := CreateTask(ctx, nonTerminalRun, "task_01", nil)
	if err != nil {
		t.Fatalf("create non-terminal Task Worktree: %v", err)
	}

	err = PruneTerminal(ctx, repoDir, location, func(runID string) bool {
		return runID == "empty-run" || runID == "valuable-run"
	})
	if err != nil {
		t.Fatalf("prune terminal: %v", err)
	}

	assertPathRemoved(t, emptyRun.Path)
	assertPathRemoved(t, emptyTask.Path)
	assertBranchRemoved(t, repoDir, emptyRun.Branch)
	assertBranchRemoved(t, repoDir, emptyTask.Branch)
	assertPathExists(t, valuableRun.Path)
	assertPathExists(t, valuableTask.Path)
	assertRunBranchExists(t, repoDir, valuableRun.Branch)
	assertRunBranchExists(t, repoDir, valuableTask.Branch)
	assertPathExists(t, nonTerminalRun.Path)
	assertPathExists(t, nonTerminalTask.Path)
	assertRunBranchExists(t, repoDir, nonTerminalRun.Branch)
	assertRunBranchExists(t, repoDir, nonTerminalTask.Branch)
}

func TestListPendingRunWorkReportsAheadRunBranches(t *testing.T) {
	ctx := context.Background()
	repoDir := initWorktreeRepo(t)
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "base.txt"), "base\n")
	gitWorktreeTest(t, repoDir, "add", "base.txt")
	gitWorktreeTest(t, repoDir, "commit", "-m", "initial")

	pending, err := ListPendingRunWork(ctx, repoDir, "main")
	if err != nil {
		t.Fatalf("list pending Run Branch work with no Run Branches: %v", err)
	}
	if len(pending) != 0 {
		t.Fatalf("expected no pending Run Branch work, got %#v", pending)
	}

	zeroBranch := BranchName("zero")
	gitWorktreeTest(t, repoDir, "branch", zeroBranch, "main")
	divergedBranch := BranchName("diverged")
	gitWorktreeTest(t, repoDir, "checkout", "-b", divergedBranch, "main")
	commitWorktreeFile(t, repoDir, "diverged.txt", "diverged\n", "diverged change")
	gitWorktreeTest(t, repoDir, "checkout", "main")
	commitWorktreeFile(t, repoDir, "main.txt", "main\n", "main change")
	fastForwardBranch := BranchName("ff")
	fastForwardPath := filepath.Join(t.TempDir(), "ff-worktree")
	gitWorktreeTest(t, repoDir, "worktree", "add", "-b", fastForwardBranch, fastForwardPath, "main")
	commitWorktreeFile(t, fastForwardPath, "ff.txt", "ff\n", "ff change")

	pending, err = ListPendingRunWork(ctx, repoDir, "main")
	if err != nil {
		t.Fatalf("list pending Run Branch work: %v", err)
	}

	if len(pending) != 2 {
		t.Fatalf("expected two pending Run Branches, got %#v", pending)
	}
	want := map[string]PendingRunWork{
		divergedBranch: {
			Branch:       divergedBranch,
			WorktreePath: "",
			AheadCommits: 1,
			FastForward:  false,
		},
		fastForwardBranch: {
			Branch:       fastForwardBranch,
			WorktreePath: canonicalPath(fastForwardPath),
			AheadCommits: 1,
			FastForward:  true,
		},
	}
	for _, got := range pending {
		if got.Branch == zeroBranch {
			t.Fatalf("expected zero-ahead Run Branch %q to be omitted, got %#v", zeroBranch, pending)
		}
		expected, ok := want[got.Branch]
		if !ok {
			t.Fatalf("unexpected pending Run Branch %#v", got)
		}
		if got != expected {
			t.Fatalf("pending Run Branch %q = %#v, want %#v", got.Branch, got, expected)
		}
		delete(want, got.Branch)
	}
	if len(want) != 0 {
		t.Fatalf("missing pending Run Branches: %#v", want)
	}
}

func TestIntegratePendingRunWorkFastForwardsBaseBranch(t *testing.T) {
	ctx := context.Background()
	repoDir := initWorktreeRepo(t)
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "base.txt"), "base\n")
	gitWorktreeTest(t, repoDir, "add", "base.txt")
	gitWorktreeTest(t, repoDir, "commit", "-m", "initial")
	runBranch := BranchName("ready")
	gitWorktreeTest(t, repoDir, "checkout", "-b", runBranch, "main")
	runSHA := commitWorktreeFile(t, repoDir, "run.txt", "run\n", "run change")
	gitWorktreeTest(t, repoDir, "checkout", "main")

	if err := IntegratePendingRunWork(ctx, repoDir, "main", runBranch); err != nil {
		t.Fatalf("integrate pending Run Branch work: %v", err)
	}

	if head := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "main")); head != runSHA {
		t.Fatalf("expected main at Run Branch SHA %s, got %s", runSHA, head)
	}
	if status := gitStatus(t, repoDir); status != "" {
		t.Fatalf("expected clean checkout after pending Run Branch integration, got %q", status)
	}
}

func TestIntegratePendingRunWorkRefusesDivergedBranch(t *testing.T) {
	ctx := context.Background()
	repoDir := initWorktreeRepo(t)
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "base.txt"), "base\n")
	gitWorktreeTest(t, repoDir, "add", "base.txt")
	gitWorktreeTest(t, repoDir, "commit", "-m", "initial")
	runBranch := BranchName("diverged")
	gitWorktreeTest(t, repoDir, "checkout", "-b", runBranch, "main")
	runSHA := commitWorktreeFile(t, repoDir, "run.txt", "run\n", "run change")
	gitWorktreeTest(t, repoDir, "checkout", "main")
	userSHA := commitWorktreeFile(t, repoDir, "user.txt", "user\n", "user change")

	err := IntegratePendingRunWork(ctx, repoDir, "main", runBranch)

	if err == nil {
		t.Fatal("expected divergence refusal, got nil")
	}
	if !strings.Contains(err.Error(), "fast-forward is impossible") || !strings.Contains(err.Error(), runBranch) {
		t.Fatalf("expected descriptive fast-forward refusal naming branch, got %v", err)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "main")); head != userSHA {
		t.Fatalf("expected main to stay at user SHA %s, got %s", userSHA, head)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", runBranch)); head != runSHA {
		t.Fatalf("expected Run Branch to stay at SHA %s, got %s", runSHA, head)
	}
	if status := gitStatus(t, repoDir); status != "" {
		t.Fatalf("expected clean checkout after divergence refusal, got %q", status)
	}
}

func TestIntegrateFastForwardsCleanCheckout(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "ff-clean")
	runSHA := fixture.commitRunChange(t, "run.txt", "run change\n")

	result, err := Integrate(ctx, fixture.ref, "main", runSHA)
	if err != nil {
		t.Fatalf("integrate: %v", err)
	}

	if result.Mode != ModeFastForwardMerge || result.Reason != "" {
		t.Fatalf("expected ff merge result, got %#v", result)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main")); head != runSHA {
		t.Fatalf("expected main at run SHA %s, got %s", runSHA, head)
	}
	if status := gitStatus(t, fixture.repoDir); status != "" {
		t.Fatalf("expected clean user checkout, got %q", status)
	}
}

func TestIntegrateFastForwardsPreservingNonOverlappingDirt(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "ff-dirty")
	runSHA := fixture.commitRunChange(t, "run.txt", "run change\n")
	mustWriteWorktreeTest(t, filepath.Join(fixture.repoDir, "user.txt"), "user dirt\n")

	result, err := Integrate(ctx, fixture.ref, "main", runSHA)
	if err != nil {
		t.Fatalf("integrate: %v", err)
	}

	if result.Mode != ModeFastForwardMerge || result.Reason != "" {
		t.Fatalf("expected ff merge result, got %#v", result)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main")); head != runSHA {
		t.Fatalf("expected main at run SHA %s, got %s", runSHA, head)
	}
	status := gitStatus(t, fixture.repoDir)
	if !strings.Contains(status, " M user.txt") {
		t.Fatalf("expected non-overlapping user dirt preserved, got %q", status)
	}
}

func TestIntegrateReturnsPendingOnOverlapAndLeavesBranchAndDirt(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "overlap")
	before := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main"))
	runSHA := fixture.commitRunChange(t, "shared.txt", "run change\n")
	mustWriteWorktreeTest(t, filepath.Join(fixture.repoDir, "shared.txt"), "user dirt\n")

	result, err := Integrate(ctx, fixture.ref, "main", runSHA)
	if err != nil {
		t.Fatalf("integrate: %v", err)
	}

	if result.Mode != ModePending || result.Reason != ReasonOverlap {
		t.Fatalf("expected pending overlap, got %#v", result)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main")); head != before {
		t.Fatalf("expected main to stay at %s, got %s", before, head)
	}
	status := gitStatus(t, fixture.repoDir)
	if !strings.Contains(status, " M shared.txt") {
		t.Fatalf("expected overlapping user dirt preserved, got %q", status)
	}
	assertNoPhantomStagedEntries(t, status)
}

func TestIntegrateReturnsPendingOnDivergenceAndLeavesBranch(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "diverged")
	runSHA := fixture.commitRunChange(t, "run.txt", "run change\n")
	mustWriteWorktreeTest(t, filepath.Join(fixture.repoDir, "user.txt"), "user commit\n")
	gitWorktreeTest(t, fixture.repoDir, "add", "user.txt")
	gitWorktreeTest(t, fixture.repoDir, "commit", "-m", "user commit")
	userSHA := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main"))

	result, err := Integrate(ctx, fixture.ref, "main", runSHA)
	if err != nil {
		t.Fatalf("integrate: %v", err)
	}

	if result.Mode != ModePending || result.Reason != ReasonDiverged {
		t.Fatalf("expected pending divergence, got %#v", result)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main")); head != userSHA {
		t.Fatalf("expected main to stay at %s, got %s", userSHA, head)
	}
	status := gitStatus(t, fixture.repoDir)
	if status != "" {
		t.Fatalf("expected clean user checkout after divergence refusal, got %q", status)
	}
	assertNoPhantomStagedEntries(t, status)
}

func TestIntegrateMovesBranchWhenTargetCheckedOutNowhere(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "branch-move")
	runSHA := fixture.commitRunChange(t, "run.txt", "run change\n")
	gitWorktreeTest(t, fixture.repoDir, "checkout", "-b", "other")

	result, err := Integrate(ctx, fixture.ref, "main", runSHA)
	if err != nil {
		t.Fatalf("integrate: %v", err)
	}

	if result.Mode != ModeBranchMove || result.Reason != "" {
		t.Fatalf("expected branch move result, got %#v", result)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main")); head != runSHA {
		t.Fatalf("expected main at run SHA %s, got %s", runSHA, head)
	}
	if branch := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "branch", "--show-current")); branch != "other" {
		t.Fatalf("expected user checkout to remain on other branch, got %q", branch)
	}
	if status := gitStatus(t, fixture.repoDir); status != "" {
		t.Fatalf("expected clean user checkout, got %q", status)
	}
}

func TestIntegrateReturnsPendingOnNonAncestryWhenTargetCheckedOutNowhere(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "non-ancestry")
	runSHA := fixture.commitRunChange(t, "run.txt", "run change\n")
	mustWriteWorktreeTest(t, filepath.Join(fixture.repoDir, "user.txt"), "user commit\n")
	gitWorktreeTest(t, fixture.repoDir, "add", "user.txt")
	gitWorktreeTest(t, fixture.repoDir, "commit", "-m", "user commit")
	userSHA := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main"))
	gitWorktreeTest(t, fixture.repoDir, "checkout", "-b", "other")

	result, err := Integrate(ctx, fixture.ref, "main", runSHA)
	if err != nil {
		t.Fatalf("integrate: %v", err)
	}

	if result.Mode != ModePending || result.Reason != ReasonNonAncestry {
		t.Fatalf("expected pending non-ancestry, got %#v", result)
	}
	if head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main")); head != userSHA {
		t.Fatalf("expected main to stay at %s, got %s", userSHA, head)
	}
	status := gitStatus(t, fixture.repoDir)
	if status != "" {
		t.Fatalf("expected clean user checkout after non-ancestry refusal, got %q", status)
	}
	assertNoPhantomStagedEntries(t, status)
}

func TestCleanupCleanRemovesCleanRunWorktreeAndBranch(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "cleanup")

	if err := CleanupClean(ctx, fixture.ref); err != nil {
		t.Fatalf("cleanup Clean: %v", err)
	}

	if _, err := os.Stat(fixture.ref.Path); !os.IsNotExist(err) {
		t.Fatalf("expected Run Worktree removed, stat err=%v", err)
	}
	if branchExists(t, fixture.repoDir, fixture.ref.Branch) {
		t.Fatalf("expected Run Branch %q deleted", fixture.ref.Branch)
	}
}

func TestCleanupCleanRemovesRunWorktreeWithUntrackedDebris(t *testing.T) {
	ctx := context.Background()
	fixture := newIntegrationFixture(t, "debris-cleanup")
	mustWriteWorktreeTest(t, filepath.Join(fixture.ref.Path, ".env.local"), "secret=1\n")
	mustMkdirWorktreeTest(t, filepath.Join(fixture.ref.Path, "node_modules", "cache"))
	mustWriteWorktreeTest(t, filepath.Join(fixture.ref.Path, "node_modules", "cache", "entry.txt"), "cache\n")

	if err := CleanupClean(ctx, fixture.ref); err != nil {
		t.Fatalf("cleanup Clean with debris: %v", err)
	}

	if _, err := os.Stat(fixture.ref.Path); !os.IsNotExist(err) {
		t.Fatalf("expected Run Worktree removed, stat err=%v", err)
	}
	if branchExists(t, fixture.repoDir, fixture.ref.Branch) {
		t.Fatalf("expected Run Branch %q deleted", fixture.ref.Branch)
	}
}

func TestPruneTerminalReapsCrashedTerminalCleanRunOnly(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	repoDir := initWorktreeRepo(t)
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "tracked.txt"), "base\n")
	gitWorktreeTest(t, repoDir, "add", "tracked.txt")
	gitWorktreeTest(t, repoDir, "commit", "-m", "initial")
	headSHA := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "HEAD"))
	location := filepath.Join(homeDir, ".roundfix", "worktrees")
	cleanRef, err := Create(ctx, CreateOptions{UserRoot: repoDir, Location: location, RunID: "terminal-clean", HeadSHA: headSHA})
	if err != nil {
		t.Fatalf("create terminal clean ref: %v", err)
	}
	keptRef, err := Create(ctx, CreateOptions{UserRoot: repoDir, Location: location, RunID: "kept", HeadSHA: headSHA})
	if err != nil {
		t.Fatalf("create kept ref: %v", err)
	}
	if err := os.RemoveAll(cleanRef.Path); err != nil {
		t.Fatalf("simulate crashed worktree removal: %v", err)
	}

	err = PruneTerminal(ctx, repoDir, location, func(runID string) bool {
		return runID == "terminal-clean"
	})
	if err != nil {
		t.Fatalf("prune terminal: %v", err)
	}

	if branchExists(t, repoDir, cleanRef.Branch) {
		t.Fatalf("expected terminal Clean Run Branch %q deleted", cleanRef.Branch)
	}
	if _, err := os.Stat(keptRef.Path); err != nil {
		t.Fatalf("expected non-terminal Run Worktree to survive, stat err=%v", err)
	}
	if !branchExists(t, repoDir, keptRef.Branch) {
		t.Fatalf("expected non-terminal Run Branch %q to survive", keptRef.Branch)
	}
}

func TestInspectTerminalRunSafe(t *testing.T) {
	ctx := context.Background()
	fixture := newTerminalRunFixture(t, "reconcile-safe")
	runHead := fixture.commitRunChange(t, "safe.txt", "safe\n")
	gitWorktreeTest(t, fixture.repoDir, "merge", "--ff-only", fixture.ref.Branch)

	result, err := InspectTerminalRun(ctx, fixture.run)
	if err != nil {
		t.Fatalf("inspect terminal Run: %v", err)
	}

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationSafe)
	if result.RunHead != runHead || result.TargetHead != runHead {
		t.Fatalf("expected both heads at %s, got Run=%q target=%q", runHead, result.RunHead, result.TargetHead)
	}
}

func TestInspectTerminalRunConcurrentSafe(t *testing.T) {
	ctx := context.Background()
	fixture := newTerminalRunFixture(t, "reconcile-concurrent-safe")
	runHead := fixture.commitRunChange(t, "safe.txt", "safe\n")
	gitWorktreeTest(t, fixture.repoDir, "merge", "--ff-only", fixture.ref.Branch)

	type inspection struct {
		result RunWorktreeReconciliation
		err    error
	}
	const inspectors = 8
	inspections := make(chan inspection, inspectors)
	for range inspectors {
		go func() {
			result, err := InspectTerminalRun(ctx, fixture.run)
			inspections <- inspection{result: result, err: err}
		}()
	}
	for range inspectors {
		got := <-inspections
		if got.err != nil {
			t.Fatalf("inspect terminal Run concurrently: %v", got.err)
		}
		assertTerminalRunReconciliation(t, got.result, fixture.run, ReconciliationSafe)
		if got.result.RunHead != runHead || got.result.TargetHead != runHead {
			t.Fatalf("expected both concurrent heads at %s, got Run=%q target=%q", runHead, got.result.RunHead, got.result.TargetHead)
		}
	}
}

func TestInspectTerminalRunUnintegrated(t *testing.T) {
	ctx := context.Background()
	fixture := newTerminalRunFixture(t, "reconcile-unintegrated")
	runHead := fixture.commitRunChange(t, "unintegrated.txt", "unintegrated\n")
	targetHead := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main"))

	result, err := InspectTerminalRun(ctx, fixture.run)
	if err != nil {
		t.Fatalf("inspect terminal Run: %v", err)
	}

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnintegrated)
	if result.RunHead != runHead || result.TargetHead != targetHead {
		t.Fatalf("expected Run head %s and target head %s, got Run=%q target=%q", runHead, targetHead, result.RunHead, result.TargetHead)
	}
}

func TestInspectTerminalRunDirty(t *testing.T) {
	tests := []struct {
		name          string
		missingTarget bool
		dirty         func(t *testing.T, fixture terminalRunFixture)
	}{
		{
			name:          "tracked change with missing target metadata",
			missingTarget: true,
			dirty: func(t *testing.T, fixture terminalRunFixture) {
				mustWriteWorktreeTest(t, filepath.Join(fixture.ref.Path, "shared.txt"), "tracked dirt\n")
			},
		},
		{
			name: "untracked change",
			dirty: func(t *testing.T, fixture terminalRunFixture) {
				mustWriteWorktreeTest(t, filepath.Join(fixture.ref.Path, "untracked.txt"), "untracked dirt\n")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			fixture := newTerminalRunFixture(t, "reconcile-dirty-"+strings.ReplaceAll(test.name, " ", "-"))
			test.dirty(t, fixture)
			if test.missingTarget {
				fixture.run.LocalBranch = ""
			}
			head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main"))

			result, err := InspectTerminalRun(ctx, fixture.run)
			if err != nil {
				t.Fatalf("inspect terminal Run: %v", err)
			}

			assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationDirty)
			if result.RunHead != head {
				t.Fatalf("expected known Run head at %s for dirty worktree, got %q", head, result.RunHead)
			}
			if test.missingTarget && result.TargetHead != "" {
				t.Fatalf("expected missing target head to stay empty, got %q", result.TargetHead)
			}
			if !test.missingTarget && result.TargetHead != head {
				t.Fatalf("expected known target head at %s for dirty worktree, got %q", head, result.TargetHead)
			}
		})
	}
}

func TestInspectTerminalRunUnknownMissingTarget(t *testing.T) {
	ctx := context.Background()
	fixture := newTerminalRunFixture(t, "reconcile-unknown-target")
	fixture.run.LocalBranch = ""
	runHead := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", fixture.ref.Branch))

	result, err := InspectTerminalRun(ctx, fixture.run)
	if err != nil {
		t.Fatalf("inspect terminal Run: %v", err)
	}

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnknown)
	if result.RunHead != runHead || result.TargetHead != "" {
		t.Fatalf("expected only Run head %s to resolve, got Run=%q target=%q", runHead, result.RunHead, result.TargetHead)
	}
}

func TestInspectTerminalRunUnknownAmbiguousRef(t *testing.T) {
	ctx := context.Background()
	fixture := newTerminalRunFixture(t, "reconcile-unknown-ambiguous")
	fixture.commitRunChange(t, "ambiguous.txt", "ambiguous\n")
	gitWorktreeTest(t, fixture.repoDir, "branch", "ambiguous-target", "main")
	gitWorktreeTest(t, fixture.repoDir, "tag", "ambiguous-target", fixture.ref.Branch)
	fixture.run.LocalBranch = "ambiguous-target"

	result, err := InspectTerminalRun(ctx, fixture.run)
	if err != nil {
		t.Fatalf("inspect terminal Run: %v", err)
	}

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnknown)
	if result.RunHead == "" {
		t.Fatal("expected Run head to remain available when target ref is ambiguous")
	}
	if result.TargetHead != "" {
		t.Fatalf("expected ambiguous target head to stay empty, got %q", result.TargetHead)
	}
}

func TestInspectTerminalRunUnknownMissingRunBranch(t *testing.T) {
	ctx := context.Background()
	fixture := newTerminalRunFixture(t, "reconcile-unknown-run-branch")
	gitWorktreeTest(t, fixture.ref.Path, "checkout", "--detach")
	gitWorktreeTest(t, fixture.repoDir, "branch", "-D", fixture.ref.Branch)

	result, err := InspectTerminalRun(ctx, fixture.run)
	if err != nil {
		t.Fatalf("inspect terminal Run: %v", err)
	}

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnknown)
	if result.RunHead != "" {
		t.Fatalf("expected missing Run Branch head to stay empty, got %q", result.RunHead)
	}
}

func TestInspectTerminalRunReleased(t *testing.T) {
	ctx := context.Background()
	fixture := newTerminalRunFixture(t, "reconcile-released")
	gitWorktreeTest(t, fixture.repoDir, "worktree", "remove", "--force", fixture.ref.Path)
	gitWorktreeTest(t, fixture.repoDir, "branch", "-D", fixture.ref.Branch)

	result, err := InspectTerminalRun(ctx, fixture.run)
	if err != nil {
		t.Fatalf("inspect terminal Run: %v", err)
	}

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationReleased)
	if result.RunHead != "" || result.TargetHead != "" {
		t.Fatalf("expected released result without heads, got Run=%q target=%q", result.RunHead, result.TargetHead)
	}
}

func TestInspectTerminalRunSafeWithoutWorktree(t *testing.T) {
	ctx := context.Background()
	fixture := newTerminalRunFixture(t, "reconcile-safe-without-worktree")
	gitWorktreeTest(t, fixture.repoDir, "worktree", "remove", "--force", fixture.ref.Path)
	head := strings.TrimSpace(gitWorktreeTest(t, fixture.repoDir, "rev-parse", "main"))

	result, err := InspectTerminalRun(ctx, fixture.run)
	if err != nil {
		t.Fatalf("inspect terminal Run: %v", err)
	}

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationSafe)
	if result.RunHead != head || result.TargetHead != head {
		t.Fatalf("expected branch-only heads at %s, got Run=%q target=%q", head, result.RunHead, result.TargetHead)
	}
}

func TestInspectTerminalRunUnsafePath(t *testing.T) {
	t.Run("symlinked worktree", func(t *testing.T) {
		fixture := newTerminalRunFixture(t, "reconcile-unsafe-symlink")
		link := filepath.Join(t.TempDir(), "linked-worktree")
		if err := os.Symlink(fixture.ref.Path, link); err != nil {
			t.Fatalf("create worktree symlink: %v", err)
		}
		fixture.run.WorkDir = link

		result, err := InspectTerminalRun(context.Background(), fixture.run)
		if err != nil {
			t.Fatalf("inspect terminal Run: %v", err)
		}

		assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnknown)
	})

	t.Run("symlinked worktree parent", func(t *testing.T) {
		fixture := newTerminalRunFixture(t, "reconcile-unsafe-parent")
		link := filepath.Join(t.TempDir(), "linked-parent")
		if err := os.Symlink(filepath.Dir(fixture.ref.Path), link); err != nil {
			t.Fatalf("create worktree parent symlink: %v", err)
		}
		fixture.run.WorkDir = filepath.Join(link, filepath.Base(fixture.ref.Path))

		result, err := InspectTerminalRun(context.Background(), fixture.run)
		if err != nil {
			t.Fatalf("inspect terminal Run: %v", err)
		}

		assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnknown)
	})

	t.Run("unclean worktree path", func(t *testing.T) {
		fixture := newTerminalRunFixture(t, "reconcile-unsafe-path")
		fixture.run.WorkDir = fixture.repoDir + string(filepath.Separator) + ".." + string(filepath.Separator) + "outside"

		result, err := InspectTerminalRun(context.Background(), fixture.run)
		if err != nil {
			t.Fatalf("inspect terminal Run: %v", err)
		}

		assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnknown)
	})

	t.Run("symlinked Git root", func(t *testing.T) {
		fixture := newTerminalRunFixture(t, "reconcile-unsafe-root")
		link := filepath.Join(t.TempDir(), "linked-root")
		if err := os.Symlink(fixture.repoDir, link); err != nil {
			t.Fatalf("create Git root symlink: %v", err)
		}
		fixture.run.GitRoot = link

		if _, err := InspectTerminalRun(context.Background(), fixture.run); err == nil {
			t.Fatal("expected symlinked recorded Git root to be rejected")
		}
	})

	t.Run("symlinked Git root parent", func(t *testing.T) {
		fixture := newTerminalRunFixture(t, "reconcile-unsafe-root-parent")
		link := filepath.Join(t.TempDir(), "linked-root-parent")
		if err := os.Symlink(filepath.Dir(fixture.repoDir), link); err != nil {
			t.Fatalf("create Git root parent symlink: %v", err)
		}
		fixture.run.GitRoot = filepath.Join(link, filepath.Base(fixture.repoDir))

		if _, err := InspectTerminalRun(context.Background(), fixture.run); err == nil {
			t.Fatal("expected recorded Git root under a symlink to be rejected")
		}
	})
}

type terminalRunFixture struct {
	integrationFixture
	run store.Run
}

func newTerminalRunFixture(t *testing.T, runID string) terminalRunFixture {
	t.Helper()
	fixture := newIntegrationFixture(t, runID)
	return terminalRunFixture{
		integrationFixture: fixture,
		run: store.Run{
			ID:          runID,
			Kind:        store.KindImplement,
			State:       store.StateIntegrationPending,
			GitRoot:     canonicalPath(fixture.repoDir),
			LocalBranch: "main",
			WorkDir:     canonicalPath(fixture.ref.Path),
			SpecSlug:    "0038-terminal-run-worktree-reconciliation",
		},
	}
}

func assertTerminalRunReconciliation(t *testing.T, result RunWorktreeReconciliation, run store.Run, want ReconciliationState) {
	t.Helper()
	if result.State != want {
		t.Fatalf("expected reconciliation state %q, got %#v", want, result)
	}
	if result.RunID != run.ID || result.Outcome != run.State {
		t.Fatalf("expected recorded Run identity and outcome, got %#v", result)
	}
	if result.Path != run.WorkDir || result.Branch != BranchName(run.ID) {
		t.Fatalf("expected recorded worktree %q and Run Branch %q, got %#v", run.WorkDir, BranchName(run.ID), result)
	}
	if result.TargetBranch != run.LocalBranch {
		t.Fatalf("expected recorded target branch %q, got %#v", run.LocalBranch, result)
	}
	if result.Reason == "" || len(result.Reason) > 160 || strings.ContainsAny(result.Reason, "\r\n") {
		t.Fatalf("expected one bounded deterministic reason, got %q", result.Reason)
	}
	switch result.State {
	case ReconciliationSafe, ReconciliationUnintegrated, ReconciliationDirty, ReconciliationUnknown, ReconciliationReleased:
	default:
		t.Fatalf("unexpected reconciliation state %q", result.State)
	}
}

type integrationFixture struct {
	repoDir string
	ref     Ref
}

func newIntegrationFixture(t *testing.T, runID string) integrationFixture {
	t.Helper()
	ctx := context.Background()
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	repoDir := initWorktreeRepo(t)
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "shared.txt"), "base\n")
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "run.txt"), "base\n")
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "user.txt"), "base\n")
	gitWorktreeTest(t, repoDir, "add", "shared.txt", "run.txt", "user.txt")
	gitWorktreeTest(t, repoDir, "commit", "-m", "initial")
	headSHA := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "HEAD"))
	location := filepath.Join(homeDir, ".roundfix", "worktrees")
	ref, err := Create(ctx, CreateOptions{UserRoot: repoDir, Location: location, RunID: runID, HeadSHA: headSHA})
	if err != nil {
		t.Fatalf("create Run Worktree: %v", err)
	}
	return integrationFixture{repoDir: repoDir, ref: ref}
}

func (fixture integrationFixture) commitRunChange(t *testing.T, path string, content string) string {
	t.Helper()
	mustWriteWorktreeTest(t, filepath.Join(fixture.ref.Path, path), content)
	gitWorktreeTest(t, fixture.ref.Path, "add", path)
	gitWorktreeTest(t, fixture.ref.Path, "commit", "-m", "run change")
	return strings.TrimSpace(gitWorktreeTest(t, fixture.ref.Path, "rev-parse", "HEAD"))
}

func commitTaskChange(t *testing.T, task TaskRef, path string, content string, message string) string {
	t.Helper()
	mustWriteWorktreeTest(t, filepath.Join(task.Path, path), content)
	gitWorktreeTest(t, task.Path, "add", path)
	messagePath := filepath.Join(t.TempDir(), "message.txt")
	mustWriteWorktreeTest(t, messagePath, message)
	gitWorktreeTest(t, task.Path, "commit", "--cleanup=verbatim", "-F", messagePath)
	return strings.TrimSpace(gitWorktreeTest(t, task.Path, "rev-parse", "HEAD"))
}

func commitWorktreeFile(t *testing.T, workDir string, path string, content string, message string) string {
	t.Helper()
	fullPath := filepath.Join(workDir, path)
	mustMkdirWorktreeTest(t, filepath.Dir(fullPath))
	mustWriteWorktreeTest(t, fullPath, content)
	gitWorktreeTest(t, workDir, "add", path)
	gitWorktreeTest(t, workDir, "commit", "-m", message)
	return strings.TrimSpace(gitWorktreeTest(t, workDir, "rev-parse", "HEAD"))
}

func gitCommitMessages(t *testing.T, workDir string, base string, branch string) []string {
	t.Helper()
	output := gitWorktreeTest(t, workDir, "log", "--reverse", "--format=%H", base+".."+branch)
	var messages []string
	for _, line := range strings.Split(output, "\n") {
		sha := strings.TrimSpace(line)
		if sha == "" {
			continue
		}
		messages = append(messages, gitCommitMessage(t, workDir, sha))
	}
	return messages
}

func gitCommitMessage(t *testing.T, workDir string, sha string) string {
	t.Helper()
	raw := gitWorktreeTest(t, workDir, "cat-file", "commit", sha)
	_, message, ok := strings.Cut(raw, "\n\n")
	if !ok {
		t.Fatalf("commit %s did not contain a raw message body", sha)
	}
	return message
}

func initWorktreeRepo(t *testing.T) string {
	t.Helper()
	repoDir := t.TempDir()
	gitWorktreeTest(t, repoDir, "init", "-b", "main")
	gitWorktreeTest(t, repoDir, "config", "user.name", "Roundfix Test")
	gitWorktreeTest(t, repoDir, "config", "user.email", "test@example.com")
	gitWorktreeTest(t, repoDir, "config", "commit.gpgsign", "false")
	return repoDir
}

func gitWorktreeTest(t *testing.T, workDir string, args ...string) string {
	t.Helper()
	cmdArgs := append([]string{"-C", workDir}, gitConfigArgsForWorktreeTest()...)
	cmdArgs = append(cmdArgs, "-c", "core.fsmonitor=false")
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Env = isolatedGitEnvForWorktreeTest()
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, string(output))
	}
	return string(output)
}

func gitConfigArgsForWorktreeTest() []string {
	return []string{
		"-c", "user.name=Roundfix Test",
		"-c", "user.email=test@example.com",
		"-c", "commit.gpgsign=false",
	}
}

func isolatedGitEnvForWorktreeTest() []string {
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

func gitStatus(t *testing.T, workDir string) string {
	t.Helper()
	return gitWorktreeTest(t, workDir, "status", "--porcelain=v1")
}

func assertNoPhantomStagedEntries(t *testing.T, status string) {
	t.Helper()
	for _, line := range strings.Split(strings.TrimRight(status, "\n"), "\n") {
		if line == "" || strings.HasPrefix(line, "?? ") {
			continue
		}
		if len(line) >= 2 && line[0] != ' ' {
			t.Fatalf("expected no staged phantom entry, got status %q", status)
		}
	}
}

func branchExists(t *testing.T, repoDir string, branch string) bool {
	t.Helper()
	return strings.TrimSpace(gitWorktreeTest(t, repoDir, "branch", "--list", branch)) != ""
}

func assertPathExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected path %q to exist: %v", path, err)
	}
}

func assertPathRemoved(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected path %q removed, stat err=%v", path, err)
	}
}

func assertRunBranchExists(t *testing.T, repoDir string, branch string) {
	t.Helper()
	if !branchExists(t, repoDir, branch) {
		t.Fatalf("expected branch %q to exist", branch)
	}
}

func assertBranchRemoved(t *testing.T, repoDir string, branch string) {
	t.Helper()
	if branchExists(t, repoDir, branch) {
		t.Fatalf("expected branch %q removed", branch)
	}
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	original := os.Stderr
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("capture stderr pipe: %v", err)
	}
	os.Stderr = writePipe
	output := make(chan string, 1)
	go func() {
		var buffer bytes.Buffer
		_, _ = buffer.ReadFrom(readPipe)
		output <- buffer.String()
	}()

	fn()

	if err := writePipe.Close(); err != nil {
		t.Fatalf("close captured stderr writer: %v", err)
	}
	os.Stderr = original
	got := <-output
	if err := readPipe.Close(); err != nil {
		t.Fatalf("close captured stderr reader: %v", err)
	}
	return got
}

func mustWriteWorktreeTest(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func mustMkdirWorktreeTest(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustReadWorktreeTest(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func hasLowerHexSuffix(value string, length int) bool {
	if len(value) < length {
		return false
	}
	suffix := value[len(value)-length:]
	for _, char := range suffix {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			return false
		}
	}
	return true
}
