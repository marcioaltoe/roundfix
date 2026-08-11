// Suite: superseded Run Branch disposition.
// Invariant: a Run Branch is disposable only after Git or later-Run evidence proves all of its work superseded.
// Boundary IN: terminal Run metadata, local Git topology, completed-Task coverage, and the disposition record.
// Boundary OUT: CLI flag parsing and report rendering, owned by internal/cli/reconcile_test.go.
package daemon

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func TestSupersededBranchIsClassifiedWhenEveryCommitIsReachable(t *testing.T) {
	t.Parallel()
	fixture := newSupersededBranchFixture(t)
	first := fixture.commitRunChange(t, "first.txt", "first\n", "feat: add first change")
	second := fixture.commitRunChange(t, "second.txt", "second\n", "feat: add second change")
	gittest.Run(t, fixture.root, "merge", "--ff-only", fixture.ref.Branch)

	got, err := ClassifySupersededBranch(context.Background(), fixture.run, nil)

	if err != nil {
		t.Fatalf("ClassifySupersededBranch() error = %v", err)
	}
	if !got.Superseded || !got.Reachable {
		t.Fatalf("disposition = %+v, want reachable superseded", got)
	}
	if !slices.Equal(got.Commits, []string{first, second}) {
		t.Fatalf("commits = %v, want [%s %s]", got.Commits, first, second)
	}
	if !slices.Equal(got.ChangedFiles, []string{"first.txt", "second.txt"}) {
		t.Fatalf("changed files = %v", got.ChangedFiles)
	}
	if !strings.Contains(got.Reason, "every Run Branch commit is reachable") {
		t.Fatalf("reason = %q, want reachability proof", got.Reason)
	}
}

func TestSupersededBranchRefusesAnUnreachableCommit(t *testing.T) {
	t.Parallel()
	fixture := newSupersededBranchFixture(t)
	unreachable := fixture.commitRunChange(t, "run-only.txt", "run only\n", "feat: retain run-only change")
	fixture.commitTargetChange(t, "target-only.txt", "target only\n", "feat: advance target")

	got, err := ClassifySupersededBranch(context.Background(), fixture.run, nil)

	if err != nil {
		t.Fatalf("ClassifySupersededBranch() error = %v", err)
	}
	if got.Superseded {
		t.Fatalf("disposition = %+v, want refusal", got)
	}
	if !strings.Contains(got.RefusalReason, "unreachable commit "+unreachable) {
		t.Fatalf("refusal reason = %q, want unreachable commit %s", got.RefusalReason, unreachable)
	}
	assertSupersededBranchSurface(t, fixture, true)
}

func TestSupersededBranchIsClassifiedWhenLaterIntegratedRunCoveredTasks(t *testing.T) {
	t.Parallel()
	fixture := newSupersededBranchFixture(t)
	fixture.commitRunChange(
		t,
		"task.txt",
		"older task work\n",
		"feat: implement task\n\nRoundfix-Spec: disposition-spec\nRoundfix-Task: task_01",
	)
	fixture.commitTargetChange(t, "replacement.txt", "replacement work\n", "feat: integrate replacement")
	later := fixture.run
	later.ID = "later-run"
	later.State = store.StateClean
	later.CreatedAt = fixture.run.CreatedAt.Add(time.Hour)

	got, err := ClassifySupersededBranch(context.Background(), fixture.run, []RunTaskCoverage{
		{Run: fixture.run, CompletedTasks: []string{"task_01"}},
		{Run: later, CompletedTasks: []string{"task_01", "qa"}},
	})

	if err != nil {
		t.Fatalf("ClassifySupersededBranch() error = %v", err)
	}
	if !got.Superseded || got.Reachable {
		t.Fatalf("disposition = %+v, want later-Run superseded with unreachable commits", got)
	}
	if !strings.Contains(got.Reason, `later integrated Run "later-run" covered Tasks [task_01]`) {
		t.Fatalf("reason = %q, want later Run and covered Tasks", got.Reason)
	}
}

