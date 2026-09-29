// Suite: Related-ADR horizon
// Invariant: SC-ADR-RELATED considers only ADRs that existed when the Spec PRD was added.
// Boundary IN: public speccheck API, accepted ADR files, and disposable Git histories
// Boundary OUT: CLI rendering and this repository's own history
package speccheck_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/speccheck"
)

const adrHorizonSlug = "adr-horizon"

func TestAnADRCommittedAfterTheSpecOpensNoRelatedGap(t *testing.T) {
	t.Parallel()

	repoRoot := newADRHorizonRepository(t)
	writeADRHorizonBase(t, repoRoot)
	commitADRHorizonFiles(t, repoRoot, "base decision")
	writeADRHorizonPRD(t, repoRoot, "")
	commitADRHorizonFiles(t, repoRoot, "spec")
	writeADRHorizonRelatedADR(t, repoRoot)
	commitADRHorizonFiles(t, repoRoot, "later decision")

	result := checkADRHorizon(t, repoRoot, repoRoot)
	requireADRHorizonFindingCount(t, result, speccheck.CodeADRRelated, 0)
}

func TestAnADRCommittedBeforeTheSpecOpensTheRelatedGap(t *testing.T) {
	t.Parallel()

	repoRoot := newADRHorizonRepository(t)
	writeADRHorizonBase(t, repoRoot)
	writeADRHorizonRelatedADR(t, repoRoot)
	commitADRHorizonFiles(t, repoRoot, "decisions")
	writeADRHorizonPRD(t, repoRoot, "")
	commitADRHorizonFiles(t, repoRoot, "spec")

	result := checkADRHorizon(t, repoRoot, repoRoot)
	requireADRHorizonFindingCount(t, result, speccheck.CodeADRRelated, 1)
}

func TestAnADRCommittedWithTheSpecOpensTheRelatedGap(t *testing.T) {
	t.Parallel()

	repoRoot := newADRHorizonRepository(t)
	writeADRHorizonBase(t, repoRoot)
	commitADRHorizonFiles(t, repoRoot, "base decision")
	writeADRHorizonPRD(t, repoRoot, "")
	writeADRHorizonRelatedADR(t, repoRoot)
	commitADRHorizonFiles(t, repoRoot, "spec and related decision")

	result := checkADRHorizon(t, repoRoot, repoRoot)
	requireADRHorizonFindingCount(t, result, speccheck.CodeADRRelated, 1)
}

func TestTheHorizonUsesTheNewestAddingCommit(t *testing.T) {
	t.Parallel()

	repoRoot := newADRHorizonRepository(t)
	writeADRHorizonBase(t, repoRoot)
	writeADRHorizonRelatedADR(t, repoRoot)
	commitADRHorizonFiles(t, repoRoot, "decisions")
	removeADRHorizonFile(t, repoRoot, "docs/adr/0002-related-decision.md")
	commitADRHorizonFiles(t, repoRoot, "remove related decision")
	writeADRHorizonPRD(t, repoRoot, "")
	commitADRHorizonFiles(t, repoRoot, "spec")
	writeADRHorizonRelatedADR(t, repoRoot)
	commitADRHorizonFiles(t, repoRoot, "restore related decision")

	result := checkADRHorizon(t, repoRoot, repoRoot)
	requireADRHorizonFindingCount(t, result, speccheck.CodeADRRelated, 0)
}

func TestAnUncommittedSpecKeepsTheFullRelatedCheck(t *testing.T) {
	t.Parallel()

	repoRoot := newADRHorizonRepository(t)
	writeADRHorizonBase(t, repoRoot)
	writeADRHorizonRelatedADR(t, repoRoot)
	commitADRHorizonFiles(t, repoRoot, "decisions")
	writeADRHorizonPRD(t, repoRoot, "")

	result := checkADRHorizon(t, repoRoot, repoRoot)
	requireADRHorizonFindingCount(t, result, speccheck.CodeADRRelated, 1)
}

func TestAnUncommittedADRIsOutsideACommittedSpecsHorizon(t *testing.T) {
	t.Parallel()

	repoRoot := newADRHorizonRepository(t)
	writeADRHorizonBase(t, repoRoot)
	commitADRHorizonFiles(t, repoRoot, "base decision")
	writeADRHorizonPRD(t, repoRoot, "")
	commitADRHorizonFiles(t, repoRoot, "spec")
	writeADRHorizonRelatedADR(t, repoRoot)

	result := checkADRHorizon(t, repoRoot, repoRoot)
	requireADRHorizonFindingCount(t, result, speccheck.CodeADRRelated, 0)
}

func TestASpecWithoutGitKeepsTheFullRelatedCheck(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	writeADRHorizonBase(t, repoRoot)
	writeADRHorizonRelatedADR(t, repoRoot)
	writeADRHorizonPRD(t, repoRoot, "")

	result := checkADRHorizon(t, repoRoot, repoRoot)
	requireADRHorizonFindingCount(t, result, speccheck.CodeADRRelated, 1)
}

func TestASpecInAnotherRepositoryKeepsTheFullRelatedCheck(t *testing.T) {
	t.Parallel()

	repoRoot := newADRHorizonRepository(t)
	writeADRHorizonBase(t, repoRoot)
	writeADRHorizonRelatedADR(t, repoRoot)
	commitADRHorizonFiles(t, repoRoot, "decisions")

	specRepoRoot := newADRHorizonRepository(t)
	writeADRHorizonPRD(t, specRepoRoot, "")
	commitADRHorizonFiles(t, specRepoRoot, "spec")

	result := checkADRHorizon(t, repoRoot, specRepoRoot)
	requireADRHorizonFindingCount(t, result, speccheck.CodeADRRelated, 1)
}

