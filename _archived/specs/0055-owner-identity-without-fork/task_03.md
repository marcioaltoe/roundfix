---
task: task_03
spec: 0055-owner-identity-without-fork
status: completed
type: backend
complexity: medium
---

# Task 03: Mark and warn about a Run created without reuse protection

## Overview

When identity capture fails at Run creation the Run records no identity and
starts anyway, silently, with PID-only protection. Make that state observable:
one startup warning and a durable marker on the Run. The silent NULL degradation
survives only for legacy rows that predate the identity column.

## Requirements

1. MUST add one additive Run column, defaulting to unset, recording that
   identity capture failed at creation.
2. MUST set the marker when capture fails and leave it unset when capture
   succeeds.
3. MUST emit exactly one warning at Run start when the marker is set, naming
   that the Run has PID-only reuse protection.
4. MUST render the marker in Run inspection output so the state is queryable
   after the fact.
5. MUST leave a legacy row — NULL identity with the marker unset — degrading
   exactly as it does today, per ADR-0044.
6. MUST NOT fail Run creation because capture failed.
7. MUST NOT change the compare-and-set terminal completion contract
   (ADR-0052).

## Subtasks

- [ ] Add the column and set it at creation on capture failure.
- [ ] Emit the single startup warning.
- [ ] Render the marker in Run inspection.
- [ ] Cover set, unset, and legacy-NULL rows.

## Acceptance Criteria

- [ ] A Run created with a failing capture carries the marker and prints one
      warning.
- [ ] A Run created with a successful capture carries no marker and prints no
      warning.
- [ ] A legacy row with NULL identity and no marker behaves as it does today.
- [ ] The warning appears once, not once per read of the Run.
- [ ] Run creation succeeds in every case above.

## Context

- interface: `internal/store/store.go`
- interface: `internal/store/store_test.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/cli_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/store/ ./internal/cli/` — expected: pass,
  including the legacy-row case.
- `make verify` — expected: exit 0.

## References

`_prd.md` → Goal 3, Story 3, Feature 3; `_techspec.md` → Build Order 3,
Data Models.

## Result

### Implementation

- Run Database schema v12 adds `owner_identity_unproven INTEGER NOT NULL
  DEFAULT 0`. New Runs set it only when they record an owner PID but owner
  identity capture produced no token; every Run query returns the durable
  marker.
- Implement startup prints one stderr warning from the newly created Run when
  the marker is set. Run reads have no warning side effect.
- `roundfix runs list` appends `owner_identity_unproven=true` to marked Runs,
  leaving unmarked output unchanged.
- Schema upgrades add the marker with its default unset, so legacy NULL
  identities remain distinguishable. Owner proof and compare-and-set terminal
  completion code were not changed.

### Focused checks

- Red signal: the selected store and CLI tests initially failed to compile on
  the missing `Run.OwnerIdentityUnproven` field and Implement capture seam.
- `rtk gofmt -w internal/store/store.go internal/store/store_test.go
  internal/cli/cli.go internal/cli/implement.go internal/cli/implement_test.go
  internal/cli/runs.go internal/cli/cli_test.go` — exit 0.
- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260731T101603Z_bdd82d74273de3d0/.gocache go test -count=1 ./internal/store -run 'TestCreateRun(PersistsOwnerIdentityAcrossRunQueries|WithoutOwnerIdentityMarksCaptureFailure)|TestOpenMigrates(V3RunDatabasePreservingRunsAndRekeyingLocks|V7RunDatabaseAddingOwnerPID|V11RunDatabaseAddingOwnerIdentityUnproven)|TestCompleteRun(WinnerAndIdenticalReplay|ConcurrentTerminalOutcomesHaveOneWinner)'`
  — exit 0, `ok roundfix/internal/store`.
- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260731T101603Z_bdd82d74273de3d0/.gocache go test -count=1 ./internal/cli -run 'TestRunImplement(ExecutesSpecEndToEnd|WarnsOnceAndMarksFailedOwnerIdentityCapture)|TestRunRunsListRendersOwnerIdentityUnprovenMarker|TestBranchIntegrityPreflightMigratesOutdatedRunDatabase|TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner'`
  — exit 0, `ok roundfix/internal/cli`.
- `rtk git diff --check` — exit 0.
- The commands under `## Verification` were not run; the Daemon owns them.

### Acceptance evidence

- Failing capture: `TestCreateRunWithoutOwnerIdentityMarksCaptureFailure` and
  `TestRunImplementWarnsOnceAndMarksFailedOwnerIdentityCapture` prove creation
  succeeds, persists the marker, and emits one PID-only warning.
- Successful capture: `TestCreateRunPersistsOwnerIdentityAcrossRunQueries` and
  `TestRunImplementExecutesSpecEndToEnd` prove the marker stays unset and no
  PID-only warning is printed.
- Legacy NULL: `TestOpenMigratesV7RunDatabaseAddingOwnerPID` and
  `TestOpenMigratesV11RunDatabaseAddingOwnerIdentityUnproven` prove migration
  leaves the marker unset; `TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner`
  preserves ADR-0044 PID-only behavior.
- Warning cardinality and read behavior: the failing-capture Implement test
  counts one warning, reads the stored Run twice, and still observes one.
- Terminal completion: the focused winner/replay and concurrent-winner tests
  preserve ADR-0052 compare-and-set behavior.

### Verification feedback attempt 1

- Inspected the Daemon diagnostic artifact at
  `/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260731T101603Z_bdd82d74273de3d0/verification/batch-003-attempt-1.log`.
  The Store package passed; the CLI failure was isolated to the retained
  Worktree list test's pre-marker stdout expectation.
- Root cause: that test creates a Run with an owner PID and no owner identity,
  so the Task 03 contract correctly marks it and renders
  `owner_identity_unproven=true`. Updated only the expected public output; the
  retained-Worktree guidance remains stderr-only.
- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260731T101603Z_bdd82d74273de3d0/.gocache go test -count=1 ./internal/cli -run '^TestRunRunsListActiveReportsRetainedWorktreesWithoutChangingStdout$'`
  — exit 0, `ok roundfix/internal/cli`.
- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260731T101603Z_bdd82d74273de3d0/.gocache go test -count=1 ./internal/cli -run '^TestRunRunsList(RendersOwnerIdentityUnprovenMarker|ActiveReportsRetainedWorktreesWithoutChangingStdout)$'`
  — exit 0, `ok roundfix/internal/cli`.
- `rtk git diff --check` — exit 0 after the feedback repair.
- The commands under `## Verification` were not rerun; the Daemon owns the
  configured retry.
