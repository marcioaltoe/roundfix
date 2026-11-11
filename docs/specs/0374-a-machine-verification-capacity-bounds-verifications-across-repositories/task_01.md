---
task: task_01
spec: 0374-a-machine-verification-capacity-bounds-verifications-across-repositories
status: pending
type: backend
complexity: high
---

# Task 01: The Machine Verification Capacity has its User Config key and its entry store

## Overview

`runs.max_verifying` becomes a User Config key, `1` by default and `0` off,
written into the User Config template. A new store file keeps the capacity's
entries under `verification-capacity/` in the Roundfix Home: an acquisition
takes a ticket, waits in ticket order, removes only entries whose owner is
proven gone, and a reader lists holders and waiters without writing. Nothing
calls the store yet; task_02 and task_03 do.

## Requirements

1. MUST characterize first: on the starting tree a User Config holding
   `runs:\n  max_verifying: 1\n` is refused by `config.Load` as an unknown
   key, and no file under `internal/store` mentions `verification-capacity`.
   Record both in the Task's `## Result`.
2. MUST add the key as `_techspec.md` → API Contract 1 says, in
   `internal/config/config.go`: `Runs.MaxVerifying`, the built-in
   `defaultRunsMaxVerifying = 1`, `runsOverlay.MaxVerifying` with tag
   `yaml:"max_verifying"` applied from User Config only, the Project Config
   removal with `warnIgnoredProjectSetting("runs.max_verifying")`, the
   refusal `runs.max_verifying must be greater than or equal to 0`, and the
   two user-template lines, word for word:
   `  # Maximum Verifications at once across repositories; 0 disables the bound.`
   and `  max_verifying: 1`, right after the `max_active` line. The project
   template stays unchanged.
3. MUST create `internal/store/machine_verification.go` with the types,
   functions, directory layout, entry schema `roundfix/machine-verification-entry/v1`
   and behavior of API Contracts 2 and 3 and Data Models: the lock file
   `lock` taken with the existing `acquireWriteLock`/`releaseWriteLock`
   helpers (no new lock implementation, no build-tagged copy), the ticket
   `sequence`, entries written by temporary file and rename, ticket order,
   removal only when `OwnerGone` answers gone, an unreadable entry skipped and
   reported but never removed or counted, the 2-second cross-process round
   with an in-process wake on every release or removal, removal of its own
   entry on a cancelled wait and on `Release` with a context the caller's
   cancellation does not cancel, bounded by 30 seconds, and the error prefix
   `wait for the Machine Verification Capacity`. With a capacity of `0` or
   less it creates no directory and no file (Invariant 6). The reader takes no
   lock, writes nothing, and leaves out an entry whose owner is proven gone
   (Invariant 5). `DescribeMachineVerificationHolder` returns exactly
   `Run <run id> of <spec> in <repository>` and
   `Delivery Queue item <spec> in <repository>`.
4. MUST NOT remove an entry by its age, a missed poll or a heartbeat, signal
   any process, or add a Run Database table, column or schema version
   (Invariant 3).
5. MUST create `internal/store/machine_verification_test.go` with
   `TestMachineVerificationCapacityServesWaitersInTicketOrder`,
   `TestMachineVerificationCapacityKeepsALaterWaiterBehindALiveWaiterAhead`
   (its input is a `waiting` entry the test writes for its own PID with ticket
   1 and no unit held; a new acquisition with capacity 1 must call `onWait`
   with that entry ahead and nobody holding, and must not be granted),
   `TestMachineVerificationCapacityNeverHoldsMoreThanItsUnits` (capacity 2:
   two grants without `onWait`, the third waits with two holders, and is
   granted after one release), `TestMachineVerificationCapacityOffWritesNothing`,
   `TestACancelledMachineVerificationWaitLeavesNoEntry` and
   `TestReadMachineVerificationLeavesOutAnOwnerProvenGoneAndWritesNothing`
   (the directory's entry names, sizes and modification times are equal
   before and after the read); and
   `internal/store/machine_verification_unix_test.go` (`//go:build unix`)
   with `TestMachineVerificationCapacityRemovesOnlyAnOwnerProvenGone`, the
   table of `_techspec.md` → Testing Approach: rows `reused` (the parent PID
   with identity `<GOOS>:0.0`) and `absent` (a PID from `reapedOwnerPID`)
   are removed and the acquisition is granted at once; rows `matching`,
   `foreign` (`foreign:0.0`) and `empty` keep the entry and the acquisition
   waits. Every top-level test calls `t.Parallel()` first; waits end on the
   observed event or `internal/testwait`, never on a sleep.
