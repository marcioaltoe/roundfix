---
spec: 0177-runs-that-fit-their-budget-and-park-honestly
status: active
created: 2026-09-28
surfaces: [backend, cli, docs]
---

# Runs that fit their budget and park honestly

## Executive Summary

Stop counting read-only `instruction:` paths as Wave collisions. Renew an
Implement Run's budget at each Task settlement. Park a Delivery Queue item whose
Run ended `BudgetExceeded` as a Run outcome. Run carry-forward's staging commits
without repository hooks. Report a line-bound multi-word phrase check against a
Markdown file at authoring.

## Project Constraints

- Identifier strategy: applicable — new blocker `run-budget-exceeded`, delivery
  Run outcome `budget-exceeded` and check code `SC-VERIFY-WRAP-FRAGILE`, each in
  its family's existing form. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the Run
  Database only; no credential and no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0014, ADR-0020, ADR-0025, ADR-0053,
  ADR-0056, ADR-0057, ADR-0080, ADR-0091, ADR-0093, ADR-0096, ADR-0097,
  ADR-0104, ADR-0113, ADR-0117, ADR-0125, ADR-0130, ADR-0133, ADR-0137,
  ADR-0148, ADR-0155, ADR-0156 and ADR-0158 hold; ADR-0038, ADR-0127,
  ADR-0159 and ADR-0160 do not apply; the Spec adds ADR-0164. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `.agents/skills/write-tasks/references/task-template.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `internal/speccheck/coherence.go`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `skills/baseline_skill_contract_test.go`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Instruction paths and Wave collisions

`declaredTaskTouches` in `internal/spec/collision.go` adds every existing
`## Context` path of a Task to its touch set, whatever its kind. The write-tasks
skill already says a Task declares every path it edits under `interface:` or
`creates:` and never under `instruction:`, so an `instruction:` path is read-only
by contract and is not evidence that the Task touches it. `declaredTaskTouches`
skips `spec.ContextKindInstruction` references. `interface:` and `creates:`
references, Verification operands (`TaskVerificationFiles`) and prior-Run
settlement paths (`priorRunTaskTouches`) keep their current behavior. A path a
Task names both as an instruction and as a Verification operand still collides
through the operand.

`spec.Collisions` is the one reader: `detectWaveCollisions` in
`internal/speccheck/coherence.go` (`SC-WAVE-COLLISION`) and
`refuseTaskPlanWaveCollisions` in `internal/daemon/task_engine.go` both call it,
so authoring and dispatch change together and neither file changes. The
write-tasks skill states that an `instruction:` path never makes two Tasks
collide.

## The renewing Implement Run Budget

Today `TaskCycle` in `internal/daemon/task_engine.go` fixes
`plan.runBudgetDeadline = plan.RunStartedAt.Add(plan.MaxRunDuration)`, wraps
the cycle in `taskAgentContext(ctx, runBudgetDeadline)`, and gives each Agent
turn, wait and QA step a context with that same fixed deadline (the calls near
the Agent invocation, the Verification wait, the repair turn and `runQAGate`).
`stopTaskCycleIfRequested` and `taskCycleResultWithBudgetOutcome` compare
`engine.deps.Now()` with it. In `internal/cli/implement.go`,
`implementRunContext` puts the same start-based deadline on the whole Run, and
`implementRunBudgetExpired(run.CreatedAt, ...)` guards worktree setup, the
post-cycle outcome, integration, push and cleanup.

The Daemon replaces the fixed deadline with one shared, mutex-guarded budget
owned by `TaskCycle`:

- It starts at `RunStartedAt + MaxRunDuration`.
- `runTaskScheduler` renews it after `integrateTaskSettlement` returns a
  settlement, completed or failed: the new deadline is `engine.deps.Now()` at
  that settlement plus `MaxRunDuration`. A settlement observed at or after the
  current deadline does not renew it, so that settlement is kept and the Run
  still ends `BudgetExceeded`.
- A watchdog goroutine that `TaskCycle` owns and stops on return cancels the
  cycle's context when the current deadline passes. No operation inside the
  cycle keeps a fixed deadline derived from an earlier renewal. The cycle's
  error still matches `context.DeadlineExceeded`, as the existing budget tests
  require.
