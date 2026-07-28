# Temporary failure and cancellation

Fresh real-shell and Task-cycle suites passed for:

- child exit 75 producing the typed Temporary Verification Failure while
  preserving the command and `exec.ExitError` diagnostic chain;
- child exit 1 with timeout, listener, database, and port text remaining
  deterministic;
- one temporary failure followed by one exclusive retry and zero Agent repair;
- initial and retry diagnostic files retaining distinct paths and contents;
- a deterministic exclusive-retry failure using the one existing Agent
  repair, then numbered attempt 2;
- repeated exit 75 exhausting the Task-scoped retry budget without a second
  retry or hidden Agent turn;
- shared attempts draining before the exclusive retry and later shared work
  not bypassing it.

The public Stop regression composes
`TestRunStopByRunIDRecordsStopRequest` with
`TestTaskCycleStopRequestWhileQueuedForVerificationStartsNoCommandAndStaysResumable`.
The durable Stop Request causes the queued Task to publish Stop, start zero
child commands, remain `in_progress`, receive no terminal settlement, and
return its worker. The already-running attempt may finish; capacity and the
Run lock drain afterward.

The adjacent direct-cancellation case and later full-capacity acquisition also
passed. The full race suite and 20 repeated runs found no goroutine, permit, or
worker leak.
