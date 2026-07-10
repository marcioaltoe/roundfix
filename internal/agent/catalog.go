package agent

import "strings"

type ModelChoice struct {
	Label       string
	Value       string
	Description string
}

var codexModelCatalog = []ModelChoice{
	{Label: "gpt-5.6-sol", Value: "gpt-5.6-sol", Description: "latest frontier agentic coding model"},
	{Label: "gpt-5.6-terra", Value: "gpt-5.6-terra", Description: "balanced everyday agentic coding"},
	{Label: "gpt-5.6-luna", Value: "gpt-5.6-luna", Description: "fast and affordable agentic coding"},
	{Label: "gpt-5.5", Value: "gpt-5.5", Description: "initial Default Agent Model"},
	{Label: "gpt-5.4", Value: "gpt-5.4", Description: "everyday coding"},
	{Label: "gpt-5.4-mini", Value: "gpt-5.4-mini", Description: "small and cost-efficient"},
	{Label: "gpt-5.3-codex-spark", Value: "gpt-5.3-codex-spark", Description: "ultra-fast coding"},
}

var claudeModelCatalog = []ModelChoice{
	{Label: "Default", Value: "default", Description: "effective configured Claude model"},
	{Label: "Opus", Value: "opus", Description: "Opus 4.8 with a 1M context window"},
	{Label: "Fable", Value: "fable", Description: "Fable 5"},
	{Label: "Sonnet", Value: "sonnet", Description: "Sonnet 5"},
	{Label: "Haiku", Value: "haiku", Description: "Haiku 4.5"},
}

func ModelCatalog(runtime string) []ModelChoice {
	switch strings.TrimSpace(runtime) {
	case "codex":
		return append([]ModelChoice(nil), codexModelCatalog...)
	case "claude":
		return append([]ModelChoice(nil), claudeModelCatalog...)
	default:
		return nil
	}
}
