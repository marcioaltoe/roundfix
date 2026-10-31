// Suite: retired catalog skills
// Invariant: no required list, dispatch, trigger, bundle, or snapshot names a Retired Skill.
// Boundary IN: embedded catalog documents and planted retired names.
// Boundary OUT: installed-copy reporting and asset synchronization.
package baseline

import (
	"reflect"
	"testing"
)

func TestNoCatalogEntryNamesARetiredSkill(t *testing.T) {
	t.Parallel()
	if got, want := RetiredSkills(), []string{"council", "design-patterns", "diagnosing-bugs", "fixing-metadata", "game-changing-features", "implement-spec", "refactoring-analysis", "tdd", "the-fool", "to-prompt"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("retired skills = %v, want %v", got, want)
	}
	activations := readUpstreamSkillDocument(t, "skill-activations.json")
	setups := readUpstreamSkillDocuments(t, "setups/*.json")
	bundles := objectsOrEmpty(activations["bundles"])
	for _, setup := range setups {
		bundles = append(bundles, objectsOrEmpty(setup["activationBundles"])...)
	}
	if findings := namedSkillFindings(RetiredSkills(), "retired", readUpstreamSkillDocuments(t, "modules/*.json"), bundles, setups); len(findings) != 0 {
		t.Fatal(findings)
	}
}

func TestARetiredSkillNameInTheCatalogIsReported(t *testing.T) {
	t.Parallel()
	for _, removed := range RetiredSkills() {
		t.Run(removed, func(t *testing.T) {
			for _, field := range []string{"required", "dispatch", "trigger", "bundle", "setup"} {
				t.Run(field, func(t *testing.T) {
					var modules, bundles, setups []document
					switch field {
					case "required":
						modules = []document{{"id": "fixture", "requiredSkills": []any{removed}}}
					case "dispatch":
						modules = []document{{"id": "fixture", "skillDispatch": []any{map[string]any{"skill": removed}}}}
					case "trigger":
						modules = []document{{"id": "fixture", "skillDispatch": []any{map[string]any{"skill": "retained", "triggers": []any{map[string]any{"id": "trigger.fixture." + removed}}}}}}
					case "bundle":
						bundles = []document{{"id": "fixture", "skills": []any{removed}}}
					case "setup":
						setups = []document{{"id": "fixture", "skills": []any{map[string]any{"name": removed}}}}
					}
					want := []string{"fixture names retired skill " + removed}
					if got := namedSkillFindings(RetiredSkills(), "retired", modules, bundles, setups); !reflect.DeepEqual(got, want) {
						t.Fatalf("findings = %v, want %v", got, want)
					}
				})
			}
		})
	}
}

func TestTheTypeScriptSnapshotsCarryTheCatalogActivationBundles(t *testing.T) {
	t.Parallel()
	activations := readUpstreamSkillDocument(t, "skill-activations.json")
	want := objectsOrEmpty(activations["bundles"])
	for _, id := range []string{"typescript", "go-cli-typescript-bun"} {
		t.Run(id, func(t *testing.T) {
			t.Parallel()
			setup := readUpstreamSkillDocument(t, "setups/"+id+".json")
			got := objectsOrEmpty(setup["activationBundles"])
			if len(got) != len(want) {
				t.Errorf("setup %s bundles = %d, want %d", id, len(got), len(want))
			}
			for _, bundle := range want {
				var matches []document
				for _, candidate := range got {
					if candidate["id"] == bundle["id"] {
						matches = append(matches, candidate)
					}
				}
				if len(matches) != 1 {
					t.Errorf("setup %s bundle %s occurs %d times, want 1", id, bundle["id"], len(matches))
					continue
				}
				if !reflect.DeepEqual(matches[0]["skills"], bundle["skills"]) {
					t.Errorf("setup %s bundle %s skills = %v, want %v", id, bundle["id"], matches[0]["skills"], bundle["skills"])
				}
			}
		})
	}
}
