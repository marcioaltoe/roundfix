package config

// Suite: rendered Run Budget configuration
// Invariant: the generated configuration explains the Implement-specific renewal beside the budget setting.
// Boundary IN: DefaultConfigYAML user-facing text.
// Boundary OUT: runtime budget behavior, covered by daemon and CLI tests.

import (
	"strings"
	"testing"
)

func TestRenderedConfigStatesTheImplementBudgetRenewal(t *testing.T) {
	const renewal = "# An Implement Run's allowance renews at each Task settlement."
	content := DefaultConfigYAML()

	if !strings.Contains(content, "budget:\n") || !strings.Contains(content, "  "+renewal+"\n") {
		t.Fatalf("DefaultConfigYAML() is missing %q beside the Run Budget:\n%s", renewal, content)
	}
}
