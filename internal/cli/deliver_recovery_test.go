// Suite: Delivery Queue item recovery.
// Invariant: retry carries only proved completed Tasks from every eligible Implement Run before the next Implement executor starts.
// Boundary IN: commandDeliveryWorkflow, the real Run Database, and real local Git worktrees.
// Boundary OUT: the Implement executor and publication services, represented by the delivery engine's existing fake boundaries.
package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestItemRecoveryCarriesAnUnresolvedRunIntoTheItemWorktree(t *testing.T) {
	t.Parallel()
	fixture := newIntegrationOrderCarryForwardFixture(t, []implementSeed{
		{id: "task_01", title: "Finish the first Task"},
		{id: "task_02", title: "Keep the failed Task pending", needs: []string{"task_01"}},
	}, []string{"task_01"})
	workflow := newItemRecoveryWorkflow(t, fixture)

	result, err := workflow.CarryForward(t.Context(), fixture.repoDir, implementTestSlug, fixture.run.LocalBranch, fixture.run.ID)
	if err != nil {
		t.Fatalf("carry unresolved Run into item worktree: %v", err)
	}
	if result.RunID != fixture.run.ID || !slices.EqualFunc(result.Runs, []delivery.CarriedRun{{RunID: fixture.run.ID, Carried: []string{"task_01"}}}, equalCarriedRun) {
		t.Fatalf("carry-forward result = %+v, want Run %q and task_01", result, fixture.run.ID)
	}
	completed := mustRead(t, implementTaskPath(fixture.repoDir, "task_01"))
	for _, want := range []string{"status: completed", fixture.run.ID, fixture.commits["task_01"]} {
		if !strings.Contains(completed, want) {
			t.Fatalf("carried task_01 does not contain %q:\n%s", want, completed)
		}
	}
	if pending := mustRead(t, implementTaskPath(fixture.repoDir, "task_02")); !strings.Contains(pending, "status: pending") {
		t.Fatalf("failed task_02 did not stay pending:\n%s", pending)
	}
}

func TestItemRecoveryFindsTheRunByItemBranchWhenNoneIsRecorded(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01"}})
	workflow := newItemRecoveryWorkflow(t, fixture)

	result, err := workflow.CarryForward(t.Context(), fixture.repoDir, implementTestSlug, fixture.run.LocalBranch, "")
	if err != nil {
		t.Fatalf("find and carry Run by item branch: %v", err)
	}
	if result.RunID != fixture.run.ID || !slices.EqualFunc(result.Runs, []delivery.CarriedRun{{RunID: fixture.run.ID, Carried: []string{"task_01"}}}, equalCarriedRun) {
		t.Fatalf("branch-selected carry-forward result = %+v, want Run %q and task_01", result, fixture.run.ID)
	}
}

func TestItemRecoveryCarriesNothingWithoutAnImplementRun(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	workflow := newItemRecoveryWorkflowForRepository(t, homeDir, repoDir)
	beforeHead := itemRecoveryHead(t, repoDir)

	result, err := workflow.CarryForward(t.Context(), repoDir, implementTestSlug, "ma/widget-flow", "")
	if err != nil {
		t.Fatalf("carry without an Implement Run: %v", err)
	}
	if result.RunID != "" || len(result.Runs) != 0 {
		t.Fatalf("carry without an Implement Run = %+v, want empty", result)
	}
	if got := itemRecoveryHead(t, repoDir); got != beforeHead {
		t.Fatalf("head after empty carry = %s, want unchanged %s", got, beforeHead)
	}
}

func TestItemRecoveryCarriesNothingForARunAlreadyCarried(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01"}})
	workflow := newItemRecoveryWorkflow(t, fixture)
	if _, err := workflow.CarryForward(t.Context(), fixture.repoDir, implementTestSlug, fixture.run.LocalBranch, fixture.run.ID); err != nil {
		t.Fatalf("first carry-forward: %v", err)
	}
	beforeHead := itemRecoveryHead(t, fixture.repoDir)

	result, err := workflow.CarryForward(t.Context(), fixture.repoDir, implementTestSlug, fixture.run.LocalBranch, fixture.run.ID)
	if err != nil {
		t.Fatalf("repeat carry-forward: %v", err)
	}
	if result.RunID != fixture.run.ID || len(result.Runs) != 0 {
		t.Fatalf("repeat carry-forward = %+v, want named Run with no Tasks", result)
	}
	if got := itemRecoveryHead(t, fixture.repoDir); got != beforeHead {
		t.Fatalf("head after repeat carry = %s, want unchanged %s", got, beforeHead)
	}
}

