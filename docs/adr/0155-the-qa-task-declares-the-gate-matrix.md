---
status: accepted
created_at: 2026-09-14T00:00:00Z
updated_at: 2026-09-14T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The `qa` Task declares the gate's matrix

The QA gate used to build its own matrix from every promise, story, criterion and
Non-Goal. Two executors of the same Spec therefore planned different rows with
different inputs. Spec 0119 needed six QA Reports, and only two of its five
failed reports came from a defect in its own work. The rest carried noise: rows
blocked in cascade, an aggregate row, a repeated audit and inputs widened after a
row ran.

The `qa` Task's Requirements, when they name what the gate verifies, are now the
complete matrix. Specs without that declaration keep a bounded default. No
declaration can waive the outside-evidence row, the Pull Request row, the
repository Verification or the frontend sweep. The gate adds no rows of its own
beyond these.

Keeping gate-derived breadth was rejected. It did not catch any defect that the
narrower matrices of Specs 0134 to 0137 missed, and every rerun paid for it. The
accepted cost is that a promise left out of the `qa` Task goes unchecked. That
moves the weight onto authoring, which Spec 0129 strengthens.
