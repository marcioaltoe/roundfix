---
task: task_07
spec: 0079-one-door-for-fleet-knowledge
status: failed
type: qa
complexity: medium
---

# Task 07: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
task_04 and task_06 — the graph's non-QA leaves — settle. This Spec's most
valuable rows are the ones a task could not assert alone: the bounded active
set as an observed number, the door's durability and non-interference as
audited behavior, and the obligations binding only on a proven pilot.

## Requirements

1. MUST run only after task_04 and task_06 settle.
2. MUST observe the swept findings directory: every active finding carries a
   valid lifecycle, the rollup relations hold bidirectionally, and the
   active count is at most fifteen — this is where the number is measured,
   per the TechSpec's relations-not-counts rule for tasks.
3. MUST observe the three checks failing their red fixtures and skipping on
   an absent artifact class, not only passing green.
4. MUST audit the pilot report: both destinations exercised, durability
   under the one-minute target or the miss explained, zero Run interference
   recorded, and the qmd advisory check result present.
5. MUST observe the obligating clauses and inbox-first skill steps landed
   only with the pilot report present — permission before obligation.
6. MUST record the Results as rows `R01` through `R07`, one per PRD Success
   Metric in order, and classify any finding by user impact with typed
   blocked-row counts per ADR-0080; a metric only measurable in fleet
   operation (pending-beyond-fourteen-days) is an environment-blocked row
   with the equivalent evidence named.

## Acceptance Criteria

- [ ] A dated QA report lands under `qa/` with a machine-readable verdict
      and typed blocked-row counts.
- [ ] Rows R01–R07 cover the seven PRD Success Metrics in order.
- [ ] The active-findings count and lifecycle coverage are recorded as
      observations with evidence.

## Verification

- `output="$(ls docs/specs/0079-one-door-for-fleet-knowledge/qa/qa-report-*.md 2>/dev/null)"; [ -n "$output" ]`
  — expected: exit 0; a dated QA report exists.
- `report="$(ls docs/specs/0079-one-door-for-fleet-knowledge/qa/qa-report-*.md | tail -1)"; grep -qE '^verdict: (pass|partial|fail)' "$report"`
  — expected: exit 0; the report carries a machine-readable verdict.
- `report="$(ls docs/specs/0079-one-door-for-fleet-knowledge/qa/qa-report-*.md | tail -1)"; grep -qE '^rows_blocked_environment: [0-9]+' "$report" && grep -qE '^rows_blocked_finding: [0-9]+' "$report"`
  — expected: exit 0; typed blocked-row counts are recorded.
- `report="$(ls docs/specs/0079-one-door-for-fleet-knowledge/qa/qa-report-*.md | tail -1)"; status=0; for id in R01 R02 R03 R04 R05 R06 R07; do grep -q "| $id |" "$report" || { echo "missing $id"; status=1; }; done; exit $status`
  — expected: exit 0; all seven Success Metric rows are present.

## References

- `_prd.md` → Success Metrics (all seven); Goals; User Stories 1–6.
- `_techspec.md` → Testing Approach; Build Order 6.
- ADR-0080, ADR-0091, ADR-0095.
