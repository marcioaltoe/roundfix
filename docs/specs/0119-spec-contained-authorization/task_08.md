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
`.agents/skills/write-tasks/SKILL.md`, and their five shipped counterparts
`skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`,
`skills/write-techspec/SKILL.md`,
`skills/write-techspec/references/techspec-template.md` and
`skills/write-tasks/SKILL.md`, plus the derived pins that
`make baseline-digests` rewrites as sanctioned regeneration, and this Task
file. Stop before any other mutation. The bounded
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
- [ ] No path under `skills/` or `.agents/skills/` outside the ten named files
      is added, modified, renamed or left untracked.
- [ ] The derived pins match a fresh regeneration.

## Context

- instruction: `.agents/skills/write-tasks/SKILL.md`
- instruction: `docs/agents/specific-repository.md`

## Verification

- `grep -q '_authorization.md' .agents/skills/write-prd/references/prd-template.md && grep -q '_authorization.md' .agents/skills/write-techspec/references/techspec-template.md` — both templates require the Spec-contained record, which they do not today.
- `grep -q '_authorization.md' .agents/skills/write-tasks/SKILL.md` — the decomposition preflight names the record it must find.
- `grep -q '_authorization.md' skills/write-prd/references/prd-template.md || exit 1; make skills-sync-check` — the shipped bundle carries the change and matches its canonical source, so it was regenerated rather than hand-edited.
- `grep -q '_authorization.md' .agents/skills/write-techspec/references/techspec-template.md || exit 1; allowed="$(mktemp)"; changed="$(mktemp)"; paths="$(mktemp)"; unexpected="$(mktemp)"; printf '%s\n' .agents/skills/write-prd/SKILL.md .agents/skills/write-prd/references/prd-template.md .agents/skills/write-techspec/SKILL.md .agents/skills/write-techspec/references/techspec-template.md .agents/skills/write-tasks/SKILL.md skills/write-prd/SKILL.md skills/write-prd/references/prd-template.md skills/write-techspec/SKILL.md skills/write-techspec/references/techspec-template.md skills/write-tasks/SKILL.md | sort > "$allowed" || exit 1; git -c core.quotepath=false status --porcelain --untracked-files=all -- skills .agents/skills > "$changed" || exit 1; cut -c4- "$changed" | sort > "$paths" || exit 1; grep -F -x -v -f "$allowed" "$paths" > "$unexpected"; test ! -s "$unexpected" || { cat "$unexpected"; exit 1; }` — exactly the ten named files may differ. The audit reads staged, unstaged and untracked paths, and an empty out-of-scope list is the passing case rather than a failure.
- `grep -q '_authorization.md' skills/write-tasks/SKILL.md || exit 1; raw="$(mktemp)"; before="$(mktemp)"; after="$(mktemp)"; make baseline-digests || exit 1; find internal/baseline -type f -exec shasum {} + > "$raw" || exit 1; sort "$raw" > "$before" || exit 1; make baseline-digests || exit 1; find internal/baseline -type f -exec shasum {} + > "$raw" || exit 1; sort "$raw" > "$after" || exit 1; diff "$before" "$after"` — the shipped skill edit landed and the sanctioned derived pins reproduce byte for byte on a second regeneration.

## References

- `_prd.md` → User Stories 1; Core Features 1, 3; Goals 1, 2.
- `_techspec.md` → System Architecture: Authoring guidance; Build Order 5.
- ADR-0149.
