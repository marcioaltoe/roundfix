---
title: "Run Event Journal: schema, append and cursor-paged reads"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 2
  - 8
  - 9
  - 10
  - 11
  - 12
  - 14
  - 15
  - 16
blocked_by:
  - 03-runevent-seam-runner-and-writersink.md
---

# Run Event Journal: schema, append and cursor-paged reads

## Parent

.scratch/roundfix-event-journal/02-run-event-journal-prd.md

## What to build

The durable Run Event Journal PRD: journal table in the Run Database with migration, per-Run monotonic cursor
allocated in the insert transaction, normalized columns for filtering (kind, source, Batch, issue ref, tool
ID/state, summary, timestamp), raw JSON payload stored as plain TEXT per ADR 0008, store APIs for
append/append-batch/list-after-cursor with limit, single-writer connection discipline (immediate transactions, WAL
enabled at creation before any reader, busy timeouts), and a data-version-based change-detection signal for pollers.

Verifiable end-to-end through store tests against a real temporary Run Database.

## Acceptance criteria

- [ ] Migration adds the journal schema; existing Run Database creation still works
- [ ] Cursors increase monotonically per Run; append to a missing Run fails clearly
- [ ] List-after-cursor returns only newer events and respects limits
- [ ] Structured payload JSON round-trips byte-exact
- [ ] Single-writer discipline: one writer connection, immediate transactions, WAL before readers, busy timeouts on all connections
- [ ] A second read-only connection pages events while the writer appends
- [ ] Change-detection signal exposed for pollers without reading rows

## Blocked by

- 03-runevent-seam-runner-and-writersink.md
