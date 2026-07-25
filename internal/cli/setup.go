package cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"roundfix/internal/agent"
	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
)

const (
	setupNodeMinimumVersion = "22.13.0"
)

var setupDeps = defaultSetupDependencies()

type setupDependencies struct {
	loadConfig    func(roundconfig.LoadOptions) (roundconfig.Loaded, error)
	nodeVersion   func(context.Context) (string, error)
	acpxVersion   func(context.Context) (string, error)
	installACPX   func(context.Context) error
	checkAdapter  func(context.Context, agent.RuntimeSpec) (agent.AdapterEvidence, error)
	probeAgent    func(context.Context, agent.ProbeRequest) error
	profileRunner agent.Runner
	lookPath      func(string) (string, error)
	exists        func(string) (bool, error)
	readFile      func(string) ([]byte, error)
	writeFile     func(string, []byte) error
	mkdirAll      func(string) error
	confirm       func(context.Context, io.Writer, string) (bool, error)
}

type setupRequest struct {
	yes     bool
	noInput bool
}

type setupRunner struct {
	req       setupRequest
	deps      setupDependencies
	health    HealthChecker
	loaded    roundconfig.Loaded
	stdout    io.Writer
	stderr    io.Writer
	failed    bool
	acpxReady bool
}

type setupFileProposal struct {
	label   string
	path    string
	before  []byte
	after   []byte
	existed bool
	changed bool
}

type setupProposal struct {
	acpx             setupFileProposal
	user             setupFileProposal
	project          *setupFileProposal
	config           roundconfig.Config
	commandOverrides map[string]string
	adapterMigration bool
}

type acpxAgentOverride struct {
	Agent   string
	Command string
	Args    []string
}

func runSetupCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("setup"))
		return exitOK
	}
	req, err := parseSetupCommand(args)
	if err != nil {
		printSetupFailure(err, stderr)
		return exitPreflight
	}
	loaded, err := setupDeps.loadConfig(roundconfig.LoadOptions{Stderr: stderr})
	if err != nil {
		runner := setupRunner{req: req, deps: setupDeps, stdout: stdout, stderr: stderr}
		runner.report("config", "failed", err.Error())
		return exitRunFailed
	}
	runner := setupRunner{
		req:    req,
		deps:   setupDeps,
		health: setupDeps.healthChecker(),
		loaded: loaded,
		stdout: stdout,
		stderr: stderr,
	}
	runner.checkNode(ctx)
	runner.checkACPX(ctx)
	if !runner.acpxReady {
		runner.reportHealthResult(CheckResult{Name: HealthCheckAdapter, Status: CheckStatusSkipped, Detail: "acpx does not meet the minimum supported version"})
		runner.report("profile readiness", "skipped", "acpx does not meet the minimum supported version")
		if runner.failed {
			return exitRunFailed
		}
		return exitOK
	}
	proposal, ok := runner.buildProposal(ctx)
	if !ok {
		return exitRunFailed
	}
	if !runner.proveProposal(ctx, &proposal) {
		return exitRunFailed
	}
	runner.persistProposal(ctx, proposal)
	if runner.failed {
		return exitRunFailed
	}
	return exitOK
}

