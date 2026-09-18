---
task: task_01
spec: 0144-a-run-stops-when-its-budget-is-spent
status: pending
type: backend
complexity: medium
---

# Task 01: Bound the Run by its configured maximum

## Overview

An Implement Run has no deadline. This slice derives one from the Run's start
when the budget is enabled, checks it where the lifecycle already checks for a
Stop Request, and ends the Run through the cancellation that path already uses.

## Requirements

1. MUST derive the Run's deadline once, as its start plus the configured maximum,
   when the budget is enabled and the maximum is positive.
2. MUST check the deadline at the points where the lifecycle already checks for a
   Stop Request, and MUST NOT interrupt an operation the Daemon cannot end
   safely.
3. MUST start no further Task once the deadline has passed.
4. MUST cancel the Run's Agent Session and owned processes through the existing
   cancellation path, keeping its proof of termination and its reporting of an
   unprovable one.
5. MUST derive no deadline and change nothing when the budget is disabled.
6. MUST name the derived value `runBudgetDeadline`, so the Task's own check and
   later readers agree.

## Subtasks

- [ ] Derive the deadline at Run start from the configured budget.
- [ ] Check it beside the existing Stop Request checks.
- [ ] End the Run through the existing cancellation.
- [ ] Cover an expiring budget, a disabled budget and a Run that finishes first.

## Acceptance Criteria

- [ ] A fixture Run whose maximum expires during its work starts no further Task
      and cancels its Agent Session.
- [ ] A fixture Run with the budget disabled derives no deadline.
- [ ] A fixture Run that finishes before its maximum is unaffected.
- [ ] Cancellation still proves termination and still reports an unprovable one.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/engine.go`

## Verification

- `grep -q "runBudgetDeadline" internal/daemon/task_engine.go` — expected: exit 0; the lifecycle carries the Run's deadline. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestTaskCycleEndsRunAtBudgetDeadline$" ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Features 1-2 and 5; User Stories 1 and 4; Goals 1 and 4;
`_techspec.md` → Implementation Design: Deriving and honouring the deadline,
Cancelling what the Daemon owns; Build Order 1; ADR-0127; ADR-0158.
