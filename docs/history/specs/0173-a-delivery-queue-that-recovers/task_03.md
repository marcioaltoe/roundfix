---
task: task_03
spec: 0173-a-delivery-queue-that-recovers
status: completed
type: backend
complexity: high
---

# Task 03: The engine can retry a parked item

## Overview

`Engine.Run` in `internal/delivery/engine.go` skips every item whose stage is `parked`, and nothing moves an item out of that stage. `runCandidate` also parks a `run-unresolved` item before it records the Run that parked it: the live Run Database holds all three wave-one items parked `run-unresolved` with an empty Run ID. The owner releases the queue in `internal/cli` without looking at the items, so an item re-staged while the owner finishes would be stranded. The item state lives in the Run Database's `delivery_queue_items` rows, written by the owner process and by the operator's retry command, which run as separate processes; the recovery facts come from the item worktree through a new boundary the CLI implements in a later Task.

## Requirements

1. MUST make `runCandidate` record `strings.TrimSpace(result.RunID)` on the item before it parks it `run-unresolved`.
2. MUST add `Store.RetryDeliveryQueueItem(ctx, gitRoot string, item DeliveryQueueItem, parkedBlocker string) (int, string, error)` to `internal/store/delivery.go`: it refuses a target stage other than `running`, `reviewing`, `gating` or `checking`; it refuses, naming the stored stage and blocker and changing nothing, unless the stored item at the item's position is `parked` with exactly `parkedBlocker`; otherwise, in one write transaction, it persists the target stage, an empty blocker, the Run ID and the candidate commits, and returns the queue's recorded owner PID and identity read in that transaction.
3. MUST add `Store.ReleaseIdleDeliveryQueueOwner(ctx, gitRoot string, pid int, identity string) (bool, error)`: in one write transaction it clears the owner and returns `true` only when the caller is the recorded owner and no item is in a stage other than `merged` or `parked`; with an advanceable item it returns `false` and keeps the claim; when the caller is not the recorded owner it returns an error naming the recorded owner.
4. MUST add to `internal/delivery/engine.go` the `ItemState`, `CarryForwardResult`, `ItemRecovery` and `RetryResult` types and `Engine.Retry(ctx, gitRoot, specSlug string) (RetryResult, error)` with the signatures in the TechSpec, and a `Recovery ItemRecovery` field on `EngineDependencies` that `Run` does not require and `Retry` does.
5. MUST make `Retry` refuse an item that is absent from the queue or not `parked`, and an item for which `UseItemBranch` returns `ErrItemWorktreeMissing`, naming the item branch.
6. MUST make `Retry` derive the re-entry stage from `InspectItem` exactly as the TechSpec's re-entry table states: for an active Spec it first calls `CarryForward` with the item worktree, slug, item branch and recorded Run ID, records the returned Run ID, inspects again, and chooses `running` when a Task is unfinished or `reviewing` otherwise, appending the item head to the candidate commits when it differs from the last one; for an archived Spec it refuses when the item head differs from the candidate head and otherwise chooses `gating` without a recorded pull request and `checking` with one.
7. MUST make `Retry` wrap a `CarryForward` error with `%w` and refuse, and MUST make every refusal happen before `RetryDeliveryQueueItem`, so a refused retry leaves the stored item unchanged.
8. MUST keep `TestDeliveryEngineParksDeclaredBlockersAndContinues`, `TestDeliveryEngineRetriesMergedItemCleanupWithoutReplayingMerge`, `TestDeliveryEngineResumesWithoutDoubleEffects` and `TestDeliveryQueueRoundTripsItemsAndReceipts` green without edits.
9. MUST put the new tests in `internal/delivery/retry_test.go` and `internal/store/delivery_retry_test.go`, with a fake `ItemRecovery` in the delivery test file, and prove that a retried archived item pushes, opens its pull request and merges exactly once.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] An item parked `run-unresolved` records the Run ID the executor returned.
- [ ] The store moves only a parked item with the expected blocker to a retry stage and returns the recorded owner; a changed blocker, a non-parked item and a terminal target stage are refused with the row unchanged.
- [ ] The idle release frees an idle queue, keeps the claim while an item is advanceable, and refuses a caller that is not the owner.
- [ ] `Retry` carries forward then re-enters `running` with a Task unfinished, re-enters `reviewing` at the current head with none, re-enters `gating` or `checking` for an archived item at its candidate head, and refuses a moved archived head, a refused carry-forward, a non-parked item and a missing branch with the item unchanged.

