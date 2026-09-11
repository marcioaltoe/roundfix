// Suite: Spec Consistency Check characterization corpus
// Invariant: report-authored regressions and repository-wide findings remain observable, while active Specs carry no errors.
// Boundary IN: public speccheck API, replay fixtures, and every active and archived repository Spec
// Boundary OUT: CLI rendering and exit-code policy
package speccheck_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

const (
	replay0058QA001  = "replay-0058-qa-001"
	replay0058QA004  = "replay-0058-qa-004"
	replay0056F001   = "replay-0056-f-001"
	replay0056F002   = "replay-0056-f-002"
	replay0060Task03 = "replay-0060-task-03"
)

func TestConstraintReaderCharacterizesGrantCitation(t *testing.T) {
	t.Parallel()

	t.Run("citation forms", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name           string
			citation       string
			recordPath     string
			wantRecordPath string
		}{
			{
				name:           "backticked legacy repository path resolves",
				citation:       "`docs/workflow/authorizations/2026-08-11-characterization.md`",
				recordPath:     "docs/workflow/authorizations/2026-08-11-characterization.md",
				wantRecordPath: "docs/workflow/authorizations/2026-08-11-characterization.md",
			},
			{
				name:           "backticked Spec-contained repository path resolves",
				citation:       "`docs/specs/0114-tooling-row/_authorization.md`",
				recordPath:     "docs/specs/0114-tooling-row/_authorization.md",
				wantRecordPath: "docs/specs/0114-tooling-row/_authorization.md",
			},
			{
				name:           "Spec-relative Markdown link resolves",
				citation:       "[_authorization.md](_authorization.md)",
				recordPath:     "docs/specs/0114-tooling-row/_authorization.md",
				wantRecordPath: "docs/specs/0114-tooling-row/_authorization.md",
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				row := "Tooling authority: applicable — express maintainer authorization recorded in " + tt.citation + "; bounded files: `Makefile`."
				record := typedConstraintAuthorization("approved", "2026-09-09", "0999-other")
				repoRoot, specsRoot, slug := writeToolingRowFixture(t, row, tt.recordPath, record)
				result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StagePRD)
				if err != nil {
					t.Fatalf("CheckStage(StagePRD): %v", err)
				}

				findings := findingsWithCode(result, speccheck.CodeToolingUnauthorized)
				if tt.wantRecordPath == "" {
					if len(findings) != 0 || hasSkip(result, speccheck.CodeToolingUnauthorized, tt.recordPath) {
						t.Fatalf("authorization observation = findings %#v, skips %#v, want no resolved record path", findings, result.Skipped)
					}
					if len(result.Findings) != 0 {
						t.Fatalf("StagePRD findings = %#v, want the unread grant to pass", result.Findings)
					}
					return
				}

				if len(findings) != 1 {
					t.Fatalf("%s findings = %#v, want exactly one resolved record", speccheck.CodeToolingUnauthorized, findings)
				}
				if !hasExactLocation(findings[0], tt.wantRecordPath, 1) {
					t.Fatalf("authorization locations = %#v, want resolved record %q", findings[0].Where, tt.wantRecordPath)
				}
			})
		}
	})

	t.Run("typed validation is keyed to the record role", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name        string
			recordPath  string
			wantUntyped bool
		}{
			{
				name:        "dated record is validated",
				recordPath:  "docs/workflow/authorizations/2026-08-11-characterization.md",
				wantUntyped: true,
			},
			{
				name:        "undated Spec record is validated",
				recordPath:  "docs/specs/0114-tooling-row/_authorization.md",
				wantUntyped: true,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				row := "Tooling authority: applicable — express maintainer authorization recorded in `" + tt.recordPath + "`; bounded files: `Makefile`."
				repoRoot, specsRoot, slug := writeToolingRowFixture(t, row, tt.recordPath, "Authorization for Spec 0114-tooling-row permits changes to Makefile.\n")
				result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StagePRD)
				if err != nil {
					t.Fatalf("CheckStage(StagePRD): %v", err)
				}

				findings := findingsWithCode(result, speccheck.CodeToolingUntyped)
				if gotUntyped := len(findings) != 0; gotUntyped != tt.wantUntyped {
					t.Fatalf("%s findings = %#v, want present = %t", speccheck.CodeToolingUntyped, findings, tt.wantUntyped)
				}
				if tt.wantUntyped && (len(findings) != 1 || !hasExactLocation(findings[0], tt.recordPath, 1)) {
					t.Fatalf("%s findings = %#v, want one finding at %q", speccheck.CodeToolingUntyped, findings, tt.recordPath)
				}
			})
		}
	})
}

