package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	"roundfix/internal/preflight"
	"roundfix/internal/reviewsource"
	"roundfix/internal/reviewsource/coderabbit"
	"roundfix/internal/rounds"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
	"roundfix/internal/watch"
	roundskills "roundfix/skills"
)

const usage = `Roundfix

Usage:
  roundfix --help
  roundfix --version
  roundfix fetch --source coderabbit --pr <number>
  roundfix resolve --pr <number> --agent <agent>
  roundfix watch --source coderabbit --pr <number> --agent <agent> --until-clean
  roundfix skills check
  roundfix skills install --target <codex|claude|opencode|all>

Commands:
  fetch      Download review issues for an Open Pull Request
  resolve    Resolve downloaded Unresolved Review Issues
  watch      Fetch and resolve in a watched loop
  skills     Check or install Roundfix agent skills

Options:
  -h, --help      Show help
  --version       Show version
`

const (
	exitOK        = 0
	exitRunFailed = 1
	exitPreflight = 2
	exitSIGINT    = 130
)

type commandRequest struct {
	name        string
	pr          string
	source      string
	agent       string
	round       string
	noInput     bool
	interactive bool
	inputShown  bool
	untilClean  bool
	maxRounds   int
	artifactDir string
	baseRepo    string
	model       string
	agentCmd    string
	headBranch  string
	headRepo    string
}

var runCommandPreflight = defaultRunCommandPreflight
var fetchReviewItems = defaultFetchReviewItems
var runAgentRuntime agent.Runner = agent.ExecRunner{}
var runVerificationGate daemon.Verifier = daemon.ExecVerifier{}
var createBatchCommit daemon.Committer = daemon.GitCommitter{}
var resolveReviewSourceIssues = defaultResolveReviewSourceIssues
var runFinalPush daemon.Pusher = daemon.GitPusher{}
var watchReviewStatus = defaultWatchReviewStatus
var watchClock watch.Clock
var watchSleeper watch.Sleeper
var inspectChangedPaths = defaultInspectChangedPaths
var collectInteractiveInput = defaultCollectInteractiveInput
var suggestCurrentPullRequest = defaultSuggestCurrentPullRequest

type validationError struct {
	message string
}

func (err validationError) Error() string {
	return err.message
}

func Run(args []string, stdout, stderr io.Writer) int {
	ctx, cleanup, interrupted := interruptContext(context.Background())
	defer cleanup()
	code := runWithContext(ctx, args, stdout, stderr)
	return exitForInterrupt(code, interrupted())
}

func RunContext(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	return runWithContext(ctx, args, stdout, stderr)
}

func runWithContext(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(args) == 0 {
		fmt.Fprint(stdout, usage)
		return exitOK
	}

	switch args[0] {
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usage)
		return exitOK
	case "--version", "version":
		fmt.Fprintf(stdout, "%s %s\n", app.Name, app.Version)
		return exitOK
	case "skills":
		return runSkillsCommand(ctx, args[1:], stdout, stderr)
	case "fetch", "resolve", "watch":
		return runOperationalCommand(ctx, args[0], args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "%s: unknown command %q\n", app.Name, args[0])
		fmt.Fprintf(stderr, "Run '%s --help' for usage.\n", app.Name)
		return exitPreflight
	}
}

func runSkillsCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("skills"))
		return exitOK
	}
	if len(args) == 0 {
		fmt.Fprint(stdout, commandUsage("skills"))
		return exitOK
	}

	switch args[0] {
	case "check":
		if commandWantsHelp(args[1:]) {
			fmt.Fprint(stdout, commandUsage("skills"))
			return exitOK
		}
		if len(args) > 1 {
			fmt.Fprintf(stderr, "%s: unexpected argument %q\n", app.Name, args[1])
			fmt.Fprintf(stderr, "Run '%s skills --help' for usage.\n", app.Name)
			return exitPreflight
		}
		diagnostics := roundskills.Check()
		if len(diagnostics) > 0 {
			fmt.Fprintf(stderr, "%s: Roundfix skills check failed:\n", app.Name)
			for _, diagnostic := range diagnostics {
				fmt.Fprintf(stderr, "  %s: %s\n", diagnostic.Path, diagnostic.Message)
			}
			return exitRunFailed
		}
		fmt.Fprintf(stdout, "Roundfix skills check passed: %s\n", strings.Join(roundskills.Names(), ", "))
		return exitOK
	case "install":
		return runSkillsInstall(ctx, args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "%s: unknown skills command %q\n", app.Name, args[0])
		fmt.Fprintf(stderr, "Run '%s skills --help' for usage.\n", app.Name)
		return exitPreflight
	}
}

func runSkillsInstall(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("skills install", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	target := fs.String("target", "all", "Skill install target: codex, claude, opencode, or all")
	dir := fs.String("dir", "", "Override target skills directory for a single target")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "%s: %v\n", app.Name, err)
		fmt.Fprintf(stderr, "Run '%s skills --help' for usage.\n", app.Name)
		return exitPreflight
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		fmt.Fprintf(stderr, "%s: unexpected argument %q\n", app.Name, remaining[0])
		fmt.Fprintf(stderr, "Run '%s skills --help' for usage.\n", app.Name)
		return exitPreflight
	}
	if strings.TrimSpace(*dir) != "" && strings.TrimSpace(*target) == "all" {
		fmt.Fprintf(stderr, "%s: --dir requires a single --target value\n", app.Name)
		return exitPreflight
	}

	targetDirs := map[string]string{}
	if strings.TrimSpace(*dir) != "" {
		targetDirs[strings.TrimSpace(*target)] = strings.TrimSpace(*dir)
	}
	result, err := roundskills.Install(ctx, roundskills.InstallRequest{
		Target:     strings.TrimSpace(*target),
		TargetDirs: targetDirs,
	})
	if err != nil {
		fmt.Fprintf(stderr, "%s: skills install failed: %v\n", app.Name, err)
		return exitPreflight
	}
	for _, installed := range result.Targets {
		fmt.Fprintf(stdout, "Installed Roundfix skills for %s: %s (%d file(s))\n", installed.Target, installed.Dir, installed.Files)
	}
	return exitOK
}

func exitForInterrupt(code int, interrupted bool) int {
	if interrupted {
		return exitSIGINT
	}
	return code
}

func interruptContext(parent context.Context) (context.Context, func(), func() bool) {
	ctx, cancel := context.WithCancel(parent)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	done := make(chan struct{})
	var interrupted atomic.Bool

	go func() {
		defer close(done)
		select {
		case <-signals:
			interrupted.Store(true)
			cancel()
		case <-ctx.Done():
		}
	}()

	cleanup := func() {
		signal.Stop(signals)
		cancel()
		<-done
	}
	return ctx, cleanup, interrupted.Load
}

