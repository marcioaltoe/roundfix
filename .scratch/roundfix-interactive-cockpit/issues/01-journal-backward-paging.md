---
title: "Journal backward paging: RunEventsBefore"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
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

- [ ] `RunEventsBefore` returns the N events strictly before the cursor, ascending, limit required positive
- [ ] Backward pages compose: walking tail→head page by page yields every event exactly once
- [ ] Works on a read-only connection while the writer appends; no long-lived read transaction
- [ ] Missing Run / empty journal behave like the forward read (clear error / empty page)
- [ ] Tests against a real temporary Run Database, including boundary cursors (first event, exact page edges)

## Blocked by

None - can start immediately
