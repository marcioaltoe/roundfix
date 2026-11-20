package flowpermissions

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/suiteguard"
)

var recordCodexRules = flag.Bool("record-codex-rules", false, "record the shipped Codex rules golden")

const recordCodexRulesCommand = "go test ./internal/flowpermissions -run '^TestCodexRulesFile$' -record-codex-rules"

// This decoder is independent of ReadCodexRules: a forbidden rendering must
// name its offending command even when the production reader refuses it.
func codexRenderedCommands(data []byte) ([]string, error) {
	var commands []string
	for _, line := range strings.Split(string(data), "\n") {
		pattern, found := strings.CutPrefix(line, "    pattern = ")
		if !found {
			continue
		}
		var positions []json.RawMessage
		if err := json.Unmarshal([]byte(strings.TrimSuffix(pattern, ",")), &positions); err != nil {
			return nil, err
		}
		expanded := []string{""}
		for _, position := range positions {
			var word string
			alternatives := []string{}
			if json.Unmarshal(position, &word) == nil {
				alternatives = append(alternatives, word)
			} else if err := json.Unmarshal(position, &alternatives); err != nil {
				return nil, err
			}
			var next []string
			for _, prefix := range expanded {
				for _, alternative := range alternatives {
					next = append(next, strings.TrimSpace(prefix+" "+alternative))
				}
			}
			expanded = next
		}
		commands = append(commands, expanded...)
	}
	return commands, nil
}

func renderedCommands(l list, harness string) ([]string, error) {
	switch harness {
	case "claude":
		var commands []string
		for _, entry := range l.entries {
			rule := ClaudeRule(entry)
			command := strings.TrimSuffix(strings.TrimPrefix(rule, "Bash("), " *)")
			commands = append(commands, command)
		}
		return commands, nil
	case "codex":
		return codexRenderedCommands(l.renderCodexRules())
	default:
		return l.plainEntries(), nil
	}
}

func forbiddenRuleFindings(commands []string, exclusions []Exclusion) []string {
	forbiddenFirst := []string{"bash", "sh", "zsh", "python", "python3", "node", "go", "make", "npm", "bun", "npx"}
	gitReads := []string{"archive", "cat-file", "check-ignore", "diff", "diff-tree", "grep", "ls-files", "merge-tree", "rev-list", "rev-parse", "show", "status"}
	var findings []string
	for _, command := range commands {
		words := strings.Fields(command)
		bad := len(words) < 2
		if len(words) > 0 {
			bad = bad || slices.Contains(forbiddenFirst, words[0]) || !slices.Contains([]string{"roundfix", "git", "gh"}, words[0])
			for _, word := range words {
				if strings.ContainsAny(word, "*?[]") {
					bad = true
				}
			}
			if words[0] == "roundfix" && len(words) != 2 {
				bad = true
			}
			if words[0] == "git" && (len(words) != 2 || !slices.Contains(gitReads, words[1])) {
				bad = true
			}
			if words[0] == "gh" && !slices.Equal(words, []string{"gh", "api"}) && !slices.Equal(words, []string{"gh", "pr", "view"}) {
				bad = true
			}
		}
		for _, exclusion := range exclusions {
			if wordPrefix(words, exclusion.Words) || wordPrefix(exclusion.Words, words) {
				bad = true
			}
		}
		if bad {
			findings = append(findings, command)
		}
	}
	return findings
}

func TestRenderingsHoldNoForbiddenRule(t *testing.T) {
	t.Parallel()
	for _, harness := range []string{"claude", "codex", "opencode"} {
		t.Run(harness, func(t *testing.T) {
			commands, err := renderedCommands(shipped, harness)
			if err != nil {
				t.Fatal(err)
			}
			if findings := forbiddenRuleFindings(commands, Exclusions()); len(findings) != 0 {
				t.Fatalf("forbidden rules: %v", findings)
			}
			if !slices.Equal(commands, PlainEntries()) {
				t.Errorf("rendering changed the command list: %v", commands)
			}
		})
	}
	for _, planted := range []struct{ name, command string }{
		{"planted_shell_is_reported", "bash exec"},
		{"planted_interpreter_is_reported", "python3 script"},
		{"planted_wildcard_subcommand_is_reported", "roundfix *"},
		{"planted_exclusion_is_reported", "roundfix permissions"},
	} {
		t.Run(planted.name, func(t *testing.T) {
			l := list{entries: append(Entries(), Entry{Words: strings.Fields(planted.command)})}
			for _, harness := range []string{"claude", "codex", "opencode"} {
				commands, err := renderedCommands(l, harness)
				if err != nil {
					t.Fatal(err)
				}
				if findings := forbiddenRuleFindings(commands, Exclusions()); !slices.Contains(findings, planted.command) {
					t.Errorf("%s did not name %q: %v", harness, planted.command, findings)
				}
			}
		})
	}
	// Check every forbidden first word and unconditionally forbidden Git writes,
	// including writes no skill currently names.
	for _, command := range []string{"sh run", "zsh run", "python script", "node script", "go test", "make all", "npm install", "bun run", "npx tool", "* spec", "git commit", "git push", "git merge", "git clone", "git fetch", "gh pr create", "gh pr merge"} {
		if len(forbiddenRuleFindings([]string{command}, Exclusions())) == 0 {
			t.Errorf("forbidden checker missed %s", command)
		}
	}
}

