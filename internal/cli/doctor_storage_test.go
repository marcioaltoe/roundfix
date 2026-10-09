// Suite: Doctor reclaimable Run storage notice
// Invariant: Doctor reports cheap reclaimability facts from a real Run Database without mutating it or changing readiness.
// Boundary IN: runDoctorCommand, a temporary Run Database, and temporary artifact directories
// Boundary OUT: GC reclamation and full storage measurement, owned by their dedicated suites
package cli

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

func TestDoctorStorageReportsReclaimableRunsAndNamesGC(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	artifactRoot := t.TempDir()
	now := time.Now().UTC()
	runStore, err := store.Open(t.Context(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database fixture: %v", err)
	}
	run := createGCTestRun(t, t.Context(), runStore, artifactRoot, "doctor-storage", 1)
	if _, err := runStore.CompleteRun(t.Context(), run.ID, store.StateClean); err != nil {
		t.Fatalf("complete reclaimable Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database fixture: %v", err)
	}
	setRunTimestamps(t, homeDir, run.ID, now.Add(-48*time.Hour), now.Add(-48*time.Hour))

	loaded := doctorStorageLoaded(homeDir, artifactRoot, 24*time.Hour)
	code, stdout, stderr := runDoctorWithRealStorage(t, loaded)

	if code != exitOK {
		t.Fatalf("Doctor exit code = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout, stderr)
	}
	for _, want := range []string{
		"storage: found",
		"Runs reclaimable: 1",
		"Run Database free bytes:",
		"next: roundfix gc",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("Doctor output %q does not contain %q", stdout, want)
		}
	}
	if strings.Contains(stdout, "next: roundfix gc compact") {
		t.Fatalf("Doctor output %q names compaction below the free-bytes threshold", stdout)
	}
}

func TestDoctorStorageResultNamesCompactAtTheFreeBytesThreshold(t *testing.T) {
	t.Parallel()
	result := doctorStorageResult(0, doctorStorageFreeBytesThreshold, 24*time.Hour, true)

	if result.Status != CheckStatusFound {
		t.Fatalf("storage status = %q, want %q", result.Status, CheckStatusFound)
	}
	want := "Runs reclaimable: 0; Run Database free bytes: 67108864; next: roundfix gc compact"
	if result.Detail != want {
		t.Fatalf("storage detail = %q, want %q", result.Detail, want)
	}
}

func TestDoctorStorageIsOKWhenNothingIsLeftToReclaim(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	runStore, err := store.Open(t.Context(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database fixture: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database fixture: %v", err)
	}

	loaded := doctorStorageLoaded(homeDir, t.TempDir(), 24*time.Hour)
	code, stdout, stderr := runDoctorWithRealStorage(t, loaded)

	if code != exitOK {
		t.Fatalf("Doctor exit code = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout, stderr)
	}
	if !strings.Contains(stdout, "storage: ok (nothing to reclaim; Runs reclaimable: 0; Run Database free bytes: 0)") {
		t.Fatalf("Doctor output = %q, want the empty storage result", stdout)
	}
}

func TestDoctorStorageDoesNotCreateAMissingRunDatabase(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	databasePath := store.DatabasePath(homeDir)
	loaded := doctorStorageLoaded(homeDir, t.TempDir(), 24*time.Hour)

	code, stdout, stderr := runDoctorWithRealStorage(t, loaded)

	if code != exitOK {
		t.Fatalf("Doctor exit code = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout, stderr)
	}
	if !strings.Contains(stdout, "storage: ok (no Run Database)") {
		t.Fatalf("Doctor output = %q, want missing Run Database result", stdout)
	}
	if _, err := os.Stat(databasePath); !os.IsNotExist(err) {
		t.Fatalf("Run Database exists after Doctor or returned unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(databasePath)); !os.IsNotExist(err) {
		t.Fatalf("Roundfix Home database directory exists after Doctor or returned unexpected error: %v", err)
	}
}

func TestDoctorStorageIsPartialWhenTheRunDatabaseCannotBeRead(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	databasePath := store.DatabasePath(homeDir)
	runStore, err := store.Open(t.Context(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database fixture: %v", err)
	}
	version, err := runStore.MigrationVersion(t.Context())
	if err != nil {
		t.Fatalf("read Run Database fixture schema version: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database fixture: %v", err)
	}
	database, err := sql.Open("sqlite", "file:"+databasePath)
	if err != nil {
		t.Fatalf("open Run Database fixture for schema mismatch: %v", err)
	}
	if _, err := database.Exec(fmt.Sprintf(`PRAGMA user_version = %d`, version+1)); err != nil {
		_ = database.Close()
		t.Fatalf("write newer Run Database schema version: %v", err)
	}
	if err := database.Close(); err != nil {
		t.Fatalf("close newer Run Database fixture: %v", err)
	}
	loaded := doctorStorageLoaded(homeDir, t.TempDir(), 24*time.Hour)

	code, stdout, stderr := runDoctorWithRealStorage(t, loaded)

	if code != exitOK {
		t.Fatalf("Doctor exit code = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout, stderr)
	}
	if !strings.Contains(stdout, "storage: partial (could not read Run Database:") {
		t.Fatalf("Doctor output = %q, want partial Run Database result", stdout)
	}
	if strings.Contains(stdout, "storage: failed") {
		t.Fatalf("Doctor storage result must never fail: %q", stdout)
	}
}

func TestDoctorStorageLeavesTheRunDatabaseBytesUnchanged(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	runStore, err := store.Open(t.Context(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database fixture: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database fixture: %v", err)
	}
	databasePath := store.DatabasePath(homeDir)
	before, err := os.ReadFile(databasePath)
	if err != nil {
		t.Fatalf("read Run Database before Doctor: %v", err)
	}
	loaded := doctorStorageLoaded(homeDir, t.TempDir(), 24*time.Hour)

	code, stdout, stderr := runDoctorWithRealStorage(t, loaded)
	if code != exitOK {
		t.Fatalf("Doctor exit code = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout, stderr)
	}
	after, err := os.ReadFile(databasePath)
	if err != nil {
		t.Fatalf("read Run Database after Doctor: %v", err)
	}
	if !bytes.Equal(after, before) {
		t.Fatal("Run Database bytes changed after Doctor")
	}
}

func TestDoctorStorageZeroRetentionKeepsEveryRun(t *testing.T) {
	t.Parallel()
	result := doctorStorageResult(0, 0, 0, true)

	want := "nothing to reclaim; Journal Retention 0 keeps every Run; Run Database free bytes: 0"
	if result.Status != CheckStatusOK || result.Detail != want {
		t.Fatalf("zero-retention storage result = %#v, want status ok and detail %q", result, want)
	}
}

func TestDoctorStorageDoesNotInspectArtifactDirectoriesOutsideGit(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	artifactRoot := t.TempDir()
	now := time.Now().UTC()
	runStore, err := store.Open(t.Context(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database fixture: %v", err)
	}
	run := createGCTestRun(t, t.Context(), runStore, artifactRoot, "outside-git", 0)
	if _, err := runStore.CompleteRun(t.Context(), run.ID, store.StateClean); err != nil {
		t.Fatalf("complete artifact-only Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database fixture: %v", err)
	}
	setRunTimestamps(t, homeDir, run.ID, now.Add(-48*time.Hour), now.Add(-48*time.Hour))
	mustMkdir(t, filepath.Join(artifactRoot, "runs", run.ID))
	loaded := doctorStorageLoaded(homeDir, artifactRoot, 24*time.Hour)
	loaded.GitRoot = ""

	results := defaultDoctorStorageResults(context.Background(), loaded)

	if len(results) != 1 || results[0].Status != CheckStatusOK {
		t.Fatalf("outside-Git storage results = %#v, want one ok result", results)
	}
	for _, want := range []string{"Runs reclaimable: 0", "artifact directories not inspected outside a Git repository"} {
		if !strings.Contains(results[0].Detail, want) {
			t.Fatalf("outside-Git storage detail %q does not contain %q", results[0].Detail, want)
		}
	}
}

func doctorStorageLoaded(homeDir string, artifactRoot string, retention time.Duration) roundconfig.Loaded {
	config := roundconfig.Builtin()
	config.Defaults.ArtifactDir = artifactRoot
	config.Store.JournalRetention = retention
	return roundconfig.Loaded{
		Config:  config,
		GitRoot: "/repo/project",
		HomeDir: homeDir,
	}
}

func runDoctorWithRealStorage(t *testing.T, loaded roundconfig.Loaded) (int, string, string) {
	t.Helper()
	checker := newDoctorFakeHealthChecker(
		CheckResult{Name: HealthCheckNode, Status: CheckStatusOK},
		CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK},
		CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK},
	)
	withDoctorFakeLoadedAndReadiness(
		t,
		checker,
		loaded,
		func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
			return profileProofResult{}
		},
	)
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.doctor.storage = defaultDoctorStorageResults
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runDoctorCommand(t.Context(), nil, &stdout, &stderr, commandEnvironmentForTest(t))
	return code, stdout.String(), stderr.String()
}
