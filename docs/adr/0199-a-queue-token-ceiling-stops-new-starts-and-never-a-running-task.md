---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A queue's token ceiling stops new starts and never a running Task

A Delivery Queue could be bounded by a deadline and a retry allowance, both
measured in time or attempts. Once ADR-0198 records the tokens each Run used,
an operator can bound a queue by what it consumed as well. Stopping a Run
when the ceiling is crossed would discard a Task's work in the middle of its
Agent Session, and a Task can spend millions of tokens in one turn, so no
check between prompts could keep the total under the ceiling anyway.

The token ceiling is a Delivery Queue Limit with the deadline's boundary:

- **What it counts.** The tokens recorded for the queue's Runs: every Run the
  queue recorded for any of its items, including Runs of earlier retries.
  Unreported prompts add nothing, and the status line says how many there are.
- **Where it acts.** Before a queued item starts and when a parked item is
  retried. At or above the ceiling, a queued item parks as
  `queue-token-ceiling` without a worktree, and a retry is refused. An item
  past `queued` continues, and a Run in progress is never stopped: its Tasks
  settle and it ends on its own.
- **How it is answered.** Only a new queue, recorded with a higher ceiling or
  none. The Pending Question says so.

## Consequences

- A queue can end above its ceiling by up to one item's Run. The ceiling
  bounds how many items start, not the exact total.
- An adapter that reports nothing can never trip the ceiling.
- Tokens outside Runs, such as the pre-PR review, do not count.
