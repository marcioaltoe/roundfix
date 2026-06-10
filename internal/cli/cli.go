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
	"roundfix/internal/runevent"
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
  roundfix init [--scope <project|user>]
  roundfix stop [<run-id>|--run-id <id>|--pr <number>]
  roundfix attach <run-id>
  roundfix skills check
  roundfix skills install --target <codex|claude|opencode|all>

Commands:
  init       Create User Config or Project Config
  fetch      Download review issues for an Open Pull Request
  resolve    Resolve downloaded Unresolved Review Issues
  watch      Fetch and resolve in a watched loop
  stop       Stop an Active Run and release its lock
  attach     Replay a Run's event timeline from the Run Database
  skills     Check or install Roundfix agent skills

Options:
  -h, --help      Show help
  -v, --version   Show version
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

// newEngineCollaborators is the single test seam for Run engine
// collaborators; orchestration itself lives in the daemon Run engine and
// receives them through an explicit dependencies struct.
var newEngineCollaborators = defaultEngineCollaborators
var watchReviewStatus = defaultWatchReviewStatus
var watchClock watch.Clock
var watchSleeper watch.Sleeper
var inspectChangedPaths = defaultInspectChangedPaths
var collectInteractiveInput = defaultCollectInteractiveInput
var suggestCurrentPullRequest = defaultSuggestCurrentPullRequest
var resolvePullRequestForStop = defaultResolvePullRequestForStop
var promptInitScope = defaultPromptInitScope

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
	case "-v", "--version", "version":
		fmt.Fprintf(stdout, "%s %s\n", app.Name, app.Version)
		return exitOK
	case "init":
		return runInitCommand(ctx, args[1:], stdout, stderr)
	case "stop":
		return runStopCommand(ctx, args[1:], stdout, stderr)
	case "attach":
		return runAttachCommand(ctx, args[1:], stdout, stderr)
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

func runInitCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("init"))
		return exitOK
	}
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	scope := fs.String("scope", "", "Config scope: project or user")
	force := fs.Bool("force", false, "Overwrite an existing config file")
	if err := fs.Parse(args); err != nil {
		printInitFailure(err, stderr)
		return exitPreflight
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		printInitFailure(validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}, stderr)
		return exitPreflight
	}

	selectedScope := strings.TrimSpace(*scope)
	if selectedScope == "" {
		var err error
		selectedScope, err = promptInitScope(ctx, stderr)
		if err != nil {
			printInitFailure(err, stderr)
			return exitPreflight
		}
	}

	result, err := roundconfig.Init(ctx, roundconfig.InitOptions{
		Scope: selectedScope,
		Force: *force,
	})
	if err != nil {
		printInitFailure(err, stderr)
		return exitPreflight
	}
	printInitSuccess(result, stdout)
	return exitOK
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

type stopRequest struct {
	runID      string
	pr         string
	headRepo   string
	headBranch string
}

func runStopCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("stop"))
		return exitOK
	}

	req, err := parseStopCommand(args)
	if err != nil {
		printStopFailure(err, stderr)
		return exitPreflight
	}
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{})
	if err != nil {
		printStopFailure(err, stderr)
		return exitPreflight
	}
	runStore, err := store.Open(ctx, loaded.HomeDir)
	if err != nil {
		printStopFailure(err, stderr)
		return exitPreflight
	}
	defer func() {
		_ = runStore.Close()
	}()

	run, err := stopTargetRun(ctx, req, loaded, runStore)
	if err != nil {
		printStopFailure(err, stderr)
		return exitPreflight
	}
	printStopSuccess(run, stdout)
	return exitOK
}

func parseStopCommand(args []string) (stopRequest, error) {
	req := stopRequest{}
	fs := flag.NewFlagSet("stop", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&req.runID, "run-id", "", "Run ID to stop")
	fs.StringVar(&req.runID, "run", "", "Run ID to stop")
	fs.StringVar(&req.pr, "pr", "", "Open Pull Request number")
	fs.StringVar(&req.headRepo, "head-repo", "", "Head Repository, owner/name")
	fs.StringVar(&req.headBranch, "head-branch", "", "PR Head Branch")
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
	if req.runID != "" && (req.pr != "" || req.headRepo != "" || req.headBranch != "") {
		return req, validationError{message: "--run-id cannot be combined with --pr, --head-repo, or --head-branch"}
	}
	if req.pr != "" {
		if err := validatePositiveInt("pr", req.pr); err != nil {
			return req, err
		}
	}
	if (req.headRepo == "") != (req.headBranch == "") {
		return req, validationError{message: "--head-repo and --head-branch must be used together"}
	}
	return req, nil
}

