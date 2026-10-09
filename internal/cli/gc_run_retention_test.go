// Suite: GC whole-Run retention through the CLI with isolated SQLite homes.
package cli

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

var gcRetentionTables = []string{"runs", "run_events", "run_agent_selections", "run_token_usage", "active_run_locks"}
var gcRetentionOtherTables = []string{"run_windows", "interactive_defaults", "delivery_queues", "delivery_queue_items", "delivery_queue_runs", "delivery_action_intents", "delivery_action_receipts"}

type gcRetentionFixture struct {
	home, root                                          string
	now                                                 time.Time
	removed, queued, worktree, recent, active, unproven store.Run
}

func seedGCRetention(t *testing.T, protected bool, mode int) gcRetentionFixture {
	t.Helper()
	ctx := context.Background()
	home, repo := withCLIWorkspace(t)
	f := gcRetentionFixture{home: home, root: filepath.Join(home, ".roundfix", "retention-artifacts"), now: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)}
	mustMkdir(t, filepath.Join(home, ".roundfix"))
	mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "store:\n  journal_retention: 0\n")
	withGCNow(t, f.now)
	s, err := store.Open(ctx, home)
	if err != nil {
		t.Fatal(err)
	}
	makeRun := func(name, root string, terminal bool) store.Run {
		run := createGCTestRun(t, ctx, s, root, name, 2)
		if _, err := s.AppendAgentSelectionAttempt(ctx, store.AgentSelectionAttemptRequest{RunID: run.ID, ScopeKind: store.AgentSelectionScopeTask, ScopeID: "task_02", Category: "backend", ProfileSource: "project", Attempt: 1, SelectionRole: store.AgentSelectionRolePreferred, Runtime: "codex", Model: "fixture", Status: store.AgentSelectionStatusClosed}); err != nil {
			t.Fatal(err)
		}
		if err := s.AppendTokenUsage(ctx, store.TokenUsageRecord{RunID: run.ID, ScopeKind: "task", ScopeID: "task_02", Session: "fixture", Basis: "unreported"}); err != nil {
			t.Fatal(err)
		}
		if terminal {
			if _, err := s.CompleteRun(ctx, run.ID, store.StateClean); err != nil {
				t.Fatal(err)
			}
		}
		writeRunArtifact(t, root, run.ID, name)
		return run
	}
	f.removed = makeRun("removed", f.root, true)
	f.queued = makeRun("queued", f.root, true)
	f.worktree = makeRun("worktree", f.root, true)
	f.recent = makeRun("recent", f.root, true)
	if protected {
		f.active = makeRun("active", f.root, false)
		f.unproven = makeRun("unproven", filepath.Join(t.TempDir(), "outside"), true)
	}
	if _, err := s.CreateDeliveryQueue(ctx, repo, []string{"fixture"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	db := openGCRetentionDB(t, home)
	gcRetentionExec(t, db, `UPDATE delivery_queue_items SET run_id = ? WHERE git_root = ?`, f.queued.ID, repo)
	gcRetentionExec(t, db, `UPDATE runs SET work_dir = ? WHERE id = ?`, repo, f.worktree.ID)
	for _, run := range []store.Run{f.removed, f.queued, f.worktree, f.unproven} {
		if run.ID != "" {
			setRunTimestamps(t, home, run.ID, f.now.Add(-40*24*time.Hour), f.now.Add(-31*24*time.Hour))
		}
	}
	setRunTimestamps(t, home, f.recent.ID, f.now.Add(-24*time.Hour), f.now.Add(-24*time.Hour))
	if f.active.ID != "" {
		setRunTimestamps(t, home, f.active.ID, f.now.Add(-40*24*time.Hour), time.Time{})
	}
	// Give each Run a substantial journal so compaction shrinks the database.
	gcRetentionExec(t, db, `UPDATE run_events SET summary = ?`, strings.Repeat("retention fixture ", 4096))
	gcRetentionExec(t, db, `INSERT INTO active_run_locks(target_kind,target_key,run_id,created_at) VALUES ('pr',?,?,?)`, f.removed.ID, f.removed.ID, f.now.Format(time.RFC3339))
	if mode != 2 {
		gcRetentionExec(t, db, fmt.Sprintf("PRAGMA auto_vacuum = %d", mode))
		gcRetentionExec(t, db, "VACUUM")
	}
	gcRetentionExec(t, db, "PRAGMA wal_checkpoint(TRUNCATE)")
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	return f
}

func openGCRetentionDB(t *testing.T, home string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(home))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	return db
}
func gcRetentionExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, err := db.Exec(query, args...); err != nil {
		t.Fatal(err)
	}
}
func gcRetentionRows(t *testing.T, home string, tables []string, runID string) map[string][]string {
	t.Helper()
	db := openGCRetentionDB(t, home)
	defer func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	}()
	result := map[string][]string{}
	for _, table := range tables {
		query := "SELECT * FROM " + table
		if runID != "" {
			column := "run_id"
			if table == "runs" {
				column = "id"
			}
			query += " WHERE " + column + " = '" + runID + "'"
		}
		rows, err := db.Query(query)
		if err != nil {
			t.Fatal(err)
		}
		cols, err := rows.Columns()
		if err != nil {
			t.Fatal(err)
		}
		result[table] = []string{}
		for rows.Next() {
			values := make([]any, len(cols))
			ptrs := make([]any, len(cols))
			for i := range values {
				ptrs[i] = &values[i]
			}
			if err := rows.Scan(ptrs...); err != nil {
				t.Fatal(err)
			}
			data, err := json.Marshal(values)
			if err != nil {
				t.Fatal(err)
			}
			result[table] = append(result[table], string(data))
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
		sort.Strings(result[table])
	}
	return result
}
func runGCRetentionCLI(t *testing.T, dry bool) string {
	t.Helper()
	args := []string{"gc"}
	if dry {
		args = append(args, "--dry-run")
	}
	var out, diagnostics bytes.Buffer
	if code := runCLIContext(t, context.Background(), args, &out, &diagnostics); code != exitOK || diagnostics.Len() != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, diagnostics.String(), out.String())
	}
	return out.String()
}
func assertGCRetentionContains(t *testing.T, output string, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(output, want) {
			t.Fatalf("missing %q in %q", want, output)
		}
	}
}
func gcRetentionFileSnapshot(t *testing.T, root string) map[string][]byte {
	t.Helper()
	files := map[string][]byte{}
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			files[path] = data
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return files
}

func TestGCDryRunReportsRunRetention(t *testing.T) {
	t.Parallel()
	f := seedGCRetention(t, false, 2)
	counts := gcRetentionRows(t, f.home, gcRetentionTables, f.removed.ID)
	before := gcRetentionFileSnapshot(t, filepath.Join(f.home, ".roundfix"))
	output := runGCRetentionCLI(t, true)
	assertGCRetentionContains(t, output, "GC dry-run", "Run Retention: 30 days", "Run Retention cutoff: "+f.now.Add(-30*24*time.Hour).Format(time.RFC3339), "Runs removable: 1", "Runs kept past the cutoff: 2 (queue-referenced 1, worktree present 1, artifact root unproven 0)", fmt.Sprintf("Rows removable: runs=1 run_events=%d run_agent_selections=%d run_token_usage=%d active_run_locks=1", len(counts["run_events"]), len(counts["run_agent_selections"]), len(counts["run_token_usage"])), "Removable Runs:\n    "+f.removed.ID)
	if gcReportInt64(t, output, "Database bytes reclaimable (estimated)") <= 0 {
		t.Fatal("no byte estimate")
	}
	after := gcRetentionFileSnapshot(t, filepath.Join(f.home, ".roundfix"))
	if !reflect.DeepEqual(before, after) {
		for path, data := range before {
			if !bytes.Equal(data, after[path]) {
				t.Errorf("changed %s: before=%d after=%d", path, len(data), len(after[path]))
			}
		}
		for path := range after {
			if _, ok := before[path]; !ok {
				t.Errorf("created %s", path)
			}
		}
		t.Fatal("dry run changed database or artifact bytes")
	}
}

func TestGCRemovesRunsPastRunRetention(t *testing.T) {
	t.Parallel()
	f := seedGCRetention(t, false, 2)
	kept := map[string]map[string][]string{}
	for _, run := range []store.Run{f.queued, f.worktree, f.recent} {
		kept[run.ID] = gcRetentionRows(t, f.home, gcRetentionTables, run.ID)
	}
	other := gcRetentionRows(t, f.home, gcRetentionOtherTables, "")
	counts := gcRetentionRows(t, f.home, gcRetentionTables, f.removed.ID)
	output := runGCRetentionCLI(t, false)
	assertGCRetentionContains(t, output, "GC complete", "Run Retention: 30 days", "Runs removed: 1", "Runs kept past the cutoff: 2 (queue-referenced 1, worktree present 1, artifact root unproven 0)", fmt.Sprintf("Rows removed: runs=1 run_events=%d run_agent_selections=%d run_token_usage=%d active_run_locks=1", len(counts["run_events"]), len(counts["run_agent_selections"]), len(counts["run_token_usage"])), "Run artifact bytes reclaimed: 7", "Removed Runs:\n    "+f.removed.ID)
	for table, rows := range gcRetentionRows(t, f.home, gcRetentionTables, f.removed.ID) {
		if len(rows) != 0 {
			t.Fatalf("%s still holds removed Run", table)
		}
	}
	if _, err := os.Lstat(filepath.Join(f.root, "runs", f.removed.ID)); !os.IsNotExist(err) {
		t.Fatalf("removed artifacts: %v", err)
	}
	for id, before := range kept {
		if !reflect.DeepEqual(before, gcRetentionRows(t, f.home, gcRetentionTables, id)) {
			t.Fatalf("kept Run %s changed", id)
		}
		assertPathExists(t, filepath.Join(f.root, "runs", id))
	}
	if !reflect.DeepEqual(other, gcRetentionRows(t, f.home, gcRetentionOtherTables, "")) {
		t.Fatal("unrelated durable tables changed")
	}
	s, err := store.OpenReader(context.Background(), f.home)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()
	sweep, found, err := s.LastRunRetentionSweep(context.Background())
	if err != nil || !found || !sweep.CompletedAt.Equal(f.now) || sweep.RetentionDays != 30 {
		t.Fatalf("sweep=%+v found=%t err=%v", sweep, found, err)
	}
}

