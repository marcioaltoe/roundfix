package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

type runsListOptions struct {
	all    bool
	active bool
}

func runRunsCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) || len(args) == 0 {
		fmt.Fprint(stdout, commandUsage("runs"))
		return exitOK
	}

	switch args[0] {
	case "list":
		return runRunsListCommand(ctx, args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "%s: unknown runs command %q\n", app.Name, args[0])
		fmt.Fprintf(stderr, "Run '%s runs --help' for usage.\n", app.Name)
		return exitPreflight
	}
}

func runRunsListCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("runs"))
		return exitOK
	}
	opts, err := parseRunsListOptions(args)
	if err != nil {
		printRunsListFailure(err, stderr)
		return exitPreflight
	}

	loaded, err := roundconfig.Load(roundconfig.LoadOptions{Stderr: stderr})
	if err != nil {
		printRunsListFailure(err, stderr)
		return exitPreflight
	}
	gitRoot := strings.TrimSpace(loaded.GitRoot)
	if !opts.all && gitRoot == "" {
		printRunsListFailure(validationError{message: "runs list requires a Git repository; pass --all to list every repository"}, stderr)
		return exitPreflight
	}

	reader, found, err := openRunsListReader(ctx, loaded.HomeDir)
	if err != nil {
		printRunsListFailure(err, stderr)
		return exitPreflight
	}
	if !found {
		printRunsList(stdout, nil, opts)
		return exitOK
	}
	defer func() {
		_ = reader.Close()
	}()

	query := store.ListRunsQuery{States: store.StatesAll}
	if opts.active {
		query.States = store.StatesActive
	}
	if !opts.all {
		query.GitRoot = gitRoot
	}
	runs, err := reader.ListRuns(ctx, query)
	if err != nil {
		printRunsListFailure(err, stderr)
		return exitPreflight
	}
	printRunsList(stdout, runs, opts)
	return exitOK
}

func parseRunsListOptions(args []string) (runsListOptions, error) {
	opts := runsListOptions{}
	fs := flag.NewFlagSet("runs list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.BoolVar(&opts.all, "all", false, "List Runs from every repository")
	fs.BoolVar(&opts.active, "active", false, "List only Active Runs")
	if err := fs.Parse(args); err != nil {
		return opts, validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return opts, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	return opts, nil
}

func openRunsListReader(ctx context.Context, homeDir string) (*store.Store, bool, error) {
	reader, err := store.OpenReader(ctx, homeDir)
	if err == nil {
		return reader, true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	return nil, false, err
}

func printRunsList(stdout io.Writer, runs []store.Run, opts runsListOptions) {
	if len(runs) == 0 {
		fmt.Fprintln(stdout, "No Runs found.")
		return
	}
	for _, run := range runs {
		if opts.all {
			fmt.Fprintf(stdout, "%s  %s  %s  %s  %s\n", run.ID, runListState(run), run.Kind, runListTarget(run), run.GitRoot)
			continue
		}
		fmt.Fprintf(stdout, "%s  %s  %s  %s\n", run.ID, runListState(run), run.Kind, runListTarget(run))
	}
}

func runListState(run store.Run) string {
	if store.IsTerminalState(run.State) {
		return run.State
	}
	return run.State + "*"
}

func runListTarget(run store.Run) string {
	if run.Kind == store.KindImplement {
		if strings.TrimSpace(run.SpecSlug) == "" {
			return ""
		}
		return "spec:" + run.SpecSlug
	}
	if strings.TrimSpace(run.PRNumber) == "" {
		return ""
	}
	return "pr:" + run.PRNumber
}

func printRunsListFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: runs list failed: %v\n", app.Name, err)
	fmt.Fprintf(stderr, "Run '%s runs list --help' for usage.\n", app.Name)
}
