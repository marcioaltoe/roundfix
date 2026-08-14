// Suite: pre-QA mechanical facts
// Invariant: written QA declarations are compared with repository facts, and only unchanged repository evidence carries.
// Boundary IN: public speccheck mechanical API, real temporary Git histories, and report carrier fixtures
// Boundary OUT: Daemon scheduling and QA verdict computation
package speccheck_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:          repoRoot,
			AuthorizationPath: "docs/workflow/authorizations/missing.md",
			TaskCommits:       []speccheck.MechanicalTaskCommit{{TaskID: "task_01", SHA: greenCommit}},
		})
		assertMechanicalSkip(t, result, speccheck.DetectorMechanicalAuthPaths, "docs/workflow/authorizations/missing.md")
		assertNoMechanicalCode(t, result, speccheck.CodeMechanicalAuthPaths)
	})
}

func TestMechanicalAuthPathsAcceptsDeclaredRegenerationOutput(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	writeMechanicalFile(t, repoRoot, "docs/workflow/authorizations/mechanical.md", mechanicalRegenerationAuthorization("command: make baseline-digests\noutputs:\n  - internal/baseline/testdata/catalog.digest\n"))
	writeMechanicalFile(t, repoRoot, "internal/baseline/testdata/catalog.digest", "generated\n")
	commit := commitMechanicalFiles(t, repoRoot, "regenerate baseline digest", "internal/baseline/testdata/catalog.digest")

	result := runMechanical(t, speccheck.MechanicalRequest{
		RepoRoot:          repoRoot,
		AuthorizationPath: "docs/workflow/authorizations/mechanical.md",
		TaskCommits: []speccheck.MechanicalTaskCommit{{
			TaskID:   "task_06",
			SHA:      commit,
			TaskFile: "docs/specs/mechanical/task_06.md",
		}},
	})

	assertNoMechanicalCode(t, result, speccheck.CodeMechanicalAuthPaths)
}

func TestMechanicalAuthPathsStillRefusesAnUndeclaredPath(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	writeMechanicalFile(t, repoRoot, "docs/workflow/authorizations/mechanical.md", mechanicalRegenerationAuthorization("command: make baseline-digests\noutputs:\n  - internal/baseline/testdata/catalog.digest\n"))
	writeMechanicalFile(t, repoRoot, "outside.txt", "outside\n")
	commit := commitMechanicalFiles(t, repoRoot, "change undeclared path", "outside.txt")

	result := runMechanical(t, speccheck.MechanicalRequest{
		RepoRoot:          repoRoot,
		AuthorizationPath: "docs/workflow/authorizations/mechanical.md",
		TaskCommits: []speccheck.MechanicalTaskCommit{{
			TaskID:   "task_06",
			SHA:      commit,
			TaskFile: "docs/specs/mechanical/task_06.md",
		}},
	})

	findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalAuthPaths)
	if len(findings) != 1 || !strings.Contains(findings[0].Detail, "outside.txt") {
		t.Fatalf("%s findings = %#v, want outside.txt", speccheck.CodeMechanicalAuthPaths, findings)
	}
}

func TestMechanicalAuthPathsRefusesInvalidRegenerationDeclaration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		declaration string
	}{
		{name: "output without command", declaration: "outputs:\n  - internal/baseline/testdata/catalog.digest\n"},
		{name: "output glob", declaration: "command: make baseline-digests\noutputs:\n  - internal/baseline/testdata/*.digest\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repoRoot := newMechanicalGitRepo(t)
			writeMechanicalFile(t, repoRoot, "docs/workflow/authorizations/mechanical.md", mechanicalRegenerationAuthorization(tt.declaration))
			writeMechanicalFile(t, repoRoot, "internal/baseline/testdata/catalog.digest", "hand edited\n")
			commit := commitMechanicalFiles(t, repoRoot, "hand edit baseline digest", "internal/baseline/testdata/catalog.digest")

			result := runMechanical(t, speccheck.MechanicalRequest{
				RepoRoot:          repoRoot,
				AuthorizationPath: "docs/workflow/authorizations/mechanical.md",
				TaskCommits: []speccheck.MechanicalTaskCommit{{
					TaskID:   "task_06",
					SHA:      commit,
					TaskFile: "docs/specs/mechanical/task_06.md",
				}},
			})

			findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalAuthPaths)
			if len(findings) != 1 || !strings.Contains(findings[0].Detail, "internal/baseline/testdata/catalog.digest") {
				t.Fatalf("%s findings = %#v, want invalid regeneration declaration refused", speccheck.CodeMechanicalAuthPaths, findings)
			}
		})
	}
}

