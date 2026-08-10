// Suite: Agent Selection catalogue characterization.
// Invariant: today's proof accepts an echoed claude model and refuses it on codex.
// Boundary IN: ProveExactSelection and the fake ACPX process boundary.
// Boundary OUT: the membership verdict added by task_03.
package agent

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

const (
	unofferedClaudeModel = "opus-9-does-not-exist"
	unofferedCodexModel  = "gpt-9-does-not-exist"
)

func TestRuntimeCatalogueReadsAdvertisedModelsWithoutAnOverride(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": honestClaudeCapabilityPayload(t),
	}))

	proof, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{
		Runtime: RuntimeSpec{
			ID:              "claude",
			Protocol:        ProtocolACP,
			Model:           "default",
			ReasoningEffort: "high",
		},
		WorkDir: harness.gitRoot,
	})
	if err != nil {
		t.Fatalf("prove selection with runtime catalogue: %v", err)
	}
	wantModels := []string{"default", "opus[1m]", "claude-fable-5[1m]", "sonnet", "haiku"}
	if !reflect.DeepEqual(proof.Catalogue.Models, wantModels) {
		t.Fatalf("catalogue models = %v, want %v", proof.Catalogue.Models, wantModels)
	}
	wantEfforts := []string{"low", "medium", "high"}
	if !reflect.DeepEqual(proof.Catalogue.Efforts, wantEfforts) {
		t.Fatalf("catalogue efforts = %v, want %v", proof.Catalogue.Efforts, wantEfforts)
	}

	invocations := readJSONInvocations(t, harness.invocationsPath)
	if len(invocations) < 2 || fakeACPXCommandKey(invocations[0]) != "sessions ensure" || containsArg(invocations[0], "--model") {
		t.Fatalf("first disposable ensure applied a model override: %#v", invocations)
	}
	if fakeACPXCommandKey(invocations[1]) != "sessions show" {
		t.Fatalf("runtime catalogue was not read immediately after the override-free ensure: %#v", invocations)
	}
	if containsCommandKey(invocations, "prompt") {
		t.Fatalf("runtime catalogue proof sent a prompt: %#v", invocations)
	}
}

func TestRuntimeCatalogueRecordsAContaminatedAdvertisement(t *testing.T) {
	t.Parallel()

	proof, err := (ACPXRunner{}).applySessionSelection(context.Background(), SessionSelectionRequest{
		Runtime: RuntimeSpec{
			ID:       "codex",
			Protocol: ProtocolACP,
			Model:    unofferedCodexModel,
		},
		Capabilities: SelectionCapabilities{
			CurrentModel: unofferedCodexModel,
			Models: []ModelCapability{{
				AdapterValue:   unofferedCodexModel,
				CanonicalModel: unofferedCodexModel,
				ModelManaged:   true,
			}},
		},
		Catalogue: RuntimeCatalogue{Models: []string{"gpt-5.6-sol", "gpt-5.5"}},
	}, nil)
	if err != nil {
		t.Fatalf("apply selection with later advertisement: %v", err)
	}
	if proof.Status != SelectionProofStatusProven {
		t.Fatalf("proof status = %q, want unchanged status %q", proof.Status, SelectionProofStatusProven)
	}
	if !proof.Catalogue.Contaminated {
		t.Fatalf("catalogue did not record later advertisement of absent model %q: %#v", unofferedCodexModel, proof.Catalogue)
	}
}

// TestSelectionCatalogueCharacterizationClaudeProvesAnUnofferedModel records
// the false-positive proof observed when claude echoes the requested model
// into its own advertised capabilities.
// Declared break: task_03 changes this success into model_not_advertised.
func TestSelectionCatalogueCharacterizationClaudeProvesAnUnofferedModel(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": contaminatedClaudeCapabilityPayload(t),
	}))

	proof, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{
		Runtime: RuntimeSpec{
			ID:              "claude",
			Protocol:        ProtocolACP,
			Model:           unofferedClaudeModel,
			ReasoningEffort: "high",
		},
		WorkDir: harness.gitRoot,
	})
	if err != nil {
		t.Fatalf("prove today's claude selection: %v", err)
	}
	if proof.Status != SelectionProofStatusProven || proof.Model != unofferedClaudeModel || proof.ReasoningEffort != "high" {
		t.Fatalf("unexpected claude proof: %#v", proof)
	}
}

// TestSelectionCatalogueCharacterizationCodexRefusesAnUnofferedModel keeps the
// adapter-owned refusal that task_03 must preserve as its fast path.
func TestSelectionCatalogueCharacterizationCodexRefusesAnUnofferedModel(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{
		"sessions ensure": 2,
	}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{
		"sessions ensure": "Cannot apply --model " + unofferedCodexModel + ": the ACP agent did not advertise that model.\nAvailable models: gpt-5.6-sol, gpt-5.5\n",
	}))

	_, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{
		Runtime: RuntimeSpec{
			ID:              "codex",
			Protocol:        ProtocolACP,
			Model:           unofferedCodexModel,
			ReasoningEffort: "high",
		},
		WorkDir: harness.gitRoot,
	})
	var modelErr *ModelNotAdvertisedError
	if !errors.As(err, &modelErr) {
		t.Fatalf("error = %T %v, want *ModelNotAdvertisedError", err, err)
	}
	if modelErr.Classification() != SelectionModelNotAdvertised {
		t.Fatalf("classification = %q, want %q", modelErr.Classification(), SelectionModelNotAdvertised)
	}
}

// contaminatedClaudeCapabilityPayload records the measured response shape:
// the requested model becomes both currentValue and an advertised option.
// Declared break: task_02 reads the honest catalogue before accepting this
// post-request payload as capability evidence.
func contaminatedClaudeCapabilityPayload(t *testing.T) string {
	t.Helper()

	models := append(honestClaudeModels(), unofferedClaudeModel)
	return sessionCapabilitySnapshotFixture(t, unofferedClaudeModel, models, "effort", "medium", []string{"low", "medium", "high"})
}

func honestClaudeCapabilityPayload(t *testing.T) string {
	t.Helper()

	return sessionCapabilitySnapshotFixture(t, "default", honestClaudeModels(), "effort", "medium", []string{"low", "medium", "high"})
}

func honestClaudeModels() []string {
	return []string{"default", "opus[1m]", "claude-fable-5[1m]", "sonnet", "haiku"}
}
