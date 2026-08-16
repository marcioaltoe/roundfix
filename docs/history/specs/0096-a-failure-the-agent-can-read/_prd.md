---
spec: 0096-a-failure-the-agent-can-read
status: archived
created: 2026-08-12
surfaces: [backend, docs]
archived: "2026-08-16"
source_slug: 0096-a-failure-the-agent-can-read
---


# A failure the Agent can read

A failing Verification can hand the Agent Session nothing at all. The command
pattern adopted across this repository redirects output to a file to keep the exit
status honest, and the check that then fails prints nothing — so the Daemon
captures empty output and the feedback prompt carries the command, the exit
status, and no cause. In one measured Task the Agent spent its repair turn
rewriting its own task file with a diagnosis it had deduced, which is reasonable
behavior given zero information. The same blindness repeats across Runs: a Task
that fails twice with an identical assertion is reported as new both times, and a
whole Run was spent reproducing a known diagnostic. And when a gate returns more
corrective work than the contract's ceiling allows, the contract says the
decomposition was wrong and does not say what to do, so the loop stops and asks a
human.

## Project Constraints

- Identifier strategy: applicable — Verification Feedback, Agent Session, Work
  Item, Run Event, and Task are glossary terms this Spec changes the reporting of,
  and a repeated-failure signal is vocabulary the glossary must own. The closing
  node checks whether the work introduced or changed a term. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. The work is diagnostic capture, event emission, and
  an authoring contract clause. Source: `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0014 gives the Daemon ownership of task
  verification and status settlement, which is the owner of every change here;
  ADR-0038 allows one Verification repair, which is the single turn this Spec
  stops wasting. The decisions extending that ownership are accounted: ADR-0020
  ranks a parsed prompt result above the runtime's exit code, ADR-0057 gives the
  Daemon exclusive ownership of Implement Task status, ADR-0056 separates Task
  Capacity from Verification Capacity, and ADR-0096 with ADR-0117 place the gate's
  mechanical stage and its checks — this Spec changes none of those. ADR-0111
  makes an unobserved Verification unknown rather than a verdict, which is exactly
  what an empty diagnostic must be reported as. ADR-0127 places process residue in the readiness diagnostic; this Spec changes what an Agent reads about its own failure and reports no residue, so it does not apply.
  This Spec's own decisions are ADR-0135, which makes an absent diagnostic a
  reported state, ADR-0136, which recognises a repeated failure by a normalised
  signature, and ADR-0137, which states what the run budget bounds where it is
  configured. ADR-0081 keeps a grant drawn around the authorized cause with its
  sanctioned regeneration following as fallout, which is how the authoring-skill
  edit and its mirror are accounted. The closing node rests on two more it follows
  rather than changes: ADR-0091 makes the QA gate a Task node of its own type, and
  ADR-0104 makes a Spec accept on evidence it did not author, which is what the
  2026-08-08 measurement supplies here. ADR-0129 lets a grant name its regeneration
  command while the tree names that command's outputs; this Spec's one authorized
  edit is a skill file whose mirror the record already declares, so it reads that
  decision rather than changing it. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the authoring contract gains the sanctioned
  exits for a reached corrective-Task ceiling. Express maintainer authorization:
  `docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`,
  granted 2026-08-12. Bounded files: `.agents/skills/write-tasks/SKILL.md`.
  Source: `docs/agents/agent-instructions.md`.

## Goals

1. A repair turn is never spent guessing at an empty diagnostic.
2. A failure that already happened is recognisable as a repetition.
3. A reached corrective-Task ceiling is a decision inside the loop's authority.

## User Stories

1. As an Agent given a failed Verification, I want an empty diagnostic reported as
   empty, so that I look for the output rather than deduce a cause.
2. As a Supervisor reading a second failed Run, I want a repeated identical
   failure named as repeated, so that I amend the task file instead of re-invoking
   and reproducing it.
3. As a Supervisor whose gate returned more corrective work than the ceiling
   allows, I want the sanctioned exits written down, so that the loop continues
   instead of stopping for a policy decision.

## Core Features

1. **An empty diagnostic says so.** A Verification that fails with no output has
   that stated explicitly in the feedback returned to the Agent Session, rather
   than an empty message, so the absence of a cause is itself information.
2. **A repeated failure is named.** When a failure's verdict and diagnostic
   signature match a prior failure of the same Work Item, that repetition is
   reported in the Run's own record and in the event stream, so the signal reaches
   a Supervisor without reading two files side by side.
3. **The ceiling names its exits.** The authoring contract states what to do when
   corrective work exceeds the ceiling: amend the technical spec and recut from
   it, or promote the excess to its own Spec with the gate failing the discovered
   story explicitly.
4. **The surface a task file came from is named.** Recovery output states which
   surface the task file was read from, so a fix is applied once rather than
   discovered by trial to be needed in two places.
5. **A vacuity refusal names what it means.** The event that reports a
   Verification refused as vacuous lists the commands that *passed* against the
   unchanged tree, and a reader takes that list for what the tool ran. It names
   the offenders explicitly, or reports every command with its own verdict, and
   points at the probe log that settles it.
6. **The run budget says what it bounds.** Whether the configured maximum Run
   duration is evaluated at Work Item boundaries is settled and stated, so a
   maintainer configuring a window knows what it promises.

## User Experience

Verification feedback that carries no diagnostic says the command produced none
and where its output went, if the command redirected it. A repeated failure is
labelled as a repetition of a named earlier attempt. Recovery output names its
surface on the same line it names the Task.

## Non-Goals / Out of Scope

- Changing how many repair turns a Task gets, or the ceiling's value.
- Parsing or interpreting Verification output beyond detecting its absence.
- Changing verdict semantics or how a Task settles.
- Rewriting the measured Verification commands in this repository; the authoring
  rule that prevents new ones belongs to the sibling authoring Spec.

## Success Metrics

- A Verification failing with no output produces feedback stating that, proven
  against the measured Task whose two diagnostic artifacts were zero bytes.
- A Task failing twice with an identical assertion is reported as repeated on the
  second failure, measured against a Spec in a repository this Spec did not build
  where that repetition cost a full Run on 2026-08-08.
- The corrective-ceiling exits are written and a graph reaching the ceiling
  proceeds under one of them without a human decision.
- Recovery output names the surface its task file came from.

## Decisions

- Absence of a diagnostic is reported as a distinct state rather than inferred by
  the Agent from an empty prompt, because an empty prompt and a prompt about
  nothing are indistinguishable to the reader who must act.

## Open Questions

- What constitutes a diagnostic signature for the repeated-failure check. A
  normalized form of the captured output is the default until answered, since raw
  equality would miss a failure whose only difference is a timestamp.
- Whether the run-budget question is answered by documentation or by changing
  where the budget is evaluated. Documentation is the default until the behavior
  is reproduced, because the single measured overrun has no established cause.
