---
task: task_02
spec: 0016-worktree-bootstrap
status: pending
type: backend
complexity: medium
---

# Task 02: Task Worktree bootstrap for concurrent Runs

## Overview

Extend bootstrap to Task Worktrees so concurrent spec Runs prepare each Task's
worktree too. A Task Worktree bootstrap failure settles only that Task failed,
preserving failure isolation for independent Tasks. Reuses the `runBootstrap`
primitive from task_01.

## Requirements

1. MUST run `worktree.bootstrap` in each newly created Task Worktree (concurrent
   Tasks) after the `worktree.copy` placement and before that Task's Agent work,
   reusing the task_01 primitive (no duplicate implementation).
2. MUST, on a Task Worktree bootstrap failure, settle that Task `failed` with
   reason `worktree bootstrap failed: <command>: <reason>` and MUST NOT affect
   other independent Tasks (failure isolation, per 0009).
3. MUST NOT run bootstrap when `worktree.bootstrap` is empty, keeping concurrent
   Run behavior byte-stable.
4. MUST bound each Task Worktree bootstrap by `worktree.bootstrap_timeout`.

## Subtasks

- [ ] Call `runBootstrap` in Task Worktree creation after CopyList placement
- [ ] Map a Task Worktree bootstrap failure to that Task settling failed
- [ ] Confirm independent Tasks are unaffected by one Task's bootstrap failure
- [ ] Tests: concurrent bootstrap success; one Task's bootstrap failure isolates

## Acceptance Criteria

- [ ] With `worktree.bootstrap` set and concurrency > 1, each Task Worktree runs the bootstrap before its Agent work.
- [ ] A Task Worktree bootstrap failure settles only that Task failed with the bootstrap-failed reason; independent Tasks still run.
- [ ] With an empty `worktree.bootstrap`, concurrent Run behavior is unchanged.
- [ ] The Task Worktree bootstrap reuses the task_01 primitive (no separate implementation).

## Verification

- `rtk go test ./internal/worktree/ ./internal/daemon/ ./internal/cli/` — expected: Task Worktree bootstrap tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 4; Core Feature 3. `_techspec.md` → Where bootstrap runs,
Failure handling, Build Order 2. ADR-0034.
