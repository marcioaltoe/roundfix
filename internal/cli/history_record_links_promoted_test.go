package cli

// Suite: promoted Record Link Repair through both history plan commands.
// Invariant: the plan names the surviving copy and approval binds that rewrite.
// Boundary IN: Delivery Base record bytes, private repositories and homes.
// Boundary OUT: complete sanitize transcript, JSON citations and Plan Digest.

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func promotedRecordLinksFixture(t *testing.T, baseline bool) (historyFixture, string) {
	t.Helper()
	var h historyFixture
	if baseline {
		h = historyFixture{repo: newBaselineUpdateRepository(t), home: t.TempDir()}
		h.env = commandEnvironmentForTest(t)
		h.env.homeDir, h.env.workDir = h.home, h.repo
		h.env.environ = withEnvValue(h.env.environ, "HOME", h.home)
	} else {
		h = newHistoryFixture(t, false)
		// Keep precisely the Characterization Repository's two missing links.
		historyWrite(t, filepath.Join(h.repo, "README.md"), "# Demo\n")
	}
	source := "docs/specs/0001-demo"
	historyWrite(t, filepath.Join(h.repo, source, "_prd.md"), "# Demo\n\nPreserved outcome.\n")
	for _, name := range []string{"promoted.md", "removed.md"} {
		historyWrite(t, filepath.Join(h.repo, source, "references", name), "# Reference\n")
	}
	historyWrite(t, filepath.Join(h.repo, "docs/specs/0002-other/_prd.md"), "# Other\n")
	revision := historyCommit(t, h.repo, "docs: demo sources")
	// Pin origin/HEAD to the local main shape without contacting a remote.
	gittest.Run(t, h.repo, "remote", "add", "origin", h.repo)
	gittest.Run(t, h.repo, "update-ref", "refs/remotes/origin/main", revision)
	gittest.Run(t, h.repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	gittest.Run(t, h.repo, "rm", "-r", source)
	historyWrite(t, filepath.Join(h.repo, "docs/references/promoted.md"), "# Reference\n")
	// The earlier archive copied satisfied-by verbatim. Do not use today's
	// builder: it rebases these entries and would erase the repair's input.
	record := fmt.Sprintf("---\nschema: %s\nspec: 0001-demo\ntitle: Demo\nstatus: archived\ncreated: 2026-10-20\narchived: 2026-10-21\ndisposition: partial\nsource: docs/specs/0001-demo\nsource_revision: %s\nqa_task: task_02\nqa_report: qa/qa-report-2026-10-21.md\nqa_verdict: partial\nunproven:\n  - the operator records it in [`references/promoted.md`](references/promoted.md)\n  - the operator follows [`references/removed.md`](references/removed.md)\nadrs: []\nsources: []\nregeneration: []\npromoted:\n  - docs/references/promoted.md\npull_request: \"91\"\n---\n\n# Demo\n\nPreserved outcome.\n", spec.ArchiveRecordSchema, revision)
	historyWrite(t, filepath.Join(h.repo, "docs/history/specs/0001-demo.md"), record)
	h.revision = historyCommit(t, h.repo, "docs: earlier archive (#91)")
	gittest.Run(t, h.repo, "update-ref", "refs/remotes/origin/main", h.revision)
	return h, revision
}

func TestHistorySanitizePointsAPromotedRecordLinkAtItsCopyTranscript(t *testing.T) {
	t.Parallel()
	if !historyRemovalTestHome(t) {
		return
	}
	h, revision := promotedRecordLinksFixture(t, false)
	before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
	code, out, stderr := historyRun(t, h, "sanitize")
	want := "history sanitize plan: 1 unit(s) pending; 0 file(s) (0 bytes) leave docs/history\n" +
		"record-links: rewrites 2 link(s) in 1 Archive Record(s)\n" +
		"cites docs/history/specs/0001-demo.md:15 names docs/history/specs/references/promoted.md; rewrites to docs/references/promoted.md\n" +
		"cites docs/history/specs/0001-demo.md:16 names docs/history/specs/references/removed.md; rewrites to `docs/specs/0001-demo/references/removed.md@" + revision + "`\n" +
		"apply with: roundfix history sanitize --apply --batch <n> (needs the annotated tag history-full at or before HEAD)\n"
	if code != 0 || stderr != "" || out != want {
		t.Fatalf("exit=%d stderr=%q\nstdout:\n%s\nwant:\n%s", code, stderr, out, want)
	}
	baselineHistoryUnchanged(t, h, before, refs, index)
}

func TestBaselineUpdatePlansAPromotedRecordLinkAtItsCopy(t *testing.T) {
	t.Parallel()
	if !historyRemovalTestHome(t) {
		return
	}
	h, revision := promotedRecordLinksFixture(t, true)
	before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
	preview, out, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 3 || stderr != "" || preview.History.Status != "pending" || preview.State != "plan_ready" {
		t.Fatalf("preview=%d %s %s", code, out, stderr)
	}
	want := []baselineUpdateHistoryCitation{
		{Path: "docs/history/specs/0001-demo.md", Line: 15, Target: "docs/history/specs/references/promoted.md", Rewrite: "docs/references/promoted.md"},
		{Path: "docs/history/specs/0001-demo.md", Line: 16, Target: "docs/history/specs/references/removed.md", Rewrite: "`docs/specs/0001-demo/references/removed.md@" + revision + "`"},
	}
	assertBaselineCitationJSON(t, out, "citations", want)
	if !reflect.DeepEqual(preview.History.Citations, want) {
		t.Fatalf("citations=%+v", preview.History.Citations)
	}
	// Change only the promotion's Rewrite in the planned report. Holding all
	// other inputs constant proves this destination itself is digest-bound.
	plan := planBaselineUpdateHistory(t.Context(), h.repo, h.env)
	original, err := plan.digest("baseline")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Report.Citations) != 2 || plan.Report.Citations[0].Rewrite != want[0].Rewrite {
		t.Fatalf("plan=%+v", plan.Report)
	}
	plan.Report.Citations[0].Rewrite = "docs/references/other.md"
	changed, err := plan.digest("baseline")
	if err != nil || changed == original {
		t.Fatalf("digest omitted promotion: %q %q %v", original, changed, err)
	}
	baselineHistoryUnchanged(t, h, before, refs, index)
}
