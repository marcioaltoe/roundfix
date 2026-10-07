---
schema: roundfix/archive-record/v1
spec: 0144-a-run-stops-when-its-budget-is-spent
title: A Run stops when its budget is spent
status: archived
created: "2026-09-18"
archived: "2026-09-18"
disposition: pass
source: docs/history/specs/0144-a-run-stops-when-its-budget-is-spent
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_05
qa_report: qa-report-2026-09-18-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0158
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "209"
delivery_commit: 4eac7597e2928e2ccdfc16b1a2779387e87fedd1
---

# A Run stops when its budget is spent

Roundfix lets a repository set a maximum Run duration, and it enforces that maximum in exactly one place: the Round watch loop, which derives a deadline from the Run's start and settles the Run as `BudgetExceeded` when it passes.
