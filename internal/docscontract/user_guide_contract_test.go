//go:build docscontract

package docscontract

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/cli"
)

// brokenLinks returns file:line: destination for missing relative Markdown links.
func brokenLinks(repoRoot, file, content string) []string {
	links := regexp.MustCompile(`\[[^\]]*\]\((<[^>]*>[^)]*|[^)]*)\)`)
	scheme := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)
	var broken []string
	var fence byte
	var fenceLength int
	for index, line := range strings.Split(content, "\n") {
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
			if _, err := os.Stat(filepath.Join(repoRoot, filepath.Dir(file), filepath.FromSlash(destination))); err != nil {
				broken = append(broken, fmt.Sprintf("%s:%d: %s", file, index+1, destination))
			}
		}
	}
	return broken
}

func userGuideFiles(t *testing.T, includeREADME bool) []string {
	t.Helper()
	root := baselineDocumentationRepoRoot()
	files, err := filepath.Glob(filepath.Join(root, "docs", "user-guide", "*.md"))
	if err != nil {
		t.Fatalf("list user guide: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no user guide files found")
	}
	companions, err := filepath.Glob(filepath.Join(root, "docs", "user-guide", "commands", "*.md"))
	if err != nil {
		t.Fatalf("list command guides: %v", err)
	}
	if includeREADME {
		files = append(files, companions...)
		files = append(files, filepath.Join(root, "README.md"))
	}
	return files
}

func TestEveryCommandIsNamedInTheUserGuide(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := cli.Run([]string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("root help exited %d: %s", code, stderr.String())
	}
	var text strings.Builder
	for _, file := range userGuideFiles(t, false) {
		text.WriteString(readContractDocument(t, file))
		text.WriteByte('\n')
	}
	paths := commandPaths(stdout.String())
	if !slices.Contains(paths, "runs causes") {
		t.Fatal("root help omitted runs causes")
	}
	if missing := undocumentedCommands(paths, text.String()); len(missing) > 0 {
		t.Fatalf("user guide does not name command paths: %v", missing)
	}
}

func TestUserGuideLinksResolve(t *testing.T) {
	root := baselineDocumentationRepoRoot()
	for _, file := range userGuideFiles(t, true) {
		relative, err := filepath.Rel(root, file)
		if err != nil {
			t.Fatalf("relative guide path: %v", err)
		}
		for _, broken := range brokenLinks(root, relative, mustRead(t, file)) {
			t.Error(broken)
		}
	}
}

func TestABrokenUserGuideLinkIsReported(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "guide"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "live.md"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	content := "[dead](missing.md)\n[live](../live.md)\n[URL](https://example.com)\n```markdown\n[fenced](also-missing.md)\n```\n"
	want := []string{"guide/index.md:1: missing.md"}
	if got := brokenLinks(root, "guide/index.md", content); !reflect.DeepEqual(got, want) {
		t.Fatalf("brokenLinks() = %v, want %v", got, want)
	}
}

func TestBrokenLinksHandlesDestinationDecorationsAndFences(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"live.md", "live file.md"} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	content := "[live](<live file.md?view=1#section> \"title\") [plain](live.md?view=1#section \"title\")\n[anchor](#section)\n[mail](mailto:user@example.com)\n~~~markdown\n[ignored](missing.md)\n```\n[still ignored](missing.md)\n~~~\n[dead](<missing.md?view=1#section> \"title\")\n"
	want := []string{"index.md:9: missing.md"}
	if got := brokenLinks(root, "index.md", content); !reflect.DeepEqual(got, want) {
		t.Fatalf("brokenLinks() = %v, want %v", got, want)
	}
}

func TestUserGuideNamesNoRefusedQAFlag(t *testing.T) {
	for _, file := range userGuideFiles(t, true) {
		for index, line := range strings.Split(mustRead(t, file), "\n") {
			for _, suffix := range strings.Split(line, "--qa")[1:] {
				if !strings.HasPrefix(suffix, "-override") {
					t.Errorf("%s:%d: refused QA flag", file, index+1)
				}
			}
		}
	}
}