func (runner *setupRunner) buildProposal(ctx context.Context) (setupProposal, bool) {
	runtime, err := runtimeForConfiguredAgent(runner.loaded.Config)
	if err != nil {
		runner.report("adapter", "failed", err.Error())
		return setupProposal{}, false
	}
	proposal := setupProposal{commandOverrides: map[string]string{}}
	evidence, adapterErr := runner.deps.checkAdapter(ctx, runtime)
	if adapterErr != nil {
		if strings.TrimSuffix(strings.TrimSpace(runtime.ID), "-custom") != "codex" || !staleCodexAdapter(adapterErr) {
			runner.reportHealthResult(setupAdapterFailureResult(adapterErr))
			return setupProposal{}, false
		}
		acpxProposal, ok := runner.readFileProposal("acpx agents override", filepath.Join(runner.loaded.HomeDir, ".acpx", "config.json"), []byte("{}\n"))
		if !ok {
			return setupProposal{}, false
		}
		hasCodex, err := acpxConfigHasAgent(acpxProposal.after, "codex")
		if err != nil {
			runner.report("acpx agents override", "failed", fmt.Sprintf("parse %s: %v", acpxProposal.path, err))
			return setupProposal{}, false
		}
		if !hasCodex {
			runner.reportHealthResult(setupAdapterFailureResult(adapterErr))
			return setupProposal{}, false
		}

		proposal.adapterMigration = true
		proposal.commandOverrides["codex"] = agent.CodexAdapterCommand()
		proposedRuntime := runtime
		proposedRuntime.Protocol = agent.ProtocolStdio
		proposedRuntime.Command = agent.CodexAdapterCommand()
		evidence, adapterErr = runner.deps.checkAdapter(ctx, proposedRuntime)
		if adapterErr != nil {
			runner.reportHealthResult(setupAdapterFailureResult(adapterErr))
			return setupProposal{}, false
		}
		proposal.acpx = acpxProposal
		runner.report("adapter", "migration proposed", adapterEvidenceDetail(evidence))
	} else {
		runner.reportHealthResult(CheckResult{Name: HealthCheckAdapter, Status: CheckStatusOK, Detail: adapterEvidenceDetail(evidence)})
	}

	if proposal.acpx.path == "" {
		var ok bool
		proposal.acpx, ok = runner.readFileProposal("acpx agents override", filepath.Join(runner.loaded.HomeDir, ".acpx", "config.json"), []byte("{}\n"))
		if !ok {
			return setupProposal{}, false
		}
	}
	overrides := runner.localACPXAgentOverrides()
	if proposal.adapterMigration {
		overrides = append(overrides, officialCodexACPXOverride())
	}
	updated, changed, err := mergeACPXAgentOverrides(proposal.acpx.after, overrides)
	if err != nil {
		runner.report("acpx agents override", "failed", fmt.Sprintf("merge %s: %v", proposal.acpx.path, err))
		return setupProposal{}, false
	}
	proposal.acpx.after = updated
	proposal.acpx.changed = changed
	for _, override := range overrides {
		if changed {
			proposal.commandOverrides[override.Agent] = acpxOverrideCommand(override)
		}
	}

	var ok bool
	proposal.user, ok = runner.readFileProposal("User Config", runner.loaded.UserConfigPath, []byte(roundconfig.DefaultConfigYAML()))
	if !ok {
		return setupProposal{}, false
	}
	if strings.TrimSpace(runner.loaded.GitRoot) != "" {
		project, ok := runner.readFileProposal("Project Config", runner.loaded.ProjectConfigPath, []byte(roundconfig.DefaultConfigYAML()))
		if !ok {
			return setupProposal{}, false
		}
		proposal.project = &project
	}
	var projectContent []byte
	if proposal.project != nil {
		projectContent = proposal.project.after
	}
	proposal.config, err = roundconfig.ResolveConfigProposal(proposal.user.after, projectContent)
	if err != nil {
		runner.report("config proposal", "failed", err.Error())
		return setupProposal{}, false
	}
	return proposal, true
}

func (runner *setupRunner) readFileProposal(label string, path string, generated []byte) (setupFileProposal, bool) {
	exists, err := runner.deps.exists(path)
	if err != nil {
		runner.report(label, "failed", fmt.Sprintf("inspect %s: %v", path, err))
		return setupFileProposal{}, false
	}
	proposal := setupFileProposal{label: label, path: path, existed: exists}
	if exists {
		content, err := runner.deps.readFile(path)
		if err != nil {
			runner.report(label, "failed", fmt.Sprintf("read %s: %v", path, err))
			return setupFileProposal{}, false
		}
		proposal.before = append([]byte(nil), content...)
		proposal.after = append([]byte(nil), content...)
		return proposal, true
	}
	proposal.after = append([]byte(nil), generated...)
	proposal.changed = true
	return proposal, true
}

