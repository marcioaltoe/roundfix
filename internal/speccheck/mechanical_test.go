// Suite: pre-QA mechanical facts
// Invariant: written QA declarations are compared with repository facts, and absent inputs are recorded as skips.
// Boundary IN: public speccheck mechanical API, real temporary Git histories, and report carrier fixtures
// Boundary OUT: Daemon scheduling, QA verdict computation, and carry-forward eligibility
package speccheck_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

func TestMechanicalAuthPaths(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	copyMechanicalFixture(t, repoRoot, "authorization.md", "docs/workflow/authorizations/mechanical.md")
	writeMechanicalFile(t, repoRoot, "Makefile", "verify:\n\t@true\n")
	greenCommit := commitMechanicalFiles(t, repoRoot, "authorized change", "Makefile")
	writeMechanicalFile(t, repoRoot, "outside.txt", "outside\n")
	redCommit := commitMechanicalFiles(t, repoRoot, "unauthorized change", "outside.txt")

	t.Run("green fixture accepts exact bounded path", func(t *testing.T) {
		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:          repoRoot,
			AuthorizationPath: "docs/workflow/authorizations/mechanical.md",
			TaskCommits: []speccheck.MechanicalTaskCommit{{
				TaskID:   "task_01",
				SHA:      greenCommit,
				TaskFile: "docs/specs/mechanical/task_01.md",
			}},
		})
		assertNoMechanicalCode(t, result, speccheck.CodeMechanicalAuthPaths)
	})

	t.Run("red fixture names every path outside the bound", func(t *testing.T) {
		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:          repoRoot,
			AuthorizationPath: "docs/workflow/authorizations/mechanical.md",
			TaskCommits: []speccheck.MechanicalTaskCommit{{
				TaskID:   "task_02",
				SHA:      redCommit,
				TaskFile: "docs/specs/mechanical/task_02.md",
			}},
		})
		findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalAuthPaths)
		if len(findings) != 1 || !strings.Contains(findings[0].Detail, "outside.txt") {
			t.Fatalf("%s findings = %#v, want outside.txt", speccheck.CodeMechanicalAuthPaths, findings)
		}
	})

	t.Run("absent authorization records a skip", func(t *testing.T) {
		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:          repoRoot,
			AuthorizationPath: "docs/workflow/authorizations/missing.md",
			TaskCommits:       []speccheck.MechanicalTaskCommit{{TaskID: "task_01", SHA: greenCommit}},
		})
		assertMechanicalSkip(t, result, speccheck.DetectorMechanicalAuthPaths, "docs/workflow/authorizations/missing.md")
		assertNoMechanicalCode(t, result, speccheck.CodeMechanicalAuthPaths)
	})
}

func TestMechanicalConsequentOrder(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	writeMechanicalFile(t, repoRoot, "cause.txt", "cause\n")
	cause := commitMechanicalFiles(t, repoRoot, "cause", "cause.txt")
	writeMechanicalFile(t, repoRoot, "fix.txt", "fix\n")
	fix := commitMechanicalFiles(t, repoRoot, "fix", "fix.txt")

	tests := []struct {
		name        string
		declaration speccheck.ConsequentFixDeclaration
		wantFinding bool
	}{
		{
			name: "green fixture keeps fix after cause",
			declaration: speccheck.ConsequentFixDeclaration{
				File: "docs/specs/mechanical/_tasks.md", Line: 12, RowHint: "R-COMMIT", CauseCommit: cause, FixCommit: fix,
			},
		},
		{
			name: "red fixture rejects fix before cause",
			declaration: speccheck.ConsequentFixDeclaration{
				File: "docs/specs/mechanical/_tasks.md", Line: 12, RowHint: "R-COMMIT", CauseCommit: fix, FixCommit: cause,
			},
			wantFinding: true,
		},
		{
			name: "red fixture rejects folded fix",
			declaration: speccheck.ConsequentFixDeclaration{
				File: "docs/specs/mechanical/_tasks.md", Line: 12, RowHint: "R-COMMIT", CauseCommit: cause, FixCommit: cause,
			},
			wantFinding: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := runMechanical(t, speccheck.MechanicalRequest{
				RepoRoot:        repoRoot,
				ConsequentFixes: []speccheck.ConsequentFixDeclaration{tt.declaration},
			})
			findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalConsequentOrder)
			if got := len(findings) > 0; got != tt.wantFinding {
				t.Fatalf("%s finding present = %t, want %t: %#v", speccheck.CodeMechanicalConsequentOrder, got, tt.wantFinding, findings)
			}
		})
	}

	t.Run("absent declaration records a skip", func(t *testing.T) {
		result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot})
		assertMechanicalSkip(t, result, speccheck.DetectorMechanicalConsequentOrder, "consequent-fix declaration")
	})
}

func TestMechanicalReportShape(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	copyMechanicalFixture(t, repoRoot, "report-green.md", "docs/specs/mechanical/qa/report-green.md")
	copyMechanicalFixture(t, repoRoot, "report-red.md", "docs/specs/mechanical/qa/report-red.md")
	copyMechanicalFixture(t, repoRoot, "evidence/pass.txt", "docs/specs/mechanical/qa/evidence/pass.txt")

	t.Run("green fixture has terminal typed rows and exact counts", func(t *testing.T) {
		result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: "docs/specs/mechanical/qa/report-green.md"})
		assertNoMechanicalCode(t, result, speccheck.CodeMechanicalReportShape)
	})

	t.Run("red fixture reports every structural defect", func(t *testing.T) {
		result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: "docs/specs/mechanical/qa/report-red.md"})
		findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalReportShape)
		if len(findings) < 3 {
			t.Fatalf("%s findings = %#v, want pending, untyped block, and count mismatch", speccheck.CodeMechanicalReportShape, findings)
		}
		for _, detail := range []string{"pending", "blocked cause", "rows_blocked_finding"} {
			if !mechanicalDetailsContain(findings, detail) {
				t.Errorf("%s findings = %#v, want detail containing %q", speccheck.CodeMechanicalReportShape, findings, detail)
			}
		}
	})

	t.Run("absent report records a skip", func(t *testing.T) {
		result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: "docs/specs/mechanical/qa/missing.md"})
		assertMechanicalSkip(t, result, speccheck.DetectorMechanicalReportShape, "docs/specs/mechanical/qa/missing.md")
		assertNoMechanicalCode(t, result, speccheck.CodeMechanicalReportShape)
	})
}

func TestMechanicalEvidencePath(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	copyMechanicalFixture(t, repoRoot, "report-green.md", "docs/specs/mechanical/qa/report-green.md")
	copyMechanicalFixture(t, repoRoot, "report-red.md", "docs/specs/mechanical/qa/report-red.md")
	copyMechanicalFixture(t, repoRoot, "report-raw.md", "docs/specs/mechanical/qa/report-raw.md")
	copyMechanicalFixture(t, repoRoot, "evidence/pass.txt", "docs/specs/mechanical/qa/evidence/pass.txt")

	green := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: "docs/specs/mechanical/qa/report-green.md"})
	assertNoMechanicalCode(t, green, speccheck.CodeMechanicalEvidencePath)
	raw := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: "docs/specs/mechanical/qa/report-raw.md"})
	assertNoMechanicalCode(t, raw, speccheck.CodeMechanicalEvidencePath)
	for _, skip := range raw.Skips {
		if skip.Detector == speccheck.DetectorMechanicalEvidencePath {
			t.Fatalf("raw evidence carrier was skipped instead of resolved: %#v", raw.Skips)
		}
	}

	red := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: "docs/specs/mechanical/qa/report-red.md"})
	findings := mechanicalFindingsWithCode(red, speccheck.CodeMechanicalEvidencePath)
	if len(findings) != 1 || !strings.Contains(findings[0].Detail, "evidence/missing.txt") {
		t.Fatalf("%s findings = %#v, want missing evidence path", speccheck.CodeMechanicalEvidencePath, findings)
	}

	missing := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: "docs/specs/mechanical/qa/missing.md"})
	assertMechanicalSkip(t, missing, speccheck.DetectorMechanicalEvidencePath, "docs/specs/mechanical/qa/missing.md")
}

func TestMechanicalReportsAllFindings(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	copyMechanicalFixture(t, repoRoot, "authorization.md", "docs/workflow/authorizations/mechanical.md")
	copyMechanicalFixture(t, repoRoot, "report-red.md", "docs/specs/mechanical/qa/report-red.md")
	copyMechanicalFixture(t, repoRoot, "evidence/pass.txt", "docs/specs/mechanical/qa/evidence/pass.txt")
	writeMechanicalFile(t, repoRoot, "outside.txt", "outside\n")
	commit := commitMechanicalFiles(t, repoRoot, "outside", "outside.txt")

	result := runMechanical(t, speccheck.MechanicalRequest{
		RepoRoot:          repoRoot,
		AuthorizationPath: "docs/workflow/authorizations/mechanical.md",
		TaskCommits:       []speccheck.MechanicalTaskCommit{{TaskID: "task_01", SHA: commit}},
		ConsequentFixes: []speccheck.ConsequentFixDeclaration{{
			File: "docs/specs/mechanical/_tasks.md", Line: 12, RowHint: "R-COMMIT", CauseCommit: commit, FixCommit: commit,
		}},
		ReportPath: "docs/specs/mechanical/qa/report-red.md",
	})
	for _, code := range []string{
		speccheck.CodeMechanicalAuthPaths,
		speccheck.CodeMechanicalConsequentOrder,
		speccheck.CodeMechanicalReportShape,
		speccheck.CodeMechanicalEvidencePath,
	} {
		if len(mechanicalFindingsWithCode(result, code)) == 0 {
			t.Errorf("Findings = %#v, want %s", result.Findings, code)
		}
	}
	if !result.Blocking {
		t.Fatal("Blocking = false, want true when findings exist")
	}
	for _, finding := range result.Findings {
		if finding.Code == "" || finding.File == "" || finding.Line < 1 || finding.Detail == "" || finding.Fix == "" {
			t.Errorf("finding omits a required typed field: %#v", finding)
		}
	}
	for _, row := range []string{"task_01", "R-COMMIT", "R01", "R02"} {
		if !mechanicalResultBlocksRow(result, row) {
			t.Errorf("Blocked = %#v, want row %s tied to its finding", result.Blocked, row)
		}
	}
}

func TestMaterializeMechanicalResult(t *testing.T) {
	t.Parallel()

	result := speccheck.MechanicalResult{
		Findings: []speccheck.MechanicalFinding{{
			Code: "QA-FIXTURE", File: "fixture.md", Line: 7, Detail: "fixture mismatch", Fix: "repair fixture", RowHint: "R02",
		}},
		Carried:  []speccheck.CarriedRow{{ID: "R01", EstablishedBy: "qa/previous.md", EstablishedHead: "abc123"}},
		Blocked:  []speccheck.BlockedRow{{ID: "R02", FindingCode: "QA-FIXTURE", WaitingOn: "fixture mismatch"}},
		Skips:    []speccheck.MechanicalSkip{{Detector: "fixture-detector", MissingArtifact: "fixture.md"}},
		Blocking: true,
	}

	var report bytes.Buffer
	if err := speccheck.WriteMechanicalResult(&report, result); err != nil {
		t.Fatalf("WriteMechanicalResult() error = %v", err)
	}
	for _, fragment := range []string{
		"### QA-FIXTURE",
		"fixture.md:7",
		"R01",
		"carried (established by: qa/previous.md; head: abc123)",
		"blocked (finding: QA-FIXTURE — waits on fixture mismatch)",
		"fixture-detector",
		"fixture.md",
	} {
		if !strings.Contains(report.String(), fragment) {
			t.Errorf("materialized report does not contain %q:\n%s", fragment, report.String())
		}
	}
	if strings.Contains(report.String(), "verdict:") {
		t.Fatalf("materializer computed a verdict:\n%s", report.String())
	}
}

func TestMechanicalCorpusNonRegression(t *testing.T) {
	t.Parallel()
	assertMechanicalCorpusBudget(t)
}

func TestCheckCorpusBudget(t *testing.T) {
	t.Parallel()
	assertMechanicalCorpusBudget(t)
}

func assertMechanicalCorpusBudget(t *testing.T) {
	t.Helper()

	entries, err := os.ReadDir(fixtureSpecRoot)
	if err != nil {
		t.Fatalf("read fixture corpus: %v", err)
	}
	repoRoot, err := filepath.Abs("testdata/repo")
	if err != nil {
		t.Fatalf("resolve fixture repository: %v", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), "_") {
			continue
		}
		result, checkErr := speccheck.Check(fixtureSpecRoot, repoRoot, entry.Name())
		if checkErr != nil {
			t.Fatalf("Check(%q) error = %v", entry.Name(), checkErr)
		}
		for _, finding := range result.Findings {
			if strings.HasPrefix(finding.Code, "QA-") {
				t.Errorf("pre-existing fixture %s gained mechanical diagnostic %#v", entry.Name(), finding)
			}
		}
	}
}

func runMechanical(t *testing.T, request speccheck.MechanicalRequest) speccheck.MechanicalResult {
	t.Helper()
	result, err := speccheck.RunMechanicalStage(context.Background(), request)
	if err != nil {
		t.Fatalf("RunMechanicalStage() error = %v", err)
	}
	return result
}

func newMechanicalGitRepo(t *testing.T) string {
	t.Helper()
	repoRoot := t.TempDir()
	runMechanicalGit(t, repoRoot, "init", "--quiet")
	runMechanicalGit(t, repoRoot, "config", "user.name", "Roundfix Test")
	runMechanicalGit(t, repoRoot, "config", "user.email", "roundfix@example.invalid")
	writeMechanicalFile(t, repoRoot, ".keep", "fixture\n")
	commitMechanicalFiles(t, repoRoot, "initial", ".keep")
	return repoRoot
}

func commitMechanicalFiles(t *testing.T, repoRoot, message string, paths ...string) string {
	t.Helper()
	args := append([]string{"add", "--"}, paths...)
	runMechanicalGit(t, repoRoot, args...)
	runMechanicalGit(t, repoRoot, "commit", "--quiet", "-m", message)
	return strings.TrimSpace(runMechanicalGit(t, repoRoot, "rev-parse", "HEAD"))
}

func runMechanicalGit(t *testing.T, repoRoot string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repoRoot}, args...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, output)
	}
	return string(output)
}

func copyMechanicalFixture(t *testing.T, repoRoot, source, destination string) {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("testdata", "mechanical", filepath.FromSlash(source)))
	if err != nil {
		t.Fatalf("read mechanical fixture %q: %v", source, err)
	}
	writeMechanicalFile(t, repoRoot, destination, string(content))
}

func writeMechanicalFile(t *testing.T, repoRoot, relative, content string) {
	t.Helper()
	path := filepath.Join(repoRoot, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture directory for %q: %v", relative, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture %q: %v", relative, err)
	}
}

func mechanicalFindingsWithCode(result speccheck.MechanicalResult, code string) []speccheck.MechanicalFinding {
	var findings []speccheck.MechanicalFinding
	for _, finding := range result.Findings {
		if finding.Code == code {
			findings = append(findings, finding)
		}
	}
	return findings
}

func assertNoMechanicalCode(t *testing.T, result speccheck.MechanicalResult, code string) {
	t.Helper()
	if findings := mechanicalFindingsWithCode(result, code); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want none", code, findings)
	}
}

func assertMechanicalSkip(t *testing.T, result speccheck.MechanicalResult, detector, missing string) {
	t.Helper()
	want := speccheck.MechanicalSkip{Detector: detector, MissingArtifact: missing}
	for _, skip := range result.Skips {
		if reflect.DeepEqual(skip, want) {
			return
		}
	}
	t.Fatalf("Skips = %#v, want %#v", result.Skips, want)
}

func mechanicalDetailsContain(findings []speccheck.MechanicalFinding, fragment string) bool {
	for _, finding := range findings {
		if strings.Contains(finding.Detail, fragment) {
			return true
		}
	}
	return false
}

func mechanicalResultBlocksRow(result speccheck.MechanicalResult, row string) bool {
	for _, blocked := range result.Blocked {
		if blocked.ID == row && blocked.FindingCode != "" && blocked.WaitingOn != "" {
			return true
		}
	}
	return false
}
