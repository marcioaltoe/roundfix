---
task: task_02
spec: 0144-a-run-stops-when-its-budget-is-spent
status: completed
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

## Result

The Task cycle now marks its own elapsed Run deadline with the existing
`BudgetExceeded` outcome and a `budgetExceededReason` that names the configured
maximum and elapsed duration. The Implement command settles that outcome before
the generic Stop Request mapping, journals the specific reason, and preserves
the non-integrated Run Worktree and Run Branch.

Acceptance evidence:

- The daemon-owned `TestBudgetExceededRunRecordsItsReason` observed a bounded
  Task cycle report `BudgetExceeded` with a reason containing both `configured
  maximum 500ms` and `elapsed`, within the existing public reason bound.
- That daemon fixture completed and committed `task_01` before the bound, then
  confirmed its `completed` status remained while interrupted `task_02`
  retained the existing `in_progress` settlement.
- `TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch` confirmed the
  Run Worktree directory and deterministic Run Branch still existed, stderr
  named the kept Run Worktree, and the command report rendered the interrupted
  Task pending rather than completed or failed.
- The focused regression set also passed the existing Stopped, Unresolved,
  disabled-budget, finishes-before-deadline, and cancellation cases.

Focused checks:

- Before the implementation, `GOCACHE=/tmp/roundfix-go-cache rtk go test
  -count=1 -run
  '^(TestBudgetExceededRunRecordsItsReason|TestRunImplementStopRequestEndsStoppedWithInterruptMapping)$'
  ./internal/cli` failed because the bounded Run settled `Stopped` with exit 0.
- `GOCACHE=/tmp/roundfix-go-cache rtk go test -count=1 -run
  '^(TestBudgetExceededRunRecordsItsReason|TestRunImplementStopRequestEndsStoppedWithInterruptMapping|TestImplementTaskStatusFailureEndsUnresolvedAndKeepsWorktree|TestTaskCycleEndsRunAtBudgetDeadline|TestTaskCycleDisabledBudgetDerivesNoDeadline|TestTaskCycleFinishesBeforeBudgetDeadline)$'
  ./internal/cli ./internal/daemon` — passed, 8 tests.
- `GOCACHE=/tmp/roundfix-go-cache rtk go test -count=1 ./internal/cli
  ./internal/daemon` — 1,486 tests passed and 3 failed. The first version of
  the new fixture used a parallel 100 ms bound and failed because its first Task
  could miss that bound under suite load; the fixture now runs non-parallel with
  500 ms and passes the focused regression. The other failures were the
  unrelated process-integration tests
  `TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner` and
  `TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion`.
- `rtk git diff --check` — passed.

Verification Feedback repair:

- Attempt 1's diagnostic artifact was present and empty. Inspection found the
  exact named regression under `internal/cli`, while the authored Verification
  intentionally searches `internal/daemon`.
- The exact test now lives beside `TaskCycle`; the CLI-level preservation test
  has its own surface-specific name.
- `GOCACHE=/tmp/roundfix-go-cache rtk go test -count=1 -run
  '^(TestBudgetExceededRunRecordsItsReason|TestTaskCycleEndsRunAtBudgetDeadline|TestTaskCycleDisabledBudgetDerivesNoDeadline|TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch|TestRunImplementStopRequestEndsStoppedWithInterruptMapping)$'
  ./internal/daemon ./internal/cli` — passed, 7 tests.

The Daemon-owned commands under `## Verification` were not run in this Agent
turn.