func maybeCollectInteractiveInput(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, stderr io.Writer) (commandRequest, error) {
	if req.noInput && req.interactive {
		return req, validationError{message: "--interactive cannot be used with --no-input"}
	}
	if !shouldOpenInteractiveInput(req) {
		return req, nil
	}
	inputReq, err := buildInteractiveInputRequest(ctx, req, loaded)
	if err != nil {
		return req, err
	}
	values, err := collectInteractiveInput(ctx, inputReq)
	if err != nil {
		return req, fmt.Errorf("Interactive Input failed: %w", err)
	}
	req = applyInteractiveValues(req, values)
	req.inputShown = true
	fmt.Fprintln(stderr, "Interactive Input collected command parameters.")
	return req, nil
}

func shouldOpenInteractiveInput(req commandRequest) bool {
	if req.noInput {
		return false
	}
	if req.interactive {
		return true
	}
	if strings.TrimSpace(req.pr) == "" {
		return true
	}
	switch req.name {
	case "resolve", "watch":
		return strings.TrimSpace(req.agent) == ""
	default:
		return false
	}
}

func buildInteractiveInputRequest(ctx context.Context, req commandRequest, loaded roundconfig.Loaded) (roundtui.InputRequest, error) {
	remembered, err := loadRememberedInteractiveDefaults(ctx, loaded.HomeDir)
	if err != nil {
		return roundtui.InputRequest{}, err
	}
	currentPR, _ := suggestCurrentPullRequest(ctx, loaded.GitRoot)
	prSuggestion := roundtui.Suggestion{Value: currentPR, Source: "current"}
	if prSuggestion.Value == "" {
		prSuggestion = roundtui.Suggestion{Value: remembered.PRNumber, Source: "remembered"}
	}
	agentSuggestion := roundtui.Suggestion{Value: req.agent, Source: "config"}
	if agentSuggestion.Value == "" {
		agentSuggestion = roundtui.Suggestion{Value: remembered.Agent, Source: "remembered"}
	}
	return roundtui.InputRequest{
		Command: req.name,
		Values: roundtui.CommandValues{
			PRNumber:     req.pr,
			ReviewSource: req.source,
			Agent:        req.agent,
			Round:        req.round,
			ArtifactDir:  req.artifactDir,
			Model:        req.model,
			MaxRounds:    req.maxRounds,
			UntilClean:   req.untilClean,
		},
		PRSuggestion:    prSuggestion,
		AgentSuggestion: agentSuggestion,
	}, nil
}

func applyInteractiveValues(req commandRequest, values roundtui.CommandValues) commandRequest {
	req.pr = strings.TrimSpace(values.PRNumber)
	req.source = strings.TrimSpace(values.ReviewSource)
	req.agent = strings.TrimSpace(values.Agent)
	req.round = strings.TrimSpace(values.Round)
	req.artifactDir = strings.TrimSpace(values.ArtifactDir)
	req.model = strings.TrimSpace(values.Model)
	if values.MaxRounds > 0 {
		req.maxRounds = values.MaxRounds
	}
	req.untilClean = values.UntilClean
	return req
}

func defaultCollectInteractiveInput(ctx context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
	return roundtui.CollectInput(ctx, req, os.Stdin, os.Stderr)
}

func loadRememberedInteractiveDefaults(ctx context.Context, homeDir string) (store.InteractiveDefaults, error) {
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		return store.InteractiveDefaults{}, err
	}
	defer func() {
		_ = runStore.Close()
	}()
	return runStore.InteractiveDefaults(ctx)
}

func rememberInteractiveDefaults(ctx context.Context, runStore *store.Store, req commandRequest) error {
	defaults := store.InteractiveDefaults{PRNumber: req.pr}
	if req.name == "resolve" || req.name == "watch" {
		defaults.Agent = req.agent
	}
	return runStore.RememberInteractiveDefaults(ctx, defaults)
}

func defaultSuggestCurrentPullRequest(ctx context.Context, gitRoot string) (string, error) {
	if strings.TrimSpace(gitRoot) == "" {
		return "", nil
	}
	cmd := exec.CommandContext(ctx, "gh", "pr", "view", "--json", "number", "--jq", ".number")
	cmd.Dir = gitRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func runOperationalCommand(ctx context.Context, name string, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage(name))
		return exitOK
	}

	loadedConfig, err := roundconfig.Load(roundconfig.LoadOptions{})
	if err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}

	req, err := parseOperationalCommand(name, args, loadedConfig.Config)
	if err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}

	req, err = maybeCollectInteractiveInput(ctx, req, loadedConfig, stderr)
	if err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}
	if err := validateCommandRequest(req); err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}

	artifactDir, err := roundconfig.ValidateArtifactDirectory(req.artifactDir, loadedConfig.GitRoot, loadedConfig.HomeDir)
	if err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}
	req.artifactDir = artifactDir

	preflightResult, err := runCommandPreflight(ctx, req, loadedConfig)
	if err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}
	switch req.name {
	case "fetch":
		return runFetchCommand(ctx, req, loadedConfig, preflightResult, stdout, stderr)
	case "resolve":
		return runResolveCommand(ctx, req, loadedConfig, preflightResult, stdout, stderr)
	case "watch":
		return runWatchCommand(ctx, req, loadedConfig, preflightResult, stdout, stderr)
	}

	fmt.Fprintf(stderr, "%s: %s command input accepted, but execution is not wired in this MVP slice\n", app.Name, req.name)
	fmt.Fprintf(stderr, "Artifact Directory: %s\n", req.artifactDir)
	fmt.Fprintf(stderr, "Git: %s@%s on %s (%d unpushed commit(s))\n", preflightResult.Git.Branch, preflightResult.Git.HEAD, preflightResult.Git.Root, preflightResult.Git.UnpushedCommits)
	fmt.Fprintf(stderr, "Open Pull Request: #%s %s %s\n", preflightResult.PullRequest.Number, preflightResult.PullRequest.HeadRepository, preflightResult.PullRequest.HeadBranch)
	if preflightResult.PushPlan.Enabled {
		fmt.Fprintf(stderr, "Final Push target: git push %s HEAD:%s\n", preflightResult.PushPlan.Remote, preflightResult.PushPlan.Branch)
	}
	fmt.Fprintln(stderr, "Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.")
	return exitRunFailed
}

