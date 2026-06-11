---
title: "Attach follow mode: cursor polling with change detection"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - completed
user_stories:
  - 1
  - 4
  - 5
  - 7
  - 8
  - 11
  - 15
blocked_by:
  - 10-attach-replay-timeline-renderer.md
---

# Attach follow mode: cursor polling with change detection

## Parent

.scratch/roundfix-event-journal/05-attach-replay-live-run-view-prd.md

## What to build

Second half of the attach/replay PRD: after replaying backlog, the attach view follows newly appended events for
Active Runs. Polling uses the journal change-detection signal so idle polls read no rows, over a read-only database
connection with short autocommit reads. The UI advances its cursor only after accepting events, so reconnects never
duplicate output. The view shows replaying-backlog versus following-live status. Detach exits the UI and leaves the
Run active; Stop Request remains a separate explicit action.

## Acceptance criteria

- [x] Follow mode appends only events newer than the replay cursor; no duplicates across reconnects
- [x] Idle polls detect no-change without reading event rows
- [x] Read-only connection; no long-lived read transactions
- [x] Backlog vs live status visible in the UI
- [x] Detach leaves the Run active; cancellation tests prove it
- [x] Polling tests use injected clocks or fake sources; no sleeps; race tests for follow goroutines

## Blocked by

- 10-attach-replay-timeline-renderer.md

## Comments

**2026-06-10 (agent):** Added `attachFollower` to `internal/cli/attach.go`: after backlog replay, non-terminal Runs
enter follow mode against the same read-only Run Database connection. Each poll reads `PRAGMA data_version` first
and only drains event rows when the version changed (idle polls read zero rows); every read is a short autocommit
query — no long-lived read transaction exists. The cursor advances only after an event is accepted, so reconnects
from the returned cursor never duplicate output (proven by a reconnect test). Events stream through the same
`RunTimeline` renderer (its `Append` now returns the rendered text so replay and follow share one renderer).
Status is explicit in the output: "Replayed backlog through cursor N…" then "Following live events. Detach with
Ctrl-C; detaching never stops the Run." A terminal transition drains a final page and exits with
"Run … reached <state>."; cancellation exits with "Detached; Run … keeps going." Pacing goes through the
`attachSleep` seam (default 250ms timer honoring ctx); tests inject immediate sleepers — no real sleeps anywhere.
Tests: fake-source follower tests (newer-than-cursor only + reconnect no-duplicates; idle change-detection with
exactly one row read across five idle polls), CLI-level detach-leaves-Run-active (cancel on first poll, Run still
non-terminal afterward), and a live end-to-end follow where a writer goroutine appends five events and completes
the Run while attach follows to Stopped with each line rendered exactly once — run under `-race`. Verification:
`rtk go vet ./...` clean, `rtk go test ./...` 201 passed in 15 packages, `rtk go test -race ./internal/cli/` 77
passed, `rtk go run ./cmd/roundfix --help` green.
