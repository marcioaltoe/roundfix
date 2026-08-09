---
task: task_03
spec: 0090-a-gate-that-could-have-failed
status: pending
type: backend
complexity: high
---

# Task 03: Refuse a Task whose gate already passes

## Overview

The Spec's centre. Before an Agent Session owner is created, the Daemon runs the
Task's own `## Verification` commands against the unchanged tree. Any command
that exits zero there proves nothing about work that has not happened, so the
Task is refused with that command named and no Agent turn is spent. The seam is
the sibling of `verifyRepositoryPrecondition`, which already runs a Verification
before the Agent and settles the Task on its result.

## Requirements

1. MUST run every command of the Task's `## Verification` against the tree as it
   stands before the Agent runs, at the point where the repository precondition
   already runs.
2. MUST settle the Task `failed` when any command exits zero, and MUST name every
   offending command in the terminal reason.
3. MUST NOT create an Agent Session or spend an Agent turn when the Task is
   refused.
4. MUST treat a probe command whose verdict could not be observed as `unknown`
   using the cause Task 02 introduced, and MUST NOT read it as evidence that the
   command is sound.
5. MUST leave a Task whose commands all exit non-zero to proceed exactly as it
   does today.
6. MUST break the characterization case Task 01 declared for the vacuous gate,
   and update that case to the new behaviour in the same commit.

## Subtasks

- [ ] Run the Task's Verification before the Agent Session owner exists.
- [ ] Classify each command and settle on a vacuous one.
- [ ] Keep the unrefused path byte-identical in behaviour.

## Acceptance Criteria

- [ ] A Task whose Verification passes against the unchanged tree settles
      `failed` with every offending command named.
- [ ] No Agent Session is created for a refused Task.
- [ ] A Task whose commands all fail against the unchanged tree reaches its Agent
      turn unchanged.
- [ ] A probe command with no observable verdict is recorded `unknown` and does
      not clear the Task.

## Rehearsal Cases

- Case: a Task carrying Spec 0089's exact `task_05` gate,
  `grep -q 'reasoning_effort: xhigh' .roundfixrc.yml`, against a fixture file
  that already contains that string; Observation: the Task settles `failed`, the
  terminal reason names that command, and the fake Agent runner records zero
  invocations.
- Case: a Task whose single Verification command exits non-zero against the
  unchanged tree; Observation: the Agent runner is invoked exactly once and the
  Task proceeds to its ordinary post-Agent Verification.
- Case: a Task whose Verification command times out during the probe;
  Observation: the outcome carries the unknown cause from Task 02, and the Task
  does not settle `completed`.

## Bounded scope

This Task may create or modify only:

- `internal/daemon/task_engine.go`
- `internal/daemon/task_engine_test.go`
- `internal/daemon/verification_probe_characterization_test.go`
- `docs/specs/0090-a-gate-that-could-have-failed/task_03.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestPreWorkProbe' -count=1 -v 2>&1 | grep -q '^--- PASS: TestPreWorkProbeRefusesATaskWhoseGateAlreadyPasses'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestPreWorkProbe' -count=1 -v 2>&1 | grep -q '^--- PASS: TestPreWorkProbeSpendsNoAgentTurnOnARefusedTask'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestPreWorkProbe' -count=1 -v 2>&1 | grep -q '^--- PASS: TestPreWorkProbeLeavesAFailingGateOnItsOrdinaryPath'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -count=1 2>&1 | grep -q '^ok'` — expected: exits 0, proving the declared break was updated.

## References

- `_prd.md` → Goal 1; Core Features, three controls.
- `_techspec.md` → Build Order 3; System Architecture.
- ADR-0109.
