---
schema: roundfix/archive-record/v1
spec: 0142-an-absent-ref-is-named-absent
title: An absent ref is named absent
status: archived
created: "2026-09-17"
archived: "2026-09-17"
disposition: pass
source: docs/history/specs/0142-an-absent-ref-is-named-absent
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_03
qa_report: qa-report-2026-09-17.md
qa_verdict: pass
unproven: []
adrs: []
sources:
  - 2026-09-17-reconcile-reports-an-absent-ref-as-ambiguous.md
regeneration: []
promoted: []
pull_request: "203"
delivery_commit: 640bd79a5fd8bc760a380b5c4c5bf09a27f08be3
---

# An absent ref is named absent

Reconciliation tells a maintainer why it cannot release a Run's debris. On this machine it told them the wrong thing: fourteen candidates were refused with `resolve local branch "ma/0118-a-task-proved-once-does-not-run-twice": short ref is ambiguous`, for a branch that does not exist at all. Git lists no ref by that name.
