---
status: pending
type: backend
---

# Task: Detect and Record Hook Refusal

Implement hook refusal detection in the Daemon's commit path.

## Work
- Add hook detection logic to commit boundary
- Classify refusal: parse stderr for hook markers
- Log with Run ID, Task ID, hook output
- Publish Run Event: `hook_refused`
- Leave staged changes in place
- Record Task as `completed`

## Verification
- `grep -r "hook_refused" internal/daemon/*_test.go | wc -l | grep -qE '^[1-9]'`

## References

- User Story 1: Hook refusal is detected and recorded
- Core Feature 1: Hook Refusal Detection

## Result
Hook refusal detection implemented and tested.