func TestItemRecoveryRefusesACarryForwardWithAMovedInput(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01"}})
	workflow := newItemRecoveryWorkflow(t, fixture)
	prdPath := filepath.Join(fixture.repoDir, "docs", "specs", implementTestSlug, "_prd.md")
	mustWrite(t, prdPath, mustRead(t, prdPath)+"\nMoved on the item branch.\n")
	beforeHead := itemRecoveryHead(t, fixture.repoDir)

	_, err := workflow.CarryForward(t.Context(), fixture.repoDir, implementTestSlug, fixture.run.LocalBranch, fixture.run.ID)
	assertItemRecoveryRefusal(t, err, fixture.run.ID, fixture.repoDir)
	if !strings.Contains(err.Error(), filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "_prd.md"))) {
		t.Fatalf("moved-input refusal = %v, want moved PRD path", err)
	}
	if got := itemRecoveryHead(t, fixture.repoDir); got != beforeHead {
		t.Fatalf("head after moved-input refusal = %s, want unchanged %s", got, beforeHead)
	}
}

func TestItemRecoveryRefusesAGoneRunWorktree(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01"}})
	workflow := newItemRecoveryWorkflow(t, fixture)
	gitImplement(t, fixture.repoDir, "worktree", "remove", "--force", fixture.ref.Path)
	beforeHead := itemRecoveryHead(t, fixture.repoDir)

	_, err := workflow.CarryForward(t.Context(), fixture.repoDir, implementTestSlug, fixture.run.LocalBranch, fixture.run.ID)
	assertItemRecoveryRefusal(t, err, fixture.run.ID, fixture.repoDir)
	if !strings.Contains(err.Error(), "Worktree is gone") {
		t.Fatalf("gone-worktree refusal = %v, want gone Worktree reason", err)
	}
	if got := itemRecoveryHead(t, fixture.repoDir); got != beforeHead {
		t.Fatalf("head after gone-worktree refusal = %s, want unchanged %s", got, beforeHead)
	}
}

func TestInspectItemListsUnfinishedTasks(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", status: string(spec.StatusCompleted)},
		{id: "task_02", status: string(spec.StatusPending), needs: []string{"task_01"}},
		{id: "task_03", status: string(spec.StatusFailed), needs: []string{"task_02"}},
	})
	workflow := newItemRecoveryWorkflowForRepository(t, homeDir, repoDir)

	state, err := workflow.InspectItem(t.Context(), repoDir, implementTestSlug)
	if err != nil {
		t.Fatalf("inspect active item: %v", err)
	}
	if state.Archived || state.Head != itemRecoveryHead(t, repoDir) || !slices.Equal(state.UnfinishedTasks, []string{"task_02", "task_03"}) {
		t.Fatalf("active item state = %+v, want task_02 and task_03 unfinished at current head", state)
	}
}

func TestInspectItemReportsAnArchivedSpec(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	workflow := newItemRecoveryWorkflowForRepository(t, homeDir, repoDir)
	source, destination, err := workflow.archivePaths(repoDir, implementTestSlug)
	if err != nil {
		t.Fatalf("resolve fixture archive paths: %v", err)
	}
	destinationPath := filepath.Join(repoDir, filepath.FromSlash(destination))
	mustMkdir(t, filepath.Dir(destinationPath))
	if err := os.Rename(filepath.Join(repoDir, filepath.FromSlash(source)), destinationPath); err != nil {
		t.Fatalf("move fixture Spec to archive: %v", err)
	}

	state, err := workflow.InspectItem(t.Context(), repoDir, implementTestSlug)
	if err != nil {
		t.Fatalf("inspect archived item: %v", err)
	}
	if !state.Archived || state.Head != itemRecoveryHead(t, repoDir) || len(state.UnfinishedTasks) != 0 {
		t.Fatalf("archived item state = %+v, want archived at current head", state)
	}
}

