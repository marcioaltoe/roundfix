package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHistorySanitizeAppliesBesideATrackedSymbolicLink(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, false)
	removed := filepath.Join(h.repo, "docs/history/findings/retired.md")
	historyWrite(t, removed, "# Retired\n")
	historyWrite(t, filepath.Join(h.repo, "AGENTS.md"), "[retired](docs/history/findings/retired.md)\n")
	if err := os.Symlink("AGENTS.md", filepath.Join(h.repo, "CLAUDE.md")); err != nil {
		t.Fatal(err)
	}
	historyCommit(t, h.repo, "Add symlink citation fixture")
	historyTag(t, h)

	code, _, stderr := historyRun(t, h, "sanitize", "--apply", "--batch", "1")
	if code != 0 || stderr != "" {
		t.Fatalf("history sanitize code=%d stderr=%q", code, stderr)
	}
	info, lstatErr := os.Lstat(filepath.Join(h.repo, "CLAUDE.md"))
	if lstatErr != nil {
		t.Fatal(lstatErr)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatal("CLAUDE.md was materialized instead of remaining a symbolic link")
	}
	target, readlinkErr := os.Readlink(filepath.Join(h.repo, "CLAUDE.md"))
	if readlinkErr != nil || target != "AGENTS.md" {
		t.Fatalf("CLAUDE.md target = %q, err=%v; want AGENTS.md", target, readlinkErr)
	}
	content, readErr := os.ReadFile(filepath.Join(h.repo, "AGENTS.md"))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(content), "docs/history/findings/retired.md@") {
		t.Fatalf("AGENTS.md does not carry the Git Locator: %q", content)
	}
}
