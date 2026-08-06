package cli

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

func TestRunGCDryRunListsEligibleRunsAndChangesNothing(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir, repoDir := withCLIWorkspace(t)
	artifactDir := filepath.Join(homeDir, "artifacts")
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), fmt.Sprintf("defaults:\n  artifact_dir: %q\nstore:\n  journal_retention: 336h\n", artifactDir))
	now := time.Date(2026, 7, 20, 12, 0, 0, 0, time.UTC)
	fixture := seedGCFixture(t, ctx, homeDir, artifactDir, now)
	withGCNow(t, now)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, ctx, []string{"gc", "--dry-run"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected gc --dry-run exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected dry-run diagnostics to stay empty, got %q", stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{
		"GC dry-run",
		"Runs eligible: 1",
		"Journal rows eligible: 2",
		"Artifact bytes reclaimable: 16",
		fixture.oldRun.ID,
		fixture.orphanID,
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected dry-run output to contain %q, got %q", want, output)
		}
	}
	for _, notWant := range []string{fixture.activeRun.ID, fixture.recentRun.ID} {
		if strings.Contains(output, notWant) {
			t.Fatalf("dry-run output listed ineligible Run %s: %q", notWant, output)
		}
	}
	assertRunEvents(t, homeDir, fixture.oldRun.ID, 2)
	assertPathExists(t, fixture.oldArtifactDir)
	assertPathExists(t, fixture.orphanDir)
	assertPathExists(t, fixture.activeArtifactDir)
	assertPathExists(t, fixture.reviewArtifactPath)
}

func TestRunGCPrunesEligibleJournalsArtifactsAndOrphans(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir, repoDir := withCLIWorkspace(t)
	artifactDir := filepath.Join(homeDir, "artifacts")
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), fmt.Sprintf("defaults:\n  artifact_dir: %q\nstore:\n  journal_retention: 336h\n", artifactDir))
	now := time.Date(2026, 7, 20, 12, 0, 0, 0, time.UTC)
	fixture := seedGCFixture(t, ctx, homeDir, artifactDir, now)
	withGCNow(t, now)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, ctx, []string{"gc"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected gc exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected gc diagnostics to stay empty, got %q", stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{
		"GC complete",
		"Runs pruned: 1",
		"Journal rows removed: 2",
		"Artifact bytes reclaimed: 16",
		fixture.oldRun.ID,
		fixture.orphanID,
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected gc output to contain %q, got %q", want, output)
		}
	}
	assertRunEvents(t, homeDir, fixture.oldRun.ID, 0)
	assertRunEvents(t, homeDir, fixture.activeRun.ID, 1)
	assertRunEvents(t, homeDir, fixture.recentRun.ID, 1)
	assertRunExists(t, homeDir, fixture.oldRun.ID)
	assertRunExists(t, homeDir, fixture.activeRun.ID)
	assertRunExists(t, homeDir, fixture.recentRun.ID)
	assertPathMissing(t, fixture.oldArtifactDir)
	assertPathMissing(t, fixture.orphanDir)
	assertPathExists(t, fixture.activeArtifactDir)
	assertPathExists(t, fixture.recentArtifactDir)
	assertPathExists(t, fixture.reviewArtifactPath)
}

func TestRunGCSkipsWhenJournalRetentionIsZero(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir, repoDir := withCLIWorkspace(t)
	artifactDir := filepath.Join(homeDir, "artifacts")
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), fmt.Sprintf("defaults:\n  artifact_dir: %q\nstore:\n  journal_retention: 0\n", artifactDir))
	now := time.Date(2026, 7, 20, 12, 0, 0, 0, time.UTC)
	fixture := seedGCFixture(t, ctx, homeDir, artifactDir, now)
	withGCNow(t, now)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, ctx, []string{"gc"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected gc zero-retention exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected gc zero-retention diagnostics to stay empty, got %q", stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{"GC skipped", "Journal Retention: 0", "No pruning performed."} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected zero-retention output to contain %q, got %q", want, output)
		}
	}
	assertRunEvents(t, homeDir, fixture.oldRun.ID, 2)
	assertPathExists(t, fixture.oldArtifactDir)
	assertPathExists(t, fixture.orphanDir)
}

