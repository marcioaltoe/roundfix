---
task: task_04
spec: 0144-a-run-stops-when-its-budget-is-spent
status: completed
type: docs
complexity: low
---

# Task 04: Tell the skill what a spent budget does

## Overview

The shipped Roundfix skill states that carry-forward accepts only two outcomes,
which stops being true here, and says nothing about a Run bounded by its budget.
This Task is an authorized tooling mutation and may change only the two bounded
files plus its own Task file.

## Requirements

1. MUST state, in the Roundfix skill, that an Implement Run is bounded by the
   configured maximum Run duration when the budget is enabled, and settles
   `BudgetExceeded` with a reason naming the maximum and the elapsed time.
2. MUST state that such a Run preserves its Run Worktree and Run Branch.
3. MUST state that Task Carry-Forward accepts `BudgetExceeded` beside `Stopped`
   and `Unresolved`, under its existing proof requirements.
4. MUST regenerate the distributed mirror so both copies stay identical.
5. MUST NOT change any path outside this Spec's bounded files, and MUST NOT
   change the skill's version or any other behavior it documents.

## Subtasks

- [x] Write the bound and its outcome into the skill.
- [x] Correct the carry-forward sentence.
- [x] Regenerate the mirror and confirm both copies match.

## Acceptance Criteria

- [ ] The skill describes the bound, the outcome and the preserved worktree.
- [ ] The carry-forward sentence names all three accepted outcomes.
- [ ] The mirror is byte-identical to the canonical copy.
- [ ] No path outside the bounded list changes.

## Context

- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `grep -q "BudgetExceeded" .agents/skills/roundfix/SKILL.md` — expected: exit 0; the canonical skill names the outcome. Before this Task it does not.
- `grep -q "BudgetExceeded" skills/roundfix/SKILL.md && diff -q .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — expected: exit 0; the mirror carries the same text and stays identical.

## References

`_prd.md` → Core Feature 4; Declared intentional breaks;
Project Constraints: Tooling authority;
`_techspec.md` → System Architecture: Shipped skill; Build Order 4;
`_authorization.md`.

## Result

- Updated the canonical Roundfix skill to document that an Implement Run with an enabled Run Budget is bounded by the configured maximum Run duration, settles `BudgetExceeded` with a reason naming the maximum and elapsed time, and preserves its Run Worktree and Run Branch.
- Updated Task Carry-Forward documentation to accept `BudgetExceeded` alongside `Stopped` and `Unresolved` under the existing proof requirements.
- Regenerated the distributed `skills/` mirror with `make skills-sync`; the canonical and mirror skill contents are identical in the resulting diff.
- Focused check: `git diff --check` passed. The changed-path listing contains only the two bounded skill files and this Task file (plus the pre-existing Task status edit).

Acceptance evidence:

- [x] The skill describes the bound, the outcome and the preserved worktree.
- [x] The carry-forward sentence names all three accepted outcomes.
- [x] The mirror was regenerated from the canonical copy and is identical.
- [x] No path outside the bounded list changed.
