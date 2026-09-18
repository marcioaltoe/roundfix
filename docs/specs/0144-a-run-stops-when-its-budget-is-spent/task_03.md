---
task: task_03
spec: 0144-a-run-stops-when-its-budget-is-spent
status: pending
type: backend
complexity: medium
---

# Task 03: Accept the outcome for carry-forward

## Overview

A Run that spent its budget after proving Tasks has work worth recovering. This
slice adds its outcome to carry-forward's accepted set, and changes nothing else
about what carry-forward proves.

## Requirements

1. MUST accept a `BudgetExceeded` Run for Task Carry-Forward, beside `Stopped`
   and `Unresolved`.
2. MUST keep refusing every terminal outcome outside that set.
3. MUST keep the existing proof requirements, the refusal of a set whose members
   cannot all be proved, and the absence of a force bypass.
4. MUST NOT change any other reconciliation classification, reason or summary
   count.

## Subtasks

- [ ] Add the outcome to the accepted set.
- [ ] Cover an accepted budget-expired Run, a refused outcome and an unprovable
      set.

## Acceptance Criteria

- [ ] A `BudgetExceeded` Run with proved Tasks is accepted.
- [ ] A `BudgetExceeded` Run whose Tasks cannot be proved is refused, as any
      other would be.
- [ ] An outcome outside the set is still refused.
- [ ] Every other reconciliation outcome is unchanged.

## Context

- interface: `internal/worktree/worktree.go`

## Verification

- `grep -q "StateBudgetExceeded" internal/worktree/worktree.go` — expected: exit 0; carry-forward knows the outcome. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestCarryForwardAcceptsBudgetExceededRun$" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Feature 4; User Story 3; Goal 3; Success Metric 3;
Regression locks;
`_techspec.md` → Implementation Design: Accepting the outcome for carry-forward;
API Contract 2; Build Order 3; ADR-0158.
