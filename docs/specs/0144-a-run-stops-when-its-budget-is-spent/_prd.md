---
spec: 0144-a-run-stops-when-its-budget-is-spent
status: active
created: 2026-09-18
surfaces: [backend, cli, docs]
---

# A Run stops when its budget is spent

Roundfix lets a repository set a maximum Run duration, and it enforces that
maximum in exactly one place: the Round watch loop, which derives a deadline
from the Run's start and settles the Run as `BudgetExceeded` when it passes.

An Implement Run reads the same setting and does nothing with it. At start it
compares the maximum against the Run Window and may print that the Run "may run
past it". After that nothing bounds the Run. An Agent Session that stalls, a
Verification that never returns, a runtime that stops answering — each of these
leaves a Run that runs until someone notices and stops it by hand.

The vocabulary for the honest end already exists: the Run Database carries
`BudgetExceeded`, produced today only by the watch loop. What is missing is the
bound itself, and a way back: Task Carry-Forward accepts only `Stopped` and
`Unresolved`, so a Run ended on its budget would strand every Task it had
already proved.

This Spec is the first slice carved from Spec 0127, which keeps the durable
orchestration, the queue preparation and the delivery order.

## Project Constraints

- Identifier strategy: applicable — the Run outcome vocabulary, Run identifiers and Run Branch names keep their spelling; this Spec produces an existing outcome in a new place and coins none. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — bounding a Run reads local configuration and cancels local Agent Sessions and processes; no credential store or network transport is involved. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — Run settlement, Task status and cancellation are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0057 applies: the Daemon owns Run and Task status, so the bound is the Daemon's to apply and record.
  ADR-0127 applies: process residue is a readiness observation, so a bounded Run proves the processes it cancelled rather than assuming them gone.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0104 applies: acceptance rests on evidence this Spec did not author, which here is the watch loop that already enforces this budget and the Run Database vocabulary that already names the outcome.
  ADR-0158 applies: a budget-expired Run settles `BudgetExceeded` and stays recoverable through carry-forward.
- Tooling authority: applicable — the Run outcome and the carry-forward contract are public CLI behavior, and the repository's hard rule ships the Roundfix skill update with such a change. Express maintainer authorization: granted 2026-09-18, recorded in [_authorization.md](_authorization.md); bounded files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned regeneration: `make skills-sync` regenerates the mirror, and `make baseline-digests` rewrites any derived pin the approved edit moves. The Daemon, the CLI and the Run Database are ordinary source that no authorization has bounded. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- A Run that exceeds its configured maximum duration ends, instead of running
  until a person notices.
- The end is named for its cause, so a reader can tell a spent budget from a
  requested stop or a failure.
- Work the Run already proved can be recovered without rerunning it.
- A repository that has not enabled the budget sees no change at all.

## User Stories

1. As a maintainer who enabled the budget, I want a stalled Run to end at its
   maximum, so that a stuck Agent Session costs one budget instead of a night.
2. As a maintainer reading a Run's outcome, I want a spent budget named as such,
   so that I do not read it as a failure of the work or as someone's stop.
3. As a maintainer whose Run ended on its budget after proving three Tasks, I
   want to carry those Tasks forward, so that the bound costs the stall and not
   the work.
4. As a maintainer who has not enabled the budget, I want Runs to behave exactly
   as they do today.

## Core Features

1. **An Implement Run is bounded by its configured maximum.** When the budget is
   enabled, the Daemon derives the Run's deadline from its start, as the watch
   loop already does, and ends the Run when the deadline passes.
2. **The end cancels what the Daemon owns.** Agent Sessions and processes the Run
   owns are cancelled through the existing cancellation path, and termination is
   proved rather than assumed.
3. **The outcome names the cause.** The Run settles `BudgetExceeded`, with a
   reason naming the configured maximum and the elapsed time. `Stopped` stays
   the Stop Command's outcome and `Unresolved` stays what it is today.
