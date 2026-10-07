package speccheck_test

import (
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func historyEntryWriteReduced(t *testing.T, root, path, front string) {
	t.Helper()
	content, err := spec.ReduceHistoryEntry([]byte(front+"\n# Retired entry\n\nOriginal observation.\n\n## Evidence\nLong full text.\n"), strings.Repeat("a", 40), path)
	if err != nil {
		t.Fatal(err)
	}
	writeFindingsArtifact(t, root, path, string(content))
}

func TestReducedFindingsKeepTheirLicensesAndMembers(t *testing.T) {
	t.Parallel()
	root := writeFindingsCarrier(t)
	const slug = "0002-record-owner"
	writeCheckArchiveRecord(t, root, "docs/history/specs", "docs/specs/"+slug, slug, strings.Repeat("a", 40))
	const rollup = "2026-10-06-rollup.md"
	const member = "2026-10-06-member.md"
	writeFindingsArtifact(t, root, "docs/findings/"+rollup, "---\nstatus: pending\nkind: rollup\nmembers:\n  - "+member+"\ncreated_at: 2026-10-06\nupdated_at: 2026-10-06\n---\n\n# Active Rollup\n")
	historyEntryWriteReduced(t, root, "docs/history/findings/2026-10-06-spec.md", "---\nstatus: done\nabsorbed_by: "+slug+"\ncreated_at: 2026-10-06\nupdated_at: 2026-10-06\n---\n")
	historyEntryWriteReduced(t, root, "docs/history/findings/"+member, "---\nstatus: done\nabsorbed_by: "+rollup+"\ncreated_at: 2026-10-06\nupdated_at: 2026-10-06\n---\n")
	const backlog = "docs/history/backlog/2026-10-06-closed.md"
	historyEntryWriteReduced(t, root, backlog, "---\ntype: fix\nstatus: closed\ncreated: 2026-10-06\nspec: null\nreason: no implementation remains\n---\n")
	historyEntryWriteReduced(t, root, "docs/history/findings/2026-10-06-closed.md", "---\nstatus: closed\nclosure_reason: no implementation remains\nclosure_evidence: "+backlog+"\ncreated_at: 2026-10-06\nupdated_at: 2026-10-06\n---\n")
	result := checkFindingsCarrier(t, root)
	if len(result.Findings) != 0 {
		t.Fatalf("reduced fixture findings = %+v", result.Findings)
	}
}

func TestReducedFindingWithABrokenLicenseStillFails(t *testing.T) {
	t.Parallel()
	root := writeFindingsCarrier(t)
	const path = "docs/history/findings/2026-10-06-broken.md"
	historyEntryWriteReduced(t, root, path, "---\nstatus: done\nabsorbed_by: missing-owner\ncreated_at: 2026-10-06\nupdated_at: 2026-10-06\n---\n")
	result := checkFindingsCarrier(t, root)
	findings := findingsWithCode(result, speccheck.CodeArchiveLicense)
	if len(findings) != 1 {
		t.Fatalf("archive license findings = %+v", findings)
	}
	requireRenderedFinding(t, result, speccheck.CodeArchiveLicense, path, 3)
	if !strings.Contains(findings[0].Summary, "missing-owner") {
		t.Fatalf("missing broken license diagnostic: %+v", findings[0])
	}
}
