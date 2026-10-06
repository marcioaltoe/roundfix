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
	if got, want := RetiredSkills(), []string{"council", "the-fool"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("retired skills = %v, want %v", got, want)
	}
	activations := readUpstreamSkillDocument(t, "skill-activations.json")
	if findings := namedSkillFindings(RetiredSkills(), "retired", readUpstreamSkillDocuments(t, "modules/*.json"), objectsOrEmpty(activations["bundles"]), readUpstreamSkillDocuments(t, "setups/*.json")); len(findings) != 0 {
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