func runFetchCommand(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, stdout, stderr io.Writer) int {
	runStore, err := store.Open(ctx, loaded.HomeDir)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	defer func() {
		_ = runStore.Close()
	}()

	run, err := runStore.CreateFetchRun(ctx, store.CreateRunRequest{
		HeadRepository: preflightResult.PullRequest.HeadRepository,
		HeadBranch:     preflightResult.PullRequest.HeadBranch,
		BaseRepository: preflightResult.PullRequest.BaseRepository,
		PRNumber:       preflightResult.PullRequest.Number,
		GitRoot:        preflightResult.Git.Root,
		LocalBranch:    preflightResult.Git.Branch,
		HeadSHA:        preflightResult.Git.HEAD,
		ArtifactDir:    req.artifactDir,
	})
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	if err := rememberInteractiveDefaults(ctx, runStore, req); err != nil {
		markRunFailed(ctx, runStore, run.ID)
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}
	printLiveRunView(stderr, req, loaded, preflightResult, run.ID, "FetchingIssues", nil, []string{"Fetching Review Source issues..."})

	items, err := fetchReviewItems(ctx, reviewsource.FetchRequest{
		Source:          req.source,
		PRNumber:        preflightResult.PullRequest.Number,
		BaseRepository:  preflightResult.PullRequest.BaseRepository,
		HeadRepository:  preflightResult.PullRequest.HeadRepository,
		HeadBranch:      preflightResult.PullRequest.HeadBranch,
		HeadSHA:         preflightResult.Git.HEAD,
		IncludeNitpicks: loaded.Config.ReviewSource.IncludeNitpicks,
	})
	if err != nil {
		markRunFailed(ctx, runStore, run.ID)
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}

	roundNumber, err := fetchRoundNumber(req.round)
	if err != nil {
		markRunFailed(ctx, runStore, run.ID)
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}
	roundResult, err := rounds.PersistRound(ctx, rounds.PersistRequest{
		ArtifactDir:    req.artifactDir,
		Source:         req.source,
		PRNumber:       preflightResult.PullRequest.Number,
		HeadRepository: preflightResult.PullRequest.HeadRepository,
		HeadBranch:     preflightResult.PullRequest.HeadBranch,
		HeadSHA:        preflightResult.Git.HEAD,
		Round:          roundNumber,
		Items:          items,
	})
	if err != nil {
		markRunFailed(ctx, runStore, run.ID)
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}

	completed, err := runStore.CompleteRun(ctx, run.ID, store.StateFetched)
	if err != nil {
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}

	fmt.Fprintf(stdout, "Fetch Run %s reached Fetched.\n", completed.ID)
	fmt.Fprintf(stdout, "Run Database: %s\n", store.DatabasePath(loaded.HomeDir))
	fmt.Fprintf(stdout, "Artifact Directory: %s\n", req.artifactDir)
	fmt.Fprintf(stdout, "Round: %03d\n", roundResult.Round)
	fmt.Fprintf(stdout, "Review Issues: %d\n", len(roundResult.IssuePaths))
	fmt.Fprintln(stdout, "Roundfix did not start an Agent, commit, or push.")
	return exitOK
}

type resolveBatchPlan struct {
	roundNumber int
	selection   rounds.SelectResult
	plan        rounds.BatchPlan
	runtime     agent.RuntimeSpec
}

type resolveBatchResult struct {
	Remaining int
}

func runResolveCommand(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, _ io.Writer, stderr io.Writer) int {
	resolvePlan, err := prepareResolveBatch(ctx, req, loaded, preflightResult)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	if err := runAgentRuntime.Probe(ctx, resolvePlan.runtime); err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}

	runStore, err := store.Open(ctx, loaded.HomeDir)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	defer func() {
		_ = runStore.Close()
	}()
	run, err := createOperationalRun(ctx, runStore, store.KindResolve, preflightResult, req.artifactDir)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	if err := rememberInteractiveDefaults(ctx, runStore, req); err != nil {
		markRunFailed(ctx, runStore, run.ID)
		printResolveRunFailure(err, stderr)
		return exitRunFailed
	}
	printLiveRunView(stderr, req, loaded, preflightResult, run.ID, "ResolvingWithAgent", resolvePlan.selection.Issues, []string{"Agent and verification output will stream below."})

	if _, err := executeResolveBatch(ctx, req, loaded, preflightResult, run.ID, resolvePlan, stderr); err != nil {
		if isStopRequest(ctx, err) {
			return completeStoppedRun(runStore, run.ID, req, preflightResult, stderr)
		}
		markRunFailed(ctx, runStore, run.ID)
		printResolveRunFailure(err, stderr)
		return exitRunFailed
	}

	completed, err := runStore.CompleteRun(ctx, run.ID, store.StateClean)
	if err != nil {
		printResolveRunFailureAfterBatchCommit(err, stderr)
		return exitRunFailed
	}
	fmt.Fprintf(stderr, "Resolve Run %s reached %s.\n", completed.ID, completed.State)
	return exitOK
}

func prepareResolveBatch(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result) (resolveBatchPlan, error) {
	roundNumber, err := resolveRoundNumber(req.round)
	if err != nil {
		return resolveBatchPlan{}, err
	}
	selection, err := rounds.SelectCompatibleIssues(ctx, rounds.SelectRequest{
		ArtifactDir:    req.artifactDir,
		PRNumber:       preflightResult.PullRequest.Number,
		HeadRepository: preflightResult.PullRequest.HeadRepository,
		HeadBranch:     preflightResult.PullRequest.HeadBranch,
		Round:          roundNumber,
	})
	if err != nil {
		return resolveBatchPlan{}, err
	}
	plan, err := rounds.PlanBatches(rounds.BatchRequest{
		Issues:    selection.Issues,
		BatchSize: loaded.Config.Resolve.BatchSize,
	})
	if err != nil {
		return resolveBatchPlan{}, err
	}
	if len(plan.Batches) == 0 {
		return resolveBatchPlan{}, fmt.Errorf("no Batch assignments were produced for selected Compatible Artifacts")
	}
	runtime, err := agent.RuntimeFor(agent.RuntimeOptions{
		Agent:           req.agent,
		CommandOverride: req.agentCmd,
		Model:           req.model,
	})
	if err != nil {
		return resolveBatchPlan{}, err
	}
	return resolveBatchPlan{
		roundNumber: roundNumber,
		selection:   selection,
		plan:        plan,
		runtime:     runtime,
	}, nil
}

