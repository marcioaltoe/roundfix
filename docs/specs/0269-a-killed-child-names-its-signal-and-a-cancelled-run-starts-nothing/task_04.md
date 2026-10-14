---
task: task_04
spec: 0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing
status: completed
type: backend
complexity: medium
---

# Task 04: A cancelled Run starts no Verification command after an immediate grant

## Overview

`acquireVerificationCapacity` returns an immediate `TryAcquire` grant without
reading the context, so a Task that published Waiting for Verification just
before its Run was cancelled can still take the freed capacity, publish
`started` and call the verifier: the second start
`TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement`
saw on Linux CI (`_techspec.md` → The second start, found). This Task applies
the queued path's post-grant cancellation check to the immediate grant and
adds two deterministic tests, one for a Verification attempt and one for the
pre-work probe that shares the grant. It is verifiable on its own through the
named tests.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [the Verification Capacity cancellation test sees a second command start](references/2026-10-13-the-verification-capacity-cancellation-test-starts-a-second-command.md):
   the path that starts the second command is the immediate grant, and a
   cancelled Run starts no Verification command after its cancellation.
2. MUST make `acquireVerificationCapacity`, after a `TryAcquire` grant, read
   the context; when it is cancelled, release the capacity and return the
   error the queued path returns,
   `acquire Verification Capacity for run %q Task %s attempt %d retry %d: %w`,
   wrapping the context's error (Invariant 5). The queued path, its Stop
   Request handling, `fairVerificationGate` and the rule that a Stop Request
   does not stop an attempt that owned capacity immediately MUST stay as they
   are.
3. MUST add `internal/daemon/verification_capacity_cancellation_test.go` with
   `TestVerificationAttemptCancelledAtTheImmediateGrantStartsNoCommand`: it
   calls `runTaskVerificationRequest` with a free gate of capacity 1, one
   command, a counting verifier and a `Publish` that records each phase and
   cancels the context when it records the Waiting phase; it asserts an error
   matching `context.Canceled`, zero verifier calls, no `started` phase, and
   that the gate grants again afterwards.
4. MUST add to the same file
   `TestPreWorkProbeCancelledAtTheImmediateGrantStartsNoProbe`: it calls
   `verifyTaskPreWork` for a Task with one Verification command, with a free
   gate and a `RunStateStore` wrapper that delegates `UpdateRunState` and then
   cancels the context; it asserts an error matching `context.Canceled`, that
   the verifier received no probe request, and that the gate grants again.
5. MUST keep both tests deterministic: each cancels at the instant the race
   needs through a callback the code already makes, and neither sleeps, polls
   or bounds by real time (ADR-0271). Each calls `t.Parallel()` first
   (ADR-0259). On the tree before this Task both MUST fail with one start.
6. MUST keep
   `TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement`
   and `TestTaskCycleStopRequestWhileQueuedForVerificationStartsNoCommandAndStaysResumable`
   byte-identical; `internal/daemon/task_engine_test.go` does not change.

## Subtasks

- [ ] Add the post-grant cancellation check to the immediate grant.
- [ ] Add the attempt test.
- [ ] Add the pre-work probe test.

## Acceptance Criteria

- [ ] Both new tests pass, and both fail on the tree without the check.
- [ ] The two existing queued-cancellation and Stop Request tests pass unchanged.
- [ ] No other daemon behavior changes: the gate tests pass unchanged.

## Consumers

- `acquireVerificationCapacity`'s cancellation contract — consumers: `./internal/daemon`

## Context