func stopTargetRun(ctx context.Context, req stopRequest, loaded roundconfig.Loaded, runStore *store.Store) (store.Run, error) {
	if req.runID != "" {
		current, found, err := runStore.Run(ctx, req.runID)
		if err != nil {
			return store.Run{}, err
		} else if !found {
			return store.Run{}, validationError{message: fmt.Sprintf("Run %q does not exist", req.runID)}
		}
		if current.State != store.StateActive {
			return store.Run{}, validationError{message: fmt.Sprintf("Run %q is already %s", req.runID, current.State)}
		}
		run, err := runStore.CompleteRun(ctx, req.runID, store.StateStopped)
		if err != nil {
			return store.Run{}, err
		}
		return run, nil
	}

	headRepo := strings.TrimSpace(req.headRepo)
	headBranch := strings.TrimSpace(req.headBranch)
	if headRepo == "" || headBranch == "" {
		pr := strings.TrimSpace(req.pr)
		if pr == "" {
			suggested, _ := suggestCurrentPullRequest(ctx, loaded.GitRoot)
			pr = strings.TrimSpace(suggested)
		}
		if pr == "" {
			return store.Run{}, validationError{message: "missing stop target; pass a Run ID, --run-id, --pr, or --head-repo with --head-branch"}
		}
		resolved, err := resolvePullRequestForStop(ctx, loaded.GitRoot, pr)
		if err != nil {
			return store.Run{}, fmt.Errorf("resolve Open Pull Request %s for stop target: %w", pr, err)
		}
		headRepo = resolved.HeadRepository
		headBranch = resolved.HeadBranch
	}

	active, found, err := runStore.ActiveRun(ctx, headRepo, headBranch)
	if err != nil {
		return store.Run{}, err
	}
	if !found {
		return store.Run{}, validationError{message: fmt.Sprintf("no Active Run exists for Head Repository %q and PR Head Branch %q", headRepo, headBranch)}
	}
	run, err := runStore.CompleteRun(ctx, active.ID, store.StateStopped)
	if err != nil {
		return store.Run{}, err
	}
	return run, nil
}

func defaultResolvePullRequestForStop(ctx context.Context, workDir string, pr string) (preflight.PullRequest, error) {
	return (preflight.GHPullRequestResolver{}).ResolvePullRequest(ctx, workDir, pr)
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

func defaultPromptInitScope(ctx context.Context, stderr io.Writer) (string, error) {
	return readInitScope(ctx, os.Stdin, stderr)
}

func readInitScope(ctx context.Context, stdin io.Reader, stderr io.Writer) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	fmt.Fprintln(stderr, "Roundfix Config Init")
	fmt.Fprintln(stderr, "Choose where to write the config.")
	fmt.Fprintln(stderr, "Press Enter to use project config.")
	fmt.Fprint(stderr, "Scope [project] (project/user): ")
	line, err := bufio.NewReader(stdin).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read init scope: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return normalizeInitScope(line)
}

