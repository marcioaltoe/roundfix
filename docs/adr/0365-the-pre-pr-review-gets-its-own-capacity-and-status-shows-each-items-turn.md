---
status: accepted
created_at: 2026-11-05T12:00:00Z
updated_at: 2026-11-06T12:00:00Z
deprecated_at: null
superseded_by: null
---

# The pre-PR review gets its own capacity, and status shows each item's turn

The Queue Owner runs the pre-PR review of items in `reviewing` beside the one
item it advances through every other stage, up to a Review Capacity set by
`delivery.review_concurrency` in User Config, `2` by default. It records each
item's Queue Turn in a Turn Record beside its console log, and `roundfix
deliver status` prints one `Turn:` line per item for an owner it proves live.

## Context

On 2026-11-05 Fiscus reported a queue of about twenty-one Specs running for
two days, with five items in `reviewing` and no review running while one
Implement Run held the owner, and `deliver status` saying nothing about the
wait. The owner advanced one item at a time through every stage, the review
inline; a Delivery Retry stored `reviewing` and the owner reached the item
only when its loop came back to it. The maintainer chose: "O deliver status
mostra 'aguardando capacidade' e a posição. O review pré-PR, que só lê, ganha
uma concorrência própria, separada da do implement."

## Decision

**A stage is the next step.** The stored stage keeps its name and meaning;
whether the step runs or waits is the item's Queue Turn.

**Reviews run in slots.** Each `Run` of the owner holds review slots and a
poller. The owner's loop hands an item that reaches `reviewing` to the slots
and goes on; the slots review waiting items in position order, up to the
Review Capacity. The poller hands every retried item at `reviewing` whose
item worktree exists to the slots every 30 seconds, so a review never waits
for another item's Implement Run. A review reads only its own item worktree
in a read-only session, builds its item binary into its own path, and the
owner never prepares the worktree of an item the slots hold. A reviewed item
continues in the owner's loop without preparing its worktree again.

**Order is kept.** A queued item starts only when no earlier item is between
its Implement Run and its merge, and never ahead of an earlier queued item,
so it still starts from a default branch holding every earlier merged item.
Archive, gate, publication, check wait and merge stay one item at a time.

**The knob is the operator's.** `delivery.review_concurrency` is a User
Config integer of at least `1`; Project Config is ignored with the warning of
the other operator-only keys. The owner reads it when it starts. `deliver
start` and `deliver status` print it on the `Limits:` line as
`review concurrency <n>`.

**The turn is a file.** The owner writes the Queue Turn of every item it has
not finished, its PID, start-time identity and Review Capacity to
`turn.json` beside its console log whenever a turn changes and on every poll.
Status reads it only for the owner it proves live and only when the PID and
identity match, and writes nothing.

**Amended by [ADR-0369](0369-the-checking-stage-waits-for-the-head-check-set-and-parks-on-any-failed-check.md) (2026-11-06).** The Turn Record of the item in `checking` also lists the members of its Head Check Set still pending, and its `Turn:` line reads `advancing at checking; waiting for <n> check(s): <check>; <check>`. Every other turn form, the file, its readers and the Review Capacity still hold.

## Considered Options

- A `waiting` stage or an item column: rejected, as in ADR-0304, because a
  schema change makes older item binaries refuse the shared Run Database.
- Deriving the turn in status from the Run Database alone: rejected, because
  nothing there says which item the owner is advancing.
- Starting the next queued item while an earlier one is in review: rejected,
  because the next item would start from a default branch without the
  earlier item, and Fiscus already pays for Specs that start from an older
  main.
- A Project Config key or a `deliver start` flag: rejected, because the
  reviews spend the operator's provider sessions on the operator's machine,
  and a flag recorded per queue needs a schema change.
- A default of `1`: rejected, because one slot leaves a batch of retried
  items in a line; two reviews beside one Implement Run of Task Capacity 2
  stay within four Agent sessions, fewer than the machine runs under the
  default ceiling of three Active Implement Runs.

## Consequences

A retried item is reviewed while another item implements, and several waiting
reviews run together. A reviewed item still waits for the owner's turn for its
archive and merge, and a fresh queue reviews one item at a time as before.
Review tokens stay outside the Queue Token Ceiling (ADR-0199), and reviews
count neither against `runs.max_active` nor a Run's Verification Capacity. An
older owner writes no Turn Record, so its status shows no `Turn:` line until
it is restarted on a binary that writes one. An older Roundfix binary,
including an item binary built from a branch older than this decision,
refuses a User Config that sets `delivery.review_concurrency` as an unknown
key; the default needs no setting.
