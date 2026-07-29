---
task: task_04
spec: 0054-tooling-task-and-verification-hygiene
status: pending
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
