---
task: task_08
spec: 0046-public-context-driven-baseline-command
status: pending
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

- [ ] Implement transaction locking and phase journaling.
- [ ] Implement staged writes and complete preimage comparison.
- [ ] Implement deterministic replacement and postimage verification.
- [ ] Implement reverse rollback and interrupted-run recovery.
- [ ] Add failure injection around every transaction phase.

## Acceptance Criteria

- [ ] No destination changes before every staged file and preimage is valid.
- [ ] Failure at each replacement or verification phase restores exact original bytes and modes.
- [ ] A simulated process interruption is recoverable before the next apply.
- [ ] Concurrent transaction attempts cannot mutate the same worktree.
- [ ] Unsafe path, symlink, mode, and stale-preimage changes fail before replacement.
- [ ] Incomplete rollback remains visible and cannot be reported as success.

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
