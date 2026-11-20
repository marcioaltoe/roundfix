package flowpermissions

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ReadOpenCode reports existing shell restrictions. Its second result is the
// restricting filename for own_rules, or the diagnostic for unreadable.
func ReadOpenCode(repoRoot string) (state State, file string) {
	state = NotNeeded
	for _, name := range []string{"opencode.json", "opencode.jsonc"} {
		if err := checkTarget(repoRoot, name); err != nil {
			return Unreadable, name + ": " + err.Error()
		}
		data, err := os.ReadFile(filepath.Join(repoRoot, name))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return Unreadable, name + ": " + err.Error()
		}
		clean, err := cleanJSONC(data)
		if err != nil {
			return Unreadable, name + ": " + err.Error()
		}
		var config map[string]json.RawMessage
		if err := json.Unmarshal(clean, &config); err != nil {
			return Unreadable, name + ": " + err.Error()
		}
		if config == nil {
			return Unreadable, name + ": expected an object"
		}
		restricted, err := openCodeShellRules(config)
		if err != nil {
			return Unreadable, name + ": " + err.Error()
		}
		if restricted && state != OwnRules {
			state, file = OwnRules, name
		}
	}
	return state, file
}

func openCodeShellRules(config map[string]json.RawMessage) (bool, error) {
	restricted := false
	if raw, exists := config["permissions"]; exists {
		var rules []struct {
			Action string `json:"action"`
			Effect string `json:"effect"`
		}
		if len(raw) == 0 || raw[0] != '[' {
			return false, fmt.Errorf("permissions: expected a list")
		}
		if err := json.Unmarshal(raw, &rules); err != nil {
			return false, fmt.Errorf("permissions: %w", err)
		}
		for _, rule := range rules {
			if (rule.Action == "shell" || rule.Action == "*") && (rule.Effect == "ask" || rule.Effect == "deny") {
				restricted = true
			}
		}
	}
	if raw, exists := config["permission"]; exists {
		var value any
		if err := json.Unmarshal(raw, &value); err != nil {
			return false, fmt.Errorf("permission: %w", err)
		}
		switch permission := value.(type) {
		case string:
			restricted = restricted || asksOrDenies(permission)
		case map[string]any:
			if bash, ok := permission["bash"]; ok {
				asks, err := shellPermission(bash)
				if err != nil {
					return false, err
				}
				restricted = restricted || asks
			}
			// A wildcard permission can restrict shell commands as well.
			if wildcard, ok := permission["*"]; ok {
				asks, err := shellPermission(wildcard)
				if err != nil {
					return false, err
				}
				restricted = restricted || asks
			}
		default:
			return false, fmt.Errorf("permission: expected a string or object")
		}
	}
	return restricted, nil
}

func asksOrDenies(value string) bool { return value == "ask" || value == "deny" }

func shellPermission(value any) (bool, error) {
	switch permission := value.(type) {
	case string:
		return asksOrDenies(permission), nil
	case map[string]any:
		restricted := false
		for _, value := range permission {
			text, ok := value.(string)
			if !ok {
				return false, fmt.Errorf("shell permission: expected string effects")
			}
			restricted = restricted || asksOrDenies(text)
		}
		return restricted, nil
	default:
		return false, fmt.Errorf("shell permission: expected a string or object")
	}
}

// Replace comments with whitespace first, so tokens separated by a comment
// cannot accidentally become a single token. Strings are copied verbatim.
func cleanJSONC(data []byte) ([]byte, error) {
	out := append([]byte(nil), data...)
	inString, escaped := false, false
	for i := 0; i < len(out); i++ {
		ch := out[i]
		if inString {
			if escaped {
				escaped = false
			} else if ch == '\\' {
				escaped = true
			} else if ch == '"' {
				inString = false
			}
			continue
		}
		if ch == '"' {
			inString = true
			continue
		}
		if ch != '/' || i+1 >= len(out) {
			continue
		}
		switch out[i+1] {
		case '/':
			out[i], out[i+1] = ' ', ' '
			i += 2
			for i < len(out) && out[i] != '\n' && out[i] != '\r' {
				out[i] = ' '
				i++
			}
			i--
		case '*':
			out[i], out[i+1] = ' ', ' '
			i += 2
			for i+1 < len(out) && !(out[i] == '*' && out[i+1] == '/') {
				if out[i] != '\n' && out[i] != '\r' {
					out[i] = ' '
				}
				i++
			}
			if i+1 >= len(out) {
				return nil, fmt.Errorf("unterminated block comment")
			}
			out[i], out[i+1] = ' ', ' '
			i++
		}
	}
	inString, escaped = false, false
	for i, ch := range out {
		if inString {
			if escaped {
				escaped = false
			} else if ch == '\\' {
				escaped = true
			} else if ch == '"' {
				inString = false
			}
			continue
		}
		if ch == '"' {
			inString = true
			continue
		}
		if ch != ',' {
			continue
		}
		next := i + 1
		for next < len(out) && (out[next] == ' ' || out[next] == '\t' || out[next] == '\r' || out[next] == '\n') {
			next++
		}
		if next < len(out) && (out[next] == ']' || out[next] == '}') {
			out[i] = ' '
		}
	}
	return out, nil
}
