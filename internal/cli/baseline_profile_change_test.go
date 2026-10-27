// Suite: Baseline profile changes
// Invariant: profile changes carry recorded answers and refuse all missing inputs together; updates never change profiles.
// Boundary IN: plan and update CLI results and temporary repository bytes.
// Boundary OUT: external skill acquisition.
package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestBaselinePlanNamesEveryMissingAnswerInOneRefusal(t *testing.T) {
	t.Parallel()
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, "seed.txt", "seed\n")
	commitBaselinePlanTestRepository(t, repo)
	var out, diagnostics bytes.Buffer
	code := RunContext(t.Context(), []string{"baseline", "plan", "--repo", repo, "--profile", "python", "--format=json"}, &out, &diagnostics)
	var result baseline.Result
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if code != exitUnverified || result.Category != "decision" {
		t.Fatalf("exit=%d result=%+v stderr=%s", code, result, &diagnostics)
	}
	for _, want := range []string{"required Baseline decisions are missing: autonomous.enabled, domain.layout, language.generated, repository.extension.enabled, secondbrain.enabled, spec.scaffold, triage.external, verification.gate, verification.incremental", "required when autonomous.enabled is true: runtime.backend, runtime.design", "instruction-preservation mode is required"} {
		if !strings.Contains(result.Message, want) {
			t.Errorf("message missing %q: %s", want, result.Message)
		}
	}
	if result.NextAction != "answer every named decision with --decision or --decision-file, including --decision preservation.mode=greenfield or preservation, and rerun" {
		t.Errorf("next action=%q", result.NextAction)
	}
}

func TestBaselineUpdateRefusesAProfileDecisionAndNamesThePlanCommand(t *testing.T) {
	t.Parallel()
	repo := newBaselineUpdateRepository(t)
	before := baselinePlanTestTree(t, repo)
	result, out, _, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--decision", "profile=python", "--no-skills", "--no-history", "--format=json")
	if code != exitPreflight || result.Category != "invalid" {
		t.Fatalf("exit=%d result=%s", code, out)
	}
	if !strings.Contains(result.Message, `"profile" is not a Baseline decision; baseline update keeps the recorded Baseline Profile go-cli-tui`) {
		t.Errorf("message=%q", result.Message)
	}
	if !strings.Contains(result.NextAction, "roundfix baseline plan --repo "+repo+" --profile python") || !strings.Contains(result.NextAction, "roundfix baseline apply") {
		t.Errorf("next action=%q", result.NextAction)
	}
	var textOut, textErr bytes.Buffer
	code = RunContext(t.Context(), []string{"baseline", "update", "--repo", repo, "--decision", "unknown=invalid", "--decision", "profile=python", "--no-skills", "--no-history"}, &textOut, &textErr)
	if code != exitPreflight || !strings.Contains(textOut.String(), result.Message) || !strings.Contains(textOut.String(), result.NextAction) {
		t.Fatalf("text refusal exit=%d stdout=%s stderr=%s", code, &textOut, &textErr)
	}
	if before != baselinePlanTestTree(t, repo) {
		t.Fatal("refusal wrote repository bytes")
	}
}

