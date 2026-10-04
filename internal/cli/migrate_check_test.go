// Suite: read-only migration check public contract
// Invariant: checking schema compatibility preserves the Run Database and creates no writer artifacts.
// Boundary IN: CLI dispatch and the real SQLite Run Database in disposable homes.
// Boundary OUT: schema migration implementation, owned by internal/store.
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

	"roundfix/internal/store"
)

func TestMigrateCheckReportsACurrentDatabase(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	setCommandEnvironmentForTest(t, homeDir, t.TempDir())
	version := seedMigrateCheckCurrentDatabase(t, homeDir)
	want := fmt.Sprintf("Run Database is at schema version %d, the version this binary supports: %s\n", version, store.DatabasePath(homeDir))
	assertMigrateCheckPreservesDatabase(t, homeDir, []string{"--check"}, exitOK, want, "")
}

func TestMigrateCheckReportsAnOlderDatabaseWithoutMigrating(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	setCommandEnvironmentForTest(t, homeDir, t.TempDir())
	supported := seedOutdatedV9RunDatabase(t, homeDir)
	want := fmt.Sprintf("roundfix: migrate check: Run Database %q has schema version 9, older than the schema version %d this binary supports; run 'roundfix migrate' to upgrade it\n", store.DatabasePath(homeDir), supported)
	assertMigrateCheckPreservesDatabase(t, homeDir, []string{"--check"}, exitPreflight, "", want)
}

func TestMigrateCheckReportsANewerDatabaseWithoutWriting(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	setCommandEnvironmentForTest(t, homeDir, t.TempDir())
	newer := seedNewerRunDatabase(t, homeDir)
	want := fmt.Sprintf("roundfix: migrate check: Run Database %q has schema version %d, newer than the schema version %d this binary supports; a newer roundfix wrote it, so use that binary or run 'roundfix upgrade'\n", store.DatabasePath(homeDir), newer, newer-1)
	assertMigrateCheckPreservesDatabase(t, homeDir, []string{"--check"}, exitPreflight, "", want)
}

func TestMigrateCheckReportsAnAbsentDatabaseAndCreatesNothing(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	setCommandEnvironmentForTest(t, homeDir, t.TempDir())
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"migrate", "--check"}, &stdout, &stderr)
	want := fmt.Sprintf("No Run Database at %s; nothing to migrate\n", store.DatabasePath(homeDir))
	if code != exitOK || stdout.String() != want || stderr.Len() != 0 {
		t.Fatalf("check exit=%d stdout=%q stderr=%q, want exit=0 stdout=%q stderr=empty", code, stdout.String(), stderr.String(), want)
	}
	if _, err := os.Stat(filepath.Dir(store.DatabasePath(homeDir))); !os.IsNotExist(err) {
		t.Fatalf("check created Roundfix Home: %v", err)
	}
	entries, err := os.ReadDir(homeDir)
	if err != nil || len(entries) != 0 {
		t.Fatalf("absent check changed home: entries=%v err=%v", entries, err)
	}
}

func TestMigrateCheckRefusesAnExtraArgument(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{{"--check", "extra"}, {"extra"}, {"--unknown"}, {"--check", "--check"}, {"--check", "--help"}, {"--help", "--check"}} {
		t.Run(fmt.Sprint(args), func(t *testing.T) {
			homeDir := t.TempDir()
			setCommandEnvironmentForTest(t, homeDir, t.TempDir())
			seedMigrateCheckCurrentDatabase(t, homeDir)
			refused := args[0]
			if refused == "--check" {
				refused = args[1]
			}
			want := fmt.Sprintf("roundfix: migrate failed: unexpected argument %q\nRun 'roundfix migrate --help' for usage.\n", refused)
			assertMigrateCheckPreservesDatabase(t, homeDir, args, exitPreflight, "", want)
		})
	}
}

