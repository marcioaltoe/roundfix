---
task: task_01
spec: 0034-release-plan
status: completed
type: backend
complexity: medium
---

# Task 01: Build the Release Plan domain and version calculator

## Overview

Establish the deterministic domain vocabulary and semantic-version behavior that every later Release Plan surface uses. This slice is independently verifiable through table-driven domain tests and creates no CLI, Git, configuration, or external-service side effects.

## Requirements

1. MUST model the versioned Release Plan schema, decision states, impact ordering, breaking marker, approval boundary, release references, and per-change evidence defined by the TechSpec.
2. MUST parse only stable `vMAJOR.MINOR.PATCH` release bases and reject malformed or pre-release values with actionable typed errors.
3. MUST calculate patch, compatible minor, version-zero breaking, and major breaking proposals exactly as defined by the Version calculation contract.
4. MUST distinguish `ready`, `approval_required`, `manual_classification_required`, and `no_release` without performing input/output or repository operations.
5. SHOULD keep domain values immutable after construction and use stdlib Go only.

## Subtasks

- [x] Define the Release Plan model and stable enum values.
- [x] Implement stable semantic-version parsing and formatting.
- [x] Implement impact ordering and proposal calculation.
- [x] Model approval questions and version-zero breaking decisions.
- [x] Add table-driven positive and negative domain tests.

## Acceptance Criteria

- [x] Stable release tags round-trip without normalization drift, while malformed and pre-release bases are rejected.
- [x] Patch, minor, version-zero breaking, and major breaking inputs produce the required proposed version and state.
- [x] `none` produces `no_release` without a proposed version.
- [x] Approval-required decisions expose the exact increment kind and proposed version needed by later renderers.
- [x] The schema identifier is `roundfix.release-plan/v1`.

## Context

- instruction: `docs/adr/0048-release-planning-is-read-only-and-confirmation-gated.md`
- interface: `CONTEXT.md`

## Verification

- `go test ./internal/releaseplan -run 'Test(ParseStableVersion|CalculateProposal|ImpactOrdering|ApprovalDecision)' -count=1` — expected: stable-version, impact, approval, and version-zero cases pass.
- `go test ./internal/releaseplan -count=1` — expected: the complete Release Plan domain package passes.

## References

- `_prd.md` → Goals 1 and 3; User Stories 1-2; Core Features 1-2 and 5-6; Success Metrics.
- `_techspec.md` → Data Models; API Contracts: stable values; Version calculation; Build Order 1.
- ADR-0048 → Release planning is read-only and confirmation-gated.

## Result

Implemented the `internal/releaseplan` domain package for the read-only Release Plan model, stable schema values, stable `vMAJOR.MINOR.PATCH` parsing, impact ordering, proposal calculation, and approval decisions. The package performs no I/O, Git, configuration, or external-service work.

Verification:

- `go test ./internal/releaseplan -run 'Test(ParseStableVersion|CalculateProposal|ImpactOrdering|ApprovalDecision)' -count=1` — passed after implementation: 45 focused domain cases.
- `go test ./internal/releaseplan -count=1` — passed after implementation: 46 package cases.
- `make verify` — passed after implementation: Go tests, Python skill tests, `roundfix skills check`, and build all completed.

Acceptance evidence:

- Stable tags round-trip and malformed/pre-release bases reject with typed errors: `TestParseStableVersion`.
- Patch, compatible minor, version-zero breaking, and major breaking proposals match the required state and proposed version: `TestCalculateProposal`.
- `none` returns `no_release` with no proposed version: `TestCalculateProposal/none_produces_no_release`.
- Approval-required proposals expose the exact increment, proposed version, and question: `TestApprovalDecision`.
- The schema identifier is pinned to `roundfix.release-plan/v1`: `TestPlanSchemaVersionAndEvidence`.

Follow-ups:

- Commit classification, manual-impact validation, Git range loading, CLI rendering, and read-only mutation audits remain in later tasks.
