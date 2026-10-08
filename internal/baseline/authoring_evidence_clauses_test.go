package baseline

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// These literals lock the authored contract, independently of the module assets.
var authoringEvidenceClauses = []struct{ id, guidance, module, rule, guide string }{
	{"clause.spec.project-constraints-06-outside-evidence", "Name a source the QA gate can check without the network before the Run starts: a document the planning change commits under `docs/references/`, or published literature whose cited passage the row quotes. Never rest the row on an artifact that exists only on one machine, on another repository's manifest, or on a lookup the Spec's authorization forbids, and declare a criterion that needs the network or an open Pull Request under Unreachable Acceptance when the row is authored.", "spec-workflow", "rule.spec.project-constraints", "spec-routing.md"},
	{"clause.spec.routing-05-task-graph", "The Task that implements a Surface Transcript asserts the whole transcript in one named test: every line of its standard output and standard error, in order, and its exit code, copied from the transcript with only its declared controlled variations; the Task's Verification runs that test by name.", "spec-workflow", "rule.spec.routing", "spec-routing.md"},
	{"clause.spec.sources-02-bound-the-group", "When a Spec's Verification or premise depends on another Spec that has not merged, name that Spec in the Task Graph manifest's `requires` and state the dependency in the TechSpec Build Order.", "spec-workflow", "rule.spec.routing", "spec-routing.md"},
}

func TestTheAuthoringEvidenceClausesCarryTheirText(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	for _, want := range authoringEvidenceClauses {
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
					if clause["enforcement"] != "mandatory" || strings.Count(clause["guidance"].(string), want.guidance) != 1 {
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

func TestTheAuthoringEvidenceClausesRenderInTheGuide(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	for _, want := range authoringEvidenceClauses {
		t.Run(want.id, func(t *testing.T) {
			asset, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/" + want.guide)
			if !ok {
				t.Fatalf("missing guide %s", want.guide)
			}
			local, err := os.ReadFile("../../docs/agents/" + want.guide)
			if err != nil {
				t.Fatal(err)
			}
			for name, data := range map[string][]byte{"Standard TypeScript golden": asset.Data, "repository guide": local} {
				t.Run(name, func(t *testing.T) {
					text := strings.Join(strings.Fields(string(data)), " ")
					literal := strings.Join(strings.Fields(want.guidance), " ")
					if count := strings.Count(text, literal); count != 1 {
						t.Fatalf("guide contains literal %d times, want exactly once", count)
					}
				})
			}
		})
	}
}

func TestAnAdopterRetainsTheAuthoringEvidenceClauses(t *testing.T) {
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
	for _, want := range authoringEvidenceClauses {
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
