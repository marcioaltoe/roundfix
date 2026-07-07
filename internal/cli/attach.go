package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/rounds"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
)

// attachReplayPageSize bounds each journal read so attach pages through
// large histories instead of loading them at once.
const attachReplayPageSize = 200

// attachTimelineLines bounds the rendered console ring during replay.
const attachTimelineLines = 300

type attachRequest struct {
	runID   string
	noInput bool
}

var errAttachPickerCanceled = errors.New("attach picker canceled")

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
	missingRunID := strings.TrimSpace(req.runID) == ""
	if missingRunID && (req.noInput || !attachInteractiveInputAvailable()) {
		printAttachFailure(validationError{message: "missing Run ID; run 'roundfix runs list' to discover Runs or pass a Run ID"}, stderr)
		return exitPreflight
	}
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{Stderr: stderr})
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

	if missingRunID {
		selectedRunID, err := pickAttachRun(ctx, reader, loaded.GitRoot, attachPickerInputReader(), stderr)
		if errors.Is(err, errAttachPickerCanceled) {
			return exitOK
		}
		if err != nil {
			printAttachFailure(err, stderr)
			return exitPreflight
		}
		req.runID = selectedRunID
	}

	run, found, err := reader.Run(ctx, req.runID)
	if err != nil {
		printAttachFailure(err, stderr)
		return exitPreflight
	}
	if !found {
		printAttachFailure(fmt.Errorf("Run %q does not exist", req.runID), stderr)
		return exitPreflight
	}

	concurrency := attachRunConcurrency(ctx, reader, run, loaded.Config.Worktree.Concurrency)
	if liveTUIEnabled(stdout) {
		return runAttachCockpit(ctx, loaded, reader, run, concurrency, stdout, stderr)
	}

	timeline := roundtui.NewRunTimeline(attachTimelineLines)
	cursor, err := replayRunEvents(ctx, reader, run.ID, 0, timeline)
	if err != nil {
		printAttachFailure(err, stderr)
		return exitPreflight
	}

	view := attachRunView(loaded, run, attachIssues(ctx, run), timeline.Lines(), concurrency)
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
func runAttachCockpit(ctx context.Context, loaded roundconfig.Loaded, reader *store.Store, run store.Run, concurrency int, stdout io.Writer, stderr io.Writer) int {
	view := attachRunView(loaded, run, attachIssues(ctx, run), nil, concurrency)
	err := roundtui.RunCockpit(ctx, stdout, roundtui.CockpitConfig{
		Mode:         roundtui.CockpitAttach,
		View:         view,
		RunID:        run.ID,
		Source:       reader,
		ColorEnabled: colorEnabled(stdout),
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
	fs.BoolVar(&req.noInput, "no-input", false, "Fail instead of opening Interactive Input")
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
	req.runID = strings.TrimSpace(req.runID)
	return req, nil
}

func pickAttachRun(ctx context.Context, reader *store.Store, gitRoot string, stdin io.Reader, stderr io.Writer) (string, error) {
	gitRoot = strings.TrimSpace(gitRoot)
	if gitRoot == "" {
		return "", validationError{message: "attach without a Run ID requires a Git repository; run 'roundfix runs list' to discover Runs"}
	}
	runs, err := reader.ListRuns(ctx, store.ListRunsQuery{GitRoot: gitRoot})
	if err != nil {
		return "", fmt.Errorf("list attach Runs: %w", err)
	}
	orderAttachPickerRuns(runs)
	if len(runs) == 0 {
		return "", validationError{message: "no Runs found for this repository; run 'roundfix runs list' to discover Runs"}
	}
	return collectAttachRunSelection(ctx, runs, stdin, stderr)
}

func orderAttachPickerRuns(runs []store.Run) {
	sort.SliceStable(runs, func(i, j int) bool {
		leftTerminal := store.IsTerminalState(runs[i].State)
		rightTerminal := store.IsTerminalState(runs[j].State)
		if leftTerminal != rightTerminal {
			return !leftTerminal
		}
		return false
	})
}

func collectAttachRunSelection(ctx context.Context, runs []store.Run, stdin io.Reader, stderr io.Writer) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if stdin == nil {
		stdin = strings.NewReader("")
	}
	if stderr == nil {
		stderr = io.Discard
	}
	fmt.Fprint(stderr, renderAttachRunPicker(runs))
	type readResult struct {
		line string
		err  error
	}
	resultCh := make(chan readResult, 1)
	go func() {
		line, err := bufio.NewReader(stdin).ReadString('\n')
		resultCh <- readResult{line: line, err: err}
	}()
	var line string
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case result := <-resultCh:
		if result.err != nil && result.err != io.EOF {
			return "", fmt.Errorf("read attach Run picker: %w", result.err)
		}
		line = result.line
	}
	choice := strings.TrimSpace(line)
	if choice == "" || isAttachPickerCancelChoice(choice) {
		return "", errAttachPickerCanceled
	}
	if index, err := strconv.Atoi(choice); err == nil {
		if index >= 1 && index <= len(runs) {
			return runs[index-1].ID, nil
		}
		return "", validationError{message: fmt.Sprintf("Run picker choice %d is out of range", index)}
	}
	for _, run := range runs {
		if choice == run.ID {
			return run.ID, nil
		}
	}
	return "", validationError{message: fmt.Sprintf("Run picker choice %q is not a listed Run", choice)}
}

