# Deterministic repair flow

Current-build integration commands:

```text
rtk go test ./internal/cli -run 'TestRunImplement.*(VerificationCapacit|DaemonStatus|TemporaryVerification|QueuedCancellation)' -count=1
rtk go test ./internal/cli ./internal/daemon -run 'Test.*(VerificationCapacit|TemporaryVerification|QueuedCancellation)' -count=20
rtk go test -race ./internal/cli ./internal/daemon ./internal/runevent ./internal/tui -run 'Test.*(VerificationCapacit|DaemonStatus|TemporaryVerification|QueuedCancellation|WaitingForVerification)' -count=1
```

All exited 0. The deterministic real-shell fixture exits non-`75`, journals
one Verification Feedback turn, uses exactly two Agent calls, records
numbered attempts 1 and 2, reacquires Verification Capacity, and emits no
exclusive-retry metadata. Prompt/Skill checks also passed and prohibit Agents
from running declared Verification or editing Task status.
