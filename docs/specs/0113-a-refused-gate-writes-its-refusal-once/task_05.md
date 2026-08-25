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
- `make test -k TestNewestReportOnly | grep -q "ok"`

## Result
Mechanical stage reads newest report only.
