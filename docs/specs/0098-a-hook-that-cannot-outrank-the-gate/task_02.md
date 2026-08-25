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
- `grep -q "completed" internal/cli/settle.go && grep -q "taskStatus == \"completed\"" internal/cli/settle.go`


## References

- User Story 2: Settle accepts completed
- Core Feature 2: Settle Recovery

## Result
Settle command accepts completed-but-uncommitted Tasks.
