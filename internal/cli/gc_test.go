package cli

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

func withGCNow(t *testing.T, now time.Time) {
	t.Helper()
	dependencies := defaultGCDependencies()
	dependencies.now = func() time.Time { return now }
	updateCommandDependenciesForTest(t, func(commandDependencies *commandDependencies) {
		commandDependencies.gc = dependencies
	})
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
