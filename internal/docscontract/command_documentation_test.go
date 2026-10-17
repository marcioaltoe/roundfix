//go:build docscontract

//verify:always

package docscontract

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/cli"
	"roundfix/internal/secretscan"
)

// commandPaths returns command paths in root help's first-seen order.
func commandPaths(help string) []string {
	var paths []string
	seen := make(map[string]bool)
	add := func(path string) {
		if !seen[path] {
			seen[path] = true
			paths = append(paths, path)
		}
	}
	for _, line := range strings.Split(help, "\n") {
		if !strings.HasPrefix(line, "  roundfix ") {
			continue
		}
		var words []string
		for _, token := range strings.Fields(strings.TrimPrefix(line, "  roundfix ")) {
			if strings.HasPrefix(token, "[") || strings.HasPrefix(token, "--") || strings.HasPrefix(token, "(") || strings.HasPrefix(token, "<") {
				if len(words) > 0 && strings.HasPrefix(token, "<") && strings.HasSuffix(token, ">") && strings.Contains(token, "|") {
					for _, alternative := range strings.Split(token[1:len(token)-1], "|") {
						add(strings.Join(words, " ") + " " + alternative)
					}
					words = nil
				}
				break
			}
			words = append(words, token)
		}
		if len(words) > 0 {
			add(strings.Join(words, " "))
		}
	}
	return paths
}

func undocumentedCommands(paths []string, text string) []string {
	var missing []string
	for _, path := range paths {
		if !strings.Contains(text, "roundfix "+path) {
			missing = append(missing, path)
		}
	}
	return missing
}

func TestCommandPathsReadSubcommandsAndAlternatives(t *testing.T) {
	help := `Usage:
  roundfix window <set|show|clear>
  roundfix archive <slug>
  roundfix baseline capabilities check [--profile <id>]
  roundfix archive <slug> [--qa-override]
  roundfix --help
roundfix ignored
  roundfix events --run <id>
  roundfix doctor (read-only)
`
	want := []string{"window set", "window show", "window clear", "archive", "baseline capabilities check", "events", "doctor"}
	if got := commandPaths(help); !reflect.DeepEqual(got, want) {
		t.Fatalf("commandPaths() = %v, want %v", got, want)
	}
}

func TestEveryCommandIsNamedInTheRoundfixSkill(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := cli.Run([]string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("root help exited %d: %s", code, stderr.String())
	}
	paths := commandPaths(stdout.String())
	if len(paths) < 40 {
		t.Fatalf("root help yielded %d command paths, want at least 40: %v", len(paths), paths)
	}
	for _, anchor := range []string{"deliver retry", "window clear", "baseline capabilities check", "runs causes"} {
		if !slices.Contains(paths, anchor) {
			t.Fatalf("root help did not yield anchor command %q: %v", anchor, paths)
		}
	}
	text := readContractDocument(t, filepath.Join(baselineDocumentationRepoRoot(), ".agents/skills/roundfix/SKILL.md"))
	if missing := undocumentedCommands(paths, text); len(missing) > 0 {
		t.Fatalf("Roundfix skill does not name command paths: %v", missing)
	}
}

func TestAnUndocumentedCommandIsReported(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := cli.Run([]string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("root help exited %d: %s", code, stderr.String())
	}
	text := readContractDocument(t, filepath.Join(baselineDocumentationRepoRoot(), ".agents/skills/roundfix/SKILL.md"))
	text = strings.ReplaceAll(text, "roundfix reopen", "")
	if got, want := undocumentedCommands(commandPaths(stdout.String()), text), []string{"reopen"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("undocumentedCommands() = %v, want %v", got, want)
	}
}

func TestTheQASecretRefusalIsDocumentedAsTheDaemonPrintsIt(t *testing.T) {
	t.Parallel()

	root := baselineDocumentationRepoRoot()
	documented := qaSecretRefusalDocumentation(t, root)
	if err := checkQASecretRefusalDocumentation(qaSecretRefusalDocuments(t, root), documented); err != nil {
		t.Fatal(err)
	}
}

