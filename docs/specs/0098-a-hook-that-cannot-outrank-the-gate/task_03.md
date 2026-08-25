---
status: pending
type: backend
---

# Task: Re-run Verification in Settle

Execute Task Verification in the selected settle surface.

## Work
- Reuse Verification logic from Implement
- Run commands verbatim (no edits)
- On pass: proceed to staging
- On fail: stop, print diagnostics

## Verification
- `grep -q "executeVerification\|runVerification" internal/cli/settle.go && grep -q "verificationCommand" internal/cli/settle.go`


## References

- Core Feature 2: Settle Recovery

## Result
Settle re-runs verification correctly.
