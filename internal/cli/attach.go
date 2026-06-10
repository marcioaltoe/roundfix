package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/rounds"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
)

// attachReplayPageSize bounds each journal read so attach pages through
// large histories instead of loading them at once.
const attachReplayPageSize = 200

// attachTimelineLines bounds the rendered console ring during replay.
const attachTimelineLines = 300

type attachRequest struct {
	runID string
}

// runAttachCommand replays a Run's event timeline from the Run Database.
// Attach is non-mutating: it opens a read-only connection and never creates
// Runs, fetches, starts Agents, commits, pushes, or resolves threads.
func runAttachCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("attach"))
		return exitOK
	}
	req, err := parseAttachCommand(args)
	if err != nil {
		printAttachFailure(err, stderr)
		return exitPreflight
	}
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{})
	if err != nil {
		printAttachFailure(err, stderr)
		return exitPreflight
	}
	reader, err := store.OpenReader(ctx, loaded.HomeDir)
	if err != nil {
		printAttachFailure(err, stderr)
		return exitPreflight
	}
	defer func() {
		_ = reader.Close()
	}()

	run, found, err := reader.Run(ctx, req.runID)
	if err != nil {
		printAttachFailure(err, stderr)
		return exitPreflight
	}
	if !found {
		printAttachFailure(fmt.Errorf("Run %q does not exist", req.runID), stderr)
		return exitPreflight
	}

	if liveTUIEnabled(stdout) {
		return runAttachCockpit(ctx, reader, run, stdout, stderr)
	}

	timeline := roundtui.NewRunTimeline(attachTimelineLines)
	cursor, err := replayRunEvents(ctx, reader, run.ID, 0, timeline)
	if err != nil {
		printAttachFailure(err, stderr)
		return exitPreflight
	}

	view := attachRunView(run, attachIssues(ctx, run), timeline.Lines())
	fmt.Fprint(stdout, roundtui.RenderLiveRunView(view))
	if store.IsTerminalState(run.State) {
		fmt.Fprintf(stdout, "Run %s reached %s; timeline replayed read-only.\n", run.ID, run.State)
		return exitOK
	}

	fmt.Fprintf(stdout, "Replayed backlog through cursor %d; Run %s is %s.\n", cursor, run.ID, run.State)
	fmt.Fprintln(stdout, "Following live events. Detach with Ctrl-C; detaching never stops the Run.")
	follower := attachFollower{
		source: reader,
		sleep:  attachSleep,
		accept: func(entry store.JournalEvent) {
			if text := timeline.Append(entry.Event); text != "" {
				fmt.Fprint(stdout, text)
			}
		},
	}
	final, _, err := follower.follow(ctx, run.ID, cursor)
	if err != nil {
		if isStopRequest(ctx, err) {
			fmt.Fprintf(stdout, "Detached; Run %s keeps going.\n", run.ID)
			return exitOK
		}
		printAttachFailure(err, stderr)
		return exitRunFailed
	}
	fmt.Fprintf(stdout, "Run %s reached %s.\n", final.ID, final.State)
	return exitOK
}

// runAttachCockpit opens the interactive cockpit in the alternate screen.
// Attach mode: q/Ctrl-C detach and never stop the Run; no stop key exists.
func runAttachCockpit(ctx context.Context, reader *store.Store, run store.Run, stdout io.Writer, stderr io.Writer) int {
	view := attachRunView(run, attachIssues(ctx, run), nil)
	err := roundtui.RunCockpit(ctx, stdout, roundtui.CockpitConfig{
		Mode:   roundtui.CockpitAttach,
		View:   view,
		RunID:  run.ID,
		Source: reader,
	})
	if err != nil && !isStopRequest(ctx, err) {
		printAttachFailure(err, stderr)
		return exitRunFailed
	}
	current, found, lookupErr := reader.Run(context.WithoutCancel(ctx), run.ID)
	if lookupErr != nil || !found {
		current = run
	}
	if store.IsTerminalState(current.State) {
		fmt.Fprintf(stdout, "Run %s reached %s.\n", current.ID, current.State)
		return exitOK
	}
	fmt.Fprintf(stdout, "Detached; Run %s keeps going.\n", current.ID)
	return exitOK
}

// replayRunEvents pages journal events after the cursor into the timeline
// and returns the last accepted cursor.
func replayRunEvents(ctx context.Context, reader *store.Store, runID string, cursor int64, timeline *roundtui.RunTimeline) (int64, error) {
	for {
		page, err := reader.RunEventsAfter(ctx, runID, cursor, attachReplayPageSize)
		if err != nil {
			return cursor, err
		}
		for _, entry := range page {
			timeline.Append(entry.Event)
			cursor = entry.Cursor
		}
		if len(page) < attachReplayPageSize {
			return cursor, nil
		}
	}
}

