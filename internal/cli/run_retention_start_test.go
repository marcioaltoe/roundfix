// Suite: automatic Run Retention with fixture homes and public command harnesses.
// Invariant: starts remain available, while protected Runs and read-only reports survive.
// Boundary IN: real SQLite, local artifacts, configuration, and command output.
// Boundary OUT: Agent execution and detached Delivery Queue owner launch.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func retentionStartStore(t *testing.T, home string) (*store.Store, roundconfig.Loaded) {
	t.Helper()
	ctx := commandContextForTest(t, context.Background())
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: home, WorkDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	s, err := store.Open(ctx, home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Error(err)
		}
	})
	return s, loaded
}

func retentionStartCall(t *testing.T, s *store.Store, loaded roundconfig.Loaded) string {
	t.Helper()
	var diagnostics bytes.Buffer
	runRetentionAtStart(commandContextForTest(t, context.Background()), s, loaded, &diagnostics)
	return diagnostics.String()
}

func assertRetentionRun(t *testing.T, s *store.Store, id string, want bool) {
	t.Helper()
	_, found, err := s.Run(t.Context(), id)
	if err != nil || found != want {
		t.Fatalf("Run %s found=%t want=%t err=%v", id, found, want, err)
	}
}

func seedOldStartRun(t *testing.T, home string, now time.Time, name string) store.Run {
	t.Helper()
	s, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	run := createGCTestRun(t, t.Context(), s, filepath.Join(home, ".roundfix", "retention-artifacts"), name, 1)
	if _, err := s.CompleteRun(t.Context(), run.ID, store.StateClean); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	setRunTimestamps(t, home, run.ID, now.Add(-40*24*time.Hour), now.Add(-31*24*time.Hour))
	return run
}

func TestRunRetentionRunsOnceADayAtRunStart(t *testing.T) {
	for _, trigger := range []string{"24 hours", "window change"} {
		t.Run(trigger, func(t *testing.T) {
			f := seedGCRetention(t, false, 2)
			s, loaded := retentionStartStore(t, f.home)
			counts := gcRetentionRows(t, f.home, gcRetentionTables, f.removed.ID)
			var rows int
			for _, values := range counts {
				rows += len(values)
			}
			before, _, _, err := s.RunRetentionStorage(t.Context())
			if err != nil {
				t.Fatal(err)
			}
			output := retentionStartCall(t, s, loaded)
			after, _, _, err := s.RunRetentionStorage(t.Context())
			if err != nil {
				t.Fatal(err)
			}
			want := fmt.Sprintf("roundfix: Run Retention removed runs=1 rows=%d database_bytes_reclaimed=%d\n", rows, before-after)
			if output != want {
				t.Fatalf("stderr=%q want=%q", output, want)
			}
			assertRetentionRun(t, s, f.removed.ID, false)
			setRunTimestamps(t, f.home, f.recent.ID, f.now.Add(-40*24*time.Hour), f.now.Add(-31*24*time.Hour))
			withGCNow(t, f.now.Add(23*time.Hour))
			if output := retentionStartCall(t, s, loaded); output != "" {
				t.Fatalf("not due: %q", output)
			}
			assertRetentionRun(t, s, f.recent.ID, true)
			if trigger == "24 hours" {
				withGCNow(t, f.now.Add(24*time.Hour))
			} else {
				loaded.Config.Store.RunRetentionDays = 7
			}
			assertGCRetentionContains(t, retentionStartCall(t, s, loaded), "Run Retention removed runs=1")
			assertRetentionRun(t, s, f.recent.ID, false)
			assertRetentionRun(t, s, f.queued.ID, true)
			assertRetentionRun(t, s, f.worktree.ID, true)
			last, found, err := s.LastRunRetentionSweep(t.Context())
			if err != nil || !found || last.RetentionDays != loaded.Config.Store.RunRetentionDays {
				t.Fatalf("last=%+v found=%t err=%v", last, found, err)
			}
		})
	}
}

func TestRunRetentionStopsOnItsBudgetAndResumes(t *testing.T) {
	home, _ := withCLIWorkspace(t)
	now := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	var runs []store.Run
	for _, name := range []string{"one", "two", "three"} {
		runs = append(runs, seedOldStartRun(t, home, now, name))
	}
	s, loaded := retentionStartStore(t, home)
	calls := 0
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.gc.now = func() time.Time {
			calls++
			if calls > 2 {
				return now.Add(2 * time.Second)
			}
			return now
		}
	})
	output := retentionStartCall(t, s, loaded)
	if output != "roundfix: Run Retention paused after its 2s budget; it continues at the next Run start\n" {
		t.Fatalf("stderr=%q", output)
	}
	removed := 0
	for _, run := range runs {
		_, found, err := s.Run(t.Context(), run.ID)
		if err != nil {
			t.Fatal(err)
		}
		if !found {
			removed++
		}
	}
	if removed != 1 {
		t.Fatalf("removed=%d want=1", removed)
	}
	if _, found, err := s.LastRunRetentionSweep(t.Context()); err != nil || found {
		t.Fatalf("completion found=%t err=%v", found, err)
	}
	withGCNow(t, now.Add(time.Minute))
	assertGCRetentionContains(t, retentionStartCall(t, s, loaded), "Run Retention removed runs=2")
	for _, run := range runs {
		assertRetentionRun(t, s, run.ID, false)
	}
	if _, found, err := s.LastRunRetentionSweep(t.Context()); err != nil || !found {
		t.Fatalf("completion found=%t err=%v", found, err)
	}
}

