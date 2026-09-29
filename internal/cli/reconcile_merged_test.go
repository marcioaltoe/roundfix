// Suite: reconcile merged-Spec terminal Runs.
// Invariant: reconcile releases a Run only when the Delivery Queue record or archived default branch proves its commits represented at the merged head.
// Boundary IN: public reconcile runner, real Git repositories, persisted Delivery Queues, Run metadata, Task status, and QA Reports.
// Boundary OUT: merged-head commit representation internals, owned by internal/worktree/merged_head_test.go.
package cli

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

const (
	reconcileMergedSpecSlug = "0172-reconcile-merged-spec"
	reconcileMergedTaskID   = "task_05"
)

func TestReconcileReleasesARunProvenByTheDeliveryMergeRecord(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	seedReconcileDeliveryMergedHead(t, fixture, reconcileMergedSpecSlug)

	report := runReconcileMergedCommand(t, fixture, "--apply")

	result := reconcileMergedResult(t, report, fixture.run.ID)
	if result.Classification != "superseded" || result.Action != "released" ||
		!strings.Contains(result.Evidence, "pull request #259 merged head") {
		t.Fatalf("merged-record apply result = %+v", result)
	}
	assertReconcilePathState(t, fixture.ref.Path, false)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, false)
}

func TestReconcileDryRunLeavesAMergedSpecRunInPlace(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	seedReconcileDeliveryMergedHead(t, fixture, reconcileMergedSpecSlug)

	report := runReconcileMergedCommand(t, fixture)

	result := reconcileMergedResult(t, report, fixture.run.ID)
	if result.Classification != "superseded" || result.Action != "would release with --apply" ||
		!strings.Contains(result.Evidence, "pull request #259 merged head") {
		t.Fatalf("merged-record dry-run result = %+v", result)
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func TestReconcileIgnoresAMergeRecordForAnotherSpec(t *testing.T) {
	t.Parallel()
	const otherSlug = "0174-another-spec"
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: otherSlug,
		taskStatus:   "completed",
	})
	seedReconcileDeliveryMergedHead(t, fixture, otherSlug)

	report := runReconcileMergedCommand(t, fixture, "--apply")

	result := reconcileMergedResult(t, report, fixture.run.ID)
	if result.Classification != "unintegrated" || result.Action != "preserve" {
		t.Fatalf("other-Spec merge-record result = %+v", result)
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func TestReconcilePreservesARunTheMergeRecordDoesNotRepresent(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "pending",
	})
	seedReconcileDeliveryMergedHead(t, fixture, reconcileMergedSpecSlug)

	report := runReconcileMergedCommand(t, fixture, "--apply")

	result := reconcileMergedResult(t, report, fixture.run.ID)
	if result.Classification != "unintegrated" || result.Action != "preserve" ||
		!strings.Contains(result.RefusalReason, fixture.taskCommit[:12]) {
		t.Fatalf("unrepresented merged-record result = %+v, want commit %s", result, fixture.taskCommit[:12])
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func TestReconcileKeepsTheApplyActionForAMergedRunWithASupersededDisposition(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
		targetExists: true,
	})
	seedReconcileDeliveryMergedHead(t, fixture, reconcileMergedSpecSlug)
	seedLaterReconcileRunCoverage(t, fixture)

	report := runReconcileMergedCommand(t, fixture)

	result := reconcileMergedResult(t, report, fixture.run.ID)
	if result.Classification != "superseded" || result.Action != "would release with --apply" ||
		!strings.Contains(result.Evidence, "pull request #259 merged head") {
		t.Fatalf("merged-head result replaced by Branch Disposition = %+v", result)
	}
}

func TestReconcileReleasesAnArchivedSpecRunWithoutARecord(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})

	report := runReconcileMergedCommand(t, fixture, "--apply")

	result := reconcileMergedResult(t, report, fixture.run.ID)
	if result.Classification != "superseded" || result.Action != "released" ||
		!strings.Contains(result.Evidence, `default branch "main"`) {
		t.Fatalf("archived default-branch apply result = %+v", result)
	}
	assertReconcilePathState(t, fixture.ref.Path, false)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, false)
}

func TestReconcileReadsMergeRecordFromTheLoadedCheckoutPath(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	alias := filepath.Join(t.TempDir(), "checkout")
	if err := os.Symlink(fixture.repoDir, alias); err != nil {
		t.Fatalf("create checkout symlink: %v", err)
	}
	setCommandEnvironmentForTest(t, fixture.homeDir, alias)
	seedReconcileDeliveryMergedHeadAtRoot(t, fixture, reconcileMergedSpecSlug, alias)

	report := runReconcileMergedCommand(t, fixture)

	result := reconcileMergedResult(t, report, fixture.run.ID)
	if result.Classification != "superseded" ||
		!strings.Contains(result.Evidence, "pull request #259 merged head") {
		t.Fatalf("loaded-checkout merge-record result = %+v", result)
	}
}

