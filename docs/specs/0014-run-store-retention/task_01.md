---
task: task_01
spec: 0014-run-store-retention
status: pending
type: backend
complexity: medium
---

# Task 01: Journal retention config and terminal-only prune in the store

## Overview

Add the retention primitive: a `store.journal_retention` config key and a store
operation that prunes the Run Event Journal of terminal Runs older than a
cutoff, returning what it removed. This is the safe foundation both the GC
Command and the preflight sweep build on — it must never touch Active Runs,
`runs` rows, or locks.

## Requirements

1. MUST add a `store` config section with `journal_retention` as a Go duration
   (default `336h`; `0` disables pruning), strict-decoded with the normal
   Project > User > builtin precedence.
2. MUST add a store operation `PruneTerminalRuns(ctx, cutoff)` that deletes
   `run_events` rows only for Runs in a terminal state with a non-empty
   `completed_at` older than the cutoff, and returns the pruned run ids and row
   count.
3. MUST NOT delete `runs` rows or `active_run_locks`, and MUST exclude Active
   Runs (non-terminal state or empty `completed_at`) from pruning.
4. MUST be a no-op when the cutoff would select nothing, and covered by tests
   seeding a mix of Active and terminal Runs at varied timestamps.

## Subtasks

- [ ] `store.journal_retention` config key (default 336h, 0 = keep everything)
- [ ] `PruneTerminalRuns(ctx, cutoff)` deleting only eligible `run_events`
- [ ] Guard: Active Runs and `runs`/locks never touched
- [ ] Table tests: terminal-past-cutoff pruned; Active + recent kept; counts

## Acceptance Criteria

- [ ] `PruneTerminalRuns` deletes `run_events` only for terminal Runs with `completed_at` older than the cutoff and returns their ids and row count.
- [ ] An Active Run (even with an old `created_at`) keeps its journal; all `runs` rows and locks survive any prune.
- [ ] `journal_retention: 0` yields a cutoff that prunes nothing.
- [ ] Config loads `store.journal_retention` with correct precedence; an unknown `store` key still fails strict validation.

## Verification

- `rtk go test ./internal/store/ ./internal/config/` — expected: prune and config tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 4; Core Features 1, 3. `_techspec.md` → Prune scope,
Build Order 1, Interfaces: `PruneTerminalRuns`. ADR-0033. Work-plan finding R4-5.
