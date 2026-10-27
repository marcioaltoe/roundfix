package baseline

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// carryRecordedDecisions takes only recorded answers reachable from the new
// profile. Explicit answers and fixed profile values drive the same closure.
func carryRecordedDecisions(root string, profile ResolvedProfile, explicit []DecisionValue, catalog *Catalog) ([]DecisionValue, []string, []Finding, error) {
	input, err := ResolveManifestInput(root, catalog)
	if errors.Is(err, ErrNoManifest) {
		return explicit, nil, nil, nil
	}
	if err != nil && input.Manifest.SchemaVersion == "" {
		return explicit, nil, []Finding{{Code: "baseline.decision.manifest-unreadable", Path: manifestPath, Message: "recorded decisions cannot be read: " + err.Error()}}, nil
	}
	values := make(map[string]any)
	candidates := make(map[string]any)
	reasons := make(map[string]string)
	for id, value := range profile.Values {
		values[id] = cloneJSONValue(value)
	}
	supplied := make(map[string]bool)
	for _, decision := range explicit {
		supplied[strings.TrimSpace(decision.ID)] = true
	}
	ids := make([]string, 0, len(input.Manifest.Decisions))
	for id := range input.Manifest.Decisions {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		recorded := input.Manifest.Decisions[id]
		if _, fixed := profile.Values[id]; fixed {
			reasons[id] = "profile " + profile.ID + " fixes it"
			continue
		}
		if supplied[id] {
			continue
		}
		declaration, known := catalog.decisions[id]
		if !known || validateDecisionValue(declaration, recorded.Value) != nil {
			reasons[id] = "the catalog refuses its recorded value"
			continue
		}
		candidates[id] = cloneJSONValue(recorded.Value)
		values[id] = cloneJSONValue(recorded.Value)
	}
	for _, decision := range explicit {
		values[strings.TrimSpace(decision.ID)] = decision.Value
	}
	selected := selectedPlanDecisions(profile, values, catalog)
	decisions := append([]DecisionValue(nil), explicit...)
	carried := []string{}
	warnings := []Finding{}
	for _, id := range ids {
		if _, fixed := profile.Values[id]; !fixed {
			if _, selects := selected[id]; !selects {
				reasons[id] = "profile " + profile.ID + " does not select it"
			}
		}
		if reason := reasons[id]; reason != "" {
			warnings = append(warnings, Finding{Code: "baseline.decision.not-carried", Path: manifestPath, Message: fmt.Sprintf("recorded decision %s is not carried: %s", id, reason)})
			continue
		}
		if value, ok := candidates[id]; ok {
			decisions = append(decisions, DecisionValue{ID: id, Value: value})
			carried = append(carried, id)
		}
	}
	return decisions, carried, warnings, nil
}

func conditionalMissingDecisionParts(profile ResolvedProfile, decisions []DecisionValue, missing []string, catalog *Catalog) []string {
	values := make(map[string]any, len(decisions))
	for _, decision := range decisions {
		values[decision.ID] = decision.Value
	}
	required := selectedPlanDecisions(profile, values, catalog)
	parts := []string{}
	for _, id := range missing {
		for _, effect := range objectsOrEmpty(catalog.decisions[id]["effects"]) {
			when, _ := objectValue(effect["when"])
			equals, ok := when["equals"]
			if !ok {
				continue
			}
			dependents := stringSet(stringsOrEmpty(effect["requireDecisions"]))
			for _, module := range stringsOrEmpty(effect["activateModules"]) {
				for _, dependent := range stringsOrEmpty(catalog.modules[module]["requiredDecisions"]) {
					dependents[dependent] = struct{}{}
				}
			}
			ids := []string{}
			for dependent := range dependents {
				if _, exists := required[dependent]; !exists && !decisionOptional(catalog.decisions[dependent]) {
					ids = append(ids, dependent)
				}
			}
			sort.Strings(ids)
			if len(ids) != 0 {
				parts = append(parts, fmt.Sprintf("required when %s is %v: %s", id, equals, strings.Join(ids, ", ")))
			}
		}
	}
	return parts
}
