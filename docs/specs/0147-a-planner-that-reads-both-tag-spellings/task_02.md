---
task: task_02
spec: 0147-a-planner-that-reads-both-tag-spellings
status: completed
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

## Result

Implemented semantic highest-version selection over parsed `VersionRef`
entries. Selection returns the original highest ref, keeps every input entry
distinct, and returns `AmbiguousHighestVersionError` with all equal-highest refs
when more than one spelling reaches the same version. The error unwraps to the
named `ErrAmbiguousHighestVersion` condition.

Focused evidence:

- Red signal: `rtk env GOCACHE=/tmp/roundfix-task02-go-cache go test -count=1 -run 'TestSelectionSpansSpellings/higher_bare_version' ./internal/releaseplan` exited 1 before implementation because `SelectHighestVersion`, `ErrAmbiguousHighestVersion`, and `AmbiguousHighestVersionError` did not exist.
- Mixed-spelling selection: `rtk env GOCACHE=/tmp/roundfix-task02-go-cache go test -count=1 -run 'TestSelectionSpansSpellings/' ./internal/releaseplan` exited 0; the `v1.2.3` and `1.3.0` case selected the original `1.3.0` ref.
- Ambiguity and inventory identity: the same focused command exited 0; the `v2.0.0` and `2.0.0` case matched `ErrAmbiguousHighestVersion` and retained both tags and their original commit SHAs in the typed error.
- Single-spelling regression: the same focused command exited 0; the prefixed-only `v1.2.3`, `v1.10.0`, and `v1.9.9` case selected `v1.10.0`.
- Parser compatibility: `rtk env GOCACHE=/tmp/roundfix-task02-go-cache go test -count=1 -run 'TestParseStableVersion|TestStableVersionAcceptsBothSpellings' ./internal/releaseplan` exited 0.
- Package-level focused check: `rtk env GOCACHE=/tmp/roundfix-task02-go-cache go test -count=1 ./internal/releaseplan` exited 1 in the pre-existing `malformed_target_version` reset-inventory case. That prior-task test still treats bare `0.0.1` as malformed even though Task 01 intentionally made the bare spelling valid; it is outside Task 02's selection slice.

Follow-up: update the stale reset-inventory malformed-target fixture in its
own owning slice so the package-level check reflects the accepted bare spelling.
