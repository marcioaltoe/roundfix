// Suite: incremental Verification decision
// Invariant: every Baseline plan selects and locally declares an incremental command independently of the complete gate.
// Boundary IN: embedded catalog declarations, alignment projection, plan refusal, and managed guide rendering.
// Boundary OUT: CLI update migration, which is exercised in internal/cli.

package baseline

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

const (
	incrementalVerificationDecisionID = "verification.incremental"
	incrementalVerificationSuggestion = "rtk make verify-incremental"
)

func TestIncrementalVerificationDecisionIsDeclaredWithoutADefault(t *testing.T) {
	t.Parallel()

	catalog := loadIncrementalVerificationCatalog(t)
	entry, found := catalog.Decision(incrementalVerificationDecisionID)
	decision := incrementalVerificationEntry(t, entry, found)
	if got, _ := decision["type"].(string); got != "string" {
		t.Fatalf("incremental decision type = %q, want string", got)
	}
	if got, _ := decision["suggestion"].(string); got != incrementalVerificationSuggestion {
		t.Fatalf("incremental decision suggestion = %q, want %q", got, incrementalVerificationSuggestion)
	}
	if _, exists := decision["default"]; exists {
		t.Fatalf("incremental decision unexpectedly declares a default: %v", decision["default"])
	}
	effects := objectsOrEmpty(decision["effects"])
	if len(effects) != 1 {
		t.Fatalf("incremental decision effects = %v, want one", effects)
	}
	bindings := objectsOrEmpty(effects[0]["renderBindings"])
	if len(bindings) != 1 || bindings[0]["artifact"] != "guide.agent-instructions" ||
		bindings[0]["template"] != "template.guide.agent-instructions" ||
		bindings[0]["token"] != incrementalVerificationDecisionID {
		t.Fatalf("incremental decision render binding = %v", bindings)
	}
}

func TestEveryBuiltInProfileRequiresTheIncrementalVerificationDecision(t *testing.T) {
	t.Parallel()

	catalog := loadIncrementalVerificationCatalog(t)
	for _, profileID := range []string{"go-cli-tui", "rust-cli", "standard-typescript-monorepo"} {
		entry, found := catalog.Profile(profileID)
		profile := incrementalVerificationEntry(t, entry, found)
		if !slices.Contains(stringsOrEmpty(profile["entryDecisions"]), incrementalVerificationDecisionID) {
			t.Errorf("profile %q entry decisions = %v, want %q", profileID, profile["entryDecisions"], incrementalVerificationDecisionID)
		}
	}
	entry, found := catalog.Module("core")
	core := incrementalVerificationEntry(t, entry, found)
	if !slices.Contains(stringsOrEmpty(core["requiredDecisions"]), incrementalVerificationDecisionID) {
		t.Fatalf("core required decisions = %v, want %q", core["requiredDecisions"], incrementalVerificationDecisionID)
	}
}

func TestIncrementalVerificationDecisionProjectsLikeTheGate(t *testing.T) {
	t.Parallel()

	repository := newAlignedTypeScriptRepository(t)
	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, loadIncrementalVerificationCatalog(t))
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	var incremental VerificationProjection
	for _, projection := range alignment.Verification {
		if projection.ID == incrementalVerificationDecisionID {
			count++
			incremental = projection
		}
	}
	if count != 1 {
		t.Fatalf("incremental projections = %d, want exactly one: %+v", count, alignment.Verification)
	}
	if incremental.Role != "incremental" || incremental.Command != "make verify-incremental" ||
		incremental.Classification != VerificationRepositoryCommand || !incremental.RepositoryExecutable ||
		incremental.DeclarationPath != "Makefile" || !strings.HasPrefix(incremental.DeclarationDigest, "sha256:") {
		t.Fatalf("incremental projection = %+v", incremental)
	}
}

