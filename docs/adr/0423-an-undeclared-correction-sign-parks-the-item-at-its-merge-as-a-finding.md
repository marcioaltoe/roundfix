---
status: accepted
created_at: 2026-11-17T19:30:00Z
updated_at: 2026-11-17T20:30:00Z
deprecated_at: null
superseded_by: null
---

# An undeclared Correction Sign parks the item at its merge as a finding

The setting that makes the declaration required is the Project Config key
`delivery.require_defect_class`, read from the default branch and never from
the item's tree. With it on, the Queue Owner parks an item at its merge with
the blocker `defect-class-undeclared`, of Park Class `finding`, while a
Correction Sign is named by no standing Correction. With it off the same step
records the signs, prints nothing and can never stop a merge.

## Context

ADR-0417 decides the requirement and leaves the key, the blocker and its Park
Class to the TechSpec. Read at `a8236a87`: the merge stage first reads the
Spec's authorization and parks `unauthorized` before it records a merge
intent; a Delivery Retry of an archived item whose head did not move returns
it to `checking`, which sets `merging` again; and the Queue Owner already
reads one delivery setting, `delivery.linux_lane`, from the Project Config at
the default branch head, so that an Item Branch cannot change it (ADR-0288).
More than a dozen tests assert the lines of `roundfix deliver status` byte for
byte.

## Decision

**The key.** `delivery.require_defect_class` is a boolean, `false` when
absent. A value in User Config is removed before validation with one warning,
and a value in Project Config that is not a boolean is refused by name.

**Where it is read.** The Queue Owner and `roundfix deliver status` read the
key from the Project Config at the default branch head as the local refs hold
it, the remote-tracking ref first, and never from the item's tree or from the
owner's checkout.

**The step.** In the `merging` stage, after the authorization is read and
before the merge intent is recorded, the Queue Owner reads the Correction
Signs of the candidate head and appends them to the corrections record.

- With the key on, a sign that no standing Correction names parks the item as
  `defect-class-undeclared: <signs>`, and a failure to read the signs or the
  record fails the step as any other step error does.
- With the key off, the step never parks, never fails and prints nothing. A
  failure is written into the signs line it appends.
- When no default branch head resolves, the step requires nothing and records
  nothing.

**The park.** Its Park Class is `finding`: the item's own work lacks an
answer only the Supervisor can give. Its next action names the five classes,
the answer `not-a-correction`, the declaring command and the Delivery Retry. The retry is counted, and it
takes the existing way back: `checking`, then `merging`, where the step runs
again.

**The status line.** With the key on, `roundfix deliver status` prints one
`Corrections:` line for each unmerged item whose branch tip holds a sign no
standing Correction names. With the key off, or when it cannot read the key
or the signs, it prints nothing.

## Considered Options

- **A new Park Class.** Rejected: ADR-0417 leaves Park Class as it is, and the
  status view prints the class as a closed value.
- **The `authorization` class.** Rejected: it names a grant only the
  maintainer gives.
- **The step at the repository gate, before the pull request.** Rejected: a
  failed check is corrected after it.
- **The key read from the owner's checkout.** Rejected: an unmerged or
  uncommitted edit would turn the requirement off, and the status view and a
  running owner could disagree.
- **A retry straight to `merging`.** Rejected: it adds a retry target, and
  `checking` proves the head again at the cost of one read of its checks.
- **A warning in the owner's log when the key is off and the signs cannot be
  read.** Rejected: a repository that never set the key would print a line it
  does not print today.

## Consequences

- A repository turns the requirement on with one commit on its default
  branch, and a Queue Owner already running obeys it at its next merge.
- A declaration made while the owner runs is read at the merge: the step
  reads the file each time.
- With the key off, every delivery of the Delivery Queue still leaves a signs
  line under Roundfix Home, which is what lets the report count `undeclared`.
- An item that reaches its merge undeclared costs one counted retry. The
  status line shows the sign before that.
