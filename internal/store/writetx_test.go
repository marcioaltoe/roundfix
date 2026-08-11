package store

import (
	"context"
	"database/sql"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
	"time"

	"roundfix/internal/runevent"
)

// TestWriteTxIsTheOnlyWriterTransaction proves the writer-transaction
// discipline this task establishes: every writer path in the store package
// opens its transaction through the withWriteTx helper — never through an
// ad-hoc BeginTx — and the machine-wide advisory lock is released on success,
// error, and cancellation.
func TestWriteTxIsTheOnlyWriterTransaction(t *testing.T) {
	t.Parallel()

	// Source discipline: every BeginTx call in the package's non-test source
	// must live inside the withWriteTx helper. A broad pattern over this
	// package fails before the work (direct BeginTx call sites) and passes
	// only once every writer routes through the helper.
	fset := token.NewFileSet()
	parsedPackages, err := parser.ParseDir(fset, ".", func(info os.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parse store package source: %v", err)
	}
	var helperStart, helperEnd token.Pos
	var helperCount int
	var beginTxCalls int
	var beginTxOutsideHelper []string
	for _, pkg := range parsedPackages {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				switch node := node.(type) {
				case *ast.FuncDecl:
					if node.Name.Name == "withWriteTx" {
						helperStart, helperEnd = node.Pos(), node.End()
						helperCount++
					}
				case *ast.CallExpr:
					selector, ok := node.Fun.(*ast.SelectorExpr)
					if !ok || selector.Sel.Name != "BeginTx" {
						return true
					}
					beginTxCalls++
					if helperStart == 0 || node.Pos() < helperStart || node.End() > helperEnd {
						beginTxOutsideHelper = append(beginTxOutsideHelper, fset.Position(node.Pos()).String())
					}
				}
				return true
			})
		}
	}
	if helperCount != 1 {
		t.Fatalf("expected exactly one withWriteTx helper, got %d", helperCount)
	}
	if beginTxCalls == 0 {
		t.Fatalf("withWriteTx helper does not open a transaction")
	}
	if len(beginTxOutsideHelper) > 0 {
		t.Fatalf("BeginTx called outside withWriteTx at: %v", beginTxOutsideHelper)
	}

	t.Run("lock releases on success", func(t *testing.T) {
		ctx := context.Background()
		homeDir := t.TempDir()
		first := openTestStore(t, ctx, homeDir)
		defer closeStore(t, first)
		if err := first.withWriteTx(ctx, "discipline write", func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, `INSERT INTO interactive_defaults (key, value, updated_at) VALUES ('discipline-success', 'ok', '2026-01-01T00:00:00Z')`)
			return err
		}); err != nil {
			t.Fatalf("first writer on success: %v", err)
		}
		second := openTestStore(t, ctx, homeDir)
		defer closeStore(t, second)
		if err := second.withWriteTx(ctx, "discipline write", func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, `INSERT INTO interactive_defaults (key, value, updated_at) VALUES ('discipline-success-2', 'ok', '2026-01-01T00:00:00Z')`)
			return err
		}); err != nil {
			t.Fatalf("second writer after success must not block on a stale lock: %v", err)
		}
	})

	t.Run("lock releases on error", func(t *testing.T) {
		ctx := context.Background()
		homeDir := t.TempDir()
		first := openTestStore(t, ctx, homeDir)
		defer closeStore(t, first)
		err := first.withWriteTx(ctx, "discipline write", func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, `INSERT INTO missing_table (x) VALUES (1)`)
			return err
		})
		if err == nil {
			t.Fatalf("expected first writer to fail")
		}
		second := openTestStore(t, ctx, homeDir)
		defer closeStore(t, second)
		if err := second.withWriteTx(ctx, "discipline write", func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, `INSERT INTO interactive_defaults (key, value, updated_at) VALUES ('discipline-error', 'ok', '2026-01-01T00:00:00Z')`)
			return err
		}); err != nil {
			t.Fatalf("second writer after error must not block on a stale lock: %v", err)
		}
	})

	t.Run("lock releases on cancellation", func(t *testing.T) {
		ctx := context.Background()
		homeDir := t.TempDir()
		first := openTestStore(t, ctx, homeDir)
		defer closeStore(t, first)
		// Open the second writer before the first holds the lock: Open runs a
		// migration that is itself a write, so it must not contend.
		second := openTestStore(t, ctx, homeDir)
		defer closeStore(t, second)

		firstCtx, cancelFirst := context.WithCancel(ctx)
		defer cancelFirst()
		held := make(chan struct{})
		firstDone := make(chan error, 1)
		go func() {
			firstDone <- first.withWriteTx(firstCtx, "discipline write", func(tx *sql.Tx) error {
				close(held)
				<-firstCtx.Done()
				return firstCtx.Err()
			})
		}()
		<-held

		secondCtx, cancelSecond := context.WithTimeout(ctx, 200*time.Millisecond)
		defer cancelSecond()
		secondErr := second.withWriteTx(secondCtx, "discipline write", func(tx *sql.Tx) error {
			return nil
		})
		if !errors.Is(secondErr, context.DeadlineExceeded) {
			t.Fatalf("expected second writer to block on the held lock until cancellation, got %v", secondErr)
		}

		cancelFirst()
		if err := <-firstDone; !errors.Is(err, context.Canceled) {
			t.Fatalf("expected first writer to return after cancellation, got %v", err)
		}
		if err := second.withWriteTx(ctx, "discipline write", func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, `INSERT INTO interactive_defaults (key, value, updated_at) VALUES ('discipline-cancel', 'ok', '2026-01-01T00:00:00Z')`)
			return err
		}); err != nil {
			t.Fatalf("writer after cancellation must not block on a stale lock: %v", err)
		}
	})
}

