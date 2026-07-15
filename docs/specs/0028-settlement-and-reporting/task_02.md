---
task: task_02
spec: 0028-settlement-and-reporting
status: completed
type: backend
complexity: high
---

# Task 02: Reclaim orphaned Active-Run locks on proven owner death

## Overview

Close the blocked-relaunch failure mode: any surface that hits an Active-Run lock whose owning process is provably dead reclaims it automatically — the Run completes Failed with a recorded reason, the reclamation is journaled, one stderr warning names the dead run id and pid, and the blocked command proceeds. A live owner keeps today's block byte-for-byte.

## Requirements

1. MUST provide a store-level reclamation operation that completes a provably-orphaned Active Run as Failed with a reason shaped like `owner process <pid> not running; lock reclaimed`, writes a Run Event, and releases the lock; it MUST be idempotent when a concurrent caller already completed the Run.
2. MUST wire the orphan check into every Active-Run lock consumption point: the Implement Command preflight, the Settle Command preflight, the review preflights including the Branch Integrity Preflight's active-run guardrail, Stop Command target resolution, and the Run-creation conflict path.
3. MUST proceed with the blocked command after reclamation, emitting one stderr warning naming the reclaimed run id and pid.
4. MUST keep the existing Active-Run error unchanged for a live owner, and never reclaim when the Run row has no recorded pid (legacy rows) or liveness reports no proof.
5. MUST leave `stop --force` behavior unchanged as the manual path for everything short of proof.

## Subtasks

- [x] Implement the idempotent store reclamation operation with its Run Event and reason
- [x] Add an orphan-aware wrapper at the lock consumption seam and apply it to all five consumption points
- [x] stderr warning wording shared across surfaces
- [x] Integration tests (buffer-captured CLI runs): a dead-pid lock lets implement, settle, and a review command proceed with the warning; a live-pid lock still blocks with today's message; a pid-less legacy lock still blocks
- [x] Store tests: reclamation is idempotent and journals exactly once

## Acceptance Criteria

- [x] A store seeded with an Active Run owned by a reaped child pid no longer blocks a new implement Run; the reclaimed Run reads Failed with the reason recorded
- [x] A lock owned by the live test process blocks with the unchanged error text
- [x] A legacy pid-less lock blocks and is never reclaimed
- [x] Two concurrent reclamation attempts settle to one terminal Run and one journal record
- [x] The full test suite passes

## Context

- interface: `internal/store/store.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/settle.go`
- interface: `internal/cli/cli.go`

## Verification

- `grep -q "ReclaimOrphanedRun" internal/store/store.go` — expected: exit 0
- `go test ./internal/store/... ./internal/cli/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 1, User Story 1, Core Feature 1, Decisions; `_techspec.md` → Build Order 2, Interfaces (ReclaimOrphanedRun), API Contracts (orphaned lock), Risks (reclamation races); ADR-0044.

## Result

Implemented orphaned Active-Run lock reclamation:

- Added `Store.ReclaimOrphanedRun`, which marks a provably orphaned non-terminal Run `Failed`, releases its Active-Run lock, and appends one `daemon.outcome` Run Event with `state`, `reason`, and `owner_pid`.
- Added shared CLI orphan handling that reclaims only when `owner_pid` exists and `ProcessAlive` proves death; live owners, pid-less legacy rows, and no-proof platforms keep the existing block path.
- Wired the check into Implement Command preflight, Settle Command preflight, Branch Integrity Preflight, Stop Command target resolution, and Run-creation conflict retry paths.

Evidence:

- Pre-change signal: `rtk grep "ReclaimOrphanedRun" internal/store/store.go` exited 1 before implementation.
- Dead-owner implement path: `TestRunImplementReclaimsDeadOwnerActiveRun` passed and asserted the blocked implement Run became `Failed` with the reclamation reason event.
- Dead-owner settle path: `TestRunSettleReclaimsDeadOwnerActiveRun` passed and asserted settle proceeded after the shared stderr warning.
- Dead-owner review path: `TestReviewFetchReclaimsDeadOwnerActiveRun` passed and asserted fetch proceeded after Branch Integrity reclaimed the lock.
- Live-owner block: `TestRunImplementPreflightRejectsActiveRunInWorkingTree` passed with the current process PID and the unchanged block text.
- Legacy pid-less block: `TestRunSettleActiveRunOnSameWorkingTreeBlocks` passed and asserted no Run Event was written and the Run stayed `Active`.
- Concurrent idempotence: `TestReclaimOrphanedRunConcurrentAttemptsJournalOnce` passed and asserted one terminal `Failed` Run plus one journal record.
- Focused gate: `rtk go test ./internal/store ./internal/cli -run 'Test(ReclaimOrphanedRun.*|RunImplementReclaimsDeadOwnerActiveRun|RunSettleReclaimsDeadOwnerActiveRun|ReviewFetchReclaimsDeadOwnerActiveRun|RunImplementPreflightRejectsActiveRunInWorkingTree|RunSettleActiveRunOnSameWorkingTreeBlocks)$'` exited 0.
- Task verification: `rtk grep -q "ReclaimOrphanedRun" internal/store/store.go && rtk go test ./internal/store/... ./internal/cli/... && rtk go build -buildvcs=false ./...` exited 0.
- Full gate: `rtk make verify` exited 0 (`go test ./...`: 1202 passed in 19 packages; `roundfix skills check`: passed; build: passed).
