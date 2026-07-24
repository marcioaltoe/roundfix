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

## Subtasks

- [ ] Add strict HTTP exception normalization.
- [ ] Derive and merge the Better Auth exception.
- [ ] Persist and re-audit both decision projections.
- [ ] Add conflict and stable-order diagnostics.
- [ ] Add human and automation parity tests.

## Acceptance Criteria

- [ ] The suggested Better Auth value produces one `GET`/`POST`
  `/api/auth/*` exception owned by Better Auth.
- [ ] Reapplying an identical provider decision creates no duplicate exception.
- [ ] A conflicting provider/HTTP pair stops before a complete Plan.
- [ ] Exception ordering and serialized identities are deterministic.
- [ ] Equivalent human and automation answers produce the same Plan Digest.

## Context

- instruction: `docs/adr/0063-repositories-own-their-http-contract.md`
- instruction: `docs/adr/0076-typed-project-decisions-render-identifier-and-authentication-policy.md`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/profile_alignment.go`
- interface: `internal/baseline/assets/decisions.json`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestDeriveBetterAuthHTTPContract|TestHTTPContractConflict|TestProjectDecisionParity|TestProjectDecisionReuse'` — expected: derivation, conflict, ordering, persistence, and interaction parity cases pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1–2 and 5; User Stories 2–3; Core Features 5–9.
- `_techspec.md` → Implementation Design: Data Models and API Contracts; Build Order 2.
- ADR-0063 → repository-owned HTTP policy.
- ADR-0076 → Better Auth decision and derived exception.
