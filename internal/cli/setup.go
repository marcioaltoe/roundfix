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
	loadConfig     func(roundconfig.LoadOptions) (roundconfig.Loaded, error)
	nodeVersion    func(context.Context) (string, error)
	acpxVersion    func(context.Context) (string, error)
	installACPX    func(context.Context) error
	probeAgent     func(context.Context, agent.RuntimeSpec) error
	lookPath       func(string) (string, error)
	exists         func(string) (bool, error)
	readFile       func(string) ([]byte, error)
	writeFile      func(string, []byte) error
	mkdirAll       func(string) error
	initACPXConfig func(context.Context, string) error
	initConfig     func(context.Context, roundconfig.InitOptions) (roundconfig.InitResult, error)
	confirm        func(context.Context, io.Writer, string) (bool, error)
}

type setupRequest struct {
	yes     bool
	noInput bool
}

type setupRunner struct {
	req       setupRequest
	deps      setupDependencies
	loaded    roundconfig.Loaded
	stdout    io.Writer
	stderr    io.Writer
	failed    bool
	acpxReady bool
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
	loaded, err := setupDeps.loadConfig(roundconfig.LoadOptions{})
	if err != nil {
		runner := setupRunner{req: req, deps: setupDeps, stdout: stdout, stderr: stderr}
		runner.report("config", "failed", err.Error())
		return exitRunFailed
	}
	runner := setupRunner{
		req:    req,
		deps:   setupDeps,
		loaded: loaded,
		stdout: stdout,
		stderr: stderr,
	}
	runner.checkNode(ctx)
	runner.checkACPX(ctx)
	runner.checkAgentProbe(ctx)
	runner.checkACPXAgentsOverride(ctx)
	runner.checkRoundfixConfig(ctx, roundconfig.InitScopeUser)
	runner.checkRoundfixConfig(ctx, roundconfig.InitScopeProject)
	if runner.failed {
		return exitRunFailed
	}
	return exitOK
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
	version, err := runner.deps.nodeVersion(ctx)
	if err != nil {
		runner.report("node", "failed", fmt.Sprintf("Node.js %s or newer is required: %v", setupNodeMinimumVersion, err))
		return
	}
	if compareVersions(version, setupNodeMinimumVersion) < 0 {
		runner.report("node", "failed", fmt.Sprintf("found %s; Node.js %s or newer is required", strings.TrimSpace(version), setupNodeMinimumVersion))
		return
	}
	runner.report("node", "ok", fmt.Sprintf("%s >= %s", strings.TrimSpace(version), setupNodeMinimumVersion))
}

