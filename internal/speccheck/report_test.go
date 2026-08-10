// Suite: Spec Consistency Check constraint reports
// Invariant: each declared constraint contradiction identifies both written sides and its concrete repair.
// Boundary IN: public speccheck API, Markdown fixtures, and text/JSON rendering
// Boundary OUT: CLI dispatch and detectors assigned to later Spec Tasks
package speccheck_test

import (
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

const fixtureSpecRoot = "testdata/repo/docs/specs"

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
			name:     "applicable tooling has no bounded files",
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

	textReport := speccheck.RenderText(result)
	for _, fragment := range []string{finding.Code, string(finding.Severity), finding.Summary, finding.Fix} {
		if !strings.Contains(textReport, fragment) {
			t.Errorf("text report does not contain %q:\n%s", fragment, textReport)
		}
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
