---
status: accepted
created_at: 2026-11-13T12:00:00Z
updated_at: 2026-11-13T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A parked item always has a way back through a Delivery Retry

A Delivery Retry returns an item parked before it has an Item Branch to
`queued`, counting a retry, and the owner then creates the branch as on a
first start. The delivery authorization of an archived Spec is read from Git
history without requiring the Spec Root in the working tree. At `checking`, a
Pull Request whose head is an earlier candidate of the item sends the item back
to `gating`, so the repository gate runs again and the candidate is pushed; any
other head still parks `review-stale`, whose next action now says what to do.

## Context

On 2026-11-13 fiscus reported an item parked `delivery-error: read
prerequisites: refresh prerequisite default branch "main": git fetch origin
main: ... Operation timed out`. The owner had not created the item's branch
yet, and `roundfix deliver retry` refused it with `use item branch "": use item
branch: branch is required`: the retry handed every park outside the
prerequisite and archived-item parks to `UseItemBranch`, which needs a branch.
When the item's worktree creation itself fails, the park also writes the
in-memory item, which has no branch, over the branch the owner had just
recorded, leaving a worktree recorded without its branch; the owner's next
`RecordDeliveryQueueItemWorktree` then refuses with `recorded branch is empty`.

On 2026-11-12 quaestor resolved a modify/delete Pull Request conflict by
removing the last file under `docs/specs/` from its item branch. Git does not
track directories, so the directory left the working tree, and the retry
stopped with `read delivery authorization: ... resolved to "<worktree>/docs/specs",
which does not exist`. Two layers needed the directory: `Authorization`
resolved the Spec Root with the strict resolver, and `ReadSpecAuthorization`
runs `git -C <spec root>`, which fails on a missing directory and would have
parked the item `unauthorized` instead. The grant itself is read from the
commit before the archive and never from the working tree.

After the operator recreated the directory, two retries parked the item
`review-stale`. The first retry had recorded the resolution commit as the new
candidate and the gate had passed, but publication had failed before the push;
the next retry of an archived item with a Pull Request re-entered at
`checking`, where the Pull Request still held the earlier candidate.

## Decision

**A branchless item is requeued.** A Delivery Retry of a parked item whose
recorded Item Branch is empty returns it to `queued`, as ADR-0193 already does
for a prerequisite park, and counts the retry. The same write clears a
worktree recorded without its branch, so the owner records a new pair. No
evidence of the item is read: it never started.

**The authorization read needs no Spec Root directory.** The delivery
authorization resolves the Spec Root with the absent-tolerant resolver, and the
Spec authorization reader resolves an absent Spec Root inside the repository
through its nearest existing ancestor, adding the missing components to the
record path it reads at the revision.

**A Pull Request behind the candidate returns the item to the gate.** At
`checking`, when the Pull Request's head is a candidate of the item other than
the current one, the owner logs it and sets the item to `gating`. The
repository gate runs again on the current candidate, and publication pushes it.
A head the item never had still parks `review-stale`, with the next action
`in <worktree>, merge into the item branch any commit its Pull Request's head
branch holds that the item branch lacks, then run roundfix deliver retry <slug>`.

**A network failure still parks.** The owner does not retry a failed fetch of
the prerequisites' default branch by itself; the park is visible and the
retry above now recovers it.

## Considered Options

- **Push the item branch from `deliver retry`.** Rejected: the retry would
  reach the network and push a candidate that may never have passed the
  repository gate, as after a `gate-failed` park of a derived-merge candidate.
- **Push from `checking` directly.** Rejected for the same reason: `checking`
  cannot tell a gated candidate from one that parked at the gate.
- **Only reword the `review-stale` park.** Rejected: the item would still need
  a manual push for a candidate the queue itself recorded.
- **Recreate the Spec Root directory before the read.** Rejected: it writes
  into the item worktree a directory the branch does not hold, and the read
  needs nothing from the working tree.
- **Stop the park from clearing the recorded branch.** Rejected: the branch a
  failed worktree creation recorded may not exist, so the retry would still
  refuse it as a missing Item Branch.
- **Retry a failed fetch inside the owner, three attempts with a 10 to 20
  second wait, as `actions/checkout` does for its own fetches.** Rejected for
  this decision: that helper wraps every network command it runs, and the
  owner has several more (the item-branch fetch, the push and every `gh`
  call); retrying one fetch alone would give a partial, misleading guarantee.
  A uniform policy for the owner's network steps is a separate decision.

## Consequences

- An item parked before its branch exists returns through `deliver retry` and
  starts as a new item would.
- A conflict resolution that empties the Spec Root no longer stops
  publication or merge.
- A candidate that never reached its Pull Request is gated again and pushed by
  the owner, at the cost of one more repository gate run.
- A `review-stale` park tells the operator what to do; status and the Pending
  Question print it.
- This decision refines ADR-0193 and ADR-0223 and supersedes none.
