---
task: task_09
spec: 0047-context-driven-guidance-composition
status: pending
type: test
complexity: high
---

# Task 09: Prove composed Baseline journeys

## Overview

Create hermetic macro journeys that prove the composed guidance contract across
greenfield, update, Readoption, Profile adaptation, formatting, apply, audit,
and empty reapply. Live Fluxus and Oraculum checks remain final QA evidence,
not Task prerequisites.

## Requirements

1. MUST cover every affected maintained Profile through real CLI planning and
   apply boundaries.
2. MUST cover legacy generic carriers, semantic redistribution, residual
   retention, and zero-residual cleanup.
3. MUST reproduce the backend-only TypeScript divergence and reviewed Profile
   adaptation without project-specific fixture branding.
4. MUST prove formatter, repository Verification recommendation, audit, and
   empty reapply composition.
5. MUST assert complete managed-entry and Upgrade Retention ledgers plus exact
   rollback after injected failure.
6. MUST leave Fluxus and Oraculum journeys to separately authorized
   `qa-gate` execution and name their required evidence.

## Subtasks

- [ ] Add all-profile greenfield and update macro fixtures.
- [ ] Add semantic redistribution and residual fixtures.
- [ ] Add Profile adaptation and universal-remediation fixtures.
- [ ] Add formatter, audit, reapply, and rollback assertions.
- [ ] Define the final live QA evidence matrix.

## Acceptance Criteria

- [ ] Every maintained Profile completes its hermetic journey with zero
  second-pass delta.
- [ ] Zero residual rules produce no generic or repository-specific carrier.
- [ ] The backend-only fixture reaches a valid repository-owned Profile without
  waivers.
- [ ] Injected failure restores all affected files, including a planned Profile
  file.
- [ ] The Spec-local QA plan names fresh Fluxus greenfield/update and Oraculum
  divergence journeys.

## Context

- instruction: `docs/adr/0073-baseline-apply-uses-a-recoverable-multi-file-transaction.md`
- interface: `internal/baseline/compatibility_corpus_test.go`
- interface: `internal/baseline/release_gate_test.go`
- interface: `internal/cli/baseline_release_gate_test.go`
- interface: `internal/baseline/testdata/parity-corpus/v1`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestGuidanceCompositionJourney|TestSemanticRedistributionJourney|TestProfileAdaptationJourney|TestBaselineReleaseGate'` — expected: greenfield, update, adaptation, formatting, rollback, audit, and reapply journeys pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1–6; User Stories 1–7; Core Features 1–20; Success Metrics.
- `_techspec.md` → Testing Approach; Build Order 8; Risks & Considerations.
- ADR-0073 → recoverable multi-file apply.
