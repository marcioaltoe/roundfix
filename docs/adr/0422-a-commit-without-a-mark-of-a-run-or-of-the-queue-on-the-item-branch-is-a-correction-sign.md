---
status: accepted
knowledge: true
created_at: 2026-11-17T19:30:00Z
updated_at: 2026-11-17T20:30:00Z
deprecated_at: null
superseded_by: null
---

# A commit without a mark of a Run or of the queue on the Item Branch is a Correction Sign

Roundfix tells a correction made by hand from its own commits by what each
commit carries, never by who made it. A commit on the first-parent line of an
Item Branch is a Correction Sign unless it carries a mark of a Run or of the
Delivery Queue, or merges the default branch. A fixed review finding is one
sign together with its fixing commit, and a Corrective Task is one sign
together with its Task commit.

## Context

ADR-0417 names three Correction Signs and leaves the commit to be defined.
The second measurement of Spec 0404 read every commit Roundfix makes on an
Item Branch at `74d0e3a9`, and the reading was repeated at `a8236a87`. A Task
commit and the QA Report commit carry a `Roundfix-Spec` trailer. The merge
that regenerates derived paths carries `Roundfix-Delivery: derived-merge`. The
archive commit and the Pull Request Record Commit carry no trailer, and the
queue accepts each by proof of its change. A carried-forward Task commit keeps
its trailers. No code reads a commit's author. Of the nine commits of pull
request #892, exactly the Supervisor's three carry none of these.

Two reviews of one Item Branch each name their first finding `F1`, and the
disposition ledger of a repository is never trimmed, so a finding is not named
by its identifier alone.

## Decision

**The commits read.** For a head `H` and the default branch head `D`, the
commits read are those of `git rev-list --first-parent D..H`. A commit that
reaches the branch only through a merge is not read.

**The commit that is a sign.** One of those commits is a Correction Sign
unless:

- it carries a `Roundfix-Spec` trailer, as a Task commit and the QA Report
  commit do;
- it carries the trailer `Roundfix-Delivery: derived-merge`;
- it is a merge whose other parents the default branch holds;
- the Delivery Queue's own proof accepts it as the archive commit of the Spec
  or as its Pull Request Record Commit;
- it is the fixing commit of a finding that is itself a sign.

**The finding that is a sign.** A fixed Review Finding Disposition whose
reviewed head is a commit of `D..H`. It is named by its finding identifier and
its reviewed head, `F<n>@<head>`.

**The Corrective Task that is a sign.** A Task of the Spec's Task Graph
numbered above its QA gate, named by its Task identifier.

**When the signs are read.** The Queue Owner reads the signs of the candidate
head each time an item reaches its merge and records them (ADR-0420);
`roundfix deliver status` reads the signs of the branch tip without recording.

## Considered Options

- **The commit's author or committer.** Rejected: a Run, the queue and the
  Supervisor commit under the same Git identity on one machine.
- **The commit's type, such as `fix`.** Rejected: a Task commit may be a
  `fix`, and a correction of the plan is a `docs` commit.
- **A required trailer on every correction commit.** Rejected by ADR-0417: a
  commit by hand runs no Roundfix command that could refuse it.
- **Every commit of `D..H`, not only the first-parent line.** Rejected: a
  merge of the default branch brings the squash commits of other Specs, which
  carry no mark.
- **The archive commit and the Pull Request Record Commit told by their
  subject.** Rejected: the queue itself accepts an archive commit an operator
  made, by its change.
- **A finding and its fixing commit as two signs.** Rejected: one cause would
  count as two `undeclared` Corrections in a repository that declares nothing.

## Consequences

- A repair made inside a merge of the default branch is no sign and is
  declared by hand, as ADR-0417 says.
- A branch merged into the Item Branch by hand is one sign, its merge commit.
- A plan change that records a new decision of the maintainer is no sign when
  it reaches the Item Branch through a merge of the default branch. Committed
  by hand on the Item Branch it is a sign, and it is declared
  `not-a-correction` with the decision as its reason (ADR-0417), which covers
  the sign and keeps it out of every defect count.
- A commit by hand that copies a `Roundfix-Spec` trailer is not seen.
