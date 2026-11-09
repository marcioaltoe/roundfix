package spec

// Suite: QA evidence belongs to the report's Run.
// Invariant: reused/imported bytes and foreign Run paths refuse; Carried Rows do not.
// Boundary IN: temporary Git repositories, report links and executable paths.
// Boundary OUT: CLI streams and Daemon repair turns.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func borrowedEvidenceRepo(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "roundfix/run-own")
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	return root, filepath.Join(root, "docs/specs/x/qa/qa-report-2026-11-09.md")
}

func borrowedEvidenceReport(body string) string {
	return "---\nverdict: pass\n---\n" + body
}

func assertBorrowedEvidence(t *testing.T, report string, options BorrowedEvidenceOptions, want []BorrowedEvidence) {
	t.Helper()
	got, err := FindBorrowedEvidence(t.Context(), report, options)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("borrowed evidence = %+v, %v; want %+v", got, err, want)
	}
	err = RequireOwnRunEvidence(t.Context(), report, options)
	var refusal BorrowedEvidenceError
	if len(want) == 0 {
		if err != nil {
			t.Fatalf("honest evidence refused: %v", err)
		}
	} else if !errors.As(err, &refusal) || !reflect.DeepEqual(refusal.Evidence, want) {
		t.Fatalf("refusal = %v; want %+v", err, want)
	}
}

func TestBorrowedEvidenceRefusesAnEarlierRunsEvidenceDirectory(t *testing.T) {
	t.Parallel()
	root, report := borrowedEvidenceRepo(t)
	qa := filepath.Dir(report)
	for _, name := range []string{"head-exclusive.execution.txt", "replay.py.txt"} {
		writeFile(t, filepath.Join(qa, "evidence/2026-11-06-gate", name), "earlier Run\n")
	}
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "-m", "earlier evidence")
	gittest.Run(t, root, "commit", "--allow-empty", "-m", "audited head")
	options := BorrowedEvidenceOptions{RunBranchPrefix: "roundfix/run-"}
	writeFile(t, report, borrowedEvidenceReport("[run](evidence/2026-11-06-gate/head-exclusive.execution.txt) [replay](evidence/2026-11-06-gate/replay.py.txt)\n"))
	assertBorrowedEvidence(t, report, options, []BorrowedEvidence{
		{4, "evidence/2026-11-06-gate/head-exclusive.execution.txt", BorrowedEvidenceReused},
		{4, "evidence/2026-11-06-gate/replay.py.txt", BorrowedEvidenceReused},
	})
	writeFile(t, report, borrowedEvidenceReport("[directory](evidence/2026-11-06-gate/)\n"))
	assertBorrowedEvidence(t, report, options, []BorrowedEvidence{
		{4, "evidence/2026-11-06-gate/head-exclusive.execution.txt", BorrowedEvidenceReused},
		{4, "evidence/2026-11-06-gate/replay.py.txt", BorrowedEvidenceReused},
	})
	writeFile(t, filepath.Join(qa, "evidence/fresh.txt"), "fresh\n")
	writeFile(t, filepath.Join(qa, "evidence/2026-11-06-gate/replay.py.txt"), "changed\n")
	writeFile(t, report, borrowedEvidenceReport("[new](evidence/fresh.txt) [changed](evidence/2026-11-06-gate/replay.py.txt) [missing](evidence/missing.txt) [elsewhere](../_prd.md)\n"))
	assertBorrowedEvidence(t, report, options, nil)
	// A corrupt index is an infrastructure error, never Borrowed Evidence.
	writeFile(t, filepath.Join(root, ".git/index"), "broken\n")
	if items, err := FindBorrowedEvidence(t.Context(), report, options); err == nil || items != nil {
		t.Fatalf("corrupt index = %+v, %v", items, err)
	}
}

