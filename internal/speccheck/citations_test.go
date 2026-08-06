// Suite: Spec citation, coverage, reference, and loop-order consistency
// Invariant: declared ADR, coverage, path, and loop-order sources report only their written inconsistencies.
// Boundary IN: public speccheck API, fixture artifacts, accepted ADR corpus, and internal/spec loader
// Boundary OUT: CLI exit-code policy and characterization corpus assigned to later Spec Tasks
package speccheck_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

const (
	loopOrderShippedClausePath   = "internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md"
	loopOrderRepositoryGuidePath = "docs/agents/autonomous-work.md"
	loopOrderBaselineModulePath  = "internal/baseline/assets/modules/autonomous-work.json"
	loopOrderCarrierSlug         = "loop-order"
)

func TestCheckLoopOrderDivergent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      string
		current   string
		divergent string
	}{
		{
			name:      "shipped clause",
			path:      loopOrderShippedClausePath,
			current:   "archive, open the Pull Request, watch until Clean, and merge",
			divergent: "merge, open the Pull Request, watch until Clean, and archive",
		},
		{
			name:      "repository guide",
			path:      loopOrderRepositoryGuidePath,
			current:   "archive, open the Pull\nRequest, watch until Clean, and merge",
			divergent: "merge, open the Pull\nRequest, watch until Clean, and archive",
		},
		{
			name:      "Baseline module asset",
			path:      loopOrderBaselineModulePath,
			current:   "archive, open the Pull Request, watch until Clean, and merge",
			divergent: "merge, open the Pull Request, watch until Clean, and archive",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoRoot := writeLoopOrderFixture(t)
			path := filepath.Join(repoRoot, filepath.FromSlash(tt.path))
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s fixture: %v", tt.name, err)
			}
			if !strings.Contains(string(content), tt.current) {
				t.Fatalf("%s fixture does not contain current order %q", tt.name, tt.current)
			}
			content = []byte(strings.Replace(string(content), tt.current, tt.divergent, 1))
			if err := os.WriteFile(path, content, 0o644); err != nil {
				t.Fatalf("write divergent %s fixture: %v", tt.name, err)
			}

			result, err := speccheck.Check(filepath.Join(repoRoot, "docs", "specs"), repoRoot, loopOrderCarrierSlug)
			if err != nil {
				t.Fatalf("Check(loop-order) error = %v", err)
			}
			finding := requireFinding(t, result, speccheck.CodeLoopOrderDivergent)
			if finding.Severity != speccheck.SeverityError {
				t.Fatalf("severity = %q, want %q", finding.Severity, speccheck.SeverityError)
			}
			for _, sourceLabel := range []string{"shipped clause", "repository guide", "Baseline module asset"} {
				if !strings.Contains(finding.Summary, sourceLabel) {
					t.Errorf("summary = %q, want source %q", finding.Summary, sourceLabel)
				}
			}
			if !strings.Contains(finding.Summary, strings.ReplaceAll(tt.divergent, "\n", " ")) {
				t.Errorf("summary = %q, want divergent order %q", finding.Summary, tt.divergent)
			}
			for _, sourcePath := range []string{
				loopOrderShippedClausePath,
				loopOrderRepositoryGuidePath,
				loopOrderBaselineModulePath,
			} {
				if !hasLocation(finding, sourcePath) {
					t.Errorf("locations = %#v, want source %q", finding.Where, sourcePath)
				}
			}
		})
	}
}

func TestCheckLoopOrderRepositoryAgrees(t *testing.T) {
	t.Parallel()

	repoRoot := characterizationRepositoryRoot(t)
	specsRoot := writeLoopOrderCarrier(t, t.TempDir())
	// SC-LOOP-ORDER-DIVERGENT reads the repository's own order statements, so
	// a dedicated carrier keeps the check independent from active and archived
	// repository Specs.
	result, err := speccheck.Check(specsRoot, repoRoot, loopOrderCarrierSlug)
	if err != nil {
		t.Fatalf("Check(repository) error = %v", err)
	}
	if findings := findingsWithCode(result, speccheck.CodeLoopOrderDivergent); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want corrected repository sources to agree", speccheck.CodeLoopOrderDivergent, findings)
	}
}

func writeLoopOrderFixture(t *testing.T) string {
	t.Helper()

	sourceRoot := characterizationRepositoryRoot(t)
	targetRoot := t.TempDir()
	for _, relative := range []string{
		loopOrderShippedClausePath,
		loopOrderRepositoryGuidePath,
		loopOrderBaselineModulePath,
	} {
		sourcePath := filepath.Join(sourceRoot, filepath.FromSlash(relative))
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			t.Fatalf("read loop-order source %q: %v", sourcePath, err)
		}
		targetPath := filepath.Join(targetRoot, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			t.Fatalf("create loop-order fixture directory: %v", err)
		}
		if err := os.WriteFile(targetPath, content, 0o644); err != nil {
			t.Fatalf("write loop-order fixture %q: %v", targetPath, err)
		}
	}

	writeLoopOrderCarrier(t, targetRoot)
	return targetRoot
}

func writeLoopOrderCarrier(t *testing.T, root string) string {
	t.Helper()
	specsRoot := filepath.Join(root, "docs", "specs")
	specDir := filepath.Join(specsRoot, loopOrderCarrierSlug)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("create fixture Spec: %v", err)
	}
	const prd = "---\nspec: " + loopOrderCarrierSlug + "\nstatus: active\n---\n\n# Loop order\n"
	if err := os.WriteFile(filepath.Join(specDir, "_prd.md"), []byte(prd), 0o644); err != nil {
		t.Fatalf("write fixture PRD: %v", err)
	}
	return specsRoot
}

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
