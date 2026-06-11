---
title: "Live resolve/watch cockpit reads the journal (ADR 0009)"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - completed
user_stories:
  - 11
  - 12
  - 13
  - 14
  - 16
blocked_by:
  - 03-interactive-cockpit-attach.md
---

# Live resolve/watch cockpit reads the journal (ADR 0009)

## Parent

.scratch/roundfix-interactive-cockpit/01-interactive-live-run-view-prd.md

## What to build

Switch the owning commands to the shared cockpit: in a TTY, resolve and watch open the same interactive cockpit as
attach, consuming the Run Event Journal over a read-only connection while the engine writes — never the live sink,
per ADR 0009. The TUI event-buffer sink delivery is retired; the fanout keeps the journal as critical sink and the
non-TTY writer path unchanged. Owning-mode keys: Ctrl-C remains Stop Request exactly as today; there is no detach
(q does nothing); the contextual footer says so. Watch runs one cockpit for the entire Run instead of constructing
a view per cycle.

## Acceptance criteria

- [x] Live resolve/watch TTY cockpit renders journal content while the engine writes (same-process writer+reader)
- [x] Ctrl-C in the owning cockpit still produces a Stop Request with today's stop semantics and exit codes
- [x] No detach in owning mode; contextual footer differs between owning and attach modes
- [x] TUI sink delivery retired; non-TTY console output byte-compatible with today across resolve and watch
- [x] One cockpit per Watch Run across all Rounds and Batches
- [x] Race tests: owning process writing while its own cockpit polls; stop during active polling

## Blocked by

- 03-interactive-cockpit-attach.md

## Comments

**2026-06-10 (agent):** Added `runUI` (`internal/cli/runui.go`), the per-command Run Event consumer wiring. TTY:
opens its own read-only journal connection (`store.OpenReader`), starts the shared cockpit in owning mode in a
goroutine, the engine's fanout becomes journal-critical only, and raw progress text goes to `io.Discard` — the
cockpit narrates the loop from the journaled daemon/agent events (ADR 0009). Non-TTY: journal critical +
`agent.WriterSink` over stderr, progress to stderr — byte-compatible, with every pre-existing CLI test passing
unchanged. Ctrl-C in the owning cockpit (a key event under raw mode, not SIGINT) calls `interruptSelf`, which
re-sends `os.Interrupt` to the process so the existing `interruptContext`/`exitForInterrupt` plumbing preserves
Stop Request semantics and exit codes by construction; q does nothing in owning mode and the contextual footer
shows `Ctrl-C stop` (vs `q detach` in attach). `executeResolveCycle` now receives the command-level `runUI`
instead of building a per-cycle view/stream/fanout, so the watch command runs **one cockpit for the entire Run**
across all Rounds and Batches (watch status/fetch prints also route through `ui.progress`). The TUI sink delivery
is fully retired: `EventBuffer` (event_sink.go), `AgentLiveStream`, the live Bubble Tea model, and their CLI
helpers are deleted along with the EventBuffer tests; render helpers the cockpit uses stay. Race coverage:
`TestOwningCockpitPollsJournalWhileOwnProcessWrites` runs a same-process writer goroutine (25 events +
CompleteRun) while the owning cockpit model polls its read-only connection, fires Ctrl-C mid-poll (Stop Request
callback observed), and renders every line plus the terminal READ-ONLY state — under `-race`. Verification:
`rtk go vet ./...` clean, `rtk go test ./...` 222 passed in 15 packages, `rtk go test -race` cli+tui 108 passed,
`rtk go run ./cmd/roundfix --help` green.
