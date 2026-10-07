---
schema: roundfix/archive-record/v1
spec: 0009-parallel-scheduling
title: Parallel Scheduling
status: archived
created: "2026-07-05"
archived: ""
disposition: pass
source: docs/history/specs/0009-parallel-scheduling
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-06.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0025
  - ADR-0026
sources: []
regeneration: []
promoted: []
---

# Parallel Scheduling

Spec Runs execute one Task at a time even when the Task Graph declares whole Waves of independent work — the dogfood cycles spent hours in sequential waits that the graph itself said were unnecessary. The 0008 worktree isolation was built precisely to host concurrency; this Spec adds the layer on top: a ready-set scheduler that runs independent Tasks simultaneously, each in its own Task Worktree, integrating settled work back onto the Run Branch through a serialized queue. It also delivers the worktree location configuration and the empty-debris cleanup the last rounds asked for.
