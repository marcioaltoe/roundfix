// Suite: core skill requirements
// Invariant: universal skills are required and dispatched by core in every rendered guide.
// Boundary IN: module requirements, exact triggers, and rendered skill-dispatch guides.
// Boundary OUT: skill installation and upstream synchronization.
package baseline

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var coreSkillTriggers = map[string]struct{ id, when string }{
	"crafting-effective-readmes": {"trigger.core.crafting-effective-readmes", "Writing or revising the repository README."},
	"typesafe-ai":                {"trigger.core.typesafe-ai", "Adding programmable semantic judgment (routing, ranking, extraction, verification) or touching TypeSafe/Jev."},
}

func coreSkillFindings(core, typescript document, guides map[string]string) []string {
	var findings []string
	for skill, want := range coreSkillTriggers {
		if !containsString(stringsOrEmpty(core["requiredSkills"]), skill) {
			findings = append(findings, "core does not require "+skill)
		}
		found := false
		for _, dispatch := range objectsOrEmpty(core["skillDispatch"]) {
			name, _ := stringValue(dispatch, "skill")
			if name != skill {
				continue
			}
			triggers := objectsOrEmpty(dispatch["triggers"])
			if len(triggers) == 1 {
				id, _ := stringValue(triggers[0], "id")
				when, _ := stringValue(triggers[0], "when")
				found = id == want.id && when == want.when
			}
		}
		if !found {
			findings = append(findings, "core lacks exact trigger "+want.id)
		}
		line := "- `" + want.id + "`: " + want.when
		for path, guide := range guides {
			present := false
			for _, renderedLine := range strings.Split(guide, "\n") {
				if strings.TrimSpace(renderedLine) == line {
					present = true
				}
			}
			if !present {
				findings = append(findings, fmt.Sprintf("%s lacks trigger %s", path, want.id))
			}
		}
	}
	if containsString(stringsOrEmpty(typescript["requiredSkills"]), "crafting-effective-readmes") {
		findings = append(findings, "typescript requires crafting-effective-readmes")
	}
	for _, dispatch := range objectsOrEmpty(typescript["skillDispatch"]) {
		name, _ := stringValue(dispatch, "skill")
		if name == "crafting-effective-readmes" {
			findings = append(findings, "typescript dispatches crafting-effective-readmes")
		}
		for _, trigger := range objectsOrEmpty(dispatch["triggers"]) {
			id, _ := stringValue(trigger, "id")
			if strings.HasSuffix(id, ".crafting-effective-readmes") {
				findings = append(findings, "typescript names trigger "+id)
			}
		}
	}
	sort.Strings(findings)
	return findings
}

func TestCoreRequiresAndDispatchesTheTypeSafeAndReadmeSkills(t *testing.T) {
	t.Parallel()
	guides := map[string]string{}
	for _, path := range []string{
		"assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md",
		"../../docs/agents/skill-dispatch.md",
	} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		guides[path] = string(data)
	}
	core := readUpstreamSkillDocument(t, "modules/core.json")
	typescript := readUpstreamSkillDocument(t, "modules/typescript.json")
	if findings := coreSkillFindings(core, typescript, guides); len(findings) != 0 {
		t.Fatal(findings)
	}
}

func TestAMissingCoreSkillTriggerIsReported(t *testing.T) {
	t.Parallel()
	for skill, want := range coreSkillTriggers {
		t.Run(skill, func(t *testing.T) {
			var required, dispatches []any
			var guide string
			for name, trigger := range coreSkillTriggers {
				required = append(required, name)
				triggers := []any{}
				if name != skill {
					triggers = append(triggers, map[string]any{"id": trigger.id, "when": trigger.when})
					guide += "- `" + trigger.id + "`: " + trigger.when + "\n"
				}
				dispatches = append(dispatches, map[string]any{"skill": name, "triggers": triggers})
			}
			core := document{"requiredSkills": required, "skillDispatch": dispatches}
			got := coreSkillFindings(core, document{}, map[string]string{"fixture": guide})
			expected := []string{"core lacks exact trigger " + want.id, "fixture lacks trigger " + want.id}
			if !reflect.DeepEqual(got, expected) {
				t.Fatalf("findings = %v, want %v", got, expected)
			}
		})
	}
}
