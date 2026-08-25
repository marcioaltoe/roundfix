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

- `make test -k TestPreconditionDetection | grep -q "ok"`

## Result
Precondition failures detected correctly.
