package cli

// Suite: both history commands repair broken Archive Record links.
// Invariant: previews do not write; apply composes removal and link rewrites.
// Boundary IN: public commands, isolated homes and temporary Git repositories.
// Boundary OUT: output, JSON, digests, and surviving record bytes.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func recordLinksFixture(t *testing.T, baseline, removal bool) (historyFixture, string, string) {
	t.Helper()
	var h historyFixture
	if removal {
		h, _, _, _, _ = historyRemovalFixture(t, baseline)
	} else if baseline {
		h = historyFixture{repo: newBaselineUpdateRepository(t), home: t.TempDir()}
		h.env = commandEnvironmentForTest(t)
		h.env.homeDir, h.env.workDir = h.home, h.repo
		h.env.environ = withEnvValue(h.env.environ, "HOME", h.home)
		mustMkdir(t, filepath.Join(h.repo, "docs/specs"))
	} else {
		h = newHistoryFixture(t, false)
	}
	historyWrite(t, filepath.Join(h.repo, "docs/specs/0003-live/_prd.md"), "# Live\n")
	source := "docs/specs/0001-old"
	historyWrite(t, filepath.Join(h.repo, source, "file.md"), "# Original\n")
	revision := historyCommit(t, h.repo, "docs: original")
	gittest.Run(t, h.repo, "rm", "-r", source)
	outcome := "[broken](file.md) and [self](0001-old.md)"
	if removal {
		outcome += " and [retired](../findings/reduced.md)"
	}
	outcome += "."
	record := spec.ArchiveRecord{Spec: "0001-old", Title: "Old", Disposition: spec.ArchiveNoQA, Source: source, SourceRevision: revision, Outcome: outcome}
	content, err := spec.RenderArchiveRecord(record)
	if err != nil {
		t.Fatal(err)
	}
	path := "docs/history/specs/0001-old.md"
	historyWrite(t, filepath.Join(h.repo, path), string(content))
	// The repair unit must skip this record, because the earlier unit removes it.
	record.Spec, record.Source, record.Disposition, record.SupersededBy = "0002-unused", "docs/specs/0002-unused", spec.ArchiveSuperseded, "0001-old"
	record.Outcome = "[broken](never.md)."
	content, err = spec.RenderArchiveRecord(record)
	if err != nil {
		t.Fatal(err)
	}
	if removal {
		historyWrite(t, filepath.Join(h.repo, "docs/history/specs/0002-unused.md"), string(content))
	}
	h.revision = historyCommit(t, h.repo, "docs: archive")
	return h, path, revision
}
func assertRecordLinksApplied(t *testing.T, h historyFixture, path, revision string, removal bool) {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(h.repo, path))
	if err != nil {
		t.Fatal(err)
	}
	for _, target := range []string{"docs/specs/0001-old/file.md", "docs/specs/0001-old"} {
		if !strings.Contains(string(content), "`"+target+"@"+revision+"`") {
			t.Fatalf("missing locator %s in %s", target, content)
		}
	}
	if removal {
		if !strings.Contains(string(content), "retired (`docs/findings/original.md@") {
			t.Fatalf("removal rewrite lost: %s", content)
		}
		assertPathMissing(t, filepath.Join(h.repo, "docs/history/specs/0002-unused.md"))
	}
}

func TestHistorySanitizeRewritesBrokenRecordLinks(t *testing.T) {
	t.Parallel()
	if !historyRemovalTestHome(t) {
		return
	}
	for _, removal := range []bool{false, true} {
		t.Run(map[bool]string{false: "only links", true: "removals and links"}[removal], func(t *testing.T) {
			h, path, revision := recordLinksFixture(t, false, removal)
			historyTag(t, h)
			before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
			code, out, stderr := historyRun(t, h, "sanitize")
			if code != 0 || stderr != "" || !strings.Contains(out, "record-links: rewrites 2 link(s) in 1 Archive Record(s)\n") {
				t.Fatalf("preview=%d %s %s", code, out, stderr)
			}
			if !removal && !strings.Contains(out, "0 file(s) (0 bytes) leave docs/history") {
				t.Fatalf("repair counted as removal: %s", out)
			}
			baselineHistoryUnchanged(t, h, before, refs, index)
			code, out, stderr = historyRun(t, h, "sanitize", "--apply", "--batch", "10")
			count := "2 citation(s) in 1 file(s)"
			if removal {
				count = "4 citation(s) in 2 file(s)"
			}
			if code != 0 || stderr != "" || !strings.Contains(out, "rewrote "+count) {
				t.Fatalf("apply=%d %s %s", code, out, stderr)
			}
			assertRecordLinksApplied(t, h, path, revision, removal)
			code, out, stderr = historyRun(t, h, "sanitize")
			if code != 0 || stderr != "" || !strings.Contains(out, "nothing pending") {
				t.Fatalf("second=%d %s %s", code, out, stderr)
			}
		})
	}
}
func TestBaselineUpdateRewritesBrokenRecordLinks(t *testing.T) {
	t.Parallel()
	if !historyRemovalTestHome(t) {
		return
	}
	for _, removal := range []bool{false, true} {
		t.Run(map[bool]string{false: "only links", true: "removals and links"}[removal], func(t *testing.T) {
			h, path, revision := recordLinksFixture(t, true, removal)
			before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
			_, textOut, textErr, textCode := baselineHistoryRun(t, h, "text")
			if textCode != 3 || textErr != "" || !strings.Contains(textOut, "- record-links: rewrites 2 link(s) in 1 Archive Record(s)\n") {
				t.Fatalf("text preview=%d %s %s", textCode, textOut, textErr)
			}
			if !removal && !strings.Contains(textOut, "0 file(s), 0 bytes leave docs/history") {
				t.Fatalf("text counted repair as removal: %s", textOut)
			}
			preview, out, stderr, code := baselineHistoryRun(t, h, "json")
			if code != 3 || stderr != "" || preview.History.Status != "pending" {
				t.Fatalf("preview=%d %s %s", code, out, stderr)
			}
			u := preview.History.Units[len(preview.History.Units)-1]
			if u.Unit != "record-links" || u.Action != "rewrite" || u.Files != 1 || u.Links != 2 || strings.Contains(out, "bytesAfter") {
				t.Fatalf("rewrite=%+v stdout=%s", u, out)
			}
			baselineHistoryUnchanged(t, h, before, refs, index)
			// A changed rewrite destination invalidates confirmation.
			p := planBaselineUpdateHistory(t.Context(), h.repo, h.env)
			digest, err := p.digest("baseline")
			if err != nil {
				t.Fatal(err)
			}
			p.Report.Citations[len(p.Report.Citations)-1].Rewrite = "changed"
			changed, err := p.digest("baseline")
			if err != nil || changed == digest {
				t.Fatal("digest omitted rewrite")
			}
			result, out, stderr, code := baselineHistoryRun(t, h, "json", "--confirm-plan", preview.PlanDigest)
			links, files, removed := 2, 1, 0
			if removal {
				links, files, removed = 4, 2, 2
			}
			a := result.History.Applied
			if code != 0 || stderr != "" || a == nil || a.RewrittenCitations != links || a.RewrittenFiles != files || a.Removed != removed || a.Reduced != 0 {
				t.Fatalf("apply=%d %s %s", code, out, stderr)
			}
			assertRecordLinksApplied(t, h, path, revision, removal)
			again, out, stderr, code := baselineHistoryRun(t, h, "json")
			if code != 0 || stderr != "" || again.History.Status != "current" {
				t.Fatalf("second=%d %s %s", code, out, stderr)
			}
		})
	}
}
