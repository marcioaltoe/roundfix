---
status: pending
type: backend
---

# Task: Write Terminal Row on Precondition Refusal

Implement gate refusal path with terminal row.

## Work
- Gate refusal → single terminal row
- Row: `| 0 | blocked | precondition |`
- Frontmatter: `rows_blocked_precondition: 1`
- Store check name and reason
- Set verdict: `fail`

## Verification
- `grep -q "rows_blocked_precondition" internal/spec/qa.go && grep -q "precondition" internal/spec/qa.go && go test -count=1 ./internal/spec 2>&1 | grep -q "^ok"`


## References

- User Story 1: Gate writes valid report
- Core Feature 1: Terminal Row Writing

## Result
Gate writes valid report on precondition refusal.
