---
schema: roundfix/archive-record/v1
spec: 0140-a-spec-traces-the-promises-it-makes
title: A Spec traces the promises it makes
status: archived
created: "2026-09-17"
archived: "2026-09-17"
disposition: pass
source: docs/history/specs/0140-a-spec-traces-the-promises-it-makes
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_05
qa_report: qa-report-2026-09-17.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0156
  - ADR-0093
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "199"
delivery_commit: 0ad1ced369bc1fc1cd094eab8e494e6862117d87
---

# A Spec traces the promises it makes

A Spec makes two kinds of promise nothing traces. It promises a measurable outcome after shipping, in its Success Metrics, and it promises a public interface, in its TechSpec's API Contracts. The Spec Consistency Check reads neither: its coverage units come from the PRD's user stories and Core Features alone. So a Spec can declare a metric no Task will ever settle, or drop the section entirely, and every stage stays silent.
