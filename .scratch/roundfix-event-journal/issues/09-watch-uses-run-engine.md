---
title: "Watch loop delegates resolve cycles to the Run engine"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - completed
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

- [x] Watch Rounds execute through the same engine operation as the resolve command
- [x] Duplicated orchestration wiring between resolve and watch removed
- [x] Watch-side Final Push gating preserved and proven by tests
- [x] Full watch loop test without network passes through the engine
- [x] Watch state machine module keeps owning polling/quiet/timeout/round policy

## Blocked by

- 07-daemon-run-engine-resolve-cycle.md
- 08-snapshot-diff-batch-commits.md

## Comments

**2026-06-10 (agent):** The structural delegation landed with issue 07: `resolveWatchBatches` calls the same
`executeResolveCycle` → `daemon.Engine.ResolveCycle` path as the resolve command, and the old duplicated
`executeResolveBatches`/`executeResolveBatch` wiring is gone (grep across the repo finds no references). The watch
state machine module (`internal/watch`) is untouched in this unit — it keeps owning polling, quiet period, review
timeout, and round policy, and calls the engine once per Round through its Resolver callback. This unit adds the
remaining proof: `TestWatchFinalPushRunsOncePerCleanRoundThroughEngine` (full watch loop with fake collaborators,
no network: fetch → engine-driven Batch commit → exactly one Final Push to the push-plan target on the clean Round)
and `TestWatchSkipsFinalPushWhenAutoPushDisabled` (auto-push off → zero pushes, "Final Push skipped" gating message,
Run still reaches Clean). The pre-existing watch table test also exercises the whole loop through the engine fakes.
Verification: `rtk go vet ./...` clean, `rtk go test ./...` 189 passed in 15 packages,
`rtk go run ./cmd/roundfix --help` green.
