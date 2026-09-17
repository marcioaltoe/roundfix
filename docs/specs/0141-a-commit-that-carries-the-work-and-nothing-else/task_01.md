---
task: task_01
spec: 0141-a-commit-that-carries-the-work-and-nothing-else
status: pending
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
