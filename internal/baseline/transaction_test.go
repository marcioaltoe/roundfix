// Suite: Baseline recoverable file transactions
// Invariant: a transaction installs every exact postimage or restores the exact bounded preimage before another transaction can mutate.
// Boundary IN: real Git-private state, real filesystem modes and bytes, advisory locking, interruption recovery, and phase failures.
// Boundary OUT: Baseline CLI apply orchestration and post-apply Baseline verification.

package baseline

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestTransactionStagesBeforeMutation(t *testing.T) {
	repo, plan := newTransactionRepository(t)
	before := snapshotVisibleTree(t, repo)

	tx := beginTestTransaction(t, repo, plan)
	tx.phaseHook = func(point transactionFaultPoint) error {
		if point.Phase == transactionPhaseStaged {
			return errors.New("injected staged failure")
		}
		return nil
	}
	if _, err := tx.Apply(context.Background()); err == nil ||
		!strings.Contains(err.Error(), "injected staged failure") {
		t.Fatalf("Apply() error = %v, want injected staged failure", err)
	}
	if err := tx.Close(); err != nil {
		t.Fatalf("Close() after staged failure: %v", err)
	}
	assertVisibleTree(t, repo, before)
}

func TestTransactionRollback(t *testing.T) {
	repo, plan := newTransactionRepository(t)
	before := snapshotVisibleTree(t, repo)

	tx := beginTestTransaction(t, repo, plan)
	tx.phaseHook = failTransactionOnce(transactionPhaseReplaced, "", errors.New("injected replacement failure"))
	if _, err := tx.Apply(context.Background()); err == nil ||
		!strings.Contains(err.Error(), "injected replacement failure") {
		t.Fatalf("Apply() error = %v, want injected replacement failure", err)
	}
	if err := tx.Close(); err != nil {
		t.Fatalf("Close() after rollback: %v", err)
	}
	assertVisibleTree(t, repo, before)
}

func TestTransactionRecovery(t *testing.T) {
	repo, plan := newTransactionRepository(t)
	before := snapshotVisibleTree(t, repo)

	interrupted := beginTestTransaction(t, repo, plan)
	if err := interrupted.stagePostimages(context.Background()); err != nil {
		t.Fatalf("stage interrupted transaction: %v", err)
	}
	if err := interrupted.revalidatePreimages(context.Background()); err != nil {
		t.Fatalf("revalidate interrupted transaction: %v", err)
	}
	if err := interrupted.replacePostimage(context.Background(), 0); err != nil {
		t.Fatalf("replace interrupted transaction path: %v", err)
	}
	if reflect.DeepEqual(snapshotVisibleTree(t, repo), before) {
		t.Fatal("interruption fixture did not mutate the worktree")
	}
	abandonTestTransaction(t, interrupted)

	recovered := beginTestTransaction(t, repo, plan)
	assertVisibleTree(t, repo, before)
	if err := recovered.Close(); err != nil {
		t.Fatalf("close recovered transaction: %v", err)
	}
}

func TestTransactionLock(t *testing.T) {
	repo, plan := newTransactionRepository(t)
	first := beginTestTransaction(t, repo, plan)
	defer func() {
		if err := first.Close(); err != nil {
			t.Errorf("close first transaction: %v", err)
		}
	}()

	second, err := BeginTransaction(context.Background(), repo, plan)
	if second != nil || !errors.Is(err, ErrTransactionLocked) {
		t.Fatalf("second BeginTransaction() = (%v, %v), want ErrTransactionLocked", second, err)
	}
}

