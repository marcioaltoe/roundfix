// Suite: layered Roundfix configuration
// Invariant: built-in, User Config, and Project Config values resolve with strict validation.
// Boundary IN: config defaults, YAML decoding, precedence, validation, and generated config text.
// Boundary OUT: operational Preflight Validation, owned by internal/preflight/preflight_test.go.
package config

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/agent"
)

var updateProfilesConfigGoldens = flag.Bool(
	"update-profiles-config-goldens",
	false,
	"re-record profiles config writer characterization goldens",
)

func TestLoadAppliesConfigPrecedence(t *testing.T) {
	t.Parallel()
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
  check_grace_period: 90s
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
runtimes:
  opencode:
    model: project-opencode
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
	if loaded.Config.Watch.CheckGracePeriod != 90*time.Second {
		t.Fatalf("expected user check grace period, got %s", loaded.Config.Watch.CheckGracePeriod)
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

func TestBuiltinRuntimeDefaults(t *testing.T) {
	t.Parallel()
	config := Builtin()

	if config.Runtimes.Codex.Model != "gpt-5.5" || config.Runtimes.Codex.ReasoningEffort != "xhigh" {
		t.Fatalf("expected built-in Codex gpt-5.5/xhigh, got %#v", config.Runtimes.Codex)
	}
	if config.Runtimes.Claude.Model != "opus" || config.Runtimes.Claude.ReasoningEffort != "" {
		t.Fatalf("expected built-in Claude opus with model-managed reasoning, got %#v", config.Runtimes.Claude)
	}
	if config.Runtimes.OpenCode.Model != "" || config.Runtimes.OpenCode.ReasoningEffort != "" {
		t.Fatalf("expected built-in OpenCode to require explicit selection, got %#v", config.Runtimes.OpenCode)
	}
}

func TestBuiltinProfilesGeneratedCodexPolicy(t *testing.T) {
	t.Parallel()
	config := Builtin()
	wantPreferred := selectionForTest("codex", "gpt-5.6-sol", "high")
	wantFallback := selectionForTest("codex", "gpt-5.5", "xhigh")
	for _, category := range []WorkCategory{CategoryGeneral, CategoryBackend, CategoryQA, CategoryReview} {
		resolved, err := ResolveProfile(config, category, nil)
		if err != nil {
			t.Fatalf("ResolveProfile(%q) error = %v", category, err)
		}
		if resolved.Profile.Preferred != wantPreferred {
			t.Fatalf("%s preferred = %#v, want %#v", category, resolved.Profile.Preferred, wantPreferred)
		}
		if len(resolved.Profile.Fallbacks) != 1 || resolved.Profile.Fallbacks[0] != wantFallback {
			t.Fatalf("%s fallbacks = %#v, want [%#v]", category, resolved.Profile.Fallbacks, wantFallback)
		}
	}
}

func TestDefaultConfigYAMLGeneratedCodexPolicy(t *testing.T) {
	t.Parallel()
	content := DefaultConfigYAML()
	if got := strings.Count(content, "model: gpt-5.6-sol"); got != 5 {
		t.Fatalf("Sol occurrence count = %d, want 5:\n%s", got, content)
	}
	if got := strings.Count(content, "model: gpt-5.5"); got != 4 {
		t.Fatalf("GPT-5.5 occurrence count = %d, want 4:\n%s", got, content)
	}
	for _, forbidden := range []string{"model: gpt-5.6-terra", "model: gpt-5.6-luna", "reasoning_effort: max"} {
		if strings.Contains(content, forbidden) {
			t.Fatalf("generated config contains operational default %q:\n%s", forbidden, content)
		}
	}
}

func TestDefaultConfigYAMLVerificationCapacity(t *testing.T) {
	t.Parallel()
	content := DefaultConfigYAML()
	for _, want := range []string{
		"# Verification command for review Batches; Spec Tasks use their task file commands.\n  verification: make verify",
		"verification:\n  # Maximum concurrent Task Verification attempts per spec Run; independent from worktree.concurrency.\n  concurrency: 1",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("generated config is missing %q:\n%s", want, content)
		}
	}
	for _, scope := range []string{InitScopeUser, InitScopeProject} {
		t.Run(scope, func(t *testing.T) {
			homeDir := t.TempDir()
			workDir := t.TempDir()
			mustMkdir(t, filepath.Join(workDir, ".git"))

			result, err := Init(context.Background(), InitOptions{
				Scope:   scope,
				HomeDir: homeDir,
				WorkDir: workDir,
			})
			if err != nil {
				t.Fatalf("generate %s config: %v", scope, err)
			}
			if got := mustRead(t, result.Path); got != content {
				t.Fatalf("generated %s config differs from DefaultConfigYAML", scope)
			}
		})
	}
}

func TestModelCatalogRetainsOfficialCodexIdentifiers(t *testing.T) {
	t.Parallel()
	catalog := agent.ModelCatalog("codex")
	for _, identifier := range []string{"gpt-5.6-sol", "gpt-5.6-terra", "gpt-5.6-luna"} {
		found := false
		for _, choice := range catalog {
			if choice.Value == identifier {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Codex Model Catalog is missing official identifier %q: %#v", identifier, catalog)
		}
	}
}

func TestAgentSelectionProfileBuiltinsResolveRequiredCategories(t *testing.T) {
	t.Parallel()
	config := Builtin()
	tests := []struct {
		name     string
		category WorkCategory
		want     AgentSelectionProfile
	}{
		{
			name:     "general",
			category: CategoryGeneral,
			want: profileForTest(
				selectionForTest("codex", "gpt-5.6-sol", "high"),
				selectionForTest("codex", "gpt-5.5", "xhigh"),
			),
		},
		{
			name:     "backend",
			category: CategoryBackend,
			want: profileForTest(
				selectionForTest("codex", "gpt-5.6-sol", "high"),
				selectionForTest("codex", "gpt-5.5", "xhigh"),
			),
		},
		{
			name:     "frontend",
			category: CategoryFrontend,
			want: profileForTest(
				selectionForTest("claude", "opus", "xhigh"),
				selectionForTest("codex", "gpt-5.6-sol", "high"),
			),
		},
		{
			name:     "qa",
			category: CategoryQA,
			want: profileForTest(
				selectionForTest("codex", "gpt-5.6-sol", "high"),
				selectionForTest("codex", "gpt-5.5", "xhigh"),
			),
		},
		{
			name:     "review",
			category: CategoryReview,
			want: profileForTest(
				selectionForTest("codex", "gpt-5.6-sol", "high"),
				selectionForTest("codex", "gpt-5.5", "xhigh"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveProfile(config, tt.category, nil)
			if err != nil {
				t.Fatalf("ResolveProfile(%q) error = %v", tt.category, err)
			}
			if got.Source != ProfileSourceBuiltIn {
				t.Fatalf("expected built-in source, got %q", got.Source)
			}
			if got.InheritedFrom != "" {
				t.Fatalf("expected no inheritance, got %q", got.InheritedFrom)
			}
			if !profilesEqual(got.Profile, tt.want) {
				t.Fatalf("ResolveProfile(%q) mismatch\nwant: %#v\ngot:  %#v", tt.category, tt.want, got.Profile)
			}
		})
	}
}

func TestAgentSelectionProfileOptionalCategoryInheritsGeneral(t *testing.T) {
	t.Parallel()
	config := Builtin()
	general, err := ResolveProfile(config, CategoryGeneral, nil)
	if err != nil {
		t.Fatalf("ResolveProfile(general) error = %v", err)
	}

	got, err := ResolveProfile(config, CategoryData, nil)
	if err != nil {
		t.Fatalf("ResolveProfile(data) error = %v", err)
	}

	if got.InheritedFrom != CategoryGeneral {
		t.Fatalf("expected data to inherit general, got %q", got.InheritedFrom)
	}
	if got.Source != general.Source {
		t.Fatalf("expected inherited source %q, got %q", general.Source, got.Source)
	}
	if !profilesEqual(got.Profile, general.Profile) {
		t.Fatalf("expected data profile to equal general profile\nwant: %#v\ngot:  %#v", general.Profile, got.Profile)
	}
	if _, stored := config.Profiles[CategoryData]; stored {
		t.Fatalf("expected optional data profile not to be stored in built-ins")
	}
}

func TestProfileResolverUsesAtomicProjectOverUserPrecedence(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
profiles:
  backend:
    preferred:
      runtime: claude
      model: user-preferred
      reasoning_effort: medium
    fallbacks:
      - runtime: codex
        model: user-fallback
        reasoning_effort: high
`)
	mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), `
profiles:
  backend:
    preferred:
      runtime: codex
      model: project-preferred
      reasoning_effort: high
    fallbacks:
      - runtime: claude
        model: project-fallback
        reasoning_effort: ""
`)

	loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
	if err != nil {
		t.Fatalf("expected config to load, got %v", err)
	}
	got, err := ResolveProfile(loaded.Config, CategoryBackend, nil)
	if err != nil {
		t.Fatalf("ResolveProfile(backend) error = %v", err)
	}

	want := profileForTest(
		selectionForTest("codex", "project-preferred", "high"),
		selectionForTest("claude", "project-fallback", ""),
	)
	if got.Source != ProfileSourceProject {
		t.Fatalf("expected project source, got %q", got.Source)
	}
	if !profilesEqual(got.Profile, want) {
		t.Fatalf("expected project profile to replace user profile atomically\nwant: %#v\ngot:  %#v", want, got.Profile)
	}
}

func TestAgentSelectionProfileRejectsInvalidConfiguredProfiles(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		config   string
		contains string
	}{
		{
			name: "empty profiles section",
			config: `
profiles: {}
`,
			contains: "profiles must define at least one Agent Selection Profile",
		},
		{
			name: "missing fallback chain",
			config: `
profiles:
  backend:
    preferred:
      runtime: codex
      model: gpt-custom
      reasoning_effort: high
`,
			contains: "profiles.backend.fallbacks is required",
		},
		{
			name: "partial preferred selection",
			config: `
profiles:
  backend:
    preferred:
      runtime: codex
      model: gpt-custom
    fallbacks:
      - runtime: codex
        model: fallback
        reasoning_effort: high
`,
			contains: "profiles.backend.preferred.reasoning_effort is required",
		},
		{
			name: "empty fallback chain",
			config: `
profiles:
  backend:
    preferred:
      runtime: codex
      model: gpt-custom
      reasoning_effort: high
    fallbacks: []
`,
			contains: "profiles.backend.fallbacks must include at least one Agent Selection",
		},
		{
			name: "duplicate tuple",
			config: `
profiles:
  backend:
    preferred:
      runtime: codex
      model: gpt-custom
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-custom
        reasoning_effort: high
`,
			contains: `profiles.backend contains duplicate Agent Selection "codex / gpt-custom / high"`,
		},
		{
			name: "unknown runtime",
			config: `
profiles:
  backend:
    preferred:
      runtime: local
      model: llama
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: fallback
        reasoning_effort: high
`,
			contains: `profiles.backend.preferred.runtime "local" is invalid; supported values: codex, claude, opencode`,
		},
		{
			name: "empty model",
			config: `
profiles:
  backend:
    preferred:
      runtime: codex
      model: ""
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: fallback
        reasoning_effort: high
`,
			contains: "profiles.backend.preferred.model must not be empty",
		},
		{
			name: "mixed legacy and profiles",
			config: `
defaults:
  agent: codex
profiles:
  backend:
    preferred:
      runtime: codex
      model: gpt-custom
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: fallback
        reasoning_effort: high
`,
			contains: "config mixes legacy runtime defaults with profiles; migrate with `roundfix profiles configure --scope user`",
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
				t.Fatalf("expected error containing %q, got %q", tt.contains, err.Error())
			}
		})
	}
}

func TestAgentSelectionProfileDistinguishesMissingReasoningFromEmptyReasoning(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
profiles:
  backend:
    preferred:
      runtime: claude
      model: custom-fable
      reasoning_effort: ""
    fallbacks:
      - runtime: codex
        model: custom-sol
        reasoning_effort: ""
`)

	loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
	if err != nil {
		t.Fatalf("expected explicit empty reasoning to load, got %v", err)
	}
	got, err := ResolveProfile(loaded.Config, CategoryBackend, nil)
	if err != nil {
		t.Fatalf("ResolveProfile(backend) error = %v", err)
	}
	if got.Profile.Preferred.ReasoningEffort != "" || got.Profile.Fallbacks[0].ReasoningEffort != "" {
		t.Fatalf("expected explicit empty reasoning to survive unchanged, got %#v", got.Profile)
	}
	if got.Profile.Preferred.Model != "custom-fable" || got.Profile.Fallbacks[0].Model != "custom-sol" {
		t.Fatalf("expected custom model strings to survive unchanged, got %#v", got.Profile)
	}
}

func TestProfileLegacyMigrationConvertsRuntimeDefaults(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
defaults:
  agent: claude
runtimes:
  claude:
    model: fable
    reasoning_effort: ""
`)

	loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
	if err != nil {
		t.Fatalf("expected legacy config to load, got %v", err)
	}
	got, err := ResolveProfile(loaded.Config, CategoryBackend, nil)
	if err != nil {
		t.Fatalf("ResolveProfile(backend) error = %v", err)
	}

	want := profileForTest(
		selectionForTest("claude", "fable", ""),
		selectionForTest("codex", "gpt-5.5", "xhigh"),
	)
	if got.Source != ProfileSourceUser {
		t.Fatalf("expected user source, got %q", got.Source)
	}
	if !profilesEqual(got.Profile, want) {
		t.Fatalf("expected legacy runtime defaults to convert into profile\nwant: %#v\ngot:  %#v", want, got.Profile)
	}
}

func TestProfileLegacyDefaultCodexKeepsDistinctBuiltInFallback(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), `
runtimes:
  claude:
    reasoning_effort: high
`)

	loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
	if err != nil {
		t.Fatalf("load partial legacy runtime config: %v", err)
	}
	resolved, err := ResolveProfile(loaded.Config, CategoryGeneral, nil)
	if err != nil {
		t.Fatalf("ResolveProfile(general): %v", err)
	}
	want := profileForTest(
		selectionForTest("codex", "gpt-5.5", "xhigh"),
		selectionForTest("codex", "gpt-5.6-sol", "high"),
	)
	if !profilesEqual(resolved.Profile, want) {
		t.Fatalf("legacy default Codex profile mismatch\nwant: %#v\n got: %#v", want, resolved.Profile)
	}
}

func TestProfilesConfigWriterCharacterization(t *testing.T) {
	t.Parallel()
	const corpusDir = "testdata/profiles_config_writer"
	tests := []struct {
		name     string
		file     string
		removals []WorkCategory
	}{
		{name: "five categories with Fallback Chains", file: "five_categories_with_fallback_chains"},
		{name: "comments above and beside entries", file: "comments_above_and_beside_entries"},
		{name: "non-default indentation", file: "non_default_indentation"},
		{name: "non-alphabetical key order", file: "non_alphabetical_key_order"},
		{name: "multiline scalar", file: "multiline_scalar"},
		{name: "YAML anchor and alias", file: "yaml_anchor_and_alias"},
		{name: "non-ASCII values", file: "non_ascii_values"},
		{name: "empty file", file: "empty_file"},
		{name: "no profiles section", file: "no_profiles_section"},
		{
			name:     "removal before top-level separator",
			file:     "removal_before_top_level_separator",
			removals: []WorkCategory{CategoryFrontend},
		},
		{
			name:     "removal with adjacent trivia",
			file:     "removal_with_adjacent_trivia",
			removals: []WorkCategory{CategoryBackend},
		},
	}
	configuredProfiles := Profiles{
		CategoryBackend: {
			Profile: profileForTest(
				selectionForTest("codex", "characterization-backend", "high"),
				selectionForTest("claude", "characterization-fallback", "xhigh"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profiles := configuredProfiles
			if len(tt.removals) > 0 {
				profiles = nil
			}
			inputPath := filepath.Join(corpusDir, tt.file+".input.yml")
			input, err := os.ReadFile(inputPath)
			if err != nil {
				t.Fatalf("read characterization config %q: %v", tt.file, err)
			}

			homeDir := t.TempDir()
			workDir := t.TempDir()
			mustMkdir(t, filepath.Join(workDir, ".git"))
			configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
			mustMkdir(t, filepath.Dir(configPath))
			if err := os.WriteFile(configPath, input, 0o644); err != nil {
				t.Fatalf("write characterization config %q: %v", tt.file, err)
			}

			if _, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
				Scope:    InitScopeUser,
				HomeDir:  homeDir,
				WorkDir:  workDir,
				Profiles: profiles,
				Removals: tt.removals,
			}); err != nil {
				t.Fatalf("write characterization config %q through current writer: %v", tt.file, err)
			}
			actual, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatalf("read written characterization config %q: %v", tt.file, err)
			}
			if _, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir}); err != nil {
				t.Fatalf("load written characterization config %q: %v", tt.file, err)
			}

			goldenPath := filepath.Join(corpusDir, tt.file+".golden.yml")
			if *updateProfilesConfigGoldens {
				if err := os.WriteFile(goldenPath, actual, 0o644); err != nil {
					t.Fatalf("update characterization golden %q: %v", tt.file, err)
				}
				return
			}
			golden, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("read characterization golden %q: %v (re-record deliberately with -update-profiles-config-goldens)", tt.file, err)
			}
			if !bytes.Equal(actual, golden) {
				t.Fatalf("characterization config %q changed:\n%s", tt.file, profilesConfigGoldenDiff(golden, actual))
			}
		})
	}
}

func TestProfilesConfigureRemovalPreservesSpacing(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		original string
		want     string
		removal  WorkCategory
	}{
		{
			name: "keeps separator before following top-level key",
			original: `profiles:
  backend:
    preferred: {runtime: codex, model: backend-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: backend-fallback, reasoning_effort: xhigh}]
  frontend:
    preferred: {runtime: claude, model: frontend-primary, reasoning_effort: xhigh}
    fallbacks: [{runtime: codex, model: frontend-fallback, reasoning_effort: high}]
  qa:
    preferred: {runtime: codex, model: qa-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: qa-fallback, reasoning_effort: xhigh}]

watch:
  max_rounds: 4
`,
			want: `profiles:
  backend:
    preferred: {runtime: codex, model: backend-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: backend-fallback, reasoning_effort: xhigh}]
  qa:
    preferred: {runtime: codex, model: qa-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: qa-fallback, reasoning_effort: xhigh}]

