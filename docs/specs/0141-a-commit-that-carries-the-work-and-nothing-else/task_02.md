---
task: task_02
spec: 0141-a-commit-that-carries-the-work-and-nothing-else
status: completed
type: backend
complexity: medium
---

# Task 02: Make a lost output fail the Task

## Overview

A refused path is a console line and an event today, and the Task settles
completed anyway. This slice separates a refusal by design from a lost output,
and stops the Task from settling completed when the commit lost one.

## Requirements

1. MUST classify each refused path as refused by design — external to the
   repository, or absent from worktree and index — or as a lost output: an
   untracked executable, or a path that crosses a symbolic link.
2. MUST refuse to settle a Task `completed` when its commit lost an output.
3. MUST carry every lost path and its reason into the Task result and the Run
   outcome, in the vocabulary the existing console line uses.
4. MUST leave a Task whose only refusals are refused-by-design settling exactly
   as today, so a Spec Root outside the repository keeps working.
5. MUST keep the existing console line and Run Event for every refusal.
6. MUST name the collected set `lostStagePaths`, so the Task's own check and
   later readers agree.

## Subtasks

- [ ] Give the dropped path its class where the reason is written.
- [ ] Collect the lost paths at the commit boundary.
- [ ] Fail the Task settlement when the set is not empty, naming paths and cause.
- [ ] Cover a lost-output Task, a refused-by-design Task and a clean Task at the
      task-cycle seam.

## Acceptance Criteria

- [ ] A Task whose work produced an untracked executable does not settle
      completed, and its result names the path and reason.
- [ ] A Task whose only refusal is an external path settles completed.
- [ ] A Task with no refusal settles exactly as today.
- [ ] Every refusal still prints its console line and publishes its event.

## Context

- interface: `internal/daemon/task_engine.go`

## Verification

- `grep -q "lostStagePaths" internal/daemon/task_engine.go` — expected: exit 0; the commit boundary collects lost outputs. Before this Task nothing does.
- `out="$(go test -count=1 -run "^TestTaskCycleLostOutputRefusesCompletion$" ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Features 2 and 4; User Stories 2-3; Goals 2 and 4;
Success Metric 2; Declared intentional breaks;
`_techspec.md` → Implementation Design: A refusal has two classes; Interfaces;
API Contract 1; Build Order 2; ADR-0035; ADR-0057; ADR-0157.

## Result

Implementation:

- `DroppedStagePath.Lost` now classifies an untracked executable and a path
  crossing a symbolic link as lost output. External and absent paths remain
  refused by design.
- The Task commit boundary publishes every refusal before settlement, collects
  the lost set as `lostStagePaths`, and settles the Task failed with every lost
  path and reason. Refused-by-design and clean Tasks keep their prior
  settlement behavior.
- `TestTaskCycleLostOutputRefusesCompletion` covers the lost-output,
  external-refusal, and clean Task-cycle outcomes, including the existing
  console line and Run Event.

Focused evidence by acceptance criterion:

- Lost output refuses completion and names its cause:
  `go test -count=1 -run 'TestTaskCycleLostOutputRefusesCompletion/lost_output_fails_with_path_and_reason' ./internal/daemon`
  passed (2 tests reported by the test runner).
- An external-only refusal settles completed:
  `go test -count=1 -run 'TestTaskCycleLostOutputRefusesCompletion/(external_refusal_remains_completed|clean_Task_settles_as_before)|TestFilterStageablePaths(PreservesOtherRefusals|RefusesUntrackedExecutableWithMode)' ./internal/daemon`
  passed (8 tests reported by the test runner).
- A clean Task settles as before: the same focused command passed the
  `clean_Task_settles_as_before` subtest and observed one Task commit with no
  refusal event.
- Every refusal retains its console line and event: the lost-output and
  external-refusal subtests assert both surfaces; the filter cases assert the
  lost/refused-by-design classification.

Additional focused checks:

- `go test -count=1 ./internal/daemon` passed (335 tests reported by the test
  runner).
- `git diff --check` passed.

The Daemon-owned commands under `## Verification` were not run in this Agent
turn.

## Carry-forward provenance

- Source Run: `run_20260917T192120Z_0e305cd0796d7680`
- Source commit: `e423da61a2a11efb7eacd7575c5437a8a51bdaea`
