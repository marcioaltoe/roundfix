---
task: task_07
spec: 0168-deliver-one-worktree-per-item
status: pending
type: backend
complexity: low
---

# Task 07: A half-removed item worktree never blocks the queue

## Overview

Corrective Task from the second pre-PR review of 2026-09-25. When `git worktree remove --force` fails part-way (for example on a read-only directory a bootstrap wrote inside the item worktree), the registration is gone but the directory remains; every later cleanup refuses with "unregistered path still exists", and `engine.Run` stops before the next queued item.

## Requirements

1. MUST remove, in the no-registration case, a leftover item path whose `.git` file points at this repository's `worktrees/` admin directory, making it writable first.
2. MUST report a failed cleanup of a merged item as a warning on that item and continue with the next queued item, instead of stopping the queue.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, a half-removed item worktree containing a read-only directory is cleaned up on the next attempt.
- [ ] A merged item whose cleanup fails does not stop the next queued item.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`
- interface: `internal/delivery/engine.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCleanupRecoversAHalfRemovedItemWorktree|TestAFailedMergedCleanupDoesNotStopTheQueue)$" ./internal/worktree ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCleanupRecoversAHalfRemovedItemWorktree TestAFailedMergedCleanupDoesNotStopTheQueue; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task neither case exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