func TestBaselinePlanCarriesTheRecordedDecisionsAcrossAProfileChange(t *testing.T) {
	t.Parallel()
	repo := newBaselinePlanTestRepository(t)
	for _, skill := range []string{"context7", "context7-cli", "exa-web-search"} {
		writeBaselinePlanTestFile(t, repo, ".agents/skills/"+skill+"/SKILL.md", "# "+skill+"\n")
	}
	writeBaselinePlanTestFile(t, repo, "Makefile", "verify:\n\t@true\nverify-incremental:\n\t@true\n")
	commitBaselinePlanTestRepository(t, repo)
	adopted, planPath := baselineApplyTestPlan(t, repo)
	var out, diagnostics bytes.Buffer
	code := RunContext(t.Context(), []string{"baseline", "apply", "--repo", repo, "--plan", planPath, "--confirm-plan", adopted.PlanDigest, "--format=json"}, &out, &diagnostics)
	if code != exitOK {
		t.Fatalf("adoption exit=%d stdout=%s stderr=%s", code, &out, &diagnostics)
	}
	commitBaselinePlanTestRepository(t, repo)
	recorded := ReadBaselineSetupManifest(t, repo)
	args := []string{"baseline", "plan", "--repo", repo, "--profile", "python", "--decision", "preservation.mode=preservation"}
	before := baselinePlanTestTree(t, repo)
	out.Reset()
	diagnostics.Reset()
	code = RunContext(t.Context(), append(args, "--format=json"), &out, &diagnostics)
	if code != exitOK {
		t.Fatalf("profile change exit=%d stdout=%s stderr=%s", code, &out, &diagnostics)
	}
	plan, err := baseline.ParsePlanDocument(out.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	ids := []string{}
	for _, decision := range plan.Decisions {
		if !reflect.DeepEqual(decision.Value, recorded.Decisions[decision.ID].Value) {
			t.Errorf("decision %s=%v recorded=%v", decision.ID, decision.Value, recorded.Decisions[decision.ID].Value)
		}
		ids = append(ids, decision.ID)
	}
	if len(ids) != len(recorded.Decisions) {
		t.Fatalf("carried=%v recorded=%v", ids, recorded.Decisions)
	}
	sort.Strings(ids)
	out.Reset()
	diagnostics.Reset()
	code = RunContext(t.Context(), args, &out, &diagnostics)
	if code != exitOK || !strings.Contains(out.String(), "Baseline Profile: python\nCarried decisions: "+strings.Join(ids, ", ")+"\n") {
		t.Fatalf("text exit=%d stdout=%s stderr=%s", code, &out, &diagnostics)
	}
	// Repository-owned drafts use the same carrying path.
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	draft, err := baseline.NewProfileAdaptationDraft("python", "python-adapted", []string{"secondbrain", "autonomous-work"}, nil, catalog)
	if err != nil {
		t.Fatal(err)
	}
	draftPath := t.TempDir() + "/profile.json"
	if err := os.WriteFile(draftPath, draft.Document, 0o600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	diagnostics.Reset()
	code = RunContext(t.Context(), []string{"baseline", "plan", "--repo", repo, "--profile-file", draftPath, "--decision", "preservation.mode=preservation", "--format=json"}, &out, &diagnostics)
	if code != exitOK {
		t.Fatalf("draft carry exit=%d stdout=%s stderr=%s", code, &out, &diagnostics)
	}
	draftPlan, err := baseline.ParsePlanDocument(out.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	for _, decision := range draftPlan.Decisions {
		if !reflect.DeepEqual(decision.Value, recorded.Decisions[decision.ID].Value) {
			t.Errorf("draft changed recorded %s", decision.ID)
		}
	}
	count := 0
	for _, warning := range draftPlan.Warnings {
		if warning.Code == "baseline.decision.not-carried" && warning.Message == "recorded decision secondbrain.enabled is not carried: profile python-adapted does not select it" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("draft missing unique not-carried warning: %+v", draftPlan.Warnings)
	}
	// An explicit flag overrides the manifest and is not named as carried.
	out.Reset()
	diagnostics.Reset()
	code = RunContext(t.Context(), append(args, "--decision", "secondbrain.enabled=true", "--format=json"), &out, &diagnostics)
	if code != exitOK {
		t.Fatalf("override exit=%d stdout=%s stderr=%s", code, &out, &diagnostics)
	}
	override, err := baseline.ParsePlanDocument(out.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if override.SetupManifest.Decisions["secondbrain.enabled"].Value != true {
		t.Fatal("recorded value overrode flag")
	}
	out.Reset()
	diagnostics.Reset()
	code = RunContext(t.Context(), []string{"baseline", "plan", "--repo", repo, "--profile", "python", "--format=json"}, &out, &diagnostics)
	if code != exitUnverified || !strings.Contains(out.String(), "instruction-preservation mode is required") {
		t.Fatalf("missing preservation exit=%d stdout=%s", code, &out)
	}
	if before != baselinePlanTestTree(t, repo) {
		t.Fatal("planning changed repository bytes")
	}
}
