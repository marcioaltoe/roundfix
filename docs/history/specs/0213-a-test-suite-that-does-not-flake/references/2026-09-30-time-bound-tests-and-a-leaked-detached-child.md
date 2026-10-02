---
type: fix
status: promoted
created: 2026-09-30
spec: 0213-a-test-suite-that-does-not-flake
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

## Addendum — 2026-10-01 — three more flakes of the same family

The entry is extended rather than a new one minted, because these flakes share
its context: tests that pass on rerun and fail under load or outside events.

- `TestDoctorAcceptsAnAdapterAtTheFloor/codex` (`internal/cli`) failed in the
  CI Verification gate of PR #317, run 36903640400 attempt 1:
  `Doctor = cli.CheckResult{Name:"adapter", Status:"failed", Detail:"codex:
  effective Codex adapter command \"/tmp/TestDoctorAcceptsAnAdapterAtTheFloorcodex1616886882/001/fake-adapter\"
  did not prove required package lineage ...`. The test writes a shell script
  and executes it at once; the attempt passed on rerun.
- `TestAssetsSyncProvenanceAndPreMutationRefusals/*` (`internal/baseline`)
  failed in the QA gate of Spec 0203 with `copy fixture
  .../assets-sync-template.../target: open .git/objects/05: no such file or
  directory`. The operator classified it as an environment failure (a
  temporary fixture deleted) and retried.
- `TestTaskBudgetCancelsAStalledTaskOneAllowanceAfterTheLastSettlement`
  (`internal/daemon`) shares the 250 ms wall-clock shape of the Run Budget test
  named above.
