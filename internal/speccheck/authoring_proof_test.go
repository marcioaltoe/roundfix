// Suite: TechSpec structural authoring proofs
// Invariant: every active TechSpec carries absence evidence, listed ADR alternatives,
// and a Coverage Cross-Read; fenced examples never satisfy a declaration.
// Boundary IN: temporary Spec artifacts through Check and CheckStage.
// Boundary OUT: semantic judgment, authored command execution, and Task consumers.
package speccheck_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

const proofAbsenceCode = speccheck.CodeAbsenceUnproven
const proofADRCode = speccheck.CodeADRAlternativesUnrelated
const proofCoverageCode = speccheck.CodeCoverageUnread

func TestAbsenceEvidenceGapFollowsTheSection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		section string
		count   int
	}{
		{"missing", "", 1},
		{"empty", "## Absence Evidence\n", 1},
		{"prose", "## Absence Evidence\n\nThe component is absent.\n", 1},
		{"none without reason", "## Absence Evidence\n\nNone.\n", 1},
		{"reason in another paragraph", "## Absence Evidence\n\nNone.\n\nNo component is new.\n", 1},
		{"missing command", "## Absence Evidence\n\n1. Detector — result: absent.\n", 1},
		{"plain command", "## Absence Evidence\n\n1. Detector — command: rg detector; result: absent.\n", 1},
		{"blank command", "## Absence Evidence\n\n1. Detector — command: ` `; result: absent.\n", 1},
		{"missing result", "## Absence Evidence\n\n1. Detector — command: `rg detector`; result: \n", 1},
		{"two incomplete items", "## Absence Evidence\n\n1. First — command: `rg first`\n2. Second — result: absent.\n", 2},
		{"fenced section", "```markdown\n## Absence Evidence\n\nNone. No component is new.\n```\n", 1},
		{"fenced item", "## Absence Evidence\n\n~~~markdown\n1. Detector — command: `rg detector`; result: absent.\n~~~\n", 1},
		{"none with reason", "## Absence Evidence\n\nNone. No component is new.\n", 0},
		{"wrapped reason", "## Absence Evidence\n\nNone.\nNo component is new.\n", 0},
		{"complete item", "## Absence Evidence\n\n1. Detector — command: `rg detector`; result: absent.\n", 0},
		{"wrapped item", "## Absence Evidence\n\n1. Detector — command: `rg detector`;\n   result: absent.\n", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := writeProofFixture(t, "", "", tt.section+"\n## Coverage Map\n\nDeviations from the PRD: None.\n")
			result := checkProofFixture(t, root, speccheck.StageTechSpec)
			requireProofCount(t, result, proofAbsenceCode, tt.count)
			for _, finding := range findingsWithCode(result, proofAbsenceCode) {
				if !hasLocation(finding, "docs/specs/proof/_techspec.md") {
					t.Errorf("locations = %#v, want TechSpec line", finding.Where)
				}
			}
		})
	}

	marker := filepath.Join(t.TempDir(), "executed")
	root := writeProofFixture(t, "", "", "## Absence Evidence\n\n1. Detector — command: `touch "+filepath.ToSlash(marker)+"`; result: not executed.\n\n## Coverage Map\n\nDeviations from the PRD: None.\n")
	for _, stage := range []speccheck.Stage{speccheck.StageTechSpec, speccheck.StageAll} {
		requireProofCount(t, checkProofFixture(t, root, stage), proofAbsenceCode, 0)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("authored command marker stat = %v, want absent", err)
	}
}

func TestADRAlternativesGapNamesEachListedADR(t *testing.T) {
	t.Parallel()

	const prdADRs = "applicable. ADR-0001 and 0003 govern this fixture."
	const techADRs = "applicable. ADR-0001 and ADR-0002 govern this fixture."
	const otherProofs = "## Absence Evidence\n\nNone. No component is new.\n\n## Coverage Map\n\nDeviations from the PRD: None.\n\n"
	tests := []struct {
		name    string
		section string
		missing []string
	}{
		{"missing section", "", []string{"0001", "0002", "0003"}},
		{"partial", "## ADR Alternatives\n\n- ADR-0001: No rejected alternative is taken.\n", []string{"0002", "0003"}},
		{"blank relation", "## ADR Alternatives\n\n- ADR-0001: \n- ADR-0002: Accounted for.\n- ADR-0003: Accounted for.\n", []string{"0001"}},
		{"fenced relation", "## ADR Alternatives\n\n```md\n- ADR-0001: Accounted for.\n```\n- ADR-0002: Accounted for.\n- ADR-0003: Accounted for.\n", []string{"0001"}},
		{"complete", "## ADR Alternatives\n\n- ADR-0001: Accounted for.\n- ADR-0002: Accounted for.\n- ADR-0003: Accounted for.\n", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := writeProofFixture(t, prdADRs, techADRs, otherProofs+tt.section)
			result := checkProofFixture(t, root, speccheck.StageTechSpec)
			requireProofCount(t, result, proofADRCode, len(tt.missing))
			for _, number := range tt.missing {
				found := false
				for _, finding := range findingsWithCode(result, proofADRCode) {
					if !strings.Contains(finding.Summary, "ADR-"+number) {
						continue
					}
					found = true
					artifact := "docs/specs/proof/_prd.md"
					if number == "0002" {
						artifact = "docs/specs/proof/_techspec.md"
					}
					if !hasExactLocation(finding, artifact, 7) || !hasLocation(finding, "docs/specs/proof/_techspec.md") {
						t.Errorf("ADR-%s locations = %#v, want obligations row and TechSpec section", number, finding.Where)
					}
				}
				if !found {
					t.Errorf("missing ADR-%s gap: %#v", number, result.Findings)
				}
			}
		})
	}
	t.Run("empty listed set ignores fenced obligations", func(t *testing.T) {
		t.Parallel()
		root := writeProofFixture(t, "", "", otherProofs)
		prd := proofConstraints("") + "## Success Metrics\n\nNone. No outcome is measured.\n"
		writeCitationFixtureFile(t, root, "docs/specs/proof/_prd.md", "```md\n## Project Constraints\n\n- Active ADR obligations: ADR-0004.\n```\n\n"+prd)
		requireProofCount(t, checkProofFixture(t, root, speccheck.StageTechSpec), proofADRCode, 0)
	})
}

func TestCoverageCrossReadGapFollowsTheDeviationsLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		section string
		count   int
	}{
		{"missing map", "", 1},
		{"missing line", "## Coverage Map\n\n- User Story 1 → detector.\n", 1},
		{"empty line", "## Coverage Map\n\nDeviations from the PRD:\n", 1},
		{"prose remainder", "## Coverage Map\n\nDeviations from the PRD: The shape changes.\n", 1},
		{"unidentified deviation", "## Coverage Map\n\nDeviations from the PRD:\n\n- Aggregate per-request records.\n", 1},
		{"contract is not a PRD unit", "## Coverage Map\n\nDeviations from the PRD:\n\n- API Contract 1 changes.\n", 1},
		{"two unidentified deviations", "## Coverage Map\n\nDeviations from the PRD:\n\n- First change.\n- Second change.\n", 2},
		{"bullet before line", "## Coverage Map\n\n- User Story 1 → detector.\n\nDeviations from the PRD:\n", 1},
		{"line outside map", "## Notes\n\nDeviations from the PRD: None.\n", 1},
		{"shorter fence does not close example", "## Coverage Map\n\n````md\nExample.\n```\nDeviations from the PRD: None.\n````\n", 1},
		{"unclosed fence", "## Coverage Map\n\n```md\nDeviations from the PRD: None.\n", 1},
		{"fenced deviation", "## Coverage Map\n\nDeviations from the PRD:\n\n~~~md\n- User Story 1 changes.\n~~~\n", 1},
		{"none", "## Coverage Map\n\nDeviations from the PRD: None.\n", 0},
		{"named PRD units", "## Coverage Map\n\nDeviations from the PRD:\n\n- User Story 1 is aggregated.\n- Core Features 2-3 use the same detector.\n- Success Metric 1 uses fixture evidence.\n", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := writeProofFixture(t, "", "", "## Absence Evidence\n\nNone. No component is new.\n\n"+tt.section)
			result := checkProofFixture(t, root, speccheck.StageTechSpec)
			requireProofCount(t, result, proofCoverageCode, tt.count)
			for _, finding := range findingsWithCode(result, proofCoverageCode) {
				if !hasLocation(finding, "docs/specs/proof/_techspec.md") {
					t.Errorf("locations = %#v, want TechSpec line", finding.Where)
				}
			}
		})
	}
}

func TestTechSpecProofRunsAtTheTechSpecStageAndInTheFullSweep(t *testing.T) {
	t.Parallel()

	root := writeProofFixture(t, "applicable. ADR-0001 governs this fixture.", "", "")
	codes := []string{proofAbsenceCode, proofADRCode, proofCoverageCode}
	for _, stage := range []speccheck.Stage{speccheck.StageTechSpec, speccheck.StageAll, speccheck.StageTasks} {
		result := checkProofFixture(t, root, stage)
		for _, code := range codes {
			requireProofCount(t, result, code, 1)
			gap := speccheck.Result{Slug: result.Slug, Findings: append([]speccheck.Finding(nil), findingsWithCode(result, code)...)}
			if !strings.Contains(speccheck.RenderText(gap, speccheck.VerificationCoverage{}), "[gap] "+code) {
				t.Errorf("non-strict report = %q, want gap", speccheck.RenderText(gap, speccheck.VerificationCoverage{}))
			}
			speccheck.PromoteGaps(&gap)
			if gap.Findings[0].Severity != speccheck.SeverityError || !strings.Contains(speccheck.RenderText(gap, speccheck.VerificationCoverage{}), "[error] "+code) {
				t.Errorf("strict %s = %#v, want error", code, gap.Findings)
			}
		}
	}
	prd := checkProofFixture(t, root, speccheck.StagePRD)
	for _, code := range codes {
		requireProofCount(t, prd, code, 0)
		if !hasSkip(prd, code, "stage prd") {
			t.Errorf("PRD skips = %#v, want %s", prd.Skipped, code)
		}
	}
	path := filepath.Join(root, "docs", "specs", "proof", "_techspec.md")
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	for _, stage := range []speccheck.Stage{speccheck.StageTechSpec, speccheck.StageAll, speccheck.StageTasks} {
		result := checkProofFixture(t, root, stage)
		for _, code := range codes {
			requireProofCount(t, result, code, 0)
			if !hasSkip(result, code, "_techspec.md") {
				t.Errorf("%s skips = %#v, want %s absent TechSpec", stage, result.Skipped, code)
			}
		}
	}
	writeCitationFixtureFile(t, root, "docs/specs/proof/_techspec.md", proofConstraints("")+"## Absence Evidence\n\nNone. No component is new.\n\n## ADR Alternatives\n\n- ADR-0001: No rejected alternative is taken.\n\n## Coverage Map\n\nDeviations from the PRD: None.\n")
	for _, stage := range []speccheck.Stage{speccheck.StageTechSpec, speccheck.StageAll, speccheck.StageTasks} {
		result := checkProofFixture(t, root, stage)
		for _, code := range codes {
			requireProofCount(t, result, code, 0)
		}
	}
	if err := os.Remove(filepath.Join(root, "docs", "specs", "proof", "_prd.md")); err != nil {
		t.Fatal(err)
	}
	for _, stage := range []speccheck.Stage{speccheck.StageTechSpec, speccheck.StageAll} {
		result := checkProofFixture(t, root, stage)
		for _, code := range codes {
			requireProofCount(t, result, code, 0)
			if !hasSkip(result, code, "_prd.md") {
				t.Errorf("%s skips = %#v, want %s absent PRD", stage, result.Skipped, code)
			}
		}
	}
}

func TestTechSpecProofFixesNameTheSectionAndTheSkillRule(t *testing.T) {
	t.Parallel()

	root := writeProofFixture(t, "applicable. ADR-0001 governs this fixture.", "", "")
	tests := []struct {
		code, section, rule string
	}{
		{proofAbsenceCode, "## Absence Evidence", "Absence evidence"},
		{proofADRCode, "## ADR Alternatives", "Rejected alternatives"},
		{proofCoverageCode, "Deviations from the PRD:", "Coverage cross-read"},
	}
	for _, tt := range tests {
		finding := requireFinding(t, checkProofFixture(t, root, speccheck.StageTechSpec), tt.code)
		for _, text := range []string{tt.section, "write-techspec", tt.rule} {
			if !strings.Contains(finding.Fix, text) {
				t.Errorf("%s fix = %q, want %q", tt.code, finding.Fix, text)
			}
		}
	}
	writeCitationFixtureFile(t, root, "docs/specs/proof/_techspec.md", proofConstraints("")+"## Absence Evidence\n\nNone. No component is new.\n\n## ADR Alternatives\n\n- ADR-0001: No rejected alternative is taken.\n\n## Coverage Map\n\nDeviations from the PRD: None.\n")
	for _, tt := range tests {
		requireProofCount(t, checkProofFixture(t, root, speccheck.StageTechSpec), tt.code, 0)
	}
}

func proofConstraints(adrs string) string {
	if adrs == "" {
		adrs = "not applicable. No ADR is cited."
	}
	return "# Structural proof fixture\n\n## Project Constraints\n\n" +
		"- Identifier strategy: not applicable. No identifier changes. Source: `docs/agents/domain.md`.\n" +
		"- Authentication and HTTP: not applicable. No transport changes. Source: `docs/agents/agent-instructions.md`.\n" +
		"- Active ADR obligations: " + adrs + " Source: `docs/agents/domain.md`.\n" +
		"- Tooling authority: not applicable. No tooling changes. Source: `docs/agents/agent-instructions.md`.\n\n"
}

func writeProofFixture(t *testing.T, prdADRs, techADRs, sections string) string {
	t.Helper()
	root := t.TempDir()
	writeCitationFixtureFile(t, root, "docs/agents/domain.md", "# Fixture domain guide\n")
	writeCitationFixtureFile(t, root, "docs/agents/agent-instructions.md", "# Fixture agent guide\n")
	writeCitationFixtureFile(t, root, "docs/specs/proof/_prd.md", proofConstraints(prdADRs)+"## Success Metrics\n\nNone. No outcome is measured.\n")
	writeCitationFixtureFile(t, root, "docs/specs/proof/_techspec.md", proofConstraints(techADRs)+"## API Contracts\n\nNone. No public interface changes.\n\n"+sections)
	return root
}

func checkProofFixture(t *testing.T, root string, stage speccheck.Stage) speccheck.Result {
	t.Helper()
	result, err := speccheck.CheckStage(filepath.Join(root, "docs", "specs"), root, "proof", stage)
	if err != nil {
		t.Fatalf("CheckStage(%q): %v", stage, err)
	}
	return result
}

func requireProofCount(t *testing.T, result speccheck.Result, code string, count int) {
	t.Helper()
	findings := findingsWithCode(result, code)
	if len(findings) != count {
		t.Fatalf("%s findings = %#v, want %d", code, findings, count)
	}
	for _, finding := range findings {
		if finding.Severity != speccheck.SeverityGap {
			t.Errorf("%s severity = %q, want gap", code, finding.Severity)
		}
	}
}