watch:
  max_rounds: 4
`,
			removal: CategoryFrontend,
		},
		{
			name: "removes only middle category trivia",
			original: `profiles:
  # general survives with its comment
  general:
    preferred: {runtime: codex, model: general-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: general-fallback, reasoning_effort: xhigh}]
  # backend comment belongs to the removed category
  backend:
    preferred: {runtime: codex, model: backend-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: backend-fallback, reasoning_effort: xhigh}]

  # qa survives with its comment
  qa:
    preferred: {runtime: codex, model: qa-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: qa-fallback, reasoning_effort: xhigh}]
`,
			want: `profiles:
  # general survives with its comment
  general:
    preferred: {runtime: codex, model: general-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: general-fallback, reasoning_effort: xhigh}]
  # qa survives with its comment
  qa:
    preferred: {runtime: codex, model: qa-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: qa-fallback, reasoning_effort: xhigh}]
`,
			removal: CategoryBackend,
		},
		{
			name: "keeps surviving category trailing blank line",
			original: `profiles:
  general:
    preferred: {runtime: codex, model: general-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: general-fallback, reasoning_effort: xhigh}]

  # backend comment belongs to the removed category
  backend:
    preferred: {runtime: codex, model: backend-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: backend-fallback, reasoning_effort: xhigh}]

  # qa survives with its comment
  qa:
    preferred: {runtime: codex, model: qa-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: qa-fallback, reasoning_effort: xhigh}]
