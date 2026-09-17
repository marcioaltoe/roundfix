---
task: task_01
spec: 0141-a-commit-that-carries-the-work-and-nothing-else
status: completed
type: backend
complexity: medium
---

# Task 01: Ask the index before refusing a file for its mode

## Overview

The stageable filter refuses every executable regular file, so a Task that edits
a tracked hook or script loses the change at the commit boundary. This slice
makes the filter ask whether Git already tracks the path, and keep it when it
does. An untracked executable is refused exactly as today.

## Requirements

1. MUST keep a path the index already tracks, whatever its file mode.
2. MUST refuse an untracked executable regular file with the reason and mode it
   reports today.
3. MUST treat an index query that cannot be answered as untracked, preserving
   today's refusal rather than staging on an unknown.
4. MUST leave the external, symbolic-link and absent refusals unchanged in
   order, reason text and payload.
5. MUST name the index query `pathTrackedInIndex`, so the Task's own check and
   later readers agree.

## Subtasks

- [ ] Add the index query beside the existing absence query.
- [ ] Ask it before the mode refusal and keep a tracked path.
- [ ] Cover tracked-executable, untracked-executable, tracked-regular and
      unanswerable-query cases at the filter's unit seam.

## Acceptance Criteria

- [ ] A tracked executable file is kept by the filter.
- [ ] An untracked executable file is refused with its mode and existing reason.
- [ ] A tracked regular file is kept, as today.
- [ ] An index query that fails refuses the path.
- [ ] The other three refusals keep their order and their text.

## Context

- interface: `internal/daemon/task_engine.go`

## Verification

- `grep -q "pathTrackedInIndex" internal/daemon/task_engine.go` — expected: exit 0; the filter asks the index. Before this Task no such query exists.
- `out="$(go test -count=1 -run "^TestFilterStageablePathsKeepsTrackedExecutable$" ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Feature 1; User Story 1; Goal 1; Success Metric 1;
`_techspec.md` → Implementation Design: A tracked path is stageable;
API Contract 2; Build Order 1; ADR-0157.

## Result

Implementation:

- Added `pathTrackedInIndex` beside the existing absence query. It asks Git's
  index with `ls-files --error-unmatch` and returns false for every query
  error, so an unknown answer never permits staging.
- `FilterStageablePaths` now keeps an index-tracked path after the external and
  symbolic-link checks and before inspecting its file mode. The executable,
  absent, symbolic-link and external refusal values remain unchanged.
- Extended the existing filter unit seam with tracked executable, tracked
  regular, untracked executable, failed index query, and ordered legacy-refusal
  cases.

Focused checks:

- `rtk env GOCACHE=/tmp/roundfix-task-0141-go-build go test -count=1 -run 'TestFilterStageablePaths' ./internal/daemon` failed before the production change because
  `TestFilterStageablePathsKeepsTrackedExecutable` received the existing
  `executable file` drop with mode `0755`.
- The same focused command passed after the production change.
- `rtk env GOCACHE=/tmp/roundfix-task-0141-go-build go test -count=1 -v -run 'TestFilterStageablePaths' ./internal/daemon` passed all filter cases and named each case in its output.
- `rtk env GOCACHE=/tmp/roundfix-task-0141-go-build go test -count=1 ./internal/daemon` passed the complete daemon package.
- `rtk git diff --check` passed.
- The Task's authored `## Verification` commands were not run; the Daemon owns
  those checks.

Acceptance evidence:

- A tracked executable is kept: `TestFilterStageablePathsKeepsTrackedExecutable` passed against a real disposable Git index.
- An untracked executable retains its reason and mode:
  `TestFilterStageablePathsRefusesUntrackedExecutableWithMode` passed for owner,
  group and other execute bits, asserting reason `executable file` and the
  exact permission mode.
- A tracked regular file remains kept:
  `TestFilterStageablePathsKeepsTrackedRegularFile` passed.
- An index query failure refuses the executable:
  `TestFilterStageablePathsRefusesExecutableWhenIndexQueryFails` passed in a
  non-repository work directory, asserting the existing reason and mode.
- The remaining refusals retain their order and text:
  `TestFilterStageablePathsPreservesOtherRefusals` passed with `external to
  repository`, `crosses a symbolic link`, and `absent from worktree and index`
  in input order and with unchanged payload fields.

## Carry-forward provenance

- Source Run: `run_20260917T192120Z_0e305cd0796d7680`
- Source commit: `c63a2a06767031ffb9d5211f58cb21f9c118014a`
