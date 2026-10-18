// Suite: merged Item Binary reclamation.
// Invariant: merged release removes Item Binaries and preserves delivery evidence.
// Boundary IN: ReleaseMergedRuns, temporary repositories, homes, and Run Databases.
// Boundary OUT: merge scheduling, covered by the delivery engine tests.
package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestMergedReleaseReclaimsTheItemBinaries(t *testing.T) {
	t.Parallel()
	f := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{archivedSlug: reconcileMergedSpecSlug, taskStatus: "completed", includeQA: true})
	s := openDeliverReleaseStore(t, t.Context(), f.homeDir)
	artifacts := filepath.Join(f.homeDir, "artifacts")
	itemDir := filepath.Join(artifacts, "delivery", reconcileMergedSpecSlug)
	mustMkdir(t, filepath.Join(itemDir, "item-binary", "digest"))
	mustWrite(t, filepath.Join(itemDir, "item-binary", "digest", "copy"), "binary bytes")
	for _, log := range []string{"item-binary-build.log", "repository-gate.log"} {
		mustWrite(t, filepath.Join(itemDir, log), "retained evidence")
	}
	workflow := &commandDeliveryWorkflow{store: s, loaded: roundconfig.Loaded{GitRoot: f.repoDir, HomeDir: f.homeDir}}
	workflow.loaded.Config.Defaults.ArtifactDir = artifacts
	if err := workflow.ReleaseMergedRuns(t.Context(), f.repoDir, mergedDeliveryItem(f, reconcileMergedSpecSlug)); err != nil {
		t.Fatal(err)
	}
	assertPathMissing(t, filepath.Join(itemDir, "item-binary"))
	for _, log := range []string{"item-binary-build.log", "repository-gate.log"} {
		assertFileContent(t, filepath.Join(itemDir, log), "retained evidence")
	}
	assertReconcilePathState(t, f.ref.Path, false)
}

func TestMergedReleaseJoinsItemBinaryCleanupAndRunErrors(t *testing.T) {
	t.Parallel()
	f := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{archivedSlug: reconcileMergedSpecSlug, taskStatus: "pending"})
	s := openDeliverReleaseStore(t, t.Context(), f.homeDir)
	artifacts := filepath.Join(f.homeDir, "artifacts")
	mustMkdir(t, artifacts)
	// A non-directory ancestor makes artifact cleanup fail on every platform.
	mustWrite(t, filepath.Join(artifacts, "delivery"), "not a directory")
	workflow := &commandDeliveryWorkflow{store: s, loaded: roundconfig.Loaded{GitRoot: f.repoDir, HomeDir: f.homeDir}}
	workflow.loaded.Config.Defaults.ArtifactDir = artifacts
	err := workflow.ReleaseMergedRuns(t.Context(), f.repoDir, mergedDeliveryItem(f, reconcileMergedSpecSlug))
	var pathErr *os.PathError
	if err == nil || !strings.Contains(err.Error(), f.run.ID) || !strings.Contains(err.Error(), "Item Binary") || !errors.As(err, &pathErr) {
		t.Fatalf("expected joined Run and artifact errors, got %v", err)
	}
	assertReconcilePathState(t, f.ref.Path, true)
}
