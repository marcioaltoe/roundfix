---
status: pending
type: backend
---

# Task: A finding blocks the rows it names

A governance finding blocks the whole matrix. Measured: one finding blocked
fifteen of nineteen rows, so a round that cost a full Agent Session reported
signal about governance and nothing about function.

## Work

- A mechanical finding blocks the rows it names. A matrix row the finding does
  not name is measured rather than blocked.
- Keep the ability to block widely: a finding that names every row still blocks
  every row. What goes away is the implicit cascade, not the reach.
- Change nothing about withholding. ADR-0096 keeps the Agent Session withheld
  when a blocking machine fact is present before a matrix exists; this Task
  changes only how blame is attributed across a matrix that already exists.
- Change no verdict rule, no row contract, and no typed blocked-cause count.
- Cover a finding that names one row in a matrix of several, asserting the
  unnamed rows are measured, and a finding that names all of them, asserting it
  still blocks all.

## References

- `_prd.md` → Goal 4, User Story 3, Core Feature 4
- `_techspec.md` → Build Order 2; Interfaces: `BlockedRows`
- ADR-0080 keeps the typed blocked-cause counts this Task does not touch

## Verification
- `grep -q "func (finding MechanicalFinding) BlockedRows" internal/speccheck/mechanical.go || exit 1; grep -q "TestFindingBlocksOnlyTheRowsItNames" internal/speccheck/mechanical_test.go || exit 1; go test -count=1 ./internal/speccheck ./internal/daemon`
