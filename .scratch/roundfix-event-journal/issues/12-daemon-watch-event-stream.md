---
title: "Daemon and Watch Run event stream: full loop coverage"
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
  - 6
  - 7
  - 8
  - 9
  - 10
  - 11
  - 12
  - 13
  - 14
  - 15
  - 16
  - 17
  - 18
blocked_by:
  - 09-watch-uses-run-engine.md
  - 11-attach-follow-mode-polling.md
---

# Daemon and Watch Run event stream: full loop coverage

## Parent

.scratch/roundfix-event-journal/06-daemon-watch-event-stream-prd.md

## What to build

The daemon and Watch Run event stream PRD: every user-meaningful state transition publishes a Run Event through the
engine's sink and the watch state machine — review status waits, quiet period, fetch start/result, compatible
artifact selection, deduplication and Batch assignment, Agent start/completion per Batch, verification pass/fail,
commit and commit-skip decisions, Final Push decisions, Review Source resolution, Stop Request, retry, timeout,
budget, and terminal outcomes. Daemon kinds use the namespaced vocabulary from the Run Event seam; payloads stay
small (IDs and counts). Attach/replay renders the full Watch loop narrative from the same journal.

## Acceptance criteria

- [ ] Watch Run outcomes provable through ordered journal event sequences (reviewing, settled, fetch, resolve, terminal)
- [ ] Stop events appended; later unsafe daemon actions absent after Stop Request
- [ ] Verification pass/fail and commit/commit-skip decisions present in the journal
- [ ] Final Push events occur only when no Unresolved Review Issues remain
- [ ] Source-resolution events daemon-owned and absent after Stop Request
- [ ] Watch Run daemon events render in the attach Live Run View timeline
- [ ] Race tests cover watch polling, Agent ownership, and attach/follow concurrency

## Blocked by

- 09-watch-uses-run-engine.md
- 11-attach-follow-mode-polling.md
