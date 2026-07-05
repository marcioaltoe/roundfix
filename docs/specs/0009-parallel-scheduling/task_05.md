---
task: task_05
spec: 0009-parallel-scheduling
status: pending
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

- [ ] Settle surface resolution order with deterministic Task paths
- [ ] Settle-time integration with conflict reporting
- [ ] Reap wiring in preflight and stop --force with reporting
- [ ] Kept-worktree fixtures for all three resolution cases

## Acceptance Criteria

- [ ] A failed concurrent Task's kept worktree settles: verification
      streams, status completes, commit lands on the Run Branch, worktree
      and branch cleaned; the Run-level state matches the shipped settle
      semantics.
- [ ] The sequential fallback (Run Worktree) and legacy fallback still pass
      the existing settle tests.
- [ ] A settle-time conflict exits 1 naming paths, everything kept.
- [ ] Force-stopping a Run with zero settled work reaps its worktrees and
      branches with the report lines; a Run with commits keeps them.
- [ ] Full suite passes.

## Verification

- `rtk go test ./internal/cli/ ./internal/worktree/` — expected: all tests
  pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 3, 7; Core Features 2, 5. `_techspec.md` → Settle
and debris, Build Order 5. ADR-0026. Round-3 finding 1.