func TestRunGCHelp(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"gc", "--help"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected gc --help exit 0, got %d", code)
	}
	for _, want := range []string{
		"roundfix gc [--dry-run]",
		"--dry-run",
		"Journal Retention",
		"never deletes Run rows or Active Run locks",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("expected help output to contain %q, got %q", want, stdout.String())
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunGCSanitizeClassifiesEveryRecordedRootAndMutatesOnlyProvenDirectories(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir, repoDir := withCLIWorkspace(t)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), "store:\n  journal_retention: 336h\n")
	now := time.Date(2026, 7, 20, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)

	activeRepository := gcSanitationRepository(t, "active")
	orphanedRepository := gcSanitationRepository(t, "orphaned")
	missingRepository := gcSanitationRepository(t, "missing")
	overriddenRepository := gcSanitationRepository(t, "overridden")
	outsideRepository := gcSanitationRepository(t, "outside")
	unsafeRepository := gcSanitationRepository(t, "unsafe")
	activeRoot := resolveGCTestArtifactRoot(t, activeRepository, homeDir)
	orphanedRoot := resolveGCTestArtifactRoot(t, orphanedRepository, homeDir)
	missingRoot := resolveGCTestArtifactRoot(t, missingRepository, homeDir)
	overriddenRoot := filepath.Join(filepath.Dir(store.DatabasePath(homeDir)), "overrides", "artifacts")
	outsideRoot := filepath.Join(t.TempDir(), "outside-artifacts")
	unsafeRoot := filepath.Dir(store.DatabasePath(homeDir))

	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run store: %v", err)
	}
	activeRun := createGCSanitationRun(t, ctx, runStore, activeRepository, activeRoot, "active", false)
	eligibleRun := createGCSanitationRun(t, ctx, runStore, orphanedRepository, orphanedRoot, "eligible", true)
	recentRun := createGCSanitationRun(t, ctx, runStore, orphanedRepository, orphanedRoot, "recent", true)
	missingRun := createGCSanitationRun(t, ctx, runStore, missingRepository, missingRoot, "missing", true)
	overriddenRun := createGCSanitationRun(t, ctx, runStore, overriddenRepository, overriddenRoot, "overridden", true)
	outsideRun := createGCSanitationRun(t, ctx, runStore, outsideRepository, outsideRoot, "outside", true)
	_ = createGCSanitationRun(t, ctx, runStore, unsafeRepository, unsafeRoot, "unsafe", true)
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run store after sanitation seed: %v", err)
	}
	for _, runID := range []string{eligibleRun.ID, missingRun.ID, overriddenRun.ID, outsideRun.ID} {
		setRunTimestamps(t, homeDir, runID, now.Add(-400*time.Hour), now.Add(-400*time.Hour))
	}
	setRunTimestamps(t, homeDir, recentRun.ID, now.Add(-time.Hour), now.Add(-time.Hour))
	setRunTimestamps(t, homeDir, activeRun.ID, now.Add(-500*time.Hour), time.Time{})

	activeDir := writeRunArtifact(t, activeRoot, activeRun.ID, "active")
	eligibleDir := writeRunArtifact(t, orphanedRoot, eligibleRun.ID, "eligible")
	recentDir := writeRunArtifact(t, orphanedRoot, recentRun.ID, "recent")
	absentDir := writeRunArtifact(t, orphanedRoot, "run_absent_sanitation", "absent")
	overriddenDir := writeRunArtifact(t, overriddenRoot, overriddenRun.ID, "overridden")
	outsideDir := writeRunArtifact(t, outsideRoot, outsideRun.ID, "outside")
	reviewArtifact := filepath.Join(orphanedRoot, "reviews", "pr-123", "round-001", "issue.md")
	mustMkdir(t, filepath.Dir(reviewArtifact))
	mustWrite(t, reviewArtifact, "review artifact")
	unsafeMarker := filepath.Join(unsafeRoot, "unsafe-marker.txt")
	mustWrite(t, unsafeMarker, "preserve unsafe root")

	var dryRunStdout bytes.Buffer
	var dryRunStderr bytes.Buffer
	code := runCLIContext(t, ctx, []string{"gc", "sanitize"}, &dryRunStdout, &dryRunStderr)
	if code != exitOK {
		t.Fatalf("expected gc sanitize dry-run exit 0, got %d stderr=%q stdout=%q", code, dryRunStderr.String(), dryRunStdout.String())
	}
	if dryRunStderr.Len() != 0 {
		t.Fatalf("expected sanitation dry-run diagnostics to stay empty, got %q", dryRunStderr.String())
	}
	for _, want := range []string{
		"GC sanitation dry-run",
		"Classification: active",
		"Classification: orphaned",
		"Classification: missing",
		"Classification: overridden",
		"Classification: outside Roundfix Home",
		"Classification: unsafe",
		"Active Runs record this Artifact Root",
		"no Active Run records this Artifact Root",
		"recorded Artifact Root does not exist",
		"overrides default",
		"is outside Roundfix Home",
		"equals Roundfix Home",
		activeRoot,
		orphanedRoot,
	} {
		if !strings.Contains(dryRunStdout.String(), want) {
			t.Fatalf("expected sanitation dry-run output to contain %q, got %q", want, dryRunStdout.String())
		}
	}
	for _, path := range []string{activeDir, eligibleDir, recentDir, absentDir, overriddenDir, outsideDir, reviewArtifact, unsafeMarker} {
		assertPathExists(t, path)
	}

	var applyStdout bytes.Buffer
	var applyStderr bytes.Buffer
	code = runCLIContext(t, ctx, []string{"gc", "sanitize", "--apply"}, &applyStdout, &applyStderr)
	if code != exitOK {
		t.Fatalf("expected gc sanitize --apply exit 0, got %d stderr=%q stdout=%q", code, applyStderr.String(), applyStdout.String())
	}
	if applyStderr.Len() != 0 {
		t.Fatalf("expected sanitation apply diagnostics to stay empty, got %q", applyStderr.String())
	}
	firstRemoved := gcReportCount(t, applyStdout.String(), "Directories removed")
	firstBytes := gcReportCount(t, applyStdout.String(), "Artifact bytes reclaimed")
	if firstRemoved == 0 || firstBytes == 0 {
		t.Fatalf("expected first sanitation apply to remove proven directories, got %q", applyStdout.String())
	}
	assertPathMissing(t, eligibleDir)
	assertPathMissing(t, absentDir)
	for _, path := range []string{activeDir, recentDir, overriddenDir, outsideDir, reviewArtifact, unsafeMarker} {
		assertPathExists(t, path)
	}

	var secondStdout bytes.Buffer
	var secondStderr bytes.Buffer
	code = runCLIContext(t, ctx, []string{"gc", "sanitize", "--apply"}, &secondStdout, &secondStderr)
	if code != exitOK {
		t.Fatalf("expected second gc sanitize --apply exit 0, got %d stderr=%q stdout=%q", code, secondStderr.String(), secondStdout.String())
	}
	secondRemoved := gcReportCount(t, secondStdout.String(), "Directories removed")
	secondBytes := gcReportCount(t, secondStdout.String(), "Artifact bytes reclaimed")
	if secondRemoved != 0 || secondBytes != 0 || secondRemoved >= firstRemoved || secondBytes >= firstBytes {
		t.Fatalf("expected idempotent sanitation relation first=(directories=%d bytes=%d) second=(directories=0 bytes=0), second output=%q", firstRemoved, firstBytes, secondStdout.String())
	}
	assertPathExists(t, reviewArtifact)
	assertPathExists(t, unsafeMarker)
}

