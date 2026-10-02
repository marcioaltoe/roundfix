// Suite: Stack rule wording
// Invariant: rendered stack guidance names its scope and preserves clause force.
// Boundary IN: embedded modules and a Standard TypeScript Monorepo Plan.
// Boundary OUT: applying the Plan and executing repository Verification.

package baseline

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type stackWordingRule struct {
	guide   string
	must    []string
	mustNot []string
}

func stackWordingFindings(rendered map[string]string, rules []stackWordingRule) []string {
	var findings []string
	normalize := func(text string) string { return strings.Join(strings.Fields(text), " ") }
	for _, rule := range rules {
		text, ok := rendered[rule.guide]
		if !ok {
			findings = append(findings, fmt.Sprintf("%s: guide is not rendered", rule.guide))
			continue
		}
		text = normalize(text)
		for _, sentence := range rule.must {
			if !strings.Contains(text, normalize(sentence)) {
				findings = append(findings, fmt.Sprintf("%s: missing sentence: %s", rule.guide, sentence))
			}
		}
		for _, sentence := range rule.mustNot {
			if strings.Contains(text, normalize(sentence)) {
				findings = append(findings, fmt.Sprintf("%s: replaced sentence remains: %s", rule.guide, sentence))
			}
		}
	}
	return findings
}

func clauseForce(module any, clauseID string) (enforcement string, found bool) {
	doc, ok := objectValue(module)
	if !ok {
		return "", false
	}
	for _, rule := range objectsOrEmpty(doc["rules"]) {
		for _, clause := range objectsOrEmpty(rule["clauses"]) {
			if clause["id"] == clauseID {
				enforcement, ok = clause["enforcement"].(string)
				return enforcement, ok
			}
		}
	}
	return "", false
}

func typeScriptAndBunWordingRule() stackWordingRule {
	return stackWordingRule{
		guide: "docs/agents/typescript-bun.md",
		must: []string{
			"These rules govern the repository's TypeScript sources and tests. Code in another language follows its own guide.",
			"These rules govern the Bun workspace: the packages under the root `package.json`. A toolchain for another language in the same repository keeps its own commands.",
			"Inside the Bun workspace, use Bun-owned commands for dependency installation, scripts, and lockfile updates, and run tests through the package's `test` script (`bun run test`), never through a bare runner such as `bun test`.",
			"Inside the Bun workspace, do not substitute another JavaScript package manager or runner (npm, pnpm, yarn, npx) or hand-edit the lockfile.",
			"A toolchain for another language in the same repository keeps its own package manager.",
			"Keep TypeScript type errors visible; never hide one to make Verification pass.",
		},
		mustNot: []string{
			"Use Bun-owned commands for dependency installation, scripts, tests, and lockfile updates.",
			"Do not substitute another package manager or hand-edit the lockfile.",
			"Keep type errors visible and preserve the repository's package and lockfile workflow.",
			"When the selected TypeScript/Bun profile treats warnings as errors, every warning reported by Verification blocks completion.",
		},
	}
}

func TestTheTypeScriptAndBunGuideSaysWhatItGoverns(t *testing.T) {
	t.Parallel()
	plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), standardTypeScriptDecisions("make verify"))
	rule := typeScriptAndBunWordingRule()
	rendered := map[string]string{rule.guide: string(planPostimage(t, plan, rule.guide).Content)}
	if findings := stackWordingFindings(rendered, []stackWordingRule{rule}); len(findings) != 0 {
		t.Fatalf("stack wording findings: %v", findings)
	}
}

func TestStackWordingCheckReportsTheWordingItReplaced(t *testing.T) {
	t.Parallel()
	const oldText = "Use Bun-owned commands for dependency installation, scripts, tests, and lockfile updates.\n" +
		"Do not substitute another package manager or hand-edit the lockfile.\n" +
		"Keep type errors visible and preserve the repository's package and lockfile workflow.\n" +
		"When the selected TypeScript/Bun profile treats warnings as errors, every warning reported by Verification blocks completion."
	rule := typeScriptAndBunWordingRule()
	t.Run("old wording", func(t *testing.T) {
		findings := stackWordingFindings(map[string]string{rule.guide: oldText}, []stackWordingRule{rule})
		var want []string
		for _, sentence := range rule.must {
			want = append(want, fmt.Sprintf("%s: missing sentence: %s", rule.guide, sentence))
		}
		for _, sentence := range rule.mustNot {
			want = append(want, fmt.Sprintf("%s: replaced sentence remains: %s", rule.guide, sentence))
		}
		if !reflect.DeepEqual(findings, want) {
			t.Fatalf("findings = %v, want %v", findings, want)
		}
	})
	t.Run("absent guide", func(t *testing.T) {
		findings := stackWordingFindings(nil, []stackWordingRule{rule})
		want := []string{rule.guide + ": guide is not rendered"}
		if !reflect.DeepEqual(findings, want) {
			t.Fatalf("findings = %v, want %v", findings, want)
		}
	})
	t.Run("wrapped sentences", func(t *testing.T) {
		text := strings.ReplaceAll(strings.Join(rule.must, "\n"), " ", "\n\t")
		if findings := stackWordingFindings(map[string]string{rule.guide: text}, []stackWordingRule{rule}); len(findings) != 0 {
			t.Fatalf("wrapped wording findings: %v", findings)
		}
	})
}

func TestTheRewordedBunAndTypeScriptClausesKeepTheirForce(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	for _, test := range []struct{ module, clause, force string }{
		{"bun", "clause.bun.use-bun-owned-commands", "mandatory"},
		{"bun", "clause.bun.prohibit-other-package-managers", "prohibited"},
		{"typescript", "clause.typescript.keep-type-errors-visible", "mandatory"},
	} {
		t.Run(test.clause, func(t *testing.T) {
			module, ok := catalog.Module(test.module)
			if !ok {
				t.Fatalf("module %s is absent", test.module)
			}
			var decoded map[string]any
			if err := json.Unmarshal(module.Data, &decoded); err != nil {
				t.Fatal(err)
			}
			force, found := clauseForce(decoded, test.clause)
			if !found || force != test.force {
				t.Fatalf("clause force = (%q, %t), want (%q, true)", force, found, test.force)
			}
			if force, found := clauseForce(decoded, "clause.does-not-exist"); found || force != "" {
				t.Fatalf("absent clause force = (%q, %t)", force, found)
			}
		})
	}
}
