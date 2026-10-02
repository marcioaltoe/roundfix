// Suite: Person-only skill dispatch
// Invariant: every profile asks the person to start skills that disable model invocation.
// Boundary IN: embedded catalog and rendered Plan postimages.
// Boundary OUT: installed skill metadata, applying a Plan, and repository Verification.

package baseline

import (
	"context"
	"slices"
	"strings"
	"testing"
)

func TestEveryBuiltInProfileTellsTheAgentToAskForAPersonOnlySkill(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	const clauseID = "clause.core.ask-the-person-to-start-a-person-only-skill"
	for id, profile := range catalog.profiles {
		t.Run(id, func(t *testing.T) {
			if !slices.Contains(stringsOrEmpty(profile["modules"]), "core") {
				t.Fatal("profile does not select core")
			}
			found := false
			for _, rule := range objectsOrEmpty(catalog.modules["core"]["rules"]) {
				if rule["id"] != "rule.core.skill-dispatch" {
					continue
				}
				for _, clause := range objectsOrEmpty(rule["clauses"]) {
					if clause["id"] == clauseID {
						found = true
						if clause["enforcement"] != "mandatory" {
							t.Errorf("person-only skill force = %v, want mandatory", clause["enforcement"])
						}
					}
				}
			}
			if !found {
				t.Fatal("core skill-dispatch rule lacks the person-only skill clause")
			}
		})
	}
}

func TestTheSkillGuidesStateThePersonOnlySkillClause(t *testing.T) {
	t.Parallel()
	const want = "- **mandatory**: When a matching skill can be started only by a person, because its metadata turns off model invocation, ask the person to run it instead of activating it or carrying out its workflow yourself."
	for _, profileID := range []string{"rust-cli", "go-cli-tui", "standard-typescript-monorepo"} {
		t.Run(profileID, func(t *testing.T) {
			decisions := planTestDecisions()
			if profileID == "standard-typescript-monorepo" {
				decisions = standardTypeScriptDecisions("make verify")
			}
			outcome, err := BuildPlan(context.Background(), PlanRequest{
				Repository: newProjectDecisionPlanRepository(t), ProfileID: profileID,
				Decisions: decisions, Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
			})
			if err != nil {
				t.Fatal(err)
			}
			if outcome.Plan == nil {
				t.Fatalf("no plan: %+v", outcome.Result)
			}
			text := strings.Join(strings.Fields(string(planPostimage(t, *outcome.Plan, "docs/agents/skill-dispatch.md").Content)), " ")
			if !strings.Contains(text, want) {
				t.Errorf("%s skill guide lacks mandatory person-only skill guidance", profileID)
			}
			if profileID == "rust-cli" && (!strings.Contains(text, "- `cut-release`:") || !strings.Contains(text, "`trigger.rust.cut-release`")) {
				t.Error("Rust skill guide lost its cut-release dispatch")
			}
		})
	}
}