func TestRunStorageReportOutsideGitRepository(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	outsideGit := t.TempDir()
	missingRepository := filepath.Join(t.TempDir(), "removed-repository")
	missingArtifactRoot := filepath.Join(t.TempDir(), "removed-artifacts")
	setCommandEnvironmentForTest(t, homeDir, outsideGit)

	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run store: %v", err)
	}
	if _, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/storage-report",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        missingRepository,
		LocalBranch:    "feature/storage-report",
		HeadSHA:        "abc123",
		ArtifactDir:    missingArtifactRoot,
	}); err != nil {
		t.Fatalf("create storage report Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run store: %v", err)
	}
	databaseBefore, err := os.ReadFile(store.DatabasePath(homeDir))
	if err != nil {
		t.Fatalf("read Run Database before command: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, ctx, []string{"storage", "report"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected storage report exit 0 outside Git, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected storage report diagnostics to stay empty, got %q", stderr.String())
	}
	for _, want := range []string{
		"Storage report",
		"Tables",
		"Repositories",
		"States",
		"Artifact Roots",
		"missing",
		missingRepository,
		missingArtifactRoot,
		"one SQLite page",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("expected storage report output to contain %q, got %q", want, stdout.String())
		}
	}
	databaseAfter, err := os.ReadFile(store.DatabasePath(homeDir))
	if err != nil {
		t.Fatalf("read Run Database after command: %v", err)
	}
	if !bytes.Equal(databaseAfter, databaseBefore) {
		t.Fatal("storage report command changed Run Database bytes")
	}
}