func TestRetriedItemRunsWithItsCompletedTasksCarried(t *testing.T) {
	t.Parallel()
	fixture := newIntegrationOrderCarryForwardFixture(t, []implementSeed{
		{id: "task_01"},
		{id: "task_02", needs: []string{"task_01"}},
	}, []string{"task_01"})
	workflow := newItemRecoveryWorkflow(t, fixture)
	queue, err := workflow.store.CreateDeliveryQueue(t.Context(), fixture.repoDir, []string{implementTestSlug})
	if err != nil {
		t.Fatalf("create retry Delivery Queue: %v", err)
	}
	if _, _, _, err := workflow.store.RecordDeliveryQueueItemWorktree(
		t.Context(),
		fixture.repoDir,
		implementTestSlug,
		fixture.run.LocalBranch,
		fixture.repoDir,
	); err != nil {
		t.Fatalf("record retry item worktree: %v", err)
	}
	if err := workflow.store.SetDeliveryQueueItemWorktreeProvisioned(t.Context(), fixture.repoDir, implementTestSlug, true); err != nil {
		t.Fatalf("record retry item provisioning: %v", err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageParked
	item.Blocker = delivery.BlockerRunUnresolved
	item.Branch = fixture.run.LocalBranch
	item.RunID = fixture.run.ID
	if err := workflow.store.UpdateDeliveryQueueItem(t.Context(), fixture.repoDir, item); err != nil {
		t.Fatalf("park retry Delivery Queue item: %v", err)
	}
	flow := &recoveryInspectingDeliveryFlow{parkTestDeliveryFlow: &parkTestDeliveryFlow{parkSlug: implementTestSlug}}
	engine := delivery.NewEngine(workflow.store, delivery.EngineDependencies{
		Workspace:    flow,
		Runner:       flow,
		Reviewer:     flow,
		Archiver:     flow,
		Gate:         flow,
		Authorizer:   flow,
		Publication:  flow,
		PullRequests: flow,
		Recovery:     workflow,
		Revalidator:  flow,
	})

	retried, err := engine.Retry(t.Context(), fixture.repoDir, implementTestSlug)
	if err != nil {
		t.Fatalf("retry parked item: %v", err)
	}
	if retried.Stage != store.DeliveryStageRunning || !slices.EqualFunc(retried.CarriedFrom.Runs, []delivery.CarriedRun{{RunID: fixture.run.ID, Carried: []string{"task_01"}}}, equalCarriedRun) {
		t.Fatalf("retry result = %+v, want task_01 carried into running", retried)
	}
	if _, err := engine.Run(t.Context(), fixture.repoDir); err != nil {
		t.Fatalf("run retried item: %v", err)
	}
	if flow.runWorkDir != fixture.repoDir || flow.statuses["task_01"] != spec.StatusCompleted || flow.statuses["task_02"] != spec.StatusPending {
		t.Fatalf("Implement executor observed workdir=%q statuses=%v, want item worktree with task_01 completed and task_02 pending", flow.runWorkDir, flow.statuses)
	}
}

func TestItemRecoveryRefusesAnUnknownRecordedRun(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	workflow := newItemRecoveryWorkflowForRepository(t, homeDir, repoDir)

	_, err := workflow.CarryForward(t.Context(), repoDir, implementTestSlug, "ma/widget-flow", "run_missing")
	if err == nil || !strings.Contains(err.Error(), "run_missing") {
		t.Fatalf("unknown recorded Run error = %v, want Run ID", err)
	}
}

func TestItemRecoveryCarriesNothingForAnUnacceptedRunOutcome(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateClean, []implementSeed{{id: "task_01"}})
	workflow := newItemRecoveryWorkflow(t, fixture)
	beforeHead := itemRecoveryHead(t, fixture.repoDir)

	result, err := workflow.CarryForward(t.Context(), fixture.repoDir, implementTestSlug, fixture.run.LocalBranch, fixture.run.ID)
	if err != nil {
		t.Fatalf("carry clean Run: %v", err)
	}
	if result.RunID != fixture.run.ID || len(result.Runs) != 0 || itemRecoveryHead(t, fixture.repoDir) != beforeHead {
		t.Fatalf("clean Run carry = %+v at head %s, want no carry at %s", result, itemRecoveryHead(t, fixture.repoDir), beforeHead)
	}
}

func TestItemRecoveryCarriesNothingWithoutCompletedTaskEvidence(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	workflow := newItemRecoveryWorkflowForRepository(t, homeDir, repoDir)
	head := itemRecoveryHead(t, repoDir)
	run, err := workflow.store.CreateRun(t.Context(), store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "ma/widget-flow",
		HeadSHA:     head,
		SpecSlug:    implementTestSlug,
		Agent:       "codex",
	})
	if err != nil {
		t.Fatalf("create evidence-free Run: %v", err)
	}
	completed, err := workflow.store.CompleteRun(t.Context(), run.ID, store.StateUnresolved)
	if err != nil {
		t.Fatalf("complete evidence-free Run: %v", err)
	}

	result, err := workflow.CarryForward(t.Context(), repoDir, implementTestSlug, run.LocalBranch, run.ID)
	if err != nil {
		t.Fatalf("carry evidence-free Run: %v", err)
	}
	if result.RunID != completed.Run.ID || len(result.Runs) != 0 || itemRecoveryHead(t, repoDir) != head {
		t.Fatalf("evidence-free carry = %+v at head %s, want no carry at %s", result, itemRecoveryHead(t, repoDir), head)
	}
}

