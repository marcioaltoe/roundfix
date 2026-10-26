// Suite: ai-memory setup and core dispatch contracts
// Invariant: every setup follows the ai-memory revision and core dispatch is conditional on MCP availability.
// Boundary IN: embedded setup snapshots and module declarations.
// Boundary OUT: network, host skill checkouts, MCP calls and skill installation.
package baseline

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

func TestEveryUpstreamSetupFollowsTheAiMemoryRevision(t *testing.T) {
	t.Parallel()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	required := []string{
		"ai-memory-durable-pages", "ai-memory-handoff", "ai-memory-learning-maintenance",
		"ai-memory-messaging", "ai-memory-retrieval", "ai-memory-routing-install",
	}
	for _, id := range catalog.SetupIDs() {
		t.Run(id, func(t *testing.T) {
			entry, ok := catalog.Setup(id)
			if !ok {
				t.Fatalf("setup %s is missing", id)
			}
			var snapshot assetsSyncSnapshot
			if err := json.Unmarshal(entry.Data, &snapshot); err != nil {
				t.Fatalf("setup %s: %v", id, err)
			}
			if snapshot.Source.Type == "github" && snapshot.Source.Ref != upstreamSetupRevision {
				t.Errorf("setup %s revision = %q, want %s", id, snapshot.Source.Ref, upstreamSetupRevision)
			}
			names := make([]string, 0, len(snapshot.Skills))
			for _, skill := range snapshot.Skills {
				names = append(names, skill.Name)
				if skill.Source.Type == "github" && skill.Source.Ref != upstreamSetupRevision {
					t.Errorf("setup %s skill %s revision = %q, want %s", id, skill.Name, skill.Source.Ref, upstreamSetupRevision)
				}
				if skill.Name == "handoff" || skill.Name == "brand-guidelines" || strings.HasPrefix(skill.Name, "gws-") {
					t.Errorf("setup %s retains retired skill %s", id, skill.Name)
				}
			}
			for _, name := range required {
				if !slices.Contains(names, name) {
					t.Errorf("setup %s lacks skill %s", id, name)
				}
			}
		})
	}
}

func TestTheCoreModuleDispatchesTheAiMemorySkills(t *testing.T) {
	t.Parallel()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range catalog.ModuleIDs() {
		entry, ok := catalog.Module(id)
		if !ok {
			t.Fatalf("module %s is missing", id)
		}
		var module document
		if err := json.Unmarshal(entry.Data, &module); err != nil {
			t.Fatalf("module %s: %v", id, err)
		}
		required := stringsOrEmpty(module["requiredSkills"])
		if slices.Contains(required, "handoff") {
			t.Errorf("module %s still requires handoff", id)
		}
		dispatches := objectsOrEmpty(module["skillDispatch"])
		for _, dispatch := range dispatches {
			if dispatch["skill"] == "handoff" {
				t.Errorf("module %s still dispatches handoff", id)
			}
		}
		if id != "core" {
			continue
		}
		for _, name := range []string{"ai-memory-handoff", "ai-memory-messaging"} {
			if !slices.Contains(required, name) {
				t.Errorf("core lacks required skill %s", name)
			}
			found := false
			for _, dispatch := range dispatches {
				if dispatch["skill"] != name {
					continue
				}
				for _, trigger := range objectsOrEmpty(dispatch["triggers"]) {
					if trigger["id"] != "trigger.core."+name {
						continue
					}
					found = true
					when, _ := stringValue(trigger, "when")
					if !strings.HasSuffix(when, "when the session exposes the ai-memory MCP tools.") {
						t.Errorf("core skill %s trigger is not conditional on MCP availability: %q", name, when)
					}
				}
			}
			if !found {
				t.Errorf("core skill %s lacks trigger.core.%s", name, name)
			}
		}
	}
}
