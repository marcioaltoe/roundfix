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

- `make test -k TestReportValidation | grep -q "ok"`

## Result
Mechanical stage validates new shape.
