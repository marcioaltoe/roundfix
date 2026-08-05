---
task: task_05
spec: 0066-run-teardown-reclaims-what-it-created
status: pending
type: chore
complexity: low
---

# Task 05: Synchronise the Roundfix Skill with the new contract

## Overview

The repository HARD RULE requires a Pull Request that changes CLI behaviour to
ship the Roundfix Skill update with it. This Spec changes what `reconcile`
offers and when Branch Integrity Preflight refuses, so the Skill must teach
both before the Spec can close. This is the authorized tooling Task.

## Requirements

1. MUST document the two new `reconcile` candidate kinds, their proofs, and
   that dry-run remains the default.
2. MUST document the Branch Integrity Preflight change: proven-superseded Run
   Branch work no longer blocks a review Run, and every other refusal stands.
3. MUST state that an unprovable termination is reported, never treated as
   success, so a reader does not mistake silence for a stopped process.
4. MUST regenerate the `skills/roundfix/**` mirror with `make skills-sync`.
5. MUST run `make baseline-digests`, then re-record the two characterization
   corpora that command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

   Neither test is in `BASELINE_DIGEST_STEPS`, so the sanctioned command reports
   "no changes" while the gate stays red. The flags are verbatim because they do
   not match their test names.
6. MUST change only `.agents/skills/roundfix/**`, `skills/roundfix/**`, this
   Task file, and the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`; any
   other path is out of scope — stop rather than widen it.
7. MUST NOT change behaviour. This Task documents what shipped.

## Subtasks

- [ ] Document the reconcile candidates and the Preflight change.
- [ ] Document the unprovable-termination reporting.
- [ ] Run `make skills-sync`, then `make baseline-digests` and both re-records.

## Acceptance Criteria

- [ ] The Skill documents both new candidate kinds and the dry-run default.
- [ ] The Skill documents the Preflight relaxation and states the other
      refusals stand.
- [ ] The Skill states an unprovable termination is never reported as success.
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
- `git diff --name-only HEAD | grep -vE "^(\.agents/skills/roundfix/|skills/roundfix/|docs/specs/0066-run-teardown-reclaims-what-it-created/task_05\.md$|internal/baseline/(assets/(setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Features 2, 4, 5.
- `_techspec.md` → Integration Points; Build Order 5.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md` → standing
  Roundfix Skill CLI synchronisation grant.
- ADR-0081.