func TestBorrowedEvidenceReadsACommittedReportAgainstItsOwnCommit(t *testing.T) {
	t.Parallel()
	for _, later := range []bool{false, true} {
		t.Run(map[bool]string{false: "same commit", true: "later report"}[later], func(t *testing.T) {
			root, report := borrowedEvidenceRepo(t)
			writeFile(t, filepath.Join(root, "initial.txt"), "initial\n")
			gittest.Run(t, root, "add", ".")
			gittest.Run(t, root, "commit", "-m", "initial")
			writeFile(t, filepath.Join(filepath.Dir(report), "evidence/run.txt"), "evidence\n")
			if later {
				gittest.Run(t, root, "add", ".")
				gittest.Run(t, root, "commit", "-m", "earlier evidence")
			}
			writeFile(t, report, borrowedEvidenceReport("[run](evidence/run.txt)\n"))
			gittest.Run(t, root, "add", ".")
			gittest.Run(t, root, "commit", "-m", "report")
			gittest.Run(t, root, "commit", "--allow-empty", "-m", "later history")
			var want []BorrowedEvidence
			if later {
				want = []BorrowedEvidence{{4, "evidence/run.txt", BorrowedEvidenceReused}}
			}
			assertBorrowedEvidence(t, report, BorrowedEvidenceOptions{RunBranchPrefix: "roundfix/run-"}, want)
			// A modified tracked report audits HEAD instead of its addition parent.
			writeFile(t, report, borrowedEvidenceReport("[run](evidence/run.txt)\nchanged\n"))
			assertBorrowedEvidence(t, report, BorrowedEvidenceOptions{RunBranchPrefix: "roundfix/run-"}, []BorrowedEvidence{{4, "evidence/run.txt", BorrowedEvidenceReused}})
		})
	}
	root, report := borrowedEvidenceRepo(t)
	writeFile(t, filepath.Join(filepath.Dir(report), "evidence/run.txt"), "root evidence\n")
	writeFile(t, report, borrowedEvidenceReport("[run](evidence/run.txt)\n"))
	options := BorrowedEvidenceOptions{RunBranchPrefix: "roundfix/run-"}
	assertBorrowedEvidence(t, report, options, nil) // unborn HEAD
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "-m", "root report")
	assertBorrowedEvidence(t, report, options, nil) // root addition has no parent
}

func TestBorrowedEvidenceExemptsACarriedRow(t *testing.T) {
	t.Parallel()
	root, report := borrowedEvidenceRepo(t)
	writeFile(t, filepath.Join(filepath.Dir(report), "evidence/run.txt"), "old\n")
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "-m", "old evidence")
	body := "\n## Results\n\n| # | Status | Evidence |\n| - | --- | --- |\n| Q1 | carried (snapshot unchanged) | [old](evidence/run.txt) |\n\n### Q1\n[detail](evidence/run.txt)\n#### More\n[detail](evidence/run.txt)\n### Q2\n[observed](evidence/run.txt)\n"
	writeFile(t, report, borrowedEvidenceReport(body))
	assertBorrowedEvidence(t, report, BorrowedEvidenceOptions{RunBranchPrefix: "roundfix/run-"}, []BorrowedEvidence{{16, "evidence/run.txt", BorrowedEvidenceReused}})
	writeFile(t, report, borrowedEvidenceReport(strings.ReplaceAll(body, "carried (snapshot unchanged)", "pass")))
	assertBorrowedEvidence(t, report, BorrowedEvidenceOptions{RunBranchPrefix: "roundfix/run-"}, []BorrowedEvidence{
		{9, "evidence/run.txt", BorrowedEvidenceReused}, {12, "evidence/run.txt", BorrowedEvidenceReused}, {14, "evidence/run.txt", BorrowedEvidenceReused}, {16, "evidence/run.txt", BorrowedEvidenceReused},
	})
}

func TestBorrowedEvidenceNamesAnExecutableOutsideTheWorktree(t *testing.T) {
	t.Parallel()
	root, report := borrowedEvidenceRepo(t)
	outside := filepath.Join(t.TempDir(), "python")
	writeFile(t, outside, "interpreter\n")
	python := filepath.Join(root, ".venv/bin/python")
	if err := os.MkdirAll(filepath.Dir(python), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, python); err != nil {
		t.Fatal(err)
	}
	options := BorrowedEvidenceOptions{RunBranchPrefix: "roundfix/run-"}
	for _, tc := range []struct {
		key  string
		want []BorrowedEvidence
	}{
		{"user_flow_executables:\n  - ./bin/tool\n  - .venv/bin/tool\n  - .venv/bin/python\n", nil},
		{"user_flow_executables:\n  - /elsewhere/bin/tool\n  - ../other/bin/tool\n", []BorrowedEvidence{{4, "/elsewhere/bin/tool", BorrowedEvidenceOutsideWorktree}, {5, "../other/bin/tool", BorrowedEvidenceOutsideWorktree}}},
		{"user_flow_executables: scalar\n", []BorrowedEvidence{{3, "user_flow_executables", BorrowedEvidenceMalformed}}},
		{"user_flow_executables: [42]\n", []BorrowedEvidence{{3, "user_flow_executables", BorrowedEvidenceMalformed}}},
		{"user_flow_executables: []\n", nil},
		{"", nil},
	} {
		writeFile(t, report, "---\nverdict: pass\n"+tc.key+"---\n")
		assertBorrowedEvidence(t, report, options, tc.want)
	}
}

func TestBorrowedEvidenceNamesAnotherRunWorktree(t *testing.T) {
	t.Parallel()
	root, report := borrowedEvidenceRepo(t)
	writeFile(t, filepath.Join(root, "initial.txt"), "initial\n")
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "-m", "initial")
	var paths []string
	for _, branch := range []string{"roundfix/run-old", "roundfix/run-own-task_01", "feature"} {
		path := filepath.Join(t.TempDir(), "worktree")
		gittest.Run(t, root, "worktree", "add", "-b", branch, path)
		path, err := filepath.EvalSymlinks(path)
		if err != nil {
			t.Fatal(err)
		}
		paths = append(paths, path)
	}
	options := BorrowedEvidenceOptions{RunBranchPrefix: "roundfix/run-"}
	evidence := filepath.Join(filepath.Dir(report), "evidence/run.txt")
	writeFile(t, report, borrowedEvidenceReport("[run](evidence/run.txt)\n"))
	writeFile(t, evidence, paths[0]+"/.venv/bin/tool\n")
	assertBorrowedEvidence(t, report, options, []BorrowedEvidence{{4, "evidence/run.txt", BorrowedEvidenceForeignWorktree + " " + paths[0]}})
	writeFile(t, evidence, root+"/bin/tool\n"+paths[1]+"/bin/tool\n"+paths[2]+"/bin/tool\n")
	assertBorrowedEvidence(t, report, options, nil)
	writeFile(t, report, borrowedEvidenceReport("[run](evidence/run.txt)\n"+paths[0]+"/.venv/bin/tool\n"))
	assertBorrowedEvidence(t, report, options, []BorrowedEvidence{{5, filepath.Base(report), BorrowedEvidenceForeignWorktree + " " + paths[0]}})
	// Explicit Run provenance also exempts the parent and Task worktrees.
	options.RunBranch = "roundfix/run-own"
	writeFile(t, report, borrowedEvidenceReport(paths[1]+"/bin/tool\n"))
	assertBorrowedEvidence(t, report, options, nil)
}

func TestBorrowedEvidenceNamesImportedEvidence(t *testing.T) {
	t.Parallel()
	root, report := borrowedEvidenceRepo(t)
	writeFile(t, filepath.Join(root, "initial.txt"), "initial\n")
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "-m", "audited head")
	head := strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
	file := filepath.Join(filepath.Dir(report), "evidence/imported.txt")
	writeFile(t, file, "imported\n")
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "-m", "import")
	commit := strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
	rel, err := filepath.Rel(root, file)
	if err != nil {
		t.Fatal(err)
	}
	options := BorrowedEvidenceOptions{RunBranchPrefix: "roundfix/run-", AuditedHead: head, ImportCommit: commit, ImportedPaths: []string{filepath.ToSlash(rel)}}
	writeFile(t, report, borrowedEvidenceReport("[imported](evidence/imported.txt)\n"))
	assertBorrowedEvidence(t, report, options, []BorrowedEvidence{{4, "evidence/imported.txt", BorrowedEvidenceImported}})
	options.ImportedPaths = nil
	assertBorrowedEvidence(t, report, options, nil)
	options.ImportedPaths = []string{filepath.ToSlash(rel)}
	writeFile(t, file, "regenerated\n")
	assertBorrowedEvidence(t, report, options, nil)
	options.AuditedHead = "not-a-commit"
	if items, err := FindBorrowedEvidence(t.Context(), report, options); err == nil || items != nil {
		t.Fatalf("bad revision = %+v, %v", items, err)
	}
}

func TestBorrowedEvidenceErrorNamesEachItem(t *testing.T) {
	t.Parallel()
	err := BorrowedEvidenceError{Evidence: []BorrowedEvidence{{9, "evidence/a.txt", BorrowedEvidenceReused}, {12, "/elsewhere/tool", BorrowedEvidenceOutsideWorktree}}}
	want := "QA Report evidence does not belong to this Run: evidence/a.txt (line 9, reused); /elsewhere/tool (line 12, outside the worktree)"
	if got := err.Error(); got != want {
		t.Fatalf("error = %q; want %q", got, want)
	}
}

func TestBorrowedEvidenceOutsideGitChecksNothing(t *testing.T) {
	t.Parallel()
	report := filepath.Join(t.TempDir(), "qa/report.md")
	writeFile(t, report, "---\nverdict: pass\nuser_flow_executables: /outside/tool\n---\n[absent](evidence/missing.txt)\n")
	options := BorrowedEvidenceOptions{RunBranchPrefix: "roundfix/run-"}
	assertBorrowedEvidence(t, report, options, nil)
	if _, err := FindBorrowedEvidence(t.Context(), report, BorrowedEvidenceOptions{}); err == nil {
		t.Fatal("empty RunBranchPrefix accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if items, err := FindBorrowedEvidence(ctx, report, options); err == nil || items != nil {
		t.Fatalf("cancelled read = %+v, %v", items, err)
	}
	if err := os.Remove(report); err != nil {
		t.Fatal(err)
	}
	if items, err := FindBorrowedEvidence(t.Context(), report, options); !errors.Is(err, os.ErrNotExist) || items != nil {
		t.Fatalf("missing report = %+v, %v", items, err)
	}
}
