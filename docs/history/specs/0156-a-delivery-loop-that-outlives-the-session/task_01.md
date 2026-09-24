---
task: task_01
spec: 0156-a-delivery-loop-that-outlives-the-session
status: completed
type: backend
complexity: medium
---

# Task 01: The queue store

## Overview

Persist the delivery queue so it survives the owner process. Items are keyed by Spec slug; each carries its stage, its blocker, and the Run, candidate commits, pull request number and merge commit it produced. Action intents and their receipts are separate rows.

## Requirements

1. MUST persist a queue with ordered items keyed by Spec slug.
2. MUST persist each item's stage, blocker reason, Run identity, candidate commits, pull request number and merge commit.
3. MUST persist action intents and receipts as separate records, so an intent without a receipt is observable.
4. MUST keep existing store data readable and existing tables unchanged.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Items, stages and blockers round-trip.
- [ ] An intent without a receipt is listed as unmatched.
- [ ] Existing store tests pass unchanged.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/store/store.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestDeliveryQueueRoundTripsItemsAndReceipts$" ./internal/store 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestDeliveryQueueRoundTripsItemsAndReceipts"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Components

## Result

Implemented schema version 14 with a repository-scoped Delivery Queue, ordered
items keyed by Spec slug, item progress fields, and separate Delivery Action
intent and receipt tables. Schema 13 migrates by creating only the new tables;
the Run Database lifecycle policy now names the owner and retention rule for
each one.

Focused checks:

- Pre-change: `GOCACHE=/tmp/roundfix-task01-gocache go test -run '^TestDeliveryQueue' ./internal/store` failed to build because the Delivery Queue API and types did not exist.
- Post-change: `GOCACHE=/tmp/roundfix-task01-gocache go test -run '^TestDeliveryQueue' ./internal/store` passed.
- `GOCACHE=/tmp/roundfix-task01-gocache go test ./internal/store` passed in 55.973 seconds.
- `GOCACHE=/tmp/roundfix-task01-gocache go vet ./internal/store` passed.
- `git diff --check` passed.

Acceptance evidence:

- Items, stages and blockers round-trip: `TestDeliveryQueueRoundTripsItemsAndReceipts` creates two ordered items, persists every requested progress field, closes and reopens the store, and compares the restored items with the originals.
- An intent without a receipt is listed as unmatched: the same test records one push receipt and leaves one merge intent without a receipt; after reopening, only the merge intent appears in `UnmatchedDeliveryActionIntents`.
- Existing store tests pass unchanged: the full `internal/store` package suite passed after the schema migration and lifecycle policy update; the migration fixture also confirms an existing Run remains readable after upgrading schema 13 to 14.

The Daemon-owned command under `## Verification` was not run.

## Carry-forward provenance

- Source Run: `run_20260924T121521Z_744099e82ec81b88`
- Source commit: `32bd39e6e62b4a801cd06bc347004c98b7f77925`
