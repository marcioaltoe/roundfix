// Suite: reconcile terminal-Run dispositions.
// Invariant: reconcile reports without mutation and changes a checkout only after an explicit flag proves the whole disposition.
// Boundary IN: public reconcile runner, Run Database metadata, local Git surfaces, and Artifact Root records.
// Boundary OUT: terminal Run classification internals, owned by internal/daemon/reconcile_test.go.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/daemon"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func TestCarryForwardAcceptsAnUnresolvedRun(t *testing.T) {
	t.Parallel()
	for _, state := range []string{store.StateStopped, store.StateUnresolved} {
		state := state
		t.Run(state, func(t *testing.T) {
			t.Parallel()
			fixture := newCarryForwardFixture(t, state, []implementSeed{{id: "task_01", title: "Build the core"}})
			beforeHead := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD"))
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
				t.Fatalf("carry-forward exit = %d, want 0 stderr=%q stdout=%q", code, stderr.String(), stdout.String())
			}
			if stderr.Len() != 0 {
				t.Fatalf("carry-forward stderr = %q, want empty", stderr.String())
			}
			afterHead := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD"))
			if afterHead == beforeHead {
				t.Fatalf("carry-forward HEAD = %s, want a fast-forward", afterHead)
			}
			task := mustRead(t, implementTaskPath(fixture.repoDir, "task_01"))
			for _, want := range []string{"status: completed", fixture.run.ID, fixture.commits["task_01"]} {
				if !strings.Contains(task, want) {
					t.Fatalf("carried task does not contain %q:\n%s", want, task)
				}
			}
			if got := mustRead(t, filepath.Join(fixture.repoDir, "src", "task_01.txt")); got != "task_01 settled\n" {
				t.Fatalf("carried implementation = %q", got)
			}
		})
	}
}

func TestCarryForwardRefusesATaskWhoseInputsMoved(t *testing.T) {
	t.Parallel()
	for _, state := range []string{store.StateStopped, store.StateUnresolved} {
		state := state
		t.Run(state, func(t *testing.T) {
			t.Parallel()
			fixture := newCarryForwardFixture(t, state, []implementSeed{{id: "task_01", title: "Build the core"}})
			prdPath := filepath.Join(fixture.repoDir, "docs", "specs", implementTestSlug, "_prd.md")
			mustWrite(t, prdPath, mustRead(t, prdPath)+"\nMoved after settlement.\n")
			beforeTask := mustRead(t, implementTaskPath(fixture.repoDir, "task_01"))
			beforeHead := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD"))
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runCLIContext(
				t,
				context.Background(),
				[]string{"reconcile", fixture.run.ID, "--carry-forward", "--format=json"},
				&stdout,
				&stderr,
			)

			if code != exitPreflight {
				t.Fatalf("moved-input exit = %d, want %d stderr=%q stdout=%q", code, exitPreflight, stderr.String(), stdout.String())
			}
			for _, want := range []string{"task_01", filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "_prd.md"))} {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("moved-input stderr does not name %q: %q", want, stderr.String())
				}
			}
			if got := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD")); got != beforeHead {
				t.Fatalf("moved-input HEAD = %s, want unchanged %s", got, beforeHead)
			}
			if got := mustRead(t, implementTaskPath(fixture.repoDir, "task_01")); got != beforeTask {
				t.Fatal("moved-input refusal changed the Task file")
			}
			if _, err := os.Stat(filepath.Join(fixture.repoDir, "src", "task_01.txt")); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("moved-input implementation stat error = %v, want not exist", err)
			}
		})
	}
}

