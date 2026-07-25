---
task: task_01
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: completed
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

- [x] Confirm the external Spec 0047 QA gate.
- [x] Add identifier-strategy declaration and validation.
- [x] Add conditional authentication-provider declaration and validation.
- [x] Connect human suggestions and strict automation input.
- [x] Add reuse, missing-input, and invalid-shape tests.

## Acceptance Criteria

- [x] UUID version 7 appears as a human keep-or-change suggestion.
- [x] Automation without `identifier.strategy` exits action-required with no
  partial Plan.
- [x] Better Auth Profiles require the complete `auth.provider` object.
- [x] Profiles without Better Auth omit the provider decision.
- [x] Compatible stored objects are reused and invalid objects fail closed.

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

## Result

The embedded catalog now declares strict `identifier.strategy` and
`auth.provider` decision types. The Standard TypeScript Monorepo Profile
selects the identifier decision and selects the authentication-provider
decision only while it retains `capability.stack.better-auth`. Human input
shows explicit UUID version 7 and complete Better Auth suggestions, while
automation still has to provide every unresolved selected value.

### Acceptance evidence

- UUID version 7 appears in the `identifier.strategy` keep-or-change prompt,
  including the complete `{"kind":"uuid-v7"}` suggestion.
- `TestIdentifierStrategyDecision` confirms that missing automation input
  returns `action_required` with a nil Plan.
- `TestProjectDecisionValidation` admits only the two identifier
  discriminators and the complete Better Auth object. Missing, conflicting,
  duplicate, unsupported, and unknown fields fail validation.
- `TestAuthProviderDecision` confirms that a Better Auth Profile requires
  `auth.provider`, while a reviewed Profile adaptation without Better Auth
  omits it and rejects attempts to retain it.
- `TestProjectDecisionPrompts` confirms that compatible Setup Manifest values
  remain visible as keep-or-change choices and invalid stored objects are not
  reused.

### Verification

- `rtk rg -q '^verdict: pass$' docs/specs/0047-context-driven-guidance-composition/qa`
  passed after the newest report,
  `qa/qa-report-2026-07-25.md`, was inspected with `verdict: pass`.
- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestIdentifierStrategyDecision|TestAuthProviderDecision|TestProjectDecisionValidation|TestProjectDecisionPrompts'`
  passed: 18 tests in 2 packages.
- `rtk make verify` passed: 2,226 Go tests, 4 skill tests, Roundfix skill
  synchronization, and the Roundfix build.

### Follow-up

Task 02 remains responsible for deterministic Better Auth-to-HTTP Contract
derivation and cross-decision conflict detection.