func TestRunRetentionAtStartWarnsAndNeverBlocks(t *testing.T) {
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	_, _ = retentionStartStore(t, home)
	db := openGCRetentionDB(t, home)
	gcRetentionExec(t, db, "DROP TABLE run_retention_sweeps")
	runner := &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}}
	withImplementCollaborators(t, runner)
	var out, diagnostics bytes.Buffer
	code := runCLI(t, []string{"implement", "--spec", implementTestSlug, "--no-input"}, &out, &diagnostics)
	if code != exitOK || runner.calls != 1 {
		t.Fatalf("exit=%d calls=%d stderr=%s", code, runner.calls, &diagnostics)
	}
	want := "roundfix: warning: Run Retention failed: read last Run Retention sweep: SQL logic error: no such table: run_retention_sweeps (1)\n"
	if strings.Count(diagnostics.String(), want) != 1 {
		t.Fatalf("stderr=%q want warning=%q", diagnostics.String(), want)
	}
}

func TestImplementStartRunsRunRetention(t *testing.T) {
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	now := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)
	old := seedOldStartRun(t, home, now, "implement-old")
	runner := &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}}
	withImplementCollaborators(t, runner)
	// Refuse Run creation if the old Run has not yet left the database.
	db := openGCRetentionDB(t, home)
	gcRetentionExec(t, db, fmt.Sprintf(`CREATE TRIGGER retention_before_implement
 BEFORE INSERT ON runs WHEN NEW.kind = 'implement'
 AND EXISTS (SELECT 1 FROM runs WHERE id = '%s')
 BEGIN SELECT RAISE(ABORT, 'old Run still exists at Run creation'); END`, old.ID))
	var out, diagnostics bytes.Buffer
	if code := runCLI(t, []string{"implement", "--spec", implementTestSlug, "--no-input"}, &out, &diagnostics); code != exitOK || runner.calls != 1 {
		t.Fatalf("exit=%d Agent calls=%d stderr=%s", code, runner.calls, &diagnostics)
	}
	assertGCRetentionContains(t, diagnostics.String(), "Run Retention removed runs=1")
}

func TestDeliverStartRunsRunRetention(t *testing.T) {
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	setImplementFixtureAuthorizationOperations(t, repo, allDeliveryOperations...)
	now := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)
	old := seedOldStartRun(t, home, now, "delivery-old")
	queued := seedOldStartRun(t, home, now, "delivery-protected")
	s, _ := retentionStartStore(t, home)
	linked := seedOldStartRun(t, home, now, "delivery-linked")
	// Seed references as the command records its queue. This exposes an
	// incorrectly ordered sweep before queue creation without a production hook.
	db := openGCRetentionDB(t, home)
	quote := func(value string) string { return "'" + strings.ReplaceAll(value, "'", "''") + "'" }
	gcRetentionExec(t, db, fmt.Sprintf(`CREATE TRIGGER retention_queue_fixture
 AFTER INSERT ON delivery_queue_items WHEN NEW.git_root = %s
 BEGIN
 UPDATE delivery_queue_items SET run_id = %s
 WHERE git_root = NEW.git_root AND spec_slug = NEW.spec_slug;
 INSERT INTO delivery_queue_runs (git_root, spec_slug, run_id)
 VALUES (NEW.git_root, NEW.spec_slug, %s);
 END`, quote(repo), quote(queued.ID), quote(linked.ID)))
	started := false
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.deliveryReadiness = readyDeliveryReadiness
		deps.startDeliveryOwner = func(ctx context.Context, _ roundconfig.Loaded, _ commandEnvironment, _, _ io.Writer) int {
			assertRetentionRun(t, s, old.ID, false)
			assertRetentionRun(t, s, queued.ID, true)
			assertRetentionRun(t, s, linked.ID, true)
			if _, found, err := s.DeliveryQueue(ctx, repo); err != nil || !found {
				t.Fatalf("new queue found=%t err=%v", found, err)
			}
			started = true
			return exitOK
		}
	})
	var out, diagnostics bytes.Buffer
	if code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &out, &diagnostics); code != exitOK || !started {
		t.Fatalf("exit=%d started=%t stderr=%s", code, started, &diagnostics)
	}
	assertGCRetentionContains(t, diagnostics.String(), "Run Retention removed runs=1")
}

