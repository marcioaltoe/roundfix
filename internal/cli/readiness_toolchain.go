package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/judge"
)

var readinessToolPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._+-]*$`)
var readinessAssignmentPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=`)

// Only the first word is interpreted; this is deliberately not a shell parser.
func readinessCommandTool(command string) string {
	for _, word := range strings.Fields(command) {
		if readinessAssignmentPattern.MatchString(word) {
			continue
		}
		switch word {
		case "cd", "export", "set", "unset", "test", "[", "if", "for", "while", "until", "case", "exec", "eval", "source", ".", ":", "true", "false", "command", "builtin":
			return ""
		}
		if readinessToolPattern.MatchString(word) {
			return word
		}
		return ""
	}
	return ""
}

func toolchainReadiness(deps readinessDependencies, loaded roundconfig.Loaded) CheckResult {
	tools := make(map[string][]string)
	var findings []readinessFinding
	addCommand := func(source, command string) {
		if strings.TrimSpace(command) == "" {
			return
		}
		tool := readinessCommandTool(command)
		if tool == "" {
			findings = append(findings, readinessFinding{"DR-TOOL-UNREAD", CheckStatusWarn,
				fmt.Sprintf("could not read the tool for %s command %q", source, command),
				"list the tools this command needs under verification.tools in Project Config"})
			return
		}
		tools[tool] = append(tools[tool], source)
	}
	addCommand("defaults.verification", loaded.Config.Defaults.Verification)
	addCommand("worktree.bootstrap", loaded.Config.Worktree.Bootstrap)
	for index, declaration := range loaded.Config.Delivery.DerivedPaths {
		addCommand(fmt.Sprintf("delivery.derived_paths[%d].regenerate", index), declaration.Regenerate)
	}
	// A missing or unreadable manifest is diagnosed by Doctor's skills line.
	if loaded.GitRoot != "" {
		data, err := os.ReadFile(filepath.Join(loaded.GitRoot, doctorSetupManifestPath))
		if err == nil {
			// Other decisions have boolean values: decode only the two command decisions.
			var raw struct {
				Decisions map[string]json.RawMessage `json:"decisions"`
			}
			if json.Unmarshal(data, &raw) == nil {
				for _, id := range []string{"verification.gate", "verification.incremental"} {
					var decision struct {
						Value string `json:"value"`
					}
					if value, ok := raw.Decisions[id]; ok && json.Unmarshal(value, &decision) == nil {
						addCommand("Setup Manifest "+id, decision.Value)
					}
				}
			}
		}
	}
	for _, tool := range loaded.Config.Verification.Tools {
		tools[tool] = append(tools[tool], "verification.tools")
	}
	names := make([]string, 0, len(tools))
	for tool := range tools {
		names = append(names, tool)
	}
	sort.Strings(names)
	var found []string
	for _, tool := range names {
		if _, err := deps.resolve(tool); err != nil {
			findings = append(findings, readinessFinding{"DR-TOOL-MISSING", CheckStatusFailed,
				tool + " is not on PATH, needed by " + strings.Join(tools[tool], ", "),
				"install " + tool + ", or change the command that names it"})
		} else {
			found = append(found, tool)
		}
	}
	detail := "no tools declared"
	if len(found) > 0 {
		detail = strings.Join(found, ", ") + " found"
	} else if len(names) > 0 {
		detail = ""
	}
	return readinessResult(HealthCheckToolchain, detail, findings)
}

// The environment is searched only for the requested name. Key values are
// reduced to presence immediately and never retained in a map or result.
func readinessKeyIsSet(environ []string, name string) bool {
	set := false
	for _, entry := range environ {
		if strings.HasPrefix(entry, name+"=") {
			set = len(entry) > len(name)+1
		}
	}
	return set
}

func environmentReadiness(deps readinessDependencies) CheckResult {
	var options string
	for _, entry := range deps.environ {
		if strings.HasPrefix(entry, "NODE_OPTIONS=") {
			options = strings.TrimPrefix(entry, "NODE_OPTIONS=")
		}
	}
	var findings []readinessFinding
	missingPaths := make(map[string]bool)
	for _, path := range agent.MissingNodePreloads(options, deps.exists) {
		if missingPaths[path] {
			continue
		}
		missingPaths[path] = true
		findings = append(findings, readinessFinding{"DR-NODE-PRELOAD-MISSING", CheckStatusWarn,
			fmt.Sprintf("NODE_OPTIONS preload %q does not exist", path),
			"remove the preload from NODE_OPTIONS where your shell sets it"})
	}
	questions, err := judge.Load()
	if err != nil {
		return CheckResult{Name: HealthCheckEnvironment, Status: CheckStatusFailed, Detail: "load spec judge transports: " + err.Error(), Err: err}
	}
	var keys []string
	seen := make(map[string]bool)
	for _, transport := range questions.Transports {
		name := transport.KeyVariable
		if !strings.HasPrefix(name, "ROUNDFIX_") || seen[name] {
			continue
		}
		seen[name] = true
		status := "not set"
		if readinessKeyIsSet(deps.environ, name) {
			status = "set"
		}
		keys = append(keys, name+" "+status)
	}
	result := readinessResult(HealthCheckEnvironment, "", findings)
	if result.Detail != "" {
		result.Detail += "; "
	}
	result.Detail += "spec judge keys: " + strings.Join(keys, ", ")
	return result
}
