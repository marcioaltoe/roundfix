package speccheck_test

import (
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

const backlogCarrierSlug = "0001-backlog-carrier"

func TestCheckBacklogUnmoved(t *testing.T) {
	t.Parallel()

	t.Run("red promoted entry never moved", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeBacklogCarrier(t)
		const entryPath = "docs/backlog/2026-08-06-still-here.md"
		writeFindingsArtifact(t, repoRoot, entryPath,
			"---\ntype: perf\nstatus: promoted\ncreated: 2026-08-06\nspec: "+backlogCarrierSlug+"\n---\n\n# Still here\n")

		result := checkBacklogCarrier(t, repoRoot)
		finding := requireRenderedFinding(t, result, speccheck.CodeBacklogUnmoved, entryPath, 2)
		if !strings.Contains(finding.Summary, "still lives in docs/backlog") {
			t.Fatalf("summary = %q, want the unmoved statement", finding.Summary)
		}
		if !strings.Contains(finding.Fix, "docs/specs/"+backlogCarrierSlug+"/references/") {
			t.Fatalf("fix = %q, want the destination path", finding.Fix)
		}
	})

	t.Run("red promoted without a Spec", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeBacklogCarrier(t)
		const entryPath = "docs/backlog/2026-08-06-no-spec.md"
		writeFindingsArtifact(t, repoRoot, entryPath,
			"---\ntype: feat\nstatus: promoted\ncreated: 2026-08-06\nspec: null\n---\n\n# No spec\n")

		result := checkBacklogCarrier(t, repoRoot)
		finding := requireRenderedFinding(t, result, speccheck.CodeBacklogUnmoved, entryPath, 2)
		if !strings.Contains(finding.Summary, "without naming its Spec") {
			t.Fatalf("summary = %q, want the missing-Spec statement", finding.Summary)
		}
	})

	t.Run("red promoted to an unresolvable Spec", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeBacklogCarrier(t)
		const entryPath = "docs/backlog/2026-08-06-ghost-spec.md"
		writeFindingsArtifact(t, repoRoot, entryPath,
			"---\ntype: fix\nstatus: promoted\ncreated: 2026-08-06\nspec: 9999-not-a-spec\n---\n\n# Ghost\n")

		result := checkBacklogCarrier(t, repoRoot)
		finding := requireRenderedFinding(t, result, speccheck.CodeBacklogUnmoved, entryPath, 4)
		if !strings.Contains(finding.Summary, "unresolvable Spec") {
			t.Fatalf("summary = %q, want the unresolvable-Spec statement", finding.Summary)
		}
	})

	t.Run("green open entry stays put", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeBacklogCarrier(t)
		writeFindingsArtifact(t, repoRoot, "docs/backlog/2026-08-06-open.md",
			"---\ntype: perf\nstatus: open\ncreated: 2026-08-06\nspec: null\n---\n\n# Open\n")

		result := checkBacklogCarrier(t, repoRoot)
		if findings := findingsWithCode(result, speccheck.CodeBacklogUnmoved); len(findings) != 0 {
			t.Fatalf("%s findings = %#v, want none", speccheck.CodeBacklogUnmoved, findings)
		}
	})

	t.Run("green declined entry stays put", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeBacklogCarrier(t)
		writeFindingsArtifact(t, repoRoot, "docs/backlog/2026-08-06-declined.md",
			"---\ntype: feat\nstatus: declined\ncreated: 2026-08-06\nspec: null\nreason: superseded\n---\n\n# Declined\n")

		result := checkBacklogCarrier(t, repoRoot)
		if findings := findingsWithCode(result, speccheck.CodeBacklogUnmoved); len(findings) != 0 {
			t.Fatalf("%s findings = %#v, want none", speccheck.CodeBacklogUnmoved, findings)
		}
	})

	t.Run("skips when no backlog directory exists", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeBacklogCarrier(t)

		result := checkBacklogCarrier(t, repoRoot)
		if findings := findingsWithCode(result, speccheck.CodeBacklogUnmoved); len(findings) != 0 {
			t.Fatalf("%s findings = %#v, want silent skip", speccheck.CodeBacklogUnmoved, findings)
		}
		if !hasSkip(result, speccheck.CodeBacklogUnmoved, "docs/backlog") {
			t.Fatalf("Skipped = %#v, want %s missing docs/backlog", result.Skipped, speccheck.CodeBacklogUnmoved)
		}
	})

	t.Run("skips when the backlog directory is empty", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeBacklogCarrier(t)
		writeFindingsArtifact(t, repoRoot, "docs/backlog/_index.md", "# Index\n")

		result := checkBacklogCarrier(t, repoRoot)
		if findings := findingsWithCode(result, speccheck.CodeBacklogUnmoved); len(findings) != 0 {
			t.Fatalf("%s findings = %#v, want silent skip", speccheck.CodeBacklogUnmoved, findings)
		}
		if !hasSkip(result, speccheck.CodeBacklogUnmoved, "docs/backlog") {
			t.Fatalf("Skipped = %#v, want %s missing docs/backlog", result.Skipped, speccheck.CodeBacklogUnmoved)
		}
	})
}

func writeBacklogCarrier(t *testing.T) string {
	t.Helper()

	repoRoot := t.TempDir()
	writeFindingsArtifact(t, repoRoot,
		filepath.ToSlash(filepath.Join("docs", "specs", backlogCarrierSlug, "_prd.md")),
		"---\nspec: "+backlogCarrierSlug+"\nstatus: active\n---\n\n# Carrier\n")
	return repoRoot
}

func checkBacklogCarrier(t *testing.T, repoRoot string) speccheck.Result {
	t.Helper()

	result, err := speccheck.Check(filepath.Join(repoRoot, "docs", "specs"), repoRoot, backlogCarrierSlug)
	if err != nil {
		t.Fatalf("Check(backlog carrier) error = %v", err)
	}
	return result
}
