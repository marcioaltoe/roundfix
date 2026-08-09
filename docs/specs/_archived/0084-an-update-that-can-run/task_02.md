---
task: task_02
spec: 0084-an-update-that-can-run
status: completed
type: backend
complexity: medium
---

# Task 02: Name the lines a refresh removes

## Overview

Computes, for every unrecorded managed region, the lines present on disk that the
refreshed rendering does not reproduce, and carries them on the portable plan
document. This is what makes replacing an unrecorded region safe to approve: the
maintainer sees what disappears before approving, and the value is inside the
Plan Digest so approval covers what was shown. Plan construction is the only
place where both the preimage and the rendered postimage exist.

## Requirements

1. MUST compute, per unrecorded managed region, the on-disk lines the refreshed
   rendering does not reproduce, comparing trimmed non-blank lines, preserving
   on-disk order, and deduplicating repeats.
2. MUST report an empty removed-line set for a region whose refresh only adds
   lines, rather than omitting the region.
3. MUST truncate the reported lines at a fixed bound and record how many were
   omitted, so one large region cannot dominate the output.
4. MUST carry the removed lines on the portable plan document only when the
   preservation mode is managed refresh, leaving greenfield and preservation plan
   JSON byte-identical.
5. MUST include the carried value in the Plan Digest, so a digest approved in a
   previous invocation covers the removed lines that were presented with it.
6. MUST round-trip through the portable plan document's strict codecs without
   loss.
7. MUST leave the preimage-bound non-managed region proof from ADR-0100 intact,
   so applying a plan still verifies that every byte outside a managed marker is
   unchanged.

## Subtasks

- [x] Compute removed lines from the preimage and the rendered postimage.
- [x] Apply order preservation, deduplication, and bounded truncation.
- [x] Carry the value on the portable plan document for the managed-refresh mode.
- [x] Include the value in the Plan Digest.
- [x] Cover a region that removes a line and a region that removes none.
- [x] Cover truncation and its reported remainder.
- [x] Cover strict-codec round-trip and unchanged non-managed region proof.

## Acceptance Criteria

- [x] A managed-refresh plan over a region whose rendering drops a specific line
      reports exactly that line as removed.
- [x] A managed-refresh plan over a region whose rendering only adds lines reports
      the region with an empty removed-line set.
- [x] A region whose removed lines exceed the bound reports the bounded list and a
      non-zero omitted count equal to the difference.
- [x] Two plans over the same repository that differ only in a removed line have
      different Plan Digests.
- [x] A greenfield plan and a preservation plan over the same fixtures produce
      byte-identical JSON to the recorded goldens.
- [x] Encoding and decoding a managed-refresh plan document preserves every
      reported region and line.

## Context

- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/preservation.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/baseline/ -run 'RemovedLines' -v > /tmp/0084-task-02-a.log 2>&1 && grep -q '^--- PASS: .*RemovedLines' /tmp/0084-task-02-a.log` — expected: exits 0, proving the removed-line cases exist and pass rather than being selected out.
- `go test ./internal/baseline/ -run 'PlanDigest' -v > /tmp/0084-task-02-b.log 2>&1 && grep -q '^--- PASS: .*PlanDigest' /tmp/0084-task-02-b.log` — expected: exits 0, proving the digest coverage case ran.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0, with the greenfield and preservation plan goldens unchanged.

## References

- `_techspec.md` → Build Order 2; Data Models; Interfaces: `UnrecordedManagedRegion`.
- `_prd.md` → Core Feature 2; User Story 3; Goal 3.
- ADR-0102, ADR-0100.

## Result

### Implementation

- Managed-refresh plan construction now compares each unrecorded Managed
  Region's on-disk and rendered bodies as trimmed non-blank lines. It keeps the
  on-disk order, removes duplicate report entries, and carries an explicit empty
  array when the refresh removes nothing.
- Each region reports at most 50 removed lines and records the remaining unique
  line count in `removedLinesTruncated`.
- The portable Plan Document carries `unrecordedManagedRegions` only in
  `managed-refresh` mode. The field participates in the existing Plan Digest
  payload and is checked by the strict Plan Document validator and codecs.
- The ADR-0100 preimage-bound non-managed-region validation path was not changed.

### Focused checks

- Pre-change signal:
  `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260808T153649Z_78746d4b80d08fc7.task_02/.gocache go test ./internal/baseline -run '^TestManagedRefreshPlanReportsRemovedLines$' -count=1`
  failed to compile because `UnrecordedManagedRegion.RemovedLines`,
  `RemovedLinesTruncated`, and the fixed bound did not exist.
- After implementation:
  `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260808T153649Z_78746d4b80d08fc7.task_02/.gocache go test ./internal/baseline -run '^TestManagedRefresh(PlanReportsRemovedLines|PlanReportsEmptyRemovedLines|PlanTruncatesRemovedLines|RemovedLinesPlanDigestAndStrictCodecRoundTrip)$' -count=1`
  passed: `ok roundfix/internal/baseline 1.830s`.
- Fresh regression check after the final test edit:
  `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260808T153649Z_78746d4b80d08fc7.task_02/.gocache go test ./internal/baseline -run '^(TestManagedRefreshPlanReportsRemovedLines|TestManagedRefreshPlanReportsEmptyRemovedLines|TestManagedRefreshPlanTruncatesRemovedLines|TestManagedRefreshRemovedLinesPlanDigestAndStrictCodecRoundTrip|TestNonManagedPlansOmitRemovedLines|TestManagedRefreshPreservesNonManagedRegionDigests|TestBaselinePlanCharacterization)$' -count=1`
  passed: `ok roundfix/internal/baseline 3.368s`.

### Acceptance evidence

1. `TestManagedRefreshPlanReportsRemovedLines` reports exactly
   `line removed by refresh` after trimming blanks and deduplicating the repeated
   on-disk line.
2. `TestManagedRefreshPlanReportsEmptyRemovedLines` ages only the recorded digest
   and asserts that the carried region has a non-nil, zero-length removed-line
   array.
3. `TestManagedRefreshPlanTruncatesRemovedLines` supplies 53 unique removed lines
   and asserts the first 50 remain in on-disk order with a truncated count of 3.
4. `TestManagedRefreshRemovedLinesPlanDigestAndStrictCodecRoundTrip` changes only
   one carried removed line and observes a different computed Plan Digest.
5. `TestNonManagedPlansOmitRemovedLines` asserts both real greenfield and
   preservation Plan Documents omit the managed-refresh fields;
   `TestBaselinePlanCharacterization` passed against the recorded plan goldens.
6. `TestManagedRefreshRemovedLinesPlanDigestAndStrictCodecRoundTrip` uses
   `MarshalPlanDocument` and `ParsePlanDocument` and asserts every carried region
   and line survives unchanged.
7. `TestManagedRefreshPreservesNonManagedRegionDigests` passed after the change,
   exercising the existing exact non-managed-region digest proof from ADR-0100.

The Daemon-owned `## Verification` commands were not run in this Agent turn.
