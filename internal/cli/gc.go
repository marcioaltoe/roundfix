package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

var gcDeps = defaultGCDependencies()

type gcDependencies struct {
	loadConfig        func(roundconfig.LoadOptions) (roundconfig.Loaded, error)
	openStore         func(context.Context, string) (*store.Store, error)
	openStoreReader   func(context.Context, string) (*store.Store, error)
	openStorageReader func(context.Context, string) (*store.Store, error)
	now               func() time.Time
}

func defaultGCDependencies() gcDependencies {
	return gcDependencies{
		loadConfig:        roundconfig.Load,
		openStore:         store.Open,
		openStoreReader:   store.OpenReader,
		openStorageReader: store.OpenStorageReader,
		now:               func() time.Time { return time.Now().UTC() },
	}
}

type gcOptions struct {
	dryRun bool
}

type gcReport struct {
	DryRun        bool
	Skipped       bool
	Retention     time.Duration
	Cutoff        time.Time
	RunIDs        []string
	OrphanIDs     []string
	JournalRows   int
	ArtifactBytes int64
}

type gcArtifactDir struct {
	runID string
	path  string
	bytes int64
}

type retentionPruneReport struct {
	RunIDs        []string
	JournalRows   int
	ArtifactBytes int64
}

func pruneRunRetention(ctx context.Context, runStore *store.Store, artifactRoot string, cutoff time.Time) (retentionPruneReport, error) {
	candidates, err := runStore.TerminalRunPruneCandidates(ctx, cutoff)
	if err != nil {
		return retentionPruneReport{}, err
	}
	runDirs, err := gcCandidateArtifactDirs(artifactRoot, candidates)
	if err != nil {
		return retentionPruneReport{}, err
	}
	pruned, err := runStore.PruneTerminalRuns(ctx, cutoff)
	if err != nil {
		return retentionPruneReport{}, err
	}
	runDirs = gcFilterArtifactDirs(runDirs, pruned.RunIDs)
	bytes := gcArtifactBytes(runDirs)
	if err := gcRemoveArtifactDirs(runDirs); err != nil {
		return retentionPruneReport{}, err
	}
	return retentionPruneReport{
		RunIDs:        pruned.RunIDs,
		JournalRows:   pruned.Events,
		ArtifactBytes: bytes,
	}, nil
}

func sweepRunRetention(ctx context.Context, runStore *store.Store, artifactRoot string, retention time.Duration, stderr io.Writer) {
	if retention <= 0 {
		return
	}
	cutoff := time.Now().UTC().Add(-retention)
	report, err := pruneRunRetention(ctx, runStore, artifactRoot, cutoff)
	if err != nil {
		fmt.Fprintf(stderr, "%s: warning: Journal Retention prune failed: %v\n", app.Name, err)
		return
	}
	if len(report.RunIDs) == 0 && report.JournalRows == 0 && report.ArtifactBytes == 0 {
		return
	}
	fmt.Fprintf(stderr, "%s: pruned Run storage runs=%d journal_rows=%d artifact_bytes=%d\n",
		app.Name,
		len(report.RunIDs),
		report.JournalRows,
		report.ArtifactBytes,
	)
}

func runGCCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	dependencies := commandDependenciesForContext(ctx).gc
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("gc"))
		return exitOK
	}
	opts, err := parseGCCommand(args)
	if err != nil {
		printPreflightFailure("gc", err, stderr)
		return exitPreflight
	}
	loadOptions, err := environment.loadOptions(stderr)
	if err != nil {
		printPreflightFailure("gc", err, stderr)
		return exitPreflight
	}
	loaded, err := dependencies.loadConfig(loadOptions)
	if err != nil {
		printPreflightFailure("gc", err, stderr)
		return exitPreflight
	}

	report, err := runGC(ctx, opts, loaded)
	if err != nil {
		fmt.Fprintf(stderr, "%s: gc failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	printGCReport(stdout, report)
	return exitOK
}

func runStorageCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("storage"))
		return exitOK
	}
	if len(args) == 0 || args[0] != "report" {
		printPreflightFailure("storage", validationError{message: "expected subcommand \"report\""}, stderr)
		return exitPreflight
	}
	if commandWantsHelp(args[1:]) {
		fmt.Fprint(stdout, commandUsage("storage"))
		return exitOK
	}
	if len(args) != 1 {
		printPreflightFailure("storage report", validationError{message: "does not accept flags or arguments"}, stderr)
		return exitPreflight
	}
	if environment.homeDirErr != nil {
		printPreflightFailure("storage report", fmt.Errorf("resolve Roundfix Home: %w", environment.homeDirErr), stderr)
		return exitPreflight
	}

	report, err := measureStorageReport(ctx, environment.homeDir)
	if err != nil {
		fmt.Fprintf(stderr, "%s: storage report failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	printStorageReport(stdout, report)
	return exitOK
}

func measureStorageReport(ctx context.Context, homeDir string) (store.StorageReport, error) {
	databasePath := store.DatabasePath(homeDir)
	_, err := os.Stat(databasePath)
	if errors.Is(err, os.ErrNotExist) {
		return store.StorageReport{DatabasePath: databasePath}, nil
	}
	if err != nil {
		return store.StorageReport{}, fmt.Errorf("stat Run Database %q: %w", databasePath, err)
	}

	runStore, err := commandDependenciesForContext(ctx).gc.openStorageReader(ctx, homeDir)
	if err != nil {
		return store.StorageReport{}, err
	}
	report, reportErr := runStore.StorageReport(ctx)
	closeErr := runStore.Close()
	if reportErr != nil {
		if closeErr != nil {
			return store.StorageReport{}, errors.Join(reportErr, fmt.Errorf("close Run Database reader: %w", closeErr))
		}
		return store.StorageReport{}, reportErr
	}
	if closeErr != nil {
		return store.StorageReport{}, fmt.Errorf("close Run Database reader: %w", closeErr)
	}
	return report, nil
}

func printStorageReport(stdout io.Writer, report store.StorageReport) {
	fmt.Fprintln(stdout, "Storage report")
	fmt.Fprintln(stdout, "Database")
	fmt.Fprintf(stdout, "  Path: %s\n", report.DatabasePath)
	fmt.Fprintf(stdout, "  File bytes: %d\n", report.DatabaseBytes)
	fmt.Fprintf(stdout, "  Table bytes: %d\n", report.DatabaseAllocatedBytes)
	fmt.Fprintf(stdout, "  Free bytes: %d\n", report.DatabaseFreeBytes)
	databaseGroupedBytes := report.DatabaseAllocatedBytes + report.DatabaseFreeBytes
	fmt.Fprintf(stdout, "  Grouped bytes: %d\n", databaseGroupedBytes)
	fmt.Fprintf(stdout, "  Reconciliation difference: %d\n", storageByteDifference(report.DatabaseBytes, databaseGroupedBytes))
	if report.ReconciliationToleranceReason == "" {
		fmt.Fprintln(stdout, "  Reconciliation tolerance: 0 bytes (no Run Database exists)")
	} else {
		fmt.Fprintf(stdout, "  Reconciliation tolerance: %d bytes (%s)\n", report.ReconciliationToleranceBytes, report.ReconciliationToleranceReason)
	}

	fmt.Fprintln(stdout, "Tables")
	fmt.Fprintln(stdout, "  TABLE\tROWS\tBYTES")
	for _, group := range report.Tables {
		fmt.Fprintf(stdout, "  %s\t%d\t%d\n", group.Table, group.Rows, group.Bytes)
	}

	fmt.Fprintln(stdout, "Repositories")
	fmt.Fprintln(stdout, "  REPOSITORY\tSTATUS\tRUN_ROWS\tRUN_ARTIFACT_BYTES")
	for _, group := range report.Repositories {
		fmt.Fprintf(stdout, "  %s\t%s\t%d\t%d\n", group.Repository, storagePathStatus(group.Missing), group.Rows, group.Bytes)
	}

	fmt.Fprintln(stdout, "States")
	fmt.Fprintln(stdout, "  STATE\tRUN_ROWS\tRUN_ARTIFACT_BYTES")
	for _, group := range report.States {
		fmt.Fprintf(stdout, "  %s\t%d\t%d\n", group.State, group.Rows, group.Bytes)
	}

	fmt.Fprintln(stdout, "Artifact Roots")
	fmt.Fprintln(stdout, "  ARTIFACT_ROOT\tSTATUS\tRUN_ROWS\tBYTES")
	var artifactGroupedBytes int64
	for _, group := range report.ArtifactRoots {
		fmt.Fprintf(stdout, "  %s\t%s\t%d\t%d\n", group.ArtifactRoot, storagePathStatus(group.Missing), group.Rows, group.Bytes)
		artifactGroupedBytes += group.Bytes
	}
	fmt.Fprintf(stdout, "  Measured bytes: %d\n", report.ArtifactBytes)
	fmt.Fprintf(stdout, "  Grouped bytes: %d\n", artifactGroupedBytes)
	fmt.Fprintf(stdout, "  Reconciliation difference: %d\n", storageByteDifference(report.ArtifactBytes, artifactGroupedBytes))
}

func storagePathStatus(missing bool) string {
	if missing {
		return "missing"
	}
	return "present"
}

func storageByteDifference(left int64, right int64) int64 {
	difference := left - right
	if difference < 0 {
		return -difference
	}
	return difference
}

func parseGCCommand(args []string) (gcOptions, error) {
	fs := flag.NewFlagSet("gc", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	dryRun := fs.Bool("dry-run", false, "Preview pruning without deleting anything")
	if err := fs.Parse(args); err != nil {
		return gcOptions{}, validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return gcOptions{}, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	return gcOptions{dryRun: *dryRun}, nil
}

func runGC(ctx context.Context, opts gcOptions, loaded roundconfig.Loaded) (gcReport, error) {
	dependencies := commandDependenciesForContext(ctx).gc
	retention := loaded.Config.Store.JournalRetention
	report := gcReport{
		DryRun:    opts.dryRun,
		Retention: retention,
	}
	if retention == 0 {
		report.Skipped = true
		return report, nil
	}

	artifactRoot, err := gcArtifactRoot(opts.dryRun, loaded)
	if err != nil {
		return gcReport{}, err
	}
	cutoff := dependencies.now().UTC().Add(-retention)
	report.Cutoff = cutoff

	dbExists, err := gcDatabaseExists(loaded.HomeDir)
	if err != nil {
		return gcReport{}, err
	}

	runIDSet := map[string]struct{}{}
	var candidates []store.PruneCandidate
	var runStore *store.Store
	if dbExists {
		runStore, err = openGCStore(ctx, opts.dryRun, loaded.HomeDir)
		if err != nil {
			return gcReport{}, err
		}
		defer func() {
			_ = runStore.Close()
		}()

		candidates, runIDSet, err = gcStoreState(ctx, runStore, cutoff)
		if err != nil {
			return gcReport{}, err
		}
	}

	runDirs, err := gcCandidateArtifactDirs(artifactRoot, candidates)
	if err != nil {
		return gcReport{}, err
	}
	orphanDirs, err := gcOrphanArtifactDirs(artifactRoot, runIDSet)
	if err != nil {
		return gcReport{}, err
	}
	report.RunIDs = gcCandidateRunIDs(candidates)
	report.OrphanIDs = gcArtifactRunIDs(orphanDirs)
	report.JournalRows = gcCandidateEventCount(candidates)
	report.ArtifactBytes = gcArtifactBytes(runDirs) + gcArtifactBytes(orphanDirs)

	if opts.dryRun {
		return report, nil
	}

	if runStore != nil {
		pruned, err := pruneRunRetention(ctx, runStore, artifactRoot, cutoff)
		if err != nil {
			return gcReport{}, err
		}
		report.RunIDs = pruned.RunIDs
		report.JournalRows = pruned.JournalRows
		report.ArtifactBytes = pruned.ArtifactBytes + gcArtifactBytes(orphanDirs)
	}
	if err := gcRemoveArtifactDirs(orphanDirs); err != nil {
		return gcReport{}, err
	}
	return report, nil
}

func gcArtifactRoot(dryRun bool, loaded roundconfig.Loaded) (string, error) {
	if dryRun {
		return roundconfig.ResolveArtifactDirectory(loaded.Config.Defaults.ArtifactDir, loaded.GitRoot, loaded.HomeDir)
	}
	return roundconfig.ValidateArtifactDirectory(loaded.Config.Defaults.ArtifactDir, loaded.GitRoot, loaded.HomeDir)
}

func gcDatabaseExists(homeDir string) (bool, error) {
	_, err := os.Stat(store.DatabasePath(homeDir))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("stat Run Database: %w", err)
	}
	return true, nil
}

func openGCStore(ctx context.Context, dryRun bool, homeDir string) (*store.Store, error) {
	dependencies := commandDependenciesForContext(ctx).gc
	if dryRun {
		return dependencies.openStoreReader(ctx, homeDir)
	}
	return dependencies.openStore(ctx, homeDir)
}

func gcStoreState(ctx context.Context, runStore *store.Store, cutoff time.Time) ([]store.PruneCandidate, map[string]struct{}, error) {
	candidates, err := runStore.TerminalRunPruneCandidates(ctx, cutoff)
	if err != nil {
		return nil, nil, err
	}
	runIDs, err := runStore.RunIDs(ctx)
	if err != nil {
		return nil, nil, err
	}
	runIDSet := make(map[string]struct{}, len(runIDs))
	for _, runID := range runIDs {
		runIDSet[runID] = struct{}{}
	}
	return candidates, runIDSet, nil
}

func gcCandidateArtifactDirs(artifactRoot string, candidates []store.PruneCandidate) ([]gcArtifactDir, error) {
	dirs := make([]gcArtifactDir, 0, len(candidates))
	for _, candidate := range candidates {
		dir, err := gcRunArtifactDir(artifactRoot, candidate.RunID)
		if err != nil {
			return nil, err
		}
		if dir.path == "" {
			continue
		}
		dirs = append(dirs, dir)
	}
	return dirs, nil
}

func gcOrphanArtifactDirs(artifactRoot string, runIDs map[string]struct{}) ([]gcArtifactDir, error) {
	runsRoot := filepath.Join(filepath.Clean(artifactRoot), "runs")
	entries, err := os.ReadDir(runsRoot)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read Run artifact root %q: %w", runsRoot, err)
	}

	orphanDirs := []gcArtifactDir{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		runID := entry.Name()
		if _, ok := runIDs[runID]; ok {
			continue
		}
		dir, err := gcRunArtifactDir(artifactRoot, runID)
		if err != nil {
			return nil, err
		}
		if dir.path == "" {
			continue
		}
		orphanDirs = append(orphanDirs, dir)
	}
	return orphanDirs, nil
}

func gcRunArtifactDir(artifactRoot string, runID string) (gcArtifactDir, error) {
	path, err := gcRunArtifactPath(artifactRoot, runID)
	if err != nil {
		return gcArtifactDir{}, err
	}
	bytes, exists, err := gcDirectoryBytes(path)
	if err != nil {
		return gcArtifactDir{}, err
	}
	if !exists {
		return gcArtifactDir{}, nil
	}
	return gcArtifactDir{runID: strings.TrimSpace(runID), path: path, bytes: bytes}, nil
}

func gcRunArtifactPath(artifactRoot string, runID string) (string, error) {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return "", errors.New("Run artifact cleanup requires a Run ID")
	}
	if filepath.Base(runID) != runID || runID == "." || runID == ".." {
		return "", fmt.Errorf("Run artifact cleanup refuses unsafe Run ID %q", runID)
	}
	runsRoot := filepath.Join(filepath.Clean(artifactRoot), "runs")
	path := filepath.Join(runsRoot, runID)
	rel, err := filepath.Rel(runsRoot, path)
	if err != nil {
		return "", fmt.Errorf("resolve Run artifact path %q: %w", path, err)
	}
	if rel == "." || rel == ".." || filepath.IsAbs(rel) || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("Run artifact cleanup refuses path outside %q: %q", runsRoot, path)
	}
	return path, nil
}

