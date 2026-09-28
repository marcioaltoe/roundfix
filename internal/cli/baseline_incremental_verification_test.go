// Suite: incremental Verification migration
// Invariant: a single-gate Setup Manifest never updates silently and can adopt only a locally declared incremental command.
// Boundary IN: baseline update and plan exits, output, Setup Manifest projection, managed guidance, and repository writes.
// Boundary OUT: catalog and alignment unit contracts, which stay in internal/baseline.

package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestBaselineUpdateNamesTheMissingIncrementalDecision(t *testing.T) {
	t.Parallel()

	repository := singleGateBaselineUpdateRepository(t, true)
	before := baselinePlanTestTree(t, repository)
	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t, context.Background(),
		"baseline", "update", "--repo", repository, "--format=json",
	)
	if code != exitUnverified || stderr != "" {
		t.Fatalf("single-gate update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if result.State != "action_required" || result.Category != "decision" ||
		len(result.NewDecisions) != 1 || result.NewDecisions[0].ID != "verification.incremental" ||
		result.NewDecisions[0].SuggestedValue != "rtk make verify-incremental" {
		t.Fatalf("single-gate update result = %+v", result)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("single-gate update changed repository bytes")
	}
}

func TestBaselineUpdateRefusesAnUndeclaredIncrementalSuggestion(t *testing.T) {
	t.Parallel()

	repository := singleGateBaselineUpdateRepository(t, false)
	before := baselinePlanTestTree(t, repository)
	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t, context.Background(),
		"baseline", "update", "--repo", repository,
		"--adopt-suggested", "--yes", "--format=json",
	)
	if code != exitUnverified || stderr != "" {
		t.Fatalf("undeclared suggestion update exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if result.State != "action_required" || result.Category != "decision" ||
		!strings.Contains(result.Message, "verification.incremental") {
		t.Fatalf("undeclared suggestion update result = %+v", result)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("undeclared suggestion update changed repository bytes")
	}
}

func TestBaselineUpdateAdoptsADeclaredIncrementalSuggestion(t *testing.T) {
	t.Parallel()

	repository := singleGateBaselineUpdateRepository(t, true)
	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t, context.Background(),
		"baseline", "update", "--repo", repository,
		"--yes", "--adopt-suggested", "--format=json",
	)
	if code != exitOK || stderr != "" || result.State != "verified" {
		t.Fatalf("declared suggestion update exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
	}

	manifest := ReadBaselineSetupManifest(t, repository)
	decision, found := manifest.Decisions["verification.incremental"]
	if !found || decision.Value != "rtk make verify-incremental" {
		t.Fatalf("incremental manifest decision = %#v, found=%v", decision, found)
	}
	projectionCount := 0
	for _, projection := range manifest.Verification {
		if projection.ID != "verification.incremental" {
			continue
		}
		projectionCount++
		if projection.Role != "incremental" || projection.Classification != baseline.VerificationRepositoryCommand ||
			!projection.RepositoryExecutable || projection.DeclarationPath != "Makefile" {
			t.Fatalf("incremental manifest projection = %+v", projection)
		}
	}
	if projectionCount != 1 {
		t.Fatalf("incremental manifest projections = %d, want one", projectionCount)
	}
	guidance, err := os.ReadFile(filepath.Join(repository, "docs", "agents", "agent-instructions.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(guidance), "The selected incremental Verification is `rtk make verify-incremental`.") {
		t.Fatalf("agent instructions do not publish incremental command:\n%s", guidance)
	}

	current, currentOut, currentErr, currentCode := runBaselineUpdateTestCommand(
		t, context.Background(),
		"baseline", "update", "--repo", repository, "--format=json",
	)
	if currentCode != exitOK || currentErr != "" || current.State != "current" || len(current.FileChanges) != 0 {
		t.Fatalf("second update exit=%d result=%+v stdout=%s stderr=%s", currentCode, current, currentOut, currentErr)
	}
}

func TestBaselinePlanWithoutIncrementalDecisionExitsThree(t *testing.T) {
	t.Parallel()

	repository := newBaselineApplyTestRepository(t)
	before := baselinePlanTestTree(t, repository)
	args := []string{
		"baseline", "plan", "--repo", repository, "--profile", "go-cli-tui",
		"--decision", "preservation.mode=greenfield",
		"--decision", "language.generated=English",
		"--decision", "verification.gate=make verify",
		"--decision", "branch.prefix=ma/",
		"--decision", "spec.scaffold=true",
		"--decision", "domain.layout=single-context",
		"--decision", "triage.external=false",
		"--decision", "autonomous.enabled=true",
		"--decision", "runtime.backend=codex gpt-5.5 xhigh",
		"--decision", "runtime.design=claude opus xhigh",
		"--decision", "secondbrain.enabled=false",
		"--decision", "repository.extension.enabled=false",
		"--format=json",
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), args, &stdout, &stderr)
	if code != exitUnverified || stderr.Len() != 0 {
		t.Fatalf("missing incremental plan exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	result, err := baseline.ParseResult(stdout.Bytes())
	if err != nil {
		t.Fatalf("parse missing incremental result: %v\n%s", err, stdout.String())
	}
	if result.State != "action_required" || result.Category != "decision" ||
		!strings.Contains(result.Message, "verification.incremental") {
		t.Fatalf("missing incremental plan result = %+v", result)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatal("missing incremental plan changed repository bytes")
	}
}

func singleGateBaselineUpdateRepository(t *testing.T, declareIncremental bool) string {
	t.Helper()
	repository := newBaselineUpdateRepository(t)
	makefile := "verify:\n\t@true\n"
	if declareIncremental {
		makefile += "verify-incremental:\n\t@true\n"
	}
	writeBaselinePlanTestFile(t, repository, "Makefile", makefile)

	manifest := ReadBaselineSetupManifest(t, repository)
	delete(manifest.Decisions, "verification.incremental")
	verification := make([]baseline.VerificationProjection, 0, len(manifest.Verification)-1)
	for _, projection := range manifest.Verification {
		if projection.ID != "verification.incremental" {
			verification = append(verification, projection)
		}
	}
	manifest.Verification = verification
	WriteBaselineSetupManifest(t, repository, manifest)
	commitBaselinePlanTestRepository(t, repository)
	return repository
}
