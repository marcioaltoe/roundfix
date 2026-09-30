// Suite: automatic release after a remote delivery merge.
// Invariant: cleanup proves merge evidence against a freshly fetched delivery-remote default branch before releasing a Run.
// Boundary IN: command delivery workflow, real local Git remotes and clones, Run metadata, and Run cleanup.
// Boundary OUT: Delivery Engine retry sequencing, owned by internal/delivery/merged_release_test.go.
package cli

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

const deliverReleaseRefreshSpecSlug = "0187-a-queue-that-recovers-without-a-supervisor"

func TestReleaseMergedRunsFetchesTheMergeCommitBeforeReleasing(t *testing.T) {
	t.Parallel()
	fixture := newDeliverReleaseRefreshFixture(t)
	assertCommitIsAbsent(t, fixture.queueClone, fixture.item.MergeCommit)

	workflow := &commandDeliveryWorkflow{store: fixture.runStore}
	if err := workflow.ReleaseMergedRuns(context.Background(), fixture.queueClone, fixture.item); err != nil {
		t.Fatalf("release merged Spec Runs after remote refresh: %v", err)
	}

	assertReconcilePathState(t, fixture.runRef.Path, false)
	assertReconcileBranchState(t, fixture.queueClone, fixture.runRef.Branch, false)
}

func TestReleaseMergedRunsUsesConfiguredDeliveryRemote(t *testing.T) {
	t.Parallel()
	fixture := newDeliverReleaseRefreshFixture(t)
	gittest.Run(t, fixture.queueClone, "remote", "rename", "origin", "delivery")
	gittest.Run(t, fixture.queueClone, "remote", "add", "origin", filepath.Join(t.TempDir(), "missing-origin.git"))
	assertCommitIsAbsent(t, fixture.queueClone, fixture.item.MergeCommit)

	workflow := configuredDeliveryWorkflow(fixture.runStore, "delivery")
	if err := workflow.ReleaseMergedRuns(context.Background(), fixture.queueClone, fixture.item); err != nil {
		t.Fatalf("release merged Spec Runs from configured delivery remote: %v", err)
	}

	assertReconcilePathState(t, fixture.runRef.Path, false)
	assertReconcileBranchState(t, fixture.queueClone, fixture.runRef.Branch, false)
}

func TestReleaseMergedRunsReportsAFailedRefresh(t *testing.T) {
	t.Parallel()
	fixture := newDeliverReleaseRefreshFixture(t)
	missingOrigin := filepath.Join(t.TempDir(), "missing-origin.git")
	gittest.Run(t, fixture.queueClone, "remote", "set-url", "origin", missingOrigin)

	workflow := &commandDeliveryWorkflow{store: fixture.runStore}
	err := workflow.ReleaseMergedRuns(context.Background(), fixture.queueClone, fixture.item)
	if err == nil || !strings.Contains(err.Error(), `refresh default branch "main"`) {
		t.Fatalf("failed refresh error = %v, want refresh default branch prefix", err)
	}

	assertReconcilePathState(t, fixture.runRef.Path, true)
	assertReconcileBranchState(t, fixture.queueClone, fixture.runRef.Branch, true)
}

func TestReleaseMergedRunsWithoutARemoteResolvesLocally(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, fixture.homeDir)
	workflow := &commandDeliveryWorkflow{store: runStore}

	if err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, mergedDeliveryItem(fixture, reconcileMergedSpecSlug)); err != nil {
		t.Fatalf("release merged Spec Runs from local evidence: %v", err)
	}

	assertReconcilePathState(t, fixture.ref.Path, false)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, false)
}

type deliverReleaseRefreshFixture struct {
	queueClone string
	runStore   *store.Store
	runRef     runworktree.Ref
	item       store.DeliveryQueueItem
}

