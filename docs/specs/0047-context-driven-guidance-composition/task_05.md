---
task: task_05
spec: 0047-context-driven-guidance-composition
status: pending
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

- [ ] Add in-memory custom Profile draft resolution.
- [ ] Validate allowed adaptation boundaries.
- [ ] Assemble the canonical Profile postimage and ledger entry.
- [ ] Bind Profile identity into Plan and apply verification.
- [ ] Add portable, stale, unsafe, and rollback tests.

## Acceptance Criteria

- [ ] A valid backend-only adaptation resolves without writing during planning.
- [ ] Universal required capabilities cannot be removed by the draft.
- [ ] The Profile file appears in the final file-change and managed-entry
  ledgers.
- [ ] Apply verifies the planned Profile bytes and Setup Manifest identity.
- [ ] Stale or conflicting target bytes produce no repository mutation.

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
