package config

import (
	"fmt"
	"strings"
)

var recommendedAuthoringDeclaration = AuthoringHarness{
	Idea:     &AuthoringEntry{Model: "fable", ReasoningEffort: "xhigh", Fallback: &AuthoringFallback{Model: "opus", ReasoningEffort: "xhigh"}},
	PRD:      &AuthoringEntry{Model: "fable", ReasoningEffort: "high", Fallback: &AuthoringFallback{Model: "opus", ReasoningEffort: "xhigh"}},
	TechSpec: &AuthoringEntry{Model: "fable", ReasoningEffort: "high", Fallback: &AuthoringFallback{Model: "opus", ReasoningEffort: "xhigh"}},
	Tasks:    &AuthoringEntry{Model: "opus", ReasoningEffort: "high", Fallback: &AuthoringFallback{Model: "sonnet", ReasoningEffort: "high"}},
	Research: &ResearchEntry{AuthoringEntry: AuthoringEntry{Model: "sonnet", ReasoningEffort: "medium"}, Tools: []string{"WebSearch", "WebFetch"}},
}

// RecommendedAuthoringDeclaration returns a fresh copy of the Recommended
// Authoring Declaration for the claude harness, including the research tools
// WebSearch and WebFetch.
func RecommendedAuthoringDeclaration() AuthoringHarness {
	idea, prd := *recommendedAuthoringDeclaration.Idea, *recommendedAuthoringDeclaration.PRD
	techspec, tasks := *recommendedAuthoringDeclaration.TechSpec, *recommendedAuthoringDeclaration.Tasks
	research := *recommendedAuthoringDeclaration.Research
	research.Tools = append([]string{}, research.Tools...)
	for _, entry := range []*AuthoringEntry{&idea, &prd, &techspec, &tasks} {
		fallback := *entry.Fallback
		entry.Fallback = &fallback
	}
	return AuthoringHarness{Idea: &idea, PRD: &prd, TechSpec: &techspec, Tasks: &tasks, Research: &research}
}

// RecommendedAuthoringYAML renders the declaration for .roundfixrc.yml,
// starting at skills: and ending with a newline, with WebSearch and WebFetch
// under research.tools.
func RecommendedAuthoringYAML() string {
	declaration := RecommendedAuthoringDeclaration()
	var block strings.Builder
	block.WriteString("skills:\n  authoring:\n    claude:\n")
	for _, phase := range []struct {
		key   string
		entry *AuthoringEntry
	}{
		{"idea", declaration.Idea},
		{"prd", declaration.PRD},
		{"techspec", declaration.TechSpec},
		{"tasks", declaration.Tasks},
		{"research", &declaration.Research.AuthoringEntry},
	} {
		fmt.Fprintf(&block, "      %s:\n        model: %s\n        reasoning_effort: %s\n", phase.key, phase.entry.Model, phase.entry.ReasoningEffort)
		if phase.key == "research" {
			fmt.Fprintf(&block, "        tools: [%s]\n", strings.Join(declaration.Research.Tools, ", "))
		}
		if fallback := phase.entry.Fallback; fallback != nil {
			fmt.Fprintf(&block, "        fallback:\n          model: %s\n          reasoning_effort: %s\n", fallback.Model, fallback.ReasoningEffort)
		}
	}
	return block.String()
}
