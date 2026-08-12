---
task: task_01
spec: 0014-run-store-retention
status: completed
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

- [x] `store.journal_retention` config key (default 336h, 0 = keep everything)
- [x] `PruneTerminalRuns(ctx, cutoff)` deleting only eligible `run_events`
- [x] Guard: Active Runs and `runs`/locks never touched
- [x] Table tests: terminal-past-cutoff pruned; Active + recent kept; counts

## Acceptance Criteria

- [x] `PruneTerminalRuns` deletes `run_events` only for terminal Runs with `completed_at` older than the cutoff and returns their ids and row count.
- [x] An Active Run (even with an old `created_at`) keeps its journal; all `runs` rows and locks survive any prune.
- [x] `journal_retention: 0` yields a cutoff that prunes nothing.
- [x] Config loads `store.journal_retention` with correct precedence; an unknown `store` key still fails strict validation.

## Verification

- `rtk go test ./internal/store/ ./internal/config/` — expected: prune and config tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 4; Core Features 1, 3. `_techspec.md` → Prune scope,
Build Order 1, Interfaces: `PruneTerminalRuns`. ADR-0033. Work-plan finding R4-5.

## Result

Implemented the retention primitive. Config now has `store.journal_retention`
with builtin default `336h`, User > Project override behavior, `0` accepted as
disabled retention, negative values rejected, and strict validation for unknown
`store` keys. The default generated config includes the new `store` section.

Added `store.PruneTerminalRuns(ctx, cutoff)` returning `PruneResult{RunIDs,
Events}`. The operation selects only terminal Runs with non-empty
`completed_at` before the cutoff, parses stored timestamps with the store's
existing parser, deletes only `run_events` rows, and leaves `runs` rows and
`active_run_locks` untouched.

Acceptance evidence:

- `TestPruneTerminalRunsDeletesOnlyEligibleJournalRows` prunes two old terminal
  Runs, returns their ids, deletes 3 journal rows, keeps recent terminal,
  Active, non-terminal, and empty-`completed_at` journals, verifies a zero
  cutoff prunes nothing, and checks all Run rows and locks survive.
- `TestPruneTerminalRunsNoOpsWhenCutoffSelectsNothing` verifies an empty prune
  result and retained journal when no terminal Run is older than the cutoff.
- `TestLoadAppliesStoreRetentionConfigHierarchy` covers builtin default, User
  override, Project override, and `journal_retention: 0`.
- `TestLoadRejectsUnknownStoreConfigKey` verifies strict validation still fails
  for an unknown `store` key.

Verification:

- `rtk go test ./internal/store/ ./internal/config/` — passed, 92 tests in 2
  packages.
- `rtk go test ./...` — passed, 794 tests in 18 packages.
- `rtk make verify` — passed: full Go suite, Roundfix skill check, and
  `go build -buildvcs=false -o bin/roundfix ./cmd/roundfix`.

Follow-up: GC Command artifact cleanup and the preflight sweep prune remain in
`task_02` and `task_03`.
