---
task: task_03
spec: 0055-owner-identity-without-fork
status: pending
type: backend
complexity: medium
---

# Task 03: Mark and warn about a Run created without reuse protection

## Overview

When identity capture fails at Run creation the Run records no identity and
starts anyway, silently, with PID-only protection. Make that state observable:
one startup warning and a durable marker on the Run. The silent NULL degradation
survives only for legacy rows that predate the identity column.

## Requirements

1. MUST add one additive Run column, defaulting to unset, recording that
   identity capture failed at creation.
2. MUST set the marker when capture fails and leave it unset when capture
   succeeds.
3. MUST emit exactly one warning at Run start when the marker is set, naming
   that the Run has PID-only reuse protection.
4. MUST render the marker in Run inspection output so the state is queryable
   after the fact.
5. MUST leave a legacy row — NULL identity with the marker unset — degrading
   exactly as it does today, per ADR-0044.
6. MUST NOT fail Run creation because capture failed.
7. MUST NOT change the compare-and-set terminal completion contract
   (ADR-0052).

## Subtasks

- [ ] Add the column and set it at creation on capture failure.
- [ ] Emit the single startup warning.
- [ ] Render the marker in Run inspection.
- [ ] Cover set, unset, and legacy-NULL rows.

## Acceptance Criteria

- [ ] A Run created with a failing capture carries the marker and prints one
      warning.
- [ ] A Run created with a successful capture carries no marker and prints no
      warning.
- [ ] A legacy row with NULL identity and no marker behaves as it does today.
- [ ] The warning appears once, not once per read of the Run.
- [ ] Run creation succeeds in every case above.

## Context

- interface: `internal/store/store.go`
- interface: `internal/store/store_test.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/cli_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/store/ ./internal/cli/` — expected: pass,
  including the legacy-row case.

## References

`_prd.md` → Goal 3, Story 3, Feature 3; `_techspec.md` → Build Order 3,
Data Models.
