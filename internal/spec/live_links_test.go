package spec

// Suite: live-page missing-link repair.
// Invariant: the history gate's live pages receive the same Git Locators.
// Boundary IN: Markdown in temporary repositories with removed source files.
// Boundary OUT: frozen plans, applied link spans and repeat planning.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestLiveLinksPlanRewritesMissingLinksInLiveDocumentsAndGuides(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	target := "docs/specs/old/_techspec.md"
	writeFile(t, filepath.Join(root, target), "# Original\n")
	revision := legacyCommit(t, root, "docs: original techspec")
	gittest.Run(t, root, "rm", "-r", "docs/specs/old")
	paths := []string{"docs/references/live.md", "docs/user-guide/guide.md", "docs/backlog/open.md", "docs/findings/pending.md", "plans/live/_prd.md"}
	for _, path := range paths {
		relative, err := filepath.Rel(filepath.Dir(path), target)
		if err != nil {
			t.Fatal(err)
		}
		writeFile(t, filepath.Join(root, path), "# Live\n\n[Current behavior]("+filepath.ToSlash(relative)+"#current-behavior)\n")
	}
	legacyCommit(t, root, "docs: retain live pages")
	plan, err := PlanLiveLinks(t.Context(), root, "plans", "plans/_archived", nil)
	if err != nil || plan.Links != len(paths) || len(plan.Files) != len(paths) || len(plan.Citations) != len(paths) {
		t.Fatalf("plan=%+v err=%v", plan, err)
	}
	locator := "`" + target + "@" + revision + "`"
	for _, path := range paths {
		want := "# Live\n\nCurrent behavior (" + locator + ")\n"
		if string(plan.Files[path]) != want {
			t.Fatalf("planned %s=%q want=%q", path, plan.Files[path], want)
		}
		before, err := os.ReadFile(filepath.Join(root, path))
		if err != nil || !strings.Contains(string(before), "[Current behavior]") {
			t.Fatalf("plan wrote %s: %s %v", path, before, err)
		}
	}
	for _, c := range plan.Citations {
		if c.Line != 3 || c.Target != target || c.Rewrite != locator {
			t.Fatalf("citation=%+v", c)
		}
	}
	links, files, err := plan.Apply(t.Context(), root)
	if err != nil || links != len(paths) || files != len(paths) {
		t.Fatalf("apply=%d %d %v", links, files, err)
	}
	for _, path := range paths {
		content, err := os.ReadFile(filepath.Join(root, path))
		if err != nil || string(content) != string(plan.Files[path]) {
			t.Fatalf("apply %s=%s %v", path, content, err)
		}
	}
	again, err := PlanLiveLinks(t.Context(), root, "plans", "plans/_archived", nil)
	if err != nil || len(again.Files) != 0 {
		t.Fatalf("second=%+v err=%v", again, err)
	}
}

func TestLiveLinksPlanLeavesRecordsCodeAndRemovedFiles(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, revision := gitLocatorRepo(t)
	untouched := "# Live\n\n`[code](absent.md)`\n```md\n[fence](missing.md)\n```\n[url](https://example.test/missing.md)\n[anchor](#part)\n[self](live.md)\n`docs/specs/demo/_prd.md@" + revision + "`\n"
	for path, content := range map[string]string{
		"docs/references/live.md":    untouched,
		"docs/references/removed.md": "[broken](missing.md)",
		"docs/adr/accepted.md":       "[not a live page](missing.md)",
		"docs/history/specs/old.md":  "[record](missing.md)",
		"plans/_archived/old.md":     "[record](missing.md)",
		"docs/references/plain.txt":  "[not Markdown](missing.md)",
	} {
		writeFile(t, filepath.Join(root, path), content)
	}
	external := filepath.Join(t.TempDir(), "outside.md")
	writeFile(t, external, "[external](missing.md)")
	if err := os.Symlink(external, filepath.Join(root, "docs/references/link.md")); err != nil {
		t.Fatal(err)
	}
	plan, err := PlanLiveLinks(t.Context(), root, "plans", "plans/_archived", map[string]bool{"docs/references/removed.md": true})
	if err != nil || plan.Links != 0 || len(plan.Files) != 0 || len(plan.Citations) != 0 {
		t.Fatalf("plan=%+v err=%v", plan, err)
	}
	if _, _, err := plan.Apply(t.Context(), root); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(filepath.Join(root, "docs/references/live.md"))
	if err != nil || string(after) != untouched {
		t.Fatalf("untouched=%q err=%v", after, err)
	}
}
