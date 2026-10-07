package gittest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPinnedHistoryMaterializesThePinnedPaths(t *testing.T) {
	t.Parallel()
	repo := filepath.Join("..", "..")
	const source = "docs/history/specs/0205-an-advisory-judge-for-spec-authoring/_techspec.md"
	before := Run(t, repo, "status", "--porcelain", "--untracked-files=all")
	root := PinnedHistory(t, repo, source)
	got, err := os.ReadFile(filepath.Join(root, source))
	if err != nil {
		t.Fatal(err)
	}
	want := Run(t, repo, "show", PinnedHistoryRevision+":"+source)
	if string(got) != want {
		t.Fatal("materialized bytes differ from pinned blob")
	}
	if _, err := os.Stat(filepath.Join(root, "docs/history/specs/0205-an-advisory-judge-for-spec-authoring/_prd.md")); !os.IsNotExist(err) {
		t.Fatalf("unrequested file materialized: %v", err)
	}
	if after := Run(t, repo, "status", "--porcelain", "--untracked-files=all"); after != before {
		t.Fatal("source repository changed")
	}
	if !strings.Contains(string(got), "Questions and thresholds") {
		t.Fatal("wrong corpus path")
	}
}
