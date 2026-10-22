package speccheck

// Suite: promoted-copy fixes for surviving Archive Records.
// Invariant: the check names the same existing copy as Record Link Repair.
// Boundary IN: historical record metadata, isolated Git and a checkout without Git.
// Boundary OUT: finding text, severity, location and exact fix lines.

import (
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func TestHistoryRetainedNamesThePromotedCopyForARecordLink(t *testing.T) {
	t.Parallel()
	if !retainedTestHome(t) {
		return
	}
	for _, git := range []bool{true, false} {
		name := "outside Git"
		if git {
			name = "inside Git"
		}
		t.Run(name, func(t *testing.T) {
			root := retainedFixture(t, git)
			source := "docs/specs/0001-demo"
			revision := strings.Repeat("a", 40)
			if git {
				for _, name := range []string{"promoted.md", "removed.md"} {
					writeReceiptFixture(t, root, source+"/references/"+name, "# Source\n")
				}
				revision = retainedCommit(t, root)
				gittest.Run(t, root, "rm", "-r", source)
			}
			writeReceiptFixture(t, root, "docs/references/promoted.md", "# Source\n")
			content, err := spec.RenderArchiveRecord(spec.ArchiveRecord{
				Spec: "0001-demo", Title: "Demo", Disposition: spec.ArchiveNoQA, Source: source, SourceRevision: revision,
				Promoted: []string{"docs/references/promoted.md"}, Outcome: "Preserved outcome.",
				Unproven: []string{"the operator records it in [`references/promoted.md`](references/promoted.md)", "the operator follows [`references/removed.md`](references/removed.md)"},
			})
			if err != nil {
				t.Fatal(err)
			}
			path := "docs/history/specs/0001-demo.md"
			writeReceiptFixture(t, root, path, string(content))
			if git {
				retainedCommit(t, root)
			}
			findings := retainedFindings(retainedCheck(t, root))
			if len(findings) != 2 {
				t.Fatalf("findings=%+v", findings)
			}
			wantFix := []string{
				"Replace the link with ../../references/promoted.md, the copy the archive promoted to docs/references/promoted.md.",
				"Restore docs/history/specs/references/removed.md or replace its link; Git Locator resolution is skipped outside Git.",
			}
			if git {
				wantFix[1] = "Replace the link with `" + source + "/references/removed.md@" + revision + "`."
			}
			links := spec.MissingRelativeLinks(root, path, content)
			for i, finding := range findings {
				if finding.Fix != wantFix[i] || finding.Severity != SeverityError || finding.Summary != path+" links to a missing file or itself: "+links[i].Target || len(finding.Where) != 1 || finding.Where[0].Path != path || finding.Where[0].Line != links[i].Line {
					t.Fatalf("finding=%+v want fix=%s", finding, wantFix[i])
				}
			}
		})
	}
}
