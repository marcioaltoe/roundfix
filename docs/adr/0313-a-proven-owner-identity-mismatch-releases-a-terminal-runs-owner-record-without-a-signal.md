---
status: accepted
created_at: 2026-10-22T12:00:00Z
updated_at: 2026-10-22T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A proven owner identity mismatch releases a terminal Run's owner record without a signal

A terminal Run keeps the PID and start identity of the process that owned it.
When that process has exited and the system has given its PID to another
process, `roundfix reconcile` reads a start identity that differs from the
recorded one, refuses to prove the owner, and lists the Run under
`preservedCandidates` on every invocation: on 2026-10-22 the maintainer's
machine showed four such Runs from 2026-09-29 to 2026-10-14. ADR-0266 already
holds that a different identity "counts as proof that the recorded process
ended". The maintainer kept the Backlog Entry's expected behavior: reconcile
treats that proof as what it is and releases the record, and an identity it
cannot read stays preserved.

## Decision

**A proven mismatch is its own error.** The owner proof distinguishes three
refusals. A live process at the recorded PID whose start identity was read,
carries the current platform's prefix and differs from the recorded identity
is a proven identity mismatch, `ErrOwnerProcessIdentityMismatch`. An identity
that cannot be read or compared (a read failure, a platform without start
identity, a token from another platform or an older format) is unreadable, as
before. A PID that is not positive or is the inspecting process's own PID is
unproven without a mismatch, as before. The mismatch wraps the unproven
error, so every caller that refuses an unproven owner, Force Stop first,
still refuses it, with the same message bytes.

**The record of an ended owner is a Stale Run Owner Record.** For a terminal
Run, a proven mismatch proves the recorded owner ended. Reconcile no longer
preserves that Run's owner: it reports one `runOwnerRecord` entry in
`processCandidates`, with the owner PID, an empty process list and the proof,
and never enumerates or inspects the tree of the process that now holds the
PID. `roundfix reconcile --apply` releases the record: in one Run Database
write it clears the Run's owner PID and identity only when they still equal
the inspected ones and the Run is still terminal, and journals a Run Event
that keeps the released PID, identity and reason. The release sends no
signal, reaches no process and changes no outcome or timestamp of the Run.
A record that changed or was already released is reported as already
released and nothing is written. The debris summary counts releases as
`runOwnerRecordsReleased` in JSON and `run-owner-records-released` in text.

**Nothing of the recorded owner's tree can remain to be attributed.** POSIX
does not reuse a PID while a process group with that ID exists ("if there
exists a process group whose process group ID is equal to that process ID,
the process ID shall not be reused by the system until the process group
lifetime ends", Base Definitions 4.17), so when the PID names a stranger the
recorded owner's group had ended. The mismatch cannot be undone either: the
recorded process does not come back. So the release needs no second proof at
apply time beyond the compare-and-set on the record itself.

**Everything short of a proven mismatch stays as it is.** An unreadable or
unproven identity, a Run recorded with no identity, an Active Run and a live
owner whose identity matches keep today's classification and actions. Force
Stop still never signals a PID whose identity does not match, and its
`--owner-identity-unreadable` flag still refuses a proven mismatch.

## Considered Options

- Keeping the preserved entry was rejected: it repeats forever and the
  operator has no command that ends it.
- Treating a mismatch as an absent owner, silently, was rejected: the record
  would be re-proven on every reconcile, and the stale PID could name yet
  another process later; the maintainer asked for a report and an explicit
  `--apply`.
- Re-proving the mismatch at apply time was rejected: the proof cannot be
  undone, and the compare-and-set already refuses a record that changed.
- Deleting the Run was rejected: Run Retention owns whole-Run removal, and the
  Run's outcome, worktree and branch evidence stay useful to reconcile
  (ADR-0053).
- Releasing from `doctor`, `runs` or another read command was rejected: those
  commands stay read-only, and the reconcile dry-run already names the record.
- Signalling the process that now holds the PID was never an option
  (ADR-0266).

## Consequences

- A terminal Run whose owner PID was reused leaves `preservedCandidates`,
  appears once as an `runOwnerRecord` process candidate, and disappears from
  every later report after `--apply`.
- After the release the Run records no owner, so `doctor`'s residue check no
  longer reads a lineage for it.
- A consumer of `processCandidates` that assumed every entry names live
  processes reads the new `kind` to tell a Run owner record apart; every other
  report field keeps its meaning.
- On a host where a stored identity token could differ from the live token of
  the same process, the record would be released while the owner lives; the
  identity prefix already refuses tokens from another platform or format, and
  no signal is ever sent on that path.
