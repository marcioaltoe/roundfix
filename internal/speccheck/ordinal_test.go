// Suite: claimed ADR ordinal consistency
// Invariant: one ADR ordinal identifies one path across the repository tree and active Spec claims.
// Boundary IN: public stage-scoped Spec checks over real Markdown fixtures
// Boundary OUT: task authoring guidance and Daemon Verification execution
package speccheck_test

import (
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

func TestOrdinalClaimedOnTree(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	writeOrdinalSpec(t, repoRoot, "claiming-spec", "docs/adr/0042-new-decision.md")
	writeCitationFixtureFile(t, repoRoot, "docs/adr/0042-existing-decision.md", "# Existing decision\n")

	result, err := speccheck.CheckStage(specsRoot, repoRoot, "claiming-spec", speccheck.StageTasks)
	if err != nil {
		t.Fatalf("CheckStage(StageTasks): %v", err)
	}
	findings := findingsWithCode(result, speccheck.CodeOrdinalClaimed)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want one", speccheck.CodeOrdinalClaimed, findings)
	}
	finding := findings[0]
	if finding.Severity != speccheck.SeverityError {
		t.Errorf("severity = %q, want %q", finding.Severity, speccheck.SeverityError)
	}
	for _, want := range []string{"0042", "docs/adr/0042-new-decision.md", "docs/adr/0042-existing-decision.md"} {
		if !strings.Contains(finding.Summary, want) {
			t.Errorf("summary = %q, want %q", finding.Summary, want)
		}
	}
}

func TestOrdinalClaimedByAnotherActiveSpec(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	writeOrdinalSpec(t, repoRoot, "first-spec", "docs/adr/0042-first-decision.md")
	writeOrdinalSpec(t, repoRoot, "second-spec", "docs/adr/0042-second-decision.md")

	result, err := speccheck.CheckStage(specsRoot, repoRoot, "first-spec", speccheck.StageTasks)
	if err != nil {
		t.Fatalf("CheckStage(StageTasks): %v", err)
	}
	findings := findingsWithCode(result, speccheck.CodeOrdinalClaimed)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want one", speccheck.CodeOrdinalClaimed, findings)
	}
	for _, want := range []string{"0042", "first-spec", "second-spec"} {
		if !strings.Contains(findings[0].Summary, want) {
			t.Errorf("summary = %q, want %q", findings[0].Summary, want)
		}
	}
}

func TestDistinctOrdinalsAreAccepted(t *testing.T) {
	t.Parallel()

	t.Run("distinct active Spec claims", func(t *testing.T) {
		repoRoot := t.TempDir()
		specsRoot := filepath.Join(repoRoot, "docs", "specs")
		writeOrdinalSpec(t, repoRoot, "first-spec", "docs/adr/0042-first-decision.md")
		writeOrdinalSpec(t, repoRoot, "second-spec", "docs/adr/0043-second-decision.md")

		result, err := speccheck.CheckStage(specsRoot, repoRoot, "first-spec", speccheck.StageTasks)
		if err != nil {
			t.Fatalf("CheckStage(StageTasks): %v", err)
		}
		if findings := findingsWithCode(result, speccheck.CodeOrdinalClaimed); len(findings) != 0 {
			t.Fatalf("%s findings = %#v, want none", speccheck.CodeOrdinalClaimed, findings)
		}
	})

	t.Run("claim fulfilled by exact tree path", func(t *testing.T) {
		repoRoot := t.TempDir()
		specsRoot := filepath.Join(repoRoot, "docs", "specs")
		const claim = "docs/adr/0042-existing-decision.md"
		writeOrdinalSpec(t, repoRoot, "fulfilled-spec", claim)
		writeCitationFixtureFile(t, repoRoot, claim, "# Existing decision\n")

		result, err := speccheck.CheckStage(specsRoot, repoRoot, "fulfilled-spec", speccheck.StageTasks)
		if err != nil {
			t.Fatalf("CheckStage(StageTasks): %v", err)
		}
		if findings := findingsWithCode(result, speccheck.CodeOrdinalClaimed); len(findings) != 0 {
			t.Fatalf("%s findings = %#v, want none", speccheck.CodeOrdinalClaimed, findings)
		}
	})
}

func writeOrdinalSpec(t *testing.T, repoRoot, slug, claim string) {
	t.Helper()

	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_prd.md", `---
spec: `+slug+`
status: active
created: 2026-09-24
surfaces: [backend]
---

# Ordinal claim fixture
`)
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_tasks.md", `---
schema: spec-tasks/v1
spec: `+slug+`
qa: declined
qa_reason: no behavioral surface in this fixture
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
---
`)
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/task_01.md", `---
task: task_01
spec: `+slug+`
status: pending
type: backend
---

# Task 01: Claim an ADR ordinal

## Context

- creates: `+"`"+claim+"`"+`

## Verification

- `+"`test -f "+claim+"`"+`
`)
}
