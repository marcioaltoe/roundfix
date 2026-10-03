package agent

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/config"
)

// The published configOptions come from Cursor's ACP model-switching forum
// response (post 9, April 14, 2026). Only public option fields are retained;
// currentValue is default[] as required by ADR-0217's session/new fixture.
// Source: https://forum.cursor.com/t/157312/9
func cursorPublishedProjection(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/cursor-session-new-published.json")
	if err != nil {
		t.Fatal(err)
	}
	var response map[string]json.RawMessage
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatal(err)
	}
	response["action"] = json.RawMessage(`"config_set"`)
	response["configId"] = json.RawMessage(`"model"`)
	response["value"] = json.RawMessage(`"default[]"`)
	projected, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	return projected
}

func TestCursorCapabilitiesReadWholeModelValues(t *testing.T) {
	for _, id := range []string{"cursor", "cursor-custom"} {
		t.Run(id, func(t *testing.T) {
			retention := RetentionFor(RuntimeSpec{ID: id, Model: "grok-4-20[thinking=true]"})
			if !retention.WholeModelValues {
				t.Fatal("Cursor retention does not read whole values")
			}
			capabilities, err := ParseSessionConfigOptions(cursorPublishedProjection(t), AdapterEvidence{Command: "cursor-agent acp"}, retention)
			if err != nil {
				t.Fatal(err)
			}
			if capabilities.CurrentModel != "default[]" || capabilities.ReasoningOption != nil || len(capabilities.Models) != 25 {
				t.Fatalf("capabilities = %+v", capabilities)
			}
			for _, model := range capabilities.Models {
				if model.AdapterValue != model.CanonicalModel || !model.ModelManaged || model.ReasoningEffort != "" {
					t.Fatalf("model was split: %+v", model)
				}
			}
			for _, value := range []string{"default[]", "grok-4-20[thinking=true]", "gpt-5.4[reasoning=medium,context=272k,fast=false]"} {
				if _, ok := modelByAdapterValue(capabilities.Models, value); !ok {
					t.Errorf("missing value %s", value)
				}
			}
		})
	}
}

func TestCursorSelectionAssignsTheExactValue(t *testing.T) {
	runtime := RuntimeSpec{ID: "cursor", Model: "grok-4-20[thinking=true]"}
	capabilities, err := ParseSessionConfigOptions(cursorPublishedProjection(t), AdapterEvidence{Command: "cursor-agent acp"}, RetentionFor(runtime))
	if err != nil {
		t.Fatal(err)
	}
	assignment, err := PlanSelectionAssignment(runtime, capabilities)
	if err != nil {
		t.Fatal(err)
	}
	want := SelectionAssignment{Runtime: "cursor", Model: runtime.Model, AdapterModel: runtime.Model, Encoding: SelectionEncodingModelManaged}
	if assignment != want {
		t.Fatalf("assignment = %+v, want %+v", assignment, want)
	}
	runtime.Model = "grok-4-20"
	_, err = PlanSelectionAssignment(runtime, capabilities)
	var unsupported *SelectionUnsupportedError
	if !errors.As(err, &unsupported) || unsupported.Kind != SelectionModelNotAdvertised {
		t.Fatalf("error = %v, want model_not_advertised", err)
	}
}

// Before the runtime change this same malformed_model_value refusal also
// applied to cursor. Other runtimes must preserve that bracket parsing.
func TestOtherRuntimesKeepBracketParsing(t *testing.T) {
	for _, id := range []string{"codex", "claude", "opencode", "claude-custom"} {
		t.Run(id, func(t *testing.T) {
			retention := RetentionFor(RuntimeSpec{ID: id})
			if retention.WholeModelValues {
				t.Fatal("whole values enabled for another runtime")
			}
			_, err := ParseSessionConfigOptions(cursorPublishedProjection(t), AdapterEvidence{Command: "fixture"}, retention)
			var invalid *CapabilityEvidenceError
			if !errors.As(err, &invalid) || !reflect.DeepEqual(invalid.Issues, []string{CapabilityIssueMalformedModelValue}) {
				t.Fatalf("error = %v, want malformed_model_value", err)
			}
		})
	}
}

func TestCursorRuntimeSpecUsesCursorAgentACP(t *testing.T) {
	runtime, err := RuntimeFor(RuntimeOptions{Agent: "cursor", Model: "default[]"})
	if err != nil {
		t.Fatal(err)
	}
	if runtime.ID != "cursor" || runtime.DisplayName != "Cursor" || runtime.Protocol != ProtocolACP || runtime.SupportedFullAccessMode != "" || runtime.FullAccessMode != "" {
		t.Fatalf("runtime = %+v", runtime)
	}
	environment := []string{"HOME=" + t.TempDir()}
	invocation, err := resolveAdapterInvocation(runtime, environment)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(invocation.argv, []string{"cursor-agent", "acp"}) {
		t.Fatalf("argv = %v", invocation.argv)
	}
	if _, err := acpxReasoningEffortConfigKey(runtime); err == nil {
		t.Fatal("Cursor accepted a reasoning effort control")
	}
	want := `unsupported Agent "unknown"; supported values: ` + strings.Join(config.SupportedRuntimes(), ", ")
	_, err = RuntimeFor(RuntimeOptions{Agent: "unknown"})
	if err == nil || err.Error() != want {
		t.Fatalf("RuntimeFor error = %v, want %s", err, want)
	}
	_, err = resolveAdapterInvocation(RuntimeSpec{ID: "unknown"}, environment)
	if err == nil || err.Error() != want {
		t.Fatalf("adapter error = %v, want %s", err, want)
	}
}
