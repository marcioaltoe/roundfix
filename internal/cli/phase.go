package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"regexp"
	"strings"

	"roundfix/internal/baseline"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/phaserecord"
)

const phaseUsage = `Usage:
  roundfix phase start <phase> [--spec <slug>] [--model <identifier>] [--effort <level>] [--author-subagent]
  roundfix phase model <phase> --spec <slug>
  roundfix phase fallback <phase> --spec <slug> --error <message>

Commands:
  start     Record in the Phase Record that this session starts a phase
  model     Print the model and effort a phase's Author Subagent starts on
  fallback  Decide and record an Authoring Fallback from the error an Author Subagent ended on

Phases:
  follow-up, idea, prd, techspec, tasks, delivery, release
  model and fallback take idea, prd, techspec or tasks

Options:
  --spec             Spec folder slug the phase works on; required for every phase
                     but follow-up and release
  --model            Model identifier the harness states for this session; recorded
                     for idea, prd, techspec and tasks
  --effort           Reasoning effort; used when CLAUDE_EFFORT is unset
  --author-subagent  The phase's Author Subagent runs the phase
  --error            Error the Author Subagent ended on; required for fallback
`

// Keep the identifier grammar aligned with the Authoring Model Declaration.
var phaseModelName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/@\[\]-]{0,199}$`)
var phaseEffortLevel = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,31}$`)

var phaseSpecSlug = regexp.MustCompile(`^[0-9]{4}(-[a-z0-9]+)+$`)

func harnessSession(lookup func(string) string) (harness, session string) {
	session = lookup("CLAUDE_CODE_SESSION_ID")
	if session == "" {
		return "", ""
	}
	return "claude", session
}

func runPhaseCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, phaseUsage)
		return exitOK
	}
	if len(args) > 0 && (args[0] == "model" || args[0] == "fallback") {
		return runPhaseAuthoringCommand(ctx, args, stdout, stderr, environment)
	}
	options, err := parsePhaseStart(args)
	entry := options.PhaseEntry
	if err == nil {
		var workDir string
		workDir, err = environment.resolveWorkDir("resolve working directory")
		if err == nil {
			entry.Repository = roundconfig.FindGitRoot(workDir)
			if entry.Repository == "" {
				err = fmt.Errorf("phase start requires a git repository working tree")
			}
		}
	}
	if err != nil {
		fmt.Fprintf(stderr, "roundfix: phase start failed: %v\nRun 'roundfix phase --help' for usage.\n", err)
		return exitPreflight
	}
	if environment.homeDirErr != nil {
		fmt.Fprintf(stderr, "roundfix: phase start failed: resolve home directory: %v\n", environment.homeDirErr)
		return exitRunFailed
	}
	entry.Harness, entry.Session = harnessSession(environment.dependencies.getenv)
	if entry.Phase.IsAuthoring() {
		// Identity is best effort and depends only on local Git objects, never a remote.
		if identity, identityErr := baseline.InspectRepositoryIdentity(ctx, entry.Repository, nil); identityErr == nil {
			entry.RepositoryID = identity.Digest
		}
		entry.Model, entry.ModelSource, entry.Effort, entry.EffortSource, entry.Runner = "", "unknown", "", "unknown", "unknown"
		if entry.Session != "" {
			entry.Runner = "session"
			if options.authorSubagent {
				entry.Runner = "author-subagent"
			}
			if options.modelSet {
				if phaseModelName.MatchString(options.Model) {
					entry.Model, entry.ModelSource = options.Model, "reported"
				} else {
					fmt.Fprintln(stderr, "roundfix: warning: phase start: --model is not a model identifier of at most 200 characters; the model is recorded as unknown")
				}
			}
			effort, source, label := environment.dependencies.getenv("CLAUDE_EFFORT"), "environment", "CLAUDE_EFFORT"
			if effort == "" {
				effort, source, label = options.Effort, "reported", "--effort"
			}
			if effort != "" || options.effortSet && label == "--effort" {
				if phaseEffortLevel.MatchString(effort) {
					entry.Effort, entry.EffortSource = effort, source
				} else {
					fmt.Fprintf(stderr, "roundfix: warning: phase start: %s is not an effort level; the effort is recorded as unknown\n", label)
				}
			}
		}
	} else if options.modelFlags {
		fmt.Fprintf(stderr, "roundfix: warning: phase start: --model, --effort and --author-subagent apply to idea, prd, techspec and tasks; ignored for %s\n", entry.Phase)
	}
	if err := phaserecord.New(environment.homeDir, environment.dependencies.phaseNow).StartPhase(ctx, entry); err != nil {
		fmt.Fprintf(stderr, "roundfix: phase start failed: %v\n", err)
		return exitRunFailed
	}
	fmt.Fprintf(stdout, "Recorded %s", entry.Phase)
	if entry.Spec != "" {
		fmt.Fprintf(stdout, " of %s", entry.Spec)
	}
	if entry.Session == "" {
		fmt.Fprintln(stdout, " without a harness session: the environment names none")
	} else {
		fmt.Fprintf(stdout, " for claude session %s", entry.Session)
		if entry.Phase.IsAuthoring() {
			runner := "the session"
			if entry.Runner == "author-subagent" {
				runner = "the Author Subagent"
			}
			fmt.Fprintf(stdout, ": model %s, effort %s, run by %s", phaseValue(entry.Model, entry.ModelSource), phaseValue(entry.Effort, entry.EffortSource), runner)
		}
		fmt.Fprintln(stdout)
	}
	return exitOK
}

