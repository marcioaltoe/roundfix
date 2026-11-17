package skills

import (
	"io/fs"
	"strings"
	"testing"
)

func TestTheReleaseReferenceKeepsTheRepositoryPublicationTrigger(t *testing.T) {
	t.Parallel()
	for name, source := range roundfixLayoutSources() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			data, err := fs.ReadFile(source, "roundfix/references/release.md")
			if err != nil {
				t.Fatal(err)
			}
			section := roundfixWhitespace(roundfixSection(string(data), "Release planning"))
			for _, phrase := range []string{
				"docs/references/",
				"Release-Preparation Commit",
				"releasePreparation: true",
				"merge of the version change",
				"version tag is pushed",
			} {
				if !strings.Contains(section, roundfixWhitespace(phrase)) {
					t.Errorf("Release planning section omits %q", phrase)
				}
			}
			if strings.Contains(section, roundfixWhitespace("Preserve the existing tag-triggered workflow")) {
				t.Error("Release planning section keeps the removed tag-triggered workflow sentence")
			}
		})
	}
}
