package config

import (
	"fmt"
	"strings"
)

var recommendedAuthoringDeclaration = AuthoringHarness{
	Idea:     &AuthoringEntry{Model: "fable", ReasoningEffort: "xhigh"},
	PRD:      &AuthoringEntry{Model: "fable", ReasoningEffort: "high"},
	TechSpec: &AuthoringEntry{Model: "fable", ReasoningEffort: "high"},
	Tasks:    &AuthoringEntry{Model: "opus", ReasoningEffort: "high"},
	Research: &ResearchEntry{AuthoringEntry: AuthoringEntry{Model: "sonnet", ReasoningEffort: "medium"}},
}

// RecommendedAuthoringDeclaration returns a fresh copy of the Recommended
// Authoring Declaration for the claude harness. It declares no research tools.
func RecommendedAuthoringDeclaration() AuthoringHarness {
	idea, prd := *recommendedAuthoringDeclaration.Idea, *recommendedAuthoringDeclaration.PRD
	techspec, tasks := *recommendedAuthoringDeclaration.TechSpec, *recommendedAuthoringDeclaration.Tasks
	research := *recommendedAuthoringDeclaration.Research
	return AuthoringHarness{Idea: &idea, PRD: &prd, TechSpec: &techspec, Tasks: &tasks, Research: &research}
}

// RecommendedAuthoringYAML renders the declaration for .roundfixrc.yml,
// starting at skills: and ending with a newline, without a tools key.
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
	}
	return block.String()
}
