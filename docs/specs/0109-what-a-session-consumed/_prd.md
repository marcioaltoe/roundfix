---
spec: 0109-what-a-session-consumed
status: active
created: 2026-08-12
surfaces: [backend, data]
---

# What a Session consumed

Roundfix persists the effective Agent Selection for every Task — runtime, model,
and reasoning effort — and records nothing about what that Session consumed.
Concurrency and reasoning effort are different knobs, and a provider's daily total
cannot separate them: concurrency changes how fast a quota window empties, while
reasoning effort changes what one unit of work costs. When a quota was exhausted
mid-Spec on 2026-08-08, the question "was it concurrency or reasoning?" could only
be answered by inference over billing aggregates, and the conclusion remains
unverified because the row linking a Task to its cost does not exist. Four
recurring questions are queries once it does, and one of them — how much of a Spec
was spent on Tasks that failed and were redone — is invisible today and is not
small.

## Project Constraints

- Identifier strategy: applicable — Agent Session, Agent Selection, Preferred
  Selection, Fallback Selection, Task Type, and Batch are glossary terms this
  Spec attaches measurements to, and consumption is new vocabulary the glossary
  must own. The closing node checks whether the work introduced or changed a term.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. Consumption is read from what the existing adapter
  session already reports, and persisted locally. Source:
  `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0111 makes an unobserved Verification
  unknown rather than a verdict, and this Spec applies the same rule to
  measurement: a consumption an adapter did not report stays observably absent and
  never becomes a zero. ADR-0099 makes retention accounting mechanical while
  classification is not, which bounds how long consumption rows are kept alongside
  the journal they accompany. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized. The work is production Go in the agent, daemon, and store
  packages plus their tests. Source: `docs/agents/agent-instructions.md`.

## Goals

1. Cost per Task is attributable to the selection that ran it.
2. A missing measurement is visibly missing, never a zero.
3. The cost of work that failed and was redone becomes visible.
4. Concurrency and reasoning effort can be told apart from recorded data.

## User Stories

1. As a Supervisor choosing a reasoning effort, I want the cost of the same model
   at different efforts, so that the choice is made on measurement rather than on
   the belief that more effort is worth it.
2. As a Supervisor sizing concurrency, I want cost per Task Type, so that I can
   tell a quota emptied by volume from one emptied by cost per unit.
3. As a Supervisor reviewing a Spec that took several Runs, I want to see what was
   spent on Tasks that failed and were redone, so that the cost of rework is a
   number rather than an impression.
4. As a maintainer reading a consumption record, I want an unreported measurement
   to be visibly absent, so that I do not read "spent nothing" where the adapter
   simply said nothing.

## Core Features

1. **Consumption recorded beside identity.** Each Agent Session's consumption is
   persisted alongside the selection identity already recorded: input and output
   quantities, with cached reads separated when the adapter distinguishes them.
2. **Attribution to the work.** Each record carries the owning Task, its Task
   Type, its Batch, whether the session ran a preferred or a fallback selection,
   and whether it was the initial turn or a verification retry.
3. **Absence that reads as absence.** A measurement the adapter did not report is
   stored and displayed as absent, never defaulted to zero.
4. **The four questions are queries.** Cost per Task Type, cost of one model at
   different reasoning efforts, cost of verification retries, and cost of Tasks
   that failed and were redone are each answerable from the recorded data without
   reconstructing them from billing.

## Non-Goals / Out of Scope

- Pricing, currency, or cost estimation. This Spec records quantities the adapter
  reports; converting them to money is a separate concern with its own inputs.
- Enforcing a budget, throttling, or refusing a Run on consumption.
- Any change to selection, fallback, or dispatch behavior.
- Reconstructing consumption for Sessions that already ran.
- A reporting surface beyond what makes the four questions answerable.

## Success Metrics

- Every completed Agent Session in a Run carries a consumption record or an
  explicit absence, with no record defaulting to zero.
- The four named questions are each answered by a query over recorded data.
- The cost of rework in a multi-Run Spec is reported as a number, against a Spec
  whose four Runs on 2026-08-08 included three that failed on task-authoring
  defects, one discarding a Task Worktree with twenty-two finished files. Source:
  a session in a repository this Spec did not build.
- A runtime whose adapter reports no consumption produces records that read as
  unmeasured rather than as free.

## Decisions

- Quantities are recorded as the adapter reports them, without normalization
  across runtimes, because a normalized unit invented here would be a number
  nobody can trace back to a provider's own accounting.

## Open Questions

- Whether the supported ACP adapters expose consumption at all. This is settled
  before any record shape is fixed; if none does, the Spec's deliverable becomes
  the explicit-absence path plus the finding that the data is unavailable.
- Whether consumption is a new event on the existing Run Event Stream or a
  separate record keyed to the session. The default until answered is the existing
  stream, since it already carries per-session identity and retention.
- Whether consumption rows follow the journal's retention or outlive it. The
  default is to follow it, because a cost with no surrounding Run context answers
  none of the four questions.
