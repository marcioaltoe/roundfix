// Suite: Backend bucket prohibition
// Invariant: the guide scopes generic buckets and accounts for the clause adopters hold.
// Boundary IN: embedded catalog, Plan postimages, and Source Baseline transition.
// Boundary OUT: applying a Plan and executing repository Verification.

package baseline

import (
	"reflect"
	"strings"
	"testing"
)

func TestTheBackendGuideScopesTheBucketProhibition(t *testing.T) {
	t.Parallel()
	plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), standardTypeScriptDecisions("make verify"))
	text := string(planPostimage(t, plan, "docs/agents/backend.md").Content)
	const want = "- **prohibited**: Do not organize backend code into generic `modules` or `services` buckets in place of the domain, application, and infrastructure layers. A domain service that lives in the domain layer is not such a bucket."
	if !strings.Contains(text, want) {
		t.Errorf("backend guide does not contain scoped prohibition %q", want)
	}
	if strings.Contains(text, "Do not introduce generic `modules` or `services` buckets as the normative backend architecture.") {
		t.Error("backend guide still contains the unscoped prohibition")
	}
}

func TestTheBucketClauseReplacesTheClauseAdoptersHold(t *testing.T) {
	t.Parallel()
	const oldID = "clause.backend.prohibit-generic-layers"
	const newID = "clause.backend.prohibit-generic-buckets"
	for _, test := range []struct {
		name              string
		removeReplacement bool
		want              ClauseDisposition
	}{
		{name: "declared replacement", want: ClauseReplaced},
		{name: "missing replacement is unaccounted", removeReplacement: true, want: ClauseUnaccounted},
	} {
		t.Run(test.name, func(t *testing.T) {
			catalog := cloneCatalogForRetentionDrift(t, mustEmbeddedCatalog(t))
			profile, err := ResolveProfile("", "standard-typescript-monorepo", catalog)
			if err != nil {
				t.Fatal(err)
			}
			activeModules, artifacts, err := resolveManagedArtifacts(catalog, profile, standardTypeScriptDecisions("make verify"), false)
			if err != nil {
				t.Fatal(err)
			}
			managed := make([]ManifestArtifact, len(artifacts))
			for index, artifact := range artifacts {
				managed[index] = ManifestArtifact{
					ID: artifact.ID, Path: artifact.Path, Kind: artifact.Kind,
					Module: artifact.Module, Template: artifact.Template,
					Version: artifact.Version, Digest: artifact.Digest,
				}
			}
			if test.removeReplacement {
				found := false
				for _, rule := range objectsOrEmpty(catalog.modules["backend"]["rules"]) {
					for _, clause := range objectsOrEmpty(rule["clauses"]) {
						if clause["id"] == newID {
							delete(clause, "replaces")
							found = true
						}
					}
				}
				if !found {
					t.Fatal("replacement clause is absent from the catalog")
				}
			}
			source, err := catalog.SourceBaseline("baseline.standard-typescript-monorepo-0.0.1")
			if err != nil {
				t.Fatal(err)
			}
			evidence, delta := classifySourceClauseTransition(source, managed, catalog, activeModules)
			if got := delta.Dispositions[oldID]; got != test.want {
				t.Errorf("old clause disposition = %q, want %q", got, test.want)
			}
			if got := delta.Dispositions[newID]; got != ClauseRetained {
				t.Errorf("new clause disposition = %q, want retained", got)
			}
			if !test.removeReplacement {
				found := false
				for _, entry := range evidence {
					if entry.FromClause == oldID {
						found = true
						if !reflect.DeepEqual(entry.Targets, []string{newID}) {
							t.Errorf("replacement targets = %v, want %q", entry.Targets, newID)
						}
					}
				}
				if !found {
					t.Error("old clause has no replacement evidence")
				}
			}
		})
	}
}