func TestUndeclaredIncrementalVerificationBlocksPlanning(t *testing.T) {
	t.Parallel()

	repository := newPlanRepository(t)
	writeInspectionFile(t, repository, "Makefile", "verify:\n\t@true\n")
	commitInspectionRepository(t, repository, "remove incremental declaration")

	catalog := loadIncrementalVerificationCatalog(t)
	profile, err := ResolveProfile("", "go-cli-tui", catalog)
	if err != nil {
		t.Fatal(err)
	}
	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: profile.ID,
		Profile:   &profile,
		Decisions: profileAlignmentDecisions(profile, planTestDecisions()),
	}, catalog)
	if err != nil {
		t.Fatal(err)
	}
	divergence, found := findProfileDivergence(alignment.Divergences, incrementalVerificationDecisionID)
	if !found || divergence.Code != "verification.command.undeclared" || !divergence.Blocking ||
		divergence.Message != `selected incremental Verification command "make verify-incremental" has no matching local declaration` {
		t.Fatalf("incremental divergence = %+v, found=%v", divergence, found)
	}

	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository:   repository,
		ProfileID:    "go-cli-tui",
		Decisions:    planTestDecisions(),
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan != nil || outcome.Result.State != "action_required" ||
		!strings.Contains(outcome.Result.Message, incrementalVerificationDecisionID) {
		t.Fatalf("undeclared incremental planning outcome = %+v", outcome)
	}
}

func TestPlanWithoutIncrementalVerificationNamesTheDecision(t *testing.T) {
	t.Parallel()

	decisions := slices.DeleteFunc(planTestDecisions(), func(decision DecisionValue) bool {
		return decision.ID == incrementalVerificationDecisionID
	})
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository:   newPlanRepository(t),
		ProfileID:    "go-cli-tui",
		Decisions:    decisions,
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan != nil || outcome.Result.State != "action_required" ||
		outcome.Result.Category != "decision" || !strings.Contains(outcome.Result.Message, incrementalVerificationDecisionID) {
		t.Fatalf("missing incremental decision outcome = %+v", outcome)
	}
}

func TestAgentInstructionsRenderTheSelectedIncrementalVerification(t *testing.T) {
	t.Parallel()

	repository := newPlanRepository(t)
	writeInspectionFile(t, repository, "Makefile", "verify:\n\t@true\nverify-fast:\n\t@true\n")
	commitInspectionRepository(t, repository, "declare selected incremental verification")
	decisions := planTestDecisions()
	for index := range decisions {
		if decisions[index].ID == incrementalVerificationDecisionID {
			decisions[index].Value = "make verify-fast"
		}
	}
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository:   repository,
		ProfileID:    "go-cli-tui",
		Decisions:    decisions,
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil {
		t.Fatalf("rendering plan returned action result: %+v", outcome.Result)
	}
	postimage := planPostimage(t, *outcome.Plan, "docs/agents/agent-instructions.md")
	if !bytesContainsLine(postimage.Content, "The selected incremental Verification is `make verify-fast`.") {
		t.Fatalf("agent instructions do not render selected incremental command:\n%s", postimage.Content)
	}
}

func TestTwoTierClausesNameTheSelectedCommands(t *testing.T) {
	t.Parallel()

	catalog := loadIncrementalVerificationCatalog(t)
	for _, test := range []struct {
		moduleID string
		clauseID string
	}{
		{moduleID: "core", clauseID: "clause.core.verification-two-tiers"},
		{moduleID: "spec-workflow", clauseID: "clause.spec.verification-two-tiers"},
	} {
		entry, found := catalog.Module(test.moduleID)
		module := incrementalVerificationEntry(t, entry, found)
		guidance := incrementalVerificationClauseGuidance(t, module, test.clauseID)
		if !strings.Contains(guidance, "selected incremental Verification") ||
			strings.Contains(guidance, "declared by the active Baseline Profile") {
			t.Errorf("%s guidance = %q", test.clauseID, guidance)
		}
	}
}

func loadIncrementalVerificationCatalog(t *testing.T) *Catalog {
	t.Helper()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	return catalog
}

func incrementalVerificationEntry(t *testing.T, entry Entry, found bool) document {
	t.Helper()
	if !found {
		t.Fatal("catalog entry not found")
	}
	var decoded document
	if err := json.Unmarshal(entry.Data, &decoded); err != nil {
		t.Fatal(err)
	}
	return decoded
}

func incrementalVerificationClauseGuidance(t *testing.T, module document, clauseID string) string {
	t.Helper()
	for _, rule := range objectsOrEmpty(module["rules"]) {
		for _, clause := range objectsOrEmpty(rule["clauses"]) {
			if clause["id"] == clauseID {
				guidance, _ := clause["guidance"].(string)
				return guidance
			}
		}
	}
	t.Fatalf("catalog module has no clause %q", clauseID)
	return ""
}

func bytesContainsLine(content []byte, line string) bool {
	for _, candidate := range strings.Split(string(content), "\n") {
		if candidate == line {
			return true
		}
	}
	return false
}