func TestRunStorageReportRejectsFlags(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"storage", "report", "--format", "json"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected unsupported storage report flag to exit %d, got %d", exitPreflight, code)
	}
	if !strings.Contains(stderr.String(), "does not accept flags or arguments") {
		t.Fatalf("expected deterministic unsupported-argument diagnostic, got %q", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout for invalid storage report input, got %q", stdout.String())
	}
}

func TestRunStorageReportWithoutDatabaseCreatesNothing(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	setCommandEnvironmentForTest(t, homeDir, t.TempDir())
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"storage", "report"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected empty storage report exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no empty-report diagnostics, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "no Run Database exists") {
		t.Fatalf("expected absent Run Database explanation, got %q", stdout.String())
	}
	if _, err := os.Stat(store.DatabasePath(homeDir)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("storage report created a Run Database, stat error=%v", err)
	}
}

func withGCNow(t *testing.T, now time.Time) {
	t.Helper()
	dependencies := defaultGCDependencies()
	dependencies.now = func() time.Time { return now }
	updateCommandDependenciesForTest(t, func(commandDependencies *commandDependencies) {
		commandDependencies.gc = dependencies
	})
}

func gcSanitationRepository(t *testing.T, name string) string {
	t.Helper()
	repository := filepath.Join(t.TempDir(), name)
	mustMkdir(t, repository)
	return repository
}

func resolveGCTestArtifactRoot(t *testing.T, repository string, homeDir string) string {
	t.Helper()
	root, err := roundconfig.ResolveArtifactDirectory("", repository, homeDir)
	if err != nil {
		t.Fatalf("resolve default Artifact Root for %q: %v", repository, err)
	}
	return root
}

func createGCSanitationRun(t *testing.T, ctx context.Context, runStore *store.Store, repository string, artifactRoot string, name string, terminal bool) store.Run {
	t.Helper()
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/" + name,
		HeadBranch:     "feature/" + name,
		BaseRepository: "owner/" + name,
		PRNumber:       name,
		GitRoot:        repository,
		LocalBranch:    "feature/" + name,
		HeadSHA:        "abc123",
		ArtifactDir:    artifactRoot,
	})
	if err != nil {
		t.Fatalf("create %s sanitation Run: %v", name, err)
	}
	if terminal {
		completed, err := runStore.CompleteRun(ctx, run.ID, store.StateClean)
		if err != nil {
			t.Fatalf("complete %s sanitation Run: %v", name, err)
		}
		return completed.Run
	}
	return run
}

func gcReportCount(t *testing.T, output string, label string) int {
	t.Helper()
	for _, line := range strings.Split(output, "\n") {
		fields := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(fields) != 2 || fields[0] != label {
			continue
		}
		count, err := strconv.Atoi(strings.TrimSpace(fields[1]))
		if err != nil {
			t.Fatalf("parse %s from %q: %v", label, line, err)
		}
		return count
	}
	t.Fatalf("missing %s in output %q", label, output)
	return 0
}

type gcFixture struct {
	oldRun             store.Run
	activeRun          store.Run
	recentRun          store.Run
	oldArtifactDir     string
	activeArtifactDir  string
	recentArtifactDir  string
	orphanID           string
	orphanDir          string
	reviewArtifactPath string
}

