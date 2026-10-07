---
schema: roundfix/archive-record/v1
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
title: A Task settles on the facts its gate will check
status: archived
created: "2026-09-30"
archived: "2026-10-01"
disposition: pass
source: docs/history/specs/0190-a-task-settles-on-the-facts-its-gate-will-check
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0182
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "308"
delivery_commit: 64aff3f72dc7d72559b31a6f9213452feb02e7e5
---

# A Task settles on the facts its gate will check

A Task settles `completed` when its own declared Verification passes. The QA gate then refuses, before it spends an Agent turn, on three machine facts: the strict Spec Consistency Check, the configured repository Verification and the authorization audit of each Task commit. A completed Task can carry any of them, and the Run learns it only at the end, after the Agent Session that could have fixed it has closed.
