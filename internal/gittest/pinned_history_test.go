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

func TestPinnedHistoryExtractsEachPathSetOnceAndReadOnly(t *testing.T) {
	t.Parallel()
	repo := filepath.Join("..", "..")
	const source = "docs/history/specs/0205-an-advisory-judge-for-spec-authoring/_techspec.md"
	const other = "docs/history/specs/0205-an-advisory-judge-for-spec-authoring/_prd.md"
	root := PinnedHistory(t, repo, source)
	if again := PinnedHistory(t, repo, source); root != again {
		t.Fatalf("same paths extracted twice: %q, %q", root, again)
	}
	if different := PinnedHistory(t, repo, other); root == different {
		t.Fatal("different paths share an extraction")
	}
	if err := os.WriteFile(filepath.Join(root, "unexpected"), []byte("mutation"), 0o644); err == nil {
		t.Fatal("shared directory accepts writes")
	}
	if err := os.WriteFile(filepath.Join(root, source), []byte("mutation"), 0o644); err == nil {
		t.Fatal("shared file accepts writes")
	}
	var firstPair string
	t.Run("first reader", func(t *testing.T) {
		firstPair = PinnedHistory(t, repo, source, other)
	})
	t.Run("later reader", func(t *testing.T) {
		pair := PinnedHistory(t, repo, source, other)
		if pair != firstPair {
			t.Fatal("extraction disappeared after its first reader")
		}
		if err := os.WriteFile(filepath.Join(pair, "unexpected"), nil, 0o644); err == nil {
			t.Fatal("later reader received a writable tree")
		}
	})
	pair := PinnedHistory(t, repo, source, other)
	if again := PinnedHistory(t, repo, other, source, source); pair != again {
		t.Fatal("path set depends on order or duplicates")
	}
}
