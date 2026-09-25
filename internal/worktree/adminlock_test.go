package worktree

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"
)

// Suite: Git worktree administration locking.
// Invariant: one repository admits one worktree administrative command at a time.
// Boundary IN: internal/worktree command dispatch and the operating-system file lock.
// Boundary OUT: Git's own worktree implementation and higher-level Run scheduling.

func TestConcurrentTaskWorktreeCreationIsSerialized(t *testing.T) {
	t.Parallel()

	const taskCount = 12
	commonDir := t.TempDir()
	userRoot := t.TempDir()
	run := Ref{
		RunID:    "concurrent-task-creation",
		Path:     filepath.Join(t.TempDir(), "run"),
		Branch:   BranchName("concurrent-task-creation"),
		UserRoot: userRoot,
	}
	runner := newConcurrentAdminRecordingRunner(commonDir, taskCount)
	defer close(runner.release)

	results := make(chan error, taskCount)
	for index := range taskCount {
		go func() {
			_, err := createTaskWithOptions(
				context.Background(),
				runner,
				run,
				fmt.Sprintf("task_%02d", index),
				TaskCreateOptions{Concurrency: taskCount},
			)
			results <- err
		}()
	}

	for range taskCount {
		waitForAdminCommandStart(t, runner.started)
		runner.release <- struct{}{}
	}
	for range taskCount {
		if err := <-results; err != nil {
			t.Fatalf("create concurrent Task Worktree: %v", err)
		}
	}
	if got := runner.maximumActive(); got != 1 {
		t.Fatalf("maximum concurrent worktree administrative commands = %d, want 1", got)
	}
	if got := runner.commandCount(); got != taskCount {
		t.Fatalf("worktree administrative command count = %d, want %d", got, taskCount)
	}
}

func TestWorktreeAdministrationIsSerializedAcrossProcesses(t *testing.T) {
	t.Parallel()

	commonDir := t.TempDir()
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	releaseHolder := sync.OnceFunc(func() { close(release) })
	defer releaseHolder()
	holderRunner := &adminLockTestRunner{
		commonDir: commonDir,
		admin: func(ctx context.Context, _ string, _ ...string) (string, error) {
			started <- struct{}{}
			select {
			case <-release:
				return "", nil
			case <-ctx.Done():
				return "", ctx.Err()
			}
		},
	}
	holderResult := make(chan error, 1)
	go func() {
		_, err := runWorktreeCommand(context.Background(), holderRunner, commonDir, "worktree", "add")
		holderResult <- err
	}()
	waitForAdminCommandStart(t, started)

	command := exec.Command(os.Args[0], "-test.run=^TestWorktreeAdministrationLockHelperProcess$")
	command.Env = append(os.Environ(),
		"ROUNDFIX_WORKTREE_ADMIN_LOCK_HELPER=1",
		"ROUNDFIX_WORKTREE_ADMIN_COMMON_DIR="+commonDir,
	)
	output, err := command.CombinedOutput()
	if err != nil {
		releaseHolder()
		<-holderResult
		t.Fatalf("run worktree administration lock helper: %v\n%s", err, output)
	}

	releaseHolder()
	if err := <-holderResult; err != nil {
		t.Fatalf("release holding worktree administrative command: %v", err)
	}
}

func TestWorktreeAdministrationWaitEndsWithTheContext(t *testing.T) {
	t.Parallel()

	commonDir := t.TempDir()
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	releaseHolder := sync.OnceFunc(func() { close(release) })
	defer releaseHolder()
	holderRunner := &adminLockTestRunner{
		commonDir: commonDir,
		admin: func(ctx context.Context, _ string, _ ...string) (string, error) {
			started <- struct{}{}
			select {
			case <-release:
				return "", nil
			case <-ctx.Done():
				return "", ctx.Err()
			}
		},
	}
	holderResult := make(chan error, 1)
	go func() {
		_, err := runWorktreeCommand(context.Background(), holderRunner, commonDir, "worktree", "add")
		holderResult <- err
	}()
	waitForAdminCommandStart(t, started)

	var waiterRan bool
	waiterRunner := &adminLockTestRunner{
		commonDir: commonDir,
		admin: func(_ context.Context, _ string, _ ...string) (string, error) {
			waiterRan = true
			return "", nil
		},
	}
	waitCtx, cancelWait := context.WithCancel(context.Background())
	cancelWait()
	_, waitErr := runWorktreeCommand(waitCtx, waiterRunner, commonDir, "worktree", "remove")
	if !errors.Is(waitErr, context.Canceled) {
		t.Fatalf("waiting worktree administration error = %v, want context canceled", waitErr)
	}
	if waiterRan {
		t.Fatal("waiting worktree administrative command ran after its context ended")
	}

	releaseHolder()
	if err := <-holderResult; err != nil {
		t.Fatalf("release holding worktree administrative command: %v", err)
	}
}

