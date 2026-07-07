---
task: task_05
spec: 0020-run-browser
status: pending
type: docs
complexity: low
---

# Task 05: Docs and skill sync for the Run Browser surface

## Overview

Close the SKILL-matches-CLI gate: document the Run Browser, the enriched
`runs list` contract (`--state`, `--limit`, columns, notes), and the updated
attach wording everywhere the CLI surface is described, then re-sync the
embedded skill bundle. The CONTEXT.md Run Browser term already exists.

## Requirements

1. MUST update the README Commands and Command Boundaries for `runs`
   (browser), `runs list` (columns, flags, notes, Active default), and
   `attach` (browser entry, unknown-id wording), truthfully against
   implemented behavior.
2. MUST update the operational usage guide's monitoring flow: Run Browser
   for humans, bounded `runs list` for agents, including the Active-only
   default and the widening flags.
3. MUST update the canonical roundfix skill (SKILL.md) with the same
   contract — including that `runs list` defaults to Active Runs and agents
   widen with `--state`/`--limit` — and re-sync the embedded bundle.
4. MUST keep help text consistent with the docs.

## Subtasks

- [ ] README Commands and Command Boundaries updates
- [ ] Usage guide monitoring flow update
- [ ] roundfix SKILL.md update and `make skills-sync`
- [ ] Drift and skills checks pass

## Acceptance Criteria

- [ ] README documents the browser keys, both flags, the Active default, the
      hidden-count notes, and the attach wording.
- [ ] The usage guide routes humans to the Run Browser and agents to the
      bounded listing.
- [ ] `make skills-sync-check` reports no drift and `roundfix skills check`
      passes.

## Verification

- `rtk make verify` — expected: fmt-check, tests, skills-sync-check, skills
  check, and build all pass.

## References

`_prd.md` → Goals; User Experience. `_techspec.md` → API Contracts; Build
Order 5. CLAUDE.md SKILL.md-matches-CLI HARD RULE. CONTEXT.md → Run Browser.
