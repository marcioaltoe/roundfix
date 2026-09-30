---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The Daemon records what a QA row observed and hands a failed pass to the next

ADR-0097 lets a QA row carry forward when its declared repository inputs are
unmoved and the report that established it recorded an evidence snapshot. No
component wrote that snapshot, and the establishing head also had to be an
ancestor of the current head. A failed gate's Run is never integrated, and Task
Carry-Forward re-commits its Tasks with new identities, so the next pass never
held the prior report nor descended from its head. Spec 0179 ran seven QA
passes whose rows declared repository inputs, and not one row was carried.

The Daemon now owns the evidence snapshot. When a QA pass closes, before the QA
Report commit, it records under `evidence_snapshots` the path and SHA-256
digest of every file each qualifying passing row declared, read from the
audited head. It replaces any value the Agent wrote. Before the next pass
builds its matrix, the Daemon also brings in the newest QA Report commit of the
same Spec that the Run's history does not contain. It copies that commit's
report and evidence byte for byte, so a failed pass on an unintegrated Run
Branch is the next pass's previous report.

The establishing head no longer has to be an ancestor. It is proven by either
of two facts: it is an ancestor of the current head, or a Daemon QA Report
commit whose first parent is that head holds the establishing report with the
same bytes. The inputs are then compared by content, as Task Carry-Forward
compares a Task's declared inputs and as a verifying trace compares hashes.
A head proven by neither fact carries nothing.

## Consequences

This refines ADR-0097 without superseding it. A row still carries only when it
passed, declared a non-empty list of repository inputs, and every input is
byte-identical and outside the changed-path set. A report with no recorded
snapshot still carries nothing. Speed still never buys a stale pass.

Imported reports join the Spec's QA directory and ride in the next QA Report
commit. The Spec keeps the history of every pass it needed, and a Run Branch
whose report reached the target is superseded under the existing rule.
