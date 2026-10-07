---
schema: roundfix/archive-record/v1
spec: 0147-a-planner-that-reads-both-tag-spellings
title: A planner that reads both tag spellings
status: archived
created: "2026-09-18"
archived: "2026-09-18"
disposition: pass
source: docs/history/specs/0147-a-planner-that-reads-both-tag-spellings
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_05
qa_report: qa-report-2026-09-18.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "222"
delivery_commit: debd76c85a69214c87943a0a8e6f07b5d3a59ee0
---

# A planner that reads both tag spellings

The Release Plan Command reads a repository's stable tags to decide what the next version is. It accepts exactly one spelling: a tag must begin with `v`, or it is rejected as malformed with the advice to "use a tag like v1.2.3".
