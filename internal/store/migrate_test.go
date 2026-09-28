// Suite: Run Database schema migration
// Invariant: schema migration is explicit, serialized, and leaves absent, current, or newer databases unchanged.
// Boundary IN: public Store open and migration APIs backed by a real SQLite Run Database.
// Boundary OUT: CLI rendering and exit codes, owned by internal/cli/migrate_test.go.
package store

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestOpenReaderRefusesAnOlderDatabaseNamingMigrate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	buildV9Fixture(t, homeDir)
	path := DatabasePath(homeDir)
	before := readDatabaseBytes(t, path)

	_, err := OpenReader(ctx, homeDir)

	var versionErr SchemaVersionError
	if !errors.As(err, &versionErr) {
		t.Fatalf("OpenReader() error = %T %v, want SchemaVersionError", err, err)
	}
	want := fmt.Sprintf(
		"Run Database %q has schema version 9, older than the schema version %d this binary supports; run 'roundfix migrate' to upgrade it",
		path,
		schemaVersion,
	)
	if err.Error() != want {
		t.Fatalf("OpenReader() error = %q, want %q", err.Error(), want)
	}
	if after := readDatabaseBytes(t, path); !bytes.Equal(after, before) {
		t.Fatal("OpenReader() changed an older Run Database")
	}
}

func TestOpenReaderRefusesANewerDatabaseNamingUpgrade(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	newerVersion := createNewerDatabase(t, homeDir)
	path := DatabasePath(homeDir)
	before := readDatabaseBytes(t, path)

	_, err := OpenReader(ctx, homeDir)

	var versionErr SchemaVersionError
	if !errors.As(err, &versionErr) {
		t.Fatalf("OpenReader() error = %T %v, want SchemaVersionError", err, err)
	}
	want := fmt.Sprintf(
		"Run Database %q has schema version %d, newer than the schema version %d this binary supports; a newer roundfix wrote it, so use that binary or run 'roundfix upgrade'",
		path,
		newerVersion,
		schemaVersion,
	)
	if err.Error() != want {
		t.Fatalf("OpenReader() error = %q, want %q", err.Error(), want)
	}
	if after := readDatabaseBytes(t, path); !bytes.Equal(after, before) {
		t.Fatal("OpenReader() changed a newer Run Database")
	}
}

func TestOpenRefusesANewerDatabaseWithTheTypedError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	newerVersion := createNewerDatabase(t, homeDir)

	_, err := Open(ctx, homeDir)

	var versionErr SchemaVersionError
	if !errors.As(err, &versionErr) {
		t.Fatalf("Open() error = %T %v, want SchemaVersionError", err, err)
	}
	if versionErr != (SchemaVersionError{
		Path:      DatabasePath(homeDir),
		Found:     newerVersion,
		Supported: schemaVersion,
	}) {
		t.Fatalf("Open() SchemaVersionError = %#v", versionErr)
	}
}

func TestConcurrentOpensMigrateAnOlderDatabaseOnce(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	buildV9Fixture(t, homeDir)

	const callers = 8
	start := make(chan struct{})
	errorsByCaller := make(chan error, callers)
	var ready sync.WaitGroup
	ready.Add(callers)
	for range callers {
		go func() {
			ready.Done()
			<-start
			runStore, err := Open(ctx, homeDir)
			if err == nil {
				err = runStore.Close()
			}
			errorsByCaller <- err
		}()
	}
	ready.Wait()
	close(start)

	for range callers {
		if err := <-errorsByCaller; err != nil {
			t.Fatalf("concurrent Open() failed: %v", err)
		}
	}
	reader, err := OpenReader(ctx, homeDir)
	if err != nil {
		t.Fatalf("OpenReader() after concurrent migration: %v", err)
	}
	defer closeStore(t, reader)
	if version, err := reader.MigrationVersion(ctx); err != nil {
		t.Fatalf("MigrationVersion() after concurrent migration: %v", err)
	} else if version != schemaVersion {
		t.Fatalf("schema version = %d, want %d", version, schemaVersion)
	}
}

func TestMigrateUpgradesAnOlderDatabase(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	buildV9Fixture(t, homeDir)

	result, err := Migrate(ctx, homeDir)

	if err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}
	want := MigrationResult{Path: DatabasePath(homeDir), From: 9, To: schemaVersion, Exists: true}
	if result != want {
		t.Fatalf("Migrate() = %#v, want %#v", result, want)
	}
	reader, err := OpenReader(ctx, homeDir)
	if err != nil {
		t.Fatalf("OpenReader() after Migrate(): %v", err)
	}
	closeStore(t, reader)
}

func TestMigrateLeavesACurrentDatabaseUnwritten(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	runStore := openTestStore(t, ctx, homeDir)
	closeStore(t, runStore)
	path := DatabasePath(homeDir)
	before := readDatabaseBytes(t, path)

	result, err := Migrate(ctx, homeDir)

	if err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}
	want := MigrationResult{Path: path, From: schemaVersion, To: schemaVersion, Exists: true}
	if result != want {
		t.Fatalf("Migrate() = %#v, want %#v", result, want)
	}
	if after := readDatabaseBytes(t, path); !bytes.Equal(after, before) {
		t.Fatal("Migrate() changed a current Run Database")
	}
}

func TestMigrateCreatesNothingForAnAbsentDatabase(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	path := DatabasePath(homeDir)

	result, err := Migrate(ctx, homeDir)

	if err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}
	want := MigrationResult{Path: path}
	if result != want {
		t.Fatalf("Migrate() = %#v, want %#v", result, want)
	}
	if _, err := os.Stat(filepath.Dir(path)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Migrate() created Run Database directory: %v", err)
	}
}

func TestMigrateRefusesANewerDatabaseWithoutWriting(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	newerVersion := createNewerDatabase(t, homeDir)
	path := DatabasePath(homeDir)
	before := readDatabaseBytes(t, path)

	result, err := Migrate(ctx, homeDir)

	if result != (MigrationResult{Path: path, From: newerVersion, To: newerVersion, Exists: true}) {
		t.Fatalf("Migrate() = %#v", result)
	}
	var versionErr SchemaVersionError
	if !errors.As(err, &versionErr) {
		t.Fatalf("Migrate() error = %T %v, want SchemaVersionError", err, err)
	}
	if after := readDatabaseBytes(t, path); !bytes.Equal(after, before) {
		t.Fatal("Migrate() changed a newer Run Database")
	}
}

func createNewerDatabase(t *testing.T, homeDir string) int {
	t.Helper()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, homeDir)
	closeStore(t, runStore)
	newerVersion := schemaVersion + 1
	db, err := sql.Open("sqlite", writerDSN(DatabasePath(homeDir)))
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

func readDatabaseBytes(t *testing.T, path string) []byte {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read Run Database %q: %v", path, err)
	}
	return content
}

func TestSchemaVersionErrorNamesNoOperationalCommand(t *testing.T) {
	t.Parallel()
	for _, found := range []int{schemaVersion - 1, schemaVersion + 1} {
		err := SchemaVersionError{Path: "run.db", Found: found, Supported: schemaVersion}
		if strings.Contains(err.Error(), "resolve, watch, or implement") {
			t.Fatalf("SchemaVersionError still names operational commands: %q", err.Error())
		}
	}
}
