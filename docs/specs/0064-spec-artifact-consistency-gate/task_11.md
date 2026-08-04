---
task: task_11
spec: 0064-spec-artifact-consistency-gate
status: pending
type: infra
complexity: low
---

# Task 11: Run the budget proof as its own gate step

## Overview

task_10 gives the sweep budget a test that measures under conditions it
declares. This slice makes the required local gate execute that proof on every
run, serially, so the budget is proven rather than merely provable.

It is a separate Task because it changes `Makefile`, and an authorized tooling
Task may mutate only its bounded files. Folding it into task_10 would fail the
tooling-authority audit that the QA gate performs against each Task commit.

## Requirements

1. MUST add a `spec-budget` target that executes the sweep-budget proof under
   the conditions task_10 declared, serially rather than inside the contended
   package sweep.
2. MUST include that step in the `verify` gate so every gate run proves the
   budget, and MUST place it where a failure stops the gate.
3. MUST NOT let a pipe hide the step's exit status — no pager, no `tail`, no
   pipeline whose last command masks a failure.
4. MUST declare the target in `.PHONY` and give it a `##` help description
   matching the file's existing style.
5. MUST change exactly one file: `Makefile`. This is the complete bounded file
   list from the Tooling authority row in `_prd.md` and `_techspec.md`, granted
   by the 2026-08-02 authorization record that names this Spec. Any other path
   is out of scope — stop rather than widen it.

## Subtasks

- [ ] Add the budget-proof target with its help description.
- [ ] Include it in `verify` where a failure stops the gate.
- [ ] Add it to `.PHONY`.

## Acceptance Criteria

- [ ] The gate step exits 0 on the current tree and measures the sweep rather
      than skipping it.
- [ ] The step fails the gate when the budget is exceeded, proven by a
      temporary limit reduction reverted within the same check.
- [ ] `verify` includes the step, and a failing step stops the gate rather
      than being reported and passed over.
- [ ] The target appears in `.PHONY` and in `make help`.
- [ ] `make verify` exits 0 on the current tree.
- [ ] `Makefile` is the only file this Task changed.

## Context

- instruction: `docs/agents/agent-instructions.md`
- interface: `Makefile`

## Verification

- `make verify` — expected: exit 0; the full gate including the budget proof.
- `grep -q "^\.PHONY:.*spec-budget" Makefile` — expected: exit 0.
- `grep -q "^verify:.*spec-budget" Makefile` — expected: exit 0; the gate
  includes the step.
- `make help | grep -q "spec-budget"` — expected: exit 0.
- `git diff --name-only HEAD | grep -v "^Makefile$" | grep -v "^docs/specs/0064-spec-artifact-consistency-gate/task_11.md$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded file and this Task file changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Feature 1.
- `_techspec.md` → Testing Approach; Integration Points.
- `qa/qa-report-2026-08-03.md` → F-001 required repair.
- `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`.
