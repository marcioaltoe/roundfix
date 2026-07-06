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

func acpxCommandEnv(overrides []string) []string {
	if len(overrides) == 0 {
		return nil
	}
	env := os.Environ()
	filtered := make([]string, 0, len(env)+len(overrides))
	for _, entry := range env {
		key, _, _ := strings.Cut(entry, "=")
		if key == codexPathEnv {
			continue
		}
		filtered = append(filtered, entry)
	}
	return append(filtered, overrides...)
}