func executeResolveBatch(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, runID string, resolvePlan resolveBatchPlan, stderr io.Writer) (resolveBatchResult, error) {
	batch := resolvePlan.plan.Batches[0]
	prompt := agent.BuildPrompt(agent.PromptRequest{
		RunID:        runID,
		Batch:        batch,
		Agent:        req.agent,
		Model:        resolvePlan.runtime.Model,
		ArtifactDir:  req.artifactDir,
		GitRoot:      preflightResult.Git.Root,
		Verification: loaded.Config.Defaults.Verification,
	})
	logPath := agent.LogPath(req.artifactDir, runID, batch.Number)

	fmt.Fprintf(stderr, "%s: resolve selected %d downloaded Unresolved Review Issue(s) from %d Compatible Artifact Round(s), assigned %d newest occurrence(s) into %d Batch(es), and associated %d older duplicate occurrence(s).\n", app.Name, len(resolvePlan.selection.Issues), len(resolvePlan.selection.Rounds), countBatchIssues(resolvePlan.plan.Batches), len(resolvePlan.plan.Batches), len(resolvePlan.plan.Duplicates))
	fmt.Fprintf(stderr, "Run: %s\n", runID)
	fmt.Fprintf(stderr, "Artifact Directory: %s\n", req.artifactDir)
	fmt.Fprintf(stderr, "Round scope: %s\n", formatRoundScope(resolvePlan.roundNumber))
	fmt.Fprintf(stderr, "Open Pull Request: #%s %s %s\n", preflightResult.PullRequest.Number, preflightResult.PullRequest.HeadRepository, preflightResult.PullRequest.HeadBranch)
	fmt.Fprintf(stderr, "Agent: %s\n", resolvePlan.runtime.DisplayName)
	fmt.Fprintf(stderr, "Batch: %03d (%d Review Issue(s))\n", batch.Number, len(batch.Issues))
	fmt.Fprintf(stderr, "Agent log: %s\n", logPath)

	if _, err := runAgentRuntime.Run(ctx, agent.ExecuteRequest{
		Runtime:      resolvePlan.runtime,
		RunID:        runID,
		Batch:        batch,
		LogPath:      logPath,
		Prompt:       prompt,
		ArtifactDir:  req.artifactDir,
		GitRoot:      preflightResult.Git.Root,
		Verification: loaded.Config.Defaults.Verification,
	}, stderr); err != nil {
		if isStopRequest(ctx, err) {
			return resolveBatchResult{}, err
		}
		_ = agent.MarkBatchFailed(batch)
		return resolveBatchResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return resolveBatchResult{}, err
	}
	if err := agent.ValidateAssignedIssuesTerminal(batch); err != nil {
		_ = agent.MarkBatchFailed(batch)
		return resolveBatchResult{}, err
	}

	fmt.Fprintln(stderr, "Agent Batch reached terminal local Review Issue statuses.")
	if err := runVerificationGate.Verify(ctx, daemon.VerifyRequest{
		WorkDir: preflightResult.Git.Root,
		Command: loaded.Config.Defaults.Verification,
		Stream:  stderr,
	}); err != nil {
		_ = agent.MarkBatchFailed(batch)
		return resolveBatchResult{}, err
	}
	fmt.Fprintf(stderr, "Verification command passed: %s\n", loaded.Config.Defaults.Verification)

	markedDuplicates, err := rounds.MarkDuplicatedAfterTerminal(ctx, resolvePlan.plan.Duplicates)
	if err != nil {
		return resolveBatchResult{}, err
	}
	if markedDuplicates > 0 {
		fmt.Fprintf(stderr, "Marked %d older duplicate Review Issue occurrence(s) as duplicated.\n", markedDuplicates)
	}

	if loaded.Config.Defaults.AutoCommit {
		message := daemon.BatchCommitMessage(batch.Number)
		if err := createBatchCommit.Commit(ctx, daemon.CommitRequest{
			WorkDir: preflightResult.Git.Root,
			Message: message,
		}); err != nil {
			return resolveBatchResult{}, err
		}
		fmt.Fprintf(stderr, "Batch commit created: %s\n", message)
	} else {
		fmt.Fprintln(stderr, "Auto-commit disabled; no Batch commit created.")
	}

	sourceIssues, err := terminalAssignedSourceIssues(batch)
	if err != nil {
		return resolveBatchResult{}, err
	}
	if len(sourceIssues) > 0 {
		if err := resolveReviewSourceIssues(ctx, reviewsource.ResolveRequest{
			Source:         req.source,
			PRNumber:       preflightResult.PullRequest.Number,
			BaseRepository: preflightResult.PullRequest.BaseRepository,
			Issues:         sourceIssues,
		}); err != nil {
			return resolveBatchResult{}, err
		}
		fmt.Fprintf(stderr, "Resolved %d Review Source thread(s).\n", len(sourceIssues))
	}

	remaining, err := remainingUnresolvedIssues(ctx, req.artifactDir, preflightResult.PullRequest)
	if err != nil {
		return resolveBatchResult{}, err
	}
	if remaining > 0 {
		fmt.Fprintf(stderr, "Final Push blocked: %d Unresolved Review Issue(s) remain.\n", remaining)
	} else if err := maybeRunFinalPush(ctx, loaded, preflightResult, loaded.Config.Defaults.AutoCommit, stderr); err != nil {
		return resolveBatchResult{}, err
	}
	return resolveBatchResult{Remaining: remaining}, nil
}

