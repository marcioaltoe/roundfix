---
task: task_03
spec: 0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity
status: pending
type: backend
complexity: high
---

# Task 03: The Queue Owner reviews items in review slots beside its one item, and keeps the start order

## Overview

The Delivery Engine reviews inline in the one loop that drives each item from
its stage to `merged` or `parked`, so an item a Delivery Retry returned to
`reviewing` waits behind the whole delivery of the item the loop is in
(`_techspec.md` → Current behavior 2 and 3); Fiscus saw five such items while
one Implement Run held the owner ("5 itens estavam em `reviewing` sem nenhum
processo de review rodando", reported by fiscus on 2026-11-05,
[references](references/2026-11-05-fiscus-delivery-queue-takes-days.md)).
After this Task each `Run` holds review slots and a poller: items in
`reviewing` are reviewed in position order up to the Review Capacity while
the loop keeps advancing one item at a time through every other stage, a
queued item waits while an earlier item is in review, and a reviewed item
continues without its worktree being prepared again.

## Requirements

1. MUST, before changing code, run
   `go test -count=1 -timeout 20m ./internal/delivery` and
   `go test -count=1 -timeout 30m -run 'Deliver|Delivery|Queue' ./internal/cli`
   and record in the Result that both pass; every one of those tests MUST
   keep passing unchanged (`_techspec.md` → Current behavior 7).