## Context

- interface: `internal/delivery/engine.go`
- interface: `internal/store/delivery.go`
- creates: `internal/delivery/retry_test.go`
- creates: `internal/store/delivery_retry_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAnUnresolvedRunParksWithItsRunID|TestRetryCarriesForwardBeforeReenteringTheRun|TestRetryReentersTheRunWhenATaskIsUnfinished|TestRetryReReviewsTheCurrentHeadWhenNoTaskIsUnfinished|TestRetryAfterArchiveReentersGatingWithoutAPullRequest|TestRetryAfterArchiveReentersCheckingWithAPullRequest|TestRetriedArchivedItemPublishesAndMergesOnce|TestRetryRefusesAnArchivedItemWhoseHeadMoved|TestRetryRefusesWhenCarryForwardRefuses|TestRetryRefusesAnItemThatIsNotParked|TestRetryRefusesAnItemWhoseBranchIsGone|TestRetryRefusesAnItemMissingFromTheQueue|TestRetryDeliveryQueueItemReentersAParkedItem|TestRetryDeliveryQueueItemRefusesAnItemThatIsNotParked|TestRetryDeliveryQueueItemRefusesAChangedBlocker|TestRetryDeliveryQueueItemRefusesATerminalStage|TestReleaseIdleDeliveryQueueOwnerReleasesAnIdleQueue|TestReleaseIdleDeliveryQueueOwnerKeepsAQueueWithAnAdvanceableItem|TestReleaseIdleDeliveryQueueOwnerRefusesAnotherOwner|TestDeliveryEngineParksDeclaredBlockersAndContinues|TestDeliveryEngineRetriesMergedItemCleanupWithoutReplayingMerge|TestDeliveryEngineResumesWithoutDoubleEffects|TestDeliveryQueueRoundTripsItemsAndReceipts)$" ./internal/delivery ./internal/store 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAnUnresolvedRunParksWithItsRunID TestRetryCarriesForwardBeforeReenteringTheRun TestRetryReentersTheRunWhenATaskIsUnfinished TestRetryReReviewsTheCurrentHeadWhenNoTaskIsUnfinished TestRetryAfterArchiveReentersGatingWithoutAPullRequest TestRetryAfterArchiveReentersCheckingWithAPullRequest TestRetriedArchivedItemPublishesAndMergesOnce TestRetryRefusesAnArchivedItemWhoseHeadMoved TestRetryRefusesWhenCarryForwardRefuses TestRetryRefusesAnItemThatIsNotParked TestRetryRefusesAnItemWhoseBranchIsGone TestRetryRefusesAnItemMissingFromTheQueue TestRetryDeliveryQueueItemReentersAParkedItem TestRetryDeliveryQueueItemRefusesAnItemThatIsNotParked TestRetryDeliveryQueueItemRefusesAChangedBlocker TestRetryDeliveryQueueItemRefusesATerminalStage TestReleaseIdleDeliveryQueueOwnerReleasesAnIdleQueue TestReleaseIdleDeliveryQueueOwnerKeepsAQueueWithAnAdvanceableItem TestReleaseIdleDeliveryQueueOwnerRefusesAnotherOwner TestDeliveryEngineParksDeclaredBlockersAndContinues TestDeliveryEngineRetriesMergedItemCleanupWithoutReplayingMerge TestDeliveryEngineResumesWithoutDoubleEffects TestDeliveryQueueRoundTripsItemsAndReceipts; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the nineteen new named tests exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Retry in the engine and store