func (runner *setupRunner) proveProposal(ctx context.Context, proposal *setupProposal) bool {
	proofRunner := runner.deps.profileRunner
	if proofRunner == nil {
		proofRunner = newEngineCollaborators().runner
	}
	if _, ok := proofRunner.(agent.SelectionProver); !ok {
		runner.report("profile readiness", "failed", "exact Agent Selection proof is unavailable")
		return false
	}
	workDir := strings.TrimSpace(runner.loaded.GitRoot)
	if workDir == "" {
		var err error
		workDir, err = os.Getwd()
		if err != nil {
			runner.report("profile readiness", "failed", fmt.Sprintf("resolve Setup proof working directory: %v", err))
			return false
		}
	}
	result := proveProfileSelectionsWithOptions(
		ctx,
		proposal.config,
		roundconfig.RequiredWorkCategories(),
		workDir,
		proofRunner,
		profileProofOptions{CommandOverrides: proposal.commandOverrides},
	)
	if result.Err != nil {
		runner.report("profile readiness", "failed", result.Err.Error())
		return false
	}
	runner.report("profile readiness", "passed", fmt.Sprintf("%d distinct Agent Selections", len(result.Proofs)))
	return true
}

func (runner *setupRunner) persistProposal(ctx context.Context, proposal setupProposal) {
	files := []*setupFileProposal{&proposal.acpx, &proposal.user}
	if proposal.project != nil {
		files = append(files, proposal.project)
	}
	if runner.req.noInput {
		for _, file := range files {
			runner.reportFileProposal(*file, "skipped")
		}
		if proposal.project == nil {
			runner.report("Project Config", "skipped", "outside a git repository")
		}
		return
	}
	if !runner.authorizeProposal(ctx, proposal, files) {
		return
	}
	for _, file := range files {
		if !file.changed {
			runner.report(file.label, "ok", file.path)
			continue
		}
		if err := runner.deps.mkdirAll(filepath.Dir(file.path)); err != nil {
			runner.report(file.label, "failed", fmt.Sprintf("create %s: %v", filepath.Dir(file.path), err))
			return
		}
		if err := runner.deps.writeFile(file.path, file.after); err != nil {
			runner.report(file.label, "failed", fmt.Sprintf("write %s: %v", file.path, err))
			return
		}
		runner.report(file.label, "installed", file.path)
		if file.label == "acpx agents override" {
			printSetupDiff(runner.stdout, file.label, file.path, file.before, file.after)
		}
	}
	if proposal.project == nil {
		runner.report("Project Config", "skipped", "outside a git repository")
	}
}

func (runner *setupRunner) authorizeProposal(ctx context.Context, proposal setupProposal, files []*setupFileProposal) bool {
	if runner.req.yes {
		return true
	}
	for _, file := range files {
		if !file.changed {
			continue
		}
		prompt := "Create " + file.label + " at " + file.path + "?"
		if file.label == "acpx agents override" {
			prompt = "Write acpx local adapter overrides to " + file.path + "?"
			if proposal.adapterMigration {
				prompt = "Migrate stale Codex adapter override in " + file.path + " to " + agent.CodexAdapterCommand() + "?"
			}
		}
		accepted, err := runner.deps.confirm(ctx, runner.stderr, prompt)
		if err != nil {
			runner.report(file.label, "failed", err.Error())
			return false
		}
		if !accepted {
			runner.report(file.label, "offered: declined", "would update "+file.path)
			return false
		}
	}
	return true
}

func (runner *setupRunner) reportFileProposal(file setupFileProposal, changedStatus string) {
	if file.changed {
		runner.report(file.label, changedStatus, "would update "+file.path)
		return
	}
	runner.report(file.label, "ok", file.path)
}

func officialCodexACPXOverride() acpxAgentOverride {
	return acpxAgentOverride{
		Agent:   "codex",
		Command: "npx",
		Args:    []string{"-y", agent.CodexAdapterPackage + "@" + agent.PinnedCodexAdapterVersion},
	}
}

func acpxOverrideCommand(override acpxAgentOverride) string {
	parts := append([]string{override.Command}, override.Args...)
	return strings.Join(parts, " ")
}

func acpxConfigHasAgent(content []byte, name string) (bool, error) {
	var config struct {
		Agents map[string]json.RawMessage `json:"agents"`
	}
	if err := json.Unmarshal(normalizeJSONContent(content), &config); err != nil {
		return false, err
	}
	_, ok := config.Agents[name]
	return ok, nil
}

func staleCodexAdapter(err error) bool {
	var lineage *agent.AdapterLineageError
	if errors.As(err, &lineage) {
		return true
	}
	var version *agent.AdapterVersionError
	return errors.As(err, &version)
}

func setupAdapterFailureResult(err error) CheckResult {
	result := CheckResult{Name: HealthCheckAdapter, Status: CheckStatusFailed, Detail: err.Error()}
	var installer interface{ InstallCommand() string }
	if errors.As(err, &installer) {
		result.NextAction = installer.InstallCommand()
	}
	return result
}

func parseSetupCommand(args []string) (setupRequest, error) {
	req := setupRequest{}
	fs := flag.NewFlagSet("setup", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.BoolVar(&req.yes, "yes", false, "Accept every offered setup change")
	fs.BoolVar(&req.noInput, "no-input", false, "Skip offered setup changes instead of prompting")
	if err := fs.Parse(args); err != nil {
		return req, validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return req, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	if req.yes && req.noInput {
		return req, validationError{message: "--yes cannot be used with --no-input"}
	}
	return req, nil
}

func (runner *setupRunner) checkNode(ctx context.Context) {
	runner.reportHealthResult(runner.health.Node(ctx))
}

func (runner *setupRunner) checkACPX(ctx context.Context) {
	result := runner.health.ACPX(ctx)
	if result.Status == CheckStatusOK {
		runner.acpxReady = true
		runner.reportHealthResult(result)
		return
	}

	detail := result.Detail
	installCommand := result.NextAction
	if installCommand == "" {
		installCommand = setupACPXInstallCommand()
	}
	if runner.req.noInput {
		runner.report("acpx", "skipped", detail)
		return
	}
	if !runner.req.yes {
		accepted, confirmErr := runner.deps.confirm(ctx, runner.stderr, "Install the minimum supported acpx with "+installCommand+"?")
		if confirmErr != nil {
			runner.report("acpx", "failed", confirmErr.Error())
			return
		}
		if !accepted {
			runner.report("acpx", "offered: declined", detail)
			return
		}
	}
	if err := runner.deps.installACPX(ctx); err != nil {
		runner.report("acpx", "failed", fmt.Sprintf("%s failed: %v", setupACPXInstallCommand(), err))
		return
	}
	runner.acpxReady = true
	runner.report("acpx", "installed", detail)
}

func (runner *setupRunner) localACPXAgentOverrides() []acpxAgentOverride {
	candidates := []acpxAgentOverride{
		{Agent: "claude", Command: "claude-agent-acp"},
		{Agent: "opencode", Command: "opencode", Args: []string{"acp"}},
	}
	overrides := make([]acpxAgentOverride, 0, len(candidates))
	for _, candidate := range candidates {
		if _, err := runner.deps.lookPath(candidate.Command); err == nil {
			overrides = append(overrides, candidate)
		}
	}
	return overrides
}

func (runner *setupRunner) report(label string, status string, detail string) {
	if status == "failed" {
		runner.failed = true
	}
	if strings.TrimSpace(detail) == "" {
		fmt.Fprintf(runner.stdout, "%s: %s\n", label, status)
		return
	}
	fmt.Fprintf(runner.stdout, "%s: %s (%s)\n", label, status, detail)
}

func (runner *setupRunner) reportHealthResult(result CheckResult) {
	runner.report(setupHealthCheckLabel(result.Name), string(result.Status), result.Detail)
}

func setupHealthCheckLabel(name string) string {
	if name == HealthCheckAgent {
		return "agent probe"
	}
	return name
}

func (deps setupDependencies) healthChecker() HealthChecker {
	return newHealthChecker(healthCheckDependencies{
		nodeVersion:  deps.nodeVersion,
		acpxVersion:  deps.acpxVersion,
		checkAdapter: deps.checkAdapter,
		probeAgent:   deps.probeAgent,
	})
}

func defaultSetupDependencies() setupDependencies {
	return setupDependencies{
		loadConfig:    roundconfig.Load,
		nodeVersion:   defaultSetupNodeVersion,
		acpxVersion:   defaultSetupACPXVersion,
		installACPX:   defaultSetupInstallACPX,
		checkAdapter:  agent.CheckAdapter,
		profileRunner: newEngineCollaborators().runner,
		probeAgent: func(ctx context.Context, req agent.ProbeRequest) error {
			return newEngineCollaborators().runner.Probe(ctx, req)
		},
		lookPath: exec.LookPath,
		exists: func(path string) (bool, error) {
			_, err := os.Stat(path)
			if errors.Is(err, os.ErrNotExist) {
				return false, nil
			}
			return err == nil, err
		},
		readFile:  os.ReadFile,
		writeFile: defaultSetupWriteFile,
		mkdirAll: func(path string) error {
			return os.MkdirAll(path, 0o755)
		},
		confirm: defaultSetupConfirm,
	}
}

func defaultSetupWriteFile(path string, content []byte) (returnErr error) {
	temp, err := os.CreateTemp(filepath.Dir(path), ".roundfix-setup-*")
	if err != nil {
		return fmt.Errorf("create temporary setup file for %q: %w", path, err)
	}
	tempPath := temp.Name()
	committed := false
	defer func() {
		if committed {
			return
		}
		if cleanupErr := os.Remove(tempPath); cleanupErr != nil && !errors.Is(cleanupErr, os.ErrNotExist) {
			returnErr = errors.Join(returnErr, fmt.Errorf("remove temporary setup file %q: %w", tempPath, cleanupErr))
		}
	}()
	if err := temp.Chmod(0o644); err != nil {
		_ = temp.Close()
		return fmt.Errorf("set temporary setup file permissions for %q: %w", path, err)
	}
	if _, err := temp.Write(content); err != nil {
		_ = temp.Close()
		return fmt.Errorf("write temporary setup file for %q: %w", path, err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temporary setup file for %q: %w", path, err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace setup file %q atomically: %w", path, err)
	}
	committed = true
	return nil
}

func defaultSetupNodeVersion(ctx context.Context) (string, error) {
	if _, err := exec.LookPath("node"); err != nil {
		return "", err
	}
	cmd := exec.CommandContext(ctx, "node", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", commandOutputError(err, output)
	}
	return strings.TrimSpace(string(output)), nil
}

func defaultSetupACPXVersion(ctx context.Context) (string, error) {
	if _, err := exec.LookPath("acpx"); err != nil {
		return "", err
	}
	cmd := exec.CommandContext(ctx, "acpx", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", commandOutputError(err, output)
	}
	return strings.TrimSpace(string(output)), nil
}

func defaultSetupInstallACPX(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "npm", "install", "-g", "acpx@"+agent.MinimumACPXVersion)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return commandOutputError(err, output)
	}
	return nil
}

func defaultSetupConfirm(ctx context.Context, stderr io.Writer, prompt string) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if _, err := fmt.Fprintf(stderr, "%s [y/N]: ", prompt); err != nil {
		return false, fmt.Errorf("write setup prompt: %w", err)
	}
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, fmt.Errorf("read setup prompt: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	return normalizeSetupYesNo(line)
}

func normalizeSetupYesNo(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "y", "yes":
		return true, nil
	case "", "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("unsupported setup response %q; enter yes or no", strings.TrimSpace(value))
	}
}

func commandOutputError(err error, output []byte) error {
	detail := strings.TrimSpace(string(output))
	if detail == "" {
		return err
	}
	return fmt.Errorf("%w: %s", err, detail)
}

func setupACPXInstallCommand() string {
	return "npm install -g acpx@" + agent.MinimumACPXVersion
}

func compareVersions(found string, minimum string) int {
	foundParts, ok := parseVersionParts(found)
	if !ok {
		return -1
	}
	minimumParts, ok := parseVersionParts(minimum)
	if !ok {
		return 0
	}
	for i := range foundParts {
		if foundParts[i] > minimumParts[i] {
			return 1
		}
		if foundParts[i] < minimumParts[i] {
			return -1
		}
	}
	return 0
}

func parseVersionParts(version string) ([3]int, bool) {
	var parts [3]int
	value := strings.TrimSpace(version)
	value = strings.TrimPrefix(value, "v")
	if fields := strings.Fields(value); len(fields) > 0 {
		value = fields[0]
	}
	rawParts := strings.Split(value, ".")
	if len(rawParts) < 2 {
		return parts, false
	}
	for i := range parts {
		if i >= len(rawParts) {
			continue
		}
		number, ok := leadingInt(rawParts[i])
		if !ok {
			return parts, false
		}
		parts[i] = number
	}
	return parts, true
}

