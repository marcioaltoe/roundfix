// Suite: Core stack wording
// Invariant: core guidance accommodates multiple languages and scopes technology skills.
// Boundary IN: embedded catalog and rendered Plan postimages or skill dispatch.
// Boundary OUT: applying the Plan and executing repository Verification.

package baseline

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func languageTriggerFindings(activations document, words map[string]string) []string {
	bundles := make(map[string][]string)
	for _, bundle := range objectsOrEmpty(activations["bundles"]) {
		id, _ := stringValue(bundle, "id")
		bundles[id] = stringsOrEmpty(bundle["skills"])
	}
	var findings []string
	for _, activation := range objectsOrEmpty(activations["activations"]) {
		id, _ := stringValue(activation, "id")
		when, _ := stringValue(activation, "when")
		bundle, _ := stringValue(activation, "bundle")
		for _, skill := range bundles[bundle] {
			if word, ok := words[skill]; ok && !strings.Contains(when, word) {
				findings = append(findings, id)
				break
			}
		}
	}
	return findings
}

func coreWordingRules() []stackWordingRule {
	return []stackWordingRule{
		{
			guide:   "docs/agents/agent-instructions.md",
			must:    []string{"Use each language's declared package manager and lockfile workflow."},
			mustNot: []string{"Use the repository's declared package manager and lockfile workflow."},
		},
		{
			guide: "docs/agents/skill-dispatch.md",
			must: []string{
				"Activate every matching required skill before governed work. When one skill has distinct active-module triggers, retain and follow each trigger. When a skill's default conflicts with a Baseline rule or a Repository-Specific Normative Rule, follow the rule.",
				"Writing or changing TypeScript tests.",
			},
			mustNot: []string{"Writing or changing tests."},
		},
	}
}

func TestTheCoreGuidesNameEachLanguageAndTheGoverningRule(t *testing.T) {
	t.Parallel()
	plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), standardTypeScriptDecisions("make verify"))
	rendered := make(map[string]string)
	for _, rule := range coreWordingRules() {
		rendered[rule.guide] = string(planPostimage(t, plan, rule.guide).Content)
	}
	if findings := stackWordingFindings(rendered, coreWordingRules()); len(findings) != 0 {
		t.Fatalf("core wording findings: %v", findings)
	}
}

func TestTheCoreWordingCheckReportsTheWordingItReplaced(t *testing.T) {
	t.Parallel()
	rules := []stackWordingRule{
		{guide: "docs/agents/agent-instructions.md", mustNot: []string{"Use the repository's declared package manager and lockfile workflow."}},
		{guide: "docs/agents/skill-dispatch.md", mustNot: []string{"Writing or changing tests."}},
	}
	rendered := map[string]string{
		"docs/agents/agent-instructions.md": "Use the repository's declared package manager and lockfile workflow.",
		"docs/agents/skill-dispatch.md":     "Writing or changing tests.",
	}
	want := []string{
		"docs/agents/agent-instructions.md: replaced sentence remains: Use the repository's declared package manager and lockfile workflow.",
		"docs/agents/skill-dispatch.md: replaced sentence remains: Writing or changing tests.",
	}
	if got := stackWordingFindings(rendered, rules); !reflect.DeepEqual(got, want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
}

func TestTheRewordedCoreClausesKeepTheirForce(t *testing.T) {
	t.Parallel()
	module, ok := mustEmbeddedCatalog(t).Module("core")
	if !ok {
		t.Fatal("core module is absent")
	}
	var decoded map[string]any
	if err := json.Unmarshal(module.Data, &decoded); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"clause.core.follow-dependency-workflow", "clause.core.activate-matching-skills"} {
		t.Run(id, func(t *testing.T) {
			if force, found := clauseForce(decoded, id); !found || force != "mandatory" {
				t.Fatalf("clause force = (%q, %t), want (mandatory, true)", force, found)
			}
		})
	}
}

func TestEveryTriggerThatDispatchesALanguageSkillNamesItsLanguage(t *testing.T) {
	t.Parallel()
	asset, ok := mustEmbeddedCatalog(t).Asset("skill-activations.json")
	if !ok {
		t.Fatal("activation file is absent")
	}
	activations, diagnostics := decodeDocument(asset.Data, asset.Path)
	if len(diagnostics) != 0 {
		t.Fatalf("activation diagnostics: %v", diagnostics)
	}
	words := map[string]string{"vitest": "TypeScript", "react": "React", "hono": "Hono"}
	if findings := languageTriggerFindings(activations, words); len(findings) != 0 {
		t.Fatalf("unscoped triggers: %v", findings)
	}
	// Count only mapped skills reached by an activation, not unused bundles.
	covered := make(map[string]bool)
	for _, activation := range objectsOrEmpty(activations["activations"]) {
		for _, bundle := range objectsOrEmpty(activations["bundles"]) {
			if bundle["id"] != activation["bundle"] {
				continue
			}
			for _, skill := range stringsOrEmpty(bundle["skills"]) {
				if _, ok := words[skill]; ok {
					covered[skill] = true
				}
			}
		}
	}
	if len(covered) < 3 {
		t.Fatalf("covered language skills = %v, want at least three", covered)
	}
}

func TestATestingTriggerThatNamesNoLanguageIsReported(t *testing.T) {
	t.Parallel()
	activations, diagnostics := decodeDocument([]byte(`{
		"bundles": [{"id": "bundle.testing", "skills": ["testing-boss", "tdd", "vitest"]}],
		"activations": [{"id": "trigger.testing", "owner": "typescript", "when": "Writing or changing tests.", "bundle": "bundle.testing"}]
	}`), "unscoped-testing.json")
	if len(diagnostics) != 0 {
		t.Fatalf("activation diagnostics: %v", diagnostics)
	}
	want := []string{"trigger.testing"}
	if got := languageTriggerFindings(activations, map[string]string{"vitest": "TypeScript"}); !reflect.DeepEqual(got, want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
}

func TestAComposedProfileDispatchesVitestOnlyForTypeScriptTests(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	for _, test := range []struct {
		name          string
		modules       []string
		must, mustNot []string
	}{
		{
			name:    "Go and TypeScript",
			modules: []string{"core", "go", "typescript"},
			must:    []string{"`trigger.testing`: Writing or changing TypeScript tests.", "`vitest`", "`trigger.go.golang-testing`: Writing or reviewing Go tests, fixtures, fuzz tests, or integration tests."},
			mustNot: []string{"Writing or changing tests."},
		},
		{
			name:    "Go only",
			modules: []string{"core", "go"},
			must:    []string{"`trigger.go.golang-testing`: Writing or reviewing Go tests, fixtures, fuzz tests, or integration tests."},
			mustNot: []string{"vitest", "trigger.testing"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			text := renderSkillDispatch(catalog, test.modules)
			rule := stackWordingRule{guide: "docs/agents/skill-dispatch.md", must: test.must, mustNot: test.mustNot}
			if findings := stackWordingFindings(map[string]string{rule.guide: text}, []stackWordingRule{rule}); len(findings) != 0 {
				t.Fatalf("composed dispatch findings: %v", findings)
			}
		})
	}
}
