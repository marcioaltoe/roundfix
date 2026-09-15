---
status: approved
granted: 2026-09-15
created: 2026-09-15
action: make the repository Verification deterministic where it runs and run the QA gate's repository Verification in the Daemon
consuming: 0139-a-suite-that-passes-where-it-runs
paths:
  - internal/daemon/task_engine.go
  - internal/daemon/task_engine_test.go
  - internal/speccheck/mechanical_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0139

On 2026-09-15 the maintainer was asked to approve this record with its exact
three paths and operations, and answered "Aprovar como proposto". That answer
approved the bounded scope above.

The same day the maintainer made two other decisions. The repository's own suite
gets fixed before Spec 0138's QA gate is retried. The Daemon, not the QA Agent,
runs the gate's repository Verification. This record lands in `main` before the
Task commit that consumes it.

## Why a governed path is unavoidable

Each of the three paths is governed because an earlier authorization bounded it
(ADR-0130), and each holds a defect this Spec has to repair where it lives:

- **`internal/daemon/task_engine.go`** holds the QA gate step. The repository
  Verification can run outside the QA Agent's sandbox only if the Daemon runs it
  there, before the Agent turn. Nothing else in the codebase owns that step.
- **`internal/daemon/task_engine_test.go`** holds the Task-cycle tests. Their
  waits use a fixed two-second deadline, which real work exceeds under load. It
  is also where the QA gate step's tests live, so the new step is tested there
  too.
- **`internal/speccheck/mechanical_test.go`** holds the historical audit
  subtest. Its expectation has been stale since the audit started reading a
  grant at its authorizing ancestor. Its skip guard also checks whether a Git
  object exists instead of whether it is reachable.

The ACPX fixture repair is in `internal/agent/acpx_runner_test.go`. No
authorization has ever bounded that file, so it needs no grant.

## Approved bounded mutation

- Run the configured repository Verification in the QA gate step before the
  Agent turn, outside the Agent sandbox, through the existing Verification
  machinery. Then:
  - record its command, verdict, exit status and diagnostics in the seeded QA
    Report;
  - tell the QA Agent that it already ran and must not be run again.
- Bind the Task-cycle test waits to the test's own deadline instead of a fixed
  wall-clock bound.
- Replace the stale historical audit expectation with a controlled fixture, and
  make the skip guard require reachability.

## Limits

- No action, operation or path beyond those above.
- No change to verdict rules, typed blocked causes, QA Report keys, Mechanical
  Refusal Codes, Verification classifications or the qa-gate skill.
- No existing assertion is weakened or deleted, except the stale historical
  expectation this Spec replaces.
- No timeout increase, retry, skip or `t.Parallel` removal as a repair.
- No Baseline module, linter, Verification configuration or `.roundfixrc.yml`
  edit, and no paid API use, release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
