---
task: task_06
spec: 0071-verification-cost
status: pending
type: infra
complexity: medium
---

# Task 06: Assert a suite-time budget

## Overview

Without a budget, a change that makes verification materially slower is
discovered several Runs later, when someone notices a cycle taking longer. This
Task asserts the suite completes within a recorded budget, so the regression
fails at the moment it is introduced.

## Requirements

1. MUST assert that the full suite completes within a recorded budget, derived
   from what the parallelised suite actually achieves rather than from the old
   baseline.
2. MUST fail when the suite exceeds the budget, naming the measured time and
   the budget.
3. MUST record the budget where a reader can find it and change it
   deliberately, never as a magic number inside a recipe.
4. MUST tolerate ordinary machine variance without flapping, so the budget
   catches regressions rather than noise.
5. MUST change only `Makefile` among protected tooling, per this Spec's Tooling
   authority row.
6. MUST leave the existing verification gate's other stages unchanged.

## Subtasks

- [ ] Derive the budget from the parallelised suite's measured time.
- [ ] Add the assertion with its failure message.
- [ ] Record the budget where it is readable and deliberately changeable.
- [ ] Confirm the gate's other stages are untouched.

## Acceptance Criteria

- [ ] The suite passes the budget assertion on the current tree.
- [ ] A deliberately slowed test trips the budget and the failure names the
      measured time and the budget.
- [ ] The budget is recorded in a readable location, not embedded as a bare
      number in a recipe line.
- [ ] Repeated runs on an unmodified tree do not flap.
- [ ] The verification gate's other stages run unchanged.
- [ ] `git status --porcelain` shows no path outside `Makefile` and this task
      file.

## Context

- instruction: `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`
- interface: `Makefile`

## Verification

- `make verify` — expected: exit 0 on the unmodified tree, budget satisfied.
- `make verify` — expected: exit 0 on a second consecutive run, proving the
  budget does not flap.
- `grep -qi 'budget' Makefile` — expected: exit 0; the assertion exists.
- `git diff --name-only HEAD -- internal/ .github/ | grep -q . && exit 1 || exit 0`
  — expected: exit 0; nothing outside the bounded path changed.

## References

- `_prd.md` → Core Features 5; Success Metrics (a deliberately slow test trips
  the budget).
- `_techspec.md` → Build Order 6; Project Constraints: Tooling authority.
