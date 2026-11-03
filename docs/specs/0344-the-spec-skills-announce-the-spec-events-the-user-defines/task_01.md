---
task: task_01
spec: 0344-the-spec-skills-announce-the-spec-events-the-user-defines
status: completed
type: docs
complexity: medium
---

# Task 01: The glossary defines the Spec Events and the PRD writers announce new and deadline events

## Overview

The glossary is the base of the process, so the two new terms land first.
Then the two skills that write a Spec's PRD gain the Spec Event Announcement
step: `write-prd` announces `new` when it writes the Spec's first PRD and
`deadline proposed` when it records a real deadline for a Spec already
announced, and `write-techspec` does the same on its correction route. This
answers the adopted Backlog Entry of 2026-11-03
([references](references/2026-11-03-the-spec-skills-announce-each-spec-event.md)),
reported by jarvis, whose events "nova" and "prazo proposto" these are.

## Requirements

1. MUST add to `CONTEXT.md`, through `domain-modeling`, right after the
   **Ignored Spec File** entry, the **Spec Event** and
   **Spec Event Announcement** entries exactly as `_techspec.md` → API Contract 7 gives
   them, each as a `**Term**:` line, its definition line and its `_Avoid_`
   line, and change no other glossary entry.
2. MUST add to `.agents/skills/write-prd/SKILL.md`, right before
   `## Advisory judgment`, a section headed `## Spec Event Announcement` that
   holds the paragraph of API Contract 1 and then the paragraph of API
   Contract 2, word for word (line wrapping is free).
3. MUST add to `.agents/skills/write-techspec/SKILL.md`, right before
   `## Advisory judgment`, the same heading with the paragraph of API
   Contract 1 and then the paragraph of API Contract 3, word for word.
4. MUST NOT name any personal assistant, tracker, channel or message block,
   and MUST NOT change the `### QA settlement` section or any other section of
   the two skills.
5. MUST copy the canonical skills with `make skills-sync`, then raise both
   version fields of `write-prd` and `write-techspec` and their mirrors with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`.
   Never hand-edit a mirror or a version field.

## Subtasks

- [ ] Write the two glossary entries through `domain-modeling`.
- [ ] Add the `write-prd` and `write-techspec` sections.
- [ ] Sync the mirrors and record the two skill versions.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines **Spec Event** and **Spec Event Announcement** with line-anchored definitions.
- [ ] Both skills, canonical and mirror, carry the generic paragraph and their `new` and `deadline proposed` events.
- [ ] The mirrors equal their canonical copies, and the skill version and documentation contracts pass.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0354-the-spec-skills-announce-spec-events-only-through-the-users-own-instructions.md`
- interface: `CONTEXT.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Consumers

- The glossary, the `write-prd` and `write-techspec` skill texts and the owned-skill version record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `grep -F -- '**Spec Event**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Spec Event\*\*(:| —( |$))' || { echo 'CONTEXT.md lacks Spec Event' >&2; exit 1; }; grep -F -- '**Spec Event Announcement**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Spec Event Announcement\*\*(:| —( |$))' || { echo 'CONTEXT.md lacks Spec Event Announcement' >&2; exit 1; }; for f in .agents/skills/write-prd/SKILL.md skills/write-prd/SKILL.md .agents/skills/write-techspec/SKILL.md skills/write-techspec/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "## Spec Event Announcement" || { printf 'missing phrase in %s: Spec Event Announcement\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "this skill never names one" || { printf 'missing phrase in %s: this skill never names one\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "An Agent inside a Daemon-assigned Run never announces." || { printf 'missing phrase in %s: An Agent inside a Daemon-assigned Run ne\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "When this skill records a proposed deadline in the PRD of a Spec already announced, announce" || { printf 'missing phrase in %s: When this skill records a proposed deadl\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/write-prd/SKILL.md skills/write-prd/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "once, when this skill writes the Spec's first" || { printf 'missing phrase in %s: once, when this skill writes the Specs f\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/write-techspec/SKILL.md skills/write-techspec/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "the feature route announces nothing, because" || { printf 'missing phrase in %s: the feature route announces nothing, bec\n' "$f" >&2; exit 1; }; done; cmp .agents/skills/write-prd/SKILL.md skills/write-prd/SKILL.md && cmp .agents/skills/write-techspec/SKILL.md skills/write-techspec/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task `CONTEXT.md` has no **Spec Event** line, so the first check fails; afterwards both terms are defined, both skills and their mirrors carry the step, and the skill version and documentation contracts pass.

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1, 2 and 4; Success Metric 3; Success Metric 4; Glossary; Decisions
- `_techspec.md` → API Contracts 1-3; API Contract 7; Invariants 1-5; Build Order 1
- ADR-0354

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, evidence-gate

Implemented the glossary entries immediately after **Ignored Spec File**, added the shared Spec Event Announcement paragraph and the route-specific `write-prd`/`write-techspec` paragraphs immediately before `## Advisory judgment`, synchronized both mirrors, and regenerated the owned-skill version record. The `### QA settlement` sections and unrelated skill sections were not changed.

Focused checks:

- `make skills-sync`: passed.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`: passed (`1 passed in 1 packages`); generated `write-prd` `0.0.18` and `write-techspec` `0.0.17` records.
- `git diff --check`: passed.
- Canonical/mirror equality checks: passed for `write-prd` and `write-techspec`.
- Required heading and glossary-anchor searches: found all expected entries.

Acceptance evidence:

- [x] `CONTEXT.md` has line-anchored **Spec Event** and **Spec Event Announcement** definitions after **Ignored Spec File**.
- [x] Both canonical skills and mirrors contain the shared paragraph plus their `new` and `deadline proposed` event paragraphs.
- [x] Mirrors equal canonical files, and the sanctioned owned-skill version recorder passed.

The task's aggregate `## Verification` command was not run; it remains for Daemon Verification. No commit, push, PR, task-status change, or Task Graph edit was made.
