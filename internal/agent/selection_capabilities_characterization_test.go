package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
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
		SelectionRetention{},
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
		SelectionRetention{},
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
		SelectionRetention{},
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

// TestCharacterizationDeclaredBreakOversizedOptionRetainsInsteadOfRefusing
// records the first break this Spec declares. Until Task 02, an advertised
// option above the bound was refused whole and answered
// too_many_option_values, missing_model_state, and contradictory_response —
// the exact triple `roundfix profiles validate --category data` returned
// against opencode 1.18.15 on 2026-08-08. The projection now accepts the
// advertised set and bounds what it retains, so the same payload projects and
// keeps the requested Agent Model.
func TestCharacterizationDeclaredBreakOversizedOptionRetainsInsteadOfRefusing(t *testing.T) {
	t.Parallel()

	fixture := oversizedModelCapabilityFixture(t, "opencode-go/kimi-k3")
	capabilities, err := ParseSessionConfigOptions(
		[]byte(fixture),
		AdapterEvidence{Command: "opencode"},
		SelectionRetention{Model: "opencode-go/kimi-k3"},
	)
	if err != nil {
		t.Fatalf("an advertised option above the bound must project: %v", err)
	}

	if !slices.ContainsFunc(capabilities.Models, func(model ModelCapability) bool {
		return model.AdapterValue == "opencode-go/kimi-k3"
	}) {
		t.Fatalf("retention dropped the requested Agent Model: %#v", capabilities.Models)
	}

	modelOption := capabilities.Options[selectCapabilityIndex(capabilities.Options, "model")]
	if modelOption.AdvertisedCount <= len(modelOption.Values) {
		t.Fatalf("advertised %d, retained %d; retention must bound an oversized option",
			modelOption.AdvertisedCount, len(modelOption.Values))
	}
	if len(modelOption.Values) != maxRetainedCapabilityValues {
		t.Fatalf("retained %d values, want the bound %d", len(modelOption.Values), maxRetainedCapabilityValues)
	}
}

// TestCharacterizationInvariantOversizedOptionStillFailsClosedAboveTheCeiling
// keeps the refusal that size alone no longer triggers: a payload larger than
// any plausible catalog is still malformed input.
func TestCharacterizationInvariantOversizedOptionStillFailsClosedAboveTheCeiling(t *testing.T) {
	t.Parallel()

	fixture := ceilingBreachingModelCapabilityFixture(t, "opencode-go/kimi-k3")
	_, err := ParseSessionConfigOptions(
		[]byte(fixture),
		AdapterEvidence{Command: "opencode"},
		SelectionRetention{Model: "opencode-go/kimi-k3"},
	)
	if err == nil {
		t.Fatal("an advertised option above the absolute ceiling must be refused")
	}
	var evidence *CapabilityEvidenceError
	if !errors.As(err, &evidence) {
		t.Fatalf("error = %v, want CapabilityEvidenceError", err)
	}
	if !strings.Contains(err.Error(), CapabilityIssueTooManyValues) {
		t.Fatalf("error %q missing issue %q", err.Error(), CapabilityIssueTooManyValues)
	}
}

// oversizedModelCapabilityFixture reproduces the shape the adopted measurement
// recorded — a model option whose values carry provider-prefixed identifiers —
// at a size above the bound rather than at the measured 417.
func oversizedModelCapabilityFixture(t *testing.T, current string, alsoAdvertised ...string) string {
	t.Helper()
	return modelCapabilityFixtureOfSize(t, maxRetainedCapabilityValues+1, current, alsoAdvertised...)
}

// ceilingBreachingModelCapabilityFixture exceeds the absolute ceiling above
// which an advertised option is refused as implausible rather than merely
// large.
func ceilingBreachingModelCapabilityFixture(t *testing.T, current string) string {
	t.Helper()
	return modelCapabilityFixtureOfSize(t, maxAdvertisedCapabilityValues+1, current)
}

func modelCapabilityFixtureOfSize(t *testing.T, size int, current string, alsoAdvertised ...string) string {
	t.Helper()

	seen := make(map[string]struct{}, size)
	values := make([]map[string]string, 0, size)
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
	for index := 0; len(values) < size; index++ {
		add(fmt.Sprintf("openrouter/vendor-%04d/model", index))
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
