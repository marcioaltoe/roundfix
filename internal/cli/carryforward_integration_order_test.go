// Suite: Task Carry-Forward integration ordering.
// Invariant: reconcile proves, reports, and applies settled Tasks in Run Branch integration order.
// Boundary IN: public reconcile runner, Run Database evidence, and real local Git worktrees.
// Boundary OUT: Task settlement queueing, owned by internal/daemon task integration tests.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func TestCarryForwardProvesTasksInTheOrderTheRunIntegratedThem(t *testing.T) {
	sharedInput := filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "_prd.md"))
	fixture := newIntegrationOrderCarryForwardFixture(
		t,
		[]implementSeed{
			{
				id:    "task_01",
				title: "Update the shared contract",
				settlementFiles: map[string]string{
					sharedInput: "---\nstatus: active\n---\n\n# PRD\n\nShared contract from task_01.\n",
				},
			},
			{id: "task_02", title: "Use the shared contract"},
		},
		[]string{"task_02", "task_01"},
	)
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
		t.Fatalf("integration-order carry-forward exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("integration-order carry-forward stderr = %q, want empty", stderr.String())
	}
	report := decodeIntegrationOrderReconcileReport(t, stdout.Bytes())
	assertCarryForwardOrder(t, report.CarryForwards, []string{"task_02", "task_01"})
	for index, taskID := range []string{"task_02", "task_01"} {
		candidate := report.CarryForwards[index]
		if candidate.Commit != fixture.commits[taskID] || candidate.Action != "carried forward" {
			t.Fatalf("carry-forward candidate %d = %+v, want Task %s commit %s carried", index, candidate, taskID, fixture.commits[taskID])
		}
		if got := mustRead(t, filepath.Join(fixture.repoDir, "src", taskID+".txt")); got != taskID+" settled\n" {
			t.Fatalf("carried %s implementation = %q", taskID, got)
		}
	}
	carriedCommits := strings.Fields(gitImplementOutput(t, fixture.repoDir, "rev-list", "--reverse", beforeHead+"..HEAD"))
	if len(carriedCommits) != 2 {
		t.Fatalf("carried commit count = %d, want 2: %v", len(carriedCommits), carriedCommits)
	}
	for index, wantTaskID := range []string{"task_02", "task_01"} {
		message := gitImplementOutput(t, fixture.repoDir, "show", "-s", "--format=%B", carriedCommits[index])
		if got := gitTrailerValue(message, "Roundfix-Task"); got != wantTaskID {
			t.Fatalf("carried commit %d Task = %q, want %q; message=%q", index, got, wantTaskID, message)
		}
	}
}

func TestCarryForwardStillRefusesAnInputMovedOnTheCheckout(t *testing.T) {
	sharedInput := filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "_prd.md"))
	fixture := newIntegrationOrderCarryForwardFixture(
		t,
		[]implementSeed{
			{
				id:    "task_01",
				title: "Update the shared contract",
				settlementFiles: map[string]string{
					sharedInput: "---\nstatus: active\n---\n\n# PRD\n\nShared contract from task_01.\n",
				},
			},
			{id: "task_02", title: "Use the shared contract"},
		},
		[]string{"task_02", "task_01"},
	)
	sharedPath := filepath.Join(fixture.repoDir, filepath.FromSlash(sharedInput))
	mustWrite(t, sharedPath, "---\nstatus: active\n---\n\n# PRD\n\nChanged on the checkout after the Run.\n")
	gitImplement(t, fixture.repoDir, "add", filepath.FromSlash(sharedInput))
	gitImplement(t, fixture.repoDir, "commit", "-m", "change shared contract after Run")
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
		t.Fatalf("moved-input carry-forward exit = %d, want %d; stderr=%q stdout=%q", code, exitPreflight, stderr.String(), stdout.String())
	}
	if !strings.Contains(stderr.String(), "carry-forward refused the whole set") || !strings.Contains(stderr.String(), sharedInput) {
		t.Fatalf("moved-input stderr = %q, want whole-set refusal naming %s", stderr.String(), sharedInput)
	}
	report := decodeIntegrationOrderReconcileReport(t, stdout.Bytes())
	assertCarryForwardOrder(t, report.CarryForwards, []string{"task_02", "task_01"})
	if got := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD")); got != beforeHead {
		t.Fatalf("moved-input carry-forward HEAD = %s, want unchanged %s", got, beforeHead)
	}
	for _, taskID := range []string{"task_02", "task_01"} {
		if _, err := os.Stat(filepath.Join(fixture.repoDir, "src", taskID+".txt")); !os.IsNotExist(err) {
			t.Fatalf("moved-input %s implementation stat error = %v, want not exist", taskID, err)
		}
	}
}

