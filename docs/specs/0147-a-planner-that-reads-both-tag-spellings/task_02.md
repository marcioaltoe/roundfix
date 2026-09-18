---
task: task_02
spec: 0147-a-planner-that-reads-both-tag-spellings
status: pending
type: backend
complexity: medium
---

# Task 02: Select across spellings and report ambiguity

## Overview

With both spellings parsed, selection must compare them as versions rather than
as text, and must notice when the highest version is reachable under more than
one ref instead of quietly returning one.

## Requirements

1. MUST select the highest reachable semantic version across both spellings, so
   `1.3.0` outranks `v1.2.3`.
2. MUST report, rather than resolve, a highest version reachable under more than
   one ref spelling, carrying both refs in the report.
3. MUST leave selection unchanged for a repository whose tags all share one
   spelling.
4. MUST NOT collapse two refs into one entry anywhere it inventories them.
5. MUST name the reported condition `ErrAmbiguousHighestVersion`, so the Task's
   own check and later readers agree.

## Subtasks

- [ ] Order by version across spellings.
- [ ] Detect and report the ambiguous highest version with both refs.
- [ ] Cover mixed-spelling selection, single-spelling selection and the
      ambiguous case.

## Acceptance Criteria

- [ ] A repository with `v1.2.3` and `1.3.0` selects `1.3.0`.
- [ ] A repository with `v2.0.0` and `2.0.0` reports the ambiguity with both
      refs.
- [ ] A single-spelling repository selects exactly what it selects today.

## Context

- interface: `internal/releaseplan/version.go`
- interface: `internal/releaseplan/errors.go`

## Verification

- `grep -rq "ErrAmbiguousHighestVersion" internal/releaseplan/` — expected: exit 0; the ambiguity has a named condition. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestSelectionSpansSpellings$" ./internal/releaseplan 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Features 2 and 5; User Story 3; Goal 3;
`_techspec.md` → Implementation Design: Selecting across spellings;
Build Order 2.
