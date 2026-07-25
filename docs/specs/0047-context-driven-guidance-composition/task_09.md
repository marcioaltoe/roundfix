---
task: task_09
spec: 0047-context-driven-guidance-composition
status: completed
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

- [x] Add all-profile greenfield and update macro fixtures.
- [x] Add semantic redistribution and residual fixtures.
- [x] Add Profile adaptation and universal-remediation fixtures.
- [x] Add formatter, audit, reapply, and rollback assertions.
- [x] Define the final live QA evidence matrix.

## Acceptance Criteria

- [x] Every maintained Profile completes its hermetic journey with zero
  second-pass delta.
- [x] Zero residual rules produce no generic or repository-specific carrier.
- [x] The backend-only fixture reaches a valid repository-owned Profile without
  waivers.
- [x] Injected failure restores all affected files, including a planned Profile
  file.
- [x] The Spec-local QA plan names fresh Fluxus greenfield/update and Oraculum
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

## Result

Implemented hermetic release-gate journeys at the existing Baseline engine and
public CLI boundaries.

- `TestGuidanceCompositionJourney` enumerates the embedded catalog's complete
  maintained Profile set (`go-cli-tui`, `rust-cli`, and
  `standard-typescript-monorepo`). Each Profile completes greenfield apply,
  update apply, external formatting, repository Verification, fresh empty
  planning, and mutation-free reapply through a built Roundfix binary.
- `TestSemanticRedistributionJourney` moves exact bytes from both legacy
  generic carriers into the active CLI semantic guide, proves zero-residual
  cleanup, separately proves canonical residual retention, and validates the
  complete managed-entry and Upgrade Retention Contract ledgers.
- `TestProfileAdaptationJourney` reproduces a generic backend-only TypeScript
  divergence, plans and applies the reviewed repository-owned Profile through
  `--profile-file`, proves the Profile ledger entry, re-audits universal
  requirements as required and satisfied, and finishes with an empty reapply.
- `TestBaselineReleaseGate` now includes semantic-guide and planned-Profile
  transaction failures. Both injected failures compare the complete visible
  repository tree with its exact preimage after rollback.
- `qa/live-journeys.md` defines the separately authorized evidence matrix for
  fresh Fluxus greenfield and update journeys plus the Oraculum divergence and
  adaptation journey. Task 09 did not read or mutate either live repository.

Verification:

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestGuidanceCompositionJourney|TestSemanticRedistributionJourney|TestProfileAdaptationJourney|TestBaselineReleaseGate'`
  passed with 34 tests across 2 packages.
- `rtk make verify` passed: 2,198 repository tests, 4 skill contract tests,
  `roundfix skills check`, and the Roundfix build.
- `rtk git diff --check` passed.

Live QA remains pending by design and requires a separately authorized
`qa-gate` run using the Spec-local evidence matrix.
