---
task: task_05
spec: 0047-context-driven-guidance-composition
status: completed
type: backend
complexity: high
---

# Task 05: Plan repository-owned Profile adaptations

## Overview

Resolve a strict custom Profile draft in memory and include its canonical
repository file in the same portable Plan as the Baseline adoption. Required
profile-specific gaps can be removed only through an explicit, catalog-valid
Profile identity.

## Requirements

1. MUST accept either one existing Profile ID or one strict Profile draft and
   reject simultaneous inputs.
2. MUST validate draft identity, safe target path, module dependencies,
   decisions, capabilities, templates, and catalog binding before planning.
3. MUST recalculate templates from confirmed modules and decisions.
4. MUST prevent removal or override of universal required capabilities.
5. MUST include the canonical Profile file as an exact preimage-bound
   repository-owned postimage.
6. MUST reject stale, conflicting, unsafe, or divergent target state without
   partial writes.

## Subtasks

- [x] Add in-memory custom Profile draft resolution.
- [x] Validate allowed adaptation boundaries.
- [x] Assemble the canonical Profile postimage and ledger entry.
- [x] Bind Profile identity into Plan and apply verification.
- [x] Add portable, stale, unsafe, and rollback tests.

## Acceptance Criteria

- [x] A valid backend-only adaptation resolves without writing during planning.
- [x] Universal required capabilities cannot be removed by the draft.
- [x] The Profile file appears in the final file-change and managed-entry
  ledgers.
- [x] Apply verifies the planned Profile bytes and Setup Manifest identity.
- [x] Stale or conflicting target bytes produce no repository mutation.

## Context

- instruction: `docs/adr/0067-custom-baseline-profiles-are-repository-owned.md`
- instruction: `docs/adr/0071-baseline-plans-are-portable-and-preimage-bound.md`
- instruction: `docs/adr/0075-profile-divergence-uses-confirmed-repository-owned-adaptation.md`
- interface: `internal/baseline/custom_profile.go`
- interface: `internal/baseline/plan.go`

## Verification

- `rtk go test -count=1 ./internal/baseline -run 'TestProfileDraftPlan|TestProfileAdaptation|TestProfileDraftRollback|TestProfileDraftStale'` — expected: validation, portable planning, apply, stale input, and recovery cases pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goal 6; User Story 7; Core Features 18–20.
- `_techspec.md` → Implementation Design: Interfaces, Data Models, and API Contracts; Build Order 4.
- ADR-0075 → confirmed repository-owned adaptation.

## Result

Implemented source-bound, in-memory Profile draft planning. Planning now accepts
exactly one existing Profile ID or strict draft, validates the draft's catalog
binding and allowed module, decision, capability, template, and path
boundaries, recalculates templates, and keeps universal capabilities outside
the adaptation surface.

The canonical repository-owned Profile is now an exact preimage-bound
postimage with `profile:<id>` managed-entry and file-change ledger records.
Portable Plan validation binds those bytes to the resolved Profile, while
apply verifies both the Profile postimage and Setup Manifest identity through
the existing recoverable transaction.

Verification:

- `rtk go test -count=1 ./internal/baseline -run 'TestProfileDraftPlan|TestProfileAdaptation|TestProfileDraftRollback|TestProfileDraftStale'` — passed, 17 tests.
- `rtk make verify` — passed: 2,170 Go tests across 22 packages, 4 skill-contract tests, Roundfix skill check, and CLI build.

Acceptance evidence:

- `TestProfileDraftPlanIncludesCanonicalRepositoryProfile` proves a
  backend-only draft resolves without creating its target during planning and
  validates in a matching clone.
- `TestProfileAdaptationCannotRemoveUniversalRequiredCapabilities` and
  `TestProfileAdaptationRejectsInvalidDraftBoundaries` prove universal
  requirements remain active and cannot be declared or overridden by a draft.
- `TestProfileDraftPlanIncludesCanonicalRepositoryProfile` checks the exact
  Profile postimage, `profile:<id>` managed entry, file change, Profile digest,
  and Setup Manifest identity.
- `TestProfileAdaptationApplyVerifiesProfileAndManifest` proves apply installs
  and reports the exact planned Profile bytes and matching Setup Manifest
  identity.
- `TestProfileDraftPlanRejectsSimultaneousAndConflictingInputs`,
  `TestProfileDraftPlanRejectsUnsafeTargetParent`,
  `TestProfileDraftStaleTargetProducesNoMutation`, and
  `TestProfileDraftRollbackRestoresMissingProfile` prove conflicts, unsafe
  targets, stale bytes, and verification failures leave repository state
  unchanged.

Follow-ups: none.
