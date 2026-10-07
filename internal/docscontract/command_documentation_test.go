//go:build docscontract

//verify:always

package docscontract

import (
	"bytes"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/cli"
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
