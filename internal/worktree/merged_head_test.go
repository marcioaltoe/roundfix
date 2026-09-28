package worktree

// Suite: merged-head terminal Run reconciliation.
// Invariant: cleanup releases a Run only when every commit absent from the merged head is represented there.
// Boundary IN: real Git objects, Run Worktrees, Run Branches, Spec Task Graphs, and QA Reports.
// Boundary OUT: Delivery Queue record loading and daemon Branch Disposition classification.

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

const mergedHeadTestSlug = "0175-cleanup-after-a-squash-merge"

func TestMergedHeadRecordReleasesARunContainedInTheMergedHead(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-contained")
	runHead := fixture.commitRunChange(t, "contained.txt", "contained\n")
	targetHead := mergedHeadTestHead(t, fixture.repoDir, "main")

	result := inspectMergedHeadTestRun(t, fixture, []MergedHead{mergedHeadTestRecord(fixture, runHead)})

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationSafe)
	if result.TargetHead != targetHead || !strings.Contains(result.Reason, "Run Branch is contained in") {
		t.Fatalf("contained merged-head result = %#v, want target head %s and containment proof", result, targetHead)
	}
}

func TestMergedHeadRecordSupersedesARedoneTaskAndAFailedQAReport(t *testing.T) {
	t.Parallel()
	fixture, mergedHead := newMergedHeadTaskAndQAReportFixture(t, "merged-head-record-supersedes")

	result := inspectMergedHeadTestRun(t, fixture, []MergedHead{mergedHeadTestRecord(fixture, mergedHead)})

	assertMergedHeadTaskAndQAReportSuperseded(t, result, "pull request #259 merged head")
}

func TestMergedHeadDefaultBranchSupersedesARedoneTaskAndAFailedQAReport(t *testing.T) {
	t.Parallel()
	fixture, _ := newMergedHeadTaskAndQAReportFixture(t, "merged-head-default-supersedes")

	result := inspectMergedHeadTestRun(t, fixture, nil)

	assertMergedHeadTaskAndQAReportSuperseded(t, result, `default branch "main"`)
}

func TestMergedHeadDefaultBranchComparesOnlyTheRunsChangedFiles(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-changed-files")
	fixture.run.LocalBranch = "delivery/merged-head-changed-files"
	fixture.commitRunChange(t, "represented.txt", "represented\n")
	commitMergedHeadSpec(t, fixture.repoDir, fixture.run.SpecSlug, true, nil, "qa-report-2026-09-28.md", map[string]string{
		"represented.txt": "represented\n",
	})
	commitWorktreeFile(t, fixture.repoDir, "default-only.txt", "later\n", "advance unrelated default content")

	result := inspectMergedHeadTestRun(t, fixture, nil)

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationSafe)
	if result.Reason != `Run Branch content is fully represented on default branch "main"` {
		t.Fatalf("content representation reason = %q", result.Reason)
	}
}

func TestMergedHeadDefaultBranchReleasesAnArchivedSpecRunAfterMainMoved(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-main-moved")
	fixture.run.LocalBranch = "delivery/merged-head-main-moved"
	commitMergedHeadTask(t, fixture.ref.Path, fixture.run.SpecSlug, "task_05", "shared.txt", "delivered\n")
	commitMergedHeadSpec(t, fixture.repoDir, fixture.run.SpecSlug, true, map[string]string{
		"task_05": "completed",
	}, "qa-report-2026-09-28.md", map[string]string{
		"shared.txt": "delivered\n",
	})
	commitWorktreeFile(t, fixture.repoDir, "shared.txt", "changed after delivery\n", "advance a delivered file")

	result := inspectMergedHeadTestRun(t, fixture, nil)

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationSuperseded)
	if !strings.Contains(result.Reason, "1 Task commit completed") {
		t.Fatalf("moved-main proof = %q, want completed Task representation", result.Reason)
	}
}

func TestMergedHeadRefusesAnUnrepresentedCommit(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-unrepresented")
	fixture.commitRunChange(t, "unique.txt", "preserve\n")
	head := mergedHeadTestHead(t, fixture.repoDir, "main")

	result := inspectMergedHeadTestRun(t, fixture, []MergedHead{mergedHeadTestRecord(fixture, head)})

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnintegrated)
	if !strings.Contains(result.Reason, "1 Run-only file") {
		t.Fatalf("unrepresented commit reason = %q", result.Reason)
	}
}

