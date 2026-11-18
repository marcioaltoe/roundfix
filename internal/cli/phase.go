package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"regexp"

	"roundfix/internal/baseline"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/phaserecord"
)

const phaseUsage = `Usage:
  roundfix phase start <phase> [--spec <slug>] [--model <identifier>] [--effort <level>] [--author-subagent]

Commands:
  start  Record in the Phase Record that this session starts a phase

Phases:
  follow-up, idea, prd, techspec, tasks, delivery, release

Options:
  --spec             Spec folder slug the phase works on; required for every phase
                     but follow-up and release
  --model            Model identifier the harness states for this session; recorded
                     for idea, prd, techspec and tasks
  --effort           Reasoning effort; used when CLAUDE_EFFORT is unset
  --author-subagent  The phase's Author Subagent runs the phase
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
