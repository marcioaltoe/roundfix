// Suite: Agent Selection catalogue characterization.
// Invariant: today's proof accepts an echoed claude model and refuses it on codex.
// Boundary IN: ProveExactSelection and the fake ACPX process boundary.
// Boundary OUT: the pre-request catalogue and membership verdict added by later Tasks.
package agent

import (
	"context"
	"errors"
	"testing"
)

const (
	unofferedClaudeModel = "opus-9-does-not-exist"
	unofferedCodexModel  = "gpt-9-does-not-exist"
)

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

	modelOptions := []any{
		map[string]string{"value": "default"},
		map[string]string{"value": "opus[1m]"},
		map[string]string{"value": "claude-fable-5[1m]"},
		map[string]string{"value": "sonnet"},
		map[string]string{"value": "haiku"},
		map[string]string{"value": unofferedClaudeModel},
	}
	return mustJSONForTest(t, map[string]any{
		"schema": "acpx.session.v1",
		"acpx": map[string]any{
			"current_model_id": unofferedClaudeModel,
			"config_options": []any{
				map[string]any{
					"id":           "model",
					"category":     "model",
					"type":         "select",
					"currentValue": unofferedClaudeModel,
					"options":      modelOptions,
				},
				map[string]any{
					"id":           "effort",
					"type":         "select",
					"currentValue": "medium",
					"options": []any{
						map[string]string{"value": "low"},
						map[string]string{"value": "medium"},
						map[string]string{"value": "high"},
					},
				},
			},
		},
	})
}
