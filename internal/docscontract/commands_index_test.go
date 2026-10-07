//go:build docscontract

//verify:always

package docscontract

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestCommandIndexNamesEveryCommandFile(t *testing.T) {
	root := baselineDocumentationRepoRoot()
	indexPath := filepath.Join(root, "docs", "user-guide", "commands.md")
	content := mustRead(t, indexPath)

	const begin = "<!-- roundfix:command-index:begin -->"
	const end = "<!-- roundfix:command-index:end -->"
	start := strings.Index(content, begin)
	finish := strings.Index(content, end)
	if start < 0 || finish < 0 || finish < start {
		t.Fatalf("command index markers are missing or out of order")
	}
	index := content[start+len(begin) : finish]

	wantPaths, err := filepath.Glob(filepath.Join(root, "docs", "user-guide", "commands", "*.md"))
	if err != nil {
		t.Fatalf("list command guides: %v", err)
	}
	want := make([]string, 0, len(wantPaths))
	for _, path := range wantPaths {
		want = append(want, filepath.Base(path))
	}
	slices.Sort(want)

	row := regexp.MustCompile(`\|\s*` + "`" + `([^` + "`" + `]+)` + "`" + `\s*\|\s*\[[^]]+\]\(commands/([^)]*\.md)\)\s*\|`)
	var got []string
	for _, match := range row.FindAllStringSubmatch(index, -1) {
		if match[1]+".md" != match[2] {
			t.Fatalf("index command %q points at %q", match[1], match[2])
		}
		got = append(got, match[2])
	}
	slices.Sort(got)
	if !slices.Equal(got, want) {
		t.Fatalf("command index files = %v, want %v", got, want)
	}
}

func TestCommandReferenceLinksResolve(t *testing.T) {
	root := baselineDocumentationRepoRoot()
	files := []string{filepath.Join(root, "docs", "user-guide", "commands.md"), filepath.Join(root, "docs", "user-guide", "usage.md")}
	companions, err := filepath.Glob(filepath.Join(root, "docs", "user-guide", "commands", "*.md"))
	if err != nil {
		t.Fatalf("list command guides: %v", err)
	}
	files = append(files, companions...)

	for _, file := range files {
		relative, err := filepath.Rel(root, file)
		if err != nil {
			t.Fatalf("relative path for %s: %v", file, err)
		}
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", relative, err)
		}
		if broken := brokenUserGuideLinks(root, relative, string(content)); len(broken) > 0 {
			t.Errorf("broken links in %s:\n%s", relative, strings.Join(broken, "\n"))
		}
	}
}

func brokenUserGuideLinks(root, file, content string) []string {
	links := regexp.MustCompile(`\[[^\]]*\]\((<[^>]*>[^)]*|[^)]*)\)`)
	scheme := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)
	var broken []string
	var fence byte
	var fenceLength int
	for lineNumber, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) >= 3 && (trimmed[0] == '`' || trimmed[0] == '~') {
			length := len(trimmed) - len(strings.TrimLeft(trimmed, string(trimmed[0])))
			if length >= 3 {
				if fence == 0 {
					fence, fenceLength = trimmed[0], length
				} else if trimmed[0] == fence && length >= fenceLength && strings.TrimSpace(trimmed[length:]) == "" {
					fence, fenceLength = 0, 0
				}
				continue
			}
		}
		if fence != 0 {
			continue
		}
		for _, match := range links.FindAllStringSubmatch(line, -1) {
			destination := strings.TrimSpace(match[1])
			if strings.HasPrefix(destination, "<") {
				end := strings.Index(destination, ">")
				destination = destination[1:end]
			} else if fields := strings.Fields(destination); len(fields) > 0 {
				destination = fields[0]
			}
			if destination == "" || strings.HasPrefix(destination, "#") || scheme.MatchString(destination) {
				continue
			}
			destination = strings.SplitN(strings.SplitN(destination, "#", 2)[0], "?", 2)[0]
			resolved := filepath.Clean(filepath.Join(filepath.Dir(file), filepath.FromSlash(destination)))
			if !strings.HasPrefix(resolved, filepath.FromSlash("docs/user-guide/")) && resolved != filepath.FromSlash("docs/user-guide") {
				continue
			}
			if _, err := os.Stat(filepath.Join(root, resolved)); err != nil {
				broken = append(broken, file+":"+itoa(lineNumber+1)+": "+destination)
			}
		}
	}
	return broken
}

func itoa(value int) string {
	return strconv.Itoa(value)
}
