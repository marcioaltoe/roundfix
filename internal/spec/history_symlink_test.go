package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestHistoryCitationsSkipATrackedSymbolicLink(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	removed := "docs/history/findings/retired.md"
	historyEntryWrite(t, root, removed, "# Retired\n")
	historyEntryWrite(t, root, "AGENTS.md", "[retired](docs/history/findings/retired.md)\n")
	if err := os.Symlink("AGENTS.md", filepath.Join(root, "CLAUDE.md")); err != nil {
		t.Fatal(err)
	}
	gittest.InitRepo(t, root, "-b", "main")
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "-m", "history fixture")
	gittest.Run(t, root, "tag", "-a", "history-full", "-m", "full history")

	citations, err := HistoryCitations(t.Context(), root, []string{removed})
	if err != nil {
		t.Fatal(err)
	}
	if len(citations) != 1 || citations[0].Path != "AGENTS.md" {
		t.Fatalf("citations = %+v, want one citation in AGENTS.md", citations)
	}
	plan, err := PlanHistoryRewrite(t.Context(), root, ".", "docs/history/specs", []string{removed}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Files) != 1 || plan.Files["AGENTS.md"] == nil {
		t.Fatalf("rewrite files = %v, want AGENTS.md only", plan.Files)
	}
	if _, _, err := plan.Apply(t.Context(), root); err != nil {
		t.Fatal(err)
	}
	info, err := os.Lstat(filepath.Join(root, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatal("CLAUDE.md was materialized instead of remaining a symbolic link")
	}
	target, err := os.Readlink(filepath.Join(root, "CLAUDE.md"))
	if err != nil || target != "AGENTS.md" {
		t.Fatalf("CLAUDE.md target = %q, err=%v; want AGENTS.md", target, err)
	}
	content, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "docs/history/findings/retired.md@") {
		t.Fatalf("AGENTS.md does not carry the Git Locator: %q", content)
	}
}

func TestHistoryRewriteCoversARegularClaudeCarrier(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	removed := "docs/history/findings/retired.md"
	historyEntryWrite(t, root, removed, "# Retired\n")
	historyEntryWrite(t, root, "AGENTS.md", "[retired](docs/history/findings/retired.md)\n")
	historyEntryWrite(t, root, "CLAUDE.md", "See [retired](docs/history/findings/retired.md).\n")
	gittest.InitRepo(t, root, "-b", "main")
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "-m", "history fixture")
	gittest.Run(t, root, "tag", "-a", "history-full", "-m", "full history")

	plan, err := PlanHistoryRewrite(t.Context(), root, "docs/specs", "docs/history/specs", []string{removed}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Files) != 2 || plan.Files["AGENTS.md"] == nil || plan.Files["CLAUDE.md"] == nil {
		t.Fatalf("rewrite files = %v, want AGENTS.md and CLAUDE.md", plan.Files)
	}
	if _, _, err := plan.Apply(t.Context(), root); err != nil {
		t.Fatal(err)
	}
	for _, carrier := range []string{"AGENTS.md", "CLAUDE.md"} {
		content, err := os.ReadFile(filepath.Join(root, carrier))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(content), "docs/history/findings/retired.md@") || strings.Contains(string(content), "](docs/history/findings/retired.md)") {
			t.Fatalf("%s does not carry the Git Locator in place of the link: %q", carrier, content)
		}
	}
}
