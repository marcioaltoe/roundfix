---
spec: 0152-one-declared-acceptance-policy
status: active
created: 2026-09-21
surfaces: [backend, cli, docs]
---

# One declared-acceptance policy

Whether the newest QA Report lets a Spec proceed is decided twice, by two
implementations that disagree.

`internal/spec/archive.go` accepts a `pass`, and accepts a `partial` when its
blocked rows are declared unreachable — the rule its own comment states: "A
partial QA Report is eligible only when its blocked rows are declared
unreachable."

`DerivedQAVerification` renders an awk program into every terminal `qa` Task,
and that program ends `exit(closed && verdicts == 1 && verdict == "pass" ? 0 : 1)`.
A qualifying `partial` fails it.

The Daemon runs the rendered command to settle the `qa` Task. So a Spec whose QA
legitimately ends `partial` with declared-unreachable blocked rows cannot settle
its gate, and a Spec that cannot settle its gate never reaches the archive that
would have accepted it. The stricter copy wins by running first, and the
documented policy never gets consulted.

The generator's own comment says the rendered command exists "so readers can see
the contract" and that "changing that rendered command does not change the
effective contract". For archive that is true. For Task settlement the rendered
command *is* the effective contract, because it is what executes.

This Spec is carved from Spec 0122 Core Feature 4, which asks for one
declared-acceptance eligibility policy for the newest QA Report across
settlement, the derived QA command, and archive.

## Project Constraints

- Identifier strategy: applicable — QA Reports keep their dated, sequenced
  names and the existing newest-report selection. No identity is minted or
  renamed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local file reads only. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080 keeps an environment-blocked row
  distinct from a failure, ADR-0091 makes the QA gate a Task node of its own
  type, ADR-0097 carries a QA row forward only on declared, unmoved evidence,
  ADR-0096 makes the gate prove machine facts before spending an agent turn,
  ADR-0117 checks a defect at the stage that can produce it, ADR-0093 checks
  Spec consistency by citation, ADR-0104 accepts on evidence a Spec did not
  author, ADR-0155 makes the `qa` Task declare the matrix, ADR-0156 makes a
  declared promise name a consuming Task, and ADR-0130 keeps a path governed
  once an authorization has bounded it. All hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the derived Verification command and the
  shipped skill are public contract surfaces. Express maintainer authorization:
  the standing grant of 2026-09-18 for the skill files, plus
  `internal/spec/archive.go` granted 2026-09-21, both recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `internal/spec/archive.go`. Sanctioned regeneration: `make skills-sync`.
  `internal/spec/task.go`, `internal/spec/qa.go`, their tests and
  `internal/cli` are ordinary source that no authorization has bounded, each
  measured by file through the governance probe rather than inferred from its
  directory. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- One implementation decides whether a QA Report is acceptable.
- A Spec that legitimately ends `partial` can settle its gate and be archived.
- Nothing that is acceptable today stops being acceptable.

## Core Features

1. **One eligibility function.** A single exported decision takes the newest QA
   Report and answers whether it is acceptable, and why not when it is not.
   Archive calls it. The derived Verification reaches the same decision.
2. **The derived command stops re-deciding.** It delegates the verdict judgement
   instead of carrying a second copy of the rule in awk, so the two cannot drift
   apart again.
3. **A qualifying partial settles.** A newest report whose verdict is `partial`
   and whose blocked rows are declared unreachable settles the `qa` Task, as it
   already would have been archived.
4. **Nothing else loosens.** `fail`, an unparseable report, a missing report, a
   `partial` carrying finding- or environment-blocked rows, a `partial` with no
   declared blocked rows, and a `partial` whose declared count exceeds the
   Spec's unreachable declarations all still refuse, with the reason each gives
   today.

   A `pass` stays accepted exactly as it is today, **including one carrying
   environment-blocked rows**. ADR-0080 keeps an environment-blocked row
   distinct from a failure, and every QA Report in this repository carries one.
   Refusing it would archive nothing ever again.

## Non-Goals / Out of Scope

- Changing which report is newest. The existing dated-and-sequenced selection
  stands.
- Changing what a QA Agent may declare, or what makes a blocked row
  unreachable. This Spec makes one rule apply in both places; it does not
  rewrite the rule.
- The settlement-semantics documentation sweep of Spec 0122 Core Feature 5,
  beyond the shipped skill this Spec's own change requires.
- Any change to `roundfix settle`'s Task-level Verification replay.

## Success Metrics

1. A fixture Spec whose newest report is a qualifying `partial` settles its `qa`
   Task, where today the derived command refuses it.
2. The same fixture archives, as it does today.
3. Every refusal in Core Feature 4 is exercised and still refuses, with an
   unchanged reason, and a `pass` carrying an environment-blocked row is still
   accepted.

## Decisions

- **Delegate rather than synchronise.** Teaching the awk to evaluate blocked-row
  declarations would produce two correct copies, which is the state that
  produced this bug. One implementation cannot disagree with itself.
- **Move toward the documented rule, not away from it.** Archive's policy is the
  one the repository wrote down and the one a partial verdict was designed for.
  The derived command is the copy that drifted, so it is the copy that changes.
- **Loosen exactly one case.** A qualifying partial is the only verdict this
  Spec makes newly acceptable to settlement. Everything else refuses exactly as
  before, and Core Feature 4 is the control that proves it.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: a change that made settlement
accept any report would satisfy Core Feature 3 while destroying the gate.

## Research basis

Both implementations were read on this repository before authoring.
`internal/spec/archive.go:143-156` accepts `pass` and a qualifying `partial`;
`DerivedQAVerification` in `internal/spec/task.go:95` renders an awk program
ending `verdict == "pass" ? 0 : 1`. The generator's comment states that the
rendered command does not change the effective contract, which holds for archive
and not for Task settlement, where the rendered command is what the Daemon
executes. Every path this Spec touches was classified by file through the
governance probe.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