func normalizeInitScope(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", roundconfig.InitScopeProject:
		return roundconfig.InitScopeProject, nil
	case roundconfig.InitScopeUser:
		return roundconfig.InitScopeUser, nil
	default:
		return "", fmt.Errorf("unsupported init scope %q; supported values: project, user", strings.TrimSpace(value))
	}
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

	printFetchSuccess(stdout, fetchSuccessView{
		RunID:          completed.ID,
		Round:          roundResult.Round,
		ReviewIssues:   len(roundResult.IssuePaths),
		RunDatabase:    store.DatabasePath(loaded.HomeDir),
		ArtifactDir:    req.artifactDir,
		ReusedRound:    roundResult.Reused,
		StartedAgent:   false,
		CreatedCommit:  false,
		CompletedPush:  false,
		ResolvedSource: false,
	})
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
	collaborators := newEngineCollaborators()
	if err := collaborators.runner.Probe(ctx, resolvePlan.runtime); err != nil {
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

	cockpitView := buildLiveRunView(req, loaded, preflightResult, run.ID, "ResolvingWithAgent", resolvePlan.selection.Issues, nil)
	for _, batch := range resolvePlan.plan.Batches {
		cockpitView.BatchSizes = append(cockpitView.BatchSizes, len(batch.Issues))
	}
	ui, err := startRunUI(ctx, cockpitView, run.ID, loaded.HomeDir, runStore, stderr)
	if err != nil {
		markRunFailed(ctx, runStore, run.ID)
		printResolveRunFailure(err, stderr)
		return exitRunFailed
	}
	defer ui.Close()

	if _, err := executeResolveCycle(ctx, req, loaded, preflightResult, run.ID, resolvePlan, collaborators, runStore, ui); err != nil {
		ui.Close()
		if isStopRequest(ctx, err) {
			return completeStoppedRun(runStore, run.ID, req, preflightResult, stderr)
		}
		markRunFailed(ctx, runStore, run.ID)
		printResolveRunFailure(err, stderr)
		return exitRunFailed
	}
	ui.Close()

	completed, err := runStore.CompleteRun(ctx, run.ID, store.StateClean)
	if err != nil {
		printResolveRunFailureAfterBatchCommit(err, stderr)
		return exitRunFailed
	}
	publishRunOutcome(ctx, runStore, completed.ID, completed.State, 0, stderr)
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

func executeResolveCycle(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, runID string, resolvePlan resolveBatchPlan, collaborators engineCollaborators, runStore *store.Store, ui *runUI) (resolveBatchResult, error) {
	fmt.Fprintf(ui.progress, "%s: resolve selected %d downloaded Unresolved Review Issue(s) from %d Compatible Artifact Round(s), assigned %d newest occurrence(s) into %d Batch(es), and associated %d older duplicate occurrence(s).\n", app.Name, len(resolvePlan.selection.Issues), len(resolvePlan.selection.Rounds), countBatchIssues(resolvePlan.plan.Batches), len(resolvePlan.plan.Batches), len(resolvePlan.plan.Duplicates))
	fmt.Fprintf(ui.progress, "Run: %s\n", runID)
	fmt.Fprintf(ui.progress, "Artifact Directory: %s\n", req.artifactDir)
	fmt.Fprintf(ui.progress, "Round scope: %s\n", formatRoundScope(resolvePlan.roundNumber))
	fmt.Fprintf(ui.progress, "Open Pull Request: #%s %s %s\n", preflightResult.PullRequest.Number, preflightResult.PullRequest.HeadRepository, preflightResult.PullRequest.HeadBranch)
	fmt.Fprintf(ui.progress, "Agent: %s\n", resolvePlan.runtime.DisplayName)

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
		return resolveBatchResult{}, err
	}

	result, cycleErr := engine.ResolveCycle(ctx, cyclePlanFrom(req, loaded, preflightResult, runID, resolvePlan))
	if cycleErr != nil {
		return resolveBatchResult{}, cycleErr
	}

	if result.Remaining > 0 {
		fmt.Fprintf(ui.progress, "Final Push blocked: %d Unresolved Review Issue(s) remain.\n", result.Remaining)
		publishPushDecision(ctx, ui.sink, runID, "blocked", fmt.Sprintf("Final Push blocked: %d Unresolved Review Issue(s) remain.", result.Remaining), result.Remaining)
	} else if err := maybeRunFinalPush(ctx, engine, ui.sink, runID, loaded, preflightResult, loaded.Config.Defaults.AutoCommit, ui.progress); err != nil {
		return resolveBatchResult{}, err
	}
	return resolveBatchResult{Remaining: result.Remaining}, nil
}

// publishPushDecision journals daemon-owned Final Push gating decisions.
func publishPushDecision(ctx context.Context, sink runevent.Sink, runID string, decision string, summary string, remaining int) {
	payload, err := json.Marshal(map[string]any{"decision": decision, "remaining": remaining})
	if err != nil {
		return
	}
	_ = sink.Publish(context.WithoutCancel(ctx), runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonPush,
		Summary: runevent.BoundSummary(summary),
		Time:    time.Now().UTC(),
		Payload: payload,
	})
}

