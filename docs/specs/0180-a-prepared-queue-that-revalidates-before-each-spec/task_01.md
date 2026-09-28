---
task: task_01
spec: 0180-a-prepared-queue-that-revalidates-before-each-spec
status: pending
type: backend
complexity: high
---

# Task 01: A queued Spec is revalidated on its starting main before its first Run

## Overview

`advanceItem` in `internal/delivery/engine.go` creates a queued item's worktree from the refreshed default branch and moves the item straight to `running`. A later Spec therefore runs against a main that earlier items of the same queue changed, with no re-check of its authored premises. Spec 0175 declares `internal/delivery/engine.go`, `internal/cli/deliver_workflow.go` and `internal/cli/carryforward.go`, and Spec 0173's merge `6fac37ea` changed all three after Spec 0175 was authored. This Task adds Delivery Revalidation between worktree creation and the first Run. It reads the Spec's artifacts and the earlier merge commits from the item worktree, and writes the verdict to the Delivery Queue item that the operator reads through `roundfix deliver status`.

## Requirements

1. MUST add to `internal/delivery/engine.go` the constants `BlockerRevalidationFailed = "revalidation-failed"` and `BlockerPremiseChanged = "premise-changed"`, the `Revalidation` type and the `ItemRevalidator` interface exactly as the TechSpec states, add `Revalidator ItemRevalidator` to `EngineDependencies` and the engine, and make both `validate` and `validateRetry` refuse an engine without one.
2. MUST make `Engine.Run` pass each item the merge commits of the items at lower positions whose stage is `merged` and whose merge commit is non-empty, in position order.
3. MUST call `Revalidate` in the `queued` branch of `advanceItem` after `CreateItemBranch` has recorded the branch and worktree and before the stage becomes `running`. It parks the item as `revalidation-failed: <code>, <code>` when findings exist, otherwise as `premise-changed: <path>, <path>` when changed premises exist, and otherwise sets `running`. A revalidation error is returned so that `Engine.Run` parks the item as `delivery-error` as today. The Runner MUST NOT be called for an item parked here.
4. MUST make `Engine.Retry`, for an item whose Spec is active on the item branch and whose recorded Run ID is empty, call `Revalidate` in the item worktree with no prior merges before carry-forward. It refuses, with an error naming the codes, while findings remain, and leaves the stored item unchanged. Changed premises MUST NOT refuse a retry.
5. MUST add `internal/cli/deliver_revalidate.go` with `strictSpecFindings(specsRoot, repoRoot, specSlug string) ([]speccheck.Finding, error)`, which composes `speccheck.Check`, `speccheck.PromoteGaps` and `speccheck.GatePrecondition(...).Findings`. The same file adds `(*commandDeliveryWorkflow).Revalidate` as the TechSpec states:
   - it resolves the Specs Root of the worktree and reports unique, sorted finding codes;
   - its premises are the production Go `interface:` paths of non-`qa` Tasks;
   - it reads `git diff --name-only <merge>^ <merge>` through `workflow.git`, and an unreadable merge commit is an error, never an empty change set.

   `newCommandDeliveryEngine` in `internal/cli/deliver_workflow.go` MUST pass the workflow as `Revalidator`.
6. MUST update the engine construction sites that the new required dependency invalidates, and only those: the helpers in `internal/delivery/engine_test.go` and `internal/delivery/retry_test.go`, and the `delivery.NewEngine` calls in `internal/cli/deliver_test.go` and `internal/cli/deliver_recovery_test.go`. Each passes a fake returning a clean `Revalidation`; for the CLI tests, that is a `Revalidate` method on `parkTestDeliveryFlow`. No existing assertion changes.
7. MUST document in the deliver section of `docs/user-guide/commands.md` and the Delivery queue section of `.agents/skills/roundfix/SKILL.md`:
   - the revalidation;
   - both blockers, by the strings `revalidation-failed` and `premise-changed`;
   - that a Delivery Retry re-runs the strict check for an item that has not run, and acknowledges a changed premise.

   Then MUST regenerate `skills/roundfix/SKILL.md` with `make skills-sync`.
