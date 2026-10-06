package agent

import (
	"encoding/json"
	"errors"
)

// lightSessionEnvironment carries only a variable reference in inline config.
// Do not include inherited configuration in diagnostics: it can hold credentials.
func lightSessionEnvironment(runtime RuntimeSpec, environment []string) ([]string, error) {
	config := map[string]any{}
	if inherited := environmentValue(environment, "OPENCODE_CONFIG_CONTENT"); inherited != "" {
		if json.Unmarshal([]byte(inherited), &config) != nil || config == nil {
			return nil, selectionPreflightError(runtime, "start runtime", errors.New("OPENCODE_CONFIG_CONTENT must be a JSON object for a light session"))
		}
	}
	object := config
	for _, key := range []string{"provider", "openrouter", "options"} {
		child, ok := object[key].(map[string]any)
		if !ok {
			child = map[string]any{}
			object[key] = child
		}
		object = child
	}
	object["apiKey"] = "{env:" + runtime.OpenRouterKeyVariable + "}"
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	return []string{"OPENROUTER_API_KEY", "OPENCODE_CONFIG_CONTENT=" + string(data)}, nil
}
