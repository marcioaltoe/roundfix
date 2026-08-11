---
task: task_09
spec: 0081-a-journal-cheap-to-write-and-keep
status: completed
type: qa
complexity: medium
---

# Task 09: Run the final QA gate

## Overview

The authored terminal gate. Its most valuable rows are the ones no Task could
assert alone: that Run start no longer tracks journal size, that parallel Runs
survive without the raised timeout carrying them, and that every consumer of
the journal behaves exactly as it did before — because the failure this Spec
could most easily cause is a silent loss of the only durable copy of a Run's
agent output.

## Requirements

1. MUST run only after task_06 and task_08 settle.
2. MUST observe Run-start wall-clock against both an empty and a
   full-retention journal, since size-independence is the first Success
   Metric and no single Task can measure it end to end.
3. MUST observe the parallel-Run scenario at the pre-change `busy_timeout`
   and confirm zero `SQLITE_BUSY`, rather than accepting a reduced count.
4. MUST verify event-write cost per thousand events against the recorded
   baseline and confirm the improvement is attributable to batching rather
   than to fewer events being produced.
5. MUST audit consumer non-regression by replaying a recorded journal through
   `events`, attach, cockpit rendering, reconcile replay detection, and `gc`,
   and confirm identical behaviour.
6. MUST confirm what the retention decision concluded and that the repository
   matches it: if no payload is shed, ADR-0008 is untouched and no operator
   documentation changed; if payloads are shed, the amendment exists and the
   lifecycle guide says so plainly.
7. MUST confirm bytes per Run before and after, with the retention shape
   stated.
8. MUST record Results as rows `R01` through `R06`, one per PRD Success Metric
   in order, and classify any finding by user impact with typed blocked-row
   counts per ADR-0080.

## Acceptance Criteria

- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] Rows R01–R06 cover the six PRD Success Metrics in order.
- [ ] The Run-start and parallel-Run observations carry measured numbers.
- [ ] The consumer replay is recorded with its evidence path.

## Verification

- `output="$(ls docs/specs/0081-a-journal-cheap-to-write-and-keep/qa/qa-report-*.md 2>/dev/null)"; [ -n "$output" ]`
  — expected: exit 0; a dated QA report exists.
- `report="$(ls docs/specs/0081-a-journal-cheap-to-write-and-keep/qa/qa-report-*.md | tail -1)"; grep -qE '^verdict: (pass|partial|fail)' "$report"`
  — expected: exit 0; the report carries a machine-readable verdict.
- `report="$(ls docs/specs/0081-a-journal-cheap-to-write-and-keep/qa/qa-report-*.md | tail -1)"; grep -qE '^rows_blocked_environment: [0-9]+' "$report" && grep -qE '^rows_blocked_finding: [0-9]+' "$report" && grep -qE '^rows_blocked_declared: [0-9]+' "$report"`
  — expected: exit 0; all three typed counts are recorded.
- `report="$(ls docs/specs/0081-a-journal-cheap-to-write-and-keep/qa/qa-report-*.md | tail -1)"; status=0; for id in R01 R02 R03 R04 R05 R06; do grep -q "| $id |" "$report" || { echo "missing $id"; status=1; }; done; exit $status`
  — expected: exit 0; all six Success Metric rows are present.

## References

- `_prd.md` → Success Metrics (all six); Goals; User Stories 1–6.
- `_techspec.md` → Testing Approach; Build Order 9.
- ADR-0080, ADR-0091, ADR-0098, ADR-0008.
