---
task: task_02
spec: 0016-worktree-bootstrap
status: completed
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

- [x] Call `runBootstrap` in Task Worktree creation after CopyList placement
- [x] Map a Task Worktree bootstrap failure to that Task settling failed
- [x] Confirm independent Tasks are unaffected by one Task's bootstrap failure
- [x] Tests: concurrent bootstrap success; one Task's bootstrap failure isolates

## Acceptance Criteria

- [x] With `worktree.bootstrap` set and concurrency > 1, each Task Worktree runs the bootstrap before its Agent work.
- [x] A Task Worktree bootstrap failure settles only that Task failed with the bootstrap-failed reason; independent Tasks still run.
- [x] With an empty `worktree.bootstrap`, concurrent Run behavior is unchanged.
- [x] The Task Worktree bootstrap reuses the task_01 primitive (no separate implementation).

## Verification

- `rtk go test ./internal/worktree/ ./internal/daemon/ ./internal/cli/` — expected: Task Worktree bootstrap tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 4; Core Feature 3. `_techspec.md` → Where bootstrap runs,
Failure handling, Build Order 2. ADR-0034.

## Result

Implemented Task Worktree bootstrap for concurrent spec Runs. Task Worktree
creation now accepts the same bootstrap spec used by Run Worktrees and runs it
after `worktree.copy` placement before the Task Agent starts. The daemon maps a
typed Task Worktree bootstrap failure to that Task settling `failed` with the
`worktree bootstrap failed: <command>: <reason>` message while preserving
independent Task execution.

Acceptance evidence:

- Concurrent bootstrap before Agent work:
  `TestCreateTaskRunsBootstrapAfterCopyInTaskWorktreeRoot` proves the shared
  worktree primitive runs in a Task Worktree root after copied files are
  present. `TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork`
  proves the CLI config path runs real Task Worktree bootstraps before each
  Task Agent and before each Task Verification command when `concurrency: 2`.
- Failure isolation:
  `TestTaskCycleTaskWorktreeBootstrapFailureIsolatesIndependentTasks` proves one
  bootstrap failure settles only that Task failed, journals the bootstrap-failed
  reason, skips that Task's Agent work and integration, and lets independent
  Tasks complete.
- Empty-bootstrap compatibility: the legacy `CreateTask` wrapper delegates with
  no bootstrap options, and the existing concurrent scheduler, daemon, CLI, and
  full-suite tests pass unchanged.
- Primitive reuse: Task Worktrees call `CreateTaskWithOptions`, which invokes
  the same `runBootstrap` primitive introduced for Run Worktrees; no second
  bootstrap implementation was added.

Verification:

- `rtk go test ./internal/worktree/ ./internal/daemon/ ./internal/cli/` passed
  (`372` tests in `3` packages).
- `rtk go test ./...` passed (`813` tests in `18` packages).
- `rtk make verify` passed: `go test ./...`, `roundfix skills check`, and
  `go build -buildvcs=false -o bin/roundfix ./cmd/roundfix`.
