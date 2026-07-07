---
task: task_04
spec: 0018-external-spec-root
status: completed
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

- [x] README Config and Command Boundaries updates
- [x] Usage guide updates
- [x] roundfix SKILL.md update and `make skills-sync`
- [x] Drift and skills checks pass

## Acceptance Criteria

- [x] README documents `specs.root` with its default, precedence, resolution,
      validation, and the external commit rule.
- [x] The usage guide no longer assumes an in-repository Spec Root.
- [x] `make skills-sync-check` reports no drift and `roundfix skills check`
      passes.

## Verification

- `rtk make verify` — expected: fmt-check, tests, skills-sync-check, skills
  check, and build all pass.

## References

`_prd.md` → Goals; User Experience. `_techspec.md` → Build Order 4; Risks
(shim removal). ADR-0035. CLAUDE.md SKILL.md-matches-CLI HARD RULE.

## Result

Implemented the docs and skill-sync slice for External Spec Root.

- README Config now documents `specs.root` with default `docs/specs`, Project
  Config over User Config over built-in precedence, relative and absolute path
  resolution, validation failures that name the resolved path, and the external
  predicate. README Command Boundaries now documents non-default startup
  reporting, external/symlink-crossing artifact drops, warning shapes, external
  Task settle-without-commit behavior, external QA Report behavior, and the
  requirement to remove temporary git shims that hid symlink pathspec failures.
- `docs/usage.md` now describes Spec execution from the resolved Spec Root
  instead of assuming an in-repository `docs/specs` root. Its remaining
  `docs/specs` reference is explicitly the default layout.
- `.agents/skills/roundfix/SKILL.md` now documents Spec Root resolution,
  worktree stability, non-default startup reporting, Interactive Input listing
  from the resolved root, external commit-boundary warnings, and archive/review
  artifact locations. Ran `rtk make skills-sync`, which regenerated
  `skills/roundfix/SKILL.md` from the canonical skill.

Evidence:

- `rtk make skills-sync-check` exited 0 with no drift output.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` passed:
  `Roundfix skill check passed: roundfix, write-idea, write-prd, write-techspec,
  write-tasks, setup-workflow, implement-task, implement-spec, brainstorming,
  council, business-analyst, archive-spec, qa-gate, evidence-gate`.
- `rtk make verify` passed. It ran `rtk go test ./...` with
  `859 passed in 18 packages`, `roundfix skills check`, and
  `rtk go build -buildvcs=false -o bin/roundfix ./cmd/roundfix`; fmt-check and
  skills-sync-check produced no failure output.
