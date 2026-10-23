package lighttier

import (
	"errors"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
)

// Selections returns one light Agent Selection per model in input order.
func Selections(models []string) []roundconfig.AgentSelection {
	selections := make([]roundconfig.AgentSelection, 0, len(models))
	for _, model := range models {
		selections = append(selections, roundconfig.AgentSelection{Runtime: "opencode", Model: "openrouter/" + model})
	}
	return selections
}

// AccessUnavailable describes a runtime that cannot serve its requested policy.
type AccessUnavailable struct {
	Selection roundconfig.AgentSelection
	Err       *agent.AccessPolicyError
}

// ServingSelections keeps factory errors for session activation to handle.
// Only a runtime's declared lack of an access mode excludes a selection.
func ServingSelections(models []string, runtimeFor func(roundconfig.AgentSelection) (agent.RuntimeSpec, error)) ([]roundconfig.AgentSelection, []AccessUnavailable) {
	var kept []roundconfig.AgentSelection
	var refused []AccessUnavailable
	for _, selection := range Selections(models) {
		runtime, err := runtimeFor(selection)
		if err == nil {
			var accessErr *agent.AccessPolicyError
			if errors.As(runtime.ValidateRequestedAccessPolicy(), &accessErr) && accessErr.Predicate == agent.AccessModeAvailablePredicate {
				refused = append(refused, AccessUnavailable{Selection: selection, Err: accessErr})
				continue
			}
		}
		kept = append(kept, selection)
	}
	return kept, refused
}
