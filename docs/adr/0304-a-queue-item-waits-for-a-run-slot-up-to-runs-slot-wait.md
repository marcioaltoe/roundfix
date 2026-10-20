---
status: accepted
created_at: 2026-10-20T22:00:00Z
updated_at: 2026-10-20T22:00:00Z
deprecated_at: null
superseded_by: null
---

# A queue item waits for a Run slot up to runs.slot_wait

On 2026-10-20 the Delivery Queue item of Spec 0289 was parked as
`delivery-error` because other repositories' Implement Runs held every slot of
the machine-wide ceiling, and it waited for a manual retry. When the only
preflight refusal of the item's `roundfix implement` is that ceiling, the item
now waits for a free slot, up to `runs.slot_wait` (User Config, a Go duration,
2h by default, `0` parks at once) or the queue's own deadline, and then parks
exactly as before.

## Decision

**Only the ceiling waits.** The queue recognizes the refusal when the child
exits with the preflight code, created no Run, and the first line under its
`Reason:` label begins with the store's ceiling message. Any other preflight
refusal parks as before.

**The wait is a Run Slot Wait on the running item.** The item keeps the stage
`running` and records `run-slot-wait: until <deadline>; <n> Run(s) hold <max>
slot(s): <run id> in <repository> for <spec>, ...` as its blocker text, which
`roundfix deliver status` prints on the item's line. No stage, column or schema
is added. The text is cleared before the Run starts and when the wait ends.

**It polls cheaply and the store decides.** Every 30 seconds the queue counts
the Active Implement Runs whose owner is alive against `runs.max_active`, the
way the child would after reclaiming orphans, without spawning anything. When
a slot looks free it starts the child again; the child's preflight and the
store's `CreateRun` transaction stay the only authority. A queue that loses
the race receives the same refusal and keeps waiting toward its original
deadline, so two waiting queues never both hold a Run beyond the ceiling.

**The deadline is the earlier of two.** The wait ends at the first refusal's
time plus `runs.slot_wait`, or at the queue's deadline when that is earlier.
At its own deadline the item parks `delivery-error` with today's reason; at the
queue's deadline it parks `queue-deadline`, so the queue deadline still wins.
A stopped owner clears the wait text; a resumed owner or a Delivery Retry
starts a new wait.

**User Config only.** `runs.slot_wait` sits beside `runs.max_active`. Only the
machine's owner, who sets the ceiling, decides how long queues on that machine
wait for one another; a `runs.slot_wait` in Project Config is ignored with the
same warning as `runs.max_active`. A negative or unparsable value fails
configuration validation.

## Considered Options

- A distinct exit code for the ceiling refusal was rejected: `implement`'s
  exit codes are public, and every preflight refusal shares `2`.
- Retrying the child on every poll was rejected: each attempt rebuilds the
  item binary, and the count it would learn is readable from the Run Database.
- A new `waiting` stage or an item column was rejected: a schema change makes
  older item binaries refuse the shared Run Database.
- A Project Config `delivery.run_slot_wait` key was rejected: the ceiling it
  waits on is machine-wide and User Config only, and a repository should not
  decide how long its queue holds a machine's capacity.
- Waiting without a bound was rejected by the maintainer's decision ("Espera
  com limite").

## Consequences

- A queue started while the machine is busy delivers when capacity returns,
  without an operator.
- A queue item can occupy its queue for up to `runs.slot_wait` before it
  parks; later items in the same queue wait behind it, as with any running
  item.
- An older Roundfix binary refuses a User Config that sets `runs.slot_wait` as
  an unknown key; the default needs no setting.
