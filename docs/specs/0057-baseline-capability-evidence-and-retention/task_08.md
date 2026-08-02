---
task: task_08
spec: 0057-baseline-capability-evidence-and-retention
status: pending
type: backend
complexity: high
---

# Task 08: Account for clauses when a Profile drifts

## Overview

A repository whose Profile or catalog digests changed under an unchanged
Baseline identifier is treated as a refresh, not an upgrade, so it bypasses
retention accounting entirely and managed Normative Clauses can disappear with
an empty ledger. This is the path ADR-0058 was written to close and the one
defect in this Spec that lets live consumer repositories lose rules silently.
It is also the only slice that turns a completing plan into a stopping one.

## Requirements

1. MUST treat a matching Baseline identifier with changed Profile or catalog
   digests as requiring a retention transition keyed by the source tuple.
2. MUST classify every previous managed clause with exactly one disposition:
   retained, moved, replaced, repository-document, repository-extension,
   reasoned-rejection, or unaccounted.
3. MUST exit action-required when any clause is unaccounted, and MUST NOT offer
   apply while any clause is unaccounted.
4. MUST NOT allow a ready update plan to carry an empty retention ledger when
   clauses changed.
5. MUST fail closed on evidence, never on uncertainty: a plan that is
   legitimately ready today stays ready, and a clause is unaccounted only when
   provably so.
6. MUST leave the fail-closed apply, digest confirmation, and preimage binding
   untouched.

## Subtasks

- [ ] Detect same-identity drift from Profile and catalog digests.
- [ ] Classify every previous clause into one disposition.
- [ ] Exit action-required with the unaccounted count.
- [ ] Withhold apply while any clause is unaccounted.
- [ ] Confirm legitimately ready plans stay ready.

## Acceptance Criteria

- [ ] A fixture with unchanged identity, changed digests, and a disappearing
      clause exits action-required and states the unaccounted count.
- [ ] No ready update plan carries an empty retention ledger when clauses
      changed.
- [ ] Apply is not offered while any clause is unaccounted.
- [ ] Every previous clause receives exactly one disposition.
- [ ] Every plan in the characterization corpus that is ready today is still
      ready, and every new action-required outcome names a provably unaccounted
      clause.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/source_contracts.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestSameIdentityDriftRequiresRetention -count=1`
  — expected: exit 0; the disappearing-clause fixture exits action-required.
- `go test ./internal/baseline -run TestReadyPlanNeverCarriesEmptyLedger -count=1`
  — expected: exit 0.
- `go test ./internal/baseline -run TestBaselinePlanCharacterization -count=1` —
  expected: exit 0; every plan ready today is still ready.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 1; Core Features 1 and 11; Success Metrics.
- `_techspec.md` → Implementation Design: Interfaces; Risks; Build Order 6.
- ADR-0058.