func runWatchCommand(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, _ io.Writer, stderr io.Writer) int {
	runtime, err := agent.RuntimeFor(agent.RuntimeOptions{
		Agent:           req.agent,
		CommandOverride: req.agentCmd,
		Model:           req.model,
	})
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	if err := runAgentRuntime.Probe(ctx, runtime); err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}

	runStore, err := store.Open(ctx, loaded.HomeDir)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	defer func() {
		_ = runStore.Close()
	}()
	run, err := createOperationalRun(ctx, runStore, store.KindWatch, preflightResult, req.artifactDir)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	if err := rememberInteractiveDefaults(ctx, runStore, req); err != nil {
		markRunFailed(ctx, runStore, run.ID)
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}

	fmt.Fprintf(stderr, "Watch Run: %s\n", run.ID)
	fmt.Fprintf(stderr, "Open Pull Request: #%s %s %s\n", preflightResult.PullRequest.Number, preflightResult.PullRequest.HeadRepository, preflightResult.PullRequest.HeadBranch)
	fmt.Fprintf(stderr, "Review Source: %s\n", req.source)
	fmt.Fprintf(stderr, "Agent: %s\n", runtime.DisplayName)
	fmt.Fprintf(stderr, "Max Rounds: %d\n", req.maxRounds)
	printLiveRunView(stderr, req, loaded, preflightResult, run.ID, "WaitingForReview", nil, []string{"Waiting for Review Source status..."})

	result, err := watch.Run(ctx, watch.Request{
		PRNumber:       preflightResult.PullRequest.Number,
		HeadSHA:        preflightResult.Git.HEAD,
		UntilClean:     req.untilClean,
		MaxRounds:      req.maxRounds,
		PollInterval:   loaded.Config.Watch.PollInterval,
		QuietPeriod:    loaded.Config.Watch.QuietPeriod,
		ReviewTimeout:  loaded.Config.Watch.ReviewTimeout,
		BudgetEnabled:  loaded.Config.Budget.Enabled,
		MaxRunDuration: loaded.Config.Budget.MaxRunDuration,
	}, watch.Dependencies{
		StatusSource: watch.StatusFunc(func(ctx context.Context, _ watch.StatusRequest) (watch.Status, error) {
			status, err := watchReviewStatus(ctx, reviewsource.WatchStatusRequest{
				Source:         req.source,
				PRNumber:       preflightResult.PullRequest.Number,
				BaseRepository: preflightResult.PullRequest.BaseRepository,
				HeadRepository: preflightResult.PullRequest.HeadRepository,
				HeadBranch:     preflightResult.PullRequest.HeadBranch,
				HeadSHA:        preflightResult.Git.HEAD,
			})
			if err == nil {
				fmt.Fprintf(stderr, "Review Source status: %s\n", status.State)
			}
			return status, err
		}),
		Fetcher: watch.FetchFunc(func(ctx context.Context, _ int) (watch.FetchResult, error) {
			return fetchWatchRound(ctx, req, loaded, preflightResult, stderr)
		}),
		Resolver: watch.ResolveFunc(func(ctx context.Context) (watch.ResolveResult, error) {
			return resolveWatchBatches(ctx, req, loaded, preflightResult, runtime, run.ID, stderr)
		}),
		Clock:   watchClock,
		Sleeper: watchSleeper,
	})
	stopped := isStopRequest(ctx, err)
	if stopped {
		result.Outcome = store.StateStopped
	} else if err != nil && result.Outcome == "" {
		result.Outcome = store.StateFailed
	}

	terminal := result.Outcome
	if terminal == "" {
		terminal = store.StateFailed
	}
	completeCtx := ctx
	var completeCancel context.CancelFunc
	if stopped {
		completeCtx, completeCancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer completeCancel()
	}
	completed, completeErr := runStore.CompleteRun(completeCtx, run.ID, terminal)
	if completeErr != nil {
		printRunFailure(req.name, completeErr, stderr)
		return exitRunFailed
	}

	fmt.Fprintf(stderr, "Watch Run %s reached %s after %d Round(s).\n", completed.ID, completed.State, result.Rounds)
	if result.Outcome == store.StateMaxRoundsReached && result.Remaining > 0 {
		fmt.Fprintf(stderr, "MaxRoundsReached with %d Unresolved Review Issue(s) remaining.\n", result.Remaining)
	}
	if result.Outcome == store.StateTimedOut {
		fmt.Fprintf(stderr, "Review Source timed out. To request another CodeRabbit review manually, comment: %s\n", result.ManualReviewCommand)
	}
	if stopped {
		printStopSummary(req, preflightResult, stderr)
		return exitOK
	}
	if err != nil {
		printWatchRunFailure(err, stderr)
		return exitRunFailed
	}
	return exitForWatchOutcome(result.Outcome)
}

func createOperationalRun(ctx context.Context, runStore *store.Store, kind string, preflightResult preflight.Result, artifactDir string) (store.Run, error) {
	return runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           kind,
		HeadRepository: preflightResult.PullRequest.HeadRepository,
		HeadBranch:     preflightResult.PullRequest.HeadBranch,
		BaseRepository: preflightResult.PullRequest.BaseRepository,
		PRNumber:       preflightResult.PullRequest.Number,
		GitRoot:        preflightResult.Git.Root,
		LocalBranch:    preflightResult.Git.Branch,
		HeadSHA:        preflightResult.Git.HEAD,
		ArtifactDir:    artifactDir,
	})
}

func fetchWatchRound(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, stderr io.Writer) (watch.FetchResult, error) {
	items, err := fetchReviewItems(ctx, reviewsource.FetchRequest{
		Source:          req.source,
		PRNumber:        preflightResult.PullRequest.Number,
		BaseRepository:  preflightResult.PullRequest.BaseRepository,
		HeadRepository:  preflightResult.PullRequest.HeadRepository,
		HeadBranch:      preflightResult.PullRequest.HeadBranch,
		HeadSHA:         preflightResult.Git.HEAD,
		IncludeNitpicks: loaded.Config.ReviewSource.IncludeNitpicks,
	})
	if err != nil {
		return watch.FetchResult{}, err
	}
	roundResult, err := rounds.PersistRound(ctx, rounds.PersistRequest{
		ArtifactDir:    req.artifactDir,
		Source:         req.source,
		PRNumber:       preflightResult.PullRequest.Number,
		HeadRepository: preflightResult.PullRequest.HeadRepository,
		HeadBranch:     preflightResult.PullRequest.HeadBranch,
		HeadSHA:        preflightResult.Git.HEAD,
		Items:          items,
	})
	if err != nil {
		return watch.FetchResult{}, err
	}
	fmt.Fprintf(stderr, "Fetched Round %03d with %d Review Issue(s).\n", roundResult.Round, len(roundResult.IssuePaths))
	return watch.FetchResult{Round: roundResult.Round, Issues: len(roundResult.IssuePaths)}, nil
}

func resolveWatchBatches(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, runtime agent.RuntimeSpec, runID string, stderr io.Writer) (watch.ResolveResult, error) {
	progress := false
	for {
		resolvePlan, err := prepareResolveBatch(ctx, req, loaded, preflightResult)
		if err != nil {
			var noArtifacts rounds.NoCompatibleArtifactsError
			if errors.As(err, &noArtifacts) {
				return watch.ResolveResult{Remaining: 0, Progress: progress}, nil
			}
			return watch.ResolveResult{}, err
		}
		resolvePlan.runtime = runtime
		result, err := executeResolveBatch(ctx, req, loaded, preflightResult, runID, resolvePlan, stderr)
		if err != nil {
			return watch.ResolveResult{}, err
		}
		progress = true
		if result.Remaining == 0 {
			return watch.ResolveResult{Remaining: 0, Progress: true}, nil
		}
	}
}

func exitForWatchOutcome(outcome string) int {
	switch outcome {
	case store.StateClean, store.StateMaxRoundsReached, store.StateStopped:
		return exitOK
	case store.StateBudgetExceeded, store.StateTimedOut, store.StateFailed:
		return exitRunFailed
	default:
		return exitRunFailed
	}
}

func completeStoppedRun(runStore *store.Store, runID string, req commandRequest, preflightResult preflight.Result, stderr io.Writer) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	completed, err := runStore.CompleteRun(ctx, runID, store.StateStopped)
	if err != nil {
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}
	fmt.Fprintf(stderr, "%s Run %s reached %s.\n", commandDisplayName(req.name), completed.ID, completed.State)
	printStopSummary(req, preflightResult, stderr)
	return exitOK
}

