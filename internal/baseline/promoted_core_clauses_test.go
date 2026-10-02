// Suite: Promoted core clauses
// Invariant: promoted rules render with force and declared replacements preserve retention.
// Boundary IN: embedded catalog and plans for temporary adopters.
// Boundary OUT: repository Verification and applying a production refresh.

package baseline

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type expectedClause struct{ module, id, force string }

func promotedCoreClauses() []expectedClause {
	return []expectedClause{
		{"core", "clause.core.regenerate-generated-files", "prohibited"},
		{"core", "clause.core.flaky-tests-block", "mandatory"},
		{"core", "clause.core.lint-warnings-block", "mandatory"},
		{"core", "clause.core.prohibit-test-only-production-hooks", "prohibited"},
		{"core", "clause.core.prohibit-editing-vendored-skills", "prohibited"},
		{"core", "clause.core.ask-before-database-mutation", "stop-and-ask"},
		{"core", "clause.core.prove-the-mutation-predicate", "mandatory"},
		{"core", "clause.core.prohibit-disguised-database-mutation", "prohibited"},
	}
}

func clauseForceFindings(catalog *Catalog, want []expectedClause) []string {
	var findings []string
	for _, clause := range want {
		asset, ok := catalog.Module(clause.module)
		if !ok {
			findings = append(findings, "missing module "+clause.module)
			continue
		}
		var module map[string]any
		if err := json.Unmarshal(asset.Data, &module); err != nil {
			findings = append(findings, fmt.Sprintf("module %s: %v", clause.module, err))
			continue
		}
		force, found := clauseForce(module, clause.id)
		if !found {
			findings = append(findings, "missing clause "+clause.id)
		} else if force != clause.force {
			findings = append(findings, fmt.Sprintf("clause %s force = %s, want %s", clause.id, force, clause.force))
		}
	}
	return findings
}

func TestTheCoreGuidesStateThePromotedRules(t *testing.T) {
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
				{guide: "docs/agents/agent-instructions.md", must: []string{
					"- **prohibited**: Do not hand-edit a generated file, such as a schema migration, a client generated from an API description, or a derived digest; change its source and run the generator that owns it.",
					"- **mandatory**: A flaky test is a blocking failure: find and fix the cause of its nondeterminism before completion, and never rerun Verification until it happens to pass.",
					"- **mandatory**: A lint warning fails Verification. Every linter the repository Verification runs uses its option that turns a warning into a failure, such as `oxlint --deny-warnings`, `biome check --error-on-warnings`, `eslint --max-warnings 0`, or `cargo clippy -- -D warnings`, and every warning Verification reports blocks completion.",
					"- **prohibited**: Do not add a production hook, flag, branch, or exported symbol that exists only for tests; test through the public entry points.",
					"- **stop-and-ask**: Stop and ask for express authorization before any statement that changes data or schema in a database other than a disposable local one. Name the exact statement and database; an authorization covers only that statement.",
					"- **mandatory**: Before an authorized database write, run a read with the same predicate and report its row count, then run the write inside an explicit transaction and confirm the affected row count before committing.",
					"- **prohibited**: Do not route a database change through a migration, script, seed, or test to avoid that authorization. The repository's committed migrate and seed commands run against a disposable local database are exempt.",
				}},
				{guide: "docs/agents/skill-dispatch.md", must: []string{
					"- **prohibited**: Do not edit a skill the repository installs from an upstream source; its upstream owns its text. State a correction as a Repository-Specific Normative Rule, which governs over the skill's default.",
				}},
			}
			if profile == "standard-typescript-monorepo" {
				rules = append(rules, stackWordingRule{
					guide:   "docs/agents/typescript-bun.md",
					mustNot: []string{"When the repository's Verification treats warnings as errors"},
				})
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

func TestThePromotedCoreClausesCarryTheirForce(t *testing.T) {
	t.Parallel()
	if findings := clauseForceFindings(mustEmbeddedCatalog(t), promotedCoreClauses()); len(findings) != 0 {
		t.Fatalf("promoted force findings: %v", findings)
	}
}

func TestAClauseForceCheckReportsAMissingOrChangedClause(t *testing.T) {
	t.Parallel()
	got := clauseForceFindings(mustEmbeddedCatalog(t), []expectedClause{
		{"core", "clause.does-not-exist", "mandatory"},
		{"core", "clause.core.fix-root-causes", "prohibited"},
	})
	want := []string{
		"missing clause clause.does-not-exist",
		"clause clause.core.fix-root-causes force = mandatory, want prohibited",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("force findings = %v, want %v", got, want)
	}
}

func TestTheWarningsClauseIsReplacedByTheLintClause(t *testing.T) {
	t.Parallel()
	const removed = "clause.bun.block-warnings-when-profile-treats-them-as-errors"
	const successor = "clause.core.lint-warnings-block"
	request, catalog := newClauseReplacementAdopter(t)
	outcome, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil || outcome.Result.State != "ready" {
		t.Fatalf("declared warnings replacement refused refresh: %+v", outcome.Result)
	}
	if delta := outcome.Plan.ClauseDelta; delta == nil || delta.Dispositions[removed] != ClauseReplaced {
		t.Fatalf("warnings replacement delta = %+v", delta)
	}
	found := false
	for _, evidence := range outcome.Plan.Retention {
		if evidence.FromClause != removed {
			continue
		}
		found = true
		if evidence.Disposition != string(ClauseReplaced) || !reflect.DeepEqual(evidence.Targets, []string{successor}) {
			t.Fatalf("warnings replacement evidence = %+v", evidence)
		}
	}
	if !found {
		t.Fatal("removed warnings clause has no retention evidence")
	}
}

func TestAnUndeclaredWarningsReplacementIsUnaccounted(t *testing.T) {
	t.Parallel()
	const removed = "clause.bun.block-warnings-when-profile-treats-them-as-errors"
	request, catalog := newClauseReplacementAdopter(t)
	mutateCatalogClause(t, catalog, "core", "clause.core.lint-warnings-block", func(clause document) bool {
		delete(clause, "replaces")
		return true
	})
	outcome, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan != nil || outcome.Result.State != "action_required" || outcome.Result.Category != "classification" {
		t.Fatalf("undeclared warnings replacement outcome = %+v", outcome)
	}
	if delta := outcome.Result.ClauseDelta; delta == nil || delta.Dispositions[removed] != ClauseUnaccounted {
		t.Fatalf("undeclared warnings replacement delta = %+v", delta)
	}
	if !strings.Contains(outcome.Result.Message, removed) {
		t.Fatalf("refusal does not name removed warnings clause: %q", outcome.Result.Message)
	}
}
