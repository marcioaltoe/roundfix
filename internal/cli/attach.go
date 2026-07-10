package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/rounds"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
)

// journalReplayPageSize bounds each journal read so consumers page through
// large histories instead of loading them at once.
const journalReplayPageSize = 200

// attachReplayPageSize is kept for attach-specific tests and helpers.
const attachReplayPageSize = journalReplayPageSize

// attachTimelineLines bounds the rendered console ring during replay.
const attachTimelineLines = 300

type attachRequest struct {
	runID   string
	noInput bool
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
	missingRunID := strings.TrimSpace(req.runID) == ""
	if missingRunID && (req.noInput || !attachInteractiveInputAvailable() || !liveTUIEnabled(stdout)) {
		printAttachFailure(validationError{message: "missing run id in non-interactive mode; pass a run id or run 'roundfix runs list --all' to discover Runs"}, stderr)
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
		// The Run Browser is machine-wide; no repository is required.
		code, err := runRunBrowserLoop(ctx, loaded, reader, stdout, stderr)
		if err != nil {
			printAttachFailure(err, stderr)
		}
		return code
	}

	run, found, err := reader.Run(ctx, req.runID)
	if err != nil {
		printAttachFailure(err, stderr)
		return exitPreflight
	}
	if !found {
		printAttachFailure(fmt.Errorf("Run %q does not exist; picker numbers are not stable Run ids — pass a run id or run 'roundfix attach' to pick interactively", req.runID), stderr)
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
		accept: func(entry store.JournalEvent) error {
			if text := timeline.Append(entry.Event); text != "" {
				fmt.Fprint(stdout, text)
			}
			return nil
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
	return replayJournalEvents(ctx, reader, runID, cursor, func(entry store.JournalEvent) error {
		timeline.Append(entry.Event)
		return nil
	})
}

func replayJournalEvents(ctx context.Context, source journalEventSource, runID string, cursor int64, accept func(store.JournalEvent) error) (int64, error) {
	for {
		page, err := source.RunEventsAfter(ctx, runID, cursor, journalReplayPageSize)
		if err != nil {
			return cursor, err
		}
		for _, entry := range page {
			if accept != nil {
				if err := accept(entry); err != nil {
					return cursor, err
				}
			}
			cursor = entry.Cursor
		}
		if len(page) < journalReplayPageSize {
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
	fs.BoolVar(&req.noInput, "no-input", false, "Fail instead of opening the Run Browser")
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

var attachInteractiveInputAvailable = defaultAttachInteractiveInputAvailable

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
		Command:         "attach",
		Repository:      run.HeadRepository,
		PRNumber:        run.PRNumber,
		HeadBranch:      run.HeadBranch,
		Agent:           displayAgent(run.Agent),
		Model:           run.Model,
		ReasoningEffort: run.ReasoningEffort,
		HEAD:            run.HeadSHA,
		RunID:           run.ID,
		PipelineState:   run.State,
		WorkDir:         run.WorkDir,
		Issues:          issues,
		Console:         console,
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

// journalEventSource is the read-only journal surface follow mode needs.
// *store.Store satisfies it.
type journalEventSource interface {
	Run(ctx context.Context, runID string) (store.Run, bool, error)
	RunEventsAfter(ctx context.Context, runID string, cursor int64, limit int) ([]store.JournalEvent, error)
	DataVersion(ctx context.Context) (int64, error)
}

// journalFollower follows newly appended Run Events for an Active Run.
// Polls use the journal's data-version change signal so idle polls read no
// rows, and every read is a short autocommit query on the read-only
// connection. The cursor advances only after an event is accepted, so
// reconnects never duplicate output.
type journalFollower struct {
	source journalEventSource
	sleep  func(ctx context.Context) error
	accept func(entry store.JournalEvent) error
}

type attachEventSource = journalEventSource
type attachFollower = journalFollower

// follow returns the terminal Run and the final cursor when the Run ends,
// or the context error when the user detaches. Detaching never mutates or
// stops the Run.
func (follower journalFollower) follow(ctx context.Context, runID string, cursor int64) (store.Run, int64, error) {
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

func (follower journalFollower) drain(ctx context.Context, runID string, cursor int64) (int64, error) {
	for {
		page, err := follower.source.RunEventsAfter(ctx, runID, cursor, journalReplayPageSize)
		if err != nil {
			return cursor, err
		}
		for _, entry := range page {
			if follower.accept != nil {
				if err := follower.accept(entry); err != nil {
					return cursor, err
				}
			}
			cursor = entry.Cursor
		}
		if len(page) < journalReplayPageSize {
			return cursor, nil
		}
	}
}
