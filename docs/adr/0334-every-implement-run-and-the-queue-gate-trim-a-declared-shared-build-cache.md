---
status: accepted
created_at: 2026-10-26T00:00:00Z
updated_at: 2026-10-26T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Every Implement Run and the queue gate trim a declared Shared Build Cache

ADR-0327 lets the GC Command trim a declared Shared Build Cache, but only when
someone runs `roundfix gc`. Between two runs the cache this repository
declares grew to 51 GB and left 157 MiB free on 2026-10-25, three days after
43.7 GB had stopped a Run on 2026-10-22. New build output enters the cache
where a fresh worktree path is compiled: in the Run and Task Worktrees of an
Implement Run and in a Delivery Queue item worktree at its repository gate.
On 2026-10-26 the maintainer approved the cycle's triage, whose first item
was this automatic trim.

## Decision

**The trim ADR-0327 defines also runs at the start of every Implement Run and
before the Delivery Queue's repository gate, under one time budget, and never
changes their outcome.**

- An Implement Run trims each declared cache right after its Run is recorded
  and before its Run Worktree is created. That covers a direct
  `roundfix implement` and the `implement --spec` child the Delivery Queue
  starts for each item.
- The Delivery Queue trims each declared cache before it runs the repository
  gate's Verification in an item worktree, the queue's only build outside a
  Run and the stage a `deliver retry` after a full-disk park resumes at.
- The trim keeps every rule of ADR-0327: it acts only above `max_size_mb`,
  removes only entries idle for at least `idle` (two hours or more), oldest
  first, rereads each entry's modification time just before removing it,
  and refuses a symbolic link or a directory without the go command's README.
  The go command tolerates a missing entry as a miss, so Runs that share the
  cache may trim it at the same time.
- One budget of thirty seconds bounds all declared caches at one trigger.
  When it runs out, the trim stops between entries, keeps what it removed,
  and the next trigger continues.
- A failure, refusal or stop never fails, parks or delays the Run or the gate
  beyond the budget. An Implement Run records one `daemon.status` Run Event
  per declared cache with the outcome and its cause; the queue writes one
  line per cache to the Queue Owner's log.

## Considered Options

- Keep the trim manual (`roundfix gc` only). Rejected: it failed twice in
  four days on the machine that delivers this repository.
- Trim from a timer inside a long-lived process. Rejected: the Daemon lives
  only as long as one Run, and Roundfix starts no background service.
- Trim before every Task. Rejected: a Run's growth is at most a few full
  builds, and a trim per Task adds its walk to each Task's start.
- Trim at the start of `resolve` and `watch` Runs too. Rejected: they build in
  the user's own checkout, whose path repeats, so they reuse cache entries
  instead of adding a full build; they still benefit from the other trims.
- Refuse to start a Run below a free-space floor. Deferred: the measured
  failures happened when space ran out during a Run, which a start floor does
  not prevent, and `roundfix doctor` now warns about low free space.

## Consequences

- A repository with no `worktree.shared_caches` declaration sees no change:
  no Run Event, no output line and no time spent.
- An Implement Run's start may take up to the budget longer when its cache is
  far above its limit; measuring 21.6 GB in about 38,000 files took 1.2 s, and
  removing 40,000 entries 5.6 s, on this repository's machine.
- `roundfix gc` keeps its report and its exit codes, and it can still trim on
  demand.
