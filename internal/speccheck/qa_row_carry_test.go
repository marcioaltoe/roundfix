// Suite: QA row carry dispositions
// Invariant: only proven, unchanged repository observations carry, with their original provenance.
// Boundary IN: mechanical stage, temporary Git histories, seeded report bytes
// Boundary OUT: Daemon prior-pass import and scheduling
package speccheck_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

const rowCarryReportPath = "docs/specs/mechanical/qa/qa-report-2026-08-11.md"

// A side branch holds the audited head and optional recording commit. The
// current branch imports the report bytes, independently recreating the input.
func rowCarrySidePass(t *testing.T, trailer string) (string, string) {
	t.Helper()
	root := newMechanicalGitRepo(t)
	runMechanicalGit(t, root, "checkout", "-q", "-b", "qa-side")
	writeMechanicalFile(t, root, "evidence.txt", "stable evidence\n")
	head := commitMechanicalFiles(t, root, "establish side evidence", "evidence.txt")
	report := mechanicalCarryReport(head, "stable evidence\n")
	writeMechanicalFile(t, root, rowCarryReportPath, report)
	if trailer != "" {
		commitMechanicalFiles(t, root, "docs: record QA\n\n"+trailer, rowCarryReportPath)
	}
	runMechanicalGit(t, root, "checkout", "-q", "main")
	writeMechanicalFile(t, root, "evidence.txt", "stable evidence\n")
	writeMechanicalFile(t, root, rowCarryReportPath, report)
	commitMechanicalFiles(t, root, "import prior bytes", "evidence.txt", rowCarryReportPath)
	return root, head
}

func rowCarryResult(t *testing.T, root string) speccheck.MechanicalResult {
	t.Helper()
	return runMechanical(t, speccheck.MechanicalRequest{RepoRoot: root, ReportPath: rowCarryReportPath})
}

func rowCarryRefusal(t *testing.T, result speccheck.MechanicalResult, reason string) {
	t.Helper()
	if len(result.Carried) != 0 || len(result.Dispositions) != 1 || result.Dispositions[0].ID != "R01" || result.Dispositions[0].Carried || result.Dispositions[0].Reason != reason {
		t.Fatalf("unexpected carry result: %#v", result)
	}
	var out bytes.Buffer
	if err := speccheck.WriteMechanicalResult(&out, result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "| R01 | re-run: "+reason+" |") {
		t.Fatal(out.String())
	}
}

func TestRowCarryAcceptsAHeadItsQAReportCommitRecorded(t *testing.T) {
	root, head := rowCarrySidePass(t, "Roundfix-Spec: mechanical")
	result := rowCarryResult(t, root)
	if len(result.Carried) != 1 || result.Carried[0].EstablishedHead != head || result.Carried[0].EstablishedBy != rowCarryReportPath || len(result.Dispositions) != 1 || !result.Dispositions[0].Carried || result.Dispositions[0].Reason != "" {
		t.Fatalf("unexpected carry result: %#v", result)
	}
}
func TestRowCarryRefusesAHeadOnlyATaskCommitRecorded(t *testing.T) {
	root, _ := rowCarrySidePass(t, "Roundfix-Spec: mechanical\nRoundfix-Task: task_01")
	rowCarryRefusal(t, rowCarryResult(t, root), speccheck.CarryReasonEstablishingHeadUnproven)
}
func TestRowCarryRefusesAnUnprovenHead(t *testing.T) {
	root, _ := rowCarrySidePass(t, "")
	rowCarryRefusal(t, rowCarryResult(t, root), speccheck.CarryReasonEstablishingHeadUnproven)
}
func TestRowCarryNamesTheMovedInput(t *testing.T) {
	root, _ := rowCarrySidePass(t, "Roundfix-Spec: mechanical")
	writeMechanicalFile(t, root, "evidence.txt", "moved\n")
	commitMechanicalFiles(t, root, "move evidence", "evidence.txt")
	rowCarryRefusal(t, rowCarryResult(t, root), speccheck.CarryReasonInputMoved+"evidence.txt")
}
func TestRowCarryNeverCarriesAlwaysObservedRows(t *testing.T) {
	for _, tc := range []struct{ provenance, kind, reason string }{
		{"repository Verification", "repository_path", speccheck.CarryReasonRepositoryVerification},
		{"Pull Request row", "repository_path", speccheck.CarryReasonPullRequestRow},
		{"criterion", "commit_range", speccheck.CarryReasonCommitRangeInput},
	} {
		t.Run(tc.reason, func(t *testing.T) {
			root := newMechanicalGitRepo(t)
			writeMechanicalFile(t, root, "evidence.txt", "stable evidence\n")
			head := commitMechanicalFiles(t, root, "establish", "evidence.txt")
			report := mechanicalCarryReport(head, "stable evidence\n")
			report = strings.ReplaceAll(report, "| Evidence |", "| Evidence | Provenance |")
			report = strings.ReplaceAll(report, "| - | --- | --- | --- | --- |", "| - | --- | --- | --- | --- | --- |")
			report = strings.ReplaceAll(report, "[evidence](../../../../evidence.txt) |", "[evidence](../../../../evidence.txt) | "+tc.provenance+" |")
			report = strings.ReplaceAll(report, "kind: repository_path", "kind: "+tc.kind)
			writeMechanicalFile(t, root, rowCarryReportPath, report)
			commitMechanicalFiles(t, root, "record", rowCarryReportPath)
			rowCarryRefusal(t, rowCarryResult(t, root), tc.reason)
		})
	}
}
func TestRowCarryKeepsTheEstablishingProvenance(t *testing.T) {
	root, head := rowCarrySidePass(t, "Roundfix-Spec: mechanical")
	report := mechanicalCarryReport(head, "stable evidence\n")
	report = strings.ReplaceAll(report, "| Evidence |", "| Evidence | Provenance |")
	report = strings.ReplaceAll(report, "| - | --- | --- | --- | --- |", "| - | --- | --- | --- | --- | --- |")
	report = strings.ReplaceAll(report, "[evidence](../../../../evidence.txt) |", "[evidence](../../../../evidence.txt) | criterion 7 |")
	// An ancestor head remains independently proven even after this report edit.
	runMechanicalGit(t, root, "checkout", "-q", "qa-side")
	writeMechanicalFile(t, root, rowCarryReportPath, report)
	commitMechanicalFiles(t, root, "provenance", rowCarryReportPath)
	newer := "docs/specs/mechanical/qa/qa-report-2026-08-12.md"
	writeMechanicalFile(t, root, newer, mechanicalCarriedReport(rowCarryReportPath, head))
	result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: root, ReportPath: newer})
	if len(result.Carried) != 1 || result.Carried[0].Provenance != "criterion 7" {
		t.Fatalf("carried: %#v", result.Carried)
	}
	var out bytes.Buffer
	if err := speccheck.WriteMechanicalResult(&out, result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "; head: "+head+") | criterion 7 |") {
		t.Fatal(out.String())
	}
}
func TestRowCarryRecordsOneDispositionPerPriorRow(t *testing.T) {
	root, _ := rowCarrySidePass(t, "Roundfix-Spec: mechanical")
	path := root + "/" + rowCarryReportPath
	report := readRowCarryFile(t, path)
	report = strings.ReplaceAll(report, "|\n\n### R01 evidence", "|\n| R02 | failed check | backend | fail | observed |\n| R03 | undeclared check | backend | pass | observed |\n\n### R01 evidence")
	writeMechanicalFile(t, root, rowCarryReportPath, report)
	// Keep ancestry valid: the independent report was altered for extra rows.
	runMechanicalGit(t, root, "merge", "--no-edit", "qa-side")
	result := rowCarryResult(t, root)
	if len(result.Dispositions) != 3 {
		t.Fatalf("dispositions: %#v", result.Dispositions)
	}
	for i, want := range []string{"R01", "R02", "R03"} {
		if result.Dispositions[i].ID != want {
			t.Fatal(result.Dispositions)
		}
	}
	if !result.Dispositions[0].Carried || result.Dispositions[1].Reason != speccheck.CarryReasonNotPass || result.Dispositions[2].Reason != speccheck.CarryReasonNoInputs {
		t.Fatal(result.Dispositions)
	}
}
func TestRowCarryWithoutDispositionsRendersTodaysBytes(t *testing.T) {
	_ = newMechanicalGitRepo(t)
	var out bytes.Buffer
	result := speccheck.MechanicalResult{Carried: []speccheck.CarriedRow{{ID: "R01", EstablishedBy: "report.md", EstablishedHead: "abc"}}}
	if err := speccheck.WriteMechanicalResult(&out, result); err != nil {
		t.Fatal(err)
	}
	want := "## Performed repairs\n\nNone.\n\n## Assigned repair failures\n\nNone.\n\n## Authorization audit inputs\n\nNone.\n\n## Mechanical findings\n\nNone.\n\n## Mechanical rows\n\n| # | Status | Provenance |\n| - | --- | --- |\n| R01 | carried (established by: report.md; head: abc) | report and head retained |\n\n## Mechanical skips\n\nNone.\n"
	if out.String() != want {
		t.Fatalf("got %q, want %q", out.String(), want)
	}
}

func readRowCarryFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func TestRowCarryRefusesARecordingCommitForAnotherSpec(t *testing.T) {
	root, _ := rowCarrySidePass(t, "Roundfix-Spec: another-spec")
	rowCarryRefusal(t, rowCarryResult(t, root), speccheck.CarryReasonEstablishingHeadUnproven)
}
func TestRowCarryRefusesDifferentRecordedBytes(t *testing.T) {
	root, _ := rowCarrySidePass(t, "Roundfix-Spec: mechanical")
	report := readRowCarryFile(t, root+"/"+rowCarryReportPath)
	writeMechanicalFile(t, root, rowCarryReportPath, report+"\nChanged report bytes.\n")
	rowCarryRefusal(t, rowCarryResult(t, root), speccheck.CarryReasonEstablishingHeadUnproven)
}
func TestRowCarryRefusesARecordingCommitWithADifferentFirstParent(t *testing.T) {
	root, _ := rowCarrySidePass(t, "")
	report := readRowCarryFile(t, root+"/"+rowCarryReportPath)
	runMechanicalGit(t, root, "checkout", "-q", "qa-side")
	writeMechanicalFile(t, root, "other.txt", "other\n")
	commitMechanicalFiles(t, root, "intervening commit", "other.txt")
	writeMechanicalFile(t, root, rowCarryReportPath, report)
	commitMechanicalFiles(t, root, "docs: record QA\n\nRoundfix-Spec: mechanical", rowCarryReportPath)
	runMechanicalGit(t, root, "checkout", "-q", "main")
	rowCarryRefusal(t, rowCarryResult(t, root), speccheck.CarryReasonEstablishingHeadUnproven)
}
func TestRowCarryRefusesAnUnreachableRecordingCommit(t *testing.T) {
	root, _ := rowCarrySidePass(t, "Roundfix-Spec: mechanical")
	runMechanicalGit(t, root, "branch", "-D", "qa-side")
	rowCarryRefusal(t, rowCarryResult(t, root), speccheck.CarryReasonEstablishingHeadUnproven)
}
func TestRowCarryRefusesANonRepositoryInput(t *testing.T) {
	root, _ := rowCarrySidePass(t, "Roundfix-Spec: mechanical")
	report := readRowCarryFile(t, root+"/"+rowCarryReportPath)
	writeMechanicalFile(t, root, rowCarryReportPath, strings.ReplaceAll(report, "kind: repository_path", "kind: live_service"))
	rowCarryRefusal(t, rowCarryResult(t, root), speccheck.CarryReasonNonRepositoryInput)
}
func TestRowCarryRefusesAnUnavailableEstablishingReport(t *testing.T) {
	root, head := rowCarrySidePass(t, "Roundfix-Spec: mechanical")
	writeMechanicalFile(t, root, rowCarryReportPath, mechanicalCarriedReport("docs/specs/mechanical/qa/missing.md", head))
	rowCarryRefusal(t, rowCarryResult(t, root), speccheck.CarryReasonEstablishingReportUnavailable)
}
func TestRowCarryRefusesAMissingSnapshot(t *testing.T) {
	root, head := rowCarrySidePass(t, "Roundfix-Spec: mechanical")
	report := mechanicalCarryReport(head, "stable evidence\n")
	start := strings.Index(report, "evidence_snapshots:")
	end := strings.Index(report[start:], "---\n") + start
	writeMechanicalFile(t, root, rowCarryReportPath, report[:start]+report[end:])
	rowCarryRefusal(t, rowCarryResult(t, root), speccheck.CarryReasonNoEvidenceSnapshot)
}
func TestRowCarryNamesADeletedInput(t *testing.T) {
	root, _ := rowCarrySidePass(t, "Roundfix-Spec: mechanical")
	runMechanicalGit(t, root, "rm", "evidence.txt")
	commitMechanicalFiles(t, root, "delete input")
	rowCarryRefusal(t, rowCarryResult(t, root), speccheck.CarryReasonInputMoved+"evidence.txt")
}
