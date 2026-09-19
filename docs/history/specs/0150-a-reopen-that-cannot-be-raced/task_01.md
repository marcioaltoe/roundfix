---
task: task_01
spec: 0150-a-reopen-that-cannot-be-raced
status: completed
type: backend
complexity: medium
---

# Task 01: Recheck the gate before writing to it

## Overview

`preflightReopen` reads the Task Graph and returns a plan; `ReopenGate` writes
from that plan. Nothing holds the tree in between, so a dependency that settles
`completed` in the gap leaves the command resetting a gate that is no longer
stale and recording dependencies that are no longer incomplete.

## Requirements

1. MUST re-derive the gate's staleness and its stale dependency set immediately
   before the replacement, through the same recovery load the preflight uses.
2. MUST refuse with exit 2 and write nothing when the gate is no longer stale,
   or when the stale dependency set differs from the one preflight saw.
3. MUST name what changed in the refusal rather than failing generically.
4. MUST NOT add a second implementation of the staleness predicate.
5. MUST keep every refusal and success Spec 0149 shipped.

## Subtasks

- [x] Re-derive and compare the plan before the write.
- [x] Add the refusal and its message.
- [x] Add tests for the changed-gate and unchanged-gate paths.

## Acceptance Criteria

- [x] A fixture whose stale dependency is completed between the preflight read
      and the write is refused, and the QA Task file is byte-identical.
- [x] A fixture whose gate is unchanged reopens exactly as before.
- [x] The reopen tests Spec 0149 shipped pass unchanged.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/reopen.go`

## Verification

- `out="$(go test -count=1 -run "^TestReopenRefusesWhenTheGateChangedBeforeTheWrite" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `out="$(go test -count=1 -v -run "^TestReopen" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReopenRefusesWhenTheGateChangedBeforeTheWrite"` — expected: exit 0; the whole reopen suite passes and the new case is among the cases that passed. `-v` is required: without it the run prints only a package summary and no test name, so the grep could never match.

## References

- [_techspec.md](_techspec.md) — The recheck

## Result

### Implementation

- `preflightReopen` and the write-time recheck now use the same
  `deriveReopenPlan` recovery load. The recheck runs before `spec.ReopenGate`
  and compares the terminal QA Task id and stale dependency set.
- A gate that is no longer stale, a changed QA Task, or a changed stale
  dependency set refuses with exit 2. The diagnostic says the gate changed
  after preflight and includes the underlying change; no Task bytes are
  written.
- The unchanged plan uses the re-derived paths and stale dependency ids for the
  existing atomic reopen operation.

### Focused checks

- Pre-change signal:
  `rtk env GOCACHE=/tmp/roundfix-go-cache go test ./internal/cli -count=1 -run 'TestReopenRefusesWhenThe(Gate|StaleDependencySet)ChangedBeforeTheWrite'`
  failed to build because `reopenFromPlan` did not exist.
- After implementation:
  `rtk env GOCACHE=/tmp/roundfix-go-cache go test ./internal/cli -count=1 -run 'TestReopen(StaleGatePreservesEvidence|RefusesWhenThe(Gate|StaleDependencySet)ChangedBeforeTheWrite)$'`
  exited 0.
- Reopen regression sweep:
  `rtk env GOCACHE=/tmp/roundfix-go-cache go test ./internal/cli -count=1 -run Reopen`
  exited 0.
- Broader package check:
  `rtk env GOCACHE=/tmp/roundfix-go-cache go test ./internal/cli -count=1`
  reached two unrelated force-stop integration failures because the sandbox
  denied process-table access (`operation not permitted`).
- `rtk git diff --check` exited 0.
- The Task's declared `## Verification` commands were not run; the Daemon owns
  them in this execution mode.

### Acceptance evidence

- `TestReopenRefusesWhenTheGateChangedBeforeTheWrite` completes the stale
  dependency after the first preflight, observes exit 2 with the named
  no-longer-stale reason, and compares the QA Task bytes before and after.
- `TestReopenRefusesWhenTheStaleDependencySetChangedBeforeTheWrite` keeps the
  gate stale while changing its dependency set, observes the named exit-2
  refusal, and compares the QA Task bytes before and after.
- `TestReopenStaleGatePreservesEvidence` exercises the unchanged plan and the
  reopen regression sweep covers the Spec 0149 reopen tests.

### Follow-ups

None for this Task slice.
