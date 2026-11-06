---
status: accepted
created_at: 2026-11-05T12:00:00Z
updated_at: 2026-11-05T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A numbered file collision with an earlier queue item parks before the Run

On 2026-11-05 fiscus reported five Specs of one Delivery Queue that all
started their database migrations at number 0052. Each Spec took the next
number of the main it knew, and the queue found the clash only at merge,
after the Run, the review, the gate and the checks, so each merge made the
others renumber by hand.

A repository now declares its numbered files in Project Config under
`delivery.numbered_files`, with the path grammar of a Derived Path
Declaration. A matching file whose name starts with digits is a Numbered File,
and its directory and the value of those digits are its place in a sequence.
Delivery Revalidation, before an item's first Run and on a retry of an item
that has not run, collects the Numbered Files the item adds: those its branch
adds since its merge base with the default branch, and those its not yet
completed Tasks declare under `creates:`. One that shares a directory and a
number with a file on the default branch, or with a file an earlier item of
the same queue adds while it is not merged, is a Numbered File Collision. The
item parks as `numbered-file-collision` in the conflict Park Class, naming
each pair, the earlier item or the default branch, and the next free number,
and the operator renumbers the plan, or delivers the earlier item first and
merges the default branch, before the retry.

## Considered Options

- A default pattern list. Rejected: names that start with digits are common
  outside sequences, such as dated Backlog Entries and Findings, so any
  default would park unrelated items; the declaration is off until a
  repository makes
  it.
- Waiting for the earlier item, as a Prerequisite Spec does. Rejected: a plan
  that names a taken number must change even after the earlier item merges,
  and a branch created before that merge keeps its old base.
- Renumbering by the queue. Rejected for now: renaming a migration also
  rewrites tool-owned files, such as a journal and its snapshots, that only
  the repository's generator owns; the queue reports the next free number and
  leaves the change to the operator or the Agent.
- Checking at merge only, as today. Rejected: it spends the Run, the review,
  the gate and the checks before the conflict appears.

## Consequences

A number a Run generates without its plan declaring it is seen only once a
branch holds it; from then on every later item is held to it. A retry of an
item whose Run already started is not checked and still meets the conflict at
merge, where `delivery.derived_paths` and `pull-request-conflict` keep their
behavior. Queue order decides precedence: the later item renumbers.
Timestamped folders, as drizzle-kit v1 writes them, are outside the
declaration.
