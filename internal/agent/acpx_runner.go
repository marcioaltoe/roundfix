package agent

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"roundfix/internal/runevent"
)

const (
	defaultACPXCommand              = "acpx"
	MinimumACPXVersion              = "0.12.0"
	CodexAdapterPackage             = "@agentclientprotocol/codex-acp"
	PinnedCodexAdapterVersion       = "1.1.5"
	legacyCodexAdapterPackage       = "@zed-industries/codex-acp"
	ClaudeAdapterPackage            = "@agentclientprotocol/claude-agent-acp"
	PinnedClaudeAdapterVersion      = "0.63.0"
	legacyClaudeCodeAdapterPackage  = "@zed-industries/" + "claude-code-acp"
	legacyClaudeAgentAdapterPackage = "@zed-industries/" + "claude-agent-acp"
	defaultCodexAdapterCommand      = "npx -y " + CodexAdapterPackage
	defaultClaudeAdapterCommand     = "npx -y " + ClaudeAdapterPackage + "@" + PinnedClaudeAdapterVersion
	adapterProbeOutputLimit         = 512
	acpxPermissionDeniedStatus      = "permissions_denied"
	acpxExitReasonAgentProtocol     = "agent/protocol error"
	acpxExitReasonTimeout           = "timeout"
	acpxExitReasonPermissionsDenied = "all permissions denied"
	acpxExitReasonUsage             = "usage error"
	acpxExitReasonMissingSession    = "missing session"
	acpxCodexSandboxUnavailable     = "codex_sandbox_full_access_unavailable"
	acpxCodexSandboxModeKey         = "sandbox_mode"
	acpxCodexFullAccessSandbox      = "danger-full-access"
	acpxCodexReasoningEffortKey     = "reasoning_effort"
	acpxGenericReasoningEffortKey   = "effort"
	acpxPreflightSessionPrefix      = "roundfix-preflight-"
	acpxPreflightSetupTimeout       = 30 * time.Second
	acpxPreflightCleanupTimeout     = 5 * time.Second
	infrastructureStderrTailLines   = 10
	infrastructureStderrTailBytes   = 1024
	infrastructureStderrDelimiter   = "\n--- acpx stderr tail ---\n"
	infrastructureStderrTruncated   = "[stderr truncated]\n"
)

var defaultAdapterCommands = map[string]string{
	"codex":    defaultCodexAdapterCommand,
	"claude":   defaultClaudeAdapterCommand,
	"opencode": "opencode",
}

var adapterInstallCommands = map[string]string{
	"claude-code-acp":  ClaudeAdapterInstallCommand(),
	"claude-agent-acp": ClaudeAdapterInstallCommand(),
	"opencode":         "npm install -g opencode-ai",
}

type adapterLineageContract struct {
	RuntimeID      string
	Package        string
	PinnedVersion  string
	LegacyPackages []string
	VersionOnly    bool
}

var adapterLineageContracts = map[string]adapterLineageContract{
	"codex": {
		RuntimeID:      "codex",
		Package:        CodexAdapterPackage,
		PinnedVersion:  PinnedCodexAdapterVersion,
		LegacyPackages: []string{legacyCodexAdapterPackage},
	},
	"claude": {
		RuntimeID:      "claude",
		Package:        ClaudeAdapterPackage,
		PinnedVersion:  PinnedClaudeAdapterVersion,
		LegacyPackages: []string{legacyClaudeCodeAdapterPackage, legacyClaudeAgentAdapterPackage},
		VersionOnly:    true,
	},
}

const (
	AdapterLineageUnknown     = "adapter_lineage_unknown"
	AdapterVersionUnsupported = "adapter_version_unsupported"
)

// AdapterEvidence is the bounded identity evidence used by readiness surfaces.
type AdapterEvidence struct {
	Command string
	Package string
	Version string
}

// CodexAdapterInstallCommand returns the deterministic official adapter action.
func CodexAdapterInstallCommand() string {
	return "npm install -g " + CodexAdapterPackage + "@" + PinnedCodexAdapterVersion
}

// CodexAdapterCommand returns the deterministic official adapter command that
// Setup persists when migrating an ACPX override.
func CodexAdapterCommand() string {
	return "npx -y " + CodexAdapterPackage + "@" + PinnedCodexAdapterVersion
}

// ClaudeAdapterInstallCommand returns the deterministic official adapter action.
func ClaudeAdapterInstallCommand() string {
	return "npm install -g " + ClaudeAdapterPackage + "@" + PinnedClaudeAdapterVersion
}

// ClaudeAdapterCommand returns the deterministic official adapter command that
// Setup persists when migrating an ACPX override.
func ClaudeAdapterCommand() string {
	return defaultClaudeAdapterCommand
}

// ACPXRunner is the acpx-backed invocation core. Later migration tasks wire
// this into Runner after Agent Session lifecycle is available.
type ACPXRunner struct {
	Command           string
	Environment       []string
	Now               func() time.Time
	warnf             func(string, ...any)
	cancelClock       cancellationClock
	stateMu           sync.Mutex
	ensuredSessions   map[string]struct{}
	sessionSelections map[string]SelectionAssignment
	codexSpawn        codexSpawnDependencies
	codexResolutions  map[string]codexSpawnResolution
}

type cancellationTimer interface {
	C() <-chan time.Time
	Stop() bool
}

type cancellationClock interface {
	NewTimer(time.Duration) cancellationTimer
}

type realCancellationClock struct{}

func (realCancellationClock) NewTimer(duration time.Duration) cancellationTimer {
	return realCancellationTimer{timer: time.NewTimer(duration)}
}

type realCancellationTimer struct {
	timer *time.Timer
}

func (timer realCancellationTimer) C() <-chan time.Time {
	return timer.timer.C
}

func (timer realCancellationTimer) Stop() bool {
	return timer.timer.Stop()
}

func (runner ACPXRunner) cancellationClock() cancellationClock {
	if runner.cancelClock != nil {
		return runner.cancelClock
	}
	return realCancellationClock{}
}

// ACPXPromptRequest carries the explicit Agent Session needed by acpx prompt
// invocations while leaving the existing Runner interface unchanged.
type ACPXPromptRequest struct {
	ExecuteRequest
	Session string
}

// BatchFailureError is returned when acpx reports an Agent/Batch-level
// failure. Callers can settle the Batch and continue under ADR 0010.
type BatchFailureError struct {
	ExitCode int
	Reason   string
	Stderr   string
	Err      error
}

func (err *BatchFailureError) Error() string {
	if err == nil {
		return ""
	}
	message := "Agent Batch failed"
	if err.ExitCode != 0 {
		message = fmt.Sprintf("%s after acpx exited with code %d", message, err.ExitCode)
	}
	if err.Reason != "" {
		message += ": " + err.Reason
	}
	if err.Err != nil {
		message += ": " + err.Err.Error()
		return message
	}
	if detail := strings.TrimSpace(err.Stderr); detail != "" {
		message += ": " + detail
	}
	return message
}