- Every checkpoint that read `plan.runBudgetDeadline` reads the current
  deadline.
- `TaskCycleResult` gains `BudgetDeadline time.Time`, the deadline in force
  when the cycle returned (zero when the budget is disabled).
- The reason keeps the form `Run Budget exceeded: configured maximum <max>;
  elapsed <duration>`. The elapsed time is measured from the renewal point,
  and after at least one settlement the reason adds `since Task <id>
  settled`. It stays within `publicOutcomeReason`'s bound.

`internal/cli/implement.go` keeps the start-based deadline for everything
before the cycle (worktree creation and bootstrap), because no Task can have
settled yet. It passes the cycle a context without that deadline and, after the
cycle, bounds integration, push and cleanup with
`cycleResult.BudgetDeadline`. Every post-cycle `implementRunBudgetExpired`
check compares against that deadline instead of `run.CreatedAt`, so a Run that
finished its graph inside the renewed allowance is never relabelled
`BudgetExceeded` because it started long ago. `DefaultConfigYAML` in
`internal/config/config.go` keeps its current budget lines, which
`TestRenderedConfig` pins, and adds the sentence `An Implement Run's allowance
renews at each Task settlement.`, which ADR-0137 requires where the budget is
configured. The watch loop in `internal/watch/watch.go` is unchanged.

`docs/adr/0164-an-implement-run-budget-renews-at-each-task-settlement.md`
records the decision, citing ADR-0158 and ADR-0137 and the three measured Runs.
`docs/user-guide/commands.md` (the `implement` Run Budget paragraph and the
`window` sentence about `budget.max_run_duration`),
`docs/user-guide/configuration.md` and the Roundfix skill's Implement Run
Budget paragraph state that the allowance renews at each Task settlement.

## A budget stop parks as a Run outcome

`internal/delivery/engine.go` adds `RunOutcomeBudgetExceeded RunOutcome =
"budget-exceeded"` and `BlockerRunBudgetExceeded = "run-budget-exceeded"`.
`runCandidate` records `strings.TrimSpace(result.RunID)` and parks a
`RunOutcomeBudgetExceeded` result with `BlockerRunBudgetExceeded`, exactly as it
parks `RunOutcomeUnresolved`. `Retry` already accepts any parked item and
carries forward from `item.RunID`, so a budget-parked item resumes without
change.

`RunSpec` in `internal/cli/deliver_workflow.go` reads the newest Implement Run of
the Spec (`latestImplementRun`) before it invokes `roundfix implement` and again
after. One function decides the result from four inputs: the invocation's exit
status and stderr, the Run read before and the Run read after.

- Exit `0` stays a Clean result with the candidate head.
- Exit `1` with an after-Run whose ID differs from the before-Run's maps
  `Unresolved` to `RunOutcomeUnresolved` and `BudgetExceeded` to
  `RunOutcomeBudgetExceeded`, carrying the Run ID and stderr as the reason.
- Every other case returns `result.failure("roundfix implement")`, which the
  engine parks as `delivery-error`: another exit status, another Run state, or
  no new Run.

The function is tested directly, without spawning a process.
`docs/user-guide/commands.md` (the `deliver` section) and the Roundfix skill's
Delivery Queue section name `run-budget-exceeded` and its retry.

## Carry-forward without repository hooks

`stageCarryForwardCandidate` in `internal/cli/carryforward.go` runs
`cherry-pick <commit>` and `commit --amend --no-edit` in the staging worktree,
and `abortCarryForwardStaging` runs `cherry-pick --abort`, all through
`reconcileGitRaw`, which applies `-c core.fsmonitor=false -c
commit.gpgSign=false`. Measured with Git 2.54.0, the cherry-pick runs
`prepare-commit-msg` and `post-commit`, and the amend runs `pre-commit`,
`prepare-commit-msg`, `commit-msg` and `post-commit`. With `-c
core.hooksPath=<an empty directory>` neither runs any hook.

