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
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	roundnotify "roundfix/internal/notify"
	"roundfix/internal/preflight"
	"roundfix/internal/reviewsource"
	"roundfix/internal/reviewsource/coderabbit"
	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
	"roundfix/internal/watch"
	runworktree "roundfix/internal/worktree"
	roundskills "roundfix/skills"
)

const usage = `Roundfix

Usage:
  roundfix --help
  roundfix --version
  roundfix fetch --source coderabbit --pr <number> [--spec <slug>]
  roundfix resolve --pr <number> [--spec <slug>]
  roundfix watch --source coderabbit --pr <number> [--spec <slug>] --until-clean
  roundfix implement --spec <slug>
  roundfix settle --spec <slug> --task <task_id>
  roundfix reconcile [run-id] [--apply] [--format <text|json>]
  roundfix release plan [--from <tag>] [--to <revision>] [--format <text|json>]
  roundfix release plan --reset-to <version> [--format <text|json>]
  roundfix spec check [<slug> ...] [--format <text|json>] [--strict]
  roundfix spec audit <slug> [--format <text|json>]
  roundfix baseline plan (--profile <id> | --profile-file <draft.json>) [--decision <id=value> ...] [--decision-file <path> ...] [--repo <path>] [--format <text|json>]
  roundfix baseline apply --plan <file> --confirm-plan <digest> [--repo <path>] [--format <text|json>]
  roundfix baseline capabilities check [--profile <id>] [--repo <path>] [--format <text|json>]
  roundfix baseline profile init --id <id> [--from <built-in-id>]
  roundfix baseline profile show <id> [--format <text|json>]
  roundfix baseline profile validate [<id>|<path>] [--format <text|json>]
  roundfix baseline skills restore --profile <id> [--skill <name> ...] [--source-dir <path>] [--confirm-plan <digest>] [--repo <path>] [--format <text|json>]
  roundfix baseline assets sync --source-dir <path> [--check] [--format <text|json>]
  roundfix profiles show [--category <category>] [--json]
  roundfix profiles configure --scope user|project [--file <path>] [--remove <category>] [--dry-run] [--yes] [--json]
  roundfix profiles validate [--category <category>] [--json]
  roundfix archive <slug>
  roundfix init [--scope <project|user>]
  roundfix setup [--yes] [--no-input]
  roundfix doctor
  roundfix gc [--dry-run]
  roundfix storage report
  roundfix upgrade [--check]
  roundfix runs
  roundfix runs list [--all] [--state <active|terminal|all>] [--limit N]
  roundfix stop [<run-id>|--run-id <id>|--pr <number>|--spec <slug>]
  roundfix attach [<run-id>] [--no-input]
  roundfix events <run-id> [--follow] [--filter <categories>]
  roundfix skills check
  roundfix skills install [--target <project|codex|claude|opencode|all>]

Commands:
  init       Create User Config or Project Config
  fetch      Download review issues for an Open Pull Request
  resolve    Resolve downloaded Unresolved Review Issues
  watch      Fetch and resolve in a watched loop
  implement  Execute a Spec's Task Graph as one Run
  settle     Verify and commit all current worktree changes for one failed Task
  reconcile  Inspect or release proven terminal spec Run worktrees
  release    Plan the next release version without mutating repository or release state
  spec       Check Spec artifact consistency; audit Spec delivery
  baseline   Plan, apply, and validate a Context-Driven Baseline
  profiles   Show Agent Selection Profiles and advisory recommendations
  archive    Archive a completed Spec
  stop       Request or force-stop an Active Run
  setup      Verify and prepare this machine for Roundfix Runs
  doctor     Diagnose this machine's readiness for Roundfix Runs
  gc         Prune old terminal Run journals and run artifacts
  storage    Report measured Run Database and Artifact Root storage
  upgrade    Upgrade the Roundfix binary from GitHub Releases
  runs       List Runs from the Run Database
  attach     Replay a Run's event timeline from the Run Database
  events     Replay or follow a Run's Supervisor event stream as JSONL
  skills     List, check, or install the bundled Roundfix skills

Options:
  -h, --help      Show help
  -v, --version   Show version
`

const (
	exitOK         = 0
	exitRunFailed  = 1
	exitPreflight  = 2
	exitUnverified = 3
	exitSIGINT     = 130
)

type commandRequest struct {
	name                 string
	arguments            []string
	pr                   string
	spec                 string
	source               string
	agent                string
	agentSet             bool
	round                string
	noInput              bool
	interactive          bool
	inputShown           bool
	untilClean           bool
	maxRounds            int
	artifactDir          string
	explicitArtifactDir  bool
	reviewRoot           string
	baseRepo             string
	model                string
	modelSet             bool
	reasoningEffort      string
	reasoningEffortSet   bool
	agentCmd             string
	agentFullAccess      bool
	noAgentConsole       bool
	headBranch           string
	headRepo             string
	skipBranchIntegrity  bool
	branchIntegrity      branchIntegrityReport
	detach               bool
	detachChild          *detachChild
	branchIntegrityActor string
}

var runCommandPreflight = defaultRunCommandPreflight
var fetchReviewItems = defaultFetchReviewItems
var newOutcomeNotifier = roundnotify.New

// newEngineCollaborators is the single test seam for Run engine
// collaborators; orchestration itself lives in the daemon Run engine and
// receives them through an explicit dependencies struct.
var newEngineCollaborators = defaultEngineCollaborators
var watchReviewEvidence = defaultWatchReviewEvidence
var watchHeadSHA = defaultWatchHeadSHA
var newReviewRequester = defaultNewReviewRequester
var listPendingRunWork = runworktree.ListPendingRunWork
var supersedingQAReport = runworktree.SupersedingQAReport
var classifyRunBranchSet = runworktree.ClassifyRunBranchSet
var integratePendingRunWork = runworktree.IntegratePendingRunWork
var refreshBranchIntegrityHead = defaultRefreshBranchIntegrityHead
var commentOnPullRequest = defaultCommentOnPullRequest
var reviewSpecGitRunner preflight.GitRunner = preflight.ExecGitRunner{}
var watchClock watch.Clock
var watchSleeper watch.Sleeper
var inspectChangedPaths = defaultInspectChangedPaths
var collectInteractiveInput = defaultCollectInteractiveInput
var suggestCurrentPullRequest = defaultSuggestCurrentPullRequest
var resolvePullRequestForStop = defaultResolvePullRequestForStop
var promptInitScope = defaultPromptInitScope
var resolveSkillsProjectRoot = defaultResolveSkillsProjectRoot
var promptProjectClaudeSkillSymlink = defaultPromptProjectClaudeSkillSymlink
var cancelStopAgentSession = defaultCancelStopAgentSession
var listRoundfixAgentSessions = defaultListRoundfixAgentSessions
var closeStopAgentSession = defaultCloseStopAgentSession
var ownerProcesses OwnerProcessController = store.NewOwnerProcessController()
var reconcileProcesses reconcileProcessController = store.NewOwnerProcessController()

type validationError struct {
	message string
}

func (err validationError) Error() string {
	return err.message
}

type commandEnvironment struct {
	homeDir        string
	homeDirErr     error
	workDir        string
	workDirErr     error
	environ        []string
	detachFD       string
	detachTempPath string
	tuiMode        string
	term           string
	columns        string
	colorMode      string
	noColor        string
	codexPath      string
	executablePath string
	branchActor    string
	dependencies   commandDependencies
}

type commandDependencies struct {
	runCommandPreflight             func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error)
	fetchReviewItems                func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error)
	newOutcomeNotifier              func(roundconfig.Config) roundnotify.Notifier
	newEngineCollaborators          func() engineCollaborators
	watchReviewEvidence             func(context.Context, reviewsource.EvidenceRequest) (reviewsource.Evidence, error)
	watchHeadSHA                    func(context.Context, string) (string, error)
	checkoutReader                  watch.CheckoutReader
	newReviewRequester              func(runevent.Sink) reviewsource.ReviewRequester
	listPendingRunWork              func(context.Context, string, string) ([]runworktree.PendingRunWork, error)
	supersedingQAReport             func(context.Context, string, string, string, string) (string, bool)
	classifyRunBranchSet            func(context.Context, string, string, string, []store.Run) (runworktree.BranchSetClassification, error)
	integratePendingRunWork         func(context.Context, string, string, string) error
	refreshBranchIntegrityHead      func(context.Context, preflight.Result) (preflight.Result, error)
	commentOnPullRequest            func(context.Context, string, string, int, string) error
	reviewSpecGitRunner             preflight.GitRunner
	watchClock                      watch.Clock
	watchSleeper                    watch.Sleeper
	inspectChangedPaths             func(context.Context, string) ([]preflight.ChangedPath, error)
	collectInteractiveInput         func(context.Context, roundtui.InputRequest) (roundtui.CommandValues, error)
	suggestCurrentPullRequest       func(context.Context, string) (string, error)
	resolvePullRequestForStop       func(context.Context, string, string) (preflight.PullRequest, error)
	promptInitScope                 func(context.Context, io.Writer) (string, error)
	resolveSkillsProjectRoot        func(context.Context, string) (string, error)
	promptProjectClaudeSkillSymlink func(context.Context, io.Writer, string, string) (bool, error)
	cancelStopAgentSession          func(context.Context, agent.RuntimeSpec, agent.SessionRef) error
	listRoundfixAgentSessions       func(context.Context, agent.RuntimeSpec, string) ([]agent.RoundfixSession, error)
	closeStopAgentSession           func(context.Context, agent.RuntimeSpec, agent.SessionRef) error
	ownerProcesses                  OwnerProcessController
	reconcileProcesses              reconcileProcessController
	fallbackConfirmationInput       func() io.Reader
	fallbackConfirmationAvailable   func(io.Writer) bool
	runsListNow                     func() time.Time
	runsInteractiveInputAvailable   func() bool
	doctor                          doctorDependencies
	runBrowserSession               func(context.Context, io.Writer, []store.Run, []store.Run) (roundtui.BrowserOutcome, error)
	browserAttachCockpit            func(context.Context, roundconfig.Loaded, *store.Store, store.Run, attachCapacities, io.Writer, io.Writer) int
	setup                           setupDependencies
	profilesConfigureInput          func() io.Reader
	confirmProfilesConfigure        func(context.Context, io.Writer, string) (bool, error)
	implementOwnerIdentity          func(context.Context) string
	createRunWorktree               func(context.Context, runworktree.CreateOptions) (runworktree.Ref, error)
	integrateRunWorktree            func(context.Context, runworktree.Ref, string, string) (runworktree.IntegrationResult, error)
	cleanupCleanRunWorktree         func(context.Context, runworktree.Ref) error
	pruneTerminalRunWorktrees       func(context.Context, string, string, runworktree.TerminalRunReconciliationStore, runworktree.TerminalRunLookup) ([]runworktree.PrunedRef, error)
	loadCommittedSpecGraph          func(context.Context, string, roundconfig.SpecsRoot, string, string) (*spec.Graph, string, error)
	detachTimeouts                  detachPhaseTimeouts
	attachInteractiveInputAvailable func() bool
	attachSleep                     func(context.Context) error
	upgrade                         upgradeDependencies
	versionFreshness                versionFreshnessDependencies
	gc                              gcDependencies
	releasePlanGitRunner            preflight.GitRunner
	releasePlanGHRunner             preflight.GHRunner
}

type commandDependenciesContextKey struct{}

func defaultCommandDependencies() commandDependencies {
	return commandDependencies{
		runCommandPreflight:             runCommandPreflight,
		fetchReviewItems:                fetchReviewItems,
		newOutcomeNotifier:              newOutcomeNotifier,
		newEngineCollaborators:          newEngineCollaborators,
		watchReviewEvidence:             watchReviewEvidence,
		watchHeadSHA:                    watchHeadSHA,
		checkoutReader:                  watch.GitCheckoutReader{},
		newReviewRequester:              newReviewRequester,
		listPendingRunWork:              listPendingRunWork,
		supersedingQAReport:             supersedingQAReport,
		classifyRunBranchSet:            classifyRunBranchSet,
		integratePendingRunWork:         integratePendingRunWork,
		refreshBranchIntegrityHead:      refreshBranchIntegrityHead,
		commentOnPullRequest:            commentOnPullRequest,
		reviewSpecGitRunner:             reviewSpecGitRunner,
		watchClock:                      watchClock,
		watchSleeper:                    watchSleeper,
		inspectChangedPaths:             inspectChangedPaths,
		collectInteractiveInput:         collectInteractiveInput,
		suggestCurrentPullRequest:       suggestCurrentPullRequest,
		resolvePullRequestForStop:       resolvePullRequestForStop,
		promptInitScope:                 promptInitScope,
		resolveSkillsProjectRoot:        resolveSkillsProjectRoot,
		promptProjectClaudeSkillSymlink: promptProjectClaudeSkillSymlink,
		cancelStopAgentSession:          cancelStopAgentSession,
		listRoundfixAgentSessions:       listRoundfixAgentSessions,
		closeStopAgentSession:           closeStopAgentSession,
		ownerProcesses:                  ownerProcesses,
		reconcileProcesses:              reconcileProcesses,
		fallbackConfirmationInput:       fallbackConfirmationInput,
		fallbackConfirmationAvailable:   fallbackConfirmationAvailable,
		runsListNow:                     runsListNow,
		runsInteractiveInputAvailable:   runsInteractiveInputAvailable,
		doctor:                          defaultDoctorDependencies(),
		runBrowserSession:               runBrowserSession,
		browserAttachCockpit:            browserAttachCockpit,
		setup:                           defaultSetupDependencies(),
		profilesConfigureInput:          func() io.Reader { return os.Stdin },
		confirmProfilesConfigure:        defaultConfirmProfilesConfigure,
		implementOwnerIdentity:          implementOwnerIdentity,
		createRunWorktree:               createRunWorktree,
		integrateRunWorktree:            integrateRunWorktree,
		cleanupCleanRunWorktree:         cleanupCleanRunWorktree,
		pruneTerminalRunWorktrees:       pruneTerminalRunWorktrees,
		loadCommittedSpecGraph:          loadCommittedSpecGraph,
		detachTimeouts:                  detachTimeouts,
		attachInteractiveInputAvailable: attachInteractiveInputAvailable,
		attachSleep:                     attachSleep,
		upgrade:                         upgradeDeps,
		versionFreshness:                versionFreshnessDeps,
		gc:                              gcDeps,
		releasePlanGitRunner:            releasePlanCommandGitRunner,
		releasePlanGHRunner:             releasePlanCommandGHRunner,
	}
}

func commandDependenciesForContext(ctx context.Context) commandDependencies {
	if ctx != nil {
		if dependencies, ok := ctx.Value(commandDependenciesContextKey{}).(commandDependencies); ok {
			return dependencies
		}
	}
	return defaultCommandDependencies()
}

func contextWithCommandDependencies(ctx context.Context, dependencies commandDependencies) context.Context {
	return context.WithValue(ctx, commandDependenciesContextKey{}, dependencies)
}

func commandEnvironmentFromProcess() commandEnvironment {
	homeDir, homeDirErr := os.UserHomeDir()
	workDir, workDirErr := os.Getwd()
	branchActor := ""
	for _, key := range []string{"GITHUB_ACTOR", "GIT_AUTHOR_NAME", "USER", "USERNAME"} {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			branchActor = value
			break
		}
	}
	return commandEnvironment{
		homeDir:        homeDir,
		homeDirErr:     homeDirErr,
		workDir:        workDir,
		workDirErr:     workDirErr,
		environ:        os.Environ(),
		detachFD:       os.Getenv(detachHandshakeFDEnv),
		detachTempPath: os.Getenv(detachConsoleTempEnv),
		tuiMode:        os.Getenv("ROUNDFIX_TUI"),
		term:           os.Getenv("TERM"),
		columns:        os.Getenv("COLUMNS"),
		colorMode:      os.Getenv("ROUNDFIX_COLOR"),
		noColor:        os.Getenv("NO_COLOR"),
		codexPath:      os.Getenv("CODEX_PATH"),
		executablePath: os.Getenv("PATH"),
		branchActor:    branchActor,
		dependencies:   defaultCommandDependencies(),
	}
}

type commandEnvironmentWriter struct {
	io.Writer
	environment commandEnvironment
}

func commandWriter(writer io.Writer, environment commandEnvironment) io.Writer {
	return commandEnvironmentWriter{Writer: writer, environment: environment}
}

func environmentForWriter(writer io.Writer) (commandEnvironment, io.Writer) {
	if wrapped, ok := writer.(commandEnvironmentWriter); ok {
		return wrapped.environment, wrapped.Writer
	}
	return commandEnvironment{}, writer
}

func (environment commandEnvironment) loadOptions(stderr io.Writer) (roundconfig.LoadOptions, error) {
	if environment.homeDirErr != nil {
		return roundconfig.LoadOptions{}, fmt.Errorf("resolve User Config home: %w", environment.homeDirErr)
	}
	if environment.workDirErr != nil {
		return roundconfig.LoadOptions{}, fmt.Errorf("resolve work directory: %w", environment.workDirErr)
	}
	return roundconfig.LoadOptions{
		HomeDir: environment.homeDir,
		WorkDir: environment.workDir,
		Stderr:  stderr,
	}, nil
}

func loadCommandConfig(environment commandEnvironment, stderr io.Writer) (roundconfig.Loaded, error) {
	options, err := environment.loadOptions(stderr)
	if err != nil {
		return roundconfig.Loaded{}, err
	}
	return roundconfig.Load(options)
}

func (environment commandEnvironment) resolveWorkDir(operation string) (string, error) {
	if environment.workDirErr != nil {
		return "", fmt.Errorf("%s: %w", operation, environment.workDirErr)
	}
	return environment.workDir, nil
}

func (environment commandEnvironment) executableDirectories(operation string) ([]string, error) {
	workDir, err := environment.resolveWorkDir(operation)
	if err != nil {
		return nil, err
	}
	directories := filepath.SplitList(environment.executablePath)
	for index, directory := range directories {
		if directory == "" {
			directory = workDir
		} else if !filepath.IsAbs(directory) {
			directory = filepath.Join(workDir, directory)
		}
		directories[index] = filepath.Clean(directory)
	}
	return directories, nil
}

func Run(args []string, stdout, stderr io.Writer) int {
	ctx, cleanup, interrupted := interruptContext(context.Background())
	defer cleanup()
	code := runWithContext(ctx, args, stdout, stderr, commandEnvironmentFromProcess())
	return exitForInterrupt(code, interrupted())
}

func RunContext(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	return runWithContext(ctx, args, stdout, stderr, commandEnvironmentFromProcess())
}

func runWithContext(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = contextWithCommandDependencies(ctx, environment.dependencies)
	stdout = commandWriter(stdout, environment)
	stderr = commandWriter(stderr, environment)
	detachChild, err := newDetachChildFromEnv(environment.detachFD, environment.detachTempPath)
	if err != nil {
		fmt.Fprintf(stderr, "%s: detach setup failed: %v\n", app.Name, err)
		return exitPreflight
	}
	if detachChild != nil {
		defer func() {
			_ = detachChild.Close()
		}()
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
		fmt.Fprintf(stdout, "%s %s\n", app.Name, app.VersionLine())
		return exitOK
	case "init":
		return runInitCommand(ctx, args[1:], stdout, stderr, environment)
	case "setup":
		return runSetupCommand(ctx, args[1:], stdout, stderr, environment)
	case "doctor":
		return runDoctorCommand(ctx, args[1:], stdout, stderr, environment)
	case "gc":
		return runGCCommand(ctx, args[1:], stdout, stderr, environment)
	case "storage":
		return runStorageCommand(ctx, args[1:], stdout, stderr, environment)
	case "upgrade":
		return runUpgradeCommand(ctx, args[1:], stdout, stderr)
	case "runs":
		return runRunsCommand(ctx, args[1:], stdout, stderr, environment)
	case "stop":
		return runStopCommand(ctx, args[1:], stdout, stderr, environment)
	case "attach":
		return runAttachCommand(ctx, args[1:], stdout, stderr, environment)
	case "events":
		return runEventsCommand(ctx, args[1:], stdout, stderr, environment)
	case "skills":
		return runSkillsCommand(ctx, args[1:], stdout, stderr, environment)
	case "fetch", "resolve", "watch":
		return runOperationalCommand(ctx, args[0], args[1:], stdout, stderr, detachChild, environment)
	case "implement":
		return runImplementCommand(ctx, args[1:], stdout, stderr, detachChild, environment)
	case "settle":
		return runSettleCommand(ctx, args[1:], stdout, stderr, environment)
	case "reconcile":
		return runReconcileCommand(ctx, args[1:], stdout, stderr, environment)
	case "release":
		return runReleaseCommand(ctx, args[1:], stdout, stderr, environment)
	case "spec":
		return runSpecCommand(ctx, args[1:], stdout, stderr, environment)
	case "baseline":
		return runBaselineCommand(ctx, args[1:], stdout, stderr, environment)
	case "profiles":
		return runProfilesCommand(ctx, args[1:], stdout, stderr, environment)
	case "archive":
		return runArchiveCommand(ctx, args[1:], stdout, stderr, environment)
	default:
		fmt.Fprintf(stderr, "%s: unknown command %q\n", app.Name, args[0])
		fmt.Fprintf(stderr, "Run '%s --help' for usage.\n", app.Name)
		return exitPreflight
	}
}

func runInitCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
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
		selectedScope, err = commandDependenciesForContext(ctx).promptInitScope(ctx, stderr)
		if err != nil {
			printInitFailure(err, stderr)
			return exitPreflight
		}
	}

	loadOptions, err := environment.loadOptions(stderr)
	if err != nil {
		printInitFailure(err, stderr)
		return exitRunFailed
	}
	result, err := roundconfig.Init(ctx, roundconfig.InitOptions{
		Scope:   selectedScope,
		HomeDir: loadOptions.HomeDir,
		WorkDir: loadOptions.WorkDir,
		Force:   *force,
	})
	if err != nil {
		printInitFailure(err, stderr)
		return exitPreflight
	}
	printInitSuccess(result, stdout)
	return exitOK
}

func runSkillsCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
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
			fmt.Fprintf(stderr, "%s: Roundfix skill check failed:\n", app.Name)
			for _, diagnostic := range diagnostics {
				fmt.Fprintf(stderr, "  %s: %s\n", diagnostic.Path, diagnostic.Message)
			}
			return exitRunFailed
		}
		fmt.Fprintf(stdout, "Roundfix skill check passed: %s\n", strings.Join(roundskills.Names(), ", "))
		return exitOK
	case "list":
		if commandWantsHelp(args[1:]) {
			fmt.Fprint(stdout, commandUsage("skills"))
			return exitOK
		}
		if len(args) > 1 {
			fmt.Fprintf(stderr, "%s: unexpected argument %q\n", app.Name, args[1])
			fmt.Fprintf(stderr, "Run '%s skills --help' for usage.\n", app.Name)
			return exitPreflight
		}
		return runSkillsList(stdout)
	case "install":
		return runSkillsInstall(ctx, args[1:], stdout, stderr, environment)
	default:
		fmt.Fprintf(stderr, "%s: unknown skills command %q\n", app.Name, args[0])
		fmt.Fprintf(stderr, "Run '%s skills --help' for usage.\n", app.Name)
		return exitPreflight
	}
}

func runSkillsList(stdout io.Writer) int {
	fmt.Fprintln(stdout, "Bundled skills (install with 'roundfix skills install'):")
	for _, name := range roundskills.Names() {
		fmt.Fprintf(stdout, "  %s\n", name)
	}
	if recommended := roundskills.Recommended(); len(recommended) > 0 {
		fmt.Fprintln(stdout, "Recommended skills (managed by your skills tooling, not shipped):")
		for _, name := range recommended {
			fmt.Fprintf(stdout, "  %s\n", name)
		}
	}
	return exitOK
}

type stopRequest struct {
	runID                   string
	pr                      string
	spec                    string
	headRepo                string
	headBranch              string
	force                   bool
	ownerIdentityUnreadable bool
}

type stopResult struct {
	Run          store.Run
	Requested    bool
	Forced       bool
	Transitioned bool
	Warnings     []cleanupWarning
}

// cleanupWarningKind is declared by each secondary cleanup warning producer
// so consumers branch on the producer's classification instead of matching
// warning text.
type cleanupWarningKind int

const (
	// cleanupWarningNotice reports non-failure cleanup activity.
	cleanupWarningNotice cleanupWarningKind = iota
	// cleanupWarningFailure reports a failed or skipped cleanup step.
	cleanupWarningFailure
)

type cleanupWarning struct {
	kind cleanupWarningKind
	text string
}

func cleanupNoticef(format string, args ...any) cleanupWarning {
	return cleanupWarning{kind: cleanupWarningNotice, text: fmt.Sprintf(format, args...)}
}

func cleanupFailuref(format string, args ...any) cleanupWarning {
	return cleanupWarning{kind: cleanupWarningFailure, text: fmt.Sprintf(format, args...)}
}

type OwnerProcessController interface {
	// ProveOwner proves the recorded owner identity without side effects so
	// Force Stop can fail closed before touching anything the Run owns.
	ProveOwner(ctx context.Context, pid int, recordedIdentity string) error
	TerminateAndWait(ctx context.Context, pid int, recordedIdentity string) error
}

type forceStopOwnerError struct {
	RunID string
	PID   int
	Step  string
	Err   error
}

func (err forceStopOwnerError) Error() string {
	return fmt.Sprintf(
		"force stop Run %s owner PID %d failed step %q: %v; Run remains Active; Active Run lock retained",
		err.RunID,
		err.PID,
		err.Step,
		err.Err,
	)
}

func (err forceStopOwnerError) Unwrap() error {
	return err.Err
}

func runStopCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("stop"))
		return exitOK
	}

	req, err := parseStopCommand(args)
	if err != nil {
		printStopFailure(err, stderr)
		return exitPreflight
	}
	loaded, err := loadCommandConfig(environment, stderr)
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

	result, err := stopTargetRun(ctx, req, loaded, runStore, stderr)
	if err != nil {
		printStopFailure(err, stderr)
		journalStopPrimaryFailure(ctx, runStore, result.Run.ID, err)
		reportSecondaryCleanupWarnings(ctx, runStore, result.Run.ID, result.Warnings, stderr)
		var ownerErr forceStopOwnerError
		if errors.As(err, &ownerErr) {
			return exitRunFailed
		}
		return exitPreflight
	}
	if result.Transitioned {
		publishTerminalCompletion(
			ctx,
			runStore,
			outcomeNotifierFromConfig(ctx, loaded.Config),
			stderr,
			store.CompleteRunResult{Run: result.Run, Transitioned: true},
			0,
		)
	}
	reportSecondaryCleanupWarnings(ctx, runStore, result.Run.ID, result.Warnings, stderr)
	printStopSuccess(result, stdout)
	return exitOK
}

func parseStopCommand(args []string) (stopRequest, error) {
	req := stopRequest{}
	fs := flag.NewFlagSet("stop", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&req.runID, "run-id", "", "Run ID to stop")
	fs.StringVar(&req.runID, "run", "", "Run ID to stop")
	fs.StringVar(&req.pr, "pr", "", "Open Pull Request number")
	fs.StringVar(&req.spec, "spec", "", "Spec slug under docs/specs/")
	fs.StringVar(&req.headRepo, "head-repo", "", "Head Repository, owner/name")
	fs.StringVar(&req.headBranch, "head-branch", "", "PR Head Branch")
	fs.BoolVar(&req.force, "force", false, "Immediately stop the target Run and release its lock")
	fs.BoolVar(&req.ownerIdentityUnreadable, "owner-identity-unreadable", false, "Permit Force Stop only when the owner identity is unreadable")
	if err := fs.Parse(hoistCommandFlags(args, stopValueFlags)); err != nil {
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
	req.pr = strings.TrimSpace(req.pr)
	req.spec = strings.TrimSpace(req.spec)
	req.headRepo = strings.TrimSpace(req.headRepo)
	req.headBranch = strings.TrimSpace(req.headBranch)
	if req.ownerIdentityUnreadable && !req.force {
		return req, validationError{message: "--owner-identity-unreadable requires --force"}
	}
	headSelector := req.headRepo != "" || req.headBranch != ""
	if req.runID != "" && (req.pr != "" || req.spec != "" || headSelector) {
		return req, validationError{message: "--run-id cannot be combined with --pr, --spec, --head-repo, or --head-branch"}
	}
	if req.spec != "" && (req.pr != "" || headSelector) {
		return req, validationError{message: "--spec cannot be combined with --pr, --head-repo, or --head-branch"}
	}
	if req.pr != "" && headSelector {
		return req, validationError{message: "--pr cannot be combined with --head-repo or --head-branch"}
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

var stopValueFlags = map[string]bool{
	"run-id":      true,
	"run":         true,
	"pr":          true,
	"spec":        true,
	"head-repo":   true,
	"head-branch": true,
}

func stopTargetRun(ctx context.Context, req stopRequest, loaded roundconfig.Loaded, runStore *store.Store, stderr io.Writer) (stopResult, error) {
	if req.runID != "" {
		current, found, err := runStore.Run(ctx, req.runID)
		if err != nil {
			return stopResult{}, err
		} else if !found {
			return stopResult{}, validationError{message: fmt.Sprintf("Run %q does not exist", req.runID)}
		}
		if req.force {
			return forceStopRun(ctx, runStore, current, loaded.Config.Worktree.Location, req.ownerIdentityUnreadable)
		}
		if reclaimed, ok, err := reclaimOrphanedActiveRun(ctx, runStore, current, stderr); err != nil {
			return stopResult{}, err
		} else if ok {
			return stopResult{Run: reclaimed}, nil
		}
		if err := runStore.RequestStop(ctx, current.ID); err != nil {
			return stopResult{}, err
		}
		return stopResult{Run: current, Requested: true}, nil
	}

	specSlug := strings.TrimSpace(req.spec)
	if specSlug != "" {
		gitRoot := strings.TrimSpace(loaded.GitRoot)
		if gitRoot == "" {
			return stopResult{}, validationError{message: "stop --spec requires running inside a git repository"}
		}
		active, found, err := runStore.ActiveSpecRun(ctx, gitRoot, specSlug)
		if err != nil {
			return stopResult{}, err
		}
		if !found {
			return stopResult{}, validationError{message: fmt.Sprintf("no Active Run exists for repository %q and Spec %q", gitRoot, specSlug)}
		}
		if req.force {
			return forceStopRun(ctx, runStore, active, loaded.Config.Worktree.Location, req.ownerIdentityUnreadable)
		}
		if reclaimed, ok, err := reclaimOrphanedActiveRun(ctx, runStore, active, stderr); err != nil {
			return stopResult{}, err
		} else if ok {
			return stopResult{Run: reclaimed}, nil
		}
		if err := runStore.RequestStop(ctx, active.ID); err != nil {
			return stopResult{}, err
		}
		return stopResult{Run: active, Requested: true}, nil
	}

	headRepo := strings.TrimSpace(req.headRepo)
	headBranch := strings.TrimSpace(req.headBranch)
	if headRepo == "" || headBranch == "" {
		pr := strings.TrimSpace(req.pr)
		if pr == "" {
			suggested, _ := commandDependenciesForContext(ctx).suggestCurrentPullRequest(ctx, loaded.GitRoot)
			pr = strings.TrimSpace(suggested)
		}
		if pr == "" {
			return stopResult{}, validationError{message: "missing stop target; pass a Run ID, --run-id, --pr, --spec, or --head-repo with --head-branch"}
		}
		resolved, err := commandDependenciesForContext(ctx).resolvePullRequestForStop(ctx, loaded.GitRoot, pr)
		if err != nil {
			return stopResult{}, fmt.Errorf("resolve Open Pull Request %s for stop target: %w", pr, err)
		}
		headRepo = resolved.HeadRepository
		headBranch = resolved.HeadBranch
	}

	active, found, err := runStore.ActiveRun(ctx, headRepo, headBranch)
	if err != nil {
		return stopResult{}, err
	}
	if !found {
		return stopResult{}, validationError{message: fmt.Sprintf("no Active Run exists for Head Repository %q and PR Head Branch %q", headRepo, headBranch)}
	}
	if req.force {
		return forceStopRun(ctx, runStore, active, loaded.Config.Worktree.Location, req.ownerIdentityUnreadable)
	}
	if reclaimed, ok, err := reclaimOrphanedActiveRun(ctx, runStore, active, stderr); err != nil {
		return stopResult{}, err
	} else if ok {
		return stopResult{Run: reclaimed}, nil
	}
	if err := runStore.RequestStop(ctx, active.ID); err != nil {
		return stopResult{}, err
	}
	return stopResult{Run: active, Requested: true}, nil
}

func forceStopRun(ctx context.Context, runStore *store.Store, active store.Run, worktreeLocation string, ownerIdentityUnreadable bool) (stopResult, error) {
	if store.IsTerminalState(active.State) {
		if ownerIdentityUnreadable {
			return stopResult{Run: active}, validationError{message: "--owner-identity-unreadable requires an Active Run whose owner identity proof is unreadable"}
		}
		if active.State == store.StateStopped {
			return stopResult{Run: active, Forced: true}, nil
		}
		return stopResult{}, store.TerminalOutcomeConflictError{
			RunID:     active.ID,
			Stored:    active.State,
			Requested: store.StateStopped,
		}
	}

	pid, ok := activeOwnerPID(active)
	if !ok {
		return stopResult{Run: active}, forceStopOwnerError{
			RunID: active.ID,
			PID:   0,
			Step:  "validate recorded owner PID",
			Err:   store.ErrOwnerProcessIdentityUnproven,
		}
	}
	// The owner proof is read-only, so it runs before anything the Run owns is
	// touched: a Run left Active with its lock retained must keep its Agent
	// Sessions intact. Once the owner is proven, PRD Core Feature 3 orders the
	// destructive steps: cancel registered Agent Sessions, then terminate the
	// owner and wait for its exit.
	terminationIdentity := active.OwnerIdentity
	if err := commandDependenciesForContext(ctx).ownerProcesses.ProveOwner(ctx, pid, active.OwnerIdentity); err != nil {
		switch {
		case ownerIdentityUnreadable && errors.Is(err, store.ErrOwnerIdentityUnreadable):
			// The explicit supervised flag authorizes the existing PID-only
			// termination proof after the identity read failed. No config,
			// environment, default, or timeout reaches this branch.
			terminationIdentity = ""
		case ownerIdentityUnreadable && errors.Is(err, store.ErrOwnerProcessIdentityUnproven):
			return stopResult{Run: active}, validationError{message: fmt.Sprintf("--owner-identity-unreadable cannot override a proven owner identity mismatch: %v", err)}
		default:
			return stopResult{Run: active}, forceStopOwnerError{
				RunID: active.ID,
				PID:   pid,
				Step:  forceStopOwnerStep(err, "prove owner process identity"),
				Err:   err,
			}
		}
	} else if ownerIdentityUnreadable {
		return stopResult{Run: active}, validationError{message: "--owner-identity-unreadable may be used only when the owner identity is unreadable"}
	}
	warnings := bestEffortForceStopAgentSessions(ctx, runStore, active)
	if err := commandDependenciesForContext(ctx).ownerProcesses.TerminateAndWait(ctx, pid, terminationIdentity); err != nil {
		return stopResult{Run: active, Warnings: warnings}, forceStopOwnerError{
			RunID: active.ID,
			PID:   pid,
			Step:  forceStopOwnerStep(err, "prove owner exit"),
			Err:   err,
		}
	}
	run, err := runStore.CompleteRun(ctx, active.ID, store.StateStopped)
	if err != nil {
		return stopResult{Run: active, Warnings: warnings}, err
	}
	if strings.TrimSpace(active.GitRoot) != "" && strings.TrimSpace(active.WorkDir) != "" {
		pruned, pruneErr := commandDependenciesForContext(ctx).pruneTerminalRunWorktrees(ctx, active.GitRoot, worktreeLocation, runStore, func(lookupCtx context.Context, runID string) (store.Run, bool, error) {
			return runStore.Run(lookupCtx, runID)
		})
		for _, ref := range pruned {
			warnings = append(warnings, cleanupNoticef("reaped terminal Worktree path=%s branch=%s", ref.Path, ref.Branch))
		}
		if pruneErr != nil {
			warnings = append(warnings, cleanupFailuref("terminal Worktree reap failed for Run %s: %v", active.ID, pruneErr))
		}
	}
	return stopResult{
		Run:          run.Run,
		Forced:       true,
		Transitioned: run.Transitioned,
		Warnings:     warnings,
	}, nil
}

// forceStopOwnerStep names the owner-control step that failed, preferring the
// controller's own step label when it reported one.
func forceStopOwnerStep(err error, fallback string) string {
	var controlErr store.OwnerProcessControlError
	if errors.As(err, &controlErr) && strings.TrimSpace(controlErr.Step) != "" {
		return controlErr.Step
	}
	return fallback
}

func bestEffortForceStopAgentSessions(ctx context.Context, runStore *store.Store, run store.Run) []cleanupWarning {
	activeScopes, err := runStore.ActiveAgentSelectionScopes(ctx, run.ID)
	if err != nil {
		return []cleanupWarning{cleanupFailuref("Agent Session cleanup registry failed for Run %s: %v", run.ID, err)}
	}
	warnings := []cleanupWarning{}
	for _, selection := range activeScopes {
		warnings = append(warnings, cleanupRegisteredAgentSession(ctx, runStore, run, selection)...)
	}
	return warnings
}

func defaultCancelStopAgentSession(ctx context.Context, runtime agent.RuntimeSpec, session agent.SessionRef) error {
	return agent.NewDefaultRunner().CancelSession(ctx, runtime, session)
}

func cleanupRegisteredAgentSession(
	ctx context.Context,
	runStore *store.Store,
	run store.Run,
	selection store.AgentSelectionAttempt,
) []cleanupWarning {
	warnings := []cleanupWarning{}
	runtime, err := agent.RuntimeFor(agent.RuntimeOptions{
		Agent:           selection.Runtime,
		Model:           selection.Model,
		ReasoningEffort: selection.ReasoningEffort,
	})
	if err != nil {
		return []cleanupWarning{cleanupFailuref(
			"Agent Session cleanup skipped for %s %s in Run %s: %v",
			selection.ScopeKind,
			selection.ScopeID,
			run.ID,
			err,
		)}
	}
	session, err := registeredAgentSessionRef(run, selection)
	if err != nil {
		return []cleanupWarning{cleanupFailuref(
			"Agent Session cleanup skipped for %s %s in Run %s: %v",
			selection.ScopeKind,
			selection.ScopeID,
			run.ID,
			err,
		)}
	}

	cancelCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	cancelErr := commandDependenciesForContext(ctx).cancelStopAgentSession(cancelCtx, runtime, session)
	cancel()
	if cancelErr != nil && !agent.IsAgentSessionAbsent(cancelErr) {
		warnings = append(warnings, cleanupFailuref(
			"Agent Session cancel failed for %s %s (%s): %v",
			selection.ScopeKind,
			selection.ScopeID,
			session.Name,
			cancelErr,
		))
	}

	closeCtx, closeCancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	closeErr := commandDependenciesForContext(ctx).closeStopAgentSession(closeCtx, runtime, session)
	closeCancel()
	if closeErr != nil && !agent.IsAgentSessionAbsent(closeErr) {
		warnings = append(warnings, cleanupFailuref(
			"Agent Session close failed for %s %s (%s): %v",
			selection.ScopeKind,
			selection.ScopeID,
			session.Name,
			closeErr,
		))
		return warnings
	}
	if err := recordClosedAgentSelection(ctx, runStore, selection); err != nil {
		warnings = append(warnings, cleanupFailuref(
			"Agent Selection closed lifecycle failed for %s %s in Run %s: %v",
			selection.ScopeKind,
			selection.ScopeID,
			run.ID,
			err,
		))
	}
	return warnings
}

func registeredAgentSessionRef(run store.Run, selection store.AgentSelectionAttempt) (agent.SessionRef, error) {
	var session agent.SessionRef
	switch selection.ScopeKind {
	case store.AgentSelectionScopeTask:
		workDir := runSessionWorkDir(run)
		if taskWorkDir, ok := taskSessionWorkDir(run, selection.ScopeID); ok {
			workDir = taskWorkDir
		}
		session = agent.SessionRefForTask(run.ID, selection.ScopeID, workDir)
	case store.AgentSelectionScopeQA:
		session = agent.SessionRefForQA(run.ID, runSessionWorkDir(run))
	case store.AgentSelectionScopeReview:
		batchID, ok := strings.CutPrefix(strings.TrimSpace(selection.ScopeID), "batch-")
		if !ok {
			return agent.SessionRef{}, fmt.Errorf("review scope ID %q must use batch-NNN", selection.ScopeID)
		}
		batchNumber, err := strconv.Atoi(batchID)
		if err != nil || batchNumber <= 0 {
			return agent.SessionRef{}, fmt.Errorf("review scope ID %q must use a positive batch number", selection.ScopeID)
		}
		session = agent.SessionRefForReview(run.ID, batchNumber, run.GitRoot)
	default:
		return agent.SessionRef{}, fmt.Errorf("scope kind %q is unsupported", selection.ScopeKind)
	}
	if strings.TrimSpace(session.Name) == "" {
		return agent.SessionRef{}, errors.New("registered Agent Session name is empty")
	}
	if selection.FallbackIndex > 0 {
		session.Name = fmt.Sprintf("%s-fallback-%02d", session.Name, selection.FallbackIndex)
	}
	return session, nil
}

func recordClosedAgentSelection(ctx context.Context, runStore *store.Store, selection store.AgentSelectionAttempt) error {
	_, err := runStore.AppendAgentSelectionAttempt(context.WithoutCancel(ctx), store.AgentSelectionAttemptRequest{
		RunID:           selection.RunID,
		ScopeKind:       selection.ScopeKind,
		ScopeID:         selection.ScopeID,
		Category:        selection.Category,
		ProfileSource:   selection.ProfileSource,
		Attempt:         selection.Attempt,
		SelectionRole:   selection.SelectionRole,
		FallbackIndex:   selection.FallbackIndex,
		Runtime:         selection.Runtime,
		Model:           selection.Model,
		ReasoningEffort: selection.ReasoningEffort,
		Status:          store.AgentSelectionStatusClosed,
	})
	if err != nil {
		return fmt.Errorf("record closed Agent Selection lifecycle: %w", err)
	}
	return nil
}

func sessionRefForDiscoveredRunSession(run store.Run, session agent.RoundfixSession) agent.SessionRef {
	workDir := runSessionWorkDir(run)
	if strings.TrimSpace(session.TaskID) != "" {
		if taskWorkDir, ok := taskSessionWorkDir(run, session.TaskID); ok {
			workDir = taskWorkDir
		}
	}
	return agent.SessionRef{Name: session.Name, WorkDir: workDir}
}

func runSessionWorkDir(run store.Run) string {
	if workDir := strings.TrimSpace(run.WorkDir); workDir != "" {
		return workDir
	}
	return strings.TrimSpace(run.GitRoot)
}

func taskSessionWorkDir(run store.Run, taskID string) (string, bool) {
	runWorkDir := strings.TrimSpace(run.WorkDir)
	gitRoot := strings.TrimSpace(run.GitRoot)
	if runWorkDir == "" || gitRoot == "" || strings.TrimSpace(taskID) == "" {
		return "", false
	}
	taskRef, err := runworktree.TaskRefFor(runworktree.Ref{
		RunID:    run.ID,
		Path:     runWorkDir,
		Branch:   runworktree.BranchName(run.ID),
		UserRoot: gitRoot,
	}, taskID)
	if err != nil {
		return "", false
	}
	return taskRef.Path, true
}

func defaultListRoundfixAgentSessions(ctx context.Context, runtime agent.RuntimeSpec, workDir string) ([]agent.RoundfixSession, error) {
	return agent.NewDefaultRunner().ListRoundfixSessions(ctx, runtime, workDir)
}

func defaultCloseStopAgentSession(ctx context.Context, runtime agent.RuntimeSpec, session agent.SessionRef) error {
	return agent.NewDefaultRunner().CloseSession(ctx, runtime, session)
}

func defaultResolvePullRequestForStop(ctx context.Context, workDir string, pr string) (preflight.PullRequest, error) {
	return (preflight.GHPullRequestResolver{}).ResolvePullRequest(ctx, workDir, pr)
}

func runSkillsInstall(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	fs := flag.NewFlagSet("skills install", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	target := fs.String("target", "project", "Skill install target: project, codex, claude, opencode, or all")
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

	targetValue := strings.TrimSpace(*target)
	if targetValue == "" {
		targetValue = "project"
	}
	targetDirs := map[string]string{}
	if strings.TrimSpace(*dir) != "" {
		targetDirs[targetValue] = strings.TrimSpace(*dir)
	}
	projectRoot := ""
	if targetValue == "project" && strings.TrimSpace(*dir) == "" {
		var err error
		workDir, workDirErr := environment.resolveWorkDir("resolve work directory")
		if workDirErr != nil {
			fmt.Fprintf(stderr, "%s: skills install failed: %v\n", app.Name, workDirErr)
			return exitPreflight
		}
		projectRoot, err = commandDependenciesForContext(ctx).resolveSkillsProjectRoot(ctx, workDir)
		if err != nil {
			fmt.Fprintf(stderr, "%s: skills install failed: %v\n", app.Name, err)
			return exitPreflight
		}
	}
	result, err := roundskills.Install(ctx, roundskills.InstallRequest{
		Target:     targetValue,
		TargetDirs: targetDirs,
		ProjectDir: projectRoot,
	})
	if err != nil {
		fmt.Fprintf(stderr, "%s: skills install failed: %v\n", app.Name, err)
		return exitPreflight
	}
	for _, installed := range result.Targets {
		fmt.Fprintf(stdout, "Installed Roundfix skill for %s: %s (%d file(s))\n", installed.Target, installed.Dir, installed.Files)
	}
	if targetValue == "project" && projectRoot != "" {
		if err := maybeCreateProjectClaudeSkillSymlink(ctx, projectRoot, stdout, stderr); err != nil {
			fmt.Fprintf(stderr, "%s: skills install failed: %v\n", app.Name, err)
			return exitPreflight
		}
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

func defaultResolveSkillsProjectRoot(ctx context.Context, cwd string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if root := findSkillsGitRoot(cwd); root != "" {
		return root, nil
	}
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	cmd.Dir = cwd
	output, err := cmd.Output()
	if err != nil {
		return "", errors.New("project skill install requires running inside a git repository or passing --dir")
	}
	root := strings.TrimSpace(string(output))
	if root == "" {
		return "", errors.New("project skill install requires running inside a git repository or passing --dir")
	}
	return root, nil
}

func findSkillsGitRoot(start string) string {
	current := filepath.Clean(start)
	for {
		if _, err := os.Stat(filepath.Join(current, ".git")); err == nil {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			return ""
		}
		current = parent
	}
}

func maybeCreateProjectClaudeSkillSymlink(ctx context.Context, projectRoot string, stdout, stderr io.Writer) error {
	claudeSkillsDir := filepath.Join(projectRoot, ".claude", "skills")
	info, err := os.Stat(claudeSkillsDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect project Claude skills directory %q: %w", claudeSkillsDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("project Claude skills path %q is not a directory", claudeSkillsDir)
	}

	targetPath := filepath.Join(projectRoot, ".agents", "skills", "roundfix")
	linkPath := filepath.Join(claudeSkillsDir, "roundfix")
	if linkInfo, err := os.Lstat(linkPath); err == nil {
		if linkInfo.Mode()&os.ModeSymlink != 0 {
			existingTarget, readErr := os.Readlink(linkPath)
			if readErr == nil && sameSkillSymlinkTarget(linkPath, existingTarget, targetPath) {
				fmt.Fprintf(stdout, "Claude skill symlink already exists: %s -> %s\n", linkPath, existingTarget)
				return nil
			}
		}
		fmt.Fprintf(stdout, "Claude skill symlink skipped: %s already exists\n", linkPath)
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect Claude skill symlink %q: %w", linkPath, err)
	}

	relativeTarget, err := filepath.Rel(claudeSkillsDir, targetPath)
	if err != nil {
		return fmt.Errorf("resolve Claude skill symlink target: %w", err)
	}
	create, err := commandDependenciesForContext(ctx).promptProjectClaudeSkillSymlink(ctx, stderr, linkPath, relativeTarget)
	if err != nil {
		return err
	}
	if !create {
		fmt.Fprintf(stdout, "Claude skill symlink skipped by user: %s\n", linkPath)
		return nil
	}
	if err := os.Symlink(relativeTarget, linkPath); err != nil {
		return fmt.Errorf("create Claude skill symlink %q -> %q: %w", linkPath, relativeTarget, err)
	}
	fmt.Fprintf(stdout, "Created Claude skill symlink: %s -> %s\n", linkPath, relativeTarget)
	return nil
}

func sameSkillSymlinkTarget(linkPath string, existingTarget string, expectedTarget string) bool {
	if !filepath.IsAbs(existingTarget) {
		existingTarget = filepath.Join(filepath.Dir(linkPath), existingTarget)
	}
	return filepath.Clean(existingTarget) == filepath.Clean(expectedTarget)
}

func defaultPromptProjectClaudeSkillSymlink(ctx context.Context, stderr io.Writer, linkPath string, target string) (bool, error) {
	return readProjectClaudeSkillSymlinkPrompt(ctx, os.Stdin, stderr, linkPath, target)
}

func readProjectClaudeSkillSymlinkPrompt(ctx context.Context, stdin io.Reader, stderr io.Writer, linkPath string, target string) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if _, err := fmt.Fprintln(stderr, "Project Claude skills directory detected."); err != nil {
		return false, fmt.Errorf("write Claude skill symlink prompt: %w", err)
	}
	if _, err := fmt.Fprintf(stderr, "Create symlink %s -> %s? [Y/n]: ", linkPath, target); err != nil {
		return false, fmt.Errorf("write Claude skill symlink prompt: %w", err)
	}
	line, err := bufio.NewReader(stdin).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, fmt.Errorf("read Claude skill symlink prompt: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	return normalizeYesNo(line)
}

func normalizeYesNo(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("unsupported symlink response %q; enter yes or no", strings.TrimSpace(value))
	}
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
	if req.noAgentConsole && req.interactive {
		return req, validationError{message: "--interactive cannot be used with --no-agent-console"}
	}
	if !shouldOpenInteractiveInput(req) {
		return req, nil
	}
	inputReq, err := buildInteractiveInputRequest(ctx, req, loaded, stderr)
	if err != nil {
		return req, err
	}
	values, err := commandDependenciesForContext(ctx).collectInteractiveInput(ctx, inputReq)
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
	if req.name == "implement" {
		// A missing Spec is the primary trigger: the built-in config default
		// fills --agent, so an empty Agent only happens when the flag
		// explicitly clears it.
		return strings.TrimSpace(req.spec) == "" || strings.TrimSpace(req.agent) == ""
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

func buildInteractiveInputRequest(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, stderr io.Writer) (roundtui.InputRequest, error) {
	var specOptions []string
	if req.name == "implement" {
		// List the picker's Specs before any Run Database access so a
		// nothing-to-implement failure leaves no side effects behind.
		resolvedSpecsRoot, err := roundconfig.ResolveSpecsRoot(loaded, loaded.GitRoot)
		if err != nil {
			return roundtui.InputRequest{}, err
		}
		options, skipped, err := implementSpecOptionsDetailed(resolvedSpecsRoot.Path)
		if err != nil {
			return roundtui.InputRequest{}, err
		}
		printSkippedSpecDiagnostics(stderr, skipped)
		specOptions = options
	}
	remembered, err := loadRememberedInteractiveDefaults(ctx, loaded.HomeDir)
	if err != nil {
		return roundtui.InputRequest{}, err
	}
	var prSuggestion roundtui.Suggestion
	if req.name != "implement" {
		currentPR, _ := commandDependenciesForContext(ctx).suggestCurrentPullRequest(ctx, loaded.GitRoot)
		prSuggestion = roundtui.Suggestion{Value: currentPR, Source: "current"}
		if prSuggestion.Value == "" {
			prSuggestion = roundtui.Suggestion{Value: remembered.PRNumber, Source: "remembered"}
		}
	}
	agentSuggestion := roundtui.Suggestion{Value: req.agent, Source: "config"}
	if agentSuggestion.Value == "" {
		agentSuggestion = roundtui.Suggestion{Value: remembered.Agent, Source: "remembered"}
	}
	return roundtui.InputRequest{
		Command: req.name,
		Values: roundtui.CommandValues{
			PRNumber:        req.pr,
			Spec:            req.spec,
			ReviewSource:    req.source,
			Agent:           req.agent,
			Round:           req.round,
			ArtifactDir:     req.artifactDir,
			Model:           req.model,
			ReasoningEffort: req.reasoningEffort,
			MaxRounds:       req.maxRounds,
			UntilClean:      req.untilClean,
		},
		PRSuggestion:      prSuggestion,
		AgentSuggestion:   agentSuggestion,
		SelectionDefaults: tuiSelectionDefaults(loaded.Config),
		SpecOptions:       specOptions,
	}, nil
}

func applyInteractiveValues(req commandRequest, values roundtui.CommandValues) commandRequest {
	req.pr = strings.TrimSpace(values.PRNumber)
	req.spec = strings.TrimSpace(values.Spec)
	req.source = strings.TrimSpace(values.ReviewSource)
	req.agent = strings.TrimSpace(values.Agent)
	req.round = strings.TrimSpace(values.Round)
	req.artifactDir = strings.TrimSpace(values.ArtifactDir)
	if model := strings.TrimSpace(values.Model); model != "" {
		req.model = model
		req.modelSet = true
	} else if !req.modelSet {
		req.model = ""
	}
	if reasoningEffort := strings.TrimSpace(values.ReasoningEffort); reasoningEffort != "" {
		req.reasoningEffort = reasoningEffort
		req.reasoningEffortSet = true
	} else if !req.reasoningEffortSet {
		req.reasoningEffort = ""
	}
	if req.modelSet || req.reasoningEffortSet {
		req.agentSet = true
	}
	if values.MaxRounds > 0 {
		req.maxRounds = values.MaxRounds
	}
	req.untilClean = values.UntilClean
	return req
}

func tuiSelectionDefaults(config roundconfig.Config) map[string]roundtui.RuntimeSelectionDefaults {
	defaults := map[string]roundtui.RuntimeSelectionDefaults{}
	for _, runtime := range []string{"codex", "claude", "opencode"} {
		runtimeDefaults, _ := config.Runtimes.DefaultsFor(runtime)
		defaults[runtime] = roundtui.RuntimeSelectionDefaults{
			Model:            runtimeDefaults.Model,
			ReasoningEffort:  runtimeDefaults.ReasoningEffort,
			ModelCatalog:     tuiModelCatalog(runtime),
			ReasoningChoices: reasoningEffortChoices(runtime),
		}
	}
	return defaults
}

func tuiModelCatalog(runtime string) []roundtui.ModelChoice {
	catalog := agent.ModelCatalog(runtime)
	choices := make([]roundtui.ModelChoice, 0, len(catalog))
	for _, choice := range catalog {
		choices = append(choices, roundtui.ModelChoice{
			Label:       choice.Label,
			Value:       choice.Value,
			Description: choice.Description,
		})
	}
	return choices
}

func reasoningEffortChoices(runtime string) []string {
	switch runtime {
	case "codex":
		return []string{"low", "medium", "high", "xhigh"}
	case "claude":
		return []string{"default", "high", "maximum"}
	default:
		return nil
	}
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
	// The spec slug is deliberately never remembered: each Run's target is
	// an explicit choice.
	defaults := store.InteractiveDefaults{PRNumber: req.pr}
	if req.name == "resolve" || req.name == "watch" || req.name == "implement" {
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

func runOperationalCommand(ctx context.Context, name string, args []string, stdout, stderr io.Writer, detachChild *detachChild, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage(name))
		return exitOK
	}
	if name != "fetch" {
		if err := validateSelectionOverrideArgs(args); err != nil {
			printPreflightFailure(name, err, stderr)
			return exitPreflight
		}
	}

	loadedConfig, err := loadCommandConfig(environment, stderr)
	if err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}
	maybeReportVersionFreshness(ctx, loadedConfig, stderr)

	req, err := parseOperationalCommand(name, args, loadedConfig.Config)
	if err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}
	req.detachChild = detachChild
	req.branchIntegrityActor = environment.branchActor
	req = applyDetachSemantics(req)

	req, err = maybeCollectInteractiveInput(ctx, req, loadedConfig, stderr)
	if err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}
	if err := validateCommandRequest(req); err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}
	if err := validateAgentConsoleDisplay(req, stderr); err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}
	if req.detach {
		return runDetachedCommand(append([]string{name}, args...), req, loadedConfig, stdout, stderr, environment.environ, environment.workDir, environment.dependencies.detachTimeouts)
	}

	explicitArtifactDir := strings.TrimSpace(req.artifactDir) != ""
	artifactDir, err := resolveArtifactDirectoryForPreflight(req.artifactDir, loadedConfig.GitRoot, loadedConfig.HomeDir)
	if err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}
	req.artifactDir = artifactDir
	req.explicitArtifactDir = explicitArtifactDir

	preflightResult, err := commandDependenciesForContext(ctx).runCommandPreflight(ctx, req, loadedConfig)
	if err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}
	specsRoot := reviewArtifactSpecsRoot(loadedConfig, preflightResult.Git.Root)
	if !reviewArtifactUsesDefaultSpecsRoot(loadedConfig.Config.Specs.Root) {
		specsRoot, err = roundconfig.ResolveSpecsRoot(loadedConfig, preflightResult.Git.Root)
		if err != nil {
			printPreflightFailure(name, err, stderr)
			return exitPreflight
		}
	}
	reviewRoot, err := resolveReviewArtifactRoot(ctx, req, preflightResult, specsRoot)
	if err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}
	req.reviewRoot = reviewRoot
	// Clean-tree validation runs before Branch Integrity so a refusal never
	// follows a fast-forward that already moved the user's branch.
	if err := requireCleanTrackedReviewTree(ctx, req.name, preflightResult); err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}
	branchIntegrity, updatedPreflightResult, err := runBranchIntegrityPreflight(ctx, req, loadedConfig, preflightResult, stderr)
	if err != nil {
		printPreflightFailure(name, err, stderr)
		return exitPreflight
	}
	req.branchIntegrity = branchIntegrity
	preflightResult = updatedPreflightResult
	switch req.name {
	case "fetch":
		return runFetchCommand(ctx, req, loadedConfig, preflightResult, stdout, stderr)
	case "resolve":
		return runResolveCommand(ctx, req, loadedConfig, preflightResult, outcomeNotifierFromConfig(ctx, loadedConfig.Config), stdout, stderr)
	case "watch":
		return runWatchCommand(ctx, req, loadedConfig, preflightResult, outcomeNotifierFromConfig(ctx, loadedConfig.Config), stdout, stderr)
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

func resolveArtifactDirectoryForPreflight(artifactDir string, gitRoot string, homeDir string) (string, error) {
	resolved, err := roundconfig.ResolveArtifactDirectory(artifactDir, gitRoot, homeDir)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(resolved)
	if errors.Is(err, os.ErrNotExist) {
		return resolved, nil
	}
	if err != nil {
		return "", fmt.Errorf("stat Artifact Directory %q: %w", resolved, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("Artifact Directory %q is not a directory", resolved)
	}
	return resolved, nil
}

type branchIntegrityReport struct {
	Pending             []runworktree.PendingRunWork
	SupersededQAReports []runworktree.PendingRunWork
	Disregarded         []branchIntegrityDisregardedWork
	Preserved           []runworktree.PendingRunWork
	Integrated          []runworktree.PendingRunWork
	ActiveRun           *store.Run
}

type branchIntegrityDisregardedWork struct {
	runworktree.PendingRunWork
	Proof string
}

type branchIntegrityPendingError struct {
	HeadBranch          string
	Pending             []runworktree.PendingRunWork
	SupersededQAReports []runworktree.PendingRunWork
}

func (err branchIntegrityPendingError) Error() string {
	var builder strings.Builder
	headBranch := strings.TrimSpace(err.HeadBranch)
	if headBranch == "" {
		headBranch = "<unknown>"
	}
	fmt.Fprintf(&builder, "Branch Integrity Preflight refused pending Run Branch work for PR Head Branch %q.", headBranch)
	if len(err.Pending) > 0 && len(err.SupersededQAReports) > 0 {
		builder.WriteString("\nPending task work:")
	}
	for _, pending := range err.Pending {
		fmt.Fprintf(&builder, "\n- branch=%s worktree=%s ahead_commits=%d integration_command=%q",
			pending.Branch, branchIntegrityWorktreePath(pending.WorktreePath), pending.AheadCommits, branchIntegrityIntegrationCommand(pending.Branch))
	}
	if len(err.Pending) > 0 {
		builder.WriteString("\nNext action: inspect each pending Run Worktree, then run the listed integration command from the repository root when it is safe.")
	}
	if len(err.SupersededQAReports) > 0 {
		if len(err.Pending) > 0 {
			builder.WriteString("\nSuperseded QA reports:")
		}
		for _, pending := range err.SupersededQAReports {
			fmt.Fprintf(&builder, "\n- branch=%s worktree=%s ahead_commits=%d superseded QA report — release with: roundfix reconcile --apply",
				pending.Branch, branchIntegrityWorktreePath(pending.WorktreePath), pending.AheadCommits)
		}
	}
	return builder.String()
}

type branchIntegrityActiveRunError struct {
	HeadRepository string
	HeadBranch     string
	Run            store.Run
}

func (err branchIntegrityActiveRunError) Error() string {
	return fmt.Sprintf(
		"Branch Integrity Preflight refused because Active Run %s is bound to Head Repository %q and PR Head Branch %q. Next action: %s. If the owning process is dead or runaway: %s.",
		err.Run.ID,
		err.HeadRepository,
		err.HeadBranch,
		branchIntegrityStopCommand(err.Run.ID, false),
		branchIntegrityStopCommand(err.Run.ID, true),
	)
}

type dirtyTrackedReviewTreeError struct {
	Root    string
	Changes []preflight.ChangedPath
}

func (err dirtyTrackedReviewTreeError) Error() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "review Run Preflight Validation refused because tracked files are dirty in checkout %s.", err.Root)
	for _, change := range err.Changes {
		fmt.Fprintf(&builder, "\n- %s %s", change.Status, change.Path)
	}
	builder.WriteString("\nNext action: stash or commit tracked changes before running Roundfix; untracked files may remain.")
	return builder.String()
}

