---
status: pending
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
