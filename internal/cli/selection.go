package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"roundfix/internal/agent"
	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
)

type AgentSelection struct {
	Runtime         string
	Model           string
	ReasoningEffort string
}

type InvocationSelection struct {
	Model              string
	ReasoningEffort    string
	ModelSet           bool
	ReasoningEffortSet bool
}

type fallbackProber interface {
	ProbeFallback(context.Context, agent.RuntimeSpec, agent.FallbackCandidateSet) (agent.FallbackSelection, bool, error)
}

var fallbackConfirmationInput = func() io.Reader { return os.Stdin }
var fallbackConfirmationAvailable = func(stderr io.Writer) bool { return liveTUIEnabled(stderr) }

type selectionFallbackReport struct {
	failure    *agent.SelectionPreflightError
	fallback   agent.FallbackSelection
	candidates agent.FallbackCandidateSet
	rerun      string
	probeErr   error
}

func (report selectionFallbackReport) Error() string {
	var text strings.Builder
	text.WriteString(report.failure.Error())
	writeFailedSelection(&text, report.failure)
	if report.probeErr != nil {
		fmt.Fprintf(&text, "\nFallback probe failed: %v", report.probeErr)
		writeProbedFallbackCandidates(&text, report.candidates)
		return text.String()
	}
	if report.rerun == "" {
		text.WriteString("\nFallback probe found no functional selection.")
		writeProbedFallbackCandidates(&text, report.candidates)
		return text.String()
	}
	writeFallbackSelection(&text, report.fallback)
	fmt.Fprintf(&text, "Re-run: %s", report.rerun)
	return text.String()
}

func (report selectionFallbackReport) Unwrap() error {
	return report.failure
}

func writeProbedFallbackCandidates(text *strings.Builder, candidates agent.FallbackCandidateSet) {
	models := strings.Join(candidates.Models, ", ")
	if models == "" {
		models = "(none)"
	}
	efforts := strings.Join(candidates.Efforts, ", ")
	if efforts == "" {
		efforts = "(model-managed only)"
	}
	fmt.Fprintf(text, "\nProbed Agent Models (newest first): %s", models)
	fmt.Fprintf(text, "\nProbed reasoning efforts (highest first): %s", efforts)
}

func writeFailedSelection(text *strings.Builder, failure *agent.SelectionPreflightError) {
	text.WriteString("\nFailed Selection:\n")
	fmt.Fprintf(text, "  ACP Runtime: %s\n", failure.Runtime)
	fmt.Fprintf(text, "  Agent Model: %s\n", failure.Model)
	fmt.Fprintf(text, "  Default Reasoning Effort: %s", displayReasoningEffort(failure.ReasoningEffort))
}

func writeFallbackSelection(text *strings.Builder, fallback agent.FallbackSelection) {
	text.WriteString("\nFallback Selection:\n")
	fmt.Fprintf(text, "  Agent Model: %s\n", fallback.Model)
	fmt.Fprintf(text, "  Default Reasoning Effort: %s\n", displayReasoningEffort(fallback.ReasoningEffort))
}

func probeRuntimeSelection(ctx context.Context, req commandRequest, runtime agent.RuntimeSpec, workDir string, runner agent.Runner, stderr io.Writer) (agent.RuntimeSpec, error) {
	err := runner.Probe(ctx, agent.ProbeRequest{Runtime: runtime, WorkDir: workDir})
	if err == nil {
		return runtime, nil
	}
	var selectionErr *agent.SelectionPreflightError
	if !errors.As(err, &selectionErr) {
		return runtime, err
	}

	candidates := fallbackCandidates(runtime)
	prober, ok := runner.(fallbackProber)
	if !ok {
		return runtime, selectionFallbackReport{
			failure:    selectionErr,
			candidates: candidates,
			probeErr:   errors.New("agent runner does not support fallback probes"),
		}
	}
	fallback, functional, probeErr := prober.ProbeFallback(ctx, runtime, candidates)
	if probeErr != nil {
		return runtime, selectionFallbackReport{failure: selectionErr, candidates: candidates, probeErr: probeErr}
	}
	if !functional {
		return runtime, selectionFallbackReport{failure: selectionErr, candidates: candidates}
	}
	report := selectionFallbackReport{
		failure:    selectionErr,
		fallback:   fallback,
		candidates: candidates,
		rerun:      fallbackRerunCommand(req, fallback),
	}
	if !selectionFailureIsNonInteractive(req, stderr) {
		confirmed, confirmErr := confirmFallbackSelection(ctx, fallbackConfirmationInput(), stderr, selectionErr, fallback)
		if confirmErr != nil {
			return runtime, fmt.Errorf("confirm fallback selection: %w", confirmErr)
		}
		if confirmed {
			runtime.Model = fallback.Model
			runtime.ReasoningEffort = fallback.ReasoningEffort
			return runtime, nil
		}
	}
	return runtime, report
}

func fallbackCandidates(runtime agent.RuntimeSpec) agent.FallbackCandidateSet {
	failedModel := strings.TrimSpace(runtime.Model)
	catalog := agent.ModelCatalog(runtime.ID)
	models := make([]string, 0, len(catalog))
	for _, choice := range catalog {
		model := strings.TrimSpace(choice.Value)
		if model == "" || model == failedModel {
			continue
		}
		models = append(models, model)
	}
	choices := reasoningEffortChoices(runtime.ID)
	efforts := make([]string, 0, len(choices))
	for index := len(choices) - 1; index >= 0; index-- {
		if effort := strings.TrimSpace(choices[index]); effort != "" {
			efforts = append(efforts, effort)
		}
	}
	return agent.FallbackCandidateSet{Models: models, Efforts: efforts}
}

