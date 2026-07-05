---
task: task_01
spec: 0004-watch-merge-readiness
status: pending
type: backend
complexity: medium
---

# Task 01: Poll first, sleep after

## Overview

Reorder the watch loop so the first Review Source status check happens
immediately on start, with `poll_interval` applying only between subsequent
checks, and skip the quiet period entirely when the review was already
settled before the Run began. Verifiable through clock-stepped watch loop
tests.

## Requirements

1. MUST perform the first status check with zero elapsed wait — no sleep
   precedes the first useful call in any Round.
2. MUST apply `poll_interval` only between consecutive status checks within
   the same wait.
3. MUST skip the quiet period when the very first status check of the Run
   reports the review already settled; a review that settles during the Run
   keeps today's quiet-period behavior.
4. MUST leave the review timeout, Max Rounds accounting, Run Budget checks,
   and every outcome mapping unchanged.

## Subtasks

- [ ] First-check-before-sleep ordering in the settled wait
- [ ] Pre-settled quiet-period skip with the mid-Run case preserved
- [ ] Clock-stepped tests asserting call ordering and elapsed time

## Acceptance Criteria

- [ ] With a fake clock, a pre-settled review reaches the fetch call with
      zero sleeps recorded; a review that settles on the third poll records
      exactly two interval sleeps before it.
- [ ] The pre-settled case records no quiet-period sleep; the settles-mid-Run
      case records exactly one.
- [ ] All existing watch outcome tests pass unchanged.

## Verification

- `rtk go test ./internal/watch/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 1; Core Feature 1; Success Metrics. `_techspec.md` →
Watch loop changes, Build Order 1. Dogfood finding 18.