func TestMechanicalAuthorizationReadsThePRDBoundedDeclaration(t *testing.T) {
	t.Parallel()
	repoRoot := t.TempDir()
	writeMechanicalFile(t, repoRoot, "docs/workflow/authorizations/mechanical.md", "# Authorization\n\n## Bounded files\n\n- `Makefile`\n- `docs/agents/spec-routing.md`\n")
	prdPath := filepath.Join(repoRoot, "docs", "specs", "mechanical", "_prd.md")
	writeMechanicalFile(t, repoRoot, "docs/specs/mechanical/_prd.md", "# PRD\n\n## Project Constraints\n\n- Tooling authority: applicable — recorded at `docs/workflow/authorizations/mechanical.md`; bounded files: `Makefile`. Source: `docs/agents/agent-instructions.md`.\n")

	path, bounded, err := speccheck.MechanicalAuthorization(repoRoot, prdPath)

	if err != nil {
		t.Fatalf("MechanicalAuthorization returned error: %v", err)
	}
	if path != "docs/workflow/authorizations/mechanical.md" {
		t.Fatalf("MechanicalAuthorization path = %q", path)
	}
	want := []string{"Makefile", "docs/agents/spec-routing.md"}
	if !reflect.DeepEqual(bounded, want) {
		t.Fatalf("MechanicalAuthorization bounded paths = %v, want %v", bounded, want)
	}

	missingPath, missingBounded, err := speccheck.MechanicalAuthorization(repoRoot, filepath.Join(repoRoot, "docs", "specs", "absent", "_prd.md"))
	if err != nil {
		t.Fatalf("MechanicalAuthorization(absent) returned error: %v", err)
	}
	if missingPath != "" || missingBounded != nil {
		t.Fatalf("MechanicalAuthorization(absent) = %q, %v; want empty presence-aware input", missingPath, missingBounded)
	}
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
		result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: "docs/specs/mechanical/qa/report-green.md"})
		assertNoMechanicalCode(t, result, speccheck.CodeMechanicalReportShape)
	})

	t.Run("red fixture reports every structural defect", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
		result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: "docs/specs/mechanical/qa/missing.md"})
		assertMechanicalSkip(t, result, speccheck.DetectorMechanicalReportShape, "docs/specs/mechanical/qa/missing.md")
		assertNoMechanicalCode(t, result, speccheck.CodeMechanicalReportShape)
	})
}

func TestMechanicalFindingsWithoutRowHintsBlockTheirRefusalCode(t *testing.T) {
	t.Parallel()
	repoRoot := newMechanicalGitRepo(t)
	reportPath := "docs/specs/mechanical/qa/report-without-results.md"
	writeMechanicalFile(t, repoRoot, reportPath, "---\nverdict: fail\nrows_blocked_environment: 0\nrows_blocked_finding: 1\nrows_blocked_declared: 0\n---\n\n# QA Report\n\n## Mechanical rows\n\n| # | Status | Provenance |\n| - | --- | --- |\n| R01 | blocked (finding: QA-FIXTURE — waits on fixture) | mechanical finding |\n")

	result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: reportPath})

	findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalReportShape)
	if len(findings) != 2 {
		t.Fatalf("%s findings = %#v, want missing Results rows and count mismatch", speccheck.CodeMechanicalReportShape, findings)
	}
	if len(result.Blocked) != 1 {
		t.Fatalf("Blocked = %#v, want one row for the unscoped refusal code", result.Blocked)
	}
	blocked := result.Blocked[0]
	if blocked.ID != speccheck.CodeMechanicalReportShape || blocked.FindingCode != speccheck.CodeMechanicalReportShape {
		t.Fatalf("Blocked[0] = %#v, want the existing refusal code as row identity and provenance", blocked)
	}
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

const (
	evidenceScratchReportPath = "docs/specs/mechanical/qa/qa-report.md"
	evidenceScratchRoot       = "docs/specs/mechanical/qa/evidence"
)

func TestEvidenceRefusesScratchStateBinary(t *testing.T) {
	t.Parallel()
	repoRoot := newMechanicalGitRepo(t)
	binaryPath := evidenceScratchRoot + "/roundfix"
	writeMechanicalFile(t, repoRoot, evidenceScratchReportPath, mechanicalEvidenceReport())
	writeMechanicalFile(t, repoRoot, binaryPath, "\x7fELF\x00fixture")
	commitMechanicalFiles(t, repoRoot, "record binary evidence", evidenceScratchReportPath, binaryPath)

	result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: evidenceScratchReportPath})

	findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalEvidencePath)
	if len(findings) != 1 || !strings.Contains(findings[0].Detail, binaryPath) ||
		!strings.Contains(findings[0].Detail, "binary") {
		t.Fatalf("%s findings = %#v, want named binary %s", speccheck.CodeMechanicalEvidencePath, findings, binaryPath)
	}
}

func TestEvidenceRefusesScratchStateGitlink(t *testing.T) {
	t.Parallel()
	repoRoot := newMechanicalGitRepo(t)
	gitlinkPath := evidenceScratchRoot + "/scratch-repository"
	writeMechanicalFile(t, repoRoot, evidenceScratchReportPath, mechanicalEvidenceReport())
	commitMechanicalFiles(t, repoRoot, "record QA report", evidenceScratchReportPath)
	target := strings.TrimSpace(runMechanicalGit(t, repoRoot, "rev-parse", "HEAD"))
	runMechanicalGit(t, repoRoot, "update-index", "--add", "--cacheinfo", "160000,"+target+","+gitlinkPath)
	runMechanicalGit(t, repoRoot, "commit", "--quiet", "-m", "record gitlink evidence")

	result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: evidenceScratchReportPath})

	findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalEvidencePath)
	if len(findings) != 1 || !strings.Contains(findings[0].Detail, gitlinkPath) ||
		!strings.Contains(findings[0].Detail, "gitlink") {
		t.Fatalf("%s findings = %#v, want named gitlink %s", speccheck.CodeMechanicalEvidencePath, findings, gitlinkPath)
	}
}

func TestEvidenceRefusesScratchStateOrdinaryArtifact(t *testing.T) {
	t.Parallel()
	repoRoot := newMechanicalGitRepo(t)
	artifactPath := evidenceScratchRoot + "/result.json"
	writeMechanicalFile(t, repoRoot, evidenceScratchReportPath, mechanicalEvidenceReport())
	writeMechanicalFile(t, repoRoot, artifactPath, "{\"result\":\"pass\"}\n")
	commitMechanicalFiles(t, repoRoot, "record readable evidence", evidenceScratchReportPath, artifactPath)

	result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: evidenceScratchReportPath})

	assertNoMechanicalCode(t, result, speccheck.CodeMechanicalEvidencePath)
}

func TestEvidenceRefusesScratchStateAbsentDirectory(t *testing.T) {
	t.Parallel()
	repoRoot := newMechanicalGitRepo(t)
	writeMechanicalFile(t, repoRoot, evidenceScratchReportPath, mechanicalEvidenceReport())
	commitMechanicalFiles(t, repoRoot, "record QA report", evidenceScratchReportPath)

	result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: evidenceScratchReportPath})

	assertNoMechanicalCode(t, result, speccheck.CodeMechanicalEvidencePath)
	assertMechanicalSkip(t, result, speccheck.DetectorMechanicalEvidencePath, evidenceScratchRoot)
}

func TestCarriable(t *testing.T) {
	t.Parallel()

	const (
		establishedHead = "established-head"
		currentHead     = "current-head"
		reportPath      = "docs/specs/mechanical/qa/qa-report-2026-08-11.md"
	)
	valid := func() (speccheck.ReportRow, []speccheck.EvidenceSnapshot, []speccheck.EvidenceSnapshot) {
		inputs := []speccheck.EvidenceInput{{Kind: speccheck.EvidenceRepositoryPath, Ref: "internal/speccheck/report.go"}}
		established := []speccheck.EvidenceSnapshot{{
			Ref: inputs[0].Ref,
			Files: []speccheck.EvidenceFile{{
				Path:   inputs[0].Ref,
				SHA256: strings.Repeat("a", 64),
			}},
		}}
		current := []speccheck.EvidenceSnapshot{{
			Ref: inputs[0].Ref,
			Files: []speccheck.EvidenceFile{{
				Path:   inputs[0].Ref,
				SHA256: strings.Repeat("a", 64),
			}},
		}}
		return speccheck.ReportRow{
			ID:               "R01",
			Status:           "pass",
			EstablishedBy:    reportPath,
			EstablishedHead:  establishedHead,
			AncestryVerified: true,
			Inputs:           inputs,
			EvidencePaths:    []string{inputs[0].Ref},
		}, established, current
	}

	refusals := []struct {
		name    string
		arrange func(*speccheck.ReportRow, *[]string, *[]speccheck.EvidenceSnapshot, *[]speccheck.EvidenceSnapshot)
	}{
		{
			name: "prior status failed refuses carry",
			arrange: func(row *speccheck.ReportRow, _ *[]string, _, _ *[]speccheck.EvidenceSnapshot) {
				row.Status = "fail"
			},
		},
		{
			name: "prior status blocked refuses carry",
			arrange: func(row *speccheck.ReportRow, _ *[]string, _, _ *[]speccheck.EvidenceSnapshot) {
				row.Status = "blocked (environment: fixture)"
			},
		},
		{
			name: "prior status skipped refuses carry",
			arrange: func(row *speccheck.ReportRow, _ *[]string, _, _ *[]speccheck.EvidenceSnapshot) {
				row.Status = "skipped"
			},
		},
		{
			name: "no declared inputs refuses carry",
			arrange: func(row *speccheck.ReportRow, _ *[]string, _, _ *[]speccheck.EvidenceSnapshot) {
				row.Inputs = nil
			},
		},
		{
			name: "changed repository path refuses carry",
			arrange: func(_ *speccheck.ReportRow, changed *[]string, _, _ *[]speccheck.EvidenceSnapshot) {
				*changed = []string{"internal/speccheck/report.go"}
			},
		},
		{
			name: "establishing head not ancestor refuses carry",
			arrange: func(row *speccheck.ReportRow, _ *[]string, _, _ *[]speccheck.EvidenceSnapshot) {
				row.AncestryVerified = false
			},
		},
		{
			name: "changed evidence content refuses carry",
			arrange: func(_ *speccheck.ReportRow, _ *[]string, _, current *[]speccheck.EvidenceSnapshot) {
				(*current)[0].Files[0].SHA256 = strings.Repeat("b", 64)
			},
		},
		{
			name: "mixed repository and elapsed inputs refuse carry",
			arrange: func(row *speccheck.ReportRow, _ *[]string, established, current *[]speccheck.EvidenceSnapshot) {
				row.Inputs = append(row.Inputs, speccheck.EvidenceInput{Kind: speccheck.EvidenceElapsedTime, Ref: "thirty-day window"})
				*established = append(*established, speccheck.EvidenceSnapshot{Ref: "thirty-day window"})
				*current = append(*current, speccheck.EvidenceSnapshot{Ref: "thirty-day window"})
			},
		},
		{
			name: "missing established snapshot refuses carry",
			arrange: func(_ *speccheck.ReportRow, _ *[]string, established, _ *[]speccheck.EvidenceSnapshot) {
				*established = nil
			},
		},
		{
			name: "missing current snapshot refuses carry",
			arrange: func(_ *speccheck.ReportRow, _ *[]string, _, current *[]speccheck.EvidenceSnapshot) {
				*current = nil
			},
		},
		{
			name: "missing establishing report citation refuses carry",
			arrange: func(row *speccheck.ReportRow, _ *[]string, _, _ *[]speccheck.EvidenceSnapshot) {
				row.EstablishedBy = ""
			},
		},
		{
			name: "cited evidence outside declared inputs refuses carry",
			arrange: func(row *speccheck.ReportRow, _ *[]string, _, _ *[]speccheck.EvidenceSnapshot) {
				row.EvidencePaths = []string{"qa/evidence/R01.txt"}
			},
		},
		{
			name: "changed path entering recursive glob refuses carry",
			arrange: func(row *speccheck.ReportRow, changed *[]string, established, current *[]speccheck.EvidenceSnapshot) {
				row.Inputs[0].Ref = "internal/speccheck/**"
				(*established)[0].Ref = row.Inputs[0].Ref
				(*current)[0].Ref = row.Inputs[0].Ref
				*changed = []string{"internal/speccheck/mechanical.go"}
			},
		},
		{
			name: "expanded path set changed refuses carry",
			arrange: func(row *speccheck.ReportRow, _ *[]string, established, current *[]speccheck.EvidenceSnapshot) {
				row.Inputs[0].Ref = "internal/speccheck/**"
				(*established)[0].Ref = row.Inputs[0].Ref
				(*current)[0].Ref = row.Inputs[0].Ref
				(*current)[0].Files = append((*current)[0].Files, speccheck.EvidenceFile{
					Path:   "internal/speccheck/mechanical.go",
					SHA256: strings.Repeat("c", 64),
				})
			},
		},
	}

	for _, tt := range refusals {
		t.Run(tt.name, func(t *testing.T) {
			row, established, current := valid()
			var changed []string
			tt.arrange(&row, &changed, &established, &current)
			if speccheck.Carriable(row, currentHead, changed, established, current) {
				t.Fatal("Carriable() = true, want false")
			}
		})
	}

	t.Run("every condition holds and carries", func(t *testing.T) {
		row, established, current := valid()
		if !speccheck.Carriable(row, currentHead, nil, established, current) {
			t.Fatal("Carriable() = false, want true")
		}
	})

	t.Run("unchanged recursive glob carries", func(t *testing.T) {
		row, _, _ := valid()
		row.Inputs[0].Ref = "internal/speccheck/**"
		files := []speccheck.EvidenceFile{
			{Path: "internal/speccheck/mechanical.go", SHA256: strings.Repeat("a", 64)},
			{Path: "internal/speccheck/report.go", SHA256: strings.Repeat("b", 64)},
		}
		established := []speccheck.EvidenceSnapshot{{Ref: row.Inputs[0].Ref, Files: append([]speccheck.EvidenceFile(nil), files...)}}
		current := []speccheck.EvidenceSnapshot{{Ref: row.Inputs[0].Ref, Files: append([]speccheck.EvidenceFile(nil), files...)}}
		if !speccheck.Carriable(row, currentHead, []string{"docs/specs/task_01.md"}, established, current) {
			t.Fatal("Carriable() = false, want unchanged recursive glob to carry")
		}
	})
}

func TestMechanicalStageCarriableCarriesUnchangedEvidenceWithCitation(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	const evidenceContent = "stable evidence\n"
	writeMechanicalFile(t, repoRoot, "evidence.txt", evidenceContent)
	establishedHead := commitMechanicalFiles(t, repoRoot, "establish evidence", "evidence.txt")
	reportPath := "docs/specs/mechanical/qa/qa-report-2026-08-11.md"
	writeMechanicalFile(t, repoRoot, reportPath, mechanicalCarryReport(establishedHead, evidenceContent))
	commitMechanicalFiles(t, repoRoot, "record QA report", reportPath)
	writeMechanicalFile(t, repoRoot, "unrelated.txt", "unrelated\n")
	commitMechanicalFiles(t, repoRoot, "change unrelated input", "unrelated.txt")

	result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: reportPath})
	if len(result.Carried) != 1 {
		t.Fatalf("Carried = %#v, want one unchanged row", result.Carried)
	}
	carried := result.Carried[0]
	if carried.ID != "R01" || carried.EstablishedBy != reportPath || carried.EstablishedHead != establishedHead {
		t.Fatalf("Carried[0] = %#v, want report %q at %q", carried, reportPath, establishedHead)
	}
	var materialized bytes.Buffer
	if err := speccheck.WriteMechanicalResult(&materialized, result); err != nil {
		t.Fatalf("WriteMechanicalResult() error = %v", err)
	}
	wantCitation := "carried (established by: " + reportPath + "; head: " + establishedHead + ")"
	if !strings.Contains(materialized.String(), wantCitation) {
		t.Fatalf("materialized carried row omits citation %q:\n%s", wantCitation, materialized.String())
	}
}

func TestMechanicalStageCarriableRefusesNonAncestorEstablishingHead(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	mainBranch := strings.TrimSpace(runMechanicalGit(t, repoRoot, "branch", "--show-current"))
	runMechanicalGit(t, repoRoot, "checkout", "--quiet", "-b", "establishing-side")
	const evidenceContent = "stable evidence\n"
	writeMechanicalFile(t, repoRoot, "evidence.txt", evidenceContent)
	nonAncestorHead := commitMechanicalFiles(t, repoRoot, "side evidence", "evidence.txt")
	runMechanicalGit(t, repoRoot, "checkout", "--quiet", mainBranch)
	writeMechanicalFile(t, repoRoot, "evidence.txt", evidenceContent)
	reportPath := "docs/specs/mechanical/qa/qa-report-2026-08-11.md"
	writeMechanicalFile(t, repoRoot, reportPath, mechanicalCarryReport(nonAncestorHead, evidenceContent))
	commitMechanicalFiles(t, repoRoot, "record unrelated QA report", "evidence.txt", reportPath)

	result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: reportPath})
	if len(result.Carried) != 0 {
		t.Fatalf("Carried = %#v, want non-ancestor evidence re-observed", result.Carried)
	}
}

func TestMechanicalStageCarriableRefusesChangedDeclaredInput(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	const evidenceContent = "stable evidence\n"
	writeMechanicalFile(t, repoRoot, "evidence.txt", evidenceContent)
	establishedHead := commitMechanicalFiles(t, repoRoot, "establish evidence", "evidence.txt")
	reportPath := "docs/specs/mechanical/qa/qa-report-2026-08-11.md"
	writeMechanicalFile(t, repoRoot, reportPath, mechanicalCarryReport(establishedHead, evidenceContent))
	commitMechanicalFiles(t, repoRoot, "record QA report", reportPath)
	writeMechanicalFile(t, repoRoot, "evidence.txt", "changed evidence\n")
	commitMechanicalFiles(t, repoRoot, "change declared evidence", "evidence.txt")

	result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: reportPath})
	if len(result.Carried) != 0 {
		t.Fatalf("Carried = %#v, want changed declared input re-observed", result.Carried)
	}
}

func TestMechanicalStageCarriablePreservesOriginalEstablishingCitation(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	const evidenceContent = "stable evidence\n"
	writeMechanicalFile(t, repoRoot, "evidence.txt", evidenceContent)
	establishedHead := commitMechanicalFiles(t, repoRoot, "establish evidence", "evidence.txt")
	originalReport := "docs/specs/mechanical/qa/qa-report-2026-08-10.md"
	writeMechanicalFile(t, repoRoot, originalReport, mechanicalCarryReport(establishedHead, evidenceContent))
	commitMechanicalFiles(t, repoRoot, "record original QA report", originalReport)
	previousReport := "docs/specs/mechanical/qa/qa-report-2026-08-11.md"
	writeMechanicalFile(t, repoRoot, previousReport, mechanicalCarriedReport(originalReport, establishedHead))
	commitMechanicalFiles(t, repoRoot, "record carried QA report", previousReport)

	result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: previousReport})
	if len(result.Carried) != 1 || result.Carried[0].EstablishedBy != originalReport || result.Carried[0].EstablishedHead != establishedHead {
		t.Fatalf("Carried = %#v, want original report and head citation", result.Carried)
	}
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

func mechanicalRegenerationAuthorization(declaration string) string {
	return "# Tooling authorization\n\n## Bounded files\n\n- `Makefile`\n\n## Sanctioned regeneration\n\n```yaml\n" + declaration + "```\n"
}

func newMechanicalGitRepo(t *testing.T) string {
	t.Helper()
	repoRoot := t.TempDir()
	runMechanicalGit(t, repoRoot, "init", "--quiet")
	runMechanicalGit(t, repoRoot, "config", "user.name", "Roundfix Test")
	runMechanicalGit(t, repoRoot, "config", "user.email", "roundfix@example.invalid")
	runMechanicalGit(t, repoRoot, "config", "commit.gpgsign", "false")
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

func mechanicalCarryReport(establishedHead, evidenceContent string) string {
	digest := fmt.Sprintf("%x", sha256.Sum256([]byte(evidenceContent)))
	return fmt.Sprintf("---\n"+
		"spec: mechanical-carrier\n"+
		"status: closed\n"+
		"verdict: pass\n"+
		"rows_blocked_environment: 0\n"+
		"rows_blocked_finding: 0\n"+
		"rows_blocked_declared: 0\n"+
		"evidence_snapshots:\n"+
		"  R01:\n"+
		"    head: %s\n"+
		"    inputs:\n"+
		"      - ref: evidence.txt\n"+
		"        files:\n"+
		"          - path: evidence.txt\n"+
		"            sha256: %s\n"+
		"---\n\n"+
		"# QA report — carry fixture\n\n"+
		"## Results\n\n"+
		"| # | Story / criterion / sweep | Actor and surface | Status | Evidence |\n"+
		"| - | --- | --- | --- | --- |\n"+
		"| R01 | Unchanged evidence | maintainer / backend | pass | [evidence](../../../../evidence.txt) |\n\n"+
		"### R01 evidence\n\n"+
		"```yaml\n"+
		"inputs:\n"+
		"  - kind: repository_path\n"+
		"    ref: evidence.txt\n"+
		"```\n", establishedHead, digest)
}

func mechanicalCarriedReport(establishedBy, establishedHead string) string {
	return fmt.Sprintf("---\n"+
		"spec: mechanical-carrier\nstatus: closed\nverdict: pass\n"+
		"rows_blocked_environment: 0\nrows_blocked_finding: 0\nrows_blocked_declared: 0\n"+
		"---\n\n# QA report — carried fixture\n\n## Results\n\n"+
		"| # | Story / criterion / sweep | Actor and surface | Status | Evidence |\n"+
		"| - | --- | --- | --- | --- |\n"+
		"| R01 | Unchanged evidence | maintainer / backend | carried (established by: %s; head: %s) | inherited |\n",
		establishedBy, establishedHead)
}

func mechanicalEvidenceReport() string {
	return "---\n" +
		"spec: mechanical-carrier\nstatus: closed\nverdict: pass\n" +
		"rows_blocked_environment: 0\nrows_blocked_finding: 0\nrows_blocked_declared: 0\n" +
		"---\n\n# QA report — evidence fixture\n\n## Results\n\n" +
		"| # | Story / criterion / sweep | Actor and surface | Status | Evidence |\n" +
		"| - | --- | --- | --- | --- |\n" +
		"| R01 | Evidence hygiene | maintainer / backend | pass | observed inline |\n"
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
