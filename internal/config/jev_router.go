package config

import "strings"

// JevRouterModel is the OpenCode model value that names the Jev Router.
const JevRouterModel = "roundfix-openrouter/typesafe/jev-router"

// IsJevRouterSelection reports whether runtime and model name the routed
// selection. It lives in config so config validation needs no agent import.
func IsJevRouterSelection(runtime, model string) bool {
	return strings.TrimSuffix(strings.TrimSpace(runtime), "-custom") == "opencode" && strings.TrimSpace(model) == JevRouterModel
}
