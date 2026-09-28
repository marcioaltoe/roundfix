---
task: task_04
spec: 0173-a-delivery-queue-that-recovers
status: pending
type: backend
complexity: high
---

# Task 04: The delivery workflow recovers an item's Run

## Overview

An Implement Run that ends Unresolved leaves its settled Task commits on its Run Branch; the item branch that `roundfix deliver` created never receives them, and `implement` refuses to start again while a carriable Run exists. In the first wave the operator fast-forwarded each item branch to its Run Branch by hand. `commandDeliveryWorkflow` in `internal/cli/deliver_workflow.go` now implements the `delivery.ItemRecovery` boundary from task_03: it inspects the item worktree and carries the item's Run forward onto the item branch with the existing Task Carry-Forward proofs. The Run and its settlement evidence come from the Run Database and the Run Worktree; the carried commits land on the item branch, which the next Implement Run and the pre-PR review read.

## Requirements

1. MUST make `commandDeliveryWorkflow` implement `delivery.ItemRecovery` and make `newCommandDeliveryEngine` pass it as `Recovery`.
2. MUST make `InspectItem` resolve the Specs Root for the item worktree with `roundconfig.ResolveSpecsRoot(workflow.loaded, workDir)`; when the Spec folder exists there, load its Task Graph and list, in Task Graph order, every Task whose status is not `completed`; when it does not, report `Archived` when the destination `archivePaths` names exists and fail otherwise; and read the head with `git rev-parse HEAD`.
3. MUST make `CarryForward` resolve the Run: a recorded Run ID must name an Implement Run of this Spec, or it fails naming the Run; with no recorded Run ID it takes the newest Implement Run of this Spec whose recorded local branch is the item branch; with none it carries nothing.
4. MUST make `CarryForward` carry nothing when the Run's outcome is not one `runworktree.CarryForwardAcceptedOutcomes` returns, when `loadReconcileTaskCoverage` finds no Task the Run settled completed, or when every such Task is already `completed` in the item worktree.
5. MUST otherwise make `CarryForward` refuse when the Run Worktree is gone or the Specs Root is external, run `inspectCarryForwards` with the item worktree as the repository resolved through `filepath.EvalSymlinks`, refuse with `carryForwardRefusalReason` when any candidate refuses, and apply the set with `applyCarryForwards`; the result names the Run and the carried Task IDs in integration order.
6. MUST make every carry-forward refusal an error whose `NextAction` names `roundfix reconcile <run-id> --carry-forward` run in the item worktree and `roundfix deliver retry <slug>`, and leave the item branch unchanged when it refuses.
7. MUST put the new tests in `internal/cli/deliver_recovery_test.go` against real Git repositories, reusing the existing carry-forward and delivery-branch fixtures of package `cli`, including one that runs `Engine.Retry` then `Engine.Run` with this workflow as `Recovery` and a fake Implement executor, and asserts that the executor ran in the item worktree with the previous Run's completed Tasks already `completed` there.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] An Unresolved Run's completed Tasks reach the item branch with carry-forward provenance while its failed Task stays pending there.
- [ ] With no recorded Run ID, the Run is found by the item branch; with no Implement Run, nothing is carried and nothing fails.
- [ ] A Run already carried is carried again as nothing, with the item head unchanged.
- [ ] A moved input and a gone Run Worktree each refuse with the item head unchanged and a next action naming both commands.
- [ ] An active Spec lists its unfinished Tasks and an archived Spec reports `Archived`.
- [ ] A retried item's next Implement Run starts with the previous Run's completed Tasks already `completed`.

## Context

- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_recovery_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestItemRecoveryCarriesAnUnresolvedRunIntoTheItemWorktree|TestItemRecoveryFindsTheRunByItemBranchWhenNoneIsRecorded|TestItemRecoveryCarriesNothingWithoutAnImplementRun|TestItemRecoveryCarriesNothingForARunAlreadyCarried|TestItemRecoveryRefusesACarryForwardWithAMovedInput|TestItemRecoveryRefusesAGoneRunWorktree|TestInspectItemListsUnfinishedTasks|TestInspectItemReportsAnArchivedSpec|TestRetriedItemRunsWithItsCompletedTasksCarried)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestItemRecoveryCarriesAnUnresolvedRunIntoTheItemWorktree TestItemRecoveryFindsTheRunByItemBranchWhenNoneIsRecorded TestItemRecoveryCarriesNothingWithoutAnImplementRun TestItemRecoveryCarriesNothingForARunAlreadyCarried TestItemRecoveryRefusesACarryForwardWithAMovedInput TestItemRecoveryRefusesAGoneRunWorktree TestInspectItemListsUnfinishedTasks TestInspectItemReportsAnArchivedSpec TestRetriedItemRunsWithItsCompletedTasksCarried; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named tests exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Item recovery in the delivery workflow
