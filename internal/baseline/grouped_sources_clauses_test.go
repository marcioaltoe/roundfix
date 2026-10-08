package baseline

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// These literals lock the authored contract, independently of the module assets.
var groupingClauses = []struct{ id, guidance, module, rule, guide string }{
	{"clause.spec.sources-01-group-by-shared-context", "One Spec may adopt several Inbox Entries, Backlog Entries, and Findings whose context is similar or complementary: they change the same component or contract, share a root cause, or one completes the other. One Spec per source is neither required nor preferred, so before minting a Spec for a source, look among the open Backlog Entries and unresolved Findings for others that share its context. A grouping suggestion from `roundfix spec judge` is advisory: adopt the suggested source or state why it stays apart.", "spec-workflow", "rule.spec.routing", "spec-routing.md"},
	{"clause.spec.sources-02-bound-the-group", "Group sources into one Spec only while the Spec fits four implementation Tasks plus its QA gate. When the grouped scope needs more, split it into Specs that each fit and give each source exactly one owning Spec; never add a source that shares no context with the Spec only to save a Spec. When a Spec's Verification or premise depends on another Spec that has not merged, name that Spec in the Task Graph manifest's `requires` and state the dependency in the TechSpec Build Order.", "spec-workflow", "rule.spec.routing", "spec-routing.md"},
	{"clause.context.inbox-03-extend-before-minting", "Before minting a Finding or Backlog Entry, including in Triage, look for an existing one that is not yet implemented and that the new observation or intent fits: an `open` Backlog Entry in `docs/backlog/` or an unresolved Finding in `docs/findings/`. Revise a fitting Backlog Entry in place under its existing name; extend a fitting Finding with a dated addendum, because a Finding is immutable history. Mint a new file only when none fits. A Triage that extends an existing entry resolves its Inbox Entry into that entry, which cites the Inbox Entry's provenance.", "context-workflow", "rule.context.docs-layout", "docs-layout.md"},
}

func TestTheGroupingClausesCarryTheirForceAndText(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	for _, want := range groupingClauses {
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
						t.Error("new clause has replaces")
					}
				}
			}
			if !found {
				t.Fatalf("missing clause %s in %s", want.id, want.rule)
			}
		})
	}
}

func TestTheGroupingClausesRenderInTheGuides(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	for _, want := range groupingClauses {
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

func TestTheGroupingClausesHaveSourceBaselineRows(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	source, err := catalog.SourceBaseline("baseline.standard-typescript-monorepo-0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range groupingClauses {
		t.Run(want.id, func(t *testing.T) {
			found := false
			for _, entry := range source.Entries {
				if entry.ID != want.id {
					continue
				}
				found = true
				if entry.Kind != "normative-clause" || entry.Enforcement != "mandatory" || entry.Carrier != "docs/agents/"+want.guide || entry.Path != "corpus/docs/agents/"+want.guide {
					t.Fatalf("Source Baseline row differs: %+v", entry)
				}
			}
			if !found {
				t.Fatalf("missing Source Baseline row %s", want.id)
			}
		})
	}
}

func TestAnAdopterRetainsTheGroupingClauses(t *testing.T) {
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
	for _, want := range groupingClauses {
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