2. MUST add `ReviewCapacity`, `TurnPollInterval` and their defaults to
   `EngineDependencies` and `Engine` as `_techspec.md` → API Contract 3
   states (the `TurnRecorder` field is task_04's), and set
   `ReviewCapacity: loaded.Config.Delivery.ReviewConcurrency` in
   `commandDeliveryEngineDependencies` (`internal/cli/deliver_workflow.go`).
3. MUST create `internal/delivery/review_slots.go` with the review slots,
   their `state(slug)` answer and the poller, reporting every take, start,
   end and poll through `Engine.reviewSlotsChanged()` in `engine.go`, which
   wakes a pass waiting on the slots, and change `Run` and `advanceItem` in
   `internal/delivery/engine.go`, exactly as API Contract 4 states: the slots
   and the poller belong to one `Run` call and end before it returns; the pass
   first continues an item the slots returned, then skips an item the slots
   hold and claims any other `reviewing` item before its prologue; advancing hands an item that reaches `reviewing` to the slots
   and moves on; a review calls `reviewCandidate` unchanged and parks a live
   error as `delivery-error: <error>`; a reviewed item at `archiving`
   continues through the stage loop alone; the poller hands provisioned
   items at `reviewing` whose worktree directory exists; a queued item waits,
   with no log line, while an earlier item is between `running` and
   `merging` or after the rule held an earlier queued item in the same pass;
   and `Run` loops passes until the slots stay empty, skipping in
   later passes the items merged when the pass began. Every goroutine is
   owned by `Run`, which waits for it, and stops on the `Run` context.
4. MUST log, once when an item is handed to the slots while every slot is
   held, exactly
   `roundfix: Delivery Queue item <slug>: waiting for a review slot; <held> of <capacity> review slot(s) held`,
   and no other new console line.
5. MUST create `internal/delivery/review_slots_test.go`, opening with the
   `Suite`, `Invariant`, `Boundary IN` and `Boundary OUT` comment lines, each
   top-level test calling `t.Parallel()` first, driving `Engine.Run` over
   `openDeliveryEngineStore` with fakes that guard their state with a mutex,
   item worktrees under `t.TempDir()` and `WorktreeProvisioned` set where an
   item is seeded past `queued`, `TurnPollInterval` of `10 * time.Millisecond`
   where a test needs the poller, and every wait through `testwait.Until`;
   with these tests:
   - `TestAReviewHandedToTheQueueStartsWhileAnotherItemImplements`: a queue
     of `implementing` at `running` (position 0) and `retried` at `reviewing`
     (position 1); the fake `RunSpec` of `implementing` blocks until the fake
     `Review` of `retried` has started, then returns a Clean result; both
     items reach `merged`; `UseItemBranch` is never called for `retried`;
     `Review` is called once per item; and `NewEngine` with no
     `TurnPollInterval` keeps `30 * time.Second`.
   - `TestReviewCapacityBoundsConcurrentReviews`: three items at `reviewing`
     with `ReviewCapacity: 2`; the fake `Review` records the largest number
     running at once and holds the first two until both have started; all
     three reach `merged`, at most two ran at once, the third started only
     after one ended, and the console log carries
     `roundfix: Delivery Queue item <third>: waiting for a review slot; 2 of 2 review slot(s) held`
     exactly once. A subtest `capacity below one` with `ReviewCapacity: 0`
     never runs two reviews at once.
   - `TestAQueuedItemWaitsForAnEarlierItemInReview`: `first` and `second`
     both `queued`; the fake `Review` of `first` blocks until released and,
     while it blocks, the fake recorded no `CreateItemBranch` for `second`;
     after release both reach `merged`, `second` after `first`, and the event
     order equals the order a capacity of `1` gives with the review inline
     today. A subtest `an earlier blocker ends during the pass` queues
     `first`, `second` and `third`: the fake prerequisite reader, when called
     for `third`, signals the fake `Review` of `first`, which waits for that
     signal and then returns findings, and the reader waits through
     `testwait.Poll` until `first` is `parked` before reporting no
     prerequisite; `second` creates its item branch before `third` does.
   - `TestAFailedReviewParksOnlyItsItem`: the fake `Review` of `broken`
     returns `errors.New("review exploded")`; `broken` parks as
     `delivery-error: run pre-PR review: review exploded` and the next item
     reaches `merged`.
   - `TestStoppingTheOwnerLeavesAnItemInReviewAtReviewing`: the fake `Review`
     of an item blocks until its context ends; cancelling the `Run` context
     makes `Run` return an error wrapping `context.Canceled` only after the
     fake `Review` returned, and the item is still at `reviewing` with an
     empty blocker.
   Add `TestTheDeliveryEngineGetsTheReviewCapacity` to a new
   `internal/cli/deliver_review_capacity_test.go` (`t.Parallel()` first, the
   four comment lines): `commandDeliveryEngineDependencies` over a
   `roundconfig.Loaded` whose `Delivery.ReviewConcurrency` is `3` returns
   `ReviewCapacity` `3`.
6. MUST run sabotages 1 to 5 of `_techspec.md` → Testing Approach on a
   temporary copy, each `go test` with `-timeout 3m`, and record which named
   test fails and on which input.
7. MUST NOT change `reviewCandidate`, `Engine.Retry`, any blocker text, Park
   Class, Delivery Stage, the Run Database, `deliver status`, the
   configuration, any document or skill, the Makefile, CI workflows, `go.mod`
   or `go.sum`.

## Subtasks

- [ ] Characterize the engine and `deliver` suites.
- [ ] Add the capacity and poll interval and wire the capacity.
- [ ] Build the review slots, the poller and the claim, and hand reviews to them from the pass.
- [ ] Add the queued-start rule and the continuation of a reviewed item.
- [ ] Add the slot tests and the wiring test, and run the five sabotages.

## Acceptance Criteria

- [ ] An item retried to `reviewing` is reviewed while another item's Implement Run runs, and its worktree is not prepared again.
- [ ] At most the Review Capacity of reviews run at once, and a waiting item logs its wait once.
- [ ] A queued item does not start while an earlier item is in review, and every existing engine and `deliver` test passes unchanged.

## Context

- interface: `internal/delivery/engine.go`
- creates: `internal/delivery/review_slots.go`
- creates: `internal/delivery/review_slots_test.go`
- interface: `internal/delivery/engine_test.go`
- interface: `internal/testwait/testwait.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_review_capacity_test.go`
- interface: `internal/config/config.go`

## Consumers

- `Engine.Run` scheduling and `EngineDependencies` — consumers: `./internal/delivery`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAReviewHandedToTheQueueStartsWhileAnotherItemImplements TestReviewCapacityBoundsConcurrentReviews TestAQueuedItemWaitsForAnEarlierItemInReview TestAFailedReviewParksOnlyItsItem TestStoppingTheOwnerLeavesAnItemInReviewAtReviewing; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 30m -v -run 'Deliver|Delivery|Queue' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheDeliveryEngineGetsTheReviewCapacity ' || { echo 'missing pass: TestTheDeliveryEngineGetsTheReviewCapacity' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/delivery ./internal/cli && GOOS=linux go build -buildvcs=false ./internal/delivery ./internal/cli` — expected: exit 0. Before this Task none of the six new tests exists, so the first pass line is missing; afterwards reviews run in slots as API Contract 4 states, every existing Delivery Engine test and every `deliver` command test passes unchanged, and both platforms build.

## References

- `_prd.md` → Goals; User Story 2; User Story 3; User Story 4; Core Feature 4; Core Feature 5; Core Feature 7; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → Current behavior 2; Current behavior 3; Current behavior 7; API Contract 3; API Contract 4; Invariants 2-6; Testing Approach; Build Order 2; ADR-0304; ADR-0225
