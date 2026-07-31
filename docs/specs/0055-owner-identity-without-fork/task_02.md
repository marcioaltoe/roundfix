---
task: task_02
spec: 0055-owner-identity-without-fork
status: pending
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