func runBranchIntegrityPreflight(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, stderr io.Writer) (branchIntegrityReport, preflight.Result, error) {
	report, err := inspectBranchIntegrity(ctx, loaded.HomeDir, preflightResult)
	if err != nil {
		return report, preflightResult, err
	}
	if req.skipBranchIntegrity {
		return report, preflightResult, nil
	}
	if report.ActiveRun != nil {
		runStore, err := store.Open(ctx, loaded.HomeDir)
		if err != nil {
			return report, preflightResult, err
		}
		for report.ActiveRun != nil {
			_, ok, reclaimErr := reclaimOrphanedActiveRun(ctx, runStore, *report.ActiveRun, stderr)
			if reclaimErr != nil {
				closeErr := runStore.Close()
				if closeErr != nil {
					return report, preflightResult, errors.Join(reclaimErr, closeErr)
				}
				return report, preflightResult, reclaimErr
			}
			if !ok {
				active := *report.ActiveRun
				closeErr := runStore.Close()
				if closeErr != nil {
					return report, preflightResult, closeErr
				}
				return report, preflightResult, branchIntegrityActiveRunError{
					HeadRepository: preflightResult.PullRequest.HeadRepository,
					HeadBranch:     preflightResult.PullRequest.HeadBranch,
					Run:            active,
				}
			}
			active, found, err := runStore.ActiveReviewRunByTarget(ctx, preflightResult.PullRequest.HeadRepository, preflightResult.PullRequest.HeadBranch)
			if err != nil {
				closeErr := runStore.Close()
				if closeErr != nil {
					return report, preflightResult, errors.Join(err, closeErr)
				}
				return report, preflightResult, err
			}
			if !found {
				report.ActiveRun = nil
				break
			}
			report.ActiveRun = &active
		}
		closeErr := runStore.Close()
		if closeErr != nil {
			return report, preflightResult, closeErr
		}
	}
	for _, disregarded := range report.Disregarded {
		fmt.Fprintf(
			stderr,
			"Branch Integrity Preflight: disregarded branch=%s worktree=%s ahead_commits=%d proof=%s; Git ref left unchanged.\n",
			disregarded.Branch,
			branchIntegrityWorktreePath(disregarded.WorktreePath),
			disregarded.AheadCommits,
			disregarded.Proof,
		)
	}
	actionable := branchIntegrityActionablePending(report.Pending, report.Disregarded)
	if len(actionable) == 0 {
		return report, preflightResult, nil
	}
	var blocked []runworktree.PendingRunWork
	for _, pending := range actionable {
		if !pending.FastForward || branchIntegrityContainsPending(report.Preserved, pending.Branch) {
			blocked = append(blocked, pending)
		}
	}
	if len(blocked) > 0 {
		return report, preflightResult, branchIntegrityPendingError{
			HeadBranch:          preflightResult.PullRequest.HeadBranch,
			Pending:             blocked,
			SupersededQAReports: report.SupersededQAReports,
		}
	}
	integrationPlan, err := branchIntegrityIntegrationPlan(ctx, preflightResult.Git.Root, preflightResult.PullRequest.HeadBranch, actionable)
	if err != nil {
		return report, preflightResult, err
	}
	for _, pending := range integrationPlan {
		if err := commandDependenciesForContext(ctx).integratePendingRunWork(ctx, preflightResult.Git.Root, preflightResult.PullRequest.HeadBranch, pending.Branch); err != nil {
			return report, preflightResult, fmt.Errorf("Branch Integrity Preflight integrate %s: %w", pending.Branch, err)
		}
		report.Integrated = append(report.Integrated, pending)
		fmt.Fprintf(stderr, "Branch Integrity Preflight: integrated %s with %s (%d commit(s), worktree %s).\n",
			pending.Branch, branchIntegrityIntegrationCommand(pending.Branch), pending.AheadCommits, branchIntegrityWorktreePath(pending.WorktreePath))
	}
	updated, err := commandDependenciesForContext(ctx).refreshBranchIntegrityHead(ctx, preflightResult)
	if err != nil {
		return report, preflightResult, err
	}
	return report, updated, nil
}

func branchIntegrityActionablePending(
	pending []runworktree.PendingRunWork,
	disregarded []branchIntegrityDisregardedWork,
) []runworktree.PendingRunWork {
	actionable := make([]runworktree.PendingRunWork, 0, len(pending))
	for _, work := range pending {
		if branchIntegrityContainsDisregarded(disregarded, work.Branch) {
			continue
		}
		actionable = append(actionable, work)
	}
	return actionable
}

func branchIntegrityContainsDisregarded(disregarded []branchIntegrityDisregardedWork, branch string) bool {
	return slices.ContainsFunc(disregarded, func(work branchIntegrityDisregardedWork) bool {
		return work.Branch == branch
	})
}

func branchIntegrityContainsPending(pending []runworktree.PendingRunWork, branch string) bool {
	return slices.ContainsFunc(pending, func(work runworktree.PendingRunWork) bool {
		return work.Branch == branch
	})
}

func requireCleanTrackedReviewTree(ctx context.Context, commandName string, preflightResult preflight.Result) error {
	if commandName != "resolve" && commandName != "watch" {
		return nil
	}
	changes, err := commandDependenciesForContext(ctx).inspectChangedPaths(ctx, preflightResult.Git.Root)
	if err != nil {
		return fmt.Errorf("inspect clean tracked checkout for review Run: %w", err)
	}
	tracked := trackedReviewChanges(changes)
	if len(tracked) == 0 {
		return nil
	}
	return dirtyTrackedReviewTreeError{
		Root:    preflightResult.Git.Root,
		Changes: tracked,
	}
}

func trackedReviewChanges(changes []preflight.ChangedPath) []preflight.ChangedPath {
	tracked := make([]preflight.ChangedPath, 0, len(changes))
	for _, change := range changes {
		if strings.TrimSpace(change.Status) == "??" {
			continue
		}
		tracked = append(tracked, change)
	}
	return tracked
}

func branchIntegrityIntegrationPlan(ctx context.Context, gitRoot string, headBranch string, pending []runworktree.PendingRunWork) ([]runworktree.PendingRunWork, error) {
	planned := append([]runworktree.PendingRunWork(nil), pending...)
	sort.SliceStable(planned, func(i, j int) bool {
		if planned[i].AheadCommits == planned[j].AheadCommits {
			return planned[i].Branch < planned[j].Branch
		}
		return planned[i].AheadCommits < planned[j].AheadCommits
	})
	for index := 1; index < len(planned); index++ {
		ancestor := planned[index-1].Branch
		descendant := planned[index].Branch
		isAncestor, err := branchIntegrityIsAncestor(ctx, gitRoot, ancestor, descendant)
		if err != nil {
			return nil, err
		}
		if !isAncestor {
			return nil, branchIntegrityPendingError{HeadBranch: headBranch, Pending: pending}
		}
	}
	return planned, nil
}

func branchIntegrityIsAncestor(ctx context.Context, gitRoot string, ancestor string, descendant string) (bool, error) {
	if _, err := (preflight.ExecGitRunner{}).RunGit(ctx, gitRoot, "merge-base", "--is-ancestor", ancestor, descendant); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return false, nil
		}
		return false, fmt.Errorf("check pending Run Branch ancestry %s..%s: %w", ancestor, descendant, err)
	}
	return true, nil
}

func inspectBranchIntegrity(ctx context.Context, homeDir string, preflightResult preflight.Result) (branchIntegrityReport, error) {
	report := branchIntegrityReport{}
	pending, err := commandDependenciesForContext(ctx).listPendingRunWork(ctx, preflightResult.Git.Root, preflightResult.PullRequest.HeadBranch)
	if err != nil {
		return report, err
	}
	if _, err := os.Stat(store.DatabasePath(homeDir)); errors.Is(err, os.ErrNotExist) {
		// No Run Database yet: no Run can attribute the branches, and no
		// Active Run can exist. Keep every pending branch conservatively
		// without creating the database on a preflight that may still fail.
		report.Pending = pending
		return report, nil
	} else if err != nil {
		return report, fmt.Errorf("inspect Run Database before Branch Integrity Preflight: %w", err)
	}
	// Open the migrating store: every caller is an operational command that
	// creates Runs, so upgrading an existing Run Database here is in-contract
	// and keeps this inspection from querying an unmigrated schema.
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		return report, err
	}
	defer func() {
		_ = runStore.Close()
	}()
	report.Pending, report.SupersededQAReports, report.Disregarded, report.Preserved, err = filterPendingRunWorkByTarget(
		ctx,
		runStore,
		pending,
		preflightResult.Git.Root,
		preflightResult.Git.HEAD,
		preflightResult.PullRequest.HeadBranch,
	)
	if err != nil {
		return report, err
	}
	// Scan the runs table instead of the lock table: Runs created with the
	// Branch Integrity bypass hold no Active Run lock but must stay visible
	// to subsequent guard checks.
	active, found, err := runStore.ActiveReviewRunByTarget(ctx, preflightResult.PullRequest.HeadRepository, preflightResult.PullRequest.HeadBranch)
	if err != nil {
		return report, err
	}
	if found {
		report.ActiveRun = &active
	}
	return report, nil
}

// filterPendingRunWorkByTarget drops Run Branches whose recorded Run belongs
// to a different branch: git topology alone cannot tell a Run Branch based on
// the PR Head Branch from one based on another feature branch. Branches with
// no Run row are kept conservatively.
func filterPendingRunWorkByTarget(
	ctx context.Context,
	runStore *store.Store,
	pending []runworktree.PendingRunWork,
	gitRoot string,
	targetHead string,
	headBranch string,
) ([]runworktree.PendingRunWork, []runworktree.PendingRunWork, []branchIntegrityDisregardedWork, []runworktree.PendingRunWork, error) {
	headBranch = strings.TrimSpace(headBranch)
	filtered := make([]runworktree.PendingRunWork, 0, len(pending))
	supersededQAReports := make([]runworktree.PendingRunWork, 0, len(pending))
	disregarded := make([]branchIntegrityDisregardedWork, 0, len(pending))
	preserved := make([]runworktree.PendingRunWork, 0, len(pending))
	preservedBranches := make(map[string]struct{}, len(pending))
	preserve := func(work runworktree.PendingRunWork) {
		if _, exists := preservedBranches[work.Branch]; exists {
			return
		}
		preservedBranches[work.Branch] = struct{}{}
		preserved = append(preserved, work)
	}
	runs, err := runStore.ListRuns(ctx, store.ListRunsQuery{GitRoot: gitRoot, States: store.StatesAll})
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("list Runs for Branch Integrity Preflight: %w", err)
	}
	runsByID := make(map[string]store.Run, len(runs))
	for _, run := range runs {
		runsByID[run.ID] = run
	}
	type classifiedSpec struct {
		branches []runworktree.PendingRunWork
	}
	bySpec := make(map[string]*classifiedSpec)
	for _, work := range pending {
		runID, runBranch := runworktree.RunIDFromBranchName(work.Branch)
		if !runBranch {
			filtered = append(filtered, work)
			continue
		}
		row, found := runsByID[runID]
		if found && strings.TrimSpace(row.LocalBranch) != headBranch && strings.TrimSpace(row.HeadBranch) != headBranch {
			continue
		}
		filtered = append(filtered, work)
		if !found || row.Kind != store.KindImplement {
			continue
		}
		if !store.IsTerminalState(row.State) {
			preserve(work)
		}
		slug := strings.TrimSpace(row.SpecSlug)
		if slug == "" {
			preserve(work)
			continue
		}
		group := bySpec[slug]
		if group == nil {
			group = &classifiedSpec{}
			bySpec[slug] = group
		}
		group.branches = append(group.branches, work)
	}

	for slug, group := range bySpec {
		classification, classifyErr := commandDependenciesForContext(ctx).classifyRunBranchSet(
			ctx,
			gitRoot,
			headBranch,
			slug,
			runs,
		)
		if classifyErr != nil {
			for _, work := range group.branches {
				preserve(work)
			}
			continue
		}
		classified := make(map[string]struct{})
		for _, branch := range classification.Preserved {
			classified[branch] = struct{}{}
			if work, ok := branchIntegrityPendingByBranch(group.branches, branch); ok {
				preserve(work)
			}
		}
		if classification.Current != "" {
			classified[classification.Current] = struct{}{}
		}
		for _, branch := range classification.Releasable {
			classified[branch] = struct{}{}
			work, ok := branchIntegrityPendingByBranch(group.branches, branch)
			if !ok {
				continue
			}
			if _, mustPreserve := preservedBranches[work.Branch]; mustPreserve {
				continue
			}
			proofReport := strings.TrimSpace(classification.CurrentReport)
			proofPrefix := "superseded by current QA Report"
			if proofReport == "" {
				proofPrefix = "superseded by target QA Report"
				var proven bool
				proofReport, proven = commandDependenciesForContext(ctx).supersedingQAReport(ctx, gitRoot, targetHead, branch, slug)
				if !proven || strings.TrimSpace(proofReport) == "" {
					preserve(work)
					continue
				}
			}
			disregarded = append(disregarded, branchIntegrityDisregardedWork{
				PendingRunWork: work,
				Proof:          fmt.Sprintf("%s %q", proofPrefix, proofReport),
			})
			supersededQAReports = append(supersededQAReports, work)
		}
		for _, work := range group.branches {
			if _, ok := classified[work.Branch]; ok {
				continue
			}
			preserve(work)
		}
	}
	filtered = branchIntegrityWithoutPending(filtered, supersededQAReports)
	sort.Slice(disregarded, func(i, j int) bool {
		return disregarded[i].Branch < disregarded[j].Branch
	})
	sort.Slice(supersededQAReports, func(i, j int) bool {
		return supersededQAReports[i].Branch < supersededQAReports[j].Branch
	})
	sort.Slice(preserved, func(i, j int) bool {
		return preserved[i].Branch < preserved[j].Branch
	})
	return filtered, supersededQAReports, disregarded, preserved, nil
}

func branchIntegrityPendingByBranch(pending []runworktree.PendingRunWork, branch string) (runworktree.PendingRunWork, bool) {
	for _, work := range pending {
		if work.Branch == branch {
			return work, true
		}
	}
	return runworktree.PendingRunWork{}, false
}

func branchIntegrityWithoutPending(
	pending []runworktree.PendingRunWork,
	removed []runworktree.PendingRunWork,
) []runworktree.PendingRunWork {
	kept := make([]runworktree.PendingRunWork, 0, len(pending))
	for _, work := range pending {
		if slices.ContainsFunc(removed, func(removedWork runworktree.PendingRunWork) bool {
			return removedWork.Branch == work.Branch
		}) {
			continue
		}
		kept = append(kept, work)
	}
	return kept
}

func defaultRefreshBranchIntegrityHead(ctx context.Context, preflightResult preflight.Result) (preflight.Result, error) {
	if len(preflightResult.Git.Root) == 0 {
		return preflightResult, nil
	}
	head, err := preflight.ExecGitRunner{}.RunGit(ctx, preflightResult.Git.Root, "rev-parse", "HEAD")
	if err != nil {
		return preflightResult, fmt.Errorf("refresh HEAD after Branch Integrity Preflight: %w", err)
	}
	preflightResult.Git.HEAD = strings.TrimSpace(head)
	return preflightResult, nil
}