func (err *BatchFailureError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

// InfrastructureError is returned when acpx reports a Roundfix/environment
// bug rather than a user-resolvable Batch failure.
type InfrastructureError struct {
	ExitCode int
	Reason   string
	Stderr   string
}

func (err *InfrastructureError) Error() string {
	if err == nil {
		return ""
	}
	message := "acpx infrastructure error"
	if err.ExitCode != 0 {
		message = fmt.Sprintf("%s after exit code %d", message, err.ExitCode)
	}
	if err.Reason != "" {
		message += ": " + err.Reason
	}
	if tail, truncated := infrastructureStderrTail(err.Stderr); tail != "" {
		message += infrastructureStderrDelimiter
		if truncated {
			message += infrastructureStderrTruncated
		}
		message += tail
	}
	return message
}

func infrastructureStderrTail(stderr string) (string, bool) {
	tail := strings.TrimSpace(stderr)
	if tail == "" {
		return "", false
	}
	truncated := false
	lines := strings.Split(tail, "\n")
	if len(lines) > infrastructureStderrTailLines {
		lines = lines[len(lines)-infrastructureStderrTailLines:]
		tail = strings.Join(lines, "\n")
		truncated = true
	}
	tail = strings.TrimSpace(tail)
	if len(tail) > infrastructureStderrTailBytes {
		tail = tail[len(tail)-infrastructureStderrTailBytes:]
		truncated = true
	}
	tail = strings.TrimSpace(tail)
	return tail, truncated
}

// ModelNotAdvertisedError is returned when acpx reports that the selected
// Agent Model is absent from the runtime-advertised model list.
type ModelNotAdvertisedError struct {
	Runtime    string
	Model      string
	Advertised []string
	Err        error
}

func (err *ModelNotAdvertisedError) Error() string {
	if err == nil {
		return ""
	}
	runtime := strings.TrimSpace(err.Runtime)
	model := strings.TrimSpace(err.Model)
	message := fmt.Sprintf("Agent Model %q not advertised by runtime %q", model, runtime)
	if len(err.Advertised) > 0 {
		message += "; advertised Agent Models: " + strings.Join(err.Advertised, ", ")
	}
	message += "; recovery: " + err.RecoveryAction()
	return message
}

func (err *ModelNotAdvertisedError) RecoveryAction() string {
	runtime := ""
	if err != nil {
		runtime = strings.TrimSpace(err.Runtime)
	}
	return selectionRecoveryAction(runtime, false)
}

func (err *ModelNotAdvertisedError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

func (err *ModelNotAdvertisedError) Classification() string {
	return SelectionModelNotAdvertised
}

type ACPXProbeError struct {
	Command        string
	FoundVersion   string
	MinimumVersion string
	Missing        bool
	Err            error
}

func (err ACPXProbeError) Error() string {
	install := acpxInstallCommand()
	if err.Missing {
		return fmt.Sprintf("acpx is required but was not found on PATH; install it with: %s", install)
	}
	if err.FoundVersion != "" {
		return fmt.Sprintf("acpx version unsupported: found %s, but Roundfix requires %s or newer; upgrade with: %s", err.FoundVersion, err.MinimumVersion, install)
	}
	message := "acpx --version failed"
	if err.Err != nil {
		message = fmt.Sprintf("%s: %v", message, err.Err)
	}
	return fmt.Sprintf("%s; install it with: %s", message, install)
}

func (err ACPXProbeError) Unwrap() error {
	return err.Err
}

type AdapterProbeError struct {
	Command    string
	Executable string
	Install    string
	Err        error
}

func (err AdapterProbeError) Error() string {
	command := strings.TrimSpace(err.Command)
	if command == "" {
		command = "adapter"
	}
	executable := strings.TrimSpace(err.Executable)
	if executable == "" || executable == command {
		return fmt.Sprintf("%s is required but was not found on PATH; install it with: %s", command, err.InstallCommand())
	}
	return fmt.Sprintf("effective adapter command %q requires %s, but it was not found on PATH; install it with: %s", command, executable, err.InstallCommand())
}

func (err AdapterProbeError) InstallCommand() string {
	if install := strings.TrimSpace(err.Install); install != "" {
		return install
	}
	return adapterInstallCommand(err.Command)
}

func (err AdapterProbeError) Unwrap() error {
	return err.Err
}

// AdapterLineageError reports an adapter that cannot prove the official
// package lineage. Raw adapter output is intentionally excluded.
type AdapterLineageError struct {
	Runtime         string
	Command         string
	Package         string
	Version         string
	RequiredPackage string
	RequiredVersion string
	Install         string
	Legacy          bool
	Err             error
}

func (err *AdapterLineageError) Error() string {
	if err == nil {
		return ""
	}
	runtimeID := adapterErrorRuntime(err.Runtime)
	requiredPackage := adapterErrorRequiredPackage(err.RequiredPackage)
	requiredVersion := adapterErrorRequiredVersion(err.RequiredVersion)
	message := fmt.Sprintf("effective %s adapter command %q did not prove required package lineage %s", runtimeDisplayName(runtimeID), strings.TrimSpace(err.Command), requiredPackage)
	if err.Legacy || isLegacyAdapterPackage(err.Package) {
		message = fmt.Sprintf("effective %s adapter command %q reported legacy package %s", runtimeDisplayName(runtimeID), strings.TrimSpace(err.Command), err.Package)
	} else if packageName := strings.TrimSpace(err.Package); packageName != "" && packageName != requiredPackage {
		message = fmt.Sprintf("effective %s adapter command %q reported unknown package %s", runtimeDisplayName(runtimeID), strings.TrimSpace(err.Command), packageName)
	}
	if strings.TrimSpace(err.Version) != "" {
		message += " version " + strings.TrimSpace(err.Version)
	}
	return fmt.Sprintf("%s; required %s %s or newer; update with: %s", message, requiredPackage, requiredVersion, err.InstallCommand())
}

func (err *AdapterLineageError) Classification() string {
	return AdapterLineageUnknown
}

func (err *AdapterLineageError) InstallCommand() string {
	if err != nil {
		if install := strings.TrimSpace(err.Install); install != "" {
			return install
		}
		if strings.TrimSpace(err.RequiredPackage) != "" && strings.TrimSpace(err.RequiredVersion) != "" {
			return adapterPackageInstallCommand(err.RequiredPackage, err.RequiredVersion)
		}
	}
	return CodexAdapterInstallCommand()
}

func (err *AdapterLineageError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

// AdapterVersionError reports an official adapter below the supported
// compatibility floor.
type AdapterVersionError struct {
	Runtime         string
	Command         string
	Package         string
	FoundVersion    string
	RequiredPackage string
	RequiredVersion string
	Install         string
}

func (err *AdapterVersionError) Error() string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("effective %s adapter command %q reported package %s version %s; required version %s or newer; update with: %s", runtimeDisplayName(adapterErrorRuntime(err.Runtime)), strings.TrimSpace(err.Command), strings.TrimSpace(err.Package), strings.TrimSpace(err.FoundVersion), adapterErrorRequiredVersion(err.RequiredVersion), err.InstallCommand())
}

func (err *AdapterVersionError) Classification() string {
	return AdapterVersionUnsupported
}

func (err *AdapterVersionError) InstallCommand() string {
	if err != nil {
		if install := strings.TrimSpace(err.Install); install != "" {
			return install
		}
		if strings.TrimSpace(err.RequiredPackage) != "" && strings.TrimSpace(err.RequiredVersion) != "" {
			return adapterPackageInstallCommand(err.RequiredPackage, err.RequiredVersion)
		}
	}
	return CodexAdapterInstallCommand()
}

type SelectionPreflightError struct {
	Runtime         string
	Model           string
	ReasoningEffort string
	Operation       string
	Err             error
}

func (err *SelectionPreflightError) Error() string {
	if err == nil {
		return ""
	}
	runtime := strings.TrimSpace(err.Runtime)
	operation := strings.TrimSpace(err.Operation)
	if operation == "" {
		operation = "validate selection"
	}
	message := fmt.Sprintf("agent selection unavailable for runtime %q with model %q and reasoning %q during %s", runtime, err.Model, err.ReasoningEffort, operation)
	if err.Err != nil {
		message += ": " + err.Err.Error()
	}
	var classified interface{ Classification() string }
	if errors.As(err.Err, &classified) {
		return message
	}
	message += "; " + selectionRecoveryGuidance(runtime, true)
	return message
}

func (err *SelectionPreflightError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

type AgentSessionCleanupError struct {
	Session string
	Err     error
}

func (err *AgentSessionCleanupError) Error() string {
	if err == nil {
		return ""
	}
	message := fmt.Sprintf("close disposable Agent Session %q", err.Session)
	if err.Err != nil {
		message += ": " + err.Err.Error()
	}
	return message + "; recovery: rerun Agent Selection readiness after the Session can be closed"
}

func (err *AgentSessionCleanupError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

func (err *AgentSessionCleanupError) Classification() string {
	return SessionCleanupFailed
}

type acpxJSONRPCMessage struct {
	Method string            `json:"method"`
	Params json.RawMessage   `json:"params"`
	Result json.RawMessage   `json:"result"`
	Error  *acpxJSONRPCError `json:"error"`
}

type acpxJSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type acpxStreamResult struct {
	output             string
	stopReason         string
	promptResultParsed bool
	err                error
}

func (runner ACPXRunner) Probe(ctx context.Context, req ProbeRequest) error {
	if err := runner.probeACPX(ctx); err != nil {
		return err
	}
	workDir := strings.TrimSpace(req.WorkDir)
	if workDir == "" {
		return nil
	}
	_, err := runner.ProveExactSelection(ctx, ProbeRequest{Runtime: req.Runtime, WorkDir: workDir})
	return err
}

func CheckAdapter(ctx context.Context, runtime RuntimeSpec) (AdapterEvidence, error) {
	return checkAdapter(ctx, runtime, os.Environ())
}

func checkAdapter(ctx context.Context, runtime RuntimeSpec, environment []string) (AdapterEvidence, error) {
	if err := ctx.Err(); err != nil {
		return AdapterEvidence{}, err
	}
	invocation, err := resolveAdapterInvocation(runtime, environment)
	if err != nil {
		return AdapterEvidence{}, err
	}
	runtimeID := strings.TrimSuffix(strings.TrimSpace(runtime.ID), "-custom")
	contract, hasLineageContract := adapterLineageContracts[runtimeID]
	if _, err := exec.LookPath(invocation.executable()); err != nil {
		install := ""
		if hasLineageContract {
			install = contract.installCommand()
		}
		return AdapterEvidence{}, AdapterProbeError{
			Command:    invocation.display(),
			Executable: invocation.executable(),
			Install:    install,
			Err:        err,
		}
	}
	evidence := AdapterEvidence{Command: invocation.display()}
	if !hasLineageContract {
		return evidence, nil
	}
	return inspectAdapter(ctx, invocation, contract, environment)
}

// resolveAdapterCommand returns the adapter command acpx will spawn for the
// selected runtime: stdio overrides first, then acpx's agents map, then
// Roundfix's built-in defaults.
func resolveAdapterCommand(runtime RuntimeSpec) (string, error) {
	return resolveAdapterCommandWithEnv(runtime, os.Environ())
}

func resolveAdapterCommandWithEnv(runtime RuntimeSpec, environment []string) (string, error) {
	invocation, err := resolveAdapterInvocation(runtime, environment)
	if err != nil {
		return "", err
	}
	return invocation.display(), nil
}

type adapterInvocation struct {
	argv []string
}

func (invocation adapterInvocation) executable() string {
	if len(invocation.argv) == 0 {
		return ""
	}
	return invocation.argv[0]
}

func (invocation adapterInvocation) display() string {
	return strings.Join(invocation.argv, " ")
}

func resolveAdapterInvocation(runtime RuntimeSpec, environment []string) (adapterInvocation, error) {
	if runtime.Protocol == ProtocolStdio {
		if invocation := newAdapterInvocation(runtime.Command, nil); len(invocation.argv) > 0 {
			return invocation, nil
		}
	}
	runtimeID := strings.TrimSpace(runtime.ID)
	if invocation, ok := configuredAdapterInvocation(runtimeID, environment); ok {
		return invocation, nil
	}
	if command, ok := defaultAdapterCommands[runtimeID]; ok {
		return newAdapterInvocation(command, nil), nil
	}
	return adapterInvocation{}, fmt.Errorf("unsupported Agent %q; supported values: codex, claude, opencode", runtimeID)
}

func configuredAdapterInvocation(runtimeID string, environment []string) (adapterInvocation, bool) {
	path, err := acpxConfigPath(environment)
	if err != nil {
		return adapterInvocation{}, false
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return adapterInvocation{}, false
	}
	var config struct {
		Agents map[string]struct {
			Command string   `json:"command"`
			Args    []string `json:"args"`
		} `json:"agents"`
	}
	if err := json.Unmarshal(content, &config); err != nil {
		return adapterInvocation{}, false
	}
	agentConfig, ok := config.Agents[runtimeID]
	if !ok {
		return adapterInvocation{}, false
	}
	invocation := newAdapterInvocation(agentConfig.Command, agentConfig.Args)
	return invocation, len(invocation.argv) > 0
}

func acpxConfigPath(environment []string) (string, error) {
	homeDir := environmentValue(environment, "HOME")
	if homeDir == "" {
		homeDir = environmentValue(environment, "USERPROFILE")
	}
	if homeDir == "" {
		return "", errors.New("explicit environment does not define HOME or USERPROFILE")
	}
	return filepath.Join(homeDir, ".acpx", "config.json"), nil
}

func adapterBinary(command string) string {
	fields := strings.Fields(strings.TrimSpace(command))
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func newAdapterInvocation(command string, args []string) adapterInvocation {
	argv := strings.Fields(strings.TrimSpace(command))
	for _, arg := range args {
		if trimmed := strings.TrimSpace(arg); trimmed != "" {
			argv = append(argv, trimmed)
		}
	}
	return adapterInvocation{argv: argv}
}

func adapterInstallCommand(command string) string {
	command = strings.TrimSpace(command)
	invocation := newAdapterInvocation(command, nil)
	if adapterBinary(command) == "codex-acp" || invocationNamesPackage(invocation, CodexAdapterPackage) {
		return CodexAdapterInstallCommand()
	}
	if binary := adapterBinary(command); binary == "claude-code-acp" || binary == "claude-agent-acp" {
		return ClaudeAdapterInstallCommand()
	}
	for _, packageName := range append([]string{ClaudeAdapterPackage}, adapterLineageContracts["claude"].LegacyPackages...) {
		if invocationNamesPackage(invocation, packageName) {
			return ClaudeAdapterInstallCommand()
		}
	}
	if install, ok := adapterInstallCommands[adapterBinary(command)]; ok {
		return install
	}
	if command == "" {
		return "install the adapter and ensure it is on PATH"
	}
	return "install " + command + " and ensure it is on PATH"
}

func inspectAdapter(ctx context.Context, invocation adapterInvocation, contract adapterLineageContract, environment []string) (AdapterEvidence, error) {
	evidence := AdapterEvidence{Command: invocation.display()}
	resolvedPackage := resolveAdapterPackage(invocation, contract)
	args := append([]string(nil), invocation.argv[1:]...)
	args = append(args, "--version")
	cmd := exec.CommandContext(ctx, invocation.executable(), args...)
	cmd.Env = append([]string(nil), environment...)
	var stdout boundedAdapterOutput
	cmd.Stdout = &stdout
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return AdapterEvidence{}, ctxErr
		}
		return AdapterEvidence{}, newAdapterLineageError(contract, evidence.Command, resolvedPackage, "", err)
	}
	if stdout.truncated {
		return AdapterEvidence{}, newAdapterLineageError(contract, evidence.Command, resolvedPackage, "", nil)
	}
	fields := strings.Fields(stdout.String())
	switch {
	case len(fields) == 2:
		evidence.Package = fields[0]
		evidence.Version = fields[1]
	case len(fields) == 1 && contract.VersionOnly:
		evidence.Package = resolvedPackage
		evidence.Version = fields[0]
	default:
		return AdapterEvidence{}, newAdapterLineageError(contract, evidence.Command, resolvedPackage, "", nil)
	}
	_, validVersion := parseAdapterVersion(evidence.Version)
	if evidence.Package != contract.Package {
		version := evidence.Version
		if !validVersion {
			version = ""
		}
		return AdapterEvidence{}, newAdapterLineageError(contract, evidence.Command, evidence.Package, version, nil)
	}
	if !validVersion && contract.VersionOnly {
		return AdapterEvidence{}, newAdapterLineageError(contract, evidence.Command, evidence.Package, "", nil)
	}
	if !validVersion || compareAdapterVersions(evidence.Version, contract.PinnedVersion) < 0 {
		return AdapterEvidence{}, newAdapterVersionError(contract, evidence)
	}
	return evidence, nil
}

func newAdapterLineageError(contract adapterLineageContract, command string, packageName string, version string, err error) *AdapterLineageError {
	return &AdapterLineageError{
		Runtime:         contract.RuntimeID,
		Command:         command,
		Package:         packageName,
		Version:         version,
		RequiredPackage: contract.Package,
		RequiredVersion: contract.PinnedVersion,
		Install:         contract.installCommand(),
		Legacy:          contract.isLegacy(packageName),
		Err:             err,
	}
}

func newAdapterVersionError(contract adapterLineageContract, evidence AdapterEvidence) *AdapterVersionError {
	return &AdapterVersionError{
		Runtime:         contract.RuntimeID,
		Command:         evidence.Command,
		Package:         evidence.Package,
		FoundVersion:    evidence.Version,
		RequiredPackage: contract.Package,
		RequiredVersion: contract.PinnedVersion,
		Install:         contract.installCommand(),
	}
}

func (contract adapterLineageContract) installCommand() string {
	return adapterPackageInstallCommand(contract.Package, contract.PinnedVersion)
}

func (contract adapterLineageContract) isLegacy(packageName string) bool {
	for _, legacyPackage := range contract.LegacyPackages {
		if packageName == legacyPackage {
			return true
		}
	}
	return false
}

func adapterPackageInstallCommand(packageName string, version string) string {
	return "npm install -g " + strings.TrimSpace(packageName) + "@" + strings.TrimSpace(version)
}

func resolveAdapterPackage(invocation adapterInvocation, contract adapterLineageContract) string {
	for _, packageName := range append([]string{contract.Package}, contract.LegacyPackages...) {
		if invocationNamesPackage(invocation, packageName) {
			return packageName
		}
	}
	executablePath, err := exec.LookPath(invocation.executable())
	if err != nil {
		return ""
	}
	resolvedPath, err := filepath.EvalSymlinks(executablePath)
	if err != nil {
		return ""
	}
	for _, packageName := range append([]string{contract.Package}, contract.LegacyPackages...) {
		if resolvedPathNamesPackage(resolvedPath, packageName) {
			return packageName
		}
	}
	return ""
}

func invocationNamesPackage(invocation adapterInvocation, packageName string) bool {
	for _, arg := range invocation.argv {
		if arg == packageName {
			return true
		}
		if strings.HasPrefix(arg, packageName+"@") && len(arg) > len(packageName)+1 {
			return true
		}
	}
	return false
}

func resolvedPathNamesPackage(path string, packageName string) bool {
	path = "/" + strings.Trim(filepath.ToSlash(filepath.Clean(path)), "/") + "/"
	packageSegment := "/node_modules/" + strings.Trim(packageName, "/") + "/"
	return strings.Contains(path, packageSegment)
}

func adapterErrorRuntime(runtimeID string) string {
	runtimeID = strings.TrimSpace(runtimeID)
	if runtimeID == "" {
		return "codex"
	}
	return runtimeID
}

func adapterErrorRequiredPackage(packageName string) string {
	if packageName = strings.TrimSpace(packageName); packageName != "" {
		return packageName
	}
	return CodexAdapterPackage
}

func adapterErrorRequiredVersion(version string) string {
	if version = strings.TrimSpace(version); version != "" {
		return version
	}
	return PinnedCodexAdapterVersion
}

func runtimeDisplayName(runtimeID string) string {
	if runtimeID == "" {
		return "Adapter"
	}
	return strings.ToUpper(runtimeID[:1]) + runtimeID[1:]
}

func isLegacyAdapterPackage(packageName string) bool {
	for _, contract := range adapterLineageContracts {
		if contract.isLegacy(packageName) {
			return true
		}
	}
	return false
}

type boundedAdapterOutput struct {
	content   []byte
	truncated bool
}

func (output *boundedAdapterOutput) Write(content []byte) (int, error) {
	remaining := adapterProbeOutputLimit - len(output.content)
	if remaining > 0 {
		length := len(content)
		if length > remaining {
			length = remaining
		}
		output.content = append(output.content, content[:length]...)
	}
	if len(content) > remaining {
		output.truncated = true
	}
	return len(content), nil
}

func (output *boundedAdapterOutput) String() string {
	return strings.TrimSpace(string(output.content))
}

func compareAdapterVersions(found string, required string) int {
	foundParts, foundOK := parseAdapterVersion(found)
	requiredParts, requiredOK := parseAdapterVersion(required)
	if !foundOK || !requiredOK {
		return -1
	}
	for index := range foundParts {
		if foundParts[index] < requiredParts[index] {
			return -1
		}
		if foundParts[index] > requiredParts[index] {
			return 1
		}
	}
	return 0
}

// SupportsACPXVersion reports whether version satisfies Roundfix's minimum
// tested acpx version. Malformed versions fail closed.
func SupportsACPXVersion(version string) bool {
	return compareAdapterVersions(version, MinimumACPXVersion) >= 0
}

func parseAdapterVersion(version string) ([3]int, bool) {
	var parsed [3]int
	parts := strings.Split(strings.TrimPrefix(strings.TrimSpace(version), "v"), ".")
	if len(parts) != len(parsed) {
		return parsed, false
	}
	for index, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil || value < 0 {
			return parsed, false
		}
		parsed[index] = value
	}
	return parsed, true
}

func (runner ACPXRunner) probeACPX(ctx context.Context) error {
	command := runner.command()
	if _, err := exec.LookPath(command); err != nil {
		return ACPXProbeError{Command: command, MinimumVersion: MinimumACPXVersion, Missing: true, Err: err}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, command, "--version")
	cmd.Env = runner.baseEnv()
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		detail := strings.TrimSpace(string(output))
		if detail == "" {
			detail = err.Error()
		}
		return ACPXProbeError{Command: command, MinimumVersion: MinimumACPXVersion, Err: errors.New(detail)}
	}
	foundVersion := strings.TrimSpace(string(output))
	if !SupportsACPXVersion(foundVersion) {
		return ACPXProbeError{Command: command, FoundVersion: displayACPXVersion(foundVersion), MinimumVersion: MinimumACPXVersion}
	}
	return nil
}

func disposablePreflightSessionName() (string, error) {
	var entropy [8]byte
	if _, err := rand.Read(entropy[:]); err != nil {
		return "", fmt.Errorf("generate disposable Agent Session name: %w", err)
	}
	return fmt.Sprintf("%s%x", acpxPreflightSessionPrefix, entropy[:]), nil
}

func (runner ACPXRunner) applyDisposableSelection(ctx context.Context, runtime RuntimeSpec, sessionName string, workDir string, codexEnv []string) error {
	args, err := acpxEnsureArgs(runtime, sessionName, workDir)
	if err != nil {
		return err
	}
	if err := runner.runACPXCommandWithEnv(ctx, args, codexEnv); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return selectionPreflightError(runtime, "set model", fmt.Errorf("ensure disposable acpx Agent Session %q with model %q: %w", sessionName, strings.TrimSpace(runtime.Model), classifyModelNotAdvertised(runtime, err)))
	}
	value := strings.TrimSpace(runtime.ReasoningEffort)
	if value == "" {
		return nil
	}
	key, err := acpxReasoningEffortConfigKey(runtime)
	if err != nil {
		return err
	}
	args, err = acpxSetConfigArgs(runtime, key, value, sessionName, workDir)
	if err != nil {
		return err
	}
	if err := runner.runACPXCommandWithEnv(ctx, args, codexEnv); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return selectionPreflightError(runtime, "set "+key, fmt.Errorf("set disposable acpx Agent Session %s %q: %w", key, value, classifyModelNotAdvertised(runtime, err)))
	}
	return nil
}

