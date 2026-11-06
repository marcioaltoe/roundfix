---
task: task_03
spec: 0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings
status: pending
type: backend
complexity: medium
---

# Task 03: A guide the repository wrote under docs/agents raises no nested-carrier conflict

## Overview

Every guide a repository writes under `docs/agents/` raises
`baseline.inventory.nested-carrier-conflict`, with or without the marker name
in its prose, because the carrier classification has no kind for a guide the
Baseline does not own; and the colon form `setup-context-driven:` quoted in
prose makes such a guide unclassifiable. This Task adds the
repository-authored guide kind, decided by a marker line rather than the
marker name, and drops its warning. It answers the second defect reported by sql-library
([references](references/2026-11-06-a-repository-profile-without-a-template-decision-fails-opaquely.md)).

## Requirements

1. MUST add to `internal/baseline/classification.go` the kind
   `carrierRepositoryAuthoredGuide` (`"repository-authored-guide"`) and the
   helper `hasMarkerLine`, and give `classifyMarkedCarrier` the branch of
   `_techspec.md` → API Contract 6 for a carrier with no managed block,
   evaluated before its existing substring checks. For every other carrier,
   `classifyMarkedCarrier` keeps its current logic and result.
2. MUST make `warningsForCarrierClassifications` drop
   `baseline.inventory.nested-carrier-conflict` for
   `carrierRepositoryAuthoredGuide` as it does for current and stale managed
   artifacts; no other warning, finding or message changes.
3. MUST add to `internal/baseline/classification_test.go`, each calling
   `t.Parallel()` first:
   - `TestARepositoryAuthoredGuideRaisesNoCarrierConflict`: one temporary
     repository per shape, planned with `buildTestPlan` (or classified with
     `classifyCarriers` as the neighbouring tests do), asserts that
     `docs/agents/sql.md` with plain prose, with the marker name in
     backticks, and with `setup-context-driven:begin` in backticks raises no
     `nested-carrier-conflict`, while `services/payments/AGENTS.md` with plain
     prose, `docs/agents/docs-layout.md` without markers (a path the catalog
     manages), `docs/agents/sql.md` holding a real begin and end marker
     with an unknown id, and `docs/agents/notes.md` holding a begin marker
     line with no end marker each still raise it;
   - `TestAMarkerLineIsACommentThatStartsALine`: `hasMarkerLine` is true for
     a begin or end comment at the start of a line or after one to three
     spaces, and false after four spaces, inside a code span within a
     sentence, and for the bare name `setup-context-driven`.
4. MUST keep passing, unchanged,
   `TestNestedCarrierRemainsByteIdenticalAndWarningOnly`,
   `TestUnclassifiableCarrierStillWarns`,
   `TestNestedCarrierWarningLeavesNestedSourcesOutOfPreservation`,
   `TestNestedCarrierWarningUnsafeAliasRemainsNonBlocking` and every other
   Baseline test.
5. MUST NOT change the inventory in `internal/baseline/repository.go`, the
   warning message text, any carrier on disk, any asset, golden or guide, or
   any file outside this Task's declarations.

## Subtasks

- [ ] Add the carrier kind, the marker-line helper and the classification branch.
- [ ] Drop the warning for the new kind.
- [ ] Add the two tests.

## Acceptance Criteria

- [ ] A guide under `docs/agents/` without a marker line at an unmanaged path raises no nested-carrier warning, whatever its prose says about the marker.
- [ ] A nested `AGENTS.md`, an unmarked guide at a managed path and a guide with an unknown real marker still warn.
- [ ] The Baseline package passes.

## Context

- instruction: `docs/adr/0372-the-baseline-names-the-archive-destination-a-profiles-missing-decision-and-a-guide-it-does-not-own.md`
- instruction: `internal/baseline/repository.go`
- instruction: `internal/baseline/preservation.go`
- instruction: `internal/baseline/preservation_test.go`
- interface: `internal/baseline/classification.go`
- interface: `internal/baseline/classification_test.go`

## Consumers

- The carrier classification and its warning filter — consumers: `./internal/baseline`

## Verification

- `for t in TestARepositoryAuthoredGuideRaisesNoCarrierConflict TestAMarkerLineIsACommentThatStartsALine; do grep -q "func $t(" internal/baseline/classification_test.go || { printf 'missing %s\n' "$t" >&2; exit 1; }; done; go test -timeout 20m -count=1 -run '^(TestARepositoryAuthoredGuideRaisesNoCarrierConflict|TestAMarkerLineIsACommentThatStartsALine|TestNestedCarrierRemainsByteIdenticalAndWarningOnly|TestUnclassifiableCarrierStillWarns|TestNestedCarrierWarningLeavesNestedSourcesOutOfPreservation|TestNestedCarrierWarningUnsafeAliasRemainsNonBlocking)$' ./internal/baseline && go test -timeout 20m -count=1 ./internal/baseline` — expected: exit 0. Before this Task the two tests do not exist, so the first check fails; afterwards a repository-authored guide raises no warning, every other carrier keeps its warning and the Baseline suite passes.

## References

- `_prd.md` → Goals; User Story 3; Core Feature 3; Success Metric 5; Acceptance evidence
- `_techspec.md` → API Contract 6; Testing Approach; Build Order 3
- ADR-0372
