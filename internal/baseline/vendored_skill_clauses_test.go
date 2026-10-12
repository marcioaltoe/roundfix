// Suite: Vendored skill clauses
// Invariant: divergence guidance preserves clause identity, force, rendering and retention.
// Boundary IN: embedded catalog, formatter goldens and a temporary Source Baseline adopter.
// Boundary OUT: live upstream skills, applying a repository refresh and Daemon Verification.

package baseline

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// These literals lock the authored contract, independently of the module assets.
var vendoredSkillClauses = []struct{ id, guidance, module, rule, enforcement string }{
	{"clause.core.activate-matching-skills", "Activate every matching required skill before governed work. When one skill has distinct active-module triggers, retain and follow each trigger. When a skill's default conflicts with a Baseline rule or a Repository-Specific Normative Rule, follow the rule. Read a skill's bundled files from the directory where the skill is installed, never from a fixed absolute path the skill names.", "core", "rule.core.skill-dispatch", "mandatory"},
	{"clause.core.prohibit-editing-vendored-skills", "Do not edit a skill the repository installs from an upstream source; its upstream owns its text, and the next sync replaces the folder. When such a skill diverges from what the repository needs, choose in this order. First, state the correction as a rule that governs over the skill's default, in the Baseline or as a Repository-Specific Normative Rule; this is the answer to a convention difference such as the package manager, the test runner, or a library choice. Second, add a repository-authored skill that supplies or overrides only what is missing and leaves the upstream skill installed. Third, replace the upstream skill with a repository-authored one only when it is useless or almost useless to the repository. Send an objective upstream defect to the upstream as a pull request; a defect never justifies the third choice.", "core", "rule.core.skill-dispatch", "prohibited"},
	{"clause.rust.read-current-docs", "Use current authoritative crate and toolchain documentation before changing Rust APIs, configuration, or dependencies. Follow the edition the crate's `Cargo.toml` declares, not the one a skill example assumes: under edition 2024, `std::env::set_var` and `std::env::remove_var` are `unsafe`, and each call needs an `unsafe` block whose `// SAFETY:` comment states why no other thread reads or writes the environment at that point.", "rust", "rule.rust.current-docs", "mandatory"},
}

func TestTheVendoredSkillClausesCarryTheDivergenceOrder(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	for _, want := range vendoredSkillClauses {
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
					if clause["enforcement"] != want.enforcement || clause["guidance"] != want.guidance {
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

func TestTheVendoredSkillClausesRenderInTheGuides(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	for _, want := range vendoredSkillClauses {
		if want.module != "core" {
			continue
		}
		t.Run(want.id, func(t *testing.T) {
			asset, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md")
			if !ok {
				t.Fatal("missing skill-dispatch guide")
			}
			line := "\n- **" + want.enforcement + "**: " + want.guidance + "\n"
			if strings.Count("\n"+string(asset.Data)+"\n", line) != 1 {
				t.Fatalf("guide lacks exactly one forced clause %s", want.id)
			}
		})
	}
}

func TestAnAdopterRetainsTheVendoredSkillClauses(t *testing.T) {
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
	for _, want := range vendoredSkillClauses {
		if want.module != "core" {
			continue
		}
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
