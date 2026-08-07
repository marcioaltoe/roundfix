// Suite: manifest-driven Baseline update command
// Invariant: update either presents a digest-bound managed refresh without writes or applies that exact plan.
// Boundary IN: CLI parsing, manifest projection, managed-refresh planning, apply, structured output, and exit categories.
// Boundary OUT: skill refresh and interactive Baseline adoption, which later Tasks own.

package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestBaselineUpdateAppliesManifestPlanAndReportsJSON(t *testing.T) {
	repository := staleBaselineUpdateRepository(t)
	before := baselinePlanTestTree(t, repository)

	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--yes", "--format=json",
	)
	if code != exitOK || stderr != "" {
		t.Fatalf("baseline update exit = %d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if result.SchemaVersion != baselineUpdateResultSchema ||
		result.Operation != "update" ||
		result.State != "verified" ||
		result.ApprovedPlanDigest == "" ||
		result.ApprovedPlanDigest != result.PlanDigest {
		t.Fatalf("baseline update result = %+v", result)
	}
	if result.PriorCatalog.Digest == "" ||
		result.CurrentCatalog.Digest == "" ||
		result.PriorCatalog.Digest == result.CurrentCatalog.Digest {
		t.Fatalf("baseline update catalog identities = prior %+v current %+v", result.PriorCatalog, result.CurrentCatalog)
	}
	if len(result.FileChanges) < 2 {
		t.Fatalf("baseline update file changes = %+v, want stale guide and manifest", result.FileChanges)
	}
	if result.Retention == nil || result.Warnings == nil || result.AdoptedSuggestions == nil {
		t.Fatalf("baseline update JSON arrays must be present: %+v", result)
	}
	if before == baselinePlanTestTree(t, repository) {
		t.Fatal("baseline update did not rewrite stale managed artifacts")
	}
}

func TestBaselineUpdateIdempotenceReportsZeroFileChanges(t *testing.T) {
	repository := staleBaselineUpdateRepository(t)

	first, firstOut, firstErr, firstCode := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--yes", "--format=json",
	)
	if firstCode != exitOK || firstErr != "" || len(first.FileChanges) == 0 {
		t.Fatalf("first baseline update exit=%d result=%+v stdout=%s stderr=%s", firstCode, first, firstOut, firstErr)
	}
	afterFirst := baselinePlanTestTree(t, repository)

	second, secondOut, secondErr, secondCode := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--yes", "--format=json",
	)
	if secondCode != exitOK || secondErr != "" {
		t.Fatalf("second baseline update exit=%d stdout=%s stderr=%s", secondCode, secondOut, secondErr)
	}
	if len(second.FileChanges) != 0 || second.State != "verified" {
		t.Fatalf("second baseline update result = %+v, want verified zero-change result", second)
	}
	if second.StatusMatrix == nil || second.StatusMatrix.Idempotence != baseline.EvidenceStatusVerified {
		t.Fatalf("second baseline update idempotence evidence = %+v", second.StatusMatrix)
	}
	if afterSecond := baselinePlanTestTree(t, repository); afterSecond != afterFirst {
		t.Fatal("second baseline update changed repository bytes")
	}
}

