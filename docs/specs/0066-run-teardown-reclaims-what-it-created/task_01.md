---
task: task_01
spec: 0066-run-teardown-reclaims-what-it-created
status: completed
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

## Result

### Implementation

- Added `TerminateTreeAndWait`, which proves the recorded Run owner before any
  signal, discovers only processes attributable to that owner, captures a live
  identity for each descendant where the host supports it, and rescans after
  termination so a descendant discovered during teardown is not omitted.
- Added one `TerminationOutcome` per discovered PID. Proven absence has an empty
  reason; identity, signal, or absence-observation failures remain unproven and
  retain the diagnostic reason.
- Bounded Unix ownership by the dedicated Run session when present, with a
  parent walk only for a live non-session owner. Linux reads `/proc`; macOS
  queries the owner's process group and retains only members of its session.
  Windows walks the process snapshot's recorded parent links. Unsupported hosts
  fail closed.
- Kept `ProveOwner` and `TerminateAndWait` as the identity and termination
  authorities. A proven owner-identity mismatch still stops before discovery or
  signalling, and the legacy empty-identity behavior is unchanged.
- Extended the real Unix helper fixture with an owner session, an exited middle
  process, and a reparented grandchild. Added distinct coverage for proven
  per-PID outcomes, an unrelated survivor, and an unobservable-absence result.

### Focused checks

- Red pre-change signal:
  `GOCACHE=<worktree>/.gocache go test ./internal/store -count=1 -run '^TestOwnerProcessControllerTerminateTree'`
  failed to compile because `TerminateTreeAndWait`, `ownedProcesses`, and
  `signalProcess` did not exist.
- `GOCACHE=<worktree>/.gocache go test ./internal/store -count=1 -run '^TestOwnerProcessController' -v -timeout=30s`
  passed. This ran the new tree cases and the existing graceful, force-kill,
  identity-match, identity-mismatch, legacy-identity, and absence cases.
- `GOCACHE=<worktree>/.gocache GOOS=linux GOARCH=amd64 go test -c ./internal/store -o /private/tmp/roundfix-store-linux.test`
  passed.
- `GOCACHE=<worktree>/.gocache GOOS=windows GOARCH=amd64 go test -c ./internal/store -o /private/tmp/roundfix-store-windows.test.exe`
  passed.
- `git diff --check` passed.

### Acceptance evidence

- Outliving grandchild: `TestOwnerProcessControllerTerminateTreeProvesOutlivingGrandchildGone`
  passed after the middle process exited; both the owner and reparented
  grandchild were observed absent.
- Proven outcome per process: the same fixture asserted exactly one outcome for
  each owned PID, with `Proven=true` and an empty reason.
- Unprovable absence: `TestOwnerProcessControllerTerminateTreeReportsUnprovenAbsence`
  passed and asserted `Proven=false`, a non-empty host diagnostic, and no
  terminated-success representation.
- Unrelated process safety: `TestOwnerProcessControllerTerminateTreeLeavesUnrelatedProcessRunning`
  passed and observed the unrelated helper still alive after tree termination.
- Identity mismatch refusal: the focused controller suite passed both
  `TestOwnerProcessControllerRefusesMismatchedOwnerIdentity` and
  `TestOwnerProcessControllerProveOwnerRefusesMismatchedIdentity`.
- Existing `internal/store` behavior: the focused existing owner-process suite
  passed unchanged. The full package and repository checks remain with Daemon
  Verification.

### Daemon verification

The commands under `## Verification` were intentionally not run in this Agent
turn. The Daemon owns those checks and the terminal Task status.
