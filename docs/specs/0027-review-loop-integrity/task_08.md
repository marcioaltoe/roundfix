---
task: task_08
spec: 0027-review-loop-integrity
status: pending
type: backend
complexity: medium
---

# Task 08: Confirm Merge-Ready through a grace window or end Clean Unverified

## Overview

Close the false-Clean race: after the Final Push, a missing Review Source check is treated like a pending one and polled through a configurable grace window. A check that appears flows through the existing merge-ready logic; a window exhausted with no check ends the Run Clean Unverified — the distinct outcome and exit code prepared earlier — with the report naming the next action.

## Requirements

1. MUST add a check grace-period duration to the watch configuration group (builtin default 5 minutes, YAML overlay, default-config rendering), following the existing watch key patterns.
2. MUST change merge-ready confirmation so a missing check keeps polling at the existing poll interval until the grace period elapses, instead of returning ready immediately.
3. MUST end the Run with the Clean Unverified outcome when the grace period exhausts with no check observed, replacing the current "treating Run as Clean" note path; the report names the outcome and the next action (confirm the pull request's Review Source check before merging).
4. MUST keep existing behavior when a check appears: success with no new Review Issues ends Clean; new Review Issues start the next Round; failure/timeout paths unchanged.
5. MUST leave until-clean-disabled watch behavior unchanged.

## Subtasks

- [ ] Add the config key across builtin defaults, overlay, and default YAML
- [ ] Rework the missing-check branch of merge-ready confirmation into the polling loop with grace-window expiry
- [ ] Thread the unverified result through the watch result to the new terminal outcome and exit code
- [ ] Update report/stderr wording for the new outcome
- [ ] Table-test the confirmation matrix with the fake clock: missing→appears-success, missing→appears-failure, missing→window exhausted, pending→timeout

## Acceptance Criteria

- [ ] A check that never appears ends the Run Clean Unverified with exit code 3 and a report naming the next action
- [ ] A check that appears late (within the window) and succeeds ends the Run Clean with exit code 0
- [ ] The grace period is configurable and defaults to 5 minutes
- [ ] The old missing-check warning path no longer exists
- [ ] The full test suite passes

## Context

- interface: `internal/watch/watch.go`
- interface: `internal/config/config.go`
- interface: `internal/cli/cli.go`

## Verification

- `rg -q "check_grace_period" internal/config` — expected: exit 0 (config key exists)
- `go test ./internal/watch/... ./internal/config/... ./internal/cli/...` — expected: all tests pass
- `go build ./...` — expected: clean build

## References

`_prd.md` → Goal 2, User Story 5, Core Feature 5; `_techspec.md` → Build Order 6, Interfaces (confirmResult), Data Models (Config), Decisions (grace period key); ADR-0043.
