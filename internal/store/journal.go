package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"roundfix/internal/runevent"
)

// JournalEvent is one persisted Run Event with its per-Run replay cursor.
// The cursor is opaque to callers: compare it by ordering and use it as the
// next replay position.
type JournalEvent struct {
	Cursor int64
	Event  runevent.RunEvent
}

// PruneResult reports the Run Event Journal rows removed for eligible Runs.
type PruneResult struct {
	RunIDs []string
	Events int
}

// PruneCandidate describes one terminal Run eligible for Run Event pruning.
type PruneCandidate struct {
	RunID  string
	Events int
}

// CompactionReconciliationToleranceReason explains why compaction results are
// compared within one SQLite page instead of against a recorded file size.
const CompactionReconciliationToleranceReason = "one SQLite page allows the compact snapshot and checkpointed Run Database to land on adjacent file boundaries"

// CompactionPreview reports the measured effect of rebuilding the Run
// Database. The private fingerprint binds the preview to the logical database
// state that produced it, so Compact can refuse a stale or forged preview
// before changing the original file.
type CompactionPreview struct {
	BytesBefore                   int64
	BytesReclaimable              int64
	BytesAfter                    int64
	ReconciliationToleranceBytes  int64
	ReconciliationToleranceReason string

	fingerprint        compactionFingerprint
	snapshotBytesAfter int64
}

// CompactionResult reports the measured file-size change after compaction.
type CompactionResult struct {
	BytesBefore    int64
	BytesReclaimed int64
	BytesAfter     int64
}

// ActiveRunCompactionError refuses compaction while one non-terminal Run
// remains in the Run Database, including Runs created without an Active Run
// lock.
type ActiveRunCompactionError struct {
	RunID string
	State string
}

func (err ActiveRunCompactionError) Error() string {
	return fmt.Sprintf("compact Run Database: Active Run %q in state %q blocks compaction", err.RunID, err.State)
}

// WriterPresentCompactionError refuses compaction when SQLite cannot prove
// that this connection is the only connection that can write the Run
// Database.
type WriterPresentCompactionError struct {
	Cause error
}

func (err WriterPresentCompactionError) Error() string {
	return fmt.Sprintf("compact Run Database: another Run Database connection can be writing: %v", err.Cause)
}

func (err WriterPresentCompactionError) Unwrap() error {
	return err.Cause
}

// CompactionCapacityError reports the measured temporary-space shortfall
// before SQLite starts rebuilding the original Run Database.
type CompactionCapacityError struct {
	RequiredBytes  int64
	AvailableBytes int64
	ShortfallBytes int64
}

func (err CompactionCapacityError) Error() string {
	return fmt.Sprintf(
		"compact Run Database: insufficient temporary capacity: required=%d available=%d shortfall=%d",
		err.RequiredBytes,
		err.AvailableBytes,
		err.ShortfallBytes,
	)
}

// CompactionPreviewStaleError refuses a preview whose database identity or
// logical page state changed before compaction acquired exclusive access.
type CompactionPreviewStaleError struct{}

func (CompactionPreviewStaleError) Error() string {
	return "compact Run Database: compaction preview is stale; preview again after all writers stop"
}

// CompactionResultMismatchError reports a completed SQLite compaction whose
// measured reclaim differs from its preview beyond the declared tolerance.
type CompactionResultMismatchError struct {
	PreviewReclaimable int64
	ResultReclaimed    int64
	ToleranceBytes     int64
}

func (err CompactionResultMismatchError) Error() string {
	return fmt.Sprintf(
		"compact Run Database: reclaimed bytes differ from preview: preview=%d reclaimed=%d tolerance=%d",
		err.PreviewReclaimable,
		err.ResultReclaimed,
		err.ToleranceBytes,
	)
}

type compactionFingerprint struct {
	databasePath string
	pageSize     int64
	pageCount    int64
	freePages    int64
	dataVersion  int64
	totalChanges int64
}

// StorageReportReconciliationToleranceReason explains the page-sized
// reconciliation tolerance. SQLite page accounting and the independent file
// stat are both page-granular but can be observed on adjacent page boundaries
// because the report deliberately takes no writer lock.
const StorageReportReconciliationToleranceReason = "one SQLite page allows page accounting and file stat to land on adjacent boundaries without a writer lock"

// StorageReport is a read-only measurement of the Run Database and every
// recorded Artifact Root. Table bytes are SQLite page allocations. Repository
// and state bytes are regular-file bytes under their recorded runs/<run-id>
// directories, so the same Run artifact is never counted twice within either
// grouping.
type StorageReport struct {
	DatabasePath                  string
	DatabaseBytes                 int64
	DatabaseAllocatedBytes        int64
	DatabaseFreeBytes             int64
	ArtifactBytes                 int64
	ReconciliationToleranceBytes  int64
	ReconciliationToleranceReason string
	Repositories                  []StorageRepositoryGroup
	States                        []StorageStateGroup
	Tables                        []StorageTableGroup
	ArtifactRoots                 []StorageArtifactRootGroup
}

// StorageRepositoryGroup reports Runs and their artifact bytes for one
// recorded repository root. Missing stays true when the recorded path no
// longer exists; the group is never omitted for that reason.
type StorageRepositoryGroup struct {
	Repository string
	Missing    bool
	Rows       int64
	Bytes      int64
}