func branchIntegrityIntegrationCommand(branch string) string {
	return "git merge --ff-only " + strings.TrimSpace(branch)
}

func branchIntegrityStopCommand(runID string, force bool) string {
	if force {
		return "roundfix stop --force " + strings.TrimSpace(runID)
	}
	return "roundfix stop " + strings.TrimSpace(runID)
}

func branchIntegrityWorktreePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return "-"
	}
	return path
}

func publishBranchIntegrityBypassAudit(ctx context.Context, req commandRequest, preflightResult preflight.Result, runStore *store.Store, runID string) error {
	if !req.skipBranchIntegrity {
		return nil
	}
	prNumber, err := strconv.Atoi(preflightResult.PullRequest.Number)
	if err != nil {
		return fmt.Errorf("parse Open Pull Request number %q for Branch Integrity bypass audit: %w", preflightResult.PullRequest.Number, err)
	}
	baseRepository := strings.TrimSpace(preflightResult.PullRequest.BaseRepository)
	if baseRepository == "" {
		return fmt.Errorf("publish Branch Integrity Preflight bypass audit comment: Base Repository is unknown for Open Pull Request #%d", prNumber)
	}
	marker := branchIntegrityBypassMarker(runID, preflightResult.PullRequest.Number)
	body := branchIntegrityBypassAuditBody(runID, preflightResult, req.branchIntegrity, marker, time.Now().UTC(), req.branchIntegrityActor)
	if err := commandDependenciesForContext(ctx).commentOnPullRequest(ctx, req.source, baseRepository, prNumber, body); err != nil {
		return fmt.Errorf("publish Branch Integrity Preflight bypass audit comment: %w", err)
	}
	journalBranchIntegrityBypass(ctx, runStore, runID, req.branchIntegrity, marker)
	return nil
}

func branchIntegrityBypassMarker(runID string, prNumber string) string {
	return coderabbit.RoundfixCommentMarker("run:"+strings.TrimSpace(runID), "bypass:branch-integrity", "pr:"+strings.TrimSpace(prNumber))
}

func branchIntegrityBypassAuditBody(runID string, preflightResult preflight.Result, report branchIntegrityReport, marker string, now time.Time, actor string) string {
	var builder strings.Builder
	builder.WriteString("Roundfix Branch Integrity Preflight bypassed.\n\n")
	fmt.Fprintf(&builder, "Run: %s\n", runID)
	fmt.Fprintf(&builder, "Actor: %s\n", branchIntegrityAuditActor(actor))
	fmt.Fprintf(&builder, "Time: %s\n", now.UTC().Format(time.RFC3339))
	fmt.Fprintf(&builder, "Open Pull Request: #%s\n", preflightResult.PullRequest.Number)
	fmt.Fprintf(&builder, "Head Repository: %s\n", preflightResult.PullRequest.HeadRepository)
	fmt.Fprintf(&builder, "PR Head Branch: %s\n\n", preflightResult.PullRequest.HeadBranch)
	builder.WriteString("Skipped guardrails:\n")
	builder.WriteString("- Pending Run Branch work\n")
	builder.WriteString("- Active Run bound to target\n\n")
	builder.WriteString("Ignored pending Run Branch work:\n")
	if len(report.Pending) == 0 {
		builder.WriteString("- none\n")
	} else {
		for _, pending := range report.Pending {
			fmt.Fprintf(&builder, "- branch=%s worktree=%s ahead_commits=%d fast_forward=%s integration_command=%q\n",
				pending.Branch,
				branchIntegrityWorktreePath(pending.WorktreePath),
				pending.AheadCommits,
				yesNo(pending.FastForward),
				branchIntegrityIntegrationCommand(pending.Branch),
			)
		}
	}
	// Superseded entries were moved out of Pending before this audit runs, so
	// rendering only Pending would let a bypass silently omit a whole class of
	// ignored branch work. An audit that hides a class is worse than no audit.
	builder.WriteString("\nIgnored superseded QA-report Run Branch work:\n")
	if len(report.SupersededQAReports) == 0 {
		builder.WriteString("- none\n")
	} else {
		for _, superseded := range report.SupersededQAReports {
			fmt.Fprintf(&builder, "- branch=%s worktree=%s ahead_commits=%d release_command=%q\n",
				superseded.Branch,
				branchIntegrityWorktreePath(superseded.WorktreePath),
				superseded.AheadCommits,
				"roundfix reconcile --apply",
			)
		}
	}
	builder.WriteString("\nIgnored Active Runs:\n")
	if report.ActiveRun == nil {
		builder.WriteString("- none\n")
	} else {
		fmt.Fprintf(&builder, "- run_id=%s state=%s stop_command=%q force_stop_command=%q\n",
			report.ActiveRun.ID,
			report.ActiveRun.State,
			branchIntegrityStopCommand(report.ActiveRun.ID, false),
			branchIntegrityStopCommand(report.ActiveRun.ID, true),
		)
	}
	return coderabbit.RoundfixCommentBody(builder.String(), marker)
}

func branchIntegrityAuditActor(actor string) string {
	if actor = strings.TrimSpace(actor); actor != "" {
		return actor
	}
	return "unknown"
}

func journalBranchIntegrityIntegrations(ctx context.Context, runStore *store.Store, runID string, integrated []runworktree.PendingRunWork) {
	if runStore == nil || len(integrated) == 0 {
		return
	}
	sink := store.JournalSink{Store: runStore}
	for _, item := range integrated {
		payload, err := json.Marshal(map[string]any{
			"event":         "branch_integrity_auto_integration",
			"branch":        item.Branch,
			"worktree_path": item.WorktreePath,
			"ahead_commits": item.AheadCommits,
			"fast_forward":  item.FastForward,
			"command":       branchIntegrityIntegrationCommand(item.Branch),
		})
		if err != nil {
			continue
		}
		_ = sink.Publish(context.WithoutCancel(ctx), runevent.RunEvent{
			RunID:   runID,
			Source:  runevent.SourceGit,
			Kind:    runevent.KindDaemonStatus,
			Summary: runevent.BoundSummary(fmt.Sprintf("Branch Integrity Preflight integrated %s.", item.Branch)),
			Time:    time.Now().UTC(),
			Payload: payload,
		})
	}
}

func journalBranchIntegrityBypass(ctx context.Context, runStore *store.Store, runID string, report branchIntegrityReport, marker string) {
	if runStore == nil {
		return
	}
	payload, err := json.Marshal(map[string]any{
		"event":              "branch_integrity_bypass",
		"marker":             marker,
		"skipped_guardrails": []string{"pending_run_branch_work", "active_run_bound_to_target"},
		"pending":            branchIntegrityPendingPayload(report.Pending),
		"active_run":         branchIntegrityActiveRunPayload(report.ActiveRun),
	})
	if err != nil {
		return
	}
	_ = (store.JournalSink{Store: runStore}).Publish(context.WithoutCancel(ctx), runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonStatus,
		Summary: runevent.BoundSummary("Branch Integrity Preflight bypass audited on the pull request."),
		Time:    time.Now().UTC(),
		Payload: payload,
	})
}

func branchIntegrityPendingPayload(pending []runworktree.PendingRunWork) []map[string]any {
	payload := make([]map[string]any, 0, len(pending))
	for _, item := range pending {
		payload = append(payload, map[string]any{
			"branch":        item.Branch,
			"worktree_path": item.WorktreePath,
			"ahead_commits": item.AheadCommits,
			"fast_forward":  item.FastForward,
			"command":       branchIntegrityIntegrationCommand(item.Branch),
		})
	}
	return payload
}

