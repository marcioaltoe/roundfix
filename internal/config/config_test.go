package config

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadAppliesConfigPrecedence(t *testing.T) {
	homeDir := t.TempDir()
	workDir := t.TempDir()
	userWorktreeLocation := filepath.Join(homeDir, "configured-user-worktrees")
	projectWorktreeLocation := filepath.Join(homeDir, "configured-project-worktrees")
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
defaults:
  agent: claude
  agent_full_access: true
  artifact_dir: user-artifacts
watch:
  max_rounds: 4
  poll_interval: 10s
implement:
  auto_push: true
worktree:
  concurrency: 4
  location: `+userWorktreeLocation+`
  copy:
    - .env.local
    - certs/dev.pem
resolve:
  batch_size: 2
specs:
  root: user-specs
`)
	mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), `
defaults:
  agent: opencode
watch:
  max_rounds: 8
implement:
  auto_push: false
worktree:
  concurrency: 1
  location: `+projectWorktreeLocation+`
  copy:
    - project.env
budget:
  max_run_duration: 3h
specs:
  root: project-specs
`)

	loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
	if err != nil {
		t.Fatalf("expected config to load, got %v", err)
	}

	if loaded.Config.Defaults.Agent != "opencode" {
		t.Fatalf("expected project config to override user config, got %q", loaded.Config.Defaults.Agent)
	}
	if loaded.Config.Defaults.ArtifactDir != "user-artifacts" {
		t.Fatalf("expected user artifact dir to survive project config, got %q", loaded.Config.Defaults.ArtifactDir)
	}
	if !loaded.Config.Defaults.AgentFullAccess {
		t.Fatal("expected user agent_full_access to survive project config")
	}
	if loaded.Config.Watch.MaxRounds != 8 {
		t.Fatalf("expected project max rounds, got %d", loaded.Config.Watch.MaxRounds)
	}
	if loaded.Config.Watch.PollInterval != 10*time.Second {
		t.Fatalf("expected user poll interval, got %s", loaded.Config.Watch.PollInterval)
	}
	if loaded.Config.Implement.AutoPush {
		t.Fatal("expected project implement.auto_push to override user default")
	}
	if len(loaded.Config.Worktree.Copy) != 1 || loaded.Config.Worktree.Copy[0] != "project.env" {
		t.Fatalf("expected project worktree.copy to override user list, got %#v", loaded.Config.Worktree.Copy)
	}
	if loaded.Config.Worktree.Concurrency != 1 {
		t.Fatalf("expected project worktree.concurrency to override user config, got %d", loaded.Config.Worktree.Concurrency)
	}
	if loaded.Config.Worktree.Location != projectWorktreeLocation {
		t.Fatalf("expected project worktree.location to override user config, got %q", loaded.Config.Worktree.Location)
	}
	if loaded.Config.Budget.MaxRunDuration != 3*time.Hour {
		t.Fatalf("expected project max run duration, got %s", loaded.Config.Budget.MaxRunDuration)
	}
	if loaded.Config.Resolve.BatchSize != 2 {
		t.Fatalf("expected user batch size, got %d", loaded.Config.Resolve.BatchSize)
	}
	if loaded.Config.Specs.Root != "project-specs" {
		t.Fatalf("expected project specs.root to override user config, got %q", loaded.Config.Specs.Root)
	}
}

func TestLoadAppliesWorktreeConfigHierarchy(t *testing.T) {
	tests := []struct {
		name            string
		userConfig      string
		projectConfig   string
		wantConcurrency int
		wantBootstrap   string
		wantTimeout     time.Duration
		wantLocation    func(homeDir, workDir string) string
	}{
		{
			name:            "builtin only",
			wantConcurrency: 2,
			wantTimeout:     10 * time.Minute,
			wantLocation: func(homeDir, workDir string) string {
				return filepath.Join(homeDir, ".roundfix", "worktrees")
			},
		},
		{
			name: "user override",
			userConfig: `
worktree:
  concurrency: 4
  location: __USER_LOCATION__
  bootstrap: make user-bootstrap
  bootstrap_timeout: 2m
`,
			wantConcurrency: 4,
			wantBootstrap:   "make user-bootstrap",
			wantTimeout:     2 * time.Minute,
			wantLocation: func(homeDir, workDir string) string {
				return filepath.Join(homeDir, "user-worktrees")
			},
		},
		{
			name: "project override",
			userConfig: `
worktree:
  concurrency: 4
  location: __USER_LOCATION__
  bootstrap: make user-bootstrap
  bootstrap_timeout: 2m
`,
			projectConfig: `
worktree:
  concurrency: 1
  location: __PROJECT_LOCATION__
  bootstrap: make project-bootstrap
  bootstrap_timeout: 30s
`,
			wantConcurrency: 1,
			wantBootstrap:   "make project-bootstrap",
			wantTimeout:     30 * time.Second,
			wantLocation: func(homeDir, workDir string) string {
				return filepath.Join(homeDir, "project-worktrees")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			workDir := t.TempDir()
			userLocation := filepath.Join(homeDir, "user-worktrees")
			projectLocation := filepath.Join(homeDir, "project-worktrees")
			mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
			mustMkdir(t, filepath.Join(workDir, ".git"))
			if strings.TrimSpace(tt.userConfig) != "" {
				content := strings.ReplaceAll(tt.userConfig, "__USER_LOCATION__", userLocation)
				mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), content)
			}
			if strings.TrimSpace(tt.projectConfig) != "" {
				content := strings.ReplaceAll(tt.projectConfig, "__PROJECT_LOCATION__", projectLocation)
				mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), content)
			}

			loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
			if err != nil {
				t.Fatalf("expected config to load, got %v", err)
			}

			if loaded.Config.Worktree.Concurrency != tt.wantConcurrency {
				t.Fatalf("expected worktree.concurrency %d, got %d", tt.wantConcurrency, loaded.Config.Worktree.Concurrency)
			}
			if loaded.Config.Worktree.Bootstrap != tt.wantBootstrap {
				t.Fatalf("expected worktree.bootstrap %q, got %q", tt.wantBootstrap, loaded.Config.Worktree.Bootstrap)
			}
			if loaded.Config.Worktree.BootstrapTimeout != tt.wantTimeout {
				t.Fatalf("expected worktree.bootstrap_timeout %s, got %s", tt.wantTimeout, loaded.Config.Worktree.BootstrapTimeout)
			}
			if want := tt.wantLocation(homeDir, workDir); loaded.Config.Worktree.Location != want {
				t.Fatalf("expected worktree.location %q, got %q", want, loaded.Config.Worktree.Location)
			}
		})
	}
}

func TestLoadAppliesAgentLogConfigHierarchy(t *testing.T) {
	tests := []struct {
		name          string
		userConfig    string
		projectConfig string
		want          bool
	}{
		{
			name: "builtin only",
			want: false,
		},
		{
			name: "user enables",
			userConfig: `
logs:
  agent: true
`,
			want: true,
		},
		{
			name: "project enables",
			projectConfig: `
logs:
  agent: true
`,
			want: true,
		},
		{
			name: "project overrides user",
			userConfig: `
logs:
  agent: true
`,
			projectConfig: `
logs:
  agent: false
`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			workDir := t.TempDir()
			mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
			mustMkdir(t, filepath.Join(workDir, ".git"))
			if strings.TrimSpace(tt.userConfig) != "" {
				mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), tt.userConfig)
			}
			if strings.TrimSpace(tt.projectConfig) != "" {
				mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), tt.projectConfig)
			}

			loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
			if err != nil {
				t.Fatalf("expected config to load, got %v", err)
			}
			if loaded.Config.Logs.Agent != tt.want {
				t.Fatalf("expected logs.agent %t, got %t", tt.want, loaded.Config.Logs.Agent)
			}
		})
	}
}

func TestLoadAppliesStoreRetentionConfigHierarchy(t *testing.T) {
	tests := []struct {
		name          string
		userConfig    string
		projectConfig string
		want          time.Duration
	}{
		{
			name: "builtin only",
			want: 336 * time.Hour,
		},
		{
			name: "user override",
			userConfig: `
store:
  journal_retention: 168h
`,
			want: 168 * time.Hour,
		},
		{
			name: "project override",
			userConfig: `
store:
  journal_retention: 168h
`,
			projectConfig: `
store:
  journal_retention: 72h
`,
			want: 72 * time.Hour,
		},
		{
			name: "project disables pruning",
			userConfig: `
store:
  journal_retention: 168h
`,
			projectConfig: `
store:
  journal_retention: 0
`,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			workDir := t.TempDir()
			mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
			mustMkdir(t, filepath.Join(workDir, ".git"))
			if strings.TrimSpace(tt.userConfig) != "" {
				mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), tt.userConfig)
			}
			if strings.TrimSpace(tt.projectConfig) != "" {
				mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), tt.projectConfig)
			}

			loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
			if err != nil {
				t.Fatalf("expected config to load, got %v", err)
			}
			if loaded.Config.Store.JournalRetention != tt.want {
				t.Fatalf("expected store.journal_retention %s, got %s", tt.want, loaded.Config.Store.JournalRetention)
			}
		})
	}
}

func TestLoadWarnsAndIgnoresDeprecatedConfigKeys(t *testing.T) {
	const warning = "config: resolve.concurrent is deprecated and ignored; use worktree.concurrency\n"

	tests := []struct {
		name            string
		userConfig      string
		projectConfig   string
		wantBatchSize   int
		wantConcurrency int
	}{
		{
			name: "user config exact 0009 corpse key",
			userConfig: `
resolve:
  concurrent: 1
`,
			wantBatchSize:   3,
			wantConcurrency: 2,
		},
		{
			name: "user config mixes deprecated and valid keys",
			userConfig: `
resolve:
  concurrent: 1
  batch_size: 5
worktree:
  concurrency: 4
`,
			wantBatchSize:   5,
			wantConcurrency: 4,
		},
		{
			name: "project config mixes deprecated and valid keys",
			projectConfig: `
resolve:
  concurrent: 1
  batch_size: 6
worktree:
  concurrency: 3
`,
			wantBatchSize:   6,
			wantConcurrency: 3,
		},
		{
			name: "user and project config warn once per load",
			userConfig: `
resolve:
  concurrent: 1
`,
			projectConfig: `
resolve:
  concurrent: 1
  batch_size: 7
`,
			wantBatchSize:   7,
			wantConcurrency: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			workDir := t.TempDir()
			mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
			mustMkdir(t, filepath.Join(workDir, ".git"))
			if strings.TrimSpace(tt.userConfig) != "" {
				mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), tt.userConfig)
			}
			if strings.TrimSpace(tt.projectConfig) != "" {
				mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), tt.projectConfig)
			}
			var stderr bytes.Buffer

			loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir, Stderr: &stderr})

			if err != nil {
				t.Fatalf("expected config to load, got %v", err)
			}
			if stderr.String() != warning {
				t.Fatalf("expected warning %q, got %q", warning, stderr.String())
			}
			if loaded.Config.Resolve.BatchSize != tt.wantBatchSize {
				t.Fatalf("expected resolve.batch_size %d, got %d", tt.wantBatchSize, loaded.Config.Resolve.BatchSize)
			}
			if loaded.Config.Worktree.Concurrency != tt.wantConcurrency {
				t.Fatalf("expected worktree.concurrency %d, got %d", tt.wantConcurrency, loaded.Config.Worktree.Concurrency)
			}
		})
	}
}

func TestLoadRejectsUnknownConfigKeys(t *testing.T) {
	const warning = "config: resolve.concurrent is deprecated and ignored; use worktree.concurrency\n"

	tests := []struct {
		name          string
		userConfig    string
		projectConfig string
		wantStderr    string
	}{
		{
			name: "user config typo",
			userConfig: `
resolve:
  concurent: 1
`,
		},
		{
			name: "project config typo",
			projectConfig: `
resolve:
  concurent: 1
`,
		},
		{
			name: "deprecated key does not bypass typo",
			userConfig: `
resolve:
  concurrent: 1
  concurent: 1
  batch_size: 4
`,
			wantStderr: warning,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			workDir := t.TempDir()
			mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
			mustMkdir(t, filepath.Join(workDir, ".git"))
			if strings.TrimSpace(tt.userConfig) != "" {
				mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), tt.userConfig)
			}
			if strings.TrimSpace(tt.projectConfig) != "" {
				mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), tt.projectConfig)
			}
			var stderr bytes.Buffer

			_, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir, Stderr: &stderr})

			if err == nil {
				t.Fatal("expected config load to fail")
			}
			if !strings.Contains(err.Error(), "resolve.concurent is not a supported config key") {
				t.Fatalf("expected strict resolve typo error, got %q", err.Error())
			}
			if stderr.String() != tt.wantStderr {
				t.Fatalf("expected stderr %q, got %q", tt.wantStderr, stderr.String())
			}
		})
	}
}

func TestLoadRejectsUnknownStoreConfigKey(t *testing.T) {
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
store:
  journal_ttl: 24h
`)

	_, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})

	if err == nil {
		t.Fatal("expected config load to fail")
	}
	if !strings.Contains(err.Error(), "store.journal_ttl is not a supported config key") {
		t.Fatalf("expected strict store key error, got %q", err.Error())
	}
}

