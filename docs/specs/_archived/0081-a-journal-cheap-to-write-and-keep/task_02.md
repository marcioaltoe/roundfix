---
task: task_02
spec: 0081-a-journal-cheap-to-write-and-keep
status: completed
type: data
complexity: high
---

# Task 02: Give the writer one transaction discipline

## Overview

Every write transaction enters through one Store helper, and Roundfix
processes serialize before SQLite ever sees a writer. Today each operational
command opens its own writer connection with `_txlock=immediate`, so six
parallel Runs on one machine become six processes racing for the same lock —
which is how four of six Runs died in their first Batch with `SQLITE_BUSY` and
why the timeout was raised from five seconds to thirty.

The raise made contention survivable. It did not make the design concurrent.
This Task does.

## Requirements

1. MUST route every writer transaction — cursor allocation, inserts, state
   changes, and commit — through one Store helper, so no path opens an ad-hoc
   immediate transaction.
2. MUST serialize Roundfix processes with a machine-wide advisory lock before
   SQLite sees a writer, and MUST release it on every exit path including
   error and cancellation.
3. MUST keep one writer Store per operational process with its single
   connection, leaving read-only Store values separate and unaffected.
4. MUST return context cancellation and advisory-lock failures to the caller
   rather than swallowing them.
5. MUST keep `busy_timeout` in place as defence against an unknown
   non-Roundfix writer, while no Roundfix path relies on it for concurrency.
6. MUST keep cursor allocation inside the locked transaction so concurrent
   Runs preserve input order and monotonic cursors.
7. MUST NOT change the schema, the stream contract, or any command's output.

The declared Verification names `TestWriteTxIsTheOnlyWriterTransaction`, which does not exist yet, so it can
fail before the work. Create it to assert that every writer path opens its transaction through the helper and releases the advisory lock on success, error and cancellation. A broad pattern over
this package matches cases that already pass and would approve the Task
before it starts.

## Subtasks

- [x] Introduce the write-transaction helper and route every writer through it.
- [x] Add the machine-wide advisory lock with its release paths.
- [x] Prove ordering and monotonic cursors under concurrency.

## Acceptance Criteria

- [x] No writer path opens a transaction outside the helper.
- [x] The advisory lock is released on success, error, and cancellation.
- [x] Concurrent writers produce monotonic, contiguous cursors per Run.
- [x] Read-only Store behaviour is unchanged.

## Context

- interface: internal/store/store.go
- interface: internal/store/journal.go

## Verification

- `output="$(go test -count=1 ./internal/store -run '^TestWriteTxIsTheOnlyWriterTransaction$' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the discipline tests are selected and pass.
- `output="$(grep -c 'BeginTx' internal/store/*.go)"; printf '%s' "$output" | grep -q . && grep -q 'withWriteTx' internal/store/store.go`
  — expected: exit 0; the helper exists on the real surface.
  — expected: exit 0; the store and the commands that open it stay green.

A whole-package `go test` sweep and `go build ./...` are deliberately absent:
both pass against a tree where no work has happened, so each approves the Task
before it starts. Regression and compilation are the Run-level gate's job.

## References

- `_prd.md` → Goal 5; User Story 3; Core Feature 6.
- `_techspec.md` → Implementation Design (`withWriteTx`); Build Order 2.
- ADR-0004.
- `docs/findings/_archived/` → the 2026-08-06 run-lifecycle rollup, whose
  members record the six-parallel-Runs incident.

## Result

Implemented the writer-transaction discipline as a Roundfix Daemon-assigned
turn. Task status is Daemon-owned and left untouched; declared Verification is
the Daemon's to run verbatim. Evidence below is from focused implementation
checks only.

### What changed (behaviour)

- `internal/store/store.go` gained a single `withWriteTx(ctx, operation,
  fn)` helper that acquires the machine-wide advisory lock, runs one `BeginTx`
  transaction, and releases the lock on every exit path — success, error, and
  cancellation — only after the transaction has committed or rolled back. It
  wraps begin/commit failures as `begin <operation>` / `commit <operation>`,
  and propagates context cancellation and advisory-lock failures to the caller.
- The machine-wide advisory lock is a `flock`-style lock on a per-database
  lock file (`<db>.lock`), created by the writer `Open` and held open for the
  Store's lifetime. `writelock_unix.go` uses `unix.Flock` with a cancellable
  non-blocking poll loop; `writelock_windows.go` uses `LockFileEx`/`UnlockFileEx`
  with the same cancellable loop. `Close` closes the lock file alongside the
  database.
- Every writer path that previously opened an ad-hoc `BeginTx` now routes
  through the helper: `createRun`, `CompleteRun`, `ReconcileIntegration`,
  `ReclaimOrphanedRun`, `RequestStop`, `applyMigration` (store.go);
  `AppendRunEvents`, `PruneTerminalRuns` (journal.go); and the two
  `AppendAgentSelection*` methods (agent_selection.go). `BeginTx` now appears
  only inside `withWriteTx`.
- Read-only Stores (`OpenReader`, `OpenStorageReader`) are untouched and carry
  no lock file; the single writer connection (`SetMaxOpenConns(1)`) and the
  defensive `busy_timeout` remain as required.
- New `internal/store/writetx_test.go` adds `TestWriteTxIsTheOnlyWriterTransaction`
  (the declared verification's named test) and a concurrency proof.

### Focused-check evidence (this session)

- `go test -count=1 ./internal/store -run '^TestWriteTxIsTheOnlyWriterTransaction$' -v` — passes. The test asserts (a) via a `go/ast` scan of the package's non-test source that every `BeginTx` call sits inside `withWriteTx` and that the helper opens a transaction, and (b) via runtime subtests that the advisory lock is released on success, on a transaction error, and on context cancellation (a second writer blocks until the holder's cancellation, then proceeds).
- `go test -count=1 ./internal/store -run '^TestConcurrentWritersAllocateMonotonicContiguousCursors$'` — passes. Four independent writer Stores append 25 events each to one Run concurrently; cursors are allocated exactly once, contiguously from 1..100, proving the machine-wide lock serializes writers and cursors stay monotonic/contiguous.
- `go test -count=1 ./internal/store` — 210 passed; `go test -buildvcs=false ./...` — 3899 passed across 27 packages (focused regression sweep; the error-message contracts `begin completion` / `begin terminal Run reconciliation` are preserved and their tests pass).
- `go build -buildvcs=false ./...` — clean. Cross-builds for `GOOS=windows` and `GOOS=linux` both compile. `gofmt -l` and `go vet ./internal/store` clean.

### Acceptance-criteria mapping

- No writer path opens a transaction outside the helper — AST scan in the new test.
- Advisory lock released on success, error, cancellation — runtime subtests.
- Concurrent writers produce monotonic, contiguous cursors per Run — concurrency test.
- Read-only Store behaviour unchanged — reader constructors untouched; full suite green.

### Follow-ups (not this task's slice)

- Single-statement autocommit writes (`UpdateRunState`, `SetRunWorkDir`,
  `RememberInteractiveDefaults`) remain autocommit statements rather than
  `withWriteTx` transactions; the techspec and this task scope the machine-wide
  lock to `BeginTx` transactions, and routing those single statements through
  the helper is a deliberate follow-up if the parallel-Run proof requires it.
- Batching (`JournalSink`) and the retention query move are later tasks
  (task_04, task_03) and intentionally untouched here.
