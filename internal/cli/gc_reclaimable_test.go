// Suite: reclaimable Run storage reporting
// Invariant: GC reports a retention-eligible Run only while journal rows or its artifact directory remain.
// Boundary IN: public GC execution, the real temporary Run Database, and temporary artifact directories
// Boundary OUT: Doctor storage notices and sanitation, owned by their dedicated suites
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"roundfix/internal/store"
)

func TestRetentionReclaimableKeepsCandidateOrder(t *testing.T) {
	t.Parallel()
	artifactRoot := t.TempDir()
	artifactRunID := "run_artifact_only"
	mustMkdir(t, filepath.Join(artifactRoot, artifactRunID))
	candidates := []store.PruneCandidate{
		{RunID: "run_events", Events: 2},
		{RunID: "run_empty"},
		{RunID: artifactRunID},
	}

	got, err := retentionReclaimable(candidates, func(runID string) (bool, error) {
		info, err := os.Stat(filepath.Join(artifactRoot, runID))
		if err != nil {
			if os.IsNotExist(err) {
				return false, nil
			}
			return false, err
		}
		return info.IsDir(), nil
	})
	if err != nil {
		t.Fatalf("find reclaimable Runs: %v", err)
	}
	if want := []string{"run_events", artifactRunID}; !slices.Equal(got, want) {
		t.Fatalf("reclaimable Runs = %v, want %v", got, want)
	}
}

func TestRetentionReclaimableReturnsArtifactInspectionError(t *testing.T) {
	t.Parallel()
	wantErr := errors.New("artifact lookup failed")

	_, err := retentionReclaimable([]store.PruneCandidate{{RunID: "run_empty"}}, func(string) (bool, error) {
		return false, wantErr
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("artifact inspection error = %v, want wrapped error %v", err, wantErr)
	}
}

func TestGCSecondRunReportsNothingPruned(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir, repoDir := withCLIWorkspace(t)
	artifactRoot := filepath.Join(homeDir, "artifacts")
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), fmt.Sprintf("defaults:\n  artifact_dir: %q\nstore:\n  journal_retention: 336h\n", artifactRoot))
	now := time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC)
	seedGCFixture(t, ctx, homeDir, artifactRoot, now)
	withGCNow(t, now)

	var firstStdout bytes.Buffer
	var firstStderr bytes.Buffer
	if code := runCLIContext(t, ctx, []string{"gc"}, &firstStdout, &firstStderr); code != exitOK {
		t.Fatalf("first gc exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, firstStderr.String(), firstStdout.String())
	}

	var secondStdout bytes.Buffer
	var secondStderr bytes.Buffer
	if code := runCLIContext(t, ctx, []string{"gc"}, &secondStdout, &secondStderr); code != exitOK {
		t.Fatalf("second gc exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, secondStderr.String(), secondStdout.String())
	}
	if !strings.Contains(secondStdout.String(), "Runs pruned: 0") {
		t.Fatalf("second gc output = %q, want Runs pruned: 0", secondStdout.String())
	}

	var dryRunStdout bytes.Buffer
	var dryRunStderr bytes.Buffer
	if code := runCLIContext(t, ctx, []string{"gc", "--dry-run"}, &dryRunStdout, &dryRunStderr); code != exitOK {
		t.Fatalf("gc --dry-run exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, dryRunStderr.String(), dryRunStdout.String())
	}
	if !strings.Contains(dryRunStdout.String(), "Runs eligible: 0") {
		t.Fatalf("gc --dry-run output = %q, want Runs eligible: 0", dryRunStdout.String())
	}
}

func TestGCReclaimsARunWhoseArtifactDirectoryOutlivedItsEvents(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir, repoDir := withCLIWorkspace(t)
	artifactRoot := filepath.Join(homeDir, "artifacts")
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), fmt.Sprintf("defaults:\n  artifact_dir: %q\nstore:\n  journal_retention: 336h\n", artifactRoot))
	now := time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)

	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run store: %v", err)
	}
	run := createGCTestRun(t, ctx, runStore, artifactRoot, "artifact-only", 0)
	if _, err := runStore.CompleteRun(ctx, run.ID, store.StateClean); err != nil {
		t.Fatalf("complete artifact-only Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run store: %v", err)
	}
	setRunTimestamps(t, homeDir, run.ID, now.Add(-400*time.Hour), now.Add(-400*time.Hour))
	runDir := writeRunArtifact(t, artifactRoot, run.ID, "surviving artifact")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := runCLIContext(t, ctx, []string{"gc"}, &stdout, &stderr); code != exitOK {
		t.Fatalf("gc exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if !strings.Contains(stdout.String(), "Runs pruned: 1") || !strings.Contains(stdout.String(), run.ID) {
		t.Fatalf("gc output = %q, want one pruned Run named %s", stdout.String(), run.ID)
	}
	assertPathMissing(t, runDir)
}

func TestRetentionSweepIsSilentWhenNothingIsReclaimed(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	artifactRoot := t.TempDir()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close Run store: %v", err)
		}
	}()
	run := createGCTestRun(t, ctx, runStore, artifactRoot, "empty-sweep", 0)
	if _, err := runStore.CompleteRun(ctx, run.ID, store.StateClean); err != nil {
		t.Fatalf("complete empty sweep Run: %v", err)
	}
	now := time.Now().UTC()
	setRunTimestamps(t, homeDir, run.ID, now.Add(-400*time.Hour), now.Add(-400*time.Hour))

	var stderr bytes.Buffer
	sweepRunRetention(ctx, runStore, artifactRoot, 336*time.Hour, &stderr)

	if strings.Contains(stderr.String(), "pruned Run storage") {
		t.Fatalf("empty retention sweep reported reclaimed storage: %q", stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("empty retention sweep stderr = %q, want silence", stderr.String())
	}
}