func selectionPreflightError(runtime RuntimeSpec, operation string, err error) error {
	return &SelectionPreflightError{
		Runtime:         strings.TrimSpace(runtime.ID),
		Model:           strings.TrimSpace(runtime.Model),
		ReasoningEffort: strings.TrimSpace(runtime.ReasoningEffort),
		Operation:       strings.TrimSpace(operation),
		Err:             err,
	}
}

func selectionRecoveryGuidance(runtime string, includeReasoning bool) string {
	return "recovery: " + selectionRecoveryAction(runtime, includeReasoning)
}

func selectionRecoveryAction(runtime string, includeReasoning bool) string {
	runtime = strings.TrimSpace(runtime)
	message := "update the ACP Runtime or adapter"
	if includeReasoning {
		message += ", choose supported Agent Model and Default Reasoning Effort values, choose an advertised Agent Model"
	} else {
		message += ", choose an advertised Agent Model"
	}
	message += ", or pass a one-Run --model override"
	if includeReasoning {
		message += " with --reasoning-effort when needed"
	}
	if includeReasoning && runtime != "" {
		message += fmt.Sprintf(`, or set runtimes.%s.reasoning_effort "" when the model manages reasoning`, runtime)
	}
	return message
}

func classifyModelNotAdvertised(runtime RuntimeSpec, err error) error {
	var infraErr *InfrastructureError
	if !errors.As(err, &infraErr) {
		return err
	}
	modelErr, ok := modelNotAdvertisedFromStderr(runtime, infraErr.Stderr, infraErr)
	if !ok {
		return err
	}
	return modelErr
}

