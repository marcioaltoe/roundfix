package config

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	userConfigRelPath               = ".roundfix/config.yml"
	projectConfigName               = ".roundfixrc.yml"
	defaultReviewSource             = "coderabbit"
	defaultReviewRequestCommand     = "@coderabbitai review"
	defaultAgent                    = "codex"
	defaultCodexModel               = "gpt-5.5"
	defaultCodexReasoningEffort     = "xhigh"
	defaultClaudeModel              = "opus"
	defaultClaudeReasoningEffort    = ""
	defaultVerification             = "make verify"
	defaultPollInterval             = 30 * time.Second
	defaultReviewTimeout            = 30 * time.Minute
	defaultCheckGracePeriod         = 5 * time.Minute
	defaultQuietPeriod              = 30 * time.Second
	defaultRunDuration              = 2 * time.Hour
	defaultJournalRetention         = 336 * time.Hour
	defaultWorktreeLocation         = "~/.roundfix/worktrees"
	defaultWorktreeConcurrency      = 2
	defaultVerificationConcurrency  = 1
	defaultWorktreeBootstrapTimeout = 10 * time.Minute
	defaultSpecsRoot                = "docs/specs"
)

const (
	InitScopeUser    = "user"
	InitScopeProject = "project"
)

type Config struct {
	Defaults     Defaults
	Runtimes     Runtimes
	Profiles     Profiles
	ReviewSource ReviewSource
	Watch        Watch
	Implement    Implement
	Notify       Notify
	Worktree     Worktree
	Verification Verification
	Budget       Budget
	Resolve      Resolve
	Logs         Logs
	Store        Store
	Specs        Specs
}

type Defaults struct {
	Agent           string
	Model           string
	AgentFullAccess bool
	AutoCommit      bool
	Verification    string
	ArtifactDir     string
}

type RuntimeDefaults struct {
	Model           string
	ReasoningEffort string
}

type Runtimes struct {
	Codex    RuntimeDefaults
	Claude   RuntimeDefaults
	OpenCode RuntimeDefaults
}

func (runtimes Runtimes) DefaultsFor(runtime string) (RuntimeDefaults, bool) {
	switch strings.TrimSpace(runtime) {
	case "codex":
		return runtimes.Codex, true
	case "claude":
		return runtimes.Claude, true
	case "opencode":
		return runtimes.OpenCode, true
	default:
		return RuntimeDefaults{}, false
	}
}

type ReviewSource struct {
	Name            string
	IncludeNitpicks bool
	RequestReview   bool
	RequestCommand  string
}

type Watch struct {
	UntilClean       bool
	MaxRounds        int
	PollInterval     time.Duration
	ReviewTimeout    time.Duration
	CheckGracePeriod time.Duration
	QuietPeriod      time.Duration
	AutoPush         bool
	PushRemote       string
	PushBranch       string
}

type Implement struct {
	AutoPush bool
}

type Notify struct {
	Enabled bool
	Command string
}

type Worktree struct {
	Concurrency      int
	Location         string
	Copy             []string
	Bootstrap        string
	BootstrapTimeout time.Duration
}

type Verification struct {
	Concurrency int
}

type Budget struct {
	Enabled        bool
	MaxRunDuration time.Duration
}

type Resolve struct {
	BatchSize int
}

type Logs struct {
	Agent bool
}

type Store struct {
	JournalRetention time.Duration
}

type Specs struct {
	Root string
}

type Loaded struct {
	Config            Config
	GitRoot           string
	HomeDir           string
	UserConfigPath    string
	ProjectConfigPath string
}

type SpecsRoot struct {
	Path        string
	External    bool
	BuiltInRoot bool
}

type LoadOptions struct {
	HomeDir string
	WorkDir string
	Stderr  io.Writer
}

type InitOptions struct {
	Scope   string
	HomeDir string
	WorkDir string
	Force   bool
}

type InitResult struct {
	Scope       string
	Path        string
	Overwritten bool
}

type durationValue struct {
	value time.Duration
}

func (duration *durationValue) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return errors.New("duration must be a scalar")
	}
	raw := node.Value
	if raw == "0" {
		duration.value = 0
		return nil
	}
	value, err := time.ParseDuration(raw)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", raw, err)
	}
	duration.value = value
	return nil
}

type configOverlay struct {
	Defaults     *defaultsOverlay     `yaml:"defaults"`
	Runtimes     *runtimesOverlay     `yaml:"runtimes"`
	Profiles     *profilesOverlay     `yaml:"profiles"`
	ReviewSource *reviewSourceOverlay `yaml:"review_source"`
	Watch        *watchOverlay        `yaml:"watch"`
	Implement    *implementOverlay    `yaml:"implement"`
	Notify       *notifyOverlay       `yaml:"notify"`
	Worktree     *worktreeOverlay     `yaml:"worktree"`
	Verification *verificationOverlay `yaml:"verification"`
	Budget       *budgetOverlay       `yaml:"budget"`
	Resolve      *resolveOverlay      `yaml:"resolve"`
	Logs         *logsOverlay         `yaml:"logs"`
	Store        *storeOverlay        `yaml:"store"`
	Specs        *specsOverlay        `yaml:"specs"`
}

type defaultsOverlay struct {
	Agent           *string `yaml:"agent"`
	AgentFullAccess *bool   `yaml:"agent_full_access"`
	AutoCommit      *bool   `yaml:"auto_commit"`
	Verification    *string `yaml:"verification"`
	ArtifactDir     *string `yaml:"artifact_dir"`
}

type runtimeDefaultsOverlay struct {
	Model           *string `yaml:"model"`
	ReasoningEffort *string `yaml:"reasoning_effort"`
}

type runtimesOverlay struct {
	Codex    *runtimeDefaultsOverlay `yaml:"codex"`
	Claude   *runtimeDefaultsOverlay `yaml:"claude"`
	OpenCode *runtimeDefaultsOverlay `yaml:"opencode"`
}

