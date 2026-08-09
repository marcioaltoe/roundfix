package config

import (
	"strings"
	"testing"
)

// Suite: OpenCode profile effort characterization
// Invariant: today's profile parser refuses a non-empty OpenCode reasoning effort.
// Boundary IN: profile fragment decoding and Agent Selection normalization.
// Boundary OUT: capability planning and acpx Agent Session execution.

func TestOpenCodeEffortCharacterizationTodayConfigurationRefusesNonEmptyReasoningEffort(t *testing.T) {
	t.Parallel()

	_, err := ParseProfilesFragment([]byte(`
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
	if err == nil {
		t.Fatal("a non-empty OpenCode reasoning effort must be refused")
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
