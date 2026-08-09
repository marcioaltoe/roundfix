package agent

import (
	"errors"
	"testing"
)

// Suite: Agent Selection effort planning characterization
// Invariant: planning preserves stable encodings and rejects models absent from advertised values.
// Boundary IN: capability parsing and PlanSelectionAssignment.
// Boundary OUT: configuration loading and acpx Agent Session lifecycle.

const measuredGrok45EffortCapabilities = `{"action":"config_set","configId":"model","value":"openrouter/x-ai/grok-4.5","configOptions":[{"id":"model","category":"model","type":"select","currentValue":"openrouter/x-ai/grok-4.5","options":[{"value":"openrouter/x-ai/grok-4.5"}]},{"id":"effort","type":"select","currentValue":"low","options":[{"value":"low"},{"value":"medium"},{"value":"high"}]}]}`

const measuredDeepSeekV4ProEffortCapabilities = `{"action":"config_set","configId":"model","value":"openrouter/deepseek/deepseek-v4-pro","configOptions":[{"id":"model","category":"model","type":"select","currentValue":"openrouter/deepseek/deepseek-v4-pro","options":[{"value":"openrouter/deepseek/deepseek-v4-pro"}]},{"id":"effort","type":"select","currentValue":"high","options":[{"value":"high"},{"value":"xhigh"}]}]}`

func TestSelectionEffortCharacterizationInvariantCodexPlansIndependentEncoding(t *testing.T) {
	t.Parallel()

	capabilities := parseMeasuredEffortCapabilities(t, measuredGrok45EffortCapabilities)
	assignment, err := PlanSelectionAssignment(
		RuntimeSpec{ID: "codex", Model: "openrouter/x-ai/grok-4.5", ReasoningEffort: "high"},
		capabilities,
	)
	if err != nil {
		t.Fatalf("plan Codex Agent Selection: %v", err)
	}
	if assignment.Encoding != SelectionEncodingIndependent {
		t.Fatalf("encoding = %q, want %q", assignment.Encoding, SelectionEncodingIndependent)
	}
	if assignment.ReasoningKey != "effort" || assignment.ReasoningValue != "high" {
		t.Fatalf("reasoning assignment = %q/%q, want effort/high", assignment.ReasoningKey, assignment.ReasoningValue)
	}
}

func TestSelectionEffortCharacterizationInvariantClaudePlansIndependentEncoding(t *testing.T) {
	t.Parallel()

	capabilities := parseMeasuredEffortCapabilities(t, measuredDeepSeekV4ProEffortCapabilities)
	assignment, err := PlanSelectionAssignment(
		RuntimeSpec{ID: "claude", Model: "openrouter/deepseek/deepseek-v4-pro", ReasoningEffort: "xhigh"},
		capabilities,
	)
	if err != nil {
		t.Fatalf("plan Claude Agent Selection: %v", err)
	}
	if assignment.Encoding != SelectionEncodingIndependent {
		t.Fatalf("encoding = %q, want %q", assignment.Encoding, SelectionEncodingIndependent)
	}
	if assignment.ReasoningKey != "effort" || assignment.ReasoningValue != "xhigh" {
		t.Fatalf("reasoning assignment = %q/%q, want effort/xhigh", assignment.ReasoningKey, assignment.ReasoningValue)
	}
}

func TestSelectionEffortCharacterizationInvariantOpenCodeEmptyEffortPlansRuntimeManagedEncoding(t *testing.T) {
	t.Parallel()

	capabilities := parseMeasuredEffortCapabilities(t, measuredGrok45EffortCapabilities)
	assignment, err := PlanSelectionAssignment(
		RuntimeSpec{ID: "opencode", Model: "openrouter/x-ai/grok-4.5"},
		capabilities,
	)
	if err != nil {
		t.Fatalf("plan OpenCode Agent Selection: %v", err)
	}
	if assignment.Encoding != SelectionEncodingRuntimeManaged {
		t.Fatalf("encoding = %q, want %q", assignment.Encoding, SelectionEncodingRuntimeManaged)
	}
}

func TestSelectionEffortCharacterizationInvariantUnadvertisedModelReturnsSelectionUnsupportedError(t *testing.T) {
	t.Parallel()

	capabilities := parseMeasuredEffortCapabilities(t, measuredGrok45EffortCapabilities)
	_, err := PlanSelectionAssignment(
		RuntimeSpec{ID: "opencode", Model: "openrouter/deepseek/deepseek-v4-flash-0731"},
		capabilities,
	)
	if err == nil {
		t.Fatal("an Agent Model absent from advertised values must not plan")
	}
	var unsupported *SelectionUnsupportedError
	if !errors.As(err, &unsupported) {
		t.Fatalf("error = %T %v, want *SelectionUnsupportedError", err, err)
	}
	if unsupported.Classification() != SelectionModelNotAdvertised {
		t.Fatalf("classification = %q, want %q", unsupported.Classification(), SelectionModelNotAdvertised)
	}
	if len(unsupported.AdvertisedModels) != 1 || unsupported.AdvertisedModels[0] != "openrouter/x-ai/grok-4.5" {
		t.Fatalf("advertised models = %#v, want only the measured Grok model", unsupported.AdvertisedModels)
	}
}

func parseMeasuredEffortCapabilities(t *testing.T, fixture string) SelectionCapabilities {
	t.Helper()

	capabilities, err := ParseSessionConfigOptions(
		[]byte(fixture),
		AdapterEvidence{Command: "opencode"},
		SelectionRetention{},
	)
	if err != nil {
		t.Fatalf("parse measured OpenCode effort capabilities: %v", err)
	}
	return capabilities
}
