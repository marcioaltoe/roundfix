---
task: task_02
spec: 0180-a-prepared-queue-that-revalidates-before-each-spec
status: completed
type: backend
complexity: high
---

# Task 02: A Delivery Queue records its limits and enforces them, and an item can carry a warning

## Overview

A Delivery Queue in `internal/store/delivery.go` records only its items and its owner, so a detached owner keeps starting Specs for as long as items remain, and `Engine.Retry` in `internal/delivery/engine.go` returns an item to the queue however many times it is asked. This Task records a deadline and a per-item retry limit with the queue, in one new Run Database schema version. The engine enforces both: no queued item starts after the deadline, and no retry passes the limit. The same schema version gives each item a durable warning field, which task_01's Delivery Revalidation writes. The limits are written once when the queue is recorded and read by the detached owner and by every `deliver retry`, in separate processes that share the Run Database.

## Requirements

1. MUST raise `schemaVersion` in `internal/store/store.go` by exactly one from its value on the Task's starting main. A fresh database MUST create `delivery_queues.deadline_unix`, `delivery_queues.max_retries` and `delivery_queue_items.retry_count`, each `INTEGER NOT NULL DEFAULT 0`, and `delivery_queue_items.warning TEXT NOT NULL DEFAULT ''`. Every older schema that already has the delivery tables MUST gain them through idempotent column-existence checks, as `deliveryWorktreeMigrationStatements` does. Tests MUST read `schemaVersion`, never a literal version.
2. MUST add `DeliveryQueueLimits` and `CreateDeliveryQueueWithLimits` to `internal/store/delivery.go` as the TechSpec states. `CreateDeliveryQueue` keeps its signature and delegates with zero limits. A negative `MaxRetries` is refused. The deadline is stored as UTC Unix seconds, where zero means none, and read back in UTC.
3. MUST add `Limits DeliveryQueueLimits` to `DeliveryQueue`, and `RetryCount int` and `Warning string` to `DeliveryQueueItem`, all read by `DeliveryQueue`. `UpdateDeliveryQueueItem` MUST persist `Warning` with the other item fields, while `RetryDeliveryQueueItem` MUST leave it unchanged.
4. MUST make `RetryDeliveryQueueItem`, keeping its signature, read the queue's `max_retries` and the item's `retry_count` in its write transaction. When the limit is non-zero and reached, it refuses with an error wrapping a new exported `ErrDeliveryRetryLimit` and changes nothing. Otherwise it increments `retry_count` in the same `UPDATE` that moves the item.
5. MUST add `BlockerQueueDeadline = "queue-deadline"` to `internal/delivery/engine.go`. `Engine.Run` MUST park an item that is still `queued` as `queue-deadline` when the queue has a deadline and `engine.clock.Now()` is not before it, without calling `CreateItemBranch`. An item in any later stage advances as before.
6. MUST make `Engine.Retry` refuse a `queue-deadline` item before `UseItemBranch`, naming the deadline and `roundfix deliver start`. It MUST refuse an item whose `RetryCount` has reached a non-zero `MaxRetries` before `UseItemBranch` or carry-forward, with an error wrapping `store.ErrDeliveryRetryLimit`. Both refusals leave the stored item unchanged.
7. MUST change no exported function signature and update only the tests this change invalidates, naming each in the Result; `TestJournalConsumerCorpusReplaysEveryConsumer` and the existing Delivery Queue migration tests stay green.
8. MUST put the store tests in `internal/store/delivery_limits_test.go` and the engine tests in `internal/delivery/limits_test.go`. The engine tests build every engine through the existing helpers `newTestDeliveryEngineWithWait` in `internal/delivery/engine_test.go` and `newRetryDeliveryEngine` in `internal/delivery/retry_test.go`, never through their own `NewEngine` call, and drive the deadline through the injected `Clock`, with no wall-clock wait.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A queue recorded with limits reads them back after reopening, and a queue recorded without limits reads zero for both.
- [ ] A negative retry limit is refused and records no queue.
- [ ] An item warning written through `UpdateDeliveryQueueItem` reads back after reopening and survives a retry, and an item never given one reads empty.
- [ ] A retry increments the item's retry count; a retry at the limit is refused and leaves the item byte-for-byte unchanged.
- [ ] A database at the previous schema version with a recorded queue migrates to `schemaVersion`, keeps its rows and gains the three columns with their defaults.
- [ ] A queued item at the deadline parks as `queue-deadline` without a worktree, a queued item before the deadline starts, and an item past `queued` advances after the deadline.
- [ ] `Engine.Retry` refuses a `queue-deadline` item and an item at its retry limit before any carry-forward, leaving each unchanged.

## Context

