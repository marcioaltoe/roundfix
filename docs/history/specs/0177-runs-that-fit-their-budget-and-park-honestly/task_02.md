---
task: task_02
spec: 0177-runs-that-fit-their-budget-and-park-honestly
status: completed
type: backend
complexity: high
---

# Task 02: The Implement Run Budget renews at each Task settlement

## Overview

`TaskCycle` in `internal/daemon/task_engine.go` fixes
`plan.runBudgetDeadline = plan.RunStartedAt.Add(plan.MaxRunDuration)`, and
`implementRunContext` in `internal/cli/implement.go` puts the same deadline on
the whole Implement Run. A serial graph therefore reaches `BudgetExceeded` with
all of its work done but its QA gate. Specs 0176 and 0173 did exactly that on
2026-09-28 under `budget.max_run_duration: 2h`. The deadline is read from the
Run's start and from the Project Config, and three readers act on it. The Daemon
cancels Agent Sessions and stops starting Tasks. The CLI settles the Run
`BudgetExceeded` and decides whether integration, push and cleanup may run. An
operator or the Delivery Queue then reads the Run outcome. The time cap
**cancels**: when the renewed deadline passes, the Daemon cancels every Agent
Session it owns, as ADR-0158 requires, and nothing is abandoned running. A
renewal that fails to reach a later reader would let a finished Run be
relabelled `BudgetExceeded`; a renewal that renews too much would let a stalled
Run run on forever.

## Requirements

1. MUST replace the fixed `runBudgetDeadline` with one mutex-guarded budget
   owned by `TaskCycle`. It starts at `RunStartedAt + MaxRunDuration`.
   `runTaskScheduler` renews it after `integrateTaskSettlement` returns a
   settlement, completed or failed, to `engine.deps.Now()` at that settlement
   plus `MaxRunDuration`. A settlement observed at or after the current deadline
   is kept but does not renew it.
2. MUST cancel the cycle's context from a watchdog goroutine that `TaskCycle`
   owns and stops on every return path, when the current deadline passes. No
   Agent turn, Verification wait, repair turn or `runQAGate` step may keep a
   fixed deadline taken from an earlier renewal, and every checkpoint that read
   `plan.runBudgetDeadline` (`stopTaskCycleIfRequested`,
   `taskCycleResultWithBudgetOutcome` and the post-Agent checks) reads the
   current deadline. The cycle's error MUST still match
   `context.DeadlineExceeded`.
3. MUST add `BudgetDeadline time.Time` to `TaskCycleResult`: the deadline in
   force when the cycle returned, zero when the budget is disabled.
4. MUST keep the `BudgetExceeded` reason form `Run Budget exceeded: configured
   maximum <max>; elapsed <duration>`, with the elapsed time measured from the
   renewal point. After at least one settlement it ends ` since Task <id>
   settled.`, and it stays within `publicOutcomeReason`'s bound.
5. MUST keep the start-based deadline in `internal/cli/implement.go` for the
   work before the cycle (Run Worktree creation and bootstrap), pass the cycle
   a context without that deadline, bound integration, push and cleanup with
   `cycleResult.BudgetDeadline`, and make every post-cycle
   `implementRunBudgetExpired` check and `implementBudgetExceededReason` use
   that deadline instead of `run.CreatedAt`.
6. MUST keep the budget lines of `DefaultConfigYAML` in
   `internal/config/config.go` that `TestRenderedConfig` pins, and add the line
   `# An Implement Run's allowance renews at each Task settlement.` beside
   them. The watch loop in `internal/watch/watch.go` MUST NOT change, and
   `.roundfixrc.yml` MUST NOT change.
7. MUST add `docs/adr/0164-an-implement-run-budget-renews-at-each-task-settlement.md`
   with `status: accepted` lifecycle front matter per `docs/agents/docs-layout.md`.
   It records the decision, the rejected alternatives (a critical-path budget
   and holding the QA gate back) and the three measured Runs, and cites
   ADR-0158 and ADR-0137. It uses the phrase
   `renews at each Task settlement`.
8. MUST state that an Implement Run's allowance renews at each Task settlement,
   using the phrase `renews at each Task settlement`, in the `implement` Run
   Budget paragraph and the `window` sentence about `budget.max_run_duration`
   of `docs/user-guide/commands.md`, in `docs/user-guide/configuration.md`, and
   in the Implement Run Budget paragraph of `.agents/skills/roundfix/SKILL.md`,
   then regenerate `skills/roundfix/SKILL.md` with `make skills-sync`.
