---
task: task_02
spec: 0168-deliver-one-worktree-per-item
status: completed
type: backend
complexity: high
---

# Task 02: Run every item in its own worktree

## Overview

`roundfix deliver` switches the user's checkout to each item branch, runs every per-item action there, and refuses to start an item on a dirty checkout. Each item instead gets its own linked worktree, and nothing the queue does reaches the user's checkout. This changes CLI-visible behaviour; Task 04 updates the shipped skill and the user guide.

## Requirements

1. MUST create, per item, a linked worktree at a path derived under `worktree.location` from the repository and the item branch, recorded on the item with the branch before any Git command creates either.
2. MUST create the worktree from the refreshed default branch, with the item branch created without an upstream, and provision it with the configured `worktree.copy` files from the user's checkout and the configured bootstrap, reusing `internal/worktree`.
3. MUST run `roundfix implement`, pre-PR review, archive, the repository gate, the authorization read, the publication plan and push with the item worktree as their working directory, keeping queue and action records keyed by the user's repository.
4. MUST NOT switch, reset, clean or inspect for cleanliness the user's checkout; an item starts while that checkout is dirty or on any branch.
5. MUST park an item without changing any checkout, leaving its worktree in place.
6. MUST print the item worktree as a fourth tab-separated field of `deliver status`, or `-` when none is recorded.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, each item runs in its own worktree created from the refreshed default branch, and its branch has no upstream.
- [ ] With real Git, a user checkout that is dirty, on a non-default branch, with an ignored file and an untracked nested repository, is unchanged in branch, HEAD, status and file contents after one item parks and another merges.
- [ ] A parked item's worktree remains at its recorded path.
- [ ] Every per-item action receives the item worktree as its working directory.
- [ ] `deliver status` prints each item's worktree.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/github.go`
- interface: `internal/worktree/worktree.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestEachItemRunsInItsOwnWorktree|TestDeliverNeverTouchesTheUserCheckout|TestParkLeavesTheItemWorktreeInPlace|TestDeliverStatusPrintsTheItemWorktree|TestDeliveryActionsRunInTheItemWorktree)$" ./internal/cli ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEachItemRunsInItsOwnWorktree TestDeliverNeverTouchesTheUserCheckout TestParkLeavesTheItemWorktreeInPlace TestDeliverStatusPrintsTheItemWorktree TestDeliveryActionsRunInTheItemWorktree; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The item worktree

## Result

Implemented the item-worktree slice without changing the Daemon-owned status:

- The workflow derives and records each item's branch and worktree before Git creates either, refreshes the default branch, creates an untracked branch in a linked worktree, and provisions configured copy files and bootstrap through `internal/worktree`.
- The Delivery Engine keeps queue and action records keyed by the user's repository while running implement, review, archive, gate, authorization, publication, push, pull-request, check, and merge operations from the recorded item worktree.
- Parking now only persists the blocker and leaves every checkout and the item worktree unchanged.
- `deliver status` prints the recorded item worktree as its fourth tab-separated field and prints `-` when none is recorded.

Focused check:

- `GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 ./internal/delivery ./internal/cli ./internal/worktree` — passed.

Acceptance evidence:

- `TestEachItemRunsInItsOwnWorktree` uses real Git to prove the worktree starts at the refreshed default-branch head, has no upstream, lives under `worktree.location`, receives copied files and bootstrap output, and is recorded with its branch.
- `TestDeliverNeverTouchesTheUserCheckout` uses real Git to run one parked item and one merged item while preserving a dirty non-default user checkout's branch, HEAD, porcelain status, tracked content, ignored content, and nested repository content; it also proves the two items record distinct worktrees.
- `TestParkLeavesTheItemWorktreeInPlace` proves a parked item's recorded directory still exists and its branch remains registered as a linked worktree.
- `TestDeliveryActionsRunInTheItemWorktree` proves every per-item workflow and GitHub boundary receives the item worktree while the Delivery Queue remains keyed only by the user's repository.
- `TestDeliverStatusPrintsTheItemWorktree` proves both the recorded-path and empty-path (`-`) status fields.

Follow-up: Task 03 remains responsible for recreating missing item worktrees on resume and removing merged item worktrees and local branches.
