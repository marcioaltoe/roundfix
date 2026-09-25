---
task: task_08
spec: 0168-deliver-one-worktree-per-item
status: pending
type: backend
complexity: low
---

# Task 08: Remove only a truly half-removed item worktree

## Overview

Third pre-PR review of 2026-09-25, data-loss finding: `CleanupItem` only asks whether the item branch is checked out anywhere, so an item worktree that is still registered but on another branch, or detached, takes the "unregistered" route and `removeUnregisteredItemWorktree` deletes it with `os.RemoveAll`, uncommitted files included, leaving a dangling registration. Every `deliver resume` retries cleanup of merged items, so this would erase a user's work in that worktree silently.

## Requirements

1. MUST remove a leftover item path only when no registered worktree matches it by canonical path and the admin directory named in its `.git` file no longer exists.
2. MUST refuse removal, keeping the path and reporting the merged item's cleanup as a warning, when the worktree is still registered, on any branch or detached.
3. MUST only remove a path equal to the one `ItemRefFor` derives from the configured location for that item.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, a registered item worktree switched to another branch keeps its uncommitted file through cleanup.
- [ ] A detached item worktree is kept.
- [ ] A truly half-removed item worktree is still recovered.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCleanupKeepsARegisteredItemWorktreeOnAnotherBranch|TestCleanupKeepsADetachedItemWorktree|TestCleanupRecoversAHalfRemovedItemWorktree)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCleanupKeepsARegisteredItemWorktreeOnAnotherBranch TestCleanupKeepsADetachedItemWorktree TestCleanupRecoversAHalfRemovedItemWorktree; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the first two cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
