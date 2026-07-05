package store

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestOpenCreatesRunDatabaseAndAppliesMigrations(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()

	store := openTestStore(t, ctx, homeDir)
	defer closeStore(t, store)

	if _, err := os.Stat(DatabasePath(homeDir)); err != nil {
		t.Fatalf("expected Run Database file at %s: %v", DatabasePath(homeDir), err)
	}
	version, err := store.MigrationVersion(ctx)
	if err != nil {
		t.Fatalf("expected migration version, got %v", err)
	}
	if version != 6 {
		t.Fatalf("expected migration version 6, got %d", version)
	}
}

func TestInteractiveDefaultsRememberLastPullRequestAndAgent(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	if defaults, err := store.InteractiveDefaults(ctx); err != nil {
		t.Fatalf("read empty defaults: %v", err)
	} else if defaults.PRNumber != "" || defaults.Agent != "" {
		t.Fatalf("expected empty defaults, got %#v", defaults)
	}

	if err := store.RememberInteractiveDefaults(ctx, InteractiveDefaults{
		PRNumber: "123",
		Agent:    "codex",
	}); err != nil {
		t.Fatalf("remember defaults: %v", err)
	}
	if err := store.RememberInteractiveDefaults(ctx, InteractiveDefaults{
		PRNumber: "456",
	}); err != nil {
		t.Fatalf("update defaults: %v", err)
	}

	defaults, err := store.InteractiveDefaults(ctx)
	if err != nil {
		t.Fatalf("read defaults: %v", err)
	}
	if defaults.PRNumber != "456" {
		t.Fatalf("expected remembered PR 456, got %q", defaults.PRNumber)
	}
	if defaults.Agent != "codex" {
		t.Fatalf("expected remembered Agent codex, got %q", defaults.Agent)
	}
}

func TestCreateFetchRunCompletesFetchedAndReleasesLock(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	run, err := store.CreateFetchRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected Fetch Run creation, got %v", err)
	}
	if run.ID == "" {
		t.Fatal("expected Run id")
	}
	if run.Kind != KindFetch {
		t.Fatalf("expected kind %q, got %q", KindFetch, run.Kind)
	}
	if run.State != StateActive {
		t.Fatalf("expected active state, got %q", run.State)
	}
	if run.CompletedAt != nil {
		t.Fatalf("expected active run without completion timestamp, got %v", run.CompletedAt)
	}

	active, found, err := store.ActiveRun(ctx, run.HeadRepository, run.HeadBranch)
	if err != nil {
		t.Fatalf("expected active run lookup, got %v", err)
	}
	if !found || active.ID != run.ID {
		t.Fatalf("expected active lock for %s, found=%v active=%#v", run.ID, found, active)
	}

	completed, err := store.CompleteRun(ctx, run.ID, StateFetched)
	if err != nil {
		t.Fatalf("expected Fetched completion, got %v", err)
	}
	if completed.State != StateFetched {
		t.Fatalf("expected Fetched state, got %q", completed.State)
	}
	if completed.CompletedAt == nil {
		t.Fatal("expected completion timestamp")
	}
	_, found, err = store.ActiveRun(ctx, run.HeadRepository, run.HeadBranch)
	if err != nil {
		t.Fatalf("expected active lookup after release, got %v", err)
	}
	if found {
		t.Fatal("expected terminal Fetched run to release Active Run lock")
	}

	second, err := store.CreateFetchRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected new Fetch Run after lock release, got %v", err)
	}
	if second.ID == run.ID {
		t.Fatal("expected second Run to have a distinct id")
	}
}

func TestCreateRunRejectsDuplicateActiveRunWithoutNewRecord(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	first, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected first Run, got %v", err)
	}
	_, err = store.CreateRun(ctx, sampleCreateRunRequest())
	var activeErr ActiveRunError
	if !errors.As(err, &activeErr) {
		t.Fatalf("expected ActiveRunError, got %T %v", err, err)
	}
	if activeErr.Existing.ID != first.ID {
		t.Fatalf("expected existing run %s, got %s", first.ID, activeErr.Existing.ID)
	}
	if !strings.Contains(err.Error(), "existing run_id="+first.ID) {
		t.Fatalf("expected existing run id in error, got %q", err.Error())
	}
	count, err := store.RunCount(ctx)
	if err != nil {
		t.Fatalf("expected run count, got %v", err)
	}
	if count != 1 {
		t.Fatalf("expected duplicate active rejection to avoid new Run records, got count %d", count)
	}
}

