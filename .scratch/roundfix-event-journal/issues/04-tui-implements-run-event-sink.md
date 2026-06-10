---
title: "Live Run View implements the Run Event sink"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 8
  - 14
  - 15
blocked_by:
  - 03-runevent-seam-runner-and-writersink.md
---

# Live Run View implements the Run Event sink

## Parent

.scratch/roundfix-event-journal/01-run-event-seam-prd.md

## What to build

Completes the Run Event seam PRD on the TUI side: the Live Run View consumes Run Events by implementing the sink
interface as a best-effort adapter with a bounded buffer and counted drops, so rendering pressure never blocks or
fails producers. Live streaming behavior in a TTY must not regress.

## Acceptance criteria

- [ ] Live Run View receives Agent output as Run Events through the sink interface
- [ ] Bounded buffer with counted drops; publication never blocks on a slow UI
- [ ] No remaining references to the removed stream-sink interface
- [ ] TUI tests prove buffer bounds and drop counting without a real terminal
- [ ] Race tests for the non-blocking adapter

## Blocked by

- 03-runevent-seam-runner-and-writersink.md