func TestConstraintsResolveSpecContainedRecord(t *testing.T) {
	t.Parallel()

	const recordPath = "docs/specs/0114-tooling-row/_authorization.md"
	row := "Tooling authority: applicable — express maintainer authorization recorded in [_authorization.md](_authorization.md); bounded files: `docs/agents/agent-instructions.md`."
	record := typedConstraintAuthorization("approved", "2026-09-09", "0999-other")
	repoRoot, specsRoot, slug := writeToolingRowFixture(t, row, recordPath, record)

	result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StagePRD)
	if err != nil {
		t.Fatalf("CheckStage(StagePRD): %v", err)
	}
	findings := findingsWithCode(result, speccheck.CodeToolingUnauthorized)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want exactly one", speccheck.CodeToolingUnauthorized, findings)
	}
	for _, location := range []speccheck.Location{
		{Path: "docs/specs/0114-tooling-row/_prd.md", Line: 8},
		{Path: recordPath, Line: 1},
	} {
		if !hasExactLocation(findings[0], location.Path, location.Line) {
			t.Errorf("authorization locations = %#v, want %#v", findings[0].Where, location)
		}
	}
}

func TestExternalSpecRootResolutionCharacterization(t *testing.T) {
	t.Parallel()

	const slug = "external-tooling-row"
	projectRoot := t.TempDir()
	specRepositoryRoot := t.TempDir()
	specsRoot := filepath.Join(specRepositoryRoot, "specs")

	writeToolingRowFile(t, projectRoot, "docs/agents/agent-instructions.md", "# Agent instructions\n")
	row := "Tooling authority: applicable — express maintainer authorization recorded in [_authorization.md](_authorization.md); bounded files: `docs/agents/agent-instructions.md`."
	writeToolingRowFile(t, specRepositoryRoot, "specs/"+slug+"/_prd.md", "# External tooling row\n\n## Project Constraints\n\n"+
		"- Identifier strategy: not applicable — no identifier change. Source: `docs/agents/agent-instructions.md`.\n"+
		"- Authentication and HTTP: not applicable — no network boundary. Source: `docs/agents/agent-instructions.md`.\n"+
		"- Active ADR obligations: not applicable — no ADR applies. Source: `docs/agents/agent-instructions.md`.\n"+
		"- "+row+" Source: `docs/agents/agent-instructions.md`.\n")
	writeToolingRowFile(
		t,
		specRepositoryRoot,
		"specs/"+slug+"/_authorization.md",
		typedConstraintAuthorization("approved", "2026-09-09", slug),
	)

	// Task 03 changes this answer by deriving the record location from the resolved Spec Root.
	operationResolution := spec.ReadSpecAuthorization(context.Background(), projectRoot, slug, "")
	if operationResolution.Outcome != spec.AuthorizationUnresolved {
		t.Fatalf("external-root operation resolution = %q, want unresolved: %#v", operationResolution.Outcome, operationResolution.Reason)
	}
	if operationResolution.Record.Source.Path != spec.AuthorizationRecordPath(slug) ||
		operationResolution.Reason.Code != spec.AuthorizationReasonUnreadableRecord {
		t.Fatalf("external-root operation resolution = %#v, want unreadable default-root record", operationResolution)
	}

	// Task 04 changes this answer by resolving the relative citation beside its carrying artifact.
	result, err := speccheck.CheckStage(specsRoot, projectRoot, slug, speccheck.StagePRD)
	if err != nil {
		t.Fatalf("CheckStage(StagePRD): %v", err)
	}
	findings := findingsWithCode(result, speccheck.CodeToolingUnapproved)
	if len(findings) != 1 || len(result.Findings) != 1 {
		t.Fatalf("external-root citation findings = %#v, want one exact-record refusal", result.Findings)
	}
	if !strings.Contains(findings[0].Summary, "does not identify exactly one authorization record") {
		t.Fatalf("external-root citation summary = %q, want exact-record refusal", findings[0].Summary)
	}
}

