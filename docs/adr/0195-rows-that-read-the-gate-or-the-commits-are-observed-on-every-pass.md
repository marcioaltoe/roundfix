---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Rows that read the gate or the commits are observed on every pass

A correction always adds a Task commit and always changes the tree the
repository Verification and the Pull Request controls see. A row whose truth
depends on those facts cannot be unmoved, whatever paths it declares. Such a
row is therefore never carried:

- a row whose provenance names the repository Verification;
- a row whose provenance names the Pull Request row;
- a row that declares a `commit_range` input, the kind a row declares when it
  reads Task commits, their authorization or their changed-file scope.

`commit_range` joins ADR-0097's input vocabulary as a kind that is never
carriable, like `live_service`.

Every row of the previous report now receives a recorded carry disposition in
the next pass's seeded report: carried, or re-run with the reason. The
reasons form a closed list, such as `input moved: <path>` or
`always observed: repository Verification`. A carried row keeps the provenance
of the row that established it, so the matrix still names every coverage
source it covers.

## Consequences

The rule rests on two things. Provenance names are the fixed identifiers the
qa-gate skill already prescribes. An audit row is recognized only through its
`commit_range` declaration. A row that audits commits but declares only
repository paths is misdeclared, and the Daemon cannot detect that. ADR-0097
already relies on declarations in the same way.

Deciding by provenance was preferred to letting the Agent choose which rows
to repeat. It is mechanical and reviewable, and it keeps the cheap
Daemon-owned facts fresh. Carrying only the pass rows that read nothing but
repository content was accepted. Recording a reason for every re-run row gives
the Agent its work list and gives the operator the reason a row was not
carried.