func TestTransactionRejectsStalePreimage(t *testing.T) {
	t.Run("unsafe path", func(t *testing.T) {
		repo, plan := newTransactionRepository(t)
		plan.Postimages[0].Path = "../escape"
		plan.PlanDigest, _ = computePlanDigest(plan)
		assertTransactionBeginFailsWithoutMutation(t, repo, plan, "unsafe")
	})

	t.Run("unsafe postimage mode", func(t *testing.T) {
		repo, plan := newTransactionRepository(t)
		plan.Postimages[0].Mode = uint32(fs.ModeSetuid | 0o644)
		plan.PlanDigest, _ = computePlanDigest(plan)
		assertTransactionBeginFailsWithoutMutation(t, repo, plan, "unsafe mode")
	})

	t.Run("bytes", func(t *testing.T) {
		repo, plan := newTransactionRepository(t)
		writeTransactionFile(t, repo, "AGENTS.md", "altered instructions\n", 0o600)
		assertTransactionBeginFailsWithoutMutation(t, repo, plan, "stale")
	})

	t.Run("mode", func(t *testing.T) {
		repo, plan := newTransactionRepository(t)
		if err := os.Chmod(filepath.Join(repo, "AGENTS.md"), 0o644); err != nil {
			t.Fatal(err)
		}
		assertTransactionBeginFailsWithoutMutation(t, repo, plan, "stale")
	})

	t.Run("symlink parent", func(t *testing.T) {
		repo, plan := newTransactionRepository(t)
		target := filepath.Join(repo, "unsafe-docs")
		if err := os.Mkdir(target, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("unsafe-docs", filepath.Join(repo, "docs")); err != nil {
			t.Fatal(err)
		}
		assertTransactionBeginFailsWithoutMutation(t, repo, plan, "symlink")
	})

	t.Run("after complete validation", func(t *testing.T) {
		repo, plan := newTransactionRepository(t)
		tx := beginTestTransaction(t, repo, plan)
		changed := false
		tx.phaseHook = func(point transactionFaultPoint) error {
			if point.Phase == transactionPhaseReplacing && point.Path == "AGENTS.md" && !changed {
				changed = true
				writeTransactionFile(t, repo, "AGENTS.md", "stale---instructions\n", 0o600)
			}
			return nil
		}
		if _, err := tx.Apply(context.Background()); err == nil ||
			!strings.Contains(strings.ToLower(err.Error()), "stale") {
			t.Fatalf("Apply() error = %v, want stale preimage refusal", err)
		}
		if err := tx.Close(); err != nil {
			t.Fatalf("close stale transaction: %v", err)
		}
		content, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != "stale---instructions\n" {
			t.Fatalf("stale user change was replaced or rolled back: %q", content)
		}
	})
}

func TestTransactionFailureMatrix(t *testing.T) {
	templateRepo, templatePlan := newTransactionRepository(t)
	points := []transactionFaultPoint{
		{Phase: transactionPhaseJournaled},
		{Phase: transactionPhaseStaged},
		{Phase: transactionPhasePreimagesValidated},
		{Phase: transactionPhaseCommitting},
	}
	for _, postimage := range templatePlan.Postimages {
		points = append(points,
			transactionFaultPoint{Phase: transactionPhaseStaging, Path: postimage.Path},
			transactionFaultPoint{Phase: transactionPhaseReplacing, Path: postimage.Path},
			transactionFaultPoint{Phase: transactionPhaseReplaced, Path: postimage.Path},
			transactionFaultPoint{Phase: transactionPhaseVerifying, Path: postimage.Path},
		)
	}
	_ = templateRepo

	for _, point := range points {
		point := point
		name := string(point.Phase)
		if point.Path != "" {
			name += "/" + strings.ReplaceAll(point.Path, "/", "_")
		}
		t.Run(name, func(t *testing.T) {
			repo, plan := newTransactionRepository(t)
			before := snapshotVisibleTree(t, repo)
			hook := failTransactionOnce(point.Phase, point.Path, fmt.Errorf("injected %s failure", point.Phase))
			if point.Phase == transactionPhaseJournaled {
				transaction, err := beginTransaction(context.Background(), repo, plan, hook)
				if transaction != nil || err == nil {
					t.Fatalf("beginTransaction() at %s = (%v, %v), want failure", point, transaction, err)
				}
				assertVisibleTree(t, repo, before)
				return
			}
			tx := beginTestTransaction(t, repo, plan)
			tx.phaseHook = hook
			if _, err := tx.Apply(context.Background()); err == nil {
				t.Fatalf("Apply() at %s succeeded, want failure", point)
			}
			if err := tx.Close(); err != nil {
				t.Fatalf("Close() at %s: %v", point, err)
			}
			assertVisibleTree(t, repo, before)
		})
	}

	t.Run("incomplete rollback remains blocking", func(t *testing.T) {
		repo, plan := newTransactionRepository(t)
		tx := beginTestTransaction(t, repo, plan)
		replaced := false
		tx.phaseHook = func(point transactionFaultPoint) error {
			if point.Phase == transactionPhaseReplaced && !replaced {
				replaced = true
				return errors.New("trigger rollback")
			}
			if point.Phase == transactionPhaseRollingBack {
				return errors.New("rollback blocked")
			}
			return nil
		}
		_, err := tx.Apply(context.Background())
		var incomplete *IncompleteRollbackError
		if !errors.As(err, &incomplete) {
			t.Fatalf("Apply() error = %v, want IncompleteRollbackError", err)
		}
		if err := tx.Close(); err != nil {
			t.Fatalf("release incomplete transaction lock: %v", err)
		}

		recovery, err := BeginTransaction(context.Background(), repo, plan)
		if err != nil {
			t.Fatalf("recover incomplete rollback: %v", err)
		}
		if err := recovery.Close(); err != nil {
			t.Fatalf("close recovered transaction: %v", err)
		}
	})
}

func newTransactionRepository(t *testing.T) (string, PlanDocument) {
	t.Helper()
	repo := newPlanRepository(t)
	writeTransactionFile(t, repo, "AGENTS.md", "original instructions\n", 0o600)
	commitInspectionRepository(t, repo, "seed transaction preimage")
	return repo, buildTestPlan(t, repo)
}

func beginTestTransaction(t *testing.T, repo string, plan PlanDocument) *fileTransaction {
	t.Helper()
	transaction, err := BeginTransaction(context.Background(), repo, plan)
	if err != nil {
		t.Fatalf("BeginTransaction(): %v", err)
	}
	tx, ok := transaction.(*fileTransaction)
	if !ok {
		t.Fatalf("BeginTransaction() type = %T, want *fileTransaction", transaction)
	}
	return tx
}

func abandonTestTransaction(t *testing.T, tx *fileTransaction) {
	t.Helper()
	if err := unlockTransactionFile(tx.lock); err != nil {
		t.Fatalf("unlock interrupted transaction: %v", err)
	}
	if err := tx.lock.Close(); err != nil {
		t.Fatalf("close interrupted transaction lock: %v", err)
	}
	tx.closed = true
}

func assertTransactionBeginFailsWithoutMutation(
	t *testing.T,
	repo string,
	plan PlanDocument,
	message string,
) {
	t.Helper()
	before := snapshotVisibleTree(t, repo)
	transaction, err := BeginTransaction(context.Background(), repo, plan)
	if transaction != nil || err == nil || !strings.Contains(strings.ToLower(err.Error()), message) {
		t.Fatalf("BeginTransaction() = (%v, %v), want error containing %q", transaction, err, message)
	}
	assertVisibleTree(t, repo, before)
}

func failTransactionOnce(
	phase transactionPhase,
	path string,
	injected error,
) func(transactionFaultPoint) error {
	failed := false
	return func(point transactionFaultPoint) error {
		if !failed && point.Phase == phase && (path == "" || point.Path == path) {
			failed = true
			return injected
		}
		return nil
	}
}

type visibleTreeEntry struct {
	Path       string
	Mode       fs.FileMode
	LinkTarget string
	Content    string
}

func snapshotVisibleTree(t *testing.T, root string) []visibleTreeEntry {
	t.Helper()
	entries := []visibleTreeEntry{}
	err := filepath.WalkDir(root, func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, name)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if relative == ".git" {
			return fs.SkipDir
		}
		if relative == "." {
			return nil
		}
		info, err := os.Lstat(name)
		if err != nil {
			return err
		}
		item := visibleTreeEntry{Path: relative, Mode: info.Mode()}
		switch {
		case info.Mode()&fs.ModeSymlink != 0:
			item.LinkTarget, err = os.Readlink(name)
		case info.Mode().IsRegular():
			var content []byte
			content, err = os.ReadFile(name)
			item.Content = string(content)
		}
		if err != nil {
			return err
		}
		entries = append(entries, item)
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot visible tree: %v", err)
	}
	sort.Slice(entries, func(left, right int) bool {
		return entries[left].Path < entries[right].Path
	})
	return entries
}

func assertVisibleTree(t *testing.T, repo string, want []visibleTreeEntry) {
	t.Helper()
	if got := snapshotVisibleTree(t, repo); !reflect.DeepEqual(got, want) {
		t.Fatalf("visible worktree changed:\ngot:  %+v\nwant: %+v", got, want)
	}
}

func writeTransactionFile(t *testing.T, root, relative, content string, mode fs.FileMode) {
	t.Helper()
	name := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(name, mode); err != nil {
		t.Fatal(err)
	}
}