func cyclePlanFrom(req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, runID string, resolvePlan resolveBatchPlan) daemon.CyclePlan {
	return daemon.CyclePlan{
		RunID:        runID,
		GitRoot:      preflightResult.Git.Root,
		ArtifactDir:  req.artifactDir,
		SourceName:   req.source,
		AgentName:    req.agent,
		Runtime:      resolvePlan.runtime,
		Verification: loaded.Config.Defaults.Verification,
		AutoCommit:   loaded.Config.Defaults.AutoCommit,
		PullRequest: daemon.PullRequestRef{
			Number:         preflightResult.PullRequest.Number,
			BaseRepository: preflightResult.PullRequest.BaseRepository,
			HeadRepository: preflightResult.PullRequest.HeadRepository,
			HeadBranch:     preflightResult.PullRequest.HeadBranch,
		},
		Batches:     resolvePlan.plan.Batches,
		Duplicates:  resolvePlan.plan.Duplicates,
		TotalIssues: len(resolvePlan.selection.Issues),
	}
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
	collaborators := newEngineCollaborators()
	if err := collaborators.runner.Probe(ctx, runtime); err != nil {
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

	// One cockpit for the entire Watch Run, across all Rounds and Batches.
	ui, err := startRunUI(ctx, buildLiveRunView(req, loaded, preflightResult, run.ID, "WaitingForReview", nil, nil), run.ID, loaded.HomeDir, runStore, stderr)
	if err != nil {
		markRunFailed(ctx, runStore, run.ID)
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}
	defer ui.Close()

	result, err := watch.Run(ctx, watch.Request{
		RunID:          run.ID,
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
				fmt.Fprintf(ui.progress, "Review Source status: %s\n", status.State)
			}
			return status, err
		}),
		Fetcher: watch.FetchFunc(func(ctx context.Context, _ int) (watch.FetchResult, error) {
			return fetchWatchRound(ctx, req, loaded, preflightResult, ui.progress)
		}),
		Resolver: watch.ResolveFunc(func(ctx context.Context) (watch.ResolveResult, error) {
			return resolveWatchBatches(ctx, req, loaded, preflightResult, runtime, run.ID, collaborators, runStore, ui)
		}),
		Clock:   watchClock,
		Sleeper: watchSleeper,
		Sink:    store.JournalSink{Store: runStore},
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
		ui.Close()
		printRunFailure(req.name, completeErr, stderr)
		return exitRunFailed
	}
	publishRunOutcome(completeCtx, runStore, completed.ID, completed.State, result.Remaining, stderr)
	ui.Close()

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
	if roundResult.Reused {
		fmt.Fprintf(stderr, "Reused Round %03d with %d Review Issue(s).\n", roundResult.Round, len(roundResult.IssuePaths))
	} else {
		fmt.Fprintf(stderr, "Fetched Round %03d with %d Review Issue(s).\n", roundResult.Round, len(roundResult.IssuePaths))
	}
	return watch.FetchResult{Round: roundResult.Round, Issues: len(roundResult.IssuePaths)}, nil
}