type reviewSourceOverlay struct {
	Name            *string             `yaml:"name"`
	IncludeNitpicks *bool               `yaml:"include_nitpicks"`
	RequestReview   *requestReviewValue `yaml:"request_review"`
	RequestCommand  *string             `yaml:"request_command"`
}

type requestReviewValue struct {
	value bool
}

func (value *requestReviewValue) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode || node.Tag != "!!bool" {
		return errors.New("review_source.request_review must be boolean: cannot unmarshal non-boolean value")
	}
	var raw bool
	if err := node.Decode(&raw); err != nil {
		return fmt.Errorf("review_source.request_review must be boolean: %w", err)
	}
	value.value = raw
	return nil
}

type watchOverlay struct {
	UntilClean       *bool          `yaml:"until_clean"`
	MaxRounds        *int           `yaml:"max_rounds"`
	PollInterval     *durationValue `yaml:"poll_interval"`
	ReviewTimeout    *durationValue `yaml:"review_timeout"`
	CheckGracePeriod *durationValue `yaml:"check_grace_period"`
	QuietPeriod      *durationValue `yaml:"quiet_period"`
	AutoPush         *bool          `yaml:"auto_push"`
	PushRemote       *string        `yaml:"push_remote"`
	PushBranch       *string        `yaml:"push_branch"`
}

type implementOverlay struct {
	AutoPush *implementAutoPushValue `yaml:"auto_push"`
}

type notifyOverlay struct {
	Enabled *bool   `yaml:"enabled"`
	Command *string `yaml:"command"`
}

type implementAutoPushValue struct {
	value bool
}

func (value *implementAutoPushValue) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode || node.Tag != "!!bool" {
		return errors.New("implement.auto_push must be boolean")
	}
	var raw bool
	if err := node.Decode(&raw); err != nil {
		return fmt.Errorf("implement.auto_push must be boolean: %w", err)
	}
	value.value = raw
	return nil
}

func (overlay *notifyOverlay) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.MappingNode {
		for index := 0; index < len(node.Content); index += 2 {
			key := node.Content[index].Value
			switch key {
			case "enabled", "command":
			default:
				return fmt.Errorf("notify.%s is not a supported config key", key)
			}
		}
	}
	type rawNotifyOverlay notifyOverlay
	var raw rawNotifyOverlay
	if err := node.Decode(&raw); err != nil {
		return err
	}
	*overlay = notifyOverlay(raw)
	return nil
}

type worktreeOverlay struct {
	Concurrency      *int           `yaml:"concurrency"`
	Location         *string        `yaml:"location"`
	Copy             *[]string      `yaml:"copy"`
	Bootstrap        *string        `yaml:"bootstrap"`
	BootstrapTimeout *durationValue `yaml:"bootstrap_timeout"`
}

type verificationOverlay struct {
	Concurrency *verificationConcurrencyValue `yaml:"concurrency"`
}

type verificationConcurrencyValue struct {
	value int
}

func (value *verificationConcurrencyValue) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode || node.Tag != "!!int" {
		return errors.New("verification.concurrency must be an integer")
	}
	var raw int
	if err := node.Decode(&raw); err != nil {
		return fmt.Errorf("verification.concurrency must be an integer: %w", err)
	}
	value.value = raw
	return nil
}

func (overlay *verificationOverlay) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.New("verification must be a mapping")
	}
	for index := 0; index < len(node.Content); index += 2 {
		key := node.Content[index].Value
		switch key {
		case "concurrency":
		default:
			return fmt.Errorf("verification.%s is not a supported config key", key)
		}
	}
	type rawVerificationOverlay verificationOverlay
	var raw rawVerificationOverlay
	if err := node.Decode(&raw); err != nil {
		return err
	}
	*overlay = verificationOverlay(raw)
	return nil
}

type budgetOverlay struct {
	Enabled        *bool          `yaml:"enabled"`
	MaxRunDuration *durationValue `yaml:"max_run_duration"`
}

type resolveOverlay struct {
	BatchSize *int `yaml:"batch_size"`
}

type logsOverlay struct {
	Agent *bool `yaml:"agent"`
}

type storeOverlay struct {
	JournalRetention *durationValue `yaml:"journal_retention"`
}

type specsOverlay struct {
	Root *string `yaml:"root"`
}

func (overlay *runtimesOverlay) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.MappingNode {
		for index := 0; index < len(node.Content); index += 2 {
			key := node.Content[index].Value
			switch key {
			case "codex", "claude", "opencode":
				if err := validateRuntimeDefaultsOverlay("runtimes."+key, node.Content[index+1]); err != nil {
					return err
				}
			default:
				return fmt.Errorf("runtimes.%s is not a supported config key", key)
			}
		}
	}
	type rawRuntimesOverlay runtimesOverlay
	var raw rawRuntimesOverlay
	if err := node.Decode(&raw); err != nil {
		return err
	}
	*overlay = runtimesOverlay(raw)
	return nil
}

func validateRuntimeDefaultsOverlay(prefix string, node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return nil
	}
	for index := 0; index < len(node.Content); index += 2 {
		key := node.Content[index].Value
		switch key {
		case "model", "reasoning_effort":
		default:
			return fmt.Errorf("%s.%s is not a supported config key", prefix, key)
		}
	}
	return nil
}

func (overlay *resolveOverlay) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.MappingNode {
		for index := 0; index < len(node.Content); index += 2 {
			key := node.Content[index].Value
			switch key {
			case "batch_size":
			default:
				return fmt.Errorf("resolve.%s is not a supported config key", key)
			}
		}
	}
	type rawResolveOverlay resolveOverlay
	var raw rawResolveOverlay
	if err := node.Decode(&raw); err != nil {
		return err
	}
	*overlay = resolveOverlay(raw)
	return nil
}

func (overlay *storeOverlay) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.MappingNode {
		for index := 0; index < len(node.Content); index += 2 {
			key := node.Content[index].Value
			switch key {
			case "journal_retention":
			default:
				return fmt.Errorf("store.%s is not a supported config key", key)
			}
		}
	}
	type rawStoreOverlay storeOverlay
	var raw rawStoreOverlay
	if err := node.Decode(&raw); err != nil {
		return err
	}
	*overlay = storeOverlay(raw)
	return nil
}