func modelNotAdvertisedFromStderr(runtime RuntimeSpec, stderr string, cause error) (*ModelNotAdvertisedError, bool) {
	if !strings.Contains(stderr, "did not advertise that model") {
		return nil, false
	}
	model := rejectedModelFromACPXStderr(stderr)
	if model == "" {
		model = strings.TrimSpace(runtime.Model)
	}
	return &ModelNotAdvertisedError{
		Runtime:    strings.TrimSpace(runtime.ID),
		Model:      model,
		Advertised: advertisedModelsFromACPXStderr(stderr),
		Err:        cause,
	}, true
}

func rejectedModelFromACPXStderr(stderr string) string {
	const marker = "--model"
	index := strings.Index(stderr, marker)
	if index < 0 {
		return ""
	}
	rest := strings.TrimSpace(stderr[index+len(marker):])
	rest = strings.TrimSpace(strings.TrimPrefix(rest, "="))
	if rest == "" {
		return ""
	}
	if rest[0] == '"' || rest[0] == '\'' {
		quote := rest[0]
		if end := strings.IndexByte(rest[1:], quote); end >= 0 {
			return strings.TrimSpace(rest[1 : 1+end])
		}
		return strings.Trim(strings.TrimSpace(rest[1:]), `"' :;,`)
	}
	if end := strings.IndexAny(rest, " \t\r\n:;,"); end >= 0 {
		rest = rest[:end]
	}
	return strings.Trim(strings.TrimSpace(rest), `"'`)
}

