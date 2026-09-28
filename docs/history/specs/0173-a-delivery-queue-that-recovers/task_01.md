---
task: task_01
spec: 0173-a-delivery-queue-that-recovers
status: completed
type: backend
complexity: medium
---

# Task 01: Carry-forward follows the Run's integration order

## Overview

`inspectCarryForwardsForRun` in `internal/cli/carryforward.go` stages and proves the settled Tasks of a Run in `graph.Tasks` order, while ADR-0026 integrated them onto the Run Branch in completion order. When a Task that settled later in the Run comes earlier in the Task Graph, its change is staged before a Task whose settlement never saw it, and that Task's declared input reads as moved. Run `run_20260925T182023Z_7e3e8c4b8bc65d93` settled `task_04` before `task_01`; `task_01` changed `CONTEXT.md`, which `task_04` declares, and `roundfix reconcile <run> --carry-forward` refused the whole set. The settled commits are read from the Run Worktree's history, the staged state is built in a temporary Worktree from the checkout, and the result is read by the operator through `reconcile` and by `implement`'s carry-forward inspection.

## Requirements

1. MUST make `carryForwardTaskCommits` report, beside the per-Task commit map, the Task IDs in the position of their first settlement commit in `git rev-list --reverse <run.HeadSHA>..HEAD` of the Run Worktree.
2. MUST make `inspectCarryForwardsForRun` iterate the settled-completed Tasks in that order; a settled-completed Task with no settlement commit keeps its refusal and follows the ordered Tasks in Task Graph order.
3. MUST make `unstagedCarryForwards` list the Tasks after a Task that could not be staged in that same integration order.
4. MUST keep every carry-forward proof and its whole-set refusal unchanged: an input changed on the checkout after the Run still refuses the whole set and leaves HEAD unchanged, and the existing `TestCarryForwardInputsResolveAgainstStagedCarries` and `TestInspectSpecCarryForwards` stay green without edits.
5. MUST state in the reconcile section of `docs/user-guide/commands.md` that carry-forward proves the candidates in the order the Run integrated them (the phrase `in the order the Run integrated them`) and that the `carryForwards` JSON array follows that order.
6. MUST put the new tests in `internal/cli/carryforward_integration_order_test.go`, building a Run whose Run Branch holds `task_02`'s settlement commit before `task_01`'s, and drive them through `roundfix reconcile <run-id> --carry-forward`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A Run whose Tasks settled out of Task Graph order, one declaring an input the other changed, carries both Tasks with exit `0`, and the carried commits follow the Run's order.
- [ ] A checkout change to that input after the Run still refuses the whole set with exit `2`, naming the input, with HEAD unchanged.
- [ ] After a conflicting Task, the Tasks reported as not evaluated follow integration order.

## Context

- interface: `internal/cli/carryforward.go`
- creates: `internal/cli/carryforward_integration_order_test.go`
- interface: `docs/user-guide/commands.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCarryForwardProvesTasksInTheOrderTheRunIntegratedThem|TestCarryForwardStillRefusesAnInputMovedOnTheCheckout|TestCarryForwardUnevaluatedTasksFollowIntegrationOrder|TestCarryForwardInputsResolveAgainstStagedCarries|TestInspectSpecCarryForwards)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCarryForwardProvesTasksInTheOrderTheRunIntegratedThem TestCarryForwardStillRefusesAnInputMovedOnTheCheckout TestCarryForwardUnevaluatedTasksFollowIntegrationOrder TestCarryForwardInputsResolveAgainstStagedCarries TestInspectSpecCarryForwards; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "in the order the Run integrated them"` — expected: exit 0; before this Task none of the three new named tests exists and the phrase is not documented, so the command fails.

## References

- [_techspec.md](_techspec.md) — Carry-forward in integration order

## Result

Implementation:

- `carryForwardTaskCommits` now returns each Task's commits and the Task IDs at
  their first settlement-commit positions in the Run Worktree's reverse
  revision walk.
- Carry-forward builds one settled-completed Task sequence from that Run
  integration order, appends Tasks without settlement commits in Task Graph
  order, and reuses the sequence for proof, reporting, staging, and conflict
  remainder accounting.
- The Reconcile Command guide states that candidates are proved in the order
  the Run integrated them and that the `carryForwards` JSON array follows that
  order.

Focused checks:

- Pre-change signal: `rtk env GOCACHE=/tmp/roundfix-0173-task01-gocache go test -count=1 -run '^TestCarryForwardProvesTasksInTheOrderTheRunIntegratedThem$' ./internal/cli`
  failed as expected with exit `1`; `task_02` was evaluated after `task_01`
  and refused because `_prd.md` appeared moved.
- `rtk env GOCACHE=/tmp/roundfix-0173-task01-gocache go test -count=1 -run '^TestCarryForward(ProvesTasksInTheOrderTheRunIntegratedThem|StillRefusesAnInputMovedOnTheCheckout|UnevaluatedTasksFollowIntegrationOrder)$' ./internal/cli`
  passed.
- `rtk env GOCACHE=/tmp/roundfix-0173-task01-gocache go test -count=1 -run '^(TestCarryForwardInputsResolveAgainstStagedCarries|TestInspectSpecCarryForwards)$' ./internal/cli`
  passed without changing those tests.
- `rtk env GOCACHE=/tmp/roundfix-0173-task01-gocache go test -count=1 -run '^(TestCarryForward.*|TestInspectSpecCarryForwards)$' ./internal/cli`
  passed.
- `rtk git diff --check` passed, and `rtk rg -n -F "in the order the Run integrated them" docs/user-guide/commands.md`
  found the required phrase beside the `carryForwards` ordering contract.

Acceptance evidence:

- `TestCarryForwardProvesTasksInTheOrderTheRunIntegratedThem` builds Task Graph
  order `task_01, task_02` and Run integration order `task_02, task_01`; the
  public reconcile command exits `0`, carries both, reports that order in
  JSON, and leaves the carried commit trailers in that order.
- `TestCarryForwardStillRefusesAnInputMovedOnTheCheckout` changes the shared
  input after the Run; reconcile exits `2`, names the input in its whole-set
  refusal, leaves HEAD unchanged, and carries neither Task.
- `TestCarryForwardUnevaluatedTasksFollowIntegrationOrder` makes the first
  integrated Task conflict; the JSON report and refusal diagnostic list the
  remaining unevaluated Tasks as `task_03, task_01`, matching Run integration
  order rather than Task Graph order.

Not run:

- The Task's declared `## Verification` command; the Roundfix Daemon owns that
  command and the terminal Task verdict.
