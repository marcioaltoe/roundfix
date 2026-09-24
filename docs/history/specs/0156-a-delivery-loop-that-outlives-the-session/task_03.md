---
task: task_03
spec: 0156-a-delivery-loop-that-outlives-the-session
status: completed
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

## Result

Implemented the Delivery Engine as a durable per-item stage machine over the
Run Database queue. It advances queued items through Run, pre-PR review (or a
recorded configured omission), archive validation, repository gating,
publication, current-head checks and squash merge. A blocker parks only its
item, so the engine continues with the next item.

The engine records an intent before each push, pull request creation and merge,
then records the observed result as its receipt. Resume reconciles an unmatched
push intent through a new read-only remote-head observation; pull request and
merge retries use the boundary's existing observe-before-mutate behavior. The
archive transition accepts only one exact Spec-folder move whose parent is the
reviewed head, and authorization for `push`, `pull_request` and `merge` is
checked before any pull request boundary call.

Focused checks:

- Pre-change: `GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/delivery` failed to compile because the Delivery Engine types and constructor did not exist.
- `GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 ./internal/delivery` passed.
- `GOCACHE=/private/tmp/roundfix-task03-gocache go vet ./internal/delivery` passed.
- `git diff --check` passed.
- `GOCACHE=/private/tmp/roundfix-task03-gocache make verify-incremental` first reached a sandbox-blocked `api.github.com` access; the same repository-required incremental gate was rerun with network access and passed, including repository-wide vet, tests, skill checks and the binary build.

Acceptance evidence:

- A two-Spec fixture queue reaches merged for both with no operator action: `TestDeliveryEngineMergesAQueueWithoutAnOperator` passed and asserts the ordered stage-boundary calls, one enabled review, one recorded `none` omission, two pull requests, two merges and no unmatched intents.
- Stopping after a push intent or a merge intent and resuming yields one pull request and one merge: `TestDeliveryEngineResumesWithoutDoubleEffects` passed for unmatched push, pull request and merge intents; each resume reconciles observed state and leaves no unmatched intent.
- A foreign commit after review parks as `review-stale`: `TestDeliveryEngineParksAStaleReview` passed, made no pull request boundary call for the stale item and merged the following item.
- A missing merge operation parks as `unauthorized` with no external call: `TestDeliveryEngineParksUnauthorizedBeforeExternalAction` passed with zero pull request boundary calls.

Additional requirement evidence:

- `TestDeliveryEngineParksDeclaredBlockersAndContinues` passed for an Unresolved Run, review findings, blocked review, failing repository gate and failing current-head checks, then merged the next item.
- `TestPullRequestBoundaryObservesRemoteHeadWithoutPushing` passed and proves push reconciliation has a read-only observation path.

The Daemon-owned commands under `## Verification` were not run.

## Carry-forward provenance

- Source Run: `run_20260924T121521Z_744099e82ec81b88`
- Source commit: `6535fe7fdba13e49b078453e06f1dcc72631bbbb`
