package worktree

// Suite: Archive Record reconciliation after squash delivery.
// Boundary IN: real Git objects and disposable Run worktrees.
// Boundary OUT: network and daemon settlement.

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func newArchiveRecordRun(t *testing.T, missing bool, overrideStatus string) (terminalRunFixture, spec.ArchiveRecord, string) {
	t.Helper()
	f := newMergedHeadTestFixture(t, "archive-record-run")
	f.run.LocalBranch = "delivery/record-item"
	gitWorktreeTest(t, f.repoDir, "branch", f.run.LocalBranch, "main")
	source := "docs/specs/" + f.run.SpecSlug
	commitMergedHeadTask(t, f.ref.Path, f.run.SpecSlug, "task_01", "feature.txt", "delivered\n")
	commitMergedHeadFiles(t, f.ref.Path, map[string]string{
		source + "/_prd.md":                    "---\nstatus: approved\n---\n\n# Record Spec\n",
		source + "/_tasks.md":                  "---\nschema: spec-tasks/v1\nqa: task_02\ngraph:\n  nodes:\n    - {id: task_01, file: task_01.md}\n    - {id: task_02, file: task_02.md}\n---\n",
		source + "/task_01.md":                 "---\ntask: task_01\nspec: " + f.run.SpecSlug + "\nstatus: completed\ntype: backend\ncomplexity: low\n---\n\n# Task\n\n## Context\n\n- interface: `declared.txt`\n- instruction: `instruction.txt`\n\n## Recorded paths\n\n- `recorded.txt`\n",
		source + "/task_02.md":                 "---\ntask: task_02\nspec: " + f.run.SpecSlug + "\nstatus: completed\ntype: qa\ncomplexity: low\n---\n",
		source + "/qa/qa-report-2026-10-06.md": "---\nverdict: pass\n---\n",
	}, "docs: source Spec")
	r := spec.ArchiveRecord{Spec: f.run.SpecSlug, Title: "Record Spec", Source: source, SourceRevision: mergedHeadTestHead(t, f.ref.Path, "HEAD"), Disposition: spec.ArchivePass, QATask: "task_02", QAReport: "qa-report-2026-10-06.md", QAVerdict: "pass"}
	if missing {
		r.SourceRevision = strings.Repeat("a", 40)
	}
	if overrideStatus != "none" {
		r.Disposition = spec.ArchiveQAOverride
		r.QAOverride = &spec.QAArchiveOverrideRecord{Approval: "operator", Reason: "environment", QATaskStatus: overrideStatus}
	}
	content, err := spec.RenderArchiveRecord(r)
	if err != nil {
		t.Fatal(err)
	}
	gitWorktreeTest(t, f.ref.Path, "rm", "-r", "--", source)
	commitMergedHeadFiles(t, f.ref.Path, map[string]string{spec.ArchiveRecordPath(spec.ArchiveDir(spec.ArchiveKindSpec), f.run.SpecSlug): string(content)}, "docs: archive record")
	gitWorktreeTest(t, f.repoDir, "merge", "--squash", f.ref.Branch)
	gitWorktreeTest(t, f.repoDir, "commit", "-m", "feat: squash delivery")
	delivery := mergedHeadTestHead(t, f.repoDir, "main")
	// Main never contained the Spec folder, even though the object exists on
	// the Run branch. Later default changes force the delivery proof.
	if out := gitWorktreeTest(t, f.repoDir, "log", "main", "--format=%H", "--", source); strings.TrimSpace(out) != "" {
		t.Fatalf("main held source folder: %s", out)
	}
	commitWorktreeFile(t, f.repoDir, "feature.txt", "later default change\n", "fix: advance main")
	return f, r, delivery
}

func TestMergeEvidenceFromTheArchiveRecordAfterASquash(t *testing.T) {
	t.Parallel()
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "source present", true: "source absent"}[missing], func(t *testing.T) {
			f, _, delivered := newArchiveRecordRun(t, missing, "none")
			head := mergedHeadTestHead(t, f.ref.Path, "HEAD")
			proof, ok := ProveDelivery(t.Context(), f.repoDir, f.run.SpecSlug, head)
			if !missing && (!ok || proof.DeliveryCommit != delivered) {
				t.Fatalf("delivery proof=%+v proven=%v", proof, ok)
			}
			result := inspectMergedHeadTestRun(t, f, nil)
			if missing {
				assertTerminalRunReconciliation(t, result, f.run, ReconciliationUnintegrated)
				if !strings.Contains(result.Reason, "source_revision") {
					t.Fatalf("missing revision reason=%q", result.Reason)
				}
				assertMergeEvidencePreserved(t, f, result)
			} else {
				assertTerminalRunReconciliation(t, result, f.run, ReconciliationSuperseded)
				if err := ApplyTerminalRun(context.Background(), &recordingTerminalRunStore{}, result); err != nil {
					t.Fatal(err)
				}
				assertPathRemoved(t, f.ref.Path)
			}
		})
	}
}

