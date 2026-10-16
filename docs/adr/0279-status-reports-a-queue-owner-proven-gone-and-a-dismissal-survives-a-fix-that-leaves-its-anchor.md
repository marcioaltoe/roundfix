---
status: accepted
created_at: 2026-10-16T00:00:00Z
updated_at: 2026-10-16T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor

On 2026-10-13 the machine rebooted and the Queue Owner died with it; for about
eleven hours `roundfix deliver status` printed the item as `checking` with no
warning, and `roundfix deliver resume` recovered it at once. On 2026-10-14 a
parked item had two findings, one to fix and one to dismiss; after the fix was
committed `roundfix review dispose F1 --dismiss` refused because HEAD no longer
equalled the reviewed head, and the only route was to reset the branch, dismiss
and reset forward. In the same delivery `roundfix deliver retry` run inside the
item worktree said that no Delivery Queue existed for that worktree instead of
naming the checkout that owns the queue.

**Status reports a recorded owner proven not running.** When the Delivery
Queue records an owner PID and at least one item can still advance (its stage
is neither `merged` nor `parked`), `deliver status` proves that owner the way
`deliver resume` and `deliver retry` already do: the PID is absent, or it is
alive under a start-time identity that differs from the one recorded. Either
proof adds one `Owner:` line naming the PID, the state and `roundfix deliver
resume`. A live PID whose identity matches, or whose identity cannot be read,
adds no line; an identity that cannot be read writes a diagnostic to stderr.
Status stays read-only and never reclaims the owner record; `resume`, `retry`
and `start` still do. A queue with no recorded owner adds no line.

**A dismissal survives a fix that leaves its anchor.** `roundfix review
dispose <id> --dismiss` keeps accepting a dismissal at the reviewed head. It
now also accepts one when the current HEAD descends from the reviewed head, as
`git merge-base --is-ancestor` proves, and the finding's recorded anchor is
unchanged: its path exists at both commits and the anchored lines, at the same
line numbers, are byte-identical (an anchor whose start line is 0 compares the
whole file). The disposition stays tied to the reviewed head, as every reader
of the ledger expects, and records the HEAD where the anchor was proven in a
new optional `currentHead` field. A finding without a recorded anchor, a head
that does not descend, or a changed, moved or deleted anchor is refused with a
reason that says which.

**A command run in an item worktree names the owning checkout.** When a
`deliver` command that reads the queue (`status`, `resume`, `retry`, `stop`)
finds no Delivery Queue for its checkout, and that checkout is the recorded
worktree of an item, compared as resolved paths (ADR-0268), the refusal adds
the item and the checkout whose queue holds it and says to run the command
there. Roundfix does not resolve the queue on its own: the queue's
configuration and its detached owner are derived from the checkout the command
runs in, so acting on another checkout's queue from an item worktree would
start an owner that cannot find it.

**Rejected.** Reporting a queue without any recorded owner would also name the
short window inside `deliver start` and `deliver retry` before the new owner
claims the queue, and the state after an explicit `deliver stop`; the observed
defect is a recorded owner that died. Proving liveness by PID alone would
report a reused PID as a running owner. Reclaiming the owner record from status
would make a read command write. Dismissing at a descendant whenever the file
is unchanged would refuse a fix to an unrelated part of the same file, and
accepting any descendant would carry a dismissal past a change to the very
lines it judged. A `review dispose --head` flag stays rejected, as ADR-0272
rejected it. Resolving the owning queue automatically from an item worktree
would act on a repository other than the one the command runs in.

## Consequences

- A rebooted machine's queue shows `Owner: Delivery Queue owner PID <pid> is
  not running; run roundfix deliver resume` until it is resumed.
- A fix and a dismissal can be recorded in either order when the fix leaves
  the dismissed finding's anchored lines alone and does not move them; a fix
  that changes or moves them requires the dismissal before the fix, as before.
- The disposition ledger gains an optional field; readers that ignore it,
  including the review correction proof and the review lineage, keep their
  answers.
- `deliver status` keeps its exit code, output order and every other line, and
  gains no JSON form.
