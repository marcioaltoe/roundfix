---
type: fix
status: open
created: 2026-10-13
spec: null
reason: null
---

# The Verification Capacity cancellation test sees a second command start

## Symptom

`TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement`
failed on Linux CI with `expected zero Verification commands for queued Task,
total starts=2`. It passed 400 of 400 on macOS and on the CI rerun.

## Where

`internal/daemon/task_engine_test.go`. `fairVerificationGate.acquireShared`
checks `ctx.Err()` before granting, so the second start comes from another
path: a retry of the active command on cancellation, or a probe outside the
gate.

## Expected

The path that starts the second command is found. If it is production code, a
cancelled Run never starts a Verification command after its cancellation; if
it is the fixture, the test synchronizes on the event it asserts.

## Evidence

PR #517, CI run 38426389923 (2026-10-13).
