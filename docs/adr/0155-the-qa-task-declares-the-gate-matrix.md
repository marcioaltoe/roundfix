---
status: accepted
created_at: 2026-09-14T00:00:00Z
updated_at: 2026-09-14T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The `qa` Task declares the gate's matrix

The QA gate built its own matrix from every promise, story, criterion and
Non-Goal, so two executors of the same Spec planned different rows with
different inputs. Spec 0119 needed six QA Reports, and only two of its five
failed reports came from a defect in its own work. The rest carried noise: rows
blocked in cascade, an aggregate row, and inputs widened after a row ran.

Each numbered `qa` Task Requirement that starts with `MUST verify` or `MUST run`
is now exactly one matrix row. Every other Requirement constrains how the gate
runs. A `qa` Task with at least one such Requirement has declared its complete
matrix. Specs without one keep a bounded default. No declaration can waive the
outside-evidence row, the Pull Request row, the repository Verification, a PRD
Unreachable Acceptance declaration or the frontend sweep. A declared row that
names one of those sources covers it, so each is planned once.

A textual rule was chosen over judging whether Requirements "name what the gate
verifies", because two Agents can judge that differently. It is also the form
Specs 0134 to 0137 already wrote. Keeping gate-derived breadth was rejected: it
caught no defect those narrower matrices missed, and every rerun paid for it.
The accepted cost is that a promise left out of the `qa` Task goes unchecked,
which moves the weight onto authoring, the area Spec 0129 strengthens.
