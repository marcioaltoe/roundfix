---
title: "Attach command: replay completed Runs through the timeline renderer"
type: HITL
category: enhancement
state: completed
labels:
  - enhancement
  - completed
user_stories:
  - 2
  - 3
  - 6
  - 9
  - 10
  - 12
  - 13
  - 14
  - 16
blocked_by:
  - 06-journal-critical-sink-in-resolve.md
---

# Attach command: replay completed Runs through the timeline renderer

## Parent

.scratch/roundfix-event-journal/05-attach-replay-live-run-view-prd.md

## What to build

First half of the attach/replay PRD: a CLI attach surface that loads the Run snapshot and replays Run Events from
the journal through one timeline renderer consuming Run Events only. Agent payloads decode from the raw ACP JSON per
ADR 0008; unknown event kinds are skipped, never fatal. Terminal Runs render read-only without waiting. Missing Run,
missing journal, and database errors surface as CLI errors before the TUI starts. Attach is non-mutating. Console
memory is bounded by a ring buffer inside the renderer; high-frequency message chunks coalesce before rendering.

HITL checkpoint: this is where Claude Code-style rendering fidelity (tool markers, diffs, durations) takes visual
shape — request a design review of the rendered timeline before the follow-mode issue builds on it.

## Acceptance criteria

- [x] Attach by Run ID replays events after cursor zero through the shared renderer
- [x] Terminal Run renders read-only and exits cleanly
- [x] Unknown Run ID / missing journal fail as CLI errors before TUI start
- [x] Attach creates no Runs, fetches nothing, starts no Agents, commits, pushes, or resolves threads
- [x] Unknown event kinds skipped on replay
- [x] Bounded console memory and chunk coalescing proven by tests with fake event sources
- [x] Design review of the rendered timeline completed

## Blocked by

- 06-journal-critical-sink-in-resolve.md

## Comments

**2026-06-10 (agent):** Added `roundfix attach <run-id>` (`internal/cli/attach.go`): parses a positional Run ID or
`--run-id`, opens the Run Database with `store.OpenReader` (read-only — attach is structurally non-mutating), loads
the Run snapshot, and pages `RunEventsAfter` from cursor zero (200/page) into the new shared renderer
`tui.RunTimeline` (`internal/tui/timeline.go`): Run Events in, console lines out, ring-bounded (300 lines for
attach), message chunks coalesced into whole lines via the StreamBuffer pending-line mechanics, unknown
kinds/undecodable payloads skipped. Output renders the existing Live Run View layout (`RenderLiveRunView`) with the
Run snapshot, Review Issues in the left pane (artifact lookup failures degrade to an empty pane, never a failure),
and the replayed timeline in the console pane; terminal Runs print "reached <state>; timeline replayed read-only."
and exit 0 without waiting. Missing Run ID, unknown Run, and missing/broken Run Database surface as
`roundfix attach failed: …` CLI errors with exit 2 before any rendering. Usage/help updated. Tests: RunTimeline
coalescing (3 chunks → 1 line), ring bound (50 → 5 newest), unknown-kind skip, tool-marker rendering from raw
payloads; CLI tests prove end-to-end replay of a journaled fake Resolve Run, non-mutation (run count and terminal
state untouched, agent prober poisoned to fail loudly if attach ever probed), unknown-Run and missing-database CLI
errors, and unknown event kinds skipped on replay against a real journal. HITL design review: rendered timeline
presented (tool markers, plan, diff/terminal markers, coalesced summary) and **approved as-is** by the user as the
baseline for follow mode. Verification: `rtk go vet ./...` clean, `rtk go test ./...` 197 passed in 15 packages,
`rtk go run ./cmd/roundfix --help` and `attach --help` green.
