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
	modelResult := doctorModelCheck(runtime, runtimeErr, adapterResult, agentResult)
	results := []CheckResult{
		checker.Node(ctx),
		checker.ACPX(ctx),
		adapterResult,
		agentResult,
		modelResult,
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
			Err:    err,
		}
	}
	return checker.Agent(ctx, agent.ProbeRequest{
		Runtime: runtime,
		WorkDir: loaded.GitRoot,
	})
}

func doctorModelCheck(runtime agent.RuntimeSpec, runtimeErr error, adapterResult CheckResult, agentResult CheckResult) CheckResult {
	// Doctor renders this check as the public "model:" line.
	result := CheckResult{Name: HealthCheckModel}
	if runtimeErr != nil {
		result.Status = CheckStatusSkipped
		result.Detail = "agent selection failed"
		return result
	}
	if adapterResult.Status == CheckStatusFailed {
		result.Status = CheckStatusSkipped
		result.Detail = "adapter failed"
		return result
	}
	if agentResult.Status == CheckStatusOK {
		result.Status = CheckStatusOK
		result.Detail = strings.TrimSpace(runtime.Model)
		return result
	}
	var modelErr *agent.ModelNotAdvertisedError
	if errors.As(agentResult.Err, &modelErr) {
		result.Status = CheckStatusFailed
		result.Detail = doctorModelNotAdvertisedDetail(runtime, modelErr)
		result.NextAction = modelErr.RecoveryAction()
		return result
	}
	result.Status = CheckStatusSkipped
	result.Detail = "agent probe failed"
	return result
}

func doctorModelNotAdvertisedDetail(runtime agent.RuntimeSpec, modelErr *agent.ModelNotAdvertisedError) string {
	model := strings.TrimSpace(runtime.Model)
	runtimeID := strings.TrimSpace(runtime.ID)
	if modelErr != nil {
		if rejected := strings.TrimSpace(modelErr.Model); rejected != "" {
			model = rejected
		}
		if rejectedRuntime := strings.TrimSpace(modelErr.Runtime); rejectedRuntime != "" {
			runtimeID = rejectedRuntime
		}
	}
	return fmt.Sprintf("Agent Model %q not advertised by runtime %q; advertised: %s", model, runtimeID, doctorAdvertisedModels(modelErr))
}

func doctorAdvertisedModels(modelErr *agent.ModelNotAdvertisedError) string {
	if modelErr == nil {
		return "unavailable"
	}
	advertised := make([]string, 0, len(modelErr.Advertised))
	for _, model := range modelErr.Advertised {
		model = strings.TrimSpace(model)
		if model != "" {
			advertised = append(advertised, model)
		}
	}
	if len(advertised) == 0 {
		return "unavailable"
	}
	return strings.Join(advertised, ", ")
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
