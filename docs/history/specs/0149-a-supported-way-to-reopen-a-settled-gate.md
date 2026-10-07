---
schema: roundfix/archive-record/v1
spec: 0149-a-supported-way-to-reopen-a-settled-gate
title: A supported way to reopen a settled gate
status: archived
created: "2026-09-19"
archived: "2026-09-19"
disposition: pass
source: docs/history/specs/0149-a-supported-way-to-reopen-a-settled-gate
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
delivery_commit: c0158ca7d8a26ec6dd6983ac53bd9829f3218cfb
---

# A supported way to reopen a settled gate

A Spec's terminal QA Task is a gate: it depends on every leaf of the Task Graph, and when it settles `completed` the recorded verdict describes the graph as it stood. If a Task below it later stops being completed — because a corrective Task is added, or because a settled Task is reopened — the gate's result no longer describes anything real, and the loader says so:
