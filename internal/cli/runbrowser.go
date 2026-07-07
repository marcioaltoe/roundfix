package cli

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
)

// runBrowserSession opens one Run Browser pass over pre-queried listings
// and reports how it closed. It is a seam: CLI tests script outcomes here
// and exercise the loop through the model, not terminal emulation.
var runBrowserSession = defaultRunBrowserSession

func defaultRunBrowserSession(ctx context.Context, stdout io.Writer, repo string, active, all []store.Run, otherActive int) (roundtui.BrowserOutcome, error) {
	return roundtui.RunBrowserSession(ctx, stdout, roundtui.NewRunBrowser(repo, active, all).WithOtherActive(otherActive))
}

// browserAttachCockpit is the cockpit step of the browser loop. It is the
// same runAttachCockpit the explicit `attach <run-id>` path uses, seamed so
// loop tests can observe the selection without opening a terminal program.
var browserAttachCockpit = runAttachCockpit

// runRunBrowserLoop drives interactive Run discovery: browser → selected
// Run's attach cockpit → refreshed browser, until the user cancels. The
// loop is read-only; cancel exits 0 with no side effects. A non-nil error
// reports the failed operation and the caller prints it in its own
// failure shape.
func runRunBrowserLoop(ctx context.Context, loaded roundconfig.Loaded, reader *store.Store, gitRoot string, stdout, stderr io.Writer) (int, error) {
	repo := filepath.Base(gitRoot)
	for {
		active, all, otherActive, err := browserRunListings(ctx, reader, gitRoot)
		if err != nil {
			return exitPreflight, err
		}
		outcome, err := runBrowserSession(ctx, stdout, repo, active, all, otherActive)
		if err != nil {
			return exitRunFailed, fmt.Errorf("run the Run Browser: %w", err)
		}
		if outcome.Cancelled {
			return exitOK, nil
		}
		run, found, err := reader.Run(ctx, outcome.RunID)
		if err != nil {
			return exitPreflight, err
		}
		if !found {
			// The Run was pruned between the listing and the selection;
			// the refreshed browser shows current data.
			continue
		}
		concurrency := attachRunConcurrency(ctx, reader, run, loaded.Config.Worktree.Concurrency)
		if code := browserAttachCockpit(ctx, loaded, reader, run, concurrency, stdout, stderr); code != exitOK {
			return code, nil
		}
	}
}

// browserRunListings is the browser's fresh store query. One all-states
// listing feeds both filters — the same single data path runs list uses —
// plus a cross-repository Active count so an empty browser can say where
// the action is.
func browserRunListings(ctx context.Context, reader *store.Store, gitRoot string) (active []store.Run, all []store.Run, otherActive int, err error) {
	all, err = reader.ListRuns(ctx, store.ListRunsQuery{GitRoot: gitRoot, States: store.StatesAll})
	if err != nil {
		return nil, nil, 0, fmt.Errorf("list Runs for the Run Browser: %w", err)
	}
	active = make([]store.Run, 0, len(all))
	for _, run := range all {
		if !store.IsTerminalState(run.State) {
			active = append(active, run)
		}
	}
	everyActive, err := reader.ListRuns(ctx, store.ListRunsQuery{States: store.StatesActive})
	if err != nil {
		return nil, nil, 0, fmt.Errorf("count Active Runs across repositories for the Run Browser: %w", err)
	}
	for _, run := range everyActive {
		if run.GitRoot != gitRoot {
			otherActive++
		}
	}
	return active, all, otherActive, nil
}

// runRunsBrowserCommand is the bare `runs` interactive entry point: it
// opens the Run Browser over the current repository's Runs.
func runRunsBrowserCommand(ctx context.Context, stdout, stderr io.Writer) int {
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{Stderr: stderr})
	if err != nil {
		printRunsBrowserFailure(err, stderr)
		return exitPreflight
	}
	gitRoot := strings.TrimSpace(loaded.GitRoot)
	if gitRoot == "" {
		printRunsBrowserFailure(validationError{message: "the Run Browser requires a Git repository; use 'roundfix runs list --all' to list every repository"}, stderr)
		return exitPreflight
	}
	reader, found, err := openRunsListReader(ctx, loaded.HomeDir)
	if err != nil {
		printRunsBrowserFailure(err, stderr)
		return exitPreflight
	}
	if !found {
		// No Run Database yet: one browser pass shows the empty state;
		// with no rows, only cancel closes it.
		if _, err := runBrowserSession(ctx, stdout, filepath.Base(gitRoot), nil, nil, 0); err != nil {
			printRunsBrowserFailure(fmt.Errorf("run the Run Browser: %w", err), stderr)
			return exitRunFailed
		}
		return exitOK
	}
	defer func() {
		_ = reader.Close()
	}()
	code, err := runRunBrowserLoop(ctx, loaded, reader, gitRoot, stdout, stderr)
	if err != nil {
		printRunsBrowserFailure(err, stderr)
	}
	return code
}

func printRunsBrowserFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: runs failed: %v\n", app.Name, err)
	fmt.Fprintf(stderr, "Run '%s runs --help' for usage.\n", app.Name)
}
