---
task: task_05
spec: 0009-parallel-scheduling
status: completed
type: backend
complexity: medium
---

# Task 05: Settle over Task Worktrees and the debris-reap wiring

## Overview

Recovery and hygiene under parallelism: `roundfix settle` finds a failed
Task's kept Task Worktree through the deterministic path (falling back to
the Run Worktree for sequential-mode failures), verifies and settles there,
and hands the commit to the same integration mechanics; the empty-debris
reap wires into the preflight sweep and `stop --force`. Verifiable through
buffer-captured CLI tests over kept-worktree fixtures.

## Requirements

1. MUST make settle resolve the target surface in order: the Task's kept
   Task Worktree at the deterministic path, else the Run's kept Run
   Worktree via `work_dir`, else the current repository (legacy behavior) —
   verifying, settling, and committing in whichever it found, then
   integrating onto the Run Branch through the queue mechanics; a
   settle-time integration conflict reports, keeps the worktree, and exits
   1 with the conflict named.
2. MUST, on successful settle-and-integrate of a Task Worktree, clean up
   that worktree and Task Branch; the Run-level integration to the user
   branch keeps the shipped settle behavior.
3. MUST wire the extended empty-debris reap (task_02) into the implement
   preflight sweep and into `stop --force`, with one stderr line naming
   each reaped path/branch.
4. MUST keep every existing settle contract line and exit code, extending
   help text only for the new resolution order.

## Subtasks

- [x] Settle surface resolution order with deterministic Task paths
- [x] Settle-time integration with conflict reporting
- [x] Reap wiring in preflight and stop --force with reporting
- [x] Kept-worktree fixtures for all three resolution cases

## Acceptance Criteria

- [x] A failed concurrent Task's kept worktree settles: verification
      streams, status completes, commit lands on the Run Branch, worktree
      and branch cleaned; the Run-level state matches the shipped settle
      semantics.
- [x] The sequential fallback (Run Worktree) and legacy fallback still pass
      the existing settle tests.
- [x] A settle-time conflict exits 1 naming paths, everything kept.
- [x] Force-stopping a Run with zero settled work reaps its worktrees and
      branches with the report lines; a Run with commits keeps them.
- [x] Full suite passes.

## Verification

- `rtk go test ./internal/cli/ ./internal/worktree/` — expected: all tests
  pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 3, 7; Core Features 2, 5. `_techspec.md` → Settle
and debris, Build Order 5. ADR-0026. Round-3 finding 1.

## Result

- Settle now resolves surfaces in order: deterministic Task Worktree, kept Run
  Worktree, then current repository. `TestRunSettleRetargetsKeptTaskWorktreeAndCleansUpAfterIntegration`
  covers the concurrent Task path; existing settle tests continue to cover the
  Run Worktree and current-repository fallbacks.
- Task Worktree settle uses the Task integration queue. On success it cleans
  the Task Worktree and Task Branch before the existing Run-level integration;
  on conflict it exits 1, names the conflicting path, and keeps the Run and
  Task surfaces.
- Empty-debris cleanup now reports reaped path and branch lines from the
  implement preflight sweep and from `stop --force`.

Acceptance evidence:

- Concurrent Task settle: `TestRunSettleRetargetsKeptTaskWorktreeAndCleansUpAfterIntegration`
  verifies streamed Verification output, completed task status, integrated user
  checkout content, Clean Run state, and cleanup of the Run/Task worktrees and
  branches.
- Fallbacks: `TestRunSettleRetargetsKeptRunWorktreeAndCleansUpAfterIntegration`
  and `TestRunSettleCommitsFailedTaskWorktreeWithDaemonMessage` still pass.
- Conflict: `TestRunSettleTaskWorktreeIntegrationConflictKeepsSurfaces` verifies
  exit 1, stderr naming `shared.txt`, unchanged Run Branch tip, and kept Run
  Worktree, Run Branch, Task Worktree, and Task Branch.
- Debris reap: `TestRunImplementPreflightReapsEmptyTerminalRunAndTaskWorktrees`
  verifies implement preflight reporting and cleanup; `TestRunStopForceReapsEmptyRunAndTaskWorktrees`
  verifies `stop --force` cleanup/reporting; `TestRunStopForceKeepsRunWorktreeWithCommits`
  verifies a Run Branch with commits is kept.

Verification:

- `rtk go test ./internal/cli/ ./internal/worktree/` passed with 262 tests.
- `rtk go test ./...` passed with 711 tests across 17 packages.
- `rtk make verify` passed: full Go tests, Roundfix skill check, and build.
