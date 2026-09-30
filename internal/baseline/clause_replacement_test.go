package baseline

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

const removedBackendClause = "rule.backend.boundary-contracts"
const successorBackendClause = "clause.backend.boundary-contracts"

func TestARemovedSourceClauseIsReplacedByItsDeclaredSuccessor(t *testing.T) {
	t.Parallel()
	request, catalog := newClauseReplacementAdopter(t)
	outcome, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil || outcome.Result.State != "ready" {
		t.Fatalf("declared replacement refused adopter refresh: %+v", outcome.Result)
	}
	if len(outcome.Plan.FileChanges) == 0 {
		t.Fatal("drifted adopter refresh has no file changes")
	}
	if delta := outcome.Plan.ClauseDelta; delta == nil || delta.Dispositions[removedBackendClause] != ClauseReplaced {
		t.Fatalf("replacement delta = %+v", delta)
	}
	found := false
	for _, evidence := range outcome.Plan.Retention {
		if evidence.FromClause != removedBackendClause {
			continue
		}
		found = true
		if len(evidence.Targets) != 1 || evidence.Targets[0] != successorBackendClause || !strings.Contains(evidence.Reason, successorBackendClause) {
			t.Fatalf("replacement evidence = %+v", evidence)
		}
	}
	if !found {
		t.Fatal("removed clause has no retention evidence")
	}
}

func TestAnUndeclaredRemovalStaysUnaccounted(t *testing.T) {
	t.Parallel()
	request, catalog := newClauseReplacementAdopter(t)
	mutateCatalogClause(t, catalog, "backend", successorBackendClause, func(clause document) bool {
		delete(clause, "replaces")
		return true
	})
	assertClauseReplacementRefused(t, request, catalog)
}

func TestAReplacementWithOtherForceStaysUnaccounted(t *testing.T) {
	t.Parallel()
	request, catalog := newClauseReplacementAdopter(t)
	mutateCatalogClause(t, catalog, "backend", successorBackendClause, func(clause document) bool {
		clause["enforcement"] = "stop-and-ask"
		return true
	})
	assertClauseReplacementRefused(t, request, catalog)
}

func assertClauseReplacementRefused(t *testing.T, request PlanRequest, catalog *Catalog) {
	t.Helper()
	outcome, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan != nil || outcome.Result.State != "action_required" || outcome.Result.Category != "classification" {
		t.Fatalf("unaccounted removal outcome = %+v", outcome)
	}
	if delta := outcome.Result.ClauseDelta; delta == nil || delta.Dispositions[removedBackendClause] != ClauseUnaccounted {
		t.Fatalf("unaccounted delta = %+v", delta)
	}
	if !strings.Contains(outcome.Result.Message, removedBackendClause) {
		t.Fatalf("refusal does not name removed clause: %q", outcome.Result.Message)
	}
}

func newClauseReplacementAdopter(t *testing.T) (PlanRequest, *Catalog) {
	t.Helper()
	catalog := cloneCatalogForRetentionDrift(t, mustEmbeddedCatalog(t))
	request := PlanRequest{
		Repository:   newProjectDecisionPlanRepository(t),
		ProfileID:    "standard-typescript-monorepo",
		Decisions:    standardTypeScriptDecisions("make verify"),
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	}
	adoption, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if adoption.Plan == nil {
		t.Fatalf("build adopter: %+v", adoption.Result)
	}
	if _, err := applyPlanWithCatalog(context.Background(), request.Repository, *adoption.Plan, adoption.Plan.PlanDigest, catalog); err != nil {
		t.Fatal(err)
	}
	manifest := adoption.Plan.SetupManifest
	if manifest.Generator.Baseline != "baseline.standard-typescript-monorepo-0.0.1" {
		t.Fatalf("adopter source baseline = %q", manifest.Generator.Baseline)
	}
	for index := range manifest.ManagedArtifacts {
		manifest.ManagedArtifacts[index].Digest = strings.Repeat("0", 64)
	}
	data, err := marshalSetupManifestBytes(manifest)
	if err != nil {
		t.Fatal(err)
	}
	writeTransactionFile(t, request.Repository, manifestPath, string(data), 0o644)
	commitInspectionRepository(t, request.Repository, "age adopter managed artifact digests")
	request.Preservation = RootPreservationRequest{Mode: PreservationModeManagedRefresh}
	return request, catalog
}

func TestClauseReplacementDeclarationsAreValidated(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		value      any
		otherClaim bool
		code       string
	}{
		{"not a list", removedBackendClause, false, "catalog.clause.replaces.invalid"},
		{"non-string", []any{42}, false, "catalog.clause.replaces.invalid"},
		{"empty ID", []any{""}, false, "catalog.clause.replaces.invalid"},
		{"repeated ID", []any{removedBackendClause, removedBackendClause}, false, "catalog.clause.replaces.invalid"},
		{"present ID in another module", []any{"clause.core.keep-root-compact"}, false, "catalog.clause.replaces.present"},
		{"two claimants", []any{removedBackendClause}, true, "catalog.clause.replaces.duplicate"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assets := cloneEmbeddedAssets(t)
			const assetPath = "modules/backend.json"
			var module document
			if err := json.Unmarshal(assets[assetPath].Data, &module); err != nil {
				t.Fatal(err)
			}
			for _, rule := range objectsOrEmpty(module["rules"]) {
				for _, clause := range objectsOrEmpty(rule["clauses"]) {
					id, _ := stringValue(clause, "id")
					if id == successorBackendClause {
						clause["replaces"] = tc.value
					} else if tc.otherClaim {
						clause["replaces"] = []any{removedBackendClause}
						tc.otherClaim = false
					}
				}
			}
			data, err := json.Marshal(module)
			if err != nil {
				t.Fatal(err)
			}
			assets[assetPath].Data = data
			_, err = LoadCatalog(assets)
			requireCatalogDiagnostic(t, err, tc.code)
		})
	}
}
