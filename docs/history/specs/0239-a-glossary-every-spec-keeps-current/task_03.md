---
task: task_03
spec: 0239-a-glossary-every-spec-keeps-current
status: completed
type: docs
complexity: medium
---

# Task 03: The authoring, QA and Roundfix skills write and check the Glossary Declaration

## Overview

The autonomous authoring route (write-prd decided autonomously, then
write-techspec, then write-tasks) never reaches `domain-modeling`, and no skill
tells an author to declare the terms a Spec introduces or a docs Task to write
them ([the adopted Backlog Entry of 2026-10-06](references/2026-10-06-the-context-driven-loop-does-not-keep-context-md-current.md)).
This Task adds `_techspec.md` → Skill text to `write-prd`, `write-techspec`,
`write-tasks`, `qa-gate` and `roundfix`, creates the `write-prd` glossary guide
whose adding commit starts the glossary horizon, syncs the mirrors and records
the raised versions (ADR-0244, ADR-0189). It follows task_01 because the skills
describe the findings and the refusal task_01 ships. It is verifiable on its
own: the skills state the route and their mirrors and versions agree.

This is an authorized tooling Task. It may change only the files in its
Context, the mirrors `make skills-sync` rewrites, and this Task file.

## Requirements

1. MUST add to each skill the content `_techspec.md` → Skill text gives it,
   under the heading it names, writing each quoted sentence exactly and without
   code spans inside it. Existing text MUST stay as it is, and no text inside
   the `### QA settlement` section of any skill may change.
2. MUST create `.agents/skills/write-prd/references/glossary.md` with the
   content Skill text lists for it, consistent with `_techspec.md` →
   Invariants 1 to 9, and link it from the `write-prd` skill's new heading.
3. MUST add a `## Glossary` section before `## Decisions` in
   `.agents/skills/write-prd/references/prd-template.md` and an optional one in
   `.agents/skills/write-techspec/references/techspec-template.md`, each with a
   comment naming the entry forms `adds`, `changes` and `not a term`.
4. MUST describe in `.agents/skills/roundfix/references/spec.md` the Glossary
   Declaration, the three codes and the glossary horizon, and in
   `.agents/skills/roundfix/references/archive.md` the refusal sentence of
   Skill text.
5. MUST raise both version fields of each changed skill's `SKILL.md`
   (`write-prd`, `write-techspec`, `write-tasks`, `qa-gate`, `roundfix`) by one
   patch step from the value on this Task's starting commit, then run
   `make skills-sync`, record the versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   run `make baseline-digests`, and name in the Result every file those
   commands rewrote.
6. MUST NOT edit the upstream `domain-modeling`, `grilling` or
   `grill-with-docs` skills.

## Subtasks

- [ ] Add the write-prd heading, guide and template section.
- [ ] Add the write-techspec, write-tasks and qa-gate headings.
- [ ] Update the Roundfix Skill's spec and archive references.
- [ ] Raise, sync and record the versions.

## Acceptance Criteria

- [ ] The authoring skills name domain-modeling and the Glossary Declaration on
      the autonomous route, and keep grilling and grill-with-docs as the
      interactive entry.
- [ ] write-tasks gives each declared term a glossary requirement in the docs
      Task, and the QA gate records the declared terms.
- [ ] The Roundfix Skill describes the three codes and the archive refusal.
- [ ] Every mirror equals its canonical skill and every raised version is
      recorded.

## Context

