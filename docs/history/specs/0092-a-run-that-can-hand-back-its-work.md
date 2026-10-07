---
schema: roundfix/archive-record/v1
spec: 0092-a-run-that-can-hand-back-its-work
title: A Run that can hand back its work
status: archived
created: "2026-08-09"
archived: "2026-08-11"
disposition: pass
source: docs/history/specs/0092-a-run-that-can-hand-back-its-work
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_07
qa_report: qa-report-2026-08-11-04.md
qa_verdict: pass
unproven: []
adrs: []
sources:
  - 2026-08-08-a-session-that-never-opened-is-a-selection-failure.md
  - 2026-08-09-a-failed-batch-erases-the-issues-it-resolved.md
  - 2026-08-09-a-stopped-run-discards-the-tasks-it-already-proved.md
regeneration: []
promoted: []
pull_request: "154"
delivery_commit: 8a33f8fdec67a82c74592ac30c1e0e91c8aa678d
---

# A Run that can hand back its work

A Run creates Tasks, Batches, Agent Sessions, Run Branches and worktrees. Each of those needs exactly one terminal disposition, and today several of them are classified by a surface that did not create them. Four measured consequences, all from the same shape.
