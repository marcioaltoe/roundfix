---
task: task_03
spec: 0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity
status: completed
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

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-concurrency, golang-context, golang-error-handling, golang-safety, golang-structs-interfaces, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, evidence-gate.

Implemented this Task's review scheduling slice:

- Added `ReviewCapacity` (minimum/default `1`) and `TurnPollInterval` (default `30s`) to the engine dependencies, and wired the operator's `Delivery.ReviewConcurrency` into the command dependencies. The new command test proves a configured capacity of `3` reaches the engine.
- Each `Run` owns the slot pool, poller, cancellation and worker join. The poller admits only provisioned `reviewing` items with an existing worktree directory. Claims protect the resume prologue; returned-item selection, hold checking and claiming are atomic. Reviews keep calling the unchanged `reviewCandidate`; live errors park only their own item with the existing blocker, and cancellation leaves `reviewing` without a blocker.
- The pass hands reviews off, continues returned `archiving` items without the workspace prologue, and waits until a pass leaves the slots empty. Later passes skip previously merged items, preserving one merged cleanup per Run. Prerequisite waits and releases retain their existing next-Run behavior.
- The queued-start check runs after prerequisites and remembers a held queued item for the rest of that pass. It emits no log line. Slot admission logs the prescribed capacity-wait line once per waiting item. Slot take, start, end and poll events report through `reviewSlotsChanged`.

Acceptance evidence:

| Criterion | Focused evidence |
| --- | --- |
| A retried review starts beside implementation; its worktree is not prepared again | `TestAReviewHandedToTheQueueStartsWhileAnotherItemImplements` holds implementation until the retried review starts, keeps that review active until the other item archives, observes one review per item and zero `UseItemBranch` calls for `retried`, and checks the default poll interval. |
| Review Capacity bounds concurrency; one capacity-wait line | `TestReviewCapacityBoundsConcurrentReviews` exercises three provisioned reviews at capacity `2`, verifies the maximum and the third review's start after a prior review ends, and checks the exact wait line occurs once. Its `capacity below one` subtest exercises dependency capacity `0` and observes maximum concurrency `1`. |
| Queued start order and unchanged existing behavior | `TestAQueuedItemWaitsForAnEarlierItemInReview` verifies no second branch during the first review and the complete event order of inline capacity-one delivery. Its `an earlier blocker ends during the pass` subtest proves `second` starts before `third` after `first` parks mid-pass. All existing delivery and matched CLI tests pass without edits. |
| Review error and owner cancellation | `TestAFailedReviewParksOnlyItsItem` checks `delivery-error: run pre-PR review: review exploded` and the next item's merge. `TestStoppingTheOwnerLeavesAnItemInReviewAtReviewing` verifies cancellation matching, review shutdown before Run returns, and the unchanged reviewing stage and empty blocker. |

Required characterization, run before changing code; both passed:

- `go test -count=1 -timeout 20m ./internal/delivery`
- `go test -count=1 -timeout 30m -run 'Deliver|Delivery|Queue' ./internal/cli`

Focused implementation checks passed:

- `go test -race -count=1 -timeout 3m ./internal/delivery` — final source and fixtures, including every existing delivery test and all five new tests.
- `go test -race -count=20 -timeout 3m -run 'TestAReviewHanded|TestReviewCapacityBounds|TestAQueuedItemWaits|TestAFailedReview|TestStoppingTheOwner' ./internal/delivery` — repeated scheduling, capacity, ordering, error and cancellation coverage.
- `go test -count=1 -timeout 5m -run 'Deliver|Delivery|Queue' ./internal/cli` — existing command tests plus the new wiring test.
- `go test -count=1 -timeout 3m -run '^TestTheDeliveryEngineGetsTheReviewCapacity$' ./internal/cli`
- `git diff --check`; byte comparison against `HEAD` confirms `Engine.Retry` and `reviewCandidate` remain unchanged. No existing test or fixture was edited.

Sabotages 1–5 were applied to separate initialized temporary repository copies. Every command used `go test -count=1 -timeout 3m -run '<named test selector>' ./internal/delivery`; failures came from the named behavioral tests, not compilation or fixture setup. Logs are under `/private/tmp/roundfix-task03-sabotages-w971ycj4/` (the `*.final.log` files cover the final implementation).

| Sabotage | Named test and input | Observed failure |
| --- | --- | --- |
| 1: inline review, with the independent poll handoff removed | `TestAReviewHandedToTheQueueStartsWhileAnotherItemImplements`; `implementing` at running before provisioned `retried` at reviewing | The retried review never starts while implementation waits; the named test fails its event wait. |
| 2: ignore capacity when starting waiting reviews | `TestReviewCapacityBoundsConcurrentReviews`; three reviewing items, capacity `2`, and its capacity-`0` subtest | Maximum concurrency becomes `3` instead of `2`/`1`; the third review starts before either earlier review ends. |
| 3: remove queued-start rule | `TestAQueuedItemWaitsForAnEarlierItemInReview`; `first` and `second` queued | `second` creates its branch during the first review and the complete event sequence differs from inline order. |
| 3, same-pass clause only | `TestAQueuedItemWaitsForAnEarlierItemInReview/an_earlier_blocker_ends_during_the_pass`; three queued items, first review parks during third's prerequisite read | `third:create` precedes `second:create`. |
| 4: prepare a worktree held by the slots | `TestAReviewHandedToTheQueueStartsWhileAnotherItemImplements`; the same running/reviewing pair | `UseItemBranch` for `retried` is called twice instead of zero times. |
| 5: park cancelled reviews using a context that permits the write | `TestStoppingTheOwnerLeavesAnItemInReviewAtReviewing`; cancel the owner after its review starts | The item becomes parked with `delivery-error: run pre-PR review: context canceled`. |

The authored Verification command and platform builds remain for the Daemon. The pre-existing `status: in_progress` was preserved, and the Task Graph was not edited. No commit, push or pull request was made.