func guardMatches(rule, command string) bool {
	pattern := strings.TrimSuffix(strings.TrimPrefix(rule, "Bash("), ")")
	pattern = strings.ReplaceAll(regexp.QuoteMeta(pattern), `\*`, `(?s:.*)`)
	return regexp.MustCompile("^" + pattern + "$").MatchString(command)
}

func TestGuardEntriesMatchTheirOption(t *testing.T) {
	t.Parallel()
	guards := Guards()
	if len(guards) != 19 {
		t.Fatalf("got %d guards, want 19", len(guards))
	}
	// The independent option matrix detects a missing, duplicate or replaced guard.
	expected := map[string][]string{
		"git archive":   {"--output", "-o", "--remote", "--exec"},
		"git diff":      {"--output", "--ext-diff"},
		"git diff-tree": {"--output", "--ext-diff"},
		"git grep":      {"--open-files-in-pager", "-O"},
		"git show":      {"--output", "--ext-diff"},
		"gh api":        {"--method", "-X", "--field", "-F", "--raw-field", "-f", "--input"},
	}
	seen := map[string]bool{}
	for _, guard := range guards {
		command := strings.Join(guard.Entry.Words, " ")
		key := command + " " + guard.Option
		if !slices.Contains(PlainEntries(), command) || !slices.Contains(expected[command], guard.Option) || seen[key] {
			t.Errorf("unexpected or duplicate Guard Entry: %+v", guard)
		}
		seen[key] = true
		want := "Bash(" + command + " *" + guard.Option + "*)"
		if !strings.HasPrefix(guard.Option, "--") {
			want = "Bash(" + command + "* " + guard.Option + "*)"
		}
		if guard.ClaudeRule != want {
			t.Errorf("guard rule = %q, want %q", guard.ClaudeRule, want)
		}
	}
	for command, options := range expected {
		for _, option := range options {
			if !seen[command+" "+option] {
				t.Errorf("missing guard %s %s", command, option)
			}
		}
	}
	for _, form := range []struct{ name, middle, suffix string }{
		{"right_after_the_command", " ", " value"},
		{"after_another_argument", " HEAD ", " value"},
		{"with_an_equals_sign", " ", "=value"},
		{"short_option_form", " ", "value"},
		{"long_option_form", " HEAD ", "=value"},
	} {
		t.Run(form.name, func(t *testing.T) {
			count := 0
			for _, guard := range guards {
				long := strings.HasPrefix(guard.Option, "--")
				if form.name == "short_option_form" && long || form.name == "long_option_form" && !long {
					continue
				}
				count++
				command := strings.Join(guard.Entry.Words, " ") + form.middle + guard.Option + form.suffix
				if !guardMatches(guard.ClaudeRule, command) {
					t.Errorf("%q did not match %q", guard.ClaudeRule, command)
				}
			}
			if count == 0 {
				t.Fatal("no guard exercised")
			}
		})
	}
	t.Run("read_without_the_option_is_not_guarded", func(t *testing.T) {
		for _, command := range []string{"git show HEAD:README.md", "git diff --name-only HEAD", "git grep -l hi HEAD", "git archive --format=tar HEAD README.md", "gh api repos/o/r/releases"} {
			for _, guard := range guards {
				if guardMatches(guard.ClaudeRule, command) {
					t.Errorf("%q guarded ordinary read %q", guard.ClaudeRule, command)
				}
			}
		}
	})
}

