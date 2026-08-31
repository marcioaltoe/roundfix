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

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func TestGateAcceptsItsOwnDeclaredTerm(t *testing.T) {
	t.Parallel()

	t.Run("own declared term becomes named gate input", func(t *testing.T) {
		t.Parallel()

		checked := checkFixture(t, "vocabulary-missing")
		precondition := speccheck.GatePrecondition(checked)
		if precondition.Blocking {
			t.Fatalf("GatePrecondition() Blocking = true, want false; findings = %#v", precondition.Findings)
		}
		if len(precondition.Findings) != 0 {
			t.Fatalf("GatePrecondition() Findings = %#v, want none", precondition.Findings)
		}
		if len(precondition.Inputs) != 1 || !strings.Contains(precondition.Inputs[0].Summary, "publish:") {
			t.Fatalf("GatePrecondition() Inputs = %#v, want pending publish: term", precondition.Inputs)
		}
	})

	t.Run("undeclared emitted term still blocks", func(t *testing.T) {
		t.Parallel()

		checked := speccheck.Result{
			Slug: "current-spec",
			Findings: []speccheck.Finding{{
				Code:     speccheck.CodeVocabularyUndocumented,
				Severity: speccheck.SeverityError,
				Summary:  `internal/example/emitter.go emits undocumented token "orphan:" absent from CONTEXT.md`,
				Where:    []speccheck.Location{{Path: "internal/example/emitter.go", Line: 7}},
				Fix:      `Document "orphan:" in CONTEXT.md.`,
			}},
		}
		precondition := speccheck.GatePrecondition(checked)
		if !precondition.Blocking || len(precondition.Findings) != 1 {
			t.Fatalf("GatePrecondition() = %#v, want undeclared term to block", precondition)
		}
		if len(precondition.Inputs) != 0 {
			t.Fatalf("GatePrecondition() Inputs = %#v, want no undeclared input", precondition.Inputs)
		}
	})

	t.Run("authoring stage still reports declared term", func(t *testing.T) {
		t.Parallel()

		repoRoot, err := filepath.Abs("testdata/repo")
		if err != nil {
			t.Fatalf("resolve fixture repository: %v", err)
		}
		checked, err := speccheck.CheckStage(fixtureSpecRoot, repoRoot, "vocabulary-missing", speccheck.StageTechSpec)
		if err != nil {
			t.Fatalf("CheckStage(StageTechSpec) error = %v", err)
		}
		findings := findingsWithCode(checked, speccheck.CodeVocabularyUndocumented)
		if len(findings) != 1 || !strings.Contains(findings[0].Summary, "publish:") {
			t.Fatalf("authoring findings = %#v, want pending publish: term", findings)
		}
	})
}

func TestGateRefusalNamesThePreconditionThatStoppedIt(t *testing.T) {
	t.Parallel()

	t.Run("a strict Spec check refusal names its check and every refusing code", func(t *testing.T) {
		t.Parallel()

		precondition := speccheck.GatePrecondition(checkFixture(t, "constraint-missing"))
		if !precondition.Blocking {
			t.Fatalf("GatePrecondition() Blocking = false, want the fixture's strict findings to stop the gate")
		}

		refusal, refused := speccheck.PreconditionRefusal(precondition)
		if !refused {
			t.Fatalf("PreconditionRefusal() refused = false, want a refusal for findings %#v", precondition.Findings)
		}
		if refusal.CheckName != speccheck.GatePreconditionCheck {
			t.Fatalf("refusal.CheckName = %q, want %q", refusal.CheckName, speccheck.GatePreconditionCheck)
		}
		for _, finding := range precondition.Findings {
			if !strings.Contains(refusal.Reason, finding.Code) {
				t.Fatalf("refusal.Reason = %q, want the refusing code %q", refusal.Reason, finding.Code)
			}
			if summary := oneLineRefusalText(finding.Summary); !strings.Contains(refusal.Reason, summary) {
				t.Fatalf("refusal.Reason = %q, want the reason %q", refusal.Reason, summary)
			}
		}
	})

	t.Run("a gate input alone refuses nothing", func(t *testing.T) {
		t.Parallel()

		precondition := speccheck.GatePrecondition(checkFixture(t, "vocabulary-missing"))
		if len(precondition.Inputs) == 0 {
			t.Fatalf("GatePrecondition() Inputs = none, want the fixture's own declared term as gate input")
		}
		refusal, refused := speccheck.PreconditionRefusal(precondition)
		if refused {
			t.Fatalf("PreconditionRefusal() = %#v, true; want no refusal when the gate is assigned repair inputs", refusal)
		}
		if refusal != (spec.PreconditionRefusal{}) {
			t.Fatalf("PreconditionRefusal() = %#v, want the zero refusal when nothing refused", refusal)
		}
	})

	t.Run("every distinct cause survives and a repeat is recorded once", func(t *testing.T) {
		t.Parallel()

		contradiction := speccheck.Finding{
			Code:     speccheck.CodeRequirementContradictory,
			Severity: speccheck.SeverityError,
			Summary:  "_prd.md requires one report per run and forbids writing one",
			Where:    []speccheck.Location{{Path: "docs/specs/gate-spec/_prd.md", Line: 12}},
			Fix:      "Resolve the contradiction.",
		}
		unmapped := speccheck.Finding{
			Code:     speccheck.CodeCoverageUnmapped,
			Severity: speccheck.SeverityError,
			Summary:  "Core Feature 2 has no TechSpec section",
			Where:    []speccheck.Location{{Path: "docs/specs/gate-spec/_techspec.md", Line: 20}},
			Fix:      "Map the feature.",
		}
		precondition := speccheck.GatePrecondition(speccheck.Result{
			Slug:     "gate-spec",
			Findings: []speccheck.Finding{contradiction, unmapped, contradiction},
		})

		refusal, refused := speccheck.PreconditionRefusal(precondition)
		if !refused {
			t.Fatalf("PreconditionRefusal() refused = false, want a refusal for three blocking findings")
		}
		want := speccheck.CodeRequirementContradictory + ": " + contradiction.Summary + "; " +
			speccheck.CodeCoverageUnmapped + ": " + unmapped.Summary
		if refusal.Reason != want {
			t.Fatalf("refusal.Reason = %q, want %q", refusal.Reason, want)
		}
	})

	t.Run("a refusal reason stays on one line", func(t *testing.T) {
		t.Parallel()

		precondition := speccheck.GatePrecondition(speccheck.Result{
			Slug: "gate-spec",
			Findings: []speccheck.Finding{{
				Code:     speccheck.CodeConstraintMissing,
				Severity: speccheck.SeverityError,
				Summary:  "_prd.md Project Constraints omits\nidentifier strategy\tand tooling authority",
				Fix:      "Account for every row.",
			}},
		})

		refusal, refused := speccheck.PreconditionRefusal(precondition)
		if !refused {
			t.Fatalf("PreconditionRefusal() refused = false, want a refusal for one blocking finding")
		}
		if strings.ContainsAny(refusal.Reason, "\n\t") {
			t.Fatalf("refusal.Reason = %q, want one line a report frontmatter can carry", refusal.Reason)
		}
		want := speccheck.CodeConstraintMissing +
			": _prd.md Project Constraints omits identifier strategy and tooling authority"
		if refusal.Reason != want {
			t.Fatalf("refusal.Reason = %q, want %q", refusal.Reason, want)
		}
	})

	t.Run("a refusal whose cause has no name still records the refusing check", func(t *testing.T) {
		t.Parallel()

		refusal, refused := speccheck.PreconditionRefusal(speccheck.GatePreconditionResult{
			Findings: []speccheck.Finding{{Severity: speccheck.SeverityError}},
			Blocking: true,
		})
		if !refused {
			t.Fatalf("PreconditionRefusal() refused = false, want an unnamed refusal recorded rather than dropped")
		}
		if refusal.CheckName != speccheck.GatePreconditionCheck {
			t.Fatalf("refusal.CheckName = %q, want %q", refusal.CheckName, speccheck.GatePreconditionCheck)
		}
		if refusal.Reason != "" {
			t.Fatalf("refusal.Reason = %q, want no invented reason", refusal.Reason)
		}
	})
}

func TestMechanicalStageStoresThePreconditionRefusalForTheReport(t *testing.T) {
	t.Parallel()

	t.Run("a refusing precondition reaches the refusal report writer", func(t *testing.T) {
		t.Parallel()

		repoRoot := newMechanicalGitRepo(t)
		blocking := speccheck.Result{
			Slug: "gate-spec",
			Findings: []speccheck.Finding{{
				Code:     speccheck.CodeRequirementContradictory,
				Severity: speccheck.SeverityError,
				Summary:  "_prd.md requires one report per run and forbids writing one",
				Where:    []speccheck.Location{{Path: "docs/specs/gate-spec/_prd.md", Line: 12}},
				Fix:      "Resolve the contradiction.",
			}},
		}

		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:     repoRoot,
			Precondition: speccheck.GatePrecondition(blocking),
		})

		if !result.PreconditionRefused {
			t.Fatalf("RunMechanicalStage() PreconditionRefused = false, want the refusal stored for report writing")
		}
		if !result.Blocking {
			t.Fatalf("RunMechanicalStage() Blocking = false, want a refusing precondition to block the gate")
		}
		if !mechanicalDetailsContain(result.Findings, blocking.Findings[0].Summary) {
			t.Fatalf("Findings = %#v, want the refusing finding still visible", result.Findings)
		}

		var report bytes.Buffer
		if err := spec.WritePreconditionRefusalReport(&report, result.PreconditionRefusal, spec.AuditorEvidence{}); err != nil {
			t.Fatalf("WritePreconditionRefusalReport(, spec.AuditorEvidence{}) error = %v", err)
		}
		for _, want := range []string{
			`precondition_check: "` + speccheck.GatePreconditionCheck + `"`,
			speccheck.CodeRequirementContradictory,
			"requires one report per run",
			"| 0 | blocked | precondition |",
		} {
			if !strings.Contains(report.String(), want) {
				t.Fatalf("refusal report does not record %q:\n%s", want, report.String())
			}
		}
	})

	t.Run("a passed precondition stores no refusal", func(t *testing.T) {
		t.Parallel()

		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:     newMechanicalGitRepo(t),
			Precondition: speccheck.GatePrecondition(checkFixture(t, "vocabulary-missing")),
		})

		if result.PreconditionRefused {
			t.Fatalf("RunMechanicalStage() PreconditionRefusal = %#v, want no refusal recorded", result.PreconditionRefusal)
		}
		if result.PreconditionRefusal != (spec.PreconditionRefusal{}) {
			t.Fatalf("RunMechanicalStage() PreconditionRefusal = %#v, want the zero refusal", result.PreconditionRefusal)
		}
	})
}

func oneLineRefusalText(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func TestGatePerformsAssignedRepairs(t *testing.T) {
	t.Parallel()

	t.Run("assigned repair is performed verified and recorded", func(t *testing.T) {
		t.Parallel()
		repoRoot := newMechanicalGitRepo(t)
		const repairPath = "docs/specs/mechanical/_prd.md"
		writeMechanicalFile(t, repoRoot, repairPath, "Tooling authority: not applicable\n")

		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:        repoRoot,
			TaskRepairPaths: []string{repairPath},
			AssignedRepairs: []speccheck.AssignedRepair{{
				ID:     "restore-tooling-authority",
				Path:   repairPath,
				Before: "Tooling authority: not applicable",
				After:  "Tooling authority: applicable",
			}},
		})

		content, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(repairPath)))
		if err != nil {
			t.Fatalf("read repaired path: %v", err)
		}
		if got := string(content); got != "Tooling authority: applicable\n" {
			t.Fatalf("repaired content = %q, want assigned replacement", got)
		}
		if len(result.Performed) != 1 || result.Performed[0].ID != "restore-tooling-authority" ||
			result.Performed[0].Path != repairPath {
			t.Fatalf("Performed = %#v, want the verified assigned repair", result.Performed)
		}
		if len(result.RepairFailures) != 0 {
			t.Fatalf("RepairFailures = %#v, want none", result.RepairFailures)
		}

		var report bytes.Buffer
		if err := speccheck.WriteMechanicalResult(&report, result); err != nil {
			t.Fatalf("WriteMechanicalResult: %v", err)
		}
		for _, want := range []string{"## Performed repairs", "restore-tooling-authority", repairPath, "verified after write"} {
			if !strings.Contains(report.String(), want) {
				t.Fatalf("mechanical report does not record %q:\n%s", want, report.String())
			}
		}
	})

	t.Run("skipped assigned repair blocks", func(t *testing.T) {
		t.Parallel()
		repoRoot := newMechanicalGitRepo(t)
		const repairPath = "CONTEXT.md"
		writeMechanicalFile(t, repoRoot, repairPath, "# Glossary\n")

		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:        repoRoot,
			TaskRepairPaths: []string{repairPath},
			AssignedRepairs: []speccheck.AssignedRepair{{
				ID:     "document-governed-path",
				Path:   repairPath,
				Before: "**Governed Path**: pending",
				After:  "**Governed Path**: documented",
			}},
		})

		if !result.Blocking || len(result.RepairFailures) != 1 || !strings.Contains(result.RepairFailures[0].Detail, "was not performed") {
			t.Fatalf("assigned repair result = %#v, want one blocking skipped-repair failure", result)
		}
		if len(result.Findings) != 0 {
			t.Fatalf("Findings = %#v, want assigned work kept out of observations", result.Findings)
		}
		if len(result.Performed) != 0 {
			t.Fatalf("Performed = %#v, want none", result.Performed)
		}
	})

	t.Run("unassigned finding is reported without a repair", func(t *testing.T) {
		t.Parallel()
		repoRoot := newMechanicalGitRepo(t)
		const reportPath = "docs/specs/mechanical/qa/report.md"
		original := "# QA Report without frontmatter or Results rows\n"
		writeMechanicalFile(t, repoRoot, reportPath, original)

		result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: reportPath})

		if findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalReportShape); len(findings) == 0 {
			t.Fatalf("%s findings = %#v, want the unassigned observation reported", speccheck.CodeMechanicalReportShape, findings)
		}
		content, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(reportPath)))
		if err != nil {
			t.Fatalf("read observed report: %v", err)
		}
		if string(content) != original {
			t.Fatalf("unassigned finding changed report to %q", content)
		}
		if len(result.Performed) != 0 {
			t.Fatalf("Performed = %#v, want none for an unassigned finding", result.Performed)
		}
	})

	t.Run("repair outside Task named paths is refused before writing", func(t *testing.T) {
		t.Parallel()
		repoRoot := newMechanicalGitRepo(t)
		const (
			allowedPath = "docs/specs/mechanical/_prd.md"
			outsidePath = "CONTEXT.md"
		)
		writeMechanicalFile(t, repoRoot, allowedPath, "allowed\n")
		writeMechanicalFile(t, repoRoot, outsidePath, "original\n")

		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:        repoRoot,
			TaskRepairPaths: []string{allowedPath},
			AssignedRepairs: []speccheck.AssignedRepair{{
				ID:     "out-of-bounds",
				Path:   outsidePath,
				Before: "original",
				After:  "changed",
			}},
		})

		if !result.Blocking || len(result.RepairFailures) != 1 || !strings.Contains(result.RepairFailures[0].Detail, "outside the Task-named repair paths") {
			t.Fatalf("out-of-bounds result = %#v, want one blocking scope failure", result)
		}
		if len(result.Findings) != 0 {
			t.Fatalf("Findings = %#v, want assigned work kept out of observations", result.Findings)
		}
		content, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(outsidePath)))
		if err != nil {
			t.Fatalf("read out-of-bounds path: %v", err)
		}
		if got := string(content); got != "original\n" {
			t.Fatalf("out-of-bounds content = %q, want unchanged", got)
		}
		if len(result.Performed) != 0 {
			t.Fatalf("Performed = %#v, want none", result.Performed)
		}
	})
}

func TestMechanicalAuthPaths(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	copyMechanicalFixture(t, repoRoot, "authorization.md", "docs/workflow/authorizations/mechanical.md")
	writeMechanicalFile(t, repoRoot, "Makefile", "verify:\n\t@true\n")
	greenCommit := commitMechanicalFiles(t, repoRoot, "authorized change", "Makefile")
	writeMechanicalFile(t, repoRoot, ".golangci.yml", "linters: {}\n")
	redCommit := commitMechanicalFiles(t, repoRoot, "unauthorized change", ".golangci.yml")

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
		if len(findings) != 1 || !strings.Contains(findings[0].Detail, ".golangci.yml") {
			t.Fatalf("%s findings = %#v, want .golangci.yml", speccheck.CodeMechanicalAuthPaths, findings)
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
	writeMechanicalResolverFixture(t, repoRoot)
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
	writeMechanicalResolverFixture(t, repoRoot)
	writeMechanicalFile(t, repoRoot, "docs/workflow/authorizations/mechanical.md", mechanicalRegenerationAuthorization("command: make baseline-digests\noutputs:\n  - internal/baseline/testdata/catalog.digest\n"))
	writeMechanicalFile(t, repoRoot, ".golangci.yml", "linters: {}\n")
	commit := commitMechanicalFiles(t, repoRoot, "change undeclared path", ".golangci.yml")

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
	if len(findings) != 1 || !strings.Contains(findings[0].Detail, ".golangci.yml") {
		t.Fatalf("%s findings = %#v, want .golangci.yml", speccheck.CodeMechanicalAuthPaths, findings)
	}
}

func TestMechanicalAuthPathsRefusesInvalidRegenerationDeclaration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		declaration string
	}{
		{name: "output without command", declaration: "outputs:\n  - internal/baseline/assets/source-baselines/index.json\n"},
		{name: "output glob", declaration: "command: make baseline-digests\noutputs:\n  - internal/baseline/assets/source-baselines/*.json\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repoRoot := newMechanicalGitRepo(t)
			writeMechanicalFile(t, repoRoot, "docs/workflow/authorizations/mechanical.md", mechanicalRegenerationAuthorization(tt.declaration))
			const derivedPath = "internal/baseline/assets/source-baselines/index.json"
			writeMechanicalFile(t, repoRoot, derivedPath, "hand edited\n")
			commit := commitMechanicalFiles(t, repoRoot, "hand edit baseline digest", derivedPath)

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
			if len(findings) != 1 || !strings.Contains(findings[0].Detail, derivedPath) {
				t.Fatalf("%s findings = %#v, want invalid regeneration declaration refused", speccheck.CodeMechanicalAuthPaths, findings)
			}
		})
	}
}

func TestAuditJudgesTheGrant(t *testing.T) {
	const (
		archiveAuthorization   = "docs/workflow/authorizations/2026-08-12-the-archive-root-under-docs.md"
		authoringAuthorization = "docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md"
		archiveHelpCommit      = "419a4661ac769ff7ee6ce5423bd795185c859d01"
		archiveCarrierCommit   = "65c51ebf2e19220ff50d25fe03be809fcdf353f0"
		verificationTaskCommit = "28acf39cc193ad490646cb5a1d23500e0c08c273"
		regeneratedCommit      = "c80e1266658929f68e8046af82f88e13392dc56d"
		verificationTaskFile   = "docs/specs/0095-a-verification-that-ran-before-anyone-believed-it/task_08.md"
		archiveHelpTaskFile    = "docs/specs/0094-one-history-root-under-docs/task_15.md"
		archiveCarrierTaskFile = "docs/specs/0094-one-history-root-under-docs/task_16.md"
	)
	repository := filepath.Clean(filepath.Join("..", ".."))

	t.Run("historical authorized asset and ordinary Go split now share one audit", func(t *testing.T) {
		if missing := firstMissingMechanicalCommit(repository, archiveHelpCommit, archiveCarrierCommit, verificationTaskCommit); missing != "" {
			t.Logf("outside evidence blocked: historical commit %s cannot be resolved", missing)
			return
		}

		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:          repository,
			AuthorizationPath: archiveAuthorization,
			TaskCommits: []speccheck.MechanicalTaskCommit{
				{
					TaskID:   "task_15",
					SHA:      archiveHelpCommit,
					TaskFile: archiveHelpTaskFile,
				},
				{
					TaskID:   "task_16",
					SHA:      archiveCarrierCommit,
					TaskFile: archiveCarrierTaskFile,
				},
			},
		})

		assertNoMechanicalCode(t, result, speccheck.CodeMechanicalAuthPaths)

		result = runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:          repository,
			AuthorizationPath: authoringAuthorization,
			TaskCommits: []speccheck.MechanicalTaskCommit{{
				TaskID:   "task_08",
				SHA:      verificationTaskCommit,
				TaskFile: verificationTaskFile,
			}},
		})

		assertNoMechanicalCode(t, result, speccheck.CodeMechanicalAuthPaths)
	})

	t.Run("historical regeneration resolves outputs absent from the grant", func(t *testing.T) {
		if missing := firstMissingMechanicalCommit(repository, regeneratedCommit); missing != "" {
			t.Logf("outside evidence blocked: historical commit %s cannot be resolved", missing)
			return
		}

		repoRoot := cloneMechanicalHistory(t, repository)
		const authorizationPath = "docs/workflow/authorizations/mechanical-history.md"
		writeMechanicalFile(t, repoRoot, authorizationPath,
			"---\nconsuming: 0094-one-history-root-under-docs\npaths:\n"+
				"  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md\n"+
				"---\n\n# Historical authorization replay\n\n## Sanctioned regeneration\n\n"+
				"```yaml\ncommand: make baseline-digests\n```\n")

		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:          repoRoot,
			AuthorizationPath: authorizationPath,
			TaskCommits: []speccheck.MechanicalTaskCommit{{
				TaskID:   "task_16",
				SHA:      regeneratedCommit,
				TaskFile: archiveCarrierTaskFile,
			}},
		})

		assertNoMechanicalCode(t, result, speccheck.CodeMechanicalAuthPaths)
	})

	t.Run("governed path outside the grant is refused by name", func(t *testing.T) {
		repoRoot := newMechanicalGitRepo(t)
		const authorizationPath = "docs/workflow/authorizations/mechanical.md"
		writeMechanicalFile(t, repoRoot, authorizationPath, "# Authorization\n\n## Bounded files\n\n- `Makefile`\n")
		writeMechanicalFile(t, repoRoot, "Makefile", "verify:\n\t@true\n")
		writeMechanicalFile(t, repoRoot, ".golangci.yml", "linters: {}\n")
		commit := commitMechanicalFiles(t, repoRoot, "change granted build and ungranted linter configuration", "Makefile", ".golangci.yml")

		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:          repoRoot,
			AuthorizationPath: authorizationPath,
			TaskCommits:       []speccheck.MechanicalTaskCommit{{TaskID: "task_01", SHA: commit}},
		})

		assertMechanicalPathEscapedGrant(t, result, ".golangci.yml", authorizationPath)
	})

	t.Run("hand edited derived value without a command is refused", func(t *testing.T) {
		repoRoot := newMechanicalGitRepo(t)
		const (
			authorizationPath = "docs/workflow/authorizations/mechanical.md"
			derivedPath       = "internal/baseline/assets/source-baselines/index.json"
		)
		writeMechanicalFile(t, repoRoot, authorizationPath, "# Authorization\n\n## Bounded files\n\n- `Makefile`\n")
		writeMechanicalFile(t, repoRoot, derivedPath, "hand edited\n")
		commit := commitMechanicalFiles(t, repoRoot, "hand edit derived value", derivedPath)

		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:          repoRoot,
			AuthorizationPath: authorizationPath,
			TaskCommits:       []speccheck.MechanicalTaskCommit{{TaskID: "task_02", SHA: commit}},
		})

		assertMechanicalPathEscapedGrant(t, result, derivedPath, authorizationPath)
	})

	t.Run("record that does not name the Spec is refused", func(t *testing.T) {
		result := checkFixture(t, "tooling-unauthorized")
		finding := requireFinding(t, result, speccheck.CodeToolingUnauthorized)
		if !strings.Contains(finding.Summary, "does not name Spec tooling-unauthorized") {
			t.Fatalf("%s summary = %q, want the unnamed Spec", speccheck.CodeToolingUnauthorized, finding.Summary)
		}
	})

	t.Run("Task commit cannot carry its own authorization", func(t *testing.T) {
		repoRoot := newMechanicalGitRepo(t)
		const authorizationPath = "docs/workflow/authorizations/mechanical.md"
		writeMechanicalFile(t, repoRoot, authorizationPath, "# Authorization\n\n## Bounded files\n\n- `Makefile`\n")
		writeMechanicalFile(t, repoRoot, "Makefile", "verify:\n\t@true\n")
		commit := commitMechanicalFiles(t, repoRoot, "fold grant into authorized change", authorizationPath, "Makefile")

		result := runMechanical(t, speccheck.MechanicalRequest{
			RepoRoot:          repoRoot,
			AuthorizationPath: authorizationPath,
			TaskCommits:       []speccheck.MechanicalTaskCommit{{TaskID: "task_03", SHA: commit}},
		})

		assertMechanicalPathEscapedGrant(t, result, authorizationPath, authorizationPath)
	})
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

func TestBlockedCauseDiagnosticNamesTheLiteral(t *testing.T) {
	t.Parallel()

	const (
		requiredLiteral = `" — waits on "`
		wrongTypeDetail = "blocked cause outside environment, finding, or declared"
	)
	tests := []struct {
		name           string
		status         string
		environment    int
		finding        int
		declared       int
		wantDetail     string
		rejectedDetail string
	}{
		{
			name:           "finding type without required literal names the literal",
			status:         "blocked (finding: QA-FIXTURE)",
			wantDetail:     requiredLiteral,
			rejectedDetail: wrongTypeDetail,
		},
		{
			name:           "unrecognised blocked cause names the three types",
			status:         "blocked (external: outage)",
			wantDetail:     wrongTypeDetail,
			rejectedDetail: requiredLiteral,
		},
		{
			name:        "environment cause keeps its typed count",
			status:      "blocked (environment: unavailable)",
			environment: 1,
		},
		{
			name:    "finding cause with required literal keeps its typed count",
			status:  "blocked (finding: QA-FIXTURE — waits on fixture)",
			finding: 1,
		},
		{
			name:     "declared cause keeps its typed count",
			status:   "blocked (declared: unavailable)",
			declared: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repoRoot := newMechanicalGitRepo(t)
			reportPath := "docs/specs/mechanical/qa/report.md"
			writeMechanicalFile(t, repoRoot, reportPath, fmt.Sprintf("---\n"+
				"verdict: fail\n"+
				"rows_blocked_environment: %d\n"+
				"rows_blocked_finding: %d\n"+
				"rows_blocked_declared: %d\n"+
				"---\n\n"+
				"# QA Report\n\n"+
				"## Results\n\n"+
				"| # | Status | Evidence |\n"+
				"| - | --- | --- |\n"+
				"| R01 | %s | observed inline |\n",
				tt.environment, tt.finding, tt.declared, tt.status))

			result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: reportPath})
			findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalReportShape)
			if tt.wantDetail == "" {
				if len(findings) != 0 {
					t.Fatalf("%s findings = %#v, want none", speccheck.CodeMechanicalReportShape, findings)
				}
				return
			}
			if len(findings) != 1 || !strings.Contains(findings[0].Detail, tt.wantDetail) {
				t.Fatalf("%s findings = %#v, want one detail containing %q", speccheck.CodeMechanicalReportShape, findings, tt.wantDetail)
			}
			if strings.Contains(findings[0].Detail, tt.rejectedDetail) {
				t.Fatalf("%s detail = %q, must differ from %q", speccheck.CodeMechanicalReportShape, findings[0].Detail, tt.rejectedDetail)
			}
		})
	}
}

func TestCountDisagreementReportsItsCause(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		rows             string
		declaredFinding  int
		wantFindings     int
		wantParseRow     string
		wantCountFinding bool
	}{
		{
			name:            "unparsed row accounts for count disagreement",
			rows:            "| R-PARSE | blocked (finding: QA-FIXTURE) | observed inline |\n",
			declaredFinding: 1,
			wantFindings:    1,
			wantParseRow:    "R-PARSE",
		},
		{
			name:             "parsed rows expose genuine count disagreement",
			rows:             "| R-PARSED | blocked (finding: QA-FIXTURE — waits on fixture) | observed inline |\n",
			declaredFinding:  2,
			wantFindings:     1,
			wantCountFinding: true,
		},
		{
			name: "unparsed row and wrong total expose both causes",
			rows: "| R-PARSE | blocked (finding: QA-FIXTURE) | observed inline |\n" +
				"| R-PARSED | blocked (finding: QA-OTHER — waits on other fixture) | observed inline |\n",
			declaredFinding:  3,
			wantFindings:     2,
			wantParseRow:     "R-PARSE",
			wantCountFinding: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repoRoot := newMechanicalGitRepo(t)
			reportPath := "docs/specs/mechanical/qa/report.md"
			writeMechanicalFile(t, repoRoot, reportPath, fmt.Sprintf("---\n"+
				"verdict: fail\n"+
				"rows_blocked_environment: 0\n"+
				"rows_blocked_finding: %d\n"+
				"rows_blocked_declared: 0\n"+
				"---\n\n"+
				"# QA Report\n\n"+
				"## Results\n\n"+
				"| # | Status | Evidence |\n"+
				"| - | --- | --- |\n%s",
				tt.declaredFinding, tt.rows))

			result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: reportPath})
			findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalReportShape)
			if len(findings) != tt.wantFindings {
				t.Fatalf("%s findings = %#v, want %d", speccheck.CodeMechanicalReportShape, findings, tt.wantFindings)
			}
			if tt.wantParseRow != "" && !mechanicalDetailsContain(findings, tt.wantParseRow) {
				t.Fatalf("%s findings = %#v, want parse failure naming row %s", speccheck.CodeMechanicalReportShape, findings, tt.wantParseRow)
			}
			if got := mechanicalDetailsContain(findings, "rows_blocked_finding"); got != tt.wantCountFinding {
				t.Fatalf("%s count finding present = %t, want %t: %#v", speccheck.CodeMechanicalReportShape, got, tt.wantCountFinding, findings)
			}
		})
	}
}

func TestMechanicalStageAcceptsThePreconditionRefusalRow(t *testing.T) {
	t.Parallel()

	t.Run("the refusal a gate writes passes the stage that reads it", func(t *testing.T) {
		t.Parallel()
		repoRoot := newMechanicalGitRepo(t)
		reportPath := "docs/specs/mechanical/qa/qa-report-2026-08-25.md"
		var refusal bytes.Buffer
		if err := spec.WritePreconditionRefusalReport(&refusal, spec.PreconditionRefusal{
			CheckName: speccheck.GatePreconditionCheck,
			Reason:    "SC-VOCABULARY-UNDOCUMENTED: undocumented emitted token",
		}, spec.AuditorEvidence{}); err != nil {
			t.Fatalf("WritePreconditionRefusalReport(, spec.AuditorEvidence{}) error = %v", err)
		}
		writeMechanicalFile(t, repoRoot, reportPath, refusal.String())

		result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: reportPath})

		assertNoMechanicalCode(t, result, speccheck.CodeMechanicalReportShape)
		if result.Blocking {
			t.Fatalf("Blocking = true, want a refusal report to leave the next run free; findings = %#v", result.Findings)
		}
	})

	t.Run("an empty Results table still refuses", func(t *testing.T) {
		t.Parallel()
		repoRoot := newMechanicalGitRepo(t)
		reportPath := "docs/specs/mechanical/qa/report-empty.md"
		writeMechanicalFile(t, repoRoot, reportPath, "---\n"+
			"verdict: fail\n"+
			"rows_blocked_environment: 0\n"+
			"rows_blocked_finding: 0\n"+
			"rows_blocked_declared: 0\n"+
			"rows_blocked_precondition: 0\n"+
			"---\n\n"+
			"# QA Report\n\n"+
			"## Results\n\n"+
			"| # | Status | Provenance |\n"+
			"| - | --- | --- |\n")

		result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: reportPath})

		findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalReportShape)
		if len(findings) != 1 || !strings.Contains(findings[0].Detail, "Results table has no report rows") {
			t.Fatalf("%s findings = %#v, want the empty-matrix refusal", speccheck.CodeMechanicalReportShape, findings)
		}
	})

	tests := []struct {
		name       string
		declared   string
		row        string
		wantDetail string
	}{
		{
			name:     "blocked beside precondition provenance is terminal",
			declared: "rows_blocked_precondition: 1\n",
			row:      "| 0 | blocked | precondition |\n",
		},
		{
			name:       "precondition provenance with a non-blocked status refuses",
			declared:   "rows_blocked_precondition: 0\n",
			row:        "| 0 | pending | precondition |\n",
			wantDetail: `row 0 has provenance precondition with status "pending"`,
		},
		{
			name:       "bare blocked without precondition provenance still refuses",
			row:        "| 0 | blocked | measured |\n",
			wantDetail: "blocked cause outside environment, finding, or declared",
		},
		{
			name:       "a declared count that no refusal row matches refuses",
			declared:   "rows_blocked_precondition: 2\n",
			row:        "| 0 | blocked | precondition |\n",
			wantDetail: "rows_blocked_precondition is 2 but the Results table contains 1 matching rows",
		},
		{
			name:       "a refusal row without its declared count refuses",
			row:        "| 0 | blocked | precondition |\n",
			wantDetail: "rows_blocked_precondition is absent from report frontmatter",
		},
		{
			name: "a report that records no refusal need not declare the count",
			row:  "| R01 | pass | measured |\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repoRoot := newMechanicalGitRepo(t)
			reportPath := "docs/specs/mechanical/qa/report.md"
			writeMechanicalFile(t, repoRoot, reportPath, "---\n"+
				"verdict: fail\n"+
				"rows_blocked_environment: 0\n"+
				"rows_blocked_finding: 0\n"+
				"rows_blocked_declared: 0\n"+
				tt.declared+
				"---\n\n"+
				"# QA Report\n\n"+
				"## Results\n\n"+
				"| # | Status | Provenance |\n"+
				"| - | --- | --- |\n"+
				tt.row)

			result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: reportPath})

			findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalReportShape)
			if tt.wantDetail == "" {
				if len(findings) != 0 {
					t.Fatalf("%s findings = %#v, want none", speccheck.CodeMechanicalReportShape, findings)
				}
				return
			}
			if len(findings) != 1 || !strings.Contains(findings[0].Detail, tt.wantDetail) {
				t.Fatalf("%s findings = %#v, want one detail containing %q", speccheck.CodeMechanicalReportShape, findings, tt.wantDetail)
			}
		})
	}
}

func TestResultsTableIsTheOnlyRowSource(t *testing.T) {
	t.Parallel()

	const frontmatter = "---\n" +
		"verdict: pass\n" +
		"rows_blocked_environment: 0\n" +
		"rows_blocked_finding: 0\n" +
		"rows_blocked_declared: 0\n" +
		"---\n\n" +
		"# QA Report\n\n"
	const results = "## Results\n\n" +
		"| # | Status | Provenance |\n" +
		"| - | --- | --- |\n" +
		"| R01 | pass | measured |\n"
	// The comparison Spec 0098's gate wrote to justify a Results row on
	// 2026-08-26. Its cells are the shape that became four blockers: a header
	// naming a case, and rows naming an observation rather than a status.
	const comparison = "| Case | Hook objection observed | Refused settle | Recovered settle |\n" +
		"| --- | --- | --- | --- |\n" +
		"| 82-line function vs 80 | function exceeds the 80-line limit | exit `1`, work staged | exit `0`, byte-identical |\n" +
		"| 2462-line generated file vs 500 | 2462 lines exceeds the 500-line limit | exit `1`, work staged | exit `0`, byte-identical |\n" +
		"| `sort()` vs `toSorted()` | use toSorted() instead of sort() | exit `1`, work staged | exit `0`, byte-identical |\n"

	tests := []struct {
		name       string
		report     string
		wantDetail string
	}{
		{
			name:   "a comparison under a row-detail subsection is not a result",
			report: frontmatter + results + "\n### Row detail — the three measured cases\n\nEach was built as its own repository:\n\n" + comparison,
		},
		{
			name:   "a comparison in prose under no heading is not a result",
			report: frontmatter + results + "\nEach was built as its own repository:\n\n" + comparison,
		},
		{
			name:   "a comparison under a later section is not a result",
			report: frontmatter + results + "\n## Findings\n\n" + comparison,
		},
		{
			name: "a defective Results row is the only finding beside a comparison",
			report: frontmatter + "## Results\n\n" +
				"| # | Status | Provenance |\n" +
				"| - | --- | --- |\n" +
				"| R01 | pending | measured |\n\n" +
				"### Row detail — the three measured cases\n\n" + comparison,
			wantDetail: "row R01 remains pending instead of carrying a terminal status",
		},
		{
			name:       "a report with no Results heading still reports the missing matrix",
			report:     frontmatter + "## Evidence\n\n" + comparison,
			wantDetail: "Results table has no report rows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repoRoot := newMechanicalGitRepo(t)
			reportPath := "docs/specs/mechanical/qa/report.md"
			writeMechanicalFile(t, repoRoot, reportPath, tt.report)

			result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: reportPath})

			findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalReportShape)
			if tt.wantDetail == "" {
				if len(findings) != 0 {
					t.Fatalf("%s findings = %#v, want none for a table outside the Results matrix", speccheck.CodeMechanicalReportShape, findings)
				}
				if result.Blocking {
					t.Fatalf("Blocking = true, want a gate's own evidence to leave the next run free; findings = %#v", result.Findings)
				}
				return
			}
			if len(findings) != 1 || !strings.Contains(findings[0].Detail, tt.wantDetail) {
				t.Fatalf("%s findings = %#v, want exactly one detail containing %q", speccheck.CodeMechanicalReportShape, findings, tt.wantDetail)
			}
		})
	}
}

// TestMechanicalStageReadsTheNewestQAReportOnly pins the property a superseded
// report violated: the stage validates the newest QA Report in the directory it
// was pointed at and ignores every older one, so a refusal a previous run wrote
// cannot block the run that supersedes it. Measured on Spec 0103 on 2026-08-14
// and on Spec 0098 on 2026-08-26, where the only exit was deleting the report.
func TestMechanicalStageReadsTheNewestQAReportOnly(t *testing.T) {
	t.Parallel()

	const directory = "docs/specs/mechanical/qa/"
	const (
		superseded = directory + "qa-report-2026-08-14.md"
		newest     = directory + "qa-report-2026-08-15.md"
	)

	tests := []struct {
		name      string
		reports   map[string]string
		requested string
		// wantRow is the row a QA-REPORT-SHAPE finding must name, empty when the
		// stage must raise none, and wantFile the report that finding must cite —
		// which report was read is proven by the finding, not by its absence.
		wantRow  string
		wantFile string
	}{
		{
			name: "a superseded refusal does not block the run that supersedes it",
			reports: map[string]string{
				superseded: mechanicalRowStatusReport("R01", "pending"),
				newest:     mechanicalRowStatusReport("R01", "pass"),
			},
			requested: superseded,
		},
		{
			name: "the newest report is the one read, not merely the older one skipped",
			reports: map[string]string{
				superseded: mechanicalRowStatusReport("R01", "pass"),
				newest:     mechanicalRowStatusReport("R02", "pending"),
			},
			requested: superseded,
			wantRow:   "R02",
			wantFile:  newest,
		},
		{
			// Raw filename order puts qa-report-2026-08-25.md above both of its
			// same-date reruns, because "." sorts above "-", so a stage that
			// sorted by name would read the defective first report of the date.
			name: "recency is the date and run sequence, not raw filename order",
			reports: map[string]string{
				directory + "qa-report-2026-08-25.md":    mechanicalRowStatusReport("R01", "pending"),
				directory + "qa-report-2026-08-25-02.md": mechanicalRowStatusReport("R02", "pending"),
				directory + "qa-report-2026-08-25-10.md": mechanicalRowStatusReport("R03", "pass"),
			},
			requested: directory + "qa-report-2026-08-25.md",
		},
		{
			name:      "a requested report that is gone yields to the newest one on disk",
			reports:   map[string]string{newest: mechanicalRowStatusReport("R02", "pending")},
			requested: directory + "qa-report-2026-08-13.md",
			wantRow:   "R02",
			wantFile:  newest,
		},
		{
			// A path outside the qa-report-*.md family has no family to be the
			// newest of, so it is read as named rather than redirected.
			name: "a report outside the QA Report family is read as named",
			reports: map[string]string{
				directory + "report.md": mechanicalRowStatusReport("R01", "pending"),
				newest:                  mechanicalRowStatusReport("R01", "pass"),
			},
			requested: directory + "report.md",
			wantRow:   "R01",
			wantFile:  directory + "report.md",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repoRoot := newMechanicalGitRepo(t)
			for path, content := range tt.reports {
				writeMechanicalFile(t, repoRoot, path, content)
			}

			result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: tt.requested})

			findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalReportShape)
			if tt.wantRow == "" {
				if len(findings) != 0 {
					t.Fatalf("%s findings = %#v, want none from a superseded report", speccheck.CodeMechanicalReportShape, findings)
				}
				if result.Blocking {
					t.Fatalf("Blocking = true, want an older report to leave the newer run free; findings = %#v", result.Findings)
				}
				return
			}
			if len(findings) != 1 || !strings.Contains(findings[0].Detail, "row "+tt.wantRow+" remains pending") {
				t.Fatalf("%s findings = %#v, want exactly one naming row %s", speccheck.CodeMechanicalReportShape, findings, tt.wantRow)
			}
			if findings[0].File != tt.wantFile {
				t.Fatalf("%s finding file = %q, want the report the stage read, %q", speccheck.CodeMechanicalReportShape, findings[0].File, tt.wantFile)
			}
		})
	}

	t.Run("a QA directory with no report keeps the requested path in its skip", func(t *testing.T) {
		t.Parallel()
		repoRoot := newMechanicalGitRepo(t)
		requested := directory + "qa-report-2026-08-15.md"

		result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: requested})

		assertMechanicalSkip(t, result, speccheck.DetectorMechanicalReportShape, requested)
		assertNoMechanicalCode(t, result, speccheck.CodeMechanicalReportShape)
	})
}

// mechanicalRowStatusReport is one QA Report whose Results table carries a
// single row with the given status: a pass row the stage accepts, a pending row
// it refuses by name. Naming the row in the refusal is what makes the report the
// stage actually read observable.
func mechanicalRowStatusReport(rowID, status string) string {
	verdict := spec.VerdictFail
	if status == "pass" {
		verdict = spec.VerdictPass
	}
	return "---\n" +
		"verdict: " + verdict + "\n" +
		"rows_blocked_environment: 0\nrows_blocked_finding: 0\nrows_blocked_declared: 0\n" +
		"---\n\n# QA Report\n\n## Results\n\n" +
		"| # | Status | Provenance |\n| - | --- | --- |\n" +
		"| " + rowID + " | " + status + " | measured |\n"
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
	writeMechanicalFile(t, repoRoot, ".golangci.yml", "linters: {}\n")
	commit := commitMechanicalFiles(t, repoRoot, "outside", ".golangci.yml")

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

func writeMechanicalResolverFixture(t *testing.T, repoRoot string) {
	t.Helper()
	writeMechanicalFile(t, repoRoot, "Makefile", "DERIVED_DIGEST_PATHS := internal/baseline/derived\n")
	writeMechanicalFile(t, repoRoot, "internal/baseline/derived/_ownership.yml", "owner: frozen\nreason: fixture\n")
	writeMechanicalFile(t, repoRoot, "internal/baseline/derived/frozen.txt", "frozen\n")
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

func firstMissingMechanicalCommit(repoRoot string, commits ...string) string {
	for _, commit := range commits {
		command := exec.Command("git", "-C", repoRoot, "cat-file", "-e", commit+"^{commit}")
		if err := command.Run(); err != nil {
			return commit
		}
	}
	return ""
}

func cloneMechanicalHistory(t *testing.T, source string) string {
	t.Helper()
	repoRoot := filepath.Join(t.TempDir(), "repository")
	command := exec.Command("git", "clone", "--quiet", "--shared", source, repoRoot)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("clone historical repository: %v: %s", err, output)
	}
	return repoRoot
}

func assertMechanicalPathEscapedGrant(
	t *testing.T,
	result speccheck.MechanicalResult,
	path string,
	grant string,
) {
	t.Helper()
	findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalAuthPaths)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want one finding for %s", speccheck.CodeMechanicalAuthPaths, findings, path)
	}
	if !strings.Contains(findings[0].Detail, path) || findings[0].File != grant {
		t.Fatalf("%s finding = %#v, want path %s escaping grant %s", speccheck.CodeMechanicalAuthPaths, findings[0], path, grant)
	}
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
