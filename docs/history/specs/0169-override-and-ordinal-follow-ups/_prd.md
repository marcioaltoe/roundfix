---
spec: 0169-override-and-ordinal-follow-ups
status: archived
created: 2026-09-25
surfaces: [backend, cli, docs]
archived: "2026-09-25"
source_slug: 0169-override-and-ordinal-follow-ups
---


# Override and ordinal follow-ups

Spec 0159 shipped the QA Archive Override and the claimed-ordinal check with
five minor gaps recorded by its second pre-PR review:

- An override of a Spec whose QA Task is failed or pending while its newest
  report says `pass` stamps `qa_override_qa_outcome: pass` and records nowhere
  that the QA Task was not completed, so the archived record reads as an
  override of a passing QA.
- The Roundfix skill still says the override "refuses when QA already
  qualifies", the rule before the fix, so an agent would not try it in exactly
  the case the fix unlocked.
- The archive-spec skill's Steps still describe stamping `qa_override: true` by
  hand, which skips the command's refusals and provenance.
- Two Tasks of one Spec creating different paths with the same ADR number are
  not reported until one file exists.
- A Spec without `_tasks.md` does not list `SC-ORDINAL-CLAIMED` among the
  detectors it skipped.

## Project Constraints

- Identifier strategy: applicable — ADR ordinals stay claimed identifiers; a
  claim is now also unique within one Spec. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0154 keeps an authorized QA Archive
  Override archive-only and truthful about the QA it waives, ADR-0093 checks
  Spec consistency by citation, ADR-0104 accepts on evidence a Spec did not
  author, ADR-0130 keeps a path governed once bounded, ADR-0155 makes the `qa`
  Task declare the matrix and ADR-0156 makes a declared promise name a consuming
  Task. This Spec's gate is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097 and
  ADR-0117. All hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer's express authorization of
  2026-09-24 for the archive override guidance and the standing grant of
  2026-09-21 for governed source, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/spec/archive.go`, `internal/spec/archive_test.go`,
  `.agents/skills/archive-spec/SKILL.md`, `skills/archive-spec/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- An override's record says exactly what was waived.
- The skills send agents through the command, with its current rule.
- An ADR number cannot be claimed twice, even inside one Spec.

## Core Features

1. **The stamp names the waived Task.** When the QA Task is not completed, the
   override stamps its status beside the report outcome.
2. **Guidance matches the command.** The Roundfix and archive-spec skills state
   that the override is refused only when a normal archive would succeed, and
   that it is performed only through `roundfix archive --qa-override`.
3. **Ordinals unique within a Spec, and a visible skip.** Two Tasks of one Spec
   claiming one number with different paths are reported, and a Spec without a
   Task Graph lists `SC-ORDINAL-CLAIMED` as skipped.

## Non-Goals / Out of Scope

- Changing when the override is accepted or refused.

## Success Metrics

1. An override of a failed QA Task with a `pass` report records the QA Task's
   status.
2. Both skills and the guide describe the current refusal rule and the command
   as the only path.
3. A same-Spec duplicate ordinal is reported at authoring, and the skip is
   listed for a Spec without `_tasks.md`.

## Recorded limits

- The mandatory rule in `docs/agents/docs-layout.md`, generated from the Baseline
  module `spec-workflow.json`, still describes stamping `qa_override: true` by
  hand. It is Baseline-owned and outside this Spec's authority; found by the
  second pre-PR review of 2026-09-25 and carried to Backlog Entry
  `docs/backlog/2026-09-25-docs-layout-still-describes-hand-stamped-overrides.md`.

## Decisions

- **Record the waived fact, not a summary.** ADR-0154 asks for a truthful
  record; the QA Task status is the fact the override waives.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: a stamp that still reads as a
passing QA, or a same-Spec duplicate that passes, would pass the happy path.

## Research basis

The five gaps and their reproductions are recorded in the archived PRD of Spec
0159, from its second pre-PR review of 2026-09-24.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
