---
spec: 0177-runs-that-fit-their-budget-and-park-honestly
status: active
created: 2026-09-28
surfaces: [backend, cli, docs]
---

# Runs that fit their budget and park honestly

On 2026-09-28 three Specs of the same wave lost a Run, or were parked, for
reasons that had nothing to do with their work:

- Spec 0176 (Run `run_20260928T111930Z_13d79b0227e06da5`) and Spec 0173 (Run
  `run_20260928T125114Z_30db2f43ddd8a501`) reached `BudgetExceeded` with every
  implementation Task completed and only the QA Task, or one Task and the QA
  Task, still pending. `budget.max_run_duration: 2h` bounds the whole Implement
  Run from its start, and both graphs were serial chains. Each stop cost a
  relaunch and a fresh QA session.
- The graphs were chains partly because of the Wave-collision rule. It counts
  every `## Context` path a Task declares, including the read-only
  `instruction:` path. Every Task that names
  `.agents/skills/implement-task/SKILL.md` therefore collides with every other
  Task, and `SC-WAVE-COLLISION` forces `needs` edges between Tasks that share
  nothing they edit.
- The delivery owner parked Spec 0175 as `delivery-error: run Implement
  executor: roundfix implement failed with exit code 1 ...` when its Run
  (`run_20260928T154312Z_63679c6540b34603`) reached `BudgetExceeded`.
  `RunSpec` in `internal/cli/deliver_workflow.go` maps only an `Unresolved` Run
  to a Run outcome. `runCandidate` in `internal/delivery/engine.go` therefore
  wraps every other non-zero exit as an executor error, and a budget stop reads
  as an infrastructure failure. `roundfix deliver retry` (Spec 0173) recovered
  the item.
- Spec 0173's first Run (`run_20260928T111838Z_19370d95d9647ce9`) lost two
  correct Tasks after two Verification attempts each. Their Verifications ran
  `grep -q "<phrase of four to six words>" <guide>`, the Agent wrapped the
  phrase across two lines of `docs/user-guide/commands.md` and `CONTEXT.md`,
  and the diagnostic artifacts were empty, so nobody could see what had missed.

A fourth defect was recorded on 2026-09-25 (B10). Carry-forward cherry-picks
and amends in a fresh staging worktree that has no installed dependencies, so a
repository's Node-based commit hooks fail with `ERR_MODULE_NOT_FOUND`. In
Fluxus, completed Tasks were redone at a cost of 13 agent-minutes.

## Project Constraints