4. **Proved work stays recoverable.** The Run Worktree and Run Branch are
   preserved, completed Tasks keep their status, and Task Carry-Forward accepts a
   `BudgetExceeded` Run beside `Stopped` and `Unresolved`, under the proof
   requirements it already applies.
5. **A disabled budget changes nothing.** With the budget disabled, no deadline
   is derived, no outcome changes, and the existing Run Window report keeps its
   wording.

## User Experience

A maintainer enables the budget and starts a Run. A stalled Agent Session no
longer runs past the maximum: the Run ends, its outcome reads `BudgetExceeded`,
and the reason names the maximum and how long the Run took. The Run Worktree and
Branch are still there, and `reconcile --carry-forward` hands back the Tasks the
Run proved.

## Non-Goals / Out of Scope

- Durable orchestration across Runs, queue preparation, delivery order, or
  progress records that survive a lost chat session. Spec 0127 keeps those.
- Changing what the Run Window means or when it refuses a Run start.
- Bounding an individual Task, a Verification command, or an Agent turn. The
  budget bounds the Run.
- Changing the watch loop, its rounds, or its own budget handling.
- Making the budget enabled by default, or changing its default value.
- Recovering a Run automatically; carry-forward stays an explicit act.

## Declared intentional breaks

- A repository with the budget enabled will see Runs end at the maximum that
  previously continued. That is the feature, and it is visible the first time a
  Run exceeds the configured duration.
- Carry-forward accepts one more terminal outcome than the shipped skill
  currently documents; the skill ships with this change.

## Regression locks

- With the budget disabled, Run settlement, outcomes and reports are unchanged.
- `Stopped` remains the Stop Command's outcome and keeps its owner-identity
  proof; `Unresolved` keeps its meaning.
- Carry-forward keeps refusing every outcome outside the accepted set, and keeps
  its existing proof requirements and its refusal of a partially provable set.
- The Run Window report keeps its current wording and timing.

## Acceptance evidence

At least one acceptance row rests on evidence this Spec did not author: the
watch loop, which already derives a deadline from a Run's start and settles
`BudgetExceeded`, and the Run Database vocabulary, which already carries that
outcome. The replay reads both and shows that this Spec applies the same rule at
the Implement path rather than inventing one. Where either source cannot be
read, the row records that reason and does not block.

## Success Metrics

1. A fixture Run with the budget enabled and a maximum shorter than its work
   settles `BudgetExceeded`, with a reason naming the maximum and the elapsed
   time.
2. A fixture Run with the budget disabled settles exactly as it does today.
3. A `BudgetExceeded` fixture Run whose earlier Tasks completed is accepted by
   carry-forward, and one whose Tasks cannot be proved is still refused.

## Decisions

- **The existing outcome, in a new place.** `BudgetExceeded` already names this
  cause; inventing another word would split one meaning in two. See ADR-0158.
- **Not `Stopped`.** That is the Stop Command's outcome, which proves owner
  identity and cancels at a person's request.
- **Recoverable by carry-forward.** A bound that stranded proved work would cost
  more than the stall it ended. See ADR-0158.
- **Slice, not portfolio.** Spec 0127 carries eleven Core Features; this Spec
  takes the Run bound and leaves the durable owner there.

## Research basis

**Secondbrain.** Consulted before authoring, index first, then a query on
bounding long-running automated jobs and on distinguishing a timeout from a
failure. The results were dominated by mirrors of this repository, which are
references rather than independent knowledge, and no source changed the design.

The repository's pending Inbox Entries were read before authoring; none is a
source for this Spec.

**Exa MCP.** Consultation was attempted, and no Exa MCP tool was available in
this session. No external source was read, so no external validation is claimed.

**Local measurement.** The watch loop's deadline and budget check, the Implement
path's start-time Run Window comparison, the absence of any deadline in the
Daemon, the Run Database outcome vocabulary and carry-forward's accepted
outcomes were each read in this repository's source before this PRD was written.

## Open Questions

None.