func TestLoadRejectsUnknownWorktreeConfigKey(t *testing.T) {
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
worktree:
  bootstrapp: make setup
`)

	_, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})

	if err == nil {
		t.Fatal("expected config load to fail")
	}
	if !strings.Contains(err.Error(), "worktree.bootstrapp is not a supported config key") {
		t.Fatalf("expected strict worktree key error, got %q", err.Error())
	}
}

func TestLoadRejectsInvalidConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   string
		contains string
	}{
		{
			name:     "invalid YAML",
			config:   "defaults:\n  agent: [",
			contains: "parse config",
		},
		{
			name: "invalid semantic value",
			config: `
watch:
  max_rounds: 0
`,
			contains: "watch.max_rounds must be greater than 0",
		},
		{
			name: "invalid duration",
			config: `
watch:
  poll_interval: soon
`,
			contains: "invalid duration",
		},
		{
			name: "negative journal retention",
			config: `
store:
  journal_retention: -1h
`,
			contains: "store.journal_retention must be greater than or equal to 0",
		},
		{
			name: "invalid implement auto push",
			config: `
implement:
  auto_push: sometimes
`,
			contains: "implement.auto_push must be boolean",
		},
		{
			name: "absolute worktree copy",
			config: `
worktree:
  copy:
    - /tmp/secret.env
`,
			contains: "worktree.copy",
		},
		{
			name: "dot-dot worktree copy",
			config: `
worktree:
  copy:
    - ../secret.env
`,
			contains: "worktree.copy",
		},
		{
			name: "invalid worktree concurrency",
			config: `
worktree:
  concurrency: 0
`,
			contains: "worktree.concurrency must be greater than 0",
		},
		{
			name: "invalid worktree bootstrap timeout",
			config: `
worktree:
  bootstrap_timeout: 0
`,
			contains: "worktree.bootstrap_timeout must be greater than 0",
		},
		{
			name: "relative worktree location",
			config: `
worktree:
  location: relative-worktrees
`,
			contains: "worktree.location must be absolute after ~ expansion",
		},
		{
			name: "unknown resolve key",
			config: `
resolve:
  concurent: 2
`,
			contains: "resolve.concurent is not a supported config key",
		},
		{
			name: "unknown specs key",
			config: `
specs:
  directory: docs/specs
`,
			contains: "specs.directory is not a supported config key",
		},
		{
			name: "empty specs root",
			config: `
specs:
  root: ""
`,
			contains: "specs.root must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			workDir := t.TempDir()
			mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
			mustMkdir(t, filepath.Join(workDir, ".git"))
			mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), tt.config)

			_, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})

			if err == nil {
				t.Fatal("expected config load to fail")
			}
			if !strings.Contains(err.Error(), tt.contains) {
				t.Fatalf("expected error to contain %q, got %q", tt.contains, err.Error())
			}
		})
	}
}

func TestLoadRejectsWorktreeLocationInsideRepository(t *testing.T) {
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
worktree:
  location: `+filepath.Join(workDir, ".roundfix", "worktrees")+`
`)

	_, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})

	if err == nil {
		t.Fatal("expected in-repository worktree.location to fail")
	}
	if !strings.Contains(err.Error(), "worktree.location must not be inside the repository tree") {
		t.Fatalf("expected in-repository location error, got %q", err.Error())
	}
}

