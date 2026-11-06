package baseline

import (
	"fmt"
	"sort"
	"strings"
)

// ProfileDecisionGap identifies an unreachable decision required by a profile's template.
type ProfileDecisionGap struct {
	Decision  string
	Template  string
	Artifact  string
	Suggested string
}

// ProfileDecisionGaps reads the catalog and profile without resolving values or writing files.
// Templates with catalog-selected alternatives are excluded to avoid conditional false positives.
func ProfileDecisionGaps(profile ResolvedProfile, catalog *Catalog) []ProfileDecisionGap {
	reachable := stringSet(profile.Decisions)
	pending := append([]string(nil), profile.Decisions...)
	for len(pending) != 0 {
		id := pending[0]
		pending = pending[1:]
		for _, effect := range objectsOrEmpty(catalog.decisions[id]["effects"]) {
			required := append([]string(nil), stringsOrEmpty(effect["requireDecisions"])...)
			for _, moduleID := range stringsOrEmpty(effect["activateModules"]) {
				required = append(required, stringsOrEmpty(catalog.modules[moduleID]["requiredDecisions"])...)
			}
			for _, dependent := range required {
				if _, seen := reachable[dependent]; !seen {
					reachable[dependent] = struct{}{}
					pending = append(pending, dependent)
				}
			}
		}
	}

	type binding struct{ artifact, template string }
	templates := make(map[binding]bool)
	for _, moduleID := range profile.Modules {
		module := catalog.modules[moduleID]
		for _, field := range []string{"rootBlocks", "supportingGuides"} {
			for _, artifact := range objectsOrEmpty(module[field]) {
				id, _ := stringValue(artifact, "id")
				template, _ := stringValue(artifact, "template")
				templates[binding{id, template}] = true
			}
		}
	}
	for _, decision := range catalog.decisions {
		for _, effect := range objectsOrEmpty(decision["effects"]) {
			for _, selection := range objectsOrEmpty(effect["selectTemplates"]) {
				artifact, _ := stringValue(selection, "artifact")
				template, _ := stringValue(selection, "template")
				for candidate := range templates {
					if candidate.artifact == artifact && candidate.template != template {
						templates[candidate] = false
					}
				}
			}
		}
	}

	var gaps []ProfileDecisionGap
	seen := make(map[ProfileDecisionGap]bool)
	for id, decision := range catalog.decisions {
		if _, reached := reachable[id]; reached || decisionOptional(decision) {
			continue
		}
		suggested := ""
		value, ok := decision["suggestion"]
		if !ok {
			value, ok = decision["default"]
		}
		if ok {
			suggested = renderDecisionValue(value)
		}
		for _, effect := range objectsOrEmpty(decision["effects"]) {
			for _, item := range objectsOrEmpty(effect["renderBindings"]) {
				artifact, _ := stringValue(item, "artifact")
				template, _ := stringValue(item, "template")
				if !templates[binding{artifact, template}] {
					continue
				}
				gap := ProfileDecisionGap{Decision: id, Template: template, Artifact: artifact, Suggested: suggested}
				if !seen[gap] {
					seen[gap] = true
					gaps = append(gaps, gap)
				}
			}
		}
	}
	sort.Slice(gaps, func(i, j int) bool {
		if gaps[i].Decision != gaps[j].Decision {
			return gaps[i].Decision < gaps[j].Decision
		}
		if gaps[i].Template != gaps[j].Template {
			return gaps[i].Template < gaps[j].Template
		}
		return gaps[i].Artifact < gaps[j].Artifact
	})
	return gaps
}

// ProfileDecisionGapMessage names each missing decision and the template requiring it.
func ProfileDecisionGapMessage(profile ResolvedProfile, gaps []ProfileDecisionGap) string {
	parts := make([]string, 0, len(gaps))
	for _, gap := range gaps {
		parts = append(parts, fmt.Sprintf("profile %s does not list decision %s, which template %s requires", profile.ID, gap.Decision, gap.Template))
	}
	return strings.Join(parts, "; ")
}

// ProfileDecisionGapRepair names the profile entries and explicit values to supply.
func ProfileDecisionGapRepair(profile ResolvedProfile, gaps []ProfileDecisionGap, command string) string {
	suggestions := make(map[string]string)
	for _, gap := range gaps {
		suggestions[gap.Decision] = gap.Suggested
	}
	ids := sortedKeys(suggestions)
	entries := make([]string, 0, len(ids))
	flags := make([]string, 0, len(ids))
	for _, id := range ids {
		entries = append(entries, fmt.Sprintf("%q", id))
		value := suggestions[id]
		if value == "" {
			value = "<value>"
		}
		flags = append(flags, fmt.Sprintf("--decision '%s=%s'", id, value))
	}
	return fmt.Sprintf("add %s to \"decisions\" in %s, then %s %s", strings.Join(entries, ", "), profile.Path, command, strings.Join(flags, " "))
}

// A static gap may still render: normalized decisions or artifact defaults can
// supply its tokens. Keep this check aligned with the renderer's value sources.
func profileDecisionGapHasRenderValue(gap ProfileDecisionGap, decisions []DecisionValue, catalog *Catalog) bool {
	if decisionValueIndex(decisions, gap.Decision) >= 0 {
		return true
	}
	defaults := artifactRenderValues(catalog, document{"id": gap.Artifact}, nil, nil, nil, nil)
	for _, effect := range objectsOrEmpty(catalog.decisions[gap.Decision]["effects"]) {
		for _, binding := range objectsOrEmpty(effect["renderBindings"]) {
			artifact, _ := stringValue(binding, "artifact")
			template, _ := stringValue(binding, "template")
			if artifact != gap.Artifact || template != gap.Template {
				continue
			}
			token, _ := stringValue(binding, "token")
			if _, exists := defaults[token]; !exists {
				return false
			}
		}
	}
	return true
}
