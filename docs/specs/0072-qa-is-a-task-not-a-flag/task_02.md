---
task: task_02
spec: 0072-qa-is-a-task-not-a-flag
status: pending
type: backend
complexity: high
---

# Task 02: Route the gate from the graph, not the request

## Overview

The Daemon stops asking the request whether to run the gate and asks the
graph. `plan.QA` derives from the loaded graph's gate declaration; a node of
type `qa`, once every dependency settles `completed`, is executed by the
existing `runQAGate` — same session kind, same before-snapshot, same report,
verdict, typed blocked-row counts, and `daemon.qa` events, per ADR-0080. A
`qa` node must never reach an Agent session: routing happens on Task Type,
and the engine refuses to schedule one as Agent work.

## Requirements

1. MUST derive the gate decision from the graph: a declared gate node means
   the cycle ends with the gate; a declined or legacy graph means it does
   not.
2. MUST execute the `qa` node through the existing gate step
   (`runQAGate`), recording its settlement in the node's task-file status
   like any Task, while report, verdict, and events stay byte-compatible
   with today's `--qa` output.
3. MUST NOT start the gate while any of its dependencies is unsettled — the
   existing withholding, now expressed by dependency semantics.
4. MUST refuse to schedule a `qa`-typed node as Agent work, with an error
   naming the node, even on a hand-edited graph.
5. MUST migrate the daemon tests that pass a QA request flag to graphs that
   declare a gate node, changing setup only — every assertion about report,
   verdict, counts, and events stays as it is.
6. MUST keep a failed Task suppressing the gate exactly as today: the gate
   node stays pending, resumable, and unreported.

## Subtasks

- [ ] Derive `plan.QA` (or its replacement) from the graph.
- [ ] Route `qa` nodes to `runQAGate`; settle their status from the verdict.
- [ ] Add the Agent-scheduling refusal for `qa` nodes.
- [ ] Migrate the gate tests from flag-driven to graph-driven setup.
- [ ] Prove the failed-Task path leaves the gate pending and resumable.

## Acceptance Criteria

- [ ] A graph ending in a gate node runs the gate after the last Task
      settles, never before, and the report and verdict match today's
      output for the same inputs.
- [ ] A graph with `qa: declined` or no declaration finishes its cycle with
      no gate step and no gate events.
- [ ] A cycle with one failed Task ends with the gate node still pending
      and no gate report.
- [ ] A `qa` node forced toward Agent scheduling fails with the node named.
- [ ] `git status --porcelain` shows no path outside `internal/daemon/`,
      `internal/spec/`, and this task file.

## Verification

- `go test ./internal/daemon -count=1 -run 'QA|Gate' -v | grep -q -- "--- PASS"`
  — expected: exit 0.
- `go test ./internal/daemon -count=1` — expected: exit 0.
- `go build -buildvcs=false ./...` — expected: exit 0.

## References

- `_prd.md` → Core Features 2, 5; Decisions (withholding stays).
- `_techspec.md` → Integration Points; Risks (a gate node must never reach
  an Agent); ADR-0080.
