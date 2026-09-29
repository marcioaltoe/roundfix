---
task: task_02
spec: 0182-delivery-that-reviews-and-retries-from-where-the-item-stands
status: pending
type: backend
complexity: high
---

# Task 02: Task Carry-Forward treats a Task completed on its target as nothing to carry

## Overview

Task Carry-Forward proves every settled Task of a Run against its target, including Tasks the target already completed. A carried Task's own file has changed on the target, so it always reads as a moved input and refuses the Run's whole set: a second `reconcile --carry-forward` refuses the operator's own first carry, the implement Preflight notes a fully carried Run as unavailable, and a Delivery Retry cannot carry a Run that overlaps work already on the item branch. This Task makes the shared inspection record a Task whose status is `completed` on the target as `already completed; nothing to carry`, never staged, proved or refused, while every remaining Task keeps every proof as one whole set. The settled Tasks come from the Run Database and the Run Branch; the completed statuses come from the target checkout's Task files; carried commits land on the checkout's branch only through the existing fast-forward.

## Requirements

1. MUST add `carryForwardCompletedAction = "already completed; nothing to carry"` in `internal/cli/carryforward.go`.
2. MUST make `inspectCarryForwards` read, per selected Run, the IDs of Tasks whose status is `completed` in the target's Spec folder `<resolvedSpecsRoot.Path>/<run.SpecSlug>`; when that folder does not exist no Task counts as completed, and any other load error MUST fail the inspection.
3. MUST make `inspectCarryForwardsForRun` record each such Task as a candidate with the completed action, its Run ID, its Task file, its settlement commit when it has exactly one, and no refusal reason, before any proof, never staging it, and never listing it among the Tasks `unstagedCarryForwards` reports after a conflict; the remaining Tasks MUST keep integration order and every existing proof.
4. MUST make `specCarryForward.wouldCarry` true only when at least one candidate is ready and no candidate carries a refusal reason, leave `carriable` counting ready candidates, and leave `carryForwardRefusalReason` unchanged.
5. MUST make `applyCarryForwards` stage only ready candidates and return before touching the checkout when none is ready, and make the Reconcile Command mark only previously ready candidates `carried forward` after a successful apply, so `roundfix reconcile <run-id> --carry-forward` exits `0` when every candidate is completed or carried.
6. MUST make `reportImplementNonCarriableCarryForwards` note only Runs whose candidates carry a refusal reason, and keep `selectImplementCarryForward` naming only ready Tasks.
7. MUST update the reconcile carry-forward and implement Preflight paragraphs of `docs/user-guide/commands.md` and of `.agents/skills/roundfix/SKILL.md` to contain `already completed; nothing to carry` and to state that such a Task never refuses the set, drop the sentence containing `effectively gets one carry-forward for overlapping work`, and run `make skills-sync`.
8. MUST replace, in the **Task Carry-Forward** entry of `CONTEXT.md`, the clause about a carried Task's own file becoming a moved input with `A Task already completed on the target is nothing to carry, so the act does not repeat itself.`
9. MUST put the new tests in `internal/cli/carryforward_completed_target_test.go` against real Git repositories through the public `reconcile` command and the implement Preflight's carry-forward inspection, reusing the package's carry-forward fixtures, and MUST keep `TestCarryForwardRefusesRatherThanCarryingASubset` and every other existing carry-forward test green without renaming any.

## Subtasks

- [ ] Read the target's completed Tasks and record them with the completed action.
- [ ] Make the carry predicate, the apply step and the reconcile report follow the new action.
- [ ] Silence the implement Preflight note for a Run with nothing to carry.
- [ ] Update the guide, the skill and its mirror, and the glossary entry.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] Over a Run whose `task_01` is completed on the checkout, `reconcile --carry-forward` exits `0`, carries only `task_02`, and its JSON reports `task_01` with action `already completed; nothing to carry`.
- [ ] Repeating the carry-forward of an already carried Run exits `0` with `HEAD` unchanged.
- [ ] With `task_01` completed on the checkout and a moved input of `task_02`, the carry-forward refuses the remaining set with `HEAD` unchanged.
- [ ] The implement Preflight names only the ready Task of a partly carried Run and prints no not-available note for a Run whose Tasks are all completed on the checkout.
- [ ] The guide, the skill and its mirror, and the glossary state the rule, and `make skills-sync-check` passes.

## Context

- interface: `internal/cli/carryforward.go`
- interface: `internal/cli/reconcile.go`
- interface: `internal/cli/implement.go`
- creates: `internal/cli/carryforward_completed_target_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `CONTEXT.md`
- instruction: `docs/adr/0170-a-task-already-completed-on-its-target-is-nothing-to-carry.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReconcileCarryForwardSkipsATaskCompletedOnTheCheckout|TestReconcileCarryForwardOfACarriedRunCarriesNothing|TestCarryForwardStillRefusesTheRemainingSetOnAMovedInput|TestImplementPreflightNamesOnlyTheTasksLeftToCarry|TestImplementPreflightIsSilentForARunWithNothingToCarry|TestCarryForwardRefusesRatherThanCarryingASubset)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReconcileCarryForwardSkipsATaskCompletedOnTheCheckout TestReconcileCarryForwardOfACarriedRunCarriesNothing TestCarryForwardStillRefusesTheRemainingSetOnAMovedInput TestImplementPreflightNamesOnlyTheTasksLeftToCarry TestImplementPreflightIsSilentForARunWithNothingToCarry TestCarryForwardRefusesRatherThanCarryingASubset; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the five new tests do not exist, so the command fails.
- `for file in docs/user-guide/commands.md .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "already completed; nothing to carry" || { printf 'missing phrase in %s: %s\n' "$file" "already completed; nothing to carry" >&2; exit 1; }; ! tr -s '[:space:]' ' ' < "$file" | grep -qF -- "effectively gets one carry-forward for overlapping work" || exit 1; done; make skills-sync-check` — expected: exit 0; before this Task the phrase is absent and the guide still carries the dropped sentence.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "A Task already completed on the target is nothing to carry, so the act does not repeat itself." || { printf 'missing phrase in %s: %s\n' CONTEXT.md "A Task already completed on the target is nothing to carry" >&2; exit 1; }` — expected: exit 0; before this Task the sentence is absent.

## References

- [_prd.md](_prd.md) — Goal 3; Core Feature 3; Success Metric 3
- [_techspec.md](_techspec.md) — A completed Task is nothing to carry; API Contracts 2–3; Testing Approach 2; Build Order 2
- ADR-0170; ADR-0053; ADR-0057; ADR-0026
