package baseline

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// These literals lock the authored contract, independently of the module assets.
var loopClauses = []struct{ id, guidance, module, rule, guide string }{
	{"clause.autonomous.loop-01-qa-once", "Before a Spec authored beside another one enters a Run or the Delivery Queue, rebase it on the default branch and run `roundfix spec check <slug> --strict --run-verification` again, because a decision record or glossary term the other Spec merged can make it fail.", "autonomous-work", "rule.autonomous.loop", "autonomous-work.md"},
	{"clause.autonomous.loop-01-qa-once", "When a corrective Task is added after the gate settled, add it to the QA Task's `needs` in the Task Graph and then reopen the gate with `roundfix reopen --spec <slug>`, because reopen returns the gate to `pending` only over a Task inside the gate's dependency closure; never edit the QA Task by hand.", "autonomous-work", "rule.autonomous.loop", "autonomous-work.md"},
	{"clause.autonomous.loop-04-verify-the-class", "A Task that changes a message, field, flag, refusal, or rule searches the repository's documentation, agent guides, and skill copies for the old wording and declares every file that must change with it; the QA gate repeats the same search. Characterize current behavior before changing it, and declare each break: name every existing test, golden, or fixture the Task changes and the contract it moves to.", "autonomous-work", "rule.autonomous.loop", "autonomous-work.md"},
}

const hermeticGoGuidance = "Keep tests hermetic: set or clear with `t.Setenv` every environment variable the code under test reads, and never read the host's credentials, home directory, or tool state. Create a Unix socket under a short directory from `os.MkdirTemp(\"\", ...)` rather than a deep `t.TempDir()`, because macOS refuses a socket path longer than 104 bytes, and never let a Verification depend on a test that can skip on the host."

func TestTheLoopClausesCarryTheirText(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	for _, want := range loopClauses {
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
					if want.id == "clause.autonomous.loop-01-qa-once" && strings.Contains(clause["guidance"].(string), "reopen the gate with `roundfix reopen --spec <slug>`; never edit the QA Task by hand") {
						t.Error("loop-01 still carries the old corrective Task instruction")
					}
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

func TestTheLoopClausesRenderInTheGuides(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	for _, want := range loopClauses {
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

func TestAnAdopterRetainsTheLoopClauses(t *testing.T) {
	request, catalog := newClauseReplacementAdopter(t)
	// The shared adopter disables autonomous work. Enable it and apply its
	// guide before aging the manifest, so the refresh accounts for loop clauses.
	for i := range request.Decisions {
		if request.Decisions[i].ID == "autonomous.enabled" {
			request.Decisions[i].Value = true
		}
	}
	request.Decisions = append(request.Decisions,
		DecisionValue{ID: "runtime.backend", Value: "codex gpt-5.5 xhigh"},
		DecisionValue{ID: "runtime.design", Value: "claude opus xhigh"},
	)
	enabled, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if enabled.Plan == nil || enabled.Result.State != "ready" {
		t.Fatalf("enable autonomous adopter: %+v", enabled.Result)
	}
	if _, err := applyPlanWithCatalog(context.Background(), request.Repository, *enabled.Plan, enabled.Plan.PlanDigest, catalog); err != nil {
		t.Fatal(err)
	}
	manifest := enabled.Plan.SetupManifest
	for i := range manifest.ManagedArtifacts {
		manifest.ManagedArtifacts[i].Digest = strings.Repeat("0", 64)
	}
	data, err := marshalSetupManifestBytes(manifest)
	if err != nil {
		t.Fatal(err)
	}
	writeTransactionFile(t, request.Repository, manifestPath, string(data), 0o644)
	commitInspectionRepository(t, request.Repository, "age autonomous adopter managed artifact digests")
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
	for _, want := range loopClauses {
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

func TestTheHermeticGoClauseRendersInTheGoGuide(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	asset, ok := catalog.Module("go")
	if !ok {
		t.Fatal("missing Go module")
	}
	var module document
	if err := json.Unmarshal(asset.Data, &module); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, rule := range objectsOrEmpty(module["rules"]) {
		if rule["id"] != "rule.go.observable-tests" {
			continue
		}
		clauses := objectsOrEmpty(rule["clauses"])
		for i, clause := range clauses {
			if clause["id"] != "clause.go.keep-tests-hermetic" {
				continue
			}
			found = true
			if clause["enforcement"] != "mandatory" || clause["guidance"] != hermeticGoGuidance {
				t.Errorf("hermetic Go clause force/text differs: %+v", clause)
			}
			if i == 0 || clauses[i-1]["id"] != "clause.go.test-observable-behavior" {
				t.Error("hermetic Go clause must follow observable behavior clause")
			}
		}
	}
	if !found {
		t.Fatal("missing hermetic Go clause")
	}
	plan := buildTestPlan(t, newPlanRepository(t))
	local, err := os.ReadFile("../../docs/agents/go.md")
	if err != nil {
		t.Fatal(err)
	}
	line := "\n- **mandatory**: " + hermeticGoGuidance + "\n"
	for name, data := range map[string][]byte{"Go adopter postimage": planPostimage(t, plan, "docs/agents/go.md").Content, "repository guide": local} {
		t.Run(name, func(t *testing.T) {
			if count := strings.Count("\n"+string(data)+"\n", line); count != 1 {
				t.Fatalf("guide contains mandatory hermetic clause %d times, want exactly once", count)
			}
		})
	}
}
