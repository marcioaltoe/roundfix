---
task: task_04
spec: 0080-cheap-detectors-run-before-the-gate
status: pending
type: chore
complexity: medium
---

# Task 04: Let a report row declare its evidence inputs

## Overview

Tooling Task one of two, under the authorization recorded at
`docs/workflow/authorizations/2026-08-06-proof-cost.md`. A row cannot be
carried forward unless it can say what it depends on, so the declaration lands
before the resolver that reads it — and lands inert, consumed by nothing yet.

The declaration is typed on purpose. ADR-0097 refuses to carry a row whose
evidence lives outside the repository, and the only way to refuse that
mechanically is for the row to name the kind of each input rather than leaving
it to be inferred.

## Requirements

1. MUST extend the QA Report row format with an `inputs:` declaration carrying
   typed entries, using the four kinds ADR-0097 names: `repository_path`,
   `external_repository`, `live_service`, and `elapsed_time`.
2. MUST instruct, in the qa-gate skill, how a row declares its inputs and when
   each kind applies — a row whose truth depends on elapsed time or an
   external repository declares that kind and is thereby never carriable.
3. MUST state that a row with no `inputs:` declaration is never carriable, so
   the declaration is opt-in and the default is to re-observe.
4. MUST keep the change additive: a report without `inputs:` behaves exactly
   as today, and no existing row, count, verdict rule, or naming contract
   changes.
5. MUST regenerate the `skills/qa-gate/**` mirror with `make skills-sync`.
6. MUST NOT implement carry-forward, read the declaration anywhere, or change
   what the gate does with it.
7. MUST change only the authorization's bounded files, their sanctioned
   deterministic digest fallout, and this task file.

## Subtasks

- [ ] Author the typed `inputs:` row format in the skill.
- [ ] Author when each input kind applies, and the no-declaration default.
- [ ] Sync the mirror and audit the changed-path scope.

## Acceptance Criteria

- [ ] The skill documents the `inputs:` declaration with all four kinds.
- [ ] The skill states that a row without `inputs:` is never carriable.
- [ ] The canonical skill and its mirror are byte-identical.
- [ ] Nothing consumes the declaration yet.
- [ ] The diff stays inside the bounded files, sanctioned fallout, and this
      task file.

## Context

- instruction: docs/workflow/authorizations/2026-08-06-proof-cost.md
- interface: .agents/skills/qa-gate/SKILL.md

## Verification

- `grep -q 'inputs:' .agents/skills/qa-gate/SKILL.md && grep -q 'repository_path' .agents/skills/qa-gate/SKILL.md && grep -q 'external_repository' .agents/skills/qa-gate/SKILL.md && grep -q 'live_service' .agents/skills/qa-gate/SKILL.md && grep -q 'elapsed_time' .agents/skills/qa-gate/SKILL.md`
  — expected: exit 0; the typed declaration and all four kinds are authored.
- `make skills-sync-check`
  — expected: exit 0; the embedded mirror matches the canonical skill.
  — expected: exit 0; nothing outside the bounded files and sanctioned fallout
  changed.
  — expected: exit 0; the derived artifacts are converged.

## References

- `_prd.md` → Core Feature 5; User Story 6.
- `_techspec.md` → Implementation Design (row input declaration); Build
  Order 4.
- ADR-0097, ADR-0081.