`,
			want: `profiles:
  general:
    preferred: {runtime: codex, model: general-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: general-fallback, reasoning_effort: xhigh}]

  # qa survives with its comment
  qa:
    preferred: {runtime: codex, model: qa-primary, reasoning_effort: high}
    fallbacks: [{runtime: claude, model: qa-fallback, reasoning_effort: xhigh}]
`,
			removal: CategoryBackend,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			workDir := t.TempDir()
			mustMkdir(t, filepath.Join(workDir, ".git"))
			configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
			mustMkdir(t, filepath.Dir(configPath))
			mustWrite(t, configPath, tt.original)

			if _, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
				Scope:    InitScopeUser,
				HomeDir:  homeDir,
				WorkDir:  workDir,
				Removals: []WorkCategory{tt.removal},
			}); err != nil {
				t.Fatalf("remove %s profile: %v", tt.removal, err)
			}
			if got := mustRead(t, configPath); got != tt.want {
				t.Fatalf("category removal changed unrelated trivia:\n%s", profilesConfigGoldenDiff([]byte(tt.want), []byte(got)))
			}
		})
	}

	t.Run("rejects removal that strands an alias before write", func(t *testing.T) {
		original := `profiles:
  backend:
    preferred:
      runtime: codex
      model: &backend_model anchored-backend
      reasoning_effort: high
    fallbacks:
      - runtime: claude
        model: backend-fallback
        reasoning_effort: xhigh

artifact_dir: *backend_model
`
		homeDir := t.TempDir()
		workDir := t.TempDir()
		mustMkdir(t, filepath.Join(workDir, ".git"))
		configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
		mustMkdir(t, filepath.Dir(configPath))
		mustWrite(t, configPath, original)

		_, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
			Scope:    InitScopeUser,
			HomeDir:  homeDir,
			WorkDir:  workDir,
			Removals: []WorkCategory{CategoryBackend},
		})
		if err == nil {
			t.Fatal("expected removal that strands an alias to fail")
		}
		if got := mustRead(t, configPath); got != original {
			t.Fatalf("failed alias removal mutated config\nwant: %q\n got: %q", original, got)
		}
	})
}

func TestProfilesConfigureMergePreservesOtherCategories(t *testing.T) {
	t.Parallel()
	t.Run("replaces one category without touching the other four", func(t *testing.T) {
		original := `profiles:
  # review stays first
  review:
    preferred:
      runtime: codex
      model: review-primary
      reasoning_effort: high
    fallbacks:
      - runtime: claude
        model: review-fallback
        reasoning_effort: xhigh
  # general keeps its comment and Fallback Chain
  general:
    preferred:
      runtime: codex
      model: general-primary
      reasoning_effort: high
    fallbacks:
      - runtime: claude
        model: general-fallback
        reasoning_effort: xhigh
  frontend:
    preferred:
      runtime: claude
      model: frontend-primary
      reasoning_effort: xhigh
    fallbacks:
      - runtime: codex
        model: frontend-fallback
        reasoning_effort: high
  qa:
    preferred:
      runtime: codex
      model: qa-primary
      reasoning_effort: high
    fallbacks:
      - runtime: claude
        model: qa-fallback
        reasoning_effort: xhigh
  backend:
    preferred:
      runtime: codex
      model: backend-primary
      reasoning_effort: high
    fallbacks:
      - runtime: claude
        model: backend-fallback
        reasoning_effort: xhigh
`
		want := strings.Replace(original, "backend-primary", "backend-updated", 1)

		homeDir := t.TempDir()
		workDir := t.TempDir()
		mustMkdir(t, filepath.Join(workDir, ".git"))
		configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
		mustMkdir(t, filepath.Dir(configPath))
		mustWrite(t, configPath, original)

		_, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
			Scope:   InitScopeUser,
			HomeDir: homeDir,
			WorkDir: workDir,
			Profiles: Profiles{
				CategoryBackend: {
					Profile: profileForTest(
						selectionForTest("codex", "backend-updated", "high"),
						selectionForTest("claude", "backend-fallback", "xhigh"),
					),
				},
			},
		})
		if err != nil {
			t.Fatalf("merge backend profile: %v", err)
		}
		if got := mustRead(t, configPath); got != want {
			t.Fatalf("profile merge changed bytes outside backend model:\n%s", profilesConfigGoldenDiff([]byte(want), []byte(got)))
		}
	})

	t.Run("replaces the whole category object", func(t *testing.T) {
		original := `profiles:
    backend:
        preferred:
            runtime: codex
            model: old-backend
            reasoning_effort: high
        fallbacks:
            - runtime: claude
              model: old-first-fallback
              reasoning_effort: xhigh
            - runtime: opencode
              model: old-second-fallback
              reasoning_effort: medium
`
		want := `profiles:
    backend:
        preferred:
            runtime: codex
            model: replacement-backend
            reasoning_effort: max
        fallbacks:
            - runtime: claude
              model: replacement-fallback
              reasoning_effort: high
`

		homeDir := t.TempDir()
		workDir := t.TempDir()
		mustMkdir(t, filepath.Join(workDir, ".git"))
		configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
		mustMkdir(t, filepath.Dir(configPath))
		mustWrite(t, configPath, original)

		_, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
			Scope:   InitScopeUser,
			HomeDir: homeDir,
			WorkDir: workDir,
			Profiles: Profiles{
				CategoryBackend: {
					Profile: profileForTest(
						selectionForTest("codex", "replacement-backend", "max"),
						selectionForTest("claude", "replacement-fallback", "high"),
					),
				},
			},
		})
		if err != nil {
			t.Fatalf("replace backend profile: %v", err)
		}
		if got := mustRead(t, configPath); got != want {
			t.Fatalf("category replacement merged fields instead of replacing the object:\n%s", profilesConfigGoldenDiff([]byte(want), []byte(got)))
		}
	})

	t.Run("removes one category without moving the others", func(t *testing.T) {
		original := `profiles:
    general:
        preferred:
            runtime: codex
            model: general-primary
            reasoning_effort: high
        fallbacks:
            - runtime: claude
              model: general-fallback
              reasoning_effort: xhigh
    backend:
        preferred:
            runtime: codex
            model: backend-primary
            reasoning_effort: high
        fallbacks:
            - runtime: claude
              model: backend-fallback
              reasoning_effort: xhigh
    qa:
        preferred:
            runtime: codex
            model: qa-primary
            reasoning_effort: high
        fallbacks:
            - runtime: claude
              model: qa-fallback
              reasoning_effort: xhigh
    review:
        preferred:
            runtime: codex
            model: review-primary
            reasoning_effort: high
        fallbacks:
            - runtime: claude
              model: review-fallback
              reasoning_effort: xhigh