func TestReconcileFailsWhenDeliveryMergeRecordsCannotBeRead(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug,
		taskStatus:   "completed",
		includeQA:    true,
	})
	seedReconcileDeliveryMergedHead(t, fixture, reconcileMergedSpecSlug)
	database, err := sql.Open("sqlite", "file:"+store.DatabasePath(fixture.homeDir))
	if err != nil {
		t.Fatalf("open Run Database for corruption fixture: %v", err)
	}
	if _, err := database.ExecContext(
		context.Background(),
		`UPDATE delivery_queue_items SET candidate_commits = 'not-json' WHERE git_root = ?`,
		fixture.repoDir,
	); err != nil {
		_ = database.Close()
		t.Fatalf("corrupt Delivery Queue candidate commits: %v", err)
	}
	if err := database.Close(); err != nil {
		t.Fatalf("close corrupted Run Database fixture: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", fixture.run.ID, "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitRunFailed || stdout.Len() != 0 ||
		!strings.Contains(stderr.String(), "read Delivery Queue merge records") {
		t.Fatalf("unreadable merge-record result: exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

type reconcileMergedFixtureOptions struct {
	archivedSlug string
	taskStatus   string
	includeQA    bool
	targetExists bool
}

type reconcileMergedFixture struct {
	homeDir    string
	repoDir    string
	run        store.Run
	ref        runworktree.Ref
	mergedHead string
	taskCommit string
}

func newReconcileMergedFixture(t *testing.T, opts reconcileMergedFixtureOptions) reconcileMergedFixture {
	t.Helper()
	homeDir, repoDir, location := newReconcileWorkspace(t)
	ctx := context.Background()
	startHead := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
	targetBranch := "delivery/" + reconcileMergedSpecSlug
	if opts.targetExists {
		gitImplement(t, repoDir, "branch", targetBranch, startHead)
	}
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open merged reconcile Run Database: %v", err)
	}
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: targetBranch,
		HeadSHA:     startHead,
		SpecSlug:    reconcileMergedSpecSlug,
		Agent:       "codex",
	})
	if err != nil {
		_ = runStore.Close()
		t.Fatalf("create merged reconcile Run: %v", err)
	}
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: repoDir,
		Location: location,
		RunID:    run.ID,
		HeadSHA:  startHead,
	})
	if err != nil {
		_ = runStore.Close()
		t.Fatalf("create merged reconcile Run Worktree: %v", err)
	}
	ref.Path, err = filepath.EvalSymlinks(ref.Path)
	if err != nil {
		_ = runStore.Close()
		t.Fatalf("resolve merged reconcile Run Worktree: %v", err)
	}
	run, err = runStore.SetRunWorkDir(ctx, run.ID, ref.Path)
	if err != nil {
		_ = runStore.Close()
		t.Fatalf("record merged reconcile Run Worktree: %v", err)
	}
	taskCommit := commitReconcileMergedTask(t, ref.Path, reconcileMergedSpecSlug, reconcileMergedTaskID)
	if opts.includeQA {
		commitReconcileMergedQAReport(t, ref.Path, reconcileMergedSpecSlug, "qa-report-2026-09-28.md", "fail")
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, store.StateStopped)
	if err != nil {
		_ = runStore.Close()
		t.Fatalf("complete merged reconcile Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close merged reconcile Run Database: %v", err)
	}

	gitImplement(t, repoDir, "checkout", "main")
	mergedHead := commitReconcileMergedSpec(t, repoDir, opts.archivedSlug, opts.taskStatus, opts.includeQA)
	return reconcileMergedFixture{
		homeDir:    homeDir,
		repoDir:    repoDir,
		run:        completed.Run,
		ref:        ref,
		mergedHead: mergedHead,
		taskCommit: taskCommit,
	}
}

func commitReconcileMergedTask(t *testing.T, workDir, slug, taskID string) string {
	t.Helper()
	mustWrite(t, filepath.Join(workDir, "feature.txt"), "delivered\n")
	gitImplement(t, workDir, "add", "feature.txt")
	gitImplement(
		t,
		workDir,
		"commit",
		"-m", fmt.Sprintf("feat: settle %s", taskID),
		"-m", fmt.Sprintf("Roundfix-Spec: %s\nRoundfix-Task: %s", slug, taskID),
	)
	return strings.TrimSpace(gitImplementOutput(t, workDir, "rev-parse", "HEAD"))
}

func commitReconcileMergedQAReport(t *testing.T, workDir, slug, reportName, verdict string) {
	t.Helper()
	report := filepath.Join("docs", "specs", slug, "qa", reportName)
	mustMkdir(t, filepath.Join(workDir, filepath.Dir(report)))
	mustWrite(t, filepath.Join(workDir, report), "QA report\n")
	gitImplement(t, workDir, "add", filepath.ToSlash(report))
	gitImplement(
		t,
		workDir,
		"commit",
		"-m", fmt.Sprintf("docs: qa report for %s (%s)\n\nRoundfix-Spec: %s", slug, verdict, slug),
	)
}

