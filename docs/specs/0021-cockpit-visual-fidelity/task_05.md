---
task: task_05
spec: 0021-cockpit-visual-fidelity
status: pending
type: docs
complexity: low
---

# Task 05: Docs and skill sync for the cockpit visual contract

## Overview

Close the SKILL-matches-CLI gate for the visual contract: update the Live Run
View descriptions where they name concrete cockpit behavior (group markers,
summary rows, indicator, empty states), and re-sync the embedded skill
bundle. Verifiable by the drift check and reading docs against the shipped
render.

## Requirements

1. MUST update the README Live Run View section and the roundfix SKILL.md
   where the cockpit's concrete rendering is described: automatic group
   collapse, summary-only event rows with full content in the Detail Modal,
   the `Live · detail` indicator, and explanatory empty states.
2. MUST note the color behavior: state colors in capable terminals, full
   meaning preserved under `ROUNDFIX_COLOR=never`/`NO_COLOR`.
3. MUST re-sync the embedded bundle so the drift check passes.

## Subtasks

- [ ] README Live Run View updates
- [ ] roundfix SKILL.md updates and `make skills-sync`
- [ ] Drift and skills checks pass

## Acceptance Criteria

- [ ] README and SKILL.md describe the shipped rendering truthfully —
      collapse rules, bounded summaries, indicator, empty states, color
      degrade.
- [ ] `make skills-sync-check` reports no drift and `roundfix skills check`
      passes.

## Verification

- `rtk make verify` — expected: fmt-check, tests, skills-sync-check, skills
  check, and build all pass.

## References

`_prd.md` → Goals; User Experience. `_techspec.md` → Build Order 5.
CLAUDE.md SKILL.md-matches-CLI HARD RULE.
