---
status: completed
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

The Implement Command now reads the repository's Run Window immediately after
opening the Run Database and before retention or Run Worktree pruning. A
cutoff that is not in the future returns a typed Preflight refusal; its reason
prints the literal cutoff and current time, and its `NextAction` names both the
force-set and clear commands. No other operational Run command gained this
check.

Acceptance evidence:

- Closed window: `TestImplementRunWindow/closed_window_refuses_before_mutating_preflight_work`
  observes exit `2`, empty stdout, the two literal instants, the standard
  no-side-effects block, the actionable commands, zero Run rows, no Run
  Worktree root, and a guard that fails if pruning starts.
- Future window: `TestImplementRunWindow/future_window_permits_Run_creation`
  drives a complete Implement Command flow and observes one Run row.
- Absent window: `TestImplementRunWindow/absent_window_permits_Run_creation`
  drives the same flow without a stored window and observes one Run row.
- Scope: diff inspection shows the enforcement is confined to the Implement
  Command; fetch, resolve, and watch paths are unchanged.

Focused checks:

- Before implementation,
  `GOCACHE=/private/tmp/roundfix-task03-go-build rtk go test ./internal/cli -run '^TestImplementRunWindow$/closed_window_refuses_before_mutating_preflight_work$'`
  failed because pruning was reached.
- After implementation, the same focused command passed (2 tests reported by
  the Go runner: parent and subtest).
- `GOCACHE=/private/tmp/roundfix-task03-go-build rtk go test ./internal/cli -run '^TestImplementRunWindow$'`
  passed (4 tests reported: parent and three subtests).
- The Task's declared `## Verification` command was not run; the Daemon owns
  that check and Task settlement.