`
		want := `profiles:
    general:
        preferred:
            runtime: codex
            model: general-primary
            reasoning_effort: high
        fallbacks:
            - runtime: claude
              model: general-fallback
              reasoning_effort: xhigh
    backend:
        preferred:
            runtime: codex
            model: backend-primary
            reasoning_effort: high
        fallbacks:
            - runtime: claude
              model: backend-fallback
              reasoning_effort: xhigh
    review:
        preferred:
            runtime: codex
            model: review-primary
            reasoning_effort: high
        fallbacks:
            - runtime: claude
              model: review-fallback
              reasoning_effort: xhigh
`

		homeDir := t.TempDir()
		workDir := t.TempDir()
		mustMkdir(t, filepath.Join(workDir, ".git"))
		configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
		mustMkdir(t, filepath.Dir(configPath))
		mustWrite(t, configPath, original)

		_, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
			Scope:    InitScopeUser,
			HomeDir:  homeDir,
			WorkDir:  workDir,
			Removals: []WorkCategory{CategoryQA},
			Profiles: Profiles{
				CategoryBackend: {
					Profile: profileForTest(
						selectionForTest("codex", "backend-primary", "high"),
						selectionForTest("claude", "backend-fallback", "xhigh"),
					),
				},
			},
		})
		if err != nil {
			t.Fatalf("remove qa profile: %v", err)
		}
		if got := mustRead(t, configPath); got != want {
			t.Fatalf("category removal moved or changed another entry:\n%s", profilesConfigGoldenDiff([]byte(want), []byte(got)))
		}
	})

	t.Run("appends an added category", func(t *testing.T) {
		original := `profiles:
  backend:
    preferred:
      runtime: codex
      model: backend-primary
      reasoning_effort: high
    fallbacks:
      - runtime: claude
        model: backend-fallback
        reasoning_effort: xhigh
`
		frontendSuffix := `  frontend:
    preferred:
      runtime: claude
      model: frontend-primary
      reasoning_effort: xhigh
    fallbacks:
      - runtime: codex
        model: frontend-fallback
        reasoning_effort: high
`

		homeDir := t.TempDir()
		workDir := t.TempDir()
		mustMkdir(t, filepath.Join(workDir, ".git"))
		configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
		mustMkdir(t, filepath.Dir(configPath))
		mustWrite(t, configPath, original)

		_, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
			Scope:   InitScopeUser,
			HomeDir: homeDir,
			WorkDir: workDir,
			Profiles: Profiles{
				CategoryFrontend: {
					Profile: profileForTest(
						selectionForTest("claude", "frontend-primary", "xhigh"),
						selectionForTest("codex", "frontend-fallback", "high"),
					),
				},
			},
		})
		if err != nil {
			t.Fatalf("add frontend profile: %v", err)
		}
		if got := mustRead(t, configPath); got != original+frontendSuffix {
			t.Fatalf("added category was not appended:\n%s", profilesConfigGoldenDiff([]byte(original+frontendSuffix), []byte(got)))
		}
	})

	t.Run("adds profiles to empty and sectionless configs", func(t *testing.T) {
		tests := []struct {
			name    string
			content string
			keeps   string
		}{
			{name: "empty file"},
			{name: "no profiles section", content: "watch:\n    max_rounds: 7\n", keeps: "max_rounds: 7"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				homeDir := t.TempDir()
				workDir := t.TempDir()
				mustMkdir(t, filepath.Join(workDir, ".git"))
				configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
				mustMkdir(t, filepath.Dir(configPath))
				mustWrite(t, configPath, tt.content)

				_, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
					Scope:   InitScopeUser,
					HomeDir: homeDir,
					WorkDir: workDir,
					Profiles: Profiles{
						CategoryBackend: {
							Profile: profileForTest(
								selectionForTest("codex", "added-backend", "high"),
								selectionForTest("claude", "added-fallback", "xhigh"),
							),
						},
					},
				})
				if err != nil {
					t.Fatalf("add backend profile: %v", err)
				}
				written := mustRead(t, configPath)
				if !strings.Contains(written, "profiles:") || !strings.Contains(written, "model: added-backend") {
					t.Fatalf("written config does not contain added backend profile:\n%s", written)
				}
				if tt.keeps != "" && !strings.Contains(written, tt.keeps) {
					t.Fatalf("written config lost %q:\n%s", tt.keeps, written)
				}
				if _, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir}); err != nil {
					t.Fatalf("load merged config: %v", err)
				}
			})
		}
	})

	t.Run("rejects a replacement that strands an alias", func(t *testing.T) {
		original := `profiles:
    backend:
        preferred:
            runtime: codex
            model: &backend_model anchored-backend
            reasoning_effort: high
        fallbacks:
            - runtime: claude
              model: backend-fallback
              reasoning_effort: xhigh
