---
task: task_02
spec: 0300-reconcile-releases-a-stale-run-owner-record
status: completed
type: backend
complexity: medium
---

# Task 02: Reconcile releases a Run owner record whose owner identity mismatch is proven

## Overview

The owner proof returns the same error for a live process whose start
identity differs from the recorded one as for PID 0 and the inspecting
process's own PID, so `roundfix reconcile` preserves forever a terminal Run
whose owner PID the system gave to another process. This Task gives the proven
mismatch its own error, reports the record as a `runOwnerRecord` process
candidate, and makes `--apply` release it from the Run Database without
sending any signal, while every identity that cannot be read or compared stays
preserved.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-21,
   [`references/2026-10-21-reconcile-preserves-forever-a-process-candidate-whose-owner-pid-now-belongs-to-another-process.md`](references/2026-10-21-reconcile-preserves-forever-a-process-candidate-whose-owner-pid-now-belongs-to-another-process.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1-6 of
   `_techspec.md`.
2. MUST add `ErrOwnerProcessIdentityMismatch` to `internal/store/process.go`
   exactly as API Contract 1 states, wrapping
   `ErrOwnerProcessIdentityUnproven`, and make `proveOwner` return it for the
   comparable mismatch only, with the message bytes it renders today. The
   own-PID and non-positive-PID refusals keep the bare unproven error, and
   every unreadable path keeps `ErrOwnerIdentityUnreadable`.
3. MUST add `ReleaseStaleRunOwnerRecord` in the new file
   `internal/store/run_owner_record.go` exactly as API Contract 2 states: one
   write transaction, a refusal for a missing or non-terminal Run, a
   compare-and-set on the Run ID, owner PID and owner identity, one
   `daemon.outcome` Run Event with the payload `event`
   `stale_run_owner_record_release` only when a row changed, and no change to
   the Run's state, `updated_at` or `completed_at`.
4. MUST make `inspectReconcileProcesses` in `internal/cli/reconcile.go` report
   a selected terminal Run whose `InspectTree` error matches
   `ErrOwnerProcessIdentityMismatch` as the `runOwnerRecord` entry of API
   Contract 3, in `processCandidates` and not in `preservedCandidates`; every
   other error and every case decided before `InspectTree` keep today's entry
   and refusal reason.
5. MUST make `applyReconcileProcesses` release each `runOwnerRecord` entry
   through `ReleaseStaleRunOwnerRecord` exactly as API Contract 4 states,
   never passing such an entry to the process controller, and add
   `RunOwnerRecordsReleased` (`runOwnerRecordsReleased`) to the debris
   summary. The text report prints the heading and summary field of API
   Contract 5, with `run-owner-records-released=<n>` right after
   `processes-applied=<n>`.
6. MUST state the hazard in the code's behavior: the release writes the
   maintainer's Run Database, so it clears only the PID and identity this
   invocation inspected, only for a terminal Run, and sends no signal to the
   process that now holds the PID, whose process tree is never enumerated. A
   dry run writes nothing.
7. MUST keep API Contract 6: Force Stop, `TerminateAndWait`,
   `SignalProvenProcess`, `SignalProvenProcessGroup`, the Delivery Queue owner
   checks in `internal/cli/deliver.go`, `internal/cli/doctor.go`,
   `internal/cli/cli.go` and `internal/store/store.go` keep their bytes. No
   existing test, golden or fixture changes.
8. MUST add the tests `_techspec.md` → Testing Approach lists, each top-level
   test calling `t.Parallel()` first, compiling for Windows, using temporary
   homes and repositories, and passing no real-time duration (build the
   controller with `NewOwnerProcessController()` and replace its process-table
   functions in the test):
   - `TestOwnerProcessControllerTellsAProvenIdentityMismatchFromAnUnreadableIdentity`
     in the new file `internal/store/process_identity_mismatch_test.go`, with
     the rows listed there, including the exact mismatch message, the own PID
     and PID 0 as unproven without a mismatch, a read failure,
     `ErrOwnerProcessUnsupported` and a foreign platform token as unreadable
     without a mismatch, and `TerminateAndWait` returning the mismatch with no
     signal;
   - `TestInspectTreeReportsAProvenMismatchOfALiveProcessWithoutSignallingIt`
     in the new file `internal/store/process_identity_mismatch_unix_test.go`
     (`//go:build unix`), using `startOwnerProcessHelper`;
   - `TestReleaseStaleRunOwnerRecordClearsOnlyTheRecordedTerminalOwner` in the
     new file `internal/store/run_owner_record_test.go`;
   - `TestRunReconcileReleasesAStaleRunOwnerRecordWithoutASignal` and
     `TestRunReconcilePreservesARunOwnerRecordWhoseIdentityIsNotProvenStale`
     in the new file `internal/cli/reconcile_run_owner_record_test.go`, with a
     fake process controller that answers each owner PID with a fixed proof
     error and records every termination request; the second test has one
     subtest each for an unreadable read failure, an uncomparable token, a
     bare unproven error, `ErrOwnerProcessUnsupported` and a Run recorded
     without an owner identity.

## Subtasks

- [ ] Give the proven mismatch its own error with today's message.
- [ ] Add the compare-and-set release of a terminal Run's owner record with its journal event.
- [ ] Report the `runOwnerRecord` candidate and release it on `--apply` without the process controller.
- [ ] Add the store, process-table and command tests.

## Acceptance Criteria

- [ ] A proven mismatch is reported as one `runOwnerRecord` candidate and released by `--apply` with no termination request, and a second `--apply` reports nothing.
- [ ] Unreadable, uncomparable, unsupported, bare unproven and identity-less records stay preserved with today's reasons.
- [ ] The mismatch keeps today's message, still matches the unproven error, and Force Stop still refuses it.
- [ ] Every existing reconcile, Force Stop and residue test named in the Verification passes unchanged.

## Consumers

- `ErrOwnerProcessIdentityMismatch` and `ReleaseStaleRunOwnerRecord` — consumers: `./internal/cli`
- The `runOwnerRecord` candidate and the `runOwnerRecordsReleased` count of the reconcile report — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0313-a-proven-owner-identity-mismatch-releases-a-terminal-runs-owner-record-without-a-signal.md`
- instruction: `docs/user-guide/commands/reconcile.md`
- instruction: `internal/store/store.go`
- instruction: `internal/store/delivery.go`
- instruction: `internal/store/journal.go`
- instruction: `internal/store/process_signal_unix.go`
- instruction: `internal/store/process_unix_test.go`
- instruction: `internal/store/store_test.go`
- instruction: `internal/cli/cli.go`
- instruction: `internal/cli/doctor.go`
- instruction: `internal/cli/deliver.go`
- instruction: `internal/cli/cli_test.go`
- instruction: `internal/cli/doctor_test.go`
- instruction: `internal/cli/reconcile_item_branch_test.go`
- instruction: `internal/cli/reconcile_staging_test.go`
- creates: `internal/store/run_owner_record.go`
- creates: `internal/store/run_owner_record_test.go`
- creates: `internal/store/process_identity_mismatch_test.go`
- creates: `internal/store/process_identity_mismatch_unix_test.go`
- creates: `internal/cli/reconcile_run_owner_record_test.go`
- interface: `internal/store/process.go`
- interface: `internal/cli/reconcile.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestOwnerProcessControllerTellsAProvenIdentityMismatchFromAnUnreadableIdentity|TestInspectTreeReportsAProvenMismatchOfALiveProcessWithoutSignallingIt|TestReleaseStaleRunOwnerRecordClearsOnlyTheRecordedTerminalOwner|TestOwnerProcessControllerProveOwnerClassifiesIdentityEvidence|TestOwnerProcessControllerRefusesMismatchedOwnerIdentity|TestOwnerProcessControllerProveOwnerRefusesMismatchedIdentity|TestRunReconcileReleasesAStaleRunOwnerRecordWithoutASignal|TestRunReconcilePreservesARunOwnerRecordWhoseIdentityIsNotProvenStale|TestRunReconcileOffersAndAppliesOwnedProcessTreesAndRunBranches|TestRunReconcileApplyPreservesProcessWhenTerminationIsUnproven|TestReconcileListsAndReleasesAnItemBranchOfAMergedSpec|TestRunForceStopOwnerIdentityUnreadableFlagRequiresUnreadableProof|TestDoctorReportsProcessResidueUnreadable)$' ./internal/store ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestOwnerProcessControllerTellsAProvenIdentityMismatchFromAnUnreadableIdentity TestInspectTreeReportsAProvenMismatchOfALiveProcessWithoutSignallingIt TestReleaseStaleRunOwnerRecordClearsOnlyTheRecordedTerminalOwner TestOwnerProcessControllerProveOwnerClassifiesIdentityEvidence TestOwnerProcessControllerRefusesMismatchedOwnerIdentity TestOwnerProcessControllerProveOwnerRefusesMismatchedIdentity TestRunReconcileReleasesAStaleRunOwnerRecordWithoutASignal TestRunReconcilePreservesARunOwnerRecordWhoseIdentityIsNotProvenStale TestRunReconcileOffersAndAppliesOwnedProcessTreesAndRunBranches TestRunReconcileApplyPreservesProcessWhenTerminationIsUnproven TestReconcileListsAndReleasesAnItemBranchOfAMergedSpec TestRunForceStopOwnerIdentityUnreadableFlagRequiresUnreadableProof TestDoctorReportsProcessResidueUnreadable; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/store ./internal/cli` — expected: exit 0. Before this Task the five new tests and `ErrOwnerProcessIdentityMismatch` do not exist, so the command fails at its first pass check. Afterwards the mismatch, live-process, release and command tests pass with the unchanged owner proof, Force Stop, reconcile and residue tests, and both packages' tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6
- `_techspec.md` → Current behavior; API Contracts 1-6; Invariants 1-6; Testing Approach; Build Order 2
- ADR-0313