func TestCarryForwardUnevaluatedTasksFollowIntegrationOrder(t *testing.T) {
	conflictPath := filepath.ToSlash(filepath.Join("src", "shared.txt"))
	fixture := newIntegrationOrderCarryForwardFixture(
		t,
		[]implementSeed{
			{id: "task_01", title: "Finish the graph's first Task"},
			{
				id:              "task_02",
				title:           "Write the shared file",
				settlementFiles: map[string]string{conflictPath: "settled by task_02\n"},
			},
			{id: "task_03", title: "Finish the graph's last Task"},
		},
		[]string{"task_02", "task_03", "task_01"},
	)
	mustMkdir(t, filepath.Join(fixture.repoDir, "src"))
	mustWrite(t, filepath.Join(fixture.repoDir, filepath.FromSlash(conflictPath)), "written on the checkout instead\n")
	gitImplement(t, fixture.repoDir, "add", filepath.FromSlash(conflictPath))
	gitImplement(t, fixture.repoDir, "commit", "-m", "add conflicting checkout content")
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
		t.Fatalf("conflicting carry-forward exit = %d, want %d; stderr=%q stdout=%q", code, exitPreflight, stderr.String(), stdout.String())
	}
	report := decodeIntegrationOrderReconcileReport(t, stdout.Bytes())
	assertCarryForwardOrder(t, report.CarryForwards, []string{"task_02", "task_03", "task_01"})
	for index, taskID := range []string{"task_03", "task_01"} {
		candidate := report.CarryForwards[index+1]
		if candidate.Action != "refuse" || !strings.Contains(candidate.RefusalReason, "Task "+taskID+" was not evaluated") {
			t.Fatalf("unevaluated candidate %d = %+v, want Task %s refusal", index, candidate, taskID)
		}
	}
	task03Index := strings.Index(stderr.String(), "Task task_03 was not evaluated")
	task01Index := strings.Index(stderr.String(), "Task task_01 was not evaluated")
	if task03Index < 0 || task01Index < 0 || task03Index >= task01Index {
		t.Fatalf("unevaluated-task stderr order = %q, want task_03 before task_01", stderr.String())
	}
	if got := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD")); got != beforeHead {
		t.Fatalf("conflicting carry-forward HEAD = %s, want unchanged %s", got, beforeHead)
	}
}

func newIntegrationOrderCarryForwardFixture(
	t *testing.T,
	graphSeeds []implementSeed,
	settlementOrder []string,
) carryForwardFixture {
	t.Helper()
	homeDir, repoDir := newImplementWorkspace(t, graphSeeds)
	ctx := context.Background()
	head := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open integration-order Run Database: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close integration-order Run Database: %v", err)
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
		t.Fatalf("create integration-order Run: %v", err)
	}
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: repoDir,
		Location: t.TempDir(),
		RunID:    run.ID,
		HeadSHA:  head,
	})
	if err != nil {
		t.Fatalf("create integration-order Run Worktree: %v", err)
	}
	ref.Path, err = filepath.EvalSymlinks(ref.Path)
	if err != nil {
		t.Fatalf("resolve integration-order Run Worktree: %v", err)
	}
	run, err = runStore.SetRunWorkDir(ctx, run.ID, ref.Path)
	if err != nil {
		t.Fatalf("record integration-order Run Worktree: %v", err)
	}
	seedsByID := make(map[string]implementSeed, len(graphSeeds))
	for _, seed := range graphSeeds {
		seedsByID[seed.id] = seed
	}
	commits := make(map[string]string, len(settlementOrder))
	for _, taskID := range settlementOrder {
		seed, ok := seedsByID[taskID]
		if !ok {
			t.Fatalf("settlement order names unknown Task %q", taskID)
		}
		implementationPath := filepath.Join("src", taskID+".txt")
		mustMkdir(t, filepath.Join(ref.Path, filepath.Dir(implementationPath)))
		mustWrite(t, filepath.Join(ref.Path, implementationPath), taskID+" settled\n")
		taskPath := implementTaskPath(ref.Path, taskID)
		if err := spec.SetStatus(taskPath, spec.StatusCompleted); err != nil {
			t.Fatalf("settle source %s: %v", taskID, err)
		}
		settlementPaths := []string{implementationPath, filepath.Join("docs", "specs", implementTestSlug, taskID+".md")}
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
			"-m", fmt.Sprintf("feat: settle %s", taskID),
			"-m", fmt.Sprintf("Roundfix-Spec: %s\nRoundfix-Task: %s", implementTestSlug, taskID),
		)
		commits[taskID] = strings.TrimSpace(gitImplementOutput(t, ref.Path, "rev-parse", "HEAD"))
		if _, err := runStore.AppendRunEvents(ctx, []runevent.RunEvent{
			{
				RunID:       run.ID,
				Source:      runevent.SourceDaemon,
				Kind:        runevent.KindDaemonVerification,
				ReviewIssue: taskID,
				Payload:     []byte(fmt.Sprintf(`{"attempt":1,"phase":"verdict","task":%q,"verdict":"passed"}`, taskID)),
			},
			{
				RunID:       run.ID,
				Source:      runevent.SourceDaemon,
				Kind:        runevent.KindDaemonTask,
				ReviewIssue: taskID,
				Payload:     []byte(fmt.Sprintf(`{"task":%q,"phase":"settled","status":"completed"}`, taskID)),
			},
		}); err != nil {
			t.Fatalf("append integration-order evidence for %s: %v", taskID, err)
		}
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, store.StateUnresolved)
	if err != nil {
		t.Fatalf("complete integration-order Run: %v", err)
	}
	return carryForwardFixture{
		homeDir: homeDir,
		repoDir: repoDir,
		run:     completed.Run,
		ref:     ref,
		commits: commits,
	}
}

func decodeIntegrationOrderReconcileReport(t *testing.T, output []byte) reconcileReport {
	t.Helper()
	var report reconcileReport
	if err := json.Unmarshal(output, &report); err != nil {
		t.Fatalf("decode integration-order reconcile report: %v\n%s", err, output)
	}
	return report
}

func assertCarryForwardOrder(t *testing.T, candidates []spec.CarryForward, want []string) {
	t.Helper()
	got := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		got = append(got, candidate.TaskID)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("carry-forward Task order = %v, want %v; candidates=%+v", got, want, candidates)
	}
}
