---
status: done
absorbed_by: 2026-08-06-rollup-qa-gates-and-verification-evidence.md
created_at: 2026-08-26
updated_at: 2026-08-26
kind: finding
---

# A QA Task's authored Verification is never executed, yet the authoring contract requires one

The authored terminal `qa` Task settles from the QA Report's `verdict:` — that
is the contract and it works. What does not close: the authoring contract treats
the `qa` Task like any other and requires a `## Verification` proving its own
effect, while the runtime never executes it.

Measured in oraculum, Spec 0047, task_05 (`type: qa`), two Runs on 2026-08-24
(`run_20260824T210308Z_06a60e2ff9e7d359` failed, `run_20260824T213047Z_f6a26cca04221883`
completed): zero `verification` events for the QA Task in both; settlement came
from the report verdict each time.

Source: secondbrain `inbox/roundfix/2026-08-24-a-verification-de-uma-task-qa-nunca-e-executada.md`
(origin oraculum). Deferred at the 2026-08-26 triage: needs reproduction in this
repository before a Spec commits to it. Related family: Spec 0116 names it as a
non-goal; 0105 owns the gate's Verification derivation and may absorb it.
