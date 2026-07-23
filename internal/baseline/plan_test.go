// Suite: Portable Baseline Plans
// Invariant: equivalent bounded repository state and normalized decisions produce one strict, portable, digest-bound plan without writes.
// Boundary IN: plan assembly, codecs, rendering, ledger projection, digest, and clone preimage validation.
// Boundary OUT: transaction apply, rollback, and interactive decision collection.

package baseline

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestPlanDocumentStrictCodecs(t *testing.T) {
	plan := buildTestPlan(t, newPlanRepository(t))
	encoded, err := MarshalPlanDocument(plan)
	if err != nil {
		t.Fatalf("marshal Plan Document: %v", err)
	}
	decoded, err := ParsePlanDocument(encoded)
	if err != nil {
		t.Fatalf("parse Plan Document: %v", err)
	}
	if !reflect.DeepEqual(decoded, plan) {
		t.Fatalf("Plan Document round trip differs")
	}

	var raw map[string]any
	if err := json.Unmarshal(encoded, &raw); err != nil {
		t.Fatal(err)
	}
	raw["unknown"] = true
	unknown, _ := json.Marshal(raw)
	if _, err := ParsePlanDocument(unknown); err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("unknown field error = %v", err)
	}
	duplicate := bytes.Replace(encoded, []byte(`"schemaVersion":`),
		[]byte(`"schemaVersion":"roundfix/baseline-plan/v1","schemaVersion":`), 1)
	if _, err := ParsePlanDocument(duplicate); err == nil || !strings.Contains(err.Error(), "duplicate JSON key") {
		t.Fatalf("duplicate key error = %v", err)
	}

	result := Result{
		SchemaVersion: ResultSchemaVersion, Operation: "plan", State: "action_required",
		Category: "decision", Message: "missing", NextAction: "answer",
		VerifiedPostimages: []Postimage{}, Warnings: []Finding{}, Recommendations: []string{},
	}
	resultJSON, err := MarshalResult(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	if _, err := ParseResult(resultJSON); err != nil {
		t.Fatalf("parse result: %v", err)
	}
}

func TestFileChangesProjectionRejectsMismatch(t *testing.T) {
	plan := buildTestPlan(t, newPlanRepository(t))
	if len(plan.FileChanges) == 0 {
		t.Fatal("plan has no file changes")
	}
	plan.FileChanges[0].ManagedEntries = append(plan.FileChanges[0].ManagedEntries, "invented")
	plan.PlanDigest, _ = computePlanDigest(plan)
	if err := ValidatePlanDocument(plan); err == nil ||
		!strings.Contains(err.Error(), "fileChanges does not match") {
		t.Fatalf("projection mismatch error = %v", err)
	}
}

func TestPlanDigestBindsExactPostimagesAndIgnoresProjection(t *testing.T) {
	plan := buildTestPlan(t, newPlanRepository(t))
	original := plan.PlanDigest

	projected := plan
	projected.FileChanges = append([]FileChange(nil), plan.FileChanges...)
	projected.FileChanges[0].Action = "invented"
	digest, err := computePlanDigest(projected)
	if err != nil {
		t.Fatal(err)
	}
	if digest != original {
		t.Fatalf("projection changed digest: got %s want %s", digest, original)
	}

	changed := plan
	changed.Postimages = append([]Postimage(nil), plan.Postimages...)
	changed.Postimages[0].Content = append([]byte(nil), changed.Postimages[0].Content...)
	changed.Postimages[0].Content = append(changed.Postimages[0].Content, 'x')
	changed.Postimages[0].ContentIdentity = planContentIdentity(changed.Postimages[0].Content)
	digest, err = computePlanDigest(changed)
	if err != nil {
		t.Fatal(err)
	}
	if digest == original {
		t.Fatal("exact postimage bytes did not change Plan Digest")
	}
}

func TestCrossClonePlanAcceptsMatchingIdentityAndPreimages(t *testing.T) {
	source := newPlanRepository(t)
	plan := buildTestPlan(t, source)
	clone := filepath.Join(t.TempDir(), "clone")
	runPlanGit(t, "", "clone", "--quiet", source, clone)
	if err := ValidatePlanRepository(context.Background(), clone, plan); err != nil {
		t.Fatalf("matching clone rejected portable plan: %v", err)
	}
	writeInspectionFile(t, clone, ".agents/skills/context7/SKILL.md", "changed\n")
	if err := ValidatePlanRepository(context.Background(), clone, plan); err == nil ||
		!strings.Contains(err.Error(), "stale") {
		t.Fatalf("changed bounded preimage error = %v", err)
	}
}

