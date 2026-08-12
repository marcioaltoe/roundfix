---
task: task_08
spec: 0046-public-context-driven-baseline-command
status: completed
type: backend
complexity: high
---

# Task 08: Recover interrupted file transactions

## Overview

Build the repository-local transaction boundary required before Baseline can
replace multiple files safely. The slice proves staging, preimage comparison,
postimage verification, rollback, and restart recovery independently of the
later apply command.

## Requirements

1. MUST acquire an exclusive worktree transaction lock and keep recovery state
   in a Git-private path rather than the visible worktree.
2. MUST journal exact preimages, stage every postimage, fsync staged files,
   revalidate the complete bounded preimage, and replace paths in deterministic
   order.
3. MUST verify every postimage and roll back changed paths in reverse order
   after any incomplete mutation.
4. MUST recover or conclusively report an interrupted transaction before a new
   transaction begins.
5. MUST preserve required file modes, refuse unsafe path changes, and surface
   incomplete rollback as a blocking execution failure.

## Subtasks

- [x] Implement transaction locking and phase journaling.
- [x] Implement staged writes and complete preimage comparison.
- [x] Implement deterministic replacement and postimage verification.
- [x] Implement reverse rollback and interrupted-run recovery.
- [x] Add failure injection around every transaction phase.

## Acceptance Criteria

- [x] No destination changes before every staged file and preimage is valid.
- [x] Failure at each replacement or verification phase restores exact original bytes and modes.
- [x] A simulated process interruption is recoverable before the next apply.
- [x] Concurrent transaction attempts cannot mutate the same worktree.
- [x] Unsafe path, symlink, mode, and stale-preimage changes fail before replacement.
- [x] Incomplete rollback remains visible and cannot be reported as success.

## Context

- instruction: `docs/adr/0073-baseline-apply-uses-a-recoverable-multi-file-transaction.md`
- interface: `internal/config/profiles.go`
- interface: `internal/cli/profiles_configure.go`

## Verification

- `rtk go test -count=1 ./internal/baseline -run 'TestTransactionStagesBeforeMutation|TestTransactionRollback|TestTransactionRecovery|TestTransactionLock|TestTransactionRejectsStalePreimage'` — expected: staging, locking, failure injection, rollback, and recovery cases pass.
- `rtk go test -count=1 ./internal/baseline -run TestTransactionFailureMatrix` — expected: every injected phase failure preserves or restores the exact repository preimage.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Stories 3 and 7; Core Features 7, 13–14, 19; Success Metrics.
- `_techspec.md` → Interfaces: Transaction; Data Models: apply transaction; Testing Approach; Build Order 5.
- ADR-0073 → recoverable multi-file transaction.

## Result

Implemented the repository-local Baseline transaction boundary. Each
transaction now holds a worktree-scoped advisory lock, stores its exact
preimages and phase journal under the worktree's Git-private directory, stages
and fsyncs every regular postimage, revalidates the complete bounded preimage,
and replaces paths in lexical order. Destination writes use fsynced
same-directory temporary files, preserve the approved mode, verify every
postimage, and journal mutation intent before each filesystem change.

Rollback restores changed paths in reverse order from the exact journaled
bytes and modes, removes transaction-created directories, and keeps an
incomplete journal when recovery cannot finish. A later transaction recovers
an interrupted journal while holding the same exclusive lock before it
validates or creates new state. Unix and Windows lock adapters use the existing
`golang.org/x/sys` platform APIs.

Acceptance evidence:

- `TestTransactionStagesBeforeMutation` fails after complete staging and
  confirms the visible tree is byte-and-mode identical to its preimage.
- `TestTransactionFailureMatrix` injects 58 failures across journaling,
  per-path staging, complete preimage validation, every pre/post replacement
  point, every postimage verification, commit, and rollback. Every recoverable
  failure preserves or restores the exact visible tree; the rollback-blocked
  case returns `IncompleteRollbackError` and leaves recovery state for the
  next transaction.
- `TestTransactionRecovery` abandons a transaction after its first
  replacement and proves the next `BeginTransaction` restores the exact tree
  before returning.
- `TestTransactionLock` proves a second worktree transaction receives
  `ErrTransactionLocked` and cannot mutate.
- `TestTransactionRejectsStalePreimage` proves traversal paths, unsafe modes,
  symlink parents, stale bytes, stale modes, and a change made after the
  complete-set validation all fail before that destination is replaced.

Verification:

- `rtk go test -count=1 ./internal/baseline -run 'TestTransactionStagesBeforeMutation|TestTransactionRollback|TestTransactionRecovery|TestTransactionLock|TestTransactionRejectsStalePreimage'`
  — passed, 11 tests.
- `rtk go test -count=1 ./internal/baseline -run TestTransactionFailureMatrix`
  — passed, 58 tests.
- `rtk make verify` — passed: both 256-test setup-context-driven suites,
  1,891 Go tests in 21 packages, setup asset validation, Roundfix skill check,
  and binary build.

Follow-ups: none.
