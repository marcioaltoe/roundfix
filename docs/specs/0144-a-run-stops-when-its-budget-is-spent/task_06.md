---
task: task_06
spec: 0144-a-run-stops-when-its-budget-is-spent
status: pending
type: backend
complexity: medium
---

# Task 06: Bound the whole Run, and say so consistently

## Overview

Independent review found the bound covering less of the Run than the Spec
promises. The deadline reaches the Task cycle but not the Worktree setup before
it, nor the integration and push after it, and the budget becomes an outcome
only on paths that already returned an error — so a cycle that finishes normally
past its deadline settles as if it had not. The shipped skill also states the
new acceptance in one place while three others still name only two outcomes.
This slice closes all four, and it is the first of the two corrective Tasks the
contract allows.

## Requirements

1. MUST apply the Run's deadline to the Worktree setup that precedes the Task
   cycle, including each Task Worktree creation and the bootstrap command, so a
   setup that outlives the budget ends with the Run.
2. MUST settle `BudgetExceeded` when the deadline has passed even if the
   scheduler and the QA step return without error.
3. MUST apply the deadline to the integration and push that follow a clean
   cycle, so neither runs unbounded after the budget is spent.
4. MUST keep every cancellation proving termination, and MUST NOT settle
   `BudgetExceeded` for a Run that finished inside its budget.
5. MUST make the shipped skill state the three accepted carry-forward outcomes
   everywhere it names them, including the preflight candidates, the exit-code
   table and the reconcile help text, and MUST regenerate the mirror.
6. MUST change no path outside this Spec's bounded governed files, its own Task
   file and ordinary source.

## Subtasks

- [ ] Carry the deadline into Worktree setup and bootstrap.
- [ ] Convert an expired budget into the outcome on the success path too.
- [ ] Carry the deadline into integration and push.
- [ ] Make every skill passage name the three outcomes, and regenerate the mirror.

## Acceptance Criteria

- [ ] A fixture whose bootstrap outlives the budget ends with the Run and settles
      `BudgetExceeded`.
- [ ] A fixture whose cycle returns without error after the deadline settles
      `BudgetExceeded`.
- [ ] A fixture whose integration would run past the deadline is bounded.
- [ ] A fixture that finishes inside its budget is unaffected.
- [ ] No passage in either skill copy names only two accepted outcomes, and the
      copies are identical.

## Context

- interface: `internal/cli/implement.go`
- interface: `internal/daemon/task_engine.go`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -run "^TestImplementRunBudgetBoundsSetupAndIntegration$" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `out="$(go test -count=1 -run "^TestTaskCycleSettlesBudgetOutcomeWithoutError$" ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `remaining="$(grep -n "Stopped\` or \`Unresolved" .agents/skills/roundfix/SKILL.md)"; test -z "$remaining"` — expected: exit 0; no passage names only two accepted outcomes. Before this Task at least one does, so the command fails.
- `mirror="$(grep -n "Stopped\` or \`Unresolved" skills/roundfix/SKILL.md)"; test -z "$mirror" && diff -q .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — expected: exit 0; the mirror carries the corrected passages and stays identical to its canonical copy. Before this Task it still names only two outcomes, so the command fails.

## References

`_prd.md` → Core Features 1-4; Goals 1-3; Regression locks;
`_techspec.md` → Implementation Design: Deriving and honouring the deadline,
Settling the outcome, Accepting the outcome for carry-forward; API Contracts 1-2;
`_authorization.md`.