func leadingInt(value string) (int, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	end := 0
	for end < len(value) && value[end] >= '0' && value[end] <= '9' {
		end++
	}
	if end == 0 {
		return 0, false
	}
	number, err := strconv.Atoi(value[:end])
	return number, err == nil
}

func mergeACPXAgentOverrides(content []byte, overrides []acpxAgentOverride) ([]byte, bool, error) {
	if len(overrides) == 0 {
		return content, false, nil
	}
	content = normalizeJSONContent(content)
	if err := json.Unmarshal(content, new(any)); err != nil {
		return nil, false, err
	}
	agentsStart, agentsEnd, found, err := findObjectPropertyValue(content, skipSpaces(content, 0), "agents")
	if err != nil {
		return nil, false, err
	}
	if !found {
		members := "  \"agents\": " + formatACPXAgentsObject(overrides, "  ")
		updated, err := insertObjectMembers(content, skipSpaces(content, 0), members)
		if err != nil {
			return nil, false, err
		}
		if err := json.Unmarshal(updated, new(any)); err != nil {
			return nil, false, err
		}
		return updated, true, nil
	}
	if agentsStart >= len(content) || content[agentsStart] != '{' {
		value := []byte(formatACPXAgentsObject(overrides, "  "))
		updated := replaceBytes(content, agentsStart, agentsEnd, value)
		if err := json.Unmarshal(updated, new(any)); err != nil {
			return nil, false, err
		}
		return updated, true, nil
	}

	updated := content
	changed := false
	appendOverrides := []acpxAgentOverride{}
	for _, override := range overrides {
		currentAgentsStart, _, _, err := findObjectPropertyValue(updated, skipSpaces(updated, 0), "agents")
		if err != nil {
			return nil, false, err
		}
		valueStart, valueEnd, agentFound, err := findObjectPropertyValue(updated, currentAgentsStart, override.Agent)
		if err != nil {
			return nil, false, err
		}
		if !agentFound {
			appendOverrides = append(appendOverrides, override)
			continue
		}
		var existing struct {
			Command string   `json:"command"`
			Args    []string `json:"args"`
		}
		if err := json.Unmarshal(updated[valueStart:valueEnd], &existing); err != nil {
			return nil, false, err
		}
		if existing.Command == override.Command && stringSlicesEqual(existing.Args, override.Args) {
			continue
		}
		updated = replaceBytes(updated, valueStart, valueEnd, []byte(formatACPXAgentValue(override, "    ")))
		changed = true
	}
	if len(appendOverrides) > 0 {
		currentAgentsStart, _, _, err := findObjectPropertyValue(updated, skipSpaces(updated, 0), "agents")
		if err != nil {
			return nil, false, err
		}
		updated, err = insertObjectMembers(updated, currentAgentsStart, formatACPXAgentMembers(appendOverrides, "    "))
		if err != nil {
			return nil, false, err
		}
		changed = true
	}
	if changed {
		if err := json.Unmarshal(updated, new(any)); err != nil {
			return nil, false, err
		}
	}
	return updated, changed, nil
}

func normalizeJSONContent(content []byte) []byte {
	if strings.TrimSpace(string(content)) == "" {
		return []byte("{}\n")
	}
	return content
}

