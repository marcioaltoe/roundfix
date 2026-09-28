// Suite: migrate command public contract
// Invariant: roundfix migrate reports one deterministic outcome and preserves non-migratable databases.
// Boundary IN: CLI dispatch, config loading, and the real SQLite Run Database.
// Boundary OUT: migration statement selection, owned by internal/store/migrate_test.go.
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

	_ "modernc.org/sqlite"
)

func TestRunsListOnAnOlderDatabaseNamesMigrate(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	setCommandEnvironmentForTest(t, homeDir, t.TempDir())
	seedOutdatedV9RunDatabase(t, homeDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"runs", "list", "--all"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("runs list exit = %d, want %d stdout=%q stderr=%q", code, exitPreflight, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "run 'roundfix migrate'") {
		t.Fatalf("runs list stderr = %q, want migrate remedy", stderr.String())
	}
}

func TestRunsListOnANewerDatabaseNamesUpgrade(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	setCommandEnvironmentForTest(t, homeDir, t.TempDir())
	newerVersion := seedNewerRunDatabase(t, homeDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"runs", "list", "--all"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("runs list exit = %d, want %d stdout=%q stderr=%q", code, exitPreflight, stdout.String(), stderr.String())
	}
	for _, want := range []string{fmt.Sprintf("schema version %d", newerVersion), "run 'roundfix upgrade'"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("runs list stderr = %q, want %q", stderr.String(), want)
		}
	}
}

func TestMigrateCommandUpgradesAnOlderDatabaseThenRunsListWorks(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	setCommandEnvironmentForTest(t, homeDir, t.TempDir())
	supportedVersion := seedOutdatedV9RunDatabase(t, homeDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"migrate"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("migrate exit = %d, want 0 stderr=%q", code, stderr.String())
	}
	want := fmt.Sprintf("Run Database migrated from schema version 9 to %d: %s\n", supportedVersion, store.DatabasePath(homeDir))
	if stdout.String() != want || stderr.Len() != 0 {
		t.Fatalf("migrate output stdout=%q stderr=%q, want stdout=%q", stdout.String(), stderr.String(), want)
	}
	stdout.Reset()
	stderr.Reset()
	code = runCLI(t, []string{"runs", "list", "--all"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("runs list after migrate exit = %d, want 0 stderr=%q", code, stderr.String())
	}
}

func TestMigrateCommandReportsACurrentDatabase(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	setCommandEnvironmentForTest(t, homeDir, t.TempDir())
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("create Run Database: %v", err)
	}
	version, err := runStore.MigrationVersion(context.Background())
	if err != nil {
		t.Fatalf("read schema version: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	path := store.DatabasePath(homeDir)
	before := readMigrateDatabaseBytes(t, path)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"migrate"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("migrate exit = %d, want 0 stderr=%q", code, stderr.String())
	}
	want := fmt.Sprintf("Run Database is already at schema version %d: %s\n", version, path)
	if stdout.String() != want || stderr.Len() != 0 {
		t.Fatalf("migrate output stdout=%q stderr=%q, want stdout=%q", stdout.String(), stderr.String(), want)
	}
	if after := readMigrateDatabaseBytes(t, path); !bytes.Equal(after, before) {
		t.Fatal("migrate command changed a current Run Database")
	}
}

func TestMigrateCommandReportsAnAbsentDatabaseAndCreatesNothing(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	setCommandEnvironmentForTest(t, homeDir, t.TempDir())
	path := store.DatabasePath(homeDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"migrate"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("migrate exit = %d, want 0 stderr=%q", code, stderr.String())
	}
	want := fmt.Sprintf("No Run Database at %s; nothing to migrate\n", path)
	if stdout.String() != want || stderr.Len() != 0 {
		t.Fatalf("migrate output stdout=%q stderr=%q, want stdout=%q", stdout.String(), stderr.String(), want)
	}
	if _, err := os.Stat(filepath.Dir(path)); !os.IsNotExist(err) {
		t.Fatalf("migrate command created Run Database directory: %v", err)
	}
}

func TestMigrateCommandRefusesANewerDatabase(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	setCommandEnvironmentForTest(t, homeDir, t.TempDir())
	seedNewerRunDatabase(t, homeDir)
	path := store.DatabasePath(homeDir)
	before := readMigrateDatabaseBytes(t, path)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"migrate"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("migrate exit = %d, want %d stdout=%q stderr=%q", code, exitPreflight, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 || !strings.Contains(stderr.String(), "run 'roundfix upgrade'") {
		t.Fatalf("migrate output stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	if after := readMigrateDatabaseBytes(t, path); !bytes.Equal(after, before) {
		t.Fatal("migrate command changed a newer Run Database")
	}
}

func TestMigrateCommandRefusesAnArgument(t *testing.T) {
	t.Parallel()
	setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"migrate", "unexpected"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("migrate exit = %d, want %d", code, exitPreflight)
	}
	if stdout.Len() != 0 || !strings.Contains(stderr.String(), `unexpected argument "unexpected"`) {
		t.Fatalf("migrate output stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func seedNewerRunDatabase(t *testing.T, homeDir string) int {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("create Run Database: %v", err)
	}
	version, err := runStore.MigrationVersion(context.Background())
	if err != nil {
		_ = runStore.Close()
		t.Fatalf("read supported schema version: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	newerVersion := version + 1
	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir)+"?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("open newer Run Database fixture: %v", err)
	}
	if _, err := db.Exec(fmt.Sprintf(`PRAGMA user_version = %d`, newerVersion)); err != nil {
		_ = db.Close()
		t.Fatalf("advance Run Database fixture: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close newer Run Database fixture: %v", err)
	}
	return newerVersion
}

func readMigrateDatabaseBytes(t *testing.T, path string) []byte {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read Run Database %q: %v", path, err)
	}
	return content
}
