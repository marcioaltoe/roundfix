---
task: task_03
spec: 0363-a-reused-owner-pid-means-the-owner-is-gone
status: pending
type: backend
complexity: medium
---

# Task 03: The lock reclaim, the Run slots and deliver start treat a reused owner PID as gone

## Overview

Three readers ask only whether a recorded owner PID is alive, so a dead
owner whose PID the system reused stays present: the orphan reclaim of an
Active Run's lock, the Delivery Queue's Run-slot holders and the release of
a dead queue owner in `roundfix deliver start`. This Task makes each ask
`store.OwnerGone` instead and names the reuse in the reclaim reason. Every
signal path and every reader that already reads a mismatch as gone stays as
it is (`_techspec.md` → System Architecture).

## Requirements

1. MUST make `reclaimOrphanedActiveRun` (`internal/cli/cli.go`) reclaim
   when `store.OwnerGone(ctx, pid, active.OwnerIdentity)` reports gone, and
   change `orphanedActiveRunReason` to take `reused` and return the two
   reasons of API Contract 3 exactly.
2. MUST make `Store.ReclaimOrphanedRun` (`internal/store/store.go`) return
   without a write unless `OwnerGone(ctx, *run.OwnerPID, run.OwnerIdentity)`
   reports gone, with the default reason of API Contract 3.
3. MUST make `RunSlots` (`internal/cli/deliver_workflow.go`) exclude a
   holder whose recorded owner `OwnerGone` reports gone, passing the Run's
   recorded owner identity; no other function of that file changes.
4. MUST move the owner release that `runDeliverStart`
   (`internal/cli/deliver.go`) performs before
   `CreateDeliveryQueueWithLimits` into
   `releaseGoneDeliveryQueueOwner(ctx, runStore, gitRoot) error` (API
   Contract 4), deciding with `OwnerGone(ctx, existing.OwnerPID, existing.OwnerIdentity)`
   and keeping today's error texts; `runDeliverStart` calls it.
5. MUST add `internal/store/reclaim_reused_owner_unix_test.go`
   (`//go:build unix`, suite header comment) with
   `TestReclaimOrphanedRunReclaimsAnOwnerWhosePIDWasReused`: an Active Run
   recorded with the test's live parent PID (`os.Getppid()`) and the
   identity `runtime.GOOS + ":0.0"` is reclaimed with reason
   `owner process <ppid> gone (pid reused); lock reclaimed`; the same Run
   recorded with the parent's real identity from `OwnerProcessIdentity` is
   left Active. It skips, naming the reason, when `OwnerProcessIdentity`
   returns `ErrOwnerProcessUnsupported`.
6. MUST add `internal/cli/owner_pid_reused_unix_test.go` (`//go:build unix`,
   suite header comment) with, each calling `t.Parallel()` first and using
   the same parent-PID technique and temporary Roundfix homes:
   - `TestImplementReclaimsAnActiveRunWhoseOwnerPIDWasReused`: through
     `reclaimOrphanedActiveRun`, the reclaim and its stderr reason with the
     mismatching identity, no reclaim with the matching one;
   - `TestRunSlotsExcludeAHolderWhoseOwnerPIDWasReused`: the holder is
     excluded with the mismatching identity and counted with the matching
     one;
   - `TestDeliverStartReleasesAQueueOwnerWhosePIDWasReused`: through
     `releaseGoneDeliveryQueueOwner`, the owner record is released with the
     mismatching identity and kept with the matching one.
   No test signals any process.
7. MUST update every existing test that calls `orphanedActiveRunReason` or
   pins the absent-owner reason to the new signature, keeping the absent
   reason's bytes; the tests of `orphan_unix_test.go`,
   `reclaim_unix_test.go`, `TestRunSlotsExcludeAHolderWhoseOwnerIsGone`, the
   Force Stop, `deliver stop`, retry, resume and reconcile tests keep
   passing.
8. MUST NOT change `forceStopRun`, `runDeliverStop`, any terminate or signal
   function, `inspectReconcileProcesses`, `deliveryOwnerLine`,
   `runDeliverRetry`, `runDeliverResume`, `ProcessAlive`, or any file
   outside this Task's declarations.

## Subtasks

- [ ] Switch the lock reclaim (CLI and store) to `OwnerGone` and name the reuse.
- [ ] Switch `RunSlots` to `OwnerGone`.
- [ ] Extract and switch the queue owner release of `deliver start`.
- [ ] Add the store and CLI tests.

## Acceptance Criteria

- [ ] A dead owner whose PID was reused frees the Active Run lock with the reuse reason, leaves the Run slots and lets `deliver start` take over the queue.
- [ ] A matching or unprovable identity keeps the owner present everywhere, and no path sends a signal.

## Context

- instruction: `docs/adr/0373-a-proven-owner-identity-mismatch-means-the-owner-is-gone-for-every-liveness-reader.md`
- instruction: `docs/adr/0044-orphaned-run-locks-are-reclaimed-on-proven-owner-death.md`
- instruction: `internal/store/process.go`
- instruction: `internal/store/reclaim_unix_test.go`
- instruction: `internal/cli/orphan_unix_test.go`
- instruction: `internal/cli/deliver_run_slot_wait_unix_test.go`
- creates: `internal/store/reclaim_reused_owner_unix_test.go`
- creates: `internal/cli/owner_pid_reused_unix_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/store/store.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver.go`

## Consumers

- The Active Run reclaim, `RunSlots` and the queue owner release — consumers: `./internal/store`, `./internal/cli`

## Verification

- `grep -q 'func TestReclaimOrphanedRunReclaimsAnOwnerWhosePIDWasReused(' internal/store/reclaim_reused_owner_unix_test.go || exit 1; go test -timeout 20m -count=1 -run '^(TestReclaimOrphanedRunReclaimsAnOwnerWhosePIDWasReused|TestReclaimOrphanedRunCompletesFailedReleasesLockAndJournalsReason|TestReclaimOrphanedRunConcurrentAttemptsJournalOnce)$' ./internal/store` — expected: exit 0. Before this Task the test file does not exist, so the first check fails; afterwards the store reclaims a reused owner and the existing reclaim tests pass.
- `for t in TestImplementReclaimsAnActiveRunWhoseOwnerPIDWasReused TestRunSlotsExcludeAHolderWhoseOwnerPIDWasReused TestDeliverStartReleasesAQueueOwnerWhosePIDWasReused; do grep -q "func $t(" internal/cli/owner_pid_reused_unix_test.go || { printf 'missing %s\n' "$t" >&2; exit 1; }; done; go test -timeout 20m -count=1 -run '^(TestImplementReclaimsAnActiveRunWhoseOwnerPIDWasReused|TestRunSlotsExcludeAHolderWhoseOwnerPIDWasReused|TestDeliverStartReleasesAQueueOwnerWhosePIDWasReused|TestRunImplementReclaimsDeadOwnerActiveRun|TestRunSettleReclaimsDeadOwnerActiveRun|TestReviewFetchReclaimsDeadOwnerActiveRun|TestRunSlotsExcludeAHolderWhoseOwnerIsGone|TestCancelReclaimsAnOrphanedRun)$' ./internal/cli` — expected: exit 0. Before this Task the test file does not exist, so the first check fails; afterwards the three readers treat a reused owner as gone and the existing reclaim and slot tests pass.

## References

- `_prd.md` → Goals; User Story 2; User Story 3; Core Feature 3; Success Metric 4; Success Metric 5
- `_techspec.md` → API Contract 3; API Contract 4; Interfaces; Testing Approach; Build Order 3
- ADR-0373, ADR-0044, ADR-0279
