package spec

// Suite: history conversion across repository text sources.
// Invariant: every listed citation is rewritten; excluded sources stay untouched.
// Boundary IN: tracked files in isolated temporary Git repositories.
// Boundary OUT: plan citations, file bytes and actual Apply counts.

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

// Keep each parallel test's Git reads in a private HOME, and register child
// cancellation and reaping before waiting so failure cannot strand the child.
func citationSourcesTestHome(t *testing.T) bool {
	t.Helper()
	const key = "ROUNDFIX_CITATION_SOURCES_TEST"
	if os.Getenv(key) == t.Name() {
		return true
	}
	ctx, cancel := context.WithCancel(t.Context())
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^"+regexp.QuoteMeta(t.Name())+"$", "-test.count=1")
	for _, entry := range gittest.IsolatedEnv() {
		if !strings.HasPrefix(entry, "HOME=") && !strings.HasPrefix(entry, key+"=") {
			cmd.Env = append(cmd.Env, entry)
		}
	}
	cmd.Env = append(cmd.Env, "HOME="+t.TempDir(), key+"="+t.Name())
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	waited := false
	t.Cleanup(func() {
		cancel()
		if cmd.Process != nil && !waited {
			_ = cmd.Wait()
		}
	})
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	err := cmd.Wait()
	waited = true
	if err != nil {
		t.Fatalf("isolated citation test: %v\n%s", err, output.String())
	}
	return false
}

func TestHistoryRewriteReachesEveryCitationSource(t *testing.T) {
	t.Parallel()
	if !citationSourcesTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	target := "docs/history/findings/x.md"
	historyEntryWrite(t, root, target, "# Retired\n")
	sources := map[string]string{
		"docs/adr/0001-current.md":      "`" + target + "`\n[old](../history/findings/x.md)\n",
		"docs/agents/guide.md":          "`../history/findings/x.md`\n",
		"docs/references/page.markdown": "[old](../history/findings/x.md)\n",
		"README.mdx":                    "[old](docs/history/findings/x.md)\n`./docs/history/findings/x.md`\n",
		"Makefile":                      "# " + target + "\n",
	}
	for path, content := range sources {
		historyEntryWrite(t, root, path, content)
	}
	revision := legacyCommit(t, root, "docs: citation sources")
	plan, err := PlanHistoryRewrite(t.Context(), root, "docs/specs", "docs/history/specs", []string{target}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Citations) != 7 || len(plan.Files) != 5 || plan.Count != 7 {
		t.Fatalf("plan=%+v", plan)
	}
	for _, c := range plan.Citations {
		if c.Rewrite == "" {
			t.Fatalf("unrewritten citation: %+v", c)
		}
	}
	count, files, err := plan.Apply(t.Context(), root)
	if err != nil || count != 7 || files != 5 {
		t.Fatalf("apply=%d,%d,%v", count, files, err)
	}
	for path := range sources {
		content, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(content), target+"@"+revision) {
			t.Fatalf("%s: %s", path, content)
		}
	}
	remaining, err := SourceHistoryCitations(t.Context(), root, []string{target})
	if err != nil || len(remaining) != 0 {
		t.Fatalf("remaining=%+v err=%v", remaining, err)
	}
}

func TestHistoryRewriteResolvesARelativeProsePath(t *testing.T) {
	t.Parallel()
	if !citationSourcesTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	target := "docs/history/findings/x.md"
	historyEntryWrite(t, root, target, "# Retired\n")
	source := "docs/history/guide.md"
	historyEntryWrite(t, root, source, "./findings/x.md\n`./findings/x.md`\n../history/findings/x.md\n`../history/findings/x.md`\nfindings/x.md\n")
	revision := legacyCommit(t, root, "docs: relative prose")
	plan, err := PlanHistoryRewrite(t.Context(), root, "docs/specs", "docs/history/specs", []string{target}, nil)
	want := strings.Repeat("`"+target+"@"+revision+"`\n", 4) + "findings/x.md\n"
	if err != nil || len(plan.Citations) != 4 || plan.Count != 4 || string(plan.Files[source]) != want {
		t.Fatalf("plan=%+v err=%v want=%s", plan, err, want)
	}
	destination := "docs/references/x.md"
	moved, err := PlanHistoryRewrite(t.Context(), root, "docs/specs", "docs/history/specs", []string{target}, map[string]string{target: destination})
	wantMoved := destination + "\n`" + destination + "`\n" + destination + "\n`" + destination + "`\nfindings/x.md\n"
	if err != nil || moved.Count != 4 || string(moved.Files[source]) != wantMoved {
		t.Fatalf("moved=%+v err=%v", moved, err)
	}

}

