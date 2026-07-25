---
task: task_08
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: completed
type: test
complexity: high
---

# Task 08: Prove project-constraint journeys

## Overview

Create hermetic macro journeys for typed decisions, rendered guidance, Spec
authoring, tooling refusal, apply, audit, and empty reapply. Separately
authorized Fluxus journeys remain fresh final QA evidence.

## Requirements

1. MUST cover human and automation decision collection for every affected
   maintained Profile.
2. MUST cover compatible reuse, unresolved input, invalid objects, derived
   exception conflict, and stable Plan Digests.
3. MUST author one new PRD and TechSpec fixture with complete Project
   Constraints and one authorized tooling fixture with exact bounded files.
4. MUST prove decomposition, execution, and QA refusal for missing or exceeded
   tooling authorization.
5. MUST complete formatter, apply, repository Verification recommendation,
   audit, and empty reapply with zero managed delta.
6. MUST define fresh Fluxus greenfield and update evidence for final
   `qa-gate`.

## Subtasks

- [ ] Add all-profile decision journey fixtures.
- [ ] Add decision reuse and conflict journeys.
- [ ] Add PRD and TechSpec authoring fixtures.
- [ ] Add tooling authorization and refusal journeys.
- [ ] Add formatter, audit, reapply, and final QA assertions.

## Acceptance Criteria

- [ ] Equivalent human and automation answers produce identical Plans.
- [ ] Missing or conflicting structured values stop without mutation.
- [ ] New Spec fixtures contain every required Project Constraint row and
  source.
- [ ] Tooling mutation outside bounded authorization fails before settlement.
- [ ] Every affected Profile completes apply, audit, and empty reapply with
  zero managed delta.
- [ ] The QA matrix requires fresh Fluxus greenfield and update evidence.

## Context

- instruction: `docs/adr/0076-typed-project-decisions-render-identifier-and-authentication-policy.md`
- instruction: `docs/adr/0077-new-specs-carry-a-readable-project-constraint-snapshot.md`
- interface: `internal/baseline/compatibility_corpus_test.go`
- interface: `internal/baseline/release_gate_test.go`
- interface: `internal/cli/baseline_release_gate_test.go`
- interface: `internal/baseline/testdata/parity-corpus/v1`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli ./skills -run 'TestProjectDecisionJourney|TestProjectConstraintJourney|TestToolingAuthorizationJourney|TestBaselineReleaseGate'` — expected: decision, authoring, refusal, formatter, apply, audit, and reapply journeys pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1–5; User Stories 1–6; Core Features 1–17; Success Metrics.
- `_techspec.md` → Testing Approach; Build Order 7; Risks & Considerations.
- ADR-0076 and ADR-0077 → final project-decision and Spec-constraint contracts.
