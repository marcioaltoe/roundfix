---
task: task_01
spec: 0034-release-plan
status: pending
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

- [ ] Define the Release Plan model and stable enum values.
- [ ] Implement stable semantic-version parsing and formatting.
- [ ] Implement impact ordering and proposal calculation.
- [ ] Model approval questions and version-zero breaking decisions.
- [ ] Add table-driven positive and negative domain tests.

## Acceptance Criteria

- [ ] Stable release tags round-trip without normalization drift, while malformed and pre-release bases are rejected.
- [ ] Patch, minor, version-zero breaking, and major breaking inputs produce the required proposed version and state.
- [ ] `none` produces `no_release` without a proposed version.
- [ ] Approval-required decisions expose the exact increment kind and proposed version needed by later renderers.
- [ ] The schema identifier is `roundfix.release-plan/v1`.

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
