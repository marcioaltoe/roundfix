// Suite: Delivery Retry refusals caused by Spec amendments.
// Invariant: only non-Task commits that moved every refused Task input add the quoted five-command recovery.
// Boundary IN: commandDeliveryWorkflow, the real Run Database, and real local Git repositories and worktrees.
// Boundary OUT: command-level stderr rendering, covered by the existing CLI preflight tests.
package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestRetryRefusalNamesTheAmendingCommitAndTheRecovery(t *testing.T) {
	t.Parallel()
	homeDir, originalRepo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	repository := filepath.Join(t.TempDir(), "item worktree's $queue")
	if err := os.Rename(originalRepo, repository); err != nil {
		t.Fatalf("move item worktree under shell-sensitive path: %v", err)
	}
	resolvedRepository, err := filepath.EvalSymlinks(repository)
	if err != nil {
		t.Fatalf("resolve item worktree: %v", err)
	}
	setCommandEnvironmentForTest(t, homeDir, resolvedRepository)
	workflow := newItemRecoveryWorkflowForRepository(t, homeDir, resolvedRepository)
	fixture := createDeliveryRetryRun(
		t,
		workflow,
		resolvedRepository,
		"ma/widget-flow",
		[]implementSeed{{id: "task_01"}},
		store.StateBudgetExceeded,
	)
	prdPath := filepath.Join(resolvedRepository, "docs", "specs", implementTestSlug, "_prd.md")
	amendments := make([]string, 0, 2)
	for _, subject := range []string{"docs: amend retry premise", "docs: clarify retry premise"} {
		mustWrite(t, prdPath, mustRead(t, prdPath)+"\n"+subject+".\n")
		gitImplement(t, resolvedRepository, "add", filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "_prd.md")))
		gitImplement(t, resolvedRepository, "commit", "-m", subject)
		amendments = append(amendments, itemRecoveryHead(t, resolvedRepository))
	}
	beforeHead := itemRecoveryHead(t, resolvedRepository)

	_, carryErr := workflow.CarryForward(
		t.Context(),
		resolvedRepository,
		implementTestSlug,
		fixture.run.LocalBranch,
		fixture.run.ID,
	)
	reason, nextAction := deliveryRetryRefusalTexts(t, carryErr)
	movedInput := filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "_prd.md"))
	wantReason := fmt.Sprintf(
		"carry-forward refused the whole set: Task task_01 declared input(s) moved: %s; amended by %s",
		movedInput,
		strings.Join(amendments, ", "),
	)
	if reason != wantReason {
		t.Fatalf("amendment refusal = %q, want %q", reason, wantReason)
	}
	quotedRepository := "'" + strings.Replace(resolvedRepository, "worktree's", "worktree'\"'\"'s", 1) + "'"
	wantNextAction := strings.Join([]string{
		fmt.Sprintf("git -C %s branch 'roundfix-amended-%s' HEAD", quotedRepository, fixture.run.ID),
		fmt.Sprintf("git -C %s reset --hard '%s^'", quotedRepository, amendments[0]),
		fmt.Sprintf("(cd %s && roundfix reconcile '%s' --carry-forward)", quotedRepository, fixture.run.ID),
		fmt.Sprintf("git -C %s cherry-pick '%s' '%s'", quotedRepository, amendments[0], amendments[1]),
		fmt.Sprintf("roundfix deliver retry '%s'", implementTestSlug),
	}, "\n")
	if nextAction != wantNextAction {
		t.Fatalf("amendment next action = %q, want %q", nextAction, wantNextAction)
	}
	if got := itemRecoveryHead(t, resolvedRepository); got != beforeHead {
		t.Fatalf("item HEAD after amendment refusal = %s, want unchanged %s", got, beforeHead)
	}
}

func TestRetryRefusalByATaskCommitKeepsTodaysText(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateBudgetExceeded, []implementSeed{{id: "task_01"}})
	workflow := newItemRecoveryWorkflow(t, fixture)
	prdPath := filepath.Join(fixture.repoDir, "docs", "specs", implementTestSlug, "_prd.md")
	mustWrite(t, prdPath, mustRead(t, prdPath)+"\nMoved by a Task commit.\n")
	gitImplement(t, fixture.repoDir, "add", filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "_prd.md")))
	gitImplement(
		t,
		fixture.repoDir,
		"commit",
		"-m", "docs: move a declared input",
		"-m", fmt.Sprintf("Roundfix-Spec: %s\nRoundfix-Task: task_02", implementTestSlug),
	)

	_, err := workflow.CarryForward(
		t.Context(),
		fixture.repoDir,
		implementTestSlug,
		fixture.run.LocalBranch,
		fixture.run.ID,
	)
	reason, nextAction := deliveryRetryRefusalTexts(t, err)
	wantReason := fmt.Sprintf(
		"carry-forward refused the whole set: Task task_01 declared input(s) moved: %s",
		filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "_prd.md")),
	)
	if reason != wantReason {
		t.Fatalf("Task-commit refusal = %q, want today's %q", reason, wantReason)
	}
	wantNextAction := todaysDeliveryCarryForwardNextAction(fixture.run.ID, fixture.repoDir)
	if nextAction != wantNextAction {
		t.Fatalf("Task-commit next action = %q, want today's %q", nextAction, wantNextAction)
	}
}

func TestRetryRefusalWithoutMovedInputsKeepsTodaysText(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateBudgetExceeded, []implementSeed{{id: "task_01"}})
	workflow := newItemRecoveryWorkflow(t, fixture)
	gitImplement(t, fixture.repoDir, "worktree", "remove", "--force", fixture.ref.Path)

	_, err := workflow.CarryForward(
		t.Context(),
		fixture.repoDir,
		implementTestSlug,
		fixture.run.LocalBranch,
		fixture.run.ID,
	)
	reason, nextAction := deliveryRetryRefusalTexts(t, err)
	wantReason := fmt.Sprintf("carry-forward Run %q Worktree is gone", fixture.run.ID)
	if reason != wantReason {
		t.Fatalf("non-moved-input refusal = %q, want today's %q", reason, wantReason)
	}
	wantNextAction := todaysDeliveryCarryForwardNextAction(fixture.run.ID, fixture.repoDir)
	if nextAction != wantNextAction {
		t.Fatalf("non-moved-input next action = %q, want today's %q", nextAction, wantNextAction)
	}
}

func deliveryRetryRefusalTexts(t *testing.T, err error) (string, string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected Delivery Retry carry-forward refusal")
	}
	var actionable interface{ NextAction() string }
	if !errors.As(err, &actionable) {
		t.Fatalf("Delivery Retry refusal %T has no next action: %v", err, err)
	}
	return err.Error(), actionable.NextAction()
}

func todaysDeliveryCarryForwardNextAction(runID, workDir string) string {
	return fmt.Sprintf(
		"run `roundfix reconcile %s --carry-forward` in item worktree %q, then run `roundfix deliver retry %s`",
		runID,
		workDir,
		implementTestSlug,
	)
}