func TestStoppedRunReleasesActiveLock(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	req := sampleCreateRunRequest()
	req.Kind = KindResolve
	run, err := store.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("expected active Resolve Run, got %v", err)
	}
	if _, err := store.CompleteRun(ctx, run.ID, StateStopped); err != nil {
		t.Fatalf("expected Stopped completion, got %v", err)
	}
	_, found, err := store.ActiveRun(ctx, run.HeadRepository, run.HeadBranch)
	if err != nil {
		t.Fatalf("expected active lookup, got %v", err)
	}
	if found {
		t.Fatal("expected Stopped terminal outcome to release Active Run lock")
	}

	second, err := store.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("expected new Run after Stopped lock release, got %v", err)
	}
	if second.ID == run.ID {
		t.Fatal("expected distinct run id")
	}
}

func TestRunLooksUpExistingRunByID(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	created, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected Run creation, got %v", err)
	}

	found, ok, err := store.Run(ctx, created.ID)
	if err != nil {
		t.Fatalf("lookup Run: %v", err)
	}
	if !ok || found.ID != created.ID {
		t.Fatalf("expected Run lookup for %s, ok=%v found=%#v", created.ID, ok, found)
	}

	_, ok, err = store.Run(ctx, "run_missing")
	if err != nil {
		t.Fatalf("lookup missing Run: %v", err)
	}
	if ok {
		t.Fatal("expected missing Run lookup")
	}
}

func TestCreateRunPersistsWorkDir(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	req := sampleCreateRunRequest()
	req.WorkDir = filepath.Join("tmp", "roundfix-home", "worktrees", "repo-id", "run-id")
	created, err := store.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("expected Run creation, got %v", err)
	}
	if created.WorkDir != req.WorkDir {
		t.Fatalf("expected created WorkDir %q, got %q", req.WorkDir, created.WorkDir)
	}

	found, ok, err := store.Run(ctx, created.ID)
	if err != nil || !ok {
		t.Fatalf("lookup persisted Run: ok=%v err=%v", ok, err)
	}
	if found.WorkDir != req.WorkDir {
		t.Fatalf("expected looked-up WorkDir %q, got %q", req.WorkDir, found.WorkDir)
	}

	active, ok, err := store.ActiveRun(ctx, req.HeadRepository, req.HeadBranch)
	if err != nil || !ok {
		t.Fatalf("lookup active Run: ok=%v err=%v", ok, err)
	}
	if active.WorkDir != req.WorkDir {
		t.Fatalf("expected active WorkDir %q, got %q", req.WorkDir, active.WorkDir)
	}

	inGitRoot, ok, err := store.ActiveRunInGitRoot(ctx, req.GitRoot)
	if err != nil || !ok {
		t.Fatalf("lookup active Run in Git root: ok=%v err=%v", ok, err)
	}
	if inGitRoot.WorkDir != req.WorkDir {
		t.Fatalf("expected Git-root active WorkDir %q, got %q", req.WorkDir, inGitRoot.WorkDir)
	}
}

func TestCreateRunAllowsDifferentHeadBranch(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	first, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected first Run, got %v", err)
	}
	secondReq := sampleCreateRunRequest()
	secondReq.HeadBranch = "feature/other-review"
	second, err := store.CreateRun(ctx, secondReq)
	if err != nil {
		t.Fatalf("expected simultaneous Run on different PR Head Branch, got %v", err)
	}
	if second.ID == first.ID {
		t.Fatal("expected distinct run ids")
	}
	count, err := store.RunCount(ctx)
	if err != nil {
		t.Fatalf("expected run count, got %v", err)
	}
	if count != 2 {
		t.Fatalf("expected two Run records, got %d", count)
	}
}

func TestCompleteRunAcceptsUnresolvedAsTerminal(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	run, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected Run creation, got %v", err)
	}
	completed, err := store.CompleteRun(ctx, run.ID, StateUnresolved)
	if err != nil {
		t.Fatalf("expected Unresolved completion, got %v", err)
	}
	if completed.State != StateUnresolved || completed.CompletedAt == nil {
		t.Fatalf("expected completed Unresolved Run, got %+v", completed)
	}
	_, found, err := store.ActiveRun(ctx, run.HeadRepository, run.HeadBranch)
	if err != nil {
		t.Fatalf("expected active lookup, got %v", err)
	}
	if found {
		t.Fatal("expected Unresolved Run to release the Active Run lock")
	}
}

func TestCompleteRunRejectsNonTerminalState(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	run, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected Run creation, got %v", err)
	}
	if _, err := store.CompleteRun(ctx, run.ID, StateActive); err == nil {
		t.Fatal("expected non-terminal completion to fail")
	}
	active, found, err := store.ActiveRun(ctx, run.HeadRepository, run.HeadBranch)
	if err != nil {
		t.Fatalf("expected active lookup, got %v", err)
	}
	if !found || active.ID != run.ID {
		t.Fatalf("expected active lock to remain after failed completion, found=%v active=%#v", found, active)
	}
}

func TestRequestStopRecordsStopRequestForActiveRun(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	now := time.Date(2026, 7, 5, 10, 30, 0, 0, time.UTC)
	runStore.now = func() time.Time { return now }

	run, err := runStore.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	requested, err := runStore.StopRequested(ctx, run.ID)
	if err != nil {
		t.Fatalf("read initial Stop Request flag: %v", err)
	}
	if requested {
		t.Fatal("expected new Active Run without a Stop Request")
	}
	if err := runStore.UpdateRunState(ctx, run.ID, StateResolvingWithAgent); err != nil {
		t.Fatalf("update active Run state: %v", err)
	}

	if err := runStore.RequestStop(ctx, run.ID); err != nil {
		t.Fatalf("request Stop: %v", err)
	}

	requested, err = runStore.StopRequested(ctx, run.ID)
	if err != nil {
		t.Fatalf("read Stop Request flag: %v", err)
	}
	if !requested {
		t.Fatal("expected Stop Request flag recorded")
	}
	active, found, err := runStore.ActiveRun(ctx, run.HeadRepository, run.HeadBranch)
	if err != nil {
		t.Fatalf("active run lookup after Stop Request: %v", err)
	}
	if !found || active.ID != run.ID {
		t.Fatalf("expected Stop Request to keep Active Run lock, found=%v active=%#v", found, active)
	}
	var recorded string
	if err := runStore.db.QueryRowContext(ctx, `SELECT stop_requested_at FROM runs WHERE id = ?`, run.ID).Scan(&recorded); err != nil {
		t.Fatalf("read stop_requested_at: %v", err)
	}
	if recorded != formatTime(now) {
		t.Fatalf("expected stop_requested_at %q, got %q", formatTime(now), recorded)
	}
}

func TestRequestStopRejectsTerminalRunWithNamedError(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	run, err := runStore.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	if _, err := runStore.CompleteRun(ctx, run.ID, StateClean); err != nil {
		t.Fatalf("complete Run: %v", err)
	}

	err = runStore.RequestStop(ctx, run.ID)

	if !errors.Is(err, ErrTerminalRunStopRequest) {
		t.Fatalf("expected ErrTerminalRunStopRequest, got %T %v", err, err)
	}
	requested, flagErr := runStore.StopRequested(ctx, run.ID)
	if flagErr != nil {
		t.Fatalf("read Stop Request flag: %v", flagErr)
	}
	if requested {
		t.Fatal("expected terminal Run to keep Stop Request flag unset")
	}
}

func openTestStore(t *testing.T, ctx context.Context, homeDir string) *Store {
	t.Helper()
	store, err := Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	return store
}

func closeStore(t *testing.T, store *Store) {
	t.Helper()
	if err := store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
}

func sampleCreateRunRequest() CreateRunRequest {
	return CreateRunRequest{
		Kind:           KindFetch,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        filepath.Join("tmp", "repo"),
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join("tmp", "repo", ".roundfix"),
	}
}

// buildV3Fixture creates a populated schema v3 Run Database via raw SQL:
// runs in several states plus one Active Run lock in the v3
// (head_repository, head_branch) shape.
func buildV3Fixture(t *testing.T, homeDir string) {
	t.Helper()
	path := DatabasePath(homeDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}
	db, err := sql.Open("sqlite", writerDSN(path))
	if err != nil {
		t.Fatalf("open fixture database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close fixture database: %v", err)
		}
	}()

	statements := []string{
		`CREATE TABLE runs (
			id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			state TEXT NOT NULL,
			head_repository TEXT NOT NULL,
			head_branch TEXT NOT NULL,
			base_repository TEXT NOT NULL DEFAULT '',
			pr_number TEXT NOT NULL,
			git_root TEXT NOT NULL,
			local_branch TEXT NOT NULL,
			head_sha TEXT NOT NULL,
			artifact_dir TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			completed_at TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE active_run_locks (
			head_repository TEXT NOT NULL,
			head_branch TEXT NOT NULL,
			run_id TEXT NOT NULL UNIQUE,
			created_at TEXT NOT NULL,
			PRIMARY KEY (head_repository, head_branch),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX idx_runs_head ON runs (head_repository, head_branch)`,
		`CREATE TABLE interactive_defaults (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE run_events (
			run_id TEXT NOT NULL,
			cursor INTEGER NOT NULL,
			batch INTEGER NOT NULL DEFAULT 0,
			source TEXT NOT NULL,
			kind TEXT NOT NULL,
			review_issue TEXT NOT NULL DEFAULT '',
			tool_id TEXT NOT NULL DEFAULT '',
			tool_state TEXT NOT NULL DEFAULT '',
			summary TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			payload TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (run_id, cursor),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, created_at, updated_at, completed_at)
		 VALUES ('run_v3_active', 'resolve', 'Active', 'owner/project', 'feature/review', 'owner/project',
			'123', 'tmp/repo', 'feature/review', 'abc123', 'tmp/repo/.roundfix',
			'2026-07-01T10:00:00Z', '2026-07-01T10:00:00Z', '')`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, created_at, updated_at, completed_at)
		 VALUES ('run_v3_clean', 'resolve', 'Clean', 'owner/project', 'feature/done', 'owner/project',
			'99', 'tmp/repo', 'feature/done', 'def456', 'tmp/repo/.roundfix',
			'2026-07-01T08:00:00Z', '2026-07-01T09:00:00Z', '2026-07-01T09:00:00Z')`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, created_at, updated_at, completed_at)
		 VALUES ('run_v3_fetched', 'fetch', 'Fetched', 'owner/other', 'feature/fetch', 'owner/other',
			'7', 'tmp/other', 'feature/fetch', 'fed789', 'tmp/other/.roundfix',
			'2026-07-01T07:00:00Z', '2026-07-01T07:30:00Z', '2026-07-01T07:30:00Z')`,
		`INSERT INTO active_run_locks (head_repository, head_branch, run_id, created_at)
		 VALUES ('owner/project', 'feature/review', 'run_v3_active', '2026-07-01T10:00:00Z')`,
		`PRAGMA user_version = 3`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("build v3 fixture: %v", err)
		}
	}
}

func TestOpenMigratesV3RunDatabasePreservingRunsAndRekeyingLocks(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	buildV3Fixture(t, homeDir)

	store := openTestStore(t, ctx, homeDir)
	defer closeStore(t, store)

	version, err := store.MigrationVersion(ctx)
	if err != nil {
		t.Fatalf("read migration version: %v", err)
	}
	if version != 6 {
		t.Fatalf("expected user_version 6 after migration, got %d", version)
	}

	count, err := store.RunCount(ctx)
	if err != nil {
		t.Fatalf("count runs: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected all 3 v3 run rows to survive, got %d", count)
	}

	clean, ok, err := store.Run(ctx, "run_v3_clean")
	if err != nil || !ok {
		t.Fatalf("expected run_v3_clean to survive, ok=%v err=%v", ok, err)
	}
	if clean.Kind != KindResolve || clean.State != StateClean || clean.PRNumber != "99" ||
		clean.HeadRepository != "owner/project" || clean.HeadBranch != "feature/done" ||
		clean.HeadSHA != "def456" || clean.ArtifactDir != "tmp/repo/.roundfix" ||
		clean.WorkDir != "" || clean.SpecSlug != "" || clean.CompletedAt == nil {
		t.Fatalf("expected run_v3_clean fields preserved, got %#v", clean)
	}

	active, found, err := store.ActiveRun(ctx, "owner/project", "feature/review")
	if err != nil {
		t.Fatalf("active run lookup after migration: %v", err)
	}
	if !found || active.ID != "run_v3_active" {
		t.Fatalf("expected re-keyed lock to keep run_v3_active active, found=%v active=%#v", found, active)
	}

	var lockCount int
	if err := store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM active_run_locks`).Scan(&lockCount); err != nil {
		t.Fatalf("count migrated locks: %v", err)
	}
	if lockCount != 1 {
		t.Fatalf("expected exactly one migrated lock row, got %d", lockCount)
	}
	var targetKind, targetKey, runID string
	if err := store.db.QueryRowContext(ctx,
		`SELECT target_kind, target_key, run_id FROM active_run_locks`).Scan(&targetKind, &targetKey, &runID); err != nil {
		t.Fatalf("read migrated lock row: %v", err)
	}
	if targetKind != "pr" || targetKey != "owner/project#feature/review" || runID != "run_v3_active" {
		t.Fatalf("expected lock re-keyed to (pr, owner/project#feature/review, run_v3_active), got (%s, %s, %s)",
			targetKind, targetKey, runID)
	}

	implementRun, err := store.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("expected implement Run creation on migrated database, got %v", err)
	}
	if implementRun.Kind != KindImplement || implementRun.SpecSlug != "0001-implement-command" {
		t.Fatalf("expected persisted implement Run with Spec slug, got %#v", implementRun)
	}
}

// buildV4Fixture creates a populated schema v4 Run Database via raw SQL:
// runs in several states plus one Active Run lock in the v4 work-target
// shape. It intentionally omits stop_requested_at so the v5 migration must
// add it without disturbing existing rows or locks.
func buildV4Fixture(t *testing.T, homeDir string) {
	t.Helper()
	path := DatabasePath(homeDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}
	db, err := sql.Open("sqlite", writerDSN(path))
	if err != nil {
		t.Fatalf("open fixture database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close fixture database: %v", err)
		}
	}()

	statements := []string{
		`CREATE TABLE runs (
			id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			state TEXT NOT NULL,
			head_repository TEXT NOT NULL DEFAULT '',
			head_branch TEXT NOT NULL DEFAULT '',
			base_repository TEXT NOT NULL DEFAULT '',
			pr_number TEXT NOT NULL DEFAULT '',
			git_root TEXT NOT NULL,
			local_branch TEXT NOT NULL,
			head_sha TEXT NOT NULL DEFAULT '',
			artifact_dir TEXT NOT NULL DEFAULT '',
			spec_slug TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			completed_at TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE active_run_locks ` + activeRunLocksColumns,
		`CREATE INDEX idx_runs_head ON runs (head_repository, head_branch)`,
		`CREATE TABLE interactive_defaults (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE run_events (
			run_id TEXT NOT NULL,
			cursor INTEGER NOT NULL,
			batch INTEGER NOT NULL DEFAULT 0,
			source TEXT NOT NULL,
			kind TEXT NOT NULL,
			review_issue TEXT NOT NULL DEFAULT '',
			tool_id TEXT NOT NULL DEFAULT '',
			tool_state TEXT NOT NULL DEFAULT '',
			summary TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			payload TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (run_id, cursor),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, spec_slug, created_at, updated_at, completed_at)
		 VALUES ('run_v4_active', 'resolve', 'ResolvingWithAgent', 'owner/project', 'feature/review', 'owner/project',
			'123', 'tmp/repo', 'feature/review', 'abc123', 'tmp/repo/.roundfix', '',
			'2026-07-01T10:00:00Z', '2026-07-01T10:05:00Z', '')`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, spec_slug, created_at, updated_at, completed_at)
		 VALUES ('run_v4_clean', 'watch', 'Clean', 'owner/project', 'feature/done', 'owner/project',
			'99', 'tmp/repo', 'feature/done', 'def456', 'tmp/repo/.roundfix', '',
			'2026-07-01T08:00:00Z', '2026-07-01T09:00:00Z', '2026-07-01T09:00:00Z')`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, spec_slug, created_at, updated_at, completed_at)
		 VALUES ('run_v4_implement', 'implement', 'Stopped', '', '', '',
			'', 'tmp/spec-repo', 'ma/spec-work', '', '', '0001-widget-flow',
			'2026-07-01T07:00:00Z', '2026-07-01T07:30:00Z', '2026-07-01T07:30:00Z')`,
		`INSERT INTO active_run_locks (target_kind, target_key, run_id, created_at)
		 VALUES ('pr', 'owner/project#feature/review', 'run_v4_active', '2026-07-01T10:00:00Z')`,
		`PRAGMA user_version = 4`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("build v4 fixture: %v", err)
		}
	}
}

func TestOpenMigratesV4RunDatabasePreservingRunsLocksAndAddingStopRequests(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	buildV4Fixture(t, homeDir)

	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)

	version, err := runStore.MigrationVersion(ctx)
	if err != nil {
		t.Fatalf("read migration version: %v", err)
	}
	if version != 6 {
		t.Fatalf("expected user_version 6 after migration, got %d", version)
	}
	count, err := runStore.RunCount(ctx)
	if err != nil {
		t.Fatalf("count runs: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected all 3 v4 run rows to survive, got %d", count)
	}
	clean, ok, err := runStore.Run(ctx, "run_v4_clean")
	if err != nil || !ok {
		t.Fatalf("expected run_v4_clean to survive, ok=%v err=%v", ok, err)
	}
	if clean.Kind != KindWatch || clean.State != StateClean || clean.PRNumber != "99" ||
		clean.HeadRepository != "owner/project" || clean.HeadBranch != "feature/done" ||
		clean.HeadSHA != "def456" || clean.WorkDir != "" || clean.SpecSlug != "" || clean.CompletedAt == nil {
		t.Fatalf("expected run_v4_clean fields preserved, got %#v", clean)
	}
	implement, ok, err := runStore.Run(ctx, "run_v4_implement")
	if err != nil || !ok {
		t.Fatalf("expected run_v4_implement to survive, ok=%v err=%v", ok, err)
	}
	if implement.Kind != KindImplement || implement.SpecSlug != "0001-widget-flow" || implement.GitRoot != "tmp/spec-repo" ||
		implement.WorkDir != "" || implement.CompletedAt == nil {
		t.Fatalf("expected implement fields preserved, got %#v", implement)
	}
	active, found, err := runStore.ActiveRun(ctx, "owner/project", "feature/review")
	if err != nil {
		t.Fatalf("active run lookup after migration: %v", err)
	}
	if !found || active.ID != "run_v4_active" || active.State != StateResolvingWithAgent {
		t.Fatalf("expected v4 active lock to survive, found=%v active=%#v", found, active)
	}
	requested, err := runStore.StopRequested(ctx, "run_v4_active")
	if err != nil {
		t.Fatalf("read migrated Stop Request flag: %v", err)
	}
	if requested {
		t.Fatal("expected migrated v4 Run to have no Stop Request")
	}
	var stopRequestedAt any
	if err := runStore.db.QueryRowContext(ctx, `SELECT stop_requested_at FROM runs WHERE id = 'run_v4_active'`).Scan(&stopRequestedAt); err != nil {
		t.Fatalf("read migrated stop_requested_at column: %v", err)
	}
	if stopRequestedAt != nil {
		t.Fatalf("expected migrated stop_requested_at NULL, got %#v", stopRequestedAt)
	}
}

// buildV5Fixture creates a populated schema v5 Run Database via raw SQL:
// persisted rows have agents and Stop Request state but no work_dir column.
func buildV5Fixture(t *testing.T, homeDir string) {
	t.Helper()
	path := DatabasePath(homeDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}
	db, err := sql.Open("sqlite", writerDSN(path))
	if err != nil {
		t.Fatalf("open fixture database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close fixture database: %v", err)
		}
	}()

	statements := []string{
		`CREATE TABLE runs (
			id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			state TEXT NOT NULL,
			head_repository TEXT NOT NULL DEFAULT '',
			head_branch TEXT NOT NULL DEFAULT '',
			base_repository TEXT NOT NULL DEFAULT '',
			pr_number TEXT NOT NULL DEFAULT '',
			git_root TEXT NOT NULL,
			local_branch TEXT NOT NULL,
			head_sha TEXT NOT NULL DEFAULT '',
			artifact_dir TEXT NOT NULL DEFAULT '',
			spec_slug TEXT NOT NULL DEFAULT '',
			agent TEXT NOT NULL DEFAULT '',
			stop_requested_at TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			completed_at TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE active_run_locks ` + activeRunLocksColumns,
		`CREATE INDEX idx_runs_head ON runs (head_repository, head_branch)`,
		`CREATE TABLE interactive_defaults (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE run_events (
			run_id TEXT NOT NULL,
			cursor INTEGER NOT NULL,
			batch INTEGER NOT NULL DEFAULT 0,
			source TEXT NOT NULL,
			kind TEXT NOT NULL,
			review_issue TEXT NOT NULL DEFAULT '',
			tool_id TEXT NOT NULL DEFAULT '',
			tool_state TEXT NOT NULL DEFAULT '',
			summary TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			payload TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (run_id, cursor),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, spec_slug, agent,
			stop_requested_at, created_at, updated_at, completed_at)
		 VALUES ('run_v5_active', 'resolve', 'Verifying', 'owner/project', 'feature/review', 'owner/project',
			'123', 'tmp/repo', 'feature/review', 'abc123', 'tmp/repo/.roundfix', '', 'codex',
			'2026-07-01T10:06:00Z', '2026-07-01T10:00:00Z', '2026-07-01T10:06:00Z', '')`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, spec_slug, agent,
			stop_requested_at, created_at, updated_at, completed_at)
		 VALUES ('run_v5_clean', 'watch', 'Clean', 'owner/project', 'feature/done', 'owner/project',
			'99', 'tmp/repo', 'feature/done', 'def456', 'tmp/repo/.roundfix', '', 'claude',
			NULL, '2026-07-01T08:00:00Z', '2026-07-01T09:00:00Z', '2026-07-01T09:00:00Z')`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, spec_slug, agent,
			stop_requested_at, created_at, updated_at, completed_at)
		 VALUES ('run_v5_implement', 'implement', 'Active', '', '', '',
			'', 'tmp/spec-repo', 'ma/spec-work', '', '', '0001-widget-flow', 'opencode',
			NULL, '2026-07-01T07:00:00Z', '2026-07-01T07:30:00Z', '')`,
		`INSERT INTO active_run_locks (target_kind, target_key, run_id, created_at)
		 VALUES ('pr', 'owner/project#feature/review', 'run_v5_active', '2026-07-01T10:00:00Z')`,
		`INSERT INTO active_run_locks (target_kind, target_key, run_id, created_at)
		 VALUES ('spec', 'tmp/spec-repo#0001-widget-flow', 'run_v5_implement', '2026-07-01T07:00:00Z')`,
		`PRAGMA user_version = 5`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("build v5 fixture: %v", err)
		}
	}
}

func TestOpenMigratesV5RunDatabasePreservingRunsLocksAndAddingWorkDir(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	buildV5Fixture(t, homeDir)

	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)

	version, err := runStore.MigrationVersion(ctx)
	if err != nil {
		t.Fatalf("read migration version: %v", err)
	}
	if version != 6 {
		t.Fatalf("expected user_version 6 after migration, got %d", version)
	}
	count, err := runStore.RunCount(ctx)
	if err != nil {
		t.Fatalf("count runs: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected all 3 v5 run rows to survive, got %d", count)
	}

	active, found, err := runStore.ActiveRun(ctx, "owner/project", "feature/review")
	if err != nil {
		t.Fatalf("active review Run lookup after migration: %v", err)
	}
	if !found || active.ID != "run_v5_active" || active.State != StateVerifying ||
		active.Agent != "codex" || active.WorkDir != "" {
		t.Fatalf("expected v5 review lock and row to survive with empty WorkDir, found=%v active=%#v", found, active)
	}
	requested, err := runStore.StopRequested(ctx, "run_v5_active")
	if err != nil {
		t.Fatalf("read migrated Stop Request flag: %v", err)
	}
	if !requested {
		t.Fatal("expected populated v5 Stop Request to survive")
	}

	implement, found, err := runStore.ActiveSpecRun(ctx, "tmp/spec-repo", "0001-widget-flow")
	if err != nil {
		t.Fatalf("active spec Run lookup after migration: %v", err)
	}
	if !found || implement.ID != "run_v5_implement" || implement.Agent != "opencode" || implement.WorkDir != "" {
		t.Fatalf("expected v5 spec lock and row to survive with empty WorkDir, found=%v implement=%#v", found, implement)
	}

	clean, ok, err := runStore.Run(ctx, "run_v5_clean")
	if err != nil || !ok {
		t.Fatalf("expected run_v5_clean to survive, ok=%v err=%v", ok, err)
	}
	if clean.Kind != KindWatch || clean.State != StateClean || clean.Agent != "claude" ||
		clean.WorkDir != "" || clean.CompletedAt == nil {
		t.Fatalf("expected run_v5_clean fields preserved with empty WorkDir, got %#v", clean)
	}

	var rawWorkDir any
	if err := runStore.db.QueryRowContext(ctx, `SELECT work_dir FROM runs WHERE id = 'run_v5_active'`).Scan(&rawWorkDir); err != nil {
		t.Fatalf("read migrated work_dir column: %v", err)
	}
	if rawWorkDir != nil {
		t.Fatalf("expected migrated work_dir NULL for legacy row, got %#v", rawWorkDir)
	}
	var lockCount int
	if err := runStore.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM active_run_locks`).Scan(&lockCount); err != nil {
		t.Fatalf("count migrated locks: %v", err)
	}
	if lockCount != 2 {
		t.Fatalf("expected both v5 locks to survive, got %d", lockCount)
	}
}

func TestCreateRunRejectsSecondActiveRunForSameSpecTarget(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	first, err := store.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("expected first implement Run, got %v", err)
	}

	_, err = store.CreateRun(ctx, sampleImplementCreateRunRequest())
	var activeErr ActiveRunError
	if !errors.As(err, &activeErr) {
		t.Fatalf("expected ActiveRunError for same Spec target, got %T %v", err, err)
	}
	if activeErr.Existing.ID != first.ID {
		t.Fatalf("expected blocking run %s, got %s", first.ID, activeErr.Existing.ID)
	}
	want := `Active Run already exists for repository "tmp/spec-repo" and Spec "0001-implement-command"; existing run_id=` +
		first.ID + ` state=Active; stop it with: roundfix stop ` + first.ID
	if err.Error() != want {
		t.Fatalf("expected spec-target error naming the work target and blocking run,\nwant %q\ngot  %q", want, err.Error())
	}

	otherSlug := sampleImplementCreateRunRequest()
	otherSlug.SpecSlug = "0002-other-feature"
	if _, err := store.CreateRun(ctx, otherSlug); err != nil {
		t.Fatalf("expected different Spec slug in same repository to pass the lock, got %v", err)
	}

	otherRepo := sampleImplementCreateRunRequest()
	otherRepo.GitRoot = "tmp/other-repo"
	if _, err := store.CreateRun(ctx, otherRepo); err != nil {
		t.Fatalf("expected same Spec slug in different repository to pass the lock, got %v", err)
	}
}

