---
task: task_02
spec: 0363-a-reused-owner-pid-means-the-owner-is-gone
status: pending
type: backend
complexity: medium
---

# Task 02: Doctor reads a reused owner PID as an owner that is gone

## Overview

`roundfix doctor` files a proven start-identity mismatch at a terminal Run's
owner PID as a process table it could not read, so `residue: partial`
never leaves a machine where the system reused a dead owner's PID. This Task
adds the store question `OwnerGone` and makes Doctor's `Residue` count a
mismatch as an owner that is gone. It answers the message reported by argus
([references](references/2026-11-06-doctor-residue-stays-partial-when-an-owner-pid-is-reused.md)).

## Requirements

1. MUST add to `internal/store/process.go` the method
   `(*OwnerProcessControl).OwnerGone` and the package function `OwnerGone`
   (default controller) exactly as `_techspec.md` → Interfaces and API
   Contract 1 define them, built on `proveOwner`; they send no signal and
   write nothing.
2. MUST add `internal/store/process_owner_gone_test.go` with
   `TestOwnerGoneTreatsOnlyAProvenMismatchOrAbsenceAsGone`, which calls
   `t.Parallel()` first and drives a controller with fake `processAbsent`
   and `processStartIdentity` (as
   `process_identity_mismatch_test.go` does) through every case of API
   Contract 1: absent, mismatch, matching, unreadable identity read,
   foreign-platform recorded token, empty recorded identity on a live PID,
   `pid <= 0` and the caller's own PID, asserting both results of each.
3. MUST change `Residue` (`internal/cli/doctor.go`) so a `ReadLineage`
   error that `errors.Is` `store.ErrOwnerProcessIdentityMismatch` adds no
   residue and no read error for that Run (API Contract 2). Every other
   branch and text stays as it is.
4. MUST add to `internal/cli/doctor_test.go`, through
   `testDoctorReportsProcessResidueCase`, each calling `t.Parallel()` first:
   - `TestDoctorReportsNoResidueForAReusedOwnerPID`: `ReadLineage` returns
     an error wrapping `store.ErrOwnerProcessIdentityMismatch` (built as
     `fmt.Errorf("owner process 4242 failed step %q: %w", "prove owner process identity", store.ErrOwnerProcessIdentityMismatch)`);
     the line is `residue: ok (no process residue found)` and no line
     contains `partial`;
   - `TestDoctorKeepsPartialForAnUnreadableOwnerBesideAReusedOne`: two
     terminal Runs, one whose read returns that mismatch and one whose read
     returns an error wrapping `store.ErrOwnerIdentityUnreadable`; the only
     residue line is `residue: partial (could not read process table for Run <id of the second Run>: ...)`
     and the first Run's ID appears in no line.
   `TestDoctorReportsProcessResidueUnreadable` and the other residue tests
   keep passing unchanged.
5. MUST NOT change any signal or terminate path, `proveOwner`'s results,
   `InspectTree`, `InspectTreeProcesses`, the Doctor output of any other
   check, or any file outside this Task's declarations.

## Subtasks

- [ ] Add `OwnerGone` and its table test.
- [ ] Skip a proven mismatch in `Residue`.
- [ ] Add the two Doctor tests.

## Acceptance Criteria

- [ ] `OwnerGone` answers gone only for an absent PID or a proven mismatch, and reports the reuse.
- [ ] Doctor prints `residue: ok (no process residue found)` when the only problem was a reused owner PID, and keeps `partial` for every unreadable or unproven case.

## Context

- instruction: `docs/adr/0373-a-proven-owner-identity-mismatch-means-the-owner-is-gone-for-every-liveness-reader.md`
- instruction: `docs/adr/0313-a-proven-owner-identity-mismatch-releases-a-terminal-runs-owner-record-without-a-signal.md`
- instruction: `internal/store/process_identity_mismatch_test.go`
- creates: `internal/store/process_owner_gone_test.go`
- interface: `internal/store/process.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/doctor_test.go`

## Consumers

- `store.OwnerGone` and the Doctor residue result — consumers: `./internal/store`, `./internal/cli`

## Verification

- `grep -q 'func TestOwnerGoneTreatsOnlyAProvenMismatchOrAbsenceAsGone(' internal/store/process_owner_gone_test.go || exit 1; go test -timeout 20m -count=1 -run '^(TestOwnerGoneTreatsOnlyAProvenMismatchOrAbsenceAsGone|TestOwnerProcessControllerTellsAProvenIdentityMismatchFromAnUnreadableIdentity)$' ./internal/store` — expected: exit 0. Before this Task the test file does not exist, so the first check fails; afterwards the owner-gone table and the existing mismatch test pass.
- `for t in TestDoctorReportsNoResidueForAReusedOwnerPID TestDoctorKeepsPartialForAnUnreadableOwnerBesideAReusedOne; do grep -q "func $t(" internal/cli/doctor_test.go || { printf 'missing %s\n' "$t" >&2; exit 1; }; done; go test -timeout 20m -count=1 -run '^(TestDoctorReportsNoResidueForAReusedOwnerPID|TestDoctorKeepsPartialForAnUnreadableOwnerBesideAReusedOne|TestDoctorReportsProcessResidueReported|TestDoctorReportsProcessResidueExcludesLiveRun|TestDoctorReportsProcessResidueEmpty|TestDoctorReportsProcessResidueUnreadable|TestDoctorReportsProcessResidueDoesNotWriteRunRecord)$' ./internal/cli` — expected: exit 0. Before this Task the two tests do not exist, so the first check fails; afterwards a reused owner gives `ok` and every other residue case keeps its line.

## References

- `_prd.md` → Goals; User Story 1; Core Feature 1; Core Feature 2; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → API Contract 1; API Contract 2; Interfaces; Testing Approach; Build Order 2
- ADR-0373, ADR-0313, ADR-0127
