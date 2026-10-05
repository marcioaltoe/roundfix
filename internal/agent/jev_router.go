package agent

import (
	"encoding/json"
	"strings"

	"roundfix/internal/config"
)

const (
	JevRouterModel            = config.JevRouterModel
	JevRouterKeyEnv           = "ROUNDFIX_OPENROUTER_API_KEY"
	JevRouterKeyMissing       = "jev_router_key_missing"
	jevRouterProviderTemplate = `{"$schema":"https://opencode.ai/config.json","provider":{"roundfix-openrouter":{"npm":"@ai-sdk/openai-compatible","name":"Roundfix OpenRouter","options":{"baseURL":"https://openrouter.ai/api/v1","apiKey":"{env:ROUNDFIX_OPENROUTER_API_KEY}"},"models":{"typesafe/jev-router":{"name":"Jev Router"}}}}}`
)

// IsJevRouterSelection reports whether runtime and model name the routed selection.
func IsJevRouterSelection(runtime, model string) bool {
	return config.IsJevRouterSelection(runtime, model)
}

// jevRouterProviderConfig keeps credentials in OpenCode's environment boundary.
func jevRouterProviderConfig(baseURL string) string {
	encoded, _ := json.Marshal(baseURL)
	return strings.Replace(jevRouterProviderTemplate, `"https://openrouter.ai/api/v1"`, string(encoded), 1)
}