func TestDeliverStatusIsUnchangedByRunRetention(t *testing.T) {
	f := seedGCRetention(t, false, 2)
	repo := commandEnvironmentForTest(t).workDir
	mustMkdir(t, filepath.Join(repo, "docs", "specs"))
	s, loaded := retentionStartStore(t, f.home)
	queue := openDeliveryQueueForCLI(t, f.home, repo)
	// Use the queue API to record the durable Run link used for token totals.
	if err := s.UpdateDeliveryQueueItem(t.Context(), repo, queue.Items[0]); err != nil {
		t.Fatal(err)
	}
	db := openGCRetentionDB(t, f.home)
	gcRetentionExec(t, db, "UPDATE runs SET work_dir = '', completed_at = ? WHERE id <> ?", f.now.Add(-31*24*time.Hour).Format(time.RFC3339), f.queued.ID)
	gcRetentionExec(t, db, "UPDATE run_token_usage SET total_tokens = 123, basis = 'turn' WHERE run_id = ?", f.queued.ID)
	status := func() string {
		var out, diagnostics bytes.Buffer
		if code := runCLI(t, []string{"deliver", "status"}, &out, &diagnostics); code != exitOK || diagnostics.Len() != 0 {
			t.Fatalf("exit=%d stderr=%s", code, &diagnostics)
		}
		return out.String()
	}
	before := status()
	assertGCRetentionContains(t, before, "123")
	assertGCRetentionContains(t, retentionStartCall(t, s, loaded), "Run Retention removed runs=3")
	if after := status(); before != after {
		t.Fatalf("before=%q after=%q", before, after)
	}
	assertRetentionRun(t, s, f.queued.ID, true)
	for _, run := range []store.Run{f.removed, f.recent, f.worktree} {
		assertRetentionRun(t, s, run.ID, false)
	}
}

func TestRunRetentionKeepsTheRunReconcileReads(t *testing.T) {
	home, repo, location := newReconcileWorkspace(t)
	run, ref := createReconcileRun(t, home, repo, location, "ma/widget-flow", store.StateFailed)
	now := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)
	setRunTimestamps(t, home, run.ID, now.Add(-40*24*time.Hour), now.Add(-31*24*time.Hour))
	s, loaded := retentionStartStore(t, home)
	if output := retentionStartCall(t, s, loaded); output != "" {
		t.Fatalf("stderr=%q", output)
	}
	assertRetentionRun(t, s, run.ID, true)
	var out, diagnostics bytes.Buffer
	if code := runCLI(t, []string{"reconcile", run.ID, "--format", "text"}, &out, &diagnostics); code != exitOK {
		t.Fatalf("exit=%d stderr=%s", code, &diagnostics)
	}
	assertGCRetentionContains(t, out.String(), "Run: "+run.ID, "worktree: "+ref.Path)
}

func TestUnknownRunNamesRunRetention(t *testing.T) {
	home, _ := withCLIWorkspace(t)
	_, _ = retentionStartStore(t, home)
	now := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)
	for _, days := range []int{30, 15, 7} {
		mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), fmt.Sprintf("store:\n  run_retention_days: %d\n", days))
		cutoffID := "run_" + now.Add(-time.Duration(days)*24*time.Hour).Format("20060102T150405Z") + "_0123456789abcdef"
		for _, id := range []string{"run_20260801T120000Z_0123456789abcdef", "run_missing", cutoffID, "run_20261008T120000Z_0123456789abcdef", "run_20260801T120000Z_invalid", "run_20260230T120000Z_a", "run_20260801T120000Z_"} {
			message := fmt.Sprintf("Run %q does not exist", id)
			if id == "run_20260801T120000Z_0123456789abcdef" {
				message += fmt.Sprintf("; Run Retention may have removed it, because it removes terminal Runs that completed more than %d days ago", days)
			}
			for _, command := range []string{"show", "events"} {
				args := []string{"events", id}
				want := "roundfix events failed: " + message + "\n"
				if command == "show" {
					args = []string{"runs", "show", id}
					want = "roundfix: runs show failed: " + message + "\nRun 'roundfix runs show --help' for usage.\n"
				}
				var out, diagnostics bytes.Buffer
				if code := runCLI(t, args, &out, &diagnostics); code != exitPreflight || out.Len() != 0 || diagnostics.String() != want {
					t.Fatalf("%v exit=%d stdout=%q stderr=%q want=%q", args, code, out.String(), diagnostics.String(), want)
				}
			}
		}
	}
}