func TestMergedHeadRefusesATaskNotCompletedAtTheMergedHead(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-task-pending")
	commit := commitMergedHeadTask(t, fixture.ref.Path, fixture.run.SpecSlug, "task_05", "task.txt", "work\n")
	_ = commit
	head := commitMergedHeadSpec(t, fixture.repoDir, fixture.run.SpecSlug, false, map[string]string{
		"task_05": "pending",
	}, "", nil)

	result := inspectMergedHeadTestRun(t, fixture, []MergedHead{mergedHeadTestRecord(fixture, head)})

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnintegrated)
	if !strings.Contains(result.Reason, "Task task_05 is not completed") || !strings.Contains(result.Reason, commit[:12]) {
		t.Fatalf("pending Task reason = %q, want Task and commit", result.Reason)
	}
}

func TestMergedHeadRefusesATaskCommitOfAnotherSpec(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-other-spec-task")
	const otherSlug = "0174-another-spec"
	commit := commitMergedHeadTask(t, fixture.ref.Path, otherSlug, "task_05", "other.txt", "work\n")
	head := commitMergedHeadSpec(t, fixture.repoDir, fixture.run.SpecSlug, false, map[string]string{
		"task_05": "completed",
	}, "", nil)

	result := inspectMergedHeadTestRun(t, fixture, []MergedHead{mergedHeadTestRecord(fixture, head)})

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnintegrated)
	if !strings.Contains(result.Reason, "it belongs to Spec "+otherSlug) || !strings.Contains(result.Reason, commit[:12]) {
		t.Fatalf("other-Spec Task reason = %q", result.Reason)
	}
}

func TestMergedHeadRefusesAQAReportTheMergedHeadDoesNotSupersede(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-unsuperseded-qa")
	commit := commitQAReport(t, fixture.ref.Path, fixture.run.SpecSlug, "qa-report-2026-09-29.md", false, "fail")
	head := commitMergedHeadSpec(t, fixture.repoDir, fixture.run.SpecSlug, true, nil, "qa-report-2026-09-28.md", nil)

	result := inspectMergedHeadTestRun(t, fixture, []MergedHead{mergedHeadTestRecord(fixture, head)})

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnintegrated)
	if !strings.Contains(result.Reason, "its QA Report is not superseded") || !strings.Contains(result.Reason, commit[:12]) {
		t.Fatalf("unsuperseded QA reason = %q", result.Reason)
	}
}

func TestMergedHeadIgnoresARecordForAnotherSpec(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-ignore-other-spec")
	runHead := fixture.commitRunChange(t, "unique.txt", "preserve\n")
	record := mergedHeadTestRecord(fixture, runHead)
	record.SpecSlug = "0174-another-spec"

	result := inspectMergedHeadTestRun(t, fixture, []MergedHead{record})

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnintegrated)
	if result.Reason != reconciliationReasonUnintegrated {
		t.Fatalf("other-Spec record reason = %q, want unchanged unintegrated reason", result.Reason)
	}
}

func TestMergedHeadFallsBackWhenRecordsDisagree(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-disagree")
	runHead := commitMergedHeadTask(t, fixture.ref.Path, fixture.run.SpecSlug, "task_05", "task.txt", "work\n")
	defaultHead := commitMergedHeadSpec(t, fixture.repoDir, fixture.run.SpecSlug, true, map[string]string{
		"task_05": "completed",
	}, "qa-report-2026-09-28.md", nil)
	records := []MergedHead{
		mergedHeadTestRecord(fixture, runHead),
		mergedHeadTestRecord(fixture, defaultHead),
	}

	result := inspectMergedHeadTestRun(t, fixture, records)

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationSuperseded)
	if !strings.Contains(result.Reason, `default branch "main"`) {
		t.Fatalf("disagreeing-record fallback reason = %q", result.Reason)
	}
}

func TestMergedHeadFallsBackWhenRecordCommitIsMissing(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-missing-record")
	commitMergedHeadTask(t, fixture.ref.Path, fixture.run.SpecSlug, "task_05", "task.txt", "work\n")
	commitMergedHeadSpec(t, fixture.repoDir, fixture.run.SpecSlug, true, map[string]string{
		"task_05": "completed",
	}, "qa-report-2026-09-28.md", nil)
	record := mergedHeadTestRecord(fixture, strings.Repeat("f", 40))

	result := inspectMergedHeadTestRun(t, fixture, []MergedHead{record})

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationSuperseded)
	if !strings.Contains(result.Reason, `default branch "main"`) {
		t.Fatalf("missing-record fallback reason = %q", result.Reason)
	}
}

