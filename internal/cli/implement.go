package cli

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
	runworktree "roundfix/internal/worktree"
)

const implementUsage = `Usage:
  roundfix implement --spec <slug> --agent <agent>

Executes the Spec's Task Graph on the current branch as one Run: Tasks run
in dependency order, each Task's Verification commands gate one commit, and
the Run pushes only when config enables implement.auto_push and the outcome
is Clean. Without --spec, Interactive Input lists the repository's active
Specs for selection.

Options:
  --spec               Spec slug under docs/specs/
  --qa                 End the Run with the qa-gate step once every Task is
                       completed; only a pass verdict lets the Run end Clean
  --agent              Agent runtime. Supported: codex, claude, opencode
  --model              Agent model override
  --agent-command      Agent command override
  --agent-full-access  Opt into Agent runtime full-access mode
  --no-agent-console   Hide Agent-source console events from non-TTY stderr
  --detach             Start a Detached Run and print attach/stop commands
  --interactive        Open Interactive Input before starting
  --no-input           Fail instead of opening Interactive Input
`

var createRunWorktree = runworktree.Create
var integrateRunWorktree = runworktree.Integrate
var cleanupCleanRunWorktree = runworktree.CleanupClean
var pruneTerminalRunWorktrees = runworktree.PruneTerminalReport

