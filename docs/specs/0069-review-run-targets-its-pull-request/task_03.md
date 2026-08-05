---
task: task_03
spec: 0069-review-run-targets-its-pull-request
status: pending
type: chore
complexity: low
---

# Task 03: Synchronise the Roundfix Skill

## Overview

This Spec adds a Preflight refusal and a terminal outcome, so the Skill must
teach both before the Spec can close. This is the authorized tooling Task.

## Requirements

1. MUST document that a review Run validates the checkout against the Pull
   Request's head branch at Preflight, and refuses with exit `2` naming both
   branches and both revisions.
2. MUST document that the refusal creates no Run and has no side effects, and
   name the recovery command shape.
3. MUST document the terminal interruption outcome for a checkout that moves
   mid-Run, stating plainly that it is not a Review Issue failure and that the
   affected issues stay unsettled.
4. MUST state that Roundfix never checks out or moves the working tree.
5. MUST regenerate the `skills/roundfix/**` mirror with `make skills-sync`.
6. MUST run `make baseline-digests`, then re-record the two characterization
   corpora that command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

7. MUST change only `.agents/skills/roundfix/**`, `skills/roundfix/**`, this
   Task file, and the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`; any
   other path is out of scope — stop rather than widen it.
8. MUST NOT change behaviour. This Task documents what shipped.

## Subtasks

- [ ] Document the refusal, the interruption, and the no-checkout rule.
- [ ] Run `make skills-sync`, then `make baseline-digests` and both re-records.

## Acceptance Criteria

- [ ] The Skill states the Preflight target validation and its exit `2`.
- [ ] The Skill states the refusal has no side effects and names recovery.
- [ ] The Skill names the interruption outcome and distinguishes it from a
      Review Issue failure.
- [ ] The Skill states Roundfix never moves the working tree.
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
- `git diff --name-only HEAD | grep -vE "^(\.agents/skills/roundfix/|skills/roundfix/|docs/specs/0069-review-run-targets-its-pull-request/task_03\.md$|internal/baseline/(assets/(setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Features 2, 4 and 5.
- `_techspec.md` → Integration Points; Build Order 3.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0081.
