---
spec: 0149-a-supported-way-to-reopen-a-settled-gate
status: active
created: 2026-09-19
surfaces: [backend, cli, docs]
---

# A supported way to reopen a settled gate

A Spec's terminal QA Task is a gate: it depends on every leaf of the Task Graph,
and when it settles `completed` the recorded verdict describes the graph as it
stood. If a Task below it later stops being completed — because a corrective
Task is added, or because a settled Task is reopened — the gate's result no
longer describes anything real, and the loader says so:

```
validate qa gate: QA gate result is invalidated for Task "task_02" because these dependencies are not completed: task_01
```

That refusal is correct. What is missing is the next step. Every command that
loads the graph — `archive`, `settle`, `implement` — refuses with that reason,
and no command resolves it. The Spec is not merely blocked; it is unreachable by
the tooling that owns it.

The only escape today is to open the QA Task file and write `status: pending` by
hand. That works, and it is exactly what a maintainer must never do: it makes a
stale gate disappear by editing the record of the gate, leaving nothing that
says a verdict was invalidated, which report it was, or why.

This Spec adds the missing operation. It is carved from Spec 0129 Core Feature
3, which asks for a supported stale-QA invalidation path with retained earlier
report and Result evidence.

## Project Constraints

- Identifier strategy: applicable — the command addresses a Spec by slug and a
  Task by its Task Graph id, and preserves both. It mints no new identity.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the command reads and writes local
  files only; it opens no network connection. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0155 makes the `qa` Task declare the
  gate's matrix and ADR-0104 makes a Spec accept on evidence it did not author;
  both hold here, and this Spec adds an operation under them rather than relaxing
  either. ADR-0156 makes each promise this Spec declares name a consuming Task,
  which its Success Metrics and API Contracts satisfy through the Task Graph. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the shipped Roundfix skill documents the CLI
  surface, so a new command reaches it. Express maintainer authorization:
  granted 2026-09-18 as a standing grant, consumed here and recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Every other path this Spec touches —
  `internal/cli`, `internal/spec`, `docs/user-guide` — is ordinary source that no
  authorization has bounded, measured through the governance probe rather than
  estimated. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- A stale QA gate has a supported operation that clears it.
- A healthy QA gate cannot be cleared by that operation.
- Nothing that recorded the prior verdict is deleted or overwritten.

## Core Features

1. **A reopen that a stale gate accepts.** `roundfix reopen --spec <slug>`
   returns the Spec's terminal QA Task to `pending` when its recorded verdict is
   invalidated by a dependency that is not completed. After it runs, the
   commands that refused on the stale gate load the Spec again.
2. **A reopen that a healthy gate refuses.** When the QA Task is not completed
   there is nothing to reopen, and when every dependency is completed the
   verdict still describes the graph. Both refuse in preflight, change no file,
   and say which condition failed.
3. **The prior verdict is kept, not erased.** The QA Report that recorded the
   invalidated verdict is left byte-identical, and the QA Task keeps its prior
   `## Result`. The reopen is recorded alongside them — which report it
   invalidated, which dependencies made it stale, and when.
4. **No Run, no push, no side effects.** The command creates no Run, writes no
   Run Event Journal entry, commits nothing and pushes nothing, matching the
   envelope `settle` already established.

## Non-Goals / Out of Scope

- Re-running the QA gate. `reopen` restores the Spec to a state where a Run can
  execute the gate again; it does not execute it.
- Deciding when a gate *should* be reopened. The command proves the gate is
  stale; whether corrective work was the right call stays a maintainer's
  judgement.
- The premise-falsification and supersession amendment of Spec 0129 Core Feature
  2, the coverage extension of Core Feature 1, and the temporal-prerequisite
  representation of Core Feature 4. They stay there.
- Any change to when the loader raises the stale-gate refusal. This Spec adds the
  operation that answers it and leaves the detection exactly as it is.

## Success Metrics

1. A fixture Spec whose QA Task is `completed` above a `pending` dependency:
   `reopen` exits 0, and a command that previously refused with the stale-gate
   reason loads the same fixture afterwards.
2. A fixture whose QA Task is `completed` above completed dependencies: `reopen`
   refuses, and every file under the Spec directory is byte-identical to what it
   was before the command ran.
3. After a successful reopen, the QA Report file is byte-identical and the QA
   Task file still contains its prior `## Result` text.

## Decisions

- **Prove staleness, do not take it on request.** A command that reopens any
  gate on request is a supported way to discard a verdict. This one refuses
  unless the loader's own condition already holds, so the operation can only
  undo a state the tooling itself calls invalid.
- **Reuse the settle envelope.** `settle` already established what a narrow
  repair command may do: mutate one Task's status, touch no Run, write no
  journal entry, never push. A second command with the same envelope is one
  contract to learn, not two.
- **Record the invalidation in the Task, keep the report untouched.** The QA
  Report is the evidence of what was verified; rewriting it to say it no longer
  counts would corrupt the thing being preserved. The Task file is where the
  gate's state lives, so that is where its invalidation belongs.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases are the point: a reopen that accepts a healthy
gate, or that rewrites a report, would satisfy Core Feature 1 while defeating
the Spec.

## Research basis

The behavior was measured on this repository, not inferred. A fixture Spec whose
QA Task is completed above a pending dependency was loaded through `archive` and
`settle`; both refused with `validate qa gate: QA gate result is invalidated for
Task "task_02" because these dependencies are not completed: task_01`, and no
other command resolves it. The governance probe classified every path this Spec
touches; only the two shipped skill files are governed, and both fall under the
standing grant of 2026-09-18.

The premise is also first-hand: during Spec 0148's delivery a corrective Task was
needed under a settled gate, and the only available action was to hand-edit the
QA Task's status back to `pending` and disclose it in the commit. That is the
gap this Spec closes.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
