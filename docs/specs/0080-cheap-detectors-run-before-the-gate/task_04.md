---
task: task_04
spec: 0080-cheap-detectors-run-before-the-gate
status: completed
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
  — expected: exit 0; the embedded mirror matches the canonical skill.
  — expected: exit 0; nothing outside the bounded files and sanctioned fallout
  changed.
  — expected: exit 0; the derived artifacts are converged.

`make skills-sync-check` is deliberately absent: it is a repository gate that
passes whenever the mirror is already in sync, which is the state every Task is
dispatched from. The Run-level gate proves the mirror after this Task settles.

## References

- `_prd.md` → Core Feature 5; User Story 6.
- `_techspec.md` → Implementation Design (row input declaration); Build
  Order 4.
- ADR-0097, ADR-0081.

## Result

Implemented the additive QA Report row declaration in the canonical qa-gate
skill and regenerated its shipped mirror. The declaration uses `kind` and
`ref` entries, explains all four ADR-0097 kinds, makes absent or empty inputs
fail closed to re-observation, and states explicitly that no gate behavior
consumes the metadata yet.

Focused checks run during implementation:

- `rtk make skills-sync` — exited 0 and regenerated the owned skill mirrors.
- `rtk shasum -a 256 .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md`
  — both files produced
  `7a3ee44d32fe9644ffd5b8acc1963abbe96e07fdbe702d69525c7c4a4892ca08`.
- `rtk sh -c 'GOCACHE=/private/tmp/roundfix-task04-gocache rtk make skills-check'`
  — exited 0; Roundfix reported all shipped skill artifacts passed.
- `rtk git -c core.fsmonitor=false diff --check` — exited 0.
- `rtk git -c core.fsmonitor=false diff --name-only` — before this Result
  update, listed only `.agents/skills/qa-gate/SKILL.md`, its
  `skills/qa-gate/SKILL.md` mirror, and this Task file.

Acceptance evidence:

- The canonical skill's `Row input declaration` documents `inputs:` entries
  with `kind` and `ref`, including `repository_path`,
  `external_repository`, `live_service`, and `elapsed_time`, and explains when
  each applies.
- The same section states that a missing or empty declaration is never
  carriable and must be re-observed; any non-repository kind, including a kind
  in a mixed list, also refuses carry-forward.
- The canonical and mirror skill files have the identical SHA-256 digest
  recorded above after `make skills-sync`.
- The diff changes only skill instructions and this Result; no Go or other
  runtime path reads `inputs:`. The skill states that reports without the
  declaration retain all existing row, count, status, verdict, and naming
  behavior.
- The changed-path audit remains within the authorization's two exact skill
  paths and this assigned Task file; regeneration produced no digest fallout.

The Task's declared `## Verification` command was not run; the Daemon owns
that verification and settlement.
