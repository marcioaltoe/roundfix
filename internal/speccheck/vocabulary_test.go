// Suite: Spec emitted-vocabulary consistency
// Invariant: every distinct token selected by a declared RE2 pattern appears in the named documentation or produces one located diagnostic.
// Boundary IN: public speccheck API, Markdown contract declarations, and repository fixture files
// Boundary OUT: CLI exit-code policy and characterization corpus assigned to later Spec Tasks
package speccheck_test

import (
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

func TestCheckVocabularyUndocumented(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "vocabulary-missing")
	findings := findingsWithCode(result, speccheck.CodeVocabularyUndocumented)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want exactly one", speccheck.CodeVocabularyUndocumented, findings)
	}
	finding := findings[0]
	if finding.Severity != speccheck.SeverityError {
		t.Fatalf("severity = %q, want %q", finding.Severity, speccheck.SeverityError)
	}
	if !strings.Contains(finding.Summary, "publish:") {
		t.Fatalf("summary = %q, want undocumented token publish:", finding.Summary)
	}
	if !hasExactLocation(finding, "docs/specs/vocabulary-missing/emitter.go", 8) {
		t.Errorf("locations = %#v, want first publish: emission at line 8", finding.Where)
	}
	for _, path := range []string{
		"docs/specs/vocabulary-missing/emitter.go",
		"docs/specs/vocabulary-missing/runbook.md",
	} {
		if !hasLocation(finding, path) {
			t.Errorf("locations = %#v, want %q", finding.Where, path)
		}
	}
}

func TestCheckVocabularySatisfied(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "vocabulary-satisfied")
	if findings := findingsWithCode(result, speccheck.CodeVocabularyUndocumented); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want none", speccheck.CodeVocabularyUndocumented, findings)
	}
}

func TestCheckVocabularyDeduplicatesRepeatedToken(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "vocabulary-missing")
	findings := findingsWithCode(result, speccheck.CodeVocabularyUndocumented)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want repeated publish: token once", speccheck.CodeVocabularyUndocumented, findings)
	}
}

func TestCheckVocabularyInvalidPatternReturnsFinding(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "vocabulary-invalid-pattern")
	finding := requireFinding(t, result, speccheck.CodeVocabularyUndocumented)
	if !strings.Contains(finding.Summary, "invalid RE2 pattern") {
		t.Fatalf("summary = %q, want invalid RE2 pattern", finding.Summary)
	}
	if !hasExactLocation(finding, "docs/specs/vocabulary-invalid-pattern/_techspec.md", 13) {
		t.Fatalf("locations = %#v, want pattern declaration line 13", finding.Where)
	}
}

func TestCheckVocabularyUnreadablePathsReturnFindings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		slug     string
		wantPath string
	}{
		{
			name:     "emitting path",
			slug:     "vocabulary-missing-emitter",
			wantPath: "docs/specs/vocabulary-missing-emitter/missing.go",
		},
		{
			name:     "documenting path",
			slug:     "vocabulary-missing-documentation",
			wantPath: "docs/specs/vocabulary-missing-documentation/missing.md",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := checkFixture(t, tt.slug)
			finding := requireFinding(t, result, speccheck.CodeVocabularyUndocumented)
			if finding.Severity != speccheck.SeverityError {
				t.Fatalf("severity = %q, want %q", finding.Severity, speccheck.SeverityError)
			}
			if !strings.Contains(finding.Summary, "unreadable") {
				t.Fatalf("summary = %q, want unreadable path", finding.Summary)
			}
			if !hasLocation(finding, "docs/specs/"+tt.slug+"/_techspec.md") {
				t.Fatalf("locations = %#v, want TechSpec declaration line", finding.Where)
			}
			if !hasLocation(finding, tt.wantPath) {
				t.Fatalf("locations = %#v, want %q", finding.Where, tt.wantPath)
			}
		})
	}
}

func TestCheckVocabularyAbsentContractIsSkipped(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "vocabulary-none")
	if findings := findingsWithCode(result, speccheck.CodeVocabularyUndocumented); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want none", speccheck.CodeVocabularyUndocumented, findings)
	}
	if !hasSkip(result, speccheck.CodeVocabularyUndocumented, "Vocabulary Contract") {
		t.Fatalf("Skipped = %#v, want vocabulary detector skip", result.Skipped)
	}
}

func TestCheckVocabularyFixTeachesDeclarationShape(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "vocabulary-missing")
	finding := requireFinding(t, result, speccheck.CodeVocabularyUndocumented)
	const declarationShape = "## Vocabulary Contract\n\n" +
		"- emits: `<repository-relative path>`\n" +
		"  pattern: `<RE2>`\n" +
		"  documented-in: `<repository-relative path>`"
	if !strings.Contains(finding.Fix, declarationShape) {
		t.Fatalf("fix = %q, want exact declaration shape %q", finding.Fix, declarationShape)
	}
	for _, field := range []string{"emits:", "pattern:", "documented-in:"} {
		if !strings.Contains(finding.Fix, field) {
			t.Errorf("fix = %q, want field %q", finding.Fix, field)
		}
	}
}

func hasExactLocation(finding speccheck.Finding, path string, line int) bool {
	for _, location := range finding.Where {
		if location.Path == path && location.Line == line {
			return true
		}
	}
	return false
}