func TestAShallowHistoryKeepsTheFullRelatedCheck(t *testing.T) {
	t.Parallel()

	sourceRoot := filepath.Join(t.TempDir(), "source")
	gittest.InitRepo(t, sourceRoot, "-b", "main")
	writeADRHorizonBase(t, sourceRoot)
	writeADRHorizonPRD(t, sourceRoot, "")
	commitADRHorizonFiles(t, sourceRoot, "spec")
	writeADRHorizonRelatedADR(t, sourceRoot)
	commitADRHorizonFiles(t, sourceRoot, "later decision")

	cloneRoot := filepath.Join(t.TempDir(), "clone")
	gittest.Run(t, "", "clone", "--depth=1", "file://"+filepath.ToSlash(sourceRoot), cloneRoot)
	gittest.Harden(t, cloneRoot)

	result := checkADRHorizon(t, cloneRoot, cloneRoot)
	requireADRHorizonFindingCount(t, result, speccheck.CodeADRRelated, 1)
}

func TestTheHorizonLeavesUnlistedCitationsChecked(t *testing.T) {
	t.Parallel()

	repoRoot := newADRHorizonRepository(t)
	writeADRHorizonBase(t, repoRoot)
	writeADRHorizonPRD(t, repoRoot, "")
	commitADRHorizonFiles(t, repoRoot, "spec")
	writeADRHorizonRelatedADR(t, repoRoot)
	commitADRHorizonFiles(t, repoRoot, "later decision")
	writeADRHorizonPRD(t, repoRoot, "\nADR-0002.\n")
	commitADRHorizonFiles(t, repoRoot, "cite later decision")

	result := checkADRHorizon(t, repoRoot, repoRoot)
	findings := findingsWithCode(result, speccheck.CodeADRUnlisted)
	if len(findings) != 1 || !strings.Contains(findings[0].Summary, "ADR-0002") {
		t.Fatalf("%s findings = %#v, want the later ADR cited but unlisted", speccheck.CodeADRUnlisted, findings)
	}
	requireADRHorizonFindingCount(t, result, speccheck.CodeADRRelated, 0)
}

func newADRHorizonRepository(t *testing.T) string {
	t.Helper()

	repoRoot := t.TempDir()
	gittest.InitRepo(t, repoRoot, "-b", "main")
	return repoRoot
}

func writeADRHorizonBase(t *testing.T, repoRoot string) {
	t.Helper()

	for _, source := range []string{
		"docs/agents/agent-instructions.md",
		"docs/agents/cli.md",
		"docs/agents/domain.md",
	} {
		writeADRHorizonFile(t, repoRoot, source, "# Fixture source\n")
	}
	writeADRHorizonFile(t, repoRoot, "docs/adr/0001-base-decision.md", `---
status: accepted
---

# Base decision
`)
}

func writeADRHorizonRelatedADR(t *testing.T, repoRoot string) {
	t.Helper()

	writeADRHorizonFile(t, repoRoot, "docs/adr/0002-related-decision.md", `---
status: accepted
---

# Related decision

This decision cites ADR-0001.
`)
}

func writeADRHorizonPRD(t *testing.T, repoRoot, references string) {
	t.Helper()

	writeADRHorizonFile(t, repoRoot, "docs/specs/"+adrHorizonSlug+"/_prd.md", `---
spec: adr-horizon
status: active
---

# ADR horizon

## Project Constraints

- Identifier strategy: not applicable — fixture only. Source: `+"`docs/agents/domain.md`"+`.
- Authentication and HTTP: not applicable — fixture only. Source: `+"`docs/agents/cli.md`"+`.
- Active ADR obligations: applicable — ADR-0001. Source: `+"`docs/agents/domain.md`"+`.
- Tooling authority: not applicable — fixture only. Source: `+"`docs/agents/agent-instructions.md`"+`.

## References
`+references)
}

func writeADRHorizonFile(t *testing.T, repoRoot, relative, content string) {
	t.Helper()

	path := filepath.Join(repoRoot, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent for %s: %v", relative, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relative, err)
	}
}

func removeADRHorizonFile(t *testing.T, repoRoot, relative string) {
	t.Helper()

	if err := os.Remove(filepath.Join(repoRoot, filepath.FromSlash(relative))); err != nil {
		t.Fatalf("remove %s: %v", relative, err)
	}
}

func commitADRHorizonFiles(t *testing.T, repoRoot, message string) {
	t.Helper()

	gittest.Run(t, repoRoot, "add", "--all")
	gittest.Run(t, repoRoot, "commit", "-m", message)
}

func checkADRHorizon(t *testing.T, repoRoot, specRepoRoot string) speccheck.Result {
	t.Helper()

	result, err := speccheck.Check(filepath.Join(specRepoRoot, "docs", "specs"), repoRoot, adrHorizonSlug)
	if err != nil {
		t.Fatalf("Check(%s) error = %v", adrHorizonSlug, err)
	}
	return result
}

func requireADRHorizonFindingCount(t *testing.T, result speccheck.Result, code string, want int) {
	t.Helper()

	if findings := findingsWithCode(result, code); len(findings) != want {
		t.Fatalf("%s findings = %#v, want %d", code, findings, want)
	}
}
