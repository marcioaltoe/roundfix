---
task: task_04
spec: 0055-owner-identity-without-fork
status: pending
type: cli
complexity: medium
---

# Task 04: Fix Stop argument order and add the supervised exit

## Overview

`roundfix stop <run-id> --force` rejects its trailing flag — the same
argument-ordering defect Spec 0042 fixed for Attach. Fix it the same way, and
add the one supervised path out of an unreadable identity, so failing closed on
ignorance has an explicit, operator-driven exit.

## Requirements

1. MUST accept Stop Command flags in any position relative to the Run ID,
   reusing the parsing Spec 0042 introduced for Attach rather than a second
   mechanism.
2. MUST add an explicit `--owner-identity-unreadable` flag that permits the stop
   only when the ownership proof returned `ErrOwnerIdentityUnreadable`.
3. MUST exit `2` and signal nothing when that flag is passed while the identity
   is readable, or while the proof returned a proven mismatch.
4. MUST NOT reach the supervised path through any configuration key, environment
   variable, default, or timeout — the flag is the only entry.
5. MUST keep the proven-mismatch refusal absolute: no flag weakens it.
6. MUST leave Stop Request semantics and force-stop signaling order unchanged.

## Subtasks

- [ ] Fix the argument ordering with the Attach parsing.
- [ ] Add the flag, gated on the unreadable classification.
- [ ] Cover both argument orders, the permitted case, and both refusals.

## Acceptance Criteria

- [ ] `roundfix stop run_x --force` and `roundfix stop --force run_x` behave
      identically.
- [ ] With an unreadable identity, the flag permits the stop.
- [ ] With a readable, matching identity, the flag exits `2` and signals nothing.
- [ ] With a proven mismatch, the flag exits `2` and signals nothing.
- [ ] Without the flag, an unreadable identity still fails closed with its own
      diagnostic.

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/cli/` — expected: pass, including both refusals.
- `go run ./cmd/roundfix stop --help` — expected: the flag is documented.

## References

`_prd.md` → Goal 2 Story 4, Goal 4 Story 5, Features 4–5; `_techspec.md` →
Build Order 4, Supervised path, Risks (the supervised flag is a real hazard).
