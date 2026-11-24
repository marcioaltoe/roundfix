//go:build docscontract

//verify:always

package docscontract

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestTheDeliverGuideAndReferenceDescribeTheRecoveryLines(t *testing.T) {
	root := baselineDocumentationRepoRoot()
	for _, tc := range []struct{ name, path string }{
		{"guide", "docs/user-guide/commands/deliver.md"},
		{"canonical_reference", ".agents/skills/roundfix/references/deliver.md"},
		{"mirror_reference", "skills/roundfix/references/deliver.md"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			text := strings.Join(strings.Fields(readContractDocument(t, filepath.Join(root, tc.path))), " ")
			for _, want := range []string{
				"Owner: Delivery Queue has no recorded owner and <n> item(s) can still advance; run roundfix deliver resume",
				"Run: <slug> no Active Run",
				"Run: <slug> <run-id>; Tasks in progress: <task>, <task>",
				"Run: <slug> <run-id>; no Task in progress",
				"Retry failed", "QA-Only Failed Run",
			} {
				if !strings.Contains(text, want) {
					t.Errorf("missing %q", want)
				}
			}
			forbidden := []string{"A proven live owner, no recorded owner, or a queue"}
			if tc.name == "guide" {
				forbidden = append(forbidden, "With no recorded owner, it starts one directly.")
			}
			for _, phrase := range forbidden {
				if strings.Contains(text, phrase) {
					t.Errorf("unexpected %q", phrase)
				}
			}
		})
	}
}
