---
title: "Watch loop delegates resolve cycles to the Run engine"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 5
  - 13
blocked_by:
  - 07-daemon-run-engine-resolve-cycle.md
  - 08-snapshot-diff-batch-commits.md
---

# Watch loop delegates resolve cycles to the Run engine

## Parent

.scratch/roundfix-event-journal/04-daemon-run-engine-prd.md

## What to build

Completes the Daemon Run engine PRD: the watch state machine keeps polling, quiet period, review timeout, and round
policy in its own module and calls the engine once per Round, eliminating the duplicated verify/commit/push wiring
between resolve and watch. Final Push gating in watch remains: only when no Unresolved Review Issues remain across
the Run and auto-push is enabled.

## Acceptance criteria

- [ ] Watch Rounds execute through the same engine operation as the resolve command
- [ ] Duplicated orchestration wiring between resolve and watch removed
- [ ] Watch-side Final Push gating preserved and proven by tests
- [ ] Full watch loop test without network passes through the engine
- [ ] Watch state machine module keeps owning polling/quiet/timeout/round policy

## Blocked by

- 07-daemon-run-engine-resolve-cycle.md
- 08-snapshot-diff-batch-commits.md
