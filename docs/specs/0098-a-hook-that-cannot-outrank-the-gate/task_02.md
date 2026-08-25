---
status: pending
type: backend
---

# Task: Extend Settle to Accept Completed Status

Update settle command contract to accept `completed` status.

## Work
- Modify settle preflight to accept `completed`
- Resolve settle surface in priority order
- Load Task file from selected surface

## Verification
- `make test -k TestSettleAcceptsCompleted | grep -q "ok"`


## References

- User Story 2: Settle accepts completed
- Core Feature 2: Settle Recovery

## Result
Settle command accepts completed-but-uncommitted Tasks.
