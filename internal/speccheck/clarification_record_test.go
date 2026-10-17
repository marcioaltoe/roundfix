// IN: Check and CheckStage against disposable Git repositories.
// OUT: real repository history, interviews, and Daemon settlement.
// Invariant: PRDs at the guide horizon record their clarification route.
package speccheck

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/gittest"
)

const clarificationFixtureSlug = "0200-clarification"
const clarificationFixturePRD = "docs/specs/" + clarificationFixtureSlug + "/_prd.md"

func clarificationRepo(t *testing.T, prd string) string {
	t.Helper()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	writeReceiptFixture(t, root, ClarificationGuidePath, "guide\n")
	receiptCommit(t, root)
	writeReceiptFixture(t, root, clarificationFixturePRD, prd)
	return root
}

func clarificationCheck(t *testing.T, root string, stage Stage) Result {
	t.Helper()
	result, err := CheckStage(filepath.Join(root, "docs/specs"), root, clarificationFixtureSlug, stage)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestClarificationRecordGapNamesAPRDWithoutTheRecord(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, prd string
		line      int
	}{
		{"no section", "# Example\n", 1},
		{"neither shape", "# Example\n\n## Decisions\n\n- The author decided.\n", 3},
		{"invalid date", "# Example\n\n## Decisions\n\n- 2026-02-30: Question — answer\n", 3},
		{"answer outside section", "- 2026-10-17: Question — answer\n\n## Decisions\n\n## Other\nAutonomous route: no person present; decisions taken from the source.\n", 3},
		{"autonomous source missing", "## Decisions\nAutonomous route: no person present; decisions taken from .\n", 1},
		{"autonomous period missing", "## Decisions\nAutonomous route: no person present; decisions taken from the source\n", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := clarificationRepo(t, tc.prd)
			receiptCommit(t, root)
			findings := receiptCodeFindings(clarificationCheck(t, root, StagePRD), CodeDecisionsUnrecorded)
			want := []Finding{{
				Code: CodeDecisionsUnrecorded, Severity: SeverityGap,
				Summary: clarificationFixturePRD + " has no Clarification Record",
				Where:   []Location{{Path: clarificationFixturePRD, Line: tc.line}},
				Fix:     "Under ## Decisions, write one - YYYY-MM-DD: <question> — <answer> line for each question a person answered, or the line Autonomous route: no person present; decisions taken from <source>., following the write-prd rule \"Clarification Record\".",
			}}
			if !reflect.DeepEqual(findings, want) {
				t.Fatalf("findings = %#v, want %#v", findings, want)
			}
		})
	}
}

func TestClarificationRecordAcceptsADatedAnswerOrTheAutonomousRoute(t *testing.T) {
	t.Parallel()
	const dated = "- 2024-02-29: Which route? — The person chose the existing route.\n"
	const autonomous = "Autonomous route: no person present; decisions taken from the accepted decision record.\n"
	for _, tc := range []struct{ name, record string }{
		{"dated", dated}, {"autonomous", autonomous}, {"both", dated + autonomous},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := clarificationRepo(t, "## Decisions\n\n"+tc.record)
			receiptCommit(t, root)
			for _, stage := range []Stage{StagePRD, StageTechSpec, StageTasks, StageAll} {
				if found := receiptCodeFindings(clarificationCheck(t, root, stage), CodeDecisionsUnrecorded); len(found) != 0 {
					t.Fatalf("stage %q findings = %#v", stage, found)
				}
			}
		})
	}
}

func TestClarificationRecordSkipsAPRDBeforeItsHorizon(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name                           string
		guide, uncommitted, missingPRD bool
	}{
		{"before guide", true, false, false},
		{"without guide", false, false, false},
		{"uncommitted", true, true, false},
		{"missing PRD", true, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			gittest.InitRepo(t, root, "-b", "main")
			if !tc.uncommitted {
				writeReceiptFixture(t, root, clarificationFixturePRD, "# Legacy\n")
				receiptCommit(t, root)
			}
			if tc.guide {
				writeReceiptFixture(t, root, ClarificationGuidePath, "guide\n")
				receiptCommit(t, root)
			}
			if tc.uncommitted {
				writeReceiptFixture(t, root, clarificationFixturePRD, "# New\n")
			}
			if tc.missingPRD {
				if err := os.Remove(filepath.Join(root, clarificationFixturePRD)); err != nil {
					t.Fatal(err)
				}
			}
			for _, stage := range []Stage{StagePRD, StageAll} {
				result := clarificationCheck(t, root, stage)
				found := receiptCodeFindings(result, CodeDecisionsUnrecorded)
				if tc.uncommitted {
					if len(found) != 1 {
						t.Fatalf("uncommitted findings = %#v, want one gap", found)
					}
					continue
				}
				if len(found) != 0 {
					t.Fatalf("skipped findings = %#v", found)
				}
				missing := "a PRD committed at or after " + ClarificationGuidePath
				if tc.missingPRD {
					missing = clarificationFixturePRD
				}
				skipped := false
				for _, skip := range result.Skipped {
					skipped = skipped || (skip.Code == CodeDecisionsUnrecorded && skip.Missing == missing)
				}
				if !skipped {
					t.Fatalf("stage %q skips = %#v, want %q", stage, result.Skipped, missing)
				}
			}
		})
	}

	t.Run("guide and PRD added together", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		gittest.InitRepo(t, root, "-b", "main")
		writeReceiptFixture(t, root, ClarificationGuidePath, "guide\n")
		writeReceiptFixture(t, root, clarificationFixturePRD, "# New\n")
		receiptCommit(t, root)
		if found := receiptCodeFindings(clarificationCheck(t, root, StagePRD), CodeDecisionsUnrecorded); len(found) != 1 {
			t.Fatalf("same-commit findings = %#v, want one gap", found)
		}
	})
	t.Run("oldest guide addition holds an intermediate PRD", func(t *testing.T) {
		t.Parallel()
		root := clarificationRepo(t, "# New\n")
		receiptCommit(t, root)
		if err := os.Remove(filepath.Join(root, ClarificationGuidePath)); err != nil {
			t.Fatal(err)
		}
		receiptCommit(t, root)
		writeReceiptFixture(t, root, ClarificationGuidePath, "reintroduced guide\n")
		receiptCommit(t, root)
		if found := receiptCodeFindings(clarificationCheck(t, root, StagePRD), CodeDecisionsUnrecorded); len(found) != 1 {
			t.Fatalf("reintroduced-guide findings = %#v, want one gap", found)
		}
	})
}

func TestClarificationRecordRunsFromThePRDStageAndFailsStrict(t *testing.T) {
	t.Parallel()
	root := clarificationRepo(t, "# New\n")
	receiptCommit(t, root)
	for _, stage := range []Stage{StagePRD, StageTechSpec, StageTasks, StageAll} {
		result := clarificationCheck(t, root, stage)
		found := receiptCodeFindings(result, CodeDecisionsUnrecorded)
		if len(found) != 1 || found[0].Severity != SeverityGap {
			t.Fatalf("stage %q findings = %#v, want one gap", stage, found)
		}
		for _, skip := range result.Skipped {
			if skip.Code == CodeDecisionsUnrecorded {
				t.Fatalf("stage %q skipped the detector: %#v", stage, skip)
			}
		}
		PromoteGaps(&result)
		if receiptCodeFindings(result, CodeDecisionsUnrecorded)[0].Severity != SeverityError {
			t.Fatalf("stage %q strict result = %#v, want error", stage, result)
		}
	}
}