func seedGCFixture(t *testing.T, ctx context.Context, homeDir string, artifactDir string, now time.Time) gcFixture {
	t.Helper()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run store: %v", err)
	}
	oldRun := createGCTestRun(t, ctx, runStore, artifactDir, "old-terminal", 2)
	activeRun := createGCTestRun(t, ctx, runStore, artifactDir, "active", 1)
	recentRun := createGCTestRun(t, ctx, runStore, artifactDir, "recent-terminal", 1)
	if _, err := runStore.CompleteRun(ctx, oldRun.ID, store.StateClean); err != nil {
		t.Fatalf("complete old Run: %v", err)
	}
	if _, err := runStore.CompleteRun(ctx, recentRun.ID, store.StateClean); err != nil {
		t.Fatalf("complete recent Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run store after seed: %v", err)
	}
	setRunTimestamps(t, homeDir, oldRun.ID, now.Add(-400*time.Hour), now.Add(-400*time.Hour))
	setRunTimestamps(t, homeDir, recentRun.ID, now.Add(-1*time.Hour), now.Add(-1*time.Hour))
	setRunTimestamps(t, homeDir, activeRun.ID, now.Add(-500*time.Hour), time.Time{})

	oldArtifactDir := writeRunArtifact(t, artifactDir, oldRun.ID, "old-artifact")
	activeArtifactDir := writeRunArtifact(t, artifactDir, activeRun.ID, "active-artifact")
	recentArtifactDir := writeRunArtifact(t, artifactDir, recentRun.ID, "recent-artifact")
	orphanID := "run_orphan_gc"
	orphanDir := writeRunArtifact(t, artifactDir, orphanID, "orph")
	reviewArtifactPath := filepath.Join(artifactDir, "reviews", "pr-123", "round-001", "issue.md")
	mustMkdir(t, filepath.Dir(reviewArtifactPath))
	mustWrite(t, reviewArtifactPath, "review artifact")

	return gcFixture{
		oldRun:             oldRun,
		activeRun:          activeRun,
		recentRun:          recentRun,
		oldArtifactDir:     oldArtifactDir,
		activeArtifactDir:  activeArtifactDir,
		recentArtifactDir:  recentArtifactDir,
		orphanID:           orphanID,
		orphanDir:          orphanDir,
		reviewArtifactPath: reviewArtifactPath,
	}
}

func createGCTestRun(t *testing.T, ctx context.Context, runStore *store.Store, artifactDir string, branch string, eventCount int) store.Run {
	t.Helper()
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/" + branch,
		BaseRepository: "owner/project",
		PRNumber:       branch,
		GitRoot:        filepath.Join("tmp", "repo"),
		LocalBranch:    "feature/" + branch,
		HeadSHA:        "abc123",
		ArtifactDir:    artifactDir,
	})
	if err != nil {
		t.Fatalf("create %s Run: %v", branch, err)
	}
	for index := 0; index < eventCount; index++ {
		_, err := runStore.AppendRunEvent(ctx, runevent.RunEvent{
			RunID:   run.ID,
			Batch:   1,
			Source:  runevent.SourceAgent,
			Kind:    runevent.KindAgentMessage,
			Summary: fmt.Sprintf("%s event %d", branch, index+1),
			Time:    time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC),
		})
		if err != nil {
			t.Fatalf("append %s Run Event %d: %v", branch, index+1, err)
		}
	}
	return run
}

func setRunTimestamps(t *testing.T, homeDir string, runID string, createdAt time.Time, completedAt time.Time) {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir))
	if err != nil {
		t.Fatalf("open Run Database for timestamp update: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close timestamp db: %v", err)
		}
	}()
	completed := ""
	if !completedAt.IsZero() {
		completed = completedAt.UTC().Format(time.RFC3339Nano)
	}
	_, err = db.Exec(`UPDATE runs SET created_at = ?, updated_at = ?, completed_at = ? WHERE id = ?`,
		createdAt.UTC().Format(time.RFC3339Nano),
		createdAt.UTC().Format(time.RFC3339Nano),
		completed,
		runID,
	)
	if err != nil {
		t.Fatalf("update Run %s timestamps: %v", runID, err)
	}
}

func writeRunArtifact(t *testing.T, artifactDir string, runID string, content string) string {
	t.Helper()
	runDir := filepath.Join(artifactDir, "runs", runID)
	path := filepath.Join(runDir, "agent", "batch-001.log")
	mustMkdir(t, filepath.Dir(path))
	mustWrite(t, path, content)
	return runDir
}

func assertRunEvents(t *testing.T, homeDir string, runID string, want int) {
	t.Helper()
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database reader: %v", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			t.Fatalf("close Run Database reader: %v", err)
		}
	}()
	events, err := reader.RunEventsAfter(context.Background(), runID, 0, 100)
	if err != nil {
		t.Fatalf("read Run Events for %s: %v", runID, err)
	}
	if len(events) != want {
		t.Fatalf("expected %d Run Events for %s, got %d", want, runID, len(events))
	}
}

func assertRunExists(t *testing.T, homeDir string, runID string) {
	t.Helper()
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database reader: %v", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			t.Fatalf("close Run Database reader: %v", err)
		}
	}()
	if _, ok, err := reader.Run(context.Background(), runID); err != nil || !ok {
		t.Fatalf("expected Run %s to survive, ok=%v err=%v", runID, ok, err)
	}
}
