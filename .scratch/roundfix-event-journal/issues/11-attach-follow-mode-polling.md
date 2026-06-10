---
title: "Attach follow mode: cursor polling with change detection"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
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

- [ ] Follow mode appends only events newer than the replay cursor; no duplicates across reconnects
- [ ] Idle polls detect no-change without reading event rows
- [ ] Read-only connection; no long-lived read transactions
- [ ] Backlog vs live status visible in the UI
- [ ] Detach leaves the Run active; cancellation tests prove it
- [ ] Polling tests use injected clocks or fake sources; no sleeps; race tests for follow goroutines

## Blocked by

- 10-attach-replay-timeline-renderer.md