func formatACPXAgentsObject(overrides []acpxAgentOverride, indent string) string {
	return "{\n" + formatACPXAgentMembers(overrides, indent+"  ") + "\n" + indent + "}"
}

func formatACPXAgentMembers(overrides []acpxAgentOverride, indent string) string {
	var builder strings.Builder
	for i, override := range overrides {
		if i > 0 {
			builder.WriteString(",\n")
		}
		builder.WriteString(indent)
		builder.WriteString(strconv.Quote(override.Agent))
		builder.WriteString(": ")
		builder.WriteString(formatACPXAgentValue(override, indent))
	}
	return builder.String()
}

func formatACPXAgentValue(override acpxAgentOverride, indent string) string {
	var builder strings.Builder
	builder.WriteString("{\n")
	builder.WriteString(indent)
	builder.WriteString("  \"command\": ")
	builder.WriteString(strconv.Quote(override.Command))
	if len(override.Args) > 0 {
		args, _ := json.Marshal(override.Args)
		builder.WriteString(",\n")
		builder.WriteString(indent)
		builder.WriteString("  \"args\": ")
		builder.Write(args)
	}
	builder.WriteString("\n")
	builder.WriteString(indent)
	builder.WriteString("}")
	return builder.String()
}

func insertObjectMembers(content []byte, objectStart int, members string) ([]byte, error) {
	if objectStart >= len(content) || content[objectStart] != '{' {
		return nil, errors.New("target JSON value is not an object")
	}
	objectEnd, err := skipJSONValue(content, objectStart)
	if err != nil {
		return nil, err
	}
	closeIndex := objectEnd - 1
	insertIndex := closeIndex
	for insertIndex > objectStart+1 && isJSONSpace(content[insertIndex-1]) {
		insertIndex--
	}
	hasMembers := strings.TrimSpace(string(content[objectStart+1:closeIndex])) != ""
	addition := "\n" + members + "\n"
	if hasMembers {
		addition = ",\n" + members
	}
	return insertBytes(content, insertIndex, []byte(addition)), nil
}

func findObjectPropertyValue(content []byte, objectStart int, key string) (int, int, bool, error) {
	objectStart = skipSpaces(content, objectStart)
	if objectStart >= len(content) || content[objectStart] != '{' {
		return 0, 0, false, errors.New("target JSON value is not an object")
	}
	i := objectStart + 1
	for {
		i = skipSpaces(content, i)
		if i >= len(content) {
			return 0, 0, false, errors.New("unterminated JSON object")
		}
		if content[i] == '}' {
			return 0, 0, false, nil
		}
		name, next, err := readJSONString(content, i)
		if err != nil {
			return 0, 0, false, err
		}
		i = skipSpaces(content, next)
		if i >= len(content) || content[i] != ':' {
			return 0, 0, false, errors.New("expected JSON object colon")
		}
		valueStart := skipSpaces(content, i+1)
		valueEnd, err := skipJSONValue(content, valueStart)
		if err != nil {
			return 0, 0, false, err
		}
		if name == key {
			return valueStart, valueEnd, true, nil
		}
		i = skipSpaces(content, valueEnd)
		if i >= len(content) {
			return 0, 0, false, errors.New("unterminated JSON object")
		}
		if content[i] == ',' {
			i++
			continue
		}
		if content[i] == '}' {
			return 0, 0, false, nil
		}
		return 0, 0, false, errors.New("expected JSON object comma or close")
	}
}

func readJSONString(content []byte, start int) (string, int, error) {
	end, err := skipJSONString(content, start)
	if err != nil {
		return "", 0, err
	}
	var value string
	if err := json.Unmarshal(content[start:end], &value); err != nil {
		return "", 0, err
	}
	return value, end, nil
}

