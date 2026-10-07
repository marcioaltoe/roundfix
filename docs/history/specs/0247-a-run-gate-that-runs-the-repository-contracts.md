---
schema: roundfix/archive-record/v1
spec: 0247-a-run-gate-that-runs-the-repository-contracts
title: A Run gate that runs the repository contracts
status: archived
created: "2026-10-07"
archived: "2026-10-07"
disposition: pass
source: docs/specs/0247-a-run-gate-that-runs-the-repository-contracts
source_revision: b1020518f761ceb60ae67f9e6802ea3d11962cd7
qa_task: task_04
qa_report: qa-report-2026-10-07-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0252
sources: []
regeneration: []
promoted: []
---

# A Run gate that runs the repository contracts

A Run's Verification and its QA gate use `make verify-changed`, the `defaults.verification` of `.roundfixrc.yml`. That target compiles no test file under the `docscontract` or `repocontract` build tags. Only `make verify-docs` and `repo-test` run those tests, and in practice that means GitHub CI. On 2026-10-07 Spec 0245's Run and QA passed, and then CI failed twice on checks the Run never ran:
