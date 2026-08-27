package cli

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"roundfix/internal/preflight"
	"roundfix/internal/store"
)

const (
	runWindowClockLayout    = "15:04"
	runWindowAbsoluteLayout = "2006-01-02T15:04"
	runWindowDisplayLayout  = "2006-01-02 15:04 MST"

	runWindowBoundsExplanation = "The Run Window bounds when a Run may start; budget.max_run_duration bounds how long one may run."
)

const windowUsage = `Usage:
  roundfix window <set|show|clear>
  roundfix window set <HH:MM|YYYY-MM-DDTHH:MM> [--force]
  roundfix window show
  roundfix window clear

The Run Window bounds when a Run may start. budget.max_run_duration bounds
how long one may run after it starts.

Commands:
  set    Set the cutoff if no Run Window stands; use --force to replace it
  show   Print the cutoff, current time, and remaining duration
  clear  Remove the Run Window and report whether one was set

Options:
  --force  Replace an existing Run Window
`

var runWindowNow = time.Now

type windowSetRequest struct {
	cutoff string
	force  bool
}

func runWindowCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("window"))
		return exitOK
	}
	if len(args) == 0 {
		printPreflightFailure("window", validationError{message: "window requires a subcommand: set, show, or clear"}, stderr)
		return exitPreflight
	}

	subcommand := args[0]
	var setRequest windowSetRequest
	switch subcommand {
	case "set":
		var err error
		setRequest, err = parseWindowSetRequest(args[1:])
		if err != nil {
			printPreflightFailure("window", err, stderr)
			return exitPreflight
		}
	case "show", "clear":
		if err := validateWindowNoArguments(subcommand, args[1:]); err != nil {
			printPreflightFailure("window", err, stderr)
			return exitPreflight
		}
	default:
		printPreflightFailure("window", validationError{message: fmt.Sprintf("unknown window command %q; use set, show, or clear", subcommand)}, stderr)
		return exitPreflight
	}

	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		printPreflightFailure("window", err, stderr)
		return exitPreflight
	}
	gitState, err := preflight.InspectGit(ctx, loaded.GitRoot, nil)
	if err != nil {
		printPreflightFailure("window", fmt.Errorf("window requires a git repository working tree: %w", err), stderr)
		return exitPreflight
	}
	runStore, err := store.Open(ctx, loaded.HomeDir)
	if err != nil {
		printPreflightFailure("window", err, stderr)
		return exitPreflight
	}
	defer func() {
		_ = runStore.Close()
	}()

	now := commandDependenciesForContext(ctx).currentRunWindowTime()
	switch subcommand {
	case "set":
		return setRunWindow(ctx, runStore, gitState.Root, setRequest, now, stdout, stderr)
	case "show":
		return showRunWindow(ctx, runStore, gitState.Root, now, stdout, stderr)
	case "clear":
		return clearRunWindow(ctx, runStore, gitState.Root, stdout, stderr)
	default:
		panic("validated Run Window subcommand became unreachable")
	}
}

func parseWindowSetRequest(args []string) (windowSetRequest, error) {
	if len(args) == 0 {
		return windowSetRequest{}, validationError{message: "window set requires <HH:MM|YYYY-MM-DDTHH:MM>"}
	}
	if len(args) > 2 {
		return windowSetRequest{}, validationError{message: fmt.Sprintf("unexpected argument %q", args[2])}
	}

	request := windowSetRequest{}
	for _, arg := range args {
		switch {
		case arg == "--force":
			if request.force {
				return windowSetRequest{}, validationError{message: "--force may be passed only once"}
			}
			request.force = true
		case strings.HasPrefix(arg, "-"):
			return windowSetRequest{}, validationError{message: fmt.Sprintf("unknown window set option %q", arg)}
		case request.cutoff == "":
			request.cutoff = arg
		default:
			return windowSetRequest{}, validationError{message: fmt.Sprintf("unexpected argument %q", arg)}
		}
	}
	if request.cutoff == "" {
		return windowSetRequest{}, validationError{message: "window set requires <HH:MM|YYYY-MM-DDTHH:MM>"}
	}
	return request, nil
}

func validateWindowNoArguments(subcommand string, args []string) error {
	if len(args) == 0 {
		return nil
	}
	return validationError{message: fmt.Sprintf("unexpected argument %q for window %s", args[0], subcommand)}
}

