---
task: task_03
spec: 0014-run-store-retention
status: pending
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

- [ ] Wire `PruneTerminalRuns` + artifact cleanup into the preflight sweep
- [ ] Gate on `journal_retention > 0`
- [ ] Best-effort semantics (warn, never fatal) + summary line
- [ ] Tests: eligible Run pruned during a new Run's preflight; prune failure is non-fatal

## Acceptance Criteria

- [ ] Starting an operational Run with a non-zero window prunes an eligible terminal Run during preflight, without a manual step.
- [ ] A prune failure during the sweep produces a warning and the Run proceeds normally.
- [ ] With `journal_retention: 0`, the sweep prunes nothing.
- [ ] The sweep prunes via the same primitive as `roundfix gc` (no separate implementation).

## Verification

- `rtk go test ./internal/cli/` — expected: preflight-sweep prune tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 5; Core Feature 4. `_techspec.md` → Self-healing prune,
Build Order 3. ADR-0033. Ties into R4-1/R4-2 sweep work (same preflight seam).
