package skills

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const roundfixSkillEntryByteBudget = 20_000

func roundfixLayoutSources() map[string]fs.FS {
	return map[string]fs.FS{
		"canonical": os.DirFS("../.agents/skills"),
		"embedded":  embedded,
	}
}

func roundfixReferenceIndexLinks(t *testing.T, source fs.FS) []string {
	t.Helper()
	data, err := fs.ReadFile(source, "roundfix/SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	const begin = "<!-- roundfix:reference-index:begin -->"
	const end = "<!-- roundfix:reference-index:end -->"
	text := string(data)
	if strings.Count(text, begin) != 1 || strings.Count(text, end) != 1 {
		t.Fatal("entry must have exactly one reference index marker pair")
	}
	_, after, _ := strings.Cut(text, begin)
	index, _, ok := strings.Cut(after, end)
	if !ok {
		t.Fatal("reference index markers are out of order")
	}
	matches := regexp.MustCompile(`(?m)^\| \[[^\]]+\]\(([^)]+)\) \|`).FindAllStringSubmatch(index, -1)
	if len(matches) == 0 {
		t.Fatal("reference index has no linked table rows")
	}
	links := make([]string, 0, len(matches))
	for _, match := range matches {
		links = append(links, match[1])
	}
	return links
}

func TestRoundfixSkillIndexNamesEveryReference(t *testing.T) {
	t.Parallel()
	for name, source := range roundfixLayoutSources() {
		t.Run(name, func(t *testing.T) {
			entries, err := fs.ReadDir(source, "roundfix/references")
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 18 {
				t.Fatalf("reference files = %d, want 18 command families", len(entries))
			}
			counts := make(map[string]int)
			for _, link := range roundfixReferenceIndexLinks(t, source) {
				counts[link]++
			}
			for _, entry := range entries {
				link := "references/" + entry.Name()
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
					t.Errorf("unexpected reference artifact %s", link)
				}
				if counts[link] != 1 {
					t.Errorf("index rows for %s = %d, want 1", link, counts[link])
				}
				delete(counts, link)
			}
			for link := range counts {
				t.Errorf("index names absent reference %s", link)
			}
		})
	}
}

func TestRoundfixSkillIndexLinksResolve(t *testing.T) {
	t.Parallel()
	for name, source := range roundfixLayoutSources() {
		t.Run(name, func(t *testing.T) {
			for _, link := range roundfixReferenceIndexLinks(t, source) {
				if !fs.ValidPath(link) || !strings.HasPrefix(link, "references/") {
					t.Errorf("invalid relative reference link %q", link)
					continue
				}
				info, err := fs.Stat(source, "roundfix/"+link)
				if err != nil {
					t.Errorf("resolve %s: %v", link, err)
				} else if !info.Mode().IsRegular() {
					t.Errorf("%s does not resolve to a regular file", link)
				}
			}
		})
	}
}

func TestRoundfixSkillEntryFileStaysWithinItsBudget(t *testing.T) {
	t.Parallel()
	for name, source := range roundfixLayoutSources() {
		t.Run(name, func(t *testing.T) {
			data, err := fs.ReadFile(source, "roundfix/SKILL.md")
			if err != nil {
				t.Fatal(err)
			}
			if len(data) > roundfixSkillEntryByteBudget {
				t.Errorf("entry bytes = %d, budget = %d", len(data), roundfixSkillEntryByteBudget)
			}
		})
	}
}

func TestInstallWritesEveryRoundfixReference(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if _, err := Install(t.Context(), InstallRequest{
		Target: "codex", TargetDirs: map[string]string{"codex": dir},
	}); err != nil {
		t.Fatal(err)
	}
	source := os.DirFS("../.agents/skills")
	links := roundfixReferenceIndexLinks(t, source)
	for _, link := range append(links, "SKILL.md") {
		want, err := fs.ReadFile(source, "roundfix/"+link)
		if err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(filepath.Join(dir, "roundfix", filepath.FromSlash(link)))
		if err != nil {
			t.Errorf("read installed %s: %v", link, err)
			continue
		}
		if !bytes.Equal(got, want) {
			t.Errorf("installed %s differs from canonical bytes", link)
		}
	}
}