func TestItemRecoveryRefusesAnExternalSpecsRoot(t *testing.T) {
	t.Parallel()
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01"}})
	_, externalSpecsRoot := newExternalSpecsRoot(t, implementTestSlug, []implementSeed{{id: "task_01"}})
	workflow := newItemRecoveryWorkflow(t, fixture)
	workflow.loaded.Config.Specs.Root = externalSpecsRoot
	beforeHead := itemRecoveryHead(t, fixture.repoDir)

	_, err := workflow.CarryForward(t.Context(), fixture.repoDir, implementTestSlug, fixture.run.LocalBranch, fixture.run.ID)
	assertItemRecoveryRefusal(t, err, fixture.run.ID, fixture.repoDir)
	if !strings.Contains(err.Error(), "external") {
		t.Fatalf("external Specs Root refusal = %v, want external reason", err)
	}
	if got := itemRecoveryHead(t, fixture.repoDir); got != beforeHead {
		t.Fatalf("head after external-root refusal = %s, want unchanged %s", got, beforeHead)
	}
}

func newItemRecoveryWorkflow(t *testing.T, fixture carryForwardFixture) *commandDeliveryWorkflow {
	t.Helper()
	return newItemRecoveryWorkflowForRepository(t, fixture.homeDir, fixture.repoDir)
}

func newItemRecoveryWorkflowForRepository(t *testing.T, homeDir, repoDir string) *commandDeliveryWorkflow {
	t.Helper()
	runStore, err := store.Open(t.Context(), homeDir)
	if err != nil {
		t.Fatalf("open item recovery Run Database: %v", err)
	}
	t.Cleanup(func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close item recovery Run Database: %v", err)
		}
	})
	return &commandDeliveryWorkflow{
		store: runStore,
		loaded: roundconfig.Loaded{
			Config:  roundconfig.Builtin(),
			GitRoot: repoDir,
			HomeDir: homeDir,
		},
		git: preflight.ExecGitRunner{},
	}
}

func assertItemRecoveryRefusal(t *testing.T, err error, runID, workDir string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected carry-forward refusal")
	}
	var actionable interface{ NextAction() string }
	if !errors.As(err, &actionable) {
		t.Fatalf("carry-forward refusal %T has no next action: %v", err, err)
	}
	next := actionable.NextAction()
	for _, want := range []string{
		"roundfix reconcile " + runID + " --carry-forward",
		workDir,
		"roundfix deliver retry " + implementTestSlug,
	} {
		if strings.Contains(err.Error(), "; amended by ") {
			want = strings.ReplaceAll(want, runID, "'"+runID+"'")
			want = strings.ReplaceAll(want, implementTestSlug, "'"+implementTestSlug+"'")
		}
		if !strings.Contains(next, want) {
			t.Fatalf("carry-forward next action = %q, want %q", next, want)
		}
	}
}

func itemRecoveryHead(t *testing.T, repository string) string {
	t.Helper()
	return strings.TrimSpace(gitImplementOutput(t, repository, "rev-parse", "HEAD"))
}

type recoveryInspectingDeliveryFlow struct {
	*parkTestDeliveryFlow
	statuses map[string]spec.Status
}

func (flow *recoveryInspectingDeliveryFlow) RunSpec(_ context.Context, workDir, specSlug string) (delivery.RunResult, error) {
	graph, err := spec.Load(filepath.Join(workDir, "docs", "specs"), specSlug)
	if err != nil {
		return delivery.RunResult{}, fmt.Errorf("load Spec at Implement executor entry: %w", err)
	}
	flow.runWorkDir = workDir
	flow.runCalls++
	flow.statuses = make(map[string]spec.Status, len(graph.Tasks))
	for _, task := range graph.Tasks {
		flow.statuses[task.ID] = task.Status
	}
	return delivery.RunResult{RunID: "run-retried", Outcome: delivery.RunOutcomeUnresolved}, nil
}
