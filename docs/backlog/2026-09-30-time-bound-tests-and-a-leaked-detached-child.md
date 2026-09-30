---
type: fix
status: open
created: 2026-09-30
spec: null
reason: null
---

# Time-bound tests fail under load, and one test leaks its detached child

## Opportunity

Two tests assert wall-clock bounds and failed under load in the v0.22.0 cycle while passing in isolation:

- `TestTaskBudgetReasonNamesTheSettlementThatRenewedIt` (`internal/daemon`), a 200 ms Run Budget measured at 266 ms on the CI runner (PR #297);
- `TestBatchClosesOnCountLingerAndImmediate/linger_closes_a_quiet_batch` (`internal/store`), in the queue owner's repository gate for 0193. It passed 10 of 10 runs in isolation.

`TestRunImplementDetachSurvivesCallerProcessGroupKill` left its detached `cli.test implement --spec 0001-widget-flow --detach` child alive for five hours after the suite ended.

## Value

Each flake parks a delivery that was otherwise ready, and the orphan holds a process slot until someone kills it by PID.

## Shape

Replace the two wall-clock assertions with event-driven or fake-clock waits, and make the detach test reap its child with a cleanup that proves the process is gone. Evidence: `docs/history/findings/2026-09-30-the-v0-22-0-queue-needed-twelve-manual-interventions.md`, section 4.
