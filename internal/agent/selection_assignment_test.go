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
