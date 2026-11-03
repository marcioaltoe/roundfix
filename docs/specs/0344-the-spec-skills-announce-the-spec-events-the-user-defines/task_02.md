---
task: task_02
spec: 0344-the-spec-skills-announce-the-spec-events-the-user-defines
status: pending
type: docs
complexity: low
---

# Task 02: The graph and archive skills announce authorized, completed and cancelled

## Overview

`write-tasks` gains the Spec Event Announcement step for `authorized`, when
it first writes the Task Graph of a Spec whose authorization record is
approved, and `archive-spec` gains it for `completed` and `cancelled`, once
the Archive Command or the Cancel Command has written the Archive Record,
with a short business summary and the record's path. This answers the
adopted Backlog Entry of 2026-11-03
([references](references/2026-11-03-the-spec-skills-announce-each-spec-event.md)),
reported by jarvis, whose events "autorizada", "concluída" and "cancelada"
these are.

## Requirements

1. MUST add to `.agents/skills/write-tasks/SKILL.md`, right after the
   `## Release decision` section and before `## Process`, a section headed
   `## Spec Event Announcement` that holds the paragraph of `_techspec.md` →
   API Contract 1 and then the paragraph of API Contract 4, word for word
   (line wrapping is free).
2. MUST add to `.agents/skills/archive-spec/SKILL.md`, right after the
   `## Archive Record` section and before `## Unarchive`, the same heading
   with the paragraph of API Contract 1 and then the paragraph of API
   Contract 5, word for word.
3. MUST NOT name any personal assistant, tracker, channel or message block,
   and MUST NOT change the `### QA settlement` section of `archive-spec` or
   any other section of the two skills.
4. MUST copy the canonical skills with `make skills-sync`, then raise both
   version fields of `write-tasks` and `archive-spec` and their mirrors with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`.
   Never hand-edit a mirror or a version field.

## Subtasks

- [ ] Add the `write-tasks` section.
- [ ] Add the `archive-spec` section.
- [ ] Sync the mirrors and record the two skill versions.

## Acceptance Criteria

- [ ] Both skills, canonical and mirror, carry the generic paragraph and their own events.
- [ ] The `### QA settlement` section of `archive-spec` is unchanged.
- [ ] The mirrors equal their canonical copies, and the skill version contracts pass.

## Context

- instruction: `docs/adr/0354-the-spec-skills-announce-spec-events-only-through-the-users-own-instructions.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Consumers

- The `write-tasks` and `archive-spec` skill texts and the owned-skill version record — consumers: `./skills`

## Verification

- `for f in .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md .agents/skills/archive-spec/SKILL.md skills/archive-spec/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "## Spec Event Announcement" || { printf 'missing phrase in %s: Spec Event Announcement\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "this skill never names one" || { printf 'missing phrase in %s: this skill never names one\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "An Agent inside a Daemon-assigned Run never announces." || { printf 'missing phrase in %s: An Agent inside a Daemon-assigned Run ne\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "once, when this skill first writes the Spec's Task Graph and the Spec's" || { printf 'missing phrase in %s: once, when this skill first writes the S\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "A run that amends an existing Task Graph announces nothing." || { printf 'missing phrase in %s: A run that amends an existing Task Graph\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/archive-spec/SKILL.md skills/archive-spec/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "the announcement adds a summary of one to three sentences written for a business reader, free of technical terms." || { printf 'missing phrase in %s: the announcement adds a summary of one t\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "When the Delivery Queue writes the record, no Agent announces; the Supervisor session announces" || { printf 'missing phrase in %s: When the Delivery Queue writes the recor\n' "$f" >&2; exit 1; }; done; cmp .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md && cmp .agents/skills/archive-spec/SKILL.md skills/archive-spec/SKILL.md && go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task neither skill has the heading, so the first check fails; afterwards both skills and their mirrors carry the step, and the skill version and mirror contracts in `./skills` pass.

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1 and 2; Success Metric 1; Success Metric 4
- `_techspec.md` → API Contract 1; API Contracts 4-5; Invariants 1-5; Build Order 2
- ADR-0354
