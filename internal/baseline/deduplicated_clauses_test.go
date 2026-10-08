package baseline

import (
	"context"
	"strings"
	"testing"
)

// Literals pin the authored contract independently of the catalog assets.
var mergedClauseTexts = []struct{ id, guidance, enforcement string }{
	{"clause.core.request-review-explicitly", "Resolve and follow the repository's pre-PR review policy before publication: `codex`, `claude`, `coderabbit`, or explicit `none`, with Codex as the built-in default, User Config over built-in defaults, and Project Config over User Config. Never re-enable an opted-out provider; CodeRabbit is optional, never a mandatory dependency. With a provider enabled, obtain its review of the current candidate, with its configured model and effort, before opening the Pull Request, and block on failed, unavailable, incomplete, absent, or stale review evidence, which never selects `none` automatically. With explicit `none`, request no review, record that review was disabled by configuration, and keep QA and required checks; disabled review is not a passing review. Neither choice waives repository or GitHub requirements, and unimplemented configuration or runtime support is never presented as a working command.", "mandatory"},
	{"clause.core.prohibit-external-research-for-local-code", "Do not use external research tools to discover or infer local repository code or behavior, and never put credentials, private client records, or proprietary source code in an external search query.", "prohibited"},
	{"clause.autonomous.loop-05-clean-is-not-evidence", "Treat a terminal Clean and a resolved status as claims, not evidence: read the diff after a Run reports Clean and, when review is enabled, confirm that the selected provider reviewed the current candidate, because silence or a provider-reported skip is not approval. Preserve historical Review Skipped and Clean Unverified meanings.", "mandatory"},
	{"clause.spec.project-constraints-05-legacy-and-ownership", "Keep completed or archived legacy Specs byte-identical.", "mandatory"},
	{"clause.secondbrain.03-decision-consultation", "Use Secondbrain for prior decisions, related project experience, and existing research.", "mandatory"},
}

var removedClauseSuccessors = []struct{ removed, successor, enforcement, guidance string }{
	{"clause.core.request-pull-request-review", "clause.core.request-review-explicitly", "mandatory", "Resolve the repository's pre-PR review policy before publication. Supported policy choices are `codex`, `claude`, `coderabbit`, and explicit `none`; Codex remains the built-in default when no applicable explicit selection is declared. Preserve User Config over built-in defaults and Project Config over User Config; an explicit project selection takes precedence. Honor the selected provider and its configured model/effort where applicable. With review enabled, obtain the configured provider's review of the current candidate before opening the Pull Request and require a complete applicable result. With explicit `none`, omit the reviewer and provider calls and record that review was disabled by configuration; otherwise authorized publication and merge may proceed after QA and required checks. Disabled review is not a passing review. This policy does not waive repository or GitHub requirements."},
	{"clause.spec.status-only-in-task", "clause.spec.tracker-artifacts", "mandatory", "Keep Task status only in the assigned Task file frontmatter. The Task Graph records topology and dependencies, not progress."},
	{"clause.secondbrain.prohibit-external-local-discovery", "clause.core.prohibit-external-research-for-local-code", "prohibited", "Do not use Exa or another external research tool to discover or infer local repository code or behavior. Use local code-search tools for that purpose. Never include credentials, private client records, or proprietary source code in external search queries."},
}

func TestTheMergedClausesCarryTheirText(t *testing.T) {
	t.Parallel()
	clauses := embeddedBaselineClauses(t)
	for _, want := range mergedClauseTexts {
		t.Run(want.id, func(t *testing.T) {
			found := false
			for _, clause := range clauses {
				if clause.ID == want.id {
					found = true
					if clause.Guidance != want.guidance || clause.Enforcement != want.enforcement {
						t.Errorf("merged clause text/force = %+v, want %+v", clause, want)
					}
				}
			}
			if !found {
				t.Error("missing merged clause")
			}
		})
	}
}

func TestEveryRemovedClauseHasOneSuccessor(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	for _, want := range removedClauseSuccessors {
		t.Run(want.removed, func(t *testing.T) {
			claims := 0
			for _, clause := range embeddedBaselineClauses(t) {
				if strings.Contains(strings.Join(strings.Fields(clause.Guidance), " "), strings.Join(strings.Fields(want.guidance), " ")) {
					t.Errorf("removed text restored in %s", clause.ID)
				}
			}
			for _, module := range catalog.modules {
				for _, rule := range objectsOrEmpty(module["rules"]) {
					for _, clause := range objectsOrEmpty(rule["clauses"]) {
						if clause["id"] == want.removed {
							t.Error("removed clause still present")
						}
						for _, replaced := range stringsOrEmpty(clause["replaces"]) {
							if replaced != want.removed {
								continue
							}
							claims++
							if clause["id"] != want.successor || clause["enforcement"] != want.enforcement {
								t.Errorf("wrong successor or force: %+v", clause)
							}
						}
					}
				}
			}
			if claims != 1 {
				t.Errorf("successor count = %d, want 1", claims)
			}
		})
	}
}

func TestAnAdopterRecordsTheRemovedClausesReplaced(t *testing.T) {
	t.Parallel()
	request, catalog := newClauseReplacementAdopter(t)
	// Exercise the optional Secondbrain clauses as well as the core and Spec clauses.
	for index := range request.Decisions {
		if request.Decisions[index].ID == "secondbrain.enabled" {
			request.Decisions[index].Value = true
		}
	}
	// Record an adopter with the optional guide before testing its Source Baseline refresh.
	enabled, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if enabled.Plan == nil || enabled.Result.State != "ready" {
		t.Fatalf("enable Secondbrain: %+v", enabled.Result)
	}
	if _, err := applyPlanWithCatalog(context.Background(), request.Repository, *enabled.Plan, enabled.Plan.PlanDigest, catalog); err != nil {
		t.Fatal(err)
	}
	manifest := enabled.Plan.SetupManifest
	for index := range manifest.ManagedArtifacts {
		manifest.ManagedArtifacts[index].Digest = strings.Repeat("0", 64)
	}
	data, err := marshalSetupManifestBytes(manifest)
	if err != nil {
		t.Fatal(err)
	}
	writeTransactionFile(t, request.Repository, manifestPath, string(data), 0o644)
	commitInspectionRepository(t, request.Repository, "age adopter with Secondbrain enabled")
	outcome, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil || outcome.Result.State != "ready" {
		t.Fatalf("Managed Refresh is not ready: %+v", outcome.Result)
	}
	delta := outcome.Plan.ClauseDelta
	if delta == nil {
		t.Fatal("refresh lacks clause accounting")
	}
	for id, disposition := range delta.Dispositions {
		if disposition == ClauseUnaccounted {
			t.Errorf("unaccounted clause %s", id)
		}
	}
	for _, want := range removedClauseSuccessors {
		t.Run(want.removed, func(t *testing.T) {
			if delta.Dispositions[want.removed] != ClauseReplaced {
				t.Errorf("disposition = %s, want replaced", delta.Dispositions[want.removed])
			}
			found := false
			for _, evidence := range outcome.Plan.Retention {
				if evidence.FromClause != want.removed {
					continue
				}
				found = true
				if evidence.Disposition != string(ClauseReplaced) || len(evidence.Targets) != 1 || evidence.Targets[0] != want.successor || !strings.Contains(evidence.Reason, want.successor) {
					t.Errorf("replacement evidence = %+v", evidence)
				}
			}
			if !found {
				t.Error("missing replacement evidence")
			}
		})
	}
}
