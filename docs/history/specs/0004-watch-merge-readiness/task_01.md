---
task: task_01
spec: 0004-watch-merge-readiness
status: completed
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

- [x] First-check-before-sleep ordering in the settled wait
- [x] Pre-settled quiet-period skip with the mid-Run case preserved
- [x] Clock-stepped tests asserting call ordering and elapsed time

## Acceptance Criteria

- [x] With a fake clock, a pre-settled review reaches the fetch call with
      zero sleeps recorded; a review that settles on the third poll records
      exactly two interval sleeps before it.
- [x] The pre-settled case records no quiet-period sleep; the settles-mid-Run
      case records exactly one.
- [x] All existing watch outcome tests pass unchanged.

## Verification

- `rtk go test ./internal/watch/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 1; Core Feature 1; Success Metrics. `_techspec.md` →
Watch loop changes, Build Order 1. Dogfood finding 18.

## Result

Implemented the watch-loop ordering slice in `internal/watch`: the settled
wait now reports how many Review Source status checks occurred, and `Run`
skips the quiet period only when round 1 settled on its first status check.
Reviews that settle after one or more polls still run the quiet period before
fetching.

Evidence:

- `TestRunSkipsQuietPeriodWhenReviewAlreadySettledAtStart` verifies the
  pre-settled path reaches fetch with no recorded sleeps and records no
  quiet-period sleep.
- `TestRunSleepsBetweenStatusChecksAndKeepsQuietPeriodWhenReviewSettlesDuringRun`
  verifies status calls happen before sleep, then after one and two
  `poll_interval` sleeps, and the final sleep list contains exactly one
  quiet-period sleep.
- Existing watch outcome assertions for timeout, Max Rounds, Run Budget, and
  Unresolved still pass under `rtk go test ./internal/watch/`.

Verification:

- `rtk go test ./internal/watch/`: passed, 7 tests.
- `rtk go test ./...`: passed, 481 tests in 16 packages.
- `rtk make verify`: passed; full Go suite, Roundfix skill check, and build
  completed successfully.

Follow-ups: none.