// TestConcurrentWritersAllocateMonotonicContiguousCursors opens independent
// writer Stores against one Run Database and proves that the machine-wide
// advisory lock serializes them: concurrent appends to one Run allocate every
// cursor exactly once, contiguously, with no gaps or duplicates.
func TestConcurrentWritersAllocateMonotonicContiguousCursors(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()

	seed := openTestStore(t, ctx, homeDir)
	run, err := seed.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	closeStore(t, seed)

	const writers = 4
	const eventsPerWriter = 25
	stores := make([]*Store, writers)
	for i := range stores {
		stores[i] = openTestStore(t, ctx, homeDir)
		defer closeStore(t, stores[i])
	}

	results := make(chan []int64, writers)
	errs := make(chan error, writers)
	for _, store := range stores {
		store := store
		go func() {
			var cursors []int64
			for i := 0; i < eventsPerWriter; i++ {
				cursor, err := store.AppendRunEvents(ctx, []runevent.RunEvent{
					sampleRunEvent(run.ID, "concurrent append"),
				})
				if err != nil {
					errs <- err
					return
				}
				cursors = append(cursors, cursor...)
			}
			results <- cursors
		}()
	}

	seen := make(map[int64]bool)
	total := 0
	for i := 0; i < writers; i++ {
		select {
		case err := <-errs:
			t.Fatalf("concurrent append: %v", err)
		case cursors := <-results:
			for _, cursor := range cursors {
				if seen[cursor] {
					t.Fatalf("cursor %d allocated more than once", cursor)
				}
				seen[cursor] = true
				total++
			}
		}
	}
	want := writers * eventsPerWriter
	if total != want {
		t.Fatalf("expected %d cursors total, got %d", want, total)
	}
	// The first cursor is 1 (MAX+1 over an empty journal), so a contiguous
	// allocation spans [1, want].
	for cursor := 1; cursor <= want; cursor++ {
		if !seen[int64(cursor)] {
			t.Fatalf("cursor %d missing: cursors are not contiguous", cursor)
		}
	}
}
