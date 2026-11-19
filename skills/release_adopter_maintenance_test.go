package skills

import (
	"io/fs"
	"os"
	"strings"
	"testing"
)

func TestTheReleaseDocumentsNameTheAdopterMaintenancePaths(t *testing.T) {
	t.Parallel()
	documents := make(map[string]string)
	for name, source := range roundfixLayoutSources() {
		data, err := fs.ReadFile(source, "roundfix/references/release.md")
		if err != nil {
			t.Fatal(err)
		}
		documents[name] = string(data)
	}
	for name, path := range map[string]string{"runbook": "../docs/user-guide/release-runbook.md", "usage": "../docs/user-guide/usage.md"} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		documents[name] = string(data)
	}
	for document, text := range documents {
		normalized := strings.Join(strings.Fields(text), " ")
		for check, phrase := range map[string]string{
			"adopter_test":                      "a repository whose `skills-lock.json` at the target records the `roundfix` skill",
			"setup_manifest":                    "the Setup Manifest `docs/agents/setup-context.json` and every path it records",
			"skills_lock":                       "`skills-lock.json` and every skill directory it records under `.agents/skills/` and `.claude/skills/`",
			"claude_md_and_specific_repository": "`CLAUDE.md`, `docs/agents/specific-repository.md`",
			"generated_subagents":               "a generated `.claude/agents/roundfix-*.md`",
			"project_config_and_context":        "`.roundfixrc.yml` and every `CONTEXT.md`",
			"base_and_target":                   "the two records read at the base and at the target",
			"roundfix_itself":                   "Roundfix's own repository authors the `roundfix` skill and records none in its lock",
		} {
			t.Run(document+"/"+check, func(t *testing.T) {
				if !strings.Contains(normalized, phrase) {
					t.Errorf("%s lacks %q", document, phrase)
				}
			})
		}
	}
}