func TestMigrateCheckReportsAReadFailureWithoutWriting(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	setCommandEnvironmentForTest(t, homeDir, t.TempDir())
	path := store.DatabasePath(homeDir)
	markerPath := filepath.Dir(path)
	marker := []byte("not a directory")
	if err := os.WriteFile(markerPath, marker, 0600); err != nil {
		t.Fatal(err)
	}
	_, statErr := os.Stat(path)
	if statErr == nil {
		t.Fatal("expected invalid database path")
	}
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"migrate", "--check"}, &stdout, &stderr)
	want := fmt.Sprintf("roundfix: migrate check failed: open Run Database reader %q: %v\n", path, statErr)
	if code != exitRunFailed || stdout.Len() != 0 || stderr.String() != want {
		t.Fatalf("check exit=%d stdout=%q stderr=%q, want exit=1 stderr=%q", code, stdout.String(), stderr.String(), want)
	}
	if after := readMigrateDatabaseBytes(t, markerPath); !bytes.Equal(after, marker) {
		t.Fatal("check changed invalid Roundfix Home")
	}
	entries, err := os.ReadDir(homeDir)
	if err != nil || len(entries) != 1 || entries[0].Name() != filepath.Base(markerPath) {
		t.Fatalf("check created unexpected files: entries=%v err=%v", entries, err)
	}
}

func TestMigrateCheckHelpDescribesReadOnlyMode(t *testing.T) {
	t.Parallel()
	setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
	for _, args := range [][]string{{"migrate", "--help"}, {"--help"}} {
		var stdout, stderr bytes.Buffer
		code := runCLI(t, args, &stdout, &stderr)
		if code != exitOK || stderr.Len() != 0 || !strings.Contains(stdout.String(), "roundfix migrate [--check]") {
			t.Fatalf("help %v exit=%d stdout=%q stderr=%q", args, code, stdout.String(), stderr.String())
		}
		if args[0] == "migrate" && !strings.Contains(stdout.String(), "without writing") {
			t.Fatalf("help lacks read-only explanation: %q", stdout.String())
		}
	}
}

func seedMigrateCheckCurrentDatabase(t *testing.T, homeDir string) int {
	t.Helper()
	reader, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatal(err)
	}
	version, versionErr := reader.MigrationVersion(context.Background())
	closeErr := reader.Close()
	if versionErr != nil || closeErr != nil {
		t.Fatalf("seed current database: version=%v close=%v", versionErr, closeErr)
	}
	return version
}

func assertMigrateCheckPreservesDatabase(t *testing.T, homeDir string, args []string, wantCode int, wantStdout, wantStderr string) {
	t.Helper()
	path := store.DatabasePath(homeDir)
	before := readMigrateDatabaseBytes(t, path)
	version := readMigrateCheckVersion(t, path)
	files := migrateCheckHomeFiles(t, filepath.Dir(path))
	var stdout, stderr bytes.Buffer
	code := runCLI(t, append([]string{"migrate"}, args...), &stdout, &stderr)
	if code != wantCode || stdout.String() != wantStdout || stderr.String() != wantStderr {
		t.Errorf("check exit=%d stdout=%q stderr=%q, want exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String(), wantCode, wantStdout, wantStderr)
	}
	if after := readMigrateDatabaseBytes(t, path); !bytes.Equal(before, after) {
		t.Error("check changed database bytes")
	}
	if after := readMigrateCheckVersion(t, path); after != version {
		t.Errorf("check changed user_version from %d to %d", version, after)
	}
	for name := range migrateCheckHomeFiles(t, filepath.Dir(path)) {
		if !files[name] && name != "roundfix.db-wal" && name != "roundfix.db-shm" {
			t.Errorf("check created unexpected file %q", name)
		}
	}
	for name := range files {
		if _, err := os.Stat(filepath.Join(filepath.Dir(path), name)); err != nil {
			t.Errorf("check removed existing file %q: %v", name, err)
		}
	}
}

func readMigrateCheckVersion(t *testing.T, path string) int {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+path+"?mode=ro&immutable=1")
	if err != nil {
		t.Fatal(err)
	}
	var version int
	queryErr := db.QueryRow("PRAGMA user_version").Scan(&version)
	closeErr := db.Close()
	if queryErr != nil || closeErr != nil {
		t.Fatalf("read user_version: query=%v close=%v", queryErr, closeErr)
	}
	return version
}

func migrateCheckHomeFiles(t *testing.T, root string) map[string]bool {
	t.Helper()
	files := make(map[string]bool)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path != root {
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			files[relative] = true
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}