Those three commands run through one helper in `internal/cli/carryforward.go`
that adds `-c core.hooksPath=<dir>` for that invocation only. `<dir>` is an
empty directory Roundfix creates and removes. No Git configuration file is
written, and the function signatures of `createCarryForwardStaging` and
`stageCarryForwardCandidate` do not change. The fast-forward merge into the
checkout in `applyCarryForwards`, and every read-only command, keep
`reconcileGitRaw`. `docs/user-guide/commands.md` (the `reconcile
--carry-forward` section) and the Roundfix skill's carry-forward paragraph
state that staging commits run without repository hooks and why.

## The wrap-fragile phrase check

`internal/speccheck/verification.go` adds `CodeVerifyWrapFragile =
"SC-VERIFY-WRAP-FRAGILE"` (`SeverityError`) and `WrapFragileVerification(task
spec.Task) []Finding`. For each top-level command of each Verification line
(split on `&&`, `||`, `;` and `|`, with a leading `!` allowed), a `grep`
invocation (optionally after `rtk`) is reported when all three hold:

- its pattern contains whitespace, where the pattern is the value of the first
  `-e` or else the first operand that is not an option;
- the pattern does not start with `^`, `#` or `|` and does not end with `$`,
  since those name a structure that cannot wrap;
- at least one file operand ends in `.md`.

A grep that reads standard input is never reported. The summary names the
phrase and the file. The fix gives the wrap-tolerant form with both filled in:
`tr -s '[:space:]' ' ' < <file> | grep -qF -- "<phrase>" || { printf 'missing
phrase in %s: %s\n' <file> "<phrase>" >&2; exit 1; }` for presence, or `! { tr
-s '[:space:]' ' ' < <file> | grep -qF -- "<phrase>"; }` for a `!`-prefixed
grep.

`detectTaskCoverageAndContextReferences` in `internal/speccheck/citations.go`
calls it for pending non-QA Tasks beside `InvertedExitVerification`, and lists
the code as skipped when `_tasks.md` is absent. `stagedDetectors` in
`coherence.go` registers it at the Tasks stage. The corpus code list in
`internal/docscontract/corpus_test.go`, the golden and its pin in
`internal/spec/archive_layout_characterization_test.go` carry it with count
`0`. The write-tasks Task template's Verification guidance teaches the form and
names the code, `skills/baseline_skill_contract_test.go` pins that guidance, and
`CONTEXT.md` defines **Wrap-Fragile Phrase Check**.

## API Contracts

1. `spec.Collisions` never reports a collision whose only shared path is an
   `instruction:` Context path.
2. An Implement Run's budget deadline is the time of its latest Task settlement
   before the deadline, or its start, plus `budget.max_run_duration`.
   `daemon.TaskCycleResult.BudgetDeadline` reports it.
3. A `BudgetExceeded` reason after a settlement reads `Run Budget exceeded:
   configured maximum <max>; elapsed <duration> since Task <id> settled.`
4. `delivery.RunOutcomeBudgetExceeded` parks an item with blocker
   `run-budget-exceeded` and its Run ID.
5. A delivery `RunSpec` result is decided from the implement exit status and
   the Run this invocation created; a Run that existed before it never decides
   it.
