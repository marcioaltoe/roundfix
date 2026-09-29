// Suite: reconcile discovery of legacy terminal Runs.
// Invariant: reconcile from the main checkout reports a Run keyed from its surviving Run Worktree.
// Boundary IN: public reconcile runner, write-mode store backfill, and real Git worktree metadata.
// Boundary OUT: implement-preflight pruning and migration-only store cases.
package cli

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func TestReconcileReportsALegacyRunKeyedFromItsRunWorktree(t *testing.T) {
	ctx := context.Background()
	homeDir, repoDir, location := newReconcileWorkspace(t)
	repoDir, err := filepath.EvalSymlinks(repoDir)
	if err != nil {
		t.Fatalf("resolve main checkout: %v", err)
	}
	sourceRoot := filepath.Join(t.TempDir(), "removed-source-checkout")
	gitImplement(t, repoDir, "worktree", "add", "-b", "feature/legacy-source", sourceRoot)
	sourceRoot, err = filepath.EvalSymlinks(sourceRoot)
	if err != nil {
		t.Fatalf("resolve source checkout: %v", err)
	}

	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run Database for legacy Run: %v", err)
	}
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     sourceRoot,
		LocalBranch: "main",
		HeadSHA:     strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD")),
		SpecSlug:    "reconcile-spec",
		Agent:       "codex",
	})
	if err != nil {
		_ = runStore.Close()
		t.Fatalf("create legacy reconciliation Run: %v", err)
	}
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: sourceRoot,
		Location: location,
		RunID:    run.ID,
		HeadSHA:  run.HeadSHA,
	})
	if err != nil {
		_ = runStore.Close()
		t.Fatalf("create surviving Run Worktree: %v", err)
	}
	run, err = runStore.SetRunWorkDir(ctx, run.ID, ref.Path)
	if err != nil {
		_ = runStore.Close()
		t.Fatalf("record surviving Run Worktree: %v", err)
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, store.StateStopped)
	if err != nil {
		_ = runStore.Close()
		t.Fatalf("complete legacy reconciliation Run: %v", err)
	}
	run = completed.Run
	if err := runStore.Close(); err != nil {
		t.Fatalf("close legacy Run Database: %v", err)
	}
	database, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir))
	if err != nil {
		t.Fatalf("open Run Database to seed legacy key: %v", err)
	}
	if _, err := database.ExecContext(ctx,
		`UPDATE runs SET repository_root = '' WHERE id = ?`, run.ID,
	); err != nil {
		_ = database.Close()
		t.Fatalf("clear legacy Run repository key: %v", err)
	}
	if err := database.Close(); err != nil {
		t.Fatalf("close Run Database after seeding legacy key: %v", err)
	}
	gitImplement(t, repoDir, "worktree", "remove", sourceRoot)

	writer, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run Database for repository-key backfill: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close Run Database after repository-key backfill: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, ctx, []string{"reconcile", "--format=json"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("legacy-key reconcile exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("legacy-key reconcile stderr = %q, want empty", stderr.String())
	}
	var report reconcileReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode legacy-key reconciliation report: %v\n%s", err, stdout.String())
	}
	if len(report.Results) != 1 || report.Results[0].RunID != run.ID {
		t.Fatalf("legacy-key reconciliation results = %+v, want Run %q", report.Results, run.ID)
	}
}
