---
status: accepted
created_at: 2026-09-17T00:00:00Z
updated_at: 2026-09-17T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A promise a Spec declares names a consuming Task

The Spec Consistency Check derives its coverage units from two sections: the
PRD's user stories and its Core Features. A Spec's Success Metrics and its
TechSpec's API Contracts are read by no rule. Measured over the twenty most
recent Specs, twelve PRDs carry no Success Metrics section at all, eight declare
metrics and three of those eight are named by no Task, while twelve TechSpecs
declare API Contracts that nothing traces.

That was survivable while the QA gate rebuilt its own matrix from every promise
it could find. ADR-0155 ended it: the gate now covers what the `qa` Task
declares. A promise no artifact traces is now unread at authoring, unmapped at
design and uncovered at the gate.

A Success Metric and an API Contract are therefore declared as numbered items,
or as an explicit `None.` with the reason none applies. Each numbered item is a
coverage unit with the obligations user stories and Core Features already carry:
it appears in the TechSpec Coverage Map, and some Task names it in References.
The check locates both sides and cites them, never inferring coverage from
resemblance, and each side is reported by the authoring stage that can
establish it.

Three alternatives were rejected.

- **Inferring coverage from Task prose.** A Task that mentions latency does not
  thereby own a latency metric. ADR-0093 keeps consistency a citation
  relationship.
- **Requiring every Spec to invent a metric.** Half the recent Specs are
  internal repairs with nothing to measure after shipping. Forcing a number
  produces ceremony, so an explicit `None.` with a reason satisfies the
  declaration.
- **Leaving it to independent review.** Review is configurable and may be
  explicitly declined, so it cannot carry an obligation every Spec must meet.

The accepted cost is that authoring now pays a declaration it used to skip in
silence, and a Spec with genuinely nothing to measure must still say so.
