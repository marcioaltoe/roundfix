---
status: accepted
created_at: 2026-11-17T19:30:00Z
updated_at: 2026-11-17T19:30:00Z
deprecated_at: null
superseded_by: null
---

# The Corrections Report counts the Tasks the plan held on the default branch and starts no fetch

The Corrections Report takes a delivered Spec, its delivery date and its
Delivery Type from the default branch as the local refs hold it, and the size
of the Spec from the Task Graph the default branch held just before the
delivery merge. It never reads the Spec's source revision, which the squash
merge leaves unreachable, so it never fetches. Its window is the Run
Retention window.

## Context

ADR-0419 fixes what the report prints and ADR-0417 that it reads this
machine's records. Read at `a8236a87`: an Archive Record holds no Task count,
only the identifier of its QA gate and a `source_revision` on the Item
Branch. The squash merge leaves that commit reachable from no branch, and the
one function that reads it, `spec.EnsureSourceRevision`, fetches the pull
request's head when the object is missing. The status view already lists
delivered Specs from the default branch head with Git alone and reads no
working tree for them.

A Spec's plan is on the default branch before its delivery starts, and the
delivery merge removes the Spec folder in the commit that adds the Archive
Record. A Corrective Task is appended on the Item Branch, by whoever
corrects, and is numbered above the QA gate.

## Decision

**A delivered Spec.** A Spec whose Archive Record a commit of the default
branch added inside the window and whose record is not a cancellation. The
commit's time is the delivery date. The records are read from the default
branch head as the local refs hold it, never from the working tree.

**The window.** From the Run Retention cutoff to now. `--since` and `--until`
narrow it by delivery date, and neither widens it.

**The size of a Spec.** The Tasks of the Task Graph in the first parent of
the commit that added the Archive Record, without the QA gate and without any
Task numbered above it. A Spec whose plan that commit does not hold has no
size: it is shown, and it enters no group.

**The moment that picks a phase's entry.** The creation of the first
Implement Run this machine's Run Database holds for the Spec, or the delivery
date when it holds none.

**What the report starts.** Git, on local objects and refs. No fetch, no
forge request, no judge and no read of a harness transcript.

**What it cannot read.** A Run Database or a corrections record that cannot be
read ends the report with exit code 1, as `roundfix runs causes` ends. Any
other source it cannot read is named as unavailable with its reason, and the
report exits 0 with the rows it has.

## Considered Options

- **The Task Graph at the `source_revision`.** Rejected: it needs the pull
  request's head from the forge once the object is pruned, and it counts the
  Corrective Tasks of the delivery as size.
- **A Task count written into the Archive Record.** Not taken: ADR-0417 adds
  nothing to the Archive Record.
- **The Task count recorded with the signs at the merge.** Rejected: a Spec
  whose delivery this machine did not record would have no size, and its row
  would leave every Task count.
- **Every delivered Spec of the repository's history.** Rejected: the phase
  entries, the declarations and the Runs behind a row are gone after Run
  Retention.
- **Exit code 1 for any unreadable source.** Rejected: an unreadable Phase
  Record leaves the Corrections readable, and the PRD asks for the rows.

## Consequences

- A plan defect is counted over the Tasks the plan held, not over the Tasks
  the delivery ran.
- A Spec planned and delivered in one commit has no size on this report.
- A checkout whose default branch ref is old shows an older window; the
  report states the window it read.
- Any clone of the repository on this machine finds the same corrections
  record and phase entries, through the Repository Identity (ADR-0421). The
  Runs are found only from the checkout that ran them, so another clone picks
  each phase's entry by the delivery date.
