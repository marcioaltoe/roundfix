package config

import (
	"bytes"
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
	userConfigRelPath    = ".roundfix/config.yml"
	projectConfigName    = ".roundfixrc.yml"
	defaultArtifactDir   = ".roundfix"
	defaultReviewSource  = "coderabbit"
	defaultAgent         = "codex"
	defaultVerification  = "make verify"
	defaultPollInterval  = 30 * time.Second
	defaultReviewTimeout = 30 * time.Minute
	defaultQuietPeriod   = 30 * time.Second
	defaultRunDuration   = 2 * time.Hour
)

type Config struct {
	Defaults     Defaults
	ReviewSource ReviewSource
	Watch        Watch
	Budget       Budget
	Resolve      Resolve
}

type Defaults struct {
	Agent        string
	Model        string
	AutoCommit   bool
	Verification string
	ArtifactDir  string
}

type ReviewSource struct {
	Name            string
	IncludeNitpicks bool
}

type Watch struct {
	UntilClean    bool
	MaxRounds     int
	PollInterval  time.Duration
	ReviewTimeout time.Duration
	QuietPeriod   time.Duration
	AutoPush      bool
	PushRemote    string
	PushBranch    string
}

type Budget struct {
	Enabled        bool
	MaxRunDuration time.Duration
}

type Resolve struct {
	BatchSize  int
	Concurrent int
}

type Loaded struct {
	Config            Config
	GitRoot           string
	HomeDir           string
	UserConfigPath    string
	ProjectConfigPath string
}

type LoadOptions struct {
	HomeDir string
	WorkDir string
}

type durationValue struct {
	value time.Duration
}

