---
status: pending
type: qa
---

# Task: QA Gate

Verify all deliverables.

## Work
- Verify hook detection works
- Verify settle accepts completed
- Verify verification re-runs
- Verify three cases resolve
- Verify invariant in docs

## Verification
- `roundfix spec check 0098-a-hook-that-cannot-outrank-the-gate --strict && make test -k TestHookStrictness`

## Result
QA gate passes all acceptance criteria.