func commitReconcileMergedSpec(t *testing.T, workDir, slug, taskStatus string, includeQA bool) string {
	t.Helper()
	if slug == "" {
		mustWrite(t, filepath.Join(workDir, "merged-head.txt"), "merged\n")
		gitImplement(t, workDir, "add", "merged-head.txt")
	} else {
		root := filepath.Join(workDir, "docs", "history", "specs", slug)
		mustMkdir(t, root)
		mustWrite(t, filepath.Join(root, "_tasks.md"), fmt.Sprintf(
			"---\nschema: spec-tasks/v1\nspec: %s\ngraph:\n  nodes:\n    - id: %s\n      file: %s.md\n      needs: []\n---\n\n# Task Graph\n",
			slug,
			reconcileMergedTaskID,
			reconcileMergedTaskID,
		))
		mustWrite(t, filepath.Join(root, reconcileMergedTaskID+".md"), fmt.Sprintf(
			"---\ntask: %s\nspec: %s\nstatus: %s\ntype: backend\ncomplexity: low\n---\n\n# Task 05\n\n## Verification\n\n- `true`\n",
			reconcileMergedTaskID,
			slug,
			taskStatus,
		))
		if includeQA {
			mustMkdir(t, filepath.Join(root, "qa"))
			mustWrite(t, filepath.Join(root, "qa", "qa-report-2026-09-29.md"), "QA report\n")
		}
		gitImplement(t, workDir, "add", filepath.ToSlash(filepath.Join("docs", "history", "specs", slug)))
	}
	gitImplement(t, workDir, "commit", "-m", "feat: record merged Spec")
	return strings.TrimSpace(gitImplementOutput(t, workDir, "rev-parse", "HEAD"))
}

func seedReconcileDeliveryMergedHead(t *testing.T, fixture reconcileMergedFixture, slug string) {
	t.Helper()
	seedReconcileDeliveryMergedHeadAtRoot(t, fixture, slug, fixture.repoDir)
}

func seedReconcileDeliveryMergedHeadAtRoot(
	t *testing.T,
	fixture reconcileMergedFixture,
	slug string,
	gitRoot string,
) {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, fixture.homeDir)
	if err != nil {
		t.Fatalf("open Delivery Queue store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close Delivery Queue store: %v", err)
		}
	}()
	queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{slug})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageMerged
	item.Branch = "roundfix/delivery-" + slug
	item.CandidateCommits = []string{fixture.run.HeadSHA, fixture.mergedHead}
	item.PullRequestNumber = "259"
	item.MergeCommit = fixture.mergedHead
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("record Delivery Queue merge: %v", err)
	}
}

func seedLaterReconcileRunCoverage(t *testing.T, fixture reconcileMergedFixture) {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, fixture.homeDir)
	if err != nil {
		t.Fatalf("open later-Run coverage store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close later-Run coverage store: %v", err)
		}
	}()
	appendCompletedTaskEvidence(t, runStore, fixture.run.ID, reconcileMergedTaskID)
	later, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     fixture.repoDir,
		LocalBranch: fixture.run.LocalBranch,
		HeadSHA:     fixture.mergedHead,
		SpecSlug:    reconcileMergedSpecSlug,
		Agent:       "codex",
	})
	if err != nil {
		t.Fatalf("create later covering Run: %v", err)
	}
	appendCompletedTaskEvidence(t, runStore, later.ID, reconcileMergedTaskID)
	if _, err := runStore.CompleteRun(ctx, later.ID, store.StateClean); err != nil {
		t.Fatalf("complete later covering Run: %v", err)
	}
}

func appendCompletedTaskEvidence(t *testing.T, runStore *store.Store, runID, taskID string) {
	t.Helper()
	payload := []byte(fmt.Sprintf(`{"task":%q,"phase":"settled","status":"completed"}`, taskID))
	if _, err := runStore.AppendRunEvents(context.Background(), []runevent.RunEvent{{
		RunID:       runID,
		Source:      runevent.SourceDaemon,
		Kind:        runevent.KindDaemonTask,
		ReviewIssue: taskID,
		Payload:     payload,
	}}); err != nil {
		t.Fatalf("append completed Task evidence for Run %s: %v", runID, err)
	}
}

func runReconcileMergedCommand(t *testing.T, fixture reconcileMergedFixture, extra ...string) reconcileReport {
	t.Helper()
	args := []string{"reconcile", fixture.run.ID}
	args = append(args, extra...)
	args = append(args, "--format=json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), args, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("reconcile merged Run exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("reconcile merged Run stderr = %q, want empty", stderr.String())
	}
	var report reconcileReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode merged reconciliation report: %v\n%s", err, stdout.String())
	}
	return report
}

func reconcileMergedResult(t *testing.T, report reconcileReport, runID string) reconcileResult {
	t.Helper()
	if len(report.Results) != 1 || report.Results[0].RunID != runID {
		t.Fatalf("merged reconciliation results = %+v, want Run %q", report.Results, runID)
	}
	return report.Results[0]
}