- instruction: `docs/adr/0244-a-spec-declares-the-domain-terms-it-introduces-and-the-check-holds-them-to-the-glossary.md`
- instruction: `docs/adr/0187-the-roundfix-skill-and-the-command-reference-are-read-one-command-at-a-time.md`
- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- instruction: `docs/user-guide/commands/spec.md`
- instruction: `docs/user-guide/commands/archive.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- creates: `.agents/skills/write-prd/references/glossary.md`
- creates: `skills/write-prd/references/glossary.md`
- interface: `.agents/skills/write-prd/references/prd-template.md`
- interface: `skills/write-prd/references/prd-template.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-techspec/references/techspec-template.md`
- interface: `skills/write-techspec/references/techspec-template.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `for pair in ".agents/skills/write-prd/SKILL.md|activate domain-modeling before writing it, even on the autonomous route" ".agents/skills/write-prd/SKILL.md|grilling or grill-with-docs stays the interactive entry" ".agents/skills/write-prd/references/glossary.md|is written by one of the Spec's own Tasks, never left for a later Spec" ".agents/skills/write-prd/references/glossary.md|SC-GLOSSARY-MISSING" ".agents/skills/write-prd/references/prd-template.md|## Glossary" ".agents/skills/write-techspec/SKILL.md|declare it in the Glossary section of the PRD or of the TechSpec" ".agents/skills/write-techspec/references/techspec-template.md|## Glossary" ".agents/skills/write-tasks/SKILL.md|carries a requirement to write the term through domain-modeling" ".agents/skills/qa-gate/SKILL.md|whether the glossary defines it after the work" ".agents/skills/roundfix/references/spec.md|SC-GLOSSARY-UNPLANNED" ".agents/skills/roundfix/references/archive.md|Archive refuses a Spec with a Glossary Gap, also under a QA Archive Override"; do file="${pair%%|*}"; phrase="${pair#*|}"; test -f "$file" || { printf 'missing file: %s\n' "$file" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; for skill in write-prd write-techspec write-tasks qa-gate roundfix; do diff -r ".agents/skills/$skill" "skills/$skill" >/dev/null || { printf 'mirror differs: skills/%s\n' "$skill" >&2; exit 1; }; done; make skills-sync-check || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills` — expected: exit 0; before this Task none of the skills names the Glossary Declaration and the glossary guide does not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 2; User Story 2; Core Feature 6; Success Metric 6
- [_techspec.md](_techspec.md) — Skill text; Invariants; API Contract 7; Testing Approach 5; Build Order 3
- ADR-0244; ADR-0187; ADR-0189; ADR-0233

## Result

Implemented the Glossary Declaration guidance for the five owned skills,
created the write-prd glossary guide and both template sections, documented
the Roundfix finding codes, horizon and archive refusal, synchronized the
mirrors, and raised each canonical skill's two version fields by one patch
step. The Task status remains Daemon-owned.

Focused checks and generated evidence:

- `make skills-sync` completed; it rewrote the mirrored files under
  `skills/qa-gate/`, `skills/roundfix/`, `skills/write-prd/`,
  `skills/write-tasks/`, and `skills/write-techspec/`, including the new
  `skills/write-prd/references/glossary.md`.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` first hit the sandbox default-cache permission
  error; the same required check passed with
  `GOCACHE=/tmp/roundfix-0239-task03-gocache`, recording the five raised
  versions in `skills/testdata/owned-skill-versions.json`.
- `make baseline-digests` passed and reported no derived-artifact changes.
- The final diff contains only the authorized skill sources, mirrors, owned
  version record, new glossary guide, and this Task's Result; the upstream
  `domain-modeling`, `grilling`, and `grill-with-docs` skills are untouched.

Acceptance evidence:

- The write-prd and write-techspec headings direct authors to
  `domain-modeling` on the autonomous route and preserve `grilling` and
  `grill-with-docs` as the interactive entry; the templates expose `adds`,
  `changes`, and `not a term` forms.
- The write-tasks heading requires each declared term to be written through
  `domain-modeling` and checked in glossary Verification text; the QA heading
  records whether each declared term is defined after the work.
- The Roundfix spec reference describes the Glossary Declaration, horizon and
  `SC-GLOSSARY-UNDECLARED`, `SC-GLOSSARY-UNPLANNED`, and
  `SC-GLOSSARY-MISSING`; the archive reference carries the required refusal
  sentence.
- The version-recording test passed after synchronization, and the mirror
  rewrites plus recorded JSON are present for all five raised skills.
