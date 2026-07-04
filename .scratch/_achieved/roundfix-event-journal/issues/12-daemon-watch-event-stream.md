---
title: "Daemon and Watch Run event stream: full loop coverage"
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

- [x] Watch Run outcomes provable through ordered journal event sequences (reviewing, settled, fetch, resolve, terminal)
- [x] Stop events appended; later unsafe daemon actions absent after Stop Request
- [x] Verification pass/fail and commit/commit-skip decisions present in the journal
- [x] Final Push events occur only when no Unresolved Review Issues remain
- [x] Source-resolution events daemon-owned and absent after Stop Request
- [x] Watch Run daemon events render in the attach Live Run View timeline
- [x] Race tests cover watch polling, Agent ownership, and attach/follow concurrency

## Blocked by

- 09-watch-uses-run-engine.md
- 11-attach-follow-mode-polling.md

## Comments

**2026-06-10 (agent):** Completed the daemon event taxonomy in runevent (`daemon.review_status`, `quiet_period`,
`fetch`, `selection`, `batch`, `verification`, `commit`, `push`, `source_resolution`, `retry` (reserved), `outcome`,
plus the existing `daemon.status`; `IsDaemonKind` gates reader rendering so unknown kinds stay skippable). The Run
engine publishes through its sink dependency at every transition — selection at cycle start, batch
started/failed/completed, verification started/passed/failed (failure event appended before the Run fails), commit
created/skipped/disabled, source resolution (daemon-owned), and Final Push pushed — with small ID/count payloads
and bounded summaries; publication is non-optional, so a critical journal failure fails the cycle. The watch state
machine gained a Sink dependency and RunID and publishes review-status polls, quiet-period waits, and fetch
started/completed. The CLI journals Final Push gating decisions (blocked/skipped) and a terminal `daemon.outcome`
event after every CompleteRun (resolve clean/failed/stopped and all watch outcomes). `RunTimeline` renders daemon
events from their summaries, so attach/replay narrates the whole loop with no Watch special-casing; the
resolve/watch console text is unchanged (engine progress writes untouched; WriterSink stays agent-only). Tests:
ordered subsequence proof for a clean Watch Run (review_status → fetch×2 → selection → batch → verification×2 →
commit → source_resolution → batch → push → outcome, exactly one pushed decision, Clean outcome last), stop run
(stop event journaled, zero unsafe daemon events after it, Stopped outcome last), failed verification (failure
event, no commit/push events, Failed outcome last), triage-only commit-skip decision journaled, and attach
rendering the Watch narrative end to end. Verification: `rtk go vet ./...` clean, `rtk go test ./...` 206 passed
in 15 packages, `rtk go test -race` across watch/cli/daemon/tui 117 passed, `rtk go run ./cmd/roundfix --help`
green.
