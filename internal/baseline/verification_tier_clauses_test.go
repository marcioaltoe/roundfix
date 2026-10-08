package baseline

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// These literals lock the authored contract, independently of the module assets.
var verificationTierClauses = []struct{ id, guidance, module, rule, guide string }{
	{"clause.core.run-selected-verification", "Outside a Run, run the selected repository Verification before a completion claim, treat every failure as blocking, and report the command plus its actionable diagnostic. In a Daemon-assigned turn the Daemon owns that claim: it runs the Task's declared Verification and the repository Verification at settlement.", "core", "rule.core.verification-selected", "agent-instructions.md"},
	{"clause.core.verification-two-tiers", "Outside a Run, use the selected incremental Verification named at the top of this guide for fast local checks; it answers whether the current change remains valid while reusing safe local state. In a Daemon-assigned turn, run focused tests of the changed packages and neither selected command, because the Daemon verifies the Task after handoff. CI must run the selected repository Verification from a fresh run; it answers whether the complete tree satisfies the repository contract. Baseline planning refuses until the repository selects and declares both commands, so neither tier is ever satisfied by omission. Execute authored Verification only on committed provenance: the Spec artifacts that carry its commands must be tracked in the repository and byte-identical to their committed bytes at the resolved revision, comparing the authored projection while excluding Daemon-owned `status` and `## Result` fields. A Run satisfies this by construction, because its Run Worktree is created from a commit. No other command checks it: `roundfix spec check --run-verification` executes the commands of the working-tree Spec in a checkout of `HEAD`, and `roundfix settle` reads the Task file of the directory it settles, so whoever runs them on an uncommitted or modified Spec first reads the commands they will execute. No approval record makes an untrusted source executable; commit the artifacts or do not execute them. Read-only checking remains available for any source.", "core", "rule.core.verification-selected", "agent-instructions.md"},
	{"clause.spec.verification-two-tiers", "Outside a Run, run the selected incremental Verification named in `docs/agents/agent-instructions.md` for each Task to answer whether the current slice remains valid before handoff; a Daemon-assigned Task hands back after focused tests. CI must run the selected repository Verification from a fresh run to answer whether the assembled tree satisfies the repository contract. A missing incremental selection is a Baseline decision to answer, never a license to skip the local tier or a waiver to repeat in each Spec.", "spec-workflow", "rule.spec.routing", "spec-routing.md"},
	{"clause.spec.verification-fails-before-the-change", "Author each Task's Verification so every command fails on the tree before the Task's change, and run `roundfix spec check <slug> --strict --run-verification` before a Run starts: strict mode also fails on gap findings, and the Daemon refuses a Task whose command already exits zero on the unchanged tree.", "spec-workflow", "rule.spec.routing", "spec-routing.md"},
}

func TestTheVerificationTierClausesCarryTheirText(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	for _, want := range verificationTierClauses {
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

func TestTheVerificationTierClausesRenderInTheGuides(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	for _, want := range verificationTierClauses {
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

func TestAnAdopterRetainsTheVerificationTierClauses(t *testing.T) {
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
	for _, want := range verificationTierClauses {
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