func skipJSONValue(content []byte, start int) (int, error) {
	start = skipSpaces(content, start)
	if start >= len(content) {
		return 0, errors.New("missing JSON value")
	}
	switch content[start] {
	case '{':
		i := start + 1
		for {
			i = skipSpaces(content, i)
			if i >= len(content) {
				return 0, errors.New("unterminated JSON object")
			}
			if content[i] == '}' {
				return i + 1, nil
			}
			_, next, err := readJSONString(content, i)
			if err != nil {
				return 0, err
			}
			i = skipSpaces(content, next)
			if i >= len(content) || content[i] != ':' {
				return 0, errors.New("expected JSON object colon")
			}
			next, err = skipJSONValue(content, i+1)
			if err != nil {
				return 0, err
			}
			i = skipSpaces(content, next)
			if i >= len(content) {
				return 0, errors.New("unterminated JSON object")
			}
			if content[i] == ',' {
				i++
				continue
			}
			if content[i] == '}' {
				return i + 1, nil
			}
			return 0, errors.New("expected JSON object comma or close")
		}
	case '[':
		i := start + 1
		for {
			i = skipSpaces(content, i)
			if i >= len(content) {
				return 0, errors.New("unterminated JSON array")
			}
			if content[i] == ']' {
				return i + 1, nil
			}
			next, err := skipJSONValue(content, i)
			if err != nil {
				return 0, err
			}
			i = skipSpaces(content, next)
			if i >= len(content) {
				return 0, errors.New("unterminated JSON array")
			}
			if content[i] == ',' {
				i++
				continue
			}
			if content[i] == ']' {
				return i + 1, nil
			}
			return 0, errors.New("expected JSON array comma or close")
		}
	case '"':
		return skipJSONString(content, start)
	default:
		i := start
		for i < len(content) && !isJSONDelimiter(content[i]) {
			i++
		}
		if i == start {
			return 0, errors.New("missing JSON scalar")
		}
		var value any
		if err := json.Unmarshal(content[start:i], &value); err != nil {
			return 0, err
		}
		return i, nil
	}
}

func skipJSONString(content []byte, start int) (int, error) {
	if start >= len(content) || content[start] != '"' {
		return 0, errors.New("expected JSON string")
	}
	escaped := false
	for i := start + 1; i < len(content); i++ {
		if escaped {
			escaped = false
			continue
		}
		switch content[i] {
		case '\\':
			escaped = true
		case '"':
			return i + 1, nil
		}
	}
	return 0, errors.New("unterminated JSON string")
}

func skipSpaces(content []byte, start int) int {
	for start < len(content) && isJSONSpace(content[start]) {
		start++
	}
	return start
}

func isJSONSpace(value byte) bool {
	return value == ' ' || value == '\n' || value == '\r' || value == '\t'
}

func isJSONDelimiter(value byte) bool {
	return isJSONSpace(value) || value == ',' || value == '}' || value == ']'
}

func insertBytes(content []byte, index int, addition []byte) []byte {
	result := make([]byte, 0, len(content)+len(addition))
	result = append(result, content[:index]...)
	result = append(result, addition...)
	result = append(result, content[index:]...)
	return result
}

func replaceBytes(content []byte, start int, end int, replacement []byte) []byte {
	result := make([]byte, 0, len(content)-(end-start)+len(replacement))
	result = append(result, content[:start]...)
	result = append(result, replacement...)
	result = append(result, content[end:]...)
	return result
}

func stringSlicesEqual(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func printSetupDiff(stdout io.Writer, label string, path string, before []byte, after []byte) {
	fmt.Fprintf(stdout, "%s diff:\n", label)
	fmt.Fprintf(stdout, "--- %s (before)\n", path)
	fmt.Fprintf(stdout, "+++ %s (after)\n", path)
	for _, line := range splitDiffLines(before) {
		fmt.Fprintf(stdout, "-%s\n", line)
	}
	for _, line := range splitDiffLines(after) {
		fmt.Fprintf(stdout, "+%s\n", line)
	}
}

func splitDiffLines(content []byte) []string {
	content = bytes.TrimRight(content, "\n")
	if len(content) == 0 {
		return []string{""}
	}
	return strings.Split(string(content), "\n")
}

func printSetupFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: setup failed: %v\n", app.Name, err)
	fmt.Fprintf(stderr, "Run '%s setup --help' for usage.\n", app.Name)
}
