package baseline

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

const cancelClauseID = "clause.spec.keep-artifacts-in-spec-folder"

const cancelClauseGuidance = "A Spec dropped before delivery is cancelled, never deleted by hand: a runtime cancel command, such as `roundfix cancel <slug> --reason <text>`, writes an Archive Record with the disposition `cancelled` that keeps the reason, the date and the Tasks completed and pending at that moment, returns each adopted Backlog Entry to `docs/backlog/` with `status: open` and each adopted Finding to `docs/findings/` with `status: pending`, and removes the Spec folder in the same change. It refuses while a Run or a delivery queue item is active for the Spec. A cancelled Spec never meets a prerequisite and never supersedes another Spec. A runtime without a cancel command must report cancellation as unsupported rather than deleting the folder."

func TestTheCancelClauseCarriesItsText(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	asset, ok := catalog.Module("spec-workflow")
	if !ok {
		t.Fatal("missing spec-workflow module")
	}
	var module document
	if err := json.Unmarshal(asset.Data, &module); err != nil {
		t.Fatal(err)
	}
	for _, rule := range objectsOrEmpty(module["rules"]) {
		for _, clause := range objectsOrEmpty(rule["clauses"]) {
			if clause["id"] != cancelClauseID {
				continue
			}
			if clause["enforcement"] != "mandatory" || strings.Count(clause["guidance"].(string), cancelClauseGuidance) != 1 {
				t.Fatalf("cancel clause force/text differs: %+v", clause)
			}
			if _, exists := clause["replaces"]; exists {
				t.Fatal("cancel clause has replaces")
			}
			return
		}
	}
	t.Fatalf("missing clause %s", cancelClauseID)
}

func TestTheCancelClauseRendersInTheDocsLayoutGuide(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	asset, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md")
	if !ok {
		t.Fatal("missing Standard TypeScript docs-layout golden")
	}
	local, err := os.ReadFile("../../docs/agents/docs-layout.md")
	if err != nil {
		t.Fatal(err)
	}
	literal := strings.Join(strings.Fields(cancelClauseGuidance), " ")
	for name, data := range map[string][]byte{"Standard TypeScript golden": asset.Data, "repository guide": local} {
		t.Run(name, func(t *testing.T) {
			text := strings.Join(strings.Fields(string(data)), " ")
			if count := strings.Count(text, literal); count != 1 {
				t.Fatalf("guide contains literal %d times, want exactly once", count)
			}
		})
	}
}

func TestAnAdopterRetainsTheCancelClause(t *testing.T) {
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
	if delta.Dispositions[cancelClauseID] != ClauseRetained {
		t.Errorf("disposition = %s, want retained", delta.Dispositions[cancelClauseID])
	}
	for _, evidence := range outcome.Plan.Retention {
		if evidence.FromClause == cancelClauseID {
			if evidence.Disposition != string(ClauseRetained) {
				t.Errorf("retention evidence = %+v", evidence)
			}
			return
		}
	}
	t.Error("missing retention evidence")
}
