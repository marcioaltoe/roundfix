---
task: task_02
spec: 0144-a-run-stops-when-its-budget-is-spent
status: pending
type: backend
complexity: medium
---

# Task 02: Settle the bounded Run as BudgetExceeded

## Overview

A Run ended on its budget must say so. This slice settles it with the outcome
the Run Database already carries for that cause, and a reason a reader can act
on.

## Requirements

1. MUST settle a Run ended by its budget as `BudgetExceeded`.
2. MUST record a reason naming the configured maximum and the elapsed time,
   within the existing public reason length.
3. MUST keep completed Tasks completed, and settle an interrupted Task the way an
   interrupted Task settles today.
4. MUST preserve the Run Worktree and the Run Branch, as every non-integrated
   outcome does.
5. MUST leave `Stopped` and `Unresolved` meaning exactly what they mean today.
6. MUST name the reason builder `budgetExceededReason`, so the Task's own check
   and later readers agree.

## Subtasks

- [ ] Settle the bounded Run with the existing outcome.
- [ ] Build the reason from the maximum and the elapsed time.
- [ ] Cover the outcome, the reason's two facts and the preserved worktree.

## Acceptance Criteria

- [ ] A bounded fixture Run settles `BudgetExceeded`.
- [ ] Its reason names the configured maximum and the elapsed time.
- [ ] Tasks that completed before the bound keep their status.
- [ ] The Run Worktree and Run Branch survive.

## Context

- interface: `internal/daemon/task_engine.go`

## Verification

- `grep -q "budgetExceededReason" internal/daemon/task_engine.go` — expected: exit 0; the settlement names its cause. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestBudgetExceededRunRecordsItsReason$" ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Feature 3; User Story 2; Goal 2; Success Metrics 1-2;
Regression locks;
`_techspec.md` → Implementation Design: Settling the outcome; API Contracts 1
and 3; Build Order 2; ADR-0057; ADR-0158.
