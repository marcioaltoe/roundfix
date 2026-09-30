// Suite: retained terminal Runs whose recorded checkout vanished.
// Invariant: a missing checkout falls back to safe repository keys without hiding unsafe paths.
// Boundary IN: retained-Run counting, filesystem validation, and the Git runner boundary.
// Boundary OUT: public runs-list rendering, owned by internal/cli/runs_list_vanished_checkout_test.go.
package worktree

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestCountRetainedTerminalRunsListsBranchesThroughTheKeyWhenTheCheckoutIsGone(t *testing.T) {
	t.Parallel()
	checkout := filepath.Join(canonicalPath(t.TempDir()), "removed-checkout")
	key := canonicalPath(t.TempDir())
	runs := []store.Run{
		terminalRunForRetainedCount("retained", checkout, key, ""),
		terminalRunForRetainedCount("released", checkout, key, ""),
	}
	runner := &retainedVanishedCheckoutGitRunner{run: func(workDir string, args []string) (string, error) {
		if workDir != key {
			return "", fmt.Errorf("unexpected Git work directory %q", workDir)
		}
		switch {
		case slices.Equal(args, []string{"rev-parse", "--show-toplevel"}):
			return key + "\n", nil
		case slices.Equal(args, retainedRunBranchListArgs()):
			return BranchName("retained") + "\n", nil
		default:
			return "", fmt.Errorf("unexpected Git arguments: %v", args)
		}
	}}

	retained, failures := countRetainedTerminalRuns(context.Background(), runner, runs)

	if retained != 1 {
		t.Fatalf("retained terminal Runs = %d, want 1", retained)
	}
	if len(failures) != 0 {
		t.Fatalf("retained terminal Run failures = %v, want none", failures)
	}
	if got := runner.countCalls(key, retainedRunBranchListArgs()); got != 1 {
		t.Fatalf("repository-key branch listings = %d, want 1; calls=%v", got, runner.calls)
	}
}

func TestCountRetainedTerminalRunsCountsOnlyAnExistingWorktreeWhenTheRepositoryIsGone(t *testing.T) {
	t.Parallel()
	checkout := filepath.Join(canonicalPath(t.TempDir()), "removed-checkout")
	key := filepath.Join(canonicalPath(t.TempDir()), "removed-repository")
	worktree := canonicalPath(t.TempDir())
	runs := []store.Run{
		terminalRunForRetainedCount("worktree-retained", checkout, key, worktree),
		terminalRunForRetainedCount("released", checkout, key, filepath.Join(canonicalPath(t.TempDir()), "removed-run-worktree")),
	}
	runner := &retainedVanishedCheckoutGitRunner{run: func(workDir string, args []string) (string, error) {
		return "", fmt.Errorf("unexpected Git call in %q: %v", workDir, args)
	}}

	retained, failures := countRetainedTerminalRuns(context.Background(), runner, runs)

	if retained != 1 {
		t.Fatalf("retained terminal Runs = %d, want 1", retained)
	}
	if len(failures) != 0 {
		t.Fatalf("retained terminal Run failures = %v, want none", failures)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("absent repositories triggered Git calls: %v", runner.calls)
	}
}

func TestCountRetainedTerminalRunsListsABareRepositoryKey(t *testing.T) {
	t.Parallel()
	checkout := filepath.Join(canonicalPath(t.TempDir()), "removed-checkout")
	key := canonicalPath(t.TempDir())
	run := terminalRunForRetainedCount("bare-retained", checkout, key, "")
	runner := &retainedVanishedCheckoutGitRunner{run: func(workDir string, args []string) (string, error) {
		if workDir != key {
			return "", fmt.Errorf("unexpected Git work directory %q", workDir)
		}
		switch {
		case slices.Equal(args, []string{"rev-parse", "--show-toplevel"}):
			return "", errors.New("this operation must be run in a work tree")
		case slices.Equal(args, []string{"rev-parse", "--absolute-git-dir"}):
			return key + "\n", nil
		case slices.Equal(args, retainedRunBranchListArgs()):
			return BranchName(run.ID) + "\n", nil
		default:
			return "", fmt.Errorf("unexpected Git arguments: %v", args)
		}
	}}

	retained, failures := countRetainedTerminalRuns(context.Background(), runner, []store.Run{run})

	if retained != 1 {
		t.Fatalf("retained terminal Runs = %d, want 1", retained)
	}
	if len(failures) != 0 {
		t.Fatalf("retained terminal Run failures = %v, want none", failures)
	}
}

func TestCountRetainedTerminalRunsStillReportsASymlinkedRecordedRoot(t *testing.T) {
	t.Parallel()
	target := canonicalPath(t.TempDir())
	root := filepath.Join(canonicalPath(t.TempDir()), "symlinked-checkout")
	if err := os.Symlink(target, root); err != nil {
		t.Fatalf("create recorded-root symlink: %v", err)
	}
	run := terminalRunForRetainedCount("symlinked", root, target, "")
	runner := &retainedVanishedCheckoutGitRunner{}

	retained, failures := countRetainedTerminalRuns(context.Background(), runner, []store.Run{run})

	if retained != 0 {
		t.Fatalf("retained terminal Runs = %d, want 0", retained)
	}
	assertSingleRetainedRootFailure(t, failures, "contains a symlink")
	if len(runner.calls) != 0 {
		t.Fatalf("symlinked recorded root triggered Git calls: %v", runner.calls)
	}
}

