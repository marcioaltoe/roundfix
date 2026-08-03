---
task: task_07
spec: 0074-git-spawn-economy
status: failed
type: qa
complexity: high
---

# Task 07: Run the final QA gate

## Overview

The authored terminal gate — present in this graph from decomposition, as
the contract Spec 0072 shipped requires. The Daemon executes `qa-gate` as
this node once task_06 settles `completed` (covering every other node
transitively), producing the dated report, verdict, and typed blocked-row
counts per ADR-0080.

Two lessons from 0072's close apply operationally: the gate's
review-readiness row should be judged against the reviewer's decision on
the exact current head, and this Spec's tooling authority is `not
applicable` by design — there is no protected-tooling surface for the
governance audit to fail.

## Acceptance Criteria

- [ ] The gate runs only after task_06 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict
      and typed blocked-row counts.

## Verification

- `ls docs/specs/0074-git-spawn-economy/qa/ | grep -q "qa-report-"` —
  expected: exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0074-git-spawn-economy/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0080, ADR-0091.

## Result

The gate executed as the graph's terminal node with no invocation requesting
it — the first authored gate to run that way, one day after Spec 0072
shipped the contract. Its verdict was `fail`, and the report names why: two
base-versus-current binary comparisons were blocked by the sandboxed Go
cache, and the Pull Request row was blocked because none was open when the
gate ran. Both are environment and sequencing conditions, not defects in the
work.

Closed under `qa_override` per the maintainer's standing direction. The
substance the blocked rows describe was proven afterwards by other means:
PR #92 opened, passed the Verification gate, was reviewed with three
findings resolved, and merged approved. The suite passes at 3,125 tests.
