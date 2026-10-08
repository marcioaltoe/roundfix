// Suite: Task Carry-Forward staging hooks.
// Invariant: carry-forward proves and applies settled Tasks without running repository hooks or changing hook configuration.
// Boundary IN: the public reconcile runner, real local Git worktrees, repository hooks, and carry-forward staging cleanup.
// Boundary OUT: Daemon settlement hooks, owned by internal/cli/settle_test.go.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/store"
)

const carryForwardHookMarkerEnv = "ROUNDFIX_TEST_CARRY_FORWARD_HOOK_MARKER"

func TestCarryForwardAppliesThroughRefusingCommitHooks(t *testing.T) {
	// Sequential: sets process-wide TMPDIR and carry-forward hook marker environment variables.
	hooks := writeCarryForwardHookFixtures(t, true)
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01", title: "Build the core"}})
	configureCarryForwardHooks(t, fixture.repoDir, hooks.directory)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", fixture.run.ID, "--carry-forward", "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitOK {
		t.Fatalf("carry-forward with refusing hooks exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if got := mustRead(t, filepath.Join(fixture.repoDir, "src", "task_01.txt")); got != "task_01 settled\n" {
		t.Fatalf("carried implementation = %q, want task_01 settlement", got)
	}
}

func TestCarryForwardProofIgnoresRefusingCommitHooks(t *testing.T) {
	// Sequential: sets process-wide TMPDIR and carry-forward hook marker environment variables.
	hooks := writeCarryForwardHookFixtures(t, true)
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01", title: "Build the core"}})
	configureCarryForwardHooks(t, fixture.repoDir, hooks.directory)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", fixture.run.ID, "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitOK {
		t.Fatalf("carry-forward proof with refusing hooks exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	var report reconcileReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode carry-forward proof: %v\n%s", err, stdout.String())
	}
	if len(report.CarryForwards) != 1 || report.CarryForwards[0].Action != carryForwardReadyAction {
		t.Fatalf("carry-forward proof = %+v, want one ready candidate", report.CarryForwards)
	}
}

func TestCarryForwardStagingRunsNoRepositoryHook(t *testing.T) {
	// Sequential: sets process-wide TMPDIR and carry-forward hook marker environment variables.
	hooks := writeCarryForwardHookFixtures(t, false)
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01", title: "Build the core"}})
	configureCarryForwardHooks(t, fixture.repoDir, hooks.directory)

	// Prove that every fixture hook is active before clearing its marker. The
	// carry-forward that follows must not add the marker back.
	gitImplement(t, fixture.repoDir, "commit", "--allow-empty", "-m", "exercise repository hooks")
	assertCarryForwardHookMarkers(t, hooks.marker, []string{"pre-commit", "prepare-commit-msg", "commit-msg", "post-commit"})
	if err := os.Remove(hooks.marker); err != nil {
		t.Fatalf("clear hook marker: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", fixture.run.ID, "--carry-forward", "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitOK {
		t.Fatalf("carry-forward with marker hooks exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if _, err := os.Stat(hooks.marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("hook marker stat error = %v, want no repository hook to run during staging", err)
	}
}

func TestCarryForwardLeavesTheCheckoutHooksPathUnchanged(t *testing.T) {
	// Sequential: sets process-wide TMPDIR and carry-forward hook marker environment variables.
	hooks := writeCarryForwardHookFixtures(t, true)
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01", title: "Build the core"}})
	configureCarryForwardHooks(t, fixture.repoDir, hooks.directory)
	before := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "config", "--local", "--get", "core.hooksPath"))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", fixture.run.ID, "--carry-forward", "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitOK {
		t.Fatalf("carry-forward while preserving hooks path exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	after := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "config", "--local", "--get", "core.hooksPath"))
	if after != before {
		t.Fatalf("checkout core.hooksPath = %q after carry-forward, want unchanged %q", after, before)
	}
}

func TestCarryForwardRemovesItsEmptyHooksDirectory(t *testing.T) {
	// Sequential: sets process-wide TMPDIR and carry-forward hook marker environment variables.
	hooks := writeCarryForwardHookFixtures(t, true)
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01", title: "Build the core"}})
	configureCarryForwardHooks(t, fixture.repoDir, hooks.directory)
	tempRoot := t.TempDir()
	t.Setenv("TMPDIR", tempRoot)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", fixture.run.ID, "--carry-forward", "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitOK {
		t.Fatalf("carry-forward hooks cleanup exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	entries, err := os.ReadDir(tempRoot)
	if err != nil {
		t.Fatalf("read carry-forward temp root: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("carry-forward temporary directories remain after cleanup: %v", entries)
	}
}

type carryForwardHookFixtures struct {
	directory string
	marker    string
}

func writeCarryForwardHookFixtures(t *testing.T, refuseCommits bool) carryForwardHookFixtures {
	t.Helper()
	root := t.TempDir()
	directory := filepath.Join(root, "hooks")
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatalf("create carry-forward hook fixtures: %v", err)
	}
	marker := filepath.Join(root, "hook-markers.txt")
	t.Setenv(carryForwardHookMarkerEnv, marker)
	for _, name := range []string{"pre-commit", "prepare-commit-msg", "commit-msg", "post-commit"} {
		exit := "exit 0\n"
		if refuseCommits && (name == "pre-commit" || name == "commit-msg") {
			exit = "exit 1\n"
		}
		content := "#!/bin/sh\n" +
			"printf '%s\\n' " + name + " >> \"$" + carryForwardHookMarkerEnv + "\"\n" +
			exit
		writeScriptFixture(t, filepath.Join(directory, name), content)
	}
	return carryForwardHookFixtures{directory: directory, marker: marker}
}

func configureCarryForwardHooks(t *testing.T, repository string, hooksDirectory string) {
	t.Helper()
	gitImplement(t, repository, "config", "--local", "core.hooksPath", hooksDirectory)
}

func assertCarryForwardHookMarkers(t *testing.T, marker string, want []string) {
	t.Helper()
	got := strings.Fields(mustRead(t, marker))
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Fatalf("repository hook markers = %v, want %v", got, want)
	}
}