// StorageStateGroup reports Runs and their artifact bytes for one Run state.
type StorageStateGroup struct {
	State string
	Rows  int64
	Bytes int64
}

// StorageTableGroup reports the row count and allocated SQLite bytes for one
// durable table. Index pages are attributed to their owning table.
type StorageTableGroup struct {
	Table string
	Rows  int64
	Bytes int64
}

// StorageArtifactRootGroup reports one unique recorded Artifact Root. Rows is
// the number of Runs that record the root; Bytes measures regular files below
// the root once, regardless of how many Runs reference it.
type StorageArtifactRootGroup struct {
	ArtifactRoot string
	Missing      bool
	Rows         int64
	Bytes        int64
}

type storageRunRecord struct {
	id           string
	repository   string
	state        string
	artifactRoot string
}

// StorageReport measures Run storage without starting a transaction, writing
// a pragma, migrating the schema, or changing any database or Artifact Root
// byte. Call it through OpenStorageReader for the byte-identical connection
// contract.
func (store *Store) StorageReport(ctx context.Context) (StorageReport, error) {
	databasePath, err := storageDatabasePath(ctx, store.db)
	if err != nil {
		return StorageReport{}, err
	}
	pageSize, err := storagePragmaInt64(ctx, store.db, "page_size")
	if err != nil {
		return StorageReport{}, err
	}
	freePages, err := storagePragmaInt64(ctx, store.db, "freelist_count")
	if err != nil {
		return StorageReport{}, err
	}
	tables, allocatedBytes, err := storageTableGroups(ctx, store.db)
	if err != nil {
		return StorageReport{}, err
	}
	runs, err := storageRunRecords(ctx, store.db)
	if err != nil {
		return StorageReport{}, err
	}

	repositories, states, artifactRoots, artifactBytes, err := storageFilesystemGroups(runs)
	if err != nil {
		return StorageReport{}, err
	}
	databaseInfo, err := os.Stat(databasePath)
	if err != nil {
		return StorageReport{}, fmt.Errorf("stat Run Database %q: %w", databasePath, err)
	}
	databaseGroupedBytes := allocatedBytes + freePages*pageSize
	databaseDifference := databaseInfo.Size() - databaseGroupedBytes
	if databaseDifference < 0 {
		databaseDifference = -databaseDifference
	}
	if databaseDifference > pageSize {
		return StorageReport{}, fmt.Errorf(
			"reconcile Run Database storage groups: file bytes=%d grouped bytes=%d difference=%d tolerance=%d",
			databaseInfo.Size(),
			databaseGroupedBytes,
			databaseDifference,
			pageSize,
		)
	}

	return StorageReport{
		DatabasePath:                  databasePath,
		DatabaseBytes:                 databaseInfo.Size(),
		DatabaseAllocatedBytes:        allocatedBytes,
		DatabaseFreeBytes:             freePages * pageSize,
		ArtifactBytes:                 artifactBytes,
		ReconciliationToleranceBytes:  pageSize,
		ReconciliationToleranceReason: StorageReportReconciliationToleranceReason,
		Repositories:                  repositories,
		States:                        states,
		Tables:                        tables,
		ArtifactRoots:                 artifactRoots,
	}, nil
}

