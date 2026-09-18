---
task: task_06
spec: 0144-a-run-stops-when-its-budget-is-spent
status: completed
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

- [x] Carry the deadline into Worktree setup and bootstrap.
- [x] Convert an expired budget into the outcome on the success path too.
- [x] Carry the deadline into integration and push.
- [x] Make every skill passage name the three outcomes, and regenerate the mirror.

## Acceptance Criteria

- [x] A fixture whose bootstrap outlives the budget ends with the Run and settles
      `BudgetExceeded`.
- [x] A fixture whose cycle returns without error after the deadline settles
      `BudgetExceeded`.
- [x] A fixture whose integration would run past the deadline is bounded.
- [x] A fixture that finishes inside its budget is unaffected.
- [x] No passage in either skill copy names only two accepted outcomes, and the
      copies are identical.

## Context

- interface: `internal/cli/implement.go`
- interface: `internal/daemon/task_engine.go`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -run "^TestImplementRunBudgetBoundsSetupAndIntegration$" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `out="$(go test -count=1 -run "^TestTaskCycleSettlesBudgetOutcomeWithoutError$" ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `tick="$(printf '\140')"; remaining="$(grep -n "Stopped${tick} or ${tick}Unresolved" .agents/skills/roundfix/SKILL.md)"; test -z "$remaining"` — expected: exit 0; no passage names only two accepted outcomes. Before this Task at least one does, so the command fails.
- `tick="$(printf '\140')"; mirror="$(grep -n "Stopped${tick} or ${tick}Unresolved" skills/roundfix/SKILL.md)"; test -z "$mirror" && diff -q .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — expected: exit 0; the mirror carries the corrected passages and stays identical to its canonical copy. Before this Task it still names only two outcomes, so the command fails.

## References

`_prd.md` → Core Features 1-4; Goals 1-3; Regression locks;
`_techspec.md` → Implementation Design: Deriving and honouring the deadline,
Settling the outcome, Accepting the outcome for carry-forward; API Contracts 1-2;
`_authorization.md`.

## Result

Implementation-ready behavior:

- The Implement Command derives one deadline from the persisted Run start and
  passes its context through Run Worktree creation, bootstrap, the Task cycle,
  integration, cleanup and optional push. Terminal settlement uses the parent
  context after the deadline cancels bounded work.
- The Task cycle applies that deadline to Task Worktree creation and bootstrap,
  and converts an elapsed deadline into `BudgetExceeded` even when the
  scheduler and QA return without error.
- Clean Run Worktree cleanup now follows the optional push, so a push cancelled
  by the Run Budget retains the Run Worktree and Run Branch for recovery.
- Every carry-forward passage names `BudgetExceeded`, `Stopped`, and
  `Unresolved`; `make skills-sync` regenerated the distributed mirror.

Focused evidence by acceptance criterion:

- Bootstrap bound: `GOCACHE=/private/tmp/roundfix-task06-go-cache rtk go test
  -count=1 -run Budget ./internal/cli ./internal/daemon` passed, including
  `TestImplementRunBudgetBoundsSetupAndIntegration/bootstrap`; the fixture
  observed context cancellation and stored `BudgetExceeded`.
- Success-path settlement: the same focused command passed
  `TestTaskCycleSettlesBudgetOutcomeWithoutError`, whose already-settled graph
  returns no scheduler error before the final budget check.
- Integration and push bounds: the same focused command passed the
  `integration` and `push` subtests; both blocking fakes observed cancellation,
  and both Runs stored `BudgetExceeded` with their Run Worktrees retained.
- Within-budget control: the same focused command passed the `inside budget`
  subtest and the existing `TestTaskCycleFinishesBeforeBudgetDeadline`; both
  retained their prior Clean/success behavior.
- Skill consistency: `make skills-sync` completed; SHA-256 for both skill copies
  is `07e6c98bef1e6b1b71a6ad0148264583c4454f4226da439d091045973660fc27`.
  Inspection of the reconcile contract, preflight candidates and exit-code
  table showed all three outcomes. `make baseline-digests` reported no derived
  changes.

Additional focused checks:

- `GOCACHE=/private/tmp/roundfix-task06-go-cache rtk go test -count=1 -run
  'Test(TaskCycleCreatesTaskWorktreesWithBootstrapBeforeAgentWork|TaskCycleTaskWorktreeBootstrapFailureIsolatesIndependentTasks|RunImplementAutoPushOutcomeMatrix|RunImplementAutoPushFailureEndsFailedAndJournalsPush|RunImplementCleanup)'
  ./internal/cli ./internal/daemon` passed 9 tests.
- `rtk git diff --check` passed.

The Daemon-owned commands in `## Verification` were not run during this Agent
turn.

Verification Feedback repair:

- Attempt 1 reached the third command but the Task parser truncated its inline
  Markdown code span at the escaped backtick, so the shell received an
  unterminated quoted command and exited `2` before `grep` ran.
- Both skill checks now create the literal backtick with `printf '\140'` and
  interpolate it into the exact search phrase. This keeps backticks out of the
  Markdown code-span body while preserving the intended byte-level check.
- `GOCACHE=/private/tmp/roundfix-task06-go-cache rtk go run ./cmd/roundfix spec
  check 0144-a-run-stops-when-its-budget-is-spent --format json` exited `0`,
  reported no findings, and reported `verification.executed: false`.
- Post-repair `rtk git diff --check` and the byte comparison of the canonical
  and mirrored skills both exited `0`.
