---
task: task_02
spec: 0055-owner-identity-without-fork
status: completed
type: backend
complexity: medium
---

# Task 02: Separate an unreadable identity from a proven mismatch

## Overview

`proveOwner` collapses two different facts into one error: the host could not
answer, and the owner is provably someone else. Only the second justifies
sending an operator to investigate PID reuse. Split them, keeping the
proven-mismatch refusal exactly as it behaves today.

## Requirements

1. MUST add an `ErrOwnerIdentityUnreadable` sentinel alongside
   `ErrOwnerProcessIdentityUnproven`.
2. MUST classify a failed identity read on a still-present process as
   unreadable, wrapping the host error, with a diagnostic naming the resource
   failure and the next action.
3. MUST classify a recorded token this platform could not have produced — no
   platform prefix, or another platform's prefix — as unreadable, never as a
   mismatch, so a Run created by the previous `ps` implementation is not reported
   as a reused PID.
4. MUST keep classifying a comparable token that differs from the live token as
   `ErrOwnerProcessIdentityUnproven`, refusing exactly as today.
5. MUST preserve the existing exited-between-checks recovery ahead of both
   classifications, so proven absence still wins over any read failure.
6. MUST keep every other existing proof outcome unchanged: reuse refusal,
   matching identity, absent-owner-as-proof, and legacy no-token degradation
   (ADR-0044).

## Subtasks

- [ ] Add the sentinel and the two classifications.
- [ ] Add the non-comparable-token rule with its own diagnostic.
- [ ] Extend the proof table with every new outcome and re-assert the preserved
      ones.

## Acceptance Criteria

- [ ] A read failure on a live process yields `ErrOwnerIdentityUnreadable` with
      the host error in the message.
- [ ] A recorded token with no platform prefix yields
      `ErrOwnerIdentityUnreadable`, not a mismatch.
- [ ] A comparable token that differs yields
      `ErrOwnerProcessIdentityUnproven` and refuses.
- [ ] A process that exits between the liveness check and the identity read is
      still proven absent.
- [ ] An empty recorded token still degrades to PID-only protection.

## Context

- interface: `internal/store/process.go`
- interface: `internal/store/process_unix_test.go`
- interface: `internal/cli/orphan_unix_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/store/ ./internal/cli/` — expected: pass,
  including ADR-0044 reclaim behavior.
- `make verify` — expected: exit 0.

## References

`_prd.md` → Goal 2, Story 2, Feature 2; `_techspec.md` → Build Order 2,
Classification at the proof, Risks (token format change).

## Result

Implemented separate ownership-proof classifications. A failed live identity
read now wraps both `ErrOwnerIdentityUnreadable` and the host error, and tells
the operator to resolve the host resource failure before retrying. Recorded
tokens without the current platform prefix now report an unreadable,
non-comparable identity; only a differently valued token with the current
platform prefix reports `ErrOwnerProcessIdentityUnproven`. The existing
absence-before-classification and empty-token PID-only paths remain in place.

Focused checks and acceptance evidence:

- Red starting point: the focused legacy-token case failed before the
  implementation because the unprefixed token still matched
  `ErrOwnerProcessIdentityUnproven`.
- Live read failure: the proof-table case `live identity read failure is
  unreadable` passed and asserted `errors.Is` for both
  `ErrOwnerIdentityUnreadable` and the synthetic `resource temporarily
  unavailable` host error. Its diagnostic contains the host error and the
  `resolve the host resource failure, then retry` next action.
- Non-comparable token: the proof-table cases for a legacy token without a
  prefix and a token with another platform's prefix passed. The real-process
  `TestOwnerProcessControllerProveOwnerDoesNotReportLegacyIdentityAsReusedPID`
  case also passed and rejected matching
  `ErrOwnerProcessIdentityUnproven`.
- Comparable mismatch: the proof-table case kept
  `ErrOwnerProcessIdentityUnproven` distinct from
  `ErrOwnerIdentityUnreadable`. The real store refusal case and
  `TestRunForceStopOwnerPIDReuseFailsClosed` passed with a current-platform
  token, leaving the live process unsignaled and the Active Run lock retained.
- Exited between checks: the proof table returned proven absence when the
  first liveness read reported present, the identity read failed, and the
  second liveness read reported absent.
- Empty recorded token: the proof table passed without reading identity, and
  `TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner` passed the
  existing ADR-0044 PID-only Force Stop flow.
- Preserved outcomes: the focused owner-process run also passed matching
  identity, initial absence, current-process refusal, graceful termination,
  and force-kill cases.
- `GOCACHE=<repo>/.gocache rtk go test ./internal/store -run
  '^(TestOwnerProcessController|TestOwnerProcessIdentity)'` passed 24 cases.
- `GOCACHE=<repo>/.gocache rtk go test ./internal/cli -run
  '^(TestRunForceStopOwnerPIDReuseFailsClosed|TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner)$'`
  passed 2 cases.
- `rtk gofmt -w internal/store/process.go internal/store/process_unix_test.go
  internal/cli/orphan_unix_test.go` exited 0, and `rtk git diff --check`
  exited 0.

The commands under `## Verification` were not run; Daemon Verification owns
those commands and the Task verdict.
