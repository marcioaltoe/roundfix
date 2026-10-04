package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"roundfix/internal/app"
	"roundfix/internal/store"
)

func runMigrateCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) && !(len(args) > 1 && args[0] == "--check") {
		fmt.Fprint(stdout, commandUsage("migrate"))
		return exitOK
	}
	if len(args) > 0 {
		if args[0] == "--check" && len(args) == 1 {
			return runMigrateCheck(ctx, stdout, stderr, environment)
		}
		unexpected := args[0]
		if unexpected == "--check" {
			unexpected = args[1]
		}
		printMigrateFailure(validationError{message: fmt.Sprintf("unexpected argument %q", unexpected)}, stderr, true)
		return exitPreflight
	}

	if environment.homeDirErr != nil {
		printMigrateFailure(fmt.Errorf("resolve home directory: %w", environment.homeDirErr), stderr, true)
		return exitPreflight
	}
	result, err := store.Migrate(ctx, environment.homeDir)
	if err != nil {
		var versionErr store.SchemaVersionError
		if errors.As(err, &versionErr) {
			printMigrateFailure(err, stderr, false)
			return exitPreflight
		}
		printMigrateFailure(err, stderr, false)
		return exitRunFailed
	}

	switch {
	case !result.Exists:
		fmt.Fprintf(stdout, "No Run Database at %s; nothing to migrate\n", result.Path)
	case result.From == result.To:
		fmt.Fprintf(stdout, "Run Database is already at schema version %d: %s\n", result.To, result.Path)
	default:
		fmt.Fprintf(stdout, "Run Database migrated from schema version %d to %d: %s\n", result.From, result.To, result.Path)
	}
	return exitOK
}

func runMigrateCheck(ctx context.Context, stdout, stderr io.Writer, environment commandEnvironment) int {
	if environment.homeDirErr != nil {
		fmt.Fprintf(stderr, "%s: migrate check failed: resolve home directory: %v\n", app.Name, environment.homeDirErr)
		return exitRunFailed
	}
	reader, err := store.OpenReader(ctx, environment.homeDir)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(stdout, "No Run Database at %s; nothing to migrate\n", store.DatabasePath(environment.homeDir))
		return exitOK
	}
	if err != nil {
		var versionErr store.SchemaVersionError
		if errors.As(err, &versionErr) {
			fmt.Fprintf(stderr, "%s: migrate check: %v\n", app.Name, err)
			return exitPreflight
		}
		fmt.Fprintf(stderr, "%s: migrate check failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	version, versionErr := reader.MigrationVersion(ctx)
	if err := errors.Join(versionErr, reader.Close()); err != nil {
		fmt.Fprintf(stderr, "%s: migrate check failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	fmt.Fprintf(stdout, "Run Database is at schema version %d, the version this binary supports: %s\n", version, store.DatabasePath(environment.homeDir))
	return exitOK
}

func printMigrateFailure(err error, stderr io.Writer, includeUsage bool) {
	fmt.Fprintf(stderr, "%s: migrate failed: %v\n", app.Name, err)
	if includeUsage {
		fmt.Fprintf(stderr, "Run '%s migrate --help' for usage.\n", app.Name)
	}
}