artifact_dir: *backend_model
`

		homeDir := t.TempDir()
		workDir := t.TempDir()
		mustMkdir(t, filepath.Join(workDir, ".git"))
		configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
		mustMkdir(t, filepath.Dir(configPath))
		mustWrite(t, configPath, original)

		_, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
			Scope:   InitScopeUser,
			HomeDir: homeDir,
			WorkDir: workDir,
			Profiles: Profiles{
				CategoryBackend: {
					Profile: profileForTest(
						selectionForTest("codex", "replacement-backend", "high"),
						selectionForTest("claude", "replacement-fallback", "xhigh"),
					),
				},
			},
		})
		if err == nil {
			t.Fatal("expected replacement that strands an alias to fail")
		}
		if got := mustRead(t, configPath); got != original {
			t.Fatalf("failed alias replacement mutated config\nwant: %q\n got: %q", original, got)
		}
	})
}

func TestEffectiveChangeSet(t *testing.T) {
	t.Parallel()
	backendProfile := profileForTest(selectionForTest("codex", "backend-model", "high"))
	frontendProfile := profileForTest(selectionForTest("claude", "frontend-model", "xhigh"))
	generalProfile := profileForTest(selectionForTest("codex", "general-model", "high"))
	choreProfile := profileForTest(selectionForTest("codex", "chore-model", "medium"))
	qaProfile := profileForTest(selectionForTest("codex", "qa-model", "high"))

	t.Run("classifies configured and absent fragment categories", func(t *testing.T) {
		existing := Profiles{CategoryBackend: {Profile: backendProfile}}
		fragment := Profiles{
			CategoryBackend:  {Profile: backendProfile},
			CategoryFrontend: {Profile: frontendProfile},
		}

		got, err := DeriveEffectiveChangeSet(existing, fragment, nil)
		if err != nil {
			t.Fatalf("derive Effective Change Set: %v", err)
		}
		want := EffectiveChangeSet{Changes: []CategoryChange{
			{Category: CategoryBackend, Kind: ChangeReplaced, Profile: backendProfile},
			{Category: CategoryFrontend, Kind: ChangeAdded, Profile: frontendProfile},
		}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Effective Change Set mismatch\nwant: %#v\n got: %#v", want, got)
		}
	})

	t.Run("classifies configured removal", func(t *testing.T) {
		got, err := DeriveEffectiveChangeSet(
			Profiles{CategoryQA: {Profile: qaProfile}},
			nil,
			[]WorkCategory{CategoryQA},
		)
		if err != nil {
			t.Fatalf("derive Effective Change Set: %v", err)
		}
		want := EffectiveChangeSet{Changes: []CategoryChange{{Category: CategoryQA, Kind: ChangeRemoved}}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Effective Change Set mismatch\nwant: %#v\n got: %#v", want, got)
		}
	})

	t.Run("reports absent removal", func(t *testing.T) {
		got, err := DeriveEffectiveChangeSet(nil, nil, []WorkCategory{CategoryReview})
		if err != nil {
			t.Fatalf("derive Effective Change Set: %v", err)
		}
		want := EffectiveChangeSet{Changes: []CategoryChange{{Category: CategoryReview, Kind: ChangeRemoved}}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Effective Change Set mismatch\nwant: %#v\n got: %#v", want, got)
		}
	})

	t.Run("rejects fragment and removal conflict", func(t *testing.T) {
		_, err := DeriveEffectiveChangeSet(
			nil,
			Profiles{CategoryBackend: {Profile: backendProfile}},
			[]WorkCategory{CategoryBackend},
		)
		if err == nil {
			t.Fatal("expected conflicting category to fail validation")
		}
		if !strings.Contains(err.Error(), string(CategoryBackend)) {
			t.Fatalf("expected validation error to name %q, got %q", CategoryBackend, err)
		}
	})

	t.Run("rejects an unsupported removal category", func(t *testing.T) {
		_, err := DeriveEffectiveChangeSet(nil, nil, []WorkCategory{"not-a-category"})
		if err == nil {
			t.Fatal("expected unsupported removal category to fail validation")
		}
		if !strings.Contains(err.Error(), "not-a-category") {
			t.Fatalf("expected validation error to name the unsupported value, got %q", err)
		}
	})

	t.Run("orders repeated derivations by category", func(t *testing.T) {
		fragment := Profiles{
			CategoryQA:      {Profile: qaProfile},
			CategoryGeneral: {Profile: generalProfile},
			CategoryChore:   {Profile: choreProfile},
		}
		removals := []WorkCategory{CategoryReview, CategoryData, CategoryReview}
		want := EffectiveChangeSet{Changes: []CategoryChange{
			{Category: CategoryGeneral, Kind: ChangeAdded, Profile: generalProfile},
			{Category: CategoryData, Kind: ChangeRemoved},
			{Category: CategoryChore, Kind: ChangeAdded, Profile: choreProfile},
			{Category: CategoryQA, Kind: ChangeAdded, Profile: qaProfile},
			{Category: CategoryReview, Kind: ChangeRemoved},
		}}

		for attempt := 0; attempt < 25; attempt++ {
			got, err := DeriveEffectiveChangeSet(nil, fragment, removals)
			if err != nil {
				t.Fatalf("derive Effective Change Set on attempt %d: %v", attempt, err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("Effective Change Set mismatch on attempt %d\nwant: %#v\n got: %#v", attempt, want, got)
			}
		}
	})
}

func profilesConfigGoldenDiff(golden []byte, actual []byte) string {
	goldenLines := strings.Split(string(golden), "\n")
	actualLines := strings.Split(string(actual), "\n")
	prefix := 0
	for prefix < len(goldenLines) && prefix < len(actualLines) && goldenLines[prefix] == actualLines[prefix] {
		prefix++
	}
	suffix := 0
	for suffix < len(goldenLines)-prefix &&
		suffix < len(actualLines)-prefix &&
		goldenLines[len(goldenLines)-1-suffix] == actualLines[len(actualLines)-1-suffix] {
		suffix++
	}

	const contextLines = 2
	contextStart := max(0, prefix-contextLines)
	goldenChangeEnd := len(goldenLines) - suffix
	actualChangeEnd := len(actualLines) - suffix
	contextEnd := min(len(goldenLines), goldenChangeEnd+contextLines)

	var diff strings.Builder
	diff.WriteString("--- golden\n+++ actual\n")
	fmt.Fprintf(
		&diff,
		"@@ -%d,%d +%d,%d @@\n",
		contextStart+1,
		contextEnd-contextStart,
		contextStart+1,
		(prefix-contextStart)+(actualChangeEnd-prefix)+min(contextLines, suffix),
	)
	for _, line := range goldenLines[contextStart:prefix] {
		fmt.Fprintf(&diff, " %s\n", line)
	}
	for _, line := range goldenLines[prefix:goldenChangeEnd] {
		fmt.Fprintf(&diff, "-%s\n", line)
	}
	for _, line := range actualLines[prefix:actualChangeEnd] {
		fmt.Fprintf(&diff, "+%s\n", line)
	}
	for _, line := range goldenLines[goldenChangeEnd:contextEnd] {
		fmt.Fprintf(&diff, " %s\n", line)
	}
	return diff.String()
}

func TestProfileResolverPreferredOverridePreservesFallbackChain(t *testing.T) {
	t.Parallel()
	config := Builtin()
	override := selectionForTest("claude", "custom-frontier", "")

	got, err := ResolveProfile(config, CategoryBackend, &override)
	if err != nil {
		t.Fatalf("ResolveProfile(backend) with override error = %v", err)
	}

	if got.Source != ProfileSourceInvocation {
		t.Fatalf("expected invocation source, got %q", got.Source)
	}
	if got.Profile.Preferred != override {
		t.Fatalf("expected preferred override %#v, got %#v", override, got.Profile.Preferred)
	}
	wantFallbacks := []AgentSelection{selectionForTest("codex", "gpt-5.5", "xhigh")}
	if !selectionsEqual(got.Profile.Fallbacks, wantFallbacks) {
		t.Fatalf("expected configured fallbacks to survive override\nwant: %#v\ngot:  %#v", wantFallbacks, got.Profile.Fallbacks)
	}
}

func TestProfileConfigAtomicWritesUserAndProjectProfiles(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
watch:
  max_rounds: 4
`)
	userProfile := profileForTest(
		selectionForTest("codex", "user-backend", "high"),
		selectionForTest("claude", "user-fallback", ""),
	)

	userResult, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
		Scope:    InitScopeUser,
		HomeDir:  homeDir,
		WorkDir:  workDir,
		Profiles: Profiles{CategoryBackend: {Profile: userProfile}},
	})
	if err != nil {
		t.Fatalf("write user profiles: %v", err)
	}
	if !userResult.Changed || userResult.Scope != InitScopeUser {
		t.Fatalf("unexpected user write result: %+v", userResult)
	}
	userContent := mustRead(t, userResult.Path)
	if !strings.Contains(userContent, "watch:") || !strings.Contains(userContent, "profiles:") {
		t.Fatalf("expected profiles write to preserve unrelated watch config, got:\n%s", userContent)
	}
	loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
	if err != nil {
		t.Fatalf("load after user profile write: %v", err)
	}
	resolved, err := ResolveProfile(loaded.Config, CategoryBackend, nil)
	if err != nil {
		t.Fatalf("resolve user backend: %v", err)
	}
	if resolved.Source != ProfileSourceUser || !profilesEqual(resolved.Profile, userProfile) {
		t.Fatalf("expected backend to resolve from user profile\nsource=%s profile=%#v", resolved.Source, resolved.Profile)
	}

	projectProfile := profileForTest(
		selectionForTest("claude", "project-backend", "medium"),
		selectionForTest("codex", "project-fallback", "max"),
	)
	projectResult, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
		Scope:    InitScopeProject,
		HomeDir:  homeDir,
		WorkDir:  workDir,
		Profiles: Profiles{CategoryBackend: {Profile: projectProfile}},
	})
	if err != nil {
		t.Fatalf("write project profiles: %v", err)
	}
	if !projectResult.Changed || projectResult.Scope != InitScopeProject {
		t.Fatalf("unexpected project write result: %+v", projectResult)
	}
	loaded, err = Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
	if err != nil {
		t.Fatalf("load after project profile write: %v", err)
	}
	resolved, err = ResolveProfile(loaded.Config, CategoryBackend, nil)
	if err != nil {
		t.Fatalf("resolve project backend: %v", err)
	}
	if resolved.Source != ProfileSourceProject || !profilesEqual(resolved.Profile, projectProfile) {
		t.Fatalf("expected backend to resolve from project profile\nsource=%s profile=%#v", resolved.Source, resolved.Profile)
	}
}

