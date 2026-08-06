---
task: task_05
spec: 0073-skill-versions-decoupled-from-the-binary
status: pending
type: chore
complexity: low
---

# Task 05: Synchronise the Roundfix Skill

## Overview

This Spec changes what `roundfix doctor` reports and what blocks a command, so
the Skill must teach the new contract before the Spec can close. This is one of
the two authorized tooling Tasks.

## Requirements

1. MUST document that owned-skill readiness is a comparison against a declared
   minimum, not a content match, and that a version at or above the minimum
   satisfies.
2. MUST document the three distinct states — satisfies, below the minimum, and
   unversioned-or-unresolvable — and state that an unreachable source is never
   reported as a missing skill.
3. MUST document that a below-minimum skill blocks and names the skill, the
   minimum, the version found, and the upgrade path.
4. MUST state plainly that third-party skills are never held to a version
   Roundfix invented for them.
5. MUST state that editing an owned skill no longer requires a regeneration
   step.
6. MUST regenerate the `skills/roundfix/**` mirror with `make skills-sync`, run
   `make baseline-digests`, and re-record the two characterization corpora that
   command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

7. MUST change only `.agents/skills/roundfix/**`, `skills/roundfix/**`, this
   Task file, and the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`.
8. MUST NOT change behaviour. This Task documents what shipped.

## Subtasks

- [ ] Document the comparison, the three states, and the third-party boundary.
- [ ] Run `make skills-sync`, then the regeneration chain.

## Acceptance Criteria

- [ ] The Skill states readiness is a comparison against a declared minimum.
- [ ] The Skill names all three states and their distinctions.
- [ ] The Skill states the four facts a below-minimum failure names.
- [ ] The Skill states third-party skills are never held to a version.
- [ ] The Skill states an owned-skill edit needs no regeneration step.
- [ ] `skills/roundfix/` is byte-identical to `.agents/skills/roundfix/`.
- [ ] `make verify` exits 0 after the regeneration chain.
- [ ] No Go source file changed.

## Context

- instruction: `.agents/skills/roundfix/SKILL.md`

## Verification

- `make skills-sync-check` — expected: exit 0; the mirror matches.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0.
- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.
- `git diff --name-only HEAD | grep -vE "^(\.agents/skills/roundfix/|skills/roundfix/|docs/specs/0073-skill-versions-decoupled-from-the-binary/task_05\.md$|internal/baseline/(assets/(setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Features 4, 5 and 6.
- `_techspec.md` → Integration Points; Build Order 5.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0081.