func (overlay *specsOverlay) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.MappingNode {
		for index := 0; index < len(node.Content); index += 2 {
			key := node.Content[index].Value
			switch key {
			case "root":
			default:
				return fmt.Errorf("specs.%s is not a supported config key", key)
			}
		}
	}
	type rawSpecsOverlay specsOverlay
	var raw rawSpecsOverlay
	if err := node.Decode(&raw); err != nil {
		return err
	}
	*overlay = specsOverlay(raw)
	return nil
}

func (overlay *worktreeOverlay) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.MappingNode {
		for index := 0; index < len(node.Content); index += 2 {
			key := node.Content[index].Value
			switch key {
			case "concurrency", "location", "copy", "bootstrap", "bootstrap_timeout":
			default:
				return fmt.Errorf("worktree.%s is not a supported config key", key)
			}
		}
	}
	type rawWorktreeOverlay worktreeOverlay
	var raw rawWorktreeOverlay
	if err := node.Decode(&raw); err != nil {
		return err
	}
	*overlay = worktreeOverlay(raw)
	return nil
}

type deprecatedConfigKey struct {
	path        []string
	name        string
	replacement string
}

var deprecatedConfigKeys = []deprecatedConfigKey{
	{
		path:        []string{"defaults", "model"},
		name:        "defaults.model",
		replacement: "profiles.<category>.preferred.model",
	},
	{
		path:        []string{"resolve", "concurrent"},
		name:        "resolve.concurrent",
		replacement: "worktree.concurrency",
	},
}

type deprecatedConfigWarnings struct {
	stderr  io.Writer
	emitted map[string]bool
}

func newDeprecatedConfigWarnings(stderr io.Writer) *deprecatedConfigWarnings {
	if stderr == nil {
		stderr = os.Stderr
	}
	return &deprecatedConfigWarnings{
		stderr:  stderr,
		emitted: map[string]bool{},
	}
}

func (warnings *deprecatedConfigWarnings) warn(key deprecatedConfigKey) {
	if warnings == nil || warnings.emitted[key.name] {
		return
	}
	warnings.emitted[key.name] = true
	fmt.Fprintf(warnings.stderr, "config: %s is deprecated and ignored; use %s\n", key.name, key.replacement)
}

func Builtin() Config {
	return Config{
		Defaults: Defaults{
			Agent:        defaultAgent,
			AutoCommit:   true,
			Verification: defaultVerification,
		},
		Runtimes: Runtimes{
			Codex: RuntimeDefaults{
				Model:           defaultCodexModel,
				ReasoningEffort: defaultCodexReasoningEffort,
			},
			Claude: RuntimeDefaults{
				Model:           defaultClaudeModel,
				ReasoningEffort: defaultClaudeReasoningEffort,
			},
		},
		Profiles: builtinProfiles(),
		ReviewSource: ReviewSource{
			Name:            defaultReviewSource,
			IncludeNitpicks: false,
			RequestReview:   false,
			RequestCommand:  defaultReviewRequestCommand,
		},
		Watch: Watch{
			UntilClean:       true,
			MaxRounds:        6,
			PollInterval:     defaultPollInterval,
			ReviewTimeout:    defaultReviewTimeout,
			CheckGracePeriod: defaultCheckGracePeriod,
			QuietPeriod:      defaultQuietPeriod,
			AutoPush:         true,
		},
		Implement: Implement{
			AutoPush: false,
		},
		Notify: Notify{
			Enabled: true,
		},
		Worktree: Worktree{
			Concurrency:      defaultWorktreeConcurrency,
			Location:         defaultWorktreeLocation,
			BootstrapTimeout: defaultWorktreeBootstrapTimeout,
		},
		Verification: Verification{
			Concurrency: defaultVerificationConcurrency,
		},
		Budget: Budget{
			Enabled:        true,
			MaxRunDuration: defaultRunDuration,
		},
		Resolve: Resolve{
			BatchSize: 3,
		},
		Logs: Logs{
			Agent: false,
		},
		Store: Store{
			JournalRetention: defaultJournalRetention,
		},
		Specs: Specs{
			Root: defaultSpecsRoot,
		},
	}
}

func Load(opts LoadOptions) (Loaded, error) {
	homeDir, err := resolveHomeDir(opts.HomeDir)
	if err != nil {
		return Loaded{}, err
	}
	workDir, err := resolveWorkDir(opts.WorkDir)
	if err != nil {
		return Loaded{}, err
	}

	loaded := Loaded{
		Config:         Builtin(),
		GitRoot:        findGitRoot(workDir),
		HomeDir:        homeDir,
		UserConfigPath: filepath.Join(homeDir, userConfigRelPath),
	}
	warnings := newDeprecatedConfigWarnings(opts.Stderr)
	if err := applyConfigFile(&loaded.Config, loaded.UserConfigPath, warnings, ProfileSourceUser); err != nil {
		return Loaded{}, err
	}

	if loaded.GitRoot != "" {
		loaded.ProjectConfigPath = filepath.Join(loaded.GitRoot, projectConfigName)
		if err := applyConfigFile(&loaded.Config, loaded.ProjectConfigPath, warnings, ProfileSourceProject); err != nil {
			return Loaded{}, err
		}
	}

	if err := Validate(loaded.Config); err != nil {
		return Loaded{}, err
	}
	worktreeLocation, err := ResolveWorktreeLocation(loaded.Config.Worktree.Location, loaded.GitRoot, loaded.HomeDir)
	if err != nil {
		return Loaded{}, err
	}
	loaded.Config.Worktree.Location = worktreeLocation
	return loaded, nil
}

// ResolveConfigProposal resolves proposed User Config and Project Config bytes
// with the same precedence and validation as Load. A nil scope is absent.
func ResolveConfigProposal(userContent []byte, projectContent []byte) (Config, error) {
	config := Builtin()
	warnings := newDeprecatedConfigWarnings(io.Discard)
	if userContent != nil {
		if err := applyConfigContent(&config, "User Config proposal", userContent, warnings, ProfileSourceUser); err != nil {
			return Config{}, err
		}
	}
	if projectContent != nil {
		if err := applyConfigContent(&config, "Project Config proposal", projectContent, warnings, ProfileSourceProject); err != nil {
			return Config{}, err
		}
	}
	if err := Validate(config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func Init(ctx context.Context, opts InitOptions) (InitResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	scope := strings.ToLower(strings.TrimSpace(opts.Scope))
	if scope != InitScopeUser && scope != InitScopeProject {
		return InitResult{}, fmt.Errorf("unsupported init scope %q; supported values: user, project", opts.Scope)
	}
	if err := ctx.Err(); err != nil {
		return InitResult{}, err
	}

	homeDir, err := resolveHomeDir(opts.HomeDir)
	if err != nil {
		return InitResult{}, err
	}
	workDir, err := resolveWorkDir(opts.WorkDir)
	if err != nil {
		return InitResult{}, err
	}

	path := filepath.Join(homeDir, userConfigRelPath)
	if scope == InitScopeProject {
		gitRoot := findGitRoot(workDir)
		if gitRoot == "" {
			return InitResult{}, errors.New("project init requires a Git root; use --scope user outside a repository")
		}
		path = filepath.Join(gitRoot, projectConfigName)
	}

	overwritten, err := writeDefaultConfig(ctx, path, opts.Force)
	if err != nil {
		return InitResult{}, err
	}
	return InitResult{Scope: scope, Path: path, Overwritten: overwritten}, nil
}

func DefaultConfigYAML() string {
	config := Builtin()
	return fmt.Sprintf(`# Roundfix config.
# User Config: ~/.roundfix/config.yml
# Project Config: <repo>/.roundfixrc.yml

defaults:
  # false keeps each ACP Runtime's normal sandbox or permission mode.
  agent_full_access: %t
  # Verification command for review Batches; Spec Tasks use their task file commands.
  verification: %s
  # Empty uses Roundfix Home artifacts/<repo-id>; set a path to override.
  artifact_dir: ""
  auto_commit: %t

profiles:
  general:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.5
        reasoning_effort: xhigh
  backend:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.5
        reasoning_effort: xhigh
  frontend:
    preferred:
      runtime: claude
      model: opus
      reasoning_effort: xhigh
    fallbacks:
      - runtime: codex
        model: gpt-5.6-sol
        reasoning_effort: high
  qa:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.5
        reasoning_effort: xhigh
  review:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.5
        reasoning_effort: xhigh

specs:
  # Directory holding Spec folders; relative paths resolve against the repository root.
  root: %q

worktree:
  # Parent directory; Roundfix always appends <repo-slug>/<run-id>.
  location: %q
  # Maximum concurrent Task Worktrees for spec Runs; 1 keeps sequential behavior.
  concurrency: %d
  # Repository-relative untracked files copied into each Run Worktree.
  # Empty copies no files.
  copy: []
  # Command run once after copy before Agent work; empty disables bootstrap.
  bootstrap: ""
  # Maximum time allowed for each worktree bootstrap command.
  bootstrap_timeout: %s

verification:
  # Maximum concurrent Task Verification attempts per spec Run; independent from worktree.concurrency.
  concurrency: %d

store:
  # Terminal Run journals older than this duration are eligible for pruning; 0 keeps everything.
  journal_retention: %s

review_source:
  name: %s
  # false excludes CodeRabbit findings whose severity is nitpick.
  include_nitpicks: %t
  # true asks the Review Source to review each pushed Round head.
  request_review: %t
  request_command: %q

watch:
  until_clean: %t
  max_rounds: %d
  poll_interval: %s
  review_timeout: %s
  check_grace_period: %s
  quiet_period: %s
  # auto_push runs only after no Unresolved Review Issues remain.
  auto_push: %t
  # Leave empty to use the branch upstream detected by Preflight Validation.
  push_remote: ""
  push_branch: ""

implement:
  # auto_push runs only after a Clean spec Run and never creates pull requests.
  auto_push: %t

notify:
  # false disables notifications entirely.
  enabled: %t
  # Shell command run on terminal outcome; empty uses the native desktop notification.
  command: %q

budget:
  enabled: %t
  max_run_duration: %s

resolve:
  batch_size: %d
`,
		config.Defaults.AgentFullAccess,
		config.Defaults.Verification,
		config.Defaults.AutoCommit,
		config.Specs.Root,
		config.Worktree.Location,
		config.Worktree.Concurrency,
		formatConfigDuration(config.Worktree.BootstrapTimeout),
		config.Verification.Concurrency,
		formatConfigDuration(config.Store.JournalRetention),
		config.ReviewSource.Name,
		config.ReviewSource.IncludeNitpicks,
		config.ReviewSource.RequestReview,
		config.ReviewSource.RequestCommand,
		config.Watch.UntilClean,
		config.Watch.MaxRounds,
		formatConfigDuration(config.Watch.PollInterval),
		formatConfigDuration(config.Watch.ReviewTimeout),
		formatConfigDuration(config.Watch.CheckGracePeriod),
		formatConfigDuration(config.Watch.QuietPeriod),
		config.Watch.AutoPush,
		config.Implement.AutoPush,
		config.Notify.Enabled,
		config.Notify.Command,
		config.Budget.Enabled,
		formatConfigDuration(config.Budget.MaxRunDuration),
		config.Resolve.BatchSize,
	)
}

func Validate(config Config) error {
	if config.Defaults.Agent != "" && !isSupportedAgent(config.Defaults.Agent) {
		return fmt.Errorf("defaults.agent %q is invalid; supported values: codex, claude, opencode", config.Defaults.Agent)
	}
	if err := validateProfiles(config.Profiles); err != nil {
		return err
	}
	if strings.TrimSpace(config.Defaults.Verification) == "" {
		return errors.New("defaults.verification must not be empty")
	}
	if config.ReviewSource.Name != defaultReviewSource {
		return fmt.Errorf("review_source.name %q is invalid; supported value: coderabbit", config.ReviewSource.Name)
	}
	if strings.TrimSpace(config.ReviewSource.RequestCommand) == "" {
		return errors.New("review_source.request_command must not be empty")
	}
	if config.Watch.MaxRounds < 1 {
		return errors.New("watch.max_rounds must be greater than 0")
	}
	if config.Watch.PollInterval <= 0 {
		return errors.New("watch.poll_interval must be greater than 0")
	}
	if config.Watch.ReviewTimeout <= 0 {
		return errors.New("watch.review_timeout must be greater than 0")
	}
	if config.Watch.CheckGracePeriod <= 0 {
		return errors.New("watch.check_grace_period must be greater than 0")
	}
	if config.Watch.QuietPeriod <= 0 {
		return errors.New("watch.quiet_period must be greater than 0")
	}
	if config.Budget.MaxRunDuration <= 0 {
		return errors.New("budget.max_run_duration must be greater than 0")
	}
	if config.Resolve.BatchSize < 1 {
		return errors.New("resolve.batch_size must be greater than 0")
	}
	if config.Store.JournalRetention < 0 {
		return errors.New("store.journal_retention must be greater than or equal to 0")
	}
	if strings.TrimSpace(config.Specs.Root) == "" {
		return errors.New("specs.root must not be empty")
	}
	if config.Worktree.Concurrency < 1 {
		return errors.New("worktree.concurrency must be greater than 0")
	}
	if config.Verification.Concurrency < 1 {
		return errors.New("verification.concurrency must be greater than 0")
	}
	if strings.TrimSpace(config.Worktree.Location) == "" {
		return errors.New("worktree.location must not be empty")
	}
	if config.Worktree.BootstrapTimeout <= 0 {
		return errors.New("worktree.bootstrap_timeout must be greater than 0")
	}
	for _, path := range config.Worktree.Copy {
		if err := validateWorktreeCopyPath(path); err != nil {
			return err
		}
	}
	if config.Watch.AutoPush && !config.Defaults.AutoCommit {
		return errors.New("watch.auto_push requires defaults.auto_commit to be true")
	}
	return nil
}

func ResolveSpecsRoot(loaded Loaded, repoRoot string) (SpecsRoot, error) {
	effectiveRepoRoot := strings.TrimSpace(repoRoot)
	if effectiveRepoRoot == "" {
		effectiveRepoRoot = strings.TrimSpace(loaded.GitRoot)
	}
	if effectiveRepoRoot == "" {
		return SpecsRoot{}, errors.New("specs.root requires a Git root")
	}
	absoluteRepoRoot, err := filepath.Abs(effectiveRepoRoot)
	if err != nil {
		return SpecsRoot{}, fmt.Errorf("resolve repository root %q: %w", effectiveRepoRoot, err)
	}

	configuredRoot := strings.TrimSpace(loaded.Config.Specs.Root)
	if configuredRoot == "" {
		return SpecsRoot{}, errors.New("specs.root must not be empty")
	}
	resolved := filepath.Clean(configuredRoot)
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(absoluteRepoRoot, resolved)
	}

	info, err := os.Stat(resolved)
	if errors.Is(err, os.ErrNotExist) {
		return SpecsRoot{}, fmt.Errorf("specs.root resolved to %q, which does not exist; create the directory or update specs.root", resolved)
	}
	if err != nil {
		return SpecsRoot{}, fmt.Errorf("stat specs.root %q: %w", resolved, err)
	}
	if !info.IsDir() {
		return SpecsRoot{}, fmt.Errorf("specs.root resolved to %q, which is not a directory; update specs.root to a directory", resolved)
	}

	evaluatedRoot, err := filepath.EvalSymlinks(resolved)
	if err != nil {
		return SpecsRoot{}, fmt.Errorf("evaluate specs.root %q: %w", resolved, err)
	}
	evaluatedRepoRoot, err := filepath.EvalSymlinks(absoluteRepoRoot)
	if err != nil {
		return SpecsRoot{}, fmt.Errorf("evaluate repository root %q: %w", absoluteRepoRoot, err)
	}

	external := !pathInsideOrSame(evaluatedRoot, evaluatedRepoRoot)
	builtInRoot := !external && filepath.Clean(resolved) == filepath.Clean(filepath.Join(absoluteRepoRoot, filepath.FromSlash(defaultSpecsRoot)))

	return SpecsRoot{
		Path:        resolved,
		External:    external,
		BuiltInRoot: builtInRoot,
	}, nil
}

func writeDefaultConfig(ctx context.Context, path string, force bool) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	overwritten := false
	if _, err := os.Stat(path); err == nil {
		if !force {
			return false, fmt.Errorf("config already exists at %q; pass --force to overwrite", path)
		}
		overwritten = true
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, fmt.Errorf("stat config %q: %w", path, err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, fmt.Errorf("create config directory %q: %w", filepath.Dir(path), err)
	}
	flags := os.O_WRONLY | os.O_CREATE
	if force {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}
	file, err := os.OpenFile(path, flags, 0o644)
	if errors.Is(err, os.ErrExist) {
		return false, fmt.Errorf("config already exists at %q; pass --force to overwrite", path)
	}
	if err != nil {
		return false, fmt.Errorf("create config %q: %w", path, err)
	}
	_, writeErr := file.WriteString(DefaultConfigYAML())
	closeErr := file.Close()
	if writeErr != nil {
		return false, fmt.Errorf("write config %q: %w", path, writeErr)
	}
	if closeErr != nil {
		return false, fmt.Errorf("write config %q: %w", path, closeErr)
	}
	return overwritten, nil
}

func formatConfigDuration(duration time.Duration) string {
	switch {
	case duration%time.Hour == 0:
		return fmt.Sprintf("%dh", int(duration/time.Hour))
	case duration%time.Minute == 0:
		return fmt.Sprintf("%dm", int(duration/time.Minute))
	case duration%time.Second == 0:
		return fmt.Sprintf("%ds", int(duration/time.Second))
	default:
		return duration.String()
	}
}

func ValidateArtifactDirectory(artifactDir string, gitRoot string, homeDir string) (string, error) {
	resolved, err := ResolveArtifactDirectory(artifactDir, gitRoot, homeDir)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(resolved)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(resolved, 0o755); err != nil {
			return "", fmt.Errorf("create Artifact Directory %q: %w", resolved, err)
		}
		info, err = os.Stat(resolved)
	}
	if err != nil {
		return "", fmt.Errorf("stat Artifact Directory %q: %w", resolved, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("Artifact Directory %q is not a directory", resolved)
	}

	temp, err := os.CreateTemp(resolved, ".roundfix-write-check-*")
	if err != nil {
		return "", fmt.Errorf("write-check Artifact Directory %q: %w", resolved, err)
	}
	tempPath := temp.Name()
	closeErr := temp.Close()
	removeErr := os.Remove(tempPath)
	if closeErr != nil {
		return "", fmt.Errorf("write-check Artifact Directory %q: %w", resolved, closeErr)
	}
	if removeErr != nil {
		return "", fmt.Errorf("remove write-check file %q: %w", tempPath, removeErr)
	}
	return resolved, nil
}

func ResolveArtifactDirectory(artifactDir string, gitRoot string, homeDir string) (string, error) {
	expanded, err := expandHome(artifactDir, homeDir)
	if err != nil {
		return "", err
	}
	if expanded == "" {
		if gitRoot == "" {
			return "", errors.New("empty artifact_dir requires a Git root")
		}
		if homeDir == "" {
			return "", errors.New("empty artifact_dir requires Roundfix Home")
		}
		return filepath.Join(homeDir, ".roundfix", "artifacts", repoID(gitRoot)), nil
	}
	if filepath.IsAbs(expanded) {
		return filepath.Clean(expanded), nil
	}
	if gitRoot == "" {
		return "", fmt.Errorf("relative artifact_dir %q requires a Git root", artifactDir)
	}
	return filepath.Join(gitRoot, expanded), nil
}

// ReviewArtifactContext identifies the review artifact root for one Open Pull
// Request. ExplicitArtifactDir must already be resolved like Artifact Directory.
type ReviewArtifactContext struct {
	ExplicitArtifactDir string
	RepoRoot            string
	SpecsRoot           string
	SpecSlug            string
	PRNumber            int
}

func ResolveReviewRoot(ctx ReviewArtifactContext) (string, error) {
	if ctx.PRNumber < 1 {
		return "", errors.New("PR number is required to resolve Review artifact root")
	}
	prDir := fmt.Sprintf("pr-%d", ctx.PRNumber)
	if explicit := strings.TrimSpace(ctx.ExplicitArtifactDir); explicit != "" {
		resolved := filepath.Join(explicit, "reviews", prDir)
		if strings.TrimSpace(ctx.RepoRoot) == "" {
			return "", errors.New("repository root is required to keep Review artifacts out of history")
		}
		if reviewRootInsideHistory(ctx.RepoRoot, resolved) {
			return "", fmt.Errorf("Review artifact root %q must not be inside repository history", resolved)
		}
		return resolved, nil
	}

	specsRoot := strings.TrimSpace(ctx.SpecsRoot)
	if specsRoot == "" {
		repoRoot := strings.TrimSpace(ctx.RepoRoot)
		if repoRoot == "" {
			return "", errors.New("Spec Root is required to resolve Review artifact root")
		}
		specsRoot = filepath.Join(repoRoot, "docs", "specs")
	} else if strings.TrimSpace(ctx.RepoRoot) == "" {
		return "", errors.New("repository root is required to keep Review artifacts out of history")
	}
	if reviewRootInsideHistory(ctx.RepoRoot, specsRoot) {
		repoRoot := strings.TrimSpace(ctx.RepoRoot)
		if repoRoot == "" {
			return "", errors.New("repository root is required to keep Review artifacts out of history")
		}
		specsRoot = filepath.Join(repoRoot, "docs", "specs")
	}
	if slug := strings.TrimSpace(ctx.SpecSlug); slug != "" && reviewSpecDirExists(specsRoot, slug) {
		return filepath.Join(specsRoot, slug, "reviews"), nil
	}
	live := filepath.Join(specsRoot, "reviews", prDir)
	if reviewPRDirExists(specsRoot, prDir) {
		return live, nil
	}
	if legacy := filepath.Join(specsRoot, "_reviews", prDir); reviewPRDirExistsAt(legacy) {
		return legacy, nil
	}
	return live, nil
}

func reviewPRDirExists(specsRoot string, prDir string) bool {
	return reviewPRDirExistsAt(filepath.Join(specsRoot, "reviews", prDir))
}

func reviewPRDirExistsAt(prDirPath string) bool {
	info, err := os.Stat(prDirPath)
	return err == nil && info.IsDir()
}

func reviewRootInsideHistory(repoRoot string, path string) bool {
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return false
	}
	return pathInsideOrSame(filepath.Clean(path), filepath.Join(filepath.Clean(repoRoot), "docs", "history"))
}

func reviewSpecDirExists(specsRoot string, slug string) bool {
	if !validReviewSpecSlug(slug) {
		return false
	}
	info, err := os.Stat(filepath.Join(specsRoot, slug))
	return err == nil && info.IsDir()
}

func validReviewSpecSlug(slug string) bool {
	if slug == "" || slug == "." || slug == ".." || filepath.IsAbs(slug) {
		return false
	}
	return !strings.ContainsAny(slug, `/\`)
}

func ResolveWorktreeLocation(location string, gitRoot string, homeDir string) (string, error) {
	expanded, err := expandHomeForKey("worktree.location", strings.TrimSpace(location), homeDir)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(expanded) == "" {
		return "", errors.New("worktree.location must not be empty")
	}
	if !filepath.IsAbs(expanded) {
		return "", errors.New("worktree.location must be absolute after ~ expansion")
	}
	resolved := filepath.Clean(expanded)
	if strings.TrimSpace(gitRoot) != "" && pathInsideOrSame(resolved, gitRoot) {
		return "", fmt.Errorf("worktree.location must not be inside the repository tree: %q", resolved)
	}
	return resolved, nil
}

func applyConfigFile(config *Config, path string, warnings *deprecatedConfigWarnings, source ProfileSource) error {
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read config %q: %w", path, err)
	}
	return applyConfigContent(config, path, content, warnings, source)
}

func applyConfigContent(config *Config, label string, content []byte, warnings *deprecatedConfigWarnings, source ProfileSource) error {
	var document yaml.Node
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	if err := decoder.Decode(&document); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf("parse config %q: %w", label, err)
	}
	stripDeprecatedConfigKeys(&document, warnings)
	if value, found := yamlValueAtPath(&document, []string{"review_source", "request_review"}); found && value.Tag == "!!null" {
		return fmt.Errorf("parse config %q: review_source.request_review must be boolean: cannot unmarshal null value", label)
	}
	hasProfiles := configHasProfilesSection(&document)
	hasLegacyRuntimeDefaults := configHasLegacyRuntimeDefaults(&document)
	if hasProfiles && hasLegacyRuntimeDefaults {
		return fmt.Errorf("parse config %q: %w", label, profileSchemaConflictError(source))
	}
	cleaned, err := encodeYAMLNode(&document)
	if err != nil {
		return fmt.Errorf("parse config %q: %w", label, err)
	}

	var overlay configOverlay
	decoder = yaml.NewDecoder(bytes.NewReader(cleaned))
	decoder.KnownFields(true)
	if err := decoder.Decode(&overlay); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf("parse config %q: %w", label, err)
	}
	applyOverlay(config, overlay)
	if overlay.Profiles != nil {
		applyProfilesOverlay(config, overlay.Profiles, source)
	} else if hasLegacyRuntimeDefaults {
		applyLegacyRuntimeProfiles(config, source)
	}
	return nil
}

func stripDeprecatedConfigKeys(document *yaml.Node, warnings *deprecatedConfigWarnings) {
	for _, key := range deprecatedConfigKeys {
		if removeYAMLPath(document, key.path) {
			warnings.warn(key)
		}
	}
}

func removeYAMLPath(node *yaml.Node, path []string) bool {
	if node == nil || len(path) == 0 {
		return false
	}
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return false
		}
		return removeYAMLPath(node.Content[0], path)
	}
	if node.Kind != yaml.MappingNode {
		return false
	}
	removed := false
	for index := 0; index+1 < len(node.Content); {
		key := node.Content[index]
		value := node.Content[index+1]
		if key.Value != path[0] {
			index += 2
			continue
		}
		if len(path) == 1 {
			node.Content = append(node.Content[:index], node.Content[index+2:]...)
			removed = true
			continue
		}
		if removeYAMLPath(value, path[1:]) {
			removed = true
		}
		index += 2
	}
	return removed
}

func yamlValueAtPath(node *yaml.Node, path []string) (*yaml.Node, bool) {
	if node == nil || len(path) == 0 {
		return nil, false
	}
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return nil, false
		}
		return yamlValueAtPath(node.Content[0], path)
	}
	if node.Kind != yaml.MappingNode {
		return nil, false
	}
	for index := 0; index+1 < len(node.Content); index += 2 {
		key := node.Content[index]
		value := node.Content[index+1]
		if key.Value != path[0] {
			continue
		}
		if len(path) == 1 {
			return value, true
		}
		return yamlValueAtPath(value, path[1:])
	}
	return nil, false
}

func encodeYAMLNode(node *yaml.Node) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	if err := encoder.Encode(node); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func applyOverlay(config *Config, overlay configOverlay) {
	if overlay.Defaults != nil {
		if overlay.Defaults.Agent != nil {
			config.Defaults.Agent = *overlay.Defaults.Agent
		}
		if overlay.Defaults.AgentFullAccess != nil {
			config.Defaults.AgentFullAccess = *overlay.Defaults.AgentFullAccess
		}
		if overlay.Defaults.AutoCommit != nil {
			config.Defaults.AutoCommit = *overlay.Defaults.AutoCommit
		}
		if overlay.Defaults.Verification != nil {
			config.Defaults.Verification = *overlay.Defaults.Verification
		}
		if overlay.Defaults.ArtifactDir != nil {
			config.Defaults.ArtifactDir = *overlay.Defaults.ArtifactDir
		}
	}
	if overlay.Runtimes != nil {
		if overlay.Runtimes.Codex != nil {
			applyRuntimeDefaultsOverlay(&config.Runtimes.Codex, *overlay.Runtimes.Codex)
		}
		if overlay.Runtimes.Claude != nil {
			applyRuntimeDefaultsOverlay(&config.Runtimes.Claude, *overlay.Runtimes.Claude)
		}
		if overlay.Runtimes.OpenCode != nil {
			applyRuntimeDefaultsOverlay(&config.Runtimes.OpenCode, *overlay.Runtimes.OpenCode)
		}
	}
	if overlay.ReviewSource != nil {
		if overlay.ReviewSource.Name != nil {
			config.ReviewSource.Name = *overlay.ReviewSource.Name
		}
		if overlay.ReviewSource.IncludeNitpicks != nil {
			config.ReviewSource.IncludeNitpicks = *overlay.ReviewSource.IncludeNitpicks
		}
		if overlay.ReviewSource.RequestReview != nil {
			config.ReviewSource.RequestReview = overlay.ReviewSource.RequestReview.value
		}
		if overlay.ReviewSource.RequestCommand != nil {
			config.ReviewSource.RequestCommand = *overlay.ReviewSource.RequestCommand
		}
	}
	if overlay.Watch != nil {
		if overlay.Watch.UntilClean != nil {
			config.Watch.UntilClean = *overlay.Watch.UntilClean
		}
		if overlay.Watch.MaxRounds != nil {
			config.Watch.MaxRounds = *overlay.Watch.MaxRounds
		}
		if overlay.Watch.PollInterval != nil {
			config.Watch.PollInterval = overlay.Watch.PollInterval.value
		}
		if overlay.Watch.ReviewTimeout != nil {
			config.Watch.ReviewTimeout = overlay.Watch.ReviewTimeout.value
		}
		if overlay.Watch.CheckGracePeriod != nil {
			config.Watch.CheckGracePeriod = overlay.Watch.CheckGracePeriod.value
		}
		if overlay.Watch.QuietPeriod != nil {
			config.Watch.QuietPeriod = overlay.Watch.QuietPeriod.value
		}
		if overlay.Watch.AutoPush != nil {
			config.Watch.AutoPush = *overlay.Watch.AutoPush
		}
		if overlay.Watch.PushRemote != nil {
			config.Watch.PushRemote = *overlay.Watch.PushRemote
		}
		if overlay.Watch.PushBranch != nil {
			config.Watch.PushBranch = *overlay.Watch.PushBranch
		}
	}
	if overlay.Implement != nil {
		if overlay.Implement.AutoPush != nil {
			config.Implement.AutoPush = overlay.Implement.AutoPush.value
		}
	}
	if overlay.Notify != nil {
		if overlay.Notify.Enabled != nil {
			config.Notify.Enabled = *overlay.Notify.Enabled
		}
		if overlay.Notify.Command != nil {
			config.Notify.Command = *overlay.Notify.Command
		}
	}
	if overlay.Worktree != nil {
		if overlay.Worktree.Concurrency != nil {
			config.Worktree.Concurrency = *overlay.Worktree.Concurrency
		}
		if overlay.Worktree.Location != nil {
			config.Worktree.Location = *overlay.Worktree.Location
		}
		if overlay.Worktree.Copy != nil {
			config.Worktree.Copy = append([]string(nil), (*overlay.Worktree.Copy)...)
		}
		if overlay.Worktree.Bootstrap != nil {
			config.Worktree.Bootstrap = *overlay.Worktree.Bootstrap
		}
		if overlay.Worktree.BootstrapTimeout != nil {
			config.Worktree.BootstrapTimeout = overlay.Worktree.BootstrapTimeout.value
		}
	}
	if overlay.Verification != nil {
		if overlay.Verification.Concurrency != nil {
			config.Verification.Concurrency = overlay.Verification.Concurrency.value
		}
	}
	if overlay.Budget != nil {
		if overlay.Budget.Enabled != nil {
			config.Budget.Enabled = *overlay.Budget.Enabled
		}
		if overlay.Budget.MaxRunDuration != nil {
			config.Budget.MaxRunDuration = overlay.Budget.MaxRunDuration.value
		}
	}
	if overlay.Resolve != nil {
		if overlay.Resolve.BatchSize != nil {
			config.Resolve.BatchSize = *overlay.Resolve.BatchSize
		}
	}
	if overlay.Logs != nil {
		if overlay.Logs.Agent != nil {
			config.Logs.Agent = *overlay.Logs.Agent
		}
	}
	if overlay.Store != nil {
		if overlay.Store.JournalRetention != nil {
			config.Store.JournalRetention = overlay.Store.JournalRetention.value
		}
	}
	if overlay.Specs != nil {
		if overlay.Specs.Root != nil {
			config.Specs.Root = *overlay.Specs.Root
		}
	}
}

func applyRuntimeDefaultsOverlay(defaults *RuntimeDefaults, overlay runtimeDefaultsOverlay) {
	if overlay.Model != nil {
		defaults.Model = *overlay.Model
	}
	if overlay.ReasoningEffort != nil {
		defaults.ReasoningEffort = *overlay.ReasoningEffort
	}
}

func resolveHomeDir(homeDir string) (string, error) {
	if homeDir != "" {
		return homeDir, nil
	}
	resolved, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve User Config home: %w", err)
	}
	return resolved, nil
}

func resolveWorkDir(workDir string) (string, error) {
	if workDir != "" {
		return workDir, nil
	}
	resolved, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve work directory: %w", err)
	}
	return resolved, nil
}

func findGitRoot(start string) string {
	current := filepath.Clean(start)
	for {
		if _, err := os.Stat(filepath.Join(current, ".git")); err == nil {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			return ""
		}
		current = parent
	}
}

func expandHome(path string, homeDir string) (string, error) {
	return expandHomeForKey("artifact_dir", path, homeDir)
}

func expandHomeForKey(key string, path string, homeDir string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path, nil
	}
	if homeDir == "" {
		return "", fmt.Errorf("%s uses ~ but home directory is unavailable", key)
	}
	if path == "~" {
		return homeDir, nil
	}
	return filepath.Join(homeDir, strings.TrimPrefix(path, "~/")), nil
}

func pathInsideOrSame(path string, root string) bool {
	path = filepath.Clean(path)
	root = filepath.Clean(root)
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel))
}

func validateWorktreeCopyPath(path string) error {
	if path == "" {
		return errors.New("worktree.copy entries must not be empty")
	}
	if filepath.IsAbs(path) {
		return fmt.Errorf("worktree.copy entry %q must be repository-relative", path)
	}
	clean := filepath.Clean(path)
	if clean == "." || clean != path || containsDotDotSegment(clean) {
		return fmt.Errorf("worktree.copy entry %q must be clean and stay inside the repository", path)
	}
	return nil
}

func containsDotDotSegment(path string) bool {
	for _, segment := range strings.Split(filepath.ToSlash(path), "/") {
		if segment == ".." {
			return true
		}
	}
	return false
}

func repoID(gitRoot string) string {
	sum := sha256.Sum256([]byte(filepath.Clean(gitRoot)))
	return hex.EncodeToString(sum[:])[:16]
}

func isSupportedAgent(agent string) bool {
	switch agent {
	case "codex", "claude", "opencode":
		return true
	default:
		return false
	}
}