func TestProfileConfigAtomicDryRunAndFailuresLeaveBytesUnchanged(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
	original := "watch:\n  max_rounds: 4\n"
	mustWrite(t, configPath, original)
	validProfiles := Profiles{
		CategoryBackend: {
			Profile: profileForTest(
				selectionForTest("codex", "dry-run-backend", "high"),
				selectionForTest("codex", "dry-run-fallback", "max"),
			),
		},
	}

	result, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
		Scope:    InitScopeUser,
		HomeDir:  homeDir,
		WorkDir:  workDir,
		Profiles: validProfiles,
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("dry-run profiles write: %v", err)
	}
	if result.Changed {
		t.Fatalf("dry-run changed = true, want false: %+v", result)
	}
	if got := mustRead(t, configPath); got != original {
		t.Fatalf("dry-run mutated config\nwant: %q\n got: %q", original, got)
	}

	invalidProfiles := Profiles{
		CategoryBackend: {
			Profile: AgentSelectionProfile{Preferred: selectionForTest("codex", "broken", "high")},
		},
	}
	if _, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
		Scope:    InitScopeUser,
		HomeDir:  homeDir,
		WorkDir:  workDir,
		Profiles: invalidProfiles,
	}); err == nil || !strings.Contains(err.Error(), "fallbacks must include at least one") {
		t.Fatalf("expected invalid profile failure, got %v", err)
	}
	if got := mustRead(t, configPath); got != original {
		t.Fatalf("invalid profile mutated config\nwant: %q\n got: %q", original, got)
	}

	mustWrite(t, configPath, "defaults:\n  agent: codex\n")
	legacyBytes := mustRead(t, configPath)
	if _, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
		Scope:    InitScopeUser,
		HomeDir:  homeDir,
		WorkDir:  workDir,
		Profiles: validProfiles,
	}); err == nil || !strings.Contains(err.Error(), "legacy runtime defaults") {
		t.Fatalf("expected legacy conflict failure, got %v", err)
	}
	if got := mustRead(t, configPath); got != legacyBytes {
		t.Fatalf("legacy conflict mutated config\nwant: %q\n got: %q", legacyBytes, got)
	}

	mustWrite(t, configPath, `
defaults:
  agent: codex
profiles:
  backend:
    preferred:
      runtime: codex
      model: existing
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: existing-fallback
        reasoning_effort: max
`)
	mixedBytes := mustRead(t, configPath)
	if _, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
		Scope:    InitScopeUser,
		HomeDir:  homeDir,
		WorkDir:  workDir,
		Profiles: validProfiles,
	}); err == nil || !strings.Contains(err.Error(), "config mixes legacy runtime defaults with profiles") {
		t.Fatalf("expected same-scope legacy/new conflict failure, got %v", err)
	}
	if got := mustRead(t, configPath); got != mixedBytes {
		t.Fatalf("same-scope conflict mutated config\nwant: %q\n got: %q", mixedBytes, got)
	}

	mustWrite(t, configPath, "watch: [")
	malformedBytes := mustRead(t, configPath)
	if _, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
		Scope:    InitScopeUser,
		HomeDir:  homeDir,
		WorkDir:  workDir,
		Profiles: validProfiles,
	}); err == nil || !strings.Contains(err.Error(), "parse config") {
		t.Fatalf("expected malformed config failure, got %v", err)
	}
	if got := mustRead(t, configPath); got != malformedBytes {
		t.Fatalf("malformed config mutated bytes\nwant: %q\n got: %q", malformedBytes, got)
	}
}

func TestNormalizeProfilesFragmentRequiresDistinctFallback(t *testing.T) {
	t.Parallel()
	preferred := selectionForTest("codex", "only-selection", "high")
	tests := []struct {
		name     string
		profile  AgentSelectionProfile
		contains string
	}{
		{
			name:     "missing fallback",
			profile:  AgentSelectionProfile{Preferred: preferred},
			contains: "one additional distinct authorized and proven Agent Selection is required",
		},
		{
			name: "duplicate fallback",
			profile: AgentSelectionProfile{
				Preferred: preferred,
				Fallbacks: []AgentSelection{preferred},
			},
			contains: "one additional distinct authorized and proven Agent Selection is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NormalizeProfilesFragment(Profiles{
				CategoryBackend: {Profile: tt.profile},
			})
			if err == nil || !strings.Contains(err.Error(), tt.contains) {
				t.Fatalf("NormalizeProfilesFragment error = %v, want %q", err, tt.contains)
			}
		})
	}
}

func TestProfileConfigAtomicParsesStrictProfileFragments(t *testing.T) {
	t.Parallel()
	profiles, err := ParseProfilesFragment([]byte(`
profiles:
  backend:
    preferred:
      runtime: codex
      model: custom-backend
      reasoning_effort: high
    fallbacks:
      - runtime: claude
        model: custom-fallback
        reasoning_effort: ""
`))
	if err != nil {
		t.Fatalf("parse wrapped profile fragment: %v", err)
	}
	if _, ok := profiles[CategoryBackend]; !ok {
		t.Fatalf("expected backend profile in parsed fragment: %#v", profiles)
	}

	_, err = ParseProfilesFragment([]byte(`
profiles:
  backend:
    preferred:
      runtime: codex
      model: custom-backend
      reasoning_effort: high
    fallbacks: []
watch:
  max_rounds: 4
`))
	if err == nil || !strings.Contains(err.Error(), "must not include other top-level keys") {
		t.Fatalf("expected strict fragment rejection, got %v", err)
	}

	_, err = ParseProfilesFragment([]byte(`
backend:
  preferred:
    runtime: codex
    model: custom-backend
    reasoning_effort: high
  fallbacks:
    - runtime: codex
      model: custom-fallback
      reasoning_effort: max
---
backend:
  preferred:
    runtime: claude
    model: hidden
    reasoning_effort: medium
  fallbacks:
    - runtime: codex
      model: hidden-fallback
      reasoning_effort: high
`))
	if err == nil || !strings.Contains(err.Error(), "multiple YAML documents") {
		t.Fatalf("expected multi-document fragment rejection, got %v", err)
	}
}

func TestWriteProfilesConfigRejectsMultiDocumentConfigWithoutMutating(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(workDir, ".git"))
	configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
	current := "watch:\n  max_rounds: 4\n---\nnotify:\n  enabled: false\n"
	mustMkdir(t, filepath.Dir(configPath))
	mustWrite(t, configPath, current)

	_, err := WriteProfilesConfig(context.Background(), ProfileConfigOptions{
		Scope:   InitScopeUser,
		HomeDir: homeDir,
		WorkDir: workDir,
		Profiles: Profiles{
			CategoryBackend: {Profile: AgentSelectionProfile{
				Preferred: AgentSelection{Runtime: "codex", Model: "backend", ReasoningEffort: "high"},
				Fallbacks: []AgentSelection{{Runtime: "codex", Model: "fallback", ReasoningEffort: "max"}},
			}},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "multiple YAML documents") {
		t.Fatalf("expected multi-document config rejection, got %v", err)
	}
	if got := mustRead(t, configPath); got != current {
		t.Fatalf("multi-document config mutated\nwant: %q\n got: %q", current, got)
	}
}

func TestBuiltinWatchDefaultsIncludeCheckGracePeriod(t *testing.T) {
	t.Parallel()
	config := Builtin()

	if config.Watch.CheckGracePeriod != 5*time.Minute {
		t.Fatalf("expected built-in check grace period 5m, got %s", config.Watch.CheckGracePeriod)
	}
	if !strings.Contains(DefaultConfigYAML(), "check_grace_period: 5m") {
		t.Fatalf("expected default config YAML to document check_grace_period, got:\n%s", DefaultConfigYAML())
	}
}

func TestBuiltinReviewSourceExcludesNitpicks(t *testing.T) {
	t.Parallel()
	config := Builtin()

	if config.ReviewSource.IncludeNitpicks {
		t.Fatal("expected built-in review source to exclude nitpicks")
	}
	if !strings.Contains(DefaultConfigYAML(), "include_nitpicks: false") {
		t.Fatalf("expected default config YAML to exclude nitpicks, got:\n%s", DefaultConfigYAML())
	}
}

func TestBuiltinReviewRequestDefaults(t *testing.T) {
	t.Parallel()
	config := Builtin()

	if config.ReviewSource.RequestReview {
		t.Fatal("expected built-in review requests to be disabled")
	}
	if config.ReviewSource.RequestCommand != "@coderabbitai review" {
		t.Fatalf("request command = %q, want %q", config.ReviewSource.RequestCommand, "@coderabbitai review")
	}
	generated := DefaultConfigYAML()
	if !strings.Contains(generated, "request_review: false") {
		t.Fatalf("expected default config YAML to disable review requests, got:\n%s", generated)
	}
	if !strings.Contains(generated, "request_command: \"@coderabbitai review\"") {
		t.Fatalf("expected default config YAML to document the request command, got:\n%s", generated)
	}
}

func TestLoadAppliesReviewRequestHierarchy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		userConfig         string
		projectConfig      string
		wantRequestReview  bool
		wantRequestCommand string
	}{
		{
			name:               "built-in request settings",
			wantRequestCommand: "@coderabbitai review",
		},
		{
			name:               "user config overrides built-in settings",
			userConfig:         "review_source:\n  request_review: true\n  request_command: user review\n",
			wantRequestReview:  true,
			wantRequestCommand: "user review",
		},
		{
			name:               "project config overrides user settings",
			userConfig:         "review_source:\n  request_review: true\n  request_command: user review\n",
			projectConfig:      "review_source:\n  request_review: false\n  request_command: project review\n",
			wantRequestCommand: "project review",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			workDir := t.TempDir()
			mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
			mustMkdir(t, filepath.Join(workDir, ".git"))
			if tt.userConfig != "" {
				mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), tt.userConfig)
			}
			if tt.projectConfig != "" {
				mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), tt.projectConfig)
			}

			loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
			if err != nil {
				t.Fatalf("load config: %v", err)
			}
			if loaded.Config.ReviewSource.RequestReview != tt.wantRequestReview {
				t.Fatalf("request_review = %t, want %t", loaded.Config.ReviewSource.RequestReview, tt.wantRequestReview)
			}
			if loaded.Config.ReviewSource.RequestCommand != tt.wantRequestCommand {
				t.Fatalf("request_command = %q, want %q", loaded.Config.ReviewSource.RequestCommand, tt.wantRequestCommand)
			}
		})
	}
}

