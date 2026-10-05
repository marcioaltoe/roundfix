package cli

// Suite: merged Spec delivery evidence at reconcile and delivery-owner boundaries.
// Invariant: an earlier Run's inherited operator commits are superseded by its
// Spec's delivery even when later default-branch work edits and renames files.
// Boundary IN: public command runner, delivery workflow, local Git and Run Store.
// Boundary OUT: GitHub and network.

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func newReconcileMergeEvidenceFixture(t *testing.T, earlier bool) reconcileMergedFixture {
	t.Helper()
	f := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug, taskStatus: "completed", includeQA: true,
	})
	// Model the item branch's operator commit inherited before the Run's Task.
	gitImplement(t, f.ref.Path, "reset", "--hard", f.run.HeadSHA)
	mustWrite(t, filepath.Join(f.ref.Path, "operator.txt"), "operator work\n")
	mustWrite(t, filepath.Join(f.ref.Path, "renamed.txt"), "operator file\n")
	gitImplement(t, f.ref.Path, "add", "operator.txt", "renamed.txt")
	gitImplement(t, f.ref.Path, "commit", "-m", "fix: operator work inherited from item branch")
	f.taskCommit = commitReconcileMergedTask(t, f.ref.Path, reconcileMergedSpecSlug, reconcileMergedTaskID)
	root := filepath.Join(f.repoDir, "docs", "history", "specs", reconcileMergedSpecSlug)
	mustWrite(t, filepath.Join(root, "_prd.md"), "archived planning\n")
	content := "operator work\n"
	if earlier {
		content = "candidate operator work differs\n"
	}
	mustWrite(t, filepath.Join(f.repoDir, "operator.txt"), content)
	mustWrite(t, filepath.Join(f.repoDir, "renamed.txt"), "operator file\n")
	mustWrite(t, filepath.Join(f.repoDir, "feature.txt"), "delivered\n")
	gitImplement(t, f.repoDir, "add", ".")
	gitImplement(t, f.repoDir, "commit", "-m", "feat: deliver archived Spec")
	f.mergedHead = strings.TrimSpace(gitImplementOutput(t, f.repoDir, "rev-parse", "HEAD"))
	if !earlier {
		mustWrite(t, filepath.Join(f.repoDir, "operator.txt"), "later default edit\n")
		gitImplement(t, f.repoDir, "mv", "renamed.txt", "later-name.txt")
		gitImplement(t, f.repoDir, "add", "operator.txt")
		gitImplement(t, f.repoDir, "commit", "-m", "fix: later Spec edits and renames inherited files")
	}
	return f
}

func TestReconcileReleasesAMergedSpecRunWhoseFilesDivergedLater(t *testing.T) {
	t.Parallel()
	f := newReconcileMergeEvidenceFixture(t, false)
	// No Delivery Queue merge record exists in this measured fixture.
	dry := reconcileMergedResult(t, runReconcileMergedCommand(t, f), f.run.ID)
	if dry.Classification != "superseded" || dry.Action != "would release with --apply" || !strings.Contains(dry.Evidence, "superseded by the delivery of Spec") {
		t.Fatalf("dry-run = %+v", dry)
	}
	assertReconcilePathState(t, f.ref.Path, true)
	assertReconcileBranchState(t, f.repoDir, f.ref.Branch, true)
	applied := reconcileMergedResult(t, runReconcileMergedCommand(t, f, "--apply"), f.run.ID)
	if applied.Classification != "superseded" || applied.Action != "released" {
		t.Fatalf("apply = %+v", applied)
	}
	assertReconcilePathState(t, f.ref.Path, false)
	assertReconcileBranchState(t, f.repoDir, f.ref.Branch, false)
}

func TestDeliverReleaseReleasesAnEarlierRunOfTheMergedSpec(t *testing.T) {
	t.Parallel()
	f := newReconcileMergeEvidenceFixture(t, true)
	// Neither ancestry nor the inherited files represent this earlier Run at
	// the recorded candidate, which is also the local delivery merge commit.
	ctx := context.Background()
	runStore := openDeliverReleaseStore(t, ctx, f.homeDir)
	workflow := &commandDeliveryWorkflow{store: runStore}
	if err := workflow.ReleaseMergedRuns(ctx, f.repoDir, mergedDeliveryItem(f, reconcileMergedSpecSlug)); err != nil {
		t.Fatal(err)
	}
	assertReconcilePathState(t, f.ref.Path, false)
	assertReconcileBranchState(t, f.repoDir, f.ref.Branch, false)
}
