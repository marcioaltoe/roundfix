---
status: accepted
created_at: 2026-08-03T00:00:00Z
updated_at: 2026-08-03T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Spec consistency is checked by citation, never by inference

The Spec Consistency Check reports only what a Spec's own artifacts say about
each other. It compares citations, declarations, and cross-references that are
already written down; it never judges whether a decision is correct, whether an
ADR is topically relevant, or whether two paragraphs of prose disagree. That
boundary is what makes the check deterministic, hermetic, and fast enough to
run on every authoring change — the properties that let it replace a
twenty-minute gate cycle for this defect class.

The boundary was measured, not assumed. Only one of eighty-nine accepted ADRs
names a repository path, so an ADR-to-code-surface index cannot work here;
ADRs are written in domain vocabulary. What does work is the ADR citation
graph: an ADR that cites an ADR the Spec already lists is a candidate the Spec
should account for. Replaying Spec 0056, whose F-001 named ADR-0055 as the
omission, the depth-one closure over its cited set — ADR-0037, ADR-0039,
ADR-0049 — is exactly ADR-0040 and ADR-0055. Two candidates, one of them the
finding. That is a report an author can read, not noise.

Findings therefore carry two severities with different consequences. An
**error** states a contradiction both of whose sides the check located, with
file and line for each: a cited ADR missing from the Active ADR row, a
Constraint row citing a source path that does not exist, a tooling row citing
an authorization record that does not list the Spec, a PRD Core Feature no
Coverage Map entry covers. A **gap** states a candidate the check surfaced but
cannot settle — the relation-closure ADRs above. Errors fail the check; gaps
report and exit zero. A gap is dismissed by writing the reason into the Spec,
which is the artifact the next reader has. Nothing about this weakens or
duplicates the QA gate: undeclared prose contradictions remain QA's, and
ADR-0080 keeps sole ownership of QA verdict semantics.
