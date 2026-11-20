package skills

import (
	"bytes"
	"io/fs"
	"strings"
	"testing"
)

func TestTheRoundfixSkillNamesTheFlowPermissionsReports(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, path, section, phrase string }{
		{"update_mode_names_next", "modes.md", "--update", "When `flowPermissions` names a harness that is `behind` or `unreadable`, name its `next` command"},
		{"update_mode_never_applies_unasked", "modes.md", "--update", "command, and never apply it unasked."},
		{"setup_names_both_reports", "setup.md", "Flow permissions", "`roundfix doctor` reports where each harness stands on its `flow-permissions:` line, and `roundfix update check` names each harness whose applied list is behind or not readable."},
		{"reports_never_apply", "setup.md", "Flow permissions", "Neither of them applies the list."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for name, source := range roundfixLayoutSources() {
				data, err := fs.ReadFile(source, "roundfix/references/"+tc.path)
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(roundfixSection(string(data), tc.section), tc.phrase) {
					t.Errorf("%s %s: missing %q", name, tc.path, tc.phrase)
				}
			}
		})
	}
	t.Run("copies_identical", func(t *testing.T) {
		t.Parallel()
		sources := roundfixLayoutSources()
		for _, path := range []string{"roundfix/SKILL.md", "roundfix/references/modes.md", "roundfix/references/setup.md"} {
			canonical, err := fs.ReadFile(sources["canonical"], path)
			if err != nil {
				t.Fatal(err)
			}
			mirrored, err := fs.ReadFile(sources["embedded"], path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(canonical, mirrored) {
				t.Errorf("%s: copies differ", path)
			}
		}
	})
}
