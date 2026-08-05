---
task: task_04
spec: 0077-a-green-check-is-not-a-review
status: pending
type: chore
complexity: low
---

# Task 04: Synchronise the Roundfix Skill

## Overview

This Spec changes when `watch --until-clean` will and will not merge, so the
Skill must teach it before the Spec can close. This is the authorized tooling
Task.

## Requirements

1. MUST document that a Review Source refusal resolves as skipped evidence and
   never verifies a head, and that `watch --until-clean` will not merge it.
2. MUST document that an unrecognised signal resolves pending, and state plainly
   that a green check is not evidence a review ran.
3. MUST document that no automatic retry exists yet, so a reader does not wait
   for one — the follow-on Spec owns it.
4. MUST regenerate the `skills/roundfix/**` mirror with `make skills-sync`.
5. MUST run `make baseline-digests`, then re-record the two characterization
   corpora that command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

6. MUST change only `.agents/skills/roundfix/**`, `skills/roundfix/**`, this
   Task file, and the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`; any
   other path is out of scope — stop rather than widen it.
7. MUST NOT change behaviour. This Task documents what shipped.

## Subtasks

- [ ] Document the refusal, pending, and no-retry contract.
- [ ] Run `make skills-sync`, then `make baseline-digests` and both re-records.

## Acceptance Criteria

- [ ] The Skill states a refusal never verifies a head and blocks the merge.
- [ ] The Skill states an unrecognised signal is pending, and that a green check
      is not evidence a review ran.
- [ ] The Skill states no automatic retry exists yet.
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
- `git diff --name-only HEAD | grep -vE "^(\.agents/skills/roundfix/|skills/roundfix/|docs/specs/0077-a-green-check-is-not-a-review/task_04\.md$|internal/baseline/(assets/(setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Features 1, 2 and 4.
- `_techspec.md` → Integration Points; Build Order 4.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0081.
