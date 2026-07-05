---
task: task_05
spec: 0008-worktree-isolation
status: pending
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

- [ ] WorkDir on the view model with fallback semantics
- [ ] Live and attach wiring
- [ ] Header/path rendering
- [ ] Cockpit and attach tests over worktree and fallback fixtures

## Acceptance Criteria

- [ ] A cockpit over a spec Run reads statuses that exist only in the
      worktree (the user-root copies stale) and renders the worktree
      truth.
- [ ] Attach on a kept-worktree Run reads it; attach on a pruned Clean Run
      falls back to the user root without error.
- [ ] Review rendering snapshots unchanged.
- [ ] Full suite passes.

## Verification

- `rtk go test ./internal/tui/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 6; Core Feature 5. `_techspec.md` → Run flows
(readers), Build Order 5. ADR-0009 (journal discipline unchanged),
ADR-0023.
