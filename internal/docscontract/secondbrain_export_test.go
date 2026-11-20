//go:build docscontract

//verify:always

package docscontract

import (
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/baseline"
)

func exportVariationAt(root string) ([]baseline.DocumentationDeparture, error) {
	content, err := os.ReadFile(filepath.Join(root, ".secondbrain-export"))
	if err != nil {
		return nil, err
	}
	facts := baseline.ExportListFacts{ProcessPaths: []string{"AGENTS.md", "CLAUDE.md", "docs/agents/", "docs/backlog/", "docs/findings/", "docs/handoffs/", "docs/specs/", "docs/history/"}}
	return baseline.ExportVariationDepartures(content, facts), nil
}

func TestThisRepositoryExportsNoProcessPath(t *testing.T) {
	departures, err := exportVariationAt(baselineDocumentationRepoRoot())
	if err != nil {
		t.Fatal(err)
	}
	if len(departures) != 0 {
		t.Fatalf("export variation departures: %+v", departures)
	}
}

func TestTheExportRuleRefusesAProcessPath(t *testing.T) {
	for _, tc := range []struct{ name, path string }{{"docs_whole", "docs/"}, {"agents_md", "AGENTS.md"}} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, ".secondbrain-export"), []byte("# reason\n"+tc.path+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			departures, err := exportVariationAt(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(departures) != 1 || departures[0].Code != "DR-EXPORT-PROCESS-PATH" {
				t.Fatalf("departures: %+v", departures)
			}
		})
	}
}
