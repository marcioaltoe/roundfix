// Suite: recorded Verification and runtime preferences
// Invariant: portable suggestions omit local filters and profiles select Agent Sessions.
// Boundary IN: embedded decisions, required runtimes, and the rendered fixture guide.
// Boundary OUT: runtime selection and repository-specific recorded decisions.

package baseline

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

func TestTheDecisionCatalogProposesNoRtkCommand(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	asset, ok := catalog.Asset("decisions.json")
	if !ok {
		t.Fatal("missing decisions.json")
	}
	var decisions document
	if err := json.Unmarshal(asset.Data, &decisions); err != nil {
		t.Fatal(err)
	}
	for _, decision := range objectsOrEmpty(decisions["decisions"]) {
		for _, field := range []string{"default", "suggestion"} {
			value, _ := decision[field].(string)
			if strings.HasPrefix(value, "rtk ") {
				t.Errorf("%s %s = %q includes a local filter", decision["id"], field, value)
			}
		}
	}
	for _, want := range []struct{ id, field, value string }{
		{"verification.gate", "default", "make verify"},
		{"verification.incremental", "suggestion", "make verify-incremental"},
	} {
		entry, found := catalog.Decision(want.id)
		decision := incrementalVerificationEntry(t, entry, found)
		if decision[want.field] != want.value {
			t.Errorf("%s %s = %v, want %q", want.id, want.field, decision[want.field], want.value)
		}
	}
}

func TestTheRuntimeDecisionsPreferTheProjectConfigTuples(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	entry, found := catalog.Module("autonomous-work")
	module := incrementalVerificationEntry(t, entry, found)
	for _, want := range []struct{ id, value string }{
		{"runtime.backend", "codex gpt-6.1-sol high"},
		{"runtime.design", "claude opus high"},
	} {
		entry, found := catalog.Decision(want.id)
		decision := incrementalVerificationEntry(t, entry, found)
		if decision["default"] != want.value {
			t.Errorf("%s default = %v, want %q", want.id, decision["default"], want.value)
		}
		if !slices.Contains(stringsOrEmpty(module["requiredDecisions"]), want.id) {
			t.Errorf("autonomous-work no longer requires %s", want.id)
		}
	}
}

func TestTheAutonomousWorkGuideDefersToAgentSelectionProfiles(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	asset, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md")
	if !ok {
		t.Fatal("missing autonomous-work golden guide")
	}
	text := strings.Join(strings.Fields(string(asset.Data)), " ")
	for _, phrase := range []string{
		"Agent Selection Profiles in Project Config choose each Agent Session's ACP Runtime, model, and reasoning effort",
		"runs on the Light Tier when it is available",
	} {
		if !strings.Contains(text, phrase) {
			t.Errorf("guide lacks %q", phrase)
		}
	}
	if strings.Contains(text, "Default backend work uses") {
		t.Error("guide still assigns a default backend runtime")
	}
}
