// Suite: Promoted Spec-workflow and TypeScript clauses
// Invariant: promoted rules render with force and retain their Source Baseline rows.
// Boundary IN: embedded catalog and plans for temporary adopters.
// Boundary OUT: repository Verification and applying a production refresh.

package baseline

import (
	"fmt"
	"reflect"
	"testing"
)

func promotedSpecAndTypeScriptClauses() []expectedClause {
	return []expectedClause{
		{"spec-workflow", "clause.spec.name-the-finding-a-task-answers", "mandatory"},
		{"spec-workflow", "clause.spec.verification-fails-before-the-change", "mandatory"},
		{"spec-workflow", "clause.spec.prohibit-tests-that-read-specs", "prohibited"},
		{"typescript", "clause.typescript.type-fixtures-from-the-schema", "mandatory"},
	}
}

func sourceBaselineRowFindings(source SourceBaseline, want []expectedClause) []string {
	rows := make(map[string]SourceBaselineEntry, len(source.Entries))
	for _, row := range source.Entries {
		rows[row.ID] = row
	}
	var findings []string
	for _, clause := range want {
		row, ok := rows[clause.id]
		if !ok {
			findings = append(findings, "missing Source Baseline row "+clause.id)
			continue
		}
		if row.Kind != "normative-clause" {
			findings = append(findings, fmt.Sprintf("row %s kind = %s, want normative-clause", clause.id, row.Kind))
		}
		if row.Enforcement != clause.force {
			findings = append(findings, fmt.Sprintf("row %s force = %s, want %s", clause.id, row.Enforcement, clause.force))
		}
	}
	return findings
}

func TestTheSpecAndTypeScriptGuidesStateThePromotedRules(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"standard-typescript-monorepo", "go-cli-tui"} {
		t.Run(profile, func(t *testing.T) {
			var plan PlanDocument
			if profile == "standard-typescript-monorepo" {
				plan = buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), standardTypeScriptDecisions("make verify"))
			} else {
				plan = buildTestPlan(t, newPlanRepository(t))
			}
			rules := []stackWordingRule{
				{guide: "docs/agents/spec-routing.md", must: []string{
					"- **mandatory**: When a Task answers a recorded finding, name that finding and its date in one of the Task's requirements.",
					"- **mandatory**: Author each Task's Verification so every command fails on the tree before the Task's change, and run `roundfix spec check <slug> --run-verification` before a Run starts: the Daemon refuses a Task whose command already exits zero on the unchanged tree.",
				}},
				{guide: "docs/agents/docs-layout.md", must: []string{
					"- **prohibited**: Do not make a test, fixture, or build step read a file under the Spec root or assume that a Spec is still active; Specs archive and may be deleted. A check whose subject is the Spec artifacts themselves is exempt.",
				}},
			}
			if profile == "standard-typescript-monorepo" {
				rules = append(rules, stackWordingRule{guide: "docs/agents/typescript-bun.md", must: []string{
					"- **mandatory**: Type a test fixture that stands for a stored row from the schema's inferred row type, with a complete base row and partial overrides, never from an untyped record.",
				}})
			}
			rendered := make(map[string]string)
			for _, rule := range rules {
				rendered[rule.guide] = string(planPostimage(t, plan, rule.guide).Content)
			}
			if findings := stackWordingFindings(rendered, rules); len(findings) != 0 {
				t.Fatalf("promoted wording findings: %v", findings)
			}
		})
	}
}

func TestThePromotedSpecAndTypeScriptClausesCarryTheirForce(t *testing.T) {
	t.Parallel()
	if findings := clauseForceFindings(mustEmbeddedCatalog(t), promotedSpecAndTypeScriptClauses()); len(findings) != 0 {
		t.Fatalf("promoted force findings: %v", findings)
	}
}

func TestEveryPromotedClauseHasItsSourceBaselineRow(t *testing.T) {
	t.Parallel()
	source, err := mustEmbeddedCatalog(t).SourceBaseline("baseline.standard-typescript-monorepo-0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	want := append(promotedCoreClauses(), promotedSpecAndTypeScriptClauses()...)
	if findings := sourceBaselineRowFindings(source, want); len(findings) != 0 {
		t.Fatalf("promoted Source Baseline findings: %v", findings)
	}
}

func TestAMissingSourceBaselineRowIsReported(t *testing.T) {
	t.Parallel()
	const id = "clause.typescript.type-fixtures-from-the-schema"
	wantClauses := []expectedClause{{"typescript", id, "mandatory"}}
	for _, tc := range []struct {
		name   string
		source SourceBaseline
		want   []string
	}{
		{name: "missing", source: SourceBaseline{}, want: []string{"missing Source Baseline row " + id}},
		{name: "wrong force", source: SourceBaseline{Entries: []SourceBaselineEntry{{ID: id, Kind: "normative-clause", Enforcement: "prohibited"}}}, want: []string{"row " + id + " force = prohibited, want mandatory"}},
		{name: "wrong kind", source: SourceBaseline{Entries: []SourceBaselineEntry{{ID: id, Kind: "recommendation", Enforcement: "mandatory"}}}, want: []string{"row " + id + " kind = recommendation, want normative-clause"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := sourceBaselineRowFindings(tc.source, wantClauses); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Source Baseline findings = %v, want %v", got, tc.want)
			}
		})
	}
}
