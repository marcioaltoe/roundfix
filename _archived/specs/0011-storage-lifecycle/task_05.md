---
task: task_05
spec: 0011-storage-lifecycle
status: completed
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

- [x] Update user docs for artifact location, `logs.agent`, and Archive Command
- [x] Update SKILL.md + agent manifest for the changed surface
- [x] `make skills-sync` to regenerate the embedded copy
- [x] Verify no skill drift and no doc/behavior mismatch

## Acceptance Criteria

- [x] Docs describe the three-branch artifact hierarchy, the opt-in `logs.agent` default, and `roundfix archive` accurately.
- [x] SKILL.md matches the shipped CLI surface; the embedded copy is regenerated.
- [x] `roundfix skills check` passes and `skills-sync-check` reports no drift.

## Verification

- `rtk make skills-sync` then `rtk go run ./cmd/roundfix skills check` — expected: check passes, no drift.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all pass.

## References

`_prd.md` → all stories (documentation). `_techspec.md` → Build Order 5. ADR-0029,
ADR-0030. CLAUDE.md skill-ownership and SKILL.md-matches-CLI gate.

## Result

- Updated `README.md`, `CONTEXT.md`, and agent workflow docs so user-facing docs describe the ADR-0029 review-artifact hierarchy, the ADR-0030 default-off `logs.agent` opt-in, and the `roundfix archive <slug>` precondition and outcome.
- Updated `.agents/skills/roundfix/SKILL.md` and `.agents/skills/roundfix/agents/openai.yaml` for the Archive Command, the review `--spec` selector, `logs.agent`, review artifact storage, and the current missing-check note; regenerated `skills/roundfix/` with `rtk make skills-sync`.
- Verification evidence: `rtk make skills-sync` passed; `rtk go run ./cmd/roundfix skills check` passed with `Roundfix skill check passed: roundfix`; `rtk make skills-sync-check` passed with no drift output; `rtk make verify` passed with `Go test: 765 passed in 17 packages`, `Roundfix skill check passed: roundfix`, and a successful build.
