---
task: task_03
spec: 0014-run-store-retention
status: completed
type: backend
complexity: low
---

# Task 03: Best-effort retention prune in the preflight sweep

## Overview

Make storage self-heal: the operational preflight sweep (which already reaps
worktree debris and closes terminal sessions) gains a best-effort retention
prune, so a developer who never runs `roundfix gc` still gets a bounded Run
Database. Failures are warnings, never fatal.

## Requirements

1. MUST call the store retention prune (and pruned-Run artifact cleanup) from
   the implement/resolve/watch preflight sweep when `journal_retention > 0`.
2. MUST be best-effort: a prune failure is one stderr warning and never blocks
   or fails the Run.
3. MUST report one optional summary line when the sweep prunes something, in the
   sweep's existing report style.
4. MUST NOT prune when `journal_retention` is `0`, and MUST reuse the same
   terminal-only, Active-safe primitive as the GC Command (no duplicate logic).

## Subtasks

- [x] Wire `PruneTerminalRuns` + artifact cleanup into the preflight sweep
- [x] Gate on `journal_retention > 0`
- [x] Best-effort semantics (warn, never fatal) + summary line
- [x] Tests: eligible Run pruned during a new Run's preflight; prune failure is non-fatal

## Acceptance Criteria

- [x] Starting an operational Run with a non-zero window prunes an eligible terminal Run during preflight, without a manual step.
- [x] A prune failure during the sweep produces a warning and the Run proceeds normally.
- [x] With `journal_retention: 0`, the sweep prunes nothing.
- [x] The sweep prunes via the same primitive as `roundfix gc` (no separate implementation).

## Verification

- `rtk go test ./internal/cli/` — expected: preflight-sweep prune tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 5; Core Feature 4. `_techspec.md` → Self-healing prune,
Build Order 3. ADR-0033. Ties into R4-1/R4-2 sweep work (same preflight seam).

## Result

Implemented the best-effort Journal Retention sweep for `implement`, `resolve`,
and `watch` Run startup. The sweep is gated by `journal_retention > 0`, calls
the shared `pruneRunRetention` helper used by `roundfix gc`, and writes either
one summary line when storage is pruned or one warning line when pruning fails.
Warnings do not block the Run.

Acceptance evidence:

- Non-zero retention: `TestRunImplementPreflightPrunesRetainedRunStorage`
  starts an implement Run with `journal_retention: 336h`; preflight prints
  `roundfix: pruned Run storage runs=1 journal_rows=2 artifact_bytes=12`,
  removes the eligible terminal Run's journal and run artifact directory, and
  keeps Active and recent Runs.
- Best-effort failure: `TestRunImplementPreflightRetentionPruneFailureIsNonFatal`
  seeds an invalid run artifact path; preflight prints
  `roundfix: warning: Journal Retention prune failed:` and the implement Run
  still exits successfully with `Clean: all 1 Task(s) completed.`
- Zero retention: `TestRunImplementPreflightRetentionZeroSkipsPrune` sets
  `journal_retention: 0`; preflight prints no retention summary and leaves the
  old terminal Run journal and artifact directory in place.
- Shared primitive: `roundfix gc` and the preflight sweep both call
  `pruneRunRetention`, which delegates terminal selection and deletion to
  `Store.PruneTerminalRuns`.

Verification evidence:

- `rtk go test ./internal/cli/` passed: `Go test: 290 passed in 1 packages`.
- `rtk go test ./...` passed: `Go test: 801 passed in 18 packages`.
- Repo gate `rtk make verify` passed: full Go suite, `roundfix skills check`,
  and build completed.
