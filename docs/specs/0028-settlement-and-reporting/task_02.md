---
task: task_02
spec: 0028-settlement-and-reporting
status: pending
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

- [ ] Implement the idempotent store reclamation operation with its Run Event and reason
- [ ] Add an orphan-aware wrapper at the lock consumption seam and apply it to all five consumption points
- [ ] stderr warning wording shared across surfaces
- [ ] Integration tests (buffer-captured CLI runs): a dead-pid lock lets implement, settle, and a review command proceed with the warning; a live-pid lock still blocks with today's message; a pid-less legacy lock still blocks
- [ ] Store tests: reclamation is idempotent and journals exactly once

## Acceptance Criteria

- [ ] A store seeded with an Active Run owned by a reaped child pid no longer blocks a new implement Run; the reclaimed Run reads Failed with the reason recorded
- [ ] A lock owned by the live test process blocks with the unchanged error text
- [ ] A legacy pid-less lock blocks and is never reclaimed
- [ ] Two concurrent reclamation attempts settle to one terminal Run and one journal record
- [ ] The full test suite passes

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
