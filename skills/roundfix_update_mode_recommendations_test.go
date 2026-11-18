package skills

import (
	"io/fs"
	"strings"
	"testing"
)

func TestTheUpdateModeOffersTheProfileCommands(t *testing.T) {
	t.Parallel()
	for name, source := range roundfixLayoutSources() {
		name, source := name, source
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			modesData, err := fs.ReadFile(source, "roundfix/references/modes.md")
			if err != nil {
				t.Fatal(err)
			}
			update := roundfixWhitespace(roundfixSection(string(modesData), "--update"))
			for _, phrase := range []string{
				"recommendations",
				"roundfix profiles check --apply --scope",
				"--yes",
				"deviation.from",
			} {
				if !strings.Contains(update, roundfixWhitespace(phrase)) {
					t.Errorf("--update section omits %q", phrase)
				}
			}

			setupData, err := fs.ReadFile(source, "roundfix/references/setup.md")
			if err != nil {
				t.Fatal(err)
			}
			setup := roundfixWhitespace(string(setupData))
			for _, phrase := range []string{"Profiles:", "recommendations"} {
				if !strings.Contains(setup, roundfixWhitespace(phrase)) {
					t.Errorf("setup reference omits %q", phrase)
				}
			}
		})
	}
}