func selectionFailureIsNonInteractive(req commandRequest, stderr io.Writer) bool {
	return req.noInput || req.detach || req.detachChild != nil || !fallbackConfirmationAvailable(stderr)
}

func confirmFallbackSelection(ctx context.Context, stdin io.Reader, stderr io.Writer, failure *agent.SelectionPreflightError, fallback agent.FallbackSelection) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	var report strings.Builder
	report.WriteString(failure.Error())
	writeFailedSelection(&report, failure)
	writeFallbackSelection(&report, fallback)
	if _, err := fmt.Fprint(stderr, report.String()); err != nil {
		return false, fmt.Errorf("write fallback report: %w", err)
	}
	if _, err := fmt.Fprintln(stderr, "A different Agent Model can consume tokens differently."); err != nil {
		return false, fmt.Errorf("write fallback caveat: %w", err)
	}
	if _, err := fmt.Fprint(stderr, "Use this Fallback Selection for this Run? [y/N]: "); err != nil {
		return false, fmt.Errorf("write fallback confirmation prompt: %w", err)
	}
	line, err := bufio.NewReader(stdin).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, fmt.Errorf("read fallback confirmation prompt: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	return normalizeFallbackConfirmation(line)
}

func normalizeFallbackConfirmation(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "y", "yes":
		return true, nil
	case "", "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("unsupported fallback response %q; enter yes or no", strings.TrimSpace(value))
	}
}

func fallbackRerunCommand(req commandRequest, fallback agent.FallbackSelection) string {
	args := withoutSelectionArgs(req.arguments)
	if req.detachChild != nil && !containsDetachArg(args) {
		args = append(args, "--detach")
	}
	args = append(args,
		"--agent", strings.TrimSpace(req.agent),
		"--model", strings.TrimSpace(fallback.Model),
		"--reasoning-effort", strings.TrimSpace(fallback.ReasoningEffort),
	)
	command := []string{app.Name, req.name}
	command = append(command, args...)
	quoted := make([]string, 0, len(command))
	for _, arg := range command {
		quoted = append(quoted, quoteCommandArg(arg))
	}
	return strings.Join(quoted, " ")
}

func withoutSelectionArgs(args []string) []string {
	result := make([]string, 0, len(args))
	for index := 0; index < len(args); index++ {
		arg := args[index]
		name, _, hasValue := strings.Cut(strings.TrimLeft(arg, "-"), "=")
		if strings.HasPrefix(arg, "-") && (name == "agent" || name == "model" || name == "reasoning-effort") {
			if !hasValue && index+1 < len(args) {
				index++
			}
			continue
		}
		result = append(result, arg)
	}
	return result
}

func containsDetachArg(args []string) bool {
	for _, arg := range args {
		if arg == "--detach" || arg == "-detach" || strings.HasPrefix(arg, "--detach=") || strings.HasPrefix(arg, "-detach=") {
			return true
		}
	}
	return false
}

func quoteCommandArg(value string) string {
	if value == "" {
		return "\"\""
	}
	for _, char := range value {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || strings.ContainsRune("-._/:=@+,%", char) {
			continue
		}
		return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
	}
	return value
}

func ResolveSelection(runtime string, defaults roundconfig.RuntimeDefaults, invocation InvocationSelection) (AgentSelection, error) {
	runtime = strings.TrimSpace(runtime)
	if err := validateAgent(runtime); err != nil {
		return AgentSelection{}, err
	}

	model := strings.TrimSpace(defaults.Model)
	if invocation.ModelSet {
		model = strings.TrimSpace(invocation.Model)
	}
	reasoningEffort := strings.TrimSpace(defaults.ReasoningEffort)
	if invocation.ReasoningEffortSet {
		reasoningEffort = strings.TrimSpace(invocation.ReasoningEffort)
	}

	missing := []string{}
	if model == "" {
		missing = append(missing, "model")
	}
	if len(missing) > 0 {
		return AgentSelection{}, fmt.Errorf("agent selection for runtime %q missing %s; configure runtimes.%s or pass invocation overrides", runtime, strings.Join(missing, " and "), runtime)
	}

	return AgentSelection{
		Runtime:         runtime,
		Model:           model,
		ReasoningEffort: reasoningEffort,
	}, nil
}

func runtimeForAgentWork(req commandRequest, config roundconfig.Config) (agent.RuntimeSpec, error) {
	defaults, _ := config.Runtimes.DefaultsFor(req.agent)
	invocation := InvocationSelection{
		Model:              req.model,
		ReasoningEffort:    req.reasoningEffort,
		ModelSet:           req.modelSet,
		ReasoningEffortSet: req.reasoningEffortSet,
	}
	selection, err := ResolveSelection(req.agent, defaults, invocation)
	if err != nil {
		return agent.RuntimeSpec{}, err
	}
	return agent.RuntimeFor(agent.RuntimeOptions{
		Agent:            selection.Runtime,
		CommandOverride:  req.agentCmd,
		Model:            selection.Model,
		ReasoningEffort:  selection.ReasoningEffort,
		EnableFullAccess: req.agentFullAccess,
	})
}

func runtimeForConfiguredAgent(config roundconfig.Config) (agent.RuntimeSpec, error) {
	agentID := strings.TrimSpace(config.Defaults.Agent)
	defaults, _ := config.Runtimes.DefaultsFor(agentID)
	selection, err := ResolveSelection(agentID, defaults, InvocationSelection{})
	if err != nil {
		return agent.RuntimeSpec{}, err
	}
	return agent.RuntimeFor(agent.RuntimeOptions{
		Agent:            selection.Runtime,
		Model:            selection.Model,
		ReasoningEffort:  selection.ReasoningEffort,
		EnableFullAccess: config.Defaults.AgentFullAccess,
	})
}
