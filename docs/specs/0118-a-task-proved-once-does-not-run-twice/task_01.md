---
status: pending
type: backend
---

# Task: Carry-forward accepts an Unresolved Run

The act that hands proved work back refuses the one outcome that strands it.
Two guards tie it to the Stopped outcome — the invocation refusal and the
candidate-set skip — while every proof downstream of them is already
outcome-agnostic. Widen the guards, and change nothing else.

## Work

- Introduce one accepted-outcome set holding Stopped and Unresolved, and make
  both guards test membership in it rather than comparing against Stopped.
- Leave every proof byte-for-byte unchanged: settlement-commit count, committed
  status, declared-input comparison, symlink crossing, checkout cleanliness,
  HEAD stability, and the repository-local Specs Root requirement. A proof that
  moves in this Task is a defect, not a refinement.
- Rewrite the guard's refusal so it names the outcome the Run actually has and
  the outcomes carry-forward accepts. The sentence claiming carry-forward
  accepts one stopped Run goes away.
- Make the remaining carry-forward diagnostics outcome-neutral. Several say
  "stopped Run" while describing evidence that has nothing to do with why the
  Run ended.
- Parameterize the four existing carry-forward tests over both accepted
  outcomes. Reusing the fixture helper's state parameter is what proves the
  proofs did not change: the same assertions must hold for Stopped and for
  Unresolved.
- Add a case for a Run whose outcome is neither, asserting the refusal names it.

## References

- `_prd.md` → Goal 2, Core Features 1 and 2, User Story 4
- `_techspec.md` → Build Order 1; Interfaces: `carryForwardAcceptedStates`
- ADR-0023 keeps a non-Clean Run's worktree and branch as the inspection
  surface; ADR-0053 keeps reconciliation proof-based

## Verification
- `grep -q "carryForwardAcceptedStates" internal/cli/reconcile.go && grep -q "TestCarryForwardAcceptsAnUnresolvedRun" internal/cli/reconcile_test.go && ! grep -q "carry-forward accepts one stopped Run" internal/cli/reconcile.go && go test -count=1 ./internal/cli -run 'TestCarryForward' 2>&1 | grep -q "^ok"`
