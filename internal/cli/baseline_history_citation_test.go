// Suite: Relocation Citation CLI contracts
// Invariant: public Baseline planning renders digest-bound citation warnings and stale confirmation never writes.
// Boundary IN: baseline plan and baseline update through their public command runners over real Git repositories.
// Boundary OUT: citation parsing details and portable-plan apply internals.

package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

const (
	cliCitationSource = "docs/adr/0040-retired.md"
	cliCitationPath   = "docs/references/layout.md"
)

func TestBaselinePlanPrintsRelocationCitationWarnings(t *testing.T) {
	t.Parallel()
	repository := newBaselineHistoryCitationRepository(t)
	want := "Warning: baseline.history.citation: " + cliCitationPath + ": line 1 cites ../adr/0040-retired.md"

	t.Run("text", func(t *testing.T) {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunContext(
			context.Background(),
			baselineHistoryCitationPlanArgs(repository, "text"),
			&stdout,
			&stderr,
		)
		if code != exitOK || stderr.Len() != 0 {
			t.Fatalf("baseline plan text exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("baseline plan text missing %q:\n%s", want, stdout.String())
		}
	})

	t.Run("json", func(t *testing.T) {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunContext(
			context.Background(),
			baselineHistoryCitationPlanArgs(repository, "json"),
			&stdout,
			&stderr,
		)
		if code != exitOK || stderr.Len() != 0 {
			t.Fatalf("baseline plan JSON exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		plan, err := baseline.ParsePlanDocument(stdout.Bytes())
		if err != nil {
			t.Fatalf("parse Baseline Plan JSON: %v\n%s", err, stdout.String())
		}
		for _, warning := range plan.Warnings {
			if warning.Code == "baseline.history.citation" &&
				warning.Path == cliCitationPath &&
				strings.Contains(warning.Message, "line 1 cites ../adr/0040-retired.md") {
				return
			}
		}
		t.Fatalf("Baseline Plan JSON warnings = %#v, want Relocation Citation warning", plan.Warnings)
	})
}

func TestBaselineUpdateRefusesADigestWhoseCitationsChanged(t *testing.T) {
	t.Parallel()
	repository := newBaselineUpdateRepository(t)
	writeBaselinePlanTestFile(t, repository, cliCitationSource, "---\nstatus: superseded\n---\n")
	writeBaselinePlanTestFile(t, repository, cliCitationPath, "See [the retired decision](../adr/0040-retired.md).\n")
	commitBaselinePlanTestRepository(t, repository)

	preview, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository, "--format=json",
	)
	if code != exitUnverified || stderr != "" || preview.State != "plan_ready" || preview.PlanDigest == "" {
		t.Fatalf("citation-bound preview exit=%d result=%+v stdout=%s stderr=%s", code, preview, stdout, stderr)
	}

	writeBaselinePlanTestFile(t, repository, cliCitationPath, "No retired decision citation.\n")
	before := baselinePlanTestTree(t, repository)
	result, stdout, stderr, code := runBaselineUpdateTestCommand(
		t,
		context.Background(),
		"baseline", "update", "--repo", repository,
		"--confirm-plan", preview.PlanDigest, "--format=json",
	)
	if code != exitUnverified || !strings.Contains(stderr, "does not match") {
		t.Fatalf("stale citation confirmation exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
	}
	if result.State != "action_required" || result.Category != "approval" ||
		result.PlanDigest == "" || result.PlanDigest == preview.PlanDigest || result.ApprovedPlanDigest != "" {
		t.Fatalf("stale citation confirmation result = %+v", result)
	}
	if after := baselinePlanTestTree(t, repository); after != before {
		t.Fatalf("stale citation confirmation changed repository bytes: before=%s after=%s", before, after)
	}
}

func newBaselineHistoryCitationRepository(t *testing.T) string {
	t.Helper()
	repository := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repository, ".agents/skills/context7/SKILL.md", "# context7\n")
	writeBaselinePlanTestFile(t, repository, ".agents/skills/exa-web-search/SKILL.md", "# exa\n")
	writeBaselinePlanTestFile(t, repository, "Makefile", "verify:\n\t@true\nverify-incremental:\n\t@true\n")
	writeBaselinePlanTestFile(t, repository, cliCitationSource, "---\nstatus: superseded\n---\n")
	writeBaselinePlanTestFile(t, repository, cliCitationPath, "See [the retired decision](../adr/0040-retired.md).\n")
	commitBaselinePlanTestRepository(t, repository)
	return repository
}

func baselineHistoryCitationPlanArgs(repository, format string) []string {
	args := baselinePlanCharacterizationArgs(repository, "greenfield")
	args[len(args)-1] = "--format=" + format
	return args
}