func gcDirectoryBytes(path string) (int64, bool, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("stat Run artifact directory %q: %w", path, err)
	}
	if !info.IsDir() {
		return 0, false, fmt.Errorf("Run artifact path %q is not a directory", path)
	}

	var bytes int64
	err = filepath.WalkDir(path, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			bytes += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0, false, fmt.Errorf("measure Run artifact directory %q: %w", path, err)
	}
	return bytes, true, nil
}

func gcRemoveArtifactDirs(dirs []gcArtifactDir) error {
	for _, dir := range dirs {
		info, err := os.Lstat(dir.path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return fmt.Errorf("stat Run artifact directory %q: %w", dir.path, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("Run artifact path %q is not a directory", dir.path)
		}
		if err := os.RemoveAll(dir.path); err != nil {
			return fmt.Errorf("remove Run artifact directory %q: %w", dir.path, err)
		}
	}
	return nil
}

func gcCandidateRunIDs(candidates []store.PruneCandidate) []string {
	runIDs := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		runIDs = append(runIDs, candidate.RunID)
	}
	return runIDs
}

func gcCandidateEventCount(candidates []store.PruneCandidate) int {
	count := 0
	for _, candidate := range candidates {
		count += candidate.Events
	}
	return count
}

func gcArtifactRunIDs(dirs []gcArtifactDir) []string {
	runIDs := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		runIDs = append(runIDs, dir.runID)
	}
	sort.Strings(runIDs)
	return runIDs
}

