---
task: task_01
spec: 0180-a-prepared-queue-that-revalidates-before-each-spec
status: pending
type: backend
complexity: high
---

# Task 01: A queued Spec is revalidated on its starting main before its first Run

## Overview

`advanceItem` in `internal/delivery/engine.go` creates a queued item's worktree from the refreshed default branch and moves the item straight to `running`. A later Spec therefore runs against a main that earlier items of the same queue changed, with no re-check of its authored premises. Spec 0175 declares `internal/delivery/engine.go`, `internal/cli/deliver_workflow.go` and `internal/cli/carryforward.go`, and Spec 0173's merge `6fac37ea` changed all three after Spec 0175 was authored. This Task adds Delivery Revalidation between worktree creation and the first Run. A strict consistency finding parks the item. A declared production Go file changed by an earlier merge only records a `premise-changed` warning, because such overlaps are common and usually harmless. The check reads the Spec's artifacts and the earlier merge commits from the item worktree. It writes to the Delivery Queue item, which the operator reads through `roundfix deliver status`, and to the delivery console log.

## Requirements

1. MUST add to `internal/delivery/engine.go`, exactly as the TechSpec states:
   - the constants `BlockerRevalidationFailed = "revalidation-failed"` and `WarningPremiseChanged = "premise-changed"`;
   - the `Revalidation` type and the `ItemRevalidator` interface;
   - `Revalidator ItemRevalidator` and `Log io.Writer` in `EngineDependencies` and the engine, where a nil `Log` discards.

   Both `validate` and `validateRetry` MUST refuse an engine without a revalidator.
2. MUST make `Engine.Run` pass each item the merge commits of the items at lower positions whose stage is `merged` and whose merge commit is non-empty, in position order.
3. MUST call `Revalidate` in the `queued` branch of `advanceItem` after task_02's deadline check and after `CreateItemBranch` has recorded the branch and worktree, and before the stage becomes `running`. It parks the item as `revalidation-failed: <code>, <code>` when findings exist and otherwise sets `running`. A revalidation error is returned so that `Engine.Run` parks the item as `delivery-error` as today. The Runner MUST NOT be called for an item parked here.
4. MUST, whenever `ChangedPremises` is non-empty, set the item's `Warning` (task_02's field) to `premise-changed: <path>, <path> (merge <sha>, <sha>)` from `ChangedPremises` and `ChangedBy`, persist it with the next step, and write `roundfix: warning: Delivery Queue item <slug>: <warning>` to `Log`. A changed premise MUST NOT park the item or delay its Run. An item with no changed premise keeps an empty `Warning`, and nothing is written to `Log`.
5. MUST make `Engine.Retry`, for an item whose Spec is active on the item branch and whose recorded Run ID is empty, call `Revalidate` in the item worktree with no prior merges before carry-forward. It refuses, with an error naming the codes, while findings remain, and leaves the stored item unchanged. A retry MUST NOT change the recorded `Warning`.
6. MUST add `internal/cli/deliver_revalidate.go` with `(*commandDeliveryWorkflow).Revalidate` as the TechSpec states:
   - it resolves the Specs Root of the worktree and reports the unique, sorted codes of task_03's `strictSpecFindings`;
   - its premises are task_03's `productionPremises` of the Spec's graph;
   - it reads the name-only diff of each prior merge against its first parent through `workflow.git`, and reports the premises named and each merge that names one;
   - an unreadable merge commit is an error, never an empty change set.

   `newCommandDeliveryEngine` in `internal/cli/deliver_workflow.go` MUST pass the workflow as `Revalidator` and `os.Stderr` as `Log`.
7. MUST make `runDeliverStatus` in `internal/cli/deliver.go` print the item rows unchanged, then one `Warning: <slug> <warning>` line for each item with a non-empty `Warning`, in position order. A queue with no warning prints exactly the output it prints today.
8. MUST update the engine construction sites that the new required dependency invalidates, and only those:
   - the helpers `newTestDeliveryEngineWithWait` in `internal/delivery/engine_test.go` and `newRetryDeliveryEngine` in `internal/delivery/retry_test.go`, and any other `NewEngine` call in `internal/delivery/engine_test.go`, which also covers task_02's `internal/delivery/limits_test.go`, since that file builds its engines only through these helpers;
   - the `delivery.NewEngine` calls in `internal/cli/deliver_test.go` and `internal/cli/deliver_recovery_test.go`.

   Each passes a fake returning a clean `Revalidation`; for the CLI tests, that is a `Revalidate` method on `parkTestDeliveryFlow`. No existing assertion changes.
9. MUST document in the deliver section of `docs/user-guide/commands.md` and the Delivery queue section of `.agents/skills/roundfix/SKILL.md`:
   - the revalidation and the `revalidation-failed` blocker;
   - the `premise-changed` warning, the `Warning:` status line and the console log line, and that the item continues to its Run;
   - that a Delivery Retry re-runs the strict check for an item that has not run.

   Then MUST regenerate `skills/roundfix/SKILL.md` with `make skills-sync`.
