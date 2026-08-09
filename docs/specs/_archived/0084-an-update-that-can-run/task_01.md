---
task: task_01
spec: 0084-an-update-that-can-run
status: completed
type: backend
complexity: high
---

# Task 01: Classify a managed region instead of blocking on it

## Overview

Replaces the managed-refresh marker check that turns any byte difference into a
blocking finding. A managed region whose bytes are not the ones the Setup
Manifest recorded becomes a reported classification carried on the preservation
plan, and planning continues. Only a duplicated managed identity stays blocking,
because a refresh has no defensible target when the same identity appears twice
in one file. This slice is verifiable through planning alone, before any output
surface changes.

## Requirements

1. MUST classify a managed region whose on-disk bytes differ from the digest the
   adopted Setup Manifest recorded as unrecorded, carrying its path, its managed
   identity, and a reason distinguishing a digest mismatch from an absent marker.
2. MUST let root-preservation planning reach a ready state in managed-refresh
   mode when every managed-region difference is classified, so a repository whose
   only obstacle was a stale recorded digest produces a plan.
3. MUST keep a duplicated managed identity within one file blocking, under a
   finding code distinct from the retired one, naming the file and the duplicated
   identity.
4. MUST retire the blocking finding code the previous check emitted, so no caller
   can keep depending on the old behavior by name.
5. MUST carry the classification on the root-preservation plan in a field that is
   absent when nothing is unrecorded.
6. MUST leave the greenfield and preservation modes behaviorally unchanged,
   including their findings, backups, and dispositions.
7. MUST leave retention accounting on this path untouched, so an unaccounted
   managed clause still blocks planning in managed-refresh mode.
8. MUST prove, by test, that an unsafe root carrier still blocks a managed
   refresh, so this change narrows only the marker condition.

## Subtasks

- [x] Introduce the unrecorded-region type and its reason values.
- [x] Replace the blocking marker check with the classifier.
- [x] Keep a duplicated managed identity blocking under a new finding code.
- [x] Carry the classification on the preservation plan.
- [x] Cover a stale-digest repository reaching a ready plan.
- [x] Cover a duplicated managed identity still blocking.
- [x] Cover an unsafe root carrier and an unaccounted clause still blocking.

## Acceptance Criteria

- [x] A managed-refresh plan against a repository whose recorded digest is stale
      reports a ready state and lists the region as unrecorded with reason
      `digest-mismatch`.
- [x] A managed-refresh plan against a repository whose managed marker is absent
      for a recorded artifact reports a ready state and lists the region as
      unrecorded with reason `marker-absent`.
- [x] A repository with the same managed identity twice in one file produces no
      applicable plan and names the file and the identity.
- [x] Searching the package for the retired finding code returns no production
      occurrence.
- [x] A managed-refresh plan against a repository with no unrecorded region omits
      the classification field entirely.
- [x] An unsafe root carrier and an unaccounted managed clause each still block a
      managed refresh.

## Context

- interface: `internal/baseline/preservation.go`
- interface: `internal/baseline/plan.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/baseline/ -run 'UnrecordedManagedRegion' -v > /tmp/0084-task-01-a.log 2>&1 && grep -q '^--- PASS: .*UnrecordedManagedRegion' /tmp/0084-task-01-a.log` — expected: exits 0, proving the classification cases exist and pass rather than being selected out.
- `go test ./internal/baseline/ -run 'AmbiguousManagedMarker' -v > /tmp/0084-task-01-b.log 2>&1 && grep -q '^--- PASS: .*AmbiguousManagedMarker' /tmp/0084-task-01-b.log` — expected: exits 0, proving the surviving blocking condition is covered.
- `grep -rn 'managed-marker.modified' internal/ > /tmp/0084-task-01-c.log 2>&1; grep -v '_test.go' /tmp/0084-task-01-c.log > /tmp/0084-task-01-d.log; test ! -s /tmp/0084-task-01-d.log` — expected: exits 0, proving the retired code has no production occurrence.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0, with the greenfield and preservation corpora passing unchanged.

## References

- `_techspec.md` → Build Order 1; Interfaces: `UnrecordedManagedRegion`,
  `classifyManagedRegions`; Testing Approach.
- `_prd.md` → Core Feature 1; User Story 2; Goals 1 and 3.
- `references/2026-08-08-the-update-refuses-six-of-the-eight-copies-it-exists-to-update.md`
  → the measured cold-start block this task removes.
- ADR-0101, ADR-0102, ADR-0058, ADR-0100.

## Result

### Implementation

- Root-preservation planning now carries typed unrecorded managed-region
  classifications with `digest-mismatch` and `marker-absent` reasons. The
  optional JSON field remains absent when the classification slice is empty.
- Managed-refresh planning reports those classifications without blocking.
  A duplicated managed identity remains blocking under
  `baseline.preservation.managed-marker.ambiguous`, with the carrier path and
  identity in the finding.
- The prior full-plan regression case now proves that an edited managed region
  reaches a ready plan. Greenfield, preservation, unsafe-carrier, and retention
  paths retain their existing behavior.

### Focused checks

- Pre-change signal:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-0084-task01-gocache go test ./internal/baseline -run 'Test(UnrecordedManagedRegion|AmbiguousManagedMarker|ManagedRefreshUnrecordedManagedRegion)' -count=1`
  failed to compile because `RootPreservationPlan.UnrecordedManagedRegions` and
  the unrecorded-region types did not exist.
- After implementation, the same filtered command passed:
  `ok roundfix/internal/baseline 1.019s`.
- Regression check:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-0084-task01-gocache go test ./internal/baseline -run 'Test(GreenfieldPlanBacksUpWithoutImport|PreservationPlanAcceptsCompleteDecisionDocument|ManagedRefreshUnsafeRootCarrierStillBlocks|ManagedRefreshUnaccountedClauseStillBlocks)' -count=1`
  passed: `ok roundfix/internal/baseline 0.826s`.
- `rtk rg -n 'baseline\.preservation\.managed-marker\.modified' internal`
  returned no matches.

### Acceptance evidence

1. `TestUnrecordedManagedRegionDigestMismatchReachesReadyPlan` asserts ready
   root-preservation state plus the exact path, managed identity, and
   `digest-mismatch` reason. The full construction seam
   `TestManagedRefreshUnrecordedManagedRegionReachesReadyPlan` also returns a
   non-nil ready plan.
2. `TestUnrecordedManagedRegionMarkerAbsentReachesReadyPlan` asserts ready state
   plus the exact path, managed identity, and `marker-absent` reason.
3. `TestAmbiguousManagedMarkerBlocksManagedRefresh` asserts blocked state and
   the new ambiguity finding with both the carrier path and duplicated identity.
4. The focused package search found no production occurrence of the retired
   finding code.
5. `TestUnrecordedManagedRegionFieldOmittedWhenCurrent` asserts an empty
   classification and verifies marshaled JSON omits
   `unrecordedManagedRegions`.
6. `TestManagedRefreshUnsafeRootCarrierStillBlocks` and
   `TestManagedRefreshUnaccountedClauseStillBlocks` each passed. The same
   regression command also passed the existing greenfield and preservation
   behavior tests.

The Daemon-owned `## Verification` commands were not run in this Agent turn.
