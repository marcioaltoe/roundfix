package baseline

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestAnUnrecordedFrontendLayoutStatesTheSuggestionAndKeepsTheSystemsClauses(t *testing.T) {
	assertFrontendLayoutPlan(t, "", "No frontend layout is recorded. The suggested `systems` layout applies until the repository records one.", true)
	// A repository-owned profile need not select the new optional decision.
	catalog := mustEmbeddedCatalog(t)
	profile, err := ResolveProfile("", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatal(err)
	}
	var selected []string
	for _, id := range profile.Decisions {
		if id != frontendLayoutDecisionID {
			selected = append(selected, id)
		}
	}
	profile.Decisions = selected
	_, artifacts, err := resolveManagedArtifacts(catalog, profile, standardTypeScriptDecisions("make verify"), false)
	if err != nil {
		t.Fatal(err)
	}
	for _, artifact := range artifacts {
		if artifact.ID == "guide.frontend" {
			if !strings.Contains(string(artifact.Body), "No frontend layout is recorded.") {
				t.Fatalf("repository-owned profile guidance = %s", artifact.Body)
			}
			return
		}
	}
	t.Fatal("repository-owned profile lost frontend guide")
}

func TestARecordedSystemsLayoutKeepsTheSystemsClauses(t *testing.T) {
	assertFrontendLayoutPlan(t, "systems", "The repository records the `systems` frontend layout.", true)
}

func TestARecordedRepositoryDefinedLayoutRendersOnlyItsOwnClause(t *testing.T) {
	assertFrontendLayoutPlan(t, "repository-defined", "The repository records its own frontend layout, stated in its repository-owned rules.", false)
}

func assertFrontendLayoutPlan(t *testing.T, layout, sentence string, systems bool) {
	t.Helper()
	decisions := standardTypeScriptDecisions("make verify")
	if layout != "" {
		decisions = append(decisions, DecisionValue{ID: frontendLayoutDecisionID, Value: layout})
	}
	plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), decisions)
	guide := string(planPostimage(t, plan, "docs/agents/frontend.md").Content)
	if !strings.Contains(guide, sentence) {
		t.Fatalf("frontend guide missing %q: %s", sentence, guide)
	}
	catalog := mustEmbeddedCatalog(t)
	count := 0
	for _, rule := range objectsOrEmpty(catalog.modules["frontend"]["rules"]) {
		for _, clause := range objectsOrEmpty(rule["clauses"]) {
			id, _ := stringValue(clause, "id")
			guidance, _ := stringValue(clause, "guidance")
			want := true
			switch id {
			case "clause.frontend.organize-by-system", "clause.frontend.public-system-boundary":
				want = systems
			case "clause.frontend.follow-recorded-layout":
				want = !systems
			}
			got := strings.Contains(guide, guidance)
			if got != want {
				t.Errorf("clause %s rendered = %v, want %v", id, got, want)
			}
			if got {
				count++
			}
		}
	}
	wantCount := 6
	if !systems {
		wantCount = 5
	}
	if count != wantCount || strings.Count(guide, "- **") != wantCount {
		t.Fatalf("rendered clauses = %d, want %d: %s", count, wantCount, guide)
	}
	if strings.Contains(guide, "{{") {
		t.Fatal("unexpanded token")
	}
}

func TestAnUnrecordedOptionalDecisionIsNotMissing(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	profile, err := ResolveProfile("", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatal(err)
	}
	decisions, missing, err := normalizePlanDecisions(profile, standardTypeScriptDecisions("make verify"), catalog)
	if err != nil || len(missing) != 0 {
		t.Fatalf("normalization missing=%v err=%v", missing, err)
	}
	for _, decision := range decisions {
		if decision.ID == frontendLayoutDecisionID {
			t.Fatal("unanswered optional decision was recorded")
		}
	}
	_, divergences, err := normalizeAlignmentDecisions(profile, standardTypeScriptDecisions("make verify"), catalog)
	if err != nil || len(divergences) != 0 {
		t.Fatalf("alignment divergences=%v err=%v", divergences, err)
	}
	var unansweredRequired []DecisionValue
	for _, decision := range standardTypeScriptDecisions("make verify") {
		if decision.ID != "verification.gate" {
			unansweredRequired = append(unansweredRequired, decision)
		}
	}
	_, divergences, err = normalizeAlignmentDecisions(profile, unansweredRequired, catalog)
	if err != nil || len(divergences) != 1 || divergences[0].ID != "verification.gate" || divergences[0].Code != "profile.decision.required" {
		t.Fatalf("required decision divergences=%v err=%v", divergences, err)
	}
}

func TestAnOptionalDecisionWithoutADefaultIsRefused(t *testing.T) {
	assets := cloneEmbeddedAssets(t)
	replaceAsset(t, assets, "decisions.json", "\"default\": \"systems\",", "")
	_, err := LoadCatalog(assets)
	assertFrontendDiagnostic(t, err, "catalog.decision.optional.invalid")
}

func TestAClauseGateOnARequiredDecisionIsRefused(t *testing.T) {
	assets := cloneEmbeddedAssets(t)
	replaceAsset(t, assets, "modules/frontend.json", `"id": "clause.frontend.organize-by-system",
          "enforcement": "mandatory",
          "appliesWhen": {"decision": "frontend.layout", "equals": "systems"}`, `"id": "clause.frontend.organize-by-system",
          "enforcement": "mandatory",
          "appliesWhen": {"decision": "domain.layout", "equals": "single-context"}`)
	_, err := LoadCatalog(assets)
	assertFrontendDiagnostic(t, err, "catalog.clause.applies-when.invalid")
}