func TestValidateRejectsEmptyReviewRequestCommand(t *testing.T) {
	t.Parallel()
	config := Builtin()
	config.ReviewSource.RequestCommand = " \t"

	err := Validate(config)

	if err == nil || !strings.Contains(err.Error(), "review_source.request_command must not be empty") {
		t.Fatalf("expected empty request command rejection, got %v", err)
	}
}

func TestLoadRejectsNonBooleanReviewRequest(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), "review_source:\n  request_review: \"true\"\n")

	_, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})

	if err == nil || !strings.Contains(err.Error(), "request_review") || !strings.Contains(err.Error(), "cannot unmarshal") {
		t.Fatalf("expected non-boolean request_review rejection, got %v", err)
	}
}

func TestLoadRejectsProjectNullReviewRequestInsteadOfInheritingUserValue(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), "review_source:\n  request_review: true\n")
	mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), "review_source:\n  request_review: null\n")

	_, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})

	if err == nil || !strings.Contains(err.Error(), "review_source.request_review must be boolean") {
		t.Fatalf("expected project null request_review rejection, got %v", err)
	}
}

func TestLoadAppliesReviewSourceIncludeNitpicksHierarchy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		userConfig    string
		projectConfig string
		want          bool
	}{
		{name: "built-in excludes nitpicks", want: false},
		{
			name:       "user config includes nitpicks",
			userConfig: "review_source:\n  include_nitpicks: true\n",
			want:       true,
		},
		{
			name:          "project config excludes user nitpicks",
			userConfig:    "review_source:\n  include_nitpicks: true\n",
			projectConfig: "review_source:\n  include_nitpicks: false\n",
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			workDir := t.TempDir()
			mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
			mustMkdir(t, filepath.Join(workDir, ".git"))
			if tt.userConfig != "" {
				mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), tt.userConfig)
			}
			if tt.projectConfig != "" {
				mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), tt.projectConfig)
			}

			loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
			if err != nil {
				t.Fatalf("load config: %v", err)
			}
			if loaded.Config.ReviewSource.IncludeNitpicks != tt.want {
				t.Fatalf("include_nitpicks = %t, want %t", loaded.Config.ReviewSource.IncludeNitpicks, tt.want)
			}
		})
	}
}

func TestLoadAppliesRuntimeConfigHierarchy(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
runtimes:
  codex:
    model: user-codex
    reasoning_effort: user-xhigh
  claude:
    model: user-claude
    reasoning_effort: user-high
  opencode:
    model: user-opencode
`)
	mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), `
runtimes:
  codex:
    model: project-codex
  opencode:
    reasoning_effort: project-opencode-effort
`)

	loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
	if err != nil {
		t.Fatalf("expected config to load, got %v", err)
	}

	if loaded.Config.Runtimes.Codex.Model != "project-codex" {
		t.Fatalf("expected project Codex model override, got %q", loaded.Config.Runtimes.Codex.Model)
	}
	if loaded.Config.Runtimes.Codex.ReasoningEffort != "user-xhigh" {
		t.Fatalf("expected user Codex reasoning to survive project model-only override, got %q", loaded.Config.Runtimes.Codex.ReasoningEffort)
	}
	if loaded.Config.Runtimes.Claude.Model != "user-claude" || loaded.Config.Runtimes.Claude.ReasoningEffort != "user-high" {
		t.Fatalf("expected user Claude defaults to survive unrelated project config, got %#v", loaded.Config.Runtimes.Claude)
	}
	if loaded.Config.Runtimes.OpenCode.Model != "user-opencode" {
		t.Fatalf("expected user OpenCode model to survive project reasoning-only override, got %q", loaded.Config.Runtimes.OpenCode.Model)
	}
	if loaded.Config.Runtimes.OpenCode.ReasoningEffort != "project-opencode-effort" {
		t.Fatalf("expected project OpenCode reasoning override, got %q", loaded.Config.Runtimes.OpenCode.ReasoningEffort)
	}
}

func TestLoadAppliesWorktreeConfigHierarchy(t *testing.T) {
	t.Parallel()
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

func TestLoadVerificationCapacityHierarchy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		userConfig    string
		projectConfig string
		want          int
	}{
		{
			name:          "omission keeps safe default despite project Task Capacity",
			projectConfig: "worktree:\n  concurrency: 7\n",
			want:          1,
		},
		{
			name:       "User Config overrides built-in",
			userConfig: "verification:\n  concurrency: 3\nworktree:\n  concurrency: 9\n",
			want:       3,
		},
		{
			name:          "omitted Project Config key inherits User Config",
			userConfig:    "verification:\n  concurrency: 3\n",
			projectConfig: "worktree:\n  concurrency: 4\n",
			want:          3,
		},
		{
			name:          "Project Config overrides User Config",
			userConfig:    "verification:\n  concurrency: 3\n",
			projectConfig: "verification:\n  concurrency: 2\n",
			want:          2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			workDir := t.TempDir()
			mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
			mustMkdir(t, filepath.Join(workDir, ".git"))
			if tt.userConfig != "" {
				mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), tt.userConfig)
			}
			if tt.projectConfig != "" {
				mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), tt.projectConfig)
			}

			loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
			if err != nil {
				t.Fatalf("load config: %v", err)
			}
			if loaded.Config.Verification.Concurrency != tt.want {
				t.Fatalf("verification.concurrency = %d, want %d", loaded.Config.Verification.Concurrency, tt.want)
			}
		})
	}
}

func TestLoadAppliesAgentLogConfigHierarchy(t *testing.T) {
	t.Parallel()
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

func TestLoadAppliesNotifyConfigHierarchy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		userConfig    string
		projectConfig string
		wantEnabled   bool
		wantCommand   string
	}{
		{
			name:        "builtin only",
			wantEnabled: true,
		},
		{
			name: "user override",
			userConfig: `
notify:
  enabled: false
  command: ntfy publish roundfix-runs
`,
			wantEnabled: false,
			wantCommand: "ntfy publish roundfix-runs",
		},
		{
			name: "project override",
			userConfig: `
notify:
  enabled: false
  command: user-notify
`,
			projectConfig: `
notify:
  enabled: true
  command: ""
`,
			wantEnabled: true,
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
			if loaded.Config.Notify.Enabled != tt.wantEnabled {
				t.Fatalf("expected notify.enabled %t, got %t", tt.wantEnabled, loaded.Config.Notify.Enabled)
			}
			if loaded.Config.Notify.Command != tt.wantCommand {
				t.Fatalf("expected notify.command %q, got %q", tt.wantCommand, loaded.Config.Notify.Command)
			}
		})
	}
}

func TestLoadAppliesStoreRetentionConfigHierarchy(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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

func TestLoadWarnsAndIgnoresDeprecatedDefaultsModel(t *testing.T) {
	t.Parallel()
	const warning = "config: defaults.model is deprecated and ignored; use profiles.<category>.preferred.model\n"

	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
defaults:
  model: user-global-model
`)
	mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), `
defaults:
  model: project-global-model
`)
	var stderr bytes.Buffer

	loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir, Stderr: &stderr})

	if err != nil {
		t.Fatalf("expected config to load, got %v", err)
	}
	if stderr.String() != warning {
		t.Fatalf("expected warning %q, got %q", warning, stderr.String())
	}
	if loaded.Config.Defaults.Model != "" {
		t.Fatalf("expected defaults.model to be ignored, got %q", loaded.Config.Defaults.Model)
	}
	if loaded.Config.Runtimes.Codex.Model != "gpt-5.5" {
		t.Fatalf("expected Codex runtime model to keep built-in value, got %q", loaded.Config.Runtimes.Codex.Model)
	}
}

func TestLoadRejectsUnknownConfigKeys(t *testing.T) {
	t.Parallel()
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

func TestLoadRejectsUnknownRuntimeConfigKeys(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		config   string
		contains string
	}{
		{
			name: "unknown runtime",
			config: `
runtimes:
  local:
    model: llama
`,
			contains: "runtimes.local is not a supported config key",
		},
		{
			name: "unknown runtime default",
			config: `
runtimes:
  codex:
    effort: high
`,
			contains: "runtimes.codex.effort is not a supported config key",
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
				t.Fatalf("expected error containing %q, got %q", tt.contains, err.Error())
			}
		})
	}
}

