package agent

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestRuntimeCatalogueBindsCanonicalVariant(t *testing.T) {
	t.Parallel()

	catalogue := RuntimeCatalogue{Models: []string{"default", "opus[1m]", "sonnet"}}
	if !catalogue.AdvertisesModel("opus") {
		t.Fatalf("catalogue models %v do not advertise canonical Agent Model %q", catalogue.Models, "opus")
	}
	if catalogue.AdvertisesModel("haiku") {
		t.Fatalf("catalogue models %v advertise absent Agent Model %q", catalogue.Models, "haiku")
	}
}

func TestProofRefusesAModelTheCatalogueDoesNotAdvertise(t *testing.T) {
	t.Parallel()

	_, err := (ACPXRunner{}).applySessionSelection(context.Background(), SessionSelectionRequest{
		Runtime: RuntimeSpec{ID: "claude", Model: unofferedClaudeModel, ReasoningEffort: "high"},
		Capabilities: SelectionCapabilities{
			CurrentModel: unofferedClaudeModel,
			Models: []ModelCapability{{
				AdapterValue:   unofferedClaudeModel,
				CanonicalModel: unofferedClaudeModel,
				ModelManaged:   true,
			}},
			ReasoningOption: &SelectCapability{ID: "effort", CurrentValue: "high", Values: []string{"low", "medium", "high"}},
		},
		Catalogue: RuntimeCatalogue{Models: honestClaudeModels()},
	}, nil)
	var modelErr *ModelNotAdvertisedError
	if !errors.As(err, &modelErr) {
		t.Fatalf("error = %T %v, want *ModelNotAdvertisedError", err, err)
	}
	if modelErr.Classification() != SelectionModelNotAdvertised {
		t.Fatalf("classification = %q, want %q", modelErr.Classification(), SelectionModelNotAdvertised)
	}
}

func TestProofRefusalNamesTheAdvertisedSet(t *testing.T) {
	t.Parallel()

	catalogue := RuntimeCatalogue{Models: honestClaudeModels()}
	_, err := (ACPXRunner{}).applySessionSelection(context.Background(), SessionSelectionRequest{
		Runtime: RuntimeSpec{ID: "claude", Model: unofferedClaudeModel, ReasoningEffort: "high"},
		Capabilities: SelectionCapabilities{
			CurrentModel: unofferedClaudeModel,
			Models: []ModelCapability{{
				AdapterValue:   unofferedClaudeModel,
				CanonicalModel: unofferedClaudeModel,
				ModelManaged:   true,
			}},
		},
		Catalogue: catalogue,
	}, nil)
	var modelErr *ModelNotAdvertisedError
	if !errors.As(err, &modelErr) {
		t.Fatalf("error = %T %v, want *ModelNotAdvertisedError", err, err)
	}
	if !reflect.DeepEqual(modelErr.Advertised, catalogue.Models) {
		t.Fatalf("advertised models = %v, want %v", modelErr.Advertised, catalogue.Models)
	}
	if want := strings.Join(catalogue.Models, ", "); !strings.Contains(err.Error(), want) {
		t.Fatalf("refusal = %q, want advertised set %q", err, want)
	}
}

func TestProofKeepsTheAdapterRefusalFastPath(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": sessionCapabilitySnapshotFixture(t, "gpt-5.6-sol", []string{"gpt-5.6-sol", "gpt-5.5"}, "reasoning_effort", "medium", []string{"medium", "high"}),
	}))
	harness.setEnv(fakeACPXExitByCall, mustJSONForTest(t, map[string]int{
		"sessions ensure model=" + unofferedCodexModel: 2,
	}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{
		"sessions ensure": "Cannot apply --model " + unofferedCodexModel + ": the ACP agent did not advertise that model.\nAvailable models: gpt-5.6-sol, gpt-5.5\n",
	}))

	_, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{
		Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: unofferedCodexModel, ReasoningEffort: "high"},
		WorkDir: harness.gitRoot,
	})
	var modelErr *ModelNotAdvertisedError
	if !errors.As(err, &modelErr) {
		t.Fatalf("error = %T %v, want *ModelNotAdvertisedError", err, err)
	}
	if modelErr.Err == nil {
		t.Fatal("adapter refusal lost its underlying adapter error")
	}
	wantAdvertised := []string{"gpt-5.6-sol", "gpt-5.5"}
	if !reflect.DeepEqual(modelErr.Advertised, wantAdvertised) {
		t.Fatalf("adapter-advertised models = %v, want %v", modelErr.Advertised, wantAdvertised)
	}
}

