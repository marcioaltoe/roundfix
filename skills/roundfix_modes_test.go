package skills

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func roundfixSection(text, heading string) string {
	start := strings.Index(text, "## "+heading)
	if start < 0 {
		return ""
	}
	body := text[start+len("## "+heading):]
	if next := strings.Index(body, "\n## "); next >= 0 {
		body = body[:next]
	}
	return strings.TrimSpace(body)
}

func TestTheRoundfixSkillReadsItsModeFromTheRequest(t *testing.T) {
	t.Parallel()
	for name, source := range roundfixLayoutSources() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			data, err := fs.ReadFile(source, "roundfix/SKILL.md")
			if err != nil {
				t.Fatal(err)
			}
			var frontmatter struct {
				Metadata map[string]any `yaml:"metadata"`
			}
			parts := strings.SplitN(string(data), "---", 3)
			if len(parts) != 3 {
				t.Fatal("skill entry has no frontmatter")
			}
			if err := yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
				t.Fatal(err)
			}
			if got := frontmatter.Metadata["opencode/slash"]; got != "true" {
				t.Fatalf("metadata opencode/slash = %#v, want %q", got, "true")
			}
			text := string(data)
			modes := roundfixSection(text, "Modes")
			want := "A request names a Skill Mode with the first word after the skill's name: `follow-up`, `status` or `unattended`, with or without a leading `--`, in every harness: `/roundfix follow-up` or `/roundfix --follow-up` in Claude Code, `$roundfix status` or `$roundfix --status` in Codex, and `@roundfix unattended`, `@roundfix --unattended` or `/roundfix unattended` in OpenCode. A Skill Mode is not the CLI command `roundfix follow-up`, which the follow-up mode runs as one of its steps. Read the mode from the request text; a request whose first word after the skill's name is anything else, such as `/roundfix watch 123`, uses the skill as the rest of this file describes."
			if !strings.HasPrefix(modes, want) {
				t.Fatalf("Modes section does not start with required paragraph")
			}
			for _, mode := range []string{"--follow-up", "--status", "--unattended"} {
				if !strings.Contains(modes, mode) {
					t.Errorf("Modes section omits %s", mode)
				}
			}
			if !strings.Contains(modes, "references/modes.md") {
				t.Error("Modes section does not link its reference")
			}
			if regexp.MustCompile(`\$(?:[0-9]+|ARGUMENTS)`).MatchString(parts[2]) {
				t.Error("skill body contains an argument placeholder")
			}
		})
	}
}
