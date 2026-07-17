package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	"roundfix/internal/spec"
)

type profileOperationalPreflightResult struct {
	Proofs   []profileProofReport        `json:"proofs,omitempty"`
	Warnings []string                    `json:"warnings,omitempty"`
	Override *roundconfig.AgentSelection `json:"override,omitempty"`
	Err      error                       `json:"-"`
}

func implementProfileCategories(graph *spec.Graph, includeQA bool) []roundconfig.WorkCategory {
	present := map[roundconfig.WorkCategory]bool{}
	if graph != nil {
		for _, task := range graph.Tasks {
			if task.Status == spec.StatusCompleted {
				continue
			}
			present[workCategoryForTaskType(task.Type)] = true
		}
	}
	categories := make([]roundconfig.WorkCategory, 0, len(present)+1)
	for _, taskType := range spec.AllowedTaskTypes() {
		category := workCategoryForTaskType(taskType)
		if present[category] {
			categories = append(categories, category)
		}
	}
	if includeQA {
		categories = append(categories, roundconfig.CategoryQA)
	}
	return categories
}

func reviewProfileCategories() []roundconfig.WorkCategory {
	return []roundconfig.WorkCategory{roundconfig.CategoryReview}
}

func workCategoryForTaskType(taskType spec.TaskType) roundconfig.WorkCategory {
	return roundconfig.WorkCategory(taskType)
}

func runProfileOperationalPreflight(ctx context.Context, req commandRequest, config roundconfig.Config, categories []roundconfig.WorkCategory, workDir string, runner agent.Runner, stderr io.Writer) (profileOperationalPreflightResult, error) {
	override, err := invocationProfileOverride(req, config)
	if err != nil {
		return profileOperationalPreflightResult{Err: err}, err
	}
	warnings := profileOperationalWarnings(override, categories)
	for _, warning := range warnings {
		fmt.Fprintln(stderr, warning)
	}
	result := proveProfileSelectionsWithOptions(ctx, config, categories, workDir, runner, profileProofOptions{
		PreferredOverride: override,
		CommandOverride:   req.agentCmd,
		EnableFullAccess:  req.agentFullAccess,
	})
	if result.Err != nil {
		return profileOperationalPreflightResult{
			Proofs:   result.Proofs,
			Warnings: warnings,
			Override: override,
			Err:      result.Err,
		}, result.Err
	}
	return profileOperationalPreflightResult{
		Proofs:   result.Proofs,
		Warnings: warnings,
		Override: override,
	}, nil
}

func invocationProfileOverride(req commandRequest, config roundconfig.Config) (*roundconfig.AgentSelection, error) {
	if !req.agentSet && !req.modelSet && !req.reasoningEffortSet {
		return nil, nil
	}
	runtime, err := runtimeForAgentWork(req, config)
	if err != nil {
		return nil, err
	}
	return &roundconfig.AgentSelection{
		Runtime:         strings.TrimSpace(req.agent),
		Model:           strings.TrimSpace(runtime.Model),
		ReasoningEffort: strings.TrimSpace(runtime.ReasoningEffort),
	}, nil
}

func profileOperationalWarnings(override *roundconfig.AgentSelection, categories []roundconfig.WorkCategory) []string {
	if override == nil || len(categories) < 2 {
		return nil
	}
	return []string{fmt.Sprintf(
		"roundfix: warning: one invocation Agent Selection override applies to Agent Work Categories %s; configured Fallback Chains are preserved.",
		formatWorkCategories(categories),
	)}
}

func runtimeForOperationalProfileRun(req commandRequest, config roundconfig.Config, categories []roundconfig.WorkCategory, override *roundconfig.AgentSelection) (agent.RuntimeSpec, error) {
	if len(categories) == 0 {
		return runtimeForAgentWork(req, config)
	}
	category := categories[0]
	if req.name == "implement" {
		category = roundconfig.CategoryGeneral
	}
	resolved, err := roundconfig.ResolveProfile(config, category, override)
	if err != nil {
		return agent.RuntimeSpec{}, err
	}
	return runtimeForProfileSelectionWithOptions(resolved.Profile.Preferred, profileProofOptions{
		CommandOverride:  req.agentCmd,
		EnableFullAccess: req.agentFullAccess,
	})
}

func operationalAgentSelectionProfiles(config roundconfig.Config, categories []roundconfig.WorkCategory, override *roundconfig.AgentSelection) (daemon.AgentSelectionProfiles, error) {
	profiles := daemon.AgentSelectionProfiles{}
	for _, category := range categories {
		resolved, err := roundconfig.ResolveProfile(config, category, override)
		if err != nil {
			return nil, err
		}
		profiles[category] = resolved
	}
	return profiles, nil
}

func operationalRuntimeFactory(req commandRequest) daemon.AgentRuntimeFactory {
	return func(selection roundconfig.AgentSelection) (agent.RuntimeSpec, error) {
		return runtimeForProfileSelectionWithOptions(selection, profileProofOptions{
			CommandOverride:  req.agentCmd,
			EnableFullAccess: req.agentFullAccess,
		})
	}
}

func formatWorkCategories(categories []roundconfig.WorkCategory) string {
	values := make([]string, 0, len(categories))
	for _, category := range categories {
		values = append(values, string(category))
	}
	return strings.Join(values, ", ")
}
