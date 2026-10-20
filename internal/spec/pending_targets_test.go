package spec

// Suite: pending creates targets and removed history inventory.
// Boundary IN: active and inactive Specs plus isolated local Git history.
// Boundary OUT: target exemptions and sorted, deduplicated removed paths.

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"roundfix/internal/gittest"
)

func TestRemovedHistoryPathsListsWhatHEADHeldAndTheTreeLacks(t *testing.T) {
	t.Parallel()
	if !citationSourcesTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	paths := []string{"docs/history/findings/worktree.md", "docs/history/backlog/deleted.md", "docs/history/adr/renamed.md", "docs/history/findings/present.md", "docs/history/reviews/outside.md"}
	for _, path := range paths {
		writeFile(t, filepath.Join(root, path), "# Same content\n")
	}
	legacyCommit(t, root, "docs: history")
	gittest.Run(t, root, "rm", paths[1], paths[4])
	gittest.Run(t, root, "mv", paths[2], "docs/history/adr/new.md")
	legacyCommit(t, root, "docs: remove and rename")
	if err := os.Remove(filepath.Join(root, paths[0])); err != nil {
		t.Fatal(err)
	}
	got, err := RemovedHistoryPaths(t.Context(), root)
	want := []string{paths[2], paths[1], paths[0]}
	if err != nil || !slices.Equal(got, want) {
		t.Fatalf("removed=%v err=%v want=%v", got, err, want)
	}
	outside := t.TempDir()
	if got, err := RemovedHistoryPaths(t.Context(), outside); err != nil || len(got) != 0 {
		t.Fatalf("outside Git=%v %v", got, err)
	}
	gittest.InitRepo(t, outside, "-b", "main")
	if got, err := RemovedHistoryPaths(t.Context(), outside); err != nil || len(got) != 0 {
		t.Fatalf("unborn HEAD=%v %v", got, err)
	}
}

func TestPendingTaskCreatesUsesOnlyActiveSpecsAndValidCreates(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "active/_prd.md"), "---\nstatus: active\n---\n")
	writeFile(t, filepath.Join(root, "inactive/_prd.md"), "---\nstatus: completed\n---\n")
	for _, name := range []string{"active", "inactive"} {
		writeFile(t, filepath.Join(root, name, "task_01.md"), "---\nstatus: in_progress\ntype: docs\n---\n## Context\n- creates: "+name+".md\n- interface: existing.md\n")
		writeFile(t, filepath.Join(root, name, "task_02.md"), "---\nstatus: invalid\ntype: docs\n---\n## Context\n- creates: invalid.md\n")
	}
	got, err := PendingTaskCreates(root)
	if err != nil || len(got) != 1 || !got["active.md"] {
		t.Fatalf("targets=%v err=%v", got, err)
	}
}
