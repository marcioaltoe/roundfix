// Suite: Roundfix skill membership and dispatch
// Invariant: every embedded setup lists Roundfix and every profile selects its dispatch trigger.
// Boundary IN: embedded setup membership and selected module dispatch entries.
// Boundary OUT: skill installation, upstream synchronization, and CLI output.

package baseline

import (
	"reflect"
	"testing"
)

func setupsWithoutRoundfix(catalog *Catalog) []string {
	var missing []string
	for _, id := range catalog.SetupIDs() {
		found := false
		for _, skill := range objectsOrEmpty(catalog.setups[id]["skills"]) {
			if skill["name"] == "roundfix" {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, id)
		}
	}
	return missing
}

func TestEverySetupListsTheRoundfixSkill(t *testing.T) {
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	if missing := setupsWithoutRoundfix(catalog); len(missing) != 0 {
		t.Fatalf("setups without the Roundfix skill: %v", missing)
	}
}

func TestASetupWithoutTheRoundfixSkillIsReported(t *testing.T) {
	catalog := &Catalog{setups: map[string]document{
		"missing-roundfix": {"skills": []any{map[string]any{"name": "qa-gate"}}},
		"has-roundfix":     {"skills": []any{map[string]any{"name": "roundfix"}}},
	}}
	if got := setupsWithoutRoundfix(catalog); !reflect.DeepEqual(got, []string{"missing-roundfix"}) {
		t.Fatalf("setups without Roundfix = %v, want [missing-roundfix]", got)
	}
}

func TestEveryProfileDispatchesTheRoundfixSkill(t *testing.T) {
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range catalog.ProfileIDs() {
		t.Run(id, func(t *testing.T) {
			modules, ok := catalog.OrderedModules(id)
			if !ok {
				t.Fatalf("profile %s has no selected modules", id)
			}
			for _, moduleID := range modules {
				for _, dispatch := range objectsOrEmpty(catalog.modules[moduleID]["skillDispatch"]) {
					if dispatch["skill"] != "roundfix" {
						continue
					}
					for _, trigger := range objectsOrEmpty(dispatch["triggers"]) {
						if trigger["id"] == "trigger.autonomous-work.roundfix" {
							return
						}
					}
				}
			}
			t.Fatal("selected modules do not dispatch the Roundfix skill")
		})
	}
}
