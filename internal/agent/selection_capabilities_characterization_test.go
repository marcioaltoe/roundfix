package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// Characterization corpus for Spec 0088-a-third-runtime-that-can-run.
//
// A test named CharacterizationInvariant pins behavior the Spec must not move;
// changing one is a regression. A test named CharacterizationToday pins
// behavior the Spec intends to break, so a later Task editing it is a declared
// break rather than a silent one.

func TestCharacterizationInvariantRetainsEveryValueAtOrBelowTheBound(t *testing.T) {
	t.Parallel()

	capabilities, err := ParseSessionConfigOptions(
		[]byte(officialAdapterCapabilityFixture()),
		AdapterEvidence{Command: "codex-acp"},
	)
	if err != nil {
		t.Fatalf("project an advertised set within the bound: %v", err)
	}

	wantModels := []string{"gpt-5.6-sol", "gpt-5.6-terra", "gpt-5.6-luna", "gpt-5.5"}
	gotModels := make([]string, 0, len(capabilities.Models))
	for _, model := range capabilities.Models {
		gotModels = append(gotModels, model.AdapterValue)
	}
	if strings.Join(gotModels, ",") != strings.Join(wantModels, ",") {
		t.Fatalf("advertised models = %v, want %v in advertised order", gotModels, wantModels)
	}

	if capabilities.ReasoningOption == nil {
		t.Fatal("advertised reasoning option must survive projection")
	}
	wantEfforts := []string{"low", "medium", "high", "xhigh", "max", "ultra"}
	if strings.Join(capabilities.ReasoningOption.Values, ",") != strings.Join(wantEfforts, ",") {
		t.Fatalf("advertised efforts = %v, want %v in advertised order", capabilities.ReasoningOption.Values, wantEfforts)
	}
}

func TestCharacterizationInvariantParsesBracketedVariantIntoCanonicalAndEffort(t *testing.T) {
	t.Parallel()

	capabilities, err := ParseSessionConfigOptions(
		[]byte(modelVariantCapabilityFixture()),
		AdapterEvidence{Command: "codex-acp"},
	)
	if err != nil {
		t.Fatalf("project a variant-encoded advertised set: %v", err)
	}

	var variant ModelCapability
	found := false
	for _, model := range capabilities.Models {
		if model.AdapterValue == "future-model[xhigh]" {
			variant = model
			found = true
		}
	}
	if !found {
		t.Fatalf("bracketed variant absent from projection: %#v", capabilities.Models)
	}
	if variant.CanonicalModel != "future-model" || variant.ReasoningEffort != "xhigh" {
		t.Fatalf("variant = %#v, want canonical %q and effort %q", variant, "future-model", "xhigh")
	}
	if variant.ModelManaged {
		t.Fatal("a variant-encoded identifier is not model-managed")
	}
}

func TestCharacterizationInvariantUnadvertisedModelIsUnsupportedNotInvalidEvidence(t *testing.T) {
	t.Parallel()

	capabilities, err := ParseSessionConfigOptions(
		[]byte(officialAdapterCapabilityFixture()),
		AdapterEvidence{Command: "codex-acp"},
	)
	if err != nil {
		t.Fatalf("project an advertised set within the bound: %v", err)
	}

	_, err = PlanSelectionAssignment(
		RuntimeSpec{ID: "codex", Model: "never-advertised", ReasoningEffort: "high"},
		capabilities,
	)
	if err == nil {
		t.Fatal("an unadvertised Agent Model must not plan")
	}

	var unsupported *SelectionUnsupportedError
	if !errors.As(err, &unsupported) {
		t.Fatalf("error = %v, want SelectionUnsupportedError", err)
	}
	var evidence *CapabilityEvidenceError
	if errors.As(err, &evidence) {
		t.Fatalf("an unadvertised Agent Model must never read as invalid capability evidence: %v", err)
	}
}

// TestCharacterizationTodayRefusesAnOversizedAdvertisedOption pins the defect
// Spec 0088 exists to remove. Measured against opencode 1.18.15 on 2026-08-08,
// a real adapter advertised 417 model values and Roundfix answered exactly the
// three issues asserted here.
func TestCharacterizationTodayRefusesAnOversizedAdvertisedOption(t *testing.T) {
	t.Parallel()

	fixture := oversizedModelCapabilityFixture(t, "opencode-go/kimi-k3")
	_, err := ParseSessionConfigOptions([]byte(fixture), AdapterEvidence{Command: "opencode"})
	if err == nil {
		t.Fatal("today an advertised option above the bound is refused whole")
	}

	var evidence *CapabilityEvidenceError
	if !errors.As(err, &evidence) {
		t.Fatalf("error = %v, want CapabilityEvidenceError", err)
	}
	for _, issue := range []string{
		CapabilityIssueTooManyValues,
		CapabilityIssueMissingModel,
		CapabilityIssueContradictoryResponse,
	} {
		if !strings.Contains(err.Error(), issue) {
			t.Fatalf("error %q missing issue %q", err.Error(), issue)
		}
	}
}

// oversizedModelCapabilityFixture reproduces the shape the adopted measurement
// recorded — a model option whose values carry provider-prefixed identifiers —
// at a size above the bound rather than at the measured 417.
func oversizedModelCapabilityFixture(t *testing.T, current string, alsoAdvertised ...string) string {
	t.Helper()

	seen := make(map[string]struct{}, maxCapabilityValues*2)
	values := make([]map[string]string, 0, maxCapabilityValues*2)
	add := func(value string) {
		if _, exists := seen[value]; exists {
			return
		}
		seen[value] = struct{}{}
		values = append(values, map[string]string{"value": value})
	}
	add(current)
	for _, value := range alsoAdvertised {
		add(value)
	}
	for index := 0; len(values) <= maxCapabilityValues; index++ {
		add(fmt.Sprintf("openrouter/vendor-%03d/model", index))
	}

	payload := map[string]any{
		"action":   "config_set",
		"configId": "model",
		"value":    current,
		"configOptions": []any{
			map[string]any{
				"id":           "model",
				"category":     "model",
				"type":         "select",
				"currentValue": current,
				"options":      values,
			},
			map[string]any{
				"id":           "effort",
				"type":         "select",
				"currentValue": "max",
				"options":      []any{map[string]string{"value": "max"}},
			},
		},
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("encode oversized capability fixture: %v", err)
	}
	return string(encoded)
}