func renderAttachRunPicker(runs []store.Run) string {
	var builder strings.Builder
	builder.WriteString("Roundfix Interactive Input\n")
	builder.WriteString("Command: attach\n")
	builder.WriteString("Runs:\n")
	for index, run := range runs {
		fmt.Fprintf(&builder, "  %d. %s  %s  %s  %s\n", index+1, run.ID, runListState(run), run.Kind, runListTarget(run))
	}
	builder.WriteString("Pick a Run by number or run id.\n")
	builder.WriteString("Press Enter to cancel.\n")
	builder.WriteString("Run: ")
	return builder.String()
}

func isAttachPickerCancelChoice(choice string) bool {
	switch strings.ToLower(strings.TrimSpace(choice)) {
	case "q", "quit", "cancel":
		return true
	default:
		return false
	}
}

var attachPickerInputReader = defaultAttachPickerInputReader
var attachInteractiveInputAvailable = defaultAttachInteractiveInputAvailable

func defaultAttachPickerInputReader() io.Reader {
	return os.Stdin
}

func defaultAttachInteractiveInputAvailable() bool {
	if strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}
	info, err := os.Stdin.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
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

func attachRunView(loaded roundconfig.Loaded, run store.Run, issues []rounds.Issue, console []string, concurrency int) roundtui.LiveRunView {
	view := roundtui.LiveRunView{
		Command:       "attach",
		Repository:    run.HeadRepository,
		PRNumber:      run.PRNumber,
		HeadBranch:    run.HeadBranch,
		HEAD:          run.HeadSHA,
		RunID:         run.ID,
		PipelineState: run.State,
		WorkDir:       run.WorkDir,
		Issues:        issues,
		Console:       console,
		// RunKind lets the cockpit's empty states explain the Run: a Fetch
		// Run writes artifacts and starts no Agent.
		RunKind: run.Kind,
		Width:   liveViewWidth(),
	}
	if run.Kind == store.KindImplement {
		specsRoot := attachSpecsRoot(loaded, run)
		view.SpecSlug = run.SpecSlug
		view.GitRoot = run.GitRoot
		view.SpecsRoot = specsRoot
		view.Concurrency = concurrency
		view.HeadBranch = run.LocalBranch
		view.Tasks = attachTasks(specsRoot, run)
	}
	return view
}

type attachConcurrencyPayload struct {
	Concurrency int `json:"concurrency"`
}

func attachRunConcurrency(ctx context.Context, reader *store.Store, run store.Run, fallback int) int {
	if run.Kind != store.KindImplement {
		return 0
	}
	cursor := int64(0)
	for {
		page, err := reader.RunEventsAfter(ctx, run.ID, cursor, attachReplayPageSize)
		if err != nil {
			break
		}
		for _, entry := range page {
			cursor = entry.Cursor
			var payload attachConcurrencyPayload
			if len(entry.Event.Payload) > 0 && json.Unmarshal(entry.Event.Payload, &payload) == nil && payload.Concurrency > 0 {
				return payload.Concurrency
			}
		}
		if len(page) < attachReplayPageSize {
			break
		}
	}
	if fallback > 0 {
		return fallback
	}
	return 0
}

// attachTasks loads the spec Run's Task Graph for the work-item pane, in
// graph order, through the Run row's execution workspace and Spec slug.
// Attach must stay usable when the Spec moved or was archived, so load
// failures render an empty pane instead of failing the command — mirroring
// attachIssues.
func attachTasks(specsRoot string, run store.Run) []spec.Task {
	graph, err := spec.Load(specsRoot, run.SpecSlug)
	if err != nil {
		return nil
	}
	return graph.Tasks
}

func attachSpecsRoot(loaded roundconfig.Loaded, run store.Run) string {
	repoRoot := strings.TrimSpace(run.GitRoot)
	if repoRoot == "" {
		repoRoot = strings.TrimSpace(loaded.GitRoot)
	}
	if repoRoot != "" {
		if resolved, err := roundconfig.ResolveSpecsRoot(loaded, repoRoot); err == nil {
			specsRoot := specsRootForWorkDir(resolved, repoRoot, attachTaskRoot(run))
			if info, statErr := os.Stat(specsRoot); statErr == nil && info.IsDir() {
				return specsRoot
			}
			return specsRootForWorkDir(resolved, repoRoot, repoRoot)
		}
	}
	return filepath.Join(attachTaskRoot(run), "docs", "specs")
}

func attachTaskRoot(run store.Run) string {
	workDir := strings.TrimSpace(run.WorkDir)
	if workDir != "" {
		if info, err := os.Stat(workDir); err == nil && info.IsDir() {
			return workDir
		}
	}
	return run.GitRoot
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