func advertisedModelsFromACPXStderr(stderr string) []string {
	const marker = "Available models:"
	index := strings.Index(stderr, marker)
	if index < 0 {
		return nil
	}
	value := stderr[index+len(marker):]
	if line, _, ok := strings.Cut(value, "\n"); ok {
		value = line
	}
	value = strings.Trim(strings.TrimSpace(value), "[]")
	if value == "" {
		return nil
	}
	var fields []string
	if strings.Contains(value, ",") {
		fields = strings.Split(value, ",")
	} else {
		fields = strings.Fields(value)
	}
	models := make([]string, 0, len(fields))
	for _, field := range fields {
		model := strings.Trim(strings.TrimSpace(field), `"'[]`)
		if model != "" {
			models = append(models, model)
		}
	}
	return models
}

func (runner ACPXRunner) closeDisposableSession(ctx context.Context, runtime RuntimeSpec, sessionName string, workDir string) error {
	defer func() {
		(&runner).clearSessionState(sessionName)
	}()
	if err := runner.CloseSession(ctx, runtime, SessionRef{Name: sessionName, WorkDir: workDir}); err != nil {
		return &AgentSessionCleanupError{Session: sessionName, Err: err}
	}
	return nil
}

func (runner *ACPXRunner) Run(ctx context.Context, req ExecuteRequest, sink runevent.Sink) (ExecuteResult, error) {
	result := ExecuteResult{LogPath: req.LogPath}
	if err := runner.PrepareSession(ctx, req, sink); err != nil {
		return result, err
	}
	if err := runner.publishStatus(ctx, req, sink, AgentWorkStartedStatus); err != nil {
		return result, err
	}
	return runner.RunPrepared(ctx, req, sink)
}

func (runner *ACPXRunner) PrepareSession(ctx context.Context, req ExecuteRequest, sink runevent.Sink) error {
	if err := validateRuntimeSelection(req.Runtime); err != nil {
		return err
	}
	if _, err := runner.codexEnvForSession(ctx, req.Runtime, req.Session.Name); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return selectionPreflightError(req.Runtime, "start runtime", err)
	}
	return runner.ensureSession(ctx, req, sink)
}

func (runner *ACPXRunner) RunPrepared(ctx context.Context, req ExecuteRequest, sink runevent.Sink) (ExecuteResult, error) {
	if assignment, ok := runner.sessionSelection(req.Session.Name); ok {
		req.Runtime.Model = assignment.AdapterModel
	}
	return runner.RunPrompt(ctx, ACPXPromptRequest{
		ExecuteRequest: req,
		Session:        req.Session.Name,
	}, sink)
}

func (runner *ACPXRunner) EndSession(ctx context.Context, runtime RuntimeSpec, session SessionRef) error {
	sessionName := strings.TrimSpace(session.Name)
	if sessionName == "" {
		return nil
	}
	if err := runner.CloseSession(ctx, runtime, session); err != nil {
		runner.warningf("close acpx Agent Session %q: %v", sessionName, err)
	}
	runner.clearSessionState(sessionName)
	return nil
}

func (runner *ACPXRunner) CloseSession(ctx context.Context, runtime RuntimeSpec, session SessionRef) error {
	sessionName := strings.TrimSpace(session.Name)
	if sessionName == "" {
		return nil
	}
	args, err := acpxCloseArgs(runtime, sessionName, session.WorkDir)
	if err != nil {
		return fmt.Errorf("build acpx close command for Agent Session %q: %w", sessionName, err)
	}
	codexEnv, err := runner.codexEnvForSession(ctx, runtime, sessionName)
	if err != nil {
		return err
	}
	if err := runner.runACPXCommandWithEnv(ctx, args, codexEnv); err != nil {
		return fmt.Errorf("close acpx Agent Session %q: %w", sessionName, err)
	}
	return nil
}

func (runner ACPXRunner) ListRoundfixSessions(ctx context.Context, runtime RuntimeSpec, workDir string) ([]RoundfixSession, error) {
	args, err := acpxListSessionsArgs(runtime, workDir)
	if err != nil {
		return nil, fmt.Errorf("build acpx sessions list command: %w", err)
	}
	output, err := runner.runACPXCommandOutput(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("list acpx Agent Sessions: %w", err)
	}
	return ParseRoundfixSessions(output), nil
}

func (runner *ACPXRunner) CancelSession(ctx context.Context, runtime RuntimeSpec, session SessionRef) error {
	sessionName := strings.TrimSpace(session.Name)
	if sessionName == "" {
		return errors.New("agent session name is required")
	}
	args, err := acpxCancelArgs(runtime, sessionName, session.WorkDir)
	if err != nil {
		return fmt.Errorf("build acpx cancel command for Agent Session %q: %w", sessionName, err)
	}
	codexEnv, err := runner.codexEnvForSession(ctx, runtime, sessionName)
	if err != nil {
		return err
	}
	if err := runner.runACPXCommandWithEnv(ctx, args, codexEnv); err != nil {
		return fmt.Errorf("cancel acpx Agent Session %q: %w", sessionName, err)
	}
	return nil
}