func newDeliverReleaseRefreshFixture(t *testing.T) deliverReleaseRefreshFixture {
	t.Helper()
	root := t.TempDir()
	seed := filepath.Join(root, "seed")
	origin := filepath.Join(root, "origin.git")
	queueClone := filepath.Join(root, "queue-clone")
	mergerClone := filepath.Join(root, "merger-clone")
	gittest.InitRepo(t, seed, "--initial-branch=main")
	mustWrite(t, filepath.Join(seed, "README.md"), "seed\n")
	gittest.Run(t, seed, "add", "README.md")
	gittest.Run(t, seed, "commit", "-m", "test: seed delivery repository")
	gittest.Run(t, "", "clone", "--bare", seed, origin)
	gittest.Harden(t, origin)
	gittest.Run(t, "", "clone", origin, queueClone)
	gittest.Harden(t, queueClone)
	gittest.PersistIdentity(t, queueClone)
	gittest.Run(t, "", "clone", origin, mergerClone)
	gittest.Harden(t, mergerClone)
	gittest.PersistIdentity(t, mergerClone)
	var err error
	queueClone, err = filepath.EvalSymlinks(queueClone)
	if err != nil {
		t.Fatalf("resolve delivery queue clone: %v", err)
	}
	mergerClone, err = filepath.EvalSymlinks(mergerClone)
	if err != nil {
		t.Fatalf("resolve delivery merger clone: %v", err)
	}

	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, t.TempDir())
	startHead := strings.TrimSpace(gittest.Run(t, queueClone, "rev-parse", "HEAD"))
	deliveryBranch := "roundfix/deliver-" + deliverReleaseRefreshSpecSlug
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     queueClone,
		LocalBranch: deliveryBranch,
		HeadSHA:     startHead,
		SpecSlug:    deliverReleaseRefreshSpecSlug,
		Agent:       "codex",
	})
	if err != nil {
		t.Fatalf("create delivery refresh Run: %v", err)
	}
	runRef, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: queueClone,
		Location: t.TempDir(),
		RunID:    run.ID,
		HeadSHA:  startHead,
	})
	if err != nil {
		t.Fatalf("create delivery refresh Run Worktree: %v", err)
	}
	runRef.Path, err = filepath.EvalSymlinks(runRef.Path)
	if err != nil {
		t.Fatalf("resolve delivery refresh Run Worktree: %v", err)
	}
	if _, err := runStore.SetRunWorkDir(ctx, run.ID, runRef.Path); err != nil {
		t.Fatalf("record delivery refresh Run Worktree: %v", err)
	}
	taskCommit := commitReconcileMergedTask(t, runRef.Path, deliverReleaseRefreshSpecSlug, reconcileMergedTaskID)
	if _, err := runStore.CompleteRun(ctx, run.ID, store.StateStopped); err != nil {
		t.Fatalf("complete delivery refresh Run: %v", err)
	}

	gittest.Run(t, queueClone, "checkout", "-b", deliveryBranch, startHead)
	gittest.Run(t, queueClone, "cherry-pick", taskCommit)
	candidateHead := commitReconcileMergedSpec(t, queueClone, deliverReleaseRefreshSpecSlug, "completed", true)
	gittest.Run(t, queueClone, "push", "origin", deliveryBranch)

	gittest.Run(t, mergerClone, "fetch", "origin", deliveryBranch)
	gittest.Run(t, mergerClone, "checkout", "main")
	gittest.Run(t, mergerClone, "merge", "--squash", "origin/"+deliveryBranch)
	gittest.Run(t, mergerClone, "commit", "-m", "feat: merge delivered Spec")
	mergeCommit := strings.TrimSpace(gittest.Run(t, mergerClone, "rev-parse", "HEAD"))
	gittest.Run(t, mergerClone, "push", "origin", "main")

	return deliverReleaseRefreshFixture{
		queueClone: queueClone,
		runStore:   runStore,
		runRef:     runRef,
		item: store.DeliveryQueueItem{
			SpecSlug:          deliverReleaseRefreshSpecSlug,
			Branch:            deliveryBranch,
			CandidateCommits:  []string{candidateHead},
			MergeCommit:       mergeCommit,
			PullRequestNumber: "282",
		},
	}
}

func assertCommitIsAbsent(t *testing.T, repository string, commit string) {
	t.Helper()
	if _, err := (preflight.ExecGitRunner{}).RunGit(
		context.Background(),
		repository,
		"rev-parse",
		"--verify",
		"--end-of-options",
		commit+"^{commit}",
	); err == nil {
		t.Fatalf("merge commit %s unexpectedly exists before cleanup refresh", commit)
	}
}

func configuredDeliveryWorkflow(runStore *store.Store, remote string) *commandDeliveryWorkflow {
	return &commandDeliveryWorkflow{
		store: runStore,
		loaded: roundconfig.Loaded{Config: roundconfig.Config{
			Watch: roundconfig.Watch{PushRemote: remote},
		}},
	}
}
