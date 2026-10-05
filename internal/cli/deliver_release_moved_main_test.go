// Suite: post-merge cleanup after the default branch moves.
// Invariant: a squash merge must contain exactly the candidate merged into its first parent before Runs are released.
// Boundary IN: ReleaseMergedRuns, disposable Git repositories, Run metadata, and cleanup.
// Boundary OUT: GitHub and Delivery Engine merge sequencing.
package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/preflight"
	"roundfix/internal/store"
)

func TestReleaseMergedRunsAcceptsASquashMergeOntoAMovedDefaultBranch(t *testing.T) {
	t.Parallel()
	fixture, item := newDeliverReleaseMovedMainFixture(t, false)
	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, fixture.homeDir)
	workflow := &commandDeliveryWorkflow{store: runStore}

	if err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, item); err != nil {
		t.Fatalf("release squash merge onto moved default branch: %v", err)
	}
	assertReconcilePathState(t, fixture.ref.Path, false)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, false)
	if _, found, err := runStore.Run(ctx, fixture.run.ID); err != nil || !found {
		t.Fatalf("released Run record: found=%v err=%v", found, err)
	}
}

func TestReleaseMergedRunsRefusesAMergeCommitThatAltersTheCandidate(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		conflict bool
	}{
		{name: "different tree"},
		{name: "conflict resolved differently", conflict: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fixture, item := newDeliverReleaseMovedMainFixture(t, tc.conflict)
			if !tc.conflict {
				mustWrite(t, filepath.Join(fixture.repoDir, "default-only.txt"), "altered merge\n")
				gitImplement(t, fixture.repoDir, "add", "default-only.txt")
				gitImplement(t, fixture.repoDir, "commit", "--amend", "--no-edit")
				item.MergeCommit = strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD"))
			}
			ctx := context.Background()
			runStore := openDeliverReleaseStore(t, ctx, fixture.homeDir)
			workflow := &commandDeliveryWorkflow{store: runStore}

			err := workflow.ReleaseMergedRuns(ctx, fixture.repoDir, item)
			want := fmt.Sprintf("release merged Spec %q Runs: candidate head %q is not represented by merge commit %q",
				item.SpecSlug, item.CandidateCommits[len(item.CandidateCommits)-1], item.MergeCommit)
			if err == nil || err.Error() != want {
				t.Fatalf("altered merge error = %v, want %s", err, want)
			}
			assertReconcilePathState(t, fixture.ref.Path, true)
			assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
		})
	}
}

func newDeliverReleaseMovedMainFixture(t *testing.T, conflict bool) (reconcileMergedFixture, store.DeliveryQueueItem) {
	t.Helper()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	base := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", fixture.mergedHead+"^1"))
	candidate := fixture.mergedHead
	if conflict {
		mustWrite(t, filepath.Join(fixture.repoDir, "shared.txt"), "candidate\n")
		gitImplement(t, fixture.repoDir, "add", "shared.txt")
		gitImplement(t, fixture.repoDir, "commit", "-m", "test: change shared candidate file")
		candidate = strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD"))
	}
	gitImplement(t, fixture.repoDir, "reset", "--hard", base)
	mustWrite(t, filepath.Join(fixture.repoDir, "default-only.txt"), "later default branch commit\n")
	gitImplement(t, fixture.repoDir, "add", "default-only.txt")
	if conflict {
		mustWrite(t, filepath.Join(fixture.repoDir, "shared.txt"), "default branch\n")
		gitImplement(t, fixture.repoDir, "add", "shared.txt")
	}
	gitImplement(t, fixture.repoDir, "commit", "-m", "test: advance default branch during delivery")
	_, mergeErr := (preflight.ExecGitRunner{}).RunGit(context.Background(), fixture.repoDir, "merge", "--squash", candidate)
	if conflict {
		if mergeErr == nil {
			t.Fatal("shared-file squash merge unexpectedly did not conflict")
		}
		unmerged := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "diff", "--name-only", "--diff-filter=U"))
		if unmerged != "shared.txt" {
			t.Fatalf("conflicted paths = %q, want shared.txt", unmerged)
		}
		mustWrite(t, filepath.Join(fixture.repoDir, "shared.txt"), "different resolution\n")
		gitImplement(t, fixture.repoDir, "add", "shared.txt")
	} else if mergeErr != nil {
		t.Fatalf("squash candidate onto moved default branch: %v", mergeErr)
	}
	gitImplement(t, fixture.repoDir, "commit", "-m", "test: squash delivered candidate")
	merge := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", "HEAD"))
	if _, err := (preflight.ExecGitRunner{}).RunGit(context.Background(), fixture.repoDir, "merge-base", "--is-ancestor", candidate, merge); err == nil {
		t.Fatal("candidate unexpectedly is an ancestor of squash merge")
	}
	candidateTree := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", candidate+"^{tree}"))
	mergeTree := strings.TrimSpace(gitImplementOutput(t, fixture.repoDir, "rev-parse", merge+"^{tree}"))
	if candidateTree == mergeTree {
		t.Fatal("moved default branch unexpectedly has the candidate tree")
	}
	item := mergedDeliveryItem(fixture, reconcileMergedSpecSlug)
	item.CandidateCommits = []string{candidate}
	item.MergeCommit = merge
	return fixture, item
}
