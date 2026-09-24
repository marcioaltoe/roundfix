---
task: task_03
spec: 0156-a-delivery-loop-that-outlives-the-session
status: pending
type: backend
complexity: high
---

# Task 03: The delivery engine

## Overview

The stage machine that carries each item from Run to merge in the order the maintainer set, binds publication to the reviewed head, parks blockers, and reconciles unmatched intents on resume.

## Requirements

1. MUST advance an item through running, reviewing, archiving, gating, publishing, checking and merging, in that order.
2. MUST skip the review call and record the configured omission when the pre-PR review policy is `none`.
3. MUST permit, between the reviewed head and publication, only an archive commit whose diff is exactly the move of the Spec folder, and MUST park the item as `review-stale` otherwise.
4. MUST record an intent before push, create-pull-request and merge, and a receipt after each.
5. MUST reconcile every intent without a receipt against observed state before retrying it.
6. MUST park the item as `unauthorized` when its authorization record lacks push, pull_request or merge, before any external action.
7. MUST park an item on an Unresolved Run, review findings, a blocked review, a failing gate or failing checks, and continue with the next item.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A two-Spec fixture queue reaches merged for both with no operator action.
- [ ] Stopping after a push intent or a merge intent and resuming yields one pull request and one merge.
- [ ] A foreign commit after review parks as review-stale.
- [ ] A missing merge operation parks as unauthorized with no external call.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/review.go`
- interface: `internal/cli/archive.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestDeliveryEngineMergesAQueueWithoutAnOperator$" ./internal/delivery 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestDeliveryEngineMergesAQueueWithoutAnOperator"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^TestDeliveryEngineResumesWithoutDoubleEffects$" ./internal/delivery 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestDeliveryEngineResumesWithoutDoubleEffects"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^TestDeliveryEngineParksAStaleReview$" ./internal/delivery 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestDeliveryEngineParksAStaleReview"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Components
