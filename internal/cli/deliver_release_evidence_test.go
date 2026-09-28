// Suite: automatic release merge evidence and repository scope.
// Invariant: automatic release mutates only Runs in the recorded repository after Git proves the merge receipt.
// Boundary IN: command delivery workflow, real Git repositories, durable Run repository keys, and Run cleanup.
// Boundary OUT: Delivery Engine sequencing, owned by internal/delivery/merged_release_test.go.
package cli

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseMergedRunsRefusesAnUnresolvedMergeCommit(t *testing.T) {
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
	item.MergeCommit = "missing-merge-commit"

	err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, item)
	if err == nil || !strings.Contains(err.Error(), item.MergeCommit) ||
		!strings.Contains(err.Error(), "does not resolve to a commit") {
		t.Fatalf("unresolved merge commit error = %v, want the bad value and refusal reason", err)
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func TestReleaseMergedRunsRefusesAnUnresolvedCandidateHead(t *testing.T) {
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
	item.CandidateCommits = append(item.CandidateCommits, "missing-candidate-head")

	err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, item)
	if err == nil || !strings.Contains(err.Error(), "missing-candidate-head") ||
		!strings.Contains(err.Error(), "does not resolve to a commit") {
		t.Fatalf("unresolved candidate head error = %v, want the bad value and refusal reason", err)
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func TestReleaseMergedRunsRefusesAMergeCommitOffTheDefaultBranch(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	gitImplement(t, fixture.repoDir, "checkout", "-b", "test/off-default-merge")
	mustWrite(t, filepath.Join(fixture.repoDir, "off-default.txt"), "not merged\n")
	gitImplement(t, fixture.repoDir, "add", "off-default.txt")
	gitImplement(t, fixture.repoDir, "commit", "-m", "test: create off-default merge receipt")
	offDefault := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD"))
	gitImplement(t, fixture.repoDir, "checkout", "main")

	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, fixture.homeDir)
	workflow := &commandDeliveryWorkflow{store: runStore}
	item := mergedDeliveryItem(fixture, reconcileMergedSpecSlug)
	item.MergeCommit = offDefault

	err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, item)
	if err == nil || !strings.Contains(err.Error(), offDefault) ||
		!strings.Contains(err.Error(), `default branch "main"`) {
		t.Fatalf("off-default merge commit error = %v, want the bad value and default branch", err)
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func TestReleaseMergedRunsRefusesAnUnrelatedCandidateHead(t *testing.T) {
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
	item.CandidateCommits = append(item.CandidateCommits, fixture.taskCommit)

	err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, item)
	if err == nil || !strings.Contains(err.Error(), fixture.taskCommit) ||
		!strings.Contains(err.Error(), "not represented by merge commit") {
		t.Fatalf("unrelated candidate head error = %v, want candidate and merge-evidence refusal", err)
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func TestReleaseMergedRunsStillReleasesWithValidEvidence(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, fixture.homeDir)
	workflow := &commandDeliveryWorkflow{store: runStore}
	mergedTree := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", fixture.mergedHead+"^{tree}"))
	mergedParent := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", fixture.mergedHead+"^"))
	squashEquivalentHead := strings.TrimSpace(gitImplementOutput(
		t,
		fixture.repoDir,
		"commit-tree",
		mergedTree,
		"-p",
		mergedParent,
		"-m",
		"test: synthesize squash-equivalent candidate",
	))
	item := mergedDeliveryItem(fixture, reconcileMergedSpecSlug)
	item.CandidateCommits = append(item.CandidateCommits, squashEquivalentHead)

	if err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, item); err != nil {
		t.Fatalf("release with valid merge evidence: %v", err)
	}
	assertReconcilePathState(t, fixture.ref.Path, false)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, false)
}

func TestAutomaticReleaseFindsARunStartedFromAnotherWorktree(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	secondCheckout := filepath.Join(t.TempDir(), "second-checkout")
	gitImplement(t, fixture.repoDir, "worktree", "add", "--detach", secondCheckout, fixture.mergedHead)
	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, fixture.homeDir)
	secondRun, secondRef := createDeliverReleaseRun(
		t,
		ctx,
		runStore,
		secondCheckout,
		t.TempDir(),
		reconcileMergedSpecSlug,
		"delivery/second-worktree-run",
		fixture.mergedHead,
		true,
	)
	workflow := &commandDeliveryWorkflow{store: runStore}

	if err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, mergedDeliveryItem(fixture, reconcileMergedSpecSlug)); err != nil {
		t.Fatalf("release merged Spec Runs from every repository worktree: %v", err)
	}
	assertReconcilePathState(t, secondRef.Path, false)
	assertReconcileBranchState(t, fixture.repoDir, secondRef.Branch, false)
	if _, found, err := runStore.Run(ctx, secondRun.ID); err != nil || !found {
		t.Fatalf("released second-worktree Run record: found=%v err=%v", found, err)
	}
}

func TestAutomaticReleaseIgnoresAnotherRepository(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	_, otherRepository, _ := newReconcileWorkspace(t)
	otherHead := strings.TrimSpace(gitImplementOutput(t, otherRepository, "rev-parse", "HEAD"))
	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, fixture.homeDir)
	_, otherRef := createDeliverReleaseRun(
		t,
		ctx,
		runStore,
		otherRepository,
		t.TempDir(),
		reconcileMergedSpecSlug,
		"delivery/other-repository-run",
		otherHead,
		true,
	)
	workflow := &commandDeliveryWorkflow{store: runStore}

	if err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, mergedDeliveryItem(fixture, reconcileMergedSpecSlug)); err != nil {
		t.Fatalf("release merged Spec Runs in selected repository: %v", err)
	}
	assertReconcilePathState(t, otherRef.Path, true)
	assertReconcileBranchState(t, otherRepository, otherRef.Branch, true)
}
