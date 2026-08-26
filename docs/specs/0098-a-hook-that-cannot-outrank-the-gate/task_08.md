---
status: completed
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
- `roundfix spec check 0098-a-hook-that-cannot-outrank-the-gate --strict && go test -count=1 ./internal/daemon ./internal/cli 2>&1 | tail -1 | grep -q "ok"`


## References

- All user stories and core features

## Result
QA gate passes all acceptance criteria.