func TestConstraintsRefuseNonOperativeGrant(t *testing.T) {
	t.Parallel()

	const recordPath = "docs/specs/0114-tooling-row/_authorization.md"
	tests := []struct {
		name      string
		status    string
		granted   string
		consuming string
		wantCode  string
		wantField string
	}{
		{
			name:      "proposed record withholds status",
			status:    "proposed",
			granted:   "null",
			consuming: "0114-tooling-row",
			wantCode:  speccheck.CodeToolingUnapproved,
			wantField: "status",
		},
		{
			name:      "approved record without grant date withholds granted",
			status:    "approved",
			granted:   "null",
			consuming: "0114-tooling-row",
			wantCode:  speccheck.CodeToolingUnapproved,
			wantField: "granted",
		},
		{
			name:      "withdrawn record withholds status",
			status:    "withdrawn",
			granted:   "null",
			consuming: "0114-tooling-row",
			wantCode:  speccheck.CodeToolingUnapproved,
			wantField: "status",
		},
		{
			name:      "different consumer keeps unauthorized refusal",
			status:    "approved",
			granted:   "2026-09-09",
			consuming: "0999-other",
			wantCode:  speccheck.CodeToolingUnauthorized,
		},
		{
			name:      "approved record for asking Spec passes",
			status:    "approved",
			granted:   "2026-09-09",
			consuming: "0114-tooling-row",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			row := "Tooling authority: applicable — express maintainer authorization recorded in `" + recordPath + "`; bounded files: `docs/agents/agent-instructions.md`."
			record := typedConstraintAuthorization(tt.status, tt.granted, tt.consuming)
			repoRoot, specsRoot, slug := writeToolingRowFixture(t, row, recordPath, record)
			result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StagePRD)
			if err != nil {
				t.Fatalf("CheckStage(StagePRD): %v", err)
			}

			if tt.wantCode == "" {
				if len(result.Findings) != 0 {
					t.Fatalf("StagePRD findings = %#v, want none", result.Findings)
				}
				return
			}
			findings := findingsWithCode(result, tt.wantCode)
			if len(findings) != 1 || len(result.Findings) != 1 {
				t.Fatalf("StagePRD findings = %#v, want exactly one %s", result.Findings, tt.wantCode)
			}
			if !hasExactLocation(findings[0], "docs/specs/0114-tooling-row/_prd.md", 8) ||
				!hasExactLocation(findings[0], recordPath, 1) {
				t.Errorf("%s locations = %#v, want citing row and record", tt.wantCode, findings[0].Where)
			}
			if tt.wantField != "" && !strings.Contains(findings[0].Summary, "field "+tt.wantField) {
				t.Errorf("%s summary = %q, want withholding field %q", tt.wantCode, findings[0].Summary, tt.wantField)
			}
		})
	}
}

