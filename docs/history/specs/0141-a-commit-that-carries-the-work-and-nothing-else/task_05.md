---
task: task_05
spec: 0141-a-commit-that-carries-the-work-and-nothing-else
status: completed
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

- [x] Restore the permission-variant test under its recorded name.
- [x] Confirm the new tracked-path tests still pass beside it.
- [ ] Run the coverage equivalence check.

## Acceptance Criteria

- [x] The restored test passes with the tracked-path rule in place.
- [ ] The coverage equivalence check reports no regression.
- [x] No coverage record file changes.

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

## Result

Implementation:

- Restored `TestFilterStageablePathsDropsRegularFileWithAnyExecutePermission`
  as a named, table-driven test over owner, group and other execute bits on an
  untracked regular file in a temporary Git repository.
- Preserved the Task 01 tracked-path tests. The restored assertions require the
  path to be omitted, the `executable file` reason, the lost-output class and
  the exact permission mode for every variant.

Focused checks:

- `rtk env GOCACHE=/private/tmp/roundfix-task-0141-05-go-cache go test -count=1 -run 'TestFilterStageablePaths(KeepsTrackedExecutable|KeepsTrackedRegularFile|RefusesUntrackedExecutableWithMode|DropsRegularFileWithAnyExecutePermission|RefusesExecutableWhenIndexQueryFails)$' ./internal/daemon` — exit 0; the restored test and Task 01's tracked, untracked and failed-index cases pass together.
- The first attempt used the default Go cache and stopped before compilation
  because the sandbox denied access to `/Users/marcio/Library/Caches/go-build`;
  the command above reran the same selection with a writable cache.
- The authored coverage-equivalence Verification was intentionally left for
  the Daemon. No coverage record was edited by this Task.
- `rtk git diff --name-only` — listed only this Task file and
  `internal/daemon/task_engine_test.go`; `docs/references/coverage-record.json`
  is unchanged.
- `rtk git diff --check` — exit 0.