func TestWriteBranchDispositionRecordSurvivesDiscard(t *testing.T) {
	t.Parallel()
	fixture := newSupersededBranchFixture(t)
	fixture.commitRunChange(t, "work.txt", "work\n", "feat: add work")
	gittest.Run(t, fixture.root, "merge", "--ff-only", fixture.ref.Branch)
	disposition, err := ClassifySupersededBranch(context.Background(), fixture.run, nil)
	if err != nil {
		t.Fatalf("classify superseded branch: %v", err)
	}
	recordPath := filepath.Join(t.TempDir(), "branch-disposition.json")

	if err := RecordAndDiscardSupersededBranch(
		context.Background(),
		recordPath,
		disposition,
		time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC),
	); err != nil {
		t.Fatalf("record and discard superseded branch: %v", err)
	}

	record, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatalf("read disposition record after discard: %v", err)
	}
	var decoded branchDispositionRecord
	if err := json.Unmarshal(record, &decoded); err != nil {
		t.Fatalf("decode disposition record: %v\n%s", err, record)
	}
	if decoded.SchemaVersion != BranchDispositionRecordSchemaVersion ||
		decoded.RunID != fixture.run.ID || decoded.Branch != fixture.ref.Branch ||
		decoded.Reason != disposition.Reason ||
		!slices.Equal(decoded.Commits, disposition.Commits) ||
		!slices.Equal(decoded.ChangedFiles, []string{"work.txt"}) {
		t.Fatalf("decoded disposition record = %+v", decoded)
	}
	assertSupersededBranchSurface(t, fixture, false)
}

type supersededBranchFixture struct {
	root string
	ref  runworktree.Ref
	run  store.Run
}

func newSupersededBranchFixture(t *testing.T) supersededBranchFixture {
	t.Helper()
	ctx := context.Background()
	root := t.TempDir()
	var err error
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("resolve Git root: %v", err)
	}
	gittest.InitRepo(t, root, "-b", "ma/target")
	gittest.PersistIdentity(t, root)
	mustWriteForTest(t, filepath.Join(root, "base.txt"), "base\n")
	gittest.Run(t, root, "add", "base.txt")
	gittest.Run(t, root, "commit", "-m", "test: seed target")
	base := strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: root,
		Location: t.TempDir(),
		RunID:    "superseded-test",
		HeadSHA:  base,
	})
	if err != nil {
		t.Fatalf("create Run Worktree: %v", err)
	}
	ref.Path, err = filepath.EvalSymlinks(ref.Path)
	if err != nil {
		t.Fatalf("resolve Run Worktree path: %v", err)
	}
	created := time.Date(2026, 8, 10, 10, 0, 0, 0, time.UTC)
	return supersededBranchFixture{
		root: root,
		ref:  ref,
		run: store.Run{
			ID:          "superseded-test",
			Kind:        store.KindImplement,
			State:       store.StateStopped,
			GitRoot:     root,
			LocalBranch: "ma/target",
			HeadSHA:     base,
			WorkDir:     ref.Path,
			SpecSlug:    "disposition-spec",
			CreatedAt:   created,
			UpdatedAt:   created,
		},
	}
}

func (fixture supersededBranchFixture) commitRunChange(t *testing.T, path, content, message string) string {
	t.Helper()
	mustWriteForTest(t, filepath.Join(fixture.ref.Path, path), content)
	gittest.Run(t, fixture.ref.Path, "add", path)
	gittest.Run(t, fixture.ref.Path, "commit", "-m", message)
	return strings.TrimSpace(gittest.Run(t, fixture.ref.Path, "rev-parse", "HEAD"))
}

func (fixture supersededBranchFixture) commitTargetChange(t *testing.T, path, content, message string) string {
	t.Helper()
	mustWriteForTest(t, filepath.Join(fixture.root, path), content)
	gittest.Run(t, fixture.root, "add", path)
	gittest.Run(t, fixture.root, "commit", "-m", message)
	return strings.TrimSpace(gittest.Run(t, fixture.root, "rev-parse", "HEAD"))
}

func assertSupersededBranchSurface(t *testing.T, fixture supersededBranchFixture, exists bool) {
	t.Helper()
	_, worktreeErr := os.Stat(fixture.ref.Path)
	worktreeExists := worktreeErr == nil
	branchExists := strings.TrimSpace(gittest.Run(t, fixture.root, "branch", "--list", fixture.ref.Branch)) != ""
	if worktreeExists != exists || branchExists != exists {
		t.Fatalf("surface exists: worktree=%t branch=%t, want both %t", worktreeExists, branchExists, exists)
	}
}
