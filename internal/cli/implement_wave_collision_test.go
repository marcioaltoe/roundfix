// Suite: Implement Command Wave Collision preflight.
// Invariant: unsafe parallel graphs create no Run; completed and serial graphs reach dispatch.
// Boundary IN: committed temporary repositories, isolated configuration, and a fake Agent.
// Boundary OUT: real Agent sessions, network access, and successful Task settlement.
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestImplementRefusesAWaveCollisionBeforeTheRun(t *testing.T) {
	t.Parallel()
	home, _, runner := newImplementWaveCollisionWorkspace(t, spec.StatusPending, 2)
	var stdout, stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	wantStderr := "Preflight failed\n\nReason:\n" +
		"  Tasks would collide in one Wave at Task Capacity 2: Task task_01 and Task task_02 share repository files: shared/notes.txt (verification command); add a `needs` edge from Task task_02 to Task task_01\n\n" +
		"No side effects:\n" +
		"  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.\n\n" +
		"Usage:\n" +
		"  Run 'roundfix implement --help' for usage.\n"
	if code != 2 || stdout.String() != "" || stderr.String() != wantStderr {
		t.Fatalf("transcript mismatch: exit=%d stdout=%q stderr=%q; want exit=2 stdout=%q stderr=%q", code, stdout.String(), stderr.String(), "", wantStderr)
	}
	assertRunCount(t, store.DatabasePath(home), 0)
	if runner.calls != 0 || len(runner.probeRequests) != 0 {
		t.Fatalf("preflight reached Agent: calls=%d probes=%d", runner.calls, len(runner.probeRequests))
	}
}

func TestImplementRunsATaskThatSharesAFileWithACompletedTask(t *testing.T) {
	t.Parallel()
	home, _, runner := newImplementWaveCollisionWorkspace(t, spec.StatusCompleted, 2)
	var stdout, stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed || strings.Contains(strings.ToLower(stderr.String()), "collid") {
		t.Fatalf("expected Run to reach scripted Agent failure: exit=%d stderr=%q", code, stderr.String())
	}
	assertRunCount(t, store.DatabasePath(home), 1)
	if runner.calls != 1 || len(runner.taskIDs) != 1 || runner.taskIDs[0] != "task_02" {
		t.Fatalf("expected only pending task_02 to dispatch: calls=%d Tasks=%v", runner.calls, runner.taskIDs)
	}
}

func TestImplementAtTaskCapacityOneLeavesCollidingTasksToTheRun(t *testing.T) {
	t.Parallel()
	home, _, runner := newImplementWaveCollisionWorkspace(t, spec.StatusPending, 1)
	var stdout, stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed || strings.Contains(strings.ToLower(stderr.String()), "collid") {
		t.Fatalf("expected serial Run to reach scripted Agent failure: exit=%d stderr=%q", code, stderr.String())
	}
	assertRunCount(t, store.DatabasePath(home), 1)
	if runner.calls == 0 {
		t.Fatal("serial Run dispatched no Task")
	}
}

func newImplementWaveCollisionWorkspace(t *testing.T, firstStatus spec.Status, capacity int) (string, string, *implementFakeRunner) {
	t.Helper()
	home, repo := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", status: string(firstStatus), verification: []string{"grep -q task_01 shared/notes.txt"}},
		{id: "task_02", verification: []string{"grep -q task_02 shared/notes.txt"}},
	})
	// Stderr names both paths. Keep "Colliding" from the test name out of
	// those paths so the no-collision assertion reads product diagnostics.
	root, err := os.MkdirTemp("", "implement-wave-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(root); err != nil {
			t.Error(err)
		}
	})
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, move := range []struct{ from, to string }{{home, filepath.Join(root, "home")}, {repo, filepath.Join(root, "repo")}} {
		if err := os.Rename(move.from, move.to); err != nil {
			t.Fatal(err)
		}
	}
	home, repo = filepath.Join(root, "home"), filepath.Join(root, "repo")
	setCommandEnvironmentForTest(t, home, repo)
	writeUserConfig(t, home, fmt.Sprintf("worktree:\n  concurrency: %d\n", capacity))
	mustMkdir(t, filepath.Join(repo, "shared"))
	mustWrite(t, filepath.Join(repo, "shared", "notes.txt"), "pending\n")
	gitImplement(t, repo, "add", "shared/notes.txt")
	gitImplement(t, repo, "commit", "-m", "seed shared Verification operand")
	// Stop at dispatch: these tests prove admission, not successful settlement.
	runner := &implementFakeRunner{
		gitRoot: repo,
		errByTask: map[string]error{
			"task_01": errors.New("scripted stop after dispatch"),
			"task_02": errors.New("scripted stop after dispatch"),
		},
	}
	withAgentRunner(t, runner)
	return home, repo, runner
}
