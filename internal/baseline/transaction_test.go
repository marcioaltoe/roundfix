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
	"path"
	"path/filepath"
	"reflect"
	"roundfix/internal/gittest"
	"sort"
	"strings"
	"testing"
)

func TestTransactionStagesBeforeMutation(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
			// Each case builds its own repository and shares no mutable
			// state, so the 57 fault points overlap instead of queueing.
			t.Parallel()
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

func TestHistoryMoveApply(t *testing.T) {
	t.Parallel()

	t.Run("moves every recorded identity", func(t *testing.T) {
		t.Parallel()
		repo, plan := newHistoryMoveTransactionRepository(t, map[string]string{
			"_archived/specs/0001-alpha/_prd.md": "alpha prd\n",
			"docs/findings/_archived/beta.md":    "beta finding\n",
		})

		tx := beginTestTransaction(t, repo, plan)
		evidence, err := tx.Apply(context.Background())
		if err != nil {
			t.Fatalf("Apply(): %v", err)
		}
		if err := tx.Close(); err != nil {
			t.Fatalf("Close(): %v", err)
		}
		if !reflect.DeepEqual(evidence.VerifiedHistoryMoves, plan.HistoryMoves) ||
			len(evidence.RefusedHistoryMoves) != 0 {
			t.Fatalf("history move evidence = %+v, want every move verified", evidence)
		}
		assertHistoryMoveState(t, repo, plan.HistoryMoves[0], "alpha prd\n", true)
		assertHistoryMoveState(t, repo, plan.HistoryMoves[1], "beta finding\n", true)
	})

	t.Run("empty ledger preserves the existing transaction result", func(t *testing.T) {
		t.Parallel()
		repo, plan := newTransactionRepository(t)
		tx := beginTestTransaction(t, repo, plan)
		evidence, err := tx.Apply(context.Background())
		if err != nil {
			t.Fatalf("Apply(): %v", err)
		}
		if err := tx.Close(); err != nil {
			t.Fatalf("Close(): %v", err)
		}
		if evidence.VerifiedHistoryMoves != nil || evidence.RefusedHistoryMoves != nil {
			t.Fatalf("empty history move evidence = %+v, want omitted ledgers", evidence)
		}
		if !verifiedPostimagesMatch(plan.Postimages, evidence.VerifiedPostimages) {
			t.Fatalf("verified postimages = %+v, want the approved postimages", evidence.VerifiedPostimages)
		}
	})
}

func TestHistoryMoveRemovesEmptiedSource(t *testing.T) {
	t.Parallel()

	t.Run("finished Review Artifact leaves no source shell or rerun finding", func(t *testing.T) {
		t.Parallel()

		repo := newPlanRepository(t)
		head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
		const reviewPath = "docs/specs/reviews/pr-201"
		reviewDir := filepath.Join(repo, filepath.FromSlash(reviewPath))
		historyPersistRound(t, reviewDir, "feature/merged", head)
		historyWriteFiles(t, repo, map[string]string{
			path.Join(reviewPath, "issues/001.md"): "finished issue\n",
		})

		plan := buildTestPlan(t, repo)
		if _, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest); err != nil {
			t.Fatalf("ApplyPlan(): %v", err)
		}
		if _, err := os.Lstat(reviewDir); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("finished Review Artifact directory error = %v, want not exist", err)
		}
		if info, err := os.Stat(filepath.Dir(reviewDir)); err != nil || !info.IsDir() {
			t.Fatalf("live Review Artifact root = (%v, %v), want surviving directory", info, err)
		}

		moves, findings, err := planHistoryMoves(context.Background(), repo)
		if err != nil {
			t.Fatalf("planHistoryMoves() rerun: %v", err)
		}
		if len(moves) != 0 || len(findings) != 0 {
			t.Fatalf("planHistoryMoves() rerun = (%#v, %#v), want no relocated Review Artifact", moves, findings)
		}
	})

	t.Run("unmoved file keeps its source directory", func(t *testing.T) {
		t.Parallel()

		const source = "docs/specs/reviews/pr-202/round-001/round.md"
		const unmoved = "docs/specs/reviews/pr-202/round-001/keep.txt"
		repo, plan := newHistoryMoveTransactionRepository(t, map[string]string{
			source: "round metadata\n",
		})
		writeTransactionFile(t, repo, unmoved, "keep me\n", 0o644)

		tx := beginTestTransaction(t, repo, plan)
		if _, err := tx.Apply(context.Background()); err != nil {
			t.Fatalf("Apply(): %v", err)
		}
		if err := tx.Close(); err != nil {
			t.Fatalf("Close(): %v", err)
		}
		content, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(unmoved)))
		if err != nil || string(content) != "keep me\n" {
			t.Fatalf("unmoved source file = %q error=%v, want original bytes", content, err)
		}
		assertHistoryMoveState(t, repo, plan.HistoryMoves[0], "round metadata\n", true)
	})
}

