// Suite: claimed ADR ordinal consistency
// Invariant: one ADR ordinal identifies one path across the repository tree and active Spec claims.
// Boundary IN: public stage-scoped Spec checks over real Markdown fixtures
// Boundary OUT: task authoring guidance and Daemon Verification execution
package speccheck_test

import (
	"fmt"
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

func TestSameSpecDuplicateOrdinalIsClaimed(t *testing.T) {
	t.Parallel()

	t.Run("duplicate ordinal with different paths", func(t *testing.T) {
		repoRoot := t.TempDir()
		specsRoot := filepath.Join(repoRoot, "docs", "specs")
		writeOrdinalSpecClaims(t, repoRoot, "claiming-spec",
			"docs/adr/0042-first-decision.md",
			"docs/adr/0042-second-decision.md",
		)

		result, err := speccheck.CheckStage(specsRoot, repoRoot, "claiming-spec", speccheck.StageTasks)
		if err != nil {
			t.Fatalf("CheckStage(StageTasks): %v", err)
		}
		findings := findingsWithCode(result, speccheck.CodeOrdinalClaimed)
		if len(findings) != 1 {
			t.Fatalf("%s findings = %#v, want one", speccheck.CodeOrdinalClaimed, findings)
		}
		for _, want := range []string{"0042", "docs/adr/0042-first-decision.md", "docs/adr/0042-second-decision.md", "claiming-spec"} {
			if !strings.Contains(findings[0].Summary, want) {
				t.Errorf("summary = %q, want %q", findings[0].Summary, want)
			}
		}
	})

	t.Run("distinct ordinals", func(t *testing.T) {
		repoRoot := t.TempDir()
		specsRoot := filepath.Join(repoRoot, "docs", "specs")
		writeOrdinalSpecClaims(t, repoRoot, "claiming-spec",
			"docs/adr/0042-first-decision.md",
			"docs/adr/0043-second-decision.md",
		)

		result, err := speccheck.CheckStage(specsRoot, repoRoot, "claiming-spec", speccheck.StageTasks)
		if err != nil {
			t.Fatalf("CheckStage(StageTasks): %v", err)
		}
		if findings := findingsWithCode(result, speccheck.CodeOrdinalClaimed); len(findings) != 0 {
			t.Fatalf("%s findings = %#v, want none", speccheck.CodeOrdinalClaimed, findings)
		}
	})
}

func TestSameSpecDuplicateBesideAFulfilledClaimIsReportedOnce(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	const fulfilledClaim = "docs/adr/0042-first-decision.md"
	writeOrdinalSpecClaims(t, repoRoot, "claiming-spec",
		fulfilledClaim,
		"docs/adr/0042-second-decision.md",
	)
	writeCitationFixtureFile(t, repoRoot, fulfilledClaim, "# First decision\n")

	result, err := speccheck.CheckStage(specsRoot, repoRoot, "claiming-spec", speccheck.StageTasks)
	if err != nil {
		t.Fatalf("CheckStage(StageTasks): %v", err)
	}
	findings := findingsWithCode(result, speccheck.CodeOrdinalClaimed)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want one", speccheck.CodeOrdinalClaimed, findings)
	}
}

func TestSameSpecDuplicateFixNamesRenumbering(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	writeOrdinalSpecClaims(t, repoRoot, "claiming-spec",
		"docs/adr/0042-first-decision.md",
		"docs/adr/0042-second-decision.md",
	)

	result, err := speccheck.CheckStage(specsRoot, repoRoot, "claiming-spec", speccheck.StageTasks)
	if err != nil {
		t.Fatalf("CheckStage(StageTasks): %v", err)
	}
	findings := findingsWithCode(result, speccheck.CodeOrdinalClaimed)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want one", speccheck.CodeOrdinalClaimed, findings)
	}
	if summary := findings[0].Summary; !strings.Contains(summary, "is claimed by both") {
		t.Errorf("summary = %q, want same-Spec duplicate message", summary)
	}
	if fix := findings[0].Fix; !strings.Contains(fix, "Renumber one Task's `creates:` path") {
		t.Errorf("fix = %q, want renumbering instruction", fix)
	}
}

