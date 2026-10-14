package cli

// Suite: Legacy conversion links through both history commands.
// Invariant: plans and writes agree, and moved text names surviving files.
// Boundary IN: disposable Git repositories, temporary homes and public commands.
// Boundary OUT: Archive Record and promoted reference bytes and destinations.

import (
	"bytes"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

const historyLinkFolder = "docs/history/specs/0001-legacy"
const historyLinkPromotion = historyLinkFolder + "/references/lesson.md"
const historyLinkSummary = "See [lesson](../../references/lesson.md) and [ADR](../../adr/0001-rule.md)."
const historyLinkCopy = "# Lesson\n\n[PRD](../history/specs/0001-legacy.md#goals), [notes](../history/specs/0001-legacy.md), [peer](peer.md?raw=1#data), [ADR](../adr/0001-rule.md).\n"

func historyLinkFixture(t *testing.T, h historyFixture) spec.LegacyConversion {
	t.Helper()
	historyWrite(t, filepath.Join(h.repo, historyLinkFolder, "_prd.md"), "---\nspec: 0001-legacy\ncreated: 2026-01-01\n---\n\n# Legacy\n\nSee [lesson](references/lesson.md) and [ADR](../../../adr/0001-rule.md).\n")
	historyWrite(t, filepath.Join(h.repo, historyLinkPromotion), "# Lesson\n\n[PRD](../_prd.md#goals), [notes](../notes.md), [peer](peer.md?raw=1#data), [ADR](../../../../adr/0001-rule.md).\n")
	historyWrite(t, filepath.Join(h.repo, historyLinkFolder, "references/peer.md"), "# Peer\n\nUnchanged bytes.\n")
	historyWrite(t, filepath.Join(h.repo, historyLinkFolder, "notes.md"), "# Notes\n")
	historyWrite(t, filepath.Join(h.repo, "docs/adr/0001-rule.md"), "# Rule\n")
	revision := historyCommit(t, h.repo, "Legacy links (#42)")
	delivery, err := spec.FindLegacyDelivery(t.Context(), h.repo, "docs/specs", "docs/history/specs", "0001-legacy")
	if err != nil {
		t.Fatal(err)
	}
	plan, err := spec.PlanLegacyConversion(spec.LegacyConversionRequest{
		RepositoryRoot: h.repo, ArchiveRoot: filepath.Join(h.repo, "docs/history/specs"),
		Slug: "0001-legacy", SourceRevision: revision, Delivery: delivery,
		Promote: []string{historyLinkPromotion, historyLinkFolder + "/references/peer.md"},
	})
	if err != nil {
		t.Fatal(err)
	}
	return plan
}

func assertHistoryLinkWrites(t *testing.T, h historyFixture, plan spec.LegacyConversion) {
	t.Helper()
	for _, output := range []struct{ path, want string }{
		{plan.RecordPath, string(plan.Rendered)},
		{"docs/references/lesson.md", historyLinkCopy},
		{"docs/references/peer.md", "# Peer\n\nUnchanged bytes.\n"},
	} {
		content, err := os.ReadFile(filepath.Join(h.repo, output.path))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(content, []byte(output.want)) {
			t.Fatalf("%s = %q, want %q", output.path, content, output.want)
		}
		if output.path == plan.RecordPath {
			record, err := spec.ParseArchiveRecord(content)
			if err != nil || record.Outcome != historyLinkSummary {
				t.Fatalf("outcome = %q, err = %v, want %q", record.Outcome, err, historyLinkSummary)
			}
			content = []byte(record.Outcome)
		}
		for _, match := range regexp.MustCompile(`\]\(([^)]+)\)`).FindAllSubmatch(content, -1) {
			destination, err := url.Parse(string(match[1]))
			if err != nil {
				t.Fatal(err)
			}
			target := filepath.Join(h.repo, filepath.Dir(output.path), filepath.FromSlash(destination.Path))
			if _, err := os.Stat(target); err != nil {
				t.Fatalf("%s destination %q does not resolve: %v", output.path, match[1], err)
			}
		}
	}
	assertPathMissing(t, filepath.Join(h.repo, historyLinkFolder))
}

func TestHistorySanitizeRebasesPromotedLinks(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, false)
	plan := historyLinkFixture(t, h)
	historyTag(t, h)
	code, out, stderr := historyRun(t, h, "sanitize", "--apply", "--batch", "1", "--promote", historyLinkPromotion, "--promote", historyLinkFolder+"/references/peer.md")
	if code != 0 || stderr != "" || !strings.Contains(out, "wrote 1 Archive Record(s)") {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	assertHistoryLinkWrites(t, h, plan)
}

func TestBaselineUpdatePromotionRebasesLinks(t *testing.T) {
	t.Parallel()
	h := baselineCitationFixture(t, false)
	plan := historyLinkFixture(t, h)
	result, out, stderr, code := baselineHistoryRun(t, h, "json", "--yes", "--promote", historyLinkPromotion, "--promote", historyLinkFolder+"/references/peer.md")
	if code != 0 || stderr != "" || result.History.Status != "applied" {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	assertHistoryLinkWrites(t, h, plan)
}