func TestHistoryMoveCollision(t *testing.T) {
	t.Parallel()

	repo, plan := newHistoryMoveTransactionRepository(t, map[string]string{
		"_archived/specs/0001-alpha/_prd.md": "alpha prd\n",
		"_archived/specs/0002-beta/_prd.md":  "beta prd\n",
	})
	collision := plan.HistoryMoves[0]
	sibling := plan.HistoryMoves[1]
	writeTransactionFile(t, repo, collision.To, "occupied destination\n", 0o644)

	tx := beginTestTransaction(t, repo, plan)
	evidence, err := tx.Apply(context.Background())
	if err != nil {
		t.Fatalf("Apply(): %v", err)
	}
	if err := tx.Close(); err != nil {
		t.Fatalf("Close(): %v", err)
	}
	if !reflect.DeepEqual(evidence.VerifiedHistoryMoves, []HistoryMove{sibling}) ||
		len(evidence.RefusedHistoryMoves) != 1 {
		t.Fatalf("history move evidence = %+v, want one verified sibling and one refusal", evidence)
	}
	refusal := evidence.RefusedHistoryMoves[0]
	if refusal.From != collision.From || refusal.To != collision.To ||
		!strings.Contains(refusal.Reason, "already exists") {
		t.Fatalf("collision refusal = %+v, want both paths and occupied reason", refusal)
	}
	collisionSource, readErr := os.ReadFile(filepath.Join(repo, filepath.FromSlash(collision.From)))
	if readErr != nil || string(collisionSource) != "alpha prd\n" ||
		transactionContentIdentity(collisionSource) != collision.ContentIdentity {
		t.Fatalf("collision source = %q error=%v, want recorded source bytes", collisionSource, readErr)
	}
	if content, readErr := os.ReadFile(filepath.Join(repo, filepath.FromSlash(collision.To))); readErr != nil ||
		string(content) != "occupied destination\n" {
		t.Fatalf("collision destination = %q error=%v, want original bytes", content, readErr)
	}
	assertHistoryMoveState(t, repo, sibling, "beta prd\n", true)

	reportRepo, reportPlan := newHistoryMoveTransactionRepository(t, map[string]string{
		"_archived/specs/0003-gamma/_prd.md": "gamma prd\n",
		"_archived/specs/0004-delta/_prd.md": "delta prd\n",
	})
	writeTransactionFile(t, reportRepo, reportPlan.HistoryMoves[0].To, "occupied report destination\n", 0o644)
	_, err = ApplyPlan(context.Background(), reportRepo, reportPlan, reportPlan.PlanDigest)
	if err == nil || !strings.Contains(err.Error(), "not every history relocation was performed") ||
		!strings.Contains(err.Error(), reportPlan.HistoryMoves[0].From) ||
		!strings.Contains(err.Error(), reportPlan.HistoryMoves[0].To) {
		t.Fatalf("ApplyPlan() error = %v, want incomplete relocation report naming both paths", err)
	}
	assertHistoryMoveState(t, reportRepo, reportPlan.HistoryMoves[1], "delta prd\n", true)
}

