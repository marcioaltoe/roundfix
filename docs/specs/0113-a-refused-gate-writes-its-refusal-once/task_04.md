---
status: pending
type: backend
---

# Task: Update Mechanical Stage Validation

Modify mechanical stage to accept terminal row.

## Work
- Empty table → refuse SC-REPORT-SHAPE
- Terminal blocked row → accept
- Validate status and provenance

## Verification
- `grep -q "provenance.*precondition\|ProvenancePrecondition" internal/speccheck/mechanical.go && go test -count=1 ./internal/speccheck 2>&1 | grep -q "^ok"`


## References

- Core Feature 2: Mechanical Stage Update

## Result
Mechanical stage validates new shape.
