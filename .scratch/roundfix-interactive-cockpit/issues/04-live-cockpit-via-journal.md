---
title: "Live resolve/watch cockpit reads the journal (ADR 0009)"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
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

- [ ] Live resolve/watch TTY cockpit renders journal content while the engine writes (same-process writer+reader)
- [ ] Ctrl-C in the owning cockpit still produces a Stop Request with today's stop semantics and exit codes
- [ ] No detach in owning mode; contextual footer differs between owning and attach modes
- [ ] TUI sink delivery retired; non-TTY console output byte-compatible with today across resolve and watch
- [ ] One cockpit per Watch Run across all Rounds and Batches
- [ ] Race tests: owning process writing while its own cockpit polls; stop during active polling

## Blocked by

- 03-interactive-cockpit-attach.md