func storageDatabasePath(ctx context.Context, db queryContextRunner) (string, error) {
	rows, err := db.QueryContext(ctx, `PRAGMA database_list`)
	if err != nil {
		return "", fmt.Errorf("list Run Database files: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var sequence int
		var name string
		var path string
		if err := rows.Scan(&sequence, &name, &path); err != nil {
			return "", fmt.Errorf("scan Run Database file: %w", err)
		}
		if name == "main" {
			if strings.TrimSpace(path) == "" {
				return "", errors.New("read Run Database path: main database has no file")
			}
			return path, nil
		}
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("iterate Run Database files: %w", err)
	}
	return "", errors.New("read Run Database path: main database is missing")
}

func storagePragmaInt64(ctx context.Context, db runQuerier, name string) (int64, error) {
	var value int64
	if err := db.QueryRowContext(ctx, `PRAGMA `+name).Scan(&value); err != nil {
		return 0, fmt.Errorf("read Run Database %s: %w", name, err)
	}
	return value, nil
}

func storageTableGroups(ctx context.Context, db *sql.DB) ([]StorageTableGroup, int64, error) {
	allocated := map[string]int64{}
	rows, err := db.QueryContext(ctx, `
SELECT COALESCE(NULLIF(schema_entry.tbl_name, ''), page.name) AS table_name,
       COALESCE(SUM(page.pgsize), 0)
FROM dbstat AS page
LEFT JOIN sqlite_schema AS schema_entry ON schema_entry.name = page.name
GROUP BY table_name
ORDER BY table_name`)
	if err != nil {
		return nil, 0, fmt.Errorf("measure Run Database table pages: %w", err)
	}
	for rows.Next() {
		var table string
		var bytes int64
		if err := rows.Scan(&table, &bytes); err != nil {
			_ = rows.Close()
			return nil, 0, fmt.Errorf("scan Run Database table pages: %w", err)
		}
		allocated[table] = bytes
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, 0, fmt.Errorf("iterate Run Database table pages: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, 0, fmt.Errorf("close Run Database table page rows: %w", err)
	}

	tableNames, err := storageTableNames(ctx, db)
	if err != nil {
		return nil, 0, err
	}
	groups := make([]StorageTableGroup, 0, len(tableNames))
	var allocatedBytes int64
	for _, table := range tableNames {
		var count int64
		query := `SELECT COUNT(*) FROM "` + strings.ReplaceAll(table, `"`, `""`) + `"`
		if err := db.QueryRowContext(ctx, query).Scan(&count); err != nil {
			return nil, 0, fmt.Errorf("count Run Database table %q: %w", table, err)
		}
		bytes := allocated[table]
		allocatedBytes += bytes
		groups = append(groups, StorageTableGroup{Table: table, Rows: count, Bytes: bytes})
	}
	return groups, allocatedBytes, nil
}

func storageTableNames(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, `
SELECT name
FROM sqlite_schema
WHERE type = 'table'
UNION SELECT 'sqlite_schema'
ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list Run Database tables: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	names := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan Run Database table: %w", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Run Database tables: %w", err)
	}
	return names, nil
}

func storageRunRecords(ctx context.Context, db *sql.DB) ([]storageRunRecord, error) {
	rows, err := db.QueryContext(ctx, `
SELECT id, git_root, state, artifact_dir
FROM runs
ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list Runs for storage report: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	runs := []storageRunRecord{}
	for rows.Next() {
		var run storageRunRecord
		if err := rows.Scan(&run.id, &run.repository, &run.state, &run.artifactRoot); err != nil {
			return nil, fmt.Errorf("scan Run for storage report: %w", err)
		}
		run.repository = filepath.Clean(run.repository)
		run.artifactRoot = filepath.Clean(run.artifactRoot)
		runs = append(runs, run)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Runs for storage report: %w", err)
	}
	return runs, nil
}

func storageFilesystemGroups(runs []storageRunRecord) ([]StorageRepositoryGroup, []StorageStateGroup, []StorageArtifactRootGroup, int64, error) {
	repositories := map[string]*StorageRepositoryGroup{}
	states := map[string]*StorageStateGroup{}
	artifactRoots := map[string]*StorageArtifactRootGroup{}

	for _, run := range runs {
		repository := repositories[run.repository]
		if repository == nil {
			exists, err := storagePathExists(run.repository)
			if err != nil {
				return nil, nil, nil, 0, fmt.Errorf("inspect recorded repository %q: %w", run.repository, err)
			}
			repository = &StorageRepositoryGroup{Repository: run.repository, Missing: !exists}
			repositories[run.repository] = repository
		}
		state := states[run.state]
		if state == nil {
			state = &StorageStateGroup{State: run.state}
			states[run.state] = state
		}
		root := artifactRoots[run.artifactRoot]
		if root == nil {
			root = &StorageArtifactRootGroup{ArtifactRoot: run.artifactRoot}
			artifactRoots[run.artifactRoot] = root
		}

		runBytes, _, err := storageRunArtifactBytes(run.artifactRoot, run.id)
		if err != nil {
			return nil, nil, nil, 0, err
		}
		repository.Rows++
		repository.Bytes += runBytes
		state.Rows++
		state.Bytes += runBytes
		root.Rows++
	}

	var artifactBytes int64
	for _, root := range artifactRoots {
		bytes, exists, err := storageDirectoryBytes(root.ArtifactRoot)
		if err != nil {
			return nil, nil, nil, 0, fmt.Errorf("measure Artifact Root %q: %w", root.ArtifactRoot, err)
		}
		root.Missing = !exists
		root.Bytes = bytes
		artifactBytes += bytes
	}

	return sortedStorageRepositories(repositories), sortedStorageStates(states), sortedStorageArtifactRoots(artifactRoots), artifactBytes, nil
}

func storagePathExists(path string) (bool, error) {
	if path == "" || path == "." {
		return false, nil
	}
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func storageRunArtifactBytes(artifactRoot string, runID string) (int64, bool, error) {
	if strings.TrimSpace(artifactRoot) == "" || artifactRoot == "." {
		return 0, false, nil
	}
	if filepath.Base(runID) != runID || runID == "." || runID == ".." {
		return 0, false, fmt.Errorf("measure Run artifacts: unsafe Run ID %q", runID)
	}
	return storageDirectoryBytes(filepath.Join(artifactRoot, "runs", runID))
}

func storageDirectoryBytes(path string) (int64, bool, error) {
	if strings.TrimSpace(path) == "" || path == "." {
		return 0, false, nil
	}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	if !info.IsDir() {
		return 0, true, fmt.Errorf("path is not a directory")
	}
	walkRoot := path
	pathInfo, err := os.Lstat(path)
	if err != nil {
		return 0, true, err
	}
	if pathInfo.Mode()&os.ModeSymlink != 0 {
		walkRoot, err = filepath.EvalSymlinks(path)
		if err != nil {
			return 0, true, err
		}
	}
	var bytes int64
	err = filepath.WalkDir(walkRoot, func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		entryInfo, err := entry.Info()
		if err != nil {
			return err
		}
		if entryInfo.Mode().IsRegular() {
			bytes += entryInfo.Size()
		}
		return nil
	})
	if err != nil {
		return 0, true, err
	}
	return bytes, true, nil
}

func sortedStorageRepositories(groups map[string]*StorageRepositoryGroup) []StorageRepositoryGroup {
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]StorageRepositoryGroup, 0, len(keys))
	for _, key := range keys {
		result = append(result, *groups[key])
	}
	return result
}