9. MUST put the new daemon tests in
   `internal/daemon/task_budget_renewal_test.go`, the new CLI tests in
   `internal/cli/implement_budget_renewal_test.go` and the new config test in
   `internal/config/budget_renewal_config_test.go`. Tests drive time through
   `engine.deps.Now` or `implementBudgetNow`, and any real-time bound uses
   `internal/testwait`, never a fixed sleep. The existing budget tests named in
   the Verification stay green, and no existing top-level test is renamed or
   removed.

## Subtasks

- [ ] Add the renewing budget, its watchdog and `BudgetDeadline` to the Daemon.
- [ ] Bound the CLI's post-cycle work with the renewed deadline.
- [ ] Add the rendered-config line, the ADR and the guide and skill text, then
      run `make skills-sync`.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A serial graph whose Tasks and QA gate each settle within 0.9 allowances
      of the previous settlement ends without a terminal outcome in the Daemon
      and `Clean` through `roundfix implement`, although its total exceeds one
      allowance.
- [ ] A settlement after the deadline is kept, does not renew, and the Run ends
      `BudgetExceeded` before the next Task starts.
- [ ] A stalled Task is cancelled one allowance after the last settlement, and
      earlier settlements keep their commits and `completed` status.
- [ ] A failed settlement renews the allowance.
- [ ] `BudgetDeadline` equals the last settlement plus the maximum, is zero when
      the budget is disabled, and bounds integration.
- [ ] The watchdog has stopped when `TaskCycle` returns.

## Context