func TestCompletedImplementRunReleasesSpecTargetLock(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	first, err := store.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("expected implement Run, got %v", err)
	}
	if _, err := store.CompleteRun(ctx, first.ID, StateStopped); err != nil {
		t.Fatalf("expected Stopped completion, got %v", err)
	}
	second, err := store.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("expected new implement Run after lock release, got %v", err)
	}
	if second.ID == first.ID {
		t.Fatal("expected distinct run id")
	}
}

func TestReviewKindActiveRunErrorTextUnchanged(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	first, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected first Run, got %v", err)
	}
	_, err = store.CreateRun(ctx, sampleCreateRunRequest())
	if err == nil {
		t.Fatal("expected duplicate review Run rejection")
	}
	want := `Active Run already exists for Head Repository "owner/project" and PR Head Branch "feature/review"; existing run_id=` +
		first.ID + ` state=Active`
	if err.Error() != want {
		t.Fatalf("expected review-path error text unchanged,\nwant %q\ngot  %q", want, err.Error())
	}
}

func TestActiveRunInGitRootFindsActiveRunsOfAnyKind(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	review, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected review Run, got %v", err)
	}
	found, ok, err := store.ActiveRunInGitRoot(ctx, review.GitRoot)
	if err != nil {
		t.Fatalf("active run in git root: %v", err)
	}
	if !ok || found.ID != review.ID {
		t.Fatalf("expected review Run active in git root, ok=%v found=%#v", ok, found)
	}

	if _, ok, err := store.ActiveRunInGitRoot(ctx, "tmp/elsewhere"); err != nil || ok {
		t.Fatalf("expected no Active Run in unrelated git root, ok=%v err=%v", ok, err)
	}

	if _, err := store.CompleteRun(ctx, review.ID, StateClean); err != nil {
		t.Fatalf("complete review Run: %v", err)
	}
	if _, ok, err := store.ActiveRunInGitRoot(ctx, review.GitRoot); err != nil || ok {
		t.Fatalf("expected completed Run to leave git root free, ok=%v err=%v", ok, err)
	}

	implementReq := sampleImplementCreateRunRequest()
	implementReq.GitRoot = review.GitRoot
	implementRun, err := store.CreateRun(ctx, implementReq)
	if err != nil {
		t.Fatalf("expected implement Run, got %v", err)
	}
	found, ok, err = store.ActiveRunInGitRoot(ctx, review.GitRoot)
	if err != nil {
		t.Fatalf("active run in git root after implement create: %v", err)
	}
	if !ok || found.ID != implementRun.ID || found.Kind != KindImplement {
		t.Fatalf("expected implement Run active in git root, ok=%v found=%#v", ok, found)
	}
}

