---
task: task_02
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: pending
type: backend
complexity: high
---

# Task 02: Derive Better Auth HTTP policy

## Overview

Project the confirmed Better Auth route exception into the repository-owned
HTTP Contract Decision through one pure deterministic normalization path.
Human and automation inputs produce the same ordered contract or one explicit
conflict.

## Requirements

1. MUST derive the provider exception from `auth.provider` before Plan digest
   calculation.
2. MUST merge exceptions by normalized owner and scope, treating an identical
   value as idempotent.
3. MUST reject duplicate methods, unsupported methods, missing rationale, and
   conflicting values while naming both decision IDs.
4. MUST persist the confirmed provider value and derived HTTP contract in the
   Setup Manifest.
5. MUST re-derive and compare compatible stored values during every update.
6. MUST preserve exact human and Decision Document parity.
7. MUST make the human update suggestion reuse the exact rationale from a
   compatible persisted Better Auth HTTP exception, so accepting both
   displayed defaults reaches a Plan while an explicitly conflicting
   Decision Document still fails closed.

## Subtasks

- [x] Add strict HTTP exception normalization.
- [x] Derive and merge the Better Auth exception.
- [x] Persist and re-audit both decision projections.
- [x] Add conflict and stable-order diagnostics.
- [x] Add human and automation parity tests.
- [ ] Reconcile the human Better Auth suggestion with a compatible persisted
  HTTP exception.

## Acceptance Criteria

- [x] The suggested Better Auth value produces one `GET`/`POST`
  `/api/auth/*` exception owned by Better Auth.
- [x] Reapplying an identical provider decision creates no duplicate exception.
- [x] A conflicting provider/HTTP pair stops before a complete Plan.
- [x] Exception ordering and serialized identities are deterministic.
- [x] Equivalent human and automation answers produce the same Plan Digest.
- [ ] A compatible Fluxus-style update reaches a Plan when the maintainer
  accepts both displayed defaults without manually editing either rationale.

## Context

- instruction: `docs/adr/0063-repositories-own-their-http-contract.md`
- instruction: `docs/adr/0076-typed-project-decisions-render-identifier-and-authentication-policy.md`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/profile_alignment.go`
- interface: `internal/baseline/assets/decisions.json`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestDeriveBetterAuthHTTPContract|TestHTTPContractConflict|TestProjectDecisionParity|TestProjectDecisionReuse'` — expected: derivation, conflict, ordering, persistence, and interaction parity cases pass.
- `rtk go test -count=1 ./internal/cli -run 'TestBetterAuthSuggestionReusesHTTPReason'`
  — expected: the update suggestion adopts the exact compatible persisted
  rationale while explicit mismatches remain conflicts.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1–2 and 5; User Stories 2–3; Core Features 5–9.
- `_techspec.md` → Implementation Design: Data Models and API Contracts; Build Order 2.
- ADR-0063 → repository-owned HTTP policy.
- ADR-0076 → Better Auth decision and derived exception.

## Result

Baseline planning now normalizes the repository-owned HTTP Contract Decision
and the selected Better Auth decision through one deterministic path before
the Plan Digest is calculated. It trims and orders supported methods, rejects
duplicates and incomplete exceptions, merges identical owner/scope entries
idempotently, and rejects conflicting projections with diagnostics that name
both `auth.provider` and `http.contract`.

The normalized provider and derived HTTP contract are persisted together in
the Setup Manifest. Plan validation re-derives the projection, and the human
update path normalizes stored decisions before offering them for reuse.

### Acceptance evidence

- `TestDeriveBetterAuthHTTPContract` proves that the suggested provider
  produces one Better Auth-owned `GET`/`POST /api/auth/*` exception, preserves
  other HTTP exceptions, orders the serialized contract, and stores both
  normalized decisions in the Setup Manifest.
- `TestDeriveBetterAuthHTTPContract` and `TestProjectDecisionReuse` prove that
  an identical stored or supplied provider exception remains one exception
  after planning, apply, stored-value inspection, and re-planning.
- `TestHTTPContractConflict` proves duplicate normalized methods, unsupported
  methods, missing rationale, and provider/HTTP conflicts stop before a
  complete Plan and name both decision IDs.
- `TestDeriveBetterAuthHTTPContract` asserts the exact normalized serialized
  identity, including scope/owner/method ordering.
- `TestProjectDecisionParity` proves equivalent human answers and reversed
  Decision Document input produce byte-identical Plans and the same Plan
  Digest.

### Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestDeriveBetterAuthHTTPContract|TestHTTPContractConflict|TestProjectDecisionParity|TestProjectDecisionReuse'`
  passed: 8 tests in 2 packages.
- `rtk make verify` passed: 2,234 Go tests in 22 packages, 4 skill-contract
  tests, the Roundfix skill synchronization check, and the Roundfix build.

### Reopened QA repair

The 2026-07-25 QA gate found that the public Fluxus update offered two
individually valid defaults whose Better Auth rationale strings differed.
Accepting both defaults therefore failed before producing a Plan. This Task is
reopened to align only the human suggestion with an already-compatible
persisted exception; explicit automation conflicts remain invalid.