func TestProofKeepsAdvertisedSelectionEncoding(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	state := selectionStateFixture(t, "model", "opus[1m]", "opus[1m]", honestClaudeModels(), "effort", "high", []string{"low", "medium", "high"})
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"set effort value=high": selectionStateFixture(t, "effort", "high", "opus[1m]", honestClaudeModels(), "effort", "high", []string{"low", "medium", "high"}),
	}))
	capabilities, err := ParseSessionConfigOptions(
		[]byte(state),
		AdapterEvidence{Command: "claude-agent-acp"},
		SelectionRetention{Model: "opus", ReasoningEffort: "high"},
	)
	if err != nil {
		t.Fatalf("parse advertised claude capabilities: %v", err)
	}
	proof, err := harness.runner.ApplySessionSelection(context.Background(), SessionSelectionRequest{
		Runtime:      RuntimeSpec{ID: "claude", Model: "opus", ReasoningEffort: "high"},
		Session:      SessionRef{Name: "roundfix-live", WorkDir: harness.gitRoot},
		Capabilities: capabilities,
		Catalogue:    RuntimeCatalogue{Models: honestClaudeModels()},
	})
	if err != nil {
		t.Fatalf("prove advertised selection: %v", err)
	}
	if proof.Assignment.Encoding != SelectionEncodingIndependent {
		t.Fatalf("encoding = %q, want %q", proof.Assignment.Encoding, SelectionEncodingIndependent)
	}
}

func TestProofAcceptsAMatchingSelectionAmongSiblings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		assignment SelectionAssignment
		state      SelectionCapabilities
	}{
		{
			name: "codex sol with independent effort",
			assignment: SelectionAssignment{
				Runtime: "codex", Model: "gpt-5.6-sol", ReasoningEffort: "high",
				AdapterModel: "gpt-5.6-sol", ReasoningKey: "reasoning_effort", ReasoningValue: "high", Encoding: SelectionEncodingIndependent,
			},
			state: matchingIndependentSelectionState("gpt-5.6-sol", "gpt-5.6-sol", "reasoning_effort", "high"),
		},
		{
			name: "opencode flash with deferred effort",
			assignment: SelectionAssignment{
				Runtime: "opencode", Model: "openrouter/deepseek/deepseek-v4-flash-0731", ReasoningEffort: "max",
				AdapterModel: "openrouter/deepseek/deepseek-v4-flash-0731", ReasoningKey: "effort", ReasoningValue: "max", Encoding: SelectionEncodingRuntimeDeferred,
			},
			state: SelectionCapabilities{
				CurrentModel: "openrouter/deepseek/deepseek-v4-flash-0731",
				Models: []ModelCapability{{
					AdapterValue: "openrouter/deepseek/deepseek-v4-flash-0731", CanonicalModel: "openrouter/deepseek/deepseek-v4-flash-0731", ModelManaged: true,
				}},
				ReasoningOption: &SelectCapability{ID: "effort", CurrentValue: "low", Values: []string{"low", "high", "max"}},
			},
		},
		{
			name: "claude opus normalizes its echoed alias",
			assignment: SelectionAssignment{
				Runtime: "claude", Model: "opus", ReasoningEffort: "high",
				AdapterModel: "opus", ReasoningKey: "effort", ReasoningValue: "high", Encoding: SelectionEncodingIndependent,
			},
			state: matchingIndependentSelectionState("opus[1m]", "opus", "effort", "high"),
		},
		{
			name: "codex luna with independent effort",
			assignment: SelectionAssignment{
				Runtime: "codex", Model: "gpt-5.6-luna", ReasoningEffort: "max",
				AdapterModel: "gpt-5.6-luna", ReasoningKey: "reasoning_effort", ReasoningValue: "max", Encoding: SelectionEncodingIndependent,
			},
			state: matchingIndependentSelectionState("gpt-5.6-luna", "gpt-5.6-luna", "reasoning_effort", "max"),
		},
		{
			name: "claude sonnet with independent effort",
			assignment: SelectionAssignment{
				Runtime: "claude", Model: "sonnet", ReasoningEffort: "xhigh",
				AdapterModel: "sonnet", ReasoningKey: "effort", ReasoningValue: "xhigh", Encoding: SelectionEncodingIndependent,
			},
			state: matchingIndependentSelectionState("sonnet", "sonnet", "effort", "xhigh"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !selectionStateMatches(tt.assignment, tt.state) {
				t.Fatalf("observed state does not prove matching assignment %#v: %#v", tt.assignment, tt.state)
			}
		})
	}
}

