---
status: pending
type: backend
---

# Task: Preflight refuses a closed window

The bound has to be a contract, not a discipline. A caller that must remember
to consult a guard is a guard that fails the first time the caller forgets —
the failure this Spec exists to remove.

## Work

- Add the check to the `implement` preflight block, after `store.Open` because
  it needs the store, and before `pruneTerminalRunWorktreeDebris` because a
  refused Run must not first mutate worktree state.
- Absent window means no check. A cutoff in the future proceeds. A passed
  cutoff refuses: no Run created, no side effects, exit `2`.
- The error implements `NextAction() string`, the existing extension point
  whose precedent is `preflight.TargetMismatch.NextAction`, so the refusal
  renders with the same shape as every other preflight failure.
- Per ADR-0133 the reason names both instants literally — the cutoff and the
  current time — so a reader never computes the delta to understand it. The
  next action names the force-set and the clear commands.
- `implement` only. `fetch`, `resolve`, and `watch` answer an already-open Pull
  Request; refusing those by clock strands a review Round.

## References

- User Story 2: The refusal names the cutoff and how to act
- Core Feature 1: A stored bound on Run creation

## Verification
- `grep -q "RunWindow" internal/cli/implement.go && grep -q "NextAction" internal/cli/implement.go && go test -count=1 ./internal/cli -run 'TestImplementRunWindow' 2>&1 | grep -q "^ok"`

## Result
A passed Run Window refuses Run creation from the Preflight, whether or not the
caller thought to check.
