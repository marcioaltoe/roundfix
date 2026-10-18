package worktree

// Suite: delivered Spec directory scope and retained dirty Run diagnostics.
// Boundary IN: disposable Git repositories, archived Tasks and Run worktrees.
// Boundary OUT: network, the real Run Database and Daemon settlement.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"roundfix/internal/spec"
)

func newDirectoryLeftoverFixture(t *testing.T, record bool) terminalRunFixture {
	t.Helper()
	f := newArchivedLeftoverFixture(t, "archived-directory-leftovers")
	root := archivedLeftoverRoot(f)
	task, err := gitBlobAtHead(t.Context(), execGitRunner{}, f.repoDir, "main", root+"/task_01.md")
	if err != nil {
		t.Fatal(err)
	}
	scope := strings.NewReplacer(
		"`declared.txt`", "`declared`",
		"`created.txt`", "`created`",
		"`deleted.txt`", "`docs/history/findings`",
		"`recorded.txt`", "`recorded`",
	).Replace(string(task))
	deleted := "docs/history/findings/old.md"
	commitMergedHeadFiles(t, f.ref.Path, map[string]string{deleted: "delivered\n"}, "docs: retain delivered finding")
	sourceHead := commitMergedHeadFiles(t, f.repoDir, map[string]string{
		root + "/task_01.md": scope,
	}, "docs: declare directory scope")
	if record {
		content, err := spec.RenderArchiveRecord(spec.ArchiveRecord{
			Spec: f.run.SpecSlug, Title: "Directory scope", Source: root,
			SourceRevision: sourceHead, Disposition: spec.ArchivePass,
			QATask: "task_qa", QAReport: "qa-report-2026-10-01.md", QAVerdict: "pass",
		})
		if err != nil {
			t.Fatal(err)
		}
		gitWorktreeTest(t, f.repoDir, "rm", "-r", "--", root)
		commitMergedHeadFiles(t, f.repoDir, map[string]string{
			spec.ArchiveRecordPath(spec.ArchiveDir(spec.ArchiveKindSpec), f.run.SpecSlug): string(content),
		}, "docs: retain archive source revision")
	}
	if err := os.Remove(filepath.Join(f.ref.Path, filepath.FromSlash(deleted))); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{"recorded/nested/leftover.txt", "declared/leftover.txt", "created/leftover.txt"} {
		fullPath := filepath.Join(f.ref.Path, filepath.FromSlash(file))
		mustMkdirWorktreeTest(t, filepath.Dir(fullPath))
		mustWriteWorktreeTest(t, fullPath, "scoped leftover\n")
	}
	return f
}

func TestMergedSpecRunWithLeftoversBeneathADeclaredDirectoryIsSuperseded(t *testing.T) {
	t.Parallel()
	for _, archive := range []struct {
		name   string
		record bool
	}{{"legacy folder", false}, {"archive record", true}} {
		t.Run(archive.name, func(t *testing.T) {
			for _, sibling := range []string{"", "docs/history/findings-other/keep.md", "recorded-other.txt"} {
				name := sibling
				if name == "" {
					name = "all paths scoped"
				}
				t.Run(name, func(t *testing.T) {
					f := newDirectoryLeftoverFixture(t, archive.record)
					if sibling != "" {
						fullPath := filepath.Join(f.ref.Path, filepath.FromSlash(sibling))
						mustMkdirWorktreeTest(t, filepath.Dir(fullPath))
						mustWriteWorktreeTest(t, fullPath, "preserve\n")
					}
					result := inspectMergedHeadTestRun(t, f, nil)
					if sibling != "" {
						assertTerminalRunReconciliation(t, result, f.run, ReconciliationDirty)
						want := reconciliationReasonDirty + "; " + sibling + " is outside the archived Spec's declared and recorded paths"
						if result.Reason != want {
							t.Fatalf("reason = %q, want %q", result.Reason, want)
						}
						assertMergeEvidencePreserved(t, f, result)
						return
					}
					assertTerminalRunReconciliation(t, result, f.run, ReconciliationSuperseded)
					if err := ApplyTerminalRun(t.Context(), &recordingTerminalRunStore{}, result); err != nil {
						t.Fatalf("apply: %v", err)
					}
					assertPathRemoved(t, f.ref.Path)
					assertBranchRemoved(t, f.repoDir, f.ref.Branch)
				})
			}
		})
	}
}

func TestMergedSpecRunNamesTheLeftoverOutsideItsScope(t *testing.T) {
	t.Parallel()
	for _, file := range []string{"unrelated.txt", strings.Repeat("é", 65) + ".txt"} {
		t.Run(file, func(t *testing.T) {
			f := newArchivedLeftoverFixture(t, "archived-outside-scope-reason")
			addDeclaredLeftovers(t, f)
			// Creating the later path first must not change byte-order selection.
			mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, file+".later"), "preserve\n")
			mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, file), "preserve\n")
			for _, records := range [][]MergedHead{nil, {mergedHeadTestRecord(f, mergedHeadTestHead(t, f.repoDir, "main"))}} {
				result := inspectMergedHeadTestRun(t, f, records)
				assertTerminalRunReconciliation(t, result, f.run, ReconciliationDirty)
				want := reconciliationReasonDirty + "; " + file + " is outside the archived Spec's declared and recorded paths"
				if file == "unrelated.txt" {
					if result.Reason != want {
						t.Fatalf("reason = %q, want %q", result.Reason, want)
					}
				} else if !strings.HasPrefix(result.Reason, reconciliationReasonDirty+"; ") ||
					!strings.HasSuffix(result.Reason, "declared and recorded paths") ||
					!strings.Contains(result.Reason, "...") ||
					len(result.Reason) > reconciliationReasonMaxBytes || !utf8.ValidString(result.Reason) {
					t.Fatalf("bounded reason = %q (%d bytes)", result.Reason, len(result.Reason))
				}
				assertMergeEvidencePreserved(t, f, result)
			}
		})
	}
}

func TestMergedSpecRunNamesTheUnrepresentedCommit(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "archived-unrepresented-commit-reason")
	commit := commitMergedHeadTask(t, f.ref.Path, f.run.SpecSlug, "task_02", "unique.txt", "preserve\n")
	addDeclaredLeftovers(t, f)
	for _, records := range [][]MergedHead{nil, {mergedHeadTestRecord(f, mergedHeadTestHead(t, f.repoDir, "main"))}} {
		result := inspectMergedHeadTestRun(t, f, records)
		assertTerminalRunReconciliation(t, result, f.run, ReconciliationDirty)
		if !strings.HasPrefix(result.Reason, reconciliationReasonDirty+"; Run commit "+commit[:12]) ||
			!strings.HasSuffix(result.Reason, "Task task_02 is not completed") || len(result.Reason) > reconciliationReasonMaxBytes {
			t.Fatalf("reason = %q", result.Reason)
		}
		assertMergeEvidencePreserved(t, f, result)
	}
}