func phaseValue(value, source string) string {
	if source == "unknown" {
		return "unknown"
	}
	return value + " (" + source + ")"
}

type phaseStartOptions struct {
	phaserecord.PhaseEntry
	modelSet, effortSet, authorSubagent, modelFlags bool
}

func parsePhaseStart(args []string) (phaseStartOptions, error) {
	if len(args) == 0 {
		return phaseStartOptions{}, fmt.Errorf("missing subcommand; use start")
	}
	if args[0] != "start" {
		return phaseStartOptions{}, fmt.Errorf("unknown subcommand %q; use start", args[0])
	}
	if len(args) < 2 {
		return phaseStartOptions{}, fmt.Errorf("missing phase; use follow-up, idea, prd, techspec, tasks, delivery or release")
	}
	phase, err := phaserecord.ParsePhase(args[1])
	if err != nil {
		return phaseStartOptions{}, err
	}
	var spec string
	var options phaseStartOptions
	flags := flag.NewFlagSet("phase start", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&spec, "spec", "", "Spec folder slug")
	flags.StringVar(&options.Model, "model", "", "Model identifier")
	flags.StringVar(&options.Effort, "effort", "", "Reasoning effort")
	flags.BoolVar(&options.authorSubagent, "author-subagent", false, "The phase's Author Subagent runs the phase")
	if err := flags.Parse(args[2:]); err != nil {
		return phaseStartOptions{}, err
	}
	if flags.NArg() > 0 {
		return phaseStartOptions{}, fmt.Errorf("unexpected argument %q", flags.Arg(0))
	}
	if spec == "" && phase != "follow-up" && phase != "release" {
		return phaseStartOptions{}, fmt.Errorf("phase %q requires --spec <slug>", phase)
	}
	specSet := false
	flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "spec":
			specSet = true
		case "model":
			options.modelSet, options.modelFlags = true, true
		case "effort":
			options.effortSet, options.modelFlags = true, true
		case "author-subagent":
			options.modelFlags = true
		}
	})
	if specSet && !phaseSpecSlug.MatchString(spec) {
		return phaseStartOptions{}, fmt.Errorf("invalid Spec slug %q; --spec must match ^[0-9]{4}(-[a-z0-9]+)+$", spec)
	}
	options.Phase, options.Spec = phase, spec
	return options, nil
}

// parsePhaseAuthoring validates both read-only model requests and fallback writes.
func parsePhaseAuthoring(args []string) (phase phaserecord.Phase, spec, message string, err error) {
	if len(args) < 2 {
		return "", "", "", fmt.Errorf("missing phase; use idea, prd, techspec or tasks")
	}
	phase = phaserecord.Phase(args[1])
	if !phase.IsAuthoring() {
		return "", "", "", fmt.Errorf("phase %q has no Author Subagent; use idea, prd, techspec or tasks", args[1])
	}
	flags := flag.NewFlagSet("phase "+args[0], flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&spec, "spec", "", "Spec folder slug")
	if args[0] == "fallback" {
		flags.StringVar(&message, "error", "", "Error the Author Subagent ended on")
	}
	if err := flags.Parse(args[2:]); err != nil {
		return "", "", "", err
	}
	if flags.NArg() > 0 {
		return "", "", "", fmt.Errorf("unexpected argument %q", flags.Arg(0))
	}
	if spec == "" {
		return "", "", "", fmt.Errorf("phase %q requires --spec <slug>", phase)
	}
	if !phaseSpecSlug.MatchString(spec) {
		return "", "", "", fmt.Errorf("invalid Spec slug %q; --spec must match ^[0-9]{4}(-[a-z0-9]+)+$", spec)
	}
	if args[0] == "fallback" && (strings.TrimSpace(message) == "" || len(message) > 4096) {
		return "", "", "", fmt.Errorf("--error must be non-empty and at most 4096 bytes")
	}
	return phase, spec, message, nil
}

func runPhaseAuthoringCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	command := args[0]
	phase, spec, message, err := parsePhaseAuthoring(args)
	var root string
	var loaded roundconfig.Loaded
	if err == nil {
		var workDir string
		workDir, err = environment.resolveWorkDir("resolve working directory")
		if err == nil {
			root = roundconfig.FindGitRoot(workDir)
			if root == "" {
				err = fmt.Errorf("phase %s requires a git repository working tree", command)
			}
		}
	}
	if err == nil {
		environment.workDir, environment.workDirErr = root, nil
		loaded, err = loadCommandConfig(environment, stderr)
	}
	if err != nil {
		fmt.Fprintf(stderr, "roundfix: phase %s failed: %v\nRun 'roundfix phase --help' for usage.\n", command, err)
		return exitPreflight
	}
	if environment.homeDirErr != nil {
		fmt.Fprintf(stderr, "roundfix: phase %s failed: resolve home directory: %v\n", command, environment.homeDirErr)
		return exitRunFailed
	}
	var entry roundconfig.AuthoringEntry
	if h := loaded.Config.Skills.Authoring.Claude; h != nil {
		var declared *roundconfig.AuthoringEntry
		switch phase {
		case "idea":
			declared = h.Idea
		case "prd":
			declared = h.PRD
		case "techspec":
			declared = h.TechSpec
		case "tasks":
			declared = h.Tasks
		}
		if declared != nil {
			entry = *declared
		}
	}
	repositoryID := ""
	// As for start, an unborn repository can be matched by its Git root.
	if identity, identityErr := baseline.InspectRepositoryIdentity(ctx, root, nil); identityErr == nil {
		repositoryID = identity.Digest
	}
	if command == "model" {
		return printPhaseModel(phase, spec, root, repositoryID, entry, environment.homeDir, stdout, stderr)
	}
	limit := phaserecord.ReadUsageLimit(message)
	reason := ""
	switch {
	case limit.Kind == "none":
		reason = "the error is not a usage limit"
	case limit.Kind == "shared":
		reason = fmt.Sprintf("a %s limit is shared by every model", limit.Name)
	case limit.Kind == "other":
		reason = fmt.Sprintf("a %s limit is not a model's usage limit", limit.Name)
	case entry.Model == "":
		reason = fmt.Sprintf("%s declares no model", phase)
	case !phaseLimitMatches(limit.Name, entry.Model):
		reason = fmt.Sprintf("the %s limit is not that of the declared model %s", limit.Name, entry.Model)
	case entry.Fallback == nil:
		reason = fmt.Sprintf("%s declares no fallback", phase)
	case phaseLimitMatches(limit.Name, entry.Fallback.Model):
		reason = fmt.Sprintf("the fallback %s is a %s model too", entry.Fallback.Model, limit.Name)
	}
	if reason != "" {
		fmt.Fprintf(stdout, "No fallback for %s of %s: %s.\n", phase, spec, reason)
		return 3
	}
	harness, session := harnessSession(environment.dependencies.getenv)
	row := phaserecord.FallbackEntry{Phase: phase, Spec: spec, Repository: root, RepositoryID: repositoryID,
		Harness: harness, Session: session, DeclaredModel: entry.Model, DeclaredEffort: entry.ReasoningEffort,
		Model: entry.Fallback.Model, Effort: entry.Fallback.ReasoningEffort, Limit: limit.Name, Message: limit.Line}
	if err := phaserecord.New(environment.homeDir, environment.dependencies.phaseNow).RecordFallback(ctx, row); err != nil {
		fmt.Fprintf(stderr, "roundfix: phase fallback failed: %v\n", err)
		return exitRunFailed
	}
	fmt.Fprintf(stdout, "%s: \"%s\"\n", phaseFallbackStart(phase, entry), limit.Line)
	return exitOK
}

func phaseLimitMatches(name, model string) bool {
	return name != "" && strings.Contains(strings.ToLower(model), strings.ToLower(name))
}

func phaseFallbackStart(phase phaserecord.Phase, entry roundconfig.AuthoringEntry) string {
	effort := "no effort"
	if entry.Fallback.ReasoningEffort != "" {
		effort = "effort " + entry.Fallback.ReasoningEffort
	}
	return fmt.Sprintf("Start roundfix-%s-author with model %s and %s, the fallback of %s", phase, entry.Fallback.Model, effort, entry.Model)
}

func printPhaseModel(phase phaserecord.Phase, spec, root, repositoryID string, entry roundconfig.AuthoringEntry, home string, stdout, stderr io.Writer) int {
	if entry.Model == "" {
		fmt.Fprintf(stdout, "No Author Subagent model for %s: the Authoring Model Declaration declares no model for it.\n", phase)
		return exitOK
	}
	rows, _, err := phaserecord.Fallbacks(home)
	if err != nil {
		fmt.Fprintf(stderr, "roundfix: phase model failed: %v\n", err)
		return exitRunFailed
	}
	for _, row := range rows {
		sameRepository := row.RepositoryID == repositoryID
		if row.RepositoryID == "" || repositoryID == "" {
			sameRepository = row.Repository == root
		}
		if row.Spec != spec || !sameRepository || !phaseLimitMatches(row.Limit, entry.Model) || entry.Fallback == nil || phaseLimitMatches(row.Limit, entry.Fallback.Model) {
			continue
		}
		fmt.Fprintf(stdout, "%s: %s of %s fell back when the %s limit ran out.\n", phaseFallbackStart(phase, entry), row.Phase, spec, row.Limit)
		return exitOK
	}
	fmt.Fprintf(stdout, "Start roundfix-%s-author with no model and no effort: the declared model %s.\n", phase, entry.Model)
	return exitOK
}
