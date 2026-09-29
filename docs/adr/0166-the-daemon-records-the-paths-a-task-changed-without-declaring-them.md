---
status: accepted
created_at: 2026-09-29T00:00:00Z
updated_at: 2026-09-29T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The Daemon records the paths a Task changed without declaring them

A Task's `## Context` names the paths it plans to edit, and the authored QA gate
audits every changed file against those declarations. An Agent that needed a
file it did not plan, usually a new test file, could not declare it: the Agent
may edit only its Task's `## Result`, never its Context. The gate then refused a
completed, verified Task on the undeclared path. That cost two extra QA Runs
across Specs 0179 and 0180, and each needed a corrective edit that only added a
line to a Task's Context.

The Daemon now records those paths itself, when it commits a completed Task. It
takes the paths the Task commit stages and removes the Task file, every
`interface:` and `creates:` path in the Task's Context, and every Governed Path.
It writes what is left into a Daemon-owned `## Recorded paths` section at the end
of the Task file and lists the same paths in the commit event. A QA Task records
nothing. The QA scope audit counts a recorded path as declared, and names it, so
the undeclared change is disclosed rather than refused.

The Daemon already writes the Task file it commits: it settles the status there
(ADR-0014, ADR-0057). The record joins that write. It never edits the authored
Context, for two reasons. An appended `creates:` entry can collide with another
Task's declaration or pass the fifty-entry Context limit, and either one makes
the Task file unparseable. And Task Carry-Forward reads a Task's declared inputs
from the Task file before its settlement commit, so a record written at
settlement cannot change what carry-forward compares.

## Consequences

A Governed Path is never recorded. A Governed Path a Task changes stays under
its authorization record, exactly as before (ADR-0130). Declaring a path is
still how a Task plans its work and reserves the path against parallel Tasks.
Recording discloses a change after the fact, and it reserves nothing.