func gcArtifactBytes(dirs []gcArtifactDir) int64 {
	var bytes int64
	for _, dir := range dirs {
		bytes += dir.bytes
	}
	return bytes
}

func gcFilterArtifactDirs(dirs []gcArtifactDir, runIDs []string) []gcArtifactDir {
	keep := make(map[string]struct{}, len(runIDs))
	for _, runID := range runIDs {
		keep[runID] = struct{}{}
	}
	filtered := make([]gcArtifactDir, 0, len(dirs))
	for _, dir := range dirs {
		if _, ok := keep[dir.runID]; ok {
			filtered = append(filtered, dir)
		}
	}
	return filtered
}

func printGCReport(stdout io.Writer, report gcReport) {
	if report.Skipped {
		fmt.Fprintln(stdout, "GC skipped")
		fmt.Fprintf(stdout, "  Journal Retention: %s\n", formatGCDuration(report.Retention))
		fmt.Fprintln(stdout, "  No pruning performed.")
		return
	}
	if report.DryRun {
		fmt.Fprintln(stdout, "GC dry-run")
		fmt.Fprintf(stdout, "  Journal Retention: %s\n", formatGCDuration(report.Retention))
		fmt.Fprintf(stdout, "  Cutoff: %s\n", report.Cutoff.Format(time.RFC3339Nano))
		fmt.Fprintf(stdout, "  Runs eligible: %d\n", len(report.RunIDs))
		fmt.Fprintf(stdout, "  Journal rows eligible: %d\n", report.JournalRows)
		fmt.Fprintf(stdout, "  Artifact bytes reclaimable: %d\n", report.ArtifactBytes)
		fmt.Fprintf(stdout, "  Orphan artifact dirs eligible: %d\n", len(report.OrphanIDs))
		printGCIDList(stdout, "Eligible Runs", report.RunIDs)
		printGCIDList(stdout, "Orphan artifact dirs", report.OrphanIDs)
		return
	}
	fmt.Fprintln(stdout, "GC complete")
	fmt.Fprintf(stdout, "  Journal Retention: %s\n", formatGCDuration(report.Retention))
	fmt.Fprintf(stdout, "  Cutoff: %s\n", report.Cutoff.Format(time.RFC3339Nano))
	fmt.Fprintf(stdout, "  Runs pruned: %d\n", len(report.RunIDs))
	fmt.Fprintf(stdout, "  Journal rows removed: %d\n", report.JournalRows)
	fmt.Fprintf(stdout, "  Artifact bytes reclaimed: %d\n", report.ArtifactBytes)
	fmt.Fprintf(stdout, "  Orphan artifact dirs removed: %d\n", len(report.OrphanIDs))
	printGCIDList(stdout, "Pruned Runs", report.RunIDs)
	printGCIDList(stdout, "Orphan artifact dirs", report.OrphanIDs)
}

func printGCIDList(stdout io.Writer, title string, ids []string) {
	if len(ids) == 0 {
		return
	}
	fmt.Fprintf(stdout, "  %s:\n", title)
	for _, id := range ids {
		fmt.Fprintf(stdout, "    %s\n", id)
	}
}

func formatGCDuration(duration time.Duration) string {
	if duration == 0 {
		return "0"
	}
	if duration%time.Hour == 0 {
		return fmt.Sprintf("%dh", int(duration/time.Hour))
	}
	if duration%time.Minute == 0 {
		return fmt.Sprintf("%dm", int(duration/time.Minute))
	}
	if duration%time.Second == 0 {
		return fmt.Sprintf("%ds", int(duration/time.Second))
	}
	return duration.String()
}