## Result

Implemented this Task's owner-record release path. The comparable identity
mismatch has its own sentinel wrapping the existing unproven error with the
same diagnostic bytes. Reconcile reports it as `runOwnerRecord`, with an
empty process list, and applies it through a terminal-state-guarded
compare-and-set in the Run Database. The release and its single
`daemon.outcome` journal event share one write transaction. The release never
enumerates the unrelated process's tree or requests termination. The report
includes the specified text heading, actions and release count.

Acceptance evidence from focused checks:

- Proven mismatch and idempotence: the new command test asserts exact text
  and JSON fields, unchanged database bytes through both dry runs, one
  signal-free release, unchanged Run outcome and timestamps, one journal
  event, no candidate or inspection on a second apply, and no doctor lineage
  read afterwards. Additional cases cover an already-released record, a
  missing Run and a database-open failure with the required actions and
  failure counts.
- Preserved records: the new command test exercises unreadable read failure,
  uncomparable token, bare unproven, unsupported and identity-less records in
  dry-run and apply modes. Each retains its exact refusal reason and database
  bytes and requests no termination; an identity-less record is never
  inspected.
- Proof compatibility: the new process-table test asserts the exact mismatch
  message, matching of both mismatch and unproven sentinels, refusal by
  `TerminateAndWait` without a signal, own PID and PID zero as bare unproven,
  read failure/unsupported read/foreign token as unreadable, and absence on
  the second liveness read as success. The Unix test inspects a real helper
  and observes it still alive after the mismatch refusal.
