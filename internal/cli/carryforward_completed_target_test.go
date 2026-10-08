// Suite: Task Carry-Forward for Tasks completed on the target.
// Invariant: a settled Task already completed on the target is nothing to carry and never refuses the remaining set.
// Boundary IN: public reconcile and implement runners, Run Database evidence, and real local Git worktrees.
// Boundary OUT: Task settlement and status ownership, covered by internal/daemon tests.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestReconcileCarryForwardSkipsATaskCompletedOnTheCheckout(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{
		{id: "task_01", title: "Build the core"},
		{id: "task_02", title: "Wire the shell", needs: []string{"task_01"}},
	})
	completeCarryForwardTaskOnCheckout(t, fixture.repoDir, "task_01")
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
		t.Fatalf("partial carry-forward exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("partial carry-forward stderr = %q, want empty", stderr.String())
	}
	report := decodeCompletedTargetReconcileReport(t, stdout.Bytes())
	if len(report.CarryForwards) != 2 {
		t.Fatalf("partial carry-forward candidates = %+v, want two", report.CarryForwards)
	}
	assertCompletedTargetCandidate(t, report.CarryForwards[0], fixture, "task_01")
	if candidate := report.CarryForwards[1]; candidate.TaskID != "task_02" || candidate.Action != "carried forward" {
		t.Fatalf("remaining carry-forward candidate = %+v, want task_02 carried forward", candidate)
	}
	if got := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD")); got == beforeHead {
		t.Fatalf("partial carry-forward HEAD = %s, want task_02 fast-forwarded", got)
	}
	if _, err := os.Stat(filepath.Join(fixture.repoDir, "src", "task_01.txt")); !os.IsNotExist(err) {
		t.Fatalf("completed target task_01 implementation stat error = %v, want not carried", err)
	}
	if got := mustRead(t, filepath.Join(fixture.repoDir, "src", "task_02.txt")); got != "task_02 settled\n" {
		t.Fatalf("carried task_02 implementation = %q", got)
	}
}

func TestReconcileCarryForwardOfACarriedRunCarriesNothing(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01", title: "Build the core"}})
	var firstStdout bytes.Buffer
	var firstStderr bytes.Buffer
	if code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", fixture.run.ID, "--carry-forward", "--format=json"},
		&firstStdout,
		&firstStderr,
	); code != exitOK {
		t.Fatalf("first carry-forward exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, firstStderr.String(), firstStdout.String())
	}
	carriedHead := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD"))
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
		t.Fatalf("repeated carry-forward exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("repeated carry-forward stderr = %q, want empty", stderr.String())
	}
	report := decodeCompletedTargetReconcileReport(t, stdout.Bytes())
	if len(report.CarryForwards) != 1 {
		t.Fatalf("repeated carry-forward candidates = %+v, want one", report.CarryForwards)
	}
	assertCompletedTargetCandidate(t, report.CarryForwards[0], fixture, "task_01")
	if got := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD")); got != carriedHead {
		t.Fatalf("repeated carry-forward HEAD = %s, want unchanged %s", got, carriedHead)
	}
}

func TestCarryForwardStillRefusesTheRemainingSetOnAMovedInput(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{
		{id: "task_01", title: "Build the core"},
		{id: "task_02", title: "Wire the shell", needs: []string{"task_01"}},
	})
	completeCarryForwardTaskOnCheckout(t, fixture.repoDir, "task_01")
	secondTaskPath := implementTaskPath(fixture.repoDir, "task_02")
	mustWrite(t, secondTaskPath, mustRead(t, secondTaskPath)+"\nMoved after settlement.\n")
	gitImplement(t, fixture.repoDir, "add", filepath.Join("docs", "specs", implementTestSlug, "task_02.md"))
	gitImplement(t, fixture.repoDir, "commit", "-m", "move task_02 input after settlement")
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
		t.Fatalf("remaining-set refusal exit = %d, want %d; stderr=%q stdout=%q", code, exitPreflight, stderr.String(), stdout.String())
	}
	if !strings.Contains(stderr.String(), "carry-forward refused the whole set") || !strings.Contains(stderr.String(), "task_02") {
		t.Fatalf("remaining-set refusal stderr = %q, want whole-set refusal naming task_02", stderr.String())
	}
	report := decodeCompletedTargetReconcileReport(t, stdout.Bytes())
	if len(report.CarryForwards) != 2 {
		t.Fatalf("remaining-set candidates = %+v, want two", report.CarryForwards)
	}
	assertCompletedTargetCandidate(t, report.CarryForwards[0], fixture, "task_01")
	if candidate := report.CarryForwards[1]; candidate.TaskID != "task_02" || candidate.Action != "refuse" || candidate.RefusalReason == "" {
		t.Fatalf("remaining-set candidate = %+v, want task_02 refusal", candidate)
	}
	if got := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD")); got != beforeHead {
		t.Fatalf("remaining-set refusal HEAD = %s, want unchanged %s", got, beforeHead)
	}
	if _, err := os.Stat(filepath.Join(fixture.repoDir, "src", "task_02.txt")); !os.IsNotExist(err) {
		t.Fatalf("refused task_02 implementation stat error = %v, want not carried", err)
	}
}

