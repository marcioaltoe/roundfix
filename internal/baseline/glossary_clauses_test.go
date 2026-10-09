package baseline

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// These literals lock the authored contract, independently of the module assets.
var glossaryClauses = []struct{ id, guidance, module, rule, guide string }{
	{"clause.context.read-domain-contract", "The repository's selected domain context and its accepted ADRs are the base of the CONTEXT-driven workflow. Read both before naming domain concepts or changing behavior, and flag conflicts instead of silently overriding repository decisions.", "context-workflow", "rule.context.domain-docs", "domain.md"},
	{"clause.domain.glossary-currency", "A Spec adds each domain term it introduces to the domain context, and revises each term it changes, through `domain-modeling` in one of its own Tasks, whose Verification names the term. It lists those terms in a `## Glossary` section of its PRD or TechSpec, with any bolded phrase it declares not a domain term. `roundfix spec check` reports a bolded term the domain context lacks and the section does not cover, and the QA gate and `roundfix archive` refuse a Spec whose declared term is still missing. When `domain-modeling` names its glossary file differently, it writes to the selected domain context. Work outside a Spec checks at its close whether it introduced, changed, or retired a term and updates the domain context when it did. Neither the check nor the update waits for human interaction; reach for `grilling` only when a term is ambiguous enough to need sharpening before it is written down.", "context-workflow", "rule.context.domain-docs", "domain.md"},
}

func TestTheGlossaryClausesCarryTheirForceAndText(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	for _, want := range glossaryClauses {
		t.Run(want.id, func(t *testing.T) {
			asset, ok := catalog.Module(want.module)
			if !ok {
				t.Fatalf("missing module %s", want.module)
			}
			var module document
			if err := json.Unmarshal(asset.Data, &module); err != nil {
				t.Fatal(err)
			}
			found := false
			for _, rule := range objectsOrEmpty(module["rules"]) {
				if rule["id"] != want.rule {
					continue
				}
				for _, clause := range objectsOrEmpty(rule["clauses"]) {
					if clause["id"] != want.id {
						continue
					}
					found = true
					if clause["enforcement"] != "mandatory" || clause["guidance"] != want.guidance {
						t.Errorf("clause %s force/text differs: %+v", want.id, clause)
					}
					if _, exists := clause["replaces"]; exists {
						t.Error("reworded clause has replaces")
					}
				}
			}
			if !found {
				t.Fatalf("missing clause %s in %s", want.id, want.rule)
			}
		})
	}
}

func TestTheGlossaryClausesRenderInTheDomainGuide(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	for _, want := range glossaryClauses {
		t.Run(want.id, func(t *testing.T) {
			asset, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/" + want.guide)
			if !ok {
				t.Fatalf("missing guide %s", want.guide)
			}
			line := "\n- **mandatory**: " + want.guidance + "\n"
			if strings.Count("\n"+string(asset.Data)+"\n", line) != 1 {
				t.Fatalf("guide lacks exactly one forced clause %s", want.id)
			}
		})
	}
}

func TestAnAdopterRetainsTheGlossaryClauses(t *testing.T) {
	t.Parallel()
	request, catalog := newClauseReplacementAdopter(t)
	outcome, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil || outcome.Result.State != "ready" {
		t.Fatalf("refresh is not ready: %+v", outcome.Result)
	}
	delta := outcome.Plan.ClauseDelta
	if delta == nil {
		t.Fatal("refresh has no clause accounting")
	}
	for id, disposition := range delta.Dispositions {
		if disposition == ClauseUnaccounted {
			t.Errorf("unaccounted clause %s", id)
		}
	}
	for _, want := range glossaryClauses {
		t.Run(want.id, func(t *testing.T) {
			if delta.Dispositions[want.id] != ClauseRetained {
				t.Errorf("disposition = %s, want retained", delta.Dispositions[want.id])
			}
			found := false
			for _, evidence := range outcome.Plan.Retention {
				if evidence.FromClause == want.id {
					found = true
					if evidence.Disposition != string(ClauseRetained) {
						t.Errorf("retention evidence = %+v", evidence)
					}
				}
			}
			if !found {
				t.Error("missing retention evidence")
			}
		})
	}
}
