package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"roundfix/internal/agent"
	"roundfix/internal/codex"
)

type CheckStatus string

const (
	CheckStatusOK      CheckStatus = "ok"
	CheckStatusSkipped CheckStatus = "skipped"
	CheckStatusFailed  CheckStatus = "failed"

	HealthCheckNode     = "node"
	HealthCheckACPX     = "acpx"
	HealthCheckAdapter  = "adapter"
	HealthCheckAgent    = "agent"
	HealthCheckProfiles = "profiles"
	HealthCheckCodex    = "codex"
)

type CheckResult struct {
	Name       string
	Status     CheckStatus
	Detail     string
	NextAction string
	Err        error
}

type HealthChecker interface {
	Node(ctx context.Context) CheckResult
	ACPX(ctx context.Context) CheckResult
	Adapter(ctx context.Context, runtime agent.RuntimeSpec) CheckResult
	Agent(ctx context.Context, req agent.ProbeRequest) CheckResult
	Codex(ctx context.Context) CheckResult
}

type healthCheckDependencies struct {
	nodeVersion    func(context.Context) (string, error)
	acpxVersion    func(context.Context) (string, error)
	checkAdapter   func(context.Context, agent.RuntimeSpec) (agent.AdapterEvidence, error)
	probeAgent     func(context.Context, agent.ProbeRequest) error
	codexInspector codexInspector
}

type runtimeHealthChecker struct {
	deps healthCheckDependencies
}

type codexInspector interface {
	Inspect(ctx context.Context) codex.Result
}

func newHealthChecker(deps healthCheckDependencies) HealthChecker {
	return runtimeHealthChecker{deps: deps}
}

func (checker runtimeHealthChecker) Node(ctx context.Context) CheckResult {
	version, err := checker.deps.nodeVersion(ctx)
	if err != nil {
		return CheckResult{
			Name:       HealthCheckNode,
			Status:     CheckStatusFailed,
			Detail:     fmt.Sprintf("Node.js %s or newer is required: %v", setupNodeMinimumVersion, err),
			NextAction: fmt.Sprintf("install Node.js %s or newer", setupNodeMinimumVersion),
		}
	}
	if compareVersions(version, setupNodeMinimumVersion) < 0 {
		return CheckResult{
			Name:       HealthCheckNode,
			Status:     CheckStatusFailed,
			Detail:     fmt.Sprintf("found %s; Node.js %s or newer is required", strings.TrimSpace(version), setupNodeMinimumVersion),
			NextAction: fmt.Sprintf("install Node.js %s or newer", setupNodeMinimumVersion),
		}
	}
	return CheckResult{
		Name:   HealthCheckNode,
		Status: CheckStatusOK,
		Detail: fmt.Sprintf("%s >= %s", strings.TrimSpace(version), setupNodeMinimumVersion),
	}
}

func (checker runtimeHealthChecker) ACPX(ctx context.Context) CheckResult {
	version, err := checker.deps.acpxVersion(ctx)
	if err == nil && strings.TrimSpace(version) == agent.PinnedACPXVersion {
		return CheckResult{
			Name:   HealthCheckACPX,
			Status: CheckStatusOK,
			Detail: agent.PinnedACPXVersion,
		}
	}

	detail := fmt.Sprintf("run %s", setupACPXInstallCommand())
	if err == nil {
		detail = fmt.Sprintf("found %s; required %s; run %s", strings.TrimSpace(version), agent.PinnedACPXVersion, setupACPXInstallCommand())
	}
	return CheckResult{
		Name:       HealthCheckACPX,
		Status:     CheckStatusFailed,
		Detail:     detail,
		NextAction: setupACPXInstallCommand(),
	}
}

func (checker runtimeHealthChecker) Adapter(ctx context.Context, runtime agent.RuntimeSpec) CheckResult {
	checkAdapter := checker.deps.checkAdapter
	if checkAdapter == nil {
		checkAdapter = agent.CheckAdapter
	}
	evidence, err := checkAdapter(ctx, runtime)
	if err == nil {
		return CheckResult{
			Name:   HealthCheckAdapter,
			Status: CheckStatusOK,
			Detail: adapterEvidenceDetail(evidence),
		}
	}
	result := CheckResult{
		Name:   HealthCheckAdapter,
		Status: CheckStatusFailed,
		Detail: err.Error(),
	}
	var installErr interface {
		error
		InstallCommand() string
	}
	if errors.As(err, &installErr) {
		result.NextAction = installErr.InstallCommand()
	}
	return result
}

func adapterEvidenceDetail(evidence agent.AdapterEvidence) string {
	command := strings.TrimSpace(evidence.Command)
	if evidence.Package == "" && evidence.Version == "" {
		return command
	}
	return fmt.Sprintf("command=%q; package=%s; version=%s", command, strings.TrimSpace(evidence.Package), strings.TrimSpace(evidence.Version))
}

func (checker runtimeHealthChecker) Agent(ctx context.Context, req agent.ProbeRequest) CheckResult {
	if err := checker.deps.probeAgent(ctx, req); err != nil {
		return CheckResult{
			Name:   HealthCheckAgent,
			Status: CheckStatusFailed,
			Detail: err.Error(),
			Err:    err,
		}
	}
	return CheckResult{
		Name:   HealthCheckAgent,
		Status: CheckStatusOK,
		Detail: req.Runtime.ID,
	}
}

func (checker runtimeHealthChecker) Codex(ctx context.Context) CheckResult {
	inspector := checker.deps.codexInspector
	if inspector == nil {
		// Doctor must inspect the codex a Run would actually spawn: CODEX_PATH
		// first (mirroring the codex-acp spawn path), then PATH inside the
		// inspector. Without the configured path, Doctor would check the wrong
		// binary.
		inspector = codex.Inspector{ConfiguredPath: os.Getenv("CODEX_PATH")}
	}
	result := inspector.Inspect(ctx)
	return CheckResult{
		Name:       HealthCheckCodex,
		Status:     CheckStatus(result.Status),
		Detail:     result.Detail,
		NextAction: result.NextAction,
	}
}