func sortedStorageStates(groups map[string]*StorageStateGroup) []StorageStateGroup {
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]StorageStateGroup, 0, len(keys))
	for _, key := range keys {
		result = append(result, *groups[key])
	}
	return result
}

func sortedStorageArtifactRoots(groups map[string]*StorageArtifactRootGroup) []StorageArtifactRootGroup {
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]StorageArtifactRootGroup, 0, len(keys))
	for _, key := range keys {
		result = append(result, *groups[key])
	}
	return result
}

// AppendRunEvent persists one Run Event, allocating the next per-Run cursor
// inside the insert transaction. Appending to a missing Run fails clearly.
func (store *Store) AppendRunEvent(ctx context.Context, event runevent.RunEvent) (int64, error) {
	cursors, err := store.AppendRunEvents(ctx, []runevent.RunEvent{event})
	if err != nil {
		return 0, err
	}
	return cursors[0], nil
}

// AppendRunEvents persists a small ordered batch of Run Events in one
// immediate transaction and returns the allocated cursors in input order.
func (store *Store) AppendRunEvents(ctx context.Context, events []runevent.RunEvent) ([]int64, error) {
	if len(events) == 0 {
		return nil, nil
	}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin Run Event append: %w", err)
	}
	defer rollbackUnlessCommitted(tx)

	cursors := make([]int64, 0, len(events))
	for _, event := range events {
		cursor, err := appendRunEvent(ctx, tx, event)
		if err != nil {
			return nil, err
		}
		cursors = append(cursors, cursor)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit Run Event append: %w", err)
	}
	return cursors, nil
}

// PruneTerminalRuns deletes Run Event Journal rows for terminal Runs completed
// before cutoff. It never deletes Run rows or Active Run locks.
func (store *Store) PruneTerminalRuns(ctx context.Context, cutoff time.Time) (PruneResult, error) {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return PruneResult{}, fmt.Errorf("begin Run Event prune: %w", err)
	}
	defer rollbackUnlessCommitted(tx)

	candidates, err := terminalRunPruneCandidates(ctx, tx, cutoff)
	if err != nil {
		return PruneResult{}, err
	}
	runIDs := pruneCandidateRunIDs(candidates)
	if len(runIDs) == 0 {
		if err := tx.Commit(); err != nil {
			return PruneResult{}, fmt.Errorf("commit Run Event prune: %w", err)
		}
		return PruneResult{}, nil
	}

	deleted, err := deleteRunEventsForRuns(ctx, tx, runIDs)
	if err != nil {
		return PruneResult{}, err
	}
	if err := tx.Commit(); err != nil {
		return PruneResult{}, fmt.Errorf("commit Run Event prune: %w", err)
	}
	return PruneResult{RunIDs: runIDs, Events: deleted}, nil
}

// PreviewCompaction measures a compact SQLite snapshot without changing the
// original Run Database. Compact accepts only a preview whose private
// fingerprint still matches after it acquires exclusive database access.
func (store *Store) PreviewCompaction(ctx context.Context) (preview CompactionPreview, err error) {
	conn, err := store.db.Conn(ctx)
	if err != nil {
		return CompactionPreview{}, fmt.Errorf("preview Run Database compaction: acquire connection: %w", err)
	}
	defer func() {
		err = errors.Join(err, conn.Close())
	}()

	before, err := readCompactionFingerprint(ctx, conn)
	if err != nil {
		return CompactionPreview{}, fmt.Errorf("preview Run Database compaction: %w", err)
	}
	bytesBefore, err := checkedMultiply(before.pageCount, before.pageSize)
	if err != nil {
		return CompactionPreview{}, fmt.Errorf("preview Run Database compaction: %w", err)
	}
	if err := store.requireTemporaryCapacity(before.databasePath, bytesBefore); err != nil {
		return CompactionPreview{}, err
	}
	bytesAfter, err := vacuumIntoSize(ctx, conn, before.databasePath)
	if err != nil {
		return CompactionPreview{}, fmt.Errorf("preview Run Database compaction: %w", err)
	}
	after, err := readCompactionFingerprint(ctx, conn)
	if err != nil {
		return CompactionPreview{}, fmt.Errorf("preview Run Database compaction: recheck writer state: %w", err)
	}
	if before != after {
		return CompactionPreview{}, WriterPresentCompactionError{
			Cause: errors.New("Run Database changed while the compact snapshot was measured"),
		}
	}
	if bytesAfter > bytesBefore {
		return CompactionPreview{}, fmt.Errorf(
			"preview Run Database compaction: compact snapshot grew from %d to %d bytes",
			bytesBefore,
			bytesAfter,
		)
	}
	return CompactionPreview{
		BytesBefore:                   bytesBefore,
		BytesReclaimable:              bytesBefore - bytesAfter,
		BytesAfter:                    bytesAfter,
		ReconciliationToleranceBytes:  before.pageSize,
		ReconciliationToleranceReason: CompactionReconciliationToleranceReason,
		fingerprint:                   before,
		snapshotBytesAfter:            bytesAfter,
	}, nil
}

