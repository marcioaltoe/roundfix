package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"roundfix/internal/agent"
	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
)

const implementUsage = `Usage:
  roundfix implement --spec <slug> --agent <agent>

Executes the Spec's Task Graph on the current branch as one Run: Tasks run
in dependency order, each Task's Verification commands gate one commit, and
the Run never pushes. Interactive Input for implement is not available yet,
so --spec is always required.

Options:
  --spec               Spec slug under docs/specs/
  --qa                 End the Run with the qa-gate step once every Task is
                       completed; only a pass verdict lets the Run end Clean
  --agent              Agent runtime. Supported: codex, claude, opencode
  --model              Agent model override
  --agent-command      Agent command override
  --agent-full-access  Opt into Agent runtime full-access mode
  --interactive        Open Interactive Input before starting; not available
                       for implement yet
  --no-input           Fail instead of opening Interactive Input
`

// runImplementCommand executes the Implement Command: Preflight Validation,
// Run creation, the Live Run View, one Task cycle over the Task Graph, and
// the terminal outcome, following the runOperationalCommand shape.
func runImplementCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("implement"))
		return exitOK
	}

	loadedConfig, err := roundconfig.Load(roundconfig.LoadOptions{})
	if err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	req, err := parseImplementCommand(args, loadedConfig.Config)
	if err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	if err := validateImplementRequest(req); err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}

	// Preflight Validation: every failure below exits 2 with one actionable
	// message, and nothing is written to the Run Database until every check
	// passed (ADR 0012, ADR 0013).
	gitState, err := preflight.InspectGit(ctx, loadedConfig.GitRoot, nil)
	if err != nil {
		printPreflightFailure("implement", fmt.Errorf("implement requires a git repository working tree: %w", err), stderr)
		return exitPreflight
	}
	graph, err := spec.Load(gitState.Root, req.spec)
	if err != nil {
		// spec.Load returns typed validation errors whose messages name the
		// offending Task or check and the next useful action.
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	if countNonCompletedTasks(graph.Tasks) == 0 && !req.qa {
		// With --qa, an all-completed graph is not a no-op: the Run
		// consists of the qa-gate step only (ADR 0015).
		counts := printImplementTaskLines(stdout, gitState.Root, graph, true)
		fmt.Fprintf(stdout, "All %d Task(s) already completed; no Run was created.\n", counts.total())
		return exitOK
	}
	if len(gitState.Dirty) > 0 {
		printPreflightFailure("implement", validationError{message: fmt.Sprintf(
			"working tree %s is not clean: %d uncommitted change(s); a failed Task attempt keeps its changes in the working tree on purpose — commit, stash, or discard them, then re-run implement",
			gitState.Root, len(gitState.Dirty))}, stderr)
		return exitPreflight
	}
	defaultBranch := preflight.DetectDefaultBranch(ctx, gitState.Root, gitState.Branch, nil)
	if defaultBranch.IsDefault(gitState.Branch) {
		printPreflightFailure("implement", validationError{message: fmt.Sprintf(
			"current branch %q is the repository default branch %q (detected via %s); spec Runs commit per Task, so switch to a work branch (git switch -c <branch>) and re-run implement",
			gitState.Branch, defaultBranch.Name, defaultBranch.Source)}, stderr)
		return exitPreflight
	}

	runStore, err := store.Open(ctx, loadedConfig.HomeDir)
	if err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	defer func() {
		_ = runStore.Close()
	}()
	if blocking, found, err := runStore.ActiveRunInGitRoot(ctx, gitState.Root); err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	} else if found {
		printPreflightFailure("implement", validationError{message: fmt.Sprintf(
			"an Active Run already holds working tree %s: run_id=%s kind=%s state=%s; stop it with: roundfix stop %s",
			gitState.Root, blocking.ID, blocking.Kind, blocking.State, blocking.ID)}, stderr)
		return exitPreflight
	}

	runtime, err := agent.RuntimeFor(agent.RuntimeOptions{
		Agent:            req.agent,
		CommandOverride:  req.agentCmd,
		Model:            req.model,
		EnableFullAccess: req.agentFullAccess,
	})
	if err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	collaborators := newEngineCollaborators()
	if err := collaborators.runner.Probe(ctx, runtime); err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}

	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     gitState.Root,
		LocalBranch: gitState.Branch,
		HeadSHA:     gitState.HEAD,
		SpecSlug:    graph.Spec.Slug,
	})
	if err != nil {
		// A lost work-target race surfaces the store's ActiveRunError as-is;
		// it already names the blocking run id and the stop command.
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}

	view := implementLiveRunView(req, loadedConfig, gitState, run.ID, graph)
	if !liveTUIEnabled(stderr) {
		fmt.Fprint(stderr, roundtui.RenderLiveRunView(view))
	}
	ui, err := startRunUI(ctx, view, run.ID, loadedConfig.HomeDir, runStore, stderr)
	if err != nil {
		markRunFailed(ctx, runStore, run.ID)
		printImplementRunFailure(err, stderr)
		return exitRunFailed
	}
	defer ui.Close()

	cycleResult, err := executeImplementCycle(ctx, gitState, run.ID, graph, req.qa, runtime, collaborators, runStore, ui)
	if err != nil {
		if isStopRequest(ctx, err) {
			code := completeStoppedRunRecord(runStore, run.ID)
			ui.Wait()
			ui.Close()
			if code != exitOK {
				printRunFailure("implement", errors.New("complete stopped Run"), stderr)
				return code
			}
			fmt.Fprintf(stderr, "Implement Run %s reached %s.\n", run.ID, store.StateStopped)
			counts := printImplementTaskLines(stdout, gitState.Root, graph, false)
			printImplementOutcomeLine(stdout, store.StateStopped, counts)
			return exitOK
		}
		markRunFailed(ctx, runStore, run.ID)
		ui.Wait()
		ui.Close()
		printImplementRunFailure(err, stderr)
		return exitRunFailed
	}

	outcome := store.StateClean
	if cycleResult.Failed > 0 || cycleResult.Skipped > 0 {
		outcome = store.StateUnresolved
	}
	// Only a pass verdict lets a QA Run end Clean: partial, fail, missing,
	// and unreadable all settle Unresolved (ADR 0015).
	if cycleResult.QAVerdict != "" && cycleResult.QAVerdict != spec.VerdictPass {
		outcome = store.StateUnresolved
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, outcome)
	if err != nil {
		ui.Close()
		printImplementRunFailure(err, stderr)
		return exitRunFailed
	}
	publishRunOutcome(ctx, runStore, completed.ID, completed.State, cycleResult.Failed+cycleResult.Skipped, stderr)
	// The cockpit stays on screen, read-only, until the user closes it.
	ui.Wait()
	ui.Close()
	fmt.Fprintf(stderr, "Implement Run %s reached %s.\n", completed.ID, completed.State)
	counts := printImplementTaskLines(stdout, gitState.Root, graph, true)
	printImplementQAVerdictLine(stdout, cycleResult)
	printImplementOutcomeLine(stdout, completed.State, counts)
	if completed.State == store.StateUnresolved {
		return exitRunFailed
	}
	return exitOK
}