func TestCarryForwardCompletedTargetIsNotReportedAsUnstagedAfterConflict(t *testing.T) {
	t.Parallel()
	conflictPath := filepath.ToSlash(filepath.Join("src", "shared.txt"))
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{
		{
			id:              "task_01",
			title:           "Write the shared file",
			settlementFiles: map[string]string{conflictPath: "settled by task_01\n"},
		},
		{id: "task_02", title: "Finish the remaining work"},
	})
	completeCarryForwardTaskOnCheckout(t, fixture.repoDir, "task_02")
	mustMkdir(t, filepath.Join(fixture.repoDir, "src"))
	mustWrite(t, filepath.Join(fixture.repoDir, filepath.FromSlash(conflictPath)), "written on the checkout instead\n")
	gitImplement(t, fixture.repoDir, "add", filepath.FromSlash(conflictPath))
	gitImplement(t, fixture.repoDir, "commit", "-m", "add conflicting target content")
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
		t.Fatalf("conflicting remaining-set exit = %d, want %d; stderr=%q stdout=%q", code, exitPreflight, stderr.String(), stdout.String())
	}
	if strings.Contains(stderr.String(), "task_02 was not evaluated") {
		t.Fatalf("conflicting remaining-set stderr = %q, completed task_02 must not be reported as unstaged", stderr.String())
	}
	report := decodeCompletedTargetReconcileReport(t, stdout.Bytes())
	if len(report.CarryForwards) != 2 {
		t.Fatalf("conflicting remaining-set candidates = %+v, want two", report.CarryForwards)
	}
	assertCompletedTargetCandidate(t, report.CarryForwards[0], fixture, "task_02")
	if candidate := report.CarryForwards[1]; candidate.TaskID != "task_01" || candidate.Action != "refuse" || candidate.RefusalReason == "" {
		t.Fatalf("conflicting remaining-set candidate = %+v, want task_01 refusal", candidate)
	}
	if got := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD")); got != beforeHead {
		t.Fatalf("conflicting remaining-set HEAD = %s, want unchanged %s", got, beforeHead)
	}
}

func TestImplementPreflightNamesOnlyTheTasksLeftToCarry(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{
		{id: "task_01", title: "Build the core"},
		{id: "task_02", title: "Wire the shell", needs: []string{"task_01"}},
	})
	completeCarryForwardTaskOnCheckout(t, fixture.repoDir, "task_01")
	runner := &implementFakeRunner{gitRoot: fixture.repoDir}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"implement", "--spec", implementTestSlug, "--no-input"},
		&stdout,
		&stderr,
	)

	if code != exitPreflight {
		t.Fatalf("partial carry-forward Preflight exit = %d, want %d; stderr=%q stdout=%q", code, exitPreflight, stderr.String(), stdout.String())
	}
	if !strings.Contains(stderr.String(), "Tasks available for Task Carry-Forward: task_02") {
		t.Fatalf("partial carry-forward Preflight stderr = %q, want task_02", stderr.String())
	}
	if strings.Contains(stderr.String(), "Task Carry-Forward: task_01") {
		t.Fatalf("partial carry-forward Preflight stderr = %q, must not name completed task_01", stderr.String())
	}
	if runner.calls != 0 {
		t.Fatalf("partial carry-forward Preflight reached %d Task Agent turn(s), want zero", runner.calls)
	}
}

func TestImplementPreflightIsSilentForARunWithNothingToCarry(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01", title: "Build the core"}})
	writeImplementSpec(t, fixture.repoDir, implementTestSlug, []implementSeed{
		{id: "task_01", title: "Build the core", status: string(spec.StatusCompleted)},
		{id: "task_02", title: "Continue the Spec", needs: []string{"task_01"}},
	})
	gitImplement(t, fixture.repoDir, "add", filepath.Join("docs", "specs", implementTestSlug))
	gitImplement(t, fixture.repoDir, "commit", "-m", "continue the Spec after carried work")
	runner := &implementFakeRunner{
		gitRoot:      fixture.repoDir,
		statusByTask: map[string]spec.Status{"task_02": spec.StatusCompleted},
	}
	withImplementCollaborators(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"implement", "--spec", implementTestSlug, "--no-input"},
		&stdout,
		&stderr,
	)

	if code != exitOK {
		t.Fatalf("nothing-to-carry Preflight exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if strings.Contains(stderr.String(), "not available for Task Carry-Forward") {
		t.Fatalf("nothing-to-carry Preflight stderr = %q, want no not-available note", stderr.String())
	}
	if runner.calls != 1 {
		t.Fatalf("nothing-to-carry Preflight Task Agent turns = %d, want one for task_02", runner.calls)
	}
}

func completeCarryForwardTaskOnCheckout(t *testing.T, repoDir string, taskID string) {
	t.Helper()
	taskPath := implementTaskPath(repoDir, taskID)
	if err := spec.SetStatus(taskPath, spec.StatusCompleted); err != nil {
		t.Fatalf("complete target %s: %v", taskID, err)
	}
	gitImplement(t, repoDir, "add", filepath.Join("docs", "specs", implementTestSlug, taskID+".md"))
	gitImplement(t, repoDir, "commit", "-m", "complete "+taskID+" on target")
}

func decodeCompletedTargetReconcileReport(t *testing.T, output []byte) reconcileReport {
	t.Helper()
	var report reconcileReport
	if err := json.Unmarshal(output, &report); err != nil {
		t.Fatalf("decode completed-target reconcile report: %v\n%s", err, output)
	}
	return report
}

func assertCompletedTargetCandidate(t *testing.T, candidate spec.CarryForward, fixture carryForwardFixture, taskID string) {
	t.Helper()
	if candidate.TaskID != taskID ||
		candidate.RunID != fixture.run.ID ||
		candidate.TaskFile != filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, taskID+".md")) ||
		candidate.Commit != fixture.commits[taskID] ||
		candidate.Action != carryForwardCompletedAction ||
		candidate.RefusalReason != "" {
		t.Fatalf("completed target candidate = %+v, want Task %s recorded as nothing to carry", candidate, taskID)
	}
}