- instruction: `docs/adr/0056-spec-runs-separate-task-and-verification-capacity.md`
- instruction: `internal/daemon/task_engine_test.go`
- interface: `internal/daemon/task_engine.go`
- creates: `internal/daemon/verification_capacity_cancellation_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestVerificationAttemptCancelledAtTheImmediateGrantStartsNoCommand|TestPreWorkProbeCancelledAtTheImmediateGrantStartsNoProbe|TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement|TestTaskCycleStopRequestWhileQueuedForVerificationStartsNoCommandAndStaysResumable|TestTaskCycleVerificationGateTryAcquireDistinguishesQueuedAttempts|TestTaskCycleExclusiveRetryCancellationRestoresFullSharedCapacity)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestVerificationAttemptCancelledAtTheImmediateGrantStartsNoCommand TestPreWorkProbeCancelledAtTheImmediateGrantStartsNoProbe TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement TestTaskCycleStopRequestWhileQueuedForVerificationStartsNoCommandAndStaysResumable TestTaskCycleVerificationGateTryAcquireDistinguishesQueuedAttempts TestTaskCycleExclusiveRetryCancellationRestoresFullSharedCapacity; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two new tests do not exist, so no PASS line is printed for them.
- `out="$(go test -count=20 -run '^(TestVerificationAttemptCancelledAtTheImmediateGrantStartsNoCommand|TestPreWorkProbeCancelledAtTheImmediateGrantStartsNoProbe|TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement)$' -v ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -c -- '--- PASS: TestVerificationAttemptCancelledAtTheImmediateGrantStartsNoCommand ' | grep -qx 20 || { printf 'the attempt test did not pass 20 times\n' >&2; exit 1; }` — expected: exit 0; before this Task the test does not exist, so the count is 0.

## References

- [_prd.md](_prd.md) — Goals 3-4; User Story 3; Core Feature 4; Success Metric 3
- [_techspec.md](_techspec.md) — Invariant 5; The second start, found; Testing Approach; Build Order 4
- [2026-10-13-the-verification-capacity-cancellation-test-starts-a-second-command.md](references/2026-10-13-the-verification-capacity-cancellation-test-starts-a-second-command.md)
- ADR-0056; ADR-0259; ADR-0271

## Result

The immediate `TryAcquire` grant checks `ctx.Err()`, releases its permit on
cancellation, and returns the queued path's acquisition error with `%w`.
The queued path, Stop Request handling and `fairVerificationGate` are unchanged.
The Daemon-owned `status: in_progress` remains untouched.

The Daemon's first Verification attempt identified the missing probe test.
The Verification Feedback repair added
`TestPreWorkProbeCancelledAtTheImmediateGrantStartsNoProbe` alongside the
attempt regression. It wraps a disposable real Run store, delegates the state
update successfully, then cancels synchronously before capacity acquisition.
It verifies the persisted Verifying state, cancellation matching, the required
acquisition error, zero verifier requests, and immediate permit reuse.
Both new tests call `t.Parallel()` first and use no sleep, polling or time budget.

### Focused evidence

- Acceptance criterion 1: `go test -count=1 -v -run
  '^Test(VerificationAttempt|PreWorkProbe)CancelledAtTheImmediateGrantStartsNo(Command|Probe)$'
  ./internal/daemon` exited 0 and printed PASS for both new tests.
- Acceptance criterion 1, before-check evidence: the same focused selection
  with `go test -overlay <temporary-overlay.json> -count=1 -v -run
  '^Test(VerificationAttempt|PreWorkProbe)CancelledAtTheImmediateGrantStartsNo(Command|Probe)$'
  ./internal/daemon` exited 1 and both new tests failed. A temporary Go overlay
  replaced only `task_engine.go`, removing exactly the four added lines of the
  immediate-grant check. The working source stayed intact and the temporary
  overlay was removed afterwards. The attempt observed one verifier request
  and a started phase. The probe received a probe-command cancellation error
  instead of the required capacity-acquisition error.
- Requirement 5 caveat: the probe already has a downstream cancellation guard
  in `ProbeCommands`, so it cannot observe one verifier start on the original
  tree when cancelled after the state update. Its acquisition-error assertion
  makes it fail deterministically without this Task's check while preserving
  the required zero-request assertion. No probe implementation or existing
  test was changed to manufacture a start.
- Acceptance criteria 2 and 3: `go test -count=1 -v -run
  '^(TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement|TestTaskCycleStopRequestWhileQueuedForVerificationStartsNoCommandAndStaysResumable|TestTaskCycleVerificationGateTryAcquireDistinguishesQueuedAttempts|TestTaskCycleExclusiveRetryCancellationRestoresFullSharedCapacity)$'
  ./internal/daemon` exited 0 after the repair; all four named tests printed
  PASS. `git diff --exit-code HEAD -- internal/daemon/task_engine_test.go`
  exited 0, confirming that the existing tests remain byte-identical.
- `git diff --check` exited 0. Only the assigned Task file,
  `internal/daemon/task_engine.go`, and the new cancellation test file changed.

The authored Verification commands, repository Verification and incremental
Verification were not rerun; the Daemon owns those checks and Task settlement.
No commit, push, pull request, Task Graph or other Task-file edit was made.