// parseImplementCommand parses the implement flags over the config defaults,
// following the parseOperationalCommand config-then-flag precedence.
func parseImplementCommand(args []string, config roundconfig.Config) (commandRequest, error) {
	req := commandRequest{
		name:            "implement",
		agent:           config.Defaults.Agent,
		model:           config.Defaults.Model,
		agentFullAccess: config.Defaults.AgentFullAccess,
	}
	fs := flag.NewFlagSet("implement", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&req.spec, "spec", "", "Spec slug under docs/specs/")
	fs.BoolVar(&req.qa, "qa", false, "End the Run with the qa-gate step once every Task is completed")
	fs.StringVar(&req.agent, "agent", req.agent, "Agent runtime")
	fs.StringVar(&req.model, "model", req.model, "Agent model override")
	fs.StringVar(&req.agentCmd, "agent-command", "", "Agent command override")
	fs.BoolVar(&req.agentFullAccess, "agent-full-access", req.agentFullAccess, "Opt into Agent runtime full-access mode")
	fs.BoolVar(&req.interactive, "interactive", false, "Open Interactive Input before starting")
	fs.BoolVar(&req.noInput, "no-input", false, "Fail instead of opening Interactive Input")
	if err := fs.Parse(args); err != nil {
		return req, validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return req, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	req.spec = strings.TrimSpace(req.spec)
	return req, nil
}

// validateImplementRequest rejects invalid flag combinations. Interactive
// Input for implement is not wired yet, so a missing --spec is a validation
// error in every mode.
func validateImplementRequest(req commandRequest) error {
	if req.noInput && req.interactive {
		return validationError{message: "--interactive cannot be used with --no-input"}
	}
	if req.spec == "" {
		return validationError{message: "missing required --spec; Interactive Input for implement is not available yet, so pass --spec <slug>"}
	}
	if req.agent == "" {
		return validationError{message: "missing required --agent; pass --agent or set defaults.agent in the config (Interactive Input for implement is not available yet)"}
	}
	return validateAgent(req.agent)
}

// executeImplementCycle wires the Run engine exactly like the resolve path
// and runs one Task cycle over the full graph.
func executeImplementCycle(ctx context.Context, gitState preflight.GitState, runID string, graph *spec.Graph, qa bool, runtime agent.RuntimeSpec, collaborators engineCollaborators, runStore *store.Store, ui *runUI) (daemon.TaskCycleResult, error) {
	fmt.Fprintf(ui.progress, "%s: implement selected Spec %s with %d Task(s); %d to execute this Run.\n", app.Name, graph.Spec.Slug, len(graph.Tasks), countNonCompletedTasks(graph.Tasks))
	fmt.Fprintf(ui.progress, "Implement Run: %s\n", runID)
	fmt.Fprintf(ui.progress, "Working tree: %s on branch %s\n", gitState.Root, gitState.Branch)
	fmt.Fprintf(ui.progress, "Agent: %s\n", runtime.DisplayName)

	engine, err := daemon.NewEngine(daemon.Dependencies{
		Runner:    collaborators.runner,
		Verifier:  collaborators.verifier,
		Committer: collaborators.committer,
		// Pusher and Source only satisfy engine construction: the Task
		// cycle never invokes them and implement never runs Final Push,
		// because spec Runs never push (ADR 0013).
		Pusher:   collaborators.pusher,
		Source:   collaborators.source,
		Runs:     runStore,
		Worktree: collaborators.worktree,
		Sink:     ui.sink,
		Progress: ui.progress,
	})
	if err != nil {
		return daemon.TaskCycleResult{}, err
	}
	return engine.TaskCycle(ctx, daemon.TaskPlan{
		RunID:   runID,
		WorkDir: gitState.Root,
		Spec:    graph.Spec,
		Tasks:   graph.Tasks,
		Runtime: runtime,
		QA:      qa,
	})
}

// implementLiveRunView builds the cockpit header for a spec Run. Tasks
// render as the Work Items of the left pane in a later slice; until then
// the cockpit shows the Run header and the Run Event timeline.
func implementLiveRunView(req commandRequest, loaded roundconfig.Loaded, gitState preflight.GitState, runID string, graph *spec.Graph) roundtui.LiveRunView {
	return roundtui.LiveRunView{
		Command:       "implement",
		Repository:    graph.Spec.Slug,
		HeadBranch:    gitState.Branch,
		ReviewSource:  "-",
		Agent:         displayAgent(req.agent),
		Model:         req.model,
		HEAD:          gitState.HEAD,
		RunID:         runID,
		PipelineState: "ResolvingWithAgent",
		BudgetState:   formatBudgetState(loaded.Config),
		GitState:      formatGitState(gitState),
		AutoCommit:    true,
		AutoPush:      false,
		LastPush:      "disabled",
		Console:       []string{"Agent and verification output will stream below."},
		Width:         liveViewWidth(),
	}
}

type implementTaskCounts struct {
	completed int
	failed    int
	skipped   int
	pending   int
}

func (counts implementTaskCounts) total() int {
	return counts.completed + counts.failed + counts.skipped + counts.pending
}

// printImplementTaskLines prints the deterministic per-Task stdout contract:
// one line per Task in graph order — `task_NN <status> — <title>` — with the
// final status re-read from the task file. cycleFinished distinguishes
// skipped (a finished cycle left the Task pending because a needed Task did
// not complete) from pending (the cycle never reached the Task).
func printImplementTaskLines(stdout io.Writer, gitRoot string, graph *spec.Graph, cycleFinished bool) implementTaskCounts {
	counts := implementTaskCounts{}
	for _, task := range graph.Tasks {
		current := task
		if err := spec.ReloadTask(gitRoot, &current); err != nil {
			// Keep the last known state when the task file is mid-write or
			// broken; the closing report never fails the command.
			current = task
		}
		status := implementDisplayStatus(current.Status, cycleFinished)
		switch status {
		case "completed":
			counts.completed++
		case "failed":
			counts.failed++
		case "skipped":
			counts.skipped++
		default:
			counts.pending++
		}
		title := strings.TrimSpace(current.Title)
		if title == "" {
			title = task.Title
		}
		fmt.Fprintf(stdout, "%s %s — %s\n", task.ID, status, title)
	}
	return counts
}

func implementDisplayStatus(status spec.Status, cycleFinished bool) string {
	switch status {
	case spec.StatusCompleted:
		return "completed"
	case spec.StatusFailed:
		return "failed"
	default:
		if cycleFinished {
			return "skipped"
		}
		return "pending"
	}
}

// printImplementQAVerdictLine prints the deterministic QA verdict line of
// the stdout contract — `qa <verdict> — <report path>`, mirroring the
// per-Task line shape — after the Task lines and before the outcome line.
// A missing report has no path, so the detail names the absence. Nothing
// is printed when the QA step did not run.
func printImplementQAVerdictLine(stdout io.Writer, result daemon.TaskCycleResult) {
	if result.QAVerdict == "" {
		return
	}
	detail := result.QAReportPath
	if detail == "" {
		detail = "no QA Report found"
	}
	fmt.Fprintf(stdout, "qa %s — %s\n", result.QAVerdict, detail)
}

// printImplementOutcomeLine prints the one terminal outcome line of the
// stdout contract, reusing the store's terminal state vocabulary.
func printImplementOutcomeLine(stdout io.Writer, outcome string, counts implementTaskCounts) {
	if outcome == store.StateClean {
		fmt.Fprintf(stdout, "Clean: all %d Task(s) completed.\n", counts.total())
		return
	}
	fmt.Fprintf(stdout, "%s: %d completed, %d failed, %d skipped, %d pending.\n", outcome, counts.completed, counts.failed, counts.skipped, counts.pending)
}

func countNonCompletedTasks(tasks []spec.Task) int {
	count := 0
	for _, task := range tasks {
		if task.Status != spec.StatusCompleted {
			count++
		}
	}
	return count
}

func printImplementRunFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: implement failed after Run start: %v\n", app.Name, err)
	fmt.Fprintln(stderr, "Roundfix did not push; completed Task commits and preserved working tree changes remain for inspection.")
}