// Compact rebuilds the Run Database only after proving this is the sole
// SQLite connection, rejecting Active Runs, revalidating the preview, and
// measuring the documented worst-case temporary-space requirement. SQLite's
// transactional VACUUM keeps the original database usable on interruption.
func (store *Store) Compact(ctx context.Context, preview CompactionPreview) (result CompactionResult, err error) {
	if err := validateCompactionPreview(preview); err != nil {
		return CompactionResult{}, err
	}
	conn, err := store.db.Conn(ctx)
	if err != nil {
		return CompactionResult{}, fmt.Errorf("compact Run Database: acquire connection: %w", err)
	}
	defer func() {
		err = errors.Join(err, conn.Close())
	}()

	fence, err := acquireCompactionFence(ctx, conn)
	if err != nil {
		return CompactionResult{}, err
	}
	defer func() {
		err = errors.Join(err, fence.release(context.WithoutCancel(ctx)))
	}()

	activeRunID, activeState, found, err := activeRunForCompaction(ctx, conn)
	if err != nil {
		return CompactionResult{}, err
	}
	if found {
		return CompactionResult{}, ActiveRunCompactionError{RunID: activeRunID, State: activeState}
	}
	current, err := readCompactionFingerprint(ctx, conn)
	if err != nil {
		return CompactionResult{}, fmt.Errorf("compact Run Database: revalidate preview: %w", err)
	}
	if current != preview.fingerprint {
		return CompactionResult{}, CompactionPreviewStaleError{}
	}
	requiredCapacity, err := checkedMultiply(preview.BytesBefore, 2)
	if err != nil {
		return CompactionResult{}, fmt.Errorf("compact Run Database: measure temporary capacity: %w", err)
	}
	if err := store.requireTemporaryCapacity(current.databasePath, requiredCapacity); err != nil {
		return CompactionResult{}, err
	}

	if _, err := conn.ExecContext(ctx, `VACUUM`); err != nil {
		return CompactionResult{}, fmt.Errorf("compact Run Database: rebuild database: %w", err)
	}
	if err := checkpointCompactedDatabase(ctx, conn); err != nil {
		return CompactionResult{}, err
	}
	databaseInfo, err := os.Stat(current.databasePath)
	if err != nil {
		return CompactionResult{}, fmt.Errorf("compact Run Database: stat result %q: %w", current.databasePath, err)
	}
	result = CompactionResult{
		BytesBefore:    preview.BytesBefore,
		BytesReclaimed: preview.BytesBefore - databaseInfo.Size(),
		BytesAfter:     databaseInfo.Size(),
	}
	difference := result.BytesReclaimed - preview.BytesReclaimable
	if difference < 0 {
		difference = -difference
	}
	if difference > preview.ReconciliationToleranceBytes {
		return result, CompactionResultMismatchError{
			PreviewReclaimable: preview.BytesReclaimable,
			ResultReclaimed:    result.BytesReclaimed,
			ToleranceBytes:     preview.ReconciliationToleranceBytes,
		}
	}
	return result, nil
}

func validateCompactionPreview(preview CompactionPreview) error {
	if preview.fingerprint.databasePath == "" || preview.fingerprint.pageSize <= 0 {
		return errors.New("compact Run Database: a preview from PreviewCompaction is required")
	}
	if preview.BytesBefore < 0 || preview.BytesReclaimable < 0 || preview.BytesAfter < 0 {
		return errors.New("compact Run Database: preview bytes must not be negative")
	}
	if preview.BytesAfter != preview.BytesBefore-preview.BytesReclaimable {
		return errors.New("compact Run Database: preview byte relation is invalid")
	}
	expectedBefore, err := checkedMultiply(preview.fingerprint.pageCount, preview.fingerprint.pageSize)
	if err != nil {
		return fmt.Errorf("compact Run Database: preview byte relation is invalid: %w", err)
	}
	if preview.BytesBefore != expectedBefore || preview.BytesAfter != preview.snapshotBytesAfter {
		return errors.New("compact Run Database: preview measurements were changed after PreviewCompaction")
	}
	if preview.ReconciliationToleranceBytes != preview.fingerprint.pageSize || preview.ReconciliationToleranceReason != CompactionReconciliationToleranceReason {
		return errors.New("compact Run Database: preview tolerance is invalid")
	}
	return nil
}

func readCompactionFingerprint(ctx context.Context, conn *sql.Conn) (compactionFingerprint, error) {
	databasePath, err := storageDatabasePath(ctx, conn)
	if err != nil {
		return compactionFingerprint{}, err
	}
	pageSize, err := storagePragmaInt64(ctx, conn, "page_size")
	if err != nil {
		return compactionFingerprint{}, err
	}
	pageCount, err := storagePragmaInt64(ctx, conn, "page_count")
	if err != nil {
		return compactionFingerprint{}, err
	}
	freePages, err := storagePragmaInt64(ctx, conn, "freelist_count")
	if err != nil {
		return compactionFingerprint{}, err
	}
	dataVersion, err := storagePragmaInt64(ctx, conn, "data_version")
	if err != nil {
		return compactionFingerprint{}, err
	}
	var totalChanges int64
	if err := conn.QueryRowContext(ctx, `SELECT total_changes()`).Scan(&totalChanges); err != nil {
		return compactionFingerprint{}, fmt.Errorf("read Run Database total changes: %w", err)
	}
	return compactionFingerprint{
		databasePath: databasePath,
		pageSize:     pageSize,
		pageCount:    pageCount,
		freePages:    freePages,
		dataVersion:  dataVersion,
		totalChanges: totalChanges,
	}, nil
}