func TestPlanDeterminismAndNoMutation(t *testing.T) {
	repo := newPlanRepository(t)
	before := testRepositoryDigest(t, repo)
	first := buildTestPlan(t, repo)
	second := buildTestPlan(t, repo)
	after := testRepositoryDigest(t, repo)
	if first.PlanDigest != second.PlanDigest {
		t.Fatalf("Plan Digests differ: %s != %s", first.PlanDigest, second.PlanDigest)
	}
	firstJSON, _ := MarshalPlanDocument(first)
	secondJSON, _ := MarshalPlanDocument(second)
	if !bytes.Equal(firstJSON, secondJSON) {
		t.Fatal("equivalent planning produced different JSON bytes")
	}
	if before != after {
		t.Fatalf("planning mutated repository: before=%s after=%s", before, after)
	}
	if strings.Contains(string(firstJSON), filepath.ToSlash(repo)) {
		t.Fatalf("portable plan contains checkout path %q", repo)
	}
	for _, postimage := range first.Postimages {
		if postimage.Path == "" || filepath.IsAbs(postimage.Path) {
			t.Fatalf("non-portable postimage path %q", postimage.Path)
		}
	}
}

func TestPlanDeterminismMatchesMaintainedManagedEntryFixture(t *testing.T) {
	fixturePath := filepath.Join(
		"..", "..", ".agents", "skills", "setup-context-driven", "assets",
		"parity-corpus", "v1", "fixtures", "greenfield-go-cli-tui.json",
	)
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		ManagedEntryLedger []struct {
			ID       string `json:"id"`
			Path     string `json:"path"`
			Kind     string `json:"kind"`
			Module   string `json:"module"`
			Template string `json:"template"`
			Version  string `json:"version"`
			Digest   string `json:"digest"`
		} `json:"managedEntryLedger"`
		PlannedByteSequence []struct {
			Path          string `json:"path"`
			AfterIdentity string `json:"afterIdentity"`
		} `json:"plannedByteSequence"`
	}
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	profile, err := ResolveProfile("", "go-cli-tui", catalog)
	if err != nil {
		t.Fatal(err)
	}
	_, artifacts, err := resolveManagedArtifacts(catalog, profile, planTestDecisions())
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != len(fixture.ManagedEntryLedger) {
		t.Fatalf("managed entry count = %d, want %d", len(artifacts), len(fixture.ManagedEntryLedger))
	}
	for index, artifact := range artifacts {
		if artifact.ID != fixture.ManagedEntryLedger[index].ID {
			t.Fatalf("managed entry %d = %q, want maintained order %q",
				index, artifact.ID, fixture.ManagedEntryLedger[index].ID)
		}
	}
	expected := make(map[string]any, len(fixture.ManagedEntryLedger))
	for _, entry := range fixture.ManagedEntryLedger {
		expected[entry.ID] = entry
	}
	for _, artifact := range artifacts {
		want, ok := expected[artifact.ID]
		if !ok {
			t.Fatalf("managed entry %q is not in maintained fixture", artifact.ID)
		}
		wantJSON, _ := json.Marshal(want)
		gotJSON, _ := json.Marshal(struct {
			ID       string `json:"id"`
			Path     string `json:"path"`
			Kind     string `json:"kind"`
			Module   string `json:"module"`
			Template string `json:"template"`
			Version  string `json:"version"`
			Digest   string `json:"digest"`
		}{
			ID: artifact.ID, Path: artifact.Path, Kind: artifact.Kind,
			Module: artifact.Module, Template: artifact.Template,
			Version: artifact.Version, Digest: artifact.Digest,
		})
		if !bytes.Equal(gotJSON, wantJSON) {
			t.Fatalf("managed entry %q differs:\ngot  %s\nwant %s", artifact.ID, gotJSON, wantJSON)
		}
	}

	plan := buildTestPlan(t, newPlanRepository(t))
	expectedPostimages := make(map[string]string)
	for _, entry := range fixture.PlannedByteSequence {
		if entry.Path != manifestPath {
			expectedPostimages[entry.Path] = entry.AfterIdentity
		}
	}
	for _, postimage := range plan.Postimages {
		want, ok := expectedPostimages[postimage.Path]
		if !ok {
			continue
		}
		if postimage.ContentIdentity != want {
			t.Fatalf("planned bytes %q identity = %s, want maintained %s",
				postimage.Path, postimage.ContentIdentity, want)
		}
		delete(expectedPostimages, postimage.Path)
	}
	if len(expectedPostimages) != 0 {
		t.Fatalf("maintained planned paths missing from plan: %v", sortedKeys(expectedPostimages))
	}
}

