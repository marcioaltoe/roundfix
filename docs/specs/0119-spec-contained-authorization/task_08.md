---
task: task_08
spec: 0119-spec-contained-authorization
status: pending
type: docs
complexity: medium
---

# Task 08: Teach the authoring skills the record's placement

## Overview

Canonical guidance alone does not change what gets authored; the skills that
write PRDs, TechSpecs and Task Graphs do. Teach them to place the authorization
record inside the Spec and to distinguish an approved grant from a proposal, so
a newly authored Spec carries its authority in the right shape. The slice is
verifiable on its own: the skill instructions and templates require the record,
and the shipped bundle matches its canonical source.

This is an authorized tooling Task. It may change only
`.agents/skills/write-prd/SKILL.md`,
`.agents/skills/write-prd/references/prd-template.md`,
`.agents/skills/write-techspec/SKILL.md`,
`.agents/skills/write-techspec/references/techspec-template.md`,
`.agents/skills/write-tasks/SKILL.md`, their five shipped counterparts under
`skills/`, the derived pins that `make baseline-digests` rewrites as sanctioned
regeneration, and this Task file. Stop before any other mutation. The bounded
set and the sanctioned regeneration come from the approved grant in
[_authorization.md](_authorization.md).

## Requirements

1. MUST require the PRD and TechSpec templates to cite the Spec-contained
   authorization record from the Tooling authority row, with the exact bounded
   repository paths.
2. MUST state in the authoring instructions that a record with a proposed state
   or an absent grant date authorizes nothing, and that decomposition is refused
   until an operative grant exists.
3. MUST state that a new or widened grant lands in target ancestry before the
   consuming squash delivery, so an author does not fold approval into the
   change it authorizes.
4. MUST NOT modify an upstream-managed skill; only the five canonical
   Roundfix-owned authoring skill files named above and their shipped
   counterparts change.
5. MUST regenerate the shipped bundle from its canonical source rather than
   editing the shipped copies by hand, and MUST regenerate the sanctioned
   derived pins afterwards.
6. MUST leave every other shipped skill byte-identical.

## Subtasks

- [ ] Require the Spec-contained record citation in both templates.
- [ ] State the non-granting states and the refusal they force in the
      authoring instructions.
- [ ] State the ancestry-before-squash rule where an author will read it.
- [ ] Regenerate the shipped bundle and the sanctioned derived pins.
- [ ] Prove every other shipped skill is unchanged.

## Acceptance Criteria

- [ ] Both templates instruct the author to cite the Spec-contained
      authorization record with exact bounded paths in the Tooling authority
      row.
- [ ] The authoring instructions state that a proposed or undated record grants
      nothing and that decomposition refuses without an operative grant.
- [ ] The ancestry-before-squash rule appears in the authoring instructions.
- [ ] Each of the five shipped counterparts is byte-identical to its canonical
      source, proven by the repository's drift check rather than by inspection.
- [ ] No skill outside the five named canonical files and their shipped
      counterparts changed.
- [ ] The derived pins match a fresh regeneration.

## Context

- instruction: `.agents/skills/write-tasks/SKILL.md`
- instruction: `docs/agents/specific-repository.md`

## Verification

- `grep -q '_authorization.md' .agents/skills/write-prd/references/prd-template.md && grep -q '_authorization.md' .agents/skills/write-techspec/references/techspec-template.md` — both templates require the Spec-contained record, which they do not today.
- `grep -q '_authorization.md' .agents/skills/write-tasks/SKILL.md` — the decomposition preflight names the record it must find.
- `grep -q '_authorization.md' skills/write-prd/references/prd-template.md || exit 1; make skills-sync-check` — the shipped bundle carries the change and matches its canonical source, so it was regenerated rather than hand-edited.
- `changed="$(git diff --name-only -- skills .agents/skills | grep -v -e '^\.agents/skills/write-\(prd\|techspec\|tasks\)/' -e '^skills/write-\(prd\|techspec\|tasks\)/')" || exit 1; test -z "$changed" || { printf '%s\n' "$changed"; exit 1; }` — no skill outside the bounded set changed.
- `grep -q '_authorization.md' skills/write-tasks/SKILL.md || exit 1; make baseline-digests || exit 1; git diff --quiet -- internal/baseline || { git --no-pager diff -- internal/baseline; exit 1; }` — the shipped skill edit landed and the sanctioned derived pins match a fresh regeneration.

## References

- `_prd.md` → User Stories 1; Core Features 1, 3; Goals 1, 2.
- `_techspec.md` → System Architecture: Authoring guidance; Build Order 5.
- ADR-0149.
