---
status: pending
type: test
---

# Task: Acceptance Verification

Verify three measured hook refusal cases.

## Work
- Case 1: 82-line function over 80-char limit
- Case 2: 2462-line file over 500-line limit
- Case 3: `sort()` instead of `toSorted()`

## Verification
```bash
make test -k TestHookRefusalRecovery | grep -q "ok.*3 subtests"
```

## Result
Three cases resolve via settle without losing work.
