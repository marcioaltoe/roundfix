package config

import (
	"path/filepath"
	"strings"
	"testing"
)

// Characterization corpus for Spec 0088-a-third-runtime-that-can-run.
//
// A test named CharacterizationInvariant pins behavior the Spec must not move.
// A test named CharacterizationToday pins behavior the Spec intends to break,
// so a later Task editing it is a declared break rather than a silent one.

// TestCharacterizationDeclaredBreakRefusesOpenCodeReasoningEffort records the
// second break this Spec declares. Until Task 03, a profile naming
// `runtime: opencode` with a non-empty `reasoning_effort` loaded and resolved
// with that effort intact. Measured against opencode 1.18.15 on 2026-08-08, the
// adapter answers ACP -32602 for every effort applied before an Agent Session's
// first prompt, so accepting one promised what no Run could keep. Configuration
// now refuses it and names the empty value as the repair.
func TestCharacterizationDeclaredBreakRefusesOpenCodeReasoningEffort(t *testing.T) {
	t.Parallel()

	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
profiles:
  data:
    preferred:
      runtime: opencode
      model: opencode-go/kimi-k3
      reasoning_effort: max
    fallbacks:
      - runtime: claude
        model: opus
        reasoning_effort: high
`)

	_, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
	if err == nil {
		t.Fatal("a non-empty opencode reasoning effort must be refused")
	}
	for _, want := range []string{
		"profiles.data.preferred.reasoning_effort",
		"must be empty for runtime \"opencode\"",
		"model-managed reasoning",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q missing %q", err.Error(), want)
		}
	}
}

// TestOpenCodeReasoningRefusalNamesTheRepairOnEveryFallbackPosition keeps the
// refusal positional-agnostic: a Fallback Chain entry is as unusable as a
// Preferred Selection.
func TestOpenCodeReasoningRefusalNamesTheRepairOnEveryFallbackPosition(t *testing.T) {
	t.Parallel()

	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
profiles:
  data:
    preferred:
      runtime: claude
      model: opus
      reasoning_effort: high
    fallbacks:
      - runtime: opencode
        model: opencode-go/kimi-k3
        reasoning_effort: max
`)

	_, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
	if err == nil {
		t.Fatal("a non-empty opencode reasoning effort must be refused in a Fallback Chain too")
	}
	if !strings.Contains(err.Error(), "model-managed reasoning") {
		t.Fatalf("error %q must name the repair", err.Error())
	}
}

func TestCharacterizationInvariantAcceptsCodexAndClaudeReasoningEffort(t *testing.T) {
	t.Parallel()

	loaded := loadCharacterizationConfig(t, `
profiles:
  backend:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: claude
        model: opus
        reasoning_effort: xhigh
`)

	resolved, err := ResolveProfile(loaded.Config, CategoryBackend, nil)
	if err != nil {
		t.Fatalf("resolve the backend profile: %v", err)
	}
	if resolved.Profile.Preferred.ReasoningEffort != "high" {
		t.Fatalf("codex reasoning effort = %q, want %q", resolved.Profile.Preferred.ReasoningEffort, "high")
	}
	if len(resolved.Profile.Fallbacks) != 1 || resolved.Profile.Fallbacks[0].ReasoningEffort != "xhigh" {
		t.Fatalf("claude fallback reasoning effort not preserved: %#v", resolved.Profile.Fallbacks)
	}
}

func TestCharacterizationInvariantAcceptsAnEmptyReasoningEffort(t *testing.T) {
	t.Parallel()

	loaded := loadCharacterizationConfig(t, `
profiles:
  data:
    preferred:
      runtime: opencode
      model: opencode-go/kimi-k3
      reasoning_effort: ""
    fallbacks:
      - runtime: claude
        model: opus
        reasoning_effort: high
`)

	resolved, err := ResolveProfile(loaded.Config, CategoryData, nil)
	if err != nil {
		t.Fatalf("an empty reasoning effort is valid model-managed state: %v", err)
	}
	if resolved.Profile.Preferred.ReasoningEffort != "" {
		t.Fatalf("reasoning effort = %q, want empty", resolved.Profile.Preferred.ReasoningEffort)
	}
}

func loadCharacterizationConfig(t *testing.T, document string) Loaded {
	t.Helper()
	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), document)

	loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
	if err != nil {
		t.Fatalf("load characterization config: %v", err)
	}
	return loaded
}
