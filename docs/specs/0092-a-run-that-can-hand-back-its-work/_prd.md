---
spec: 0092-a-run-that-can-hand-back-its-work
status: active
created: 2026-08-09
surfaces: [backend, cli]
---

# A Run that can hand back its work

A Run creates Tasks, Batches, Agent Sessions, Run Branches and worktrees. Each of
those needs exactly one terminal disposition, and today several of them are
classified by a surface that did not create them. Four measured consequences,
all from the same shape.

**A session that never opened is read as work already started.** On 2026-08-08
the Codex quota was exhausted. The adapter printed its usage limit and exited;
the whole Run lasted nineteen seconds. Roundfix had already emitted
`AGENT_WORK_STARTED`, so the exit was classified as a Batch failure after work
began, and the configured `codex → claude` Fallback Chain became ineligible. The
guard is right in intent — never swap models over partially finished work — and
wrong at the boundary. A Spec stopped for four days on a quota that a proven
alternative could have absorbed in seconds.

**A Batch that fails on one issue erases the others.** `MarkBatchFailed`
overwrites every Review Issue in the Batch. On pull request 143 round 001 that
turned twenty recorded resolutions into `failed`; each carries triage notes
describing the work it did, and a spot check confirmed one had landed in the
tree. The obvious repair — skip issues already in a Terminal status, as the
sibling `SettleAssignedIssues` already does — was attempted on 2026-08-09 and
reverted, because six tests encode the opposite contract: a Batch whose Agent
failed must leave its issues `failed` so the Run ends `Unresolved` and retries.
Preserving the outcomes makes the same Run reach `Clean`, reporting success for a
Run whose Agent crashed. The defect is not the marking; it is that the Run's
outcome is derived from the Batch's failure rather than from whether unresolved
work remains.

**A stopped Run discards the Tasks it proved.** Implementing Spec 0089
re-executed Task 01 in four separate Runs and Task 02 in three, against unchanged
inputs, because each Run stopped later in the graph and the checkout still read
`status: pending` for work that had been committed inside a Run Worktree.

**And then it blocks the next Run.** `reconcile` classifies those branches
`unintegrated` and preserves them; Branch Integrity Preflight then refuses to
create any Run while one exists. Three such branches refused two consecutive
`roundfix resolve` invocations on 2026-08-09 until they were removed with
`git branch -D`. The suggested `git merge --ff-only` could not apply, because the
branches had diverged behind the work that superseded them.

## Project Constraints

- Identifier strategy: not applicable — no new persisted entity. Runs, Batches,
  Tasks and Run Branches keep the identities they already carry.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Run lifecycle, git plumbing and
  local storage only. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0010 lets a Batch settle and the Run
  continue, which this Spec builds on rather than replaces. ADR-0050 keeps
  Fallback Chains inactive until Run creation and forbids preflight
  substitution; the eligibility boundary this Spec moves lives after that point
  and must not weaken it. ADR-0069 also cites ADR-0050 and does not apply here:
  it governs Baseline semantic analysis, which this Spec does not touch.
  ADR-0091 makes the QA gate a Task node of its own type, which is why this
  Spec's graph carries one. ADR-0096 has that gate prove machine facts before
  spending an Agent turn, and this Spec's gate follows it — every disposition
  row is a command run, not Agent work. ADR-0104 requires an acceptance row on
  evidence this Spec did not author. This Spec adds ADR-0113, ADR-0114 and
  ADR-0115.
  ADR-0117 places a check with the stage that can produce its defect; it does not change what this Spec delivers, and it moves where this Spec's gate rows run only once Spec 0093 ships. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the change lives in `internal/daemon`,
  `internal/agent`, `internal/rounds` and `internal/spec`, all ordinary source.
  Source: `docs/agents/agent-instructions.md`.

## Goals

1. A Fallback Chain activates when the selected runtime could not open a
   session, because opening a session is not Agent work.
2. A Batch that fails preserves the outcomes its Agent already recorded, and the
   Run still reports `Unresolved` while unresolved work remains.
3. A Run that stops can hand its settled Tasks back to the checkout, on evidence
   that their inputs have not moved.
4. A Run Branch that has been superseded has a disposition Roundfix can perform,
   so a stopped Run stops blocking every Run after it.
5. No Run reports `Clean` on work it did not finish.

## Core Features

- **A work-started boundary that means what it says.** The signal that makes a
  Fallback Selection ineligible moves to the first Agent turn that could have
  changed something, rather than to the attempt to open a session. A runtime
  that refuses to serve — quota, auth, unavailable adapter — is a selection
  failure, and the chain exists for exactly that.
- **Settlement separated from outcome.** A Batch records what its Agent
  achieved, including on a Batch that failed. The Run's outcome is derived from
  whether unresolved work remains, not from whether a Batch reported failure.
  Both statements can be true at once: the Agent crashed, and seventeen issues
  are resolved.
- **Settled Tasks carried forward on evidence.** A stopped Run's Tasks whose
  Verification passed against inputs that have not since changed can be handed
  to the checkout, with the Run Worktree's own commits as the record. The
  carry-forward is an explicit maintainer action, not a silent one.
- **A disposition for a superseded branch.** A Run Branch whose work is already
  present in the target — or whose Run was superseded — can be discarded through
  Roundfix, with the evidence of what it held. `reconcile` keeps its read-only
  contract; the disposition is a separate, named act.

## Non-Goals / Out of Scope

- Making `reconcile` integrate. Its read-only contract is deliberate and stays.
- Retrying a failed Task automatically. This Spec makes prior work reusable; it
  does not decide to re-run anything.
- Changing what a Fallback Chain contains or how it is configured. Only when it
  becomes eligible moves.
- A one-Run override that expresses a chain. The workaround gap recorded in the
  backlog is real and stays out of scope; this Spec removes the need to reach
  for the workaround.

## Decisions

- A Task is carried forward only when its Verification passed and its declared
  inputs are unchanged since it settled. Anything weaker risks reporting a Task
  complete against inputs it never saw, which is worse than re-executing it.
- Carrying forward is an explicit act, not automatic. The failure mode of
  getting it wrong is silent and expensive; the failure mode of asking is one
  command.
- Discarding a superseded Run Branch requires the evidence of what it held to be
  written first, so the act is auditable after the branch is gone.
