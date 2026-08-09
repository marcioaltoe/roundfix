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

// Values are the identifiers @agentclientprotocol/claude-agent-acp advertises,
// with the bracketed context suffix removed as the capability parser does.
// Verified against adapter 0.63.0 on 2026-08-07. The previous list offered
// claude-opus-5 and claude-opus-4-8, which no adapter advertises, and omitted
// sonnet, haiku, and default, which it does.
var claudeModelCatalog = []ModelChoice{
	{Label: "opus", Value: "opus", Description: "Opus 5 with a 1M context window; preferred design and frontend model"},
	{Label: "claude-fable-5", Value: "claude-fable-5", Description: "Fable 5; most capable for the hardest work, at roughly 25x the latency"},
	{Label: "sonnet", Value: "sonnet", Description: "Sonnet 5; efficient for routine tasks"},
	{Label: "haiku", Value: "haiku", Description: "Haiku 4.5; fastest for quick answers"},
	{Label: "default", Value: "default", Description: "adapter default; currently Opus 5 with a 1M context window"},
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
