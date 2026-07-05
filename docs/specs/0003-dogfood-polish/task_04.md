---
task: task_04
spec: 0003-dogfood-polish
status: pending
type: backend
complexity: low
---

# Task 04: Include a bounded stderr tail in infrastructure errors

## Overview

The first acpx production failure surfaced only "Internal error" — the
adapter's actual complaint sat unseen in captured stderr. Infrastructure
errors from the agent layer must carry a bounded tail of the failing tool's
stderr so the message alone names the next action. Verifiable through
error-formatting unit tests.

## Requirements

1. MUST include the trailing portion of captured stderr in the
   infrastructure error's `Error()` string, bounded (last 10 lines or 1 KiB,
   whichever is smaller), whitespace-trimmed, and clearly delimited.
2. MUST omit the tail cleanly when stderr is empty (no dangling delimiters).
3. MUST not echo process environment or expand the error beyond the bound —
   the full stderr remains available where it is already stored.
4. SHOULD keep the existing message prefix stable so current error-matching
   tests need only additive updates.

## Subtasks

- [ ] Bounded-tail formatting in the infrastructure error
- [ ] Empty-stderr and oversized-stderr table tests
- [ ] Call-site audit: every producer fills the stderr field

## Acceptance Criteria

- [ ] A synthetic failure with multi-line stderr yields a message ending in
      the delimited tail; a 100-line stderr is truncated to the bound with a
      truncation marker.
- [ ] Empty stderr reproduces today's message exactly.
- [ ] The acpx runner's ensure/prompt/close paths all populate stderr
      (asserted via the fake rig).

## Verification

- `rtk go test ./internal/agent/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 4; Core Feature 4. `_techspec.md` → Interfaces
(InfrastructureError), Build Order 4, Risks (bounded, secret-agnostic).
Dogfood finding 27.
