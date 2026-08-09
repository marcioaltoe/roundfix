package config

import (
	"testing"
)

// Suite: OpenCode profile effort characterization
// Invariant: profile parsing preserves a non-empty OpenCode reasoning effort.
// Boundary IN: profile fragment decoding and Agent Selection normalization.
// Boundary OUT: capability planning and acpx Agent Session execution.

func TestOpenCodeEffortAcceptedConfigurationLoadsAndResolvesNonEmptyReasoningEffort(t *testing.T) {
	t.Parallel()

	profiles, err := ParseProfilesFragment([]byte(`
profiles:
  data:
    preferred:
      runtime: opencode
      model: openrouter/x-ai/grok-4.5
      reasoning_effort: high
    fallbacks:
      - runtime: claude
        model: opus
        reasoning_effort: xhigh
`))
	if err != nil {
		t.Fatalf("parse profiles with a non-empty OpenCode reasoning effort: %v", err)
	}

	resolved, err := ResolveProfile(Config{Profiles: profiles}, CategoryData, nil)
	if err != nil {
		t.Fatalf("resolve data profile: %v", err)
	}
	if got := resolved.Profile.Preferred.ReasoningEffort; got != "high" {
		t.Fatalf("OpenCode reasoning effort = %q, want %q", got, "high")
	}
}