func (runner *ACPXRunner) ensureSession(ctx context.Context, req ExecuteRequest, sink runevent.Sink) error {
	sessionName := strings.TrimSpace(req.Session.Name)
	if sessionName == "" {
		return errors.New("agent session is required")
	}
	workDir := strings.TrimSpace(req.GitRoot)
	if workDir == "" {
		return errors.New("agent working directory is required")
	}
	if runner.sessionEnsured(sessionName) {
		return nil
	}
	adapter, err := checkAdapter(ctx, req.Runtime, runner.baseEnv())
	if err != nil {
		return err
	}
	codexEnv, err := runner.codexEnvForSession(ctx, req.Runtime, sessionName)
	if err != nil {
		return err
	}
	session := SessionRef{Name: sessionName, WorkDir: workDir}
	capabilities, err := runner.startSessionSelection(ctx, req.Runtime, session, adapter, codexEnv, false)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return selectionPreflightError(req.Runtime, "apply advertised selection", err)
	}
	proof, err := runner.applySessionSelection(ctx, SessionSelectionRequest{Runtime: req.Runtime, Session: session, Capabilities: capabilities}, codexEnv)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return selectionPreflightError(req.Runtime, "apply advertised selection", err)
	}
	if err := runner.applyFullAccess(ctx, req, sink, codexEnv); err != nil {
		return err
	}
	if err := runner.publishStatus(ctx, req, sink, AgentSessionStartedStatus); err != nil {
		return err
	}
	runner.markSessionEnsured(sessionName, proof.Assignment)
	return nil
}

func (runner *ACPXRunner) applyFullAccess(ctx context.Context, req ExecuteRequest, sink runevent.Sink, codexEnv []string) error {
	mode := strings.TrimSpace(req.Runtime.FullAccessMode)
	if mode == "" {
		return nil
	}
	sessionName := strings.TrimSpace(req.Session.Name)
	args, err := acpxSetModeArgs(req.Runtime, mode, sessionName, req.GitRoot)
	if err != nil {
		return err
	}
	if err := runner.runACPXCommandWithEnv(ctx, args, codexEnv); err != nil {
		return fmt.Errorf("set acpx Agent Session mode %q: %w", mode, err)
	}
	if req.Runtime.ID != "codex" || mode != "full-access" {
		return nil
	}
	args, err = acpxSetConfigArgs(req.Runtime, acpxCodexSandboxModeKey, acpxCodexFullAccessSandbox, sessionName, req.GitRoot)
	if err != nil {
		return err
	}
	if err := runner.runACPXCommandWithEnv(ctx, args, codexEnv); err != nil {
		if isCodexSandboxUnavailable(err) {
			runner.warningf("codex full-access sandbox preset unavailable for Agent Session %q: %v", sessionName, err)
			if publishErr := runner.publishStatus(ctx, req, sink, acpxCodexSandboxUnavailable); publishErr != nil {
				return publishErr
			}
			return nil
		}
		return fmt.Errorf("set acpx Codex sandbox preset %q: %w", acpxCodexFullAccessSandbox, err)
	}
	return nil
}

func (runner *ACPXRunner) RunPrompt(ctx context.Context, req ACPXPromptRequest, sink runevent.Sink) (ExecuteResult, error) {
	result := ExecuteResult{LogPath: req.LogPath}
	if err := validateACPXPromptRequest(req); err != nil {
		return result, err
	}
	if sink == nil {
		sink = runevent.Discard
	}
	logWriter := io.Writer(io.Discard)
	if strings.TrimSpace(req.LogPath) != "" {
		if err := os.MkdirAll(filepath.Dir(req.LogPath), 0o755); err != nil {
			return result, fmt.Errorf("create Agent log directory: %w", err)
		}
		logFile, err := os.OpenFile(req.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return result, fmt.Errorf("open Agent log %q: %w", req.LogPath, err)
		}
		defer func() {
			_ = logFile.Close()
		}()
		logWriter = logFile
	}

	if err := ctx.Err(); err != nil {
		return result, StopError{LogPath: req.LogPath, Err: err}
	}
	codexEnv, err := runner.codexEnvForSession(ctx, req.Runtime, req.Session)
	if err != nil {
		return result, err
	}
	args, err := acpxPromptArgs(req)
	if err != nil {
		return result, err
	}
	cmd := exec.Command(runner.command(), args...)
	cmd.Env = runner.commandEnv(codexEnv)
	cmd.Dir = req.GitRoot
	cmd.Stdin = strings.NewReader(req.Prompt)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return result, fmt.Errorf("open acpx stdout: %w", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return result, fmt.Errorf("start acpx prompt: %w", err)
	}

	streamCh := make(chan acpxStreamResult, 1)
	waitCh := make(chan error, 1)
	go func() {
		// Drain stdout to EOF before Wait: Wait closes the pipe, so calling
		// it while the reader is mid-read races into "file already closed"
		// and can drop tail output (observed on slow CI runners). Both
		// channels are buffered, so this goroutine never blocks on delivery.
		streamCh <- runner.readPromptStream(ctx, req.ExecuteRequest, sink, stdout, logWriter)
		waitCh <- cmd.Wait()
	}()

	var waitErr error
	select {
	case waitErr = <-waitCh:
	case <-ctx.Done():
		forceClosed := runner.cancelPrompt(context.WithoutCancel(ctx), req, waitCh, codexEnv)
		stream := waitForACPXStream(streamCh, stopGrace(req.StopGrace))
		result.Output = stream.output
		result.StopReason = stream.stopReason
		_ = runner.publishStatus(context.WithoutCancel(ctx), req.ExecuteRequest, sink, "stopped")
		return result, StopError{LogPath: req.LogPath, Output: result.Output, Killed: forceClosed, Err: ctx.Err()}
	}
	stream := <-streamCh
	result.Output = stream.output
	result.StopReason = stream.stopReason
	if waitErr != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return result, StopError{LogPath: req.LogPath, Output: result.Output, Err: ctxErr}
		}
		exitCode, ok := commandExitCode(waitErr)
		if !ok {
			return result, fmt.Errorf("wait for acpx prompt: %w", waitErr)
		}
		if stream.promptResultParsed && exitCode != 130 {
			result.TransportAnomaly = acpxTransportAnomaly(exitCode, stderr.String())
			return result, nil
		}
		return result, runner.mapExitCode(ctx, req.ExecuteRequest, sink, exitCode, stderr.String(), result.Output)
	}
	if stream.err != nil {
		return result, &BatchFailureError{Reason: stream.err.Error(), Stderr: stderr.String()}
	}
	if result.StopReason == "" {
		return result, &BatchFailureError{Reason: "missing session/prompt stop reason", Stderr: stderr.String()}
	}
	return result, nil
}

func (runner ACPXRunner) readPromptStream(ctx context.Context, req ExecuteRequest, sink runevent.Sink, stdout io.Reader, logFile io.Writer) acpxStreamResult {
	var output bytes.Buffer
	var stopReason string
	var promptResultParsed bool
	var streamErr error
	reader := bufio.NewReader(stdout)
	for {
		line, readErr := reader.ReadBytes('\n')
		if len(line) > 0 {
			if _, err := logFile.Write(line); err != nil && streamErr == nil {
				streamErr = fmt.Errorf("write Agent log: %w", err)
			}
			if _, err := output.Write(line); err != nil && streamErr == nil {
				streamErr = fmt.Errorf("capture acpx stdout: %w", err)
			}
			if err := runner.handleStdoutLine(ctx, req, sink, line, &stopReason, &promptResultParsed); err != nil && streamErr == nil {
				streamErr = err
			}
		}
		if readErr != nil {
			if !errors.Is(readErr, io.EOF) && streamErr == nil {
				streamErr = fmt.Errorf("read acpx stdout: %w", readErr)
			}
			break
		}
	}
	return acpxStreamResult{output: output.String(), stopReason: stopReason, promptResultParsed: promptResultParsed, err: streamErr}
}

