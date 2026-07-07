---
task: task_04
spec: 0017-run-discovery
status: pending
type: docs
complexity: low
---

# Task 04: Docs and skill sync for the Run discovery surface

## Overview

Close the SKILL-matches-CLI gate for the new surface: document `runs list` and
the Attach picker everywhere the CLI surface is described, and re-sync the
embedded skill bundle. Verifiable by the skills drift check and by reading the
updated docs against the shipped behavior.

## Requirements

1. MUST document `runs list` (columns, `--all`, `--active`, empty result,
   exit codes) and the no-argument Attach behavior in the README Commands and
   Command Boundaries sections, truthfully against implemented behavior.
2. MUST update the operational usage guide's monitoring flow to discover Runs
   with `runs list` and the Attach picker instead of requiring a captured run
   id.
3. MUST update the canonical roundfix skill (SKILL.md) with the new command
   surface and re-sync the embedded bundle so the drift check passes.
4. MUST keep top-level `--help` and the `runs`/`attach` usage text consistent
   with the docs.

## Subtasks

- [ ] README Commands and Command Boundaries entries for `runs list` and the
      Attach picker
- [ ] Usage guide monitoring flow update
- [ ] roundfix SKILL.md update and `make skills-sync`
- [ ] Drift and skills checks pass

## Acceptance Criteria

- [ ] README documents the listing columns, both flags, the empty-result
      contract, and the no-argument Attach behavior.
- [ ] The usage guide shows Run discovery without a captured run id.
- [ ] `make skills-sync-check` reports no drift and `roundfix skills check`
      passes.

## Verification

- `rtk make verify` — expected: fmt-check, tests, skills-sync-check, skills
  check, and build all pass.

## References

`_prd.md` → Goals; Success Metrics. `_techspec.md` → Build Order 4. CLAUDE.md
SKILL.md-matches-CLI HARD RULE.
