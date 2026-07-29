---
task: task_04
spec: 0054-tooling-task-and-verification-hygiene
status: completed
type: backend
complexity: high
---

# Task 04: Prove the repository is green before a gate-bound Task starts

## Overview

A Task whose Verification is the repository-wide gate is unsatisfiable when
the repository is already red for unrelated reasons, and today that only
surfaces after an Agent has done its work — the Task settles failed with a
reason that implies its own work is broken. Check the precondition before the
Agent Session is created and fail with the real cause instead.

## Requirements

1. MUST run the repository-wide gate once, before creating the Agent Session,
   for a Task whose declared Verification includes that gate.
2. MUST settle such a Task with a precondition reason naming the failing
   check and a bounded excerpt of its output when the gate fails on entry,
   and MUST create no Agent Session for it.
3. MUST distinguish the precondition failure from a post-Agent Verification
   failure in the Task reason, the Run Event, and the reported outcome, so a
   red repository never reads as broken Task work.
4. MUST NOT run the precondition for a Task whose Verification does not
   include the repository-wide gate.
5. MUST NOT consume the single post-Agent Verification repair on a Task that
   never started Agent work.
6. SHOULD reuse the existing Verification execution path rather than adding a
   second command runner.

## Subtasks

- [ ] Detect, per Task, whether its declared Verification includes the
      repository-wide gate.
- [ ] Run the gate once before Agent Session creation for those Tasks.
- [ ] Settle a red-on-entry Task with a distinct precondition reason and no
      Agent Session.
- [ ] Publish the precondition outcome as its own Run Event kind or reason.

## Acceptance Criteria

- [ ] Given a repository failing the gate for an unrelated reason, a
      gate-bound Task settles with a precondition reason naming the failing
      check, and no Agent Session is created for it.
- [ ] Given a green repository, the same Task runs its Agent normally and
      the gate runs after Agent work exactly as today.
- [ ] A Task whose Verification does not include the repository-wide gate
      never pays the precondition run.
- [ ] A precondition failure is distinguishable from a post-Agent
      Verification failure in the Task reason and the Run Event Stream.
- [ ] The single post-Agent Verification repair remains available to a Task
      that did start Agent work.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/spec/task.go`
- interface: `internal/daemon/task_engine_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/daemon/ ./internal/spec/` — expected: pass, including the red-on-entry precondition and the unchanged green path.

## References

`_prd.md` → User Story 5, Core Feature 7; `_techspec.md` → Build Order 5,
Interfaces: green-on-entry precondition; ADR-0038, ADR-0057.

## Result

Implemented the repository-green precondition in the Task worker before Agent
Session ownership. The Implement Command now passes the configured repository
Verification command into the Task cycle; Tasks that declare that exact
command run it once through the existing Verification Capacity, command
runner, diagnostics, and Run Event path. A command failure settles the Task
without Agent preparation or Verification Feedback, with a
`repository not green on entry` reason that names the command, exit status,
diagnostic path, and a maximum 1,024-byte output tail.

Run Events classify this failure as `precondition` with reason
`repository_not_green_on_entry`. Existing temporary-failure metadata remains
unchanged for post-Agent Verification.

Acceptance evidence:

- Red repository / no Agent Session: focused test
  `TestTaskCycleRepositoryGatePreconditionFailureStartsNoAgentSession` proves
  one gate invocation, failed Task settlement, no Agent preparation or work,
  the bounded output tail, and the distinct Task outcome reason.
- Green repository / unchanged post-Agent gate:
  `TestTaskCycleRepositoryGatePreconditionPassesBeforeAgentAndPostVerification`
  proves the order `precondition → Agent → full declared Verification →
  commit`, including the repository gate after Agent work.
- Non-gate Task:
  `TestTaskCycleWithoutRepositoryGateSkipsPrecondition` proves the only
  Verification invocation occurs after Agent work.
- Distinct Run Event Stream outcome: the red-entry test projects the emitted
  Verification event through `ProjectStreamEvent` and asserts classification
  `precondition` plus reason `repository_not_green_on_entry`; the Task outcome
  asserts the separate `repository not green on entry` reason.
- Repair remains available:
  `TestTaskCycleRepositoryGatePreconditionDoesNotConsumeVerificationRepair`
  proves a passing precondition followed by deterministic post-Agent failure
  still returns once to the same Agent flow and reruns the full declared
  Verification.

Focused checks:

- `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache go test -run 'TestTaskCycle(RepositoryGatePrecondition|WithoutRepositoryGate|VerificationFailureRepairsSameSessionAndRerunsFullSequence|TemporaryVerificationFlowPassesExclusiveRetryWithoutAgentRepair)' ./internal/daemon`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-runevent-gocache go test -run 'TestProjectStreamEvent.*Verification|TestProjectStreamEvent' ./internal/runevent`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-cli-gocache go test -run '^$' ./internal/cli`
  — package compiled; no tests selected.
- The first focused invocation without a cache override did not reach
  compilation because the sandbox denied the host Go cache. The exact focused
  check was rerun with the writable temporary cache above.

The Daemon-owned commands in `## Verification` were not run.
