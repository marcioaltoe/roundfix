---
task: task_03
spec: 0177-runs-that-fit-their-budget-and-park-honestly
status: completed
type: backend
complexity: medium
---

# Task 03: A budget-stopped Run parks as a Run outcome

## Overview

`RunSpec` in `internal/cli/deliver_workflow.go` runs `roundfix implement` as a
child process and then reads the newest Implement Run of the Spec. It maps only
an `Unresolved` Run to `delivery.RunOutcomeUnresolved`; every other non-zero
exit becomes `result.failure("roundfix implement")`. `runCandidate` in
`internal/delivery/engine.go` wraps that error as `run Implement executor: ...`,
and `Engine.Run` parks the item `delivery-error: ...`. On 2026-09-28 Spec 0175's
Run `run_20260928T154312Z_63679c6540b34603` reached `BudgetExceeded` and was
parked that way, so a budget stop read as an infrastructure failure, and the
item's recorded Run ID was lost to the error path. The outcome is read from
two sources: the child's exit status, and the Run Database written by that
child. The operator reads the blocker in `roundfix deliver status`, and
`roundfix deliver retry` reads the item's Run ID to carry settled Tasks forward.
`latestImplementRun` returns whatever Run is newest, so a child that exits
before creating a Run (a Preflight refusal exits `2`) would otherwise be judged
by an older Run.

## Requirements

1. MUST add `RunOutcomeBudgetExceeded RunOutcome = "budget-exceeded"` and
   `BlockerRunBudgetExceeded = "run-budget-exceeded"` to
   `internal/delivery/engine.go`. `runCandidate` MUST record
   `strings.TrimSpace(result.RunID)` on the item and park a
   `RunOutcomeBudgetExceeded` result with `BlockerRunBudgetExceeded`, as it
   already does for `RunOutcomeUnresolved`, so later queue items continue.
2. MUST keep an executor error parked as `delivery-error: run Implement
   executor: ...` and every other `runCandidate` branch unchanged. `Retry`
   MUST resume a `run-budget-exceeded` item through `ItemRecovery.CarryForward`
   from its recorded Run ID, with no change to `Retry` required.
3. MUST make `RunSpec` read the newest Implement Run of the Spec
   (`latestImplementRun`) before it invokes `roundfix implement` and again
   after, and decide the result in one function of
   `internal/cli/deliver_workflow.go` from the exit status, the stderr, the Run
   read before and the Run read after:
   - exit `0` stays a Clean result with the candidate head;
   - exit `1` with an after-Run whose ID differs from the before-Run's maps
     `Unresolved` to `RunOutcomeUnresolved` and `BudgetExceeded` to
     `RunOutcomeBudgetExceeded`, carrying the Run ID and the trimmed stderr as
     the reason;
   - every other case returns `result.failure("roundfix implement")`.
4. MUST test that function directly, without spawning a process, and MUST NOT
   add a production-only test hook.
5. MUST state in the `deliver` section of `docs/user-guide/commands.md` and the
   Delivery Queue section of `.agents/skills/roundfix/SKILL.md` that a Run
   ending `BudgetExceeded` parks its item `run-budget-exceeded` with its Run ID
   and that `roundfix deliver retry` resumes it, using the literal
   `run-budget-exceeded`, then regenerate `skills/roundfix/SKILL.md` with
   `make skills-sync`.
6. MUST put the new delivery tests in `internal/delivery/budget_park_test.go`
   and the new CLI tests in `internal/cli/deliver_budget_park_test.go`, each
   negative case a test of its own, and MUST NOT rename or remove an existing
   top-level test.

## Subtasks

- [ ] Add the outcome and blocker and park them in `runCandidate`.
- [ ] Read the Run before and after the invocation and decide the result in one
      function.
- [ ] Document the blocker in the guide and the skill, then run
      `make skills-sync`.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A `RunOutcomeBudgetExceeded` result parks the item `run-budget-exceeded`
      with its Run ID, and the next queued item still merges.
- [ ] A retry of that item carries forward from its recorded Run.
- [ ] An executor error still parks `delivery-error: run Implement executor:
      ...`.
- [ ] Exit `1` with a new `BudgetExceeded` Run maps to the budget outcome, with
      a new `Unresolved` Run to the unresolved outcome, and with a `Failed` Run,
      no new Run or the Run that existed before the invocation to an error.

## Context