func branchIntegrityActiveRunPayload(active *store.Run) map[string]any {
	if active == nil {
		return nil
	}
	return map[string]any{
		"run_id":             active.ID,
		"state":              active.State,
		"stop_command":       branchIntegrityStopCommand(active.ID, false),
		"force_stop_command": branchIntegrityStopCommand(active.ID, true),
	}
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

	run, err := createFetchRun(ctx, runStore, req, preflightResult, stderr)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	journalBranchIntegrityIntegrations(ctx, runStore, run.ID, req.branchIntegrity.Integrated)
	if err := publishBranchIntegrityBypassAudit(ctx, req, preflightResult, runStore, run.ID); err != nil {
		markRunFailed(ctx, runStore, run.ID)
		printBranchIntegrityAuditFailure(req.name, run.ID, err, stderr)
		return exitPreflight
	}
	if err := ensureReviewRunArtifactDirectory(req, loaded, preflightResult); err != nil {
		markRunFailed(ctx, runStore, run.ID)
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}
	if err := rememberInteractiveDefaults(ctx, runStore, req); err != nil {
		markRunFailed(ctx, runStore, run.ID)
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}
	printLiveRunView(stderr, req, loaded, preflightResult, run.ID, "FetchingIssues", nil, []string{"Fetching Review Source issues..."})

	items, err := commandDependenciesForContext(ctx).fetchReviewItems(ctx, reviewsource.FetchRequest{
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
		ReviewRoot:     req.reviewRoot,
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

func createFetchRun(ctx context.Context, runStore *store.Store, req commandRequest, preflightResult preflight.Result, stderr io.Writer) (store.Run, error) {
	createReq := store.CreateRunRequest{
		Kind:           store.KindFetch,
		HeadRepository: preflightResult.PullRequest.HeadRepository,
		HeadBranch:     preflightResult.PullRequest.HeadBranch,
		BaseRepository: preflightResult.PullRequest.BaseRepository,
		PRNumber:       preflightResult.PullRequest.Number,
		GitRoot:        preflightResult.Git.Root,
		LocalBranch:    preflightResult.Git.Branch,
		HeadSHA:        preflightResult.Git.HEAD,
		ArtifactDir:    req.artifactDir,
		WorkDir:        preflightResult.Git.Root,
	}
	return createReviewRun(ctx, runStore, req, createReq, stderr)
}

type resolveBatchPlan struct {
	roundNumber     int
	selection       rounds.SelectResult
	plan            rounds.BatchPlan
	runtime         agent.RuntimeSpec
	agentSelections daemon.AgentSelectionProfiles
	runtimeFactory  daemon.AgentRuntimeFactory
}

type resolveBatchResult struct {
	Remaining     int
	CommitCreated bool
}

func runResolveCommand(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, notifier roundnotify.Notifier, stdout, stderr io.Writer) int {
	resolvePlan, err := prepareResolveBatch(ctx, req, loaded, preflightResult)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	collaborators := commandDependenciesForContext(ctx).newEngineCollaborators()
	categories := reviewProfileCategories()
	profilePreflight, err := runProfileOperationalPreflight(ctx, req, loaded.Config, categories, preflightResult.Git.Root, collaborators.runner, stderr)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	resolvePlan.runtime, err = runtimeForOperationalProfileRun(req, loaded.Config, categories, profilePreflight.Override)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	resolvePlan.agentSelections, err = operationalAgentSelectionProfiles(loaded.Config, categories, profilePreflight.Override)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	resolvePlan.runtimeFactory = operationalRuntimeFactory(req)
	req = requestWithRuntimeSelection(req, resolvePlan.runtime)

	runStore, err := store.Open(ctx, loaded.HomeDir)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	defer func() {
		_ = runStore.Close()
	}()
	sweepRunRetention(ctx, runStore, req.artifactDir, loaded.Config.Store.JournalRetention, stderr)
	run, err := createOperationalRun(ctx, runStore, store.KindResolve, req, preflightResult, resolvePlan.runtime, stderr)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	journalBranchIntegrityIntegrations(ctx, runStore, run.ID, req.branchIntegrity.Integrated)
	if err := publishBranchIntegrityBypassAudit(ctx, req, preflightResult, runStore, run.ID); err != nil {
		markRunFailedAndNotify(ctx, runStore, run.ID, notifier, stderr)
		printBranchIntegrityAuditFailure(req.name, run.ID, err, stderr)
		return exitPreflight
	}
	if err := ensureReviewRunArtifactDirectory(req, loaded, preflightResult); err != nil {
		markRunFailedAndNotify(ctx, runStore, run.ID, notifier, stderr)
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}
	if err := req.reportDetachedRunCreated(run.ID); err != nil {
		markRunFailedAndNotify(ctx, runStore, run.ID, notifier, stderr)
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}
	session := agent.SessionRefForRun(run.ID, preflightResult.Git.Root)
	sessionForClose := session
	if len(resolvePlan.agentSelections) > 0 {
		sessionForClose = agent.SessionRef{}
	}
	if err := rememberInteractiveDefaults(ctx, runStore, req); err != nil {
		closeAgentSession(ctx, collaborators.runner, resolvePlan.runtime, sessionForClose, run.ID, runStore)
		markRunFailedAndNotify(ctx, runStore, run.ID, notifier, stderr)
		printResolveRunFailure(err, stderr)
		return exitRunFailed
	}
	cockpitView := buildLiveRunView(req, loaded, preflightResult, run.ID, "ResolvingWithAgent", resolvePlan.selection.Issues, nil, stderr)
	cockpitView.WorkDir = preflightResult.Git.Root
	if !liveTUIEnabled(stderr) {
		plainView := buildLiveRunView(req, loaded, preflightResult, run.ID, "ResolvingWithAgent", resolvePlan.selection.Issues, []string{"Agent and verification output will stream below."}, stderr)
		plainView.WorkDir = preflightResult.Git.Root
		fmt.Fprint(stderr, roundtui.RenderLiveRunView(plainView))
	}
	for _, batch := range resolvePlan.plan.Batches {
		cockpitView.BatchSizes = append(cockpitView.BatchSizes, len(batch.Issues))
	}
	ui, err := startRunUI(ctx, cockpitView, run.ID, loaded.HomeDir, runStore, stderr, req.noAgentConsole)
	if err != nil {
		closeAgentSession(ctx, collaborators.runner, resolvePlan.runtime, sessionForClose, run.ID, runStore)
		markRunFailedAndNotify(ctx, runStore, run.ID, notifier, stderr)
		printResolveRunFailure(err, stderr)
		return exitRunFailed
	}
	defer ui.Close()

	writeGuard := reviewRunTargetGuard(ctx, run)
	cycleResult, err := executeResolveCycle(ctx, req, loaded, preflightResult, run.ID, session, resolvePlan, collaborators, runStore, ui, writeGuard)
	if err != nil {
		if isStopRequest(ctx, err) {
			closeAgentSession(ctx, collaborators.runner, resolvePlan.runtime, sessionForClose, run.ID, runStore)
			code := completeStoppedRunRecord(runStore, run.ID, notifier, stderr)
			ui.Wait()
			ui.Close()
			if code != exitOK {
				printRunFailure(req.name, errors.New("complete stopped Run"), stderr)
				return code
			}
			fmt.Fprintf(stderr, "%s Run %s reached %s.\n", commandDisplayName(req.name), run.ID, store.StateStopped)
			printStopSummary(ctx, req, preflightResult, stderr)
			printReviewIssueReport(stdout, store.StateStopped, 1, true, reviewIssueReportData(context.WithoutCancel(ctx), req, preflightResult, resolvePlan.selection.Issues, stderr))
			return exitOK
		}
		var checkoutMoved watch.CheckoutMovedError
		if errors.As(err, &checkoutMoved) {
			closeAgentSession(ctx, collaborators.runner, resolvePlan.runtime, sessionForClose, run.ID, runStore)
			completed, completeErr := runStore.CompleteRun(context.WithoutCancel(ctx), run.ID, store.StateCheckoutMoved)
			if completeErr != nil {
				ui.Close()
				printResolveRunFailure(completeErr, stderr)
				return exitRunFailed
			}
			publishTerminalCompletion(context.WithoutCancel(ctx), runStore, notifier, stderr, completed, cycleResult.Remaining)
			ui.Wait()
			ui.Close()
			fmt.Fprintf(stderr, "Resolve Run %s reached %s: %v\n", completed.ID, completed.State, checkoutMoved)
			printReviewIssueReport(stdout, completed.State, 1, true, reviewIssueReportData(context.WithoutCancel(ctx), req, preflightResult, resolvePlan.selection.Issues, stderr))
			return exitRunFailed
		}
		closeAgentSession(ctx, collaborators.runner, resolvePlan.runtime, sessionForClose, run.ID, runStore)
		markRunFailedAndNotify(ctx, runStore, run.ID, notifier, stderr)
		ui.Wait()
		ui.Close()
		printResolveRunFailure(err, stderr)
		return exitRunFailed
	}

	outcome := store.StateClean
	if cycleResult.Remaining > 0 {
		outcome = store.StateUnresolved
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, outcome)
	if err != nil {
		ui.Close()
		closeAgentSession(ctx, collaborators.runner, resolvePlan.runtime, sessionForClose, run.ID, runStore)
		printResolveRunFailureAfterBatchCommit(err, stderr)
		return exitRunFailed
	}
	closeAgentSession(ctx, collaborators.runner, resolvePlan.runtime, sessionForClose, completed.ID, runStore)
	publishTerminalCompletion(ctx, runStore, notifier, stderr, completed, cycleResult.Remaining)
	// The cockpit stays on screen, read-only, until the user closes it.
	ui.Wait()
	ui.Close()
	fmt.Fprintf(stderr, "Resolve Run %s reached %s.\n", completed.ID, completed.State)
	if completed.State == store.StateUnresolved {
		fmt.Fprintf(stderr, "%d Unresolved Review Issue(s) remain; failed issues are retried by the next fetched Round.\n", cycleResult.Remaining)
		printAgentCheckoutChangesNotice(stderr)
	}
	printReviewIssueReport(stdout, completed.State, 1, true, reviewIssueReportData(context.WithoutCancel(ctx), req, preflightResult, resolvePlan.selection.Issues, stderr))
	if completed.State == store.StateUnresolved {
		return exitRunFailed
	}
	return exitOK
}

// completeStoppedRunRecord finishes a stopped Run in the Run Database and
// journals the outcome, without printing: callers print after the cockpit
// leaves the screen.
func completeStoppedRunRecord(runStore *store.Store, runID string, notifier roundnotify.Notifier, stderr io.Writer) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	completed, err := runStore.CompleteRun(ctx, runID, store.StateStopped)
	if err != nil {
		return exitRunFailed
	}
	publishTerminalCompletion(ctx, runStore, notifier, stderr, completed, 0)
	return exitOK
}

type reviewIssueReportCounts struct {
	resolved   int
	invalid    int
	duplicated int
	failed     int
	unresolved int
}

type reviewIssueReport struct {
	runIssues             []rounds.Issue
	cumulativeIssues      []rounds.Issue
	cumulativeUnavailable bool
}

func printReviewSkippedReport(stdout io.Writer, reason string, nextAction string) {
	fmt.Fprintf(stdout, "Review Source: skipped — reason: %s\n", reason)
	fmt.Fprintf(stdout, "Next action: %s\n", nextAction)
}

func formatReviewWaitProgress(progress watch.WaitProgress) string {
	return fmt.Sprintf(
		"Review Source status: %s; phase=%s; expected_head=%s; started_at=%s; deadline=%s; evidence_kind=%s; retry=%s",
		progress.Evidence.State,
		progress.Phase,
		progress.ExpectedHeadSHA,
		progress.StartedAt.Format(time.RFC3339),
		progress.Deadline.Format(time.RFC3339),
		progress.Evidence.Kind,
		progress.RetryStatus,
	)
}

// printReviewIssueReport prints the deterministic stdout contract for review
// Runs after the terminal outcome is known.
func printReviewIssueReport(stdout io.Writer, outcome string, roundsCompleted int, reviewIssuesKnown bool, report reviewIssueReport) {
	if !reviewIssuesKnown {
		fmt.Fprintln(stdout, "Review Issues: unknown — fetch did not complete.")
		return
	}
	runCounts := reviewIssueReportCounts{}
	for index, current := range report.runIssues {
		status := reviewIssueDisplayStatus(current.Status)
		runCounts.add(status)
		fmt.Fprintf(stdout, "issue %03d %s — %s%s\n", index+1, status, strings.TrimSpace(current.Title), reviewIssueReasonSuffix(current, status))
	}
	fmt.Fprintf(stdout, "This Run (%s after %d Round(s)): %d resolved, %d invalid, %d duplicated, %d failed, %d unresolved.\n",
		reviewIssueOutcomeDisplay(outcome), roundsCompleted, runCounts.resolved, runCounts.invalid, runCounts.duplicated, runCounts.failed, runCounts.unresolved)
	if report.cumulativeUnavailable {
		fmt.Fprintln(stdout, "Pull Request cumulative: unavailable.")
		return
	}
	cumulativeCounts := countReviewIssueReportStatuses(report.cumulativeIssues)
	fmt.Fprintf(stdout, "Pull Request cumulative: %d resolved, %d invalid, %d duplicated, %d failed, %d unresolved.\n",
		cumulativeCounts.resolved, cumulativeCounts.invalid, cumulativeCounts.duplicated, cumulativeCounts.failed, cumulativeCounts.unresolved)
}

func (counts *reviewIssueReportCounts) add(status string) {
	switch status {
	case rounds.StatusResolved:
		counts.resolved++
	case rounds.StatusInvalid:
		counts.invalid++
	case rounds.StatusDuplicated:
		counts.duplicated++
	case rounds.StatusFailed:
		counts.failed++
	case "unresolved":
		counts.unresolved++
	}
}

func countReviewIssueReportStatuses(issues []rounds.Issue) reviewIssueReportCounts {
	counts := reviewIssueReportCounts{}
	for _, issue := range issues {
		counts.add(reviewIssueDisplayStatus(issue.Status))
	}
	return counts
}

func reviewIssueDisplayStatus(status string) string {
	switch status {
	case rounds.StatusResolved, rounds.StatusInvalid, rounds.StatusFailed, rounds.StatusDuplicated:
		return status
	default:
		return "unresolved"
	}
}

func reviewIssueOutcomeDisplay(outcome string) string {
	if outcome == store.StateCleanUnverified {
		return "Clean Unverified"
	}
	return outcome
}

func reviewIssueReasonSuffix(issue rounds.Issue, status string) string {
	switch status {
	case rounds.StatusFailed, rounds.StatusInvalid, "unresolved":
	default:
		return ""
	}
	reason := strings.Join(strings.Fields(issue.TerminalReason), " ")
	if reason == "" {
		return ""
	}
	return " — reason: " + reason
}

func reviewIssueReportData(ctx context.Context, req commandRequest, preflightResult preflight.Result, runIssues []rounds.Issue, stderr io.Writer) reviewIssueReport {
	refreshedRunIssues := refreshReviewIssueReportIssues(runIssues)
	cumulativeIssues, err := loadReviewIssueReportIssues(ctx, req, preflightResult)
	if err != nil {
		if stderr != nil {
			fmt.Fprintf(stderr, "%s: Pull Request cumulative Review Issue report unavailable: %v\n", app.Name, err)
		}
		return reviewIssueReport{runIssues: refreshedRunIssues, cumulativeUnavailable: true}
	}
	return reviewIssueReport{runIssues: refreshedRunIssues, cumulativeIssues: cumulativeIssues}
}

func refreshReviewIssueReportIssues(issues []rounds.Issue) []rounds.Issue {
	refreshed := make([]rounds.Issue, 0, len(issues))
	for _, listed := range issues {
		current := listed
		if parsed, err := rounds.ParseIssue(listed.Path); err == nil {
			current = parsed
			if strings.TrimSpace(current.Title) == "" {
				current.Title = listed.Title
			}
		}
		refreshed = append(refreshed, current)
	}
	return refreshed
}

func loadReviewIssueReportIssues(ctx context.Context, req commandRequest, preflightResult preflight.Result) ([]rounds.Issue, error) {
	root := req.reviewRoot
	if strings.TrimSpace(root) == "" {
		root = filepath.Join(req.artifactDir, "reviews", "pr-"+preflightResult.PullRequest.Number)
	}
	roundDirs, err := filepath.Glob(filepath.Join(root, "round-*"))
	if err != nil {
		return nil, fmt.Errorf("find Round artifacts in %q: %w", root, err)
	}
	sort.Strings(roundDirs)
	issues := []rounds.Issue{}
	for _, roundDir := range roundDirs {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		issuePaths, err := filepath.Glob(filepath.Join(roundDir, "issue_*.md"))
		if err != nil {
			return nil, fmt.Errorf("find Review Issue artifacts in %q: %w", roundDir, err)
		}
		sort.Strings(issuePaths)
		for _, issuePath := range issuePaths {
			issue, err := rounds.ParseIssue(issuePath)
			if err != nil {
				return nil, err
			}
			if issue.PRNumber != preflightResult.PullRequest.Number ||
				issue.HeadRepository != preflightResult.PullRequest.HeadRepository ||
				issue.HeadBranch != preflightResult.PullRequest.HeadBranch {
				continue
			}
			issues = append(issues, issue)
		}
	}
	return issues, nil
}

func prepareResolveBatch(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result) (resolveBatchPlan, error) {
	roundNumber, err := resolveRoundNumber(req.round)
	if err != nil {
		return resolveBatchPlan{}, err
	}
	selection, err := rounds.SelectCompatibleIssues(ctx, rounds.SelectRequest{
		ArtifactDir:    req.artifactDir,
		ReviewRoot:     req.reviewRoot,
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
	runtime, err := runtimeForAgentWork(req, loaded.Config)
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

func executeResolveCycle(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, runID string, session agent.SessionRef, resolvePlan resolveBatchPlan, collaborators engineCollaborators, runStore *store.Store, ui *runUI, writeGuard daemon.WriteBoundaryGuard) (resolveBatchResult, error) {
	fmt.Fprintf(ui.progress, "%s: resolve selected %d downloaded Unresolved Review Issue(s) from %d Compatible Artifact Round(s), assigned %d newest occurrence(s) into %d Batch(es), and associated %d older duplicate occurrence(s).\n", app.Name, len(resolvePlan.selection.Issues), len(resolvePlan.selection.Rounds), countBatchIssues(resolvePlan.plan.Batches), len(resolvePlan.plan.Batches), len(resolvePlan.plan.Duplicates))
	fmt.Fprintf(ui.progress, "Run: %s\n", runID)
	fmt.Fprintf(ui.progress, "User checkout: %s on branch %s\n", preflightResult.Git.Root, preflightResult.Git.Branch)
	fmt.Fprintf(ui.progress, "Artifact Directory: %s\n", req.artifactDir)
	fmt.Fprintf(ui.progress, "Round scope: %s\n", formatRoundScope(resolvePlan.roundNumber))
	fmt.Fprintf(ui.progress, "Open Pull Request: #%s %s %s\n", preflightResult.PullRequest.Number, preflightResult.PullRequest.HeadRepository, preflightResult.PullRequest.HeadBranch)
	fmt.Fprintf(ui.progress, "Agent: %s\n", resolvePlan.runtime.DisplayName)
	fmt.Fprintf(ui.progress, "Agent Model: %s\n", resolvePlan.runtime.Model)
	fmt.Fprintf(ui.progress, "Default Reasoning Effort: %s\n", displayReasoningEffort(resolvePlan.runtime.ReasoningEffort))

	engine, err := newResolveEngine(collaborators, runStore, ui, writeGuard)
	if err != nil {
		return resolveBatchResult{}, err
	}

	result, cycleErr := engine.ResolveCycle(ctx, cyclePlanFrom(req, loaded, preflightResult, runID, preflightResult.Git.Root, session, resolvePlan))
	commitCreated := false
	for _, batch := range result.Batches {
		if batch.Committed {
			commitCreated = true
			break
		}
	}
	if cycleErr != nil {
		return resolveBatchResult{Remaining: result.Remaining, CommitCreated: commitCreated}, cycleErr
	}
	if result.Remaining > 0 {
		fmt.Fprintf(ui.progress, "Final Push blocked: %d Unresolved Review Issue(s) remain.\n", result.Remaining)
		publishPushDecision(ctx, ui.sink, runID, "blocked", fmt.Sprintf("Final Push blocked: %d Unresolved Review Issue(s) remain.", result.Remaining), result.Remaining)
		return resolveBatchResult{Remaining: result.Remaining, CommitCreated: commitCreated}, nil
	}
	reviewCommit := reviewsource.ArtifactCommit{}
	if req.name != "watch" || !req.untilClean {
		reviewCommit, err = maybeCommitReviewArtifacts(ctx, req, loaded, preflightResult, collaborators.committer, ui.sink, runID, resolvePlan.roundNumber, ui.progress, writeGuard)
		if err != nil {
			return resolveBatchResult{}, err
		}
	}
	pushed, err := maybeRunFinalPush(ctx, engine, ui.sink, runID, loaded, preflightResult, preflightResult.Git.Root, commitCreated || reviewCommit.CommitSHA != "", ui.progress)
	if err != nil {
		return resolveBatchResult{}, err
	}
	if req.name == "resolve" && pushed {
		if err := maybeRequestReview(ctx, req, loaded, preflightResult, runID, ui.sink); err != nil {
			return resolveBatchResult{}, err
		}
	}
	return resolveBatchResult{Remaining: result.Remaining, CommitCreated: commitCreated}, nil
}

func newResolveEngine(collaborators engineCollaborators, runStore *store.Store, ui *runUI, writeGuard daemon.WriteBoundaryGuard) (*daemon.Engine, error) {
	return daemon.NewEngine(daemon.Dependencies{
		Runner:     collaborators.runner,
		Verifier:   collaborators.verifier,
		Committer:  collaborators.committer,
		Pusher:     collaborators.pusher,
		Source:     collaborators.source,
		Runs:       runStore,
		WriteGuard: writeGuard,
		Worktree:   collaborators.worktree,
		Sink:       ui.sink,
		Progress:   ui.progress,
	})
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

func cyclePlanFrom(req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, runID string, gitRoot string, session agent.SessionRef, resolvePlan resolveBatchPlan) daemon.CyclePlan {
	return daemon.CyclePlan{
		RunID:           runID,
		Session:         session,
		GitRoot:         gitRoot,
		ArtifactDir:     req.artifactDir,
		ReviewRoot:      req.reviewRoot,
		AgentLogs:       loaded.Config.Logs.Agent,
		SourceName:      req.source,
		AgentName:       req.agent,
		Runtime:         resolvePlan.runtime,
		AgentSelections: resolvePlan.agentSelections,
		RuntimeFactory:  resolvePlan.runtimeFactory,
		Verification:    loaded.Config.Defaults.Verification,
		AutoCommit:      loaded.Config.Defaults.AutoCommit,
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

func runWatchCommand(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, notifier roundnotify.Notifier, stdout, stderr io.Writer) int {
	collaborators := commandDependenciesForContext(ctx).newEngineCollaborators()
	categories := reviewProfileCategories()
	profilePreflight, err := runProfileOperationalPreflight(ctx, req, loaded.Config, categories, preflightResult.Git.Root, collaborators.runner, stderr)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	runtime, err := runtimeForOperationalProfileRun(req, loaded.Config, categories, profilePreflight.Override)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	agentSelections, err := operationalAgentSelectionProfiles(loaded.Config, categories, profilePreflight.Override)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	runtimeFactory := operationalRuntimeFactory(req)
	req = requestWithRuntimeSelection(req, runtime)

	runStore, err := store.Open(ctx, loaded.HomeDir)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	defer func() {
		_ = runStore.Close()
	}()
	sweepRunRetention(ctx, runStore, req.artifactDir, loaded.Config.Store.JournalRetention, stderr)
	run, err := createOperationalRun(ctx, runStore, store.KindWatch, req, preflightResult, runtime, stderr)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	journalBranchIntegrityIntegrations(ctx, runStore, run.ID, req.branchIntegrity.Integrated)
	if err := publishBranchIntegrityBypassAudit(ctx, req, preflightResult, runStore, run.ID); err != nil {
		markRunFailedAndNotify(ctx, runStore, run.ID, notifier, stderr)
		printBranchIntegrityAuditFailure(req.name, run.ID, err, stderr)
		return exitPreflight
	}
	if err := ensureReviewRunArtifactDirectory(req, loaded, preflightResult); err != nil {
		markRunFailedAndNotify(ctx, runStore, run.ID, notifier, stderr)
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}
	if err := req.reportDetachedRunCreated(run.ID); err != nil {
		markRunFailedAndNotify(ctx, runStore, run.ID, notifier, stderr)
		printRunFailure(req.name, err, stderr)
		return exitRunFailed
	}
	session := agent.SessionRefForRun(run.ID, preflightResult.Git.Root)
	sessionForClose := session
	if len(agentSelections) > 0 {
		sessionForClose = agent.SessionRef{}
	}
	if err := rememberInteractiveDefaults(ctx, runStore, req); err != nil {
		closeAgentSession(ctx, collaborators.runner, runtime, sessionForClose, run.ID, runStore)
		markRunFailedAndNotify(ctx, runStore, run.ID, notifier, stderr)
		printWatchRunFailure(err, stderr)
		return exitRunFailed
	}

	fmt.Fprintf(stderr, "Watch Run: %s\n", run.ID)
	fmt.Fprintf(stderr, "User checkout: %s\n", preflightResult.Git.Root)
	fmt.Fprintf(stderr, "Open Pull Request: #%s %s %s\n", preflightResult.PullRequest.Number, preflightResult.PullRequest.HeadRepository, preflightResult.PullRequest.HeadBranch)
	fmt.Fprintf(stderr, "Review Source: %s\n", req.source)
	fmt.Fprintf(stderr, "Agent: %s\n", runtime.DisplayName)
	fmt.Fprintf(stderr, "Agent Model: %s\n", runtime.Model)
	fmt.Fprintf(stderr, "Default Reasoning Effort: %s\n", displayReasoningEffort(runtime.ReasoningEffort))
	fmt.Fprintf(stderr, "Max Rounds: %d\n", req.maxRounds)

	// One cockpit for the entire Watch Run, across all Rounds and Batches.
	cockpitView := buildLiveRunView(req, loaded, preflightResult, run.ID, "WaitingForReview", nil, nil, stderr)
	cockpitView.WorkDir = preflightResult.Git.Root
	if !liveTUIEnabled(stderr) {
		plainView := buildLiveRunView(req, loaded, preflightResult, run.ID, "WaitingForReview", nil, []string{"Waiting for Review Source status..."}, stderr)
		plainView.WorkDir = preflightResult.Git.Root
		fmt.Fprint(stderr, roundtui.RenderLiveRunView(plainView))
	}
	ui, err := startRunUI(ctx, cockpitView, run.ID, loaded.HomeDir, runStore, stderr, req.noAgentConsole)
	if err != nil {
		closeAgentSession(ctx, collaborators.runner, runtime, sessionForClose, run.ID, runStore)
		markRunFailedAndNotify(ctx, runStore, run.ID, notifier, stderr)
		printWatchRunFailure(err, stderr)
		return exitRunFailed
	}
	defer ui.Close()

	watchReportIssues := []rounds.Issue{}
	var requester reviewsource.ReviewRequester
	if loaded.Config.ReviewSource.RequestReview {
		requester = commandDependenciesForContext(ctx).newReviewRequester(store.JournalSink{Store: runStore})
	}
	writeGuard := reviewRunTargetGuard(ctx, run)
	result, err := watch.Run(ctx, watch.Request{
		RunID:            run.ID,
		PRNumber:         preflightResult.PullRequest.Number,
		BaseRepository:   preflightResult.PullRequest.BaseRepository,
		HeadSHA:          preflightResult.Git.HEAD,
		RequestReview:    loaded.Config.ReviewSource.RequestReview,
		ReviewCommand:    loaded.Config.ReviewSource.RequestCommand,
		UntilClean:       req.untilClean,
		MaxRounds:        req.maxRounds,
		PollInterval:     loaded.Config.Watch.PollInterval,
		QuietPeriod:      loaded.Config.Watch.QuietPeriod,
		ReviewTimeout:    loaded.Config.Watch.ReviewTimeout,
		CheckGracePeriod: loaded.Config.Watch.CheckGracePeriod,
		BudgetEnabled:    loaded.Config.Budget.Enabled,
		MaxRunDuration:   loaded.Config.Budget.MaxRunDuration,
	}, watch.Dependencies{
		StopRequests:    runStore,
		ReviewRequester: requester,
		WriteGuard:      writeGuard,
		ReviewEvidence: watch.ReviewEvidenceFunc(func(ctx context.Context, evidenceReq watch.ReviewEvidenceRequest) (reviewsource.Evidence, error) {
			return commandDependenciesForContext(ctx).watchReviewEvidence(ctx, reviewsource.EvidenceRequest{
				Source:          req.source,
				PRNumber:        evidenceReq.PRNumber,
				BaseRepository:  preflightResult.PullRequest.BaseRepository,
				HeadRepository:  preflightResult.PullRequest.HeadRepository,
				HeadBranch:      preflightResult.PullRequest.HeadBranch,
				ExpectedHeadSHA: evidenceReq.ExpectedHeadSHA,
			})
		}),
		Artifacts: watch.ArtifactPublishFunc(func(ctx context.Context, artifactReq watch.ArtifactPublishRequest) (watch.ArtifactPublication, error) {
			commit, err := maybeCommitReviewArtifacts(
				ctx,
				req,
				loaded,
				preflightResult,
				collaborators.committer,
				ui.sink,
				run.ID,
				artifactReq.Round,
				ui.progress,
				writeGuard,
			)
			if err != nil || commit.CommitSHA == "" {
				return watch.ArtifactPublication{Commit: commit}, err
			}
			engine, err := newResolveEngine(collaborators, runStore, ui, writeGuard)
			if err != nil {
				return watch.ArtifactPublication{}, err
			}
			if _, err := maybeRunFinalPush(ctx, engine, ui.sink, run.ID, loaded, preflightResult, preflightResult.Git.Root, true, ui.progress); err != nil {
				return watch.ArtifactPublication{}, err
			}
			evidence, inherited, err := inheritReviewArtifactEvidence(ctx, reviewArtifactEvidenceRequest{
				Source:         req.source,
				PRNumber:       preflightResult.PullRequest.Number,
				BaseRepository: preflightResult.PullRequest.BaseRepository,
				HeadRepository: preflightResult.PullRequest.HeadRepository,
				HeadBranch:     preflightResult.PullRequest.HeadBranch,
				GitRoot:        preflightResult.Git.Root,
				Commit:         commit,
				ParentEvidence: artifactReq.ParentEvidence,
				ParentHeadSHA:  artifactReq.ParentHeadSHA,
			})
			if err != nil {
				fmt.Fprintf(ui.progress, "Review artifact Evidence inheritance unavailable; falling back to head polling: %v\n", err)
				inherited = false
			}
			if !inherited {
				evidence = reviewsource.Evidence{}
			}
			return watch.ArtifactPublication{Commit: commit, Evidence: evidence}, nil
		}),
		Fetcher: watch.FetchFunc(func(ctx context.Context, _ int) (watch.FetchResult, error) {
			fetchResult, issues, err := fetchWatchRound(ctx, req, loaded, preflightResult, ui.progress)
			if err == nil {
				watchReportIssues = append(watchReportIssues, issues...)
			}
			return fetchResult, err
		}),
		Resolver: watch.ResolveFunc(func(ctx context.Context) (watch.ResolveResult, error) {
			return resolveWatchBatches(ctx, req, loaded, preflightResult, runtime, agentSelections, runtimeFactory, run.ID, session, collaborators, runStore, ui, writeGuard)
		}),
		Clock:   commandDependenciesForContext(ctx).watchClock,
		Sleeper: commandDependenciesForContext(ctx).watchSleeper,
		Sink:    store.JournalSink{Store: runStore},
		Progress: func(progress watch.WaitProgress) {
			fmt.Fprintln(ui.progress, formatReviewWaitProgress(progress))
		},
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
		closeAgentSession(ctx, collaborators.runner, runtime, sessionForClose, run.ID, runStore)
		printRunFailure(req.name, completeErr, stderr)
		return exitRunFailed
	}
	terminalContext := watchTerminalCompletionContext(req, completed.Run, result)
	var cleanupWarnings []cleanupWarning
	if result.Outcome == store.StateReviewSkipped {
		publishTerminalCompletionWithContext(
			completeCtx,
			runStore,
			notifier,
			stderr,
			completed,
			terminalContext,
		)
		cleanupWarnings = bestEffortForceStopAgentSessions(completeCtx, runStore, completed.Run)
	} else {
		closeAgentSession(completeCtx, collaborators.runner, runtime, sessionForClose, completed.ID, runStore)
		publishTerminalCompletionWithContext(completeCtx, runStore, notifier, stderr, completed, terminalContext)
	}
	if result.Outcome == store.StateCheckoutMoved {
		fmt.Fprintf(stderr, "%s: %s Next: %s\n", store.StateCheckoutMoved, result.TerminalReason, result.NextAction)
	}
	// The cockpit stays on screen, read-only, until the user closes it.
	ui.Wait()
	ui.Close()

	fmt.Fprintf(stderr, "Watch Run %s reached %s after %d Round(s).\n", completed.ID, completed.State, result.Rounds)
	if result.Outcome == store.StateReviewSkipped {
		fmt.Fprintf(stderr, "Review Skipped: %s\n", result.TerminalReason)
		fmt.Fprintf(stderr, "Next: %s\n", result.NextAction)
	}
	reportSecondaryCleanupWarnings(completeCtx, runStore, completed.ID, cleanupWarnings, stderr)
	if result.Outcome == store.StateMaxRoundsReached && result.Remaining > 0 {
		fmt.Fprintf(stderr, "MaxRoundsReached with %d Unresolved Review Issue(s) remaining.\n", result.Remaining)
		printAgentCheckoutChangesNotice(stderr)
	}
	if result.Outcome == store.StateUnresolved {
		fmt.Fprintf(stderr, "Unresolved: the last Round settled nothing; %d Unresolved Review Issue(s) remain for developer attention.\n", result.Remaining)
		printAgentCheckoutChangesNotice(stderr)
	}
	if result.Outcome == store.StateTimedOut {
		if result.TerminalReasonIsDiagnostic {
			fmt.Fprintln(stderr, result.TerminalReason)
		}
		fmt.Fprintf(stderr, "Review Source timed out. To request another CodeRabbit review manually, comment: %s\n", result.ManualReviewCommand)
	}
	if result.Outcome == store.StateCleanUnverified {
		fmt.Fprintln(stderr, "CleanUnverified: Merge-Ready was not confirmed because no accepted Review Source Evidence appeared within the grace period. Next: confirm the pull request's Review Source Evidence before merging.")
	}
	if result.Outcome == store.StateReviewSkipped {
		printReviewSkippedReport(stdout, result.TerminalReason, result.NextAction)
	} else {
		printReviewIssueReport(stdout, completed.State, result.Rounds, result.ReviewIssuesKnown, reviewIssueReportData(context.WithoutCancel(ctx), req, preflightResult, watchReportIssues, stderr))
	}
	if stopped {
		printStopSummary(ctx, req, preflightResult, stderr)
		return exitOK
	}
	if err != nil {
		printWatchRunFailure(err, stderr)
		return exitRunFailed
	}
	return exitForWatchOutcome(result.Outcome)
}

func createOperationalRun(ctx context.Context, runStore *store.Store, kind string, req commandRequest, preflightResult preflight.Result, runtime agent.RuntimeSpec, stderr io.Writer) (store.Run, error) {
	createReq := store.CreateRunRequest{
		Kind:            kind,
		HeadRepository:  preflightResult.PullRequest.HeadRepository,
		HeadBranch:      preflightResult.PullRequest.HeadBranch,
		BaseRepository:  preflightResult.PullRequest.BaseRepository,
		PRNumber:        preflightResult.PullRequest.Number,
		GitRoot:         preflightResult.Git.Root,
		LocalBranch:     preflightResult.Git.Branch,
		HeadSHA:         preflightResult.Git.HEAD,
		ArtifactDir:     req.artifactDir,
		WorkDir:         preflightResult.Git.Root,
		Agent:           runtime.ID,
		Model:           runtime.Model,
		ReasoningEffort: runtime.ReasoningEffort,
	}
	return createReviewRun(ctx, runStore, req, createReq, stderr)
}

func ensureReviewRunArtifactDirectory(req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result) error {
	_, err := roundconfig.ValidateArtifactDirectory(req.artifactDir, preflightResult.Git.Root, loaded.HomeDir)
	return err
}

func reviewRunTargetGuard(ctx context.Context, run store.Run) *watch.TargetGuard {
	return &watch.TargetGuard{
		WorkDir: run.GitRoot,
		Target: watch.Checkout{
			Branch:   run.HeadBranch,
			Revision: run.HeadSHA,
		},
		Checkouts: commandDependenciesForContext(ctx).checkoutReader,
	}
}

func createReviewRun(ctx context.Context, runStore *store.Store, req commandRequest, createReq store.CreateRunRequest, stderr io.Writer) (store.Run, error) {
	createReq.OwnerPID = os.Getpid()
	createReq.OwnerIdentity = currentOwnerIdentity(ctx)
	if req.skipBranchIntegrity && req.branchIntegrity.ActiveRun != nil {
		return runStore.CreateRunSkippingActiveLock(ctx, createReq)
	}
	return createRunReclaimingOrphan(ctx, runStore, stderr, func() (store.Run, error) {
		return runStore.CreateRun(ctx, createReq)
	})
}

func createRunReclaimingOrphan(ctx context.Context, runStore *store.Store, stderr io.Writer, create func() (store.Run, error)) (store.Run, error) {
	run, err := create()
	if err == nil {
		return run, nil
	}
	var activeErr store.ActiveRunError
	if !errors.As(err, &activeErr) {
		return store.Run{}, err
	}
	if _, ok, reclaimErr := reclaimOrphanedActiveRun(ctx, runStore, activeErr.Existing, stderr); reclaimErr != nil {
		return store.Run{}, reclaimErr
	} else if !ok {
		return store.Run{}, err
	}
	return create()
}

func reclaimOrphanedActiveRun(ctx context.Context, runStore *store.Store, active store.Run, stderr io.Writer) (store.Run, bool, error) {
	if store.IsTerminalState(active.State) {
		return active, false, nil
	}
	pid, ok := activeOwnerPID(active)
	if !ok || store.ProcessAlive(pid) {
		return active, false, nil
	}
	reason := orphanedActiveRunReason(pid)
	printForceStopAgentSessionWarnings(stderr, bestEffortForceStopAgentSessions(ctx, runStore, active))
	if err := runStore.ReclaimOrphanedRun(ctx, active, reason); err != nil {
		return active, false, fmt.Errorf("reclaim orphaned Active Run %s: %w", active.ID, err)
	}
	reclaimed, found, err := runStore.Run(ctx, active.ID)
	if err != nil {
		return active, false, fmt.Errorf("read reclaimed Active Run %s: %w", active.ID, err)
	}
	if !found || reclaimed.State != store.StateFailed {
		return active, false, nil
	}
	printOrphanedActiveRunReclaimed(stderr, reclaimed.ID, reason)
	return reclaimed, true, nil
}

func printForceStopAgentSessionWarnings(stderr io.Writer, warnings []cleanupWarning) {
	if stderr == nil {
		return
	}
	for _, warning := range warnings {
		fmt.Fprintf(stderr, "%s: %s\n", app.Name, warning.text)
	}
}

func activeOwnerPID(run store.Run) (int, bool) {
	if run.OwnerPID == nil || *run.OwnerPID <= 0 {
		return 0, false
	}
	return *run.OwnerPID, true
}

// currentOwnerIdentity returns this process's start-time identity token for
// the Run row, or "" when the platform cannot provide one. An absent token
// degrades Force Stop to the legacy PID-only owner proof.
func currentOwnerIdentity(ctx context.Context) string {
	identity, err := store.OwnerProcessIdentity(ctx, os.Getpid())
	if err != nil {
		return ""
	}
	return identity
}

func orphanedActiveRunReason(pid int) string {
	return fmt.Sprintf("owner process %d not running; lock reclaimed", pid)
}

func printOrphanedActiveRunReclaimed(stderr io.Writer, runID string, reason string) {
	if stderr == nil {
		return
	}
	fmt.Fprintf(stderr, "%s: reclaimed orphaned Active Run %s: %s\n", app.Name, runID, reason)
}

func requestWithRuntimeSelection(req commandRequest, runtime agent.RuntimeSpec) commandRequest {
	req.model = runtime.Model
	req.reasoningEffort = runtime.ReasoningEffort
	return req
}

func worktreeBootstrapSpec(config roundconfig.Config) runworktree.BootstrapSpec {
	return runworktree.BootstrapSpec{
		Command: config.Worktree.Bootstrap,
		Timeout: config.Worktree.BootstrapTimeout,
	}
}

func newBootstrapOutputWriter(ctx context.Context, runID string, runStore *store.Store, stderr io.Writer) io.Writer {
	return &bootstrapRunWriter{
		ctx:    ctx,
		runID:  runID,
		stderr: stderr,
		sink:   store.JournalSink{Store: runStore},
		mu:     &sync.Mutex{},
	}
}

type bootstrapRunWriter struct {
	ctx    context.Context
	runID  string
	stderr io.Writer
	sink   runevent.Sink
	mu     *sync.Mutex
}

func (writer *bootstrapRunWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if writer.mu != nil {
		writer.mu.Lock()
		defer writer.mu.Unlock()
	}
	if writer.stderr != nil {
		n, err := writer.stderr.Write(p)
		if err != nil {
			return n, err
		}
		if n != len(p) {
			return n, io.ErrShortWrite
		}
	}
	if writer.sink == nil {
		return len(p), nil
	}
	payload, err := json.Marshal(map[string]string{"output": string(p)})
	if err != nil {
		return len(p), nil
	}
	ctx := writer.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	if err := writer.sink.Publish(context.WithoutCancel(ctx), runevent.RunEvent{
		RunID:   writer.runID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonStatus,
		Summary: runevent.BoundSummary(string(p)),
		Time:    time.Now().UTC(),
		Payload: payload,
	}); err != nil {
		return len(p), err
	}
	return len(p), nil
}

func fetchWatchRound(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, stderr io.Writer) (watch.FetchResult, []rounds.Issue, error) {
	items, err := commandDependenciesForContext(ctx).fetchReviewItems(ctx, reviewsource.FetchRequest{
		Source:          req.source,
		PRNumber:        preflightResult.PullRequest.Number,
		BaseRepository:  preflightResult.PullRequest.BaseRepository,
		HeadRepository:  preflightResult.PullRequest.HeadRepository,
		HeadBranch:      preflightResult.PullRequest.HeadBranch,
		HeadSHA:         preflightResult.Git.HEAD,
		IncludeNitpicks: loaded.Config.ReviewSource.IncludeNitpicks,
	})
	if err != nil {
		return watch.FetchResult{}, nil, err
	}
	roundResult, err := rounds.PersistRound(ctx, rounds.PersistRequest{
		ArtifactDir:    req.artifactDir,
		ReviewRoot:     req.reviewRoot,
		Source:         req.source,
		PRNumber:       preflightResult.PullRequest.Number,
		HeadRepository: preflightResult.PullRequest.HeadRepository,
		HeadBranch:     preflightResult.PullRequest.HeadBranch,
		HeadSHA:        preflightResult.Git.HEAD,
		Items:          items,
	})
	if err != nil {
		return watch.FetchResult{}, nil, err
	}
	if roundResult.Reused {
		fmt.Fprintf(stderr, "Reused Round %03d with %d Review Issue(s).\n", roundResult.Round, len(roundResult.IssuePaths))
	} else {
		fmt.Fprintf(stderr, "Fetched Round %03d with %d Review Issue(s).\n", roundResult.Round, len(roundResult.IssuePaths))
	}
	reportIssues := make([]rounds.Issue, 0, len(roundResult.IssuePaths))
	for _, path := range roundResult.IssuePaths {
		issue, err := rounds.ParseIssue(path)
		if err != nil {
			return watch.FetchResult{}, nil, err
		}
		reportIssues = append(reportIssues, issue)
	}
	return watch.FetchResult{Round: roundResult.Round, Issues: len(roundResult.IssuePaths)}, reportIssues, nil
}

// resolveWatchBatches runs exactly one resolve cycle for the current
// Round. Failed Review Issues are not retried inside the same Round: the
// next fetched Round re-downloads their still-open Review Source threads
// as fresh occurrences. Progress means the cycle settled at least one
// selected issue, so the watch loop can stop Rounds that change nothing.
func resolveWatchBatches(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, runtime agent.RuntimeSpec, agentSelections daemon.AgentSelectionProfiles, runtimeFactory daemon.AgentRuntimeFactory, runID string, session agent.SessionRef, collaborators engineCollaborators, runStore *store.Store, ui *runUI, writeGuard daemon.WriteBoundaryGuard) (watch.ResolveResult, error) {
	resolvePlan, err := prepareResolveBatch(ctx, req, loaded, preflightResult)
	if err != nil {
		var noArtifacts rounds.NoCompatibleArtifactsError
		if errors.As(err, &noArtifacts) {
			return watch.ResolveResult{Remaining: 0, Progress: false}, nil
		}
		return watch.ResolveResult{}, err
	}
	resolvePlan.runtime = runtime
	resolvePlan.agentSelections = agentSelections
	resolvePlan.runtimeFactory = runtimeFactory
	result, err := executeResolveCycle(ctx, req, loaded, preflightResult, runID, session, resolvePlan, collaborators, runStore, ui, writeGuard)
	if err != nil {
		return watch.ResolveResult{}, err
	}
	headSHA := ""
	if req.untilClean && result.Remaining == 0 {
		headSHA, err = commandDependenciesForContext(ctx).watchHeadSHA(ctx, preflightResult.Git.Root)
		if err != nil {
			return watch.ResolveResult{}, err
		}
	}
	return watch.ResolveResult{
		Remaining: result.Remaining,
		Progress:  result.Remaining < len(resolvePlan.selection.Issues),
		HeadSHA:   headSHA,
	}, nil
}

func exitForWatchOutcome(outcome string) int {
	switch outcome {
	case store.StateClean, store.StateMaxRoundsReached, store.StateStopped:
		return exitOK
	case store.StateCleanUnverified:
		return exitUnverified
	case store.StateReviewSkipped:
		return exitUnverified
	case store.StateBudgetExceeded, store.StateTimedOut, store.StateFailed, store.StateUnresolved:
		return exitRunFailed
	default:
		return exitRunFailed
	}
}

func printStopSummary(ctx context.Context, req commandRequest, preflightResult preflight.Result, stderr io.Writer) {
	fmt.Fprintln(stderr, "Stop Request preserved local work and stopped before any later verification, commit, push, fetch, or Review Source mutation.")
	changes, err := commandDependenciesForContext(ctx).inspectChangedPaths(ctx, preflightResult.Git.Root)
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
	return agent.IsStopError(err) || errors.Is(err, daemon.ErrStopRequested) || errors.Is(err, watch.ErrStopRequested) || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || (ctx != nil && ctx.Err() != nil)
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
	fmt.Fprint(stderr, roundtui.RenderLiveRunView(buildLiveRunView(req, loaded, preflightResult, runID, pipelineState, issues, console, stderr)))
}

func buildLiveRunView(req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, runID string, pipelineState string, issues []rounds.Issue, console []string, output io.Writer) roundtui.LiveRunView {
	return roundtui.LiveRunView{
		Command:         req.name,
		Repository:      preflightResult.PullRequest.HeadRepository,
		PRNumber:        preflightResult.PullRequest.Number,
		HeadBranch:      preflightResult.PullRequest.HeadBranch,
		ReviewSource:    displayReviewSource(req.source),
		Agent:           displayAgent(req.agent),
		Model:           req.model,
		ReasoningEffort: req.reasoningEffort,
		HEAD:            preflightResult.Git.HEAD,
		RunID:           runID,
		PipelineState:   pipelineState,
		BudgetState:     formatBudgetState(loaded.Config),
		GitState:        formatGitState(preflightResult.Git),
		CurrentRound:    0,
		MaxRounds:       req.maxRounds,
		AutoCommit:      loaded.Config.Defaults.AutoCommit,
		AutoPush:        preflightResult.PushPlan.Enabled,
		LastPush:        lastPushState(preflightResult.PushPlan),
		Issues:          issues,
		Console:         console,
		Width:           liveViewWidth(output),
	}
}

func liveTUIEnabled(output io.Writer) bool {
	environment, output := environmentForWriter(output)
	switch strings.ToLower(strings.TrimSpace(environment.tuiMode)) {
	case "always", "1", "true", "yes", "on":
		return true
	case "never", "0", "false", "no", "off":
		return false
	}
	if strings.EqualFold(environment.term, "dumb") {
		return false
	}
	file, ok := output.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func liveViewWidth(output io.Writer) int {
	environment, _ := environmentForWriter(output)
	width, err := strconv.Atoi(strings.TrimSpace(environment.columns))
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
		name:            name,
		arguments:       append([]string(nil), args...),
		source:          config.ReviewSource.Name,
		agent:           config.Defaults.Agent,
		round:           "all",
		untilClean:      config.Watch.UntilClean,
		maxRounds:       config.Watch.MaxRounds,
		artifactDir:     config.Defaults.ArtifactDir,
		agentFullAccess: config.Defaults.AgentFullAccess,
	}
	if name == "fetch" {
		req.round = "auto"
	}

	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&req.pr, "pr", "", "Open Pull Request number")
	fs.StringVar(&req.spec, "spec", "", "Spec slug under docs/specs/")
	fs.BoolVar(&req.noInput, "no-input", false, "Fail instead of opening Interactive Input")
	fs.BoolVar(&req.interactive, "interactive", false, "Open Interactive Input before starting")
	fs.StringVar(&req.artifactDir, "artifact-dir", req.artifactDir, "Artifact Directory")
	fs.StringVar(&req.baseRepo, "base-repo", "", "Base repository, owner/name")
	fs.StringVar(&req.headRepo, "head-repo", "", "Head Repository, owner/name")
	fs.StringVar(&req.headBranch, "head-branch", "", "PR Head Branch")
	fs.BoolVar(&req.skipBranchIntegrity, "skip-branch-integrity", false, "Skip pending Run Branch and Active Run guardrails after publishing a PR audit comment")

	switch name {
	case "fetch":
		fs.StringVar(&req.source, "source", req.source, "Review Source")
		fs.StringVar(&req.round, "round", "auto", "Round number or auto")
	case "resolve":
		fs.BoolVar(&req.detach, "detach", false, "Start a Detached Run and print attach/stop commands")
		fs.StringVar(&req.agent, "agent", req.agent, "Agent runtime")
		fs.StringVar(&req.model, "model", req.model, "Agent model override")
		fs.StringVar(&req.reasoningEffort, "reasoning-effort", req.reasoningEffort, "Default reasoning effort override")
		fs.StringVar(&req.agentCmd, "agent-command", "", "Agent command override")
		fs.BoolVar(&req.agentFullAccess, "agent-full-access", req.agentFullAccess, "Opt into Agent runtime full-access mode")
		fs.BoolVar(&req.noAgentConsole, "no-agent-console", false, "Hide Agent-source console events from non-TTY stderr")
		fs.StringVar(&req.round, "round", "all", "Round number or all")
	case "watch":
		fs.BoolVar(&req.detach, "detach", false, "Start a Detached Run and print attach/stop commands")
		fs.StringVar(&req.source, "source", req.source, "Review Source")
		fs.StringVar(&req.agent, "agent", req.agent, "Agent runtime")
		fs.StringVar(&req.model, "model", req.model, "Agent model override")
		fs.StringVar(&req.reasoningEffort, "reasoning-effort", req.reasoningEffort, "Default reasoning effort override")
		fs.StringVar(&req.agentCmd, "agent-command", "", "Agent command override")
		fs.BoolVar(&req.agentFullAccess, "agent-full-access", req.agentFullAccess, "Opt into Agent runtime full-access mode")
		fs.BoolVar(&req.noAgentConsole, "no-agent-console", false, "Hide Agent-source console events from non-TTY stderr")
		fs.BoolVar(&req.untilClean, "until-clean", req.untilClean, "Repeat until no Unresolved Review Issues remain and accepted Review Source Evidence confirms the pushed head")
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
	recordSelectionFlagPresence(fs, &req)
	if err := validateExplicitSelectionFlags(req); err != nil {
		return req, err
	}

	return req, nil
}

func recordSelectionFlagPresence(fs *flag.FlagSet, req *commandRequest) {
	fs.Visit(func(flag *flag.Flag) {
		switch flag.Name {
		case "agent":
			req.agentSet = true
		case "model":
			req.modelSet = true
		case "reasoning-effort":
			req.reasoningEffortSet = true
		}
	})
}

func validateExplicitSelectionFlags(req commandRequest) error {
	presence := invocationSelectionFlagPresence{
		agent:           req.agentSet,
		model:           req.modelSet,
		reasoningEffort: req.reasoningEffortSet,
	}
	if err := presence.validate(); err != nil {
		return err
	}
	if presence.empty() {
		return nil
	}
	if strings.TrimSpace(req.agent) == "" {
		return validationError{message: "--agent cannot be empty in a complete one-Run Agent Selection override"}
	}
	if req.modelSet && strings.TrimSpace(req.model) == "" {
		return validationError{message: "--model cannot be empty in a complete one-Run Agent Selection override"}
	}
	return nil
}

type invocationSelectionFlagPresence struct {
	agent           bool
	model           bool
	reasoningEffort bool
}

func selectionFlagPresence(args []string) invocationSelectionFlagPresence {
	var presence invocationSelectionFlagPresence
	for index := 0; index < len(args); index++ {
		arg := args[index]
		if arg == "--" {
			break
		}
		if !strings.HasPrefix(arg, "-") {
			continue
		}
		name, _, hasValue := strings.Cut(strings.TrimLeft(arg, "-"), "=")
		switch name {
		case "agent":
			presence.agent = true
		case "model":
			presence.model = true
		case "reasoning-effort":
			presence.reasoningEffort = true
		}
		if !hasValue && selectionFlagConsumesValue(name) && index+1 < len(args) {
			index++
		}
	}
	return presence
}

func selectionFlagConsumesValue(name string) bool {
	switch name {
	case "agent-full-access", "detach", "interactive", "no-agent-console", "no-input", "qa", "skip-branch-integrity", "until-clean":
		return false
	default:
		return true
	}
}

func validateSelectionOverrideArgs(args []string) error {
	return selectionFlagPresence(args).validate()
}

func (presence invocationSelectionFlagPresence) empty() bool {
	return !presence.agent && !presence.model && !presence.reasoningEffort
}

func (presence invocationSelectionFlagPresence) validate() error {
	if presence.empty() || (presence.agent && presence.model && presence.reasoningEffort) {
		return nil
	}
	return validationError{message: "--agent, --model, and --reasoning-effort must be provided together for a one-Run Agent Selection override; omit all three to use Agent Selection Profiles"}
}

func displayReasoningEffort(reasoningEffort string) string {
	reasoningEffort = strings.TrimSpace(reasoningEffort)
	if reasoningEffort == "" {
		return "model-managed"
	}
	return reasoningEffort
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

func validateAgentConsoleDisplay(req commandRequest, stderr io.Writer) error {
	if req.noAgentConsole && liveTUIEnabled(stderr) {
		if req.detach || req.detachChild != nil {
			return nil
		}
		return validationError{message: "--no-agent-console cannot be used with the interactive cockpit"}
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

type reviewSpecRequest struct {
	ExplicitSlug string
	RepoRoot     string
	SpecsRoot    string
	HeadSHA      string
}

func specsRootForWorkDir(resolved roundconfig.SpecsRoot, repoRoot string, workDir string) string {
	if resolved.External {
		return resolved.Path
	}
	if strings.TrimSpace(workDir) == "" {
		return resolved.Path
	}
	rel, err := filepath.Rel(repoRoot, resolved.Path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return resolved.Path
	}
	return filepath.Join(workDir, rel)
}

func specsRootIsDefault(repoRoot string, resolved roundconfig.SpecsRoot) bool {
	return filepath.Clean(resolved.Path) == filepath.Join(filepath.Clean(repoRoot), "docs", "specs")
}

func reportNonDefaultSpecsRoot(stderr io.Writer, repoRoot string, resolved roundconfig.SpecsRoot) {
	if specsRootIsDefault(repoRoot, resolved) {
		return
	}
	fmt.Fprintf(stderr, "Spec Root: %s\n", resolved.Path)
}

func reviewArtifactSpecsRoot(loaded roundconfig.Loaded, repoRoot string) roundconfig.SpecsRoot {
	return roundconfig.SpecsRoot{Path: filepath.Join(repoRoot, loaded.Config.Specs.Root)}
}

func reviewArtifactUsesDefaultSpecsRoot(configured string) bool {
	return strings.TrimSpace(configured) == "docs/specs"
}

func resolveReviewArtifactRoot(ctx context.Context, req commandRequest, preflightResult preflight.Result, specsRoot roundconfig.SpecsRoot) (string, error) {
	prNumber, err := strconv.Atoi(req.pr)
	if err != nil {
		return "", fmt.Errorf("parse Open Pull Request number %q: %w", req.pr, err)
	}
	explicitArtifactDir := ""
	if req.explicitArtifactDir {
		explicitArtifactDir = req.artifactDir
	}
	specSlug := ""
	if explicitArtifactDir == "" {
		specSlug, err = reviewSpecSlug(ctx, reviewSpecRequest{
			ExplicitSlug: req.spec,
			RepoRoot:     preflightResult.Git.Root,
			SpecsRoot:    specsRoot.Path,
			HeadSHA:      preflightResult.Git.HEAD,
		}, commandDependenciesForContext(ctx).reviewSpecGitRunner)
		if err != nil {
			return "", err
		}
	}
	return roundconfig.ResolveReviewRoot(roundconfig.ReviewArtifactContext{
		ExplicitArtifactDir: explicitArtifactDir,
		RepoRoot:            preflightResult.Git.Root,
		SpecsRoot:           specsRoot.Path,
		SpecSlug:            specSlug,
		PRNumber:            prNumber,
	})
}

func reviewSpecSlug(ctx context.Context, req reviewSpecRequest, runner preflight.GitRunner) (string, error) {
	specsRoot := strings.TrimSpace(req.SpecsRoot)
	if specsRoot == "" {
		specsRoot = filepath.Join(req.RepoRoot, "docs", "specs")
	}
	if slug := strings.TrimSpace(req.ExplicitSlug); slug != "" {
		if reviewSpecFolderExists(specsRoot, slug) {
			return slug, nil
		}
		return "", nil
	}
	message, err := runner.RunGit(ctx, req.RepoRoot, "log", "-1", "--format=%B", req.HeadSHA)
	if err != nil {
		return "", fmt.Errorf("read PR head commit trailers: %w", err)
	}
	slug := newestRoundfixSpecTrailer(message)
	if slug == "" || !reviewSpecFolderExists(specsRoot, slug) {
		return "", nil
	}
	return slug, nil
}

func newestRoundfixSpecTrailer(message string) string {
	slug := ""
	for _, line := range strings.Split(message, "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok || strings.TrimSpace(key) != "Roundfix-Spec" {
			continue
		}
		if candidate := strings.TrimSpace(value); candidate != "" {
			slug = candidate
		}
	}
	return slug
}

func reviewSpecFolderExists(specsRoot string, slug string) bool {
	if !validReviewSpecSlug(slug) {
		return false
	}
	info, err := os.Stat(filepath.Join(specsRoot, slug))
	return err == nil && info.IsDir()
}

func validReviewSpecSlug(slug string) bool {
	if slug == "" || slug == "." || slug == ".." || filepath.IsAbs(slug) {
		return false
	}
	return !strings.ContainsAny(slug, `/\`)
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
		RequestReview:          loaded.Config.ReviewSource.RequestReview,
	})
}

func defaultFetchReviewItems(ctx context.Context, req reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
	if req.Source != reviewsource.SourceCodeRabbit {
		return nil, fmt.Errorf("unsupported Review Source %q; supported value: coderabbit", req.Source)
	}
	return coderabbit.Client{}.FetchReviews(ctx, req)
}

func defaultCommentOnPullRequest(ctx context.Context, source string, repository string, prNumber int, body string) error {
	if source != reviewsource.SourceCodeRabbit {
		return fmt.Errorf("unsupported Review Source %q; supported value: coderabbit", source)
	}
	return coderabbit.GHClient{}.CommentOnPullRequest(ctx, repository, prNumber, body)
}

type engineCollaborators struct {
	runner       agent.Runner
	verifier     daemon.Verifier
	committer    daemon.Committer
	pusher       daemon.Pusher
	source       daemon.ReviewSourceResolver
	worktree     daemon.WorktreeSnapshotter
	priorChanges daemon.PriorChangedResolver
}

func defaultEngineCollaborators() engineCollaborators {
	return engineCollaborators{
		runner:       agent.NewDefaultRunner(),
		verifier:     daemon.ExecVerifier{},
		committer:    daemon.GitCommitter{},
		pusher:       daemon.GitPusher{},
		source:       defaultReviewSourceResolver{},
		worktree:     daemon.GitWorktreeSnapshotter{},
		priorChanges: daemon.GitPriorChangedResolver{},
	}
}

type defaultReviewSourceResolver struct{}

func (defaultReviewSourceResolver) ResolveIssues(ctx context.Context, req reviewsource.ResolveRequest) error {
	if req.Source != reviewsource.SourceCodeRabbit {
		return fmt.Errorf("unsupported Review Source %q; supported value: coderabbit", req.Source)
	}
	return coderabbit.Client{}.ResolveIssues(ctx, req)
}

func (defaultReviewSourceResolver) ResolveIssue(ctx context.Context, req reviewsource.IssueResolveRequest) error {
	if req.Source != reviewsource.SourceCodeRabbit {
		return fmt.Errorf("unsupported Review Source %q; supported value: coderabbit", req.Source)
	}
	return coderabbit.Client{}.ResolveIssue(ctx, req)
}

func (defaultReviewSourceResolver) ReplyToIssue(ctx context.Context, req reviewsource.IssueCommentRequest) (reviewsource.IssueCommentResult, error) {
	if req.Source != reviewsource.SourceCodeRabbit {
		return reviewsource.IssueCommentResult{}, fmt.Errorf("unsupported Review Source %q; supported value: coderabbit", req.Source)
	}
	return coderabbit.Client{}.ReplyToIssue(ctx, req)
}

func defaultWatchReviewEvidence(ctx context.Context, req reviewsource.EvidenceRequest) (reviewsource.Evidence, error) {
	if req.Source != reviewsource.SourceCodeRabbit {
		return reviewsource.Evidence{}, fmt.Errorf("unsupported Review Source %q; supported value: coderabbit", req.Source)
	}
	return coderabbit.Client{}.Evidence(ctx, req)
}

func defaultNewReviewRequester(sink runevent.Sink) reviewsource.ReviewRequester {
	return coderabbit.Client{Sink: sink}
}

func maybeRequestReview(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, runID string, sink runevent.Sink) error {
	if !loaded.Config.ReviewSource.RequestReview {
		return nil
	}
	requester := commandDependenciesForContext(ctx).newReviewRequester(sink)
	if requester == nil {
		return nil
	}
	headSHA, err := commandDependenciesForContext(ctx).watchHeadSHA(ctx, preflightResult.Git.Root)
	if err != nil {
		return fmt.Errorf("detect pushed head for Review Source request: %w", err)
	}
	if _, err := requester.RequestReview(ctx, reviewsource.ReviewRequest{
		RunID:          runID,
		BaseRepository: preflightResult.PullRequest.BaseRepository,
		PRNumber:       preflightResult.PullRequest.Number,
		HeadSHA:        headSHA,
		Command:        loaded.Config.ReviewSource.RequestCommand,
	}); err != nil {
		return fmt.Errorf("request Review Source review for pushed head: %w", err)
	}
	return nil
}

func defaultWatchHeadSHA(ctx context.Context, gitRoot string) (string, error) {
	head, err := preflight.ExecGitRunner{}.RunGit(ctx, gitRoot, "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("detect pushed HEAD: %w", err)
	}
	head = strings.TrimSpace(head)
	if head == "" {
		return "", errors.New("detect pushed HEAD: git returned an empty SHA")
	}
	return head, nil
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

// maybeCommitReviewArtifacts creates the ADR-0036 review artifact commit on
// the user's checkout: one docs commit, separate from Batch fix commits,
// carrying every dirty path under the Run's resolved review artifact root.
// Roots outside the repository — an explicit external Artifact Directory, an
// external Spec Root, or a root reached through a symbolic link — are never
// staged; the Run proceeds without the commit.
func maybeCommitReviewArtifacts(ctx context.Context, req commandRequest, loaded roundconfig.Loaded, preflightResult preflight.Result, committer daemon.Committer, sink runevent.Sink, runID string, roundNumber int, stderr io.Writer, writeGuard daemon.WriteBoundaryGuard) (reviewsource.ArtifactCommit, error) {
	if !loaded.Config.Defaults.AutoCommit {
		return reviewsource.ArtifactCommit{}, nil
	}
	specsRoot := reviewArtifactSpecsRoot(loaded, preflightResult.Git.Root)
	if !reviewArtifactUsesDefaultSpecsRoot(loaded.Config.Specs.Root) {
		resolved, err := roundconfig.ResolveSpecsRoot(loaded, preflightResult.Git.Root)
		if err != nil {
			return reviewsource.ArtifactCommit{}, err
		}
		specsRoot = resolved
	}
	reviewRoot, err := resolveReviewArtifactRoot(ctx, req, preflightResult, specsRoot)
	if err != nil {
		return reviewsource.ArtifactCommit{}, err
	}
	relative, stageable := stageableReviewRoot(preflightResult.Git.Root, reviewRoot)
	if !stageable {
		message := fmt.Sprintf("Review artifacts kept outside the repository (%s); no review artifact commit created.", reviewRoot)
		fmt.Fprintln(stderr, message)
		publishReviewArtifactCommitDecision(ctx, sink, runID, "skipped", message)
		return reviewsource.ArtifactCommit{}, nil
	}
	output, err := commandDependenciesForContext(ctx).reviewSpecGitRunner.RunGit(ctx, preflightResult.Git.Root, "status", "--porcelain", "--", relative)
	if err != nil {
		return reviewsource.ArtifactCommit{}, fmt.Errorf("inspect review artifact changes under %q: %w", relative, err)
	}
	if strings.TrimSpace(output) == "" {
		return reviewsource.ArtifactCommit{}, nil
	}
	parentSHA, err := commandDependenciesForContext(ctx).reviewSpecGitRunner.RunGit(ctx, preflightResult.Git.Root, "rev-parse", "HEAD")
	if err != nil {
		return reviewsource.ArtifactCommit{}, fmt.Errorf("read review artifact commit parent: %w", err)
	}
	parentSHA = strings.TrimSpace(parentSHA)
	if parentSHA == "" {
		return reviewsource.ArtifactCommit{}, errors.New("read review artifact commit parent: git returned an empty SHA")
	}
	message := daemon.ReviewArtifactsCommitMessage(roundNumber, req.pr)
	if writeGuard != nil {
		if err := writeGuard.Check(ctx); err != nil {
			return reviewsource.ArtifactCommit{}, fmt.Errorf("guard review artifact commit: %w", err)
		}
	}
	if err := committer.Commit(ctx, daemon.CommitRequest{
		WorkDir: preflightResult.Git.Root,
		Message: message,
		Paths:   []string{relative},
	}); err != nil {
		return reviewsource.ArtifactCommit{}, fmt.Errorf("create review artifact commit: %w", err)
	}
	if recorder, ok := writeGuard.(interface {
		RecordWriteBoundary(context.Context) error
	}); ok {
		if err := recorder.RecordWriteBoundary(ctx); err != nil {
			return reviewsource.ArtifactCommit{}, fmt.Errorf("record review artifact commit target: %w", err)
		}
	}
	commitSHA, err := commandDependenciesForContext(ctx).reviewSpecGitRunner.RunGit(ctx, preflightResult.Git.Root, "rev-parse", "HEAD")
	if err != nil {
		return reviewsource.ArtifactCommit{}, fmt.Errorf("read created review artifact commit: %w", err)
	}
	commitSHA = strings.TrimSpace(commitSHA)
	if commitSHA == "" || commitSHA == parentSHA {
		return reviewsource.ArtifactCommit{}, errors.New("create review artifact commit: HEAD did not advance")
	}
	fmt.Fprintf(stderr, "Review artifacts commit created: %s\n", message)
	publishReviewArtifactCommitDecision(ctx, sink, runID, "created", fmt.Sprintf("Review artifacts commit created: %s", message))
	return reviewsource.ArtifactCommit{
		CommitSHA:  commitSHA,
		ParentSHA:  parentSHA,
		ReviewRoot: reviewRoot,
		Message:    message,
	}, nil
}

type reviewArtifactEvidenceRequest struct {
	Source         string
	PRNumber       string
	BaseRepository string
	HeadRepository string
	HeadBranch     string
	GitRoot        string
	Commit         reviewsource.ArtifactCommit
	ParentEvidence reviewsource.Evidence
	ParentHeadSHA  string
}

func inheritReviewArtifactEvidence(ctx context.Context, req reviewArtifactEvidenceRequest) (reviewsource.Evidence, bool, error) {
	commit := req.Commit
	if !validReviewArtifactEvidenceRequest(req) {
		return reviewsource.Evidence{}, false, nil
	}

	currentHead, err := commandDependenciesForContext(ctx).reviewSpecGitRunner.RunGit(ctx, req.GitRoot, "rev-parse", "HEAD")
	if err != nil {
		return reviewsource.Evidence{}, false, fmt.Errorf("prove review artifact current head: %w", err)
	}
	if strings.TrimSpace(currentHead) != commit.CommitSHA {
		return reviewsource.Evidence{}, false, nil
	}

	parentLine, err := commandDependenciesForContext(ctx).reviewSpecGitRunner.RunGit(ctx, req.GitRoot, "rev-list", "--parents", "-n", "1", commit.CommitSHA)
	if err != nil {
		return reviewsource.Evidence{}, false, fmt.Errorf("prove review artifact parent: %w", err)
	}
	parents := strings.Fields(parentLine)
	if len(parents) != 2 || parents[0] != commit.CommitSHA || parents[1] != commit.ParentSHA {
		return reviewsource.Evidence{}, false, nil
	}

	subject, err := commandDependenciesForContext(ctx).reviewSpecGitRunner.RunGit(ctx, req.GitRoot, "show", "-s", "--format=%s", commit.CommitSHA)
	if err != nil {
		return reviewsource.Evidence{}, false, fmt.Errorf("prove review artifact commit message: %w", err)
	}
	if strings.TrimSpace(subject) != commit.Message {
		return reviewsource.Evidence{}, false, nil
	}

	relativeRoot, stageable := stageableReviewRoot(req.GitRoot, commit.ReviewRoot)
	if !stageable {
		return reviewsource.Evidence{}, false, nil
	}
	diffOutput, err := commandDependenciesForContext(ctx).reviewSpecGitRunner.RunGit(
		ctx,
		req.GitRoot,
		"diff",
		"--name-only",
		"--no-renames",
		"-z",
		commit.ParentSHA,
		commit.CommitSHA,
		"--",
	)
	if err != nil {
		return reviewsource.Evidence{}, false, fmt.Errorf("prove review artifact diff: %w", err)
	}
	paths := nulSeparatedPaths(diffOutput)
	if len(paths) == 0 {
		return reviewsource.Evidence{}, false, nil
	}
	if !reviewArtifactPathsWithinRoot(paths, relativeRoot) {
		return reviewsource.Evidence{}, false, nil
	}

	refreshed, err := commandDependenciesForContext(ctx).watchReviewEvidence(ctx, reviewsource.EvidenceRequest{
		Source:          req.Source,
		PRNumber:        req.PRNumber,
		BaseRepository:  req.BaseRepository,
		HeadRepository:  req.HeadRepository,
		HeadBranch:      req.HeadBranch,
		ExpectedHeadSHA: commit.ParentSHA,
	})
	if err != nil || !exactVerifiedEvidence(refreshed, commit.ParentSHA) {
		return reviewsource.Evidence{}, false, nil
	}
	return reviewsource.Evidence{
		State:           reviewsource.EvidenceVerified,
		Kind:            reviewsource.EvidenceKindArtifactOnlyDescendant,
		Identity:        "daemon_review_artifact:" + commit.CommitSHA,
		ExpectedHeadSHA: commit.CommitSHA,
		ObservedHeadSHA: commit.CommitSHA,
		ParentHeadSHA:   commit.ParentSHA,
		Conclusion:      "inherited",
		Detail:          reviewsource.BoundEvidenceDetail("Inherited verified parent Evidence for the exact Daemon review-artifact commit."),
	}, true, nil
}

func validReviewArtifactEvidenceRequest(req reviewArtifactEvidenceRequest) bool {
	commit := req.Commit
	return commit.CommitSHA != "" &&
		commit.ParentSHA != "" &&
		commit.ReviewRoot != "" &&
		commit.Message != "" &&
		req.ParentHeadSHA == commit.ParentSHA &&
		exactVerifiedEvidence(req.ParentEvidence, commit.ParentSHA)
}

func reviewArtifactPathsWithinRoot(paths []string, relativeRoot string) bool {
	rootPrefix := relativeRoot + "/"
	for _, path := range paths {
		clean := filepath.ToSlash(filepath.Clean(path))
		if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || !strings.HasPrefix(clean, rootPrefix) {
			return false
		}
	}
	return true
}

func exactVerifiedEvidence(evidence reviewsource.Evidence, headSHA string) bool {
	return evidence.State == reviewsource.EvidenceVerified &&
		evidence.ExpectedHeadSHA == headSHA &&
		evidence.ObservedHeadSHA == headSHA
}

func nulSeparatedPaths(output string) []string {
	output = strings.TrimSuffix(output, "\x00")
	if output == "" {
		return nil
	}
	return strings.Split(output, "\x00")
}

// stageableReviewRoot reports the repo-relative review root when it can be
// staged: inside the repository working tree and not reached through a
// symbolic link component.
func stageableReviewRoot(gitRoot string, reviewRoot string) (string, bool) {
	relative, err := filepath.Rel(gitRoot, reviewRoot)
	if err != nil || relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", false
	}
	current := gitRoot
	for _, part := range strings.Split(relative, string(filepath.Separator)) {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			// Missing segments cannot cross a symlink; git status decides next.
			break
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", false
		}
	}
	return filepath.ToSlash(relative), true
}

// publishReviewArtifactCommitDecision journals the review artifact commit
// decision as a Daemon-source commit event.
func publishReviewArtifactCommitDecision(ctx context.Context, sink runevent.Sink, runID string, decision string, summary string) {
	payload, err := json.Marshal(map[string]any{"decision": decision, "artifacts": "review"})
	if err != nil {
		return
	}
	_ = sink.Publish(context.WithoutCancel(ctx), runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonCommit,
		Summary: runevent.BoundSummary(summary),
		Time:    time.Now().UTC(),
		Payload: payload,
	})
}

func maybeRunFinalPush(ctx context.Context, engine *daemon.Engine, sink runevent.Sink, runID string, loaded roundconfig.Loaded, preflightResult preflight.Result, workDir string, batchCommitCreated bool, stderr io.Writer) (bool, error) {
	if !preflightResult.PushPlan.Enabled {
		fmt.Fprintln(stderr, "Final Push skipped: auto-push disabled or no push target configured.")
		publishPushDecision(ctx, sink, runID, "skipped", "Final Push skipped: auto-push disabled or no push target configured.", 0)
		return false, nil
	}
	if preflightResult.PushPlan.Force {
		return false, errors.New("final push rejected: force-push is not allowed in the MVP")
	}
	if !loaded.Config.Defaults.AutoCommit {
		fmt.Fprintln(stderr, "Final Push skipped: auto-commit disabled.")
		publishPushDecision(ctx, sink, runID, "skipped", "Final Push skipped: auto-commit disabled.", 0)
		return false, nil
	}
	if preflightResult.Git.UnpushedCommits == 0 && !batchCommitCreated {
		fmt.Fprintln(stderr, "Final Push skipped: no local commits are waiting for the PR Head Branch.")
		publishPushDecision(ctx, sink, runID, "skipped", "Final Push skipped: no local commits are waiting for the PR Head Branch.", 0)
		return false, nil
	}
	if err := engine.FinalPush(ctx, daemon.FinalPushRequest{
		RunID:   runID,
		WorkDir: workDir,
		Remote:  preflightResult.PushPlan.Remote,
		Branch:  preflightResult.PushPlan.Branch,
	}); err != nil {
		return false, fmt.Errorf("final push: %w", err)
	}
	fmt.Fprintf(stderr, "Final Push completed: git push %s HEAD:%s\n", preflightResult.PushPlan.Remote, preflightResult.PushPlan.Branch)
	return true, nil
}

func markRunFailed(ctx context.Context, runStore *store.Store, runID string) {
	completed, ok := completeFailedRun(ctx, runStore, runID)
	if ok {
		publishTerminalCompletion(ctx, runStore, nil, io.Discard, completed, 0)
	}
}

func markRunFailedAndNotify(ctx context.Context, runStore *store.Store, runID string, notifier roundnotify.Notifier, stderr io.Writer) {
	completed, ok := completeFailedRun(ctx, runStore, runID)
	if !ok {
		return
	}
	publishTerminalCompletion(ctx, runStore, notifier, stderr, completed, 0)
}

func completeFailedRun(ctx context.Context, runStore *store.Store, runID string) (store.CompleteRunResult, bool) {
	completeCtx := withoutCancelOrBackground(ctx)
	completed, err := runStore.CompleteRun(completeCtx, runID, store.StateFailed)
	if err != nil {
		return store.CompleteRunResult{}, false
	}
	return completed, true
}

func closeAgentSession(ctx context.Context, runner agent.Runner, runtime agent.RuntimeSpec, session agent.SessionRef, runID string, runStore *store.Store) {
	if runner == nil || runStore == nil || strings.TrimSpace(session.Name) == "" {
		return
	}
	closeCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	_ = runner.EndSession(closeCtx, runtime, session)
	publishAgentSessionStatus(closeCtx, store.JournalSink{Store: runStore}, runID, agent.AgentSessionClosedStatus)
}

func publishAgentSessionStatus(ctx context.Context, sink runevent.Sink, runID string, status string) {
	payload, err := json.Marshal(struct {
		Status string `json:"status"`
	}{Status: status})
	if err != nil {
		return
	}
	_ = sink.Publish(context.WithoutCancel(ctx), runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentStatus,
		Summary: runevent.BoundSummary("SESSION " + strings.ToUpper(status) + "\n"),
		Time:    time.Now().UTC(),
		Payload: payload,
	})
}

// publishRunOutcome appends the terminal outcome event after CompleteRun so
// terminal states are provable from the journal alone.
func watchTerminalCompletionContext(req commandRequest, run store.Run, result watch.Result) terminalCompletionContext {
	reviewIssuesKnown := result.ReviewIssuesKnown
	terminal := terminalCompletionContext{
		Remaining:         result.Remaining,
		Reason:            result.TerminalReason,
		NextAction:        result.NextAction,
		ReviewIssuesKnown: &reviewIssuesKnown,
		Evidence:          result.Evidence,
		VerifiedHeadSHA:   result.VerifiedHeadSHA,
	}
	if run.ID != "" {
		terminal.AttachCommand = "roundfix attach " + run.ID
	}
	if req.detachChild != nil && req.artifactDir != "" && run.ID != "" {
		terminal.ConsoleLog = detachedConsoleLogPath(req.artifactDir, run.ID)
	}
	return terminal
}

type terminalCompletionContext struct {
	Remaining         int
	Reason            string
	NextAction        string
	ReviewIssuesKnown *bool
	ConsoleLog        string
	AttachCommand     string
	Evidence          reviewsource.Evidence
	VerifiedHeadSHA   string
}

func normalizedTerminalCompletionContext(state string, terminal terminalCompletionContext) terminalCompletionContext {
	terminal.Reason = boundTerminalContextText(terminal.Reason)
	terminal.NextAction = boundTerminalContextText(terminal.NextAction)
	if state == store.StateClean {
		return terminal
	}
	defaultReason, defaultAction := terminalCompletionDefaults(state)
	if terminal.Reason == "" {
		terminal.Reason = defaultReason
	}
	if terminal.NextAction == "" {
		terminal.NextAction = defaultAction
	}
	return terminal
}

func terminalCompletionDefaults(state string) (string, string) {
	switch state {
	case store.StateFetched:
		return "The fetch Run completed without resolving Review Issues.", "Run resolve or watch when the Review Issues are ready for Agent work."
	case store.StateStopped:
		return "A Stop Request ended the Run.", "Inspect the preserved work before starting another Run."
	case store.StateCleanUnverified:
		return "Merge-Ready was not confirmed for the completed Run.", "Confirm the pull request's Review Source Evidence before merging."
	case store.StateReviewSkipped:
		return "The Review Source explicitly skipped the review.", "Reduce or split the pull request, then request another Review Source review."
	case store.StateMaxRoundsReached:
		return "The configured maximum number of Rounds was reached.", "Review the remaining Review Issues before deciding whether to start another Run."
	case store.StateBudgetExceeded:
		return "The Run Budget was exhausted.", "Inspect the Run Event Stream before starting another Run with an appropriate budget."
	case store.StateTimedOut:
		return "The Run timed out before reaching Clean.", "Inspect the Run Event Stream, restore the missing prerequisite, and start another Run."
	case store.StateFailed:
		return "The Run failed before it could complete.", "Inspect the diagnostics, correct the failure, and start another Run."
	case store.StateCheckoutMoved:
		return "The checkout no longer matches the target recorded at Preflight.", "Restore the PR Head Branch checkout, then start another Run."
	case store.StateIntegrationPending:
		return "Completed work could not be integrated into the target branch.", "Inspect the retained Run Worktree and follow the reported integration command."
	case store.StateUnresolved:
		return "The Run completed with Unresolved Review Issues.", "Review the remaining Review Issues before starting another Run."
	default:
		return "The Run reached a non-Clean outcome.", "Inspect the Run Event Stream before deciding the next recovery step."
	}
}

func boundTerminalContextText(text string) string {
	return reviewsource.BoundEvidenceDetail(strings.Join(strings.Fields(text), " "))
}

func terminalEvidenceHead(evidence reviewsource.Evidence) string {
	if evidence.ObservedHeadSHA != "" {
		return evidence.ObservedHeadSHA
	}
	return evidence.ExpectedHeadSHA
}

func publishRunOutcome(ctx context.Context, runStore *store.Store, runID string, state string, terminal terminalCompletionContext, stderr io.Writer) {
	terminal = normalizedTerminalCompletionContext(state, terminal)
	payload, err := json.Marshal(runevent.OutcomePayload{
		State:             state,
		Remaining:         terminal.Remaining,
		Reason:            terminal.Reason,
		NextAction:        terminal.NextAction,
		ReviewIssuesKnown: terminal.ReviewIssuesKnown,
		ConsoleLog:        terminal.ConsoleLog,
		AttachCommand:     terminal.AttachCommand,
		EvidenceState:     string(terminal.Evidence.State),
		EvidenceKind:      string(terminal.Evidence.Kind),
		EvidenceHeadSHA:   terminalEvidenceHead(terminal.Evidence),
		VerifiedHeadSHA:   terminal.VerifiedHeadSHA,
	})
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

func publishTerminalCompletion(
	ctx context.Context,
	runStore *store.Store,
	notifier roundnotify.Notifier,
	stderr io.Writer,
	completed store.CompleteRunResult,
	remaining int,
) {
	publishTerminalCompletionWithContext(ctx, runStore, notifier, stderr, completed, terminalCompletionContext{Remaining: remaining})
}

func publishTerminalCompletionWithContext(
	ctx context.Context,
	runStore *store.Store,
	notifier roundnotify.Notifier,
	stderr io.Writer,
	completed store.CompleteRunResult,
	terminal terminalCompletionContext,
) {
	if !completed.Transitioned {
		return
	}
	terminal = normalizedTerminalCompletionContext(completed.State, terminal)
	publishRunOutcome(ctx, runStore, completed.ID, completed.State, terminal, stderr)
	notifyTerminalOutcome(ctx, runStore, notifier, stderr, completed.Run, terminal)
}

func journalStopPrimaryFailure(ctx context.Context, runStore *store.Store, runID string, primary error) {
	if runStore == nil || strings.TrimSpace(runID) == "" || primary == nil {
		return
	}
	reason := primary.Error()
	payload, err := json.Marshal(map[string]string{
		"event":  "force_stop_primary_failure",
		"reason": reason,
	})
	if err != nil {
		return
	}
	_ = (store.JournalSink{Store: runStore}).Publish(withoutCancelOrBackground(ctx), runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonStatus,
		Summary: runevent.BoundSummary("Force Stop primary failure: " + reason),
		Time:    time.Now().UTC(),
		Payload: payload,
	})
}

func reportSecondaryCleanupWarnings(
	ctx context.Context,
	runStore *store.Store,
	runID string,
	warnings []cleanupWarning,
	stderr io.Writer,
) {
	for _, warning := range warnings {
		if warning.kind != cleanupWarningFailure {
			fmt.Fprintf(stderr, "%s: %s\n", app.Name, warning.text)
			continue
		}
		summary := "Secondary cleanup warning: " + warning.text
		fmt.Fprintf(stderr, "%s: %s\n", app.Name, summary)
		if runStore == nil || strings.TrimSpace(runID) == "" {
			continue
		}
		payload, err := json.Marshal(map[string]string{
			"event":  "secondary_cleanup_warning",
			"reason": warning.text,
		})
		if err != nil {
			continue
		}
		_ = (store.JournalSink{Store: runStore}).Publish(withoutCancelOrBackground(ctx), runevent.RunEvent{
			RunID:   runID,
			Source:  runevent.SourceDaemon,
			Kind:    runevent.KindDaemonStatus,
			Summary: runevent.BoundSummary(summary),
			Time:    time.Now().UTC(),
			Payload: payload,
		})
	}
}

func warnCleanRunWorktreeCleanupFailed(ctx context.Context, runStore *store.Store, runID string, worktreePath string, cleanupErr error, stderr io.Writer) {
	if cleanupErr == nil {
		return
	}
	reason := cleanupErr.Error()
	summary := fmt.Sprintf("%s: Run Worktree cleanup failed; kept %s: %s", app.Name, worktreePath, reason)
	fmt.Fprintln(stderr, summary)
	payload, err := json.Marshal(map[string]string{
		"path":   worktreePath,
		"reason": reason,
	})
	if err != nil {
		return
	}
	_ = (store.JournalSink{Store: runStore}).Publish(context.WithoutCancel(ctx), runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonStatus,
		Summary: runevent.BoundSummary(summary),
		Time:    time.Now().UTC(),
		Payload: payload,
	})
}

func outcomeNotifierFromConfig(ctx context.Context, config roundconfig.Config) roundnotify.Notifier {
	factory := commandDependenciesForContext(ctx).newOutcomeNotifier
	if factory == nil {
		return nil
	}
	return factory(config)
}

const outcomeNotificationTimeout = 30 * time.Second

func notifyTerminalOutcome(
	ctx context.Context,
	runStore *store.Store,
	notifier roundnotify.Notifier,
	stderr io.Writer,
	run store.Run,
	terminal terminalCompletionContext,
) {
	if notifier == nil {
		return
	}
	notifyCtx, cancel := context.WithTimeout(withoutCancelOrBackground(ctx), outcomeNotificationTimeout)
	defer cancel()
	receipt, err := notifier.Notify(notifyCtx, outcomeFromRun(run, terminal))
	if err != nil {
		receipt.Status = roundnotify.StatusFailed
	}
	if receipt.Status == "" {
		receipt.Status = roundnotify.StatusSkipped
	}
	if receipt.Route == "" {
		receipt.Route = roundnotify.RouteDisabled
	}
	if receipt.CompletedAt.IsZero() {
		receipt.CompletedAt = time.Now().UTC()
	}
	journalOutcomeNotificationReceipt(notifyCtx, runStore, run.ID, receipt, err, stderr)
	if err != nil {
		reportOutcomeNotificationFailure(err, stderr)
	}
}

func outcomeFromRun(run store.Run, terminal terminalCompletionContext) roundnotify.Outcome {
	return roundnotify.Outcome{
		RunID:             run.ID,
		State:             run.State,
		Kind:              run.Kind,
		Target:            outcomeTarget(run),
		Reason:            terminal.Reason,
		ConsoleLog:        terminal.ConsoleLog,
		AttachCommand:     terminal.AttachCommand,
		ReviewIssuesKnown: terminal.ReviewIssuesKnown,
		NextAction:        terminal.NextAction,
	}
}

func outcomeTarget(run store.Run) string {
	if run.Kind == store.KindImplement {
		return "spec:" + strings.TrimSpace(run.SpecSlug)
	}
	if strings.TrimSpace(run.PRNumber) != "" {
		return "pr:" + strings.TrimSpace(run.PRNumber)
	}
	return ""
}

func reportOutcomeNotificationFailure(err error, stderr io.Writer) {
	reason := strings.TrimSpace(err.Error())
	if reason == "" {
		reason = "unknown error"
	}
	warning := fmt.Sprintf("%s: outcome notification failed: %s", app.Name, reason)
	if stderr != nil {
		fmt.Fprintln(stderr, warning)
	}
}

func journalOutcomeNotificationReceipt(
	ctx context.Context,
	runStore *store.Store,
	runID string,
	receipt roundnotify.NotificationReceipt,
	notifyErr error,
	stderr io.Writer,
) {
	if runStore == nil {
		return
	}
	reason := ""
	if notifyErr != nil {
		reason = boundTerminalContextText(notifyErr.Error())
	}
	eventName := "outcome_notification_" + string(receipt.Status)
	payload, err := json.Marshal(runevent.NotificationReceiptPayload{
		Event:       eventName,
		Route:       string(receipt.Route),
		Status:      string(receipt.Status),
		CompletedAt: receipt.CompletedAt,
		Reason:      reason,
	})
	if err != nil {
		return
	}
	summary := fmt.Sprintf("Outcome notification %s via %s.", receipt.Status, receipt.Route)
	if err := (store.JournalSink{Store: runStore}).Publish(withoutCancelOrBackground(ctx), runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonStatus,
		Summary: runevent.BoundSummary(summary),
		Time:    receipt.CompletedAt,
		Payload: payload,
	}); err != nil && stderr != nil {
		fmt.Fprintf(stderr, "Warning: notification receipt event not journaled: %v\n", err)
	}
}

func withoutCancelOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return context.WithoutCancel(ctx)
}

func commandUsage(name string) string {
	switch name {
	case "events":
		return `Usage:
  roundfix events <run-id> [--follow] [--filter <categories>]

Replays one Run's Supervisor event stream from the Run Database as
roundfix-events/v1 JSONL. With --follow, continues until the Run reaches a
terminal state or the command is interrupted.

Options:
  --follow  Continue after replay and exit after the terminal event drains
  --filter  Comma-separated categories: task-status,batch,verification,outcome,agent-selection
`
	case "attach":
		return `Usage:
  roundfix attach [<run-id>] [--no-input]

Replays the Run's event timeline from the Run Database, read-only.
Attach never creates Runs, fetches, starts Agents, commits, pushes, or
resolves Review Source threads. Without a Run ID at an interactive
terminal, attach opens the Run Browser to pick a Run; leaving the Live
Run View returns to a refreshed browser until you quit.

Options:
  --run-id    Run ID to attach to (same as the positional argument)
  --no-input  Fail instead of opening the Run Browser when no Run ID is passed
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
	case "setup":
		return `Usage:
  roundfix setup [--yes] [--no-input]

Checks Node.js and the minimum supported acpx version, proves the effective
official Codex and Claude adapters and generated Agent Selection
profile readiness, then offers acpx local adapter overrides, User Config, and
Project Config. Proposed profiles are exact-proved before writing. Legacy
Codex and Claude override migration requires authorization. Each check prints
one deterministic report line with ok, installed, skipped, offered: declined,
or failed.

Options:
  --yes       Accept every offered install or file change
  --no-input  Skip offered changes instead of prompting
`
	case "doctor":
		return `Usage:
  roundfix doctor

Diagnoses this machine's readiness for Runs. Checks Node.js, the minimum
supported acpx version, the effective adapter, the required
Agent Selection Profiles, the Repository Skill Set, and codex runtime hygiene.
The aggregate profiles: line exact-proves every distinct tuple. The skills:
line compares Roundfix-owned artifacts with the running binary and external
artifacts with skills-lock.json. Each failure reports its next action.
Doctor is offline, read-only, and mutates nothing.
`
	case "gc":
		return `Usage:
  roundfix gc [--dry-run]
  roundfix gc sanitize [--apply]

Prunes Run Event Journal rows for terminal Runs older than Journal Retention,
then removes their run artifact directories and orphaned runs/<id> directories
under the resolved run artifact root. It never deletes Run rows or Active Run locks,
and it never removes Review artifacts outside the run artifact root.

The sanitize subcommand discovers every recorded Artifact Root from the
machine-wide Run Database, classifies it with durable and filesystem evidence,
and preserves every ambiguous path. It is a dry-run unless --apply is explicit.

Options:
  --dry-run  List the Runs, journal rows, and artifact bytes that would be pruned without changing anything
  --apply    With sanitize, remove only proven retention-eligible or absent Run artifact directories
`
	case "storage":
		return `Usage:
  roundfix storage report

Commands:
  report  Measure bytes and row counts by repository, state, table, and Artifact Root.

The report is read-only, accepts no flags, and needs no Git repository. It
reads the machine-wide Run Database and recorded Artifact Roots from Roundfix
Home without migrating, locking for writes, or changing any byte.
`
	case "upgrade":
		return `Usage:
  roundfix upgrade [--check]

Resolves the latest Roundfix release for this platform through the GitHub CLI.
Without --check, downloads the matching asset, verifies it, and atomically
replaces the current executable. If no releases exist, reports that cleanly.

Options:
  --check  Report the latest release outcome without installing it
`
	case "runs":
		return `Usage:
  roundfix runs
  roundfix runs list [--all] [--state <active|terminal|all>] [--limit N]

Lists Runs from the Run Database newest first. By default the listing is
scoped to the current repository and shows the 20 newest Active Runs. When
the state filter or the bound hides Runs, one trailing stderr note names the
hidden count and the widening flag. Without a subcommand at an interactive
terminal, runs opens the Run Browser; non-interactive contexts must use
'roundfix runs list'.

Commands:
  list  Print run id, state, kind, target, agent, start time (UTC), duration,
        and local branch columns

Options:
  --all    List Runs from every repository and include the repository column
  --state  Filter by Run state: active (default), terminal, or all
  --limit  Print at most N matching Runs, newest first; 0 lists all (default 20)
`
	case "fetch":
		return `Usage:
  roundfix fetch --source coderabbit --pr <number> [--spec <slug>] [--round <number|auto>] [--no-input]

Behavior:
  Runs Branch Integrity Preflight before fetching. Fetch executes in the
  user's checkout, creates no Run Worktree, starts no Agent, and never commits
  or pushes.

Options:
  --source       Review Source. Supported: coderabbit
  --pr           Open Pull Request number
  --spec         Spec slug under docs/specs/
  --round        Round number or auto
  --artifact-dir Artifact Directory
  --base-repo    Explicit base repository, owner/name
  --head-repo    Explicit Head Repository, owner/name
  --head-branch  Explicit PR Head Branch
  --skip-branch-integrity
                 Skip pending Run Branch and Active Run guardrails after publishing a PR audit comment
  --interactive  Open Interactive Input before starting
  --no-input     Fail instead of opening Interactive Input
`
	case "resolve":
		return `Usage:
  roundfix resolve --pr <number> [--spec <slug>] [--round <number|all>] [--no-input]
  roundfix resolve --pr <number> --agent <agent> --model <model> --reasoning-effort <effort> [--spec <slug>]

Behavior:
  Runs Branch Integrity Preflight and clean tracked checkout validation before
  Agent work. Resolve executes in the user's checkout and creates no Run
  Worktree. Omit all Agent Selection flags to use the review profile. A one-Run
  override requires --agent, --model, and --reasoning-effort together.

Options:
  --pr           Open Pull Request number
  --spec         Spec slug under docs/specs/
  --agent        Agent runtime. Supported: codex, claude, opencode
  --model        Agent model override
  --reasoning-effort Default reasoning effort override
  --agent-command Agent command override
  --agent-full-access Opt into Agent runtime full-access mode
  --no-agent-console Hide Agent-source console events from non-TTY stderr
  --detach       Start a Detached Run and print attach/stop commands
  --round        Round number or all
  --artifact-dir Artifact Directory
  --base-repo    Explicit base repository, owner/name
  --head-repo    Explicit Head Repository, owner/name
  --head-branch  Explicit PR Head Branch
  --skip-branch-integrity
                 Skip pending Run Branch and Active Run guardrails after publishing a PR audit comment
  --interactive  Open Interactive Input before starting
  --no-input     Fail instead of opening Interactive Input
`
	case "watch":
		return `Usage:
  roundfix watch --source coderabbit --pr <number> [--spec <slug>] [--until-clean] [--max-rounds <number>] [--no-input]
  roundfix watch --source coderabbit --pr <number> --agent <agent> --model <model> --reasoning-effort <effort> [--spec <slug>]

Behavior:
  Runs Branch Integrity Preflight and clean tracked checkout validation before
  Agent work. Watch executes in the user's checkout and creates no Run
  Worktree. With --until-clean, Clean requires accepted Review Source Evidence
  on the pushed head. The only check-or-status route to a verified head is a
  recognised review-completed current-head CodeRabbit check or commit status;
  a current-head CodeRabbit APPROVED review is also accepted, with zero
  unresolved CodeRabbit threads.
  An unrecognised signal resolves to pending even when its check
  conclusion is success:
  a green check is not evidence that a review ran.
  An explicit Review Source refusal resolves to skipped;
  watch will not merge that head or clear it for merge. If no accepted
  Evidence appears within the grace period, watch ends CleanUnverified
  and exits 3. Omit all Agent Selection flags to use the review profile.
  A one-Run override requires --agent, --model, and --reasoning-effort
  together.

Options:
  --source       Review Source. Supported: coderabbit
  --pr           Open Pull Request number
  --spec         Spec slug under docs/specs/
  --agent        Agent runtime. Supported: codex, claude, opencode
  --model        Agent model override
  --reasoning-effort Default reasoning effort override
  --agent-command Agent command override
  --agent-full-access Opt into Agent runtime full-access mode
  --no-agent-console Hide Agent-source console events from non-TTY stderr
  --detach       Start a Detached Run and print attach/stop commands
  --until-clean  Repeat until no Unresolved Review Issues remain and accepted Review Source Evidence confirms the pushed head
  --max-rounds   Maximum Review Source rounds
  --artifact-dir Artifact Directory
  --base-repo    Explicit base repository, owner/name
  --head-repo    Explicit Head Repository, owner/name
  --head-branch  Explicit PR Head Branch
  --skip-branch-integrity
                 Skip pending Run Branch and Active Run guardrails after publishing a PR audit comment
  --interactive  Open Interactive Input before starting
  --no-input     Fail instead of opening Interactive Input
`
	case "implement":
		return implementUsage
	case "settle":
		return settleUsage
	case "reconcile":
		return `Usage:
  roundfix reconcile [run-id] [--apply] [--format <text|json>]

Inspects one terminal spec Run in the current repository, or every terminal
spec Run in the current repository when no Run ID is supplied. The default is
a read-only report. --apply releases only entries classified and revalidated
safe or superseded during the current invocation; dirty, unintegrated, unknown,
and already released entries remain successful preserved results.

Options:
  --apply   Release freshly revalidated safe and superseded Run Worktrees and Run Branches
  --format  Output format: text (default) or json
`
	case "release":
		return `Usage:
  roundfix release plan [--from <tag>] [--to <revision>] [--impact <none|patch|minor|major> --reason <text>] [--format <text|json>]
  roundfix release plan --reset-to <version> [--format <text|json>]

Commands:
  plan  Analyze committed changes and propose the next semantic version.

Release planning is read-only: it creates no Run, reads no Roundfix config,
and never edits files, refs, tags, packages, releases, remotes, or
configuration. Range planning stays local; reset planning reads complete Git
tag and paginated GitHub Release inventory.
`
	case "release plan":
		return `Usage:
  roundfix release plan [--from <tag>] [--to <revision>] [--impact <none|patch|minor|major> --reason <text>] [--format <text|json>]
  roundfix release plan --reset-to <version> [--format <text|json>]

Builds a read-only Release Plan from a stable vMAJOR.MINOR.PATCH base through
a committed target revision. --from defaults to the latest reachable stable
tag; --to defaults to committed HEAD. --reset-to inventories every local and
remote stable tag and every paginated GitHub Release for a clean committed
HEAD, binds them to a plan digest, and exposes no deletion action.

Decision states:
  ready                           Patch release can proceed without version approval.
  approval_required               Minor, major, or breaking version decision needs approval.
  manual_classification_required  Ambiguous commits need --impact and --reason.
  no_release                      Maintenance-only changes require no release.

Exit codes:
  0  ready or no_release
  2  invalid flags, dirty tree, invalid range, or repository failure
  3  approval_required or manual_classification_required

Options:
  --from    Stable release tag to use as the base, for example v1.2.3
  --to      Target revision to analyze; defaults to HEAD
  --impact  Manual impact for ambiguous changes: none, patch, minor, or major
  --reason  Non-empty reason required with --impact
  --reset-to
            Stable reset target; incompatible with --from, --to, --impact,
            and --reason
  --format  Output format: text or json (default text)

The command creates no Run, reads no Roundfix configuration, and never mutates
files, refs, tags, remotes, packages, releases, or configuration. Reset mode
uses read-only Git and GitHub inventory calls; any deletion requires separate
explicit post-QA authority.
`
	case "baseline":
		return `Usage:
  roundfix baseline [--repo <path>] [--format <text|json>]
  roundfix baseline plan (--profile <id> | --profile-file <draft.json>) [--decision <id=value> ...] [--decision-file <path> ...] [--repo <path>] [--format <text|json>]
  roundfix baseline apply --plan <file> --confirm-plan <digest> [--repo <path>] [--format <text|json>]
  roundfix baseline capabilities check [--profile <id>] [--repo <path>] [--format <text|json>]
  roundfix baseline profile init --id <id> [--from <built-in-id>]
  roundfix baseline profile show <id> [--format <text|json>]
  roundfix baseline profile validate [<id>|<path>] [--format <text|json>]
  roundfix baseline skills restore --profile <id> [--skill <name> ...] [--source-dir <path>] [--confirm-plan <digest>] [--repo <path>] [--format <text|json>]
  roundfix baseline assets sync --source-dir <path> [--check] [--format <text|json>]

The root command guides one interactive adoption or update from repository
preflight through Baseline verification. Numbered linear prompts collect one
Baseline Profile and repository decisions, then present one consolidated
Change Plan with file changes first. Repository mutation occurs only after
explicit confirmation of the displayed Plan Digest.

Commands:
  plan     Automation: emit a portable, digest-bound Baseline Plan without prompting or writing.
  apply    Automation: apply and verify exactly one approved portable Baseline Plan without prompting.
  capabilities  Re-check Profile capability evidence without decisions, prompts, or writes.
  profile  Author, inspect, and validate built-in or repository-owned Baseline Profiles.
  skills   Preview or apply immutable external Repository Skill Set restoration.
  assets   Check or refresh Go-owned canonical Baseline setup snapshots.

The interactive root command refuses redirected or absent terminal input.
Automation must use baseline plan followed by baseline apply with the exact
approved Plan Digest.

Repository-owned profiles live only under
.roundfix/baseline/profiles/<id>.json and may reference only entries in the
embedded Baseline catalog.
`
	case "baseline capabilities check":
		return `Usage:
  roundfix baseline capabilities check [--profile <id>] [--repo <path>] [--format <text|json>]

Re-checks Repository Capability evidence through the same evaluator and
divergence renderer used by Baseline planning. The command accepts and resolves
no decisions, never writes repository or journal bytes, never executes a
candidate or repository command, and never uses the network.

JSON output uses roundfix/baseline-capability-recheck/v1.

When --profile is omitted, the command resolves the current Baseline Profile
from a valid Setup Manifest. A repository without either source returns a named
Profile error instead of an empty result.

Exit codes:
  0  capability evidence evaluated with no blocking divergence
  1  output failure
  2  invalid arguments, repository failure, or no resolvable Baseline Profile
  3  capability evidence evaluated with a blocking divergence

Options:
  --profile  Built-in or repository-owned Baseline Profile; defaults to the current Setup Manifest
  --repo     Git worktree or a path inside it (default current directory)
  --format   Output format: text or json (default text)
`
	case "baseline plan":
		return `Usage:
  roundfix baseline plan (--profile <id> | --profile-file <draft.json>) [--decision <id=value> ...] [--decision-file <path> ...] [--repo <path>] [--format <text|json>]

Builds the complete portable roundfix/baseline-plan/v1 document from
clone-stable Git lineage, bounded repository preimages, one selected Baseline
Profile ID or strict repository-owned Profile draft, and normalized decisions.
The two Profile inputs are mutually exclusive and produce the same normalized
Plan when they resolve to the same Profile. JSON is the portable apply input;
text is a concise file-level projection. Missing decisions return a
roundfix/baseline-result/v1 next action without a partial plan.

The command never prompts, writes repository bytes, executes
repository-defined commands, or uses the network. Select instruction handling
with --decision preservation.mode=greenfield|preservation.

Exit codes:
  0  complete Baseline Plan emitted
  2  invalid arguments, Git/repository failure, or unsafe bounded carrier
  3  a decision, manual classification, or repository-alignment action is required

Options:
  --profile       Built-in or repository-owned Baseline Profile
  --profile-file  Strict repository-owned Profile draft; mutually exclusive with --profile
  --decision      Decision as id=value; repeat for multiple answers
  --decision-file Strict Decision Document path; repeat to merge inputs
  --repo          Git worktree or a path inside it (default current directory)
  --format        Output format: text or json (default text)
`
	case "baseline apply":
		return `Usage:
  roundfix baseline apply --plan <file> --confirm-plan <digest> [--repo <path>] [--format <text|json>]

Strictly parses one portable roundfix/baseline-plan/v1 document, confirms its
exact Plan Digest, validates clone-stable Git lineage and every bounded
preimage, then applies only the supplied postimages through the recoverable
transaction. An exact empty reapply is a verified success.

Baseline verification checks managed postimages, immutable backups, carrier
relationships, Setup Manifest identity, retention accounting, and resolved
audit state. Repository formatter and Verification commands are reported as
recommendations and are never run.

Exit codes:
  0  approved plan applied or already applied, and Baseline verification passed
  1  apply, verification, output, rollback, or recovery failure
  2  invalid arguments, plan schema, or unsafe repository
  3  confirmation mismatch, stale preimage, or unrelated Git lineage

Options:
  --plan          Portable roundfix/baseline-plan/v1 JSON file
  --confirm-plan  Exact approved Plan Digest from the supplied document
  --repo          Git worktree or a path inside it (default current directory)
  --format        Output format: text or json (default text)
`
	case "baseline profile":
		return `Usage:
  roundfix baseline profile init --id <id> [--from <built-in-id>]
  roundfix baseline profile show <id> [--format <text|json>]
  roundfix baseline profile validate [<id>|<path>] [--format <text|json>]

Commands:
  init      Create one repository-owned profile from an embedded built-in profile.
  show      Resolve one built-in or repository-owned profile.
  validate  Validate one profile ID, one repository profile path, or all repository profiles.
`
	case "baseline profile init":
		return `Usage:
  roundfix baseline profile init --id <id> [--from <built-in-id>]

Creates .roundfix/baseline/profiles/<id>.json exclusively. The declaration
copies allowed embedded entry IDs from the selected built-in source; it does
not compose profiles, copy assets, or accept executable or remote content.

Options:
  --id    Required lowercase repository-owned Baseline Profile ID
  --from  Embedded built-in Baseline Profile (default go-cli-tui)
`
	case "baseline profile show":
		return `Usage:
  roundfix baseline profile show <id> [--format <text|json>]

Resolves one built-in or repository-owned Baseline Profile against the
embedded catalog and prints its normalized composition and digest.

Options:
  --format  Output format: text or json (default text)
`
	case "baseline profile validate":
		return `Usage:
  roundfix baseline profile validate [<id>|<path>] [--format <text|json>]

Validates one built-in or repository-owned Baseline Profile by ID, one direct
file under .roundfix/baseline/profiles, or every repository-owned profile when
no target is supplied. The command reads no user-scoped profile catalog.

Options:
  --format  Output format: text or json (default text)
`
	case "baseline skills restore":
		return `Usage:
  roundfix baseline skills restore --profile <id> [--skill <name> ...] [--source-dir <path>] [--confirm-plan <digest>] [--repo <path>] [--format <text|json>]

Previews or applies exact external Repository Skill Set restoration for one
built-in Baseline Profile. Sources are grouped by immutable provider,
repository, and commit provenance. Source bytes, portable tree digests,
skills-lock.json adapter compatibility, targets, and the complete preimage are
validated before mutation.

A non-empty preview exits 3 and returns its exact Plan Digest. Apply requires
that digest through --confirm-plan and uses the recoverable Baseline
transaction to update selected skill files and skills-lock.json atomically.
An empty restoration is an idempotent exit 0.

Exit codes:
  0  selected skills already match, or the confirmed restoration was applied
  1  source acquisition, proof, apply, output, rollback, or recovery failure
  2  invalid arguments, profile, skill, lock schema, source, or unsafe target
  3  confirmation is required or does not match the current Change Plan
  130 operation canceled

Options:
  --profile       Required built-in Baseline Profile
  --skill         External profile skill to restore; repeat to select multiple,
                  or omit to restore every drifted external profile skill
  --source-dir    Declared offline Git checkout or bare object store containing
                  every selected skill's exact immutable commit
  --confirm-plan  Exact lowercase Plan Digest returned by the current preview
  --repo          Git worktree or a path inside it (default current directory)
  --format        Output format: text or json (default text)
`
	case "baseline assets sync":
		return `Usage:
  roundfix baseline assets sync --source-dir <path> [--check] [--format <text|json>]

Checks or refreshes the Go-owned canonical Baseline setup snapshots from an
explicit canonical setups directory. The source must be a clean Git checkout
with a portable GitHub origin, a full immutable HEAD commit, committed setup
documents, safe source-relative paths, and complete regular-file skill trees
whose working bytes match the declared commit.

Check mode is read-only and reports whether every canonical snapshot is
current. A non-empty refresh first validates the generated catalog in memory,
then uses the recoverable Baseline transaction to update only
internal/baseline/assets/setups. It never installs skills, writes the canonical
source, or reads the installed setup-context-driven skill at runtime.

Exit codes:
  0  snapshots are current, or a canonical refresh completed
  1  check-mode drift, refresh, output, rollback, or recovery failure
  2  invalid arguments, source provenance, path, tree, or catalog compatibility
  130 operation canceled

Options:
  --source-dir  Required canonical setups directory inside the immutable source checkout
  --check       Report snapshot drift without writing canonical assets
  --format      Output format: text or json (default text)
`
	case "profiles":
		return `Usage:
  roundfix profiles show [--category <category>] [--json]
  roundfix profiles configure --scope user|project [--file <path>] [--remove <category>] [--dry-run] [--yes] [--json]
  roundfix profiles validate [--category <category>] [--json]

Commands:
  show       Render effective Agent Selection Profiles and advisory recommendations.
  configure  Write complete Agent Selection Profiles after validation and confirmation.
  validate   Prove effective Agent Selection Profiles through disposable sessions.

Profile recommendations are advisory. They never route selections or mutate
configuration unless configure is explicitly invoked and confirmed.
`
	case "profiles show":
		return `Usage:
  roundfix profiles show [--category <category>] [--json]

Renders the effective Agent Selection Profile source, Preferred Selection,
Fallback Chain, and advisory top-five recommendations for one category or all
categories. Recommendations are read-only guidance and never change routing.

Options:
  --category  Agent Work Category: general, backend, frontend, data, infra, docs, test, chore, qa, or review
  --json      Print roundfix/profiles/v1 JSON
`
	case "profiles configure":
		return `Usage:
  roundfix profiles configure --scope user|project [--file <path>] [--remove <category>] [--dry-run] [--yes] [--json]

Adds or replaces complete Agent Selection Profiles in User Config or Project
Config. --remove declares a category removal and may be repeated. Without
--file or --remove, collects one profile through Interactive Input. The command
proves every exact Agent Selection it writes before confirmation. --dry-run
performs the same proof and change summary without writing.

Options:
  --scope    Required config scope: user or project
  --file     Strict profile fragment YAML; omitted opens Interactive Input
  --remove   Agent Work Category to remove; repeatable
  --dry-run  Validate and render the normalized result without writing
  --yes      Write without confirmation after validation
  --json     Print deterministic JSON report
`
	case "profiles validate":
		return `Usage:
  roundfix profiles validate [--category <category>] [--json]

Read-only validation resolves effective profiles, deduplicates exact Agent
Selections, proves each exact tuple through a disposable ACP Runtime session,
and closes every disposable session on success or error.

Options:
  --category  Agent Work Category: general, backend, frontend, data, infra, docs, test, chore, qa, or review
  --json      Print deterministic validation JSON
`
	case "archive":
		return archiveUsage
	case "spec check":
		return specCheckUsage
	case "spec audit":
		return specAuditUsage
	case "stop":
		return `Usage:
  roundfix stop <run-id>
  roundfix stop --run-id <id>
  roundfix stop --pr <number>
  roundfix stop --spec <slug>
  roundfix stop --head-repo <owner/name> --head-branch <branch>

Options:
  --run-id      Active Run ID to stop
  --run         Alias for --run-id
  --pr          Open Pull Request number used to find the Active Run
  --spec        Spec slug used to find the Active Run in the current repository
  --head-repo   Explicit Head Repository, owner/name
  --head-branch Explicit PR Head Branch
  --force       Force Stop a dead, stuck, or runaway Run after proving owner exit
  --owner-identity-unreadable
                Permit Force Stop only after owner identity proof fails as unreadable

Default stop is graceful: it records a Stop Request and the Run stops after
the current Work Item settles. During a Review Source wait, the owner observes
the request by the next configured poll boundary and runs no later fetch,
check, commit, push, or Review Source mutation.

Force Stop cancels registered active Agent Sessions, terminates the recorded
owner, and reports Stopped only after owner exit is proven. It then releases
the Active Run lock and may reap provably empty kept Worktrees and branches.
If exit proof fails, the Run remains Active; its Active Run lock stays retained.
Inspect it with 'roundfix runs list --state active', resolve the reported
owner-process failure, then retry 'roundfix stop --force <run-id>'.

Repeating Force Stop for an already Stopped Run reports the existing outcome
without repeating cleanup. A different terminal outcome is rejected and
preserved. Cleanup failures are secondary warnings after the primary failure.
`
	case "skills":
		return `Usage:
  roundfix skills list
  roundfix skills check
  roundfix skills install [--target <project|codex|claude|opencode|all>] [--dir <path>]

Commands:
  list       List the bundled Roundfix skills and recommended external skills
  check      Validate shipped Roundfix skill artifacts
  install    Install the bundled Roundfix skills into this project or compatible Agent skill directories

Options:
  --target   Install target. Supported: project, codex, claude, opencode, all
             project writes <repo>/.agents/skills and can link .claude/skills
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
	fmt.Fprintf(stderr, "%s\n", style.cyan("Reason:"))
	fmt.Fprintf(stderr, "  %v\n\n", err)
	printPreflightNoSideEffects(stderr, style)
	var actionable interface{ NextAction() string }
	if errors.As(err, &actionable) {
		if nextAction := strings.TrimSpace(actionable.NextAction()); nextAction != "" {
			fmt.Fprintf(stderr, "%s\n", style.cyan("Next action:"))
			fmt.Fprintf(stderr, "  %s\n\n", nextAction)
		}
	}
	fmt.Fprintf(stderr, "%s\n", style.cyan("Usage:"))
	fmt.Fprintf(stderr, "  Run '%s %s --help' for usage.\n", app.Name, name)
}

func printPreflightNoSideEffects(stderr io.Writer, style terminalStyle) {
	fmt.Fprintf(stderr, "%s\n", style.cyan("No side effects:"))
	fmt.Fprintln(stderr, "  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.")
	fmt.Fprintln(stderr)
}

func printBranchIntegrityAuditFailure(name string, runID string, err error, stderr io.Writer) {
	style := styleFor(stderr)
	fmt.Fprintf(stderr, "%s\n\n", style.red(style.bold("Preflight failed")))
	fmt.Fprintf(stderr, "%s\n", style.cyan("Reason:"))
	fmt.Fprintf(stderr, "  publish Branch Integrity Preflight bypass audit for Run %s: %v\n\n", runID, err)
	fmt.Fprintf(stderr, "%s\n", style.cyan("No further side effects:"))
	fmt.Fprintf(stderr, "  Roundfix created Run %s and attempted the mandatory audit comment.\n", runID)
	fmt.Fprintln(stderr, "  Roundfix did not fetch Review Source issues, start an Agent, commit, or push.")
	fmt.Fprintln(stderr)
	fmt.Fprintf(stderr, "%s\n", style.cyan("Next:"))
	fmt.Fprintln(stderr, "  Fix pull request comment publishing, or rerun without --skip-branch-integrity and resolve the Branch Integrity Preflight refusal.")
	fmt.Fprintln(stderr)
	fmt.Fprintf(stderr, "%s\n", style.cyan("Usage:"))
	fmt.Fprintf(stderr, "  Run '%s %s --help' for usage.\n", app.Name, name)
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

func printStopSuccess(result stopResult, stdout io.Writer) {
	run := result.Run
	style := styleFor(stdout)
	if result.Requested {
		fmt.Fprintf(stdout, "%s\n\n", style.green(style.bold("Roundfix Stop Request recorded")))
		fmt.Fprintf(stdout, "%s\n", style.cyan("Run:"))
		printStopRunFields(stdout, run)
		fmt.Fprintln(stdout)
		fmt.Fprintf(stdout, "%s\n", style.cyan("Next:"))
		fmt.Fprintln(stdout, "  Stop Request recorded; the Run stops after the current Work Item settles.")
		fmt.Fprintf(stdout, "  If the owning process is dead or runaway, run '%s stop --force %s'.\n\n", app.Name, run.ID)
		fmt.Fprintf(stdout, "%s\n", style.cyan("No repository side effects:"))
		fmt.Fprintln(stdout, "  Roundfix recorded the Stop Request in the Run Database only.")
		return
	}
	if result.Forced {
		fmt.Fprintf(stdout, "%s\n\n", style.green(style.bold("Roundfix Run force-stopped")))
		fmt.Fprintf(stdout, "%s\n", style.cyan("Run:"))
		printStopRunFields(stdout, run)
		fmt.Fprintln(stdout)
		fmt.Fprintf(stdout, "%s\n", style.cyan("Result:"))
		fmt.Fprintln(stdout, "  Force Stop proved the recorded owner process exited, completed the Run as Stopped, and released its Active Run locks.")
		fmt.Fprintln(stdout)
		fmt.Fprintf(stdout, "%s\n", style.cyan("No user work side effects:"))
		fmt.Fprintln(stdout, "  Roundfix did not edit user files, commit, push, fetch, or resolve Review Source threads.")
		return
	}
	fmt.Fprintf(stdout, "%s\n\n", style.green(style.bold("Roundfix Run stopped")))
	fmt.Fprintf(stdout, "%s\n", style.cyan("Run:"))
	printStopRunFields(stdout, run)
	fmt.Fprintln(stdout)
	fmt.Fprintf(stdout, "%s\n", style.cyan("No repository side effects:"))
	fmt.Fprintln(stdout, "  Roundfix released the Active Run lock without editing files, committing, pushing, fetching, or resolving Review Source threads.")
}

func printStopRunFields(stdout io.Writer, run store.Run) {
	fmt.Fprintf(stdout, "  ID: %s\n", run.ID)
	fmt.Fprintf(stdout, "  State: %s\n", run.State)
	fmt.Fprintf(stdout, "  Kind: %s\n", run.Kind)
	if run.PRNumber != "" || run.HeadRepository != "" || run.HeadBranch != "" {
		fmt.Fprintf(stdout, "  PR: #%s %s\n", run.PRNumber, run.HeadRepository)
		fmt.Fprintf(stdout, "  Branch: %s\n", run.HeadBranch)
	}
	if run.SpecSlug != "" {
		fmt.Fprintf(stdout, "  Spec: %s\n", run.SpecSlug)
	}
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
	printAgentCheckoutChangesNotice(stderr)
}

func printWatchRunFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: watch failed after Run start: %v\n", app.Name, err)
	fmt.Fprintln(stderr, "Roundfix completed the Watch Run with a terminal failure; inspect local artifacts and Run output before retrying.")
}

func printAgentCheckoutChangesNotice(stderr io.Writer) {
	fmt.Fprintln(stderr, "The tracked checkout was clean at Preflight; inspect current tracked and untracked changes before retrying.")
}
