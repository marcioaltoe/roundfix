---
schema: roundfix/archive-record/v1
spec: 0150-a-reopen-that-cannot-be-raced
title: A reopen that cannot be raced
status: archived
created: "2026-09-19"
archived: "2026-09-19"
disposition: pass
source: docs/history/specs/0150-a-reopen-that-cannot-be-raced
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_04
qa_report: qa-report-2026-09-19-02.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
delivery_commit: c23f089a22a1f3af0783407f63b70bb07c78a28e
---

# A reopen that cannot be raced

Spec 0149 shipped `roundfix reopen` with two defects its third review round found and its corrective ceiling would not let it fix. Both are carried here rather than dropped.