// runImplementCommand executes the Implement Command: Preflight Validation,
// Run creation, the Live Run View, one Task cycle over the Task Graph, and
// the terminal outcome, following the runOperationalCommand shape.
func runImplementCommand(ctx context.Context, args []string, stdout, stderr io.Writer, detachChild *detachChild) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("implement"))
		return exitOK
	}

	loadedConfig, err := roundconfig.Load(roundconfig.LoadOptions{Stderr: stderr})
	if err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	maybeReportVersionFreshness(ctx, loadedConfig, stderr)
	req, err := parseImplementCommand(args, loadedConfig.Config)
	if err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	req.detachChild = detachChild
	req = applyDetachSemantics(req)
	req, err = maybeCollectInteractiveInput(ctx, req, loadedConfig, stderr)
	if err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	if err := validateImplementRequest(req); err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	if err := validateAgentConsoleDisplay(req, stderr); err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	if req.detach {
		return runDetachedCommand(append([]string{"implement"}, args...), req, loadedConfig, stdout, stderr)
	}
	outcomeNotifier := outcomeNotifierFromConfig(loadedConfig.Config)

	// Preflight Validation: every failure below exits 2 with one actionable
	// message, and nothing is written to the Run Database until every check
	// passed (ADR 0012, ADR 0013).
	gitState, err := preflight.InspectGit(ctx, loadedConfig.GitRoot, nil)
	if err != nil {
		printPreflightFailure("implement", fmt.Errorf("implement requires a git repository working tree: %w", err), stderr)
		return exitPreflight
	}
	artifactDir, err := roundconfig.ValidateArtifactDirectory(req.artifactDir, gitState.Root, loadedConfig.HomeDir)
	if err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	req.artifactDir = artifactDir
	resolvedSpecsRoot, err := roundconfig.ResolveSpecsRoot(loadedConfig, gitState.Root)
	if err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	checkoutSpecsRoot := specsRootForWorkDir(resolvedSpecsRoot, gitState.Root, gitState.Root)
	graph, err := spec.Load(checkoutSpecsRoot, req.spec)
	if err != nil {
		// spec.Load returns typed validation errors whose messages name the
		// offending Task or check and the next useful action.
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	if countNonCompletedTasks(graph.Tasks) == 0 && !req.qa {
		// With --qa, an all-completed graph is not a no-op: the Run
		// consists of the qa-gate step only (ADR 0015).
		counts := printImplementTaskLines(stdout, checkoutSpecsRoot, graph, true)
		fmt.Fprintf(stdout, "All %d Task(s) already completed; no Run was created.\n", counts.total())
		return exitOK
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
	sweepRunRetention(ctx, runStore, req.artifactDir, loadedConfig.Config.Store.JournalRetention, stderr)
	if err := pruneTerminalRunWorktreeDebris(ctx, gitState.Root, loadedConfig.Config.Worktree.Location, runtime, runStore, stderr); err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	if blocking, found, err := runStore.ActiveRunInGitRoot(ctx, gitState.Root); err != nil {
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	} else if found {
		printPreflightFailure("implement", validationError{message: fmt.Sprintf(
			"an Active Run already holds working tree %s: run_id=%s kind=%s state=%s; stop it with: roundfix stop %s",
			gitState.Root, blocking.ID, blocking.Kind, blocking.State, blocking.ID)}, stderr)
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
		Agent:       runtime.ID,
	})
	if err != nil {
		// A lost work-target race surfaces the store's ActiveRunError as-is;
		// it already names the blocking run id and the stop command.
		printPreflightFailure("implement", err, stderr)
		return exitPreflight
	}
	if err := req.reportDetachedRunCreated(run.ID); err != nil {
		markRunFailedAndNotify(ctx, runStore, run.ID, outcomeNotifier, stderr)
		printImplementRunFailure(err, stderr)
		return exitRunFailed
	}
	if len(gitState.Dirty) > 0 {
		fmt.Fprintf(stderr, "%s: note: working tree %s has %d uncommitted change(s); implement will run in a Run Worktree, and overlapping local changes end the Run Integration Pending.\n", app.Name, gitState.Root, len(gitState.Dirty))
	}
	reportNonDefaultSpecsRoot(stderr, gitState.Root, resolvedSpecsRoot)
	runRef, err := createRunWorktree(ctx, runworktree.CreateOptions{
		UserRoot:        gitState.Root,
		Location:        loadedConfig.Config.Worktree.Location,
		RunID:           run.ID,
		HeadSHA:         gitState.HEAD,
		CopyList:        loadedConfig.Config.Worktree.Copy,
		Bootstrap:       worktreeBootstrapSpec(loadedConfig.Config),
		BootstrapOutput: newBootstrapOutputWriter(ctx, run.ID, runStore, stderr),
	})
	if err != nil {
		markRunFailedAndNotify(ctx, runStore, run.ID, outcomeNotifier, stderr)
		printImplementRunFailure(err, stderr)
		return exitRunFailed
	}
	run, err = runStore.SetRunWorkDir(ctx, run.ID, runRef.Path)
	if err != nil {
		markRunFailedAndNotify(ctx, runStore, run.ID, outcomeNotifier, stderr)
		printImplementRunFailureWithWorktree(err, runRef.Path, stderr)
		return exitRunFailed
	}
	executionSpecsRoot := specsRootForWorkDir(resolvedSpecsRoot, gitState.Root, runRef.Path)
	executionGraph, err := spec.Load(executionSpecsRoot, graph.Spec.Slug)
	if err != nil {
		markRunFailedAndNotify(ctx, runStore, run.ID, outcomeNotifier, stderr)
		printImplementRunFailureWithWorktree(fmt.Errorf("load Spec from Run Worktree: %w", err), runRef.Path, stderr)
		return exitRunFailed
	}
	fmt.Fprintf(stderr, "Run Worktree: %s\n", runRef.Path)
	session := agent.SessionRefForRun(run.ID, runRef.Path)
	if err := rememberInteractiveDefaults(ctx, runStore, req); err != nil {
		closeAgentSession(ctx, collaborators.runner, runtime, session, run.ID, runStore)
		markRunFailedAndNotify(ctx, runStore, run.ID, outcomeNotifier, stderr)
		printImplementRunFailureWithWorktree(err, runRef.Path, stderr)
		return exitRunFailed
	}

	view := implementLiveRunView(req, loadedConfig, gitState, run.ID, runRef.Path, executionSpecsRoot, executionGraph)
	if !liveTUIEnabled(stderr) {
		fmt.Fprint(stderr, roundtui.RenderLiveRunView(view))
	}
	ui, err := startRunUI(ctx, view, run.ID, loadedConfig.HomeDir, runStore, stderr, req.noAgentConsole)
	if err != nil {
		closeAgentSession(ctx, collaborators.runner, runtime, session, run.ID, runStore)
		markRunFailedAndNotify(ctx, runStore, run.ID, outcomeNotifier, stderr)
		printImplementRunFailure(err, stderr)
		return exitRunFailed
	}
	defer ui.Close()

	cycleResult, err := executeImplementCycle(ctx, gitState, runRef, session, executionSpecsRoot, executionGraph, req.artifactDir, loadedConfig.Config.Logs.Agent, req.qa, loadedConfig.Config.Worktree.Concurrency, loadedConfig.Config.Worktree.Copy, worktreeBootstrapSpec(loadedConfig.Config), newBootstrapOutputWriter(ctx, run.ID, runStore, ui.progress), runtime, collaborators, runStore, ui)
	if err != nil {
		if isStopRequest(ctx, err) {
			closeAgentSession(ctx, collaborators.runner, runtime, session, run.ID, runStore)
			code := completeStoppedRunRecord(runStore, run.ID, outcomeNotifier, stderr)
			ui.Wait()
			ui.Close()
			if code != exitOK {
				printRunFailure("implement", errors.New("complete stopped Run"), stderr)
				return code
			}
			fmt.Fprintf(stderr, "Implement Run %s reached %s.\n", run.ID, store.StateStopped)
			printKeptRunWorktree(stderr, runRef.Path)
			report, counts := renderImplementTaskLines(executionSpecsRoot, executionGraph, false)
			fmt.Fprint(stdout, report)
			printImplementOutcomeLine(stdout, store.StateStopped, counts)
			return exitOK
		}
		closeAgentSession(ctx, collaborators.runner, runtime, session, run.ID, runStore)
		markRunFailedAndNotify(ctx, runStore, run.ID, outcomeNotifier, stderr)
		ui.Wait()
		ui.Close()
		printImplementRunFailureWithWorktree(err, runRef.Path, stderr)
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
	report, counts := renderImplementTaskLines(executionSpecsRoot, executionGraph, true)
	integrationCommand := ""
	if outcome == store.StateClean {
		integration, err := integrateCleanImplementRun(ctx, runRef, gitState.Branch)
		if err != nil {
			closeAgentSession(ctx, collaborators.runner, runtime, session, run.ID, runStore)
			markRunFailedAndNotify(ctx, runStore, run.ID, outcomeNotifier, stderr)
			ui.Wait()
			ui.Close()
			printImplementRunFailureWithWorktree(err, runRef.Path, stderr)
			return exitRunFailed
		}
		if integration.Mode == runworktree.ModePending {
			outcome = store.StateIntegrationPending
			integrationCommand = implementIntegrationCommand(runRef)
		} else if err := cleanupCleanRunWorktree(ctx, runRef); err != nil {
			warnCleanRunWorktreeCleanupFailed(ctx, runStore, run.ID, runRef.Path, err, stderr)
		}
	}
	pushResult := implementPushResult{}
	if outcome == store.StateClean {
		pushResult, err = maybeRunImplementAutoPush(ctx, gitState, loadedConfig.Config, collaborators, runStore, ui, run.ID, stderr)
		if err != nil {
			closeAgentSession(ctx, collaborators.runner, runtime, session, run.ID, runStore)
			markRunFailedAndNotify(ctx, runStore, run.ID, outcomeNotifier, stderr)
			ui.Wait()
			ui.Close()
			printImplementRunPushFailure(err, stderr)
			return exitRunFailed
		}
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, outcome)
	if err != nil {
		ui.Close()
		closeAgentSession(ctx, collaborators.runner, runtime, session, run.ID, runStore)
		printImplementRunFailure(err, stderr)
		return exitRunFailed
	}
	closeAgentSession(ctx, collaborators.runner, runtime, session, completed.ID, runStore)
	publishRunOutcome(ctx, runStore, completed.ID, completed.State, cycleResult.Failed+cycleResult.Skipped, stderr)
	notifyTerminalOutcome(ctx, runStore, outcomeNotifier, stderr, completed)
	// The cockpit stays on screen, read-only, until the user closes it.
	ui.Wait()
	ui.Close()
	fmt.Fprintf(stderr, "Implement Run %s reached %s.\n", completed.ID, completed.State)
	if completed.State != store.StateClean {
		printKeptRunWorktree(stderr, runRef.Path)
	}
	if completed.State == store.StateIntegrationPending {
		fmt.Fprintf(stderr, "Integration command: %s\n", integrationCommand)
	}
	fmt.Fprint(stdout, report)
	printImplementQAVerdictLine(stdout, cycleResult)
	printImplementOutcomeLineWithCommand(stdout, completed.State, counts, integrationCommand)
	if pushResult.pushed {
		fmt.Fprintf(stdout, "pushed %s/%s\n", pushResult.remote, pushResult.branch)
	}
	if completed.State == store.StateUnresolved || completed.State == store.StateIntegrationPending {
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
		artifactDir:     config.Defaults.ArtifactDir,
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
	fs.BoolVar(&req.noAgentConsole, "no-agent-console", false, "Hide Agent-source console events from non-TTY stderr")
	fs.BoolVar(&req.detach, "detach", false, "Start a Detached Run and print attach/stop commands")
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

// validateImplementRequest rejects invalid flag combinations after
// Interactive Input had its chance to fill the gaps.
func validateImplementRequest(req commandRequest) error {
	if req.noInput && req.interactive {
		return validationError{message: "--interactive cannot be used with --no-input"}
	}
	if req.spec == "" {
		return missingRequiredFlag(req, "spec")
	}
	if req.agent == "" {
		return missingRequiredFlag(req, "agent")
	}
	return validateAgent(req.agent)
}

// implementSpecOptions lists the active Spec slugs the Interactive Input
// Spec picker offers from the resolved Spec Root. An empty list fails with
// the fix instead of opening an empty picker: there is nothing to implement.
func implementSpecOptions(specsRoot string) ([]string, error) {
	options, _, err := implementSpecOptionsDetailed(specsRoot)
	return options, err
}

func implementSpecOptionsDetailed(specsRoot string) ([]string, []spec.SkippedSpec, error) {
	active, skipped, err := spec.ListActiveDetailed(specsRoot)
	if err != nil {
		return nil, nil, err
	}
	if len(active) == 0 {
		return nil, skipped, validationError{message: "no active Specs to implement: a Spec is eligible when docs/specs/<slug>/_prd.md has frontmatter status: active; create or activate one, then re-run implement"}
	}
	options := make([]string, 0, len(active))
	for _, entry := range active {
		options = append(options, entry.Slug)
	}
	return options, skipped, nil
}

func printSkippedSpecDiagnostics(stderr io.Writer, skipped []spec.SkippedSpec) {
	for _, entry := range skipped {
		fmt.Fprintf(stderr, "skipped %s: %s\n", entry.Dir, entry.Reason)
	}
}

// executeImplementCycle wires the Run engine exactly like the resolve path
// and runs one Task cycle over the full graph.
func executeImplementCycle(ctx context.Context, gitState preflight.GitState, runRef runworktree.Ref, session agent.SessionRef, specsRoot string, graph *spec.Graph, artifactDir string, agentLogs bool, qa bool, concurrency int, copyList []string, bootstrap runworktree.BootstrapSpec, bootstrapOutput io.Writer, runtime agent.RuntimeSpec, collaborators engineCollaborators, runStore *store.Store, ui *runUI) (daemon.TaskCycleResult, error) {
	runID := runRef.RunID
	fmt.Fprintf(ui.progress, "%s: implement selected Spec %s with %d Task(s); %d to execute this Run.\n", app.Name, graph.Spec.Slug, len(graph.Tasks), countNonCompletedTasks(graph.Tasks))
	fmt.Fprintf(ui.progress, "Implement Run: %s\n", runID)
	fmt.Fprintf(ui.progress, "User checkout: %s on branch %s\n", gitState.Root, gitState.Branch)
	fmt.Fprintf(ui.progress, "Run Worktree: %s on branch %s\n", runRef.Path, runRef.Branch)
	fmt.Fprintf(ui.progress, "Agent: %s\n", runtime.DisplayName)

	engine, err := daemon.NewEngine(daemon.Dependencies{
		Runner:    collaborators.runner,
		Verifier:  collaborators.verifier,
		Committer: collaborators.committer,
		// Pusher and Source only satisfy engine construction: the Task
		// cycle never invokes them. The CLI performs the optional
		// Clean-only spec push after the cycle settles (ADR 0021).
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
		RunID:           runID,
		Session:         session,
		WorkDir:         runRef.Path,
		RunWorktree:     runRef,
		SpecsRoot:       specsRoot,
		ArtifactDir:     artifactDir,
		AgentLogs:       agentLogs,
		Spec:            graph.Spec,
		Tasks:           graph.Tasks,
		Runtime:         runtime,
		QA:              qa,
		Concurrency:     concurrency,
		CopyList:        copyList,
		Bootstrap:       bootstrap,
		BootstrapOutput: bootstrapOutput,
	})
}

type implementPushResult struct {
	pushed bool
	remote string
	branch string
}

func maybeRunImplementAutoPush(ctx context.Context, gitState preflight.GitState, config roundconfig.Config, collaborators engineCollaborators, runStore *store.Store, ui *runUI, runID string, stderr io.Writer) (implementPushResult, error) {
	if !config.Implement.AutoPush {
		return implementPushResult{}, nil
	}
	remote := strings.TrimSpace(gitState.UpstreamRemote)
	branch := strings.TrimSpace(gitState.UpstreamBranch)
	if remote == "" || branch == "" {
		summary := fmt.Sprintf("Spec Run push skipped: branch %s has no upstream; set an upstream or push manually.", gitState.Branch)
		fmt.Fprintln(stderr, summary)
		publishPushDecision(ctx, ui.sink, runID, "skipped", summary, 0)
		return implementPushResult{}, nil
	}
	engine, err := daemon.NewEngine(daemon.Dependencies{
		Runner:    collaborators.runner,
		Verifier:  collaborators.verifier,
		Committer: collaborators.committer,
		Pusher:    collaborators.pusher,
		Source:    collaborators.source,
		Runs:      runStore,
		Worktree:  collaborators.worktree,
		Sink:      ui.sink,
		Progress:  ui.progress,
	})
	if err != nil {
		return implementPushResult{}, err
	}
	if err := engine.FinalPush(ctx, daemon.FinalPushRequest{
		RunID:   runID,
		WorkDir: gitState.Root,
		Remote:  remote,
		Branch:  branch,
	}); err != nil {
		summary := fmt.Sprintf("Spec Run push failed: git push %s HEAD:%s: %v", remote, branch, err)
		publishPushDecision(ctx, ui.sink, runID, "failed", summary, 0)
		return implementPushResult{}, fmt.Errorf("push Clean spec Run: %w", err)
	}
	fmt.Fprintf(ui.progress, "Spec Run push completed: git push %s HEAD:%s\n", remote, branch)
	return implementPushResult{pushed: true, remote: remote, branch: branch}, nil
}

// implementLiveRunView builds the Live Run View for a spec Run: Tasks are
// the Work Items of the left pane, in Task Graph order, located through the
// git root and Spec slug so the cockpit refreshes their statuses from the
// task files.
func implementLiveRunView(req commandRequest, loaded roundconfig.Loaded, gitState preflight.GitState, runID string, workDir string, specsRoot string, graph *spec.Graph) roundtui.LiveRunView {
	return roundtui.LiveRunView{
		Command:       "implement",
		RunKind:       store.KindImplement,
		SpecSlug:      graph.Spec.Slug,
		GitRoot:       gitState.Root,
		WorkDir:       workDir,
		SpecsRoot:     specsRoot,
		Tasks:         graph.Tasks,
		HeadBranch:    gitState.Branch,
		Agent:         displayAgent(req.agent),
		Model:         req.model,
		HEAD:          gitState.HEAD,
		RunID:         runID,
		PipelineState: "ResolvingWithAgent",
		Concurrency:   loaded.Config.Worktree.Concurrency,
		BudgetState:   formatBudgetState(loaded.Config),
		GitState:      formatGitState(gitState),
		AutoCommit:    true,
		AutoPush:      loaded.Config.Implement.AutoPush,
		LastPush:      implementPushState(loaded.Config.Implement.AutoPush),
		Console:       []string{"Agent and verification output will stream below."},
		Width:         liveViewWidth(),
	}
}

func implementPushState(enabled bool) string {
	if enabled {
		return "pending"
	}
	return "disabled"
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
func printImplementTaskLines(stdout io.Writer, specsRoot string, graph *spec.Graph, cycleFinished bool) implementTaskCounts {
	report, counts := renderImplementTaskLines(specsRoot, graph, cycleFinished)
	fmt.Fprint(stdout, report)
	return counts
}

func renderImplementTaskLines(specsRoot string, graph *spec.Graph, cycleFinished bool) (string, implementTaskCounts) {
	counts := implementTaskCounts{}
	var report bytes.Buffer
	for _, task := range graph.Tasks {
		current := task
		if err := spec.ReloadTask(specsRoot, &current); err != nil {
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
		fmt.Fprintf(&report, "%s %s — %s\n", task.ID, status, title)
	}
	return report.String(), counts
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
	printImplementOutcomeLineWithCommand(stdout, outcome, counts, "")
}

func printImplementOutcomeLineWithCommand(stdout io.Writer, outcome string, counts implementTaskCounts, integrationCommand string) {
	if outcome == store.StateClean {
		fmt.Fprintf(stdout, "Clean: all %d Task(s) completed.\n", counts.total())
		return
	}
	if outcome == store.StateIntegrationPending && strings.TrimSpace(integrationCommand) != "" {
		fmt.Fprintf(stdout, "%s: %d completed, %d failed, %d skipped, %d pending; integrate with %s\n", outcome, counts.completed, counts.failed, counts.skipped, counts.pending, integrationCommand)
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

func printImplementRunFailureWithWorktree(err error, worktreePath string, stderr io.Writer) {
	printImplementRunFailure(err, stderr)
	printKeptRunWorktree(stderr, worktreePath)
}

func printImplementRunPushFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: implement failed after Run start: %v\n", app.Name, err)
	fmt.Fprintln(stderr, "Roundfix did not complete the spec Run push; completed Task commits and preserved working tree changes remain for inspection.")
}

func printKeptRunWorktree(stderr io.Writer, worktreePath string) {
	worktreePath = strings.TrimSpace(worktreePath)
	if worktreePath == "" {
		return
	}
	fmt.Fprintf(stderr, "Run Worktree kept: %s\n", worktreePath)
}

func integrateCleanImplementRun(ctx context.Context, ref runworktree.Ref, targetBranch string) (runworktree.IntegrationResult, error) {
	runSHA, err := gitHEAD(ctx, ref.Path)
	if err != nil {
		return runworktree.IntegrationResult{}, fmt.Errorf("read Run Worktree HEAD before integration: %w", err)
	}
	result, err := integrateRunWorktree(ctx, ref, targetBranch, runSHA)
	if err != nil {
		return runworktree.IntegrationResult{}, err
	}
	return result, nil
}

func implementIntegrationCommand(ref runworktree.Ref) string {
	return "git merge --ff-only " + ref.Branch
}

func gitHEAD(ctx context.Context, workDir string) (string, error) {
	output, err := preflight.ExecGitRunner{}.RunGit(ctx, workDir, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	head := strings.TrimSpace(output)
	if head == "" {
		return "", errors.New("git returned an empty HEAD")
	}
	return head, nil
}

func pruneTerminalRunWorktreeDebris(ctx context.Context, gitRoot string, location string, runtime agent.RuntimeSpec, runStore *store.Store, stderr io.Writer) error {
	pruned, err := pruneTerminalRunWorktrees(ctx, gitRoot, location, func(runID string) bool {
		run, found, err := runStore.Run(ctx, runID)
		return err == nil && found && store.IsTerminalState(run.State)
	})
	if err != nil {
		return err
	}
	for _, ref := range pruned {
		fmt.Fprintf(stderr, "%s: reaped terminal Worktree path=%s branch=%s\n", app.Name, ref.Path, ref.Branch)
	}
	sweepTerminalRunSessions(ctx, runtime, gitRoot, runStore, stderr)
	return nil
}

func sweepTerminalRunSessions(ctx context.Context, runtime agent.RuntimeSpec, gitRoot string, runStore *store.Store, stderr io.Writer) {
	sweepCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	sessions, err := listRoundfixAgentSessions(sweepCtx, runtime, gitRoot)
	if err != nil {
		return
	}
	sort.SliceStable(sessions, func(i, j int) bool {
		return sessions[i].Name < sessions[j].Name
	})
	seen := map[string]struct{}{}
	for _, session := range sessions {
		if strings.TrimSpace(session.Name) == "" {
			continue
		}
		if _, ok := seen[session.Name]; ok {
			continue
		}
		seen[session.Name] = struct{}{}
		run, found, err := runStore.Run(sweepCtx, session.RunID)
		if err != nil || !found || !store.IsTerminalState(run.State) {
			continue
		}
		ref := sessionRefForDiscoveredRunSession(run, session)
		if err := closeStopAgentSession(sweepCtx, runtime, ref); err != nil {
			fmt.Fprintf(stderr, "%s: could not close session %s: %v\n", app.Name, ref.Name, err)
			continue
		}
		fmt.Fprintf(stderr, "%s: closed session %s\n", app.Name, ref.Name)
	}
}