func TestConstraintsAcceptHonestProposalDeclaration(t *testing.T) {
	t.Parallel()

	t.Run("honest proposed mutation", func(t *testing.T) {
		t.Parallel()

		const recordPath = "docs/specs/0114-tooling-row/_authorization.md"
		row := "Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `docs/agents/agent-instructions.md`."
		record := typedConstraintAuthorization("proposed", "null", "0114-tooling-row")
		repoRoot, specsRoot, slug := writeToolingRowFixture(t, row, recordPath, record)
		result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StagePRD)
		if err != nil {
			t.Fatalf("CheckStage(StagePRD): %v", err)
		}
		if len(result.Findings) != 0 {
			t.Fatalf("StagePRD findings = %#v, want honest proposal to pass", result.Findings)
		}
	})

	t.Run("approved narrow record wins over proposed record", func(t *testing.T) {
		t.Parallel()

		const (
			narrowPath   = "docs/specs/0114-tooling-row/references/narrow-authorization.md"
			proposalPath = "docs/specs/0114-tooling-row/_authorization.md"
		)
		row := "Tooling authority: applicable — express maintainer authorization is recorded in [the narrow grant](references/narrow-authorization.md); bounded files: `docs/agents/agent-instructions.md`. The broader [_authorization.md](_authorization.md) remains proposed."
		repoRoot, specsRoot, slug := writeToolingRowFixture(t, row, narrowPath, typedConstraintAuthorization("approved", "2026-09-09", "0114-tooling-row"))
		writeToolingRowFile(t, repoRoot, proposalPath, typedConstraintAuthorization("proposed", "null", "0114-tooling-row"))
		result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StagePRD)
		if err != nil {
			t.Fatalf("CheckStage(StagePRD): %v", err)
		}
		if len(result.Findings) != 0 {
			t.Fatalf("StagePRD findings = %#v, want approved narrow grant selected", result.Findings)
		}
	})

	t.Run("ambiguous express authorization refuses", func(t *testing.T) {
		t.Parallel()

		const (
			firstPath  = "docs/specs/0114-tooling-row/references/first-authorization.md"
			secondPath = "docs/specs/0114-tooling-row/references/second-authorization.md"
		)
		row := "Tooling authority: applicable — express maintainer authorization is recorded in [the first grant](references/first-authorization.md) and [the second grant](references/second-authorization.md); bounded files: `docs/agents/agent-instructions.md`."
		record := typedConstraintAuthorization("approved", "2026-09-09", "0114-tooling-row")
		repoRoot, specsRoot, slug := writeToolingRowFixture(t, row, firstPath, record)
		writeToolingRowFile(t, repoRoot, secondPath, record)
		result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StagePRD)
		if err != nil {
			t.Fatalf("CheckStage(StagePRD): %v", err)
		}
		findings := findingsWithCode(result, speccheck.CodeToolingUnapproved)
		if len(findings) != 1 || len(result.Findings) != 1 {
			t.Fatalf("StagePRD findings = %#v, want one ambiguous-record refusal", result.Findings)
		}
		for _, recordPath := range []string{firstPath, secondPath} {
			if !hasExactLocation(findings[0], recordPath, 1) {
				t.Errorf("ambiguous record locations = %#v, want %q", findings[0].Where, recordPath)
			}
		}
		if !strings.Contains(findings[0].Summary, "exactly one") {
			t.Errorf("ambiguous record summary = %q, want exact-one refusal", findings[0].Summary)
		}
	})

	t.Run("dated legacy grant remains accepted", func(t *testing.T) {
		t.Parallel()

		const recordPath = "docs/workflow/authorizations/2026-08-09-characterization.md"
		row := "Tooling authority: applicable — express maintainer authorization recorded in `" + recordPath + "`; bounded files: `docs/agents/agent-instructions.md`."
		record := "# Authorization — bounded tooling change\n\n## Consuming Spec\n\n- 0114-tooling-row\n\n## Authorized Paths\n\n- `docs/agents/agent-instructions.md`\n"
		repoRoot, specsRoot, slug := writeToolingRowFixture(t, row, recordPath, record)
		result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StagePRD)
		if err != nil {
			t.Fatalf("CheckStage(StagePRD): %v", err)
		}
		if len(result.Findings) != 0 {
			t.Fatalf("StagePRD findings = %#v, want dated legacy grant to pass", result.Findings)
		}
	})
}

func TestStrictCheckRefusesMissingImplementAuthority(t *testing.T) {
	t.Parallel()

	const recordPath = "docs/specs/0114-tooling-row/_authorization.md"
	row := "Tooling authority: applicable — express maintainer authorization recorded in [_authorization.md](_authorization.md); bounded files: `docs/agents/agent-instructions.md`."
	repoRoot, specsRoot, slug := writeToolingRowFixture(t, row, recordPath, typedConstraintAuthorization("approved", "2026-09-09", "0114-tooling-row"))
	prdPath := filepath.Join(repoRoot, "docs", "specs", slug, "_prd.md")
	prd, err := os.ReadFile(prdPath)
	if err != nil {
		t.Fatalf("read fixture PRD: %v", err)
	}
	writeToolingRowFile(t, repoRoot, "docs/specs/0114-tooling-row/_prd.md", "---\nstatus: active\n---\n\n"+string(prd))
	writeToolingRowFile(t, repoRoot, "docs/specs/0114-tooling-row/_tasks.md", `---
schema: spec-tasks/v1
spec: 0114-tooling-row
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
---

# Task Graph
`)
	writeToolingRowFile(t, repoRoot, "docs/specs/0114-tooling-row/task_01.md", `---
task: task_01
spec: 0114-tooling-row
status: pending
type: backend
---

# Exercise the grant

## Verification

`+"- `true` — expected: passes.\n")

	withoutImplement, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StageTasks)
	if err != nil {
		t.Fatalf("CheckStage(StageTasks) without implement: %v", err)
	}
	speccheck.PromoteGaps(&withoutImplement)
	findings := findingsWithCode(withoutImplement, speccheck.CodeToolingUnapproved)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want one missing-implement refusal", speccheck.CodeToolingUnapproved, findings)
	}
	for _, want := range []string{"implement", recordPath} {
		if !strings.Contains(findings[0].Summary, want) {
			t.Errorf("missing-implement summary = %q, want %q", findings[0].Summary, want)
		}
	}

	writeToolingRowFile(t, repoRoot, recordPath, typedConstraintAuthorizationWithOperations(
		"approved", "2026-09-09", "0114-tooling-row", "implement",
	))
	withImplement, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StageTasks)
	if err != nil {
		t.Fatalf("CheckStage(StageTasks) with implement: %v", err)
	}
	if findings := findingsWithCode(withImplement, speccheck.CodeToolingUnapproved); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want implement authority accepted", speccheck.CodeToolingUnapproved, findings)
	}
}