- Guarded persistence: the new store test checks required inputs, missing
  and Active Run refusals, no mutation or event for a different PID or
  identity, trimming, clearing only the owner fields, unchanged state and
  timestamps, exact journal source/kind/summary/payload, idempotence, and
  rollback of the owner clear when journal insertion fails.
- Existing behavior: focused owner-proof, Force Stop, residue, process and
  branch reconciliation, item-branch release and staging tests passed
  unchanged. Byte comparison against HEAD confirms that
  `internal/store/store.go`, `internal/store/delivery.go`,
  `internal/store/process_signal_unix.go`, `internal/cli/cli.go`,
  `internal/cli/doctor.go`, `internal/cli/deliver.go`, `TerminateAndWait` and
  every existing test retain their bytes. The complete declared Verification
  remains pending with the Daemon.

Focused checks executed:

1. Before implementation,
   `go test ./internal/store -run '^TestOwnerProcessControllerTellsAProvenIdentityMismatchFromAnUnreadableIdentity$' -count=1`
   exited 1 because `ErrOwnerProcessIdentityMismatch` was undefined.
2. After implementation,
   `go test ./internal/store ./internal/cli -run '^(TestOwnerProcessControllerTellsAProvenIdentityMismatchFromAnUnreadableIdentity|TestInspectTreeReportsAProvenMismatchOfALiveProcessWithoutSignallingIt|TestReleaseStaleRunOwnerRecordClearsOnlyTheRecordedTerminalOwner|TestRunReconcileReleasesAStaleRunOwnerRecordWithoutASignal|TestRunReconcilePreservesARunOwnerRecordWhoseIdentityIsNotProvenStale)$' -count=1 -v`
   exited 0; all five named tests and their subtests passed.
3. `go test ./internal/store -run '^TestOwnerProcessController.*(Identity|ProveOwner)' -count=1`
   exited 0.
4. `go test ./internal/cli -run '^(TestRunForceStopOwnerIdentityUnreadable|TestRunReconcile.*(Process|Branch)|TestDoctorReportsProcessResidue|TestReconcile.*Staging|TestReconcileListsAndReleasesAnItemBranch)' -count=1`
   exited 0.
5. `GOOS=windows go test -c -o /tmp/roundfix-0300-task02-store.test.exe ./internal/store`
   and
   `GOOS=windows go test -c -o /tmp/roundfix-0300-task02-cli.test.exe ./internal/cli`
   both exited 0. These compile the Windows tests; they do not execute them.
6. `git diff --check` exited 0. All five new top-level tests begin with
   `t.Parallel()` and introduce no real-time durations.

Execution environment: the sandbox refused all commands because the assigned
writable root contains a symlink. Approved execution outside the sandbox
resolved this same Run checkout to `/Volumes/Brain-SSD/cache/.roundfix/worktrees/`.
Tests used disposable homes and repositories; no maintainer Run Database was
read or changed. The task file already carried the Daemon's `in_progress`
change at the initial Git-status check; its status was left untouched.

No declared Verification command, repository Verification, commit, push or
pull-request operation was run. No Task Graph or other Task file was edited.

## Carry-forward provenance

- Source Run: `run_20261023T022528Z_d6dffdb8e615fb09`
- Source commit: `3bbf7e3ffeeca697fcdbbcb6b08b7239d1d172bd`
