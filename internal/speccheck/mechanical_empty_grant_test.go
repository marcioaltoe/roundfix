// Suite: mechanical audit of an explicit empty authorization grant.
// Invariant: an empty grant rejects every Governed Path and permits ordinary changed paths.
// Boundary IN: public mechanical API, real temporary Git history, and the authorization reader.
// Boundary OUT: Daemon scheduling and QA verdict computation.
package speccheck_test

import (
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

const mechanicalEmptyGrantSlug = "mechanical-empty-grant"

func TestMechanicalAuthPathsRefusesAGovernedChangeUnderAnEmptyGrant(t *testing.T) {
	t.Parallel()

	result := runMechanicalEmptyGrantChange(t, "Makefile")
	findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalAuthPaths)
	if len(findings) != 1 || findings[0].RowHint != "task_01" || !strings.Contains(findings[0].Detail, "Makefile") {
		t.Fatalf("%s findings = %#v, want one Task finding", speccheck.CodeMechanicalAuthPaths, findings)
	}
}

func TestMechanicalAuthPathsAcceptsAnOrdinaryChangeUnderAnEmptyGrant(t *testing.T) {
	t.Parallel()

	result := runMechanicalEmptyGrantChange(t, "ordinary.txt")
	assertNoMechanicalCode(t, result, speccheck.CodeMechanicalAuthPaths)
}

func runMechanicalEmptyGrantChange(t *testing.T, changedPath string) speccheck.MechanicalResult {
	t.Helper()

	repoRoot := newMechanicalGitRepo(t)
	const authorizationPath = "docs/specs/" + mechanicalEmptyGrantSlug + "/_authorization.md"
	writeMechanicalFile(t, repoRoot, authorizationPath, mechanicalEmptyGrantAuthorization())
	writeMechanicalFile(t, repoRoot, changedPath, "before\n")
	target := commitMechanicalFiles(t, repoRoot, "record empty grant", authorizationPath, changedPath)

	writeMechanicalFile(t, repoRoot, changedPath, "after\n")
	consumer := commitMechanicalFiles(t, repoRoot, "change "+changedPath, changedPath)
	return runMechanical(t, speccheck.MechanicalRequest{
		RepoRoot:               repoRoot,
		AuthorizationPath:      authorizationPath,
		ConsumingSpec:          mechanicalEmptyGrantSlug,
		DeliveryTargetRevision: target,
		TaskCommits: []speccheck.MechanicalTaskCommit{{
			TaskID: "task_01",
			SHA:    consumer,
		}},
	})
}

func mechanicalEmptyGrantAuthorization() string {
	return "---\n" +
		"status: approved\n" +
		"granted: 2026-09-30\n" +
		"action: grant operations without bounded paths\n" +
		"consuming: " + mechanicalEmptyGrantSlug + "\n" +
		"paths: []\n" +
		"operations:\n" +
		"  - implement\n" +
		"---\n"
}