- interface: `internal/store/store.go`
- interface: `internal/store/delivery.go`
- creates: `internal/store/delivery_limits_test.go`
- interface: `internal/delivery/engine.go`
- creates: `internal/delivery/limits_test.go`
- interface: `internal/store/delivery_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestDeliveryQueueRecordsItsLimits|TestDeliveryQueueWithoutLimitsRecordsNone|TestDeliveryQueueRefusesANegativeRetryLimit|TestRetryIncrementsTheItemRetryCount|TestRetryAtTheLimitIsRefusedAndLeavesTheItemUnchanged|TestThePreviousSchemaGainsTheLimitColumns|TestDeliveryQueueItemWarningRoundTrips|TestOpenMigratesDeliveryQueueAddingWorktreeProvisioning|TestOpenMigratesV14DeliveryQueueAddingOwnerAndItemBranch|TestRetryDeliveryQueueItemReentersAParkedItem|TestJournalConsumerCorpusReplaysEveryConsumer|TestAQueuedItemParksAtTheDeadlineWithoutAWorktree|TestAQueuedItemStartsBeforeTheDeadline|TestAnItemPastQueuedAdvancesAfterTheDeadline|TestRetryRefusesAQueueDeadlineItem|TestRetryRefusesAnItemAtItsRetryLimitBeforeCarryForward|TestRetryCarriesForwardBeforeReenteringTheRun)$" ./internal/store ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDeliveryQueueRecordsItsLimits TestDeliveryQueueWithoutLimitsRecordsNone TestDeliveryQueueRefusesANegativeRetryLimit TestRetryIncrementsTheItemRetryCount TestRetryAtTheLimitIsRefusedAndLeavesTheItemUnchanged TestThePreviousSchemaGainsTheLimitColumns TestDeliveryQueueItemWarningRoundTrips TestOpenMigratesDeliveryQueueAddingWorktreeProvisioning TestOpenMigratesV14DeliveryQueueAddingOwnerAndItemBranch TestRetryDeliveryQueueItemReentersAParkedItem TestJournalConsumerCorpusReplaysEveryConsumer TestAQueuedItemParksAtTheDeadlineWithoutAWorktree TestAQueuedItemStartsBeforeTheDeadline TestAnItemPastQueuedAdvancesAfterTheDeadline TestRetryRefusesAQueueDeadlineItem TestRetryRefusesAnItemAtItsRetryLimitBeforeCarryForward TestRetryCarriesForwardBeforeReenteringTheRun; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the twelve new named tests exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Queue limits and the item warning in the store and the engine
- `_prd.md` → Goal 4; Core Feature 4; Success Metric 3
- `_techspec.md` → API Contract 6; API Contract 8; Testing Approach 2

## Result

Implemented schema version 21 (one version after the starting version 20),
durable queue limits, item retry counts and warnings, transactional retry-limit
enforcement, and the Delivery Engine deadline and retry refusals. Existing
exported function signatures remain unchanged; `CreateDeliveryQueue` delegates
to the new `CreateDeliveryQueueWithLimits` entry point.

Acceptance evidence:

- `TestDeliveryQueueRecordsItsLimits` and
  `TestDeliveryQueueWithoutLimitsRecordsNone` reopen the Run Database and read
  the UTC-second deadline and retry limit, including both zero values.
- `TestDeliveryQueueRefusesANegativeRetryLimit` observes the validation error
  and confirms that no queue was recorded.
- `TestDeliveryQueueItemWarningRoundTrips` reopens before and after retrying,
  proves the warning survives, and proves an unwarned item remains empty.
- `TestRetryIncrementsTheItemRetryCount` observes the atomic increment;
  `TestRetryAtTheLimitIsRefusedAndLeavesTheItemUnchanged` matches
  `ErrDeliveryRetryLimit` and compares the complete item before and after the
  refusal.
- `TestThePreviousSchemaGainsTheLimitColumns` uses `schemaVersion - 1`, keeps
  the recorded queue row through migration, and observes zero defaults for the
  deadline, retry limit, retry count and warning at `schemaVersion`.
- `TestAQueuedItemParksAtTheDeadlineWithoutAWorktree`,
  `TestAQueuedItemStartsBeforeTheDeadline` and
  `TestAnItemPastQueuedAdvancesAfterTheDeadline` drive the injected Clock and
  cover the three deadline boundaries without wall-clock waits.
- `TestRetryRefusesAQueueDeadlineItem` and
  `TestRetryRefusesAnItemAtItsRetryLimitBeforeCarryForward` observe no
  `UseItemBranch` or recovery call, match the required error information and
  compare the complete stored item before and after each refusal.

Focused checks:

- The initial store test run failed to compile on the missing
  `CreateDeliveryQueueWithLimits`, `DeliveryQueueLimits` and `Limits` symbols,
  providing the pre-change signal.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run
  '^(TestDeliveryQueueRecordsItsLimits|TestDeliveryQueueWithoutLimitsRecordsNone|TestDeliveryQueueRefusesANegativeRetryLimit|TestRetryIncrementsTheItemRetryCount|TestRetryAtTheLimitIsRefusedAndLeavesTheItemUnchanged|TestThePreviousSchemaGainsTheLimitColumns|TestDeliveryQueueItemWarningRoundTrips)$'
  ./internal/store` passed.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run
  '^(TestAQueuedItemParksAtTheDeadlineWithoutAWorktree|TestAQueuedItemStartsBeforeTheDeadline|TestAnItemPastQueuedAdvancesAfterTheDeadline|TestRetryRefusesAQueueDeadlineItem|TestRetryRefusesAnItemAtItsRetryLimitBeforeCarryForward)$'
  ./internal/delivery` passed.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run
  '^(TestOpenMigratesDeliveryQueueAddingItemWorktree|TestOpenMigratesDeliveryQueueAddingWorktreeProvisioning|TestOpenMigratesV14DeliveryQueueAddingOwnerAndItemBranch)$'
  ./internal/store` passed.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1
  ./internal/store` and the corresponding `./internal/delivery` command passed;
  the store package run includes
  `TestJournalConsumerCorpusReplaysEveryConsumer`.
- `rtk make verify-incremental` first encountered the sandbox's denied process
  table in two force-stop integration tests. The host-access rerun passed
  `go vet ./...`, `go test -parallel 16 ./...`, the focused skills test, the
  Roundfix skill check and the build.

The schema increase invalidated only the historical downgrade setup used by
`TestOpenMigratesDeliveryQueueAddingItemWorktree`,
`TestOpenMigratesDeliveryQueueAddingWorktreeProvisioning` and
`TestOpenMigratesV14DeliveryQueueAddingOwnerAndItemBranch`. Their fixtures now
drop the version-21 columns before declaring their historical schema version;
their assertions were not weakened or otherwise changed.

The authored `## Verification` command was not run; the Daemon owns that
verification and settlement step.