func TestGCKeepsProtectedRuns(t *testing.T) {
	t.Parallel()
	f := seedGCRetention(t, true, 2)
	protected := []store.Run{f.queued, f.worktree, f.active, f.unproven, f.recent}
	before := map[string]map[string][]string{}
	for _, run := range protected {
		before[run.ID] = gcRetentionRows(t, f.home, gcRetentionTables, run.ID)
	}
	output := runGCRetentionCLI(t, false)
	assertGCRetentionContains(t, output, "Runs removed: 1", "Runs kept past the cutoff: 3 (queue-referenced 1, worktree present 1, artifact root unproven 1)")
	for _, run := range protected {
		if !reflect.DeepEqual(before[run.ID], gcRetentionRows(t, f.home, gcRetentionTables, run.ID)) {
			t.Fatalf("protected Run %s changed", run.ID)
		}
		assertPathExists(t, filepath.Join(run.ArtifactDir, "runs", run.ID))
	}
}

func TestGCKeepsARunUnderAnUnprovenArtifactRoot(t *testing.T) {
	t.Parallel()
	f := seedGCRetention(t, true, 2)
	path := filepath.Join(f.unproven.ArtifactDir, "runs", f.unproven.ID)
	before := gcRetentionFileSnapshot(t, path)
	rows := gcRetentionRows(t, f.home, gcRetentionTables, f.unproven.ID)
	output := runGCRetentionCLI(t, false)
	assertGCRetentionContains(t, output, "artifact root unproven 1")
	if !reflect.DeepEqual(before, gcRetentionFileSnapshot(t, path)) || !reflect.DeepEqual(rows, gcRetentionRows(t, f.home, gcRetentionTables, f.unproven.ID)) {
		t.Fatal("unproven Run changed")
	}
}

func TestGCCompactsAfterRunRetention(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		mode    int
		active  bool
		wording string
	}{
		{"incremental", 2, false, "Compaction: incremental ("},
		{"conversion", 0, false, "Compaction: full (converted to incremental)"},
		{"full mode", 1, false, "Compaction: full\n"},
		{"active blocks full", 0, true, "Compaction: skipped ("},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := seedGCRetention(t, tc.active, tc.mode)
			before := gcDatabaseFileSize(t, f.home)
			output := runGCRetentionCLI(t, false)
			assertGCRetentionContains(t, output, tc.wording)
			if !tc.active {
				if gcReportInt64(t, output, "Database bytes after") >= gcReportInt64(t, output, "Database bytes before") {
					t.Fatal("logical database did not shrink")
				}
				if gcDatabaseFileSize(t, f.home) >= before {
					t.Fatal("database file did not shrink")
				}
			} else {
				assertGCRetentionContains(t, output, f.active.ID)
			}
			if !tc.active {
				repeated := runGCRetentionCLI(t, false)
				assertGCRetentionContains(t, repeated, "Runs removed: 0", "Rows removed: runs=0 run_events=0 run_agent_selections=0 run_token_usage=0 active_run_locks=0", "Run artifact bytes reclaimed: 0", "Compaction: not needed")
			}
		})
	}
}