func TestCarryForwardRefusesRatherThanCarryingASubset(t *testing.T) {
	t.Parallel()
	for _, state := range []string{store.StateStopped, store.StateUnresolved} {
		state := state
		t.Run(state, func(t *testing.T) {
			t.Parallel()
			fixture := newCarryForwardFixture(t, state, []implementSeed{
				{id: "task_01", title: "Build the core"},
				{id: "task_02", title: "Wire the shell", needs: []string{"task_01"}},
			})
			secondTaskPath := implementTaskPath(fixture.repoDir, "task_02")
			mustWrite(t, secondTaskPath, mustRead(t, secondTaskPath)+"\nMoved after settlement.\n")
			beforeFirst := mustRead(t, implementTaskPath(fixture.repoDir, "task_01"))
			beforeHead := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD"))
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runCLIContext(
				t,
				context.Background(),
				[]string{"reconcile", fixture.run.ID, "--carry-forward", "--format=json"},
				&stdout,
				&stderr,
			)

			if code != exitPreflight {
				t.Fatalf("mixed-set exit = %d, want %d stderr=%q stdout=%q", code, exitPreflight, stderr.String(), stdout.String())
			}
			if !strings.Contains(stderr.String(), "task_02") || !strings.Contains(stderr.String(), filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "task_02.md"))) {
				t.Fatalf("mixed-set stderr = %q, want task_02 and its moved input", stderr.String())
			}
			if got := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD")); got != beforeHead {
				t.Fatalf("mixed-set HEAD = %s, want unchanged %s", got, beforeHead)
			}
			if got := mustRead(t, implementTaskPath(fixture.repoDir, "task_01")); got != beforeFirst {
				t.Fatal("mixed-set refusal carried task_01")
			}
			for _, taskID := range []string{"task_01", "task_02"} {
				if _, err := os.Stat(filepath.Join(fixture.repoDir, "src", taskID+".txt")); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("mixed-set %s implementation stat error = %v, want not exist", taskID, err)
				}
			}
		})
	}
}

func TestCarryForwardWithoutTheFlagReportsAndChangesNothing(t *testing.T) {
	t.Parallel()
	for _, state := range []string{store.StateStopped, store.StateUnresolved} {
		state := state
		t.Run(state, func(t *testing.T) {
			t.Parallel()
			fixture := newCarryForwardFixture(t, state, []implementSeed{{id: "task_01", title: "Build the core"}})
			beforeHead := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD"))
			beforeTask := mustRead(t, implementTaskPath(fixture.repoDir, "task_01"))
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
				t.Fatalf("carry-forward report exit = %d, want 0 stderr=%q stdout=%q", code, stderr.String(), stdout.String())
			}
			var report reconcileReport
			if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
				t.Fatalf("decode carry-forward report: %v\n%s", err, stdout.String())
			}
			if len(report.CarryForwards) != 1 || report.CarryForwards[0].TaskID != "task_01" || report.CarryForwards[0].Action != "would carry forward with --carry-forward" {
				t.Fatalf("carry-forward report = %+v", report.CarryForwards)
			}
			if got := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD")); got != beforeHead {
				t.Fatalf("read-only carry-forward HEAD = %s, want %s", got, beforeHead)
			}
			if got := mustRead(t, implementTaskPath(fixture.repoDir, "task_01")); got != beforeTask {
				t.Fatal("read-only carry-forward report changed the Task file")
			}
		})
	}
}

func TestCarryForwardRefusesAnUnacceptedRunOutcomeByName(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateClean, []implementSeed{{id: "task_01", title: "Build the core"}})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", fixture.run.ID, "--carry-forward", "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitPreflight {
		t.Fatalf("unaccepted-outcome exit = %d, want %d stderr=%q stdout=%q", code, exitPreflight, stderr.String(), stdout.String())
	}
	for _, want := range []string{store.StateClean, store.StateStopped, store.StateUnresolved} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("unaccepted-outcome stderr does not name %q: %q", want, stderr.String())
		}
	}
}

type carryForwardFixture struct {
	homeDir string
	repoDir string
	run     store.Run
	ref     runworktree.Ref
	commits map[string]string
}

