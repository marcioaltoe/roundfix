---
task: task_09
spec: 0064-spec-artifact-consistency-gate
status: failed
type: qa
complexity: high
---

# Task 09: Run the final QA gate

## Overview

The authored terminal gate, present in this graph from decomposition as the
contract Spec 0072 shipped requires. The Daemon executes `qa-gate` as this node
once task_08 settles `completed`, which transitively covers every other node,
producing the dated report, verdict, and typed blocked-row counts per ADR-0080.

This Spec's gate has an unusual advantage: the artifact-governance rows it
would normally have to audit by hand are the very thing the Spec built. The
gate should run `spec check` against this Spec and treat a clean report as the
governance evidence, then judge whether the check's own artifacts agree.

## Requirements

1. MUST run only after task_08 settles `completed`.
2. MUST exercise the command as a user does — against this Spec, against every
   active Spec, in both formats, and through the `make verify` gate — rather
   than by reading tests.
3. MUST verify the four replay fixtures report the codes their QA reports
   describe, since that is the Spec's first Success Metric.
4. MUST record the archived-corpus gap count as the observable false-positive
   rate the third Success Metric asks for.
5. MUST classify a finding by user impact and record typed blocked-row counts.

## Acceptance Criteria

- [ ] The gate runs only after task_08 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records the exit code observed for a clean Spec, a Spec with
      an error, and a Spec whose only finding is a gap.
- [ ] The report records the measured sweep duration against the Spec's
      sub-second budget.

## Verification

- `ls docs/specs/0064-spec-artifact-consistency-gate/qa/ | grep -q "qa-report-"`
  — expected: exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0064-spec-artifact-consistency-gate/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0080, ADR-0091, ADR-0093.