func setRunWindow(
	ctx context.Context,
	runStore *store.Store,
	gitRoot string,
	request windowSetRequest,
	now time.Time,
	stdout, stderr io.Writer,
) int {
	cutoff, err := resolveRunWindowCutoff(request.cutoff, now)
	if err != nil {
		printPreflightFailure("window", err, stderr)
		return exitPreflight
	}
	if !cutoff.After(now) {
		err := validationError{message: fmt.Sprintf(
			"cutoff %s must be in the future; the current time is %s",
			formatRunWindowInstant(cutoff, now.Location()),
			formatRunWindowInstant(now, now.Location()),
		)}
		printPreflightFailure("window", err, stderr)
		return exitPreflight
	}

	window, written, err := runStore.SetRunWindow(ctx, gitRoot, cutoff, request.force)
	if err != nil {
		printPreflightFailure("window", fmt.Errorf("set Run Window: %w", err), stderr)
		return exitPreflight
	}
	if !written {
		fmt.Fprintf(stdout, "Run Window already set for %s; unchanged without --force.\n", gitRoot)
		fmt.Fprintf(stdout, "Cutoff: %s\n", formatRunWindowInstant(window.CutoffAt, now.Location()))
		return exitOK
	}
	if request.force {
		fmt.Fprintf(stdout, "Run Window replaced for %s.\n", gitRoot)
	} else {
		fmt.Fprintf(stdout, "Run Window set for %s.\n", gitRoot)
	}
	cutoffAt := formatRunWindowInstant(window.CutoffAt, now.Location())
	fmt.Fprintf(stdout, "Cutoff: %s\n", cutoffAt)
	return exitOK
}

func showRunWindow(ctx context.Context, runStore *store.Store, gitRoot string, now time.Time, stdout, stderr io.Writer) int {
	window, found, err := runStore.RunWindowFor(ctx, gitRoot)
	if err != nil {
		printPreflightFailure("window", fmt.Errorf("show Run Window: %w", err), stderr)
		return exitPreflight
	}
	if !found {
		fmt.Fprintf(stdout, "No Run Window is set for %s.\n", gitRoot)
		fmt.Fprintln(stdout, runWindowBoundsExplanation)
		return exitOK
	}

	fmt.Fprintf(stdout, "Run Window for %s\n", gitRoot)
	fmt.Fprintf(stdout, "Cutoff: %s\n", formatRunWindowInstant(window.CutoffAt, now.Location()))
	fmt.Fprintf(stdout, "Current time: %s\n", formatRunWindowInstant(now, now.Location()))
	fmt.Fprintf(stdout, "Remaining: %s\n", window.CutoffAt.Sub(now).Round(time.Second))
	fmt.Fprintln(stdout, runWindowBoundsExplanation)
	return exitOK
}

func clearRunWindow(ctx context.Context, runStore *store.Store, gitRoot string, stdout, stderr io.Writer) int {
	removed, err := runStore.ClearRunWindow(ctx, gitRoot)
	if err != nil {
		printPreflightFailure("window", fmt.Errorf("clear Run Window: %w", err), stderr)
		return exitPreflight
	}
	if removed {
		fmt.Fprintf(stdout, "Run Window cleared for %s.\n", gitRoot)
	} else {
		fmt.Fprintf(stdout, "No Run Window was set for %s.\n", gitRoot)
	}
	return exitOK
}

func resolveRunWindowCutoff(value string, now time.Time) (time.Time, error) {
	location := now.Location()
	if wallClock, err := time.ParseInLocation(runWindowClockLayout, value, location); err == nil {
		cutoff := time.Date(now.Year(), now.Month(), now.Day(), wallClock.Hour(), wallClock.Minute(), 0, 0, location)
		if !cutoff.After(now) {
			cutoff = time.Date(now.Year(), now.Month(), now.Day()+1, wallClock.Hour(), wallClock.Minute(), 0, 0, location)
		}
		return cutoff, nil
	}
	if cutoff, err := time.ParseInLocation(runWindowAbsoluteLayout, value, location); err == nil {
		return cutoff, nil
	}
	return time.Time{}, validationError{message: fmt.Sprintf("cutoff %q must match HH:MM or YYYY-MM-DDTHH:MM", value)}
}

func formatRunWindowInstant(instant time.Time, location *time.Location) string {
	return instant.In(location).Format(runWindowDisplayLayout)
}
