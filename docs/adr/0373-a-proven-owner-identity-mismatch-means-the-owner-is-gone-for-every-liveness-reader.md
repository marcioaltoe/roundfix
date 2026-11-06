---
status: accepted
created_at: 2026-11-06T16:00:00Z
updated_at: 2026-11-06T16:00:00Z
deprecated_at: null
superseded_by: null
---

# A proven owner identity mismatch means the owner is gone for every liveness reader

On 2026-11-06 argus reported that `roundfix doctor` 0.136.0 ends with
`residue: partial` in every repository, naming three terminal Runs of other
repositories with `could not read process table for Run <id>: owner process
<pid> failed step "prove owner process identity": owner process identity
cannot be proven: live process start identity does not match the recorded
owner identity`. The PIDs belonged to unrelated macOS processes; the recorded
owners had died and the system had reused their PIDs. No `reconcile` or `gc`
run changes the result, so the warning never leaves.

ADR-0313 made a proven mismatch its own error and let `reconcile --apply`
release the owner record of a terminal Run. Doctor's residue check, which
reaches the same proof through `InspectTreeProcesses`, still files every
error as an unreadable process table. Three other readers ask only whether
the PID is alive and never read its identity, so a reused PID keeps a dead
owner present: the orphan reclaim of an Active Run's lock, the Delivery
Queue's Run-slot count and `roundfix deliver start`'s takeover of a queue
whose owner died.

POSIX.1-2017, section 4.14, says "A process ID shall not be reused by the
system until the process lifetime ends", and a process group ID not until
the group's lifetime ends. A live process at the recorded PID whose start
identity differs is therefore proof that the recorded owner ended and that
the process group it led is empty.

## Decision

**A proven start-identity mismatch at a recorded owner PID means the owner is
gone, wherever Roundfix asks whether a recorded owner is still there; every
answer short of that proof keeps today's behavior.**

- One store function answers the question: the owner is gone when the PID
  is absent, or when it is alive and its live start identity was read,
  comparable and different from the recorded one; it says whether the PID was
  reused. An unreadable identity, a recorded token from another platform, an
  empty recorded identity on a live PID, and the caller's own PID are not
  proof, and the owner counts as present.
- Doctor's residue check treats a mismatch as an owner that is gone: the Run
  contributes no residue and no read failure, and the check does not walk the
  process tree of the PID's new holder. `partial` stays for a process table
  that cannot be read and for a Run whose recorded identity is unproven.
  Doctor still releases nothing (ADR-0127).
- The orphan reclaim of an Active Run's lock, the Delivery Queue's Run-slot
  holders and `roundfix deliver start` use the same function, so a reused PID
  frees the lock, excludes the holder and lets a new queue take over, as an
  absent PID already does (ADR-0044). The reclaim reason names the reuse:
  `owner process <pid> gone (pid reused); lock reclaimed`; an absent owner
  keeps `owner process <pid> not running; lock reclaimed`.
- Nothing that signals changes: Force Stop, `deliver stop` and every
  terminate path still refuse a PID whose identity does not match and send it
  nothing. `reconcile` (ADR-0313) and `deliver status`, `retry` and `resume`
  (ADR-0279), which already read a mismatch as gone, keep their behavior.

## Considered Options

- Reporting a reused owner as a `found` or `warn` residue line. Rejected: the
  process at the PID is not Roundfix's, and the owner record it points at is
  harmless; `reconcile --apply` already releases it.
- Walking the new holder's process tree after a mismatch. Rejected: its
  process group and session belong to an unrelated process, so every member
  would be misreported as Roundfix residue.
- Treating an unreadable identity as gone too. Rejected: a permission error
  or a vanished read proves nothing, and the lock reclaim would then free a
  slot an owner still holds.
- Releasing the owner record from Doctor. Rejected by ADR-0313: Doctor and
  the other read commands stay read-only.

## Consequences

- `residue:` reads `ok` when the only problem was reused owner PIDs, and
  `partial` keeps its meaning of an incomplete inspection.
- A lock, a Run slot or a Delivery Queue held by a dead owner whose PID was
  reused is freed by the same paths that free an absent owner.
- The function costs one start-identity read for a live PID with a recorded
  identity; an absent PID costs what it costs today.