func typedConstraintAuthorization(status, granted, consuming string) string {
	return "---\n" +
		"status: " + status + "\n" +
		"granted: " + granted + "\n" +
		"action: implement the bounded tooling change\n" +
		"consuming: " + consuming + "\n" +
		"paths:\n" +
		"  - docs/agents/agent-instructions.md\n" +
		"---\n"
}

func typedConstraintAuthorizationWithOperations(status, granted, consuming string, operations ...string) string {
	var record strings.Builder
	fmt.Fprintf(&record, "---\nstatus: %s\ngranted: %s\naction: implement the bounded tooling change\nconsuming: %s\npaths:\n  - docs/agents/agent-instructions.md\noperations:\n", status, granted, consuming)
	for _, operation := range operations {
		fmt.Fprintf(&record, "  - %s\n", operation)
	}
	record.WriteString("---\n")
	return record.String()
}

func TestCheckReplay0060Task03RefusesWorkIndependentVerification(t *testing.T) {
	t.Parallel()

	findingPath := archivedSpeccheckPath(spec.ArchiveKindFinding, "2026-07-31-a-rehearsal-task-can-settle-completed-without-rehearsing.md")
	result := checkFixture(t, replay0060Task03)
	finding := requireReplayFinding(t, findingPath, result, "SC-VERIFY-WORK-INDEPENDENT", "cannot distinguish Task work from no work")
	assertReplayLocations(t, findingPath, finding,
		speccheck.Location{Path: "docs/specs/" + replay0060Task03 + "/task_03.md", Line: 44},
	)
	report := speccheck.RenderText(result, speccheck.VerificationCoverage{})
	for _, want := range []string{
		"SC-VERIFY-WORK-INDEPENDENT",
		"docs/specs/" + replay0060Task03 + "/task_03.md:44",
		"fix: Add a declared Verification command that asserts this Task's own effect.",
	} {
		if !strings.Contains(report, want) {
			t.Errorf("replay of %s: text finding does not contain %q:\n%s", findingPath, want, report)
		}
	}
	provenance := readReplayFile(t, replay0060Task03, "README.md")
	provenancePath := historicalReplayArchivePath(spec.ArchiveKindFinding, "2026-07-31-a-rehearsal-task-can-settle-completed-without-rehearsing.md")
	for _, want := range []string{provenancePath, "exact Verification commands", "status is `pending`"} {
		if !strings.Contains(provenance, want) {
			t.Errorf("replay provenance does not contain %q:\n%s", want, provenance)
		}
	}
}

func TestCheckReplay0060Task03RefusesContradictoryRequirementsAndUndeclaredRehearsal(t *testing.T) {
	t.Parallel()

	const findingPath = "docs/findings/2026-07-31-a-rehearsal-task-can-settle-completed-without-rehearsing.md"
	result := checkFixture(t, replay0060Task03)

	contradiction := requireReplayFinding(t, findingPath, result, speccheck.CodeRequirementContradictory, "commit")
	assertReplayLocations(t, findingPath, contradiction,
		speccheck.Location{Path: "docs/specs/" + replay0060Task03 + "/task_03.md", Line: 13},
		speccheck.Location{Path: "docs/specs/" + replay0060Task03 + "/task_03.md", Line: 15},
	)

	rehearsal := requireReplayFinding(t, findingPath, result, speccheck.CodeRehearsalUndeclared, "Rehearsal Cases")
	assertReplayLocations(t, findingPath, rehearsal,
		speccheck.Location{Path: "docs/specs/" + replay0060Task03 + "/task_03.md", Line: 9},
	)
}

