---
task: task_02
spec: 0090-a-gate-that-could-have-failed
status: pending
type: backend
complexity: medium
---

# Task 02: Separate a command's verdict from the runner's sight

## Overview

A Verification command that ran and exited non-zero has a verdict. One that timed
out, ran partially, or could not be executed has none. Today both arrive as the
same failure. This Task gives the outcome an explicit unknown cause so the
difference reaches a Task's terminal reason, and so the probe in Task 03 can tell
"this gate did not already pass" from "we could not find out".

## Requirements

1. MUST add an unknown cause to the Verification attempt outcome, carrying the
   reason and the diagnostic path.
2. MUST set that cause only when the runner could not observe a verdict, and
   never together with a command failure.
3. MUST surface the cause in the Task's terminal reason with wording a
   maintainer can act on, distinct from a command that ran and failed.
4. MUST keep a Task whose Verification could not be observed settling `failed`;
   this Task introduces no new `spec.Status`, per ADR-0111.
5. MUST break the characterization case Task 01 declared for the unobserved
   verdict, and update that case to the new behaviour in the same commit.

## Subtasks

- [ ] Add the unknown cause to the outcome and its error type.
- [ ] Classify runner errors that are not command verdicts.
- [ ] Give the terminal reason its own wording.

## Acceptance Criteria

- [ ] An unobservable Verification produces an outcome carrying the unknown
      cause and no command failure.
- [ ] A command that ran and exited non-zero still produces a command failure
      and no unknown cause.
- [ ] The two terminal reasons are distinguishable by text.
- [ ] The Task still settles `failed` in both cases.

## Bounded scope

This Task may create or modify only:

- `internal/daemon/daemon.go`
- `internal/daemon/engine.go`
- `internal/daemon/daemon_test.go`
- `internal/daemon/engine_test.go`
- `internal/daemon/verification_probe_characterization_test.go`
- `docs/specs/0090-a-gate-that-could-have-failed/task_02.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestVerificationUnknownCause' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestVerificationUnknownCauseIsSetOnlyWhenNoVerdictWasObserved'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestVerificationUnknownCause' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestVerificationUnknownCauseAndCommandFailureAreMutuallyExclusive'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestVerificationProbeCharacterization' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestVerificationProbeCharacterization'` — expected: exits 0, proving the declared break was updated to the new behaviour rather than left failing. A whole-package sweep would pass with the work absent; this names the case that must change.
- `grep -q 'UnknownCause' internal/daemon/engine.go` — expected: exits 0. This string does not exist in the file before this Task.

## References

- `_prd.md` → Goal 4.
- `_techspec.md` → Build Order 2; Data Models.
- ADR-0111.

