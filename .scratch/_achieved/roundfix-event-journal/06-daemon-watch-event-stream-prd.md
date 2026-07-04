---
title: "PRD: daemon and Watch Run event stream"
type: PRD
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
blocked_by:
  - 01-run-event-seam-prd.md
  - 02-run-event-journal-prd.md
  - 03-journaled-agent-stream-prd.md
  - 04-daemon-run-engine-prd.md
  - 05-attach-replay-live-run-view-prd.md
---

# PRD: daemon and Watch Run event stream

## Problem Statement

After Agent stream events are journaled and attach/replay works, Watch Runs still need complete daemon-level event coverage. A Watch Run is more than Agent output: it waits for CodeRabbit status, fetches Rounds, resolves Batches, verifies, commits, pushes, handles Stop Requests, and terminates on clean, timed out, budget-limited, stopped, failed, or max-rounds outcomes.

The user needs the durable cockpit to explain the entire review-resolution loop, not just the child Agent stream.

## Solution

Make daemon and Watch Run orchestration publish structured Run Events for every meaningful state transition and decision. The Run Event Journal should become the timeline for the full Roundfix loop: Review Source waits, quiet period, fetch results, compatible artifact decisions, Batch assignment, Agent execution, verification, commit, Final Push, source resolution, Stop Request, retry, timeout, and terminal outcome.

The attach/replay Live Run View should show a complete Run narrative from the same durable event stream.

## User Stories

1. As a developer, I want Watch Run status waits to be journaled, so that I know whether Roundfix is waiting on CodeRabbit or local work.
2. As a developer, I want quiet period waits to be journaled, so that apparent idle time is explained.
3. As a developer, I want fetch start and fetch result events, so that each Round's downloaded issue count is visible.
4. As a developer, I want compatible artifact selection events, so that Resolve scope is explainable.
5. As a developer, I want deduplication and Batch assignment events, so that duplicate Review Issues are transparent.
6. As a developer, I want Agent start and completion events per Batch, so that the cockpit shows progress across Batches.
7. As a developer, I want verification start, pass, and fail events, so that Batch outcomes are auditable.
8. As a developer, I want commit decisions journaled, so that I know when the daemon created a Batch commit or skipped one.
9. As a developer, I want Final Push decisions journaled, so that I know why Roundfix pushed or refused to push.
10. As a developer, I want Review Source resolution events journaled, so that GitHub mutations are traceable.
11. As a developer, I want Stop Request events journaled, so that interrupted Runs explain what was skipped and what local work was preserved.
12. As a developer, I want timeout and budget events journaled, so that terminal outcomes are clear.
13. As a developer, I want max-rounds terminal events journaled, so that the configured review policy is visible.
14. As a developer, I want attach/replay to show the whole Watch loop, so that I do not need to read separate logs.
15. As a developer, I want daemon event payloads to be concise, so that long watch histories remain navigable.
16. As a developer, I want errors to include next actions, so that a failed Watch Run is actionable from the cockpit.
17. As a daemon, I want event publishing to be non-optional for state transitions, so that the Run Database remains the operational source of truth.
18. As a test author, I want the watch loop's terminal outcomes to be provable through journal events, so that UI confidence does not depend on visual snapshots alone.

## Implementation Decisions

- Treat daemon event publication as part of the Run state contract. Any user-meaningful state transition should append a Run Event.
- Use the Run Event Journal for Watch Run events and Resolve Run daemon events. Do not introduce a second event system.
- Keep event payloads structured but small. Large Agent output belongs to Agent stream events; daemon events should summarize decisions and include IDs/counts.
- Define a stable set of daemon event kinds for review status, fetch, resolve, batch, verification, commit, commit skip, push, source resolution, stop, retry, timeout, budget, and terminal outcome, using the namespaced kind vocabulary from the Run Event seam.
- Keep Final Push and Review Source mutation events daemon-owned. Child Agents must never publish those as if they performed the mutation.
- Attach/replay should consume these daemon events without special casing Watch Runs outside the existing Live Run View state model.
- If a daemon operation fails after Run start, append a failure event before completing the Run as failed when possible.
- If context cancellation or Stop Request occurs, append stop events and then skip later verification, commits, pushes, fetches, and Review Source mutation.
- Do not add a remote daemon process protocol in this PRD unless it is already necessary for existing command execution. The durable Run Database remains the primary interface.
- Publish daemon events through the Run engine's Run Event sink dependency, so publication is enforced at the engine seam introduced by the Daemon Run engine PRD and testable with a fake sink. Watch-loop waits and round policy events publish through the same sink from the watch state machine.

## Testing Decisions

- Good tests should prove user-visible event sequences for Watch Run outcomes, not private helper calls.
- Watch loop tests should assert that reviewing, settled, fetch, resolve, and terminal outcome events are appended in order.
- Stop tests should assert that stop events are appended and later unsafe daemon actions are absent.
- Verification/commit tests should assert pass/fail and commit-skip decisions are present in the journal.
- Push tests should assert Final Push events occur only when no Unresolved Review Issues remain.
- Review Source resolution tests should assert source-resolution events are daemon-owned and absent after Stop Request.
- Attach/replay tests should verify Watch Run daemon events render in the Live Run View timeline.
- Race tests are required because Watch polling, Agent process ownership, and attach/follow behavior involve concurrency.

## Out of Scope

- Changing Review Source fetch semantics.
- Changing deduplication rules.
- Changing Final Push policy.
- Child Agents committing, pushing, or resolving Review Source threads.
- A remote multi-client daemon socket.
- Event analytics, metrics dashboards, or long-term retention policy.

## Further Notes

This is the integration slice that makes the durable cockpit explain the whole Roundfix loop. It should come after the journal, Agent stream persistence, and attach/replay foundations are already working.