func TestProofStillRejectsAGenuineEffectiveMismatch(t *testing.T) {
	t.Parallel()

	assignment := SelectionAssignment{
		Runtime: "claude", Model: "opus", ReasoningEffort: "high",
		AdapterModel: "opus", ReasoningKey: "effort", ReasoningValue: "high", Encoding: SelectionEncodingIndependent,
	}
	state := matchingIndependentSelectionState("opus[1m]", "opus", "effort", "xhigh")
	if selectionStateMatches(assignment, state) {
		t.Fatalf("observed reasoning %q proved requested reasoning %q", state.ReasoningOption.CurrentValue, assignment.ReasoningValue)
	}

	err := effectiveSelectionError(assignment, state)
	if err.Classification() != EffectiveSelectionMismatch {
		t.Fatalf("classification = %q, want %q", err.Classification(), EffectiveSelectionMismatch)
	}
	if err.EffectiveModel != "opus" || err.EffectiveReasoning != "xhigh" {
		t.Fatalf("effective selection = %q/%q, want opus/xhigh", err.EffectiveModel, err.EffectiveReasoning)
	}
}

func matchingIndependentSelectionState(adapterModel string, canonicalModel string, reasoningKey string, reasoningValue string) SelectionCapabilities {
	return SelectionCapabilities{
		CurrentModel: adapterModel,
		Models: []ModelCapability{{
			AdapterValue: adapterModel, CanonicalModel: canonicalModel, ModelManaged: true,
		}},
		ReasoningOption: &SelectCapability{
			ID: reasoningKey, CurrentValue: reasoningValue, Values: []string{"low", "medium", "high", "xhigh", "max"},
		},
	}
}

func TestSelectionRuntimeDeferredStateMatchesBeforeEffortApplication(t *testing.T) {
	t.Parallel()

	assignment, state := runtimeDeferredSelectionState(t)
	if !selectionStateMatches(assignment, state) {
		t.Fatalf("effective state does not prove deferred assignment %#v", assignment)
	}
}

func TestSelectionRuntimeDeferredStateRejectsDifferentCurrentModel(t *testing.T) {
	t.Parallel()

	assignment, state := runtimeDeferredSelectionState(t)
	state.CurrentModel = "openrouter/deepseek/deepseek-v4-pro"
	if selectionStateMatches(assignment, state) {
		t.Fatalf("a different current model proved deferred assignment %#v", assignment)
	}
}

func runtimeDeferredSelectionState(t *testing.T) (SelectionAssignment, SelectionCapabilities) {
	t.Helper()

	state := parseMeasuredEffortCapabilities(t, measuredGrok45EffortCapabilities)
	assignment, err := PlanSelectionAssignment(
		RuntimeSpec{ID: "opencode", Model: "openrouter/x-ai/grok-4.5", ReasoningEffort: "high"},
		state,
	)
	if err != nil {
		t.Fatalf("plan deferred OpenCode Agent Selection: %v", err)
	}
	if state.ReasoningOption == nil || state.ReasoningOption.CurrentValue == assignment.ReasoningValue {
		t.Fatalf("fixture reasoning value must remain unapplied: state = %#v, assignment = %#v", state.ReasoningOption, assignment)
	}
	return assignment, state
}
