---
status: accepted
created_at: 2026-09-29T00:00:00Z
updated_at: 2026-10-06T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The pre-PR Pull Request row never decides a qualifying partial

The authored QA gate is a Task of the Spec's own graph (ADR-0088, ADR-0091), and
it runs before any Pull Request exists. The Pull Request row, which every matrix
must carry, is therefore always recorded as
`blocked (environment: no open Pull Request)`. ADR-0080 lets an
environment-blocked row reach `pass` when the report carries equivalent
evidence. A qualifying declared `partial`, however, required
`rows_blocked_environment: 0`, so a Spec with a real Unreachable Acceptance
declaration could never qualify. Spec 0179 needed a boilerplate fourth
declaration and an extra QA Requirement to get past this, at the cost of two
extra QA Runs.

A row counts as the pre-PR Pull Request row when both of the following hold:

- its status is exactly `blocked (environment: no open Pull Request)`;
- its provenance names the Pull Request row.

That row no longer counts against a qualifying `partial`. Every other
environment-blocked row keeps today's rule. The report's frontmatter count is
unchanged, and the exception is derived from the Results rows the mechanical
stage already validates.

## Consequences

The exception covers the one row that no pre-PR gate can observe. Its cause is
the order ADR-0088 fixed, not the situation of a particular Run. Every other
environment block still has to be declared for a `partial` to qualify. A row
with the same status but another source in its provenance is an ordinary
environment block. A report written without a provenance column gains nothing
from the exception.

**One policy (2026-10-06).** Extended by ADR-0240: a partial whose only unmet row is the pre-PR Pull Request row qualifies without an Unreachable Acceptance declaration, and a provenance item that starts with the Pull Request row followed by a space, `:` or `(` names it. Every other part of this decision stands.