func TestWorktreeAdministrationOfDifferentRepositoriesRunsConcurrently(t *testing.T) {
	t.Parallel()

	started := make(chan struct{}, 2)
	release := make(chan struct{})
	releaseAll := sync.OnceFunc(func() { close(release) })
	defer releaseAll()
	result := make(chan error, 2)
	for range 2 {
		runner := &adminLockTestRunner{
			commonDir: t.TempDir(),
			admin: func(ctx context.Context, _ string, _ ...string) (string, error) {
				started <- struct{}{}
				select {
				case <-release:
					return "", nil
				case <-ctx.Done():
					return "", ctx.Err()
				}
			},
		}
		go func() {
			_, err := runWorktreeCommand(context.Background(), runner, runner.commonDir, "worktree", "prune")
			result <- err
		}()
	}

	waitForAdminCommandStart(t, started)
	waitForAdminCommandStart(t, started)
	releaseAll()
	for range 2 {
		if err := <-result; err != nil {
			t.Fatalf("run worktree administration for a distinct repository: %v", err)
		}
	}
}

func TestWorktreeAdministrationLockHelperProcess(t *testing.T) {
	if os.Getenv("ROUNDFIX_WORKTREE_ADMIN_LOCK_HELPER") != "1" {
		return
	}
	commonDir := os.Getenv("ROUNDFIX_WORKTREE_ADMIN_COMMON_DIR")
	if commonDir == "" {
		t.Fatal("helper common directory is required")
	}

	commandRan := false
	runner := &adminLockTestRunner{
		commonDir: commonDir,
		admin: func(_ context.Context, _ string, _ ...string) (string, error) {
			commandRan = true
			return "", nil
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := runWorktreeCommand(ctx, runner, commonDir, "worktree", "move")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("helper lock wait error = %v, want context deadline exceeded", err)
	}
	if commandRan {
		t.Fatal("helper worktree administrative command overlapped the holding process")
	}
}

type adminLockTestRunner struct {
	commonDir string
	admin     func(context.Context, string, ...string) (string, error)
}

func (runner *adminLockTestRunner) Run(ctx context.Context, workDir string, args ...string) (string, error) {
	if slices.Equal(args, []string{"rev-parse", "--path-format=absolute", "--git-common-dir"}) {
		return runner.commonDir + "\n", nil
	}
	if len(args) >= 2 && args[0] == "worktree" {
		return runner.admin(ctx, workDir, args...)
	}
	return "", fmt.Errorf("unexpected Git arguments: %v", args)
}

type concurrentAdminRecordingRunner struct {
	commonDir string
	want      int

	mu          sync.Mutex
	resolved    int
	active      int
	maxActive   int
	commands    int
	allResolved chan struct{}
	started     chan struct{}
	release     chan struct{}
}

func newConcurrentAdminRecordingRunner(commonDir string, want int) *concurrentAdminRecordingRunner {
	return &concurrentAdminRecordingRunner{
		commonDir:   commonDir,
		want:        want,
		allResolved: make(chan struct{}),
		started:     make(chan struct{}, want),
		release:     make(chan struct{}),
	}
}

func (runner *concurrentAdminRecordingRunner) Run(ctx context.Context, _ string, args ...string) (string, error) {
	switch {
	case len(args) == 2 && args[0] == "rev-parse":
		return "base-sha\n", nil
	case slices.Equal(args, []string{"rev-parse", "--path-format=absolute", "--git-common-dir"}):
		runner.mu.Lock()
		runner.resolved++
		if runner.resolved == runner.want {
			close(runner.allResolved)
		}
		runner.mu.Unlock()
		select {
		case <-runner.allResolved:
			return runner.commonDir + "\n", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	case len(args) >= 2 && args[0] == "worktree" && args[1] == "add":
		runner.mu.Lock()
		runner.active++
		runner.commands++
		if runner.active > runner.maxActive {
			runner.maxActive = runner.active
		}
		runner.mu.Unlock()
		runner.started <- struct{}{}
		select {
		case <-runner.release:
			runner.mu.Lock()
			runner.active--
			runner.mu.Unlock()
			return "", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	default:
		return "", fmt.Errorf("unexpected Git arguments: %v", args)
	}
}

func (runner *concurrentAdminRecordingRunner) maximumActive() int {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	return runner.maxActive
}

func (runner *concurrentAdminRecordingRunner) commandCount() int {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	return runner.commands
}

func waitForAdminCommandStart(t *testing.T, started <-chan struct{}) {
	t.Helper()
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	select {
	case <-started:
	case <-timer.C:
		t.Fatal("timed out waiting for a worktree administrative command to start")
	}
}
