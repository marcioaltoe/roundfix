---
title: "Attach command: replay completed Runs through the timeline renderer"
type: HITL
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
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

- [ ] Attach by Run ID replays events after cursor zero through the shared renderer
- [ ] Terminal Run renders read-only and exits cleanly
- [ ] Unknown Run ID / missing journal fail as CLI errors before TUI start
- [ ] Attach creates no Runs, fetches nothing, starts no Agents, commits, pushes, or resolves threads
- [ ] Unknown event kinds skipped on replay
- [ ] Bounded console memory and chunk coalescing proven by tests with fake event sources
- [ ] Design review of the rendered timeline completed

## Blocked by

- 06-journal-critical-sink-in-resolve.md