func printStopSummary(req commandRequest, preflightResult preflight.Result, stderr io.Writer) {
	fmt.Fprintln(stderr, "Stop Request preserved local work and stopped before any later verification, commit, push, fetch, or Review Source mutation.")
	changes, err := inspectChangedPaths(context.Background(), preflightResult.Git.Root)
	if err != nil {
		fmt.Fprintf(stderr, "Changed paths after Stop Request: unavailable: %v\n", err)
		return
	}
	if len(changes) == 0 {
		fmt.Fprintln(stderr, "Changed paths after Stop Request: none")
		return
	}
	fmt.Fprintln(stderr, "Changed paths after Stop Request:")
	for _, change := range changes {
		fmt.Fprintf(stderr, "  %s %s\n", change.Status, change.Path)
	}
	if req.artifactDir != "" {
		fmt.Fprintf(stderr, "Artifact Directory preserved: %s\n", req.artifactDir)
	}
}

func isStopRequest(ctx context.Context, err error) bool {
	return agent.IsStopError(err) || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || (ctx != nil && ctx.Err() != nil)
}

func defaultInspectChangedPaths(ctx context.Context, gitRoot string) ([]preflight.ChangedPath, error) {
	state, err := preflight.InspectGit(ctx, gitRoot, preflight.ExecGitRunner{})
	if err != nil {
		return nil, err
	}
	return state.Dirty, nil
}

func printLiveRunView(stderr io.Writer, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, runID string, pipelineState string, issues []rounds.Issue, console []string) {
	fmt.Fprint(stderr, roundtui.RenderLiveRunView(roundtui.LiveRunView{
		Repository:    preflightResult.PullRequest.HeadRepository,
		PRNumber:      preflightResult.PullRequest.Number,
		HeadBranch:    preflightResult.PullRequest.HeadBranch,
		ReviewSource:  displayReviewSource(req.source),
		Agent:         displayAgent(req.agent),
		HEAD:          preflightResult.Git.HEAD,
		RunID:         runID,
		PipelineState: pipelineState,
		BudgetState:   formatBudgetState(loaded.Config),
		GitState:      formatGitState(preflightResult.Git),
		CurrentRound:  0,
		MaxRounds:     req.maxRounds,
		AutoCommit:    loaded.Config.Defaults.AutoCommit,
		AutoPush:      preflightResult.PushPlan.Enabled,
		LastPush:      lastPushState(preflightResult.PushPlan),
		Issues:        issues,
		Console:       console,
	}))
}

func displayReviewSource(source string) string {
	switch source {
	case reviewsource.SourceCodeRabbit:
		return "CodeRabbit"
	case "":
		return "-"
	default:
		return source
	}
}

func displayAgent(agentName string) string {
	switch agentName {
	case "codex":
		return "Codex"
	case "claude":
		return "Claude"
	case "opencode":
		return "OpenCode"
	case "":
		return "-"
	default:
		return agentName
	}
}

func formatBudgetState(config roundconfig.Config) string {
	if !config.Budget.Enabled {
		return "off"
	}
	return fmt.Sprintf("0s / %s", config.Budget.MaxRunDuration)
}

func formatGitState(git preflight.GitState) string {
	dirty := "clean"
	if len(git.Dirty) > 0 {
		dirty = fmt.Sprintf("%d changed path(s)", len(git.Dirty))
	}
	upstream := "no upstream"
	if git.UpstreamRemote != "" || git.UpstreamBranch != "" {
		upstream = strings.Trim(git.UpstreamRemote+"/"+git.UpstreamBranch, "/")
	}
	return fmt.Sprintf("%s, %d unpushed commit(s), upstream %s", dirty, git.UnpushedCommits, upstream)
}

func lastPushState(plan preflight.PushPlan) string {
	if !plan.Enabled {
		return "disabled"
	}
	return "pending"
}