func TestTaskCompletionFollowsTheArchiveRecord(t *testing.T) {
	t.Parallel()
	for _, status := range []string{"none", "", "pending", "failed"} {
		t.Run("override status "+status, func(t *testing.T) {
			for _, missing := range []bool{false, true} {
				f, _, _ := newArchiveRecordRun(t, missing, status)
				head := mergedHeadTestHead(t, f.repoDir, "main")
				if !taskCompletedAtMergedHead(t.Context(), execGitRunner{}, f.repoDir, head, f.run.SpecSlug, "task_01") {
					t.Fatal("non-QA Task not completed by record")
				}
				wantQA := status == "none" || status == ""
				if got := taskCompletedAtMergedHead(t.Context(), execGitRunner{}, f.repoDir, head, f.run.SpecSlug, "task_02"); got != wantQA {
					t.Fatalf("QA completed=%v want=%v", got, wantQA)
				}
			}
		})
	}
}

func TestLeftoversReadTaskContextAtTheSourceRevision(t *testing.T) {
	t.Parallel()
	for _, file := range []string{"declared.txt", "recorded.txt", "instruction.txt", "outside.txt"} {
		t.Run(file, func(t *testing.T) {
			f, _, _ := newArchiveRecordRun(t, false, "none")
			mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, file), "leftover\n")
			result := inspectMergedHeadTestRun(t, f, nil)
			if file == "declared.txt" || file == "recorded.txt" {
				assertTerminalRunReconciliation(t, result, f.run, ReconciliationSuperseded)
				if err := ApplyTerminalRun(t.Context(), &recordingTerminalRunStore{}, result); err != nil {
					t.Fatal(err)
				}
				assertPathRemoved(t, f.ref.Path)
			} else {
				assertTerminalRunReconciliation(t, result, f.run, ReconciliationDirty)
				assertMergeEvidencePreserved(t, f, result)
			}
		})
	}
}

func TestLeftoversKeepTheWorktreeWithoutTheSourceRevision(t *testing.T) {
	t.Parallel()
	for _, file := range []string{"declared.txt", "docs/specs/" + mergedHeadTestSlug + "/leftover.md"} {
		t.Run(file, func(t *testing.T) {
			f, _, _ := newArchiveRecordRun(t, true, "none")
			if err := os.MkdirAll(filepath.Dir(filepath.Join(f.ref.Path, file)), 0755); err != nil {
				t.Fatal(err)
			}
			mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, file), "preserve\n")
			result := inspectMergedHeadTestRun(t, f, nil)
			assertTerminalRunReconciliation(t, result, f.run, ReconciliationDirty)
			if !strings.Contains(result.Reason, "source_revision") {
				t.Fatalf("reason=%q", result.Reason)
			}
			assertMergeEvidencePreserved(t, f, result)
		})
	}
}

func TestQASupersessionReadsTheArchiveRecord(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		missing bool
		report  string
		content string
		want    bool
	}{
		{"newer report", false, "qa-report-2026-10-05.md", "---\nverdict: fail\n---\n", true},
		{"identical report", false, "qa-report-2026-10-06.md", "---\nverdict: pass\n---\n", true},
		{"different report", false, "qa-report-2026-10-06.md", "---\nverdict: fail\n---\n", false},
		{"source absent", true, "qa-report-2026-10-05.md", "---\nverdict: fail\n---\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, r, _ := newArchiveRecordRun(t, tc.missing, "none")
			gitWorktreeTest(t, f.ref.Path, "reset", "--hard", "main~2")
			commitMergedHeadFiles(t, f.ref.Path, map[string]string{r.Source + "/qa/" + tc.report: tc.content}, "docs: qa report for "+f.run.SpecSlug+" (fail)\n\nRoundfix-Spec: "+f.run.SpecSlug)
			report, proven := SupersedingQAReport(t.Context(), f.repoDir, mergedHeadTestHead(t, f.repoDir, "main"), mergedHeadTestHead(t, f.ref.Path, "HEAD"), f.run.SpecSlug)
			if proven != tc.want || (tc.want && report != r.Source+"/qa/"+r.QAReport) {
				t.Fatalf("report=%s proven=%v want=%v", report, proven, tc.want)
			}
		})
	}
}
