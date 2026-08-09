---
task: task_04
spec: 0090-a-gate-that-could-have-failed
status: pending
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

- `GOCACHE="$PWD/.gocache" go test ./internal/runevent -run '^TestVerificationProbeEvent' -count=1 -v 2>&1 | grep -q '^--- PASS: TestVerificationProbeEventProjectsVacuousAndUnknownDistinctly'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestPreWorkProbePublishes' -count=1 -v 2>&1 | grep -q '^--- PASS: TestPreWorkProbePublishesEveryOffendingCommand'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestPreWorkProbePublishes' -count=1 -v 2>&1 | grep -q '^--- PASS: TestPreWorkProbePublishesNothingForAClearedTask'` — expected: exits 0.
- `grep -q 'verification_vacuous' internal/runevent/event.go && grep -q 'verification_unknown' internal/runevent/event.go` — expected: exits 0. Neither string exists before this Task.

## References

- `_prd.md` → Goal 2; Core Features, recorded rather than narrated evidence.
- `_techspec.md` → Build Order 4; API Contracts.
