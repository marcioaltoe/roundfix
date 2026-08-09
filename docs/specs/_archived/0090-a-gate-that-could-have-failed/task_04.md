---
task: task_04
spec: 0090-a-gate-that-could-have-failed
status: completed
type: backend
complexity: medium
---

# Task 04: Publish the probe's finding where a reader can see it

## Overview

A refusal that only reaches a terminal reason is a sentence in a file. The Run
Event stream is where a maintainer, the Supervisor, and any later audit read what
the Daemon actually did, so the probe publishes there: which commands it ran,
which were vacuous, and which it could not observe. This is the half of the Spec
that makes a recorded Result distinguishable from a claimed one.

## Requirements

1. MUST publish a Daemon Verification event when the probe refuses a Task,
   carrying the Task identifier and every offending command.
2. MUST publish a Daemon Verification event when a probe command's verdict could
   not be observed, carrying the command, the reason, and the diagnostic path.
3. MUST use the existing Daemon Verification event kind so current consumers of
   the stream keep working without change.
4. MUST give the two situations distinct classifications, so a reader can filter
   one without the other.
5. MUST NOT publish a per-command event for a Task the probe cleared; silence is
   the ordinary path.

## Subtasks

- [ ] Add the two classifications.
- [ ] Publish on refusal and on an unobserved probe command.
- [ ] Project both through the event stream.

## Acceptance Criteria

- [ ] A refused Task produces an event naming every offending command.
- [ ] An unobserved probe command produces an event carrying its reason and
      diagnostic path.
- [ ] The two classifications are distinct strings.
- [ ] A cleared Task produces no probe event.

## Bounded scope

This Task may create or modify only:

- `internal/runevent/event.go`
- `internal/runevent/stream.go`
- `internal/runevent/event_test.go`
- `internal/runevent/stream_test.go`
- `internal/daemon/task_engine.go`
- `internal/daemon/task_engine_test.go`
- `docs/specs/0090-a-gate-that-could-have-failed/task_04.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/runevent -run '^TestVerificationProbeEvent' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestVerificationProbeEventProjectsVacuousAndUnknownDistinctly'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestPreWorkProbePublishes' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestPreWorkProbePublishesEveryOffendingCommand'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestPreWorkProbePublishes' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestPreWorkProbePublishesNothingForAClearedTask'` — expected: exits 0.
- `grep -q 'verification_vacuous' internal/runevent/event.go && grep -q 'verification_unknown' internal/runevent/event.go` — expected: exits 0. Neither string exists before this Task.

## References

- `_prd.md` → Goal 2; Core Features, recorded rather than narrated evidence.
- `_techspec.md` → Build Order 4; API Contracts.

## Result

Published pre-work probe findings through the existing
`daemon.verification` event kind. A refusal produces one aggregate
`verification_vacuous` record with the Task ID and every command that exited
zero against the unchanged tree. Each unobserved command produces a
`verification_unknown` record with its Task ID, command, runner reason, and
diagnostic path. The stable Run Event stream projects those fields only for the
new classifications, so ordinary Verification records keep their existing
bounded shape. A probe with only observed non-zero commands publishes neither
classification.

Focused-check evidence:

- Red baseline: `GOCACHE="$PWD/.gocache" rtk go test ./internal/runevent ./internal/daemon -run '^$'` failed to compile after the tests were added because the two classifications and the projected command evidence fields did not exist.
- `GOCACHE="$PWD/.gocache" rtk go test ./internal/runevent ./internal/daemon` passed: 234 tests across 2 packages.
- `rtk git -c core.fsmonitor=false diff --check` exited 0.

Acceptance evidence:

- Refusal names every offending command:
  `TestPreWorkProbePublishesEveryOffendingCommand` passed with two vacuous
  commands separated by an observed non-zero command, and asserted one
  aggregate `daemon.verification` event carrying both commands and the Task ID.
- Unknown observations carry actionable evidence:
  `TestPreWorkProbePublishesUnknownCommandReasonAndDiagnosticPath` passed and
  asserted the command, the runner's reason, and the retained diagnostic path.
- The classifications are distinct strings and both project through the stable
  stream: `TestVerificationProbeEventProjectsVacuousAndUnknownDistinctly`
  passed for `verification_vacuous` and `verification_unknown`.
- Cleared Tasks stay silent:
  `TestPreWorkProbePublishesNothingForAClearedTask` passed after the Task took
  its ordinary Agent, post-Agent Verification, settlement, and commit path with
  no probe-classified event.

The commands under `## Verification` were not run; the Daemon owns them.
