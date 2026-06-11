---
title: "Timeline window and Follow Mode state machine"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - completed
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

- [x] Sliding window bounded at the internal constant; paging evicts from the opposite edge; whole journal reachable both ways
- [x] No duplicate or missing lines across window slides (event→lines rendering stable at boundaries)
- [x] FOLLOWING auto-advances the tail; SCROLLED freezes the viewport and counts new events below
- [x] End/G (and manually reaching the bottom) resumes FOLLOWING; terminal Runs never enter FOLLOWING
- [x] State + new-event count exposed for the status bar (REPLAYING / FOLLOWING / SCROLLED · N below / terminal state)
- [x] All behavior proven with fake event sources and injected pacing; no sleeps; race tests where goroutines exist

## Blocked by

- 01-journal-backward-paging.md

## Comments

**2026-06-10 (agent):** Built `tui.TimelineViewport` (`internal/tui/viewport.go`): a synchronous sliding window of
journal events (defaults 500 events / 100 per page, injectable bounds for tests) over a `TimelineSource` interface
satisfied by `*store.Store`. Replay loads the newest window in one `RunEventsBefore(MaxInt64)` read; entries keep
rendered *text* per event and the window flattens text before line-splitting, so partial message chunks coalesce
across event boundaries exactly like the streaming renderer — proven by a chunks-split-across-three-events test.
Scroll is line-based; hitting the top edge pages backward (prepend shifts the scroll by exactly the lines added,
then bottom eviction — which never moves top offsets — sheds entries it can't still see), hitting the bottom edge
pages forward with head eviction and delta-exact scroll correction. Follow Mode states: REPLAYING on construction;
FOLLOWING pins the tail and appends on Poll with the window bound enforced; SCROLLED freezes the viewport — Poll
only counts arrivals beyond the window tail (`State()` exposes state + N-below for the status bar); End/G
(`JumpToTail`) reloads the newest window and resumes; manually scrolling to the bottom at the journal tail also
resumes; `SetTerminal` Runs never enter FOLLOWING and keep full scrollback. Walk tests prove every one of 120
events becomes visible in both a tail→head and a head→tail traversal with a 30-event window (the backward-paging
scroll-adjustment bug this caught was fixed by separating the prepend shift from bottom eviction); unknown kinds
render nothing, without blank lines. The viewport owns no goroutines (driver polls), so the race criterion is
satisfied structurally — the package still runs green under `-race`. Verification: `rtk go vet ./...` clean,
`rtk go test -race ./internal/tui/` 22 passed, full `rtk go test ./...` 218 passed in 15 packages,
`rtk go run ./cmd/roundfix --help` green.
