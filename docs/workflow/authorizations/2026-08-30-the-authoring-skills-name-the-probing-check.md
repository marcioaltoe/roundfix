---
granted: 2026-08-26
action: Make the four Roundfix-owned authoring skills name the probing form of the Spec Consistency Check, and record the auditing binary in the QA Report contract.
consuming: 0116-a-verdict-that-states-its-own-scope
paths:
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-techspec/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/qa-gate/SKILL.md
---

# Tooling authorization — the authoring skills name the probing check (2026-08-26)

On 2026-08-26 the maintainer granted authorization to edit the four
Roundfix-owned authoring skills for Spec 0116. The grant was recorded in that
Spec's PRD prose and nowhere else. This record moves it to the directory the
changed-path audit actually reads, without widening it.

## Why this record exists separately

`docs/specs/0116-a-verdict-that-states-its-own-scope/_prd.md` states the grant,
its date, and its bounded files. `SC-TOOLING-UNAUTHORIZED` and the QA gate's
`QA-AUTH-PATHS` resolve a grant from `docs/workflow/authorizations/`, not from
the citing artifact, so a Spec whose authorization exists only in its own prose
reaches its gate and is refused there. Spec 0118 was refused for the adjacent
reason on 2026-08-27: its record existed but named a path in prose that its
`paths:` frontmatter omitted.

## What this covers

`--run-verification` appears zero times across every skill in this repository,
so an author following the shipped guidance runs the form of the Spec
Consistency Check that does not execute the authored Verification commands.
Measured in this repository on 2026-08-25/26: five of eight authored commands in
Spec 0098 and six of six in Spec 0113 were vacuous or non-hermetic past a check
that reported `No findings.`, costing three Unresolved Runs.

The QA gate skill additionally owns the QA Report template, which records the
audited commit but not the binary that produced the verdict.

## Authorized paths

- `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/write-tasks/SKILL.md`, limited to the step that instructs the
  author to run the Spec Consistency Check: naming the probing form and stating
  what a clean non-probing verdict does not cover.
- `.agents/skills/qa-gate/SKILL.md`, limited to the same check instruction plus
  the QA Report template's record of the auditing binary.

The generated copies under `skills/` are rewritten by the declared
`make skills-sync` and are listed above because the changed-path audit reads
this frontmatter rather than this prose. They remain sanctioned fallout of the
authorized source edit under ADR-0081; a hand-edited generated copy is still an
unauthorized mutation.

## Bounded by purpose

This grant covers how the skills instruct an author to run the check, and the
QA Report's record of its own producer. It does not authorize changing what the
Verification prober classifies, the Daemon's pre-work probe, the QA gate's
verdict rules, its row contract, which checks it runs, or any other skill in
the repository.

## Consuming Spec

This authorization is consumed by Spec
`0116-a-verdict-that-states-its-own-scope`.

## Commit choreography

This record lands as its own commit, before any commit that edits a skill.
