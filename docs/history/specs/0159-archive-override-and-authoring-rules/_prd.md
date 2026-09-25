---
spec: 0159-archive-override-and-authoring-rules
status: archived
created: 2026-09-24
surfaces: [backend, cli, docs]
archived: "2026-09-25"
source_slug: 0159-archive-override-and-authoring-rules
---


# Archive override and authoring rules

Delivery 3b of the restructured queue, second half. Three gaps remain between
what the workflow promises and what it lets an author or maintainer do.

**The QA Archive Override exists only on paper.** The Baseline guidance and the
archive skill describe archiving a Spec despite failed or missing QA when the
maintainer authorizes it, stamped `qa_override: true`. `roundfix archive`
accepts no such input: a parser probe on 2026-09-08 answered
`unexpected argument "--qa-override"`. The documented recovery has no command.

**Two Specs can claim the same ordinal.** A Task names the ADR it will create
with `creates:`. Nothing checks that the number is free, and with deliveries now
authored in parallel two Specs can each plan the same number. The collision surfaces
only at merge, as two decisions under one number.

**Settlement and authoring guidance disagree or are missing.** The QA gate, the
archive skill and the Roundfix skill each describe settlement in their own
words. The authoring skill does not yet describe the two declarations Spec 0158
adds, temporal prerequisites, property-shaped acceptance, or whether a new test
class runs in the repository gate at all.

This Spec carries Spec 0122 Core Features 5 and 8 and Spec 0129 Core Features 4
and 5. The semantic-review half of Spec 0129 Core Feature 1 moves to the review
providers delivery, which owns reviewer selection.

## Project Constraints

- Identifier strategy: applicable — ADR ordinals become claimed identifiers: a
  Task's `creates:` path under `docs/adr/` reserves its number against the tree
  and every other active Spec. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080 accepts an environment-blocked
  row, ADR-0093 checks Spec consistency by citation, ADR-0104 accepts on evidence
  a Spec did not author, ADR-0130 keeps a path governed once bounded, ADR-0154
  keeps an authorized QA Archive Override archive-only, ADR-0155 makes the `qa`
  Task declare the matrix and ADR-0156 makes a declared promise name a consuming
  Task. This Spec's gate is bound by ADR-0091, ADR-0096, ADR-0097 and ADR-0117.
  All hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-24 for the skill content listed in
  [_authorization.md](_authorization.md), and the standing grant of 2026-09-21
  for the governed source this slice needs; bounded files:
  `internal/spec/archive.go`, `internal/spec/archive_test.go`,
  `internal/cli/cli_test.go`, `internal/speccheck/coherence.go`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-tasks/references/task-template.md`,
  `skills/write-tasks/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`,
  `skills/qa-gate/SKILL.md`, `.agents/skills/archive-spec/SKILL.md`,
  `skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`. Sanctioned regeneration: `make skills-sync`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A maintainer-authorized archive despite failed or missing QA has a command,
  and the command records exactly what was waived and on whose authority.
- An ADR ordinal cannot be claimed twice.
- Settlement reads the same in every skill that describes it, and the authoring
  skill describes every declaration a Task can make.

## Core Features

1. **QA Archive Override.** `roundfix archive <slug> --qa-override --approval
   <source> --reason <text>` archives a Spec whose non-QA Tasks are all
   completed although its QA is failed, missing or otherwise ineligible. It
   stamps `qa_override: true` with the approval source, the reason, the actual
   QA outcome and the revision archived. It never changes the QA Task's status
   or its report's verdict. It is refused when QA already qualifies, when a
   non-QA Task is not completed, and when approval or reason is missing.
2. **Claimed ordinals.** Spec check refuses, as `SC-ORDINAL-CLAIMED`, a Task
   whose `creates:` path under `docs/adr/` uses a number that another file on
   the tree already holds or that another active Spec's Task also creates.
3. **One settlement table.** The QA gate, archive and Roundfix skills carry the
   same settlement table: pass, qualifying declared partial,
   environment-blocked, failed, missing and override, each with what it settles
   and what it archives.
4. **Complete authoring guidance.** The task-writing skill describes
   `verification: independent`, `precondition_repairs`, claimed ordinals,
   temporal prerequisites as named prerequisites that never invent release
   authority, property-shaped acceptance, test seams, narrow Spec commits, and
   the rule that a newly required test class must run under the repository
   gate.

## Non-Goals / Out of Scope

- Semantic review of Spec decisions, which belongs to the review providers
  delivery.
- Any change to what an override means for publication: review, required
  checks, merge and release stay independently enforced.
- Ordinals of anything other than ADRs.

## Success Metrics

1. An authorized override archives a Spec with failed QA, its stamp names
   approval, reason, QA outcome and revision, and the QA Task and report are
   byte-identical before and after.
2. A Task creating an ADR number already held on the tree or by another active
   Spec fails Spec check with `SC-ORDINAL-CLAIMED`.
3. The settlement table is identical across the three skills and their
   distributed mirrors.

## Recorded limits

The corrective ceiling of two Tasks was spent on the defects the first pre-PR
review of 2026-09-24 found. The second review found five minor gaps, carried to
Spec 0169:

- An override of a Spec whose QA Task is failed or pending while its newest
  report says `pass` stamps `qa_override_qa_outcome: pass` and records nowhere
  that the QA Task was not completed.
- The Roundfix skill still says the override "refuses when QA already
  qualifies", the rule before Task 05.
- The archive-spec skill's Steps still describe stamping `qa_override: true` by
  hand instead of running the command.
- Two Tasks of one Spec creating different paths with the same ADR number are
  not reported until one of the files exists.
- A Spec without `_tasks.md` does not list `SC-ORDINAL-CLAIMED` among the
  detectors it skipped.

## Decisions

- **Archive-only, and only when needed.** An override on a Spec whose QA already
  qualifies is refused, so no normal archive ever carries a synthetic override.
- **Ordinals are claimed at authoring.** Checking at Spec check time finds the
  collision before any Run spends work on it.
- **One table, copied, checked.** The same text in three skills is enforced by a
  repository contract test rather than trusted to review.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: an override that flipped a QA
verdict, or an ordinal check that passed a duplicate, would pass the happy path
while breaking the evidence it exists to protect.

## Research basis

`Archive` in `internal/spec/archive.go` requires every Task completed and a
qualifying newest QA Report, with no override input. The override policy text
is in `internal/baseline/assets/modules/spec-workflow.json` and
`skills/archive-spec/SKILL.md`. The ordinal collision was met while authoring
Specs 0158 and 0159 in parallel on 2026-09-24, and Spec 0122's TechSpec records
the parser probe of 2026-09-08.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
