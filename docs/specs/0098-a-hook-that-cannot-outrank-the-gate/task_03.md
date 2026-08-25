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
- `make test -k TestSettleVerification | grep -q "ok"`


## References

- Core Feature 2: Settle Recovery

## Result
Settle re-runs verification correctly.