func TestCountRetainedTerminalRunsStillReportsARecordedRootThatIsNotADirectory(t *testing.T) {
	t.Parallel()
	root := filepath.Join(canonicalPath(t.TempDir()), "checkout-file")
	if err := os.WriteFile(root, []byte("not a directory\n"), 0o644); err != nil {
		t.Fatalf("write recorded-root file: %v", err)
	}
	run := terminalRunForRetainedCount("not-a-directory", root, "", "")

	retained, failures := countRetainedTerminalRuns(context.Background(), &retainedVanishedCheckoutGitRunner{}, []store.Run{run})

	if retained != 0 {
		t.Fatalf("retained terminal Runs = %d, want 0", retained)
	}
	assertSingleRetainedRootFailure(t, failures, "is not a real directory")
}

func TestCountRetainedTerminalRunsStillReportsARecordedRootThatIsNotTheRepositoryRoot(t *testing.T) {
	t.Parallel()
	root := canonicalPath(t.TempDir())
	otherRoot := canonicalPath(t.TempDir())
	run := terminalRunForRetainedCount("not-the-root", root, "", "")
	runner := &retainedVanishedCheckoutGitRunner{run: func(_ string, args []string) (string, error) {
		if slices.Equal(args, []string{"rev-parse", "--show-toplevel"}) {
			return otherRoot + "\n", nil
		}
		return "", fmt.Errorf("unexpected Git arguments: %v", args)
	}}

	retained, failures := countRetainedTerminalRuns(context.Background(), runner, []store.Run{run})

	if retained != 0 {
		t.Fatalf("retained terminal Runs = %d, want 0", retained)
	}
	assertSingleRetainedRootFailure(t, failures, "is not the repository root")
}

func TestCountRetainedTerminalRunsStillReportsARecordedRootGitFailure(t *testing.T) {
	t.Parallel()
	root := canonicalPath(t.TempDir())
	run := terminalRunForRetainedCount("git-failure", root, "", "")
	runner := &retainedVanishedCheckoutGitRunner{run: func(_ string, args []string) (string, error) {
		if slices.Equal(args, []string{"rev-parse", "--show-toplevel"}) {
			return "", errors.New("git inspection failed")
		}
		return "", fmt.Errorf("unexpected Git arguments: %v", args)
	}}

	retained, failures := countRetainedTerminalRuns(context.Background(), runner, []store.Run{run})

	if retained != 0 {
		t.Fatalf("retained terminal Runs = %d, want 0", retained)
	}
	assertSingleRetainedRootFailure(t, failures, "git inspection failed")
}

func TestCountRetainedTerminalRunsReportsAnInvalidExistingRepositoryKey(t *testing.T) {
	t.Parallel()
	checkout := filepath.Join(canonicalPath(t.TempDir()), "removed-checkout")
	key := canonicalPath(t.TempDir())
	otherRoot := canonicalPath(t.TempDir())
	run := terminalRunForRetainedCount("invalid-key", checkout, key, "")
	runner := &retainedVanishedCheckoutGitRunner{run: func(_ string, args []string) (string, error) {
		switch {
		case slices.Equal(args, []string{"rev-parse", "--show-toplevel"}):
			return otherRoot + "\n", nil
		case slices.Equal(args, []string{"rev-parse", "--absolute-git-dir"}):
			return filepath.Join(key, ".git") + "\n", nil
		default:
			return "", fmt.Errorf("unexpected Git arguments: %v", args)
		}
	}}

	retained, failures := countRetainedTerminalRuns(context.Background(), runner, []store.Run{run})

	if retained != 0 {
		t.Fatalf("retained terminal Runs = %d, want 0", retained)
	}
	assertSingleRetainedRootFailure(t, failures, "is not the repository root")
}

type retainedVanishedCheckoutGitCall struct {
	workDir string
	args    []string
}

type retainedVanishedCheckoutGitRunner struct {
	run   func(workDir string, args []string) (string, error)
	calls []retainedVanishedCheckoutGitCall
}

func (runner *retainedVanishedCheckoutGitRunner) Run(_ context.Context, workDir string, args ...string) (string, error) {
	callArgs := append([]string(nil), args...)
	runner.calls = append(runner.calls, retainedVanishedCheckoutGitCall{workDir: workDir, args: callArgs})
	if runner.run == nil {
		return "", fmt.Errorf("unexpected Git call in %q: %v", workDir, args)
	}
	return runner.run(workDir, callArgs)
}

func (runner *retainedVanishedCheckoutGitRunner) countCalls(workDir string, args []string) int {
	count := 0
	for _, call := range runner.calls {
		if call.workDir == workDir && slices.Equal(call.args, args) {
			count++
		}
	}
	return count
}

func terminalRunForRetainedCount(id, gitRoot, repositoryRoot, workDir string) store.Run {
	return store.Run{
		ID:             id,
		Kind:           store.KindImplement,
		State:          store.StateStopped,
		GitRoot:        gitRoot,
		RepositoryRoot: repositoryRoot,
		WorkDir:        workDir,
	}
}

func retainedRunBranchListArgs() []string {
	return []string{"for-each-ref", "--format=%(refname:short)", "refs/heads/" + runBranchPrefix + "*"}
}

func assertSingleRetainedRootFailure(t *testing.T, failures []error, want string) {
	t.Helper()
	if len(failures) != 1 {
		t.Fatalf("retained terminal Run failures = %v, want one", failures)
	}
	if !strings.Contains(failures[0].Error(), want) {
		t.Fatalf("retained terminal Run failure = %q, want text %q", failures[0], want)
	}
}
