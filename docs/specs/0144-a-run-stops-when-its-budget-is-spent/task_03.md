---
task: task_03
spec: 0144-a-run-stops-when-its-budget-is-spent
status: completed
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

- [x] Add the outcome to the accepted set.
- [x] Cover an accepted budget-expired Run, a refused outcome and an unprovable
      set.

## Acceptance Criteria

- [x] A `BudgetExceeded` Run with proved Tasks is accepted.
- [x] A `BudgetExceeded` Run whose Tasks cannot be proved is refused, as any
      other would be.
- [x] An outcome outside the set is still refused.
- [x] Every other reconciliation outcome is unchanged.

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

## Result

Implementation:

- `internal/worktree` now owns the accepted Task Carry-Forward outcome list:
  `Stopped`, `Unresolved`, and `BudgetExceeded`.
- Both the explicit Reconcile Command guard and the Spec-scoped carry-forward
  query read that list. The existing Task proofs, all-or-nothing refusal, and
  apply path are unchanged.
- Reconciliation coverage exercises a proved `BudgetExceeded` Run, moved-input
  and mixed-set refusal for that outcome, and refusal of `Clean` while naming
  all three accepted outcomes.

Focused checks:

- The pre-change focused regression check refused the `BudgetExceeded` fixture
  with exit 2, and the new worktree policy test could not compile because
  `CarryForwardAcceptedOutcomes` did not exist.
- `GOCACHE=/tmp/roundfix-task03-go-cache rtk go test -count=1 -run
  'TestCarryForward(Accepts|Refuses|Without)' ./internal/cli` passed 16 tests.
- `GOCACHE=/tmp/roundfix-task03-go-cache rtk go test -count=1 -run
  'TestCarryForwardAcceptsBudgetExceededRun|TestApplyTerminalRunOutcomeRemainsUnchanged'
  ./internal/worktree` passed 6 tests.
- The sandboxed `GOCACHE=/tmp/roundfix-task03-go-cache rtk make
  verify-incremental` reached the complete Go suite and failed only because two
  existing force-stop integration tests could not read the process table.
  `rtk env GOCACHE=/tmp/roundfix-task03-go-cache make verify-incremental`
  then passed outside that process sandbox, including all Go packages, skill
  checks, and the build.

Acceptance evidence:

- `TestCarryForwardAcceptsBudgetExceededRun` drives the public Reconcile Command
  with a proved Task, observes exit 0 and an empty stderr, and proves that the
  checkout fast-forwards with the implementation, completed status, Run ID,
  and settlement commit.
- `TestCarryForwardRefusesATaskWhoseInputsMoved/BudgetExceeded` and
  `TestCarryForwardRefusesRatherThanCarryingASubset/BudgetExceeded` observe exit
  2, unchanged checkout and Task files, and no carried implementation when one
  Task proof fails.
- `TestCarryForwardRefusesAnUnacceptedRunOutcomeByName` still refuses `Clean`
  and names `Stopped`, `Unresolved`, and `BudgetExceeded` as the accepted set.
- No reconciliation classification, reason, or summary code changed.
  `TestApplyTerminalRunOutcomeRemainsUnchanged` passed, and the incremental gate
  passed after the final implementation edit.
