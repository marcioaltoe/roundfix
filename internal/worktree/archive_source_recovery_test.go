package worktree

// Suite: leftovers and archived QA read through a squash-dropped revision.
// Boundary IN: disposable Run worktrees and a main-only clone with a bare forge.
// Boundary OUT: network and daemon settlement.

import (
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func TestLeftoversRecoverASquashDroppedRevision(t *testing.T) {
	t.Parallel()
	for _, number := range []string{"7", "", "unreachable"} {
		t.Run(map[string]string{"7": "record", "": "adding subject", "unreachable": "unreachable forge"}[number], func(t *testing.T) {
			f, r, _ := newArchiveRecordRun(t, false, "none")
			r.PullRequest = number
			if number == "unreachable" {
				r.PullRequest = "7"
			}
			content, err := spec.RenderArchiveRecord(r)
			if err != nil {
				t.Fatal(err)
			}
			recordPath := spec.ArchiveRecordPath(spec.ArchiveDir(spec.ArchiveKindSpec), r.Spec)
			// Rewrite only this fixture's squash commit so the subject fallback names 7.
			gittest.Run(t, f.repoDir, "reset", "--hard", "main~1")
			mustWriteWorktreeTest(t, filepath.Join(f.repoDir, recordPath), string(content))
			gittest.Run(t, f.repoDir, "add", recordPath)
			gittest.Run(t, f.repoDir, "commit", "--amend", "-m", "feat: delivery (#7)")
			clone, _ := gittest.PullRequestClone(t, f.repoDir, f.ref.Branch, r.SourceRevision)
			if number == "unreachable" {
				gittest.Run(t, clone, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "absent.git"))
			}
			// A Run based before delivery has no copy of the source objects either.
			ref, err := Create(t.Context(), CreateOptions{UserRoot: clone, Location: t.TempDir(), RunID: "source-recovery", HeadSHA: strings.TrimSpace(gittest.Run(t, clone, "rev-parse", "main^"))})
			if err != nil {
				t.Fatal(err)
			}
			f.repoDir, f.ref = clone, ref
			f.run.GitRoot, f.run.WorkDir = canonicalPath(clone), canonicalPath(ref.Path)
			f.run.LocalBranch = "main"
			f.run.ID = ref.RunID
			mustWriteWorktreeTest(t, filepath.Join(ref.Path, "declared.txt"), "scoped leftover\n")
			result := inspectMergedHeadTestRun(t, f, nil)
			if number == "unreachable" {
				assertTerminalRunReconciliation(t, result, f.run, ReconciliationDirty)
				if !strings.Contains(result.Reason, "source_revision") {
					t.Fatalf("unavailable reason=%q", result.Reason)
				}
				assertMergeEvidencePreserved(t, f, result)
				return
			}
			assertTerminalRunReconciliation(t, result, f.run, ReconciliationSuperseded)
			gittest.Run(t, clone, "cat-file", "-e", r.SourceRevision+"^{commit}")
			report, err := newestQAReportAtHead(t.Context(), execGitRunner{}, clone, "main", r.Spec)
			if err != nil || report != r.Source+"/qa/"+r.QAReport {
				t.Fatalf("report=%q err=%v", report, err)
			}
		})
	}
}
