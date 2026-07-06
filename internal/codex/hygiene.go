package codex

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

const (
	BinaryName          = "codex"
	QuarantineAttribute = "com.apple.quarantine"
	ReinstallNextAction = "reinstall codex with the official curl installer into ~/.local/bin, then set CODEX_PATH to that binary"
)

type Status string

const (
	StatusOK      Status = "ok"
	StatusSkipped Status = "skipped"
	StatusFailed  Status = "failed"
)

type Result struct {
	Status     Status
	Detail     string
	NextAction string
	Hygiene    CodexHygiene
}

type CodexHygiene struct {
	Path        string
	Quarantined bool
	Accepted    bool
}

type QuarantineProbe interface {
	Quarantined(ctx context.Context, path string) (bool, error)
}

type AcceptanceProbe interface {
	Accepted(ctx context.Context, path string) (bool, error)
}

type LookPathFunc func(string) (string, error)

type Inspector struct {
	ConfiguredPath string
	GOOS           string
	LookPath       LookPathFunc
	Quarantine     QuarantineProbe
	Acceptance     AcceptanceProbe
}

func (inspector Inspector) Inspect(ctx context.Context) Result {
	goos := inspector.GOOS
	if strings.TrimSpace(goos) == "" {
		goos = runtime.GOOS
	}
	if goos != "darwin" {
		return Result{
			Status: StatusSkipped,
			Detail: "not-applicable on " + goos,
		}
	}

	path, err := inspector.resolve()
	if err != nil {
		return Result{
			Status:     StatusFailed,
			Detail:     err.Error(),
			NextAction: ReinstallNextAction,
		}
	}

	quarantined, err := inspector.quarantineProbe().Quarantined(ctx, path)
	if err != nil {
		return Result{
			Status: StatusFailed,
			Detail: fmt.Sprintf("inspect codex quarantine at %s: %v", path, err),
			Hygiene: CodexHygiene{
				Path: path,
			},
		}
	}
	accepted, err := inspector.acceptanceProbe().Accepted(ctx, path)
	if err != nil {
		return Result{
			Status: StatusFailed,
			Detail: fmt.Sprintf("verify codex signature at %s: %v", path, err),
			Hygiene: CodexHygiene{
				Path:        path,
				Quarantined: quarantined,
			},
		}
	}

	hygiene := CodexHygiene{
		Path:        path,
		Quarantined: quarantined,
		Accepted:    accepted,
	}
	if quarantined {
		return Result{
			Status:     StatusFailed,
			Detail:     fmt.Sprintf("%s is quarantined with %s", path, QuarantineAttribute),
			NextAction: ReinstallNextAction,
			Hygiene:    hygiene,
		}
	}
	if !accepted {
		return Result{
			Status:     StatusFailed,
			Detail:     fmt.Sprintf("%s has an invalid or missing code signature", path),
			NextAction: ReinstallNextAction,
			Hygiene:    hygiene,
		}
	}
	return Result{
		Status:  StatusOK,
		Detail:  fmt.Sprintf("%s has a valid code signature and no %s attribute", path, QuarantineAttribute),
		Hygiene: hygiene,
	}
}

func (inspector Inspector) resolve() (string, error) {
	configuredPath := strings.TrimSpace(inspector.ConfiguredPath)
	if configuredPath != "" {
		return configuredPath, nil
	}
	lookPath := inspector.LookPath
	if lookPath == nil {
		lookPath = exec.LookPath
	}
	path, err := lookPath(BinaryName)
	if err != nil {
		return "", fmt.Errorf("resolve codex on PATH: %w", err)
	}
	return path, nil
}

func (inspector Inspector) quarantineProbe() QuarantineProbe {
	if inspector.Quarantine != nil {
		return inspector.Quarantine
	}
	return xattrQuarantineProbe{}
}

func (inspector Inspector) acceptanceProbe() AcceptanceProbe {
	if inspector.Acceptance != nil {
		return inspector.Acceptance
	}
	return codesignAcceptanceProbe{}
}

type xattrQuarantineProbe struct{}

func (xattrQuarantineProbe) Quarantined(ctx context.Context, path string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	cmd := exec.CommandContext(ctx, "xattr", path)
	output, err := cmd.CombinedOutput()
	if ctxErr := ctx.Err(); ctxErr != nil {
		return false, ctxErr
	}
	if err != nil {
		return false, commandOutputError("run xattr", err, output)
	}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) == QuarantineAttribute {
			return true, nil
		}
	}
	return false, nil
}

// codesignAcceptanceProbe treats a codex with a valid code signature as
// acceptable. It deliberately does NOT use `spctl --assess`, which rejects any
// signed CLI binary that is not a notarized app bundle ("the code is valid but
// does not seem to be an app") — codex is OpenAI-signed and never Apple-
// notarized, so spctl would reject every codex. The real XProtect trigger is
// the quarantine attribute, checked separately.
type codesignAcceptanceProbe struct{}

func (codesignAcceptanceProbe) Accepted(ctx context.Context, path string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	cmd := exec.CommandContext(ctx, "codesign", "--verify", path)
	output, err := cmd.CombinedOutput()
	if ctxErr := ctx.Err(); ctxErr != nil {
		return false, ctxErr
	}
	if err == nil {
		return true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return false, nil
	}
	return false, commandOutputError("run codesign verify", err, output)
}

func commandOutputError(operation string, err error, output []byte) error {
	detail := strings.TrimSpace(string(output))
	if detail == "" {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return fmt.Errorf("%s: %w: %s", operation, err, detail)
}