func TestHistoryMoveRollback(t *testing.T) {
	t.Parallel()

	t.Run("failure after a partial relocation restores every source", func(t *testing.T) {
		t.Parallel()
		repo, plan := newHistoryMoveTransactionRepository(t, map[string]string{
			"_archived/specs/0001-alpha/_prd.md": "alpha prd\n",
			"_archived/specs/0002-beta/_prd.md":  "beta prd\n",
		})
		tx := beginTestTransaction(t, repo, plan)
		tx.phaseHook = failTransactionOnce(
			transactionPhaseReplacing,
			plan.HistoryMoves[1].To,
			errors.New("injected history move failure"),
		)
		if _, err := tx.Apply(context.Background()); err == nil ||
			!strings.Contains(err.Error(), "injected history move failure") {
			t.Fatalf("Apply() error = %v, want injected history move failure", err)
		}
		if err := tx.Close(); err != nil {
			t.Fatalf("Close(): %v", err)
		}
		assertHistoryMoveState(t, repo, plan.HistoryMoves[0], "alpha prd\n", false)
		assertHistoryMoveState(t, repo, plan.HistoryMoves[1], "beta prd\n", false)
	})

	t.Run("destination identity mismatch fails and restores the recorded bytes", func(t *testing.T) {
		t.Parallel()
		repo, plan := newHistoryMoveTransactionRepository(t, map[string]string{
			"_archived/specs/0001-alpha/_prd.md": "alpha prd\n",
		})
		move := plan.HistoryMoves[0]
		tx := beginTestTransaction(t, repo, plan)
		tx.phaseHook = func(point transactionFaultPoint) error {
			if point.Phase == transactionPhaseVerifying && point.Path == move.To {
				writeTransactionFile(t, repo, move.To, "tampered destination\n", 0o644)
			}
			return nil
		}
		if _, err := tx.Apply(context.Background()); err == nil ||
			!strings.Contains(err.Error(), "content identity") {
			t.Fatalf("Apply() error = %v, want destination identity failure", err)
		}
		if err := tx.Close(); err != nil {
			t.Fatalf("Close(): %v", err)
		}
		assertHistoryMoveState(t, repo, move, "alpha prd\n", false)
	})

	t.Run("failure after source pruning recreates source directories", func(t *testing.T) {
		t.Parallel()

		const source = "docs/specs/reviews/pr-203/round-001/round.md"
		repo, plan := newHistoryMoveTransactionRepository(t, map[string]string{
			source: "round metadata\n",
		})
		move := plan.HistoryMoves[0]
		tx := beginTestTransaction(t, repo, plan)
		tx.phaseHook = failTransactionOnce(
			transactionPhaseVerifying,
			move.To,
			errors.New("injected failure after source pruning"),
		)
		if _, err := tx.Apply(context.Background()); err == nil ||
			!strings.Contains(err.Error(), "injected failure after source pruning") {
			t.Fatalf("Apply() error = %v, want injected post-pruning failure", err)
		}
		if err := tx.Close(); err != nil {
			t.Fatalf("Close(): %v", err)
		}
		assertHistoryMoveState(t, repo, move, "round metadata\n", false)
		if info, err := os.Stat(filepath.Join(repo, filepath.FromSlash(path.Dir(source)))); err != nil || !info.IsDir() {
			t.Fatalf("restored source directory = (%v, %v), want directory", info, err)
		}
	})

	t.Run("interrupted relocation recovers from the journal", func(t *testing.T) {
		t.Parallel()
		repo, plan := newHistoryMoveTransactionRepository(t, map[string]string{
			"_archived/specs/0001-alpha/_prd.md": "alpha prd\n",
		})
		move := plan.HistoryMoves[0]
		interrupted := beginTestTransaction(t, repo, plan)
		if err := interrupted.revalidatePreimages(context.Background()); err != nil {
			t.Fatalf("revalidate interrupted transaction: %v", err)
		}
		if refusal, err := interrupted.relocateHistoryMove(context.Background(), 0); err != nil || refusal != nil {
			t.Fatalf("relocate interrupted history move = (%+v, %v)", refusal, err)
		}
		assertHistoryMoveState(t, repo, move, "alpha prd\n", true)
		abandonTestTransaction(t, interrupted)

		recovered := beginTestTransaction(t, repo, plan)
		if err := recovered.Close(); err != nil {
			t.Fatalf("close recovered transaction: %v", err)
		}
		assertHistoryMoveState(t, repo, move, "alpha prd\n", false)
	})

	t.Run("corrupted history move sidecar is rejected on rollback", func(t *testing.T) {
		t.Parallel()
		repo, plan := newHistoryMoveTransactionRepository(t, map[string]string{
			"_archived/specs/0001-alpha/_prd.md": "alpha prd\n",
		})
		move := plan.HistoryMoves[0]
		interrupted := beginTestTransaction(t, repo, plan)
		if err := interrupted.revalidatePreimages(context.Background()); err != nil {
			t.Fatalf("revalidate interrupted transaction: %v", err)
		}
		if refusal, err := interrupted.relocateHistoryMove(context.Background(), 0); err != nil || refusal != nil {
			t.Fatalf("relocate interrupted history move = (%+v, %v)", refusal, err)
		}
		if err := os.Remove(filepath.Join(repo, filepath.FromSlash(move.To))); err != nil {
			t.Fatalf("remove interrupted history move destination: %v", err)
		}
		sidecar := filepath.Join(interrupted.stateDir, historyMoveContentName(move.Ordinal))
		if err := os.WriteFile(sidecar, []byte("corrupted sidecar bytes\n"), 0o600); err != nil {
			t.Fatalf("corrupt history move sidecar: %v", err)
		}
		abandonTestTransaction(t, interrupted)

		recovered, err := BeginTransaction(context.Background(), repo, plan)
		if recovered != nil {
			t.Fatalf("BeginTransaction() returned a transaction for a corrupted sidecar")
		}
		if err == nil || !strings.Contains(err.Error(), "sidecar for ordinal 0 does not match the journal") {
			t.Fatalf("BeginTransaction() error = %v, want sidecar mismatch", err)
		}
		if _, lerr := os.Lstat(filepath.Join(repo, filepath.FromSlash(move.From))); !errors.Is(lerr, fs.ErrNotExist) {
			t.Fatalf("corrupted sidecar restored source=%q lerr=%v, want source absent", move.From, lerr)
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

func newHistoryMoveTransactionRepository(t *testing.T, sources map[string]string) (string, PlanDocument) {
	t.Helper()
	repo, plan := newTransactionRepository(t)
	fromPaths := make([]string, 0, len(sources))
	for from, content := range sources {
		writeTransactionFile(t, repo, from, content, 0o644)
		fromPaths = append(fromPaths, from)
	}
	sort.Strings(fromPaths)
	plan.HistoryMoves = make([]HistoryMove, len(fromPaths))
	for index, from := range fromPaths {
		content := sources[from]
		plan.HistoryMoves[index] = HistoryMove{
			Ordinal:         index,
			From:            from,
			To:              historyMoveTestDestination(from),
			ContentIdentity: transactionContentIdentity([]byte(content)),
		}
	}
	var err error
	plan.PlanDigest, err = computePlanDigest(plan)
	if err != nil {
		t.Fatalf("compute history move Plan Digest: %v", err)
	}
	return repo, plan
}

func historyMoveTestDestination(from string) string {
	base := strings.TrimPrefix(from, "_archived/")
	base = strings.Replace(base, "docs/findings/_archived/", "findings/", 1)
	return path.Join("docs/history", base)
}

func assertHistoryMoveState(t *testing.T, repo string, move HistoryMove, content string, moved bool) {
	t.Helper()
	from := filepath.Join(repo, filepath.FromSlash(move.From))
	to := filepath.Join(repo, filepath.FromSlash(move.To))
	wantPresent, wantAbsent := to, from
	if !moved {
		wantPresent, wantAbsent = from, to
	}
	got, err := os.ReadFile(wantPresent)
	if err != nil || string(got) != content || transactionContentIdentity(got) != move.ContentIdentity {
		t.Fatalf("history move present path %q = %q error=%v, want recorded identity", wantPresent, got, err)
	}
	if _, err := os.Lstat(wantAbsent); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("history move absent path %q error = %v, want not exist", wantAbsent, err)
	}
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
