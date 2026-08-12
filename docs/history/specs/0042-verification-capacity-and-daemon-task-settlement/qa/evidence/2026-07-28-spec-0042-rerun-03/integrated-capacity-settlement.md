# Integrated capacity and settlement

The current disposable CLI macro suite passed:

```text
rtk env GOCACHE=/private/tmp/roundfix-qa-0042-rerun03-gocache \
  GOFLAGS=-buildvcs=false go test ./internal/cli \
  -run 'Test(RunImplement.*(VerificationCapacit|DaemonStatus|TemporaryVerification|QueuedCancellation)|RunStopByRunIDRecordsStopRequest|AttachAcceptsDocumentedFlagOrders|ParseAttachCommandAcceptsFlagsInAnyPosition|AttachSpecRun.*(Capacit|Legacy|Verification))' \
  -count=1
```

The suite uses disposable Git repositories, real local shell Verification
processes, and a fake only at the external ACP boundary. It observed two
simultaneous Agent turns with Task Capacity 2, a maximum of one active gate
with Verification Capacity 1, and `waiting` before `started`.

Agent-authored `completed` and `failed` values both normalized to
`in_progress`, reached real Daemon Verification, preserved Result content,
and settled only from the Daemon verdict. A deterministic non-75 failure
released capacity, received exactly one same-Session Verification Feedback
turn, and reacquired capacity for numbered attempt 2.

The current CLI and Daemon capacity, retry, Stop, and cancellation cases then
passed 20 consecutive runs:

```text
rtk env GOCACHE=/private/tmp/roundfix-qa-0042-rerun03-gocache \
  GOFLAGS=-buildvcs=false go test ./internal/cli ./internal/daemon \
  -run 'Test.*(VerificationCapacit|TemporaryVerification|QueuedCancellation|StopRequestWhileQueued)' \
  -count=20
```

The CLI package completed in 39.501 seconds and the Daemon package in 3.801
seconds with exit 0. Current stdout/exit assertions, Task/worktree state,
journal ordering, and diagnostic identities all passed.
