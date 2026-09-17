---
task: task_05
spec: 0141-a-commit-that-carries-the-work-and-nothing-else
status: pending
type: test
complexity: low
---

# Task 05: Restore the permission-variant coverage the rule kept true

## Overview

Task 01 replaced the test that pinned the refusal of a regular file carrying any
execute permission. The behavior it pinned still holds for an untracked file, so
the repository's coverage record reads the replacement as a lost test. This
slice restores it beside the new cases.

## Requirements

1. MUST restore a test named
   `TestFilterStageablePathsDropsRegularFileWithAnyExecutePermission` that
   asserts the refusal for owner, group and other execute permissions.
2. MUST keep every test Task 01 added.
3. MUST NOT edit the repository's coverage record, which is a Governed Path this
   Spec has no authority over.
4. MUST NOT weaken an assertion to make a test pass.

## Subtasks

- [ ] Restore the permission-variant test under its recorded name.
- [ ] Confirm the new tracked-path tests still pass beside it.
- [ ] Run the coverage equivalence check.

## Acceptance Criteria

- [ ] The restored test passes with the tracked-path rule in place.
- [ ] The coverage equivalence check reports no regression.
- [ ] No coverage record file changes.

## Context

- interface: `internal/daemon/task_engine_test.go`

## Verification

- `grep -q "func TestFilterStageablePathsDropsRegularFileWithAnyExecutePermission" internal/daemon/task_engine_test.go` — expected: exit 0; the recorded test exists again. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestFilterStageablePathsDropsRegularFileWithAnyExecutePermission$" ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `go test -count=1 -run "^TestCoverageEquivalence$" ./internal/spec` — expected: exit 0; no coverage regression remains. Before this Task it fails on the missing test.

## References

`_prd.md` → Core Feature 1; Regression locks;
`_techspec.md` → Implementation Design: A tracked path is stageable;
Testing Approach 1; ADR-0157.
