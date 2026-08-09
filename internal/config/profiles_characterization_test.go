package config

import (
	"path/filepath"
	"testing"
)

// Characterization corpus for Spec 0088-a-third-runtime-that-can-run.
//
// A test named CharacterizationInvariant pins behavior the Spec must not move.
// A test named CharacterizationToday pins behavior the Spec intends to break,
// so a later Task editing it is a declared break rather than a silent one.

// TestCharacterizationDeclaredBreakAcceptsOpenCodeReasoningEffort records a
// later declared break to Spec 0088's corpus. Spec 0089 replaces the refusal
// because runtime_deferred now represents an effort that Preflight proves as
// advertised and the Run applies after warming the Agent Session.
func TestCharacterizationDeclaredBreakAcceptsOpenCodeReasoningEffort(t *testing.T) {
	t.Parallel()

	loaded := loadCharacterizationConfig(t, `
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

	resolved, err := ResolveProfile(loaded.Config, CategoryData, nil)
	if err != nil {
		t.Fatalf("resolve data profile: %v", err)
	}
	if got := resolved.Profile.Preferred.ReasoningEffort; got != "max" {
		t.Fatalf("OpenCode reasoning effort = %q, want %q", got, "max")
	}
}

// TestCharacterizationDeclaredBreakAcceptsOpenCodeReasoningEffortInFallbackChain
// keeps the later acceptance positional-agnostic: a Fallback Selection retains
// the same configured effort as a Preferred Selection.
func TestCharacterizationDeclaredBreakAcceptsOpenCodeReasoningEffortInFallbackChain(t *testing.T) {
	t.Parallel()

	loaded := loadCharacterizationConfig(t, `
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

	resolved, err := ResolveProfile(loaded.Config, CategoryData, nil)
	if err != nil {
		t.Fatalf("resolve data profile: %v", err)
	}
	if len(resolved.Profile.Fallbacks) != 1 {
		t.Fatalf("Fallback Chain length = %d, want 1", len(resolved.Profile.Fallbacks))
	}
	if got := resolved.Profile.Fallbacks[0].ReasoningEffort; got != "max" {
		t.Fatalf("OpenCode fallback reasoning effort = %q, want %q", got, "max")
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
