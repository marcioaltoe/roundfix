---
task: task_06
spec: 0141-a-commit-that-carries-the-work-and-nothing-else
status: pending
type: backend
complexity: medium
---

# Task 06: Close the two gaps the review found

## Overview

Independent review found two ways the new rules read a path more coarsely than
they promise. The index question matches a directory prefix, so an untracked
executable that replaced a tracked directory is treated as tracked. And the
worktree snapshot collapses a new untracked directory to the directory itself,
so excluding the Verification's window can discard a file the QA Agent wrote
into the same new directory. This slice closes both, and it is the second and
final corrective Task the contract allows.

## Requirements

1. MUST treat a path as tracked only when the index holds that exact path, never
   when the path merely names a directory whose contents are tracked.
2. MUST keep refusing an untracked executable that replaced a tracked directory,
   with the reason and mode the filter reports today.
3. MUST exclude from the QA Report commit only what the repository Verification
   actually wrote, so a file the QA Agent writes into a directory the
   Verification created still reaches the commit.
4. MUST keep the report the mechanical stage seeded and its evidence in the
   commit, and leave a Spec whose Verification writes nothing unchanged.
5. MUST NOT change the shared snapshot contract for other callers unless the
   change is proved by their own tests.

## Subtasks

- [ ] Make the index question exact.
- [ ] Cover the directory-replacement case at the filter's unit seam.
- [ ] Make the Verification window compare files rather than collapsed
      directories.
- [ ] Cover an Agent file written into a directory the Verification created.

## Acceptance Criteria

- [ ] An untracked executable that replaced a tracked directory is refused with
      its mode.
- [ ] A tracked executable is still kept.
- [ ] A QA Report commit keeps an Agent file written into a directory the
      repository Verification created, and still excludes the Verification's own
      file in that directory.
- [ ] A Spec whose Verification writes nothing produces the same commit as
      before this Task.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/daemon.go`

## Verification

- `out="$(go test -count=1 -run "^TestFilterStageablePathsRefusesExecutableReplacingATrackedDirectory$" ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `out="$(go test -count=1 -run "^TestQAReportCommitKeepsAgentEvidenceInANewDirectory$" ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `grep -q "TestFilterStageablePathsRefusesExecutableReplacingATrackedDirectory" internal/daemon/task_engine_test.go && go test -count=1 -run "^TestFilterStageablePaths" ./internal/daemon` — expected: exit 0; every filter case, including the tracked one, holds after the question becomes exact. Before this Task the directory-replacement case does not exist, so the command fails.

## References

`_prd.md` → Core Features 1 and 3; Goals 1 and 3; Success Metrics 1 and 3;
Regression locks;
`_techspec.md` → Implementation Design: A tracked path is stageable, The QA
Report commit excludes the Verification's writes; API Contracts 2-3; ADR-0157.
