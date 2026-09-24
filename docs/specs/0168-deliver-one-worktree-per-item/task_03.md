---
task: task_03
spec: 0168-deliver-one-worktree-per-item
status: pending
type: backend
complexity: medium
---

# Task 03: Park, resume and merge without the restore machinery

## Overview

With each item in its own worktree, resume must find or rebuild that worktree, merge must remove it, and the starting-branch, untracked-baseline and park-cleanup machinery Spec 0161 added to protect the user's checkout has nothing left to protect. Today a deleted starting branch makes park fail forever, so every resume replays the item's stage. This changes CLI-visible behaviour; Task 04 updates the shipped skill and the user guide.

## Requirements

1. MUST, on resume, use the recorded worktree when Git lists it with the recorded branch checked out.
2. MUST recreate a missing recorded worktree at its recorded path from the recorded branch.
3. MUST park the item with the blocker `item-worktree-missing` when the recorded worktree and the recorded branch are both gone, instead of replaying its stage.
4. MUST, after the merged stage is persisted, remove the item worktree and delete the local item branch; a failed removal leaves the item merged and is retried on the next resume.
5. MUST remove the starting branch from the item's reads and writes, `ParkItem`'s reset, clean and switch with its untracked-path tracking, and the clean-checkout checks in branch creation and reuse.
6. MUST replace `TestAParkLeavesACleanCheckout`, `TestParkKeepsUntrackedFilesTheItemDidNotCreate`, `TestParkNeverTouchesACheckoutTheItemRefused` and `TestParkRestoresTheRecordedStartingBranchAfterACrash`, which pin the retired behaviour, and keep the branch and archive guarantees named in the Verification passing.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, a resume uses a recorded worktree that exists.
- [ ] With real Git, a resume recreates a removed worktree from the recorded branch.
- [ ] With real Git, a resume parks with `item-worktree-missing` when the branch is gone too, and a later resume does not replay the stage.
- [ ] With real Git, a merged item leaves no worktree and no local item branch.
- [ ] The kept branch and archive guarantees still pass.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/store/delivery.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestResumeUsesTheRecordedItemWorktree|TestResumeRecreatesAMissingItemWorktree|TestResumeParksWhenTheItemBranchIsGone|TestAMergedItemLeavesNoWorktreeOrBranch|TestItemBranchHasNoUpstream|TestEachDeliveryGetsItsOwnBranch|TestResumeReusesTheRecordedItemBranch|TestResumeAcceptsARealArchiveCommit|TestResumeRefusesAnArchiveCommitWithExtraChanges)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestResumeUsesTheRecordedItemWorktree TestResumeRecreatesAMissingItemWorktree TestResumeParksWhenTheItemBranchIsGone TestAMergedItemLeavesNoWorktreeOrBranch TestItemBranchHasNoUpstream TestEachDeliveryGetsItsOwnBranch TestResumeReusesTheRecordedItemBranch TestResumeAcceptsARealArchiveCommit TestResumeRefusesAnArchiveCommitWithExtraChanges; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the first four named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Park, resume and merge; Retired