func TestMergedHeadKeepsADirtyRunWorktree(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-dirty")
	runHead := fixture.commitRunChange(t, "contained.txt", "contained\n")
	mustWriteWorktreeTest(t, filepath.Join(fixture.ref.Path, "dirty.txt"), "preserve\n")

	result := inspectMergedHeadTestRun(t, fixture, []MergedHead{mergedHeadTestRecord(fixture, runHead)})

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationDirty)
	if result.Reason != reconciliationReasonDirty {
		t.Fatalf("dirty Run reason = %q", result.Reason)
	}
}

func TestMergedHeadKeepsTheUnintegratedReasonWithoutASource(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-no-source")
	fixture.commitRunChange(t, "unique.txt", "preserve\n")

	result := inspectMergedHeadTestRun(t, fixture, nil)

	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnintegrated)
	if result.Reason != reconciliationReasonUnintegrated {
		t.Fatalf("no-source reason = %q, want %q", result.Reason, reconciliationReasonUnintegrated)
	}
}

func TestMergedHeadApplyReleasesTheRunWorktreeAndBranch(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-apply")
	runHead := fixture.commitRunChange(t, "contained.txt", "contained\n")
	result := inspectMergedHeadTestRun(t, fixture, []MergedHead{mergedHeadTestRecord(fixture, runHead)})

	if err := ApplyTerminalRun(context.Background(), &recordingTerminalRunStore{}, result); err != nil {
		t.Fatalf("apply merged-head reconciliation: %v", err)
	}

	assertPathRemoved(t, fixture.ref.Path)
	assertBranchRemoved(t, fixture.repoDir, fixture.ref.Branch)
}

func TestMergedHeadApplyRefusesAChangedRecord(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-apply-record-changed")
	runHead := fixture.commitRunChange(t, "contained.txt", "contained\n")
	result := inspectMergedHeadTestRun(t, fixture, []MergedHead{mergedHeadTestRecord(fixture, runHead)})
	result.evidence.merged[0].Head = mergedHeadTestHead(t, fixture.repoDir, "main")

	err := ApplyTerminalRun(context.Background(), &recordingTerminalRunStore{}, result)

	if err == nil || !strings.Contains(err.Error(), "merged-head record changed") {
		t.Fatalf("changed-record apply error = %v", err)
	}
	assertPathExists(t, fixture.ref.Path)
	assertRunBranchExists(t, fixture.repoDir, fixture.ref.Branch)
}

func TestMergedHeadApplyRefusesAMovedFallbackHead(t *testing.T) {
	t.Parallel()
	fixture := newMergedHeadTestFixture(t, "merged-head-apply-fallback-moved")
	fixture.run.LocalBranch = "delivery/stable-target"
	gitWorktreeTest(t, fixture.repoDir, "branch", fixture.run.LocalBranch, "main")
	fixture.commitRunChange(t, "represented.txt", "represented\n")
	commitMergedHeadSpec(t, fixture.repoDir, fixture.run.SpecSlug, true, nil, "qa-report-2026-09-28.md", map[string]string{
		"represented.txt": "represented\n",
	})
	result := inspectMergedHeadTestRun(t, fixture, nil)
	commitWorktreeFile(t, fixture.repoDir, "later.txt", "later\n", "move fallback head")

	err := ApplyTerminalRun(context.Background(), &recordingTerminalRunStore{}, result)

	if err == nil || !strings.Contains(err.Error(), "evidence is stale") {
		t.Fatalf("moved-fallback apply error = %v", err)
	}
	assertPathExists(t, fixture.ref.Path)
	assertRunBranchExists(t, fixture.repoDir, fixture.ref.Branch)
}

func newMergedHeadTestFixture(t *testing.T, runID string) terminalRunFixture {
	t.Helper()
	fixture := newTerminalRunFixture(t, runID)
	fixture.run.SpecSlug = mergedHeadTestSlug
	return fixture
}

func newMergedHeadTaskAndQAReportFixture(t *testing.T, runID string) (terminalRunFixture, string) {
	t.Helper()
	fixture := newMergedHeadTestFixture(t, runID)
	fixture.run.LocalBranch = "delivery/" + runID
	commitMergedHeadTask(t, fixture.ref.Path, fixture.run.SpecSlug, "task_05", "feature.txt", "delivered\n")
	commitQAReport(t, fixture.ref.Path, fixture.run.SpecSlug, "qa-report-2026-09-28.md", false, "fail")
	head := commitMergedHeadSpec(t, fixture.repoDir, fixture.run.SpecSlug, true, map[string]string{
		"task_05": "completed",
	}, "qa-report-2026-09-29.md", map[string]string{
		"feature.txt": "delivered\n",
	})
	return fixture, head
}

