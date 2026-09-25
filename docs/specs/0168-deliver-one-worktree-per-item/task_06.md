---
task: task_06
spec: 0168-deliver-one-worktree-per-item
status: pending
type: backend
complexity: medium
---

# Task 06: Provision every item worktree, and let migrated merged items rest

## Overview

Corrective Task from the pre-PR review of 2026-09-25. `worktree.copy` and bootstrap run only in `CreateItem`: a worktree recreated on resume, or one whose owner was killed mid-bootstrap, is used without copied files or bootstrap, so the Run loses files like `.env` and the gate can park the item for no real reason. And every `engine.Run` cleans up `merged` items; a merged item migrated from v18 has no recorded worktree, so cleanup fails and blocks the queue.

## Requirements

1. MUST record, on the item, whether its worktree finished provisioning, and run copy and bootstrap whenever a worktree is recreated or reused without completed provisioning.
2. MUST skip worktree removal for a merged item with no recorded worktree, deleting its local branch only when it is not checked out anywhere.
3. MUST keep a resume of an all-merged queue a no-op.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, a recreated item worktree carries the copied files and the bootstrap output.
- [ ] A worktree whose provisioning was interrupted is completed on reuse.
- [ ] A migrated merged item without a worktree does not block the next item.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/cli/deliver_workflow.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRecreatedItemWorktreeIsProvisioned|TestUnfinishedProvisioningIsCompletedOnReuse|TestMigratedMergedItemDoesNotBlockResume)$" ./internal/worktree ./internal/delivery ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRecreatedItemWorktreeIsProvisioned TestUnfinishedProvisioningIsCompletedOnReuse TestMigratedMergedItemDoesNotBlockResume; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
