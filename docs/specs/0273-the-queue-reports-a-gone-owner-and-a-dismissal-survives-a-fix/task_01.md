---
task: task_01
spec: 0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
status: completed
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

## Result

Implemented the Task 01 slice for Daemon Verification and settlement.

- Added `deliveryOwnerLine` to the existing reader-only status path, after
  `Warning:` lines and before `Park:` lines. It checks for a recorded PID and
  an item that can advance, probes liveness before proving the recorded
  identity, and rechecks liveness after a successful proof to cover an owner
  exiting during the proof. Proven absence and identity mismatch name
  `roundfix deliver resume`; other proof errors and a missing controller
  produce only the prescribed stderr diagnostic.
- Added all seven required tests in
  `internal/cli/deliver_owner_status_test.go`, plus an exact transcript test
  for placement among Warning, Park, limits, usage and Pending Question
  lines. Every top-level test and subtest calls `t.Parallel()` first. Gone
  owners use PID `1<<30`; live-owner cases use the test's own PID with the
  existing fake controller. No test signals an owner.
- Updated the status paragraph in the deliver guide and ran its sanctioned
  Behavior Surface Record command. The generator changed only that guide's
  fingerprint in `docs/references/behavior-surfaces.json`.

### Acceptance evidence

| Acceptance criterion | Implementation and focused-check evidence |
| --- | --- |
| Absent or differently identified owner prints one recovery line | `TestDeliverStatusReproducesSurfaceTranscriptOwnerNotRunning` checks the whole Surface Transcript 1, empty stderr and exit 0 for PID `1<<30`; `TestDeliverStatusNamesAnOwnerWithADifferentIdentity` checks the identity-mismatch line and recorded proof arguments. Both ran in the focused status selection. |
| Proven owner, no owner or no advanceable item prints no line; uncheckable identity writes a diagnostic | The three required `TestDeliverStatusPrintsNoOwnerLine...` tests and `TestDeliverStatusReportsAnUncheckableOwnerOnStderr` ran in the focused selection. Diagnostic subtests cover unreadable identity, another error and a missing controller, with exact stdout/stderr and exit 0. |
| Status preserves the owner record, exit code and existing lines | `TestDeliverStatusLeavesTheOwnerRecord` reads the entire queue before and after status and compares it unchanged for absent, mismatched, proven and unreadable owners. Exact transcript tests cover exit 0 and line order; the existing status tests also ran unchanged. Source review confirms `store.OpenReader` remains in use and the helper has no store writer or reclamation call. |
| Guide names the line and its states | The deliver guide contains the exact required line, identity-mismatch wording, suppression conditions, stderr diagnostic, exit 0 and read-only recovery behavior. The record command exited 0 and regenerated its fingerprint. |

### Focused checks

- Before the implementation,
  `go test -count=1 ./internal/cli -run '^TestDeliverStatusReproducesSurfaceTranscriptOwnerNotRunning$'`
  exited 1: stdout had the checking item, limits and usage but lacked the
  Owner line; exit 0 and empty stderr already matched.
- After the implementation,
  `go test -count=1 ./internal/cli -run '^TestDeliverStatus'` exited 0
  (`ok roundfix/internal/cli`). This ran the new owner tests and existing
  status tests together.
- Required regeneration:
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited 0 (`ok roundfix/internal/docscontract`).
- `git diff --check` exited 0. Both Go files were formatted with `gofmt`.

The shell sandbox rejected the configured symlinked writable root before
execution; the approved escalation mechanism permitted the scoped file work
and focused checks in the assigned checkout.

The Task's declared Verification, all-deliver selection and repository gates
remain for the Daemon. Task status and its authored Verification are
unchanged; no other Task or Task Graph was edited, and no commit, push or
pull request was made. The feedback investigation below records a separately
owned fixture blocker.


### Verification Feedback — attempt 1

Inspected the Daemon artifact
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261016T175249Z_782526ccf3bf55ac/verification/batch-001-attempt-1.log`.
The repository gate stopped in the existing
`TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate/stdout`: its probe
ended by SIGKILL instead of the fixture's intended exit 2. No owner-status
test failure was reported.

Source inspection found that `newDeliveryItemFixture` still builds by copying
onto `bin/roundfix`, and the fallback test rebuilds that same path in its
second subtest. ADR-0281 documents this exact intermittent failure and the
required replacement-by-rename repair. The failing test and its production
workflow are unchanged by Task 01, and their execution does not call
`runDeliverStatus` or `deliveryOwnerLine`.

The local report
`~/Library/Logs/DiagnosticReports/roundfix-2026-10-16-150208.ips` records a
`roundfix` child of `cli.test` terminated in namespace `CODESIGNING` with
indicator `Taskgated Invalid Signature`. This supports the failure class in
ADR-0281; the gate log does not include a child PID, so a direct PID join to
this specific crash report is unavailable.

Ran one focused check:
`go test -count=1 ./internal/cli -run '^(TestDeliverStatus.*|TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate)$'`.
It exited 0 (`ok roundfix/internal/cli`). This checks Task 01's status behavior
alongside the reported test; a single passing execution does not settle the
intermittent repository-gate failure. `git diff --check` also exited 0.

No in-scope implementation defect was identified. The fixture repair is
explicitly assigned to Spec
`0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced`, Task 02,
which declares both fixture files and the rebuild regression test. That
separately owned repair remains a follow-up; no fixture, production workflow,
Verification configuration, Task status or other Task file was changed here.
Only this Result was updated during the feedback turn. The full configured
Verification remains with the Daemon, and its intermittent blocker remains
unresolved by this slice. No gate was rerun and no commit, push or pull
request was made.

## Carry-forward provenance

- Source Run: `run_20261016T175249Z_782526ccf3bf55ac`
- Source commit: `7cdf03d3e43a1d151d5931ca5e73086951d9b7ba`