func (duration *durationValue) UnmarshalYAML(node *yaml.Node) error {
	var raw string
	if err := node.Decode(&raw); err != nil {
		return err
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
	ReviewSource *reviewSourceOverlay `yaml:"review_source"`
	Watch        *watchOverlay        `yaml:"watch"`
	Budget       *budgetOverlay       `yaml:"budget"`
	Resolve      *resolveOverlay      `yaml:"resolve"`
}

type defaultsOverlay struct {
	Agent        *string `yaml:"agent"`
	Model        *string `yaml:"model"`
	AutoCommit   *bool   `yaml:"auto_commit"`
	Verification *string `yaml:"verification"`
	ArtifactDir  *string `yaml:"artifact_dir"`
}

type reviewSourceOverlay struct {
	Name            *string `yaml:"name"`
	IncludeNitpicks *bool   `yaml:"include_nitpicks"`
}

type watchOverlay struct {
	UntilClean    *bool          `yaml:"until_clean"`
	MaxRounds     *int           `yaml:"max_rounds"`
	PollInterval  *durationValue `yaml:"poll_interval"`
	ReviewTimeout *durationValue `yaml:"review_timeout"`
	QuietPeriod   *durationValue `yaml:"quiet_period"`
	AutoPush      *bool          `yaml:"auto_push"`
	PushRemote    *string        `yaml:"push_remote"`
	PushBranch    *string        `yaml:"push_branch"`
}

type budgetOverlay struct {
	Enabled        *bool          `yaml:"enabled"`
	MaxRunDuration *durationValue `yaml:"max_run_duration"`
}

type resolveOverlay struct {
	BatchSize  *int `yaml:"batch_size"`
	Concurrent *int `yaml:"concurrent"`
}

func Builtin() Config {
	return Config{
		Defaults: Defaults{
			Agent:        defaultAgent,
			AutoCommit:   true,
			Verification: defaultVerification,
		},
		ReviewSource: ReviewSource{
			Name:            defaultReviewSource,
			IncludeNitpicks: true,
		},
		Watch: Watch{
			UntilClean:    true,
			MaxRounds:     6,
			PollInterval:  defaultPollInterval,
			ReviewTimeout: defaultReviewTimeout,
			QuietPeriod:   defaultQuietPeriod,
			AutoPush:      true,
		},
		Budget: Budget{
			Enabled:        true,
			MaxRunDuration: defaultRunDuration,
		},
		Resolve: Resolve{
			BatchSize:  3,
			Concurrent: 1,
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
	if err := applyConfigFile(&loaded.Config, loaded.UserConfigPath); err != nil {
		return Loaded{}, err
	}

	if loaded.GitRoot != "" {
		loaded.ProjectConfigPath = filepath.Join(loaded.GitRoot, projectConfigName)
		if err := applyConfigFile(&loaded.Config, loaded.ProjectConfigPath); err != nil {
			return Loaded{}, err
		}
	}

	if err := Validate(loaded.Config); err != nil {
		return Loaded{}, err
	}
	return loaded, nil
}

func Validate(config Config) error {
	if config.Defaults.Agent != "" && !isSupportedAgent(config.Defaults.Agent) {
		return fmt.Errorf("defaults.agent %q is invalid; supported values: codex, claude, opencode", config.Defaults.Agent)
	}
	if strings.TrimSpace(config.Defaults.Verification) == "" {
		return errors.New("defaults.verification must not be empty")
	}
	if config.ReviewSource.Name != defaultReviewSource {
		return fmt.Errorf("review_source.name %q is invalid; supported value: coderabbit", config.ReviewSource.Name)
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
	if config.Watch.QuietPeriod <= 0 {
		return errors.New("watch.quiet_period must be greater than 0")
	}
	if config.Budget.MaxRunDuration <= 0 {
		return errors.New("budget.max_run_duration must be greater than 0")
	}
	if config.Resolve.BatchSize < 1 {
		return errors.New("resolve.batch_size must be greater than 0")
	}
	if config.Resolve.Concurrent < 1 {
		return errors.New("resolve.concurrent must be greater than 0")
	}
	if config.Watch.AutoPush && !config.Defaults.AutoCommit {
		return errors.New("watch.auto_push requires defaults.auto_commit to be true")
	}
	return nil
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
		return filepath.Join(gitRoot, defaultArtifactDir), nil
	}
	if filepath.IsAbs(expanded) {
		return filepath.Clean(expanded), nil
	}
	if gitRoot == "" {
		return "", fmt.Errorf("relative artifact_dir %q requires a Git root", artifactDir)
	}
	return filepath.Join(gitRoot, expanded), nil
}

func applyConfigFile(config *Config, path string) error {
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read config %q: %w", path, err)
	}

	var overlay configOverlay
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)
	if err := decoder.Decode(&overlay); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf("parse config %q: %w", path, err)
	}
	applyOverlay(config, overlay)
	return nil
}

func applyOverlay(config *Config, overlay configOverlay) {
	if overlay.Defaults != nil {
		if overlay.Defaults.Agent != nil {
			config.Defaults.Agent = *overlay.Defaults.Agent
		}
		if overlay.Defaults.Model != nil {
			config.Defaults.Model = *overlay.Defaults.Model
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
	if overlay.ReviewSource != nil {
		if overlay.ReviewSource.Name != nil {
			config.ReviewSource.Name = *overlay.ReviewSource.Name
		}
		if overlay.ReviewSource.IncludeNitpicks != nil {
			config.ReviewSource.IncludeNitpicks = *overlay.ReviewSource.IncludeNitpicks
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
		if overlay.Resolve.Concurrent != nil {
			config.Resolve.Concurrent = *overlay.Resolve.Concurrent
		}
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
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path, nil
	}
	if homeDir == "" {
		return "", errors.New("artifact_dir uses ~ but home directory is unavailable")
	}
	if path == "~" {
		return homeDir, nil
	}
	return filepath.Join(homeDir, strings.TrimPrefix(path, "~/")), nil
}

func isSupportedAgent(agent string) bool {
	switch agent {
	case "codex", "claude", "opencode":
		return true
	default:
		return false
	}
}
