// Suite: automatic release after a delivery merge.
// Invariant: delivery releases only terminal Runs of the merged Spec that the recorded candidate head proves represented.
// Boundary IN: command delivery workflow, real Git repositories, Run metadata, merged-head inspection, and cleanup.
// Boundary OUT: Delivery Engine sequencing, owned by internal/delivery/merged_release_test.go.
package cli

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func TestDeliverReleaseRemovesEveryProvenRunOfTheMergedSpec(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, fixture.homeDir)
	contained, containedRef := createDeliverReleaseRun(
		t,
		ctx,
		runStore,
		fixture.repoDir,
		filepath.Dir(fixture.ref.Path),
		reconcileMergedSpecSlug,
		"delivery/contained-run",
		fixture.mergedHead,
		true,
	)
	workflow := &commandDeliveryWorkflow{store: runStore}

	err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, mergedDeliveryItem(fixture, reconcileMergedSpecSlug))
	if err != nil {
		t.Fatalf("release merged Spec Runs: %v", err)
	}

	assertReconcilePathState(t, fixture.ref.Path, false)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, false)
	assertReconcilePathState(t, containedRef.Path, false)
	assertReconcileBranchState(t, fixture.repoDir, containedRef.Branch, false)
	for _, runID := range []string{fixture.run.ID, contained.ID} {
		_, found, err := runStore.Run(ctx, runID)
		if err != nil || !found {
			t.Fatalf("released Run %q record: found=%v err=%v", runID, found, err)
		}
	}
}

func TestDeliverReleaseKeepsAnUnrepresentedRun(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "pending",
	})
	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, fixture.homeDir)
	workflow := &commandDeliveryWorkflow{store: runStore}

	err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, mergedDeliveryItem(fixture, reconcileMergedSpecSlug))
	if err == nil || !strings.Contains(err.Error(), fixture.run.ID) ||
		!strings.Contains(err.Error(), fixture.taskCommit[:12]) {
		t.Fatalf("unrepresented Run release error = %v, want Run ID and commit", err)
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func TestDeliverReleaseLeavesAnActiveRunAlone(t *testing.T) {
	t.Parallel()
	homeDir, repoDir, _ := newReconcileWorkspace(t)
	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, homeDir)
	head := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "delivery/active-run",
		HeadSHA:     head,
		SpecSlug:    reconcileMergedSpecSlug,
		Agent:       "codex",
	})
	if err != nil {
		t.Fatalf("create Active Run: %v", err)
	}
	workflow := &commandDeliveryWorkflow{store: runStore}
	item := store.DeliveryQueueItem{
		SpecSlug:          reconcileMergedSpecSlug,
		Branch:            "roundfix/deliver-" + reconcileMergedSpecSlug,
		CandidateCommits:  []string{head},
		MergeCommit:       head,
		PullRequestNumber: "259",
	}

	err = workflow.ReleaseMergedRuns(ctx, repoDir, item)
	if err == nil || !strings.Contains(err.Error(), run.ID) || !strings.Contains(err.Error(), "Active") {
		t.Fatalf("Active Run release error = %v, want Run ID and Active reason", err)
	}
	got, found, readErr := runStore.Run(ctx, run.ID)
	if readErr != nil || !found || got.State != store.StateActive {
		t.Fatalf("Active Run after release: found=%v state=%q err=%v", found, got.State, readErr)
	}
}

func TestDeliverReleaseNeverTouchesAnotherSpecsRun(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, fixture.homeDir)
	workflow := &commandDeliveryWorkflow{store: runStore}

	err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, mergedDeliveryItem(fixture, "0174-another-spec"))
	if err != nil {
		t.Fatalf("release another Spec's Runs: %v", err)
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func TestDeliverReleaseRefusesAnItemWithoutAMergeCommit(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, fixture.homeDir)
	workflow := &commandDeliveryWorkflow{store: runStore}
	item := mergedDeliveryItem(fixture, reconcileMergedSpecSlug)
	item.MergeCommit = ""

	err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, item)
	if err == nil || !strings.Contains(err.Error(), "merge commit is required") {
		t.Fatalf("missing merge commit error = %v", err)
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func TestDeliverReleaseRefusesAnItemWithoutACandidateHead(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, fixture.homeDir)
	workflow := &commandDeliveryWorkflow{store: runStore}
	item := mergedDeliveryItem(fixture, reconcileMergedSpecSlug)
	item.CandidateCommits = nil

	err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, item)
	if err == nil || !strings.Contains(err.Error(), "candidate head is required") {
		t.Fatalf("missing candidate head error = %v", err)
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func openDeliverReleaseStore(t *testing.T, ctx context.Context, homeDir string) *store.Store {
	t.Helper()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open delivery release Run Database: %v", err)
	}
	t.Cleanup(func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close delivery release Run Database: %v", err)
		}
	})
	return runStore
}

func createDeliverReleaseRun(
	t *testing.T,
	ctx context.Context,
	runStore *store.Store,
	repository string,
	location string,
	specSlug string,
	targetBranch string,
	head string,
	terminal bool,
) (store.Run, runworktree.Ref) {
	t.Helper()
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repository,
		LocalBranch: targetBranch,
		HeadSHA:     head,
		SpecSlug:    specSlug,
		Agent:       "codex",
	})
	if err != nil {
		t.Fatalf("create delivery release Run: %v", err)
	}
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: repository,
		Location: location,
		RunID:    run.ID,
		HeadSHA:  head,
	})
	if err != nil {
		t.Fatalf("create delivery release Run Worktree: %v", err)
	}
	ref.Path, err = filepath.EvalSymlinks(ref.Path)
	if err != nil {
		t.Fatalf("resolve delivery release Run Worktree: %v", err)
	}
	run, err = runStore.SetRunWorkDir(ctx, run.ID, ref.Path)
	if err != nil {
		t.Fatalf("record delivery release Run Worktree: %v", err)
	}
	if terminal {
		completed, err := runStore.CompleteRun(ctx, run.ID, store.StateStopped)
		if err != nil {
			t.Fatalf("complete delivery release Run: %v", err)
		}
		run = completed.Run
	}
	return run, ref
}

func mergedDeliveryItem(fixture reconcileMergedFixture, specSlug string) store.DeliveryQueueItem {
	return store.DeliveryQueueItem{
		SpecSlug:          specSlug,
		Branch:            "roundfix/deliver-" + specSlug,
		CandidateCommits:  []string{fixture.run.HeadSHA, fixture.mergedHead},
		MergeCommit:       fixture.mergedHead,
		PullRequestNumber: "259",
	}
}