func TestBaselineUpdateNoManifestRequiresAdoptionWithoutWrites(t *testing.T) {
	repository := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repository, "README.md", "unadopted repository\n")
	commitBaselinePlanTestRepository(t, repository)
	before := baselinePlanTestTree(t, repository)

	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--format=json",
	)
	if code != exitUnverified || stderr != "" {
		t.Fatalf("no-manifest update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if result.State != "action_required" || result.Category != "adoption" ||
		!strings.Contains(result.NextAction, "roundfix baseline") ||
		!strings.Contains(strings.ToLower(result.NextAction), "adoption") {
		t.Fatalf("no-manifest update result = %+v", result)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("no-manifest update changed repository bytes")
	}
}

func TestBaselineUpdateNewDecisionRequiresActionWithoutWrites(t *testing.T) {
	repository := incompleteBaselineUpdateRepository(t, "secondbrain.enabled")
	before := baselinePlanTestTree(t, repository)

	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--format=json",
	)
	if code != exitUnverified || stderr != "" {
		t.Fatalf("new-decision update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if result.State != "action_required" || result.Category != "decision" ||
		len(result.NewDecisions) != 1 || result.NewDecisions[0].ID != "secondbrain.enabled" ||
		!strings.Contains(result.Message, "secondbrain.enabled") {
		t.Fatalf("new-decision update result = %+v", result)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("new-decision refusal changed repository bytes")
	}
}

func TestBaselineUpdateAdoptsEverySuggestedDecision(t *testing.T) {
	repository := incompleteBaselineUpdateRepository(t, "secondbrain.enabled")

	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository,
		"--adopt-suggested", "--yes", "--format=json",
	)
	if code != exitOK || stderr != "" {
		t.Fatalf("suggestion-adopting update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if result.State != "verified" || len(result.AdoptedSuggestions) != 1 ||
		result.AdoptedSuggestions[0].ID != "secondbrain.enabled" ||
		result.AdoptedSuggestions[0].SuggestedValue != true {
		t.Fatalf("suggestion-adopting result = %+v", result)
	}
}

func TestBaselineUpdatePresentsPlanAndConfirmsPreviousDigest(t *testing.T) {
	repository := staleBaselineUpdateRepository(t)
	before := baselinePlanTestTree(t, repository)

	planned, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--format=text",
	)
	if code != exitUnverified || stderr != "" {
		t.Fatalf("unconfirmed update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	for _, want := range []string{
		"Baseline update: plan ready",
		"Prior catalog:",
		"Current catalog:",
		"File changes:",
		"Retention evidence:",
		"Plan Digest:",
		"Next action:",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("unconfirmed update output missing %q:\n%s", want, stdout)
		}
	}
	if planned.PlanDigest == "" || planned.ApprovedPlanDigest != "" {
		t.Fatalf("unconfirmed update result = %+v", planned)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("unconfirmed update changed repository bytes")
	}

	confirmed, confirmOut, confirmErr, confirmCode := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository,
		"--confirm-plan", planned.PlanDigest, "--format=json",
	)
	if confirmCode != exitOK || confirmErr != "" || confirmed.ApprovedPlanDigest != planned.PlanDigest {
		t.Fatalf("previous-digest confirmation exit=%d result=%+v stdout=%s stderr=%s", confirmCode, confirmed, confirmOut, confirmErr)
	}
}

func TestBaselineUpdateRejectsMutuallyExclusiveConfirmationForms(t *testing.T) {
	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--yes", "--confirm-plan", "sha256:reviewed", "--format=json",
	)
	if code != exitPreflight || !strings.Contains(stderr, "mutually exclusive") {
		t.Fatalf("mutually exclusive update exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
	}
	if result.State != "failed" || result.Category != "invalid" {
		t.Fatalf("mutually exclusive update result = %+v", result)
	}
}

func TestBaselineUpdateRejectsDigestOtherThanCurrentPlanWithoutWrites(t *testing.T) {
	repository := staleBaselineUpdateRepository(t)
	before := baselinePlanTestTree(t, repository)
	const staleDigest = "sha256:0000000000000000000000000000000000000000000000000000000000000000"

	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository,
		"--confirm-plan", staleDigest, "--format=json",
	)
	if code != exitUnverified || !strings.Contains(stderr, "does not match") {
		t.Fatalf("stale confirmation exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
	}
	if result.State != "action_required" || result.Category != "approval" ||
		result.PlanDigest == "" || result.PlanDigest == staleDigest || result.ApprovedPlanDigest != "" {
		t.Fatalf("stale confirmation result = %+v", result)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("stale confirmation changed repository bytes")
	}
}

func TestBaselineUpdateHelpNamesNonInteractiveContract(t *testing.T) {
	_, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--help",
	)
	if code != exitOK || stderr != "" {
		t.Fatalf("baseline update help exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	for _, want := range []string{
		"roundfix baseline update",
		"--yes | --confirm-plan <digest>",
		"--adopt-suggested",
		baselineUpdateResultSchema,
		"without prompting or invoking a semantic analyzer",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("baseline update help missing %q:\n%s", want, stdout)
		}
	}
}

func TestBaselineUpdateExitCategoriesAndNoACPRuntimeDependency(t *testing.T) {
	repository := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repository, "README.md", "unadopted repository\n")
	commitBaselinePlanTestRepository(t, repository)

	var stderr bytes.Buffer
	code := RunContext(
		context.Background(),
		[]string{"baseline", "update", "--repo", repository, "--format=json"},
		failingWriter{err: errors.New("injected baseline update output failure")},
		&stderr,
	)
	if code != exitRunFailed || !strings.Contains(stderr.String(), "injected baseline update output failure") {
		t.Fatalf("output failure exit=%d stderr=%s", code, stderr.String())
	}

	adopted := newBaselineUpdateRepository(t)
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	_, stdout, canceledErr, canceledCode := runBaselineUpdateTestCommand(
		t,
		canceled,
		"baseline", "update", "--repo", adopted, "--yes", "--format=json",
	)
	if canceledCode != exitSIGINT || !strings.Contains(canceledErr, "context canceled") {
		t.Fatalf("canceled update exit=%d stdout=%s stderr=%s", canceledCode, stdout, canceledErr)
	}
}