func commandDisplayName(name string) string {
	if name == "" {
		return "Run"
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func parseOperationalCommand(name string, args []string, config roundconfig.Config) (commandRequest, error) {
	req := commandRequest{
		name:        name,
		source:      config.ReviewSource.Name,
		agent:       config.Defaults.Agent,
		round:       "all",
		untilClean:  config.Watch.UntilClean,
		maxRounds:   config.Watch.MaxRounds,
		artifactDir: config.Defaults.ArtifactDir,
		model:       config.Defaults.Model,
	}
	if name == "fetch" {
		req.round = "auto"
	}

	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&req.pr, "pr", "", "Open Pull Request number")
	fs.BoolVar(&req.noInput, "no-input", false, "Fail instead of opening Interactive Input")
	fs.BoolVar(&req.interactive, "interactive", false, "Open Interactive Input before starting")
	fs.StringVar(&req.artifactDir, "artifact-dir", req.artifactDir, "Artifact Directory")
	fs.StringVar(&req.baseRepo, "base-repo", "", "Base repository, owner/name")
	fs.StringVar(&req.headRepo, "head-repo", "", "Head Repository, owner/name")
	fs.StringVar(&req.headBranch, "head-branch", "", "PR Head Branch")

	switch name {
	case "fetch":
		fs.StringVar(&req.source, "source", req.source, "Review Source")
		fs.StringVar(&req.round, "round", "auto", "Round number or auto")
	case "resolve":
		fs.StringVar(&req.agent, "agent", req.agent, "Agent runtime")
		fs.StringVar(&req.model, "model", req.model, "Agent model override")
		fs.StringVar(&req.agentCmd, "agent-command", "", "Agent command override")
		fs.StringVar(&req.round, "round", "all", "Round number or all")
	case "watch":
		fs.StringVar(&req.source, "source", req.source, "Review Source")
		fs.StringVar(&req.agent, "agent", req.agent, "Agent runtime")
		fs.StringVar(&req.model, "model", req.model, "Agent model override")
		fs.StringVar(&req.agentCmd, "agent-command", "", "Agent command override")
		fs.BoolVar(&req.untilClean, "until-clean", req.untilClean, "Repeat until no Unresolved Review Issues remain")
		fs.IntVar(&req.maxRounds, "max-rounds", req.maxRounds, "Maximum Review Source rounds")
	default:
		return req, validationError{message: fmt.Sprintf("unknown command %q", name)}
	}

	err := fs.Parse(args)
	if err != nil {
		return req, validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return req, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}

	return req, nil
}

func validateCommandRequest(req commandRequest) error {
	if req.pr == "" {
		return missingRequiredFlag(req, "pr")
	}
	if err := validatePositiveInt("pr", req.pr); err != nil {
		return err
	}
	if req.source != "" && req.source != "coderabbit" {
		return validationError{message: fmt.Sprintf("unsupported Review Source %q; supported value: coderabbit", req.source)}
	}

	switch req.name {
	case "fetch":
		return validateRoundValue(req.round, true)
	case "resolve":
		if req.agent == "" {
			return missingRequiredFlag(req, "agent")
		}
		if err := validateAgent(req.agent); err != nil {
			return err
		}
		return validateRoundValue(req.round, false)
	case "watch":
		if req.agent == "" {
			return missingRequiredFlag(req, "agent")
		}
		if err := validateAgent(req.agent); err != nil {
			return err
		}
		if req.maxRounds < 1 {
			return validationError{message: "--max-rounds must be greater than 0"}
		}
	}
	return nil
}

func missingRequiredFlag(req commandRequest, flagName string) error {
	if req.noInput {
		return validationError{message: fmt.Sprintf("missing required --%s because --no-input disables Interactive Input", flagName)}
	}
	if req.inputShown {
		return validationError{message: fmt.Sprintf("Interactive Input did not collect required --%s; enter a value or pass --%s", flagName, flagName)}
	}
	return validationError{message: fmt.Sprintf("missing required --%s; pass --%s or use --interactive when Interactive Input is available", flagName, flagName)}
}

func validatePositiveInt(flagName string, value string) error {
	number, err := strconv.Atoi(value)
	if err != nil || number < 1 {
		return validationError{message: fmt.Sprintf("--%s must be a positive integer", flagName)}
	}
	return nil
}

func validateRoundValue(value string, allowAuto bool) error {
	if value == "all" && !allowAuto {
		return nil
	}
	if value == "auto" && allowAuto {
		return nil
	}
	if err := validatePositiveInt("round", value); err != nil {
		expected := "positive integer"
		if allowAuto {
			expected += " or auto"
		} else {
			expected += " or all"
		}
		return validationError{message: fmt.Sprintf("--round must be a %s", expected)}
	}
	return nil
}

func validateAgent(agent string) error {
	switch agent {
	case "codex", "claude", "opencode":
		return nil
	default:
		return validationError{message: fmt.Sprintf("unsupported Agent %q; supported values: codex, claude, opencode", agent)}
	}
}

func defaultRunCommandPreflight(ctx context.Context, req commandRequest, loaded roundconfig.Loaded) (preflight.Result, error) {
	return preflight.Run(ctx, preflight.Request{
		Command:                req.name,
		WorkDir:                loaded.GitRoot,
		ArtifactDir:            req.artifactDir,
		PRNumber:               req.pr,
		ExplicitBaseRepository: req.baseRepo,
		ExplicitHeadBranch:     req.headBranch,
		ExplicitHeadRepository: req.headRepo,
		AutoCommit:             loaded.Config.Defaults.AutoCommit,
		AutoPush:               loaded.Config.Watch.AutoPush,
		PushRemote:             loaded.Config.Watch.PushRemote,
		PushBranch:             loaded.Config.Watch.PushBranch,
	})
}

func defaultFetchReviewItems(ctx context.Context, req reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
	if req.Source != reviewsource.SourceCodeRabbit {
		return nil, fmt.Errorf("unsupported Review Source %q; supported value: coderabbit", req.Source)
	}
	return coderabbit.Client{}.FetchReviews(ctx, req)
}

func defaultResolveReviewSourceIssues(ctx context.Context, req reviewsource.ResolveRequest) error {
	if req.Source != reviewsource.SourceCodeRabbit {
		return fmt.Errorf("unsupported Review Source %q; supported value: coderabbit", req.Source)
	}
	return coderabbit.Client{}.ResolveIssues(ctx, req)
}

func defaultWatchReviewStatus(ctx context.Context, req reviewsource.WatchStatusRequest) (reviewsource.WatchStatus, error) {
	if req.Source != reviewsource.SourceCodeRabbit {
		return reviewsource.WatchStatus{}, fmt.Errorf("unsupported Review Source %q; supported value: coderabbit", req.Source)
	}
	return coderabbit.Client{}.WatchStatus(ctx, req)
}

func fetchRoundNumber(value string) (int, error) {
	if value == "" || value == "auto" {
		return 0, nil
	}
	roundNumber, err := strconv.Atoi(value)
	if err != nil || roundNumber < 1 {
		return 0, fmt.Errorf("--round must be a positive integer or auto")
	}
	return roundNumber, nil
}

func resolveRoundNumber(value string) (int, error) {
	if value == "" || value == "all" {
		return 0, nil
	}
	roundNumber, err := strconv.Atoi(value)
	if err != nil || roundNumber < 1 {
		return 0, fmt.Errorf("--round must be a positive integer or all")
	}
	return roundNumber, nil
}

func formatRoundScope(roundNumber int) string {
	if roundNumber == 0 {
		return "all"
	}
	return fmt.Sprintf("%03d", roundNumber)
}

func countBatchIssues(batches []rounds.Batch) int {
	count := 0
	for _, batch := range batches {
		count += len(batch.Issues)
	}
	return count
}

func terminalAssignedSourceIssues(batch rounds.Batch) ([]reviewsource.ResolvedIssue, error) {
	issues := make([]reviewsource.ResolvedIssue, 0, len(batch.Issues))
	for _, assigned := range batch.Issues {
		issue, err := rounds.ParseIssue(assigned.Path)
		if err != nil {
			return nil, err
		}
		if issue.Status != rounds.StatusResolved && issue.Status != rounds.StatusInvalid {
			continue
		}
		if strings.TrimSpace(issue.SourceRef) == "" {
			continue
		}
		issues = append(issues, reviewsource.ResolvedIssue{
			FilePath:  issue.Path,
			Status:    issue.Status,
			SourceRef: issue.SourceRef,
		})
	}
	return issues, nil
}

func remainingUnresolvedIssues(ctx context.Context, artifactDir string, pullRequest preflight.PullRequest) (int, error) {
	selection, err := rounds.SelectCompatibleIssues(ctx, rounds.SelectRequest{
		ArtifactDir:    artifactDir,
		PRNumber:       pullRequest.Number,
		HeadRepository: pullRequest.HeadRepository,
		HeadBranch:     pullRequest.HeadBranch,
	})
	if err != nil {
		var noArtifacts rounds.NoCompatibleArtifactsError
		if errors.As(err, &noArtifacts) {
			return 0, nil
		}
		return 0, err
	}
	return len(selection.Issues), nil
}

func maybeRunFinalPush(ctx context.Context, loaded roundconfig.Loaded, preflightResult preflight.Result, batchCommitCreated bool, stderr io.Writer) error {
	if !preflightResult.PushPlan.Enabled {
		fmt.Fprintln(stderr, "Final Push skipped: auto-push disabled or no push target configured.")
		return nil
	}
	if preflightResult.PushPlan.Force {
		return errors.New("Final Push rejected: force-push is not allowed in the MVP")
	}
	if !loaded.Config.Defaults.AutoCommit {
		fmt.Fprintln(stderr, "Final Push skipped: auto-commit disabled.")
		return nil
	}
	if preflightResult.Git.UnpushedCommits == 0 && !batchCommitCreated {
		fmt.Fprintln(stderr, "Final Push skipped: no local commits are waiting for the PR Head Branch.")
		return nil
	}
	if err := runFinalPush.Push(ctx, daemon.PushRequest{
		WorkDir: preflightResult.Git.Root,
		Remote:  preflightResult.PushPlan.Remote,
		Branch:  preflightResult.PushPlan.Branch,
	}); err != nil {
		return err
	}
	fmt.Fprintf(stderr, "Final Push completed: git push %s HEAD:%s\n", preflightResult.PushPlan.Remote, preflightResult.PushPlan.Branch)
	return nil
}

func markRunFailed(ctx context.Context, runStore *store.Store, runID string) {
	_, _ = runStore.CompleteRun(ctx, runID, store.StateFailed)
}

func commandUsage(name string) string {
	switch name {
	case "fetch":
		return `Usage:
  roundfix fetch --source coderabbit --pr <number> [--round <number|auto>] [--no-input]

Options:
  --source       Review Source. Supported: coderabbit
  --pr           Open Pull Request number
  --round        Round number or auto
  --artifact-dir Artifact Directory
  --base-repo    Explicit base repository, owner/name
  --head-repo    Explicit Head Repository, owner/name
  --head-branch  Explicit PR Head Branch
  --interactive  Open Interactive Input before starting
  --no-input     Fail instead of opening Interactive Input
`
	case "resolve":
		return `Usage:
  roundfix resolve --pr <number> --agent <agent> [--round <number|all>] [--no-input]

Options:
  --pr           Open Pull Request number
  --agent        Agent runtime. Supported: codex, claude, opencode
  --model        Agent model override
  --agent-command Agent command override
  --round        Round number or all
  --artifact-dir Artifact Directory
  --base-repo    Explicit base repository, owner/name
  --head-repo    Explicit Head Repository, owner/name
  --head-branch  Explicit PR Head Branch
  --interactive  Open Interactive Input before starting
  --no-input     Fail instead of opening Interactive Input
`
	case "watch":
		return `Usage:
  roundfix watch --source coderabbit --pr <number> --agent <agent> [--until-clean] [--max-rounds <number>] [--no-input]

Options:
  --source       Review Source. Supported: coderabbit
  --pr           Open Pull Request number
  --agent        Agent runtime. Supported: codex, claude, opencode
  --model        Agent model override
  --agent-command Agent command override
  --until-clean  Repeat until no Unresolved Review Issues remain
  --max-rounds   Maximum Review Source rounds
  --artifact-dir Artifact Directory
  --base-repo    Explicit base repository, owner/name
  --head-repo    Explicit Head Repository, owner/name
  --head-branch  Explicit PR Head Branch
  --interactive  Open Interactive Input before starting
  --no-input     Fail instead of opening Interactive Input
`
	case "skills":
		return `Usage:
  roundfix skills check
  roundfix skills install [--target <codex|claude|opencode|all>] [--dir <path>]

Commands:
  check      Validate shipped Roundfix skill artifacts
  install    Install Roundfix skills into Codex, Claude Code, or OpenCode-compatible directories

Options:
  --target   Install target. Supported: codex, claude, opencode, all
  --dir      Override the target skills directory for a single target
`
	default:
		return strings.TrimRight(usage, "\n") + "\n"
	}
}

func commandWantsHelp(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "help" {
			return true
		}
	}
	return false
}