func (runner *setupRunner) checkACPX(ctx context.Context) {
	version, err := runner.deps.acpxVersion(ctx)
	if err == nil && strings.TrimSpace(version) == agent.PinnedACPXVersion {
		runner.acpxReady = true
		runner.report("acpx", "ok", agent.PinnedACPXVersion)
		return
	}

	detail := fmt.Sprintf("run %s", setupACPXInstallCommand())
	if err == nil {
		detail = fmt.Sprintf("found %s; required %s; run %s", strings.TrimSpace(version), agent.PinnedACPXVersion, setupACPXInstallCommand())
	}
	if runner.req.noInput {
		runner.report("acpx", "skipped", detail)
		return
	}
	if !runner.req.yes {
		accepted, confirmErr := runner.deps.confirm(ctx, runner.stderr, "Install pinned acpx with "+setupACPXInstallCommand()+"?")
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

func (runner *setupRunner) checkAgentProbe(ctx context.Context) {
	if !runner.acpxReady {
		runner.report("agent probe", "skipped", "acpx is not at the pinned version")
		return
	}
	runtime, err := agent.RuntimeFor(agent.RuntimeOptions{
		Agent:            runner.loaded.Config.Defaults.Agent,
		Model:            runner.loaded.Config.Defaults.Model,
		EnableFullAccess: runner.loaded.Config.Defaults.AgentFullAccess,
	})
	if err != nil {
		runner.report("agent probe", "failed", err.Error())
		return
	}
	if err := runner.deps.probeAgent(ctx, runtime); err != nil {
		runner.report("agent probe", "failed", err.Error())
		return
	}
	runner.report("agent probe", "ok", runtime.ID)
}

func (runner *setupRunner) checkACPXAgentsOverride(ctx context.Context) {
	overrides := runner.localACPXAgentOverrides()
	if len(overrides) == 0 {
		runner.report("acpx agents override", "ok", "no local adapter binaries found on PATH")
		return
	}

	configPath := filepath.Join(runner.loaded.HomeDir, ".acpx", "config.json")
	exists, err := runner.deps.exists(configPath)
	if err != nil {
		runner.report("acpx agents override", "failed", fmt.Sprintf("inspect %s: %v", configPath, err))
		return
	}
	content := []byte("{}\n")
	if exists {
		content, err = runner.deps.readFile(configPath)
		if err != nil {
			runner.report("acpx agents override", "failed", fmt.Sprintf("read %s: %v", configPath, err))
			return
		}
	}
	missing, err := missingACPXAgentOverrides(content, overrides)
	if err != nil {
		runner.report("acpx agents override", "failed", fmt.Sprintf("parse %s: %v", configPath, err))
		return
	}
	if len(missing) == 0 {
		runner.report("acpx agents override", "ok", configPath)
		return
	}
	detail := fmt.Sprintf("would update %s", configPath)
	if runner.req.noInput {
		runner.report("acpx agents override", "skipped", detail)
		return
	}
	if !runner.req.yes {
		accepted, confirmErr := runner.deps.confirm(ctx, runner.stderr, "Write acpx local adapter overrides to "+configPath+"?")
		if confirmErr != nil {
			runner.report("acpx agents override", "failed", confirmErr.Error())
			return
		}
		if !accepted {
			runner.report("acpx agents override", "offered: declined", detail)
			return
		}
	}
	if !exists {
		if err := runner.deps.initACPXConfig(ctx, configPath); err != nil {
			runner.report("acpx agents override", "failed", fmt.Sprintf("initialize %s: %v", configPath, err))
			return
		}
		content, err = runner.deps.readFile(configPath)
		if err != nil {
			runner.report("acpx agents override", "failed", fmt.Sprintf("read initialized %s: %v", configPath, err))
			return
		}
	}
	updated, changed, err := mergeACPXAgentOverrides(content, missing)
	if err != nil {
		runner.report("acpx agents override", "failed", fmt.Sprintf("merge %s: %v", configPath, err))
		return
	}
	if !changed {
		runner.report("acpx agents override", "ok", configPath)
		return
	}
	if err := runner.deps.mkdirAll(filepath.Dir(configPath)); err != nil {
		runner.report("acpx agents override", "failed", fmt.Sprintf("create %s: %v", filepath.Dir(configPath), err))
		return
	}
	if err := runner.deps.writeFile(configPath, updated); err != nil {
		runner.report("acpx agents override", "failed", fmt.Sprintf("write %s: %v", configPath, err))
		return
	}
	runner.report("acpx agents override", "installed", configPath)
	printSetupDiff(runner.stdout, "acpx agents override", configPath, content, updated)
}

func (runner *setupRunner) checkRoundfixConfig(ctx context.Context, scope string) {
	label := "User Config"
	path := runner.loaded.UserConfigPath
	if scope == roundconfig.InitScopeProject {
		label = "Project Config"
		path = runner.loaded.ProjectConfigPath
		if strings.TrimSpace(runner.loaded.GitRoot) == "" {
			runner.report(label, "skipped", "outside a git repository")
			return
		}
	}
	exists, err := runner.deps.exists(path)
	if err != nil {
		runner.report(label, "failed", fmt.Sprintf("inspect %s: %v", path, err))
		return
	}
	if exists {
		runner.report(label, "ok", path)
		return
	}
	detail := fmt.Sprintf("would create %s", path)
	if runner.req.noInput {
		runner.report(label, "skipped", detail)
		return
	}
	if !runner.req.yes {
		accepted, confirmErr := runner.deps.confirm(ctx, runner.stderr, "Create "+label+" at "+path+"?")
		if confirmErr != nil {
			runner.report(label, "failed", confirmErr.Error())
			return
		}
		if !accepted {
			runner.report(label, "offered: declined", detail)
			return
		}
	}
	result, err := runner.deps.initConfig(ctx, roundconfig.InitOptions{Scope: scope})
	if err != nil {
		runner.report(label, "failed", err.Error())
		return
	}
	runner.report(label, "installed", result.Path)
}

func (runner *setupRunner) localACPXAgentOverrides() []acpxAgentOverride {
	candidates := []acpxAgentOverride{
		{Agent: "codex", Command: "codex-acp"},
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

func defaultSetupDependencies() setupDependencies {
	return setupDependencies{
		loadConfig:  roundconfig.Load,
		nodeVersion: defaultSetupNodeVersion,
		acpxVersion: defaultSetupACPXVersion,
		installACPX: defaultSetupInstallACPX,
		probeAgent: func(ctx context.Context, runtime agent.RuntimeSpec) error {
			return newEngineCollaborators().runner.Probe(ctx, runtime)
		},
		lookPath: exec.LookPath,
		exists: func(path string) (bool, error) {
			_, err := os.Stat(path)
			if errors.Is(err, os.ErrNotExist) {
				return false, nil
			}
			return err == nil, err
		},
		readFile: os.ReadFile,
		writeFile: func(path string, content []byte) error {
			return os.WriteFile(path, content, 0o644)
		},
		mkdirAll: func(path string) error {
			return os.MkdirAll(path, 0o755)
		},
		initACPXConfig: defaultSetupInitACPXConfig,
		initConfig:     roundconfig.Init,
		confirm:        defaultSetupConfirm,
	}
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
	cmd := exec.CommandContext(ctx, "npm", "install", "-g", "acpx@"+agent.PinnedACPXVersion)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return commandOutputError(err, output)
	}
	return nil
}

func defaultSetupInitACPXConfig(ctx context.Context, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create acpx config directory: %w", err)
	}
	cmd := exec.CommandContext(ctx, "acpx", "config", "init")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return commandOutputError(err, output)
	}
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("acpx config init did not create %s: %w", path, err)
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
	return "npm install -g acpx@" + agent.PinnedACPXVersion
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

func missingACPXAgentOverrides(content []byte, overrides []acpxAgentOverride) ([]acpxAgentOverride, error) {
	content = normalizeJSONContent(content)
	var config struct {
		Agents map[string]struct {
			Command string   `json:"command"`
			Args    []string `json:"args"`
		} `json:"agents"`
	}
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, err
	}
	missing := []acpxAgentOverride{}
	for _, override := range overrides {
		existing, ok := config.Agents[override.Agent]
		if !ok || existing.Command != override.Command || !stringSlicesEqual(existing.Args, override.Args) {
			missing = append(missing, override)
		}
	}
	return missing, nil
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