func TestLoadRejectsUnknownStoreConfigKey(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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

func TestLoadRejectsUnknownNotifyConfigKey(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
notify:
  channel: desktop
`)

	_, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})

	if err == nil {
		t.Fatal("expected config load to fail")
	}
	if !strings.Contains(err.Error(), "notify.channel is not a supported config key") {
		t.Fatalf("expected strict notify key error, got %q", err.Error())
	}
}

func TestLoadRejectsInvalidConfig(t *testing.T) {
	t.Parallel()
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
			name: "invalid check grace period",
			config: `
watch:
  check_grace_period: 0
`,
			contains: "watch.check_grace_period must be greater than 0",
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

func TestValidateVerificationCapacity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		config   string
		contains string
	}{
		{
			name:     "zero",
			config:   "verification:\n  concurrency: 0\n",
			contains: "verification.concurrency must be greater than 0",
		},
		{
			name:     "negative",
			config:   "verification:\n  concurrency: -1\n",
			contains: "verification.concurrency must be greater than 0",
		},
		{
			name:     "float",
			config:   "verification:\n  concurrency: 1.5\n",
			contains: "verification.concurrency must be an integer",
		},
		{
			name:     "string",
			config:   "verification:\n  concurrency: many\n",
			contains: "verification.concurrency must be an integer",
		},
		{
			name:     "unknown key",
			config:   "verification:\n  concurrent: 2\n",
			contains: "verification.concurrent is not a supported config key",
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
				t.Fatalf("expected error containing %q, got %q", tt.contains, err.Error())
			}
		})
	}
}

func TestLoadRejectsWorktreeLocationInsideRepository(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	if !strings.Contains(content, "agent_full_access: false") ||
		!strings.Contains(content, `artifact_dir: ""`) || !strings.Contains(content, "Roundfix Home artifacts/<repo-id>") ||
		!strings.Contains(content, "profiles:") || !strings.Contains(content, "model: gpt-5.6-sol") ||
		!strings.Contains(content, "model: gpt-5.5") || !strings.Contains(content, "model: opus") ||
		!strings.Contains(content, "fallbacks:") ||
		!strings.Contains(content, "specs:") || !strings.Contains(content, `root: "docs/specs"`) ||
		!strings.Contains(content, "worktree:") || !strings.Contains(content, `location: "~/.roundfix/worktrees"`) ||
		!strings.Contains(content, "concurrency: 2") || !strings.Contains(content, "copy: []") ||
		!strings.Contains(content, "store:") || !strings.Contains(content, "journal_retention: 336h") ||
		!strings.Contains(content, "implement:") || !strings.Contains(content, "auto_push: false") ||
		!strings.Contains(content, "notify:") || !strings.Contains(content, "enabled: true") ||
		!strings.Contains(content, `command: ""`) ||
		!strings.Contains(content, "max_run_duration: 2h") {
		t.Fatalf("expected default config content, got %s", content)
	}
	if strings.Contains(content, "resolve.concurrent") || strings.Contains(content, "  concurrent:") {
		t.Fatalf("expected generated config to omit resolve.concurrent, got %s", content)
	}
	for _, forbidden := range []string{"defaults:\n  agent:", "runtimes:"} {
		if strings.Contains(content, forbidden) {
			t.Fatalf("expected generated config to omit legacy selection key %q, got %s", forbidden, content)
		}
	}
	if _, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir}); err != nil {
		t.Fatalf("expected generated User Config to load, got %v", err)
	}
}

func TestProfileGeneratedConfigUsesCompleteProfilesSchema(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(workDir, ".git"))

	result, err := Init(context.Background(), InitOptions{
		Scope:   InitScopeProject,
		HomeDir: homeDir,
		WorkDir: workDir,
	})
	if err != nil {
		t.Fatalf("init project config: %v", err)
	}
	content := mustRead(t, result.Path)
	for _, want := range []string{
		"profiles:",
		"general:",
		"backend:",
		"qa:",
		"review:",
		"model: gpt-5.6-sol",
		"model: gpt-5.5",
		"  frontend:\n" +
			"    preferred:\n" +
			"      runtime: claude\n" +
			"      model: opus\n" +
			"      reasoning_effort: xhigh\n" +
			"    fallbacks:\n" +
			"      - runtime: codex\n" +
			"        model: gpt-5.6-sol\n" +
			"        reasoning_effort: high\n",
		"fallbacks:",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("generated config missing %q:\n%s", want, content)
		}
	}
	for _, forbidden := range []string{"defaults:\n  agent:", "runtimes:", "defaults.model"} {
		if strings.Contains(content, forbidden) {
			t.Fatalf("generated config contains legacy selection key %q:\n%s", forbidden, content)
		}
	}
	loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
	if err != nil {
		t.Fatalf("load generated config: %v", err)
	}
	for _, category := range []WorkCategory{CategoryGeneral, CategoryBackend, CategoryFrontend, CategoryQA, CategoryReview} {
		resolved, err := ResolveProfile(loaded.Config, category, nil)
		if err != nil {
			t.Fatalf("resolve generated profile %s: %v", category, err)
		}
		if resolved.Source != ProfileSourceProject {
			t.Fatalf("%s source = %q, want project", category, resolved.Source)
		}
		if len(resolved.Profile.Fallbacks) == 0 {
			t.Fatalf("%s generated profile has no fallback: %#v", category, resolved.Profile)
		}
	}
}

func TestInitCreatesProjectConfig(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	if content := mustRead(t, path); !strings.Contains(content, "profiles:") ||
		!strings.Contains(content, "model: gpt-5.6-sol") ||
		!strings.Contains(content, "model: opus") ||
		strings.Contains(content, "agent: claude") ||
		strings.Contains(content, "runtimes:") {
		t.Fatalf("expected default config to replace old content, got %s", content)
	}
}

func TestValidateArtifactDirectoryResolvesAndCreatesPaths(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

func selectionForTest(runtime string, model string, reasoningEffort string) AgentSelection {
	return AgentSelection{
		Runtime:         runtime,
		Model:           model,
		ReasoningEffort: reasoningEffort,
	}
}

func profileForTest(preferred AgentSelection, fallbacks ...AgentSelection) AgentSelectionProfile {
	return AgentSelectionProfile{
		Preferred: preferred,
		Fallbacks: append([]AgentSelection(nil), fallbacks...),
	}
}

func profilesEqual(left AgentSelectionProfile, right AgentSelectionProfile) bool {
	return left.Preferred == right.Preferred && selectionsEqual(left.Fallbacks, right.Fallbacks)
}

func selectionsEqual(left []AgentSelection, right []AgentSelection) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

// TestConfiguredWorkCategoriesCoversDefinedOptionalCategories keeps profile
// readiness aligned with what the effective configuration actually defines.
// See ADR-0107.
func TestConfiguredWorkCategoriesCoversDefinedOptionalCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		optional []WorkCategory
		want     []WorkCategory
	}{
		{
			name: "no optional category yields exactly the required categories",
			want: RequiredWorkCategories(),
		},
		{
			name:     "one defined optional category joins the required ones in order",
			optional: []WorkCategory{CategoryData},
			want:     []WorkCategory{CategoryGeneral, CategoryBackend, CategoryFrontend, CategoryData, CategoryQA, CategoryReview},
		},
		{
			name:     "several defined optional categories keep catalog order",
			optional: []WorkCategory{CategoryChore, CategoryInfra},
			want:     []WorkCategory{CategoryGeneral, CategoryBackend, CategoryFrontend, CategoryInfra, CategoryChore, CategoryQA, CategoryReview},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			config := Builtin()
			for _, category := range tt.optional {
				config.Profiles[category] = ProfileEntry{
					Profile: AgentSelectionProfile{
						Preferred: AgentSelection{Runtime: "opencode", Model: "opencode-go/kimi-k3"},
						Fallbacks: []AgentSelection{{Runtime: "claude", Model: "opus", ReasoningEffort: "high"}},
					},
					Source: ProfileSourceProject,
				}
			}

			got := ConfiguredWorkCategories(config)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("configured categories = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConfiguredWorkCategoriesOmitsInheritedCategories keeps a category that
// resolves only by inheritance out of readiness: it contributes no distinct
// Agent Selection tuple to prove.
func TestConfiguredWorkCategoriesOmitsInheritedCategories(t *testing.T) {
	t.Parallel()

	config := Builtin()
	resolved, err := ResolveProfile(config, CategoryDocs, nil)
	if err != nil {
		t.Fatalf("an optional category must still resolve by inheritance: %v", err)
	}
	if resolved.InheritedFrom != CategoryGeneral {
		t.Fatalf("inherited from = %q, want %q", resolved.InheritedFrom, CategoryGeneral)
	}

	for _, category := range ConfiguredWorkCategories(config) {
		if category == CategoryDocs {
			t.Fatal("an inherited category must not be enumerated for readiness")
		}
	}
}