func TestGCRefusesAnUnsupportedRunRetention(t *testing.T) {
	t.Parallel()
	home, repo := withCLIWorkspace(t)
	mustMkdir(t, filepath.Join(home, ".roundfix"))
	mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "store:\n  run_retention_days: 10\n")
	mustWrite(t, filepath.Join(repo, "marker"), "unchanged")
	before := gcRetentionFileSnapshot(t, home)
	for _, args := range [][]string{{"gc", "--dry-run"}, {"gc"}} {
		var out, diagnostics bytes.Buffer
		code := runCLIContext(t, context.Background(), args, &out, &diagnostics)
		if code != exitPreflight || out.Len() != 0 {
			t.Fatalf("exit=%d stdout=%q", code, out.String())
		}
		assertGCRetentionContains(t, diagnostics.String(), "Preflight failed", fmt.Sprintf("parse config %q: store.run_retention_days must be 7, 15 or 30", filepath.Join(home, ".roundfix", "config.yml")))
	}
	if !reflect.DeepEqual(before, gcRetentionFileSnapshot(t, home)) {
		t.Fatal("refusal changed Home")
	}
}

func TestRunDatabaseSweepStopsAtBudgetBoundaries(t *testing.T) {
	t.Parallel()
	for _, phase := range []string{"removal", "compaction"} {
		t.Run(phase, func(t *testing.T) {
			f := seedGCRetention(t, false, 2)
			calls := 0
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
				deps.gc.now = func() time.Time {
					calls++
					limit := 1
					if phase == "compaction" {
						limit = 2
					}
					if calls > limit {
						return f.now.Add(3 * time.Second)
					}
					return f.now
				}
			})
			ctx := commandContextForTest(t, context.Background())
			loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: f.home, WorkDir: t.TempDir()})
			if err != nil {
				t.Fatal(err)
			}
			s, err := store.Open(ctx, f.home)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := s.Close(); err != nil {
					t.Error(err)
				}
			}()
			report, err := sweepRunDatabase(ctx, s, loaded, runRetentionOptions{budget: 2 * time.Second})
			if err != nil || !report.paused {
				t.Fatalf("report=%+v err=%v", report, err)
			}
			if _, found, err := s.LastRunRetentionSweep(ctx); err != nil || found {
				t.Fatalf("completion recorded: found=%t err=%v", found, err)
			}
			if phase == "removal" {
				if len(report.runIDs) != 0 {
					t.Fatalf("removed Runs=%v", report.runIDs)
				}
				assertPathExists(t, filepath.Join(f.root, "runs", f.removed.ID))
			} else {
				if len(report.runIDs) != 1 {
					t.Fatalf("removed Runs=%v", report.runIDs)
				}
				_, free, _, err := s.RunRetentionStorage(ctx)
				if err != nil || free == 0 {
					t.Fatalf("free pages=%d err=%v", free, err)
				}
			}
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) { deps.gc.now = func() time.Time { return f.now } })
			resumed, err := sweepRunDatabase(commandContextForTest(t, context.Background()), s, loaded, runRetentionOptions{budget: 2 * time.Second})
			if err != nil || resumed.paused {
				t.Fatalf("resume=%+v err=%v", resumed, err)
			}
			if _, found, err := s.LastRunRetentionSweep(ctx); err != nil || !found {
				t.Fatalf("resume completion: found=%t err=%v", found, err)
			}
		})
	}
}

func TestRunDatabaseBudgetNeverConvertsDefaultMode(t *testing.T) {
	t.Parallel()
	f := seedGCRetention(t, false, 0)
	ctx := commandContextForTest(t, context.Background())
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: f.home, WorkDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	s, err := store.Open(ctx, f.home)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			t.Error(err)
		}
	}()
	report, err := sweepRunDatabase(ctx, s, loaded, runRetentionOptions{budget: 2 * time.Second})
	if err != nil || report.paused || len(report.runIDs) != 1 {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	_, free, mode, err := s.RunRetentionStorage(ctx)
	if err != nil || mode != 0 || free == 0 {
		t.Fatalf("mode=%d free=%d err=%v", mode, free, err)
	}
	if _, found, err := s.LastRunRetentionSweep(ctx); err != nil || !found {
		t.Fatalf("completion: found=%t err=%v", found, err)
	}
}

func TestGCDryRunReadsCommittedWAL(t *testing.T) {
	t.Parallel()
	f := seedGCRetention(t, false, 2)
	db := openGCRetentionDB(t, f.home)
	gcRetentionExec(t, db, "UPDATE runs SET completed_at = ? WHERE id = ?", f.now.Add(-31*24*time.Hour).Format(time.RFC3339), f.recent.ID)
	snapshot := func() map[string][]byte {
		files := gcRetentionFileSnapshot(t, f.root)
		for _, path := range []string{store.DatabasePath(f.home), store.DatabasePath(f.home) + "-wal"} {
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			files[path] = content
		}
		return files
	}
	before := snapshot()
	output := runGCRetentionCLI(t, true)
	assertGCRetentionContains(t, output, "Runs removable: 2", f.removed.ID, f.recent.ID)
	if !reflect.DeepEqual(before, snapshot()) {
		t.Fatal("dry run changed database, WAL, or artifacts")
	}
}
