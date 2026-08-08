---
task: task_02
spec: 0084-an-update-that-can-run
status: pending
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

- [ ] Compute removed lines from the preimage and the rendered postimage.
- [ ] Apply order preservation, deduplication, and bounded truncation.
- [ ] Carry the value on the portable plan document for the managed-refresh mode.
- [ ] Include the value in the Plan Digest.
- [ ] Cover a region that removes a line and a region that removes none.
- [ ] Cover truncation and its reported remainder.
- [ ] Cover strict-codec round-trip and unchanged non-managed region proof.

## Acceptance Criteria

- [ ] A managed-refresh plan over a region whose rendering drops a specific line
      reports exactly that line as removed.
- [ ] A managed-refresh plan over a region whose rendering only adds lines reports
      the region with an empty removed-line set.
- [ ] A region whose removed lines exceed the bound reports the bounded list and a
      non-zero omitted count equal to the difference.
- [ ] Two plans over the same repository that differ only in a removed line have
      different Plan Digests.
- [ ] A greenfield plan and a preservation plan over the same fixtures produce
      byte-identical JSON to the recorded goldens.
- [ ] Encoding and decoding a managed-refresh plan document preserves every
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