func assertFrontendDiagnostic(t *testing.T, err error, code string) {
	t.Helper()
	var validation *ValidationError
	if !errors.As(err, &validation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	for _, diagnostic := range validation.Diagnostics {
		if diagnostic.Code == code {
			return
		}
	}
	t.Fatalf("missing diagnostic %s: %v", code, err)
}

func TestAClauseARecordedLayoutTurnsOffIsAReasonedRejection(t *testing.T) {
	repository, prior := frontendLayoutManagedRepository(t)
	decisions := append(standardTypeScriptDecisions("make verify"), DecisionValue{ID: frontendLayoutDecisionID, Value: "repository-defined"})
	// Preserve both prior systems clauses' artifacts while recording the new value.
	prior.Decisions[frontendLayoutDecisionID] = ManifestDecision{Value: "repository-defined"}
	data, err := json.Marshal(prior)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repository, manifestPath), data, 0o644); err != nil {
		t.Fatal(err)
	}
	outcome, err := BuildPlan(context.Background(), PlanRequest{Repository: repository, ProfileID: "standard-typescript-monorepo", Decisions: decisions, Preservation: RootPreservationRequest{Mode: PreservationModeManagedRefresh}})
	if err != nil || outcome.Plan == nil {
		t.Fatalf("layout transition: %v %+v", err, outcome.Result)
	}
	for _, id := range []string{"clause.frontend.organize-by-system", "clause.frontend.public-system-boundary"} {
		found := false
		for _, evidence := range outcome.Plan.Retention {
			if evidence.FromClause != id {
				continue
			}
			found = true
			if evidence.Disposition != "reasoned-rejection" || !reflect.DeepEqual(evidence.Targets, []string{frontendLayoutDecisionID}) || evidence.Reason != "The repository recorded frontend.layout = repository-defined; this clause applies only when it is systems." {
				t.Errorf("retention %s = %+v", id, evidence)
			}
		}
		if !found {
			t.Errorf("no retention for %s", id)
		}
	}
	if outcome.Plan.ClauseDelta == nil || outcome.Plan.ClauseDelta.Counts[ClauseUnaccounted] != 0 {
		t.Fatalf("delta = %+v", outcome.Plan.ClauseDelta)
	}
}

func TestAClauseMissingWithoutARecordedDecisionIsStillUnaccounted(t *testing.T) {
	repository, _ := frontendLayoutManagedRepository(t)
	catalog, err := LoadCatalog(cloneEmbeddedAssets(t))
	if err != nil {
		t.Fatal(err)
	}
	// Remove the catalog clause, not merely its rendering, with no recorded gate reason.
	for _, rule := range objectsOrEmpty(catalog.modules["frontend"]["rules"]) {
		var retained []any
		for _, clause := range objectsOrEmpty(rule["clauses"]) {
			if clause["id"] != "clause.frontend.organize-by-system" {
				retained = append(retained, map[string]any(clause))
			}
		}
		rule["clauses"] = retained
	}
	outcome, err := buildPlanWithCatalog(context.Background(), PlanRequest{Repository: repository, ProfileID: "standard-typescript-monorepo", Decisions: standardTypeScriptDecisions("make verify"), Preservation: RootPreservationRequest{Mode: PreservationModeManagedRefresh}}, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan != nil {
		t.Fatal("unaccounted removal produced an applicable plan")
	}
	if !strings.Contains(outcome.Result.Message, "unaccounted") {
		t.Fatalf("result = %+v", outcome.Result)
	}
	if outcome.Result.ClauseDelta == nil || outcome.Result.ClauseDelta.Dispositions["clause.frontend.organize-by-system"] != ClauseUnaccounted {
		t.Fatalf("unaccounted delta = %+v", outcome.Result.ClauseDelta)
	}
}

func frontendLayoutManagedRepository(t *testing.T) (string, SetupManifest) {
	t.Helper()
	repository := newProjectDecisionPlanRepository(t)
	plan := buildProjectDecisionPlan(t, repository, standardTypeScriptDecisions("make verify"))
	if _, err := ApplyPlan(context.Background(), repository, plan, plan.PlanDigest); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(repository, manifestPath))
	if err != nil {
		t.Fatal(err)
	}
	var manifest SetupManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	return repository, manifest
}

func TestTheProfileStatesTheLayoutSuggestionTheCatalogDefaults(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	architecture, _ := objectValue(catalog.profiles["standard-typescript-monorepo"]["architecture"])
	frontend, _ := objectValue(architecture["frontend"])
	suggested, _ := objectValue(frontend["suggested"])
	if frontend["layoutDecision"] != frontendLayoutDecisionID || suggested["organization"] != catalog.decisions[frontendLayoutDecisionID]["default"] || suggested["publicBoundary"] != "public system boundary" || suggested["internalImports"] != "direct" {
		t.Fatalf("frontend architecture = %+v", frontend)
	}
	ids := stringsOrEmpty(catalog.profiles["standard-typescript-monorepo"]["entryDecisions"])
	for index, id := range ids {
		if id == frontendLayoutDecisionID {
			if index == 0 || index+1 >= len(ids) || ids[index-1] != "http.contract" || ids[index+1] != "auth.provider" {
				t.Fatalf("entry decisions = %v", ids)
			}
			return
		}
	}
	t.Fatal("profile does not select frontend layout")
}
