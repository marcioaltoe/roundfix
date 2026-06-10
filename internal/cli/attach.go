package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"

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

	timeline := roundtui.NewRunTimeline(attachTimelineLines)
	if _, err := replayRunEvents(ctx, reader, run.ID, 0, timeline); err != nil {
		printAttachFailure(err, stderr)
		return exitPreflight
	}

	view := attachRunView(run, attachIssues(ctx, run), timeline.Lines())
	fmt.Fprint(stdout, roundtui.RenderLiveRunView(view))
	if store.IsTerminalState(run.State) {
		fmt.Fprintf(stdout, "Run %s reached %s; timeline replayed read-only.\n", run.ID, run.State)
	} else {
		fmt.Fprintf(stdout, "Run %s is %s; timeline replayed read-only. Detaching never stops the Run.\n", run.ID, run.State)
	}
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
