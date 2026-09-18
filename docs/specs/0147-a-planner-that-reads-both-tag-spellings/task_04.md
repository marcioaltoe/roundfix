---
task: task_04
spec: 0147-a-planner-that-reads-both-tag-spellings
status: pending
type: docs
complexity: low
---

# Task 04: Tell the skill which spellings the planner reads

## Overview

The shipped Roundfix skill describes the Release Plan Command. It must say which
tag spellings the command accepts and what it does when one version carries
both. This Task is an authorized tooling mutation and may change only the two
bounded files plus its own Task file.

## Requirements

1. MUST state that the command accepts a stable tag with or without the `v`
   prefix.
2. MUST state that a highest version reachable under both spellings refuses in
   preflight, and name the selector that resolves it.
3. MUST state that the proposed version keeps the spelling of the tag it was
   selected from.
4. MUST regenerate the distributed mirror so both copies stay identical.
5. MUST NOT change any path outside the bounded list, the skill's version, or
   any other behavior the skill documents.

## Subtasks

- [ ] Write the accepted spellings into the skill's Release Plan section.
- [ ] Write the ambiguity refusal and its selector.
- [ ] Regenerate the mirror and confirm the copies match.

## Acceptance Criteria

- [ ] The skill names both accepted spellings.
- [ ] The skill states the refusal and its selector.
- [ ] The mirror is byte-identical to the canonical copy.
- [ ] No path outside the bounded list changes.

## Context

- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `grep -q "MAJOR.MINOR.PATCH" .agents/skills/roundfix/SKILL.md && grep -q "ambiguous" .agents/skills/roundfix/SKILL.md` — expected: exit 0; the canonical skill names the spellings and the refusal. Before this Task it does not.
- `grep -q "MAJOR.MINOR.PATCH" skills/roundfix/SKILL.md && diff -q .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — expected: exit 0; the mirror names the accepted spellings and stays identical to its canonical copy. Before this Task it does not name them.

## References

`_prd.md` → Core Features 1 and 3; Project Constraints: Tooling authority;
`_techspec.md` → System Architecture: Shipped skill; Build Order 4;
`_authorization.md`.