func resolveWatchBatches(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, runtime agent.RuntimeSpec, runID string, collaborators engineCollaborators, runStore *store.Store, ui *runUI) (watch.ResolveResult, error) {
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
		result, err := executeResolveCycle(ctx, req, loaded, preflightResult, runID, resolvePlan, collaborators, runStore, ui)
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
	publishRunOutcome(ctx, runStore, completed.ID, completed.State, 0, stderr)
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

type fetchSuccessView struct {
	RunID          string
	Round          int
	ReviewIssues   int
	RunDatabase    string
	ArtifactDir    string
	ReusedRound    bool
	StartedAgent   bool
	CreatedCommit  bool
	CompletedPush  bool
	ResolvedSource bool
}

func printFetchSuccess(stdout io.Writer, view fetchSuccessView) {
	style := styleFor(stdout)
	fmt.Fprintf(stdout, "%s\n\n", style.green(style.bold("Fetch complete")))
	fmt.Fprintf(stdout, "%s\n", style.cyan("Result:"))
	fmt.Fprintf(stdout, "  Run: %s reached Fetched\n", view.RunID)
	fmt.Fprintf(stdout, "  Round: %03d\n", view.Round)
	if view.ReviewIssues == 0 {
		fmt.Fprintln(stdout, "  Review Issues: none")
	} else {
		fmt.Fprintf(stdout, "  Review Issues: %d\n", view.ReviewIssues)
	}
	if view.ReusedRound {
		fmt.Fprintln(stdout, "  Artifacts: reused existing matching Round")
	} else {
		fmt.Fprintln(stdout, "  Artifacts: created new Round")
	}
	fmt.Fprintln(stdout)
	fmt.Fprintf(stdout, "%s\n", style.cyan("Files:"))
	fmt.Fprintf(stdout, "  Run Database: %s\n", view.RunDatabase)
	fmt.Fprintf(stdout, "  Artifact Directory: %s\n", view.ArtifactDir)
	fmt.Fprintln(stdout)
	fmt.Fprintf(stdout, "%s\n", style.cyan("No side effects:"))
	if !view.StartedAgent && !view.CreatedCommit && !view.CompletedPush && !view.ResolvedSource {
		fmt.Fprintln(stdout, "  Roundfix did not start an Agent, commit, or push.")
		fmt.Fprintln(stdout, "  Roundfix did not resolve Review Source threads.")
		return
	}
	fmt.Fprintf(stdout, "  Agent started: %s\n", yesNo(view.StartedAgent))
	fmt.Fprintf(stdout, "  Commit created: %s\n", yesNo(view.CreatedCommit))
	fmt.Fprintf(stdout, "  Push completed: %s\n", yesNo(view.CompletedPush))
	fmt.Fprintf(stdout, "  Review Source resolved: %s\n", yesNo(view.ResolvedSource))
}

func printLiveRunView(stderr io.Writer, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, runID string, pipelineState string, issues []rounds.Issue, console []string) {
	fmt.Fprint(stderr, roundtui.RenderLiveRunView(buildLiveRunView(req, loaded, preflightResult, runID, pipelineState, issues, console)))
}

func buildLiveRunView(req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, runID string, pipelineState string, issues []rounds.Issue, console []string) roundtui.LiveRunView {
	return roundtui.LiveRunView{
		Command:       req.name,
		Repository:    preflightResult.PullRequest.HeadRepository,
		PRNumber:      preflightResult.PullRequest.Number,
		HeadBranch:    preflightResult.PullRequest.HeadBranch,
		ReviewSource:  displayReviewSource(req.source),
		Agent:         displayAgent(req.agent),
		Model:         req.model,
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
		Width:         liveViewWidth(),
	}
}

func liveTUIEnabled(output io.Writer) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("ROUNDFIX_TUI"))) {
	case "always", "1", "true", "yes", "on":
		return true
	case "never", "0", "false", "no", "off":
		return false
	}
	if strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}
	file, ok := output.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func liveViewWidth() int {
	width, err := strconv.Atoi(strings.TrimSpace(os.Getenv("COLUMNS")))
	if err == nil && width >= 80 {
		return width
	}
	return 100
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
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

type engineCollaborators struct {
	runner    agent.Runner
	verifier  daemon.Verifier
	committer daemon.Committer
	pusher    daemon.Pusher
	source    daemon.ReviewSourceResolver
	worktree  daemon.WorktreeSnapshotter
}

func defaultEngineCollaborators() engineCollaborators {
	return engineCollaborators{
		runner:    agent.DefaultRunner{},
		verifier:  daemon.ExecVerifier{},
		committer: daemon.GitCommitter{},
		pusher:    daemon.GitPusher{},
		source:    daemon.ReviewSourceResolverFunc(defaultResolveReviewSourceIssues),
		worktree:  daemon.GitWorktreeSnapshotter{},
	}
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

func maybeRunFinalPush(ctx context.Context, engine *daemon.Engine, sink runevent.Sink, runID string, loaded roundconfig.Loaded, preflightResult preflight.Result, batchCommitCreated bool, stderr io.Writer) error {
	if !preflightResult.PushPlan.Enabled {
		fmt.Fprintln(stderr, "Final Push skipped: auto-push disabled or no push target configured.")
		publishPushDecision(ctx, sink, runID, "skipped", "Final Push skipped: auto-push disabled or no push target configured.", 0)
		return nil
	}
	if preflightResult.PushPlan.Force {
		return errors.New("Final Push rejected: force-push is not allowed in the MVP")
	}
	if !loaded.Config.Defaults.AutoCommit {
		fmt.Fprintln(stderr, "Final Push skipped: auto-commit disabled.")
		publishPushDecision(ctx, sink, runID, "skipped", "Final Push skipped: auto-commit disabled.", 0)
		return nil
	}
	if preflightResult.Git.UnpushedCommits == 0 && !batchCommitCreated {
		fmt.Fprintln(stderr, "Final Push skipped: no local commits are waiting for the PR Head Branch.")
		publishPushDecision(ctx, sink, runID, "skipped", "Final Push skipped: no local commits are waiting for the PR Head Branch.", 0)
		return nil
	}
	if err := engine.FinalPush(ctx, daemon.FinalPushRequest{
		RunID:   runID,
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
	completeCtx := context.WithoutCancel(ctx)
	if _, err := runStore.CompleteRun(completeCtx, runID, store.StateFailed); err == nil {
		publishRunOutcome(completeCtx, runStore, runID, store.StateFailed, 0, io.Discard)
	}
}

// publishRunOutcome appends the terminal outcome event after CompleteRun so
// terminal states are provable from the journal alone.
func publishRunOutcome(ctx context.Context, runStore *store.Store, runID string, state string, remaining int, stderr io.Writer) {
	payload, err := json.Marshal(map[string]any{"state": state, "remaining": remaining})
	if err != nil {
		return
	}
	if err := (store.JournalSink{Store: runStore}).Publish(context.WithoutCancel(ctx), runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonOutcome,
		Summary: fmt.Sprintf("Run reached %s.", state),
		Time:    time.Now().UTC(),
		Payload: payload,
	}); err != nil {
		fmt.Fprintf(stderr, "Warning: terminal outcome event not journaled: %v\n", err)
	}
}

func commandUsage(name string) string {
	switch name {
	case "attach":
		return `Usage:
  roundfix attach <run-id>

Replays the Run's event timeline from the Run Database, read-only.
Attach never creates Runs, fetches, starts Agents, commits, pushes, or
resolves Review Source threads.

Options:
  --run-id  Run ID to attach to (same as the positional argument)
`
	case "init":
		return `Usage:
  roundfix init [--scope <project|user>] [--force]

Options:
  --scope  Config scope. Supported: project, user
           project writes <repo>/.roundfixrc.yml
           user writes ~/.roundfix/config.yml
  --force  Overwrite an existing config file
`
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
	case "stop":
		return `Usage:
  roundfix stop <run-id>
  roundfix stop --run-id <id>
  roundfix stop --pr <number>
  roundfix stop --head-repo <owner/name> --head-branch <branch>

Options:
  --run-id      Active Run ID to stop
  --run         Alias for --run-id
  --pr          Open Pull Request number used to find the Active Run
  --head-repo   Explicit Head Repository, owner/name
  --head-branch Explicit PR Head Branch
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

func printInitSuccess(result roundconfig.InitResult, stdout io.Writer) {
	style := styleFor(stdout)
	action := "created"
	if result.Overwritten {
		action = "updated"
	}
	fmt.Fprintf(stdout, "%s\n\n", style.green(style.bold("Roundfix config "+action)))
	fmt.Fprintf(stdout, "%s\n", style.cyan("Scope:"))
	fmt.Fprintf(stdout, "  %s\n\n", result.Scope)
	fmt.Fprintf(stdout, "%s\n", style.cyan("Path:"))
	fmt.Fprintf(stdout, "  %s\n\n", result.Path)
	fmt.Fprintf(stdout, "%s\n", style.cyan("Next:"))
	fmt.Fprintln(stdout, "  roundfix fetch --pr <number>")
}

func printInitFailure(err error, stderr io.Writer) {
	style := styleFor(stderr)
	fmt.Fprintf(stderr, "%s\n\n", style.red(style.bold("Init failed")))
	fmt.Fprintf(stderr, "%s\n", style.cyan("Reason:"))
	fmt.Fprintf(stderr, "  %v\n\n", err)
	fmt.Fprintf(stderr, "%s\n", style.cyan("Usage:"))
	fmt.Fprintf(stderr, "  Run '%s init --help' for usage.\n", app.Name)
}

func printStopSuccess(run store.Run, stdout io.Writer) {
	style := styleFor(stdout)
	fmt.Fprintf(stdout, "%s\n\n", style.green(style.bold("Roundfix Run stopped")))
	fmt.Fprintf(stdout, "%s\n", style.cyan("Run:"))
	fmt.Fprintf(stdout, "  ID: %s\n", run.ID)
	fmt.Fprintf(stdout, "  State: %s\n", run.State)
	fmt.Fprintf(stdout, "  Kind: %s\n", run.Kind)
	fmt.Fprintf(stdout, "  PR: #%s %s\n", run.PRNumber, run.HeadRepository)
	fmt.Fprintf(stdout, "  Branch: %s\n\n", run.HeadBranch)
	fmt.Fprintf(stdout, "%s\n", style.cyan("No repository side effects:"))
	fmt.Fprintln(stdout, "  Roundfix released the Active Run lock without editing files, committing, pushing, fetching, or resolving Review Source threads.")
}

func printStopFailure(err error, stderr io.Writer) {
	style := styleFor(stderr)
	fmt.Fprintf(stderr, "%s\n\n", style.red(style.bold("Stop failed")))
	fmt.Fprintf(stderr, "%s\n", style.cyan("Reason:"))
	fmt.Fprintf(stderr, "  %v\n\n", err)
	fmt.Fprintf(stderr, "%s\n", style.cyan("Usage:"))
	fmt.Fprintf(stderr, "  Run '%s stop --help' for usage.\n", app.Name)
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
