// Suite: upstream skill names
// Invariant: catalog skills follow upstream renames and modules require their dispatched skills.
// Boundary IN: module requirements and triggers, activation bundles, setup membership, and profile setup selection.
// Boundary OUT: upstream synchronization and production catalog diagnostics.
package baseline

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"reflect"
	"testing"
)

var upstreamSkillRenames = []struct{ old, current string }{
	{"context7", "context7-cli"},
	{"feature-systems-pattern", "app-renderer-systems"},
	{"rust", "rust-expert"},
}

func upstreamSkillNameFindings(modules, bundles, setups []document) []string {
	var findings []string
	check := func(owner, name string) {
		for _, rename := range upstreamSkillRenames {
			if name == rename.old {
				findings = append(findings, fmt.Sprintf("%s names renamed skill %s", owner, name))
			}
		}
	}
	for _, module := range modules {
		id, _ := stringValue(module, "id")
		for _, name := range stringsOrEmpty(module["requiredSkills"]) {
			check(id, name)
		}
		for _, dispatch := range objectsOrEmpty(module["skillDispatch"]) {
			name, ok := stringValue(dispatch, "skill")
			if !ok {
				name, _ = stringValue(dispatch, "id")
			}
			check(id, name)
			for _, rename := range upstreamSkillRenames {
				if name != rename.old && name != rename.current {
					continue
				}
				want := "trigger." + id + "." + rename.current
				triggers := objectsOrEmpty(dispatch["triggers"])
				if len(triggers) == 0 {
					findings = append(findings, fmt.Sprintf("%s dispatch %s has no renamed trigger", id, name))
				}
				for _, trigger := range triggers {
					got, _ := stringValue(trigger, "id")
					if got != want {
						findings = append(findings, fmt.Sprintf("%s trigger %s, want %s", id, got, want))
					}
				}
			}
		}
	}
	for _, bundle := range bundles {
		id, _ := stringValue(bundle, "id")
		for _, name := range stringsOrEmpty(bundle["skills"]) {
			check(id, name)
		}
	}
	for _, setup := range setups {
		id, _ := stringValue(setup, "id")
		for _, skill := range objectsOrEmpty(setup["skills"]) {
			name, _ := stringValue(skill, "name")
			check(id, name)
		}
	}
	return findings
}

func dispatchedSkillRequirementFindings(modules []document) []string {
	var findings []string
	for _, module := range modules {
		id, _ := stringValue(module, "id")
		required := stringsOrEmpty(module["requiredSkills"])
		for _, dispatch := range objectsOrEmpty(module["skillDispatch"]) {
			name, ok := stringValue(dispatch, "skill")
			if !ok {
				name, _ = stringValue(dispatch, "id")
			}
			if !containsString(required, name) {
				findings = append(findings, fmt.Sprintf("%s dispatches unrequired skill %s", id, name))
			}
		}
	}
	return findings
}

func readUpstreamSkillDocument(t *testing.T, name string) document {
	t.Helper()
	content, err := fs.ReadFile(embeddedAssets, name)
	if err != nil {
		t.Fatal(err)
	}
	var doc document
	if err := json.Unmarshal(content, &doc); err != nil {
		t.Fatal(err)
	}
	return doc
}

func readUpstreamSkillDocuments(t *testing.T, pattern string) []document {
	t.Helper()
	names, err := fs.Glob(embeddedAssets, pattern)
	if err != nil {
		t.Fatal(err)
	}
	if len(names) == 0 {
		t.Fatalf("no catalog documents match %s", pattern)
	}
	var docs []document
	for _, name := range names {
		docs = append(docs, readUpstreamSkillDocument(t, name))
	}
	return docs
}

func TestNoCatalogEntryNamesASkillRenamedUpstream(t *testing.T) {
	t.Parallel()
	modules := readUpstreamSkillDocuments(t, "modules/*.json")
	activations := readUpstreamSkillDocument(t, "skill-activations.json")
	setups := readUpstreamSkillDocuments(t, "setups/*.json")
	if findings := upstreamSkillNameFindings(modules, objectsOrEmpty(activations["bundles"]), setups); len(findings) != 0 {
		t.Fatal(findings)
	}
}

func TestARenamedSkillNameInTheCatalogIsReported(t *testing.T) {
	t.Parallel()
	for _, rename := range upstreamSkillRenames {
		t.Run(rename.old, func(t *testing.T) {
			for _, field := range []string{"required", "dispatch", "bundle", "setup", "trigger"} {
				t.Run(field, func(t *testing.T) {
					var modules, bundles, setups []document
					switch field {
					case "required":
						modules = []document{{"id": "fixture", "requiredSkills": []any{rename.old}}}
					case "dispatch":
						modules = []document{{"id": "fixture", "skillDispatch": []any{map[string]any{"skill": rename.old}}}}
					case "bundle":
						bundles = []document{{"id": "bundle.fixture", "skills": []any{rename.old}}}
					case "setup":
						setups = []document{{"id": "fixture", "skills": []any{map[string]any{"name": rename.old}}}}
					case "trigger":
						modules = []document{{"id": "fixture", "skillDispatch": []any{map[string]any{"skill": rename.current, "triggers": []any{map[string]any{"id": "trigger.fixture." + rename.old}}}}}}
					}
					if findings := upstreamSkillNameFindings(modules, bundles, setups); len(findings) == 0 {
						t.Fatalf("did not report %s in %s", rename.old, field)
					}
				})
			}
		})
	}
}

func TestTheGoCLITUIProfileTakesTheGoSetup(t *testing.T) {
	t.Parallel()
	profile := readUpstreamSkillDocument(t, "profiles/go-cli-tui.json")
	if got, _ := stringValue(profile, "setup"); got != "go" {
		t.Fatalf("profile setup = %s, want go", got)
	}
	setup := readUpstreamSkillDocument(t, "setups/go.json")
	names := map[string]bool{}
	for _, skill := range objectsOrEmpty(setup["skills"]) {
		name, _ := stringValue(skill, "name")
		names[name] = true
	}
	for _, name := range []string{"bubbletea", "tui-design"} {
		if !names[name] {
			t.Errorf("go omits %s", name)
		}
	}
}

func TestEveryModuleRequiresEverySkillItDispatches(t *testing.T) {
	t.Parallel()
	if findings := dispatchedSkillRequirementFindings(readUpstreamSkillDocuments(t, "modules/*.json")); len(findings) != 0 {
		t.Fatal(findings)
	}
}

func TestADispatchedSkillNoModuleRequiresIsReported(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"skill", "id"} {
		t.Run(field, func(t *testing.T) {
			modules := []document{{"id": "fixture", "requiredSkills": []any{"required"}, "skillDispatch": []any{map[string]any{field: "unrequired"}}}}
			want := []string{"fixture dispatches unrequired skill unrequired"}
			if got := dispatchedSkillRequirementFindings(modules); !reflect.DeepEqual(got, want) {
				t.Fatalf("findings = %v, want %v", got, want)
			}
		})
	}
}