func newCarryForwardFixture(t *testing.T, state string, seeds []implementSeed) carryForwardFixture {
	t.Helper()
	homeDir, repoDir := newImplementWorkspace(t, seeds)
	ctx := context.Background()
	head := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open carry-forward Run Database: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close carry-forward Run Database: %v", err)
		}
	}()
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "ma/widget-flow",
		HeadSHA:     head,
		SpecSlug:    implementTestSlug,
		Agent:       "codex",
	})
	if err != nil {
		t.Fatalf("create carry-forward Run: %v", err)
	}
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: repoDir,
		Location: t.TempDir(),
		RunID:    run.ID,
		HeadSHA:  head,
	})
	if err != nil {
		t.Fatalf("create carry-forward Run Worktree: %v", err)
	}
	ref.Path, err = filepath.EvalSymlinks(ref.Path)
	if err != nil {
		t.Fatalf("resolve carry-forward Run Worktree: %v", err)
	}
	run, err = runStore.SetRunWorkDir(ctx, run.ID, ref.Path)
	if err != nil {
		t.Fatalf("record carry-forward Run Worktree: %v", err)
	}
	commits := make(map[string]string, len(seeds))
	for _, seed := range seeds {
		implementationPath := filepath.Join("src", seed.id+".txt")
		mustMkdir(t, filepath.Join(ref.Path, filepath.Dir(implementationPath)))
		mustWrite(t, filepath.Join(ref.Path, implementationPath), seed.id+" settled\n")
		taskPath := implementTaskPath(ref.Path, seed.id)
		if err := spec.SetStatus(taskPath, spec.StatusCompleted); err != nil {
			t.Fatalf("settle source %s: %v", seed.id, err)
		}
		settlementPaths := []string{implementationPath, filepath.Join("docs", "specs", implementTestSlug, seed.id+".md")}
		for path, content := range seed.settlementFiles {
			mustMkdir(t, filepath.Join(ref.Path, filepath.Dir(path)))
			mustWrite(t, filepath.Join(ref.Path, filepath.FromSlash(path)), content)
			settlementPaths = append(settlementPaths, filepath.FromSlash(path))
		}
		slices.Sort(settlementPaths)
		gitImplement(t, ref.Path, append([]string{"add"}, settlementPaths...)...)
		gitImplement(
			t,
			ref.Path,
			"commit",
			"-m", fmt.Sprintf("feat: settle %s", seed.id),
			"-m", fmt.Sprintf("Roundfix-Spec: %s\nRoundfix-Task: %s", implementTestSlug, seed.id),
		)
		commits[seed.id] = strings.TrimSpace(gitImplementOutput(t, ref.Path, "rev-parse", "HEAD"))
		if _, err := runStore.AppendRunEvents(ctx, []runevent.RunEvent{
			{
				RunID:       run.ID,
				Source:      runevent.SourceDaemon,
				Kind:        runevent.KindDaemonVerification,
				ReviewIssue: seed.id,
				Payload:     []byte(fmt.Sprintf(`{"attempt":1,"phase":"verdict","task":%q,"verdict":"passed"}`, seed.id)),
			},
			{
				RunID:       run.ID,
				Source:      runevent.SourceDaemon,
				Kind:        runevent.KindDaemonTask,
				ReviewIssue: seed.id,
				Payload:     []byte(fmt.Sprintf(`{"task":%q,"phase":"settled","status":"completed"}`, seed.id)),
			},
		}); err != nil {
			t.Fatalf("append carry-forward evidence for %s: %v", seed.id, err)
		}
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, state)
	if err != nil {
		t.Fatalf("stop carry-forward Run: %v", err)
	}
	return carryForwardFixture{homeDir: homeDir, repoDir: repoDir, run: completed.Run, ref: ref, commits: commits}
}

func TestReconcileDiscardWritesTheRecordBeforeRemoving(t *testing.T) {
	t.Parallel()
	homeDir, repoDir, location := newReconcileWorkspace(t)
	run, ref := createReconcileRun(t, homeDir, repoDir, location, "ma/widget-flow", store.StateStopped)
	first := commitReconcileDispositionChange(t, ref.Path, "first.txt", "first\n", "feat: add first change")
	second := commitReconcileDispositionChange(t, ref.Path, "second.txt", "second\n", "feat: add second change")
	gitImplement(t, repoDir, "merge", "--ff-only", ref.Branch)
	recordPath := filepath.Join(builtinArtifactDirForRepo(t, repoDir), "runs", run.ID, "branch-disposition.json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", run.ID, "--discard-superseded", "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitOK {
		t.Fatalf("discard exit = %d, want 0 stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("discard stderr = %q, want empty", stderr.String())
	}
	var report reconcileReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode discard report: %v\n%s", err, stdout.String())
	}
	if len(report.Results) != 1 {
		t.Fatalf("discard results = %+v, want one", report.Results)
	}
	result := report.Results[0]
	if result.Classification != "superseded" || result.Action != "discarded" ||
		!strings.Contains(result.Evidence, "branch record: "+recordPath) {
		t.Fatalf("discard result = %+v", result)
	}
	record, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatalf("read branch record after removal: %v", err)
	}
	for _, want := range []string{daemon.BranchDispositionRecordSchemaVersion, first, second, "first.txt", "second.txt"} {
		if !bytes.Contains(record, []byte(want)) {
			t.Fatalf("branch record does not contain %q: %s", want, record)
		}
	}
	assertReconcilePathState(t, ref.Path, false)
	assertReconcileBranchState(t, repoDir, ref.Branch, false)
}