func TestTheQASecretRefusalContractRejectsTheParenthesizedForm(t *testing.T) {
	t.Parallel()

	root := baselineDocumentationRepoRoot()
	documented := qaSecretRefusalDocumentation(t, root)
	parenthesized := strings.Replace(documented, ":<line>: <kind>", ":<line> (<kind>)", 1)
	err := checkQASecretRefusalDocumentation([]qaSecretRefusalDocument{{
		path:    "docs/user-guide/commands/implement.md",
		content: parenthesized,
	}}, documented)
	if err == nil {
		t.Fatal("parenthesized QA secret refusal was accepted")
	}
	if !strings.Contains(err.Error(), "docs/user-guide/commands/implement.md") || !strings.Contains(err.Error(), parenthesized) {
		t.Fatalf("parenthesized QA secret refusal failure does not name the old sentence: %v", err)
	}
}

type qaSecretRefusalDocument struct {
	path    string
	content string
}

func qaSecretRefusalDocumentation(t *testing.T, root string) string {
	t.Helper()
	const prefix = "QA evidence carries a secret"
	var literals []string
	if err := filepath.WalkDir(filepath.Join(root, "internal", "daemon"), func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(file) != ".go" || strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}
		source, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read %s: %w", file, err)
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), file, source, 0)
		if err != nil {
			return fmt.Errorf("parse %s: %w", file, err)
		}
		ast.Inspect(parsed, func(node ast.Node) bool {
			literal, ok := node.(*ast.BasicLit)
			if !ok || literal.Kind != token.STRING {
				return true
			}
			value, err := strconv.Unquote(literal.Value)
			if err != nil {
				t.Fatalf("unquote %s: %v", file, err)
			}
			if strings.HasPrefix(value, prefix) {
				literals = append(literals, value)
			}
			return true
		})
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(literals) != 1 {
		t.Fatalf("found %d Daemon QA secret refusal literals, want exactly one", len(literals))
	}

	sample := secretscan.Match{Path: "<path>", Line: 1, Kind: "<kind>"}.String()
	documented, ok := strings.CutPrefix(strings.Replace(sample, "<path>:1:", "<path>:<line>:", 1), "<path>:<line>:")
	if !ok {
		t.Fatalf("Secret Match.String() did not produce the expected location: %q", sample)
	}
	return literals[0] + "<path>:<line>:" + documented
}

func qaSecretRefusalDocuments(t *testing.T, root string) []qaSecretRefusalDocument {
	t.Helper()
	var documents []qaSecretRefusalDocument
	directories := []string{
		"docs/user-guide",
		"docs/agents",
		".agents/skills",
		"skills",
	}
	for _, directory := range directories {
		if err := filepath.WalkDir(filepath.Join(root, directory), func(file string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || filepath.Ext(file) != ".md" {
				return nil
			}
			content, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read %s: %w", file, err)
			}
			relative, err := filepath.Rel(root, file)
			if err != nil {
				return fmt.Errorf("relativize %s: %w", file, err)
			}
			documents = append(documents, qaSecretRefusalDocument{
				path:    filepath.ToSlash(relative),
				content: string(content),
			})
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	content, err := os.ReadFile(filepath.Join(root, "CONTEXT.md"))
	if err != nil {
		t.Fatalf("read CONTEXT.md: %v", err)
	}
	documents = append(documents, qaSecretRefusalDocument{path: "CONTEXT.md", content: string(content)})
	sort.Slice(documents, func(i, j int) bool { return documents[i].path < documents[j].path })
	return documents
}

func checkQASecretRefusalDocumentation(documents []qaSecretRefusalDocument, documented string) error {
	const prefix = "QA evidence carries a secret"
	required := map[string]int{
		".agents/skills/roundfix/references/implement.md": 0,
		"skills/roundfix/references/implement.md":         0,
		"docs/user-guide/commands/implement.md":           0,
	}
	for _, document := range documents {
		normalized := strings.Join(strings.Fields(document.content), " ")
		for offset := 0; ; {
			index := strings.Index(normalized[offset:], prefix)
			if index < 0 {
				break
			}
			index += offset
			occurrence := normalized[index:]
			if !strings.HasPrefix(occurrence, documented) {
				return fmt.Errorf("%s: found QA secret refusal %q, want prefix %q", document.path, occurrence, documented)
			}
			if _, ok := required[document.path]; ok {
				required[document.path]++
			}
			offset = index + len(prefix)
		}
	}
	for path, count := range required {
		if count != 1 {
			return fmt.Errorf("%s: found %d QA secret refusal occurrences, want exactly one", path, count)
		}
	}
	return nil
}