6. MUST create `internal/config/machine_verification_config_test.go` with
   `TestRunsMaxVerifyingDefaultsToOne`,
   `TestRunsMaxVerifyingAcceptsZeroAndRefusesANegativeValue`,
   `TestRunsMaxVerifyingInProjectConfigIsIgnoredWithAWarning` (asserting the
   warning line of API Contract 1 and that the User Config value stands) and
   `TestTheUserConfigTemplateWritesRunsMaxVerifying` (both template lines in
   `DefaultConfigYAML`, neither in `DefaultProjectConfigYAML`, and the user
   template loads with `MaxVerifying` `1`). Every existing test of
   `internal/config` passes unchanged.
7. MUST, as the Task's last step, re-record the Behavior Surface Record,
   because the set of configuration keys changed, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   The Task changes no Governed Path.

## Subtasks

- [ ] Characterize the unknown key and the absent store.
- [ ] Add the key, its validation, its Project Config warning and the template lines.
- [ ] Write the entry store and its reader.
- [ ] Add the store and config tests.
- [ ] Re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] `runs.max_verifying` loads from User Config with default `1`, accepts `0`, refuses a negative value and is ignored with a warning in Project Config.
- [ ] Waiters are granted in ticket order, never beyond the capacity, and an entry is removed only when its owner is proven gone.
- [ ] A capacity of `0` writes nothing, and the reader writes nothing.

## Context

- instruction: `docs/adr/0384-a-machine-verification-capacity-bounds-the-verifications-of-every-repository.md`
- instruction: `docs/adr/0373-a-proven-owner-identity-mismatch-means-the-owner-is-gone-for-every-liveness-reader.md`
- instruction: `internal/store/process.go`
- instruction: `internal/store/writelock_unix.go`
- instruction: `internal/store/writelock_windows.go`
- instruction: `internal/store/reclaim_unix_test.go`
- instruction: `internal/store/reclaim_reused_owner_unix_test.go`
- instruction: `internal/config/run_slot_wait_test.go`
- instruction: `internal/testwait/testwait.go`
- interface: `internal/config/config.go`
- creates: `internal/config/machine_verification_config_test.go`
- creates: `internal/store/machine_verification.go`
- creates: `internal/store/machine_verification_test.go`
- creates: `internal/store/machine_verification_unix_test.go`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `store.AcquireMachineVerification`, `store.ReadMachineVerification` and `store.DescribeMachineVerificationHolder` — consumers: `./internal/store`
- `config.Runs.MaxVerifying` and the user template — consumers: `./internal/config`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestMachineVerificationCapacityServesWaitersInTicketOrder|TestMachineVerificationCapacityKeepsALaterWaiterBehindALiveWaiterAhead|TestMachineVerificationCapacityNeverHoldsMoreThanItsUnits|TestMachineVerificationCapacityOffWritesNothing|TestACancelledMachineVerificationWaitLeavesNoEntry|TestReadMachineVerificationLeavesOutAnOwnerProvenGoneAndWritesNothing|TestMachineVerificationCapacityRemovesOnlyAnOwnerProvenGone)$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestMachineVerificationCapacityServesWaitersInTicketOrder TestMachineVerificationCapacityKeepsALaterWaiterBehindALiveWaiterAhead TestMachineVerificationCapacityNeverHoldsMoreThanItsUnits TestMachineVerificationCapacityOffWritesNothing TestACancelledMachineVerificationWaitLeavesNoEntry TestReadMachineVerificationLeavesOutAnOwnerProvenGoneAndWritesNothing TestMachineVerificationCapacityRemovesOnlyAnOwnerProvenGone; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/store && GOOS=windows go vet ./internal/store` — expected: exit 0. Before this Task the seven tests do not exist, so their pass lines are missing; afterwards they pass, and the store, its tests included, compiles for Windows.

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestRunsMaxVerifyingDefaultsToOne|TestRunsMaxVerifyingAcceptsZeroAndRefusesANegativeValue|TestRunsMaxVerifyingInProjectConfigIsIgnoredWithAWarning|TestTheUserConfigTemplateWritesRunsMaxVerifying)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRunsMaxVerifyingDefaultsToOne TestRunsMaxVerifyingAcceptsZeroAndRefusesANegativeValue TestRunsMaxVerifyingInProjectConfigIsIgnoredWithAWarning TestTheUserConfigTemplateWritesRunsMaxVerifying; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/config && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — expected: exit 0. Before this Task the four tests do not exist; afterwards they pass, every existing configuration test passes, and the Behavior Surface Record matches the new key set.

## References

- `_prd.md` → Goal 3; Goal 5; User Story 4; User Story 5; Core Feature 1; Core Feature 4; Core Feature 5; Core Feature 9; Success Metric 1; Success Metric 2
- `_techspec.md` → Current behavior 4 and 6; Interfaces; Data Models; API Contracts 1-3; Invariants 3-6; Testing Approach; Build Order 1
- ADR-0384
