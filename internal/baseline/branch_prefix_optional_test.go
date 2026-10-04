package baseline

import (
	"slices"
	"strings"
	"testing"
)

func TestAnUnrecordedBranchPrefixStatesTheCommitTypeRule(t *testing.T) {
	decisions := standardTypeScriptDecisions("make verify")
	decisions = slices.DeleteFunc(decisions, func(d DecisionValue) bool { return d.ID == "branch.prefix" })
	plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), decisions)
	guide := string(planPostimage(t, plan, "docs/agents/agent-instructions.md").Content)
	want := "No branch prefix is recorded. Name new work branches `<type>/<description>`,\nwhere `<type>` is the work's Conventional Commit type, as the branch rule\nbelow states. Tool-owned Run and Task branches follow their tool's documented\nnamespace."
	if !strings.Contains(guide, want+"\n\n") || strings.Contains(guide, "The branch-prefix pattern is") {
		t.Fatalf("unrecorded prefix guide = %s", guide)
	}
	for _, d := range plan.Decisions {
		if d.ID == "branch.prefix" {
			t.Fatal("unrecorded prefix was recorded")
		}
	}
}

func TestARecordedBranchPrefixKeepsItsSentence(t *testing.T) {
	for _, value := range []string{"<type>/", "ma/"} {
		t.Run(value, func(t *testing.T) {
			decisions := standardTypeScriptDecisions("make verify")
			for i := range decisions {
				if decisions[i].ID == "branch.prefix" {
					decisions[i].Value = value
				}
			}
			plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), decisions)
			guide := string(planPostimage(t, plan, "docs/agents/agent-instructions.md").Content)
			want := "The branch-prefix pattern is `" + value + "`; `<type>` is replaced by the\nwork's purpose, never used literally. Use `<type>/` as the portable decision\nvalue. Legacy personal-prefix values must be revised through Baseline and do\nnot override the purpose-based branch rule below. Tool-owned Run and Task\nbranches follow their tool's documented namespace."
			if !strings.Contains(guide, want+"\n\n") || strings.Contains(guide, "No branch prefix is recorded.") {
				t.Fatalf("recorded prefix guide = %s", guide)
			}
		})
	}
}

func TestBranchPrefixIsAnOptionalDecision(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	if !decisionOptional(catalog.decisions["branch.prefix"]) {
		t.Fatal("branch.prefix is required")
	}
	if slices.Contains(stringsOrEmpty(catalog.modules["core"]["requiredDecisions"]), "branch.prefix") {
		t.Fatal("core requires branch.prefix")
	}
	for _, id := range catalog.ProfileIDs() {
		t.Run(id, func(t *testing.T) {
			profile, err := ResolveProfile("", id, catalog)
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Contains(profile.Decisions, "branch.prefix") {
				t.Fatal("profile no longer selects branch.prefix")
			}
		})
	}
}
