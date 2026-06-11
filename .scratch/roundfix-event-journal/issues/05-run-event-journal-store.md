---
title: "Run Event Journal: schema, append and cursor-paged reads"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - completed
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

- [x] Migration adds the journal schema; existing Run Database creation still works
- [x] Cursors increase monotonically per Run; append to a missing Run fails clearly
- [x] List-after-cursor returns only newer events and respects limits
- [x] Structured payload JSON round-trips byte-exact
- [x] Single-writer discipline: one writer connection, immediate transactions, WAL before readers, busy timeouts on all connections
- [x] A second read-only connection pages events while the writer appends
- [x] Change-detection signal exposed for pollers without reading rows

## Blocked by

- 03-runevent-seam-runner-and-writersink.md

## Comments

**2026-06-10 (agent):** Migration v3 adds `run_events` (composite PK `run_id, cursor`; normalized columns batch,
source, kind, review_issue, tool_id, tool_state, summary, created_at; payload as plain TEXT per ADR 0008; FK to
`runs` with cascade). Store APIs in `internal/store/journal.go`: `AppendRunEvent`, `AppendRunEvents` (one immediate
transaction, cursors returned in input order), `RunEventsAfter(runID, cursor, limit)` (strictly-newer, ascending,
positive limit required), and `DataVersion` (SQLite `PRAGMA data_version` change signal for pollers). Cursor is
allocated as per-Run MAX+1 inside the insert transaction via `INSERT … RETURNING`. Connection discipline via DSN:
writer `_txlock=immediate`, `busy_timeout(5000)`, `journal_mode(WAL)` at creation, `foreign_keys(1)`, single
connection; new `OpenReader` opens `mode=ro` with busy timeout and never migrates. Tests prove WAL + busy timeout
active, reopen-with-existing-database works, monotonic per-Run cursors (independent across Runs, batch append
continues the sequence), clear missing-Run error, cursor paging with limits, byte-exact payload round-trip with
normalized columns preserved, a read-only connection paging contiguous cursors while a writer goroutine appends
(no sleeps; read-only writes rejected), and data_version changing after a writer commit. Verification:
`rtk go vet` clean, `rtk go test -race ./internal/store/` 15 passed, full `rtk go test ./...` 173 passed in 15
packages, `rtk go run ./cmd/roundfix --help` green.
