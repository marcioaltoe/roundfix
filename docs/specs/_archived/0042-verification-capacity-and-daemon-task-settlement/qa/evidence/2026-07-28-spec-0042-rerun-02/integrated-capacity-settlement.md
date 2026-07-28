# Integrated capacity and settlement

The current disposable CLI macro suite passed on build `ffd6852`:

```text
rtk go test ./internal/cli -run \
  'Test(RunImplement.*(VerificationCapacit|DaemonStatus|TemporaryVerification|QueuedCancellation)|RunStopByRunIDRecordsStopRequest|AttachAcceptsDocumentedFlagOrders)' \
  -count=1
```

It uses disposable Git repositories, real local shell Verification processes,
and a fake only at the external ACP boundary. The suite observed:

- Task Capacity 2 overlaps Agent work while Verification Capacity 1 never
  exceeds one active attempt;
- waiting precedes started for every attempt;
- Agent-authored `completed` and `failed` values both reach real Daemon
  Verification and only Daemon settlement writes the terminal status;
- a one-time exit 75 receives one exclusive retry, keeps distinct diagnostic
  identities, uses no Agent repair, and can end Clean;
- repeated exit 75 ends non-clean with no second retry and preserves the Task
  Worktree;
- a deterministic non-75 failure gets exactly one Verification Feedback turn
  and numbered attempt 2;
- queued cancellation starts no child and retains resumable Task state;
- normal stdout and exit codes retain the public contract.

The same capacity, temporary, cancellation, and Stop cases passed 20
consecutive runs across `internal/cli` and `internal/daemon`. Current config,
prompt, typed-error, event-projection, and Task-status focused suites also
passed.
