---
task: task_05
spec: 0008-worktree-isolation
status: completed
type: frontend
complexity: medium
---

# Task 05: Point the Live Run View and Attach at the execution workspace

## Overview

The cockpit, plain renderer, attach, and end-of-run reports must read Task
state from where execution actually happens: the Run Worktree, with a clean
fallback to the user root for terminal Runs whose worktree is gone.
Verifiable through synchronous cockpit tests and attach fixtures over both
locations.

## Requirements

1. MUST give the Live Run View a `WorkDir` distinct from `GitRoot`: task
   status refresh, detail-modal task reads, and the implement report lines
   read from `WorkDir` when set, falling back to `GitRoot` when empty or
   missing on disk (mid-write tolerance unchanged).
2. MUST wire it live (implement/resolve set the worktree path) and in
   Attach (the Run row's `work_dir`, fallback for legacy and pruned Runs).
3. MUST show the Run Worktree path in the Run header/target block for
   worktree-backed Runs.
4. MUST keep review-Run rendering and all cockpit key semantics
   byte-stable.

## Subtasks

- [x] WorkDir on the view model with fallback semantics
- [x] Live and attach wiring
- [x] Header/path rendering
- [x] Cockpit and attach tests over worktree and fallback fixtures

## Acceptance Criteria

- [x] A cockpit over a spec Run reads statuses that exist only in the
      worktree (the user-root copies stale) and renders the worktree
      truth.
- [x] Attach on a kept-worktree Run reads it; attach on a pruned Clean Run
      falls back to the user root without error.
- [x] Review rendering snapshots unchanged.
- [x] Full suite passes.

## Verification

- `rtk go test ./internal/tui/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 6; Core Feature 5. `_techspec.md` → Run flows
(readers), Build Order 5. ADR-0009 (journal discipline unchanged),
ADR-0023.

## Result

Implemented the Live Run View `WorkDir` reader path and wired it through
implement, resolve/watch live views, and Attach. Task status refresh and Task
detail reads now prefer an existing Run Worktree and fall back to the user
checkout when `work_dir` is empty or pruned. Worktree-backed views render the
stored Run Worktree path.

Evidence:

- Cockpit WorkDir truth: `TestCockpitSpecRunReadsTaskStatusAndDetailFromWorkDir`
  passes; it keeps stale user-root Task files pending while rendering the
  completed status and detail body that exist only in the worktree.
- Cockpit fallback: `TestCockpitSpecRunFallsBackToGitRootWhenWorkDirIsGone`
  passes; a missing WorkDir falls back to the user root without surfacing an
  error.
- Attach WorkDir and pruned fallback:
  `TestAttachSpecRunReadsTasksFromKeptWorkDir` and
  `TestAttachSpecRunFallsBackToGitRootWhenCleanWorkDirIsPruned` pass.
- Review rendering snapshots stayed byte-stable under
  `TestCockpitRenderSnapshots`; the new path only appears when `WorkDir` is
  set.

Verification:

- `rtk go test ./internal/tui/ ./internal/cli/` passed: 342 tests in 2
  packages.
- `rtk go test ./...` passed: 682 tests in 17 packages.
