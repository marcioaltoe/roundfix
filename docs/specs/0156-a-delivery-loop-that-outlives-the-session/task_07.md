---
task: task_07
spec: 0156-a-delivery-loop-that-outlives-the-session
status: pending
type: backend
complexity: low
---

# Task 07: A migration ladder that applies from every earlier version

## Overview

Corrective Task from the QA gate of 2026-09-24: the repository Verification failed before the matrix was built. `migrateV13ToV14Statements` creates `delivery_queues` with `owner_pid` and `owner_identity` already present, and `migrateV14ToV15Statements` then adds both columns again, so opening any Run Database older than version 14 fails with "duplicate column name: owner_pid".

## Requirements

1. MUST make a Run Database at every earlier schema version migrate to version 15 in one pass without error.
2. MUST make a version 14 database, whose `delivery_queues` lacks the owner columns, gain them when migrated to version 15.
3. MUST leave a freshly created version 15 database byte-for-byte equivalent in schema to one reached by migration.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The existing migration tests from versions 11 and 12, and the preflight migration test, pass again.
- [ ] A version 14 database migrates to version 15 and its queue table carries the owner columns.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/store/store.go`
- interface: `internal/store/delivery.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestOpenMigratesV11RunDatabaseAddingOwnerIdentityUnproven|TestOpenMigratesV12RunDatabaseAddingRunWindows|TestOpenMigratesV14DeliveryQueueAddingOwner|TestBranchIntegrityPreflightMigratesOutdatedRunDatabase)$" ./internal/store ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestOpenMigratesV11RunDatabaseAddingOwnerIdentityUnproven TestOpenMigratesV12RunDatabaseAddingRunWindows TestOpenMigratesV14DeliveryQueueAddingOwner TestBranchIntegrityPreflightMigratesOutdatedRunDatabase; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task three of the named tests fail with "duplicate column name: owner_pid" and the fourth does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
