package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"roundfix/internal/app"
	"roundfix/internal/store"
)

func runMigrateCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("migrate"))
		return exitOK
	}
	if len(args) > 0 {
		printMigrateFailure(validationError{message: fmt.Sprintf("unexpected argument %q", args[0])}, stderr, true)
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

func printMigrateFailure(err error, stderr io.Writer, includeUsage bool) {
	fmt.Fprintf(stderr, "%s: migrate failed: %v\n", app.Name, err)
	if includeUsage {
		fmt.Fprintf(stderr, "Run '%s migrate --help' for usage.\n", app.Name)
	}
}
