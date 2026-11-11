---
status: accepted
created_at: 2026-11-10T23:30:00Z
updated_at: 2026-11-10T23:30:00Z
deprecated_at: null
superseded_by: null
---

# A Machine Verification Capacity bounds the Verifications of every repository

On 2026-11-10 the Delivery Queue of vortex stopped three times in
`run-unresolved` on time-sensitive tests its Tasks never touched, each time
while a Run of another repository verified on the same twelve-core machine;
the same tree passed alone in 285 s. Roundfix bounded whole Runs
(`runs.max_active`) and the Verifications inside one Run
(`verification.concurrency`), but nothing bounded the Verifications of
different repositories against each other. The maintainer wants the projects'
Agents to keep running in parallel and chose a machine capacity that only the
verification step waits for, `1` by default.

## Decision

**One User Config key bounds the machine.** `runs.max_verifying` is the
Machine Verification Capacity: how many Verifications Roundfix runs at once on
the machine, across every repository that shares the Roundfix Home. It is an
integer, `1` by default, and `0` turns the bound off, as `runs.max_active: 0`
does. Only User Config sets it; a Project Config value is ignored with the
warning of the other operator-only keys. `roundfix init --scope user` and
`roundfix setup` write it into the User Config they create.

**It wraps the per-Run capacity and gates exactly the Daemon's and the
queue's Verifications.** Every Verification attempt the Daemon starts takes
its unit of the Run's own Verification Capacity first and then one unit of
the machine's, and gives the machine unit back first: the pre-work probe, the
repository precondition, each Task attempt with the repository Verification
appended at settlement, the exclusive retry after exit `75`, the Integrated
Repository Verification and the QA Task's repository Verification the Daemon
runs before the QA Agent. The Delivery Queue's repository gate takes one unit
for its repository Verification and its Linux Lane together. Nothing a
process runs inside an Agent session is gated: the QA Agent's own suites and
an Agent's focused tests run whenever the Agent runs them. `roundfix spec
check --run-verification`, `roundfix settle`, the item binary build and the
derived-file regeneration of the queue, Worktree Bootstrap, and the
Verification of review Runs (`fetch`, `resolve`, `watch`) are not gated.

**The capacity is a set of entry files, not a table.** Each waiter and holder
writes one entry under `verification-capacity/` in the Roundfix Home, naming
its ticket, state, Run or Delivery Queue item, repository, Spec, step, owner
PID and start-time identity. A process changes the directory only while it
holds the directory's lock file, with the same advisory lock the Run
Database writer uses. Tickets come from a sequence file inside the same lock,
so waiters are served in arrival order: a waiter takes a unit only when the
holders and the waiters ahead of it leave one free.

**A holder that dies is reclaimed by proof, never by elapsed time.** Before it
counts, an acquirer removes every entry whose recorded owner is gone by the
one store rule of ADR-0373: an absent PID, or a live PID whose start identity
was read and differs. An unreadable identity, a foreign token or an empty
recorded identity counts as present. A cancelled or stopped wait removes its
own entry, and a release removes it with a context that the Run's
cancellation does not cancel.

**Waiting is visible and costs no budget.** `roundfix runs list`, `roundfix
deliver status` and `roundfix doctor` name each holder and waiter with its
step and place in line, and the waiting Run journals a Run Event when its
wait begins. The Run Budget does not run while any Verification of the Run
waits for the Machine Verification Capacity; the wait moves the deadline
later by its own length. A Stop Request during the wait starts no command and
leaves the Task resumable, as a wait for the Run's own capacity does.

## Considered Options

- `runs.max_active: 1`, rejected by the maintainer: it stops the projects'
  Agents from working in parallel.
- A Run Database table, rejected as in ADR-0304 and ADR-0365: a schema change
  makes every older Roundfix binary on the machine refuse the shared Run
  Database, and the entries are short-lived state that needs no history.
- An advisory lock per slot whose release by the kernel frees a dead holder:
  rejected, because it gives no arrival order, so a waiter can lose every
  race, and a reader still needs the holder's recorded identity to name it.
- Reclaiming an entry after a timeout or a missed heartbeat: rejected, because
  a Verification can legitimately run for an hour, and a live holder whose
  slot is taken lets two suites collide again.
- Taking every machine unit for the exclusive retry after exit `75`: rejected,
  because that retry would wait for every other repository's Verification to
  drain; it takes the Run's whole capacity and one machine unit, and with the
  default of `1` it runs alone on the machine anyway.
- Gating the QA Agent's suites: impossible from the Daemon, which never sees
  the commands an Agent runs.
- A core budget exported to the Verification commands: complementary, not a
  substitute; it lightens the Verifications the capacity still lets run
  together, and builds on the capacity value recorded here.

## Consequences

- After the update, every machine verifies one repository at a time unless
  its User Config raises or turns off `runs.max_verifying`; the Agents of
  every Run keep working in parallel.
- A Run can wait for another repository's whole Verification; its Run Budget
  does not shrink while it waits, and an absolute deadline such as the queue
  deadline still applies to an item that has not started.
- A live holder that never ends, such as a hung Verification command, keeps
  its unit until its Run is stopped; the surfaces name it so the operator can
  run `roundfix stop`.
- An older Roundfix binary takes no unit and refuses a User Config that sets
  `runs.max_verifying` as an unknown key; the default needs no setting.
- The capacity is per Roundfix Home: a process with another `HOME` has its own
  capacity.
