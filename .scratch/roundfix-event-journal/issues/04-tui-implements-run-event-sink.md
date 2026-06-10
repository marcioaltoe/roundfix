---
title: "Live Run View implements the Run Event sink"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - completed
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

- [x] Live Run View receives Agent output as Run Events through the sink interface
- [x] Bounded buffer with counted drops; publication never blocks on a slow UI
- [x] No remaining references to the removed stream-sink interface
- [x] TUI tests prove buffer bounds and drop counting without a real terminal
- [x] Race tests for the non-blocking adapter

## Blocked by

- 03-runevent-seam-runner-and-writersink.md

## Comments

**2026-06-10 (agent):** Added `tui.EventBuffer` (`internal/tui/event_sink.go`): a best-effort `runevent.Sink`
with a bounded channel (512 default, injectable capacity), one delivery goroutine, non-blocking enqueue with
atomic drop counting, and drain-on-Close. `AgentLiveStream` now implements `Publish` through the buffer and
exposes `DroppedEvents()`; delivery goes to the unexported `handleAgentUpdate` (TTY → Bubble Tea program,
non-TTY → console text). The CLI's interim `liveViewSink` from issue 03 is deleted — `newAgentEventSink` returns
the live view itself when live, `agent.WriterSink` otherwise, so live streaming behavior is unchanged. Grep
confirms zero references to the removed `StreamUpdateSink`/`HandleAgentUpdate`/`publishStreamUpdate`. New tests
prove ordered delivery, never-blocking publication with counted drops under a blocked renderer, unknown kinds
skipped without drop counts, post-Close drops, and delivered+dropped accounting under 8 concurrent publishers —
all without a terminal and run under `-race`. Verification: `rtk go vet ./...` clean,
`rtk go test -race` on tui/cli/agent/runevent 97 passed, full `rtk go test ./...` 166 passed in 15 packages,
`rtk go run ./cmd/roundfix --help` green.
