---
task: task_03
spec: 0078-roundfix-asks-for-the-review
status: completed
type: backend
complexity: medium
---

# Task 03: Refuse a pair that would strand the Run

## Overview

Two settings decide whether a pushed head ever gets reviewed: whether the
Review Source reviews pushes on its own, and whether Roundfix asks. Three of
the four combinations are coherent to a reader and only two of them work.

This slice adds the configuration and the Preflight Validation that decides the
pair before the Run starts. It is independent of task_01: the refusal must be
provable with no requester in existence.

## Requirements

1. MUST add `review_source.request_review` (bool, default `false`) and
   `review_source.request_command` (string, default `@coderabbitai review`),
   with Project Config over User Config over built-in precedence.
2. MUST derive `pushTriggersReview` from the repository's `.coderabbit.yaml` as
   `reviews.auto_review.enabled != false AND
   reviews.auto_review.auto_incremental_review != false AND
   reviews.auto_review.auto_pause_after_reviewed_commits == 0`, treating an
   absent or unreadable file as Review Source defaults, including the
   finite-pause default of `5`.
3. MUST refuse `resolve` and `watch` in Preflight Validation when
   `pushTriggersReview` equals `request_review`, exiting `2` before any Agent
   Session, Review Source mutation, commit, or push.
4. MUST name, in the refusal, the file read, the values it read, and the
   deterministic next action for each of the two refused rows.
5. MUST exempt `fetch`, which publishes nothing and pushes nothing.
6. MUST reject the built-in `request_review=false` setting when the Review
   Source's default finite pause applies; the operator must either enable
   Roundfix requests or explicitly configure zero-pause automatic reviews.
7. MUST NOT write, repair, or generate `.coderabbit.yaml`; Roundfix reads it and
   refuses, and volume control outside this repository's code is a maintainer
   decision.

## Subtasks

- [ ] Add both config keys with validation and precedence.
- [ ] Add the `.coderabbit.yaml` inspection and the `pushTriggersReview`
      predicate.
- [ ] Wire the refusal into Preflight Validation for `resolve` and `watch`.
- [ ] Table-test all four rows plus the absent-file default.

## Acceptance Criteria

- [ ] `pushTriggersReview` true with `request_review` false runs when
      `auto_pause_after_reviewed_commits` is explicitly `0`.
- [ ] `pushTriggersReview` false with `request_review` true runs.
- [ ] `pushTriggersReview` false with `request_review` false refuses with exit
      `2`, naming the stall.
- [ ] `pushTriggersReview` true with `request_review` true refuses with exit
      `2`, naming the duplicate review.
- [ ] An absent `.coderabbit.yaml` is treated as Review Source defaults, so
      built-in `request_review=false` refuses the potentially stranded Run.
- [ ] `fetch` runs in every one of the four combinations.
- [ ] Each refusal names the file, the values read, and one next action.

## Context

- interface: `internal/preflight/preflight.go`
- interface: `internal/config/config.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/preflight ./internal/config -list 'Review|Request|Coheren|Config' | grep '^Test' > /dev/null && go test ./internal/preflight ./internal/config -count=1 -run 'Review|Request|Coheren|Config'`
  — expected: exit 0; the refusal tests ran and passed.
- `go test ./internal/preflight ./internal/config ./internal/cli -count=1`
  — expected: exit 0.
- `go test -parallel 16 ./...`
  — expected: exit 0.
- `git diff --quiet HEAD -- .coderabbit.yaml`
  — expected: exit 0; this Task reads that file and never writes it.

## References

- `_prd.md` → Core Features 3 and 4; Success Metric 3.
- `_techspec.md` → API Contracts; Build Order 3.

## Result

### Implementation

- Added the two layered Review Source settings. The built-in configuration
  keeps requests disabled and uses `@coderabbitai review`; User Config and
  Project Config override both values in that order. Configuration loading
  rejects a non-boolean `request_review`, and validation rejects a blank
  `request_command`.
