---
task: task_01
spec: 0066-run-teardown-reclaims-what-it-created
status: pending
type: backend
complexity: high
---

# Task 01: Terminate the tree and prove each process gone

## Overview

Terminating a Run reaches the recorded PID and stops there. Four `acpx`
processes from a QA fixture were found still running after three days and six
hours, reparented to init, pointing at worktrees that no longer existed —
nothing had ever told them to stop.

This slice terminates the descendants Roundfix started and returns one outcome
per process, so an unprovable termination is visible rather than silent.

## Requirements

1. MUST terminate the descendants Roundfix started for a Run, including a child
   that outlives its immediate parent.
2. MUST return one outcome per process recording whether absence was observed
   and, when it was not, why.
3. MUST report an unprovable termination as unproven, never as success. ADR-0044
   reclaims orphaned locks by reading that distinction, so a host that cannot
   answer must never look like a terminated process.
4. MUST bound the walk by recorded ownership. Terminating a process Roundfix did
   not start is out of scope and MUST NOT happen.
5. MUST preserve `ProveOwner`'s existing refusal on a proven identity mismatch,
   and MUST NOT widen the `--owner-identity-unreadable` last resort.
6. MUST leave every existing exported symbol in `internal/store` behaving as it
   does today.

## Subtasks

- [ ] Walk the descendants Roundfix recorded ownership for.
- [ ] Terminate and observe absence per process.
- [ ] Return per-process outcomes with reasons.
- [ ] Add the outliving-grandchild fixture and the unprovable case.

## Acceptance Criteria

- [ ] A fixture starting a grandchild that outlives its immediate parent leaves
      no descendant running after termination.
- [ ] Each terminated process returns an outcome recording proven absence.
- [ ] A process whose absence cannot be observed returns unproven with a
      non-empty reason, and is never reported as terminated.
- [ ] A process Roundfix did not start is never signalled, proven by a fixture
      that starts an unrelated process and asserts it survives.
- [ ] The existing identity-mismatch refusal still refuses.
- [ ] Existing `internal/store` tests pass unchanged.

## Context

- interface: `internal/store/process.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/store -count=1 -run 'Terminate|Tree|Owner|Unproven' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the termination tests ran and passed.
- `go test ./internal/store -count=1` — expected: exit 0.
- `go test ./internal/store -count=20 -run 'TerminateTree' 2>&1 | grep -qE "FAIL|fatal error" && exit 1 || exit 0`
  — expected: exit 0; twenty repetitions with no failure, since process
  lifecycle is where flakiness hides.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Features 1 and 2; Goals; Success Metric 1.
- `_techspec.md` → Interfaces; Build Order 1; Risks & Considerations.
- ADR-0044, ADR-0052.
