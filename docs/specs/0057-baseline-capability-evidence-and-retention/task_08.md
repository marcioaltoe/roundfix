---
task: task_08
spec: 0057-baseline-capability-evidence-and-retention
status: completed
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

## Result

Implemented same-identity retention accounting in `internal/baseline/` without
changing the fail-closed apply transaction, digest confirmation, or preimage
binding. Planning now binds the previous managed-clause inventory to the exact
Baseline, Profile digest, and catalog digest tuple; builds one exhaustive
clause delta; and returns action-required without a Plan when an exact source
inventory proves that any previous clause is unaccounted. A ready drift plan
carries both its non-empty retention evidence and clause delta, and portable
Plan validation rejects empty, incomplete, or unaccounted deltas.

Focused checks:

- The new focused tests first failed against the pre-change behavior: the
  disappearing-clause fixture returned ready, and the accounted drift fixture
  carried an empty retention ledger. This established the intended regression
  signal after rerunning with a writable `GOCACHE` (the default user cache is
  unavailable in this sandbox).
- `rtk env GOCACHE=/private/tmp/roundfix-task08-gocache go test
  ./internal/baseline -run
  'Test(SameIdentityDriftRequiresRetention|ReadyPlanNeverCarriesEmptyLedger|PlanDocument|ProfileDraftPlan)$'
  -count=1` passed (`ok roundfix/internal/baseline`, 1.703s).
- `rtk env GOCACHE=/private/tmp/roundfix-task08-gocache go test
  ./internal/baseline -run
  'Test(EmbeddedCatalog|PlanDocumentStrictCodecs|PlanDocumentIncludesMaintainedUpgradeRetention|SameIdentityDriftRequiresRetention|ReadyPlanNeverCarriesEmptyLedger)$'
  -count=1` passed (`ok roundfix/internal/baseline`, 1.622s), covering catalog
  source capture, strict codecs, the existing maintained upgrade path, both new
  drift outcomes, and rejection of an externally supplied empty delta.
- `rtk git diff --check` passed.

Acceptance criterion evidence:

- `TestSameIdentityDriftRequiresRetention` applies an initial Baseline, keeps
  its identity, changes the catalog digest, removes one known clause, and
  asserts action-required with `1 unaccounted clause` and that clause's stable
  ID in the message.
- `TestReadyPlanNeverCarriesEmptyLedger` changes a known retained clause and
  asserts a ready Plan has one retained ledger entry and one retained clause
  disposition. It also proves Plan validation rejects a non-nil empty delta.
- The unaccounted fixture asserts `outcome.Plan == nil`, so no apply operation
  is offered while the delta contains an unaccounted clause.
- Classification walks the exact previous managed-clause inventory once into a
  map keyed by clause ID and a seven-disposition count map. The fixture asserts
  both previous clauses receive exactly one disposition and the retained plus
  unaccounted counts are exact.
- Digest-only drift with unchanged managed artifacts preserves the existing
  ready path. Changed artifacts are classified only when an exact source tuple
  supplies positive evidence; action-required messages enumerate only clauses
  proven absent by stable identity and enforcement. The full characterization
  corpus remains for Daemon Verification and was not run in this Agent turn.
- `rtk git -c core.fsmonitor=false status --porcelain` showed only this task
  file and files under `internal/baseline/`.

The Task's declared `## Verification` commands were not run; the Daemon owns
that gate.
