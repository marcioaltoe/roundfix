---
status: done
created_at: 2026-09-08
updated_at: 2026-09-08
kind: finding
spec: 0127-durable-unattended-spec-workflow
---

# Implement — the declared Run duration is not applied to execution (2026-09-08)

The project declares `budget.max_run_duration: 2h`, but the inspected Implement
path uses the value for presentation and a Run Window warning rather than a
Run deadline. Unattended delivery must not describe this setting as an enforced
bound until its execution path is corrected and tested.

## 1. Presentation and enforcement have different consumers

- Symptom / evidence: `internal/cli/implement.go` warns about cutoff overrun and renders the budget label, then passes the original context into `executeImplementCycle`. Production searches for `MaxRunDuration` and `BudgetEnabled` locate enforcement in `watch.Run`, not the Implement executor.
- Root cause: the existing review Watch budget is not propagated as an enforced deadline through Spec execution.
- Action / suggestion: Spec 0127 must enforce an approved Run deadline before its first unattended Spec delivery, propagate cancellation to the owned Agent/process tree, and preserve recovery evidence rather than declaring Clean on expiry. Run Window continues to bound new starts, not silently supply the missing execution deadline.

Evidence is static inspection on 2026-09-08. No two-hour experiment, paid Agent
invocation, Run or timeout reproduction was executed. The acceptance must use
a controlled clock and a real isolated process boundary rather than this
inspection as its final proof.

## Addendum — 2026-09-08 — Implementation owner

The maintainer selected this remaining intent for the implementation queue.
[0127-durable-unattended-spec-workflow](../_prd.md) is its primary owner.
The source moves once into that Spec's reference index; other Specs link this
owned copy. Its lifecycle status records adoption, not implementation or QA
completion. Governed mutation and execution still require the consuming grant.