func TestCheckReplay0058QA001FromReport(t *testing.T) {
	t.Parallel()

	reportPath := archivedSpeccheckPath(spec.ArchiveKindSpec, "0058-npm-trusted-publishing-and-release-preflight", "qa", "qa-report-2026-07-31.md")
	result := checkFixture(t, replay0058QA001)
	finding := requireReplayFinding(t, reportPath, result, speccheck.CodeCoverageUnmapped, "Core Feature 2")
	assertReplayLocations(t, reportPath, finding,
		speccheck.Location{Path: "docs/specs/" + replay0058QA001 + "/_prd.md", Line: 17},
		speccheck.Location{Path: "docs/specs/" + replay0058QA001 + "/_techspec.md", Line: 12},
	)
}

func TestCheckReplay0058QA004FromReport(t *testing.T) {
	t.Parallel()

	reportPath := archivedSpeccheckPath(spec.ArchiveKindSpec, "0058-npm-trusted-publishing-and-release-preflight", "qa", "qa-report-2026-08-01.md")
	result := checkFixture(t, replay0058QA004)
	finding := requireReplayFinding(t, reportPath, result, speccheck.CodeVocabularyUndocumented, "publish:")
	assertReplayLocations(t, reportPath, finding,
		speccheck.Location{Path: "docs/specs/" + replay0058QA004 + "/workflow.sh", Line: 4},
		speccheck.Location{Path: "docs/specs/" + replay0058QA004 + "/runbook.md", Line: 1},
	)
	emitted := readReplayFile(t, replay0058QA004, "workflow.sh")
	documented := readReplayFile(t, replay0058QA004, "runbook.md")
	for _, token := range []string{"identity:", "publish:", "registry:", "runtime:", "undetermined:"} {
		if !strings.Contains(emitted, token) {
			t.Errorf("replay of %s: workflow does not emit %q", reportPath, token)
		}
	}
	if strings.Contains(documented, "publish:") {
		t.Errorf("replay of %s: runbook unexpectedly documents publish:", reportPath)
	}
	for _, token := range []string{"identity:", "registry:", "runtime:", "undetermined:"} {
		if !strings.Contains(documented, token) {
			t.Errorf("replay of %s: runbook does not document %q", reportPath, token)
		}
	}
}

func TestCheckReplay0056F001FromReport(t *testing.T) {
	t.Parallel()

	reportPath := archivedSpeccheckPath(spec.ArchiveKindSpec, "0056-profiles-configure-merge-semantics", "qa", "qa-report-2026-08-01.md")
	result := checkFixture(t, replay0056F001)

	unlisted := requireReplayFinding(t, reportPath, result, speccheck.CodeADRUnlisted, "ADR-0086")
	assertReplayLocations(t, reportPath, unlisted,
		speccheck.Location{Path: "docs/specs/" + replay0056F001 + "/_techspec.md", Line: 12},
		speccheck.Location{Path: "docs/specs/" + replay0056F001 + "/_prd.md", Line: 11},
	)

	related := requireReplayFinding(t, reportPath, result, speccheck.CodeADRRelated, "ADR-0055")
	assertReplayLocations(t, reportPath, related,
		speccheck.Location{Path: "docs/adr/0055-exact-capability-proof.md", Line: 7},
		speccheck.Location{Path: "docs/specs/" + replay0056F001 + "/_prd.md", Line: 11},
	)
	relatedADR := readFixtureRepositoryFile(t, "docs", "adr", "0055-exact-capability-proof.md")
	for _, citation := range []string{"ADR-0039", "ADR-0049"} {
		if !strings.Contains(relatedADR, citation) {
			t.Errorf("replay of %s: related ADR does not cite %s", reportPath, citation)
		}
	}
}

func TestCheckReplay0056F002FromReport(t *testing.T) {
	t.Parallel()

	reportPath := archivedSpeccheckPath(spec.ArchiveKindSpec, "0056-profiles-configure-merge-semantics", "qa", "qa-report-2026-08-01.md")
	result := checkFixture(t, replay0056F002)
	finding := requireReplayFinding(t, reportPath, result, speccheck.CodeCoverageUnmapped, "Core Feature 6")
	assertReplayLocations(t, reportPath, finding,
		speccheck.Location{Path: "docs/specs/" + replay0056F002 + "/_prd.md", Line: 16},
		speccheck.Location{Path: "docs/specs/" + replay0056F002 + "/_techspec.md", Line: 14},
	)
	if techSpec := readReplayFile(t, replay0056F002, "_techspec.md"); !strings.Contains(techSpec, "only added or replaced categories") {
		t.Errorf("replay of %s: TechSpec does not record the narrowed proof scope", reportPath)
	}
}

