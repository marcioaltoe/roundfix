---
status: completed
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

## Result

Implementation:

- Added `carryForwardAcceptedStates` with `Stopped` and `Unresolved`; the
  invocation guard and candidate filter now use membership in that set.
- Reworded the outcome refusal to report the selected Run's actual outcome and
  both accepted outcomes. Other carry-forward diagnostics now refer to a Run
  without assuming why it ended.
- Parameterized the four existing carry-forward behavior tests through the
  fixture's Run-state argument and added a `Clean` negative case.

Acceptance evidence:

- Accepted-outcome set and both guards: `rtk rg -n
  'carryForwardAcceptedStates|slices.Contains\(carryForwardAcceptedStates'
  internal/cli/reconcile.go` found the set and both membership checks.
- Unchanged proofs: `rtk git diff --unified=0 --
  internal/cli/reconcile.go` showed no logic changes to settlement-commit count,
  committed status, declared-input comparison, symlink crossing, checkout
  cleanliness, HEAD stability, or the repository-local Specs Root requirement;
  changes outside the two guards were diagnostic text only.
- Outcome-neutral diagnostics: `rtk rg -n 'stopped Run|carry-forward accepts one
  stopped Run' internal/cli/reconcile.go` exited `1` with no matches.
- Shared proof behavior and named refusal: `rtk env
  GOCACHE=/private/tmp/roundfix-task-0118-01-gocache go test -count=1
  ./internal/cli -run
  '^TestCarryForward(AcceptsAnUnresolvedRun|RefusesATaskWhoseInputsMoved|RefusesRatherThanCarryingASubset|WithoutTheFlagReportsAndChangesNothing|RefusesAnUnacceptedRunOutcomeByName)$'`
  exited `0` with `ok roundfix/internal/cli`; the four behavior tests each ran
  for `Stopped` and `Unresolved`, and the negative case required `Clean`,
  `Stopped`, and `Unresolved` in the refusal.
- Formatting and patch hygiene: `rtk git diff --check` exited `0`.

Pre-change signal:

- Before the production edit, the focused `Unresolved` subtest exited `1`: the
  command returned Preflight exit `2` with `Run ... is not Stopped;
  carry-forward accepts one stopped Run`.

The Daemon-owned Verification command was not run during this Agent turn.
