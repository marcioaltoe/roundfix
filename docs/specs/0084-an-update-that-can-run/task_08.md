---
task: task_08
spec: 0084-an-update-that-can-run
status: pending
type: docs
complexity: low
---

# Task 08: State the two obligations where a Task author reads them

## Overview

A clause in the catalog governs every repository that adopts the Baseline; it
does not reach the two skills a Supervisor actually reads while decomposing a
Spec and running its gate. This slice states the outside-evidence obligation and
the glossary check inside `write-tasks` and `qa-gate`, where the Rehearsal Cases
rule and the QA row contract already live, so the obligation is visible at the
moment it applies.

## Requirements

1. MUST state, in the task-authoring skill, that at least one acceptance row per
   Spec rests on evidence originating outside the Spec's own artifacts and records
   that evidence's origin, placed with the existing Rehearsal Cases contract.
2. MUST state, in the QA gate skill, that the report records the origin of that
   row's external evidence, and that a row whose external evidence cannot be
   obtained is recorded as blocked with its reason rather than dropped.
3. MUST state, in the task-authoring skill, the glossary check at the close of a
   Spec, without duplicating the catalog clause's wording so the two cannot drift
   into contradiction.
4. MUST NOT add a new required section, frontmatter field, or refusal condition to
   either skill beyond what the clauses oblige.
5. MUST edit the authoritative skill sources rather than the generated copies, and
   sync the generated copies with the sanctioned command.
6. MUST change only the paths the authorization bounds.

## Subtasks

- [ ] State the outside-evidence obligation in the task-authoring skill.
- [ ] State the external-evidence recording and blocked-row handling in the QA
      gate skill.
- [ ] State the glossary check in the task-authoring skill.
- [ ] Sync the generated skill copies with the sanctioned command.

## Acceptance Criteria

- [ ] The task-authoring skill states the outside-evidence obligation adjacent to
      its Rehearsal Cases contract.
- [ ] The QA gate skill states how the external-evidence row's origin is recorded
      and how an unobtainable one is blocked.
- [ ] The task-authoring skill states the glossary check without restating the
      catalog clause verbatim.
- [ ] Neither skill gains a new required section or frontmatter field.
- [ ] The generated skill copies match their authoritative sources.

## Context

- instruction: `docs/workflow/authorizations/2026-08-08-evidence-from-outside-the-spec.md`
- instruction: `docs/agents/skill-dispatch.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`

## Verification

- `grep -q 'outside' .agents/skills/write-tasks/SKILL.md` — expected: exits 0, proving the obligation is stated in the authoritative source.
- `grep -q 'outside' .agents/skills/qa-gate/SKILL.md` — expected: exits 0, proving the QA gate states the recording rule.
- `diff -r .agents/skills/write-tasks skills/write-tasks > /tmp/0084-task-08-a.log 2>&1` — expected: exits 0, proving the generated copy matches its authoritative source.
- `diff -r .agents/skills/qa-gate skills/qa-gate > /tmp/0084-task-08-b.log 2>&1` — expected: exits 0, proving the generated copy matches its authoritative source.
- `go run ./cmd/roundfix skills check > /tmp/0084-task-08-c.log 2>&1` — expected: exits 0.

## References

- `_techspec.md` → Build Order 8.
- `_prd.md` → Core Features 8 and 9; User Story 7; Goal 5.
- ADR-0104.
