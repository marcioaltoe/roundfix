package spec

// Suite: Archive Record link repair.
// Invariant: only broken relative links and self-links get reachable locators.
// Boundary IN: Markdown and isolated local Git histories.
// Boundary OUT: planned bytes and applied surviving records.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestMissingRelativeLinksFindsBrokenAndSelfLinks(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path := "docs/history/specs/demo.md"
	writeFile(t, filepath.Join(root, path), "record")
	writeFile(t, filepath.Join(root, "docs/references/exists.md"), "existing")
	content := []byte("[broken](missing.md#part)\n[self](demo.md)\n[good](../../references/exists.md)\n`[code](absent.md)`\n```md\n[fenced](no.md)\n```\n[url](https://example.test/no.md)\n`docs/no.md@" + strings.Repeat("a", 40) + "`\n![image](gone.png)\n[ref]: other.md \"title\"\n[anchor](#part)\n")
	links := MissingRelativeLinks(root, path, content)
	if len(links) != 4 {
		t.Fatalf("links=%+v", links)
	}
	for i, want := range []struct {
		target string
		line   int
		self   bool
	}{{"missing.md", 1, false}, {"demo.md", 2, true}, {"gone.png", 10, false}, {"other.md", 11, false}} {
		link := links[i]
		if link.Path != path || link.Target != "docs/history/specs/"+want.target || link.Line != want.line || link.Self != want.self {
			t.Fatalf("link=%+v want=%+v", link, want)
		}
	}
}

func TestPlanRecordLinksResolvesAgainstTheRecordThenItsSource(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, revision := gitLocatorRepo(t)
	for _, path := range []string{"docs/history/specs/first.md", "docs/specs/0001-demo/first.md", "docs/specs/0001-demo/second.md"} {
		writeFile(t, filepath.Join(root, path), "# Original\n")
	}
	revision = legacyCommit(t, root, "docs: sources")
	for _, path := range []string{"docs/history/specs/first.md", "docs/specs/0001-demo"} {
		gittest.Run(t, root, "rm", "-r", path)
	}
	record := ArchiveRecord{Spec: "0001-demo", Title: "Demo", Disposition: ArchiveNoQA, Source: "docs/specs/0001-demo", SourceRevision: revision, PullRequest: "42", Outcome: "[first](first.md) [second](second.md) [self](0001-demo.md) [absent](absent.md) [folder](../../specs/0001-demo). Unchanged `docs/specs/0001-demo`."}
	rendered, err := RenderArchiveRecord(record)
	if err != nil {
		t.Fatal(err)
	}
	path := "docs/history/specs/0001-demo.md"
	writeFile(t, filepath.Join(root, path), string(rendered))
	legacyCommit(t, root, "docs: archive")
	before, _ := os.ReadFile(filepath.Join(root, path))
	plan, err := PlanRecordLinks(t.Context(), root, "docs/history/specs", nil)
	if err != nil || plan.Links != 5 || len(plan.Files) != 1 {
		t.Fatalf("plan=%+v err=%v", plan, err)
	}
	for _, want := range []string{"`docs/history/specs/first.md@" + revision + "`", "`docs/specs/0001-demo/second.md@" + revision + "`", "`docs/specs/0001-demo@" + revision + "`", "`docs/history/specs/absent.md` (pull request #42)"} {
		if !strings.Contains(string(plan.Files[path]), want) {
			t.Fatalf("missing %s in %s", want, plan.Files[path])
		}
	}
	unchanged, _ := os.ReadFile(filepath.Join(root, path))
	if string(before) != string(unchanged) {
		t.Fatal("planning wrote record")
	}
	skipped, err := PlanRecordLinks(t.Context(), root, "docs/history/specs", map[string]bool{path: true})
	if err != nil || len(skipped.Files) != 0 {
		t.Fatalf("skip=%+v %v", skipped, err)
	}
	links, files, err := plan.Apply(t.Context(), root)
	if err != nil || links != 5 || files != 1 {
		t.Fatalf("apply=%d %d %v", links, files, err)
	}
	applied, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseArchiveRecord(applied)
	if err != nil || parsed.Source != record.Source || !strings.Contains(string(applied), "Unchanged `docs/specs/0001-demo`") {
		t.Fatalf("repair altered metadata or prose: %s err=%v", applied, err)
	}
	again, err := PlanRecordLinks(t.Context(), root, "docs/history/specs", nil)
	if err != nil || again.Links != 0 {
		t.Fatalf("second=%+v %v", again, err)
	}
}
