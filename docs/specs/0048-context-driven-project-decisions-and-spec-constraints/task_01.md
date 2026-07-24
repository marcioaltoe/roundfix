---
task: task_01
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: pending
type: backend
complexity: high
---

# Task 01: Model typed project decisions

## Overview

Add strict identifier-strategy and authentication-provider decision contracts
to the embedded Baseline catalog. This slice starts only after the newest Spec
0047 QA Report passes and makes unresolved values explicit for humans and
automation.

## Requirements

1. MUST verify that the newest Spec 0047 QA Report has `verdict: pass` before
   implementation begins.
2. MUST admit only the confirmed discriminated identifier variants and strict
   fields.
3. MUST select `auth.provider` only when the resolved Profile retains the
   Better Auth capability.
4. MUST expose UUID version 7 and the complete Better Auth value as visible
   human suggestions, never automation answers.
5. MUST require explicit automation values for every unresolved selected
   decision.
6. MUST reuse compatible Setup Manifest values and reject incomplete,
   conflicting, or unknown fields.

## Subtasks

- [ ] Confirm the external Spec 0047 QA gate.
- [ ] Add identifier-strategy declaration and validation.
- [ ] Add conditional authentication-provider declaration and validation.
- [ ] Connect human suggestions and strict automation input.
- [ ] Add reuse, missing-input, and invalid-shape tests.

## Acceptance Criteria

- [ ] UUID version 7 appears as a human keep-or-change suggestion.
- [ ] Automation without `identifier.strategy` exits action-required with no
  partial Plan.
- [ ] Better Auth Profiles require the complete `auth.provider` object.
- [ ] Profiles without Better Auth omit the provider decision.
- [ ] Compatible stored objects are reused and invalid objects fail closed.

## Context

- instruction: `docs/adr/0076-typed-project-decisions-render-identifier-and-authentication-policy.md`
- interface: `internal/baseline/assets/decisions.json`
- interface: `internal/baseline/custom_profile.go`
- interface: `internal/baseline/plan.go`
- interface: `internal/cli/baseline_human.go`

## Verification

- `rtk rg -q '^verdict: pass$' docs/specs/0047-context-driven-guidance-composition/qa` — expected: the newest Spec 0047 QA Report inspected by the Task carries passing evidence before this Task can settle.
- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestIdentifierStrategyDecision|TestAuthProviderDecision|TestProjectDecisionValidation|TestProjectDecisionPrompts'` — expected: strict values, conditional selection, suggestions, reuse, and missing-input cases pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1–2 and 5; User Stories 1 and 3; Core Features 1–7.
- `_techspec.md` → System Architecture; Implementation Design: Interfaces and Data Models; Build Order 1.
- ADR-0076 → typed identifier and authentication decisions.
