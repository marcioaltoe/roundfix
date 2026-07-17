---
task: task_02
spec: 0034-release-plan
status: completed
type: backend
complexity: high
---

# Task 02: Classify commits and maintenance-only changes

## Overview

Turn normalized committed changes into conservative, reproducible Release Plan evidence. This slice proves automatic Conventional Commit minimums, maintenance-only no-release decisions, ambiguity handling, and validated manual classification without reading Git directly.

## Requirements

1. MUST classify breaking markers, `feat`, `fix`, and `perf` evidence from complete commit subjects and bodies using the minimum impacts defined by the TechSpec.
2. MUST leave other or malformed commit types without an automatic increment instead of guessing their semantic impact.
3. MUST classify an otherwise ambiguous commit as `none` only when every changed path stays inside the documented maintenance-only boundary.
4. MUST require manual classification for ambiguous commits that touch shipped, configured, distributed, release, skill, CLI, or public compatibility surfaces.
5. MUST validate that manual impact has a non-empty reason and never falls below the automatic minimum.
6. MUST retain evidence for every commit and select the highest required impact independent of commit order.
7. SHOULD express classification failures as actionable domain errors suitable for CLI mapping.

## Subtasks

- [x] Parse the supported Conventional Commit evidence subset.
- [x] Detect breaking markers in subjects and footers.
- [x] Encode the maintenance-only changed-path boundary.
- [x] Aggregate mixed commit evidence by maximum impact.
- [x] Validate manual impact and reason inputs.
- [x] Add ambiguous, maintenance-only, mixed-order, and downgrade-rejection tests.

## Acceptance Criteria

- [x] Breaking evidence outranks compatible features, features outrank fixes, and fixes outrank no-release evidence in every commit order.
- [x] Documentation-, test-, fixture-, planning-, and CI-only changes qualify for `none` only within the documented boundary.
- [x] Ambiguous shipped-surface changes produce `manual_classification_required` and identify the blocking commits.
- [x] A valid manual classification records its reason, participates in maximum-impact selection, and does not imply version approval.
- [x] Manual impact below an automatic minimum is rejected without producing a partial plan.

## Context

- instruction: `docs/adr/0048-release-planning-is-read-only-and-confirmation-gated.md`
- interface: `cog.toml`
- interface: `docs/findings/2026-07-16-release-version-strategy-and-approval-gates.md`

## Verification

- `go test ./internal/releaseplan -run 'Test(ClassifyCommit|ClassifyChanges|MaintenanceOnly|ValidateManualImpact)' -count=1` — expected: automatic, ambiguous, maintenance-only, mixed-order, and manual cases pass.
- `go test ./internal/releaseplan -count=1` — expected: all classifier and domain tests pass together.

## References

- `_prd.md` → Goals 1-3; User Stories 1, 3, and 5; Core Features 2-6; Success Metrics.
- `_techspec.md` → Data Models: ChangeEvidence; API Contracts: automatic classification and maintenance-only boundary; Build Order 2.
- ADR-0048 → Release planning is read-only and confirmation-gated.

## Result

Implemented commit classification in `internal/releaseplan` without Git, repository, configuration, or external-service access. The classifier now parses the supported Conventional Commit subset, detects subject and footer breaking markers, applies the documented maintenance-only path boundary, retains evidence for every commit, aggregates maximum impact independent of commit order, and validates manual classifications through typed errors.

Verification:

- `go test ./internal/releaseplan -run 'Test(ClassifyCommit|ClassifyChanges|MaintenanceOnly|ValidateManualImpact)' -count=1` — passed after implementation: 47 focused classifier cases.
- `go test ./internal/releaseplan -count=1` — passed after implementation: 93 package cases.
- `make verify` — passed after implementation: Go tests, Python skill tests, `roundfix skills check`, and build all completed.

Acceptance evidence:

- Breaking, feature, fix/perf, and no-release ordering is covered in `TestClassifyChanges`.
- Documentation planning paths, tests, fixtures, and non-release CI boundaries are covered in `TestMaintenanceOnly`.
- Ambiguous shipped-surface commits return `manual_classification_required` and blocking commit SHAs in `TestClassifyChanges/ambiguous_shipped_surface_requires_manual_classification`.
- Manual classifications record the reason, participate in max-impact selection, and do not set approval in `TestClassifyChanges/valid_manual_classification_resolves_ambiguity`.
- Manual downgrades and invalid manual input return typed errors with a zero result in `TestValidateManualImpact` and `TestClassifyChangesRejectsInvalidManualImpactWithoutPartialResult`.

Follow-ups:

- Git range loading, proposal construction from classification output, CLI exit-code mapping, and renderer behavior remain in later tasks.