func TestReconcileDiscardRefusesAnUnreachableCommit(t *testing.T) {
	t.Parallel()
	homeDir, repoDir, location := newReconcileWorkspace(t)
	run, ref := createReconcileRun(t, homeDir, repoDir, location, "ma/widget-flow", store.StateStopped)
	unreachable := commitReconcileDispositionChange(t, ref.Path, "run-only.txt", "run only\n", "feat: retain run-only change")
	commitReconcileDispositionChange(t, repoDir, "target-only.txt", "target only\n", "feat: advance target")
	recordPath := filepath.Join(builtinArtifactDirForRepo(t, repoDir), "runs", run.ID, "branch-disposition.json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", run.ID, "--discard-superseded", "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitPreflight {
		t.Fatalf("refused discard exit = %d, want %d stderr=%q stdout=%q", code, exitPreflight, stderr.String(), stdout.String())
	}
	var report reconcileReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode refused discard report: %v\n%s", err, stdout.String())
	}
	if len(report.Results) != 1 ||
		!strings.Contains(report.Results[0].RefusalReason, "unreachable commit "+unreachable) {
		t.Fatalf("refused discard results = %+v", report.Results)
	}
	assertReconcilePathState(t, ref.Path, true)
	assertReconcileBranchState(t, repoDir, ref.Branch, true)
	if _, err := os.Stat(recordPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("refused discard record stat error = %v, want not exist", err)
	}
}

func TestReconcileDiscardKeepsSurfaceWhenRecordCannotBeWritten(t *testing.T) {
	t.Parallel()
	homeDir, repoDir, location := newReconcileWorkspace(t)
	run, ref := createReconcileRun(t, homeDir, repoDir, location, "ma/widget-flow", store.StateStopped)
	commitReconcileDispositionChange(t, ref.Path, "work.txt", "work\n", "feat: add work")
	gitImplement(t, repoDir, "merge", "--ff-only", ref.Branch)
	artifactRoot := builtinArtifactDirForRepo(t, repoDir)
	if err := os.MkdirAll(artifactRoot, 0o755); err != nil {
		t.Fatalf("create Artifact Root: %v", err)
	}
	if err := os.WriteFile(filepath.Join(artifactRoot, "runs"), []byte("not a directory\n"), 0o644); err != nil {
		t.Fatalf("block branch record directory: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", run.ID, "--discard-superseded", "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitRunFailed {
		t.Fatalf("record failure exit = %d, want %d stderr=%q stdout=%q", code, exitRunFailed, stderr.String(), stdout.String())
	}
	if !strings.Contains(stderr.String(), "reconciliation operation(s) failed") {
		t.Fatalf("record failure stderr = %q", stderr.String())
	}
	assertReconcilePathState(t, ref.Path, true)
	assertReconcileBranchState(t, repoDir, ref.Branch, true)
}

func TestReconcileWithoutTheFlagRemovesNothing(t *testing.T) {
	t.Parallel()
	homeDir, repoDir, location := newReconcileWorkspace(t)
	run, ref := createReconcileRun(t, homeDir, repoDir, location, "ma/widget-flow", store.StateStopped)
	commitReconcileDispositionChange(t, ref.Path, "work.txt", "work\n", "feat: add work")
	gitImplement(t, repoDir, "merge", "--ff-only", ref.Branch)
	recordPath := filepath.Join(builtinArtifactDirForRepo(t, repoDir), "runs", run.ID, "branch-disposition.json")
	beforeDatabase := readReconcileBytes(t, store.DatabasePath(homeDir))
	beforeGit := reconcileGitSurface(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", run.ID, "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitOK {
		t.Fatalf("report exit = %d, want 0 stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report reconcileReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode report: %v\n%s", err, stdout.String())
	}
	if len(report.Results) != 1 || report.Results[0].Classification != "superseded" ||
		report.Results[0].Action == "discarded" {
		t.Fatalf("read-only disposition report = %+v", report.Results)
	}
	if got := readReconcileBytes(t, store.DatabasePath(homeDir)); !bytes.Equal(got, beforeDatabase) {
		t.Fatal("reconcile without discard flag changed the Run Database")
	}
	if got := reconcileGitSurface(t, repoDir); got != beforeGit {
		t.Fatalf("reconcile without discard flag changed Git state\nbefore:\n%s\nafter:\n%s", beforeGit, got)
	}
	assertReconcilePathState(t, ref.Path, true)
	assertReconcileBranchState(t, repoDir, ref.Branch, true)
	if _, err := os.Stat(recordPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("read-only branch record stat error = %v, want not exist", err)
	}
}

func commitReconcileDispositionChange(t *testing.T, workDir, path, content, message string) string {
	t.Helper()
	mustWrite(t, filepath.Join(workDir, path), content)
	gitImplement(t, workDir, "add", path)
	gitImplement(t, workDir, "commit", "-m", message)
	return strings.TrimSpace(gitImplementOutput(t, workDir, "rev-parse", "HEAD"))
}