func validateACPXPromptRequest(req ACPXPromptRequest) error {
	if strings.TrimSpace(req.Session) == "" {
		return errors.New("Agent Session is required")
	}
	if strings.TrimSpace(req.GitRoot) == "" {
		return errors.New("Agent working directory is required")
	}
	return nil
}

func validateRuntimeSelection(runtime RuntimeSpec) error {
	if strings.TrimSpace(runtime.Model) == "" {
		return errors.New("agent model is required")
	}
	if strings.TrimSpace(runtime.ReasoningEffort) == "" {
		return nil
	}
	if _, err := acpxReasoningEffortConfigKey(runtime); err != nil {
		return err
	}
	return nil
}

func (runner ACPXRunner) command() string {
	if strings.TrimSpace(runner.Command) != "" {
		return strings.TrimSpace(runner.Command)
	}
	return defaultACPXCommand
}

func (runner *ACPXRunner) codexEnvForSession(ctx context.Context, runtime RuntimeSpec, sessionName string) ([]string, error) {
	if strings.TrimSpace(runtime.ID) != "codex" || runtime.Protocol == ProtocolStdio {
		return nil, nil
	}
	key := strings.TrimSpace(sessionName)
	if key == "" {
		key = "codex"
	}
	unlock := runner.lockState()
	if runner.codexResolutions != nil {
		if resolution, ok := runner.codexResolutions[key]; ok {
			unlock()
			return resolution.env, nil
		}
	}
	unlock()
	dependencies := runner.codexSpawn
	if dependencies.getenv == nil {
		environment := runner.baseEnv()
		dependencies.getenv = func(key string) string {
			return environmentValue(environment, key)
		}
	}
	resolution, err := dependencies.resolve(ctx)
	if err != nil {
		return nil, err
	}
	unlock = runner.lockState()
	defer unlock()
	if runner.codexResolutions == nil {
		runner.codexResolutions = map[string]codexSpawnResolution{}
	}
	runner.codexResolutions[key] = resolution
	return resolution.env, nil
}

func (runner *ACPXRunner) sessionEnsured(sessionName string) bool {
	unlock := runner.lockState()
	defer unlock()
	_, ok := runner.ensuredSessions[sessionName]
	return ok
}

func (runner *ACPXRunner) markSessionEnsured(sessionName string, selection SelectionAssignment) {
	unlock := runner.lockState()
	defer unlock()
	if runner.ensuredSessions == nil {
		runner.ensuredSessions = map[string]struct{}{}
	}
	if runner.sessionSelections == nil {
		runner.sessionSelections = map[string]SelectionAssignment{}
	}
	runner.ensuredSessions[sessionName] = struct{}{}
	runner.sessionSelections[sessionName] = selection
}

func (runner *ACPXRunner) sessionSelection(sessionName string) (SelectionAssignment, bool) {
	unlock := runner.lockState()
	defer unlock()
	selection, ok := runner.sessionSelections[sessionName]
	return selection, ok
}

func (runner *ACPXRunner) clearSessionState(sessionName string) {
	unlock := runner.lockState()
	defer unlock()
	delete(runner.ensuredSessions, sessionName)
	delete(runner.sessionSelections, sessionName)
	delete(runner.codexResolutions, sessionName)
}

func (runner *ACPXRunner) lockState() func() {
	runner.stateMu.Lock()
	return runner.stateMu.Unlock
}

func acpxInstallCommand() string {
	return "npm install -g acpx@" + MinimumACPXVersion
}

func displayACPXVersion(version string) string {
	if strings.TrimSpace(version) == "" {
		return "<empty>"
	}
	return version
}

func (runner ACPXRunner) runACPXCommand(ctx context.Context, args []string) error {
	_, err := runner.runACPXCommandOutput(ctx, args)
	return err
}

func (runner ACPXRunner) runACPXCommandOutput(ctx context.Context, args []string) (string, error) {
	return runner.runACPXCommandOutputWithEnv(ctx, args, nil)
}

func (runner ACPXRunner) runACPXCommandWithEnv(ctx context.Context, args []string, env []string) error {
	_, err := runner.runACPXCommandOutputWithEnv(ctx, args, env)
	return err
}

func (runner ACPXRunner) runACPXCommandOutputWithEnv(ctx context.Context, args []string, env []string) (string, error) {
	cmd := exec.CommandContext(ctx, runner.command(), args...)
	cmd.Env = runner.commandEnv(env)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return stdout.String(), ctxErr
		}
		if exitCode, ok := commandExitCode(err); ok {
			return stdout.String(), &InfrastructureError{ExitCode: exitCode, Reason: "acpx command failed", Stderr: stderr.String()}
		}
		return stdout.String(), fmt.Errorf("run acpx command: %w", err)
	}
	return stdout.String(), nil
}

// acpx CLI grammar: --cwd, --format, --json-strict, --approve-all, --model,
// and --agent are program-level globals and must appear before the agent
// name / subcommand. Subcommands accept only their own options (prompt,
// cancel, set-mode, set take -s <session>; sessions ensure takes --name;
// sessions close takes the session name as a positional argument). Every
// session-scoped invocation carries the global --cwd so session resolution
// is deterministic regardless of the Roundfix process cwd.

func acpxPromptArgs(req ACPXPromptRequest) ([]string, error) {
	agentArgs, err := acpxAgentArgs(req.Runtime)
	if err != nil {
		return nil, err
	}
	args := []string{
		"--cwd", strings.TrimSpace(req.GitRoot),
		"--format", "json",
		"--json-strict",
		"--approve-all",
	}
	if model := strings.TrimSpace(req.Runtime.Model); model != "" {
		args = append(args, "--model", model)
	}
	args = append(args, agentArgs...)
	args = append(args, "prompt", "-s", strings.TrimSpace(req.Session), "-f", "-")
	return args, nil
}

