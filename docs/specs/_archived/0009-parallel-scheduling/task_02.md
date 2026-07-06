---
task: task_02
spec: 0009-parallel-scheduling
status: completed
type: backend
complexity: high
---

# Task 02: Task Worktree lifecycle and the serialized integration queue

## Overview

Extend `internal/worktree` with the per-Task layer: Task Worktrees created
from the Run Branch tip as siblings of the Run Worktree, the ADR-0026
integration (fast-forward when possible, cherry-pick otherwise, conflict →
reported with the Run Branch unmoved), success cleanup, and the
empty-debris reap for provably valueless kept worktrees. No scheduler
wiring yet. Verifiable alone over hermetic temp repositories.

## Requirements

1. MUST create Task Worktrees at the sibling path
   `<location>/<repo-slug>/<run-id>.<task_id>/` on branch
   `roundfix/run-<run-id>-<task_id>` from the Run Branch tip, with the
   copy-list provisioning applied; never nested inside another worktree.
2. MUST implement `IntegrateTask` for a single caller-serialized queue:
   fast-forward the Run Branch when the task branch's base equals the
   current tip; otherwise cherry-pick the Task's commit(s) onto the Run
   Branch preserving message and trailers byte-for-byte; on conflict,
   `cherry-pick --abort`, leave the Run Branch at its pre-attempt tip, and
   return a conflict result naming the conflicting paths.
3. MUST implement success cleanup (remove Task Worktree and Task Branch
   after integration) while failed Tasks keep both as inspection surfaces.
4. MUST extend the terminal-Run reap (used by the preflight sweep and
   `stop --force`): kept worktrees and branches — Run or Task — of terminal
   Runs whose branch has no commits beyond its base are removed;
   anything with commits stays untouched.
5. MUST keep every operation on the package's context-first runner with
   hermetic test discipline, and never use plumbing ref updates.

## Subtasks

- [x] Task Worktree creation with sibling paths and provisioning
- [x] IntegrateTask: ff, cherry-pick, conflict-abort matrix
- [x] Success cleanup vs kept-on-failure
- [x] Empty-debris reap extension
- [x] Hermetic matrix test suite

## Acceptance Criteria

- [x] Two Task Worktrees from one tip: first integrates by ff, second by
      cherry-pick, both commits on the Run Branch in completion order with
      messages and trailers byte-identical to the originals.
- [x] Induced conflict (two Tasks editing one file): second integration
      returns the conflict result naming the path, the Run Branch tip is
      byte-unmoved, and the conflicting Task Worktree survives.
- [x] Success cleanup removes worktree and branch; failure keeps both.
- [x] Reap matrix: empty terminal Run and Task branches reaped; branches
      with commits kept; non-terminal Runs never touched.
- [x] Full suite passes; no wiring outside the package.

## Verification

- `rtk go test ./internal/worktree/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 3, 7; Core Features 2, 3, 5. `_techspec.md` →
Interfaces, Paths and naming, Build Order 2. ADR-0024 (unchanged outer
protocol), ADR-0026. Round-3 finding 1.

## Result

- Task Worktree creation evidence: `TestTaskWorktreesIntegrateFirstByFastForwardThenCherryPick` creates two Task Worktrees as siblings of the Run Worktree, verifies neither is nested, and verifies copy-list provisioning into the Task Worktree.
- Integration evidence: `TestTaskWorktreesIntegrateFirstByFastForwardThenCherryPick` proves first Task integration uses fast-forward and second uses cherry-pick, then compares raw commit-object messages/trailers byte-for-byte in completion order.
- Conflict evidence: `TestIntegrateTaskReturnsConflictAndLeavesRunBranchUnmoved` induces a conflict on `shared.txt`, verifies the conflict result names the path, verifies the Run Branch tip is unchanged after `cherry-pick --abort`, and verifies the conflicting Task Worktree and branch remain.
- Cleanup evidence: `TestCleanupTaskRemovesTaskWorktreeAndBranch` verifies successful Task cleanup removes both Task Worktree and Task Branch; the conflict test verifies failed/conflicting Tasks keep both inspection surfaces.
- Reap evidence: `TestPruneTerminalReapsOnlyEmptyTerminalRunAndTaskBranches` verifies empty terminal Run and Task worktrees/branches are removed, branches with commits are kept, and non-terminal Runs are untouched.
- Scope evidence: implementation is contained to `internal/worktree` and its tests, with no scheduler or CLI wiring changes.
- Verification: `rtk go test ./internal/worktree/` passed with 15 tests; `rtk go test ./...` passed with 698 tests across 17 packages.