- interface: `internal/delivery/engine.go`
- creates: `internal/delivery/budget_park_test.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_budget_park_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestDeliveryParksABudgetExceededRunAsARunOutcome|TestDeliveryContinuesAfterABudgetParkedItem|TestDeliveryRetryResumesABudgetParkedItemFromItsRun|TestDeliveryStillParksAnExecutorErrorAsADeliveryError|TestDeliveryEngineParksDeclaredBlockersAndContinues|TestDeliveryRunResultMapsABudgetExceededRun|TestDeliveryRunResultMapsAnUnresolvedRun|TestDeliveryRunResultMapsACleanExit|TestDeliveryRunResultRefusesAFailedRun|TestDeliveryRunResultRefusesAnExitWithoutANewRun|TestDeliveryRunResultIgnoresARunThatExistedBeforeTheInvocation)$" ./internal/delivery ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDeliveryParksABudgetExceededRunAsARunOutcome TestDeliveryContinuesAfterABudgetParkedItem TestDeliveryRetryResumesABudgetParkedItemFromItsRun TestDeliveryStillParksAnExecutorErrorAsADeliveryError TestDeliveryEngineParksDeclaredBlockersAndContinues TestDeliveryRunResultMapsABudgetExceededRun TestDeliveryRunResultMapsAnUnresolvedRun TestDeliveryRunResultMapsACleanExit TestDeliveryRunResultRefusesAFailedRun TestDeliveryRunResultRefusesAnExitWithoutANewRun TestDeliveryRunResultIgnoresARunThatExistedBeforeTheInvocation; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "run-budget-exceeded" docs/user-guide/commands.md && grep -q "run-budget-exceeded" .agents/skills/roundfix/SKILL.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the new named tests exists and neither guide names `run-budget-exceeded`, so the command fails.

## References

- `_prd.md` → Goal 3; Core Feature 3; Success Metric 3; Decisions.
- `_techspec.md` → A budget stop parks as a Run outcome; API Contracts 4-5;
  Testing Approach 3; ADR-0020; ADR-0053; ADR-0113; ADR-0158.

## Result

### Implementation

- Added the `budget-exceeded` delivery Run outcome and
  `run-budget-exceeded` blocker. The Delivery Engine records the trimmed Run ID
  before parking this outcome and continues with later queue items.
- Changed the command-backed delivery workflow to read the newest Implement Run
  before and after invoking `roundfix implement`. A direct decision function
  maps only exit `1` plus a newly created `Unresolved` or `BudgetExceeded` Run;
  clean exits retain the candidate head and every other shape remains an
  executor failure.
- Kept retry generic: it carries forward from the parked item's recorded Run
  ID without a production change to `Retry`.
- Documented the operator-visible blocker and retry behavior in the user guide
  and canonical Roundfix skill, then regenerated the mirrored skill with
  `make skills-sync` (exit `0`).

### Focused checks

- Red signal: the new focused delivery and CLI tests initially failed to build
  because `RunOutcomeBudgetExceeded`, `BlockerRunBudgetExceeded`, and
  `deliveryRunResult` did not exist.
- `GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1
  ./internal/delivery -run
  '^(TestDeliveryParksABudgetExceededRunAsARunOutcome|TestDeliveryContinuesAfterABudgetParkedItem|TestDeliveryRetryResumesABudgetParkedItemFromItsRun|TestDeliveryStillParksAnExecutorErrorAsADeliveryError)$'`
  exited `0`.
- `GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1
  ./internal/cli -run '^TestDeliveryRunResult'` exited `0`.
- `GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1
  ./internal/delivery ./internal/cli` reached a green delivery package; the
  first CLI run was environment-blocked when two unrelated force-stop tests
  could not inspect the process table. Rerunning the full CLI package with
  process-table permission exited `0` in `76.613s`.
- `GOCACHE=/private/tmp/roundfix-task03-gocache make verify-incremental`
  exited `0`, including repository-wide `go vet`, `go test`, skill checks, and
  the build.
- The Task's declared `## Verification` command was not run; the Daemon owns
  that gate.

### Acceptance evidence

- Budget parking and queue continuation: focused tests
  `TestDeliveryParksABudgetExceededRunAsARunOutcome` and
  `TestDeliveryContinuesAfterABudgetParkedItem` exercised the real queue store
  and observed the `run-budget-exceeded` blocker, trimmed Run ID, and merged
  next item.
- Retry carry-forward: `TestDeliveryRetryResumesABudgetParkedItemFromItsRun`
  observed `ItemRecovery.CarryForward` receive the parked item's recorded Run
  ID.
- Executor-error preservation:
  `TestDeliveryStillParksAnExecutorErrorAsADeliveryError` observed the blocker
  prefix `delivery-error: run Implement executor:`.
- CLI decision boundary: the six `TestDeliveryRunResult...` tests directly
  exercised clean, new `BudgetExceeded`, new `Unresolved`, failed, no-new-Run,
  and pre-existing-Run cases without spawning a process.