func TestCodexRulesFile(t *testing.T) {
	// Sequential: the record flag writes the shared Codex golden.
	if *recordCodexRules {
		suiteguard.DeclareSanctionedRegeneration(recordCodexRulesCommand)
		if err := os.WriteFile("testdata/roundfix-flow.rules.golden", RenderCodexRules(), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	data := RenderCodexRules()
	t.Run("golden", func(t *testing.T) {
		golden, err := os.ReadFile("testdata/roundfix-flow.rules.golden")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(data, golden) {
			t.Fatal("the Codex rules file changed; record it with: " + recordCodexRulesCommand)
		}
	})
	t.Run("every_entry_allowed", func(t *testing.T) {
		commands, err := codexRenderedCommands(data)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(commands, PlainEntries()) {
			t.Errorf("allowed commands = %v, want %v", commands, PlainEntries())
		}
		if got := strings.Count(string(data), "    decision = \"allow\","); got != 4 {
			t.Errorf("got %d allow blocks, want 4", got)
		}
	})
	t.Run("no_exclusion_rule", func(t *testing.T) {
		commands, err := codexRenderedCommands(data)
		if err != nil {
			t.Fatal(err)
		}
		if findings := forbiddenRuleFindings(commands, Exclusions()); len(findings) > 0 {
			t.Errorf("forbidden rules: %v", findings)
		}
	})
	t.Run("no_guard", func(t *testing.T) {
		if bytes.Contains(data, []byte("ask")) || bytes.Contains(data, []byte("Bash(")) {
			t.Error("Codex rendering contains Claude guard syntax")
		}
		for _, guard := range Guards() {
			if bytes.Contains(data, []byte(strconv.Quote(guard.Option))) {
				t.Errorf("Codex rendering contains option %s", guard.Option)
			}
		}
	})
	t.Run("reads_back_every_entry", func(t *testing.T) {
		entries, ok := ReadCodexRules(data)
		if !ok || !slices.Equal(entries, PlainEntries()) {
			t.Errorf("ReadCodexRules = %v, %v", entries, ok)
		}
		withEmptyLines := strings.ReplaceAll(string(data), "\n", "\n\n")
		entries, ok = ReadCodexRules([]byte(withEmptyLines))
		if !ok || !slices.Equal(entries, PlainEntries()) {
			t.Errorf("ReadCodexRules with empty lines = %v, %v", entries, ok)
		}
		older := list{entries: []Entry{{Words: []string{"roundfix", "spec"}}, {Words: []string{"git", "status"}}, {Words: []string{"gh", "pr", "view"}}}}
		entries, ok = ReadCodexRules(older.renderCodexRules())
		if !ok || !slices.Equal(entries, older.plainEntries()) {
			t.Errorf("older rendered list = %v, %v", entries, ok)
		}
	})
	t.Run("other_content_is_not_ok", func(t *testing.T) {
		valid := string(data)
		for i, invalid := range []string{
			"", strings.Replace(valid, "# Roundfix", "# Edited", 1),
			strings.Replace(valid, "# To remove it:", "# Remove:", 1),
			valid + "load(\"other.rules\")\n", valid + "# unrelated comment\n",
			strings.Replace(valid, `decision = "allow"`, `decision = "ask"`, 1),
			strings.Replace(valid, `decision = "allow"`, `decision = "deny"`, 1),
			strings.Replace(valid, ")\n", "", 1), valid + "prefix_rule(\n", valid + ")\n",
			strings.Replace(valid, `"archive", `, `"*", `, 1),
			strings.Replace(valid, `pattern = ["roundfix", [`, `pattern = ["bash", [`, 1),
			strings.Replace(valid, `pattern = ["gh", "api"]`, `pattern = ["gh", "api", "--method", "POST"]`, 1),
			strings.Replace(valid, `pattern = ["gh", "api"]`, `pattern = ["gh", null]`, 1),
			strings.Replace(valid, `pattern = ["gh", "api"]`, `pattern = ["gh", ["api"]]`, 1),
			strings.Replace(valid, `pattern = ["gh", "api"]`, `pattern = ["gh","api"]`, 1),
			strings.Replace(valid, `pattern = ["gh", "api"]`, `pattern = ["gh", "api"] # extra`, 1),
			strings.Replace(valid, "Roundfix Flow Command List: a Git read the skills name", "other justification", 1),
		} {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				if entries, ok := ReadCodexRules([]byte(invalid)); ok || entries != nil {
					t.Errorf("accepted invalid content: entries %v, ok %v", entries, ok)
				}
			})
		}
	})
}
