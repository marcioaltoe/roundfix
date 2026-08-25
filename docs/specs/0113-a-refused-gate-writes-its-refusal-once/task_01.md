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
```bash
grep -q "| 0 | blocked | precondition |" /tmp/qa-report*.md
```

## Result
Gate writes valid report on precondition refusal.