func TestHistoryRewriteWritesAPlainLocatorOutsideMarkdown(t *testing.T) {
	t.Parallel()
	if !citationSourcesTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	target := "docs/history/findings/x.md"
	historyEntryWrite(t, root, target, "# Retired\n")
	for _, path := range []string{"Makefile", "sample.go"} {
		historyEntryWrite(t, root, path, "// "+target+"\n// [old](./"+target+")\n// [unrelated](other.md)\n")
	}
	revision := legacyCommit(t, root, "docs: plain locators")
	plan, err := PlanHistoryRewrite(t.Context(), root, "docs/specs", "docs/history/specs", []string{target}, nil)
	if err != nil {
		t.Fatal(err)
	}
	locator := target + "@" + revision
	for _, c := range plan.Citations {
		if c.Rewrite != locator {
			t.Fatalf("citation=%+v", c)
		}
	}
	for _, path := range []string{"Makefile", "sample.go"} {
		want := "// " + locator + "\n// [old](" + locator + ")\n// [unrelated](other.md)\n"
		if string(plan.Files[path]) != want {
			t.Fatalf("%s=%s want=%s", path, plan.Files[path], want)
		}
	}
	if plan.Count != 4 || len(plan.Citations) != 4 {
		t.Fatalf("plan=%+v", plan)
	}
}

func TestCitationSourcesSkipSkillTreesTestdataAndBinaryFiles(t *testing.T) {
	t.Parallel()
	if !citationSourcesTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	target := "docs/history/findings/x.md"
	historyEntryWrite(t, root, target, "# Retired\n")
	excluded := []string{".agents/skills/demo/SKILL.md", "skills/demo/SKILL.md", "nested/testdata/input.md", "binary.dat"}
	for _, path := range excluded {
		historyEntryWrite(t, root, path, target+"\n")
	}
	binary := target + strings.Repeat(" ", 7999-len(target)) + "\x00\n"
	historyEntryWrite(t, root, "binary.dat", binary)
	historyEntryWrite(t, root, "README.md", target+"\n")
	historyEntryWrite(t, root, ".gitignore", "ignored.txt\n")
	// Git can represent a symlink even on Windows hosts without symlink privileges.
	if err := os.Symlink("README.md", filepath.Join(root, "link.md")); err != nil {
		if err := os.WriteFile(filepath.Join(root, "link.md"), []byte("README.md"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	legacyCommit(t, root, "docs: excluded sources")

	// Use README's existing blob for a tracked symlink index entry; its bytes are never scanned.
	blob := strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD:README.md"))
	gittest.Run(t, root, "update-index", "--cacheinfo", "120000,"+blob+",link.md")
	historyEntryWrite(t, root, "untracked.txt", target+"\n")
	historyEntryWrite(t, root, "ignored.txt", target+"\n")
	tracked, err := CitationSources(t.Context(), root, false)
	if err != nil {
		t.Fatal(err)
	}
	all, err := CitationSources(t.Context(), root, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range append(excluded, "link.md", "ignored.txt") {
		for _, got := range all {
			if got == path {
				t.Fatalf("excluded source %s listed: %v", path, all)
			}
		}
	}
	want := append(append([]string(nil), tracked...), "untracked.txt")
	// untracked.txt sorts last in this fixture.
	if !reflect.DeepEqual(all, want) {
		t.Fatalf("tracked=%v all=%v", tracked, all)
	}
	plan, err := PlanHistoryRewrite(t.Context(), root, "docs/specs", "docs/history/specs", []string{target}, nil)
	if err != nil || len(plan.Files) != 1 || plan.Files["README.md"] == nil {
		t.Fatalf("plan=%+v err=%v", plan, err)
	}

	if _, _, err := plan.Apply(t.Context(), root); err != nil {
		t.Fatal(err)
	}
	for _, path := range excluded {
		want := target + "\n"
		if path == "binary.dat" {
			want = binary
		}
		content, err := os.ReadFile(filepath.Join(root, path))
		if err != nil || string(content) != want {
			t.Fatalf("excluded %s changed: %q err=%v", path, content, err)
		}
	}
}
