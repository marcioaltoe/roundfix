---
task: task_05
spec: 0011-storage-lifecycle
status: pending
type: docs
complexity: low
---

# Task 05: Docs and skill sync

## Overview

Bring the docs and the Roundfix Skill in line with the storage/lifecycle
changes: the new review-artifact default, the `logs.agent` opt-in, and the
Archive Command. Runs after the behavior tasks so the documented contract
matches shipped behavior, and closes the mandatory SKILL.md-matches-CLI gate.

## Requirements

1. MUST update user-facing docs to describe the review-artifact location
   hierarchy, the `logs.agent` opt-in default, and the `roundfix archive`
   command and its precondition.
2. MUST update `.agents/skills/roundfix/SKILL.md` (and its agent manifest) for
   any changed command surface — the Archive Command and the `--spec` selector —
   and regenerate the embedded copy with `make skills-sync`.
3. MUST keep CONTEXT.md glossary usage consistent (Archive Command already
   added) and reference ADR-0029 and ADR-0030 where the decisions are described.
4. MUST leave no drift: `roundfix skills check` and `skills-sync-check` pass.

## Subtasks

- [ ] Update user docs for artifact location, `logs.agent`, and Archive Command
- [ ] Update SKILL.md + agent manifest for the changed surface
- [ ] `make skills-sync` to regenerate the embedded copy
- [ ] Verify no skill drift and no doc/behavior mismatch

## Acceptance Criteria

- [ ] Docs describe the three-branch artifact hierarchy, the opt-in `logs.agent` default, and `roundfix archive` accurately.
- [ ] SKILL.md matches the shipped CLI surface; the embedded copy is regenerated.
- [ ] `roundfix skills check` passes and `skills-sync-check` reports no drift.

## Verification

- `rtk make skills-sync` then `rtk go run ./cmd/roundfix skills check` — expected: check passes, no drift.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all pass.

## References

`_prd.md` → all stories (documentation). `_techspec.md` → Build Order 5. ADR-0029,
ADR-0030. CLAUDE.md skill-ownership and SKILL.md-matches-CLI gate.