- interface: `internal/daemon/task_engine.go`
- creates: `internal/daemon/task_budget_renewal_test.go`
- interface: `internal/cli/implement.go`
- creates: `internal/cli/implement_budget_renewal_test.go`
- interface: `internal/config/config.go`
- creates: `internal/config/budget_renewal_config_test.go`
- creates: `docs/adr/0164-an-implement-run-budget-renews-at-each-task-settlement.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/configuration.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTaskBudgetRenewsAtEachTaskSettlement|TestTaskBudgetStartsTheQAGateWithARenewedAllowance|TestTaskBudgetSettlementAfterTheDeadlineDoesNotRenew|TestTaskBudgetCancelsAStalledTaskOneAllowanceAfterTheLastSettlement|TestTaskBudgetAFailedSettlementRenewsTheAllowance|TestTaskCycleReportsTheRenewedBudgetDeadline|TestTaskCycleReportsNoBudgetDeadlineWhenTheBudgetIsDisabled|TestTaskBudgetReasonNamesTheSettlementThatRenewedIt|TestTaskBudgetWatchdogStopsWhenTheCycleReturns|TestTaskCycleEndsRunAtBudgetDeadline|TestBudgetExceededRunRecordsItsReason|TestTaskCycleRefusesToStartWorkPastItsBudget|TestTaskCycleSettlesBudgetOutcomeWithoutError|TestTaskCycleDisabledBudgetDerivesNoDeadline|TestTaskCycleFinishesBeforeBudgetDeadline|TestImplementSerialRunLongerThanItsBudgetEndsClean|TestImplementIntegrationRunsUnderTheRenewedDeadline|TestImplementEndsBudgetExceededWhenTheRenewedDeadlinePasses|TestBudgetClockCrossesAfterTheSettledTask|TestBudgetExceededKeepsWorkSettledBeforeTheBudget|TestImplementRunBudgetBoundsSetupAndIntegration|TestRenderedConfigStatesTheImplementBudgetRenewal|TestRenderedConfig|TestRunWithoutStopRequestKeepsRunBudgetBehavior|TestRunTransientReviewEvidenceExhaustsRunBudgetBeforeTimeout)$" ./internal/daemon ./internal/cli ./internal/config ./internal/watch 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTaskBudgetRenewsAtEachTaskSettlement TestTaskBudgetStartsTheQAGateWithARenewedAllowance TestTaskBudgetSettlementAfterTheDeadlineDoesNotRenew TestTaskBudgetCancelsAStalledTaskOneAllowanceAfterTheLastSettlement TestTaskBudgetAFailedSettlementRenewsTheAllowance TestTaskCycleReportsTheRenewedBudgetDeadline TestTaskCycleReportsNoBudgetDeadlineWhenTheBudgetIsDisabled TestTaskBudgetReasonNamesTheSettlementThatRenewedIt TestTaskBudgetWatchdogStopsWhenTheCycleReturns TestTaskCycleEndsRunAtBudgetDeadline TestBudgetExceededRunRecordsItsReason TestTaskCycleRefusesToStartWorkPastItsBudget TestTaskCycleSettlesBudgetOutcomeWithoutError TestTaskCycleDisabledBudgetDerivesNoDeadline TestTaskCycleFinishesBeforeBudgetDeadline TestImplementSerialRunLongerThanItsBudgetEndsClean TestImplementIntegrationRunsUnderTheRenewedDeadline TestImplementEndsBudgetExceededWhenTheRenewedDeadlinePasses TestBudgetClockCrossesAfterTheSettledTask TestBudgetExceededKeepsWorkSettledBeforeTheBudget TestImplementRunBudgetBoundsSetupAndIntegration TestRenderedConfigStatesTheImplementBudgetRenewal TestRenderedConfig TestRunWithoutStopRequestKeepsRunBudgetBehavior TestRunTransientReviewEvidenceExhaustsRunBudgetBeforeTimeout; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "^status: accepted" docs/adr/0164-an-implement-run-budget-renews-at-each-task-settlement.md && grep -q "ADR-0158" docs/adr/0164-an-implement-run-budget-renews-at-each-task-settlement.md && tr -s '[:space:]' ' ' < docs/adr/0164-an-implement-run-budget-renews-at-each-task-settlement.md | grep -qF -- "renews at each Task settlement" && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "renews at each Task settlement" && tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "renews at each Task settlement" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "renews at each Task settlement" && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the new named tests exists, ADR-0164 is absent and no guide says the allowance renews at each Task settlement, so the command fails.

## References

- `_prd.md` → Goal 2; Core Feature 2; Success Metric 2; Decisions.
- `_techspec.md` → The renewing Implement Run Budget; API Contracts 2-3;
  Testing Approach 2; ADR-0014; ADR-0057; ADR-0125; ADR-0137; ADR-0158.

## Result

Implemented one mutex-guarded Task-cycle budget whose watchdog follows deadline
renewals, cancels the shared cycle context at the current deadline, and stops
before `TaskCycle` returns. Every timely completed or failed Task settlement,
including the authored QA Task settlement, renews the allowance; a settlement
observed at or after the deadline remains settled without extending it.
`TaskCycleResult` now reports the active deadline and its renewing Task so the
CLI can bound integration, push, cleanup, expiry checks, and the public reason
from the renewal point while retaining the start-based setup deadline.

Added the rendered-config explanation, accepted ADR-0164, both user-guide
statements, and the canonical Roundfix skill statement. Ran the sanctioned
`make skills-sync`; the embedded skill is byte-identical to its canonical
source. The watch loop and `.roundfixrc.yml` remain unchanged.

Focused checks:

- `go test -race -count=1 -run '^TestTaskBudget' ./internal/daemon` with a
  task-scoped `GOCACHE`: passed.
- `go test -count=1 -run '^TestTask(Budget|CycleReports)' ./internal/daemon`
  with a task-scoped `GOCACHE`: passed.
- The three new Implement renewal tests in
  `internal/cli/implement_budget_renewal_test.go`: passed together.
- `TestBudgetClockCrossesAfterTheSettledTask`,
  `TestBudgetExceededKeepsWorkSettledBeforeTheBudget`,
  `TestImplementRunBudgetBoundsSetupAndIntegration`, `TestRenderedConfig`, and
  the two unchanged watch-budget tests: passed in focused runs.
- `TestRenderedConfigStatesTheImplementBudgetRenewal`: passed.
- `make verify-incremental` with a task-scoped `GOCACHE`: the sandboxed run
  reached the full suite but could not inspect the process table for two
  existing force-stop integration tests; the permission-enabled rerun passed
  formatting, vet, every Go package, skill checks, and the build.
- Phrase/reference inspection found `renews at each Task settlement` in both
  required command-guide locations, configuration, ADR-0164, and the Roundfix
  skill; ADR-0164 cites ADR-0158 and ADR-0137. `diff -r
  .agents/skills/roundfix skills/roundfix` and `git diff --check` passed.

Acceptance evidence:

- A serial graph renews across two 0.9-allowance Task settlements, the QA gate
  starts and settles under renewed allowances, and the public Implement flow
  ends Clean after more than one total allowance.
- A settlement observed at the deadline stays completed, retains the original
  deadline, ends `BudgetExceeded`, and leaves the next Task pending.
- A stalled Task is cancelled by the watchdog one allowance after the prior
  settlement while that earlier Task keeps its commit and `completed` status.
- A failed Task settlement advances the reported deadline by one allowance.
- The reported deadline equals the latest settlement plus the maximum, is zero
  with a disabled budget, and is the deadline observed by integration.
- The watchdog-stop regression check observes no Task-cycle watchdog goroutine
  after `TaskCycle` returns.