func TestPlanDocumentMissingDecisionsReturnsResultWithoutPartialPlan(t *testing.T) {
	repo := newPlanRepository(t)
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Preservation: RootPreservationRequest{
			Mode: PreservationModeGreenfield,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan != nil || outcome.Result.State != "action_required" ||
		outcome.Result.Category != "decision" || outcome.Result.NextAction == "" {
		t.Fatalf("missing-decision outcome = %+v", outcome)
	}
	if _, err := MarshalResult(outcome.Result); err != nil {
		t.Fatalf("missing-decision result schema: %v", err)
	}
}

func TestPlanDocumentIncludesMaintainedUpgradeRetention(t *testing.T) {
	repo := newPlanRepository(t)
	writeInspectionFile(t, repo, manifestPath, `{
  "schemaVersion": 1,
  "generator": {
    "skill": "setup-context-driven",
    "version": 1,
    "baseline": "baseline.managed-v2"
  },
  "managedArtifacts": []
}
`)
	plan := buildTestPlan(t, repo)
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	contract, err := catalog.UpgradeRetentionContract("transition.managed-v2-to-portable-v3")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Retention) != len(contract.PriorClauses) {
		t.Fatalf("retention entries = %d, want %d", len(plan.Retention), len(contract.PriorClauses))
	}
	for index, entry := range plan.Retention {
		if entry.FromClause != contract.Accounting[index].FromClause ||
			entry.Enforcement == "" ||
			entry.Disposition != contract.Accounting[index].Disposition ||
			!reflect.DeepEqual(entry.Targets, contract.Accounting[index].Targets) ||
			entry.Reason != contract.Accounting[index].Reason {
			t.Fatalf("retention entry %d = %+v, want contract accounting %+v",
				index, entry, contract.Accounting[index])
		}
	}
}

func TestPlanDocumentRejectsUnknownManifestRetentionWithoutPartialPlan(t *testing.T) {
	repo := newPlanRepository(t)
	writeInspectionFile(t, repo, manifestPath, `{
  "schemaVersion": 1,
  "generator": {
    "skill": "setup-context-driven",
    "version": 1,
    "baseline": "baseline.unknown"
  },
  "managedArtifacts": []
}
`)
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  planTestDecisions(),
		Preservation: RootPreservationRequest{
			Mode: PreservationModeGreenfield,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan != nil || outcome.Result.State != "action_required" ||
		outcome.Result.Category != "classification" || outcome.Result.NextAction == "" {
		t.Fatalf("unknown retention outcome = %+v", outcome)
	}
}

func buildTestPlan(t *testing.T, repo string) PlanDocument {
	t.Helper()
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  planTestDecisions(),
		Preservation: RootPreservationRequest{
			Mode: PreservationModeGreenfield,
		},
	})
	if err != nil {
		t.Fatalf("build plan: %v", err)
	}
	if outcome.Plan == nil {
		t.Fatalf("build plan returned result: %+v", outcome.Result)
	}
	return *outcome.Plan
}

func planTestDecisions() []DecisionValue {
	return []DecisionValue{
		{ID: "language.generated", Value: "English"},
		{ID: "verification.gate", Value: "make verify"},
		{ID: "spec.scaffold", Value: true},
		{ID: "domain.layout", Value: "single-context"},
		{ID: "triage.external", Value: false},
		{ID: "autonomous.enabled", Value: true},
		{ID: "runtime.backend", Value: "codex gpt-5.5 xhigh"},
		{ID: "runtime.design", Value: "claude opus xhigh"},
		{ID: "secondbrain.enabled", Value: false},
		{ID: "repository.extension.enabled", Value: false},
	}
}

func newPlanRepository(t *testing.T) string {
	t.Helper()
	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, ".agents/skills/context7/SKILL.md", "# context7\n")
	writeInspectionFile(t, repo, ".agents/skills/exa-web-search/SKILL.md", "# exa\n")
	writeInspectionFile(t, repo, "Makefile", "verify:\n\t@true\n")
	commitInspectionRepository(t, repo, "seed portable plan")
	return repo
}

func runPlanGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = dir
	command.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_OPTIONAL_LOCKS=0")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
}

func testRepositoryDigest(t *testing.T, repo string) string {
	t.Helper()
	command := exec.Command("git", "status", "--short", "--untracked-files=all")
	command.Dir = repo
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git status: %v\n%s", err, output)
	}
	return string(output)
}
