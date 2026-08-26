---
status: pending
type: backend
---

# Task: Detect Precondition Failure

Detect which precondition failed and why.

## Work
- Parse `spec check --strict` output
- Extract error codes
- Extract check name and reason
- Store for report writing

## Verification
- `grep -q "PreconditionRefusal\|preconditionRefusal" internal/speccheck/mechanical.go && go test -count=1 ./internal/speccheck 2>&1 | grep -q "^ok"`


## References

- User Story 2: Precondition captured
- Core Feature 1: Terminal Row Writing

## Result
Precondition failures detected correctly.
