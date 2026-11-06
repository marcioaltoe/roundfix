---
task: task_03
spec: 0363-a-reused-owner-pid-means-the-owner-is-gone
status: completed
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

## Result

Skills loaded: implement-task, coding-guidelines, golang-cli, golang-testing, testing-boss, no-workarounds, systematic-debugging, agentic-cli-design, golang-error-handling, golang-context, evidence-gate, tech-writer.

Implementation:

- CLI and store orphan reclamation now pass the recorded identity to `OwnerGone`. A proven mismatch uses exactly `owner process <pid> gone (pid reused); lock reclaimed`; absent owners retain `owner process <pid> not running; lock reclaimed`, including the store's default reason.
- `RunSlots` excludes owners proven gone without reclaiming their Run records. No other function in `deliver_workflow.go` changed.
- `runDeliverStart` calls the extracted `releaseGoneDeliveryQueueOwner` before queue creation. The helper uses `OwnerGone`, retains the existing compare-and-release operation, and preserves read/release errors and the owner-changed diagnostic.
- Both new Unix test files have suite headers, use the live parent PID and real start identity, and skip with the unsupported-identity reason only for `ErrOwnerProcessUnsupported`. Tests use disposable homes and never signal the parent or any other process. Existing tests contain no calls to `orphanedActiveRunReason`; their absent-owner reason assertions remain byte-identical.

Acceptance evidence:

- Reused owner: the new store regression checks Failed state, completion time, released lock and the exact journal reason using the default reason argument. The CLI reclaim regression checks the released Spec lock, journal and exact stderr warning. The slot regression checks exclusion while the Run remains unchanged. The queue regression checks that only the owner record is cleared and that a matching owner can then claim it.
- Present or unprovable owner: every new regression checks a matching identity and a foreign-platform identity; store reclaim, CLI reclaim and slots also check an empty identity. Reclaim keeps the record, lock and journal unchanged and prints nothing; slots retain the holder; queue start retains the entire queue. The existing `OwnerGone` unit matrix also covers unreadable/unsupported identities, process-table errors, the caller PID and nonpositive PIDs, with a signal callback that fails the test if invoked.
- Signal and already-correct readers: production diff inspection confirms Force Stop, `deliver stop`, terminate/signal functions, reconcile inspection, owner status, retry/resume and `ProcessAlive` are unchanged. Focused compatibility tests below exercise those existing boundaries.

Focused checks (all final checks exited 0):

- Before the production edit: `rtk proxy go test -timeout 2m -count=1 -run '^TestReclaimOrphanedRunReclaimsAnOwnerWhosePIDWasReused$' ./internal/store` exited 1 as expected: the reused-identity case left the Run Active with its lock present.
- `rtk proxy go test -timeout 3m -count=1 -run '^(TestReclaimOrphanedRunReclaimsAnOwnerWhosePIDWasReused|TestOwnerGoneTreatsOnlyAProvenMismatchOrAbsenceAsGone)$' ./internal/store` passed.
- `rtk proxy go test -timeout 3m -count=1 -run '^(TestImplementReclaimsAnActiveRunWhoseOwnerPIDWasReused|TestRunSlotsExcludeAHolderWhoseOwnerPIDWasReused|TestDeliverStartReleasesAQueueOwnerWhosePIDWasReused)$' ./internal/cli` passed. Its first attempt found a fixture API typo (`ActiveImplementRun`); the fixture now uses the existing `ActiveSpecRun` API.
- `rtk proxy go test -timeout 3m -count=1 -run '^TestReclaimOrphanedRun(CompletesFailedReleasesLockAndJournalsReason|ConcurrentAttemptsJournalOnce)$' ./internal/store` passed, preserving absent-owner reasons and single-event concurrent reclaim.
- `rtk proxy go test -timeout 5m -count=1 -run '^(TestRunForceStopOwnerPIDReuseFailsClosed|TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion|TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner|TestRunSlotsExcludeAHolderWhoseOwnerIsGone|TestDeliverCommandStopsAndResumesPersistedQueue|TestResumeReleasesAStaleOwner|TestDeliverRetryStartsAnOwnerWhenNoneIsRunning|TestDeliverRetryHandsTheItemToALiveOwner|TestDeliverRetryReclaimsAStaleOwnerAndStartsOne|TestRunReconcileReleasesAStaleRunOwnerRecordWithoutASignal|TestRunReconcilePreservesARunOwnerRecordWhoseIdentityIsNotProvenStale|TestRunImplementReclaimsDeadOwnerActiveRun|TestRunSettleReclaimsDeadOwnerActiveRun|TestReviewFetchReclaimsDeadOwnerActiveRun|TestCancelReclaimsAnOrphanedRun)$' ./internal/cli` passed.
- `GOOS=windows rtk proxy go build ./internal/store ./internal/cli` and `GOOS=linux rtk proxy go build ./internal/store ./internal/cli` succeeded.
- `gofmt` applied to the six changed Go files; `rtk proxy git diff --check` passed.

Declared Task Verification and repository Verification were not run; they remain Daemon-owned. The pre-existing `status: in_progress` is unchanged. No Task Graph, other Task, commit, push or pull request was changed or created. No follow-up work identified.