- Identifier strategy: applicable — the only new identifiers are the Delivery
  Queue blocker `run-budget-exceeded`, the Run outcome value `budget-exceeded`
  on the delivery boundary and the Spec Consistency Check code
  `SC-VERIFY-WRAP-FRAGILE`. Each follows the existing kebab-case form of its
  family. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the Run
  Database only; no credential and no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0158 settles a budget-expired
  Implement Run `BudgetExceeded`, cancels its Agent Sessions and keeps it
  recoverable through carry-forward; ADR-0137 says changing where the budget is
  evaluated needs its own evidence, which the three measured Runs now supply;
  ADR-0014 and ADR-0057 keep Verification and Task status Daemon-owned;
  ADR-0025 schedules ready Tasks concurrently and ADR-0056 separates Task and
  Verification capacity, which the collision rule gates; ADR-0020 lets a parsed
  result outrank a process exit code, the precedent for reading the Run outcome
  instead of the implement exit status; ADR-0113 derives a Run outcome from
  unresolved work, not from a failed Batch; ADR-0053 keeps terminal Run
  reconciliation proof-based; ADR-0148 keeps one Verification prober for
  authoring and dispatch; ADR-0133 makes a diagnostic name the literal it
  requires; ADR-0125 forbids executing a file a test process may still hold
  open. ADR-0093 checks Spec consistency by citation, ADR-0104 accepts on
  evidence a Spec did not author, ADR-0117 checks a defect at the stage that
  can produce it, ADR-0130 keeps a path governed once bounded, ADR-0155 makes
  the `qa` Task declare the matrix and ADR-0156 makes a declared promise name a
  consuming Task. ADR-0038, ADR-0127, ADR-0159 and ADR-0160 cite listed ADRs
  but govern the Verification repair, process residue, independent
  Verification and the red repository gate, which this Spec does not touch, so
  they do not apply. This Spec adds ADR-0164.
  This Spec's gate is bound by ADR-0080, ADR-0091, ADR-0096 and ADR-0097. All
  hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer authorized continuing with the
  next wave after the v0.18.0 release in chat on 2026-09-28 ("Após o release,
  pode continuar com as implementações da onda seguinte"); the Go test files
  ride the standing grant of 2026-09-21 for governed source and the skill files
  ride the standing grant of 2026-09-18 for keeping the shipped skills true to
  the CLI, recorded in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `.agents/skills/write-tasks/references/task-template.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `internal/speccheck/coherence.go`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `skills/baseline_skill_contract_test.go`. Sanctioned regeneration:
  `make skills-sync`. `.roundfixrc.yml` is not touched: its
  `budget.max_run_duration: 2h` keeps its value. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A Task's read-only instruction file never forces two Tasks into a chain.
- A Run whose Tasks each settle within the budget is never stopped by the
  budget, however long the graph, while a stalled Run still ends within one
  allowance.
- A Run that ends `BudgetExceeded` under the Delivery Queue parks as a Run
  outcome with its Run ID, never as a delivery error.
- Carry-forward succeeds in a repository whose commit hooks cannot run in the
  staging worktree.
- A Verification phrase check survives Markdown wrapping, and a failing one
  names the phrase and the file it missed.

## Core Features

1. **Instruction paths never collide.** The Wave-collision rule stops counting
   `instruction:` Context paths as touches. A Task still declares every path it
   edits under `interface:` or `creates:`, and those, Verification operands and
   prior-Run settlement paths still collide. `roundfix spec check`
   (`SC-WAVE-COLLISION`) and the Daemon's pre-dispatch refusal share
   `spec.Collisions`, so both change together.
2. **The Implement Run Budget renews at each Task settlement.** An Implement
   Run's deadline is the time of its most recent Task settlement, or its start
   when no Task has settled, plus `budget.max_run_duration`. A settlement
   observed after the deadline passed does not renew it. The QA gate starts
   with the allowance its last dependency's settlement renewed. The work that
   follows the cycle, which is integration, push and cleanup, runs under the
   renewed deadline. When the deadline passes, the Run still settles
   `BudgetExceeded` under ADR-0158. Its reason names the configured maximum and
   the time elapsed since the settlement that renewed it. A watch Run keeps its
   bound from the Run's start.
3. **A budget stop parks as a Run outcome.** When `roundfix implement` exits
   `1` and the newest Implement Run of the Spec is a Run this invocation
   created and that Run ended `BudgetExceeded`, the Delivery Queue parks the
   item `run-budget-exceeded` with that Run's ID. `roundfix deliver retry`
   resumes it through carry-forward like a `run-unresolved` item. An executor
   failure keeps its `delivery-error` blocker, and a Run that existed before the
   invocation never decides the outcome.
4. **Carry-forward stages without repository hooks.** The Git commands
   carry-forward runs inside its staging worktree that can run a repository
   hook (`cherry-pick`, `cherry-pick --abort` and `commit --amend`) run with
   hooks disabled for that invocation only. The carried commits already passed
   the Daemon's Verification and the repository's commit hooks when the Daemon
   settled them in their Run Worktree. Carry-forward then changes only the
   carried Task file's provenance. The checkout's Git configuration is never
   written.
5. **Phrase checks survive wrapping and name what they missed.** `roundfix spec
   check` reports `SC-VERIFY-WRAP-FRAGILE` for a pending non-QA Task whose
   Verification runs `grep` with a multi-word, unanchored pattern against a
   Markdown file. The finding names the phrase and the file, and its fix gives
   the wrap-tolerant form, which prints the phrase and the file when it misses.
   The write-tasks Task template teaches that form.

## Non-Goals / Out of Scope

- Changing the value of `budget.max_run_duration` in `.roundfixrc.yml`, adding
  a configuration key, or changing how a watch Run's budget is evaluated.
- Holding the QA gate back when the remaining budget looks short. A renewed
  allowance makes that estimate unnecessary, and it could only end a Run early.
- Parking `Stopped`, `TimedOut` or `Failed` Runs under new blockers; no
  measured case asks for it, and `Failed` is an infrastructure failure.
- Changing the Daemon's own Task commits, which keep running repository hooks
  and keep the `hook_refused` classification.
- Rewriting archived Specs whose Verifications use the line-bound form.

## Success Metrics

1. Two Tasks with no `needs` edge that share only an `instruction:` path
   produce no `SC-WAVE-COLLISION`, and two that share an `interface:` path
   still produce one.
2. A serial Task Graph whose total duration is 1.8 times
   `budget.max_run_duration`, and in which every Task settles within one
   allowance of the previous settlement, ends `Clean`. A Task that runs past one
   allowance since the last settlement ends the Run `BudgetExceeded` with every
   earlier settlement kept.
3. An Implement Run that ends `BudgetExceeded` under the Delivery Queue leaves
   its item `parked` with blocker `run-budget-exceeded` and a non-empty Run ID,
   and an executor failure still leaves `delivery-error: run Implement
   executor: ...`.
4. Carry-forward of a proved Task succeeds in a repository whose `pre-commit`
   and `commit-msg` hooks exit `1`, and none of the `pre-commit`,
   `prepare-commit-msg`, `commit-msg` or `post-commit` hooks runs during
   staging.
5. `roundfix spec check` reports one `SC-VERIFY-WRAP-FRAGILE` for a Task whose
   Verification is `grep -q "records no QA row" docs/user-guide/commands.md`,
   and none for the wrap-tolerant form, an anchored pattern, a single word or a
   non-Markdown file, while the active corpus golden records `0` for the code.

## Recorded limits

- A Run that keeps settling Tasks can run longer than one allowance in total.
  That is the point of the renewal; the Run Window crossing report still
  compares one allowance with the time left before the cutoff.
- `git worktree add` for the staging worktree still runs a repository
  `post-checkout` hook. Spec 0175 moves that command into
  `internal/worktree/staging.go`, and no measured failure involves it.
- The self-reporting phrase form is taught and suggested, not required: the
  wrap-tolerant form without the message is not a finding.

## Decisions

- **Renew the budget at each settlement, not per critical path.** A deadline
  scaled by the length of the critical path would let a stall at the first Task
  run for the whole scaled budget. A renewing deadline ends a stall one
  allowance after the last progress, which is the purpose ADR-0158 gave the
  budget, and it never stops a serial graph whose Tasks each settle in time.
  Holding the QA gate back would still leave a Spec that needs a relaunch.
- **Any settlement renews.** A failed Task also ends an Agent Session and
  advances the Run, and each Task settles at most once per Run, so the renewal
  stays bounded by the size of the graph.
- **Read the Run outcome, not the exit code.** An exit status of `1` covers
  both `Unresolved` and `BudgetExceeded`. The Run record names which one
  happened, so the delivery boundary reads it, the way ADR-0020 prefers a
  parsed result to an exit code.
- **Bypass hooks in staging instead of classifying the refusal.** A classified
  refusal would still stop carry-forward on work that is already proved. The
  staging commits are Roundfix's own bookkeeping over commits whose hooks
  already ran, and the checkout receives them only by a fast-forward merge.
- **Detect the line-bound form at authoring.** The wrapped phrase is an
  authoring defect, so ADR-0117 puts the check where the author can fix it. A
  Daemon heuristic over shell output could not tell which member of a chain
  missed.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: an `interface:` collision that
disappears with the instruction path, a stalled Task that the renewal keeps
alive, a Run that existed before the invocation deciding a park, a hook bypass
that leaks into the checkout's configuration, or a finding on the wrap-tolerant
form would each pass a happy-path test.

The outside-evidence row rests on sources this Spec did not author:

- The installed Git itself. It is measured running `prepare-commit-msg` and
  `post-commit` for a `cherry-pick`, and `pre-commit`, `prepare-commit-msg`,
  `commit-msg` and `post-commit` for a `commit --amend`. With
  `-c core.hooksPath=<an empty directory>` it runs none of them. Git 2.54.0
  showed exactly this on 2026-09-28.
- The live Run Database records of the three measured Runs, read only.
- The Verifications of Spec 0173's `task_01` and `task_02` as authored in
  commit `f0faa780`, before they were rewritten to the wrap-tolerant form. The
  new detector must report them when they are copied into a disposable Spec
  Root.

## Research basis

The budget and phrase-check evidence was recorded as Backlog Entries on
2026-09-28, and the carry-forward defect on 2026-09-25 from the secondbrain
inbox
(`inbox/roundfix/2026-09-10-events-aborta-e-carry-forward-quebra-no-worktree.md`,
defect 2). The adopted sources are indexed in
[references/_index.md](references/_index.md). The Spec 0175 delivery-error park
has no Backlog Entry. The Delivery Queue item record and Run
`run_20260928T154312Z_63679c6540b34603` in the Run Database are its evidence.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
