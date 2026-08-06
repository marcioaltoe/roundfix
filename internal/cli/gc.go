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
	dryRun   bool
	sanitize bool
	apply    bool
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

type gcSanitationClassification string

const (
	gcSanitationActive      gcSanitationClassification = "active"
	gcSanitationOrphaned    gcSanitationClassification = "orphaned"
	gcSanitationMissing     gcSanitationClassification = "missing"
	gcSanitationOverridden  gcSanitationClassification = "overridden"
	gcSanitationOutsideHome gcSanitationClassification = "outside Roundfix Home"
	gcSanitationUnsafe      gcSanitationClassification = "unsafe"
)

type gcSanitationCandidate struct {
	dir      gcArtifactDir
	evidence string
}

type gcSanitationPreservation struct {
	path   string
	reason string
}

type gcSanitationRootReport struct {
	path           string
	classification gcSanitationClassification
	evidence       string
	candidates     []gcSanitationCandidate
	preserved      []gcSanitationPreservation
}

type gcSanitationReport struct {
	apply         bool
	retention     time.Duration
	cutoff        time.Time
	roots         []gcSanitationRootReport
	directories   int
	artifactBytes int64
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

	if opts.sanitize {
		report, err := runGCSanitation(ctx, opts, loaded)
		if err != nil {
			fmt.Fprintf(stderr, "%s: gc sanitize failed: %v\n", app.Name, err)
			return exitRunFailed
		}
		printGCSanitationReport(stdout, report)
		return exitOK
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
	if len(args) > 0 && args[0] == "sanitize" {
		fs := flag.NewFlagSet("gc sanitize", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		apply := fs.Bool("apply", false, "Remove proven eligible and absent Run artifact directories")
		if err := fs.Parse(args[1:]); err != nil {
			return gcOptions{}, validationError{message: err.Error()}
		}
		if remaining := fs.Args(); len(remaining) > 0 {
			return gcOptions{}, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
		}
		return gcOptions{sanitize: true, apply: *apply}, nil
	}

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

func runGCSanitation(ctx context.Context, opts gcOptions, loaded roundconfig.Loaded) (gcSanitationReport, error) {
	report := gcSanitationReport{
		apply:     opts.apply,
		retention: loaded.Config.Store.JournalRetention,
	}
	if report.retention > 0 {
		report.cutoff = commandDependenciesForContext(ctx).gc.now().UTC().Add(-report.retention)
	}

	databaseExists, err := gcDatabaseExists(loaded.HomeDir)
	if err != nil {
		return gcSanitationReport{}, err
	}
	if !databaseExists {
		return report, nil
	}

	runStore, err := commandDependenciesForContext(ctx).gc.openStoreReader(ctx, loaded.HomeDir)
	if err != nil {
		return gcSanitationReport{}, err
	}
	roots, discoverErr := store.DiscoverArtifactRoots(ctx, runStore)
	closeErr := runStore.Close()
	if discoverErr != nil {
		if closeErr != nil {
			return gcSanitationReport{}, errors.Join(discoverErr, fmt.Errorf("close Run Database reader: %w", closeErr))
		}
		return gcSanitationReport{}, discoverErr
	}
	if closeErr != nil {
		return gcSanitationReport{}, fmt.Errorf("close Run Database reader: %w", closeErr)
	}

	allRunIDs := map[string]string{}
	for _, root := range roots {
		for _, run := range root.Runs {
			allRunIDs[run.ID] = root.Path
		}
	}
	for _, root := range roots {
		rootReport := classifyGCSanitationRoot(root, loaded.HomeDir)
		if rootReport.classification == gcSanitationActive || rootReport.classification == gcSanitationOrphaned {
			candidates, preserved, err := gcSanitationArtifactDirs(root, allRunIDs, report.cutoff, report.retention)
			if err != nil {
				rootReport.classification = gcSanitationUnsafe
				rootReport.evidence = fmt.Sprintf("preserved because directory ownership could not be proven: %v", err)
				rootReport.preserved = []gcSanitationPreservation{{path: root.Path, reason: rootReport.evidence}}
			} else {
				rootReport.candidates = candidates
				rootReport.preserved = preserved
			}
		}
		report.roots = append(report.roots, rootReport)
		for _, candidate := range rootReport.candidates {
			report.directories++
			report.artifactBytes += candidate.dir.bytes
		}
	}

	if !opts.apply {
		return report, nil
	}
	for _, root := range report.roots {
		for _, candidate := range root.candidates {
			if err := gcRemoveArtifactDirs([]gcArtifactDir{candidate.dir}); err != nil {
				return gcSanitationReport{}, err
			}
		}
	}
	return report, nil
}

func classifyGCSanitationRoot(root store.ArtifactRoot, homeDir string) gcSanitationRootReport {
	report := gcSanitationRootReport{path: root.Path}
	path := strings.TrimSpace(root.Path)
	if path == "" || !filepath.IsAbs(path) || filepath.Clean(path) != path {
		report.classification = gcSanitationUnsafe
		report.evidence = fmt.Sprintf("preserved because recorded Artifact Root %q is not a clean absolute path", root.Path)
		return report
	}
	for _, run := range root.Runs {
		if _, err := gcRunArtifactPath(path, run.ID); err != nil {
			report.classification = gcSanitationUnsafe
			report.evidence = fmt.Sprintf("preserved because Run %q has unsafe Artifact Directory evidence: %v", run.ID, err)
			return report
		}
	}

	roundfixHome := filepath.Dir(store.DatabasePath(homeDir))
	insideRoundfixHome, err := gcPathWithin(roundfixHome, path)
	if err != nil {
		report.classification = gcSanitationUnsafe
		report.evidence = fmt.Sprintf("preserved because Artifact Root containment could not be proven: %v", err)
		return report
	}
	if filepath.Clean(path) == filepath.Clean(roundfixHome) {
		report.classification = gcSanitationUnsafe
		report.evidence = fmt.Sprintf("preserved because Artifact Root equals Roundfix Home %q", roundfixHome)
		return report
	}
	if !insideRoundfixHome {
		report.classification = gcSanitationOutsideHome
		report.evidence = fmt.Sprintf("preserved because recorded Artifact Root is outside Roundfix Home %q", roundfixHome)
		return report
	}

	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		report.classification = gcSanitationMissing
		report.evidence = "preserved because the recorded Artifact Root does not exist"
		return report
	}
	if err != nil {
		report.classification = gcSanitationUnsafe
		report.evidence = fmt.Sprintf("preserved because the recorded Artifact Root cannot be inspected: %v", err)
		return report
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		report.classification = gcSanitationUnsafe
		report.evidence = "preserved because the recorded Artifact Root is not a physical directory"
		return report
	}
	runsRoot := filepath.Join(path, "runs")
	runsInfo, err := os.Lstat(runsRoot)
	if err == nil && (runsInfo.Mode()&os.ModeSymlink != 0 || !runsInfo.IsDir()) {
		report.classification = gcSanitationUnsafe
		report.evidence = fmt.Sprintf("preserved because Run artifact root %q is not a physical directory", runsRoot)
		return report
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		report.classification = gcSanitationUnsafe
		report.evidence = fmt.Sprintf("preserved because Run artifact root %q cannot be inspected: %v", runsRoot, err)
		return report
	}

	repositories := gcArtifactRootRepositories(root)
	for _, repository := range repositories {
		defaultRoot, err := roundconfig.ResolveArtifactDirectory("", repository, homeDir)
		if err != nil {
			report.classification = gcSanitationUnsafe
			report.evidence = fmt.Sprintf("preserved because the default Artifact Root for repository %q cannot be proven: %v", repository, err)
			return report
		}
		if path != defaultRoot {
			report.classification = gcSanitationOverridden
			report.evidence = fmt.Sprintf("preserved because recorded Artifact Root overrides default %q for repository %q", defaultRoot, repository)
			return report
		}
	}

	activeRunIDs := []string{}
	for _, run := range root.Runs {
		if !store.IsTerminalState(run.State) {
			activeRunIDs = append(activeRunIDs, run.ID)
		}
	}
	if len(activeRunIDs) > 0 {
		sort.Strings(activeRunIDs)
		report.classification = gcSanitationActive
		report.evidence = fmt.Sprintf("Active Runs record this Artifact Root: %s", strings.Join(activeRunIDs, ", "))
		return report
	}

	report.classification = gcSanitationOrphaned
	report.evidence = fmt.Sprintf("no Active Run records this Artifact Root; terminal Runs recorded: %d", len(root.Runs))
	return report
}

func gcPathWithin(root string, path string) (bool, error) {
	relative, err := filepath.Rel(filepath.Clean(root), filepath.Clean(path))
	if err != nil {
		return false, err
	}
	return relative == "." || (!filepath.IsAbs(relative) && relative != ".." && !strings.HasPrefix(relative, ".."+string(os.PathSeparator))), nil
}

func gcArtifactRootRepositories(root store.ArtifactRoot) []string {
	set := map[string]struct{}{}
	for _, run := range root.Runs {
		set[run.Repository] = struct{}{}
	}
	repositories := make([]string, 0, len(set))
	for repository := range set {
		repositories = append(repositories, repository)
	}
	sort.Strings(repositories)
	return repositories
}

func gcSanitationArtifactDirs(root store.ArtifactRoot, allRunIDs map[string]string, cutoff time.Time, retention time.Duration) ([]gcSanitationCandidate, []gcSanitationPreservation, error) {
	recordedRuns := map[string]store.ArtifactRootRun{}
	for _, run := range root.Runs {
		recordedRuns[run.ID] = run
	}
	runsRoot := filepath.Join(root.Path, "runs")
	entries, err := os.ReadDir(runsRoot)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("read Run artifact root %q: %w", runsRoot, err)
	}

	candidates := []gcSanitationCandidate{}
	preserved := []gcSanitationPreservation{}
	for _, entry := range entries {
		path := filepath.Join(runsRoot, entry.Name())
		info, err := os.Lstat(path)
		if err != nil {
			return nil, nil, fmt.Errorf("inspect Run artifact path %q: %w", path, err)
		}
		if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			preserved = append(preserved, gcSanitationPreservation{path: path, reason: "not a physical Run artifact directory"})
			continue
		}

		run, recordedHere := recordedRuns[entry.Name()]
		if recordedHere {
			if !gcSanitationRunEligible(run, cutoff, retention) {
				preserved = append(preserved, gcSanitationPreservation{path: path, reason: fmt.Sprintf("Run %s is not retention-eligible", run.ID)})
				continue
			}
			dir, err := gcRunArtifactDir(root.Path, run.ID)
			if err != nil {
				return nil, nil, err
			}
			candidates = append(candidates, gcSanitationCandidate{dir: dir, evidence: fmt.Sprintf("terminal Run %s completed before cutoff", run.ID)})
			continue
		}

		if recordedRoot, recorded := allRunIDs[entry.Name()]; recorded {
			preserved = append(preserved, gcSanitationPreservation{
				path:   path,
				reason: fmt.Sprintf("Run %s is recorded under Artifact Root %q", entry.Name(), recordedRoot),
			})
			continue
		}
		dir, err := gcRunArtifactDir(root.Path, entry.Name())
		if err != nil {
			return nil, nil, err
		}
		candidates = append(candidates, gcSanitationCandidate{dir: dir, evidence: fmt.Sprintf("Run %s is absent from durable Run metadata", entry.Name())})
	}
	return candidates, preserved, nil
}