func mergedHeadTestRecord(fixture terminalRunFixture, head string) MergedHead {
	return MergedHead{
		SpecSlug:     fixture.run.SpecSlug,
		TargetBranch: fixture.run.LocalBranch,
		Head:         head,
		MergeCommit:  head,
		PullRequest:  "259",
	}
}

func inspectMergedHeadTestRun(t *testing.T, fixture terminalRunFixture, records []MergedHead) RunWorktreeReconciliation {
	t.Helper()
	result, err := InspectTerminalRunMerged(context.Background(), fixture.run, records)
	if err != nil {
		t.Fatalf("inspect terminal Run with merged heads: %v", err)
	}
	return result
}

func assertMergedHeadTaskAndQAReportSuperseded(t *testing.T, result RunWorktreeReconciliation, source string) {
	t.Helper()
	if result.State != ReconciliationSuperseded {
		t.Fatalf("merged-head state = %q, want superseded: %#v", result.State, result)
	}
	for _, fragment := range []string{source, "1 Task commit completed", "1 QA Report commit superseded"} {
		if !strings.Contains(result.Reason, fragment) {
			t.Fatalf("merged-head reason = %q, want fragment %q", result.Reason, fragment)
		}
	}
}

func commitMergedHeadTask(
	t *testing.T,
	workDir string,
	slug string,
	taskID string,
	file string,
	content string,
) string {
	t.Helper()
	return commitMergedHeadFiles(t, workDir, map[string]string{file: content}, fmt.Sprintf(
		"feat: settle %s\n\nRoundfix-Spec: %s\nRoundfix-Task: %s",
		taskID,
		slug,
		taskID,
	))
}

func commitMergedHeadSpec(
	t *testing.T,
	workDir string,
	slug string,
	archived bool,
	statuses map[string]string,
	reportName string,
	extra map[string]string,
) string {
	t.Helper()
	root := filepath.ToSlash(filepath.Join("docs", "specs", slug))
	if archived {
		root = filepath.ToSlash(filepath.Join(filepath.FromSlash(spec.ArchiveDir(spec.ArchiveKindSpec)), slug))
	}
	files := make(map[string]string, len(extra)+len(statuses)+2)
	for file, content := range extra {
		files[file] = content
	}
	taskIDs := make([]string, 0, len(statuses))
	for taskID := range statuses {
		taskIDs = append(taskIDs, taskID)
	}
	sort.Strings(taskIDs)
	var manifest strings.Builder
	manifest.WriteString("---\nschema: spec-tasks/v1\nspec: " + slug + "\ngraph:\n  nodes:\n")
	for _, taskID := range taskIDs {
		manifest.WriteString("    - id: " + taskID + "\n      file: " + taskID + ".md\n      needs: []\n")
		files[filepath.ToSlash(filepath.Join(root, taskID+".md"))] = fmt.Sprintf(
			"---\ntask: %s\nspec: %s\nstatus: %s\ntype: backend\ncomplexity: low\n---\n\n# %s\n\n## Verification\n\n- `true`\n",
			taskID,
			slug,
			statuses[taskID],
			taskID,
		)
	}
	manifest.WriteString("---\n\n# Task Graph\n")
	files[filepath.ToSlash(filepath.Join(root, "_tasks.md"))] = manifest.String()
	if reportName != "" {
		files[filepath.ToSlash(filepath.Join(root, "qa", reportName))] = "QA report\n"
	}
	return commitMergedHeadFiles(t, workDir, files, "feat: record merged Spec")
}

func commitMergedHeadFiles(t *testing.T, workDir string, files map[string]string, message string) string {
	t.Helper()
	paths := make([]string, 0, len(files))
	for file := range files {
		paths = append(paths, file)
	}
	sort.Strings(paths)
	for _, file := range paths {
		fullPath := filepath.Join(workDir, filepath.FromSlash(file))
		mustMkdirWorktreeTest(t, filepath.Dir(fullPath))
		mustWriteWorktreeTest(t, fullPath, files[file])
	}
	gitWorktreeTest(t, workDir, append([]string{"add", "--"}, paths...)...)
	messagePath := filepath.Join(t.TempDir(), "message.txt")
	mustWriteWorktreeTest(t, messagePath, message)
	gitWorktreeTest(t, workDir, "commit", "--cleanup=verbatim", "-F", messagePath)
	return mergedHeadTestHead(t, workDir, "HEAD")
}

func mergedHeadTestHead(t *testing.T, workDir string, revision string) string {
	t.Helper()
	return strings.TrimSpace(gitWorktreeTest(t, workDir, "rev-parse", revision))
}
