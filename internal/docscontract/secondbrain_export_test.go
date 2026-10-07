//go:build docscontract

package docscontract

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const secondbrainHistoryExclusion = "!docs/history/specs/*/qa/"

func TestSecondbrainExportFollowsTheHistoryForm(t *testing.T) {
	root := baselineDocumentationRepoRoot()
	if _, err := os.Stat(filepath.Join(root, ".secondbrain-export")); errors.Is(err, os.ErrNotExist) {
		t.Skip("repository has no .secondbrain-export")
	} else if err != nil {
		t.Fatal(err)
	}
	if err := checkSecondbrainExport(root); err != nil {
		t.Fatal(err)
	}
}

func TestSecondbrainExportRuleRefusesBothMismatches(t *testing.T) {
	for _, test := range []struct {
		name       string
		legacy     bool
		exclusions bool
	}{
		{name: "exclusions without legacy folders", exclusions: true},
		{name: "legacy folder without exclusions", legacy: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			if test.legacy {
				if err := os.MkdirAll(filepath.Join(root, "docs", "history", "specs", "old-spec"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, "docs", "history", "specs", "old-spec", "_prd.md"), []byte("---\n---\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			content := "docs/\n"
			if test.exclusions {
				content += secondbrainHistoryExclusion + "\n"
			}
			if err := os.WriteFile(filepath.Join(root, ".secondbrain-export"), []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := checkSecondbrainExport(root); err == nil {
				t.Fatal("checkSecondbrainExport succeeded for a mismatched fixture")
			}
		})
	}
}

func checkSecondbrainExport(root string) error {
	export, err := os.ReadFile(filepath.Join(root, ".secondbrain-export"))
	if err != nil {
		return fmt.Errorf("read .secondbrain-export: %w", err)
	}
	entries, err := os.ReadDir(filepath.Join(root, "docs", "history", "specs"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read history specs: %w", err)
	}
	legacy := false
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(root, "docs", "history", "specs", entry.Name(), "_prd.md")); err == nil {
			legacy = true
			break
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("inspect legacy folder %q: %w", entry.Name(), err)
		}
	}
	hasExclusion := false
	for _, line := range strings.Split(string(export), "\n") {
		if strings.TrimSpace(line) == secondbrainHistoryExclusion {
			hasExclusion = true
		}
	}
	if legacy && !hasExclusion {
		return fmt.Errorf("legacy archive folders require %q", secondbrainHistoryExclusion)
	}
	if !legacy {
		for _, line := range strings.Split(string(export), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "!docs/history") {
				return fmt.Errorf("history exclusion %q remains without legacy archive folders", strings.TrimSpace(line))
			}
		}
	}
	return nil
}
