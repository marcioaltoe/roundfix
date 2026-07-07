---
task: task_04
spec: 0018-external-spec-root
status: pending
type: docs
complexity: low
---

# Task 04: Docs and skill sync for the External Spec Root surface

## Overview

Close the SKILL-matches-CLI gate for External Spec Root: document `specs.root`
and the external-artifact commit rule everywhere the config and command
surfaces are described, and re-sync the embedded skill bundle. Verifiable by
the drift check and by reading the docs against shipped behavior.

## Requirements

1. MUST document `specs.root` in the README Config section (default,
   precedence, relative/absolute resolution, validation failure) and the
   external-artifact commit rule in Command Boundaries, truthfully against
   implemented behavior.
2. MUST update the operational usage guide where it assumes Spec artifacts
   live under `docs/specs` inside the repository.
3. MUST update the canonical roundfix skill (SKILL.md) with the Spec Root
   behavior — resolution, worktree stability, startup report, and the
   commit-boundary warnings — and re-sync the embedded bundle.
4. MUST note in the docs that interim git shims papering over the symlink
   failure must be removed once this ships.

## Subtasks

- [ ] README Config and Command Boundaries updates
- [ ] Usage guide updates
- [ ] roundfix SKILL.md update and `make skills-sync`
- [ ] Drift and skills checks pass

## Acceptance Criteria

- [ ] README documents `specs.root` with its default, precedence, resolution,
      validation, and the external commit rule.
- [ ] The usage guide no longer assumes an in-repository Spec Root.
- [ ] `make skills-sync-check` reports no drift and `roundfix skills check`
      passes.

## Verification

- `rtk make verify` — expected: fmt-check, tests, skills-sync-check, skills
  check, and build all pass.

## References

`_prd.md` → Goals; User Experience. `_techspec.md` → Build Order 4; Risks
(shim removal). ADR-0035. CLAUDE.md SKILL.md-matches-CLI HARD RULE.
