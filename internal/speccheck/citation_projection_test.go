// Suite: Spec citation authored projection
// Invariant: only author-owned Spec text creates unlisted ADR obligations, with original locations.
// Boundary IN: public speccheck API and plain temporary repository directories
// Boundary OUT: Git-backed ADR horizon behavior and corpus characterization fixtures
package speccheck_test

import (
	"path/filepath"
	"testing"

	"roundfix/internal/speccheck"
)

const citationProjectionSlug = "citation-projection"

func TestASpecCitationInATaskResultIsNotAnObligation(t *testing.T) {
	t.Parallel()

	repoRoot := newCitationProjectionFixture(t)
	writeCitationFixtureFile(t, repoRoot, citationProjectionSpecPath("task_01.md"), `# Task

## Requirements

No decision citation here.

## Result

The implementation fixture mentions ADR-0002.
`)

	requireNoProjectionUnlistedFinding(t, checkCitationProjection(t, repoRoot))
}

func TestASpecCitationInADaemonSectionIsNotAnObligation(t *testing.T) {
	t.Parallel()

	for _, heading := range []string{
		"## Recorded paths",
		"## Carry-forward provenance",
	} {
		heading := heading
		t.Run(heading, func(t *testing.T) {
			t.Parallel()

			repoRoot := newCitationProjectionFixture(t)
			writeCitationFixtureFile(t, repoRoot, citationProjectionSpecPath("task_01.md"), "# Task\n\n"+heading+"\n\nThe Daemon fixture mentions ADR-0002.\n")

			requireNoProjectionUnlistedFinding(t, checkCitationProjection(t, repoRoot))
		})
	}
}

func TestASpecCitationInAQAReportIsNotAnObligation(t *testing.T) {
	t.Parallel()

	for _, relative := range []string{
		"qa/qa-report-2026-09-29.md",
		"qa/evidence/user-flow.md",
	} {
		relative := relative
		t.Run(relative, func(t *testing.T) {
			t.Parallel()

			repoRoot := newCitationProjectionFixture(t)
			writeCitationFixtureFile(t, repoRoot, citationProjectionSpecPath(relative), "# QA fixture\n\nThe gate fixture mentions ADR-0002.\n")

			requireNoProjectionUnlistedFinding(t, checkCitationProjection(t, repoRoot))
		})
	}
}

func TestAnAuthoredSpecCitationStillMustBeListed(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name     string
		relative string
		content  string
		wantLine int
	}{
		{
			name:     "Task Requirements",
			relative: "task_01.md",
			content:  "# Task\n\n## Requirements\n\nThe authored requirement cites ADR-0002.\n",
			wantLine: 5,
		},
		{
			name:     "TechSpec",
			relative: "_techspec.md",
			content:  "# TechSpec\n\n## Design\n\nThe authored design cites ADR-0002.\n",
			wantLine: 5,
		},
		{
			name:     "adopted reference",
			relative: "references/source.md",
			content:  "# Adopted source\n\nThe authored source cites ADR-0002.\n",
			wantLine: 3,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoRoot := newCitationProjectionFixture(t)
			writeCitationFixtureFile(t, repoRoot, citationProjectionSpecPath(tt.relative), tt.content)

			requireProjectionUnlistedFinding(t, checkCitationProjection(t, repoRoot), citationProjectionSpecPath(tt.relative), tt.wantLine)
		})
	}
}

func TestAnAuthoredSectionAfterAResultIsStillRead(t *testing.T) {
	t.Parallel()

	repoRoot := newCitationProjectionFixture(t)
	const taskPath = "task_01.md"
	writeCitationFixtureFile(t, repoRoot, citationProjectionSpecPath(taskPath), `# Task

## Result

The implementation fixture mentions ADR-0002.

## Requirements

The authored requirement cites ADR-0002.
`)

	requireProjectionUnlistedFinding(t, checkCitationProjection(t, repoRoot), citationProjectionSpecPath(taskPath), 9)
}

func newCitationProjectionFixture(t *testing.T) string {
	t.Helper()

	repoRoot := t.TempDir()
	for _, source := range []string{
		"docs/agents/agent-instructions.md",
		"docs/agents/cli.md",
		"docs/agents/domain.md",
	} {
		writeCitationFixtureFile(t, repoRoot, source, "# Fixture source\n")
	}
	writeCitationFixtureFile(t, repoRoot, "docs/adr/0001-listed-decision.md", "---\nstatus: accepted\n---\n\n# Listed decision\n")
	writeCitationFixtureFile(t, repoRoot, "docs/adr/0002-projection-target.md", "---\nstatus: accepted\n---\n\n# Projection target\n")
	writeCitationFixtureFile(t, repoRoot, citationProjectionSpecPath("_prd.md"), `---
spec: citation-projection
status: active
---

# Citation projection

## Project Constraints

- Identifier strategy: not applicable — fixture only. Source: `+"`docs/agents/domain.md`"+`.
- Authentication and HTTP: not applicable — fixture only. Source: `+"`docs/agents/cli.md`"+`.
- Active ADR obligations: applicable — ADR-0001. Source: `+"`docs/agents/domain.md`"+`.
- Tooling authority: not applicable — fixture only. Source: `+"`docs/agents/agent-instructions.md`"+`.
`)
	return repoRoot
}

func citationProjectionSpecPath(relative string) string {
	return filepath.ToSlash(filepath.Join("docs", "specs", citationProjectionSlug, filepath.FromSlash(relative)))
}

func checkCitationProjection(t *testing.T, repoRoot string) speccheck.Result {
	t.Helper()

	result, err := speccheck.Check(filepath.Join(repoRoot, "docs", "specs"), repoRoot, citationProjectionSlug)
	if err != nil {
		t.Fatalf("Check(%q): %v", citationProjectionSlug, err)
	}
	return result
}

func requireNoProjectionUnlistedFinding(t *testing.T, result speccheck.Result) {
	t.Helper()

	if findings := findingsWithCode(result, speccheck.CodeADRUnlisted); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want non-authorial citation ignored", speccheck.CodeADRUnlisted, findings)
	}
}

func requireProjectionUnlistedFinding(t *testing.T, result speccheck.Result, citationPath string, citationLine int) {
	t.Helper()

	findings := findingsWithCode(result, speccheck.CodeADRUnlisted)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want exactly one authored citation", speccheck.CodeADRUnlisted, findings)
	}
	finding := findings[0]
	const prdPath = "docs/specs/citation-projection/_prd.md"
	wantSummary := citationPath + " cites ADR-0002 (Projection target), but " + prdPath + " does not list it under Active ADR obligations"
	if finding.Summary != wantSummary {
		t.Fatalf("summary = %q, want %q", finding.Summary, wantSummary)
	}
	wantFix := "List ADR-0002 in the Active ADR obligations row in " + prdPath + " using the recognised ADR-NNNN form, or remove the stale citation."
	if finding.Fix != wantFix {
		t.Fatalf("fix = %q, want %q", finding.Fix, wantFix)
	}
	if !hasExactLocation(finding, citationPath, citationLine) {
		t.Fatalf("locations = %#v, want %s:%d", finding.Where, citationPath, citationLine)
	}
}
