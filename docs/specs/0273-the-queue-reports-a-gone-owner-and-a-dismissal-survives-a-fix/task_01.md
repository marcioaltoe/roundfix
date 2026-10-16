---
task: task_01
spec: 0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
status: pending
type: backend
complexity: medium
---

# Task 01: Deliver status names a recorded Queue Owner that is not running

## Overview

`roundfix deliver status` prints nothing about the Queue Owner, so a queue
whose recorded owner died with a reboot reads as if it were advancing
(`_techspec.md` → Current behavior). This Task adds the `Owner:` line of API
Contract 1, proved with the start-time identity check `deliver resume` and
`deliver retry` already use, and documents it in the `deliver` user guide. It
is verifiable on its own through `roundfix deliver status` on a temporary home
and repository.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-14,
   [deliver status does not say the queue owner is gone](references/2026-10-14-deliver-status-does-not-say-the-queue-owner-is-gone.md),
   except its optional login or launchd resume, which stays out of scope.
2. MUST implement API Contract 1 in `runDeliverStatus` through a helper
   `deliveryOwnerLine`: only when the queue records an owner PID and at least
   one item's stage is neither `merged` nor `parked`; `store.ProcessAlive`
   first, then the command dependency `ownerProcesses.ProveOwner` with the
   recorded identity; `is not running` for an absent PID, also when it exits
   between the probe and the proof; `has a different process identity` for
   `store.ErrOwnerProcessIdentityUnproven`; no line for a proven live owner;
   no line and the stderr diagnostic
   `roundfix: Delivery Queue owner PID <pid> could not be checked: <error>`
   for any other error or a missing controller. The line is
   `Owner: Delivery Queue owner PID <pid> <state>; run roundfix deliver resume`,
   printed after the `Warning:` lines and before the `Park:` lines.
3. MUST keep status read-only: it opens the Run Database with the reader, never
   calls `ReleaseDeliveryQueueOwner` or any other write, and the recorded
   owner PID and identity are unchanged after it runs. Its exit code stays 0
   and every existing line keeps its bytes and order (Invariants 1-3).
4. MUST NOT change the `deliver` help text, `deliver resume`, `deliver retry`,
   `deliver start` or `deliver stop`.
5. MUST add, in the new `internal/cli/deliver_owner_status_test.go`, each test
   calling `t.Parallel()` first: `TestDeliverStatusReproducesSurfaceTranscriptOwnerNotRunning`,
   which asserts Surface Transcript 1 whole (every stdout and stderr line in
   order and exit 0) for one item at `checking` and the recorded PID `1<<30`;
   `TestDeliverStatusNamesAnOwnerWithADifferentIdentity` (the test's own PID
   and a fake controller returning `store.ErrOwnerProcessIdentityUnproven`);
   `TestDeliverStatusPrintsNoOwnerLineForAProvenOwner` (the test's own PID, a
   fake controller returning `nil`);
   `TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner`;
   `TestDeliverStatusPrintsNoOwnerLineWhenNoItemCanAdvance` (PID `1<<30`,
   every item `merged` or `parked`);
   `TestDeliverStatusReportsAnUncheckableOwnerOnStderr` (a fake controller
   returning another error: exit 0, the stderr diagnostic, no `Owner:` line);
   and `TestDeliverStatusLeavesTheOwnerRecord` (the PID and identity read
   back unchanged after status). No test uses a reaped child's PID as a gone
   owner. Every existing `deliver` test MUST pass unchanged.
6. MUST describe in `docs/user-guide/commands/deliver.md`, where it describes
   what `deliver status` prints, the `Owner:` line with the exact text
   `Owner: Delivery Queue owner PID <pid> is not running; run roundfix deliver resume`,
   the different-identity state, when no line is printed, and that status
   reclaims nothing; then re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Add `deliveryOwnerLine` and call it from `runDeliverStatus`.
- [ ] Add the seven tests.
- [ ] Describe the `Owner:` line in the `deliver` guide.
- [ ] Re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A recorded owner whose PID is absent, or alive under another identity, prints one `Owner:` line naming `roundfix deliver resume`.
- [ ] A proven owner, no recorded owner, or no advanceable item prints no `Owner:` line, and an unreadable identity writes only a stderr diagnostic.
- [ ] Status leaves the owner record as it found it and keeps its exit code and other lines.
- [ ] The `deliver` guide names the line and its states.

## Context

- instruction: `docs/adr/0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md`
- instruction: `docs/adr/0044-orphaned-run-locks-are-reclaimed-on-proven-owner-death.md`
- instruction: `internal/store/process.go`
- instruction: `internal/store/delivery.go`
- interface: `internal/cli/deliver.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/references/behavior-surfaces.json`
- creates: `internal/cli/deliver_owner_status_test.go`

## Consumers

- `deliver status` standard output and standard error — consumers: `./internal/cli`
- Behavior Surface Record of the `deliver` guide — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDeliverStatusReproducesSurfaceTranscriptOwnerNotRunning|TestDeliverStatusNamesAnOwnerWithADifferentIdentity|TestDeliverStatusPrintsNoOwnerLineForAProvenOwner|TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner|TestDeliverStatusPrintsNoOwnerLineWhenNoItemCanAdvance|TestDeliverStatusReportsAnUncheckableOwnerOnStderr|TestDeliverStatusLeavesTheOwnerRecord)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliverStatusReproducesSurfaceTranscriptOwnerNotRunning TestDeliverStatusNamesAnOwnerWithADifferentIdentity TestDeliverStatusPrintsNoOwnerLineForAProvenOwner TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner TestDeliverStatusPrintsNoOwnerLineWhenNoItemCanAdvance TestDeliverStatusReportsAnUncheckableOwnerOnStderr TestDeliverStatusLeavesTheOwnerRecord; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Deliver' ./internal/cli || exit 1; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'Owner: Delivery Queue owner PID <pid> is not running; run roundfix deliver resume' || { printf 'deliver guide lacks the Owner line\n' >&2; exit 1; }; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task the seven tests do not exist and the guide lacks the line, so the command fails; afterwards they and every existing `Deliver` test pass and the Behavior Surface Record is current.

## References

- `_prd.md` → Goals; User Story 1; Core Feature 1; Core Feature 4; Success Metric 1; User Experience
- `_techspec.md` → API Contract 1; API Contract 4; Surface Transcript 1; Invariants 1, 2 and 3; Current behavior; Testing Approach; Build Order 1
- ADR-0279; ADR-0044; ADR-0266