func TestCreateRunValidatesRequiredFieldsByKind(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	tests := []struct {
		name    string
		request func() CreateRunRequest
		wantErr string
	}{
		{
			name: "implement missing Spec slug",
			request: func() CreateRunRequest {
				req := sampleImplementCreateRunRequest()
				req.SpecSlug = ""
				return req
			},
			wantErr: "Spec slug is required to create a Run",
		},
		{
			name: "implement missing Git root",
			request: func() CreateRunRequest {
				req := sampleImplementCreateRunRequest()
				req.GitRoot = ""
				return req
			},
			wantErr: "Git root is required to create a Run",
		},
		{
			name: "implement missing local branch",
			request: func() CreateRunRequest {
				req := sampleImplementCreateRunRequest()
				req.LocalBranch = ""
				return req
			},
			wantErr: "local branch is required to create a Run",
		},
		{
			name: "review missing pull request",
			request: func() CreateRunRequest {
				req := sampleCreateRunRequest()
				req.PRNumber = ""
				return req
			},
			wantErr: "pull request is required to create a Run",
		},
		{
			name: "review missing Head Repository",
			request: func() CreateRunRequest {
				req := sampleCreateRunRequest()
				req.HeadRepository = ""
				return req
			},
			wantErr: "Head Repository is required to create a Run",
		},
		{
			name: "review missing HEAD",
			request: func() CreateRunRequest {
				req := sampleCreateRunRequest()
				req.HeadSHA = ""
				return req
			},
			wantErr: "HEAD is required to create a Run",
		},
		{
			name: "unknown kind",
			request: func() CreateRunRequest {
				req := sampleCreateRunRequest()
				req.Kind = "deploy"
				return req
			},
			wantErr: `Run kind "deploy" is invalid`,
		},
		{
			name: "empty kind",
			request: func() CreateRunRequest {
				req := sampleCreateRunRequest()
				req.Kind = ""
				return req
			},
			wantErr: "Run kind is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := store.CreateRun(ctx, tt.request())
			if err == nil {
				t.Fatalf("expected validation error %q, got nil", tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func sampleImplementCreateRunRequest() CreateRunRequest {
	return CreateRunRequest{
		Kind:        KindImplement,
		GitRoot:     "tmp/spec-repo",
		LocalBranch: "ma/implement-spec",
		SpecSlug:    "0001-implement-command",
		ArtifactDir: "tmp/spec-repo/.roundfix",
	}
}

func TestUpdateRunStateRejectsTerminalStatesAndMissingRuns(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	run, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run: %v", err)
	}

	if err := store.UpdateRunState(ctx, run.ID, StateResolvingWithAgent); err != nil {
		t.Fatalf("expected intermediate state update, got %v", err)
	}
	updated, _, err := store.Run(ctx, run.ID)
	if err != nil {
		t.Fatalf("lookup run: %v", err)
	}
	if updated.State != StateResolvingWithAgent {
		t.Fatalf("expected ResolvingWithAgent, got %q", updated.State)
	}

	if err := store.UpdateRunState(ctx, run.ID, StateClean); err == nil {
		t.Fatal("expected terminal state rejection; terminal outcomes go through CompleteRun")
	}
	if err := store.UpdateRunState(ctx, "run_missing", StateVerifying); err == nil {
		t.Fatal("expected missing Run rejection")
	}
}
