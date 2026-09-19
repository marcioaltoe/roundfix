---
spec: 0151-a-supersession-the-archive-can-see
status: active
created: 2026-09-19
surfaces: [backend, cli, docs]
---

# A supersession the archive can see

Spec 0128 asked the release planner to accept bare stable tags. Every one of its
Core Features shipped — in Spec 0147, which was authored as that whole content
rather than as a slice of it. Nothing in 0128 remains to build.

It is still in `docs/specs/`, active, because nothing can move it. Measured:

```
$ roundfix archive 0128-release-planning-with-bare-stable-tags
Reason:
  Task Graph manifest ".../_tasks.md": file does not exist; run the
  write-tasks workflow to create the Task Graph
```

Archive proves a Spec is finished by reading completed Tasks and a passing gate.
A Spec whose content another Spec delivered was never decomposed, so it has
neither, and it cannot acquire them: writing a Task Graph for work that is
already shipped would be a fiction, and running a gate over it would prove
nothing.

Today the only record is a `## Disposition` paragraph someone wrote by hand into
the PRD. It is accurate and it is invisible: no command reads it, no lifecycle
state reflects it, and the Spec keeps appearing as active work.

This Spec is carved from Spec 0129 Core Feature 2, which asks for premise
falsification and supersession recorded as an amendment with prior Results and
evidence preserved.

## Project Constraints

- Identifier strategy: applicable — an amendment names both Specs by slug and
  preserves their identities. Nothing is renumbered and no identity is minted.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local file reads and writes only.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0155 makes the `qa` Task declare the
  gate's matrix, ADR-0104 makes a Spec accept on evidence it did not author, and
  ADR-0156 makes each declared promise name a consuming Task. All three hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — a new command changes the public CLI surface,
  which the shipped skill documents. Express maintainer authorization: granted
  2026-09-18 as a standing grant, consumed here and recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. `internal/cli`, `internal/spec` and
  `docs/user-guide` are ordinary source that no authorization has bounded,
  measured through the governance probe. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## Goals

- A Spec whose content another Spec delivered can leave the active queue.
- The record of why it left is written where a command can read it.
- Nothing that already existed is deleted to make room for that record.

## Core Features

1. **Supersession is recorded as an amendment.** A command writes, into the
   superseded Spec, which Spec delivered its content, on what date, and the
   reason given. The amendment is added; the PRD, the TechSpec and any evidence
   the Spec already holds are left exactly as they were.
2. **An amended Spec can be archived without a Task Graph.** Archive accepts a
   recorded supersession in place of completed Tasks and a passing gate, because
   those prove a Spec finished its own work and this Spec's work was finished
   elsewhere. Every other archive precondition still applies.
3. **The superseding Spec must exist.** A supersession pointing at a slug that
   is neither active nor archived is refused, as is a Spec superseding itself.
   A Spec that already carries a supersession is refused rather than amended
   twice.
4. **A Spec with its own Task Graph keeps proving itself.** Supersession is for
   a Spec that was never decomposed. Where a Task Graph exists, archive's
   existing evidence rules are unchanged.

## Non-Goals / Out of Scope

- Premise falsification, which Spec 0129 Core Feature 2 also covers. This Spec
  records that content shipped elsewhere, not that a premise turned out false.
- Revalidating consumers after their prerequisites land, which stays in Spec
  0129.
- Deleting, rewriting or merging the superseded Spec's content into the
  superseding one. The folder moves to history intact.
- Any change to archive's evidence rules for a Spec that has a Task Graph.

## Success Metrics

1. A fixture Spec with no Task Graph, superseded by an existing Spec, archives
   successfully; the same fixture without the supersession still refuses with
   the missing-Task-Graph reason.
2. After the amendment, every file the fixture already had is byte-identical
   apart from the one that carries the amendment.
3. A supersession naming an unknown slug, naming itself, or applied twice is
   refused, and the Spec directory is byte-identical afterwards.

## Decisions

- **Amend, do not annotate.** A paragraph a human writes into a PRD is a note.
  An amendment a command writes and a command reads is a record. The difference
  is whether the lifecycle can see it.
- **Supersession replaces the evidence, not the standard.** Archive still
  demands proof that a Spec is finished. For a superseded Spec the proof is that
  another Spec delivered its content, named and dated, rather than a gate it
  could never run.
- **Refuse a second supersession.** A Spec whose content moved once has a single
  true answer to where it went. Amending twice would leave the record ambiguous
  at exactly the moment someone is reading it to find out.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases matter most: a supersession that accepts an
unknown slug, or an archive that stops checking anything once an amendment is
present, would satisfy the happy path while making the record worthless.

## Research basis

The blocker was measured on this repository, not inferred: `roundfix archive`
on Spec 0128 refuses on the missing Task Graph, and 0128's own PRD carries a
hand-written Disposition paragraph explaining that Spec 0147 delivered its
content and that no supported operation retires it. That paragraph is the
evidence this Spec exists to replace with a record.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
