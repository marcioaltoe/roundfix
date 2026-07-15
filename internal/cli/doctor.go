package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"

	"roundfix/internal/agent"
	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
)

var doctorDeps = defaultDoctorDependencies()

type doctorDependencies struct {
	loadConfig    func(roundconfig.LoadOptions) (roundconfig.Loaded, error)
	healthChecker func(roundconfig.Loaded) HealthChecker
}

func defaultDoctorDependencies() doctorDependencies {
	return doctorDependencies{
		loadConfig: roundconfig.Load,
		healthChecker: func(roundconfig.Loaded) HealthChecker {
			return setupDeps.healthChecker()
		},
	}
}

func runDoctorCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("doctor"))
		return exitOK
	}
	if err := parseDoctorCommand(args); err != nil {
		printDoctorFailure(err, stderr)
		return exitPreflight
	}

	loaded, err := doctorDeps.loadConfig(roundconfig.LoadOptions{Stderr: stderr})
	if err != nil {
		printDoctorFailure(err, stderr)
		return exitRunFailed
	}

	checker := doctorDeps.healthChecker(loaded)
	runtime, runtimeErr := runtimeForConfiguredAgent(loaded.Config)
	adapterResult := doctorAdapterCheck(ctx, checker, runtime, runtimeErr)
	var agentResult CheckResult
	if runtimeErr == nil && adapterResult.Status == CheckStatusFailed {
		agentResult = CheckResult{
			Name:   HealthCheckAgent,
			Status: CheckStatusSkipped,
			Detail: "adapter failed",
		}
	} else {
		agentResult = doctorAgentCheck(ctx, checker, loaded)
	}
	results := []CheckResult{
		checker.Node(ctx),
		checker.ACPX(ctx),
		adapterResult,
		agentResult,
		checker.Codex(ctx),
	}

	failed := false
	for _, result := range results {
		printDoctorResult(stdout, result)
		if result.Status == CheckStatusFailed {
			failed = true
		}
	}
	if failed {
		return exitRunFailed
	}
	return exitOK
}

func parseDoctorCommand(args []string) error {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	if err := fs.Parse(args); err != nil {
		return validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	return nil
}

func doctorAdapterCheck(ctx context.Context, checker HealthChecker, runtime agent.RuntimeSpec, runtimeErr error) CheckResult {
	if runtimeErr != nil {
		return CheckResult{
			Name:   HealthCheckAdapter,
			Status: CheckStatusFailed,
			Detail: runtimeErr.Error(),
		}
	}
	return checker.Adapter(ctx, runtime)
}

func doctorAgentCheck(ctx context.Context, checker HealthChecker, loaded roundconfig.Loaded) CheckResult {
	runtime, err := runtimeForConfiguredAgent(loaded.Config)
	if err != nil {
		return CheckResult{
			Name:   HealthCheckAgent,
			Status: CheckStatusFailed,
			Detail: err.Error(),
		}
	}
	return checker.Agent(ctx, agent.ProbeRequest{
		Runtime: runtime,
		WorkDir: loaded.GitRoot,
	})
}

func printDoctorResult(stdout io.Writer, result CheckResult) {
	detail := strings.TrimSpace(result.Detail)
	if result.Status == CheckStatusFailed {
		nextAction := strings.TrimSpace(result.NextAction)
		if nextAction != "" {
			if detail == "" {
				detail = "next: " + nextAction
			} else {
				detail += "; next: " + nextAction
			}
		}
	}
	if detail == "" {
		fmt.Fprintf(stdout, "%s: %s\n", result.Name, result.Status)
		return
	}
	fmt.Fprintf(stdout, "%s: %s (%s)\n", result.Name, result.Status, detail)
}

func printDoctorFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: doctor failed: %v\n", app.Name, err)
	fmt.Fprintf(stderr, "Run '%s doctor --help' for usage.\n", app.Name)
}
