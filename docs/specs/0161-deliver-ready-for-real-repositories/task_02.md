---
task: task_02
spec: 0161-deliver-ready-for-real-repositories
status: pending
type: backend
complexity: medium
---

# Task 02: Untracked, per-delivery item branches

## Overview

The item branch is created tracking `origin/<default>`, so with `implement.auto_push: true` a Clean Run pushes to the default branch; its name is fixed per slug and created with `git switch -c`, so re-delivery and resume after a crash park with "branch already exists".

## Requirements

1. MUST create the item branch with no upstream.
2. MUST give the branch a per-delivery suffix and record it on the item before any Git command creates it.
3. MUST reuse the recorded branch on resume when it already exists.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, the item branch has no upstream.
- [ ] Two deliveries of one slug get different branches.
- [ ] A resume after a crash between creation and the stage write reuses the branch.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/store/delivery.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestItemBranchHasNoUpstream|TestEachDeliveryGetsItsOwnBranch|TestResumeReusesTheRecordedItemBranch)$" ./internal/cli ./internal/store 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestItemBranchHasNoUpstream TestEachDeliveryGetsItsOwnBranch TestResumeReusesTheRecordedItemBranch; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Branches