func printPreflightFailure(name string, err error, stderr io.Writer) {
	style := styleFor(stderr)
	fmt.Fprintf(stderr, "%s\n\n", style.red(style.bold("Preflight failed")))
	var dirtyErr preflight.DirtyWorktreeError
	if errors.As(err, &dirtyErr) {
		printDirtyWorktreeFailure(dirtyErr, stderr, style)
	} else {
		fmt.Fprintf(stderr, "%s\n", style.cyan("Reason:"))
		fmt.Fprintf(stderr, "  %v\n\n", err)
	}
	printPreflightNoSideEffects(stderr, style)
	fmt.Fprintf(stderr, "%s\n", style.cyan("Usage:"))
	fmt.Fprintf(stderr, "  Run '%s %s --help' for usage.\n", app.Name, name)
}

func printDirtyWorktreeFailure(err preflight.DirtyWorktreeError, stderr io.Writer, style terminalStyle) {
	fmt.Fprintf(stderr, "%s\n", style.cyan("Reason:"))
	fmt.Fprintln(stderr, "  Dirty worktree outside the Artifact Directory.")
	if strings.TrimSpace(err.ArtifactDir) != "" {
		fmt.Fprintf(stderr, "  Artifact Directory: %s\n", style.dim(err.ArtifactDir))
	}
	fmt.Fprintln(stderr)
	fmt.Fprintf(stderr, "%s\n", style.yellow("Resolve before running Roundfix:"))
	fmt.Fprintln(stderr, "  commit, stash, or remove these changes:")
	for _, change := range err.Changes {
		fmt.Fprintf(stderr, "  %s %s\n", change.Status, change.Path)
	}
	fmt.Fprintln(stderr)
}

func printPreflightNoSideEffects(stderr io.Writer, style terminalStyle) {
	fmt.Fprintf(stderr, "%s\n", style.cyan("No side effects:"))
	fmt.Fprintln(stderr, "  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.")
	fmt.Fprintln(stderr)
}

func printRunFailure(name string, err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: %s failed after Run start: %v\n", app.Name, name, err)
	fmt.Fprintln(stderr, "Roundfix did not start an Agent, commit, or push.")
}

func printResolveRunFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: resolve failed after Run start: %v\n", app.Name, err)
	fmt.Fprintln(stderr, "Roundfix did not commit, push, or resolve Review Source threads.")
}

func printResolveRunFailureAfterBatchCommit(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: resolve failed after Run start: %v\n", app.Name, err)
	fmt.Fprintln(stderr, "Roundfix did not complete Final Push; local Batch commit and artifact changes remain for inspection.")
}

func printWatchRunFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: watch failed after Run start: %v\n", app.Name, err)
	fmt.Fprintln(stderr, "Roundfix completed the Watch Run with a terminal failure; inspect local artifacts and Run output before retrying.")
}
