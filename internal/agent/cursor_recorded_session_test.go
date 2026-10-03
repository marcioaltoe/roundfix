// Suite: recorded Cursor session.
// Invariant: a real Cursor session/new advertises a grok model value that a cursor selection assigns whole, as model_managed.
// Boundary IN: the recorded configOptions of one session/new on 2026-10-03, account fields removed, and the capability parser and assignment planner.
// Boundary OUT: Cursor service, login and process execution; the test reads one data file and writes nothing.
package agent

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// cursorRecordedProjection reads the configOptions that cursor-agent acp
// (2026.07.09-a3815c0) returned for session/new in an empty directory, and
// frames them as the config_set projection the parser reads.
func cursorRecordedProjection(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/cursor-session-recorded.json")
	if err != nil {
		t.Fatal(err)
	}
	var response map[string]json.RawMessage
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatal(err)
	}
	if _, ok := response["configOptions"]; !ok {
		t.Fatal("recording has no configOptions")
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

func TestCursorRecordedSessionAdvertisesGrok(t *testing.T) {
	retention := RetentionFor(RuntimeSpec{ID: "cursor"})
	capabilities, err := ParseSessionConfigOptions(cursorRecordedProjection(t), AdapterEvidence{Command: "cursor-agent acp"}, retention)
	if err != nil {
		t.Fatal(err)
	}
	if capabilities.CurrentModel != "default[]" || capabilities.ReasoningOption != nil {
		t.Fatalf("capabilities = %+v", capabilities)
	}
	var grok []string
	for _, model := range capabilities.Models {
		if model.AdapterValue != model.CanonicalModel || !model.ModelManaged || model.ReasoningEffort != "" {
			t.Fatalf("model was split: %+v", model)
		}
		if strings.HasPrefix(model.AdapterValue, "grok") {
			grok = append(grok, model.AdapterValue)
		}
	}
	if len(grok) == 0 {
		t.Fatal("the recorded session advertises no grok value")
	}
	runtime := RuntimeSpec{ID: "cursor", Model: grok[0]}
	assignment, err := PlanSelectionAssignment(runtime, capabilities)
	if err != nil {
		t.Fatal(err)
	}
	want := SelectionAssignment{Runtime: "cursor", Model: grok[0], AdapterModel: grok[0], Encoding: SelectionEncodingModelManaged}
	if assignment != want {
		t.Fatalf("assignment = %+v, want %+v", assignment, want)
	}
}
