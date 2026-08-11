---
task: task_04
spec: 0084-an-update-that-can-run
status: completed
type: test
complexity: medium
---

# Task 04: Prove the update converges on a second run

## Overview

Adds the property that tells a finished sweep from an unfinished one: after an
applied refresh, running the update again against an unchanged catalog reports
the repository current and proposes nothing. The property has no test today, and
it is the only check that fails if Setup Manifest republication ever regresses —
which is precisely the regression that made the cold-start block permanent in the
first place.

## Requirements

1. MUST build an adopted repository copy on the filesystem, age its Setup Manifest
   so the recorded digests no longer describe its managed regions, and leave the
   managed regions themselves untouched.
2. MUST apply a managed refresh to that copy and assert the apply verifies.
3. MUST assert the second plan over the same copy and the same catalog reports the
   repository current and proposes zero file changes.
4. MUST assert the republished Setup Manifest's recorded digests describe the
   bytes on disk after the apply, so convergence is proven by the record and not
   only by the reported state.
5. MUST assert that every byte outside a managed marker is identical before and
   after the apply, so convergence is not bought by rewriting authored prose.
6. MUST fail if manifest republication is removed: the test must distinguish a
   converged repository from one whose second run proposes the same work again.

## Subtasks

- [x] Build the aged-manifest fixture from an adopted copy.
- [x] Apply the refresh and assert verification.
- [x] Assert the second run reports current with zero file changes.
- [x] Assert the republished manifest describes the on-disk bytes.
- [x] Assert non-managed regions are byte-identical across the apply.
- [x] Demonstrate the negative: a copy whose manifest is not republished does not
      report current.

## Rehearsal Cases

- Case: an adopted copy whose Setup Manifest records digests that no longer
  describe its untouched managed regions; Observation: the first run reaches a
  ready plan rather than an action-required state.
- Case: the same copy after an approved apply; Observation: the second run
  reports the repository current and proposes zero file changes.
- Case: the same copy after an approved apply, reading the Setup Manifest
  directly; Observation: every recorded managed-artifact digest equals the digest
  of the corresponding on-disk region.
- Case: a copy whose Setup Manifest is deliberately left un-republished after the
  refresh; Observation: the second run proposes the same changes again, proving
  the assertion can fail.
- Case: the same copy, comparing every region outside a managed marker before and
  after the apply; Observation: every region digest is unchanged.

## Acceptance Criteria

- [x] The aged-manifest copy reaches a ready plan on the first run.
- [x] Applying that plan verifies.
- [x] The second run over the applied copy reports the repository current with
      zero proposed file changes.
- [x] Every managed-artifact digest recorded in the applied copy's Setup Manifest
      equals the digest of the corresponding on-disk region.
- [x] Every region outside a managed marker has an identical digest before and
      after the apply.
- [x] The negative case, with republication suppressed, does not report current.

## Context

- interface: `internal/baseline/preservation.go`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/apply.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/baseline/ -run 'Converge' -v > /tmp/0084-task-04-a.log 2>&1 && grep -q '^--- PASS: .*Converge' /tmp/0084-task-04-a.log` — expected: exits 0, proving the convergence cases exist and pass rather than being selected out.
- `go test ./internal/baseline/ -run 'Converge' -count=2 > /tmp/0084-task-04-b.log 2>&1` — expected: exits 0, proving the fixture does not depend on state left by a previous run.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0.

## References

- `_techspec.md` → Build Order 4; Testing Approach.
- `_prd.md` → Core Feature 4; User Story 4; Goal 2; Success Metrics.
- `references/2026-08-08-the-update-refuses-six-of-the-eight-copies-it-exists-to-update.md`
  → the measured second-run behavior this task turns into a test.
- ADR-0103, ADR-0100.

## Result

### Implementation

- Added filesystem-backed Managed Refresh convergence coverage to the existing
  plan/apply suite. The fixture adopts a repository, changes every recorded
  managed-artifact digest in its Setup Manifest without changing any carrier,
  then plans and applies the refresh with one catalog instance.
- The positive case reads the republished Setup Manifest from disk, independently
  hashes every corresponding on-disk Managed Region, compares exact digests for
  every region outside managed markers, and requires the second plan to contain
  zero file changes.
- The negative companion restores the aged Setup Manifest after apply to model
  suppressed republication, then requires the second plan to repeat the first
  file-change set.

### Focused checks

- Pre-change signal:
  `rtk rg -n 'ManagedRefresh.*Converge|Converge.*ManagedRefresh' internal/baseline/*_test.go`
  exited 1 with no matches, confirming the Managed Refresh convergence property
  had no owning package test.
- After implementation:
  `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260808T153649Z_78746d4b80d08fc7.task_04/.gocache go test ./internal/baseline -run '^(TestManagedRefreshConvergesAfterManifestRepublication|TestManagedRefreshDoesNotConvergeWithoutManifestRepublication|TestManagedRefreshPreservesNonManagedRegionDigests|TestManagedRefreshPlanReportsEmptyRemovedLines)$' -count=1`
  passed: `ok roundfix/internal/baseline 2.485s`.
- `rtk gofmt -d internal/baseline/plan_test.go` exited 0 with no output.
- `rtk git -c core.fsmonitor=false diff --check` exited 0 with no output.

### Acceptance evidence

1. `newManagedRefreshConvergenceFixture` changes every Setup Manifest
   managed-artifact digest to a stale value, proves all managed-carrier bytes
   stayed identical, and `TestManagedRefreshConvergesAfterManifestRepublication`
   requires the first Managed Refresh outcome to contain a ready plan.
2. The same test applies the plan and requires state `verified` plus verified
   approved postimages in the Result Status Matrix.
3. The positive test rebuilds the plan against the same repository and catalog
   and requires zero `FileChanges`, the planner condition the update command
   reports as `current`.
4. `assertSetupManifestMatchesManagedRegions` reads the applied copy's Setup
   Manifest directly and requires every recorded artifact digest to equal an
   independently computed digest of its one on-disk Managed Region.
5. The positive test records every outside-marker region digest before apply and
   requires an identical path-and-order digest map after apply, including blank
   outside-marker spans.
6. `TestManagedRefreshDoesNotConvergeWithoutManifestRepublication` restores the
   aged manifest after the real apply, rejects a zero-change second plan, and
   requires that plan to propose exactly the first plan's file changes again.

The Daemon-owned commands under `## Verification` were not run in this Agent
turn.
