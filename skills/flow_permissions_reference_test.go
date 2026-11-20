package skills

import (
	"bytes"
	"io/fs"
	"strings"
	"testing"
)

func TestTheRoundfixSkillDescribesFlowPermissions(t *testing.T) {
	t.Parallel()
	tests := []struct{ name, sentence string }{
		{"see_the_list", "Run `roundfix permissions` to see the Flow Command List for each harness this repository holds; it writes nothing."},
		{"one_question_then_apply", "Before applying, show the preview and put one question to the user; apply only on that answer, with `--apply <digest>`, in every permission mode."},
		{"never_on_its_own_or_from_a_run", "Never apply or remove the list on your own, and never from a Run."},
		{"silence_is_never_that_yes", "A recommendation, silence or a missing answer is never that yes."},
		{"trust_rule", "When a listed command is still asked about, tell the user that a harness applies a project's allow rules only after the folder is trusted there."},
		{"not_a_security_boundary", "The list is not a security boundary."},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for name, source := range roundfixLayoutSources() {
				data, err := fs.ReadFile(source, "roundfix/references/setup.md")
				if err != nil {
					t.Fatal(err)
				}
				section := roundfixSection(string(data), "Flow permissions")
				if !strings.Contains(section, tc.sentence) {
					t.Errorf("%s: missing %q", name, tc.sentence)
				}
				if strings.Index(string(data), "## Flow permissions") >= strings.Index(string(data), "## Update check") {
					t.Errorf("%s: Flow permissions must precede Update check", name)
				}
			}
		})
	}
	t.Run("index_row_names_the_command", func(t *testing.T) {
		t.Parallel()
		const row = "| [setup](references/setup.md) | `init`, `setup`, `migrate`, `doctor`, `upgrade`, `update`, `skills`, `permissions` | Initializing config, checking readiness, upgrading, installing skills, or showing the Flow Command List with `roundfix permissions`. |"
		for name, source := range roundfixLayoutSources() {
			data, err := fs.ReadFile(source, "roundfix/SKILL.md")
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), row) {
				t.Errorf("%s: setup row omits permissions", name)
			}
		}
	})
	t.Run("copies_identical", func(t *testing.T) {
		t.Parallel()
		sources := roundfixLayoutSources()
		for _, path := range []string{"roundfix/SKILL.md", "roundfix/references/setup.md"} {
			canonical, err := fs.ReadFile(sources["canonical"], path)
			if err != nil {
				t.Fatal(err)
			}
			embedded, err := fs.ReadFile(sources["embedded"], path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(canonical, embedded) {
				t.Errorf("%s: copies differ", path)
			}
		}
	})
}