func parseAttachCommand(args []string) (attachRequest, error) {
	req := attachRequest{}
	fs := flag.NewFlagSet("attach", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&req.runID, "run-id", "", "Run ID to attach to")
	fs.StringVar(&req.runID, "run", "", "Run ID to attach to")
	if err := fs.Parse(args); err != nil {
		return req, validationError{message: err.Error()}
	}
	remaining := fs.Args()
	if len(remaining) > 1 {
		return req, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[1])}
	}
	if len(remaining) == 1 {
		if req.runID != "" {
			return req, validationError{message: "pass Run ID either as an argument or with --run-id, not both"}
		}
		req.runID = strings.TrimSpace(remaining[0])
	}
	if strings.TrimSpace(req.runID) == "" {
		return req, validationError{message: "a Run ID is required to attach"}
	}
	return req, nil
}

// attachIssues loads the Run's Review Issues for the left pane. Attach must
// stay usable when artifacts moved or were cleaned, so lookup failures
// render an empty pane instead of failing the command.
func attachIssues(ctx context.Context, run store.Run) []rounds.Issue {
	selection, err := rounds.SelectCompatibleIssues(ctx, rounds.SelectRequest{
		ArtifactDir:    run.ArtifactDir,
		PRNumber:       run.PRNumber,
		HeadRepository: run.HeadRepository,
		HeadBranch:     run.HeadBranch,
	})
	if err != nil {
		return nil
	}
	return selection.Issues
}

func attachRunView(run store.Run, issues []rounds.Issue, console []string) roundtui.LiveRunView {
	return roundtui.LiveRunView{
		Command:       "attach",
		Repository:    run.HeadRepository,
		PRNumber:      run.PRNumber,
		HeadBranch:    run.HeadBranch,
		HEAD:          run.HeadSHA,
		RunID:         run.ID,
		PipelineState: run.State,
		Issues:        issues,
		Console:       console,
		Width:         liveViewWidth(),
	}
}

func printAttachFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "roundfix attach failed: %v\n", err)
}

// attachPollInterval paces follow-mode polls between change checks.
const attachPollInterval = 250 * time.Millisecond

// attachSleep is the follow-mode pacing seam; tests inject an immediate
// sleeper so polling behavior is provable without real waits.
var attachSleep = defaultAttachSleep

func defaultAttachSleep(ctx context.Context) error {
	timer := time.NewTimer(attachPollInterval)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// attachEventSource is the read-only journal surface follow mode needs.
// *store.Store satisfies it.
type attachEventSource interface {
	Run(ctx context.Context, runID string) (store.Run, bool, error)
	RunEventsAfter(ctx context.Context, runID string, cursor int64, limit int) ([]store.JournalEvent, error)
	DataVersion(ctx context.Context) (int64, error)
}

// attachFollower follows newly appended Run Events for an Active Run.
// Polls use the journal's data-version change signal so idle polls read no
// rows, and every read is a short autocommit query on the read-only
// connection. The cursor advances only after an event is accepted, so
// reconnects never duplicate output.
type attachFollower struct {
	source attachEventSource
	sleep  func(ctx context.Context) error
	accept func(entry store.JournalEvent)
}

// follow returns the terminal Run and the final cursor when the Run ends,
// or the context error when the user detaches. Detaching never mutates or
// stops the Run.
func (follower attachFollower) follow(ctx context.Context, runID string, cursor int64) (store.Run, int64, error) {
	lastVersion := int64(-1)
	for {
		version, err := follower.source.DataVersion(ctx)
		if err != nil {
			return store.Run{}, cursor, err
		}
		if version != lastVersion {
			lastVersion = version
			cursor, err = follower.drain(ctx, runID, cursor)
			if err != nil {
				return store.Run{}, cursor, err
			}
			run, found, err := follower.source.Run(ctx, runID)
			if err != nil {
				return store.Run{}, cursor, err
			}
			if !found {
				return store.Run{}, cursor, fmt.Errorf("Run %q disappeared while following", runID)
			}
			if store.IsTerminalState(run.State) {
				cursor, err = follower.drain(ctx, runID, cursor)
				if err != nil {
					return store.Run{}, cursor, err
				}
				return run, cursor, nil
			}
		}
		if err := follower.sleep(ctx); err != nil {
			return store.Run{}, cursor, err
		}
	}
}

func (follower attachFollower) drain(ctx context.Context, runID string, cursor int64) (int64, error) {
	for {
		page, err := follower.source.RunEventsAfter(ctx, runID, cursor, attachReplayPageSize)
		if err != nil {
			return cursor, err
		}
		for _, entry := range page {
			follower.accept(entry)
			cursor = entry.Cursor
		}
		if len(page) < attachReplayPageSize {
			return cursor, nil
		}
	}
}
