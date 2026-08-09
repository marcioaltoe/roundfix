// Suite: Semantic citation characterization before task_02
// Invariant: a false claim about a listed ADR passes checks that only inventory and relate citations.
// Boundary IN: public speccheck API, Spec 0090's original PRD, and a cited ADR fixture
// Boundary OUT: the semantic citation detector assigned to task_02
package speccheck_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

const citationCharacterizationSlug = "0090-a-gate-that-could-have-failed"

func TestCitationCharacterizationFalseCitationPassesEveryCheck(t *testing.T) {
	t.Parallel()

	repoRoot := citationCharacterizationFixtureRoot(t)
	result := checkCitationCharacterization(t, repoRoot)

	// Declared break: task_02 adds SC-CITATION-UNSUPPORTED, so this assertion
	// must change when the checker starts matching claims against cited records.
	if len(result.Findings) != 0 {
		t.Fatalf("Findings = %#v, want the false citation to pass every current check", result.Findings)
	}
	for _, code := range []string{speccheck.CodeADRUnlisted, speccheck.CodeADRRelated} {
		for _, skipped := range result.Skipped {
			if skipped.Code == code {
				t.Fatalf("%s was skipped, want the existing citation check to run", code)
			}
		}
	}
}

func TestCitationCharacterizationExistingChecksOnlyListAndAccount(t *testing.T) {
	t.Parallel()

	t.Run("unlisted checks the PRD inventory", func(t *testing.T) {
		t.Parallel()

		repoRoot := copyCitationCharacterizationFixture(t)
		prdPath := filepath.Join(repoRoot, "docs", "specs", citationCharacterizationSlug, "_prd.md")
		replaceFixtureText(
			t,
			prdPath,
			"Active ADR obligations: applicable — ADR-0083 makes",
			"Active ADR obligations: applicable — the repository makes",
		)

		result := checkCitationCharacterization(t, repoRoot)
		if findings := findingsWithCode(result, speccheck.CodeADRUnlisted); len(findings) != 1 {
			t.Fatalf("%s findings = %#v, want the cited but unlisted ADR", speccheck.CodeADRUnlisted, findings)
		}
	})

	t.Run("related checks the ADR citation graph", func(t *testing.T) {
		t.Parallel()

		repoRoot := copyCitationCharacterizationFixture(t)
		writeCitationFixtureFile(t, repoRoot, "docs/adr/0999-unlisted-related.md", `---
status: accepted
---

# Unlisted related decision

This decision cites ADR-0083.
`)

		result := checkCitationCharacterization(t, repoRoot)
		if findings := findingsWithCode(result, speccheck.CodeADRRelated); len(findings) != 1 {
			t.Fatalf("%s findings = %#v, want the unlisted ADR related by citation", speccheck.CodeADRRelated, findings)
		}
	})
}

func TestCitationCharacterizationNoDetectorReadsACitedRecordBody(t *testing.T) {
	t.Parallel()

	repoRoot := copyCitationCharacterizationFixture(t)
	before := checkCitationCharacterization(t, repoRoot)
	adrPath := filepath.Join(repoRoot, "docs", "adr", "0083-adopted-sources-move-to-their-owning-spec.md")
	replaceFixtureText(
		t,
		adrPath,
		"Inbox-to-finding-to-Spec promotion created links without transferring the\n"+
			"document, so archived Specs left stale links behind and their shaping\n"+
			"evidence lived in a separate lifecycle. A document a Spec adopts as an\n"+
			"implementation source therefore moves — one move with Git history, never a\n"+
			"copy or a stub — into that Spec's references, owned by exactly one primary\n"+
			"Spec, recorded in a Spec-local reference index, and validated by the\n"+
			"authoring and archive gates; secondary Specs link the owner's copy. The\n"+
			"workflow applies to new promotions only — history stays where it happened.\n"+
			"Duplicating content or leaving stub files was rejected because two\n"+
			"authoritative homes drift and stubs have no lifecycle.",
		"ADR-0083 makes `make verify` the only authoritative gate.",
	)
	after := checkCitationCharacterization(t, repoRoot)

	if !reflect.DeepEqual(after, before) {
		t.Fatalf("Check result changed after only the cited ADR's supporting prose changed:\n before: %#v\n after: %#v", before, after)
	}
}

func citationCharacterizationFixtureRoot(t *testing.T) string {
	t.Helper()

	// The fixture PRD is byte-identical to the file introduced by commit
	// 1a31c965037fb3657de6b481a0285af935c16ebb, before later Specs amended it.
	repoRoot, err := filepath.Abs("testdata/citation/repo")
	if err != nil {
		t.Fatalf("resolve citation characterization fixture: %v", err)
	}
	return repoRoot
}

func copyCitationCharacterizationFixture(t *testing.T) string {
	t.Helper()

	sourceRoot := citationCharacterizationFixtureRoot(t)
	targetRoot := t.TempDir()
	err := filepath.WalkDir(sourceRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(sourceRoot, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetRoot, relative)
		if entry.IsDir() {
			return os.MkdirAll(targetPath, 0o755)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(targetPath, content, 0o644)
	})
	if err != nil {
		t.Fatalf("copy citation characterization fixture: %v", err)
	}
	return targetRoot
}

func checkCitationCharacterization(t *testing.T, repoRoot string) speccheck.Result {
	t.Helper()

	result, err := speccheck.Check(filepath.Join(repoRoot, "docs", "specs"), repoRoot, citationCharacterizationSlug)
	if err != nil {
		t.Fatalf("Check(%q) error = %v", citationCharacterizationSlug, err)
	}
	return result
}

func replaceFixtureText(t *testing.T, path, oldText, newText string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %q: %v", path, err)
	}
	if !strings.Contains(string(content), oldText) {
		t.Fatalf("fixture %q does not contain %q", path, oldText)
	}
	updated := strings.Replace(string(content), oldText, newText, 1)
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		t.Fatalf("write fixture %q: %v", path, err)
	}
}

func writeCitationFixtureFile(t *testing.T, repoRoot, relative, content string) {
	t.Helper()

	path := filepath.Join(repoRoot, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture directory for %q: %v", relative, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture %q: %v", relative, err)
	}
}