6. Carry-forward's staging `cherry-pick`, `cherry-pick --abort` and `commit
   --amend` run with `-c core.hooksPath=<empty directory>` and write no Git
   configuration.
7. `roundfix spec check` emits `SC-VERIFY-WRAP-FRAGILE` as an error naming the
   phrase and the file, with the wrap-tolerant form as its fix.

## Coverage Map

- Goal 1 → Instruction paths and Wave collisions; API Contract 1.
- Goal 2 → The renewing Implement Run Budget; API Contracts 2-3.
- Goal 3 → A budget stop parks as a Run outcome; API Contracts 4-5.
- Goal 4 → Carry-forward without repository hooks; API Contract 6.
- Goal 5 → The wrap-fragile phrase check; API Contract 7.
- Core Feature 1 → Instruction paths and Wave collisions.
- Core Feature 2 → The renewing Implement Run Budget.
- Core Feature 3 → A budget stop parks as a Run outcome.
- Core Feature 4 → Carry-forward without repository hooks.
- Core Feature 5 → The wrap-fragile phrase check.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 4.
- Success Metric 5 → Testing Approach 5.
- API Contract 1 → Instruction paths and Wave collisions.
- API Contracts 2-3 → The renewing Implement Run Budget.
- API Contracts 4-5 → A budget stop parks as a Run outcome.
- API Contract 6 → Carry-forward without repository hooks.
- API Contract 7 → The wrap-fragile phrase check.

## Integration Points

- **Spec 0175.** Its Tasks edit `internal/cli/carryforward.go` (moving staging
  creation into `internal/worktree/staging.go`), `internal/delivery/engine.go`
  and `internal/cli/deliver_workflow.go`. It is ahead of this Spec in the
  queue, so this Spec's Tasks build on whatever 0175 merged. The functions they
  name (`stageCarryForwardCandidate`, `abortCarryForwardStaging`,
  `runCandidate`, `RunSpec`) survive that Spec.
- **Spec 0173.** Its `roundfix deliver retry` is the recovery a
  `run-budget-exceeded` item uses.
- **Spec 0170.** Added the Task-stage detector pattern, the corpus golden entry
  and the skill contract pins this Spec follows for `SC-VERIFY-WRAP-FRAGILE`.

## Testing Approach

1. **Collisions.** Two unordered Tasks sharing only an `instruction:` path
   produce no collision in `spec.Collisions` or in `roundfix spec check`. The
   same pair sharing an `interface:` path still collides, and a path named as an
   instruction by one Task and as a Verification operand by both still
   collides.
2. **Budget.** A fake `Now` advances 0.9 allowances per Task through a serial
   graph and its QA gate, which ends Clean. A settlement after the deadline is
   kept but does not renew, and the Run ends `BudgetExceeded`. A stalled Task is
   cancelled one real allowance after the last settlement. A failed settlement
   renews. `BudgetDeadline` is reported, and is zero when the budget is disabled.
   The CLI ends Clean when the total exceeds one allowance, runs integration
   under the renewed deadline, and ends `BudgetExceeded` after it. The existing
   daemon, CLI and watch budget tests stay green unchanged.
3. **Delivery.** The engine parks `run-budget-exceeded` with the Run ID, later
   items continue, a retry carries forward from that Run, and an executor error
   stays `delivery-error`. The CLI mapping covers each input shape, including a
   stale Run.
4. **Carry-forward.** A fixture repository with refusing `pre-commit` and
   `commit-msg` hooks and marker-writing `prepare-commit-msg` and `post-commit`
   hooks proves and applies a carry-forward with no marker written. The
   checkout's `core.hooksPath` is unchanged afterwards, and the existing
   conflict and operational-failure tests stay green.
5. **Phrase check.** The line-bound form is reported with its phrase, file and
   fix; the wrap-tolerant form, an anchored pattern, a heading, a single word,
   a stdin grep and a `.json` operand are not; a completed Task is not; `--strict`
   keeps it an error; the corpus golden records `0`.
6. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. Instruction paths and Wave collisions (depends on: none).
2. The renewing Implement Run Budget (depends on: none).
3. A budget stop parks as a Run outcome (depends on: 2).
4. Carry-forward without repository hooks (depends on: 3).
5. The wrap-fragile phrase check (depends on: none).
6. Terminal QA (depends on: 1, 2, 3, 4, 5).

## Risks & Considerations

- **Shared guides.** Tasks 2, 3 and 4 each edit `docs/user-guide/commands.md`
  and the Roundfix skill, so they form a chain. Tasks 1 and 5 share no file
  with them or with each other and run in the first Wave beside Task 2. No
  non-QA Task declares an `instruction:` path, so the binary that runs this
  Spec, which still counts instruction paths, sees no collision.
- **A serial tail under today's budget.** This Spec's own chain of three Tasks
  and its QA gate runs under the start-based budget it replaces. A
  `BudgetExceeded` stop recovers through `roundfix deliver retry`.
- **Watchdog ownership.** The watchdog goroutine is owned by `TaskCycle` and
  stops on every return path; a named test proves it has stopped once
  `TaskCycle` returns.
