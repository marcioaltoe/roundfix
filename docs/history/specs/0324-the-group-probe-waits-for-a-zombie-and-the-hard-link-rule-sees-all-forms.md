---
schema: roundfix/archive-record/v1
spec: 0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms
title: The group probe waits for a zombie, and the Hard Link rule sees all forms
status: archived
created: "2026-10-26"
archived: "2026-10-27"
disposition: partial
source: docs/specs/0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms
source_revision: ba32f39db32e218cc973bf9e2ec777ae545a4cce
qa_task: task_03
qa_report: qa-report-2026-10-26.md
qa_verdict: partial
unproven:
  - the operator's reading of the queue records and Run logs of the Specs delivered after this one, recorded session evidence.
adrs:
  - ADR-0324
sources:
  - 2026-10-26-the-unreaped-group-test-fails-under-the-full-suite-on-macos.md
  - 2026-10-26-the-hard-link-rule-misses-a-dot-import-and-a-function-value.md
regeneration: []
promoted: []
pull_request: "699"
---

# The group probe waits for a zombie, and the Hard Link rule sees all forms

Every runnable criterion passed, including 48,000 loaded iterations with zero failures and all six sabotage probes. The Delivery Base reproduced 167 EPERM failures. This is a qualifying declared partial: later-Run acceptance remains unproven, XNU network lookup is prohibited by authorization, and no Pull Request exists. The offline plan identifies the two adopted Backlog references as candidates; the existing audit reference already retains the reusable mechanism, so no promotion or archive write was performed.
