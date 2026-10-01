// Suite: upstream removals and setup names
// Invariant: catalog entries omit removed skills and built-in profiles follow upstream setup names.
// Boundary IN: embedded modules, bundles, profiles, and setup snapshots; planted invalid entries.
// Boundary OUT: upstream synchronization, skill installation, and production diagnostics.
package baseline

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var upstreamRemovedSkills = []string{
	"review", "triage", "resolving-merge-conflicts", "lesson-learned",
	"accessibility", "best-practices", "performance", "web-quality-audit",
	"firecrawl-developer-index", "firecrawl-monitor", "firecrawl-research-index",
	"golang-benchmark", "golang-continuous-integration", "golang-database",
	"golang-modernize", "golang-observability", "golang-performance",
}

func removedSkillFindings(modules, bundles, setups []document) []string {
	var findings []string
	check := func(owner, name string) {
		for _, removed := range upstreamRemovedSkills {
			if name == removed {
				findings = append(findings, fmt.Sprintf("%s names removed skill %s", owner, name))
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
			for _, trigger := range objectsOrEmpty(dispatch["triggers"]) {
				triggerID, _ := stringValue(trigger, "id")
				for _, removed := range upstreamRemovedSkills {
					if strings.HasSuffix(triggerID, "."+removed) {
						check(id, removed)
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

var renamedSetups = map[string]string{
	"go-cli-tui": "go", "rust-cli": "rust", "standard-typescript-monorepo": "typescript",
}

func setupNameFindings(profiles []document, setups []document) []string {
	var findings []string
	for _, profile := range profiles {
		id, _ := stringValue(profile, "id")
		if want, listed := renamedSetups[id]; listed {
			got, _ := stringValue(profile, "setup")
			if got != want {
				findings = append(findings, fmt.Sprintf("%s setup %s, want %s", id, got, want))
			}
		}
	}
	for _, setup := range setups {
		id, _ := stringValue(setup, "id")
		switch id {
		case "go-cli", "go-tui", "rust-cli", "typescript-bun":
			findings = append(findings, fmt.Sprintf("%s is a retired setup", id))
		}
		source, _ := objectValue(setup["source"])
		if source["type"] == "github" && source["path"] != "setups/"+id+".txt" {
			findings = append(findings, fmt.Sprintf("%s source.path = %v, want setups/%s.txt", id, source["path"], id))
		}
		owned := false
		for _, skill := range objectsOrEmpty(setup["skills"]) {
			skillSource, _ := objectValue(skill["source"])
			if skill["name"] == "roundfix" && skillSource["type"] == "repo" && skillSource["name"] == "roundfix" {
				owned = true
			}
		}
		if !owned {
			findings = append(findings, fmt.Sprintf("%s omits the Roundfix-owned roundfix entry", id))
		}
	}
	return findings
}

func TestNoCatalogEntryNamesASkillRemovedUpstream(t *testing.T) {
	t.Parallel()
	activations := readUpstreamSkillDocument(t, "skill-activations.json")
	if findings := removedSkillFindings(readUpstreamSkillDocuments(t, "modules/*.json"), objectsOrEmpty(activations["bundles"]), readUpstreamSkillDocuments(t, "setups/*.json")); len(findings) != 0 {
		t.Fatal(findings)
	}
}

func TestARemovedSkillNameInTheCatalogIsReported(t *testing.T) {
	t.Parallel()
	for _, removed := range upstreamRemovedSkills {
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
					want := []string{"fixture names removed skill " + removed}
					if got := removedSkillFindings(modules, bundles, setups); !reflect.DeepEqual(got, want) {
						t.Fatalf("findings = %v, want %v", got, want)
					}
				})
			}
		})
	}
}

func TestEveryBuiltInProfileTakesItsRenamedUpstreamSetup(t *testing.T) {
	t.Parallel()
	if findings := setupNameFindings(readUpstreamSkillDocuments(t, "profiles/*.json"), readUpstreamSkillDocuments(t, "setups/*.json")); len(findings) != 0 {
		t.Fatal(findings)
	}
}

func TestARetiredOrUnownedSetupIsReported(t *testing.T) {
	t.Parallel()
	owned := []any{map[string]any{"name": "roundfix", "source": map[string]any{"type": "repo", "name": "roundfix"}}}
	for _, retired := range []string{"go-cli", "go-tui", "rust-cli", "typescript-bun"} {
		t.Run(retired, func(t *testing.T) {
			got := setupNameFindings(nil, []document{{"id": retired, "skills": owned}})
			want := []string{retired + " is a retired setup"}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("findings = %v, want %v", got, want)
			}
		})
	}
	for _, tc := range []struct {
		name             string
		profiles, setups []document
		want             []string
	}{
		{name: "wrong profile setup", profiles: []document{{"id": "go-cli-tui", "setup": "go-tui"}}, want: []string{"go-cli-tui setup go-tui, want go"}},
		{name: "wrong source path", setups: []document{{"id": "go", "source": map[string]any{"type": "github", "path": "setups/go-tui.txt"}, "skills": owned}}, want: []string{"go source.path = setups/go-tui.txt, want setups/go.txt"}},
		{name: "missing owned entry", setups: []document{{"id": "go"}}, want: []string{"go omits the Roundfix-owned roundfix entry"}},
		{name: "external roundfix", setups: []document{{"id": "go", "skills": []any{map[string]any{"name": "roundfix", "source": map[string]any{"type": "github"}}}}}, want: []string{"go omits the Roundfix-owned roundfix entry"}},
		{name: "later composed setup", profiles: []document{{"id": "composed", "setup": "composed"}}, setups: []document{{"id": "composed", "skills": owned}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := setupNameFindings(tc.profiles, tc.setups); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("findings = %v, want %v", got, tc.want)
			}
		})
	}
}