func TestFulfilledOrdinalClaimIsNotBlamedForALaterClaim(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	const fulfilledClaim = "docs/adr/0042-first-decision.md"
	writeOrdinalSpec(t, repoRoot, "first-spec", fulfilledClaim)
	writeCitationFixtureFile(t, repoRoot, fulfilledClaim, "# First decision\n")
	writeOrdinalSpec(t, repoRoot, "later-spec", "docs/adr/0042-later-decision.md")

	firstResult, err := speccheck.CheckStage(specsRoot, repoRoot, "first-spec", speccheck.StageTasks)
	if err != nil {
		t.Fatalf("CheckStage(first-spec, StageTasks): %v", err)
	}
	if findings := findingsWithCode(firstResult, speccheck.CodeOrdinalClaimed); len(findings) != 0 {
		t.Fatalf("fulfilled claim findings = %#v, want none", findings)
	}

	laterResult, err := speccheck.CheckStage(specsRoot, repoRoot, "later-spec", speccheck.StageTasks)
	if err != nil {
		t.Fatalf("CheckStage(later-spec, StageTasks): %v", err)
	}
	for _, finding := range findingsWithCode(laterResult, speccheck.CodeOrdinalClaimed) {
		if strings.Contains(finding.Summary, "already held by "+fulfilledClaim) {
			return
		}
	}
	t.Fatalf("later claim findings = %#v, want tree collision with %s",
		findingsWithCode(laterResult, speccheck.CodeOrdinalClaimed), fulfilledClaim)
}

func TestOrdinalCheckSkipsAnUnloadableActiveSpec(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	writeOrdinalSpec(t, repoRoot, "checked-spec", "docs/adr/0042-checked-decision.md")
	writeCitationFixtureFile(t, repoRoot, "docs/specs/broken-neighbour/_prd.md", `---
spec: broken-neighbour
status: active
created: 2026-09-24
surfaces: [backend]
---

# Broken neighbour
`)
	writeCitationFixtureFile(t, repoRoot, "docs/specs/broken-neighbour/_tasks.md", `---
schema: spec-tasks/v1
spec: broken-neighbour
qa: declined
qa_reason: no behavioral surface in this fixture
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
---
`)

	if _, err := speccheck.CheckStage(specsRoot, repoRoot, "checked-spec", speccheck.StageTasks); err != nil {
		t.Fatalf("CheckStage(StageTasks) with unloadable active Spec: %v", err)
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
	writeOrdinalSpecClaims(t, repoRoot, slug, claim)
}

func writeOrdinalSpecClaims(t *testing.T, repoRoot, slug string, claims ...string) {
	t.Helper()

	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_prd.md", `---
spec: `+slug+`
status: active
created: 2026-09-24
surfaces: [backend]
---

# Ordinal claim fixture
`)
	var nodes strings.Builder
	for index := range claims {
		id := fmt.Sprintf("task_%02d", index+1)
		fmt.Fprintf(&nodes, "    - id: %s\n      file: %s.md\n      needs: []\n", id, id)
	}
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_tasks.md", `---
schema: spec-tasks/v1
spec: `+slug+`
qa: declined
qa_reason: no behavioral surface in this fixture
graph:
  nodes:
`+nodes.String()+`---
`)
	for index, claim := range claims {
		id := fmt.Sprintf("task_%02d", index+1)
		writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/"+id+".md", `---
task: `+id+`
spec: `+slug+`
status: pending
type: backend
---

# Claim an ADR ordinal

## Context

- creates: `+"`"+claim+"`"+`

## Verification

- `+"`test -f "+claim+"`"+`
`)
	}
}
