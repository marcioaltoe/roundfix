package agent

import "strings"

const (
	JevRouterModel          = "roundfix-openrouter/typesafe/jev-router"
	JevRouterKeyEnv         = "ROUNDFIX_OPENROUTER_API_KEY"
	JevRouterKeyMissing     = "jev_router_key_missing"
	jevRouterProviderConfig = `{"$schema":"https://opencode.ai/config.json","provider":{"roundfix-openrouter":{"npm":"@ai-sdk/openai-compatible","name":"Roundfix OpenRouter","options":{"baseURL":"https://openrouter.ai/api/v1","apiKey":"{env:ROUNDFIX_OPENROUTER_API_KEY}"},"models":{"typesafe/jev-router":{"name":"Jev Router"}}}}}`
)

// IsJevRouterSelection reports whether runtime and model name the routed selection.
func IsJevRouterSelection(runtime, model string) bool {
	return strings.TrimSuffix(strings.TrimSpace(runtime), "-custom") == "opencode" && strings.TrimSpace(model) == JevRouterModel
}
