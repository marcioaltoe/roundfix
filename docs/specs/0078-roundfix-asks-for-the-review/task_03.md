---
task: task_03
spec: 0078-roundfix-asks-for-the-review
status: pending
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
   `auto_review.enabled != false AND auto_review.auto_incremental_review != false`,
   treating an absent or unreadable file as Review Source defaults, so the
   value is true.
3. MUST refuse `resolve` and `watch` in Preflight Validation when
   `pushTriggersReview` equals `request_review`, exiting `2` before any Agent
   Session, Review Source mutation, commit, or push.
4. MUST name, in the refusal, the file read, the values it read, and the
   deterministic next action for each of the two refused rows.
5. MUST exempt `fetch`, which publishes nothing and pushes nothing.
6. MUST keep every existing Run working unchanged under the built-in default,
   where `request_review` is false and the Review Source defaults apply.
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

- [ ] `pushTriggersReview` true with `request_review` false runs, unchanged
      from today.
- [ ] `pushTriggersReview` false with `request_review` true runs.
- [ ] `pushTriggersReview` false with `request_review` false refuses with exit
      `2`, naming the stall.
- [ ] `pushTriggersReview` true with `request_review` true refuses with exit
      `2`, naming the duplicate review.
- [ ] An absent `.coderabbit.yaml` is treated as Review Source defaults, so the
      built-in configuration still runs.
- [ ] `fetch` runs in every one of the four combinations.
- [ ] Each refusal names the file, the values read, and one next action.

## Context

- interface: `internal/preflight/preflight.go`
- interface: `internal/config/config.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/preflight ./internal/config -count=1 -run 'Review|Request|Coheren|Config' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the refusal tests ran and passed.
- `go test ./internal/preflight ./internal/config ./internal/cli -count=1`
  — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `git diff --name-only HEAD | grep -q "^\.coderabbit\.yaml$" && exit 1 || exit 0`
  — expected: exit 0; this Task reads that file and never writes it.

## References

- `_prd.md` → Core Features 3 and 4; Success Metric 3.
- `_techspec.md` → API Contracts; Build Order 3.
