---
title: "Timeline window and Follow Mode state machine"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 1
  - 2
  - 3
  - 4
  - 5
  - 14
  - 15
  - 16
blocked_by:
  - 01-journal-backward-paging.md
---

# Timeline window and Follow Mode state machine

## Parent

.scratch/roundfix-interactive-cockpit/01-interactive-live-run-view-prd.md

## What to build

Evolve `RunTimeline` from a tail-only line ring into the cockpit's viewport engine: a bounded sliding window of
~500 events (internal constant) rendered to lines, scrolled by line, paged by event when the viewport hits a window
edge — backward via the new store read, forward via the existing one. On top of it, the Follow Mode state machine:
REPLAYING while loading backlog, FOLLOWING at the bottom (tail auto-advances), SCROLLED when the user scrolled up
(viewport frozen, count of new events below, End/G resumes), and a terminal read-only state that never follows.
Scrolling never affects the Run — the reader merely stops paging forward. Everything provable with fake event
sources, no terminal.

## Acceptance criteria

- [ ] Sliding window bounded at the internal constant; paging evicts from the opposite edge; whole journal reachable both ways
- [ ] No duplicate or missing lines across window slides (event→lines rendering stable at boundaries)
- [ ] FOLLOWING auto-advances the tail; SCROLLED freezes the viewport and counts new events below
- [ ] End/G (and manually reaching the bottom) resumes FOLLOWING; terminal Runs never enter FOLLOWING
- [ ] State + new-event count exposed for the status bar (REPLAYING / FOLLOWING / SCROLLED · N below / terminal state)
- [ ] All behavior proven with fake event sources and injected pacing; no sleeps; race tests where goroutines exist

## Blocked by

- 01-journal-backward-paging.md