func acpxEnsureArgs(runtime RuntimeSpec, sessionName string, workDir string) ([]string, error) {
	args, err := acpxGlobalArgsWithModel(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(args, "sessions", "ensure", "--name", strings.TrimSpace(sessionName)), nil
}

func acpxSetModeArgs(runtime RuntimeSpec, mode string, sessionName string, workDir string) ([]string, error) {
	args, err := acpxGlobalArgs(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(args, "set-mode", strings.TrimSpace(mode), "-s", strings.TrimSpace(sessionName)), nil
}

func acpxSetConfigArgs(runtime RuntimeSpec, key string, value string, sessionName string, workDir string) ([]string, error) {
	args, err := acpxGlobalArgs(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(args, "set", strings.TrimSpace(key), strings.TrimSpace(value), "-s", strings.TrimSpace(sessionName)), nil
}

func acpxCancelArgs(runtime RuntimeSpec, sessionName string, workDir string) ([]string, error) {
	args, err := acpxGlobalArgs(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(args, "cancel", "-s", strings.TrimSpace(sessionName)), nil
}

func acpxCloseArgs(runtime RuntimeSpec, sessionName string, workDir string) ([]string, error) {
	args, err := acpxGlobalArgs(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(args, "sessions", "close", strings.TrimSpace(sessionName)), nil
}

func acpxListSessionsArgs(runtime RuntimeSpec, workDir string) ([]string, error) {
	args, err := acpxGlobalArgs(runtime, workDir)
	if err != nil {
		return nil, err
	}
	return append(args, "sessions", "list"), nil
}

func acpxGlobalArgs(runtime RuntimeSpec, workDir string) ([]string, error) {
	workDir = strings.TrimSpace(workDir)
	if workDir == "" {
		return nil, errors.New("Agent working directory is required")
	}
	agentArgs, err := acpxAgentArgs(runtime)
	if err != nil {
		return nil, err
	}
	return append([]string{"--cwd", workDir}, agentArgs...), nil
}

func acpxGlobalArgsWithModel(runtime RuntimeSpec, workDir string) ([]string, error) {
	workDir = strings.TrimSpace(workDir)
	if workDir == "" {
		return nil, errors.New("Agent working directory is required")
	}
	model := strings.TrimSpace(runtime.Model)
	if model == "" {
		return nil, errors.New("agent model is required")
	}
	agentArgs, err := acpxAgentArgs(runtime)
	if err != nil {
		return nil, err
	}
	return append([]string{"--cwd", workDir, "--model", model}, agentArgs...), nil
}

func acpxAgentArgs(runtime RuntimeSpec) ([]string, error) {
	if runtime.Protocol == ProtocolStdio {
		command := strings.TrimSpace(runtime.Command)
		if command == "" {
			return nil, errors.New("Agent command override is required")
		}
		return []string{"--agent", command}, nil
	}
	agent := strings.TrimSpace(runtime.ID)
	if agent == "" {
		return nil, errors.New("ACP Runtime id is required")
	}
	return []string{agent}, nil
}

func acpxReasoningEffortConfigKey(runtime RuntimeSpec) (string, error) {
	switch strings.TrimSuffix(strings.TrimSpace(runtime.ID), "-custom") {
	case "codex":
		return acpxCodexReasoningEffortKey, nil
	case "claude", "opencode":
		return acpxGenericReasoningEffortKey, nil
	case "":
		return "", errors.New("ACP Runtime id is required for Agent reasoning effort")
	default:
		return "", fmt.Errorf("unsupported ACP Runtime %q for Agent reasoning effort", runtime.ID)
	}
}

func (runner ACPXRunner) handleStdoutLine(ctx context.Context, req ExecuteRequest, sink runevent.Sink, line []byte, stopReason *string, promptResultParsed *bool) error {
	var message acpxJSONRPCMessage
	if err := json.Unmarshal(line, &message); err != nil {
		return fmt.Errorf("parse acpx stdout JSON-RPC line: %w", err)
	}
	if message.Error != nil {
		return fmt.Errorf("acpx JSON-RPC error %d: %s", message.Error.Code, message.Error.Message)
	}
	if message.Method == acpMethodSessionUpdate {
		update, ok, err := streamUpdateFromSessionUpdatePayload(line)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		event := newAgentRunEvent(req, update, json.RawMessage(append([]byte(nil), line...)), eventClock(runner.Now)())
		if err := sink.Publish(ctx, event); err != nil {
			return fmt.Errorf("publish Run Events: %w", err)
		}
		return nil
	}
	if len(message.Result) == 0 {
		return nil
	}
	var response struct {
		StopReason string `json:"stopReason"`
	}
	if err := json.Unmarshal(message.Result, &response); err != nil {
		return fmt.Errorf("parse acpx session/prompt response: %w", err)
	}
	if response.StopReason != "" {
		*stopReason = response.StopReason
		*promptResultParsed = true
	}
	return nil
}

func acpxTransportAnomaly(exitCode int, stderr string) string {
	message := fmt.Sprintf("acpx exited with exit code %d after parsed session/prompt result", exitCode)
	if tail, truncated := infrastructureStderrTail(stderr); tail != "" {
		message += infrastructureStderrDelimiter
		if truncated {
			message += infrastructureStderrTruncated
		}
		message += tail
	}
	return message
}

func (runner ACPXRunner) mapExitCode(ctx context.Context, req ExecuteRequest, sink runevent.Sink, exitCode int, stderr string, output string) error {
	switch exitCode {
	case 0:
		return nil
	case 1:
		if modelErr, ok := modelNotAdvertisedFromStderr(req.Runtime, stderr, nil); ok {
			return &BatchFailureError{ExitCode: exitCode, Reason: acpxExitReasonAgentProtocol, Stderr: stderr, Err: modelErr}
		}
		return &BatchFailureError{ExitCode: exitCode, Reason: acpxExitReasonAgentProtocol, Stderr: stderr}
	case 3:
		return &BatchFailureError{ExitCode: exitCode, Reason: acpxExitReasonTimeout, Stderr: stderr}
	case 5:
		update := StreamUpdate{Kind: StreamUpdateStatus, Status: acpxPermissionDeniedStatus}
		event := newAgentRunEvent(req, update, marshalStatusPayload(acpxPermissionDeniedStatus), eventClock(runner.Now)())
		if err := sink.Publish(ctx, event); err != nil {
			return fmt.Errorf("publish acpx permission-denied Run Event: %w", err)
		}
		return &BatchFailureError{ExitCode: exitCode, Reason: acpxExitReasonPermissionsDenied, Stderr: stderr}
	case 2:
		return &InfrastructureError{ExitCode: exitCode, Reason: acpxExitReasonUsage, Stderr: stderr}
	case 4:
		return &InfrastructureError{ExitCode: exitCode, Reason: acpxExitReasonMissingSession, Stderr: stderr}
	case 130:
		return StopError{LogPath: req.LogPath, Output: output, Err: context.Canceled}
	default:
		return &InfrastructureError{ExitCode: exitCode, Reason: "unexpected acpx exit code", Stderr: stderr}
	}
}

func (runner ACPXRunner) cancelPrompt(ctx context.Context, req ACPXPromptRequest, waitCh <-chan error, codexEnv []string) bool {
	grace := stopGrace(req.StopGrace)
	cancelCtx, cancel := context.WithTimeout(ctx, grace)
	defer cancel()
	if args, err := acpxCancelArgs(req.Runtime, req.Session, req.GitRoot); err == nil {
		if err := runner.runACPXCommandWithEnv(cancelCtx, args, codexEnv); err != nil {
			runner.warningf("cancel acpx Agent Session %q: %v", req.Session, err)
		}
	} else {
		runner.warningf("build acpx cancel command for Agent Session %q: %v", req.Session, err)
	}
	clock := runner.cancellationClock()
	timer := clock.NewTimer(grace)
	defer timer.Stop()
	select {
	case <-waitCh:
		return false
	case <-timer.C():
		closeCtx, closeCancel := context.WithTimeout(ctx, grace)
		defer closeCancel()
		if args, err := acpxCloseArgs(req.Runtime, req.Session, req.GitRoot); err == nil {
			if err := runner.runACPXCommandWithEnv(closeCtx, args, codexEnv); err != nil {
				runner.warningf("close acpx Agent Session %q: %v", req.Session, err)
			}
		} else {
			runner.warningf("build acpx close command for Agent Session %q: %v", req.Session, err)
		}
		closeTimer := clock.NewTimer(grace)
		defer closeTimer.Stop()
		select {
		case <-waitCh:
		case <-closeTimer.C():
		}
		return true
	}
}

func waitForACPXStream(streamCh <-chan acpxStreamResult, grace time.Duration) acpxStreamResult {
	timer := time.NewTimer(grace)
	defer timer.Stop()
	select {
	case stream := <-streamCh:
		return stream
	case <-timer.C:
		return acpxStreamResult{}
	}
}

func (runner ACPXRunner) publishStatus(ctx context.Context, req ExecuteRequest, sink runevent.Sink, status string) error {
	if sink == nil {
		sink = runevent.Discard
	}
	update := StreamUpdate{Kind: StreamUpdateStatus, Status: status}
	event := newAgentRunEvent(req, update, marshalStatusPayload(status), eventClock(runner.Now)())
	if err := sink.Publish(ctx, event); err != nil {
		return fmt.Errorf("publish acpx Agent status %q: %w", status, err)
	}
	return nil
}

func (runner ACPXRunner) warningf(format string, args ...any) {
	if runner.warnf != nil {
		runner.warnf(format, args...)
		return
	}
	log.Printf(format, args...)
}

func isCodexSandboxUnavailable(err error) bool {
	text := strings.ToLower(err.Error())
	if !strings.Contains(text, acpxCodexSandboxModeKey) && !strings.Contains(text, acpxCodexFullAccessSandbox) {
		return false
	}
	return strings.Contains(text, "unavailable") ||
		strings.Contains(text, "unknown") ||
		strings.Contains(text, "invalid") ||
		strings.Contains(text, "unsupported") ||
		strings.Contains(text, "does not advertise config option")
}

func commandExitCode(err error) (int, bool) {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), true
	}
	return 0, false
}
