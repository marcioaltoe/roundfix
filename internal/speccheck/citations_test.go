// Suite: Spec citation, coverage, and reference consistency
// Invariant: declared ADR, PRD coverage, Task, and path references report only their written inconsistencies.
// Boundary IN: public speccheck API, fixture Markdown, accepted ADR corpus, and internal/spec loader
// Boundary OUT: CLI exit-code policy and characterization corpus assigned to later Spec Tasks
package speccheck_test

import (
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

func TestCheckADRUnlisted(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "citation-dirty")
	findings := findingsWithCode(result, speccheck.CodeADRUnlisted)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want exactly one", speccheck.CodeADRUnlisted, findings)
	}
	finding := findings[0]
	if finding.Severity != speccheck.SeverityError {
		t.Fatalf("severity = %q, want %q", finding.Severity, speccheck.SeverityError)
	}
	for _, path := range []string{
		"docs/specs/citation-dirty/_prd.md",
		"docs/specs/citation-dirty/_techspec.md",
	} {
		if !hasLocation(finding, path) {
			t.Errorf("locations = %#v, want %q", finding.Where, path)
		}
	}
}

func TestCheckADRClosureDepthOne(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "citation-dirty")
	findings := findingsWithCode(result, speccheck.CodeADRRelated)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want related ADR once", speccheck.CodeADRRelated, findings)
	}
	if findings[0].Severity != speccheck.SeverityGap {
		t.Fatalf("severity = %q, want %q", findings[0].Severity, speccheck.SeverityGap)
	}
	if !strings.Contains(findings[0].Summary, "ADR-0102") {
		t.Fatalf("summary = %q, want ADR-0102", findings[0].Summary)
	}
	for _, listed := range []string{"ADR-0100", "ADR-0101"} {
		for _, finding := range findings {
			if strings.HasPrefix(finding.Summary, listed+" ") {
				t.Errorf("listed %s reported as related: %#v", listed, finding)
			}
		}
	}
}

func TestCheckCoverageUnmapped(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "coverage-unmapped")
	findings := findingsWithCode(result, speccheck.CodeCoverageUnmapped)
	if len(findings) != 4 {
		t.Fatalf("%s findings = %#v, want one per four Core Features", speccheck.CodeCoverageUnmapped, findings)
	}
	for _, finding := range findings {
		if !strings.Contains(finding.Summary, "Core Feature") {
			t.Errorf("unexpected uncovered unit: %q", finding.Summary)
		}
	}
}

func TestCheckCoverageRange(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "coverage-range")
	if findings := findingsWithCode(result, speccheck.CodeCoverageUnmapped); len(findings) != 0 {
		t.Fatalf("range produced findings = %#v, want none", findings)
	}
	if findings := findingsWithCode(result, speccheck.CodeCoverageUntasked); len(findings) != 0 {
		t.Fatalf("Task reference range produced findings = %#v, want none", findings)
	}
}

func TestCheckCoverageUntasked(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "coverage-untasked")
	findings := findingsWithCode(result, speccheck.CodeCoverageUntasked)
	if len(findings) != 1 || !strings.Contains(findings[0].Summary, "Core Feature 4") {
		t.Fatalf("%s findings = %#v, want only Core Feature 4", speccheck.CodeCoverageUntasked, findings)
	}
}

func TestCheckReferenceUnresolved(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "reference-unresolved")
	finding := requireFinding(t, result, speccheck.CodeReferenceUnresolved)
	const missingPath = "missing/guide.md"
	if !strings.Contains(finding.Summary, missingPath) {
		t.Fatalf("summary = %q, want missing path %q", finding.Summary, missingPath)
	}
	if !hasLocation(finding, "docs/specs/reference-unresolved/task_01.md") {
		t.Errorf("locations = %#v, want declaring Task line", finding.Where)
	}
	if !hasLocation(finding, missingPath) {
		t.Errorf("locations = %#v, want unresolved path", finding.Where)
	}
}

func TestCheckReferenceIndexUnresolved(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "reference-index-unresolved")
	finding := requireFinding(t, result, speccheck.CodeReferenceUnresolved)
	if !strings.Contains(finding.Summary, "missing-note.md") {
		t.Fatalf("summary = %q, want missing-note.md", finding.Summary)
	}
	if !hasLocation(finding, "docs/specs/reference-index-unresolved/references/_index.md") {
		t.Errorf("locations = %#v, want declaring reference index line", finding.Where)
	}
}

func TestCheckCoverageUntaskedSkippedWithoutTaskGraph(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "no-taskgraph")
	if findings := findingsWithCode(result, speccheck.CodeCoverageUntasked); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want none", speccheck.CodeCoverageUntasked, findings)
	}
	if !hasSkip(result, speccheck.CodeCoverageUntasked, "_tasks.md") {
		t.Fatalf("Skipped = %#v, want %s missing _tasks.md", result.Skipped, speccheck.CodeCoverageUntasked)
	}
}

func TestCheckCoverageUnmappedSkippedWithoutTechSpec(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "no-techspec")
	if findings := findingsWithCode(result, speccheck.CodeCoverageUnmapped); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want none", speccheck.CodeCoverageUnmapped, findings)
	}
	if !hasSkip(result, speccheck.CodeCoverageUnmapped, "_techspec.md") {
		t.Fatalf("Skipped = %#v, want %s missing _techspec.md", result.Skipped, speccheck.CodeCoverageUnmapped)
	}
}

func TestCheckADRCorpusAbsentSkipsADRDetectors(t *testing.T) {
	t.Parallel()

	repoRoot, err := filepath.Abs("testdata/no-adr-repo")
	if err != nil {
		t.Fatalf("resolve fixture repository: %v", err)
	}
	result, err := speccheck.Check(filepath.Join(repoRoot, "docs", "specs"), repoRoot, "no-adr")
	if err != nil {
		t.Fatalf("Check(no-adr) error = %v", err)
	}
	for _, code := range []string{speccheck.CodeADRUnlisted, speccheck.CodeADRRelated} {
		if !hasSkip(result, code, "docs/adr") {
			t.Errorf("Skipped = %#v, want %s missing docs/adr", result.Skipped, code)
		}
	}
}

func TestCheckCitationCoverageErrorLocations(t *testing.T) {
	t.Parallel()

	for _, slug := range []string{
		"citation-dirty",
		"coverage-unmapped",
		"coverage-untasked",
		"reference-unresolved",
		"reference-index-unresolved",
	} {
		slug := slug
		t.Run(slug, func(t *testing.T) {
			t.Parallel()

			result := checkFixture(t, slug)
			for _, finding := range result.Findings {
				if finding.Code == speccheck.CodeADRRelated {
					continue
				}
				if finding.Severity != speccheck.SeverityError {
					t.Errorf("%s severity = %q, want %q", finding.Code, finding.Severity, speccheck.SeverityError)
				}
				if len(finding.Where) < 2 {
					t.Errorf("%s locations = %#v, want both sides", finding.Code, finding.Where)
				}
			}
		})
	}
}

func findingsWithCode(result speccheck.Result, code string) []speccheck.Finding {
	var findings []speccheck.Finding
	for _, finding := range result.Findings {
		if finding.Code == code {
			findings = append(findings, finding)
		}
	}
	return findings
}
