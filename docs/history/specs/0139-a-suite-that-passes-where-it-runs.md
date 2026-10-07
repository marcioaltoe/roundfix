---
schema: roundfix/archive-record/v1
spec: 0139-a-suite-that-passes-where-it-runs
title: A suite that passes where it runs
status: archived
created: "2026-09-15"
archived: "2026-09-16"
disposition: pass
source: docs/history/specs/0139-a-suite-that-passes-where-it-runs
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_05
qa_report: qa-report-2026-09-16.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0096
sources: []
regeneration: []
promoted: []
pull_request: "191"
delivery_commit: 368b92a90ef0da08b848429dcc54396880b49473
---

# A suite that passes where it runs

The repository Verification passes in CI and fails on a maintainer's macOS machine at the same tree. Spec 0138's QA gate failed on it, and so does every later Spec's gate on this machine, whatever the Spec changes. On the unchanged base, the same `make verify` failed in four packages no Spec had touched. Each failure has a proven cause outside the code being verified:
