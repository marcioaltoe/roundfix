// Suite: Spec Consistency Check constraint reports
// Invariant: each declared constraint contradiction identifies both written sides and its concrete repair.
// Boundary IN: public speccheck API, Markdown fixtures, and text/JSON rendering
// Boundary OUT: CLI dispatch and detectors assigned to later Spec Tasks
package speccheck_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

const fixtureSpecRoot = "testdata/repo/docs/specs"

func TestToolingRowStatesApplicability(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		row        string
		recordPath string
		record     string
		wantCode   string
	}{
		{
			name: "template no mutation wording",
			row:  "Tooling authority: applicable — no protected tooling mutation proposed or authorized.",
		},
		{
			name:     "declared mutation without bounded files",
			row:      "Tooling authority: applicable — protected tooling mutation proposed and authorized.",
			wantCode: speccheck.CodeToolingUnbounded,
		},
		{
			name:       "authorization record omits Spec",
			row:        "Tooling authority: applicable — protected tooling mutation authorized at `docs/workflow/authorizations/2026-08-11-other-spec.md`; bounded files: `Makefile`.",
			recordPath: "docs/workflow/authorizations/2026-08-11-other-spec.md",
			record:     "---\ngranted: 2026-08-11\naction: mutate protected tooling\npaths:\n  - Makefile\nconsuming: 0115-other-spec\n---\n",
			wantCode:   speccheck.CodeToolingUnauthorized,
		},
		{
			name:       "authorization record states grant only in prose",
			row:        "Tooling authority: applicable — protected tooling mutation authorized at `docs/workflow/authorizations/2026-08-11-prose-grant.md`; bounded files: `Makefile`.",
			recordPath: "docs/workflow/authorizations/2026-08-11-prose-grant.md",
			record:     "Authorization for Spec 0114-tooling-row permits changes to Makefile.\n",
			wantCode:   speccheck.CodeToolingUntyped,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoRoot, specsRoot, slug := writeToolingRowFixture(t, tt.row, tt.recordPath, tt.record)
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
			if findings := findingsWithCode(result, tt.wantCode); len(findings) != 1 {
				t.Fatalf("%s findings = %#v, want exactly one", tt.wantCode, findings)
			}
			if len(result.Findings) != 1 {
				t.Fatalf("StagePRD findings = %#v, want only %s", result.Findings, tt.wantCode)
			}
		})
	}
}

func writeToolingRowFixture(t *testing.T, toolingRow, recordPath, record string) (string, string, string) {
	t.Helper()

	const slug = "0114-tooling-row"
	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	writeToolingRowFile(t, repoRoot, "docs/agents/agent-instructions.md", "# Agent instructions\n")
	writeToolingRowFile(t, repoRoot, "docs/specs/"+slug+"/_prd.md", "# Tooling row\n\n## Project Constraints\n\n"+
		"- Identifier strategy: not applicable — no identifier change. Source: `docs/agents/agent-instructions.md`.\n"+
		"- Authentication and HTTP: not applicable — no network boundary. Source: `docs/agents/agent-instructions.md`.\n"+
		"- Active ADR obligations: not applicable — no ADR applies. Source: `docs/agents/agent-instructions.md`.\n"+
		"- "+toolingRow+" Source: `docs/agents/agent-instructions.md`.\n")
	if recordPath != "" {
		writeToolingRowFile(t, repoRoot, recordPath, record)
	}
	return repoRoot, specsRoot, slug
}

func writeToolingRowFile(t *testing.T, repoRoot, relativePath, content string) {
	t.Helper()

	path := filepath.Join(repoRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create directory for %s: %v", relativePath, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relativePath, err)
	}
}

func TestCheckConstraintAndTooling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		slug     string
		wantCode string
	}{
		{
			name:     "missing required PRD row",
			slug:     "constraint-missing",
			wantCode: speccheck.CodeConstraintMissing,
		},
		{
			name:     "applicability without reason",
			slug:     "constraint-unreasoned",
			wantCode: speccheck.CodeConstraintUnreasoned,
		},
		{
			name:     "cited source does not exist",
			slug:     "constraint-source",
			wantCode: speccheck.CodeConstraintSource,
		},
		{
			name:     "authorization record omits Spec",
			slug:     "tooling-unauthorized",
			wantCode: speccheck.CodeToolingUnauthorized,
		},
		{
			name:     "declared tooling mutation has no bounded files",
			slug:     "tooling-unbounded",
			wantCode: speccheck.CodeToolingUnbounded,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := checkFixture(t, tt.slug)
			finding := requireFinding(t, result, tt.wantCode)
			if finding.Severity != speccheck.SeverityError {
				t.Fatalf("finding %s severity = %q, want %q", tt.wantCode, finding.Severity, speccheck.SeverityError)
			}
			if strings.Contains(finding.Summary, "\n") || strings.TrimSpace(finding.Summary) == "" {
				t.Fatalf("finding %s summary must be one non-empty line: %q", tt.wantCode, finding.Summary)
			}
			if strings.TrimSpace(finding.Fix) == "" {
				t.Fatalf("finding %s has no concrete fix", tt.wantCode)
			}
		})
	}
}

func TestCheckConstraintSourceNamesMissingPath(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "constraint-source")
	finding := requireFinding(t, result, speccheck.CodeConstraintSource)
	const missingPath = "docs/agents/missing-guide.md"
	if !strings.Contains(finding.Summary, missingPath) {
		t.Fatalf("summary = %q, want missing path %q", finding.Summary, missingPath)
	}
	if !hasLocation(finding, missingPath) {
		t.Fatalf("locations = %#v, want missing path %q", finding.Where, missingPath)
	}
}

func TestCheckToolingUnauthorizedLocatesSpecAndRecord(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "tooling-unauthorized")
	finding := requireFinding(t, result, speccheck.CodeToolingUnauthorized)
	for _, path := range []string{
		"docs/specs/tooling-unauthorized/_prd.md",
		"docs/workflow/authorizations/other-spec.md",
	} {
		if !hasLocation(finding, path) {
			t.Errorf("locations = %#v, want %q", finding.Where, path)
		}
	}
}

func TestCheckToolingUntypedReportsOnlyTheProseRecord(t *testing.T) {
	t.Parallel()

	untyped := checkFixture(t, "tooling-untyped")
	finding := requireFinding(t, untyped, speccheck.CodeToolingUntyped)
	const record = "docs/workflow/authorizations/2026-08-11-untyped-grant.md"
	if !hasLocation(finding, record) {
		t.Fatalf("locations = %#v, want %q", finding.Where, record)
	}

	typed := checkFixture(t, "tooling-typed")
	for _, found := range typed.Findings {
		if found.Code == speccheck.CodeToolingUntyped {
			t.Fatalf("typed authorization record reported %#v", found)
		}
	}
}

func TestCheckSkipMissingTechSpec(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "no-techspec")
	if len(result.Findings) != 0 {
		t.Fatalf("Findings = %#v, want none", result.Findings)
	}
	for _, code := range []string{
		speccheck.CodeConstraintMissing,
		speccheck.CodeConstraintUnreasoned,
		speccheck.CodeConstraintSource,
		speccheck.CodeToolingUnauthorized,
		speccheck.CodeToolingUnbounded,
	} {
		if !hasSkip(result, code, "_techspec.md") {
			t.Errorf("Skipped = %#v, want %s missing _techspec.md", result.Skipped, code)
		}
	}
}

func TestCheckErrorLocations(t *testing.T) {
	t.Parallel()

	for _, slug := range []string{
		"constraint-missing",
		"constraint-unreasoned",
		"constraint-source",
		"tooling-unauthorized",
		"tooling-unbounded",
	} {
		slug := slug
		t.Run(slug, func(t *testing.T) {
			t.Parallel()

			result := checkFixture(t, slug)
			for _, finding := range result.Findings {
				if finding.Severity != speccheck.SeverityError {
					continue
				}
				if len(finding.Where) < 2 {
					t.Errorf("%s locations = %#v, want at least two", finding.Code, finding.Where)
				}
				for _, location := range finding.Where {
					if filepath.IsAbs(location.Path) || location.Path == "" {
						t.Errorf("%s location path = %q, want repository-relative", finding.Code, location.Path)
					}
					if location.Line < 1 {
						t.Errorf("%s location line = %d, want 1-based", finding.Code, location.Line)
					}
				}
			}
		})
	}
}

func TestCheckCleanFixture(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "clean")
	if len(result.Findings) != 0 {
		t.Fatalf("Findings = %#v, want none", result.Findings)
	}
}

func TestRenderResultTextAndJSON(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "tooling-unauthorized")
	finding := requireFinding(t, result, speccheck.CodeToolingUnauthorized)

	textReport := speccheck.RenderText(result, speccheck.VerificationCoverage{Ran: true, Commands: 2})
	for _, fragment := range []string{finding.Code, string(finding.Severity), finding.Summary, finding.Fix} {
		if !strings.Contains(textReport, fragment) {
			t.Errorf("text report does not contain %q:\n%s", fragment, textReport)
		}
	}
	if strings.Contains(textReport, "Authored Verification commands") {
		t.Fatalf("finding report contains clean-verdict coverage:\n%s", textReport)
	}
	for _, location := range finding.Where {
		want := location.Path + ":" + strconv.Itoa(location.Line)
		if !strings.Contains(textReport, want) {
			t.Errorf("text report does not contain location %q:\n%s", want, textReport)
		}
	}

	jsonReport, err := speccheck.RenderJSON(result)
	if err != nil {
		t.Fatalf("RenderJSON() error = %v", err)
	}
	var document struct {
		Schema   string                      `json:"schema"`
		Slug     string                      `json:"slug"`
		Findings []speccheck.Finding         `json:"findings"`
		Skipped  []speccheck.SkippedDetector `json:"skipped"`
	}
	if err := json.Unmarshal(jsonReport, &document); err != nil {
		t.Fatalf("RenderJSON() returned invalid JSON: %v\n%s", err, jsonReport)
	}
	if document.Schema != speccheck.SchemaVersion {
		t.Errorf("schema = %q, want %q", document.Schema, speccheck.SchemaVersion)
	}
	if document.Slug != result.Slug || len(document.Findings) != len(result.Findings) {
		t.Errorf("JSON document = %#v, want slug %q and %d findings", document, result.Slug, len(result.Findings))
	}
}

func TestVerdictLineStatesProbeCoverage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		coverage speccheck.VerificationCoverage
		want     string
	}{
		{
			name: "probe ran",
			coverage: speccheck.VerificationCoverage{
				Ran:      true,
				Commands: 2,
			},
			want: "No findings. Authored Verification commands executed: 2.",
		},
		{
			name:     "probe did not run",
			coverage: speccheck.VerificationCoverage{},
			want:     "No findings. Authored Verification commands were not executed.",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			report := speccheck.RenderText(speccheck.Result{Slug: "clean"}, tt.coverage)
			lines := strings.Split(report, "\n")
			if len(lines) < 2 {
				t.Fatalf("RenderText() = %q, want a verdict line", report)
			}
			if lines[1] != tt.want {
				t.Fatalf("verdict line = %q, want %q", lines[1], tt.want)
			}
		})
	}
}

func TestRenderVocabularySkipDoesNotLookLikeFinding(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "vocabulary-none")
	if !hasSkip(result, speccheck.CodeVocabularyUndocumented, "Vocabulary Contract") {
		t.Fatalf("Skipped = %#v, want vocabulary detector skip", result.Skipped)
	}

	textReport := speccheck.RenderText(result, speccheck.VerificationCoverage{})
	if strings.Contains(textReport, speccheck.CodeVocabularyUndocumented) {
		t.Fatalf("text skip masquerades as a %s finding:\n%s", speccheck.CodeVocabularyUndocumented, textReport)
	}
	for _, fragment := range []string{"Skipped:", "vocabulary documentation detector", "Vocabulary Contract"} {
		if !strings.Contains(textReport, fragment) {
			t.Errorf("text skip does not contain %q:\n%s", fragment, textReport)
		}
	}

	jsonReport, err := speccheck.RenderJSON(result)
	if err != nil {
		t.Fatalf("RenderJSON() error = %v", err)
	}
	if !strings.Contains(string(jsonReport), `"code":"`+speccheck.CodeVocabularyUndocumented+`"`) {
		t.Fatalf("JSON skip lost stable detector code: %s", jsonReport)
	}

	findingReport := speccheck.RenderText(checkFixture(t, "vocabulary-missing"), speccheck.VerificationCoverage{})
	if !strings.Contains(findingReport, "[error] "+speccheck.CodeVocabularyUndocumented+":") {
		t.Fatalf("vocabulary finding lost its diagnostic code:\n%s", findingReport)
	}
}

func checkFixture(t *testing.T, slug string) speccheck.Result {
	t.Helper()

	repoRoot, err := filepath.Abs("testdata/repo")
	if err != nil {
		t.Fatalf("resolve fixture repository: %v", err)
	}
	result, err := speccheck.Check(filepath.Join(repoRoot, "docs", "specs"), repoRoot, slug)
	if err != nil {
		t.Fatalf("Check(%q) error = %v", slug, err)
	}
	return result
}

func requireFinding(t *testing.T, result speccheck.Result, code string) speccheck.Finding {
	t.Helper()

	for _, finding := range result.Findings {
		if finding.Code == code {
			return finding
		}
	}
	t.Fatalf("Findings = %#v, want code %s", result.Findings, code)
	return speccheck.Finding{}
}

func hasLocation(finding speccheck.Finding, path string) bool {
	for _, location := range finding.Where {
		if location.Path == path && location.Line > 0 {
			return true
		}
	}
	return false
}

func hasSkip(result speccheck.Result, code string, missing string) bool {
	for _, skipped := range result.Skipped {
		if skipped.Code == code && strings.Contains(skipped.Missing, missing) {
			return true
		}
	}
	return false
}
