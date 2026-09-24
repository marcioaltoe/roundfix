---
task: task_08
spec: 0156-a-delivery-loop-that-outlives-the-session
status: pending
type: backend
complexity: high
---

# Task 08: Publish only the reviewed head, from each item's own branch

## Overview

Corrective Task from the pre-PR review of 2026-09-24. Every item publishes from the current checkout's branch, so a parked item's commits and open pull request are picked up and squash-merged by the next item. The push sends the checkout's live `HEAD`, and compares it with the reviewed head only after publishing, so an unreviewed commit, or local `main`, can be pushed. Pending checks park an item forever. An unmatched push intent against a moved remote wedges the queue. A finished queue can never be replaced. Any non-blocker error aborts every later item. An already merged pull request is accepted without checking the head it merged.

## Requirements

1. MUST give each item its own branch, created from the refreshed default branch before its Run, recorded on the item, and used for every later action on that item.
2. MUST refuse to adopt an open pull request that another item recorded, or whose head branch is not the item's branch.
3. MUST push the recorded reviewed SHA to the item's recorded branch (`<sha>:refs/heads/<branch>`), never the checkout's `HEAD`, and refuse when that branch is the base branch.
4. MUST wait, with a bounded timeout, while checks are pending or not yet reported, and park only on a failed or cancelled check or on the timeout.
5. MUST, for a push intent without a receipt, skip the push only when the remote branch already holds the expected SHA, and otherwise retry it.
6. MUST let `deliver start` replace a queue whose items are all merged or parked and which has no live owner, and still refuse while a live or unfinished queue exists.
7. MUST park an item on a persistent error with its reason, and continue with the next item.
8. MUST record an already merged pull request as merged only when its merged head is the reviewed head, and otherwise park the item as `review-stale`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Two items publish from two branches, and a parked item's pull request is never merged by another item.
- [ ] The pushed object is the recorded SHA even when the checkout moved, and the base branch is refused.
- [ ] Pending checks are awaited; an erroring item parks and the queue continues; a finished queue is replaced.
- [ ] An unmatched push intent against a moved remote retries; an already merged pull request at another head parks.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/github.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/store/delivery.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestEachItemPublishesFromItsOwnBranch|TestAPullRequestOfAnotherItemIsNeverReused|TestPushPublishesTheRecordedHeadNotTheCheckout|TestPushRefusesTheBaseBranch|TestPendingChecksAreAwaitedNotParked|TestUnmatchedPushIntentRetriesWhenTheRemoteDiffers|TestATerminalQueueIsReplacedByANewStart|TestAnItemErrorParksAndTheQueueContinues|TestAnAlreadyMergedPullRequestMustCarryTheReviewedHead)$" ./internal/delivery ./internal/store ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEachItemPublishesFromItsOwnBranch TestAPullRequestOfAnotherItemIsNeverReused TestPushPublishesTheRecordedHeadNotTheCheckout TestPushRefusesTheBaseBranch TestPendingChecksAreAwaitedNotParked TestUnmatchedPushIntentRetriesWhenTheRemoteDiffers TestATerminalQueueIsReplacedByANewStart TestAnItemErrorParksAndTheQueueContinues TestAnAlreadyMergedPullRequestMustCarryTheReviewedHead; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
