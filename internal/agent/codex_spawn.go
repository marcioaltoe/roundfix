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
	filtered := make([]string, 0, len(base)+len(overrides))
	for _, entry := range base {
		key, _, _ := strings.Cut(entry, "=")
		if key == codexPathEnv || key == claudeNestedGuardEnv {
			continue
		}
		filtered = append(filtered, entry)
	}
	return append(filtered, overrides...)
}

func (runner ACPXRunner) commandEnv(overrides []string) []string {
	return acpxCommandEnv(runner.baseEnv(), overrides)
}

func (runner ACPXRunner) baseEnv() []string {
	if runner.Environment == nil {
		return os.Environ()
	}
	return append([]string(nil), runner.Environment...)
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
