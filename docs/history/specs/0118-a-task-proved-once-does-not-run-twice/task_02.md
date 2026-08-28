---
status: completed
type: backend
---

# Task: A Spec-scoped carry-forward query

The Implement Command needs to ask what a Spec's prior Runs would hand back.
The inspection that answers it already exists inside the reconcile command
path, bound to that command's own Run selection. Extract it so a second caller
reaches the same proofs instead of a second copy of them.

## Work

- Move the carry-forward inspection, its per-candidate proof, its settlement
  commit reader, and its refusal-reason assembly into their own file in the
  same package. This is a relocation, not a rewrite: behavior is identical and
  the reconcile path keeps calling the same logic.
- Add the Spec-scoped entry point. It takes the Run Database, the repository,
  the resolved Specs Root, and a Spec slug, and reports per prior terminal Run
  what carry-forward would do.
- Scope the query to the current repository, to the named Spec, and to the
  accepted-outcome set task_01 established.
- Skip a Run whose recorded Run Worktree is absent instead of failing the whole
  inspection. A released Run must never block a caller that is only asking.
  The reconcile path keeps its existing behavior for an explicitly named Run,
  where failing is right because the caller named that Run.
- Report how many candidates passed every proof, so a caller can act on the
  count without re-deriving it.
- Cover the query directly: a Spec with no prior Runs, one carriable Run, a Run
  whose Run Worktree is gone, and a Run whose Tasks all refuse.

## References

- `_prd.md` → Goal 3, Core Features 3 and 4
- `_techspec.md` → Build Order 2; Interfaces: `inspectSpecCarryForwards`,
  `specCarryForward`
- ADR-0053 keeps the proofs unchanged while the selection widens

## Verification
- `test -f internal/cli/carryforward.go && grep -q "func inspectSpecCarryForwards" internal/cli/carryforward.go && grep -q "TestInspectSpecCarryForwards" internal/cli/carryforward_test.go && go test -count=1 ./internal/cli -run 'TestInspectSpecCarryForwards|TestCarryForward' 2>&1 | grep -q "^ok"`

## Result

Implementation:

- Relocated the accepted-outcome set, carry-forward inspection and candidate
  proofs, settlement-commit reader, refusal assembly, and their proof-only Git
  helpers to `internal/cli/carryforward.go`. The Reconcile Command still calls
  `inspectCarryForwards` and keeps the existing explicit-Run failure behavior.
- Added `inspectSpecCarryForwards` and `specCarryForward`. The query lists
  terminal Runs from the Run Database for the current repository, keeps only
  Implement Runs for the named Spec in the Stopped-or-Unresolved set, skips a
  missing recorded Run Worktree, and applies the relocated proof pipeline to
  each surviving Run.
- Added `specCarryForward.carriable`, which counts candidates whose inspection
  reached the carry-forward action.
- Added direct coverage for no prior Runs, one carriable Run, a released Run
  Worktree, and an all-refusing Run. The carriable fixture also seeds Runs from
  another repository, another Spec, and an unaccepted outcome to prove the
  query's three scopes.

Focused checks:

- Red signal: `rtk env GOCACHE=/tmp/roundfix-task02-go-cache go test -count=1 ./internal/cli -run '^TestInspectSpecCarryForwards$'` failed to build before implementation because `inspectSpecCarryForwards` was undefined.
- `rtk env GOCACHE=/tmp/roundfix-task02-go-cache go test -count=1 ./internal/cli -run '^(TestInspectSpecCarryForwards|TestCarryForwardExplicitRunStillRefusesAMissingRunWorktree)$'` exited 0.
- `rtk env GOCACHE=/tmp/roundfix-task02-go-cache go test -count=1 ./internal/cli -run '^TestCarryForward'` exited 0, exercising the existing Reconcile Command carry-forward behaviors after relocation.

Acceptance evidence:

1. `internal/cli/carryforward.go` owns the relocated inspection, candidate
   proof, settlement reader, and refusal assembly; the existing
   `TestCarryForward*` suite exits 0 against the Reconcile Command path.
2. `TestInspectSpecCarryForwards/one_carriable_Run_is_reported_with_its_proven_candidate_count`
   calls the new Run Database/repository/Specs Root/Spec-slug entry point and
   observes its per-Run result.
3. That same subtest excludes the seeded other-repository, other-Spec, and
   Clean Runs while retaining the Unresolved Run.
4. `TestInspectSpecCarryForwards/a_Run_whose_recorded_Run_Worktree_is_gone_is_skipped`
   observes an empty result after releasing the Run Worktree, while
   `TestCarryForwardExplicitRunStillRefusesAMissingRunWorktree` observes a
   nonzero explicit reconcile result naming the selected Run.
5. The carriable and all-refusing subtests observe counts of one and zero,
   respectively, through `specCarryForward.carriable`.
6. All four required direct-query subtests exited 0 in the focused query check.

The Daemon-owned `## Verification` command was not run in this Agent turn.
