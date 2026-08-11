---
task: task_08
spec: 0080-cheap-detectors-run-before-the-gate
status: failed
type: qa
complexity: medium
---

# Task 08: Run the final QA gate

## Overview

The authored terminal gate, and the first one that runs against the machinery
this Spec builds. Its most valuable rows are the ones no Task can assert alone:
that a defect a machine can find is now found in seconds, that a re-gated round
costs less than a first round, and — above all — that carry-forward never
carried something it should have re-observed.

## Requirements

1. MUST run only after task_05 and task_07 settle.
2. MUST observe the mechanical stage finding a blocking fact and withholding
   the Agent Session, measuring wall-clock from gate start to reported
   failure, since seconds-not-a-round is the Spec's first Success Metric.
3. MUST observe a re-gated round after a corrective Task costing materially
   less than the first, and attribute the reduction to carried rows rather
   than to chance.
4. MUST attempt to defeat carry-forward rather than confirm it: exercise at
   least one row whose evidence moved and confirm it re-observed, and one row
   with a non-repository input and confirm it never carried.
5. MUST confirm no verdict regression across the fixture corpus: for every
   Spec exercised, the mechanical stage plus audit reaches the same or a
   stricter verdict than the audit alone.
6. MUST confirm every mechanical detector reports a citation — file, line, and
   fix — for each failure it raises, and records a skip naming its detector and
   missing artifact when its input is absent.
7. MUST confirm the local incremental tier stays within the clause's declared
   budget and that CI still runs the complete fresh gate.
8. MUST record Results as rows `R01` through `R05`, one per PRD Success Metric
   in order, and classify any finding by user impact with typed blocked-row
   counts per ADR-0080.

## Acceptance Criteria

- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] Rows R01–R05 cover the five PRD Success Metrics in order.
- [ ] The carry-forward defeat attempts are recorded with their outcomes.
- [ ] The time-to-blocking-verdict observation carries a measured number.

## Verification

- `output="$(ls docs/specs/0080-cheap-detectors-run-before-the-gate/qa/qa-report-*.md 2>/dev/null)"; [ -n "$output" ]`
  — expected: exit 0; a dated QA report exists.
- `report="$(ls docs/specs/0080-cheap-detectors-run-before-the-gate/qa/qa-report-*.md | tail -1)"; grep -qE '^verdict: (pass|partial|fail)' "$report"`
  — expected: exit 0; the report carries a machine-readable verdict.
- `report="$(ls docs/specs/0080-cheap-detectors-run-before-the-gate/qa/qa-report-*.md | tail -1)"; grep -qE '^rows_blocked_environment: [0-9]+' "$report" && grep -qE '^rows_blocked_finding: [0-9]+' "$report" && grep -qE '^rows_blocked_declared: [0-9]+' "$report"`
  — expected: exit 0; all three typed counts are recorded.
- `report="$(ls docs/specs/0080-cheap-detectors-run-before-the-gate/qa/qa-report-*.md | tail -1)"; status=0; for id in R01 R02 R03 R04 R05; do grep -q "| $id |" "$report" || { echo "missing $id"; status=1; }; done; exit $status`
  — expected: exit 0; all five Success Metric rows are present.

## References

- `_prd.md` → Success Metrics (all five); Goals; User Stories 1–6.
- `_techspec.md` → Testing Approach; Build Order 8.
- ADR-0080, ADR-0091, ADR-0096, ADR-0097.
