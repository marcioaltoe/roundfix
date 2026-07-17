---
task: task_02
spec: 0034-release-plan
status: pending
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

- [ ] Parse the supported Conventional Commit evidence subset.
- [ ] Detect breaking markers in subjects and footers.
- [ ] Encode the maintenance-only changed-path boundary.
- [ ] Aggregate mixed commit evidence by maximum impact.
- [ ] Validate manual impact and reason inputs.
- [ ] Add ambiguous, maintenance-only, mixed-order, and downgrade-rejection tests.

## Acceptance Criteria

- [ ] Breaking evidence outranks compatible features, features outrank fixes, and fixes outrank no-release evidence in every commit order.
- [ ] Documentation-, test-, fixture-, planning-, and CI-only changes qualify for `none` only within the documented boundary.
- [ ] Ambiguous shipped-surface changes produce `manual_classification_required` and identify the blocking commits.
- [ ] A valid manual classification records its reason, participates in maximum-impact selection, and does not imply version approval.
- [ ] Manual impact below an automatic minimum is rejected without producing a partial plan.

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
