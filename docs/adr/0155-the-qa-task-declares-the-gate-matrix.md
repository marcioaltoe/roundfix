---
status: accepted
created_at: 2026-09-14T00:00:00Z
updated_at: 2026-09-15T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The `qa` Task declares the gate's matrix

The QA gate used to build its own matrix from every promise, story, criterion and
Non-Goal, so two executors of the same Spec covered different sources. Spec 0119
needed six QA Reports, and only two of its five failed reports came from a
defect in its own work. The rest carried noise: rows blocked in cascade, an
aggregate row, and inputs widened after a row ran.

A numbered `qa` Task Requirement that starts with `MUST verify` or `MUST run` now
declares a coverage source. When such a Requirement verifies every Acceptance
Criterion of named Tasks, each criterion is a source. A `qa` Task with at least
one such Requirement has declared its matrix; a Spec without one keeps a bounded
default. Every matrix also covers the outside-evidence row, the Pull Request
row, the repository Verification, each PRD Unreachable Acceptance declaration
and, when declared, the frontend sweep. Each row names the sources it covers.
Every source must appear in some row, and no row may cover anything else. How
rows group sources is left to the executor.

Three alternatives were rejected.

- **One row per Requirement.** Three review rounds each found a new place where
  that rule needed interpretation. A rule on the covered set does not.
- **Judging whether Requirements "name what the gate verifies".** Two Agents can
  judge that differently.
- **Keeping gate-derived breadth.** It caught no defect that the narrower
  matrices of Specs 0134 to 0137 missed, and every rerun paid for it.

The accepted cost is that a promise left out of the `qa` Task goes unchecked.
That moves the weight onto authoring, which Spec 0129 strengthens.