8. MUST put the engine tests in `internal/delivery/revalidate_test.go`. The workflow tests go in `internal/cli/deliver_revalidate_test.go` and run against real Git repositories, with an origin whose default branch receives the earlier merge. They may start from the `clean` fixture Spec that `newSpecCheckWorkspace` copies.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A finding parks the item as `revalidation-failed` and a changed premise parks it as `premise-changed`, both before any Run and with the worktree recorded; a clean item runs; a revalidation error parks it as `delivery-error`.
- [ ] The revalidator receives the merge commits of earlier merged items in position order, and never those of parked or unfinished items.
- [ ] A retry of an item with no Run refuses while findings remain, leaving the item unchanged; it proceeds once the check is clean, and it proceeds for a `premise-changed` item.
- [ ] Against real Git, a declared file that an earlier merge removed yields `SC-REF-UNRESOLVED` only after that merge. A declared production Go file changed by that merge is named, while a changed test, guide or undeclared file is not. An unreadable merge commit is an error.
- [ ] An engine without a revalidator refuses to run and to retry.

## Context

- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/engine_test.go`
- interface: `internal/delivery/retry_test.go`
- creates: `internal/delivery/revalidate_test.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_revalidate.go`
- creates: `internal/cli/deliver_revalidate_test.go`
- interface: `internal/cli/deliver_test.go`
- interface: `internal/cli/deliver_recovery_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRevalidationParksAnItemWhoseStartingMainFailsTheCheck|TestRevalidationParksAnItemWhosePremiseChanged|TestRevalidationReceivesEarlierMergeCommitsInQueueOrder|TestACleanRevalidationStartsTheRun|TestARevalidationErrorParksTheItemAsADeliveryError|TestRetryOfAnItemWithoutARunRefusesWhileFindingsRemain|TestRetryOfAnItemWithoutARunProceedsOnceTheCheckIsClean|TestRetryAcknowledgesAChangedPremise|TestEngineRefusesToRunWithoutARevalidator|TestDeliveryEngineMergesAQueueWithoutAnOperator|TestRetryRefusesAnItemMissingFromTheQueue|TestRevalidateReportsAnUnresolvedDeclarationOnlyAfterTheMergeRemovedIt|TestRevalidateNamesADeclaredProductionFileAnEarlierMergeChanged|TestRevalidateIgnoresChangedTestsGuidesAndUndeclaredFiles|TestRevalidateRefusesAnUnreadableMergeCommit|TestDeliverNeverTouchesTheUserCheckout|TestParkLeavesTheItemWorktreeInPlace|TestRetriedItemRunsWithItsCompletedTasksCarried)$" ./internal/delivery ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRevalidationParksAnItemWhoseStartingMainFailsTheCheck TestRevalidationParksAnItemWhosePremiseChanged TestRevalidationReceivesEarlierMergeCommitsInQueueOrder TestACleanRevalidationStartsTheRun TestARevalidationErrorParksTheItemAsADeliveryError TestRetryOfAnItemWithoutARunRefusesWhileFindingsRemain TestRetryOfAnItemWithoutARunProceedsOnceTheCheckIsClean TestRetryAcknowledgesAChangedPremise TestEngineRefusesToRunWithoutARevalidator TestDeliveryEngineMergesAQueueWithoutAnOperator TestRetryRefusesAnItemMissingFromTheQueue TestRevalidateReportsAnUnresolvedDeclarationOnlyAfterTheMergeRemovedIt TestRevalidateNamesADeclaredProductionFileAnEarlierMergeChanged TestRevalidateIgnoresChangedTestsGuidesAndUndeclaredFiles TestRevalidateRefusesAnUnreadableMergeCommit TestDeliverNeverTouchesTheUserCheckout TestParkLeavesTheItemWorktreeInPlace TestRetriedItemRunsWithItsCompletedTasksCarried; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "revalidation-failed" && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "premise-changed" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "revalidation-failed" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "premise-changed" && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the thirteen new named tests exists and neither blocker is documented, so the command fails.

## References

- [_techspec.md](_techspec.md) — Strict consistency and Delivery Revalidation
- `_prd.md` → Goal 3; Core Feature 3; Success Metric 2
- `_techspec.md` → API Contracts 3-4; Testing Approach 1; Testing Approach 6

## Result