func TestLoadExpandsWorktreeLocationHome(t *testing.T) {
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
worktree:
  location: ~/roundfix-worktrees
`)

	loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})

	if err != nil {
		t.Fatalf("expected home-expanded worktree.location to load, got %v", err)
	}
	if want := filepath.Join(homeDir, "roundfix-worktrees"); loaded.Config.Worktree.Location != want {
		t.Fatalf("expected worktree.location %q, got %q", want, loaded.Config.Worktree.Location)
	}
}

func TestInitCreatesUserConfig(t *testing.T) {
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(workDir, ".git"))

	result, err := Init(context.Background(), InitOptions{
		Scope:   InitScopeUser,
		HomeDir: homeDir,
		WorkDir: workDir,
	})
	if err != nil {
		t.Fatalf("expected init to create User Config, got %v", err)
	}
	expectedPath := filepath.Join(homeDir, ".roundfix", "config.yml")
	if result.Scope != InitScopeUser || result.Path != expectedPath {
		t.Fatalf("expected user result at %q, got %#v", expectedPath, result)
	}
	content := mustRead(t, expectedPath)
	if !strings.Contains(content, "agent: codex") || !strings.Contains(content, "agent_full_access: false") ||
		!strings.Contains(content, `artifact_dir: ""`) || !strings.Contains(content, "Roundfix Home artifacts/<repo-id>") ||
		!strings.Contains(content, "specs:") || !strings.Contains(content, `root: "docs/specs"`) ||
		!strings.Contains(content, "worktree:") || !strings.Contains(content, `location: "~/.roundfix/worktrees"`) ||
		!strings.Contains(content, "concurrency: 2") || !strings.Contains(content, "copy: []") ||
		!strings.Contains(content, "store:") || !strings.Contains(content, "journal_retention: 336h") ||
		!strings.Contains(content, "implement:") || !strings.Contains(content, "auto_push: false") ||
		!strings.Contains(content, "max_run_duration: 2h") {
		t.Fatalf("expected default config content, got %s", content)
	}
	if strings.Contains(content, "resolve.concurrent") || strings.Contains(content, "  concurrent:") {
		t.Fatalf("expected generated config to omit resolve.concurrent, got %s", content)
	}
	if _, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir}); err != nil {
		t.Fatalf("expected generated User Config to load, got %v", err)
	}
}

func TestInitCreatesProjectConfig(t *testing.T) {
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(workDir, ".git"))

	result, err := Init(context.Background(), InitOptions{
		Scope:   InitScopeProject,
		HomeDir: homeDir,
		WorkDir: workDir,
	})
	if err != nil {
		t.Fatalf("expected init to create Project Config, got %v", err)
	}
	expectedPath := filepath.Join(workDir, ".roundfixrc.yml")
	if result.Scope != InitScopeProject || result.Path != expectedPath {
		t.Fatalf("expected project result at %q, got %#v", expectedPath, result)
	}
	if _, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir}); err != nil {
		t.Fatalf("expected generated Project Config to load, got %v", err)
	}
}

func TestInitRejectsProjectScopeOutsideGitRoot(t *testing.T) {
	_, err := Init(context.Background(), InitOptions{
		Scope:   InitScopeProject,
		HomeDir: t.TempDir(),
		WorkDir: t.TempDir(),
	})
	if err == nil {
		t.Fatal("expected project init outside Git root to fail")
	}
	if !strings.Contains(err.Error(), "requires a Git root") {
		t.Fatalf("expected Git root guidance, got %q", err.Error())
	}
}

func TestInitDoesNotOverwriteWithoutForce(t *testing.T) {
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(workDir, ".git"))
	path := filepath.Join(workDir, ".roundfixrc.yml")
	mustWrite(t, path, "defaults:\n  agent: claude\n")

	_, err := Init(context.Background(), InitOptions{
		Scope:   InitScopeProject,
		HomeDir: homeDir,
		WorkDir: workDir,
	})
	if err == nil {
		t.Fatal("expected existing config to fail without force")
	}
	if !strings.Contains(err.Error(), "--force") {
		t.Fatalf("expected force guidance, got %q", err.Error())
	}
	if content := mustRead(t, path); !strings.Contains(content, "agent: claude") {
		t.Fatalf("expected existing config to remain, got %s", content)
	}
}

func TestInitForceOverwritesExistingConfig(t *testing.T) {
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(workDir, ".git"))
	path := filepath.Join(workDir, ".roundfixrc.yml")
	mustWrite(t, path, "defaults:\n  agent: claude\n")

	result, err := Init(context.Background(), InitOptions{
		Scope:   InitScopeProject,
		HomeDir: homeDir,
		WorkDir: workDir,
		Force:   true,
	})
	if err != nil {
		t.Fatalf("expected force init to overwrite config, got %v", err)
	}
	if !result.Overwritten {
		t.Fatalf("expected overwritten result, got %#v", result)
	}
	if content := mustRead(t, path); !strings.Contains(content, "agent: codex") || strings.Contains(content, "agent: claude") {
		t.Fatalf("expected default config to replace old content, got %s", content)
	}
}

func TestValidateArtifactDirectoryResolvesAndCreatesPaths(t *testing.T) {
	homeDir := t.TempDir()
	gitRoot := t.TempDir()

	defaultPath, err := ValidateArtifactDirectory("", gitRoot, homeDir)
	if err != nil {
		t.Fatalf("expected default artifact dir to validate, got %v", err)
	}
	if defaultPath != filepath.Join(homeDir, ".roundfix", "artifacts", repoID(gitRoot)) {
		t.Fatalf("expected default artifact dir under Roundfix Home, got %q", defaultPath)
	}
	assertDir(t, defaultPath)

	relativePath, err := ValidateArtifactDirectory("reviews", gitRoot, homeDir)
	if err != nil {
		t.Fatalf("expected relative artifact dir to validate, got %v", err)
	}
	if relativePath != filepath.Join(gitRoot, "reviews") {
		t.Fatalf("expected relative artifact dir under git root, got %q", relativePath)
	}
	assertDir(t, relativePath)

	homePath, err := ValidateArtifactDirectory("~/roundfix-artifacts", gitRoot, homeDir)
	if err != nil {
		t.Fatalf("expected home artifact dir to validate, got %v", err)
	}
	if homePath != filepath.Join(homeDir, "roundfix-artifacts") {
		t.Fatalf("expected home artifact dir expansion, got %q", homePath)
	}
	assertDir(t, homePath)

	absoluteConfigPath := filepath.Join(t.TempDir(), "configured-artifacts")
	absolutePath, err := ValidateArtifactDirectory(absoluteConfigPath, gitRoot, homeDir)
	if err != nil {
		t.Fatalf("expected absolute artifact dir to validate, got %v", err)
	}
	if absolutePath != absoluteConfigPath {
		t.Fatalf("expected absolute artifact dir unchanged, got %q", absolutePath)
	}
	assertDir(t, absolutePath)
}

func TestValidateArtifactDirectoryRejectsInvalidPaths(t *testing.T) {
	homeDir := t.TempDir()
	gitRoot := t.TempDir()

	if _, err := ValidateArtifactDirectory("", "", homeDir); err == nil {
		t.Fatal("expected empty artifact dir without git root to fail")
	}
	if _, err := ValidateArtifactDirectory("reviews", "", homeDir); err == nil {
		t.Fatal("expected relative artifact dir without git root to fail")
	}

	filePath := filepath.Join(gitRoot, "artifact-file")
	mustWrite(t, filePath, "not a directory")
	_, err := ValidateArtifactDirectory(filePath, gitRoot, homeDir)
	if err == nil {
		t.Fatal("expected file artifact dir to fail")
	}
	if !strings.Contains(err.Error(), "is not a directory") {
		t.Fatalf("expected not-directory error, got %q", err.Error())
	}
}

func TestResolveSpecsRoot(t *testing.T) {
	repoRoot := t.TempDir()
	mustMkdir(t, filepath.Join(repoRoot, "docs", "specs"))
	mustMkdir(t, filepath.Join(repoRoot, "configured-specs"))
	absoluteInternal := filepath.Join(repoRoot, "absolute-specs")
	mustMkdir(t, absoluteInternal)
	externalRoot := filepath.Join(t.TempDir(), "external-specs")
	mustMkdir(t, externalRoot)

	tests := []struct {
		name     string
		root     string
		wantPath string
		wantExt  bool
	}{
		{
			name:     "builtin default is internal docs specs",
			root:     Builtin().Specs.Root,
			wantPath: filepath.Join(repoRoot, "docs", "specs"),
		},
		{
			name:     "relative root resolves against repository root",
			root:     "configured-specs",
			wantPath: filepath.Join(repoRoot, "configured-specs"),
		},
		{
			name:     "absolute root is used as is",
			root:     absoluteInternal,
			wantPath: absoluteInternal,
		},
		{
			name:     "absolute external root is classified external",
			root:     externalRoot,
			wantPath: externalRoot,
			wantExt:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loaded := Loaded{Config: Builtin(), GitRoot: repoRoot}
			loaded.Config.Specs.Root = tt.root

			got, err := ResolveSpecsRoot(loaded, repoRoot)
			if err != nil {
				t.Fatalf("ResolveSpecsRoot() error = %v", err)
			}
			if got.Path != tt.wantPath {
				t.Fatalf("ResolveSpecsRoot().Path = %q, want %q", got.Path, tt.wantPath)
			}
			if got.External != tt.wantExt {
				t.Fatalf("ResolveSpecsRoot().External = %t, want %t", got.External, tt.wantExt)
			}
		})
	}
}

func TestResolveSpecsRootUsesLoadedDefault(t *testing.T) {
	homeDir := t.TempDir()
	repoRoot := t.TempDir()
	mustMkdir(t, filepath.Join(repoRoot, ".git"))
	mustMkdir(t, filepath.Join(repoRoot, "docs", "specs"))

	loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: repoRoot})
	if err != nil {
		t.Fatalf("expected config to load, got %v", err)
	}
	got, err := ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		t.Fatalf("ResolveSpecsRoot() error = %v", err)
	}
	wantPath := filepath.Join(repoRoot, "docs", "specs")
	if got.Path != wantPath {
		t.Fatalf("ResolveSpecsRoot().Path = %q, want %q", got.Path, wantPath)
	}
	if got.External {
		t.Fatal("expected default Spec Root to be internal")
	}
}

func TestResolveSpecsRootClassifiesExternalSymlink(t *testing.T) {
	repoRoot := t.TempDir()
	externalRoot := filepath.Join(t.TempDir(), "external-specs")
	mustMkdir(t, filepath.Join(repoRoot, "docs"))
	mustMkdir(t, externalRoot)
	symlinkPath := filepath.Join(repoRoot, "docs", "specs")
	if err := os.Symlink(externalRoot, symlinkPath); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	loaded := Loaded{Config: Builtin(), GitRoot: repoRoot}

	got, err := ResolveSpecsRoot(loaded, repoRoot)
	if err != nil {
		t.Fatalf("ResolveSpecsRoot() error = %v", err)
	}
	if got.Path != symlinkPath {
		t.Fatalf("ResolveSpecsRoot().Path = %q, want %q", got.Path, symlinkPath)
	}
	if !got.External {
		t.Fatal("expected symlinked external Spec Root to be classified external")
	}
}

func TestResolveSpecsRootRejectsInvalidRoots(t *testing.T) {
	repoRoot := t.TempDir()
	filePath := filepath.Join(repoRoot, "spec-file")
	mustWrite(t, filePath, "not a directory")

	tests := []struct {
		name     string
		root     string
		resolved string
		contains string
	}{
		{
			name:     "missing relative root",
			root:     "missing-specs",
			resolved: filepath.Join(repoRoot, "missing-specs"),
			contains: "does not exist",
		},
		{
			name:     "file root",
			root:     filePath,
			resolved: filePath,
			contains: "is not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loaded := Loaded{Config: Builtin(), GitRoot: repoRoot}
			loaded.Config.Specs.Root = tt.root

			_, err := ResolveSpecsRoot(loaded, repoRoot)
			if err == nil {
				t.Fatal("expected ResolveSpecsRoot to fail")
			}
			if !strings.Contains(err.Error(), tt.contains) {
				t.Fatalf("expected error to contain %q, got %q", tt.contains, err.Error())
			}
			if !strings.Contains(err.Error(), tt.resolved) {
				t.Fatalf("expected error to name resolved path %q, got %q", tt.resolved, err.Error())
			}
		})
	}
}

func TestResolveReviewRoot(t *testing.T) {
	repoRoot := t.TempDir()
	specSlug := "0001-widget-flow"
	mustMkdir(t, filepath.Join(repoRoot, "docs", "specs", specSlug))
	externalSpecsRoot := filepath.Join(t.TempDir(), "external-specs")
	mustMkdir(t, filepath.Join(externalSpecsRoot, specSlug))
	explicitDir := filepath.Join(t.TempDir(), "artifacts")

	tests := []struct {
		name string
		ctx  ReviewArtifactContext
		want string
	}{
		{
			name: "explicit artifact directory keeps existing layout",
			ctx: ReviewArtifactContext{
				ExplicitArtifactDir: explicitDir,
				RepoRoot:            repoRoot,
				SpecSlug:            specSlug,
				PRNumber:            123,
			},
			want: filepath.Join(explicitDir, "reviews", "pr-123"),
		},
		{
			name: "existing spec stores rounds under spec reviews",
			ctx: ReviewArtifactContext{
				RepoRoot: repoRoot,
				SpecSlug: specSlug,
				PRNumber: 123,
			},
			want: filepath.Join(repoRoot, "docs", "specs", specSlug, "reviews"),
		},
		{
			name: "spec-less stores rounds under in-repo review root",
			ctx: ReviewArtifactContext{
				RepoRoot: repoRoot,
				PRNumber: 123,
			},
			want: filepath.Join(repoRoot, "docs", "specs", "_reviews", "pr-123"),
		},
		{
			name: "existing spec under external root stores rounds under external spec reviews",
			ctx: ReviewArtifactContext{
				RepoRoot:  repoRoot,
				SpecsRoot: externalSpecsRoot,
				SpecSlug:  specSlug,
				PRNumber:  123,
			},
			want: filepath.Join(externalSpecsRoot, specSlug, "reviews"),
		},
		{
			name: "unknown spec falls back to spec-less root",
			ctx: ReviewArtifactContext{
				RepoRoot: repoRoot,
				SpecSlug: "9999-missing",
				PRNumber: 123,
			},
			want: filepath.Join(repoRoot, "docs", "specs", "_reviews", "pr-123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveReviewRoot(tt.ctx)
			if err != nil {
				t.Fatalf("ResolveReviewRoot() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("ResolveReviewRoot() = %q, want %q", got, tt.want)
			}
		})
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustWrite(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func assertDir(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %s to be a directory", path)
	}
}
