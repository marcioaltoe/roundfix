---
title: "Stop: terminate Runs without hiding unfinished work"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 45
  - 46
  - 47
  - 48
  - 49
  - 50
blocked_by:
  - 08-agent-runtime-batch-resolution.md
  - 11-watch-review-round-loop.md
---

# Stop: terminate Runs without hiding unfinished work

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

Implement Stop Request and SIGINT behavior for Active Runs. Stopping must be a
terminal path that preserves unfinished Agent work, releases the Active Run
lock, and avoids any later verification, commit, push, fetch, or Review Source
mutation.

## Acceptance criteria

- [ ] A Stop Request with no active Agent marks the Run `Stopped`, releases the
      Active Run lock, and exits cleanly.
- [ ] A Stop Request with an active Agent sends graceful termination, waits up to
      the configured grace period, then kills the process if needed.
- [ ] Stop handling persists available Agent output and reports changed paths and
      statuses left in the worktree.
- [ ] Stop handling does not run verification, commit, push, fetch more issues,
      or resolve source threads after the Stop Request.
- [ ] A later `resolve` or `watch` creates a new Run, and dirty preserved work
      blocks that new Run during Preflight Validation.
- [ ] Roundfix-controlled clean stops exit `0`; SIGINT exits `130`.

## Blocked by

- 08-agent-runtime-batch-resolution.md
- 11-watch-review-round-loop.md
