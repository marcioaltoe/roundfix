package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
	runworktree "roundfix/internal/worktree"
)

const (
	runsListStateActive   = "active"
	runsListStateTerminal = "terminal"
	runsListStateAll      = "all"

	// runsListDefaultLimit bounds the default listing to the newest
	// matching Runs so the agent surface never grows without limit.
	runsListDefaultLimit = 20
)

type runsListOptions struct {
	all   bool
	state string
	limit int
}

// runsListNow is the listing clock seam; tests pin it so durations and
// running-elapsed columns are byte-stable.
var runsListNow = time.Now

// runsInteractiveInputAvailable gates the bare `runs` interactive path;
// tests override it to prove the non-interactive contract.
var runsInteractiveInputAvailable = defaultAttachInteractiveInputAvailable

func runRunsCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("runs"))
		return exitOK
	}
	if len(args) == 0 {
		if !runsInteractiveInputAvailable() || !liveTUIEnabled(stdout) {
			fmt.Fprintf(stderr, "%s: runs requires a subcommand in non-interactive mode; use 'roundfix runs list'\n", app.Name)
			return exitPreflight
		}
		return runRunsBrowserCommand(ctx, stdout, stderr)
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
		printRunsList(stdout, nil, opts, runsListNow())
		return exitOK
	}
	defer func() {
		_ = reader.Close()
	}()

	// One unbounded all-states query feeds both the visible rows and the
	// hidden counts, so the trailing note is exact without a second data
	// path.
	query := store.ListRunsQuery{States: store.StatesAll}
	if !opts.all {
		query.GitRoot = gitRoot
	}
	runs, err := reader.ListRuns(ctx, query)
	if err != nil {
		printRunsListFailure(err, stderr)
		return exitPreflight
	}
	retainedWorktrees, retainedInspectionFailures := runworktree.CountRetainedTerminalRuns(ctx, runs)
	for _, inspectErr := range retainedInspectionFailures {
		fmt.Fprintf(stderr, "%s: warning: %v\n", app.Name, inspectErr)
	}
	matching := filterRunsListState(runs, opts.state)
	visible := matching
	if opts.limit > 0 && len(matching) > opts.limit {
		visible = matching[:opts.limit]
	}
	printRunsList(stdout, visible, opts, runsListNow())
	note := runsListRetainedWorktreeNote(retainedWorktrees)
	if note == "" {
		note = runsListHiddenNote(opts.state, len(runs), len(matching), len(visible))
	}
	if note != "" {
		fmt.Fprintln(stderr, note)
	}
	return exitOK
}

func parseRunsListOptions(args []string) (runsListOptions, error) {
	opts := runsListOptions{}
	fs := flag.NewFlagSet("runs list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.BoolVar(&opts.all, "all", false, "List Runs from every repository")
	fs.StringVar(&opts.state, "state", runsListStateActive, "Filter by Run state: active, terminal, or all")
	fs.IntVar(&opts.limit, "limit", runsListDefaultLimit, "Print at most N matching Runs; 0 lists all")
	if err := fs.Parse(args); err != nil {
		return opts, validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return opts, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	switch opts.state {
	case runsListStateActive, runsListStateTerminal, runsListStateAll:
	default:
		return opts, validationError{message: fmt.Sprintf("unknown --state %q; use active, terminal, or all", opts.state)}
	}
	if opts.limit < 0 {
		return opts, validationError{message: fmt.Sprintf("--limit must be 0 or a positive count, got %d", opts.limit)}
	}
	return opts, nil
}

func filterRunsListState(runs []store.Run, state string) []store.Run {
	if state == runsListStateAll {
		return runs
	}
	wantTerminal := state == runsListStateTerminal
	matching := make([]store.Run, 0, len(runs))
	for _, run := range runs {
		if store.IsTerminalState(run.State) == wantTerminal {
			matching = append(matching, run)
		}
	}
	return matching
}

func runsListRetainedWorktreeNote(retained int) string {
	if retained <= 0 {
		return ""
	}
	worktree := "Run Worktree"
	if retained != 1 {
		worktree += "s"
	}
	return fmt.Sprintf(
		"(%d terminal %s retained; run 'roundfix reconcile' to inspect)",
		retained,
		worktree,
	)
}

// runsListHiddenNote names what the listing hid and the flag that widens
// the view. Exactly one note prints; when both the bound and the state
// filter hide Runs, the bound wins because it truncates Runs the caller
// asked for.
func runsListHiddenNote(state string, total, matching, visible int) string {
	if older := matching - visible; older > 0 {
		return fmt.Sprintf("(%d older Run(s) hidden; use --limit 0)", older)
	}
	hidden := total - matching
	if hidden <= 0 {
		return ""
	}
	switch state {
	case runsListStateActive:
		return fmt.Sprintf("(%d terminal Run(s) hidden; use --state all)", hidden)
	case runsListStateTerminal:
		return fmt.Sprintf("(%d active Run(s) hidden; use --state all)", hidden)
	default:
		return ""
	}
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

func printRunsList(stdout io.Writer, runs []store.Run, opts runsListOptions, now time.Time) {
	if len(runs) == 0 {
		fmt.Fprintln(stdout, "No Runs found.")
		return
	}
	for _, run := range runs {
		fields := roundtui.FormatRunRow(run, now, false, opts.all)
		fmt.Fprintln(stdout, strings.Join(fields, "  "))
	}
}

func printRunsListFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: runs list failed: %v\n", app.Name, err)
	fmt.Fprintf(stderr, "Run '%s runs list --help' for usage.\n", app.Name)
}