func (store *Store) requireTemporaryCapacity(databasePath string, requiredBytes int64) error {
	if store.temporaryCapacity == nil {
		return errors.New("compact Run Database: temporary-capacity measurement is unavailable")
	}
	availableBytes, err := store.temporaryCapacity(filepath.Dir(databasePath))
	if err != nil {
		return fmt.Errorf("compact Run Database: measure temporary capacity beside %q: %w", databasePath, err)
	}
	if availableBytes >= requiredBytes {
		return nil
	}
	return CompactionCapacityError{
		RequiredBytes:  requiredBytes,
		AvailableBytes: availableBytes,
		ShortfallBytes: requiredBytes - availableBytes,
	}
}

func vacuumIntoSize(ctx context.Context, conn *sql.Conn, databasePath string) (size int64, err error) {
	temporary, err := os.CreateTemp(filepath.Dir(databasePath), ".roundfix-compaction-preview-*")
	if err != nil {
		return 0, fmt.Errorf("create compact snapshot beside Run Database: %w", err)
	}
	temporaryPath := temporary.Name()
	defer func() {
		removeErr := os.Remove(temporaryPath)
		if removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
			err = errors.Join(err, fmt.Errorf("remove compact snapshot %q: %w", temporaryPath, removeErr))
		}
	}()
	if err := temporary.Close(); err != nil {
		return 0, fmt.Errorf("close empty compact snapshot %q: %w", temporaryPath, err)
	}
	if _, err := conn.ExecContext(ctx, `VACUUM INTO ?`, temporaryPath); err != nil {
		return 0, fmt.Errorf("build compact snapshot: %w", err)
	}
	info, err := os.Stat(temporaryPath)
	if err != nil {
		return 0, fmt.Errorf("stat compact snapshot %q: %w", temporaryPath, err)
	}
	return info.Size(), nil
}

type compactionFence struct {
	conn        *sql.Conn
	busyTimeout int64
}

func acquireCompactionFence(ctx context.Context, conn *sql.Conn) (*compactionFence, error) {
	busyTimeout, err := storagePragmaInt64(ctx, conn, "busy_timeout")
	if err != nil {
		return nil, fmt.Errorf("compact Run Database: read lock wait: %w", err)
	}
	fence := &compactionFence{conn: conn, busyTimeout: busyTimeout}
	if _, err := conn.ExecContext(ctx, `PRAGMA busy_timeout = 0`); err != nil {
		return nil, fmt.Errorf("compact Run Database: disable lock wait: %w", err)
	}
	var lockingMode string
	if err := conn.QueryRowContext(ctx, `PRAGMA main.locking_mode = EXCLUSIVE`).Scan(&lockingMode); err != nil {
		return nil, errors.Join(
			fmt.Errorf("compact Run Database: request exclusive locking mode: %w", err),
			fence.release(context.WithoutCancel(ctx)),
		)
	}
	if !strings.EqualFold(lockingMode, "exclusive") {
		return nil, errors.Join(
			fmt.Errorf("compact Run Database: SQLite returned locking mode %q", lockingMode),
			fence.release(context.WithoutCancel(ctx)),
		)
	}
	if _, err := conn.ExecContext(ctx, `BEGIN EXCLUSIVE`); err != nil {
		return nil, errors.Join(
			WriterPresentCompactionError{Cause: err},
			fence.release(context.WithoutCancel(ctx)),
		)
	}
	if _, err := conn.ExecContext(ctx, `ROLLBACK`); err != nil {
		return nil, errors.Join(
			fmt.Errorf("compact Run Database: confirm exclusive access: %w", err),
			fence.release(context.WithoutCancel(ctx)),
		)
	}
	return fence, nil
}

func (fence *compactionFence) release(ctx context.Context) error {
	var releaseErrors []error
	var lockingMode string
	if err := fence.conn.QueryRowContext(ctx, `PRAGMA main.locking_mode = NORMAL`).Scan(&lockingMode); err != nil {
		releaseErrors = append(releaseErrors, fmt.Errorf("restore Run Database locking mode: %w", err))
	} else if !strings.EqualFold(lockingMode, "normal") {
		releaseErrors = append(releaseErrors, fmt.Errorf("restore Run Database locking mode: SQLite returned %q", lockingMode))
	}
	var userVersion int64
	if err := fence.conn.QueryRowContext(ctx, `PRAGMA main.user_version`).Scan(&userVersion); err != nil {
		releaseErrors = append(releaseErrors, fmt.Errorf("release Run Database exclusive lock: %w", err))
	}
	if _, err := fence.conn.ExecContext(ctx, fmt.Sprintf(`PRAGMA busy_timeout = %d`, fence.busyTimeout)); err != nil {
		releaseErrors = append(releaseErrors, fmt.Errorf("restore Run Database lock wait: %w", err))
	}
	return errors.Join(releaseErrors...)
}

