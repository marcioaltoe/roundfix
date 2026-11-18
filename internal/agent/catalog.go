package agent

import "strings"

type ModelChoice struct {
	Label       string
	Value       string
	Description string
}

var codexModelCatalog = []ModelChoice{
	{Label: "gpt-6.1-sol", Value: "gpt-6.1-sol", Description: "current Codex workhorse; built-in implementation default"},
	{Label: "gpt-6-astra", Value: "gpt-6-astra", Description: "frontier model for the most demanding work; drains quota fastest"},
	{Label: "gpt-6.1-terra", Value: "gpt-6.1-terra", Description: "fast, efficient writing model"},
	{Label: "gpt-6-sol", Value: "gpt-6-sol", Description: "previous workhorse"},
	{Label: "gpt-6-luna", Value: "gpt-6-luna", Description: "fast and affordable; built-in review default"},
	{Label: "gpt-5.6-sol", Value: "gpt-5.6-sol", Description: "older workhorse"},
	{Label: "gpt-5.6-terra", Value: "gpt-5.6-terra", Description: "older balanced model"},
	{Label: "gpt-5.6-luna", Value: "gpt-5.6-luna", Description: "older fast model"},
}

// Values are the identifiers @agentclientprotocol/claude-agent-acp advertises,
// with the bracketed context suffix removed as the capability parser does.
// Recorded from adapter 0.107.3 on 2026-11-17.
var claudeModelCatalog = []ModelChoice{
	{Label: "opus", Value: "opus", Description: "Opus 5.5; design and frontend default"},
	{Label: "sonnet", Value: "sonnet", Description: "Sonnet 5.5; efficient for routine tasks"},
	{Label: "fable", Value: "fable", Description: "alias for Fable 5.2; most capable for the hardest work, at the highest latency and quota cost"},
	{Label: "haiku", Value: "haiku", Description: "Haiku 5.5; fastest and cheapest"},
	{Label: "default", Value: "default", Description: "adapter default; currently Opus 5.5"},
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
