---
title: "Journal backward paging: RunEventsBefore"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - completed
user_stories:
  - 1
  - 15
blocked_by: []
---

# Journal backward paging: RunEventsBefore

## Parent

.scratch/roundfix-interactive-cockpit/01-interactive-live-run-view-prd.md

## What to build

The store read that scrollback stands on: `RunEventsBefore(runID, cursor, limit)` returning the `limit` events
immediately before the cursor in ascending order, symmetric to `RunEventsAfter`, with the same read-only
short-autocommit discipline. A cursor of zero (or below the first event) returns nothing; paging backward from the
tail must walk the entire journal without duplicates or gaps at page boundaries, including while a writer appends.

## Acceptance criteria

- [x] `RunEventsBefore` returns the N events strictly before the cursor, ascending, limit required positive
- [x] Backward pages compose: walking tail→head page by page yields every event exactly once
- [x] Works on a read-only connection while the writer appends; no long-lived read transaction
- [x] Missing Run / empty journal behave like the forward read (clear error / empty page)
- [x] Tests against a real temporary Run Database, including boundary cursors (first event, exact page edges)

## Blocked by

None - can start immediately

## Comments

**2026-06-10 (agent):** Added `Store.RunEventsBefore(runID, cursor, limit)` to `internal/store/journal.go`: selects
`cursor < ?` descending with LIMIT, then reverses to ascending, so a page is always "the limit events immediately
before the cursor, oldest first". Same validation surface as the forward read (Run ID required, positive limit
required, unknown Run → empty page) and the same short autocommit single-query read discipline — no transaction.
Tests against real temporary Run Databases prove: the two events immediately before a cursor come back ascending
with correct summaries; boundary cursors 0 and 1 return empty and cursor 2 returns exactly the first event; a
tail→head walk with page size 5 over 23 events (partial head page) yields every cursor exactly once with no gaps;
and a read-only `OpenReader` connection serves stable backward pages over the existing prefix while a writer
goroutine appends 20 events past the tail. Verification: `rtk go vet ./...` clean,
`rtk go test -race ./internal/store/` 20 passed, full `rtk go test ./...` 210 passed in 15 packages,
`rtk go run ./cmd/roundfix --help` green.