func gcSanitationRunEligible(run store.ArtifactRootRun, cutoff time.Time, retention time.Duration) bool {
	return retention > 0 && store.IsTerminalState(run.State) && run.CompletedAt != nil && run.CompletedAt.Before(cutoff)
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

func printGCSanitationReport(stdout io.Writer, report gcSanitationReport) {
	if report.apply {
		fmt.Fprintln(stdout, "GC sanitation complete")
	} else {
		fmt.Fprintln(stdout, "GC sanitation dry-run")
	}
	fmt.Fprintf(stdout, "  Journal Retention: %s\n", formatGCDuration(report.retention))
	if !report.cutoff.IsZero() {
		fmt.Fprintf(stdout, "  Cutoff: %s\n", report.cutoff.Format(time.RFC3339Nano))
	}
	for _, root := range report.roots {
		fmt.Fprintf(stdout, "  Artifact Root: %s\n", root.path)
		fmt.Fprintf(stdout, "    Classification: %s\n", root.classification)
		fmt.Fprintf(stdout, "    Evidence: %s\n", root.evidence)
		for _, candidate := range root.candidates {
			action := "would remove"
			if report.apply {
				action = "removed"
			}
			fmt.Fprintf(stdout, "    %s: %s (%s)\n", action, candidate.dir.path, candidate.evidence)
		}
		for _, preserved := range root.preserved {
			fmt.Fprintf(stdout, "    preserved: %s (%s)\n", preserved.path, preserved.reason)
		}
	}
	if report.apply {
		fmt.Fprintf(stdout, "  Directories removed: %d\n", report.directories)
		fmt.Fprintf(stdout, "  Artifact bytes reclaimed: %d\n", report.artifactBytes)
		return
	}
	fmt.Fprintf(stdout, "  Directories reclaimable: %d\n", report.directories)
	fmt.Fprintf(stdout, "  Artifact bytes reclaimable: %d\n", report.artifactBytes)
	fmt.Fprintln(stdout, "  No directories removed; rerun with --apply to mutate.")
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
