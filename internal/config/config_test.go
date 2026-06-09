package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadAppliesConfigPrecedence(t *testing.T) {
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
defaults:
  agent: claude
  artifact_dir: user-artifacts
watch:
  max_rounds: 4
  poll_interval: 10s
resolve:
  batch_size: 2
`)
	mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), `
defaults:
  agent: opencode
watch:
  max_rounds: 8
budget:
  max_run_duration: 3h
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
	if loaded.Config.Watch.MaxRounds != 8 {
		t.Fatalf("expected project max rounds, got %d", loaded.Config.Watch.MaxRounds)
	}
	if loaded.Config.Watch.PollInterval != 10*time.Second {
		t.Fatalf("expected user poll interval, got %s", loaded.Config.Watch.PollInterval)
	}
	if loaded.Config.Budget.MaxRunDuration != 3*time.Hour {
		t.Fatalf("expected project max run duration, got %s", loaded.Config.Budget.MaxRunDuration)
	}
	if loaded.Config.Resolve.BatchSize != 2 {
		t.Fatalf("expected user batch size, got %d", loaded.Config.Resolve.BatchSize)
	}
	if loaded.Config.Resolve.Concurrent != 1 {
		t.Fatalf("expected built-in concurrent default, got %d", loaded.Config.Resolve.Concurrent)
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

func TestValidateArtifactDirectoryResolvesAndCreatesPaths(t *testing.T) {
	homeDir := t.TempDir()
	gitRoot := t.TempDir()

	defaultPath, err := ValidateArtifactDirectory("", gitRoot, homeDir)
	if err != nil {
		t.Fatalf("expected default artifact dir to validate, got %v", err)
	}
	if defaultPath != filepath.Join(gitRoot, ".roundfix") {
		t.Fatalf("expected default artifact dir under git root, got %q", defaultPath)
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
