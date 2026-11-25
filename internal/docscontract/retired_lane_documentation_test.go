//go:build docscontract

//verify:always

package docscontract

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestTheGuideAndTheSkillListTheRetiredLaneKeys(t *testing.T) {
	root := baselineDocumentationRepoRoot()
	for _, tc := range []struct{ name, path string }{
		{"guide", "docs/user-guide/configuration.md"},
		{"canonical_reference", ".agents/skills/roundfix/references/setup.md"},
		{"mirror_reference", "skills/roundfix/references/setup.md"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			text := readContractDocument(t, filepath.Join(root, tc.path))
			for _, key := range []string{"delivery.linux_lane", "delivery.linux_lane_image"} {
				warning := "config: " + key + " is deprecated and ignored; the Linux Lane was retired, so remove it"
				if !strings.Contains(text, key) || !strings.Contains(text, warning) {
					t.Errorf("missing key or warning for %s", key)
				}
			}
		})
	}
}
