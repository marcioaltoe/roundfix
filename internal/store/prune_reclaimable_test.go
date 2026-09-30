// Suite: reclaimable Run Event Journal pruning
// Invariant: retention reports and locks only Runs whose eligible journals still contain rows.
// Boundary IN: the real temporary Run Database and its machine-wide write lock
// Boundary OUT: artifact directory reclamation, owned by internal/cli/gc_reclaimable_test.go
package store

import (
	"context"
	"slices"
	"testing"
	"time"
)

func TestPruneTerminalRunsSkipsCandidatesWithoutEvents(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	cutoff := time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC)
	withEvents := createPruneCandidate(t, ctx, runStore, cutoff.Add(-2*time.Hour), "with-events", 1)
	empty := createPruneCandidate(t, ctx, runStore, cutoff.Add(-time.Hour), "empty", 0)

	result, err := runStore.PruneTerminalRuns(ctx, cutoff)
	if err != nil {
		t.Fatalf("prune terminal Runs: %v", err)
	}
	if !slices.Equal(result.EligibleRunIDs, []string{withEvents.ID, empty.ID}) {
		t.Fatalf("eligible Run IDs = %v, want [%s %s]", result.EligibleRunIDs, withEvents.ID, empty.ID)
	}
	if !slices.Equal(result.RunIDs, []string{withEvents.ID}) {
		t.Fatalf("pruned Run IDs = %v, want [%s]", result.RunIDs, withEvents.ID)
	}
	if result.Events != 1 {
		t.Fatalf("pruned Run Events = %d, want 1", result.Events)
	}
	if got := countRunEvents(t, ctx, runStore, withEvents.ID); got != 0 {
		t.Fatalf("Run %s retained %d events, want 0", withEvents.ID, got)
	}
}

func TestPruneTerminalRunsNeedsNoWriteLockWhenNoCandidateHasEvents(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)

	cutoff := time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC)
	empty := createPruneCandidate(t, ctx, runStore, cutoff.Add(-time.Hour), "empty", 0)

	holder, err := openWriteLockFile(DatabasePath(homeDir))
	if err != nil {
		t.Fatalf("open independent write lock: %v", err)
	}
	defer func() {
		_ = releaseWriteLock(holder)
		_ = holder.Close()
	}()
	if err := acquireWriteLock(ctx, holder); err != nil {
		t.Fatalf("acquire independent write lock: %v", err)
	}

	pruneCtx, cancel := context.WithTimeout(ctx, 250*time.Millisecond)
	defer cancel()
	result, err := runStore.PruneTerminalRuns(pruneCtx, cutoff)
	if err != nil {
		t.Fatalf("prune empty terminal Run while write lock is held: %v", err)
	}
	if !slices.Equal(result.EligibleRunIDs, []string{empty.ID}) {
		t.Fatalf("eligible Run IDs = %v, want [%s]", result.EligibleRunIDs, empty.ID)
	}
	if len(result.RunIDs) != 0 || result.Events != 0 {
		t.Fatalf("empty prune result = %#v, want no reclaimed Run storage", result)
	}
}

func createPruneCandidate(
	t *testing.T,
	ctx context.Context,
	runStore *Store,
	completedAt time.Time,
	branch string,
	eventCount int,
) Run {
	t.Helper()
	runStore.now = func() time.Time { return completedAt.Add(-time.Hour) }
	request := sampleCreateRunRequest()
	request.HeadBranch = "feature/" + branch
	request.PRNumber = branch
	run, err := runStore.CreateRun(ctx, request)
	if err != nil {
		t.Fatalf("create %s Run: %v", branch, err)
	}
	for range eventCount {
		if _, err := runStore.AppendRunEvent(ctx, sampleRunEvent(run.ID, branch)); err != nil {
			t.Fatalf("append %s Run Event: %v", branch, err)
		}
	}
	runStore.now = func() time.Time { return completedAt }
	completed, err := runStore.CompleteRun(ctx, run.ID, StateClean)
	if err != nil {
		t.Fatalf("complete %s Run: %v", branch, err)
	}
	return completed.Run
}