- Added read-only `.coderabbit.yaml` inspection under the Git root. Both
  `reviews.auto_review.enabled` and
  `reviews.auto_review.auto_incremental_review` default to `true`, while
  `reviews.auto_review.auto_pause_after_reviewed_commits` defaults to `5`.
  Automatic reviews therefore count as available after every push only when
  the pause is explicitly disabled with `0`. Roundfix does not create or
  modify the file.
- Added the equality refusal to operational Preflight Validation after Git
  inspection and before Open Pull Request resolution or push planning.
  `resolve` and `watch` now reject both the stranded and duplicate-review rows;
  `fetch` bypasses the check. The CLI passes the effective `request_review`
  value into Preflight and maps the refusal to its existing exit `2` contract.

### Focused checks

- Pre-change signal:
  `GOCACHE=$PWD/.gocache rtk go test ./internal/config ./internal/preflight -count=1 -run 'Test(BuiltinReviewRequestDefaults|LoadAppliesReviewRequestHierarchy|ValidateRejectsEmptyReviewRequestCommand|RunEnforcesReviewRequestCoherence|RunExemptsFetchFromReviewRequestCoherence|RunTreatsUnreadableCodeRabbitConfigAsDefaults)$'`
  failed because the new configuration fields and Preflight input did not
  exist.
- Task-specific post-edit check:
  `GOCACHE=$PWD/.gocache rtk go test ./internal/config ./internal/preflight ./internal/cli -count=1 -run 'Test(BuiltinReviewRequestDefaults|LoadAppliesReviewRequestHierarchy|ValidateRejectsEmptyReviewRequestCommand|LoadRejectsNonBooleanReviewRequest|RunEnforcesReviewRequestCoherence|RunExemptsFetchFromReviewRequestCoherence|RunTreatsUnreadableCodeRabbitConfigAsDefaults|RunResolveReviewRequestCoherenceRefusalExitsTwoBeforeRunCreation)$'`
  passed 21 tests across three packages.
- `GOCACHE=$PWD/.gocache rtk go test ./internal/config -count=1` passed
  169 tests.
- `GOCACHE=$PWD/.gocache rtk go test ./internal/preflight -count=1` passed
  38 tests.
- `GOCACHE=$PWD/.gocache rtk go test ./internal/cli -count=1` passed 966
  tests.
- `rtk git diff --check` exited `0` with no diagnostics.
- The Task's declared `## Verification` commands were not run; the Daemon owns
  them.

### Acceptance evidence

1. The explicit zero-pause row passes with `pushTriggersReview=true` and
   `request_review=false`; the default finite-pause row does not masquerade as
   automatic coverage for every pushed head.
2. `TestRunEnforcesReviewRequestCoherence/resolve_runs_when_pushes_do_not_trigger_reviews_and_asking_is_enabled`
   passes with `pushTriggersReview=false` and `request_review=true`.
3. The `resolve` and `watch` stranded-Run rows both reject
   `pushTriggersReview=false` with `request_review=false`, naming the wait for
   a review nobody requests and directing the operator to enable
   `review_source.request_review` in Project Config.
4. The `resolve` and `watch` duplicate-review rows both reject
   `pushTriggersReview=true` with `request_review=true`. The CLI boundary test
   observes exit `2`, no stdout, and no Run Database creation.
5. The absent- and unreadable-file rows resolve the two booleans to `true` and
   the finite pause to its default `5`, producing `pushTriggersReview=false`.
6. `TestRunExemptsFetchFromReviewRequestCoherence` passes all four predicate
   combinations.
7. Both refusal rows assert the absolute `.coderabbit.yaml` path, the resolved
   `auto_review.enabled`, `auto_incremental_review`, and
   `auto_pause_after_reviewed_commits` values, `pushTriggersReview`,
   `review_source.request_review`, and one row-specific Project Config next
   action.

### Follow-ups

None for this Task slice.
