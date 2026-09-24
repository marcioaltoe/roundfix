---
task: task_01
spec: 0156-a-delivery-loop-that-outlives-the-session
status: pending
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