func activeRunForCompaction(ctx context.Context, conn *sql.Conn) (runID string, state string, found bool, err error) {
	rows, err := conn.QueryContext(ctx, `
SELECT id, state
FROM runs
ORDER BY created_at, id`)
	if err != nil {
		return "", "", false, fmt.Errorf("compact Run Database: list Active Runs: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var candidateRunID string
		var candidateState string
		if err := rows.Scan(&candidateRunID, &candidateState); err != nil {
			return "", "", false, fmt.Errorf("compact Run Database: scan Active Run: %w", err)
		}
		if !IsTerminalState(candidateState) {
			return candidateRunID, candidateState, true, nil
		}
	}
	if err := rows.Err(); err != nil {
		return "", "", false, fmt.Errorf("compact Run Database: iterate Active Runs: %w", err)
	}
	return "", "", false, nil
}

func checkpointCompactedDatabase(ctx context.Context, conn *sql.Conn) error {
	var busy int
	var logFrames int
	var checkpointedFrames int
	if err := conn.QueryRowContext(ctx, `PRAGMA wal_checkpoint(TRUNCATE)`).Scan(&busy, &logFrames, &checkpointedFrames); err != nil {
		return fmt.Errorf("compact Run Database: checkpoint compacted database: %w", err)
	}
	if busy != 0 {
		return fmt.Errorf(
			"compact Run Database: checkpoint compacted database remained busy: log=%d checkpointed=%d",
			logFrames,
			checkpointedFrames,
		)
	}
	return nil
}

func checkedMultiply(left int64, right int64) (int64, error) {
	if left < 0 || right < 0 {
		return 0, errors.New("byte count must not be negative")
	}
	if left != 0 && right > math.MaxInt64/left {
		return 0, fmt.Errorf("byte count overflow: %d * %d", left, right)
	}
	return left * right, nil
}

// TerminalRunPruneCandidates lists terminal Runs whose completed_at is before
// cutoff without mutating the Run Database.
func (store *Store) TerminalRunPruneCandidates(ctx context.Context, cutoff time.Time) ([]PruneCandidate, error) {
	return terminalRunPruneCandidates(ctx, store.db, cutoff)
}

type queryContextRunner interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func terminalRunPruneCandidates(ctx context.Context, querier queryContextRunner, cutoff time.Time) ([]PruneCandidate, error) {
	rows, err := querier.QueryContext(ctx, `
SELECT r.id, r.completed_at, COUNT(e.run_id)
FROM runs r
LEFT JOIN run_events e ON e.run_id = r.id
WHERE r.state IN (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
  AND TRIM(r.completed_at) <> ''
GROUP BY r.id, r.completed_at
ORDER BY r.completed_at, r.id`,
		StateFetched,
		StateStopped,
		StateClean,
		StateCleanUnverified,
		StateReviewSkipped,
		StateMaxRoundsReached,
		StateBudgetExceeded,
		StateTimedOut,
		StateFailed,
		StateCheckoutMoved,
		StateIntegrationPending,
		StateUnresolved,
	)
	if err != nil {
		return nil, fmt.Errorf("select terminal Runs for Run Event prune: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	candidates := []PruneCandidate{}
	for rows.Next() {
		var candidate PruneCandidate
		var completedAtRaw string
		if err := rows.Scan(&candidate.RunID, &completedAtRaw, &candidate.Events); err != nil {
			return nil, fmt.Errorf("scan terminal Run for Run Event prune: %w", err)
		}
		completedAt, err := parseTime(completedAtRaw)
		if err != nil {
			return nil, fmt.Errorf("read Run %q completion time for Run Event prune: %w", candidate.RunID, err)
		}
		if completedAt.Before(cutoff) {
			candidates = append(candidates, candidate)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate terminal Runs for Run Event prune: %w", err)
	}
	return candidates, nil
}

func pruneCandidateRunIDs(candidates []PruneCandidate) []string {
	runIDs := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		runIDs = append(runIDs, candidate.RunID)
	}
	return runIDs
}

func deleteRunEventsForRuns(ctx context.Context, tx *sql.Tx, runIDs []string) (int, error) {
	if len(runIDs) == 0 {
		return 0, nil
	}
	placeholders := strings.TrimRight(strings.Repeat("?,", len(runIDs)), ",")
	args := make([]any, 0, len(runIDs))
	for _, runID := range runIDs {
		args = append(args, runID)
	}
	result, err := tx.ExecContext(ctx, `DELETE FROM run_events WHERE run_id IN (`+placeholders+`)`, args...)
	if err != nil {
		return 0, fmt.Errorf("delete Run Events for terminal Runs: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("read Run Event prune result: %w", err)
	}
	return int(affected), nil
}

func appendRunEvent(ctx context.Context, tx *sql.Tx, event runevent.RunEvent) (int64, error) {
	runID := strings.TrimSpace(event.RunID)
	if runID == "" {
		return 0, errors.New("append Run Event: Run ID is required")
	}
	var exists int
	err := tx.QueryRowContext(ctx, `SELECT 1 FROM runs WHERE id = ?`, runID).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("append Run Event: Run %q does not exist", runID)
	}
	if err != nil {
		return 0, fmt.Errorf("append Run Event: check Run %q: %w", runID, err)
	}

	// MAX+1 inside the insert transaction is race-free under the
	// single-writer rule and never reuses a cursor (no pruning yet).
	row := tx.QueryRowContext(ctx, `
INSERT INTO run_events (
	run_id, cursor, batch, source, kind, review_issue,
	tool_id, tool_state, summary, created_at, payload
)
VALUES (?, (SELECT COALESCE(MAX(cursor), 0) + 1 FROM run_events WHERE run_id = ?), ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING cursor`,
		runID,
		runID,
		event.Batch,
		string(event.Source),
		string(event.Kind),
		event.ReviewIssue,
		event.ToolID,
		event.ToolState,
		event.Summary,
		formatTime(event.Time),
		string(event.Payload),
	)
	var cursor int64
	if err := row.Scan(&cursor); err != nil {
		return 0, fmt.Errorf("insert Run Event for Run %q: %w", runID, err)
	}
	return cursor, nil
}

// RunEventsAfter lists events for one Run with cursors strictly greater
// than the given cursor, oldest first, bounded by limit.
func (store *Store) RunEventsAfter(ctx context.Context, runID string, cursor int64, limit int) ([]JournalEvent, error) {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return nil, errors.New("list Run Events: Run ID is required")
	}
	if limit <= 0 {
		return nil, errors.New("list Run Events: a positive limit is required")
	}
	rows, err := store.db.QueryContext(ctx, `
SELECT cursor, batch, source, kind, review_issue, tool_id, tool_state, summary, created_at, payload
FROM run_events
WHERE run_id = ? AND cursor > ?
ORDER BY cursor ASC
LIMIT ?`,
		runID,
		cursor,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list Run Events for Run %q: %w", runID, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	events := []JournalEvent{}
	for rows.Next() {
		entry, err := scanJournalEvent(rows, runID)
		if err != nil {
			return nil, err
		}
		events = append(events, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Run Events for Run %q: %w", runID, err)
	}
	return events, nil
}

func scanJournalEvent(rows *sql.Rows, runID string) (JournalEvent, error) {
	var entry JournalEvent
	var source string
	var kind string
	var createdAt string
	var payload string
	if err := rows.Scan(
		&entry.Cursor,
		&entry.Event.Batch,
		&source,
		&kind,
		&entry.Event.ReviewIssue,
		&entry.Event.ToolID,
		&entry.Event.ToolState,
		&entry.Event.Summary,
		&createdAt,
		&payload,
	); err != nil {
		return JournalEvent{}, fmt.Errorf("scan Run Event for Run %q: %w", runID, err)
	}
	parsedAt, err := parseTime(createdAt)
	if err != nil {
		return JournalEvent{}, err
	}
	entry.Event.RunID = runID
	entry.Event.Source = runevent.Source(source)
	entry.Event.Kind = runevent.Kind(kind)
	entry.Event.Time = parsedAt
	if payload != "" {
		entry.Event.Payload = []byte(payload)
	}
	return entry, nil
}

// RunEventsBefore lists events for one Run with cursors strictly less than
// the given cursor, oldest first, bounded by limit: the limit events
// immediately before the cursor. Scrollback pages the journal backward with
// it, symmetric to RunEventsAfter and under the same short autocommit
// read discipline.
func (store *Store) RunEventsBefore(ctx context.Context, runID string, cursor int64, limit int) ([]JournalEvent, error) {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return nil, errors.New("list Run Events: Run ID is required")
	}
	if limit <= 0 {
		return nil, errors.New("list Run Events: a positive limit is required")
	}
	rows, err := store.db.QueryContext(ctx, `
SELECT cursor, batch, source, kind, review_issue, tool_id, tool_state, summary, created_at, payload
FROM run_events
WHERE run_id = ? AND cursor < ?
ORDER BY cursor DESC
LIMIT ?`,
		runID,
		cursor,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list Run Events for Run %q: %w", runID, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	events := []JournalEvent{}
	for rows.Next() {
		entry, err := scanJournalEvent(rows, runID)
		if err != nil {
			return nil, err
		}
		events = append(events, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Run Events for Run %q: %w", runID, err)
	}
	slices.Reverse(events)
	return events, nil
}

// JournalSink adapts the Run Database to the Run Event sink interface. It
// registers as a critical sink: an append failure after Run start must fail
// the Run, never be swallowed.
type JournalSink struct {
	Store *Store
}

func (sink JournalSink) Publish(ctx context.Context, event runevent.RunEvent) error {
	_, err := sink.Store.AppendRunEvent(ctx, event)
	return err
}

// DataVersion exposes SQLite's data_version for this connection: it changes
// when another connection commits, so pollers can detect new events without
// reading rows.
func (store *Store) DataVersion(ctx context.Context) (int64, error) {
	var version int64
	if err := store.db.QueryRowContext(ctx, `PRAGMA data_version`).Scan(&version); err != nil {
		return 0, fmt.Errorf("read Run Database data version: %w", err)
	}
	return version, nil
}
