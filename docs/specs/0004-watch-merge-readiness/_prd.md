---
spec: 0004-watch-merge-readiness
status: active
created: 2026-07-05
surfaces: [cli]
---

# Watch Merge Readiness

The review dogfood proved the watch loop works end to end — and exposed where
it wastes time and stops short. It sleeps before its first useful check even
when the review finished long ago; it declares Clean the moment the local
Review Issue queue empties, while the pull request still shows the Review
Source's status check re-reviewing the just-pushed commit; its stdout is
empty, leaving the exit code as the only machine-readable result; and its
stderr buries the Daemon's milestones under the full Agent console. This Spec
makes watch end exactly when the pull request is truly ready for the
developer's merge decision, and makes both ends of the pipe honest.

## Goals

- A watch Run against an already-reviewed pull request starts useful work
  immediately — no timer runs ahead of the first status check.
- A watch Run that ends Clean means the pull request is merge-ready: the
  Review Source's status check on the final pushed commit succeeded with no
  new Review Issues. See ADR-0019.
- Watch and resolve report deterministically on stdout, shaped like the
  Implement Command's report.
- A supervising process can follow a Run's Daemon milestones on stderr
  without filtering out Agent console noise.

## User Stories

1. As a developer starting watch on a pull request whose review finished long
   ago, I want the first Review Source status check to happen immediately and
   the quiet period to be skipped when the review was already settled, so
   that resolution starts in seconds, not after dead timers.
2. As a developer relying on `--until-clean`, I want the Run to keep watching
   after the Final Push until the Review Source's status check on the pushed
   commit reports success with no new Review Issues, so that Clean means
   ready for squash/merge, not just an empty local queue.
3. As a developer or script consuming watch and resolve output, I want one
   deterministic stdout report — per-Review-Issue lines plus one outcome
   line — so that the result is machine-readable beyond the exit code.
4. As a developer supervising a Run from another process, I want a flag that
   silences Agent console events on stderr while keeping Daemon milestones,
   so that following progress does not require fragile pattern filtering.

## Core Features

1. **Poll first, sleep after.** The watch loop performs its first Review
   Source status check immediately on start; `poll_interval` applies only
   between subsequent checks. The quiet period is skipped when the review was
   already settled before the Run began.
2. **Merge-readiness termination.** With `--until-clean`, after a Final Push
   the Run polls the Review Source's status check on the pushed head commit
   (bounded by the existing Max Rounds and review timeout): success with no
   new Review Issues ends the Run Clean; a failing or re-reviewing check that
   yields new Review Issues starts the next Round; the existing bounds and
   their outcomes are unchanged. See ADR-0019.
3. **Deterministic stdout for watch and resolve.** One line per Review Issue
   in Round/fetch order — identifier, final status, title — then one outcome
   line naming the terminal state, Rounds used, and counts. Diagnostics stay
   on stderr; exit codes are unchanged.
4. **Agent console suppression.** A flag on the operational commands
   (`--no-agent-console`) drops Agent-source events from the stderr stream,
   keeping every Daemon event; the Run Event Journal is untouched — this is a
   display filter only.

## User Experience

Same commands, same flags plus one (`--no-agent-console`), same exit codes.
The visible changes: watch starts acting immediately when there is work
waiting; watch may keep running briefly after the Final Push while the head
check settles (bounded); stdout now carries the report scripts can consume.

## Non-Goals / Out of Scope

- New Review Sources or webhook triggers (work-plan item 3).
- TUI changes (next Spec).
- Retry budgets or escalation (work-plan item 7).
- Changing Round semantics, Max Rounds, the Run Budget, or Final Push
  mechanics.

## Success Metrics

- A watch dogfood against a pre-reviewed pull request reaches its first fetch
  in under one poll interval from start.
- A watch Run that ends Clean shows the Review Source check green on the
  final head commit at the moment of completion (validated live in the next
  review dogfood).
- Watch and resolve stdout reports are byte-deterministic for equal inputs
  and non-empty for every terminal outcome.
- A supervisor using `--no-agent-console` sees only Daemon-source lines on
  stderr for a full Run.

## Decisions

- Merge-readiness is the default `--until-clean` semantics, not an opt-in —
  the old behavior was a weaker promise of the same intent. See ADR-0019.
- The check source reads the GitHub status check / check run for the pull
  request's head commit through the existing GitHub CLI integration; no new
  authentication surface.
- Check polling reuses `watch.poll_interval` and the review timeout; no new
  config keys.
- The suppression flag is display-only; the Journal always records everything
  (ADR-0008 untouched).

## Open Questions

None.
