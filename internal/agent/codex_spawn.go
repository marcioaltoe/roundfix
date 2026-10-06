package agent

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"roundfix/internal/codex"
)

const codexPathEnv = "CODEX_PATH"

type codexSpawnDependencies struct {
	goos       string
	getenv     func(string) string
	lookPath   codex.LookPathFunc
	quarantine codex.QuarantineProbe
	acceptance codex.AcceptanceProbe
}

type codexSpawnResolution struct {
	env []string
}

type codexSpawnError struct {
	result codex.Result
}

func (err codexSpawnError) Error() string {
	detail := strings.TrimSpace(err.result.Detail)
	if detail == "" {
		detail = "codex failed hygiene inspection"
	}
	if nextAction := strings.TrimSpace(err.result.NextAction); nextAction != "" {
		return fmt.Sprintf("codex runtime is not safe for acpx launch: %s; next: %s", detail, nextAction)
	}
	return fmt.Sprintf("codex runtime is not safe for acpx launch: %s", detail)
}

func (deps codexSpawnDependencies) resolve(ctx context.Context) (codexSpawnResolution, error) {
	goos := strings.TrimSpace(deps.goos)
	if goos == "" {
		goos = runtime.GOOS
	}
	if goos != "darwin" {
		return codexSpawnResolution{}, nil
	}

	var firstFailure *codex.Result
	if configuredPath := strings.TrimSpace(deps.configuredPath()); configuredPath != "" {
		result := deps.inspect(ctx, goos, configuredPath)
		if result.Status == codex.StatusOK {
			return newCodexSpawnResolution(result.Hygiene.Path, configuredPath), nil
		}
		firstFailure = &result
	}

	path, err := deps.pathCodex()
	if err != nil {
		if firstFailure != nil {
			return codexSpawnResolution{}, codexSpawnError{result: *firstFailure}
		}
		return codexSpawnResolution{}, fmt.Errorf("resolve clean codex for acpx launch: %w; next: %s", err, codex.ReinstallNextAction)
	}
	if firstFailure != nil && path == codexResultPath(*firstFailure) {
		return codexSpawnResolution{}, codexSpawnError{result: *firstFailure}
	}

	result := deps.inspect(ctx, goos, path)
	if result.Status == codex.StatusOK {
		return newCodexSpawnResolution(result.Hygiene.Path, path), nil
	}
	if firstFailure != nil {
		return codexSpawnResolution{}, codexSpawnError{result: *firstFailure}
	}
	return codexSpawnResolution{}, codexSpawnError{result: result}
}

func (deps codexSpawnDependencies) configuredPath() string {
	if deps.getenv != nil {
		return deps.getenv(codexPathEnv)
	}
	return os.Getenv(codexPathEnv)
}

func (deps codexSpawnDependencies) pathCodex() (string, error) {
	lookPath := deps.lookPath
	if lookPath == nil {
		lookPath = exec.LookPath
	}
	path, err := lookPath(codex.BinaryName)
	if err != nil {
		return "", fmt.Errorf("resolve codex on PATH: %w", err)
	}
	return path, nil
}

func (deps codexSpawnDependencies) inspect(ctx context.Context, goos string, path string) codex.Result {
	return codex.Inspector{
		ConfiguredPath: path,
		GOOS:           goos,
		LookPath:       deps.lookPath,
		Quarantine:     deps.quarantine,
		Acceptance:     deps.acceptance,
	}.Inspect(ctx)
}

func newCodexSpawnResolution(inspectedPath string, fallbackPath string) codexSpawnResolution {
	path := strings.TrimSpace(inspectedPath)
	if path == "" {
		path = strings.TrimSpace(fallbackPath)
	}
	return codexSpawnResolution{env: []string{codexPathEnv + "=" + path}}
}

func codexResultPath(result codex.Result) string {
	return strings.TrimSpace(result.Hygiene.Path)
}

// claudeNestedGuardEnv is Claude Code's nested-session guard variable. A
// Claude-driven orchestrator exports it, and an inherited copy makes the
// spawned claude runtime refuse to start ("cannot be launched inside another
// Claude Code session"). acpx-spawned Agent runtimes are independent
// processes, not nested sessions, so the guard never applies to them.
const claudeNestedGuardEnv = "CLAUDECODE"

// acpxCommandEnv builds the acpx child environment from the explicit base:
// variables that must not leak into Agent runtimes are removed (the codex
// hygiene path, re-added per session through overrides, and Claude Code's
// nested-session guard), then per-runtime overrides are appended.
func acpxCommandEnv(base []string, overrides []string) []string {
	// An override without an equals sign removes the named inherited variable.
	replaced := make(map[string]bool, len(overrides))
	for _, entry := range overrides {
		key, _, _ := strings.Cut(entry, "=")
		replaced[key] = true
	}
	filtered := make([]string, 0, len(base)+len(overrides))
	for _, entry := range base {
		key, _, _ := strings.Cut(entry, "=")
		if key == codexPathEnv || key == claudeNestedGuardEnv || replaced[key] {
			continue
		}
		filtered = append(filtered, entry)
	}
	for _, entry := range overrides {
		if strings.Contains(entry, "=") {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func (runner *ACPXRunner) commandEnv(overrides []string) []string {
	return acpxCommandEnv(runner.baseEnv(), overrides)
}

func (runner *ACPXRunner) baseEnv() []string {
	var environment []string
	if runner.Environment == nil {
		environment = os.Environ()
	} else {
		environment = append([]string(nil), runner.Environment...)
	}
	for index := len(environment) - 1; index >= 0; index-- {
		key, value, found := strings.Cut(environment[index], "=")
		if !found || key != "NODE_OPTIONS" {
			continue
		}
		kept, dropped, ok := agentNodeOptions(value, func(path string) bool {
			_, err := os.Stat(path)
			// Only confirmed absence authorizes removing a preload. Permission
			// and other inspection failures do not prove that the path is gone.
			return !os.IsNotExist(err)
		})
		if !ok {
			return environment
		}
		runner.noticeNodePreloads(dropped)
		if kept != "" {
			environment[index] = "NODE_OPTIONS=" + kept
			return environment
		}
		// Remove shadowed entries too, so an older value cannot become active
		// after removing the last entry.
		filtered := environment[:0]
		for _, entry := range environment {
			entryKey, _, _ := strings.Cut(entry, "=")
			if entryKey != "NODE_OPTIONS" {
				filtered = append(filtered, entry)
			}
		}
		return filtered
	}
	return environment
}

func environmentValue(environment []string, key string) string {
	for index := len(environment) - 1; index >= 0; index-- {
		entryKey, value, ok := strings.Cut(environment[index], "=")
		if ok && entryKey == key {
			return value
		}
	}
	return ""
}
