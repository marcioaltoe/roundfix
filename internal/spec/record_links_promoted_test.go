package spec

// Suite: promoted Archive Record link repair.
// Invariant: a listed surviving copy wins; other missing targets keep locators.
// Boundary IN: committed records from an earlier archive and private Git history.
// Boundary OUT: plan citations, planned bytes and applied record bytes.

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestPlanRecordLinksPointsAPromotedTargetAtItsCopy(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	source := "docs/specs/0001-demo"
	for _, name := range []string{"promoted.md", "removed.md"} {
		writeFile(t, filepath.Join(root, source, "references", name), "# Original\n")
	}
	revision := legacyCommit(t, root, "docs: source")
	gittest.Run(t, root, "rm", "-r", source)
	writeFile(t, filepath.Join(root, "docs/references/promoted.md"), "# Original\n")
	record := ArchiveRecord{Spec: "0001-demo", Title: "Demo", Disposition: ArchivePartial, Source: source, SourceRevision: revision, Promoted: []string{"docs/references/promoted.md"}, Unproven: []string{
		"the operator records it in [`references/promoted.md`](references/promoted.md)",
		"the operator follows [`references/removed.md`](references/removed.md)",
	}, Outcome: "Preserved outcome."}
	content, err := RenderArchiveRecord(record)
	if err != nil {
		t.Fatal(err)
	}
	path := "docs/history/specs/0001-demo.md"
	writeFile(t, filepath.Join(root, path), string(content))
	legacyCommit(t, root, "docs: earlier archive")
	plan, err := PlanRecordLinks(t.Context(), root, "docs/history/specs", nil)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Links != 2 || len(plan.Files) != 1 || len(plan.Citations) != 2 {
		t.Fatalf("plan=%+v", plan)
	}
	if plan.Citations[0].Rewrite != "docs/references/promoted.md" || plan.Citations[1].Rewrite != "`"+source+"/references/removed.md@"+revision+"`" {
		t.Fatalf("citations=%+v", plan.Citations)
	}
	for _, want := range []string{"[`references/promoted.md`](../../references/promoted.md)", "`" + source + "/references/removed.md@" + revision + "`"} {
		if !strings.Contains(string(plan.Files[path]), want) {
			t.Fatalf("missing %q in %s", want, plan.Files[path])
		}
	}
	before, err := os.ReadFile(filepath.Join(root, path))
	if err != nil || !bytes.Equal(before, content) {
		t.Fatalf("planning wrote record: %v", err)
	}
	links, files, err := plan.Apply(t.Context(), root)
	if err != nil || links != 2 || files != 1 {
		t.Fatalf("apply=%d %d %v", links, files, err)
	}
	applied, err := os.ReadFile(filepath.Join(root, path))
	if err != nil || !bytes.Equal(applied, plan.Files[path]) {
		t.Fatalf("applied bytes differ: %v", err)
	}
	again, err := PlanRecordLinks(t.Context(), root, "docs/history/specs", nil)
	if err != nil || again.Links != 0 || len(again.Citations) != 0 {
		t.Fatalf("again=%+v err=%v", again, err)
	}
}

func TestRecordLinkPromotionRequiresTheListedCopy(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"listed copy", "absent copy", "unlisted copy", "outside source", "self link", "directory copy"} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			copyPath := "docs/references/promoted.md"
			record := ArchiveRecord{Source: "docs/specs/0001-demo", Promoted: []string{copyPath}}
			link := MissingLink{Path: "docs/history/specs/0001-demo.md", Destination: "references/promoted.md", Target: "docs/history/specs/references/promoted.md"}
			if name != "absent copy" && name != "directory copy" {
				writeFile(t, filepath.Join(root, copyPath), "copy")
			}
			switch name {
			case "unlisted copy":
				record.Promoted = nil
			case "outside source":
				link.Destination = "../0002-other/promoted.md"
			case "self link":
				link.Self = true
			case "directory copy":
				if err := os.MkdirAll(filepath.Join(root, copyPath), 0755); err != nil {
					t.Fatal(err)
				}
			}
			got, ok := RecordLinkPromotion(root, link, record)
			if name == "listed copy" {
				if !ok || got != copyPath {
					t.Fatalf("promotion=%q %v", got, ok)
				}
			} else if ok || got != "" {
				t.Fatalf("unexpected promotion=%q %v", got, ok)
			}
		})
	}
}