func TestCheckReplayReadmeProvenance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		slug       string
		reportPath string
	}{
		{
			name:       "0058 QA-001 report",
			slug:       replay0058QA001,
			reportPath: historicalReplayArchivePath(spec.ArchiveKindSpec, "0058-npm-trusted-publishing-and-release-preflight", "qa", "qa-report-2026-07-31.md"),
		},
		{
			name:       "0058 QA-004 report",
			slug:       replay0058QA004,
			reportPath: historicalReplayArchivePath(spec.ArchiveKindSpec, "0058-npm-trusted-publishing-and-release-preflight", "qa", "qa-report-2026-08-01.md"),
		},
		{
			name:       "0056 F-001 report",
			slug:       replay0056F001,
			reportPath: historicalReplayArchivePath(spec.ArchiveKindSpec, "0056-profiles-configure-merge-semantics", "qa", "qa-report-2026-08-01.md"),
		},
		{
			name:       "0056 F-002 report",
			slug:       replay0056F002,
			reportPath: historicalReplayArchivePath(spec.ArchiveKindSpec, "0056-profiles-configure-merge-semantics", "qa", "qa-report-2026-08-01.md"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(fixtureSpecRoot, tt.slug, "README.md")
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read replay provenance %q: %v", path, err)
			}
			for _, required := range []string{tt.reportPath, "authored from the report", "not recovered from Git"} {
				if !strings.Contains(string(content), required) {
					t.Errorf("%s does not record %q", path, required)
				}
			}
		})
	}
}

// historicalReplayArchivePath pins provenance authored before the docs/history
// relocation; unlike replay source lookup, it must not follow ArchiveDir.
func historicalReplayArchivePath(kind spec.ArchiveKind, elements ...string) string {
	parts := make([]string, 0, len(elements)+2)
	parts = append(parts, "_archived", string(kind))
	parts = append(parts, elements...)
	return filepath.ToSlash(filepath.Join(parts...))
}

func requireReplayFinding(t *testing.T, reportPath string, result speccheck.Result, code, summaryFragment string) speccheck.Finding {
	t.Helper()

	findings := findingsWithCode(result, code)
	if len(findings) != 1 {
		t.Fatalf("replay of %s: %s findings = %#v, want exactly one", reportPath, code, findings)
	}
	if !strings.Contains(findings[0].Summary, summaryFragment) {
		t.Fatalf("replay of %s: summary = %q, want %q", reportPath, findings[0].Summary, summaryFragment)
	}
	return findings[0]
}

func assertReplayLocations(t *testing.T, reportPath string, finding speccheck.Finding, want ...speccheck.Location) {
	t.Helper()

	for _, location := range want {
		if !hasExactLocation(finding, location.Path, location.Line) {
			t.Errorf("replay of %s: %s locations = %#v, want %#v", reportPath, finding.Code, finding.Where, location)
		}
	}
}

func readReplayFile(t *testing.T, slug, name string) string {
	t.Helper()

	return readFixtureRepositoryFile(t, "docs", "specs", slug, name)
}

func readFixtureRepositoryFile(t *testing.T, pathElements ...string) string {
	t.Helper()

	path := filepath.Join(append([]string{"testdata", "repo"}, pathElements...)...)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read replay fixture %q: %v", path, err)
	}
	return string(content)
}

func activeCorpusPRD(content []byte) ([]byte, error) {
	const statusPrefix = "\nstatus:"
	statusStart := bytes.Index(content, []byte(statusPrefix))
	if statusStart < 0 {
		return nil, errors.New("frontmatter has no status field")
	}
	statusStart++
	statusEnd := bytes.IndexByte(content[statusStart:], '\n')
	if statusEnd < 0 {
		return nil, errors.New("frontmatter status field has no line ending")
	}
	statusEnd += statusStart

	active := make([]byte, 0, len(content))
	active = append(active, content[:statusStart]...)
	active = append(active, "status: active"...)
	active = append(active, content[statusEnd:]...)
	return active, nil
}

func copyCorpusFile(sourcePath, targetPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()
	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(target, source); err != nil {
		target.Close()
		return err
	}
	return target.Close()
}
