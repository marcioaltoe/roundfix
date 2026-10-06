// Package openrouterkey names each stage's OpenRouter variables in preference order.
package openrouterkey

import "strings"

const (
	Shared    = "ROUNDFIX_OPENROUTER_API_KEY"
	Judge     = "ROUNDFIX_OPENROUTER_JUDGE_API_KEY"
	Implement = "ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY"
)

type Stage string

const (
	StageJudge     Stage = "judge"
	StageImplement Stage = "implement"
)

// Variables returns a fresh list of the stage key followed by the shared key.
func Variables(stage Stage) []string {
	switch stage {
	case StageJudge:
		return []string{Judge, Shared}
	case StageImplement:
		return []string{Implement, Shared}
	default:
		return nil
	}
}

// Select returns only the first set variable's name; the last duplicate wins.
func Select(environ []string, stage Stage) (string, bool) {
	variables := Variables(stage)
	for _, variable := range variables {
		set := false
		for _, entry := range environ {
			name, value, found := strings.Cut(entry, "=")
			if found && name == variable {
				set = value != ""
			}
		}
		if set {
			return variable, true
		}
	}
	return "", false
}
