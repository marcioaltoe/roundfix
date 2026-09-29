---
status: accepted
created_at: 2026-09-28T00:00:00Z
updated_at: 2026-09-28T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An Implement Run Budget renews at each Task settlement

An Implement Run's allowance renews at each Task settlement because the Run
Budget must stop stalled work without stopping a serial graph that continues to
make progress. The renewed allowance covers the next Task or QA gate and the
post-cycle integration, push, and cleanup work; a settlement observed after the
current deadline remains settled but does not renew the allowance.

## Evidence

On 2026-09-28, Runs `run_20260928T111930Z_13d79b0227e06da5` for Spec 0176 and
`run_20260928T125114Z_30db2f43ddd8a501` for Spec 0173 reached
`BudgetExceeded` after useful serial progress under a two-hour maximum. Run
`run_20260928T154312Z_63679c6540b34603` for Spec 0175 reached the same outcome
and was then misclassified by its delivery owner. These three measured Runs
show that a start-based allowance can stop progressing graphs at an arbitrary
Task boundary.

## Considered options

A critical-path budget was rejected because it would grant the full scaled
allowance to a stall at the first Task. Holding the QA gate back was rejected
because it would end a progressing Run early and still require another Run to
finish the Spec.

## Consequences

The Daemon owns one renewing deadline, cancels all Agent Sessions it owns when
that deadline passes, and reports `BudgetExceeded` with the elapsed time since
the renewal point. This preserves the cancellation and recovery contract in
[ADR-0158](0158-a-budget-expired-run-settles-budget-exceeded-and-stays-recoverable.md).
The rendered configuration and user guides state the renewal beside
`budget.max_run_duration`, as required by
[ADR-0137](0137-the-run-budget-is-explained-where-it-is-configured.md).

