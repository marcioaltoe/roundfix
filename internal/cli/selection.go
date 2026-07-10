package cli

import (
	"fmt"
	"strings"

	"roundfix/internal/agent"
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
	if reasoningEffort == "" {
		missing = append(missing, "reasoning_effort")
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
