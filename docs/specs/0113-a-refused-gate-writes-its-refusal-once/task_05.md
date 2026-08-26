---
status: pending
type: backend
---

# Task: Read Newest Report Only

Read only newest report, ignore superseded.

## Work
- List all `qa-report-*.md`
- Sort by filename (newest first)
- Select first (newest)
- Ignore all older

## Verification
- `grep -q "NewestQAReport" internal/speccheck/mechanical.go && go test -count=1 ./internal/speccheck ./internal/spec 2>&1 | grep -q "^ok"`


## References

- Core Feature 3: Newest Report Only

## Result
Mechanical stage reads newest report only.
