package skills_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTheGuidesStateTheUnpinnedOwnedLockEntry(t *testing.T) {
	t.Parallel()

	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}

	cases := []struct {
		name  string
		paths []string
		terms []string
	}{
		{
			name:  "baseline",
			paths: []string{".agents/skills/roundfix/references/baseline.md", "skills/roundfix/references/baseline.md"},
			terms: []string{"- lock ref removed <name>", "skills.lockUnpinned"},
		},
		{
			name:  "setup",
			paths: []string{".agents/skills/roundfix/references/setup.md", "skills/roundfix/references/setup.md"},
			terms: []string{"Removed skills-lock.json ref: <name>", "DR-SKILL-LOCK-STALE", "roundfix skills install --target project"},
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			for _, path := range testCase.paths {
				content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
				if err != nil {
					t.Fatalf("read %s: %v", path, err)
				}
				normalized := strings.Join(strings.Fields(string(content)), " ")
				for _, term := range testCase.terms {
					if !strings.Contains(normalized, strings.Join(strings.Fields(term), " ")) {
						t.Errorf("%s is missing %q", path, term)
					}
				}
			}
		})
	}
}
