---
task: task_05
spec: 0065-loop-order-and-verification-honesty
status: pending
type: chore
complexity: low
---

# Task 05: Synchronise the authoring and CLI Skills

## Overview

This Spec adds four `SC-*` rules that refuse Task Graphs which author cleanly
today, and it settles the loop order. Two Skills must teach that before the
Spec can close: `write-tasks`, which authors the graphs the rules now refuse,
and `roundfix`, whose `spec check` contract gained the rules.

This is the authorized tooling Task.

## Requirements

1. MUST document, in `write-tasks`, that a Task's Verification must be able to
   fail when no work was done, with the refused shape named: only
   repository-wide gates plus a clean-tree check.
2. MUST document, in `write-tasks`, that mutually unsatisfiable requirements
   are refused, and the section a rehearsal Task uses to declare its cases and
   their observation.
3. MUST document, in `roundfix`, the four new `SC-*` identifiers as part of the
   `spec check` contract.
4. MUST state the settled loop order in both Skills identically to the sources
   task_01 corrected, so task_04's rule passes over them.
5. MUST regenerate both mirrors with `make skills-sync`.
6. MUST run `make baseline-digests`, then re-record the two characterization
   corpora that command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

7. MUST change only `.agents/skills/write-tasks/**`,
   `.agents/skills/roundfix/**`, their `skills/**` mirrors, this Task file, and
   the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`.
8. MUST NOT change behaviour. This Task documents what shipped.

## Subtasks

- [ ] Document the three authoring rules in `write-tasks`.
- [ ] Document the four `SC-*` identifiers in `roundfix`.
- [ ] Run `make skills-sync`, then `make baseline-digests` and both re-records.

## Acceptance Criteria

- [ ] `write-tasks` names the refused Verification shape.
- [ ] `write-tasks` names the contradiction rule and the rehearsal declaration.
- [ ] `roundfix` lists the four new `SC-*` identifiers.
- [ ] Both Skills state the settled loop order identically.
- [ ] Both mirrors are byte-identical to their canonical Skills.
- [ ] `make verify` exits 0 after the regeneration chain.
- [ ] No Go source file changed.

## Context

- instruction: `.agents/skills/write-tasks/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`

## Verification

- `make skills-sync-check` — expected: exit 0; both mirrors match.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0.
- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.
- `git diff --name-only HEAD | grep -vE "^(\.agents/skills/(write-tasks|roundfix)/|skills/(write-tasks|roundfix)/|docs/specs/0065-loop-order-and-verification-honesty/task_05\.md$|internal/baseline/(assets/(modules|setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Features 1, 3, 4
  and 5.
- `_techspec.md` → Integration Points; Build Order 4.
- `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0081.
