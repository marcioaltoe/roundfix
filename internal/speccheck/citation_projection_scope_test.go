// Suite: Spec citation authored projection scope
// Invariant: only top-level task_*.md files lose non-authorial sections.
// Boundary IN: public speccheck API and plain temporary repository directories
// Boundary OUT: Git-backed ADR horizon behavior and corpus characterization fixtures
package speccheck_test

import "testing"

func TestAReferenceNamedLikeATaskIsReadInFull(t *testing.T) {
	t.Parallel()

	repoRoot := newCitationProjectionFixture(t)
	const referencePath = "references/task_example.md"
	writeCitationFixtureFile(t, repoRoot, citationProjectionSpecPath(referencePath), `# Adopted source

## Result

The authored source cites ADR-0002.
`)

	requireProjectionUnlistedFinding(t, checkCitationProjection(t, repoRoot), citationProjectionSpecPath(referencePath), 5)
}

func TestANestedFileNamedLikeATaskIsReadInFull(t *testing.T) {
	t.Parallel()

	repoRoot := newCitationProjectionFixture(t)
	const notesPath = "notes/task_02.md"
	writeCitationFixtureFile(t, repoRoot, citationProjectionSpecPath(notesPath), `# Authored notes

## Recorded paths

The authored notes cite ADR-0002.
`)

	requireProjectionUnlistedFinding(t, checkCitationProjection(t, repoRoot), citationProjectionSpecPath(notesPath), 5)
}

func TestTheSpecsOwnTaskFileKeepsItsProjection(t *testing.T) {
	t.Parallel()

	repoRoot := newCitationProjectionFixture(t)
	writeCitationFixtureFile(t, repoRoot, citationProjectionSpecPath("task_01.md"), `# Task

## Result

The implementation fixture mentions ADR-0002.
`)

	requireNoProjectionUnlistedFinding(t, checkCitationProjection(t, repoRoot))
}