10. MUST put the engine tests in `internal/delivery/revalidate_test.go`, writing `Log` to a buffer. The workflow and status tests go in `internal/cli/deliver_revalidate_test.go`. The workflow tests run against real Git repositories, with an origin whose default branch receives the earlier merge, and may start from the `clean` fixture Spec that `newSpecCheckWorkspace` copies.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A finding parks the item as `revalidation-failed` before any Run, with the worktree recorded, and a revalidation error parks it as `delivery-error`.
- [ ] A changed premise records the `premise-changed` warning naming the files and the merge, writes it to the log and reaches the Run. An item with no overlap reaches the Run with an empty warning and no log line.
- [ ] The revalidator receives the merge commits of earlier merged items in position order, and never those of parked or unfinished items.
- [ ] A retry of an item with no Run refuses while findings remain, leaving the item unchanged. It proceeds once the check is clean, and it keeps the recorded warning.
- [ ] Against real Git, a declared file that an earlier merge removed yields `SC-REF-UNRESOLVED` only after that merge. A declared production Go file changed by that merge is named with the merge, while a changed test, guide or undeclared file is not. An unreadable merge commit is an error.
- [ ] `deliver status` prints a `Warning:` line for an item with a warning and none for a queue without one.
- [ ] An engine without a revalidator refuses to run and to retry.

## Context

- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/engine_test.go`
- interface: `internal/delivery/retry_test.go`
- creates: `internal/delivery/revalidate_test.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_revalidate.go`
- creates: `internal/cli/deliver_revalidate_test.go`
- interface: `internal/cli/deliver_test.go`
- interface: `internal/cli/deliver_recovery_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRevalidationParksAnItemWhoseStartingMainFailsTheCheck|TestAChangedPremiseWarnsAndContinues|TestNoChangedPremiseRecordsNoWarning|TestRevalidationReceivesEarlierMergeCommitsInQueueOrder|TestARevalidationErrorParksTheItemAsADeliveryError|TestRetryOfAnItemWithoutARunRefusesWhileFindingsRemain|TestRetryOfAnItemWithoutARunProceedsOnceTheCheckIsClean|TestRetryKeepsTheRecordedPremiseWarning|TestEngineRefusesToRunWithoutARevalidator|TestRevalidateReportsAnUnresolvedDeclarationOnlyAfterTheMergeRemovedIt|TestRevalidateNamesADeclaredProductionFileAndTheMergeThatChangedIt|TestRevalidateIgnoresChangedTestsGuidesAndUndeclaredFiles|TestRevalidateRefusesAnUnreadableMergeCommit|TestDeliverStatusPrintsAnItemWarning|TestDeliverStatusPrintsNoWarningLineWithoutAWarning|TestDeliveryEngineMergesAQueueWithoutAnOperator|TestRetryRefusesAnItemMissingFromTheQueue|TestAQueuedItemParksAtTheDeadlineWithoutAWorktree|TestDeliverNeverTouchesTheUserCheckout|TestParkLeavesTheItemWorktreeInPlace|TestRetriedItemRunsWithItsCompletedTasksCarried|TestDeliverStatusPrintsTheItemWorktree)$" ./internal/delivery ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRevalidationParksAnItemWhoseStartingMainFailsTheCheck TestAChangedPremiseWarnsAndContinues TestNoChangedPremiseRecordsNoWarning TestRevalidationReceivesEarlierMergeCommitsInQueueOrder TestARevalidationErrorParksTheItemAsADeliveryError TestRetryOfAnItemWithoutARunRefusesWhileFindingsRemain TestRetryOfAnItemWithoutARunProceedsOnceTheCheckIsClean TestRetryKeepsTheRecordedPremiseWarning TestEngineRefusesToRunWithoutARevalidator TestRevalidateReportsAnUnresolvedDeclarationOnlyAfterTheMergeRemovedIt TestRevalidateNamesADeclaredProductionFileAndTheMergeThatChangedIt TestRevalidateIgnoresChangedTestsGuidesAndUndeclaredFiles TestRevalidateRefusesAnUnreadableMergeCommit TestDeliverStatusPrintsAnItemWarning TestDeliverStatusPrintsNoWarningLineWithoutAWarning TestDeliveryEngineMergesAQueueWithoutAnOperator TestRetryRefusesAnItemMissingFromTheQueue TestAQueuedItemParksAtTheDeadlineWithoutAWorktree TestDeliverNeverTouchesTheUserCheckout TestParkLeavesTheItemWorktreeInPlace TestRetriedItemRunsWithItsCompletedTasksCarried TestDeliverStatusPrintsTheItemWorktree; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "revalidation-failed" && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "premise-changed" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "revalidation-failed" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "premise-changed" && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the fifteen new named tests exists and neither the blocker nor the warning is documented, so the command fails.

## References

- [_techspec.md](_techspec.md) — Strict consistency and Delivery Revalidation
- `_prd.md` → Goal 3; Core Feature 3; Success Metric 2
- `_techspec.md` → API Contracts 3-4; Testing Approach 1; Testing Approach 6

## Result
