package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"regexp"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/phaserecord"
)

const phaseUsage = `Usage:
  roundfix phase start <phase> [--spec <slug>]

Commands:
  start  Record in the Phase Record that this session starts a phase

Phases:
  follow-up, idea, prd, techspec, tasks, delivery, release

Options:
  --spec  Spec folder slug the phase works on; required for every phase
          but follow-up and release
`

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
	entry, err := parsePhaseStart(args)
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
		fmt.Fprintf(stdout, " for claude session %s\n", entry.Session)
	}
	return exitOK
}

func parsePhaseStart(args []string) (phaserecord.PhaseEntry, error) {
	if len(args) == 0 {
		return phaserecord.PhaseEntry{}, fmt.Errorf("missing subcommand; use start")
	}
	if args[0] != "start" {
		return phaserecord.PhaseEntry{}, fmt.Errorf("unknown subcommand %q; use start", args[0])
	}
	if len(args) < 2 {
		return phaserecord.PhaseEntry{}, fmt.Errorf("missing phase; use follow-up, idea, prd, techspec, tasks, delivery or release")
	}
	phase, err := phaserecord.ParsePhase(args[1])
	if err != nil {
		return phaserecord.PhaseEntry{}, err
	}
	var spec string
	flags := flag.NewFlagSet("phase start", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&spec, "spec", "", "Spec folder slug")
	if err := flags.Parse(args[2:]); err != nil {
		return phaserecord.PhaseEntry{}, err
	}
	if flags.NArg() > 0 {
		return phaserecord.PhaseEntry{}, fmt.Errorf("unexpected argument %q", flags.Arg(0))
	}
	if spec == "" && phase != "follow-up" && phase != "release" {
		return phaserecord.PhaseEntry{}, fmt.Errorf("phase %q requires --spec <slug>", phase)
	}
	specSet := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "spec" {
			specSet = true
		}
	})
	if specSet && !phaseSpecSlug.MatchString(spec) {
		return phaserecord.PhaseEntry{}, fmt.Errorf("invalid Spec slug %q; --spec must match ^[0-9]{4}(-[a-z0-9]+)+$", spec)
	}
	return phaserecord.PhaseEntry{Phase: phase, Spec: spec}, nil
}