func runBaselineUpdateTestCommand(
	t *testing.T,
	ctx context.Context,
	args ...string,
) (baselineUpdateResult, string, string, int) {
	t.Helper()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(ctx, args, &stdout, &stderr)
	var result baselineUpdateResult
	if index := strings.Index(stdout.String(), "{"); index >= 0 {
		if err := json.Unmarshal(stdout.Bytes()[index:], &result); err != nil {
			t.Fatalf("decode baseline update result: %v\n%s", err, stdout.String())
		}
	} else if digest := baselineUpdateTextValue(stdout.String(), "Plan Digest:"); digest != "" {
		result.PlanDigest = digest
	}
	return result, stdout.String(), stderr.String(), code
}

func baselineUpdateTextValue(output, label string) string {
	for _, line := range strings.Split(output, "\n") {
		if value, found := strings.CutPrefix(line, label); found {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func newBaselineUpdateRepository(t *testing.T) string {
	t.Helper()
	repository := newBaselineApplyTestRepository(t)
	plan, _ := baselineApplyTestPlan(t, repository)
	if _, err := baseline.ApplyPlan(context.Background(), repository, plan, plan.PlanDigest); err != nil {
		t.Fatalf("adopt Baseline update fixture: %v", err)
	}
	commitBaselinePlanTestRepository(t, repository)
	return repository
}

func staleBaselineUpdateRepository(t *testing.T) string {
	t.Helper()
	repository := newBaselineUpdateRepository(t)
	manifest := readBaselineUpdateManifest(t, repository)
	manifest.CatalogDigest = "sha256:" + strings.Repeat("0", 64)

	artifactIndex := -1
	for index, artifact := range manifest.ManagedArtifacts {
		if artifact.Kind == "guide" {
			artifactIndex = index
			break
		}
	}
	if artifactIndex < 0 {
		t.Fatal("adopted fixture has no managed guide")
	}
	artifact := manifest.ManagedArtifacts[artifactIndex]
	const staleBody = "Stale managed guidance from the prior catalog.\n"
	path := filepath.Join(repository, filepath.FromSlash(artifact.Path))
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read managed guide %q: %v", artifact.Path, err)
	}
	begin := "<!-- setup-context-driven:begin id=" + artifact.ID + " version=" + artifact.Version + " -->"
	end := "<!-- setup-context-driven:end id=" + artifact.ID + " -->"
	start := bytes.Index(content, []byte(begin))
	finish := bytes.Index(content, []byte(end))
	if start < 0 || finish < start {
		t.Fatalf("managed guide %q lacks artifact markers for %q", artifact.Path, artifact.ID)
	}
	finish += len(end)
	replacement := []byte(begin + "\n\n" + staleBody + "\n" + end)
	content = append(append(append([]byte(nil), content[:start]...), replacement...), content[finish:]...)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write stale managed guide %q: %v", artifact.Path, err)
	}
	sum := sha256.Sum256([]byte(staleBody))
	manifest.ManagedArtifacts[artifactIndex].Digest = hex.EncodeToString(sum[:])
	writeBaselineUpdateManifest(t, repository, manifest)
	commitBaselinePlanTestRepository(t, repository)
	return repository
}

func incompleteBaselineUpdateRepository(t *testing.T, decisionID string) string {
	t.Helper()
	repository := newBaselineUpdateRepository(t)
	manifest := readBaselineUpdateManifest(t, repository)
	if _, exists := manifest.Decisions[decisionID]; !exists {
		t.Fatalf("adopted fixture manifest has no decision %q", decisionID)
	}
	delete(manifest.Decisions, decisionID)
	writeBaselineUpdateManifest(t, repository, manifest)
	commitBaselinePlanTestRepository(t, repository)
	return repository
}

func readBaselineUpdateManifest(t *testing.T, repository string) baseline.SetupManifest {
	t.Helper()
	path := filepath.Join(repository, "docs", "agents", "setup-context.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read Setup Manifest: %v", err)
	}
	var manifest baseline.SetupManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("decode Setup Manifest: %v", err)
	}
	return manifest
}

func writeBaselineUpdateManifest(t *testing.T, repository string, manifest baseline.SetupManifest) {
	t.Helper()
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("encode Setup Manifest: %v", err)
	}
	data = append(data, '\n')
	path := filepath.Join(repository, "docs", "agents", "setup-context.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write Setup Manifest: %v", err)
	}
}
