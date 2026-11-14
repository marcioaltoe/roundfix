package skills

import (
	"io/fs"
	"strings"
	"testing"
)

func TestTheRoundfixSkillOffersTheRemovalOfAUserLevelCopy(t *testing.T) {
	t.Parallel()
	phrases := []struct{ name, text string }{
		{"reads_five_directories", "they read the five user-level directories in the `skills check` order, follow symbolic links and name a folder once under the first directory that reaches it"},
		{"names_hidden_headers", "a copy of a declared phase's skill in Claude Code's user-level directory, because it hides the Declared Skill Header"},
		{"keeps_exit_codes", "`skills check` warns on stderr and keeps its exit code, Doctor adds `user-skills: warn` after `skills:`, and the Update Check adds `User skills:` and `userSkills`"},
		{"names_removal_first", "each names the removal command, `roundfix skills remove --target <claude|codex|opencode> --apply`, with `--dir <path>` for the other two directories, before the refresh command"},
		{"changes_nothing_unasked", "None of them changes a copy; run a removal or a refresh only when the user asks"},
		{"writes_outside_the_repository", "because both write outside the repository, a refresh through the link the report names"},
		{"removal_lists_first", "Without `--apply`, `roundfix skills remove` only lists what it would remove, and it never touches the project's copies"},
	}
	for name, source := range roundfixLayoutSources() {
		t.Run(name, func(t *testing.T) {
			data, err := fs.ReadFile(source, "roundfix/references/setup.md")
			if err != nil {
				t.Fatal(err)
			}
			text := strings.Join(strings.Fields(string(data)), " ")
			for _, phrase := range phrases {
				t.Run(phrase.name, func(t *testing.T) {
					if count := strings.Count(text, strings.Join(strings.Fields(phrase.text), " ")); count != 1 {
						t.Errorf("phrase count=%d, want exactly one: %s", count, phrase.text)
					}
				})
			}
		})
	}
}
